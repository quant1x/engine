package features

import (
	"encoding/json"
	"fmt"
	"testing"

	"gitee.com/quant1x/gox/api"
)

func TestNewBreaksThrough(t *testing.T) {
	code := "sh600600"
	code = "sz002043"
	code = "sz000638"
	code = "600105"
	code = "sh603193"
	date := "20231012"
	v := NewKLineBox(code, date)
	data, _ := json.Marshal(v)
	text := api.Bytes2String(data)
	fmt.Println(text)
}
