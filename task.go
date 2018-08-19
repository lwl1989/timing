package timing

import (
	"time"
	"github.com/google/uuid"
	"fmt"
)


//need to do task has interface Job
type Task struct {
	Job     Job
	Uuid    string
	RunTime int64   //UnixNanoTime
	Spacing int64	//spacing sencond
	EndTime int64   //UnixNanoTime
	Number  int //exec number
}

//callback function has interface Job
type FuncJob func()

func (f FuncJob) Run() { f() }

type Job interface {
	Run()
}

func getTaskWithFunc(unixTime int64, f func()) *Task {
	return &Task{
		Job:     FuncJob(f),
		RunTime: unixTime,
		Uuid:uuid.New().String(),
	}
}

func getTaskWithFuncSpacingNumber(spacing int64, number int, f func()) *Task {
	return &Task{
		Job:     FuncJob(f),
		RunTime: time.Now().UnixNano()+spacing,
		Spacing: spacing,
		Number:  number,
		EndTime: time.Now().UnixNano()+ int64(number)*spacing*int64(time.Second),
		Uuid:uuid.New().String(),
	}
}
func getTaskWithFuncSpacing(spacing int64, endTime int64, f func()) *Task {
	return &Task{
		Job:     FuncJob(f),
		RunTime: time.Now().UnixNano()+ int64(time.Second)*spacing,
		Spacing: spacing,
		EndTime: endTime,
		Uuid:uuid.New().String(),
	}
}

func (task *Task) toString() string {
	return fmt.Sprintf("uuid: %s, runTime %d, spaceing %d, endTime　%d, number %d",task.Uuid,task.RunTime,task.Spacing,task.EndTime,task.Number)
}