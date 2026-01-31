# uber-go/zap

> Source: https://pkg.go.dev/go.uber.org/zap
> Fetched: 2026-01-30T23:49:29.212075+00:00
> Content-Hash: 44d71d5bd7ca74f6
> Type: html

---

Overview

¶

Choosing a Logger

Configuring Zap

Extending Zap

Frequently Asked Questions

Package zap provides fast, structured, leveled logging.

For applications that log in the hot path, reflection-based serialization
and string formatting are prohibitively expensive - they're CPU-intensive
and make many small allocations. Put differently, using json.Marshal and
fmt.Fprintf to log tons of interface{} makes your application slow.

Zap takes a different approach. It includes a reflection-free,
zero-allocation JSON encoder, and the base Logger strives to avoid
serialization overhead and allocations wherever possible. By building the
high-level SugaredLogger on that foundation, zap lets users choose when
they need to count every allocation and when they'd prefer a more familiar,
loosely typed API.

Choosing a Logger

¶

In contexts where performance is nice, but not critical, use the
SugaredLogger. It's 4-10x faster than other structured logging packages and
supports both structured and printf-style logging. Like log15 and go-kit,
the SugaredLogger's structured logging APIs are loosely typed and accept a
variadic number of key-value pairs. (For more advanced use cases, they also
accept strongly typed fields - see the SugaredLogger.With documentation for
details.)

sugar := zap.NewExample().Sugar()
defer sugar.Sync()
sugar.Infow("failed to fetch URL",
  "url", "http://example.com",
  "attempt", 3,
  "backoff", time.Second,
)
sugar.Infof("failed to fetch URL: %s", "http://example.com")

By default, loggers are unbuffered. However, since zap's low-level APIs
allow buffering, calling Sync before letting your process exit is a good
habit.

In the rare contexts where every microsecond and every allocation matter,
use the Logger. It's even faster than the SugaredLogger and allocates far
less, but it only supports strongly-typed, structured logging.

logger := zap.NewExample()
defer logger.Sync()
logger.Info("failed to fetch URL",
  zap.String("url", "http://example.com"),
  zap.Int("attempt", 3),
  zap.Duration("backoff", time.Second),
)

Choosing between the Logger and SugaredLogger doesn't need to be an
application-wide decision: converting between the two is simple and
inexpensive.

logger := zap.NewExample()
defer logger.Sync()
sugar := logger.Sugar()
plain := sugar.Desugar()

Configuring Zap

¶

The simplest way to build a Logger is to use zap's opinionated presets:
NewExample, NewProduction, and NewDevelopment. These presets build a logger
with a single function call:

logger, err := zap.NewProduction()
if err != nil {
  log.Fatalf("can't initialize zap logger: %v", err)
}
defer logger.Sync()

Presets are fine for small projects, but larger projects and organizations
naturally require a bit more customization. For most users, zap's Config
struct strikes the right balance between flexibility and convenience. See
the package-level BasicConfiguration example for sample code.

More unusual configurations (splitting output between files, sending logs
to a message queue, etc.) are possible, but require direct use of
go.uber.org/zap/zapcore. See the package-level AdvancedConfiguration
example for sample code.

Extending Zap

¶

The zap package itself is a relatively thin wrapper around the interfaces
in go.uber.org/zap/zapcore. Extending zap to support a new encoding (e.g.,
BSON), a new log sink (e.g., Kafka), or something more exotic (perhaps an
exception aggregation service, like Sentry or Rollbar) typically requires
implementing the zapcore.Encoder, zapcore.WriteSyncer, or zapcore.Core
interfaces. See the zapcore documentation for details.

Similarly, package authors can use the high-performance Encoder and Core
implementations in the zapcore package to build their own loggers.

Frequently Asked Questions

¶

An FAQ covering everything from installation errors to design decisions is
available at

https://github.com/uber-go/zap/blob/master/FAQ.md

.

Example (AdvancedConfiguration)

¶

package main

import (
	"io"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	// The bundled Config struct only supports the most common configuration
	// options. More complex needs, like splitting logs between multiple files
	// or writing to non-file outputs, require use of the zapcore package.
	//
	// In this example, imagine we're both sending our logs to Kafka and writing
	// them to the console. We'd like to encode the console output and the Kafka
	// topics differently, and we'd also like special treatment for
	// high-priority logs.

	// First, define our level-handling logic.
	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})
	lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.ErrorLevel
	})

	// Assume that we have clients for two Kafka topics. The clients implement
	// zapcore.WriteSyncer and are safe for concurrent use. (If they only
	// implement io.Writer, we can use zapcore.AddSync to add a no-op Sync
	// method. If they're not safe for concurrent use, we can add a protecting
	// mutex with zapcore.Lock.)
	topicDebugging := zapcore.AddSync(io.Discard)
	topicErrors := zapcore.AddSync(io.Discard)

	// High-priority output should also go to standard error, and low-priority
	// output should also go to standard out.
	consoleDebugging := zapcore.Lock(os.Stdout)
	consoleErrors := zapcore.Lock(os.Stderr)

	// Optimize the Kafka output for machine consumption and the console output
	// for human operators.
	kafkaEncoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())

	// Join the outputs, encoders, and level-handling functions into
	// zapcore.Cores, then tee the four cores together.
	core := zapcore.NewTee(
		zapcore.NewCore(kafkaEncoder, topicErrors, highPriority),
		zapcore.NewCore(consoleEncoder, consoleErrors, highPriority),
		zapcore.NewCore(kafkaEncoder, topicDebugging, lowPriority),
		zapcore.NewCore(consoleEncoder, consoleDebugging, lowPriority),
	)

	// From a zapcore.Core, it's easy to construct a Logger.
	logger := zap.New(core)
	defer logger.Sync()
	logger.Info("constructed a logger")
}

Share

Format

Run

Example (BasicConfiguration)

¶

package main

import (
	"encoding/json"

	"go.uber.org/zap"
)

func main() {
	// For some users, the presets offered by the NewProduction, NewDevelopment,
	// and NewExample constructors won't be appropriate. For most of those
	// users, the bundled Config struct offers the right balance of flexibility
	// and convenience. (For more complex needs, see the AdvancedConfiguration
	// example.)
	//
	// See the documentation for Config and zapcore.EncoderConfig for all the
	// available options.
	rawJSON := []byte(`{
	  "level": "debug",
	  "encoding": "json",
	  "outputPaths": ["stdout", "/tmp/logs"],
	  "errorOutputPaths": ["stderr"],
	  "initialFields": {"foo": "bar"},
	  "encoderConfig": {
	    "messageKey": "message",
	    "levelKey": "level",
	    "levelEncoder": "lowercase"
	  }
	}`)

	var cfg zap.Config
	if err := json.Unmarshal(rawJSON, &cfg); err != nil {
		panic(err)
	}
	logger := zap.Must(cfg.Build())
	defer logger.Sync()

	logger.Info("logger construction succeeded")
}

Output:

{"level":"info","message":"logger construction succeeded","foo":"bar"}

Share

Format

Run

Example (Presets)

¶

package main

import (
	"time"

	"go.uber.org/zap"
)

func main() {
	// Using zap's preset constructors is the simplest way to get a feel for the
	// package, but they don't allow much customization.
	logger := zap.NewExample() // or NewProduction, or NewDevelopment
	defer logger.Sync()

	const url = "http://example.com"

	// In most circumstances, use the SugaredLogger. It's 4-10x faster than most
	// other structured logging packages and has a familiar, loosely-typed API.
	sugar := logger.Sugar()
	sugar.Infow("Failed to fetch URL.",
		// Structured context as loosely typed key-value pairs.
		"url", url,
		"attempt", 3,
		"backoff", time.Second,
	)
	sugar.Infof("Failed to fetch URL: %s", url)

	// In the unusual situations where every microsecond matters, use the
	// Logger. It's even faster than the SugaredLogger, but only supports
	// structured logging.
	logger.Info("Failed to fetch URL.",
		// Structured context as strongly typed fields.
		zap.String("url", url),
		zap.Int("attempt", 3),
		zap.Duration("backoff", time.Second),
	)
}

Output:

{"level":"info","msg":"Failed to fetch URL.","url":"http://example.com","attempt":3,"backoff":"1s"}
{"level":"info","msg":"Failed to fetch URL: http://example.com"}
{"level":"info","msg":"Failed to fetch URL.","url":"http://example.com","attempt":3,"backoff":"1s"}

Share

Format

Run

Index

¶

Constants

func CombineWriteSyncers(writers ...zapcore.WriteSyncer) zapcore.WriteSyncer

func DictObject(val ...Field) zapcore.ObjectMarshaler

func LevelFlag(name string, defaultLevel zapcore.Level, usage string) *zapcore.Level

func NewDevelopmentEncoderConfig() zapcore.EncoderConfig

func NewProductionEncoderConfig() zapcore.EncoderConfig

func NewStdLog(l *Logger) *log.Logger

func NewStdLogAt(l *Logger, level zapcore.Level) (*log.Logger, error)

func Open(paths ...string) (zapcore.WriteSyncer, func(), error)

func RedirectStdLog(l *Logger) func()

func RedirectStdLogAt(l *Logger, level zapcore.Level) (func(), error)

func RegisterEncoder(name string, constructor func(zapcore.EncoderConfig) (zapcore.Encoder, error)) error

func RegisterSink(scheme string, factory func(*url.URL) (Sink, error)) error

func ReplaceGlobals(logger *Logger) func()

type AtomicLevel

func NewAtomicLevel() AtomicLevel

func NewAtomicLevelAt(l zapcore.Level) AtomicLevel

func ParseAtomicLevel(text string) (AtomicLevel, error)

func (lvl AtomicLevel) Enabled(l zapcore.Level) bool

func (lvl AtomicLevel) Level() zapcore.Level

func (lvl AtomicLevel) MarshalText() (text []byte, err error)

func (lvl AtomicLevel) ServeHTTP(w http.ResponseWriter, r *http.Request)

func (lvl AtomicLevel) SetLevel(l zapcore.Level)

func (lvl AtomicLevel) String() string

func (lvl *AtomicLevel) UnmarshalText(text []byte) error

type Config

func NewDevelopmentConfig() Config

func NewProductionConfig() Config

func (cfg Config) Build(opts ...Option) (*Logger, error)

type Field

func Any(key string, value interface{}) Field

func Array(key string, val zapcore.ArrayMarshaler) Field

func Binary(key string, val []byte) Field

func Bool(key string, val bool) Field

func Boolp(key string, val *bool) Field

func Bools(key string, bs []bool) Field

func ByteString(key string, val []byte) Field

func ByteStrings(key string, bss [][]byte) Field

func Complex128(key string, val complex128) Field

func Complex128p(key string, val *complex128) Field

