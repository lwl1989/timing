package timer

import (
	"context"
	"errors"
	"runtime"
	"time"
)

//default Job
type taskJob struct {
	Fn      func()
	finish  chan interface{}
	replies map[string]func(reply Reply)
	ctx     context.Context
	cancel  context.CancelFunc
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

func (j *taskJob) loop() {

	task := j.GetTask()
	loop := false
	now := time.Now().UnixNano()
	spacing := task.GetSpacing()
	if task.GetRunNumber() > 1 {
		task.SetRunNumber(task.GetRunNumber() - 1)
		loop = true
	} else if task.GetEndTime() > now && spacing > 0 {
		loop = true
	}

	if loop {
		if spacing > 0 {
			//must use now time
			//task.SetRuntime(task.GetRunTime() + task.GetSpacing())
			task.SetRuntime(now + spacing)
			TS.addTaskChannel(task)
		}
	}

}

//run job and catch error
func (j *taskJob) Run() {
	if f, ok := j.replies["start"]; ok {
		f(GetReply(j.Task, "-1", "start running", nil))
	}

	go func() {
		defer func() {
			e := recover()
			if e != nil {
				if f, ok := j.replies["error"]; ok {
					f(getDefaultErrorReply(j.Task, panicToError(e)))
				}
			} else {
				j.finish <- true
			}
		}()
		j.Fn()
	}()
	for {
		select {
		//获取到终止信号
		case <-j.ctx.Done():
			if f, ok := j.replies["stop"]; ok {
				f(GetReply(j.Task, "-1", "stop running", nil))
			}
			j.close()
			return
		case <-j.finish:
			if f, ok := j.replies["finish"]; ok {
				f(getDefaultSuccessReply(j.Task))
			}
			j.loop()
			j.close()
			return
		}
	}
}

func (j *taskJob) close() {
	runtime.Goexit()
}

func panicToError(r interface{}) error {
	var err error
	//check exactly what the panic was and create error.
	switch x := r.(type) {
	case string:
		err = errors.New(x)
	case error:
		err = x
	default:
		err = errors.New("UnKnow panic")
	}
	return err
}
