package factors

import (
	"gitee.com/quant1x/data/exchange"
	"gitee.com/quant1x/engine/utils"
	"gitee.com/quant1x/num"
	"gitee.com/quant1x/pandas"
	. "gitee.com/quant1x/pandas/formula"
)

// QuShiFanZhuan 趋势反转
type QuShiFanZhuan struct {
	QSFZ bool    // 反转信号
	CP   float64 // 股价涨幅
	CV   float64 // 成交量涨幅
	VP   float64 // 价量比
	VP3  float64 // 3日价量比
	VP5  float64 // 5日价量比
}

// 趋势反转
func computeQuShiFanZhuan(date string, OPEN, CLOSE, HIGH, LOW, VOL pandas.Series) *QuShiFanZhuan {
	CURRBARSCOUNT := utils.IndexReverse(OPEN)
	// {趋势反转, V1.0.7, 2023-09-15}
	// MV5:=MA(VOL,5);
	MV5 := MA(VOL, 5)
	// LB0:VOL/REF(MV5,1),NODRAW;
	R1MV5 := REF(MV5, 1)
	LB0 := VOL.Div(R1MV5)
	// FIX:=IFF(CURRBARSCOUNT=1,FROMOPEN/TOTALFZNUM,1);
	minutes := exchange.Minutes(date)
	FIX := IFF(CURRBARSCOUNT.Eq(1), float64(minutes)/float64(exchange.CN_DEFAULT_TOTALFZNUM), 1.00)
	// LB:LB0/FIX,NODRAW;
	LB := LB0.Div(FIX)
	// NVOL:LB*REF(MV5,1),NODRAW;
	NVOL := R1MV5.Mul(LB)
	// CVOL:VOL,NODRAW;
	// XVOL:=NVOL;
	XVOL := NVOL
	// CVX:VOL/REF(VOL,1),NODRAW;
	// QSCV:XVOL/REF(VOL,1),NODRAW;
	//cv := VOL.Div(REF(VOL, 1))
	cv := XVOL.Div(REF(VOL, 1))
	// QSCP:(CLOSE/REF(CLOSE,1)-1)*100;
	cp := CLOSE.Div(REF(CLOSE, 1)).Sub(1.00).Mul(100)
	//cp := CLOSE.Div(REF(CLOSE, 1))
	//cp = cp.Sub(1)
	//fmt.Println(cp)
	// QSVP:QSCP/QSCV;
	vp := cp.Div(cv)
	// QSVP3:MA(QSVP,3);
	vp3 := MA(vp, 3)
	// QSVP5:MA(QSVP,5);
	vp5 := MA(vp, 5)
	// VP20:=MA(QSVP,20);
	// B:CROSS(QSVP,QSVP3),NODRAW;
	// DRAWICON(B,CLOSE,1);
	vpBuy := CROSS(vp, vp3)
	fz := num.AnyToBool(vpBuy.IndexOf(-1))
	qsfz := QuShiFanZhuan{
		QSFZ: fz,
		CV:   utils.Float64IndexOf(cv, -1),
		CP:   utils.Float64IndexOf(cp, -1),
		VP:   utils.Float64IndexOf(vp, -1),
		VP3:  utils.Float64IndexOf(vp3, -1),
		VP5:  utils.Float64IndexOf(vp5, -1),
	}
	return &qsfz
}
