package services

import (
	"gitee.com/quant1x/data/level1"
	"gitee.com/quant1x/engine/cache"
	"gitee.com/quant1x/engine/factors"
	"gitee.com/quant1x/gox/logger"
	"gitee.com/quant1x/gox/runtime"
)

// 任务 - 交易日数据缓存重置
func jobGlobalReset() {
	defer runtime.IgnorePanic("")
	logger.Info("系统初始化...")
	logger.Info("清理过期的更新状态文件...")
	_ = cleanExpiredStateFiles()
	logger.Info("清理过期的更新状态文件...OK")
	level1.ReOpen()
	logger.Info("重置系统缓存...")
	factors.SwitchDate(cache.DefaultCanReadDate())
	logger.Info("重置系统缓存...OK")

	logger.Info("系统初始化...OK")
}
