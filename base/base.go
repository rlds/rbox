//
//  base.go
//  base
//
//  Created by 吴道睿 on 2018/4/6.
//  Copyright © 2018年 吴道睿. All rights reserved.
//
package base

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/rlds/rbox/base/def"
	. "github.com/rlds/rbox/base/def"
	"github.com/rlds/rbox/base/util"
	"github.com/rlds/rlog"
)

/*
   工具开发框架对外接口
*/
type (
	RBox = rbox
	rbox struct {
		box    Box       //具体的box
		worker Worker    //具体执行器
		cfg    BoxConfig //配置信息

		//每个任务进行编号
		taskid    string //任务id编号
		systime_u int64  //系统时间
	}

	//配置信息的定义结构
	BoxConfig struct {
		BoxInfo
		LogDir                 string //
		SelfHttpServerHost     string //http模式时开启的地址和端口
		ShowServerPath         string //界面操作层路径
		NatsServerList         string //nats通信组件server地址列表 ',' 分割多个
		NatsServerUserName     string //nats 接入用户名
		NatsServerUserPassword string //nats 接入密码
	}

	//具体工具的接口定义，每个工具的实现为此接口的实现
	Box interface {
		DoWork(taskid string, input InputData) (err error)
		Output(taskid string) BoxOutPut
	}

	BoxClient interface {
		Call(in def.RequestIn, hres *def.BoxOutPut) (err error)
		Status(in def.RequestIn, hres *def.BoxOutPut) (err error)
		Ping(in string, out *string) bool
		Close() error
	}
	//执行器
	Worker interface {
		//注册至webserver
		Register()
		Run()
	}
)

const (
	OutputType_Markdown     = "markdown"
	OutputType_Json_default = "json;default"
	OutputType_Html         = "html"
	OutputRetuen_Success    = "0"
	OutputRetuen_Error      = "-1"

	ModeType_HTTP    = "http"    //http 模式
	ModeType_Command = "command" //command 命令行模式
	ModeType_Nats    = "nats"    //nats 模式
	ModeType_Rpc     = "rpc"     //rpc 模式
)

var (
	gbox       rbox
	box        Box
	boxmode    string
	ModTypeErr = errors.New("mode type error")
)

func (m *rbox) setMode(mode, subbox, input string) (err error) {
	boxmode = mode
	switch mode {
	case ModeType_Command:
		{
			isCommand = true
			var cw commandModeWorker
			cw.input.SubBoxName = subbox
			if len(input) > 3 { //允许 input 为空，但不能为错误的json
				cw.input.Data = make(map[string]interface{})
				err = json.Unmarshal([]byte(input), &cw.input.Data)
				if err != nil {
					return
				}
			}
			gbox.worker = &cw
		}
	case ModeType_HTTP:
		{
			isCommand = false
			var hw httpModeWorker
			gbox.worker = &hw
		}
	case ModeType_Nats:
		{
			isCommand = false
			var nw natsModeWorker
			gbox.worker = &nw
		}
	case ModeType_Rpc:
		{
			isCommand = false
			var rp rpcModeWorker
			gbox.worker = &rp
		}
	default:
		{
			err = ModTypeErr
		}
	}
	return
}

//返回当前执行模式名称和是否命令行模式
func GetRunMode() (string, bool) {
	return boxmode, isCommand
}

//设置工具类别名称和登录参数信息
func SetBoxConfig(cfg BoxConfig) (err error) {
	gbox.cfg = cfg
	gbox.cfg.Mode = boxmode
	err = gbox.checkcfg()
	return
}

//配置参数的检查
func (m *rbox) checkcfg() error {
	//检查配置

	//检查日志

	return nil
}

func ParamToMapEg(cfg BoxConfig) string {
	res := ""
	for _, sub := range cfg.SubBox {
		ret := []string{}
		for _, pam := range sub.Params {
			ret = append(ret, `"`+pam.Name+`":`+pam.ValueType)
		}
		if len(sub.SubName) > 0 {
			res += " \n\t -subbox " + sub.SubName + " -input '{" + strings.Join(ret, ",") + "}'"
		} else {
			res += "{" + strings.Join(ret, ",") + "}"
		}
	}
	return res
}

func paramSubBox(sub BoxConfig) string {
	ret := []string{}
	for _, nm := range sub.SubBox {
		ret = append(ret, nm.SubName)
	}
	return strings.Join(ret, "|")
}

// Init 初始化
func Init() {
	//输出设置信息
	var (
		mode   string
		input  string
		logdir string
		subbox string
	)
	flag.StringVar(&logdir, "log", gbox.cfg.LogDir, "日志输出文件夹路径")
	flag.StringVar(&mode, "mode", "", "运行模式(http,command,nats) eg: -mode http")
	flag.StringVar(&subbox, "subbox", "", "子功能模块名称 eg: \n\t -subbox "+paramSubBox(gbox.cfg))
	flag.StringVar(&input, "input", "", "输入参数信息json格式 注意值类型需要一致 eg: "+ParamToMapEg(gbox.cfg)+"")
	flag.Parse()
	//模式设置
	err := gbox.setMode(mode, subbox, input)
	if err != nil {
		Log("init error:", err)
		os.Exit(2)
	}
	// 服务模式启动日志
	if !isCommand {
		// 日志文件路径以输入为准
		gbox.cfg.LogDir = logdir
		util.TestAndCreateDir(gbox.cfg.LogDir)
		rlog.LogInit(3, gbox.cfg.LogDir, MaxLogLen_m, 1)
	}
}

//注册工具组件
func RegisterBox(b Box) {
	box = b
}

//执行接口
func Run() {
	gbox.worker.Run()
}

var modeClientError = fmt.Errorf("mode client error")

// NewBoxClient 创建box客户端
func NewBoxClient(box *def.BoxInfo) (bc BoxClient, err error) {
	switch box.Mode {
	case "http":
		{
			bc, err = NewHTTPClient(box)
		}
	case "rpc":
		{
			bc, err = NewRpcClient(box)
		}
	// case "nats":
	// 	{
	// 		bc = nil
	// 	}
	default:
		{
			err = fmt.Errorf(box.Mode + " mode client error")
		}
	}
	return
}
