package factors

import "gitee.com/quant1x/gotdx/quotes"

// Swift 快速接口
type Swift interface {
	Checkout(securityCode, date string)   // 捡出指定日期的缓存数据
	Update(cacheDate, featureDate string) // 更新数据
	Repair(cacheDate, featureDate string) // 回补数据
	Increase(snapshot quotes.Snapshot)    // 增量计算, 用快照增量计算特征
}
