package dfcf

import (
	"encoding/json"
	"fmt"
	"gitee.com/quant1x/gotdx/proto"
	"gitee.com/quant1x/gotdx/trading"
	"gitee.com/quant1x/gox/http"
	urlpkg "net/url"
	"strings"
)

const (
	urlFinanceReports               = "https://datacenter-web.eastmoney.com/api/data/v1/get"
	EastmoneyFinanceReportsPageSize = 100
)

// PreviewQuarterlyReport 业绩预报详情
type PreviewQuarterlyReport struct {
	SecurityCode        string  `name:"证券代码" json:"SECURITY_CODE"`
	SecurityName        string  `name:"证券名称" json:"SECURITY_NAME"`
	OrgCode             string  `name:"机构代码" json:"ORG_CODE"`
	NoticeDate          string  `name:"公告日期" json:"NOTICE_DATE"`
	ReportDate          string  `name:"报告日期" json:"REPORT_DATE"`
	PredictFinanceCode  string  `name:"预告代码" json:"PREDICT_FINANCE_CODE"`
	PredictFinance      string  `name:"预告财报" json:"PREDICT_FINANCE"`
	PredictAmtLower     float64 `name:"预计营收下限" json:"PREDICT_AMT_LOWER"`
	PredictAmtUpper     float64 `name:"预计营收上限" json:"PREDICT_AMT_UPPER"`
	AddAmpLower         float64 `name:"增长下限" json:"ADD_AMP_LOWER"`
	AddAmpUpper         float64 `name:"增长下限" json:"ADD_AMP_UPPER"`
	PredictContent      string  `name:"预告内容" json:"PREDICT_CONTENT"`
	ChangeReasonExplain string  `name:"改变原因" json:"CHANGE_REASON_EXPLAIN"`
	PredictType         string  `name:"预计类型" json:"PREDICT_TYPE"`
	PreyearSamePeriod   float64 `name:"上年同期" json:"PREYEAR_SAME_PERIOD"`
	TradeMarket         string  `name:"所在交易所" json:"TRADE_MARKET"`
	TradeMarketCode     string  `name:"场内交易代码" json:"TRADE_MARKET_CODE"`
	SecurityType        string  `name:"证券类型" json:"SECURITY_TYPE"`
	SecurityTypeCode    string  `name:"证券类型代码" json:"SECURITY_TYPE_CODE"`
	IncreaseJz          float64 `name:"增长均值" json:"INCREASE_JZ"`
	ForecastJz          float64 `name:"预测均值" json:"FORECAST_JZ"`
	ForecastState       string  `name:"预测状态" json:"FORECAST_STATE"`
	IsLatest            string  `name:"是否最新财报" json:"IS_LATEST"`
	PredictRatioLower   float64 `name:"预测增速下限" json:"PREDICT_RATIO_LOWER"`
	PredictRatioUpper   float64 `name:"预测增速上限" json:"PREDICT_RATIO_UPPER"`
	PredictHbmean       float64 `name:"预计每股盈利?" json:"PREDICT_HBMEAN"`
}

func (f PreviewQuarterlyReport) GetDate() string {
	return f.NoticeDate
}

func (f PreviewQuarterlyReport) GetSecurityCode() string {
	return f.SecurityCode
}

