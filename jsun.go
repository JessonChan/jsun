// All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package jsun

import (
	"errors"
	"reflect"

	"github.com/JessonChan/jsun/json"
)

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
var errUnsupported = errors.New("unsupported json name style")

func SetDefaultStyle(style JsonNameStyle) {
	if style > UnderScoreStyle {
		panic(errUnsupported)
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
				Err:  errUnsupported,
			})
		}
	}
	json.JsonNameConverter = styleNameFunc[style]
	return json.Marshal(v)
}
