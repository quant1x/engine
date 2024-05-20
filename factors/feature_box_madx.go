package factors

import (
	"gitee.com/quant1x/engine/utils"
	"gitee.com/quant1x/pandas"
	. "gitee.com/quant1x/pandas/formula"
)

type JuXianDongXiang struct {
	Dm0       float64
	Dm1       float64
	Dm2       float64
	Diverging float64
	B         bool
	//S         bool
}

// 多空趋势
func computeJuXianDongXiang(OPEN, CLOSE, HIGH, LOW pandas.Series) *JuXianDongXiang {
	//{均线动向, V1.0.3, 2023-09-18}
	//P0:=3;
	P0 := 3
	//P1:=5;
	P1 := 5
	//P2:=10;
	P2 := 10
	//P3:=20;
	P3 := 20
	//MX:=CLOSE;
	MX := CLOSE
	//MA0:=MA(MX,P0);
	MA0 := MA(MX, P0)
	//MA1:=MA(MX,P1);
	MA1 := MA(MX, P1)
	//MA2:=MA(MX,P2);
	MA2 := MA(MX, P2)
	//MA3:=MA(MX,P3);
	MA3 := MA(MX, P3)
	//DM0:MA0-MA3;
	DM0 := MA0.Sub(MA3)
	//DM1:MA1-MA3;
	DM1 := MA1.Sub(MA3)
	//DM2:MA2-MA3;
	DM2 := MA2.Sub(MA3)
	//X0:=DM0-REF(DM0,1);
	X0 := DM0.Sub(REF(DM0, 1))
	//X1:=DM1-REF(DM1,1);
	X1 := DM1.Sub(REF(DM1, 1))
	//X2:=DM2-REF(DM2,1);
	X2 := DM2.Sub(REF(DM2, 1))
	//DIVERGING:X0+X1+X2;
	DIVERGING := X0.Add(X1).Add(X2)
	//B:X0>0 AND X1>0 AND X2>0,NODRAW;
	B := X0.Gt(0).And(X1.Gt(0)).And(X2.Gt(0))
	//DRAWICON(B,0.01,1);
	madx := JuXianDongXiang{
		Dm0:       utils.Float64IndexOf(DM0, -1),
		Dm1:       utils.Float64IndexOf(DM1, -1),
		Dm2:       utils.Float64IndexOf(DM2, -1),
		Diverging: utils.Float64IndexOf(DIVERGING, -1),
		B:         utils.BoolIndexOf(B, -1),
	}
	return &madx
}
