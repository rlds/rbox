// 
//  sys.go
//  rbox
//
//  Created by 吴道睿 on 2018/4/7.
//  Copyright © 2018年 吴道睿. All rights reserved.
//
package main

import (
	"sort"
	"time"

	"./taskid"
	"github.com/rlds/rbox/base"
	"github.com/rlds/rbox/base/def"
)

/*
   系统信息
*/
type (
	//工具的信息
	boxInfo struct {
		def.BoxInfo
		isAlive  bool   //是否活着可用
		connType string //接入方式 nats 或 http
	}

	//工具大类别信息
	boxGroupInfo struct {
		Name     string //与boxInfo.Group 相等
		ShowName string //展示给使用者看的文字标题
		Des      string //大类别的描述
		Owner    string //拥有者 此类别最高权限管理者列表
		Type     string //大类别的类型  //公共类型（无需权限校验），私有类型（需要权限校验），管理类型（不对普通用户开放）
		Status   string //大类别的状态
		//		  Idx        int      //排序编号
		mbif map[string]*boxInfo //此类别下工具列表
	}

	//系统
	webServerInfo struct { //webserver服务信息
		SelfHost   string //本机对外服务的ip及端口
		DomainHost string //外部可访问域名地址
	}

	//用户登录的信息
	loginInfo struct {
		Logined bool   //登录成功
		UserID  string //系统认可的用户唯一性id
		Name    string //用户名称
		//ShowName  string      //展示名
		Des      string //用户的描述
		Status   string //用户状态
		UserType string //用户类型
		Sex      string //性别  m(男)/f(女)
	}

	//
	athResInfo struct {
		Allow  bool   //是否允许访问
		ResStr string //不允许的原因
	}

	//登录鉴权接口
	regPro interface {
		LogIn(map[string]string) (uf loginInfo)
	}

	//权限管理接口
	athPro interface {
		CheckGroup(gp string) (ari athResInfo)
		CHeckBox(gp, bx string) (ari athResInfo)
	}

	//数据处理器接口
	dataPro interface {
		//数据处理器的初始化配置信息
		DataOption(cm map[string]string) error

		//数据写接口
		Set(key, value []byte)

		//数据获取接口
		Get(key []byte) (value []byte, err error)

		//数据删除接口
		Delete(key []byte) (err error)

		//数据处理器关闭
		Close() (err error)
	}

	//心跳检查接口
	hbeatPro interface {
		//放入需要进行心跳检查的box指针
		AddGroup(bg *boxGroupInfo) (err error)

		//心跳检查开始
		CheckStart()
	}

	//配置处理接口
	confPro interface {
		GetAllConf() (cf map[string]string, err error)
		Get(key string) (val string, err error)
		Set(key, value string) (err error)
		Delete(key string) (err error)
	}

	//用户session信息
	sessPro interface {
		//为用户分配sessionid
		Create(user loginInfo) (ssid string, exp time.Time)

		//检测session信息
		CheckSession(ssid string) (user *loginInfo, err error)
	}

	//生成任务id接口
	taskidPro interface {
		//获得一个任务id
		GetTid() string
	}

	//系统结构信息
	sysdata struct {
		/*
		  第一期先这么实现 这里需要重新定义
		*/
		mbgf     map[string]*boxGroupInfo //已接入攻击信息存储
		rscPath  string                   //资源模版文件夹路径位置
		confPath string                   //配置信息存储路径
		winfo    webServerInfo            //webserver的信息

		ss  sessPro   //session信息接口
		cf  confPro   //配置处理接口
		dp  dataPro   //数据的操作
		rg  regPro    //登录授权接口
		ath athPro    //权限管理接口
		hb  hbeatPro  //心跳检查接口
		tp  taskidPro //任务id生成器
	}
)

//
var (
	sysd sysdata
)

//系统初始化准备
func (s *sysdata) Init() {
	setTmplFilePath()

	s.mbgf = make(map[string]*boxGroupInfo)

	//相关接口的初始化注册初始化
	s.tp = taskid.NewTaskID()

	s.hb = NewWhbeat()

	go s.checkeAlive()
	//http客户端
	httpClinetmapInit()
}

