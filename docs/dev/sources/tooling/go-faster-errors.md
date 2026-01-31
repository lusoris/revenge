# go-faster/errors

> Source: https://pkg.go.dev/github.com/go-faster/errors
> Fetched: 2026-01-30T23:49:17.337680+00:00
> Content-Hash: c7216c69628ac00b
> Type: html

---

Overview

¶

Package errors implements functions to manipulate errors.

This package expands "errors" with stack traces and explicit error
wrapping.

Example

¶

package main

import (
	"fmt"
	"time"
)

// MyError is an error implementation that includes a time and message.
type MyError struct {
	When time.Time
	What string
}

func (e MyError) Error() string {
	return fmt.Sprintf("%v: %v", e.When, e.What)
}

func oops() error {
	return MyError{
		time.Date(1989, 3, 15, 22, 30, 0, 0, time.UTC),
		"the file system has gone away",
	}
}

func main() {
	if err := oops(); err != nil {
		fmt.Println(err)
	}
}

Output:

1989-03-15 22:30:00 +0000 UTC: the file system has gone away

Share

Format

Run

Index

¶

func As(err error, target interface{}) bool

func DisableTrace()

func Errorf(format string, a ...interface{}) error

func FormatError(f Formatter, s fmt.State, verb rune)

func Into[T error](err error) (val T, ok bool)

func Is(err, target error) bool

func Join(errs ...error) error

func Must[T any](val T, err error) T

func New(text string) error

func Opaque(err error) error

func Trace() bool

func Unwrap(err error) error

func Wrap(err error, message string) error

func Wrapf(err error, format string, a ...interface{}) error

type Formatter

type Frame

func Caller(skip int) Frame

func Cause(err error) (f Frame, r bool)

func (f Frame) Format(p Printer)

func (f Frame) Location() (function, file string, line int)

type Printer

type Wrapper

Examples

¶

Package

As

FormatError

Into

Must

New

New (Errorf)

Constants

¶

This section is empty.

Variables

¶

This section is empty.

Functions

¶

func

As

¶

func As(err

error

, target interface{})

bool

As finds the first error in err's chain that matches target, and if so, sets
target to that error value and returns true. Otherwise, it returns false.

The chain consists of err itself followed by the sequence of errors obtained by
repeatedly calling Unwrap.

An error matches target if the error's concrete value is assignable to the value
pointed to by target, or if the error has a method As(interface{}) bool such that
As(target) returns true. In the latter case, the As method is responsible for
setting target.

An error type might provide an As method so it can be treated as if it were a
different error type.

As panics if target is not a non-nil pointer to either a type that implements
error, or to any interface type.

Example

¶

package main

import (
	"fmt"
	"os"

	"github.com/go-faster/errors"
)

func main() {
	_, err := os.Open("non-existing")
	if err != nil {
		var pathError *os.PathError
		if errors.As(err, &pathError) {
			fmt.Println("Failed at path:", pathError.Path)
		}
	}

}

Output:

Failed at path: non-existing

Share

Format

Run

func

DisableTrace

¶

func DisableTrace()

DisableTrace disables capturing caller frames.

func

Errorf

¶

func Errorf(format

string

, a ...interface{})

error

Errorf creates new error with format.

func

FormatError

¶

func FormatError(f

Formatter

, s

fmt

.

State

, verb

rune

)

FormatError calls the FormatError method of f with an errors.Printer
configured according to s and verb, and writes the result to s.

Example

¶

package main

import (
	"fmt"

	"github.com/go-faster/errors"
)

type MyError2 struct {
	Message string
	frame   errors.Frame
}

func (m *MyError2) Error() string {
	return m.Message
}

func (m *MyError2) Format(f fmt.State, c rune) { // implements fmt.Formatter
	errors.FormatError(m, f, c)
}

func (m *MyError2) FormatError(p errors.Printer) error { // implements errors.Formatter
	p.Print(m.Message)
	if p.Detail() {
		m.frame.Format(p)
	}
	return nil
}

func main() {
	err := &MyError2{Message: "oops", frame: errors.Caller(1)}
	fmt.Printf("%v\n", err)
	fmt.Println()
	fmt.Printf("%+v\n", err)
}

Share

Format

Run

func

Into

¶

added in

v0.6.0

func Into[T

error

](err

error

) (val T, ok

bool

)

Into finds the first error in err's chain that matches target type T, and if so, returns it.

Into is type-safe alternative to As.

Example

¶

package main

import (
	"fmt"
	"os"

	"github.com/go-faster/errors"
)

func main() {
	_, err := os.Open("non-existing")
	if err != nil {
		if pathError, ok := errors.Into[*os.PathError](err); ok {
			fmt.Println("Failed at path:", pathError.Path)
		}
	}

}

Output:

Failed at path: non-existing

Share

Format

Run

func

Is

¶

func Is(err, target

error

)

bool

