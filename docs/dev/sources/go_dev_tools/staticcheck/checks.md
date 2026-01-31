# staticcheck Checks

> Source: https://staticcheck.io/docs/checks/
> Fetched: 2026-01-31T11:06:28.937433+00:00
> Content-Hash: 85b17d2b15fd4655
> Type: html

---

# Checks

Explanations for all checks in Staticcheck

Check| Short description  
---|---  
SA| `staticcheck`  
SA1| Various misuses of the standard library  
SA1000| Invalid regular expression  
SA1001| Invalid template  
SA1002| Invalid format in `time.Parse`  
SA1003| Unsupported argument to functions in `encoding/binary`  
SA1004| Suspiciously small untyped constant in `time.Sleep`  
SA1005| Invalid first argument to `exec.Command`  
SA1006| `Printf` with dynamic first argument and no further arguments  
SA1007| Invalid URL in `net/url.Parse`  
SA1008| Non-canonical key in `http.Header` map  
SA1010| `(*regexp.Regexp).FindAll` called with `n == 0`, which will always return zero results  
SA1011| Various methods in the `strings` package expect valid UTF-8, but invalid input is provided  
SA1012| A nil `context.Context` is being passed to a function, consider using `context.TODO` instead  
SA1013| `io.Seeker.Seek` is being called with the whence constant as the first argument, but it should be the second  
SA1014| Non-pointer value passed to `Unmarshal` or `Decode`  
SA1015| Using `time.Tick` in a way that will leak. Consider using `time.NewTicker`, and only use `time.Tick` in tests, commands and endless functions  
SA1016| Trapping a signal that cannot be trapped  
SA1017| Channels used with `os/signal.Notify` should be buffered  
SA1018| `strings.Replace` called with `n == 0`, which does nothing  
SA1019| Using a deprecated function, variable, constant or field  
SA1020| Using an invalid host:port pair with a `net.Listen`-related function  
SA1021| Using `bytes.Equal` to compare two `net.IP`  
SA1023| Modifying the buffer in an `io.Writer` implementation  
SA1024| A string cutset contains duplicate characters  
SA1025| It is not possible to use `(*time.Timer).Reset`’s return value correctly  
SA1026| Cannot marshal channels or functions  
SA1027| Atomic access to 64-bit variable must be 64-bit aligned  
SA1028| `sort.Slice` can only be used on slices  
SA1029| Inappropriate key in call to `context.WithValue`  
SA1030| Invalid argument in call to a `strconv` function  
SA1031| Overlapping byte slices passed to an encoder  
SA1032| Wrong order of arguments to `errors.Is`  
SA2| Concurrency issues  
SA2000| `sync.WaitGroup.Add` called inside the goroutine, leading to a race condition  
SA2001| Empty critical section, did you mean to defer the unlock?  
SA2002| Called `testing.T.FailNow` or `SkipNow` in a goroutine, which isn’t allowed  
SA2003| Deferred `Lock` right after locking, likely meant to defer `Unlock` instead  
SA3| Testing issues  
SA3000| `TestMain` doesn’t call `os.Exit`, hiding test failures  
SA3001| Assigning to `b.N` in benchmarks distorts the results  
SA4| Code that isn't really doing anything  
SA4000| Binary operator has identical expressions on both sides  
SA4001| `&*x` gets simplified to `x`, it does not copy `x`  
SA4003| Comparing unsigned values against negative values is pointless  
SA4004| The loop exits unconditionally after one iteration  
SA4005| Field assignment that will never be observed. Did you mean to use a pointer receiver?  
SA4006| A value assigned to a variable is never read before being overwritten. Forgotten error check or dead code?  
SA4008| The variable in the loop condition never changes, are you incrementing the wrong variable?  
SA4009| A function argument is overwritten before its first use  
SA4010| The result of `append` will never be observed anywhere  
SA4011| Break statement with no effect. Did you mean to break out of an outer loop?  
SA4012| Comparing a value against NaN even though no value is equal to NaN  
SA4013| Negating a boolean twice (`!!b`) is the same as writing `b`. This is either redundant, or a typo.  
SA4014| An if/else if chain has repeated conditions and no side-effects; if the condition didn’t match the first time, it won’t match the second time, either  
SA4015| Calling functions like `math.Ceil` on floats converted from integers doesn’t do anything useful  
SA4016| Certain bitwise operations, such as `x ^ 0`, do not do anything useful  
SA4017| Discarding the return values of a function without side effects, making the call pointless  
SA4018| Self-assignment of variables  
SA4019| Multiple, identical build constraints in the same file  
SA4020| Unreachable case clause in a type switch  
SA4021| `x = append(y)` is equivalent to `x = y`  
SA4022| Comparing the address of a variable against nil  
SA4023| Impossible comparison of interface value with untyped nil  
SA4024| Checking for impossible return value from a builtin function  
SA4025| Integer division of literals that results in zero  
SA4026| Go constants cannot express negative zero  
SA4027| `(*net/url.URL).Query` returns a copy, modifying it doesn’t change the URL  
SA4028| `x % 1` is always zero  
SA4029| Ineffective attempt at sorting slice  
SA4030| Ineffective attempt at generating random number  
SA4031| Checking never-nil value against nil  
SA4032| Comparing `runtime.GOOS` or `runtime.GOARCH` against impossible value  
SA5| Correctness issues  
SA5000| Assignment to nil map  
SA5001| Deferring `Close` before checking for a possible error  
SA5002| The empty for loop (`for {}`) spins and can block the scheduler  
SA5003| Defers in infinite loops will never execute  
SA5004| `for { select { ...` with an empty default branch spins  
SA5005| The finalizer references the finalized object, preventing garbage collection  
SA5007| Infinite recursive call  
SA5008| Invalid struct tag  
SA5009| Invalid Printf call  
SA5010| Impossible type assertion  
SA5011| Possible nil pointer dereference  
SA5012| Passing odd-sized slice to function expecting even size  
SA6| Performance issues  
SA6000| Using `regexp.Match` or related in a loop, should use `regexp.Compile`  
SA6001| Missing an optimization opportunity when indexing maps by byte slices  
SA6002| Storing non-pointer values in `sync.Pool` allocates memory  
SA6003| Converting a string to a slice of runes before ranging over it  
SA6005| Inefficient string comparison with `strings.ToLower` or `strings.ToUpper`  
SA6006| Using io.WriteString to write `[]byte`  
SA9| Dubious code constructs that have a high probability of being wrong  
SA9001| Defers in range loops may not run when you expect them to  
SA9002| Using a non-octal `os.FileMode` that looks like it was meant to be in octal.  
SA9003| Empty body in an if or else branch  
SA9004| Only the first constant has an explicit type  
SA9005| Trying to marshal a struct with no public fields nor custom marshaling  
SA9006| Dubious bit shifting of a fixed size integer value  
SA9007| Deleting a directory that shouldn’t be deleted  
SA9008| `else` branch of a type assertion is probably not reading the right value  
SA9009| Ineffectual Go compiler directive  
S| `simple`  
S1| Code simplifications  
S1000| Use plain channel send or receive instead of single-case select  
S1001| Replace for loop with call to copy  
S1002| Omit comparison with boolean constant  
S1003| Replace call to `strings.Index` with `strings.Contains`  
S1004| Replace call to `bytes.Compare` with `bytes.Equal`  
S1005| Drop unnecessary use of the blank identifier  
S1006| Use `for { ... }` for infinite loops  
S1007| Simplify regular expression by using raw string literal  
S1008| Simplify returning boolean expression  
S1009| Omit redundant nil check on slices, maps, and channels  
S1010| Omit default slice index  
S1011| Use a single `append` to concatenate two slices  
S1012| Replace `time.Now().Sub(x)` with `time.Since(x)`  
S1016| Use a type conversion instead of manually copying struct fields  
S1017| Replace manual trimming with `strings.TrimPrefix`  
S1018| Use `copy` for sliding elements  
S1019| Simplify `make` call by omitting redundant arguments  
S1020| Omit redundant nil check in type assertion  
S1021| Merge variable declaration and assignment  
S1023| Omit redundant control flow  
S1024| Replace `x.Sub(time.Now())` with `time.Until(x)`  
S1025| Don’t use `fmt.Sprintf("%s", x)` unnecessarily  
S1028| Simplify error construction with `fmt.Errorf`  
S1029| Range over the string directly  
S1030| Use `bytes.Buffer.String` or `bytes.Buffer.Bytes`  
S1031| Omit redundant nil check around loop  
S1032| Use `sort.Ints(x)`, `sort.Float64s(x)`, and `sort.Strings(x)`  
S1033| Unnecessary guard around call to `delete`  
S1034| Use result of type assertion to simplify cases  
S1035| Redundant call to `net/http.CanonicalHeaderKey` in method call on `net/http.Header`  
S1036| Unnecessary guard around map access  
S1037| Elaborate way of sleeping  
S1038| Unnecessarily complex way of printing formatted string  
S1039| Unnecessary use of `fmt.Sprint`  
S1040| Type assertion to current type  
ST| `stylecheck`  
ST1| Stylistic issues  
ST1000| Incorrect or missing package comment  
ST1001| Dot imports are discouraged  
ST1003| Poorly chosen identifier  
ST1005| Incorrectly formatted error string  
ST1006| Poorly chosen receiver name  
ST1008| A function’s error value should be its last return value  
ST1011| Poorly chosen name for variable of type `time.Duration`  
ST1012| Poorly chosen name for error variable  
ST1013| Should use constants for HTTP error codes, not magic numbers  
ST1015| A switch’s default case should be the first or last case  
ST1016| Use consistent method receiver names  
ST1017| Don’t use Yoda conditions  
ST1018| Avoid zero-width and control characters in string literals  
ST1019| Importing the same package multiple times  
ST1020| The documentation of an exported function should start with the function’s name  
ST1021| The documentation of an exported type should start with type’s name  
ST1022| The documentation of an exported variable or constant should start with variable’s name  
ST1023| Redundant type in variable declaration  
QF| `quickfix`  
QF1| Quickfixes  
QF1001| Apply De Morgan’s law  
QF1002| Convert untagged switch to tagged switch  
QF1003| Convert if/else-if chain to tagged switch  
QF1004| Use `strings.ReplaceAll` instead of `strings.Replace` with `n == -1`  
QF1005| Expand call to `math.Pow`  
QF1006| Lift `if`+`break` into loop condition  
QF1007| Merge conditional assignment into variable declaration  
QF1008| Omit embedded fields from selector expression  
QF1009| Use `time.Time.Equal` instead of `==` operator  
QF1010| Convert slice of bytes to string when printing it  
QF1011| Omit redundant type from variable declaration  
QF1012| Use `fmt.Fprintf(x, ...)` instead of `x.Write(fmt.Sprintf(...))`  
  
