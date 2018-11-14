//  
//  GetTidLogsMain.go
//  main
//
//  Created by 金山 on 2017-09-27 15:58:58
//
package main

import(
	"wacaispider/mbox"
)

func main(){
	var box GetTidLogsBox
	mbox.Init()
	mbox.RegisterBox(&box)
	mbox.Run()
}
