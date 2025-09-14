package config

import (
	"fmt"
	"testing"

	"gitee.com/quant1x/gox/api"
	"gitee.com/quant1x/pkg/yaml"
)

func TestTimeRange(t *testing.T) {
	text := " 09:30:00 ~ 14:56:30 "
	text = " 14:56:30 - 09:30:00 "
	var tr TimeRange
	tr.Parse(text)
	fmt.Println(tr)
	text = "09:15:00~09:26:59,09:15:00~09:19:59,09:25:00~11:29:59,13:00:00~14:59:59,09:00:00~09:14:59"
	var ts TradingSession
	ts.Parse(text)
	fmt.Println(ts)
	fmt.Println(ts.IsTrading())
}

type tt struct {
	Session    TimeRange `yaml:"session"`
	ChangeRate float64   `yaml:"change_rate" default:"0.01"`
}

func TestTimeRangeWithParse(t *testing.T) {
	text := `time: "09:50:00~09:50:59,10:50:00~10:50:59"`
	text = `
session: "09:50:00~09:50:59"
name: "buyiwang"
`
	bytes := api.String2Bytes(text)
	v := tt{}
	err := yaml.Unmarshal(bytes, &v)
	fmt.Println(err, v)
}
