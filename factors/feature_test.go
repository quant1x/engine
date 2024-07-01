package factors

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
	"unsafe"
)

//func TestNewDataBuilder(t *testing.T) {
//	date := "2023-07-04"
//	v := NewDataBuilder("test", date, nil)
//	fmt.Println(v)
//}

// typelinks2 for 1.7 ~
//
//go:linkname typelinks2 reflect.typelinks
func typelinks2() (sections []unsafe.Pointer, offset [][]int32)

//go:linkname resolveTypeOff reflect.resolveTypeOff
func resolveTypeOff(rtype unsafe.Pointer, off int32) unsafe.Pointer

var typeMaps map[string]reflect.Type
var packages map[string]map[string]reflect.Type

type emptyInterface struct {
	typ  unsafe.Pointer
	word unsafe.Pointer
}

func loadGoTypes() {
	fia := reflect.TypeOf((*Feature)(nil)).Elem()
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
				typeName := loadedType.String()
				//fmt.Println(typeName, "==>", loadedType.PkgPath())
				typeMaps[typeName] = loadedType
				pkgTypes[loadedType.Name()] = loadedType
				//if loadedType.Implements(Feature) {
				//	fmt.Println(typeName, "==>", loadedType.PkgPath())
				//}
				pkgPath := loadedType.PkgPath()
				if strings.HasPrefix(pkgPath, "gitee.com/quant1x/engine") {
					//fmt.Println(typeName, "==>", loadedType.PkgPath())
					if reflect.PtrTo(loadedType).Implements(fia) {
						fmt.Println("found", pkgPath, "==>", loadedType.Name())
					}
				}
			}
		}
	}
}

func TestInterface(t *testing.T) {
	typeMaps = make(map[string]reflect.Type)
	packages = make(map[string]map[string]reflect.Type)
	loadGoTypes()
}
