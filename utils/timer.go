package utils

// 定时器封装

import (
	"time"
)

// 实例化存储每个具体定时器的相关内容
type TimerSt struct {
	Fun Timer
	Iterval int64 // seconds
	Start time.Time // 启动时间
	ExitFlag chan bool // 定时器停止标记
	Status int // 0 stop 1 start 2 run
}

// 用于实现该接口的调用Run执行定时内容
type Timer interface {
	Run(<- chan bool) // 参数channel接收到数据时该定时任务退出
}

type timerService struct {
	fun map[string]*TimerSt
}

var g_timerService *timerService

// start timer service
func TimerStartService() {
	go func() {
		for {
			time.Sleep(time.Nanosecond * 5)
			// 定时检查是否有新的任务加入
			
		}
	}()
}

// init timer
func TimerInit() {
	
}

// start single timer

// add timer
// stop single timer
// delete single timer
// stop all

