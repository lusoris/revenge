# golang.org/x/time/rate

> Source: https://pkg.go.dev/golang.org/x/time/rate
> Fetched: 2026-01-30T23:49:33.665074+00:00
> Content-Hash: 2e884c5f5caaf61d
> Type: html

---

Overview

¶

Package rate provides a rate limiter.

Index

¶

Constants

type Limit

func Every(interval time.Duration) Limit

type Limiter

func NewLimiter(r Limit, b int) *Limiter

func (lim *Limiter) Allow() bool

func (lim *Limiter) AllowN(t time.Time, n int) bool

func (lim *Limiter) Burst() int

func (lim *Limiter) Limit() Limit

func (lim *Limiter) Reserve() *Reservation

func (lim *Limiter) ReserveN(t time.Time, n int) *Reservation

func (lim *Limiter) SetBurst(newBurst int)

func (lim *Limiter) SetBurstAt(t time.Time, newBurst int)

func (lim *Limiter) SetLimit(newLimit Limit)

func (lim *Limiter) SetLimitAt(t time.Time, newLimit Limit)

func (lim *Limiter) Tokens() float64

func (lim *Limiter) TokensAt(t time.Time) float64

func (lim *Limiter) Wait(ctx context.Context) (err error)

func (lim *Limiter) WaitN(ctx context.Context, n int) (err error)

type Reservation

func (r *Reservation) Cancel()

func (r *Reservation) CancelAt(t time.Time)

func (r *Reservation) Delay() time.Duration

func (r *Reservation) DelayFrom(t time.Time) time.Duration

func (r *Reservation) OK() bool

type Sometimes

func (s *Sometimes) Do(f func())

Examples

¶

Sometimes (Every)

Sometimes (First)

Sometimes (Interval)

Sometimes (Mix)

Sometimes (Once)

Constants

¶

View Source

const Inf =

Limit

(

math

.

MaxFloat64

)

Inf is the infinite rate limit; it allows all events (even if burst is zero).

View Source

const InfDuration =

time

.

Duration

(

math

.

MaxInt64

)

InfDuration is the duration returned by Delay when a Reservation is not OK.

Variables

¶

This section is empty.

Functions

¶

This section is empty.

Types

¶

type

Limit

¶

type Limit

float64

Limit defines the maximum frequency of some events.
Limit is represented as number of events per second.
A zero Limit allows no events.

func

Every

¶

func Every(interval

time

.

Duration

)

Limit

Every converts a minimum time interval between events to a Limit.

type

Limiter

¶

type Limiter struct {

// contains filtered or unexported fields

}

A Limiter controls how frequently events are allowed to happen.
It implements a "token bucket" of size b, initially full and refilled
at rate r tokens per second.
Informally, in any large enough time interval, the Limiter limits the
rate to r tokens per second, with a maximum burst size of b events.
As a special case, if r == Inf (the infinite rate), b is ignored.
See

https://en.wikipedia.org/wiki/Token_bucket

for more about token buckets.

The zero value is a valid Limiter, but it will reject all events.
Use NewLimiter to create non-zero Limiters.

Limiter has three main methods, Allow, Reserve, and Wait.
Most callers should use Wait.

Each of the three methods consumes a single token.
They differ in their behavior when no token is available.
If no token is available, Allow returns false.
If no token is available, Reserve returns a reservation for a future token
and the amount of time the caller must wait before using it.
If no token is available, Wait blocks until one can be obtained
or its associated context.Context is canceled.

The methods AllowN, ReserveN, and WaitN consume n tokens.

Limiter is safe for simultaneous use by multiple goroutines.

func

NewLimiter

¶

func NewLimiter(r

Limit

, b

int

) *

Limiter

NewLimiter returns a new Limiter that allows events up to rate r and permits
bursts of at most b tokens.

func (*Limiter)

Allow

¶

func (lim *

Limiter

) Allow()

bool

Allow reports whether an event may happen now.

func (*Limiter)

AllowN

¶

func (lim *

Limiter

) AllowN(t

time

.

Time

, n

int

)

bool

AllowN reports whether n events may happen at time t.
Use this method if you intend to drop / skip events that exceed the rate limit.
Otherwise use Reserve or Wait.

func (*Limiter)

Burst

¶

func (lim *

Limiter

) Burst()

int

