package base

import (
	"fmt"
	"net"
	"net/rpc"
	"time"

	"github.com/rlds/rbox/base/def"
)

// BoxRpcClient 客户端信息
type BoxRpcClient struct {
	box       *def.BoxInfo
	rpcClient *rpc.Client
}

// NewRpcClient 客户端初始化
func NewRpcClient(box *def.BoxInfo) (bclt *BoxRpcClient, err error) {
	bclt = new(BoxRpcClient)
	bclt.box = box
	err = bclt.init()
	return
}

func (box *BoxRpcClient) init() error {
	gobRegister()
	Log("start dial", box.box.ModeInfo)
	conn, err := net.DialTimeout("tcp", box.box.ModeInfo, time.Minute*2)
	//conn, err := net.Dial("tcp", box.box.ModeInfo) //, time.Minute*10)
	if err != nil {
		return fmt.Errorf("ConnectError: %s", err.Error())
	}
	// encBuf := bufio.NewWriter(conn)
	// codec := &gobClientCodec{conn, gob.NewDecoder(conn), gob.NewEncoder(encBuf), encBuf}
	box.rpcClient = rpc.NewClient(conn)
	Log("box conn:", box.box.ModeInfo, " at:", conn.LocalAddr())
	return nil
}

// Call rpc 调研访问功能
func (box *BoxRpcClient) Call(in def.RequestIn, hres *def.BoxOutPut) (err error) {
	in.Input.IsSync = in.Input.IsSync || box.box.IsSync
	err = box.call("RpcWorker.Call", in, hres)
	return
}

// Status rpc 查询状态
func (box *BoxRpcClient) Status(in def.RequestIn, hres *def.BoxOutPut) (err error) {
	err = box.call("RpcWorker.Status", in, hres)
	return
}

// Ping 心跳监测
func (box *BoxRpcClient) Ping(in string, out *string) (ok bool) {
	err := box.call("RpcWorker.Ping", in, out)
	ok = err == nil && *out == "ok"
	return
}

// Close 关闭连接
func (box *BoxRpcClient) Close() error {
	if nil != box.rpcClient {
		return box.rpcClient.Close()
	}
	return nil
}

func (box *BoxRpcClient) call(callName string, in interface{}, hres interface{}) (err error) {
	tryOnce := true
TRY:
	err = box.rpcClient.Call(callName, in, hres)
	if err != nil && tryOnce {
		tryOnce = false
		err = box.init()
		goto TRY
	}
	return
}
