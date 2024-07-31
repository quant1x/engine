package dfcf

import (
	"fmt"
	"testing"
)

func Test_rawBillBoard(t *testing.T) {
	date := "20240730"
	v, n, err := rawBillBoardList(date, 1, lhbBuy)
	fmt.Println(v, n, err)
}
