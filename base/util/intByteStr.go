//
//  intByteStr.go
//  base
//
//  Created by 吴道睿 on 2018/4/12.
//  Copyright © 2018年 吴道睿. All rights reserved.
//
package util

import(
    "encoding/json"
	"strconv"
)

/*
   字符串、数字、结构体的互相转换
*/
// ObjToBytes 结构体转json
func ObjToBytes(da interface{})[]byte{
	data,_:=json.Marshal(da)
	return data
}

// ObjToStr 结构体转json
func ObjToStr(da interface{})string{
	data,_:=json.Marshal(da)
	return string(data)
}

// StrToInt 字符串转数字
func StrToInt(s string)int{
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return i
}

// IntToStr 
func IntToStr(i int)string{
    return strconv.Itoa(i)
}

// Uint64ToByte Uint64ToByte
func Uint64ToByte(u uint64)(b []byte){
	b=make([]byte,8)
	b[7]=(byte)(u)
	b[6]=(byte)(u>>8 )
	b[5]=(byte)(u>>16)
	b[4]=(byte)(u>>24)
	b[3]=(byte)(u>>32)
	b[2]=(byte)(u>>40)
	b[1]=(byte)(u>>48)
	b[0]=(byte)(u>>56)
	return
}

// ByteToUint64 ByteToUint64
func ByteToUint64(b []byte)(r uint64){
	r =(uint64)(b[7])
	r|=(uint64)(b[6]) <<8
	r|=(uint64)(b[5]) <<16
	r|=(uint64)(b[4]) <<24
	r|=(uint64)(b[3]) <<32
	r|=(uint64)(b[2]) <<40
	r|=(uint64)(b[1]) <<48
	r|=(uint64)(b[0]) <<56
	return
}

