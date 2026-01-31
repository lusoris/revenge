# Jaeger Go Client

> Source: https://pkg.go.dev/github.com/jaegertracing/jaeger-client-go
> Fetched: 2026-01-30T23:55:26.936789+00:00
> Content-Hash: 1b39469ec59d4527
> Type: html

---

Overview

¶

Package jaeger implements an OpenTracing (

http://opentracing.io

) Tracer.

For integration instructions please refer to the README:

https://github.com/uber/jaeger-client-go/blob/master/README.md

Index

¶

Constants

Variables

func BuildJaegerProcessThrift(span *Span) *j.Process

func BuildJaegerThrift(span *Span) *j.Span

func BuildZipkinThrift(s *Span) *z.Span

func ConvertLogsToJaegerTags(logFields []log.Field) []*j.Tag

func EnableFirehose(s *Span)

func NewTracer(serviceName string, sampler Sampler, reporter Reporter, ...) (opentracing.Tracer, io.Closer)

func SelfRef(ctx SpanContext) opentracing.SpanReference

type AdaptiveSamplerUpdater

func (u *AdaptiveSamplerUpdater) Update(sampler SamplerV2, strategy interface{}) (SamplerV2, error)

type BinaryPropagator

func NewBinaryPropagator(tracer *Tracer) *BinaryPropagator

func (p *BinaryPropagator) Extract(abstractCarrier interface{}) (SpanContext, error)

func (p *BinaryPropagator) Inject(sc SpanContext, abstractCarrier interface{}) error

type ConstSampler

func NewConstSampler(sample bool) *ConstSampler

func (s *ConstSampler) Close()

func (s *ConstSampler) Equal(other Sampler) bool

func (s *ConstSampler) IsSampled(id TraceID, operation string) (bool, []Tag)

func (s *ConstSampler) OnCreateSpan(span *Span) SamplingDecision

func (s *ConstSampler) OnFinishSpan(span *Span) SamplingDecision

func (s *ConstSampler) OnSetOperationName(span *Span, operationName string) SamplingDecision

func (s *ConstSampler) OnSetTag(span *Span, key string, value interface{}) SamplingDecision

func (s *ConstSampler) String() string

type ContribObserver

type ContribSpanObserver

type ExtractableZipkinSpan

type Extractor

type GuaranteedThroughputProbabilisticSampler

func NewGuaranteedThroughputProbabilisticSampler(lowerBound, samplingRate float64) (*GuaranteedThroughputProbabilisticSampler, error)

func (s *GuaranteedThroughputProbabilisticSampler) Close()

func (s *GuaranteedThroughputProbabilisticSampler) Equal(other Sampler) bool

func (s *GuaranteedThroughputProbabilisticSampler) IsSampled(id TraceID, operation string) (bool, []Tag)

func (s GuaranteedThroughputProbabilisticSampler) String() string

type HeadersConfig

func (c *HeadersConfig) ApplyDefaults() *HeadersConfig

type InMemoryReporter

func NewInMemoryReporter() *InMemoryReporter

func (r *InMemoryReporter) Close()

func (r *InMemoryReporter) GetSpans() []opentracing.Span

func (r *InMemoryReporter) Report(span *Span)

func (r *InMemoryReporter) Reset()

func (r *InMemoryReporter) SpansSubmitted() int

type InjectableZipkinSpan

type Injector

type Logger

type Metrics

func NewMetrics(factory metrics.Factory, globalTags map[string]string) *Metrics

func NewNullMetrics() *Metrics

type Observer

deprecated

type PerOperationSampler

func NewAdaptiveSampler(strategies *sampling.PerOperationSamplingStrategies, maxOperations int) (*PerOperationSampler, error)

func NewPerOperationSampler(params PerOperationSamplerParams) *PerOperationSampler

func (s *PerOperationSampler) Close()

func (s *PerOperationSampler) Equal(other Sampler) bool

func (s *PerOperationSampler) IsSampled(id TraceID, operation string) (bool, []Tag)

func (s *PerOperationSampler) OnCreateSpan(span *Span) SamplingDecision

func (s *PerOperationSampler) OnFinishSpan(span *Span) SamplingDecision

func (s *PerOperationSampler) OnSetOperationName(span *Span, operationName string) SamplingDecision

func (s *PerOperationSampler) OnSetTag(span *Span, key string, value interface{}) SamplingDecision

func (s *PerOperationSampler) String() string

type PerOperationSamplerParams

type ProbabilisticSampler

func NewProbabilisticSampler(samplingRate float64) (*ProbabilisticSampler, error)

func (s *ProbabilisticSampler) Close()

func (s *ProbabilisticSampler) Equal(other Sampler) bool

func (s *ProbabilisticSampler) IsSampled(id TraceID, operation string) (bool, []Tag)

func (s *ProbabilisticSampler) OnCreateSpan(span *Span) SamplingDecision

func (s *ProbabilisticSampler) OnFinishSpan(span *Span) SamplingDecision

func (s *ProbabilisticSampler) OnSetOperationName(span *Span, operationName string) SamplingDecision

func (s *ProbabilisticSampler) OnSetTag(span *Span, key string, value interface{}) SamplingDecision

func (s *ProbabilisticSampler) SamplingRate() float64

func (s *ProbabilisticSampler) String() string

func (s *ProbabilisticSampler) Update(samplingRate float64) error

type ProbabilisticSamplerUpdater

func (u *ProbabilisticSamplerUpdater) Update(sampler SamplerV2, strategy interface{}) (SamplerV2, error)

type Process

type ProcessSetter

type RateLimitingSampler

func NewRateLimitingSampler(maxTracesPerSecond float64) *RateLimitingSampler

func (s *RateLimitingSampler) Close()

func (s *RateLimitingSampler) Equal(other Sampler) bool

func (s *RateLimitingSampler) IsSampled(id TraceID, operation string) (bool, []Tag)

func (s *RateLimitingSampler) OnCreateSpan(span *Span) SamplingDecision

func (s *RateLimitingSampler) OnFinishSpan(span *Span) SamplingDecision

func (s *RateLimitingSampler) OnSetOperationName(span *Span, operationName string) SamplingDecision

func (s *RateLimitingSampler) OnSetTag(span *Span, key string, value interface{}) SamplingDecision

func (s *RateLimitingSampler) String() string

func (s *RateLimitingSampler) Update(maxTracesPerSecond float64)

type RateLimitingSamplerUpdater

func (u *RateLimitingSamplerUpdater) Update(sampler SamplerV2, strategy interface{}) (SamplerV2, error)

type Reference

type RemotelyControlledSampler

func NewRemotelyControlledSampler(serviceName string, opts ...SamplerOption) *RemotelyControlledSampler

func (s *RemotelyControlledSampler) Close()

func (s *RemotelyControlledSampler) Equal(other Sampler) bool

func (s *RemotelyControlledSampler) IsSampled(id TraceID, operation string) (bool, []Tag)

func (s *RemotelyControlledSampler) OnCreateSpan(span *Span) SamplingDecision

func (s *RemotelyControlledSampler) OnFinishSpan(span *Span) SamplingDecision

func (s *RemotelyControlledSampler) OnSetOperationName(span *Span, operationName string) SamplingDecision

func (s *RemotelyControlledSampler) OnSetTag(span *Span, key string, value interface{}) SamplingDecision

func (s *RemotelyControlledSampler) Sampler() SamplerV2

func (s *RemotelyControlledSampler) UpdateSampler()

type Reporter

func NewCompositeReporter(reporters ...Reporter) Reporter

func NewLoggingReporter(logger Logger) Reporter

func NewNullReporter() Reporter

func NewRemoteReporter(sender Transport, opts ...ReporterOption) Reporter

type ReporterOption

type Sampler

type SamplerOption

type SamplerOptionsFactory

func (SamplerOptionsFactory) InitialSampler(sampler Sampler) SamplerOption

func (SamplerOptionsFactory) Logger(logger Logger) SamplerOption

func (SamplerOptionsFactory) MaxOperations(maxOperations int) SamplerOption

func (SamplerOptionsFactory) Metrics(m *Metrics) SamplerOption

func (SamplerOptionsFactory) OperationNameLateBinding(enable bool) SamplerOption

func (SamplerOptionsFactory) SamplingRefreshInterval(samplingRefreshInterval time.Duration) SamplerOption

func (SamplerOptionsFactory) SamplingServerURL(samplingServerURL string) SamplerOption

func (SamplerOptionsFactory) SamplingStrategyFetcher(fetcher SamplingStrategyFetcher) SamplerOption

func (SamplerOptionsFactory) SamplingStrategyParser(parser SamplingStrategyParser) SamplerOption

func (SamplerOptionsFactory) Updaters(updaters ...SamplerUpdater) SamplerOption

type SamplerUpdater

type SamplerV2

type SamplerV2Base

func (SamplerV2Base) Close()

func (SamplerV2Base) Equal(other Sampler) bool

func (SamplerV2Base) IsSampled(id TraceID, operation string) (sampled bool, tags []Tag)

type SamplingDecision

type SamplingStrategyFetcher

type SamplingStrategyParser

type Span

func (s *Span) BaggageItem(key string) string

func (s *Span) Context() opentracing.SpanContext

func (s *Span) Duration() time.Duration

func (s *Span) Finish()

func (s *Span) FinishWithOptions(options opentracing.FinishOptions)

func (s *Span) Log(ld opentracing.LogData)

func (s *Span) LogEvent(event string)

func (s *Span) LogEventWithPayload(event string, payload interface{})

func (s *Span) LogFields(fields ...log.Field)

