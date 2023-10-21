package cache

import (
	"gitee.com/quant1x/gox/api"
	"gitee.com/quant1x/gox/util/homedir"
	"github.com/creasty/defaults"
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
	BaseDir string           `yaml:"basedir"`
	Runtime RuntimeParameter `yaml:"runtime"`
}

type RuntimeParameter struct {
	Pprof PprofParameter `yaml:"pprof"`
}

type PprofParameter struct {
	Port int `yaml:"port" default:"6060"` // pprof web端口
}

// GetConfigFilename 获取配置文件路径
func GetConfigFilename() string {
	return quant1XConfigFilename
}

// 加载配置文件
func loadConfig() (config Quant1XConfig) {
	_ = defaults.Set(&config)
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
			err = yaml.Unmarshal(dataBytes, &config)
			if err != nil {
				continue
			}
			config.BaseDir = strings.TrimSpace(config.BaseDir)
			if len(config.BaseDir) > 0 {
				quant1XConfigFilename = filename
			}
			break
		}
	}
	return
}
