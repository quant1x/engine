package dfcf

import (
	"fmt"
	"testing"
)

func TestGetQuarterlyReports(t *testing.T) {
	v, n, err := GetQuarterlyReports()
	fmt.Println(v)
	fmt.Println(n)
	fmt.Println(err)
}