func (s *Span) LogKV(alternatingKeyValues ...interface{})

func (s *Span) Logs() []opentracing.LogRecord

func (s *Span) OperationName() string

func (s *Span) References() []opentracing.SpanReference

func (s *Span) Release()

func (s *Span) Retain() *Span

func (s *Span) SetBaggageItem(key, value string) opentracing.Span

func (s *Span) SetOperationName(operationName string) opentracing.Span

func (s *Span) SetTag(key string, value interface{}) opentracing.Span

func (s *Span) SpanContext() SpanContext

func (s *Span) StartTime() time.Time

func (s *Span) String() string

func (s *Span) Tags() opentracing.Tags

func (s *Span) Tracer() opentracing.Tracer

type SpanAllocator

type SpanContext

func ContextFromString(value string) (SpanContext, error)

func NewSpanContext(traceID TraceID, spanID, parentID SpanID, sampled bool, ...) SpanContext

func (c *SpanContext) CopyFrom(ctx *SpanContext)

func (c SpanContext) ExtendedSamplingState(key interface{}, initValue func() interface{}) interface{}

func (c SpanContext) Flags() byte

func (c SpanContext) ForeachBaggageItem(handler func(k, v string) bool)

func (c SpanContext) IsDebug() bool

func (c SpanContext) IsFirehose() bool

func (c SpanContext) IsSampled() bool

func (c SpanContext) IsSamplingFinalized() bool

func (c SpanContext) IsValid() bool

func (c SpanContext) ParentID() SpanID

func (c SpanContext) SetFirehose()

func (c SpanContext) SpanID() SpanID

func (c SpanContext) String() string

func (c SpanContext) TraceID() TraceID

func (c SpanContext) WithBaggageItem(key, value string) SpanContext

type SpanID

func SpanIDFromString(s string) (SpanID, error)

func (s SpanID) String() string

type SpanObserver

deprecated

type Tag

func NewTag(key string, value interface{}) Tag

type TextMapPropagator

func NewHTTPHeaderPropagator(headerKeys *HeadersConfig, metrics Metrics) *TextMapPropagator

func NewTextMapPropagator(headerKeys *HeadersConfig, metrics Metrics) *TextMapPropagator

func (p *TextMapPropagator) Extract(abstractCarrier interface{}) (SpanContext, error)

func (p *TextMapPropagator) Inject(sc SpanContext, abstractCarrier interface{}) error

type TraceID

func TraceIDFromString(s string) (TraceID, error)

func (t TraceID) IsValid() bool

func (t TraceID) String() string

type Tracer

func (t *Tracer) Close() error

func (t *Tracer) Extract(format interface{}, carrier interface{}) (opentracing.SpanContext, error)

func (t *Tracer) Inject(ctx opentracing.SpanContext, format interface{}, carrier interface{}) error

func (t *Tracer) Sampler() SamplerV2

func (t *Tracer) StartSpan(operationName string, options ...opentracing.StartSpanOption) opentracing.Span

func (t *Tracer) Tags() []opentracing.Tag

type TracerOption

type TracerOptionsFactory

func (TracerOptionsFactory) BaggageRestrictionManager(mgr baggage.RestrictionManager) TracerOption

func (TracerOptionsFactory) ContribObserver(observer ContribObserver) TracerOption

func (TracerOptionsFactory) CustomHeaderKeys(headerKeys *HeadersConfig) TracerOption

func (TracerOptionsFactory) DebugThrottler(throttler throttler.Throttler) TracerOption

func (TracerOptionsFactory) Extractor(format interface{}, extractor Extractor) TracerOption

func (TracerOptionsFactory) Gen128Bit(gen128Bit bool) TracerOption

func (TracerOptionsFactory) HighTraceIDGenerator(highTraceIDGenerator func() uint64) TracerOption

func (TracerOptionsFactory) HostIPv4(hostIPv4 uint32) TracerOption

func (TracerOptionsFactory) Injector(format interface{}, injector Injector) TracerOption

func (TracerOptionsFactory) Logger(logger Logger) TracerOption

func (TracerOptionsFactory) MaxLogsPerSpan(maxLogsPerSpan int) TracerOption

func (TracerOptionsFactory) MaxTagValueLength(maxTagValueLength int) TracerOption

func (TracerOptionsFactory) Metrics(m *Metrics) TracerOption

func (TracerOptionsFactory) NoDebugFlagOnForcedSampling(noDebugFlagOnForcedSampling bool) TracerOption

func (t TracerOptionsFactory) Observer(observer Observer) TracerOption

func (TracerOptionsFactory) PoolSpans(poolSpans bool) TracerOption

func (TracerOptionsFactory) RandomNumber(randomNumber func() uint64) TracerOption

func (TracerOptionsFactory) Tag(key string, value interface{}) TracerOption

func (TracerOptionsFactory) TimeNow(timeNow func() time.Time) TracerOption

func (TracerOptionsFactory) ZipkinSharedRPCSpan(zipkinSharedRPCSpan bool) TracerOption

type Transport

func NewUDPTransport(hostPort string, maxPacketSize int) (Transport, error)

func NewUDPTransportWithParams(params UDPTransportParams) (Transport, error)

type UDPTransportParams

Constants

¶

View Source

const (

// JaegerClientVersion is the version of the client library reported as Span tag.

JaegerClientVersion = "Go-2.30.0"

// JaegerClientVersionTagKey is the name of the tag used to report client version.

JaegerClientVersionTagKey = "jaeger.version"

// JaegerDebugHeader is the name of HTTP header or a TextMap carrier key which,

// if found in the carrier, forces the trace to be sampled as "debug" trace.

// The value of the header is recorded as the tag on the root span, so that the

// trace can be found in the UI using this value as a correlation ID.

JaegerDebugHeader = "jaeger-debug-id"

// JaegerBaggageHeader is the name of the HTTP header that is used to submit baggage.

// It differs from TraceBaggageHeaderPrefix in that it can be used only in cases where

// a root span does not exist.

JaegerBaggageHeader = "jaeger-baggage"

// TracerHostnameTagKey used to report host name of the process.

TracerHostnameTagKey = "hostname"

// TracerIPTagKey used to report ip of the process.

TracerIPTagKey = "ip"

// TracerUUIDTagKey used to report UUID of the client process.

TracerUUIDTagKey = "client-uuid"

// SamplerTypeTagKey reports which sampler was used on the root span.

SamplerTypeTagKey = "sampler.type"

// SamplerParamTagKey reports the parameter of the sampler, like sampling probability.

SamplerParamTagKey = "sampler.param"

// TraceContextHeaderName is the http header name used to propagate tracing context.

// This must be in lower-case to avoid mismatches when decoding incoming headers.

TraceContextHeaderName = "uber-trace-id"

// TracerStateHeaderName is deprecated.

// Deprecated: use TraceContextHeaderName

TracerStateHeaderName =

TraceContextHeaderName

// TraceBaggageHeaderPrefix is the prefix for http headers used to propagate baggage.

// This must be in lower-case to avoid mismatches when decoding incoming headers.

TraceBaggageHeaderPrefix = "uberctx-"

// SamplerTypeConst is the type of sampler that always makes the same decision.

SamplerTypeConst = "const"

// SamplerTypeRemote is the type of sampler that polls Jaeger agent for sampling strategy.

SamplerTypeRemote = "remote"

// SamplerTypeProbabilistic is the type of sampler that samples traces

// with a certain fixed probability.

SamplerTypeProbabilistic = "probabilistic"

// SamplerTypeRateLimiting is the type of sampler that samples

// only up to a fixed number of traces per second.

SamplerTypeRateLimiting = "ratelimiting"

// SamplerTypeLowerBound is the type of sampler that samples

// at least a fixed number of traces per second.

SamplerTypeLowerBound = "lowerbound"

// DefaultUDPSpanServerHost is the default host to send the spans to, via UDP

DefaultUDPSpanServerHost = "localhost"

// DefaultUDPSpanServerPort is the default port to send the spans to, via UDP

DefaultUDPSpanServerPort = 6831

// DefaultSamplingServerPort is the default port to fetch sampling config from, via http

DefaultSamplingServerPort = 5778

// DefaultMaxTagValueLength is the default max length of byte array or string allowed in the tag value.

DefaultMaxTagValueLength = 256
)

View Source

const SpanContextFormat formatKey =

iota

SpanContextFormat is a constant used as OpenTracing Format.
Requires *SpanContext as carrier.
This format is intended for interop with TChannel or other Zipkin-like tracers.

View Source

const ZipkinSpanFormat = "zipkin-span-format"

ZipkinSpanFormat is an OpenTracing carrier format constant

Variables

¶

View Source

var (

// DefaultSamplingServerURL is the default url to fetch sampling config from, via http

DefaultSamplingServerURL =

fmt

.

Sprintf

("http://127.0.0.1:%d/sampling",

DefaultSamplingServerPort

)
)

View Source

var NullLogger = &nullLogger{}

NullLogger is implementation of the Logger interface that delegates to default `log` package

View Source

var ReporterOptions reporterOptions

ReporterOptions is a factory for all available ReporterOption's

View Source

var StdLogger = &stdLogger{}

StdLogger is implementation of the Logger interface that delegates to default `log` package

Functions

¶

func

BuildJaegerProcessThrift

¶

func BuildJaegerProcessThrift(span *

Span

) *

j

.

Process

BuildJaegerProcessThrift creates a thrift Process type.
TODO: (breaking change) move to internal package.

func

BuildJaegerThrift

¶

func BuildJaegerThrift(span *

Span

) *

j

.

Span

BuildJaegerThrift builds jaeger span based on internal span.
TODO: (breaking change) move to internal package.

