package def

// DataList 数据列表返回
type DataList struct {
	GetNum int
	MaxId  uint64
	Data   []DataItem
}

// DataItem 数据单元
type DataItem struct {
	Key []byte // 名字
	Val []byte // 值
	Id  uint64 // 编号
	//	AddTime uint64 // 数据添加时间
}

// LevelDbOptions 配置参数信息
type LevelDbOptions struct {
	BlockCacheCapacity   int // opt.MiB * 800,  //
	BlockSize            int // opt.KiB * 800,
	CompactionTableSize  int // opt.MiB * 200,  //单个数据文件的最大大小
	CompactionTotalSize  int // opt.MiB * 1000,
	IteratorSamplingRate int // opt.MiB * 200,
	WriteBuffer          int // opt.MiB * 400, //日志的最大大小
}

type DbInfo struct {
	CreateTime uint64 // 创建时间
	DataPath   string // 库路径
	DbName     string // 库名称
	//  其他信息使用健值对存储
	Description string         // 库描述
	DbType      int            // 库类型
	Status      int            // 库的状态
	DbOpt       LevelDbOptions // 库设置信息
}

// ListInfo 列表信息
// 定义： datakey idxkey 采用自增 uint64数值 占用8字节
// 自增值存储于 sysdb中
//
type ListInfo struct {
	AutoId      uint64 // 自增序号
	DataPre     []byte // 数据索引前缀
	IdxPre      []byte // 键值索引前缀
	DbName      string // 库名称
	Name        string // 列表名称 需要listname在db中唯一
	Description string // 列表描述
	ListType    int    // 列表类型
	Status      int    // 列表的状态
}

// 返回数据结构
type ReturnData struct {
	Code    int
	Message string
	Data    interface{}
}
