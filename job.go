package timer

import (
    "fmt"
    "github.com/pkg/errors"
    "runtime"
    "time"
)

//default Job
type TaskJob struct {
    Fn func()
    err     chan  error
    done    chan  bool
    stop    chan  bool
    replies map[string]func(reply Reply)
    Task TaskInterface
}

func (j *TaskJob) GetTask() TaskInterface  {
    return j.Task
}

func (j *TaskJob) SetTask(Task TaskInterface) {
    j.Task = Task
}

func (j *TaskJob) OnStart(f func(reply Reply))  {
    j.replies["start"] = f
}

func (j *TaskJob) OnStop(f func(reply Reply))  {
    j.replies["stop"] = f
}

func (j *TaskJob) OnFinish(f func(reply Reply))  {
    j.replies["finish"] = f
}

func (j *TaskJob) OnError(f func(reply Reply))  {
    j.replies["error"] = f
}

func (j *TaskJob) run() {

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

func (j *TaskJob) Stop(){
    j.stop <- true
}

//run job and catch error
func (j *TaskJob) Run(){
    if f,ok := j.replies["start"]; ok {
        f(Reply{})
    }

    go j.run()

    for {
        select {
        case e := <-j.err:
            if f,ok := j.replies["error"]; ok {
                reply := Reply{
                    Code:500,
                    Msg:e.Error(),
                    Err:e,
                }
                f(reply)
            }
            return
        case <-j.done:
            if f,ok := j.replies["finish"]; ok {
                f(successResult)
            }

            task := j.GetTask()
            loop := false
            now := time.Now().UnixNano()
            if task.GetRunNumber() > 1 {
                task.SetRunNumber(task.GetRunNumber() - 1)
                loop = true
            }else if task.GetEndTime() > now {
                loop = true
            }
            if loop {
                spacing := task.GetSpacing()
                fmt.Println(spacing)
                if spacing > 0 {
                    //not use old time
                    //task.SetRuntime(task.GetRunTime() + task.GetSpacing())
                    task.SetRuntime(now + spacing)
                    TS.addTaskChannel(task)
                }
            }


            //TS.tasks.Range(func(key, value interface{}) bool {
            //    return  true
            //})
            return
        case <-j.stop:
            //todo:
            runtime.Goexit()
        }
    }
}