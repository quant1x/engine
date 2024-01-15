package command

import "gitee.com/quant1x/gox/runtime"

var (
	Application = "stock"
	MinVersion  = "0.0.1" // 版本号
	barIndex    = 1
)

func init() {
	Application = runtime.ApplicationName()
}

// UpdateApplicationName 更新应用(Application)名称
func UpdateApplicationName(name string) {
	Application = name
}

// UpdateApplicationVersion 更新版本号
func UpdateApplicationVersion(v string) {
	MinVersion = v
}
