package factors

import (
	"context"
	"gitee.com/quant1x/engine/cache"
	"gitee.com/quant1x/engine/datasets/base"
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
type History struct {
	cache.DataSummary `dataframe:"-"`
	Date              string         `name:"日期" dataframe:"date"`           // 日期, 数据落地的日期
	Code              string         `name:"代码" dataframe:"code"`           // 代码
	MA3               float64        `name:"3日均价" dataframe:"ma3"`          // 3日均价
	MV3               float64        `name:"3日均量" dataframe:"mv3"`          // 3日均量
	MA5               float64        `name:"5日均价" dataframe:"ma5"`          // 5日均价
	MV5               float64        `name:"5日均量" dataframe:"mv5"`          // 5日均量
	MA10              float64        `name:"10日均价" dataframe:"ma10"`        // 10日均价
	MV10              float64        `name:"10日均量" dataframe:"mv10"`        // 10日均量
	MA20              float64        `name:"20日均价" dataframe:"ma20"`        // 20日均价
	MV20              float64        `name:"20日均量" dataframe:"mv20"`        // 20日均量
	QSFZ              bool           `name:"QSFZ: 反转信号" dataframe:"qsfz"`   // QSFZ: 反转信号
	CP                float64        `name:"QSFZ: 股价涨幅" dataframe:"cp"`     // QSFZ: 股价涨幅
	CV                float64        `name:"QSFZ: 成交量涨幅" dataframe:"cv"`    // QSFZ: 成交量涨幅
	VP                float64        `name:"QSFZ: 价量比" dataframe:"vp"`      // QSFZ: 价量比
	VP3               float64        `name:"QSFZ: 3日价量比" dataframe:"vp_3"`  // QSFZ: 3日价量比
	VP5               float64        `name:"QSFZ: 5日价量比" dataframe:"vp_5"`  // QSFZ: 5日价量比
	Payloads          IncompleteData `name:"payloads" dataframe:"payloads"` // 扩展的半成品数据
	Last              CompleteData   `name:"last" dataframe:"last"`         // 上一个交易日的数据
	UpdateTime        string         `name:"更新时间" dataframe:"update_time"`  // 更新时间
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
	if len(klines) == 0 {
		return
	}
	df := pandas.LoadStructs(klines)
	var (
		OPEN  = df.ColAsNDArray("open")
		CLOSE = df.ColAsNDArray("close")
		HIGH  = df.ColAsNDArray("high")
		LOW   = df.ColAsNDArray("low")
		VOL   = df.ColAsNDArray("volume")
		//AMOUNT = df.ColAsNDArray("amount")
	)
	//MA3        float64 // 3日均价
	ma3 := MA(CLOSE, 3)
	this.MA3 = SeriesIndexOf(ma3, -1)
	//	MV3        float64 // 3日均量
	mv3 := MA(VOL, 3)
	this.MV3 = SeriesIndexOf(mv3, -1)
	//	MA5        float64 // 5日均价
	ma5 := MA(CLOSE, 5)
	this.MA5 = SeriesIndexOf(ma5, -1)
	//	MV5        float64 // 5日均量
	mv5 := MA(VOL, 5)
	this.MV5 = SeriesIndexOf(mv5, -1)
	//	MA10       float64 // 10日均价
	ma10 := MA(CLOSE, 10)
	this.MA10 = SeriesIndexOf(ma10, -1)
	//	MV10       float64 // 10日均量
	mv10 := MA(VOL, 10)
	this.MV10 = SeriesIndexOf(mv10, -1)
	//	MA20       float64 // 20日均价
	ma20 := MA(CLOSE, 20)
	this.MA20 = SeriesIndexOf(ma20, -1)
	//	MV20       float64 // 20日均量
	mv20 := MA(VOL, 20)
	this.MV20 = SeriesIndexOf(mv20, -1)
	// 扩展数据 修复
	{
		// hous_no1
		this.Payloads.No1.Repair(securityCode, cacheDate, featureDate, false)
		this.Last.No1.Repair(securityCode, cacheDate, featureDate, true)
	}
	_ = code
	_ = OPEN
	_ = CLOSE
	_ = HIGH
	_ = LOW
	_ = VOL
	_ = complete
}

func (this *History) Increase(snapshot quotes.Snapshot) Feature {
	//TODO implement me
	panic("implement me")
}

// GetMV5 前5日分钟均量
func (this *History) GetMV5() float64 {
	minutes := trading.Minutes(this.GetDate())
	if minutes < 1 {
		minutes = 1
	}
	return this.MV5 / float64(minutes)
}
