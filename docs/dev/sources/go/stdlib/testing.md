# Go testing

> Source: https://pkg.go.dev/testing
> Fetched: 2026-01-31T10:55:16.086761+00:00
> Content-Hash: 96c39f370b708e52
> Type: html

---

### Overview ¶

  * Benchmarks
  * b.N-style benchmarks
  * Examples
  * Fuzzing
  * Skipping
  * Subtests and Sub-benchmarks
  * Main



Package testing provides support for automated testing of Go packages. It is intended to be used in concert with the "go test" command, which automates execution of any function of the form 
    
    
    func TestXxx(*testing.T)
    

where Xxx does not start with a lowercase letter. The function name serves to identify the test routine. 

Within these functions, use T.Error, T.Fail or related methods to signal failure. 

To write a new test suite, create a file that contains the TestXxx functions as described here, and give that file a name ending in "_test.go". The file will be excluded from regular package builds but will be included when the "go test" command is run. 

The test file can be in the same package as the one being tested, or in a corresponding package with the suffix "_test". 

If the test file is in the same package, it may refer to unexported identifiers within the package, as in this example: 
    
    
    package abs
    
    import "testing"
    
    func TestAbs(t *testing.T) {
        got := Abs(-1)
        if got != 1 {
            t.Errorf("Abs(-1) = %d; want 1", got)
        }
    }
    

If the file is in a separate "_test" package, the package being tested must be imported explicitly and only its exported identifiers may be used. This is known as "black box" testing. 
    
    
    package abs_test
    
    import (
    	"testing"
    
    	"path_to_pkg/abs"
    )
    
    func TestAbs(t *testing.T) {
        got := abs.Abs(-1)
        if got != 1 {
            t.Errorf("Abs(-1) = %d; want 1", got)
        }
    }
    

