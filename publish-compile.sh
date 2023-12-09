#!/bin/sh
# 获取当前路径, 用于返回
#p0=`pwd`
# 获取脚本所在路径, 防止后续操作在非项目路径
#p1=$(cd $(dirname $0);pwd)

#golang
GOVERSION=$(go env GOVERSION)
GOOS=$(go env GOOS)
GOARCH=$(go env GOARCH)
echo "----------------< go env >----------------"
echo "   GOOS: ${GOOS}"
echo " GOARCH: ${GOARCH}"
echo "version: ${GOVERSION:2}"

export GO111MODULE=auto
export GOPRIVATE=gitee.com
export GOPROXY=https://goproxy.cn,direct

echo "----------------< project env >----------------"
version=$(git describe --tags `git rev-list --tags --max-count=1`)
version=${version:1}
echo "version: ${version}"

function compile() {
    BIN=./bin
    app=$1
    EXT=$2
    echo "正在编译应用:${app} => ${BIN}/${app}${EXT}..."
    env GOOS=$GOOS GOARCH=$GOARCH go build -ldflags "-s -w -X 'main.MinVersion=${version}'" -o ${BIN}/${app}${EXT} gitee.com/quant1x/engine
    echo "正在编译应用:${app} => ${BIN}/${app}${EXT}...OK"
}

#compile "stock" ""

# 返回当前路径
#cd $p0