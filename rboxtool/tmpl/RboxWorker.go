//  
//  {{.BoxConf.Name}}RboxWorker.go
//  main
//
//  Created by {{.BoxConf.Author}} on {{.Time}}
//
package main

import(
	"github.com/rlds/rbox/base"
	."github.com/rlds/rbox/base/def"
	"sync"
	"time"
	   
	/*
		下面为需要添加引入的包
	 */
	//------------------- 需要添加代码的位置 -----------------
	//markdown
	//"github.com/microcosm-cc/bluemonday" //去除html中不安全的代码
    "gopkg.in/russross/blackfriday.v2"
	//-----------------------------------------------------
)

//执行结构
type {{.BoxConf.Name}}Box struct{
	/*
	    注意此结构为全局结构
	    若每次任务有共同数据信息注意区分不同部分防止数据互串
	*/
	taskIdInfoMap            sync.Map   // 任务信息结果记录
	modeName              string     // 任务执行模式
	isCommandMode         bool       // 是否是命令行模式(只执行一次一个任务)
	cleanTaskTimeStep     int64      // 清理任务记录信息的时间间隔 单位秒
	taskInfoStoreTimeStep int64      // 任务信息存放最低多久时间 单位秒
	lastCleanStartTime    int64      // 最后一次执行清理数据的时间
}

type taskResData struct{
	BoxOutPut
	startTime int64
	endTime   int64
}

const(
	  C_CleanTaskTimeStep     =  60 * 3 //3分钟
	  C_TaskInfoStoreTimeStep =  60 * 5 //5分钟
	  )

//执行任务
func (l *taskResData)Run(in map[string]string){
	/*
	 任务开始的一些设置
	 */
	l.IsSync = false
	l.Status = "PROGRESS"  //执行中
	l.startTime = time.Now().Unix()
	l.Data = l.Status
	
	/*
	 这里添加任务执行代码
	 */
	
	//------------------- 需要添加代码的位置 -----------------
	rets := "   执行结果"
	//-----------------------------------------------------
	
	//任务执行完成的处理
	l.Data = rets
	l.IsSync = true
	l.Status = "COMPLETE"
	l.endTime = time.Now().Unix()
	base.Log(in," ret:",len(rets))
}

func (g *{{.BoxConf.Name}}Box)Init()bool{
	g.modeName , g.isCommandMode = base.GetRunMode()
	g.cleanTaskTimeStep     = C_CleanTaskTimeStep
	g.taskInfoStoreTimeStep = C_TaskInfoStoreTimeStep
	g.lastCleanStartTime    = time.Now().Unix()
	if g.isCommandMode {
		go g.cleanTaskInfo()
	}
	/*
	   下面这里执行box的初始化
	*/
	//------------------- 需要添加代码的位置 -----------------
	
	
	//-----------------------------------------------------
	return true
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
func (g *{{.BoxConf.Name}}Box)DoWork(taskid string,input map[string]string)(err error){
	/*
	 这里可以先执行然后存储结果
	 也可以在此记录输入在其他步骤中执行
	 还可以移步执行存储移步执行关键指针最终存储结果
	 */
	_ ,ok := g.taskIdInfoMap.Load(taskid) 
	if !ok { //
		tskdo := new (taskResData)
		
		//任务开始的一些设置处理 异步执行的任务需要关注
		tskdo.Data = "#   start"
		g.taskIdInfoMap.Store(taskid,tskdo)
		if g.isCommandMode {
			tskdo.Run(input)
		}else{
			go tskdo.Run(input)
		}
	}
    return 
}

/*
    参数说明:
	taskid  任务的编号id
    说明：
        外部获得任务的结果
*/
func (g *{{.BoxConf.Name}}Box)Output(taskid string)(m BoxOutPut){
	
	// 返回结果格式
	// 允许自定义格式但需要外部支持展示使用
	m.Type      = base.OutputType_Markdown
	
	//返回状态码 OutputRetuen_Success 表示成功 其他表示执行失败
	//m.Code      = base.OutputRetuen_Error 
	m.Code      = base.OutputRetuen_Success
	
	// 本次结果返回的描述文本信息
	m.ReturnMsg = ""
	
	//任务id 任务的执行结果需要根据任务id进行提取输出
	m.TaskId    = taskid
	tidInf,ok := g.taskIdInfoMap.Load(taskid)
	if tidInf != nil && ok {
		tidin := tidInf.(*taskResData)
		/*
		 下面为任务结果的一些处理
		 */
		
		/*
		 这里输出box执行结果
		 */
		m.IsSync = tidin.IsSync 
		m.Status = tidin.Status
		if tidin.Data != nil {
			resDat := tidin.Data.(string)
			if g.isCommandMode {
				m.Data = resDat
			}else{
				/*
				    其他模式情况下当数据量过大或涉及到外部计算量过大时的处理
				 */
				//------------------- 需要添加代码的位置 -----------------
				//m.Data = resDat
				m.Type      = "html"
				//base.Log(taskid,"markdown chg")
				m.Data = string(blackfriday.Run([]byte(resDat)))
				//unsafe := blackfriday.Run([]byte(resDat))
				//m.Data  = string(bluemonday.UGCPolicy().SanitizeBytes(unsafe))
				//-----------------------------------------------------
			}
		}else{
			m.Data   = "    data err"
			m.IsSync = true
		}
	}else{
		m.IsSync = true
		m.Data   = "    inmap error"
	}
    return
}


//自动清理过期的数据信息
func (g *{{.BoxConf.Name}}Box)cleanTaskInfo(){
	clean_step := time.Second * time.Duration( g.cleanTaskTimeStep )
	for {
		time.Sleep( clean_step )
		g.lastCleanStartTime  = time.Now().Unix()
		g.taskIdInfoMap.Range(g.rangeTidInfoDo)
	}
}

//清理数据
func (g *{{.BoxConf.Name}}Box)rangeTidInfoDo(key ,val interface{})bool{
	tidinf := val.(*taskResData)
	if (tidinf.endTime + g.taskInfoStoreTimeStep ) > g.lastCleanStartTime {
		//删除任务的记录信息
		g.taskIdInfoMap.Delete(key)
	}
	return true
}
