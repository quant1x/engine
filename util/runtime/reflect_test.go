package runtime

import (
	"fmt"
	"gitee.com/quant1x/engine/features"
	"reflect"
	"testing"
)

func TestReflect(t *testing.T) {
	//lazyInit()
	fia := reflect.TypeOf((*features.Feature)(nil)).Elem()
	v := FindImplements(fia)
	fmt.Println(v)
}
