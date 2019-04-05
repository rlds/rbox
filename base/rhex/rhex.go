package rhex

//  rhex.go
//  rbox
//
//  Created by 吴道睿 on 2018/4/12.
//  Copyright © 2018年 吴道睿. All rights reserved.
//

import (
	"bytes"
	"sync"
)

/*
 n进制计数器
 目前实现的是 自增  可自动增加进位
 加数运算           不可增加位数
*/

const (
	//baseStr = "!0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ_abcdefghijklmnopqrstuvwxyz"
	//baseStr = "!#$,.0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ_abcdefghijklmnopqrstuvwxyz"
	//baseStr = "0123456789"  //使用十进制来测试
	//baseStr = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"  //62
	baseStr = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijkmnpqrstuvwxyz" //60
	//id反序用于保证生成id时间最近的序号小
	//baseStr = "zyxwvutsrqponmlkjihgfedcbaZYXWVUTSRQPONMLKJIHGFEDCBA9876543210"

	// RetPaoqi 初始化数据又抛弃
	RetPaoqi = 3
	// RetYueJie 数值容量超过了初始化位数
	RetYueJie = 2
)

var (
	_jinzhi   int    //进制数
	_jinzhi64 uint64 //进制数62位
	rBaseStr  = []byte(baseStr)
)

func init() {
	_jinzhi = len(baseStr)
	_jinzhi64 = (uint64)(_jinzhi)
}

// RHex64w RHex64w
type RHex64w struct {
	B  []int  //数值
	W  int    //数字位数
	BB []byte //基准字符串
	L  *sync.RWMutex
	J  bool //是否允许自动增加位数
}

// StrInit 初始化
//    w     位数
//    bs    初始化数值
//    j     是否自动进位
func (h *RHex64w) StrInit(w int, bs []byte, j bool) {
	h.W = w
	h.B = make([]int, w)
	h.BB = []byte(baseStr)
	h.L = new(sync.RWMutex)
	h.J = j
	for i := 0; i < w && i < len(bs); i++ {
		h.B[i] = bytes.IndexByte(h.BB, bs[i])
	}
	return
}

// Add 并发执行量
// 212000  4核 1G 内存
// 412000      2G 内存
func (h *RHex64w) Add() (res []byte) {
	h.L.Lock()
	h.carr(0)
	for i := 0; i < h.W; i++ {
		res = append(res, h.BB[h.B[i]])
	}
	h.L.Unlock()
	return
}

// AddBytes n进制数加
func (h *RHex64w) AddBytes(str []byte) (res []byte, ret int) {
	h.L.Lock()
	dd := make([]int, h.W)
	cc := 0
	slen := len(str)
	if slen > h.W {
		slen = h.W
		ret = RetPaoqi //有抛弃的数据
	}
	for i := 0; i < slen; i++ {
		dd[i] = bytes.IndexByte(h.BB, str[i])
		t := dd[i] + h.B[i] + cc
		h.B[i] = t % _jinzhi
		cc = t / _jinzhi
	}

	if slen < h.W {
		h.B[slen] += cc
	} else if cc > 0 {
		ret = RetYueJie
	}

	for i := 0; i < h.W; i++ {
		res = append(res, h.BB[h.B[i]])
	}
	h.L.Unlock()
	return
}

// Adduint64 n进制数与10进制相加
func (h *RHex64w) Adduint64(u uint64) (res []byte, ret int) {
	h.L.Lock()
	dd := make([]int, h.W)
	for i := 0; i < h.W; i++ {
		dd[i] = int(u % _jinzhi64)
		u = u / _jinzhi64
	}
	aaa := 0
	for i := 0; i < h.W; i++ {
		t := h.B[i] + dd[i] + aaa
		h.B[i] = t % _jinzhi
		aaa = t / _jinzhi
	}
	if aaa > 0 {
		ret = RetYueJie
	}
	for i := 0; i < h.W; i++ {
		res = append(res, h.BB[h.B[i]])
	}
	h.L.Unlock()
	return
}

//进位计算 只有 Add使用 不加锁
func (h *RHex64w) carr(n int) {
	if n >= h.W {
		if h.J {
			h.W += 1
			h.B = append(h.B, 0)
			h.BB = append(h.BB, 0)
		} else {
			return
		}
	}
	h.B[n]++
	if h.B[n] >= _jinzhi {
		h.B[n] = 0
		h.carr(n + 1)
	}
}

//IntToNbyte 输入数字必须小于进位数
func IntToNbyte(i int) byte {
	return rBaseStr[i]
}

//ByteToInt 查询字符对应的进制数值
func ByteToInt(b byte) int {
	for i := 0; i < _jinzhi; i++ {
		if b == rBaseStr[i] {
			return i
		}
	}
	return -1
}

//GetHexStr 返回基础字符串
func GetHexStr() string {
	return baseStr
}

//ToBytes 输出数值从低位到高位
func (h *RHex64w) ToBytes() (res []byte) {
	h.L.RLock()
	for i := 0; i < h.W; i++ {
		res = append(res, h.BB[h.B[i]])
	}
	h.L.RUnlock()
	return
}

//ToRelBytes 返回真实的数值串 从高位到低位
func (h *RHex64w) ToRelBytes() (res []byte) {
	h.L.RLock()
	for i := h.W - 1; i >= 0; i-- {
		res = append(res, h.BB[h.B[i]])
	}
	h.L.RUnlock()
	return
}

//ReSet 重置
func (h *RHex64w) ReSet() {
	h.L.Lock()
	for i := 0; i < h.W; i++ {
		h.B[i] = 0
	}
	h.L.Unlock()
}
