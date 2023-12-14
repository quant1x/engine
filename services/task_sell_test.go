package services

import (
	"fmt"
	"testing"
)

func TestTaskSell_getEarlierDate(t *testing.T) {
	v := getEarlierDate(1)
	fmt.Println(v)
}
