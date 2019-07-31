package db

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/rlds/rbox/base"
	"github.com/rlds/rbox/base/def"
	butil "github.com/rlds/rbox/base/util"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/util"
)

// DataList 数据列表返回
type DataList = def.DataList

// DataItem 数据单元
type DataItem = def.DataItem

// LevelDbOptions 配置参数信息
type LevelDbOptions = def.LevelDbOptions

// DbInfo  库的信息
// 用于使用这些信息重启库
type DbInfo struct {
	db     *leveldb.DB  // 库指针
	autoId uint64       // 库自增序号
	mu     sync.RWMutex // 自增序号锁
	def.DbInfo
}

// ListInfo 列表信息
// 定义： datakey idxkey 采用自增 uint64数值 占用8字节
// 自增值存储于 sysdb中
//
type ListInfo struct {
	db        *leveldb.DB  // 库指针
	mu        sync.RWMutex // 自增序号锁
	autoIdPre []byte       // 自增序号索引key
	def.ListInfo
}

var defaultLevelDbOption = LevelDbOptions{
	BlockCacheCapacity:   opt.MiB * 800, //
	BlockSize:            opt.KiB * 800,
	CompactionTableSize:  opt.MiB * 200, //单个数据文件的最大大小
	CompactionTotalSize:  opt.MiB * 1000,
	IteratorSamplingRate: opt.MiB * 200,
	WriteBuffer:          opt.MiB * 400, //日志的最大大小
}

// LevelDbList 列表的处理
type LevelDbList struct {
	sysPath  string
	sysdb    *leveldb.DB          //全局存储库
	dbinfo   map[string]*DbInfo   //已有库的信息记录
	listinfo map[string]*ListInfo //列表信息 key 为 db+listname 需要listname在db中唯一
}

// Init 初始化加载
// 初始化步骤
// 打开系统库，
// 恢复库信息
// 恢复列表信息
// 检测最大自增序号
func (ll *LevelDbList) Init(name, sysPath string, opts interface{}) (err error) {
	ll.dbinfo = make(map[string]*DbInfo)
	ll.listinfo = make(map[string]*ListInfo)

	var lvbopt *LevelDbOptions
	if opts == nil {
		lvbopt = &defaultLevelDbOption
	} else {
		lvbopt = opts.(*LevelDbOptions)
	}
	ll.sysPath = sysPath
	ll.sysdb, err = leveldb.OpenFile(sysPath+name, &opt.Options{
		BlockCacheCapacity:   lvbopt.BlockCacheCapacity, //
		BlockSize:            lvbopt.BlockSize,
		CompactionTableSize:  lvbopt.CompactionTableSize, //单个数据文件的最大大小
		CompactionTotalSize:  lvbopt.CompactionTotalSize,
		IteratorSamplingRate: lvbopt.IteratorSamplingRate,
		WriteBuffer:          lvbopt.WriteBuffer, //日志的最大大小
	})

	if err != nil {
		return
	}

	// 打开其他库
	ll.checkDb()

	// 加载所有列表
	ll.checkDbList()
	return nil
}

func dbSeek(db *leveldb.DB, seekkey []byte, dofunc func(key, val []byte) bool) {
	if db == nil {
		return
	}
	rang := util.BytesPrefix(seekkey)
	iter := db.NewIterator(rang, nil)
	keyLen := len(seekkey)
	for iter.Next() {
		// 用于控制返回数量 当前未控制
		if dofunc(iter.Key()[keyLen:], iter.Value()) {
			break
		}
	}
	iter.Release()
}

func (dbi *DbInfo) getAutoId() {
	idKey := butil.Uint64ToByte(1)
	data, err := dbi.db.Get(idKey, nil)
	if err == nil {
		dbi.autoId = butil.ByteToUint64(data)
	} else {
		dbi.autoId = AutoIdStartNum
	}
}

