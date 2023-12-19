package features

import (
	"gitee.com/quant1x/engine/cache"
	"gitee.com/quant1x/engine/datasets/base"
	"gitee.com/quant1x/engine/market"
	"gitee.com/quant1x/engine/models"
	"gitee.com/quant1x/engine/smart"
	"gitee.com/quant1x/engine/utils"
	"gitee.com/quant1x/gotdx/proto"
	"gitee.com/quant1x/gotdx/securities"
	"gitee.com/quant1x/gotdx/trading"
	"gitee.com/quant1x/gox/num"
	"gitee.com/quant1x/pandas"
	. "gitee.com/quant1x/pandas/formula"
	"gitee.com/quant1x/pandas/stat"
)

// ExchangeKLine K线特征
type ExchangeKLine struct {
	Date                string  // 日期
	Code                string  // 证券代码
	Shape               uint64  // K线形态
	MV3                 float64 // 3日均量
	MA3                 float64 // 3日均价
	MV5                 float64 // 5日均量
	MA5                 float64 // 5日均价
	MA10                float64 // 10日均价
	MA20                float64 // 20日均价
	VolumeRatio         float64 // 成交量比
	TurnoverRate        float64 // 换手率
	AmplitudeRatio      float64 // 振幅
	AveragePrice        float64 // 均价线
	Change5             float64 // 5日涨幅
	Change10            float64 // 10日涨幅
	InitialPrice        float64 // 短线底部(Short-Term Bottom),股价最近一次上串5日均线
	ShortIntensity      float64 // 短线强度
	ShortIntensityDiff  float64 // 短线强度增幅
	MediumIntensity     float64 // 中线强度
	MediumIntensityDiff float64 // 中线强度增幅
	Vix                 float64 // 波动率
	Sentiment           float64 // 情绪值
	Consistent          int     // 情绪一致
}

