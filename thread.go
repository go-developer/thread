// Package thread ...
//
// Author: go_developer@163.com<张德满>
//
// Description: 协程调度主入口
//
// File: thread.go
//
// Version: 1.0.0
//
// Date: 2021/02/06 16:54:44
package thread

import (
	"bytes"
	"context"
	"fmt"
	"runtime"
	"sync"
	"time"
)

// Dispatch 协程调度器实例
var Dispatch *dispatch

func init() {
	Dispatch = &dispatch{
		lock:               &sync.RWMutex{},
		goroutineLockTable: make(map[string]chan int),
		goroutineListTable: make(map[string]chan IGoRoutine),
	}
}

// dispatch 协程调度
//
// Author : go_developer@163.com<张德满>
//
// Date : 2021/02/06 16:55:45
type dispatch struct {
	lock               *sync.RWMutex              // 锁
	goroutineLockTable map[string]chan int        // 协程锁channel
	goroutineListTable map[string]chan IGoRoutine // 协程任务队列
}

// Run 执行一个协程
//
// Author : go_developer@163.com<张德满>
//
// Date : 2021/02/06 17:05:10
func (d *dispatch) Run(goroutine IGoRoutine) {
	d.lock.Lock()
	defer d.lock.Unlock()
	if _, exist := d.goroutineListTable[goroutine.GetGoRoutineName()]; !exist {
		macGoroutineCnt := goroutine.GetMaxGoRoutineCnt()
		if macGoroutineCnt <= 0 {
			macGoroutineCnt = defaultGoroutineCount
		}

		d.goroutineListTable[goroutine.GetGoRoutineName()] = make(chan IGoRoutine, macGoroutineCnt)
		// 启动消费者
		go d.consumer(d.goroutineListTable[goroutine.GetGoRoutineName()])
	}
	if _, exist := d.goroutineLockTable[goroutine.GetGoRoutineName()]; !exist {
		d.goroutineLockTable[goroutine.GetGoRoutineName()] = make(chan int, goroutine.GetMaxGoRoutineCnt())
	}
	d.goroutineListTable[goroutine.GetGoRoutineName()] <- goroutine
}

// GetGoroutineCount 获取 goroutine 数量信息
//
// Author : go_developer@163.com<张德满>
//
// Date : 2021/02/06 21:25:02
func (d *dispatch) GetGoroutineCount() map[string]GoroutineInfo {
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	d.lock.RLock()
	defer d.lock.RUnlock()

	result := make(map[string]GoroutineInfo)
	for name, lockChan := range d.goroutineLockTable {
		result[name] = GoroutineInfo{
			Time:            currentTime,
			GoroutineName:   name,
			MaxRunCount:     cap(lockChan),
			RunCount:        len(lockChan),
			WaitingCount:    len(d.goroutineListTable[name]),
			MaxWaitingCount: cap(d.goroutineListTable[name]),
		}
	}
	return result
}

// getGoroutineLock 获取协程锁
//
// Author : go_developer@163.com<张德满>
//
// Date : 2021/02/06 17:26:25
func (d *dispatch) getGoroutineLock(goroutine IGoRoutine) {
	d.lock.RLock()
	defer d.lock.RUnlock()
	d.goroutineLockTable[goroutine.GetGoRoutineName()] <- 1
}

// releaseGoroutineLock 释放锁
//
// Author : go_developer@163.com<张德满>
//
// Date : 2021/02/06 18:15:39
func (d *dispatch) releaseGoroutineLock(goroutine IGoRoutine) {
	d.lock.RLock()
	defer d.lock.RUnlock()
	<-d.goroutineLockTable[goroutine.GetGoRoutineName()]
}

// consumer 对一个协程任务启动消费者
//
// Author : go_developer@163.com<张德满>
//
// Date : 2021/02/06 18:03:43
func (d *dispatch) consumer(taskChannel chan IGoRoutine) {
	for goroutine := range taskChannel {
		// 获取锁
		d.getGoroutineLock(goroutine)

		timeout := goroutine.GetMaxExecuteTime()
		if timeout <= 0 {
			timeout = defaultGoroutineRuntime
		}
		ctx, cancelFunc := context.WithDeadline(context.Background(), time.Unix(timeout+time.Now().Unix(), 0))
		// 执行逻辑
		go func(ctx context.Context, cancelFunc context.CancelFunc, goroutine IGoRoutine) {
			startTime := time.Now().Unix()

			finishChan := make(chan int, 1)
			go func(ctx context.Context, finishChan chan int) {
				fmt.Println("进入了真实任务逻辑")
				// 析构函数
				defer goroutine.Recycle()
				// 释放锁
				defer d.releaseGoroutineLock(goroutine)
				// 捕获panic,防止因为协程异常导致进程挂掉
				defer func() {
					if r := recover(); nil != r {
						panicTrace := d.panicTrace()
						goroutine.PanicCallback(panicTrace)
					}
				}()
				//执行协程
				if err := goroutine.Execute(); nil != err {
					// 触发失败回调`
					goroutine.FailCallback(err)
				} else {
					// 触发成功回调
					goroutine.SuccessCallback()
				}
				finishChan <- 1
			}(ctx, finishChan)

			for {
				select {
				case <-ctx.Done():
					// 接收到信号
					fmt.Printf("%s, 接收到完成信号 : %d \n", time.Now().Format("2006-01-02 15:04:05"), time.Now().Unix()-startTime)
					cancelFunc()
					return
				case <-time.After(time.Second * time.Duration(timeout)):
					// 超时
					fmt.Printf("%s 超时自动结束 %d \n", time.Now().Format("2006-01-02 15:04:05"), time.Now().Unix()-startTime)
					cancelFunc()
					return
				case <-finishChan:
					// 任务预期时间内已经提前完成
					return
				}
			}
		}(ctx, cancelFunc, goroutine)
	}
}

// panicTrace 捕获执行的异常信息
//
// Author : go_developer@163.com<张德满>
//
// Date : 2021/02/06 18:25:35
func (d *dispatch) panicTrace() []byte {
	s := []byte("/src/runtime/panic.go")
	e := []byte("\ngoroutine ")
	line := []byte("\n")
	stack := make([]byte, 40960)
	length := runtime.Stack(stack, true)
	start := bytes.Index(stack, s)
	stack = stack[start:length]
	start = bytes.Index(stack, line) + 1
	stack = stack[start:]
	end := bytes.LastIndex(stack, line)
	if end != -1 {
		stack = stack[:end]
	}
	end = bytes.Index(stack, e)
	if end != -1 {
		stack = stack[:end]
	}
	stack = bytes.TrimRight(stack, "\n")
	return stack
}
