package trader

import "gitee.com/quant1x/gox/num"

const (
	InvalidFee = float64(-1) // 无效的费用
)

// TradeFee 交易费用
type TradeFee struct {
	SecurityCode  string  // 证券代码
	Price         float64 // 价格
	Volume        int     // 数量
	StampDutyFee  float64 // 印花税, 双向, 默认单向, 费率0.1%
	TransferFee   float64 // 过户费, 双向, 默认是0.06%
	CommissionFee float64 // 券商佣金, 按照成交金额计算, 双向, 0.025%
	StockFee      float64 // 股票市值
	TotalFee      float64 // 总费用
}

// Calc 计算费用
func (f *TradeFee) calculate(direction Direction) float64 {
	volume := float64(f.Volume)
	amount := volume * f.Price
	// 1. 印花税, 按照成交金额计算, 买入没有, 卖出, 0.1%
	if direction == BUY {
		f.StampDutyFee = num.Decimal(amount * traderConfig.StampDutyRateForBuy)
	} else if direction == SELL {
		f.StampDutyFee = num.Decimal(amount * traderConfig.StampDutyRateForSell)
	} else {
		// 返回一个无效的费用常量
		return InvalidFee
	}
	// 2. 过户费, 按照股票数量, 双向, 0.06%
	f.TransferFee = num.Decimal(volume * traderConfig.TransferRate)
	// 3. 券商佣金, 按照成交金额计算, 双向, 0.025%
	f.CommissionFee = num.Decimal(volume * f.Price * traderConfig.CommissionRate)
	if f.CommissionFee < traderConfig.CommissionMin {
		// 不足最低佣金, 要补齐
		f.CommissionFee = traderConfig.CommissionMin
	}
	// 4. 股票市值
	f.StockFee = num.Decimal(volume * f.Price)
	// 5. 计算总费用
	f.TotalFee = f.StampDutyFee + f.TransferFee + f.CommissionFee + f.StockFee
	return f.TotalFee
}

// Calc 通过股价和委托量计算费用
func (f *TradeFee) Calc(price float64, volume int) float64 {
	f.Price = price
	f.Volume = volume
	return f.calculate(BUY)
}

func (f *TradeFee) evaluate(fund float64) {

}

// Estimate 估算可买股票数量
func (f *TradeFee) Estimate(fund, price float64) {

}
