package trader

import (
	"gitee.com/quant1x/engine/config"
	"gitee.com/quant1x/gox/coroutine"
	"gitee.com/quant1x/gox/logger"
	"gitee.com/quant1x/gox/num"
)

var (
	onceAccount            coroutine.PeriodicOnce
	accountTheoreticalFund = float64(0.00)
	accountRemainingCash   = float64(0.00)
)

func lazyInitFundPool() {
	accountTheoreticalFund, accountRemainingCash = calculateTheoreticalFund()
}

// 计算理论上可用的资金
func calculateTheoreticalFund() (theoretical, cash float64) {
	theoretical = InvalidFee
	cash = InvalidFee
	// 1. 查询 总资产和可用
	// acc_total, acc_cash = self.account_available()
	acc, err := QueryAccount()
	if err != nil {
		return
	}
	// 2. 查询持仓可卖的股票 TODO: 如果确定了可卖出的市值, 怎么保证当日必须卖出?
	// positions = self.query_positions()
	positions, err := QueryHolding()
	if err != nil {
		return
	}
	can_use_amount := 0.00
	// 3. 设置一个可卖股票的市值波动范围, 这里暂定10%
	vix := 0.10
	acc_value := 0.00 // 总市值
	for _, position := range positions {
		acc_value += position.MarketValue
		if position.CanUseVolume < 1 {
			continue
		}
		// 计算现价
		market_price := position.MarketValue / float64(position.Volume)
		// 累计可卖的市值: 可卖数量 * 市价
		can_use_value := market_price * float64(position.CanUseVolume)
		can_use_amount += can_use_value * (1 - vix)
	}
	acc_value = num.Decimal(acc_value)
	can_use_amount = num.Decimal(can_use_amount)
	// 4. 确定可用资金总量: 账户可以资金 + 当日可卖出的总市值 - 预留现金
	can_use_cash := acc.Cash + can_use_amount - traderConfig.KeepCash
	// 5. 计算预留仓位, 给下一个交易日留position_ratio仓位
	reserve_cash := num.Decimal(acc.TotalAsset * traderConfig.PositionRatio)
	// 6. 计算当日可用仓位: 可用资金总量 - 预留资金总量
	available := can_use_cash - reserve_cash
	logger.Warnf("账户资金: 可用=%.02f, 市值=%.02f, 预留=%.02f, 可买=%.02f, 可卖=%.02f", acc.Cash, acc_value, reserve_cash, available, can_use_amount)
	// 7. 如果当日可用金额大于资金账户的可用金额, 输出风险提示
	if available > acc.Cash {
		logger.Warnf("!!! 持仓占比[{}%], 已超过可总仓位的[{}%], 必须在收盘前择机降低仓位, 以免影响下一个交易日的买入操作 !!!", num.Decimal(100*(acc_value/acc.TotalAsset)),
			num.Decimal(100*(1-traderConfig.PositionRatio)))
	}
	// 8. 重新修订可用金额
	available = (acc.TotalAsset - traderConfig.KeepCash) * traderConfig.PositionRatio
	if available > acc.Cash {
		available = acc.Cash
	}
	theoretical = available
	cash = acc.Cash
	return theoretical, cash
}

// CalculateAvailableFund 计算一只股票的可动用资金量
//
// 参数:
//
//	totalBalance: 总账户余额
//	marketValue: 股票当前市值
//	averageCost: 平均买入成本（用于计算盈亏）
//	commission: 交易佣金率
//
// 返回值:
//
//	availableFund: 可动用资金量
func CalculateAvailableFund(strategyParameter *config.StrategyParameter) float64 {
	onceAccount.Do(lazyInitFundPool)
	if strategyParameter.Total < 1 {
		return InvalidFee
	}
	// 1. 检查可用资金
	if accountTheoreticalFund <= InvalidFee {
		return InvalidFee
	}
	// 2. 计算策略的可用资金, 总可用资金*策略权重
	strategy_fund := accountTheoreticalFund * strategyParameter.Weight
	single_funds_available := num.Decimal(strategy_fund / float64(strategyParameter.Total))
	// 3. 检查策略的可用资金范围
	if single_funds_available > strategyParameter.FeeMax {
		single_funds_available = strategyParameter.FeeMax
	} else if single_funds_available < strategyParameter.FeeMin {
		return InvalidFee
	}
	// 4. 检查可用资金的最大值和最小值
	if single_funds_available > traderConfig.BuyAmountMax {
		single_funds_available = traderConfig.BuyAmountMax
	} else if single_funds_available < traderConfig.BuyAmountMin {
		return InvalidFee
	}
	return single_funds_available
}
