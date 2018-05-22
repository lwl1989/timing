package timing

import (
	"time"
	"log"
	"github.com/google/uuid"
	"os"
	"fmt"
)

//only exec one timer cron
type OnceCron struct {
	tasks  []*Task
	add    chan *Task
	remove chan string
	stop   chan struct{}
	Logger *log.Logger
}

//return a new Cron
func NewCron() *OnceCron {
	return &OnceCron{
		tasks:  make([]*Task, 0),
		add:    make(chan *Task),
		stop:   make(chan struct{}),
		remove: make(chan string),
		Logger: log.New(os.Stdout,"[cron]: ",log.Ldate|log.Ltime|log.Lshortfile),
	}
}

//add spacing time job to list
func (one *OnceCron) AddFuncSpace(unixTime int64, endTime int64, f func()) {
	task := getTaskWithFuncSpacing(unixTime, endTime, f)
	one.tasks = append(one.tasks, task)
	one.add <- task
}

//add func to list
func (one *OnceCron) AddFunc(unixTime int64, f func()) {
	task := getTaskWithFunc(unixTime, f)
	one.tasks = append(one.tasks, task)
	one.add <- task
}

//add a task to list
func (one *OnceCron) AddTask(task *Task) string {
	if task.Spacing > 0 && task.RunTime == 0 {
		task.RunTime = time.Now().Unix() + task.Spacing
	}
	if task.Uuid == "" {
		task.Uuid = uuid.New().String()
	}
	one.tasks = append(one.tasks, task)
	one.add <- task
	return task.Uuid
}

//export tasks
func (one *OnceCron) export() []*Task {
	return one.tasks
}

//stop tasks
func (one *OnceCron) StopOnce(uuidStr string) {
	one.remove <- uuidStr
}

//run Cron
func (one *OnceCron) Start() {
	//初始化的時候加入一個一年的長定時器,間隔1小時執行一次
	task := getTaskWithFuncSpacing(3600, time.Now().Add(time.Hour * 24 * 365).Unix(), func() {
		log.Println("It's a Hour timer!")
	})
	one.tasks = append(one.tasks, task)
	go one.run()
}

//run Cron
func (one *OnceCron) Stop() {
	one.stop <- struct{}{}
}

//run task list
//if is empty, run a year timer task
func (one *OnceCron) run() {

	for {

		now := time.Now()
		task, key := one.GetTask()
		i64 := task.RunTime - now.Unix()

		var d time.Duration
		if i64 < 0 {
			one.tasks[key].RunTime = now.Unix()
			one.doAndReset(key)
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
				one.doAndReset(key)
				if task != nil {
					//fmt.Println(one.tasks[key])
					go task.Job.Run()
					timer.Stop()
				}

				//if add task
			case <-one.add:
				timer.Stop()
				// remove task with remove uuid
			case uuidstr := <-one.remove:
				one.removeTask(uuidstr)
				timer.Stop()
				//if get a stop single exit
			case <-one.stop:
				timer.Stop()
				return
			}

			break
		}
	}
}

//return a task and key In task list
func (one *OnceCron) GetTask() (task *Task, tempKey int) {

	min := one.tasks[0].RunTime
	tempKey = 0

	for key, task := range one.tasks {

		if min <= task.RunTime {
			continue
		}
		if min > task.RunTime {
			tempKey = key

			min = task.RunTime
			continue
		}
	}

	task = one.tasks[tempKey]

	return task, tempKey
}

//if add a new task and runtime < now task runtime
// stop now timer and again
func (one *OnceCron) doAndReset(key int) {
	fmt.Println(len(one.tasks),key)
	//null pointer
	if key < len(one.tasks) {

		nowTask := one.tasks[key]
		one.tasks = append(one.tasks[:key], one.tasks[key+1:]...)

		if nowTask.Spacing > 0 {
			nowTask.RunTime += nowTask.Spacing
			fmt.Println(nowTask)
			if nowTask.Number > 1 {
				nowTask.Number --
				one.tasks = append(one.tasks, nowTask)
				//one.Logger.Println("addTask",nowTask.toString())
			} else if nowTask.EndTime >= nowTask.RunTime {
				one.tasks = append(one.tasks, nowTask)
				//one.Logger.Println("addTask",nowTask.toString())
			}
		}

	}
}

func (one *OnceCron) removeTask(uuidStr string) {
	for key, task := range one.tasks {
		if task.Uuid == uuidStr {
			one.tasks = append(one.tasks[:key], one.tasks[key+1:]...)
			break
		}
	}
}


