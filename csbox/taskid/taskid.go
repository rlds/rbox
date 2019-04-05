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

	"github.com/rlds/rbox/base/rhex"
)

// TaskID wid生成器 生成字符串id
type TaskID struct {
	mcd     *rhex.RHex64w
	bPidpre []byte
	lasth   int //小时时间值
	mu      sync.RWMutex
}

// NewTaskID 用于生成任务id
func NewTaskID() (w *TaskID) {
	w = new(TaskID)
	w.mcd = new(rhex.RHex64w)
	w.mcd.StrInit(6, []byte("000000"), false)
	w.bPidpre = make([]byte, 8)
	w.bPidpre[0] = 'w'
	w.bPidpre[1] = 's'
	w.setStart()
	go w.sysTimer()
	return
}

// GetTid 生成任务id
func (w *TaskID) GetTid() string {
	btid := w.bPidpre
	btid = append(btid, w.mcd.Add()...)
	return string(btid)
}

func (w *TaskID) setStart() {
	t := time.Now()
	y, m, d := t.Date() //天数值初始化
	h, mm, cs := t.Clock()
	w.bPidpre[2] = rhex.IntToNbyte(y - 2000)
	w.bPidpre[3] = rhex.IntToNbyte(int(m))
	w.bPidpre[4] = rhex.IntToNbyte(d)
	w.bPidpre[5] = rhex.IntToNbyte(h)
	w.bPidpre[6] = rhex.IntToNbyte(mm)
	w.bPidpre[7] = rhex.IntToNbyte(cs)
}

// 进行时间的处理
func (w *TaskID) sysTimer() {
	for {
		t := time.Now()
		h, m, cs := t.Clock()
		w.mu.RLock() //锁住
		if w.lasth != h {
			w.lasth = h
			if h == 0 { //天的数据清0
				y, m, d := t.Date()
				w.bPidpre[2] = rhex.IntToNbyte(y - 2000)
				w.bPidpre[3] = rhex.IntToNbyte(int(m))
				w.bPidpre[4] = rhex.IntToNbyte(d)
			}
		}
		//时分秒
		w.bPidpre[5] = rhex.IntToNbyte(h)
		w.bPidpre[6] = rhex.IntToNbyte(m)
		w.bPidpre[7] = rhex.IntToNbyte(cs)
		w.mu.RUnlock() //解锁
		time.Sleep(time.Second)
	}
}
