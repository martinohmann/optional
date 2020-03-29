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
	o1 := Of("foo")
	o2 := Of("bar")
	o3 := Empty()

	predicate := func(val interface{}) bool {
		return val == "bar"
	}

	assert.Equal(t, Empty(), o1.Filter(predicate))
	assert.Equal(t, o2, o2.Filter(predicate))
	assert.Equal(t, Empty(), o3.Filter(predicate))
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
		default:
			return nil
		}
	}

	o1 := Of("foo")
	o2 := Of("bar")
	o3 := Of("baz")
	o4 := Of("qux")
	o5 := Empty()
	o6 := Of("nil-optional")

	assert.Equal(t, Of("bar"), o1.FlatMap(mapper))
	assert.Equal(t, Empty(), o2.FlatMap(mapper))
	assert.Equal(t, Of("qux"), o3.FlatMap(mapper))
	assert.Equal(t, Empty(), o4.FlatMap(mapper))
	assert.Equal(t, Empty(), o5.FlatMap(mapper))
	assert.Panics(t, func() { o6.FlatMap(mapper) })
}

func TestGet(t *testing.T) {
	o1 := Empty()
	o2 := Of("foo")

	assert.Panics(t, func() { o1.Get() })
	assert.Equal(t, "foo", o2.Get())
}

func TestIfPresent(t *testing.T) {
	o1 := Empty()
	o2 := Of("foo")

	var val interface{}
	calls := 0

	action := func(value interface{}) {
		calls++
		val = value
	}

	o1.IfPresent(action)

	assert.Equal(t, 0, calls)

	o2.IfPresent(action)

	assert.Equal(t, 1, calls)
	assert.Equal(t, "foo", val)
}

func TestIfPresentOrElse(t *testing.T) {
	o1 := Empty()
	o2 := Of("foo")

	var val interface{}
	calls, emptyCalls := 0, 0

	action := func(value interface{}) {
		calls++
		val = value
	}

	emptyAction := func() {
		emptyCalls++
	}

	o1.IfPresentOrElse(action, emptyAction)

	assert.Equal(t, 0, calls)
	assert.Equal(t, 1, emptyCalls)

	o2.IfPresentOrElse(action, emptyAction)

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
		default:
			return nil
		}
	}

	o1 := Of("foo")
	o2 := Of("bar")
	o3 := Of("baz")
	o4 := Empty()

	assert.Equal(t, Of(Of("bar")), o1.Map(mapper))
	assert.Equal(t, Of("baz"), o2.Map(mapper))
	assert.Equal(t, Empty(), o3.Map(mapper))
	assert.Equal(t, Empty(), o4.Map(mapper))
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

func TestOrElseGet(t *testing.T) {
	assert.Equal(t, "bar", Of("bar").OrElseGet(func() interface{} { return "foo" }))
	assert.Equal(t, "foo", Empty().OrElseGet(func() interface{} { return "foo" }))
	assert.Equal(t, nil, Empty().OrElseGet(func() interface{} { return nil }))
}

func TestOrElsePanic(t *testing.T) {
	assert.Equal(t, "bar", Of("bar").OrElsePanic("some message"))
	assert.Panics(t, func() { Empty().OrElsePanic() })
	assert.Panics(t, func() { Empty().OrElsePanic("some message") })
}

func TestString(t *testing.T) {
	assert.Equal(t, `Optional("foo")`, Of("foo").String())
	assert.Equal(t, `Optional.Empty`, Empty().String())
}
