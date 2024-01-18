package trader

import (
	"fmt"
	"testing"
)

func TestGetOrderDates(t *testing.T) {
	list := GetLocalOrderDates()
	fmt.Println(list)
}
