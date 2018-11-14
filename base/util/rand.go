//
//  rand.go
//  base
//
//  Created by 吴道睿 on 2018/4/8.
//  Copyright © 2018年 吴道睿. All rights reserved.
//
package util


import(
    "math/rand"
	"time"
)

var (
	BaseRandInt = 4  //用于产生一个int数字
	BaseRandUint64  = int64(1024 * 1024 * 1024 * 1024) //用于产生一个uint64数字的
	BaseString  = []byte(`~!@#$%^&*1234567890-_+=qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM()[]{};:"'\|,<.>/?`)
	baseStringLen = len(BaseString)
)

//设置用来产生默认随机字符串的基础字符串 需要在需要使用前进行设置
func SetBaseString(basestr string){
	BaseString = []byte(basestr)
	baseStringLen = len(BaseString)
	return
}

func GetInt8()int{
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(256)
}

func GetInt(n int)int{
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(n)
}

//获得一个随机数字用于 Rsa
func GetRsaInt()int{
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(BaseRandInt) * 256
}


//获得一个随机数字 产生一个随机组 或随机id
func GetUint64()uint64{
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return uint64(r.Int63n(BaseRandUint64))
}

//获得一个随机字节数组 输入长度
func GetBytes(n int)[]byte{
	bs := make([]byte,n)
	for i:= 0 ;i< n;i++{
		bs[i] = byte(GetInt8())
	}
	return bs
}

//随机bool值
func RandBool()bool{
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(20) >= 10
}

//产生一个随机字符串 输入长度
func GetString(n int)string{
	bs := make([]byte,n)
	for i:= 0 ;i< n;i++{
		bs[i] = BaseString[GetInt(baseStringLen)]
	}
	return string(bs)
}
