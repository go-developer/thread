// Package thread ...
//
// Description : 基础常量定义
//
// Author : go_developer@163.com<张德满>
//
// Date : 2021-02-07 11:46 上午
package thread

const (
	// 默认的运行的协程数量
	defaultGoroutineCount = 2048
	// 默认超时时间 30s
	defaultGoroutineRuntime = 30
)

// GoroutineInfo 对外输出的协程信息
type GoroutineInfo struct {
	GoroutineName   string `json:"goroutine_name"`    // 协程名称
	MaxRunCount     int    `json:"max_run_count"`     // 最大同时运行数量
	WaitingCount    int    `json:"waiting_count"`     // 当前等待数量
	MaxWaitingCount int    `json:"max_waiting_count"` // 最大等待任务数
	RunCount        int    `json:"run_count"`         // 当前运行数量
}
