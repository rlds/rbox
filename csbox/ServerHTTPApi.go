//
//  ServerHTTPApi.go
//  rbox
//
//  Created by 吴道睿 on 2018/4/18.
//  Copyright © 2018年 吴道睿. All rights reserved.
//
package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"text/template"

	"github.com/rlds/rbox/base/def"
	"github.com/rlds/rbox/base/util"
)

const (
	BoxOutPutCodeOk    = "0"
	BoxOutPutCallBoxOk = "CallBox_Ok"
	BoxOutPutRegBoxOk  = "RegBox_Ok"

	CallBoxCodeInputErr100  = "100"
	CallBoxCodeInputErr101  = "101"
	CallBoxCodeInputErr102  = "102"
	CallBoxCodeInputErr103  = "103"
	CallBox_ResTypeInputErr = "InputErr"

	CallBoxCodeInputErr104 = "104"
	CallBoxCodeInputErr105 = "105"
	CallBoxResTypeCallErr  = "CallErr"

	RegBoxCodeInputErr106  = "106"
	RegBoxCodeInputErr107  = "107"
	RegBoxCodeInputErr108  = "108"
	RegBoxCodeInputErr109  = "109"
	RegBox_ResTypeInputErr = "RegboxErr"

	CallBoxCodeBoxRetErr110 = "110"
	CallBoxResTypeBoxRetErr = "BoxRetErr"
)

var (
	jsCSSHtmlTmplFileDirPath = "../static"
	indexTmlPth              = jsCSSHtmlTmplFileDirPath + "/mindex.tmpl"
	grpltTmlPth              = jsCSSHtmlTmplFileDirPath + "/grouplist.tmpl"
	boxifTmlPth              = jsCSSHtmlTmplFileDirPath + "/boxinfo.tmpl"
	boxslTmlPth              = jsCSSHtmlTmplFileDirPath + "/boxsel.tmpl"
	boxsmTmlPth              = jsCSSHtmlTmplFileDirPath + "/boxsel_m.tmpl"
	boxpmTmlPth              = jsCSSHtmlTmplFileDirPath + "/boxparam.tmpl"
)

//  当外部输入路径时保证更改路径正确
func setTmplFilePath() {
	indexTmlPth = jsCSSHtmlTmplFileDirPath + "/mindex.tmpl"
	grpltTmlPth = jsCSSHtmlTmplFileDirPath + "/grouplist.tmpl"
	boxifTmlPth = jsCSSHtmlTmplFileDirPath + "/boxinfo.tmpl"
	boxslTmlPth = jsCSSHtmlTmplFileDirPath + "/boxsel.tmpl"
	boxsmTmlPth = jsCSSHtmlTmplFileDirPath + "/boxsel_m.tmpl"
	boxpmTmlPth = jsCSSHtmlTmplFileDirPath + "/boxparam.tmpl"
}

func cssjs(w http.ResponseWriter, r *http.Request) {
	//
	filedat := util.GetAllFileData(jsCSSHtmlTmplFileDirPath + r.URL.Path)
	arr := strings.Split(r.URL.Path, ".")
	contype := ""
	if len(arr) > 0 {
		switch arr[1] {
		case "css":
			{
				contype = "text/css"
			}
		case "js":
			{
				contype = "application/javascript"
			}
		}
	}
	w.Header().Set("Content-Type", contype)
	w.Write(filedat)
}

/*
 页面处理
*/
func pages(w http.ResponseWriter, r *http.Request) {
	doarr := strings.Split(r.URL.Path, "/")
	do := ""
	if len(doarr) > 1 {
		do = doarr[1]
		switch do {
		case "call":
			{ //调用box
				callbox(w, r)
			}
		case "status":
			{
				taskRes(w, r)
			}
		case "login":
			{ //登录
				login(w, r)
			}
		case "regbox":
			{ //注册box
				regbox(w, r)
			}
		case "boxinfo":
			{ //获得box信息
				boxinfo(w, r)
			}
		case "updategroup":
			{ //设置group的信息
				updategroup(w, r)
			}
		case "ok.htm", "ping.htm", "ok.html", "ping.html":
			{
				w.Write([]byte("ok"))
			}
		default:
			{
				arr := strings.Split(r.URL.Path, ".")
				contype := ""
				var filedat []byte
				if len(arr) > 1 {
					switch arr[1] {
					case "css":
						{
							contype = "text/css"
						}
					case "js":
						{
							contype = "application/javascript"
						}
					case "html":
						{
							pagetmldo(arr[0], w, r)
							return
						}
					case "woff2":
						{
							contype = "font/woff2"
						}
					case "wasm":
						{
							contype = "application/wasm"
						}
					case "ico":
						{

						}
					default:
						{
							Log("default", r.URL.Path)
						}
					}
					filedat = util.GetAllFileData(jsCSSHtmlTmplFileDirPath + r.URL.Path)
					/*
					 页面及模版的处理
					*/
					w.Header().Set("Content-Type", contype)
					w.Write(filedat)
				}
			}
		}
	}
}

const (
	serverInfo = "webserver/v1.0.1/wdr"
)

