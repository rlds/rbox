//
//  bworkerConfig.go
//  main
//
//  Created by wdr on 2018-09-04 10:39:31
//
package bworker

import (
	"os"
	"rlds/mbox"
)

/*
	以下代码由boxtool自动产生，
	若不清楚了解每条设置参数和含义不建议修改
*/
func Init(cfg mbox.BoxConfig) {
	//var cfg mbox.BoxConfig

	err := mbox.SetBoxConfig(cfg)
	if err != nil {
		mbox.Log("box:", cfg.Name, " init error:", err)
		os.Exit(1)
	}
}
