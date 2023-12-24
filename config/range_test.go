package config

import (
	"fmt"
	"math"
	"testing"
)

func TestParseRange(t *testing.T) {
	v := math.SmallestNonzeroFloat64
	fmt.Println(v)
	v = math.MaxFloat64
	fmt.Println(v)
	fmt.Println(-10000000000.00 < v)

	rn := NumberRange{}
	ok := rn.Validate(100)
	fmt.Println(ok)
}
