package jsun

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"reflect"
	"sync"
)

var typeCache sync.Map

type JsonNameStyle int

var DebugMessage bool

const (
	LowerCamelStyle = JsonNameStyle(0)
	UpperCamelStyle = JsonNameStyle(1)
	UnderScoreStyle = JsonNameStyle(2)
)

var styleNameFunc = []func(string) string{
	func(s string) string {
		bs := []rune(s)
		if 'A' <= bs[0] && bs[0] <= 'z' {
			bs[0] += 'a' - 'A'
		}
		return string(bs)
	},
	func(s string) string { return s },
	func(s string) string {
		bs := make([]rune, 0, 2*len(s))
		for _, s := range s {
			if 'A' <= s && s <= 'Z' {
				s += 'a' - 'A'
				bs = append(bs, '_')
			}
			bs = append(bs, s)
		}
		if bs[0] == '_' {
			s = string(bs[1:])
		} else {
			s = string(bs)
		}
		return s
	}}

var defaultStyle = LowerCamelStyle
var unsupportedErr = errors.New("unsupported json name style")

func SetDefaultStyle(style JsonNameStyle) {
	if style > UnderScoreStyle {
		panic(unsupportedErr)
	}
	defaultStyle = style
}

func Marshal(v interface{}, styles ...JsonNameStyle) ([]byte, error) {
	style := defaultStyle
	if len(styles) > 0 {
		style = styles[0]
		if style > UnderScoreStyle {
			panic(json.MarshalerError{
				Type: reflect.TypeOf(v),
				Err:  unsupportedErr,
			})
		}
	}
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		rv = reflect.Indirect(rv)
	}
	key := fmt.Sprintf("%s%d", rv.Type(), style)
	typ, find := typeCache.Load(key)
	if find == false {
		typ = buildType(rv.Type(), style)
		typeCache.Store(key, typ)
	}
	nv := reflect.New(typ.(reflect.Type))
	copyValue(nv.Elem(), rv)
	return json.Marshal(nv.Interface())
}

// dst 在这个特定的情况下，不会存在指针类型，也不存在其它如不可导出情况
func copyValue(dst, src reflect.Value) {
	for i := 0; i < src.NumField(); i++ {
		name := src.Type().Field(i).Name
		dstField := dst.FieldByName(name)
		field := src.Field(i)
		if field.Kind() == reflect.Ptr {
			field = field.Elem()
		}
		if field.Type() == dstField.Type() {
			dstField.Set(field)
			continue
		}
		if field.Kind() == reflect.Struct {
			copyValue(dstField, field)
		} else {
			dstField.Set(field)
		}
	}
}

func buildType(typ reflect.Type, style JsonNameStyle) reflect.Type {
	var fs []reflect.StructField
	visitType(typ, 0, &fs, "", "", style)
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

func visitType(typ reflect.Type, level int, fs *[]reflect.StructField, name string, tag string, style JsonNameStyle) {
	if DebugMessage {
		log.Println(repeat(level), typ.Kind(), name, tag)
	}
	if typ.Implements(marshalerType) {
		if typ.Kind() == reflect.Ptr {
			typ = typ.Elem()
		}
		*fs = append(*fs, reflect.StructField{Name: name, Type: typ, Tag: reflect.StructTag(tag)})
		return
	}
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if typ.Kind() != reflect.Struct || typ.Implements(marshalerType) {
		*fs = append(*fs, reflect.StructField{Name: name, Type: typ, Tag: reflect.StructTag(tag)})
		return
	}
	var nfs []reflect.StructField
	for i := 0; i < typ.NumField(); i++ {
		jt := typ.Field(i).Tag.Get("json")
		fn := typ.Field(i).Name
		if jt == "" {
			jt = styleNameFunc[style](fn)
		}
		jt = fmt.Sprintf(`json:"%s"`, jt)
		visitType(typ.Field(i).Type, level+1, &nfs, fn, jt, style)
	}

	if name == "" {
		*fs = nfs
	} else {
		*fs = append(*fs, reflect.StructField{
			Name: name,
			Type: reflect.StructOf(nfs),
			Tag:  reflect.StructTag(tag),
		})
		if DebugMessage {
			log.Println(repeat(level), typ.Kind(), name, tag, "build")
		}
	}
}
