package timer

import (
    "github.com/pkg/errors"
    "runtime"
    "time"
)

//default Job
type taskJob struct {
	Fn      func()
	err     chan error
	done    chan bool
	stop    chan bool
	replies map[string]func(reply Reply)
	Task    TaskInterface
}

func (j *taskJob) GetTask() TaskInterface {
	return j.Task
}

func (j *taskJob) SetTask(Task TaskInterface) {
	j.Task = Task
}

func (j *taskJob) OnStart(f func(reply Reply)) {
	j.replies["start"] = f
}

func (j *taskJob) OnStop(f func(reply Reply)) {
	j.replies["stop"] = f
}

func (j *taskJob) OnFinish(f func(reply Reply)) {
	j.replies["finish"] = f
}

func (j *taskJob) OnError(f func(reply Reply)) {
	j.replies["error"] = f
}

func (j *taskJob) run() {

	isPanic := false
	defer func() {
		if x := recover(); x != nil {
			err := errors.Errorf("job error with panic:%v", x)
			j.err <- err
			isPanic = true
			return
		}
	}()

	defer func() {
		if !isPanic {
			j.done <- true
		}
	}()
	j.Fn()
}

func (j *taskJob) Stop() {
	j.stop <- true
}

//run job and catch error
func (j *taskJob) Run() {
	if f, ok := j.replies["start"]; ok {
		f(Reply{})
	}

	go j.run()
	for {
		select {
		case e := <-j.err:
			if f, ok := j.replies["error"]; ok {
				reply := Reply{
					Code: 500,
					Msg:  e.Error(),
					Err:  e,
				}
				f(reply)
			}
			return
		case <-j.done:
			if f, ok := j.replies["finish"]; ok {
				f(successResult)
			}

			task := j.GetTask()
			loop := false
			now := time.Now().UnixNano()
			if task.GetRunNumber() > 1 {
				task.SetRunNumber(task.GetRunNumber() - 1)
				loop = true
			} else if task.GetEndTime() > now {
				loop = true
			}

			if loop {
				spacing := task.GetSpacing()
				if spacing > 0 {
					//must use now time
					//task.SetRuntime(task.GetRunTime() + task.GetSpacing())
					task.SetRuntime(now + spacing)
					TS.addTaskChannel(task)
				}
			}else {
                j.close(false)
            }
			return
		case <-j.stop:
            j.close(true)
		}
	}
}

func (j *taskJob) close(exit bool) {
	close(j.done)
	close(j.err)
	close(j.stop)
	if exit {
        runtime.Goexit()
    }
}