## SA – `staticcheck`

The SA category of checks, codenamed `staticcheck`, contains all checks that are concerned with the correctness of code.

### SA1 – Various misuses of the standard library

Checks in this category deal with misuses of the standard library. This tends to involve incorrect function arguments or violating other invariants laid out by the standard library's documentation.

#### SA1000 - Invalid regular expression

Available since
    2017.1

#### SA1001 - Invalid template

Available since
    2017.1

#### SA1002 - Invalid format in `time.Parse`

Available since
    2017.1

#### SA1003 - Unsupported argument to functions in `encoding/binary`

The `encoding/binary` package can only serialize types with known sizes. This precludes the use of the `int` and `uint` types, as their sizes differ on different architectures. Furthermore, it doesn’t support serializing maps, channels, strings, or functions.

Before Go 1.8, `bool` wasn’t supported, either.

Available since
    2017.1

#### SA1004 - Suspiciously small untyped constant in `time.Sleep`

The `time`.Sleep function takes a `time.Duration` as its only argument. Durations are expressed in nanoseconds. Thus, calling `time.Sleep(1)` will sleep for 1 nanosecond. This is a common source of bugs, as sleep functions in other languages often accept seconds or milliseconds.

The `time` package provides constants such as `time.Second` to express large durations. These can be combined with arithmetic to express arbitrary durations, for example `5 * time.Second` for 5 seconds.

If you truly meant to sleep for a tiny amount of time, use `n * time.Nanosecond` to signal to Staticcheck that you did mean to sleep for some amount of nanoseconds.

Available since
    2017.1

#### SA1005 - Invalid first argument to `exec.Command`

`os/exec` runs programs directly (using variants of the fork and exec system calls on Unix systems). This shouldn’t be confused with running a command in a shell. The shell will allow for features such as input redirection, pipes, and general scripting. The shell is also responsible for splitting the user’s input into a program name and its arguments. For example, the equivalent to
    
    
    ls / /tmp
    

would be
    
    
    exec.Command("ls", "/", "/tmp")
    

If you want to run a command in a shell, consider using something like the following – but be aware that not all systems, particularly Windows, will have a `/bin/sh` program:
    
    
    exec.Command("/bin/sh", "-c", "ls | grep Awesome")
    

Available since
    2017.1

#### SA1006 - `Printf` with dynamic first argument and no further arguments

Using `fmt.Printf` with a dynamic first argument can lead to unexpected output. The first argument is a format string, where certain character combinations have special meaning. If, for example, a user were to enter a string such as
    
    
    Interest rate: 5%
    

and you printed it with
    
    
    fmt.Printf(s)
    

it would lead to the following output:
    
    
    Interest rate: 5%!(NOVERB).
    

Similarly, forming the first parameter via string concatenation with user input should be avoided for the same reason. When printing user input, either use a variant of `fmt.Print`, or use the `%s` Printf verb and pass the string as an argument.

Available since
    2017.1

#### SA1007 - Invalid URL in `net/url.Parse`

Available since
    2017.1

#### SA1008 - Non-canonical key in `http.Header` map

Keys in `http.Header` maps are canonical, meaning they follow a specific combination of uppercase and lowercase letters. Methods such as `http.Header.Add` and `http.Header.Del` convert inputs into this canonical form before manipulating the map.

When manipulating `http.Header` maps directly, as opposed to using the provided methods, care should be taken to stick to canonical form in order to avoid inconsistencies. The following piece of code demonstrates one such inconsistency:
    
    
    h := http.Header{}
    h["etag"] = []string{"1234"}
    h.Add("etag", "5678")
    fmt.Println(h)
    
    // Output:
    // map[Etag:[5678] etag:[1234]]
    

The easiest way of obtaining the canonical form of a key is to use `http.CanonicalHeaderKey`.

Available since
    2017.1

#### SA1010 - `(*regexp.Regexp).FindAll` called with `n == 0`, which will always return zero results

If `n >= 0`, the function returns at most `n` matches/submatches. To return all results, specify a negative number.

Available since
    2017.1

#### SA1011 - Various methods in the `strings` package expect valid UTF-8, but invalid input is provided

Available since
    2017.1

#### SA1012 - A nil `context.Context` is being passed to a function, consider using `context.TODO` instead

Available since
    2017.1

#### SA1013 - `io.Seeker.Seek` is being called with the whence constant as the first argument, but it should be the second

Available since
    2017.1

#### SA1014 - Non-pointer value passed to `Unmarshal` or `Decode`

Available since
    2017.1

#### SA1015 - Using `time.Tick` in a way that will leak. Consider using `time.NewTicker`, and only use `time.Tick` in tests, commands and endless functions

Before Go 1.23, `time.Ticker`s had to be closed to be able to be garbage collected. Since `time.Tick` doesn’t make it possible to close the underlying ticker, using it repeatedly would leak memory.

Go 1.23 fixes this by allowing tickers to be collected even if they weren’t closed.

Available since
    2017.1

#### SA1016 - Trapping a signal that cannot be trapped

