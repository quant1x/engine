package dfcf

import (
	"encoding/json"
	"fmt"
	"gitee.com/quant1x/gotdx/proto"
	"gitee.com/quant1x/gotdx/trading"
	"gitee.com/quant1x/gox/api"
	"gitee.com/quant1x/gox/fastjson"
	"gitee.com/quant1x/gox/http"
	"gitee.com/quant1x/gox/logger"
	urlpkg "net/url"
	"strings"
)

const (
	kErrorCapitalBase = 90000 // 股本异常错误码基础值

	kUrlEastmoneyGdfxHoldingAnalyse     = "https://datacenter-web.eastmoney.com/api/data/v1/get"
	EastmoneyGdfxHoldingAnalysePageSize = 500
)

type HoldNumChangeState = int

const (
	HoldNumDampened       HoldNumChangeState = -1 // 减少
	HoldNumUnChanged      HoldNumChangeState = 0  // 不变
	HoldNumNewlyAdded     HoldNumChangeState = 1  // 新进/新增
	HoldNumIncrease       HoldNumChangeState = 2  // 增加
	HoldNumUnknownChanges HoldNumChangeState = -9 // 未知变化
)

// CirculatingShareholder Top10CirculatingShareholders
type CirculatingShareholder struct {
	SecurityCode     string  `dataframe:"security_code"`       // 证券代码
	SecurityName     string  `dataframe:"security_name"`       // 证券名称
	EndDate          string  `dataframe:"end_date"`            // 报告日期
	UpdateDate       string  `dataframe:"update_date"`         // 更新日期
	HolderType       string  `dataframe:"holder_type"`         // 股东类型
	HolderName       string  `dataframe:"holder_name"`         // 股东名称
	IsHoldOrg        string  `dataframe:"is_holdorg"`          // 股东是否机构
	HolderRank       int     `dataframe:"holder_rank"`         // 股东排名
	HoldNum          int     `dataframe:"hold_num"`            // 期末持股-数量
	FreeHoldNumRatio float64 `dataframe:"free_hold_num_ratio"` // 期末持股-比例
	HoldNumChange    int     `dataframe:"hold_num_change"`     // 期末持股-持股变动
	HoldChangeName   string  `dataframe:"change_name"`         // 期末持股-变化状态
	HoldChangeState  int     `dataframe:"change_state"`        // 期末持股-变化状态
	HoldChangeRatio  float64 `dataframe:"change_ratio"`        // 期末持股-持股变化比例
	HoldRatio        float64 `dataframe:"hold_ratio"`          // 期末持股-持股变动
	HoldRatioChange  float64 `dataframe:"hold_ratio_change"`   // 期末持股-数量变化比例
}

