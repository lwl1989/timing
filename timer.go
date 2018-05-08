package timing

import (
	"time"
	"log"
	"github.com/google/uuid"
)

//only exec one timer cron
type OnceCron struct {
	tasks []*Task
	add   chan *Task
	remove chan string
	stop  chan struct{}
}

//return a new Cron
func NewCron() *OnceCron {
	return &OnceCron{
		tasks: make([]*Task, 0),
		add:   make(chan *Task),
		stop:  make(chan struct{}),
	}
}

//add spacing time job to list
func (one *OnceCron) AddFuncSpace(unixTime int64,endTime int64, f func()) {
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
	if task.Spacing > 0  && task.RunTime == 0{
		task.RunTime = time.Now().Unix()+task.Spacing
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

//run Cron
func (one *OnceCron) Start() {
	go one.run()
}

//run Cron
func (one *OnceCron) Stop() {
	one.stop <- struct{}{}
}

//return A task after a year to exec
func (one *OnceCron) sleep() (task *Task) {
	rs := FuncJob(func() {
		log.Println("It's a Year timer!")
	})
	return &Task{
		Job:     rs,
		RunTime: time.Now().Unix() + 365*3600*24,
	}
}

//run task list
//if is empty, run a year timer task
func (one *OnceCron) run() {

	for {

		now := time.Now()
		var task *Task
		key := -1
		if len(one.tasks) == 0 {
			task = one.sleep()
		} else {
			task, key = one.GetTask()
		}

		i64 := task.RunTime - now.Unix()

		var d time.Duration
		if i64 < 0 {
			d = time.Microsecond * 10
		} else {
			d = time.Unix(task.RunTime, 0).Sub(now)
		}

		timer := time.NewTimer(d)

		//catch a chan and do something
		for {
			select {
			//if time has expired do task and shift key if is task list
			case <-timer.C:
				go task.Job.Run()
				one.resetTask(key)

			case <-one.add:
				timer.Stop()

			//remove uuid
			case uuidstr:= <-one.remove:
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
//if is null key === -1
func (one *OnceCron) GetTask() (task *Task, tempKey int) {

	min := one.tasks[0].RunTime
	tempKey = -1

	for key, task := range one.tasks {
		if task.RunTime <= min {
			tempKey = key
			break
		}
		if min > task.RunTime {
			tempKey = key
			min = task.RunTime
			continue
		}
	}

	if tempKey == -1 {
		return nil, -1
	}

	task = one.tasks[tempKey]

	return task, tempKey
}

//if add a new task and runtime < now task runtime
// stop now timer and again
func (one *OnceCron) resetTask(key int) {
	if key != -1 {

		nowTask := one.tasks[key]
		log.Println(nowTask)
		one.tasks = append(one.tasks[:key], one.tasks[key+1:]...)

		if nowTask.Spacing > 0 {
			nowTask.RunTime += nowTask.Spacing
			if nowTask.Number > 1 {
				nowTask.Number --
				go one.AddTask(nowTask)
			}else if nowTask.EndTime >= nowTask.RunTime {
				go one.AddTask(nowTask)
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