Not all signals can be intercepted by a process. Specifically, on UNIX-like systems, the `syscall.SIGKILL` and `syscall.SIGSTOP` signals are never passed to the process, but instead handled directly by the kernel. It is therefore pointless to try and handle these signals.

Available since
    2017.1

#### SA1017 - Channels used with `os/signal.Notify` should be buffered

The `os/signal` package uses non-blocking channel sends when delivering signals. If the receiving end of the channel isn’t ready and the channel is either unbuffered or full, the signal will be dropped. To avoid missing signals, the channel should be buffered and of the appropriate size. For a channel used for notification of just one signal value, a buffer of size 1 is sufficient.

Available since
    2017.1

#### SA1018 - `strings.Replace` called with `n == 0`, which does nothing

With `n == 0`, zero instances will be replaced. To replace all instances, use a negative number, or use `strings.ReplaceAll`.

Available since
    2017.1

#### SA1019 - Using a deprecated function, variable, constant or field

Available since
    2017.1

#### SA1020 - Using an invalid host:port pair with a `net.Listen`-related function

Available since
    2017.1

#### SA1021 - Using `bytes.Equal` to compare two `net.IP`

A `net.IP` stores an IPv4 or IPv6 address as a slice of bytes. The length of the slice for an IPv4 address, however, can be either 4 or 16 bytes long, using different ways of representing IPv4 addresses. In order to correctly compare two `net.IP`s, the `net.IP.Equal` method should be used, as it takes both representations into account.

Available since
    2017.1

#### SA1023 - Modifying the buffer in an `io.Writer` implementation

`Write` must not modify the slice data, even temporarily.

Available since
    2017.1

#### SA1024 - A string cutset contains duplicate characters

The `strings.TrimLeft` and `strings.TrimRight` functions take cutsets, not prefixes. A cutset is treated as a set of characters to remove from a string. For example,
    
    
    strings.TrimLeft("42133word", "1234")
    

will result in the string `"word"` – any characters that are 1, 2, 3 or 4 are cut from the left of the string.

In order to remove one string from another, use `strings.TrimPrefix` instead.

Available since
    2017.1

#### SA1025 - It is not possible to use `(*time.Timer).Reset`’s return value correctly

Available since
    2019.1

#### SA1026 - Cannot marshal channels or functions

Available since
    2019.2

#### SA1027 - Atomic access to 64-bit variable must be 64-bit aligned

On ARM, x86-32, and 32-bit MIPS, it is the caller’s responsibility to arrange for 64-bit alignment of 64-bit words accessed atomically. The first word in a variable or in an allocated struct, array, or slice can be relied upon to be 64-bit aligned.

You can use the structlayout tool to inspect the alignment of fields in a struct.

Available since
    2019.2

#### SA1028 - `sort.Slice` can only be used on slices

The first argument of `sort.Slice` must be a slice.

Available since
    2020.1

#### SA1029 - Inappropriate key in call to `context.WithValue`

The provided key must be comparable and should not be of type `string` or any other built-in type to avoid collisions between packages using context. Users of `WithValue` should define their own types for keys.

To avoid allocating when assigning to an `interface{}`, context keys often have concrete type `struct{}`. Alternatively, exported context key variables’ static type should be a pointer or interface.

Available since
    2020.1

#### SA1030 - Invalid argument in call to a `strconv` function

This check validates the format, number base and bit size arguments of the various parsing and formatting functions in `strconv`.

Available since
    2021.1

#### SA1031 - Overlapping byte slices passed to an encoder

In an encoding function of the form `Encode(dst, src)`, `dst` and `src` were found to reference the same memory. This can result in `src` bytes being overwritten before they are read, when the encoder writes more than one byte per `src` byte.

Available since
    2024.1

#### SA1032 - Wrong order of arguments to `errors.Is`

The first argument of the function `errors.Is` is the error that we have and the second argument is the error we’re trying to match against. For example:
    
    
    if errors.Is(err, io.EOF) { ... }
    

This check detects some cases where the two arguments have been swapped. It flags any calls where the first argument is referring to a package-level error variable, such as
    
    
    if errors.Is(io.EOF, err) { /* this is wrong */ }

Available since
    2024.1

### SA2 – Concurrency issues

Checks in this category find concurrency bugs.

#### SA2000 - `sync.WaitGroup.Add` called inside the goroutine, leading to a race condition

Available since
    2017.1

#### SA2001 - Empty critical section, did you mean to defer the unlock?

Empty critical sections of the kind
    
    
    mu.Lock()
    mu.Unlock()
    

are very often a typo, and the following was intended instead:
    
    
    mu.Lock()
    defer mu.Unlock()
    

Do note that sometimes empty critical sections can be useful, as a form of signaling to wait on another goroutine. Many times, there are simpler ways of achieving the same effect. When that isn’t the case, the code should be amply commented to avoid confusion. Combining such comments with a `//lint:ignore` directive can be used to suppress this rare false positive.

Available since
    2017.1

#### SA2002 - Called `testing.T.FailNow` or `SkipNow` in a goroutine, which isn’t allowed

Available since
    2017.1

#### SA2003 - Deferred `Lock` right after locking, likely meant to defer `Unlock` instead

Available since
    2017.1

### SA3 – Testing issues

Checks in this category find issues in tests and benchmarks.

#### SA3000 - `TestMain` doesn’t call `os.Exit`, hiding test failures

Test executables (and in turn `go test`) exit with a non-zero status code if any tests failed. When specifying your own `TestMain` function, it is your responsibility to arrange for this, by calling `os.Exit` with the correct code. The correct code is returned by `(*testing.M).Run`, so the usual way of implementing `TestMain` is to end it with `os.Exit(m.Run())`.

Available since
    2017.1

#### SA3001 - Assigning to `b.N` in benchmarks distorts the results

The testing package dynamically sets `b.N` to improve the reliability of benchmarks and uses it in computations to determine the duration of a single operation. Benchmark code must not alter `b.N` as this would falsify results.

Available since
    2017.1

### SA4 – Code that isn't really doing anything

Checks in this category point out code that doesn't have any meaningful effect on a program's execution. Usually this means that the programmer thought the code would do one thing while in reality it does something else.

#### SA4000 - Binary operator has identical expressions on both sides

Available since
    2017.1

#### SA4001 - `&*x` gets simplified to `x`, it does not copy `x`

Available since
    2017.1

#### SA4003 - Comparing unsigned values against negative values is pointless

Available since
    2017.1

#### SA4004 - The loop exits unconditionally after one iteration

Available since
    2017.1

#### SA4005 - Field assignment that will never be observed. Did you mean to use a pointer receiver?

Available since
    2021.1

#### SA4006 - A value assigned to a variable is never read before being overwritten. Forgotten error check or dead code?

Available since
    2017.1

#### SA4008 - The variable in the loop condition never changes, are you incrementing the wrong variable?

For example:
    
    
    for i := 0; i < 10; j++ { ... }
    

This may also occur when a loop can only execute once because of unconditional control flow that terminates the loop. For example, when a loop body contains an unconditional break, return, or panic:
    
    
    func f() {
    	panic("oops")
    }
    func g() {
    	for i := 0; i < 10; i++ {
    		// f unconditionally calls panic, which means "i" is
    		// never incremented.
    		f()
    	}
    }

Available since
    2017.1

#### SA4009 - A function argument is overwritten before its first use

Available since
    2017.1

#### SA4010 - The result of `append` will never be observed anywhere

Available since
    2017.1

#### SA4011 - Break statement with no effect. Did you mean to break out of an outer loop?

Available since
    2017.1

#### SA4012 - Comparing a value against NaN even though no value is equal to NaN

Available since
    2017.1

#### SA4013 - Negating a boolean twice (`!!b`) is the same as writing `b`. This is either redundant, or a typo.

Available since
    2017.1

#### SA4014 - An if/else if chain has repeated conditions and no side-effects; if the condition didn’t match the first time, it won’t match the second time, either

Available since
    2017.1

#### SA4015 - Calling functions like `math.Ceil` on floats converted from integers doesn’t do anything useful

Available since
    2017.1

