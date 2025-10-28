package dfcf

import (
	"fmt"
	urlpkg "net/url"

	"gitee.com/quant1x/data/exchange"
	"gitee.com/quant1x/gox/api"
	"gitee.com/quant1x/gox/http"
	"gitee.com/quant1x/gox/logger"
	"gitee.com/quant1x/pkg/fastjson"
)

// 数据来源: https://data.eastmoney.com/bbsj/yjbb/301381.html
const (
	urlQuarterlyReportAll               = "https://datacenter-web.eastmoney.com/api/data/v1/get"
	EastmoneyQuarterlyReportAllPageSize = 50 // 一页最大50条
)

// QuarterlyReport 财报
type QuarterlyReport struct {
	SecurityCode       string  `json:"CODE"`                 // 证券代码
	ReportDate         string  `json:"REPORTDATE"`           // 报告日期
	NoticeDate         string  `json:"NOTICE_DATE"`          // 最新公告日期
	UpdateDate         string  `json:"UPDATE_DATE"`          // 更新日期
	SecuCode           string  `json:"SECUCODE"`             // 证券代码
	BasicEPS           float64 `json:"BASIC_EPS"`            // 每股收益
	DeductBasicEPS     float64 `json:"DEDUCT_BASIC_EPS"`     // 每股收益(扣除)
	TotalOperateIncome float64 `json:"TOTAL_OPERATE_INCOME"` // 营业总收入
	ParentNetprofit    float64 `json:"PARENT_NETPROFIT"`     // 净利润
	WeightAvgRoe       float64 `json:"WEIGHTAVG_ROE"`        // 净资产收益率
	YSTZ               float64 `json:"YSTZ"`                 // 营业总收入同比增长
	SJLTZ              float64 `json:"SJLTZ"`                // 净利润同比增长
	BPS                float64 `json:"BPS"`                  // 每股净资产
	MGJYXJJE           float64 `json:"MGJYXJJE"`             // 每股经营现金流量(元)
	XSMLL              float64 `json:"XSMLL"`                // 销售毛利率(%)
	YSHZ               float64 `json:"YSHZ"`
	SJLHZ              float64 `json:"SJLHZ"`
	ASSIGNDSCRPT       float64 `json:"ASSIGNDSCRPT"`
	PAYYEAR            float64 `json:"PAYYEAR"`
	PUBLISHNAME        string  `json:"PUBLISHNAME"`
	ZXGXL              float64 `json:"ZXGXL"`
	ORGCODE            string  `json:"ORG_CODE"`
	TRADEMARKETZJG     string  `json:"TRADE_MARKET_ZJG"`
	IsNew              string  `json:"ISNEW"`
	QDATE              string  `name:"报告期" json:"QDATE"` // 季报期
	DATATYPE           string  `json:"DATATYPE"`
	DATAYEAR           string  `json:"DATAYEAR"`
	DATEMMDD           string  `json:"DATEMMDD"`
	EITIME             string  `json:"EITIME"`
	TRADEMARKETCODE    string  `json:"TRADE_MARKET_CODE"`
	TRADEMARKET        string  `json:"TRADE_MARKET"` //市场
	SECURITYTYPECODE   string  `json:"SECURITY_TYPE_CODE"`
	SECURITYTYPE       string  `json:"SECURITY_TYPE"`
	SECURITYCODE       string  `json:"SECURITY_CODE"`      // 证券代码
	SECURITYNAMEABBR   string  `json:"SECURITY_NAME_ABBR"` // 证券名称
}

// GetQuarterlyReports 分页获取季报数据
func GetQuarterlyReports(pageNumber ...int) (reports []QuarterlyReport, pages int, err error) {
	pageNo := 1
	if len(pageNumber) > 0 {
		pageNo = pageNumber[0]
	}
	_, qEnd := api.GetQuarterDay(4)
	beginDate := exchange.FixTradeDate(qEnd)
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
		"filter": {fmt.Sprintf("(REPORTDATE>='%s')", beginDate)},
		//"filter": {fmt.Sprintf("(REPORTDATE<='%s')(SECURITY_CODE=\"301381\")", beginDate)},
		//"filter": {fmt.Sprintf("(REPORTDATE>='%s')(REPORTDATE<'%s')", beginDate, "2023-03-31")},
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
	pages = result.GetInt("pages")
	if len(list) > 0 {
		for _, v := range list {
			v.GetStringBytes()
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
