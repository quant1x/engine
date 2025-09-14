package dfcf

import (
	"fmt"
	urlpkg "net/url"
	"sync"

	"github.com/quant1x/engine/cache"
	"github.com/quant1x/exchange"
	"github.com/quant1x/pkg/fastjson"
	"github.com/quant1x/x/api"
	"github.com/quant1x/x/http"
	"github.com/quant1x/x/logger"
)

// QuarterlyReports 分页获取季报数据
func QuarterlyReports(featureDate string, pageNumber ...int) (reports []QuarterlyReport, pages int, err error) {
	pageNo := 1
	if len(pageNumber) > 0 {
		pageNo = pageNumber[0]
	}
	qBegin, qEnd := api.GetQuarterDayByDate(featureDate)
	quarterBeginDate := exchange.FixTradeDate(qBegin)
	quarterEndDate := exchange.FixTradeDate(qEnd)
	params := urlpkg.Values{
		//"callback":    {"jQuery1123043614175387302234_1685785566671"},
		//"sortColumns": {"UPDATE_DATE,SECURITY_CODE"},
		"sortColumns": {"REPORTDATE,SECURITY_CODE"},
		"sortTypes":   {"-1,1"},
		"pageSize":    {fmt.Sprint(EastmoneyQuarterlyReportAllPageSize)},
		"pageNumber":  {fmt.Sprintf("%d", pageNo)},
		"reportName":  {"RPT_LICO_FN_CPD"},
		"columns":     {"ALL"},
		"filter":      {fmt.Sprintf("(REPORTDATE='%s')", quarterEndDate)},
	}
	_ = quarterBeginDate
	_ = quarterEndDate

	url := urlQuarterlyReportAll + "?" + params.Encode()
	data, err := http.Get(url)
	//fmt.Println(api.Bytes2String(data))
	obj, err := fastjson.ParseBytes(data)
	if err != nil {
		logger.Errorf("%+v\n", err)
		return
	}

	result := obj.Get("result")
	list := result.GetArray("data")
	pages = result.GetInt("pages")
	if len(list) > 0 {
		for _, v := range list {
			report := QuarterlyReport{
				SecuCode:           v.GetString("SECUCODE"),
				UpdateDate:         v.GetString("UPDATE_DATE"),
				ReportDate:         v.GetString("REPORTDATE"),
				BasicEPS:           v.GetFloat64("BASIC_EPS"),
				DeductBasicEPS:     v.GetFloat64("DEDUCT_BASIC_EPS"),
				BPS:                v.GetFloat64("BPS"),
				NoticeDate:         v.GetString("NOTICE_DATE"),
				IsNew:              v.GetString("ISNEW"),
				ORGCODE:            v.GetString("ORG_CODE"),
				TRADEMARKETZJG:     v.GetString("TRADE_MARKET_ZJG"),
				QDATE:              v.GetString("QDATE"),
				DATATYPE:           v.GetString("DATATYPE"),
				DATAYEAR:           v.GetString("DATAYEAR"),
				DATEMMDD:           v.GetString("DATEMMDD"),
				EITIME:             v.GetString("EITIME"),
				SECURITYCODE:       v.GetString("SECURITY_CODE"),
				SECURITYNAMEABBR:   v.GetString("SECURITY_NAME_ABBR"),
				TRADEMARKETCODE:    v.GetString("TRADE_MARKET_CODE"),
				TRADEMARKET:        v.GetString("TRADE_MARKET"),
				SECURITYTYPECODE:   v.GetString("SECURITY_TYPE_CODE"),
				SECURITYTYPE:       v.GetString("SECURITY_TYPE"),
				TotalOperateIncome: v.GetFloat64("TOTAL_OPERATE_INCOME"),
				ParentNetprofit:    v.GetFloat64("PARENT_NETPROFIT"),
				WeightAvgRoe:       v.GetFloat64("WEIGHTAVG_ROE"),
				YSTZ:               v.GetFloat64("YSTZ"),
				SJLTZ:              v.GetFloat64("SJLTZ"),
				MGJYXJJE:           v.GetFloat64("MGJYXJJE"),
				XSMLL:              v.GetFloat64("XSMLL"),
				YSHZ:               v.GetFloat64("YSHZ"),
				SJLHZ:              v.GetFloat64("SJLHZ"),
				ASSIGNDSCRPT:       v.GetFloat64("ASSIGNDSCRPT"),
				PAYYEAR:            v.GetFloat64("PAYYEAR"),
				PUBLISHNAME:        v.GetString("PUBLISHNAME"),
				ZXGXL:              v.GetFloat64("ZXGXL"),
			}
			// 截取市场编码，截取股票编码，市场编码+股票编码拼接作为主键
			securityCode := exchange.CorrectSecurityCode(report.SecuCode)
			report.SecurityCode = securityCode
			reports = append(reports, report)
		}
	}
	return
}

