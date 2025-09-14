package indicators

import (
	"github.com/quant1x/pandas"
	. "github.com/quant1x/pandas/formula"
)

// CDTD 抄底逃顶
func CDTD(df pandas.DataFrame) pandas.DataFrame {
	var (
		//OPEN  = df.ColAsNDArray("open")   // 开盘价
		CLOSE = df.ColAsNDArray("close") // 收盘价
		HIGH  = df.ColAsNDArray("high")  // 最高价
		LOW   = df.ColAsNDArray("low")   // 最低价
		//VOL   = df.ColAsNDArray("volume") // 成交量
		//DATALEN = df.Nrow()                 // 数据长度
	)
	N1 := 3
	N2 := 9
	N3 := 27
	N4 := 5
	HV3 := HHV(HIGH, N3)
	LV3 := LLV(LOW, N3)
	HHVLLV3 := HV3.Sub(LV3)
	CLLV3 := CLOSE.Sub(LV3)
	//RSV1:=(CLOSE-LLV(LOW,N2))/(HHV(HIGH,N2)-LLV(LOW,N2))*100;
	LV2 := LLV(LOW, N2)
	HV2 := HHV(HIGH, N2)
	rsv11 := CLOSE.Sub(LV2)
	rsv12 := HV2.Sub(LV2)
	RSV1 := rsv11.Div(rsv12).Mul(100)
	// RSV2:=CLLV3/HHVLLV3*100;
	RSV2 := CLLV3.Div(HHVLLV3).Mul(100)
	//RSV3:=SMA(RSV2,N4,1);
	RSV3 := SMA(RSV2, N4, 1)
	//WEN:=N1*RSV3-2*SMA(RSV3,N1,1);
	RSV4 := SMA(RSV3, N1, 1)
	WEN := RSV3.Mul(N1).Sub(RSV4.Mul(2))
	J1 := SMA(RSV1, N1, 1)
	J2 := SMA(J1, N1, 1)
	//W1 := SMA(RSV2, N1, 1)
	//W2 := SMA(W1, N1, 1)
	//强弱界线:49,DOTLINE,LINETHICK1,COLOR9966CC;
	//QIANGRUO := 49.00
	//顶:100, DOTLINE,COLORCCFF00;
	//DING := 100.00
	//底:0, DOTLINE, COLORRED;
	//DI := 0.00
	//趋势线:WEN,LINETHICK1,COLORFF84FF;
	Trend := WEN
	//卖出:=CROSS(J2,J1) AND J2>85;
	S1 := CROSS(J2, J1)
	S2 := J2.Gt(85)
	S := S1.And(S2)
	//卖点预警线:90,DOTLINE,LINETHICK1,COLORBLUE;
	//DRAWICON(卖出, 70, 2);
	//买入:=趋势线<REF(趋势线,1) AND 趋势线<=5;
	B := Trend.Lt(REF(Trend, 1)).And(Trend.Lte(5))
	//买点预警线:10,DOTLINE,LINETHICK1,COLORWHITE;
	//DRAWICON(买入, 30, 1);
	//G1:=W1{,LINETHICK2,COLORWHITE};
	//G2:=W2{,LINETHICK2,COLORCYAN};
	WS := Trend.Gte(85.00)
	//STICKLINE(趋势线>=85,100,趋势线,5,1),COLORGREEN;
	//STICKLINE(趋势线<=5,0,趋势线,5,1),COLORYELLOW;
	WB := Trend.Lte(5.00)
	//STICKLINE(COUNT(趋势线<REF(趋势线,1) AND 趋势线<=5,2)=2,0,20,8,0),COLORRED;
	//STICKLINE(CROSS(J2,J1) AND J2>=85,100,80,8,0),COLORGREEN;
	df = pandas.NewDataFrame(df.Col("date"), df.Col("close"))
	OS1 := pandas.NewSeriesWithType(pandas.SERIES_TYPE_BOOL, "S1", S1)
	OS2 := pandas.NewSeriesWithType(pandas.SERIES_TYPE_BOOL, "S2", S2)
	OS := pandas.NewSeriesWithType(pandas.SERIES_TYPE_BOOL, "S", S)
	OB := pandas.NewSeriesWithType(pandas.SERIES_TYPE_BOOL, "B", B)
	OWS := pandas.NewSeriesWithType(pandas.SERIES_TYPE_BOOL, "WS", WS)
	OWB := pandas.NewSeriesWithType(pandas.SERIES_TYPE_BOOL, "WB", WB)
	df = df.Join(OB, OS, OS1, OS2, OWS, OWB)
	return df
}
