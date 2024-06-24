package factors

import (
	"encoding/json"
	"fmt"
	"gitee.com/quant1x/engine/cache"
	"gitee.com/quant1x/exchange"
	"gitee.com/quant1x/gox/api"
	"testing"
)

func TestBox_basic(t *testing.T) {
	code := "002766"
	date := "2024-06-24"
	cacheDate, featureDate := cache.CorrectDate(date)
	code = exchange.CorrectSecurityCode(code)
	box := NewBox(cacheDate, code)
	box.Update(code, cacheDate, featureDate, true)
	data, _ := json.Marshal(box)
	text := api.Bytes2String(data)
	fmt.Println(text)
}
