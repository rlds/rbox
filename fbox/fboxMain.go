//  
//  fboxMain.go
//  fbox
//
//  Created by 吴道睿 on 2018-11-14 17:51:14
//
package fbox

import(
	"github.com/rlds/rbox/base"
)

func Run(){
	var box fboxBox
	base.Init()
	box.Init()
	base.RegisterBox(&box)
	base.Run()
}