// HoldingAnalyse 持股分析
type HoldingAnalyse struct {
	SECUCODE                string  `json:"SECUCODE"`                // "股票代码",
	SECURITY_NAME           string  `json:"SECURITY_NAME_ABBR"`      // "股票简称",
	END_DATE                string  `json:"END_DATE"`                // "报告期",
	UPDATE_DATE             string  `json:"UPDATE_DATE"`             // "公告日",
	HOLDER_TYPE             string  `json:"HOLDER_TYPE"`             // "股东类型",
	HOLDER_NEWTYPE          string  `json:"HOLDER_NEWTYPE"`          // "-",
	HOLDER_NAME             string  `json:"HOLDER_NAME"`             // "股东名称",
	IS_HOLDORG              string  `json:"IS_HOLDORG"`              // "是否机构",
	HOLDER_RANK             int     `json:"HOLDER_RANK"`             // "股东排名",
	HOLD_NUM                int     `json:"HOLD_NUM"`                // "期末持股-数量",
	FREE_HOLDNUM_RATIO      float64 `json:"FREE_HOLDNUM_RATIO"`      // "-",
	HOLD_NUM_CHANGE_NAME    string  `json:"HOLDNUM_CHANGE_NAME"`     // "-",
	XZCHANGE                int     `json:"XZCHANGE"`                // "期末持股-数量变化",
	CHANGE_RATIO            float64 `json:"CHANGE_RATIO"`            // "-",
	HOLDER_STATE            string  `json:"HOLDER_STATE"`            // "-",
	REPORT_DATE_NAME        string  `json:"REPORT_DATE_NAME"`        // ---华丽的分割线---
	HOLDER_MARKET_CAP       float64 `json:"HOLDER_MARKET_CAP"`       // "期末持股-流通市值",
	HOLD_RATIO              float64 `json:"HOLD_RATIO"`              // "-",
	SECURITY_CODE           string  `json:"SECURITY_CODE"`           // "股票代码简写",
	HOLD_CHANGE             string  `json:"HOLD_CHANGE"`             // "-",
	HOLD_RATIO_CHANGE       float64 `json:"HOLD_RATIO_CHANGE"`       // "期末持股-数量变化比例",
	ORG_CODE                string  `json:"ORG_CODE"`                // "-",
	HOLDER_CODE             string  `json:"HOLDER_CODE"`             // "-",
	SECURITY_TYPE_CODE      string  `json:"SECURITY_TYPE_CODE"`      // "-",
	SHARES_TYPE             string  `json:"SHARES_TYPE"`             // "-",
	HOLDER_NEW              string  `json:"HOLDER_NEW"`              // "-",
	FREE_RATIO_QOQ          string  `json:"FREE_RATIO_QOQ"`          // "-",
	HOLDER_STATEE           string  `json:"HOLDER_STATEE"`           // "-",
	IS_REPORT               string  `json:"IS_REPORT"`               // "-",
	HOLDER_CODE_OLD         string  `json:"HOLDER_CODE_OLD"`         // "-",
	IS_MAX_REPORT_DATE      string  `json:"IS_MAX_REPORTDATE"`       // "-",
	COOPERATION_HOLDER_MARK string  `json:"COOPERATION_HOLDER_MARK"` // "-",
	MXID                    string  `json:"MXID"`                    // "-",
	LISTING_STATE           string  `json:"LISTING_STATE"`           // "-",
	NEW_CHANGE_RATIO        string  `json:"NEW_CHANGE_RATIO"`        // "-",
	HOLD_NUM_CHANGE         string  `json:"HOLD_NUM_CHANGE"`         // "期末持股-持股变动",
}

// FreeHoldingAnalyse 东方财富网-数据中心-股东分析-股东持股明细-十大流通股东
//
//	https://data.eastmoney.com/gdfx/HoldingAnalyse.html
func getFreeHoldingAnalyse(pageNumber ...int) ([]HoldingAnalyse, int, error) {
	pageNo := 1
	pages := 0
	if len(pageNumber) > 0 {
		pageNo = pageNumber[0]
	}
	pageSize := EastmoneyGdfxHoldingAnalysePageSize
	_, qEnd := api.GetQuarterDay(9)
	endDate := trading.FixTradeDate(qEnd)
	params := urlpkg.Values{
		"sortColumns": {"UPDATE_DATE,SECURITY_CODE,HOLDER_RANK"},
		"sortTypes":   {"-1,1,1"},
		"pageSize":    {fmt.Sprintf("%d", pageSize)},
		"pageNumber":  {fmt.Sprintf("%d", pageNo)},
		"reportName":  {"RPT_F10_EH_FREEHOLDERS"},
		"columns":     {"ALL"},
		"source":      {"WEB"},
		"client":      {"WEB"},
		"filter":      {fmt.Sprintf("(END_DATE>='%s')", endDate)},
	}
	url := kUrlEastmoneyGdfxHoldingAnalyse + "?" + params.Encode()
	data, err := http.HttpGet(url)
	//fmt.Println(api.Bytes2String(data))
	var holds = []HoldingAnalyse{}
	obj, err := fastjson.ParseBytes(data)
	if err != nil {
		logger.Errorf("%+v\n", err)
		return holds, pages, err
	}
	result := obj.Get("result")
	if result == nil {
		return holds, pages, nil
	}
	pages = result.GetInt("pages")
	tData := result.Get("data")
	if tData == nil {
		return holds, pages, nil
	}
	text := tData.String()
	err = json.Unmarshal(api.String2Bytes(text), &holds)
	for i := 0; i < len(holds); i++ {
		holds[i].END_DATE = trading.FixTradeDate(holds[i].END_DATE)
		holds[i].UPDATE_DATE = trading.FixTradeDate(holds[i].UPDATE_DATE)
	}
	return holds, pages, err
}