// 每次取2个id 用于列表自行使用
func (dbi *DbInfo) getListPreId() (id uint64) {
	dbi.mu.Lock()
	defer dbi.mu.Unlock()
	id = dbi.autoId
	dbi.autoId += 2
	dbi.setAutoId()
	return
}

func (dbi *DbInfo) setAutoId() {
	idKey := butil.Uint64ToByte(1)
	dbi.db.Put(idKey, butil.Uint64ToByte(dbi.autoId), nil)
}

// 每次使用2个id
func (li *ListInfo) getPreKeys(id uint64) {
	li.DataPre = butil.Uint64ToByte(id)
	li.IdxPre = butil.Uint64ToByte(id + 1)
}

// 检查建立库信息
func (ll *LevelDbList) checkDb() {
	dbinfoKey := butil.Uint64ToByte(3)
	dbSeek(ll.sysdb, dbinfoKey, func(key, val []byte) bool {
		dbi := new(DbInfo)
		err := json.Unmarshal(val, dbi)
		if err == nil {
			// 库迁移时数据路径可能会变更需要更新为启动设定路径
			dbi.DataPath = ll.sysPath
			dopt := ll.getDbOption(dbi.DbName)
			err = ll.addDb(dbi, dopt)
			if err != nil {
				base.Log(err, "---", butil.ObjToStr(dopt), "---", butil.ObjToStr(dbi))
				return false
			}
			dbi.getAutoId()
			ll.dbinfo[dbi.DbName] = dbi
		} else {
			// 记录日志
		}
		return false
	})
}

func (ll *LevelDbList) getDbOption(dbName string) (dopt *LevelDbOptions) {
	dopt = new(LevelDbOptions)
	key := butil.Uint64ToByte(4)
	key = append(key, []byte(dbName)...)
	dat, err := ll.sysdb.Get(key, nil)
	if err == nil {
		err = json.Unmarshal(dat, dopt)
	}
	if err != nil {
		dopt = &defaultLevelDbOption
	}
	return
}

// 存储库的记录信息
func (ll *LevelDbList) saveDbinfo(dbi *DbInfo, dopt *LevelDbOptions) {
	dat, err := json.Marshal(dbi)
	if err == nil {
		key := butil.Uint64ToByte(3)
		key = append(key, []byte(dbi.DbName)...)
		ll.sysdb.Put(key, dat, nil)
	}
	// 存储设置信息
	dat, err = json.Marshal(dopt)
	if err == nil {
		key := butil.Uint64ToByte(4)
		key = append(key, []byte(dbi.DbName)...)
		ll.sysdb.Put(key, dat, nil)
	}
}

// 列表信息
func (ll *LevelDbList) checkDbList() {
	// autoKey
	key := butil.Uint64ToByte(3)
	for dbname, dbi := range ll.dbinfo {
		dbSeek(dbi.db, key, func(ikey, ival []byte) bool {
			li := new(ListInfo)
			err := json.Unmarshal(ival, li)
			if err == nil {
				li.db = dbi.db
				li.setAutoIdPre()
				li.getAutoId()
				ll.listinfo[dbname+li.Name] = li
			}
			return false
		})
	}
}

// 添加库 不做存在检查
func (ll *LevelDbList) addDb(dbi *DbInfo, dbopt *LevelDbOptions) (err error) {
	ll.dbinfo[dbi.DbName] = dbi
	if nil == dbopt {
		dbopt = &defaultLevelDbOption
	}
	if len(dbi.DataPath) < 5 {
		dbi.DataPath = ll.sysPath
	}
	dbi.DbOpt = defaultLevelDbOption
	dbi.db, err = leveldb.OpenFile(dbi.DataPath+dbi.DbName, &opt.Options{
		BlockCacheCapacity:   dbopt.BlockCacheCapacity, //
		BlockSize:            dbopt.BlockSize,
		CompactionTableSize:  dbopt.CompactionTableSize, //单个数据文件的最大大小
		CompactionTotalSize:  dbopt.CompactionTotalSize,
		IteratorSamplingRate: dbopt.IteratorSamplingRate,
		WriteBuffer:          dbopt.WriteBuffer, //日志的最大大小
	})
	ll.saveDbinfo(dbi, dbopt)
	base.Log("open db ok:", dbi.DataPath+dbi.DbName)
	return
}

