package config

import (
	"embed"
	"fmt"
	"gitee.com/quant1x/gox/api"
	"gitee.com/quant1x/gox/logger"
	"gitee.com/quant1x/gox/util/homedir"
	"gitee.com/quant1x/pkg/defaults"
	"gitee.com/quant1x/pkg/yaml"
	"os"
	"strings"
)

const (
	// ResourcesPath 资源路径
	ResourcesPath = "resources"
)

//go:embed resources/*
var resources embed.FS

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
	BaseDir string           `yaml:"basedir"` // 基础路径
	Runtime RuntimeParameter `yaml:"runtime"` // 运行时参数
	Rules   RuleParameter    `yaml:"rules"`   // 规则参数
	Order   OrderParameter   `yaml:"order"`   // 订单参数
	Trader  TraderParameter  `yaml:"trader"`  // 预览交易参数
}

// RuntimeParameter 运行时配置参数
type RuntimeParameter struct {
	Pprof   PprofParameter `yaml:"pprof"`
	Crontab map[string]Job `yaml:"crontab"`
}

// GetConfigFilename 获取配置文件路径
func GetConfigFilename() string {
	return quant1XConfigFilename
}

var (
	// GlobalConfig engine配置信息
	GlobalConfig Quant1XConfig
)

// LoadConfig 加载配置文件
func LoadConfig() (config Quant1XConfig, found bool) {
	_ = defaults.Set(&config)
	for _, v := range listConfigFile {
		filename, err := homedir.Expand(v)
		if err != nil {
			continue
		}
		if api.FileExist(filename) {
			dataBytes, err := os.ReadFile(filename)
			if err != nil {
				logger.Error(err)
				continue
			}
			err = yaml.Unmarshal(dataBytes, &config)
			if err != nil {
				logger.Error(err)
				continue
			}
			config.BaseDir = strings.TrimSpace(config.BaseDir)
			if len(config.BaseDir) > 0 {
				quant1XConfigFilename = filename
			}
			fixTradingSession(&config)
			found = true
			break
		}
	}
	return
}

// ReadConfig 读取配置文件
func ReadConfig(rootPath string) (config Quant1XConfig) {
	_ = defaults.Set(&config)
	target := GetConfigFilename()
	if !api.FileExist(target) {
		target = rootPath + "/" + configFilename
		target, _ = homedir.Expand(target)
		filename := fmt.Sprintf("%s/%s", ResourcesPath, configFilename)
		_ = api.Export(resources, filename, target)
	}
	if api.FileExist(target) {
		dataBytes, err := os.ReadFile(target)
		if err != nil {
			logger.Error(err)
			return
		}
		err = yaml.Unmarshal(dataBytes, &config)
		if err != nil {
			logger.Error(err)
			return
		}
		fixTradingSession(&config)
	}
	return
}

func fixTradingSession(config *Quant1XConfig) {
	// 增加配置的二次加工
	// 1. 重置交易时段字段
	config.Trader.Head.Session = ParseTradingSession(config.Trader.Head.Time)
	config.Trader.Tail.Session = ParseTradingSession(config.Trader.Tail.Time)
	config.Trader.Tick.Session = ParseTradingSession(config.Trader.Tick.Time)
	config.Trader.Sell.Session = ParseTradingSession(config.Trader.Sell.Time)
	config.Trader.CancelSession = ParseTradingSession(config.Trader.ReservedOfCancel)
}
