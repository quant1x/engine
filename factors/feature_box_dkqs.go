package factors

import (
	"gitee.com/quant1x/engine/utils"
	"gitee.com/quant1x/pandas"
	. "gitee.com/quant1x/pandas/formula"
)

// DuoKongQuShi 多空趋势
type DuoKongQuShi struct {
	Col float64 // 多空能量
	K0  float64
	K   float64
	D   float64
	B   bool
	S   bool
}

// 多空趋势
func computeDuoKongQuShi(OPEN, CLOSE, HIGH, LOW pandas.Series) *DuoKongQuShi {
	//{多空趋势, V1.1.2, 2023-09-13}
	//{量能柱}
	//SCALE:=100;
	SCALE := 100
	//N:=3;
	N := 3
	//NN:=MIN(BARSCOUNT(CLOSE),N);
	//NN := MIN(BARSCOUNT(CLOSE), N)
	NN := N
	//CNN:=REF(CLOSE,NN);
	CNN := REF(CLOSE, NN)
	//CDIFF:=CLOSE-CNN,COLORSTICK;
	CDIFF := CLOSE.Sub(CNN)
	//FF:=CDIFF/CNN;
	FF := CDIFF.Div(CNN)
	//MADK:SCALE*FF,COLORSTICK;
	MADK := FF.Mul(SCALE)
	//
	//{多空趋势}
	//MAXH:=MAX(HIGH,REF(CLOSE,1));
	//MINL:=MIN(LOW,REF(CLOSE,1));
	MINL := MIN(LOW, REF(CLOSE, 1))
	//
	//DIFFK:=IFF(OPEN>=CLOSE,OPEN-CLOSE,OPEN-LOW);
	DIFFK := IFF(OPEN.Gte(CLOSE), OPEN.Sub(CLOSE), OPEN.Sub(LOW))
	//TZ1:=OPEN-REF(CLOSE,1)-DIFFK;
	TZ1 := OPEN.Sub(REF(CLOSE, 1)).Sub(DIFFK)
	//TZ2:=TZ1/OPEN;
	TZ2 := TZ1.Div(OPEN)
	//K0:SCALE*TZ2;
	K0 := TZ2.Mul(SCALE)
	//K:ABS(K0),DOTLINE;
	K := ABS(K0)
	//DIFFD:=IFF(OPEN>=CLOSE,HIGH-OPEN,HIGH-CLOSE);
	DIFFD := IFF(OPEN.Gte(CLOSE), HIGH.Sub(OPEN), HIGH.Sub(CLOSE))
	//TD1:=CLOSE-MINL+DIFFD;
	TD1 := CLOSE.Sub(MINL).Add(DIFFD)
	//TD2:=TD1/CLOSE;
	TD2 := TD1.Div(CLOSE)
	//D:SCALE*TD2;
	D := TD2.Mul(SCALE)
	//B:CROSS(D,K);
	B := CROSS(D, K)
	//S:CROSS(K,D);
	S := CROSS(K, D)
	//DRAWICON(B,20,1);
	//DRAWICON(S,20,2);
	dkqs := DuoKongQuShi{
		Col: utils.Float64IndexOf(MADK, -1),
		K0:  utils.Float64IndexOf(K0, -1),
		K:   utils.Float64IndexOf(K, -1),
		D:   utils.Float64IndexOf(D, -1),
		B:   utils.BoolIndexOf(B, -1),
		S:   utils.BoolIndexOf(S, -1),
	}
	return &dkqs
}
