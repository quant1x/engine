package trader

import (
	"fmt"
	"testing"
)

func TestQueryAccount(t *testing.T) {
	info, err := QueryAccount()
	fmt.Println(info, err)
}

func TestQueryHolding(t *testing.T) {
	info, err := QueryHolding()
	fmt.Println(info, err)
}