func

BuildZipkinThrift

¶

func BuildZipkinThrift(s *

Span

) *

z

.

Span

BuildZipkinThrift builds thrift span based on internal span.
TODO: (breaking change) move to transport/zipkin and make private.

func

ConvertLogsToJaegerTags

¶

func ConvertLogsToJaegerTags(logFields []

log

.

Field

) []*

j

.

Tag

ConvertLogsToJaegerTags converts log Fields into jaeger tags.

func

EnableFirehose

¶

func EnableFirehose(s *

Span

)

EnableFirehose enables firehose flag on the span context

func

NewTracer

¶

func NewTracer(
	serviceName

string

,
	sampler

Sampler

,
	reporter

Reporter

,
	options ...

TracerOption

,
) (

opentracing

.

Tracer

,

io

.

Closer

)

NewTracer creates Tracer implementation that reports tracing to Jaeger.
The returned io.Closer can be used in shutdown hooks to ensure that the internal
queue of the Reporter is drained and all buffered spans are submitted to collectors.
TODO (breaking change) return *Tracer only, without closer.

func

SelfRef

¶

func SelfRef(ctx

SpanContext

)

opentracing

.

SpanReference

SelfRef creates an opentracing compliant SpanReference from a jaeger
SpanContext. This is a factory function in order to encapsulate jaeger specific
types.

Types

¶

type

AdaptiveSamplerUpdater

¶

type AdaptiveSamplerUpdater struct {

MaxOperations

int

OperationNameLateBinding

bool

}

AdaptiveSamplerUpdater is used by RemotelyControlledSampler to parse sampling configuration.
Fields have the same meaning as in PerOperationSamplerParams.

func (*AdaptiveSamplerUpdater)

Update

¶

func (u *

AdaptiveSamplerUpdater

) Update(sampler

SamplerV2

, strategy interface{}) (

SamplerV2

,

error

)

Update implements Update of SamplerUpdater.

type

BinaryPropagator

¶

type BinaryPropagator struct {

// contains filtered or unexported fields

}

BinaryPropagator is a combined Injector and Extractor for Binary format

func

NewBinaryPropagator

¶

func NewBinaryPropagator(tracer *

Tracer

) *

BinaryPropagator

NewBinaryPropagator creates a combined Injector and Extractor for Binary format

func (*BinaryPropagator)

Extract

¶

func (p *

BinaryPropagator

) Extract(abstractCarrier interface{}) (

SpanContext

,

error

)

Extract implements Extractor of BinaryPropagator

func (*BinaryPropagator)

Inject

¶

func (p *

BinaryPropagator

) Inject(
	sc

SpanContext

,
	abstractCarrier interface{},
)

error

Inject implements Injector of BinaryPropagator

type

ConstSampler

¶

type ConstSampler struct {

Decision

bool

// contains filtered or unexported fields

}

ConstSampler is a sampler that always makes the same decision.

func

NewConstSampler

¶

func NewConstSampler(sample

bool

) *

ConstSampler

NewConstSampler creates a ConstSampler.

func (*ConstSampler)

Close

¶

func (s *

ConstSampler

) Close()

Close implements Close() of Sampler.

func (*ConstSampler)

Equal

¶

func (s *

ConstSampler

) Equal(other

Sampler

)

bool

Equal implements Equal() of Sampler.

func (*ConstSampler)

IsSampled

¶

func (s *

ConstSampler

) IsSampled(id

TraceID

, operation

string

) (

bool

, []

Tag

)

IsSampled implements IsSampled() of Sampler.

func (*ConstSampler)

OnCreateSpan

¶

func (s *ConstSampler) OnCreateSpan(span *

Span

)

SamplingDecision

func (*ConstSampler)

OnFinishSpan

¶

func (s *ConstSampler) OnFinishSpan(span *

Span

)

SamplingDecision

func (*ConstSampler)

OnSetOperationName

¶

func (s *ConstSampler) OnSetOperationName(span *

Span

, operationName

string

)

SamplingDecision

func (*ConstSampler)

OnSetTag

¶

func (s *ConstSampler) OnSetTag(span *

Span

, key

string

, value interface{})

SamplingDecision

func (*ConstSampler)

String

¶

func (s *

ConstSampler

) String()

string

String is used to log sampler details.

type

ContribObserver

¶

type ContribObserver interface {

// Create and return a span observer. Called when a span starts.

// If the Observer is not interested in the given span, it must return (nil, false).

// E.g :

//     func StartSpan(opName string, opts ...opentracing.StartSpanOption) {

//         var sp opentracing.Span

//         sso := opentracing.StartSpanOptions{}

//         if spanObserver, ok := Observer.OnStartSpan(span, opName, sso); ok {

//             // we have a valid SpanObserver

//         }

//         ...

//     }

OnStartSpan(sp

opentracing

.

Span

, operationName

string

, options

opentracing

.

StartSpanOptions

) (

ContribSpanObserver

,

bool

)
}

ContribObserver can be registered with the Tracer to receive notifications
about new Spans. Modelled after github.com/opentracing-contrib/go-observer.

type

ContribSpanObserver

¶

type ContribSpanObserver interface {

OnSetOperationName(operationName

string

)

OnSetTag(key

string

, value interface{})

OnFinish(options

opentracing

.

FinishOptions

)

}

ContribSpanObserver is created by the Observer and receives notifications
about other Span events. This interface is meant to match
github.com/opentracing-contrib/go-observer, via duck typing, without
directly importing the go-observer package.

type

ExtractableZipkinSpan

¶

type ExtractableZipkinSpan interface {

TraceID()

uint64

SpanID()

uint64

ParentID()

uint64

Flags()

byte

}

ExtractableZipkinSpan is a type of Carrier used for integration with Zipkin-aware
RPC frameworks (like TChannel). It does not support baggage, only trace IDs.

type

Extractor

¶

type Extractor interface {

// Extract decodes a SpanContext instance from the given `carrier`,

// or (nil, opentracing.ErrSpanContextNotFound) if no context could

// be found in the `carrier`.

Extract(carrier interface{}) (

SpanContext

,

error

)
}

Extractor is responsible for extracting SpanContext instances from a
format-specific "carrier" object. Typically the extraction will take place
on the server side of an RPC boundary, but message queues and other IPC
mechanisms are also reasonable places to use an Extractor.

type

GuaranteedThroughputProbabilisticSampler

¶

type GuaranteedThroughputProbabilisticSampler struct {

// contains filtered or unexported fields

}

GuaranteedThroughputProbabilisticSampler is a sampler that leverages both ProbabilisticSampler and
RateLimitingSampler. The RateLimitingSampler is used as a guaranteed lower bound sampler such that
every operation is sampled at least once in a time interval defined by the lowerBound. ie a lowerBound
of 1.0 / (60 * 10) will sample an operation at least once every 10 minutes.

The ProbabilisticSampler is given higher priority when tags are emitted, ie. if IsSampled() for both
samplers return true, the tags for ProbabilisticSampler will be used.

func

NewGuaranteedThroughputProbabilisticSampler

¶

func NewGuaranteedThroughputProbabilisticSampler(
	lowerBound, samplingRate

float64

,
) (*

GuaranteedThroughputProbabilisticSampler

,

error

)

NewGuaranteedThroughputProbabilisticSampler returns a delegating sampler that applies both
ProbabilisticSampler and RateLimitingSampler.

func (*GuaranteedThroughputProbabilisticSampler)

Close

¶

func (s *

GuaranteedThroughputProbabilisticSampler

) Close()

Close implements Close() of Sampler.

func (*GuaranteedThroughputProbabilisticSampler)

Equal

¶

func (s *

GuaranteedThroughputProbabilisticSampler

) Equal(other

Sampler

)

bool

Equal implements Equal() of Sampler.

func (*GuaranteedThroughputProbabilisticSampler)

IsSampled

¶

func (s *

GuaranteedThroughputProbabilisticSampler

) IsSampled(id

TraceID

, operation

string

) (

bool

, []

Tag

)

IsSampled implements IsSampled() of Sampler.

func (GuaranteedThroughputProbabilisticSampler)

String

¶

func (s

GuaranteedThroughputProbabilisticSampler

) String()

string

type

HeadersConfig

¶

type HeadersConfig struct {

// JaegerDebugHeader is the name of HTTP header or a TextMap carrier key which,

// if found in the carrier, forces the trace to be sampled as "debug" trace.

// The value of the header is recorded as the tag on the root span, so that the

// trace can be found in the UI using this value as a correlation ID.

JaegerDebugHeader

string

`yaml:"jaegerDebugHeader"`

// JaegerBaggageHeader is the name of the HTTP header that is used to submit baggage.

// It differs from TraceBaggageHeaderPrefix in that it can be used only in cases where

// a root span does not exist.

JaegerBaggageHeader

string

`yaml:"jaegerBaggageHeader"`

// TraceContextHeaderName is the http header name used to propagate tracing context.

// This must be in lower-case to avoid mismatches when decoding incoming headers.

TraceContextHeaderName

string

`yaml:"TraceContextHeaderName"`

// TraceBaggageHeaderPrefix is the prefix for http headers used to propagate baggage.

// This must be in lower-case to avoid mismatches when decoding incoming headers.

TraceBaggageHeaderPrefix

string

`yaml:"traceBaggageHeaderPrefix"`
}

HeadersConfig contains the values for the header keys that Jaeger will use.
These values may be either custom or default depending on whether custom
values were provided via a configuration.

func (*HeadersConfig)

ApplyDefaults

¶

func (c *

HeadersConfig

) ApplyDefaults() *

