package base

import (
	"gitee.com/quant1x/gotdx"
	"gitee.com/quant1x/gotdx/proto"
	"gitee.com/quant1x/gotdx/quotes"
	"gitee.com/quant1x/gotdx/trading"
	"gitee.com/quant1x/gox/logger"
	"gitee.com/quant1x/pandas/stat"
)

const (
	kTradingFirstTime        = "09:25"      // 第一个时间
	kTradingStartTime        = "09:30"      // 开盘时间
	kTradingFinalBiddingTime = "14:57"      // 尾盘集合竞价时间
	kTradingLastTime         = "15:00"      // 最后一个时间
	TickDefaultStartDate     = "2023-01-01" // 分笔成交最早的日期
)

// Transaction 获取指定日期的历史成交数据
func Transaction(securityCode, tradeDate string) []quotes.TickTransaction {
	securityCode = proto.CorrectSecurityCode(securityCode)
	tdxApi := gotdx.GetTdxApi()
	offset := uint16(quotes.TDX_TRANSACTION_MAX)
	start := uint16(0)
	history := make([]quotes.TickTransaction, 0)
	hs := make([]quotes.TransactionReply, 0)
	date := trading.FixTradeDate(tradeDate, TDX_FORMAT_PROTOCOL_DATE)
	iDate := stat.AnyToInt64(date)
	for {
		var data *quotes.TransactionReply
		var err error
		retryTimes := 0
		for retryTimes < quotes.DefaultRetryTimes {
			data, err = tdxApi.GetHistoryTransactionData(securityCode, uint32(iDate), start, offset)
			if err == nil && data != nil {
				break
			}
			retryTimes++
		}
		if err != nil {
			logger.Errorf("code=%s, tradeDate=%s, error=%s", securityCode, tradeDate, err.Error())
			return []quotes.TickTransaction{}
		}
		if data == nil || data.Count == 0 {
			break
		}
		// 历史成交记录是按照时间排序
		//data.List = stat.Reverse(data.List)
		hs = append(hs, *data)
		if data.Count < offset {
			break
		}
		start += offset
	}
	// 这里需要反转一下
	hs = stat.Reverse(hs)
	for _, v := range hs {
		history = append(history, v.List...)
	}

	return history
}

//// CountInflow 统计指定日期的内外盘
//func CountInflow(list []quotes.TickTransaction, securityCode string, date string) (summary cache.TurnoverDataSummary) {
//	if len(list) == 0 {
//		return
//	}
//	securityCode = proto.CorrectSecurityCode(securityCode)
//	for _, v := range list {
//		tm := v.Time
//		direction := int32(v.BuyOrSell)
//		price := v.Price
//		vol := int64(v.Vol)
//		// 统计内外盘数据
//		if direction == quotes.TDX_TICK_BUY {
//			// 买入
//			summary.OuterVolume += vol
//			summary.OuterAmount += float64(vol) * price
//		} else if direction == quotes.TDX_TICK_SELL {
//			// 卖出
//			summary.InnerVolume += vol
//			summary.InnerAmount += float64(vol) * price
//		} else {
//			// 可能存在中性盘2, 最近又发现有类型是3, 暂时还是按照中性盘来处理
//			vn := vol
//			buyOffset := vn / 2
//			sellOffset := vn - buyOffset
//			// 买入
//			summary.OuterVolume += buyOffset
//			summary.OuterAmount += float64(buyOffset) * price
//			// 卖出
//			summary.InnerVolume += sellOffset
//			summary.InnerAmount += float64(sellOffset) * price
//		}
//		// 计算开盘竞价数据
//		if tm >= kTradingFirstTime && tm < kTradingStartTime {
//			summary.OpenVolume += vol
//		}
//		// 计算收盘竞价数据
//		if tm > kTradingFinalBiddingTime && tm <= kTradingLastTime {
//			summary.CloseVolume += vol
//		}
//	}
//	f10 := smart.GetL5F10(securityCode, date)
//	if f10 != nil {
//		summary.OpenTurnZ = f10.TurnZ(summary.OpenVolume)
//		summary.CloseTurnZ = f10.TurnZ(summary.CloseVolume)
//	}
//
//	return
//}
