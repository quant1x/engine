package factors

import (
	"context"
	"gitee.com/quant1x/engine/cache"
	"gitee.com/quant1x/engine/datasets/tdxweb"
	"gitee.com/quant1x/engine/market"
	"gitee.com/quant1x/gotdx/securities"
	"gitee.com/quant1x/gox/api"
	"gitee.com/quant1x/gox/num"
	"gitee.com/quant1x/pandas/stat"
)

const (
	CacheL5KeyF10 = "f10"
)

// F10 证券基本面
type F10 struct {
	cache.DataSummary    `dataframe:"-"`
	Date                 string  `name:"日期" dataframe:"Date"`                    // 日期
	Code                 string  `name:"代码" dataframe:"Code"`                    // 证券代码
	SecurityName         string  `name:"名称" dataframe:"Name"`                    // 证券名称
	SubNew               bool    `name:"次新股" dataframe:"SubNew"`                 // 是否次新股
	VolUnit              int     `name:"每手" dataframe:"VolUnit"`                 // 每手单位
	DecimalPoint         int     `name:"小数点" dataframe:"DecimalPoint"`           // 小数点
	IpoDate              string  `name:"上市日期" dataframe:"IpoDate"`               // 上市日期
	UpdateDate           string  `name:"更新日期" dataframe:"UpdateDate"`            // 更新日期
	TotalCapital         float64 `name:"总股本" dataframe:"TotalCapital"`           // 总股本
	Capital              float64 `name:"流通股本" dataframe:"Capital"`               // 流通股本
	FreeCapital          float64 `name:"自由流通股本" dataframe:"FreeCapital"`         // 自由流通股本
	Top10Capital         float64 `name:"前十大流通股东总股本" dataframe:"Top10Capital"`    // 前十大流通股东股本
	Top10Change          float64 `name:"前十大流通股东总股本变化" dataframe:"Top10Change"`   //前十大流通股东股本变化
	ChangeCapital        float64 `name:"前十大流通股东持仓变化" dataframe:"ChangeCapital"`  // 前十大流通股东持仓变化
	IncreaseRatio        float64 `name:"当期增持比例" dataframe:"IncreaseRatio"`       // 当期增持比例
	ReductionRatio       float64 `name:"当期减持比例" dataframe:"ReductionRatio"`      // 当期减持比例
	QuarterlyYearQuarter string  `name:"季报期" dataframe:"quarterly_year_quarter"` // 当前市场处于哪个季报期, 用于比较个股的季报数据是否存在拖延的情况
	QDate                string  `name:"新报告期" dataframe:"qdate"`                 // 最新报告期
	TotalOperateIncome   float64 `name:"营业总收入" dataframe:"TotalOperateIncome"`   // 当期营业总收入
	BPS                  float64 `name:"每股净资产" dataframe:"BPS"`                  // 每股净资产
	BasicEPS             float64 `name:"每股收益" dataframe:"BasicEPS"`              // 每股收益
	DeductBasicEPS       float64 `name:"每股收益(扣除)" dataframe:"DeductBasicEPS"`    // 每股收益(扣除)
	SafetyScore          int     `name:"安全分" dataframe:"SafetyScore"`            // 通达信安全分
	Increases            int     `name:"增持" dataframe:"Increases"`               // 公告-增持
	Reduces              int     `name:"减持" dataframe:"Reduces"`                 // 公告-减持
	Risk                 int     `name:"风险数" dataframe:"Risk"`                   // 公告-风险数
	RiskKeywords         string  `name:"风险关键词" dataframe:"RiskKeywords"`         // 公告-风险关键词
	UpdateTime           string  `name:"更新时间" dataframe:"update_time"`           // 更新时间
	State                uint64  `name:"样本状态" dataframe:"样本状态"`
}

func NewF10(date, code string) *F10 {
	summary := mapFeatures[FeatureF10]
	v := F10{
		DataSummary:  summary,
		Date:         date,
		Code:         code,
		SecurityName: securities.GetStockName(code),
		VolUnit:      100,
		DecimalPoint: 2,
		SubNew:       market.IsSubNewStock(code),
	}
	securityInfo, ok := securities.CheckoutSecurityInfo(code)
	if ok {
		v.VolUnit = int(securityInfo.VolUnit)
		v.DecimalPoint = int(securityInfo.DecimalPoint)
		v.SecurityName = securityInfo.Name
	}
	return &v
}

func (this *F10) GetDate() string {
	return this.Date
}

func (this *F10) GetSecurityCode() string {
	return this.Code
}

func (this *F10) Factory(date string, code string) Feature {
	v := NewF10(date, code)
	return v
}

func (this *F10) Init(ctx context.Context, date string) error {
	loadQuarterlyReports(this.GetDate())
	_ = ctx
	_ = date
	return nil
}

func (this *F10) FromHistory(history History) Feature {
	_ = history
	return this
}

func (this *F10) Update(code, cacheDate, featureDate string, complete bool) {
	securityCode := this.GetSecurityCode()

	// 1. 基本信息
	securityInfo := checkoutSecurityBasicInfo(securityCode, featureDate)
	_ = api.Copy(this, &securityInfo)
	// 2. 前十大流通股股东
	shareHolder := checkoutShareHolder(securityCode, featureDate)
	_ = api.Copy(this, shareHolder)
	// 3. 上市公司公告
	notice := getOneNotice(securityCode, featureDate)
	_ = api.Copy(this, &notice)
	// 4. 季报
	this.QuarterlyYearQuarter = getQuarterlyYearQuarter(featureDate)
	report := getQuarterlyReportSummary(securityCode, featureDate)
	_ = api.Copy(this, &report)

	// 5. 安全分
	safetyScore := tdxweb.GetSafetyScore(securityCode)
	this.SafetyScore = safetyScore

	this.UpdateTime = GetTimestamp()
	this.State |= this.Kind()

	_ = complete
}

func (this *F10) Repair(code, cacheDate, featureDate string, complete bool) {
	securityCode := code

	// 1. 基本信息
	securityInfo := checkoutSecurityBasicInfo(securityCode, featureDate)
	_ = api.Copy(this, &securityInfo)
	// 2. 前十大流通股股东
	shareHolder := checkoutShareHolder(securityCode, featureDate)
	_ = api.Copy(this, shareHolder)
	// 3. 上市公司公告
	notice := getOneNotice(securityCode, featureDate)
	_ = api.Copy(this, &notice)
	// 4. 季报
	report := getQuarterlyReportSummary(securityCode, featureDate)
	_ = api.Copy(this, &report)

	// 5. 安全分
	//safetyScore := tdxweb.GetSafetyScore(securityCode)
	//this.SafetyScore = safetyScore

	this.UpdateTime = GetTimestamp()
	this.State |= this.Kind()

	_ = complete
}

func (this *F10) Increase(snapshot QuoteSnapshot) Feature {
	//TODO implement me
	panic("implement me")
}

func (this *F10) ValidateSample() error {
	if this.State > 0 {
		return nil
	}
	return ErrInvalidFeatureSample
}

func (this *F10) TurnZ(v any) float64 {
	if this.FreeCapital == 0 {
		return 0.00
	}
	n := stat.AnyToFloat64(v)
	turnoverRateZ := num.ChangeRate(this.FreeCapital, n)
	turnoverRateZ *= 10000
	turnoverRateZ = num.Decimal(turnoverRateZ)
	return turnoverRateZ
}
