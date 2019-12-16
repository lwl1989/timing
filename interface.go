package timer

type IJob interface {
    Run()
    GetTask() TaskInterface
    SetTask(Task TaskInterface)
    OnStart(f func(reply Reply))
    OnStop(f func(reply Reply)) //todo: sign how todo
    OnFinish(f func(reply Reply))
    OnError(f func(reply Reply))
    Stop()
}

type TaskInterface interface {
    TaskGetInterface
    TaskSetInterface
}

type TaskSetInterface interface {
    SetJob(job IJob) TaskSetInterface
    SetRuntime(runtime int64) TaskSetInterface
    SetUuid(uuid string) TaskSetInterface
    SetSpacing(spacing int64) TaskSetInterface
    SetEndTime(endTime int64) TaskSetInterface
    SetRunNumber(number int) TaskSetInterface
    SetStatus(status int) TaskSetInterface
}

type TaskGetInterface interface{
    RunJob()
    GetJob()  IJob
    GetUuid() string
    GetRunTime() int64
    GetSpacing() int64
    GetEndTime() int64
    GetRunNumber() int
    GetStatus() int
}

type TaskLogInterface interface {
    Println(v ...interface{})
}
