# go-faster/errors

> Source: https://pkg.go.dev/github.com/go-faster/errors
> Fetched: 2026-01-31T10:56:37.168694+00:00
> Content-Hash: 63cc57618d1c983f
> Type: html

---

### Overview ¶

Package errors implements functions to manipulate errors. 

This package expands "errors" with stack traces and explicit error wrapping. 

Example ¶
    
    
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
    

Share Format Run

### Index ¶

  * func As(err error, target interface{}) bool
  * func DisableTrace()
  * func Errorf(format string, a ...interface{}) error
  * func FormatError(f Formatter, s fmt.State, verb rune)
  * func Into[T error](err error) (val T, ok bool)
  * func Is(err, target error) bool
  * func Join(errs ...error) error
  * func Must[T any](val T, err error) T
  * func New(text string) error
  * func Opaque(err error) error
  * func Trace() bool
  * func Unwrap(err error) error
  * func Wrap(err error, message string) error
  * func Wrapf(err error, format string, a ...interface{}) error
  * type Formatter
  * type Frame
  *     * func Caller(skip int) Frame
    * func Cause(err error) (f Frame, r bool)
  *     * func (f Frame) Format(p Printer)
    * func (f Frame) Location() (function, file string, line int)
  * type Printer
  * type Wrapper



### Examples ¶

  * Package
  * As
  * FormatError
  * Into
  * Must
  * New
  * New (Errorf)



### Constants ¶

This section is empty.

### Variables ¶

This section is empty.

### Functions ¶

