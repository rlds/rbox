//  
//  fboxConfig.go
//  main
//
//  Created by 吴道睿 on 2018-11-14 17:51:14
//
package main

import (
	"github.com/rlds/rbox/base"
	."github.com/rlds/rbox/base/def"
	"os"
)

/*
	以下代码由rboxtool自动产生，
	若不清楚了解每条设置参数和含义不建议修改
*/
func init(){
	var cfg base.BoxConfig
	
	//工具类别 字符串（中文字符未测试）
	cfg.Group                  = `FboxGroup`
	
	//工具名称
	cfg.Name                   = `fbox`
	
	//这个工具展示为什么名字
	cfg.ShowName               = `f-工具`
	
	//作者
	cfg.Author                 = `吴道睿`
	
	//工具的描述
	cfg.Description            = `获取master所有相关日志`
	
	//http模式时开启的地址和端口
	cfg.SelfHttpServerHost     = ":9879"
	
    //操作界面访问的路径
    cfg.ModeInfo               = "http://127.0.0.1:9879"
    
	//界面操作层路径
	cfg.ShowServerPath         = "http://127.0.0.1:9888/regbox"
	
	//nats通信组件server地址列表 ',' 分割
	cfg.NatsServerList         = "" 
	

	//nats 接入用户名
	cfg.NatsServerUserName     = ""
	
	//nats 接入密码
	cfg.NatsServerUserPassword = ""
	
	//日志存储文件夹
	cfg.LogDir                 = "./log"
	
	//工具的输入参数描述 //附加输入参数数组
	cfg.Params = []BoxParam { 
		
		{
		Name : "taskid",
		Label : "任务Tid",
		Type : "text",
		Hint : "任务id",
		Reg : "",
		Value : "",
		Idx : "1"},
	    
	}
    
	cfg.Version = "1.0.1"
	
	err := base.SetBoxConfig(cfg)
	if err != nil {
		base.Log("box:",cfg.Name," init error:",err)
		os.Exit(1)
	}
}
