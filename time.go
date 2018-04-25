package timing

/**/

//need to do task has interface Job
type Task struct {
	Job     Job
	RunTime int64
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
	}
}
