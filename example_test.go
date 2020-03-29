package optional_test

import (
	"fmt"
	"strings"

	"github.com/martinohmann/optional"
)

type Foo struct {
	Name string
}

func Example() {
	value := optional.Of(&Foo{Name: "foo"}).
		Map(func(value interface{}) interface{} {
			value.(*Foo).Name += "bar"
			return value
		}).
		Filter(func(value interface{}) bool {
			return !strings.Contains(value.(*Foo).Name, "bar")
		}).
		OrElse(&Foo{Name: "qux"})

	fmt.Printf("%#v\n", value)

	// output:
	// &optional_test.Foo{Name:"qux"}
}
