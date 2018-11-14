package main

import "rlds/mbox"

// Woker 任务执行实体
type Woker struct{}

// Run 任务执行体
func (w *Woker) Run(in map[string]string) (resdata, datatype string) {
	mbox.Log("indata:", in)
	resdata = "*  return test"
	datatype = "MarkDown"
	return
}
