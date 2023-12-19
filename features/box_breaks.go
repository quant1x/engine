package features

import (
	"gitee.com/quant1x/engine/cache"
	"gitee.com/quant1x/engine/datasets/base"
	"gitee.com/quant1x/engine/models"
	"gitee.com/quant1x/engine/utils"
	"gitee.com/quant1x/gotdx/proto"
	"gitee.com/quant1x/gotdx/trading"
	"gitee.com/quant1x/pandas"
	. "gitee.com/quant1x/pandas/formula"
	"gitee.com/quant1x/pandas/stat"
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
	QSFZ           bool    // 趋势反转
	QSCP           float64 // 趋势反转
	QSCV           float64 // 趋势反转
	QSVP           float64 // 趋势反转
	QSVP3          float64 // 趋势反转
	QSVP5          float64 // 趋势反转
	DkCol          float64 // dkqs: 能量柱, 通达信分时指标DkCol
	DkD            float64 // dkqs: 多头力量
	DkK            float64 // dkqs: 空头力量
	DkB            bool    // dkqs: 买入
	DkS            bool    // dkqs: 卖出
	DxDivergence   float64 // madx: 综合发散度评估值
	DxDm0          float64 // madx: 超短线均线发散度
	DxDm1          float64 // madx:   短线均线发散度
	DxDm2          float64 // madx:   中线均线发散度
	DxB            bool    // madx: 买入
}

// NewKLineBox 构建有效突破数据
func NewKLineBox(code, date string) *KLineBox {
	securityCode := proto.CorrectSecurityCode(code)
	tradeDate := trading.FixTradeDate(date)
	klines := base.CheckoutKLines(securityCode, tradeDate)
	if len(klines) < models.KLineMin {
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
	doublePeriod := int(stat.AnyToInt64(v1))
	//BLH:=IFF(BLN=0,BOXH,REF(BOXH,BLN)),DOTLINE;
	BLH := IFF(BLN.Eq(0), BOXH, REF(BOXH, BLN))
	//BLL:=IFF(BLN=0,BOXL,REF(BOXL,BLN)),DOTLINE;
	BLL := IFF(BLN.Eq(0), BOXL, REF(BOXL, BLN))
	//倍量H:IFF(BLN=0,REF(BLH,1),BLH),DOTLINE;
	dvH := IFF(BLN.Eq(0), REF(BLH, 1), BLH)
	doubleHigh := utils.SeriesIndexOf(dvH, -1)
	//倍量L:IFF(BLN=0,REF(BLL,1),BLL),DOTLINE;
	dvL := IFF(BLN.Eq(0), REF(BLL, 1), BLL)
	doubleLow := utils.SeriesIndexOf(dvL, -1)
	//倍量压力:IFF(BLN=0,HIGH,HHV(HIGH,BLN)),DOTLINE;
	//
	//MA3:MA(CLOSE,3),COLORYELLOW;
	MA3 := MA(CLOSE, 3)
	//
	//{绘制买入信号}
	//B:CROSS(CLOSE,倍量H),COLORRED;
	B := CROSS(CLOSE, dvH)
	v2 := B.IndexOf(-1)
	buy := stat.AnyToBool(v2)
	//B1:CLOSE>倍量H AND REF(CLOSE,1)<MA3,COLORRED;
	//DRAWICON(B,LOW*ICON_B_RATIO,1);
	//
	//{绘制卖出信号}
	//S:CROSS(MA3,CLOSE);
	S := CROSS(MA3, CLOSE)
	v3 := S.IndexOf(-1)
	sell := stat.AnyToBool(v3)
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
		TendencyPeriod: int(stat.AnyToInt64(tendencyPeriod)),
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
	return &box
}

func (k *KLineBox) Kind() cache.Kind {
	return FeatureBreaksThroughBox
}

func (k *KLineBox) Name() string {
	return "突破平台"
}
