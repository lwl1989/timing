# 任务定时器 v2.1.1

# 新特性

1. 返回任务状态
2. list转换成内置锁map
3. 支持手动任务终止

tips: 间隔时间需要严格按照nanoSecond计算

# 新结构体

任务执行结果
```
type Reply struct {
    Code int64 `json:"code"`
    Msg string `json:"msg"`
    Err error   `json:"err"`
    Ts  Task  `json:"task"`
}
```

抽象的任务
```
type taskJob struct {
    Fn func()
    err     chan  error
    done    chan  bool
    stop    chan  bool
    replies map[string]func(reply Reply)
    Task TaskInterface
}
```

# 快速入门

    go get -u github.com/lwl1989/timing
    
```
    cron := GetTaskScheduler()
    
    go cron.Start()
    //添加一次性任务
    cron.AddFunc(time.Now().UnixNano()+int64(time.Second*1), func() {
    	fmt.Println("one second after")
    })
    //添加循环任务，间隔时间为1秒，截止时间为10秒后
    cron.AddFuncSpace(int64(time.Second*1), time.Now().UnixNano()+int64(time.Second*10), func() {
        fmt.Println("one second after")
    })
    //添加指定执行次数任务，并指定每次间隔时间
    cron.AddFuncSpaceNumber(int64(time.Second*1), 10, func() {
        fmt.Println("number 10")
    })
```



# 自定义执行事件，支持的事件有

1. OnStart(f func(reply Reply))
2. OnStop(f func(reply Reply))
3. OnFinish(f func(reply Reply))
4. OnError(f func(reply Reply)) 

example:
```
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
    }
    t1.GetJob().OnStart(f1)
    cron.AddTask(t1)
```