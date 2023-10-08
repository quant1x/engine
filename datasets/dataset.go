package datasets

import (
	"gitee.com/quant1x/gotdx/quotes"
)

type DataKind = uint64

const (
	BaseKLine            DataKind = 1 << iota // 基础数据-基础K线
	BaseTransaction                           // 基础数据-历史成交
	BaseMinutes                               // 基础数据-分时数据
	BaseXdxr                                  // 基础数据-除权除息
	BaseQuarterlyReports                      // 基础数据-季报
	BaseSafeScore                             // 基础数据-安全分
	BaseAggregationData  DataKind = 1 << 63
	BaseTest             DataKind = 0x8000000000000000
)

// DataSet 数据层, 数据集接口 smart
type DataSet interface {
	Kind() DataKind                         // 类型
	Name() string                           // 特征名称
	Key() string                            // 缓存关键字
	Init(barIndex *int, date string) error  // 初始化, 加载配置信息
	Filename(date, code string) string      // 缓存文件名
	Update(cacheDate, featureDate string)   // 更新数据
	Repair(cacheDate, featureDate string)   // 回补数据
	Increase(snapshot quotes.Snapshot)      // 增量计算, 用快照增量计算特征
	Clone(date string, code string) DataSet // 克隆一个DataSet
}

// DataCache 基础的数据缓存
type DataCache struct {
	Date     string // 日期
	Code     string // 证券代码
	filename string // 文件名
}

type DataSetCache struct {
	Type DataKind
	Key  string
	Name string
}

var (
	mapDataSets = map[DataKind]DataSetCache{
		BaseXdxr:             {Type: BaseXdxr, Key: "xdxr", Name: "除权除息"},
		BaseKLine:            {Type: BaseKLine, Key: "day", Name: "日K线"},
		BaseTransaction:      {Type: BaseTransaction, Key: "trans", Name: "成交数据"},
		BaseMinutes:          {Type: BaseMinutes, Key: "minutes", Name: "分时数据"},
		BaseQuarterlyReports: {Type: BaseQuarterlyReports, Key: "reports", Name: "季报"},
		BaseSafeScore:        {Type: BaseSafeScore, Key: "safescore", Name: "安全分"},
	}
)
