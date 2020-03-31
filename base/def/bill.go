package def

type Pinfo struct {
	Pid     string //客户代码
	Email   string //通知报警的邮箱
	AddTime string //添加时间
	Balance int64  //精确到分
	Status  string //可用状态
}

type Pcharge struct {
	Pid    string //客户代码
	Tid    string //充值任务id包含时间信息
	Oper   string //变更人信息
	Des    string //充值描述
	Time   int64  // 时间
	Charge int64  //数额  价格精确到分
}

type BillInfo struct {
	Pid   string // 客户代码
	Tid   string // 充值任务id
	Oper  string // 修改人信息
	Key   string // 产品代码
	Des   string // 描述说明信息
	Price int64  // 价格 单位精确到分
}
