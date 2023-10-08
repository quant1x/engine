package tdxweb

import (
	"fmt"
	"testing"
)

func TestGetSafetyScore(t *testing.T) {
	code := "sh510050"
	code = "sh600105"
	code = "sh000001"
	v := GetSafetyScore(code)
	fmt.Println(v)
}
