package factors

import (
	"encoding/json"
	"fmt"
	"testing"

	"gitee.com/quant1x/engine/cache"
	"gitee.com/quant1x/exchange"
	"gitee.com/quant1x/gox/api"
)

func TestFeatureSecuritiesMarginTrading(t *testing.T) {
	code := "600580"
	code = exchange.CorrectSecurityCode(code)
	date := "2025-03-07"
	cacheDate, featureDate := cache.CorrectDate(date)
	feature := NewSecuritiesMarginTrading(date, code)
	feature.Init(nil, featureDate)
	feature.Update(code, cacheDate, featureDate, false)
	data, _ := json.Marshal(feature)
	text := api.Bytes2String(data)
	fmt.Println(text)
}
