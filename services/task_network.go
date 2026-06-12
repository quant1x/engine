package services

import (
	"github.com/quant1x/data/level1/quotes"
	"github.com/quant1x/gox/logger"
)

// 网络配置重置
func jobResetNetwork() {
	logger.Infof("刷新服务器列表...")
	quotes.BestIP()
	logger.Infof("刷新服务器列表...OK")
}
