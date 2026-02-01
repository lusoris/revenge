# Go errors

> Source: https://pkg.go.dev/errors
> Fetched: 2026-02-01T11:41:23.641355+00:00
> Content-Hash: f5e951ab354c6ce4
> Type: html

---

### Overview ¶

Package errors implements functions to manipulate errors. 

The New function creates errors whose only content is a text message. 

An error e wraps another error if e's type has one of the methods 
    
    
    Unwrap() error
    Unwrap() []error
    

If e.Unwrap() returns a non-nil error w or a slice containing w, then we say that e wraps w. A nil error returned from e.Unwrap() indicates that e does not wrap any error. It is invalid for an Unwrap method to return an []error containing a nil error value. 

An easy way to create wrapped errors is to call [fmt.Errorf](/fmt#Errorf) and apply the %w verb to the error argument: 
    
    
    wrapsErr := fmt.Errorf("... %w ...", ..., err, ...)
    

Successive unwrapping of an error creates a tree. The Is and As functions inspect an error's tree by examining first the error itself followed by the tree of each of its children in turn (pre-order, depth-first traversal). 

See <https://go.dev/blog/go1.13-errors> for a deeper discussion of the philosophy of wrapping and when to wrap. 

Is examines the tree of its first argument looking for an error that matches the second. It reports whether it finds a match. It should be used in preference to simple equality checks: 
    
    
    if errors.Is(err, fs.ErrExist)
    

is preferable to 
    
    
    if err == fs.ErrExist
    

because the former will succeed if err wraps [io/fs.ErrExist](/io/fs#ErrExist). 

As examines the tree of its first argument looking for an error that can be assigned to its second argument, which must be a pointer. If it succeeds, it performs the assignment and returns true. Otherwise, it returns false. The form 
    
    
    var perr *fs.PathError
    if errors.As(err, &perr) {
    	fmt.Println(perr.Path)
    }
    

is preferable to 
    
    
    if perr, ok := err.(*fs.PathError); ok {
    	fmt.Println(perr.Path)
    }
    

because the former will succeed if err wraps an [*io/fs.PathError](/io/fs#PathError). 

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

  * Variables
  * func As(err error, target any) bool
  * func Is(err, target error) bool
  * func Join(errs ...error) error
  * func New(text string) error
  * func Unwrap(err error) error



### Examples ¶

  * Package
  * As
  * Is
  * Join
  * New
  * New (Errorf)
  * Unwrap



### Constants ¶

This section is empty.

### Variables ¶

[View Source](https://cs.opensource.google/go/go/+/go1.25.6:src/errors/errors.go;l=90)
    
    
    var ErrUnsupported = New("unsupported operation")

ErrUnsupported indicates that a requested operation cannot be performed, because it is unsupported. For example, a call to [os.Link](/os#Link) when using a file system that does not support hard links. 

Functions and methods should not return this error but should instead return an error including appropriate context that satisfies 
    
    
    errors.Is(err, errors.ErrUnsupported)
    

either by directly wrapping ErrUnsupported or by implementing an Is method. 

Functions and methods should document the cases in which an error wrapping this will be returned. 

### Functions ¶

####  func [As](https://cs.opensource.google/go/go/+/go1.25.6:src/errors/wrap.go;l=97) ¶ added in go1.13
    
    
    func As(err [error](/builtin#error), target [any](/builtin#any)) [bool](/builtin#bool)

As finds the first error in err's tree that matches target, and if one is found, sets target to that error value and returns true. Otherwise, it returns false. 

The tree consists of err itself, followed by the errors obtained by repeatedly calling its Unwrap() error or Unwrap() []error method. When err wraps multiple errors, As examines err followed by a depth-first traversal of its children. 

An error matches target if the error's concrete value is assignable to the value pointed to by target, or if the error has a method As(any) bool such that As(target) returns true. In the latter case, the As method is responsible for setting target. 

An error type might provide an As method so it can be treated as if it were a different error type. 

As panics if target is not a non-nil pointer to either a type that implements error, or to any interface type. 

Example ¶
    
    
    package main
    
    import (
    	"errors"
    	"fmt"
    	"io/fs"
    	"os"
    )
    
    func main() {
    	if _, err := os.Open("non-existing"); err != nil {
    		var pathError *fs.PathError
    		if errors.As(err, &pathError) {
    			fmt.Println("Failed at path:", pathError.Path)
    		} else {
    			fmt.Println(err)
    		}
    	}
    
    }
    
    
    
    Output:
    
    Failed at path: non-existing
    

Share Format Run

####  func [Is](https://cs.opensource.google/go/go/+/go1.25.6:src/errors/wrap.go;l=44) ¶ added in go1.13
    
    
    func Is(err, target [error](/builtin#error)) [bool](/builtin#bool)

Is reports whether any error in err's tree matches target. 

The tree consists of err itself, followed by the errors obtained by repeatedly calling its Unwrap() error or Unwrap() []error method. When err wraps multiple errors, Is examines err followed by a depth-first traversal of its children. 

An error is considered to match a target if it is equal to that target or if it implements a method Is(error) bool such that Is(target) returns true. 

An error type might provide an Is method so it can be treated as equivalent to an existing error. For example, if MyError defines 
    
    
    func (m MyError) Is(target error) bool { return target == fs.ErrExist }
    

then Is(MyError{}, fs.ErrExist) returns true. See [syscall.Errno.Is](/syscall#Errno.Is) for an example in the standard library. An Is method should only shallowly compare err and the target and not call Unwrap on either. 

Example ¶
    
    
    package main
    
    import (
    	"errors"
    	"fmt"
    	"io/fs"
    	"os"
    )
    
    func main() {
    	if _, err := os.Open("non-existing"); err != nil {
    		if errors.Is(err, fs.ErrNotExist) {
    			fmt.Println("file does not exist")
    		} else {
    			fmt.Println(err)
    		}
    	}
    
    }
    
    
    
    Output:
    
    file does not exist
    

Share Format Run

####  func [Join](https://cs.opensource.google/go/go/+/go1.25.6:src/errors/join.go;l=19) ¶ added in go1.20
    
    
    func Join(errs ...[error](/builtin#error)) [error](/builtin#error)

Join returns an error that wraps the given errors. Any nil error values are discarded. Join returns nil if every value in errs is nil. The error formats as the concatenation of the strings obtained by calling the Error method of each element of errs, with a newline between each string. 

A non-nil error returned by Join implements the Unwrap() []error method. 

Example ¶
    
    
    package main
    
    import (
    	"errors"
    	"fmt"
    )
    
    func main() {
    	err1 := errors.New("err1")
    	err2 := errors.New("err2")
    	err := errors.Join(err1, err2)
    	fmt.Println(err)
    	if errors.Is(err, err1) {
    		fmt.Println("err is err1")
    	}
    	if errors.Is(err, err2) {
    		fmt.Println("err is err2")
    	}
    	fmt.Println(err.(interface{ Unwrap() []error }).Unwrap())
    }
    
    
    
    Output:
    
    err1
    err2
    err is err1
    err is err2
    [err1 err2]
    

Share Format Run

####  func [New](https://cs.opensource.google/go/go/+/go1.25.6:src/errors/errors.go;l=64) ¶
    
    
    func New(text [string](/builtin#string)) [error](/builtin#error)

New returns an error that formats as the given text. Each call to New returns a distinct error value even if the text is identical. 

Example ¶
    
    
    package main
    
    import (
    	"errors"
    	"fmt"
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

####  func [Unwrap](https://cs.opensource.google/go/go/+/go1.25.6:src/errors/wrap.go;l=17) ¶ added in go1.13
    
    
    func Unwrap(err [error](/builtin#error)) [error](/builtin#error)

Unwrap returns the result of calling the Unwrap method on err, if err's type contains an Unwrap method returning error. Otherwise, Unwrap returns nil. 

Unwrap only calls a method of the form "Unwrap() error". In particular Unwrap does not unwrap errors returned by Join. 

Example ¶
    
    
    package main
    
    import (
    	"errors"
    	"fmt"
    )
    
    func main() {
    	err1 := errors.New("error1")
    	err2 := fmt.Errorf("error2: [%w]", err1)
    	fmt.Println(err2)
    	fmt.Println(errors.Unwrap(err2))
    }
    
    
    
    Output:
    
    error2: [error1]
    error1
    

Share Format Run

### Types ¶

This section is empty.