Is reports whether any error in err's chain matches target.

The chain consists of err itself followed by the sequence of errors obtained by
repeatedly calling Unwrap.

An error is considered to match a target if it is equal to that target or if
it implements a method Is(error) bool such that Is(target) returns true.

An error type might provide an Is method so it can be treated as equivalent
to an existing error. For example, if MyError defines

func (m MyError) Is(target error) bool { return target == fs.ErrExist }

then Is(MyError{}, fs.ErrExist) returns true. See syscall.Errno.Is for
an example in the standard library.

func

Join

¶

added in

v0.7.0

func Join(errs ...

error

)

error

Join returns an error that wraps the given errors.
Any nil error values are discarded.
Join returns nil if every value in errs is nil.
The error formats as the concatenation of the strings obtained
by calling the Error method of each element of errs, with a newline
between each string.

A non-nil error returned by Join implements the Unwrap() []error method.

Available only for go 1.20 or superior.

func

Must

¶

added in

v0.6.0

func Must[T

any

](val T, err

error

) T

Must is a generic helper, like template.Must, that wraps a call to a function returning (T, error)
and panics if the error is non-nil.

Example

¶

package main

import (
	"fmt"
	"net/url"

	"github.com/go-faster/errors"
)

func main() {
	r := errors.Must(url.Parse(`https://google.com`))
	fmt.Println(r.String())

}

Output:

https://google.com

Share

Format

Run

func

New

¶

func New(text

string

)

error

New returns an error that formats as the given text.

The returned error contains a Frame set to the caller's location and
implements Formatter to show this information when printed with details.

Example

¶

package main

import (
	"fmt"

	"github.com/go-faster/errors"
)

func main() {
	err := errors.New("emit macho dwarf: elf header corrupted")
	if err != nil {
		fmt.Print(err)
	}
}

Output:

emit macho dwarf: elf header corrupted

Share

Format

Run

Example (Errorf)

¶

The fmt package's Errorf function lets us use the package's formatting
features to create descriptive error messages.

package main

import (
	"fmt"
)

func main() {
	const name, id = "bimmler", 17
	err := fmt.Errorf("user %q (id %d) not found", name, id)
	if err != nil {
		fmt.Print(err)
	}
}

Output:

user "bimmler" (id 17) not found

Share

Format

Run

func

Opaque

¶

func Opaque(err

error

)

error

Opaque returns an error with the same error formatting as err
but that does not match err and cannot be unwrapped.

func

Trace

¶

func Trace()

bool

Trace reports whether caller stack capture is enabled.

func

Unwrap

¶

func Unwrap(err

error

)

error

Unwrap returns the result of calling the Unwrap method on err, if err's
type contains an Unwrap method returning error.
Otherwise, Unwrap returns nil.

func

Wrap

¶

func Wrap(err

error

, message

string

)

error

Wrap error with message and caller.

func

Wrapf

¶

func Wrapf(err

error

, format

string

, a ...interface{})

error

Wrapf wraps error with formatted message and caller.

Types

¶

type

Formatter

¶

type Formatter interface {

error

// FormatError prints the receiver's first error and returns the next error in

// the error chain, if any.

FormatError(p

Printer

) (next

error

)
}

A Formatter formats error messages.

type

Frame

¶

type Frame struct {

// contains filtered or unexported fields

}

A Frame contains part of a call stack.

func

Caller

¶

func Caller(skip

int

)

Frame

Caller returns a Frame that describes a frame on the caller's stack.
The argument skip is the number of frames to skip over.
Caller(0) returns the frame for the caller of Caller.

func

Cause

¶

added in

v0.6.0

func Cause(err

error

) (f

Frame

, r

bool

)

Cause returns first recorded Frame.

func (Frame)

Format

¶

func (f

Frame

) Format(p

Printer

)

Format prints the stack as error detail.
It should be called from an error's Format implementation
after printing any other error detail.

func (Frame)

Location

¶

added in

v0.6.0

func (f

Frame

) Location() (function, file

string

, line

int

)

Location reports the file, line, and function of a frame.

The returned function may be "" even if file and line are not.

type

Printer

¶

type Printer interface {

// Print appends args to the message output.

Print(args ...interface{})

// Printf writes a formatted string.

Printf(format

string

, args ...interface{})

// Detail reports whether error detail is requested.

// After the first call to Detail, all text written to the Printer

// is formatted as additional detail, or ignored when

// detail has not been requested.

// If Detail returns false, the caller can avoid printing the detail at all.

Detail()

bool

}

A Printer formats error messages.

The most common implementation of Printer is the one provided by package fmt
during Printf (as of Go 1.13). Localization packages such as golang.org/x/text/message
typically provide their own implementations.

type

Wrapper

¶

type Wrapper interface {

// Unwrap returns the next error in the error chain.

// If there is no next error, Unwrap returns nil.

Unwrap()

error

}

A Wrapper provides context around another error.