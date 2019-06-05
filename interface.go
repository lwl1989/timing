package timing

/**
    Job     Job
    Uuid    string
	RunTime int64   //UnixNanoTime
	Spacing int64	//spacing sencond
	EndTime int64   //UnixNanoTime
	Number  int //exec number
 */
type TaskInterface interface {
    TaskGetInterface
    TaskSetInterface
}

type TaskSetInterface interface {
    SetJob(job Job) TaskSetInterface
    SetRuntime(runtime int64) TaskSetInterface
    SetUuid(uuid string) TaskSetInterface
    SetSpacing(spacing int64) TaskSetInterface
    SetEndTime(endTime int64) TaskSetInterface
    SetRunNumber(number int) TaskSetInterface
}

type TaskGetInterface interface{
    RunJob()
    GetJob()  Job
    GetUuid() string
    GetRunTime() int64
    GetSpacing() int64
    GetEndTime() int64
    GetRunNumber() int
}

type TaskLogInterface interface {
    Println(v ...interface{})
}

func (task *Task) SetJob(job Job) TaskSetInterface {
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


func (task *Task) RunJob() {
    task.GetJob().Run()
}
func (task *Task) GetJob()  Job {
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