HeadersConfig

ApplyDefaults sets missing configuration keys to default values

type

InMemoryReporter

¶

type InMemoryReporter struct {

// contains filtered or unexported fields

}

InMemoryReporter is used for testing, and simply collects spans in memory.

func

NewInMemoryReporter

¶

func NewInMemoryReporter() *

InMemoryReporter

NewInMemoryReporter creates a reporter that stores spans in memory.
NOTE: the Tracer should be created with options.PoolSpans = false.

func (*InMemoryReporter)

Close

¶

func (r *

InMemoryReporter

) Close()

Close implements Close() method of Reporter

func (*InMemoryReporter)

GetSpans

¶

func (r *

InMemoryReporter

) GetSpans() []

opentracing

.

Span

GetSpans returns accumulated spans as a copy of the buffer.

func (*InMemoryReporter)

Report

¶

func (r *

InMemoryReporter

) Report(span *

Span

)

Report implements Report() method of Reporter by storing the span in the buffer.

func (*InMemoryReporter)

Reset

¶

func (r *

InMemoryReporter

) Reset()

Reset clears all accumulated spans.

func (*InMemoryReporter)

SpansSubmitted

¶

func (r *

InMemoryReporter

) SpansSubmitted()

int

SpansSubmitted returns the number of spans accumulated in the buffer.

type

InjectableZipkinSpan

¶

type InjectableZipkinSpan interface {

SetTraceID(traceID

uint64

)

SetSpanID(spanID

uint64

)

SetParentID(parentID

uint64

)

SetFlags(flags

byte

)

}

InjectableZipkinSpan is a type of Carrier used for integration with Zipkin-aware
RPC frameworks (like TChannel). It does not support baggage, only trace IDs.

type

Injector

¶

type Injector interface {

// Inject takes `SpanContext` and injects it into `carrier`. The actual type

// of `carrier` depends on the `format` passed to `Tracer.Inject()`.

//

// Implementations may return opentracing.ErrInvalidCarrier or any other

// implementation-specific error if injection fails.

Inject(ctx

SpanContext

, carrier interface{})

error

}

Injector is responsible for injecting SpanContext instances in a manner suitable
for propagation via a format-specific "carrier" object. Typically the
injection will take place across an RPC boundary, but message queues and
other IPC mechanisms are also reasonable places to use an Injector.

type

Logger

¶

type Logger interface {

// Error logs a message at error priority

Error(msg

string

)

// Infof logs a message at info priority

Infof(msg

string

, args ...interface{})
}

Logger provides an abstract interface for logging from Reporters.
Applications can provide their own implementation of this interface to adapt
reporters logging to whatever logging library they prefer (stdlib log,
logrus, go-logging, etc).

type

Metrics

¶

type Metrics struct {

// Number of traces started by this tracer as sampled

TracesStartedSampled

metrics

.

Counter

`metric:"traces" tags:"state=started,sampled=y" help:"Number of traces started by this tracer as sampled"`

// Number of traces started by this tracer as not sampled

TracesStartedNotSampled

metrics

.

Counter

`metric:"traces" tags:"state=started,sampled=n" help:"Number of traces started by this tracer as not sampled"`

// Number of traces started by this tracer with delayed sampling

TracesStartedDelayedSampling

metrics

.

Counter

`metric:"traces" tags:"state=started,sampled=n" help:"Number of traces started by this tracer with delayed sampling"`

// Number of externally started sampled traces this tracer joined

TracesJoinedSampled

metrics

.

Counter

`metric:"traces" tags:"state=joined,sampled=y" help:"Number of externally started sampled traces this tracer joined"`

// Number of externally started not-sampled traces this tracer joined

TracesJoinedNotSampled

metrics

.

Counter

`metric:"traces" tags:"state=joined,sampled=n" help:"Number of externally started not-sampled traces this tracer joined"`

// Number of sampled spans started by this tracer

SpansStartedSampled

metrics

.

Counter

`metric:"started_spans" tags:"sampled=y" help:"Number of spans started by this tracer as sampled"`

// Number of not sampled spans started by this tracer

SpansStartedNotSampled

metrics

.

Counter

`metric:"started_spans" tags:"sampled=n" help:"Number of spans started by this tracer as not sampled"`

// Number of spans with delayed sampling started by this tracer

SpansStartedDelayedSampling

metrics

.

Counter

`metric:"started_spans" tags:"sampled=delayed" help:"Number of spans started by this tracer with delayed sampling"`

// Number of spans finished by this tracer

SpansFinishedSampled

metrics

.

Counter

`metric:"finished_spans" tags:"sampled=y" help:"Number of sampled spans finished by this tracer"`

// Number of spans finished by this tracer

SpansFinishedNotSampled

metrics

.

Counter

`metric:"finished_spans" tags:"sampled=n" help:"Number of not-sampled spans finished by this tracer"`

// Number of spans finished by this tracer

SpansFinishedDelayedSampling

metrics

.

Counter

`metric:"finished_spans" tags:"sampled=delayed" help:"Number of spans with delayed sampling finished by this tracer"`

// Number of errors decoding tracing context

DecodingErrors

metrics

.

Counter

`metric:"span_context_decoding_errors" help:"Number of errors decoding tracing context"`

// Number of spans successfully reported

ReporterSuccess

metrics

.

Counter

`metric:"reporter_spans" tags:"result=ok" help:"Number of spans successfully reported"`

// Number of spans not reported due to a Sender failure

ReporterFailure

metrics

.

Counter

`metric:"reporter_spans" tags:"result=err" help:"Number of spans not reported due to a Sender failure"`

// Number of spans dropped due to internal queue overflow

ReporterDropped

metrics

.

Counter

`metric:"reporter_spans" tags:"result=dropped" help:"Number of spans dropped due to internal queue overflow"`

// Current number of spans in the reporter queue

ReporterQueueLength

metrics

.

Gauge

`metric:"reporter_queue_length" help:"Current number of spans in the reporter queue"`

// Number of times the Sampler succeeded to retrieve sampling strategy

SamplerRetrieved

metrics

.

Counter

`metric:"sampler_queries" tags:"result=ok" help:"Number of times the Sampler succeeded to retrieve sampling strategy"`

// Number of times the Sampler failed to retrieve sampling strategy

SamplerQueryFailure

metrics

.

Counter

`metric:"sampler_queries" tags:"result=err" help:"Number of times the Sampler failed to retrieve sampling strategy"`

// Number of times the Sampler succeeded to retrieve and update sampling strategy

SamplerUpdated

metrics

.

Counter

``

/* 127-byte string literal not displayed */

// Number of times the Sampler failed to update sampling strategy

SamplerUpdateFailure

metrics

.

Counter

`metric:"sampler_updates" tags:"result=err" help:"Number of times the Sampler failed to update sampling strategy"`

// Number of times baggage was successfully written or updated on spans.

BaggageUpdateSuccess

metrics

.

Counter

`metric:"baggage_updates" tags:"result=ok" help:"Number of times baggage was successfully written or updated on spans"`

// Number of times baggage failed to write or update on spans.

BaggageUpdateFailure

metrics

.

Counter

`metric:"baggage_updates" tags:"result=err" help:"Number of times baggage failed to write or update on spans"`

// Number of times baggage was truncated as per baggage restrictions.

BaggageTruncate

metrics

.

Counter

`metric:"baggage_truncations" help:"Number of times baggage was truncated as per baggage restrictions"`

// Number of times baggage restrictions were successfully updated.

BaggageRestrictionsUpdateSuccess

metrics

.

Counter

`metric:"baggage_restrictions_updates" tags:"result=ok" help:"Number of times baggage restrictions were successfully updated"`

// Number of times baggage restrictions failed to update.

BaggageRestrictionsUpdateFailure

metrics

.

Counter

`metric:"baggage_restrictions_updates" tags:"result=err" help:"Number of times baggage restrictions failed to update"`

// Number of times debug spans were throttled.

ThrottledDebugSpans

metrics

.

Counter

`metric:"throttled_debug_spans" help:"Number of times debug spans were throttled"`

// Number of times throttler successfully updated.

ThrottlerUpdateSuccess

metrics

.

Counter

`metric:"throttler_updates" tags:"result=ok" help:"Number of times throttler successfully updated"`

// Number of times throttler failed to update.

ThrottlerUpdateFailure

metrics

.

Counter

`metric:"throttler_updates" tags:"result=err" help:"Number of times throttler failed to update"`
}

Metrics is a container of all stats emitted by Jaeger tracer.

func

NewMetrics

¶

func NewMetrics(factory

metrics

.

Factory

, globalTags map[

string

]

string

) *

Metrics

NewMetrics creates a new Metrics struct and initializes it.

func

NewNullMetrics

¶

func NewNullMetrics() *

Metrics

NewNullMetrics creates a new Metrics struct that won't report any metrics.

type

Observer

deprecated

type Observer interface {

OnStartSpan(operationName

string

, options

opentracing

.

StartSpanOptions

)

SpanObserver

}

Observer can be registered with the Tracer to receive notifications about
new Spans.

Deprecated: use jaeger.ContribObserver instead.

type

PerOperationSampler

¶

type PerOperationSampler struct {

sync

.

RWMutex

// contains filtered or unexported fields

}

PerOperationSampler is a delegating sampler that applies GuaranteedThroughputProbabilisticSampler
on a per-operation basis.

func

NewAdaptiveSampler

¶

func NewAdaptiveSampler(strategies *

sampling

.

PerOperationSamplingStrategies

, maxOperations

int

) (*

PerOperationSampler

,

error

)

