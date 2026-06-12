package factors

import (
	"errors"

	"github.com/quant1x/engine/cache"
)

const (
	InvalidPeriod = -1              // 无效的周期
	InvalidWeight = float64(-99.99) // 无效的权重值
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
	FeatureF10                       = baseFeature + 1 // 特征数据-基本面
	FeatureHistory                   = baseFeature + 2 // 特征数据-历史
	FeatureNo1                       = baseFeature + 3 // 特征数据-1号策略
	FeatureMisc                      = baseFeature + 4 // 特征数据-Misc
	FeatureBreaksThroughBox          = baseFeature + 5 // 特征数据-box
	FeatureKLineShap                 = baseFeature + 6 // 特征数据-K线形态等
	FeatureInvestmentSentimentMaster = baseFeature + 7 // 狩猎者-情绪周期
	FeatureSecuritiesMarginTrading   = baseFeature + 8 // 融资融券
)

var (
	__mapFeatures = map[cache.Kind]cache.DataSummary{
		FeatureF10:                       cache.Summary(FeatureF10, cacheL5KeyF10, "基本面", cache.DefaultDataProvider),
		FeatureHistory:                   cache.Summary(FeatureHistory, cacheL5KeyHistory, "历史数据", cache.DefaultDataProvider),
		FeatureMisc:                      cache.Summary(FeatureMisc, cacheL5KeyMisc, "交易数据集合", cache.DefaultDataProvider),
		FeatureBreaksThroughBox:          cache.Summary(FeatureBreaksThroughBox, cacheL5KeyBox, "有效突破平台", cache.DefaultDataProvider),
		FeatureInvestmentSentimentMaster: cache.Summary(FeatureInvestmentSentimentMaster, cacheL5KeyInvestmentSentimentMaster, "情绪大师", cache.DefaultDataProvider),
		FeatureSecuritiesMarginTrading:   cache.Summary(FeatureSecuritiesMarginTrading, cacheL5KeySecuritiesMarginTrading, "融资融券", cache.DefaultDataProvider),
	}
)

var (
	ErrInvalidFeatureSample = errors.New("无效的特征数据样本")
)
