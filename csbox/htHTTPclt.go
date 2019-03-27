//
//  htHTTPclt.go
//  rbox
//
//  Created by 吴道睿 on 2018/4/20.
//  Copyright © 2018年 吴道睿. All rights reserved.
//
package main

/*
   http客户端
*/
import (
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

var (
	requestimeout     = time.Second * 10
	requestRetryTimes = 3
	lc                sync.RWMutex
	hcleintmap        map[string]*http.Client
)

func httpClinetmapInit() {
	hcleintmap = make(map[string]*http.Client)
}

func setHostClt(host string, clt *http.Client) {
	lc.Lock()
	hcleintmap[host] = clt
	lc.Unlock()
}

func httpClientInit() (ts *http.Client) {
	ts = &http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				conn, err := net.DialTimeout(netw, addr, requestimeout)
				if err != nil {
					return nil, err
				}
				conn.SetDeadline(time.Now().Add(requestimeout)) //time.Second * d_time))
				//tcp_conn := conn.(*net.TCPConn)
				//tcp_conn.SetKeepAlive(false)
				return conn, nil
			},
			//DisableKeepAlives: true,
			ResponseHeaderTimeout: requestimeout,
		},
	}
	return
}

/*
   http请求客户端  GET
*/
// HTTPGet Get
func HTTPGet(urlstr string) (body []byte, err error) {
	reTimes := 0
	var (
		r      *http.Request
		resp   *http.Response
		result []byte
	)
	r, err = http.NewRequest("GET", urlstr, nil)
	if err != nil {
		return
	}

	r.Header.Add("Connection", "close")

RECONNECT:

	dhost := r.URL.Host
	htc, ok := hcleintmap[dhost]
	if !ok || htc == nil {
		htc = httpClientInit()
		setHostClt(dhost, htc)
	}

	ok = false
	if resp, err = htc.Do(r); err == nil {
		result, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return
		}
		resp.Body.Close()
		body = result
		ok = true
	} else {
		Log("err2[", err, "]")
		time.Sleep(time.Second)
		if reTimes < requestRetryTimes {
			reTimes++
			Log("尝试第[", reTimes, "] 次重新连接 url:", urlstr)
			goto RECONNECT
		}
	}
	return
}

// HTTPPost   POST urlencode
func HTTPPost(urlstr, data string) (body []byte, err error) {
	reTimes := 0
	var (
		r      *http.Request
		resp   *http.Response
		result []byte
	)
	r, err = http.NewRequest("POST", urlstr, strings.NewReader(data))
	if err != nil {
		return
	}
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")

RECONNECT:
	dhost := r.URL.Host
	htc, ok := hcleintmap[dhost]
	if !ok || htc == nil {
		htc = httpClientInit()
		setHostClt(dhost, htc)
	}
	ok = false
	if resp, err = htc.Do(r); err == nil {
		result, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			Log("err1[", err, "]")
			return
		}
		resp.Body.Close()
		body = result
		ok = true
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

// HTTPPostJSON json 格式数据发布
func HTTPPostJSON(urlstr, data string) (body []byte, err error) {
	reTimes := 0
	var (
		r      *http.Request
		resp   *http.Response
		result []byte
	)

	r, err = http.NewRequest("POST", urlstr, strings.NewReader(data))
	if err != nil {
		return
	}
	r.Header.Add("Content-Type", "application/json; charset=UTF-8")

RECONNECT:
	dhost := r.URL.Host
	htc, ok := hcleintmap[dhost]
	if !ok || htc == nil {
		htc = httpClientInit()
		setHostClt(dhost, htc)
	}

	ok = false
	if resp, err = htc.Do(r); err == nil {
		result, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return
		}
		resp.Body.Close()
		body = result
		ok = true
	} else {
		if reTimes < requestRetryTimes {
			reTimes++
			Log("尝试第[", reTimes, "] 次重新连接")
			goto RECONNECT
		}
	}
	return
}