#### SA4016 - Certain bitwise operations, such as `x ^ 0`, do not do anything useful

Available since
    2017.1

#### SA4017 - Discarding the return values of a function without side effects, making the call pointless

Available since
    2017.1

#### SA4018 - Self-assignment of variables

Available since
    2017.1

#### SA4019 - Multiple, identical build constraints in the same file

Available since
    2017.1

#### SA4020 - Unreachable case clause in a type switch

In a type switch like the following
    
    
    type T struct{}
    func (T) Read(b []byte) (int, error) { return 0, nil }
    
    var v interface{} = T{}
    
    switch v.(type) {
    case io.Reader:
        // ...
    case T:
        // unreachable
    }
    

the second case clause can never be reached because `T` implements `io.Reader` and case clauses are evaluated in source order.

Another example:
    
    
    type T struct{}
    func (T) Read(b []byte) (int, error) { return 0, nil }
    func (T) Close() error { return nil }
    
    var v interface{} = T{}
    
    switch v.(type) {
    case io.Reader:
        // ...
    case io.ReadCloser:
        // unreachable
    }
    

Even though `T` has a `Close` method and thus implements `io.ReadCloser`, `io.Reader` will always match first. The method set of `io.Reader` is a subset of `io.ReadCloser`. Thus it is impossible to match the second case without matching the first case.

#### Structurally equivalent interfaces

A special case of the previous example are structurally identical interfaces. Given these declarations
    
    
    type T error
    type V error
    
    func doSomething() error {
        err, ok := doAnotherThing()
        if ok {
            return T(err)
        }
    
        return U(err)
    }
    

the following type switch will have an unreachable case clause:
    
    
    switch doSomething().(type) {
    case T:
        // ...
    case V:
        // unreachable
    }
    

`T` will always match before V because they are structurally equivalent and therefore `doSomething()`’s return value implements both.

Available since
    2019.2

#### SA4021 - `x = append(y)` is equivalent to `x = y`

Available since
    2019.2

#### SA4022 - Comparing the address of a variable against nil

Code such as `if &x == nil` is meaningless, because taking the address of a variable always yields a non-nil pointer.

Available since
    2020.1

#### SA4023 - Impossible comparison of interface value with untyped nil

Under the covers, interfaces are implemented as two elements, a type T and a value V. V is a concrete value such as an int, struct or pointer, never an interface itself, and has type T. For instance, if we store the int value 3 in an interface, the resulting interface value has, schematically, (T=int, V=3). The value V is also known as the interface’s dynamic value, since a given interface variable might hold different values V (and corresponding types T) during the execution of the program.

An interface value is nil only if the V and T are both unset, (T=nil, V is not set), In particular, a nil interface will always hold a nil type. If we store a nil pointer of type *int inside an interface value, the inner type will be *int regardless of the value of the pointer: (T=*int, V=nil). Such an interface value will therefore be non-nil even when the pointer value V inside is nil.

This situation can be confusing, and arises when a nil value is stored inside an interface value such as an error return:
    
    
    func returnsError() error {
        var p *MyError = nil
        if bad() {
            p = ErrBad
        }
        return p // Will always return a non-nil error.
    }
    

If all goes well, the function returns a nil p, so the return value is an error interface value holding (T=*MyError, V=nil). This means that if the caller compares the returned error to nil, it will always look as if there was an error even if nothing bad happened. To return a proper nil error to the caller, the function must return an explicit nil:
    
    
    func returnsError() error {
        if bad() {
            return ErrBad
        }
        return nil
    }
    

It’s a good idea for functions that return errors always to use the error type in their signature (as we did above) rather than a concrete type such as `*MyError`, to help guarantee the error is created correctly. As an example, `os.Open` returns an error even though, if not nil, it’s always of concrete type *os.PathError.

Similar situations to those described here can arise whenever interfaces are used. Just keep in mind that if any concrete value has been stored in the interface, the interface will not be nil. For more information, see The Laws of Reflection at <https://golang.org/doc/articles/laws_of_reflection.html>.

This text has been copied from <https://golang.org/doc/faq#nil_error>, licensed under the Creative Commons Attribution 3.0 License.

Available since
    2020.2

#### SA4024 - Checking for impossible return value from a builtin function

Return values of the `len` and `cap` builtins cannot be negative.

See <https://golang.org/pkg/builtin/#len> and <https://golang.org/pkg/builtin/#cap>.

Example:
    
    
    if len(slice) < 0 {
        fmt.Println("unreachable code")
    }
    

Available since
    2021.1

#### SA4025 - Integer division of literals that results in zero

When dividing two integer constants, the result will also be an integer. Thus, a division such as `2 / 3` results in `0`. This is true for all of the following examples:
    
    
    _ = 2 / 3
    const _ = 2 / 3
    const _ float64 = 2 / 3
    _ = float64(2 / 3)
    

Staticcheck will flag such divisions if both sides of the division are integer literals, as it is highly unlikely that the division was intended to truncate to zero. Staticcheck will not flag integer division involving named constants, to avoid noisy positives.

Available since
    2021.1

#### SA4026 - Go constants cannot express negative zero

In IEEE 754 floating point math, zero has a sign and can be positive or negative. This can be useful in certain numerical code.

Go constants, however, cannot express negative zero. This means that the literals `-0.0` and `0.0` have the same ideal value (zero) and will both represent positive zero at runtime.

To explicitly and reliably create a negative zero, you can use the `math.Copysign` function: `math.Copysign(0, -1)`.

Available since
    2021.1

#### SA4027 - `(*net/url.URL).Query` returns a copy, modifying it doesn’t change the URL

`(*net/url.URL).Query` parses the current value of `net/url.URL.RawQuery` and returns it as a map of type `net/url.Values`. Subsequent changes to this map will not affect the URL unless the map gets encoded and assigned to the URL’s `RawQuery`.

As a consequence, the following code pattern is an expensive no-op: `u.Query().Add(key, value)`.

Available since
    2021.1

#### SA4028 - `x % 1` is always zero

Available since
    2022.1

#### SA4029 - Ineffective attempt at sorting slice

`sort.Float64Slice`, `sort.IntSlice`, and `sort.StringSlice` are types, not functions. Doing `x = sort.StringSlice(x)` does nothing, especially not sort any values. The correct usage is `sort.Sort(sort.StringSlice(x))` or `sort.StringSlice(x).Sort()`, but there are more convenient helpers, namely `sort.Float64s`, `sort.Ints`, and `sort.Strings`.

Available since
    2022.1

#### SA4030 - Ineffective attempt at generating random number

Functions in the `math/rand` package that accept upper limits, such as `Intn`, generate random numbers in the half-open interval [0,n). In other words, the generated numbers will be `>= 0` and `< n` – they don’t include `n`. `rand.Intn(1)` therefore doesn’t generate `0` or `1`, it always generates `0`.

Available since
    2022.1

#### SA4031 - Checking never-nil value against nil

Available since
    2022.1

#### SA4032 - Comparing `runtime.GOOS` or `runtime.GOARCH` against impossible value

Available since
    2024.1

### SA5 – Correctness issues

Checks in this category find assorted bugs and crashes.

#### SA5000 - Assignment to nil map

Available since
    2017.1

#### SA5001 - Deferring `Close` before checking for a possible error

Available since
    2017.1

#### SA5002 - The empty for loop (`for {}`) spins and can block the scheduler

Available since
    2017.1

#### SA5003 - Defers in infinite loops will never execute

Defers are scoped to the surrounding function, not the surrounding block. In a function that never returns, i.e. one containing an infinite loop, defers will never execute.

Available since
    2017.1

#### SA5004 - `for { select { ...` with an empty default branch spins

Available since
    2017.1

#### SA5005 - The finalizer references the finalized object, preventing garbage collection

A finalizer is a function associated with an object that runs when the garbage collector is ready to collect said object, that is when the object is no longer referenced by anything.

