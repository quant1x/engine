package labs

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/quant1x/engine/factors"
)

func TestReflect(t *testing.T) {
	//lazyInit()
	fia := reflect.TypeOf((*factors.Feature)(nil)).Elem()
	v := FindImplements(fia)
	fmt.Println(v)
}
