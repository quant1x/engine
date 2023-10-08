package features

import "gitee.com/quant1x/gotdx/quotes"

const (
	CacheL5KeyBox = "cache/box"
)

// Box 平台特征数据
type Box struct {
	Code           string  `name:"证券代码" dataframe:"证券代码"`                 // 证券代码
	Date           string  `name:"数据日期" dataframe:"数据日期"`                 // 数据日期
	DoubletPeriod  int     `name:"倍量周期" dataframe:"倍量周期"`                 // 倍量周期
	DoubleHigh     float64 `name:"倍量最高" dataframe:"倍量最高"`                 // 倍量最高
	DoubleLow      float64 `name:"倍量最低" dataframe:"倍量最低"`                 // 倍量最低
	Buy            bool    `name:"买入信号" dataframe:"买入信号"`                 // 买入信号
	HalfPeriod     int     `name:"半量周期" dataframe:"半量周期"`                 // 半量周期
	HalfHigh       float64 `name:"半量最高" dataframe:"半量最高"`                 // 半量最高
	HalfLow        float64 `name:"半量最低" dataframe:"半量最低"`                 // 半量最低
	Sell           bool    `name:"卖出信号" dataframe:"卖出信号"`                 // 卖出信号
	TendencyPeriod int     `name:"趋势周期" dataframe:"趋势周期"`                 // 趋势周期
	QSFZ           bool    `name:"QSFZ" dataframe:"QSFZ"`                 // QSFZ
	QSCP           float64 `name:"QSCP" dataframe:"QSCP"`                 // QSFZ: QSCP
	QSCV           float64 `name:"QSCV" dataframe:"QSCV"`                 // QSFZ: QSCV
	QSVP           float64 `name:"QSVP" dataframe:"QSVP"`                 // QSFZ: QSVP
	QSVP3          float64 `name:"QSVP3" dataframe:"QSVP3"`               // QSFZ: QSVP3
	QSVP5          float64 `name:"QSVP5" dataframe:"QSVP5"`               // QSFZ: QSVP5
	DkCol          float64 `name:"DkCol" dataframe:"DkCol"`               // dkqs: 能量柱
	DkD            float64 `name:"dkd" dataframe:"dkd"`                   // dkqs: 多头力量
	DkK            float64 `name:"dkk" dataframe:"dkk"`                   // dkqs: 空头力量
	DkB            bool    `name:"dkb" dataframe:"dkb"`                   // dkqs: buy
	DkS            bool    `name:"dks" dataframe:"dks"`                   // dkqs: sell
	DxDivergence   float64 `name:"dxdivergence" dataframe:"dxdivergence"` // madx: 综合发散度评估值
	DxDm0          float64 `name:"dxdm0" dataframe:"dxdm0"`               // madx: 超短线均线发散度
	DxDm1          float64 `name:"dxdm1" dataframe:"dxdm1"`               // madx:   短线均线发散度
	DxDm2          float64 `name:"dxdm2" dataframe:"dxdm2"`               // madx:   中线均线发散度
	DxB            bool    `name:"dxb" dataframe:"dxb"`                   // madx: 买入
	State          uint64  `name:"样本状态" dataframe:"样本状态"`                 // 样本状态
}

func NewBox(date, code string) *Box {
	v := Box{
		Date: date,
		Code: code,
	}
	return &v
}

func (b *Box) Factory(date string, code string) Feature {
	v := NewBox(date, code)
	return v
}

func (b *Box) Kind() FeatureKind {
	return FeatureBreaksThroughBox
}

func (b *Box) FeatureName() string {
	return mapFeatures[b.Kind()].Name
}

func (b *Box) Key() string {
	return mapFeatures[b.Kind()].Key
}

func (b *Box) Init() error {
	//TODO implement me
	panic("implement me")
}

func (b *Box) GetDate() string {
	return b.Date
}

func (b *Box) GetSecurityCode() string {
	return b.Code
}

func (b *Box) FromHistory(history History) Feature {
	//TODO implement me
	panic("implement me")
}

func (b *Box) Update(cacheDate, featureDate string) {
	//TODO implement me
	panic("implement me")
}

func (b *Box) Repair(code, cacheDate, featureDate string, complete bool) {
	//TODO implement me
	panic("implement me")
}

func (b *Box) Increase(snapshot quotes.Snapshot) Feature {
	//TODO implement me
	panic("implement me")
}
