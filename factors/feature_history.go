package factors

import (
	"context"
	"gitee.com/quant1x/engine/cache"
	"gitee.com/quant1x/engine/datasets/base"
	"gitee.com/quant1x/engine/utils"
	"gitee.com/quant1x/gotdx/proto"
	"gitee.com/quant1x/gotdx/quotes"
	"gitee.com/quant1x/gotdx/trading"
	"gitee.com/quant1x/gox/api"
	"gitee.com/quant1x/pandas"
	. "gitee.com/quant1x/pandas/formula"
)

// IncompleteData 不完整的数据
type IncompleteData struct {
	No1 HousNo1
}

type CompleteData struct {
	No1 HousNo1
}

const (
	CacheL5KeyHistory = "history"
)

// History 历史整合数据
//
//	记录重要的截止上一个交易日的数据
type History struct {
	cache.DataSummary `dataframe:"-"`
	Date              string         `name:"日期" dataframe:"date"`           // 日期, 数据落地的日期
	Code              string         `name:"代码" dataframe:"code"`           // 代码
	MA2               float64        `name:"2日均价" dataframe:"ma2"`          // 2日均价
	MA3               float64        `name:"3日均价" dataframe:"ma3"`          // 3日均价
	MV3               float64        `name:"3日均量" dataframe:"mv3"`          // 3日均量
	MA4               float64        `name:"4日均价" dataframe:"ma4"`          // 4日均价
	MA5               float64        `name:"5日均价" dataframe:"ma5"`          // 5日均价
	MV5               float64        `name:"5日均量" dataframe:"mv5"`          // 5日均量
	MA10              float64        `name:"10日均价" dataframe:"ma10"`        // 10日均价
	MV10              float64        `name:"10日均量" dataframe:"mv10"`        // 10日均量
	MA20              float64        `name:"20日均价" dataframe:"ma20"`        // 20日均价
	MV20              float64        `name:"20日均量" dataframe:"mv20"`        // 20日均量
	HIGH              float64        `name:"最高" dataframe:"high"`           // 昨日最高
	LOW               float64        `name:"最低" dataframe:"low"`            // 昨日最低
	AveragePrice      float64        `name:"均价" dataframe:"average_price"`  // 昨日均价
	Payloads          IncompleteData `name:"payloads" dataframe:"payloads"` // 扩展的半成品数据
	Last              CompleteData   `name:"last" dataframe:"last"`         // 上一个交易日的数据
	UpdateTime        string         `name:"更新时间" dataframe:"update_time"`  // 更新时间
	State             uint64         `name:"样本状态" dataframe:"样本状态"`
}

func NewHistory(date, code string) *History {
	summary := mapFeatures[FeatureHistory]
	v := History{
		DataSummary: summary,
		Date:        date,
		Code:        code,
	}
	return &v
}

func (this *History) GetDate() string {
	return this.Date
}

func (this *History) GetSecurityCode() string {
	return this.Code
}

func (this *History) Factory(date string, code string) Feature {
	v := NewHistory(date, code)
	return v
}

func (this *History) Init(ctx context.Context, date string) error {
	_ = ctx
	_ = date
	return nil
}

func (this *History) FromHistory(history History) Feature {
	_ = api.Copy(this, &history)
	return this
}

func (this *History) Update(code, cacheDate, featureDate string, complete bool) {
	this.Repair(code, cacheDate, featureDate, complete)
}

func (this *History) Repair(code, cacheDate, featureDate string, complete bool) {
	securityCode := proto.CorrectSecurityCode(this.Code)
	tradeDate := trading.FixTradeDate(featureDate)
	klines := base.CheckoutKLines(securityCode, tradeDate)
	if len(klines) < cache.KLineMin {
		return
	}
	df := pandas.LoadStructs(klines)
	var (
		OPEN   = df.ColAsNDArray("open")
		CLOSE  = df.ColAsNDArray("close")
		HIGH   = df.ColAsNDArray("high")
		LOW    = df.ColAsNDArray("low")
		VOL    = df.ColAsNDArray("volume")
		AMOUNT = df.ColAsNDArray("amount")
	)
	ma2 := MA(CLOSE, 2)
	this.MA2 = utils.SeriesIndexOf(ma2, -1)
	//MA3        float64 // 3日均价
	ma3 := MA(CLOSE, 3)
	this.MA3 = utils.SeriesIndexOf(ma3, -1)
	//	MV3        float64 // 3日均量
	mv3 := MA(VOL, 3)
	this.MV3 = utils.SeriesIndexOf(mv3, -1)
	ma4 := MA(CLOSE, 4)
	this.MA4 = utils.SeriesIndexOf(ma4, -1)
	//	MA5        float64 // 5日均价
	ma5 := MA(CLOSE, 5)
	this.MA5 = utils.SeriesIndexOf(ma5, -1)
	//	MV5        float64 // 5日均量
	mv5 := MA(VOL, 5)
	this.MV5 = utils.SeriesIndexOf(mv5, -1)
	//	MA10       float64 // 10日均价
	ma10 := MA(CLOSE, 10)
	this.MA10 = utils.SeriesIndexOf(ma10, -1)
	//	MV10       float64 // 10日均量
	mv10 := MA(VOL, 10)
	this.MV10 = utils.SeriesIndexOf(mv10, -1)
	//	MA20       float64 // 20日均价
	ma20 := MA(CLOSE, 20)
	this.MA20 = utils.SeriesIndexOf(ma20, -1)
	//	MV20       float64 // 20日均量
	mv20 := MA(VOL, 20)
	this.MV20 = utils.SeriesIndexOf(mv20, -1)
	this.HIGH = utils.SeriesIndexOf(HIGH, -1)
	this.LOW = utils.SeriesIndexOf(LOW, -1)
	ap := AMOUNT.Div(VOL)
	this.AveragePrice = utils.SeriesIndexOf(ap, -1)
	// 扩展数据 修复
	{
		// hous_no1
		this.Payloads.No1.Repair(securityCode, cacheDate, featureDate, false)
		this.Last.No1.Repair(securityCode, cacheDate, featureDate, true)
	}
	_ = OPEN
	this.UpdateTime = GetTimestamp()
	this.State |= this.Kind()
}

func (this *History) Increase(snapshot quotes.Snapshot) Feature {
	//TODO implement me
	panic("implement me")
}

func (this *History) ValidateSample() error {
	if this.State > 0 {
		return nil
	}
	return ErrInvalidFeatureSample
}

// GetMV5 前5日分钟均量
func (this *History) GetMV5() float64 {
	//minutes := trading.Minutes(this.GetDate())
	//if minutes < 1 {
	//	minutes = 1
	//}
	minutes := trading.CN_DEFAULT_TOTALFZNUM
	return this.MV5 / float64(minutes)
}
