package features

import (
	"gitee.com/quant1x/engine/cache"
	"gitee.com/quant1x/engine/datasets/dfcf"
	"gitee.com/quant1x/gox/api"
	"gitee.com/quant1x/gox/logger"
)

var (
	__mapQuarterlyReports = map[string]dfcf.QuarterlyReport{}
)

func loadQuarterlyReports(date string) {
	var allReports []dfcf.QuarterlyReport
	filename := cache.ReportsFilename(date)
	err := api.CsvToSlices(filename, &allReports)
	if err != nil {
		logger.Errorf("cache %s failed, error: %+v", filename, err)
	}
	if len(allReports) > 0 {
		for _, v := range allReports {
			__mapQuarterlyReports[v.SecurityCode] = v
		}
	}
}

type quarterlyReportSummary struct {
	BPS      float64
	BasicEPS float64
}

func getQuarterlyReportSummary(securityCode string) quarterlyReportSummary {
	var summary quarterlyReportSummary
	v, ok := __mapQuarterlyReports[securityCode]
	if ok {
		summary.BPS = v.BPS
		summary.BasicEPS = v.BasicEPS
	}
	return summary
}
