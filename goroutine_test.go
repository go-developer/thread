package thread

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"
	"time"
)

// TestGORoutine 单元测试
//
// Author : go_developer@163.com<张德满>
//
// Date : 2021/02/06 18:39:39
func TestGORoutine(t *testing.T) {
	for i := 0; i < 100; i++ {
		Dispatch.Run(NewTestGoroutine(i))
		byteData, _ := json.Marshal(Dispatch.GetGoroutineCount())
		fmt.Println("协程信息:", string(byteData))
	}
	for {
	}
}

// NewTestGoroutine ...
func NewTestGoroutine(num int) IGoRoutine {
	return &testGoroutine{
		num: num,
	}
}

type testGoroutine struct {
	num int
}

func (tg *testGoroutine) Execute() error {
	if (tg.num % 3) == 0 {
		return nil
	}
	if tg.num%3 == 1 {
		return errors.New("模拟业务失败")
	}
	if tg.num%3 == 2 {
		panic("模拟程序panic")
	}
	time.Sleep(time.Duration(tg.num % 10))
	return nil
}

func (tg *testGoroutine) FailCallback(err error) {
	fmt.Printf("%d : 协程执行失败,失败原因: %s \n", tg.num, err.Error())
}

func (tg *testGoroutine) SuccessCallback() {
	fmt.Printf("%d : 协程执行成功 \n", tg.num)
}

func (tg *testGoroutine) PanicCallback(data []byte) {
	fmt.Printf("%d : 协程执行异常 : %s \n", tg.num, string(data))
}

func (tg *testGoroutine) GetMaxExecuteTime() int64 {
	return 10
}

func (tg *testGoroutine) GetMaxGoRoutineCnt() int64 {
	return 10
}

func (tg *testGoroutine) GetGoRoutineName() string {
	return "test-task"
}
