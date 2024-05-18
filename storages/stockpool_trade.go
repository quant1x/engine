package storages

import (
	"gitee.com/quant1x/engine/config"
	"gitee.com/quant1x/engine/models"
	"gitee.com/quant1x/engine/trader"
	"gitee.com/quant1x/exchange"
	"gitee.com/quant1x/gox/logger"
)

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
	// 1. 判断是否交易日
	if !exchange.DateIsTradingDay() {
		// 非交易日
		logger.Errorf("%s[%d]: 非交易日, 放弃", model.Name(), model.Code())
		return false
	}
	// 2. 获取策略参数
	strategyParameter := config.GetStrategyParameterByCode(model.Code())
	if strategyParameter == nil || !strategyParameter.BuyEnable() {
		logger.Errorf("%s[%d]: 策略未配置或不买入, 放弃", model.Name(), model.Code())
		return false
	}
	// 3. 判断交易时段
	if !strategyParameter.Session.IsTrading() {
		// 非交易时段
		logger.Errorf("%s[%d]: 非交易时段, 放弃", model.Name(), model.Code())
		return false
	}
	// 4. 矫正交易日期
	tradeDate := exchange.FixTradeDate(date)
	direction := trader.BUY
	// 5. 统计指定交易日的策略已执行买入的标的数量
	numberOfStrategy := CountStrategyOrders(tradeDate, model, direction)
	if numberOfStrategy >= strategyParameter.Total {
		logger.Errorf("%s %s: 计划买入=%d, 已完成=%d. ", tradeDate, model.Name(), strategyParameter.Total, numberOfStrategy)
		return true
	}
	// 6. 策略是否盘中实时订单
	isTickOrder := strategyParameter.Flag == models.OrderFlagTick
	// 7. 策略最大可交易标的配额余量
	//remainQuota := strategyParameter.Total - numberOfStrategy
	length := len(list)
	// 8. 统计有多少标的可以买入, 留给后面统一计算可用资金
	totalTarget := numberOfStrategy
	var traderTargets []*StockPool
	for i := 0; i < length && totalTarget < strategyParameter.Total; i++ {
		v := &(list[i])
		// 8.1 非交易日记录, 忽略
		if v.Date != tradeDate {
			continue
		}
		// 8.2 策略编号不一致或者非可买入, 忽略
		if v.StrategyCode != model.Code() || v.OrderStatus != 1 {
			continue
		}
		// 8.3. 检查是否禁止买入
		if !trader.CheckForBuy(v.Code) {
			// 禁止卖出, 则返回
			logger.Infof("%s[%d]: %s ProhibitForSelling", model.Name(), model.Code(), v.Code)
			continue
		}
		// 8.4 检查买入已完成状态
		ok := CheckOrderState(date, model, v.Code, direction)
		if ok {
			logger.Errorf("%s[%d]: %s 已买入, 放弃", model.Name(), model.Code(), v.Code)
			continue
		}
		totalTarget += 1
		traderTargets = append(traderTargets, v)
	}
	// 9. 计算单只标的最多可用多少资金量
	quotaForTheNumberOfTargets := strategyParameter.Total
	// 9.1 非实时订单的head和tail类型订单, 用不超过配额数的订单数重新核定单一交易标的的可用金额
	if !isTickOrder {
		// 非实时动态新增的标的, 要求必须一次性导入股票池
		quotaForTheNumberOfTargets = totalTarget
	} else {
		// 每日实时订单的总数存在随时可新增的情况, 不能动态调整可用资金总数
		// 只能固定金额买入
	}
	// 9.2 判断可交易标的数量
	if quotaForTheNumberOfTargets < 1 {
		logger.Errorf("%s[%d]: 可交易标的数为0, 放弃", model.Name(), model.Code())
		return false
	}
	// 9.3 调用接口计算单只标的可用资金量
	singleFundsAvailable := trader.CalculateAvailableFundsForSingleTarget(quotaForTheNumberOfTargets, strategyParameter.Weight, strategyParameter.FeeMax, strategyParameter.FeeMin)
	if singleFundsAvailable <= trader.InvalidFee {
		logger.Errorf("%s[%d]: 可用资金为0, 放弃", model.Name(), model.Code())
		return false
	}
	// 10. 遍历订单
	length = len(traderTargets)
	for i := 0; i < length && numberOfStrategy < quotaForTheNumberOfTargets; i++ {
		v := traderTargets[i]
		// 10.1 非交易日记录, 忽略
		if v.Date != tradeDate {
			continue
		}
		// 10.2 策略编号不一致或者非可买入, 忽略
		if v.StrategyCode != model.Code() || v.OrderStatus != 1 {
			continue
		}
		// 确定完整的证券代码
		securityCode := v.Code
		// 10.3 检查买入已完成状态
		ok := CheckOrderState(date, model, securityCode, direction)
		if ok {
			// 已买入的标的, 记录日志, 跳过
			logger.Errorf("%s[%d]: %s 已买入, 放弃", model.Name(), model.Code(), securityCode)
			continue
		}
		// 策略执行交易数+1
		numberOfStrategy += 1
		// 10.4 执行交易指令之前先推送订单已完成状态, 防止意外中断设置已交易信号而重复委托买入
		_ = PushOrderState(date, model, securityCode, direction)
		// 10.5 启用价格笼子的计算方法
		price := trader.CalculatePriceCage(*strategyParameter, direction, v.Buy)
		// 10.6 计算买入费用
		tradeFee := trader.EvaluateFeeForBuy(securityCode, singleFundsAvailable, price)
		if tradeFee.Volume <= trader.InvalidVolume {
			logger.Errorf("%s[%d]: %s 可买数量为0, 放弃", model.Name(), model.Code(), securityCode)
			continue
		}
		// 10.7 执行买入
		orderId, err := trader.PlaceOrder(direction, model, securityCode, trader.FIX_PRICE, tradeFee.Price, tradeFee.Volume)
		v.Status |= StrategyOrderPlaced
		if err != nil || orderId < 0 {
			// 设置标的订单ID为无效
			v.OrderId = trader.InvalidOrderId
			// 设定订单状态为委托失败
			v.Status |= StrategyOrderFailed
			logger.Errorf("%s[%d]: %s 下单失败, error=%+v", model.Name(), model.Code(), securityCode, err)
			continue
		}
		// 10.8 保存订单ID
		v.OrderId = orderId
		v.Status |= StrategyOrderSucceeded
	}
	return numberOfStrategy >= quotaForTheNumberOfTargets
}
