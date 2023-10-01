package shse

import (
	"fmt"
	"testing"
)

func TestGetList(t *testing.T) {
	v, _ := GetSecurityList()
	fmt.Println(v)
}
