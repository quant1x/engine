package factors

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"testing"

	"gitee.com/quant1x/data/exchange"
	"gitee.com/quant1x/data/level1/securities"
	"gitee.com/quant1x/engine/cache"
	"gitee.com/quant1x/engine/utils"
	"gitee.com/quant1x/gox/api"
	"gitee.com/quant1x/pandas"
	"gitee.com/quant1x/pandas/formula"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
)

func TestKLine(t *testing.T) {
	code := "002882"
	securityCode := exchange.CorrectSecurityCode(code)
	var wides []SecurityFeature
	filename := cache.KLineFilename(securityCode)
	err := api.CsvToSlices(filename, &wides)
	if err != nil || len(wides) == 0 {
		return
	}
	wides = wides[len(wides)-1:]
	df := pandas.LoadStructs(wides)
	fmt.Println(df)
}

func TestKLineWide(t *testing.T) {
	code := "002882"
	code = "600178"
	securityCode := exchange.CorrectSecurityCode(code)
	var wides []SecurityFeature
	filename := cache.WideFilename(securityCode)
	err := api.CsvToSlices(filename, &wides)
	if err != nil || len(wides) == 0 {
		return
	}
	//wides = wides[len(wides)-1:]
	df := pandas.LoadStructs(wides)
	fmt.Println(df)
}

func TestDataSetWide_pullWideByDate(t *testing.T) {
	code := "sz301129"
	code = "002857"
	code = "603230"
	date := "2024-06-28"
	securityCode := exchange.CorrectSecurityCode(code)
	lines := pullWideByDate(securityCode, date)
	df := pandas.LoadStructs(lines)
	fmt.Println(df)
}

func TestWideTableValuate(t *testing.T) {
	const MAX_ROWS = 21
	code := "002615"
	code = "300461"
	code = "000004"
	code = "002173"
	code = "002766"
	//code = "603586"
	code = "300717"
	code = "300107"
	code = "002085"
	code = "002823"
	code = "002857"
	code = "300947"
	code = "301397"
	code = "300561"
	code = "000702"
	code = "605577"
	code = "603230"
	code = "880866"
	date := "2024-06-28"
	code = exchange.CorrectSecurityCode(code)
	lines := CheckoutWideTableByDate(code, date)
	if len(lines) > MAX_ROWS {
		lines = lines[len(lines)-MAX_ROWS:]
	}
	df := pandas.LoadStructs(lines)
	dates := df.Col("date").Strings()
	CLOSE := df.ColAsNDArray("close")
	low := df.ColAsNDArray("low")
	vol := df.ColAsNDArray("volume")
	amt := df.ColAsNDArray("amount")
	ap := amt.Div(vol)
	sAp := pandas.SeriesWithName("均价", ap.Float64s())

	ov := df.ColAsNDArray("outer_volume")
	oa := df.ColAsNDArray("outer_amount")
	abp := oa.Div(ov)
	sAbp := pandas.SeriesWithName("买入均价", abp.Float64s())

	iv := df.ColAsNDArray("inner_volume")
	ia := df.ColAsNDArray("inner_amount")
	asp := ia.Div(iv)
	sAsp := pandas.SeriesWithName("卖出均价", asp.Float64s())
	bs := abp.Div(asp)
	sBs := pandas.SeriesWithName("买卖均价比", bs.Float64s())
	r1Low := formula.REF(low, 1)
	lp := low.Sub(r1Low).Div(r1Low)
	sLow := pandas.SeriesWithName("创新低速度环比", lp.Float64s())

	df = df.Select([]string{"date", "low", "amount", "volume"})
	df = df.Join(sAp, sAbp, sAsp, sBs, sLow)
	name := securities.GetStockName(code)
	fmt.Printf("%s, %s\n", code, name)
	fmt.Println(df)
	// create a new line instance
	line := charts.NewLine()
	// set some global options like Title/Legend/ToolTip or anything else
	line.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{Theme: types.ThemeWesteros}),
		charts.WithTitleOpts(opts.Title{
			Title:    fmt.Sprintf("%s(%s)", name, code),
			Subtitle: "买卖双方意愿强度观测",
		}))

	// Put data into instance
	line.SetXAxis(dates).
		AddSeries("收盘价", toLineItems(CLOSE.Float64s())).
		AddSeries("均价线", toLineItems(ap.Float64s())).
		AddSeries("买入均价", toLineItems(abp.Float64s())).
		AddSeries("卖出均价", toLineItems(asp.Float64s())).
		SetSeriesOptions(charts.WithLineChartOpts(opts.LineChart{Smooth: opts.Bool(true)}))
	line.SetGlobalOptions(
		charts.WithXAxisOpts(opts.XAxis{
			AxisLabel: &opts.AxisLabel{
				Show:         opts.Bool(true),
				Interval:     "0",
				Rotate:       45, // 减小旋转角度
				ShowMinLabel: opts.Bool(true),
				ShowMaxLabel: opts.Bool(true),
				FontSize:     10, // 减小标签字体大小
			},
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Scale: opts.Bool(true),
		}),
	)
	filename := "wide-" + code + ".html"
	currentPath, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	fmt.Println(currentPath)
	filename = filepath.Join(currentPath, filename)
	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	line.PageTitle = fmt.Sprintf("%s(%s)", name, code)
	line.Render(f)
	err = utils.OpenURL("file://" + filename)
	fmt.Println(err)
}

// generate random data for line chart
func generateLineItems() []opts.LineData {
	items := make([]opts.LineData, 0)
	for i := 0; i < 7; i++ {
		items = append(items, opts.LineData{Value: rand.Intn(300)})
	}
	return items
}

func toLineItems(data []float64) []opts.LineData {
	items := make([]opts.LineData, 0)
	for i := 0; i < len(data); i++ {
		items = append(items, opts.LineData{Value: data[i]})
	}
	return items
}
