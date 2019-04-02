//
//  fboxConfig.go
//  fbox
//
//  Created by 吴道睿 on 2018-11-14 17:51:14
//
package fbox

import (
	"os"

	"github.com/rlds/rbox/base"
	// . "github.com/rlds/rbox/base/def"
)

// 变量值检测替换
func reSet(src *string, dest string, deft string) {
	if dest != "" {
		*src = dest
	} else {
		*src = deft
	}
}

/*
	以下代码由rboxtool自动产生，
	若不清楚了解每条设置参数和含义不建议修改
*/
func ResetConf(icfg base.BoxConfig) {
	var cfg base.BoxConfig

	//工具类别 字符串（中文字符未测试）
	reSet(&cfg.Group, icfg.Group, `FboxGroup`)

	//工具名称
	reSet(&cfg.Name, icfg.Name, `fbox`)

	//这个工具展示为什么名字
	reSet(&cfg.ShowName, icfg.ShowName, `f-工具`)

	//作者
	reSet(&cfg.Author, icfg.Author, `吴道睿`)

	//工具的描述
	reSet(&cfg.Description, icfg.Description, `获取master所有相关日志`)

	//http模式时开启的地址和端口
	reSet(&cfg.SelfHttpServerHost, icfg.SelfHttpServerHost, ":9879")

	//操作界面访问的路径
	reSet(&cfg.ModeInfo, icfg.ModeInfo, "http://127.0.0.1:9879")

	//界面操作层路径
	reSet(&cfg.ShowServerPath, icfg.ShowServerPath, "http://127.0.0.1:9888/regbox")

	//nats通信组件server地址列表 ',' 分割
	reSet(&cfg.NatsServerList, icfg.NatsServerList, "")

	//nats 接入用户名
	reSet(&cfg.NatsServerUserName, icfg.NatsServerUserName, "")

	//nats 接入密码
	reSet(&cfg.NatsServerUserPassword, icfg.NatsServerUserPassword, "")

	//日志存储文件夹
	reSet(&cfg.LogDir, icfg.LogDir, "./log")

	//工具的输入参数描述 //附加输入参数数组
	if len(icfg.SubBox) > 0 {
		cfg.SubBox = icfg.SubBox
	}
	// else {
	// 	cfg.SubBox = []BoxParam{
	// 		{
	// 			Name:  "taskid",
	// 			Label: "任务Tid",
	// 			Type:  "text",
	// 			Hint:  "任务id",
	// 			Reg:   "",
	// 			Value: "",
	// 			Idx:   "1"},
	// 	}
	// }

	reSet(&cfg.Version, icfg.Version, "1.0.1")

	cfg.IsSync = icfg.IsSync

	err := base.SetBoxConfig(cfg)
	if err != nil {
		base.Log("box:", cfg.Name, " init error:", err)
		os.Exit(1)
	}
}
