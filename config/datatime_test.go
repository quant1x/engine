package config

import (
	"fmt"
	"testing"
)

func TestTimeRange(t *testing.T) {
	text := " 09:30:00 ~ 14:56:30 "
	text = " 14:56:30 - 09:30:00 "
	tr := ParseTimeRange(text)
	fmt.Println(tr)
	text = "09:15:00~09:19:59,09:25:00~11:29:59,13:00:00~14:59:59"
	ts := ParseTradingSession(text)
	fmt.Println(ts)
	fmt.Println(ts.IsTrading())
}
