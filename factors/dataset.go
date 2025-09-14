package factors

import (
	"github.com/quant1x/engine/cache"
	"github.com/quant1x/engine/config"
	"github.com/quant1x/gotdx/quotes"
	"github.com/quant1x/x/logger"
)

const (
	baseKind cache.Kind = 0
)

const (
	BaseXdxr                = cache.PluginMaskBaseData | (baseKind + 1)  // 基础数据-除权除息
	BaseKLine               = cache.PluginMaskBaseData | (baseKind + 2)  // 基础数据-基础K线
	BaseTransaction         = cache.PluginMaskBaseData | (baseKind + 3)  // 基础数据-历史成交
	BaseMinutes             = cache.PluginMaskBaseData | (baseKind + 4)  // 基础数据-分时数据
	BaseQuarterlyReports    = cache.PluginMaskBaseData | (baseKind + 5)  // 基础数据-季报
	BaseSafetyScore         = cache.PluginMaskBaseData | (baseKind + 6)  // 基础数据-安全分
	BaseWideKLine           = cache.PluginMaskBaseData | (baseKind + 7)  // 基础数据-宽表
	BasePerformanceForecast = cache.PluginMaskBaseData | (baseKind + 8)  // 基础数据-业绩预告
	BaseChipDistribution    = cache.PluginMaskBaseData | (baseKind + 9)  // 基础数据-筹码分布
	BaseKLineMinute         = cache.PluginMaskBaseData | (baseKind + 10) // 基础数据-基础分钟级别K线
)

// DataSet 数据层, 数据集接口 smart
//
//	数据集是基础数据, 应当遵循结构简单, 尽量减小缓存的文件数量, 加载迅速
//	检索的规则是按日期和代码进行查询
type DataSet interface {
	cache.Manifest
	// Clone 克隆一个DataSet, 是所有写操作的基础
	Clone(featureDate, securityCode string) DataSet
	// Update 更新数据
	Update(featureDate string) error
	// Repair 回补数据
	Repair(featureDate string) error
	// Increase 增量计算, 用快照增量计算特征
	Increase(snapshot quotes.Snapshot) error
}

var (
	// 数据集集合
	__mapDataSets = map[cache.Kind]cache.DataSummary{
		BaseXdxr:                cache.Summary(BaseXdxr, "xdxr", "除权除息", cache.DefaultDataProvider),
		BaseKLine:               cache.Summary(BaseKLine, "day", "日K线", cache.DefaultDataProvider),
		BaseTransaction:         cache.Summary(BaseTransaction, "trans", "成交数据", cache.DefaultDataProvider, "默认最早日期"+config.GetDataConfig().Trans.BeginDate),
		BaseMinutes:             cache.Summary(BaseMinutes, "minutes", "分时数据", cache.DefaultDataProvider),
		BaseQuarterlyReports:    cache.Summary(BaseQuarterlyReports, "reports", "季报", cache.DefaultDataProvider),
		BaseSafetyScore:         cache.Summary(BaseSafetyScore, "safetyscore", "安全分", cache.DefaultDataProvider),
		BaseWideKLine:           cache.Summary(BaseWideKLine, "wide", "宽表", cache.DefaultDataProvider),
		BasePerformanceForecast: cache.Summary(BasePerformanceForecast, "forecast", "业绩预告", cache.DefaultDataProvider),
		BaseChipDistribution:    cache.Summary(BaseChipDistribution, "chips", "筹码分布", cache.DefaultDataProvider),
		BaseKLineMinute:         cache.Summary(BaseKLineMinute, "min", "分钟级K线", cache.DefaultDataProvider, "支持1min,5min,15min,30min,60min"),
	}
)

func GetDataDescript(kind cache.Kind) cache.DataSummary {
	v, ok := __mapDataSets[kind]
	if !ok {
		logger.Fatalf("类型不存在, name=%d", kind)
	}
	return v
}

type Manifest struct {
	cache.DataSummary
	Date string
	Code string
}

func (m Manifest) GetDate() string {
	return m.Date
}

func (m Manifest) GetSecurityCode() string {
	return m.Code
}

type Properties struct {
	date string
	code string
}

func NewProperties(featureDate, securityCode string) Properties {
	return Properties{
		date: featureDate,
		code: securityCode,
	}
}

func (p Properties) GetDate() string {
	return p.date
}

func (p Properties) GetSecurityCode() string {
	return p.code
}
