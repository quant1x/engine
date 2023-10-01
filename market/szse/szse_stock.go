package szse

import (
	"fmt"
	"gitee.com/quant1x/engine/internal/functions"
	"gitee.com/quant1x/gox/http"
	urlpkg "net/url"
)

const (
	kUrlMarketSzseCodeList = "http://www.szse.cn/api/report/ShowReport/data"
)

func GetStockList() {
	timestamp := functions.Timestamp()
	params := urlpkg.Values{
		"SHOWTYPE":     {"JSON"},
		"CATALOGID":    {"1815_stock_snapshot"},
		"TABKEY":       {"tab1"},
		"txtBeginDate": {"2023-06-16"},
		"txtEndDate":   {"2023-06-16"},
		"archiveDate":  {"2021-06-01"},
		"random":       {fmt.Sprintf("0.%d", timestamp)},
		"PAGENO":       {"2"},
		"PAGESIZE":     {"100"},
		"tab1PAGESIZE": {"100"},
	}
	header := map[string]any{
		"Referer": "http://www.szse.cn/market/trend/index.html",
		//"Cookie":  "ba17301551dcbaf9_gdp_user_key=; gdp_user_id=gioenc-1aadbe52,d720,54c2,9271,29b7b6dda362; ba17301551dcbaf9_gdp_session_id_2960a971-ddff-48be-827f-6eb99e891735=true; ba17301551dcbaf9_gdp_session_id_4602911b-d360-4f09-a438-3d40bca228d7=true; ba17301551dcbaf9_gdp_session_id_86910507-dd22-4b31-9749-e3a3b18eae25=true; ba17301551dcbaf9_gdp_session_id_1bb67914-f729-4e41-b75c-d09a6b0d7873=true; JSESSIONID=7052255EE4B2357019E75B7B09D6D571; ba17301551dcbaf9_gdp_session_id=2a1f157c-1605-45e8-a68b-e4196d04b2af; ba17301551dcbaf9_gdp_session_id_2a1f157c-1605-45e8-a68b-e4196d04b2af=true; ba17301551dcbaf9_gdp_sequence_ids={\"globalKey\":42,\"VISIT\":6,\"PAGE\":14,\"VIEW_CHANGE\":2,\"CUSTOM\":3,\"VIEW_CLICK\":21}",
	}
	url := kUrlMarketSzseCodeList + "?" + params.Encode()
	data, _, _ := http.Request(url, "get", header)
	fmt.Println(string(data))
	//if err != nil {
	//	return nil, err
	//}
	////fmt.Println(string(data), tm, err)
	////fmt.Println(string(data))
	//var raw rawShangHaiSecurities
	//err = json.Unmarshal(data, &raw)
	//if err != nil {
	//	return
	//}
	//for _, vs := range raw.List {
	//	arr := []string{}
	//	for _, v := range vs {
	//		arr = append(arr, stat.AnyToString(v))
	//	}
	//	var info sseSecurityEntity
	//	_ = api.Convert(arr, &info)
	//	list = append(list, info)
	//}
}
