package models

import (
	"fmt"
	"testing"
)

func TestBatchSnapShot(t *testing.T) {
	data := BatchSnapShot([]string{"sz003042"})
	fmt.Printf("%+v\n", data)
}
