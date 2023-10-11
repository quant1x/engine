package factors

import (
	"gitee.com/quant1x/engine/cache"
)

type FeatureCache struct {
	Type FeatureKind
	Key  string
	Name string
}

const (
	baseFeature FeatureKind = cache.PluginMaskFeature // 特征类型基础编码
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
	mapFeatures = map[FeatureKind]FeatureCache{
		FeatureHistory: {Type: FeatureHistory, Key: CacheL5KeyHistory, Name: "历史特征数据"},
		FeatureF10:     {Type: FeatureF10, Key: CacheL5KeyF10, Name: "基本面"},
		FeatureHousNo1: {Type: FeatureHousNo1, Key: "", Name: "1号策略数据"},
		FeatureHousNo2: {Type: FeatureHousNo2, Key: "", Name: "2号策略数据"},
		//FeatureBreaksThroughBox: {Type: FeatureBreaksThroughBox, Key: CacheL5KeyBox, Name: "平台"},
	}
)
