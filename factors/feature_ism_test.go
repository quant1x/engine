package factors

import (
	"encoding/json"
	"fmt"
	"gitee.com/quant1x/engine/cache"
	"gitee.com/quant1x/gox/api"
	"testing"
)

func TestISM(t *testing.T) {
	code := "000737"
	date := "2024-03-21"
	cacheDate, featureDate := cache.CorrectDate(date)
	s8 := NewInvestmentSentimentMaster(cacheDate, code)
	s8.Update(code, cacheDate, featureDate, true)
	data, _ := json.Marshal(s8)
	text := api.Bytes2String(data)
	fmt.Println(text)
}
