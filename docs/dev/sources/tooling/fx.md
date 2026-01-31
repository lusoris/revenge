# Uber fx

> Source: https://pkg.go.dev/go.uber.org/fx
> Fetched: 2026-01-30T23:48:21.671338+00:00
> Content-Hash: 24e92f5448ab505e
> Type: html

---

Overview

¶

Basic usage

Testing Fx Applications

Parameter Structs

Result Structs

Named Values

Optional Dependencies

Value Groups

Soft Value Groups

Value group flattening

Unexported fields

Package fx is a framework that makes it easy to build applications out of
reusable, composable modules.

Fx applications use dependency injection to eliminate globals without the
tedium of manually wiring together function calls. Unlike other approaches
to dependency injection, Fx works with plain Go functions: you don't need
to use struct tags or embed special types, so Fx automatically works well
with most Go packages.

Basic usage

¶

Basic usage is explained in the package-level example.
If you're new to Fx, start there!

Advanced features, including named instances, optional parameters,
and value groups, are explained in this section further down.

Testing Fx Applications

¶

To test functions that use the Lifecycle type or to write end-to-end tests
of your Fx application, use the helper functions and types provided by the
go.uber.org/fx/fxtest package.

Parameter Structs

¶

Fx constructors declare their dependencies as function parameters. This can
quickly become unreadable if the constructor has a lot of dependencies.

func NewHandler(users *UserGateway, comments *CommentGateway, posts *PostGateway, votes *VoteGateway, authz *AuthZGateway) *Handler {
	// ...
}

To improve the readability of constructors like this, create a struct that
lists all the dependencies as fields and change the function to accept that
struct instead. The new struct is called a parameter struct.

Fx has first class support for parameter structs: any struct embedding
fx.In gets treated as a parameter struct, so the individual fields in the
struct are supplied via dependency injection. Using a parameter struct, we
can make the constructor above much more readable:

type HandlerParams struct {
	fx.In

	Users    *UserGateway
	Comments *CommentGateway
	Posts    *PostGateway
	Votes    *VoteGateway
	AuthZ    *AuthZGateway
}

func NewHandler(p HandlerParams) *Handler {
	// ...
}

Though it's rarelly necessary to mix the two, constructors can receive any
combination of parameter structs and parameters.

func NewHandler(p HandlerParams, l *log.Logger) *Handler {
	// ...
}

Result Structs

¶

Result structs are the inverse of parameter structs.
These structs represent multiple outputs from a
single function as fields. Fx treats all structs embedding fx.Out as result
structs, so other constructors can rely on the result struct's fields
directly.

Without result structs, we sometimes have function definitions like this:

func SetupGateways(conn *sql.DB) (*UserGateway, *CommentGateway, *PostGateway, error) {
	// ...
}

With result structs, we can make this both more readable and easier to
modify in the future:

type Gateways struct {
	fx.Out

	Users    *UserGateway
	Comments *CommentGateway
	Posts    *PostGateway
}

func SetupGateways(conn *sql.DB) (Gateways, error) {
	// ...
}

Named Values

¶

Some use cases require the application container to hold multiple values of
the same type.

A constructor that produces a result struct can tag any field with
`name:".."` to have the corresponding value added to the graph under the
specified name. An application may contain at most one unnamed value of a
given type, but may contain any number of named values of the same type.

type ConnectionResult struct {
	fx.Out

	ReadWrite *sql.DB `name:"rw"`
	ReadOnly  *sql.DB `name:"ro"`
}

func ConnectToDatabase(...) (ConnectionResult, error) {
	// ...
	return ConnectionResult{ReadWrite: rw, ReadOnly:  ro}, nil
}

Similarly, a constructor that accepts a parameter struct can tag any field
with `name:".."` to have the corresponding value injected by name.

type GatewayParams struct {
	fx.In

	WriteToConn  *sql.DB `name:"rw"`
	ReadFromConn *sql.DB `name:"ro"`
}

func NewCommentGateway(p GatewayParams) (*CommentGateway, error) {
	// ...
}

Note that both the name AND type of the fields on the
parameter struct must match the corresponding result struct.

Optional Dependencies

¶

Constructors often have optional dependencies on some types: if those types are
missing, they can operate in a degraded state. Fx supports optional
dependencies via the `optional:"true"` tag to fields on parameter structs.

type UserGatewayParams struct {
	fx.In

	Conn  *sql.DB
	Cache *redis.Client `optional:"true"`
}

If an optional field isn't available in the container, the constructor
receives the field's zero value.

func NewUserGateway(p UserGatewayParams, log *log.Logger) (*UserGateway, error) {
	if p.Cache == nil {
		log.Print("Caching disabled")
	}
	// ...
}

Constructors that declare optional dependencies MUST gracefully handle
situations in which those dependencies are absent.

The optional tag also allows adding new dependencies without breaking
existing consumers of the constructor.

The optional tag may be combined with the name tag to declare a named
value dependency optional.

type GatewayParams struct {
	fx.In

	WriteToConn  *sql.DB `name:"rw"`
	ReadFromConn *sql.DB `name:"ro" optional:"true"`
}

func NewCommentGateway(p GatewayParams, log *log.Logger) (*CommentGateway, error) {
	if p.ReadFromConn == nil {
		log.Print("Warning: Using RW connection for reads")
		p.ReadFromConn = p.WriteToConn
	}
	// ...
}

Value Groups

¶

To make it easier to produce and consume many values of the same type, Fx
supports named, unordered collections called value groups.

Constructors can send values into value groups by returning a result struct
tagged with `group:".."`.

type HandlerResult struct {
	fx.Out

	Handler Handler `group:"server"`
}

func NewHelloHandler() HandlerResult {
	// ...
}

func NewEchoHandler() HandlerResult {
	// ...
}

Any number of constructors may provide values to this named collection, but
the ordering of the final collection is unspecified.

Value groups require parameter and result structs to use fields with
different types: if a group of constructors each returns type T, parameter
structs consuming the group must use a field of type []T.

Parameter structs can request a value group by using a field of type []T
tagged with `group:".."`.
This will execute all constructors that provide a value to
that group in an unspecified order, then collect all the results into a
single slice.

type ServerParams struct {
	fx.In

	Handlers []Handler `group:"server"`
}

func NewServer(p ServerParams) *Server {
	server := newServer()
	for _, h := range p.Handlers {
		server.Register(h)
	}
	return server
}

Note that values in a value group are unordered. Fx makes no guarantees
about the order in which these values will be produced.

Soft Value Groups

¶

By default, when a constructor declares a dependency on a value group,
all values provided to that value group are eagerly instantiated.
That is undesirable for cases where an optional component wants to
constribute to a value group, but only if it was actually used
by the rest of the application.

A soft value group can be thought of as a best-attempt at populating the
group with values from constructors that have already run. In other words,
if a constructor's output type is only consumed by a soft value group,
it will not be run.

Note that Fx randomizes the order of values in the value group,
so the slice of values may not match the order in which constructors
were run.

To declare a soft relationship between a group and its constructors, use
the `soft` option on the input group tag (`group:"[groupname],soft"`).
This option is only valid for input parameters.

type Params struct {
	fx.In

	Handlers []Handler `group:"server,soft"`
	Logger   *zap.Logger
}

func NewServer(p Params) *Server {
	// ...
}

With such a declaration, a constructor that provides a value to the 'server'
value group will be called only if there's another instantiated component
that consumes the results of that constructor.

