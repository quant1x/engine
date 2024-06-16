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
		minVersion = fixVersion(latest)
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
			minVersion = fixVersion(arr[1])
		}
	}
	return minVersion
}

// 去掉版本号前的字符v或V
func fixVersion(version string) string {
	latest := strings.TrimSpace(version)
	if latest[0] == 'v' || latest[0] == 'V' {
		latest = latest[1:]
	}
	return latest
}
