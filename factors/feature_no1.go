package factors

import (
	"context"

	"gitee.com/quant1x/engine/cache"
	"gitee.com/quant1x/engine/datasource/base"
	"gitee.com/quant1x/engine/utils"
	"gitee.com/quant1x/exchange"
	"gitee.com/quant1x/pandas"
	. "gitee.com/quant1x/pandas/formula"
)

type HousNo1 struct {
	cache.DataSummary `dataframe:"-"`
	MA5               float64 `dataframe:"ma5"`
	MA10              float64 `dataframe:"ma10"`
	MA20              float64 `dataframe:"ma20"`
}

func (this *HousNo1) GetDate() string {
	//TODO implement me
	panic("implement me")
}

func (this *HousNo1) GetSecurityCode() string {
	//TODO implement me
	panic("implement me")
}

func (this *HousNo1) Init(ctx context.Context, date string) error {
	//TODO implement me
	panic("implement me")
}

func (this *HousNo1) Factory(date string, code string) Feature {
	//TODO implement me
	panic("implement me")
}

func (this *HousNo1) FromHistory(history History) Feature {
	return this
}

func (this *HousNo1) Update(code, cacheDate, featureDate string, complete bool) {
	//TODO implement me
	panic("implement me")
}

func (this *HousNo1) Repair(code, cacheDate, featureDate string, complete bool) {
	securityCode := exchange.CorrectSecurityCode(code)
	tradeDate := exchange.FixTradeDate(featureDate)
	klines := base.CheckoutKLines(securityCode, tradeDate)
	if len(klines) < cache.KLineMin {
		return
	}
	df := pandas.LoadStructs(klines)
	var (
		//OPEN  = df.ColAsNDArray("open")
		CLOSE = df.ColAsNDArray("close")
		//HIGH  = df.ColAsNDArray("high")
		//LOW   = df.ColAsNDArray("low")
		//VOL   = df.ColAsNDArray("volume")
		//AMOUNT = df.ColAsNDArray("amount")
	)
	offset := 1
	if complete {
		offset = 0
	}
	ma5 := MA(CLOSE, 5-offset)
	this.MA5 = utils.SeriesIndexOf(ma5, -1)
	ma10 := MA(CLOSE, 10-offset)
	this.MA10 = utils.SeriesIndexOf(ma10, -1)
	ma20 := MA(CLOSE, 20-offset)
	this.MA20 = utils.SeriesIndexOf(ma20, -1)
	_ = df
}

func (this *HousNo1) Increase(snapshot QuoteSnapshot) Feature {
	tmp := HousNo1{}

	tmp.MA5 = (this.MA5*4 + snapshot.Price) / 5
	tmp.MA10 = (this.MA10*9 + snapshot.Price) / 10
	tmp.MA20 = (this.MA20*19 + snapshot.Price) / 20
	return &tmp
}

func (this *HousNo1) ValidateSample() error {
	if this.MA20 > 0 {
		return nil
	}
	return ErrInvalidFeatureSample
}
