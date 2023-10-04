package features

import (
	"fmt"
	"reflect"
	"testing"
	"unsafe"
)

func TestNewDataBuilder(t *testing.T) {
	date := "2023-07-04"
	v := NewDataBuilder("test", date, nil)
	fmt.Println(v)
}

// typelinks2 for 1.7 ~
//
//go:linkname typelinks2 reflect.typelinks
func typelinks2() (sections []unsafe.Pointer, offset [][]int32)

//go:linkname resolveTypeOff reflect.resolveTypeOff
func resolveTypeOff(rtype unsafe.Pointer, off int32) unsafe.Pointer

var types map[string]reflect.Type
var packages map[string]map[string]reflect.Type

type emptyInterface struct {
	typ  unsafe.Pointer
	word unsafe.Pointer
}

func loadGoTypes() {
	var obj interface{} = reflect.TypeOf(0)
	sections, offset := typelinks2()
	for i, offs := range offset {
		rodata := sections[i]
		for _, off := range offs {
			(*emptyInterface)(unsafe.Pointer(&obj)).word = resolveTypeOff(unsafe.Pointer(rodata), off)
			typ := obj.(reflect.Type)
			if typ.Kind() == reflect.Ptr && typ.Elem().Kind() == reflect.Struct {
				loadedType := typ.Elem()
				pkgTypes := packages[loadedType.PkgPath()]
				if pkgTypes == nil {
					pkgTypes = map[string]reflect.Type{}
					packages[loadedType.PkgPath()] = pkgTypes
				}
				types[loadedType.String()] = loadedType
				pkgTypes[loadedType.Name()] = loadedType
			}
		}
	}
}

func TestInterface(t *testing.T) {
	types = make(map[string]reflect.Type)
	packages = make(map[string]map[string]reflect.Type)
	loadGoTypes()
}
