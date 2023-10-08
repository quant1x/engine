package cache

import (
	"fmt"
	"testing"
)

func TestFilename(t *testing.T) {
	date := "2023-09-28"
	code := "sh600105"
	filename := QuarterlyReportFilename(code, date)
	fmt.Println(filename)
}
