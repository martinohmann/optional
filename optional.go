package optional

import "fmt"

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

// FlatMap applies the *Optional-bearing mapper func to the optional value (if
// present) and returns the resulting *Optional value, otherwise returns an
// empty *Optional. Panics if the mapper func returns nil or a value which is
// not of type *Optional.
func (o *Optional) FlatMap(mapper MapFunc) *Optional {
	if o.IsEmpty() {
		return o
	}

	value := mapper(o.value)
	if value == nil {
		panic("optional.FlatMap: map func returned nil value")
	}

	opt, ok := value.(*Optional)
	if !ok {
		panic(fmt.Sprintf("optional.FlatMap: expected map func to return *Optional, got %T", value))
	}

	return opt
}

// Get returns the optional value if present, otherwise it panics.
func (o *Optional) Get() interface{} {
	if o.value == nil {
		panic("opional.Get: nil value")
	}

	return o.value
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

// FlatMap applies the mapper func to the optional value (if present) and
// returns a new *Optional wrapping the result, otherwise returns an empty
// *Optional.
func (o *Optional) Map(mapper MapFunc) *Optional {
	if o.IsEmpty() {
		return o
	}

	return OfNilable(mapper(o.value))
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

// OrElseGet returns the value of the original *Optional if present, otherwise
// returns the value produced by the supplier func.
func (o *Optional) OrElseGet(supplier SupplyFunc) interface{} {
	if o.IsPresent() {
		return o.value
	}

	return supplier()
}

// OrElsePanic returns the value of the original *Optional if present, otherwise
// panics. The optional message parameter can be used to provided a custom
// panic message.
func (o *Optional) OrElsePanic(message ...string) interface{} {
	if o.IsPresent() {
		return o.value
	}

	msg := "nil value"
	if len(message) > 0 {
		msg = message[0]
	}

	panic(fmt.Sprintf("optional.OrElsePanic: %s", msg))
}

// Of returns an *Optional describing the given non-nil value. Panics if
// value is nil.
func Of(value interface{}) *Optional {
	if value == nil {
		panic("optional.Of: nil value")
	}

	return &Optional{value: value}
}

// OfNilable returns an *Optional describing the given value if it is non-nil,
// otherwise returns an empty *Optional.
func OfNilable(value interface{}) *Optional {
	if value == nil {
		return Empty()
	}

	return Of(value)
}

// String implements fmt.Stringer.
func (o *Optional) String() string {
	if o.IsPresent() {
		return fmt.Sprintf("Optional(%#v)", o.value)
	}

	return "Optional.Empty"
}