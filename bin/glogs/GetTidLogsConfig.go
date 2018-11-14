//  
//  GetTidLogsConfig.go
//  main
//
//  Created by 金山 on 2017-09-27 15:58:58
//
package main

import (
	"wacaispider/mbox"
	"os"
	."wacaispider/mbox/boxdef"
)

/*
	以下代码由boxtool自动产生，
	若不清楚了解每条设置参数和含义不建议修改
*/
func init(){
	var cfg mbox.BoxConfig
	
	//工具类别
	cfg.Group                  = "master"
	
	//工具名称
	cfg.Name                   = "GetTidLogs"
	
	//这个工具展示为什么名字
	cfg.ShowName               = "tid日志查询"
	
	//作者
	cfg.Author                 = "金山"
	
	//工具的描述
	cfg.Description            = "获取master所有相关日志"
	
	//http模式时开启的地址和端口
	cfg.SelfHttpServerHost     = ":9876"
	
	//界面操作层路径
	cfg.ShowServerPath         = "http://localhost:9888/regbox"
	
	//nats通信组件server地址列表 ',' 分割
	cfg.NatsServerList         = "" 
	
	//操作界面访问的路径
	cfg.ModeInfo               = "http://localhost:9876"
	
	//nats 接入用户名
	cfg.NatsServerUserName     = ""
	
	//nats 接入密码
	cfg.NatsServerUserPassword = ""
	
	//日志存储文件夹
	cfg.LogDir                 = "./log"
	
	//工具的输入参数描述 //附加输入参数数组
	cfg.Params = []BoxParam { {
		
		Name : "tid",
		Label : "任务Tid",
		Type : "text",
		Hint : "17位任务id",
		Reg : "",
		Value : "",
		Idx : "1"},
	    
	}
	
	cfg.Version = "1.0.1"
	
	err := mbox.SetBoxConfig(cfg)
	if err != nil {
		mbox.Log("box:",cfg.Name," init error:",err)
		os.Exit(1)
	}
}
