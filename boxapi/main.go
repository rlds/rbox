//
//  main.go
//  rbox
//
//  Created by 吴道睿 on 2018/9/7.
//  Copyright © 2018年 吴道睿. All rights reserved.
//
package main

import (
	"flag"
	"net/http"
	"time"

	"github.com/rlds/rlog"
)

/*
   webserver展示层
*/
func main() {
	httpAddr := ""
	logdir := ""
	maxLogFileByte := uint64(1024 * 1024 * 100)
	flag.StringVar(&httpAddr, "host", ":9888", "启动ShowServer的地址")
	flag.StringVar(&logdir, "log", "../log", "日志输出文件夹")
	// flag.StringVar(&jsCSSHtmlTmplFileDirPath, "static", "../static", "静态模版文件和js、css所在文件夹路径")
	flag.Parse()
	rlog.LogInit(3, logdir, maxLogFileByte, 1)
	// Log("jsCSSHtmlTmplFileDirPath :", jsCSSHtmlTmplFileDirPath)
	sysd.Init()
	webServerStart(httpAddr)
}

func webServerStart(httpAddr string) {
	http.HandleFunc("/", pages)
	svr := http.Server{
		Addr:           httpAddr,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		Handler:        http.DefaultServeMux,
		MaxHeaderBytes: 1 << 20,
	}
	svr.SetKeepAlivesEnabled(true)
	if err := svr.ListenAndServe(); nil != err {
		Log("ListenAndServe error:", err.Error())
	}
}

// Log 日志输出
func Log(arg ...interface{}) {
	rlog.V(1).Info(arg...)
}
