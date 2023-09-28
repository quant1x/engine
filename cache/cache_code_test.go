package cache

import (
	"fmt"
	"testing"
)

func TestCacheIdPath(t *testing.T) {
	code := "600072"
	code = "000001"
	v := CacheIdPath(code)
	fmt.Println(v)
}
