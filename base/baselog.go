//
//  baselog.go
//  base
//
//  Created by 吴道睿 on 2018/4/6.
//  Copyright © 2018年 吴道睿. All rights reserved.
//
package base

import (
	"fmt"

	"github.com/rlds/rlog"
)

/*
   统一的日志输出
*/
const (
	MaxLogLenM = 1024 * 1024 * 1000
)

var isCommand bool

// Log 日志输出
func Log(args ...interface{}) {
	if isCommand {
		fmt.Println(args...)
	} else {
		rlog.V(1).Info(args...)
	}
}
