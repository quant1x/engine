package factors

import (
	"gitee.com/quant1x/data/exchange"
	"gitee.com/quant1x/engine/cache"
	"gitee.com/quant1x/engine/datasource/base"
	"gitee.com/quant1x/engine/indicators"
	"gitee.com/quant1x/engine/utils"
	"gitee.com/quant1x/num"
	"gitee.com/quant1x/pandas"
	. "gitee.com/quant1x/pandas/formula"
)

// KLineBox 有效突破(BreaksThrough)平台
type KLineBox struct {
	Code           string  // 证券代码
	Date           string  // 数据日期
	DoubletPeriod  int     // box: 倍量周期
	DoubleHigh     float64 // box: 倍量最高
	DoubleLow      float64 // box: 倍量最低
	Buy            bool    // box: 买入
	HalfPeriod     int     // box: 半量周期
	HalfHigh       float64 // box: 半量最高
	HalfLow        float64 // box: 半量最低
	Sell           bool    // box: 卖出
	TendencyPeriod int     // box: 趋势周期
	QSFZ           bool    // qsfz: 信号
	QSCP           float64 // qsfz: cp
	QSCV           float64 // qsfz: cv
	QSVP           float64 // qsfz: vp
	QSVP3          float64 // qsfz: vp3
	QSVP5          float64 // qsfz: vp5
	DkCol          float64 // dkqs: 能量柱, 通达信分时指标DkCol
	DkD            float64 // dkqs: 多头力量
	DkK            float64 // dkqs: 空头力量
	DkB            bool    // dkqs: 买入
	DkS            bool    // dkqs: 卖出
	DxDivergence   float64 // madx: 综合发散度评估值
	DxDm0          float64 // madx: 均线发散度-超短线
	DxDm1          float64 // madx: 均线发散度-短线
	DxDm2          float64 // madx: 均线发散度-中线
	DxB            bool    // madx: 买入
	DxBN           int     // madx: 连续DxB信号周期数
	SarPos         int     // 坐标位置
	SarBull        bool    // 当前多空
	SarAf          float64 // 加速因子(Acceleration Factor)
	SarEp          float64 // 极值点(Extreme Point)
	SarSar         float64 // SAR[Pos]
	SarHigh        float64 // pos周期最高价
	SarLow         float64 // pos周期最低价
	SarPeriod      int     // 周期数, 上涨趋势, 周期数大于0, 下跌趋势, 周期数小于0, 绝对值就是已过多少天
}

