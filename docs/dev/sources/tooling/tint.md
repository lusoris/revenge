# lmittmann/tint

> Source: https://pkg.go.dev/github.com/lmittmann/tint
> Fetched: 2026-01-31T10:56:46.371785+00:00
> Content-Hash: f1d8722427984693
> Type: html

---

### Overview ¶

  * Customize Attributes
  * Automatically Enable Colors
  * Windows Support



Package tint implements a zero-dependency [slog.Handler](/log/slog#Handler) that writes tinted (colorized) logs. The output format is inspired by the [zerolog.ConsoleWriter](https://pkg.go.dev/github.com/rs/zerolog#ConsoleWriter) and [slog.TextHandler](/log/slog#TextHandler). 

The output format can be customized using Options, which is a drop-in replacement for [slog.HandlerOptions](/log/slog#HandlerOptions). 

#### Customize Attributes ¶

Options.ReplaceAttr can be used to alter or drop attributes. If set, it is called on each non-group attribute before it is logged. See [slog.HandlerOptions](/log/slog#HandlerOptions) for details. 

Create a new logger with a custom TRACE level: 
    
    
    const LevelTrace = slog.LevelDebug - 4
    
    w := os.Stderr
    logger := slog.New(tint.NewHandler(w, &tint.Options{
    	Level: LevelTrace,
    	ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
    		if a.Key == slog.LevelKey && len(groups) == 0 {
    			level, ok := a.Value.Any().(slog.Level)
    			if ok && level <= LevelTrace {
    				return tint.Attr(13, slog.String(a.Key, "TRC"))
    			}
    		}
    		return a
    	},
    }))
    

Create a new logger that doesn't write the time: 
    
    
    w := os.Stderr
    logger := slog.New(
    	tint.NewHandler(w, &tint.Options{
    		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
    			if a.Key == slog.TimeKey && len(groups) == 0 {
    				return slog.Attr{}
    			}
    			return a
    		},
    	}),
    )
    

Create a new logger that writes all errors in red: 
    
    
    w := os.Stderr
    logger := slog.New(
    	tint.NewHandler(w, &tint.Options{
    		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
    			if a.Value.Kind() == slog.KindAny {
    				if _, ok := a.Value.Any().(error); ok {
    					return tint.Attr(9, a)
    				}
    			}
    			return a
    		},
    	}),
    )
    

#### Automatically Enable Colors ¶

Colors are enabled by default. Use the Options.NoColor field to disable color output. To automatically enable colors based on terminal capabilities, use e.g., the [go-isatty](https://pkg.go.dev/github.com/mattn/go-isatty) package: 
    
    
    w := os.Stderr
    logger := slog.New(
    	tint.NewHandler(w, &tint.Options{
    		NoColor: !isatty.IsTerminal(w.Fd()),
    	}),
    )
    

#### Windows Support ¶

Color support on Windows can be added by using e.g., the [go-colorable](https://pkg.go.dev/github.com/mattn/go-colorable) package: 
    
    
    w := os.Stderr
    logger := slog.New(
    	tint.NewHandler(colorable.NewColorable(w), nil),
    )
    

Example ¶
    
    
    package main
    
    import (
    	"errors"
    	"log/slog"
    	"os"
    	"time"
    
    	"github.com/lmittmann/tint"
    )
    
    func main() {
    	w := os.Stderr
    	logger := slog.New(tint.NewHandler(w, &tint.Options{
    		Level:      slog.LevelDebug,
    		TimeFormat: time.Kitchen,
    	}))
    	logger.Info("Starting server", "addr", ":8080", "env", "production")
    	logger.Debug("Connected to DB", "db", "myapp", "host", "localhost:5432")
    	logger.Warn("Slow request", "method", "GET", "path", "/users", "duration", 497*time.Millisecond)
    	logger.Error("DB connection lost", tint.Err(errors.New("connection reset")), "db", "myapp")
    }
    

Share Format Run

Example (RedErrors) ¶

Create a new logger that writes all errors in red: 
    
    
    package main
    
    import (
    	"errors"
    	"log/slog"
    	"os"
    
    	"github.com/lmittmann/tint"
    )
    
    func main() {
    	w := os.Stderr
    	logger := slog.New(tint.NewHandler(w, &tint.Options{
    		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
    			if a.Value.Kind() == slog.KindAny {
    				if _, ok := a.Value.Any().(error); ok {
    					return tint.Attr(9, a)
    				}
    			}
    			return a
    		},
    	}))
    	logger.Error("DB connection lost", "error", errors.New("connection reset"), "db", "myapp")
    }
    

Share Format Run

Example (TraceLevel) ¶

Create a new logger with a custom TRACE level: 
    
    
    package main
    
    import (
    	"context"
    	"log/slog"
    	"os"
    	"time"
    
    	"github.com/lmittmann/tint"
    )
    
    func main() {
    	const LevelTrace = slog.LevelDebug - 4
    
    	w := os.Stderr
    	logger := slog.New(tint.NewHandler(w, &tint.Options{
    		Level: LevelTrace,
    		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
    			if a.Key == slog.LevelKey && len(groups) == 0 {
    				level, ok := a.Value.Any().(slog.Level)
    				if ok && level <= LevelTrace {
    					return tint.Attr(13, slog.String(a.Key, "TRC"))
    				}
    			}
    			return a
    		},
    	}))
    	logger.Log(context.Background(), LevelTrace, "DB query", "query", "SELECT * FROM users", "duration", 543*time.Microsecond)
    }
    

