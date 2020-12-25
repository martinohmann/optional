package optional

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEquals(t *testing.T) {
	assert.True(t, Empty().Equals(Empty()))
	assert.True(t, Of("foo").Equals(Of("foo")))
	assert.False(t, Of("foo").Equals("foo"))
	assert.False(t, Of("foo").Equals(Of("bar")))
}

func TestFilter(t *testing.T) {
	predicate := func(val interface{}) bool {
		return val == "bar"
	}

	assert.Equal(t, Empty(), Of("foo").Filter(predicate))
	assert.Equal(t, Of("bar"), Of("bar").Filter(predicate))
	assert.Equal(t, Empty(), Empty().Filter(predicate))
}

func TestFlatMap(t *testing.T) {
	mapper := func(value interface{}) interface{} {
		switch value {
		case "foo":
			return Of("bar")
		case "bar":
			return Empty()
		case "baz":
			return "qux"
		case "nil-optional":
			var o *Optional
			return o
		case "ptr-optional":
			return &Optional{}
		case "zero-optional":
			var o Optional
			return o
		default:
			return nil
		}
	}

	assert.Equal(t, Of("bar"), Of("foo").FlatMap(mapper))
	assert.Equal(t, Empty(), Of("bar").FlatMap(mapper))
	assert.Equal(t, Of("qux"), Of("baz").FlatMap(mapper))
	assert.Equal(t, Empty(), Of("qux").FlatMap(mapper))
	assert.Equal(t, Empty(), Empty().FlatMap(mapper))
	assert.Equal(t, Empty(), Of("zero-optional").FlatMap(mapper))
	assert.Equal(t, Empty(), Of("ptr-optional").FlatMap(mapper))
	assert.PanicsWithValue(t, "optional.FlatMap: mapper func returned nil *Optional", func() {
		Of("nil-optional").FlatMap(mapper)
	})
}

func TestGet(t *testing.T) {
	assert.Equal(t, "foo", Of("foo").Get())
	assert.PanicsWithValue(t, "optional.Get: optional has no value", func() {
		Empty().Get()
	})
}

func TestGetInto(t *testing.T) {
	assert.Panics(t, func() {
		var s string
		Empty().GetInto(&s)
	})

	var s string
	Of("foo").GetInto(&s)
	assert.Equal(t, "foo", s)

	assert.PanicsWithValue(t, "optional.GetInto: dst must be a pointer type", func() {
		var s string
		Of("foo").GetInto(s)
	})

	assert.PanicsWithValue(t, "optional.GetInto: value of type string is not assignable to type int", func() {
		var i int
		Of("foo").GetInto(&i)
	})

	var is []int
	Of([]int{1, 2, 3}).GetInto(&is)
	assert.Equal(t, []int{1, 2, 3}, is)
}

func TestIfPresent(t *testing.T) {
	var val interface{}
	calls := 0

	action := func(value interface{}) {
		calls++
		val = value
	}

	Empty().IfPresent(action)

	assert.Equal(t, 0, calls)

	Of("foo").IfPresent(action)

	assert.Equal(t, 1, calls)
	assert.Equal(t, "foo", val)
}

func TestIfPresentOrElse(t *testing.T) {
	var val interface{}
	calls, emptyCalls := 0, 0

	action := func(value interface{}) {
		calls++
		val = value
	}

	emptyAction := func() {
		emptyCalls++
	}

	Empty().IfPresentOrElse(action, emptyAction)

	assert.Equal(t, 0, calls)
	assert.Equal(t, 1, emptyCalls)

	Of("foo").IfPresentOrElse(action, emptyAction)

	assert.Equal(t, 1, calls)
	assert.Equal(t, 1, emptyCalls)
	assert.Equal(t, "foo", val)
}

func TestMap(t *testing.T) {
	mapper := func(value interface{}) interface{} {
		switch value {
		case "foo":
			return Of("bar")
		case "bar":
			return "baz"
		case "baz":
			var o *Optional
			return o
		default:
			return nil
		}
	}

	assert.Equal(t, Of(Of("bar")), Of("foo").Map(mapper))
	assert.Equal(t, Of("baz"), Of("bar").Map(mapper))
	assert.Equal(t, Empty(), Of("baz").Map(mapper))
	assert.Equal(t, Empty(), Of("qux").Map(mapper))
	assert.Equal(t, Empty(), Empty().Map(mapper))
}

func TestOf_NilPanics(t *testing.T) {
	assert.PanicsWithValue(t, "optional.Of: value must not be nil", func() { Of(nil) })
	assert.PanicsWithValue(t, "optional.Of: value must not be nil", func() {
		var s *string
		Of(s)
	})
	assert.PanicsWithValue(t, "optional.Of: value must not be nil", func() {
		var f func()
		Of(f)
	})
}

func TestOfNilable(t *testing.T) {
	assert.Equal(t, Empty(), OfNilable(nil))

	var s *string
	assert.Equal(t, Empty(), OfNilable(s))
}

func TestOr(t *testing.T) {
	assert.Equal(t, Of("bar"), Of("bar").Or(func() interface{} { return "foo" }))
	assert.Equal(t, Of("foo"), Empty().Or(func() interface{} { return "foo" }))
	assert.Panics(t, func() { Empty().Or(func() interface{} { return nil }) })
}

func TestOrElse(t *testing.T) {
	assert.Equal(t, "bar", Of("bar").OrElse("foo"))
	assert.Equal(t, "foo", Empty().OrElse("foo"))
}

func TestOrElseInto(t *testing.T) {
	var s string
	Empty().OrElseInto("baz", &s)
	assert.Equal(t, "baz", s)
}

func TestOrElseGet(t *testing.T) {
	assert.Equal(t, "bar", Of("bar").OrElseGet(func() interface{} { return "foo" }))
	assert.Equal(t, "foo", Empty().OrElseGet(func() interface{} { return "foo" }))
	assert.Equal(t, nil, Empty().OrElseGet(func() interface{} { return nil }))
}

func TestOrElseGetInto(t *testing.T) {
	var s string
	Empty().OrElseGetInto(func() interface{} { return "baz" }, &s)
	assert.Equal(t, "baz", s)
}

func TestOrElsePanic(t *testing.T) {
	assert.Equal(t, "bar", Of("bar").OrElsePanic("some message"))
	assert.PanicsWithValue(
		t,
		"some message",
		func() { Empty().OrElsePanic("some message") },
	)
}

func TestOrElsePanicInto(t *testing.T) {
	var s string
	Of("baz").OrElsePanicInto("some message", &s)
	assert.Equal(t, "baz", s)

	assert.Panics(t, func() {
		var s string
		Empty().OrElsePanicInto("some message", &s)
	})
}

func TestString(t *testing.T) {
	v := struct{ name string }{name: "foo"}
	assert.Equal(t, `Optional(&struct { name string }{name:"foo"})`, Of(&v).String())
	assert.Equal(t, `Optional("foo")`, Of("foo").String())
	assert.Equal(t, `Optional(42)`, Of(42).String())
	assert.Equal(t, `Optional.Empty`, Empty().String())
}
