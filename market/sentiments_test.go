package market

import (
	"fmt"
	"testing"
)

func TestSentiment(t *testing.T) {
	s, c := IndexSentiment("sh880521")
	fmt.Println(s, c)
}