Burst returns the maximum burst size. Burst is the maximum number of tokens
that can be consumed in a single call to Allow, Reserve, or Wait, so higher
Burst values allow more events to happen at once.
A zero Burst allows no events, unless limit == Inf.

func (*Limiter)

Limit

¶

func (lim *

Limiter

) Limit()

Limit

Limit returns the maximum overall event rate.

func (*Limiter)

Reserve

¶

func (lim *

Limiter

) Reserve() *

Reservation

Reserve is shorthand for ReserveN(time.Now(), 1).

func (*Limiter)

ReserveN

¶

func (lim *

Limiter

) ReserveN(t

time

.

Time

, n

int

) *

Reservation

ReserveN returns a Reservation that indicates how long the caller must wait before n events happen.
The Limiter takes this Reservation into account when allowing future events.
The returned Reservation’s OK() method returns false if n exceeds the Limiter's burst size.
Usage example:

r := lim.ReserveN(time.Now(), 1)
if !r.OK() {
  // Not allowed to act! Did you remember to set lim.burst to be > 0 ?
  return
}
time.Sleep(r.Delay())
Act()

Use this method if you wish to wait and slow down in accordance with the rate limit without dropping events.
If you need to respect a deadline or cancel the delay, use Wait instead.
To drop or skip events exceeding rate limit, use Allow instead.

func (*Limiter)

SetBurst

¶

func (lim *

Limiter

) SetBurst(newBurst

int

)

SetBurst is shorthand for SetBurstAt(time.Now(), newBurst).

func (*Limiter)

SetBurstAt

¶

func (lim *

Limiter

) SetBurstAt(t

time

.

Time

, newBurst

int

)

SetBurstAt sets a new burst size for the limiter.

func (*Limiter)

SetLimit

¶

func (lim *

Limiter

) SetLimit(newLimit

Limit

)

SetLimit is shorthand for SetLimitAt(time.Now(), newLimit).

func (*Limiter)

SetLimitAt

¶

func (lim *

Limiter

) SetLimitAt(t

time

.

Time

, newLimit

Limit

)

SetLimitAt sets a new Limit for the limiter. The new Limit, and Burst, may be violated
or underutilized by those which reserved (using Reserve or Wait) but did not yet act
before SetLimitAt was called.

func (*Limiter)

Tokens

¶

func (lim *

Limiter

) Tokens()

float64

Tokens returns the number of tokens available now.

func (*Limiter)

TokensAt

¶

func (lim *

Limiter

) TokensAt(t

time

.

Time

)

float64

TokensAt returns the number of tokens available at time t.

func (*Limiter)

Wait

¶

func (lim *

Limiter

) Wait(ctx

context

.

Context

) (err

error

)

Wait is shorthand for WaitN(ctx, 1).

func (*Limiter)

WaitN

¶

func (lim *

Limiter

) WaitN(ctx

context

.

Context

, n

int

) (err

error

)

WaitN blocks until lim permits n events to happen.
It returns an error if n exceeds the Limiter's burst size, the Context is
canceled, or the expected wait time exceeds the Context's Deadline.
The burst limit is ignored if the rate limit is Inf.

type

Reservation

¶

type Reservation struct {

// contains filtered or unexported fields

}

A Reservation holds information about events that are permitted by a Limiter to happen after a delay.
A Reservation may be canceled, which may enable the Limiter to permit additional events.

func (*Reservation)

Cancel

¶

func (r *

Reservation

) Cancel()

Cancel is shorthand for CancelAt(time.Now()).

func (*Reservation)

CancelAt

¶

func (r *

Reservation

) CancelAt(t

time

.

Time

)

CancelAt indicates that the reservation holder will not perform the reserved action
and reverses the effects of this Reservation on the rate limit as much as possible,
considering that other reservations may have already been made.

func (*Reservation)

Delay

¶

func (r *

Reservation

) Delay()

time

.

Duration

Delay is shorthand for DelayFrom(time.Now()).

func (*Reservation)

DelayFrom

¶

func (r *

Reservation

) DelayFrom(t

time

.

Time

)

time

.

Duration

DelayFrom returns the duration for which the reservation holder must wait
before taking the reserved action.  Zero duration means act immediately.
InfDuration means the limiter cannot grant the tokens requested in this
Reservation within the maximum wait time.

func (*Reservation)

