# cenkalti/backoff

> Source: https://pkg.go.dev/github.com/cenkalti/backoff/v4
> Fetched: 2026-01-30T23:49:35.960865+00:00
> Content-Hash: f73f1a6cb319256b
> Type: html

---

Overview

¶

Package backoff implements backoff algorithms for retrying operations.

Use Retry function for retrying operations that may fail.
If Retry does not meet your needs,
copy/paste the function into your project and modify as you wish.

There is also Ticker type similar to time.Ticker.
You can use it if you need to work with channels.

See Examples section below for usage examples.

Index

¶

Constants

Variables

func Permanent(err error) error

func Retry(o Operation, b BackOff) error

func RetryNotify(operation Operation, b BackOff, notify Notify) error

func RetryNotifyWithData[T any](operation OperationWithData[T], b BackOff, notify Notify) (T, error)

func RetryNotifyWithTimer(operation Operation, b BackOff, notify Notify, t Timer) error

func RetryNotifyWithTimerAndData[T any](operation OperationWithData[T], b BackOff, notify Notify, t Timer) (T, error)

func RetryWithData[T any](o OperationWithData[T], b BackOff) (T, error)

type BackOff

func WithMaxRetries(b BackOff, max uint64) BackOff

type BackOffContext

func WithContext(b BackOff, ctx context.Context) BackOffContext

type Clock

type ConstantBackOff

func NewConstantBackOff(d time.Duration) *ConstantBackOff

func (b *ConstantBackOff) NextBackOff() time.Duration

func (b *ConstantBackOff) Reset()

type ExponentialBackOff

func NewExponentialBackOff(opts ...ExponentialBackOffOpts) *ExponentialBackOff

func (b *ExponentialBackOff) GetElapsedTime() time.Duration

func (b *ExponentialBackOff) NextBackOff() time.Duration

func (b *ExponentialBackOff) Reset()

type ExponentialBackOffOpts

func WithClockProvider(clock Clock) ExponentialBackOffOpts

func WithInitialInterval(duration time.Duration) ExponentialBackOffOpts

func WithMaxElapsedTime(duration time.Duration) ExponentialBackOffOpts

func WithMaxInterval(duration time.Duration) ExponentialBackOffOpts

func WithMultiplier(multiplier float64) ExponentialBackOffOpts

func WithRandomizationFactor(randomizationFactor float64) ExponentialBackOffOpts

func WithRetryStopDuration(duration time.Duration) ExponentialBackOffOpts

type Notify

type Operation

type OperationWithData

type PermanentError

func (e *PermanentError) Error() string

func (e *PermanentError) Is(target error) bool

func (e *PermanentError) Unwrap() error

type StopBackOff

func (b *StopBackOff) NextBackOff() time.Duration

func (b *StopBackOff) Reset()

type Ticker

func NewTicker(b BackOff) *Ticker

func NewTickerWithTimer(b BackOff, timer Timer) *Ticker

func (t *Ticker) Stop()

type Timer

type ZeroBackOff

func (b *ZeroBackOff) NextBackOff() time.Duration

func (b *ZeroBackOff) Reset()

Examples

¶

Retry

Ticker

Constants

¶

View Source

const (

DefaultInitialInterval     = 500 *

time

.

Millisecond

DefaultRandomizationFactor = 0.5

DefaultMultiplier          = 1.5

DefaultMaxInterval         = 60 *

time

.

Second

DefaultMaxElapsedTime      = 15 *

time

.

Minute

)

Default values for ExponentialBackOff.

View Source

const Stop

time

.

Duration

= -1

Stop indicates that no more retries should be made for use in NextBackOff().

Variables

¶

View Source

var SystemClock = systemClock{}

SystemClock implements Clock interface that uses time.Now().

Functions

¶

func

Permanent

¶

func Permanent(err

error

)

error

Permanent wraps the given err in a *PermanentError.

func

Retry

¶

func Retry(o

Operation

, b

BackOff

)

error

Retry the operation o until it does not return error or BackOff stops.
o is guaranteed to be run at least once.

