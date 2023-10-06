package features

import (
	"gitee.com/quant1x/gotdx/quotes"
	"gitee.com/quant1x/gotdx/trading"
)

type FeatureKind = uint64

const (
// FeatureBaseXdxr FeatureKind = 0 // 基础数据-除权除息
)

// 登记所有的特征数据
const (
	FeatureHistory          = 1 << iota // 特征数据-历史
	FeatureF10                          // 特征数据-基本面
	FeatureKLineShap                    // 特征数据-K线形态等
	FeatureMovingAverage                // 特征数据-移动平均线
	FeatureBreaksThroughBox             // 特征数据-有效突破平台
	_                                   // 预留1
	_                                   // 预留2
	_                                   // 预留3
	_                                   // 预留4
	_                                   // 预留5
	_                                   // 预留6
	_                                   // 预留7
	_                                   // 预留8
	FeatureHousNo1                      // 侯总1号策略
	FeatureHousNo2                      // 侯总2号策略
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
		FeatureHistory: {Type: FeatureHistory, Key: CacheL5KeyHistory, Name: "历史特征数据"},
		FeatureF10:     {Type: FeatureF10, Key: "cache/f10", Name: "基本面"},
		FeatureHousNo1: {Type: FeatureHousNo1, Key: "", Name: "1号策略数据"},
		FeatureHousNo2: {Type: FeatureHousNo2, Key: "", Name: "2号策略数据"},
	}
)

// Feature 特征
type Feature interface {
	Factory(date string, code string) Feature                  // 工厂
	Kind() FeatureKind                                         // 类型
	Name() string                                              // 特征名称
	Key() string                                               // 缓存关键字
	Init() error                                               // 初始化, 加载配置信息
	GetDate() string                                           // 日期
	GetSecurityCode() string                                   // 证券代码
	FromHistory(history History) Feature                       // 从历史数据加载
	Update(cacheDate, featureDate string)                      // 更新数据
	Repair(code, cacheDate, featureDate string, complete bool) // 回补数据
	Increase(snapshot quotes.Snapshot) Feature                 // 增量计算, 用快照增量计算特征
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
