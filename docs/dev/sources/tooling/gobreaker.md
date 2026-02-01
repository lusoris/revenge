# sony/gobreaker

> Source: https://pkg.go.dev/github.com/sony/gobreaker
> Fetched: 2026-02-01T11:42:53.967209+00:00
> Content-Hash: ef508e9d6e1f42a7
> Type: html

---

### Overview ¶

Package gobreaker implements the Circuit Breaker pattern. See <https://msdn.microsoft.com/en-us/library/dn589784.aspx>.

### Index ¶

- Variables
- type CircuitBreaker
-     * func NewCircuitBreaker(st Settings) *CircuitBreaker
-     * func (cb *CircuitBreaker) Counts() Counts
  - func (cb *CircuitBreaker) Execute(req func() (interface{}, error)) (interface{}, error)
  - func (cb *CircuitBreaker) Name() string
  - func (cb *CircuitBreaker) State() State
- type Counts
- type Settings
- type State
-     * func (s State) String() string
- type TwoStepCircuitBreaker
-     * func NewTwoStepCircuitBreaker(st Settings) *TwoStepCircuitBreaker
-     * func (tscb *TwoStepCircuitBreaker) Allow() (done func(success bool), err error)
  - func (tscb *TwoStepCircuitBreaker) Counts() Counts
  - func (tscb *TwoStepCircuitBreaker) Name() string
  - func (tscb *TwoStepCircuitBreaker) State() State

### Constants ¶

This section is empty.

### Variables ¶

