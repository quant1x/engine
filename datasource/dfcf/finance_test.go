package dfcf

import (
	"encoding/json"
	"fmt"
	"gitee.com/quant1x/gox/api"
	"testing"
)

func TestFinanceReports(t *testing.T) {
	date := "2022-09-30"
	date = "2023-03-01"
	reports, _, _, err := FinanceReports(date)
	if err != nil {
		return
	}
	data, err := json.Marshal(reports)
	fmt.Println(api.Bytes2String(data))
}
