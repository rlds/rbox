//
//  boxHbeat.go
//  rbox
//
//  Created by 吴道睿 on 2018/4/12.
//  Copyright © 2018年 吴道睿. All rights reserved.
//
package main

// Whbeat 健康检测
type Whbeat struct {
}

// NewWhbeat 心跳检测
func NewWhbeat() *Whbeat {
	return &Whbeat{}
}

//AddGroup 放入需要进行心跳检查的box指针
func (w *Whbeat) AddGroup(bg *boxGroupInfo) (err error) {
	return
}

//CheckStart 心跳检查开始
func (w *Whbeat) CheckStart() {
	return
}
