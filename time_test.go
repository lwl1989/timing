

package timing

import (
	"time"
	"fmt"
	"testing"
)

//test add Func
func Test_AddFunc(t *testing.T) {
	cron := NewCron()

	cron.Start()

	go cron.AddFunc(time.Now().UnixNano()+1*(int64(time.Second)), func() {
		fmt.Println("one second after")
	})

	go cron.AddFunc(time.Now().UnixNano()+2*(int64(time.Second)), func() {
		fmt.Println("one second after, task second")
	})

	go cron.AddFunc(time.Now().UnixNano()+10*(int64(time.Second)), func() {
		fmt.Println("ten second after")
	})

	time.Sleep(20*time.Second)
}

//test add space task func
func Test_AddFuncSpace(t *testing.T) {
	cron := NewScheduler()

	cron.Start()

	go cron.AddFuncSpace(1, time.Now().UnixNano()+100*(int64(time.Second)), func() {
		fmt.Println("one second after")
	})

	go cron.AddFuncSpace(2, time.Now().UnixNano()+100*(int64(time.Second)),func() {
		fmt.Println("one second after, task second")
	})

	go cron.AddFunc(time.Now().UnixNano()+10*(int64(time.Second)), func() {
		fmt.Println("ten second after")
	})
	time.Sleep(20*time.Second)
}

//test add Task and timing add Task
func Test_AddTask(t *testing.T) {
	cron := NewCron()
	cron.Start()

	go cron.AddTask(&Task{
		Job:FuncJob(func() {
			fmt.Println("hello cron")
		}),
		RunTime:time.Now().UnixNano()+2,
	})


	go cron.AddTask(&Task{
		Job:FuncJob(func() {
			fmt.Println("hello cron1")
		}),
		RunTime:time.Now().UnixNano()+3,
	})

	go cron.AddTask(&Task{
		Job: FuncJob(func() {
			fmt.Println("hello cron2")
		}),
		RunTime: time.Now().UnixNano() + 4,
	})

	time.Sleep(20*time.Second)
}