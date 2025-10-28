package factors

import (
	"context"

	"gitee.com/quant1x/data/exchange"
	"gitee.com/quant1x/engine/cache"
	"gitee.com/quant1x/engine/datasource/base"
	"gitee.com/quant1x/engine/utils"
	"gitee.com/quant1x/num"
	"gitee.com/quant1x/pandas"
	. "gitee.com/quant1x/pandas/formula"
)

const (
	cacheL5KeyInvestmentSentimentMaster = "ism"
	S8PeriodOfLimitUp                   = 30
)

// InvestmentSentimentMaster 情绪大师
type InvestmentSentimentMaster struct {
	cache.DataSummary `dataframe:"-"`
	Date              string  `name:"日期" dataframe:"日期"`     // 数据日期
	Code              string  `name:"证券代码" dataframe:"证券代码"` // 证券代码
	PN                int     `name:"观察周期" dataframe:"pn"`
	ZDF               float64 `name:"涨跌幅" dataframe:"zdf"`
	R1CLOSE           float64 `name:"昨收盘" dataframe:"r1Close"`
	ZTJ               float64 `name:"涨停价" dataframe:"ztj"`
	CZT               bool    `name:"涨停" dataframe:"czt"`
	BN                int     `name:"板数" dataframe:"bn"`
	FZT               int     `name:"距离首次涨停" dataframe:"fzt"`
	TN                int     `name:"天数" dataframe:"tn"`
	TIAN              int     `name:"天" dataframe:"tian"`
	BAN               int     `name:"板" dataframe:"ban"`
	ZHANG             int     `name:"涨" dataframe:"zhang"`
	PING              int     `name:"平" dataframe:"ping"`
	DIE               int     `name:"跌" dataframe:"die"`
	OH                float64 `name:"周期内最高价" dataframe:"oh"`
	COH               bool    `name:"是否周期内最高价" dataframe:"coh"`
	OHN               int     `name:"最高价周期" dataframe:"ohn"`
	OL                float64 `name:"周期内最低价" dataframe:"ol"`
	COL               bool    `name:"是否周期内最低价" dataframe:"col"`
	OLN               int     `name:"最低价周期" dataframe:"oln"`
	OHV               float64 `name:"最高价量" dataframe:"ohv"`
	OHBL              float64 `name:"最高价量比" dataframe:"ohbl"`
	OLV               float64 `name:"最低价量" dataframe:"olv"`
	OLBL              float64 `name:"最低价量比" dataframe:"olbl"`
	UpdateTime        string  `name:"更新时间" dataframe:"update_time"` // 更新时间
	State             uint64  `name:"样本状态" dataframe:"样本状态"`        // 样本状态
}

func NewInvestmentSentimentMaster(date, code string) *InvestmentSentimentMaster {
	summary := __mapFeatures[FeatureInvestmentSentimentMaster]
	v := InvestmentSentimentMaster{
		DataSummary: summary,
		Date:        date,
		Code:        code,
	}
	return &v
}

func (this *InvestmentSentimentMaster) Factory(date string, code string) Feature {
	v := NewInvestmentSentimentMaster(date, code)
	return v
}

func (this *InvestmentSentimentMaster) GetDate() string {
	return this.Date
}

func (this *InvestmentSentimentMaster) GetSecurityCode() string {
	return this.Code
}

func (this *InvestmentSentimentMaster) Init(ctx context.Context, date string) error {
	_ = ctx
	_ = date
	return nil
}

