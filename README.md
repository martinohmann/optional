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

## Usage

Check out the [example](example_test.go).

## License

The source code of optional is released under the MIT License. See the bundled
LICENSE file for details.
