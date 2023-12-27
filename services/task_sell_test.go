package services

import (
	"fmt"
	"gitee.com/quant1x/engine/config"
	"gitee.com/quant1x/engine/models"
	"reflect"
	"testing"
)

func TestTaskSell_getEarlierDate(t *testing.T) {
	v := getEarlierDate(1)
	fmt.Println(v)
}

func TestTaskSell_LastSession(t *testing.T) {
	sellStrategyCode := models.ModelOneSizeFitsAllSells
	// 1. 获取117号策略(卖出)
	sellRule := config.GetStrategyParameterByCode(sellStrategyCode)
	if sellRule == nil {
		return
	}
	timestamp := "14:55:00"
	v := sellRule.Session.IsTodayLastSession(timestamp)
	fmt.Println(v)
}

type MyStruct struct {
	Field1 string
	Field2 int
}

func TestStructPtr(t *testing.T) {
	// 创建一个结构体切片
	structSlice := []MyStruct{
		{"first", 1},
		{"second", 2},
		{"third", 3},
	}

	// 获取切片的反射值对象
	sliceValue := reflect.ValueOf(structSlice)
	// 遍历切片
	for i := 0; i < sliceValue.Len(); i++ {
		// 获取切片元素
		elem := sliceValue.Index(i)
		// 获取元素的地址
		elemAddr := elem.Addr().Interface()
		// 打印元素的地址
		fmt.Printf("Element %d address: %v\n", i, elemAddr)
	}
}

func Test_checkoutCanSellStockList(t *testing.T) {
	v := checkoutCanSellStockList(117)
	fmt.Println(v)
}