func (this *InvestmentSentimentMaster) Update(code, cacheDate, featureDate string, whole bool) {
	securityCode := exchange.CorrectSecurityCode(code)
	this.Date = exchange.FixTradeDate(cacheDate)
	this.Code = securityCode
	tradeDate := exchange.FixTradeDate(featureDate)
	klines := base.CheckoutKLines(securityCode, tradeDate)
	if len(klines) < cache.KLineMin {
		return
	}
	df := pandas.LoadStructs(klines)
	var (
		DATE  = df.Col("date")
		OPEN  = df.ColAsNDArray("open")
		CLOSE = df.ColAsNDArray("close")
		HIGH  = df.ColAsNDArray("high")
		LOW   = df.ColAsNDArray("low")
		VOL   = df.ColAsNDArray("volume")
	)
	//PN:30,COLORWHITE,NODRAW;
	PN := S8PeriodOfLimitUp
	//R1CLOSE:=REF(CLOSE,1);
	R1CLOSE := REF(CLOSE, 1)
	//CST:=NOT(NAMELIKE('S') OR NAMELIKE('*S')) AND VOL>1;
	//ZDF:=IFF(INBLOCK('创业板'), 0.2, IFF(INBLOCK('科创板'),0.2, IFF(INBLOCK('ST板块'), 0.05, IFF(INBLOCK('北证A股'),0.3,0.1))));
	ZDF := exchange.MarketLimit(securityCode)
	//ZTJ:=ZTPRICE(R1CLOSE,ZDF);
	ZTJ := R1CLOSE.Mul(1.00+ZDF).Apply2(func(idx int, v any) any {
		f := v.(float64)
		return num.Decimal(f)
	}, true)
	//CZT:=CLOSE=ZTJ;
	CZT := CLOSE.Gte(ZTJ)
	//BN:COUNT(CZT,PN);
	BN := COUNT(CZT, PN)
	//FTZ:=BARSSINCEN(CZT,PN);
	FTZ := BARSSINCEN(CZT, PN)
	//TN:FTZ+1,COLORWHITE,NODRAW;
	TN := FTZ.Add(1)
	//天:TN,COLORYELLOW,NODRAW;
	TIAN := TN
	//板:BN,COLORRED,NODRAW;
	BAN := BN
	//CUP:=CLOSE>R1CLOSE;
	CUP := CLOSE.Gt(R1CLOSE)
	//CPING:=CLOSE=R1CLOSE;
	CPING := CLOSE.Eq(R1CLOSE)
	//CDOWN:=CLOSE<R1CLOSE;
	CDOWN := CLOSE.Lt(R1CLOSE)
	//涨:COUNT(CUP,TN),COLORRED,NODRAW;
	ZHANG := COUNT(CUP, TN)
	//平:COUNT(CPING,TN),COLORWHITE,NODRAW;
	PING := COUNT(CPING, TN)
	//跌:COUNT(CDOWN,TN),COLORGREEN,NODRAW;
	DIE := COUNT(CDOWN, TN)

	//{首板以来的最高价}
	//OH:HHV(HIGH,TN),COLORRED,DOTLINE;
	OH := HHV(HIGH, TN)
	//{首板以来最高价到现在的周期}
	//OHN:BARSLAST(HIGH=OH),COLORYELLOW,NODRAW;
	COH := HIGH.Eq(OH)
	OHN := BARSLAST(COH)
	//{首板以来的最低价}
	//OL:LLV(LOW,TN),COLORGREEN,DOTLINE;
	OL := LLV(LOW, TN)
	//{首板以来最低价到现在的周期}
	//OLN:BARSLAST(LOW=OL),COLORYELLOW,NODRAW;
	COL := LOW.Eq(OL)
	OLN := BARSLAST(COL)
	//{首板以来最高价当日的成交量}
	//OHV:REF(VOL,OHN),COLORWHITE,NODRAW;
	OHV := REF(VOL, OHN)
	//{今天成交量和最高价当日成交量的比值}
	//OHBL:VOL/OHV,COLORWHITE,NODRAW;
	OHBL := VOL.Div(OHV)
	//{首板以来最低价当日的成交量}
	//OLV:REF(VOL,OLN),COLORWHITE,NODRAW;
	OLV := REF(VOL, OLN)
	//{今天成交量和最低价当日成交量的比值}
	//OLBL:VOL/OLV,COLORWHITE,NODRAW;
	OLBL := VOL.Div(OLV)
	{
		// 特征数据采集
		this.PN = PN
		this.ZDF = ZDF
		this.R1CLOSE = utils.Float64IndexOf(CLOSE, -1)
		this.ZTJ = utils.Float64IndexOf(ZTJ, -1)
		this.CZT = utils.BoolIndexOf(CZT, -1)
		this.BN = utils.IntegerIndexOf(BN, -1)
		this.FZT = utils.IntegerIndexOf(FTZ, -1)
		this.TN = utils.IntegerIndexOf(TN, -1)
		this.TIAN = utils.IntegerIndexOf(TIAN, -1)
		this.BAN = utils.IntegerIndexOf(BAN, -1)
		this.ZHANG = utils.IntegerIndexOf(ZHANG, -1)
		this.PING = utils.IntegerIndexOf(PING, -1)
		this.DIE = utils.IntegerIndexOf(DIE, -1)
		this.OH = utils.Float64IndexOf(OH, -1)
		this.COH = utils.BoolIndexOf(COH, -1)
		this.OHN = utils.IntegerIndexOf(OHN, -1)
		this.OL = utils.Float64IndexOf(OL, -1)
		this.COL = utils.BoolIndexOf(COL, -1)
		this.OLN = utils.IntegerIndexOf(OLN, -1)
		this.OHV = utils.Float64IndexOf(OHV, -1)
		this.OHBL = utils.Float64IndexOf(OHBL, -1)
		this.OLV = utils.Float64IndexOf(OLV, -1)
		this.OLBL = utils.Float64IndexOf(OLBL, -1)
	}
	//{
	//	// 调试
	//	df = pandas.NewDataFrame(DATE, CLOSE, ZTJ, CZT, BN, FTZ, TN)
	//	fmt.Println(df)
	//}
	this.UpdateTime = GetTimestamp()
	this.State |= this.Kind()
	_ = DATE
	_ = OPEN
}

func (this *InvestmentSentimentMaster) Repair(securityCode, cacheDate, featureDate string, whole bool) {
	this.Update(securityCode, cacheDate, featureDate, whole)
}

func (this *InvestmentSentimentMaster) FromHistory(history History) Feature {
	_ = history
	return this
}

func (this *InvestmentSentimentMaster) Increase(snapshot QuoteSnapshot) Feature {
	_ = snapshot
	return this
}

func (this *InvestmentSentimentMaster) ValidateSample() error {
	if this.State > 0 {
		return nil
	}
	return ErrInvalidFeatureSample
}
