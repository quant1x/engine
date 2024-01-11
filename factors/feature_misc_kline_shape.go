package factors

import (
	"gitee.com/quant1x/engine/cache"
	"gitee.com/quant1x/engine/utils"
	"gitee.com/quant1x/exchange"
	"gitee.com/quant1x/gox/num"
	"gitee.com/quant1x/pandas"
	. "gitee.com/quant1x/pandas/formula"
	"gitee.com/quant1x/pandas/stat"
	"math"
)

// ShapeType K线形态类型
type ShapeType = uint64

const (
	KLineShapeYiZi                 ShapeType = 1 << iota // 1一字
	KLineShapeLimitUp                                    // 2涨停板
	KLineShapeLimitDown                                  // 3跌停板
	KLineShapeNoChanged                                  // 4平盘
	KLineShapeUp                                         // 5上涨
	KLineShapeDown                                       // 6下跌
	KLineShapeDoji                                       // 7十字星
	KLineShapeLongUpShadow                               // 8长上影线
	KLineShapeLongDownShadow                             // 9长下影线
	KLineShapeDaYangLine                                 // 10大阳线
	KLineShapeDaYinLine                                  // 11大阴线
	KLineShapeZhongYangLine                              // 12中阳线
	KLineShapeZhongYinLine                               // 13中阴线
	KLineShapeXiaoYangLine                               // 14小阳线
	KLineShapeXiaoYinLine                                // 15小阴线
	KLineShapeNotLimitUp                                 // 16炸板, 曾涨停
	KLineShapeNotLimitDown                               // 17曾跌停
	KLineShapeShrinkToHalf                               // 18第3日缩量到一半
	KLineShapeGrabChipInFinalRound                       // 19尾盘抢筹
	KLineShapeChant                                      // 20吟唱, 蓄力 配合大阳线
	KLineShapeStepBackDailyMA                            // 21日线-连续2个周期回踩均线
	KLineShapeStepBackWeeklyMA                           // 22周线-连续2个周期回踩均线
	KLineShapeStepBackMonthlyMA                          // 23月线-连续2个周期回踩均线
)

// K线实体大小
const (
	kBoxLarge      = float64(4.50) // 大
	kBoxMedium     = float64(1.50) // 中
	kBoxSmall      = float64(0.00) // 小
	kBoxDojiHeight = float64(0.50) // 十字星K线实体最大比例
)

