package dfcf

import (
	"encoding/json"
	"fmt"
	"gitee.com/quant1x/engine/cache"
	"gitee.com/quant1x/gotdx/proto"
	"gitee.com/quant1x/gotdx/trading"
	"gitee.com/quant1x/gox/api"
	"gitee.com/quant1x/gox/http"
	"gitee.com/quant1x/gox/logger"
	urlpkg "net/url"
)

// 前十大流通股东 https://emweb.securities.eastmoney.com/PC_HSF10/ShareholderResearch/Index?type=web&code=sh600822#sdltgd-0
// 数据接口 https://emweb.securities.eastmoney.com/PC_HSF10/ShareholderResearch/PageAjax?code=SH600822
//
// 十大股东 https://datacenter-web.eastmoney.com/api/data/v1/get?callback=jQuery112308928698030561204_1687518793792&sortColumns=RANK&sortTypes=1&pageSize=10&pageNumber=1&reportName=RPT_DMSK_HOLDERS&columns=ALL&source=WEB&client=WEB&filter=(SECURITY_CODE%3D%22600115%22)(END_DATE%3D%272023-03-31%27)
// callback: jQuery112308928698030561204_1687518793792
//sortColumns: RANK
//sortTypes: 1
//pageSize: 10
//pageNumber: 1
//reportName: RPT_DMSK_HOLDERS
//columns: ALL
//source: WEB
//client: WEB
//filter: (SECURITY_CODE="600115")(END_DATE='2023-03-31')

