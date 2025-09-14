package shse

import (
	"encoding/json"
	"fmt"
	urlpkg "net/url"

	"github.com/quant1x/engine/utils"
	"github.com/quant1x/num"
	"github.com/quant1x/x/api"
	"github.com/quant1x/x/http"
)

const (
	// 数据来源: http://www.sse.com.cn/market/price/report/
	kUrlMarketSseCodeList = "http://yunhq.sse.com.cn:32041/v1/sh1/list/exchange/equity"
)

type rawShangHaiSecurities struct {
	Date  int     `json:"date"`
	Time  int     `json:"time"`
	Total int     `json:"total"`
	Begin int     `json:"begin"`
	End   int     `json:"end"`
	List  [][]any `json:"list"`
}

type sseSecurityEntity struct {
	Code         string  `array:"0"`
	Name         string  `array:"1"`
	Open         float64 `array:"2"`
	High         float64 `array:"3"`
	Low          float64 `array:"4"`
	Last         float64 `array:"5"`
	PrevClose    float64 `array:"6"`
	ChangeRate   float64 `array:"7"`
	Volume       int64   `array:"8"`
	Amount       float64 `array:"9"`
	TradePhase   string  `array:"10"`
	Change       float64 `array:"11"`
	AmpRate      float64 `array:"12"`
	CpxxSubType  string  `array:"13"`
	CpxxProdusta string  `array:"14"`
}

// GetSecurityList 获取证券代码列表
func GetSecurityList() (list []sseSecurityEntity, err error) {
	timestamp := utils.Timestamp()
	params := urlpkg.Values{
		"_":            {fmt.Sprintf("%d", timestamp)},
		"isPagination": {"false"},
		"begin":        {"0"},
		"end":          {"5000"},
		"select":       {"code,name,open,high,low,last,prev_close,chg_rate,volume,amount,tradephase,change,amp_rate,cpxxsubtype,cpxxprodusta"},
		//"select": {"code,name"},
	}
	header := map[string]any{
		"Referer": "http://www.sse.com.cn/",
		//"Cookie":  "ba17301551dcbaf9_gdp_user_key=; gdp_user_id=gioenc-1aadbe52,d720,54c2,9271,29b7b6dda362; ba17301551dcbaf9_gdp_session_id_2960a971-ddff-48be-827f-6eb99e891735=true; ba17301551dcbaf9_gdp_session_id_4602911b-d360-4f09-a438-3d40bca228d7=true; ba17301551dcbaf9_gdp_session_id_86910507-dd22-4b31-9749-e3a3b18eae25=true; ba17301551dcbaf9_gdp_session_id_1bb67914-f729-4e41-b75c-d09a6b0d7873=true; JSESSIONID=7052255EE4B2357019E75B7B09D6D571; ba17301551dcbaf9_gdp_session_id=2a1f157c-1605-45e8-a68b-e4196d04b2af; ba17301551dcbaf9_gdp_session_id_2a1f157c-1605-45e8-a68b-e4196d04b2af=true; ba17301551dcbaf9_gdp_sequence_ids={\"globalKey\":42,\"VISIT\":6,\"PAGE\":14,\"VIEW_CHANGE\":2,\"CUSTOM\":3,\"VIEW_CLICK\":21}",
	}
	url := kUrlMarketSseCodeList + "?" + params.Encode()
	//data, _, err := http.Request(url, "get", header)
	data, err := http.Get(url, header)
	if err != nil {
		return nil, err
	}
	//fmt.Println(string(data), tm, err)
	//fmt.Println(string(data))
	var raw rawShangHaiSecurities
	err = json.Unmarshal(data, &raw)
	if err != nil {
		return
	}
	for _, vs := range raw.List {
		arr := []string{}
		for _, v := range vs {
			arr = append(arr, num.AnyToString(v))
		}
		var info sseSecurityEntity
		_ = api.Convert(arr, &info)
		list = append(list, info)
	}
	return
}