// FinanceReports 获取哪天开始的财报数据
// https://data.eastmoney.com/bbsj/202303/yjyg.html?type=increase
// https://data.eastmoney.com/bbsj/202303/yjbb.html
func FinanceReports(date string, pageNumber ...int) (reports []PreviewQuarterlyReport, pages, originalRecords int, err error) {
	pageNo := 1
	pageSize := EastmoneyFinanceReportsPageSize
	if len(pageNumber) > 0 {
		pageNo = pageNumber[0]
	}
	beginDate := trading.FixTradeDate(date)
	params := urlpkg.Values{
		"page_size":   {fmt.Sprintf("%d", pageSize)},
		"page_index":  {fmt.Sprintf("%d", pageNo)},
		"sortColumns": {"NOTICE_DATE", "SECURITY_CODE"},
		"sortTypes":   {"-1", "-1"},
		"reportName":  {"RPT_PUBLIC_OP_NEWPREDICT"},
		"columns":     {"ALL"},
		"filter":      {fmt.Sprintf("(REPORT_DATE>='%s')", beginDate)},
		//"callback": {"jQuery1123017378221241641723_1683879528446"},
	}
	url := urlFinanceReports + "?" + params.Encode()
	data, err := http.Get(url)
	if err != nil {
		return
	}
	//fmt.Println(api.Bytes2String(data))
	var raw rawFinanceReport
	err = json.Unmarshal(data, &raw)
	if err != nil {
		return
	}
	if !raw.Success || raw.Code != 0 || len(raw.Result.Data) == 0 {
		err = ErrNoticeNotFound
		return
	}
	//pages = functions.GetPages(pageSize, raw.Result.Count)
	pages = raw.Result.Pages
	originalRecords = len(raw.Result.Data)
	for _, v := range raw.Result.Data {
		securityName := strings.TrimSpace(v.SECURITY_NAME_ABBR)
		if len(v.SECUCODE) == 0 || len(securityName) == 0 {
			continue
		}
		code := strings.TrimSpace(v.SECUCODE)
		_, mflag, msymbol := proto.DetectMarket(code)
		securityCode := mflag + msymbol

		report := PreviewQuarterlyReport{
			SecurityCode: securityCode,
			SecurityName: securityName,
			ReportDate:   v.REPORT_DATE,
			NoticeDate:   v.NOTICE_DATE,
			OrgCode:      v.ORG_CODE,
			//PredictFinanceCode  string  `name:"预告代码" json:"PREDICT_FINANCE_CODE"`
			PredictFinanceCode: v.PREDICT_FINANCE_CODE,
			//PredictFinance      string  `name:"预告财报" json:"PREDICT_FINANCE"`
			PredictFinance: v.PREDICT_FINANCE,
			//PredictAmtLower     float64 `name:"预计营收下限" json:"PREDICT_AMT_LOWER"`
			PredictAmtLower: v.PREDICT_AMT_LOWER,
			//PredictAmtUpper     float64 `name:"预计营收上限" json:"PREDICT_AMT_UPPER"`
			PredictAmtUpper: v.PREDICT_AMT_UPPER,
			//AddAmpLower         float64 `name:"增长下限" json:"ADD_AMP_LOWER"`
			AddAmpLower: v.ADD_AMP_LOWER,
			//AddAmpUpper         float64 `name:"增长下限" json:"ADD_AMP_UPPER"`
			AddAmpUpper: v.ADD_AMP_UPPER,
			//PredictContent      string  `name:"预告内容" json:"PREDICT_CONTENT"`
			PredictContent: v.PREDICT_CONTENT,
			//ChangeReasonExplain string  `name:"改变原因" json:"CHANGE_REASON_EXPLAIN"`
			ChangeReasonExplain: v.CHANGE_REASON_EXPLAIN,
			//PredictType         string  `name:"预计类型" json:"PREDICT_TYPE"`
			PredictType: v.PREDICT_TYPE,
			//PreyearSamePeriod   float64 `name:"上年同期" json:"PREYEAR_SAME_PERIOD"`
			PreyearSamePeriod: v.PREYEAR_SAME_PERIOD,
			//TradeMarket         string  `name:"所在交易所" json:"TRADE_MARKET"`
			TradeMarket: v.TRADE_MARKET,
			//TradeMarketCode     string  `name:"场内交易代码" json:"TRADE_MARKET_CODE"`
			TradeMarketCode: v.TRADE_MARKET_CODE,
			//SecurityType        string  `name:"证券类型" json:"SECURITY_TYPE"`
			SecurityType: v.SECURITY_TYPE,
			//SecurityTypeCode    string  `name:"证券类型代码" json:"SECURITY_TYPE_CODE"`
			SecurityTypeCode: v.SECURITY_TYPE_CODE,
			//IncreaseJz          float64 `name:"增长均值" json:"INCREASE_JZ"`
			IncreaseJz: v.INCREASE_JZ,
			//ForecastJz          float64 `name:"预测均值" json:"FORECAST_JZ"`
			ForecastJz: v.FORECAST_JZ,
			//ForecastState       string  `name:"预测状态" json:"FORECAST_STATE"`
			ForecastState: v.FORECAST_STATE,
			//IsLatest            string  `name:"是否最新财报" json:"IS_LATEST"`
			IsLatest: v.IS_LATEST,
			//PredictRatioLower   float64 `name:"预测增速下限" json:"PREDICT_RATIO_LOWER"`
			PredictRatioLower: v.PREDICT_RATIO_LOWER,
			//PredictRatioUpper   float64 `name:"预测增速上限" json:"PREDICT_RATIO_UPPER"`
			PredictRatioUpper: v.PREDICT_RATIO_UPPER,
			//PredictHbmean       float64 `name:"预计每股盈利?" json:"PREDICT_HBMEAN"`
			PredictHbmean: v.PREDICT_HBMEAN,
		}
		reports = append(reports, report)
	}
	return
}

