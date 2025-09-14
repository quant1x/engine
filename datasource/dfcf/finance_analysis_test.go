package dfcf

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/quant1x/x/api"
)

func TestGetQuarterlyReports(t *testing.T) {
	v, n, err := GetQuarterlyReports()
	fmt.Println(v)
	data, _ := json.Marshal(v)
	text := api.Bytes2String(data)
	fmt.Println(text)
	fmt.Println(n)
	fmt.Println(err)
}
