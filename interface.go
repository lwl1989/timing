package timer


type TaskInterface interface {
    TaskGetInterface
    TaskSetInterface
}

type TaskSetInterface interface {
    SetRuntime(runtime int64) TaskSetInterface
    SetUuid(uuid string) TaskSetInterface
    SetSpacing(spacing int64) TaskSetInterface
    SetEndTime(endTime int64) TaskSetInterface
    SetRunNumber(number int) TaskSetInterface
    SetStatus(status int) TaskSetInterface
}

type TaskGetInterface interface{
    RunJob()
    GetJob()  *taskJob
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
