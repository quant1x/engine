package indicators

import (
	"github.com/quant1x/pandas"
	. "github.com/quant1x/pandas/formula"
)

// MA4X 缠论时间窗口
func MA4X(df pandas.DataFrame, N int) pandas.DataFrame {
	var (
		CLOSE = df.ColAsNDArray("close")
		HIGH  = df.ColAsNDArray("high")
		LOW   = df.ColAsNDArray("low")
		//VOL   = df.ColAsNDArray("volume")
		//DATE  = df.ColAsNDArray("date")
	)
	//length := CLOSE.Len()
	//N1 := N
	N2 := 5
	N3 := 2
	//N1:=6;
	//重心:(2*C+H+L)/4,COLOR00FFFF,LINETHICK0;
	ZX := CLOSE.Mul(2).Add(HIGH).Add(LOW).Div(4)
	//SJ:=WMA((重心-LLV(L,5))/(HHV(H,5)-LLV(L,5))*100,2);
	LLV5 := LLV(LOW, N2)
	HHV5 := HHV(HIGH, N2)

	sj1 := ZX.Sub(LLV5)
	sj2 := HHV5.Sub(LLV5)
	sj3 := sj1.Div(sj2).Mul(100)
	SJ := WMA(sj3, N3)
	//ZJ:=WMA(0.618*REF(SJ,1)+0.382*SJ,2);
	zj1 := REF(SJ, 1).Mul(0.618)
	zj2 := SJ.Mul(0.382)
	zj3 := zj1.Add(zj2)
	ZJ := WMA(zj3, N3)
	//B:CROSS(SJ,ZJ) AND SJ<30,COLORRED,NODRAW;
	b1 := CROSS(SJ, ZJ)
	//fmt.Println(b1)
	b2 := SJ.Lt(30)
	//fmt.Println(b2)
	B := b2.And(b1)
	//DRAWTEXT(B,L-0.1,'←低吸'),COLOR00FF00;
	//S:CROSS(ZJ,SJ) AND SJ>70,COLORGREEN,NODRAW;
	s1 := CROSS(ZJ, SJ)
	s2 := SJ.Gt(70)
	S := s2.And(s1)
	//DRAWTEXT(S,H+0.1,'←高抛'),COLOR0077FF;
	OB := pandas.NewSeriesWithType(pandas.SERIES_TYPE_BOOL, "B", B)
	OS := pandas.NewSeriesWithType(pandas.SERIES_TYPE_BOOL, "S", S)
	//df := pandas.NewDataFrame(OB, OS)
	df = pandas.NewDataFrame(df.Col("date"), df.Col("close"))
	df = df.Join(ZX, SJ, ZJ, OB, OS)
	return df
}
