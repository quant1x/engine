package realtime

import (
	"gitee.com/quant1x/engine/factors"
	"gitee.com/quant1x/pandas"
	. "gitee.com/quant1x/pandas/formula"
)

// MovingAverageConvergenceDivergence 平滑异同移动平均线(Moving Average Convergence Divergence)
//
//	返回最新一条macd的5个值
//	DIF:EMA(CLOSE,SHORT)-EMA(CLOSE,LONG);
//	DEA:EMA(DIF,MID)
//	MACD:(DIF-DEA)*2,COLORSTICK;
func MovingAverageConvergenceDivergence(CLOSE pandas.Series, pShort, pLong, pMid int) (SHORT, LONG, DIF, DEA, MACD float64) {
	//dif1 := df1["d1"].(float64)
	//	dif2 := df1["d2"].(float64)
	//	dif1 = EmaIncr(lastClose, dif1, AlphaOfEMA(SHORT))
	//	dif2 = EmaIncr(lastClose, dif2, AlphaOfEMA(LONG))
	//	dif := dif1 - dif2
	//	lastDif := df1["DIF"].(float64)
	//	lastDea := df1["DEA"].(float64)
	//	fmt.Println("lastDif", lastDif)
	//	alpha := AlphaOfEMA(MID)
	//	fmt.Println("dea-alpha:", alpha)
	//	//t1 := -0.446
	//	//fmt.Println("xx:", (1-alpha)*lastDif+alpha*dif)
	//	dea := EmaIncr(dif, lastDea, AlphaOfEMA(MID))
	//	macd := (dif - dea) * 2
	s := EMA(CLOSE, pShort)
	l := EMA(CLOSE, pLong)
	dif := s.Sub(l)
	dea := EMA(dif, pMid)
	macd := dif.Sub(dea).Mul(2)
	SHORT = s.IndexOf(-1).(float64)
	LONG = l.IndexOf(-1).(float64)
	DIF = dif.IndexOf(-1).(float64)
	DEA = dea.IndexOf(-1).(float64)
	MACD = macd.IndexOf(-1).(float64)
	return
}

// IncrementalMovingAverageConvergenceDivergence 增量的MACD
//
//	price 为现价
//	lastShort, lastLong, lastDea, 缓存的短,长周期的ema, 以及最后一条dea
//	pShort, pLong, pMid, 为短期,长期和中期周期数
func IncrementalMovingAverageConvergenceDivergence(price, lastShort, lastLong, lastDea float64, pShort, pLong, pMid int) (DIF, DEA, MACD float64) {
	s := IncrementalExponentialMovingAverage(price, lastShort, AlphaOfExponentialMovingAverage(pShort))
	l := IncrementalExponentialMovingAverage(price, lastLong, AlphaOfExponentialMovingAverage(pLong))
	dif := s - l
	dea := IncrementalExponentialMovingAverage(dif, lastDea, AlphaOfExponentialMovingAverage(pMid))
	macd := (dif - dea) * 2
	DIF = dif
	DEA = dea
	MACD = macd
	return
}

// DynamicMovingAverageConvergenceDivergence 动态的MACD
//
//	返回当前值以及最高值和最低值
//	price 为现价
//	lastShort, lastLong, lastDea, 缓存的短,长周期的ema, 以及最后一条dea
//	pShort, pLong, pMid, 为短期,长期和中期周期数
func DynamicMovingAverageConvergenceDivergence(snapshot factors.QuoteSnapshot, lastShort, lastLong, lastDea float64, pShort, pLong, pMid int) (macd, macdHigh, macdLow float64) {
	_, _, macd = IncrementalMovingAverageConvergenceDivergence(snapshot.Price, lastShort, lastLong, lastDea, pShort, pLong, pMid)
	_, _, macdHigh = IncrementalMovingAverageConvergenceDivergence(snapshot.High, lastShort, lastLong, lastDea, pShort, pLong, pMid)
	_, _, macdLow = IncrementalMovingAverageConvergenceDivergence(snapshot.Low, lastShort, lastLong, lastDea, pShort, pLong, pMid)
	return
}
