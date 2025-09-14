package services

import (
	"github.com/quant1x/gotdx/quotes"
	"github.com/quant1x/x/logger"
)

// 网络配置重置
func jobResetNetwork() {
	logger.Infof("刷新服务器列表...")
	quotes.BestIP()
	logger.Infof("刷新服务器列表...OK")
}
