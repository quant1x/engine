package strategies

import (
	"embed"
	"fmt"
	"gitee.com/quant1x/engine/cache"
	"gitee.com/quant1x/gox/api"
	"gitee.com/quant1x/gox/util/homedir"
	"github.com/creasty/defaults"
	"gopkg.in/yaml.v3"
	"os"
)

const (
	TenThousand = 1e4               // 万
	Million     = 100 * TenThousand // 百万
	Billion     = 100 * Million     // 1亿
)

var (
	globalFilterRules = RuleParameter{}
	globalOrderRules  = OrderParameter{}
)

func init() {
	// 初始化配置
	cfg := ReadConfig()
	// 加载规则参数
	_ = api.Copy(&globalFilterRules, &cfg.Rules)
	globalFilterRules.CapitalMin *= Billion
	globalFilterRules.CapitalMax *= Billion
	globalFilterRules.MaxReduceAmount *= TenThousand

	// 加载订单参数
	_ = api.Copy(&globalOrderRules, &cfg.Order)
}

// Quant1XConfig Quant1X基础配置
type Quant1XConfig struct {
	BaseDir string         `yaml:"basedir" default:"~/.quant1x"`
	Rules   RuleParameter  `yaml:"rules"` // 规则参数
	Order   OrderParameter `yaml:"order"` // 订单参数
}

const (
	CACHE_CONFIG_PATH     = "~/.quant1x"
	CACHE_CONFIG_FILENAME = "quant1x.yaml"
)

var (
	// ResourcesPath 资源路径
	ResourcesPath = "resources"
)

//go:embed resources/*
var resources embed.FS

// ReadConfig 读取配置文件
func ReadConfig() (config Quant1XConfig) {
	_ = defaults.Set(&config)
	target := cache.GetConfigFilename()
	if !api.FileExist(target) {
		target = CACHE_CONFIG_PATH + "/" + CACHE_CONFIG_FILENAME
		target, _ = homedir.Expand(target)
		filename := fmt.Sprintf("%s/%s", ResourcesPath, CACHE_CONFIG_FILENAME)
		_ = api.Export(resources, filename, target)
	}
	if api.FileExist(target) {
		dataBytes, err := os.ReadFile(target)
		if err != nil {
			return
		}
		err = yaml.Unmarshal(dataBytes, &config)
		if err != nil {
			return
		}
	}
	return
}
