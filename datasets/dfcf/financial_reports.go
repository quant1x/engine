package dfcf

import (
	"fmt"
	"gitee.com/quant1x/gotdx/proto"
	"gitee.com/quant1x/gotdx/trading"
	"gitee.com/quant1x/gox/api"
	"gitee.com/quant1x/gox/fastjson"
	"gitee.com/quant1x/gox/http"
	"gitee.com/quant1x/gox/logger"
	urlpkg "net/url"
)

// QuarterlyReports 分页获取季报数据
func QuarterlyReports(featureData string, pageNumber ...int) (reports []QuarterlyReport, pages int, err error) {
	pageNo := 1
	if len(pageNumber) > 0 {
		pageNo = pageNumber[0]
	}
	//_, qEnd := api.GetQuarterDay(4)
	//beginDate := trading.FixTradeDate(qEnd)
	//beginDate = "2022-12-31"
	qBegin, qEnd := api.GetQuarterDayByDate(featureData)
	quarterBeginDate := trading.FixTradeDate(qBegin)
	quarterEndDate := trading.FixTradeDate(qEnd)
	params := urlpkg.Values{
		//"callback":    {"jQuery1123043614175387302234_1685785566671"},
		//"sortColumns": {"UPDATE_DATE,SECURITY_CODE"},
		"sortColumns": {"REPORTDATE,SECURITY_CODE"},
		"sortTypes":   {"-1,1"},
		"pageSize":    {fmt.Sprint(EastmoneyQuarterlyReportAllPageSize)},
		"pageNumber":  {fmt.Sprintf("%d", pageNo)},
		"reportName":  {"RPT_LICO_FN_CPD"},
		"columns":     {"ALL"},
		//"filter":      {"(REPORTDATE>='2023-03-31')"},
		"filter": {fmt.Sprintf("(REPORTDATE='%s')", quarterEndDate)},
		//"filter": {fmt.Sprintf("(REPORTDATE<='%s')(SECURITY_CODE=\"301381\")", beginDate)},
		//"filter": {fmt.Sprintf("(REPORTDATE>='%s')(REPORTDATE<'%s')", quarterBeginDate, quarterEndDate)},
	}
	_ = quarterBeginDate
	_ = quarterEndDate

	url := urlQuarterlyReportAll + "?" + params.Encode()
	data, err := http.HttpGet(url)
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
				DedtctBasicEPS:     v.GetFloat64("DEDUCT_BASIC_EPS"),
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
				TOTALOPERATEINCOME: v.GetFloat64("TOTAL_OPERATE_INCOME"),
				PARENTNETPROFIT:    v.GetFloat64("PARENT_NETPROFIT"),
				WEIGHTAVGROE:       v.GetFloat64("WEIGHTAVG_ROE"),
				YSTZ:               v.GetFloat64("YSTZ"),
				SJLTZ:              v.GetFloat64("SJLTZ"),
				MGJYXJJE:           v.GetFloat64("MGJYXJJE"),
				XSMLL:              v.GetFloat64("XSMLL"),
				YSHZ:               v.GetFloat64("YSHZ"),
				SJLHZ:              v.GetFloat64("SJLHZ"),
				ASSIGNDSCRPT:       v.GetFloat64("ASSIGNDSCRPT"),
				PAYYEAR:            v.GetFloat64("PAYYEAR"),
				PUBLISHNAME:        v.GetFloat64("PUBLISHNAME"),
				ZXGXL:              v.GetFloat64("ZXGXL"),
			}
			// 截取市场编码，截取股票编码，市场编码+股票编码拼接作为主键
			securityCode := proto.CorrectSecurityCode(report.SecuCode)
			report.SecurityCode = securityCode
			reports = append(reports, report)
		}
	}
	return
}

