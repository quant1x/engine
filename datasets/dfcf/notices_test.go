package dfcf

import (
	"encoding/json"
	"fmt"
	"gitee.com/quant1x/gox/api"
	"testing"
)

func TestStockNoticeReport(t *testing.T) {
	notices, _, err := StockNotices("600105", "20230828", "20230928", 1)
	if err != nil {
		return
	}
	data, err := json.Marshal(notices)
	fmt.Println(api.Bytes2String(data))
}
