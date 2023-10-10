package features

import (
	"gitee.com/quant1x/engine/datasets/tdxweb"
	"gitee.com/quant1x/engine/market"
	"gitee.com/quant1x/gotdx/quotes"
	"gitee.com/quant1x/gotdx/securities"
	"gitee.com/quant1x/gox/api"
	"gitee.com/quant1x/gox/num"
	"gitee.com/quant1x/pandas/stat"
)

const (
	CacheL5KeyF10 = "cache/f10"
)

// F10 证券基本面
type F10 struct {
	Date           string  `name:"日期"`           // 日期
	Code           string  `name:"代码"`           // 代码
	Name           string  `name:"名称"`           // 名称
	SubNew         bool    `name:"次新股"`          // 是否次新股
	VolUnit        int     `name:"每手"`           // 每手单位
	DecimalPoint   int     `name:"小数点"`          // 小数点
	IpoDate        string  `name:"上市日期"`         // 上市日期
	UpdateDate     string  `name:"更新日期"`         // 更新日期
	TotalCapital   float64 `name:"总股本"`          // 总股本
	Capital        float64 `name:"流通股本"`         // 流通股本
	FreeCapital    float64 `name:"自由流通股本"`       // 自由流通股本
	Top10Capital   float64 `name:"前十大流通股东总股本"`   // 前十大流通股东股本
	Top10Change    float64 `name:"前十大流通股东总股本变化"` //前十大流通股东股本变化
	ChangeCapital  float64 `name:"前十大流通股东持仓变化"`  // 前十大流通股东持仓变化
	IncreaseRatio  float64 `name:"当期增持比例"`       // 当期增持比例
	ReductionRatio float64 `name:"当期减持比例"`       // 当期减持比例
	BPS            float64 `name:"每股净资产"`        // 每股净资产
	BasicEPS       float64 `name:"每股收益"`         // 每股收益
	SafetyScore    int     `name:"安全分"`          // 通达信安全分
	Increases      int     `name:"增持"`           // 公告-增持
	Reduces        int     `name:"减持"`           // 公告-减持
	Risk           int     `name:"风险数"`          // 公告-风险数
	RiskKeywords   string  `name:"风险关键词"`        // 公告-风险关键词
}

//func init() {
//	err := cache.Register(&F10{})
//	if err != nil {
//		panic(err)
//	}
//}

func NewF10(date, code string) *F10 {
	v := F10{
		Date:         date,
		Code:         code,
		Name:         securities.GetStockName(code),
		VolUnit:      100,
		DecimalPoint: 2,
		SubNew:       market.IsSubNewStock(code),
	}
	securityInfo, ok := securities.CheckoutSecurityInfo(code)
	if ok {
		v.VolUnit = int(securityInfo.VolUnit)
		v.DecimalPoint = int(securityInfo.DecimalPoint)
		v.Name = securityInfo.Name
	}
	return &v
}

func (f *F10) Factory(date string, code string) Feature {
	v := NewF10(date, code)
	return v
}

func (f *F10) Kind() FeatureKind {
	return FeatureF10
}

func (f *F10) FeatureName() string {
	return mapFeatures[f.Kind()].Name
}

func (f *F10) Key() string {
	return mapFeatures[f.Kind()].Key
}

func (f *F10) Init(barIndex *int, date string) error {
	loadQuarterlyReports(f.GetDate())
	return nil
}

func (f *F10) GetDate() string {
	return f.Date
}

func (f *F10) GetSecurityCode() string {
	return f.Code
}

func (f *F10) FromHistory(history History) Feature {
	_ = history
	return f
}

func (f *F10) Update(cacheDate, featureDate string) {
	securityCode := f.GetSecurityCode()

	// 1. 基本信息
	securityInfo := checkoutSecurityBasicInfo(securityCode, featureDate)
	_ = api.Copy(f, &securityInfo)
	// 2. 前十大流通股股东
	shareHolder := checkoutShareHolder(securityCode, featureDate)
	_ = api.Copy(f, shareHolder)
	// 3. 上市公司公告
	notice := getOneNotice(securityCode, featureDate)
	_ = api.Copy(f, &notice)
	// 4. 季报
	report := getQuarterlyReportSummary(securityCode)
	_ = api.Copy(f, &report)

	// 5. 安全分
	safetyScore := tdxweb.GetSafetyScore(securityCode)
	f.SafetyScore = safetyScore

	_ = cacheDate
}

func (f *F10) Repair(code, cacheDate, featureDate string, complete bool) {
	securityCode := code

	// 1. 基本信息
	securityInfo := checkoutSecurityBasicInfo(securityCode, featureDate)
	_ = api.Copy(f, &securityInfo)
	// 2. 前十大流通股股东
	shareHolder := checkoutShareHolder(securityCode, featureDate)
	_ = api.Copy(f, shareHolder)
	// 3. 上市公司公告
	notice := getOneNotice(securityCode, featureDate)
	_ = api.Copy(f, &notice)
	// 4. 季报
	report := getQuarterlyReportSummary(securityCode)
	_ = api.Copy(f, &report)

	// 5. 安全分
	//safetyScore := tdxweb.GetSafetyScore(securityCode)
	//f.SafetyScore = safetyScore

	_ = code
	_ = cacheDate
	_ = complete
}

func (f *F10) Increase(snapshot quotes.Snapshot) Feature {
	//TODO implement me
	panic("implement me")
}

func (f *F10) ChangingOverDate(date string) {
	//TODO implement me
	panic("implement me")
}

func (f *F10) TurnZ(v any) float64 {
	if f.FreeCapital == 0 {
		return 0.00
	}
	n := stat.AnyToFloat64(v)
	turnoverRateZ := num.ChangeRate(f.FreeCapital, n)
	turnoverRateZ *= 10000
	turnoverRateZ = num.Decimal(turnoverRateZ)
	return turnoverRateZ
}
