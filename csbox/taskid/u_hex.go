//
//  u_hex.go
//  rbox
//
//  Created by 吴道睿 on 2018/4/12.
//  Copyright © 2018年 吴道睿. All rights reserved.
//
package taskid

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
	//base_str = "!0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ_abcdefghijklmnopqrstuvwxyz"
	//base_str = "!#$,.0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ_abcdefghijklmnopqrstuvwxyz"
	//base_str = "0123456789"  //使用十进制来测试
	base_str = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	//id反序用于保证生成id时间最近的序号小
	//base_str = "zyxwvutsrqponmlkjihgfedcbaZYXWVUTSRQPONMLKJIHGFEDCBA9876543210"
	RET_PAOQI  = 3   //初始化数据又抛弃
	RET_YUEJIE = 2   //数值容量超过了初始化位数
	RET_ERR    = 111 //返回错误
	RET_OK     = 1   //正常返回
)

var (
	JINZHI     int    //进制数
	JINZHI_64  uint64 //进制数62位
	b_base_str = []byte(base_str)
)

func init() {
	JINZHI = len(base_str)
	JINZHI_64 = (uint64)(JINZHI)
}

type Hex64_w struct {
	B  []int  //数值
	W  int    //数字位数
	BB []byte //基准字符串
	L  *sync.RWMutex
	J  bool //是否允许自动增加位数
}

/*
   初始化
   w     位数
   bs    初始化数值
   j     是否自动进位
*/
func (h *Hex64_w) StrInit(w int, bs []byte, j bool) {
	h.W = w
	h.B = make([]int, w)
	h.BB = []byte(base_str)
	h.L = new(sync.RWMutex)
	h.J = j
	for i := 0; i < w && i < len(bs); i++ {
		h.B[i] = bytes.IndexByte(h.BB, bs[i])
	}
	return
}

/*
 //并发执行量
 212000  4核 1G 内存
 412000      2G 内存
*/
func (h *Hex64_w) Add() (res []byte) {
	h.L.Lock()
	h.carr(0)
	for i := 0; i < h.W; i++ {
		res = append(res, h.BB[h.B[i]])
	}
	h.L.Unlock()
	return
}

//n进制数加
func (h *Hex64_w) AddBytes(str []byte) (res []byte, ret int) {
	h.L.Lock()
	dd := make([]int, h.W)
	cc := 0
	s_len := len(str)
	if s_len > h.W {
		s_len = h.W
		ret = RET_PAOQI //有抛弃的数据
	}
	for i := 0; i < s_len; i++ {
		dd[i] = bytes.IndexByte(h.BB, str[i])
		t := dd[i] + h.B[i] + cc
		h.B[i] = t % JINZHI
		cc = t / JINZHI
	}

	if s_len < h.W {
		h.B[s_len] += cc
	} else if cc > 0 {
		ret = RET_YUEJIE
	}

	for i := 0; i < h.W; i++ {
		res = append(res, h.BB[h.B[i]])
	}
	h.L.Unlock()
	return
}

//n进制数与10进制相加
func (h *Hex64_w) Adduint64(u uint64) (res []byte, ret int) {
	h.L.Lock()
	dd := make([]int, h.W)
	for i := 0; i < h.W; i++ {
		dd[i] = int(u % JINZHI_64)
		u = u / JINZHI_64
	}
	aaa := 0
	for i := 0; i < h.W; i++ {
		t := h.B[i] + dd[i] + aaa
		h.B[i] = t % JINZHI
		aaa = t / JINZHI
	}
	if aaa > 0 {
		ret = RET_YUEJIE
	}
	for i := 0; i < h.W; i++ {
		res = append(res, h.BB[h.B[i]])
	}
	h.L.Unlock()
	return
}

//进位计算 只有 Add使用 不加锁
func (h *Hex64_w) carr(n int) {
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
	if h.B[n] >= JINZHI {
		h.B[n] = 0
		h.carr(n + 1)
	}
}

//输入数字必须小于进位数
func IntToNbyte(i int) byte {
	return b_base_str[i]
}

//查询字符对应的进制数值
func ByteToInt(b byte) int {
	for i := 0; i < JINZHI; i++ {
		if b == b_base_str[i] {
			return i
		}
	}
	return -1
}

//返回基础字符串
func GetHexStr() string {
	return base_str
}

//输出数值从低位到高位
func (h *Hex64_w) ToBytes() (res []byte) {
	h.L.RLock()
	for i := 0; i < h.W; i++ {
		res = append(res, h.BB[h.B[i]])
	}
	h.L.RUnlock()
	return
}

//返回真实的数值串 从高位到低位
func (h *Hex64_w) ToRelBytes() (res []byte) {
	h.L.RLock()
	for i := h.W - 1; i >= 0; i-- {
		res = append(res, h.BB[h.B[i]])
	}
	h.L.RUnlock()
	return
}

//重置
func (h *Hex64_w) ReSet() {
	h.L.Lock()
	for i := 0; i < h.W; i++ {
		h.B[i] = 0
	}
	h.L.Unlock()
}
