//
//  httpclt.go
//  base
//
//  Created by 吴道睿 on 2018/4/26.
//  Copyright © 2018年 吴道睿. All rights reserved.
//
package mbox

import(
	   "time"
	   "net/http"
	   "io/ioutil"
	   "net"
	   "strings"
)
/*
    用于http模式注册模块
*/
var (
	 requestimeout      = time.Second * 10
	 requestRetryTimes  = 3
)

func RegBoxPost(regpath,data string)(result []byte,err error){
	re_times := 0
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
	htc:= &http.Client{
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
	if resp,err = htc.Do(r) ;err == nil{
		result, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			Log("err1[", err,"]")
			return
		}
		resp.Body.Close()
	}else{
		Log("err2[",err,"]")
		if re_times < requestRetryTimes {
			re_times ++
			Log("尝试第[",re_times,"] 次重新连接")
			goto RECONNECT
		}
	}
	return
}
