package factors

import (
	"strconv"
	"time"

	"github.com/quant1x/engine/datasource/base"
	"github.com/quant1x/exchange"
	"github.com/quant1x/gotdx"
	"github.com/quant1x/gotdx/quotes"
	"github.com/quant1x/num"
	"github.com/quant1x/x/logger"
)

// 获取财务数据
func getFinanceInfo(securityCode, featureDate string) (capital, totalCapital float64, ipoDate, updateDate string) {
	basicDate := uint32(num.AnyToInt64(exchange.MARKET_CN_FIRST_DATE))
	for i := 0; i < quotes.DefaultRetryTimes; i++ {
		securityCode := exchange.CorrectSecurityCode(securityCode)
		tdxApi := gotdx.GetTdxApi()
		info, err := tdxApi.GetFinanceInfo(securityCode)
		if err != nil {
			logger.Error(err)
			gotdx.ReOpen()
		}
		if info != nil {
			if info.LiuTongGuBen > 0 && info.ZongGuBen > 0 {
				capital = info.LiuTongGuBen
				totalCapital = info.ZongGuBen
			}
			if info.IPODate >= basicDate {
				ipoDate = strconv.FormatInt(int64(info.IPODate), 10)
				ipoDate = exchange.FixTradeDate(ipoDate)
			} else {
				ipoDate = getIpoDate(securityCode, featureDate)
			}
			if info.UpdatedDate >= basicDate {
				updateDate = strconv.FormatInt(int64(info.UpdatedDate), 10)
				updateDate = exchange.FixTradeDate(updateDate)
			}
			break
		} else if i+1 < quotes.DefaultRetryTimes {
			time.Sleep(time.Millisecond * 10)
		}
	}
	return
}

func getIpoDate(securityCode, featureDate string) (ipoDate string) {
	// IPO日期不存在, 从日K线第一条记录获取
	kls := base.CheckoutKLines(securityCode, featureDate)
	if len(kls) > 0 {
		ipoDate = kls[0].Date
	}
	return
}
