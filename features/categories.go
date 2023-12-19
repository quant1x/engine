package features

import (
	"gitee.com/quant1x/engine/cache"
)

const (
	baseFeature cache.Kind = cache.PluginMaskFeature + 10 // 特征类型基础编码
)

// 登记所有的特征数据
const (
	//FeatureHistory          = baseFeature + 1 // 特征数据-历史
	//FeatureNo1              = baseFeature + 2 // 特征数据-1号策略
	FeatureExchange         = baseFeature + 3 // exchange
	FeatureBreaksThroughBox = baseFeature + 4 // 特征数据-有效突破平台
	FeatureKLineShap        = baseFeature + 5 // 特征数据-K线形态等
	//FeatureHalf             = baseFeature + 0 // 半成品数据, 用于部分增量计算
	//FeatureMovingAverage    = baseFeature + 4 // 特征数据-移动平均线
)

var (
	__mapFeatures = map[cache.Kind]cache.DataSummary{
		//FeatureHistory: cache.Summary(FeatureHistory, factors.CacheL5KeyHistory, "历史特征数据", cache.DefaultDataProvider),
		//FeatureNo1:     cache.Summary(FeatureNo1, "", "1号策略数据", cache.DefaultDataProvider),
		//FeatureHalf:             cache.Summary(FeatureHalf, CacheL5KeyHalf, "半成数据集合", cache.DefaultDataProvider),
		FeatureExchange:         cache.Summary(FeatureExchange, CacheL5KeyExchange, "交易数据集合", cache.DefaultDataProvider),
		FeatureBreaksThroughBox: cache.Summary(FeatureBreaksThroughBox, CacheL5KeyBox, "有效突破平台", cache.DefaultDataProvider),
	}
)
