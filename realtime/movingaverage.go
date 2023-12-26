package realtime

import (
	"gitee.com/quant1x/engine/factors"
	"gitee.com/quant1x/gox/num"
	"gitee.com/quant1x/pandas/formula"
	"gitee.com/quant1x/pandas/stat"
)

// MovingAverage 计算均线范围
func MovingAverage(CLOSE, HIGH, LOW stat.Series, PN int) (ma, half, maMax, maMin stat.Series) {
	half = formula.MA(CLOSE, PN-1)
	r1Half := formula.REF(half, 1)
	maMax = r1Half.Mul(PN - 1).Add(HIGH).Div(PN)
	maMin = r1Half.Mul(PN - 1).Add(LOW).Div(PN)
	ma = r1Half.Mul(PN - 1).Add(CLOSE).Div(PN)
	return ma, half, maMax, maMin
}

// IncrementalMovingAverage 增量计算移动平均线
//
//	period 周期数
//	previousHalfValue 前period-1的平均值
//	price 现价
func IncrementalMovingAverage(previousHalfValue float64, period int, price float64) float64 {
	n := float64(period)
	sum := previousHalfValue*(n-1) + price
	ma := num.Decimal(sum / n)
	return ma
}

// DynamicMovingAverage 增量计算移动平均线的范围
//
//	period 周期数
//	previousHalfValue 前period-1的平均值
//	price 现价
func DynamicMovingAverage(previousHalfValue float64, period int, snapshot factors.QuoteSnapshot) (ma, maHigh, maLow float64) {
	ma = IncrementalMovingAverage(previousHalfValue, period, snapshot.Price)
	maHigh = IncrementalMovingAverage(previousHalfValue, period, snapshot.High)
	maLow = IncrementalMovingAverage(previousHalfValue, period, snapshot.Low)
	return ma, maHigh, maLow
}