func Complex128s(key string, nums []complex128) Field

func Complex64(key string, val complex64) Field

func Complex64p(key string, val *complex64) Field

func Complex64s(key string, nums []complex64) Field

func Dict(key string, val ...Field) Field

func Duration(key string, val time.Duration) Field

func Durationp(key string, val *time.Duration) Field

func Durations(key string, ds []time.Duration) Field

func Error(err error) Field

func Errors(key string, errs []error) Field

func Float32(key string, val float32) Field

func Float32p(key string, val *float32) Field

func Float32s(key string, nums []float32) Field

func Float64(key string, val float64) Field

func Float64p(key string, val *float64) Field

func Float64s(key string, nums []float64) Field

func Inline(val zapcore.ObjectMarshaler) Field

func Int(key string, val int) Field

func Int16(key string, val int16) Field

func Int16p(key string, val *int16) Field

func Int16s(key string, nums []int16) Field

func Int32(key string, val int32) Field

func Int32p(key string, val *int32) Field

func Int32s(key string, nums []int32) Field

func Int64(key string, val int64) Field

func Int64p(key string, val *int64) Field

func Int64s(key string, nums []int64) Field

func Int8(key string, val int8) Field

func Int8p(key string, val *int8) Field

func Int8s(key string, nums []int8) Field

func Intp(key string, val *int) Field

func Ints(key string, nums []int) Field

func NamedError(key string, err error) Field

func Namespace(key string) Field

func Object(key string, val zapcore.ObjectMarshaler) Field

func ObjectValues[T any, P ObjectMarshalerPtr[T]](key string, values []T) Field

func Objects[T zapcore.ObjectMarshaler](key string, values []T) Field

func Reflect(key string, val interface{}) Field

func Skip() Field

func Stack(key string) Field

func StackSkip(key string, skip int) Field

func String(key string, val string) Field

func Stringer(key string, val fmt.Stringer) Field

func Stringers[T fmt.Stringer](key string, values []T) Field

func Stringp(key string, val *string) Field

func Strings(key string, ss []string) Field

func Time(key string, val time.Time) Field

func Timep(key string, val *time.Time) Field

func Times(key string, ts []time.Time) Field

func Uint(key string, val uint) Field

func Uint16(key string, val uint16) Field

func Uint16p(key string, val *uint16) Field

func Uint16s(key string, nums []uint16) Field

func Uint32(key string, val uint32) Field

func Uint32p(key string, val *uint32) Field

func Uint32s(key string, nums []uint32) Field

func Uint64(key string, val uint64) Field

func Uint64p(key string, val *uint64) Field

func Uint64s(key string, nums []uint64) Field

func Uint8(key string, val uint8) Field

func Uint8p(key string, val *uint8) Field

func Uint8s(key string, nums []uint8) Field

func Uintp(key string, val *uint) Field

func Uintptr(key string, val uintptr) Field

func Uintptrp(key string, val *uintptr) Field

func Uintptrs(key string, us []uintptr) Field

func Uints(key string, nums []uint) Field

type LevelEnablerFunc

func (f LevelEnablerFunc) Enabled(lvl zapcore.Level) bool

type Logger

func L() *Logger

func Must(logger *Logger, err error) *Logger

func New(core zapcore.Core, options ...Option) *Logger

func NewDevelopment(options ...Option) (*Logger, error)

func NewExample(options ...Option) *Logger

func NewNop() *Logger

func NewProduction(options ...Option) (*Logger, error)

func (log *Logger) Check(lvl zapcore.Level, msg string) *zapcore.CheckedEntry

func (log *Logger) Core() zapcore.Core

func (log *Logger) DPanic(msg string, fields ...Field)

func (log *Logger) Debug(msg string, fields ...Field)

func (log *Logger) Error(msg string, fields ...Field)

func (log *Logger) Fatal(msg string, fields ...Field)

func (log *Logger) Info(msg string, fields ...Field)

func (log *Logger) Level() zapcore.Level

func (log *Logger) Log(lvl zapcore.Level, msg string, fields ...Field)

func (log *Logger) Name() string

func (log *Logger) Named(s string) *Logger

func (log *Logger) Panic(msg string, fields ...Field)

func (log *Logger) Sugar() *SugaredLogger

func (log *Logger) Sync() error

func (log *Logger) Warn(msg string, fields ...Field)

func (log *Logger) With(fields ...Field) *Logger

func (log *Logger) WithLazy(fields ...Field) *Logger

func (log *Logger) WithOptions(opts ...Option) *Logger

type ObjectMarshalerPtr

type Option

func AddCaller() Option

func AddCallerSkip(skip int) Option

func AddStacktrace(lvl zapcore.LevelEnabler) Option

func Development() Option

func ErrorOutput(w zapcore.WriteSyncer) Option

func Fields(fs ...Field) Option

func Hooks(hooks ...func(zapcore.Entry) error) Option

func IncreaseLevel(lvl zapcore.LevelEnabler) Option

func OnFatal(action zapcore.CheckWriteAction) Option

deprecated

func WithCaller(enabled bool) Option

func WithClock(clock zapcore.Clock) Option

func WithFatalHook(hook zapcore.CheckWriteHook) Option

func WithPanicHook(hook zapcore.CheckWriteHook) Option

func WrapCore(f func(zapcore.Core) zapcore.Core) Option

type SamplingConfig

type Sink

type SugaredLogger

func S() *SugaredLogger

func (s *SugaredLogger) DPanic(args ...interface{})

func (s *SugaredLogger) DPanicf(template string, args ...interface{})

func (s *SugaredLogger) DPanicln(args ...interface{})

func (s *SugaredLogger) DPanicw(msg string, keysAndValues ...interface{})

func (s *SugaredLogger) Debug(args ...interface{})

func (s *SugaredLogger) Debugf(template string, args ...interface{})

func (s *SugaredLogger) Debugln(args ...interface{})

func (s *SugaredLogger) Debugw(msg string, keysAndValues ...interface{})

func (s *SugaredLogger) Desugar() *Logger

func (s *SugaredLogger) Error(args ...interface{})

func (s *SugaredLogger) Errorf(template string, args ...interface{})

func (s *SugaredLogger) Errorln(args ...interface{})

func (s *SugaredLogger) Errorw(msg string, keysAndValues ...interface{})

func (s *SugaredLogger) Fatal(args ...interface{})

func (s *SugaredLogger) Fatalf(template string, args ...interface{})

func (s *SugaredLogger) Fatalln(args ...interface{})

func (s *SugaredLogger) Fatalw(msg string, keysAndValues ...interface{})

func (s *SugaredLogger) Info(args ...interface{})

func (s *SugaredLogger) Infof(template string, args ...interface{})

func (s *SugaredLogger) Infoln(args ...interface{})

func (s *SugaredLogger) Infow(msg string, keysAndValues ...interface{})

func (s *SugaredLogger) Level() zapcore.Level

func (s *SugaredLogger) Log(lvl zapcore.Level, args ...interface{})

func (s *SugaredLogger) Logf(lvl zapcore.Level, template string, args ...interface{})

func (s *SugaredLogger) Logln(lvl zapcore.Level, args ...interface{})

func (s *SugaredLogger) Logw(lvl zapcore.Level, msg string, keysAndValues ...interface{})

func (s *SugaredLogger) Named(name string) *SugaredLogger

func (s *SugaredLogger) Panic(args ...interface{})

func (s *SugaredLogger) Panicf(template string, args ...interface{})

func (s *SugaredLogger) Panicln(args ...interface{})

func (s *SugaredLogger) Panicw(msg string, keysAndValues ...interface{})

func (s *SugaredLogger) Sync() error

func (s *SugaredLogger) Warn(args ...interface{})

func (s *SugaredLogger) Warnf(template string, args ...interface{})

func (s *SugaredLogger) Warnln(args ...interface{})

func (s *SugaredLogger) Warnw(msg string, keysAndValues ...interface{})

func (s *SugaredLogger) With(args ...interface{}) *SugaredLogger

func (s *SugaredLogger) WithLazy(args ...interface{}) *SugaredLogger

func (s *SugaredLogger) WithOptions(opts ...Option) *SugaredLogger

Examples

¶

Package (AdvancedConfiguration)

Package (BasicConfiguration)

Package (Presets)

AtomicLevel

AtomicLevel (Config)

Dict

DictObject

Logger.Check

Logger.Named

Namespace

NewStdLog

Object

ObjectValues

Objects

RedirectStdLog

ReplaceGlobals

WrapCore (Replace)

WrapCore (Wrap)

Constants

¶

View Source

const (

// DebugLevel logs are typically voluminous, and are usually disabled in

// production.

DebugLevel =

zapcore

.

DebugLevel

// InfoLevel is the default logging priority.

InfoLevel =

zapcore

.

InfoLevel

// WarnLevel logs are more important than Info, but don't need individual

// human review.

WarnLevel =

zapcore

.

WarnLevel

// ErrorLevel logs are high-priority. If an application is running smoothly,

// it shouldn't generate any error-level logs.

ErrorLevel =

zapcore

.

ErrorLevel

// DPanicLevel logs are particularly important errors. In development the

// logger panics after writing the message.

DPanicLevel =

zapcore

.

DPanicLevel

// PanicLevel logs a message, then panics.

PanicLevel =

zapcore

.

PanicLevel

// FatalLevel logs a message, then calls os.Exit(1).

FatalLevel =

zapcore

.

FatalLevel

)

Variables

¶

This section is empty.

Functions

¶

func

CombineWriteSyncers

¶

func CombineWriteSyncers(writers ...

zapcore

.

WriteSyncer

)

zapcore

.

WriteSyncer

CombineWriteSyncers is a utility that combines multiple WriteSyncers into a
single, locked WriteSyncer. If no inputs are supplied, it returns a no-op
WriteSyncer.

It's provided purely as a convenience; the result is no different from
using zapcore.NewMultiWriteSyncer and zapcore.Lock individually.

func

DictObject

¶

added in

v1.27.1

func DictObject(val ...

Field

)

zapcore

.

ObjectMarshaler

DictObject constructs a

zapcore.ObjectMarshaler

with the given list of fields.
The resulting object marshaler can be used as input to

Object

,

Objects

, or
any other functions that expect an object marshaler.

Example

¶

