package factors

import (
	"context"
	"gitee.com/quant1x/engine/cache"
	"gitee.com/quant1x/engine/datasource/dfcf"
	"gitee.com/quant1x/exchange"
	"gitee.com/quant1x/gotdx/quotes"
	"gitee.com/quant1x/gox/api"
	"gitee.com/quant1x/gox/logger"
	"gitee.com/quant1x/gox/progressbar"
)

// DataPreviewReport 业绩预告(Listed Companies'Performance Forecast)
type DataPreviewReport struct {
	Manifest
}

func init() {
	summary := mapDataSets[BasePerformanceForecast]
	_ = cache.Register(&DataPreviewReport{Manifest: Manifest{DataSummary: summary}})
}

func (this *DataPreviewReport) Clone(date string, code string) DataSet {
	summary := mapDataSets[BasePerformanceForecast]
	var dest = DataPreviewReport{
		Manifest: Manifest{
			DataSummary: summary,
			Date:        date,
			Code:        code,
		},
	}
	return &dest
}

func (this *DataPreviewReport) Init(ctx context.Context, date string) error {
	barIndex := 1
	value, ok := ctx.Value(cache.KBarIndex).(int)
	if ok {
		barIndex = value
	}
	barIndex++
	modName := "业绩预告"
	logger.Info(modName + ", 任务开始启动...")

	var allReports []dfcf.PreviewQuarterlyReport
	// 确定更新日期
	qBegin, _ := api.GetQuarterDayByDate(date, 0)
	quarterBeginDate := exchange.FixTradeDate(qBegin)
	list, pages, _, _ := dfcf.FinanceReports(quarterBeginDate, 1)
	if pages < 1 || len(list) == 0 {
		return nil
	}
	allReports = append(allReports, list...)

	barSub := progressbar.NewBar(barIndex, "评估["+modName+"]", pages)
	for pageNo := 2; pageNo < pages+1; pageNo++ {
		barSub.Add(1)
		list, pages, count, err := dfcf.FinanceReports(quarterBeginDate, pageNo)
		if err != nil || pages < 1 {
			logger.Error(err)
			break
		}
		if count == 0 {
			break
		}
		allReports = append(allReports, list...)
		if count < dfcf.EastmoneyFinanceReportsPageSize {
			break
		}
	}
	if len(allReports) > 0 {
		filename := cache.PreviewReportFilename(quarterBeginDate)
		_ = api.SlicesToCsv(filename, allReports)
	}
	logger.Info(modName + ", 任务开始结束...")
	return nil
}

func (this *DataPreviewReport) Update(date string) {
	_ = date
}

func (this *DataPreviewReport) Repair(date string) {
	_ = date
}

func (this *DataPreviewReport) Increase(snapshot quotes.Snapshot) {
	_ = snapshot
}

func (this *DataPreviewReport) Print(code string, date ...string) {
	_ = code
	_ = date
}
