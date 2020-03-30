package optional_test

import (
	"errors"
	"fmt"
	"net"

	"github.com/martinohmann/optional"
)

func ExampleOptional_Filter() {
	value := optional.Of("foo").
		Filter(filterOutFoo).
		OrElse("bar")

	fmt.Println(value)

	// output:
	// bar
}

func ExampleOptional_Map() {
	value := optional.Of(42).
		Map(addTwo).
		Get()

	fmt.Println(value)

	// output:
	// 44
}

func ExampleOptional_IfPresentOrElse() {
	optional.Of(42).
		IfPresentOrElse(
			func(value interface{}) {
				fmt.Printf("optional has value: %v\n", value)
			},
			func() {
				fmt.Println("optional is empty")
			},
		)

	// output:
	// optional has value: 42
}

func ExampleOptional_OrElseGet() {
	value := optional.Empty().
		OrElseGet(func() interface{} {
			return "some value"
		})

	fmt.Println(value)

	// output:
	// some value
}

func ExampleOptional_OrElseInto() {
	var err error

	optional.OfNilable(err).
		OrElseInto(errors.New("some error"), &err)

	fmt.Println(err)

	// output:
	// some error
}

func ExampleOfNilable() {
	ips := []string{
		"127.0.0.1",
		"foobar",
		"1.1.1.1",
		"::1",
		"foo:bar:baz",
		"2606:4700:4700::1111",
	}

	for _, ip := range ips {
		value := optional.OfNilable(net.ParseIP(ip)).
			Map(func(ip interface{}) interface{} {
				return ip.(net.IP).String()
			}).
			OrElse("invalid")

		fmt.Println(value)
	}

	// output:
	// 127.0.0.1
	// invalid
	// 1.1.1.1
	// ::1
	// invalid
	// 2606:4700:4700::1111
}

func filterOutFoo(value interface{}) bool {
	return value != "foo"
}

func addTwo(value interface{}) interface{} {
	return value.(int) + 2
}
