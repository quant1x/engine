package factors

import (
	"encoding/json"
	"fmt"
	"gitee.com/quant1x/engine/cache"
	"gitee.com/quant1x/exchange"
	"gitee.com/quant1x/gox/api"
	"testing"
)

func TestHistory(t *testing.T) {
	code := "300956"
	code = "300093"
	code = "000603"
	date := "2024-05-22"
	cacheDate, featureDate := cache.CorrectDate(date)
	code = exchange.CorrectSecurityCode(code)
	history := NewHistory(cacheDate, code)
	history.Update(code, cacheDate, featureDate, true)
	data, _ := json.Marshal(history)
	text := api.Bytes2String(data)
	fmt.Println(text)
}
