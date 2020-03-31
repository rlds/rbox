//
//  fboxMain.go
//  fbox
//
//  Created by 吴道睿 on 2018-11-14 17:51:14
//
package fbox

import (
	"github.com/rlds/rbox/base"
	"github.com/rlds/rbox/base/def"
)

// 最小存储时间
const minTaskResultStoreTime = 10
const maxTaskResultStoreTime = 60 * 60

// SetMaxTaskResultStoreTime 设置全局任务结果的最大存储时间
// 需要在 Run之前调用
// 设置值需要在 (10,3600) 区间内
func SetMaxTaskResultStoreTime(tms int64) {
	if tms < minTaskResultStoreTime || tms > maxTaskResultStoreTime {
		return
	}
	_CleanTaskTimeStep = (tms)/2 - 1
	_TaskInfoStoreTimeStep = tms
}

// Run 开始执行
// 原则上一个进程是可以执行多个的
func Run() {
	var box fboxBox
	base.Init()
	box.Init()
	base.RegisterBox(&box)
	base.Run()
}

// RegisterFunc 注册执行函数
func RegisterFunc(f TFunc) {
	taskFunc = f
}

// TFunc 外部注册函数的定义
type TFunc func(def.InputData) (string, interface{})
