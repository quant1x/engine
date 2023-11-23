package labs

import (
	"fmt"
	"gitee.com/quant1x/engine/factors"
	"reflect"
	"strings"
	"sync"
	"unsafe"
)

//go:linkname typelinks2 reflect.typelinks
func typelinks2() (sections []unsafe.Pointer, offset [][]int32)

//go:linkname resolveTypeOff reflect.resolveTypeOff
func resolveTypeOff(rtype unsafe.Pointer, off int32) unsafe.Pointer

var (
	typeOnce sync.Once
	types    = map[string]reflect.Type{}
	packages = map[string]map[string]reflect.Type{}
)

type emptyInterface struct {
	typ  unsafe.Pointer
	word unsafe.Pointer
}

func loadGoTypes() {
	fia := reflect.TypeOf((*factors.Feature)(nil)).Elem()
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
				//types[typeString] = loadedType
				pkgPath := loadedType.PkgPath()
				typeName := loadedType.Name()
				//typeString := loadedType.String()
				//fmt.Println(pkgPath, "=>", typeString, "=>", typeName)
				structName := fmt.Sprintf("%s.%s", pkgPath, typeName)
				types[structName] = loadedType
				pkgTypes[loadedType.Name()] = loadedType
				if strings.HasPrefix(pkgPath, "gitee.com/quant1x/engine") {
					//fmt.Println(structName, "==>", loadedType.PkgPath())
					if reflect.PtrTo(loadedType).Implements(fia) {
						fmt.Println("found", pkgPath, "==>", loadedType.String(), "==>", structName)
					}
				}
			}
		}
	}
}

func lazyInit() {
	loadGoTypes()
}

func FindImplements(u reflect.Type) (list []reflect.Type) {
	typeOnce.Do(lazyInit)
	for name, t := range types {
		//fmt.Println(name)
		if reflect.PtrTo(t).Implements(u) {
			list = append(list, t)
		}
		_ = name
	}
	return
}
