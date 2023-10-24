package command

import (
	"os"
	"path/filepath"
)

var (
	Application = "stock"
	MinVersion  = "0.0.1" // 版本号
	barIndex    = 1
)

func init() {
	path, _ := os.Executable()
	_, exec := filepath.Split(path)
	Application = exec
}

// UpdateApplicationName 更新应用(Application)名称
func UpdateApplicationName(name string) {
	Application = name
}

// UpdateApplicationVersion 更新版本号
func UpdateApplicationVersion(v string) {
	MinVersion = v
}
