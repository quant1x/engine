package runtime

import (
	"fmt"
	"gitee.com/quant1x/engine/factors"
	"reflect"
	"testing"
)

func TestReflect(t *testing.T) {
	//lazyInit()
	fia := reflect.TypeOf((*factors.Feature)(nil)).Elem()
	v := FindImplements(fia)
	fmt.Println(v)
}
