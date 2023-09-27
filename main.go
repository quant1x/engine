package main

import (
	"fmt"
	"gitee.com/quant1x/engine/factors"
	"gitee.com/quant1x/engine/features/base"
)

func main() {
	var feature factors.Feature
	feature = new(base.KLine)
	fmt.Println(feature.Kind())
}
