package config

// DataParameter 数据源参数
type DataParameter struct {
	Trans   HistoricalTradingDataParameter `name:"历史成交数据" yaml:"trans"`
	Feature FeatureParameter               `name:"feature" yaml:"feature"`
}

// GetDataConfig 取得数据配置
func GetDataConfig() DataParameter {
	return GlobalConfig.Data
}

// HistoricalTradingDataParameter 历史成交数据参数
type HistoricalTradingDataParameter struct {
	BeginDate string `name:"默认开始日期" yaml:"begin_date" default:"2023-10-01"`
}

// FeatureParameter 特征参数
type FeatureParameter struct {
	Tendency       int         `yaml:"tendency" default:"0"`                           // 策略是趋势主导还是股价主导, 默认是0, 0-股价主导,1-趋势主导,2-股价或趋势
	Wave           FeatureWave `name:"波浪" yaml:"wave"`                                 // 波浪
	CrossStarRatio float64     `name:"十字星实体占比" yaml:"cross_star_ratio" default:"0.50"` // 判断十字星, K线实体(OPEN-CLOSE)在K线长度(HIGH-LOW)中的占比
}

// FeatureWave 特征 - 波浪
type FeatureWave struct {
	Fields           FeatureWaveFields `name:"波浪检测字段" yaml:"fields"`                        // K线检测字段
	Periods          int               `name:"周期数" yaml:"periods" default:"89"`             // 波浪检测K线周期数, 默认89天
	ReferencePeriods int               `name:"均线参照周期" yaml:"reference_periods" default:"5"` // 趋势转变参考均线的周期数, 默认是5日均线
}

// 波浪的数据字段
type FeatureWaveFields struct {
	Peak   string `yaml:"peak" default:"close"`   // K线检测 - 波峰字段, 默认是收盘价
	Valley string `yaml:"valley" default:"close"` // K线检测 - 波谷字段, 默认是收盘价
}
