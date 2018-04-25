

package timing

import (
"time"
"fmt"
)

//test add Func
func addFuncTest() {
	cron := NewCron()

	go cron.Start()

	cron.AddFunc(time.Now().Unix()+1, func() {
		fmt.Println("one second after")
	})

	cron.AddFunc(time.Now().Unix()+1, func() {
		fmt.Println("one second after, task second")
	})

	cron.AddFunc(time.Now().Unix()+10, func() {
		fmt.Println("ten second after")
	})
}


//test add Task and timing add Task
func addTaskTest() {
	cron := NewCron()
	go cron.Start()

	cron.AddTask(&Task{
		Job:FuncJob(func() {
			fmt.Println("hello cron")
		}),
		RunTime:time.Now().Unix()+2,
	})


	cron.AddTask(&Task{
		Job:FuncJob(func() {
			fmt.Println("hello cron1")
		}),
		RunTime:time.Now().Unix()+3,
	})

	cron.AddTask(&Task{
		Job: FuncJob(func() {
			fmt.Println("hello cron2")
		}),
		RunTime: time.Now().Unix() + 4,
	})

	timer := time.NewTimer(10 * time.Second)
	for {
		select {
		case <-timer.C:
			cron.AddTask(&Task{
				Job: FuncJob(func() {
					fmt.Println("hello cron2")
				}),
				RunTime: time.Now().Unix() + 1,
			})
		}
		break
	}
}