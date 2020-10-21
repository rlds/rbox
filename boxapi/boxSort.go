//
//  mapsort.go
//  rbox
//
//  Created by 吴道睿 on 2018/4/11.
//  Copyright © 2018年 吴道睿. All rights reserved.
//
package main

type GroupList []*boxGroupInfo

func (p GroupList) Len() int           { return len(p) }
func (p GroupList) Less(i, j int) bool { return p[i].Name > p[j].Name }
func (p GroupList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

type BoxList []*boxInfo

func (p BoxList) Len() int           { return len(p) }
func (p BoxList) Less(i, j int) bool { return p[i].Name > p[j].Name }
func (p BoxList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
