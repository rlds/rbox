#!/bin/sh
#
#  fboxKeepStart.sh
#
#  Created by 吴道睿 on 2018-11-14 17:51:14
#

while :
do
#  3秒检测一次
    sleep 3
    RESULT=`ps -e|grep 'fbox'|sed -e "/grep/d"` 
    if [ -z "$RESULT" ];then 
       # 启动进程 默认http模式启动 
       (./fbox -mode http >>err.dat &)
       echo "$(date)  _restart_" >> fboxRestart.log
    fi
done