If o returns a *PermanentError, the operation is not retried, and the
wrapped error is returned.

Retry sleeps the goroutine for the duration returned by BackOff after a
failed operation returns.

Example

¶

// An operation that may fail.
operation := func() error {
	return nil // or an error
}

err := Retry(operation, NewExponentialBackOff())
if err != nil {
	// Handle error.
	return
}

// Operation is successful.

func

RetryNotify

¶

func RetryNotify(operation

Operation

, b

BackOff

, notify

Notify

)

error

RetryNotify calls notify function with the error and wait duration
for each failed attempt before sleep.

func

RetryNotifyWithData

¶

added in

v4.2.0

func RetryNotifyWithData[T

any

](operation

OperationWithData

[T], b

BackOff

, notify

Notify

) (T,

error

)

RetryNotifyWithData is like RetryNotify but returns data in the response too.

func

RetryNotifyWithTimer

¶

func RetryNotifyWithTimer(operation

Operation

, b

BackOff

, notify

Notify

, t

Timer

)

error

RetryNotifyWithTimer calls notify function with the error and wait duration using the given Timer
for each failed attempt before sleep.
A default timer that uses system timer is used when nil is passed.

func

RetryNotifyWithTimerAndData

¶

added in

v4.2.0

func RetryNotifyWithTimerAndData[T

any

](operation

OperationWithData

[T], b

BackOff

, notify

Notify

, t

Timer

) (T,

error

)

RetryNotifyWithTimerAndData is like RetryNotifyWithTimer but returns data in the response too.

func

RetryWithData

¶

added in

v4.2.0

func RetryWithData[T

any

](o

OperationWithData

[T], b

BackOff

) (T,

error

)

RetryWithData is like Retry but returns data in the response too.

Types

¶

type

BackOff

¶

type BackOff interface {

// NextBackOff returns the duration to wait before retrying the operation,

// or backoff. Stop to indicate that no more retries should be made.

//

// Example usage:

//

// 	duration := backoff.NextBackOff();

// 	if (duration == backoff.Stop) {

// 		// Do not retry operation.

// 	} else {

// 		// Sleep for duration and retry operation.

// 	}

//

NextBackOff()

time

.

Duration

// Reset to initial state.

Reset()
}

BackOff is a backoff policy for retrying an operation.

func

WithMaxRetries

¶

func WithMaxRetries(b

BackOff

, max

uint64

)

BackOff

WithMaxRetries creates a wrapper around another BackOff, which will
return Stop if NextBackOff() has been called too many times since
the last time Reset() was called

Note: Implementation is not thread-safe.

type

BackOffContext

¶

type BackOffContext interface {

BackOff

Context()

context

.

Context

}

BackOffContext is a backoff policy that stops retrying after the context
is canceled.

func

WithContext

¶

func WithContext(b

BackOff

, ctx

context

.

Context

)

BackOffContext

WithContext returns a BackOffContext with context ctx

ctx must not be nil

type

Clock

¶

type Clock interface {

Now()

time

.

Time

}

Clock is an interface that returns current time for BackOff.

type

ConstantBackOff

¶

type ConstantBackOff struct {

Interval

time

.

Duration

}

ConstantBackOff is a backoff policy that always returns the same backoff delay.
This is in contrast to an exponential backoff policy,
which returns a delay that grows longer as you call NextBackOff() over and over again.

func

NewConstantBackOff

¶

func NewConstantBackOff(d

time

.

Duration

) *

ConstantBackOff

func (*ConstantBackOff)

NextBackOff

¶

func (b *

ConstantBackOff

) NextBackOff()

time

.

Duration

func (*ConstantBackOff)

Reset

¶

func (b *

ConstantBackOff

) Reset()

type

ExponentialBackOff

¶

type ExponentialBackOff struct {

InitialInterval

time

.

Duration

RandomizationFactor

float64

Multiplier

float64

MaxInterval

time

.

Duration

// After MaxElapsedTime the ExponentialBackOff returns Stop.

// It never stops if MaxElapsedTime == 0.

MaxElapsedTime

time

.

Duration

Stop

time

.

Duration

Clock

Clock

// contains filtered or unexported fields

}

