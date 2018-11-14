//
//  main.go
//  mbox
//
//  Created by 吴道睿 on 2018/9/7.
//  Copyright © 2018年 吴道睿. All rights reserved.
//
package main

import(
    "net/http"
	"time"
	"fmt"
	"flag"
)

func init(){
	//Log("init")
}

/*
   webserver展示层
*/
func main(){
	httpAddr := ""
	flag.StringVar(&httpAddr,"host",":9888","启动webserver的地址")
	flag.StringVar(&js_css_html_tmpl_file_dir_path,"staticdir","../static","静态模版文件和js、css所在文件夹路径")
	flag.Parse()
	println(js_css_html_tmpl_file_dir_path)
	sysd.Init()
	webServerStart(httpAddr)
}

func webServerStart(httpAddr string){
	http.HandleFunc("/",pages)
	svr := http.Server{
		Addr:           httpAddr,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		Handler:        http.DefaultServeMux,
		MaxHeaderBytes: 1 << 20,
	}
	svr.SetKeepAlivesEnabled(true)
	if err := svr.ListenAndServe(); nil != err {
		Log("ListenAndServe error:%s", err.Error())
	}
}

func Log(arg...interface{}){
	fmt.Println(arg...)
}
