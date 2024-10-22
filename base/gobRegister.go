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
	gob.Register(BoxConfig{})
	gob.Register(map[string]uint64{})
	gob.Register(map[string][]uint64{})
	gob.Register(map[string][]string{})
	gob.Register(map[string]string{})
	gob.Register([]def.Pinfo{})
	gob.Register([]*def.Pinfo{})
	gob.Register([]def.Pcharge{})
	gob.Register([]*def.Pcharge{})
	gob.Register([]def.BillInfo{})
	gob.Register([]*def.BillInfo{})

	// gob.Register(&def.ReturnData{})
}