If the finalizer references the object, however, it will always remain as the final reference to that object, preventing the garbage collector from collecting the object. The finalizer will never run, and the object will never be collected, leading to a memory leak. That is why the finalizer should instead use its first argument to operate on the object. That way, the number of references can temporarily go to zero before the object is being passed to the finalizer.

Available since
    2017.1

#### SA5007 - Infinite recursive call

A function that calls itself recursively needs to have an exit condition. Otherwise it will recurse forever, until the system runs out of memory.

This issue can be caused by simple bugs such as forgetting to add an exit condition. It can also happen “on purpose”. Some languages have tail call optimization which makes certain infinite recursive calls safe to use. Go, however, does not implement TCO, and as such a loop should be used instead.

Available since
    2017.1

#### SA5008 - Invalid struct tag

Available since
    2019.2

#### SA5009 - Invalid Printf call

Available since
    2019.2

#### SA5010 - Impossible type assertion

Some type assertions can be statically proven to be impossible. This is the case when the method sets of both arguments of the type assertion conflict with each other, for example by containing the same method with different signatures.

The Go compiler already applies this check when asserting from an interface value to a concrete type. If the concrete type misses methods from the interface, or if function signatures don’t match, then the type assertion can never succeed.

This check applies the same logic when asserting from one interface to another. If both interface types contain the same method but with different signatures, then the type assertion can never succeed, either.

Available since
    2020.1

#### SA5011 - Possible nil pointer dereference

A pointer is being dereferenced unconditionally, while also being checked against nil in another place. This suggests that the pointer may be nil and dereferencing it may panic. This is commonly a result of improperly ordered code or missing return statements. Consider the following examples:
    
    
    func fn(x *int) {
        fmt.Println(*x)
    
        // This nil check is equally important for the previous dereference
        if x != nil {
            foo(*x)
        }
    }
    
    func TestFoo(t *testing.T) {
        x := compute()
        if x == nil {
            t.Errorf("nil pointer received")
        }
    
        // t.Errorf does not abort the test, so if x is nil, the next line will panic.
        foo(*x)
    }
    

Staticcheck tries to deduce which functions abort control flow. For example, it is aware that a function will not continue execution after a call to `panic` or `log.Fatal`. However, sometimes this detection fails, in particular in the presence of conditionals. Consider the following example:
    
    
    func Log(msg string, level int) {
        fmt.Println(msg)
        if level == levelFatal {
            os.Exit(1)
        }
    }
    
    func Fatal(msg string) {
        Log(msg, levelFatal)
    }
    
    func fn(x *int) {
        if x == nil {
            Fatal("unexpected nil pointer")
        }
        fmt.Println(*x)
    }
    

Staticcheck will flag the dereference of `x`, even though it is perfectly safe. Staticcheck is not able to deduce that a call to Fatal will exit the program. For the time being, the easiest workaround is to modify the definition of Fatal like so:
    
    
    func Fatal(msg string) {
        Log(msg, levelFatal)
        panic("unreachable")
    }
    

We also hard-code functions from common logging packages such as logrus. Please file an issue if we’re missing support for a popular package.

Available since
    2020.1

#### SA5012 - Passing odd-sized slice to function expecting even size

Some functions that take slices as parameters expect the slices to have an even number of elements. Often, these functions treat elements in a slice as pairs. For example, `strings.NewReplacer` takes pairs of old and new strings, and calling it with an odd number of elements would be an error.

Available since
    2020.2

### SA6 – Performance issues

Checks in this category find code that can be trivially made faster.

#### SA6000 - Using `regexp.Match` or related in a loop, should use `regexp.Compile`

Available since
    2017.1

#### SA6001 - Missing an optimization opportunity when indexing maps by byte slices

Map keys must be comparable, which precludes the use of byte slices. This usually leads to using string keys and converting byte slices to strings.

Normally, a conversion of a byte slice to a string needs to copy the data and causes allocations. The compiler, however, recognizes `m[string(b)]` and uses the data of `b` directly, without copying it, because it knows that the data can’t change during the map lookup. This leads to the counter-intuitive situation that
    
    
    k := string(b)
    println(m[k])
    println(m[k])
    

will be less efficient than
    
    
    println(m[string(b)])
    println(m[string(b)])
    

because the first version needs to copy and allocate, while the second one does not.

For some history on this optimization, check out commit f5f5a8b6209f84961687d993b93ea0d397f5d5bf in the Go repository.

Available since
    2017.1

#### SA6002 - Storing non-pointer values in `sync.Pool` allocates memory

A `sync.Pool` is used to avoid unnecessary allocations and reduce the amount of work the garbage collector has to do.

When passing a value that is not a pointer to a function that accepts an interface, the value needs to be placed on the heap, which means an additional allocation. Slices are a common thing to put in sync.Pools, and they’re structs with 3 fields (length, capacity, and a pointer to an array). In order to avoid the extra allocation, one should store a pointer to the slice instead.

See the comments on <https://go-review.googlesource.com/c/go/+/24371> that discuss this problem.

Available since
    2017.1

#### SA6003 - Converting a string to a slice of runes before ranging over it

You may want to loop over the runes in a string. Instead of converting the string to a slice of runes and looping over that, you can loop over the string itself. That is,
    
    
    for _, r := range s {}
    

and
    
    
    for _, r := range []rune(s) {}
    

will yield the same values. The first version, however, will be faster and avoid unnecessary memory allocations.

Do note that if you are interested in the indices, ranging over a string and over a slice of runes will yield different indices. The first one yields byte offsets, while the second one yields indices in the slice of runes.

Available since
    2017.1

#### SA6005 - Inefficient string comparison with `strings.ToLower` or `strings.ToUpper`

Converting two strings to the same case and comparing them like so
    
    
    if strings.ToLower(s1) == strings.ToLower(s2) {
        ...
    }
    

is significantly more expensive than comparing them with `strings.EqualFold(s1, s2)`. This is due to memory usage as well as computational complexity.

`strings.ToLower` will have to allocate memory for the new strings, as well as convert both strings fully, even if they differ on the very first byte. strings.EqualFold, on the other hand, compares the strings one character at a time. It doesn’t need to create two intermediate strings and can return as soon as the first non-matching character has been found.

For a more in-depth explanation of this issue, see <https://blog.digitalocean.com/how-to-efficiently-compare-strings-in-go/>

Available since
    2019.2

#### SA6006 - Using io.WriteString to write `[]byte`

Using io.WriteString to write a slice of bytes, as in
    
    
    io.WriteString(w, string(b))
    

is both unnecessary and inefficient. Converting from `[]byte` to `string` has to allocate and copy the data, and we could simply use `w.Write(b)` instead.

Available since
    2024.1

### SA9 – Dubious code constructs that have a high probability of being wrong

Checks in this category find code that is probably wrong. Unlike checks in the other `SA` categories, checks in `SA9` have a slight chance of reporting false positives. However, even false positives will point at code that is confusing and that should probably be refactored.

#### SA9001 - Defers in range loops may not run when you expect them to

Available since
    2017.1

#### SA9002 - Using a non-octal `os.FileMode` that looks like it was meant to be in octal.

Available since
    2017.1

#### SA9003 - Empty body in an if or else branch non-default

Available since
    2017.1

#### SA9004 - Only the first constant has an explicit type

In a constant declaration such as the following:
    
    
    const (
        First byte = 1
        Second     = 2
    )
    

the constant Second does not have the same type as the constant First. This construct shouldn’t be confused with
    
    
    const (
        First byte = iota
        Second
    )
    

where `First` and `Second` do indeed have the same type. The type is only passed on when no explicit value is assigned to the constant.

When declaring enumerations with explicit values it is therefore important not to write
    
    
    const (
          EnumFirst EnumType = 1
          EnumSecond         = 2
          EnumThird          = 3
    )
    

This discrepancy in types can cause various confusing behaviors and bugs.

#### Wrong type in variable declarations

