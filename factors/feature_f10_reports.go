package factors

import (
	"gitee.com/quant1x/engine/cache"
	"gitee.com/quant1x/engine/datasets/dfcf"
	"gitee.com/quant1x/gotdx/proto"
	"gitee.com/quant1x/gox/api"
	"gitee.com/quant1x/gox/logger"
)

var (
	__mapQuarterlyReports = map[string]dfcf.QuarterlyReport{}
)

func loadQuarterlyReports(date string) {
	var allReports []dfcf.QuarterlyReport
	_, qEnd := api.GetQuarterDayByDate(date)
	filename := cache.ReportsFilename(qEnd)
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

func getQuarterlyReportSummary(securityCode, date string) quarterlyReportSummary {
	var summary quarterlyReportSummary
	if proto.AssertIndexBySecurityCode(securityCode) {
		return summary
	}
	v, ok := __mapQuarterlyReports[securityCode]
	if ok {
		summary.BPS = v.BPS
		summary.BasicEPS = v.BasicEPS
		return summary
	}
	q := dfcf.GetCacheQuarterlyReportsBySecurityCode(securityCode, date)
	if q != nil {
		summary.BPS = q.BPS
		summary.BasicEPS = q.BasicEPS
	}
	return summary
}