func NewHandlerAndLogger() (Handler, *zap.Logger) {
	// ...
}

func NewHandler() Handler {
	// ...
}

fx.Provide(
	fx.Annotate(NewHandlerAndLogger, fx.ResultTags(`group:"server"`)),
	fx.Annotate(NewHandler, fx.ResultTags(`group:"server"`)),
)

NewHandlerAndLogger will be called because the Logger is consumed by the
application, but NewHandler will not be called because it's only consumed
by the soft value group.

Value group flattening

¶

By default, values of type T produced to a value group are consumed as []T.

type HandlerResult struct {
	fx.Out

	Handler Handler `group:"server"`
}

type ServerParams struct {
	fx.In

	Handlers []Handler `group:"server"`
}

This means that if the producer produces []T,
the consumer must consume [][]T.

There are cases where it's desirable
for the producer (the fx.Out) to produce multiple values ([]T),
and for the consumer (the fx.In) consume them as a single slice ([]T).
Fx offers flattened value groups for this purpose.

To provide multiple values for a group from a result struct, produce a
slice and use the `,flatten` option on the group tag. This indicates that
each element in the slice should be injected into the group individually.

type HandlerResult struct {
	fx.Out

	Handler []Handler `group:"server,flatten"`
	// Consumed as []Handler in ServerParams.
}

Unexported fields

¶

By default, a type that embeds fx.In may not have any unexported fields. The
following will return an error if used with Fx.

type Params struct {
	fx.In

	Logger *zap.Logger
	mu     sync.Mutex
}

If you have need of unexported fields on such a type, you may opt-into
ignoring unexported fields by adding the ignore-unexported struct tag to the
fx.In. For example,

type Params struct {
	fx.In `ignore-unexported:"true"`

	Logger *zap.Logger
	mu     sync.Mutex
}

Example

¶

package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
)

// NewLogger constructs a logger. It's just a regular Go function, without any
// special relationship to Fx.
//
// Since it returns a *log.Logger, Fx will treat NewLogger as the constructor
// function for the standard library's logger. (We'll see how to integrate
// NewLogger into an Fx application in the main function.) Since NewLogger
// doesn't have any parameters, Fx will infer that loggers don't depend on any
// other types - we can create them from thin air.
//
// Fx calls constructors lazily, so NewLogger will only be called only if some
// other function needs a logger. Once instantiated, the logger is cached and
// reused - within the application, it's effectively a singleton.
//
// By default, Fx applications only allow one constructor for each type. See
// the documentation of the In and Out types for ways around this restriction.
func NewLogger() *log.Logger {
	logger := log.New(os.Stdout, "" /* prefix */, 0 /* flags */)
	logger.Print("Executing NewLogger.")
	return logger
}

// NewHandler constructs a simple HTTP handler. Since it returns an
// http.Handler, Fx will treat NewHandler as the constructor for the
// http.Handler type.
//
// Like many Go functions, NewHandler also returns an error. If the error is
// non-nil, Go convention tells the caller to assume that NewHandler failed
// and the other returned values aren't safe to use. Fx understands this
// idiom, and assumes that any function whose last return value is an error
// follows this convention.
//
// Unlike NewLogger, NewHandler has formal parameters. Fx will interpret these
// parameters as dependencies: in order to construct an HTTP handler,
// NewHandler needs a logger. If the application has access to a *log.Logger
// constructor (like NewLogger above), it will use that constructor or its
// cached output and supply a logger to NewHandler. If the application doesn't
// know how to construct a logger and needs an HTTP handler, it will fail to
// start.
//
// Functions may also return multiple objects. For example, we could combine
// NewHandler and NewLogger into a single function:
//
//	func NewHandlerAndLogger() (*log.Logger, http.Handler, error)
//
// Fx also understands this idiom, and would treat NewHandlerAndLogger as the
// constructor for both the *log.Logger and http.Handler types. Just like
// constructors for a single type, NewHandlerAndLogger would be called at most
// once, and both the handler and the logger would be cached and reused as
// necessary.
func NewHandler(logger *log.Logger) (http.Handler, error) {
	logger.Print("Executing NewHandler.")
	return http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		logger.Print("Got a request.")
	}), nil
}

// NewMux constructs an HTTP mux. Like NewHandler, it depends on *log.Logger.
// However, it also depends on the Fx-specific Lifecycle interface.
//
// A Lifecycle is available in every Fx application. It lets objects hook into
// the application's start and stop phases. In a non-Fx application, the main
// function often includes blocks like this:
//
//	srv, err := NewServer() // some long-running network server
//	if err != nil {
//	  log.Fatalf("failed to construct server: %v", err)
//	}
//	// Construct other objects as necessary.
//	go srv.Start()
//	defer srv.Stop()
//
// In this example, the programmer explicitly constructs a bunch of objects,
// crashing the program if any of the constructors encounter unrecoverable
// errors. Once all the objects are constructed, we start any background
// goroutines and defer cleanup functions.
//
// Fx removes the manual object construction with dependency injection. It
// replaces the inline goroutine spawning and deferred cleanups with the
// Lifecycle type.
//
// Here, NewMux makes an HTTP mux available to other functions. Since
// constructors are called lazily, we know that NewMux won't be called unless
// some other function wants to register a handler. This makes it easy to use
// Fx's Lifecycle to start an HTTP server only if we have handlers registered.
func NewMux(lc fx.Lifecycle, logger *log.Logger) *http.ServeMux {
	logger.Print("Executing NewMux.")
	// First, we construct the mux and server. We don't want to start the server
	// until all handlers are registered.
	mux := http.NewServeMux()
	server := &http.Server{
		Addr:    "127.0.0.1:8080",
		Handler: mux,
	}
	// If NewMux is called, we know that another function is using the mux. In
	// that case, we'll use the Lifecycle type to register a Hook that starts
	// and stops our HTTP server.
	//
	// Hooks are executed in dependency order. At startup, NewLogger's hooks run
	// before NewMux's. On shutdown, the order is reversed.
	//
	// Returning an error from OnStart hooks interrupts application startup. Fx
	// immediately runs the OnStop portions of any successfully-executed OnStart
	// hooks (so that types which started cleanly can also shut down cleanly),
	// then exits.
	//
	// Returning an error from OnStop hooks logs a warning, but Fx continues to
	// run the remaining hooks.
	lc.Append(fx.Hook{
		// To mitigate the impact of deadlocks in application startup and
		// shutdown, Fx imposes a time limit on OnStart and OnStop hooks. By
		// default, hooks have a total of 15 seconds to complete. Timeouts are
		// passed via Go's usual context.Context.
		OnStart: func(context.Context) error {
			logger.Print("Starting HTTP server.")
			ln, err := net.Listen("tcp", server.Addr)
			if err != nil {
				return err
			}
			go server.Serve(ln)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Print("Stopping HTTP server.")
			return server.Shutdown(ctx)
		},
	})

	return mux
}

// Register mounts our HTTP handler on the mux.
//
// Register is a typical top-level application function: it takes a generic
// type like ServeMux, which typically comes from a third-party library, and
// introduces it to a type that contains our application logic. In this case,
// that introduction consists of registering an HTTP handler. Other typical
// examples include registering RPC procedures and starting queue consumers.
//
// Fx calls these functions invocations, and they're treated differently from
// the constructor functions above. Their arguments are still supplied via
// dependency injection and they may still return an error to indicate
// failure, but any other return values are ignored.
//
// Unlike constructors, invocations are called eagerly. See the main function
// below for details.
func Register(mux *http.ServeMux, h http.Handler) {
	mux.Handle("/", h)
}

