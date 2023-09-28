package cache

import (
	"gitee.com/quant1x/gox/api"
	"gitee.com/quant1x/gox/util/homedir"
	"gopkg.in/yaml.v3"
	"os"
	"strings"
)

// 配置常量
const (
	// 配置文件名
	configFilename = "quant1x.yaml"
)

var (
	quant1XConfigFilename string = "" // 配置文件完整路径
	listConfigFile               = []string{
		"~/." + configFilename,
		"~/.quant1x/" + configFilename,
		"~/runtime/etc/" + configFilename,
	}
)

// Quant1XConfig Quant1X基础配置
type Quant1XConfig struct {
	BaseDir string `json:"basedir"`
}

// 搜索配置文件
func searchConfig() (baseDir string) {
	for _, v := range listConfigFile {
		filename, err := homedir.Expand(v)
		if err != nil {
			continue
		}
		if api.FileExist(filename) {
			dataBytes, err := os.ReadFile(filename)
			if err != nil {
				continue
			}
			config := Quant1XConfig{}
			err = yaml.Unmarshal(dataBytes, &config)
			if err != nil {
				continue
			}
			config.BaseDir = strings.TrimSpace(config.BaseDir)
			if len(config.BaseDir) > 0 {
				baseDir = config.BaseDir
				quant1XConfigFilename = filename
			}
			break
		}
	}
	return
}

// GetConfigFilename 获取配置文件路径
func GetConfigFilename() string {
	return quant1XConfigFilename
}