####  func [As](https://github.com/go-faster/errors/blob/v0.7.1/wrap.go#L134) ¶
    
    
    func As(err [error](/builtin#error), target interface{}) [bool](/builtin#bool)

As finds the first error in err's chain that matches target, and if so, sets target to that error value and returns true. Otherwise, it returns false. 

The chain consists of err itself followed by the sequence of errors obtained by repeatedly calling Unwrap. 

An error matches target if the error's concrete value is assignable to the value pointed to by target, or if the error has a method As(interface{}) bool such that As(target) returns true. In the latter case, the As method is responsible for setting target. 

An error type might provide an As method so it can be treated as if it were a different error type. 

As panics if target is not a non-nil pointer to either a type that implements error, or to any interface type. 

Example ¶
    
    
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
    

Share Format Run

####  func [DisableTrace](https://github.com/go-faster/errors/blob/v0.7.1/trace.go#L32) ¶
    
    
    func DisableTrace()

DisableTrace disables capturing caller frames. 

####  func [Errorf](https://github.com/go-faster/errors/blob/v0.7.1/format.go#L42) ¶
    
    
    func Errorf(format [string](/builtin#string), a ...interface{}) [error](/builtin#error)

Errorf creates new error with format. 

####  func [FormatError](https://github.com/go-faster/errors/blob/v0.7.1/adaptor.go#L17) ¶
    
    
    func FormatError(f Formatter, s [fmt](/fmt).[State](/fmt#State), verb [rune](/builtin#rune))

FormatError calls the FormatError method of f with an errors.Printer configured according to s and verb, and writes the result to s. 

Example ¶
    
    
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
    

Share Format Run

####  func [Into](https://github.com/go-faster/errors/blob/v0.7.1/into.go#L8) ¶ added in v0.6.0
    
    
    func Into[T [error](/builtin#error)](err [error](/builtin#error)) (val T, ok [bool](/builtin#bool))

Into finds the first error in err's chain that matches target type T, and if so, returns it. 

Into is type-safe alternative to As. 

Example ¶
    
    
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
    

Share Format Run

####  func [Is](https://github.com/go-faster/errors/blob/v0.7.1/wrap.go#L114) ¶
    
    
    func Is(err, target [error](/builtin#error)) [bool](/builtin#bool)

Is reports whether any error in err's chain matches target. 

The chain consists of err itself followed by the sequence of errors obtained by repeatedly calling Unwrap. 

An error is considered to match a target if it is equal to that target or if it implements a method Is(error) bool such that Is(target) returns true. 

An error type might provide an Is method so it can be treated as equivalent to an existing error. For example, if MyError defines 
    
    
    func (m MyError) Is(target error) bool { return target == fs.ErrExist }
    

then Is(MyError{}, fs.ErrExist) returns true. See syscall.Errno.Is for an example in the standard library. 

####  func [Join](https://github.com/go-faster/errors/blob/v0.7.1/join_go120.go#L18) ¶ added in v0.7.0
    
    
    func Join(errs ...[error](/builtin#error)) [error](/builtin#error)

Join returns an error that wraps the given errors. Any nil error values are discarded. Join returns nil if every value in errs is nil. The error formats as the concatenation of the strings obtained by calling the Error method of each element of errs, with a newline between each string. 

A non-nil error returned by Join implements the Unwrap() []error method. 

Available only for go 1.20 or superior. 

####  func [Must](https://github.com/go-faster/errors/blob/v0.7.1/must.go#L7) ¶ added in v0.6.0
    
    
    func Must[T [any](/builtin#any)](val T, err [error](/builtin#error)) T

Must is a generic helper, like template.Must, that wraps a call to a function returning (T, error) and panics if the error is non-nil. 

Example ¶
    
    
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
    

Share Format Run

####  func [New](https://github.com/go-faster/errors/blob/v0.7.1/errors.go#L22) ¶
    
    
    func New(text [string](/builtin#string)) [error](/builtin#error)

New returns an error that formats as the given text. 

The returned error contains a Frame set to the caller's location and implements Formatter to show this information when printed with details. 

Example ¶
    
    
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
    

Share Format Run

Example (Errorf) ¶

The fmt package's Errorf function lets us use the package's formatting features to create descriptive error messages. 
    
    
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
    

Share Format Run

####  func [Opaque](https://github.com/go-faster/errors/blob/v0.7.1/wrap.go#L21) ¶
    
    
    func Opaque(err [error](/builtin#error)) [error](/builtin#error)

Opaque returns an error with the same error formatting as err but that does not match err and cannot be unwrapped. 

####  func [Trace](https://github.com/go-faster/errors/blob/v0.7.1/trace.go#L35) ¶
    
    
    func Trace() [bool](/builtin#bool)

Trace reports whether caller stack capture is enabled. 

####  func [Unwrap](https://github.com/go-faster/errors/blob/v0.7.1/wrap.go#L40) ¶
    
    
    func Unwrap(err [error](/builtin#error)) [error](/builtin#error)

Unwrap returns the result of calling the Unwrap method on err, if err's type contains an Unwrap method returning error. Otherwise, Unwrap returns nil. 

####  func [Wrap](https://github.com/go-faster/errors/blob/v0.7.1/wrap.go#L81) ¶
    
    
    func Wrap(err [error](/builtin#error), message [string](/builtin#string)) [error](/builtin#error)

Wrap error with message and caller. 

####  func [Wrapf](https://github.com/go-faster/errors/blob/v0.7.1/wrap.go#L90) ¶
    
    
    func Wrapf(err [error](/builtin#error), format [string](/builtin#string), a ...interface{}) [error](/builtin#error)

Wrapf wraps error with formatted message and caller. 

### Types ¶

####  type [Formatter](https://github.com/go-faster/errors/blob/v0.7.1/format.go#L13) ¶
    
    
    type Formatter interface {
    	[error](/builtin#error)
    
    	// FormatError prints the receiver's first error and returns the next error in
    	// the error chain, if any.
    	FormatError(p Printer) (next [error](/builtin#error))
    }

A Formatter formats error messages. 

####  type [Frame](https://github.com/go-faster/errors/blob/v0.7.1/frame.go#L12) ¶
    
    
    type Frame struct {
    	// contains filtered or unexported fields
    }

A Frame contains part of a call stack. 

####  func [Caller](https://github.com/go-faster/errors/blob/v0.7.1/frame.go#L22) ¶
    
    
    func Caller(skip [int](/builtin#int)) Frame

Caller returns a Frame that describes a frame on the caller's stack. The argument skip is the number of frames to skip over. Caller(0) returns the frame for the caller of Caller. 

####  func [Cause](https://github.com/go-faster/errors/blob/v0.7.1/wrap.go#L45) ¶ added in v0.6.0
    
    
    func Cause(err [error](/builtin#error)) (f Frame, r [bool](/builtin#bool))

Cause returns first recorded Frame. 

####  func (Frame) [Format](https://github.com/go-faster/errors/blob/v0.7.1/frame.go#L46) ¶
    
    
    func (f Frame) Format(p Printer)

Format prints the stack as error detail. It should be called from an error's Format implementation after printing any other error detail. 

####  func (Frame) [Location](https://github.com/go-faster/errors/blob/v0.7.1/frame.go#L31) ¶ added in v0.6.0
    
    
    func (f Frame) Location() (function, file [string](/builtin#string), line [int](/builtin#int))

Location reports the file, line, and function of a frame. 

The returned function may be "" even if file and line are not. 

####  type [Printer](https://github.com/go-faster/errors/blob/v0.7.1/format.go#L26) ¶
    
    
    type Printer interface {
    	// Print appends args to the message output.
    	Print(args ...interface{})
    
    	// Printf writes a formatted string.
    	Printf(format [string](/builtin#string), args ...interface{})
    
    	// Detail reports whether error detail is requested.
    	// After the first call to Detail, all text written to the Printer
    	// is formatted as additional detail, or ignored when
    	// detail has not been requested.
    	// If Detail returns false, the caller can avoid printing the detail at all.
    	Detail() [bool](/builtin#bool)
    }

A Printer formats error messages. 

The most common implementation of Printer is the one provided by package fmt during Printf (as of Go 1.13). Localization packages such as golang.org/x/text/message typically provide their own implementations. 

####  type [Wrapper](https://github.com/go-faster/errors/blob/v0.7.1/wrap.go#L13) ¶
    
    
    type Wrapper interface {
    	// Unwrap returns the next error in the error chain.
    	// If there is no next error, Unwrap returns nil.
    	Unwrap() [error](/builtin#error)
    }

A Wrapper provides context around another error. 
