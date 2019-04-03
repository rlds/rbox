//
//  httpclt.go
//  base
//
//  Created by 吴道睿 on 2018/4/26.
//  Copyright © 2018年 吴道睿. All rights reserved.
//
package base

import (
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/rlds/rbox/base/def"
	"github.com/rlds/rbox/base/util"
)

/*
   用于http模式注册模块
*/
var (
	requestimeout     = time.Second * 10
	requestRetryTimes = 3
)

// RegBoxPost 注册
func RegBoxPost(regpath, data string) (result []byte, err error) {
	reTimes := 0
	var (
		r    *http.Request
		resp *http.Response
	)
	r, err = http.NewRequest("POST", regpath, strings.NewReader(data))
	if err != nil {
		return
	}
	r.Header.Add("Content-Type", "application/json; charset=UTF-8")

RECONNECT:
	htc := &http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				conn, err := net.DialTimeout(netw, addr, requestimeout)
				if err != nil {
					return nil, err
				}
				conn.SetDeadline(time.Now().Add(requestimeout))
				return conn, nil
			},
			ResponseHeaderTimeout: requestimeout,
		},
	}
	if resp, err = htc.Do(r); err == nil {
		result, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			Log("err1[", err, "]")
			return
		}
		resp.Body.Close()
	} else {
		Log("err2[", err, "]")
		if reTimes < requestRetryTimes {
			reTimes++
			Log("尝试第[", reTimes, "] 次重新连接")
			goto RECONNECT
		}
	}
	return
}

// BoxHTTPClient 客户端信息
type BoxHTTPClient struct {
	box        *def.BoxInfo
	callPath   string
	statusPath string
	pingPath   string
	httpClient *http.Client
}

func NewHTTPClient(box *def.BoxInfo) (clt *BoxHTTPClient, err error) {
	clt = new(BoxHTTPClient)
	clt.box = box
	err = clt.init()
	return
}

func (clt *BoxHTTPClient) init() error {
	clt.httpClient = &http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				conn, err := net.DialTimeout(netw, addr, requestimeout)
				if err != nil {
					return nil, err
				}
				conn.SetDeadline(time.Now().Add(requestimeout))
				return conn, nil
			},
			ResponseHeaderTimeout: requestimeout,
		},
	}
	clt.callPath = clt.box.ModeInfo + "/call/" + clt.box.Group + "/" + clt.box.Name
	clt.statusPath = clt.box.ModeInfo + "/taskRes/" + clt.box.Group + "/" + clt.box.Name
	clt.pingPath = clt.box.ModeInfo + "/ping"
	return nil
}

// Call 功能调用
func (clt *BoxHTTPClient) Call(in def.RequestIn, hres *def.BoxOutPut) error {
	return clt.post(clt.callPath, util.ObjToStr(in), hres)
}

// Status 状态查询
func (clt *BoxHTTPClient) Status(in def.RequestIn, hres *def.BoxOutPut) error {
	return clt.post(clt.statusPath, util.ObjToStr(in), hres)
}

// Ping 心跳监测
func (clt *BoxHTTPClient) Ping(in string, out *string) bool {
	return clt.postPing()
}
func (clt *BoxHTTPClient) post(urlpath, indat string, out *def.BoxOutPut) (err error) {
	reTimes := 0
	var (
		r    *http.Request
		resp *http.Response
		body []byte
	)
	r, err = http.NewRequest("POST", urlpath, strings.NewReader(indat))
	if err != nil {
		return
	}
	r.Header.Add("Content-Type", "application/json; charset=UTF-8")
RECONNECT:
	if resp, err = clt.httpClient.Do(r); err == nil {
		body, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return
		}
		resp.Body.Close()
		err = json.Unmarshal(body, out)
	} else {
		if reTimes < requestRetryTimes {
			reTimes++
			Log("尝试第[", reTimes, "] 次重新连接", urlpath)
			goto RECONNECT
		}
	}
	return
}

func (clt *BoxHTTPClient) postPing() (ok bool) {
	reTimes := 0
	var (
		r    *http.Request
		resp *http.Response
		body []byte
		err  error
	)
	r, err = http.NewRequest("POST", clt.pingPath, nil)
	if err != nil {
		return
	}
	r.Header.Add("Content-Type", "application/json; charset=UTF-8")
RECONNECT:
	if resp, err = clt.httpClient.Do(r); err == nil {
		body, err = ioutil.ReadAll(resp.Body)
		if err == nil {
			resp.Body.Close()
			ok = string(body) == "ok"
		}
	} else {
		if reTimes < requestRetryTimes {
			reTimes++
			Log("尝试第[", reTimes, "] 次重新连接", clt.pingPath)
			goto RECONNECT
		}
	}
	return
}
