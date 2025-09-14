package indicators

import (
	"github.com/quant1x/num"
	"github.com/quant1x/pandas"
	. "github.com/quant1x/pandas/formula"
)

// Platform 平台整理
func Platform(df pandas.DataFrame) pandas.DataFrame {
	var (
		OPEN  = df.ColAsNDArray("open")  // 开盘价
		CLOSE = df.ColAsNDArray("close") // 收盘价
		//HIGH    = df.ColAsNDArray("high")   // 最高价
		//LOW     = df.ColAsNDArray("low")    // 最低价
		VOL     = df.ColAsNDArray("volume") // 成交量
		BV      = df.ColAsNDArray("bv")     // 买入量
		SV      = df.ColAsNDArray("sv")     // 卖出量
		BA      = df.ColAsNDArray("ba")     // 买入金额
		SA      = df.ColAsNDArray("sa")     // 卖出金额
		DATALEN = df.Nrow()                 // 数据长度
	)
	//{T02: 平台, V1.0.0 2023-03-08}
	//BL1:=VOL/REF(VOL,1);
	LB := VOL.Div(REF(VOL, 1))
	//STH:=MAX(OPEN,CLOSE);
	STH := MAX(OPEN, CLOSE)
	//STL:=MIN(OPEN,CLOSE);
	STL := MIN(OPEN, CLOSE)
	//BLN1:=BARSLAST(BL1>=2.00),NODRAW;
	BLN := BARSLAST(LB.Gte(2.00))
	//倍量周期:BLN1,NODRAW;
	//{为HHV修复BLN1的值,需要+1}
	//BLN:=IFF(BLN1>=0,BLN1+1,BLN1);
	bln1 := IFF(BLN.Gte(0.00), BLN.Add(1), BLN)
	//倍量H:REF(STH,BLN);
	BLH := REF(STH, bln1)
	//倍量L:REF(STL,BLN);
	BLL := REF(STL, bln1)
	//BLVH:=HHV(VOL,BLN),NODRAW;
	BLVH := HHV(VOL, bln1)
	BLZC := LLV(STL, bln1)
	//BLVHN:=BARSLAST(VOL=BLVH),NODRAW;
	//BLVHN := BARSLAST(VOL.Eq(BLVH))
	//BLVL:=LLV(VOL,BLVHN),NODRAW;
	//BLVL := LLV(VOL, BLVHN)
	//BLVLN:=BARSLAST(VOL=BLVL),NODRAW;
	//BLVLN := BARSLAST(VOL.Eq(BLVL))
	//
	//SL:VOL/BLVH,NODRAW;
	SL := VOL.Div(BLVH)
	//SL21:=BARSLASTCOUNT(SL<=0.50),NODRAW;
	//SL20 := BARSSINCEN(SL.Lte(0.50), BLVHN)
	SL21 := BARSLASTCOUNT(SL.Lte(0.50))
	//{为REF修复SL21,需要-1}
	//SL2:=IFF(SL21>0,SL21-1,DRAWNULL),NODRAW;
	nan := num.Repeat(num.NaN(), DATALEN)
	null := pandas.ToSeries[num.DType](nan...)
	SL2 := IFF(SL21.Gt(0), SL21.Sub(1), null)
	//SL2 := SL21.Sub(1)
	//SL2 = SL21
	//SL3 := REF(VOL, SL2)
	SL3 := REF(VOL, SL2)
	//SLN:=BARSLAST(SL3=VOL);
	SLN := BARSLAST(SL3.Eq(VOL))
	//SLN := SL21
	//
	//缩量周期:SLN,NODRAW,COLORYELLOW;
	//缩量H:REF(STH,SLN+1),DOTLINE,COLORLIRED;
	SLH := REF(STH, SLN.Add(1))
	//缩量L:REF(STL,SLN),DOTLINE,COLORLIGREEN;
	SLL := REF(STL, SLN)

	oBlN := pandas.NewSeriesWithType(pandas.SERIES_TYPE_DTYPE, "BLN", BLN)
	oBlH := pandas.NewSeriesWithType(pandas.SERIES_TYPE_DTYPE, "BLH", BLH)
	oBlL := pandas.NewSeriesWithType(pandas.SERIES_TYPE_DTYPE, "BLL", BLL)
	oBlzc := pandas.NewSeriesWithType(pandas.SERIES_TYPE_DTYPE, "BLZC", BLZC)
	oSL := pandas.NewSeriesWithType(pandas.SERIES_TYPE_DTYPE, "SL", SL)
	oBLVH := pandas.NewSeriesWithType(pandas.SERIES_TYPE_DTYPE, "BLVH", BLVH)
	//oBLVHN := pandas.NewSeriesWithType(pandas.SERIES_TYPE_DTYPE, "BLVHN", BLVHN)
	oSlN := pandas.NewSeriesWithType(pandas.SERIES_TYPE_DTYPE, "SLN", SLN)
	//oSlNx := pandas.NewSeriesWithType(pandas.SERIES_TYPE_DTYPE, "SLNx", SL20)
	oSlH := pandas.NewSeriesWithType(pandas.SERIES_TYPE_DTYPE, "SLH", SLH)
	oSlL := pandas.NewSeriesWithType(pandas.SERIES_TYPE_DTYPE, "SLL", SLL)

	//B:倍量H<缩量H AND 倍量支撑=REF(倍量支撑,1) AND CROSS(缩量L,倍量支撑),COLORRED,NODRAW;
	x1 := BLH.Lte(SLH)
	//x2 := BLL.Lte(SLL)
	x2 := BLZC.Eq(REF(BLZC, 1))
	x3 := CROSS(SLL, BLZC)
	//x6 := REF(LOW, 1).Lt(BLZC)
	//x6 := REF(LOW, 1).Lt(BLZC)
	b1 := x1.And(x2).And(x3)
	oB1 := pandas.NewSeriesWithType(pandas.SERIES_TYPE_BOOL, "B1", b1)

	//B1:CROSS(CLOSE,缩量L),NODRAW;
	b2 := CROSS(CLOSE, SLL)
	oB2 := pandas.NewSeriesWithType(pandas.SERIES_TYPE_BOOL, "B2", b2)
	//B2:CROSS(CLOSE,倍量H) AND CROSS(CLOSE,倍量L),NODRAW;
	b3 := CROSS(CLOSE, BLH).And(CROSS(CLOSE, BLL))
	oB3 := pandas.NewSeriesWithType(pandas.SERIES_TYPE_BOOL, "B3", b3)
	mb := BA.Div(BV).Div(100)
	ms := SA.Div(SV).Div(100)
	oMb := pandas.NewSeriesWithType(pandas.SERIES_TYPE_DTYPE, "mb", mb)
	oMs := pandas.NewSeriesWithType(pandas.SERIES_TYPE_DTYPE, "ms", ms)
	df = df.Join(oBLVH, oBlzc, oBlN, oBlH, oBlL, oSlN, oSlH, oSlL, oSL, oB1, oB2, oB3, oMb, oMs)
	return df
}
