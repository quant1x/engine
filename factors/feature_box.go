package factors

import (
	"context"

	"gitee.com/quant1x/engine/cache"
	"gitee.com/quant1x/engine/indicators"
	"gitee.com/quant1x/gox/logger"
)

const (
	cacheL5KeyBox = "box"
)

// Box 平台特征数据
type Box struct {
	cache.DataSummary `dataframe:"-"`
	Code              string  `name:"证券代码" dataframe:"证券代码"`                 // 证券代码
	Date              string  `name:"数据日期" dataframe:"数据日期"`                 // 数据日期
	DoubletPeriod     int     `name:"倍量周期" dataframe:"倍量周期"`                 // 倍量周期
	DoubleHigh        float64 `name:"倍量最高" dataframe:"倍量最高"`                 // 倍量最高
	DoubleLow         float64 `name:"倍量最低" dataframe:"倍量最低"`                 // 倍量最低
	Buy               bool    `name:"买入信号" dataframe:"买入信号"`                 // 买入信号
	HalfPeriod        int     `name:"半量周期" dataframe:"半量周期"`                 // 半量周期
	HalfHigh          float64 `name:"半量最高" dataframe:"半量最高"`                 // 半量最高
	HalfLow           float64 `name:"半量最低" dataframe:"半量最低"`                 // 半量最低
	Sell              bool    `name:"卖出信号" dataframe:"卖出信号"`                 // 卖出信号
	TendencyPeriod    int     `name:"趋势周期" dataframe:"趋势周期"`                 // 趋势周期
	QSFZ              bool    `name:"QSFZ" dataframe:"QSFZ"`                 // QSFZ
	QSCP              float64 `name:"QSCP" dataframe:"QSCP"`                 // QSFZ: QSCP
	QSCV              float64 `name:"QSCV" dataframe:"QSCV"`                 // QSFZ: QSCV
	QSVP              float64 `name:"QSVP" dataframe:"QSVP"`                 // QSFZ: QSVP
	QSVP3             float64 `name:"QSVP3" dataframe:"QSVP3"`               // QSFZ: QSVP3
	QSVP5             float64 `name:"QSVP5" dataframe:"QSVP5"`               // QSFZ: QSVP5
	DkCol             float64 `name:"DkCol" dataframe:"DkCol"`               // dkqs: 能量柱
	DkD               float64 `name:"dkd" dataframe:"dkd"`                   // dkqs: 多头力量
	DkK               float64 `name:"dkk" dataframe:"dkk"`                   // dkqs: 空头力量
	DkB               bool    `name:"dkb" dataframe:"dkb"`                   // dkqs: buy
	DkS               bool    `name:"dks" dataframe:"dks"`                   // dkqs: sell
	DxDivergence      float64 `name:"dxdivergence" dataframe:"dxdivergence"` // madx: 综合发散度评估值
	DxDm0             float64 `name:"dxdm0" dataframe:"dxdm0"`               // madx: 超短线均线发散度
	DxDm1             float64 `name:"dxdm1" dataframe:"dxdm1"`               // madx: 短线均线发散度
	DxDm2             float64 `name:"dxdm2" dataframe:"dxdm2"`               // madx: 中线均线发散度
	DxB               bool    `name:"dxb" dataframe:"dxb"`                   // madx: 买入
	DxBN              int     `name:"dxbn" dataframe:"dxbn"`                 // madx: 连续DxB信号周期数
	SarPos            int     `name:"sarPos" dataframe:"sarPos"`             // 坐标位置
	SarBull           bool    `name:"sarBull" dataframe:"sarBull"`           // 当前多空
	SarAf             float64 `name:"sarAf" dataframe:"sarAf"`               // 加速因子(Acceleration Factor)
	SarEp             float64 `name:"sarEp" dataframe:"sarEp"`               // 极值点(Extreme Point)
	SarSar            float64 `name:"sarSar" dataframe:"sarSar"`             // SAR[Pos]
	SarHigh           float64 `name:"sarHigh" dataframe:"sarHigh"`           // pos周期最高价
	SarLow            float64 `name:"sarLow" dataframe:"sarLow"`             // pos周期最低价
	SarPeriod         int     `name:"sarPeriod" dataframe:"sarPeriod"`       // 周期数, 上涨趋势, 周期数大于0, 下跌趋势, 周期数小于0, 绝对值就是已过多少天
	UpdateTime        string  `name:"更新时间" dataframe:"update_time"`          // 更新时间
	State             uint64  `name:"样本状态" dataframe:"样本状态"`                 // 样本状态
}

