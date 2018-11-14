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
	//函数注册
	fbox.RegisterFunc(worker)
	//重设初始配置信息
	fbox.ResetConf(getConf())
	//开始执行
	fbox.Run()
}
/*
   输入是如下结构的map
   in = map[string]string{ {{range .BoxConf.Params}} "{{.Name}}":"", {{end}} }
	返回需要为文本,格式可以为markdown，json，html，imgstring=base64(bin(img))
	并指明结果的文本类型
*/
// 需要执行的函数体
func worker(in map[string]string)(rets string,restype string){
	//
	// todo :
	// 在下面添加功能代码
	// 这里是主体功能的入口和结果返回的地方
	// 

	// restype = "markdown|json|html|imgstring"
	return 
}