package main

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	logger := zap.NewExample()
	defer logger.Sync()

	// Use DictObject to create zapcore.ObjectMarshaler implementations from Field arrays,
	// then use the Object and Objects field constructors to turn them back into a Field.

	logger.Debug("worker received job",
		zap.Object("w1",
			zap.DictObject(
				zap.Int("id", 402000),
				zap.String("description", "compress image data"),
				zap.Int("priority", 3),
			),
		))

	d1 := 68 * time.Millisecond
	d2 := 79 * time.Millisecond
	d3 := 57 * time.Millisecond

	logger.Info("worker status checks",
		zap.Objects("job batch enqueued",
			[]zapcore.ObjectMarshaler{
				zap.DictObject(
					zap.String("worker", "w1"),
					zap.Int("load", 419),
					zap.Duration("latency", d1),
				),
				zap.DictObject(
					zap.String("worker", "w2"),
					zap.Int("load", 520),
					zap.Duration("latency", d2),
				),
				zap.DictObject(
					zap.String("worker", "w3"),
					zap.Int("load", 310),
					zap.Duration("latency", d3),
				),
			},
		))
}

Output:

{"level":"debug","msg":"worker received job","w1":{"id":402000,"description":"compress image data","priority":3}}
{"level":"info","msg":"worker status checks","job batch enqueued":[{"worker":"w1","load":419,"latency":"68ms"},{"worker":"w2","load":520,"latency":"79ms"},{"worker":"w3","load":310,"latency":"57ms"}]}

Share

Format

Run

func

LevelFlag

¶

func LevelFlag(name

string

, defaultLevel

zapcore

.

Level

, usage

string

) *

zapcore

.

Level

LevelFlag uses the standard library's flag.Var to declare a global flag
with the specified name, default, and usage guidance. The returned value is
a pointer to the value of the flag.

If you don't want to use the flag package's global state, you can use any
non-nil *Level as a flag.Value with your own *flag.FlagSet.

func

NewDevelopmentEncoderConfig

¶

func NewDevelopmentEncoderConfig()

zapcore

.

EncoderConfig

NewDevelopmentEncoderConfig returns an opinionated EncoderConfig for
development environments.

Messages encoded with this configuration will use Zap's console encoder
intended to print human-readable output.
It will print log messages with the following information:

The log level (e.g. "INFO", "ERROR").

The time in ISO8601 format (e.g. "2017-01-01T12:00:00Z").

The message passed to the log statement.

If available, a short path to the file and line number
where the log statement was issued.
The logger configuration determines whether this field is captured.

If available, a stacktrace from the line
where the log statement was issued.
The logger configuration determines whether this field is captured.

By default, the following formats are used for different types:

Time is formatted in ISO8601 format (e.g. "2017-01-01T12:00:00Z").

Duration is formatted as a string (e.g. "1.234s").

You may change these by setting the appropriate fields in the returned
object.
For example, use the following to change the time encoding format:

cfg := zap.NewDevelopmentEncoderConfig()
cfg.EncodeTime = zapcore.ISO8601TimeEncoder

func

NewProductionEncoderConfig

¶

func NewProductionEncoderConfig()

zapcore

.

EncoderConfig

NewProductionEncoderConfig returns an opinionated EncoderConfig for
production environments.

Messages encoded with this configuration will be JSON-formatted
and will have the following keys by default:

"level": The logging level (e.g. "info", "error").

"ts": The current time in number of seconds since the Unix epoch.

"msg": The message passed to the log statement.

"caller": If available, a short path to the file and line number
where the log statement was issued.
The logger configuration determines whether this field is captured.

"stacktrace": If available, a stack trace from the line
where the log statement was issued.
The logger configuration determines whether this field is captured.

By default, the following formats are used for different types:

Time is formatted as floating-point number of seconds since the Unix
epoch.

Duration is formatted as floating-point number of seconds.

You may change these by setting the appropriate fields in the returned
object.
For example, use the following to change the time encoding format:

cfg := zap.NewProductionEncoderConfig()
cfg.EncodeTime = zapcore.ISO8601TimeEncoder

func

NewStdLog

¶

func NewStdLog(l *

Logger

) *

log

.

Logger

NewStdLog returns a *log.Logger which writes to the supplied zap Logger at
InfoLevel. To redirect the standard library's package-global logging
functions, use RedirectStdLog instead.

Example

¶

package main

import (
	"go.uber.org/zap"
)

func main() {
	logger := zap.NewExample()
	defer logger.Sync()

	std := zap.NewStdLog(logger)
	std.Print("standard logger wrapper")
}

Output:

{"level":"info","msg":"standard logger wrapper"}

Share

Format

Run

func

NewStdLogAt

¶

added in

v1.7.0

func NewStdLogAt(l *

Logger

, level

zapcore

.

Level

) (*

log

.

Logger

,

error

)

NewStdLogAt returns *log.Logger which writes to supplied zap logger at
required level.

func

Open

¶

func Open(paths ...

string

) (

zapcore

.

WriteSyncer

, func(),

error

)

Open is a high-level wrapper that takes a variadic number of URLs, opens or
creates each of the specified resources, and combines them into a locked
WriteSyncer. It also returns any error encountered and a function to close
any opened files.

Passing no URLs returns a no-op WriteSyncer. Zap handles URLs without a
scheme and URLs with the "file" scheme. Third-party code may register
factories for other schemes using RegisterSink.

URLs with the "file" scheme must use absolute paths on the local
filesystem. No user, password, port, fragments, or query parameters are
allowed, and the hostname must be empty or "localhost".

Since it's common to write logs to the local filesystem, URLs without a
scheme (e.g., "/var/log/foo.log") are treated as local file paths. Without
a scheme, the special paths "stdout" and "stderr" are interpreted as
os.Stdout and os.Stderr. When specified without a scheme, relative file
paths also work.

func

RedirectStdLog

¶

func RedirectStdLog(l *

Logger

) func()

RedirectStdLog redirects output from the standard library's package-global
logger to the supplied logger at InfoLevel. Since zap already handles caller
annotations, timestamps, etc., it automatically disables the standard
library's annotations and prefixing.

It returns a function to restore the original prefix and flags and reset the
standard library's output to os.Stderr.

Example

¶

package main

import (
	"log"

	"go.uber.org/zap"
)

func main() {
	logger := zap.NewExample()
	defer logger.Sync()

	undo := zap.RedirectStdLog(logger)
	defer undo()

	log.Print("redirected standard library")
}

Output:

{"level":"info","msg":"redirected standard library"}

Share

Format

Run

func

RedirectStdLogAt

¶

added in

v1.8.0

func RedirectStdLogAt(l *

Logger

, level

zapcore

.

Level

) (func(),

error

)

RedirectStdLogAt redirects output from the standard library's package-global
logger to the supplied logger at the specified level. Since zap already
handles caller annotations, timestamps, etc., it automatically disables the
standard library's annotations and prefixing.

It returns a function to restore the original prefix and flags and reset the
standard library's output to os.Stderr.

func

RegisterEncoder

¶

func RegisterEncoder(name

string

, constructor func(

zapcore

.

EncoderConfig

) (

zapcore

.

Encoder

,

error

))

error

RegisterEncoder registers an encoder constructor, which the Config struct
can then reference. By default, the "json" and "console" encoders are
registered.

Attempting to register an encoder whose name is already taken returns an
error.

func

RegisterSink

¶

added in

v1.9.0

func RegisterSink(scheme

string

, factory func(*

url

.

URL

) (

Sink

,

error

))

error

RegisterSink registers a user-supplied factory for all sinks with a
particular scheme.

All schemes must be ASCII, valid under section 0.1 of

RFC 3986

(

https://tools.ietf.org/html/rfc3983#section-3.1

), and must not already
have a factory registered. Zap automatically registers a factory for the
"file" scheme.

func

ReplaceGlobals

¶

func ReplaceGlobals(logger *

Logger

) func()

ReplaceGlobals replaces the global Logger and SugaredLogger, and returns a
function to restore the original values. It's safe for concurrent use.

Example

¶

package main

import (
	"go.uber.org/zap"
)

func main() {
	logger := zap.NewExample()
	defer logger.Sync()

	undo := zap.ReplaceGlobals(logger)
	defer undo()

	zap.L().Info("replaced zap's global loggers")
}

Output:

{"level":"info","msg":"replaced zap's global loggers"}

Share

Format

Run

Types

¶

type

AtomicLevel

¶

type AtomicLevel struct {

// contains filtered or unexported fields

}

An AtomicLevel is an atomically changeable, dynamic logging level. It lets
you safely change the log level of a tree of loggers (the root logger and
any children created by adding context) at runtime.

The AtomicLevel itself is an http.Handler that serves a JSON endpoint to
alter its level.

AtomicLevels must be created with the NewAtomicLevel constructor to allocate
their internal atomic pointer.

Example

¶

package main

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	atom := zap.NewAtomicLevel()

	// To keep the example deterministic, disable timestamps in the output.
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = ""

	logger := zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		zapcore.Lock(os.Stdout),
		atom,
	))
	defer logger.Sync()

	logger.Info("info logging enabled")

	atom.SetLevel(zap.ErrorLevel)
	logger.Info("info logging disabled")
}

Output:

{"level":"info","msg":"info logging enabled"}

Share

Format

Run

Example (Config)

¶

package main

import (
	"encoding/json"

	"go.uber.org/zap"
)

func main() {
	// The zap.Config struct includes an AtomicLevel. To use it, keep a
	// reference to the Config.
	rawJSON := []byte(`{
		"level": "info",
		"outputPaths": ["stdout"],
		"errorOutputPaths": ["stderr"],
		"encoding": "json",
		"encoderConfig": {
			"messageKey": "message",
			"levelKey": "level",
			"levelEncoder": "lowercase"
		}
	}`)
	var cfg zap.Config
	if err := json.Unmarshal(rawJSON, &cfg); err != nil {
		panic(err)
	}
	logger := zap.Must(cfg.Build())
	defer logger.Sync()

	logger.Info("info logging enabled")

	cfg.Level.SetLevel(zap.ErrorLevel)
	logger.Info("info logging disabled")
}

Output:

{"level":"info","message":"info logging enabled"}

Share

Format

Run

func

NewAtomicLevel

¶

func NewAtomicLevel()

AtomicLevel

NewAtomicLevel creates an AtomicLevel with InfoLevel and above logging
enabled.

func

NewAtomicLevelAt

¶

added in

v1.3.0

func NewAtomicLevelAt(l

zapcore

.

Level

)

AtomicLevel

NewAtomicLevelAt is a convenience function that creates an AtomicLevel
and then calls SetLevel with the given level.

func

ParseAtomicLevel

¶

added in

v1.21.0

func ParseAtomicLevel(text

string

) (

AtomicLevel

,

error

)

ParseAtomicLevel parses an AtomicLevel based on a lowercase or all-caps ASCII
representation of the log level. If the provided ASCII representation is
invalid an error is returned.

This is particularly useful when dealing with text input to configure log
levels.

func (AtomicLevel)

Enabled

¶

func (lvl

AtomicLevel

) Enabled(l

zapcore

.

Level

)

bool

Enabled implements the zapcore.LevelEnabler interface, which allows the
AtomicLevel to be used in place of traditional static levels.

