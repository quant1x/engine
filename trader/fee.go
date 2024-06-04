package trader

import (
	"fmt"
	"gitee.com/quant1x/engine/config"
	"gitee.com/quant1x/num"
	"math"
)

const (
	InvalidFee            = float64(0) // 无效的费用
	InvalidVolume         = int(0)     // 无效的股数
	UnknownVolume         = int(1)     // 未知的股数
	InvalidOrderId        = int(-1)    // 无效的订单ID
	MinimumNumberOfOrders = int(2)     // 最小订单数为2, 预防一个策略资金都打到一个标的上面
)

// 价格笼子
//
//	价格笼子是买卖股票申报价格限制的一种制度
//	对于主板, 本次新增2%有效申报价格范围要求, 同时增加10个申报价格最小变动单位的安排
//	A股最小交易变动单位是0.01元，10个也就是0.1元
//	买入价取两者高值，卖出价取两者低值.
const (
	validDeclarationPriceRange  = float64(0.02) // 价格限制比例
	minimumPriceFluctuationUnit = float64(0.10) // 价格浮动最大值
)

// CalculatePriceCage 计算价格笼子
func CalculatePriceCage(strategyParameter config.StrategyParameter, direction Direction, price float64) float64 {
	priceLimit := 0.000
	var priceCage, minimumPriceFluctuation float64
	if direction == BUY {
		priceCage = price * (1 + strategyParameter.PriceCageRatio)
		minimumPriceFluctuation = price + strategyParameter.MinimumPriceFluctuationUnit
		priceLimit = max(priceCage, minimumPriceFluctuation)
	} else {
		priceCage = price * (1 - strategyParameter.PriceCageRatio)
		minimumPriceFluctuation = price - strategyParameter.MinimumPriceFluctuationUnit
		priceLimit = min(priceCage, minimumPriceFluctuation)
	}
	return priceLimit
}

// calculate_price_limit_for_buy 计算合适的买入价格
//
//	价格笼子, +2%和+0.10哪个大
//	目前使用, 当前价格+0.05
func calculate_price_limit_for_buy(last_price, price_cage_ratio, minimum_price_fluctuation_unit float64) float64 {
	// 价格笼子, +2%和+0.10哪个大
	priceLimit := max(last_price*(1+price_cage_ratio), last_price+minimum_price_fluctuation_unit)
	// 当前价格+0.05
	priceLimit = last_price + 0.05
	// 最后修订价格
	priceLimit = num.Decimal(priceLimit)
	return priceLimit
}

// CalculateBuyPriceLimit 计算合适的买入价格
//
//	价格笼子, +2%和+0.10哪个大
//	目前使用, 当前价格+0.05
//
// Deprecated: 推荐使用 CalculatePriceCage [wangfeng on 2024/3/15 08:35]
func CalculateBuyPriceLimit(price float64) float64 {
	priceLimit := calculate_price_limit_for_buy(price, validDeclarationPriceRange, minimumPriceFluctuationUnit)
	return priceLimit
}

// calculate_price_limit_for_sell 计算合适的卖出价格
//
//	价格笼子, -2%和-0.10哪个大
//	目前使用, 当前价格-0.01
func calculate_price_limit_for_sell(last_price, price_cage_ratio, minimum_price_fluctuation_unit float64) float64 {
	// 价格笼子, -2%和-0.10哪个小
	priceLimit := min(last_price*(1-price_cage_ratio), last_price-minimum_price_fluctuation_unit)
	// 当前价格-0.01
	priceLimit = last_price - 0.01
	// 最后修订价格
	priceLimit = num.Decimal(priceLimit)
	return priceLimit
}

// CalculateSellPriceLimit 计算卖出价格笼子
//
// Deprecated: 推荐使用 CalculatePriceCage [wangfeng on 2024/3/15 08:35]
func CalculateSellPriceLimit(price float64) float64 {
	priceLimit := calculate_price_limit_for_sell(price, validDeclarationPriceRange, minimumPriceFluctuationUnit)
	return priceLimit
}

// 计算买入总费用
//
//	@param direction 交易方向
//	@param price 价格
//	@param volume 数量
//	@param align 费用是否对齐, 即四舍五入. direction=SELL的时候, align必须是true
func calculate_transaction_fee(direction Direction, price float64, volume int, align bool) (TotalFee, StampDutyFee, TransferFee, CommissionFee, MarketValue float64) {
	if volume < 1 {
		return InvalidFee, 0, 0, 0, 0
	}
	vol := float64(volume)
	amount := vol * price
	// 1. 印花税, 按照成交金额计算, 买入没有, 卖出, 0.1%
	_stamp_duty_fee := amount
	if direction == BUY {
		_stamp_duty_fee *= traderParameter.StampDutyRateForBuy
	} else if direction == SELL {
		_stamp_duty_fee *= traderParameter.StampDutyRateForSell
	} else {
		return InvalidFee, 0, 0, 0, 0
	}
	if align {
		_stamp_duty_fee = num.Decimal(_stamp_duty_fee)
	}
	// 2. 过户费, 按照股票数量, 双向, 0.06%
	_transfer_fee := vol * traderParameter.TransferRate
	if align {
		_transfer_fee = num.Decimal(_transfer_fee)
	}
	// 3. 券商佣金, 按照成交金额计算, 双向, 0.025%
	_commission_fee := amount * traderParameter.CommissionRate
	if align {
		_commission_fee = num.Decimal(_commission_fee)
	}
	if align && _commission_fee < traderParameter.CommissionMin {
		_commission_fee = traderParameter.CommissionMin
	}
	// 4. 股票市值
	_marketValue := amount
	if align {
		_marketValue = num.Decimal(_marketValue)
	}
	// 5. 计算费用
	_fee := _stamp_duty_fee + _transfer_fee + _commission_fee
	// 6. 计算总费用
	_total_fee := _fee
	if direction == BUY {
		// 买入操作, 加上股票市值
		_total_fee += _marketValue
	} else {
		// 卖出操作, 股票市值减去费用
		_marketValue -= _fee
	}
	return _total_fee, _stamp_duty_fee, _transfer_fee, _commission_fee, _marketValue
}

