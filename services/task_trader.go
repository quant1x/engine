package services

import (
	"gitee.com/quant1x/engine/trader"
	"gitee.com/quant1x/gotdx/trading"
	"gitee.com/quant1x/gox/api"
	"gitee.com/quant1x/gox/logger"
)

// 同步委托订单
func jobSyncTraderOrders() {
	// 非交易日直接退出
	if !trading.DateIsTradingDay() {
		return
	}
	logger.Info("同步交易订单...")
	defer logger.Info("同步交易订单...OK")
	list, err := trader.QueryOrders()
	if err != nil || len(list) == 0 {
		logger.Info("同步交易订单...今日未操作")
		return
	}
	name := trader.GetOrderFilename()
	_ = api.SlicesToCsv(name, list)
}