// 添加列表 不做存在检查
func (ll *LevelDbList) addList(li *ListInfo) (err error) {
	oli, ok := ll.listinfo[li.DbName+li.Name]
	if ok { // 列表存在
		oli.Description = li.Description
		oli.ListType = li.ListType
		oli.Status = li.Status
		oli.saveInfo()
		return
	}
	// 检查库是否存在
	dbi, ok := ll.dbinfo[li.DbName]
	if !ok {
		dbi = new(DbInfo)
		dbi.DbName = li.DbName
		dbi.Description = getDateStr()
		err = ll.AddDb(dbi, nil)
		if err != nil {
			err = ErrDBNOTFUNDERROR
			return
		}
	}
	// 库在
	li.db = dbi.db
	li.setAutoIdPre()
	li.getAutoId()
	// 设置前缀值
	if len(li.DataPre) < 9 {
		preid := dbi.getListPreId()
		li.getPreKeys(preid)
	}
	// 存储信息
	li.saveInfo()
	ll.listinfo[li.DbName+li.Name] = li
	//base.Log("li.DbName+li.Name = ", li.DbName+li.Name, *li)
	return
}

//  存储信息
func (li *ListInfo) saveInfo() {
	key := butil.Uint64ToByte(3)
	key = append(key, []byte(li.Name)...)
	dat, err := json.Marshal(li)
	if err == nil {
		li.db.Put(key, dat, nil)
	}
}

// 设置前缀
func (li *ListInfo) setAutoIdPre() {
	li.autoIdPre = butil.Uint64ToByte(2)
	li.autoIdPre = append(li.autoIdPre, []byte(li.Name)...)
}

// 获取自增id
func (li *ListInfo) getAutoId() {
	data, err := li.db.Get(li.autoIdPre, nil)
	if err == nil {
		li.AutoId = butil.ByteToUint64(data)
	} else {
		li.AutoId = 0
	}
}

// 直接返回 []byte 序号从1开始
func (li *ListInfo) getListId() (bid []byte) {
	li.mu.Lock()
	li.AutoId++
	bid = butil.Uint64ToByte(li.AutoId)
	// 存储自动id
	li.db.Put(li.autoIdPre, bid, nil)
	li.mu.Unlock()
	return
}

// 存储数据
func (li *ListInfo) addData(key, val []byte) (err error) {
	likey := li.DataPre
	likey = append(likey, key...)
	dat, err := li.db.Get(likey, nil)
	if len(dat) >= 8 {
		err = li.db.Put(likey, setValue(dat[:8], val), nil)
	} else {
		idxKey := li.IdxPre
		bidx := li.getListId()
		idxKey = append(idxKey, bidx...)
		li.db.Put(idxKey, key, nil)
		err = li.db.Put(likey, setValue(bidx, val), nil)
	}
	return
}

func (li *ListInfo) getDataById(id uint64) (dif DataItem, err error) {
	idxKey := li.IdxPre
	bidx := butil.Uint64ToByte(id)
	idxKey = append(idxKey, bidx...)
	dkey, err := li.db.Get(idxKey, nil)
	if err == nil && len(dkey) > 0 {
		return li.getData(dkey)
	}
	return
}

// 取数据
func (li *ListInfo) getData(key []byte) (dif DataItem, err error) {
	likey := li.DataPre
	likey = append(likey, key...)
	val, err := li.db.Get(likey, nil)
	if len(val) >= 8 {
		dif = getDataItem(key, val)
	}
	return
}

