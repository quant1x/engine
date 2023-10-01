package market

import (
	"fmt"
	"testing"
)

func TestIsSubNewStock(t *testing.T) {
	code := "603052"
	v := IsSubNewStock(code)
	fmt.Println(v)
}