// 财报原始数据结构
type rawFinanceReport struct {
	Version string `json:"version"`
	Result  struct {
		Pages int `json:"pages"`
		Data  []struct {
			SECUCODE              string  `json:"SECUCODE"`
			SECURITY_CODE         string  `json:"SECURITY_CODE"`
			SECURITY_NAME_ABBR    string  `json:"SECURITY_NAME_ABBR"`
			ORG_CODE              string  `json:"ORG_CODE"`
			NOTICE_DATE           string  `json:"NOTICE_DATE"`
			REPORT_DATE           string  `json:"REPORT_DATE"`
			PREDICT_FINANCE_CODE  string  `json:"PREDICT_FINANCE_CODE"`
			PREDICT_FINANCE       string  `json:"PREDICT_FINANCE"`
			PREDICT_AMT_LOWER     float64 `json:"PREDICT_AMT_LOWER"`
			PREDICT_AMT_UPPER     float64 `json:"PREDICT_AMT_UPPER"`
			ADD_AMP_LOWER         float64 `json:"ADD_AMP_LOWER"`
			ADD_AMP_UPPER         float64 `json:"ADD_AMP_UPPER"`
			PREDICT_CONTENT       string  `json:"PREDICT_CONTENT"`
			CHANGE_REASON_EXPLAIN string  `json:"CHANGE_REASON_EXPLAIN"`
			PREDICT_TYPE          string  `json:"PREDICT_TYPE"`
			PREYEAR_SAME_PERIOD   float64 `json:"PREYEAR_SAME_PERIOD"`
			TRADE_MARKET          string  `json:"TRADE_MARKET"`
			TRADE_MARKET_CODE     string  `json:"TRADE_MARKET_CODE"`
			SECURITY_TYPE         string  `json:"SECURITY_TYPE"`
			SECURITY_TYPE_CODE    string  `json:"SECURITY_TYPE_CODE"`
			INCREASE_JZ           float64 `json:"INCREASE_JZ"`
			FORECAST_JZ           float64 `json:"FORECAST_JZ"`
			FORECAST_STATE        string  `json:"FORECAST_STATE"`
			IS_LATEST             string  `json:"IS_LATEST"`
			PREDICT_RATIO_LOWER   float64 `json:"PREDICT_RATIO_LOWER"`
			PREDICT_RATIO_UPPER   float64 `json:"PREDICT_RATIO_UPPER"`
			PREDICT_HBMEAN        float64 `json:"PREDICT_HBMEAN"`
		} `json:"data"`
		Count int `json:"count"`
	} `json:"result"`
	Success bool   `json:"success"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}
