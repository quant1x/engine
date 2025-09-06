package dfcf

import (
	"encoding/json"
	"fmt"
	"testing"

	"gitee.com/quant1x/gox/api"
)

func TestQuarterlyReports(t *testing.T) {
	date := "20230928"

	list, pages, err := QuarterlyReports(date, 1)
	fmt.Println(list)
	fmt.Println(pages)
	fmt.Println(err)
}

func TestGetCacheQuarterlyReportsBySecurityCode(t *testing.T) {
	date := "20231027"
	code := "sh600178"
	v := GetCacheQuarterlyReportsBySecurityCode(code, date)
	fmt.Println(v)
	data, _ := json.Marshal(v)
	text := api.Bytes2String(data)
	fmt.Println(text)
}
