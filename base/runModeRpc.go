package base

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"io"
	"net"
	"net/rpc"
	"strings"
	"time"

	"github.com/rlds/rbox/base/def"
	. "github.com/rlds/rbox/base/util"
)

/*
   rpc模式接口
*/

type RpcWorker struct {
	Name string
}

func newWorker() *RpcWorker {
	return &RpcWorker{"test"}
}

func (w *RpcWorker) Call(in def.RequestIn, hres *def.BoxOutPut) error {
	Log("call T:", in.TaskId, " F:", in.From, " C:", in.Call)
	box.DoWork(in.TaskId, in.Input)
	*hres = box.Output(in.TaskId)
	hres.TaskId = in.TaskId
	return nil
}

func (w *RpcWorker) Status(in def.RequestIn, hres *def.BoxOutPut) error {
	Log("status T:", in.TaskId, " F:", in.From, " C:", in.Call)
	*hres = box.Output(in.TaskId)
	hres.TaskId = in.TaskId
	return nil
}

func (w *RpcWorker) Ping(in string, out *string) error {
	*out = "ok"
	return nil
}

func timeoutCoder(f func(interface{}) error, e interface{}, msg string) error {
	echan := make(chan error, 1)
	go func() { echan <- f(e) }()
	select {
	case e := <-echan:
		return e
	case <-time.After(time.Minute):
		return fmt.Errorf("Timeout %s", msg)
	}
}

type gobServerCodec struct {
	rwc io.ReadWriteCloser
	dec *gob.Decoder
	enc *gob.Encoder
	// dec    *json.Decoder
	// enc    *json.Encoder
	encBuf *bufio.Writer
	closed bool
}

func (c *gobServerCodec) ReadRequestHeader(r *rpc.Request) error {
	return timeoutCoder(c.dec.Decode, r, "server read request header")
}

func (c *gobServerCodec) ReadRequestBody(body interface{}) error {
	return timeoutCoder(c.dec.Decode, body, "server read request body")
}

func (c *gobServerCodec) WriteResponse(r *rpc.Response, body interface{}) (err error) {
	if err = timeoutCoder(c.enc.Encode, r, "server write response"); err != nil {
		if c.encBuf.Flush() == nil {
			Log("serv rpc: gob error encoding response:", err)
			c.Close()
		}
		return
	}
	if err = timeoutCoder(c.enc.Encode, body, "server write response body"); err != nil {
		if c.encBuf.Flush() == nil {
			Log("serv rpc: gob error encoding body:", err)
			c.Close()
		}
		return
	}
	return c.encBuf.Flush()
}

func (c *gobServerCodec) Close() error {
	if c.closed {
		// Only call c.rwc.Close once; otherwise the semantics are undefined.
		return nil
	}
	c.closed = true
	return c.rwc.Close()
}

type rpcModeWorker struct {
}

func checkRpcModeHost(host string) string {
	arr := strings.SplitN(host, "//", 2)
	if len(arr) == 2 {
		return arr[1]
	}
	return host
}

// 注册
func (r *rpcModeWorker) Register() {
	gbox.cfg.Mode = "rpc"
	//Log(gbox.cfg)
	gbox.cfg.BoxInfo.ModeInfo = checkRpcModeHost(gbox.cfg.BoxInfo.ModeInfo)
	//Log(gbox.cfg)
	res, err := RegBoxPost(gbox.cfg.ShowServerPath, ObjToStr(gbox.cfg.BoxInfo))
	if err == nil {
		Log("regok:", string(res))
	} else {
		Log("regerr:", err)
	}
}

// 启动服务
func (r *rpcModeWorker) Run() {
	Log("rpc mode start")

	gobRegister()

	rpc.Register(newWorker())
	host := gbox.cfg.SelfHttpServerHost
	l, e := net.Listen("tcp", host)
	if e != nil {
		Log("Error: listen ", host, " error:", e)
		return
	}

	// 开始注册
	r.Register()

	// 阻塞运行
	for {
		conn, err := l.Accept()
		if err != nil {
			Log("Error: accept rpc connection", err.Error())
			continue
		}
		go func(conn net.Conn) {
			buf := bufio.NewWriter(conn)
			srv := &gobServerCodec{
				rwc: conn,
				dec: gob.NewDecoder(conn),
				enc: gob.NewEncoder(buf),
				// dec:    json.NewDecoder(conn),
				// enc:    json.NewEncoder(buf),
				encBuf: buf,
			}
			err = rpc.ServeRequest(srv)
			if err != nil {
				Log("Error: server rpc request", err.Error())
			}
			srv.Close()
		}(conn)
	}
}