The most obvious issue with such incorrect enumerations expresses itself as a compile error:
    
    
    package pkg
    
    const (
        EnumFirst  uint8 = 1
        EnumSecond       = 2
    )
    
    func fn(useFirst bool) {
        x := EnumSecond
        if useFirst {
            x = EnumFirst
        }
    }
    

fails to compile with
    
    
    ./const.go:11:5: cannot use EnumFirst (type uint8) as type int in assignment
    

#### Losing method sets

A more subtle issue occurs with types that have methods and optional interfaces. Consider the following:
    
    
    package main
    
    import "fmt"
    
    type Enum int
    
    func (e Enum) String() string {
        return "an enum"
    }
    
    const (
        EnumFirst  Enum = 1
        EnumSecond      = 2
    )
    
    func main() {
        fmt.Println(EnumFirst)
        fmt.Println(EnumSecond)
    }
    

This code will output
    
    
    an enum
    2
    

as `EnumSecond` has no explicit type, and thus defaults to `int`.

Available since
    2019.1

#### SA9005 - Trying to marshal a struct with no public fields nor custom marshaling

The `encoding/json` and `encoding/xml` packages only operate on exported fields in structs, not unexported ones. It is usually an error to try to (un)marshal structs that only consist of unexported fields.

This check will not flag calls involving types that define custom marshaling behavior, e.g. via `MarshalJSON` methods. It will also not flag empty structs.

Available since
    2019.2

#### SA9006 - Dubious bit shifting of a fixed size integer value

Bit shifting a value past its size will always clear the value.

For instance:
    
    
    v := int8(42)
    v >>= 8
    

will always result in 0.

This check flags bit shifting operations on fixed size integer values only. That is, int, uint and uintptr are never flagged to avoid potential false positives in somewhat exotic but valid bit twiddling tricks:
    
    
    // Clear any value above 32 bits if integers are more than 32 bits.
    func f(i int) int {
        v := i >> 32
        v = v << 32
        return i-v
    }
    

Available since
    2020.2

#### SA9007 - Deleting a directory that shouldn’t be deleted

It is virtually never correct to delete system directories such as /tmp or the user’s home directory. However, it can be fairly easy to do by mistake, for example by mistakenly using `os.TempDir` instead of `ioutil.TempDir`, or by forgetting to add a suffix to the result of `os.UserHomeDir`.

Writing
    
    
    d := os.TempDir()
    defer os.RemoveAll(d)
    

in your unit tests will have a devastating effect on the stability of your system.

This check flags attempts at deleting the following directories:

  * os.TempDir
  * os.UserCacheDir
  * os.UserConfigDir
  * os.UserHomeDir



Available since
    2022.1

#### SA9008 - `else` branch of a type assertion is probably not reading the right value

When declaring variables as part of an `if` statement (like in `if foo := ...; foo {`), the same variables will also be in the scope of the `else` branch. This means that in the following example
    
    
    if x, ok := x.(int); ok {
        // ...
    } else {
        fmt.Printf("unexpected type %T", x)
    }
    

`x` in the `else` branch will refer to the `x` from `x, ok :=`; it will not refer to the `x` that is being type-asserted. The result of a failed type assertion is the zero value of the type that is being asserted to, so `x` in the else branch will always have the value `0` and the type `int`.

Available since
    2022.1

#### SA9009 - Ineffectual Go compiler directive

A potential Go compiler directive was found, but is ineffectual as it begins with whitespace.

Available since
    2024.1

## S – `simple`

The S category of checks, codenamed `simple`, contains all checks that are concerned with simplifying code.

### S1 – Code simplifications

Checks in this category find code that is unnecessarily complex and that can be trivially simplified.

#### S1000 - Use plain channel send or receive instead of single-case select

Select statements with a single case can be replaced with a simple send or receive.

**Before:**
    
    
     select {
    case x := <-ch:
        fmt.Println(x)
    }

**After:**
    
    
     x := <-ch
    fmt.Println(x)

Available since
    2017.1

#### S1001 - Replace for loop with call to copy

Use `copy()` for copying elements from one slice to another. For arrays of identical size, you can use simple assignment.

**Before:**
    
    
     for i, x := range src {
        dst[i] = x
    }

**After:**
    
    
     copy(dst, src)

Available since
    2017.1

#### S1002 - Omit comparison with boolean constant

**Before:**
    
    
     if x == true {}

**After:**
    
    
     if x {}

Available since
    2017.1

#### S1003 - Replace call to `strings.Index` with `strings.Contains`

**Before:**
    
    
     if strings.Index(x, y) != -1 {}

**After:**
    
    
     if strings.Contains(x, y) {}

Available since
    2017.1

#### S1004 - Replace call to `bytes.Compare` with `bytes.Equal`

**Before:**
    
    
     if bytes.Compare(x, y) == 0 {}

**After:**
    
    
     if bytes.Equal(x, y) {}

Available since
    2017.1

#### S1005 - Drop unnecessary use of the blank identifier

In many cases, assigning to the blank identifier is unnecessary.

**Before:**
    
    
     for _ = range s {}
    x, _ = someMap[key]
    _ = <-ch

**After:**
    
    
     for range s{}
    x = someMap[key]
    <-ch

Available since
    2017.1

#### S1006 - Use `for { ... }` for infinite loops

For infinite loops, using `for { ... }` is the most idiomatic choice.

Available since
    2017.1

#### S1007 - Simplify regular expression by using raw string literal

Raw string literals use backticks instead of quotation marks and do not support any escape sequences. This means that the backslash can be used freely, without the need of escaping.

Since regular expressions have their own escape sequences, raw strings can improve their readability.

**Before:**
    
    
     regexp.Compile("\\A(\\w+) profile: total \\d+\\n\\z")

**After:**
    
    
     regexp.Compile(`\A(\w+) profile: total \d+\n\z`)

Available since
    2017.1

#### S1008 - Simplify returning boolean expression

**Before:**
    
    
     if <expr> {
        return true
    }
    return false

**After:**
    
    
     return <expr>

Available since
    2017.1

#### S1009 - Omit redundant nil check on slices, maps, and channels

The `len` function is defined for all slices, maps, and channels, even nil ones, which have a length of zero. It is not necessary to check for nil before checking that their length is not zero.

**Before:**
    
    
     if x != nil && len(x) != 0 {}

**After:**
    
    
     if len(x) != 0 {}

Available since
    2017.1

#### S1010 - Omit default slice index

When slicing, the second index defaults to the length of the value, making `s[n:len(s)]` and `s[n:]` equivalent.

Available since
    2017.1

#### S1011 - Use a single `append` to concatenate two slices

**Before:**
    
    
     for _, e := range y {
        x = append(x, e)
    }
    
    for i := range y {
        x = append(x, y[i])
    }
    
    for i := range y {
        v := y[i]
        x = append(x, v)
    }

**After:**
    
    
     x = append(x, y...)
    x = append(x, y...)
    x = append(x, y...)

Available since
    2017.1

#### S1012 - Replace `time.Now().Sub(x)` with `time.Since(x)`

The `time.Since` helper has the same effect as using `time.Now().Sub(x)` but is easier to read.

**Before:**
    
    
     time.Now().Sub(x)

**After:**
    
    
     time.Since(x)

Available since
    2017.1

#### S1016 - Use a type conversion instead of manually copying struct fields

Two struct types with identical fields can be converted between each other. In older versions of Go, the fields had to have identical struct tags. Since Go 1.8, however, struct tags are ignored during conversions. It is thus not necessary to manually copy every field individually.

**Before:**
    
    
     var x T1
    y := T2{
        Field1: x.Field1,
        Field2: x.Field2,
    }

**After:**
    
    
     var x T1
    y := T2(x)

Available since
    2017.1

#### S1017 - Replace manual trimming with `strings.TrimPrefix`

Instead of using `strings.HasPrefix` and manual slicing, use the `strings.TrimPrefix` function. If the string doesn’t start with the prefix, the original string will be returned. Using `strings.TrimPrefix` reduces complexity, and avoids common bugs, such as off-by-one mistakes.

