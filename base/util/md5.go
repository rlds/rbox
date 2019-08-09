//
//  md5.go
//  base
//
//  Created by 吴道睿 on 2018/4/7.
//  Copyright © 2018年 吴道睿. All rights reserved.
//

package util

import (
	"crypto/md5"
	"encoding/hex"
	"hash/crc32"
)

// GetMd5Str MD5
func GetMd5Str(b []byte) string {
	h := md5.New()
	h.Write(b)
	return hex.EncodeToString(h.Sum(nil))
}

func GetMd5Byte(b []byte) []byte {
	h := md5.New()
	h.Write(b)
	return h.Sum(nil)
}

func GetBMd5Int(key []byte) uint32 {
	if len(key) < 64 {
		var scratch [64]byte
		copy(scratch[:], key)
		return crc32.ChecksumIEEE(scratch[:len(key)])
	}
	return crc32.ChecksumIEEE(key)
}