// 删除数据
func (li *ListInfo) delData(key []byte) (err error) {
	likey := li.DataPre
	likey = append(likey, key...)
	val, err := li.db.Get(likey, nil)
	if err == nil {
		li.db.Delete(likey, nil)
		if len(val) > 7 {
			likey = li.IdxPre
			likey = append(likey, val[:8]...)
			li.db.Delete(likey, nil)
		}
	}
	return
}

func setValue(numPos, val []byte) (res []byte) {
	res = append(res, numPos...)
	res = append(res, val...)
	return
}

// 建立数据单元结构
func getDataItem(key, val []byte) (dif DataItem) {
	dif.Key = key
	dif.Id = butil.ByteToUint64(val[:8])
	dif.Val = val[8:]
	return
}

// 查找数据
func (li *ListInfo) findData(key []byte, flag int) (dli DataList, err error) {
	switch flag {
	case KeyOnly: //根据key 获取一个 等同  GetData
		{
			var dif DataItem
			dif, err = li.getData(key)
			if err == nil {
				dli.Data = append(dli.Data, dif)
			}
		}
	case KeyIndex:
		{
			skey := setValue(li.DataPre, key)
			dbSeek(li.db, skey, func(ikey, ival []byte) bool {
				fkey := key
				fkey = append(fkey, ikey...)
				dif := getDataItem(fkey, ival)
				dli.Data = append(dli.Data, dif)
				return false
			})
		}
	case KeyContens:
		{
			dbSeek(li.db, li.DataPre, func(ikey, ival []byte) bool {
				if bytes.Contains(ikey, key) {
					dif := getDataItem(ikey, ival)
					dli.Data = append(dli.Data, dif)
				}
				return false
			})
		}
	case ValueContent:
		{
			dbSeek(li.db, li.DataPre, func(ikey, ival []byte) bool {
				if bytes.Contains(ival, key) {
					dif := getDataItem(ikey, ival)
					dli.Data = append(dli.Data, dif)
				}
				return false
			})
		}
	}
	dli.GetNum = len(dli.Data)
	return
}

// 删除数据
func (li *ListInfo) deleteData(key []byte, flag int) (err error) {
	switch flag {
	case KeyOnly: //根据key 获取一个 等同  GetData
		{
			err = li.delData(key)
		}
	case KeyIndex:
		{
			skey := setValue(li.DataPre, key)
			dbSeek(li.db, skey, func(ikey, ival []byte) bool {
				fkey := key
				fkey = append(fkey, ikey...)
				err = li.delData(fkey)
				return false
			})
		}
	case KeyContens:
		{
			dbSeek(li.db, li.DataPre, func(ikey, ival []byte) bool {
				if bytes.Contains(ikey, key) {
					err = li.delData(ikey)
				}
				return false
			})
		}
	case ValueContent:
		{
			dbSeek(li.db, li.DataPre, func(ikey, ival []byte) bool {
				if bytes.Contains(ival, key) {
					err = li.delData(ikey)
				}
				return false
			})
		}
	}
	return
}

// AddData 添加数据
func (ll *LevelDbList) AddData(dbName, listName string, key, val []byte) (err error) {
	li, ok := ll.listinfo[dbName+listName]
	if !ok {
		//err = ErrLISTNOTFOUND
		//return
		// 库是否存在
		li = new(ListInfo)
		li.DbName = dbName
		li.Name = listName
		err = ll.addList(li)
		if err != nil {
			return
		}
	}
	//base.Log("li ==", *li)
	err = li.addData(key, val)
	return
}

// AddDb 添加库
func (ll *LevelDbList) AddDb(di *DbInfo, options *LevelDbOptions) (err error) {
	db, ok := ll.dbinfo[di.DbName]
	if !ok {
		// 不存在建立及打开
		err = ll.addDb(di, options)
		di.getAutoId()
	} else {
		if options == nil {
			options = &defaultLevelDbOption
		} else {
			db.DbOpt = *options
		}
		// 信息更新
		ll.saveDbinfo(di, options)
		db.Description = di.Description
		db.DbType = di.DbType
		db.Status = di.Status
	}
	return
}

