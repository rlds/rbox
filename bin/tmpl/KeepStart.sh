#!/bin/sh
#
#  {{.BoxConf.Name}}KeepStart.sh
#
#  Created by {{.BoxConf.Author}} on {{.Time}}
#

while :
do
#  3秒检测一次
    sleep 3
    RESULT=`ps -e|grep '{{.BoxConf.Name}}'|sed -e "/grep/d"` 
    if [ -z "$RESULT" ];then 
       # 启动进程 默认http模式启动 
       (./{{.BoxConf.Name}} -mode http >>err.dat &)
       echo "$(date)  _restart_" >> {{.BoxConf.Name}}Restart.log
    fi
done
