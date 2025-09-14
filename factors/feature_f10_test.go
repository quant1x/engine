package factors

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/quant1x/engine/cache"
	"github.com/quant1x/x/api"
)

func TestF10(t *testing.T) {
	date := "2024-04-10"
	q := getQuarterlyYearQuarter(date)
	fmt.Println(q)
	cacheDate, featureDate := cache.CorrectDate(date)
	//cacheDate := "2023-09-28"
	//featureDate := date
	code := "sh600105"
	code = "sh000001"
	code = "sh600859"
	code = "sz002685"
	code = "sh603158"
	code = "sh600178"
	code = "sh880941"
	code = "sh600016"
	f10 := NewF10(cacheDate, code)
	//barIndex := 1
	ctx := context.Background()
	f10.Init(ctx, featureDate)
	f10.Repair(code, cacheDate, featureDate, true)
	data, _ := json.Marshal(f10)
	text := api.Bytes2String(data)
	fmt.Println(text)
}
