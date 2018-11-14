#!/bin/sh
#
#  GetTidLogsKeepStart.sh
#
#  Created by 金山 on 2017-09-27 15:58:58
#

while :
do
#  3秒检测一次
    sleep 3
    RESULT=`ps -e|grep 'GetTidLogs'|sed -e "/grep/d"` 
    if [ -z "$RESULT" ];then 
       # 启动进程 默认http模式启动 
       (./GetTidLogs -mode http >>err.dat &)
       echo "$(date)  _restart_" >> GetTidLogsRestart.log
    fi
done