func (AtomicLevel)

Level

¶

func (lvl

AtomicLevel

) Level()

zapcore

.

Level

Level returns the minimum enabled log level.

func (AtomicLevel)

MarshalText

¶

added in

v1.3.0

func (lvl

AtomicLevel

) MarshalText() (text []

byte

, err

error

)

MarshalText marshals the AtomicLevel to a byte slice. It uses the same
text representation as the static zapcore.Levels ("debug", "info", "warn",
"error", "dpanic", "panic", and "fatal").

func (AtomicLevel)

ServeHTTP

¶

func (lvl

AtomicLevel

) ServeHTTP(w

http

.

ResponseWriter

, r *

http

.

Request

)

GET

PUT

ServeHTTP is a simple JSON endpoint that can report on or change the current
logging level.

GET

¶

The GET request returns a JSON description of the current logging level like:

{"level":"info"}

PUT

¶

The PUT request changes the logging level. It is perfectly safe to change the
logging level while a program is running. Two content types are supported:

Content-Type: application/x-www-form-urlencoded

With this content type, the level can be provided through the request body or
a query parameter. The log level is URL encoded like:

level=debug

The request body takes precedence over the query parameter, if both are
specified.

This content type is the default for a curl PUT request. Following are two
example curl requests that both set the logging level to debug.

curl -X PUT localhost:8080/log/level?level=debug
curl -X PUT localhost:8080/log/level -d level=debug

For any other content type, the payload is expected to be JSON encoded and
look like:

{"level":"info"}

An example curl request could look like this:

curl -X PUT localhost:8080/log/level -H "Content-Type: application/json" -d '{"level":"debug"}'

func (AtomicLevel)

SetLevel

¶

func (lvl

AtomicLevel

) SetLevel(l

zapcore

.

Level

)

SetLevel alters the logging level.

func (AtomicLevel)

String

¶

added in

v1.4.0

func (lvl

AtomicLevel

) String()

string

String returns the string representation of the underlying Level.

func (*AtomicLevel)

UnmarshalText

¶

func (lvl *

AtomicLevel

) UnmarshalText(text []

byte

)

error

UnmarshalText unmarshals the text to an AtomicLevel. It uses the same text
representations as the static zapcore.Levels ("debug", "info", "warn",
"error", "dpanic", "panic", and "fatal").

type

Config

¶

type Config struct {

// Level is the minimum enabled logging level. Note that this is a dynamic

// level, so calling Config.Level.SetLevel will atomically change the log

// level of all loggers descended from this config.

Level

AtomicLevel

`json:"level" yaml:"level"`

// Development puts the logger in development mode, which changes the

// behavior of DPanicLevel and takes stacktraces more liberally.

Development

bool

`json:"development" yaml:"development"`

// DisableCaller stops annotating logs with the calling function's file

// name and line number. By default, all logs are annotated.

DisableCaller

bool

`json:"disableCaller" yaml:"disableCaller"`

// DisableStacktrace completely disables automatic stacktrace capturing. By

// default, stacktraces are captured for WarnLevel and above logs in

// development and ErrorLevel and above in production.

DisableStacktrace

bool

`json:"disableStacktrace" yaml:"disableStacktrace"`

// Sampling sets a sampling policy. A nil SamplingConfig disables sampling.

Sampling *

SamplingConfig

`json:"sampling" yaml:"sampling"`

// Encoding sets the logger's encoding. Valid values are "json" and

// "console", as well as any third-party encodings registered via

// RegisterEncoder.

Encoding

string

`json:"encoding" yaml:"encoding"`

// EncoderConfig sets options for the chosen encoder. See

// zapcore.EncoderConfig for details.

EncoderConfig

zapcore

.

EncoderConfig

`json:"encoderConfig" yaml:"encoderConfig"`

// OutputPaths is a list of URLs or file paths to write logging output to.

// See Open for details.

OutputPaths []

string

`json:"outputPaths" yaml:"outputPaths"`

// ErrorOutputPaths is a list of URLs to write internal logger errors to.

// The default is standard error.

//

// Note that this setting only affects internal errors; for sample code that

// sends error-level logs to a different location from info- and debug-level

// logs, see the package-level AdvancedConfiguration example.

ErrorOutputPaths []

string

`json:"errorOutputPaths" yaml:"errorOutputPaths"`

// InitialFields is a collection of fields to add to the root logger.

InitialFields map[

string

]interface{} `json:"initialFields" yaml:"initialFields"`
}

Config offers a declarative way to construct a logger. It doesn't do
anything that can't be done with New, Options, and the various
zapcore.WriteSyncer and zapcore.Core wrappers, but it's a simpler way to
toggle common options.

Note that Config intentionally supports only the most common options. More
unusual logging setups (logging to network connections or message queues,
splitting output between multiple files, etc.) are possible, but require
direct use of the zapcore package. For sample code, see the package-level
BasicConfiguration and AdvancedConfiguration examples.

For an example showing runtime log level changes, see the documentation for
AtomicLevel.

func

NewDevelopmentConfig

¶

func NewDevelopmentConfig()

Config

NewDevelopmentConfig builds a reasonable default development logging
configuration.
Logging is enabled at DebugLevel and above, and uses a console encoder.
Logs are written to standard error.
Stacktraces are included on logs of WarnLevel and above.
DPanicLevel logs will panic.

See

NewDevelopmentEncoderConfig

for information
on the default encoder configuration.

func

NewProductionConfig

¶

func NewProductionConfig()

Config

NewProductionConfig builds a reasonable default production logging
configuration.
Logging is enabled at InfoLevel and above, and uses a JSON encoder.
Logs are written to standard error.
Stacktraces are included on logs of ErrorLevel and above.
DPanicLevel logs will not panic, but will write a stacktrace.

Sampling is enabled at 100:100 by default,
meaning that after the first 100 log entries
with the same level and message in the same second,
it will log every 100th entry
with the same level and message in the same second.
You may disable this behavior by setting Sampling to nil.

See

NewProductionEncoderConfig

for information
on the default encoder configuration.

func (Config)

Build

¶

func (cfg

Config

) Build(opts ...

Option

) (*

Logger

,

error

)

Build constructs a logger from the Config and Options.

type

Field

¶

added in

v1.8.0

type Field =

zapcore

.

Field

Field is an alias for Field. Aliasing this type dramatically
improves the navigability of this package's API documentation.

func

Any

¶

func Any(key

string

, value interface{})

Field

Any takes a key and an arbitrary value and chooses the best way to represent
them as a field, falling back to a reflection-based approach only if
necessary.

Since byte/uint8 and rune/int32 are aliases, Any can't differentiate between
them. To minimize surprises, []byte values are treated as binary blobs, byte
values are treated as uint8, and runes are always treated as integers.

func

Array

¶

func Array(key

string

, val

zapcore

.

ArrayMarshaler

)

Field

Array constructs a field with the given key and ArrayMarshaler. It provides
a flexible, but still type-safe and efficient, way to add array-like types
to the logging context. The struct's MarshalLogArray method is called lazily.

func

Binary

¶

func Binary(key

string

, val []

byte

)

Field

Binary constructs a field that carries an opaque binary blob.

Binary data is serialized in an encoding-appropriate format. For example,
zap's JSON encoder base64-encodes binary blobs. To log UTF-8 encoded text,
use ByteString.

func

Bool

¶

func Bool(key

string

, val

bool

)

Field

Bool constructs a field that carries a bool.

func

Boolp

¶

added in

v1.13.0

func Boolp(key

string

, val *

bool

)

Field

Boolp constructs a field that carries a *bool. The returned Field will safely
and explicitly represent `nil` when appropriate.

func

Bools

¶

func Bools(key

string

, bs []

bool

)

Field

Bools constructs a field that carries a slice of bools.

func

ByteString

¶

func ByteString(key

string

, val []

byte

)

Field