func NewBox(date, code string) *Box {
	summary := __mapFeatures[FeatureBreaksThroughBox]
	v := Box{
		DataSummary: summary,
		Date:        date,
		Code:        code,
	}
	return &v
}

func (this *Box) GetDate() string {
	return this.Date
}

func (this *Box) GetSecurityCode() string {
	return this.Code
}

func (this *Box) Factory(date string, code string) Feature {
	v := NewBox(date, code)
	return v
}

func (this *Box) Init(ctx context.Context, date string) error {
	_ = ctx
	_ = date
	return nil
}

func (this *Box) FromHistory(history History) Feature {
	_ = history
	return this
}

func (this *Box) Update(code, cacheDate, featureDate string, complete bool) {
	cover := NewKLineBox(code, featureDate)
	if cover == nil {
		logger.Errorf("code[%s, %s] kline not found", code, featureDate)
		return
	}
	info := this
	info.Date = cover.Date
	info.DoubletPeriod = cover.DoubletPeriod
	info.DoubleHigh = cover.DoubleHigh
	info.DoubleLow = cover.DoubleLow
	info.Buy = cover.Buy
	info.HalfPeriod = cover.HalfPeriod
	info.HalfHigh = cover.HalfHigh
	info.HalfLow = cover.HalfLow
	info.Sell = cover.Sell
	info.TendencyPeriod = cover.TendencyPeriod

	// 趋势反转
	info.QSFZ = cover.QSFZ
	info.QSCP = cover.QSCP
	info.QSCV = cover.QSCV
	info.QSVP = cover.QSVP
	info.QSVP3 = cover.QSVP3
	info.QSVP5 = cover.QSVP5

	// 多空趋势
	info.DkCol = cover.DkCol
	info.DkD = cover.DkD
	info.DkK = cover.DkK
	info.DkB = cover.DkB
	info.DkS = cover.DkS

	// 均线发散
	info.DxDm0 = cover.DxDm0
	info.DxDm1 = cover.DxDm1
	info.DxDm2 = cover.DxDm2
	info.DxDivergence = cover.DxDivergence
	info.DxB = cover.DxB
	info.DxBN = cover.DxBN

	// SAR
	info.SarPos = cover.SarPos
	info.SarBull = cover.SarBull
	info.SarAf = cover.SarAf
	info.SarEp = cover.SarEp
	info.SarSar = cover.SarSar
	info.SarHigh = cover.SarHigh
	info.SarLow = cover.SarLow
	info.SarPeriod = cover.SarPeriod
	// 样本状态
	info.State |= cover.Kind()
	info.UpdateTime = GetTimestamp()
	_ = cacheDate
	_ = complete
}

func (this *Box) Repair(code, cacheDate, featureDate string, complete bool) {
	this.Update(code, cacheDate, featureDate, complete)
}

func (this *Box) Increase(snapshot QuoteSnapshot) Feature {
	_ = snapshot
	return this
}

// ValidateSample 验证样本数据
func (this *Box) ValidateSample() error {
	if this.State > 0 {
		return nil
	}
	return ErrInvalidFeatureSample
}

// SarIncr SAR增量计算
func (this *Box) SarIncr(snapshot QuoteSnapshot) indicators.FeatureSar {
	sar := indicators.FeatureSar{
		Pos:    this.SarPos,
		Bull:   this.SarBull,
		Af:     this.SarAf,
		Ep:     this.SarEp,
		Sar:    this.SarSar,
		High:   this.SarHigh,
		Low:    this.SarLow,
		Period: this.SarPeriod,
	}
	return sar.Incr(snapshot.High, snapshot.Low)
}
