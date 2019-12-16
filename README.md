# timing

A one-off timed task.

### Thanks

[robfig/cron](https://github.com/robfig/cron)

### ChangeLog

[2018-05-22]  Add Export

[2018-05-22]  Renamed Tasks to TaskScheduler

[2019-02-23]  Add task interface and can rewrite the logic


### Quick Start

```
scheduler := NewScheduler()

scheduler.Start()

scheduler.AddFunc(time.Now().Unix()+1, func() {
	fmt.Println("one second after")
})

scheduler.AddTask(&Task{
	Job:FuncJob(func() {
    		fmt.Println("hello task2")
    }),
    RunTime:time.Now().Unix()+4,
})

//block it
```

add task
```
//10 seconds print one
scheduler.AddFuncSpace(10, func() {
	fmt.Println("one second after")
})


scheduler.AddTask(&Task{
	Job:FuncJob(func() {
    		fmt.Println("hello task")
    }),
    Spacing:4 //4 seconds send one
})

scheduler.AddTask(&Task{
	Job:FuncJob(func() {
    		fmt.Println("hello task2")
    }),
    Spacing:4 //4 seconds send one
    Number: 5 //exec 5 num go stop
})

scheduler.AddTask(&Task{
	Job:FuncJob(func() {
    		fmt.Println("hello task3")
    }),
    Spacing:4 //4 seconds send one
    EndTime: 1999999999   // at 199999999 go stop
})
```

### task interface

If you need rewrite a new task

implement TaskInterface

```
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
```

### next

1. distributed the cron task


### Example

Order expire, change order status => fail
Timing of the send message

and more