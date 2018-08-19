

package timing

import (
	"time"
	"fmt"
	"testing"
)

//test add Func
func Test_AddFunc(t *testing.T) {
	cron := NewCron()

	go cron.Start()

	cron.AddFunc(time.Now().UnixNano()+1, func() {
		fmt.Println("one second after")
	})

	cron.AddFunc(time.Now().UnixNano()+1, func() {
		fmt.Println("one second after, task second")
	})

	cron.AddFunc(time.Now().UnixNano()+10, func() {
		fmt.Println("ten second after")
	})
}

//test add space task func
func Test_AddFuncSpace(t *testing.T) {
	cron := NewScheduler()

	go cron.Start()

	cron.AddFuncSpace(1, time.Now().UnixNano()+10, func() {
		fmt.Println("one second after")
	})

	cron.AddFuncSpace(1, time.Now().UnixNano()+20,func() {
		fmt.Println("one second after, task second")
	})

	cron.AddFunc(time.Now().UnixNano()+10, func() {
		fmt.Println("ten second after")
	})
}

//test add Task and timing add Task
func Test_AddTask(t *testing.T) {
	cron := NewCron()
	go cron.Start()

	cron.AddTask(&Task{
		Job:FuncJob(func() {
			fmt.Println("hello cron")
		}),
		RunTime:time.Now().UnixNano()+2,
	})


	cron.AddTask(&Task{
		Job:FuncJob(func() {
			fmt.Println("hello cron1")
		}),
		RunTime:time.Now().UnixNano()+3,
	})

	cron.AddTask(&Task{
		Job: FuncJob(func() {
			fmt.Println("hello cron2")
		}),
		RunTime: time.Now().UnixNano() + 4,
	})

	timer := time.NewTimer(10 * time.Second)
	for {
		select {
		case <-timer.C:
			cron.AddTask(&Task{
				Job: FuncJob(func() {
					fmt.Println("hello cron2")
				}),
				RunTime: time.Now().UnixNano() + 1,
			})
		}
		break
	}
}