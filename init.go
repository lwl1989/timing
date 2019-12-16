package timer

import (
    "log"
    "os"
    "sync"
)

var successResult = Reply{
    Code:200,
    Msg:"操作成功",
    Err:nil,
}

var TS *TaskScheduler

func init()  {
    TS = &TaskScheduler{
        tasks:  new(sync.Map),
        running:new(sync.Map),
        add:    make(chan TaskInterface),
        stop:   make(chan struct{}),
        remove: make(chan string),
        Logger: log.New(os.Stdout, "[Control]: ", log.Ldate|log.Ltime|log.Lshortfile),
    }
}

func GetTaskScheduler() *TaskScheduler {
    return TS
}
