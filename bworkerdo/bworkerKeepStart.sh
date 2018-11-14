#!/bin/sh
#
#  bworkerKeepStart.sh
#
#  Created by wdr on 2018-09-04 10:39:31
#

while :
do
#  3秒检测一次
    sleep 3
    RESULT=`ps -e|grep 'bworker'|sed -e "/grep/d"` 
    if [ -z "$RESULT" ];then 
       # 启动进程 默认http模式启动 
       (./bworker -mode http >>err.dat &)
       echo "$(date)  _restart_" >> bworkerRestart.log
    fi
done
