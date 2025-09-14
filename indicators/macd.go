package indicators

import (
	"gitee.com/quant1x/pandas"
	. "gitee.com/quant1x/pandas/formula"
)

// MACD 指标
//
//	DIF:EMA(CLOSE,SHORT)-EMA(CLOSE,LONG);
//	输出DIF:收盘价的SHORT日指数移动平均-收盘价的LONG日指数移动平均
//	DEA:EMA(DIF,MID);
//	输出DEA:DIF的MID日指数移动平均
//	MACD:(DIF-DEA)*2,COLORSTICK;
//	输出平滑异同平均线:(DIF-DEA)*2,COLORSTICK
//	系统默认12, 26, 9
//	这里采用5,13,3
func MACD(df pandas.DataFrame, SHORT, LONG, MID int) pandas.DataFrame {
	var (
		CLOSE = df.ColAsNDArray("close")
		//HIGH  = df.ColAsNDArray("high")
		//LOW   = df.ColAsNDArray("low")
	)
	//DIF:EMA(CLOSE,SHORT)-EMA(CLOSE,LONG);
	DIF := EMA(CLOSE, SHORT).Sub(EMA(CLOSE, LONG))
	//DEA:EMA(DIF,MID)
	DEA := EMA(DIF, MID)
	//MACD:(DIF-DEA)*2,COLORSTICK;
	MACD := DIF.Sub(DEA).Mul(2)
	return pandas.NewDataFrame(DIF, DEA, MACD)
}

func macd(data pandas.Series, SHORT, LONG, MID int) pandas.DataFrame {
	//DIF:EMA(CLOSE,SHORT)-EMA(CLOSE,LONG);
	DIF := EMA(data, SHORT).Sub(EMA(data, LONG))
	//DEA:EMA(DIF,MID)
	DEA := EMA(DIF, MID)
	//MACD:(DIF-DEA)*2,COLORSTICK;
	MACD := DIF.Sub(DEA).Mul(2)
	return pandas.NewDataFrame(DIF, DEA, MACD)
}
