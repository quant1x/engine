package factors

import (
	"fmt"
	"testing"
)

func TestGetL5History(t *testing.T) {
	code := "sz000713"
	v := GetL5History(code)
	fmt.Println(v)
}

func TestMisc(t *testing.T) {
	code := "sh000001"
	v := GetL5Misc(code)
	fmt.Println(v)
}
