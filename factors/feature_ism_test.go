package factors

import (
	"encoding/json"
	"fmt"
	"gitee.com/quant1x/engine/cache"
	"gitee.com/quant1x/gox/api"
	"testing"
)

func TestS8(t *testing.T) {
	code := "002632"
	date := "2024-04-10"
	cacheDate, featureDate := cache.CorrectDate(date)
	s8 := NewInvestmentSentimentMaster(cacheDate, code)
	s8.Update(code, cacheDate, featureDate, true)
	data, _ := json.Marshal(s8)
	text := api.Bytes2String(data)
	fmt.Println(text)
}
