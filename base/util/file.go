//
//  file.go
//  base
//
//  Created by 吴道睿 on 2018/4/8.
//  Copyright © 2018年 吴道睿. All rights reserved.
//
package util

import (
	"bytes"
	"os"
)

// GetAllFileData 读取完整文件内容（小文件）
func GetAllFileData(filepath string) []byte {
	f, err := os.Open(filepath)
	if err != nil {
		return nil
	}
	var n int64

	if fi, err2 := f.Stat(); err2 == nil {
		if size := fi.Size(); size < 1e9 {
			n = size
		}
	} else {
		return nil
	}
	buf := bytes.NewBuffer(make([]byte, 0, n+bytes.MinRead))
	defer buf.Reset()
	_, err = buf.ReadFrom(f)
	f.Close()
	if err != nil {
		return nil
	}
	return buf.Bytes()
}

//SaveReplaceFile 存储并替换文件内容
func SaveReplaceFile(path string, data []byte) error {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	f.Write(data)
	f.Close()
	return nil
}

//TestAndCreateDir  检测并创建文件夹
func TestAndCreateDir(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		err = os.MkdirAll(path, 0755)
		return err == nil
	}
	return false
}

func MkAlldir(dir string) bool {
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		//rlog.V(1).Info(HC+"建立文件夹出现错误["+err.Error()+"]"+EC)
		return false
	}
	return true
}

func DelFile(path string) error {
	return os.Remove(path)
}

func DelDir(path string) error {
	return os.RemoveAll(path)
}

func GetSubPath(dir, path string) string {
	return path[len(dir):]
}