func FreeHoldingAnalyse(pageNumber ...int) ([]CirculatingShareholder, int, error) {
	pageNo := 1
	pages := 0
	if len(pageNumber) > 0 {
		pageNo = pageNumber[0]
	}
	shareholders := []CirculatingShareholder{}
	pageSize := EastmoneyGdfxHoldingAnalysePageSize
	_, qEnd := api.GetQuarterDay(9)
	endDate := trading.FixTradeDate(qEnd)
	params := urlpkg.Values{
		"sortColumns": {"UPDATE_DATE,SECURITY_CODE,HOLDER_RANK"},
		"sortTypes":   {"-1,1,1"},
		"pageSize":    {fmt.Sprintf("%d", pageSize)},
		"pageNumber":  {fmt.Sprintf("%d", pageNo)},
		"reportName":  {"RPT_F10_EH_FREEHOLDERS"},
		"columns":     {"ALL"},
		"source":      {"WEB"},
		"client":      {"WEB"},
		"filter":      {fmt.Sprintf("(END_DATE>='%s')", endDate)},
	}
	url := kUrlEastmoneyGdfxHoldingAnalyse + "?" + params.Encode()
	data, err := http.HttpGet(url)
	var holds = []HoldingAnalyse{}
	obj, err := fastjson.ParseBytes(data)
	if err != nil {
		logger.Errorf("%+v\n", err)
		return shareholders, pages, err
	}
	result := obj.Get("result")
	if result == nil {
		return shareholders, pages, nil
	}
	pages = result.GetInt("pages")
	tData := result.Get("data")
	if tData == nil {
		return shareholders, pages, nil
	}
	text := tData.String()
	err = json.Unmarshal(api.String2Bytes(text), &holds)
	for i := 0; i < len(holds); i++ {
		holds[i].END_DATE = trading.FixTradeDate(holds[i].END_DATE)
		holds[i].UPDATE_DATE = trading.FixTradeDate(holds[i].UPDATE_DATE)
	}
	for _, v := range holds {
		shareholder := CirculatingShareholder{
			//SecurityCode   string `dataframe:"security_code"`    // 证券代码
			SecurityCode: v.SECUCODE,
			//SecurityName   string `dataframe:"security_name"`    // 证券名称
			SecurityName: v.SECURITY_NAME,
			//EndDate        string `dataframe:"end_date"`         // 报告日期
			EndDate: v.END_DATE,
			//UpdateDate     string `dataframe:"update_date"`      // 更新日期
			UpdateDate: v.UPDATE_DATE,
			//HolderType     string `dataframe:"holder_type"`      // 股东类型
			HolderType: v.HOLDER_NEWTYPE,
			//HolderName     string `dataframe:"holder_name"`      // 股东名称
			HolderName: v.HOLDER_NAME,
			//IsHolderOrg    string `dataframe:"is_holder_org"`    // 股东是否机构
			IsHoldOrg: v.IS_HOLDORG,
			//HolderRank     int    `dataframe:"holder_rank"`      // 股东排名
			HolderRank: v.HOLDER_RANK,
			//HoldNum        int    `dataframe:"hold_num"`         // 期末持股-数量
			HoldNum: v.HOLD_NUM,
			//FreeHoldNumRatio float64 `dataframe:"hold_num_ratio"`  // 期末持股-比例
			FreeHoldNumRatio: v.FREE_HOLDNUM_RATIO,
			//HoldNumChange  string `dataframe:"hold_num_change"`  // 期末持股-持股变动
			HoldNumChange: v.XZCHANGE,
			//HoldChangeName string `dataframe:"hold_change_name"` // 期末持股-变化状态
			HoldChangeName: v.HOLD_NUM_CHANGE_NAME,
			//HoldChangeRatio  string  `dataframe:"change_ratio"`    // 期末持股-变化比例
			HoldChangeRatio: v.CHANGE_RATIO,
			//HoldRatio        float64 `dataframe:"hold_ratio"`          // 期末持股-持股变动
			HoldRatio: v.HOLD_RATIO,
			//HoldRatioChange  float64 `dataframe:"hold_ratio_change"`   // "期末持股-数量变化比例",
			HoldRatioChange: v.HOLD_RATIO_CHANGE,
		}
		// 修订证券代码
		_, mfalg, mcode := proto.DetectMarket(shareholder.SecurityCode)
		shareholder.SecurityCode = mfalg + mcode
		//HoldChangeState  int     `dataframe:"change_state"`        // 期末持股-变化状态
		switch v.HOLD_NUM_CHANGE_NAME {
		case "新进":
			shareholder.HoldChangeState = HoldNumNewlyAdded
		case "增加":
			shareholder.HoldChangeState = HoldNumIncrease
		case "减少":
			shareholder.HoldChangeState = HoldNumDampened
		case "不变":
			shareholder.HoldChangeState = HoldNumUnChanged
		default: // 未知变化报警
			shareholder.HoldChangeState = HoldNumUnknownChanges
			warning := fmt.Sprintf("%s: %s, 变化状态未知: %s", v.SECURITY_NAME, v.SECUCODE, v.HOLD_NUM_CHANGE_NAME)
			logger.Warnf(warning)
		}
		shareholders = append(shareholders, shareholder)
	}
	return shareholders, pages, err
}

