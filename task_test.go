

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

	cron.AddFunc(time.Now().UnixNano()+int64(time.Second*1), func() {
		fmt.Println("one second after")
	})

	cron.AddFunc(time.Now().UnixNano()+int64(time.Second*1), func() {
		fmt.Println("one second after, task second")
	})

	cron.AddFunc(time.Now().UnixNano()+int64(time.Second*10), func() {
		fmt.Println("ten second after")
	})

    timer := time.NewTimer(11 * time.Second)
    for {
        select {
        case <-timer.C:
            fmt.Println("over")
        }
        break
    }
}

//test add space task func
func Test_AddFuncSpace(t *testing.T) {
	cron := NewScheduler()

	go cron.Start()

	cron.AddFuncSpace(1, time.Now().UnixNano()+int64(time.Second*1), func() {
		fmt.Println("one second after")
	})

	cron.AddFuncSpace(1, time.Now().UnixNano()+int64(time.Second*20),func() {
		fmt.Println("one second after, task second")
	})

	cron.AddFunc(time.Now().UnixNano()+int64(time.Second*10), func() {
		fmt.Println("ten second after")
	})
    timer := time.NewTimer(11 * time.Second)
    for {
        select {
        case <-timer.C:
            fmt.Println("over")
        }
        break
    }
}

//test add Task and timing add Task
func Test_AddTask(t *testing.T) {
	cron := NewCron()
	go cron.Start()

	cron.AddTask(&Task{
		Job:FuncJob(func() {
			fmt.Println("hello cron")
		}),
		RunTime:time.Now().UnixNano()+int64(time.Second*2),
	})


	cron.AddTask(&Task{
		Job:FuncJob(func() {
			fmt.Println("hello cron1")
		}),
		RunTime:time.Now().UnixNano()+int64(time.Second*3),
	})

	cron.AddTask(&Task{
		Job: FuncJob(func() {
			fmt.Println("hello cron2 loop")
		}),
		RunTime: time.Now().UnixNano() + int64(time.Second*4),
		Spacing: 1,
		EndTime: time.Now().UnixNano() + 9*(int64(time.Second)),
	})

	timer := time.NewTimer(10 * time.Second)
	for {
		select {
		case <-timer.C:
		    fmt.Println("over")
		}
		break
	}
}