Share Format Run

### Index ¶

  * func Attr(color uint8, attr slog.Attr) slog.Attr
  * func Err(err error) slog.Attr
  * func NewHandler(w io.Writer, opts *Options) slog.Handler
  * type Options



### Examples ¶

  * Package
  * Package (RedErrors)
  * Package (TraceLevel)



### Constants ¶

This section is empty.

### Variables ¶

This section is empty.

### Functions ¶

####  func [Attr](https://github.com/lmittmann/tint/blob/v1.1.2/handler.go#L751) ¶ added in v1.1.0
    
    
    func Attr(color [uint8](/builtin#uint8), attr [slog](/log/slog).[Attr](/log/slog#Attr)) [slog](/log/slog).[Attr](/log/slog#Attr)

Attr returns a tinted (colorized) [slog.Attr](/log/slog#Attr) that will be written in the specified color by the tint.Handler. When used with any other [slog.Handler](/log/slog#Handler), it behaves as a plain [slog.Attr](/log/slog#Attr). 

Use the uint8 color value to specify the color of the attribute: 

  * 0-7: standard ANSI colors
  * 8-15: high intensity ANSI colors
  * 16-231: 216 colors (6×6×6 cube)
  * 232-255: grayscale from dark to light in 24 steps



See <https://en.wikipedia.org/wiki/ANSI_escape_code#8-bit>

####  func [Err](https://github.com/lmittmann/tint/blob/v1.1.2/handler.go#L735) ¶
    
    
    func Err(err [error](/builtin#error)) [slog](/log/slog).[Attr](/log/slog#Attr)

Err returns a tinted (colorized) [slog.Attr](/log/slog#Attr) that will be written in red color by the tint.Handler. When used with any other [slog.Handler](/log/slog#Handler), it behaves as 
    
    
    slog.Any("err", err)
    

####  func [NewHandler](https://github.com/lmittmann/tint/blob/v1.1.2/handler.go#L157) ¶
    
    
    func NewHandler(w [io](/io).[Writer](/io#Writer), opts *Options) [slog](/log/slog).[Handler](/log/slog#Handler)

NewHandler creates a [slog.Handler](/log/slog#Handler) that writes tinted logs to Writer w, using the default options. If opts is nil, the default options are used. 

### Types ¶

####  type [Options](https://github.com/lmittmann/tint/blob/v1.1.2/handler.go#L128) ¶
    
    
    type Options struct {
    	// Enable source code location (Default: false)
    	AddSource [bool](/builtin#bool)
    
    	// Minimum level to log (Default: slog.LevelInfo)
    	Level [slog](/log/slog).[Leveler](/log/slog#Leveler)
    
    	// ReplaceAttr is called to rewrite each non-group attribute before it is logged.
    	// See <https://pkg.go.dev/log/slog#HandlerOptions> for details.
    	ReplaceAttr func(groups [][string](/builtin#string), attr [slog](/log/slog).[Attr](/log/slog#Attr)) [slog](/log/slog).[Attr](/log/slog#Attr)
    
    	// Time format (Default: time.StampMilli)
    	TimeFormat [string](/builtin#string)
    
    	// Disable color (Default: false)
    	NoColor [bool](/builtin#bool)
    }

Options for a slog.Handler that writes tinted logs. A zero Options consists entirely of default values. 

Options can be used as a drop-in replacement for [slog.HandlerOptions](/log/slog#HandlerOptions). 
