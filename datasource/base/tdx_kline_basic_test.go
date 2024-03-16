package base

import (
	"fmt"
	"testing"
)

func TestUpdateAllBasicKLine(t *testing.T) {
	code := "sh000001"
	data := UpdateAllBasicKLine(code)
	fmt.Println(data[len(data)-1])
	_ = data
}
