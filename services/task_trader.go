package services

import (
	"github.com/quant1x/engine/trader"
	"github.com/quant1x/exchange"
	"github.com/quant1x/x/api"
	"github.com/quant1x/x/logger"
	"github.com/quant1x/x/runtime"
)

// 同步委托订单
func jobSyncTraderOrders() {
	defer runtime.IgnorePanic("")
	// 非交易日直接退出
	if !exchange.DateIsTradingDay() {
		return
	}
	name := trader.GetOrderFilename()
	// 检查文件最后修改时间, 如果文件存在, 且时间在收盘之后, 则跳过同步
	stat, err := api.GetFileStat(name)
	if err == nil && stat != nil {
		modTime := stat.LastWriteTime.Format(exchange.CN_SERVERTIME_FORMAT)
		if modTime >= exchange.CN_CallAuctionPmEnd {
			return
		}
	}
	logger.Info("同步交易订单...")
	defer logger.Info("同步交易订单...OK")
	list, err := trader.QueryOrders()
	if err != nil || len(list) == 0 {
		logger.Info("同步交易订单...今日未操作")
		return
	}
	_ = api.SlicesToCsv(name, list)
}
