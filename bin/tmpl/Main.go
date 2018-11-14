//  
//  {{.BoxConf.Name}}Main.go
//  main
//
//  Created by {{.BoxConf.Author}} on {{.Time}}
//
package main

import(
	"{{.MboxPackagePath}}"
)

func main(){
	var box {{.BoxConf.Name}}Box
	mbox.Init()
	mbox.RegisterBox(&box)
	mbox.Run()
}
