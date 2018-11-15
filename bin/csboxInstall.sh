#!/bin/sh
#
# 自动复制并安装csbox至指定文件夹下
#

installPath=$1
if [ -z "$installPath" ];then 
   echo "请输入将要安装csbox的文件夹路径"
else
   echo "csbox将安装的文件夹路径为："${installPath}
   echo "开始清理已安装部分"
   rm -rf ${installPath}/bin/csbox
   rm -rf ${installPath}/static
   echo "清理完成"
   mkdir -p {${installPath}/log,${installPath}/bin}
   cp csbox ${installPath}/bin/
   cp -rf ../csbox/static ${installPath}
   echo "安装完成可以开始使用了 csbox路径为："${installPath}/bin/csbox
fi