ExponentialBackOff is a backoff implementation that increases the backoff
period for each retry attempt using a randomization function that grows exponentially.

NextBackOff() is calculated using the following formula:

randomized interval =
    RetryInterval * (random value in range [1 - RandomizationFactor, 1 + RandomizationFactor])

In other words NextBackOff() will range between the randomization factor
percentage below and above the retry interval.

For example, given the following parameters:

RetryInterval = 2
RandomizationFactor = 0.5
Multiplier = 2

the actual backoff period used in the next retry attempt will range between 1 and 3 seconds,
multiplied by the exponential, that is, between 2 and 6 seconds.

Note: MaxInterval caps the RetryInterval and not the randomized interval.

If the time elapsed since an ExponentialBackOff instance is created goes past the
MaxElapsedTime, then the method NextBackOff() starts returning backoff.Stop.

The elapsed time can be reset by calling Reset().

Example: Given the following default arguments, for 10 tries the sequence will be,
and assuming we go over the MaxElapsedTime on the 10th try:

Request #  RetryInterval (seconds)  Randomized Interval (seconds)

 1          0.5                     [0.25,   0.75]
 2          0.75                    [0.375,  1.125]
 3          1.125                   [0.562,  1.687]
 4          1.687                   [0.8435, 2.53]
 5          2.53                    [1.265,  3.795]
 6          3.795                   [1.897,  5.692]
 7          5.692                   [2.846,  8.538]
 8          8.538                   [4.269, 12.807]
 9         12.807                   [6.403, 19.210]
10         19.210                   backoff.Stop

Note: Implementation is not thread-safe.

func

NewExponentialBackOff

¶

func NewExponentialBackOff(opts ...

ExponentialBackOffOpts

) *

ExponentialBackOff

NewExponentialBackOff creates an instance of ExponentialBackOff using default values.

func (*ExponentialBackOff)

GetElapsedTime

¶

func (b *

ExponentialBackOff

) GetElapsedTime()

time

.

Duration

GetElapsedTime returns the elapsed time since an ExponentialBackOff instance
is created and is reset when Reset() is called.

The elapsed time is computed using time.Now().UnixNano(). It is
safe to call even while the backoff policy is used by a running
ticker.

func (*ExponentialBackOff)

NextBackOff

¶

func (b *

ExponentialBackOff

) NextBackOff()

time

.

Duration

NextBackOff calculates the next backoff interval using the formula:

Randomized interval = RetryInterval * (1 ± RandomizationFactor)

func (*ExponentialBackOff)

Reset

¶

func (b *

ExponentialBackOff

) Reset()

Reset the interval back to the initial retry interval and restarts the timer.
Reset must be called before using b.

type

ExponentialBackOffOpts

¶

added in

v4.3.0

type ExponentialBackOffOpts func(*

ExponentialBackOff

)

ExponentialBackOffOpts is a function type used to configure ExponentialBackOff options.

func

WithClockProvider

¶

added in

v4.3.0

func WithClockProvider(clock

Clock

)

ExponentialBackOffOpts

WithClockProvider sets the clock used to measure time.

func

WithInitialInterval

¶

added in

v4.3.0

func WithInitialInterval(duration

time

.

Duration

)

ExponentialBackOffOpts

WithInitialInterval sets the initial interval between retries.

func

WithMaxElapsedTime

¶

added in

v4.3.0

func WithMaxElapsedTime(duration

time

.

Duration

)

ExponentialBackOffOpts

WithMaxElapsedTime sets the maximum total time for retries.

func

WithMaxInterval

¶

added in

v4.3.0

func WithMaxInterval(duration

time

.

Duration

)

ExponentialBackOffOpts

WithMaxInterval sets the maximum interval between retries.

func

WithMultiplier

¶

added in

v4.3.0

func WithMultiplier(multiplier

float64

)

ExponentialBackOffOpts

WithMultiplier sets the multiplier for increasing the interval after each retry.

func

WithRandomizationFactor

¶

added in

