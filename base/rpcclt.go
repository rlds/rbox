package base

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"io"
	"net"
	"net/rpc"
	"time"

	"github.com/rlds/rbox/base/def"
)

/*
   rpc模式的客户端
*/

type gobClientCodec struct {
	rwc io.ReadWriteCloser
	dec *gob.Decoder
	enc *gob.Encoder
	// dec    *json.Decoder
	// enc    *json.Encoder
	encBuf *bufio.Writer
}

func (c *gobClientCodec) WriteRequest(r *rpc.Request, body interface{}) (err error) {
	if err = timeoutCoder(c.enc.Encode, r, "client write request"); err != nil {
		return
	}
	if err = timeoutCoder(c.enc.Encode, body, "client write request body"); err != nil {
		return
	}
	return c.encBuf.Flush()
}

func (c *gobClientCodec) ReadResponseHeader(r *rpc.Response) error {
	return c.dec.Decode(r)
}

func (c *gobClientCodec) ReadResponseBody(body interface{}) error {
	return c.dec.Decode(body)
}

func (c *gobClientCodec) Close() error {
	return c.rwc.Close()
}

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
	conn, err := net.DialTimeout("tcp", box.box.ModeInfo, time.Minute*10)
	if err != nil {
		return fmt.Errorf("ConnectError: %s", err.Error())
	}
	encBuf := bufio.NewWriter(conn)
	codec := &gobClientCodec{conn, gob.NewDecoder(conn), gob.NewEncoder(encBuf), encBuf}
	// codec := &gobClientCodec{conn, json.NewDecoder(conn), json.NewEncoder(encBuf), encBuf}
	box.rpcClient = rpc.NewClientWithCodec(codec)
	Log("box conn:", box.box.ModeInfo, " at:", conn.LocalAddr())
	return nil
}

// Call rpc 调研访问功能
func (box *BoxRpcClient) Call(in def.RequestIn, hres *def.BoxOutPut) (err error) {
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
