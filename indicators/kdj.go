package indicators

import (
	"github.com/quant1x/pandas"
	. "github.com/quant1x/pandas/formula"
)

// KDJ 指标
//
//	RSV:=(CLOSE-LLV(LOW,N))/(HHV(HIGH,N)-LLV(LOW,N))*100;
//	RSV赋值:(收盘价-N日内最低价的最低值)/(N日内最高价的最高值-N日内最低价的最低值)*100
//	K:EMA(RSV,M1,1);
//	输出K:RSV的M1日[1日权重]移动平均
//	D:EMA(K,M2,1);
//	输出D:K的M2日[1日权重]移动平均
//	J:3*K-2*D;
//	输出J:3*K-2*D
func KDJ(df pandas.DataFrame, N, M1, M2 int) pandas.DataFrame {
	//CLOSE, HIGH, LOW stat.Series
	var (
		CLOSE = df.ColAsNDArray("close")
		HIGH  = df.ColAsNDArray("high")
		LOW   = df.ColAsNDArray("low")
	)

	// 计算N周期的最高价序列
	x01 := HHV(HIGH, N)
	// 计算N周期的最低价序列
	x02 := LLV(LOW, N)

	// CLOSE-LLV(LOW,N)
	x11 := CLOSE.Sub(x02)
	// (HHV(HIGH,N)-LLV(LOW,N))*100
	x12 := x01.Sub(x02)

	//(CLOSE-LLV(LOW,N))/(HHV(HIGH,N)-LLV(LOW,N))*100
	RSV := x11.Div(x12).Mul(100)
	//K:EMA(RSV,M1,1);
	K := EMA(RSV, M1*2-1)
	//D:EMA(K,M2,1)
	D := EMA(K, M2*2-1)
	// 3*K-2*D;
	J := K.Mul(3).Sub(D.Mul(2))

	return pandas.NewDataFrame(K, D, J)
}