ByteString constructs a field that carries UTF-8 encoded text as a []byte.
To log opaque binary blobs (which aren't necessarily valid UTF-8), use
Binary.

func

ByteStrings

¶

func ByteStrings(key

string

, bss [][]

byte

)

Field

ByteStrings constructs a field that carries a slice of []byte, each of which
must be UTF-8 encoded text.

func

Complex128

¶

func Complex128(key

string

, val

complex128

)

Field

Complex128 constructs a field that carries a complex number. Unlike most
numeric fields, this costs an allocation (to convert the complex128 to
interface{}).

func

Complex128p

¶

added in

v1.13.0

func Complex128p(key

string

, val *

complex128

)

Field

Complex128p constructs a field that carries a *complex128. The returned Field will safely
and explicitly represent `nil` when appropriate.

func

Complex128s

¶

func Complex128s(key

string

, nums []

complex128

)

Field

Complex128s constructs a field that carries a slice of complex numbers.

func

Complex64

¶

func Complex64(key

string

, val

complex64

)

Field

Complex64 constructs a field that carries a complex number. Unlike most
numeric fields, this costs an allocation (to convert the complex64 to
interface{}).

func

Complex64p

¶

added in

v1.13.0

func Complex64p(key

string

, val *

complex64

)

Field

Complex64p constructs a field that carries a *complex64. The returned Field will safely
and explicitly represent `nil` when appropriate.

func

Complex64s

¶

func Complex64s(key

string

, nums []

complex64

)

Field

Complex64s constructs a field that carries a slice of complex numbers.

func

Dict

¶

added in

v1.26.0

func Dict(key

string

, val ...

Field

)

Field

Dict constructs a field containing the provided key-value pairs.
It acts similar to

Object

, but with the fields specified as arguments.

Example

¶

package main

import (
	"go.uber.org/zap"
)

func main() {
	logger := zap.NewExample()
	defer logger.Sync()

	logger.Info("login event",
		zap.Dict("event",
			zap.Int("id", 123),
			zap.String("name", "jane"),
			zap.String("status", "pending")))
}

Output:

{"level":"info","msg":"login event","event":{"id":123,"name":"jane","status":"pending"}}

Share

Format

Run

func

Duration

¶

func Duration(key

string

, val

time

.

Duration

)

Field

Duration constructs a field with the given key and value. The encoder
controls how the duration is serialized.

func

Durationp

¶

added in

v1.13.0

func Durationp(key

string

, val *

time

.

Duration

)

Field

Durationp constructs a field that carries a *time.Duration. The returned Field will safely
and explicitly represent `nil` when appropriate.

func

Durations

¶

func Durations(key

string

, ds []

time

.

Duration

)

Field

Durations constructs a field that carries a slice of time.Durations.

func

Error

¶

func Error(err

error

)

Field

Error is shorthand for the common idiom NamedError("error", err).

func

Errors

¶

func Errors(key

string

, errs []

error

)

Field

Errors constructs a field that carries a slice of errors.

func

Float32

¶

func Float32(key

string

, val

float32

)

Field

Float32 constructs a field that carries a float32. The way the
floating-point value is represented is encoder-dependent, so marshaling is
necessarily lazy.

func

Float32p

¶

added in

v1.13.0

func Float32p(key

string

, val *

float32

)

Field

Float32p constructs a field that carries a *float32. The returned Field will safely
and explicitly represent `nil` when appropriate.

func

Float32s

¶

func Float32s(key

string

, nums []

float32

)

Field

Float32s constructs a field that carries a slice of floats.

func

Float64

¶

func Float64(key

string

, val

float64

)

Field

Float64 constructs a field that carries a float64. The way the
floating-point value is represented is encoder-dependent, so marshaling is
necessarily lazy.

func

Float64p

¶

added in

v1.13.0

func Float64p(key

string

, val *

float64

)

Field

Float64p constructs a field that carries a *float64. The returned Field will safely
and explicitly represent `nil` when appropriate.

func

Float64s

¶

func Float64s(key

string

, nums []

float64

)

Field

Float64s constructs a field that carries a slice of floats.

func

Inline

¶

added in

v1.17.0

func Inline(val

zapcore

.

ObjectMarshaler

)

Field

Inline constructs a Field that is similar to Object, but it
will add the elements of the provided ObjectMarshaler to the
current namespace.

func

Int

¶

func Int(key

string

, val

int

)

Field

Int constructs a field with the given key and value.

func

Int16

¶

func Int16(key

string

, val

int16

)

Field

Int16 constructs a field with the given key and value.

func

Int16p

¶

added in

v1.13.0

func Int16p(key

string

, val *

int16

)

Field

Int16p constructs a field that carries a *int16. The returned Field will safely
and explicitly represent `nil` when appropriate.

func

Int16s

¶

func Int16s(key

string

, nums []

int16

)

Field

Int16s constructs a field that carries a slice of integers.

func

Int32

¶

func Int32(key

string

, val

int32

)

Field

Int32 constructs a field with the given key and value.

func

Int32p

¶

added in

v1.13.0

func Int32p(key

string

, val *

int32

)

Field

Int32p constructs a field that carries a *int32. The returned Field will safely
and explicitly represent `nil` when appropriate.

func

Int32s

¶

func Int32s(key

string

, nums []

int32

)

Field

Int32s constructs a field that carries a slice of integers.

func

Int64

¶

func Int64(key

string

, val

int64

)

Field

Int64 constructs a field with the given key and value.

func

Int64p

¶

added in

v1.13.0

func Int64p(key

string

, val *

int64

)

Field

Int64p constructs a field that carries a *int64. The returned Field will safely
and explicitly represent `nil` when appropriate.

func

Int64s

¶

func Int64s(key

string

, nums []

int64

)

Field

Int64s constructs a field that carries a slice of integers.

func

Int8

¶

func Int8(key

string

, val

int8

)

Field

Int8 constructs a field with the given key and value.

func

Int8p

¶

added in

v1.13.0

func Int8p(key

string

, val *

int8

)

Field

Int8p constructs a field that carries a *int8. The returned Field will safely
and explicitly represent `nil` when appropriate.

func

Int8s

¶

func Int8s(key

string

, nums []

int8

)

Field

Int8s constructs a field that carries a slice of integers.

func

Intp

¶

added in

v1.13.0

func Intp(key

string

, val *

int

)

Field

Intp constructs a field that carries a *int. The returned Field will safely
and explicitly represent `nil` when appropriate.

func

Ints

¶

func Ints(key

string

, nums []

int

)

Field

Ints constructs a field that carries a slice of integers.

func

NamedError

¶

func NamedError(key

string

, err

error

)

Field

NamedError constructs a field that lazily stores err.Error() under the
provided key. Errors which also implement fmt.Formatter (like those produced
by github.com/pkg/errors) will also have their verbose representation stored
under key+"Verbose". If passed a nil error, the field is a no-op.

For the common case in which the key is simply "error", the Error function
is shorter and less repetitive.

func

Namespace

¶

func Namespace(key

string

)

Field

Namespace creates a named, isolated scope within the logger's context. All
subsequent fields will be added to the new namespace.

This helps prevent key collisions when injecting loggers into sub-components
or third-party libraries.

Example

¶

package main

import (
	"go.uber.org/zap"
)

func main() {
	logger := zap.NewExample()
	defer logger.Sync()

	logger.With(
		zap.Namespace("metrics"),
		zap.Int("counter", 1),
	).Info("tracked some metrics")
}

Output:

{"level":"info","msg":"tracked some metrics","metrics":{"counter":1}}

Share

Format

Run

func

Object

¶

func Object(key

string

, val

zapcore

.

ObjectMarshaler

)

Field

Object constructs a field with the given key and ObjectMarshaler. It
provides a flexible, but still type-safe and efficient, way to add map- or
struct-like user-defined types to the logging context. The struct's
MarshalLogObject method is called lazily.

Example

¶

package main

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type addr struct {
	IP   string
	Port int
}

type request struct {
	URL    string
	Listen addr
	Remote addr
}

func (a addr) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("ip", a.IP)
	enc.AddInt("port", a.Port)
	return nil
}

func (r *request) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("url", r.URL)
	zap.Inline(r.Listen).AddTo(enc)
	return enc.AddObject("remote", r.Remote)
}

func main() {
	logger := zap.NewExample()
	defer logger.Sync()

	req := &request{
		URL:    "/test",
		Listen: addr{"127.0.0.1", 8080},
		Remote: addr{"127.0.0.1", 31200},
	}
	logger.Info("new request, in nested object", zap.Object("req", req))
	logger.Info("new request, inline", zap.Inline(req))
}

Output:

{"level":"info","msg":"new request, in nested object","req":{"url":"/test","ip":"127.0.0.1","port":8080,"remote":{"ip":"127.0.0.1","port":31200}}}
{"level":"info","msg":"new request, inline","url":"/test","ip":"127.0.0.1","port":8080,"remote":{"ip":"127.0.0.1","port":31200}}

Share

Format

Run

func

ObjectValues

¶

added in

v1.22.0

func ObjectValues[T

any

, P

ObjectMarshalerPtr

[T]](key

string

, values []T)

Field

ObjectValues constructs a field with the given key, holding a list of the
provided objects, where pointers to these objects can be marshaled by Zap.

Note that pointers to these objects must implement zapcore.ObjectMarshaler.
That is, if you're trying to marshal a []Request, the MarshalLogObject
method must be declared on the *Request type, not the value (Request).
If it's on the value, use Objects.

Given an object that implements MarshalLogObject on the pointer receiver,
you can log a slice of those objects with ObjectValues like so:

type Request struct{ ... }
func (r *Request) MarshalLogObject(enc zapcore.ObjectEncoder) error

var requests []Request = ...
logger.Info("sending requests", zap.ObjectValues("requests", requests))

If instead, you have a slice of pointers of such an object, use the Objects
field constructor.

var requests []*Request = ...
logger.Info("sending requests", zap.Objects("requests", requests))

Example

¶

package main

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type addr struct {
	IP   string
	Port int
}

type request struct {
	URL    string
	Listen addr
	Remote addr
}

func (a addr) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("ip", a.IP)
	enc.AddInt("port", a.Port)
	return nil
}

func (r *request) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("url", r.URL)
	zap.Inline(r.Listen).AddTo(enc)
	return enc.AddObject("remote", r.Remote)
}

func main() {
	logger := zap.NewExample()
	defer logger.Sync()

	// Use the ObjectValues field constructor when you have a list of
	// objects that do not implement zapcore.ObjectMarshaler directly,
	// but on their pointer receivers.
	logger.Debug("starting tunnels",
		zap.ObjectValues("addrs", []request{
			{
				URL:    "/foo",
				Listen: addr{"127.0.0.1", 8080},
				Remote: addr{"123.45.67.89", 4040},
			},
			{
				URL:    "/bar",
				Listen: addr{"127.0.0.1", 8080},
				Remote: addr{"127.0.0.1", 31200},
			},
		}))
}

Output:

{"level":"debug","msg":"starting tunnels","addrs":[{"url":"/foo","ip":"127.0.0.1","port":8080,"remote":{"ip":"123.45.67.89","port":4040}},{"url":"/bar","ip":"127.0.0.1","port":8080,"remote":{"ip":"127.0.0.1","port":31200}}]}

Share

Format

Run

func

Objects

¶

added in

v1.22.0

func Objects[T

zapcore

.

ObjectMarshaler

](key

string

, values []T)

Field

Objects constructs a field with the given key, holding a list of the
provided objects that can be marshaled by Zap.

Note that these objects must implement zapcore.ObjectMarshaler directly.
That is, if you're trying to marshal a []Request, the MarshalLogObject
method must be declared on the Request type, not its pointer (*Request).
If it's on the pointer, use ObjectValues.

Given an object that implements MarshalLogObject on the value receiver, you
can log a slice of those objects with Objects like so:

type Author struct{ ... }
func (a Author) MarshalLogObject(enc zapcore.ObjectEncoder) error

var authors []Author = ...
logger.Info("loading article", zap.Objects("authors", authors))

Similarly, given a type that implements MarshalLogObject on its pointer
receiver, you can log a slice of pointers to that object with Objects like
so:

type Request struct{ ... }
func (r *Request) MarshalLogObject(enc zapcore.ObjectEncoder) error

var requests []*Request = ...
logger.Info("sending requests", zap.Objects("requests", requests))

If instead, you have a slice of values of such an object, use the
ObjectValues constructor.

var requests []Request = ...
logger.Info("sending requests", zap.ObjectValues("requests", requests))

Example

¶

package main

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type addr struct {
	IP   string
	Port int
}

func (a addr) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("ip", a.IP)
	enc.AddInt("port", a.Port)
	return nil
}

func main() {
	logger := zap.NewExample()
	defer logger.Sync()

	// Use the Objects field constructor when you have a list of objects,
	// all of which implement zapcore.ObjectMarshaler.
	logger.Debug("opening connections",
		zap.Objects("addrs", []addr{
			{IP: "123.45.67.89", Port: 4040},
			{IP: "127.0.0.1", Port: 4041},
			{IP: "192.168.0.1", Port: 4042},
		}))
}

Output:

{"level":"debug","msg":"opening connections","addrs":[{"ip":"123.45.67.89","port":4040},{"ip":"127.0.0.1","port":4041},{"ip":"192.168.0.1","port":4042}]}

Share

Format

Run

func

Reflect

¶

func Reflect(key

string

, val interface{})

Field

Reflect constructs a field with the given key and an arbitrary object. It uses
an encoding-appropriate, reflection-based function to lazily serialize nearly
any object into the logging context, but it's relatively slow and
allocation-heavy. Outside tests, Any is always a better choice.