func pagetmldo(gbox string, w http.ResponseWriter, r *http.Request) {
	arr := strings.Split(gbox, "/")
	group, box := "", ""
	il := len(arr)
	if il == 2 {
		group = arr[1]
	} else if il > 2 {
		group, box = arr[1], arr[2]
	}
	w.Header().Set("Content-Type", "text/html")
	s1, err := template.ParseFiles(indexTmlPth, grpltTmlPth, boxifTmlPth, boxslTmlPth, boxsmTmlPth, boxpmTmlPth)
	if err != nil {
		//当模版文件不存在时会报错
		w.Write([]byte(err.Error()))
		return
	}
	/*
	  s1.ExecuteTemplate(w, "grouplist", nil)
	  s1.ExecuteTemplate(w, "boxsel", nil)
	  s1.ExecuteTemplate(w, "boxsel_m", nil)
	  s1.ExecuteTemplate(w, "boxparam", nil)
	  s1.ExecuteTemplate(w, "boxinfo", nil)
	*/
	tmplmap := make(map[string]interface{})
	//Log(indexTmlPth)
	grouplist := sysd.GetUserdGroup()
	tmplmap["Grouplist"] = grouplist
	tmplmap["Server"] = serverInfo
	if grouplist != nil && len(grouplist) > 0 && len(group) > 0 && len(group) < 30 {
		boxl := sysd.GetGroupUsedBox(group)
		boxnum := len(boxl)
		if boxnum < 1 {
			goto DOEND
		} else if boxnum > 5 {
			tmplmap["Boxsel_m"] = boxl
			tmplmap["Bm"] = true
		} else {
			tmplmap["Boxsel"] = boxl
			tmplmap["Bm"] = false
		}
		if len(box) > 0 && len(box) < 30 {
			box, err := sysd.GetCallBox(group, box)
			if err == nil {
				tmplmap["Hsb"] = len(box.SubBox) > 1
				tmplmap["SubBoxs"] = box.SubBox
				//tmplmap["Boxparam"] = box.Params
				tmplmap["Boxinfo"] = box
			}
		}
	}
DOEND:
	err = s1.Execute(w, tmplmap)
	if err != nil {
		w.Write([]byte(err.Error()))
		Log("Execute err:", err)
	}
}

func regbox(w http.ResponseWriter, r *http.Request) {
	/* 读取消息体
	   结构必须为box信息结构 否则报错
	   返回结构信息结构统一
	*/
	var bp def.BoxOutPut
	var box def.BoxInfo
	body, err := ioutil.ReadAll(r.Body)
	if err != nil || len(body) < 10 {
		bp.Type = RegBox_ResTypeInputErr
		bp.Code = RegBoxCodeInputErr106
		bp.ReturnMsg = "RegInPutErr"
		goto RegEnd
	}
	err = json.Unmarshal(body, &box)
	if err != nil {
		bp.Type = RegBox_ResTypeInputErr
		bp.Code = RegBoxCodeInputErr107
		bp.ReturnMsg = err.Error()
		goto RegEnd
	}
	box.SetModeInfo(r.RemoteAddr)
	//box 初步检测信息正确性
	if box.InfoOk() {
		//添加box
		err = sysd.AddBox(box)
		if err == nil {
			//添加完成正常返回
			bp.Type = BoxOutPutRegBoxOk
			bp.Code = BoxOutPutCodeOk
			bp.ReturnMsg = ""
		} else {
			bp.Type = RegBox_ResTypeInputErr
			bp.Code = RegBoxCodeInputErr109
			bp.ReturnMsg = err.Error()
		}
	} else {
		bp.Type = RegBox_ResTypeInputErr
		bp.Code = RegBoxCodeInputErr108
		bp.ReturnMsg = "BoxChkErr"
		goto RegEnd
	}

RegEnd:
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.Write(util.ObjToBytes(bp))
}

//获取box信息
func boxinfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	group := r.FormValue("group")
	var b []byte
	if len(group) > 0 {
		b = util.ObjToBytes(sysd.GetGroupUsedBox(group))
	} else {
		b = util.ObjToBytes(sysd.GetUserdGroup())
	}
	w.Write(b)
}

//更新分组信息
func updategroup(w http.ResponseWriter, r *http.Request) {
	var bp def.BoxOutPut
	var bgi boxGroupInfo

	body, err := ioutil.ReadAll(r.Body)
	if err != nil || len(body) < 10 {
		bp.Type = RegBox_ResTypeInputErr
		bp.Code = RegBoxCodeInputErr106
		bp.ReturnMsg = "RegInPutErr"
		goto UpdateEnd
	}

	err = json.Unmarshal(body, &bgi)
	if err != nil {
		bp.Type = RegBox_ResTypeInputErr
		bp.Code = RegBoxCodeInputErr107
		bp.ReturnMsg = err.Error()
		goto UpdateEnd
	}

	//信息更新成功
	err = sysd.AddBoxGroup(bgi)
	if err == nil {
		bp.Type = BoxOutPutRegBoxOk
		bp.Code = BoxOutPutCodeOk
		bp.ReturnMsg = ""
	} else {
		bp.Type = RegBox_ResTypeInputErr
		bp.Code = RegBoxCodeInputErr106
		bp.ReturnMsg = err.Error()
	}

UpdateEnd:
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.Write(util.ObjToBytes(bp))
}

