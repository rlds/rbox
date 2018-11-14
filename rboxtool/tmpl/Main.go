//  
//  {{.BoxConf.Name}}Main.go
//  main
//
//  Created by {{.BoxConf.Author}} on {{.Time}}
//
package main

import(
	"github.com/rlds/rbox/base"
)

func main(){
	var box {{.BoxConf.Name}}Box
	base.Init()
	box.Init()
	base.RegisterBox(&box)
	base.Run()
}
