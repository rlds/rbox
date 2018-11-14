//
//  runModeCommand.go
//  base
//
//  Created by 吴道睿 on 2018/4/6.
//  Copyright © 2018年 吴道睿. All rights reserved.
//
package base

/*
	命令行模式的执行
*/
type commandModeWorker struct {
	input map[string]string
}

func (c *commandModeWorker) Register() {
	Log("------------------------------------")
	Log("工具名称 :", gbox.cfg.Group, gbox.cfg.Name)
	Log("作者信息 :", gbox.cfg.Author)
	Log("执行模式 :", gbox.cfg.Mode)
	Log("功能描述 :", gbox.cfg.Description)
	Log("------------------------------------\n")
}

//命令行功能执行
func (c *commandModeWorker) Run() {
	c.Register()
	task_id := "commandtask"
	box.DoWork(task_id, c.input)
	res := box.Output(task_id)
	if res.Code == OutputRetuen_Success {
		Log("返回结果格式:", res.Type)
		Log("命令执行成功结果如下:\n")
		if res.Data != nil {
			Log(res.Data.(string))
		}
	} else {
		Log("命令执行错误:")
		Log("错误代码：", res.Code)
		Log("错误提示：", res.ReturnMsg)
	}
	return
}
