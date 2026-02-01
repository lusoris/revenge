# Jaeger Go Client

> Source: https://pkg.go.dev/github.com/jaegertracing/jaeger-client-go
> Fetched: 2026-02-01T11:50:10.008742+00:00
> Content-Hash: 5768f484eeb1cf71
> Type: html

---

### Overview ¶

Package jaeger implements an OpenTracing (<http://opentracing.io>) Tracer.

For integration instructions please refer to the README:

<https://github.com/uber/jaeger-client-go/blob/master/README.md>

### Index ¶

- Constants
- Variables
- func BuildJaegerProcessThrift(span *Span)*j.Process
- func BuildJaegerThrift(span *Span)*j.Span
- func BuildZipkinThrift(s *Span)*z.Span
- func ConvertLogsToJaegerTags(logFields []log.Field) []*j.Tag
- func EnableFirehose(s *Span)
- func NewTracer(serviceName string, sampler Sampler, reporter Reporter, ...) (opentracing.Tracer, io.Closer)
- func SelfRef(ctx SpanContext) opentracing.SpanReference
- type AdaptiveSamplerUpdater
-     * func (u *AdaptiveSamplerUpdater) Update(sampler SamplerV2, strategy interface{}) (SamplerV2, error)
- type BinaryPropagator
-     * func NewBinaryPropagator(tracer *Tracer) *BinaryPropagator
-     * func (p *BinaryPropagator) Extract(abstractCarrier interface{}) (SpanContext, error)
  - func (p *BinaryPropagator) Inject(sc SpanContext, abstractCarrier interface{}) error
- type ConstSampler
-     * func NewConstSampler(sample bool) *ConstSampler
-     * func (s *ConstSampler) Close()
  - func (s *ConstSampler) Equal(other Sampler) bool
  - func (s *ConstSampler) IsSampled(id TraceID, operation string) (bool, []Tag)
  - func (s *ConstSampler) OnCreateSpan(span*Span) SamplingDecision
  - func (s *ConstSampler) OnFinishSpan(span*Span) SamplingDecision
  - func (s *ConstSampler) OnSetOperationName(span*Span, operationName string) SamplingDecision
  - func (s *ConstSampler) OnSetTag(span*Span, key string, value interface{}) SamplingDecision
  - func (s *ConstSampler) String() string
- type ContribObserver
- type ContribSpanObserver
- type ExtractableZipkinSpan
- type Extractor
- type GuaranteedThroughputProbabilisticSampler
-     * func NewGuaranteedThroughputProbabilisticSampler(lowerBound, samplingRate float64) (*GuaranteedThroughputProbabilisticSampler, error)
-     * func (s *GuaranteedThroughputProbabilisticSampler) Close()
  - func (s *GuaranteedThroughputProbabilisticSampler) Equal(other Sampler) bool
  - func (s *GuaranteedThroughputProbabilisticSampler) IsSampled(id TraceID, operation string) (bool, []Tag)
  - func (s GuaranteedThroughputProbabilisticSampler) String() string
- type HeadersConfig
-     * func (c *HeadersConfig) ApplyDefaults() *HeadersConfig
- type InMemoryReporter
-     * func NewInMemoryReporter() *InMemoryReporter
-     * func (r *InMemoryReporter) Close()
  - func (r *InMemoryReporter) GetSpans() []opentracing.Span
  - func (r *InMemoryReporter) Report(span*Span)
  - func (r *InMemoryReporter) Reset()
  - func (r *InMemoryReporter) SpansSubmitted() int
- type InjectableZipkinSpan
- type Injector
- type Logger
- type Metrics
-     * func NewMetrics(factory metrics.Factory, globalTags map[string]string) *Metrics
  - func NewNullMetrics() *Metrics
- type Observerdeprecated
- type PerOperationSampler
-     * func NewAdaptiveSampler(strategies *sampling.PerOperationSamplingStrategies, maxOperations int) (*PerOperationSampler, error)
  - func NewPerOperationSampler(params PerOperationSamplerParams) *PerOperationSampler
-     * func (s *PerOperationSampler) Close()
  - func (s *PerOperationSampler) Equal(other Sampler) bool
  - func (s *PerOperationSampler) IsSampled(id TraceID, operation string) (bool, []Tag)
  - func (s *PerOperationSampler) OnCreateSpan(span*Span) SamplingDecision
  - func (s *PerOperationSampler) OnFinishSpan(span*Span) SamplingDecision
  - func (s *PerOperationSampler) OnSetOperationName(span*Span, operationName string) SamplingDecision
  - func (s *PerOperationSampler) OnSetTag(span*Span, key string, value interface{}) SamplingDecision
  - func (s *PerOperationSampler) String() string
- type PerOperationSamplerParams
- type ProbabilisticSampler
-     * func NewProbabilisticSampler(samplingRate float64) (*ProbabilisticSampler, error)
-     * func (s *ProbabilisticSampler) Close()
  - func (s *ProbabilisticSampler) Equal(other Sampler) bool
  - func (s *ProbabilisticSampler) IsSampled(id TraceID, operation string) (bool, []Tag)
  - func (s *ProbabilisticSampler) OnCreateSpan(span*Span) SamplingDecision
  - func (s *ProbabilisticSampler) OnFinishSpan(span*Span) SamplingDecision
  - func (s *ProbabilisticSampler) OnSetOperationName(span*Span, operationName string) SamplingDecision
  - func (s *ProbabilisticSampler) OnSetTag(span*Span, key string, value interface{}) SamplingDecision
  - func (s *ProbabilisticSampler) SamplingRate() float64
  - func (s *ProbabilisticSampler) String() string
  - func (s *ProbabilisticSampler) Update(samplingRate float64) error
- type ProbabilisticSamplerUpdater
-     * func (u *ProbabilisticSamplerUpdater) Update(sampler SamplerV2, strategy interface{}) (SamplerV2, error)
- type Process
- type ProcessSetter
- type RateLimitingSampler
-     * func NewRateLimitingSampler(maxTracesPerSecond float64) *RateLimitingSampler
-     * func (s *RateLimitingSampler) Close()
  - func (s *RateLimitingSampler) Equal(other Sampler) bool
  - func (s *RateLimitingSampler) IsSampled(id TraceID, operation string) (bool, []Tag)
  - func (s *RateLimitingSampler) OnCreateSpan(span*Span) SamplingDecision
  - func (s *RateLimitingSampler) OnFinishSpan(span*Span) SamplingDecision
  - func (s *RateLimitingSampler) OnSetOperationName(span*Span, operationName string) SamplingDecision
  - func (s *RateLimitingSampler) OnSetTag(span*Span, key string, value interface{}) SamplingDecision
  - func (s *RateLimitingSampler) String() string
  - func (s *RateLimitingSampler) Update(maxTracesPerSecond float64)
- type RateLimitingSamplerUpdater
-     * func (u *RateLimitingSamplerUpdater) Update(sampler SamplerV2, strategy interface{}) (SamplerV2, error)
- type Reference
- type RemotelyControlledSampler
-     * func NewRemotelyControlledSampler(serviceName string, opts ...SamplerOption) *RemotelyControlledSampler
-     * func (s *RemotelyControlledSampler) Close()
  - func (s *RemotelyControlledSampler) Equal(other Sampler) bool
  - func (s *RemotelyControlledSampler) IsSampled(id TraceID, operation string) (bool, []Tag)
  - func (s *RemotelyControlledSampler) OnCreateSpan(span*Span) SamplingDecision
  - func (s *RemotelyControlledSampler) OnFinishSpan(span*Span) SamplingDecision
  - func (s *RemotelyControlledSampler) OnSetOperationName(span*Span, operationName string) SamplingDecision
  - func (s *RemotelyControlledSampler) OnSetTag(span*Span, key string, value interface{}) SamplingDecision
  - func (s *RemotelyControlledSampler) Sampler() SamplerV2
  - func (s *RemotelyControlledSampler) UpdateSampler()
- type Reporter
-     * func NewCompositeReporter(reporters ...Reporter) Reporter
  - func NewLoggingReporter(logger Logger) Reporter
  - func NewNullReporter() Reporter
  - func NewRemoteReporter(sender Transport, opts ...ReporterOption) Reporter
- type ReporterOption
- type Sampler
- type SamplerOption
- type SamplerOptionsFactory
-     * func (SamplerOptionsFactory) InitialSampler(sampler Sampler) SamplerOption
  - func (SamplerOptionsFactory) Logger(logger Logger) SamplerOption
  - func (SamplerOptionsFactory) MaxOperations(maxOperations int) SamplerOption
  - func (SamplerOptionsFactory) Metrics(m *Metrics) SamplerOption
  - func (SamplerOptionsFactory) OperationNameLateBinding(enable bool) SamplerOption
  - func (SamplerOptionsFactory) SamplingRefreshInterval(samplingRefreshInterval time.Duration) SamplerOption
  - func (SamplerOptionsFactory) SamplingServerURL(samplingServerURL string) SamplerOption
  - func (SamplerOptionsFactory) SamplingStrategyFetcher(fetcher SamplingStrategyFetcher) SamplerOption
  - func (SamplerOptionsFactory) SamplingStrategyParser(parser SamplingStrategyParser) SamplerOption
  - func (SamplerOptionsFactory) Updaters(updaters ...SamplerUpdater) SamplerOption
- type SamplerUpdater
- type SamplerV2
- type SamplerV2Base
-     * func (SamplerV2Base) Close()
  - func (SamplerV2Base) Equal(other Sampler) bool
  - func (SamplerV2Base) IsSampled(id TraceID, operation string) (sampled bool, tags []Tag)
- type SamplingDecision
- type SamplingStrategyFetcher
- type SamplingStrategyParser
- type Span
-     * func (s *Span) BaggageItem(key string) string
  - func (s *Span) Context() opentracing.SpanContext
  - func (s *Span) Duration() time.Duration
  - func (s *Span) Finish()
  - func (s *Span) FinishWithOptions(options opentracing.FinishOptions)
  - func (s *Span) Log(ld opentracing.LogData)
  - func (s *Span) LogEvent(event string)
  - func (s *Span) LogEventWithPayload(event string, payload interface{})
  - func (s *Span) LogFields(fields ...log.Field)
  - func (s *Span) LogKV(alternatingKeyValues ...interface{})
  - func (s *Span) Logs() []opentracing.LogRecord
  - func (s *Span) OperationName() string
  - func (s *Span) References() []opentracing.SpanReference
  - func (s *Span) Release()
  - func (s *Span) Retain()*Span
  - func (s *Span) SetBaggageItem(key, value string) opentracing.Span
  - func (s *Span) SetOperationName(operationName string) opentracing.Span
  - func (s *Span) SetTag(key string, value interface{}) opentracing.Span
  - func (s *Span) SpanContext() SpanContext
  - func (s *Span) StartTime() time.Time
  - func (s *Span) String() string
  - func (s *Span) Tags() opentracing.Tags
  - func (s *Span) Tracer() opentracing.Tracer
- type SpanAllocator
- type SpanContext
-     * func ContextFromString(value string) (SpanContext, error)
  - func NewSpanContext(traceID TraceID, spanID, parentID SpanID, sampled bool, ...) SpanContext
-     * func (c *SpanContext) CopyFrom(ctx *SpanContext)
  - func (c SpanContext) ExtendedSamplingState(key interface{}, initValue func() interface{}) interface{}
  - func (c SpanContext) Flags() byte
  - func (c SpanContext) ForeachBaggageItem(handler func(k, v string) bool)
  - func (c SpanContext) IsDebug() bool
  - func (c SpanContext) IsFirehose() bool
  - func (c SpanContext) IsSampled() bool
  - func (c SpanContext) IsSamplingFinalized() bool
  - func (c SpanContext) IsValid() bool
  - func (c SpanContext) ParentID() SpanID
  - func (c SpanContext) SetFirehose()
  - func (c SpanContext) SpanID() SpanID
  - func (c SpanContext) String() string
  - func (c SpanContext) TraceID() TraceID
  - func (c SpanContext) WithBaggageItem(key, value string) SpanContext
- type SpanID
-     * func SpanIDFromString(s string) (SpanID, error)
-     * func (s SpanID) String() string
- type SpanObserverdeprecated
- type Tag
-     * func NewTag(key string, value interface{}) Tag
- type TextMapPropagator
-     * func NewHTTPHeaderPropagator(headerKeys *HeadersConfig, metrics Metrics) *TextMapPropagator
  - func NewTextMapPropagator(headerKeys *HeadersConfig, metrics Metrics)*TextMapPropagator
-     * func (p *TextMapPropagator) Extract(abstractCarrier interface{}) (SpanContext, error)
  - func (p *TextMapPropagator) Inject(sc SpanContext, abstractCarrier interface{}) error
- type TraceID
-     * func TraceIDFromString(s string) (TraceID, error)
-     * func (t TraceID) IsValid() bool
  - func (t TraceID) String() string
- type Tracer
-     * func (t *Tracer) Close() error
  - func (t *Tracer) Extract(format interface{}, carrier interface{}) (opentracing.SpanContext, error)
  - func (t *Tracer) Inject(ctx opentracing.SpanContext, format interface{}, carrier interface{}) error
  - func (t *Tracer) Sampler() SamplerV2
  - func (t *Tracer) StartSpan(operationName string, options ...opentracing.StartSpanOption) opentracing.Span
  - func (t *Tracer) Tags() []opentracing.Tag
- type TracerOption
- type TracerOptionsFactory
-     * func (TracerOptionsFactory) BaggageRestrictionManager(mgr baggage.RestrictionManager) TracerOption
  - func (TracerOptionsFactory) ContribObserver(observer ContribObserver) TracerOption
  - func (TracerOptionsFactory) CustomHeaderKeys(headerKeys *HeadersConfig) TracerOption
  - func (TracerOptionsFactory) DebugThrottler(throttler throttler.Throttler) TracerOption
  - func (TracerOptionsFactory) Extractor(format interface{}, extractor Extractor) TracerOption
  - func (TracerOptionsFactory) Gen128Bit(gen128Bit bool) TracerOption
  - func (TracerOptionsFactory) HighTraceIDGenerator(highTraceIDGenerator func() uint64) TracerOption
  - func (TracerOptionsFactory) HostIPv4(hostIPv4 uint32) TracerOption
  - func (TracerOptionsFactory) Injector(format interface{}, injector Injector) TracerOption
  - func (TracerOptionsFactory) Logger(logger Logger) TracerOption
  - func (TracerOptionsFactory) MaxLogsPerSpan(maxLogsPerSpan int) TracerOption
  - func (TracerOptionsFactory) MaxTagValueLength(maxTagValueLength int) TracerOption
  - func (TracerOptionsFactory) Metrics(m *Metrics) TracerOption
  - func (TracerOptionsFactory) NoDebugFlagOnForcedSampling(noDebugFlagOnForcedSampling bool) TracerOption
  - func (t TracerOptionsFactory) Observer(observer Observer) TracerOption
  - func (TracerOptionsFactory) PoolSpans(poolSpans bool) TracerOption
  - func (TracerOptionsFactory) RandomNumber(randomNumber func() uint64) TracerOption
  - func (TracerOptionsFactory) Tag(key string, value interface{}) TracerOption
  - func (TracerOptionsFactory) TimeNow(timeNow func() time.Time) TracerOption
  - func (TracerOptionsFactory) ZipkinSharedRPCSpan(zipkinSharedRPCSpan bool) TracerOption
- type Transport
-     * func NewUDPTransport(hostPort string, maxPacketSize int) (Transport, error)
  - func NewUDPTransportWithParams(params UDPTransportParams) (Transport, error)
- type UDPTransportParams

### Constants ¶

[View Source](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/constants.go#L23)

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
     TracerStateHeaderName = TraceContextHeaderName
    
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

[View Source](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/interop.go#L28)

    const SpanContextFormat formatKey = [iota](/builtin#iota)

SpanContextFormat is a constant used as OpenTracing Format. Requires *SpanContext as carrier. This format is intended for interop with TChannel or other Zipkin-like tracers.

[View Source](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/zipkin.go#L22)

    const ZipkinSpanFormat = "zipkin-span-format"

ZipkinSpanFormat is an OpenTracing carrier format constant

### Variables ¶

[View Source](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/constants.go#L103)

    var (
     // DefaultSamplingServerURL is the default url to fetch sampling config from, via http
     DefaultSamplingServerURL = [fmt](/fmt).[Sprintf](/fmt#Sprintf)("http://127.0.0.1:%d/sampling", DefaultSamplingServerPort)
    )

[View Source](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/logger.go#L48)

    var NullLogger = &nullLogger{}

NullLogger is implementation of the Logger interface that delegates to default `log` package

[View Source](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/reporter_options.go#L27)

    var ReporterOptions reporterOptions

ReporterOptions is a factory for all available ReporterOption's

[View Source](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/logger.go#L34)

    var StdLogger = &stdLogger{}

StdLogger is implementation of the Logger interface that delegates to default `log` package

### Functions ¶

#### func [BuildJaegerProcessThrift](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/jaeger_thrift_span.go#L51) ¶

    func BuildJaegerProcessThrift(span *Span) *[j](/github.com/uber/jaeger-client-go/thrift-gen/jaeger).[Process](/github.com/uber/jaeger-client-go/thrift-gen/jaeger#Process)

BuildJaegerProcessThrift creates a thrift Process type. TODO: (breaking change) move to internal package.

#### func [BuildJaegerThrift](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/jaeger_thrift_span.go#L28) ¶

    func BuildJaegerThrift(span *Span) *[j](/github.com/uber/jaeger-client-go/thrift-gen/jaeger).[Span](/github.com/uber/jaeger-client-go/thrift-gen/jaeger#Span)

BuildJaegerThrift builds jaeger span based on internal span. TODO: (breaking change) move to internal package.

#### func [BuildZipkinThrift](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/zipkin_thrift_span.go#L44) ¶

    func BuildZipkinThrift(s *Span) *[z](/github.com/uber/jaeger-client-go/thrift-gen/zipkincore).[Span](/github.com/uber/jaeger-client-go/thrift-gen/zipkincore#Span)

BuildZipkinThrift builds thrift span based on internal span. TODO: (breaking change) move to transport/zipkin and make private.

#### func [ConvertLogsToJaegerTags](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/jaeger_tag.go#L28) ¶

    func ConvertLogsToJaegerTags(logFields [][log](/github.com/opentracing/opentracing-go/log).[Field](/github.com/opentracing/opentracing-go/log#Field)) []*[j](/github.com/uber/jaeger-client-go/thrift-gen/jaeger).[Tag](/github.com/uber/jaeger-client-go/thrift-gen/jaeger#Tag)

ConvertLogsToJaegerTags converts log Fields into jaeger tags.

#### func [EnableFirehose](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/span.go#L499) ¶

    func EnableFirehose(s *Span)

EnableFirehose enables firehose flag on the span context

#### func [NewTracer](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/tracer.go#L79) ¶

    func NewTracer(
     serviceName [string](/builtin#string),
     sampler Sampler,
     reporter Reporter,
     options ...TracerOption,
    ) ([opentracing](/github.com/opentracing/opentracing-go).[Tracer](/github.com/opentracing/opentracing-go#Tracer), [io](/io).[Closer](/io#Closer))

NewTracer creates Tracer implementation that reports tracing to Jaeger. The returned io.Closer can be used in shutdown hooks to ensure that the internal queue of the Reporter is drained and all buffered spans are submitted to collectors. TODO (breaking change) return *Tracer only, without closer.

#### func [SelfRef](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/tracer.go#L488) ¶

    func SelfRef(ctx SpanContext) [opentracing](/github.com/opentracing/opentracing-go).[SpanReference](/github.com/opentracing/opentracing-go#SpanReference)

SelfRef creates an opentracing compliant SpanReference from a jaeger SpanContext. This is a factory function in order to encapsulate jaeger specific types.

### Types ¶

#### type [AdaptiveSamplerUpdater](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/sampler_remote.go#L272) ¶

    type AdaptiveSamplerUpdater struct {
     MaxOperations            [int](/builtin#int)
     OperationNameLateBinding [bool](/builtin#bool)
    }

AdaptiveSamplerUpdater is used by RemotelyControlledSampler to parse sampling configuration. Fields have the same meaning as in PerOperationSamplerParams.

#### func (*AdaptiveSamplerUpdater) [Update](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/sampler_remote.go#L278) ¶

    func (u *AdaptiveSamplerUpdater) Update(sampler SamplerV2, strategy interface{}) (SamplerV2, [error](/builtin#error))

Update implements Update of SamplerUpdater.

#### type [BinaryPropagator](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/propagation.go#L95) ¶

    type BinaryPropagator struct {
     // contains filtered or unexported fields
    }

BinaryPropagator is a combined Injector and Extractor for Binary format

#### func [NewBinaryPropagator](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/propagation.go#L101) ¶

    func NewBinaryPropagator(tracer *Tracer) *BinaryPropagator

NewBinaryPropagator creates a combined Injector and Extractor for Binary format

#### func (*BinaryPropagator) [Extract](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/propagation.go#L225) ¶

    func (p *BinaryPropagator) Extract(abstractCarrier interface{}) (SpanContext, [error](/builtin#error))

Extract implements Extractor of BinaryPropagator

#### func (*BinaryPropagator) [Inject](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/propagation.go#L177) ¶

    func (p *BinaryPropagator) Inject(
     sc SpanContext,
     abstractCarrier interface{},
    ) [error](/builtin#error)

Inject implements Injector of BinaryPropagator

#### type [ConstSampler](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/sampler.go#L53) ¶

    type ConstSampler struct {
     Decision [bool](/builtin#bool)
     // contains filtered or unexported fields
    }

ConstSampler is a sampler that always makes the same decision.

#### func [NewConstSampler](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/sampler.go#L60) ¶

    func NewConstSampler(sample [bool](/builtin#bool)) *ConstSampler

NewConstSampler creates a ConstSampler.

#### func (*ConstSampler) [Close](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/sampler.go#L79) ¶

    func (s *ConstSampler) Close()

Close implements Close() of Sampler.

#### func (*ConstSampler) [Equal](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/sampler.go#L84) ¶

    func (s *ConstSampler) Equal(other Sampler) [bool](/builtin#bool)

Equal implements Equal() of Sampler.

#### func (*ConstSampler) [IsSampled](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/sampler.go#L74) ¶

    func (s *ConstSampler) IsSampled(id TraceID, operation [string](/builtin#string)) ([bool](/builtin#bool), []Tag)

IsSampled implements IsSampled() of Sampler.

#### func (*ConstSampler) [OnCreateSpan](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/sampler_v2.go#L75) ¶

    func (s *ConstSampler) OnCreateSpan(span *Span) SamplingDecision

#### func (*ConstSampler) [OnFinishSpan](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/sampler_v2.go#L89) ¶

    func (s *ConstSampler) OnFinishSpan(span *Span) SamplingDecision

#### func (*ConstSampler) [OnSetOperationName](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/sampler_v2.go#L80) ¶

    func (s *ConstSampler) OnSetOperationName(span *Span, operationName [string](/builtin#string)) SamplingDecision

#### func (*ConstSampler) [OnSetTag](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/sampler_v2.go#L85) ¶

    func (s *ConstSampler) OnSetTag(span *Span, key [string](/builtin#string), value interface{}) SamplingDecision

#### func (*ConstSampler) [String](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/sampler.go#L92) ¶

    func (s *ConstSampler) String() [string](/builtin#string)

String is used to log sampler details.

#### type [ContribObserver](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/contrib_observer.go#L23) ¶

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
     OnStartSpan(sp [opentracing](/github.com/opentracing/opentracing-go).[Span](/github.com/opentracing/opentracing-go#Span), operationName [string](/builtin#string), options [opentracing](/github.com/opentracing/opentracing-go).[StartSpanOptions](/github.com/opentracing/opentracing-go#StartSpanOptions)) (ContribSpanObserver, [bool](/builtin#bool))
    }

ContribObserver can be registered with the Tracer to receive notifications about new Spans. Modelled after github.com/opentracing-contrib/go-observer.

#### type [ContribSpanObserver](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/contrib_observer.go#L42) ¶

    type ContribSpanObserver interface {
     OnSetOperationName(operationName [string](/builtin#string))
     OnSetTag(key [string](/builtin#string), value interface{})
     OnFinish(options [opentracing](/github.com/opentracing/opentracing-go).[FinishOptions](/github.com/opentracing/opentracing-go#FinishOptions))
    }

ContribSpanObserver is created by the Observer and receives notifications about other Span events. This interface is meant to match github.com/opentracing-contrib/go-observer, via duck typing, without directly importing the go-observer package.

#### type [ExtractableZipkinSpan](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/zipkin.go#L26) ¶

    type ExtractableZipkinSpan interface {
     TraceID() [uint64](/builtin#uint64)
     SpanID() [uint64](/builtin#uint64)
     ParentID() [uint64](/builtin#uint64)
     Flags() [byte](/builtin#byte)
    }

ExtractableZipkinSpan is a type of Carrier used for integration with Zipkin-aware RPC frameworks (like TChannel). It does not support baggage, only trace IDs.

#### type [Extractor](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/propagation.go#L47) ¶

    type Extractor interface {
     // Extract decodes a SpanContext instance from the given `carrier`,
     // or (nil, opentracing.ErrSpanContextNotFound) if no context could
     // be found in the `carrier`.
     Extract(carrier interface{}) (SpanContext, [error](/builtin#error))
    }

Extractor is responsible for extracting SpanContext instances from a format-specific "carrier" object. Typically the extraction will take place on the server side of an RPC boundary, but message queues and other IPC mechanisms are also reasonable places to use an Extractor.

#### type [GuaranteedThroughputProbabilisticSampler](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/sampler.go#L249) ¶

    type GuaranteedThroughputProbabilisticSampler struct {
     // contains filtered or unexported fields
    }

GuaranteedThroughputProbabilisticSampler is a sampler that leverages both ProbabilisticSampler and RateLimitingSampler. The RateLimitingSampler is used as a guaranteed lower bound sampler such that every operation is sampled at least once in a time interval defined by the lowerBound. ie a lowerBound of 1.0 / (60 * 10) will sample an operation at least once every 10 minutes.

The ProbabilisticSampler is given higher priority when tags are emitted, ie. if IsSampled() for both samplers return true, the tags for ProbabilisticSampler will be used.

#### func [NewGuaranteedThroughputProbabilisticSampler](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/sampler.go#L259) ¶

    func NewGuaranteedThroughputProbabilisticSampler(
     lowerBound, samplingRate [float64](/builtin#float64),
    ) (*GuaranteedThroughputProbabilisticSampler, [error](/builtin#error))

NewGuaranteedThroughputProbabilisticSampler returns a delegating sampler that applies both ProbabilisticSampler and RateLimitingSampler.

#### func (*GuaranteedThroughputProbabilisticSampler) [Close](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/sampler.go#L302) ¶

    func (s *GuaranteedThroughputProbabilisticSampler) Close()

Close implements Close() of Sampler.

#### func (*GuaranteedThroughputProbabilisticSampler) [Equal](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/sampler.go#L308) ¶

    func (s *GuaranteedThroughputProbabilisticSampler) Equal(other Sampler) [bool](/builtin#bool)

Equal implements Equal() of Sampler.

#### func (*GuaranteedThroughputProbabilisticSampler) [IsSampled](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/sampler.go#L292) ¶

    func (s *GuaranteedThroughputProbabilisticSampler) IsSampled(id TraceID, operation [string](/builtin#string)) ([bool](/builtin#bool), []Tag)

IsSampled implements IsSampled() of Sampler.

#### func (GuaranteedThroughputProbabilisticSampler) [String](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/sampler.go#L323) ¶

    func (s GuaranteedThroughputProbabilisticSampler) String() [string](/builtin#string)

#### type [HeadersConfig](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/header.go#L20) ¶

    type HeadersConfig struct {
     // JaegerDebugHeader is the name of HTTP header or a TextMap carrier key which,
     // if found in the carrier, forces the trace to be sampled as "debug" trace.
     // The value of the header is recorded as the tag on the root span, so that the
     // trace can be found in the UI using this value as a correlation ID.
     JaegerDebugHeader [string](/builtin#string) `yaml:"jaegerDebugHeader"`
    
     // JaegerBaggageHeader is the name of the HTTP header that is used to submit baggage.
     // It differs from TraceBaggageHeaderPrefix in that it can be used only in cases where
     // a root span does not exist.
     JaegerBaggageHeader [string](/builtin#string) `yaml:"jaegerBaggageHeader"`
    
     // TraceContextHeaderName is the http header name used to propagate tracing context.
     // This must be in lower-case to avoid mismatches when decoding incoming headers.
     TraceContextHeaderName [string](/builtin#string) `yaml:"TraceContextHeaderName"`
    
     // TraceBaggageHeaderPrefix is the prefix for http headers used to propagate baggage.
     // This must be in lower-case to avoid mismatches when decoding incoming headers.
     TraceBaggageHeaderPrefix [string](/builtin#string) `yaml:"traceBaggageHeaderPrefix"`
    }

HeadersConfig contains the values for the header keys that Jaeger will use. These values may be either custom or default depending on whether custom values were provided via a configuration.

#### func (*HeadersConfig) [ApplyDefaults](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/header.go#L42) ¶

    func (c *HeadersConfig) ApplyDefaults() *HeadersConfig

ApplyDefaults sets missing configuration keys to default values

#### type [InMemoryReporter](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/reporter.go#L83) ¶

    type InMemoryReporter struct {
     // contains filtered or unexported fields
    }

InMemoryReporter is used for testing, and simply collects spans in memory.

#### func [NewInMemoryReporter](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/reporter.go#L90) ¶

    func NewInMemoryReporter() *InMemoryReporter

NewInMemoryReporter creates a reporter that stores spans in memory. NOTE: the Tracer should be created with options.PoolSpans = false.

#### func (*InMemoryReporter) [Close](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/reporter.go#L105) ¶

    func (r *InMemoryReporter) Close()

Close implements Close() method of Reporter

#### func (*InMemoryReporter) [GetSpans](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/reporter.go#L117) ¶

    func (r *InMemoryReporter) GetSpans() [][opentracing](/github.com/opentracing/opentracing-go).[Span](/github.com/opentracing/opentracing-go#Span)

GetSpans returns accumulated spans as a copy of the buffer.

#### func (*InMemoryReporter) [Report](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/reporter.go#L97) ¶

    func (r *InMemoryReporter) Report(span *Span)

Report implements Report() method of Reporter by storing the span in the buffer.

#### func (*InMemoryReporter) [Reset](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/reporter.go#L126) ¶

    func (r *InMemoryReporter) Reset()

Reset clears all accumulated spans.

#### func (*InMemoryReporter) [SpansSubmitted](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/reporter.go#L110) ¶

    func (r *InMemoryReporter) SpansSubmitted() [int](/builtin#int)

SpansSubmitted returns the number of spans accumulated in the buffer.

#### type [InjectableZipkinSpan](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/zipkin.go#L35) ¶

    type InjectableZipkinSpan interface {
     SetTraceID(traceID [uint64](/builtin#uint64))
     SetSpanID(spanID [uint64](/builtin#uint64))
     SetParentID(parentID [uint64](/builtin#uint64))
     SetFlags(flags [byte](/builtin#byte))
    }

InjectableZipkinSpan is a type of Carrier used for integration with Zipkin-aware RPC frameworks (like TChannel). It does not support baggage, only trace IDs.

#### type [Injector](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/propagation.go#L34) ¶

    type Injector interface {
     // Inject takes `SpanContext` and injects it into `carrier`. The actual type
     // of `carrier` depends on the `format` passed to `Tracer.Inject()`.
     //
     // Implementations may return opentracing.ErrInvalidCarrier or any other
     // implementation-specific error if injection fails.
     Inject(ctx SpanContext, carrier interface{}) [error](/builtin#error)
    }

Injector is responsible for injecting SpanContext instances in a manner suitable for propagation via a format-specific "carrier" object. Typically the injection will take place across an RPC boundary, but message queues and other IPC mechanisms are also reasonable places to use an Injector.

#### type [Logger](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/logger.go#L25) ¶

    type Logger interface {
     // Error logs a message at error priority
     Error(msg [string](/builtin#string))
    
     // Infof logs a message at info priority
     Infof(msg [string](/builtin#string), args ...interface{})
    }

Logger provides an abstract interface for logging from Reporters. Applications can provide their own implementation of this interface to adapt reporters logging to whatever logging library they prefer (stdlib log, logrus, go-logging, etc).

#### type [Metrics](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/metrics.go#L22) ¶

    type Metrics struct {
     // Number of traces started by this tracer as sampled
     TracesStartedSampled [metrics](/github.com/uber/jaeger-lib/metrics).[Counter](/github.com/uber/jaeger-lib/metrics#Counter) `metric:"traces" tags:"state=started,sampled=y" help:"Number of traces started by this tracer as sampled"`
    
     // Number of traces started by this tracer as not sampled
     TracesStartedNotSampled [metrics](/github.com/uber/jaeger-lib/metrics).[Counter](/github.com/uber/jaeger-lib/metrics#Counter) `metric:"traces" tags:"state=started,sampled=n" help:"Number of traces started by this tracer as not sampled"`
    
     // Number of traces started by this tracer with delayed sampling
     TracesStartedDelayedSampling [metrics](/github.com/uber/jaeger-lib/metrics).[Counter](/github.com/uber/jaeger-lib/metrics#Counter) `metric:"traces" tags:"state=started,sampled=n" help:"Number of traces started by this tracer with delayed sampling"`
    
     // Number of externally started sampled traces this tracer joined
     TracesJoinedSampled [metrics](/github.com/uber/jaeger-lib/metrics).[Counter](/github.com/uber/jaeger-lib/metrics#Counter) `metric:"traces" tags:"state=joined,sampled=y" help:"Number of externally started sampled traces this tracer joined"`
    
     // Number of externally started not-sampled traces this tracer joined
     TracesJoinedNotSampled [metrics](/github.com/uber/jaeger-lib/metrics).[Counter](/github.com/uber/jaeger-lib/metrics#Counter) `metric:"traces" tags:"state=joined,sampled=n" help:"Number of externally started not-sampled traces this tracer joined"`
    
     // Number of sampled spans started by this tracer
     SpansStartedSampled [metrics](/github.com/uber/jaeger-lib/metrics).[Counter](/github.com/uber/jaeger-lib/metrics#Counter) `metric:"started_spans" tags:"sampled=y" help:"Number of spans started by this tracer as sampled"`
    
     // Number of not sampled spans started by this tracer
     SpansStartedNotSampled [metrics](/github.com/uber/jaeger-lib/metrics).[Counter](/github.com/uber/jaeger-lib/metrics#Counter) `metric:"started_spans" tags:"sampled=n" help:"Number of spans started by this tracer as not sampled"`
    
     // Number of spans with delayed sampling started by this tracer
     SpansStartedDelayedSampling [metrics](/github.com/uber/jaeger-lib/metrics).[Counter](/github.com/uber/jaeger-lib/metrics#Counter) `metric:"started_spans" tags:"sampled=delayed" help:"Number of spans started by this tracer with delayed sampling"`
    
     // Number of spans finished by this tracer
     SpansFinishedSampled [metrics](/github.com/uber/jaeger-lib/metrics).[Counter](/github.com/uber/jaeger-lib/metrics#Counter) `metric:"finished_spans" tags:"sampled=y" help:"Number of sampled spans finished by this tracer"`
    
     // Number of spans finished by this tracer
     SpansFinishedNotSampled [metrics](/github.com/uber/jaeger-lib/metrics).[Counter](/github.com/uber/jaeger-lib/metrics#Counter) `metric:"finished_spans" tags:"sampled=n" help:"Number of not-sampled spans finished by this tracer"`
    
     // Number of spans finished by this tracer
     SpansFinishedDelayedSampling [metrics](/github.com/uber/jaeger-lib/metrics).[Counter](/github.com/uber/jaeger-lib/metrics#Counter) `metric:"finished_spans" tags:"sampled=delayed" help:"Number of spans with delayed sampling finished by this tracer"`
    
     // Number of errors decoding tracing context
     DecodingErrors [metrics](/github.com/uber/jaeger-lib/metrics).[Counter](/github.com/uber/jaeger-lib/metrics#Counter) `metric:"span_context_decoding_errors" help:"Number of errors decoding tracing context"`
    
     // Number of spans successfully reported
     ReporterSuccess [metrics](/github.com/uber/jaeger-lib/metrics).[Counter](/github.com/uber/jaeger-lib/metrics#Counter) `metric:"reporter_spans" tags:"result=ok" help:"Number of spans successfully reported"`
    
     // Number of spans not reported due to a Sender failure
     ReporterFailure [metrics](/github.com/uber/jaeger-lib/metrics).[Counter](/github.com/uber/jaeger-lib/metrics#Counter) `metric:"reporter_spans" tags:"result=err" help:"Number of spans not reported due to a Sender failure"`
    
     // Number of spans dropped due to internal queue overflow
     ReporterDropped [metrics](/github.com/uber/jaeger-lib/metrics).[Counter](/github.com/uber/jaeger-lib/metrics#Counter) `metric:"reporter_spans" tags:"result=dropped" help:"Number of spans dropped due to internal queue overflow"`
    
     // Current number of spans in the reporter queue
     ReporterQueueLength [metrics](/github.com/uber/jaeger-lib/metrics).[Gauge](/github.com/uber/jaeger-lib/metrics#Gauge) `metric:"reporter_queue_length" help:"Current number of spans in the reporter queue"`
    
     // Number of times the Sampler succeeded to retrieve sampling strategy
     SamplerRetrieved [metrics](/github.com/uber/jaeger-lib/metrics).[Counter](/github.com/uber/jaeger-lib/metrics#Counter) `metric:"sampler_queries" tags:"result=ok" help:"Number of times the Sampler succeeded to retrieve sampling strategy"`
    
     // Number of times the Sampler failed to retrieve sampling strategy
     SamplerQueryFailure [metrics](/github.com/uber/jaeger-lib/metrics).[Counter](/github.com/uber/jaeger-lib/metrics#Counter) `metric:"sampler_queries" tags:"result=err" help:"Number of times the Sampler failed to retrieve sampling strategy"`
    
     // Number of times the Sampler succeeded to retrieve and update sampling strategy
     SamplerUpdated [metrics](/github.com/uber/jaeger-lib/metrics).[Counter](/github.com/uber/jaeger-lib/metrics#Counter) `` /* 127-byte string literal not displayed */
    
     // Number of times the Sampler failed to update sampling strategy
     SamplerUpdateFailure [metrics](/github.com/uber/jaeger-lib/metrics).[Counter](/github.com/uber/jaeger-lib/metrics#Counter) `metric:"sampler_updates" tags:"result=err" help:"Number of times the Sampler failed to update sampling strategy"`
    
     // Number of times baggage was successfully written or updated on spans.
     BaggageUpdateSuccess [metrics](/github.com/uber/jaeger-lib/metrics).[Counter](/github.com/uber/jaeger-lib/metrics#Counter) `metric:"baggage_updates" tags:"result=ok" help:"Number of times baggage was successfully written or updated on spans"`
    
     // Number of times baggage failed to write or update on spans.
     BaggageUpdateFailure [metrics](/github.com/uber/jaeger-lib/metrics).[Counter](/github.com/uber/jaeger-lib/metrics#Counter) `metric:"baggage_updates" tags:"result=err" help:"Number of times baggage failed to write or update on spans"`
    
     // Number of times baggage was truncated as per baggage restrictions.
     BaggageTruncate [metrics](/github.com/uber/jaeger-lib/metrics).[Counter](/github.com/uber/jaeger-lib/metrics#Counter) `metric:"baggage_truncations" help:"Number of times baggage was truncated as per baggage restrictions"`
    
     // Number of times baggage restrictions were successfully updated.
     BaggageRestrictionsUpdateSuccess [metrics](/github.com/uber/jaeger-lib/metrics).[Counter](/github.com/uber/jaeger-lib/metrics#Counter) `metric:"baggage_restrictions_updates" tags:"result=ok" help:"Number of times baggage restrictions were successfully updated"`
    
     // Number of times baggage restrictions failed to update.
     BaggageRestrictionsUpdateFailure [metrics](/github.com/uber/jaeger-lib/metrics).[Counter](/github.com/uber/jaeger-lib/metrics#Counter) `metric:"baggage_restrictions_updates" tags:"result=err" help:"Number of times baggage restrictions failed to update"`
    
     // Number of times debug spans were throttled.
     ThrottledDebugSpans [metrics](/github.com/uber/jaeger-lib/metrics).[Counter](/github.com/uber/jaeger-lib/metrics#Counter) `metric:"throttled_debug_spans" help:"Number of times debug spans were throttled"`
    
     // Number of times throttler successfully updated.
     ThrottlerUpdateSuccess [metrics](/github.com/uber/jaeger-lib/metrics).[Counter](/github.com/uber/jaeger-lib/metrics#Counter) `metric:"throttler_updates" tags:"result=ok" help:"Number of times throttler successfully updated"`
    
     // Number of times throttler failed to update.
     ThrottlerUpdateFailure [metrics](/github.com/uber/jaeger-lib/metrics).[Counter](/github.com/uber/jaeger-lib/metrics#Counter) `metric:"throttler_updates" tags:"result=err" help:"Number of times throttler failed to update"`
    }

Metrics is a container of all stats emitted by Jaeger tracer.

#### func [NewMetrics](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/metrics.go#L109) ¶

    func NewMetrics(factory [metrics](/github.com/uber/jaeger-lib/metrics).[Factory](/github.com/uber/jaeger-lib/metrics#Factory), globalTags map[[string](/builtin#string)][string](/builtin#string)) *Metrics

NewMetrics creates a new Metrics struct and initializes it.

#### func [NewNullMetrics](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/metrics.go#L117) ¶

    func NewNullMetrics() *Metrics

NewNullMetrics creates a new Metrics struct that won't report any metrics.

#### type [Observer](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/observer.go#L23) deprecated

    type Observer interface {
     OnStartSpan(operationName [string](/builtin#string), options [opentracing](/github.com/opentracing/opentracing-go).[StartSpanOptions](/github.com/opentracing/opentracing-go#StartSpanOptions)) SpanObserver
    }

Observer can be registered with the Tracer to receive notifications about new Spans.

Deprecated: use jaeger.ContribObserver instead.

#### type [PerOperationSampler](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/sampler.go#L331) ¶

    type PerOperationSampler struct {
     [sync](/sync).[RWMutex](/sync#RWMutex)
     // contains filtered or unexported fields
    }

PerOperationSampler is a delegating sampler that applies GuaranteedThroughputProbabilisticSampler on a per-operation basis.

#### func [NewAdaptiveSampler](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/sampler.go#L345) ¶

    func NewAdaptiveSampler(strategies *[sampling](/github.com/uber/jaeger-client-go/thrift-gen/sampling).[PerOperationSamplingStrategies](/github.com/uber/jaeger-client-go/thrift-gen/sampling#PerOperationSamplingStrategies), maxOperations [int](/builtin#int)) (*PerOperationSampler, [error](/builtin#error))

NewAdaptiveSampler returns a new PerOperationSampler. Deprecated: please use NewPerOperationSampler.

#### func [NewPerOperationSampler](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/sampler.go#L370) ¶

    func NewPerOperationSampler(params PerOperationSamplerParams) *PerOperationSampler

NewPerOperationSampler returns a new PerOperationSampler.

#### func (*PerOperationSampler) [Close](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/sampler.go#L455) ¶

    func (s *PerOperationSampler) Close()

Close invokes Close on all underlying samplers.

#### func (*PerOperationSampler) [Equal](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/sampler.go#L483) ¶

    func (s *PerOperationSampler) Equal(other Sampler) [bool](/builtin#bool)

Equal is not used. TODO (breaking change) remove this in the future

#### func (*PerOperationSampler) [IsSampled](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/sampler.go#L393) ¶

    func (s *PerOperationSampler) IsSampled(id TraceID, operation [string](/builtin#string)) ([bool](/builtin#bool), []Tag)

IsSampled is not used and only exists to match Sampler V1 API. TODO (breaking change) remove when upgrading everything to SamplerV2

#### func (*PerOperationSampler) [OnCreateSpan](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/sampler.go#L408) ¶

    func (s *PerOperationSampler) OnCreateSpan(span *Span) SamplingDecision

OnCreateSpan implements OnCreateSpan of SamplerV2.

#### func (*PerOperationSampler) [OnFinishSpan](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/sampler.go#L425) ¶

    func (s *PerOperationSampler) OnFinishSpan(span *Span) SamplingDecision

OnFinishSpan implements OnFinishSpan of SamplerV2.

#### func (*PerOperationSampler) [OnSetOperationName](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/sampler.go#L414) ¶

    func (s *PerOperationSampler) OnSetOperationName(span *Span, operationName [string](/builtin#string)) SamplingDecision

OnSetOperationName implements OnSetOperationName of SamplerV2.

#### func (*PerOperationSampler) [OnSetTag](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/sampler.go#L420) ¶

    func (s *PerOperationSampler) OnSetTag(span *Span, key [string](/builtin#string), value interface{}) SamplingDecision

OnSetTag implements OnSetTag of SamplerV2.

#### func (*PerOperationSampler) [String](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/sampler.go#L464) ¶

    func (s *PerOperationSampler) String() [string](/builtin#string)

#### type [PerOperationSamplerParams](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/sampler.go#L353) ¶

    type PerOperationSamplerParams struct {
     // Max number of operations that will be tracked. Other operations will be given default strategy.
     MaxOperations [int](/builtin#int)
    
     // Opt-in feature for applications that require late binding of span name via explicit call to SetOperationName.
     // When this feature is enabled, the sampler will return retryable=true from OnCreateSpan(), thus leaving
     // the sampling decision as non-final (and the span as writeable). This may lead to degraded performance
     // in applications that always provide the correct span name on trace creation.
     //
     // For backwards compatibility this option is off by default.
     OperationNameLateBinding [bool](/builtin#bool)
    
     // Initial configuration of the sampling strategies (usually retrieved from the backend by Remote Sampler).
     Strategies *[sampling](/github.com/uber/jaeger-client-go/thrift-gen/sampling).[PerOperationSamplingStrategies](/github.com/uber/jaeger-client-go/thrift-gen/sampling#PerOperationSamplingStrategies)
    }

PerOperationSamplerParams defines parameters when creating PerOperationSampler.

#### type [ProbabilisticSampler](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/sampler.go#L100) ¶

    type ProbabilisticSampler struct {
     // contains filtered or unexported fields
    }

ProbabilisticSampler is a sampler that randomly samples a certain percentage of traces.

#### func [NewProbabilisticSampler](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/sampler.go#L115) ¶

    func NewProbabilisticSampler(samplingRate [float64](/builtin#float64)) (*ProbabilisticSampler, [error](/builtin#error))

NewProbabilisticSampler creates a sampler that randomly samples a certain percentage of traces specified by the samplingRate, in the range between 0.0 and 1.0.

It relies on the fact that new trace IDs are 63bit random numbers themselves, thus making the sampling decision without generating a new random number, but simply calculating if traceID < (samplingRate * 2^63). TODO remove the error from this function for next major release

#### func (*ProbabilisticSampler) [Close](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/sampler.go#L149) ¶

    func (s *ProbabilisticSampler) Close()

Close implements Close() of Sampler.

#### func (*ProbabilisticSampler) [Equal](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/sampler.go#L154) ¶

    func (s *ProbabilisticSampler) Equal(other Sampler) [bool](/builtin#bool)

Equal implements Equal() of Sampler.

#### func (*ProbabilisticSampler) [IsSampled](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/sampler.go#L144) ¶

    func (s *ProbabilisticSampler) IsSampled(id TraceID, operation [string](/builtin#string)) ([bool](/builtin#bool), []Tag)

IsSampled implements IsSampled() of Sampler.

#### func (*ProbabilisticSampler) [OnCreateSpan](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/sampler_v2.go#L75) ¶

    func (s *ProbabilisticSampler) OnCreateSpan(span *Span) SamplingDecision

#### func (*ProbabilisticSampler) [OnFinishSpan](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/sampler_v2.go#L89) ¶

    func (s *ProbabilisticSampler) OnFinishSpan(span *Span) SamplingDecision

#### func (*ProbabilisticSampler) [OnSetOperationName](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/sampler_v2.go#L80) ¶

    func (s *ProbabilisticSampler) OnSetOperationName(span *Span, operationName [string](/builtin#string)) SamplingDecision

#### func (*ProbabilisticSampler) [OnSetTag](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/sampler_v2.go#L85) ¶

    func (s *ProbabilisticSampler) OnSetTag(span *Span, key [string](/builtin#string), value interface{}) SamplingDecision

#### func (*ProbabilisticSampler) [SamplingRate](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/sampler.go#L139) ¶

    func (s *ProbabilisticSampler) SamplingRate() [float64](/builtin#float64)

SamplingRate returns the sampling probability this sampled was constructed with.

#### func (*ProbabilisticSampler) [String](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/sampler.go#L171) ¶

    func (s *ProbabilisticSampler) String() [string](/builtin#string)

String is used to log sampler details.

#### func (*ProbabilisticSampler) [Update](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/sampler.go#L162) ¶

    func (s *ProbabilisticSampler) Update(samplingRate [float64](/builtin#float64)) [error](/builtin#error)

Update modifies in-place the sampling rate. Locking must be done externally.

#### type [ProbabilisticSamplerUpdater](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/sampler_remote.go#L222) ¶

    type ProbabilisticSamplerUpdater struct{}

ProbabilisticSamplerUpdater is used by RemotelyControlledSampler to parse sampling configuration.

#### func (*ProbabilisticSamplerUpdater) [Update](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/sampler_remote.go#L225) ¶

    func (u *ProbabilisticSamplerUpdater) Update(sampler SamplerV2, strategy interface{}) (SamplerV2, [error](/builtin#error))

Update implements Update of SamplerUpdater.

#### type [Process](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/process.go#L18) ¶

    type Process struct {
     Service [string](/builtin#string)
     UUID    [string](/builtin#string)
     Tags    []Tag
    }

Process holds process specific metadata that's relevant to this client.

#### type [ProcessSetter](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/process.go#L27) ¶

    type ProcessSetter interface {
     SetProcess(process Process)
    }

ProcessSetter sets a process. This can be used by any class that requires the process to be set as part of initialization. See internal/throttler/remote/throttler.go for an example.

#### type [RateLimitingSampler](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/sampler.go#L181) ¶

    type RateLimitingSampler struct {
     // contains filtered or unexported fields
    }

RateLimitingSampler samples at most maxTracesPerSecond. The distribution of sampled traces follows burstiness of the service, i.e. a service with uniformly distributed requests will have those requests sampled uniformly as well, but if requests are bursty, especially sub-second, then a number of sequential requests can be sampled each second.

#### func [NewRateLimitingSampler](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/sampler.go#L189) ¶

    func NewRateLimitingSampler(maxTracesPerSecond [float64](/builtin#float64)) *RateLimitingSampler

NewRateLimitingSampler creates new RateLimitingSampler.

#### func (*RateLimitingSampler) [Close](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/sampler.go#L223) ¶

    func (s *RateLimitingSampler) Close()

Close does nothing.

#### func (*RateLimitingSampler) [Equal](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/sampler.go#L228) ¶

    func (s *RateLimitingSampler) Equal(other Sampler) [bool](/builtin#bool)

Equal compares with another sampler.

#### func (*RateLimitingSampler) [IsSampled](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/sampler.go#L210) ¶

    func (s *RateLimitingSampler) IsSampled(id TraceID, operation [string](/builtin#string)) ([bool](/builtin#bool), []Tag)

IsSampled implements IsSampled() of Sampler.

#### func (*RateLimitingSampler) [OnCreateSpan](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/sampler_v2.go#L75) ¶

    func (s *RateLimitingSampler) OnCreateSpan(span *Span) SamplingDecision

#### func (*RateLimitingSampler) [OnFinishSpan](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/sampler_v2.go#L89) ¶

    func (s *RateLimitingSampler) OnFinishSpan(span *Span) SamplingDecision

#### func (*RateLimitingSampler) [OnSetOperationName](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/sampler_v2.go#L80) ¶

    func (s *RateLimitingSampler) OnSetOperationName(span *Span, operationName [string](/builtin#string)) SamplingDecision

#### func (*RateLimitingSampler) [OnSetTag](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/sampler_v2.go#L85) ¶

    func (s *RateLimitingSampler) OnSetTag(span *Span, key [string](/builtin#string), value interface{}) SamplingDecision

#### func (*RateLimitingSampler) [String](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/sampler.go#L236) ¶

    func (s *RateLimitingSampler) String() [string](/builtin#string)

String is used to log sampler details.

#### func (*RateLimitingSampler) [Update](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/sampler.go#L216) ¶

    func (s *RateLimitingSampler) Update(maxTracesPerSecond [float64](/builtin#float64))

Update reconfigures the rate limiter, while preserving its accumulated balance. Locking must be done externally.

#### type [RateLimitingSamplerUpdater](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/sampler_remote.go#L247) ¶

    type RateLimitingSamplerUpdater struct{}

RateLimitingSamplerUpdater is used by RemotelyControlledSampler to parse sampling configuration.

#### func (*RateLimitingSamplerUpdater) [Update](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/sampler_remote.go#L250) ¶

    func (u *RateLimitingSamplerUpdater) Update(sampler SamplerV2, strategy interface{}) (SamplerV2, [error](/builtin#error))

Update implements Update of SamplerUpdater.

#### type [Reference](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/reference.go#L20) ¶

    type Reference struct {
     Type    [opentracing](/github.com/opentracing/opentracing-go).[SpanReferenceType](/github.com/opentracing/opentracing-go#SpanReferenceType)
     Context SpanContext
    }

Reference represents a causal reference to other Spans (via their SpanContext).

#### type [RemotelyControlledSampler](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/sampler_remote.go#L63) ¶

    type RemotelyControlledSampler struct {
     [sync](/sync).[RWMutex](/sync#RWMutex) // used to serialize access to samplerOptions.sampler
     // contains filtered or unexported fields
    }

RemotelyControlledSampler is a delegating sampler that polls a remote server for the appropriate sampling strategy, constructs a corresponding sampler and delegates to it for sampling decisions.

#### func [NewRemotelyControlledSampler](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/sampler_remote.go#L77) ¶

    func NewRemotelyControlledSampler(
     serviceName [string](/builtin#string),
     opts ...SamplerOption,
    ) *RemotelyControlledSampler

NewRemotelyControlledSampler creates a sampler that periodically pulls the sampling strategy from an HTTP sampling server (e.g. jaeger-agent).

#### func (*RemotelyControlledSampler) [Close](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/sampler_remote.go#L126) ¶

    func (s *RemotelyControlledSampler) Close()

Close implements Close() of Sampler.

#### func (*RemotelyControlledSampler) [Equal](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/sampler_remote.go#L139) ¶

    func (s *RemotelyControlledSampler) Equal(other Sampler) [bool](/builtin#bool)

Equal implements Equal() of Sampler.

#### func (*RemotelyControlledSampler) [IsSampled](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/sampler_remote.go#L93) ¶

    func (s *RemotelyControlledSampler) IsSampled(id TraceID, operation [string](/builtin#string)) ([bool](/builtin#bool), []Tag)

IsSampled implements IsSampled() of Sampler. TODO (breaking change) remove when Sampler V1 is removed

#### func (*RemotelyControlledSampler) [OnCreateSpan](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/sampler_remote.go#L98) ¶

    func (s *RemotelyControlledSampler) OnCreateSpan(span *Span) SamplingDecision

OnCreateSpan implements OnCreateSpan of SamplerV2.

#### func (*RemotelyControlledSampler) [OnFinishSpan](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/sampler_remote.go#L119) ¶

    func (s *RemotelyControlledSampler) OnFinishSpan(span *Span) SamplingDecision

OnFinishSpan implements OnFinishSpan of SamplerV2.

#### func (*RemotelyControlledSampler) [OnSetOperationName](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/sampler_remote.go#L105) ¶

    func (s *RemotelyControlledSampler) OnSetOperationName(span *Span, operationName [string](/builtin#string)) SamplingDecision

OnSetOperationName implements OnSetOperationName of SamplerV2.

#### func (*RemotelyControlledSampler) [OnSetTag](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/sampler_remote.go#L112) ¶

    func (s *RemotelyControlledSampler) OnSetTag(span *Span, key [string](/builtin#string), value interface{}) SamplingDecision

OnSetTag implements OnSetTag of SamplerV2.

#### func (*RemotelyControlledSampler) [Sampler](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/sampler_remote.go#L164) ¶

    func (s *RemotelyControlledSampler) Sampler() SamplerV2

Sampler returns the currently active sampler.

#### func (*RemotelyControlledSampler) [UpdateSampler](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/sampler_remote.go#L177) ¶

    func (s *RemotelyControlledSampler) UpdateSampler()

UpdateSampler forces the sampler to fetch sampling strategy from backend server. This function is called automatically on a timer, but can also be safely called manually, e.g. from tests.

#### type [Reporter](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/reporter.go#L30) ¶

    type Reporter interface {
     // Report submits a new span to collectors, possibly asynchronously and/or with buffering.
     // If the reporter is processing Span asynchronously then it needs to Retain() the span,
     // and then Release() it when no longer needed, to avoid span data corruption.
     Report(span *Span)
    
     // Close does a clean shutdown of the reporter, flushing any traces that may be buffered in memory.
     Close()
    }

Reporter is called by the tracer when a span is completed to report the span to the tracing collector.

#### func [NewCompositeReporter](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/reporter.go#L144) ¶

    func NewCompositeReporter(reporters ...Reporter) Reporter

NewCompositeReporter creates a reporter that ignores all reported spans.

#### func [NewLoggingReporter](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/reporter.go#L66) ¶

    func NewLoggingReporter(logger Logger) Reporter

NewLoggingReporter creates a reporter that logs all reported spans to provided logger.

#### func [NewNullReporter](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/reporter.go#L45) ¶

    func NewNullReporter() Reporter

NewNullReporter creates a no-op reporter that ignores all reported spans.

#### func [NewRemoteReporter](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/reporter.go#L211) ¶

    func NewRemoteReporter(sender Transport, opts ...ReporterOption) Reporter

NewRemoteReporter creates a new reporter that sends spans out of process by means of Sender. Calls to Report(Span) return immediately (side effect: if internal buffer is full the span is dropped). Periodically the transport buffer is flushed even if it hasn't reached max packet size. Calls to Close() block until all spans reported prior to the call to Close are flushed.

#### type [ReporterOption](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/reporter_options.go#L24) ¶

    type ReporterOption func(c *reporterOptions)

ReporterOption is a function that sets some option on the reporter.

#### type [Sampler](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/sampler.go#L32) ¶

    type Sampler interface {
     // IsSampled decides whether a trace with given `id` and `operation`
     // should be sampled. This function will also return the tags that
     // can be used to identify the type of sampling that was applied to
     // the root span. Most simple samplers would return two tags,
     // sampler.type and sampler.param, similar to those used in the Configuration
     IsSampled(id TraceID, operation [string](/builtin#string)) (sampled [bool](/builtin#bool), tags []Tag)
    
     // Close does a clean shutdown of the sampler, stopping any background
     // go-routines it may have started.
     Close()
    
     // Equal checks if the `other` sampler is functionally equivalent
     // to this sampler.
     // TODO (breaking change) remove this function. See PerOperationSampler.Equals for explanation.
     Equal(other Sampler) [bool](/builtin#bool)
    }

Sampler decides whether a new trace should be sampled or not.

#### type [SamplerOption](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/sampler_remote_options.go#L24) ¶

    type SamplerOption func(options *samplerOptions)

SamplerOption is a function that sets some option on the sampler

#### type [SamplerOptionsFactory](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/sampler_remote_options.go#L33) ¶

    type SamplerOptionsFactory struct{}

SamplerOptionsFactory is a factory for all available SamplerOption's. The type acts as a namespace for factory functions. It is public to make the functions discoverable via godoc. Recommended to be used via global SamplerOptions variable.

    var SamplerOptions SamplerOptionsFactory

SamplerOptions is a factory for all available SamplerOption's.

#### func (SamplerOptionsFactory) [InitialSampler](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/sampler_remote_options.go#L73) ¶

    func (SamplerOptionsFactory) InitialSampler(sampler Sampler) SamplerOption

InitialSampler creates a SamplerOption that sets the initial sampler to use before a remote sampler is created and used.

#### func (SamplerOptionsFactory) [Logger](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/sampler_remote_options.go#L80) ¶

    func (SamplerOptionsFactory) Logger(logger Logger) SamplerOption

Logger creates a SamplerOption that sets the logger used by the sampler.

#### func (SamplerOptionsFactory) [MaxOperations](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/sampler_remote_options.go#L57) ¶

    func (SamplerOptionsFactory) MaxOperations(maxOperations [int](/builtin#int)) SamplerOption

MaxOperations creates a SamplerOption that sets the maximum number of operations the sampler will keep track of.

#### func (SamplerOptionsFactory) [Metrics](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/sampler_remote_options.go#L49) ¶

    func (SamplerOptionsFactory) Metrics(m *Metrics) SamplerOption

Metrics creates a SamplerOption that initializes Metrics on the sampler, which is used to emit statistics.

#### func (SamplerOptionsFactory) [OperationNameLateBinding](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/sampler_remote_options.go#L65) ¶

    func (SamplerOptionsFactory) OperationNameLateBinding(enable [bool](/builtin#bool)) SamplerOption

OperationNameLateBinding creates a SamplerOption that sets the respective field in the PerOperationSamplerParams.

#### func (SamplerOptionsFactory) [SamplingRefreshInterval](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/sampler_remote_options.go#L96) ¶

    func (SamplerOptionsFactory) SamplingRefreshInterval(samplingRefreshInterval [time](/time).[Duration](/time#Duration)) SamplerOption

SamplingRefreshInterval creates a SamplerOption that sets how often the sampler will poll local agent for the appropriate sampling strategy.

#### func (SamplerOptionsFactory) [SamplingServerURL](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/sampler_remote_options.go#L88) ¶

    func (SamplerOptionsFactory) SamplingServerURL(samplingServerURL [string](/builtin#string)) SamplerOption

SamplingServerURL creates a SamplerOption that sets the sampling server url of the local agent that contains the sampling strategies.

#### func (SamplerOptionsFactory) [SamplingStrategyFetcher](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/sampler_remote_options.go#L103) ¶

    func (SamplerOptionsFactory) SamplingStrategyFetcher(fetcher SamplingStrategyFetcher) SamplerOption

SamplingStrategyFetcher creates a SamplerOption that initializes sampling strategy fetcher.

#### func (SamplerOptionsFactory) [SamplingStrategyParser](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/sampler_remote_options.go#L110) ¶

    func (SamplerOptionsFactory) SamplingStrategyParser(parser SamplingStrategyParser) SamplerOption

SamplingStrategyParser creates a SamplerOption that initializes sampling strategy parser.

#### func (SamplerOptionsFactory) [Updaters](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/sampler_remote_options.go#L117) ¶

    func (SamplerOptionsFactory) Updaters(updaters ...SamplerUpdater) SamplerOption

Updaters creates a SamplerOption that initializes sampler updaters.

#### type [SamplerUpdater](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/sampler_remote.go#L56) ¶

    type SamplerUpdater interface {
     Update(sampler SamplerV2, strategy interface{}) (modified SamplerV2, err [error](/builtin#error))
    }

SamplerUpdater is used by RemotelyControlledSampler to apply sampling strategies, retrieved from remote config server, to the current sampler. The updater can modify the sampler in-place if sampler supports it, or create a new one.

If the strategy does not contain configuration for the sampler in question, updater must return modifiedSampler=nil to give other updaters a chance to inspect the sampling strategy response.

RemotelyControlledSampler invokes the updaters while holding a lock on the main sampler.

#### type [SamplerV2](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/sampler_v2.go#L26) ¶

    type SamplerV2 interface {
     OnCreateSpan(span *Span) SamplingDecision
     OnSetOperationName(span *Span, operationName [string](/builtin#string)) SamplingDecision
     OnSetTag(span *Span, key [string](/builtin#string), value interface{}) SamplingDecision
     OnFinishSpan(span *Span) SamplingDecision
    
     // Close does a clean shutdown of the sampler, stopping any background
     // go-routines it may have started.
     Close()
    }

SamplerV2 is an extension of the V1 samplers that allows sampling decisions be made at different points of the span lifecycle.

#### type [SamplerV2Base](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/sampler_v2.go#L56) ¶

    type SamplerV2Base struct{}

SamplerV2Base can be used by V2 samplers to implement dummy V1 methods. Supporting V1 API is required because Tracer configuration only accepts V1 Sampler for backwards compatibility reasons. TODO (breaking change) remove this in the next major release

#### func (SamplerV2Base) [Close](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/sampler_v2.go#L64) ¶

    func (SamplerV2Base) Close()

Close implements Close of Sampler.

#### func (SamplerV2Base) [Equal](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/sampler_v2.go#L67) ¶

    func (SamplerV2Base) Equal(other Sampler) [bool](/builtin#bool)

Equal implements Equal of Sampler.

#### func (SamplerV2Base) [IsSampled](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/sampler_v2.go#L59) ¶

    func (SamplerV2Base) IsSampled(id TraceID, operation [string](/builtin#string)) (sampled [bool](/builtin#bool), tags []Tag)

IsSampled implements IsSampled of Sampler.

#### type [SamplingDecision](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/sampler_v2.go#L18) ¶

    type SamplingDecision struct {
     Sample    [bool](/builtin#bool)
     Retryable [bool](/builtin#bool)
     Tags      []Tag
    }

SamplingDecision is returned by the V2 samplers.

#### type [SamplingStrategyFetcher](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/sampler_remote.go#L37) ¶

    type SamplingStrategyFetcher interface {
     Fetch(service [string](/builtin#string)) ([][byte](/builtin#byte), [error](/builtin#error))
    }

SamplingStrategyFetcher is used to fetch sampling strategy updates from remote server.

#### type [SamplingStrategyParser](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/sampler_remote.go#L43) ¶

    type SamplingStrategyParser interface {
     Parse(response [][byte](/builtin#byte)) (interface{}, [error](/builtin#error))
    }

SamplingStrategyParser is used to parse sampling strategy updates. The output object should be of the type that is recognized by the SamplerUpdaters.

#### type [Span](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/span.go#L28) ¶

    type Span struct {
     [sync](/sync).[RWMutex](/sync#RWMutex)
     // contains filtered or unexported fields
    }

Span implements opentracing.Span

#### func (*Span) [BaggageItem](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/span.go#L333) ¶

    func (s *Span) BaggageItem(key [string](/builtin#string)) [string](/builtin#string)

BaggageItem implements BaggageItem() of opentracing.SpanContext

#### func (*Span) [Context](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/span.go#L380) ¶

    func (s *Span) Context() [opentracing](/github.com/opentracing/opentracing-go).[SpanContext](/github.com/opentracing/opentracing-go#SpanContext)

Context implements opentracing.Span API

#### func (*Span) [Duration](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/span.go#L147) ¶

    func (s *Span) Duration() [time](/time).[Duration](/time#Duration)

Duration returns span duration

#### func (*Span) [Finish](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/span.go#L342) ¶

    func (s *Span) Finish()

Finish implements opentracing.Span API After finishing the Span object it returns back to the allocator unless the reporter retains it again, so after that, the Span object should no longer be used because it won't be valid anymore.

#### func (*Span) [FinishWithOptions](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/span.go#L347) ¶

    func (s *Span) FinishWithOptions(options [opentracing](/github.com/opentracing/opentracing-go).[FinishOptions](/github.com/opentracing/opentracing-go#FinishOptions))

FinishWithOptions implements opentracing.Span API

#### func (*Span) [Log](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/span.go#L243) ¶

    func (s *Span) Log(ld [opentracing](/github.com/opentracing/opentracing-go).[LogData](/github.com/opentracing/opentracing-go#LogData))

Log implements opentracing.Span API

#### func (*Span) [LogEvent](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/span.go#L233) ¶

    func (s *Span) LogEvent(event [string](/builtin#string))

LogEvent implements opentracing.Span API

#### func (*Span) [LogEventWithPayload](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/span.go#L238) ¶

    func (s *Span) LogEventWithPayload(event [string](/builtin#string), payload interface{})

LogEventWithPayload implements opentracing.Span API

#### func (*Span) [LogFields](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/span.go#L198) ¶

    func (s *Span) LogFields(fields ...[log](/github.com/opentracing/opentracing-go/log).[Field](/github.com/opentracing/opentracing-go/log#Field))

LogFields implements opentracing.Span API

#### func (*Span) [LogKV](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/span.go#L217) ¶

    func (s *Span) LogKV(alternatingKeyValues ...interface{})

LogKV implements opentracing.Span API

#### func (*Span) [Logs](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/span.go#L165) ¶

    func (s *Span) Logs() [][opentracing](/github.com/opentracing/opentracing-go).[LogRecord](/github.com/opentracing/opentracing-go#LogRecord)

Logs returns micro logs for span

#### func (*Span) [OperationName](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/span.go#L398) ¶

    func (s *Span) OperationName() [string](/builtin#string)

OperationName allows retrieving current operation name.

#### func (*Span) [References](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/span.go#L178) ¶

    func (s *Span) References() [][opentracing](/github.com/opentracing/opentracing-go).[SpanReference](/github.com/opentracing/opentracing-go#SpanReference)

References returns references for this span

#### func (*Span) [Release](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/span.go#L412) ¶

    func (s *Span) Release()

Release decrements object counter and return to the allocator manager when counter will below zero

#### func (*Span) [Retain](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/span.go#L405) ¶

    func (s *Span) Retain() *Span

Retain increases object counter to increase the lifetime of the object

#### func (*Span) [SetBaggageItem](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/span.go#L325) ¶

    func (s *Span) SetBaggageItem(key, value [string](/builtin#string)) [opentracing](/github.com/opentracing/opentracing-go).[Span](/github.com/opentracing/opentracing-go#Span)

SetBaggageItem implements SetBaggageItem() of opentracing.SpanContext. The call is proxied via tracer.baggageSetter to allow policies to be applied before allowing to set/replace baggage keys. The setter eventually stores a new SpanContext with extended baggage:

       span.context = span.context.WithBaggageItem(key, value)
    
    See SpanContext.WithBaggageItem() for explanation why it's done this way.
    

#### func (*Span) [SetOperationName](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/span.go#L85) ¶

    func (s *Span) SetOperationName(operationName [string](/builtin#string)) [opentracing](/github.com/opentracing/opentracing-go).[Span](/github.com/opentracing/opentracing-go#Span)

SetOperationName sets or changes the operation name.

#### func (*Span) [SetTag](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/span.go#L99) ¶

    func (s *Span) SetTag(key [string](/builtin#string), value interface{}) [opentracing](/github.com/opentracing/opentracing-go).[Span](/github.com/opentracing/opentracing-go#Span)

SetTag implements SetTag() of opentracing.Span

#### func (*Span) [SpanContext](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/span.go#L133) ¶

    func (s *Span) SpanContext() SpanContext

SpanContext returns span context

#### func (*Span) [StartTime](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/span.go#L140) ¶

    func (s *Span) StartTime() [time](/time).[Time](/time#Time)

StartTime returns span start time

#### func (*Span) [String](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/span.go#L391) ¶

    func (s *Span) String() [string](/builtin#string)

#### func (*Span) [Tags](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/span.go#L154) ¶

    func (s *Span) Tags() [opentracing](/github.com/opentracing/opentracing-go).[Tags](/github.com/opentracing/opentracing-go#Tags)

Tags returns tags for span

#### func (*Span) [Tracer](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/span.go#L387) ¶

    func (s *Span) Tracer() [opentracing](/github.com/opentracing/opentracing-go).[Tracer](/github.com/opentracing/opentracing-go#Tracer)

Tracer implements opentracing.Span API

#### type [SpanAllocator](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/span_allocator.go#L20) ¶

    type SpanAllocator interface {
     Get() *Span
     Put(*Span)
    }

SpanAllocator abstraction of managing span allocations

#### type [SpanContext](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/span_context.go#L49) ¶

    type SpanContext struct {
     // contains filtered or unexported fields
    }

SpanContext represents propagated span identity and state

#### func [ContextFromString](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/span_context.go#L226) ¶

    func ContextFromString(value [string](/builtin#string)) (SpanContext, [error](/builtin#error))

ContextFromString reconstructs the Context encoded in a string

#### func [NewSpanContext](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/span_context.go#L285) ¶

    func NewSpanContext(traceID TraceID, spanID, parentID SpanID, sampled [bool](/builtin#bool), baggage map[[string](/builtin#string)][string](/builtin#string)) SpanContext

NewSpanContext creates a new instance of SpanContext

#### func (*SpanContext) [CopyFrom](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/span_context.go#L301) ¶

    func (c *SpanContext) CopyFrom(ctx *SpanContext)

CopyFrom copies data from ctx into this context, including span identity and baggage. TODO This is only used by interop.go. Remove once TChannel Go supports OpenTracing.

#### func (SpanContext) [ExtendedSamplingState](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/span_context.go#L200) ¶

    func (c SpanContext) ExtendedSamplingState(key interface{}, initValue func() interface{}) interface{}

ExtendedSamplingState returns the custom state object for a given key. If the value for this key does not exist, it is initialized via initValue function. This state can be used by samplers (e.g. x.PrioritySampler).

#### func (SpanContext) [Flags](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/span_context.go#L270) ¶

    func (c SpanContext) Flags() [byte](/builtin#byte)

Flags returns the bitmap containing such bits as 'sampled' and 'debug'.

#### func (SpanContext) [ForeachBaggageItem](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/span_context.go#L169) ¶

    func (c SpanContext) ForeachBaggageItem(handler func(k, v [string](/builtin#string)) [bool](/builtin#bool))

ForeachBaggageItem implements ForeachBaggageItem() of opentracing.SpanContext

#### func (SpanContext) [IsDebug](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/span_context.go#L184) ¶

    func (c SpanContext) IsDebug() [bool](/builtin#bool)

IsDebug indicates whether sampling was explicitly requested by the service.

#### func (SpanContext) [IsFirehose](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/span_context.go#L194) ¶

    func (c SpanContext) IsFirehose() [bool](/builtin#bool)

IsFirehose indicates whether the firehose flag was set

#### func (SpanContext) [IsSampled](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/span_context.go#L179) ¶

    func (c SpanContext) IsSampled() [bool](/builtin#bool)

IsSampled returns whether this trace was chosen for permanent storage by the sampling mechanism of the tracer.

#### func (SpanContext) [IsSamplingFinalized](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/span_context.go#L189) ¶

    func (c SpanContext) IsSamplingFinalized() [bool](/builtin#bool)

IsSamplingFinalized indicates whether the sampling decision has been finalized.

#### func (SpanContext) [IsValid](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/span_context.go#L205) ¶ added in v1.4.1

    func (c SpanContext) IsValid() [bool](/builtin#bool)

IsValid indicates whether this context actually represents a valid trace.

#### func (SpanContext) [ParentID](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/span_context.go#L265) ¶

    func (c SpanContext) ParentID() SpanID

ParentID returns the parent span ID of this span context

#### func (SpanContext) [SetFirehose](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/span_context.go#L210) ¶

    func (c SpanContext) SetFirehose()

SetFirehose enables firehose mode for this trace.

#### func (SpanContext) [SpanID](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/span_context.go#L260) ¶

    func (c SpanContext) SpanID() SpanID

SpanID returns the span ID of this span context

#### func (SpanContext) [String](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/span_context.go#L214) ¶

    func (c SpanContext) String() [string](/builtin#string)

#### func (SpanContext) [TraceID](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/span_context.go#L255) ¶

    func (c SpanContext) TraceID() TraceID

TraceID returns the trace ID of this span context

#### func (SpanContext) [WithBaggageItem](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/span_context.go#L326) ¶

    func (c SpanContext) WithBaggageItem(key, value [string](/builtin#string)) SpanContext

WithBaggageItem creates a new context with an extra baggage item. Delete a baggage item if provided blank value.

The SpanContext is designed to be immutable and passed by value. As such, it cannot contain any locks, and should only hold immutable data, including baggage. Another reason for why baggage is immutable is when the span context is passed as a parent when starting a new span. The new span's baggage cannot affect the parent span's baggage, so the child span either needs to take a copy of the parent baggage (which is expensive and unnecessary since baggage rarely changes in the life span of a trace), or it needs to do a copy-on-write, which is the approach taken here.

#### type [SpanID](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/span_context.go#L46) ¶

    type SpanID [uint64](/builtin#uint64)

SpanID represents unique 64bit identifier of a span

#### func [SpanIDFromString](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/span_context.go#L409) ¶

    func SpanIDFromString(s [string](/builtin#string)) (SpanID, [error](/builtin#error))

SpanIDFromString creates a SpanID from a hexadecimal string

#### func (SpanID) [String](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/span_context.go#L404) ¶

    func (s SpanID) String() [string](/builtin#string)

#### type [SpanObserver](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/observer.go#L31) deprecated

    type SpanObserver interface {
     OnSetOperationName(operationName [string](/builtin#string))
     OnSetTag(key [string](/builtin#string), value interface{})
     OnFinish(options [opentracing](/github.com/opentracing/opentracing-go).[FinishOptions](/github.com/opentracing/opentracing-go#FinishOptions))
    }

SpanObserver is created by the Observer and receives notifications about other Span events.

Deprecated: use jaeger.ContribSpanObserver instead.

#### type [Tag](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/span.go#L73) ¶

    type Tag struct {
     // contains filtered or unexported fields
    }

Tag is a simple key value wrapper. TODO (breaking change) deprecate in the next major release, use opentracing.Tag instead.

#### func [NewTag](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/span.go#L80) ¶

    func NewTag(key [string](/builtin#string), value interface{}) Tag

NewTag creates a new Tag. TODO (breaking change) deprecate in the next major release, use opentracing.Tag instead.

#### type [TextMapPropagator](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/propagation.go#L55) ¶

    type TextMapPropagator struct {
     // contains filtered or unexported fields
    }

TextMapPropagator is a combined Injector and Extractor for TextMap format

#### func [NewHTTPHeaderPropagator](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/propagation.go#L77) ¶

    func NewHTTPHeaderPropagator(headerKeys *HeadersConfig, metrics Metrics) *TextMapPropagator

NewHTTPHeaderPropagator creates a combined Injector and Extractor for HTTPHeaders format

#### func [NewTextMapPropagator](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/propagation.go#L63) ¶

    func NewTextMapPropagator(headerKeys *HeadersConfig, metrics Metrics) *TextMapPropagator

NewTextMapPropagator creates a combined Injector and Extractor for TextMap format

#### func (*TextMapPropagator) [Extract](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/propagation.go#L131) ¶

    func (p *TextMapPropagator) Extract(abstractCarrier interface{}) (SpanContext, [error](/builtin#error))

Extract implements Extractor of TextMapPropagator

#### func (*TextMapPropagator) [Inject](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/propagation.go#L109) ¶

    func (p *TextMapPropagator) Inject(
     sc SpanContext,
     abstractCarrier interface{},
    ) [error](/builtin#error)

Inject implements Injector of TextMapPropagator

#### type [TraceID](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/span_context.go#L41) ¶

    type TraceID struct {
     High, Low [uint64](/builtin#uint64)
    }

TraceID represents unique 128bit identifier of a trace

#### func [TraceIDFromString](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/span_context.go#L376) ¶

    func TraceIDFromString(s [string](/builtin#string)) (TraceID, [error](/builtin#error))

TraceIDFromString creates a TraceID from a hexadecimal string

#### func (TraceID) [IsValid](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/span_context.go#L398) ¶

    func (t TraceID) IsValid() [bool](/builtin#bool)

IsValid checks if the trace ID is valid, i.e. not zero.

#### func (TraceID) [String](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/span_context.go#L368) ¶

    func (t TraceID) String() [string](/builtin#string)

#### type [Tracer](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/tracer.go#L37) ¶

    type Tracer struct {
     // contains filtered or unexported fields
    }

Tracer implements opentracing.Tracer.

#### func (*Tracer) [Close](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/tracer.go#L371) ¶

    func (t *Tracer) Close() [error](/builtin#error)

Close releases all resources used by the Tracer and flushes any remaining buffered spans.

#### func (*Tracer) [Extract](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/tracer.go#L355) ¶

    func (t *Tracer) Extract(
     format interface{},
     carrier interface{},
    ) ([opentracing](/github.com/opentracing/opentracing-go).[SpanContext](/github.com/opentracing/opentracing-go#SpanContext), [error](/builtin#error))

Extract implements Extract() method of opentracing.Tracer

#### func (*Tracer) [Inject](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/tracer.go#L343) ¶

    func (t *Tracer) Inject(ctx [opentracing](/github.com/opentracing/opentracing-go).[SpanContext](/github.com/opentracing/opentracing-go#SpanContext), format interface{}, carrier interface{}) [error](/builtin#error)

Inject implements Inject() method of opentracing.Tracer

#### func (*Tracer) [Sampler](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/tracer.go#L481) ¶

    func (t *Tracer) Sampler() SamplerV2

Sampler returns the sampler given to the tracer at creation.

#### func (*Tracer) [StartSpan](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/tracer.go#L200) ¶

    func (t *Tracer) StartSpan(
     operationName [string](/builtin#string),
     options ...[opentracing](/github.com/opentracing/opentracing-go).[StartSpanOption](/github.com/opentracing/opentracing-go#StartSpanOption),
    ) [opentracing](/github.com/opentracing/opentracing-go).[Span](/github.com/opentracing/opentracing-go#Span)

StartSpan implements StartSpan() method of opentracing.Tracer.

#### func (*Tracer) [Tags](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/tracer.go#L385) ¶

    func (t *Tracer) Tags() [][opentracing](/github.com/opentracing/opentracing-go).[Tag](/github.com/opentracing/opentracing-go#Tag)

Tags returns a slice of tracer-level tags.

#### type [TracerOption](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/tracer_options.go#L28) ¶

    type TracerOption func(tracer *Tracer)

TracerOption is a function that sets some option on the tracer

#### type [TracerOptionsFactory](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/tracer_options.go#L34) ¶

    type TracerOptionsFactory struct{}

TracerOptionsFactory is a struct that defines functions for all available TracerOption's.

    var TracerOptions TracerOptionsFactory

TracerOptions is a factory for all available TracerOption's.

#### func (TracerOptionsFactory) [BaggageRestrictionManager](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/tracer_options.go#L191) ¶

    func (TracerOptionsFactory) BaggageRestrictionManager(mgr [baggage](/github.com/uber/jaeger-client-go/internal/baggage).[RestrictionManager](/github.com/uber/jaeger-client-go/internal/baggage#RestrictionManager)) TracerOption

BaggageRestrictionManager registers BaggageRestrictionManager.

#### func (TracerOptionsFactory) [ContribObserver](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/tracer_options.go#L127) ¶

    func (TracerOptionsFactory) ContribObserver(observer ContribObserver) TracerOption

ContribObserver registers a ContribObserver.

#### func (TracerOptionsFactory) [CustomHeaderKeys](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/tracer_options.go#L53) ¶

    func (TracerOptionsFactory) CustomHeaderKeys(headerKeys *HeadersConfig) TracerOption

CustomHeaderKeys allows to override default HTTP header keys used to propagate tracing context.

#### func (TracerOptionsFactory) [DebugThrottler](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/tracer_options.go#L198) ¶

    func (TracerOptionsFactory) DebugThrottler(throttler [throttler](/github.com/uber/jaeger-client-go/internal/throttler).[Throttler](/github.com/uber/jaeger-client-go/internal/throttler#Throttler)) TracerOption

DebugThrottler registers a Throttler for debug spans.

#### func (TracerOptionsFactory) [Extractor](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/tracer_options.go#L115) ¶

    func (TracerOptionsFactory) Extractor(format interface{}, extractor Extractor) TracerOption

Extractor registers an Extractor for given format.

#### func (TracerOptionsFactory) [Gen128Bit](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/tracer_options.go#L134) ¶

    func (TracerOptionsFactory) Gen128Bit(gen128Bit [bool](/builtin#bool)) TracerOption

Gen128Bit enables generation of 128bit trace IDs.

#### func (TracerOptionsFactory) [HighTraceIDGenerator](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/tracer_options.go#L149) ¶

    func (TracerOptionsFactory) HighTraceIDGenerator(highTraceIDGenerator func() [uint64](/builtin#uint64)) TracerOption

HighTraceIDGenerator allows to override define ID generator.

#### func (TracerOptionsFactory) [HostIPv4](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/tracer_options.go#L101) ¶

    func (TracerOptionsFactory) HostIPv4(hostIPv4 [uint32](/builtin#uint32)) TracerOption

HostIPv4 creates a TracerOption that identifies the current service/process. If not set, the factory method will obtain the current IP address. The TracerOption is deprecated; the tracer will attempt to automatically detect the IP.

Deprecated.

#### func (TracerOptionsFactory) [Injector](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/tracer_options.go#L108) ¶

    func (TracerOptionsFactory) Injector(format interface{}, injector Injector) TracerOption

Injector registers a Injector for given format.

#### func (TracerOptionsFactory) [Logger](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/tracer_options.go#L45) ¶

    func (TracerOptionsFactory) Logger(logger Logger) TracerOption

Logger creates a TracerOption that gives the tracer a Logger.

#### func (TracerOptionsFactory) [MaxLogsPerSpan](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/tracer_options.go#L168) ¶

    func (TracerOptionsFactory) MaxLogsPerSpan(maxLogsPerSpan [int](/builtin#int)) TracerOption

MaxLogsPerSpan limits the number of Logs in a span (if set to a nonzero value). If a span has more logs than this value, logs are dropped as necessary (and replaced with a log describing how many were dropped).

About half of the MaxLogsPerSpan logs kept are the oldest logs, and about half are the newest logs.

#### func (TracerOptionsFactory) [MaxTagValueLength](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/tracer_options.go#L156) ¶

    func (TracerOptionsFactory) MaxTagValueLength(maxTagValueLength [int](/builtin#int)) TracerOption

MaxTagValueLength sets the limit on the max length of tag values.

#### func (TracerOptionsFactory) [Metrics](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/tracer_options.go#L38) ¶

    func (TracerOptionsFactory) Metrics(m *Metrics) TracerOption

Metrics creates a TracerOption that initializes Metrics on the tracer, which is used to emit statistics.

#### func (TracerOptionsFactory) [NoDebugFlagOnForcedSampling](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/tracer_options.go#L142) ¶

    func (TracerOptionsFactory) NoDebugFlagOnForcedSampling(noDebugFlagOnForcedSampling [bool](/builtin#bool)) TracerOption

NoDebugFlagOnForcedSampling turns off setting the debug flag in the trace context when the trace is force-started via sampling=1 span tag.

#### func (TracerOptionsFactory) [Observer](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/tracer_options.go#L122) ¶

    func (t TracerOptionsFactory) Observer(observer Observer) TracerOption

Observer registers an Observer.

#### func (TracerOptionsFactory) [PoolSpans](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/tracer_options.go#L86) ¶

    func (TracerOptionsFactory) PoolSpans(poolSpans [bool](/builtin#bool)) TracerOption

PoolSpans creates a TracerOption that tells the tracer whether it should use an object pool to minimize span allocations. This should be used with care, only if the service is not running any async tasks that can access parent spans after those spans have been finished.

#### func (TracerOptionsFactory) [RandomNumber](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/tracer_options.go#L76) ¶

    func (TracerOptionsFactory) RandomNumber(randomNumber func() [uint64](/builtin#uint64)) TracerOption

RandomNumber creates a TracerOption that gives the tracer a thread-safe random number generator function for generating trace IDs.

#### func (TracerOptionsFactory) [Tag](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/tracer_options.go#L184) ¶

    func (TracerOptionsFactory) Tag(key [string](/builtin#string), value interface{}) TracerOption

Tag adds a tracer-level tag that will be added to all spans.

#### func (TracerOptionsFactory) [TimeNow](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/tracer_options.go#L68) ¶

    func (TracerOptionsFactory) TimeNow(timeNow func() [time](/time).[Time](/time#Time)) TracerOption

TimeNow creates a TracerOption that gives the tracer a function used to generate timestamps for spans.

#### func (TracerOptionsFactory) [ZipkinSharedRPCSpan](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/tracer_options.go#L177) ¶

    func (TracerOptionsFactory) ZipkinSharedRPCSpan(zipkinSharedRPCSpan [bool](/builtin#bool)) TracerOption

ZipkinSharedRPCSpan enables a mode where server-side span shares the span ID from the client span from the incoming request, for compatibility with Zipkin's "one span per RPC" model.

#### type [Transport](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/transport.go#L24) ¶

    type Transport interface {
     // Append converts the span to the wire representation and adds it
     // to sender's internal buffer.  If the buffer exceeds its designated
     // size, the transport should call Flush() and return the number of spans
     // flushed, otherwise return 0. If error is returned, the returned number
     // of spans is treated as failed span, and reported to metrics accordingly.
     Append(span *Span) ([int](/builtin#int), [error](/builtin#error))
    
     // Flush submits the internal buffer to the remote server. It returns the
     // number of spans flushed. If error is returned, the returned number of
     // spans is treated as failed span, and reported to metrics accordingly.
     Flush() ([int](/builtin#int), [error](/builtin#error))
    
     [io](/io).[Closer](/io#Closer)
    }

Transport abstracts the method of sending spans out of process. Implementations are NOT required to be thread-safe; the RemoteReporter is expected to only call methods on the Transport from the same go-routine.

#### func [NewUDPTransport](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/transport_udp.go#L104) ¶

    func NewUDPTransport(hostPort [string](/builtin#string), maxPacketSize [int](/builtin#int)) (Transport, [error](/builtin#error))

NewUDPTransport creates a reporter that submits spans to jaeger-agent. TODO: (breaking change) move to transport/ package.

#### func [NewUDPTransportWithParams](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/transport_udp.go#L70) ¶

    func NewUDPTransportWithParams(params UDPTransportParams) (Transport, [error](/builtin#error))

NewUDPTransportWithParams creates a reporter that submits spans to jaeger-agent. TODO: (breaking change) move to transport/ package.

#### type [UDPTransportParams](https://github.com/jaegertracing/jaeger-client-go/blob/v2.30.0/transport_udp.go#L64) ¶

    type UDPTransportParams struct {
     [utils](/github.com/uber/jaeger-client-go/utils).[AgentClientUDPParams](/github.com/uber/jaeger-client-go/utils#AgentClientUDPParams)
    }

UDPTransportParams allows specifying options for initializing a UDPTransport. An instance of this struct should be passed to NewUDPTransportWithParams.
  *[↑]: Back to Top
  *[v]: View this template
  *[t]: Discuss this template
  *[e]: Edit this template
