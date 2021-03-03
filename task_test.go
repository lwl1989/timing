package timer

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/grestful/utils"
	"log"
	"os"
	"testing"
	"time"
)

func TestTime(t *testing.T) {
	startTime := "2019-08-11 21:30:01"
	endTime := "2019-08-11 21:45:00"
	cal := "2019-08-12 21:30:07"

	t1, err := time.ParseInLocation("2006-01-02 15:04:05", cal, time.Local)
	s, err1 := time.ParseInLocation("2006-01-02 15:04:05", startTime, time.Local)
	e, err2 := time.ParseInLocation("2006-01-02 15:04:05", endTime, time.Local)

	sTime := time.Now().Format("2006-01-02") + " " + s.Format("15:04:05")
	eTime := time.Now().Format("2006-01-02") + " " + e.Format("15:04:05")
	cal = time.Now().Format("2006-01-02") + " " + t1.Format("15:04:05")

	fmt.Println(t1, err, err1, err2, sTime, eTime)
	if cal >= sTime && cal <= eTime {
		fmt.Println(cal)
	}

}

//test add Func
func Test_AddFunc(t *testing.T) {
	cron := GetTaskScheduler()

	cron.Start()
	cron.AddFunc(time.Now().UnixNano()+int64(time.Second*1), func() {
		fmt.Println("one second after")
	})

	cron.AddFunc(time.Now().UnixNano()+int64(time.Second*1), func() {
		fmt.Println("one second after 1, task second")
	})

	cron.AddFunc(time.Now().UnixNano()+int64(time.Second*10), func() {
		fmt.Println("ten second after")
	})

	cron.tasks.Range(func(key, value interface{}) bool {
		fmt.Println(key, value)
		return true
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
	cron := GetTaskScheduler()
	log.SetOutput(os.Stdout)
	go cron.Start()

	cron.AddFuncSpace(int64(time.Second*1), time.Now().UnixNano()+int64(time.Second*10), func() {
		fmt.Println("one second after")
	})

	cron.AddFuncSpace(int64(time.Second*1), time.Now().UnixNano()+int64(time.Second*10), func() {
		fmt.Println("one second after, task second")
	})

	cron.AddFunc(time.Now().UnixNano()+int64(time.Second*10), func() {
		fmt.Println("ten second after")
	})

	cron.AddFuncSpaceNumber(int64(time.Second*1), 10, func() {
		fmt.Println("number 10")
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
	cron := GetTaskScheduler()
	go cron.Start()

	cron.AddTask(&Task{
		Job: getJob(func() {
			fmt.Println("hello cron")
		}),
		RunTime: time.Now().UnixNano() + int64(time.Second*2),
	})

	cron.AddTask(&Task{
		Job: getJob(func() {
			fmt.Println("hello cron1")
		}),
		RunTime: time.Now().UnixNano() + int64(time.Second*3),
	})

	cron.AddTask(&Task{
		Job: getJob(func() {
			fmt.Println("hello cron2")
		}),
		RunTime: time.Now().UnixNano() + +int64(time.Second*4),
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

func TestJobEvent(t *testing.T) {
	cron := GetTaskScheduler()
	cron.Start()
	f := func() {
		fmt.Println("now is run job")
		time.Sleep(1 * time.Second)
		fmt.Println("now job success")
	}
	t1 := &Task{
		Job:     getJob(f),
		RunTime: time.Now().UnixNano() + int64(time.Second)*1,
		Spacing: int64(3 * time.Second),
		EndTime: time.Now().UnixNano() + int64(time.Second*20),
		Uuid:    "123",
	}
	f1 := func(reply Reply) {
		fmt.Println(reply)
		fmt.Println("It's reply")
		log.Println("task uuid:" + reply.Ts.GetUuid() + " run start")
		log.Println("task uuid:" + reply.Ts.GetUuid() + " start time" + utils.GetTimeString())
	}
	t1.GetJob().OnStart(f1)
	t1.GetJob().OnFinish(func(reply Reply) {
		log.Println("task uuid:" + reply.Ts.GetUuid() + "success")
		log.Println("task uuid:" + reply.Ts.GetUuid() + " finish time" + utils.GetTimeString())
	})
	cron.AddTask(t1)

	timer := time.NewTimer(10 * time.Second)
	for {
		select {
		case <-timer.C:
			fmt.Println("over")
		}
		break
	}
}

func TestTaskLoop(t *testing.T) {
	fmt.Println(uuid.New().String())
	fmt.Println(uuid.New().String())
	fmt.Println(uuid.New().String())
	fmt.Println(uuid.New().String())
	fmt.Println(uuid.New().String())

	//f := func() {
	//	fmt.Println("hello")
	//}
	//t1 := &Task{
	//	Job:     getJob(f),
	//	RunTime: time.Now().UnixNano() + int64(time.Second)*1,
	//	Spacing: int64(3 * time.Second),
	//	EndTime: time.Now().UnixNano() + int64(time.Second*20),
	//	Uuid:    "123",
	//}
}

func TestJobStopEvent(t *testing.T) {
	cron := GetTaskScheduler()
	cron.Start()
	f := func() {
		fmt.Println("now is run job")
		time.Sleep(1 * time.Second)
		fmt.Println("now job success")
	}
	t1 := &Task{
		Job:     getJob(f),
		RunTime: time.Now().UnixNano() + int64(time.Second)*1,
		Spacing: int64(2 * time.Second),
		EndTime: time.Now().UnixNano() + int64(time.Second*20),
		Uuid:    "123",
	}
	f1 := func(reply Reply) {
		log.Println("task uuid:" + reply.Ts.GetUuid() + " stop time" + utils.GetTimeString())
	}
	t1.GetJob().OnStop(f1)
	cron.AddTask(t1)

	go func() {
		t2 := time.NewTimer(2 * time.Second)
		<-t2.C
		cron.StopOnce("123")
	}()

	timer := time.NewTimer(10 * time.Second)
	for {
		select {
		case <-timer.C:
			fmt.Println("over")
		}
		break
	}
}
