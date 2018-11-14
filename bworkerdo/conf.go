package main

import (
	"rlds/mbox"
	"rlds/mbox/boxdef"
)

func setConf() (cfg mbox.BoxConfig) {
	//工具类别
	cfg.Group = "master"

	//工具名称
	cfg.Name = `bworker`

	//这个工具展示为什么名字
	cfg.ShowName = `tid日志查询`

	//作者
	cfg.Author = `wdr`

	//工具的描述
	cfg.Description = `获取master所有相关日志`

	//http模式时开启的地址和端口
	cfg.SelfHttpServerHost = ":9879"

	//操作界面访问的路径
	cfg.ModeInfo = "http://192.168.4.225:9879"

	//界面操作层路径
	cfg.ShowServerPath = "http://172.16.48.207:9888/regbox"

	//nats通信组件server地址列表 ',' 分割
	cfg.NatsServerList = ""

	//nats 接入用户名
	cfg.NatsServerUserName = ""

	//nats 接入密码
	cfg.NatsServerUserPassword = ""

	//日志存储文件夹
	cfg.LogDir = "./log"

	//工具的输入参数描述 //附加输入参数数组
	cfg.Params = []boxdef.BoxParam{
		{
			Name:  "tid",
			Label: "任务Tid",
			Type:  "text",
			Hint:  "17位任务id",
			Reg:   "",
			Value: "",
			Idx:   "1"},
	}
	cfg.Version = "1.0.1"
	return
}
