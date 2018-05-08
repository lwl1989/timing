# timing

A one-off timed task.

### Thanks

[robfig/cron](https://github.com/robfig/cron)

### Quick Start

```
cron := NewCron()

go cron.Start()

cron.AddFunc(time.Now().Unix()+1, func() {
	fmt.Println("one second after")
})

cron.AddTask(&Task{
	Job:FuncJob(func() {
    		fmt.Println("hello cron2")
    }),
    RunTime:time.Now().Unix()+4,
})

//block it
```

add cron task
```
//10 seconds print one
cron.AddFuncSpace(10, func() {
	fmt.Println("one second after")
})


cron.AddTask(&Task{
	Job:FuncJob(func() {
    		fmt.Println("hello cron2")
    }),
    Spacing:4 //4 seconds send one
})

cron.AddTask(&Task{
	Job:FuncJob(func() {
    		fmt.Println("hello cron2")
    }),
    Spacing:4 //4 seconds send one
    Number: 5 //exec 5 num go stop
})

cron.AddTask(&Task{
	Job:FuncJob(func() {
    		fmt.Println("hello cron2")
    }),
    Spacing:4 //4 seconds send one
    EndTime: 1999999999   // at 199999999 go stop
})
```

### next

1. Add export task and import task
2. Save the task list to the persistence layer.


### example

1. when user order expire, change order status => fail
2. Timing of the message like[fcm-message](https://github.com/lwl1989/TTTask)

and......