package tid

import (
	"sync"
	"time"

	"github.com/rlds/rbox/base/rhex"
	"github.com/rlds/rbox/base/util"
)

// TaskID wid生成器 生成字符串id
type TaskID struct {
	mcd        *rhex.RHex64w
	bPidpre    []byte
	lasth      int //小时时间值
	weeknum    int
	weeknumstr string
	mu         sync.RWMutex
}

// NewTaskID 用于生成任务id
func NewTaskID(prec string) (w *TaskID) {
	pre := []byte(prec)
	w = new(TaskID)
	w.mcd = new(rhex.RHex64w)
	w.mcd.StrInit(6, []byte("01se80"), false)
	w.bPidpre = make([]byte, 8)
	w.bPidpre[0] = pre[0]
	w.bPidpre[1] = pre[1]
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
		w.mu.Lock() //锁住
		if w.lasth != h {
			w.lasth = h
			if h == 0 { //天的数据清0
				y, m, d := t.Date()
				w.bPidpre[2] = rhex.IntToNbyte(y - 2000)
				w.bPidpre[3] = rhex.IntToNbyte(int(m))
				w.bPidpre[4] = rhex.IntToNbyte(d)
				_, wn := t.ISOWeek()
				w.weeknum = y*100 + wn
				w.weeknumstr = util.IntToStr(w.weeknum)
			}
		}
		//时分秒
		w.bPidpre[5] = rhex.IntToNbyte(h)
		w.bPidpre[6] = rhex.IntToNbyte(m)
		w.bPidpre[7] = rhex.IntToNbyte(cs)
		w.mu.Unlock() //解锁
		time.Sleep(time.Second)
	}
}

func (w *TaskID) GetWeekNum() int {
	return w.weeknum
}

func (w *TaskID) GetWeekStr() string {
	return w.weeknumstr
}
