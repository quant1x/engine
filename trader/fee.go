package trader

import (
	"gitee.com/quant1x/gox/logger"
	"gitee.com/quant1x/gox/num"
	"math"
)

const (
	InvalidFee    = float64(-1) // 无效的费用
	UnknownVolume = int(1)      // 未知费用
)

// 计算买入总费用
func calculate_fee(direction Direction, price float64, volume int, align bool) (float64, float64, float64, float64, float64) {
	if volume < 1 {
		return InvalidFee, 0, 0, 0, 0
	}
	vol := float64(volume)
	amount := vol * price
	// 1. 印花税, 按照成交金额计算, 买入没有, 卖出, 0.1%
	_stamp_duty_fee := amount
	if direction == BUY {
		_stamp_duty_fee *= traderConfig.StampDutyRateForBuy
	} else if direction == SELL {
		_stamp_duty_fee *= traderConfig.StampDutyRateForSell
	} else {
		return InvalidFee, 0, 0, 0, 0
	}
	if align {
		_stamp_duty_fee = num.Decimal(_stamp_duty_fee)
	}
	// 2. 过户费, 按照股票数量, 双向, 0.06%
	_transfer_fee := vol * traderConfig.TransferRate
	if align {
		_transfer_fee = num.Decimal(_transfer_fee)
	}
	// 3. 券商佣金, 按照成交金额计算, 双向, 0.025%
	_commission_fee := amount * traderConfig.CommissionRate
	if align {
		_commission_fee = num.Decimal(_commission_fee)
	}
	if align && _commission_fee < traderConfig.CommissionMin {
		_commission_fee = traderConfig.CommissionMin
	}
	// 4. 股票市值
	_stock_fee := amount
	if align {
		_stock_fee = num.Decimal(_stock_fee)
	}
	// 5. 计算总费用
	_fee := _stamp_duty_fee + _transfer_fee + _commission_fee + _stock_fee
	return _fee, _stamp_duty_fee, _transfer_fee, _commission_fee, _stock_fee
}

// EvaluateFeeForBuy 评估费用
func EvaluateFeeForBuy(securityCode string, fund, price float64) *TradeFee {
	f := TradeFee{
		SecurityCode: securityCode,
		Price:        price,
	}
	f.Volume = f.Estimate(fund, price)
	return &f
}

// TradeFee 交易费用
type TradeFee struct {
	Direction     Direction // 交易方向
	SecurityCode  string    // 证券代码
	Price         float64   // 价格
	Volume        int       // 数量
	StampDutyFee  float64   // 印花税, 双向, 默认单向, 费率0.1%
	TransferFee   float64   // 过户费, 双向, 默认是0.06%
	CommissionFee float64   // 券商佣金, 按照成交金额计算, 双向, 0.025%
	StockFee      float64   // 股票市值
	TotalFee      float64   // 总费用
}

//// Calc 计算费用
//func (f *TradeFee) calculate(direction Direction) float64 {
//	volume := float64(f.Volume)
//	amount := volume * f.Price
//	// 1. 印花税, 按照成交金额计算, 买入没有, 卖出, 0.1%
//	if direction == BUY {
//		f.StampDutyFee = num.Decimal(amount * traderConfig.StampDutyRateForBuy)
//	} else if direction == SELL {
//		f.StampDutyFee = num.Decimal(amount * traderConfig.StampDutyRateForSell)
//	} else {
//		// 返回一个无效的费用常量
//		return InvalidFee
//	}
//	// 2. 过户费, 按照股票数量, 双向, 0.06%
//	f.TransferFee = num.Decimal(volume * traderConfig.TransferRate)
//	// 3. 券商佣金, 按照成交金额计算, 双向, 0.025%
//	f.CommissionFee = num.Decimal(volume * f.Price * traderConfig.CommissionRate)
//	if f.CommissionFee < traderConfig.CommissionMin {
//		// 不足最低佣金, 要补齐
//		f.CommissionFee = traderConfig.CommissionMin
//	}
//	// 4. 股票市值
//	f.StockFee = num.Decimal(volume * f.Price)
//	// 5. 计算总费用
//	f.TotalFee = f.StampDutyFee + f.TransferFee + f.CommissionFee + f.StockFee
//	return f.TotalFee
//}
//
//// Calc 通过股价和委托量计算费用
//func (f *TradeFee) Calc(price float64, volume int) float64 {
//	f.Price = price
//	f.Volume = volume
//	return f.calculate(BUY)
//}
//
//func (f *TradeFee) evaluate(fund float64) {
//
//}

// Estimate 估算可买股票数量
func (f *TradeFee) Estimate(fund, price float64) int {
	direction := BUY
	f.Direction = direction
	f.Price = price
	// 1. 计算每股费用
	_fee, _, _, _, _ := calculate_fee(direction, price, UnknownVolume, false)
	// 2. 计算股数
	_vol := fund / _fee
	// 3. 换算成手数
	_vol = math.Floor(_vol / 100)
	// 4. 转成整数
	f.Volume = int(_vol) * 100
	// 5. 重新计算
	f.TotalFee, f.StampDutyFee, f.TransferFee, f.CommissionFee, f.StockFee = calculate_fee(f.Direction, f.Price, f.Volume, true)
	if _fee > fund {
		// 如果费用超了, 则减去1手(100股)
		f.Volume -= 100
		f.TotalFee, f.StampDutyFee, f.TransferFee, f.CommissionFee, f.StockFee = calculate_fee(f.Direction, f.Price, f.Volume, true)
	}
	logger.Infof("trader: code=%s: 综合费用=%.02f, 委托价格=%.02f, 数量=%d, 其中印花说=%.02f, 过户费=%.02f, 佣金=%.02f, 股票=%.02f", f.SecurityCode, f.TotalFee, price,
		f.Volume, f.StampDutyFee, f.TransferFee, f.CommissionFee, f.StockFee)
	return f.Volume
}
