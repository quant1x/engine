package market

import (
	"fmt"
	"testing"
)

func TestGetCodeList(t *testing.T) {
	codes := GetCodeList()
	fmt.Println(len(codes))
	codes = GetStockCodeList()
	fmt.Println(len(codes))
}