NewAdaptiveSampler returns a new PerOperationSampler.
Deprecated: please use NewPerOperationSampler.

func

NewPerOperationSampler

¶

func NewPerOperationSampler(params

PerOperationSamplerParams

) *

PerOperationSampler

NewPerOperationSampler returns a new PerOperationSampler.

func (*PerOperationSampler)

Close

¶

func (s *

PerOperationSampler

) Close()

Close invokes Close on all underlying samplers.

func (*PerOperationSampler)

Equal

¶

func (s *

PerOperationSampler

) Equal(other

Sampler

)

bool

Equal is not used.
TODO (breaking change) remove this in the future

func (*PerOperationSampler)

IsSampled

¶

func (s *

PerOperationSampler

) IsSampled(id

TraceID

, operation

string

) (

bool

, []

Tag

)

IsSampled is not used and only exists to match Sampler V1 API.
TODO (breaking change) remove when upgrading everything to SamplerV2

func (*PerOperationSampler)

OnCreateSpan

¶

func (s *

PerOperationSampler

) OnCreateSpan(span *

Span

)

SamplingDecision

OnCreateSpan implements OnCreateSpan of SamplerV2.

func (*PerOperationSampler)

OnFinishSpan

¶

func (s *

PerOperationSampler

) OnFinishSpan(span *

Span

)

SamplingDecision

OnFinishSpan implements OnFinishSpan of SamplerV2.

func (*PerOperationSampler)

OnSetOperationName

¶

func (s *

PerOperationSampler

) OnSetOperationName(span *

Span

, operationName

string

)

SamplingDecision

OnSetOperationName implements OnSetOperationName of SamplerV2.

func (*PerOperationSampler)

OnSetTag

¶

func (s *

PerOperationSampler

) OnSetTag(span *

Span

, key

string

, value interface{})

SamplingDecision

OnSetTag implements OnSetTag of SamplerV2.

func (*PerOperationSampler)

String

¶

func (s *

PerOperationSampler

) String()

string

type

PerOperationSamplerParams

¶

type PerOperationSamplerParams struct {

// Max number of operations that will be tracked. Other operations will be given default strategy.

MaxOperations

int

// Opt-in feature for applications that require late binding of span name via explicit call to SetOperationName.

// When this feature is enabled, the sampler will return retryable=true from OnCreateSpan(), thus leaving

// the sampling decision as non-final (and the span as writeable). This may lead to degraded performance

// in applications that always provide the correct span name on trace creation.

//

// For backwards compatibility this option is off by default.

OperationNameLateBinding

bool

// Initial configuration of the sampling strategies (usually retrieved from the backend by Remote Sampler).

Strategies *

sampling

.

PerOperationSamplingStrategies

}

PerOperationSamplerParams defines parameters when creating PerOperationSampler.

type

ProbabilisticSampler

¶

type ProbabilisticSampler struct {

// contains filtered or unexported fields

}

ProbabilisticSampler is a sampler that randomly samples a certain percentage
of traces.

func

NewProbabilisticSampler

¶

func NewProbabilisticSampler(samplingRate

float64

) (*

ProbabilisticSampler

,

error

)

NewProbabilisticSampler creates a sampler that randomly samples a certain percentage of traces specified by the
samplingRate, in the range between 0.0 and 1.0.

It relies on the fact that new trace IDs are 63bit random numbers themselves, thus making the sampling decision
without generating a new random number, but simply calculating if traceID < (samplingRate * 2^63).
TODO remove the error from this function for next major release

func (*ProbabilisticSampler)

Close

¶

func (s *

ProbabilisticSampler

) Close()

Close implements Close() of Sampler.

func (*ProbabilisticSampler)

Equal

¶

func (s *

ProbabilisticSampler

) Equal(other

Sampler

)

bool

Equal implements Equal() of Sampler.

func (*ProbabilisticSampler)

IsSampled

¶

func (s *

ProbabilisticSampler

) IsSampled(id

TraceID

, operation

string

) (

bool

, []

Tag

)

IsSampled implements IsSampled() of Sampler.

func (*ProbabilisticSampler)

OnCreateSpan

¶

func (s *ProbabilisticSampler) OnCreateSpan(span *

Span

)

SamplingDecision

func (*ProbabilisticSampler)

OnFinishSpan

¶

func (s *ProbabilisticSampler) OnFinishSpan(span *

Span

)

SamplingDecision

func (*ProbabilisticSampler)

OnSetOperationName

¶

func (s *ProbabilisticSampler) OnSetOperationName(span *

Span

, operationName

string

)

SamplingDecision

func (*ProbabilisticSampler)

OnSetTag

¶

func (s *ProbabilisticSampler) OnSetTag(span *

Span

, key

string

, value interface{})

SamplingDecision

func (*ProbabilisticSampler)

SamplingRate

¶

func (s *

ProbabilisticSampler

) SamplingRate()

float64

SamplingRate returns the sampling probability this sampled was constructed with.

func (*ProbabilisticSampler)

String

¶

func (s *

ProbabilisticSampler

) String()

string

String is used to log sampler details.

func (*ProbabilisticSampler)

Update

¶

func (s *

ProbabilisticSampler

) Update(samplingRate

float64

)

error

Update modifies in-place the sampling rate. Locking must be done externally.

type

ProbabilisticSamplerUpdater

¶

type ProbabilisticSamplerUpdater struct{}

ProbabilisticSamplerUpdater is used by RemotelyControlledSampler to parse sampling configuration.

func (*ProbabilisticSamplerUpdater)

Update

¶

func (u *

ProbabilisticSamplerUpdater

) Update(sampler

SamplerV2

, strategy interface{}) (

SamplerV2

,

error

)

Update implements Update of SamplerUpdater.

type

Process

¶

type Process struct {

Service

string

UUID

string

Tags    []

Tag

}

Process holds process specific metadata that's relevant to this client.

type

ProcessSetter

¶

type ProcessSetter interface {

SetProcess(process

Process

)

}

ProcessSetter sets a process. This can be used by any class that requires
the process to be set as part of initialization.
See internal/throttler/remote/throttler.go for an example.

type

RateLimitingSampler

¶

type RateLimitingSampler struct {

// contains filtered or unexported fields

}

RateLimitingSampler samples at most maxTracesPerSecond. The distribution of sampled traces follows
burstiness of the service, i.e. a service with uniformly distributed requests will have those
requests sampled uniformly as well, but if requests are bursty, especially sub-second, then a
number of sequential requests can be sampled each second.

func

NewRateLimitingSampler

¶

func NewRateLimitingSampler(maxTracesPerSecond

float64

) *

RateLimitingSampler

NewRateLimitingSampler creates new RateLimitingSampler.

func (*RateLimitingSampler)

Close

¶

func (s *

RateLimitingSampler

) Close()

Close does nothing.

func (*RateLimitingSampler)

Equal

¶

func (s *

RateLimitingSampler

) Equal(other

Sampler

)

bool

Equal compares with another sampler.

func (*RateLimitingSampler)

IsSampled

¶

func (s *

RateLimitingSampler

) IsSampled(id

TraceID

, operation

string

) (

bool

, []

Tag

)

IsSampled implements IsSampled() of Sampler.

func (*RateLimitingSampler)

OnCreateSpan

¶

func (s *RateLimitingSampler) OnCreateSpan(span *

Span

)

SamplingDecision

func (*RateLimitingSampler)

OnFinishSpan

¶

func (s *RateLimitingSampler) OnFinishSpan(span *

Span

)

SamplingDecision

func (*RateLimitingSampler)

OnSetOperationName

¶

func (s *RateLimitingSampler) OnSetOperationName(span *

Span

, operationName

string

)

SamplingDecision

func (*RateLimitingSampler)

OnSetTag

¶

func (s *RateLimitingSampler) OnSetTag(span *

Span

, key

string

, value interface{})

SamplingDecision

func (*RateLimitingSampler)

String

¶

func (s *

RateLimitingSampler

) String()

string

String is used to log sampler details.

func (*RateLimitingSampler)

Update

¶

func (s *

RateLimitingSampler

) Update(maxTracesPerSecond

float64

)

Update reconfigures the rate limiter, while preserving its accumulated balance.
Locking must be done externally.

type

RateLimitingSamplerUpdater

¶

type RateLimitingSamplerUpdater struct{}

RateLimitingSamplerUpdater is used by RemotelyControlledSampler to parse sampling configuration.

func (*RateLimitingSamplerUpdater)

Update

¶

func (u *

RateLimitingSamplerUpdater

) Update(sampler

SamplerV2

, strategy interface{}) (

SamplerV2

,

error

)

Update implements Update of SamplerUpdater.

type

Reference

¶

type Reference struct {

Type

opentracing

.

SpanReferenceType

Context

SpanContext

}

Reference represents a causal reference to other Spans (via their SpanContext).

type

RemotelyControlledSampler

¶

type RemotelyControlledSampler struct {

sync

.

RWMutex

// used to serialize access to samplerOptions.sampler

// contains filtered or unexported fields

}

RemotelyControlledSampler is a delegating sampler that polls a remote server
for the appropriate sampling strategy, constructs a corresponding sampler and
delegates to it for sampling decisions.

func

NewRemotelyControlledSampler

¶

func NewRemotelyControlledSampler(
	serviceName

string

,
	opts ...

SamplerOption

,
) *

RemotelyControlledSampler

NewRemotelyControlledSampler creates a sampler that periodically pulls
the sampling strategy from an HTTP sampling server (e.g. jaeger-agent).

func (*RemotelyControlledSampler)

Close

¶

func (s *

RemotelyControlledSampler

) Close()

Close implements Close() of Sampler.

