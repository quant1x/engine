package storages

import (
	"gitee.com/quant1x/engine/cache"
	"gitee.com/quant1x/engine/config"
	"gitee.com/quant1x/engine/models"
	"gitee.com/quant1x/engine/trader"
	"gitee.com/quant1x/exchange"
	"gitee.com/quant1x/gox/api"
	"gitee.com/quant1x/gox/logger"
	"path/filepath"
	"sync"
	"time"
)

const (
	filenameStockPool = "stock_pool.csv"
)

var (
	poolMutex sync.Mutex
)

// 股票池文件
func getStockPoolFilename() string {
	filename := filepath.Join(cache.GetQmtCachePath(), filenameStockPool)
	return filename
}

func getStockPoolFromCache() (list []StockPool) {
	filename := getStockPoolFilename()
	err := api.CsvToSlices(filename, &list)
	_ = err
	return
}

func saveStockPoolToCache(list []StockPool) {
	filename := getStockPoolFilename()
	err := api.SlicesToCsv(filename, list)
	_ = err
	return
}

func stockPoolMerge(model models.Strategy, date string, orders []models.Statistics, topN int) {
	poolMutex.Lock()
	defer poolMutex.Unlock()
	localStockPool := getStockPoolFromCache()
	//targets := []StockPool{}
	cacheStatistics := map[string]*StockPool{}
	tradeDate := exchange.FixTradeDate(date)
	for i, v := range orders {
		sp := StockPool{
			Status:       StrategyHit,
			Date:         v.Date,
			Code:         v.Code,
			Name:         v.Name,
			Buy:          v.Price,
			StrategyCode: model.Code(),
			StrategyName: model.Name(),
			OrderStatus:  0, // 股票池订单状态默认是0
			Active:       v.Active,
			Speed:        v.Speed,
			CreateTime:   v.UpdateTime,
			UpdateTime:   v.UpdateTime,
		}
		if i < topN {
			//  如果是前排个股标志可以买入
			sp.OrderStatus = 1
		}
		//targets = append(targets, sp)
		cacheStatistics[sp.Key()] = &sp
	}
	count := len(localStockPool)
	now := time.Now()
	updateTime := now.Format(cache.TimeStampMilli)
	for i := 0; i < count; i++ {
		local := &(localStockPool[i])
		// 1. 非当日的跳过
		if local.Date != tradeDate {
			continue
		}
		v, found := cacheStatistics[local.Key()]
		if found {
			// 相同日期, 策略和证券代码, 视为重复
			// 找到了, 标记为已存在
			v.Status = StrategyAlreadyExists
			//local.OrderStatus = v.OrderStatus
			continue
		}
		// 没找到, 做召回处理
		local.Status.Set(StrategyCancel, true)
		local.UpdateTime = updateTime
	}
	var newList []StockPool
	for _, v := range cacheStatistics {
		if v.Status == StrategyAlreadyExists {
			continue
		}
		v.UpdateTime = updateTime
		logger.Infof("%s[%d]: buy queue append %s", model.Name(), model.Code(), v.Code)
		newList = append(newList, *v)
	}
	if len(newList) > 0 {
		localStockPool = append(localStockPool, newList...)
		logger.Infof("检查是否需要委托下单...")
		checkOrderForBuy(localStockPool, model, date)
		logger.Infof("检查是否需要委托下单...OK")
		saveStockPoolToCache(localStockPool)
	}
}

// 策略订单是否已完成
func strategyOrderIsFinished(model models.Strategy) bool {
	strategyId := model.Code()
	strategyName := models.QmtStrategyName(model)
	tradeRule := config.GetStrategyParameterByCode(strategyId)
	if tradeRule == nil || !tradeRule.BuyEnable() {
		return true
	}
	orders, err := trader.QueryOrders()
	if err != nil {
		return true
	}
	total := 0
	for _, v := range orders {
		if v.StrategyName == strategyName && v.OrderType == trader.STOCK_BUY {
			total++
		}
	}
	return total >= tradeRule.Total
}

// 检查买入订单, 条件满足则买入
func checkOrderForBuy(list []StockPool, model models.Strategy, date string) bool {
	tradeDate := exchange.FixTradeDate(date)
	strategyParameter := config.GetStrategyParameterByCode(model.Code())
	if strategyParameter != nil && strategyParameter.BuyEnable() {
		direction := trader.BUY
		numberOfStrategy := CountStrategyOrders(tradeDate, model, direction)
		if numberOfStrategy >= strategyParameter.Total {
			logger.Errorf("%s %s: 计划买入=%d, 已完成=%d. ", tradeDate, model.Name(), strategyParameter.Total, numberOfStrategy)
			return true
		}
		length := len(list)
		for i := 0; i < length && numberOfStrategy < strategyParameter.Total; i++ {
			v := &(list[i])
			if v.Date != tradeDate {
				logger.Errorf("订单日期不匹配: order[%s], trading[%s]", v.Date, tradeDate)
				continue
			}
			if v.StrategyCode == model.Code() && v.OrderStatus == 1 {
				numberOfStrategy += 1
				securityCode := v.Code
				// 1. 检查是否禁止买入
				if !trader.CheckForBuy(securityCode) {
					// 禁止卖出, 则返回
					logger.Infof("%s[%d]: %s ProhibitForSelling", model.Name(), model.Code(), securityCode)
					continue
				}
				// 暂时不用价格笼子, 只向上浮动0.05
				price := trader.CalculatePriceCage(*strategyParameter, direction, v.Buy)
				//price = v.Buy + 0.05
				// 2. 检查买入已完成状态
				ok := CheckOrderState(date, model, securityCode, direction)
				if ok {
					logger.Errorf("%s[%d]: %s 已买入, 放弃", model.Name(), model.Code(), securityCode)
					continue
				}
				// 3. 首先推送订单已完成状态
				_ = PushOrderState(date, model, securityCode, direction)
				if !exchange.DateIsTradingDay() {
					// 非交易日
					logger.Errorf("%s[%d]: %s 非交易日, 放弃", model.Name(), model.Code(), securityCode)
					continue
				}
				if !strategyParameter.Session.IsTrading() {
					// 非交易时段
					logger.Errorf("%s[%d]: %s 非交易时段, 放弃", model.Name(), model.Code(), securityCode)
					continue
				}
				// 4. 执行买入
				fund := trader.CalculateAvailableFund(strategyParameter)
				if fund <= trader.InvalidFee {
					logger.Errorf("%s[%d]: %s 可用资金为0, 放弃", model.Name(), model.Code(), securityCode)
					continue
				}
				// 5. 计算买入费用
				tradeFee := trader.EvaluateFeeForBuy(securityCode, fund, price)
				if tradeFee.Volume <= trader.InvalidVolume {
					logger.Errorf("%s[%d]: %s 可买数量为0, 放弃", model.Name(), model.Code(), securityCode)
					continue
				}
				// 6. 执行买入
				orderId, err := trader.PlaceOrder(direction, model, securityCode, trader.FIX_PRICE, tradeFee.Price, tradeFee.Volume)
				v.Status |= StrategyOrderPlaced
				if err != nil || orderId < 0 {
					v.OrderId = -1
					v.Status |= StrategyOrderFailed
					logger.Errorf("%s[%d]: %s 下单失败, error=%+v", model.Name(), model.Code(), securityCode, err)
					continue
				}
				// 7. 保存订单ID
				v.OrderId = orderId
				v.Status |= StrategyOrderSucceeded
			}
		}
		return numberOfStrategy >= strategyParameter.Total
	}
	return false
}
