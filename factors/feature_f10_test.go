package factors

import (
	"context"
	"encoding/json"
	"fmt"
	"gitee.com/quant1x/engine/cache"
	"gitee.com/quant1x/gox/api"
	"testing"
)

func TestF10(t *testing.T) {
	date := "2023-10-17"
	cacheDate, featureDate := cache.CorrectDate(date)
	//cacheDate := "2023-09-28"
	//featureDate := date
	code := "sh600105"
	code = "sh000001"
	code = "sh600859"
	f10 := NewF10(cacheDate, code)
	//barIndex := 1
	ctx := context.Background()
	f10.Init(ctx, featureDate, code)
	f10.Repair(code, cacheDate, featureDate, true)
	data, _ := json.Marshal(f10)
	text := api.Bytes2String(data)
	fmt.Println(text)
}