**Before:**
    
    
     if strings.HasPrefix(str, prefix) {
        str = str[len(prefix):]
    }

**After:**
    
    
     str = strings.TrimPrefix(str, prefix)

Available since
    2017.1

#### S1018 - Use `copy` for sliding elements

`copy()` permits using the same source and destination slice, even with overlapping ranges. This makes it ideal for sliding elements in a slice.

**Before:**
    
    
     for i := 0; i < n; i++ {
        bs[i] = bs[offset+i]
    }

**After:**
    
    
     copy(bs[:n], bs[offset:])

Available since
    2017.1

#### S1019 - Simplify `make` call by omitting redundant arguments

The `make` function has default values for the length and capacity arguments. For channels, the length defaults to zero, and for slices, the capacity defaults to the length.

Available since
    2017.1

#### S1020 - Omit redundant nil check in type assertion

**Before:**
    
    
     if _, ok := i.(T); ok && i != nil {}

**After:**
    
    
     if _, ok := i.(T); ok {}

Available since
    2017.1

#### S1021 - Merge variable declaration and assignment

**Before:**
    
    
     var x uint
    x = 1

**After:**
    
    
     var x uint = 1

Available since
    2017.1

#### S1023 - Omit redundant control flow

Functions that have no return value do not need a return statement as the final statement of the function.

Switches in Go do not have automatic fallthrough, unlike languages like C. It is not necessary to have a break statement as the final statement in a case block.

Available since
    2017.1

#### S1024 - Replace `x.Sub(time.Now())` with `time.Until(x)`

The `time.Until` helper has the same effect as using `x.Sub(time.Now())` but is easier to read.

**Before:**
    
    
     x.Sub(time.Now())

**After:**
    
    
     time.Until(x)

Available since
    2017.1

#### S1025 - Don’t use `fmt.Sprintf("%s", x)` unnecessarily

In many instances, there are easier and more efficient ways of getting a value’s string representation. Whenever a value’s underlying type is a string already, or the type has a String method, they should be used directly.

Given the following shared definitions
    
    
    type T1 string
    type T2 int
    
    func (T2) String() string { return "Hello, world" }
    
    var x string
    var y T1
    var z T2
    

we can simplify
    
    
    fmt.Sprintf("%s", x)
    fmt.Sprintf("%s", y)
    fmt.Sprintf("%s", z)
    

to
    
    
    x
    string(y)
    z.String()
    

Available since
    2017.1

#### S1028 - Simplify error construction with `fmt.Errorf`

**Before:**
    
    
     errors.New(fmt.Sprintf(...))

**After:**
    
    
     fmt.Errorf(...)

Available since
    2017.1

#### S1029 - Range over the string directly

Ranging over a string will yield byte offsets and runes. If the offset isn’t used, this is functionally equivalent to converting the string to a slice of runes and ranging over that. Ranging directly over the string will be more performant, however, as it avoids allocating a new slice, the size of which depends on the length of the string.

**Before:**
    
    
     for _, r := range []rune(s) {}

**After:**
    
    
     for _, r := range s {}

Available since
    2017.1

#### S1030 - Use `bytes.Buffer.String` or `bytes.Buffer.Bytes`

`bytes.Buffer` has both a `String` and a `Bytes` method. It is almost never necessary to use `string(buf.Bytes())` or `[]byte(buf.String())` – simply use the other method.

The only exception to this are map lookups. Due to a compiler optimization, `m[string(buf.Bytes())]` is more efficient than `m[buf.String()]`.

Available since
    2017.1

#### S1031 - Omit redundant nil check around loop

You can use range on nil slices and maps, the loop will simply never execute. This makes an additional nil check around the loop unnecessary.

**Before:**
    
    
     if s != nil {
        for _, x := range s {
            ...
        }
    }

**After:**
    
    
     for _, x := range s {
        ...
    }

Available since
    2017.1

#### S1032 - Use `sort.Ints(x)`, `sort.Float64s(x)`, and `sort.Strings(x)`

The `sort.Ints`, `sort.Float64s` and `sort.Strings` functions are easier to read than `sort.Sort(sort.IntSlice(x))`, `sort.Sort(sort.Float64Slice(x))` and `sort.Sort(sort.StringSlice(x))`.

**Before:**
    
    
     sort.Sort(sort.StringSlice(x))

**After:**
    
    
     sort.Strings(x)

Available since
    2019.1

#### S1033 - Unnecessary guard around call to `delete`

Calling `delete` on a nil map is a no-op.

Available since
    2019.2

#### S1034 - Use result of type assertion to simplify cases

Available since
    2019.2

#### S1035 - Redundant call to `net/http.CanonicalHeaderKey` in method call on `net/http.Header`

The methods on `net/http.Header`, namely `Add`, `Del`, `Get` and `Set`, already canonicalize the given header name.

Available since
    2020.1

#### S1036 - Unnecessary guard around map access

When accessing a map key that doesn’t exist yet, one receives a zero value. Often, the zero value is a suitable value, for example when using append or doing integer math.

The following
    
    
    if _, ok := m["foo"]; ok {
        m["foo"] = append(m["foo"], "bar")
    } else {
        m["foo"] = []string{"bar"}
    }
    

can be simplified to
    
    
    m["foo"] = append(m["foo"], "bar")
    

and
    
    
    if _, ok := m2["k"]; ok {
        m2["k"] += 4
    } else {
        m2["k"] = 4
    }
    

can be simplified to
    
    
    m["k"] += 4
    

Available since
    2020.1

#### S1037 - Elaborate way of sleeping

Using a select statement with a single case receiving from the result of `time.After` is a very elaborate way of sleeping that can much simpler be expressed with a simple call to time.Sleep.

Available since
    2020.1

#### S1038 - Unnecessarily complex way of printing formatted string

Instead of using `fmt.Print(fmt.Sprintf(...))`, one can use `fmt.Printf(...)`.

Available since
    2020.1

#### S1039 - Unnecessary use of `fmt.Sprint`

Calling `fmt.Sprint` with a single string argument is unnecessary and identical to using the string directly.

Available since
    2020.1

#### S1040 - Type assertion to current type

The type assertion `x.(SomeInterface)`, when `x` already has type `SomeInterface`, can only fail if `x` is nil. Usually, this is left-over code from when `x` had a different type and you can safely delete the type assertion. If you want to check that `x` is not nil, consider being explicit and using an actual `if x == nil` comparison instead of relying on the type assertion panicking.

Available since
    2021.1

## ST – `stylecheck`

The ST category of checks, codenamed `stylecheck`, contains all checks that are concerned with stylistic issues.

### ST1 – Stylistic issues

