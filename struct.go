package timer

import "sync"

type Reply struct {
    Code int64 `json:"code"`
    Msg string `json:"msg"`
    Err error   `json:"err"`
    Ts  Task  `json:"task"`
}

type TaskScheduler struct {
    tasks   *sync.Map
    running *sync.Map
    add     chan TaskInterface
    remove  chan string
    stop    chan struct{}
    Logger  TaskLogInterface
}

//need to do task has interface Job
type Task struct {
    Job     *taskJob    `json:"job"`
    Uuid    string  `json:"uuid"`
    RunTime int64   `json:"run_time"`//UnixNanoTime
    Spacing int64	`json:"spacing"`//spacing sencond
    EndTime int64   `json:"end_time"`//UnixNanoTime
    Number  int `json:"number"`//exec number
    Status  int `json:"status"`
}


//
//func GetSuccessResult(msg string)  Reply {
//    if msg != "" {
//        successResult.Msg = msg
//    }
//    return successResult
//}
//
//func GetErrorResult(code int64, msg string, err error) Reply  {
//    return Reply{
//        Code:code,
//        Msg:msg,
//        Err:err,
//    }
//}