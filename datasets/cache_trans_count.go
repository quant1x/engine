package datasets

import (
	"gitee.com/quant1x/engine/cache"
	"gitee.com/quant1x/engine/datasets/base"
	"gitee.com/quant1x/engine/smart"
	"gitee.com/quant1x/gotdx/proto"
	"gitee.com/quant1x/gotdx/quotes"
)

//// InflowCount 统计内外盘
//func InflowCount(df pandas.DataFrame, securityCode string) (summary cache.TurnoverDataSummary) {
//	securityCode = proto.CorrectSecurityCode(securityCode)
//	if df.Nrow() > 0 {
//		for i := 0; i < df.Nrow(); i++ {
//			m := df.IndexOf(i)
//			tm := stat.AnyToString(m["time"])
//			direction := stat.AnyToInt32(m["buyorsell"])
//			price := stat.AnyToFloat64(m["price"])
//			vol := stat.AnyToInt64(m["vol"])
//
//			// 统计内外盘数据
//			if direction == quotes.TDX_TICK_BUY {
//				// 买入
//				summary.OuterVolume += vol
//				summary.OuterAmount += float64(vol) * price
//			} else if direction == quotes.TDX_TICK_SELL {
//				// 卖出
//				summary.InnerVolume += vol
//				summary.InnerAmount += float64(vol) * price
//			} else {
//				// 可能存在中性盘2, 最近又发现有类型是3, 暂时还是按照中性盘来处理
//				vn := vol
//				buyOffset := vn / 2
//				sellOffset := vn - buyOffset
//				// 买入
//				summary.OuterVolume += buyOffset
//				summary.OuterAmount += float64(buyOffset) * price
//				// 卖出
//				summary.InnerVolume += sellOffset
//				summary.InnerAmount += float64(sellOffset) * price
//			}
//			// 计算开盘竞价数据
//			if tm >= TradingFirstTime && tm < TradingStartTime {
//				summary.OpenVolume += vol
//			}
//			// 计算收盘竞价数据
//			if tm > kTradingFinalBiddingTime && tm <= kTradingLastTime {
//				summary.CloseVolume += vol
//			}
//		}
//		f10 := smart.GetL5F10(securityCode)
//		if f10 != nil {
//			freeCapital := f10.FreeCapital
//			summary.OpenTurnZ = 10000.00 * (float64(summary.OpenVolume) / freeCapital)
//			summary.CloseTurnZ = 10000.00 * (float64(summary.CloseVolume) / freeCapital)
//		}
//	}
//	return
//}

// CountInflow 统计指定日期的内外盘
func CountInflow(list []quotes.TickTransaction, securityCode string, date string) (summary cache.TurnoverDataSummary) {
	if len(list) == 0 {
		return
	}
	securityCode = proto.CorrectSecurityCode(securityCode)
	for _, v := range list {
		tm := v.Time
		direction := int32(v.BuyOrSell)
		price := v.Price
		vol := int64(v.Vol)
		// 统计内外盘数据
		if direction == quotes.TDX_TICK_BUY {
			// 买入
			summary.OuterVolume += vol
			summary.OuterAmount += float64(vol) * price
		} else if direction == quotes.TDX_TICK_SELL {
			// 卖出
			summary.InnerVolume += vol
			summary.InnerAmount += float64(vol) * price
		} else {
			// 可能存在中性盘2, 最近又发现有类型是3, 暂时还是按照中性盘来处理
			vn := vol
			buyOffset := vn / 2
			sellOffset := vn - buyOffset
			// 买入
			summary.OuterVolume += buyOffset
			summary.OuterAmount += float64(buyOffset) * price
			// 卖出
			summary.InnerVolume += sellOffset
			summary.InnerAmount += float64(sellOffset) * price
		}
		// 计算开盘竞价数据
		if tm >= base.TradingFirstTime && tm < base.TradingStartTime {
			summary.OpenVolume += vol
		}
		// 计算收盘竞价数据
		if tm > base.TradingFinalBiddingTime && tm <= base.TradingLastTime {
			summary.CloseVolume += vol
		}
	}
	f10 := smart.GetL5F10(securityCode, date)
	if f10 != nil {
		summary.OpenTurnZ = f10.TurnZ(summary.OpenVolume)
		summary.CloseTurnZ = f10.TurnZ(summary.CloseVolume)
	}

	return
}
