package base

import (
	"net"
	"net/rpc"

	"github.com/rlds/rbox/base/def"
	"github.com/rlds/rbox/base/util"
)

// RpcWorker rpc mode1
//    rpc模式接口
//    此模式下tcp为短链接方式使用
type RpcWorker struct{}

func newWorker() *RpcWorker {
	return &RpcWorker{}
}

// Call 执行请求
func (w *RpcWorker) Call(in def.RequestIn, hres *def.BoxOutPut) error {
	Log("call T:", in.TaskId, " F:", in.From, " C:", in.Call) //, " in:", in.Input)
	//Log("call T:", in.TaskId, " F:", in.From, " C:", in.Call, " in:", in.Input)
	// in.Input.TaskId = in.TaskId
	box.DoWork(in.TaskId, in.Input)
	*hres = box.Output(in.TaskId)
	hres.TaskId = in.TaskId
	return nil
}

// Status 获取状态
func (w *RpcWorker) Status(in def.RequestIn, hres *def.BoxOutPut) error {
	Log("status T:", in.TaskId, " F:", in.From, " C:", in.Call)
	*hres = box.Output(in.TaskId)
	hres.TaskId = in.TaskId
	return nil
}

// Ping 心跳监测
func (w *RpcWorker) Ping(in string, out *string) error {
	*out = "ok"
	return nil
}

type rpcModeWorker struct {
}

// Register 注册
func (r *rpcModeWorker) Register() {
	gbox.cfg.Mode = "rpc"
	//Log(gbox.cfg)
	gbox.cfg.BoxInfo.ModeInfo = checkRpcModeHost(gbox.cfg.BoxInfo.ModeInfo)
	//Log(gbox.cfg)
	res, err := RegBoxPost(gbox.cfg.ShowServerPath, util.ObjToStr(gbox.cfg.BoxInfo))
	if err == nil {
		Log("regok:", string(res))
	} else {
		Log("regerr:", err)
	}
}

// Run 启动服务
func (r *rpcModeWorker) Run() {
	Log("rpc mode start")
	gobRegister()

	rpc.Register(newWorker())
	host := gbox.cfg.SelfHttpServerHost

	tcpAddr, err := net.ResolveTCPAddr("tcp", host)
	if err != nil {
		Log("Error: listen ", host, " error:", err)
		return
	}
	l, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		Log("Error: listen ", host, " error:", err)
		return
	}

	// 开始注册
	r.Register()
	// 阻塞运行
	for {
		conn, err := l.Accept()
		if err != nil {
			continue
		}
		Log("ConnLocalAddr:", conn.RemoteAddr())
		go rpc.ServeConn(conn)
	}
}
