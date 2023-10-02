package features

import "gitee.com/quant1x/gotdx/quotes"

// Quicker 迅速的数据接口
//
//	Deprecated: 废弃的接口
type Quicker interface {
	Update(cacheDate, featureDate string) // 更新数据
	Repair(cacheDate, featureDate string) // 回补数据
	Increase(snapshot quotes.Snapshot)    // 增量计算, 用快照增量计算特征
}
