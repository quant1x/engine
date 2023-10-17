package factors

import (
	"gitee.com/quant1x/engine/cache"
	"gitee.com/quant1x/gotdx/quotes"
)

// Feature 特征
type Feature interface {
	cache.Trait
	cache.Future
	//cache.Swift

	// Factory 工厂
	Factory(date string, code string) Feature

	// FromHistory 从历史数据加载
	FromHistory(history History) Feature

	// Increase 增量计算
	//	用快照增量计算特征
	Increase(snapshot quotes.Snapshot) Feature
}

// Weight 权重数据类型为64, 实际容纳63个
type Weight = uint64