// EvaluateFeeForBuy 评估买入总费用
func EvaluateFeeForBuy(securityCode string, fund, price float64) *TradeFee {
	f := TradeFee{
		SecurityCode: securityCode,
		Price:        price,
		Volume:       UnknownVolume,
		Direction:    BUY,
	}
	f.Volume = f.CalculateNumToBuy(fund, price)
	return &f
}

// EvaluateFeeForSell 评估卖出费用
func EvaluateFeeForSell(securityCode string, price float64, volume int) *TradeFee {
	f := TradeFee{
		SecurityCode: securityCode,
		Price:        price,
		Volume:       volume,
		Direction:    SELL,
	}
	f.TotalFee, f.StampDutyFee, f.TransferFee, f.CommissionFee, f.MarketValue = calculate_transaction_fee(f.Direction, f.Price, f.Volume, true)
	return &f
}

// EvaluatePriceForSell 评估卖出价格
//
//	固定收益率>0, 重新评估卖出价格, 以确保卖出后的净利润不低于固定收益率
func EvaluatePriceForSell(securityCode string, price float64, volume int, fixedYield float64) *TradeFee {
	f := TradeFee{
		SecurityCode: securityCode,
		Price:        price,
		Volume:       volume,
		Direction:    SELL,
	}
	f.TotalFee, f.StampDutyFee, f.TransferFee, f.CommissionFee, f.MarketValue = calculate_transaction_fee(f.Direction, f.Price, f.Volume, true)
	if fixedYield > 0 {
		fee := (f.TotalFee-f.TransferFee)*(1+fixedYield) + f.TransferFee
		amount := float64(volume) * price
		marketValue := amount * (1 + fixedYield)
		totalFee := fee + marketValue
		fixedPrice := totalFee / float64(volume)
		f.Price = num.Decimal(fixedPrice)
		f.TotalFee, f.StampDutyFee, f.TransferFee, f.CommissionFee, f.MarketValue = calculate_transaction_fee(f.Direction, f.Price, f.Volume, true)
	}
	return &f
}

// TradeFee 交易费用
type TradeFee struct {
	Direction     Direction // 交易方向
	SecurityCode  string    // 证券代码
	Price         float64   // 价格
	Volume        int       // 数量
	StampDutyFee  float64   // 印花税, 按照成交金额计算, 双向, 默认单向, 费率0.1%
	TransferFee   float64   // 过户费, 按照股票数量, 双向, 默认是0.06%
	CommissionFee float64   // 券商佣金, 按照成交金额计算, 双向, 0.025%
	MarketValue   float64   // 股票市值
	TotalFee      float64   // 支出总费用
}

func (f *TradeFee) log() {
	fmt.Printf("trader[%s]: code=%s, 综合费用=%.02f, 委托价格=%.02f, 数量=%d, 其中印花说=%.02f, 过户费=%.02f, 佣金=%.02f, 股票=%.02f", f.Direction, f.SecurityCode, f.TotalFee, f.Price,
		f.Volume, f.StampDutyFee, f.TransferFee, f.CommissionFee, f.MarketValue)
}

// CalculateNumToBuy 估算可买股票数量
//
//	评估买入操作涉及的所有费用
//	返回100股的整数倍
func (f *TradeFee) CalculateNumToBuy(fund, price float64) int {
	f.Direction = BUY
	f.Price = price
	// 1. 计算每股费用
	_fee, _, _, _, _ := calculate_transaction_fee(f.Direction, f.Price, UnknownVolume, false)
	if _fee == InvalidFee {
		return InvalidVolume
	}
	// 2. 计算股数
	_vol := fund / _fee
	// 3. 换算成手数
	_vol = math.Floor(_vol / 100)
	// 4. 转成整数
	f.Volume = int(_vol) * 100
	// 5. 重新计算
	f.TotalFee, f.StampDutyFee, f.TransferFee, f.CommissionFee, f.MarketValue = calculate_transaction_fee(f.Direction, f.Price, f.Volume, true)
	if f.TotalFee == InvalidFee {
		return InvalidVolume
	} else if _fee > fund {
		// 如果费用超了, 则减去1手(100股)
		f.Volume -= 100
		// 重新计算交易费用
		f.TotalFee, f.StampDutyFee, f.TransferFee, f.CommissionFee, f.MarketValue = calculate_transaction_fee(f.Direction, f.Price, f.Volume, true)
	}
	return f.Volume
}

// CalculateFundFromSell 计算卖出股票后净收益
func (f *TradeFee) CalculateFundFromSell(price float64, volume int) float64 {
	f.Direction = SELL
	f.Price = price
	f.Volume = volume
	f.TotalFee, f.StampDutyFee, f.TransferFee, f.CommissionFee, f.MarketValue = calculate_transaction_fee(f.Direction, f.Price, f.Volume, true)
	if f.TotalFee == InvalidFee {
		return InvalidFee
	}
	return f.MarketValue
}
