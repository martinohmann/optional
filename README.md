# optional

[![Build Status](https://travis-ci.com/martinohmann/optional.svg?branch=master)](https://travis-ci.com/martinohmann/optional)
[![GoDoc](https://godoc.org/github.com/martinohmann/optional?status.svg)](https://godoc.org/github.com/martinohmann/optional)
![GitHub](https://img.shields.io/github/license/martinohmann/optional?color=orange)

Optional provides a container for values that may or may not contain a
non-nil value. It is intended for use as a method return type where there is a
clear need to represent "no result" and where using nil is likely to cause
errors.

## Installation

```
go get -u github.com/martinohmann/optional
```

## Usage Examples

### Value filtering


```go
value := optional.Of("foo").
  Filter(filterOutFoo).
  OrElse("bar")

fmt.Println(value)

// output:
// bar
```

### Value mapping

```go
value := optional.Of(42).
  Map(addTwo).
  Get()

fmt.Println(value)

// output:
// 44

```

### Presence actions

```go
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

```

### Value receiver by reference

```go
var err error

optional.OfNilable(err).
  OrElseInto(errors.New("some error"), &err)

fmt.Println(err)

// output:
// some error
```

Check out the [`examples_test.go`](examples_test.go) for more usage examples.

## License

The source code of optional is released under the MIT License. See the bundled
LICENSE file for details.