// QuarterlyReportsBySecurityCode 分页获取季报数据
func QuarterlyReportsBySecurityCode(securityCode, date string, diffQuarters int, pageNumber ...int) (reports []QuarterlyReport) {
	pageNo := 1
	if len(pageNumber) > 0 {
		pageNo = pageNumber[0]
	}
	_, _, code := exchange.DetectMarket(securityCode)
	quarterEndDate := exchange.FixTradeDate(date)
	//_, _, qEnd := api.GetQuarterByDate(date, diffQuarters)
	//quarterEndDate = trading.FixTradeDate(qEnd)
	params := urlpkg.Values{
		"sortColumns": {"REPORTDATE,SECURITY_CODE"},
		"sortTypes":   {"-1,1"},
		"pageSize":    {fmt.Sprint(EastmoneyQuarterlyReportAllPageSize)},
		"pageNumber":  {fmt.Sprintf("%d", pageNo)},
		"reportName":  {"RPT_LICO_FN_CPD"},
		"columns":     {"ALL"},
		"filter":      {fmt.Sprintf("(SECURITY_CODE=\"%s\")(REPORTDATE='%s')", code, quarterEndDate)},
	}

	url := urlQuarterlyReportAll + "?" + params.Encode()
	data, err := http.Get(url)
	//fmt.Println(api.Bytes2String(data))
	obj, err := fastjson.ParseBytes(data)
	if err != nil {
		logger.Errorf("%+v\n", err)
		return
	}

	result := obj.Get("result")
	list := result.GetArray("data")
	pages := result.GetInt("pages")
	_ = pages
	if len(list) > 0 {
		for _, v := range list {
			report := QuarterlyReport{
				SecuCode:           v.GetString("SECUCODE"),
				UpdateDate:         v.GetString("UPDATE_DATE"),
				ReportDate:         v.GetString("REPORTDATE"),
				BasicEPS:           v.GetFloat64("BASIC_EPS"),
				DeductBasicEPS:     v.GetFloat64("DEDUCT_BASIC_EPS"),
				BPS:                v.GetFloat64("BPS"),
				NoticeDate:         v.GetString("NOTICE_DATE"),
				IsNew:              v.GetString("ISNEW"),
				ORGCODE:            v.GetString("ORG_CODE"),
				TRADEMARKETZJG:     v.GetString("TRADE_MARKET_ZJG"),
				QDATE:              v.GetString("QDATE"),
				DATATYPE:           v.GetString("DATATYPE"),
				DATAYEAR:           v.GetString("DATAYEAR"),
				DATEMMDD:           v.GetString("DATEMMDD"),
				EITIME:             v.GetString("EITIME"),
				SECURITYCODE:       v.GetString("SECURITY_CODE"),
				SECURITYNAMEABBR:   v.GetString("SECURITY_NAME_ABBR"),
				TRADEMARKETCODE:    v.GetString("TRADE_MARKET_CODE"),
				TRADEMARKET:        v.GetString("TRADE_MARKET"),
				SECURITYTYPECODE:   v.GetString("SECURITY_TYPE_CODE"),
				SECURITYTYPE:       v.GetString("SECURITY_TYPE"),
				TotalOperateIncome: v.GetFloat64("TOTAL_OPERATE_INCOME"),
				ParentNetprofit:    v.GetFloat64("PARENT_NETPROFIT"),
				WeightAvgRoe:       v.GetFloat64("WEIGHTAVG_ROE"),
				YSTZ:               v.GetFloat64("YSTZ"),
				SJLTZ:              v.GetFloat64("SJLTZ"),
				MGJYXJJE:           v.GetFloat64("MGJYXJJE"),
				XSMLL:              v.GetFloat64("XSMLL"),
				YSHZ:               v.GetFloat64("YSHZ"),
				SJLHZ:              v.GetFloat64("SJLHZ"),
				ASSIGNDSCRPT:       v.GetFloat64("ASSIGNDSCRPT"),
				PAYYEAR:            v.GetFloat64("PAYYEAR"),
				PUBLISHNAME:        v.GetString("PUBLISHNAME"),
				ZXGXL:              v.GetFloat64("ZXGXL"),
			}
			// 截取市场编码，截取股票编码，市场编码+股票编码拼接作为主键
			securityCode := exchange.CorrectSecurityCode(report.SecuCode)
			report.SecurityCode = securityCode
			reports = append(reports, report)
		}
	}
	return
}