//// QuarterlyReports 分页获取季报数据
//func QuarterlyReports(securityCode, beginDate, endDate string, diffQuarters int, pageNumber ...int) (reports []QuarterlyReport, pages int, err error) {
//	pageNo := 1
//	if len(pageNumber) > 0 {
//		pageNo = pageNumber[0]
//	}
//	diff := 1
//	if len(diffQuarters) > 0 {
//		diff = diffQuarters[0]
//	}
//	_, _, last := api.GetQuarterByDate(date, diff)
//	_, qEnd := api.GetQuarterDay(4)
//	api.GetQuarterDayByDate()
//	beginDate := trading.FixTradeDate(qEnd)
//	beginDate = "2022-12-31"
//	params := urlpkg.Values{
//		//"callback":    {"jQuery1123043614175387302234_1685785566671"},
//		//"sortColumns": {"UPDATE_DATE,SECURITY_CODE"},
//		"sortColumns": {"REPORTDATE,SECURITY_CODE"},
//		"sortTypes":   {"-1,1"},
//		"pageSize":    {fmt.Sprint(EastmoneyQuarterlyReportAllPageSize)},
//		"pageNumber":  {fmt.Sprintf("%d", pageNo)},
//		"reportName":  {"RPT_LICO_FN_CPD"},
//		"columns":     {"ALL"},
//		//"filter":      {"(REPORTDATE>='2023-03-31')"},
//		"filter": {fmt.Sprintf("(REPORTDATE>='%s')", beginDate)},
//		//"filter": {fmt.Sprintf("(REPORTDATE<='%s')(SECURITY_CODE=\"301381\")", beginDate)},
//		//"filter": {fmt.Sprintf("(REPORTDATE>='%s')(REPORTDATE<'%s')", beginDate, "2023-03-31")},
//	}
//
//	url := urlQuarterlyReportAll + "?" + params.Encode()
//	data, err := http.HttpGet(url)
//	//fmt.Println(api.Bytes2String(data))
//	obj, err := fastjson.ParseBytes(data)
//	if err != nil {
//		logger.Errorf("%+v\n", err)
//		return
//	}
//
//	result := obj.Get("result")
//	list := result.GetArray("data")
//	pages = result.GetInt("pages")
//	if len(list) > 0 {
//		for _, v := range list {
//			report := QuarterlyReport{
//				SecuCode:           v.GetString("SECUCODE"),
//				UpdateDate:         v.GetString("UPDATE_DATE"),
//				ReportDate:         v.GetString("REPORTDATE"),
//				BasicEPS:           v.GetFloat64("BASIC_EPS"),
//				DedtctBasicEPS:     v.GetFloat64("DEDUCT_BASIC_EPS"),
//				BPS:                v.GetFloat64("BPS"),
//				NoticeDate:         v.GetString("NOTICE_DATE"),
//				IsNew:              v.GetString("ISNEW"),
//				ORGCODE:            v.GetString("ORG_CODE"),
//				TRADEMARKETZJG:     v.GetString("TRADE_MARKET_ZJG"),
//				QDATE:              v.GetString("QDATE"),
//				DATATYPE:           v.GetString("DATATYPE"),
//				DATAYEAR:           v.GetString("DATAYEAR"),
//				DATEMMDD:           v.GetString("DATEMMDD"),
//				EITIME:             v.GetString("EITIME"),
//				SECURITYCODE:       v.GetString("SECURITY_CODE"),
//				SECURITYNAMEABBR:   v.GetString("SECURITY_NAME_ABBR"),
//				TRADEMARKETCODE:    v.GetString("TRADE_MARKET_CODE"),
//				TRADEMARKET:        v.GetString("TRADE_MARKET"),
//				SECURITYTYPECODE:   v.GetString("SECURITY_TYPE_CODE"),
//				SECURITYTYPE:       v.GetString("SECURITY_TYPE"),
//				TOTALOPERATEINCOME: v.GetFloat64("TOTAL_OPERATE_INCOME"),
//				PARENTNETPROFIT:    v.GetFloat64("PARENT_NETPROFIT"),
//				WEIGHTAVGROE:       v.GetFloat64("WEIGHTAVG_ROE"),
//				YSTZ:               v.GetFloat64("YSTZ"),
//				SJLTZ:              v.GetFloat64("SJLTZ"),
//				MGJYXJJE:           v.GetFloat64("MGJYXJJE"),
//				XSMLL:              v.GetFloat64("XSMLL"),
//				YSHZ:               v.GetFloat64("YSHZ"),
//				SJLHZ:              v.GetFloat64("SJLHZ"),
//				ASSIGNDSCRPT:       v.GetFloat64("ASSIGNDSCRPT"),
//				PAYYEAR:            v.GetFloat64("PAYYEAR"),
//				PUBLISHNAME:        v.GetFloat64("PUBLISHNAME"),
//				ZXGXL:              v.GetFloat64("ZXGXL"),
//			}
//			// 截取市场编码，截取股票编码，市场编码+股票编码拼接作为主键
//			securityCode := proto.CorrectSecurityCode(report.SecuCode)
//			report.SecurityCode = securityCode
//			reports = append(reports, report)
//		}
//	}
//	return
//}
////
////// GetCacheQuarterlyReports (GetCacheFinancialReports) 获取上市公司财务季报 Quarterly Reports
////func GetCacheQuarterlyReports(securityCode, date string, diffQuarters ...int) (list []QuarterlyReport) {
////	diff := 1
////	if len(diffQuarters) > 0 {
////		diff = diffQuarters[0]
////	}
////	_, _, last := api.GetQuarterByDate(date, diff)
////	filename := cache.Top10HoldersFilename(securityCode, last)
////	if api.FileExist(filename) {
////		err := api.CsvToSlices(filename, &list)
////		if err == nil && len(list) > 0 {
////			return
////		}
////	}
////	tmpList := QuarterlyReports(securityCode, last)
////	if len(tmpList) > 0 {
////		list = tmpList
////	}
////	if len(list) > 0 {
////		_ = api.SlicesToCsv(filename, list)
////	}
////	return
////}
