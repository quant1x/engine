package config

// DataParameter 数据源参数
type DataParameter struct {
	Trans HistoricalTradingDataParameter `name:"历史成交数据" yaml:"trans"`
}

// HistoricalTradingDataParameter 历史成交数据参数
type HistoricalTradingDataParameter struct {
	BeginDate string `name:"默认开始日期" yaml:"begin_date" default:"2023-10-01"`
}

// GetDataConfig 取得数据配置
func GetDataConfig() DataParameter {
	return GlobalConfig.DataSource
}
