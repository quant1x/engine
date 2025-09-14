package factors

import (
	"encoding/json"
	"fmt"
	"testing"

	"gitee.com/quant1x/engine/cache"
	"gitee.com/quant1x/exchange"
	"gitee.com/quant1x/gox/api"
)

func TestHistory(t *testing.T) {
	code := "300956"
	code = "300093"
	code = "301389"
	code = "000751"
	code = "301129"
	code = "002962"
	code = "600600"
	date := "2025-05-30"
	cacheDate, featureDate := cache.CorrectDate(date)
	code = exchange.CorrectSecurityCode(code)
	history := NewHistory(cacheDate, code)
	history.Update(code, cacheDate, featureDate, true)
	data, _ := json.Marshal(history)
	text := api.Bytes2String(data)
	fmt.Println(text)
}
