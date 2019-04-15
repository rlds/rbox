package base

import (
	"encoding/gob"

	"github.com/rlds/rbox/base/def"
)

//var gobObjList []interface{}

// GobRegister gob注册
func GobRegister(obj interface{}) {
	//gobObjList = append(gobObjList, obj)
	gob.Register(obj)
}

func gobRegister() {
	gob.Register(&def.DataList{})
	gob.Register(&def.DataItem{})
	gob.Register(&def.LevelDbOptions{})
	gob.Register(&def.DbInfo{})
	gob.Register(&def.ListInfo{})
	gob.Register(&[]*def.DbInfo{})
	gob.Register(&[]*def.ListInfo{})
	gob.Register(&def.ReturnData{})
	gob.Register(&def.TimeCounterResult{})
	gob.Register(&[]*def.TimeCounterResult{})
	gob.Register(&def.TimeCountRes{})
	gob.Register(&[]def.TimeCountRes{})
	gob.Register(&[]*def.TimeCountRes{})
	// gob.Register(&def.ReturnData{})
}
