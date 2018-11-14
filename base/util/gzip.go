//
//  gzip.go
//  base
//
//  Created by 吴道睿 on 2018/4/6.
//  Copyright © 2018年 吴道睿. All rights reserved.
//

package util

import (
    "bytes"
    "compress/gzip"
	"io/ioutil"
)

func GzipEncode(in []byte) ([]byte, error) {
    var (
        buffer bytes.Buffer
        out    []byte
        err    error
    )
    writer := gzip.NewWriter(&buffer)
    _, err = writer.Write(in)
    if err != nil {
        writer.Close()
        return out, err
    }
    err = writer.Close()
    if err != nil {
        return out, err
    }
    return buffer.Bytes(), nil
}

func GzipDecode(in []byte) ([]byte, error) {
    reader, err := gzip.NewReader(bytes.NewReader(in))
    if err != nil {
        var out []byte
        return out, err
    }
    defer reader.Close()

    return ioutil.ReadAll(reader)
}
