package timing

import (
	"time"
	"log"
	"github.com/google/uuid"
	"os"
	"fmt"
)

//compatible old name
type OnceCron struct {
	*TaskScheduler
}

//only exec cron timer cron
type TaskScheduler struct {
	tasks  []*Task
	swap   []*Task
	add    chan *Task
	remove chan string
	stop   chan struct{}
	Logger *log.Logger
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
		tasks:  make([]*Task, 0),
		swap:   make([]*Task, 0),
		add:    make(chan *Task),
		stop:   make(chan struct{}),
		remove: make(chan string),
		Logger: log.New(os.Stdout, "[Control]: ", log.Ldate|log.Ltime|log.Lshortfile),
		lock:	false,
	}
}

//add spacing time job to list with number
func (scheduler *TaskScheduler) AddFuncSpaceNumber(spaceTime int64, number int, f func()) {
	task := getTaskWithFuncSpacingNumber(spaceTime, number, f)
	scheduler.addTask(task)
}
//add spacing time job to list with endTime
func (scheduler *TaskScheduler) AddFuncSpace(spaceTime int64, endTime int64, f func()) {
	task := getTaskWithFuncSpacing(spaceTime, endTime, f)
	scheduler.addTask(task)
}

//add func to list
func (scheduler *TaskScheduler) AddFunc(unixTime int64, f func()) {
	task := getTaskWithFunc(unixTime, f)
	scheduler.addTask(task)
}

//add a task to list
func (scheduler *TaskScheduler) AddTask(task *Task) string {
	if task.RunTime != 0 {
		if task.RunTime < 10000000000 {
			task.RunTime = task.RunTime * int64(time.Second)
		}
		if task.RunTime < time.Now().UnixNano() {
			//延遲1秒
			task.RunTime = time.Now().UnixNano() + int64(time.Second)
		}
	} else {
		if task.Spacing > 0 {
			task.RunTime = time.Now().UnixNano() + task.Spacing * int64(time.Second)
		}else{
			scheduler.Logger.Println("error too add task! Runtime error")
			return ""
		}
	}

	if task.Uuid == "" {
		task.Uuid = uuid.New().String()
	}
	return scheduler.addTask(task)
}

//if lock add to swap
func (scheduler *TaskScheduler) addTask(task *Task) string  {
	if scheduler.lock {
		fmt.Println("locked")
		scheduler.swap = append(scheduler.swap, task)
		scheduler.add <- task
	} else{
		scheduler.tasks = append(scheduler.tasks, task)
		scheduler.add <- task
	}

	return task.Uuid
}

//export tasks
func (scheduler *TaskScheduler) Export() []*Task {
	return scheduler.tasks
}

//stop tasks
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

//run Cron
func (scheduler *TaskScheduler) Stop() {
	scheduler.stop <- struct{}{}
}

//run task list
//if is empty, run a year timer task
func (scheduler *TaskScheduler) run() {

	for {

		now := time.Now()
		task, key := scheduler.GetTask()
		i64 := task.RunTime - now.UnixNano()

		var d time.Duration
		if i64 < 0 {
			scheduler.tasks[key].RunTime = now.UnixNano()
			if task != nil {
				go task.Job.Run()
			}
			scheduler.doAndReset(key)
			continue
		} else {
			sec := task.RunTime / int64(time.Second)
			nsec := task.RunTime % int64(time.Second)

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
					go task.Job.Run()
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
func (scheduler *TaskScheduler) GetTask() (task *Task, tempKey int) {
	scheduler.Lock()
	defer scheduler.UnLock()

	min := scheduler.tasks[0].RunTime
	tempKey = 0

	for key, task := range scheduler.tasks {

		if min <= task.RunTime {
			continue
		}
		if min > task.RunTime {
			tempKey = key

			min = task.RunTime
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

		if nowTask.Spacing > 0 {
			nowTask.RunTime += nowTask.Spacing * int64(time.Second)
			if nowTask.Number > 1 {
				nowTask.Number --
				scheduler.tasks = append(scheduler.tasks, nowTask)
			} else if nowTask.EndTime >= nowTask.RunTime {
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
		if task.Uuid == uuidStr {
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
		scheduler.swap = make([]*Task, 0)
	}

}