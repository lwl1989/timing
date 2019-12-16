package timer

import (
    "github.com/google/uuid"
    "time"
    "fmt"
)

func (task *Task) RunJob() {
    task.GetJob().Run()
}

func (task *Task) GetJob()  IJob {
    return task.Job
}

func (task *Task) GetUuid() string {
    return task.Uuid
}

func (task *Task) GetRunTime() int64 {
    return task.RunTime
}

func (task *Task) GetSpacing() int64 {
    return task.Spacing
}

func (task *Task) GetEndTime() int64 {
    return task.EndTime
}

func (task *Task) GetRunNumber() int {
    return task.Number
}

func (task *Task) SetJob(job IJob) TaskSetInterface {
    task.Job = job
    return task
}

func (task *Task)  SetRuntime(runtime int64) TaskSetInterface {
    // if is unix second
    if runtime < 100000000000 {
        runtime = runtime * 1000
    }
    task.RunTime = runtime
    return task
}

func (task *Task)  SetUuid(uuid string) TaskSetInterface {
    task.Uuid = uuid
    return task
}

func (task *Task)  SetSpacing(spacing int64) TaskSetInterface {
    task.Spacing = spacing
    return task
}

func (task *Task)  SetEndTime(endTime int64) TaskSetInterface {
    task.EndTime = endTime
    return task
}

func (task *Task)  SetRunNumber(number int) TaskSetInterface {
    task.Number = number
    return task
}

func (task *Task)  SetStatus(status int) TaskSetInterface {
    task.Status = status
    return task
}

func (task *Task)  GetStatus() int {
    return task.Status
}

//get a new Job
func getJob(f func()) *TaskJob {
    return &TaskJob{
        Fn:f,
        stop:   make(chan bool,1),
        err:    make(chan error,1),
        done:   make(chan bool,1),
        replies:make(map[string]func(reply Reply)),
    }
}

//get task with func
func getTaskWithFunc(unixTime int64, f func()) *Task {
    return &Task{
        Job:    getJob(f),
        RunTime: unixTime,
        Uuid:	 uuid.New().String(),
    }
}

//get task with func and spacing
func getTaskWithFuncSpacingNumber(spacing int64, number int, f func()) *Task {
    return &Task{
        Job:    getJob(f),
        RunTime: time.Now().UnixNano()+spacing,
        Spacing: spacing,
        Number:  number,
        EndTime: time.Now().UnixNano()+ int64(number)*spacing*int64(time.Second),
        Uuid:	 uuid.New().String(),
    }
}

//get task with spacing
func getTaskWithFuncSpacing(spacing int64, endTime int64, f func()) *Task {
    return &Task{
        Job:    getJob(f),
        RunTime: time.Now().UnixNano()+ int64(time.Second)*spacing,
        Spacing: spacing,
        EndTime: endTime,
        Uuid:	 uuid.New().String(),
    }
}

//task toString
func (task *Task) toString() string {
    return fmt.Sprintf("uuid: %s, runTime %d, spaceing %d, endTime　%d, number %d",task.Uuid,task.RunTime,task.Spacing,task.EndTime,task.Number)
}