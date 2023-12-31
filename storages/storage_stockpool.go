package storages

import (
	"fmt"
	"gitee.com/quant1x/engine/cache"
	"gitee.com/quant1x/engine/config"
	"gitee.com/quant1x/engine/models"
	"gitee.com/quant1x/engine/trader"
	"gitee.com/quant1x/gotdx/trading"
	"gitee.com/quant1x/gox/api"
	"gitee.com/quant1x/gox/logger"
	"sync"
	"time"
)

const (
	filenameStockPool = "stock_pool.csv"
)

var (
	poolMutex sync.Mutex
)

func GetStockPool() (list []StockPool) {
	filename := fmt.Sprintf("%s/%s", getQmtCachePath(), filenameStockPool)
	err := api.CsvToSlices(filename, &list)
	_ = err
	return
}

func SaveStockPool(list []StockPool) {
	filename := fmt.Sprintf("%s/%s", getQmtCachePath(), filenameStockPool)
	err := api.SlicesToCsv(filename, list)
	_ = err
	return
}

func stockPoolMerge(model models.Strategy, date string, orders []models.Statistics, topN int) {
	poolMutex.Lock()
	defer poolMutex.Unlock()
	localStockPool := GetStockPool()
	//targets := []StockPool{}
	cacheStatistics := map[string]*StockPool{}
	tradeDate := trading.FixTradeDate(date)
	for i, v := range orders {
		sp := StockPool{
			//Status         StrategyStatus `name:"策略状态" dataframe:"status"`
			Status: StrategyHit,
			//Date           string         `name:"信号日期" dataframe:"date"`
			Date: v.Date,
			//Code           string         `name:"证券代码" dataframe:"code"`
			Code: v.Code,
			//Name           string         `name:"证券名称" dataframe:"name"`
			Name: v.Name,
			//TurnZ          float64        `name:"开盘换手Z" dataframe:"turn_z"`
			TurnZ: v.TurnZ,
			//Rate           float64        `name:"涨跌幅%" dataframe:"rate"`
			Rate: v.UpRate,
			//Buy            float64        `name:"委托价格" dataframe:"buy"`
			Buy: v.Price,
			//Sell           float64        `name:"目标价格" dataframe:"sell"`
			//Sell: v.Price,
			//StrategyCode   int            `name:"策略编码" dataframe:"strategy_code"`
			StrategyCode: model.Code(),
			//StrategyName   string         `name:"策略名称" dataframe:"strategy_name"`
			StrategyName: model.Name(),
			//Rules          uint64         `name:"规则" dataframe:"rules"`
			//BlockType      string         `name:"板块类型" dataframe:"block_type"`
			//BlockType: v.BlockType,
			//BlockCode      string         `name:"板块代码" dataframe:"block_code"`
			//BlockName      string         `name:"板块名称" dataframe:"block_name"`
			//BlockRate      float64        `name:"板块涨幅%" dataframe:"block_rate"`
			//BlockTop       int            `name:"板块排名" dataframe:"block_top"`
			//BlockRank      int            `name:"个股排名" dataframe:"block_rank"`
			//BlockZhangTing string         `name:"板块涨停数" dataframe:"block_zhangting"`
			//BlockDescribe  string         `name:"涨/跌/平" dataframe:"block_describe"`
			//BlockTopCode   string         `name:"领涨股代码" dataframe:"block_top_code"`
			//BlockTopName   string         `name:"领涨股名称" dataframe:"block_top_name"`
			//BlockTopName: v.BlockName,
			//BlockTopRate   float64        `name:"领涨股涨幅%" dataframe:"block_top_rate"`
			//BlockTopRate: v.BlockRate,
			//Tendency       string         `name:"短线趋势" dataframe:"tendency"`
			//Tendency: v.Tendency,
			OrderStatus: 0, // 默认订单状态是0
			Active:      v.Active,
			Speed:       v.Speed,
			//CreateTime     string         `name:"创建时间" dataframe:"create_time"`
			CreateTime: v.UpdateTime,
			//UpdateTime     string         `name:"更新时间" dataframe:"update_time"`
			UpdateTime: v.UpdateTime,
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
	newList := []StockPool{}
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
		checkOrderForBuy(localStockPool, model, date)
		SaveStockPool(localStockPool)
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
	tradeDate := trading.FixTradeDate(date)
	strategyParameter := config.GetStrategyParameterByCode(model.Code())
	if strategyParameter != nil && strategyParameter.BuyEnable() {
		direction := trader.BUY
		numberOfStrategy := CountStrategyOrders(tradeDate, model, direction)
		if numberOfStrategy >= strategyParameter.Total {
			logger.Warnf("%s %s: 计划买入=%d, 已完成=%d. ", tradeDate, model.Name(), strategyParameter.Total, numberOfStrategy)
			return true
		}
		length := len(list)
		for i := 0; i < length && numberOfStrategy < strategyParameter.Total; i++ {
			v := &(list[i])
			if v.Date != tradeDate {
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
				price := v.Buy + 0.05
				// 2. 检查买入已完成状态
				ok := CheckOrderState(date, model, securityCode, direction)
				if ok {
					logger.Errorf("%s[%d]: %s 已买入, 放弃", model.Name(), model.Code(), securityCode)
					continue
				}
				// 3. 首先推送订单已完成状态
				_ = PushOrderState(date, model, securityCode, direction)
				if !trading.DateIsTradingDay() {
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
				orderId, err := trader.PlaceOrder(direction, model, securityCode, tradeFee.Price, tradeFee.Volume)
				if err != nil {
					logger.Errorf("%s[%d]: %s 下单失败, error=%+v", model.Name(), model.Code(), securityCode, err)
					continue
				}
				// 7. 保存订单ID
				v.OrderId = orderId
			}
		}
		return numberOfStrategy >= strategyParameter.Total
	}
	return false
}
