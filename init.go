package timer

import (
	"log"
	"os"
	"sync"
)

var TS *TaskScheduler
var stop chan string

func init() {
	TS = &TaskScheduler{
		tasks:   new(sync.Map),
		running: new(sync.Map),
		add:     make(chan TaskInterface), //添加新任务信号
		stop:    make(chan struct{}),      //终止主进程
		remove:  make(chan string),		   //添加子进程
		Logger:  log.New(os.Stdout, "[Control]: ", log.Ldate|log.Ltime|log.Lshortfile),
	}
	stop = make(chan string) //手动关闭正在执行的任务信号
}

func GetTaskScheduler() *TaskScheduler {
	return TS
}