type rawStockHolder struct {
	Version string `json:"version"`
	Result  struct {
		Pages int `json:"pages"`
		Data  []struct {
			SECUCODE                string  `json:"SECUCODE"`
			SECURITY_CODE           string  `json:"SECURITY_CODE"`
			ORG_CODE                string  `json:"ORG_CODE"`
			END_DATE                string  `json:"END_DATE"`
			HOLDER_NAME             string  `json:"HOLDER_NAME"`
			HOLD_NUM                int64   `json:"HOLD_NUM"`
			FREE_HOLDNUM_RATIO      float64 `json:"FREE_HOLDNUM_RATIO"`
			HOLD_NUM_CHANGE         string  `json:"HOLD_NUM_CHANGE"`
			CHANGE_RATIO            float64 `json:"CHANGE_RATIO"`
			IS_HOLDORG              string  `json:"IS_HOLDORG"`
			HOLDER_RANK             int     `json:"HOLDER_RANK"`
			SECURITY_NAME_ABBR      string  `json:"SECURITY_NAME_ABBR"`
			HOLDER_CODE             string  `json:"HOLDER_CODE"`
			SECURITY_TYPE_CODE      string  `json:"SECURITY_TYPE_CODE"`
			HOLDER_STATE            string  `json:"HOLDER_STATE"`
			HOLDER_MARKET_CAP       float64 `json:"HOLDER_MARKET_CAP"`
			HOLD_RATIO              float64 `json:"HOLD_RATIO"`
			HOLD_CHANGE             string  `json:"HOLD_CHANGE"`
			HOLD_RATIO_CHANGE       float64 `json:"HOLD_RATIO_CHANGE"`
			HOLDER_TYPE             string  `json:"HOLDER_TYPE"`
			SHARES_TYPE             string  `json:"SHARES_TYPE"`
			UPDATE_DATE             string  `json:"UPDATE_DATE"`
			REPORTDATENAME          string  `json:"REPORT_DATE_NAME"`
			REPORT_DATE_NAME        string  `json:"HOLDER_NEW"`
			FREE_RATIO_QOQ          string  `json:"FREE_RATIO_QOQ"`
			HOLDER_STATEE           string  `json:"HOLDER_STATEE"`
			IS_REPORT               string  `json:"IS_REPORT"`
			HOLDER_CODE_OLD         string  `json:"HOLDER_CODE_OLD"`
			HOLDER_NEWTYPE          string  `json:"HOLDER_NEWTYPE"`
			HOLDNUM_CHANGE_NAME     string  `json:"HOLDNUM_CHANGE_NAME"`
			IS_MAX_REPORTDATE       string  `json:"IS_MAX_REPORTDATE"`
			COOPERATION_HOLDER_MARK string  `json:"COOPERATION_HOLDER_MARK"`
			MXID                    string  `json:"MXID"`
			LISTING_STATE           string  `json:"LISTING_STATE"`
			XZCHANGE                int     `json:"XZCHANGE"`
			NEW_CHANGE_RATIO        string  `json:"NEW_CHANGE_RATIO"`
		} `json:"data"`
		Count int `json:"count"`
	} `json:"result"`
	Success bool   `json:"success"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// 前十大流通股东 https://data.eastmoney.com/gdfx/stock/600115.html
// 数据接口 https://datacenter-web.eastmoney.com/api/data/v1/get
//callback: jQuery112308928698030561204_1687518793796
//sortColumns: HOLDER_RANK
//sortTypes: 1
//pageSize: 10
//pageNumber: 1
//reportName: RPT_F10_EH_FREEHOLDERS
//columns: ALL
//source: WEB
//client: WEB
//filter: (SECURITY_CODE="600115")(END_DATE='2023-03-31')

const (
	urlTop10ShareHolder = "https://datacenter-web.eastmoney.com/api/data/v1/get"
)

func ShareHolder(securityCode, date string, diffQuarters ...int) (list []CirculatingShareholder) {
	_, _, code := proto.DetectMarket(securityCode)
	quarterEndDate := trading.FixTradeDate(date)
	_, _, qEnd := api.GetQuarterByDate(date, diffQuarters...)
	quarterEndDate = trading.FixTradeDate(qEnd)
	params := urlpkg.Values{
		"sortColumns": {"HOLDER_RANK"},
		"sortTypes":   {"1"},
		"pageSize":    {"10"},
		"pageNumber":  {"1"},
		"reportName":  {"RPT_F10_EH_FREEHOLDERS"},
		"columns":     {"ALL"},
		"source":      {"WEB"},
		"client":      {"WEB"},
		"filter":      {fmt.Sprintf("(SECURITY_CODE=\"%s\")(END_DATE='%s')", code, quarterEndDate)},
	}

	url := urlTop10ShareHolder + "?" + params.Encode()
	data, err := http.HttpGet(url)
	//fmt.Println(api.Bytes2String(data))
	if err != nil {
		return
	}
	var raw rawStockHolder
	err = json.Unmarshal(data, &raw)
	if err != nil || raw.Result.Count == 0 || len(raw.Result.Data) == 0 {
		return
	}
	for _, v := range raw.Result.Data {
		shareholder := CirculatingShareholder{
			SecurityCode:     v.SECUCODE,
			SecurityName:     v.SECURITY_NAME_ABBR,
			EndDate:          trading.FixTradeDate(v.END_DATE),
			UpdateDate:       trading.FixTradeDate(v.UPDATE_DATE),
			HolderType:       v.HOLDER_NEWTYPE,
			HolderName:       v.HOLDER_NAME,
			IsHoldOrg:        v.IS_HOLDORG,
			HolderRank:       v.HOLDER_RANK,
			HoldNum:          int(v.HOLD_NUM),
			FreeHoldNumRatio: v.FREE_HOLDNUM_RATIO,
			HoldNumChange:    v.XZCHANGE,
			HoldChangeName:   v.HOLDNUM_CHANGE_NAME,
			HoldChangeRatio:  v.CHANGE_RATIO,
			HoldRatio:        v.HOLD_RATIO,
			HoldRatioChange:  v.HOLD_RATIO_CHANGE,
		}
		// 修订证券代码
		_, mfalg, mcode := proto.DetectMarket(shareholder.SecurityCode)
		shareholder.SecurityCode = mfalg + mcode
		//HoldChangeState  int     `dataframe:"change_state"`        // 期末持股-变化状态
		switch v.HOLDNUM_CHANGE_NAME {
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
			warning := fmt.Sprintf("%s: %s, 变化状态未知: %s", v.SECURITY_NAME_ABBR, v.SECUCODE, v.HOLDNUM_CHANGE_NAME)
			logger.Warnf(warning)
		}
		list = append(list, shareholder)
	}
	api.SliceSort(list, func(a, b CirculatingShareholder) bool {
		return a.HolderRank < b.HolderRank
	})
	return
}

// cacheShareHolder 获取流动股东数据
func cacheShareHolder(securityCode, date string, diffQuarters ...int) (list []CirculatingShareholder) {
	diff := 1
	if len(diffQuarters) > 0 {
		diff = diffQuarters[0]
	}
	_, _, last := api.GetQuarterByDate(date, diff)
	filename := cache.Top10HoldersFilename(securityCode, last)
	if api.FileExist(filename) {
		err := api.CsvToSlices(filename, &list)
		if err == nil && len(list) > 0 {
			return
		}
	}
	tmpList := ShareHolder(securityCode, last)
	if len(tmpList) > 0 {
		list = tmpList
	}
	if len(list) > 0 {
		_ = api.SlicesToCsv(filename, list)
	}
	return
}

// GetCacheShareHolder 获取流动股东数据
func GetCacheShareHolder(securityCode, date string, diffQuarters ...int) (list []CirculatingShareholder) {
	diff := 1
	if len(diffQuarters) > 0 {
		diff = diffQuarters[0]
	}
	for ; diff < 4; diff++ {
		tmpList := cacheShareHolder(securityCode, date, diff)
		if len(tmpList) == 0 {
			continue
		}
		list = tmpList
		break
	}
	return
}
