package factors

import (
	"encoding/json"
	"fmt"
	"gitee.com/quant1x/gox/api"
	"testing"
)

func TestNewExchangeKLine(t *testing.T) {
	code := "sh880652"
	code = "sh603099"
	code = "sh600354"
	//code = "sh603029"
	code = "sz002679"
	code = "sh600313"
	code = "002553"
	code = "880482"
	code = "sz300377"
	code = "sz000828"
	date := "2023-11-13"
	ek := NewExchangeKLine(code, date)
	data, _ := json.Marshal(ek)
	text := api.Bytes2String(data)
	fmt.Println(text)
}
