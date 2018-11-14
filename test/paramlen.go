//
//  paramlen.go
//  mbox
//
//  Created by 吴道睿 on 2018/5/16.
//  Copyright © 2018年 吴道睿. All rights reserved.
//
package main

import(
    "fmt"	   
)	   

func main(){
	var data = []string{"a","b","c",
		"d",
		"e",     
	}
	
	fmt.Println(len(data))
	
	for i,v := range data {
		fmt.Println(i," =",v)
	}
}
