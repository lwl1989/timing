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
```

执行结果
```
=== RUN   TestJobEvent
2021/03/03 18:03:52 task uuid:test run start
2021/03/03 18:03:52 task uuid:test start time2021-03-03 18:03:52
now is run job
now job success
2021/03/03 18:03:53 task uuid:testsuccess
2021/03/03 18:03:53 task uuid:test finish time2021-03-03 18:03:53
2021/03/03 18:03:56 task uuid:test run start
2021/03/03 18:03:56 task uuid:test start time2021-03-03 18:03:56
now is run job
now job success
2021/03/03 18:03:57 task uuid:testsuccess
2021/03/03 18:03:57 task uuid:test finish time2021-03-03 18:03:57
2021/03/03 18:04:00 task uuid:test run start
2021/03/03 18:04:00 task uuid:test start time2021-03-03 18:04:00
now is run job
over
--- PASS: TestJobEvent (10.00s)
```

手动停止正在执行的任务：
```
cron := GetTaskScheduler()
	cron.Start()
	f := func() {
		fmt.Println("now is run job")
		time.Sleep(1 * time.Second)
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
```

执行结果
```
=== RUN   TestJobStopEvent
now is run job
2021/03/03 18:07:22 task uuid:123 stop time2021-03-03 18:07:22
over
--- PASS: TestJobStopEvent (10.00s)
```