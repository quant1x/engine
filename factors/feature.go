package factors

import (
	"gitee.com/quant1x/engine/cache"
	"gitee.com/quant1x/gotdx/quotes"
)

// Feature 特征
type Feature interface {
	cache.Trait
	cache.Swift
	// Factory 工厂
	Factory(date string, code string) Feature
	// FromHistory 从历史数据加载
	FromHistory(history History) Feature
	//Update(code, cacheDate, featureDate string, complete bool) // 更新数据
	//Repair(code, cacheDate, featureDate string, complete bool) // 回补数据
	Increase(snapshot quotes.Snapshot) Feature // 增量计算, 用快照增量计算特征
}

// Weight 权重数据类型为64, 实际容纳63个
type Weight = uint64
