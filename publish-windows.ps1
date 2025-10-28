# engine 编译脚本
# author: wangfeng
# since: 2023-09-12

$commit_id = (git rev-list --tags --max-count=1) | Out-String
$command = "git describe --tags ${commit_id}"
$repo = (Invoke-Expression $command) | Out-String
$version = $repo.Substring(1).Trim()
Write-Output "engine version: ${version}"
$repo = (go list -m gitee.com/quant1x/data/level1) | Out-String
$gotdx_version = ($repo -split " ")[1]
$gotdx_version = $gotdx_version.Substring(1).Trim()
Write-Output "gotdx version: $gotdx_version"

$BIN = "./bin"
$APP = "stock"
$EXT = ".exe"
$repo = "gitee.com/quant1x/engine"
$GOOS = "windows"
$GOARCH = $env:PROCESSOR_ARCHITECTURE
$GOARCH = $GOARCH.Trim().ToLower()
$command = "go env -w GOOS=${GOOS} GOARCH=${GOARCH}; go build -ldflags `"-s  -w -X 'main.MinVersion=$version' -X 'main.tdxVersion=$gotdx_version'`" -o ${BIN}/${APP}${EXT} ${repo}"
Write-Output $command
Invoke-Expression $command
