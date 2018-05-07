package timing

import (
	"time"
	"log"
)

//only exec one timer cron
type OnceCron struct {
	tasks []*Task
	add   chan *Task
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
func (one *OnceCron) AddFuncSpace(unixTime int64, f func()) {
	task := getTaskWithFuncSpacing(unixTime, f)
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
func (one *OnceCron) AddTask(task *Task) {
	if task.Spacing > 0  && task.RunTime == 0{
		task.RunTime = time.Now().Unix()+task.Spacing
	}
	one.tasks = append(one.tasks, task)
	one.add <- task
}

func (one *OnceCron) resetCronTask(key int) {
	if one.tasks[key].Spacing > 0 {
		one.tasks[key].RunTime = one.tasks[key].RunTime + one.tasks[key].Spacing
	}else{
		one.tasks = append(one.tasks[:key], one.tasks[key+1:]...)
	}
	one.add <- one.tasks[key]
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
				if key != -1 {
					nowTask := one.tasks[key]
					if nowTask.Spacing > 0 {
						one.resetCronTask(key)
					}else{
						one.tasks = append(one.tasks[:key], one.tasks[key+1:]...)
					}
				}
				go task.Job.Run()

				//if add a new task and runtime < now task runtime
				// stop now timer and again
			case t := <-one.add:
				if t.RunTime < task.RunTime {
					timer.Stop()
				}

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
