# slog-multi

> Source: https://pkg.go.dev/github.com/samber/slog-multi
> Fetched: 2026-01-31T11:02:47.410231+00:00
> Content-Hash: f0e745219473fd6d
> Type: html

---

### Index ¶

  * func AttrKindIs(args ...any) func(ctx context.Context, r slog.Record) bool
  * func AttrValueIs(args ...any) func(ctx context.Context, r slog.Record) bool
  * func Failover() func(...slog.Handler) slog.Handler
  * func Fanout(handlers ...slog.Handler) slog.Handler
  * func LevelIs(levels ...slog.Level) func(ctx context.Context, r slog.Record) bool
  * func LevelIsNot(levels ...slog.Level) func(ctx context.Context, r slog.Record) bool
  * func MessageContains(part string) func(ctx context.Context, r slog.Record) bool
  * func MessageIs(msg string) func(ctx context.Context, r slog.Record) bool
  * func MessageIsNot(msg string) func(ctx context.Context, r slog.Record) bool
  * func MessageNotContains(part string) func(ctx context.Context, r slog.Record) bool
  * func NewHandleInlineHandler(...) slog.Handler
  * func NewInlineHandler(...) slog.Handler
  * func Pool() func(...slog.Handler) slog.Handler
  * func RecoverHandlerError(recovery RecoveryFunc) func(slog.Handler) slog.Handler
  * func Router() *router
  * type EnabledInlineMiddleware
  *     * func (h *EnabledInlineMiddleware) Enabled(ctx context.Context, level slog.Level) bool
    * func (h *EnabledInlineMiddleware) Handle(ctx context.Context, record slog.Record) error
    * func (h *EnabledInlineMiddleware) WithAttrs(attrs []slog.Attr) slog.Handler
    * func (h *EnabledInlineMiddleware) WithGroup(name string) slog.Handler
  * type FailoverHandler
  *     * func (h *FailoverHandler) Enabled(ctx context.Context, l slog.Level) bool
    * func (h *FailoverHandler) Handle(ctx context.Context, r slog.Record) error
    * func (h *FailoverHandler) WithAttrs(attrs []slog.Attr) slog.Handler
    * func (h *FailoverHandler) WithGroup(name string) slog.Handler
  * type FanoutHandler
  *     * func (h *FanoutHandler) Enabled(ctx context.Context, l slog.Level) bool
    * func (h *FanoutHandler) Handle(ctx context.Context, r slog.Record) error
    * func (h *FanoutHandler) WithAttrs(attrs []slog.Attr) slog.Handler
    * func (h *FanoutHandler) WithGroup(name string) slog.Handler
  * type FirstMatchHandler
  *     * func FirstMatch(handlers ...*RoutableHandler) *FirstMatchHandler
  *     * func (h *FirstMatchHandler) Enabled(ctx context.Context, l slog.Level) bool
    * func (h *FirstMatchHandler) Handle(ctx context.Context, r slog.Record) error
    * func (h *FirstMatchHandler) WithAttrs(attrs []slog.Attr) slog.Handler
    * func (h *FirstMatchHandler) WithGroup(name string) slog.Handler
  * type HandleInlineHandler
  *     * func (h *HandleInlineHandler) Enabled(ctx context.Context, level slog.Level) bool
    * func (h *HandleInlineHandler) Handle(ctx context.Context, record slog.Record) error
    * func (h *HandleInlineHandler) WithAttrs(attrs []slog.Attr) slog.Handler
    * func (h *HandleInlineHandler) WithGroup(name string) slog.Handler
  * type HandleInlineMiddleware
  *     * func (h *HandleInlineMiddleware) Enabled(ctx context.Context, level slog.Level) bool
    * func (h *HandleInlineMiddleware) Handle(ctx context.Context, record slog.Record) error
    * func (h *HandleInlineMiddleware) WithAttrs(attrs []slog.Attr) slog.Handler
    * func (h *HandleInlineMiddleware) WithGroup(name string) slog.Handler
  * type HandlerErrorRecovery
  *     * func (h *HandlerErrorRecovery) Enabled(ctx context.Context, l slog.Level) bool
    * func (h *HandlerErrorRecovery) Handle(ctx context.Context, record slog.Record) error
    * func (h *HandlerErrorRecovery) WithAttrs(attrs []slog.Attr) slog.Handler
    * func (h *HandlerErrorRecovery) WithGroup(name string) slog.Handler
  * type InlineHandler
  *     * func (h *InlineHandler) Enabled(ctx context.Context, level slog.Level) bool
    * func (h *InlineHandler) Handle(ctx context.Context, record slog.Record) error
    * func (h *InlineHandler) WithAttrs(attrs []slog.Attr) slog.Handler
    * func (h *InlineHandler) WithGroup(name string) slog.Handler
  * type InlineMiddleware
  *     * func (h *InlineMiddleware) Enabled(ctx context.Context, level slog.Level) bool
    * func (h *InlineMiddleware) Handle(ctx context.Context, record slog.Record) error
    * func (h *InlineMiddleware) WithAttrs(attrs []slog.Attr) slog.Handler
    * func (h *InlineMiddleware) WithGroup(name string) slog.Handler
  * type Middleware
  *     * func NewEnabledInlineMiddleware(enabledFunc func(ctx context.Context, level slog.Level, ...) bool) Middleware
    * func NewHandleInlineMiddleware(handleFunc func(ctx context.Context, record slog.Record, ...) error) Middleware
    * func NewInlineMiddleware(enabledFunc func(ctx context.Context, level slog.Level, ...) bool, ...) Middleware
    * func NewWithAttrsInlineMiddleware(...) Middleware
    * func NewWithGroupInlineMiddleware(withGroupFunc func(name string, next func(string) slog.Handler) slog.Handler) Middleware
  * type PipeBuilder
  *     * func Pipe(middlewares ...Middleware) *PipeBuilder
  *     * func (h *PipeBuilder) Handler(handler slog.Handler) slog.Handler
    * func (h *PipeBuilder) Pipe(middleware Middleware) *PipeBuilder
  * type PoolHandler
  *     * func (h *PoolHandler) Enabled(ctx context.Context, l slog.Level) bool
    * func (h *PoolHandler) Handle(ctx context.Context, r slog.Record) error
    * func (h *PoolHandler) WithAttrs(attrs []slog.Attr) slog.Handler
    * func (h *PoolHandler) WithGroup(name string) slog.Handler
  * type RecoveryFunc
  * type RoutableHandler
  *     * func (h *RoutableHandler) Enabled(ctx context.Context, l slog.Level) bool
    * func (h *RoutableHandler) Handle(ctx context.Context, r slog.Record) error
    * func (h *RoutableHandler) WithAttrs(attrs []slog.Attr) slog.Handler
    * func (h *RoutableHandler) WithGroup(name string) slog.Handler
  * type WithAttrsInlineMiddleware
  *     * func (h *WithAttrsInlineMiddleware) Enabled(ctx context.Context, level slog.Level) bool
    * func (h *WithAttrsInlineMiddleware) Handle(ctx context.Context, record slog.Record) error
    * func (h *WithAttrsInlineMiddleware) WithAttrs(attrs []slog.Attr) slog.Handler
    * func (h *WithAttrsInlineMiddleware) WithGroup(name string) slog.Handler
  * type WithGroupInlineMiddleware
  *     * func (h *WithGroupInlineMiddleware) Enabled(ctx context.Context, level slog.Level) bool
    * func (h *WithGroupInlineMiddleware) Handle(ctx context.Context, record slog.Record) error
    * func (h *WithGroupInlineMiddleware) WithAttrs(attrs []slog.Attr) slog.Handler
    * func (h *WithGroupInlineMiddleware) WithGroup(name string) slog.Handler