v4.3.0

func WithRandomizationFactor(randomizationFactor

float64

)

ExponentialBackOffOpts

WithRandomizationFactor sets the randomization factor to add jitter to intervals.

func

WithRetryStopDuration

¶

added in

v4.3.0

func WithRetryStopDuration(duration

time

.

Duration

)

ExponentialBackOffOpts

WithRetryStopDuration sets the duration after which retries should stop.

type

Notify

¶

type Notify func(

error

,

time

.

Duration

)

Notify is a notify-on-error function. It receives an operation error and
backoff delay if the operation failed (with an error).

NOTE that if the backoff policy stated to stop retrying,
the notify function isn't called.

type

Operation

¶

type Operation func()

error

An Operation is executing by Retry() or RetryNotify().
The operation will be retried using a backoff policy if it returns an error.

type

OperationWithData

¶

added in

v4.2.0

type OperationWithData[T

any

] func() (T,

error

)

An OperationWithData is executing by RetryWithData() or RetryNotifyWithData().
The operation will be retried using a backoff policy if it returns an error.

type

PermanentError

¶

type PermanentError struct {

Err

error

}

PermanentError signals that the operation should not be retried.

func (*PermanentError)

Error

¶

func (e *

PermanentError

) Error()

string

func (*PermanentError)

Is

¶

added in

v4.1.0

func (e *

PermanentError

) Is(target

error

)

bool

func (*PermanentError)

Unwrap

¶

func (e *

PermanentError

) Unwrap()

error

type

StopBackOff

¶

type StopBackOff struct{}

StopBackOff is a fixed backoff policy that always returns backoff.Stop for
NextBackOff(), meaning that the operation should never be retried.

func (*StopBackOff)

NextBackOff

¶

func (b *

StopBackOff

) NextBackOff()

time

.

Duration

func (*StopBackOff)

Reset

¶

func (b *

StopBackOff

) Reset()

type

Ticker

¶

type Ticker struct {

C <-chan

time

.

Time

// contains filtered or unexported fields

}

Ticker holds a channel that delivers `ticks' of a clock at times reported by a BackOff.

Ticks will continue to arrive when the previous operation is still running,
so operations that take a while to fail could run in quick succession.

Example

¶

// An operation that may fail.
operation := func() error {
	return nil // or an error
}

ticker := NewTicker(NewExponentialBackOff())

var err error

// Ticks will continue to arrive when the previous operation is still running,
// so operations that take a while to fail could run in quick succession.
for range ticker.C {
	if err = operation(); err != nil {
		log.Println(err, "will retry...")
		continue
	}

	ticker.Stop()
	break
}

if err != nil {
	// Operation has failed.
	return
}

// Operation is successful.

func

NewTicker

¶

func NewTicker(b

BackOff

) *

Ticker

NewTicker returns a new Ticker containing a channel that will send
the time at times specified by the BackOff argument. Ticker is
guaranteed to tick at least once.  The channel is closed when Stop
method is called or BackOff stops. It is not safe to manipulate the
provided backoff policy (notably calling NextBackOff or Reset)
while the ticker is running.

func

NewTickerWithTimer

¶

func NewTickerWithTimer(b

BackOff

, timer

Timer

) *

Ticker

NewTickerWithTimer returns a new Ticker with a custom timer.
A default timer that uses system timer is used when nil is passed.

func (*Ticker)

Stop

¶

func (t *

Ticker

) Stop()

Stop turns off a ticker. After Stop, no more ticks will be sent.

type

Timer

¶

type Timer interface {

Start(duration

time

.

Duration

)

Stop()

C() <-chan

time

.

Time

}

type

ZeroBackOff

¶

type ZeroBackOff struct{}

ZeroBackOff is a fixed backoff policy whose backoff time is always zero,
meaning that the operation is retried immediately without waiting, indefinitely.

func (*ZeroBackOff)

NextBackOff

¶

func (b *

ZeroBackOff

) NextBackOff()

time

.

Duration

func (*ZeroBackOff)

Reset

¶

func (b *

ZeroBackOff

) Reset()