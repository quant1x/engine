package factors

import (
	"gitee.com/quant1x/gotdx/trading"
)

type FeatureKind = uint64

const (
	FeatureBaseKLine FeatureKind = 0
)

const (
	FeatureF10              FeatureKind = 1 << iota // 基本面
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

// Feature 特征
type Feature interface {
	Kind() FeatureKind                    // 类型
	Name() string                         // 特征名称
	Init() error                          // 初始化, 加载配置信息
	Update(cacheDate, featureDate string) // 更新数据
	Repair(cacheDate, featureDate string) // 回补数据
	Increase(snapshot Snapshot)           // 增量计算, 用快照增量计算特征
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