### Constants ¶

This section is empty.

### Variables ¶

This section is empty.

### Functions ¶

####  func [AttrKindIs](https://github.com/samber/slog-multi/blob/v1.7.0/router_predicate.go#L205) ¶ added in v1.7.0
    
    
    func AttrKindIs(args ...[any](/builtin#any)) func(ctx [context](/context).[Context](/context#Context), r [slog](/log/slog).[Record](/log/slog#Record)) [bool](/builtin#bool)

AttrKindIs returns a function that checks if the record has an attribute with the given key and type. Example usage: 
    
    
    r := slogmulti.Router().
        Add(consoleHandler, slogmulti.AttrKindIs("user_id", slog.KindString)).
        Add(fileHandler, slogmulti.AttrKindIs("user_id", slog.KindString)).
        Handler()
    

Args: 
    
    
    key: The attribute key to check
    ty: The attribute type to check
    

Returns: 
    
    
    A function that checks if the record has an attribute with the given key and type
    

####  func [AttrValueIs](https://github.com/samber/slog-multi/blob/v1.7.0/router_predicate.go#L160) ¶ added in v1.7.0
    
    
    func AttrValueIs(args ...[any](/builtin#any)) func(ctx [context](/context).[Context](/context#Context), r [slog](/log/slog).[Record](/log/slog#Record)) [bool](/builtin#bool)

AttrValueIs returns a function that checks if the record has all specified attributes with exact values. Example usage: 
    
    
    r := slogmulti.Router().
        Add(consoleHandler, slogmulti.AttrValueIs("scope", "influx")).
        Add(fileHandler, slogmulti.AttrValueIs("env", "production", "region", "us-east")).
        Handler()
    

Args: 
    
    
    args: Pairs of attribute key (string) and expected value (any)
    

Returns: 
    
    
    A function that checks if the record has all specified attributes with exact values
    

####  func [Failover](https://github.com/samber/slog-multi/blob/v1.7.0/failover.go#L41) ¶ added in v0.4.0
    
    
    func Failover() func(...[slog](/log/slog).[Handler](/log/slog#Handler)) [slog](/log/slog).[Handler](/log/slog#Handler)

Failover creates a failover handler factory function. This function returns a closure that can be used to create failover handlers with different sets of handlers. 

Example usage: 
    
    
    handler := slogmulti.Failover()(
        primaryHandler,   // First choice
        secondaryHandler, // Fallback if primary fails
        backupHandler,    // Last resort
    )
    logger := slog.New(handler)
    

Returns: 
    
    
    A function that creates FailoverHandler instances with the provided handlers
    

####  func [Fanout](https://github.com/samber/slog-multi/blob/v1.7.0/multi.go#L43) ¶ added in v0.3.0
    
    
    func Fanout(handlers ...[slog](/log/slog).[Handler](/log/slog#Handler)) [slog](/log/slog).[Handler](/log/slog#Handler)

Fanout creates a new FanoutHandler that distributes records to multiple slog.Handler instances. If exactly one handler is provided, it returns that handler unmodified. If you pass a FanoutHandler as an argument, its handlers are flattened into the new FanoutHandler. This function is the primary entry point for creating a multi-handler setup. 

Example usage: 
    
    
    handler := slogmulti.Fanout(
        slog.NewJSONHandler(os.Stdout, nil),
        slogdatadog.NewDatadogHandler(...),
    )
    logger := slog.New(handler)
    

Args: 
    
    
    handlers: Variable number of slog.Handler instances to distribute logs to
    

Returns: 
    
    
    A slog.Handler that forwards all operations to the provided handlers
    

####  func [LevelIs](https://github.com/samber/slog-multi/blob/v1.7.0/router_predicate.go#L24) ¶ added in v1.5.0
    
    
    func LevelIs(levels ...[slog](/log/slog).[Level](/log/slog#Level)) func(ctx [context](/context).[Context](/context#Context), r [slog](/log/slog).[Record](/log/slog#Record)) [bool](/builtin#bool)

LevelIs returns a function that checks if the record level is in the given levels. Example usage: 
    
    
    r := slogmulti.Router().
        Add(consoleHandler, slogmulti.LevelIs(slog.LevelInfo)).
        Add(fileHandler, slogmulti.LevelIs(slog.LevelError)).
        Handler()
    

Args: 
    
    
    levels: The levels to match
    

Returns: 
    
    
    A function that checks if the record level is in the given levels
    

####  func [LevelIsNot](https://github.com/samber/slog-multi/blob/v1.7.0/router_predicate.go#L50) ¶ added in v1.5.0
    
    
    func LevelIsNot(levels ...[slog](/log/slog).[Level](/log/slog#Level)) func(ctx [context](/context).[Context](/context#Context), r [slog](/log/slog).[Record](/log/slog#Record)) [bool](/builtin#bool)

LevelIsNot returns a function that checks if the record level is not in the given levels. Example usage: 
    
    
    r := slogmulti.Router().
        Add(consoleHandler, slogmulti.LevelIsNot(slog.LevelInfo)).
        Add(fileHandler, slogmulti.LevelIsNot(slog.LevelError)).
        Handler()
    

Args: 
    
    
    levels: The levels to check
    

Returns: 
    
    
    A function that checks if the record level is not in the given levels
    

####  func [MessageContains](https://github.com/samber/slog-multi/blob/v1.7.0/router_predicate.go#L118) ¶ added in v1.5.0
    
    
    func MessageContains(part [string](/builtin#string)) func(ctx [context](/context).[Context](/context#Context), r [slog](/log/slog).[Record](/log/slog#Record)) [bool](/builtin#bool)

MessageContains returns a function that checks if the record message contains the given part. Example usage: 
    
    
    r := slogmulti.Router().
        Add(consoleHandler, slogmulti.MessageContains("database error")).
        Add(fileHandler, slogmulti.MessageContains("database error")).
        Handler()
    

Args: 
    
    
    part: The part to check
    

Returns: 
    
    
    A function that checks if the record message contains the given part
    

####  func [MessageIs](https://github.com/samber/slog-multi/blob/v1.7.0/router_predicate.go#L76) ¶ added in v1.5.0
    
    
    func MessageIs(msg [string](/builtin#string)) func(ctx [context](/context).[Context](/context#Context), r [slog](/log/slog).[Record](/log/slog#Record)) [bool](/builtin#bool)

MessageIs returns a function that checks if the record message is equal to the given message. Example usage: 
    
    
    r := slogmulti.Router().
        Add(consoleHandler, slogmulti.MessageIs("database error")).
        Add(fileHandler, slogmulti.MessageIs("database error")).
        Handler()
    

Args: 
    
    
    msg: The message to check
    

Returns: 
    
    
    A function that checks if the record message is equal to the given message
    

####  func [MessageIsNot](https://github.com/samber/slog-multi/blob/v1.7.0/router_predicate.go#L97) ¶ added in v1.5.0
    
    
    func MessageIsNot(msg [string](/builtin#string)) func(ctx [context](/context).[Context](/context#Context), r [slog](/log/slog).[Record](/log/slog#Record)) [bool](/builtin#bool)

MessageIsNot returns a function that checks if the record message is not equal to the given message. Example usage: 
    
    
    r := slogmulti.Router().
        Add(consoleHandler, slogmulti.MessageIsNot("database error")).
        Add(fileHandler, slogmulti.MessageIsNot("database error")).
        Handler()
    

Args: 
    
    
    msg: The message to check
    

Returns: 
    
    
    A function that checks if the record message is not equal to the given message
    

####  func [MessageNotContains](https://github.com/samber/slog-multi/blob/v1.7.0/router_predicate.go#L139) ¶ added in v1.5.0
    
    
    func MessageNotContains(part [string](/builtin#string)) func(ctx [context](/context).[Context](/context#Context), r [slog](/log/slog).[Record](/log/slog#Record)) [bool](/builtin#bool)

MessageNotContains returns a function that checks if the record message does not contain the given part. Example usage: 
    
    
    r := slogmulti.Router().
        Add(consoleHandler, slogmulti.MessageNotContains("database error")).
        Add(fileHandler, slogmulti.MessageNotContains("database error")).
        Handler()
    

Args: 
    
    
    part: The part to check
    

Returns: 
    
    
    A function that checks if the record message does not contain the given part
    

####  func [NewHandleInlineHandler](https://github.com/samber/slog-multi/blob/v1.7.0/handler_inline_handle.go#L10) ¶ added in v1.3.0
    
    
    func NewHandleInlineHandler(handleFunc func(ctx [context](/context).[Context](/context#Context), groups [][string](/builtin#string), attrs [][slog](/log/slog).[Attr](/log/slog#Attr), record [slog](/log/slog).[Record](/log/slog#Record)) [error](/builtin#error)) [slog](/log/slog).[Handler](/log/slog#Handler)

NewHandleInlineHandler is a shortcut to a middleware that implements only the `Handle` method. 

####  func [NewInlineHandler](https://github.com/samber/slog-multi/blob/v1.7.0/handler_inline.go#L10) ¶ added in v1.3.0
    
    
    func NewInlineHandler(
    	enabledFunc func(ctx [context](/context).[Context](/context#Context), groups [][string](/builtin#string), attrs [][slog](/log/slog).[Attr](/log/slog#Attr), level [slog](/log/slog).[Level](/log/slog#Level)) [bool](/builtin#bool),
    	handleFunc func(ctx [context](/context).[Context](/context#Context), groups [][string](/builtin#string), attrs [][slog](/log/slog).[Attr](/log/slog#Attr), record [slog](/log/slog).[Record](/log/slog#Record)) [error](/builtin#error),
    ) [slog](/log/slog).[Handler](/log/slog#Handler)

NewInlineHandler is a shortcut to a handler that implements all methods. 

####  func [Pool](https://github.com/samber/slog-multi/blob/v1.7.0/pool.go#L48) ¶ added in v0.4.0
    
    
    func Pool() func(...[slog](/log/slog).[Handler](/log/slog#Handler)) [slog](/log/slog).[Handler](/log/slog#Handler)

Pool creates a load balancing handler factory function. This function returns a closure that can be used to create pool handlers with different sets of handlers for load balancing. 

The pool uses a round-robin strategy with randomization to distribute log records evenly across all available handlers. This is useful for: - Increasing logging throughput by parallelizing handler operations - Providing redundancy by having multiple handlers process the same records - Load balancing across multiple logging destinations 

Example usage: 
    
    
    handler := slogmulti.Pool()(
        handler1, // Will receive ~33% of records
        handler2, // Will receive ~33% of records
        handler3, // Will receive ~33% of records
    )
    logger := slog.New(handler)
    

Returns: 
    
    
    A function that creates PoolHandler instances with the provided handlers
    

####  func [RecoverHandlerError](https://github.com/samber/slog-multi/blob/v1.7.0/recover.go#L52) ¶ added in v1.4.0
    
    
    func RecoverHandlerError(recovery RecoveryFunc) func([slog](/log/slog).[Handler](/log/slog#Handler)) [slog](/log/slog).[Handler](/log/slog#Handler)

RecoverHandlerError creates a middleware that adds error recovery to a slog.Handler. This function returns a closure that can be used to wrap handlers with recovery logic. 

The recovery handler provides fault tolerance by: 1. Catching panics from the underlying handler 2. Catching errors returned by the underlying handler 3. Calling the recovery function with the error details 4. Propagating the original error to maintain logging semantics 

Example usage: 
    
    
    recovery := slogmulti.RecoverHandlerError(func(ctx context.Context, record slog.Record, err error) {
        fmt.Printf("Logging error: %v\n", err)
    })
    safeHandler := recovery(riskyHandler)
    logger := slog.New(safeHandler)
    

Args: 
    
    
    recovery: The function to call when an error or panic occurs
    

Returns: 
    
    
    A function that wraps handlers with recovery logic
    

####  func [Router](https://github.com/samber/slog-multi/blob/v1.7.0/router.go#L31) ¶ added in v0.5.0
    
    
    func Router() *router

Router creates a new router instance for building conditional log routing. This function is the entry point for creating a routing configuration. 

Example usage: 
    
    
    r := slogmulti.Router().
        Add(consoleHandler, slogmulti.LevelIs(slog.LevelInfo)).
        Add(fileHandler, slogmulti.LevelIs(slog.LevelError)).
        Handler()
    

Returns: 
    
    
    A new router instance ready for configuration
    

### Types ¶

####  type [EnabledInlineMiddleware](https://github.com/samber/slog-multi/blob/v1.7.0/middleware_inline_enabled.go#L28) ¶
    
    
    type EnabledInlineMiddleware struct {
    	// contains filtered or unexported fields
    }

####  func (*EnabledInlineMiddleware) [Enabled](https://github.com/samber/slog-multi/blob/v1.7.0/middleware_inline_enabled.go#L35) ¶
    
    
    func (h *EnabledInlineMiddleware) Enabled(ctx [context](/context).[Context](/context#Context), level [slog](/log/slog).[Level](/log/slog#Level)) [bool](/builtin#bool)

Implements slog.Handler 

####  func (*EnabledInlineMiddleware) [Handle](https://github.com/samber/slog-multi/blob/v1.7.0/middleware_inline_enabled.go#L40) ¶
    
    
    func (h *EnabledInlineMiddleware) Handle(ctx [context](/context).[Context](/context#Context), record [slog](/log/slog).[Record](/log/slog#Record)) [error](/builtin#error)

Implements slog.Handler 

####  func (*EnabledInlineMiddleware) [WithAttrs](https://github.com/samber/slog-multi/blob/v1.7.0/middleware_inline_enabled.go#L45) ¶
    
    
    func (h *EnabledInlineMiddleware) WithAttrs(attrs [][slog](/log/slog).[Attr](/log/slog#Attr)) [slog](/log/slog).[Handler](/log/slog#Handler)

Implements slog.Handler 

####  func (*EnabledInlineMiddleware) [WithGroup](https://github.com/samber/slog-multi/blob/v1.7.0/middleware_inline_enabled.go#L50) ¶
    
    
    func (h *EnabledInlineMiddleware) WithGroup(name [string](/builtin#string)) [slog](/log/slog).[Handler](/log/slog#Handler)

Implements slog.Handler 

####  type [FailoverHandler](https://github.com/samber/slog-multi/blob/v1.7.0/failover.go#L19) ¶ added in v0.4.0
    
    
    type FailoverHandler struct {
    	// contains filtered or unexported fields
    }

FailoverHandler implements a high-availability logging pattern. It attempts to forward log records to handlers in order until one succeeds. This is useful for scenarios where you want primary and backup logging destinations. 

@TODO: implement round robin strategy for load balancing across multiple handlers 

####  func (*FailoverHandler) [Enabled](https://github.com/samber/slog-multi/blob/v1.7.0/failover.go#L64) ¶ added in v0.4.0
    
    
    func (h *FailoverHandler) Enabled(ctx [context](/context).[Context](/context#Context), l [slog](/log/slog).[Level](/log/slog#Level)) [bool](/builtin#bool)

Enabled checks if any of the underlying handlers are enabled for the given log level. This method implements the slog.Handler interface requirement. 

The handler is considered enabled if at least one of its child handlers is enabled for the specified level. This ensures that if any handler can process the log, the failover handler will attempt to distribute it. 

Args: 
    
    
    ctx: The context for the logging operation
    l: The log level to check
    

Returns: 
    
    
    true if at least one handler is enabled for the level, false otherwise
    

####  func (*FailoverHandler) [Handle](https://github.com/samber/slog-multi/blob/v1.7.0/failover.go#L88) ¶ added in v0.4.0
    
    
    func (h *FailoverHandler) Handle(ctx [context](/context).[Context](/context#Context), r [slog](/log/slog).[Record](/log/slog#Record)) [error](/builtin#error)

Handle attempts to process a log record using handlers in priority order. This method implements the slog.Handler interface requirement. 

This implements a "fail-fast" strategy where the first successful handler prevents further attempts, making it efficient for high-availability scenarios. 

Args: 
    
    
    ctx: The context for the logging operation
    r: The log record to process
    

Returns: 
    
    
    nil if any handler successfully processed the record, or the last error encountered
    

####  func (*FailoverHandler) [WithAttrs](https://github.com/samber/slog-multi/blob/v1.7.0/failover.go#L119) ¶ added in v0.4.0
    
    
    func (h *FailoverHandler) WithAttrs(attrs [][slog](/log/slog).[Attr](/log/slog#Attr)) [slog](/log/slog).[Handler](/log/slog#Handler)

WithAttrs creates a new FailoverHandler with additional attributes added to all child handlers. This method implements the slog.Handler interface requirement. 

The method creates new handler instances for each child handler with the additional attributes, ensuring that the attributes are properly propagated to all handlers in the failover chain. 

Args: 
    
    
    attrs: The attributes to add to all handlers
    

Returns: 
    
    
    A new FailoverHandler with the attributes added to all child handlers
    

####  func (*FailoverHandler) [WithGroup](https://github.com/samber/slog-multi/blob/v1.7.0/failover.go#L141) ¶ added in v0.4.0
    
    
    func (h *FailoverHandler) WithGroup(name [string](/builtin#string)) [slog](/log/slog).[Handler](/log/slog#Handler)

WithGroup creates a new FailoverHandler with a group name applied to all child handlers. This method implements the slog.Handler interface requirement. 

The method follows the same pattern as the standard slog implementation: - If the group name is empty, returns the original handler unchanged - Otherwise, creates new handler instances for each child handler with the group name 

Args: 
    
    
    name: The group name to apply to all handlers
    

Returns: 
    
    
    A new FailoverHandler with the group name applied to all child handlers,
    or the original handler if the group name is empty
    

####  type [FanoutHandler](https://github.com/samber/slog-multi/blob/v1.7.0/multi.go#L18) ¶ added in v0.3.0
    
    
    type FanoutHandler struct {
    	// contains filtered or unexported fields
    }

FanoutHandler distributes log records to multiple slog.Handler instances in parallel. It implements the slog.Handler interface and forwards all logging operations to all registered handlers that are enabled for the given log level. 

####  func (*FanoutHandler) [Enabled](https://github.com/samber/slog-multi/blob/v1.7.0/multi.go#L76) ¶ added in v0.3.0
    
    
    func (h *FanoutHandler) Enabled(ctx [context](/context).[Context](/context#Context), l [slog](/log/slog).[Level](/log/slog#Level)) [bool](/builtin#bool)

Enabled checks if any of the underlying handlers are enabled for the given log level. This method implements the slog.Handler interface requirement. 

The handler is considered enabled if at least one of its child handlers is enabled for the specified level. This ensures that if any handler can process the log, the fanout handler will attempt to distribute it. 

Args: 
    
    
    ctx: The context for the logging operation
    l: The log level to check
    

Returns: 
    
    
    true if at least one handler is enabled for the level, false otherwise
    

####  func (*FanoutHandler) [Handle](https://github.com/samber/slog-multi/blob/v1.7.0/multi.go#L107) ¶ added in v0.3.0
    
    
    func (h *FanoutHandler) Handle(ctx [context](/context).[Context](/context#Context), r [slog](/log/slog).[Record](/log/slog#Record)) [error](/builtin#error)

Handle distributes a log record to all enabled handlers. This method implements the slog.Handler interface requirement. 

The method: 1. Iterates through all registered handlers 2. Checks if each handler is enabled for the record's level 3. For enabled handlers, calls their Handle method with a cloned record 4. Collects any errors that occur during handling 5. Returns a combined error if any handlers failed 

Note: Each handler receives a cloned record to prevent interference between handlers. This ensures that one handler cannot modify the record for other handlers. 

Args: 
    
    
    ctx: The context for the logging operation
    r: The log record to distribute
    

Returns: 
    
    
    An error if any handler failed to process the record, nil otherwise
    

####  func (*FanoutHandler) [WithAttrs](https://github.com/samber/slog-multi/blob/v1.7.0/multi.go#L138) ¶ added in v0.3.0
    
    
    func (h *FanoutHandler) WithAttrs(attrs [][slog](/log/slog).[Attr](/log/slog#Attr)) [slog](/log/slog).[Handler](/log/slog#Handler)

WithAttrs creates a new FanoutHandler with additional attributes added to all child handlers. This method implements the slog.Handler interface requirement. 

The method creates new handler instances for each child handler with the additional attributes, ensuring that the attributes are properly propagated to all handlers in the fanout chain. 

Args: 
    
    
    attrs: The attributes to add to all handlers
    

Returns: 
    
    
    A new FanoutHandler with the attributes added to all child handlers
    

####  func (*FanoutHandler) [WithGroup](https://github.com/samber/slog-multi/blob/v1.7.0/multi.go#L160) ¶ added in v0.3.0
    
    
    func (h *FanoutHandler) WithGroup(name [string](/builtin#string)) [slog](/log/slog).[Handler](/log/slog#Handler)

WithGroup creates a new FanoutHandler with a group name applied to all child handlers. This method implements the slog.Handler interface requirement. 

The method follows the same pattern as the standard slog implementation: - If the group name is empty, returns the original handler unchanged - Otherwise, creates new handler instances for each child handler with the group name 

Args: 
    
    
    name: The group name to apply to all handlers
    

Returns: 
    
    
    A new FanoutHandler with the group name applied to all child handlers,
    or the original handler if the group name is empty
    

####  type [FirstMatchHandler](https://github.com/samber/slog-multi/blob/v1.7.0/firstmatch.go#L14) ¶ added in v1.7.0
    
    
    type FirstMatchHandler struct {
    	// contains filtered or unexported fields
    }

####  func [FirstMatch](https://github.com/samber/slog-multi/blob/v1.7.0/firstmatch.go#L18) ¶ added in v1.7.0
    
    
    func FirstMatch(handlers ...*RoutableHandler) *FirstMatchHandler

####  func (*FirstMatchHandler) [Enabled](https://github.com/samber/slog-multi/blob/v1.7.0/firstmatch.go#L33) ¶ added in v1.7.0
    
    
    func (h *FirstMatchHandler) Enabled(ctx [context](/context).[Context](/context#Context), l [slog](/log/slog).[Level](/log/slog#Level)) [bool](/builtin#bool)

Enabled checks if any of the underlying handlers are enabled for the given log level. This method implements the slog.Handler interface requirement. See FanoutHandler.WithAttrs for details. 

####  func (*FirstMatchHandler) [Handle](https://github.com/samber/slog-multi/blob/v1.7.0/firstmatch.go#L52) ¶ added in v1.7.0
    
    
    func (h *FirstMatchHandler) Handle(ctx [context](/context).[Context](/context#Context), r [slog](/log/slog).[Record](/log/slog#Record)) [error](/builtin#error)

Handle distributes a log record to the first matching handler. This method implements the slog.Handler interface requirement. 

The method: 1. Iterates through each child handler. 2. Checks if the handler's predicates match the record. 3. If a match is found, it checks if the handler is enabled for the record's level. 4. If enabled, it forwards the record to that handler and returns. 5. If no handlers match, it returns nil. 

####  func (*FirstMatchHandler) [WithAttrs](https://github.com/samber/slog-multi/blob/v1.7.0/firstmatch.go#L72) ¶ added in v1.7.0
    
    
    func (h *FirstMatchHandler) WithAttrs(attrs [][slog](/log/slog).[Attr](/log/slog#Attr)) [slog](/log/slog).[Handler](/log/slog#Handler)

WithAttrs creates a new FirstMatchHandler with additional attributes added to all child handlers. This method implements the slog.Handler interface requirement. See FanoutHandler.WithAttrs for details. 

####  func (*FirstMatchHandler) [WithGroup](https://github.com/samber/slog-multi/blob/v1.7.0/firstmatch.go#L82) ¶ added in v1.7.0
    
    
    func (h *FirstMatchHandler) WithGroup(name [string](/builtin#string)) [slog](/log/slog).[Handler](/log/slog#Handler)

WithGroup creates a new FirstMatchHandler with a group name applied to all child handlers. This method implements the slog.Handler interface requirement. See FanoutHandler.WithGroup for details. 

####  type [HandleInlineHandler](https://github.com/samber/slog-multi/blob/v1.7.0/handler_inline_handle.go#L20) ¶ added in v1.3.0
    
    
    type HandleInlineHandler struct {
    	// contains filtered or unexported fields
    }

####  func (*HandleInlineHandler) [Enabled](https://github.com/samber/slog-multi/blob/v1.7.0/handler_inline_handle.go#L27) ¶ added in v1.3.0
    
    
    func (h *HandleInlineHandler) Enabled(ctx [context](/context).[Context](/context#Context), level [slog](/log/slog).[Level](/log/slog#Level)) [bool](/builtin#bool)

Implements slog.Handler 

####  func (*HandleInlineHandler) [Handle](https://github.com/samber/slog-multi/blob/v1.7.0/handler_inline_handle.go#L32) ¶ added in v1.3.0
    
    
    func (h *HandleInlineHandler) Handle(ctx [context](/context).[Context](/context#Context), record [slog](/log/slog).[Record](/log/slog#Record)) [error](/builtin#error)

Implements slog.Handler 

####  func (*HandleInlineHandler) [WithAttrs](https://github.com/samber/slog-multi/blob/v1.7.0/handler_inline_handle.go#L37) ¶ added in v1.3.0
    
    
    func (h *HandleInlineHandler) WithAttrs(attrs [][slog](/log/slog).[Attr](/log/slog#Attr)) [slog](/log/slog).[Handler](/log/slog#Handler)

Implements slog.Handler 

####  func (*HandleInlineHandler) [WithGroup](https://github.com/samber/slog-multi/blob/v1.7.0/handler_inline_handle.go#L50) ¶ added in v1.3.0
    
    
    func (h *HandleInlineHandler) WithGroup(name [string](/builtin#string)) [slog](/log/slog).[Handler](/log/slog#Handler)

Implements slog.Handler 

####  type [HandleInlineMiddleware](https://github.com/samber/slog-multi/blob/v1.7.0/middleware_inline_handle.go#L28) ¶
    
    
    type HandleInlineMiddleware struct {
    	// contains filtered or unexported fields
    }

####  func (*HandleInlineMiddleware) [Enabled](https://github.com/samber/slog-multi/blob/v1.7.0/middleware_inline_handle.go#L34) ¶
    
    
    func (h *HandleInlineMiddleware) Enabled(ctx [context](/context).[Context](/context#Context), level [slog](/log/slog).[Level](/log/slog#Level)) [bool](/builtin#bool)

Implements slog.Handler 

####  func (*HandleInlineMiddleware) [Handle](https://github.com/samber/slog-multi/blob/v1.7.0/middleware_inline_handle.go#L39) ¶
    
    
    func (h *HandleInlineMiddleware) Handle(ctx [context](/context).[Context](/context#Context), record [slog](/log/slog).[Record](/log/slog#Record)) [error](/builtin#error)

Implements slog.Handler 

####  func (*HandleInlineMiddleware) [WithAttrs](https://github.com/samber/slog-multi/blob/v1.7.0/middleware_inline_handle.go#L44) ¶
    
    
    func (h *HandleInlineMiddleware) WithAttrs(attrs [][slog](/log/slog).[Attr](/log/slog#Attr)) [slog](/log/slog).[Handler](/log/slog#Handler)

Implements slog.Handler 

####  func (*HandleInlineMiddleware) [WithGroup](https://github.com/samber/slog-multi/blob/v1.7.0/middleware_inline_handle.go#L49) ¶
    
    
    func (h *HandleInlineMiddleware) WithGroup(name [string](/builtin#string)) [slog](/log/slog).[Handler](/log/slog#Handler)

Implements slog.Handler 

####  type [HandlerErrorRecovery](https://github.com/samber/slog-multi/blob/v1.7.0/recover.go#L21) ¶ added in v1.4.0
    
    
    type HandlerErrorRecovery struct {
    	// contains filtered or unexported fields
    }

HandlerErrorRecovery wraps a slog.Handler to provide panic and error recovery. It catches both panics and errors from the underlying handler and calls a recovery function to handle them gracefully. 

####  func (*HandlerErrorRecovery) [Enabled](https://github.com/samber/slog-multi/blob/v1.7.0/recover.go#L72) ¶ added in v1.4.0
    
    
    func (h *HandlerErrorRecovery) Enabled(ctx [context](/context).[Context](/context#Context), l [slog](/log/slog).[Level](/log/slog#Level)) [bool](/builtin#bool)

Enabled checks if the underlying handler is enabled for the given log level. This method implements the slog.Handler interface requirement. 

Args: 
    
    
    ctx: The context for the logging operation
    l: The log level to check
    

Returns: 
    
    
    true if the underlying handler is enabled for the level, false otherwise
    

####  func (*HandlerErrorRecovery) [Handle](https://github.com/samber/slog-multi/blob/v1.7.0/recover.go#L90) ¶ added in v1.4.0
    
    
    func (h *HandlerErrorRecovery) Handle(ctx [context](/context).[Context](/context#Context), record [slog](/log/slog).[Record](/log/slog#Record)) [error](/builtin#error)

Handle processes a log record with error recovery. This method implements the slog.Handler interface requirement. 

This ensures that logging errors don't crash the application while still allowing the error to be handled appropriately by the calling code. 

Args: 
    
    
    ctx: The context for the logging operation
    record: The log record to process
    

Returns: 
    
    
    The error from the underlying handler (never nil if an error occurred)
    

####  func (*HandlerErrorRecovery) [WithAttrs](https://github.com/samber/slog-multi/blob/v1.7.0/recover.go#L120) ¶ added in v1.4.0
    
    
    func (h *HandlerErrorRecovery) WithAttrs(attrs [][slog](/log/slog).[Attr](/log/slog#Attr)) [slog](/log/slog).[Handler](/log/slog#Handler)

WithAttrs creates a new HandlerErrorRecovery with additional attributes. This method implements the slog.Handler interface requirement. 

Args: 
    
    
    attrs: The attributes to add to the underlying handler
    

Returns: 
    
    
    A new HandlerErrorRecovery with the additional attributes
    

####  func (*HandlerErrorRecovery) [WithGroup](https://github.com/samber/slog-multi/blob/v1.7.0/recover.go#L141) ¶ added in v1.4.0
    
    
    func (h *HandlerErrorRecovery) WithGroup(name [string](/builtin#string)) [slog](/log/slog).[Handler](/log/slog#Handler)

WithGroup creates a new HandlerErrorRecovery with a group name. This method implements the slog.Handler interface requirement. 

The method follows the same pattern as the standard slog implementation: - If the group name is empty, returns the original handler unchanged - Otherwise, creates a new handler with the group name applied to the underlying handler 

Args: 
    
    
    name: The group name to apply to the underlying handler
    

Returns: 
    
    
    A new HandlerErrorRecovery with the group name, or the original handler if the name is empty
    

####  type [InlineHandler](https://github.com/samber/slog-multi/blob/v1.7.0/handler_inline.go#L31) ¶ added in v1.3.0
    
    
    type InlineHandler struct {
    	// contains filtered or unexported fields
    }

####  func (*InlineHandler) [Enabled](https://github.com/samber/slog-multi/blob/v1.7.0/handler_inline.go#L39) ¶ added in v1.3.0
    
    
    func (h *InlineHandler) Enabled(ctx [context](/context).[Context](/context#Context), level [slog](/log/slog).[Level](/log/slog#Level)) [bool](/builtin#bool)

Implements slog.Handler 

####  func (*InlineHandler) [Handle](https://github.com/samber/slog-multi/blob/v1.7.0/handler_inline.go#L44) ¶ added in v1.3.0
    
    
    func (h *InlineHandler) Handle(ctx [context](/context).[Context](/context#Context), record [slog](/log/slog).[Record](/log/slog#Record)) [error](/builtin#error)

Implements slog.Handler 

####  func (*InlineHandler) [WithAttrs](https://github.com/samber/slog-multi/blob/v1.7.0/handler_inline.go#L49) ¶ added in v1.3.0
    
    
    func (h *InlineHandler) WithAttrs(attrs [][slog](/log/slog).[Attr](/log/slog#Attr)) [slog](/log/slog).[Handler](/log/slog#Handler)

Implements slog.Handler 

####  func (*InlineHandler) [WithGroup](https://github.com/samber/slog-multi/blob/v1.7.0/handler_inline.go#L63) ¶ added in v1.3.0
    
    
    func (h *InlineHandler) WithGroup(name [string](/builtin#string)) [slog](/log/slog).[Handler](/log/slog#Handler)

Implements slog.Handler 

####  type [InlineMiddleware](https://github.com/samber/slog-multi/blob/v1.7.0/middleware_inline.go#L45) ¶
    
    
    type InlineMiddleware struct {
    	// contains filtered or unexported fields
    }

####  func (*InlineMiddleware) [Enabled](https://github.com/samber/slog-multi/blob/v1.7.0/middleware_inline.go#L54) ¶
    
    
    func (h *InlineMiddleware) Enabled(ctx [context](/context).[Context](/context#Context), level [slog](/log/slog).[Level](/log/slog#Level)) [bool](/builtin#bool)

Implements slog.Handler 

####  func (*InlineMiddleware) [Handle](https://github.com/samber/slog-multi/blob/v1.7.0/middleware_inline.go#L59) ¶
    
    
    func (h *InlineMiddleware) Handle(ctx [context](/context).[Context](/context#Context), record [slog](/log/slog).[Record](/log/slog#Record)) [error](/builtin#error)

Implements slog.Handler 

####  func (*InlineMiddleware) [WithAttrs](https://github.com/samber/slog-multi/blob/v1.7.0/middleware_inline.go#L64) ¶
    
    
    func (h *InlineMiddleware) WithAttrs(attrs [][slog](/log/slog).[Attr](/log/slog#Attr)) [slog](/log/slog).[Handler](/log/slog#Handler)

Implements slog.Handler 

####  func (*InlineMiddleware) [WithGroup](https://github.com/samber/slog-multi/blob/v1.7.0/middleware_inline.go#L74) ¶
    
    
    func (h *InlineMiddleware) WithGroup(name [string](/builtin#string)) [slog](/log/slog).[Handler](/log/slog#Handler)

Implements slog.Handler 

####  type [Middleware](https://github.com/samber/slog-multi/blob/v1.7.0/middleware.go#L28) ¶
    
    
    type Middleware func([slog](/log/slog).[Handler](/log/slog#Handler)) [slog](/log/slog).[Handler](/log/slog#Handler)

Middleware is a function type that transforms one slog.Handler into another. It follows the standard middleware pattern where a function takes a handler and returns a new handler that wraps the original with additional functionality. 

Middleware functions can be used to: - Transform log records (e.g., add timestamps, modify levels) - Filter records based on conditions - Add context or attributes to records - Implement cross-cutting concerns like error recovery or sampling 

Example usage: 
    
    
      gdprMiddleware := NewGDPRMiddleware()
      sink := slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{})
    
    	 logger := slog.New(
    		slogmulti.
    			Pipe(gdprMiddleware).
    			// ...
    			Handler(sink),
    	  )
    

####  func [NewEnabledInlineMiddleware](https://github.com/samber/slog-multi/blob/v1.7.0/middleware_inline_enabled.go#L10) ¶
    
    
    func NewEnabledInlineMiddleware(enabledFunc func(ctx [context](/context).[Context](/context#Context), level [slog](/log/slog).[Level](/log/slog#Level), next func([context](/context).[Context](/context#Context), [slog](/log/slog).[Level](/log/slog#Level)) [bool](/builtin#bool)) [bool](/builtin#bool)) Middleware

NewEnabledInlineMiddleware is shortcut to a middleware that implements only the `Enable` method. 

####  func [NewHandleInlineMiddleware](https://github.com/samber/slog-multi/blob/v1.7.0/middleware_inline_handle.go#L10) ¶
    
    
    func NewHandleInlineMiddleware(handleFunc func(ctx [context](/context).[Context](/context#Context), record [slog](/log/slog).[Record](/log/slog#Record), next func([context](/context).[Context](/context#Context), [slog](/log/slog).[Record](/log/slog#Record)) [error](/builtin#error)) [error](/builtin#error)) Middleware

NewHandleInlineMiddleware is a shortcut to a middleware that implements only the `Handle` method. 

####  func [NewInlineMiddleware](https://github.com/samber/slog-multi/blob/v1.7.0/middleware_inline.go#L10) ¶
    
    
    func NewInlineMiddleware(
    	enabledFunc func(ctx [context](/context).[Context](/context#Context), level [slog](/log/slog).[Level](/log/slog#Level), next func([context](/context).[Context](/context#Context), [slog](/log/slog).[Level](/log/slog#Level)) [bool](/builtin#bool)) [bool](/builtin#bool),
    	handleFunc func(ctx [context](/context).[Context](/context#Context), record [slog](/log/slog).[Record](/log/slog#Record), next func([context](/context).[Context](/context#Context), [slog](/log/slog).[Record](/log/slog#Record)) [error](/builtin#error)) [error](/builtin#error),
    	withAttrsFunc func(attrs [][slog](/log/slog).[Attr](/log/slog#Attr), next func([][slog](/log/slog).[Attr](/log/slog#Attr)) [slog](/log/slog).[Handler](/log/slog#Handler)) [slog](/log/slog).[Handler](/log/slog#Handler),
    	withGroupFunc func(name [string](/builtin#string), next func([string](/builtin#string)) [slog](/log/slog).[Handler](/log/slog#Handler)) [slog](/log/slog).[Handler](/log/slog#Handler),
    ) Middleware

NewInlineMiddleware is a shortcut to a middleware that implements all methods. 

####  func [NewWithAttrsInlineMiddleware](https://github.com/samber/slog-multi/blob/v1.7.0/middleware_inline_with_attrs.go#L10) ¶
    
    
    func NewWithAttrsInlineMiddleware(withAttrsFunc func(attrs [][slog](/log/slog).[Attr](/log/slog#Attr), next func([][slog](/log/slog).[Attr](/log/slog#Attr)) [slog](/log/slog).[Handler](/log/slog#Handler)) [slog](/log/slog).[Handler](/log/slog#Handler)) Middleware

NewWithAttrsInlineMiddleware is a shortcut to a middleware that implements only the `WithAttrs` method. 

####  func [NewWithGroupInlineMiddleware](https://github.com/samber/slog-multi/blob/v1.7.0/middleware_inline_with_group.go#L10) ¶
    
    
    func NewWithGroupInlineMiddleware(withGroupFunc func(name [string](/builtin#string), next func([string](/builtin#string)) [slog](/log/slog).[Handler](/log/slog#Handler)) [slog](/log/slog).[Handler](/log/slog#Handler)) Middleware

NewWithGroupInlineMiddleware is a shortcut to a middleware that implements only the `WithAttrs` method. 

####  type [PipeBuilder](https://github.com/samber/slog-multi/blob/v1.7.0/pipe.go#L10) ¶
    
    
    type PipeBuilder struct {
    	// contains filtered or unexported fields
    }

PipeBuilder provides a fluent API for building middleware chains. It allows you to compose multiple middleware functions that will be applied to log records in the order they are added (last-in, first-out). 

####  func [Pipe](https://github.com/samber/slog-multi/blob/v1.7.0/pipe.go#L39) ¶
    
    
    func Pipe(middlewares ...Middleware) *PipeBuilder

Pipe creates a new PipeBuilder with the provided middleware functions. This function is the entry point for building middleware chains. 

Middleware functions are applied in reverse order (last-in, first-out), which means the last middleware added will be the first one applied to incoming records. This allows for intuitive composition where you can think of the chain as "transform A, then transform B, then send to handler". 

Example usage: 
    
    
    handler := slogmulti.Pipe(
        RewriteLevel(slog.LevelWarn, slog.LevelInfo),
        RewriteMessage("prefix: %s"),
        RedactPII(),
    ).Handler(finalHandler)
    

Args: 
    
    
    middlewares: Variable number of middleware functions to chain together
    

Returns: 
    
    
    A new PipeBuilder instance ready for further configuration
    

####  func (*PipeBuilder) [Handler](https://github.com/samber/slog-multi/blob/v1.7.0/pipe.go#L71) ¶
    
    
    func (h *PipeBuilder) Handler(handler [slog](/log/slog).[Handler](/log/slog#Handler)) [slog](/log/slog).[Handler](/log/slog#Handler)

Handler creates a slog.Handler by applying all middleware to the provided handler. This method finalizes the middleware chain and returns a handler that can be used with slog.New(). 

This LIFO approach ensures that the middleware chain is applied in the intuitive order: the first middleware in the chain is applied first to incoming records. 

Args: 
    
    
    handler: The final slog.Handler that will receive the transformed records
    

Returns: 
    
    
    A slog.Handler that applies all middleware transformations before forwarding to the final handler
    

####  func (*PipeBuilder) [Pipe](https://github.com/samber/slog-multi/blob/v1.7.0/pipe.go#L53) ¶
    
    
    func (h *PipeBuilder) Pipe(middleware Middleware) *PipeBuilder

Pipe adds an additional middleware to the chain. This method provides a fluent API for building middleware chains incrementally. 

Args: 
    
    
    middleware: The middleware function to add to the chain
    

Returns: 
    
    
    The PipeBuilder instance for method chaining
    

####  type [PoolHandler](https://github.com/samber/slog-multi/blob/v1.7.0/pool.go#L19) ¶ added in v0.4.0
    
    
    type PoolHandler struct {
    	// contains filtered or unexported fields
    }

PoolHandler implements a load balancing strategy for logging handlers. It distributes log records across multiple handlers using a round-robin approach with randomization to ensure even distribution and prevent hot-spotting. 

####  func (*PoolHandler) [Enabled](https://github.com/samber/slog-multi/blob/v1.7.0/pool.go#L72) ¶ added in v0.4.0
    
    
    func (h *PoolHandler) Enabled(ctx [context](/context).[Context](/context#Context), l [slog](/log/slog).[Level](/log/slog#Level)) [bool](/builtin#bool)

Enabled checks if any of the underlying handlers are enabled for the given log level. This method implements the slog.Handler interface requirement. 

The handler is considered enabled if at least one of its child handlers is enabled for the specified level. This ensures that if any handler can process the log, the pool handler will attempt to distribute it. 

Args: 
    
    
    ctx: The context for the logging operation
    l: The log level to check
    

Returns: 
    
    
    true if at least one handler is enabled for the level, false otherwise
    

####  func (*PoolHandler) [Handle](https://github.com/samber/slog-multi/blob/v1.7.0/pool.go#L96) ¶ added in v0.4.0
    
    
    func (h *PoolHandler) Handle(ctx [context](/context).[Context](/context#Context), r [slog](/log/slog).[Record](/log/slog#Record)) [error](/builtin#error)

Handle distributes a log record to a handler selected using round-robin with randomization. This method implements the slog.Handler interface requirement. 

This approach ensures even distribution of load while providing fault tolerance through the failover behavior when a handler is unavailable. 

Args: 
    
    
    ctx: The context for the logging operation
    r: The log record to distribute
    

Returns: 
    
    
    nil if any handler successfully processed the record, or the last error encountered
    

####  func (*PoolHandler) [WithAttrs](https://github.com/samber/slog-multi/blob/v1.7.0/pool.go#L135) ¶ added in v0.4.0
    
    
    func (h *PoolHandler) WithAttrs(attrs [][slog](/log/slog).[Attr](/log/slog#Attr)) [slog](/log/slog).[Handler](/log/slog#Handler)

WithAttrs creates a new PoolHandler with additional attributes added to all child handlers. This method implements the slog.Handler interface requirement. 

The method creates new handler instances for each child handler with the additional attributes, ensuring that the attributes are properly propagated to all handlers in the pool. 

Args: 
    
    
    attrs: The attributes to add to all handlers
    

Returns: 
    
    
    A new PoolHandler with the attributes added to all child handlers
    

####  func (*PoolHandler) [WithGroup](https://github.com/samber/slog-multi/blob/v1.7.0/pool.go#L157) ¶ added in v0.4.0
    
    
    func (h *PoolHandler) WithGroup(name [string](/builtin#string)) [slog](/log/slog).[Handler](/log/slog#Handler)

WithGroup creates a new PoolHandler with a group name applied to all child handlers. This method implements the slog.Handler interface requirement. 

The method follows the same pattern as the standard slog implementation: - If the group name is empty, returns the original handler unchanged - Otherwise, creates new handler instances for each child handler with the group name 

Args: 
    
    
    name: The group name to apply to all handlers
    

Returns: 
    
    
    A new PoolHandler with the group name applied to all child handlers,
    or the original handler if the group name is empty
    

####  type [RecoveryFunc](https://github.com/samber/slog-multi/blob/v1.7.0/recover.go#L13) ¶ added in v1.4.0
    
    
    type RecoveryFunc func(ctx [context](/context).[Context](/context#Context), record [slog](/log/slog).[Record](/log/slog#Record), err [error](/builtin#error))

RecoveryFunc is a callback function that handles errors and panics from logging handlers. It receives the context, the log record that caused the error, and the error itself. This function can be used to log the error, send alerts, or perform any other error handling logic without affecting the main application flow. 

####  type [RoutableHandler](https://github.com/samber/slog-multi/blob/v1.7.0/router.go#L101) ¶ added in v0.5.0
    
    
    type RoutableHandler struct {
    	// contains filtered or unexported fields
    }

RoutableHandler wraps a slog.Handler with conditional matching logic. It only forwards records to the underlying handler if all predicates return true. This enables sophisticated routing scenarios like level-based or attribute-based routing. 

@TODO: implement round robin strategy for load balancing across multiple handlers 

####  func (*RoutableHandler) [Enabled](https://github.com/samber/slog-multi/blob/v1.7.0/router.go#L125) ¶ added in v0.5.0
    
    
    func (h *RoutableHandler) Enabled(ctx [context](/context).[Context](/context#Context), l [slog](/log/slog).[Level](/log/slog#Level)) [bool](/builtin#bool)

Enabled checks if the underlying handler is enabled for the given log level. This method implements the slog.Handler interface requirement. 

Args: 
    
    
    ctx: The context for the logging operation
    l: The log level to check
    

Returns: 
    
    
    true if the underlying handler is enabled for the level, false otherwise
    

####  func (*RoutableHandler) [Handle](https://github.com/samber/slog-multi/blob/v1.7.0/router.go#L140) ¶ added in v0.5.0
    
    
    func (h *RoutableHandler) Handle(ctx [context](/context).[Context](/context#Context), r [slog](/log/slog).[Record](/log/slog#Record)) [error](/builtin#error)

Handle processes a log record if all predicates return true. This method implements the slog.Handler interface requirement. 

Args: 
    
    
    ctx: The context for the logging operation
    r: The log record to process
    

Returns: 
    
    
    An error if the underlying handler failed to process the record, nil otherwise
    

####  func (*RoutableHandler) [WithAttrs](https://github.com/samber/slog-multi/blob/v1.7.0/router.go#L181) ¶ added in v0.5.0
    
    
    func (h *RoutableHandler) WithAttrs(attrs [][slog](/log/slog).[Attr](/log/slog#Attr)) [slog](/log/slog).[Handler](/log/slog#Handler)

WithAttrs creates a new RoutableHandler with additional attributes. This method implements the slog.Handler interface requirement. 

The method properly handles attribute accumulation within the current group context, ensuring that attributes are correctly applied to records when they are processed. 

Args: 
    
    
    attrs: The attributes to add to the handler
    

Returns: 
    
    
    A new RoutableHandler with the additional attributes
    

####  func (*RoutableHandler) [WithGroup](https://github.com/samber/slog-multi/blob/v1.7.0/router.go#L205) ¶ added in v0.5.0
    
    
    func (h *RoutableHandler) WithGroup(name [string](/builtin#string)) [slog](/log/slog).[Handler](/log/slog#Handler)

WithGroup creates a new RoutableHandler with a group name. This method implements the slog.Handler interface requirement. 

The method follows the same pattern as the standard slog implementation: - If the group name is empty, returns the original handler unchanged - Otherwise, creates a new handler with the group name added to the group hierarchy 

Args: 
    
    
    name: The group name to apply to the handler
    

Returns: 
    
    
    A new RoutableHandler with the group name, or the original handler if the name is empty
    

####  type [WithAttrsInlineMiddleware](https://github.com/samber/slog-multi/blob/v1.7.0/middleware_inline_with_attrs.go#L28) ¶
    
    
    type WithAttrsInlineMiddleware struct {
    	// contains filtered or unexported fields
    }

####  func (*WithAttrsInlineMiddleware) [Enabled](https://github.com/samber/slog-multi/blob/v1.7.0/middleware_inline_with_attrs.go#L34) ¶
    
    
    func (h *WithAttrsInlineMiddleware) Enabled(ctx [context](/context).[Context](/context#Context), level [slog](/log/slog).[Level](/log/slog#Level)) [bool](/builtin#bool)

Implements slog.Handler 

####  func (*WithAttrsInlineMiddleware) [Handle](https://github.com/samber/slog-multi/blob/v1.7.0/middleware_inline_with_attrs.go#L39) ¶
    
    
    func (h *WithAttrsInlineMiddleware) Handle(ctx [context](/context).[Context](/context#Context), record [slog](/log/slog).[Record](/log/slog#Record)) [error](/builtin#error)

Implements slog.Handler 

####  func (*WithAttrsInlineMiddleware) [WithAttrs](https://github.com/samber/slog-multi/blob/v1.7.0/middleware_inline_with_attrs.go#L44) ¶
    
    
    func (h *WithAttrsInlineMiddleware) WithAttrs(attrs [][slog](/log/slog).[Attr](/log/slog#Attr)) [slog](/log/slog).[Handler](/log/slog#Handler)

Implements slog.Handler 

####  func (*WithAttrsInlineMiddleware) [WithGroup](https://github.com/samber/slog-multi/blob/v1.7.0/middleware_inline_with_attrs.go#L49) ¶
    
    
    func (h *WithAttrsInlineMiddleware) WithGroup(name [string](/builtin#string)) [slog](/log/slog).[Handler](/log/slog#Handler)

Implements slog.Handler 

####  type [WithGroupInlineMiddleware](https://github.com/samber/slog-multi/blob/v1.7.0/middleware_inline_with_group.go#L28) ¶
    
    
    type WithGroupInlineMiddleware struct {
    	// contains filtered or unexported fields
    }

####  func (*WithGroupInlineMiddleware) [Enabled](https://github.com/samber/slog-multi/blob/v1.7.0/middleware_inline_with_group.go#L34) ¶
    
    
    func (h *WithGroupInlineMiddleware) Enabled(ctx [context](/context).[Context](/context#Context), level [slog](/log/slog).[Level](/log/slog#Level)) [bool](/builtin#bool)

Implements slog.Handler 

####  func (*WithGroupInlineMiddleware) [Handle](https://github.com/samber/slog-multi/blob/v1.7.0/middleware_inline_with_group.go#L39) ¶
    
    
    func (h *WithGroupInlineMiddleware) Handle(ctx [context](/context).[Context](/context#Context), record [slog](/log/slog).[Record](/log/slog#Record)) [error](/builtin#error)

Implements slog.Handler 

####  func (*WithGroupInlineMiddleware) [WithAttrs](https://github.com/samber/slog-multi/blob/v1.7.0/middleware_inline_with_group.go#L44) ¶
    
    
    func (h *WithGroupInlineMiddleware) WithAttrs(attrs [][slog](/log/slog).[Attr](/log/slog#Attr)) [slog](/log/slog).[Handler](/log/slog#Handler)

Implements slog.Handler 

####  func (*WithGroupInlineMiddleware) [WithGroup](https://github.com/samber/slog-multi/blob/v1.7.0/middleware_inline_with_group.go#L49) ¶
    
    
    func (h *WithGroupInlineMiddleware) WithGroup(name [string](/builtin#string)) [slog](/log/slog).[Handler](/log/slog#Handler)

Implements slog.Handler 
  *[v]: View this template
  *[t]: Discuss this template
  *[e]: Edit this template