The rules contained in this category are primarily derived from the [Go wiki](https://go.dev/wiki/CodeReviewComments) and represent community consensus.

Some checks are very pedantic and disabled by default. You may want to [tweak which checks from this category run](/docs/configuration/options/#checks), based on your project's needs.

#### ST1000 - Incorrect or missing package comment non-default

Packages must have a package comment that is formatted according to the guidelines laid out in <https://go.dev/wiki/CodeReviewComments#package-comments>.

Available since
    2019.1

#### ST1001 - Dot imports are discouraged

Dot imports that aren’t in external test packages are discouraged.

The `dot_import_whitelist` option can be used to whitelist certain imports.

Quoting Go Code Review Comments:

> The `import .` form can be useful in tests that, due to circular dependencies, cannot be made part of the package being tested:
>     
>     
>     package foo_test
>     
>     import (
>         "bar/testutil" // also imports "foo"
>         . "foo"
>     )
>     
> 
> In this case, the test file cannot be in package foo because it uses `bar/testutil`, which imports `foo`. So we use the `import .` form to let the file pretend to be part of package foo even though it is not. Except for this one case, do not use `import .` in your programs. It makes the programs much harder to read because it is unclear whether a name like `Quux` is a top-level identifier in the current package or in an imported package.

Available since
    2019.1
Options
    

  * [dot_import_whitelist](/docs/configuration/options/#dot_import_whitelist)



#### ST1003 - Poorly chosen identifier non-default

Identifiers, such as variable and package names, follow certain rules.

See the following links for details:

  * <https://go.dev/doc/effective_go#package-names>
  * <https://go.dev/doc/effective_go#mixed-caps>
  * <https://go.dev/wiki/CodeReviewComments#initialisms>
  * <https://go.dev/wiki/CodeReviewComments#variable-names>



Available since
    2019.1
Options
    

  * [initialisms](/docs/configuration/options/#initialisms)



#### ST1005 - Incorrectly formatted error string

Error strings follow a set of guidelines to ensure uniformity and good composability.

Quoting Go Code Review Comments:

> Error strings should not be capitalized (unless beginning with proper nouns or acronyms) or end with punctuation, since they are usually printed following other context. That is, use `fmt.Errorf("something bad")` not `fmt.Errorf("Something bad")`, so that `log.Printf("Reading %s: %v", filename, err)` formats without a spurious capital letter mid-message.

Available since
    2019.1

#### ST1006 - Poorly chosen receiver name

Quoting Go Code Review Comments:

> The name of a method’s receiver should be a reflection of its identity; often a one or two letter abbreviation of its type suffices (such as “c” or “cl” for “Client”). Don’t use generic names such as “me”, “this” or “self”, identifiers typical of object-oriented languages that place more emphasis on methods as opposed to functions. The name need not be as descriptive as that of a method argument, as its role is obvious and serves no documentary purpose. It can be very short as it will appear on almost every line of every method of the type; familiarity admits brevity. Be consistent, too: if you call the receiver “c” in one method, don’t call it “cl” in another.

Available since
    2019.1

#### ST1008 - A function’s error value should be its last return value

A function’s error value should be its last return value.

Available since
    2019.1

#### ST1011 - Poorly chosen name for variable of type `time.Duration`

`time.Duration` values represent an amount of time, which is represented as a count of nanoseconds. An expression like `5 * time.Microsecond` yields the value `5000`. It is therefore not appropriate to suffix a variable of type `time.Duration` with any time unit, such as `Msec` or `Milli`.

Available since
    2019.1

#### ST1012 - Poorly chosen name for error variable

Error variables that are part of an API should be called `errFoo` or `ErrFoo`.

Available since
    2019.1

#### ST1013 - Should use constants for HTTP error codes, not magic numbers

HTTP has a tremendous number of status codes. While some of those are well known (200, 400, 404, 500), most of them are not. The `net/http` package provides constants for all status codes that are part of the various specifications. It is recommended to use these constants instead of hard-coding magic numbers, to vastly improve the readability of your code.

Available since
    2019.1
Options
    

  * [http_status_code_whitelist](/docs/configuration/options/#http_status_code_whitelist)



#### ST1015 - A switch’s default case should be the first or last case

Available since
    2019.1

#### ST1016 - Use consistent method receiver names non-default

Available since
    2019.1

#### ST1017 - Don’t use Yoda conditions

Yoda conditions are conditions of the kind `if 42 == x`, where the literal is on the left side of the comparison. These are a common idiom in languages in which assignment is an expression, to avoid bugs of the kind `if (x = 42)`. In Go, which doesn’t allow for this kind of bug, we prefer the more idiomatic `if x == 42`.

Available since
    2019.2

#### ST1018 - Avoid zero-width and control characters in string literals

Available since
    2019.2

#### ST1019 - Importing the same package multiple times

Go allows importing the same package multiple times, as long as different import aliases are being used. That is, the following bit of code is valid:
    
    
    import (
        "fmt"
        fumpt "fmt"
        format "fmt"
        _ "fmt"
    )
    

However, this is very rarely done on purpose. Usually, it is a sign of code that got refactored, accidentally adding duplicate import statements. It is also a rarely known feature, which may contribute to confusion.

Do note that sometimes, this feature may be used intentionally (see for example <https://github.com/golang/go/commit/3409ce39bfd7584523b7a8c150a310cea92d879d>) – if you want to allow this pattern in your code base, you’re advised to disable this check.

Available since
    2020.1

#### ST1020 - The documentation of an exported function should start with the function’s name non-default

Doc comments work best as complete sentences, which allow a wide variety of automated presentations. The first sentence should be a one-sentence summary that starts with the name being declared.

If every doc comment begins with the name of the item it describes, you can use the `doc` subcommand of the `go` tool and run the output through grep.

See <https://go.dev/doc/effective_go#commentary> for more information on how to write good documentation.

Available since
    2020.1

#### ST1021 - The documentation of an exported type should start with type’s name non-default

Doc comments work best as complete sentences, which allow a wide variety of automated presentations. The first sentence should be a one-sentence summary that starts with the name being declared.

If every doc comment begins with the name of the item it describes, you can use the `doc` subcommand of the `go` tool and run the output through grep.

See <https://go.dev/doc/effective_go#commentary> for more information on how to write good documentation.

Available since
    2020.1

#### ST1022 - The documentation of an exported variable or constant should start with variable’s name non-default

Doc comments work best as complete sentences, which allow a wide variety of automated presentations. The first sentence should be a one-sentence summary that starts with the name being declared.

If every doc comment begins with the name of the item it describes, you can use the `doc` subcommand of the `go` tool and run the output through grep.

See <https://go.dev/doc/effective_go#commentary> for more information on how to write good documentation.

Available since
    2020.1

#### ST1023 - Redundant type in variable declaration non-default

Available since
    2021.1

## QF – `quickfix`

The QF category of checks, codenamed `quickfix`, contains checks that are used as part of _gopls_ for automatic refactorings. In the context of gopls, diagnostics of these checks will usually show up as hints, sometimes as information-level diagnostics.

### QF1 – Quickfixes

#### QF1001 - Apply De Morgan’s law

Available since
    2021.1

#### QF1002 - Convert untagged switch to tagged switch

An untagged switch that compares a single variable against a series of values can be replaced with a tagged switch.

**Before:**
    
    
     switch {
    case x == 1 || x == 2, x == 3:
        ...
    case x == 4:
        ...
    default:
        ...
    }

**After:**
    
    
     switch x {
    case 1, 2, 3:
        ...
    case 4:
        ...
    default:
        ...
    }

Available since
    2021.1

#### QF1003 - Convert if/else-if chain to tagged switch

A series of if/else-if checks comparing the same variable against values can be replaced with a tagged switch.

**Before:**
    
    
     if x == 1 || x == 2 {
        ...
    } else if x == 3 {
        ...
    } else {
        ...
    }

**After:**
    
    
     switch x {
    case 1, 2:
        ...
    case 3:
        ...
    default:
        ...
    }

Available since
    2021.1

#### QF1004 - Use `strings.ReplaceAll` instead of `strings.Replace` with `n == -1`

Available since
    2021.1

#### QF1005 - Expand call to `math.Pow`

Some uses of `math.Pow` can be simplified to basic multiplication.

**Before:**
    
    
     math.Pow(x, 2)

**After:**
    
    
     x * x

Available since
    2021.1

#### QF1006 - Lift `if`+`break` into loop condition

**Before:**
    
    
     for {
        if done {
            break
        }
        ...
    }

**After:**
    
    
     for !done {
        ...
    }

Available since
    2021.1

#### QF1007 - Merge conditional assignment into variable declaration

**Before:**
    
    
     x := false
    if someCondition {
        x = true
    }

**After:**
    
    
     x := someCondition

Available since
    2021.1

#### QF1008 - Omit embedded fields from selector expression

Available since
    2021.1

#### QF1009 - Use `time.Time.Equal` instead of `==` operator

Available since
    2021.1

#### QF1010 - Convert slice of bytes to string when printing it

Available since
    2021.1

#### QF1011 - Omit redundant type from variable declaration

Available since
    2021.1

#### QF1012 - Use `fmt.Fprintf(x, ...)` instead of `x.Write(fmt.Sprintf(...))`

Available since
    2022.1
  *[v]: View this template
  *[t]: Discuss this template
  *[e]: Edit this template
