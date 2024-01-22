#!/bin/sh
# 获取当前路径, 用于返回
p0=`pwd`
# 获取脚本所在路径, 防止后续操作在非项目路径
p1=$(cd $(dirname $0);pwd)
cd $p1

source ./publish-compile.sh

# darwin arm64
GOOS=darwin
GOARCH=arm64
app=stock
ext=

compile $app $ext

cd $p0