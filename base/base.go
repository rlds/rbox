//
//  base.go
//  base
//
//  Created by 吴道睿 on 2018/4/6.
//  Copyright © 2018年 吴道睿. All rights reserved.
//
package base

import(
	"flag"
	"os"
	"errors"
	"encoding/json"
	"github.com/rlds/rlog"
	"github.com/rlds/rbox/base/util"
	."github.com/rlds/rbox/base/def"
)

/*
   工具开发框架对外接口
 */
type (
	RBox = rbox  
	rbox struct{
		box        Box      //具体的box
		worker     Worker   //具体执行器
		
		cfg        BoxConfig //配置信息
		
		//每个任务进行编号
		taskid     string   //任务id编号
		systime_u  int64    //系统时间
	}
    
	
	//配置信息的定义结构
	BoxConfig struct{
		BoxInfo
		LogDir                 string      //
		SelfHttpServerHost     string      //http模式时开启的地址和端口
		ShowServerPath         string      //界面操作层路径
		NatsServerList         string      //nats通信组件server地址列表 ',' 分割多个
		NatsServerUserName     string      //nats 接入用户名
		NatsServerUserPassword string      //nats 接入密码
	}
	
	
	//具体工具的接口定义，每个工具的实现为此接口的实现
	Box interface{
		DoWork(taskid string,input map[string]string)(err error)
		Output(taskid string)(BoxOutPut)
	}

	//真正要实现的接口
	/*
	BoxWorker interface{
		DoWork(taskid string,input map[string]string)(err error)
		Output(taskid string)(BoxOutPut)
	}
	*/
	//执行器
	Worker interface{
		//注册至webserver
		Register()
		Run()
	}
)

const(
    OutputType_Markdown     = "markdown"
	OutputType_Json_default = "json;default"
	OutputType_Html         = "html"
	OutputRetuen_Success    = "0"
	OutputRetuen_Error      = "-1"

	ModeType_HTTP    = "http"    //http 模式
	ModeType_Command = "command" //command 命令行模式
	ModeType_Nats    = "nats"    //nats 模式
)

var (
	gbox rbox
	box  Box
	boxmode string
	ModTypeErr = errors.New("mode type error")
)

func (m *rbox)setMode(mode,input string)(err error){
	boxmode = mode
	switch mode {
		case ModeType_Command:{
			isCommand = true
			var cw commandModeWorker 
			cw.input = make(map[string]string)
			err = json.Unmarshal([]byte(input),&cw.input)
			if err != nil {
				return
			}
			gbox.worker = &cw
		}
		case ModeType_HTTP:{
			isCommand = false
			var hw httpModeWorker 
			gbox.worker = &hw
		}
		case ModeType_Nats:{
			isCommand = false
			var nw natsModeWorker 
			gbox.worker = &nw
		}
		default:{
			err = ModTypeErr
		}
	}
	return
}

//返回当前执行模式名称和是否命令行模式
func GetRunMode()(string,bool){
	return boxmode,isCommand
}

//设置工具类别名称和登录参数信息
func SetBoxConfig(cfg BoxConfig)(err error){
	gbox.cfg      = cfg
	gbox.cfg.Mode = boxmode
	err = gbox.checkcfg()
	return
}

//配置参数的检查
func (m *rbox)checkcfg()error{
	//检查配置
	
	//检查日志
	
	return nil
}

// Init 初始化
func Init(){
	var (
	     mode  string
		 input string
	)
	flag.StringVar(&mode,"mode","","运行模式(http,command,nats)")
	flag.StringVar(&input,"input","","输入参数信息json格式")
	flag.Parse()
	//模式设置
	err := gbox.setMode(mode,input)
	if err != nil {
		Log("init error:",err)
		os.Exit(2)
	}
	// 服务模式启动日志
	if !isCommand {
		util.TestAndCreateDir(gbox.cfg.LogDir)
		rlog.LogInit(3, gbox.cfg.LogDir , MaxLogLen_m , 1)
	}
}

//注册工具组件
func RegisterBox(b Box){
	box = b
}

//执行接口
func Run(){
	gbox.worker.Run()
}

