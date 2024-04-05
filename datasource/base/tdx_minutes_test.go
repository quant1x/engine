package base

import (
	"fmt"
	"testing"
)

func TestGetMinutes(t *testing.T) {
	code := "sh000001"
	code = "sh510050"
	code = "sh600105"
	date := "2023-09-28"
	v := GetMinutes(code, date)
	fmt.Println(v)
}

func TestLoadMinutes(t *testing.T) {
	code := "sh000001"
	code = "sh510050"
	code = "sh600105"
	date := "2024-01-02"
	v := LoadMinutes(code, date)
	fmt.Println(v)
}