func (*RemotelyControlledSampler)

Equal

¶

func (s *

RemotelyControlledSampler

) Equal(other

Sampler

)

bool

Equal implements Equal() of Sampler.

func (*RemotelyControlledSampler)

IsSampled

¶

func (s *

RemotelyControlledSampler

) IsSampled(id

TraceID

, operation

string

) (

bool

, []

Tag

)

IsSampled implements IsSampled() of Sampler.
TODO (breaking change) remove when Sampler V1 is removed

func (*RemotelyControlledSampler)

OnCreateSpan

¶

func (s *

RemotelyControlledSampler

) OnCreateSpan(span *

Span

)

SamplingDecision

OnCreateSpan implements OnCreateSpan of SamplerV2.

func (*RemotelyControlledSampler)

OnFinishSpan

¶

func (s *

RemotelyControlledSampler

) OnFinishSpan(span *

Span

)

SamplingDecision

OnFinishSpan implements OnFinishSpan of SamplerV2.

func (*RemotelyControlledSampler)

OnSetOperationName

¶

func (s *

RemotelyControlledSampler

) OnSetOperationName(span *

Span

, operationName

string

)

SamplingDecision

OnSetOperationName implements OnSetOperationName of SamplerV2.

func (*RemotelyControlledSampler)

OnSetTag

¶

func (s *

RemotelyControlledSampler

) OnSetTag(span *

Span

, key

string

, value interface{})

SamplingDecision

OnSetTag implements OnSetTag of SamplerV2.

func (*RemotelyControlledSampler)

Sampler

¶

func (s *

RemotelyControlledSampler

) Sampler()

SamplerV2

Sampler returns the currently active sampler.

func (*RemotelyControlledSampler)

UpdateSampler

¶

func (s *

RemotelyControlledSampler

) UpdateSampler()

UpdateSampler forces the sampler to fetch sampling strategy from backend server.
This function is called automatically on a timer, but can also be safely called manually, e.g. from tests.

type

Reporter

¶

type Reporter interface {

// Report submits a new span to collectors, possibly asynchronously and/or with buffering.

// If the reporter is processing Span asynchronously then it needs to Retain() the span,

// and then Release() it when no longer needed, to avoid span data corruption.

Report(span *

Span

)

// Close does a clean shutdown of the reporter, flushing any traces that may be buffered in memory.

Close()
}

Reporter is called by the tracer when a span is completed to report the span to the tracing collector.

func

NewCompositeReporter

¶

func NewCompositeReporter(reporters ...

Reporter

)

Reporter

NewCompositeReporter creates a reporter that ignores all reported spans.

func

NewLoggingReporter

¶

func NewLoggingReporter(logger

Logger

)

Reporter

NewLoggingReporter creates a reporter that logs all reported spans to provided logger.

func

NewNullReporter

¶

func NewNullReporter()

Reporter

NewNullReporter creates a no-op reporter that ignores all reported spans.

func

NewRemoteReporter

¶

func NewRemoteReporter(sender

Transport

, opts ...

ReporterOption

)

Reporter

NewRemoteReporter creates a new reporter that sends spans out of process by means of Sender.
Calls to Report(Span) return immediately (side effect: if internal buffer is full the span is dropped).
Periodically the transport buffer is flushed even if it hasn't reached max packet size.
Calls to Close() block until all spans reported prior to the call to Close are flushed.

type

ReporterOption

¶

type ReporterOption func(c *reporterOptions)

ReporterOption is a function that sets some option on the reporter.

type

Sampler

¶

type Sampler interface {

// IsSampled decides whether a trace with given `id` and `operation`

// should be sampled. This function will also return the tags that

// can be used to identify the type of sampling that was applied to

// the root span. Most simple samplers would return two tags,

// sampler.type and sampler.param, similar to those used in the Configuration

IsSampled(id

TraceID

, operation

string

) (sampled

bool

, tags []

Tag

)

// Close does a clean shutdown of the sampler, stopping any background

// go-routines it may have started.

Close()

// Equal checks if the `other` sampler is functionally equivalent

// to this sampler.

// TODO (breaking change) remove this function. See PerOperationSampler.Equals for explanation.

Equal(other

Sampler

)

bool

}

Sampler decides whether a new trace should be sampled or not.

type

SamplerOption

¶

type SamplerOption func(options *samplerOptions)

SamplerOption is a function that sets some option on the sampler

type

SamplerOptionsFactory

¶

type SamplerOptionsFactory struct{}

SamplerOptionsFactory is a factory for all available SamplerOption's.
The type acts as a namespace for factory functions. It is public to
make the functions discoverable via godoc. Recommended to be used
via global SamplerOptions variable.

var SamplerOptions

SamplerOptionsFactory

SamplerOptions is a factory for all available SamplerOption's.

func (SamplerOptionsFactory)

InitialSampler

¶

func (

SamplerOptionsFactory

) InitialSampler(sampler

Sampler

)

SamplerOption

InitialSampler creates a SamplerOption that sets the initial sampler
to use before a remote sampler is created and used.

func (SamplerOptionsFactory)

Logger

¶

func (

SamplerOptionsFactory

) Logger(logger

Logger

)

SamplerOption

Logger creates a SamplerOption that sets the logger used by the sampler.

func (SamplerOptionsFactory)

MaxOperations

¶

func (

SamplerOptionsFactory

) MaxOperations(maxOperations

int

)

SamplerOption

MaxOperations creates a SamplerOption that sets the maximum number of
operations the sampler will keep track of.

func (SamplerOptionsFactory)

Metrics

¶

func (

SamplerOptionsFactory

) Metrics(m *

Metrics

)

SamplerOption

Metrics creates a SamplerOption that initializes Metrics on the sampler,
which is used to emit statistics.

func (SamplerOptionsFactory)

OperationNameLateBinding

¶

func (

SamplerOptionsFactory

) OperationNameLateBinding(enable

bool

)

SamplerOption

OperationNameLateBinding creates a SamplerOption that sets the respective
field in the PerOperationSamplerParams.

func (SamplerOptionsFactory)

SamplingRefreshInterval

¶

func (

SamplerOptionsFactory

) SamplingRefreshInterval(samplingRefreshInterval

time

.

Duration

)

SamplerOption

SamplingRefreshInterval creates a SamplerOption that sets how often the
sampler will poll local agent for the appropriate sampling strategy.

func (SamplerOptionsFactory)

SamplingServerURL

¶

func (

SamplerOptionsFactory

) SamplingServerURL(samplingServerURL

string

)

SamplerOption

SamplingServerURL creates a SamplerOption that sets the sampling server url
of the local agent that contains the sampling strategies.

func (SamplerOptionsFactory)

SamplingStrategyFetcher

¶

func (

SamplerOptionsFactory

) SamplingStrategyFetcher(fetcher

SamplingStrategyFetcher

)

SamplerOption

SamplingStrategyFetcher creates a SamplerOption that initializes sampling strategy fetcher.

func (SamplerOptionsFactory)

SamplingStrategyParser

¶

func (

SamplerOptionsFactory

) SamplingStrategyParser(parser

SamplingStrategyParser

)

SamplerOption

SamplingStrategyParser creates a SamplerOption that initializes sampling strategy parser.

func (SamplerOptionsFactory)

Updaters

¶

func (

SamplerOptionsFactory

) Updaters(updaters ...

SamplerUpdater

)

SamplerOption

Updaters creates a SamplerOption that initializes sampler updaters.

type

SamplerUpdater

¶

type SamplerUpdater interface {

Update(sampler

SamplerV2

, strategy interface{}) (modified

SamplerV2

, err

error

)

}

SamplerUpdater is used by RemotelyControlledSampler to apply sampling strategies,
retrieved from remote config server, to the current sampler. The updater can modify
the sampler in-place if sampler supports it, or create a new one.

If the strategy does not contain configuration for the sampler in question,
updater must return modifiedSampler=nil to give other updaters a chance to inspect
the sampling strategy response.

RemotelyControlledSampler invokes the updaters while holding a lock on the main sampler.

type

SamplerV2

¶

type SamplerV2 interface {

OnCreateSpan(span *

Span

)

SamplingDecision

OnSetOperationName(span *

Span

, operationName

string

)

SamplingDecision

OnSetTag(span *

Span

, key

string

, value interface{})

SamplingDecision

OnFinishSpan(span *

Span

)

SamplingDecision

// Close does a clean shutdown of the sampler, stopping any background

// go-routines it may have started.

Close()
}

SamplerV2 is an extension of the V1 samplers that allows sampling decisions
be made at different points of the span lifecycle.

type

SamplerV2Base

¶

type SamplerV2Base struct{}

SamplerV2Base can be used by V2 samplers to implement dummy V1 methods.
Supporting V1 API is required because Tracer configuration only accepts V1 Sampler
for backwards compatibility reasons.
TODO (breaking change) remove this in the next major release

func (SamplerV2Base)

Close

¶

func (

SamplerV2Base

) Close()

Close implements Close of Sampler.

func (SamplerV2Base)

Equal

¶

func (

SamplerV2Base

) Equal(other

Sampler

)

bool

Equal implements Equal of Sampler.

func (SamplerV2Base)

IsSampled

¶

func (

SamplerV2Base

) IsSampled(id

TraceID

, operation

string

) (sampled

bool

, tags []

Tag

)

IsSampled implements IsSampled of Sampler.

type

SamplingDecision

¶

type SamplingDecision struct {

Sample

bool

Retryable

bool

Tags      []

Tag

}

SamplingDecision is returned by the V2 samplers.

type

SamplingStrategyFetcher

¶

