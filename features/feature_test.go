package features

import (
	"fmt"
	"testing"
)

func TestNewDataBuilder(t *testing.T) {
	date := "2023-07-04"
	v := NewDataBuilder("test", date, nil)
	fmt.Println(v)
}
