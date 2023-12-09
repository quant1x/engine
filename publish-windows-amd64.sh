#!/bin/sh
# 获取当前路径, 用于返回
p0=`pwd`
# 获取脚本所在路径, 防止后续操作在非项目路径
p1=$(cd $(dirname $0);pwd)

source ./publish-compile.sh

# windows amd64
GOOS=windows
GOARCH=amd64
app=stock
ext=.exe

compile $app $ext

cd $p0