OK

¶

func (r *

Reservation

) OK()

bool

OK returns whether the limiter can provide the requested number of tokens
within the maximum wait time.  If OK is false, Delay returns InfDuration, and
Cancel does nothing.

type

Sometimes

¶

added in

v0.2.0

type Sometimes struct {

First

int

// if non-zero, the first N calls to Do will run f.

Every

int

// if non-zero, every Nth call to Do will run f.

Interval

time

.

Duration

// if non-zero and Interval has elapsed since f's last run, Do will run f.

// contains filtered or unexported fields

}

Example: logging with rate limiting

Sometimes will perform an action occasionally.  The First, Every, and
Interval fields govern the behavior of Do, which performs the action.
A zero Sometimes value will perform an action exactly once.

Example: logging with rate limiting

¶

var sometimes = rate.Sometimes{First: 3, Interval: 10*time.Second}
func Spammy() {
        sometimes.Do(func() { log.Info("here I am!") })
}

Example (Every)

¶

package main

import (
	"fmt"

	"golang.org/x/time/rate"
)

func main() {
	s := rate.Sometimes{Every: 2}
	s.Do(func() { fmt.Println("1") })
	s.Do(func() { fmt.Println("2") })
	s.Do(func() { fmt.Println("3") })
}

Output:

1
3

Share

Format

Run

Example (First)

¶

package main

import (
	"fmt"

	"golang.org/x/time/rate"
)

func main() {
	s := rate.Sometimes{First: 2}
	s.Do(func() { fmt.Println("1") })
	s.Do(func() { fmt.Println("2") })
	s.Do(func() { fmt.Println("3") })
}

Output:

1
2

Share

Format

Run

Example (Interval)

¶

package main

import (
	"fmt"
	"time"

	"golang.org/x/time/rate"
)

func main() {
	s := rate.Sometimes{Interval: 1 * time.Second}
	s.Do(func() { fmt.Println("1") })
	s.Do(func() { fmt.Println("2") })
	time.Sleep(1 * time.Second)
	s.Do(func() { fmt.Println("3") })
}

Output:

1
3

Share

Format

Run

Example (Mix)

¶

package main

import (
	"fmt"
	"time"

	"golang.org/x/time/rate"
)

func main() {
	s := rate.Sometimes{
		First:    2,
		Every:    2,
		Interval: 2 * time.Second,
	}
	s.Do(func() { fmt.Println("1 (First:2)") })
	s.Do(func() { fmt.Println("2 (First:2)") })
	s.Do(func() { fmt.Println("3 (Every:2)") })
	time.Sleep(2 * time.Second)
	s.Do(func() { fmt.Println("4 (Interval)") })
	s.Do(func() { fmt.Println("5 (Every:2)") })
	s.Do(func() { fmt.Println("6") })
}

Output:

1 (First:2)
2 (First:2)
3 (Every:2)
4 (Interval)
5 (Every:2)

Share

Format

Run

Example (Once)

¶

package main

import (
	"fmt"

	"golang.org/x/time/rate"
)

func main() {
	// The zero value of Sometimes behaves like sync.Once, though less efficiently.
	var s rate.Sometimes
	s.Do(func() { fmt.Println("1") })
	s.Do(func() { fmt.Println("2") })
	s.Do(func() { fmt.Println("3") })
}

Output:

1

Share

Format

Run

func (*Sometimes)

Do

¶

added in

v0.2.0

func (s *

Sometimes

) Do(f func())

Do runs the function f as allowed by First, Every, and Interval.

The model is a union (not intersection) of filters.  The first call to Do
always runs f.  Subsequent calls to Do run f if allowed by First or Every or
Interval.

A non-zero First:N causes the first N Do(f) calls to run f.

A non-zero Every:M causes every Mth Do(f) call, starting with the first, to
run f.

A non-zero Interval causes Do(f) to run f if Interval has elapsed since
Do last ran f.

Specifying multiple filters produces the union of these execution streams.
For example, specifying both First:N and Every:M causes the first N Do(f)
calls and every Mth Do(f) call, starting with the first, to run f.  See
Examples for more.

If Do is called multiple times simultaneously, the calls will block and run
serially.  Therefore, Do is intended for lightweight operations.

Because a call to Do may block until f returns, if f causes Do to be called,
it will deadlock.