// NewExchangeKLine 构建制定日期的K线数据
func NewExchangeKLine(code, date string) *ExchangeKLine {
	securityCode := proto.CorrectSecurityCode(code)
	klines := base.CheckoutKLines(securityCode, date)
	if len(klines) < models.KLineMin {
		return nil
	}
	digits := 2
	securityInfo, ok := securities.CheckoutSecurityInfo(securityCode)
	if ok {
		digits = int(securityInfo.DecimalPoint)
	}
	ek := &ExchangeKLine{
		Date: date,
		Code: securityCode,
	}
	df := pandas.LoadStructs(klines)
	var (
		OPEN   = df.ColAsNDArray("open")
		CLOSE  = df.ColAsNDArray("close")
		HIGH   = df.ColAsNDArray("high")
		LOW    = df.ColAsNDArray("low")
		VOL    = df.ColAsNDArray("volume")
		AMOUNT = df.ColAsNDArray("amount")
		UP     = df.ColAsNDArray("up")
		DOWN   = df.ColAsNDArray("down")
	)
	// 2. 计算3日均量
	mv3 := MA(VOL, 3)
	ma3 := MA(CLOSE, 3)
	ek.MV3 = utils.SeriesIndexOf(mv3, -1)
	ek.MA3 = utils.SeriesIndexOf(ma3, -1)
	// 2. 计算5日分钟均量
	mv5 := MA(VOL, 5) // 5日均量
	mv5m := mv5.Div(trading.CN_DEFAULT_TOTALFZNUM)
	ma5Volume := utils.SeriesIndexOf(mv5m, -1)
	ek.MV5 = num.Decimal(ma5Volume, digits)
	ma5 := MA(CLOSE, 5)
	ek.MA5 = utils.SeriesIndexOf(ma5, -1)
	ma10 := MA(CLOSE, 10)
	ek.MA10 = utils.SeriesIndexOf(ma10, -1)
	ma20 := MA(CLOSE, 20)
	ek.MA20 = utils.SeriesIndexOf(ma20, -1)
	// 3. 隔日成交量放大比例
	vr := VOL.Div(REF(VOL, 1))
	volumeChangeRate := utils.SeriesIndexOf(vr, -1)
	ek.VolumeRatio = num.Decimal(volumeChangeRate, digits)
	// 换手率
	f10 := smart.GetL5F10(securityCode)
	if f10 != nil {
		turnoverRate := VOL.Div(f10.Capital).Mul(100.00)
		ek.TurnoverRate = num.Decimal(utils.SeriesIndexOf(turnoverRate, -1))
	}
	// 振幅, 这里只比对最高价和最低价的幅度, 不参考前一天的收盘价
	ar := HIGH.Div(LOW).Sub(1.00).Mul(100.00)
	ek.AmplitudeRatio = utils.SeriesIndexOf(ar, -1)

	// 4. 当日K线图形概要
	shape := KLineShape(df, securityCode)
	ek.Shape = shape
	// 均价线
	averagePrice := AMOUNT.Div(VOL)
	ap := utils.SeriesIndexOf(averagePrice, -1)
	ek.AveragePrice = num.Decimal(ap, digits)
	// 5. 计算阶段涨幅
	chg5 := CLOSE.Div(REF(CLOSE, 5)).Sub(1.00).Mul(100)
	chg10 := CLOSE.Div(REF(CLOSE, 10)).Sub(1.00).Mul(100)
	change5 := utils.SeriesIndexOf(chg5, -1)
	change10 := utils.SeriesIndexOf(chg10, -1)
	ek.Change5 = num.Decimal(change5, digits)
	ek.Change10 = num.Decimal(change10, digits)

	// 6. 多空强度
	//QDMA5:=MA(CLOSE,5);
	//ma5 := MA(CLOSE, 5)
	//QDMA10:=MA(CLOSE,10);
	//ma10 := MA(CLOSE, 10)
	//QDMA20:=MA(CLOSE,20);
	//ma20 := MA(CLOSE, 20)
	//QDV1:=QDMA5*100/QDMA20-100;
	//qdv1 := ma5.Div(ma20).Sub(1.00).Mul(100)
	qdv1 := utils.SeriesChangeRate(ma20, ma5)
	//QDV2:=QDMA10*100/QDMA20-100;
	qdv2 := utils.SeriesChangeRate(ma20, ma10)
	//QD1:=ABS(QDV1-REF(QDV1,1));
	//qd1 := ABS(qdv1.Sub(REF(qdv1, 1)))
	qd1 := qdv1.Sub(REF(qdv1, 1))
	//QD2:=ABS(QDV2-REF(QDV2,1));
	//qd2 := ABS(qdv2.Sub(REF(qdv2, 1)))
	qd2 := qdv2.Sub(REF(qdv2, 1))
	//QDCD:=CLOSE*100/REF(CLOSE,1)-100;
	//超短线:QDCD-REF(QDCD,1),NODRAW;
	//短线强度:QD1-REF(QD1,1),NODRAW;
	shortIntensity := qd1
	//中线强度:QD2-REF(QD2,1),NODRAW;
	mediumIntensity := qd2
	ek.ShortIntensity = utils.SeriesIndexOf(shortIntensity, -1)
	shortIntensityDiff := qd1.Sub(REF(qd1, 1))
	ek.ShortIntensityDiff = utils.SeriesIndexOf(shortIntensityDiff, -1)
	ek.MediumIntensity = utils.SeriesIndexOf(mediumIntensity, -1)
	mediumIntensityDiff := qd2.Sub(REF(qd2, 1))
	ek.MediumIntensityDiff = utils.SeriesIndexOf(mediumIntensityDiff, -1)

	// 7. 波动率
	N := 3
	M := 21
	//TYPICAL_PRICE:=(OPEN+CLOSE+HIGH+LOW)/4;
	TYPICAL_PRICE := OPEN.Add(CLOSE).Add(HIGH).Add(LOW).Div(4.00)
	//METHOD:=EMA(TYPICAL_PRICE,N);
	METHOD := EMA(TYPICAL_PRICE, N)
	//APPLY_TO:=TYPICAL_PRICE;
	APPLY_TO := TYPICAL_PRICE
	//MYSTDDEV:=SQRT(SUM(POW(APPLY_TO-METHOD,2),N)/N);
	var1 := APPLY_TO.Sub(METHOD)
	var2 := var1.Mul(var1)
	var3 := SUM(var2, N)
	var4 := SQRT(var3.Div(N))
	MYSTDDEV := stat.NewSeries[stat.DType](var4...)
	//MAXB:=HHV(MYSTDDEV,M);
	MAXB := HHV(MYSTDDEV, M)
	//MINB:=LLV(MYSTDDEV,M);
	MINB := LLV(MYSTDDEV, M)
	//VIX:100*(MYSTDDEV-MINB)/(MAXB-MINB);
	//10,DOTLINE,COLORGREEN;
	//50,DOTLINE,COLORYELLOW;
	vix := MYSTDDEV.Sub(MINB).Div(MAXB.Sub(MINB)).Mul(100)
	ek.Vix = utils.SeriesIndexOf(vix, -1)

	// 8. 短线底部(Short-Term Bottom),股价最近一次上串5日均线
	closeCrossMa5 := CROSS(CLOSE, ma5)
	crossPeriod := BARSLAST(closeCrossMa5)
	crossPrice := REF(OPEN, crossPeriod)
	ek.InitialPrice = num.Decimal(utils.SeriesIndexOf(crossPrice, -1))

	// 情绪值 情绪一致
	up := utils.SeriesIndexOf(UP, -1)
	down := utils.SeriesIndexOf(DOWN, -1)
	ek.Sentiment, ek.Consistent = market.SecuritySentiment(up, down)

	return ek
}

func (k *ExchangeKLine) Kind() cache.Kind {
	return FeatureKLineShap
}

func (k *ExchangeKLine) Name() string {
	return "K线形态"
}