func main() {
	app := fx.New(
		// Provide all the constructors we need, which teaches Fx how we'd like to
		// construct the *log.Logger, http.Handler, and *http.ServeMux types.
		// Remember that constructors are called lazily, so this block doesn't do
		// much on its own.
		fx.Provide(
			NewLogger,
			NewHandler,
			NewMux,
		),
		// Since constructors are called lazily, we need some invocations to
		// kick-start our application. In this case, we'll use Register. Since it
		// depends on an http.Handler and *http.ServeMux, calling it requires Fx
		// to build those types using the constructors above. Since we call
		// NewMux, we also register Lifecycle hooks to start and stop an HTTP
		// server.
		fx.Invoke(Register),

		// This is optional. With this, you can control where Fx logs
		// its events. In this case, we're using a NopLogger to keep
		// our test silent. Normally, you'll want to use an
		// fxevent.ZapLogger or an fxevent.ConsoleLogger.
		fx.WithLogger(
			func() fxevent.Logger {
				return fxevent.NopLogger
			},
		),
	)

	// In a typical application, we could just use app.Run() here. Since we
	// don't want this example to run forever, we'll use the more-explicit Start
	// and Stop.
	startCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := app.Start(startCtx); err != nil {
		log.Fatal(err)
	}

	// Normally, we'd block here with <-app.Done(). Instead, we'll make an HTTP
	// request to demonstrate that our server is running.
	if _, err := http.Get("http://localhost:8080/"); err != nil {
		log.Fatal(err)
	}

	stopCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := app.Stop(stopCtx); err != nil {
		log.Fatal(err)
	}

}

Output:

Executing NewLogger.
Executing NewMux.
Executing NewHandler.
Starting HTTP server.
Got a request.
Stopping HTTP server.

Share

Format

Run

Index

¶

Constants

Variables

func Annotate(t interface{}, anns ...Annotation) interface{}

func Self() any

func ValidateApp(opts ...Option) error

func VisualizeError(err error) (string, error)

type Annotated

func (a Annotated) String() string

type Annotation

func As(interfaces ...interface{}) Annotation

func From(interfaces ...interface{}) Annotation

func OnStart(onStart interface{}) Annotation

func OnStop(onStop interface{}) Annotation

func ParamTags(tags ...string) Annotation

func ResultTags(tags ...string) Annotation

type App

func New(opts ...Option) *App

func (app *App) Done() <-chan os.Signal

func (app *App) Err() error

func (app *App) Run()

func (app *App) Start(ctx context.Context) (err error)

func (app *App) StartTimeout() time.Duration

func (app *App) Stop(ctx context.Context) (err error)

func (app *App) StopTimeout() time.Duration

func (app *App) Wait() <-chan ShutdownSignal

type DotGraph

type ErrorHandler

type Hook

func StartHook[T HookFunc](start T) Hook

func StartStopHook[T1 HookFunc, T2 HookFunc](start T1, stop T2) Hook

func StopHook[T HookFunc](stop T) Hook

type HookFunc

type In

type Lifecycle

type Option

func Decorate(decorators ...interface{}) Option

func Error(errs ...error) Option

func ErrorHook(funcs ...ErrorHandler) Option

func Extract(target interface{}) Option

deprecated

func Invoke(funcs ...interface{}) Option

func Logger(p Printer) Option

func Module(name string, opts ...Option) Option

func Options(opts ...Option) Option

func Populate(targets ...interface{}) Option

func Provide(constructors ...interface{}) Option

func RecoverFromPanics() Option

func Replace(values ...interface{}) Option

func StartTimeout(v time.Duration) Option

func StopTimeout(v time.Duration) Option

func Supply(values ...interface{}) Option

func WithLogger(constructor interface{}) Option

type Out

type Printer

type ShutdownOption

func ExitCode(code int) ShutdownOption

func ShutdownTimeout(timeout time.Duration) ShutdownOption

deprecated

type ShutdownSignal

func (sig ShutdownSignal) String() string

type Shutdowner

Examples

¶

Package

Error

Populate

Constants

¶

View Source

const DefaultTimeout = 15 *

time

.

Second

DefaultTimeout is the default timeout for starting or stopping an
application. It can be configured with the

StartTimeout

and

StopTimeout

options.

View Source

const Version = "1.24.0"

Version is exported for runtime compatibility checks.

Variables

¶

View Source

var NopLogger =

WithLogger

(func()

fxevent

.

Logger

{ return

fxevent

.

NopLogger

})

NopLogger disables the application's log output.

Note that this makes some failures difficult to debug,
since no errors are printed to console.
Prefer to log to an in-memory buffer instead.

View Source

var Private = privateOption{}

Private is an option that can be passed as an argument to

Provide

or

Supply

to
restrict access to the constructors being provided. Specifically,
corresponding constructors can only be used within the current module
or modules the current module contains. Other modules that contain this
module won't be able to use the constructor.

For example, the following would fail because the app doesn't have access
to the inner module's constructor.

fx.New(
	fx.Module("SubModule", fx.Provide(func() int { return 0 }, fx.Private)),
	fx.Invoke(func(a int) {}),
)

Functions

¶

func

Annotate

¶

added in

v1.15.0

func Annotate(t interface{}, anns ...

Annotation

) interface{}

Variadic functions

Annotate lets you annotate a function's parameters and returns
without you having to declare separate struct definitions for them.

For example,

func NewGateway(ro, rw *db.Conn) *Gateway { ... }
fx.Provide(
  fx.Annotate(
    NewGateway,
    fx.ParamTags(`name:"ro" optional:"true"`, `name:"rw"`),
    fx.ResultTags(`name:"foo"`),
  ),
)

Is equivalent to,

type params struct {
  fx.In

  RO *db.Conn `name:"ro" optional:"true"`
  RW *db.Conn `name:"rw"`
}

type result struct {
  fx.Out

  GW *Gateway `name:"foo"`
}

fx.Provide(func(p params) result {
   return result{GW: NewGateway(p.RO, p.RW)}
})

Using the same annotation multiple times is invalid.
For example, the following will fail with an error:

