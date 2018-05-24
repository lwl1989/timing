# timing

A one-off timed task.

### Thanks

[robfig/cron](https://github.com/robfig/cron)

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

### next

1. Add export task and import task
2. Save the task list to the persistence layer.


### example

1. when user order expire, change order status => fail
2. Timing of the message like[fcm-message](https://github.com/lwl1989/TTTask)

and......