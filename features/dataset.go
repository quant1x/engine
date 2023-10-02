package features

import "gitee.com/quant1x/gotdx/quotes"

// DataSet 数据层, 数据集接口
type DataSet interface {
	Kind() FeatureKind                      // 类型
	Name() string                           // 特征名称
	Key() string                            // 缓存关键字
	Filename(date, code string) string      // 缓存文件名
	Update(cacheDate, featureDate string)   // 更新数据
	Repair(cacheDate, featureDate string)   // 回补数据
	Increase(snapshot quotes.Snapshot)      // 增量计算, 用快照增量计算特征
	Clone(date string, code string) DataSet // 克隆一个DataSet
}

type DataCache struct {
	Date     string // 日期
	Code     string // 证券代码
	filename string // 文件名
}
