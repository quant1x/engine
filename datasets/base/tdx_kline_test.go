package base

import (
	"fmt"
	"testing"
)

func TestGetAllBasicKLine(t *testing.T) {
	code := "sh000001"
	code = "sz002043"
	code = "sz000567"
	code = "sz300580"
	df := UpdateAllBasicKLine(code)
	fmt.Println(df)
	_ = df
}