fx.Provide(
  fx.Annotate(
    NewGateWay,
    fx.ParamTags(`name:"ro" optional:"true"`),
    fx.ParamTags(`name:"rw"), // ERROR: ParamTags was already used above
    fx.ResultTags(`name:"foo"`)
  )
)

If more tags are given than the number of parameters/results, only
the ones up to the number of parameters/results will be applied.

Variadic functions

¶

If the provided function is variadic, Annotate treats its parameter as a
slice. For example,

fx.Annotate(func(w io.Writer, rs ...io.Reader) {
  // ...
}, ...)

Is equivalent to,

fx.Annotate(func(w io.Writer, rs []io.Reader) {
  // ...
}, ...)

You can use variadic parameters with Fx's value groups.
For example,

fx.Annotate(func(mux *http.ServeMux, handlers ...http.Handler) {
  // ...
}, fx.ParamTags(``, `group:"server"`))

If we provide the above to the application,
any constructor in the Fx application can inject its HTTP handlers
by using

Annotate

,

Annotated

, or

Out

.

fx.Annotate(
  func(..) http.Handler { ... },
  fx.ResultTags(`group:"server"`),
)

fx.Annotated{
  Target: func(..) http.Handler { ... },
  Group:  "server",
}

func

Self

¶

added in

v1.22.0

func Self()

any

Self returns a special value that can be passed to

As

to indicate
that a type should be provided as its original type, in addition to whatever other
types it gets provided as via other

As

annotations.

For example,

fx.Provide(
  fx.Annotate(
    bytes.NewBuffer,
    fx.As(new(io.Writer)),
    fx.As(fx.Self()),
  )
)

Is equivalent to,

fx.Provide(
  bytes.NewBuffer,
  func(b *bytes.Buffer) io.Writer {
    return b
  },
)

in that it provides the same *bytes.Buffer instance
as both a *bytes.Buffer and an io.Writer.

func

ValidateApp

¶

added in

v1.13.0

func ValidateApp(opts ...

Option

)

error

ValidateApp validates that supplied graph would run and is not missing any dependencies. This
method does not invoke actual input functions.

func

VisualizeError

¶

added in

v1.7.0

func VisualizeError(err

error

) (

string

,

error

)

VisualizeError returns the visualization of the error if available.

Note that VisualizeError does not yet recognize

Decorate

and

Replace

.

Types

¶

type

Annotated

¶

added in

v1.9.0

type Annotated struct {

// If specified, this will be used as the name for all non-error values returned

// by the constructor. For more information on named values, see the documentation

// for the fx.Out type.

//

// A name option may not be provided if a group option is provided.

Name

string

// If specified, this will be used as the group name for all non-error values returned

// by the constructor. For more information on value groups, see the package documentation.

//

// A group option may not be provided if a name option is provided.

//

// Similar to group tags, the group name may be followed by a `,flatten`

// option to indicate that each element in the slice returned by the

// constructor should be injected into the value group individually.

Group

string

// Target is the constructor or value being annotated with fx.Annotated.

Target interface{}
}

Annotated annotates a constructor provided to Fx with additional options.

For example,

func NewReadOnlyConnection(...) (*Connection, error)

fx.Provide(fx.Annotated{
  Name: "ro",
  Target: NewReadOnlyConnection,
})

Is equivalent to,

type result struct {
  fx.Out

  Connection *Connection `name:"ro"`
}

fx.Provide(func(...) (result, error) {
  conn, err := NewReadOnlyConnection(...)
  return result{Connection: conn}, err
})

Annotated cannot be used with constructors which produce fx.Out objects.
When used with

Supply

, Target is a value instead of a constructor.

This type represents a less powerful version of the

Annotate

construct;
prefer

Annotate

where possible.

func (Annotated)

String

¶

added in

v1.10.0

func (a

Annotated

) String()

string

type

Annotation

¶

added in

v1.15.0

type Annotation interface {

// contains filtered or unexported methods

}

Annotation specifies how to wrap a target for

Annotate

.
It can be used to set up additional options for a constructor,
or with

Supply

, for a value.

func

As

¶

added in

v1.15.0

func As(interfaces ...interface{})

Annotation

As is an Annotation that annotates the result of a function (i.e. a
constructor) to be provided as another interface.

For example, the following code specifies that the return type of
bytes.NewBuffer (bytes.Buffer) should be provided as io.Writer type:

fx.Provide(
  fx.Annotate(bytes.NewBuffer, fx.As(new(io.Writer)))
)

In other words, the code above is equivalent to:

fx.Provide(func() io.Writer {
  return bytes.NewBuffer()
  // provides io.Writer instead of *bytes.Buffer
})

Note that the bytes.Buffer type is provided as an io.Writer type, so this
constructor does NOT provide both bytes.Buffer and io.Writer type; it just
provides io.Writer type.

When multiple values are returned by the annotated function, each type
gets mapped to corresponding positional result of the annotated function.

For example,

func a() (bytes.Buffer, bytes.Buffer) {
  ...
}
fx.Provide(
  fx.Annotate(a, fx.As(new(io.Writer), new(io.Reader)))
)

Is equivalent to,

fx.Provide(func() (io.Writer, io.Reader) {
  w, r := a()
  return w, r
}

As entirely replaces the default return types of a function. In order
to maintain the original return types when using As, see

Self

.

As annotation cannot be used in a function that returns an

Out

struct as a return type.

func

From

¶

added in

v1.19.0

func From(interfaces ...interface{})

Annotation

From is an

Annotation

that annotates the parameter(s) for a function (i.e. a
constructor) to be accepted from other provided types. It is analogous to the

As

for parameter types to the constructor.

For example,

type Runner interface { Run() }
func NewFooRunner() *FooRunner // implements Runner
func NewRunnerWrap(r Runner) *RunnerWrap

fx.Provide(
  fx.Annotate(
    NewRunnerWrap,
    fx.From(new(*FooRunner)),
  ),
)

Is equivalent to,

fx.Provide(func(r *FooRunner) *RunnerWrap {
  // need *FooRunner instead of Runner
  return NewRunnerWrap(r)
})

When the annotated function takes in multiple parameters, each type gets
mapped to corresponding positional parameter of the annotated function

For example,

func NewBarRunner() *BarRunner // implements Runner
func NewRunnerWraps(r1 Runner, r2 Runner) *RunnerWraps

fx.Provide(
  fx.Annotate(
    NewRunnerWraps,
    fx.From(new(*FooRunner), new(*BarRunner)),
  ),
)

Is equivalent to,

fx.Provide(func(r1 *FooRunner, r2 *BarRunner) *RunnerWraps {
  return NewRunnerWraps(r1, r2)
})

From annotation cannot be used in a function that takes an

In

struct as a
parameter.

func

OnStart

¶

added in

v1.18.0

func OnStart(onStart interface{})

Annotation

OnStart is an Annotation that appends an OnStart Hook to the application
Lifecycle when that function is called. This provides a way to create
Lifecycle OnStart (see Lifecycle type documentation) hooks without building a
function that takes a dependency on the Lifecycle type.

fx.Provide(
	fx.Annotate(
		NewServer,
		fx.OnStart(func(ctx context.Context, server Server) error {
			return server.Listen(ctx)
		}),
	)
)

Which is functionally the same as:

fx.Provide(
   func(lifecycle fx.Lifecycle, p Params) Server {
     server := NewServer(p)
     lifecycle.Append(fx.Hook{
	      OnStart: func(ctx context.Context) error {
		    return server.Listen(ctx)
	      },
     })
	 return server
   }
 )

It is also possible to use OnStart annotation with other parameter and result
annotations, provided that the parameter of the function passed to OnStart
matches annotated parameters and results.

For example, the following is possible:

fx.Provide(
	fx.Annotate(
		func (a A) B {...},
		fx.ParamTags(`name:"A"`),
		fx.ResultTags(`name:"B"`),
		fx.OnStart(func (p OnStartParams) {...}),
	),
)

As long as OnStartParams looks like the following and has no other dependencies
besides Context or Lifecycle:

type OnStartParams struct {
	fx.In
	FieldA A `name:"A"`
	FieldB B `name:"B"`
}

Only one OnStart annotation may be applied to a given function at a time,
however functions may be annotated with other types of lifecycle Hooks, such
as OnStop. The hook function passed into OnStart cannot take any arguments
outside of the annotated constructor's existing dependencies or results, except
a context.Context.

func

OnStop

¶

added in

v1.18.0

func OnStop(onStop interface{})

Annotation

OnStop is an Annotation that appends an OnStop Hook to the application
Lifecycle when that function is called. This provides a way to create
Lifecycle OnStop (see Lifecycle type documentation) hooks without building a
function that takes a dependency on the Lifecycle type.

fx.Provide(
	fx.Annotate(
		NewServer,
		fx.OnStop(func(ctx context.Context, server Server) error {
			return server.Shutdown(ctx)
		}),
	)
)

Which is functionally the same as:

fx.Provide(
   func(lifecycle fx.Lifecycle, p Params) Server {
     server := NewServer(p)
     lifecycle.Append(fx.Hook{
	      OnStop: func(ctx context.Context) error {
		    return server.Shutdown(ctx)
	      },
     })
	 return server
   }
 )

It is also possible to use OnStop annotation with other parameter and result
annotations, provided that the parameter of the function passed to OnStop
matches annotated parameters and results.

For example, the following is possible:

fx.Provide(
	fx.Annotate(
		func (a A) B {...},
		fx.ParamTags(`name:"A"`),
		fx.ResultTags(`name:"B"`),
		fx.OnStop(func (p OnStopParams) {...}),
	),
)

As long as OnStopParams looks like the following and has no other dependencies
besides Context or Lifecycle:

type OnStopParams struct {
	fx.In
	FieldA A `name:"A"`
	FieldB B `name:"B"`
}

Only one OnStop annotation may be applied to a given function at a time,
however functions may be annotated with other types of lifecycle Hooks, such
as OnStart. The hook function passed into OnStop cannot take any arguments
outside of the annotated constructor's existing dependencies or results, except
a context.Context.

func

ParamTags

¶

added in

v1.15.0

func ParamTags(tags ...

string

)

Annotation

ParamTags is an Annotation that annotates the parameter(s) of a function.

When multiple tags are specified, each tag is mapped to the corresponding
positional parameter.
For example, the following will refer to a named database connection,
and the default, unnamed logger:

fx.Annotate(func(log *log.Logger, conn *sql.DB) *Handler {
	// ...
}, fx.ParamTags("", `name:"ro"`))

ParamTags cannot be used in a function that takes an fx.In struct as a
parameter.

func

ResultTags

¶

added in

v1.15.0

func ResultTags(tags ...

string

)

Annotation

ResultTags is an Annotation that annotates the result(s) of a function.
When multiple tags are specified, each tag is mapped to the corresponding
positional result.

For example, the following will produce a named database connection.

fx.Annotate(func() (*sql.DB, error) {
	// ...
}, fx.ResultTags(`name:"ro"`))

ResultTags cannot be used on a function that returns an fx.Out struct.

type

App

¶

type App struct {

// contains filtered or unexported fields

}

An App is a modular application built around dependency injection. Most
users will only need to use the New constructor and the all-in-one Run
convenience method. In more unusual cases, users may need to use the Err,
Start, Done, and Stop methods by hand instead of relying on Run.

New

creates and initializes an App. All applications begin with a
constructor for the Lifecycle type already registered.

In addition to that built-in functionality, users typically pass a handful
of

Provide

options and one or more

Invoke

options. The Provide options
teach the application how to instantiate a variety of types, and the Invoke
options describe how to initialize the application.

When created, the application immediately executes all the functions passed
via Invoke options. To supply these functions with the parameters they
need, the application looks for constructors that return the appropriate
types; if constructors for any required types are missing or any
invocations return an error, the application will fail to start (and Err
will return a descriptive error message).

Once all the invocations (and any required constructors) have been called,
New returns and the application is ready to be started using Run or Start.
On startup, it executes any OnStart hooks registered with its Lifecycle.
OnStart hooks are executed one at a time, in order, and must all complete
within a configurable deadline (by default, 15 seconds). For details on the
order in which OnStart hooks are executed, see the documentation for the
Start method.

At this point, the application has successfully started up. If started via
Run, it will continue operating until it receives a shutdown signal from
Done (see the

App.Done

documentation for details); if started explicitly via
Start, it will operate until the user calls Stop. On shutdown, OnStop hooks
execute one at a time, in reverse order, and must all complete within a
configurable deadline (again, 15 seconds by default).

func

New

¶

func New(opts ...

Option

) *

App

New creates and initializes an App, immediately executing any functions
registered via

Invoke

options. See the documentation of the App struct for
details on the application's initialization, startup, and shutdown logic.

func (*App)

Done

¶

func (app *

App

) Done() <-chan

os

.

Signal

Done returns a channel of signals to block on after starting the
application. Applications listen for the SIGINT and SIGTERM signals; during
development, users can send the application SIGTERM by pressing Ctrl-C in
the same terminal as the running process.

Alternatively, a signal can be broadcast to all done channels manually by
using the Shutdown functionality (see the

Shutdowner

documentation for details).

func (*App)

Err

¶

func (app *

App

) Err()

error

Err returns any error encountered during New's initialization. See the
documentation of the New method for details, but typical errors include
missing constructors, circular dependencies, constructor errors, and
invocation errors.

Most users won't need to use this method, since both Run and Start
short-circuit if initialization failed.

func (*App)

Run

¶

func (app *

App

) Run()

Run starts the application, blocks on the signals channel, and then
gracefully shuts the application down. It uses

DefaultTimeout

to set a
deadline for application startup and shutdown, unless the user has
configured different timeouts with the

StartTimeout

or

StopTimeout

options.
It's designed to make typical applications simple to run.
The minimal Fx application looks like this:

fx.New().Run()

All of Run's functionality is implemented in terms of the exported
Start, Done, and Stop methods. Applications with more specialized needs
can use those methods directly instead of relying on Run.

After the application has started,
it can be shut down by sending a signal or calling [Shutdowner.Shutdown].
On successful shutdown, whether initiated by a signal or by the user,
Run will return to the caller, allowing it to exit cleanly.
Run will exit with a non-zero status code
if startup or shutdown operations fail,
or if the

Shutdowner

supplied a non-zero exit code.

func (*App)

Start

¶

func (app *

App

) Start(ctx

context

.

Context

) (err

error

)

Start kicks off all long-running goroutines, like network servers or
message queue consumers. It does this by interacting with the application's
Lifecycle.

By taking a dependency on the Lifecycle type, some of the user-supplied
functions called during initialization may have registered start and stop
hooks. Because initialization calls constructors serially and in dependency
order, hooks are naturally registered in serial and dependency order too.

Start executes all OnStart hooks registered with the application's
Lifecycle, one at a time and in order. This ensures that each constructor's
start hooks aren't executed until all its dependencies' start hooks
complete. If any of the start hooks return an error, Start short-circuits,
calls Stop, and returns the inciting error.

Note that Start short-circuits immediately if the New constructor
encountered any errors in application initialization.

func (*App)

StartTimeout

¶

added in

v1.5.0

func (app *

App

) StartTimeout()

time

.

Duration

StartTimeout returns the configured startup timeout.
This defaults to

DefaultTimeout

, and can be changed with the

StartTimeout

option.

func (*App)

Stop

¶

func (app *

App

) Stop(ctx

context

.

Context

) (err

error

)

Stop gracefully stops the application. It executes any registered OnStop
hooks in reverse order, so that each constructor's stop hooks are called
before its dependencies' stop hooks.

If the application didn't start cleanly, only hooks whose OnStart phase was
called are executed. However, all those hooks are executed, even if some
fail.

func (*App)

StopTimeout

¶

added in

v1.5.0

func (app *

App

) StopTimeout()

time

.

Duration

StopTimeout returns the configured shutdown timeout.
This defaults to

DefaultTimeout

, and can be changed with the

StopTimeout

option.

func (*App)

Wait

¶

added in

v1.19.0

func (app *

App

) Wait() <-chan

ShutdownSignal

Wait returns a channel of

ShutdownSignal

to block on after starting the
application and function, similar to

App.Done

, but with a minor difference:
if the app was shut down via [Shutdowner.Shutdown],
the exit code (if provied via

ExitCode

) will be available
in the

ShutdownSignal

struct.
Otherwise, the signal that was received will be set.

type

DotGraph

¶

added in

v1.7.0

type DotGraph

string

DotGraph contains a DOT language visualization of the dependency graph in
an Fx application. It is provided in the container by default at
initialization. On failure to build the dependency graph, it is attached
to the error and if possible, colorized to highlight the root cause of the
failure.

Note that DotGraph does not yet recognize

Decorate

and

Replace

.

type

ErrorHandler

¶

added in

v1.7.0

type ErrorHandler interface {

HandleError(

error

)

}

ErrorHandler handles Fx application startup errors.
Register these with

ErrorHook

.
If specified, and the application fails to start up,
the failure will still cause a crash,
but you'll have a chance to log the error or take some other action.

type

Hook

¶

type Hook struct {

OnStart func(

context

.

Context

)

error

OnStop  func(

context

.

Context

)

error

// contains filtered or unexported fields

}

A Hook is a pair of start and stop callbacks, either of which can be nil.
If a Hook's OnStart callback isn't executed (because a previous OnStart
failure short-circuited application startup), its OnStop callback won't be
executed.

func

StartHook

¶

added in

v1.19.0

func StartHook[T

HookFunc

](start T)

Hook

StartHook returns a new Hook with start as its [Hook.OnStart] function,
wrapping its signature as needed. For example, given the following function:

func myfunc() {
  fmt.Println("hook called")
}

then calling:

lifecycle.Append(StartHook(myfunc))

is functionally equivalent to calling:

lifecycle.Append(fx.Hook{
  OnStart: func(context.Context) error {
    myfunc()
    return nil
  },
})

The same is true for all functions that satisfy the HookFunc constraint.
Note that any context.Context parameter or error return will be propagated
as expected. If propagation is not intended, users should instead provide a
closure that discards the undesired value(s), or construct a Hook directly.

func

StartStopHook

¶

added in

v1.19.0

func StartStopHook[T1

HookFunc

, T2

HookFunc

](start T1, stop T2)

Hook

StartStopHook returns a new Hook with start as its [Hook.OnStart] function
and stop as its [Hook.OnStop] function, independently wrapping the signature
of each as needed.

func

StopHook

¶

added in

v1.19.0

func StopHook[T

HookFunc

](stop T)

Hook

StopHook returns a new Hook with stop as its [Hook.OnStop] function,
wrapping its signature as needed. For example, given the following function:

func myfunc() {
  fmt.Println("hook called")
}

then calling:

lifecycle.Append(StopHook(myfunc))

is functionally equivalent to calling:

lifecycle.Append(fx.Hook{
  OnStop: func(context.Context) error {
    myfunc()
    return nil
  },
})

The same is true for all functions that satisfy the HookFunc constraint.
Note that any context.Context parameter or error return will be propagated
as expected. If propagation is not intended, users should instead provide a
closure that discards the undesired value(s), or construct a Hook directly.

type

HookFunc

¶

added in

v1.19.0

type HookFunc interface {
	~func() | ~func()

error

| ~func(

context

.

Context

) | ~func(

context

.

Context

)

error

}

A HookFunc is a function that can be used as a

Hook

.

type

In

¶

type In =

dig

.

In

In can be embedded into a struct to mark it as a parameter struct.
This allows it to make use of advanced dependency injection features.
See package documentation for more information.

It's recommended that shared modules use a single parameter struct to
provide a forward-compatible API:
adding new optional fields to a struct is backward-compatible,
so modules can evolve as needs change.

type

Lifecycle

¶

type Lifecycle interface {

Append(

Hook

)

}

Lifecycle allows constructors to register callbacks that are executed on
application start and stop. See the documentation for App for details on Fx
applications' initialization, startup, and shutdown logic.

type

Option

¶

type Option interface {

fmt

.

Stringer

// contains filtered or unexported methods

}

An Option specifies the behavior of the application.
This is the primary means by which you interface with Fx.

Zero or more options are specified at startup with

New

.
Options cannot be changed once an application has been initialized.
Options may be grouped into a single option using the

Options

function.
A group of options providing a logical unit of functionality
may use

Module

to name that functionality
and scope certain operations to within that module.

func

Decorate

¶

added in

v1.17.0

func Decorate(decorators ...interface{})

Option

Decorator functions

Decorator scope

Decorate specifies one or more decorator functions to an Fx application.

Decorator functions

¶

Decorator functions let users augment objects in the graph.
They can take in zero or more dependencies that must be provided to the
application with fx.Provide, and produce one or more values that can be used
by other fx.Provide and fx.Invoke calls.

fx.Decorate(func(log *zap.Logger) *zap.Logger {
  return log.Named("myapp")
})
fx.Invoke(func(log *zap.Logger) {
  log.Info("hello")
  // Output:
  // {"level": "info","logger":"myapp","msg":"hello"}
})

The following decorator accepts multiple dependencies from the graph,
augments and returns one of them.

fx.Decorate(func(log *zap.Logger, cfg *Config) *zap.Logger {
  return log.Named(cfg.Name)
})

Similar to fx.Provide, functions passed to fx.Decorate may optionally return
an error as their last result.
If a decorator returns a non-nil error, it will halt application startup.

fx.Decorate(func(conn *sql.DB, cfg *Config) (*sql.DB, error) {
  if err := conn.Ping(); err != nil {
    return sql.Open("driver-name", cfg.FallbackDB)
  }
  return conn, nil
})

Decorators support both, fx.In and fx.Out structs, similar to fx.Provide and
fx.Invoke.

type Params struct {
  fx.In

  Client usersvc.Client `name:"readOnly"`
}

type Result struct {
  fx.Out

  Client usersvc.Client `name:"readOnly"`
}

fx.Decorate(func(p Params) Result {
  ...
})

Decorators can be annotated with the fx.Annotate function, but not with the
fx.Annotated type. Refer to documentation on fx.Annotate() to learn how to
use it for annotating functions.

fx.Decorate(
  fx.Annotate(
    func(client usersvc.Client) usersvc.Client {
      // ...
    },
    fx.ParamTags(`name:"readOnly"`),
    fx.ResultTags(`name:"readOnly"`),
  ),
)

Decorators support augmenting, filtering, or replacing value groups.
To decorate a value group, expect the entire value group slice and produce
the new slice.

type HandlerParam struct {
  fx.In

  Log      *zap.Logger
  Handlers []Handler `group:"server"
}

type HandlerResult struct {
  fx.Out

  Handlers []Handler `group:"server"
}

fx.Decorate(func(p HandlerParam) HandlerResult {
  var r HandlerResult
  for _, handler := range p.Handlers {
    r.Handlers = append(r.Handlers, wrapWithLogger(p.Log, handler))
  }
  return r
}),

Decorators can not add new values to the graph,
only modify or replace existing ones.
Types returned by a decorator that are not already in the graph
will be ignored.

Decorator scope

¶

Modifications made to the Fx graph with fx.Decorate are scoped to the
deepest fx.Module inside which the decorator was specified.

fx.Module("mymodule",
  fx.Decorate(func(log *zap.Logger) *zap.Logger {
    return log.Named("myapp")
  }),
  fx.Invoke(func(log *zap.Logger) {
    log.Info("decorated logger")
    // Output:
    // {"level": "info","logger":"myapp","msg":"decorated logger"}
  }),
),
fx.Invoke(func(log *zap.Logger) {
  log.Info("plain logger")
  // Output:
  // {"level": "info","msg":"plain logger"}
}),

Decorations specified in the top-level fx.New call apply across the
application and chain with module-specific decorators.

fx.New(
  // ...
  fx.Decorate(func(log *zap.Logger) *zap.Logger {
    return log.With(zap.Field("service", "myservice"))
  }),
  // ...
  fx.Invoke(func(log *zap.Logger) {
    log.Info("outer decorator")
    // Output:
    // {"level": "info","service":"myservice","msg":"outer decorator"}
  }),
  // ...
  fx.Module("mymodule",
    fx.Decorate(func(log *zap.Logger) *zap.Logger {
      return log.Named("myapp")
    }),
    fx.Invoke(func(log *zap.Logger) {
      log.Info("inner decorator")
      // Output:
      // {"level": "info","logger":"myapp","service":"myservice","msg":"inner decorator"}
    }),
  ),
)

func

Error

¶

added in

v1.6.0

func Error(errs ...

error

)

Option

Error registers any number of errors with the application to short-circuit
startup. If more than one error is given, the errors are combined into a
single error.

Similar to invocations, errors are applied in order. All Provide and Invoke
options registered before or after an Error option will not be applied.

Example

¶

package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"go.uber.org/fx"
)

func main() {
	// A module that provides a HTTP server depends on
	// the $PORT environment variable. If the variable
	// is unset, the module returns an fx.Error option.
	newHTTPServer := func() fx.Option {
		port := os.Getenv("PORT")
		if port == "" {
			return fx.Error(errors.New("$PORT is not set"))
		}
		return fx.Provide(&http.Server{
			Addr: fmt.Sprintf("127.0.0.1:%s", port),
		})
	}

	app := fx.New(
		fx.NopLogger,
		newHTTPServer(),
		fx.Invoke(func(s *http.Server) error { return s.ListenAndServe() }),
	)

	fmt.Println(app.Err())

}

Output:

$PORT is not set

Share

Format

Run

func

ErrorHook

¶

added in

v1.7.0

func ErrorHook(funcs ...

ErrorHandler

)

Option

ErrorHook registers error handlers that implement error handling functions.
They are executed on invoke failures. Passing multiple ErrorHandlers appends
the new handlers to the application's existing list.

func

Extract

deprecated

func Extract(target interface{})

Option

Extract fills the given struct with values from the dependency injection
container on application initialization. The target MUST be a pointer to a
struct. Only exported fields will be filled.

Deprecated: Use Populate instead.

func

Invoke

¶

func Invoke(funcs ...interface{})

Option

Invoke registers functions that are executed eagerly on application start.
Arguments for these invocations are built using the constructors registered
by Provide. Passing multiple Invoke options appends the new invocations to
the application's existing list.

Unlike constructors, invocations are always executed, and they're always
run in order. Invocations may have any number of returned values.
If the final returned object is an error, it indicates whether the operation
was successful.
All other returned values are discarded.

Invokes registered in [Module]s are run before the ones registered at the
scope of the parent. Invokes within the same Module is run in the order
they were provided. For example,

fx.New(
	fx.Invoke(func3),
	fx.Module("someModule",
		fx.Invoke(func1),
		fx.Invoke(func2),
	),
	fx.Invoke(func4),
)

invokes func1, func2, func3, func4 in that order.

Typically, invoked functions take a handful of high-level objects (whose
constructors depend on lower-level objects) and introduce them to each
other. This kick-starts the application by forcing it to instantiate a
variety of types.

To see an invocation in use, read through the package-level example. For
advanced features, including optional parameters and named instances, see
the documentation of the In and Out types.

func

Logger

¶

func Logger(p

Printer

)

Option

Logger redirects the application's log output to the provided printer.

Prefer to use

WithLogger

instead.

func

Module

¶

added in

v1.17.0

func Module(name

string

, opts ...

Option

)

Option

Module is a named group of zero or more fx.Options.

A Module scopes the effect of certain operations to within the module.
For more information, see

Decorate

,

Replace

, or

Invoke

.

Module allows packages to bundle sophisticated functionality into easy-to-use
logical units.
For example, a logging package might export a simple option like this:

package logging

var Module = fx.Module("logging",
	fx.Provide(func() *log.Logger {
		return log.New(os.Stdout, "", 0)
	}),
	// ...
)

A shared all-in-one microservice package could use Module to bundle
all required components of a microservice:

package server

var Module = fx.Module("server",
	logging.Module,
	metrics.Module,
	tracing.Module,
	rpc.Module,
)

When new global functionality is added to the service ecosystem,
it can be added to the shared module with minimal churn for users.

Use the all-in-one pattern sparingly.
It limits the flexibility available to the application.

func

Options

¶

func Options(opts ...

Option

)

Option

Options bundles a group of options together into a single option.

Use Options to group together options that don't belong in a

Module

.

var loggingAndMetrics = fx.Options(
	logging.Module,
	metrics.Module,
	fx.Invoke(func(logger *log.Logger) {
		app.globalLogger = logger
	}),
)

func

Populate

¶

added in

v1.4.0

func Populate(targets ...interface{})

Option

Populate sets targets with values from the dependency injection container
during application initialization. All targets must be pointers to the
values that must be populated. Pointers to structs that embed In are
supported, which can be used to populate multiple values in a struct.

Annotating each pointer with ParamTags is also supported as a shorthand
to passing a pointer to a struct that embeds In with field tags. For example:

var a A
 var b B
 fx.Populate(
	fx.Annotate(
			&a,
			fx.ParamTags(`name:"A"`)
 	),
	fx.Annotate(
			&b,
			fx.ParamTags(`name:"B"`)
 	)
 )

Code above is equivalent to the following:

type Target struct {
	fx.In

	a A `name:"A"`
	b B `name:"B"`
}
var target Target
...
fx.Populate(&target)

This is most helpful in unit tests: it lets tests leverage Fx's automatic
constructor wiring to build a few structs, but then extract those structs
for further testing.

Example

¶

package main

import (
	"context"
	"fmt"

	"go.uber.org/fx"
)

func main() {
	// Some external module that provides a user name.
	type Username string
	UserModule := fx.Provide(func() Username { return "john" })

	// We want to use Fx to wire up our constructors, but don't actually want to
	// run the application - we just want to yank out the user name.
	//
	// This is common in unit tests, and is even easier with the fxtest
	// package's RequireStart and RequireStop helpers.
	var user Username
	app := fx.New(
		UserModule,
		fx.NopLogger, // silence test output
		fx.Populate(&user),
	)
	if err := app.Start(context.Background()); err != nil {
		panic(err)
	}
	defer app.Stop(context.Background())

	fmt.Println(user)

}

Output:

john

Share

Format

Run

func

Provide

¶

func Provide(constructors ...interface{})

Option

Provide registers any number of constructor functions, teaching the
application how to instantiate various types. The supplied constructor
function(s) may depend on other types available in the application, must
return one or more objects, and may return an error. For example:

// Constructs type *C, depends on *A and *B.
func(*A, *B) *C

// Constructs type *C, depends on *A and *B, and indicates failure by
// returning an error.
func(*A, *B) (*C, error)

// Constructs types *B and *C, depends on *A, and can fail.
func(*A) (*B, *C, error)

The order in which constructors are provided doesn't matter, and passing
multiple Provide options appends to the application's collection of
constructors. Constructors are called only if one or more of their returned
types are needed, and their results are cached for reuse (so instances of a
type are effectively singletons within an application). Taken together,
these properties make it perfectly reasonable to Provide a large number of
constructors even if only a fraction of them are used.

See the documentation of the In and Out types for advanced features,
including optional parameters and named instances.

See the documentation for

Private

for restricting access to constructors.

Constructor functions should perform as little external interaction as
possible, and should avoid spawning goroutines. Things like server listen
loops, background timer loops, and background processing goroutines should
instead be managed using Lifecycle callbacks.

func

RecoverFromPanics

¶

added in

v1.19.0

func RecoverFromPanics()

Option

RecoverFromPanics causes panics that occur in functions given to

Provide

,

Decorate

, and

Invoke

to be recovered from.
This error can be retrieved as any other error, by using (*App).Err().

func

Replace

¶

added in

v1.17.0

func Replace(values ...interface{})

Option

Replace Caveats

Replace provides instantiated values for graph modification as if
they had been provided using a decorator with fx.Decorate.
The most specific type of each value (as determined by reflection) is used.

Refer to the documentation on fx.Decorate to see how graph modifications
work with fx.Module.

This serves a purpose similar to what fx.Supply does for fx.Provide.

For example, given,

var log *zap.Logger = ...

The following two forms are equivalent.

fx.Replace(log)

fx.Decorate(
	func() *zap.Logger {
		return log
	},
)

Replace panics if a value (or annotation target) is an untyped nil or an error.

Replace Caveats

¶

As mentioned above, Replace uses the most specific type of the provided
value. For interface values, this refers to the type of the implementation,
not the interface. So if you try to replace an io.Writer, fx.Replace will
use the type of the implementation.

var stderr io.Writer = os.Stderr
fx.Replace(stderr)

Is equivalent to,

fx.Decorate(func() *os.File { return os.Stderr })

This is typically NOT what you intended. To replace the io.Writer in the
container with the value above, we need to use the fx.Annotate function with
the fx.As annotation.

fx.Replace(
	fx.Annotate(os.Stderr, fx.As(new(io.Writer)))
)

func

StartTimeout

¶

added in

v1.5.0

func StartTimeout(v

time

.

Duration

)

Option

StartTimeout changes the application's start timeout.
This controls the total time that all

OnStart

hooks have to complete.
If the timeout is exceeded, the application will fail to start.

Defaults to

DefaultTimeout

.

func

StopTimeout

¶

added in

v1.5.0

func StopTimeout(v

time

.

Duration

)

Option

StopTimeout changes the application's stop timeout.
This controls the total time that all

OnStop

hooks have to complete.
If the timeout is exceeded, the application will exit early.

Defaults to

DefaultTimeout

.

func

Supply

¶

added in

v1.12.0

func Supply(values ...interface{})

Option

Supply Caveats

Supply provides instantiated values for dependency injection as if
they had been provided using a constructor that simply returns them.
The most specific type of each value (as determined by reflection) is used.

This serves a purpose similar to what fx.Replace does for fx.Decorate.

For example, given:

type (
	TypeA struct{}
	TypeB struct{}
	TypeC struct{}
)

var a, b, c = &TypeA{}, TypeB{}, &TypeC{}

The following two forms are equivalent:

fx.Supply(a, b, fx.Annotated{Target: c})

fx.Provide(
	func() *TypeA { return a },
	func() TypeB { return b },
	fx.Annotated{Target: func() *TypeC { return c }},
)

Supply panics if a value (or annotation target) is an untyped nil or an error.

Private

can be used to restrict access to supplied values.

Supply Caveats

¶

As mentioned above, Supply uses the most specific type of the provided
value. For interface values, this refers to the type of the implementation,
not the interface. So if you supply an http.Handler, fx.Supply will use the
type of the implementation.

var handler http.Handler = http.HandlerFunc(f)
fx.Supply(handler)

Is equivalent to,

fx.Provide(func() http.HandlerFunc { return f })

This is typically NOT what you intended. To supply the handler above as an
http.Handler, we need to use the fx.Annotate function with the fx.As
annotation.

fx.Supply(
	fx.Annotate(handler, fx.As(new(http.Handler))),
)

func

WithLogger

¶

added in

v1.14.0

func WithLogger(constructor interface{})

Option

WithLogger specifies the

fxevent.Logger

used by Fx to log its own events
(e.g. a constructor was provided, a function was invoked, etc.).

The argument to this is a constructor with one of the following return
types:

fxevent.Logger
(fxevent.Logger, error)

The constructor may depend on any other types provided to the application.
For example,

WithLogger(func(logger *zap.Logger) fxevent.Logger {
  return &fxevent.ZapLogger{Logger: logger}
})

If specified, Fx will construct the logger and log all its events to the
specified logger.

If Fx fails to build the logger, or no logger is specified, it will fall back to

fxevent.ConsoleLogger

configured to write to stderr.

type

Out

¶

type Out =

dig

.

Out

Out is the inverse of In: it marks a struct as a result struct so that
it can be used with advanced dependency injection features.
See package documentation for more information.

It's recommended that shared modules use a single result struct to
provide a forward-compatible API:
adding new fields to a struct is backward-compatible,
so modules can produce more outputs as they grow.

type

Printer

¶

type Printer interface {

Printf(

string

, ...interface{})

}

Printer is the interface required by Fx's logging backend. It's implemented
by most loggers, including the one bundled with the standard library.

Note, this will be deprecated in a future release.
Prefer to use

fxevent.Logger

instead.

type

ShutdownOption

¶

added in

v1.9.0

type ShutdownOption interface {

// contains filtered or unexported methods

}

ShutdownOption provides a way to configure properties of the shutdown
process. Currently, no options have been implemented.

func

ExitCode

¶

added in

v1.19.0

func ExitCode(code

int

)

ShutdownOption

ExitCode is a

ShutdownOption

that may be passed to the Shutdown method of the

Shutdowner

interface.
The given integer exit code will be broadcasted to any receiver waiting
on a

ShutdownSignal

from the [Wait] method.

func

ShutdownTimeout

deprecated

added in

v1.19.0

func ShutdownTimeout(timeout

time

.

Duration

)

ShutdownOption

ShutdownTimeout is a

ShutdownOption

that allows users to specify a timeout
for a given call to Shutdown method of the

Shutdowner

interface. As the
Shutdown method will block while waiting for a signal receiver relay
goroutine to stop.

Deprecated: This option has no effect. Shutdown is not a blocking operation.

type

ShutdownSignal

¶

added in

v1.19.0

type ShutdownSignal struct {

Signal

os

.

Signal

ExitCode

int

}

ShutdownSignal represents a signal to be written to Wait or Done.
Should a user call the Shutdown method via the Shutdowner interface with
a provided ExitCode, that exit code will be populated in the ExitCode field.

Should the application receive an operating system signal,
the Signal field will be populated with the received os.Signal.

func (ShutdownSignal)

String

¶

added in

v1.19.0

func (sig

ShutdownSignal

) String()

string

String will render a ShutdownSignal type as a string suitable for printing.

type

Shutdowner

¶

added in

v1.9.0

type Shutdowner interface {

Shutdown(...

ShutdownOption

)

error

}

Shutdowner provides a method that can manually trigger the shutdown of the
application by sending a signal to all open Done channels. Shutdowner works
on applications using Run as well as Start, Done, and Stop. The Shutdowner is
provided to all Fx applications.