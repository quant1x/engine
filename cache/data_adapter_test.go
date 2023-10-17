package cache

import (
	"fmt"
	"slices"
	"testing"
)

func TestPlugins(t *testing.T) {
	k := Kind(0)
	v := k & 1
	fmt.Println(v)
	k1 := PluginMaskBaseData | 2
	v = k1 & PluginMaskBaseData
	fmt.Println(v == PluginMaskBaseData)

	list := []uint{3, 1, 2, 5, 4}
	slices.Sort(list)
	fmt.Println(list)
}