var (
	mutexReports sync.RWMutex
	mapReports   = map[string][]QuarterlyReport{}
	chanReports  = make(chan int, 1)
)

// 获取指定个股和周期的季报
func cacheQuarterlyReportsBySecurityCode(securityCode, date string, diffQuarters ...int) *QuarterlyReport {
	diff := 1
	if len(diffQuarters) > 0 {
		diff = diffQuarters[0]
	}
	_, _, last := api.GetQuarterByDate(date, diff)
	filename := cache.ReportsFilename(last)
	var allReports []QuarterlyReport

	mutexReports.Lock()
	defer mutexReports.Unlock()
	allReports, ok := mapReports[filename]
	if !ok && api.FileExist(filename) {
		_ = api.CsvToSlices(filename, &allReports)
		if len(allReports) > 0 {
			mapReports[filename] = allReports
		}
	}

	//chanReports <- 1
	if len(allReports) == 0 {
		if diff > 1 {
			_, _, date = api.GetQuarterByDate(date, diff-1)
		}
		reports, pages, _ := QuarterlyReports(date)
		if pages < 2 || len(reports) == 0 {
			return nil
		}
		allReports = append(allReports, reports...)
		for pageNo := 2; pageNo < pages+1; pageNo++ {
			list, pages, err := QuarterlyReports(date, pageNo)
			if err != nil || pages < 1 {
				logger.Error(err)
				break
			}
			count := len(list)
			if count == 0 {
				break
			}
			allReports = append(allReports, list...)
			if count < EastmoneyQuarterlyReportAllPageSize {
				break
			}
		}
		if len(allReports) > 0 {
			mapReports[filename] = allReports
			err := api.SlicesToCsv(filename, allReports)
			if err != nil {
				logger.Errorf("cache %s failed, error: %+v", filename, err)
			}
		}
	}

	for _, v := range allReports {
		if v.SecurityCode == securityCode {
			return &v
		}
	}

	return nil
}

// GetCacheQuarterlyReportsBySecurityCode 获取上市公司财务季报 Quarterly Reports
func GetCacheQuarterlyReportsBySecurityCode(securityCode, date string, diffQuarters ...int) *QuarterlyReport {
	diff := 1
	if len(diffQuarters) > 0 {
		diff = diffQuarters[0]
	}
	for ; diff < 4; diff++ {
		report := cacheQuarterlyReportsBySecurityCode(securityCode, date, diff)
		if report == nil {
			continue
		}
		return report
	}
	return nil
}
