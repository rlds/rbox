//
//  fmtPrntTest.go
//  mbox
//
//  Created by 吴道睿 on 2018/5/26.
//  Copyright © 2018年 吴道睿. All rights reserved.
//
package main

import(
    "fmt"
)

func main(){
    //fmt.Println = Log
	
	fmt_f := &fmt.Println
	log_f := &Log
	
	fmt_f = &Log
	
	fmt.Println("cve")
	
	log_f("logf","dd")
}

func Log(a ...interface{}) (n int, err error){
	fmt.Println("Log",a...)
	return
}
