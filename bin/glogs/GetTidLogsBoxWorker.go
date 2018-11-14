//  
//  GetTidLogsBoxWorker.go
//  main
//
//  Created by 金山 on 2017-09-27 15:58:58
//
package main

import(
	"wacaispider/mbox"
	."wacaispider/mbox/boxdef"
)

//执行结构
type GetTidLogsBox struct{
	/*
	    注意此结构为全局结构
	    若每次任务有共同数据信息注意区分不同部分防止数据互串
	*/
}

/*
    参数说明：
    taskid    任务的编号id
    input     任务的输入参数
    注意:
 	    此函数仅接收输入信息需要在这里对数据进行处理
        数据的最终获得结果由 Output() 给出
        需要主动记录任务id及最终结果以备输出
*/
func (g *GetTidLogsBox)DoWork(taskid string,input map[string]string)(err error){
    
    return 
}

/*
    参数说明:
	taskid  任务的编号id
    说明：
        外部获得任务的结果
*/
func (g *GetTidLogsBox)Output(taskid string)(m BoxOutPut){
	
	// 返回结果格式
	// 允许自定义格式但需要外部支持展示使用
    m.Type      = mbox.OutputType_Markdown
	
	//返回状态码 OutputRetuen_Success 表示成功 其他表示执行失败
    //m.Code      = mbox.OutputRetuen_Error 
    m.Code      = mbox.OutputRetuen_Success
	
	// 本次结果返回的描述文本信息
    m.ReturnMsg = ""
	
	//任务id 任务的执行结果需要根据任务id进行提取输出
    m.TaskId    = taskid
    
	//要返回的数据信息内容
    //m.Data = data
    return
}

