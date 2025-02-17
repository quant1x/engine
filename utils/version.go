package utils

import (
	"strings"
)

const (
	InvalidVersion = "0.0.0"
)

// CurrentVersion 开发中获取版本号
func CurrentVersion() string {
	minVersion := InvalidVersion
	latest, err := git.Exec("describe", "--tags", "--abbrev=0")
	if err == nil {
		minVersion = NormalizeVersion(latest)
	}
	return minVersion
}

// RequireVersion 依赖模块版本号
//
//	通过 go list -m 命令获取
func RequireVersion(module string) string {
	minVersion := InvalidVersion
	mod, err := shell("go", "list", "-m", module)
	if err == nil {
		arr := strings.Split(mod, " ")
		if len(arr) >= 2 {
			minVersion = NormalizeVersion(arr[1])
		}
	}
	return minVersion
}

// NormalizeVersion 表示将版本号格式化为标准形式
//
//	去掉版本号前的字符v或V
func NormalizeVersion(version string) string {
	latest := strings.TrimSpace(version)
	for len(latest) > 0 && (latest[0] == 'v' || latest[0] == 'V') {
		latest = latest[1:]
	}
	return latest
}
