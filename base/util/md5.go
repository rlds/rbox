//
//  md5.go
//  base
//
//  Created by 吴道睿 on 2018/4/7.
//  Copyright © 2018年 吴道睿. All rights reserved.
//

package util


import(
       "crypto/md5"
       "encoding/hex"
)

//
func GetMd5Str(b []byte)string{
    h := md5.New()
    h.Write(b)
    return hex.EncodeToString(h.Sum(nil))
}