// KLineShape K线形态
func KLineShape(df pandas.DataFrame, securityCode string) (shape ShapeType) {
	if df.Nrow() <= cache.KLineMin {
		return
	}
	var (
		VOLs   = df.ColAsNDArray("volume")
		HIGHs  = df.ColAsNDArray("high")
		OPENs  = df.ColAsNDArray("open")
		CLOSEs = df.ColAsNDArray("close")
	)

	// 1. 捡出基本数据
	m := df.IndexOf(-1)
	// 1.1 最高价
	tmpValue, _ := m["high"]
	HIGH := stat.AnyToFloat64(tmpValue)
	// 1.2 最低价
	tmpValue, _ = m["low"]
	LOW := stat.AnyToFloat64(tmpValue)
	// 1.3 开盘价
	tmpValue, _ = m["open"]
	OPEN := stat.AnyToFloat64(tmpValue)
	// 1.4 收盘价
	tmpValue, _ = m["close"]
	CLOSE := stat.AnyToFloat64(tmpValue)
	// 1.5 昨日收盘价
	//tmpValue, _ = m["last_close"]
	//LAST_CLOSE := stat.AnyToFloat64(tmpValue)
	LAST_CLOSE := utils.SeriesIndexOf(CLOSEs, -2)

	// 2. 判断基本形态
	if CLOSE > LAST_CLOSE {
		// 2.1 上涨
		shape |= KLineShapeUp
	} else if CLOSE < LAST_CLOSE {
		// 2.2 下跌
		shape |= KLineShapeDown
	} else {
		// 2.3 平盘
		shape |= KLineShapeNoChanged
	}

	// 3. 一字
	if OPEN == CLOSE && HIGH == LOW && OPEN == HIGH {
		// 3.1 一字
		shape |= KLineShapeYiZi
	}
	// 4 判断是否涨跌停板
	limitRate := exchange.MarketLimit(securityCode)
	lastClose := num.Decimal(LAST_CLOSE)
	lastHigh := num.Decimal(HIGH)
	lastLow := num.Decimal(LOW)
	priceLimitUp := num.Decimal(lastClose * (1.000 + limitRate))
	priceLimitDown := num.Decimal(lastClose * (1.000 - limitRate))
	price := num.Decimal(CLOSE)
	if price == priceLimitUp {
		// 4.1 涨停板
		shape |= KLineShapeLimitUp
	} else if lastHigh == priceLimitUp {
		// 4.2 炸板
		shape |= KLineShapeNotLimitUp
	} else if price == priceLimitDown {
		// 4.2 跌停板
		shape |= KLineShapeLimitDown
	} else if lastLow == priceLimitDown {
		// 4.4 曾跌停
		shape |= KLineShapeNotLimitDown
	}
	// 5. 十字星
	boxHeight := num.ChangeRate(LAST_CLOSE, math.Abs(OPEN-CLOSE)) * 100
	doji := math.Abs(boxHeight) < kBoxDojiHeight && HIGH > LOW
	if doji {
		shape |= KLineShapeDoji
	}
	// 6. 计算K线实体
	boxTop := CLOSE
	boxBottom := OPEN
	// 6.1 确定K线实体顶部和底部
	if CLOSE < OPEN {
		boxTop, boxBottom = boxBottom, boxTop
	}
	// 7. 计算影线长度
	shadowTop := HIGH - boxTop
	shadowBottom := boxBottom - LOW
	if shadowTop > shadowBottom {
		// 7.1 长上影线
		shape |= KLineShapeLongUpShadow
	} else if shadowTop < shadowBottom {
		// 7.2 长下影线
		shape |= KLineShapeLongDownShadow
	}
	// 8. 判断K线的大小
	changeRate := num.NetChangeRate(OPEN, CLOSE)
	if changeRate >= kBoxLarge {
		// 8.1 大阳线
		shape |= KLineShapeDaYangLine
	} else if changeRate <= -kBoxLarge {
		// 8.2 大阴线
		shape |= KLineShapeDaYinLine
	} else if changeRate >= kBoxMedium {
		// 8.3 中阳线
		shape |= KLineShapeZhongYangLine
	} else if changeRate <= -kBoxMedium {
		// 8.4 中阴线
		shape |= KLineShapeZhongYinLine
	} else if changeRate >= kBoxSmall {
		// 8.5 小阳线
		shape |= KLineShapeXiaoYangLine
	} else if changeRate <= -kBoxSmall {
		// 8.6 小阴线
		shape |= KLineShapeXiaoYinLine
	}
	// 9. 成交量判断
	hv3 := HHV(VOLs, 5) // 5日内最大的量
	hh3 := HHV(HIGHs, 5)
	hx1 := REF(VOLs, 2).Eq(hv3)
	hx2 := REF(HIGHs, 2).Eq(hh3)
	hvn := BARSLAST(hx1.And(hx2))
	x1 := hvn.Eq(0) // __+__
	x2 := hv3.Div(VOLs).Gte(2.00)
	x := x1.And(x2)
	//DATE := df.Col("date")
	//df1 := pandas.NewDataFrame(DATE, hv3, hvn, x1, x2, x)
	//fmt.Println(df1)
	if x.IndexOf(-1).(bool) {
		shape |= KLineShapeShrinkToHalf
	}
	// 10. 尾盘抢筹
	tmpValue, _ = m["open_turnz"]
	openTurnZ := stat.AnyToFloat64(tmpValue)
	tmpValue, _ = m["close_turnz"]
	closeTurnZ := stat.AnyToFloat64(tmpValue)
	if openTurnZ > 0 && closeTurnZ > 0 && closeTurnZ > openTurnZ {
		shape |= KLineShapeGrabChipInFinalRound
	}
	// 11. 吟唱, 蓄力
	//LZL:VOL/REF(VOL,1) -1,NODRAW;
	//BOXH:=MAX(OPEN,CLOSE);
	//BOXL:=MIN(OPEN,CLOSE);
	//C1:REF(LZL>=1,1),NODRAW;
	//C2:CLOSE>=REF((BOXH+BOXL)/2,1),NODRAW;
	//
	lzl := VOLs.Div(REF(VOLs, 1)).Sub(1.00) // 立桩量
	boxHalf := ABS(CLOSEs.Add(OPENs)).Div(2.00)
	xc1 := REF(lzl, 1).Gte(1.00)
	xc2 := CLOSEs.Gte(REF(boxHalf, 1))
	xc3 := lzl.Lt(0).And(HIGHs.Lte(REF(HIGHs, 1)))
	if xc1.And(xc2).And(xc3).IndexOf(-1).(bool) {
		shape |= KLineShapeChant
	}
	//DATE := df.Col("date")
	//df1 := pandas.NewDataFrame(DATE, xc1, xc2, xc3)
	//fmt.Println(df1)

	// 12. 趋势改变
	//MA5:=MA(CLOSE,5);
	//X1:=CLOSE>=MA5;
	//C3:BARSLASTCOUNT(X1),NODRAW;
	//B:(REF(C3,1)>=2) AND LOW<=MA5,NODRAW;
	//DRAWICON(B,50,1);
	//X2:=CLOSE<MA5;
	//C4:BARSLASTCOUNT(X2),NODRAW;
	//S:(REF(C4,1)>=2) AND HIGH>MA5,NODRAW;
	//DRAWICON(S,50,2);
	//DRAWICON(CLOSE<MA5,90,8);

	// 吊顶线
	//OUT:HIGH=MAX(OPEN,CLOSE)&&
	//HIGH-LOW>3*(HIGH-MIN(OPEN,CLOSE))&&
	//CLOSE>MA(CLOSE,5);

	// 倒转锤头
	// KSTAR:MIN(OPEN,CLOSE)=LOW&&
	//HIGH-LOW>3*(MAX(OPEN,CLOSE)-LOW)&&
	//CLOSE<MA(CLOSE,5);
	return shape
}
