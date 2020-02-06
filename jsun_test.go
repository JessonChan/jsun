// All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package jsun

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

type S struct {
	Num  int
	Unit U `json:"unit"`
}

type U struct {
	UnitDesc    string
	UnitCountry string
}

type Date struct {
	timeUnix int64
}

func (d *Date) MarshalJSON() ([]byte, error) {
	return time.Unix(d.timeUnix, 0).MarshalJSON()
}

func Test_visitStruct(t *testing.T) {
	a := struct {
		CarName string
		CarSize S
		Zero    string
		Create  time.Time
		LogDate *Date
	}{
		CarName: "Car",
		CarSize: S{
			Num: 4,
			Unit: U{
				UnitDesc:    "meter",
				UnitCountry: "China",
			},
		},
		Zero:    "zero",
		Create:  time.Now(),
		LogDate: &Date{time.Now().Unix()},
	}
	// json.Marshal(&a)
	// visitStruct(&a)
	bs, err := Marshal(a)
	fmt.Println(err, string(bs))
	bs, err = json.Marshal(a)
	fmt.Println(err, string(bs))
	bs, err = Marshal(a, UnderScoreStyle)
	fmt.Println(err, string(bs))
	bs, err = Marshal(a, UpperCamelStyle)
	fmt.Println(err, string(bs))
}
