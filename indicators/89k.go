package indicators

import (
	"github.com/quant1x/num"
	"github.com/quant1x/pandas"
	. "github.com/quant1x/pandas/formula"
)

// F89K 89k测试
func F89K(df pandas.DataFrame, N int) pandas.DataFrame {
	var (
		//OPEN  = df.ColAsNDArray("open")
		CLOSE = df.ColAsNDArray("close")
		HIGH  = df.ColAsNDArray("high")
		LOW   = df.ColAsNDArray("low")
		VOL   = df.ColAsNDArray("volume")
		//DATE  = df.ColAsNDArray("date")
	)
	length := CLOSE.Len()
	N1 := N
	// {89K图形}
	//
	//MA89:=MA(CLOSE,N1);
	//MA89 := MA(CLOSE, N1)
	//{计算N1日内最低价}
	//L89:=LLV(LOW,N1);
	//L89 := LLV(LOW, N1)
	//{计算N1日内最高价}
	//H89:=HHV(HIGH,N1);
	H89 := HHV(HIGH, N1)
	//
	//{确定⑤}
	//K5:L89,NODRAW,COLORGREEN;
	//K5 := L89
	//{确定⑦}
	//N7:BARSLAST(H89=HIGH),NODRAW,COLORLIGRAY;
	x7 := EQ(H89, HIGH)
	N7 := BARSLAST(pandas.ToSeries[bool](x7...))
	//T7:=REF(HIGH,N7);
	//T7 := REF(HIGH, N7)
	//K7:T7,NODRAW,COLORRED;
	//K7 := T7
	//{确定⑧}
	//T8:=LLV(LOW,N7),NODRAW,COLORGREEN;
	T8 := LLV(LOW, N7)
	//N8:BARSLAST(T8=LOW AND N1>N7),NODRAW,COLORLIGRAY;
	x81 := EQ(T8, LOW)
	xn1 := num.Repeat[num.DType](num.DType(N1), length)
	x82 := CompareGt(xn1, N7)
	x8 := AND(x81, x82)
	N8 := BARSLAST(pandas.ToSeries[bool](x8...))
	//K8:T8,NODRAW,COLORGREEN;
	//K8 := T8
	//{确定⑨}
	//T9:=HHV(HIGH,N8);
	T9 := HHV(HIGH, N8)
	//N9X:=BARSLAST(T9=HIGH AND N7>N8);
	//x91 := EQ(T9, HIGH)
	x91 := T9.Eq(HIGH)
	//x92 := CompareGt(N7, N8)
	x92 := N7.Gt(N8)
	//x9 := AND(x91, x92)
	x9 := x91.And(x92)
	N9X := BARSLAST(x9)
	//N9:IFF(N9X=0,N9X+1,N9X),NODRAW,COLORLIGRAY;
	//x911 := EQ2(N9X, 0)
	x911 := N9X.Eq(0)
	//var n9x stat.Series
	//n9x = stat.NDArray[stat.DType](N9X)
	N9 := IFF(x911, N9X.Add(1), N9X)
	//K9:REF(HIGH,N9X),NODRAW,COLORRED;
	K9 := REF(HIGH, N9X)
	//{确定⑩}
	//K10:LLV(LOW,N9),NODRAW,COLORGREEN;
	K10 := LLV(LOW, N9)
	//N10:BARSLAST(K10=LOW AND N8>N9),NODRAW,COLORLIGRAY;
	//x101 := EQ(K10, LOW)
	x101 := K10.Eq(LOW)
	//x102 := CompareGt(N8, N9)
	x102 := N8.Gt(N9)
	//x10 := AND(x101, x102)
	x10 := x101.And(x102)
	N10 := BARSLAST(x10)
	//
	//{比对周期长度}
	//C_N:=5;
	//C_N := 5
	//{量价关系最低校对比率}
	//C_S:=0.191;
	C_S := 0.191
	//{涂画与股价的纵向比率}
	//C_PX:=0.002;
	//C_PX := 0.002
	//{真阳线}
	//{C_ISMALE:=(CLOSE > REF(CLOSE,1)) AND (CLOSE > OPEN);}
	//C_ISMALE:=CLOSE > REF(CLOSE,1);
	C_ISMALE := CLOSE.Gt(REF(CLOSE, 1))
	//{成交量较上一个周期放大}
	//C_VOL:= VOL>REF(VOL,1);
	C_VOL := VOL.Gt(REF(VOL, 1))
	//{成交量均线周期}
	//VOL_PERIOD:=5;
	VOL_PERIOD := 5
	//{成交量比例}
	//VOLSCALE:=1+C_S;
	VOLSCALE := 1 + C_S
	//{高股价或指数的计算方法, 比MAVOL5高出C_S/10且比前一日方法}
	//X_INDEX:=VOL>=MA(VOL, VOL_PERIOD)*(1 + C_S/10);
	bl := 1.00 + C_S/10.00
	x_ma0 := MA(VOL, VOL_PERIOD)
	x_ma := x_ma0.Mul(bl)
	X_INDEX := VOL.Gte(x_ma)
	//{一般股价的计算方法}
	//X_GENERAL:=VOL>=x_ma0*VOLSCALE;
	X_GENERAL := VOL.Gte(x_ma0.Mul(VOLSCALE))
	//{指数类或者高股价类的成交量不太可能像个股那样成倍放量, 这里做一个降级处理}
	//X:=IFF(CLOSE>=100, X_INDEX, X_GENERAL) AND C_ISMALE AND C_VOL;
	xx1 := IFF(CLOSE.Gte(100), X_INDEX, X_GENERAL)
	xx2 := C_ISMALE.And(C_VOL)
	X := xx1.And(xx2)
	//{确定X上一次确立成立是在哪一天}
	//DN:=BARSLAST(X);
	DN := BARSLAST(X)
	//{放量上攻作为一个小阶段的起点, 该起点K线的最低价作为止盈止损线}
	//ZS_VOL:=REF(VOL,DN);
	//ZS_VOL := REF(VOL, DN)
	//ZS_LOW:=REF(LOW,DN);
	ZS_LOW := REF(LOW, DN)
	//ZS:ZS_LOW,NODRAW,COLORLIGRAY;
	//ZS := ZS_VOL
	//N11:BARSLAST(CLOSE<ZS_LOW AND N9>N10),NODRAW,COLORLIGRAY;
	//N11 := BARSLAST(CLOSE.Lt(ZS_LOW).Eq(N9.Gt(N10)))
	//{股价突破止损线}
	//{C0:BARSSINCEN(CLOSE<ZS_LOW,DN),NODRAW;}
	//{C0:FINDHIGHBARS(CLOSE,1,DN-1,1),NODRAW;}
	//C00:=BARSLASTCOUNT(CLOSE<ZS_LOW),NODRAW;
	//C01:=BARSLASTCOUNT(CLOSE>=ZS_LOW),NODRAW;
	//C02:=IFF(C01>0,REF(C00,1),C01),NODRAW;
	//C03:=REF(CLOSE,C02-1),NODRAW;
	//C1:BARSLAST(CLOSE>ZS_LOW AND N9>N10),NODRAW,COLORLIGRAY;
	c11 := CLOSE.Gt(ZS_LOW)
	c12 := N9.Gt(N10)
	c1c := c11.And(c12)
	C1 := BARSLAST(c1c)
	//{股价跌破止损线}
	//C2:BARSLAST(CLOSE<ZS_LOW AND N9>N10),NODRAW,COLORLIGRAY;
	c21 := CLOSE.Lt(ZS_LOW)
	c22 := N9.Gt(N10)
	C2 := BARSLAST(c21.And(c22))
	//C3:N10>0,NODRAW;
	//C3 := CompareGt(N10, 0)
	C3 := N10.Gt(0)
	//B0:=C1=0 AND C2=1 AND C3;
	//B01 := EQ2(C1, 0)
	B01 := C1.Eq(0)
	//B02 := EQ2(C2, 1)
	B02 := C2.Eq(1)
	//B03 := AND(B01, B02)
	B03 := B01.And(B02)
	//AND C3
	//B0 := AND(B03, C3)
	B0 := B03.And(C3)
	//B:50*B0,COLORYELLOW;
	OB := pandas.NewSeriesWithType(pandas.SERIES_TYPE_BOOL, "B", B0)
	OZS := pandas.NewSeriesWithType(pandas.SERIES_TYPE_DTYPE, "ZS", ZS_LOW)
	//ODT := pandas.NewSeries(stat.SERIES_TYPE_STRING, "date", DATE.Values().([]string))
	//return pandas.NewDataFrame(ODT, ON9, OK9, ON10, OK10, OB)
	ODN := pandas.NewSeriesWithType(pandas.SERIES_TYPE_INT64, "N", DN)
	_ = K9
	return pandas.NewDataFrame(df.Col("date"), df.Col("open"), df.Col("close"), df.Col("high"), df.Col("low"), OZS, OB, ODN)
}
