package indicators

import (
	"github.com/quant1x/pandas"
	. "github.com/quant1x/pandas/formula"
)

// RSI 指标
//
//	LC:=REF(CLOSE,1);
//	LC赋值:1日前的收盘价
//	RSI1:SMA(MAX(CLOSE-LC,0),N1,1)/SMA(ABS(CLOSE-LC),N1,1)*100;
//	输出RSI1:收盘价-LC和0的较大值的N1日[1日权重]移动平均/收盘价-LC的绝对值的N1日[1日权重]移动平均*100
//	RSI2:SMA(MAX(CLOSE-LC,0),N2,1)/SMA(ABS(CLOSE-LC),N2,1)*100;
//	输出RSI2:收盘价-LC和0的较大值的N2日[1日权重]移动平均/收盘价-LC的绝对值的N2日[1日权重]移动平均*100
//	RSI3:SMA(MAX(CLOSE-LC,0),N3,1)/SMA(ABS(CLOSE-LC),N3,1)*100;
//	输出RSI3:收盘价-LC和0的较大值的N3日[1日权重]移动平均/收盘价-LC的绝对值的N3日[1日权重]移动平均*100
//
//	系统默认参数6,12,24
func RSI(df pandas.DataFrame, N1, N2, N3 int) pandas.DataFrame {
	var (
		CLOSE = df.ColAsNDArray("close")
		//HIGH  = df.ColAsNDArray("high")
		//LOW   = df.ColAsNDArray("low")
	)
	//	LC:=REF(CLOSE,1);
	LC := REF(CLOSE, 1)
	cls := CLOSE.Sub(LC)
	//	RSI1:SMA(MAX(CLOSE-LC,0),N1,1)/SMA(ABS(CLOSE-LC),N1,1)*100;
	RSI1 := SMA(MAX(cls, 0), N1, 1).Div(SMA(ABS(cls), N1, 1)).Mul(100)
	//	RSI2:SMA(MAX(CLOSE-LC,0),N2,1)/SMA(ABS(CLOSE-LC),N2,1)*100;
	RSI2 := SMA(MAX(cls, 0), N2, 1).Div(SMA(ABS(cls), N2, 1)).Mul(100)
	//	RSI3:SMA(MAX(CLOSE-LC,0),N3,1)/SMA(ABS(CLOSE-LC),N3,1)*100;
	RSI3 := SMA(MAX(cls, 0), N3, 1).Div(SMA(ABS(cls), N3, 1)).Mul(100)

	return pandas.NewDataFrame(RSI1, RSI2, RSI3)
}