/*
 登录处理
*/
func login(w http.ResponseWriter, r *http.Request) {

}

/*
 调用工具的功能进行传递参数
 处理结果返回
*/

/*
 调用具体工具 页面使用ajax进行访问
 输入为json格式
 输出为json格式
*/
func callbox(w http.ResponseWriter, r *http.Request) {
	var bp def.BoxOutPut
	/*
	 权限校验
	*/
	group, boxnm := "", ""
	var box *boxInfo
	var arr []string
	var indat def.InputData
	indat.Data = make(map[string]interface{})
	//params := make(map[string]string)
	bp.TaskId = sysd.tp.GetTid()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		bp.Type = CallBox_ResTypeInputErr
		bp.Code = CallBoxCodeInputErr100
		bp.ReturnMsg = err.Error()
		goto CALLEND
	}

	if len(body) > 5 {
		err = json.Unmarshal(body, &indat)
		if err != nil {
			bp.Type = CallBox_ResTypeInputErr
			bp.Code = CallBoxCodeInputErr101
			bp.ReturnMsg = err.Error()
			goto CALLEND
		}
	}

	r.ParseForm()
	if indat.Data == nil {
		indat.Data = make(map[string]interface{})
	}
	for k, v := range r.Form {
		if len(v) > 1 {
			indat.Data[k] = util.ObjToStr(v)
		} else {
			indat.Data[k] = v
		}
	}
	Log("header:", util.ObjToStr(r.Header))
	//   params["xRealIP"]       = r.Header.Get("X-Real-Ip")
	//   params["xForwardedFor"] = r.Header.Get("X-Forwarded-For")
	indat.Data["remoteAddr"] = r.RemoteAddr
	//------------------------------
	arr = strings.Split(r.URL.Path, "/")
	if len(arr) > 3 {
		//Log(arr)
		group, boxnm = arr[2], arr[3]
		box, err = sysd.GetCallBox(group, boxnm)
		if err != nil {
			bp.Type = CallBox_ResTypeInputErr
			bp.Code = CallBoxCodeInputErr102
			bp.ReturnMsg = err.Error()
			goto CALLEND
		}

		var callin def.RequestIn
		callin.TaskId = bp.TaskId
		callin.From = serverInfo
		callin.Input = indat
		callin.Call = boxnm
		//执行
		bp = box.DoWork(callin)
	} else {
		bp.Type = CallBox_ResTypeInputErr
		bp.Code = CallBoxCodeInputErr105
		bp.ReturnMsg = "input err:box name null"
	}
	//------------------------------
CALLEND:
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.Write(util.ObjToBytes(bp))
}

/*
 调用具体工具 页面使用ajax进行访问
 输入为json格式
 输出为json格式
*/
func taskRes(w http.ResponseWriter, r *http.Request) {
	var bp def.BoxOutPut
	/*
	 权限校验
	*/
	group, boxnm := "", ""
	var box *boxInfo
	var arr []string
	var indat def.InputData
	//params := make(map[string]string)
	//bp.TaskId = sysd.tp.GetTid()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		bp.Type = CallBox_ResTypeInputErr
		bp.Code = CallBoxCodeInputErr100
		bp.ReturnMsg = err.Error()
		goto CALLEND
	}

	if len(body) > 5 {
		err = json.Unmarshal(body, &indat)
		if err != nil {
			bp.Type = CallBox_ResTypeInputErr
			bp.Code = CallBoxCodeInputErr101
			bp.ReturnMsg = err.Error()
			goto CALLEND
		}
	}

	r.ParseForm()
	if indat.Data == nil {
		indat.Data = make(map[string]interface{})
	}
	for k, v := range r.Form {
		if len(v) > 1 {
			indat.Data[k] = util.ObjToStr(v)
		} else {
			indat.Data[k] = v
		}
	}

	//------------------------------
	arr = strings.Split(r.URL.Path, "/")
	if len(arr) > 3 {
		//Log(arr)
		group, boxnm = arr[2], arr[3]
		box, err = sysd.GetCallBox(group, boxnm)
		if err != nil {
			bp.Type = CallBox_ResTypeInputErr
			bp.Code = CallBoxCodeInputErr102
			bp.ReturnMsg = err.Error()
			goto CALLEND
		}

		var callin def.RequestIn
		callin.TaskId = r.FormValue("TaskId") //bp.TaskId
		callin.From = serverInfo
		callin.Input = indat
		callin.Call = boxnm
		//执行
		bp = box.TaskRes(callin)
	} else {
		bp.Type = CallBox_ResTypeInputErr
		bp.Code = CallBoxCodeInputErr105
		bp.ReturnMsg = "input err:box name null"
	}
	//------------------------------
CALLEND:
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.Write(util.ObjToBytes(bp))
}
