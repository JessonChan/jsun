# jsun
this is a go json tool like encode/json but allow you to use camelName/CamelName or under_score name style

```go
package main

import (
	"fmt"

	"github.com/JessonChan/jsun"
)

type Person struct {
	Id        int
	FirstName string
	LastName  string `json:"name"`
}

func main() {
	p := Person{
		Id:        1,
		FirstName: "FirstName",
		LastName:  "LastName",
	}
	b, _ := jsun.Marshal(p)
    // or
	// b, _ := jsun.Marshal(p,jsun.LowerCamelStyle)
	fmt.Println(string(b))
}

``` 

output

> {"id":1,"firstName":"FirstName","name":"LastName"}


```go
package main

import (
	"fmt"

	"github.com/JessonChan/jsun"
)

type Person struct {
	Id        int
	FirstName string
	LastName  string `json:"name"`
}

func main() {
	p := Person{
		Id:        1,
		FirstName: "FirstName",
		LastName:  "LastName",
	}
	b, _ := jsun.Marshal(p, jsun.UnderScoreStyle)
	fmt.Println(string(b))
}
```

output

> {"id":1,"first_name":"FirstName","name":"LastName"}


```go
package main

import (
	"fmt"

	"github.com/JessonChan/jsun"
)

type Person struct {
	Id        int
	FirstName string
	LastName  string `json:"name"`
}

func main() {
	p := Person{
		Id:        1,
		FirstName: "FirstName",
		LastName:  "LastName",
	}
	b, _ := jsun.Marshal(p, jsun.UpperCamelStyle)
	fmt.Println(string(b))
}

``` 

output

> {"id":1,"FirstName":"FirstName","name":"LastName"}
