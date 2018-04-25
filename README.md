# timing


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