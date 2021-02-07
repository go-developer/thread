// Package thread ...
//
// Author: go_developer@163.com<张德满>
//
// Description: 管理go协程
//
// File: i_goroutine.go
//
// Version: 1.0.0
//
// Date: 2021/02/06 16:13:42
package thread

// IGoRoutine 定义接口约束
//
// Author : go_developer@163.com<张德满>
//
// Date : 2021/02/06 16:14:45
type IGoRoutine interface {
	// Execute 执行协程主逻辑
	//
	// Author : go_developer@163.com<张德满>
	//
	// Date : 2021/02/06 16:15:58
	Execute() error

	// FailCallback 执行失败时的回调(业务逻辑执行失败)
	//
	// Author : go_developer@163.com<张德满>
	//
	// Date : 2021/02/06 16:16:21
	FailCallback(err error)

	// SuccessCallback 执行成功时的回调
	//
	// Author : go_developer@163.com<张德满>
	//
	// Date : 2021/02/06 16:17:11
	SuccessCallback()

	// PanicCallback 逻辑发生panic时的回调
	//
	// Author : go_developer@163.com<张德满>
	//
	// Date : 2021/02/06 16:17:45
	PanicCallback(paniceTrace []byte)

	// GetMaxExecuteTime 获取协程最大执行时间,单位: 秒,设置成0表示不限时间
	//
	// Author : go_developer@163.com<张德满>
	//
	// Date : 2021/02/06 16:18:09
	GetMaxExecuteTime() int64

	// GetMaxGoRoutineCnt 获取当前协程允许同时调度的任务数量,设置成0代表不限制
	//
	// Author : go_developer@163.com<张德满>
	//
	// Date : 2021/02/06 16:29:44
	GetMaxGoRoutineCnt() int

	// GetGoRoutineName 获取goroutine任务名,需要全局唯一,不要冲突
	//
	// Author : go_developer@163.com<张德满>
	//
	// Date : 2021/02/06 16:33:18
	GetGoRoutineName() string
}
