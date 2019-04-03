//
//  boxWorker.go
//  rbox
//
//  Created by 吴道睿 on 2018/5/13.
//  Copyright © 2018年 吴道睿. All rights reserved.
//
package main

import (
	"github.com/rlds/rbox/base"
	"github.com/rlds/rbox/base/def"
)

func (b *boxInfo) Init() (err error) {
	if b.isAlive {
		b.boxClient, err = base.NewBoxClient(&b.BoxInfo)
	}
	return
}

/*
   开始执行任务
*/
func (b *boxInfo) DoWork(indat def.RequestIn) (rt def.BoxOutPut) {
	var err error
	if b.isAlive {
		switch b.connType {
		case "http", "rpc":
			{
				err = b.boxClient.Call(indat, &rt)
			}
		case "nats":
			{
				rt = b.natsMode(indat)
			}
		default:
			{
				rt.Type = CallBoxResTypeCallErr
				rt.Code = CallBoxCodeInputErr103
				rt.ReturnMsg = "connType err"
				rt.Status = "COMPLETE"
			}
		}
	} else {
		rt.Type = CallBoxResTypeCallErr
		rt.Code = CallBoxCodeInputErr104
		rt.ReturnMsg = "not alive"
		rt.Status = "COMPLETE"
	}
	if err != nil {
		rt.Type = CallBoxResTypeBoxRetErr
		rt.Code = CallBoxCodeBoxRetErr110
		rt.ReturnMsg = "BoxRetError"
		rt.Status = "COMPLETE"
	}
	return
}

/*
 开始执行任务
*/
func (b *boxInfo) TaskRes(indat def.RequestIn) (rt def.BoxOutPut) {
	indat.Input.IsSync = b.IsSync
	var err error
	if b.isAlive {
		switch b.connType {
		case "http", "rpc":
			{
				err = b.boxClient.Status(indat, &rt)
			}
		case "nats":
			{
				rt = b.natsModeTaskRes(indat)
			}
		default:
			{
				rt.Type = CallBoxResTypeCallErr
				rt.Code = CallBoxCodeInputErr103
				rt.ReturnMsg = "connType err"
			}
		}
	} else {
		rt.Type = CallBoxResTypeCallErr
		rt.Code = CallBoxCodeInputErr104
		rt.ReturnMsg = "not alive"
		rt.Status = "COMPLETE"
	}
	if err != nil {
		rt.Type = CallBoxResTypeBoxRetErr
		rt.Code = CallBoxCodeBoxRetErr110
		rt.ReturnMsg = "BoxRetError"
		rt.Status = "COMPLETE"
	}
	return
}

//nats 模式访问执行
func (b *boxInfo) natsMode(indat def.RequestIn) (rt def.BoxOutPut) {
	//natstopic := b.ModeInfo + "." + b.Group + "." + b.Name
	//
	//println(natstopic)
	return
}

//nats 模式访问执行
func (b *boxInfo) natsModeTaskRes(indat def.RequestIn) (rt def.BoxOutPut) {
	//natstopic := b.ModeInfo + "." + b.Group + "." + b.Name
	//
	//println(natstopic)
	return
}
