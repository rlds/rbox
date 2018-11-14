//  
//  fboxMain.go
//  main
//
//  Created by 吴道睿 on 2018-11-14 17:51:14
//
package main

import(
	"github.com/rlds/rbox/fbox"
)

func main(){
	fbox.RegisterFunc(worker)
	fbox.Run()
}


func worker(in map[string]string)string{
	fmt.Println(in)
	return "* worker Test"
}