[View Source](https://github.com/sony/gobreaker/blob/v1.0.0/gobreaker.go#L22)

    var (
     // ErrTooManyRequests is returned when the CB state is half open and the requests count is over the cb maxRequests
     ErrTooManyRequests = [errors](/errors).[New](/errors#New)("too many requests")
     // ErrOpenState is returned when the CB state is open
     ErrOpenState = [errors](/errors).[New](/errors#New)("circuit breaker is open")
    )

### Functions ¶

This section is empty.

### Types ¶

#### type [CircuitBreaker](https://github.com/sony/gobreaker/blob/v1.0.0/gobreaker.go#L117) ¶

    type CircuitBreaker struct {
     // contains filtered or unexported fields
    }

CircuitBreaker is a state machine to prevent sending requests that are likely to fail.

#### func [NewCircuitBreaker](https://github.com/sony/gobreaker/blob/v1.0.0/gobreaker.go#L141) ¶

    func NewCircuitBreaker(st Settings) *CircuitBreaker

NewCircuitBreaker returns a new CircuitBreaker configured with the given Settings.

#### func (*CircuitBreaker) [Counts](https://github.com/sony/gobreaker/blob/v1.0.0/gobreaker.go#L216) ¶ added in v0.5.0

    func (cb *CircuitBreaker) Counts() Counts

Counts returns internal counters

#### func (*CircuitBreaker) [Execute](https://github.com/sony/gobreaker/blob/v1.0.0/gobreaker.go#L228) ¶

    func (cb *CircuitBreaker) Execute(req func() (interface{}, [error](/builtin#error))) (interface{}, [error](/builtin#error))

Execute runs the given request if the CircuitBreaker accepts it. Execute returns an error instantly if the CircuitBreaker rejects the request. Otherwise, Execute returns the result of the request. If a panic occurs in the request, the CircuitBreaker handles it as an error and causes the same panic again.

#### func (*CircuitBreaker) [Name](https://github.com/sony/gobreaker/blob/v1.0.0/gobreaker.go#L201) ¶

    func (cb *CircuitBreaker) Name() [string](/builtin#string)

Name returns the name of the CircuitBreaker.

#### func (*CircuitBreaker) [State](https://github.com/sony/gobreaker/blob/v1.0.0/gobreaker.go#L206) ¶

    func (cb *CircuitBreaker) State() State

State returns the current state of the CircuitBreaker.

#### type [Counts](https://github.com/sony/gobreaker/blob/v1.0.0/gobreaker.go#L47) ¶

    type Counts struct {
     Requests             [uint32](/builtin#uint32)
     TotalSuccesses       [uint32](/builtin#uint32)
     TotalFailures        [uint32](/builtin#uint32)
     ConsecutiveSuccesses [uint32](/builtin#uint32)
     ConsecutiveFailures  [uint32](/builtin#uint32)
    }

Counts holds the numbers of requests and their successes/failures. CircuitBreaker clears the internal Counts either on the change of the state or at the closed-state intervals. Counts ignores the results of the requests sent before clearing.

#### type [Settings](https://github.com/sony/gobreaker/blob/v1.0.0/gobreaker.go#L106) ¶

    type Settings struct {
     Name          [string](/builtin#string)
     MaxRequests   [uint32](/builtin#uint32)
     Interval      [time](/time).[Duration](/time#Duration)
     Timeout       [time](/time).[Duration](/time#Duration)
     ReadyToTrip   func(counts Counts) [bool](/builtin#bool)
     OnStateChange func(name [string](/builtin#string), from State, to State)
     IsSuccessful  func(err [error](/builtin#error)) [bool](/builtin#bool)
    }

Settings configures CircuitBreaker:

Name is the name of the CircuitBreaker.

MaxRequests is the maximum number of requests allowed to pass through when the CircuitBreaker is half-open. If MaxRequests is 0, the CircuitBreaker allows only 1 request.

Interval is the cyclic period of the closed state for the CircuitBreaker to clear the internal Counts. If Interval is less than or equal to 0, the CircuitBreaker doesn't clear internal Counts during the closed state.

Timeout is the period of the open state, after which the state of the CircuitBreaker becomes half-open. If Timeout is less than or equal to 0, the timeout value of the CircuitBreaker is set to 60 seconds.

ReadyToTrip is called with a copy of Counts whenever a request fails in the closed state. If ReadyToTrip returns true, the CircuitBreaker will be placed into the open state. If ReadyToTrip is nil, default ReadyToTrip is used. Default ReadyToTrip returns true when the number of consecutive failures is more than 5.

OnStateChange is called whenever the state of the CircuitBreaker changes.

IsSuccessful is called with the error returned from a request. If IsSuccessful returns true, the error is counted as a success. Otherwise the error is counted as a failure. If IsSuccessful is nil, default IsSuccessful is used, which returns false for all non-nil errors.

#### type [State](https://github.com/sony/gobreaker/blob/v1.0.0/gobreaker.go#L13) ¶

    type State [int](/builtin#int)

State is a type that represents a state of CircuitBreaker.

    const (
     StateClosed State = [iota](/builtin#iota)
     StateHalfOpen
     StateOpen
    )

These constants are states of CircuitBreaker.

#### func (State) [String](https://github.com/sony/gobreaker/blob/v1.0.0/gobreaker.go#L30) ¶

    func (s State) String() [string](/builtin#string)

String implements stringer interface.

#### type [TwoStepCircuitBreaker](https://github.com/sony/gobreaker/blob/v1.0.0/gobreaker.go#L136) ¶

    type TwoStepCircuitBreaker struct {
     // contains filtered or unexported fields
    }

TwoStepCircuitBreaker is like CircuitBreaker but instead of surrounding a function with the breaker functionality, it only checks whether a request can proceed and expects the caller to report the outcome in a separate step using a callback.

#### func [NewTwoStepCircuitBreaker](https://github.com/sony/gobreaker/blob/v1.0.0/gobreaker.go#L183) ¶

    func NewTwoStepCircuitBreaker(st Settings) *TwoStepCircuitBreaker

NewTwoStepCircuitBreaker returns a new TwoStepCircuitBreaker configured with the given Settings.

#### func (*TwoStepCircuitBreaker) [Allow](https://github.com/sony/gobreaker/blob/v1.0.0/gobreaker.go#L265) ¶

    func (tscb *TwoStepCircuitBreaker) Allow() (done func(success [bool](/builtin#bool)), err [error](/builtin#error))

Allow checks if a new request can proceed. It returns a callback that should be used to register the success or failure in a separate step. If the circuit breaker doesn't allow requests, it returns an error.

#### func (*TwoStepCircuitBreaker) [Counts](https://github.com/sony/gobreaker/blob/v1.0.0/gobreaker.go#L258) ¶ added in v0.5.0

    func (tscb *TwoStepCircuitBreaker) Counts() Counts

Counts returns internal counters

#### func (*TwoStepCircuitBreaker) [Name](https://github.com/sony/gobreaker/blob/v1.0.0/gobreaker.go#L248) ¶

    func (tscb *TwoStepCircuitBreaker) Name() [string](/builtin#string)

Name returns the name of the TwoStepCircuitBreaker.

#### func (*TwoStepCircuitBreaker) [State](https://github.com/sony/gobreaker/blob/v1.0.0/gobreaker.go#L253) ¶

    func (tscb *TwoStepCircuitBreaker) State() State

State returns the current state of the TwoStepCircuitBreaker.
