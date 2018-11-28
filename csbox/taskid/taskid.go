package taskid

//
//  taskid.go
//  rbox
//
//  Created by 吴道睿 on 2018/4/12.
//  Copyright © 2018年 吴道睿. All rights reserved.
//

import (
	"sync"
	"time"
)

// TaskID wid生成器 生成字符串id
type TaskID struct {
	mcd     *Hex64_w
	bPidpre []byte
	lasth   int //小时时间值
	mu      sync.RWMutex
}

// NewTaskID 用于生成任务id
func NewTaskID() (w *TaskID) {
	w = new(TaskID)
	w.mcd = new(Hex64_w)
	w.mcd.StrInit(6, []byte("000000"), false)
	w.bPidpre = make([]byte, 8)
	w.bPidpre[0] = 'w'
	w.bPidpre[1] = 's'
	go w.sysTimer()
	return
}

// GetTid 生成任务id
func (w *TaskID) GetTid() string {
	btid := w.bPidpre
	btid = append(btid, w.mcd.Add()...)
	return string(btid)
}

// 进行时间的处理
func (w *TaskID) sysTimer() {
	y, m, d := time.Now().Date() //天数值初始化
	w.bPidpre[2] = IntToNbyte(y - 2000)
	w.bPidpre[3] = IntToNbyte(int(m))
	w.bPidpre[4] = IntToNbyte(d)
	for {
		t := time.Now()
		h, m, cs := t.Clock()
		w.mu.RLock() //锁住
		if w.lasth != h {
			w.lasth = h
			if h == 0 { //天的数据清0
				y, m, d := t.Date()
				w.bPidpre[2] = IntToNbyte(y - 2000)
				w.bPidpre[3] = IntToNbyte(int(m))
				w.bPidpre[4] = IntToNbyte(d)
			}
		}
		//时分秒
		w.bPidpre[5] = IntToNbyte(h)
		w.bPidpre[6] = IntToNbyte(m)
		w.bPidpre[7] = IntToNbyte(cs)
		w.mu.RUnlock() //解锁
		time.Sleep(time.Second)
	}
}
