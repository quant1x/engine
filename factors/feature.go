package factors

import (
	"gitee.com/quant1x/gotdx/quotes"
)

type FeatureKind = uint64

// Feature 特征
type Feature interface {
	Factory(date string, code string) Feature                  // 工厂
	Kind() FeatureKind                                         // 类型
	FeatureName() string                                       // 特征名称
	Key() string                                               // 缓存关键字
	Init(barIndex *int, date string) error                     // 初始化, 加载配置信息
	GetDate() string                                           // 日期
	GetSecurityCode() string                                   // 证券代码
	FromHistory(history History) Feature                       // 从历史数据加载
	Update(code, cacheDate, featureDate string, complete bool) // 更新数据
	Repair(code, cacheDate, featureDate string, complete bool) // 回补数据
	Increase(snapshot quotes.Snapshot) Feature                 // 增量计算, 用快照增量计算特征
	//Assign(target *cachel5.CacheAdapter) bool   // 赋值
}

// Weight 权重数据类型为64, 实际容纳63个
type Weight = uint64