// FreeHoldingDetail 拉取近期的
func freeHoldingDetail() []CirculatingShareholder {
	pageNo := 1
	holds := []CirculatingShareholder{}
	for {
		list, pages, err := getFreeHoldingAnalyse(pageNo)
		if err != nil || pages < 1 {
			logger.Error(err)
			break
		}
		count := len(list)
		if count == 0 {
			break
		}
		for _, v := range list {
			shareholder := CirculatingShareholder{
				//SecurityCode   string `dataframe:"security_code"`    // 证券代码
				SecurityCode: strings.TrimSpace(v.SECUCODE),
				//SecurityName   string `dataframe:"security_name"`    // 证券名称
				SecurityName: strings.TrimSpace(v.SECURITY_NAME),
				//EndDate        string `dataframe:"end_date"`         // 报告日期
				EndDate: strings.TrimSpace(v.END_DATE),
				//UpdateDate     string `dataframe:"update_date"`      // 更新日期
				UpdateDate: strings.TrimSpace(v.UPDATE_DATE),
				//HolderType     string `dataframe:"holder_type"`      // 股东类型
				HolderType: v.HOLDER_NEWTYPE,
				//HolderName     string `dataframe:"holder_name"`      // 股东名称
				HolderName: v.HOLDER_NAME,
				//IsHolderOrg    string `dataframe:"is_holder_org"`    // 股东是否机构
				IsHoldOrg: v.IS_HOLDORG,
				//HolderRank     int    `dataframe:"holder_rank"`      // 股东排名
				HolderRank: v.HOLDER_RANK,
				//HoldNum        int    `dataframe:"hold_num"`         // 期末持股-数量
				HoldNum: v.HOLD_NUM,
				//FreeHoldNumRatio float64 `dataframe:"hold_num_ratio"`  // 期末持股-比例
				FreeHoldNumRatio: v.FREE_HOLDNUM_RATIO,
				//HoldNumChange  string `dataframe:"hold_num_change"`  // 期末持股-持股变动
				HoldNumChange: v.XZCHANGE,
				//HoldChangeName string `dataframe:"hold_change_name"` // 期末持股-变化状态
				HoldChangeName: v.HOLD_NUM_CHANGE_NAME,
				//HoldChangeRatio  string  `dataframe:"change_ratio"`    // 期末持股-变化比例
				HoldChangeRatio: v.CHANGE_RATIO,
				//HoldRatio        float64 `dataframe:"hold_ratio"`          // 期末持股-持股变动
				HoldRatio: v.HOLD_RATIO,
				//HoldRatioChange  float64 `dataframe:"hold_ratio_change"`   // "期末持股-数量变化比例",
				HoldRatioChange: v.HOLD_RATIO_CHANGE,
			}
			// 修订证券代码
			_, mfalg, mcode := proto.DetectMarket(shareholder.SecurityCode)
			shareholder.SecurityCode = mfalg + mcode
			//HoldChangeState  int     `dataframe:"change_state"`        // 期末持股-变化状态
			switch v.HOLD_NUM_CHANGE_NAME {
			case "新进":
				shareholder.HoldChangeState = HoldNumNewlyAdded
			case "增加":
				shareholder.HoldChangeState = HoldNumIncrease
			case "减少":
				shareholder.HoldChangeState = HoldNumDampened
			case "不变":
				shareholder.HoldChangeState = HoldNumUnChanged
			default: // 未知变化报警
				shareholder.HoldChangeState = HoldNumUnknownChanges
				warning := fmt.Sprintf("%s: %s, 变化状态未知: %s", v.SECURITY_NAME, v.SECUCODE, v.HOLD_NUM_CHANGE_NAME)
				logger.Warnf(warning)
			}
			holds = append(holds, shareholder)
		}
		if count < EastmoneyGdfxHoldingAnalysePageSize {
			break
		}
		pageNo += 1
		break
	}
	return holds
}