//
func (s *sysdata) AddBoxGroup(bg boxGroupInfo) (err error) {
	bgf, ok := s.mbgf[bg.Name]
	if !ok {
		bgf = new(boxGroupInfo)
		*bgf = bg
		bgf.mbif = make(map[string]*boxInfo)
		err = s.hb.AddGroup(bgf)
	}
	//存在则进行数据更新
	bgf.ShowName = bg.ShowName
	bgf.Des = bg.Des
	bgf.Owner = bg.Owner
	bgf.Type = bg.Type
	bgf.Status = bg.Status
	return
}

//添加一个工具
func (s *sysdata) AddBox(bx def.BoxInfo) (err error) {
	bgf, ok := s.mbgf[bx.Group]
	if !ok {
		bgf = new(boxGroupInfo)
		bgf.Name = bx.Group
		bgf.ShowName = bx.Group
		bgf.Des = "此类别为未被定义的类别，需要联系管理员处理"
		bgf.mbif = make(map[string]*boxInfo)
		err = s.hb.AddGroup(bgf)
		if err != nil {
			return
		}
		s.mbgf[bx.Group] = bgf
	}

	bif, ok := bgf.mbif[bx.Name]
	if !ok {
		bif = new(boxInfo)
		bif.BoxInfo = bx
		bgf.mbif[bx.Name] = bif
	}
	//--信息更新
	bif.Group = bx.Group
	bif.Name = bx.Name
	bif.ShowName = bx.ShowName
	bif.Params = bx.Params
	bif.Author = bx.Author
	bif.Description = bx.Description
	bif.Mode = bx.Mode
	bif.Version = bx.Version
	bif.ApiVersion = bx.ApiVersion
	bif.ModeInfo = bx.ModeInfo

	//重新探测进行连接

	bif.isAlive = true
	bif.connType = bx.Mode
	if bx.Mode == "http" && len(bx.ModeInfo) < 10 {
		bx.ModeInfo = "http://localhost" + bx.ModeInfo
	}
	return
}

//获取可以使用的组 按name排序
func (s *sysdata) GetUserdGroup() (gl GroupList) {
	gl = make(GroupList, len(s.mbgf))
	i := 0
	for _, v := range s.mbgf {
		gl[i] = v
		i++
	}
	sort.Sort(sort.Reverse(gl))
	return

}

//获取组下可用的工具列表 按name排序
func (s *sysdata) GetGroupUsedBox(group string) (bl BoxList) {
	bgf, ok := s.mbgf[group]
	if !ok {
		return
	}
	//bl = make(BoxList, len(bgf.mbif))
	//i := 0
	for _, v := range bgf.mbif {
		if v.isAlive {
			bl = append(bl, v)
		}
		//bl[i] = v
		//i++
	}
	sort.Sort(sort.Reverse(bl))
	return
}

//获得要访问的box
func (s *sysdata) GetCallBox(group, box string) (bx *boxInfo, err error) {
	bgf, ok := s.mbgf[group]
	if ok {
		bx, ok = bgf.mbif[box]
		if ok {
			return
		}
	}
	//错误流程
	err = GROUP_BOX_ERR
	return
}

func (s *sysdata) checkeAlive() {
	sleepTime := time.Second * 3 * 60
	for {
		for _, bgf := range s.mbgf {
			for _, v := range bgf.mbif {
				if v.connType == "http" {
					v.isAlive = pingBox(v.ModeInfo + "/ping")
					if !v.isAlive {
						Log(v.Name, " 心跳检测失效，将不可用")
					}
				}
			}
		}
		time.Sleep(sleepTime)
	}
}

func pingBox(urlstr string) bool {
	data, err := base.RegBoxPost(urlstr, "")
	Log(urlstr, "res:", string(data))
	if err == nil && string(data) == "ok" {
		return true
	}
	return false
}
