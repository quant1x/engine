package factors

import (
	"github.com/quant1x/engine/datasource/base"
	"github.com/quant1x/exchange"
)

// AdjustmentExecutor 复权执行接口
type AdjustmentExecutor interface {
	Apply(factor func(p float64) float64)
	GetDate() string
}

// ApplyAdjustment 计算前复权 假定缓存中的记录都是截至当日的前一个交易日已经前复权
//
// 参数:
//   - securityCode 证券代码
//   - data 实现AdjustmentExecutor接口的数据集
//   - startDate 表示已经除权的日期
//
// 返回: 无
func ApplyAdjustment[E any](securityCode string, data []E, startDate string) {
	rows := len(data)
	if rows == 0 {
		return
	}
	startDate = exchange.FixTradeDate(startDate)
	// 复权之前, 假定当前缓存之中的数据都是复权过的数据
	// 那么就应该只拉取缓存最后1条记录之后的除权除息记录进行复权
	// 前复权adjust
	dividends := base.GetCacheXdxrList(securityCode)
	for i := 0; i < len(dividends); i++ {
		xdxr := dividends[i]
		if xdxr.Category != 1 {
			// 忽略非除权信息
			continue
		}
		xdxrDate := exchange.FixTradeDate(xdxr.Date)
		if xdxrDate <= startDate {
			// 忽略除权数据在新数据之前的除权记录
			continue
		}
		factor := xdxr.Adjust()
		for j := 0; j < rows; j++ {
			kl, ok := any(&data[j]).(AdjustmentExecutor)
			if !ok {
				continue
			}
			barCurrentDate := kl.GetDate()
			if barCurrentDate > xdxrDate {
				break
			}
			if barCurrentDate < xdxrDate {
				//kl.Open = factor(kl.Open)
				//kl.Close = factor(kl.Close)
				//kl.High = factor(kl.High)
				//kl.Low = factor(kl.Low)
				//// 成交量复权
				//// 1. 计算均价线
				//maPrice := kl.Amount / kl.Volume
				//// 2. 均价线复权
				//// 通达信中可能存在没有量复权的情况, 需要在系统设置中的"设置1"勾选分析图中成交量复权
				//maPrice = factor(maPrice)
				//// 3. 以成交金额为基准, 用复权均价线计算成交量
				//kl.Volume = kl.Amount / maPrice
				kl.Apply(factor)
			}
			if barCurrentDate == xdxrDate {
				break
			}
		}
	}
}
