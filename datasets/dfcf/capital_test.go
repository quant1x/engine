package dfcf

import (
	"fmt"
	"testing"
	"time"
)

func TestCapitalChange(t *testing.T) {
	code := "600115"
	v := CapitalChange(code)
	fmt.Println(v)
	fmt.Println(time.Local)
}
