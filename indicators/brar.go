package indicators

import (
	"github.com/quant1x/pandas"
	. "github.com/quant1x/pandas/formula"
)

// BRAR 情绪指标
//
//	BR:SUM(MAX(0,HIGH-REF(CLOSE,1)),N)/SUM(MAX(0,REF(CLOSE,1)-LOW),N)*100;
//	输出BR:0和最高价-1日前的收盘价的较大值的N日累和/0和1日前的收盘价-最低价的较大值的N日累和*100
//	AR:SUM(HIGH-OPEN,N)/SUM(OPEN-LOW,N)*100;
//	输出AR:最高价-开盘价的N日累和/开盘价-最低价的N日累和*100
func BRAR(df pandas.DataFrame, N any) pandas.DataFrame {
	var (
		OPEN  = df.ColAsNDArray("open")
		CLOSE = df.ColAsNDArray("close")
		HIGH  = df.ColAsNDArray("high")
		LOW   = df.ColAsNDArray("low")
	)
	//BR:SUM(MAX(0,HIGH-REF(CLOSE,1)),N)/SUM(MAX(0,REF(CLOSE,1)-LOW),N)*100;
	c1 := REF(CLOSE, 1)
	BR := SUM(MAX(HIGH.Sub(c1), 0), N).Div(SUM(MAX(c1.Sub(LOW), 0), N)).Mul(100)
	//AR:SUM(HIGH-OPEN,N)/SUM(OPEN-LOW,N)*100;
	AR := SUM(HIGH.Sub(OPEN), N).Div(SUM(OPEN.Sub(LOW), N)).Mul(100)
	return pandas.NewDataFrame(BR, AR)
}