For more detail, run [go help test](https://pkg.go.dev/cmd/go#hdr-Test_packages) and [go help testflag](https://pkg.go.dev/cmd/go#hdr-Testing_flags). 

#### Benchmarks ¶

Functions of the form 
    
    
    func BenchmarkXxx(*testing.B)
    

are considered benchmarks, and are executed by the "go test" command when its -bench flag is provided. Benchmarks are run sequentially. 

For a description of the testing flags, see [go help testflag](https://pkg.go.dev/cmd/go#hdr-Testing_flags). 

A sample benchmark function looks like this: 
    
    
    func BenchmarkRandInt(b *testing.B) {
        for b.Loop() {
            rand.Int()
        }
    }
    

The output 
    
    
    BenchmarkRandInt-8   	68453040	        17.8 ns/op
    

means that the body of the loop ran 68453040 times at a speed of 17.8 ns per loop. 

Only the body of the loop is timed, so benchmarks may do expensive setup before calling b.Loop, which will not be counted toward the benchmark measurement: 
    
    
    func BenchmarkBigLen(b *testing.B) {
        big := NewBig()
        for b.Loop() {
            big.Len()
        }
    }
    

If a benchmark needs to test performance in a parallel setting, it may use the RunParallel helper function; such benchmarks are intended to be used with the go test -cpu flag: 
    
    
    func BenchmarkTemplateParallel(b *testing.B) {
        templ := template.Must(template.New("test").Parse("Hello, {{.}}!"))
        b.RunParallel(func(pb *testing.PB) {
            var buf bytes.Buffer
            for pb.Next() {
                buf.Reset()
                templ.Execute(&buf, "World")
            }
        })
    }
    

A detailed specification of the benchmark results format is given in <https://go.dev/design/14313-benchmark-format>. 

There are standard tools for working with benchmark results at [golang.org/x/perf/cmd](/golang.org/x/perf/cmd). In particular, [golang.org/x/perf/cmd/benchstat](/golang.org/x/perf/cmd/benchstat) performs statistically robust A/B comparisons. 

#### b.N-style benchmarks ¶

Prior to the introduction of B.Loop, benchmarks were written in a different style using B.N. For example: 
    
    
    func BenchmarkRandInt(b *testing.B) {
        for range b.N {
            rand.Int()
        }
    }
    

In this style of benchmark, the benchmark function must run the target code b.N times. The benchmark function is called multiple times with b.N adjusted until the benchmark function lasts long enough to be timed reliably. This also means any setup done before the loop may be run several times. 

If a benchmark needs some expensive setup before running, the timer should be explicitly reset: 
    
    
    func BenchmarkBigLen(b *testing.B) {
        big := NewBig()
        b.ResetTimer()
        for range b.N {
            big.Len()
        }
    }
    

New benchmarks should prefer using B.Loop, which is more robust and more efficient. 

#### Examples ¶

The package also runs and verifies example code. Example functions may include a concluding line comment that begins with "Output:" and is compared with the standard output of the function when the tests are run. (The comparison ignores leading and trailing space.) These are examples of an example: 
    
    
    func ExampleHello() {
        fmt.Println("hello")
        // Output: hello
    }
    
    func ExampleSalutations() {
        fmt.Println("hello, and")
        fmt.Println("goodbye")
        // Output:
        // hello, and
        // goodbye
    }
    

The comment prefix "Unordered output:" is like "Output:", but matches any line order: 
    
    
    func ExamplePerm() {
        for _, value := range Perm(5) {
            fmt.Println(value)
        }
        // Unordered output: 4
        // 2
        // 1
        // 3
        // 0
    }
    

Example functions without output comments are compiled but not executed. 

The naming convention to declare examples for the package, a function F, a type T and method M on type T are: 
    
    
    func Example() { ... }
    func ExampleF() { ... }
    func ExampleT() { ... }
    func ExampleT_M() { ... }
    

Multiple example functions for a package/type/function/method may be provided by appending a distinct suffix to the name. The suffix must start with a lower-case letter. 
    
    
    func Example_suffix() { ... }
    func ExampleF_suffix() { ... }
    func ExampleT_suffix() { ... }
    func ExampleT_M_suffix() { ... }
    

The entire test file is presented as the example when it contains a single example function, at least one other function, type, variable, or constant declaration, and no test or benchmark functions. 

#### Fuzzing ¶

'go test' and the testing package support fuzzing, a testing technique where a function is called with randomly generated inputs to find bugs not anticipated by unit tests. 

Functions of the form 
    
    
    func FuzzXxx(*testing.F)
    

are considered fuzz tests. 

For example: 
    
    
    func FuzzHex(f *testing.F) {
      for _, seed := range [][]byte{{}, {0}, {9}, {0xa}, {0xf}, {1, 2, 3, 4}} {
        f.Add(seed)
      }
      f.Fuzz(func(t *testing.T, in []byte) {
        enc := hex.EncodeToString(in)
        out, err := hex.DecodeString(enc)
        if err != nil {
          t.Fatalf("%v: decode: %v", in, err)
        }
        if !bytes.Equal(in, out) {
          t.Fatalf("%v: not equal after round trip: %v", in, out)
        }
      })
    }
    

A fuzz test maintains a seed corpus, or a set of inputs which are run by default, and can seed input generation. Seed inputs may be registered by calling F.Add or by storing files in the directory testdata/fuzz/<Name> (where <Name> is the name of the fuzz test) within the package containing the fuzz test. Seed inputs are optional, but the fuzzing engine may find bugs more efficiently when provided with a set of small seed inputs with good code coverage. These seed inputs can also serve as regression tests for bugs identified through fuzzing. 

The function passed to F.Fuzz within the fuzz test is considered the fuzz target. A fuzz target must accept a *T parameter, followed by one or more parameters for random inputs. The types of arguments passed to F.Add must be identical to the types of these parameters. The fuzz target may signal that it's found a problem the same way tests do: by calling T.Fail (or any method that calls it like T.Error or T.Fatal) or by panicking. 

When fuzzing is enabled (by setting the -fuzz flag to a regular expression that matches a specific fuzz test), the fuzz target is called with arguments generated by repeatedly making random changes to the seed inputs. On supported platforms, 'go test' compiles the test executable with fuzzing coverage instrumentation. The fuzzing engine uses that instrumentation to find and cache inputs that expand coverage, increasing the likelihood of finding bugs. If the fuzz target fails for a given input, the fuzzing engine writes the inputs that caused the failure to a file in the directory testdata/fuzz/<Name> within the package directory. This file later serves as a seed input. If the file can't be written at that location (for example, because the directory is read-only), the fuzzing engine writes the file to the fuzz cache directory within the build cache instead. 

When fuzzing is disabled, the fuzz target is called with the seed inputs registered with F.Add and seed inputs from testdata/fuzz/<Name>. In this mode, the fuzz test acts much like a regular test, with subtests started with F.Fuzz instead of T.Run. 

See <https://go.dev/doc/fuzz> for documentation about fuzzing. 

#### Skipping ¶

Tests or benchmarks may be skipped at run time with a call to T.Skip or B.Skip: 
    
    
    func TestTimeConsuming(t *testing.T) {
        if testing.Short() {
            t.Skip("skipping test in short mode.")
        }
        ...
    }
    

The T.Skip method can be used in a fuzz target if the input is invalid, but should not be considered a failing input. For example: 
    
    
    func FuzzJSONMarshaling(f *testing.F) {
        f.Fuzz(func(t *testing.T, b []byte) {
            var v interface{}
            if err := json.Unmarshal(b, &v); err != nil {
                t.Skip()
            }
            if _, err := json.Marshal(v); err != nil {
                t.Errorf("Marshal: %v", err)
            }
        })
    }
    

#### Subtests and Sub-benchmarks ¶

The T.Run and B.Run methods allow defining subtests and sub-benchmarks, without having to define separate functions for each. This enables uses like table-driven benchmarks and creating hierarchical tests. It also provides a way to share common setup and tear-down code: 
    
    
    func TestFoo(t *testing.T) {
        // <setup code>
        t.Run("A=1", func(t *testing.T) { ... })
        t.Run("A=2", func(t *testing.T) { ... })
        t.Run("B=1", func(t *testing.T) { ... })
        // <tear-down code>
    }
    

Each subtest and sub-benchmark has a unique name: the combination of the name of the top-level test and the sequence of names passed to Run, separated by slashes, with an optional trailing sequence number for disambiguation. 

The argument to the -run, -bench, and -fuzz command-line flags is an unanchored regular expression that matches the test's name. For tests with multiple slash-separated elements, such as subtests, the argument is itself slash-separated, with expressions matching each name element in turn. Because it is unanchored, an empty expression matches any string. For example, using "matching" to mean "whose name contains": 
    
    
    go test -run ''        # Run all tests.
    go test -run Foo       # Run top-level tests matching "Foo", such as "TestFooBar".
    go test -run Foo/A=    # For top-level tests matching "Foo", run subtests matching "A=".
    go test -run /A=1      # For all top-level tests, run subtests matching "A=1".
    go test -fuzz FuzzFoo  # Fuzz the target matching "FuzzFoo"
    

The -run argument can also be used to run a specific value in the seed corpus, for debugging. For example: 
    
    
    go test -run=FuzzFoo/9ddb952d9814
    

The -fuzz and -run flags can both be set, in order to fuzz a target but skip the execution of all other tests. 

Subtests can also be used to control parallelism. A parent test will only complete once all of its subtests complete. In this example, all tests are run in parallel with each other, and only with each other, regardless of other top-level tests that may be defined: 
    
    
    func TestGroupedParallel(t *testing.T) {
        for _, tc := range tests {
            t.Run(tc.Name, func(t *testing.T) {
                t.Parallel()
                ...
            })
        }
    }
    

Run does not return until parallel subtests have completed, providing a way to clean up after a group of parallel tests: 
    
    
    func TestTeardownParallel(t *testing.T) {
        // This Run will not return until the parallel tests finish.
        t.Run("group", func(t *testing.T) {
            t.Run("Test1", parallelTest1)
            t.Run("Test2", parallelTest2)
            t.Run("Test3", parallelTest3)
        })
        // <tear-down code>
    }
    

#### Main ¶

It is sometimes necessary for a test or benchmark program to do extra setup or teardown before or after it executes. It is also sometimes necessary to control which code runs on the main thread. To support these and other cases, if a test file contains a function: 
    
    
    func TestMain(m *testing.M)
    

then the generated test will call TestMain(m) instead of running the tests or benchmarks directly. TestMain runs in the main goroutine and can do whatever setup and teardown is necessary around a call to m.Run. m.Run will return an exit code that may be passed to [os.Exit](/os#Exit). If TestMain returns, the test wrapper will pass the result of m.Run to [os.Exit](/os#Exit) itself. 

When TestMain is called, flag.Parse has not been run. If TestMain depends on command-line flags, including those of the testing package, it should call [flag.Parse](/flag#Parse) explicitly. Command line flags are always parsed by the time test or benchmark functions run. 

A simple implementation of TestMain is: 
    
    
    func TestMain(m *testing.M) {
    	// call flag.Parse() here if TestMain uses flags
    	m.Run()
    }
    

TestMain is a low-level primitive and should not be necessary for casual testing needs, where ordinary test functions suffice. 

### Index ¶

  * func AllocsPerRun(runs int, f func()) (avg float64)
  * func CoverMode() string
  * func Coverage() float64
  * func Init()
  * func Main(matchString func(pat, str string) (bool, error), tests []InternalTest, ...)
  * func RegisterCover(c Cover)
  * func RunBenchmarks(matchString func(pat, str string) (bool, error), ...)
  * func RunExamples(matchString func(pat, str string) (bool, error), examples []InternalExample) (ok bool)
  * func RunTests(matchString func(pat, str string) (bool, error), tests []InternalTest) (ok bool)
  * func Short() bool
  * func Testing() bool
  * func Verbose() bool
  * type B
  *     * func (c *B) Attr(key, value string)
    * func (c *B) Chdir(dir string)
    * func (c *B) Cleanup(f func())
    * func (c *B) Context() context.Context
    * func (b *B) Elapsed() time.Duration
    * func (c *B) Error(args ...any)
    * func (c *B) Errorf(format string, args ...any)
    * func (c *B) Fail()
    * func (c *B) FailNow()
    * func (c *B) Failed() bool
    * func (c *B) Fatal(args ...any)
    * func (c *B) Fatalf(format string, args ...any)
    * func (c *B) Helper()
    * func (c *B) Log(args ...any)
    * func (c *B) Logf(format string, args ...any)
    * func (b *B) Loop() bool
    * func (c *B) Name() string
    * func (c *B) Output() io.Writer
    * func (b *B) ReportAllocs()
    * func (b *B) ReportMetric(n float64, unit string)
    * func (b *B) ResetTimer()
    * func (b *B) Run(name string, f func(b *B)) bool
    * func (b *B) RunParallel(body func(*PB))
    * func (b *B) SetBytes(n int64)
    * func (b *B) SetParallelism(p int)
    * func (c *B) Setenv(key, value string)
    * func (c *B) Skip(args ...any)
    * func (c *B) SkipNow()
    * func (c *B) Skipf(format string, args ...any)
    * func (c *B) Skipped() bool
    * func (b *B) StartTimer()
    * func (b *B) StopTimer()
    * func (c *B) TempDir() string
  * type BenchmarkResult
  *     * func Benchmark(f func(b *B)) BenchmarkResult
  *     * func (r BenchmarkResult) AllocedBytesPerOp() int64
    * func (r BenchmarkResult) AllocsPerOp() int64
    * func (r BenchmarkResult) MemString() string
    * func (r BenchmarkResult) NsPerOp() int64
    * func (r BenchmarkResult) String() string
  * type Cover
  * type CoverBlock
  * type F
  *     * func (f *F) Add(args ...any)
    * func (c *F) Attr(key, value string)
    * func (c *F) Chdir(dir string)
    * func (c *F) Cleanup(f func())
    * func (c *F) Context() context.Context
    * func (c *F) Error(args ...any)
    * func (c *F) Errorf(format string, args ...any)
    * func (f *F) Fail()
    * func (c *F) FailNow()
    * func (c *F) Failed() bool
    * func (c *F) Fatal(args ...any)
    * func (c *F) Fatalf(format string, args ...any)
    * func (f *F) Fuzz(ff any)
    * func (f *F) Helper()
    * func (c *F) Log(args ...any)
    * func (c *F) Logf(format string, args ...any)
    * func (c *F) Name() string
    * func (c *F) Output() io.Writer
    * func (c *F) Setenv(key, value string)
    * func (c *F) Skip(args ...any)
    * func (c *F) SkipNow()
    * func (c *F) Skipf(format string, args ...any)
    * func (f *F) Skipped() bool
    * func (c *F) TempDir() string
  * type InternalBenchmark
  * type InternalExample
  * type InternalFuzzTarget
  * type InternalTest
  * type M
  *     * func MainStart(deps testDeps, tests []InternalTest, benchmarks []InternalBenchmark, ...) *M
  *     * func (m *M) Run() (code int)
  * type PB
  *     * func (pb *PB) Next() bool
  * type T
  *     * func (c *T) Attr(key, value string)
    * func (t *T) Chdir(dir string)
    * func (c *T) Cleanup(f func())
    * func (c *T) Context() context.Context
    * func (t *T) Deadline() (deadline time.Time, ok bool)
    * func (c *T) Error(args ...any)
    * func (c *T) Errorf(format string, args ...any)
    * func (c *T) Fail()
    * func (c *T) FailNow()
    * func (c *T) Failed() bool
    * func (c *T) Fatal(args ...any)
    * func (c *T) Fatalf(format string, args ...any)
    * func (c *T) Helper()
    * func (c *T) Log(args ...any)
    * func (c *T) Logf(format string, args ...any)
    * func (c *T) Name() string
    * func (c *T) Output() io.Writer
    * func (t *T) Parallel()
    * func (t *T) Run(name string, f func(t *T)) bool
    * func (t *T) Setenv(key, value string)
    * func (c *T) Skip(args ...any)
    * func (c *T) SkipNow()
    * func (c *T) Skipf(format string, args ...any)
    * func (c *T) Skipped() bool
    * func (c *T) TempDir() string
  * type TB



### Examples ¶

  * B.Loop
  * B.ReportMetric
  * B.ReportMetric (Parallel)
  * B.RunParallel



### Constants ¶

This section is empty.

### Variables ¶

This section is empty.

### Functions ¶

####  func [AllocsPerRun](https://cs.opensource.google/go/go/+/go1.25.6:src/testing/allocs.go;l=20) ¶ added in go1.1
    
    
    func AllocsPerRun(runs [int](/builtin#int), f func()) (avg [float64](/builtin#float64))

AllocsPerRun returns the average number of allocations during calls to f. Although the return value has type float64, it will always be an integral value. 

To compute the number of allocations, the function will first be run once as a warm-up. The average number of allocations over the specified number of runs will then be measured and returned. 

AllocsPerRun sets [runtime.GOMAXPROCS](/runtime#GOMAXPROCS) to 1 during its measurement and will restore it before returning. 

####  func [CoverMode](https://cs.opensource.google/go/go/+/go1.25.6:src/testing/testing.go;l=710) ¶ added in go1.8
    
    
    func CoverMode() [string](/builtin#string)

CoverMode reports what the test coverage mode is set to. The values are "set", "count", or "atomic". The return value will be empty if test coverage is not enabled. 

####  func [Coverage](https://cs.opensource.google/go/go/+/go1.25.6:src/testing/newcover.go;l=54) ¶ added in go1.4
    
    
    func Coverage() [float64](/builtin#float64)

Coverage reports the current code coverage as a fraction in the range [0, 1]. If coverage is not enabled, Coverage returns 0. 

When running a large set of sequential test cases, checking Coverage after each one can be useful for identifying which test cases exercise new code paths. It is not a replacement for the reports generated by 'go test -cover' and 'go tool cover'. 

####  func [Init](https://cs.opensource.google/go/go/+/go1.25.6:src/testing/testing.go;l=439) ¶ added in go1.13
    
    
    func Init()

Init registers testing flags. These flags are automatically registered by the "go test" command before running test functions, so Init is only needed when calling functions such as Benchmark without using "go test". 

Init is not safe to call concurrently. It has no effect if it was already called. 

####  func [Main](https://cs.opensource.google/go/go/+/go1.25.6:src/testing/testing.go;l=2166) ¶
    
    
    func Main(matchString func(pat, str [string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error)), tests []InternalTest, benchmarks []InternalBenchmark, examples []InternalExample)

Main is an internal function, part of the implementation of the "go test" command. It was exported because it is cross-package and predates "internal" packages. It is no longer used by "go test" but preserved, as much as possible, for other systems that simulate "go test" using Main, but Main sometimes cannot be updated as new functionality is added to the testing package. Systems simulating "go test" should be updated to use MainStart. 

####  func [RegisterCover](https://cs.opensource.google/go/go/+/go1.25.6:src/testing/cover.go;l=36) ¶ added in go1.2
    
    
    func RegisterCover(c Cover)

RegisterCover records the coverage data accumulators for the tests. NOTE: This function is internal to the testing infrastructure and may change. It is not covered (yet) by the Go 1 compatibility guidelines. 

####  func [RunBenchmarks](https://cs.opensource.google/go/go/+/go1.25.6:src/testing/benchmark.go;l=682) ¶
    
    
    func RunBenchmarks(matchString func(pat, str [string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error)), benchmarks []InternalBenchmark)

RunBenchmarks is an internal function but exported because it is cross-package; it is part of the implementation of the "go test" command. 

####  func [RunExamples](https://cs.opensource.google/go/go/+/go1.25.6:src/testing/example.go;l=24) ¶
    
    
    func RunExamples(matchString func(pat, str [string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error)), examples []InternalExample) (ok [bool](/builtin#bool))

RunExamples is an internal function but exported because it is cross-package; it is part of the implementation of the "go test" command. 

####  func [RunTests](https://cs.opensource.google/go/go/+/go1.25.6:src/testing/testing.go;l=2433) ¶
    
    
    func RunTests(matchString func(pat, str [string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error)), tests []InternalTest) (ok [bool](/builtin#bool))

RunTests is an internal function but exported because it is cross-package; it is part of the implementation of the "go test" command. 

####  func [Short](https://cs.opensource.google/go/go/+/go1.25.6:src/testing/testing.go;l=679) ¶
    
    
    func Short() [bool](/builtin#bool)

Short reports whether the -test.short flag is set. 

####  func [Testing](https://cs.opensource.google/go/go/+/go1.25.6:src/testing/testing.go;l=703) ¶ added in go1.21.0
    
    
    func Testing() [bool](/builtin#bool)

Testing reports whether the current code is being run in a test. This will report true in programs created by "go test", false in programs created by "go build". 

####  func [Verbose](https://cs.opensource.google/go/go/+/go1.25.6:src/testing/testing.go;l=715) ¶ added in go1.1
    
    
    func Verbose() [bool](/builtin#bool)

Verbose reports whether the -test.v flag is set. 

### Types ¶

####  type [B](https://cs.opensource.google/go/go/+/go1.25.6:src/testing/benchmark.go;l=94) ¶
    
    
    type B struct {
    	N [int](/builtin#int)
    	// contains filtered or unexported fields
    }

B is a type passed to Benchmark functions to manage benchmark timing and control the number of iterations. 

A benchmark ends when its Benchmark function returns or calls any of the methods B.FailNow, B.Fatal, B.Fatalf, B.SkipNow, B.Skip, or B.Skipf. Those methods must be called only from the goroutine running the Benchmark function. The other reporting methods, such as the variations of B.Log and B.Error, may be called simultaneously from multiple goroutines. 

Like in tests, benchmark logs are accumulated during execution and dumped to standard output when done. Unlike in tests, benchmark logs are always printed, so as not to hide output whose existence may be affecting benchmark results. 

####  func (*B) [Attr](https://cs.opensource.google/go/go/+/go1.25.6:src/testing/testing.go;l=1509) ¶ added in go1.25.0
    
    
    func (c *B) Attr(key, value [string](/builtin#string))

Attr emits a test attribute associated with this test. 

The key must not contain whitespace. The value must not contain newlines or carriage returns. 

The meaning of different attribute keys is left up to continuous integration systems and test frameworks. 

Test attributes are emitted immediately in the test log, but they are intended to be treated as unordered. 

####  func (*B) [Chdir](https://cs.opensource.google/go/go/+/go1.25.6:src/testing/testing.go;l=1453) ¶ added in go1.24.0
    
    
    func (c *B) Chdir(dir [string](/builtin#string))

Chdir calls [os.Chdir](/os#Chdir) and uses Cleanup to restore the current working directory to its original value after the test. On Unix, it also sets PWD environment variable for the duration of the test. 

Because Chdir affects the whole process, it cannot be used in parallel tests or tests with parallel ancestors. 

####  func (*B) [Cleanup](https://cs.opensource.google/go/go/+/go1.25.6:src/testing/testing.go;l=1287) ¶ added in go1.14
    
    
    func (c *B) Cleanup(f func())

Cleanup registers a function to be called when the test (or subtest) and all its subtests complete. Cleanup functions will be called in last added, first called order. 

####  func (*B) [Context](https://cs.opensource.google/go/go/+/go1.25.6:src/testing/testing.go;l=1494) ¶ added in go1.24.0
    
    
    func (c *B) Context() [context](/context).[Context](/context#Context)

Context returns a context that is canceled just before Cleanup-registered functions are called. 

Cleanup functions can wait for any resources that shut down on [context.Context.Done](/context#Context.Done) before the test or benchmark completes. 

####  func (*B) [Elapsed](https://cs.opensource.google/go/go/+/go1.25.6:src/testing/benchmark.go;l=364) ¶ added in go1.20
    
    
    func (b *B) Elapsed() [time](/time).[Duration](/time#Duration)

Elapsed returns the measured elapsed time of the benchmark. The duration reported by Elapsed matches the one measured by B.StartTimer, B.StopTimer, and B.ResetTimer. 

####  func (*B) [Error](https://cs.opensource.google/go/go/+/go1.25.6:src/testing/testing.go;l=1195) ¶
    
    
    func (c *B) Error(args ...[any](/builtin#any))

Error is equivalent to Log followed by Fail. 

####  func (*B) [Errorf](https://cs.opensource.google/go/go/+/go1.25.6:src/testing/testing.go;l=1202) ¶
    
    
    func (c *B) Errorf(format [string](/builtin#string), args ...[any](/builtin#any))

Errorf is equivalent to Logf followed by Fail. 

####  func (*B) [Fail](https://cs.opensource.google/go/go/+/go1.25.6:src/testing/testing.go;l=952) ¶
    
    
    func (c *B) Fail()

Fail marks the function as having failed but continues execution. 

####  func (*B) [FailNow](https://cs.opensource.google/go/go/+/go1.25.6:src/testing/testing.go;l=987) ¶
    
    
    func (c *B) FailNow()

FailNow marks the function as having failed and stops its execution by calling [runtime.Goexit](/runtime#Goexit) (which then runs all deferred calls in the current goroutine). Execution will continue at the next test or benchmark. FailNow must be called from the goroutine running the test or benchmark function, not from other goroutines created during the test. Calling FailNow does not stop those other goroutines. 

####  func (*B) [Failed](https://cs.opensource.google/go/go/+/go1.25.6:src/testing/testing.go;l=966) ¶
    
    
    func (c *B) Failed() [bool](/builtin#bool)

Failed reports whether the function has failed. 

####  func (*B) [Fatal](https://cs.opensource.google/go/go/+/go1.25.6:src/testing/testing.go;l=1209) ¶
    
    
    func (c *B) Fatal(args ...[any](/builtin#any))

Fatal is equivalent to Log followed by FailNow. 

####  func (*B) [Fatalf](https://cs.opensource.google/go/go/+/go1.25.6:src/testing/testing.go;l=1216) ¶
    
    
    func (c *B) Fatalf(format [string](/builtin#string), args ...[any](/builtin#any))

Fatalf is equivalent to Logf followed by FailNow. 

####  func (*B) [Helper](https://cs.opensource.google/go/go/+/go1.25.6:src/testing/testing.go;l=1263) ¶ added in go1.9
    
    
    func (c *B) Helper()

Helper marks the calling function as a test helper function. When printing file and line information, that function will be skipped. Helper may be called simultaneously from multiple goroutines. 

####  func (*B) [Log](https://cs.opensource.google/go/go/+/go1.25.6:src/testing/testing.go;l=1178) ¶
    
    
    func (c *B) Log(args ...[any](/builtin#any))

Log formats its arguments using default formatting, analogous to [fmt.Println](/fmt#Println), and records the text in the error log. For tests, the text will be printed only if the test fails or the -test.v flag is set. For benchmarks, the text is always printed to avoid having performance depend on the value of the -test.v flag. It is an error to call Log after a test or benchmark returns. 

####  func (*B) [Logf](https://cs.opensource.google/go/go/+/go1.25.6:src/testing/testing.go;l=1189) ¶
    
    
    func (c *B) Logf(format [string](/builtin#string), args ...[any](/builtin#any))

Logf formats its arguments according to the format, analogous to [fmt.Printf](/fmt#Printf), and records the text in the error log. A final newline is added if not provided. For tests, the text will be printed only if the test fails or the -test.v flag is set. For benchmarks, the text is always printed to avoid having performance depend on the value of the -test.v flag. It is an error to call Logf after a test or benchmark returns. 

####  func (*B) [Loop](https://cs.opensource.google/go/go/+/go1.25.6:src/testing/benchmark.go;l=497) ¶ added in go1.24.0
    
    
    func (b *B) Loop() [bool](/builtin#bool)

Loop returns true as long as the benchmark should continue running. 

A typical benchmark is structured like: 
    
    
    func Benchmark(b *testing.B) {
    	... setup ...
    	for b.Loop() {
    		... code to measure ...
    	}
    	... cleanup ...
    }
    

Loop resets the benchmark timer the first time it is called in a benchmark, so any setup performed prior to starting the benchmark loop does not count toward the benchmark measurement. Likewise, when it returns false, it stops the timer so cleanup code is not measured. 

Within the body of a "for b.Loop() { ... }" loop, arguments to and results from function calls within the loop are kept alive, preventing the compiler from fully optimizing away the loop body. Currently, this is implemented by disabling inlining of functions called in a b.Loop loop. This applies only to calls syntactically between the curly braces of the loop, and the loop condition must be written exactly as "b.Loop()". Optimizations are performed as usual in any functions called by the loop. 

After Loop returns false, b.N contains the total number of iterations that ran, so the benchmark may use b.N to compute other average metrics. 

Prior to the introduction of Loop, benchmarks were expected to contain an explicit loop from 0 to b.N. Benchmarks should either use Loop or contain a loop to b.N, but not both. Loop offers more automatic management of the benchmark timer, and runs each benchmark function only once per measurement, whereas b.N-based benchmarks must run the benchmark function (and any associated setup and cleanup) several times. 

Example ¶
    
    
    package main
    
    import (
    	"math/rand/v2"
    	"testing"
    )
    
    // ExBenchmark shows how to use b.Loop in a benchmark.
    //
    // (If this were a real benchmark, not an example, this would be named
    // BenchmarkSomething.)
    func ExBenchmark(b *testing.B) {
    	// Generate a large random slice to use as an input.
    	// Since this is done before the first call to b.Loop(),
    	// it doesn't count toward the benchmark time.
    	input := make([]int, 128<<10)
    	for i := range input {
    		input[i] = rand.Int()
    	}
    
    	// Perform the benchmark.
    	for b.Loop() {
    		// Normally, the compiler would be allowed to optimize away the call
    		// to sum because it has no side effects and the result isn't used.
    		// However, inside a b.Loop loop, the compiler ensures function calls
    		// aren't optimized away.
    		sum(input)
    	}
    
    	// Outside the loop, the timer is stopped, so we could perform
    	// cleanup if necessary without affecting the result.
    }
    
    func sum(data []int) int {
    	total := 0
    	for _, value := range data {
    		total += value
    	}
    	return total
    }
    
    func main() {
    	testing.Benchmark(ExBenchmark)
    }
    

Share Format Run

####  func (*B) [Name](https://cs.opensource.google/go/go/+/go1.25.6:src/testing/testing.go;l=938) ¶ added in go1.8
    
    
    func (c *B) Name() [string](/builtin#string)

Name returns the name of the running (sub-) test or benchmark. 

The name will include the name of the test along with the names of any nested sub-tests. If two sibling sub-tests have the same name, Name will append a suffix to guarantee the returned name is unique. 

####  func (*B) [Output](https://cs.opensource.google/go/go/+/go1.25.6:src/testing/testing.go;l=1105) ¶ added in go1.25.0
    
    
    func (c *B) Output() [io](/io).[Writer](/io#Writer)

Output returns a Writer that writes to the same test output stream as TB.Log. The output is indented like TB.Log lines, but Output does not add source locations or newlines. The output is internally line buffered, and a call to TB.Log or the end of the test will implicitly flush the buffer, followed by a newline. After a test function and all its parents return, neither Output nor the Write method may be called. 

####  func (*B) [ReportAllocs](https://cs.opensource.google/go/go/+/go1.25.6:src/testing/benchmark.go;l=192) ¶ added in go1.1
    
    
    func (b *B) ReportAllocs()

ReportAllocs enables malloc statistics for this benchmark. It is equivalent to setting -test.benchmem, but it only affects the benchmark function that calls ReportAllocs. 

####  func (*B) [ReportMetric](https://cs.opensource.google/go/go/+/go1.25.6:src/testing/benchmark.go;l=381) ¶ added in go1.13
    
    
    func (b *B) ReportMetric(n [float64](/builtin#float64), unit [string](/builtin#string))

ReportMetric adds "n unit" to the reported benchmark results. If the metric is per-iteration, the caller should divide by b.N, and by convention units should end in "/op". ReportMetric overrides any previously reported value for the same unit. ReportMetric panics if unit is the empty string or if unit contains any whitespace. If unit is a unit normally reported by the benchmark framework itself (such as "allocs/op"), ReportMetric will override that metric. Setting "ns/op" to 0 will suppress that built-in metric. 

Example ¶
    
    
    package main
    
    import (
    	"cmp"
    	"slices"
    	"testing"
    )
    
    func main() {
    	// This reports a custom benchmark metric relevant to a
    	// specific algorithm (in this case, sorting).
    	testing.Benchmark(func(b *testing.B) {
    		var compares int64
    		for b.Loop() {
    			s := []int{5, 4, 3, 2, 1}
    			slices.SortFunc(s, func(a, b int) int {
    				compares++
    				return cmp.Compare(a, b)
    			})
    		}
    		// This metric is per-operation, so divide by b.N and
    		// report it as a "/op" unit.
    		b.ReportMetric(float64(compares)/float64(b.N), "compares/op")
    		// This metric is per-time, so divide by b.Elapsed and
    		// report it as a "/ns" unit.
    		b.ReportMetric(float64(compares)/float64(b.Elapsed().Nanoseconds()), "compares/ns")
    	})
    }
    

Share Format Run

Example (Parallel) ¶
    
    
    package main
    
    import (
    	"cmp"
    	"slices"
    	"sync/atomic"
    	"testing"
    )
    
    func main() {
    	// This reports a custom benchmark metric relevant to a
    	// specific algorithm (in this case, sorting) in parallel.
    	testing.Benchmark(func(b *testing.B) {
    		var compares atomic.Int64
    		b.RunParallel(func(pb *testing.PB) {
    			for pb.Next() {
    				s := []int{5, 4, 3, 2, 1}
    				slices.SortFunc(s, func(a, b int) int {
    					// Because RunParallel runs the function many
    					// times in parallel, we must increment the
    					// counter atomically to avoid racing writes.
    					compares.Add(1)
    					return cmp.Compare(a, b)
    				})
    			}
    		})
    
    		// NOTE: Report each metric once, after all of the parallel
    		// calls have completed.
    
    		// This metric is per-operation, so divide by b.N and
    		// report it as a "/op" unit.
    		b.ReportMetric(float64(compares.Load())/float64(b.N), "compares/op")
    		// This metric is per-time, so divide by b.Elapsed and
    		// report it as a "/ns" unit.
    		b.ReportMetric(float64(compares.Load())/float64(b.Elapsed().Nanoseconds()), "compares/ns")
    	})
    }
    

Share Format Run

####  func (*B) [ResetTimer](https://cs.opensource.google/go/go/+/go1.25.6:src/testing/benchmark.go;l=166) ¶
    
    
    func (b *B) ResetTimer()

ResetTimer zeroes the elapsed benchmark time and memory allocation counters and deletes user-reported metrics. It does not affect whether the timer is running. 

####  func (*B) [Run](https://cs.opensource.google/go/go/+/go1.25.6:src/testing/benchmark.go;l=803) ¶ added in go1.7
    
    
    func (b *B) Run(name [string](/builtin#string), f func(b *B)) [bool](/builtin#bool)

Run benchmarks f as a subbenchmark with the given name. It reports whether there were any failures. 

A subbenchmark is like any other benchmark. A benchmark that calls Run at least once will not be measured itself and will be called once with N=1. 

####  func (*B) [RunParallel](https://cs.opensource.google/go/go/+/go1.25.6:src/testing/benchmark.go;l=945) ¶ added in go1.3
    
    
    func (b *B) RunParallel(body func(*PB))

RunParallel runs a benchmark in parallel. It creates multiple goroutines and distributes b.N iterations among them. The number of goroutines defaults to GOMAXPROCS. To increase parallelism for non-CPU-bound benchmarks, call B.SetParallelism before RunParallel. RunParallel is usually used with the go test -cpu flag. 

The body function will be run in each goroutine. It should set up any goroutine-local state and then iterate until pb.Next returns false. It should not use the B.StartTimer, B.StopTimer, or B.ResetTimer functions, because they have global effect. It should also not call B.Run. 

RunParallel reports ns/op values as wall time for the benchmark as a whole, not the sum of wall time or CPU time over each parallel goroutine. 

Example ¶
    
    
    package main
    
    import (
    	"bytes"
    	"testing"
    	"text/template"
    )
    
    func main() {
    	// Parallel benchmark for text/template.Template.Execute on a single object.
    	testing.Benchmark(func(b *testing.B) {
    		templ := template.Must(template.New("test").Parse("Hello, {{.}}!"))
    		// RunParallel will create GOMAXPROCS goroutines
    		// and distribute work among them.
    		b.RunParallel(func(pb *testing.PB) {
    			// Each goroutine has its own bytes.Buffer.
    			var buf bytes.Buffer
    			for pb.Next() {
    				// The loop body is executed b.N times total across all goroutines.
    				buf.Reset()
    				templ.Execute(&buf, "World")
    			}
    		})
    	})
    }
    

Share Format Run

####  func (*B) [SetBytes](https://cs.opensource.google/go/go/+/go1.25.6:src/testing/benchmark.go;l=187) ¶
    
    
    func (b *B) SetBytes(n [int64](/builtin#int64))

SetBytes records the number of bytes processed in a single operation. If this is called, the benchmark will report ns/op and MB/s. 

####  func (*B) [SetParallelism](https://cs.opensource.google/go/go/+/go1.25.6:src/testing/benchmark.go;l=989) ¶ added in go1.3
    
    
    func (b *B) SetParallelism(p [int](/builtin#int))

SetParallelism sets the number of goroutines used by B.RunParallel to p*GOMAXPROCS. There is usually no need to call SetParallelism for CPU-bound benchmarks. If p is less than 1, this call will have no effect. 

####  func (*B) [Setenv](https://cs.opensource.google/go/go/+/go1.25.6:src/testing/testing.go;l=1428) ¶ added in go1.17
    
    
    func (c *B) Setenv(key, value [string](/builtin#string))

Setenv calls [os.Setenv](/os#Setenv) and uses Cleanup to restore the environment variable to its original value after the test. 

Because Setenv affects the whole process, it cannot be used in parallel tests or tests with parallel ancestors. 

####  func (*B) [Skip](https://cs.opensource.google/go/go/+/go1.25.6:src/testing/testing.go;l=1223) ¶ added in go1.1
    
    
    func (c *B) Skip(args ...[any](/builtin#any))

Skip is equivalent to Log followed by SkipNow. 

####  func (*B) [SkipNow](https://cs.opensource.google/go/go/+/go1.25.6:src/testing/testing.go;l=1244) ¶ added in go1.1
    
    
    func (c *B) SkipNow()

SkipNow marks the test as having been skipped and stops its execution by calling [runtime.Goexit](/runtime#Goexit). If a test fails (see Error, Errorf, Fail) and is then skipped, it is still considered to have failed. Execution will continue at the next test or benchmark. See also FailNow. SkipNow must be called from the goroutine running the test, not from other goroutines created during the test. Calling SkipNow does not stop those other goroutines. 

####  func (*B) [Skipf](https://cs.opensource.google/go/go/+/go1.25.6:src/testing/testing.go;l=1230) ¶ added in go1.1
    
    
    func (c *B) Skipf(format [string](/builtin#string), args ...[any](/builtin#any))

Skipf is equivalent to Logf followed by SkipNow. 

####  func (*B) [Skipped](https://cs.opensource.google/go/go/+/go1.25.6:src/testing/testing.go;l=1254) ¶ added in go1.1
    
    
    func (c *B) Skipped() [bool](/builtin#bool)

Skipped reports whether the test was skipped. 

####  func (*B) [StartTimer](https://cs.opensource.google/go/go/+/go1.25.6:src/testing/benchmark.go;l=138) ¶
    
    
    func (b *B) StartTimer()

StartTimer starts timing a test. This function is called automatically before a benchmark starts, but it can also be used to resume timing after a call to B.StopTimer. 

####  func (*B) [StopTimer](https://cs.opensource.google/go/go/+/go1.25.6:src/testing/benchmark.go;l=151) ¶
    
    
    func (b *B) StopTimer()

StopTimer stops timing a test. This can be used to pause the timer while performing steps that you don't want to measure. 

####  func (*B) [TempDir](https://cs.opensource.google/go/go/+/go1.25.6:src/testing/testing.go;l=1321) ¶ added in go1.15
    
    
    func (c *B) TempDir() [string](/builtin#string)

TempDir returns a temporary directory for the test to use. The directory is automatically removed when the test and all its subtests complete. Each subsequent call to TempDir returns a unique directory; if the directory creation fails, TempDir terminates the test by calling Fatal. 

####  type [BenchmarkResult](https://cs.opensource.google/go/go/+/go1.25.6:src/testing/benchmark.go;l=531) ¶
    
    
    type BenchmarkResult struct {
    	N         [int](/builtin#int)           // The number of iterations.
    	T         [time](/time).[Duration](/time#Duration) // The total time taken.
    	Bytes     [int64](/builtin#int64)         // Bytes processed in one iteration.
    	MemAllocs [uint64](/builtin#uint64)        // The total number of memory allocations.
    	MemBytes  [uint64](/builtin#uint64)        // The total number of bytes allocated.
    
    	// Extra records additional metrics reported by ReportMetric.
    	Extra map[[string](/builtin#string)][float64](/builtin#float64)
    }

BenchmarkResult contains the results of a benchmark run. 

####  func [Benchmark](https://cs.opensource.google/go/go/+/go1.25.6:src/testing/benchmark.go;l=1003) ¶
    
    
    func Benchmark(f func(b *B)) BenchmarkResult

Benchmark benchmarks a single function. It is useful for creating custom benchmarks that do not use the "go test" command. 

If f depends on testing flags, then Init must be used to register those flags before calling Benchmark and before calling [flag.Parse](/flag#Parse). 

If f calls Run, the result will be an estimate of running all its subbenchmarks that don't call Run in sequence in a single benchmark. 

####  func (BenchmarkResult) [AllocedBytesPerOp](https://cs.opensource.google/go/go/+/go1.25.6:src/testing/benchmark.go;l=578) ¶ added in go1.1
    
    
    func (r BenchmarkResult) AllocedBytesPerOp() [int64](/builtin#int64)

AllocedBytesPerOp returns the "B/op" metric, which is calculated as r.MemBytes / r.N. 

####  func (BenchmarkResult) [AllocsPerOp](https://cs.opensource.google/go/go/+/go1.25.6:src/testing/benchmark.go;l=566) ¶ added in go1.1
    
    
    func (r BenchmarkResult) AllocsPerOp() [int64](/builtin#int64)

AllocsPerOp returns the "allocs/op" metric, which is calculated as r.MemAllocs / r.N. 

####  func (BenchmarkResult) [MemString](https://cs.opensource.google/go/go/+/go1.25.6:src/testing/benchmark.go;l=660) ¶ added in go1.1
    
    
    func (r BenchmarkResult) MemString() [string](/builtin#string)

MemString returns r.AllocedBytesPerOp and r.AllocsPerOp in the same format as 'go test'. 

####  func (BenchmarkResult) [NsPerOp](https://cs.opensource.google/go/go/+/go1.25.6:src/testing/benchmark.go;l=543) ¶
    
    
    func (r BenchmarkResult) NsPerOp() [int64](/builtin#int64)

NsPerOp returns the "ns/op" metric. 

####  func (BenchmarkResult) [String](https://cs.opensource.google/go/go/+/go1.25.6:src/testing/benchmark.go;l=595) ¶
    
    
    func (r BenchmarkResult) String() [string](/builtin#string)

String returns a summary of the benchmark results. It follows the benchmark result line format from <https://golang.org/design/14313-benchmark-format>, not including the benchmark name. Extra metrics override built-in metrics of the same name. String does not include allocs/op or B/op, since those are reported by BenchmarkResult.MemString. 

####  type [Cover](https://cs.opensource.google/go/go/+/go1.25.6:src/testing/cover.go;l=26) ¶ added in go1.2
    
    
    type Cover struct {
    	Mode            [string](/builtin#string)
    	Counters        map[[string](/builtin#string)][][uint32](/builtin#uint32)
    	Blocks          map[[string](/builtin#string)][]CoverBlock
    	CoveredPackages [string](/builtin#string)
    }

Cover records information about test coverage checking. NOTE: This struct is internal to the testing infrastructure and may change. It is not covered (yet) by the Go 1 compatibility guidelines. 

####  type [CoverBlock](https://cs.opensource.google/go/go/+/go1.25.6:src/testing/cover.go;l=15) ¶ added in go1.2
    
    
    type CoverBlock struct {
    	Line0 [uint32](/builtin#uint32) // Line number for block start.
    	Col0  [uint16](/builtin#uint16) // Column number for block start.
    	Line1 [uint32](/builtin#uint32) // Line number for block end.
    	Col1  [uint16](/builtin#uint16) // Column number for block end.
    	Stmts [uint16](/builtin#uint16) // Number of statements included in this block.
    }

CoverBlock records the coverage data for a single basic block. The fields are 1-indexed, as in an editor: The opening line of the file is number 1, for example. Columns are measured in bytes. NOTE: This struct is internal to the testing infrastructure and may change. It is not covered (yet) by the Go 1 compatibility guidelines. 

####  type [F](https://cs.opensource.google/go/go/+/go1.25.6:src/testing/fuzz.go;l=69) ¶ added in go1.18
    
    
    type F struct {
    	// contains filtered or unexported fields
    }

F is a type passed to fuzz tests. 

Fuzz tests run generated inputs against a provided fuzz target, which can find and report potential bugs in the code being tested. 

A fuzz test runs the seed corpus by default, which includes entries provided by F.Add and entries in the testdata/fuzz/<FuzzTestName> directory. After any necessary setup and calls to F.Add, the fuzz test must then call F.Fuzz to provide the fuzz target. See the testing package documentation for an example, and see the F.Fuzz and F.Add method documentation for details. 

*F methods can only be called before F.Fuzz. Once the test is executing the fuzz target, only *T methods can be used. The only *F methods that are allowed in the F.Fuzz function are F.Failed and F.Name. 

####  func (*F) [Add](https://cs.opensource.google/go/go/+/go1.25.6:src/testing/fuzz.go;l=153) ¶ added in go1.18
    
    
    func (f *F) Add(args ...[any](/builtin#any))

Add will add the arguments to the seed corpus for the fuzz test. This will be a no-op if called after or within the fuzz target, and args must match the arguments for the fuzz target. 

####  func (*F) [Attr](https://cs.opensource.google/go/go/+/go1.25.6:src/testing/testing.go;l=1509) ¶ added in go1.25.0
    
    
    func (c *F) Attr(key, value [string](/builtin#string))

Attr emits a test attribute associated with this test. 

The key must not contain whitespace. The value must not contain newlines or carriage returns. 

The meaning of different attribute keys is left up to continuous integration systems and test frameworks. 

Test attributes are emitted immediately in the test log, but they are intended to be treated as unordered. 

####  func (*F) [Chdir](https://cs.opensource.google/go/go/+/go1.25.6:src/testing/testing.go;l=1453) ¶ added in go1.24.0
    
    
    func (c *F) Chdir(dir [string](/builtin#string))

Chdir calls [os.Chdir](/os#Chdir) and uses Cleanup to restore the current working directory to its original value after the test. On Unix, it also sets PWD environment variable for the duration of the test. 

Because Chdir affects the whole process, it cannot be used in parallel tests or tests with parallel ancestors. 

####  func (*F) [Cleanup](https://cs.opensource.google/go/go/+/go1.25.6:src/testing/testing.go;l=1287) ¶ added in go1.18
    
    
    func (c *F) Cleanup(f func())

Cleanup registers a function to be called when the test (or subtest) and all its subtests complete. Cleanup functions will be called in last added, first called order. 

####  func (*F) [Context](https://cs.opensource.google/go/go/+/go1.25.6:src/testing/testing.go;l=1494) ¶ added in go1.24.0
    
    
    func (c *F) Context() [context](/context).[Context](/context#Context)

Context returns a context that is canceled just before Cleanup-registered functions are called. 

Cleanup functions can wait for any resources that shut down on [context.Context.Done](/context#Context.Done) before the test or benchmark completes. 

####  func (*F) [Error](https://cs.opensource.google/go/go/+/go1.25.6:src/testing/testing.go;l=1195) ¶ added in go1.18
    
    
    func (c *F) Error(args ...[any](/builtin#any))

Error is equivalent to Log followed by Fail. 

####  func (*F) [Errorf](https://cs.opensource.google/go/go/+/go1.25.6:src/testing/testing.go;l=1202) ¶ added in go1.18
    
    
    func (c *F) Errorf(format [string](/builtin#string), args ...[any](/builtin#any))

Errorf is equivalent to Logf followed by Fail. 

####  func (*F) [Fail](https://cs.opensource.google/go/go/+/go1.25.6:src/testing/fuzz.go;l=129) ¶ added in go1.18
    
    
    func (f *F) Fail()

Fail marks the function as having failed but continues execution. 

####  func (*F) [FailNow](https://cs.opensource.google/go/go/+/go1.25.6:src/testing/testing.go;l=987) ¶ added in go1.18
    
    
    func (c *F) FailNow()

FailNow marks the function as having failed and stops its execution by calling [runtime.Goexit](/runtime#Goexit) (which then runs all deferred calls in the current goroutine). Execution will continue at the next test or benchmark. FailNow must be called from the goroutine running the test or benchmark function, not from other goroutines created during the test. Calling FailNow does not stop those other goroutines. 

####  func (*F) [Failed](https://cs.opensource.google/go/go/+/go1.25.6:src/testing/testing.go;l=966) ¶ added in go1.18
    
    
    func (c *F) Failed() [bool](/builtin#bool)

Failed reports whether the function has failed. 

####  func (*F) [Fatal](https://cs.opensource.google/go/go/+/go1.25.6:src/testing/testing.go;l=1209) ¶ added in go1.18
    
    
    func (c *F) Fatal(args ...[any](/builtin#any))

Fatal is equivalent to Log followed by FailNow. 

####  func (*F) [Fatalf](https://cs.opensource.google/go/go/+/go1.25.6:src/testing/testing.go;l=1216) ¶ added in go1.18
    
    
    func (c *F) Fatalf(format [string](/builtin#string), args ...[any](/builtin#any))

Fatalf is equivalent to Logf followed by FailNow. 

####  func (*F) [Fuzz](https://cs.opensource.google/go/go/+/go1.25.6:src/testing/fuzz.go;l=211) ¶ added in go1.18
    
    
    func (f *F) Fuzz(ff [any](/builtin#any))

Fuzz runs the fuzz function, ff, for fuzz testing. If ff fails for a set of arguments, those arguments will be added to the seed corpus. 

ff must be a function with no return value whose first argument is *T and whose remaining arguments are the types to be fuzzed. For example: 
    
    
    f.Fuzz(func(t *testing.T, b []byte, i int) { ... })
    

The following types are allowed: []byte, string, bool, byte, rune, float32, float64, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64. More types may be supported in the future. 

ff must not call any *F methods, e.g. F.Log, F.Error, F.Skip. Use the corresponding *T method instead. The only *F methods that are allowed in the F.Fuzz function are F.Failed and F.Name. 

This function should be fast and deterministic, and its behavior should not depend on shared state. No mutable input arguments, or pointers to them, should be retained between executions of the fuzz function, as the memory backing them may be mutated during a subsequent invocation. ff must not modify the underlying data of the arguments provided by the fuzzing engine. 

When fuzzing, F.Fuzz does not return until a problem is found, time runs out (set with -fuzztime), or the test process is interrupted by a signal. F.Fuzz should be called exactly once, unless F.Skip or F.Fail is called beforehand. 

####  func (*F) [Helper](https://cs.opensource.google/go/go/+/go1.25.6:src/testing/fuzz.go;l=103) ¶ added in go1.18
    
    
    func (f *F) Helper()

Helper marks the calling function as a test helper function. When printing file and line information, that function will be skipped. Helper may be called simultaneously from multiple goroutines. 

####  func (*F) [Log](https://cs.opensource.google/go/go/+/go1.25.6:src/testing/testing.go;l=1178) ¶ added in go1.18
    
    
    func (c *F) Log(args ...[any](/builtin#any))

Log formats its arguments using default formatting, analogous to [fmt.Println](/fmt#Println), and records the text in the error log. For tests, the text will be printed only if the test fails or the -test.v flag is set. For benchmarks, the text is always printed to avoid having performance depend on the value of the -test.v flag. It is an error to call Log after a test or benchmark returns. 

####  func (*F) [Logf](https://cs.opensource.google/go/go/+/go1.25.6:src/testing/testing.go;l=1189) ¶ added in go1.18
    
    
    func (c *F) Logf(format [string](/builtin#string), args ...[any](/builtin#any))

Logf formats its arguments according to the format, analogous to [fmt.Printf](/fmt#Printf), and records the text in the error log. A final newline is added if not provided. For tests, the text will be printed only if the test fails or the -test.v flag is set. For benchmarks, the text is always printed to avoid having performance depend on the value of the -test.v flag. It is an error to call Logf after a test or benchmark returns. 

####  func (*F) [Name](https://cs.opensource.google/go/go/+/go1.25.6:src/testing/testing.go;l=938) ¶ added in go1.18
    
    
    func (c *F) Name() [string](/builtin#string)

Name returns the name of the running (sub-) test or benchmark. 

The name will include the name of the test along with the names of any nested sub-tests. If two sibling sub-tests have the same name, Name will append a suffix to guarantee the returned name is unique. 

####  func (*F) [Output](https://cs.opensource.google/go/go/+/go1.25.6:src/testing/testing.go;l=1105) ¶ added in go1.25.0
    
    
    func (c *F) Output() [io](/io).[Writer](/io#Writer)

Output returns a Writer that writes to the same test output stream as TB.Log. The output is indented like TB.Log lines, but Output does not add source locations or newlines. The output is internally line buffered, and a call to TB.Log or the end of the test will implicitly flush the buffer, followed by a newline. After a test function and all its parents return, neither Output nor the Write method may be called. 

####  func (*F) [Setenv](https://cs.opensource.google/go/go/+/go1.25.6:src/testing/testing.go;l=1428) ¶ added in go1.18
    
    
    func (c *F) Setenv(key, value [string](/builtin#string))

Setenv calls [os.Setenv](/os#Setenv) and uses Cleanup to restore the environment variable to its original value after the test. 

Because Setenv affects the whole process, it cannot be used in parallel tests or tests with parallel ancestors. 

####  func (*F) [Skip](https://cs.opensource.google/go/go/+/go1.25.6:src/testing/testing.go;l=1223) ¶ added in go1.18
    
    
    func (c *F) Skip(args ...[any](/builtin#any))

Skip is equivalent to Log followed by SkipNow. 

####  func (*F) [SkipNow](https://cs.opensource.google/go/go/+/go1.25.6:src/testing/testing.go;l=1244) ¶ added in go1.18
    
    
    func (c *F) SkipNow()

SkipNow marks the test as having been skipped and stops its execution by calling [runtime.Goexit](/runtime#Goexit). If a test fails (see Error, Errorf, Fail) and is then skipped, it is still considered to have failed. Execution will continue at the next test or benchmark. See also FailNow. SkipNow must be called from the goroutine running the test, not from other goroutines created during the test. Calling SkipNow does not stop those other goroutines. 

####  func (*F) [Skipf](https://cs.opensource.google/go/go/+/go1.25.6:src/testing/testing.go;l=1230) ¶ added in go1.18
    
    
    func (c *F) Skipf(format [string](/builtin#string), args ...[any](/builtin#any))

Skipf is equivalent to Logf followed by SkipNow. 

####  func (*F) [Skipped](https://cs.opensource.google/go/go/+/go1.25.6:src/testing/fuzz.go;l=140) ¶ added in go1.18
    
    
    func (f *F) Skipped() [bool](/builtin#bool)

Skipped reports whether the test was skipped. 

####  func (*F) [TempDir](https://cs.opensource.google/go/go/+/go1.25.6:src/testing/testing.go;l=1321) ¶ added in go1.18
    
    
    func (c *F) TempDir() [string](/builtin#string)

TempDir returns a temporary directory for the test to use. The directory is automatically removed when the test and all its subtests complete. Each subsequent call to TempDir returns a unique directory; if the directory creation fails, TempDir terminates the test by calling Fatal. 

####  type [InternalBenchmark](https://cs.opensource.google/go/go/+/go1.25.6:src/testing/benchmark.go;l=76) ¶
    
    
    type InternalBenchmark struct {
    	Name [string](/builtin#string)
    	F    func(b *B)
    }

InternalBenchmark is an internal type but exported because it is cross-package; it is part of the implementation of the "go test" command. 

####  type [InternalExample](https://cs.opensource.google/go/go/+/go1.25.6:src/testing/example.go;l=15) ¶
    
    
    type InternalExample struct {
    	Name      [string](/builtin#string)
    	F         func()
    	Output    [string](/builtin#string)
    	Unordered [bool](/builtin#bool)
    }

####  type [InternalFuzzTarget](https://cs.opensource.google/go/go/+/go1.25.6:src/testing/fuzz.go;l=49) ¶ added in go1.18
    
    
    type InternalFuzzTarget struct {
    	Name [string](/builtin#string)
    	Fn   func(f *F)
    }

InternalFuzzTarget is an internal type but exported because it is cross-package; it is part of the implementation of the "go test" command. 

####  type [InternalTest](https://cs.opensource.google/go/go/+/go1.25.6:src/testing/testing.go;l=1767) ¶
    
    
    type InternalTest struct {
    	Name [string](/builtin#string)
    	F    func(*T)
    }

InternalTest is an internal type but exported because it is cross-package; it is part of the implementation of the "go test" command. 

####  type [M](https://cs.opensource.google/go/go/+/go1.25.6:src/testing/testing.go;l=2171) ¶ added in go1.4
    
    
    type M struct {
    	// contains filtered or unexported fields
    }

M is a type passed to a TestMain function to run the actual tests. 

####  func [MainStart](https://cs.opensource.google/go/go/+/go1.25.6:src/testing/testing.go;l=2213) ¶ added in go1.4
    
    
    func MainStart(deps testDeps, tests []InternalTest, benchmarks []InternalBenchmark, fuzzTargets []InternalFuzzTarget, examples []InternalExample) *M

MainStart is meant for use by tests generated by 'go test'. It is not meant to be called directly and is not subject to the Go 1 compatibility document. It may change signature from release to release. 

####  func (*M) [Run](https://cs.opensource.google/go/go/+/go1.25.6:src/testing/testing.go;l=2234) ¶ added in go1.4
    
    
    func (m *M) Run() (code [int](/builtin#int))

Run runs the tests. It returns an exit code to pass to os.Exit. The exit code is zero when all tests pass, and non-zero for any kind of failure. For machine readable test results, parse the output of 'go test -json'. 

####  type [PB](https://cs.opensource.google/go/go/+/go1.25.6:src/testing/benchmark.go;l=909) ¶ added in go1.3
    
    
    type PB struct {
    	// contains filtered or unexported fields
    }

A PB is used by RunParallel for running parallel benchmarks. 

####  func (*PB) [Next](https://cs.opensource.google/go/go/+/go1.25.6:src/testing/benchmark.go;l=917) ¶ added in go1.3
    
    
    func (pb *PB) Next() [bool](/builtin#bool)

Next reports whether there are more iterations to execute. 

####  type [T](https://cs.opensource.google/go/go/+/go1.25.6:src/testing/testing.go;l=925) ¶
    
    
    type T struct {
    	// contains filtered or unexported fields
    }

T is a type passed to Test functions to manage test state and support formatted test logs. 

A test ends when its Test function returns or calls any of the methods T.FailNow, T.Fatal, T.Fatalf, T.SkipNow, T.Skip, or T.Skipf. Those methods, as well as the T.Parallel method, must be called only from the goroutine running the Test function. 

The other reporting methods, such as the variations of T.Log and T.Error, may be called simultaneously from multiple goroutines. 

####  func (*T) [Attr](https://cs.opensource.google/go/go/+/go1.25.6:src/testing/testing.go;l=1509) ¶ added in go1.25.0
    
    
    func (c *T) Attr(key, value [string](/builtin#string))

Attr emits a test attribute associated with this test. 

The key must not contain whitespace. The value must not contain newlines or carriage returns. 

The meaning of different attribute keys is left up to continuous integration systems and test frameworks. 

Test attributes are emitted immediately in the test log, but they are intended to be treated as unordered. 

####  func (*T) [Chdir](https://cs.opensource.google/go/go/+/go1.25.6:src/testing/testing.go;l=1760) ¶ added in go1.24.0
    
    
    func (t *T) Chdir(dir [string](/builtin#string))

Chdir calls [os.Chdir](/os#Chdir) and uses Cleanup to restore the current working directory to its original value after the test. On Unix, it also sets PWD environment variable for the duration of the test. 

Because Chdir affects the whole process, it cannot be used in parallel tests or tests with parallel ancestors. 

####  func (*T) [Cleanup](https://cs.opensource.google/go/go/+/go1.25.6:src/testing/testing.go;l=1287) ¶ added in go1.14
    
    
    func (c *T) Cleanup(f func())

Cleanup registers a function to be called when the test (or subtest) and all its subtests complete. Cleanup functions will be called in last added, first called order. 

####  func (*T) [Context](https://cs.opensource.google/go/go/+/go1.25.6:src/testing/testing.go;l=1494) ¶ added in go1.24.0
    
    
    func (c *T) Context() [context](/context).[Context](/context#Context)

Context returns a context that is canceled just before Cleanup-registered functions are called. 

Cleanup functions can wait for any resources that shut down on [context.Context.Done](/context#Context.Done) before the test or benchmark completes. 

####  func (*T) [Deadline](https://cs.opensource.google/go/go/+/go1.25.6:src/testing/testing.go;l=2059) ¶ added in go1.15
    
    
    func (t *T) Deadline() (deadline [time](/time).[Time](/time#Time), ok [bool](/builtin#bool))

Deadline reports the time at which the test binary will have exceeded the timeout specified by the -timeout flag. 

The ok result is false if the -timeout flag indicates “no timeout” (0). 

####  func (*T) [Error](https://cs.opensource.google/go/go/+/go1.25.6:src/testing/testing.go;l=1195) ¶
    
    
    func (c *T) Error(args ...[any](/builtin#any))

Error is equivalent to Log followed by Fail. 

####  func (*T) [Errorf](https://cs.opensource.google/go/go/+/go1.25.6:src/testing/testing.go;l=1202) ¶
    
    
    func (c *T) Errorf(format [string](/builtin#string), args ...[any](/builtin#any))

Errorf is equivalent to Logf followed by Fail. 

####  func (*T) [Fail](https://cs.opensource.google/go/go/+/go1.25.6:src/testing/testing.go;l=952) ¶
    
    
    func (c *T) Fail()

Fail marks the function as having failed but continues execution. 

####  func (*T) [FailNow](https://cs.opensource.google/go/go/+/go1.25.6:src/testing/testing.go;l=987) ¶
    
    
    func (c *T) FailNow()

FailNow marks the function as having failed and stops its execution by calling [runtime.Goexit](/runtime#Goexit) (which then runs all deferred calls in the current goroutine). Execution will continue at the next test or benchmark. FailNow must be called from the goroutine running the test or benchmark function, not from other goroutines created during the test. Calling FailNow does not stop those other goroutines. 

####  func (*T) [Failed](https://cs.opensource.google/go/go/+/go1.25.6:src/testing/testing.go;l=966) ¶
    
    
    func (c *T) Failed() [bool](/builtin#bool)

Failed reports whether the function has failed. 

####  func (*T) [Fatal](https://cs.opensource.google/go/go/+/go1.25.6:src/testing/testing.go;l=1209) ¶
    
    
    func (c *T) Fatal(args ...[any](/builtin#any))

Fatal is equivalent to Log followed by FailNow. 

####  func (*T) [Fatalf](https://cs.opensource.google/go/go/+/go1.25.6:src/testing/testing.go;l=1216) ¶
    
    
    func (c *T) Fatalf(format [string](/builtin#string), args ...[any](/builtin#any))

Fatalf is equivalent to Logf followed by FailNow. 

####  func (*T) [Helper](https://cs.opensource.google/go/go/+/go1.25.6:src/testing/testing.go;l=1263) ¶ added in go1.9
    
    
    func (c *T) Helper()

Helper marks the calling function as a test helper function. When printing file and line information, that function will be skipped. Helper may be called simultaneously from multiple goroutines. 

####  func (*T) [Log](https://cs.opensource.google/go/go/+/go1.25.6:src/testing/testing.go;l=1178) ¶
    
    
    func (c *T) Log(args ...[any](/builtin#any))

Log formats its arguments using default formatting, analogous to [fmt.Println](/fmt#Println), and records the text in the error log. For tests, the text will be printed only if the test fails or the -test.v flag is set. For benchmarks, the text is always printed to avoid having performance depend on the value of the -test.v flag. It is an error to call Log after a test or benchmark returns. 

####  func (*T) [Logf](https://cs.opensource.google/go/go/+/go1.25.6:src/testing/testing.go;l=1189) ¶
    
    
    func (c *T) Logf(format [string](/builtin#string), args ...[any](/builtin#any))

Logf formats its arguments according to the format, analogous to [fmt.Printf](/fmt#Printf), and records the text in the error log. A final newline is added if not provided. For tests, the text will be printed only if the test fails or the -test.v flag is set. For benchmarks, the text is always printed to avoid having performance depend on the value of the -test.v flag. It is an error to call Logf after a test or benchmark returns. 

####  func (*T) [Name](https://cs.opensource.google/go/go/+/go1.25.6:src/testing/testing.go;l=938) ¶ added in go1.8
    
    
    func (c *T) Name() [string](/builtin#string)

Name returns the name of the running (sub-) test or benchmark. 

The name will include the name of the test along with the names of any nested sub-tests. If two sibling sub-tests have the same name, Name will append a suffix to guarantee the returned name is unique. 

####  func (*T) [Output](https://cs.opensource.google/go/go/+/go1.25.6:src/testing/testing.go;l=1105) ¶ added in go1.25.0
    
    
    func (c *T) Output() [io](/io).[Writer](/io#Writer)

Output returns a Writer that writes to the same test output stream as TB.Log. The output is indented like TB.Log lines, but Output does not add source locations or newlines. The output is internally line buffered, and a call to TB.Log or the end of the test will implicitly flush the buffer, followed by a newline. After a test function and all its parents return, neither Output nor the Write method may be called. 

####  func (*T) [Parallel](https://cs.opensource.google/go/go/+/go1.25.6:src/testing/testing.go;l=1663) ¶
    
    
    func (t *T) Parallel()

Parallel signals that this test is to be run in parallel with (and only with) other parallel tests. When a test is run multiple times due to use of -test.count or -test.cpu, multiple instances of a single test never run in parallel with each other. 

####  func (*T) [Run](https://cs.opensource.google/go/go/+/go1.25.6:src/testing/testing.go;l=1948) ¶ added in go1.7
    
    
    func (t *T) Run(name [string](/builtin#string), f func(t *T)) [bool](/builtin#bool)

Run runs f as a subtest of t called name. It runs f in a separate goroutine and blocks until f returns or calls t.Parallel to become a parallel test. Run reports whether f succeeded (or at least did not fail before calling t.Parallel). 

Run may be called simultaneously from multiple goroutines, but all such calls must return before the outer test function for t returns. 

####  func (*T) [Setenv](https://cs.opensource.google/go/go/+/go1.25.6:src/testing/testing.go;l=1749) ¶ added in go1.17
    
    
    func (t *T) Setenv(key, value [string](/builtin#string))

Setenv calls os.Setenv(key, value) and uses Cleanup to restore the environment variable to its original value after the test. 

Because Setenv affects the whole process, it cannot be used in parallel tests or tests with parallel ancestors. 

####  func (*T) [Skip](https://cs.opensource.google/go/go/+/go1.25.6:src/testing/testing.go;l=1223) ¶ added in go1.1
    
    
    func (c *T) Skip(args ...[any](/builtin#any))

Skip is equivalent to Log followed by SkipNow. 

####  func (*T) [SkipNow](https://cs.opensource.google/go/go/+/go1.25.6:src/testing/testing.go;l=1244) ¶ added in go1.1
    
    
    func (c *T) SkipNow()

SkipNow marks the test as having been skipped and stops its execution by calling [runtime.Goexit](/runtime#Goexit). If a test fails (see Error, Errorf, Fail) and is then skipped, it is still considered to have failed. Execution will continue at the next test or benchmark. See also FailNow. SkipNow must be called from the goroutine running the test, not from other goroutines created during the test. Calling SkipNow does not stop those other goroutines. 

####  func (*T) [Skipf](https://cs.opensource.google/go/go/+/go1.25.6:src/testing/testing.go;l=1230) ¶ added in go1.1
    
    
    func (c *T) Skipf(format [string](/builtin#string), args ...[any](/builtin#any))

Skipf is equivalent to Logf followed by SkipNow. 

####  func (*T) [Skipped](https://cs.opensource.google/go/go/+/go1.25.6:src/testing/testing.go;l=1254) ¶ added in go1.1
    
    
    func (c *T) Skipped() [bool](/builtin#bool)

Skipped reports whether the test was skipped. 

####  func (*T) [TempDir](https://cs.opensource.google/go/go/+/go1.25.6:src/testing/testing.go;l=1321) ¶ added in go1.15
    
    
    func (c *T) TempDir() [string](/builtin#string)

TempDir returns a temporary directory for the test to use. The directory is automatically removed when the test and all its subtests complete. Each subsequent call to TempDir returns a unique directory; if the directory creation fails, TempDir terminates the test by calling Fatal. 

####  type [TB](https://cs.opensource.google/go/go/+/go1.25.6:src/testing/testing.go;l=881) ¶ added in go1.2
    
    
    type TB interface {
    	Attr(key, value [string](/builtin#string))
    	Cleanup(func())
    	Error(args ...[any](/builtin#any))
    	Errorf(format [string](/builtin#string), args ...[any](/builtin#any))
    	Fail()
    	FailNow()
    	Failed() [bool](/builtin#bool)
    	Fatal(args ...[any](/builtin#any))
    	Fatalf(format [string](/builtin#string), args ...[any](/builtin#any))
    	Helper()
    	Log(args ...[any](/builtin#any))
    	Logf(format [string](/builtin#string), args ...[any](/builtin#any))
    	Name() [string](/builtin#string)
    	Setenv(key, value [string](/builtin#string))
    	Chdir(dir [string](/builtin#string))
    	Skip(args ...[any](/builtin#any))
    	SkipNow()
    	Skipf(format [string](/builtin#string), args ...[any](/builtin#any))
    	Skipped() [bool](/builtin#bool)
    	TempDir() [string](/builtin#string)
    	Context() [context](/context).[Context](/context#Context)
    	Output() [io](/io).[Writer](/io#Writer)
    	// contains filtered or unexported methods
    }

TB is the interface common to T, B, and F. 
