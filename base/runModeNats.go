//
//  natsMode.go
//  base
//
//  Created by 吴道睿 on 2018/4/6.
//  Copyright © 2018年 吴道睿. All rights reserved.
//
package base

import(
	"time"
	nats "github.com/nats-io/go-nats"
	."./def"
)

/*
   暂未测试通过
*/

/*
    nats模式的执行
 */
const(
	NatsMaxReconTimes = 5
	NatsReconWaitTime = 1 * time.Minute
)

var (
	//topic key 前缀
    natsTopicPreKey  = "base."
)

type natsModeWorker struct{
	clt *nats.Conn
}

//连接至server
func (n *natsModeWorker)conn()(err error){
	n.clt,err = nats.Connect(gbox.cfg.NatsServerList,
					nats.UserInfo(gbox.cfg.NatsServerUserName, gbox.cfg.NatsServerUserPassword),
					nats.MaxReconnects(NatsMaxReconTimes),
					nats.ReconnectWait(NatsReconWaitTime))
	return
}

func (n *natsModeWorker)Register(){
	var boxinfo BoxInfo
	boxinfo = gbox.cfg.BoxInfo
	Log(boxinfo)
	Log(gbox.cfg.Group,gbox.cfg.Name)
	//发出自己的消息
}

//启动执行
func (n *natsModeWorker)Run(){
	n.conn()
	n.Register()
	
	//启动消息监听
	
	return
}
