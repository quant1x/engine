package cache

import (
	"gitee.com/quant1x/gox/logger"
	"gitee.com/quant1x/gox/util/homedir"
	"os"
)

const (
	// 目录权限
	cacheDirMode os.FileMode = 0755
	// 文件权限
	cacheFileMode os.FileMode = 0644
	// 调试开关
	debug = false
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
	// 搜索配置文件
	configPath := searchConfig()
	if len(configPath) > 0 {
		__path, err := homedir.Expand(configPath)
		if err == nil {
			configPath = __path
		} else {
			panic(err)
		}
	} else {
		configPath = cacheRootPath
	}
	// 校验配置文件的路径
	__path, err := homedir.Expand(configPath)
	if err != nil {
		panic(err)
	}
	cacheRootPath = __path
	// 创建根路径
	if err := os.MkdirAll(cacheRootPath, cacheDirMode); err != nil {
		panic(err)
	}
	// 创建日志路径
	cacheLogPath = cacheRootPath + "/logs"
	__logsPath, err := homedir.Expand(cacheLogPath)
	if err != nil {
		panic(err)
	}
	cacheLogPath = __logsPath
	if err := os.MkdirAll(cacheLogPath, cacheDirMode); err != nil {
		panic(err)
	}
	logger.InitLogger(cacheLogPath, logger.INFO)
	// 创建var路径
	cacheVariablePath = cacheRootPath + "/var"
	__varPath, err := homedir.Expand(cacheVariablePath)
	if err != nil {
		panic(err)
	}
	cacheVariablePath = __varPath
	if err := os.MkdirAll(cacheVariablePath, cacheDirMode); err != nil {
		panic(err)
	}
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