If encoding fails (e.g., trying to serialize a map[int]string to JSON), Reflect
includes the error message in the final log output.

func

Skip

¶

func Skip()

Field

Skip constructs a no-op field, which is often useful when handling invalid
inputs in other Field constructors.

func

Stack

¶

func Stack(key

string

)

Field

Stack constructs a field that stores a stacktrace of the current goroutine
under provided key. Keep in mind that taking a stacktrace is eager and
expensive (relatively speaking); this function both makes an allocation and
takes about two microseconds.

func

StackSkip

¶

added in

v1.16.0

func StackSkip(key

string

, skip

int

)

Field

StackSkip constructs a field similarly to Stack, but also skips the given
number of frames from the top of the stacktrace.

func

String

¶

func String(key

string

, val

string

)

Field

String constructs a field with the given key and value.

func

Stringer

¶

func Stringer(key

string

, val

fmt

.

Stringer

)

Field

Stringer constructs a field with the given key and the output of the value's
String method. The Stringer's String method is called lazily.

func

Stringers

¶

added in

v1.23.0

func Stringers[T

fmt

.

Stringer

](key

string

, values []T)

Field

Stringers constructs a field with the given key, holding a list of the
output provided by the value's String method

Given an object that implements String on the value receiver, you
can log a slice of those objects with Objects like so:

type Request struct{ ... }
func (a Request) String() string

var requests []Request = ...
logger.Info("sending requests", zap.Stringers("requests", requests))

Note that these objects must implement fmt.Stringer directly.
That is, if you're trying to marshal a []Request, the String method
must be declared on the Request type, not its pointer (*Request).

func

Stringp

¶

added in

v1.13.0

func Stringp(key

string

, val *

string

)

Field

Stringp constructs a field that carries a *string. The returned Field will safely
and explicitly represent `nil` when appropriate.

func

Strings

¶

func Strings(key

string

, ss []

string

)

Field

Strings constructs a field that carries a slice of strings.

func

Time

¶

func Time(key

string

, val

time

.

Time

)

Field

Time constructs a Field with the given key and value. The encoder
controls how the time is serialized.

func

Timep

¶

added in

v1.13.0

func Timep(key

string

, val *

time

.

Time

)

Field

Timep constructs a field that carries a *time.Time. The returned Field will safely
and explicitly represent `nil` when appropriate.

func

Times

¶

func Times(key

string

, ts []

time

.

Time

)

Field

Times constructs a field that carries a slice of time.Times.

func

Uint

¶

func Uint(key

string

, val

uint

)

Field

Uint constructs a field with the given key and value.

func

Uint16

¶

func Uint16(key

string

, val

uint16

)

Field

Uint16 constructs a field with the given key and value.

func

Uint16p

¶

added in

v1.13.0

func Uint16p(key

string

, val *

uint16

)

Field

Uint16p constructs a field that carries a *uint16. The returned Field will safely
and explicitly represent `nil` when appropriate.

func

Uint16s

¶

func Uint16s(key

string

, nums []

uint16

)

Field

Uint16s constructs a field that carries a slice of unsigned integers.

func

Uint32

¶

func Uint32(key

string

, val

uint32

)

Field

Uint32 constructs a field with the given key and value.

func

Uint32p

¶

added in

v1.13.0

func Uint32p(key

string

, val *

uint32

)

Field

Uint32p constructs a field that carries a *uint32. The returned Field will safely
and explicitly represent `nil` when appropriate.

func

Uint32s

¶

func Uint32s(key

string

, nums []

uint32

)

Field

Uint32s constructs a field that carries a slice of unsigned integers.

func

Uint64

¶

func Uint64(key

string

, val

uint64

)

Field

Uint64 constructs a field with the given key and value.

func

Uint64p

¶

added in

v1.13.0

func Uint64p(key

string

, val *

uint64

)

Field

Uint64p constructs a field that carries a *uint64. The returned Field will safely
and explicitly represent `nil` when appropriate.

func

Uint64s

¶

func Uint64s(key

string

, nums []

uint64

)

Field

Uint64s constructs a field that carries a slice of unsigned integers.

func

Uint8

¶

func Uint8(key

string

, val

uint8

)

Field

Uint8 constructs a field with the given key and value.

func

Uint8p

¶

added in

v1.13.0

func Uint8p(key

string

, val *

uint8

)

Field

Uint8p constructs a field that carries a *uint8. The returned Field will safely
and explicitly represent `nil` when appropriate.

func

Uint8s

¶

func Uint8s(key

string

, nums []

uint8

)

Field

Uint8s constructs a field that carries a slice of unsigned integers.

func

Uintp

¶

added in

v1.13.0

func Uintp(key

string

, val *

uint

)

Field

Uintp constructs a field that carries a *uint. The returned Field will safely
and explicitly represent `nil` when appropriate.

func

Uintptr

¶

func Uintptr(key

string

, val

uintptr

)

Field

Uintptr constructs a field with the given key and value.

func

Uintptrp

¶

added in

v1.13.0

func Uintptrp(key

string

, val *

uintptr

)

Field

Uintptrp constructs a field that carries a *uintptr. The returned Field will safely
and explicitly represent `nil` when appropriate.

func

Uintptrs

¶

func Uintptrs(key

string

, us []

uintptr

)

Field

Uintptrs constructs a field that carries a slice of pointer addresses.

func

Uints

¶

func Uints(key

string

, nums []

uint

)

Field

Uints constructs a field that carries a slice of unsigned integers.

type

LevelEnablerFunc

¶

type LevelEnablerFunc func(

zapcore

.

Level

)

bool

LevelEnablerFunc is a convenient way to implement zapcore.LevelEnabler with
an anonymous function.

It's particularly useful when splitting log output between different
outputs (e.g., standard error and standard out). For sample code, see the
package-level AdvancedConfiguration example.

func (LevelEnablerFunc)

Enabled

¶

func (f

LevelEnablerFunc

) Enabled(lvl

zapcore

.

Level

)

bool

Enabled calls the wrapped function.

type

Logger

¶

type Logger struct {

// contains filtered or unexported fields

}

A Logger provides fast, leveled, structured logging. All methods are safe
for concurrent use.

The Logger is designed for contexts in which every microsecond and every
allocation matters, so its API intentionally favors performance and type
safety over brevity. For most applications, the SugaredLogger strikes a
better balance between performance and ergonomics.

func

L

¶

func L() *

Logger

L returns the global Logger, which can be reconfigured with ReplaceGlobals.
It's safe for concurrent use.

func

Must

¶

added in

v1.22.0

func Must(logger *

Logger

, err

error

) *

Logger

Must is a helper that wraps a call to a function returning (*Logger, error)
and panics if the error is non-nil. It is intended for use in variable
initialization such as:

var logger = zap.Must(zap.NewProduction())

func

New

¶

func New(core

zapcore

.

Core

, options ...

Option

) *

Logger

New constructs a new Logger from the provided zapcore.Core and Options. If
the passed zapcore.Core is nil, it falls back to using a no-op
implementation.

This is the most flexible way to construct a Logger, but also the most
verbose. For typical use cases, the highly-opinionated presets
(NewProduction, NewDevelopment, and NewExample) or the Config struct are
more convenient.

For sample code, see the package-level AdvancedConfiguration example.

func

NewDevelopment

¶

func NewDevelopment(options ...

Option

) (*

Logger

,

error

)

NewDevelopment builds a development Logger that writes DebugLevel and above
logs to standard error in a human-friendly format.

It's a shortcut for NewDevelopmentConfig().Build(...Option).

func

NewExample

¶

added in

v1.5.0

func NewExample(options ...

Option

) *

Logger

NewExample builds a Logger that's designed for use in zap's testable
examples. It writes DebugLevel and above logs to standard out as JSON, but
omits the timestamp and calling function to keep example output
short and deterministic.

func

NewNop

¶

func NewNop() *

Logger

NewNop returns a no-op Logger. It never writes out logs or internal errors,
and it never runs user-defined hooks.

Using WithOptions to replace the Core or error output of a no-op Logger can
re-enable logging.

func

NewProduction

¶

func NewProduction(options ...

Option

) (*

Logger

,

error

)

NewProduction builds a sensible production Logger that writes InfoLevel and
above logs to standard error as JSON.

It's a shortcut for NewProductionConfig().Build(...Option).

func (*Logger)

Check

¶

func (log *

Logger

) Check(lvl

zapcore

.

Level

, msg

string

) *

zapcore

.

CheckedEntry

Check returns a CheckedEntry if logging a message at the specified level
is enabled. It's a completely optional optimization; in high-performance
applications, Check can help avoid allocating a slice to hold fields.

Example

¶

package main

import (
	"go.uber.org/zap"
)

func main() {
	logger := zap.NewExample()
	defer logger.Sync()

	if ce := logger.Check(zap.DebugLevel, "debugging"); ce != nil {
		// If debug-level log output isn't enabled or if zap's sampling would have
		// dropped this log entry, we don't allocate the slice that holds these
		// fields.
		ce.Write(
			zap.String("foo", "bar"),
			zap.String("baz", "quux"),
		)
	}

}

Output:

{"level":"debug","msg":"debugging","foo":"bar","baz":"quux"}

Share

Format

Run

func (*Logger)

Core

¶

func (log *

Logger

) Core()

zapcore

.

Core

Core returns the Logger's underlying zapcore.Core.

func (*Logger)

DPanic

¶

func (log *

Logger

) DPanic(msg

string

, fields ...

Field

)

DPanic logs a message at DPanicLevel. The message includes any fields
passed at the log site, as well as any fields accumulated on the logger.

If the logger is in development mode, it then panics (DPanic means
"development panic"). This is useful for catching errors that are
recoverable, but shouldn't ever happen.

func (*Logger)

Debug

¶

func (log *

Logger

) Debug(msg

string

, fields ...

Field

)

Debug logs a message at DebugLevel. The message includes any fields passed
at the log site, as well as any fields accumulated on the logger.

func (*Logger)

Error

¶

func (log *

Logger

) Error(msg

string

, fields ...

Field

)

Error logs a message at ErrorLevel. The message includes any fields passed
at the log site, as well as any fields accumulated on the logger.

func (*Logger)

Fatal

¶

func (log *

Logger

) Fatal(msg

string

, fields ...

Field

)

Fatal logs a message at FatalLevel. The message includes any fields passed
at the log site, as well as any fields accumulated on the logger.

The logger then calls os.Exit(1), even if logging at FatalLevel is
disabled.

func (*Logger)

Info

¶

func (log *

Logger

) Info(msg

string

, fields ...

Field

)

Info logs a message at InfoLevel. The message includes any fields passed
at the log site, as well as any fields accumulated on the logger.

