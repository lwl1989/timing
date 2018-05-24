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
	add    chan *Task
	remove chan string
	stop   chan struct{}
	Logger *log.Logger
}

//return old name with OnceCron
func NewCron() *OnceCron {
	return &OnceCron{
		&TaskScheduler{
			tasks:  make([]*Task, 0),
			add:    make(chan *Task),
			stop:   make(chan struct{}),
			remove: make(chan string),
			Logger: log.New(os.Stdout, "[Control]: ", log.Ldate|log.Ltime|log.Lshortfile),
		},
	}
}

//return a Controller Scheduler
func NewScheduler() *TaskScheduler {
	return &TaskScheduler{
		tasks:  make([]*Task, 0),
		add:    make(chan *Task),
		stop:   make(chan struct{}),
		remove: make(chan string),
		Logger: log.New(os.Stdout, "[Control]: ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

//add spacing time job to list
func (scheduler *TaskScheduler) AddFuncSpace(spaceTime int64, endTime int64, f func()) {
	task := getTaskWithFuncSpacing(spaceTime, endTime, f)
	scheduler.tasks = append(scheduler.tasks, task)
	scheduler.add <- task
}

//add func to list
func (scheduler *TaskScheduler) AddFunc(unixTime int64, f func()) {
	task := getTaskWithFunc(unixTime, f)
	scheduler.tasks = append(scheduler.tasks, task)
	scheduler.add <- task
}

//add a task to list
func (scheduler *TaskScheduler) AddTask(task *Task) string {
	if task.Spacing > 0 && task.RunTime == 0 {
		task.RunTime = time.Now().Unix() + task.Spacing
	}
	if task.Uuid == "" {
		task.Uuid = uuid.New().String()
	}
	scheduler.tasks = append(scheduler.tasks, task)
	scheduler.add <- task
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
	task := getTaskWithFuncSpacing(3600, time.Now().Add(time.Hour * 24 * 365).Unix(), func() {
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
		i64 := task.RunTime - now.Unix()

		var d time.Duration
		if i64 < 0 {
			scheduler.tasks[key].RunTime = now.Unix()
			scheduler.doAndReset(key)
			continue
		} else {
			d = time.Unix(task.RunTime, 0).Sub(now)
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
	fmt.Println(len(scheduler.tasks), key)
	//null pointer
	if key < len(scheduler.tasks) {

		nowTask := scheduler.tasks[key]
		scheduler.tasks = append(scheduler.tasks[:key], scheduler.tasks[key+1:]...)

		if nowTask.Spacing > 0 {
			nowTask.RunTime += nowTask.Spacing
			fmt.Println(nowTask)
			if nowTask.Number > 1 {
				nowTask.Number --
				scheduler.tasks = append(scheduler.tasks, nowTask)
				//scheduler.Logger.Println("addTask",nowTask.toString())
			} else if nowTask.EndTime >= nowTask.RunTime {
				scheduler.tasks = append(scheduler.tasks, nowTask)
				//scheduler.Logger.Println("addTask",nowTask.toString())
			}
		}

	}
}

func (scheduler *TaskScheduler) removeTask(uuidStr string) {
	for key, task := range scheduler.tasks {
		if task.Uuid == uuidStr {
			scheduler.tasks = append(scheduler.tasks[:key], scheduler.tasks[key+1:]...)
			break
		}
	}
}
