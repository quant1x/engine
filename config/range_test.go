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

func TestNumberRange(t *testing.T) {
	text := ""
	fmt.Println("==>", text)
	var r NumberRange
	err := r.Parse(text)
	fmt.Println(r, err)
	fmt.Println(-1000000 < r.Min())
	text = "1"
	fmt.Println("==>", text)
	err = r.Parse(text)
	fmt.Println(r, err)
	text = "1~"
	fmt.Println("==>", text)
	err = r.Parse(text)
	fmt.Println(r, err)
	text = "~1"
	fmt.Println("==>", text)
	err = r.Parse(text)
	fmt.Println(r, err)
	text = "1~2"
	fmt.Println("==>", text)
	err = r.Parse(text)
	fmt.Println(r, err)
	text = "1~2~3"
	fmt.Println("==>", text)
	err = r.Parse(text)
	fmt.Println(r, err)
}
