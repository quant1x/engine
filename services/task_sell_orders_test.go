package services

import (
	"fmt"
	"testing"
)

func TestTaskSell_getEarlierDate(t *testing.T) {
	v := getEarlierDate(1)
	fmt.Println(v)
}

func Test_checkoutCanSellStockList(t *testing.T) {
	v := CheckoutCanSellStockList(117)
	fmt.Println(v)
}