type SamplingStrategyFetcher interface {

Fetch(service

string

) ([]

byte

,

error

)

}

SamplingStrategyFetcher is used to fetch sampling strategy updates from remote server.

type

SamplingStrategyParser

¶

type SamplingStrategyParser interface {

Parse(response []

byte

) (interface{},

error

)

}

SamplingStrategyParser is used to parse sampling strategy updates. The output object
should be of the type that is recognized by the SamplerUpdaters.

type

Span

¶

type Span struct {

sync

.

RWMutex

// contains filtered or unexported fields

}

Span implements opentracing.Span

func (*Span)

BaggageItem

¶

func (s *

Span

) BaggageItem(key

string

)

string

BaggageItem implements BaggageItem() of opentracing.SpanContext

func (*Span)

Context

¶

func (s *

Span

) Context()

opentracing

.

SpanContext

Context implements opentracing.Span API

func (*Span)

Duration

¶

func (s *

Span

) Duration()

time

.

Duration

Duration returns span duration

func (*Span)

Finish

¶

func (s *

Span

) Finish()

Finish implements opentracing.Span API
After finishing the Span object it returns back to the allocator unless the reporter retains it again,
so after that, the Span object should no longer be used because it won't be valid anymore.

func (*Span)

FinishWithOptions

¶

func (s *

Span

) FinishWithOptions(options

opentracing

.

FinishOptions

)

FinishWithOptions implements opentracing.Span API

func (*Span)

Log

¶

func (s *

Span

) Log(ld

opentracing

.

LogData

)

Log implements opentracing.Span API

func (*Span)

LogEvent

¶

func (s *

Span

) LogEvent(event

string

)

LogEvent implements opentracing.Span API

func (*Span)

LogEventWithPayload

¶

func (s *

Span

) LogEventWithPayload(event

string

, payload interface{})

LogEventWithPayload implements opentracing.Span API

func (*Span)

LogFields

¶

func (s *

Span

) LogFields(fields ...

log

.

Field

)

LogFields implements opentracing.Span API

func (*Span)

LogKV

¶

func (s *

Span

) LogKV(alternatingKeyValues ...interface{})

LogKV implements opentracing.Span API

func (*Span)

Logs

¶

func (s *

Span

) Logs() []

opentracing

.

LogRecord

Logs returns micro logs for span

func (*Span)

OperationName

¶

func (s *

Span

) OperationName()

string

OperationName allows retrieving current operation name.

func (*Span)

References

¶

func (s *

Span

) References() []

opentracing

.

SpanReference

References returns references for this span

func (*Span)

Release

¶

func (s *

Span

) Release()

Release decrements object counter and return to the
allocator manager  when counter will below zero

func (*Span)

Retain

¶

func (s *

Span

) Retain() *

Span

Retain increases object counter to increase the lifetime of the object

func (*Span)

SetBaggageItem

¶

func (s *

Span

) SetBaggageItem(key, value

string

)

opentracing

.

Span

SetBaggageItem implements SetBaggageItem() of opentracing.SpanContext.
The call is proxied via tracer.baggageSetter to allow policies to be applied
before allowing to set/replace baggage keys.
The setter eventually stores a new SpanContext with extended baggage:

span.context = span.context.WithBaggageItem(key, value)

See SpanContext.WithBaggageItem() for explanation why it's done this way.

func (*Span)

SetOperationName

¶

func (s *

Span

) SetOperationName(operationName

string

)

opentracing

.

Span

SetOperationName sets or changes the operation name.

func (*Span)

SetTag

¶

func (s *

Span

) SetTag(key

string

, value interface{})

opentracing

.

Span

SetTag implements SetTag() of opentracing.Span

func (*Span)

SpanContext

¶

func (s *

Span

) SpanContext()

SpanContext

SpanContext returns span context

func (*Span)

StartTime

¶

func (s *

Span

) StartTime()

time

.

Time

StartTime returns span start time

func (*Span)

String

¶

func (s *

Span

) String()

string

func (*Span)

Tags

¶

func (s *

Span

) Tags()

opentracing

.

Tags

Tags returns tags for span

func (*Span)

Tracer

¶

func (s *

Span

) Tracer()

opentracing

.

Tracer

Tracer implements opentracing.Span API

type

SpanAllocator

¶

type SpanAllocator interface {

Get() *

Span

Put(*

Span

)

}

SpanAllocator abstraction of managing span allocations

type

SpanContext

¶

type SpanContext struct {

// contains filtered or unexported fields

}

SpanContext represents propagated span identity and state

func

ContextFromString

¶

func ContextFromString(value

string

) (

SpanContext

,

error

)

ContextFromString reconstructs the Context encoded in a string

func

NewSpanContext

¶

func NewSpanContext(traceID

TraceID

, spanID, parentID

SpanID

, sampled

bool

, baggage map[

string

]

string

)

SpanContext

NewSpanContext creates a new instance of SpanContext

func (*SpanContext)

CopyFrom

¶

func (c *

SpanContext

) CopyFrom(ctx *

SpanContext

)

CopyFrom copies data from ctx into this context, including span identity and baggage.
TODO This is only used by interop.go. Remove once TChannel Go supports OpenTracing.

func (SpanContext)

ExtendedSamplingState

¶

func (c

SpanContext

) ExtendedSamplingState(key interface{}, initValue func() interface{}) interface{}

ExtendedSamplingState returns the custom state object for a given key. If the value for this key does not exist,
it is initialized via initValue function. This state can be used by samplers (e.g. x.PrioritySampler).

func (SpanContext)

Flags

¶

func (c

SpanContext

) Flags()

byte

Flags returns the bitmap containing such bits as 'sampled' and 'debug'.

func (SpanContext)

ForeachBaggageItem

¶

func (c

SpanContext

) ForeachBaggageItem(handler func(k, v

string

)

bool

)

ForeachBaggageItem implements ForeachBaggageItem() of opentracing.SpanContext

func (SpanContext)

IsDebug

¶

func (c

SpanContext

) IsDebug()

bool

IsDebug indicates whether sampling was explicitly requested by the service.

func (SpanContext)

IsFirehose

¶

func (c

SpanContext

) IsFirehose()

bool

IsFirehose indicates whether the firehose flag was set

func (SpanContext)

IsSampled

¶

func (c

SpanContext

) IsSampled()

bool

IsSampled returns whether this trace was chosen for permanent storage
by the sampling mechanism of the tracer.

func (SpanContext)

IsSamplingFinalized

¶

func (c

SpanContext

) IsSamplingFinalized()

bool

IsSamplingFinalized indicates whether the sampling decision has been finalized.

func (SpanContext)

IsValid

¶

added in

v1.4.1

func (c

SpanContext

) IsValid()

bool

IsValid indicates whether this context actually represents a valid trace.

func (SpanContext)

ParentID

¶

func (c

SpanContext

) ParentID()

SpanID

ParentID returns the parent span ID of this span context

func (SpanContext)

SetFirehose

¶

func (c

SpanContext

) SetFirehose()

SetFirehose enables firehose mode for this trace.

func (SpanContext)

SpanID

¶

func (c

SpanContext

) SpanID()

SpanID

SpanID returns the span ID of this span context

func (SpanContext)

String

¶

func (c

SpanContext

) String()

string

func (SpanContext)

TraceID

¶

func (c

SpanContext

) TraceID()

TraceID

TraceID returns the trace ID of this span context

func (SpanContext)

WithBaggageItem

¶

func (c

SpanContext

) WithBaggageItem(key, value

string

)

SpanContext

WithBaggageItem creates a new context with an extra baggage item.
Delete a baggage item if provided blank value.

The SpanContext is designed to be immutable and passed by value. As such,
it cannot contain any locks, and should only hold immutable data, including baggage.
Another reason for why baggage is immutable is when the span context is passed
as a parent when starting a new span. The new span's baggage cannot affect the parent
span's baggage, so the child span either needs to take a copy of the parent baggage
(which is expensive and unnecessary since baggage rarely changes in the life span of
a trace), or it needs to do a copy-on-write, which is the approach taken here.

type

SpanID

¶

type SpanID

uint64

SpanID represents unique 64bit identifier of a span

func

SpanIDFromString

¶

func SpanIDFromString(s

string

) (

SpanID

,

error

)

SpanIDFromString creates a SpanID from a hexadecimal string

func (SpanID)

String

¶

func (s

SpanID

) String()

string

type

SpanObserver

deprecated

type SpanObserver interface {

OnSetOperationName(operationName

string

)

OnSetTag(key

string

, value interface{})

OnFinish(options

opentracing

.

FinishOptions

)

}

SpanObserver is created by the Observer and receives notifications about
other Span events.

Deprecated: use jaeger.ContribSpanObserver instead.

type

Tag

¶

type Tag struct {

// contains filtered or unexported fields

}

Tag is a simple key value wrapper.
TODO (breaking change) deprecate in the next major release, use opentracing.Tag instead.

func

NewTag

¶

func NewTag(key

string

, value interface{})

Tag

NewTag creates a new Tag.
TODO (breaking change) deprecate in the next major release, use opentracing.Tag instead.

type

TextMapPropagator

¶

type TextMapPropagator struct {

// contains filtered or unexported fields

}

TextMapPropagator is a combined Injector and Extractor for TextMap format

func

NewHTTPHeaderPropagator

¶

func NewHTTPHeaderPropagator(headerKeys *

HeadersConfig

, metrics

Metrics

) *

TextMapPropagator

NewHTTPHeaderPropagator creates a combined Injector and Extractor for HTTPHeaders format

func

NewTextMapPropagator

¶