func (*Logger)

Level

¶

added in

v1.24.0

func (log *

Logger

) Level()

zapcore

.

Level

Level reports the minimum enabled level for this logger.

For NopLoggers, this is

zapcore.InvalidLevel

.

func (*Logger)

Log

¶

added in

v1.22.0

func (log *

Logger

) Log(lvl

zapcore

.

Level

, msg

string

, fields ...

Field

)

Log logs a message at the specified level. The message includes any fields
passed at the log site, as well as any fields accumulated on the logger.
Any Fields that require  evaluation (such as Objects) are evaluated upon
invocation of Log.

func (*Logger)

Name

¶

added in

v1.25.0

func (log *

Logger

) Name()

string

Name returns the Logger's underlying name,
or an empty string if the logger is unnamed.

func (*Logger)

Named

¶

func (log *

Logger

) Named(s

string

) *

Logger

Named adds a new path segment to the logger's name. Segments are joined by
periods. By default, Loggers are unnamed.

Example

¶

package main

import (
	"go.uber.org/zap"
)

func main() {
	logger := zap.NewExample()
	defer logger.Sync()

	// By default, Loggers are unnamed.
	logger.Info("no name")

	// The first call to Named sets the Logger name.
	main := logger.Named("main")
	main.Info("main logger")

	// Additional calls to Named create a period-separated path.
	main.Named("subpackage").Info("sub-logger")
}

Output:

{"level":"info","msg":"no name"}
{"level":"info","logger":"main","msg":"main logger"}
{"level":"info","logger":"main.subpackage","msg":"sub-logger"}

Share

Format

Run

func (*Logger)

Panic

¶

func (log *

Logger

) Panic(msg

string

, fields ...

Field

)

Panic logs a message at PanicLevel. The message includes any fields passed
at the log site, as well as any fields accumulated on the logger.

The logger then panics, even if logging at PanicLevel is disabled.

func (*Logger)

Sugar

¶

func (log *

Logger

) Sugar() *

SugaredLogger

Sugar wraps the Logger to provide a more ergonomic, but slightly slower,
API. Sugaring a Logger is quite inexpensive, so it's reasonable for a
single application to use both Loggers and SugaredLoggers, converting
between them on the boundaries of performance-sensitive code.

func (*Logger)

Sync

¶

func (log *

Logger

) Sync()

error

Sync calls the underlying Core's Sync method, flushing any buffered log
entries. Applications should take care to call Sync before exiting.

func (*Logger)

Warn

¶

func (log *

Logger

) Warn(msg

string

, fields ...

Field

)

Warn logs a message at WarnLevel. The message includes any fields passed
at the log site, as well as any fields accumulated on the logger.

func (*Logger)

With

¶

func (log *

Logger

) With(fields ...

Field

) *

Logger

With creates a child logger and adds structured context to it. Fields added
to the child don't affect the parent, and vice versa. Any fields that
require evaluation (such as Objects) are evaluated upon invocation of With.

func (*Logger)

WithLazy

¶

added in

v1.26.0

func (log *

Logger

) WithLazy(fields ...

Field

) *

Logger

WithLazy creates a child logger and adds structured context to it lazily.

The fields are evaluated only if the logger is further chained with [With]
or is written to with any of the log level methods.
Until that occurs, the logger may retain references to objects inside the fields,
and logging will reflect the state of an object at the time of logging,
not the time of WithLazy().

WithLazy provides a worthwhile performance optimization for contextual loggers
when the likelihood of using the child logger is low,
such as error paths and rarely taken branches.

Similar to [With], fields added to the child don't affect the parent, and vice versa.

func (*Logger)

WithOptions

¶

func (log *

Logger

) WithOptions(opts ...

Option

) *

Logger

WithOptions clones the current Logger, applies the supplied Options, and
returns the resulting Logger. It's safe to use concurrently.

type

ObjectMarshalerPtr

¶

added in

v1.24.0

type ObjectMarshalerPtr[T

any

] interface {
	*T

zapcore

.

ObjectMarshaler

}

ObjectMarshalerPtr is a constraint that specifies that the given type
implements zapcore.ObjectMarshaler on a pointer receiver.

type

Option

¶

type Option interface {

// contains filtered or unexported methods

}

An Option configures a Logger.

func

AddCaller

¶

func AddCaller()

Option

AddCaller configures the Logger to annotate each message with the filename,
line number, and function name of zap's caller. See also WithCaller.

func

AddCallerSkip

¶

func AddCallerSkip(skip

int

)

Option

AddCallerSkip increases the number of callers skipped by caller annotation
(as enabled by the AddCaller option). When building wrappers around the
Logger and SugaredLogger, supplying this Option prevents zap from always
reporting the wrapper code as the caller.

func

AddStacktrace

¶

func AddStacktrace(lvl

zapcore

.

LevelEnabler

)

Option

AddStacktrace configures the Logger to record a stack trace for all messages at
or above a given level.

func

Development

¶

func Development()

Option

Development puts the logger in development mode, which makes DPanic-level
logs panic instead of simply logging an error.

func

ErrorOutput

¶

func ErrorOutput(w

zapcore

.

WriteSyncer

)

Option

ErrorOutput sets the destination for errors generated by the Logger. Note
that this option only affects internal errors; for sample code that sends
error-level logs to a different location from info- and debug-level logs,
see the package-level AdvancedConfiguration example.

The supplied WriteSyncer must be safe for concurrent use. The Open and
zapcore.Lock functions are the simplest ways to protect files with a mutex.

func

Fields

¶

func Fields(fs ...

Field

)

Option

Fields adds fields to the Logger.

func

Hooks

¶

func Hooks(hooks ...func(

zapcore

.

Entry

)

error

)

Option

Hooks registers functions which will be called each time the Logger writes
out an Entry. Repeated use of Hooks is additive.

Hooks are useful for simple side effects, like capturing metrics for the
number of emitted logs. More complex side effects, including anything that
requires access to the Entry's structured fields, should be implemented as
a zapcore.Core instead. See zapcore.RegisterHooks for details.

func

IncreaseLevel

¶

added in

v1.14.0

func IncreaseLevel(lvl

zapcore

.

LevelEnabler

)

Option

IncreaseLevel increase the level of the logger. It has no effect if
the passed in level tries to decrease the level of the logger.

func

OnFatal

deprecated

added in

v1.16.0

func OnFatal(action

zapcore

.

CheckWriteAction

)

Option

OnFatal sets the action to take on fatal logs.

Deprecated: Use

WithFatalHook

instead.

func

WithCaller

¶

added in

v1.15.0

func WithCaller(enabled

bool

)

Option

WithCaller configures the Logger to annotate each message with the filename,
line number, and function name of zap's caller, or not, depending on the
value of enabled. This is a generalized form of AddCaller.

func

WithClock

¶

added in

v1.18.0

func WithClock(clock

zapcore

.

Clock

)

Option

WithClock specifies the clock used by the logger to determine the current
time for logged entries. Defaults to the system clock with time.Now.

func

WithFatalHook

¶

added in

v1.22.0

func WithFatalHook(hook

zapcore

.

CheckWriteHook

)

Option

WithFatalHook sets a CheckWriteHook to run on fatal logs.
Zap will call this hook after writing a log statement with a Fatal level.

For example, the following builds a logger that will exit the current
goroutine after writing a fatal log message, but it will not exit the
program.

zap.New(core, zap.WithFatalHook(zapcore.WriteThenGoexit))

It is important that the provided CheckWriteHook stops the control flow at
the current statement to meet expectations of callers of the logger.
We recommend calling os.Exit or runtime.Goexit inside custom hooks at
minimum.

func

WithPanicHook

¶

added in

v1.27.0

func WithPanicHook(hook

zapcore

.

CheckWriteHook

)

Option

WithPanicHook sets a CheckWriteHook to run on Panic/DPanic logs.
Zap will call this hook after writing a log statement with a Panic/DPanic level.

For example, the following builds a logger that will exit the current
goroutine after writing a Panic/DPanic log message, but it will not start a panic.

zap.New(core, zap.WithPanicHook(zapcore.WriteThenGoexit))

This is useful for testing Panic/DPanic log output.

func

WrapCore

¶

func WrapCore(f func(

zapcore

.

Core

)

zapcore

.

Core

)

Option

WrapCore wraps or replaces the Logger's underlying zapcore.Core.

Example (Replace)

¶

package main

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	// Replacing a Logger's core can alter fundamental behaviors.
	// For example, it can convert a Logger to a no-op.
	nop := zap.WrapCore(func(zapcore.Core) zapcore.Core {
		return zapcore.NewNopCore()
	})

	logger := zap.NewExample()
	defer logger.Sync()

	logger.Info("working")
	logger.WithOptions(nop).Info("no-op")
	logger.Info("original logger still works")
}

Output:

{"level":"info","msg":"working"}
{"level":"info","msg":"original logger still works"}

Share

Format

Run

Example (Wrap)

¶

package main

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	// Wrapping a Logger's core can extend its functionality. As a trivial
	// example, it can double-write all logs.
	doubled := zap.WrapCore(func(c zapcore.Core) zapcore.Core {
		return zapcore.NewTee(c, c)
	})

	logger := zap.NewExample()
	defer logger.Sync()

	logger.Info("single")
	logger.WithOptions(doubled).Info("doubled")
}

Output:

{"level":"info","msg":"single"}
{"level":"info","msg":"doubled"}
{"level":"info","msg":"doubled"}

Share

Format

Run

type

SamplingConfig

¶

type SamplingConfig struct {

Initial

int

`json:"initial" yaml:"initial"`

Thereafter

int

`json:"thereafter" yaml:"thereafter"`

Hook       func(

zapcore

.

Entry

,

zapcore

.

SamplingDecision

) `json:"-" yaml:"-"`

}

SamplingConfig sets a sampling strategy for the logger. Sampling caps the
global CPU and I/O load that logging puts on your process while attempting
to preserve a representative subset of your logs.

If specified, the Sampler will invoke the Hook after each decision.

Values configured here are per-second. See zapcore.NewSamplerWithOptions for
details.

type

Sink

¶

added in

v1.9.0

type Sink interface {

zapcore

.

WriteSyncer

io

.

Closer

}

Sink defines the interface to write to and close logger destinations.

type

SugaredLogger

¶

type SugaredLogger struct {

// contains filtered or unexported fields

}

A SugaredLogger wraps the base Logger functionality in a slower, but less
verbose, API. Any Logger can be converted to a SugaredLogger with its Sugar
method.

Unlike the Logger, the SugaredLogger doesn't insist on structured logging.
For each log level, it exposes four methods:

methods named after the log level for log.Print-style logging

methods ending in "w" for loosely-typed structured logging

methods ending in "f" for log.Printf-style logging

