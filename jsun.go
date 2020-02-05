package jsun

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"sync"
)

var typeCache sync.Map

func Marshal(v interface{}) ([]byte, error) {
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		rv = reflect.Indirect(rv)
	}
	typ, _ := typeCache.LoadOrStore(rv, buildType(rv.Type()))
	nv := reflect.New(typ.(reflect.Type))
	fmt.Printf("%v\n", typ.(reflect.Type))
	copyValue(nv.Elem(), rv)
	return json.Marshal(nv.Interface())
}

func copyValue(dst, src reflect.Value) {
	for i := 0; i < src.NumField(); i++ {
		name := src.Type().Field(i).Name
		if dst.FieldByName(name).IsValid() && dst.FieldByName(name).CanSet() {
			field := src.Field(i)
			if field.Kind() == reflect.Ptr {
				field = field.Elem()
			}
			if field.Kind() == reflect.Struct {
				copyValue(dst.FieldByName(name), field)
			} else {
				dst.FieldByName(name).Set(field)
			}
		}
	}
}

func buildType(typ reflect.Type) reflect.Type {
	var fs []reflect.StructField
	visitType(typ, 0, &fs, "", "", true)
	return reflect.StructOf(fs)
}

func repeat(n int) string {
	s := ""
	for i := 0; i < n+1; i++ {
		s = s + "--"
	}
	return s
}

var marshalerType = reflect.TypeOf((*json.Marshaler)(nil)).Elem()

func visitType(typ reflect.Type, level int, fs *[]reflect.StructField, name string, tag string, isStruct bool) {
	fmt.Println(repeat(level), typ.Kind(), name, tag, isStruct)
	if typ.Kind() == reflect.Ptr {
		visitType(typ.Elem(), level+1, fs, name, tag, false)
	}
	if typ.Kind() != reflect.Struct {
		*fs = append(*fs, reflect.StructField{Name: name, Type: typ, Tag: reflect.StructTag(tag)})
		return
	}
	if typ.Implements(marshalerType) {
		*fs = append(*fs, reflect.StructField{Name: name, Type: typ, Tag: reflect.StructTag(tag)})
		return
	}
	var nfs []reflect.StructField
	for i := 0; i < typ.NumField(); i++ {
		jt := typ.Field(i).Tag.Get("json")
		fn := typ.Field(i).Name
		if jt == "" {
			jt = strings.ToLower(fn)
		}
		jt = fmt.Sprintf(`json:"%s"`, jt)
		// fmt.Println("-----", fn, jt)
		visitType(typ.Field(i).Type, level+1, &nfs, fn, jt, true)
	}

	if isStruct {
		if name == "" {
			*fs = nfs
		} else {
			*fs = append(*fs, reflect.StructField{
				Name: name,
				Type: reflect.StructOf(nfs),
				Tag:  reflect.StructTag(tag),
			})
			fmt.Println(repeat(level), typ.Kind(), name, tag, "build")
		}
	}
}
