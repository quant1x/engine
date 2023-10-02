package features

import (
	"gitee.com/quant1x/engine/features/base"
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

// History 历史整合数据
type History struct {
	Date       string         // 日期, 数据落地的日期
	Code       string         // 代码
	MA3        float64        // 3日均价
	MV3        float64        // 3日均量
	MA5        float64        // 5日均价
	MV5        float64        // 5日均量
	MA10       float64        // 10日均价
	MV10       float64        // 10日均量
	MA20       float64        // 20日均价
	MV20       float64        // 20日均量
	QSFZ       bool           // QSFZ: 反转信号
	CP         float64        // QSFZ: 股价涨幅
	CV         float64        // QSFZ: 成交量涨幅
	VP         float64        // QSFZ: 价量比
	VP3        float64        // QSFZ: 3日价量比
	VP5        float64        // QSFZ: 5日价量比
	Payloads   IncompleteData // 扩展的半成品数据
	Last       CompleteData   // 上一个交易日的数据
	UpdateTime string         // 更新时间
}

func NewHistory(date, code string) *History {
	v := History{
		Date: date,
		Code: code,
	}
	return &v
}

func (h *History) Factory(date string, code string) Feature {
	v := NewHistory(date, code)
	return v
}

func (h *History) Kind() FeatureKind {
	return FeatureHistory
}

func (h *History) Name() string {
	return mapFeatures[h.Kind()].Name
}

func (h *History) Key() string {
	return mapFeatures[h.Kind()].Key
}

func (h *History) Init() error {
	//TODO implement me
	panic("implement me")
}

func (h *History) GetDate() string {
	return h.Date
}

func (h *History) GetSecurityCode() string {
	return h.Code
}

func (h *History) FromHistory(history History) Feature {
	_ = api.Copy(h, &history)
	return h
}

func (h *History) Update(cacheDate, featureDate string) {
	//TODO implement me
	panic("implement me")
}

func (h *History) Repair(code, cacheDate, featureDate string, complete bool) {
	securityCode := proto.CorrectSecurityCode(h.Code)
	tradeDate := trading.FixTradeDate(featureDate)
	klines := base.CheckoutKLines(securityCode, tradeDate)
	if len(klines) == 0 {
		return
	}
	//digits := 2
	//securityInfo, ok := securities.CheckoutSecurityInfo(securityCode)
	//if ok {
	//	digits = int(securityInfo.DecimalPoint)
	//}
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
	h.MA3 = SeriesIndexOf(ma3, -1)
	//	MV3        float64 // 3日均量
	mv3 := MA(VOL, 3)
	h.MV3 = SeriesIndexOf(mv3, -1)
	//	MA5        float64 // 5日均价
	ma5 := MA(CLOSE, 5)
	h.MA5 = SeriesIndexOf(ma5, -1)
	//	MV5        float64 // 5日均量
	mv5 := MA(VOL, 5)
	h.MV5 = SeriesIndexOf(mv5, -1)
	//	MA10       float64 // 10日均价
	ma10 := MA(CLOSE, 10)
	h.MA10 = SeriesIndexOf(ma10, -1)
	//	MV10       float64 // 10日均量
	mv10 := MA(VOL, 10)
	h.MV10 = SeriesIndexOf(mv10, -1)
	//	MA20       float64 // 20日均价
	ma20 := MA(CLOSE, 20)
	h.MA20 = SeriesIndexOf(ma20, -1)
	//	MV20       float64 // 20日均量
	mv20 := MA(VOL, 20)
	h.MV20 = SeriesIndexOf(mv20, -1)

	// 扩展数据 修复
	{
		//no1 := HousNo1{
		//	MA5:  h.MA5,
		//	MA10: h.MA10,
		//	MA20: h.MA20,
		//}
		//_ = api.Copy(&h.Payloads.No1, &no1)
		h.Payloads.No1.Repair(securityCode, cacheDate, featureDate, false)
		h.Last.No1.Repair(securityCode, cacheDate, featureDate, true)
	}
	_ = OPEN
	_ = CLOSE
	_ = HIGH
	_ = LOW
	_ = VOL
}

func (h *History) Increase(snapshot quotes.Snapshot) Feature {
	//TODO implement me
	panic("implement me")
}