// NewKLineBox 构建有效突破数据
func NewKLineBox(code, date string) *KLineBox {
	securityCode := exchange.CorrectSecurityCode(code)
	tradeDate := exchange.FixTradeDate(date)
	klines := base.CheckoutKLines(securityCode, tradeDate)
	if len(klines) < cache.KLineMin {
		return nil
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
		VOL   = df.ColAsNDArray("volume")
		HIGH  = df.ColAsNDArray("high")
		LOW   = df.ColAsNDArray("low")
		//AMOUNT = df.ColAsNDArray("amount")
	)
	//{T05: 有效突破平台, V1.0.9 2023-07-22}
	//{倍量1, 以5日均量线的2倍计算}
	//ICON_B_RATIO:=0.999;
	//ICON_S_RATIO:=1.029;
	//VRATIO:=2.00;
	VRATIO := 2.00
	//BL1:=VOL/REF(VOL,1);
	BL1 := VOL.Div(REF(VOL, 1))
	//BOXH:=MAX(OPEN,CLOSE);
	BOXH := MAX(OPEN, CLOSE)
	//BOXL:=MIN(OPEN,CLOSE);
	BOXL := MIN(OPEN, CLOSE)
	//BLN1:=BARSLAST(BL1>=VRATIO AND CLOSE>OPEN),NODRAW;
	BLN1 := BARSLAST(BL1.Gte(VRATIO).And(CLOSE.Gt(OPEN)))
	//倍量周期:BLN1,NODRAW;
	//{为HHV修复BLN1的值,需要+1}
	//BLN:=IFF(BLN1>=0,BLN1,BLN1),NODRAW;
	BLN := BLN1
	v1 := BLN.IndexOf(-1)
	doublePeriod := int(num.AnyToInt64(v1))
	//BLH:=IFF(BLN=0,BOXH,REF(BOXH,BLN)),DOTLINE;
	BLH := IFF(BLN.Eq(0), BOXH, REF(BOXH, BLN))
	//BLL:=IFF(BLN=0,BOXL,REF(BOXL,BLN)),DOTLINE;
	BLL := IFF(BLN.Eq(0), BOXL, REF(BOXL, BLN))
	//倍量H:IFF(BLN=0,REF(BLH,1),BLH),DOTLINE;
	dvH := IFF(BLN.Eq(0), REF(BLH, 1), BLH)
	doubleHigh := utils.Float64IndexOf(dvH, -1)
	//倍量L:IFF(BLN=0,REF(BLL,1),BLL),DOTLINE;
	dvL := IFF(BLN.Eq(0), REF(BLL, 1), BLL)
	doubleLow := utils.Float64IndexOf(dvL, -1)
	//倍量压力:IFF(BLN=0,HIGH,HHV(HIGH,BLN)),DOTLINE;
	//
	//MA3:MA(CLOSE,3),COLORYELLOW;
	MA3 := MA(CLOSE, 3)
	//
	//{绘制买入信号}
	//B:CROSS(CLOSE,倍量H),COLORRED;
	B := CROSS(CLOSE, dvH)
	v2 := B.IndexOf(-1)
	buy := num.AnyToBool(v2)
	//B1:CLOSE>倍量H AND REF(CLOSE,1)<MA3,COLORRED;
	//DRAWICON(B,LOW*ICON_B_RATIO,1);
	//
	//{绘制卖出信号}
	//S:CROSS(MA3,CLOSE);
	S := CROSS(MA3, CLOSE)
	v3 := S.IndexOf(-1)
	sell := num.AnyToBool(v3)
	//DRAWICON(S,HIGH*ICON_S_RATIO,2);
	//df = df.Join(BLN, dvH, dvL, B, S)
	//fmt.Println(df)

	// 多空趋势周期
	//MA5:=MA(CLOSE,5)
	MA5 := MA(CLOSE, 5)
	//D:=MA5>REF(MA5,1) AND CLOSE>=MA5;
	D := MA5.Gt(REF(MA5, 1)).And(CLOSE.Gte(MA5))
	//K:=MA5<=REF(MA5,1) OR CLOSE<MA5
	K := MA5.Lte(REF(MA5, 1)).Or(CLOSE.Lt(MA5))
	//FD:=BARSLASTCOUNT(D);
	FD := BARSLASTCOUNT(D)
	//FK:=BARSLASTCOUNT(K);
	FK := BARSLASTCOUNT(K)
	//多空:IFF(FD>0,FD,-1*FK),NODRAW
	DK := IFF(FD.Gt(0), FD, FK.Mul(-1))
	tendencyPeriod := DK.IndexOf(-1)

	box := KLineBox{
		Code:           securityCode,
		Date:           tradeDate,
		DoubletPeriod:  doublePeriod,
		DoubleHigh:     doubleHigh,
		DoubleLow:      doubleLow,
		Buy:            buy,
		Sell:           sell,
		TendencyPeriod: int(num.AnyToInt64(tendencyPeriod)),
	}
	// 附加 趋势反转
	qsfz := computeQuShiFanZhuan(tradeDate, OPEN, CLOSE, HIGH, LOW, VOL)
	box.QSFZ = qsfz.QSFZ
	box.QSCP = qsfz.CP
	box.QSCV = qsfz.CV
	box.QSVP = qsfz.VP
	box.QSVP3 = qsfz.VP3
	box.QSVP5 = qsfz.VP5

	// 附加 多空趋势
	dkqs := computeDuoKongQuShi(OPEN, CLOSE, HIGH, LOW)
	box.DkCol = dkqs.Col
	box.DkD = dkqs.D
	box.DkK = dkqs.K0
	box.DkB = dkqs.B
	box.DkS = dkqs.S
	// 附加 均线动向
	madx := computeJuXianDongXiang(OPEN, CLOSE, HIGH, LOW)
	box.DxDivergence = madx.Diverging
	box.DxDm0 = madx.Dm0
	box.DxDm1 = madx.Dm1
	box.DxDm2 = madx.Dm2
	box.DxB = madx.B
	box.DxBN = madx.BN

	// SAR
	sars := indicators.SAR(HIGH.Float64s(), LOW.Float64s())
	if len(sars) > 0 {
		sar := sars[len(sars)-1]
		box.SarPos = sar.Pos
		box.SarBull = sar.Bull
		box.SarAf = sar.Af
		box.SarEp = sar.Ep
		box.SarSar = sar.Sar
		box.SarHigh = sar.High
		box.SarLow = sar.Low
		box.SarPeriod = sar.Period
	}
	return &box
}

func (k *KLineBox) Kind() cache.Kind {
	return FeatureBreaksThroughBox
}

func (k *KLineBox) Name() string {
	return "突破平台"
}
