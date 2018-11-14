//
//  boxWorker.go
//  rbox
//
//  Created by 吴道睿 on 2018/5/13.
//  Copyright © 2018年 吴道睿. All rights reserved.
//
package main

import (
	"encoding/json"
	"rlds/mbox/boxdef"
	"rlds/mbox/util"
)

/*
   开始执行任务
*/
func (b *boxInfo) DoWork(indat boxdef.RequestIn) (rt boxdef.BoxOutPut) {
	if b.isAlive {
		switch b.connType {
		case "http":
			{
				rt = b.httpMode(indat)
			}
		case "nats":
			{
				rt = b.natsMode(indat)
			}
		default:
			{
				rt.Type = CallBox_ResType_CallErr
				rt.Code = CallBox_Code_InputErr_103
				rt.ReturnMsg = "connType err"
			}
		}
	} else {
		rt.Type = CallBox_ResType_CallErr
		rt.Code = CallBox_Code_InputErr_104
		rt.ReturnMsg = "not alive"
	}
	return
}

/*
 开始执行任务
*/
func (b *boxInfo) TaskRes(indat boxdef.RequestIn) (rt boxdef.BoxOutPut) {
	if b.isAlive {
		switch b.connType {
		case "http":
			{
				rt = b.httpModeTaskRes(indat)
			}
		case "nats":
			{
				rt = b.natsModeTaskRes(indat)
			}
		default:
			{
				rt.Type = CallBox_ResType_CallErr
				rt.Code = CallBox_Code_InputErr_103
				rt.ReturnMsg = "connType err"
			}
		}
	} else {
		rt.Type = CallBox_ResType_CallErr
		rt.Code = CallBox_Code_InputErr_104
		rt.ReturnMsg = "not alive"
	}
	return
}

//http模式访问执行
func (b *boxInfo) httpMode(indat boxdef.RequestIn) (rt boxdef.BoxOutPut) {
	urlpath := b.ModeInfo + "/call/" + b.Group + "/" + b.Name
	println(urlpath)
	retb, err := HttpPostJson(urlpath, util.ObjToStr(indat))
	if err == nil {
		err = json.Unmarshal(retb, &rt)
		if err != nil {
			rt.Type = CallBox_ResType_BoxRetErr
			rt.Code = CallBox__Code_BoxRetErr_110
			rt.ReturnMsg = "BoxRetError"
		}

		//成功返回  不修改返回状态
		//rt.Type      = BoxOutPut_CallBox_Ok
		//rt.Code      = BoxOutPut_Code_Ok

	} else {
		rt.Type = CallBox_ResType_BoxRetErr
		rt.Code = CallBox__Code_BoxRetErr_110
		rt.ReturnMsg = "BoxRetError"
	}
	return
}

//http模式访问执行
func (b *boxInfo) httpModeTaskRes(indat boxdef.RequestIn) (rt boxdef.BoxOutPut) {
	urlpath := b.ModeInfo + "/taskRes/" + b.Group + "/" + b.Name
	println(urlpath)
	retb, err := HttpPostJson(urlpath, util.ObjToStr(indat))
	if err == nil {
		err = json.Unmarshal(retb, &rt)
		if err != nil {
			rt.Type = CallBox_ResType_BoxRetErr
			rt.Code = CallBox__Code_BoxRetErr_110
			rt.ReturnMsg = "BoxRetError"
		}

		//成功返回  不修改返回状态
		//rt.Type      = BoxOutPut_CallBox_Ok
		//rt.Code      = BoxOutPut_Code_Ok

	} else {
		rt.Type = CallBox_ResType_BoxRetErr
		rt.Code = CallBox__Code_BoxRetErr_110
		rt.ReturnMsg = "BoxRetError"
	}
	return
}

//nats 模式访问执行
func (b *boxInfo) natsMode(indat boxdef.RequestIn) (rt boxdef.BoxOutPut) {
	natstopic := b.ModeInfo + "." + b.Group + "." + b.Name
	//
	println(natstopic)
	return
}

//nats 模式访问执行
func (b *boxInfo) natsModeTaskRes(indat boxdef.RequestIn) (rt boxdef.BoxOutPut) {
	natstopic := b.ModeInfo + "." + b.Group + "." + b.Name
	//

	println(natstopic)
	return
}

//连接检测
func (b *boxInfo) conn() bool {
	return true
}
