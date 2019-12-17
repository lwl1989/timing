package timer

import (
    "fmt"
    "github.com/google/uuid"
    "log"
    "time"
)


//add spacing time job to list with number
func (scheduler *TaskScheduler) AddFuncSpaceNumber(spaceTime int64, number int, f func()) {
    task := getTaskWithFuncSpacingNumber(spaceTime, number, f)
    scheduler.AddTask(task)
}
//add spacing time job to list with endTime
//spaceTime is nano time
func (scheduler *TaskScheduler) AddFuncSpace(spaceTime int64, endTime int64, f func()) {
    task := getTaskWithFuncSpacing(spaceTime, endTime, f)
    scheduler.AddTask(task)
}

//add func to list
func (scheduler *TaskScheduler) AddFunc(unixTime int64, f func()) {
    task := getTaskWithFunc(unixTime, f)
    scheduler.AddTask(task)
}

func (scheduler *TaskScheduler) AddTaskInterface(task TaskInterface) {
    scheduler.addTaskChannel(task)
}

//add a task to list
func (scheduler *TaskScheduler) AddTask(task *Task) string {
    now := time.Now().UnixNano()
    if task.RunTime != 0 {
        if task.RunTime < 9999999999 {
            task.RunTime = task.RunTime * int64(time.Second)
        }
        if task.RunTime <= now {

            if task.Spacing > 0 {
                task.RunTime = now + task.Spacing
            }else {
                //延遲1秒
                task.RunTime = now + int64(time.Second)
            }
        }
    } else {
        if task.Spacing > 0 {
            task.RunTime = now + task.Spacing
        }else{
            scheduler.Logger.Println("error too add task! Runtime error")
            return ""
        }
    }

    if task.Uuid == "" {
        task.Uuid = uuid.New().String()
    }
    scheduler.addTaskChannel(task)
    return task.GetUuid()
}

func (scheduler *TaskScheduler) addTask(task TaskInterface) string  {
    scheduler.tasks.Store(task.GetUuid(), task)
    return task.GetUuid()
}

func  (scheduler *TaskScheduler)  addTaskChannel(task TaskInterface) {
    scheduler.add <- task
}
//new export
func (scheduler *TaskScheduler) ExportInterface() []TaskInterface {
    tasks := make([]TaskInterface, 0)
    scheduler.tasks.Range(func(key, value interface{}) bool {
        switch value.(type) {
        case TaskInterface:
            tasks = append(tasks, value.(TaskInterface))
        }

        return true
    })
    return tasks
}
//compatible old export tasks
func (scheduler *TaskScheduler) Export() []*Task {
    tasks := make([]*Task,0)
    scheduler.tasks.Range(func(key, value interface{}) bool {
        switch value.(type) {
        case *Task:
            tasks = append(tasks, value.(*Task))
        }

        return true
    })
    return tasks
}

//stop task with uuid
func (scheduler *TaskScheduler) StopOnce(uuidStr string) {
    scheduler.remove <- uuidStr
}

//run Cron
func (scheduler *TaskScheduler) Start() {
    //初始化的時候加入一個一年的長定時器,間隔1小時執行一次
    task := getTaskWithFuncSpacing(int64(3600*time.Second), time.Now().Add(time.Hour * 24 * 365).UnixNano(), func() {
        log.Println("It's a Hour timer!")
    })
    scheduler.tasks.Store(task.Uuid, task)
    go scheduler.run()
}

//stop all
func (scheduler *TaskScheduler) Stop() {
    scheduler.stop <- struct{}{}
}

//run task list
//if is empty, run a year timer task
func (scheduler *TaskScheduler) run() {

    for {
        now := time.Now()
        task := scheduler.GetTask()
        var d time.Duration
        if task == nil {
            d = 1 * time.Second
        }else{
            task.GetJob().SetTask(task)
            runTime := task.GetRunTime()
            i64 := runTime - now.UnixNano()
            if i64 < 0 {
                if task.GetSpacing() > 0 {
                    task.SetRuntime(now.UnixNano() + task.GetSpacing())
                }else{
                    task.SetRuntime(now.UnixNano() + int64(time.Second))
                }
                task.SetStatus(1)
                go task.RunJob()
                continue
            } else {
                sec := runTime / int64(time.Second)
                nsec := runTime % int64(time.Second)

                d = time.Unix(sec, nsec).Sub(now)
            }
        }

        fmt.Println(d)
        timer := time.NewTimer(d)

        //catch a chan and do something
        for {
            select {
            //if time has expired do task and shift key if is task list
            case <-timer.C:
                //not get remove,just run
                scheduler.removeTask(task.GetUuid())
                go task.RunJob()
                timer.Stop()
                //if add task
            case t1 := <-scheduler.add:
                scheduler.addTask(t1)
                timer.Stop()
                // remove task with remove uuid
            case uuidStr := <-scheduler.remove:
                scheduler.removeTask(uuidStr)
                timer.Stop()
                //if get a stop single exit
            case <-scheduler.stop:
                timer.Stop()
                return
            }

            break
        }
    }
}

//return a task and key In task list
func (scheduler *TaskScheduler) GetTask() (task TaskInterface) {
    var min int64 = 0
    scheduler.tasks.Range(func(key, value interface{}) bool {
        switch value.(type) {
            case TaskInterface:

            t := value.(TaskInterface)
            runTime := t.GetRunTime()
            if min == 0 {
                min = runTime
                task = t
            } else {
                if min > runTime {
                    min = runTime
                    task = t
                }
            }
        }

        return true
    })

    //if task != nil {
    //    scheduler.removeTask(task.GetUuid())
    //}
    return task
}

//remove task by uuid
func (scheduler *TaskScheduler) removeTask(uuidStr string) {
    scheduler.tasks.Delete(uuidStr)
}


