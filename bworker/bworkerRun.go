//
//  bworkerMain.go
//  main
//
//  Created by wdr on 2018-09-04 10:39:31
//
package bworker

import (
	"rlds/mbox"
)

func Run(cfg mbox.BoxConfig, worker Bworker) {
	bworker = worker
	Init(cfg)
	var box bworkerBox
	mbox.Init()
	box.Init()
	mbox.RegisterBox(&box)
	mbox.Run()
}
