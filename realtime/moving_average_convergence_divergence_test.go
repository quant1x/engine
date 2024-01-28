package realtime

import (
	"fmt"
	"gitee.com/quant1x/engine/datasource/base"
	"gitee.com/quant1x/engine/models"
	"gitee.com/quant1x/pandas"
	. "gitee.com/quant1x/pandas/formula"
	"testing"
	"time"
)

func TestIncrementalMovingAverageConvergenceDivergence(t *testing.T) {
	code := "002528"
	date := "20240126"
	klines := base.CheckoutKLines(code, date)
	if len(klines) == 0 {
		panic("no data")
	}
	df := pandas.LoadStructs(klines)
	if df.Nrow() == 0 {
		panic("加载k线失败")
	}
	var (
		DATE  = df.Col("date")
		CLOSE = df.ColAsNDArray("close")
	)

	SHORT := 12
	LONG := 26
	MID := 9

	short := EMA(CLOSE, SHORT)
	long := EMA(CLOSE, LONG)
	DIF := short.Sub(long)
	DEA := EMA(DIF, MID)
	MACD := DIF.Sub(DEA).Mul(2)

	df = pandas.NewDataFrame(
		DATE,
		pandas.NewSeriesWithoutType("short", short),
		pandas.NewSeriesWithoutType("long", long),
		pandas.NewSeriesWithoutType("DIF", DIF),
		pandas.NewSeriesWithoutType("DEA", DEA),
		pandas.NewSeriesWithoutType("MACD", MACD),
	)
	fmt.Println(df)
	fmt.Println("==============================================================================================================")
	lastClose := CLOSE.IndexOf(-1).(float64)
	fmt.Println("确定最新收盘价:", lastClose)
	df1 := df.IndexOf(-2)
	fmt.Println(df1)
	//date0 := df1["date"]
	lastShort := df1["short"].(float64)
	lastLong := df1["long"].(float64)
	lastDif := df1["DIF"].(float64)
	lastDea := df1["DEA"].(float64)
	fmt.Println("lastDif", lastDif)
	dif, dea, macd := IncrementalMovingAverageConvergenceDivergence(lastClose, lastShort, lastLong, lastDea, SHORT, LONG, MID)

	fmt.Println("date:", DATE.IndexOf(-1))
	fmt.Println(" dif:", dif)
	fmt.Println(" dea:", dea)
	fmt.Println("macd:", macd)
}

func TestDynamicMovingAverageConvergenceDivergence(t *testing.T) {
	code := "002528"
	date := "20240126"
	klines := base.CheckoutKLines(code, date)
	if len(klines) == 0 {
		panic("no data")
	}
	df := pandas.LoadStructs(klines)
	if df.Nrow() == 0 {
		panic("加载k线失败")
	}
	var (
		DATE  = df.Col("date")
		CLOSE = df.ColAsNDArray("close")
	)

	SHORT := 12
	LONG := 26
	MID := 9

	short := EMA(CLOSE, SHORT)
	long := EMA(CLOSE, LONG)
	DIF := short.Sub(long)
	DEA := EMA(DIF, MID)
	MACD := DIF.Sub(DEA).Mul(2)

	df = pandas.NewDataFrame(
		DATE,
		pandas.NewSeriesWithoutType("short", short),
		pandas.NewSeriesWithoutType("long", long),
		pandas.NewSeriesWithoutType("DIF", DIF),
		pandas.NewSeriesWithoutType("DEA", DEA),
		pandas.NewSeriesWithoutType("MACD", MACD),
	)
	fmt.Println(df)
	fmt.Println("==============================================================================================================")
	lastClose := CLOSE.IndexOf(-1).(float64)
	fmt.Println("确定最新收盘价:", lastClose)
	df1 := df.IndexOf(-2)
	fmt.Println(df1)
	//date0 := df1["date"]
	lastShort := df1["short"].(float64)
	lastLong := df1["long"].(float64)
	lastDif := df1["DIF"].(float64)
	lastDea := df1["DEA"].(float64)
	fmt.Println("lastDif", lastDif)
	snapshot := models.GetTick(code)
	now := time.Now()
	macd, macdHigh, macdLow := DynamicMovingAverageConvergenceDivergence(*snapshot, lastShort, lastLong, lastDea, SHORT, LONG, MID)
	fmt.Println("   times:", time.Since(now))
	fmt.Println("    date:", DATE.IndexOf(-1))
	fmt.Println("    macd:", macd)
	fmt.Println("macdHigh:", macdHigh)
	fmt.Println(" macdLow:", macdLow)
}
