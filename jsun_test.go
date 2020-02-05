package jsun

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

type S struct {
	Num  int
	Unit U `json:"unitInChina"`
}

type U struct {
	Desc    string
	Country string
}

func Test_visitStruct(t *testing.T) {
	a := struct {
		Name   string
		Size   S
		Zero   string
		Create time.Time
	}{
		Name: "Car",
		Size: S{
			Num: 4,
			Unit: U{
				Desc:    "meter",
				Country: "China",
			},
		},
		Zero:   "zero",
		Create: time.Now(),
	}
	// json.Marshal(&a)
	// visitStruct(&a)
	bs, err := Marshal(a)
	fmt.Println(err, string(bs))
	bs, err = json.Marshal(a)
	fmt.Println(err, string(bs))
}
