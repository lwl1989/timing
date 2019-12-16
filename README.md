# v2.1.1 版本定时器

# 新特性

1. job自定义
2. 返回任务状态
3. list转换成内置锁map
4. 支持手动任务终止


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
type TaskJob struct {
    Fn func()
    err     chan  error
    done    chan  bool
    stop    chan  bool
    replies map[string]func(reply Reply)
    Task TaskInterface
}
```

# 快速入门

```
    go get -u github.com/lwl1989/timing

    timing.GetTaskScheduler().AddFuncSpace(1, time.Now().UnixNano()+int64(time.Second*1), func() {
        fmt.Println("one second after")
        //do more logic
    })
```
