#!/bin/sh
# 获取当前路径, 用于返回
#p0=`pwd`
# 获取脚本所在路径, 防止后续操作在非项目路径
#p1=$(cd $(dirname $0);pwd)

#COLOR_NORMAL="\033[0m"
#COLOR_GREEN="\033[1;32m"
#COLOR_YELLOW="\033[1;33m"
#COLOR_RED="\033[1;33m"
#COLOR_GREY="\033[1;30m"
#echo "${COLOR_GREEN} 绿色 ${COLOR_NORMAL}"
#echo "${COLOR_RED} 红色 ${COLOR_NORMAL}"

#golang
echo "----------------< go env >----------------"
GOVERSION=$(go env GOVERSION)
GOOS=$(go env GOOS)
GOARCH=$(go env GOARCH)
echo "   GOOS: ${GOOS}"
echo " GOARCH: ${GOARCH}"
echo "version: ${GOVERSION:2}"

export GO111MODULE=auto
export GOPRIVATE=gitee.com
export GOPROXY=https://goproxy.cn,direct

echo "----------------< project >----------------"
module=$(awk 'NR==1 {print}' go.mod)
module=`echo $module | awk '{split($0,a," ");print a[2]}'`
echo " go mod: ${module}"
tag=$(git describe --tags `git rev-list --tags --max-count=1`)
version=${tag:1}
echo "version: ${version}"
last_commit=`git rev-parse HEAD`
author=`git log ${tag} --pretty=format:"%an"|sed -n 1p`
echo " author: ${author}"

function compile() {
    echo "----------------< compile >----------------"
    BIN=./bin
    app=$1
    EXT=$2
    echo "   GOOS: ${GOOS}"
    echo " GOARCH: ${GOARCH}"
    echo "正在编译应用:${app} => ${BIN}/${app}${EXT}..."
    env GOOS=$GOOS GOARCH=$GOARCH go build -ldflags "-s -w -X 'main.MinVersion=${version}'" -o ${BIN}/${app}${EXT} ${module}
    echo "正在编译应用:${app} => ${BIN}/${app}${EXT}...OK"
}

#compile "stock" ""

# 返回当前路径
#cd $p0