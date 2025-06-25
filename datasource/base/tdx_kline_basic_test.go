package base

import (
	"fmt"
	"testing"
)

func TestUpdateAllBasicKLine(t *testing.T) {
	code := "sz300773"
	data := UpdateAllBasicKLine(code)
	fmt.Println(data[0])
	_ = data
}