func NewTextMapPropagator(headerKeys *

HeadersConfig

, metrics

Metrics

) *

TextMapPropagator

NewTextMapPropagator creates a combined Injector and Extractor for TextMap format

func (*TextMapPropagator)

Extract

¶

func (p *

TextMapPropagator

) Extract(abstractCarrier interface{}) (

SpanContext

,

error

)

Extract implements Extractor of TextMapPropagator

func (*TextMapPropagator)

Inject

¶

func (p *

TextMapPropagator

) Inject(
	sc

SpanContext

,
	abstractCarrier interface{},
)

error

Inject implements Injector of TextMapPropagator

type

TraceID

¶

type TraceID struct {

High, Low

uint64

}

TraceID represents unique 128bit identifier of a trace

func

TraceIDFromString

¶

func TraceIDFromString(s

string

) (

TraceID

,

error

)

TraceIDFromString creates a TraceID from a hexadecimal string

func (TraceID)

IsValid

¶

func (t

TraceID

) IsValid()

bool

IsValid checks if the trace ID is valid, i.e. not zero.

func (TraceID)

String

¶

func (t

TraceID

) String()

string

type

Tracer

¶

type Tracer struct {

// contains filtered or unexported fields

}

Tracer implements opentracing.Tracer.

func (*Tracer)

Close

¶

func (t *

Tracer

) Close()

error

Close releases all resources used by the Tracer and flushes any remaining buffered spans.

func (*Tracer)

Extract

¶

func (t *

Tracer

) Extract(
	format interface{},
	carrier interface{},
) (

opentracing

.

SpanContext

,

error

)

Extract implements Extract() method of opentracing.Tracer

func (*Tracer)

Inject

¶

func (t *

Tracer

) Inject(ctx

opentracing

.

SpanContext

, format interface{}, carrier interface{})

error

Inject implements Inject() method of opentracing.Tracer

func (*Tracer)

Sampler

¶

func (t *

Tracer

) Sampler()

SamplerV2

Sampler returns the sampler given to the tracer at creation.

func (*Tracer)

StartSpan

¶

func (t *

Tracer

) StartSpan(
	operationName

string

,
	options ...

opentracing

.

StartSpanOption

,
)

opentracing

.

Span

StartSpan implements StartSpan() method of opentracing.Tracer.

func (*Tracer)

Tags

¶

func (t *

Tracer

) Tags() []

opentracing

.

Tag

Tags returns a slice of tracer-level tags.

type

TracerOption

¶

type TracerOption func(tracer *

Tracer

)

TracerOption is a function that sets some option on the tracer

type

TracerOptionsFactory

¶

type TracerOptionsFactory struct{}

TracerOptionsFactory is a struct that defines functions for all available TracerOption's.

var TracerOptions

TracerOptionsFactory

TracerOptions is a factory for all available TracerOption's.

func (TracerOptionsFactory)

BaggageRestrictionManager

¶

func (

TracerOptionsFactory

) BaggageRestrictionManager(mgr

baggage

.

RestrictionManager

)

TracerOption

BaggageRestrictionManager registers BaggageRestrictionManager.

func (TracerOptionsFactory)

ContribObserver

¶

func (

TracerOptionsFactory

) ContribObserver(observer

ContribObserver

)

TracerOption

ContribObserver registers a ContribObserver.

func (TracerOptionsFactory)

CustomHeaderKeys

¶

func (

TracerOptionsFactory

) CustomHeaderKeys(headerKeys *

HeadersConfig

)

TracerOption

CustomHeaderKeys allows to override default HTTP header keys used to propagate
tracing context.

func (TracerOptionsFactory)

DebugThrottler

¶

func (

TracerOptionsFactory

) DebugThrottler(throttler

throttler

.

Throttler

)

TracerOption

DebugThrottler registers a Throttler for debug spans.

func (TracerOptionsFactory)

Extractor

¶

func (

TracerOptionsFactory

) Extractor(format interface{}, extractor

Extractor

)

TracerOption

Extractor registers an Extractor for given format.

func (TracerOptionsFactory)

Gen128Bit

¶

func (

TracerOptionsFactory

) Gen128Bit(gen128Bit

bool

)

TracerOption

Gen128Bit enables generation of 128bit trace IDs.

func (TracerOptionsFactory)

HighTraceIDGenerator

¶

func (

TracerOptionsFactory

) HighTraceIDGenerator(highTraceIDGenerator func()

uint64

)

TracerOption

HighTraceIDGenerator allows to override define ID generator.

func (TracerOptionsFactory)

HostIPv4

¶

func (

TracerOptionsFactory

) HostIPv4(hostIPv4

uint32

)

TracerOption

HostIPv4 creates a TracerOption that identifies the current service/process.
If not set, the factory method will obtain the current IP address.
The TracerOption is deprecated; the tracer will attempt to automatically detect the IP.

Deprecated.

func (TracerOptionsFactory)

Injector

¶

func (

TracerOptionsFactory

) Injector(format interface{}, injector

Injector

)

TracerOption

Injector registers a Injector for given format.

func (TracerOptionsFactory)

Logger

¶

func (

TracerOptionsFactory

) Logger(logger

Logger

)

TracerOption

Logger creates a TracerOption that gives the tracer a Logger.

func (TracerOptionsFactory)

MaxLogsPerSpan

¶

func (

TracerOptionsFactory

) MaxLogsPerSpan(maxLogsPerSpan

int

)

TracerOption

MaxLogsPerSpan limits the number of Logs in a span (if set to a nonzero
value). If a span has more logs than this value, logs are dropped as
necessary (and replaced with a log describing how many were dropped).

About half of the MaxLogsPerSpan logs kept are the oldest logs, and about
half are the newest logs.

func (TracerOptionsFactory)

MaxTagValueLength

¶

func (

TracerOptionsFactory

) MaxTagValueLength(maxTagValueLength

int

)

TracerOption

MaxTagValueLength sets the limit on the max length of tag values.

func (TracerOptionsFactory)

Metrics

¶

func (

TracerOptionsFactory

) Metrics(m *

Metrics

)

TracerOption

Metrics creates a TracerOption that initializes Metrics on the tracer,
which is used to emit statistics.

func (TracerOptionsFactory)

NoDebugFlagOnForcedSampling

¶

func (

TracerOptionsFactory

) NoDebugFlagOnForcedSampling(noDebugFlagOnForcedSampling

bool

)

TracerOption

NoDebugFlagOnForcedSampling turns off setting the debug flag in the trace context
when the trace is force-started via sampling=1 span tag.

func (TracerOptionsFactory)

Observer

¶

func (t

TracerOptionsFactory

) Observer(observer

Observer

)

TracerOption

Observer registers an Observer.

func (TracerOptionsFactory)

PoolSpans

¶

func (

TracerOptionsFactory

) PoolSpans(poolSpans

bool

)

TracerOption

PoolSpans creates a TracerOption that tells the tracer whether it should use
an object pool to minimize span allocations.
This should be used with care, only if the service is not running any async tasks
that can access parent spans after those spans have been finished.

func (TracerOptionsFactory)

RandomNumber

¶

func (

TracerOptionsFactory

) RandomNumber(randomNumber func()

uint64

)

TracerOption

RandomNumber creates a TracerOption that gives the tracer
a thread-safe random number generator function for generating trace IDs.

func (TracerOptionsFactory)

Tag

¶

func (

TracerOptionsFactory

) Tag(key

string

, value interface{})

TracerOption

Tag adds a tracer-level tag that will be added to all spans.

func (TracerOptionsFactory)

TimeNow

¶

func (

TracerOptionsFactory

) TimeNow(timeNow func()

time

.

Time

)

TracerOption

TimeNow creates a TracerOption that gives the tracer a function
used to generate timestamps for spans.

func (TracerOptionsFactory)

ZipkinSharedRPCSpan

¶

func (

TracerOptionsFactory

) ZipkinSharedRPCSpan(zipkinSharedRPCSpan

bool

)

TracerOption

ZipkinSharedRPCSpan enables a mode where server-side span shares the span ID
from the client span from the incoming request, for compatibility with Zipkin's
"one span per RPC" model.

type

Transport

¶

type Transport interface {

// Append converts the span to the wire representation and adds it

// to sender's internal buffer.  If the buffer exceeds its designated

// size, the transport should call Flush() and return the number of spans

// flushed, otherwise return 0. If error is returned, the returned number

// of spans is treated as failed span, and reported to metrics accordingly.

Append(span *

Span

) (

int

,

error

)

// Flush submits the internal buffer to the remote server. It returns the

// number of spans flushed. If error is returned, the returned number of

// spans is treated as failed span, and reported to metrics accordingly.

Flush() (

int

,

error

)

io

.

Closer

}

Transport abstracts the method of sending spans out of process.
Implementations are NOT required to be thread-safe; the RemoteReporter
is expected to only call methods on the Transport from the same go-routine.

func

NewUDPTransport

¶

func NewUDPTransport(hostPort

string

, maxPacketSize

int

) (

Transport

,

error

)

NewUDPTransport creates a reporter that submits spans to jaeger-agent.
TODO: (breaking change) move to transport/ package.

func

NewUDPTransportWithParams

¶

func NewUDPTransportWithParams(params

UDPTransportParams

) (

Transport

,

error

)

NewUDPTransportWithParams creates a reporter that submits spans to jaeger-agent.
TODO: (breaking change) move to transport/ package.

type

UDPTransportParams

¶

type UDPTransportParams struct {

utils

.

AgentClientUDPParams

}

UDPTransportParams allows specifying options for initializing a UDPTransport. An instance of this struct should
be passed to NewUDPTransportWithParams.