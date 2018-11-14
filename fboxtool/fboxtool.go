//
//  boxtool.go
//  mbox
//
//  Created by 吴道睿 on 2018/4/26.
//  Copyright © 2018年 吴道睿. All rights reserved.
//
package main

import (
	"flag"
	"fmt"
	"os"
	"text/template"
	"time"

	"github.com/rlds/rbox/base"
	"github.com/rlds/rbox/base/util"

	"gopkg.in/yaml.v2"
)
/*
   功能：自动生成初始代码
   计划提供成为服务
*/
type toolConfig struct {
	TmplFileDirPath string
	BuildBoxConf    base.BoxConfig
}

func main() {
	conf, outdir := "", ""
	flag.StringVar(&conf, "conf", "", "基础公共配置参数")
	flag.StringVar(&outdir, "codedir", "", "生成的代码存放文件夹位置(需要保证文件夹是可以访问的)")
	flag.Parse()
	if conf == "" {
		fmt.Println("配置参数不能为空！")
		return
	}

	data := util.GetAllFileData(conf)
	var config toolConfig
	if len(data) > 0 {
		err := yaml.Unmarshal(data, &config)
		if err == nil {
			createDo(outdir, config)
		} else {
			fmt.Println("config error:", err)
		}
	} else {
		fmt.Println("config error.")
	}
}

const (
	boxMainTmplFileName        = "Main.go"
	boxMakefileTmplFileName    = "makefile"
	boxConfigTmplFileName      = "Config.go"
	boxKeepstartShTmplFileName = "KeepStart.sh"
)

// 生成文件
func createDo(outfiledirpath string, cfg toolConfig) {

	tlcfg := make(map[string]interface{})
	tlcfg["BoxConf"] = cfg.BuildBoxConf
	tlcfg["Time"] = time.Now().Format("2006-01-02 15:04:05")

	tmplDirPath := cfg.TmplFileDirPath

	// 主函数文件
	outfilename := outfiledirpath + "/" + cfg.BuildBoxConf.BoxInfo.Name + boxMainTmplFileName
	createBoxFile(tmplDirPath+"/"+boxMainTmplFileName, outfilename, tlcfg)

	// 工具的配置信息
	outfilename = outfiledirpath + "/" + cfg.BuildBoxConf.BoxInfo.Name + boxConfigTmplFileName
	createBoxFile(tmplDirPath+"/"+boxConfigTmplFileName, outfilename, tlcfg)

	// 执行脚本
	outfilename = outfiledirpath + "/" + cfg.BuildBoxConf.BoxInfo.Name + boxKeepstartShTmplFileName
	createBoxFile(tmplDirPath+"/"+boxKeepstartShTmplFileName, outfilename, tlcfg)

	// makefile
	outfilename = outfiledirpath + "/" + boxMakefileTmplFileName
	createBoxFile(tmplDirPath+"/"+boxMakefileTmplFileName, outfilename, tlcfg)
}

func createBoxFile(tmplfilename, outfile string, tlcfg map[string]interface{}) {
	s1, err := template.ParseFiles(tmplfilename)
	if err != nil {
		fmt.Println(tmplfilename, " parse err:", err)
		return
	}

	// 打开创建文件
	// 覆盖
	ofile, err2 := os.OpenFile(outfile, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, os.ModePerm)
	if err2 != nil {
		fmt.Println(outfile, " file create err:", err2)
		return
	}
	s1.Execute(ofile, tlcfg)
	ofile.Close()
}
