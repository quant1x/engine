package factors

import (
	"gitee.com/quant1x/engine/cache"
	dfcf2 "gitee.com/quant1x/engine/datasource/dfcf"
	"gitee.com/quant1x/gotdx/proto"
	"gitee.com/quant1x/gox/api"
	"gitee.com/quant1x/gox/logger"
)

var (
	__mapQuarterlyReports = map[string]dfcf2.QuarterlyReport{}
)

func loadQuarterlyReports(date string) {
	var allReports []dfcf2.QuarterlyReport
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

func getQuarterlyYearQuarter(date string) string {
	q, _, _ := api.GetQuarterByDate(date, 1)
	return q
}

// 季报概要
type quarterlyReportSummary struct {
	QDate              string
	BPS                float64
	BasicEPS           float64
	TotalOperateIncome float64
	DeductBasicEPS     float64
}

func (q *quarterlyReportSummary) Assign(v dfcf2.QuarterlyReport) {
	q.BPS = v.BPS
	q.BasicEPS = v.BasicEPS
	q.TotalOperateIncome = v.TotalOperateIncome
	q.DeductBasicEPS = v.DeductBasicEPS
	q.QDate = v.QDATE
}

func getQuarterlyReportSummary(securityCode, date string) quarterlyReportSummary {
	var summary quarterlyReportSummary
	if proto.AssertIndexBySecurityCode(securityCode) {
		return summary
	}
	v, ok := __mapQuarterlyReports[securityCode]
	if ok {
		summary.Assign(v)
		return summary
	}
	q := dfcf2.GetCacheQuarterlyReportsBySecurityCode(securityCode, date)
	if q != nil {
		summary.Assign(*q)
	}
	return summary
}