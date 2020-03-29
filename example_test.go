package optional_test

import (
	"fmt"
	"net"

	"github.com/martinohmann/optional"
)

func Example() {
	ips := []string{
		"127.0.0.1",
		"foobar",
		"1.1.1.1",
		"::1",
		"foo:bar:baz",
		"2606:4700:4700::1111",
	}

	for _, ip := range ips {
		value := parseIP(ip).
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

func parseIP(s string) *optional.Optional {
	ip := net.ParseIP(s)

	return optional.OfNilable(ip)
}
