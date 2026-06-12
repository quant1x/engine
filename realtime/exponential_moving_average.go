package realtime

import (
	"github.com/quant1x/engine/factors"
	"github.com/quant1x/pandas/formula"
)

// AlphaOfExponentialMovingAverage 计算EMA的alpha值
func AlphaOfExponentialMovingAverage(period int) float64 {
	return formula.AlphaOfEMA(period)
}

// IncrementalExponentialMovingAverage 增量计算 指数移动平均线
func IncrementalExponentialMovingAverage(now, last, alpha float64) float64 {
	return formula.EmaIncr(now, last, alpha)
}

// DynamicExponentialMovingAverage 动态EMA
//
//	返回当前值以及最高值和最低值
func DynamicExponentialMovingAverage(snapshot factors.QuoteSnapshot, last, alpha float64) (ema, emaHigh, emaLow float64) {
	ema = IncrementalExponentialMovingAverage(snapshot.Price, last, alpha)
	emaHigh = IncrementalExponentialMovingAverage(snapshot.High, last, alpha)
	emaLow = IncrementalExponentialMovingAverage(snapshot.Low, last, alpha)
	return
}
