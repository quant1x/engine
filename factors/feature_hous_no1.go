package factors

import (
	"context"
	"gitee.com/quant1x/engine/datasets/base"
	"gitee.com/quant1x/gotdx/proto"
	"gitee.com/quant1x/gotdx/quotes"
	"gitee.com/quant1x/gotdx/trading"
	"gitee.com/quant1x/gox/api"
	"gitee.com/quant1x/pandas"
	. "gitee.com/quant1x/pandas/formula"
)

type HousNo1 struct {
	featureManifest `dataframe:"-"`
	MA5             float64 `dataframe:"ma5"`
	MA10            float64 `dataframe:"ma10"`
	MA20            float64 `dataframe:"ma20"`
}

func (f *HousNo1) Init(ctx context.Context, date string) error {
	//TODO implement me
	panic("implement me")
}

func (f *HousNo1) Factory(date string, code string) Feature {
	//TODO implement me
	panic("implement me")
}

//func (f *HousNo1) Kind() cache.Kind {
//	return FeatureHousNo1
//}
//
//func (f *HousNo1) Key() string {
//	return mapFeatures[f.Kind()].Key()
//}
//
//func (f *HousNo1) Name() string {
//	return mapFeatures[f.Kind()].Name()
//}
//
//func (f *HousNo1) Owner() string {
//	return mapFeatures[f.Kind()].Owner()
//}
//
//func (f *HousNo1) Usage() string {
//	return mapFeatures[f.Kind()].Name()
//}

//func (f *HousNo1) Init(ctx context.Context, date string) error {
//	return nil
//}
//
//func (f *HousNo1) GetDate() string {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (f *HousNo1) GetSecurityCode() string {
//	//TODO implement me
//	panic("implement me")
//}

func (f *HousNo1) FromHistory(history History) Feature {
	no1 := history.Payloads.No1
	_ = api.Copy(f, &no1)
	return f
}

func (f *HousNo1) Update(code, cacheDate, featureDate string, complete bool) {
	//TODO implement me
	panic("implement me")
}

func (f *HousNo1) Repair(code, cacheDate, featureDate string, complete bool) {
	securityCode := proto.CorrectSecurityCode(code)
	tradeDate := trading.FixTradeDate(featureDate)
	klines := base.CheckoutKLines(securityCode, tradeDate)
	if len(klines) == 0 {
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
	f.MA5 = SeriesIndexOf(ma5, -1)
	ma10 := MA(CLOSE, 10-offset)
	f.MA10 = SeriesIndexOf(ma10, -1)
	ma20 := MA(CLOSE, 20-offset)
	f.MA20 = SeriesIndexOf(ma20, -1)
	_ = df
}

func (f *HousNo1) Increase(snapshot quotes.Snapshot) Feature {
	tmp := HousNo1{}

	tmp.MA5 = (f.MA5*4 + snapshot.Price) / 5
	tmp.MA10 = (f.MA10*9 + snapshot.Price) / 10
	tmp.MA20 = (f.MA20*19 + snapshot.Price) / 20
	return &tmp
}
