package cache

import (
	"os"

	"gitee.com/quant1x/engine/config"
	"gitee.com/quant1x/gox/logger"
	"gitee.com/quant1x/gox/util/homedir"
)

const (
	// 目录权限
	cacheDirMode os.FileMode = 0755
	// 文件权限
	cacheFileMode os.FileMode = 0644
	// 文件替换模式, 会用到os.TRUNC
	cacheReplace = os.O_CREATE | os.O_RDWR | os.O_TRUNC
	// 更新
	cacheUpdate = os.O_CREATE | os.O_WRONLY
)

var (
	// 根路径
	cacheRootPath = "~/.quant1x"
	// cacheLogPath 日志路径
	cacheLogPath = cacheRootPath + "/logs"
	// var 路径
	cacheVariablePath = cacheRootPath + "/var"
)

func init() {
	initCache()
}

func initCache() {
	// 加载配置文件
	tmpConfig, found := config.LoadConfig()
	// 搜索配置文件
	baseDir := tmpConfig.BaseDir
	if len(baseDir) > 0 {
		__path, err := homedir.Expand(baseDir)
		if err == nil {
			baseDir = __path
		} else {
			logger.Fatalf("%+v", err)
		}
	} else {
		baseDir = cacheRootPath
	}
	// 校验配置文件的路径
	__path, err := homedir.Expand(baseDir)
	if err != nil {
		logger.Fatalf("%+v", err)
	}
	cacheRootPath = __path
	// 创建根路径
	if err := os.MkdirAll(cacheRootPath, cacheDirMode); err != nil {
		logger.Fatalf("%+v", err)
	}
	// 创建日志路径
	cacheLogPath = cacheRootPath + "/logs"
	__logsPath, err := homedir.Expand(cacheLogPath)
	if err != nil {
		logger.Fatalf("%+v", err)
	}
	cacheLogPath = __logsPath
	if err := os.MkdirAll(cacheLogPath, cacheDirMode); err != nil {
		logger.Fatalf("%+v", err)
	}
	logger.InitLogger(cacheLogPath, logger.INFO)
	// 创建var路径
	cacheVariablePath = cacheRootPath + "/var"
	__varPath, err := homedir.Expand(cacheVariablePath)
	if err != nil {
		logger.Fatalf("%+v", err)
	}
	cacheVariablePath = __varPath
	if err := os.MkdirAll(cacheVariablePath, cacheDirMode); err != nil {
		logger.Fatalf("%+v", err)
	}
	// 检查配置文件并加载配置
	if !found {
		config.GlobalConfig = config.ReadConfig(GetRootPath())
	} else {
		config.GlobalConfig = tmpConfig
	}

	// 启动性能分析
	config.StartPprof()
	initMiniQmt()
}

// Reset 重置日志记录器
func Reset() {
	initCache()
}

// GetRootPath 获取缓存根路径
func GetRootPath() string {
	return cacheRootPath
}

// GetLoggerPath 获取日志路径
func GetLoggerPath() string {
	return cacheLogPath
}

// GetVariablePath 获取VAR路径
func GetVariablePath() string {
	return cacheVariablePath
}
