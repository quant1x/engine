package factors

import (
	"errors"
	"gitee.com/quant1x/engine/cache"
	"time"
)

// Trait 基础的特性
//
//	这也是一个特征, 为啥起这个名字, 自己可以脑补 哈哈~
type Trait interface {
	// FromHistory 从历史数据加载
	FromHistory(history History) Feature
	// Increase 增量计算
	//	用快照增量计算特征
	Increase(snapshot QuoteSnapshot) Feature
	// ValidateSample 验证样本
	ValidateSample() error
}

// Feature 特征
type Feature interface {
	cache.Manifest
	cache.Future
	// Factory 工厂
	Factory(date string, code string) Feature
	Trait
}

// Weight 权重数据类型为64, 实际容纳63个
type Weight = uint64

const (
	baseFeature cache.Kind = cache.PluginMaskFeature // 特征类型基础编码
)

// 登记所有的特征数据
const (
	FeatureF10              = baseFeature + 1    // 特征数据-基本面
	FeatureHistory          = baseFeature + 2    // 特征数据-历史
	FeatureKLineShap        = baseFeature + 3    // 特征数据-K线形态等
	FeatureMovingAverage    = baseFeature + 4    // 特征数据-移动平均线
	FeatureBreaksThroughBox = baseFeature + 5    // 特征数据-有效突破平台
	featureHous             = baseFeature + 1000 // 侯总策略编码号段
	FeatureHousNo1          = featureHous + 1    // 侯总1号策略
	FeatureHousNo2          = featureHous + 2    // 侯总2号策略
)

var (
	mapFeatures = map[cache.Kind]cache.DataSummary{
		FeatureHistory: cache.Summary(FeatureHistory, CacheL5KeyHistory, "历史特征数据", cache.DefaultDataProvider),
		FeatureF10:     cache.Summary(FeatureF10, CacheL5KeyF10, "基本面", cache.DefaultDataProvider),
		FeatureHousNo1: cache.Summary(FeatureHousNo1, "", "1号策略数据", cache.DefaultDataProvider),
		FeatureHousNo2: cache.Summary(FeatureHousNo2, "", "2号策略数据", cache.DefaultDataProvider),
	}
)

var (
	ErrInvalidFeatureSample = errors.New("无效的特征数据样本")
)

// GetTimestamp 时间戳
//
//	格式: YYYY-MM-DD hh:mm:ss.SSS
func GetTimestamp() string {
	now := time.Now()
	return now.Format(cache.TimeStampMilli)
}
