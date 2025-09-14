package shse

import (
	"fmt"
	"math"
	"math/rand"
	urlpkg "net/url"
	"time"

	"gitee.com/quant1x/gox/http"
)

const (
	// http://www.sse.com.cn/market/sseindex/indexlist/
	urlSSEIndex = "https://query.sse.com.cn/commonSoaQuery.do"
	//urlSSEIndex = "http://query.sse.com.cn/commonQuery.do"
)

// IndexList 上海指数列表
func IndexList() {
	//isPagination=false&sqlId=DB_SZZSLB_ZSLB&_=1682986576599
	now := time.Now()
	timestamp := now.UnixMilli()
	cbNum := int(math.Floor(rand.Float64() * (100000000 + 1)))
	// 37039113
	params := urlpkg.Values{
		"jsonCallBack": {fmt.Sprintf("jsonpCallback%d", cbNum)},
		"isPagination": {"false"},
		"sqlId":        {"DB_SZZSLB_ZSLB"},
		"_":            {fmt.Sprintf("%d", timestamp)},
	}
	header := map[string]any{
		"Referer": "http://www.sse.com.cn/",
		//"Cookie":  "ba17301551dcbaf9_gdp_user_key=; gdp_user_id=gioenc-1aadbe52,d720,54c2,9271,29b7b6dda362; ba17301551dcbaf9_gdp_session_id_2960a971-ddff-48be-827f-6eb99e891735=true; ba17301551dcbaf9_gdp_session_id_4602911b-d360-4f09-a438-3d40bca228d7=true; ba17301551dcbaf9_gdp_session_id_86910507-dd22-4b31-9749-e3a3b18eae25=true; ba17301551dcbaf9_gdp_session_id_1bb67914-f729-4e41-b75c-d09a6b0d7873=true; JSESSIONID=7052255EE4B2357019E75B7B09D6D571; ba17301551dcbaf9_gdp_session_id=2a1f157c-1605-45e8-a68b-e4196d04b2af; ba17301551dcbaf9_gdp_session_id_2a1f157c-1605-45e8-a68b-e4196d04b2af=true; ba17301551dcbaf9_gdp_sequence_ids={\"globalKey\":42,\"VISIT\":6,\"PAGE\":14,\"VIEW_CHANGE\":2,\"CUSTOM\":3,\"VIEW_CLICK\":21}",
	}
	url := urlSSEIndex + "?"
	//url += fmt.Sprintf("jsonCallBack=jsonpCallback%d&", cbNum)
	//url += "isPagination=false&"
	//url += "sqlId=DB_SZZSLB_ZSLB&"
	url += params.Encode()
	data, tm, err := http.Request(url, "get", "", header)
	fmt.Println(string(data), tm, err)

}
