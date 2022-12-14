// All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package jsun

import (
	"errors"
)

type JsonNameStyle int

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

var errUnsupported = errors.New("unsupported json name style")

// JsonNameConverter 用来最终确定json的风格，最终在 encode.go 中起作用
var JsonNameConverter = func(name string) string { return name }

func SetDefaultStyle(style JsonNameStyle) {
	if style > UnderScoreStyle {
		panic(errUnsupported)
	}
	JsonNameConverter = styleNameFunc[style]
}
