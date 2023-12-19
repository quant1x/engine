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
