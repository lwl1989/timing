package timer

import (
    "github.com/pkg/errors"
    "runtime"
    "fmt"
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


    j.Fn()

    defer func() {
        if !isPanic {
            j.done <- true
        }
    }()

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
            task.SetRuntime(task.GetRunTime() + task.GetSpacing())
            TS.addTask(task)
            TS.tasks.Range(func(key, value interface{}) bool {
                fmt.Println(key, value)
                return  true
            })
            return
        case <-j.stop:
            //todo:
            runtime.Goexit()
        }
    }
}