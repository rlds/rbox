//
//  def.go
//  base
//
//  Created by 吴道睿 on 2018/4/20.
//  Copyright © 2018年 吴道睿. All rights reserved.
//
package def

/*
   基础的结构定义
*/

type (
	//配置信息的定义结构
	BoxInfo struct {
		Group    string //工具分组代码
		Name     string //工具名称
		ShowName string //工具展示名称
		//工具的输入参数描述
		SubBox      []SubBoxInfo //输入参数数组
		Author      string       //作者
		Description string       //工具的描述
		Mode        string       //启动模式
		Version     string       //工具版本
		ApiVersion  string       //使用的api版本

		//对应模式的连接方式
		//若为http则为http地址端口
		//若为nats则为nats监听topic
		ModeInfo string
	}

	SubBoxInfo struct {
		SubName string     // 子功能模块名
		Label   string     // 子功能模块名称
		Des     string     // 子功能模块功能描述
		Params  []BoxParam //输入参数数组
	}

	/*
	   工具的输出参数描述结构
	*/
	BoxParam struct {
		Hint      string `json:"hint"`      //此参数用户输入时的帮助说明文字
		Label     string `json:"label"`     //产生输入框的名字标题
		Name      string `json:"name"`      //需要输入参数的 参数名
		Type      string `json:"type"`      //用户需要输入的类型 密码
		Reg       string `json:"reg"`       //前端校验输入参数的正则
		Value     string `json:"value"`     //参数的值 用于单选 复选框时
		ValueType string `json:"valuetype"` //输入参数值类型
		Idx       string `json:"idx"`       //排序编号 1,1 表示 1行1列 用','分割行列
	}

	//输出结果
	BoxOutPut struct {
		Type      string      //输出结果的类型
		Code      string      //返回码
		ReturnMsg string      //返回的错误描述，只有Code!=0时会给出此输出
		TaskId    string      //任务编号
		IsSync    bool        //是否同步返回结果 若为异步返回结果则需调用状态获取拿到结果
		Status    string      //异步执行时任务状态
		Data      interface{} //输出结果的内容
	}

	InputData struct {
		//		RemoteInfo string                 // 请求来源信息 默认传入请求ip地址
		SubBoxName string                 // 数据类型
		Data       map[string]interface{} // 数据内容
	}

	//call 的信息
	RequestIn struct {
		From    string //从那来的任务 描述请求者信息
		Call    string //需要调用那个功能模块
		TaskId  string //任务id编号
		Session string //用于权限校验的session 暂时不用
		//接受格式类型 取消不考虑这个问题
		//AcceptType string          //接受的结果返回格式 markdown|json;default|html|text
		Input InputData //任务请求参数
	}
)

// InfoOk 检测基础信息是否完整
func (b *BoxInfo) InfoOk() bool {
	if b.Group == "" || b.Name == "" || b.Mode == "" || b.ModeInfo == "" {
		return false
	}
	return true
}
