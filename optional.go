package optional

import (
	"fmt"
	"reflect"
)

// empty is a sentinel value for empty optionals.
var emptyOptional = &Optional{}

type (
	// ActionFunc is the type of funcs invoked by IfPresent and IfPresentOrElse
	// if an *Optional has a value.
	ActionFunc func(interface{})

	// EmptyActionFunc is the type of funcs type invoked by IfPresentOrElse if
	// an *Optional has no value.
	EmptyActionFunc func()

	// MapFunc is the type of funcs invoked by FlatMap and Map on the optional
	// value (if present) and returns the mapped value.
	MapFunc func(interface{}) interface{}

	// PredicateFunc is the type of funcs invoked by Filter on the optional
	// value (if present). If it returns true, the value will be kept.
	PredicateFunc func(interface{}) bool

	// SupplyFunc is the type of funcs invoked by Or and OrElseGet if the
	// optional has no value.
	SupplyFunc func() interface{}
)

// Optional is a container type which may or may not have a value.
type Optional struct {
	value interface{}
}

// Empty returns an empty *Optional.
func Empty() *Optional {
	return emptyOptional
}

// Equals returns true if other is equal to o. Equality is implied if:
// 1) other is o (pointer equality)
// 2) other and o are both of type *Optional and have the same value.
func (o *Optional) Equals(other interface{}) bool {
	if other == o {
		return true
	}

	opt, ok := other.(*Optional)
	if !ok {
		return false
	}

	return opt.value == o.value
}

// Filter matches the optional value (if present) against predicate and returns
// an *Optional describing the matched value, otherwise returns an empty
// *Optional.
func (o *Optional) Filter(predicate PredicateFunc) *Optional {
	if o.IsEmpty() || predicate(o.value) {
		return o
	}

	return Empty()
}

// FlatMap applies the mapper func to the optional value (if present) and
// returns a new *Optional wrapping the result, otherwise returns an empty
// *Optional. If result itself is already an *Optional, it will not be wrapped
// again but instead returned as is. Returning a nil *Optional from the mapper
// func will cause a panic.
func (o *Optional) FlatMap(mapper MapFunc) *Optional {
	if o.IsEmpty() {
		return o
	}

	switch val := mapper(o.value).(type) {
	case *Optional:
		if val == nil {
			panic("optional.FlatMap: mapper func returned nil *Optional")
		}

		return val
	default:
		return OfNilable(val)
	}
}

// Get returns the optional value if present, otherwise it panics.
func (o *Optional) Get() interface{} {
	if o.IsEmpty() {
		panic("optional.Get: optional has no value")
	}

	return o.value
}

// GetInto updates dst with the value of the *Optional if present, otherwise it
// panics. If dst is not a pointer or has a different type than the value
// wrapped by the *Optional, GetInto will panic as well.
func (o *Optional) GetInto(dst interface{}) {
	into(o.Get(), dst)
}

// IfPresent invokes action with the optional value if it is present.
func (o *Optional) IfPresent(action ActionFunc) {
	if o.IsPresent() {
		action(o.value)
	}
}

// IfPresentOrElse invokes action with the optional value if it is present,
// otherwise invokes emptyAction.
func (o *Optional) IfPresentOrElse(action ActionFunc, emptyAction EmptyActionFunc) {
	if o.IsPresent() {
		action(o.value)
	} else {
		emptyAction()
	}
}

// IsEmpty returns true if no value is present, otherwise false. IsEmpty is the
// opposite of IsPresent.
func (o *Optional) IsEmpty() bool {
	return o.value == nil
}

// IsPresent returns true if a value is present, otherwise false. IsPresent is
// the opposite of IsEmpty.
func (o *Optional) IsPresent() bool {
	return o.value != nil
}

// Map applies the mapper func to the optional value (if present) and returns a
// new *Optional wrapping the result, otherwise returns an empty *Optional.
func (o *Optional) Map(mapper MapFunc) *Optional {
	if o.IsEmpty() {
		return o
	}

	return OfNilable(mapper(o.value))
}

// Of returns an *Optional describing the given non-nil value. Panics if
// value is nil.
func Of(value interface{}) *Optional {
	if isNil(value) {
		panic("optional.Of: value must not be nil")
	}

	return &Optional{value}
}

// OfNilable returns an *Optional describing the given value if it is non-nil,
// otherwise returns an empty *Optional.
func OfNilable(value interface{}) *Optional {
	if isNil(value) {
		return Empty()
	}

	return &Optional{value}
}

// Or returns the original *Optional if a value is present, otherwise returns
// an *Optional wrapping the result of the supplier func. Panics if the
// supplier func returns nil.
func (o *Optional) Or(supplier SupplyFunc) *Optional {
	if o.IsPresent() {
		return o
	}

	return Of(supplier())
}

// OrElse returns the value of the original *Optional if present, otherwise
// returns other.
func (o *Optional) OrElse(other interface{}) interface{} {
	if o.IsPresent() {
		return o.value
	}

	return other
}

// OrElseInto updates dst with the value of the *Optional if present, otherwise
// with the value of other. If dst is not a pointer or has a different type
// than the value wrapped by the *Optional or other, OrElseInto will panic.
func (o *Optional) OrElseInto(other, dst interface{}) {
	into(o.OrElse(other), dst)
}

// OrElseGet returns the value of the original *Optional if present, otherwise
// returns the value produced by the supplier func.
func (o *Optional) OrElseGet(supplier SupplyFunc) interface{} {
	if o.IsPresent() {
		return o.value
	}

	return supplier()
}

// OrElseGetInto updates dst with the value of the *Optional if present,
// otherwise with the value of produced by the supplier func. If dst is not a
// pointer or has a different type than the value wrapped by the *Optional or
// other, OrElseGetInto will panic.
func (o *Optional) OrElseGetInto(supplier SupplyFunc, dst interface{}) {
	into(o.OrElseGet(supplier), dst)
}

// OrElsePanic returns the value of the original *Optional if present,
// otherwise panics with the given message.
func (o *Optional) OrElsePanic(message string) interface{} {
	if o.IsPresent() {
		return o.value
	}

	panic(message)
}

// OrElsePanicInto updates dst with the value of the *Optional if present,
// otherwise panics with the given message. If dst is not a pointer or has a
// different type than the value wrapped by the *Optional or other,
// OrElsePanicInto will panic as well.
func (o *Optional) OrElsePanicInto(message string, dst interface{}) {
	into(o.OrElsePanic(message), dst)
}

// String implements fmt.Stringer.
func (o *Optional) String() string {
	if o.IsPresent() {
		return fmt.Sprintf("Optional(%#v)", o.value)
	}

	return "Optional.Empty"
}

// isNil returns true if value is a typed or untyped nil value.
func isNil(value interface{}) bool {
	if value == nil {
		return true
	}

	switch reflect.TypeOf(value).Kind() {
	case reflect.Ptr, reflect.Map, reflect.Array, reflect.Chan, reflect.Slice:
		return reflect.ValueOf(value).IsNil()
	default:
		return false
	}
}

// into writes the value of src into dst. If dst is not a pointer or has a
// different type than src, this will panic.
func into(src, dst interface{}) {
	dstVal := reflect.Indirect(reflect.ValueOf(dst))
	srcVal := reflect.ValueOf(src)

	dstVal.Set(srcVal)
}
