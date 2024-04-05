package tracker

import (
	"gitee.com/quant1x/engine/config"
	"gitee.com/quant1x/engine/datasource"
	"gitee.com/quant1x/engine/factors"
	"gitee.com/quant1x/engine/models"
	"gitee.com/quant1x/engine/permissions"
	"gitee.com/quant1x/engine/utils"
	"gitee.com/quant1x/exchange"
	"gitee.com/quant1x/gox/api"
	"gitee.com/quant1x/gox/logger"
	"gitee.com/quant1x/gox/runtime"
	"sort"
	"time"
)

// Tracker 盘中跟踪
func Tracker(strategyNumbers ...uint64) {
	tickDataAdapter := &datasource.RealTickDataAdapter{}
	updateInRealTime, status := exchange.CanUpdateInRealtime()
	isTrading := updateInRealTime && status == exchange.ExchangeTrading
	if !runtime.Debug() && !isTrading {
		// 非调试且非交易时段返回
		return
	}

	for {
		//stockCodes := radar.ScanSectorForTick(barIndex)

		tickDataAdapter.SyncAllSnapshots()

		for _, strategyNumber := range strategyNumbers {
			model, err := models.CheckoutStrategy(strategyNumber)
			if err != nil || model == nil {
				continue
			}
			err = permissions.CheckPermission(model)
			if err != nil {
				logger.Error(err)
				continue
			}
			strategyParameter := config.GetStrategyParameterByCode(strategyNumber)
			if strategyParameter == nil {
				continue
			}
			if strategyParameter.Session.IsTrading() {
				snapshotTracker(tickDataAdapter, model, strategyParameter)
			} else {
				if runtime.Debug() {
					snapshotTracker(tickDataAdapter, model, strategyParameter)
				} else {
					break
				}
			}
		}
		time.Sleep(time.Second * 1)
	}
}

func snapshotTracker(provider datasource.TickDataProvider, model models.Strategy, tradeRule *config.StrategyParameter) {
	if tradeRule == nil {
		return
	}
	stockCodes := tradeRule.StockList()
	if len(stockCodes) == 0 {
		return
	}
	var stockSnapshots []factors.QuoteSnapshot
	stockCount := len(stockCodes)

	progressManager := utils.NewProgressBarManager("执行["+model.Name()+"全市场扫描]", stockCount)
	progressManager.Start()
	defer progressManager.Wait()

	for start := 0; start < stockCount; start++ {
		progressManager.Update(1)

		code := stockCodes[start]
		securityCode := exchange.CorrectSecurityCode(code)
		if exchange.AssertIndexBySecurityCode(securityCode) {
			continue
		}
		v := provider.GetTickFromMemory(securityCode)
		if v != nil {
			snapshot := provider.QuoteSnapshotFromProtocol(*v)
			stockSnapshots = append(stockSnapshots, snapshot)
		}
	}
	if len(stockSnapshots) == 0 {
		return
	}
	// 过滤不符合条件的个股
	stockSnapshots = api.Filter(stockSnapshots, func(snapshot factors.QuoteSnapshot) bool {
		err := model.Filter(tradeRule.Rules, snapshot)
		return err == nil
	})
	// 结果集排序
	sortedStatus := model.Sort(stockSnapshots)
	if sortedStatus == models.SortDefault || sortedStatus == models.SortNotExecuted {
		// 默认排序或者排序未执行, 使用默认排序
		sort.Slice(stockSnapshots, func(i, j int) bool {
			a := stockSnapshots[i]
			b := stockSnapshots[j]
			if a.OpenTurnZ > b.OpenTurnZ {
				return true
			}
			return a.OpenTurnZ == b.OpenTurnZ && a.OpeningChangeRate > b.OpeningChangeRate
		})
	}
	// 输出表格
	OutputTable(model, stockSnapshots)
}