// GetDbs 取库信息
func (ll *LevelDbList) GetDbs() (dbis []*def.DbInfo) {
	for _, dbi := range ll.dbinfo {
		dbis = append(dbis, &dbi.DbInfo)
	}
	return
}

// EditList 更新list信息
func (ll *LevelDbList) EditList(li *ListInfo) error {
	return ll.addList(li)
}

// DeleteDb 删除库
func (ll *LevelDbList) DeleteDb(dbName string) (err error) {
	dbi, ok := ll.dbinfo[dbName]
	if ok {
		// 删除索引
		delete(ll.dbinfo, dbName)
		for k := range ll.listinfo {
			if strings.HasPrefix(k, dbName) {
				delete(ll.listinfo, k)
			}
		}
		// 库关闭
		dbi.db.Close()
	}
	return
}

// GetDbLists 取得库中的列表
func (ll *LevelDbList) GetDbLists(dbName string) (lis []*def.ListInfo) {
	for k, li := range ll.listinfo {
		if strings.HasPrefix(k, dbName) {
			lis = append(lis, &li.ListInfo)
		}
	}
	return
}

// GetDataById 根据id查找数据
func (ll *LevelDbList) GetDataById(dbName, listName string, start, getNum uint64) (data def.DataList, err error) {
	if getNum <= 0 {
		return
	}
	li, ok := ll.listinfo[dbName+listName]
	if !ok {
		err = ErrLISTNOTFOUND
		return
	}
	if start > li.AutoId {
		start = li.AutoId
	}
	n := uint64(0)
	for {
		di, err := li.getDataById(start)
		if err == nil {
			data.Data = append(data.Data, di)
			n++
		}
		if start == 0 || n >= getNum {
			break
		}
		start--
	}
	// dbSeek(li.db, li.IdxPre, func(ikey, ival []byte) bool {
	// 	di, err := li.getData(ival)
	// 	if err == nil {
	// 		data.Data = append(data.Data, di)
	// 		n++
	// 	}
	// 	return n >= getNum
	// })
	data.MaxId = li.AutoId
	data.GetNum = len(data.Data)
	return
}

// GetData 取数据
func (ll *LevelDbList) GetData(dbName, listName string, key []byte) (itm def.DataItem, err error) {
	li, ok := ll.listinfo[dbName+listName]
	if !ok {
		err = ErrLISTNOTFOUND
		return
	}
	itm, err = li.getData(key)
	return
}

// FindData 数据查找
func (ll *LevelDbList) FindData(dbName, listName string, key []byte, flag int) (dli def.DataList, err error) {
	li, ok := ll.listinfo[dbName+listName]
	if !ok {
		err = ErrLISTNOTFOUND
		return
	}
	dli, err = li.findData(key, flag)
	return
}

// DeleteData 删除数据
func (ll *LevelDbList) DeleteData(dbName, listName string, key []byte, flag int) (err error) {
	li, ok := ll.listinfo[dbName+listName]
	if !ok {
		err = ErrLISTNOTFOUND
		return
	}
	err = li.deleteData(key, flag)
	return
}

// RangeAll 遍历
func (ll *LevelDbList) RangeAll(dbname string) (err error) {
	dbi, ok := ll.dbinfo[dbname]
	if !ok {
		err = ErrDBNOTFUNDERROR
		return
	}
	dbSeek(dbi.db, []byte{0}, func(key, val []byte) bool {
		fmt.Println("rangeAll key:", string(key), " keyB:", key, " val:", string(val), " valB:", val)
		return false
	})
	dbSeek(ll.sysdb, nil, func(key, val []byte) bool {
		fmt.Println("rangeAll key:", string(key), " keyB:", key, " val:", string(val), " valB:", val)
		return false
	})
	return
}
