package timing

import (
	"time"
	"log"
	"github.com/google/uuid"
	"os"
)

//compatible old name
type OnceCron struct {
	*TaskScheduler
}

//only exec cron timer cron
type TaskScheduler struct {
	tasks  []TaskInterface
	swap   []TaskInterface
	add    chan TaskInterface
	remove chan string
	stop   chan struct{}
	Logger TaskLogInterface
	lock	bool
}


type Lock interface {
	Lock()
	UnLock()
}

//return old name with OnceCron
func NewCron() *OnceCron {
	return &OnceCron{
		TaskScheduler:NewScheduler(),
	}
}

//return a Controller Scheduler
func NewScheduler() *TaskScheduler {
	return &TaskScheduler{
		tasks:  make([]TaskInterface, 0),
		swap:   make([]TaskInterface, 0),
		add:    make(chan TaskInterface),
		stop:   make(chan struct{}),
		remove: make(chan string),
		Logger: log.New(os.Stdout, "[Control]: ", log.Ldate|log.Ltime|log.Lshortfile),
		lock:	false,
	}
}

//add spacing time job to list with number
func (scheduler *TaskScheduler) AddFuncSpaceNumber(spaceTime int64, number int, f func()) {
	task := getTaskWithFuncSpacingNumber(spaceTime, number, f)
	scheduler.AddTask(task)
}
//add spacing time job to list with endTime
func (scheduler *TaskScheduler) AddFuncSpace(spaceTime int64, endTime int64, f func()) {
	task := getTaskWithFuncSpacing(spaceTime, endTime, f)
	scheduler.AddTask(task)
}

//add func to list
func (scheduler *TaskScheduler) AddFunc(unixTime int64, f func()) {
	task := getTaskWithFunc(unixTime, f)
	scheduler.AddTask(task)
}

func (scheduler *TaskScheduler) AddTaskInterface(task TaskInterface) {
    scheduler.addTask(task)
}
//add a task to list
func (scheduler *TaskScheduler) AddTask(task *Task) string {
	if task.RunTime < 100000000000 {
		task.RunTime = task.RunTime * int64(time.Second)
	}
	if task.RunTime < time.Now().UnixNano() {
		//延遲1秒
		task.RunTime = time.Now().UnixNano() + int64(time.Second)
	}

	//if task.Spacing > 0 {
	//		task.RunTime = time.Now().UnixNano() + int64(time.Second)
	//}else{
	//	scheduler.Logger.Println("error too add task! Runtime error")
	//	return ""
	//}


	if task.Uuid == "" {
		task.Uuid = uuid.New().String()
	}
	return scheduler.addTask(task)
}

//if lock add to swap
func (scheduler *TaskScheduler) addTask(task TaskInterface) string  {
	if scheduler.lock {
		scheduler.swap = append(scheduler.swap, task)
		scheduler.add <- task
	} else{
		scheduler.tasks = append(scheduler.tasks, task)
		scheduler.add <- task
	}

	return task.GetUuid()
}
//new export
func (scheduler *TaskScheduler) ExportInterface() []TaskInterface {
    return scheduler.tasks
}
//compatible old export tasks
func (scheduler *TaskScheduler) Export() []*Task {
    task := make([]*Task,0)
    for _,v := range scheduler.tasks {
        task = append(task, v.(*Task))
    }
	return task
}

//stop task with uuid
func (scheduler *TaskScheduler) StopOnce(uuidStr string) {
	scheduler.remove <- uuidStr
}

//run Cron
func (scheduler *TaskScheduler) Start() {
	//初始化的時候加入一個一年的長定時器,間隔1小時執行一次
	task := getTaskWithFuncSpacing(3600, time.Now().Add(time.Hour * 24 * 365).UnixNano(), func() {
		log.Println("It's a Hour timer!")
	})
	scheduler.tasks = append(scheduler.tasks, task)
	go scheduler.run()
}

//stop all
func (scheduler *TaskScheduler) Stop() {
	scheduler.stop <- struct{}{}
}

//run task list
//if is empty, run a year timer task
func (scheduler *TaskScheduler) run() {

	for {

		now := time.Now()
		task, key := scheduler.GetTask()
		//0 阻断
		if key < 0 {
			continue
		}
		runTime := task.GetRunTime()
		i64 := runTime - now.UnixNano()

		var d time.Duration
		if i64 < 0 {
			scheduler.tasks[key].SetRuntime(now.UnixNano())
			if task != nil {
				go task.RunJob()
			}
			scheduler.doAndReset(key)
			continue
		} else {
			sec := runTime / int64(time.Second)
			nsec := runTime % int64(time.Second)

			d = time.Unix(sec, nsec).Sub(now)
		}

		timer := time.NewTimer(d)

		//catch a chan and do something
		for {
			select {
			//if time has expired do task and shift key if is task list
			case <-timer.C:
				scheduler.doAndReset(key)
				if task != nil {
					//fmt.Println(scheduler.tasks[key])
					go task.RunJob()
					timer.Stop()
				}

				//if add task
			case <-scheduler.add:
				timer.Stop()
				// remove task with remove uuid
			case uuidstr := <-scheduler.remove:
				scheduler.removeTask(uuidstr)
				timer.Stop()
				//if get a stop single exit
			case <-scheduler.stop:
				timer.Stop()
				return
			}

			break
		}
	}
}

//return a task and key In task list
func (scheduler *TaskScheduler) GetTask() (task TaskGetInterface, tempKey int) {
	scheduler.Lock()
	defer scheduler.UnLock()

	if len(scheduler.tasks) < 1 {
		return nil,-1
	}
	min := scheduler.tasks[0].GetRunTime()
	tempKey = 0

	for key, task := range scheduler.tasks {
        tTime := task.GetRunTime()
		if min <= tTime {
			continue
		}
		if min > tTime {
			tempKey = key

			min = tTime
			continue
		}
	}

	task = scheduler.tasks[tempKey]

	return task, tempKey
}

//if add a new task and runtime < now task runtime
// stop now timer and again
func (scheduler *TaskScheduler) doAndReset(key int) {
	scheduler.Lock()
	defer scheduler.UnLock()
	//null pointer
	if key < len(scheduler.tasks) {

		nowTask := scheduler.tasks[key]
		scheduler.tasks = append(scheduler.tasks[:key], scheduler.tasks[key+1:]...)

		if nowTask.GetSpacing() > 0 {
            tTime := nowTask.GetRunTime()
			nowTask.SetRuntime(nowTask.GetSpacing() * int64(time.Second) + tTime)
			number := nowTask.GetRunNumber()
			if number > 1 {
				nowTask.SetRunNumber(number - 1)
				scheduler.tasks = append(scheduler.tasks, nowTask)
			} else if nowTask.GetEndTime() >= tTime {
				scheduler.tasks = append(scheduler.tasks, nowTask)
			}
		}

	}
}


//remove task by uuid
func (scheduler *TaskScheduler) removeTask(uuidStr string) {
	scheduler.Lock()
	defer scheduler.UnLock()
	for key, task := range scheduler.tasks {
		if task.GetUuid() == uuidStr {
			scheduler.tasks = append(scheduler.tasks[:key], scheduler.tasks[key+1:]...)
			break
		}
	}
}

//lock task []
func (scheduler *TaskScheduler) Lock() {
	scheduler.lock = true
}

//unlock task []
func (scheduler *TaskScheduler) UnLock() {
	scheduler.lock = false
	if len(scheduler.swap) > 0 {
		for _, task := range scheduler.swap {
			scheduler.tasks = append(scheduler.tasks, task)
		}
		scheduler.swap = make([]TaskInterface, 0)
	}

}
