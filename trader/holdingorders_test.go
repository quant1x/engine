package trader

import (
	"encoding/json"
	"fmt"
	"gitee.com/quant1x/gox/api"
	"testing"
)

func Test_lazyLoadHoldingOrder(t *testing.T) {
	lazyLoadHoldingOrder()
	data, err := json.Marshal(holdingOrders)
	fmt.Println(data, err)
	if err != nil {
		return
	}
	text := api.Bytes2String(data)
	fmt.Println(text)
}
