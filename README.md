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
    
timer := time.NewTimer(10*time.Second)
for {
    select {
    case  <-timer.C:
    	cron.AddTask(&Task{
    		Job:FuncJob(func() {
    				fmt.Println("hello cron2")
    		}),
    		RunTime:time.Now().Unix()+1,
    	})
    }
    break
}
```

### next

1. Add export task and import task
2. Save the task list to the persistence layer.


### example

1. when user order expire, change order status => fail
2. Timing of the message like[fcm-message](https://github.com/lwl1989/TTTask)
and......