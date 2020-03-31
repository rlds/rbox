//
//  runModeHttp.go
//  base
//
//  Created by 吴道睿 on 2018/4/6.
//  Copyright © 2018年 吴道睿. All rights reserved.
//
package base

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	. "github.com/rlds/rbox/base/def"
	. "github.com/rlds/rbox/base/util"
)

/*
   测试可用，未优化
*/

/*
   http模式的执行
*/
type (
	httpModeWorker struct {
	}
)

//首次运行注册至webserver 通知开始服务了
func (h *httpModeWorker) Register() {
	gbox.cfg.Mode = "http"
	Log(gbox.cfg)
	res, err := RegBoxPost(gbox.cfg.ShowServerPath, ObjToStr(gbox.cfg.BoxInfo))
	if err == nil {
		Log("regok:", string(res))
	} else {
		Log("regerr:", err)
	}
}

//http server启动
func (h *httpModeWorker) Run() {
	boxUrlPath := "/" + gbox.cfg.Group + "/" + gbox.cfg.Name
	Log("boxUrlPath:", boxUrlPath)
	http.HandleFunc(boxUrlPath, h.boxServ)
	http.HandleFunc("/call"+boxUrlPath, h.boxServ)
	http.HandleFunc("/taskRes"+boxUrlPath, h.boxTaskRes)
	http.HandleFunc("/ok.htm", h.ping)
	http.HandleFunc("/ping", h.ping)
	http.HandleFunc("/about", h.about)
	svr := http.Server{
		Addr:           gbox.cfg.SelfHttpServerHost,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		Handler:        http.DefaultServeMux,
		MaxHeaderBytes: 1 << 20,
	}

	svr.SetKeepAlivesEnabled(true)
	Log("httpModeStart:", gbox.cfg.SelfHttpServerHost)
	//注册执行
	h.Register()

	//开启监听
	if err := svr.ListenAndServe(); nil != err {
		Log("ListenAndServe error:", err.Error())
	}
}

//输入 content-type: application/json
//http模式接口的入口处理
func (h *httpModeWorker) boxServ(w http.ResponseWriter, r *http.Request) {
	var hres BoxOutPut
	var hi RequestIn
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		hres.Code = OutputRetuen_Error
		hres.ReturnMsg = err.Error()
	} else {
		err = json.Unmarshal(body, &hi)
		if err != nil {
			hres.Code = OutputRetuen_Error
			hres.ReturnMsg = err.Error()
		} else {
			Log("call T:", hi.TaskId, " F:", hi.From, " C:", hi.Call) //, " I:", hi.Input)
			// Log("call T:", hi.TaskId, " F:", hi.From, " C:", hi.Call, " I:", hi.Input)
			// hi.Input.TaskId = hi.TaskId
			box.DoWork(hi.TaskId, hi.Input)
			hres = box.Output(hi.TaskId)
			hres.TaskId = hi.TaskId
		}
	}
	b, _ := json.Marshal(hres)
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
	Log(hi.TaskId, " :", string(b))
	return
}

/*
   异步执行时获取结果
*/
func (h *httpModeWorker) boxTaskRes(w http.ResponseWriter, r *http.Request) {
	var hres BoxOutPut
	var hi RequestIn
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		hres.Code = OutputRetuen_Error
		hres.ReturnMsg = err.Error()
	} else {
		err = json.Unmarshal(body, &hi)
		if err != nil {
			hres.Code = OutputRetuen_Error
			hres.ReturnMsg = err.Error()
		} else {
			Log("res T:", hi.TaskId, " F:", hi.From, " C:", hi.Call)
			hres = box.Output(hi.TaskId)
			hres.TaskId = hi.TaskId
		}
	}
	b, _ := json.Marshal(hres)
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
	Log(hi.TaskId, " :", string(b))
}

func (h *httpModeWorker) about(w http.ResponseWriter, r *http.Request) {
	b, _ := json.Marshal(gbox.cfg.BoxInfo)
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

//心跳检测的响应
func (h *httpModeWorker) ping(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

/*
var cltMap map[string]*http.Client

func HttpModeCallBox(group, boxname string, indata InputData) (rtype string, data interface{}) {
	key := group + boxname
	if cltMap == nil {
		cltMap = make(map[string]*http.Client)
	}

}
*/