methods ending in "ln" for log.Println-style logging

For example, the methods for InfoLevel are:

Info(...any)           Print-style logging
Infow(...any)          Structured logging (read as "info with")
Infof(string, ...any)  Printf-style logging
Infoln(...any)         Println-style logging

func

S

¶

func S() *

SugaredLogger

S returns the global SugaredLogger, which can be reconfigured with
ReplaceGlobals. It's safe for concurrent use.

func (*SugaredLogger)

DPanic

¶

func (s *

SugaredLogger

) DPanic(args ...interface{})

DPanic logs the provided arguments at

DPanicLevel

.
In development, the logger then panics. (See

DPanicLevel

for details.)
Spaces are added between arguments when neither is a string.

func (*SugaredLogger)

DPanicf

¶

func (s *

SugaredLogger

) DPanicf(template

string

, args ...interface{})

DPanicf formats the message according to the format specifier
and logs it at

DPanicLevel

.
In development, the logger then panics. (See

DPanicLevel

for details.)

func (*SugaredLogger)

DPanicln

¶

added in

v1.22.0

func (s *

SugaredLogger

) DPanicln(args ...interface{})

DPanicln logs a message at

DPanicLevel

.
In development, the logger then panics. (See

DPanicLevel

for details.)
Spaces are always added between arguments.

func (*SugaredLogger)

DPanicw

¶

func (s *

SugaredLogger

) DPanicw(msg

string

, keysAndValues ...interface{})

DPanicw logs a message with some additional context. In development, the
logger then panics. (See DPanicLevel for details.) The variadic key-value
pairs are treated as they are in With.

func (*SugaredLogger)

Debug

¶

func (s *

SugaredLogger

) Debug(args ...interface{})

Debug logs the provided arguments at

DebugLevel

.
Spaces are added between arguments when neither is a string.

func (*SugaredLogger)

Debugf

¶

func (s *

SugaredLogger

) Debugf(template

string

, args ...interface{})

Debugf formats the message according to the format specifier
and logs it at

DebugLevel

.

func (*SugaredLogger)

Debugln

¶

added in

v1.22.0

func (s *

SugaredLogger

) Debugln(args ...interface{})

Debugln logs a message at

DebugLevel

.
Spaces are always added between arguments.

func (*SugaredLogger)

Debugw

¶

func (s *

SugaredLogger

) Debugw(msg

string

, keysAndValues ...interface{})

Debugw logs a message with some additional context. The variadic key-value
pairs are treated as they are in With.

When debug-level logging is disabled, this is much faster than

s.With(keysAndValues).Debug(msg)

func (*SugaredLogger)

Desugar

¶

func (s *

SugaredLogger

) Desugar() *

Logger

Desugar unwraps a SugaredLogger, exposing the original Logger. Desugaring
is quite inexpensive, so it's reasonable for a single application to use
both Loggers and SugaredLoggers, converting between them on the boundaries
of performance-sensitive code.

func (*SugaredLogger)

Error

¶

func (s *

SugaredLogger

) Error(args ...interface{})

Error logs the provided arguments at

ErrorLevel

.
Spaces are added between arguments when neither is a string.

func (*SugaredLogger)

Errorf

¶

func (s *

SugaredLogger

) Errorf(template

string

, args ...interface{})

Errorf formats the message according to the format specifier
and logs it at

ErrorLevel

.

func (*SugaredLogger)

Errorln

¶

added in

v1.22.0

func (s *

SugaredLogger

) Errorln(args ...interface{})

Errorln logs a message at

ErrorLevel

.
Spaces are always added between arguments.

func (*SugaredLogger)

Errorw

¶

func (s *

SugaredLogger

) Errorw(msg

string

, keysAndValues ...interface{})

Errorw logs a message with some additional context. The variadic key-value
pairs are treated as they are in With.

func (*SugaredLogger)

Fatal

¶

func (s *

SugaredLogger

) Fatal(args ...interface{})

Fatal constructs a message with the provided arguments and calls os.Exit.
Spaces are added between arguments when neither is a string.

func (*SugaredLogger)

Fatalf

¶

func (s *

SugaredLogger

) Fatalf(template

string

, args ...interface{})

Fatalf formats the message according to the format specifier
and calls os.Exit.

func (*SugaredLogger)

Fatalln

¶

added in

v1.22.0

func (s *

SugaredLogger

) Fatalln(args ...interface{})

Fatalln logs a message at

FatalLevel

and calls os.Exit.
Spaces are always added between arguments.

func (*SugaredLogger)

Fatalw

¶

func (s *

SugaredLogger

) Fatalw(msg

string

, keysAndValues ...interface{})

Fatalw logs a message with some additional context, then calls os.Exit. The
variadic key-value pairs are treated as they are in With.

func (*SugaredLogger)

Info

¶

func (s *

SugaredLogger

) Info(args ...interface{})

Info logs the provided arguments at

InfoLevel

.
Spaces are added between arguments when neither is a string.

func (*SugaredLogger)

Infof

¶

func (s *

SugaredLogger

) Infof(template

string

, args ...interface{})

Infof formats the message according to the format specifier
and logs it at

InfoLevel

.

func (*SugaredLogger)

Infoln

¶

added in

v1.22.0

func (s *

SugaredLogger

) Infoln(args ...interface{})

Infoln logs a message at

InfoLevel

.
Spaces are always added between arguments.

func (*SugaredLogger)

Infow

¶

func (s *

SugaredLogger

) Infow(msg

string

, keysAndValues ...interface{})

Infow logs a message with some additional context. The variadic key-value
pairs are treated as they are in With.

func (*SugaredLogger)

Level

¶

added in

v1.24.0

func (s *

SugaredLogger

) Level()

zapcore

.

Level

Level reports the minimum enabled level for this logger.

For NopLoggers, this is

zapcore.InvalidLevel

.

func (*SugaredLogger)

Log

¶

added in

v1.27.0

func (s *

SugaredLogger

) Log(lvl

zapcore

.

Level

, args ...interface{})

Log logs the provided arguments at provided level.
Spaces are added between arguments when neither is a string.

func (*SugaredLogger)

Logf

¶

added in

v1.27.0

func (s *

SugaredLogger

) Logf(lvl

zapcore

.

Level

, template

string

, args ...interface{})

Logf formats the message according to the format specifier
and logs it at provided level.

func (*SugaredLogger)

Logln

¶

added in

v1.27.0

func (s *

SugaredLogger

) Logln(lvl

zapcore

.

Level

, args ...interface{})

Logln logs a message at provided level.
Spaces are always added between arguments.

func (*SugaredLogger)

Logw

¶

added in

v1.27.0

func (s *

SugaredLogger

) Logw(lvl

zapcore

.

Level

, msg

string

, keysAndValues ...interface{})

Logw logs a message with some additional context. The variadic key-value
pairs are treated as they are in With.

func (*SugaredLogger)

Named

¶

func (s *

SugaredLogger

) Named(name

string

) *

SugaredLogger

Named adds a sub-scope to the logger's name. See Logger.Named for details.

func (*SugaredLogger)

Panic

¶

func (s *

SugaredLogger

) Panic(args ...interface{})

Panic constructs a message with the provided arguments and panics.
Spaces are added between arguments when neither is a string.

func (*SugaredLogger)

Panicf

¶

func (s *

SugaredLogger

) Panicf(template

string

, args ...interface{})

Panicf formats the message according to the format specifier
and panics.

func (*SugaredLogger)

Panicln

¶

added in

v1.22.0

func (s *

SugaredLogger

) Panicln(args ...interface{})

Panicln logs a message at

PanicLevel

and panics.
Spaces are always added between arguments.

func (*SugaredLogger)

Panicw

¶

func (s *

SugaredLogger

) Panicw(msg

string

, keysAndValues ...interface{})

Panicw logs a message with some additional context, then panics. The
variadic key-value pairs are treated as they are in With.

func (*SugaredLogger)

Sync

¶

func (s *

SugaredLogger

) Sync()

error

Sync flushes any buffered log entries.

func (*SugaredLogger)

Warn

¶

func (s *

SugaredLogger

) Warn(args ...interface{})

Warn logs the provided arguments at

WarnLevel

.
Spaces are added between arguments when neither is a string.

func (*SugaredLogger)

Warnf

¶

func (s *

SugaredLogger

) Warnf(template

string

, args ...interface{})

Warnf formats the message according to the format specifier
and logs it at

WarnLevel

.

func (*SugaredLogger)

Warnln

¶

added in

v1.22.0

func (s *

SugaredLogger

) Warnln(args ...interface{})

Warnln logs a message at

WarnLevel

.
Spaces are always added between arguments.

func (*SugaredLogger)

Warnw

¶

func (s *

SugaredLogger

) Warnw(msg

string

, keysAndValues ...interface{})

Warnw logs a message with some additional context. The variadic key-value
pairs are treated as they are in With.

func (*SugaredLogger)

With

¶

func (s *

SugaredLogger

) With(args ...interface{}) *

SugaredLogger

With adds a variadic number of fields to the logging context. It accepts a
mix of strongly-typed Field objects and loosely-typed key-value pairs. When
processing pairs, the first element of the pair is used as the field key
and the second as the field value.

For example,

sugaredLogger.With(
   "hello", "world",
   "failure", errors.New("oh no"),
   Stack(),
   "count", 42,
   "user", User{Name: "alice"},
)

is the equivalent of

unsugared.With(
  String("hello", "world"),
  String("failure", "oh no"),
  Stack(),
  Int("count", 42),
  Object("user", User{Name: "alice"}),
)

Note that the keys in key-value pairs should be strings. In development,
passing a non-string key panics. In production, the logger is more
forgiving: a separate error is logged, but the key-value pair is skipped
and execution continues. Passing an orphaned key triggers similar behavior:
panics in development and errors in production.

func (*SugaredLogger)

WithLazy

¶

added in

v1.27.0

func (s *

SugaredLogger

) WithLazy(args ...interface{}) *

SugaredLogger

WithLazy adds a variadic number of fields to the logging context lazily.
The fields are evaluated only if the logger is further chained with [With]
or is written to with any of the log level methods.
Until that occurs, the logger may retain references to objects inside the fields,
and logging will reflect the state of an object at the time of logging,
not the time of WithLazy().

Similar to [With], fields added to the child don't affect the parent,
and vice versa. Also, the keys in key-value pairs should be strings. In development,
passing a non-string key panics, while in production it logs an error and skips the pair.
Passing an orphaned key has the same behavior.

func (*SugaredLogger)

WithOptions

¶

added in

v1.22.0

func (s *

SugaredLogger

) WithOptions(opts ...

Option

) *

SugaredLogger

WithOptions clones the current SugaredLogger, applies the supplied Options,
and returns the result. It's safe to use concurrently.