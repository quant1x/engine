package features

import (
	"gitee.com/quant1x/gotdx/trading"
)

type FeatureKind = uint64

const (
	FeatureBaseXdxr FeatureKind = 0 // 基础数据-除权除息
)

const (
	FeatureBaseKLine        FeatureKind = 1 << iota // 基础数据-基础K线
	FeatureBaseTransaction                          // 基础数据-历史成交
	FeatureBaseMinutes                              // 基础数据-分时数据
	FeatureHistory                                  // 历史特征数据
	FeatureF10                                      // 基本面
	FeatureKLineShap                                // K线形态等
	FeatureMovingAverage                            // 移动平均线
	FeatureBreaksThroughBox                         // 有效突破平台
	_                                               // 预留1
	_                                               // 预留2
	_                                               // 预留3
	_                                               // 预留4
	_                                               // 预留5
	_                                               // 预留6
	_
	_
)

type FeatureCache struct {
	Type FeatureKind
	Key  string
	Name string
}

const (
	CacheL5KeyHistory = "cache/history"
)

var (
	mapFeatures = map[FeatureKind]FeatureCache{
		FeatureBaseXdxr:        {Type: FeatureBaseXdxr, Key: "xdxr", Name: "除权除息"},
		FeatureBaseKLine:       {Type: FeatureBaseKLine, Key: "day", Name: "日K线"},
		FeatureBaseTransaction: {Type: FeatureBaseTransaction, Key: "trans", Name: "成交数据"},
		FeatureBaseMinutes:     {Type: FeatureBaseMinutes, Key: "minutes", Name: "分时数据"},
		FeatureHistory:         {Type: FeatureHistory, Key: CacheL5KeyHistory, Name: "历史特征数据"},
		FeatureF10:             {Type: FeatureF10, Key: "cache/f10", Name: "基本面"},
	}
)

// Feature 特征
type Feature interface {
	Factory(date string, code string) Feature // 工厂
	Kind() FeatureKind                        // 类型
	Name() string                             // 特征名称
	Key() string                              // 缓存关键字
	Init() error                              // 初始化, 加载配置信息
	GetDate() string                          //  日期
	GetSecurityCode() string                  // 证券代码
	FromHistory(history History) Feature      // 从历史数据加载
	Update(cacheDate, featureDate string)     // 更新数据
	Repair(cacheDate, featureDate string)     // 回补数据
	Increase(snapshot Snapshot)               // 增量计算, 用快照增量计算特征
}

// DataBuilder 数据构建器
type DataBuilder struct {
	Name          string // 名称
	CacheDate     string // 缓存文件日期
	ResourcesDate string // 源数据日期, 一般来说源数据日期要比缓存文件的日期早一个交易日
	Build         func(allCodes []string)
}

func NewDataBuilder(name, date string, build func(allCodes []string)) *DataBuilder {
	//lastDay := trading.LastTradeDate()
	//cacheDate := trading.FixTradeDate(date)
	//if cacheDate >= lastDay {
	//	cacheDate = lastDay
	//}
	//
	cacheDate := trading.FixTradeDate(date)
	dates := trading.LastNDate(cacheDate, 1)
	if len(dates) == 0 {
		return nil
	}
	resourcesDate := dates[0]
	cacheDate = trading.NextTradeDate(resourcesDate)
	builder := &DataBuilder{
		Name:          name,
		CacheDate:     cacheDate,
		ResourcesDate: resourcesDate,
		Build:         build,
	}
	return builder
}

func (this *DataBuilder) Execute(allCodes []string) {
	this.Build(allCodes)
}
