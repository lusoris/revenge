# Prometheus Go Client

> Source: https://pkg.go.dev/github.com/prometheus/client_golang/prometheus
> Fetched: 2026-01-31T11:02:25.347019+00:00
> Content-Hash: 1d8e3c15cb921db4
> Type: html

---

### Overview ¶

  * A Basic Example
  * Metrics
  * Custom Collectors and constant Metrics
  * Advanced Uses of the Registry
  * HTTP Exposition
  * Pushing to the Pushgateway
  * Graphite Bridge
  * Other Means of Exposition



Package prometheus is the core instrumentation package. It provides metrics primitives to instrument code for monitoring. It also offers a registry for metrics. Sub-packages allow to expose the registered metrics via HTTP (package promhttp) or push them to a Pushgateway (package push). There is also a sub-package promauto, which provides metrics constructors with automatic registration. 

All exported functions and methods are safe to be used concurrently unless specified otherwise. 

#### A Basic Example ¶

As a starting point, a very basic usage example: 
    
    
    package main
    
    import (
    	"log"
    	"net/http"
    
    	"github.com/prometheus/client_golang/prometheus"
    	"github.com/prometheus/client_golang/prometheus/promhttp"
    )
    
    type metrics struct {
    	cpuTemp  prometheus.Gauge
    	hdFailures *prometheus.CounterVec
    }
    
    func NewMetrics(reg prometheus.Registerer) *metrics {
    	m := &metrics{
    		cpuTemp: prometheus.NewGauge(prometheus.GaugeOpts{
    			Name: "cpu_temperature_celsius",
    			Help: "Current temperature of the CPU.",
    		}),
    		hdFailures: prometheus.NewCounterVec(
    			prometheus.CounterOpts{
    				Name: "hd_errors_total",
    				Help: "Number of hard-disk errors.",
    			},
    			[]string{"device"},
    		),
    	}
    	reg.MustRegister(m.cpuTemp)
    	reg.MustRegister(m.hdFailures)
    	return m
    }
    
    func main() {
    	// Create a non-global registry.
    	reg := prometheus.NewRegistry()
    
    	// Create new metrics and register them using the custom registry.
    	m := NewMetrics(reg)
    	// Set values for the new created metrics.
    	m.cpuTemp.Set(65.3)
    	m.hdFailures.With(prometheus.Labels{"device":"/dev/sda"}).Inc()
    
    	// Expose metrics and custom registry via an HTTP server
    	// using the HandleFor function. "/metrics" is the usual endpoint for that.
    	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg}))
    	log.Fatal(http.ListenAndServe(":8080", nil))
    }
    

This is a complete program that exports two metrics, a Gauge and a Counter, the latter with a label attached to turn it into a (one-dimensional) vector. It register the metrics using a custom registry and exposes them via an HTTP server on the /metrics endpoint. 

#### Metrics ¶

The number of exported identifiers in this package might appear a bit overwhelming. However, in addition to the basic plumbing shown in the example above, you only need to understand the different metric types and their vector versions for basic usage. Furthermore, if you are not concerned with fine-grained control of when and how to register metrics with the registry, have a look at the promauto package, which will effectively allow you to ignore registration altogether in simple cases. 

Above, you have already touched the Counter and the Gauge. There are two more advanced metric types: the Summary and Histogram. A more thorough description of those four metric types can be found in the Prometheus docs: <https://prometheus.io/docs/concepts/metric_types/>

In addition to the fundamental metric types Gauge, Counter, Summary, and Histogram, a very important part of the Prometheus data model is the partitioning of samples along dimensions called labels, which results in metric vectors. The fundamental types are GaugeVec, CounterVec, SummaryVec, and HistogramVec. 

While only the fundamental metric types implement the Metric interface, both the metrics and their vector versions implement the Collector interface. A Collector manages the collection of a number of Metrics, but for convenience, a Metric can also “collect itself”. Note that Gauge, Counter, Summary, and Histogram are interfaces themselves while GaugeVec, CounterVec, SummaryVec, and HistogramVec are not. 

To create instances of Metrics and their vector versions, you need a suitable …Opts struct, i.e. GaugeOpts, CounterOpts, SummaryOpts, or HistogramOpts. 

#### Custom Collectors and constant Metrics ¶

While you could create your own implementations of Metric, most likely you will only ever implement the Collector interface on your own. At a first glance, a custom Collector seems handy to bundle Metrics for common registration (with the prime example of the different metric vectors above, which bundle all the metrics of the same name but with different labels). 

There is a more involved use case, too: If you already have metrics available, created outside of the Prometheus context, you don't need the interface of the various Metric types. You essentially want to mirror the existing numbers into Prometheus Metrics during collection. An own implementation of the Collector interface is perfect for that. You can create Metric instances “on the fly” using NewConstMetric, NewConstHistogram, and NewConstSummary (and their respective Must… versions). NewConstMetric is used for all metric types with just a float64 as their value: Counter, Gauge, and a special “type” called Untyped. Use the latter if you are not sure if the mirrored metric is a Counter or a Gauge. Creation of the Metric instance happens in the Collect method. The Describe method has to return separate Desc instances, representative of the “throw-away” metrics to be created later. NewDesc comes in handy to create those Desc instances. Alternatively, you could return no Desc at all, which will mark the Collector “unchecked”. No checks are performed at registration time, but metric consistency will still be ensured at scrape time, i.e. any inconsistencies will lead to scrape errors. Thus, with unchecked Collectors, the responsibility to not collect metrics that lead to inconsistencies in the total scrape result lies with the implementer of the Collector. While this is not a desirable state, it is sometimes necessary. The typical use case is a situation where the exact metrics to be returned by a Collector cannot be predicted at registration time, but the implementer has sufficient knowledge of the whole system to guarantee metric consistency. 

The Collector example illustrates the use case. You can also look at the source code of the processCollector (mirroring process metrics), the goCollector (mirroring Go metrics), or the expvarCollector (mirroring expvar metrics) as examples that are used in this package itself. 

If you just need to call a function to get a single float value to collect as a metric, GaugeFunc, CounterFunc, or UntypedFunc might be interesting shortcuts. 

#### Advanced Uses of the Registry ¶

While MustRegister is the by far most common way of registering a Collector, sometimes you might want to handle the errors the registration might cause. As suggested by the name, MustRegister panics if an error occurs. With the Register function, the error is returned and can be handled. 

An error is returned if the registered Collector is incompatible or inconsistent with already registered metrics. The registry aims for consistency of the collected metrics according to the Prometheus data model. Inconsistencies are ideally detected at registration time, not at collect time. The former will usually be detected at start-up time of a program, while the latter will only happen at scrape time, possibly not even on the first scrape if the inconsistency only becomes relevant later. That is the main reason why a Collector and a Metric have to describe themselves to the registry. 

So far, everything we did operated on the so-called default registry, as it can be found in the global DefaultRegisterer variable. With NewRegistry, you can create a custom registry, or you can even implement the Registerer or Gatherer interfaces yourself. The methods Register and Unregister work in the same way on a custom registry as the global functions Register and Unregister on the default registry. 

There are a number of uses for custom registries: You can use registries with special properties, see NewPedanticRegistry. You can avoid global state, as it is imposed by the DefaultRegisterer. You can use multiple registries at the same time to expose different metrics in different ways. You can use separate registries for testing purposes. 

Also note that the DefaultRegisterer comes registered with a Collector for Go runtime metrics (via NewGoCollector) and a Collector for process metrics (via NewProcessCollector). With a custom registry, you are in control and decide yourself about the Collectors to register. 

#### HTTP Exposition ¶

The Registry implements the Gatherer interface. The caller of the Gather method can then expose the gathered metrics in some way. Usually, the metrics are served via HTTP on the /metrics endpoint. That's happening in the example above. The tools to expose metrics via HTTP are in the promhttp sub-package. 

#### Pushing to the Pushgateway ¶

Function for pushing to the Pushgateway can be found in the push sub-package. 

#### Graphite Bridge ¶

Functions and examples to push metrics from a Gatherer to Graphite can be found in the graphite sub-package. 

#### Other Means of Exposition ¶

More ways of exposing metrics can easily be added by following the approaches of the existing implementations. 

### Index ¶

  * Constants
  * Variables
  * func BuildFQName(namespace, subsystem, name string) string
  * func DescribeByCollect(c Collector, descs chan<- *Desc)
  * func ExponentialBuckets(start, factor float64, count int) []float64
  * func ExponentialBucketsRange(minBucket, maxBucket float64, count int) []float64
  * func LinearBuckets(start, width float64, count int) []float64
  * func MakeLabelPairs(desc *Desc, labelValues []string) []*dto.LabelPair
  * func MustRegister(cs ...Collector)
  * func NewPidFileFn(pidFilePath string) func() (int, error)
  * func Register(c Collector) error
  * func Unregister(c Collector) bool
  * func WriteToTextfile(filename string, g Gatherer) error
  * type AlreadyRegisteredError
  *     * func (err AlreadyRegisteredError) Error() string
  * type Collector
  *     * func NewBuildInfoCollector() Collectordeprecated
    * func NewExpvarCollector(exports map[string]*Desc) Collectordeprecated
    * func NewGoCollector(opts ...func(o *internal.GoCollectorOptions)) Collectordeprecated
    * func NewProcessCollector(opts ProcessCollectorOpts) Collectordeprecated
    * func WrapCollectorWith(labels Labels, c Collector) Collector
    * func WrapCollectorWithPrefix(prefix string, c Collector) Collector
  * type CollectorFunc
  *     * func (f CollectorFunc) Collect(ch chan<- Metric)
    * func (f CollectorFunc) Describe(ch chan<- *Desc)
  * type ConstrainableLabels
  * type ConstrainedLabel
  * type ConstrainedLabels
  * type Counter
  *     * func NewCounter(opts CounterOpts) Counter
  * type CounterFunc
  *     * func NewCounterFunc(opts CounterOpts, function func() float64) CounterFunc
  * type CounterOpts
  * type CounterVec
  *     * func NewCounterVec(opts CounterOpts, labelNames []string) *CounterVec
  *     * func (v *CounterVec) CurryWith(labels Labels) (*CounterVec, error)
    * func (v *CounterVec) GetMetricWith(labels Labels) (Counter, error)
    * func (v *CounterVec) GetMetricWithLabelValues(lvs ...string) (Counter, error)
    * func (v *CounterVec) MustCurryWith(labels Labels) *CounterVec
    * func (v *CounterVec) With(labels Labels) Counter
    * func (v *CounterVec) WithLabelValues(lvs ...string) Counter
  * type CounterVecOpts
  * type Desc
  *     * func NewDesc(fqName, help string, variableLabels []string, constLabels Labels) *Desc
    * func NewInvalidDesc(err error) *Desc
  *     * func (d *Desc) String() string
  * type Exemplar
  * type ExemplarAdder
  * type ExemplarObserver
  * type Gatherer
  * type GathererFunc
  *     * func (gf GathererFunc) Gather() ([]*dto.MetricFamily, error)
  * type Gatherers
  *     * func (gs Gatherers) Gather() ([]*dto.MetricFamily, error)
  * type Gauge
  *     * func NewGauge(opts GaugeOpts) Gauge
  * type GaugeFunc
  *     * func NewGaugeFunc(opts GaugeOpts, function func() float64) GaugeFunc
  * type GaugeOpts
  * type GaugeVec
  *     * func NewGaugeVec(opts GaugeOpts, labelNames []string) *GaugeVec
  *     * func (v *GaugeVec) CurryWith(labels Labels) (*GaugeVec, error)
    * func (v *GaugeVec) GetMetricWith(labels Labels) (Gauge, error)
    * func (v *GaugeVec) GetMetricWithLabelValues(lvs ...string) (Gauge, error)
    * func (v *GaugeVec) MustCurryWith(labels Labels) *GaugeVec
    * func (v *GaugeVec) With(labels Labels) Gauge
    * func (v *GaugeVec) WithLabelValues(lvs ...string) Gauge
  * type GaugeVecOpts
  * type Histogram
  *     * func NewHistogram(opts HistogramOpts) Histogram
  * type HistogramOpts
  * type HistogramVec
  *     * func NewHistogramVec(opts HistogramOpts, labelNames []string) *HistogramVec
  *     * func (v *HistogramVec) CurryWith(labels Labels) (ObserverVec, error)
    * func (v *HistogramVec) GetMetricWith(labels Labels) (Observer, error)
    * func (v *HistogramVec) GetMetricWithLabelValues(lvs ...string) (Observer, error)
    * func (v *HistogramVec) MustCurryWith(labels Labels) ObserverVec
    * func (v *HistogramVec) With(labels Labels) Observer
    * func (v *HistogramVec) WithLabelValues(lvs ...string) Observer
  * type HistogramVecOpts
  * type LabelConstraint
  * type Labels
  * type Metric
  *     * func MustNewConstHistogram(desc *Desc, count uint64, sum float64, buckets map[float64]uint64, ...) Metric
    * func MustNewConstHistogramWithCreatedTimestamp(desc *Desc, count uint64, sum float64, buckets map[float64]uint64, ...) Metric
    * func MustNewConstMetric(desc *Desc, valueType ValueType, value float64, labelValues ...string) Metric
    * func MustNewConstMetricWithCreatedTimestamp(desc *Desc, valueType ValueType, value float64, ct time.Time, ...) Metric
    * func MustNewConstNativeHistogram(desc *Desc, count uint64, sum float64, ...) Metric
    * func MustNewConstSummary(desc *Desc, count uint64, sum float64, quantiles map[float64]float64, ...) Metric
    * func MustNewConstSummaryWithCreatedTimestamp(desc *Desc, count uint64, sum float64, quantiles map[float64]float64, ...) Metric
    * func MustNewMetricWithExemplars(m Metric, exemplars ...Exemplar) Metric
    * func NewConstHistogram(desc *Desc, count uint64, sum float64, buckets map[float64]uint64, ...) (Metric, error)
    * func NewConstHistogramWithCreatedTimestamp(desc *Desc, count uint64, sum float64, buckets map[float64]uint64, ...) (Metric, error)
    * func NewConstMetric(desc *Desc, valueType ValueType, value float64, labelValues ...string) (Metric, error)
    * func NewConstMetricWithCreatedTimestamp(desc *Desc, valueType ValueType, value float64, ct time.Time, ...) (Metric, error)
    * func NewConstNativeHistogram(desc *Desc, count uint64, sum float64, ...) (Metric, error)
    * func NewConstSummary(desc *Desc, count uint64, sum float64, quantiles map[float64]float64, ...) (Metric, error)
    * func NewConstSummaryWithCreatedTimestamp(desc *Desc, count uint64, sum float64, quantiles map[float64]float64, ...) (Metric, error)
    * func NewInvalidMetric(desc *Desc, err error) Metric
    * func NewMetricWithExemplars(m Metric, exemplars ...Exemplar) (Metric, error)
    * func NewMetricWithTimestamp(t time.Time, m Metric) Metric
  * type MetricVec
  *     * func NewMetricVec(desc *Desc, newMetric func(lvs ...string) Metric) *MetricVec
  *     * func (m *MetricVec) Collect(ch chan<- Metric)
    * func (m *MetricVec) CurryWith(labels Labels) (*MetricVec, error)
    * func (m *MetricVec) Delete(labels Labels) bool
    * func (m *MetricVec) DeleteLabelValues(lvs ...string) bool
    * func (m *MetricVec) DeletePartialMatch(labels Labels) int
    * func (m *MetricVec) Describe(ch chan<- *Desc)
    * func (m *MetricVec) GetMetricWith(labels Labels) (Metric, error)
    * func (m *MetricVec) GetMetricWithLabelValues(lvs ...string) (Metric, error)
    * func (m *MetricVec) Reset()
  * type MultiError
  *     * func (errs *MultiError) Append(err error)
    * func (errs MultiError) Error() string
    * func (errs MultiError) MaybeUnwrap() error
  * type MultiTRegistry
  *     * func NewMultiTRegistry(tGatherers ...TransactionalGatherer) *MultiTRegistry
  *     * func (r *MultiTRegistry) Gather() (mfs []*dto.MetricFamily, done func(), err error)
  * type Observer
  * type ObserverFunc
  *     * func (f ObserverFunc) Observe(value float64)
  * type ObserverVec
  * type Opts
  * type ProcessCollectorOpts
  * type Registerer
  *     * func WrapRegistererWith(labels Labels, reg Registerer) Registerer
    * func WrapRegistererWithPrefix(prefix string, reg Registerer) Registerer
  * type Registry
  *     * func NewPedanticRegistry() *Registry
    * func NewRegistry() *Registry
  *     * func (r *Registry) Collect(ch chan<- Metric)
    * func (r *Registry) Describe(ch chan<- *Desc)
    * func (r *Registry) Gather() ([]*dto.MetricFamily, error)
    * func (r *Registry) MustRegister(cs ...Collector)
    * func (r *Registry) Register(c Collector) error
    * func (r *Registry) Unregister(c Collector) bool
  * type Summary
  *     * func NewSummary(opts SummaryOpts) Summary
  * type SummaryOpts
  * type SummaryVec
  *     * func NewSummaryVec(opts SummaryOpts, labelNames []string) *SummaryVec
  *     * func (v *SummaryVec) CurryWith(labels Labels) (ObserverVec, error)
    * func (v *SummaryVec) GetMetricWith(labels Labels) (Observer, error)
    * func (v *SummaryVec) GetMetricWithLabelValues(lvs ...string) (Observer, error)
    * func (v *SummaryVec) MustCurryWith(labels Labels) ObserverVec
    * func (v *SummaryVec) With(labels Labels) Observer
    * func (v *SummaryVec) WithLabelValues(lvs ...string) Observer
  * type SummaryVecOpts
  * type Timer
  *     * func NewTimer(o Observer) *Timer
  *     * func (t *Timer) ObserveDuration() time.Duration
    * func (t *Timer) ObserveDurationWithExemplar(exemplar Labels) time.Duration
  * type TransactionalGatherer
  *     * func ToTransactionalGatherer(g Gatherer) TransactionalGatherer
  * type UnconstrainedLabels
  * type UntypedFunc
  *     * func NewUntypedFunc(opts UntypedOpts, function func() float64) UntypedFunc
  * type UntypedOpts
  * type ValueType
  *     * func (v ValueType) ToDTO() *dto.MetricType



### Examples ¶

  * AlreadyRegisteredError
  * Collector
  * CollectorFunc
  * CounterVec
  * Gatherers
  * Gauge
  * GaugeFunc (ConstLabels)
  * GaugeFunc (Simple)
  * GaugeVec
  * Histogram
  * MetricVec
  * NewConstHistogram
  * NewConstHistogram (WithExemplar)
  * NewConstHistogramWithCreatedTimestamp
  * NewConstMetricWithCreatedTimestamp
  * NewConstSummary
  * NewConstSummaryWithCreatedTimestamp
  * NewExpvarCollector
  * NewMetricWithTimestamp
  * Register
  * Registry (Grouping)
  * Summary
  * SummaryVec
  * Timer
  * Timer (Complex)
  * Timer (Gauge)
  * WrapCollectorWith



### Constants ¶

[View Source](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/summary.go#L72)
    
    
    const (
    	// DefMaxAge is the default duration for which observations stay
    	// relevant.
    	DefMaxAge [time](/time).[Duration](/time#Duration) = 10 * [time](/time).[Minute](/time#Minute)
    	// DefAgeBuckets is the default number of buckets used to calculate the
    	// age of observations.
    	DefAgeBuckets = 5
    	// DefBufCap is the standard buffer size for collecting Summary observations.
    	DefBufCap = 500
    )

Default values for SummaryOpts. 

[View Source](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/histogram.go#L278)
    
    
    const DefNativeHistogramZeroThreshold = 2.938735877055719e-39

DefNativeHistogramZeroThreshold is the default value for NativeHistogramZeroThreshold in the HistogramOpts. 

The value is 2^-128 (or 0.5*2^-127 in the actual IEEE 754 representation), which is a bucket boundary at all possible resolutions. 

[View Source](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/value.go#L240)
    
    
    const ExemplarMaxRunes = 128

ExemplarMaxRunes is the max total number of runes allowed in exemplar labels. 

[View Source](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/histogram.go#L283)
    
    
    const NativeHistogramZeroThresholdZero = -1

NativeHistogramZeroThresholdZero can be used as NativeHistogramZeroThreshold in the HistogramOpts to create a zero bucket of width zero, i.e. a zero bucket that only receives observations of precisely zero. 

### Variables ¶

[View Source](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/registry.go#L54)
    
    
    var (
    	DefaultRegisterer Registerer = defaultRegistry
    	DefaultGatherer   Gatherer   = defaultRegistry
    )

DefaultRegisterer and DefaultGatherer are the implementations of the Registerer and Gatherer interface a number of convenience functions in this package act on. Initially, both variables point to the same Registry, which has a process collector (currently on Linux only, see NewProcessCollector) and a Go collector (see NewGoCollector, in particular the note about stop-the-world implication with Go versions older than 1.9) already registered. This approach to keep default instances as global state mirrors the approach of other packages in the Go standard library. Note that there are caveats. Change the variables with caution and only if you understand the consequences. Users who want to avoid global state altogether should not use the convenience functions and act on custom instances instead. 

[View Source](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/value.go#L42)
    
    
    var (
    	CounterMetricTypePtr = func() *[dto](/github.com/prometheus/client_model/go).[MetricType](/github.com/prometheus/client_model/go#MetricType) { d := [dto](/github.com/prometheus/client_model/go).[MetricType_COUNTER](/github.com/prometheus/client_model/go#MetricType_COUNTER); return &d }()
    	GaugeMetricTypePtr   = func() *[dto](/github.com/prometheus/client_model/go).[MetricType](/github.com/prometheus/client_model/go#MetricType) { d := [dto](/github.com/prometheus/client_model/go).[MetricType_GAUGE](/github.com/prometheus/client_model/go#MetricType_GAUGE); return &d }()
    	UntypedMetricTypePtr = func() *[dto](/github.com/prometheus/client_model/go).[MetricType](/github.com/prometheus/client_model/go#MetricType) { d := [dto](/github.com/prometheus/client_model/go).[MetricType_UNTYPED](/github.com/prometheus/client_model/go#MetricType_UNTYPED); return &d }()
    )

[View Source](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/histogram.go#L271)
    
    
    var DefBuckets = [][float64](/builtin#float64){.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10}

DefBuckets are the default Histogram buckets. The default buckets are tailored to broadly measure the response time (in seconds) of a network service. Most likely, however, you will be required to define buckets customized to your use case. 

[View Source](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/vnext.go#L23)
    
    
    var V2 = v2{}

V2 is a struct that can be referenced to access experimental API that might be present in v2 of client golang someday. It offers extended functionality of v1 with slightly changed API. It is acceptable to use some pieces from v1 and e.g `prometheus.NewGauge` and some from v2 e.g. `prometheus.V2.NewDesc` in the same codebase. 

### Functions ¶

####  func [BuildFQName](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/metric.go#L107) ¶
    
    
    func BuildFQName(namespace, subsystem, name [string](/builtin#string)) [string](/builtin#string)

BuildFQName joins the given three name components by "_". Empty name components are ignored. If the name parameter itself is empty, an empty string is returned, no matter what. Metric implementations included in this library use this function internally to generate the fully-qualified metric name from the name component in their Opts. Users of the library will only need this function if they implement their own Metric or instantiate a Desc (with NewDesc) directly. 

####  func [DescribeByCollect](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/collector.go#L87) ¶ added in v0.9.0
    
    
    func DescribeByCollect(c Collector, descs chan<- *Desc)

DescribeByCollect is a helper to implement the Describe method of a custom Collector. It collects the metrics from the provided Collector and sends their descriptors to the provided channel. 

If a Collector collects the same metrics throughout its lifetime, its Describe method can simply be implemented as: 
    
    
    func (c customCollector) Describe(ch chan<- *Desc) {
    	DescribeByCollect(c, ch)
    }
    

However, this will not work if the metrics collected change dynamically over the lifetime of the Collector in a way that their combined set of descriptors changes as well. The shortcut implementation will then violate the contract of the Describe method. If a Collector sometimes collects no metrics at all (for example vectors like CounterVec, GaugeVec, etc., which only collect metrics after a metric with a fully specified label set has been accessed), it might even get registered as an unchecked Collector (cf. the Register method of the Registerer interface). Hence, only use this shortcut implementation of Describe if you are certain to fulfill the contract. 

The Collector example demonstrates a use of DescribeByCollect. 

####  func [ExponentialBuckets](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/histogram.go#L315) ¶
    
    
    func ExponentialBuckets(start, factor [float64](/builtin#float64), count [int](/builtin#int)) [][float64](/builtin#float64)

ExponentialBuckets creates 'count' regular buckets, where the lowest bucket has an upper bound of 'start' and each following bucket's upper bound is 'factor' times the previous bucket's upper bound. The final +Inf bucket is not counted and not included in the returned slice. The returned slice is meant to be used for the Buckets field of HistogramOpts. 

The function panics if 'count' is 0 or negative, if 'start' is 0 or negative, or if 'factor' is less than or equal 1. 

####  func [ExponentialBucketsRange](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/histogram.go#L339) ¶ added in v0.12.1
    
    
    func ExponentialBucketsRange(minBucket, maxBucket [float64](/builtin#float64), count [int](/builtin#int)) [][float64](/builtin#float64)

ExponentialBucketsRange creates 'count' buckets, where the lowest bucket is 'min' and the highest bucket is 'max'. The final +Inf bucket is not counted and not included in the returned slice. The returned slice is meant to be used for the Buckets field of HistogramOpts. 

The function panics if 'count' is 0 or negative, if 'min' is 0 or negative. 

####  func [LinearBuckets](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/histogram.go#L295) ¶
    
    
    func LinearBuckets(start, width [float64](/builtin#float64), count [int](/builtin#int)) [][float64](/builtin#float64)

LinearBuckets creates 'count' regular buckets, each 'width' wide, where the lowest bucket has an upper bound of 'start'. The final +Inf bucket is not counted and not included in the returned slice. The returned slice is meant to be used for the Buckets field of HistogramOpts. 

The function panics if 'count' is zero or negative. 

####  func [MakeLabelPairs](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/value.go#L217) ¶ added in v0.12.1
    
    
    func MakeLabelPairs(desc *Desc, labelValues [][string](/builtin#string)) []*[dto](/github.com/prometheus/client_model/go).[LabelPair](/github.com/prometheus/client_model/go#LabelPair)

MakeLabelPairs is a helper function to create protobuf LabelPairs from the variable and constant labels in the provided Desc. The values for the variable labels are defined by the labelValues slice, which must be in the same order as the corresponding variable labels in the Desc. 

This function is only needed for custom Metric implementations. See MetricVec example. 

####  func [MustRegister](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/registry.go#L176) ¶
    
    
    func MustRegister(cs ...Collector)

MustRegister registers the provided Collectors with the DefaultRegisterer and panics if any error occurs. 

MustRegister is a shortcut for DefaultRegisterer.MustRegister(cs...). See there for more details. 

####  func [NewPidFileFn](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/process_collector.go#L167) ¶ added in v0.12.1
    
    
    func NewPidFileFn(pidFilePath [string](/builtin#string)) func() ([int](/builtin#int), [error](/builtin#error))

NewPidFileFn returns a function that retrieves a pid from the specified file. It is meant to be used for the PidFn field in ProcessCollectorOpts. 

####  func [Register](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/registry.go#L167) ¶
    
    
    func Register(c Collector) [error](/builtin#error)

Register registers the provided Collector with the DefaultRegisterer. 

Register is a shortcut for DefaultRegisterer.Register(c). See there for more details. 

Example ¶
    
    
    package main
    
    import (
    	"fmt"
    	"net/http"
    
    	"github.com/prometheus/client_golang/prometheus"
    	"github.com/prometheus/client_golang/prometheus/promhttp"
    )
    
    func main() {
    	// Imagine you have a worker pool and want to count the tasks completed.
    	taskCounter := prometheus.NewCounter(prometheus.CounterOpts{
    		Subsystem: "worker_pool",
    		Name:      "completed_tasks_total",
    		Help:      "Total number of tasks completed.",
    	})
    	// This will register fine.
    	if err := prometheus.Register(taskCounter); err != nil {
    		fmt.Println(err)
    	} else {
    		fmt.Println("taskCounter registered.")
    	}
    	// Don't forget to tell the HTTP server about the Prometheus handler.
    	// (In a real program, you still need to start the HTTP server...)
    	http.Handle("/metrics", promhttp.Handler())
    
    	// Now you can start workers and give every one of them a pointer to
    	// taskCounter and let it increment it whenever it completes a task.
    	taskCounter.Inc() // This has to happen somewhere in the worker code.
    
    	// But wait, you want to see how individual workers perform. So you need
    	// a vector of counters, with one element for each worker.
    	taskCounterVec := prometheus.NewCounterVec(
    		prometheus.CounterOpts{
    			Subsystem: "worker_pool",
    			Name:      "completed_tasks_total",
    			Help:      "Total number of tasks completed.",
    		},
    		[]string{"worker_id"},
    	)
    
    	// Registering will fail because we already have a metric of that name.
    	if err := prometheus.Register(taskCounterVec); err != nil {
    		fmt.Println("taskCounterVec not registered:", err)
    	} else {
    		fmt.Println("taskCounterVec registered.")
    	}
    
    	// To fix, first unregister the old taskCounter.
    	if prometheus.Unregister(taskCounter) {
    		fmt.Println("taskCounter unregistered.")
    	}
    
    	// Try registering taskCounterVec again.
    	if err := prometheus.Register(taskCounterVec); err != nil {
    		fmt.Println("taskCounterVec not registered:", err)
    	} else {
    		fmt.Println("taskCounterVec registered.")
    	}
    	// Bummer! Still doesn't work.
    
    	// Prometheus will not allow you to ever export metrics with
    	// inconsistent help strings or label names. After unregistering, the
    	// unregistered metrics will cease to show up in the /metrics HTTP
    	// response, but the registry still remembers that those metrics had
    	// been exported before. For this example, we will now choose a
    	// different name. (In a real program, you would obviously not export
    	// the obsolete metric in the first place.)
    	taskCounterVec = prometheus.NewCounterVec(
    		prometheus.CounterOpts{
    			Subsystem: "worker_pool",
    			Name:      "completed_tasks_by_id",
    			Help:      "Total number of tasks completed.",
    		},
    		[]string{"worker_id"},
    	)
    	if err := prometheus.Register(taskCounterVec); err != nil {
    		fmt.Println("taskCounterVec not registered:", err)
    	} else {
    		fmt.Println("taskCounterVec registered.")
    	}
    	// Finally it worked!
    
    	// The workers have to tell taskCounterVec their id to increment the
    	// right element in the metric vector.
    	taskCounterVec.WithLabelValues("42").Inc() // Code from worker 42.
    
    	// Each worker could also keep a reference to their own counter element
    	// around. Pick the counter at initialization time of the worker.
    	myCounter := taskCounterVec.WithLabelValues("42") // From worker 42 initialization code.
    	myCounter.Inc()                                   // Somewhere in the code of that worker.
    
    	// Note that something like WithLabelValues("42", "spurious arg") would
    	// panic (because you have provided too many label values). If you want
    	// to get an error instead, use GetMetricWithLabelValues(...) instead.
    	notMyCounter, err := taskCounterVec.GetMetricWithLabelValues("42", "spurious arg")
    	if err != nil {
    		fmt.Println("Worker initialization failed:", err)
    	}
    	if notMyCounter == nil {
    		fmt.Println("notMyCounter is nil.")
    	}
    
    	// A different (and somewhat tricky) approach is to use
    	// ConstLabels. ConstLabels are pairs of label names and label values
    	// that never change. Each worker creates and registers an own Counter
    	// instance where the only difference is in the value of the
    	// ConstLabels. Those Counters can all be registered because the
    	// different ConstLabel values guarantee that each worker will increment
    	// a different Counter metric.
    	counterOpts := prometheus.CounterOpts{
    		Subsystem:   "worker_pool",
    		Name:        "completed_tasks",
    		Help:        "Total number of tasks completed.",
    		ConstLabels: prometheus.Labels{"worker_id": "42"},
    	}
    	taskCounterForWorker42 := prometheus.NewCounter(counterOpts)
    	if err := prometheus.Register(taskCounterForWorker42); err != nil {
    		fmt.Println("taskCounterVForWorker42 not registered:", err)
    	} else {
    		fmt.Println("taskCounterForWorker42 registered.")
    	}
    	// Obviously, in real code, taskCounterForWorker42 would be a member
    	// variable of a worker struct, and the "42" would be retrieved with a
    	// GetId() method or something. The Counter would be created and
    	// registered in the initialization code of the worker.
    
    	// For the creation of the next Counter, we can recycle
    	// counterOpts. Just change the ConstLabels.
    	counterOpts.ConstLabels = prometheus.Labels{"worker_id": "2001"}
    	taskCounterForWorker2001 := prometheus.NewCounter(counterOpts)
    	if err := prometheus.Register(taskCounterForWorker2001); err != nil {
    		fmt.Println("taskCounterVForWorker2001 not registered:", err)
    	} else {
    		fmt.Println("taskCounterForWorker2001 registered.")
    	}
    
    	taskCounterForWorker2001.Inc()
    	taskCounterForWorker42.Inc()
    	taskCounterForWorker2001.Inc()
    
    	// Yet another approach would be to turn the workers themselves into
    	// Collectors and register them. See the Collector example for details.
    
    }
    
    
    
    Output:
    
    taskCounter registered.
    taskCounterVec not registered: a previously registered descriptor with the same fully-qualified name as Desc{fqName: "worker_pool_completed_tasks_total", help: "Total number of tasks completed.", constLabels: {}, variableLabels: {worker_id}} has different label names or a different help string
    taskCounter unregistered.
    taskCounterVec not registered: a previously registered descriptor with the same fully-qualified name as Desc{fqName: "worker_pool_completed_tasks_total", help: "Total number of tasks completed.", constLabels: {}, variableLabels: {worker_id}} has different label names or a different help string
    taskCounterVec registered.
    Worker initialization failed: inconsistent label cardinality: expected 1 label values but got 2 in []string{"42", "spurious arg"}
    notMyCounter is nil.
    taskCounterForWorker42 registered.
    taskCounterForWorker2001 registered.
    

Share Format Run

####  func [Unregister](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/registry.go#L185) ¶
    
    
    func Unregister(c Collector) [bool](/builtin#bool)

Unregister removes the registration of the provided Collector from the DefaultRegisterer. 

Unregister is a shortcut for DefaultRegisterer.Unregister(c). See there for more details. 

####  func [WriteToTextfile](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/registry.go#L593) ¶ added in v0.9.1
    
    
    func WriteToTextfile(filename [string](/builtin#string), g Gatherer) [error](/builtin#error)

WriteToTextfile calls Gather on the provided Gatherer, encodes the result in the Prometheus text format, and writes it to a temporary file. Upon success, the temporary file is renamed to the provided filename. 

This is intended for use with the textfile collector of the node exporter. Note that the node exporter expects the filename to be suffixed with ".prom". 

### Types ¶

####  type [AlreadyRegisteredError](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/registry.go#L205) ¶
    
    
    type AlreadyRegisteredError struct {
    	ExistingCollector, NewCollector Collector
    }

AlreadyRegisteredError is returned by the Register method if the Collector to be registered has already been registered before, or a different Collector that collects the same metrics has been registered before. Registration fails in that case, but you can detect from the kind of error what has happened. The error contains fields for the existing Collector and the (rejected) new Collector that equals the existing one. This can be used to find out if an equal Collector has been registered before and switch over to using the old one, as demonstrated in the example. 

Example ¶
    
    
    package main
    
    import (
    	"errors"
    
    	"github.com/prometheus/client_golang/prometheus"
    )
    
    func main() {
    	reqCounter := prometheus.NewCounter(prometheus.CounterOpts{
    		Name: "requests_total",
    		Help: "The total number of requests served.",
    	})
    	if err := prometheus.Register(reqCounter); err != nil {
    		are := &prometheus.AlreadyRegisteredError{}
    		if errors.As(err, are) {
    			// A counter for that metric has been registered before.
    			// Use the old counter from now on.
    			reqCounter = are.ExistingCollector.(prometheus.Counter)
    		} else {
    			// Something else went wrong!
    			panic(err)
    		}
    	}
    	reqCounter.Inc()
    }
    

Share Format Run

####  func (AlreadyRegisteredError) [Error](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/registry.go#L209) ¶
    
    
    func (err AlreadyRegisteredError) Error() [string](/builtin#string)

####  type [Collector](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/collector.go#L27) ¶
    
    
    type Collector interface {
    	// Describe sends the super-set of all possible descriptors of metrics
    	// collected by this Collector to the provided channel and returns once
    	// the last descriptor has been sent. The sent descriptors fulfill the
    	// consistency and uniqueness requirements described in the Desc
    	// documentation.
    	//
    	// It is valid if one and the same Collector sends duplicate
    	// descriptors. Those duplicates are simply ignored. However, two
    	// different Collectors must not send duplicate descriptors.
    	//
    	// Sending no descriptor at all marks the Collector as “unchecked”,
    	// i.e. no checks will be performed at registration time, and the
    	// Collector may yield any Metric it sees fit in its Collect method.
    	//
    	// This method idempotently sends the same descriptors throughout the
    	// lifetime of the Collector. It may be called concurrently and
    	// therefore must be implemented in a concurrency safe way.
    	//
    	// If a Collector encounters an error while executing this method, it
    	// must send an invalid descriptor (created with NewInvalidDesc) to
    	// signal the error to the registry.
    	Describe(chan<- *Desc)
    	// Collect is called by the Prometheus registry when collecting
    	// metrics. The implementation sends each collected metric via the
    	// provided channel and returns once the last metric has been sent. The
    	// descriptor of each sent metric is one of those returned by Describe
    	// (unless the Collector is unchecked, see above). Returned metrics that
    	// share the same descriptor must differ in their variable label
    	// values.
    	//
    	// This method may be called concurrently and must therefore be
    	// implemented in a concurrency safe way. Blocking occurs at the expense
    	// of total performance of rendering all registered metrics. Ideally,
    	// Collector implementations support concurrent readers.
    	Collect(chan<- Metric)
    }

Collector is the interface implemented by anything that can be used by Prometheus to collect metrics. A Collector has to be registered for collection. See Registerer.Register. 

The stock metrics provided by this package (Gauge, Counter, Summary, Histogram, Untyped) are also Collectors (which only ever collect one metric, namely itself). An implementer of Collector may, however, collect multiple metrics in a coordinated fashion and/or create metrics on the fly. Examples for collectors already implemented in this library are the metric vectors (i.e. collection of multiple instances of the same Metric but with different label values) like GaugeVec or SummaryVec, and the ExpvarCollector. 

Example ¶
    
    
    package main
    
    import (
    	"log"
    	"net/http"
    
    	"github.com/prometheus/client_golang/prometheus"
    	"github.com/prometheus/client_golang/prometheus/promhttp"
    )
    
    // ClusterManager is an example for a system that might have been built without
    // Prometheus in mind. It models a central manager of jobs running in a
    // cluster. Thus, we implement a custom Collector called
    // ClusterManagerCollector, which collects information from a ClusterManager
    // using its provided methods and turns them into Prometheus Metrics for
    // collection.
    //
    // An additional challenge is that multiple instances of the ClusterManager are
    // run within the same binary, each in charge of a different zone. We need to
    // make use of wrapping Registerers to be able to register each
    // ClusterManagerCollector instance with Prometheus.
    type ClusterManager struct {
    	Zone string
    	// Contains many more fields not listed in this example.
    }
    
    // ReallyExpensiveAssessmentOfTheSystemState is a mock for the data gathering a
    // real cluster manager would have to do. Since it may actually be really
    // expensive, it must only be called once per collection. This implementation,
    // obviously, only returns some made-up data.
    func (c *ClusterManager) ReallyExpensiveAssessmentOfTheSystemState() (
    	oomCountByHost map[string]int, ramUsageByHost map[string]float64,
    ) {
    	// Just example fake data.
    	oomCountByHost = map[string]int{
    		"foo.example.org": 42,
    		"bar.example.org": 2001,
    	}
    	ramUsageByHost = map[string]float64{
    		"foo.example.org": 6.023e23,
    		"bar.example.org": 3.14,
    	}
    	return
    }
    
    // ClusterManagerCollector implements the Collector interface.
    type ClusterManagerCollector struct {
    	ClusterManager *ClusterManager
    }
    
    // Descriptors used by the ClusterManagerCollector below.
    var (
    	oomCountDesc = prometheus.NewDesc(
    		"clustermanager_oom_crashes_total",
    		"Number of OOM crashes.",
    		[]string{"host"}, nil,
    	)
    	ramUsageDesc = prometheus.NewDesc(
    		"clustermanager_ram_usage_bytes",
    		"RAM usage as reported to the cluster manager.",
    		[]string{"host"}, nil,
    	)
    )
    
    // Describe is implemented with DescribeByCollect. That's possible because the
    // Collect method will always return the same two metrics with the same two
    // descriptors.
    func (cc ClusterManagerCollector) Describe(ch chan<- *prometheus.Desc) {
    	prometheus.DescribeByCollect(cc, ch)
    }
    
    // Collect first triggers the ReallyExpensiveAssessmentOfTheSystemState. Then it
    // creates constant metrics for each host on the fly based on the returned data.
    //
    // Note that Collect could be called concurrently, so we depend on
    // ReallyExpensiveAssessmentOfTheSystemState to be concurrency-safe.
    func (cc ClusterManagerCollector) Collect(ch chan<- prometheus.Metric) {
    	oomCountByHost, ramUsageByHost := cc.ClusterManager.ReallyExpensiveAssessmentOfTheSystemState()
    	for host, oomCount := range oomCountByHost {
    		ch <- prometheus.MustNewConstMetric(
    			oomCountDesc,
    			prometheus.CounterValue,
    			float64(oomCount),
    			host,
    		)
    	}
    	for host, ramUsage := range ramUsageByHost {
    		ch <- prometheus.MustNewConstMetric(
    			ramUsageDesc,
    			prometheus.GaugeValue,
    			ramUsage,
    			host,
    		)
    	}
    }
    
    // NewClusterManager first creates a Prometheus-ignorant ClusterManager
    // instance. Then, it creates a ClusterManagerCollector for the just created
    // ClusterManager. Finally, it registers the ClusterManagerCollector with a
    // wrapping Registerer that adds the zone as a label. In this way, the metrics
    // collected by different ClusterManagerCollectors do not collide.
    func NewClusterManager(zone string, reg prometheus.Registerer) *ClusterManager {
    	c := &ClusterManager{
    		Zone: zone,
    	}
    	cc := ClusterManagerCollector{ClusterManager: c}
    	prometheus.WrapRegistererWith(prometheus.Labels{"zone": zone}, reg).MustRegister(cc)
    	return c
    }
    
    func main() {
    	// Since we are dealing with custom Collector implementations, it might
    	// be a good idea to try it out with a pedantic registry.
    	reg := prometheus.NewPedanticRegistry()
    
    	// Construct cluster managers. In real code, we would assign them to
    	// variables to then do something with them.
    	NewClusterManager("db", reg)
    	NewClusterManager("ca", reg)
    
    	// Add the standard process and Go metrics to the custom registry.
    	reg.MustRegister(
    		prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}),
    		prometheus.NewGoCollector(),
    	)
    
    	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
    	log.Fatal(http.ListenAndServe(":8080", nil))
    }
    

Share Format Run

####  func [NewBuildInfoCollector](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/build_info_collector.go#L22) deprecated added in v0.9.4
    
    
    func NewBuildInfoCollector() Collector

NewBuildInfoCollector is the obsolete version of collectors.NewBuildInfoCollector. See there for documentation. 

Deprecated: Use collectors.NewBuildInfoCollector instead. 

####  func [NewExpvarCollector](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/expvar_collector.go#L29) deprecated
    
    
    func NewExpvarCollector(exports map[[string](/builtin#string)]*Desc) Collector

NewExpvarCollector is the obsolete version of collectors.NewExpvarCollector. See there for documentation. 

Deprecated: Use collectors.NewExpvarCollector instead. 

Example ¶
    
    
    expvarCollector := prometheus.NewExpvarCollector(map[string]*prometheus.Desc{
    	"memstats": prometheus.NewDesc(
    		"expvar_memstats",
    		"All numeric memstats as one metric family. Not a good role-model, actually... ;-)",
    		[]string{"type"}, nil,
    	),
    	"lone-int": prometheus.NewDesc(
    		"expvar_lone_int",
    		"Just an expvar int as an example.",
    		nil, nil,
    	),
    	"http-request-map": prometheus.NewDesc(
    		"expvar_http_request_total",
    		"How many http requests processed, partitioned by status code and http method.",
    		[]string{"code", "method"}, nil,
    	),
    })
    prometheus.MustRegister(expvarCollector)
    
    // The Prometheus part is done here. But to show that this example is
    // doing anything, we have to manually export something via expvar.  In
    // real-life use-cases, some library would already have exported via
    // expvar what we want to re-export as Prometheus metrics.
    expvar.NewInt("lone-int").Set(42)
    expvarMap := expvar.NewMap("http-request-map")
    var (
    	expvarMap1, expvarMap2                             expvar.Map
    	expvarInt11, expvarInt12, expvarInt21, expvarInt22 expvar.Int
    )
    expvarMap1.Init()
    expvarMap2.Init()
    expvarInt11.Set(3)
    expvarInt12.Set(13)
    expvarInt21.Set(11)
    expvarInt22.Set(212)
    expvarMap1.Set("POST", &expvarInt11)
    expvarMap1.Set("GET", &expvarInt12)
    expvarMap2.Set("POST", &expvarInt21)
    expvarMap2.Set("GET", &expvarInt22)
    expvarMap.Set("404", &expvarMap1)
    expvarMap.Set("200", &expvarMap2)
    // Results in the following expvar map:
    // "http-request-count": {"200": {"POST": 11, "GET": 212}, "404": {"POST": 3, "GET": 13}}
    
    // Let's see what the scrape would yield, but exclude the memstats metrics.
    metricStrings := []string{}
    metric := dto.Metric{}
    metricChan := make(chan prometheus.Metric)
    go func() {
    	expvarCollector.Collect(metricChan)
    	close(metricChan)
    }()
    for m := range metricChan {
    	if !strings.Contains(m.Desc().String(), "expvar_memstats") {
    		metric.Reset()
    		m.Write(&metric)
    		metricStrings = append(metricStrings, toNormalizedJSON(&metric))
    	}
    }
    sort.Strings(metricStrings)
    for _, s := range metricStrings {
    	fmt.Println(s)
    }
    
    
    
    Output:
    
    {"label":[{"name":"code","value":"200"},{"name":"method","value":"GET"}],"untyped":{"value":212}}
    {"label":[{"name":"code","value":"200"},{"name":"method","value":"POST"}],"untyped":{"value":11}}
    {"label":[{"name":"code","value":"404"},{"name":"method","value":"GET"}],"untyped":{"value":13}}
    {"label":[{"name":"code","value":"404"},{"name":"method","value":"POST"}],"untyped":{"value":3}}
    {"untyped":{"value":42}}
    

####  func [NewGoCollector](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/go_collector_latest.go#L167) deprecated
    
    
    func NewGoCollector(opts ...func(o *[internal](/github.com/prometheus/client_golang@v1.23.2/prometheus/internal).[GoCollectorOptions](/github.com/prometheus/client_golang@v1.23.2/prometheus/internal#GoCollectorOptions))) Collector

NewGoCollector is the obsolete version of collectors.NewGoCollector. See there for documentation. 

Deprecated: Use collectors.NewGoCollector instead. 

####  func [NewProcessCollector](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/process_collector.go#L62) deprecated
    
    
    func NewProcessCollector(opts ProcessCollectorOpts) Collector

NewProcessCollector is the obsolete version of collectors.NewProcessCollector. See there for documentation. 

Deprecated: Use collectors.NewProcessCollector instead. 

####  func [WrapCollectorWith](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/wrap.go#L97) ¶ added in v1.23.0
    
    
    func WrapCollectorWith(labels Labels, c Collector) Collector

WrapCollectorWith returns a Collector wrapping the provided Collector. The wrapped Collector will add the provided Labels to all Metrics it collects (as ConstLabels). The Metrics collected by the unmodified Collector must not duplicate any of those labels. 

WrapCollectorWith can be useful to work with multiple instances of a third party library that does not expose enough flexibility on the lifecycle of its registered metrics. For example, let's say you have a foo.New(reg Registerer) constructor that registers metrics but never unregisters them, and you want to create multiple instances of foo.Foo with different labels. The way to achieve that, is to create a new Registry, pass it to foo.New, then use WrapCollectorWith to wrap that Registry with the desired labels and register that as a collector in your main Registry. Then you can un-register the wrapped collector effectively un-registering the metrics registered by foo.New. 

Example ¶

Using WrapCollectorWith to un-register metrics registered by a third party lib. newThirdPartyLibFoo illustrates a constructor from a third-party lib that does not expose any way to un-register metrics. 
    
    
    reg := prometheus.NewRegistry()
    
    // We want to create two instances of thirdPartyLibFoo, each one wrapped with
    // its "instance" label.
    firstReg := prometheus.NewRegistry()
    _ = newThirdPartyLibFoo(firstReg)
    firstCollector := prometheus.WrapCollectorWith(prometheus.Labels{"instance": "first"}, firstReg)
    reg.MustRegister(firstCollector)
    
    secondReg := prometheus.NewRegistry()
    _ = newThirdPartyLibFoo(secondReg)
    secondCollector := prometheus.WrapCollectorWith(prometheus.Labels{"instance": "second"}, secondReg)
    reg.MustRegister(secondCollector)
    
    // So far we have illustrated that we can create two instances of thirdPartyLibFoo,
    // wrapping each one's metrics with some const label.
    // This is something we could've achieved by doing:
    // newThirdPartyLibFoo(prometheus.WrapRegistererWith(prometheus.Labels{"instance": "first"}, reg))
    metricFamilies, err := reg.Gather()
    if err != nil {
    	panic("unexpected behavior of registry")
    }
    fmt.Println("Both instances:")
    fmt.Println(toNormalizedJSON(sanitizeMetricFamily(metricFamilies[0])))
    
    // Now we want to unregister first Foo's metrics, and then register them again.
    // This is not possible by passing a wrapped Registerer to newThirdPartyLibFoo,
    // because we have already lost track of the registered Collectors,
    // however since we've collected Foo's metrics in it's own Registry, and we have registered that
    // as a specific Collector, we can now de-register them:
    unregistered := reg.Unregister(firstCollector)
    if !unregistered {
    	panic("unexpected behavior of registry")
    }
    
    metricFamilies, err = reg.Gather()
    if err != nil {
    	panic("unexpected behavior of registry")
    }
    fmt.Println("First unregistered:")
    fmt.Println(toNormalizedJSON(sanitizeMetricFamily(metricFamilies[0])))
    
    // Now we can create another instance of Foo with {instance: "first"} label again.
    firstRegAgain := prometheus.NewRegistry()
    _ = newThirdPartyLibFoo(firstRegAgain)
    firstCollectorAgain := prometheus.WrapCollectorWith(prometheus.Labels{"instance": "first"}, firstRegAgain)
    reg.MustRegister(firstCollectorAgain)
    
    metricFamilies, err = reg.Gather()
    if err != nil {
    	panic("unexpected behavior of registry")
    }
    fmt.Println("Both again:")
    fmt.Println(toNormalizedJSON(sanitizeMetricFamily(metricFamilies[0])))
    
    
    
    Output:
    
    Both instances:
    {"name":"foo","help":"Registered forever.","type":"GAUGE","metric":[{"label":[{"name":"instance","value":"first"}],"gauge":{"value":1}},{"label":[{"name":"instance","value":"second"}],"gauge":{"value":1}}]}
    First unregistered:
    {"name":"foo","help":"Registered forever.","type":"GAUGE","metric":[{"label":[{"name":"instance","value":"second"}],"gauge":{"value":1}}]}
    Both again:
    {"name":"foo","help":"Registered forever.","type":"GAUGE","metric":[{"label":[{"name":"instance","value":"first"}],"gauge":{"value":1}},{"label":[{"name":"instance","value":"second"}],"gauge":{"value":1}}]}
    

####  func [WrapCollectorWithPrefix](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/wrap.go#L108) ¶ added in v1.23.0
    
    
    func WrapCollectorWithPrefix(prefix [string](/builtin#string), c Collector) Collector

WrapCollectorWithPrefix returns a Collector wrapping the provided Collector. The wrapped Collector will add the provided prefix to the name of all Metrics it collects. 

See the documentation of WrapCollectorWith for more details on the use case. 

####  type [CollectorFunc](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/collectorfunc.go#L20) ¶ added in v1.22.0
    
    
    type CollectorFunc func(chan<- Metric)

CollectorFunc is a convenient way to implement a Prometheus Collector without interface boilerplate. This implementation is based on DescribeByCollect method. familiarize yourself to it before using. 

Example ¶

Using CollectorFunc that registers the metric info for the HTTP requests. 
    
    
    desc := prometheus.NewDesc(
    	"http_requests_info",
    	"Information about the received HTTP requests.",
    	[]string{"code", "method"},
    	nil,
    )
    
    // Example 1: 42 GET requests with 200 OK status code.
    collector := prometheus.CollectorFunc(func(ch chan<- prometheus.Metric) {
    	ch <- prometheus.MustNewConstMetric(
    		desc,
    		prometheus.CounterValue, // Metric type: Counter
    		42,                      // Value
    		"200",                   // Label value: HTTP status code
    		"GET",                   // Label value: HTTP method
    	)
    
    	// Example 2: 15 POST requests with 404 Not Found status code.
    	ch <- prometheus.MustNewConstMetric(
    		desc,
    		prometheus.CounterValue,
    		15,
    		"404",
    		"POST",
    	)
    })
    
    prometheus.MustRegister(collector)
    
    // Just for demonstration, let's check the state of the metric by registering
    // it with a custom registry and then let it collect the metrics.
    
    reg := prometheus.NewRegistry()
    reg.MustRegister(collector)
    
    metricFamilies, err := reg.Gather()
    if err != nil || len(metricFamilies) != 1 {
    	panic("unexpected behavior of custom test registry")
    }
    
    fmt.Println(toNormalizedJSON(sanitizeMetricFamily(metricFamilies[0])))
    
    
    
    Output:
    
    {"name":"http_requests_info","help":"Information about the received HTTP requests.","type":"COUNTER","metric":[{"label":[{"name":"code","value":"200"},{"name":"method","value":"GET"}],"counter":{"value":42}},{"label":[{"name":"code","value":"404"},{"name":"method","value":"POST"}],"counter":{"value":15}}]}
    

####  func (CollectorFunc) [Collect](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/collectorfunc.go#L23) ¶ added in v1.22.0
    
    
    func (f CollectorFunc) Collect(ch chan<- Metric)

Collect calls the defined CollectorFunc function with the provided Metrics channel 

####  func (CollectorFunc) [Describe](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/collectorfunc.go#L28) ¶ added in v1.22.0
    
    
    func (f CollectorFunc) Describe(ch chan<- *Desc)

Describe sends the descriptor information using DescribeByCollect 

####  type [ConstrainableLabels](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/labels.go#L56) ¶ added in v1.15.0
    
    
    type ConstrainableLabels interface {
    	// contains filtered or unexported methods
    }

ConstrainableLabels is an interface that allows creating of labels that can be optionally constrained. 
    
    
    prometheus.V2().NewCounterVec(CounterVecOpts{
      CounterOpts: {...}, // Usual CounterOpts fields
      VariableLabels: []ConstrainedLabels{
        {Name: "A"},
        {Name: "B", Constraint: func(v string) string { ... }},
      },
    })
    

####  type [ConstrainedLabel](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/labels.go#L41) ¶ added in v1.15.0
    
    
    type ConstrainedLabel struct {
    	Name       [string](/builtin#string)
    	Constraint LabelConstraint
    }

ConstrainedLabels represents a label name and its constrain function to normalize label values. This type is commonly used when constructing metric vector Collectors. 

####  type [ConstrainedLabels](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/labels.go#L64) ¶ added in v1.15.0
    
    
    type ConstrainedLabels []ConstrainedLabel

ConstrainedLabels represents a collection of label name -> constrain function to normalize label values. This type is commonly used when constructing metric vector Collectors. 

####  type [Counter](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/counter.go#L35) ¶
    
    
    type Counter interface {
    	Metric
    	Collector
    
    	// Inc increments the counter by 1. Use Add to increment it by arbitrary
    	// non-negative values.
    	Inc()
    	// Add adds the given value to the counter. It panics if the value is <
    	// 0.
    	Add([float64](/builtin#float64))
    }

Counter is a Metric that represents a single numerical value that only ever goes up. That implies that it cannot be used to count items whose number can also go down, e.g. the number of currently running goroutines. Those "counters" are represented by Gauges. 

A Counter is typically used to count requests served, tasks completed, errors occurred, etc. 

To create Counter instances, use NewCounter. 

####  func [NewCounter](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/counter.go#L87) ¶
    
    
    func NewCounter(opts CounterOpts) Counter

NewCounter creates a new Counter based on the provided CounterOpts. 

The returned implementation also implements ExemplarAdder. It is safe to perform the corresponding type assertion. 

The returned implementation tracks the counter value in two separate variables, a float64 and a uint64. The latter is used to track calls of the Inc method and calls of the Add method with a value that can be represented as a uint64. This allows atomic increments of the counter with optimal performance. (It is common to have an Inc call in very hot execution paths.) Both internal tracking values are added up in the Write method. This has to be taken into account when it comes to precision and overflow behavior. 

####  type [CounterFunc](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/counter.go#L336) ¶
    
    
    type CounterFunc interface {
    	Metric
    	Collector
    }

CounterFunc is a Counter whose value is determined at collect time by calling a provided function. 

To create CounterFunc instances, use NewCounterFunc. 

####  func [NewCounterFunc](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/counter.go#L351) ¶
    
    
    func NewCounterFunc(opts CounterOpts, function func() [float64](/builtin#float64)) CounterFunc

NewCounterFunc creates a new CounterFunc based on the provided CounterOpts. The value reported is determined by calling the given function from within the Write method. Take into account that metric collection may happen concurrently. If that results in concurrent calls to Write, like in the case where a CounterFunc is directly registered with Prometheus, the provided function must be concurrency-safe. The function should also honor the contract for a Counter (values only go up, not down), but compliance will not be checked. 

Check out the ExampleGaugeFunc examples for the similar GaugeFunc. 

####  type [CounterOpts](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/counter.go#L61) ¶
    
    
    type CounterOpts Opts

CounterOpts is an alias for Opts. See there for doc comments. 

####  type [CounterVec](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/counter.go#L188) ¶
    
    
    type CounterVec struct {
    	*MetricVec
    }

CounterVec is a Collector that bundles a set of Counters that all share the same Desc, but have different values for their variable labels. This is used if you want to count the same thing partitioned by various dimensions (e.g. number of HTTP requests, partitioned by response code and method). Create instances with NewCounterVec. 

Example ¶
    
    
    httpReqs := prometheus.NewCounterVec(
    	prometheus.CounterOpts{
    		Name: "http_requests_total",
    		Help: "How many HTTP requests processed, partitioned by status code and HTTP method.",
    	},
    	[]string{"code", "method"},
    )
    prometheus.MustRegister(httpReqs)
    
    httpReqs.WithLabelValues("404", "POST").Add(42)
    
    // If you have to access the same set of labels very frequently, it
    // might be good to retrieve the metric only once and keep a handle to
    // it. But beware of deletion of that metric, see below!
    m := httpReqs.WithLabelValues("200", "GET")
    for i := 0; i < 1000000; i++ {
    	m.Inc()
    }
    // Delete a metric from the vector. If you have previously kept a handle
    // to that metric (as above), future updates via that handle will go
    // unseen (even if you re-create a metric with the same label set
    // later).
    httpReqs.DeleteLabelValues("200", "GET")
    // Same thing with the more verbose Labels syntax.
    httpReqs.Delete(prometheus.Labels{"method": "GET", "code": "200"})
    
    // Just for demonstration, let's check the state of the counter vector
    // by registering it with a custom registry and then let it collect the
    // metrics.
    reg := prometheus.NewRegistry()
    reg.MustRegister(httpReqs)
    
    metricFamilies, err := reg.Gather()
    if err != nil || len(metricFamilies) != 1 {
    	panic("unexpected behavior of custom test registry")
    }
    
    fmt.Println(toNormalizedJSON(sanitizeMetricFamily(metricFamilies[0])))
    
    
    
    Output:
    
    {"name":"http_requests_total","help":"How many HTTP requests processed, partitioned by status code and HTTP method.","type":"COUNTER","metric":[{"label":[{"name":"code","value":"404"},{"name":"method","value":"POST"}],"counter":{"value":42,"createdTimestamp":"1970-01-01T00:00:10Z"}}]}
    

####  func [NewCounterVec](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/counter.go#L194) ¶
    
    
    func NewCounterVec(opts CounterOpts, labelNames [][string](/builtin#string)) *CounterVec

NewCounterVec creates a new CounterVec based on the provided CounterOpts and partitioned by the given label names. 

####  func (*CounterVec) [CurryWith](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/counter.go#L314) ¶ added in v0.9.0
    
    
    func (v *CounterVec) CurryWith(labels Labels) (*CounterVec, [error](/builtin#error))

CurryWith returns a vector curried with the provided labels, i.e. the returned vector has those labels pre-set for all labeled operations performed on it. The cardinality of the curried vector is reduced accordingly. The order of the remaining labels stays the same (just with the curried labels taken out of the sequence – which is relevant for the (GetMetric)WithLabelValues methods). It is possible to curry a curried vector, but only with labels not yet used for currying before. 

The metrics contained in the CounterVec are shared between the curried and uncurried vectors. They are just accessed differently. Curried and uncurried vectors behave identically in terms of collection. Only one must be registered with a given registry (usually the uncurried version). The Reset method deletes all metrics, even if called on a curried vector. 

####  func (*CounterVec) [GetMetricWith](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/counter.go#L268) ¶
    
    
    func (v *CounterVec) GetMetricWith(labels Labels) (Counter, [error](/builtin#error))

GetMetricWith returns the Counter for the given Labels map (the label names must match those of the variable labels in Desc). If that label map is accessed for the first time, a new Counter is created. Implications of creating a Counter without using it and keeping the Counter for later use are the same as for GetMetricWithLabelValues. 

An error is returned if the number and names of the Labels are inconsistent with those of the variable labels in Desc (minus any curried labels). 

This method is used for the same purpose as GetMetricWithLabelValues(...string). See there for pros and cons of the two methods. 

####  func (*CounterVec) [GetMetricWithLabelValues](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/counter.go#L248) ¶
    
    
    func (v *CounterVec) GetMetricWithLabelValues(lvs ...[string](/builtin#string)) (Counter, [error](/builtin#error))

GetMetricWithLabelValues returns the Counter for the given slice of label values (same order as the variable labels in Desc). If that combination of label values is accessed for the first time, a new Counter is created. 

It is possible to call this method without using the returned Counter to only create the new Counter but leave it at its starting value 0. See also the SummaryVec example. 

Keeping the Counter for later use is possible (and should be considered if performance is critical), but keep in mind that Reset, DeleteLabelValues and Delete can be used to delete the Counter from the CounterVec. In that case, the Counter will still exist, but it will not be exported anymore, even if a Counter with the same label values is created later. 

An error is returned if the number of label values is not the same as the number of variable labels in Desc (minus any curried labels). 

Note that for more than one label value, this method is prone to mistakes caused by an incorrect order of arguments. Consider GetMetricWith(Labels) as an alternative to avoid that type of mistake. For higher label numbers, the latter has a much more readable (albeit more verbose) syntax, but it comes with a performance overhead (for creating and processing the Labels map). See also the GaugeVec example. 

####  func (*CounterVec) [MustCurryWith](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/counter.go#L324) ¶ added in v0.9.0
    
    
    func (v *CounterVec) MustCurryWith(labels Labels) *CounterVec

MustCurryWith works as CurryWith but panics where CurryWith would have returned an error. 

####  func (*CounterVec) [With](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/counter.go#L293) ¶
    
    
    func (v *CounterVec) With(labels Labels) Counter

With works as GetMetricWith, but panics where GetMetricWithLabels would have returned an error. Not returning an error allows shortcuts like 
    
    
    myVec.With(prometheus.Labels{"code": "404", "method": "GET"}).Add(42)
    

####  func (*CounterVec) [WithLabelValues](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/counter.go#L281) ¶
    
    
    func (v *CounterVec) WithLabelValues(lvs ...[string](/builtin#string)) Counter

WithLabelValues works as GetMetricWithLabelValues, but panics where GetMetricWithLabelValues would have returned an error. Not returning an error allows shortcuts like 
    
    
    myVec.WithLabelValues("404", "GET").Add(42)
    

####  type [CounterVecOpts](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/counter.go#L66) ¶ added in v1.15.0
    
    
    type CounterVecOpts struct {
    	CounterOpts
    
    	// VariableLabels are used to partition the metric vector by the given set
    	// of labels. Each label value will be constrained with the optional Constraint
    	// function, if provided.
    	VariableLabels ConstrainableLabels
    }

CounterVecOpts bundles the options to create a CounterVec metric. It is mandatory to set CounterOpts, see there for mandatory fields. VariableLabels is optional and can safely be left to its default value. 

####  type [Desc](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/desc.go#L45) ¶
    
    
    type Desc struct {
    	// contains filtered or unexported fields
    }

Desc is the descriptor used by every Prometheus Metric. It is essentially the immutable meta-data of a Metric. The normal Metric implementations included in this package manage their Desc under the hood. Users only have to deal with Desc if they use advanced features like the ExpvarCollector or custom Collectors and Metrics. 

Descriptors registered with the same registry have to fulfill certain consistency and uniqueness criteria if they share the same fully-qualified name: They must have the same help string and the same label names (aka label dimensions) in each, constLabels and variableLabels, but they must differ in the values of the constLabels. 

Descriptors that share the same fully-qualified names and the same label values of their constLabels are considered equal. 

Use NewDesc to create new Desc instances. 

####  func [NewDesc](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/desc.go#L78) ¶
    
    
    func NewDesc(fqName, help [string](/builtin#string), variableLabels [][string](/builtin#string), constLabels Labels) *Desc

NewDesc allocates and initializes a new Desc. Errors are recorded in the Desc and will be reported on registration time. variableLabels and constLabels can be nil if no such labels should be set. fqName must not be empty. 

variableLabels only contain the label names. Their label values are variable and therefore not part of the Desc. (They are managed within the Metric.) 

For constLabels, the label values are constant. Therefore, they are fully specified in the Desc. See the Collector example for a usage pattern. 

####  func [NewInvalidDesc](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/desc.go#L179) ¶
    
    
    func NewInvalidDesc(err [error](/builtin#error)) *Desc

NewInvalidDesc returns an invalid descriptor, i.e. a descriptor with the provided error set. If a collector returning such a descriptor is registered, registration will fail with the provided error. NewInvalidDesc can be used by a Collector to signal inability to describe itself. 

####  func (*Desc) [String](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/desc.go#L185) ¶
    
    
    func (d *Desc) String() [string](/builtin#string)

####  type [Exemplar](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/metric.go#L225) ¶ added in v1.13.0
    
    
    type Exemplar struct {
    	Value  [float64](/builtin#float64)
    	Labels Labels
    	// Optional.
    	// Default value (time.Time{}) indicates its empty, which should be
    	// understood as time.Now() time at the moment of creation of metric.
    	Timestamp [time](/time).[Time](/time#Time)
    }

Exemplar is easier to use, user-facing representation of *dto.Exemplar. 

####  type [ExemplarAdder](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/counter.go#L56) ¶ added in v0.12.1
    
    
    type ExemplarAdder interface {
    	AddWithExemplar(value [float64](/builtin#float64), exemplar Labels)
    }

ExemplarAdder is implemented by Counters that offer the option of adding a value to the Counter together with an exemplar. Its AddWithExemplar method works like the Add method of the Counter interface but also replaces the currently saved exemplar (if any) with a new one, created from the provided value, the current time as timestamp, and the provided labels. Empty Labels will lead to a valid (label-less) exemplar. But if Labels is nil, the current exemplar is left in place. AddWithExemplar panics if the value is < 0, if any of the provided labels are invalid, or if the provided labels contain more than 128 runes in total. 

####  type [ExemplarObserver](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/observer.go#L62) ¶ added in v0.12.1
    
    
    type ExemplarObserver interface {
    	ObserveWithExemplar(value [float64](/builtin#float64), exemplar Labels)
    }

ExemplarObserver is implemented by Observers that offer the option of observing a value together with an exemplar. Its ObserveWithExemplar method works like the Observe method of an Observer but also replaces the currently saved exemplar (if any) with a new one, created from the provided value, the current time as timestamp, and the provided Labels. Empty Labels will lead to a valid (label-less) exemplar. But if Labels is nil, the current exemplar is left in place. ObserveWithExemplar panics if any of the provided labels are invalid or if the provided labels contain more than 128 runes in total. 

####  type [Gatherer](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/registry.go#L140) ¶
    
    
    type Gatherer interface {
    	// Gather calls the Collect method of the registered Collectors and then
    	// gathers the collected metrics into a lexicographically sorted slice
    	// of uniquely named MetricFamily protobufs. Gather ensures that the
    	// returned slice is valid and self-consistent so that it can be used
    	// for valid exposition. As an exception to the strict consistency
    	// requirements described for metric.Desc, Gather will tolerate
    	// different sets of label names for metrics of the same metric family.
    	//
    	// Even if an error occurs, Gather attempts to gather as many metrics as
    	// possible. Hence, if a non-nil error is returned, the returned
    	// MetricFamily slice could be nil (in case of a fatal error that
    	// prevented any meaningful metric collection) or contain a number of
    	// MetricFamily protobufs, some of which might be incomplete, and some
    	// might be missing altogether. The returned error (which might be a
    	// MultiError) explains the details. Note that this is mostly useful for
    	// debugging purposes. If the gathered protobufs are to be used for
    	// exposition in actual monitoring, it is almost always better to not
    	// expose an incomplete result and instead disregard the returned
    	// MetricFamily protobufs in case the returned error is non-nil.
    	Gather() ([]*[dto](/github.com/prometheus/client_model/go).[MetricFamily](/github.com/prometheus/client_model/go#MetricFamily), [error](/builtin#error))
    }

Gatherer is the interface for the part of a registry in charge of gathering the collected metrics into a number of MetricFamilies. The Gatherer interface comes with the same general implication as described for the Registerer interface. 

####  type [GathererFunc](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/registry.go#L190) ¶
    
    
    type GathererFunc func() ([]*[dto](/github.com/prometheus/client_model/go).[MetricFamily](/github.com/prometheus/client_model/go#MetricFamily), [error](/builtin#error))

GathererFunc turns a function into a Gatherer. 

####  func (GathererFunc) [Gather](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/registry.go#L193) ¶
    
    
    func (gf GathererFunc) Gather() ([]*[dto](/github.com/prometheus/client_model/go).[MetricFamily](/github.com/prometheus/client_model/go#MetricFamily), [error](/builtin#error))

Gather implements Gatherer. 

####  type [Gatherers](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/registry.go#L743) ¶
    
    
    type Gatherers []Gatherer

Gatherers is a slice of Gatherer instances that implements the Gatherer interface itself. Its Gather method calls Gather on all Gatherers in the slice in order and returns the merged results. Errors returned from the Gather calls are all returned in a flattened MultiError. Duplicate and inconsistent Metrics are skipped (first occurrence in slice order wins) and reported in the returned error. 

Gatherers can be used to merge the Gather results from multiple Registries. It also provides a way to directly inject existing MetricFamily protobufs into the gathering by creating a custom Gatherer with a Gather method that simply returns the existing MetricFamily protobufs. Note that no registration is involved (in contrast to Collector registration), so obviously registration-time checks cannot happen. Any inconsistencies between the gathered MetricFamilies are reported as errors by the Gather method, and inconsistent Metrics are dropped. Invalid parts of the MetricFamilies (e.g. syntactically invalid metric or label names) will go undetected. 

Example ¶
    
    
    package main
    
    import (
    	"bytes"
    	"fmt"
    	"strings"
    
    	dto "github.com/prometheus/client_model/go"
    	"github.com/prometheus/common/expfmt"
    	"github.com/prometheus/common/model"
    
    	"github.com/prometheus/client_golang/prometheus"
    )
    
    func main() {
    	reg := prometheus.NewRegistry()
    	temp := prometheus.NewGaugeVec(
    		prometheus.GaugeOpts{
    			Name: "temperature_kelvin",
    			Help: "Temperature in Kelvin.",
    		},
    		[]string{"location"},
    	)
    	reg.MustRegister(temp)
    	temp.WithLabelValues("outside").Set(273.14)
    	temp.WithLabelValues("inside").Set(298.44)
    
    	parser := expfmt.NewTextParser(model.UTF8Validation)
    
    	text := `
    # TYPE humidity_percent gauge
    # HELP humidity_percent Humidity in %.
    humidity_percent{location="outside"} 45.4
    humidity_percent{location="inside"} 33.2
    # TYPE temperature_kelvin gauge
    # HELP temperature_kelvin Temperature in Kelvin.
    temperature_kelvin{location="somewhere else"} 4.5
    `
    
    	parseText := func() ([]*dto.MetricFamily, error) {
    		parsed, err := parser.TextToMetricFamilies(strings.NewReader(text))
    		if err != nil {
    			return nil, err
    		}
    		var result []*dto.MetricFamily
    		for _, mf := range parsed {
    			result = append(result, mf)
    		}
    		return result, nil
    	}
    
    	gatherers := prometheus.Gatherers{
    		reg,
    		prometheus.GathererFunc(parseText),
    	}
    
    	gathering, err := gatherers.Gather()
    	if err != nil {
    		fmt.Println(err)
    	}
    
    	out := &bytes.Buffer{}
    	for _, mf := range gathering {
    		if _, err := expfmt.MetricFamilyToText(out, mf); err != nil {
    			panic(err)
    		}
    	}
    	fmt.Print(out.String())
    	fmt.Println("----------")
    
    	// Note how the temperature_kelvin metric family has been merged from
    	// different sources. Now try
    	text = `
    # TYPE humidity_percent gauge
    # HELP humidity_percent Humidity in %.
    humidity_percent{location="outside"} 45.4
    humidity_percent{location="inside"} 33.2
    # TYPE temperature_kelvin gauge
    # HELP temperature_kelvin Temperature in Kelvin.
    # Duplicate metric:
    temperature_kelvin{location="outside"} 265.3
     # Missing location label (note that this is undesirable but valid):
    temperature_kelvin 4.5
    `
    
    	gathering, err = gatherers.Gather()
    	if err != nil {
    		// We expect error collected metric "temperature_kelvin" { label:<name:"location" value:"outside" > gauge:<value:265.3 > } was collected before with the same name and label values
    		// We cannot assert it because of https://github.com/golang/protobuf/issues/1121
    		if strings.HasPrefix(err.Error(), `collected metric "temperature_kelvin" `) {
    			fmt.Println("Found duplicated metric `temperature_kelvin`")
    		} else {
    			fmt.Print(err)
    		}
    	}
    	// Note that still as many metrics as possible are returned:
    	out.Reset()
    	for _, mf := range gathering {
    		if _, err := expfmt.MetricFamilyToText(out, mf); err != nil {
    			panic(err)
    		}
    	}
    	fmt.Print(out.String())
    
    }
    
    
    
    Output:
    
    # HELP humidity_percent Humidity in %.
    # TYPE humidity_percent gauge
    humidity_percent{location="inside"} 33.2
    humidity_percent{location="outside"} 45.4
    # HELP temperature_kelvin Temperature in Kelvin.
    # TYPE temperature_kelvin gauge
    temperature_kelvin{location="inside"} 298.44
    temperature_kelvin{location="outside"} 273.14
    temperature_kelvin{location="somewhere else"} 4.5
    ----------
    Found duplicated metric `temperature_kelvin`
    # HELP humidity_percent Humidity in %.
    # TYPE humidity_percent gauge
    humidity_percent{location="inside"} 33.2
    humidity_percent{location="outside"} 45.4
    # HELP temperature_kelvin Temperature in Kelvin.
    # TYPE temperature_kelvin gauge
    temperature_kelvin 4.5
    temperature_kelvin{location="inside"} 298.44
    temperature_kelvin{location="outside"} 273.14
    

Share Format Run

####  func (Gatherers) [Gather](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/registry.go#L746) ¶
    
    
    func (gs Gatherers) Gather() ([]*[dto](/github.com/prometheus/client_model/go).[MetricFamily](/github.com/prometheus/client_model/go#MetricFamily), [error](/builtin#error))

Gather implements Gatherer. 

####  type [Gauge](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/gauge.go#L32) ¶
    
    
    type Gauge interface {
    	Metric
    	Collector
    
    	// Set sets the Gauge to an arbitrary value.
    	Set([float64](/builtin#float64))
    	// Inc increments the Gauge by 1. Use Add to increment it by arbitrary
    	// values.
    	Inc()
    	// Dec decrements the Gauge by 1. Use Sub to decrement it by arbitrary
    	// values.
    	Dec()
    	// Add adds the given value to the Gauge. (The value can be negative,
    	// resulting in a decrease of the Gauge.)
    	Add([float64](/builtin#float64))
    	// Sub subtracts the given value from the Gauge. (The value can be
    	// negative, resulting in an increase of the Gauge.)
    	Sub([float64](/builtin#float64))
    
    	// SetToCurrentTime sets the Gauge to the current Unix time in seconds.
    	SetToCurrentTime()
    }

Gauge is a Metric that represents a single numerical value that can arbitrarily go up and down. 

A Gauge is typically used for measured values like temperatures or current memory usage, but also "counts" that can go up and down, like the number of running goroutines. 

To create Gauge instances, use NewGauge. 

Example ¶
    
    
    package main
    
    import (
    	"github.com/prometheus/client_golang/prometheus"
    )
    
    func main() {
    	opsQueued := prometheus.NewGauge(prometheus.GaugeOpts{
    		Namespace: "our_company",
    		Subsystem: "blob_storage",
    		Name:      "ops_queued",
    		Help:      "Number of blob storage operations waiting to be processed.",
    	})
    	prometheus.MustRegister(opsQueued)
    
    	// 10 operations queued by the goroutine managing incoming requests.
    	opsQueued.Add(10)
    	// A worker goroutine has picked up a waiting operation.
    	opsQueued.Dec()
    	// And once more...
    	opsQueued.Dec()
    }
    

Share Format Run

####  func [NewGauge](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/gauge.go#L78) ¶
    
    
    func NewGauge(opts GaugeOpts) Gauge

NewGauge creates a new Gauge based on the provided GaugeOpts. 

The returned implementation is optimized for a fast Set method. If you have a choice for managing the value of a Gauge via Set vs. Inc/Dec/Add/Sub, pick the former. For example, the Inc method of the returned Gauge is slower than the Inc method of a Counter returned by NewCounter. This matches the typical scenarios for Gauges and Counters, where the former tends to be Set-heavy and the latter Inc-heavy. 

####  type [GaugeFunc](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/gauge.go#L290) ¶
    
    
    type GaugeFunc interface {
    	Metric
    	Collector
    }

GaugeFunc is a Gauge whose value is determined at collect time by calling a provided function. 

To create GaugeFunc instances, use NewGaugeFunc. 

Example (ConstLabels) ¶
    
    
    package main
    
    import (
    	"fmt"
    
    	"github.com/prometheus/client_golang/prometheus"
    )
    
    func main() {
    	// primaryDB and secondaryDB represent two example *sql.DB connections we want to instrument.
    	var primaryDB, secondaryDB interface {
    		Stats() struct{ OpenConnections int }
    	}
    
    	if err := prometheus.Register(prometheus.NewGaugeFunc(
    		prometheus.GaugeOpts{
    			Namespace:   "mysql",
    			Name:        "connections_open",
    			Help:        "Number of mysql connections open.",
    			ConstLabels: prometheus.Labels{"destination": "primary"},
    		},
    		func() float64 { return float64(primaryDB.Stats().OpenConnections) },
    	)); err == nil {
    		fmt.Println(`GaugeFunc 'connections_open' for primary DB connection registered with labels {destination="primary"}`)
    	}
    
    	if err := prometheus.Register(prometheus.NewGaugeFunc(
    		prometheus.GaugeOpts{
    			Namespace:   "mysql",
    			Name:        "connections_open",
    			Help:        "Number of mysql connections open.",
    			ConstLabels: prometheus.Labels{"destination": "secondary"},
    		},
    		func() float64 { return float64(secondaryDB.Stats().OpenConnections) },
    	)); err == nil {
    		fmt.Println(`GaugeFunc 'connections_open' for secondary DB connection registered with labels {destination="secondary"}`)
    	}
    
    	// Note that we can register more than once GaugeFunc with same metric name
    	// as long as their const labels are consistent.
    
    }
    
    
    
    Output:
    
    GaugeFunc 'connections_open' for primary DB connection registered with labels {destination="primary"}
    GaugeFunc 'connections_open' for secondary DB connection registered with labels {destination="secondary"}
    

Share Format Run

Example (Simple) ¶
    
    
    package main
    
    import (
    	"fmt"
    	"runtime"
    
    	"github.com/prometheus/client_golang/prometheus"
    )
    
    func main() {
    	if err := prometheus.Register(prometheus.NewGaugeFunc(
    		prometheus.GaugeOpts{
    			Subsystem: "runtime",
    			Name:      "goroutines_count",
    			Help:      "Number of goroutines that currently exist.",
    		},
    		func() float64 { return float64(runtime.NumGoroutine()) },
    	)); err == nil {
    		fmt.Println("GaugeFunc 'goroutines_count' registered.")
    	}
    	// Note that the count of goroutines is a gauge (and not a counter) as
    	// it can go up and down.
    
    }
    
    
    
    Output:
    
    GaugeFunc 'goroutines_count' registered.
    

Share Format Run

####  func [NewGaugeFunc](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/gauge.go#L304) ¶
    
    
    func NewGaugeFunc(opts GaugeOpts, function func() [float64](/builtin#float64)) GaugeFunc

NewGaugeFunc creates a new GaugeFunc based on the provided GaugeOpts. The value reported is determined by calling the given function from within the Write method. Take into account that metric collection may happen concurrently. Therefore, it must be safe to call the provided function concurrently. 

NewGaugeFunc is a good way to create an “info” style metric with a constant value of 1. Example: <https://github.com/prometheus/common/blob/8558a5b7db3c84fa38b4766966059a7bd5bfa2ee/version/info.go#L36-L56>

####  type [GaugeOpts](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/gauge.go#L56) ¶
    
    
    type GaugeOpts Opts

GaugeOpts is an alias for Opts. See there for doc comments. 

####  type [GaugeVec](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/gauge.go#L146) ¶
    
    
    type GaugeVec struct {
    	*MetricVec
    }

GaugeVec is a Collector that bundles a set of Gauges that all share the same Desc, but have different values for their variable labels. This is used if you want to count the same thing partitioned by various dimensions (e.g. number of operations queued, partitioned by user and operation type). Create instances with NewGaugeVec. 

Example ¶
    
    
    package main
    
    import (
    	"github.com/prometheus/client_golang/prometheus"
    )
    
    func main() {
    	opsQueued := prometheus.NewGaugeVec(
    		prometheus.GaugeOpts{
    			Namespace: "our_company",
    			Subsystem: "blob_storage",
    			Name:      "ops_queued",
    			Help:      "Number of blob storage operations waiting to be processed, partitioned by user and type.",
    		},
    		[]string{
    			// Which user has requested the operation?
    			"user",
    			// Of what type is the operation?
    			"type",
    		},
    	)
    	prometheus.MustRegister(opsQueued)
    
    	// Increase a value using compact (but order-sensitive!) WithLabelValues().
    	opsQueued.WithLabelValues("bob", "put").Add(4)
    	// Increase a value with a map using WithLabels. More verbose, but order
    	// doesn't matter anymore.
    	opsQueued.With(prometheus.Labels{"type": "delete", "user": "alice"}).Inc()
    }
    

Share Format Run

####  func [NewGaugeVec](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/gauge.go#L152) ¶
    
    
    func NewGaugeVec(opts GaugeOpts, labelNames [][string](/builtin#string)) *GaugeVec

NewGaugeVec creates a new GaugeVec based on the provided GaugeOpts and partitioned by the given label names. 

####  func (*GaugeVec) [CurryWith](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/gauge.go#L268) ¶ added in v0.9.0
    
    
    func (v *GaugeVec) CurryWith(labels Labels) (*GaugeVec, [error](/builtin#error))

CurryWith returns a vector curried with the provided labels, i.e. the returned vector has those labels pre-set for all labeled operations performed on it. The cardinality of the curried vector is reduced accordingly. The order of the remaining labels stays the same (just with the curried labels taken out of the sequence – which is relevant for the (GetMetric)WithLabelValues methods). It is possible to curry a curried vector, but only with labels not yet used for currying before. 

The metrics contained in the GaugeVec are shared between the curried and uncurried vectors. They are just accessed differently. Curried and uncurried vectors behave identically in terms of collection. Only one must be registered with a given registry (usually the uncurried version). The Reset method deletes all metrics, even if called on a curried vector. 

####  func (*GaugeVec) [GetMetricWith](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/gauge.go#L222) ¶
    
    
    func (v *GaugeVec) GetMetricWith(labels Labels) (Gauge, [error](/builtin#error))

GetMetricWith returns the Gauge for the given Labels map (the label names must match those of the variable labels in Desc). If that label map is accessed for the first time, a new Gauge is created. Implications of creating a Gauge without using it and keeping the Gauge for later use are the same as for GetMetricWithLabelValues. 

An error is returned if the number and names of the Labels are inconsistent with those of the variable labels in Desc (minus any curried labels). 

This method is used for the same purpose as GetMetricWithLabelValues(...string). See there for pros and cons of the two methods. 

####  func (*GaugeVec) [GetMetricWithLabelValues](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/gauge.go#L202) ¶
    
    
    func (v *GaugeVec) GetMetricWithLabelValues(lvs ...[string](/builtin#string)) (Gauge, [error](/builtin#error))

GetMetricWithLabelValues returns the Gauge for the given slice of label values (same order as the variable labels in Desc). If that combination of label values is accessed for the first time, a new Gauge is created. 

It is possible to call this method without using the returned Gauge to only create the new Gauge but leave it at its starting value 0. See also the SummaryVec example. 

Keeping the Gauge for later use is possible (and should be considered if performance is critical), but keep in mind that Reset, DeleteLabelValues and Delete can be used to delete the Gauge from the GaugeVec. In that case, the Gauge will still exist, but it will not be exported anymore, even if a Gauge with the same label values is created later. See also the CounterVec example. 

An error is returned if the number of label values is not the same as the number of variable labels in Desc (minus any curried labels). 

Note that for more than one label value, this method is prone to mistakes caused by an incorrect order of arguments. Consider GetMetricWith(Labels) as an alternative to avoid that type of mistake. For higher label numbers, the latter has a much more readable (albeit more verbose) syntax, but it comes with a performance overhead (for creating and processing the Labels map). 

####  func (*GaugeVec) [MustCurryWith](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/gauge.go#L278) ¶ added in v0.9.0
    
    
    func (v *GaugeVec) MustCurryWith(labels Labels) *GaugeVec

MustCurryWith works as CurryWith but panics where CurryWith would have returned an error. 

####  func (*GaugeVec) [With](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/gauge.go#L247) ¶
    
    
    func (v *GaugeVec) With(labels Labels) Gauge

With works as GetMetricWith, but panics where GetMetricWithLabels would have returned an error. Not returning an error allows shortcuts like 
    
    
    myVec.With(prometheus.Labels{"code": "404", "method": "GET"}).Add(42)
    

####  func (*GaugeVec) [WithLabelValues](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/gauge.go#L235) ¶
    
    
    func (v *GaugeVec) WithLabelValues(lvs ...[string](/builtin#string)) Gauge

WithLabelValues works as GetMetricWithLabelValues, but panics where GetMetricWithLabelValues would have returned an error. Not returning an error allows shortcuts like 
    
    
    myVec.WithLabelValues("404", "GET").Add(42)
    

####  type [GaugeVecOpts](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/gauge.go#L61) ¶ added in v1.15.0
    
    
    type GaugeVecOpts struct {
    	GaugeOpts
    
    	// VariableLabels are used to partition the metric vector by the given set
    	// of labels. Each label value will be constrained with the optional Constraint
    	// function, if provided.
    	VariableLabels ConstrainableLabels
    }

GaugeVecOpts bundles the options to create a GaugeVec metric. It is mandatory to set GaugeOpts, see there for mandatory fields. VariableLabels is optional and can safely be left to its default value. 

####  type [Histogram](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/histogram.go#L249) ¶
    
    
    type Histogram interface {
    	Metric
    	Collector
    
    	// Observe adds a single observation to the histogram. Observations are
    	// usually positive or zero. Negative observations are accepted but
    	// prevent current versions of Prometheus from properly detecting
    	// counter resets in the sum of observations. (The experimental Native
    	// Histograms handle negative observations properly.) See
    	// <https://prometheus.io/docs/practices/histograms/#count-and-sum-of-observations>
    	// for details.
    	Observe([float64](/builtin#float64))
    }

A Histogram counts individual observations from an event or sample stream in configurable static buckets (or in dynamic sparse buckets as part of the experimental Native Histograms, see below for more details). Similar to a Summary, it also provides a sum of observations and an observation count. 

On the Prometheus server, quantiles can be calculated from a Histogram using the histogram_quantile PromQL function. 

Note that Histograms, in contrast to Summaries, can be aggregated in PromQL (see the documentation for detailed procedures). However, Histograms require the user to pre-define suitable buckets, and they are in general less accurate. (Both problems are addressed by the experimental Native Histograms. To use them, configure a NativeHistogramBucketFactor in the HistogramOpts. They also require a Prometheus server v2.40+ with the corresponding feature flag enabled.) 

The Observe method of a Histogram has a very low performance overhead in comparison with the Observe method of a Summary. 

To create Histogram instances, use NewHistogram. 

Example ¶
    
    
    temps := prometheus.NewHistogram(prometheus.HistogramOpts{
    	Name:    "pond_temperature_celsius",
    	Help:    "The temperature of the frog pond.", // Sorry, we can't measure how badly it smells.
    	Buckets: prometheus.LinearBuckets(20, 5, 5),  // 5 buckets, each 5 centigrade wide.
    })
    
    // Simulate some observations.
    for i := 0; i < 1000; i++ {
    	temps.Observe(30 + math.Floor(120*math.Sin(float64(i)*0.1))/10)
    }
    
    // Just for demonstration, let's check the state of the histogram by
    // (ab)using its Write method (which is usually only used by Prometheus
    // internally).
    metric := &dto.Metric{}
    temps.Write(metric)
    
    fmt.Println(toNormalizedJSON(sanitizeMetric(metric)))
    
    
    
    Output:
    
    {"histogram":{"sampleCount":"1000","sampleSum":29969.50000000001,"bucket":[{"cumulativeCount":"192","upperBound":20},{"cumulativeCount":"366","upperBound":25},{"cumulativeCount":"501","upperBound":30},{"cumulativeCount":"638","upperBound":35},{"cumulativeCount":"816","upperBound":40}],"createdTimestamp":"1970-01-01T00:00:10Z"}}
    

####  func [NewHistogram](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/histogram.go#L523) ¶
    
    
    func NewHistogram(opts HistogramOpts) Histogram

NewHistogram creates a new Histogram based on the provided HistogramOpts. It panics if the buckets in HistogramOpts are not in strictly increasing order. 

The returned implementation also implements ExemplarObserver. It is safe to perform the corresponding type assertion. Exemplars are tracked separately for each bucket. 

####  type [HistogramOpts](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/histogram.go#L365) ¶
    
    
    type HistogramOpts struct {
    	// Namespace, Subsystem, and Name are components of the fully-qualified
    	// name of the Histogram (created by joining these components with
    	// "_"). Only Name is mandatory, the others merely help structuring the
    	// name. Note that the fully-qualified name of the Histogram must be a
    	// valid Prometheus metric name.
    	Namespace [string](/builtin#string)
    	Subsystem [string](/builtin#string)
    	Name      [string](/builtin#string)
    
    	// Help provides information about this Histogram.
    	//
    	// Metrics with the same fully-qualified name must have the same Help
    	// string.
    	Help [string](/builtin#string)
    
    	// ConstLabels are used to attach fixed labels to this metric. Metrics
    	// with the same fully-qualified name must have the same label names in
    	// their ConstLabels.
    	//
    	// ConstLabels are only used rarely. In particular, do not use them to
    	// attach the same labels to all your metrics. Those use cases are
    	// better covered by target labels set by the scraping Prometheus
    	// server, or by one specific metric (e.g. a build_info or a
    	// machine_role metric). See also
    	// <https://prometheus.io/docs/instrumenting/writing_exporters/#target-labels-not-static-scraped-labels>
    	ConstLabels Labels
    
    	// Buckets defines the buckets into which observations are counted. Each
    	// element in the slice is the upper inclusive bound of a bucket. The
    	// values must be sorted in strictly increasing order. There is no need
    	// to add a highest bucket with +Inf bound, it will be added
    	// implicitly. If Buckets is left as nil or set to a slice of length
    	// zero, it is replaced by default buckets. The default buckets are
    	// DefBuckets if no buckets for a native histogram (see below) are used,
    	// otherwise the default is no buckets. (In other words, if you want to
    	// use both regular buckets and buckets for a native histogram, you have
    	// to define the regular buckets here explicitly.)
    	Buckets [][float64](/builtin#float64)
    
    	// If NativeHistogramBucketFactor is greater than one, so-called sparse
    	// buckets are used (in addition to the regular buckets, if defined
    	// above). A Histogram with sparse buckets will be ingested as a Native
    	// Histogram by a Prometheus server with that feature enabled (requires
    	// Prometheus v2.40+). Sparse buckets are exponential buckets covering
    	// the whole float64 range (with the exception of the “zero” bucket, see
    	// NativeHistogramZeroThreshold below). From any one bucket to the next,
    	// the width of the bucket grows by a constant
    	// factor. NativeHistogramBucketFactor provides an upper bound for this
    	// factor (exception see below). The smaller
    	// NativeHistogramBucketFactor, the more buckets will be used and thus
    	// the more costly the histogram will become. A generally good trade-off
    	// between cost and accuracy is a value of 1.1 (each bucket is at most
    	// 10% wider than the previous one), which will result in each power of
    	// two divided into 8 buckets (e.g. there will be 8 buckets between 1
    	// and 2, same as between 2 and 4, and 4 and 8, etc.).
    	//
    	// Details about the actually used factor: The factor is calculated as
    	// 2^(2^-n), where n is an integer number between (and including) -4 and
    	// 8. n is chosen so that the resulting factor is the largest that is
    	// still smaller or equal to NativeHistogramBucketFactor. Note that the
    	// smallest possible factor is therefore approx. 1.00271 (i.e. 2^(2^-8)
    	// ). If NativeHistogramBucketFactor is greater than 1 but smaller than
    	// 2^(2^-8), then the actually used factor is still 2^(2^-8) even though
    	// it is larger than the provided NativeHistogramBucketFactor.
    	//
    	// NOTE: Native Histograms are still an experimental feature. Their
    	// behavior might still change without a major version
    	// bump. Subsequently, all NativeHistogram... options here might still
    	// change their behavior or name (or might completely disappear) without
    	// a major version bump.
    	NativeHistogramBucketFactor [float64](/builtin#float64)
    	// All observations with an absolute value of less or equal
    	// NativeHistogramZeroThreshold are accumulated into a “zero” bucket.
    	// For best results, this should be close to a bucket boundary. This is
    	// usually the case if picking a power of two. If
    	// NativeHistogramZeroThreshold is left at zero,
    	// DefNativeHistogramZeroThreshold is used as the threshold. To
    	// configure a zero bucket with an actual threshold of zero (i.e. only
    	// observations of precisely zero will go into the zero bucket), set
    	// NativeHistogramZeroThreshold to the NativeHistogramZeroThresholdZero
    	// constant (or any negative float value).
    	NativeHistogramZeroThreshold [float64](/builtin#float64)
    
    	// The next three fields define a strategy to limit the number of
    	// populated sparse buckets. If NativeHistogramMaxBucketNumber is left
    	// at zero, the number of buckets is not limited. (Note that this might
    	// lead to unbounded memory consumption if the values observed by the
    	// Histogram are sufficiently wide-spread. In particular, this could be
    	// used as a DoS attack vector. Where the observed values depend on
    	// external inputs, it is highly recommended to set a
    	// NativeHistogramMaxBucketNumber.) Once the set
    	// NativeHistogramMaxBucketNumber is exceeded, the following strategy is
    	// enacted:
    	//  - First, if the last reset (or the creation) of the histogram is at
    	//    least NativeHistogramMinResetDuration ago, then the whole
    	//    histogram is reset to its initial state (including regular
    	//    buckets).
    	//  - If less time has passed, or if NativeHistogramMinResetDuration is
    	//    zero, no reset is performed. Instead, the zero threshold is
    	//    increased sufficiently to reduce the number of buckets to or below
    	//    NativeHistogramMaxBucketNumber, but not to more than
    	//    NativeHistogramMaxZeroThreshold. Thus, if
    	//    NativeHistogramMaxZeroThreshold is already at or below the current
    	//    zero threshold, nothing happens at this step.
    	//  - After that, if the number of buckets still exceeds
    	//    NativeHistogramMaxBucketNumber, the resolution of the histogram is
    	//    reduced by doubling the width of the sparse buckets (up to a
    	//    growth factor between one bucket to the next of 2^(2^4) = 65536,
    	//    see above).
    	//  - Any increased zero threshold or reduced resolution is reset back
    	//    to their original values once NativeHistogramMinResetDuration has
    	//    passed (since the last reset or the creation of the histogram).
    	NativeHistogramMaxBucketNumber  [uint32](/builtin#uint32)
    	NativeHistogramMinResetDuration [time](/time).[Duration](/time#Duration)
    	NativeHistogramMaxZeroThreshold [float64](/builtin#float64)
    
    	// NativeHistogramMaxExemplars limits the number of exemplars
    	// that are kept in memory for each native histogram. If you leave it at
    	// zero, a default value of 10 is used. If no exemplars should be kept specifically
    	// for native histograms, set it to a negative value. (Scrapers can
    	// still use the exemplars exposed for classic buckets, which are managed
    	// independently.)
    	NativeHistogramMaxExemplars [int](/builtin#int)
    	// NativeHistogramExemplarTTL is only checked once
    	// NativeHistogramMaxExemplars is exceeded. In that case, the
    	// oldest exemplar is removed if it is older than NativeHistogramExemplarTTL.
    	// Otherwise, the older exemplar in the pair of exemplars that are closest
    	// together (on an exponential scale) is removed.
    	// If NativeHistogramExemplarTTL is left at its zero value, a default value of
    	// 5m is used. To always delete the oldest exemplar, set it to a negative value.
    	NativeHistogramExemplarTTL [time](/time).[Duration](/time#Duration)
    	// contains filtered or unexported fields
    }

HistogramOpts bundles the options for creating a Histogram metric. It is mandatory to set Name to a non-empty string. All other fields are optional and can safely be left at their zero value, although it is strongly encouraged to set a Help string. 

####  type [HistogramVec](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/histogram.go#L1173) ¶
    
    
    type HistogramVec struct {
    	*MetricVec
    }

HistogramVec is a Collector that bundles a set of Histograms that all share the same Desc, but have different values for their variable labels. This is used if you want to count the same thing partitioned by various dimensions (e.g. HTTP request latencies, partitioned by status code and method). Create instances with NewHistogramVec. 

####  func [NewHistogramVec](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/histogram.go#L1179) ¶
    
    
    func NewHistogramVec(opts HistogramOpts, labelNames [][string](/builtin#string)) *HistogramVec

NewHistogramVec creates a new HistogramVec based on the provided HistogramOpts and partitioned by the given label names. 

####  func (*HistogramVec) [CurryWith](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/histogram.go#L1291) ¶ added in v0.9.0
    
    
    func (v *HistogramVec) CurryWith(labels Labels) (ObserverVec, [error](/builtin#error))

CurryWith returns a vector curried with the provided labels, i.e. the returned vector has those labels pre-set for all labeled operations performed on it. The cardinality of the curried vector is reduced accordingly. The order of the remaining labels stays the same (just with the curried labels taken out of the sequence – which is relevant for the (GetMetric)WithLabelValues methods). It is possible to curry a curried vector, but only with labels not yet used for currying before. 

The metrics contained in the HistogramVec are shared between the curried and uncurried vectors. They are just accessed differently. Curried and uncurried vectors behave identically in terms of collection. Only one must be registered with a given registry (usually the uncurried version). The Reset method deletes all metrics, even if called on a curried vector. 

####  func (*HistogramVec) [GetMetricWith](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/histogram.go#L1245) ¶
    
    
    func (v *HistogramVec) GetMetricWith(labels Labels) (Observer, [error](/builtin#error))

GetMetricWith returns the Histogram for the given Labels map (the label names must match those of the variable labels in Desc). If that label map is accessed for the first time, a new Histogram is created. Implications of creating a Histogram without using it and keeping the Histogram for later use are the same as for GetMetricWithLabelValues. 

An error is returned if the number and names of the Labels are inconsistent with those of the variable labels in Desc (minus any curried labels). 

This method is used for the same purpose as GetMetricWithLabelValues(...string). See there for pros and cons of the two methods. 

####  func (*HistogramVec) [GetMetricWithLabelValues](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/histogram.go#L1225) ¶
    
    
    func (v *HistogramVec) GetMetricWithLabelValues(lvs ...[string](/builtin#string)) (Observer, [error](/builtin#error))

GetMetricWithLabelValues returns the Histogram for the given slice of label values (same order as the variable labels in Desc). If that combination of label values is accessed for the first time, a new Histogram is created. 

It is possible to call this method without using the returned Histogram to only create the new Histogram but leave it at its starting value, a Histogram without any observations. 

Keeping the Histogram for later use is possible (and should be considered if performance is critical), but keep in mind that Reset, DeleteLabelValues and Delete can be used to delete the Histogram from the HistogramVec. In that case, the Histogram will still exist, but it will not be exported anymore, even if a Histogram with the same label values is created later. See also the CounterVec example. 

An error is returned if the number of label values is not the same as the number of variable labels in Desc (minus any curried labels). 

Note that for more than one label value, this method is prone to mistakes caused by an incorrect order of arguments. Consider GetMetricWith(Labels) as an alternative to avoid that type of mistake. For higher label numbers, the latter has a much more readable (albeit more verbose) syntax, but it comes with a performance overhead (for creating and processing the Labels map). See also the GaugeVec example. 

####  func (*HistogramVec) [MustCurryWith](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/histogram.go#L1301) ¶ added in v0.9.0
    
    
    func (v *HistogramVec) MustCurryWith(labels Labels) ObserverVec

MustCurryWith works as CurryWith but panics where CurryWith would have returned an error. 

####  func (*HistogramVec) [With](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/histogram.go#L1270) ¶
    
    
    func (v *HistogramVec) With(labels Labels) Observer

With works as GetMetricWith but panics where GetMetricWithLabels would have returned an error. Not returning an error allows shortcuts like 
    
    
    myVec.With(prometheus.Labels{"code": "404", "method": "GET"}).Observe(42.21)
    

####  func (*HistogramVec) [WithLabelValues](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/histogram.go#L1258) ¶
    
    
    func (v *HistogramVec) WithLabelValues(lvs ...[string](/builtin#string)) Observer

WithLabelValues works as GetMetricWithLabelValues, but panics where GetMetricWithLabelValues would have returned an error. Not returning an error allows shortcuts like 
    
    
    myVec.WithLabelValues("404", "GET").Observe(42.21)
    

####  type [HistogramVecOpts](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/histogram.go#L508) ¶ added in v1.15.0
    
    
    type HistogramVecOpts struct {
    	HistogramOpts
    
    	// VariableLabels are used to partition the metric vector by the given set
    	// of labels. Each label value will be constrained with the optional Constraint
    	// function, if provided.
    	VariableLabels ConstrainableLabels
    }

HistogramVecOpts bundles the options to create a HistogramVec metric. It is mandatory to set HistogramOpts, see there for mandatory fields. VariableLabels is optional and can safely be left to its default value. 

####  type [LabelConstraint](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/labels.go#L36) ¶ added in v1.17.0
    
    
    type LabelConstraint func([string](/builtin#string)) [string](/builtin#string)

LabelConstraint normalizes label values. 

####  type [Labels](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/labels.go#L33) ¶
    
    
    type Labels map[[string](/builtin#string)][string](/builtin#string)

Labels represents a collection of label name -> value mappings. This type is commonly used with the With(Labels) and GetMetricWith(Labels) methods of metric vector Collectors, e.g.: 
    
    
    myVec.With(Labels{"code": "404", "method": "GET"}).Add(42)
    

The other use-case is the specification of constant label pairs in Opts or to create a Desc. 

####  type [Metric](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/metric.go#L33) ¶
    
    
    type Metric interface {
    	// Desc returns the descriptor for the Metric. This method idempotently
    	// returns the same descriptor throughout the lifetime of the
    	// Metric. The returned descriptor is immutable by contract. A Metric
    	// unable to describe itself must return an invalid descriptor (created
    	// with NewInvalidDesc).
    	Desc() *Desc
    	// Write encodes the Metric into a "Metric" Protocol Buffer data
    	// transmission object.
    	//
    	// Metric implementations must observe concurrency safety as reads of
    	// this metric may occur at any time, and any blocking occurs at the
    	// expense of total performance of rendering all registered
    	// metrics. Ideally, Metric implementations should support concurrent
    	// readers.
    	//
    	// While populating dto.Metric, it is the responsibility of the
    	// implementation to ensure validity of the Metric protobuf (like valid
    	// UTF-8 strings or syntactically valid metric and label names). It is
    	// recommended to sort labels lexicographically. Callers of Write should
    	// still make sure of sorting if they depend on it.
    	Write(*[dto](/github.com/prometheus/client_model/go).[Metric](/github.com/prometheus/client_model/go#Metric)) [error](/builtin#error)
    }

A Metric models a single sample value with its meta data being exported to Prometheus. Implementations of Metric in this package are Gauge, Counter, Histogram, Summary, and Untyped. 

####  func [MustNewConstHistogram](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/histogram.go#L1386) ¶
    
    
    func MustNewConstHistogram(
    	desc *Desc,
    	count [uint64](/builtin#uint64),
    	sum [float64](/builtin#float64),
    	buckets map[[float64](/builtin#float64)][uint64](/builtin#uint64),
    	labelValues ...[string](/builtin#string),
    ) Metric

MustNewConstHistogram is a version of NewConstHistogram that panics where NewConstHistogram would have returned an error. 

####  func [MustNewConstHistogramWithCreatedTimestamp](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/histogram.go#L1427) ¶ added in v1.20.0
    
    
    func MustNewConstHistogramWithCreatedTimestamp(
    	desc *Desc,
    	count [uint64](/builtin#uint64),
    	sum [float64](/builtin#float64),
    	buckets map[[float64](/builtin#float64)][uint64](/builtin#uint64),
    	ct [time](/time).[Time](/time#Time),
    	labelValues ...[string](/builtin#string),
    ) Metric

MustNewConstHistogramWithCreatedTimestamp is a version of NewConstHistogramWithCreatedTimestamp that panics where NewConstHistogramWithCreatedTimestamp would have returned an error. 

####  func [MustNewConstMetric](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/value.go#L126) ¶
    
    
    func MustNewConstMetric(desc *Desc, valueType ValueType, value [float64](/builtin#float64), labelValues ...[string](/builtin#string)) Metric

MustNewConstMetric is a version of NewConstMetric that panics where NewConstMetric would have returned an error. 

####  func [MustNewConstMetricWithCreatedTimestamp](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/value.go#L163) ¶ added in v1.17.0
    
    
    func MustNewConstMetricWithCreatedTimestamp(desc *Desc, valueType ValueType, value [float64](/builtin#float64), ct [time](/time).[Time](/time#Time), labelValues ...[string](/builtin#string)) Metric

MustNewConstMetricWithCreatedTimestamp is a version of NewConstMetricWithCreatedTimestamp that panics where NewConstMetricWithCreatedTimestamp would have returned an error. 

####  func [MustNewConstNativeHistogram](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/histogram.go#L1969) ¶ added in v1.21.0
    
    
    func MustNewConstNativeHistogram(
    	desc *Desc,
    	count [uint64](/builtin#uint64),
    	sum [float64](/builtin#float64),
    	positiveBuckets, negativeBuckets map[[int](/builtin#int)][int64](/builtin#int64),
    	zeroBucket [uint64](/builtin#uint64),
    	nativeHistogramSchema [int32](/builtin#int32),
    	nativeHistogramZeroThreshold [float64](/builtin#float64),
    	createdTimestamp [time](/time).[Time](/time#Time),
    	labelValues ...[string](/builtin#string),
    ) Metric

MustNewConstNativeHistogram is a version of NewConstNativeHistogram that panics where NewConstNativeHistogram would have returned an error. 

####  func [MustNewConstSummary](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/summary.go#L776) ¶
    
    
    func MustNewConstSummary(
    	desc *Desc,
    	count [uint64](/builtin#uint64),
    	sum [float64](/builtin#float64),
    	quantiles map[[float64](/builtin#float64)][float64](/builtin#float64),
    	labelValues ...[string](/builtin#string),
    ) Metric

MustNewConstSummary is a version of NewConstSummary that panics where NewConstMetric would have returned an error. 

####  func [MustNewConstSummaryWithCreatedTimestamp](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/summary.go#L817) ¶ added in v1.20.0
    
    
    func MustNewConstSummaryWithCreatedTimestamp(
    	desc *Desc,
    	count [uint64](/builtin#uint64),
    	sum [float64](/builtin#float64),
    	quantiles map[[float64](/builtin#float64)][float64](/builtin#float64),
    	ct [time](/time).[Time](/time#Time),
    	labelValues ...[string](/builtin#string),
    ) Metric

MustNewConstSummaryWithCreatedTimestamp is a version of NewConstSummaryWithCreatedTimestamp that panics where NewConstSummaryWithCreatedTimestamp would have returned an error. 

####  func [MustNewMetricWithExemplars](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/metric.go#L270) ¶ added in v1.13.0
    
    
    func MustNewMetricWithExemplars(m Metric, exemplars ...Exemplar) Metric

MustNewMetricWithExemplars is a version of NewMetricWithExemplars that panics where NewMetricWithExemplars would have returned an error. 

####  func [NewConstHistogram](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/histogram.go#L1362) ¶
    
    
    func NewConstHistogram(
    	desc *Desc,
    	count [uint64](/builtin#uint64),
    	sum [float64](/builtin#float64),
    	buckets map[[float64](/builtin#float64)][uint64](/builtin#uint64),
    	labelValues ...[string](/builtin#string),
    ) (Metric, [error](/builtin#error))

NewConstHistogram returns a metric representing a Prometheus histogram with fixed values for the count, sum, and bucket counts. As those parameters cannot be changed, the returned value does not implement the Histogram interface (but only the Metric interface). Users of this package will not have much use for it in regular operations. However, when implementing custom Collectors, it is useful as a throw-away metric that is generated on the fly to send it to Prometheus in the Collect method. 

buckets is a map of upper bounds to cumulative counts, excluding the +Inf bucket. The +Inf bucket is implicit, and its value is equal to the provided count. 

NewConstHistogram returns an error if the length of labelValues is not consistent with the variable labels in Desc or if Desc is invalid. 

Example ¶
    
    
    desc := prometheus.NewDesc(
    	"http_request_duration_seconds",
    	"A histogram of the HTTP request durations.",
    	[]string{"code", "method"},
    	prometheus.Labels{"owner": "example"},
    )
    
    // Create a constant histogram from values we got from a 3rd party telemetry system.
    h := prometheus.MustNewConstHistogram(
    	desc,
    	4711, 403.34,
    	map[float64]uint64{25: 121, 50: 2403, 100: 3221, 200: 4233},
    	"200", "get",
    )
    
    // Just for demonstration, let's check the state of the histogram by
    // (ab)using its Write method (which is usually only used by Prometheus
    // internally).
    metric := &dto.Metric{}
    h.Write(metric)
    fmt.Println(toNormalizedJSON(metric))
    
    
    
    Output:
    
    {"label":[{"name":"code","value":"200"},{"name":"method","value":"get"},{"name":"owner","value":"example"}],"histogram":{"sampleCount":"4711","sampleSum":403.34,"bucket":[{"cumulativeCount":"121","upperBound":25},{"cumulativeCount":"2403","upperBound":50},{"cumulativeCount":"3221","upperBound":100},{"cumulativeCount":"4233","upperBound":200}]}}
    

Example (WithExemplar) ¶
    
    
    desc := prometheus.NewDesc(
    	"http_request_duration_seconds",
    	"A histogram of the HTTP request durations.",
    	[]string{"code", "method"},
    	prometheus.Labels{"owner": "example"},
    )
    
    // Create a constant histogram from values we got from a 3rd party telemetry system.
    h := prometheus.MustNewConstHistogram(
    	desc,
    	4711, 403.34,
    	map[float64]uint64{25: 121, 50: 2403, 100: 3221, 200: 4233},
    	"200", "get",
    )
    
    // Wrap const histogram with exemplars for each bucket.
    exemplarTs, _ := time.Parse(time.RFC850, "Monday, 02-Jan-06 15:04:05 GMT")
    exemplarLabels := prometheus.Labels{"testName": "testVal"}
    h = prometheus.MustNewMetricWithExemplars(
    	h,
    	prometheus.Exemplar{Labels: exemplarLabels, Timestamp: exemplarTs, Value: 24.0},
    	prometheus.Exemplar{Labels: exemplarLabels, Timestamp: exemplarTs, Value: 42.0},
    	prometheus.Exemplar{Labels: exemplarLabels, Timestamp: exemplarTs, Value: 89.0},
    	prometheus.Exemplar{Labels: exemplarLabels, Timestamp: exemplarTs, Value: 157.0},
    )
    
    // Just for demonstration, let's check the state of the histogram by
    // (ab)using its Write method (which is usually only used by Prometheus
    // internally).
    metric := &dto.Metric{}
    h.Write(metric)
    fmt.Println(toNormalizedJSON(metric))
    
    
    
    Output:
    
    {"label":[{"name":"code","value":"200"},{"name":"method","value":"get"},{"name":"owner","value":"example"}],"histogram":{"sampleCount":"4711","sampleSum":403.34,"bucket":[{"cumulativeCount":"121","upperBound":25,"exemplar":{"label":[{"name":"testName","value":"testVal"}],"value":24,"timestamp":"2006-01-02T15:04:05Z"}},{"cumulativeCount":"2403","upperBound":50,"exemplar":{"label":[{"name":"testName","value":"testVal"}],"value":42,"timestamp":"2006-01-02T15:04:05Z"}},{"cumulativeCount":"3221","upperBound":100,"exemplar":{"label":[{"name":"testName","value":"testVal"}],"value":89,"timestamp":"2006-01-02T15:04:05Z"}},{"cumulativeCount":"4233","upperBound":200,"exemplar":{"label":[{"name":"testName","value":"testVal"}],"value":157,"timestamp":"2006-01-02T15:04:05Z"}}]}}
    

####  func [NewConstHistogramWithCreatedTimestamp](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/histogram.go#L1401) ¶ added in v1.20.0
    
    
    func NewConstHistogramWithCreatedTimestamp(
    	desc *Desc,
    	count [uint64](/builtin#uint64),
    	sum [float64](/builtin#float64),
    	buckets map[[float64](/builtin#float64)][uint64](/builtin#uint64),
    	ct [time](/time).[Time](/time#Time),
    	labelValues ...[string](/builtin#string),
    ) (Metric, [error](/builtin#error))

NewConstHistogramWithCreatedTimestamp does the same thing as NewConstHistogram but sets the created timestamp. 

Example ¶
    
    
    desc := prometheus.NewDesc(
    	"http_request_duration_seconds",
    	"A histogram of the HTTP request durations.",
    	[]string{"code", "method"},
    	prometheus.Labels{"owner": "example"},
    )
    
    createdTs := time.Unix(1719670764, 123)
    h := prometheus.MustNewConstHistogramWithCreatedTimestamp(
    	desc,
    	4711, 403.34,
    	map[float64]uint64{25: 121, 50: 2403, 100: 3221, 200: 4233},
    	createdTs,
    	"200", "get",
    )
    
    // Just for demonstration, let's check the state of the histogram by
    // (ab)using its Write method (which is usually only used by Prometheus
    // internally).
    metric := &dto.Metric{}
    h.Write(metric)
    fmt.Println(toNormalizedJSON(metric))
    
    
    
    Output:
    
    {"label":[{"name":"code","value":"200"},{"name":"method","value":"get"},{"name":"owner","value":"example"}],"histogram":{"sampleCount":"4711","sampleSum":403.34,"bucket":[{"cumulativeCount":"121","upperBound":25},{"cumulativeCount":"2403","upperBound":50},{"cumulativeCount":"3221","upperBound":100},{"cumulativeCount":"4233","upperBound":200}],"createdTimestamp":"2024-06-29T14:19:24.000000123Z"}}
    

####  func [NewConstMetric](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/value.go#L105) ¶
    
    
    func NewConstMetric(desc *Desc, valueType ValueType, value [float64](/builtin#float64), labelValues ...[string](/builtin#string)) (Metric, [error](/builtin#error))

NewConstMetric returns a metric with one fixed value that cannot be changed. Users of this package will not have much use for it in regular operations. However, when implementing custom Collectors, it is useful as a throw-away metric that is generated on the fly to send it to Prometheus in the Collect method. NewConstMetric returns an error if the length of labelValues is not consistent with the variable labels in Desc or if Desc is invalid. 

####  func [NewConstMetricWithCreatedTimestamp](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/value.go#L136) ¶ added in v1.17.0
    
    
    func NewConstMetricWithCreatedTimestamp(desc *Desc, valueType ValueType, value [float64](/builtin#float64), ct [time](/time).[Time](/time#Time), labelValues ...[string](/builtin#string)) (Metric, [error](/builtin#error))

NewConstMetricWithCreatedTimestamp does the same thing as NewConstMetric, but generates Counters with created timestamp set and returns an error for other metric types. 

Example ¶
    
    
    // Here we have a metric that is reported by an external system.
    // Besides providing the value, the external system also provides the
    // timestamp when the metric was created.
    desc := prometheus.NewDesc(
    	"time_since_epoch_seconds",
    	"Current epoch time in seconds.",
    	nil, nil,
    )
    
    timeSinceEpochReportedByExternalSystem := time.Date(2009, time.November, 10, 23, 0, 0, 12345678, time.UTC)
    epoch := time.Unix(0, 0).UTC()
    s := prometheus.MustNewConstMetricWithCreatedTimestamp(
    	desc, prometheus.CounterValue, float64(timeSinceEpochReportedByExternalSystem.Unix()), epoch,
    )
    
    metric := &dto.Metric{}
    s.Write(metric)
    fmt.Println(toNormalizedJSON(metric))
    
    
    
    Output:
    
    {"counter":{"value":1257894000,"createdTimestamp":"1970-01-01T00:00:00Z"}}
    

####  func [NewConstNativeHistogram](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/histogram.go#L1913) ¶ added in v1.21.0
    
    
    func NewConstNativeHistogram(
    	desc *Desc,
    	count [uint64](/builtin#uint64),
    	sum [float64](/builtin#float64),
    	positiveBuckets, negativeBuckets map[[int](/builtin#int)][int64](/builtin#int64),
    	zeroBucket [uint64](/builtin#uint64),
    	schema [int32](/builtin#int32),
    	zeroThreshold [float64](/builtin#float64),
    	createdTimestamp [time](/time).[Time](/time#Time),
    	labelValues ...[string](/builtin#string),
    ) (Metric, [error](/builtin#error))

NewConstNativeHistogram returns a metric representing a Prometheus native histogram with fixed values for the count, sum, and positive/negative/zero bucket counts. As those parameters cannot be changed, the returned value does not implement the Histogram interface (but only the Metric interface). Users of this package will not have much use for it in regular operations. However, when implementing custom OpenTelemetry Collectors, it is useful as a throw-away metric that is generated on the fly to send it to Prometheus in the Collect method. 

zeroBucket counts all (positive and negative) observations in the zero bucket (with an absolute value less or equal the current threshold). positiveBuckets and negativeBuckets are separate maps for negative and positive observations. The map's value is an int64, counting observations in that bucket. The map's key is the index of the bucket according to the used Schema. Index 0 is for an upper bound of 1 in positive buckets and for a lower bound of -1 in negative buckets. NewConstNativeHistogram returns an error if 

  * the length of labelValues is not consistent with the variable labels in Desc or if Desc is invalid.
  * the schema passed is not between 8 and -4
  * the sum of counts in all buckets including the zero bucket does not equal the count if sum is not NaN (or exceeds the count if sum is NaN)



See <https://opentelemetry.io/docs/specs/otel/compatibility/prometheus_and_openmetrics/#exponential-histograms> for more details about the conversion from OTel to Prometheus. 

####  func [NewConstSummary](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/summary.go#L752) ¶
    
    
    func NewConstSummary(
    	desc *Desc,
    	count [uint64](/builtin#uint64),
    	sum [float64](/builtin#float64),
    	quantiles map[[float64](/builtin#float64)][float64](/builtin#float64),
    	labelValues ...[string](/builtin#string),
    ) (Metric, [error](/builtin#error))

NewConstSummary returns a metric representing a Prometheus summary with fixed values for the count, sum, and quantiles. As those parameters cannot be changed, the returned value does not implement the Summary interface (but only the Metric interface). Users of this package will not have much use for it in regular operations. However, when implementing custom Collectors, it is useful as a throw-away metric that is generated on the fly to send it to Prometheus in the Collect method. 

quantiles maps ranks to quantile values. For example, a median latency of 0.23s and a 99th percentile latency of 0.56s would be expressed as: 
    
    
    map[float64]float64{0.5: 0.23, 0.99: 0.56}
    

NewConstSummary returns an error if the length of labelValues is not consistent with the variable labels in Desc or if Desc is invalid. 

Example ¶
    
    
    desc := prometheus.NewDesc(
    	"http_request_duration_seconds",
    	"A summary of the HTTP request durations.",
    	[]string{"code", "method"},
    	prometheus.Labels{"owner": "example"},
    )
    
    // Create a constant summary from values we got from a 3rd party telemetry system.
    s := prometheus.MustNewConstSummary(
    	desc,
    	4711, 403.34,
    	map[float64]float64{0.5: 42.3, 0.9: 323.3},
    	"200", "get",
    )
    
    // Just for demonstration, let's check the state of the summary by
    // (ab)using its Write method (which is usually only used by Prometheus
    // internally).
    metric := &dto.Metric{}
    s.Write(metric)
    fmt.Println(toNormalizedJSON(metric))
    
    
    
    Output:
    
    {"label":[{"name":"code","value":"200"},{"name":"method","value":"get"},{"name":"owner","value":"example"}],"summary":{"sampleCount":"4711","sampleSum":403.34,"quantile":[{"quantile":0.5,"value":42.3},{"quantile":0.9,"value":323.3}]}}
    

####  func [NewConstSummaryWithCreatedTimestamp](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/summary.go#L791) ¶ added in v1.20.0
    
    
    func NewConstSummaryWithCreatedTimestamp(
    	desc *Desc,
    	count [uint64](/builtin#uint64),
    	sum [float64](/builtin#float64),
    	quantiles map[[float64](/builtin#float64)][float64](/builtin#float64),
    	ct [time](/time).[Time](/time#Time),
    	labelValues ...[string](/builtin#string),
    ) (Metric, [error](/builtin#error))

NewConstSummaryWithCreatedTimestamp does the same thing as NewConstSummary but sets the created timestamp. 

Example ¶
    
    
    desc := prometheus.NewDesc(
    	"http_request_duration_seconds",
    	"A summary of the HTTP request durations.",
    	[]string{"code", "method"},
    	prometheus.Labels{"owner": "example"},
    )
    
    // Create a constant summary with created timestamp set
    createdTs := time.Unix(1719670764, 123)
    s := prometheus.MustNewConstSummaryWithCreatedTimestamp(
    	desc,
    	4711, 403.34,
    	map[float64]float64{0.5: 42.3, 0.9: 323.3},
    	createdTs,
    	"200", "get",
    )
    
    // Just for demonstration, let's check the state of the summary by
    // (ab)using its Write method (which is usually only used by Prometheus
    // internally).
    metric := &dto.Metric{}
    s.Write(metric)
    fmt.Println(toNormalizedJSON(metric))
    
    
    
    Output:
    
    {"label":[{"name":"code","value":"200"},{"name":"method","value":"get"},{"name":"owner","value":"example"}],"summary":{"sampleCount":"4711","sampleSum":403.34,"quantile":[{"quantile":0.5,"value":42.3},{"quantile":0.9,"value":323.3}],"createdTimestamp":"2024-06-29T14:19:24.000000123Z"}}
    

####  func [NewInvalidMetric](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/metric.go#L138) ¶
    
    
    func NewInvalidMetric(desc *Desc, err [error](/builtin#error)) Metric

NewInvalidMetric returns a metric whose Write method always returns the provided error. It is useful if a Collector finds itself unable to collect a metric and wishes to report an error to the registry. 

####  func [NewMetricWithExemplars](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/metric.go#L244) ¶ added in v1.13.0
    
    
    func NewMetricWithExemplars(m Metric, exemplars ...Exemplar) (Metric, [error](/builtin#error))

NewMetricWithExemplars returns a new Metric wrapping the provided Metric with given exemplars. Exemplars are validated. 

Only last applicable exemplar is injected from the list. For example for Counter it means last exemplar is injected. For Histogram, it means last applicable exemplar for each bucket is injected. For a Native Histogram, all valid exemplars are injected. 

NewMetricWithExemplars works best with MustNewConstMetric and MustNewConstHistogram, see example. 

####  func [NewMetricWithTimestamp](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/metric.go#L170) ¶ added in v0.9.0
    
    
    func NewMetricWithTimestamp(t [time](/time).[Time](/time#Time), m Metric) Metric

NewMetricWithTimestamp returns a new Metric wrapping the provided Metric in a way that it has an explicit timestamp set to the provided Time. This is only useful in rare cases as the timestamp of a Prometheus metric should usually be set by the Prometheus server during scraping. Exceptions include mirroring metrics with given timestamps from other metric sources. 

NewMetricWithTimestamp works best with MustNewConstMetric, MustNewConstHistogram, and MustNewConstSummary, see example. 

Currently, the exposition formats used by Prometheus are limited to millisecond resolution. Thus, the provided time will be rounded down to the next full millisecond value. 

Example ¶
    
    
    desc := prometheus.NewDesc(
    	"temperature_kelvin",
    	"Current temperature in Kelvin.",
    	nil, nil,
    )
    
    // Create a constant gauge from values we got from an external
    // temperature reporting system. Those values are reported with a slight
    // delay, so we want to add the timestamp of the actual measurement.
    temperatureReportedByExternalSystem := 298.15
    timeReportedByExternalSystem := time.Date(2009, time.November, 10, 23, 0, 0, 12345678, time.UTC)
    s := prometheus.NewMetricWithTimestamp(
    	timeReportedByExternalSystem,
    	prometheus.MustNewConstMetric(
    		desc, prometheus.GaugeValue, temperatureReportedByExternalSystem,
    	),
    )
    
    // Just for demonstration, let's check the state of the gauge by
    // (ab)using its Write method (which is usually only used by Prometheus
    // internally).
    metric := &dto.Metric{}
    s.Write(metric)
    fmt.Println(toNormalizedJSON(metric))
    
    
    
    Output:
    
    {"gauge":{"value":298.15},"timestampMs":"1257894000012"}
    

####  type [MetricVec](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/vec.go#L36) ¶
    
    
    type MetricVec struct {
    	// contains filtered or unexported fields
    }

MetricVec is a Collector to bundle metrics of the same name that differ in their label values. MetricVec is not used directly but as a building block for implementations of vectors of a given metric type, like GaugeVec, CounterVec, SummaryVec, and HistogramVec. It is exported so that it can be used for custom Metric implementations. 

To create a FooVec for custom Metric Foo, embed a pointer to MetricVec in FooVec and initialize it with NewMetricVec. Implement wrappers for GetMetricWithLabelValues and GetMetricWith that return (Foo, error) rather than (Metric, error). Similarly, create a wrapper for CurryWith that returns (*FooVec, error) rather than (*MetricVec, error). It is recommended to also add the convenience methods WithLabelValues, With, and MustCurryWith, which panic instead of returning errors. See also the MetricVec example. 

Example ¶
    
    
    package main
    
    import (
    	"fmt"
    
    	"google.golang.org/protobuf/proto"
    
    	dto "github.com/prometheus/client_model/go"
    
    	"github.com/prometheus/client_golang/prometheus"
    )
    
    // Info implements an info pseudo-metric, which is modeled as a Gauge that
    // always has a value of 1. In practice, you would just use a Gauge directly,
    // but for this example, we pretend it would be useful to have a “native”
    // implementation.
    type Info struct {
    	desc       *prometheus.Desc
    	labelPairs []*dto.LabelPair
    }
    
    func (i Info) Desc() *prometheus.Desc {
    	return i.desc
    }
    
    func (i Info) Write(out *dto.Metric) error {
    	out.Label = i.labelPairs
    	out.Gauge = &dto.Gauge{Value: proto.Float64(1)}
    	return nil
    }
    
    // InfoVec is the vector version for Info. As an info metric never changes, we
    // wouldn't really need to wrap GetMetricWithLabelValues and GetMetricWith
    // because Info has no additional methods compared to the vanilla Metric that
    // the unwrapped MetricVec methods return. However, to demonstrate all there is
    // to do to fully implement a vector for a custom Metric implementation, we do
    // it in this example anyway.
    type InfoVec struct {
    	*prometheus.MetricVec
    }
    
    func NewInfoVec(name, help string, labelNames []string) *InfoVec {
    	desc := prometheus.NewDesc(name, help, labelNames, nil)
    	return &InfoVec{
    		MetricVec: prometheus.NewMetricVec(desc, func(lvs ...string) prometheus.Metric {
    			if len(lvs) != len(labelNames) {
    				panic("inconsistent label cardinality")
    			}
    			return Info{desc: desc, labelPairs: prometheus.MakeLabelPairs(desc, lvs)}
    		}),
    	}
    }
    
    func (v *InfoVec) GetMetricWithLabelValues(lvs ...string) (Info, error) {
    	metric, err := v.MetricVec.GetMetricWithLabelValues(lvs...)
    	return metric.(Info), err
    }
    
    func (v *InfoVec) GetMetricWith(labels prometheus.Labels) (Info, error) {
    	metric, err := v.MetricVec.GetMetricWith(labels)
    	return metric.(Info), err
    }
    
    func (v *InfoVec) WithLabelValues(lvs ...string) Info {
    	i, err := v.GetMetricWithLabelValues(lvs...)
    	if err != nil {
    		panic(err)
    	}
    	return i
    }
    
    func (v *InfoVec) With(labels prometheus.Labels) Info {
    	i, err := v.GetMetricWith(labels)
    	if err != nil {
    		panic(err)
    	}
    	return i
    }
    
    func (v *InfoVec) CurryWith(labels prometheus.Labels) (*InfoVec, error) {
    	vec, err := v.MetricVec.CurryWith(labels)
    	if vec != nil {
    		return &InfoVec{vec}, err
    	}
    	return nil, err
    }
    
    func (v *InfoVec) MustCurryWith(labels prometheus.Labels) *InfoVec {
    	vec, err := v.CurryWith(labels)
    	if err != nil {
    		panic(err)
    	}
    	return vec
    }
    
    func main() {
    	infoVec := NewInfoVec(
    		"library_version_info",
    		"Versions of the libraries used in this binary.",
    		[]string{"library", "version"},
    	)
    
    	infoVec.WithLabelValues("prometheus/client_golang", "1.7.1")
    	infoVec.WithLabelValues("k8s.io/client-go", "0.18.8")
    
    	// Just for demonstration, let's check the state of the InfoVec by
    	// registering it with a custom registry and then let it collect the
    	// metrics.
    	reg := prometheus.NewRegistry()
    	reg.MustRegister(infoVec)
    
    	metricFamilies, err := reg.Gather()
    	if err != nil || len(metricFamilies) != 1 {
    		panic("unexpected behavior of custom test registry")
    	}
    	fmt.Println(toNormalizedJSON(metricFamilies[0]))
    
    }
    
    
    
    Output:
    
    {"name":"library_version_info","help":"Versions of the libraries used in this binary.","type":"GAUGE","metric":[{"label":[{"name":"library","value":"k8s.io/client-go"},{"name":"version","value":"0.18.8"}],"gauge":{"value":1}},{"label":[{"name":"library","value":"prometheus/client_golang"},{"name":"version","value":"1.7.1"}],"gauge":{"value":1}}]}
    

Share Format Run

####  func [NewMetricVec](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/vec.go#L47) ¶ added in v0.12.1
    
    
    func NewMetricVec(desc *Desc, newMetric func(lvs ...[string](/builtin#string)) Metric) *MetricVec

NewMetricVec returns an initialized metricVec. 

####  func (*MetricVec) [Collect](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/vec.go#L127) ¶
    
    
    func (m *MetricVec) Collect(ch chan<- Metric)

Collect implements Collector. 

####  func (*MetricVec) [CurryWith](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/vec.go#L149) ¶ added in v0.12.1
    
    
    func (m *MetricVec) CurryWith(labels Labels) (*MetricVec, [error](/builtin#error))

CurryWith returns a vector curried with the provided labels, i.e. the returned vector has those labels pre-set for all labeled operations performed on it. The cardinality of the curried vector is reduced accordingly. The order of the remaining labels stays the same (just with the curried labels taken out of the sequence – which is relevant for the (GetMetric)WithLabelValues methods). It is possible to curry a curried vector, but only with labels not yet used for currying before. 

The metrics contained in the MetricVec are shared between the curried and uncurried vectors. They are just accessed differently. Curried and uncurried vectors behave identically in terms of collection. Only one must be registered with a given registry (usually the uncurried version). The Reset method deletes all metrics, even if called on a curried vector. 

Note that CurryWith is usually not called directly but through a wrapper around MetricVec, implementing a vector for a specific Metric implementation, for example GaugeVec. 

####  func (*MetricVec) [Delete](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/vec.go#L95) ¶
    
    
    func (m *MetricVec) Delete(labels Labels) [bool](/builtin#bool)

Delete deletes the metric where the variable labels are the same as those passed in as labels. It returns true if a metric was deleted. 

It is not an error if the number and names of the Labels are inconsistent with those of the VariableLabels in Desc. However, such inconsistent Labels can never match an actual metric, so the method will always return false in that case. 

This method is used for the same purpose as DeleteLabelValues(...string). See there for pros and cons of the two methods. 

####  func (*MetricVec) [DeleteLabelValues](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/vec.go#L74) ¶
    
    
    func (m *MetricVec) DeleteLabelValues(lvs ...[string](/builtin#string)) [bool](/builtin#bool)

DeleteLabelValues removes the metric where the variable labels are the same as those passed in as labels (same order as the VariableLabels in Desc). It returns true if a metric was deleted. 

It is not an error if the number of label values is not the same as the number of VariableLabels in Desc. However, such inconsistent label count can never match an actual metric, so the method will always return false in that case. 

Note that for more than one label value, this method is prone to mistakes caused by an incorrect order of arguments. Consider Delete(Labels) as an alternative to avoid that type of mistake. For higher label numbers, the latter has a much more readable (albeit more verbose) syntax, but it comes with a performance overhead (for creating and processing the Labels map). See also the CounterVec example. 

####  func (*MetricVec) [DeletePartialMatch](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/vec.go#L113) ¶ added in v1.13.0
    
    
    func (m *MetricVec) DeletePartialMatch(labels Labels) [int](/builtin#int)

DeletePartialMatch deletes all metrics where the variable labels contain all of those passed in as labels. The order of the labels does not matter. It returns the number of metrics deleted. 

Note that curried labels will never be matched if deleting from the curried vector. To match curried labels with DeletePartialMatch, it must be called on the base vector. 

####  func (*MetricVec) [Describe](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/vec.go#L124) ¶
    
    
    func (m *MetricVec) Describe(ch chan<- *Desc)

Describe implements Collector. 

####  func (*MetricVec) [GetMetricWith](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/vec.go#L238) ¶
    
    
    func (m *MetricVec) GetMetricWith(labels Labels) (Metric, [error](/builtin#error))

GetMetricWith returns the Metric for the given Labels map (the label names must match those of the variable labels in Desc). If that label map is accessed for the first time, a new Metric is created. Implications of creating a Metric without using it and keeping the Metric for later use are the same as for GetMetricWithLabelValues. 

An error is returned if the number and names of the Labels are inconsistent with those of the variable labels in Desc (minus any curried labels). 

This method is used for the same purpose as GetMetricWithLabelValues(...string). See there for pros and cons of the two methods. 

Note that GetMetricWith is usually not called directly but through a wrapper around MetricVec, implementing a vector for a specific Metric implementation, for example GaugeVec. 

####  func (*MetricVec) [GetMetricWithLabelValues](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/vec.go#L212) ¶
    
    
    func (m *MetricVec) GetMetricWithLabelValues(lvs ...[string](/builtin#string)) (Metric, [error](/builtin#error))

GetMetricWithLabelValues returns the Metric for the given slice of label values (same order as the variable labels in Desc). If that combination of label values is accessed for the first time, a new Metric is created (by calling the newMetric function provided during construction of the MetricVec). 

It is possible to call this method without using the returned Metric to only create the new Metric but leave it in its initial state. 

Keeping the Metric for later use is possible (and should be considered if performance is critical), but keep in mind that Reset, DeleteLabelValues and Delete can be used to delete the Metric from the MetricVec. In that case, the Metric will still exist, but it will not be exported anymore, even if a Metric with the same label values is created later. 

An error is returned if the number of label values is not the same as the number of variable labels in Desc (minus any curried labels). 

Note that for more than one label value, this method is prone to mistakes caused by an incorrect order of arguments. Consider GetMetricWith(Labels) as an alternative to avoid that type of mistake. For higher label numbers, the latter has a much more readable (albeit more verbose) syntax, but it comes with a performance overhead (for creating and processing the Labels map). 

Note that GetMetricWithLabelValues is usually not called directly but through a wrapper around MetricVec, implementing a vector for a specific Metric implementation, for example GaugeVec. 

####  func (*MetricVec) [Reset](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/vec.go#L130) ¶
    
    
    func (m *MetricVec) Reset()

Reset deletes all metrics in this vector. 

####  type [MultiError](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/registry.go#L215) ¶
    
    
    type MultiError [][error](/builtin#error)

MultiError is a slice of errors implementing the error interface. It is used by a Gatherer to report multiple errors during MetricFamily gathering. 

####  func (*MultiError) [Append](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/registry.go#L232) ¶ added in v0.9.0
    
    
    func (errs *MultiError) Append(err [error](/builtin#error))

Append appends the provided error if it is not nil. 

####  func (MultiError) [Error](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/registry.go#L219) ¶
    
    
    func (errs MultiError) Error() [string](/builtin#string)

Error formats the contained errors as a bullet point list, preceded by the total number of errors. Note that this results in a multi-line string. 

####  func (MultiError) [MaybeUnwrap](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/registry.go#L242) ¶
    
    
    func (errs MultiError) MaybeUnwrap() [error](/builtin#error)

MaybeUnwrap returns nil if len(errs) is 0. It returns the first and only contained error as error if len(errs is 1). In all other cases, it returns the MultiError directly. This is helpful for returning a MultiError in a way that only uses the MultiError if needed. 

####  type [MultiTRegistry](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/registry.go#L999) ¶ added in v1.13.0
    
    
    type MultiTRegistry struct {
    	// contains filtered or unexported fields
    }

MultiTRegistry is a TransactionalGatherer that joins gathered metrics from multiple transactional gatherers. 

It is caller responsibility to ensure two registries have mutually exclusive metric families, no deduplication will happen. 

####  func [NewMultiTRegistry](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/registry.go#L1004) ¶ added in v1.13.0
    
    
    func NewMultiTRegistry(tGatherers ...TransactionalGatherer) *MultiTRegistry

NewMultiTRegistry creates MultiTRegistry. 

####  func (*MultiTRegistry) [Gather](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/registry.go#L1011) ¶ added in v1.13.0
    
    
    func (r *MultiTRegistry) Gather() (mfs []*[dto](/github.com/prometheus/client_model/go).[MetricFamily](/github.com/prometheus/client_model/go#MetricFamily), done func(), err [error](/builtin#error))

Gather implements TransactionalGatherer interface. 

####  type [Observer](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/observer.go#L18) ¶ added in v0.9.0
    
    
    type Observer interface {
    	Observe([float64](/builtin#float64))
    }

Observer is the interface that wraps the Observe method, which is used by Histogram and Summary to add observations. 

####  type [ObserverFunc](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/observer.go#L35) ¶ added in v0.9.0
    
    
    type ObserverFunc func([float64](/builtin#float64))

The ObserverFunc type is an adapter to allow the use of ordinary functions as Observers. If f is a function with the appropriate signature, ObserverFunc(f) is an Observer that calls f. 

This adapter is usually used in connection with the Timer type, and there are two general use cases: 

The most common one is to use a Gauge as the Observer for a Timer. See the "Gauge" Timer example. 

The more advanced use case is to create a function that dynamically decides which Observer to use for observing the duration. See the "Complex" Timer example. 

####  func (ObserverFunc) [Observe](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/observer.go#L38) ¶ added in v0.9.0
    
    
    func (f ObserverFunc) Observe(value [float64](/builtin#float64))

Observe calls f(value). It implements Observer. 

####  type [ObserverVec](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/observer.go#L43) ¶ added in v0.9.0
    
    
    type ObserverVec interface {
    	GetMetricWith(Labels) (Observer, [error](/builtin#error))
    	GetMetricWithLabelValues(lvs ...[string](/builtin#string)) (Observer, [error](/builtin#error))
    	With(Labels) Observer
    	WithLabelValues(...[string](/builtin#string)) Observer
    	CurryWith(Labels) (ObserverVec, [error](/builtin#error))
    	MustCurryWith(Labels) ObserverVec
    
    	Collector
    }

ObserverVec is an interface implemented by `HistogramVec` and `SummaryVec`. 

####  type [Opts](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/metric.go#L68) ¶
    
    
    type Opts struct {
    	// Namespace, Subsystem, and Name are components of the fully-qualified
    	// name of the Metric (created by joining these components with
    	// "_"). Only Name is mandatory, the others merely help structuring the
    	// name. Note that the fully-qualified name of the metric must be a
    	// valid Prometheus metric name.
    	Namespace [string](/builtin#string)
    	Subsystem [string](/builtin#string)
    	Name      [string](/builtin#string)
    
    	// Help provides information about this metric.
    	//
    	// Metrics with the same fully-qualified name must have the same Help
    	// string.
    	Help [string](/builtin#string)
    
    	// ConstLabels are used to attach fixed labels to this metric. Metrics
    	// with the same fully-qualified name must have the same label names in
    	// their ConstLabels.
    	//
    	// ConstLabels are only used rarely. In particular, do not use them to
    	// attach the same labels to all your metrics. Those use cases are
    	// better covered by target labels set by the scraping Prometheus
    	// server, or by one specific metric (e.g. a build_info or a
    	// machine_role metric). See also
    	// <https://prometheus.io/docs/instrumenting/writing_exporters/#target-labels-not-static-scraped-labels>
    	ConstLabels Labels
    	// contains filtered or unexported fields
    }

Opts bundles the options for creating most Metric types. Each metric implementation XXX has its own XXXOpts type, but in most cases, it is just an alias of this type (which might change when the requirement arises.) 

It is mandatory to set Name to a non-empty string. All other fields are optional and can safely be left at their zero value, although it is strongly encouraged to set a Help string. 

####  type [ProcessCollectorOpts](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/process_collector.go#L39) ¶ added in v0.9.0
    
    
    type ProcessCollectorOpts struct {
    	// PidFn returns the PID of the process the collector collects metrics
    	// for. It is called upon each collection. By default, the PID of the
    	// current process is used, as determined on construction time by
    	// calling os.Getpid().
    	PidFn func() ([int](/builtin#int), [error](/builtin#error))
    	// If non-empty, each of the collected metrics is prefixed by the
    	// provided string and an underscore ("_").
    	Namespace [string](/builtin#string)
    	// If true, any error encountered during collection is reported as an
    	// invalid metric (see NewInvalidMetric). Otherwise, errors are ignored
    	// and the collected metrics will be incomplete. (Possibly, no metrics
    	// will be collected at all.) While that's usually not desired, it is
    	// appropriate for the common "mix-in" of process metrics, where process
    	// metrics are nice to have, but failing to collect them should not
    	// disrupt the collection of the remaining metrics.
    	ReportErrors [bool](/builtin#bool)
    }

ProcessCollectorOpts defines the behavior of a process metrics collector created with NewProcessCollector. 

####  type [Registerer](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/registry.go#L96) ¶
    
    
    type Registerer interface {
    	// Register registers a new Collector to be included in metrics
    	// collection. It returns an error if the descriptors provided by the
    	// Collector are invalid or if they — in combination with descriptors of
    	// already registered Collectors — do not fulfill the consistency and
    	// uniqueness criteria described in the documentation of metric.Desc.
    	//
    	// If the provided Collector is equal to a Collector already registered
    	// (which includes the case of re-registering the same Collector), the
    	// returned error is an instance of AlreadyRegisteredError, which
    	// contains the previously registered Collector.
    	//
    	// A Collector whose Describe method does not yield any Desc is treated
    	// as unchecked. Registration will always succeed. No check for
    	// re-registering (see previous paragraph) is performed. Thus, the
    	// caller is responsible for not double-registering the same unchecked
    	// Collector, and for providing a Collector that will not cause
    	// inconsistent metrics on collection. (This would lead to scrape
    	// errors.)
    	Register(Collector) [error](/builtin#error)
    	// MustRegister works like Register but registers any number of
    	// Collectors and panics upon the first registration that causes an
    	// error.
    	MustRegister(...Collector)
    	// Unregister unregisters the Collector that equals the Collector passed
    	// in as an argument.  (Two Collectors are considered equal if their
    	// Describe method yields the same set of descriptors.) The function
    	// returns whether a Collector was unregistered. Note that an unchecked
    	// Collector cannot be unregistered (as its Describe method does not
    	// yield any descriptor).
    	//
    	// Note that even after unregistering, it will not be possible to
    	// register a new Collector that is inconsistent with the unregistered
    	// Collector, e.g. a Collector collecting metrics with the same name but
    	// a different help string. The rationale here is that the same registry
    	// instance must only collect consistent metrics throughout its
    	// lifetime.
    	Unregister(Collector) [bool](/builtin#bool)
    }

Registerer is the interface for the part of a registry in charge of registering and unregistering. Users of custom registries should use Registerer as type for registration purposes (rather than the Registry type directly). In that way, they are free to use custom Registerer implementation (e.g. for testing purposes). 

####  func [WrapRegistererWith](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/wrap.go#L46) ¶ added in v0.9.0
    
    
    func WrapRegistererWith(labels Labels, reg Registerer) Registerer

WrapRegistererWith returns a Registerer wrapping the provided Registerer. Collectors registered with the returned Registerer will be registered with the wrapped Registerer in a modified way. The modified Collector adds the provided Labels to all Metrics it collects (as ConstLabels). The Metrics collected by the unmodified Collector must not duplicate any of those labels. Wrapping a nil value is valid, resulting in a no-op Registerer. 

WrapRegistererWith provides a way to add fixed labels to a subset of Collectors. It should not be used to add fixed labels to all metrics exposed. See also <https://prometheus.io/docs/instrumenting/writing_exporters/#target-labels-not-static-scraped-labels>

Conflicts between Collectors registered through the original Registerer with Collectors registered through the wrapping Registerer will still be detected. Any AlreadyRegisteredError returned by the Register method of either Registerer will contain the ExistingCollector in the form it was provided to the respective registry. 

The Collector example demonstrates a use of WrapRegistererWith. 

####  func [WrapRegistererWithPrefix](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/wrap.go#L74) ¶ added in v0.9.0
    
    
    func WrapRegistererWithPrefix(prefix [string](/builtin#string), reg Registerer) Registerer

WrapRegistererWithPrefix returns a Registerer wrapping the provided Registerer. Collectors registered with the returned Registerer will be registered with the wrapped Registerer in a modified way. The modified Collector adds the provided prefix to the name of all Metrics it collects. Wrapping a nil value is valid, resulting in a no-op Registerer. 

WrapRegistererWithPrefix is useful to have one place to prefix all metrics of a sub-system. To make this work, register metrics of the sub-system with the wrapping Registerer returned by WrapRegistererWithPrefix. It is rarely useful to use the same prefix for all metrics exposed. In particular, do not prefix metric names that are standardized across applications, as that would break horizontal monitoring, for example the metrics provided by the Go collector (see NewGoCollector) and the process collector (see NewProcessCollector). (In fact, those metrics are already prefixed with "go_" or "process_", respectively.) 

Conflicts between Collectors registered through the original Registerer with Collectors registered through the wrapping Registerer will still be detected. Any AlreadyRegisteredError returned by the Register method of either Registerer will contain the ExistingCollector in the form it was provided to the respective registry. 

####  type [Registry](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/registry.go#L260) ¶
    
    
    type Registry struct {
    	// contains filtered or unexported fields
    }

Registry registers Prometheus collectors, collects their metrics, and gathers them into MetricFamilies for exposition. It implements Registerer, Gatherer, and Collector. The zero value is not usable. Create instances with NewRegistry or NewPedanticRegistry. 

Registry implements Collector to allow it to be used for creating groups of metrics. See the Grouping example for how this can be done. 

Example (Grouping) ¶

This example shows how to use multiple registries for registering and unregistering groups of metrics. 
    
    
    package main
    
    import (
    	"math/rand"
    	"strconv"
    	"time"
    
    	"github.com/prometheus/client_golang/prometheus"
    )
    
    func main() {
    	// Create a global registry.
    	globalReg := prometheus.NewRegistry()
    
    	// Spawn 10 workers, each of which will have their own group of metrics.
    	for i := 0; i < 10; i++ {
    		// Create a new registry for each worker, which acts as a group of
    		// worker-specific metrics.
    		workerReg := prometheus.NewRegistry()
    		globalReg.Register(workerReg)
    
    		go func(workerID int) {
    			// Once the worker is done, it can unregister itself.
    			defer globalReg.Unregister(workerReg)
    
    			workTime := prometheus.NewCounter(prometheus.CounterOpts{
    				Name: "worker_total_work_time_milliseconds",
    				ConstLabels: prometheus.Labels{
    					// Generate a label unique to this worker so its metric doesn't
    					// collide with the metrics from other workers.
    					"worker_id": strconv.Itoa(workerID),
    				},
    			})
    			workerReg.MustRegister(workTime)
    
    			start := time.Now()
    			time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
    			workTime.Add(float64(time.Since(start).Milliseconds()))
    		}(i)
    	}
    }
    

Share Format Run

####  func [NewPedanticRegistry](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/registry.go#L85) ¶
    
    
    func NewPedanticRegistry() *Registry

NewPedanticRegistry returns a registry that checks during collection if each collected Metric is consistent with its reported Desc, and if the Desc has actually been registered with the registry. Unchecked Collectors (those whose Describe method does not yield any descriptors) are excluded from the check. 

Usually, a Registry will be happy as long as the union of all collected Metrics is consistent and valid even if some metrics are not consistent with their own Desc or a Desc provided by their registered Collector. Well-behaved Collectors and Metrics will only provide consistent Descs. This Registry is useful to test the implementation of Collectors and Metrics. 

####  func [NewRegistry](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/registry.go#L67) ¶
    
    
    func NewRegistry() *Registry

NewRegistry creates a new vanilla Registry without any Collectors pre-registered. 

####  func (*Registry) [Collect](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/registry.go#L575) ¶ added in v1.14.0
    
    
    func (r *Registry) Collect(ch chan<- Metric)

Collect implements Collector. 

####  func (*Registry) [Describe](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/registry.go#L563) ¶ added in v1.14.0
    
    
    func (r *Registry) Describe(ch chan<- *Desc)

Describe implements Collector. 

####  func (*Registry) [Gather](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/registry.go#L412) ¶
    
    
    func (r *Registry) Gather() ([]*[dto](/github.com/prometheus/client_model/go).[MetricFamily](/github.com/prometheus/client_model/go#MetricFamily), [error](/builtin#error))

Gather implements Gatherer. 

####  func (*Registry) [MustRegister](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/registry.go#L403) ¶
    
    
    func (r *Registry) MustRegister(cs ...Collector)

MustRegister implements Registerer. 

####  func (*Registry) [Register](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/registry.go#L270) ¶
    
    
    func (r *Registry) Register(c Collector) [error](/builtin#error)

Register implements Registerer. 

####  func (*Registry) [Unregister](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/registry.go#L366) ¶
    
    
    func (r *Registry) Unregister(c Collector) [bool](/builtin#bool)

Unregister implements Registerer. 

####  type [Summary](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/summary.go#L54) ¶
    
    
    type Summary interface {
    	Metric
    	Collector
    
    	// Observe adds a single observation to the summary. Observations are
    	// usually positive or zero. Negative observations are accepted but
    	// prevent current versions of Prometheus from properly detecting
    	// counter resets in the sum of observations. See
    	// <https://prometheus.io/docs/practices/histograms/#count-and-sum-of-observations>
    	// for details.
    	Observe([float64](/builtin#float64))
    }

A Summary captures individual observations from an event or sample stream and summarizes them in a manner similar to traditional summary statistics: 1. sum of observations, 2. observation count, 3. rank estimations. 

A typical use-case is the observation of request latencies. By default, a Summary provides the median, the 90th and the 99th percentile of the latency as rank estimations. However, the default behavior will change in the upcoming v1.0.0 of the library. There will be no rank estimations at all by default. For a sane transition, it is recommended to set the desired rank estimations explicitly. 

Note that the rank estimations cannot be aggregated in a meaningful way with the Prometheus query language (i.e. you cannot average or add them). If you need aggregatable quantiles (e.g. you want the 99th percentile latency of all queries served across all instances of a service), consider the Histogram metric type. See the Prometheus documentation for more details. 

To create Summary instances, use NewSummary. 

Example ¶
    
    
    temps := prometheus.NewSummary(prometheus.SummaryOpts{
    	Name:       "pond_temperature_celsius",
    	Help:       "The temperature of the frog pond.",
    	Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
    })
    
    // Simulate some observations.
    for i := 0; i < 1000; i++ {
    	temps.Observe(30 + math.Floor(120*math.Sin(float64(i)*0.1))/10)
    }
    
    // Just for demonstration, let's check the state of the summary by
    // (ab)using its Write method (which is usually only used by Prometheus
    // internally).
    metric := &dto.Metric{}
    temps.Write(metric)
    
    fmt.Println(toNormalizedJSON(sanitizeMetric(metric)))
    
    
    
    Output:
    
    {"summary":{"sampleCount":"1000","sampleSum":29969.50000000001,"quantile":[{"quantile":0.5,"value":31.1},{"quantile":0.9,"value":41.3},{"quantile":0.99,"value":41.9}],"createdTimestamp":"1970-01-01T00:00:10Z"}}
    

####  func [NewSummary](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/summary.go#L182) ¶
    
    
    func NewSummary(opts SummaryOpts) Summary

NewSummary creates a new Summary based on the provided SummaryOpts. 

####  type [SummaryOpts](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/summary.go#L88) ¶
    
    
    type SummaryOpts struct {
    	// Namespace, Subsystem, and Name are components of the fully-qualified
    	// name of the Summary (created by joining these components with
    	// "_"). Only Name is mandatory, the others merely help structuring the
    	// name. Note that the fully-qualified name of the Summary must be a
    	// valid Prometheus metric name.
    	Namespace [string](/builtin#string)
    	Subsystem [string](/builtin#string)
    	Name      [string](/builtin#string)
    
    	// Help provides information about this Summary.
    	//
    	// Metrics with the same fully-qualified name must have the same Help
    	// string.
    	Help [string](/builtin#string)
    
    	// ConstLabels are used to attach fixed labels to this metric. Metrics
    	// with the same fully-qualified name must have the same label names in
    	// their ConstLabels.
    	//
    	// Due to the way a Summary is represented in the Prometheus text format
    	// and how it is handled by the Prometheus server internally, “quantile”
    	// is an illegal label name. Construction of a Summary or SummaryVec
    	// will panic if this label name is used in ConstLabels.
    	//
    	// ConstLabels are only used rarely. In particular, do not use them to
    	// attach the same labels to all your metrics. Those use cases are
    	// better covered by target labels set by the scraping Prometheus
    	// server, or by one specific metric (e.g. a build_info or a
    	// machine_role metric). See also
    	// <https://prometheus.io/docs/instrumenting/writing_exporters/#target-labels-not-static-scraped-labels>
    	ConstLabels Labels
    
    	// Objectives defines the quantile rank estimates with their respective
    	// absolute error. If Objectives[q] = e, then the value reported for q
    	// will be the φ-quantile value for some φ between q-e and q+e.  The
    	// default value is an empty map, resulting in a summary without
    	// quantiles.
    	Objectives map[[float64](/builtin#float64)][float64](/builtin#float64)
    
    	// MaxAge defines the duration for which an observation stays relevant
    	// for the summary. Only applies to pre-calculated quantiles, does not
    	// apply to _sum and _count. Must be positive. The default value is
    	// DefMaxAge.
    	MaxAge [time](/time).[Duration](/time#Duration)
    
    	// AgeBuckets is the number of buckets used to exclude observations that
    	// are older than MaxAge from the summary. A higher number has a
    	// resource penalty, so only increase it if the higher resolution is
    	// really required. For very high observation rates, you might want to
    	// reduce the number of age buckets. With only one age bucket, you will
    	// effectively see a complete reset of the summary each time MaxAge has
    	// passed. The default value is DefAgeBuckets.
    	AgeBuckets [uint32](/builtin#uint32)
    
    	// BufCap defines the default sample stream buffer size.  The default
    	// value of DefBufCap should suffice for most uses. If there is a need
    	// to increase the value, a multiple of 500 is recommended (because that
    	// is the internal buffer size of the underlying package
    	// "github.com/bmizerany/perks/quantile").
    	BufCap [uint32](/builtin#uint32)
    	// contains filtered or unexported fields
    }

SummaryOpts bundles the options for creating a Summary metric. It is mandatory to set Name to a non-empty string. While all other fields are optional and can safely be left at their zero value, it is recommended to set a help string and to explicitly set the Objectives field to the desired value as the default value will change in the upcoming v1.0.0 of the library. 

####  type [SummaryVec](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/summary.go#L552) ¶
    
    
    type SummaryVec struct {
    	*MetricVec
    }

SummaryVec is a Collector that bundles a set of Summaries that all share the same Desc, but have different values for their variable labels. This is used if you want to count the same thing partitioned by various dimensions (e.g. HTTP request latencies, partitioned by status code and method). Create instances with NewSummaryVec. 

Example ¶
    
    
    temps := prometheus.NewSummaryVec(
    	prometheus.SummaryOpts{
    		Name:       "pond_temperature_celsius",
    		Help:       "The temperature of the frog pond.",
    		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
    	},
    	[]string{"species"},
    )
    
    // Simulate some observations.
    for i := 0; i < 1000; i++ {
    	temps.WithLabelValues("litoria-caerulea").Observe(30 + math.Floor(120*math.Sin(float64(i)*0.1))/10)
    	temps.WithLabelValues("lithobates-catesbeianus").Observe(32 + math.Floor(100*math.Cos(float64(i)*0.11))/10)
    }
    
    // Create a Summary without any observations.
    temps.WithLabelValues("leiopelma-hochstetteri")
    
    // Just for demonstration, let's check the state of the summary vector
    // by registering it with a custom registry and then let it collect the
    // metrics.
    reg := prometheus.NewRegistry()
    reg.MustRegister(temps)
    
    metricFamilies, err := reg.Gather()
    if err != nil || len(metricFamilies) != 1 {
    	panic("unexpected behavior of custom test registry")
    }
    
    fmt.Println(toNormalizedJSON(sanitizeMetricFamily(metricFamilies[0])))
    
    
    
    Output:
    
    {"name":"pond_temperature_celsius","help":"The temperature of the frog pond.","type":"SUMMARY","metric":[{"label":[{"name":"species","value":"leiopelma-hochstetteri"}],"summary":{"sampleCount":"0","sampleSum":0,"quantile":[{"quantile":0.5,"value":"NaN"},{"quantile":0.9,"value":"NaN"},{"quantile":0.99,"value":"NaN"}],"createdTimestamp":"1970-01-01T00:00:10Z"}},{"label":[{"name":"species","value":"lithobates-catesbeianus"}],"summary":{"sampleCount":"1000","sampleSum":31956.100000000017,"quantile":[{"quantile":0.5,"value":32.4},{"quantile":0.9,"value":41.4},{"quantile":0.99,"value":41.9}],"createdTimestamp":"1970-01-01T00:00:10Z"}},{"label":[{"name":"species","value":"litoria-caerulea"}],"summary":{"sampleCount":"1000","sampleSum":29969.50000000001,"quantile":[{"quantile":0.5,"value":31.1},{"quantile":0.9,"value":41.3},{"quantile":0.99,"value":41.9}],"createdTimestamp":"1970-01-01T00:00:10Z"}}]}
    

####  func [NewSummaryVec](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/summary.go#L562) ¶
    
    
    func NewSummaryVec(opts SummaryOpts, labelNames [][string](/builtin#string)) *SummaryVec

NewSummaryVec creates a new SummaryVec based on the provided SummaryOpts and partitioned by the given label names. 

Due to the way a Summary is represented in the Prometheus text format and how it is handled by the Prometheus server internally, “quantile” is an illegal label name. NewSummaryVec will panic if this label name is used. 

####  func (*SummaryVec) [CurryWith](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/summary.go#L679) ¶ added in v0.9.0
    
    
    func (v *SummaryVec) CurryWith(labels Labels) (ObserverVec, [error](/builtin#error))

CurryWith returns a vector curried with the provided labels, i.e. the returned vector has those labels pre-set for all labeled operations performed on it. The cardinality of the curried vector is reduced accordingly. The order of the remaining labels stays the same (just with the curried labels taken out of the sequence – which is relevant for the (GetMetric)WithLabelValues methods). It is possible to curry a curried vector, but only with labels not yet used for currying before. 

The metrics contained in the SummaryVec are shared between the curried and uncurried vectors. They are just accessed differently. Curried and uncurried vectors behave identically in terms of collection. Only one must be registered with a given registry (usually the uncurried version). The Reset method deletes all metrics, even if called on a curried vector. 

####  func (*SummaryVec) [GetMetricWith](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/summary.go#L633) ¶
    
    
    func (v *SummaryVec) GetMetricWith(labels Labels) (Observer, [error](/builtin#error))

GetMetricWith returns the Summary for the given Labels map (the label names must match those of the variable labels in Desc). If that label map is accessed for the first time, a new Summary is created. Implications of creating a Summary without using it and keeping the Summary for later use are the same as for GetMetricWithLabelValues. 

An error is returned if the number and names of the Labels are inconsistent with those of the variable labels in Desc (minus any curried labels). 

This method is used for the same purpose as GetMetricWithLabelValues(...string). See there for pros and cons of the two methods. 

####  func (*SummaryVec) [GetMetricWithLabelValues](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/summary.go#L613) ¶
    
    
    func (v *SummaryVec) GetMetricWithLabelValues(lvs ...[string](/builtin#string)) (Observer, [error](/builtin#error))

GetMetricWithLabelValues returns the Summary for the given slice of label values (same order as the variable labels in Desc). If that combination of label values is accessed for the first time, a new Summary is created. 

It is possible to call this method without using the returned Summary to only create the new Summary but leave it at its starting value, a Summary without any observations. 

Keeping the Summary for later use is possible (and should be considered if performance is critical), but keep in mind that Reset, DeleteLabelValues and Delete can be used to delete the Summary from the SummaryVec. In that case, the Summary will still exist, but it will not be exported anymore, even if a Summary with the same label values is created later. See also the CounterVec example. 

An error is returned if the number of label values is not the same as the number of variable labels in Desc (minus any curried labels). 

Note that for more than one label value, this method is prone to mistakes caused by an incorrect order of arguments. Consider GetMetricWith(Labels) as an alternative to avoid that type of mistake. For higher label numbers, the latter has a much more readable (albeit more verbose) syntax, but it comes with a performance overhead (for creating and processing the Labels map). See also the GaugeVec example. 

####  func (*SummaryVec) [MustCurryWith](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/summary.go#L689) ¶ added in v0.9.0
    
    
    func (v *SummaryVec) MustCurryWith(labels Labels) ObserverVec

MustCurryWith works as CurryWith but panics where CurryWith would have returned an error. 

####  func (*SummaryVec) [With](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/summary.go#L658) ¶
    
    
    func (v *SummaryVec) With(labels Labels) Observer

With works as GetMetricWith, but panics where GetMetricWithLabels would have returned an error. Not returning an error allows shortcuts like 
    
    
    myVec.With(prometheus.Labels{"code": "404", "method": "GET"}).Observe(42.21)
    

####  func (*SummaryVec) [WithLabelValues](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/summary.go#L646) ¶
    
    
    func (v *SummaryVec) WithLabelValues(lvs ...[string](/builtin#string)) Observer

WithLabelValues works as GetMetricWithLabelValues, but panics where GetMetricWithLabelValues would have returned an error. Not returning an error allows shortcuts like 
    
    
    myVec.WithLabelValues("404", "GET").Observe(42.21)
    

####  type [SummaryVecOpts](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/summary.go#L157) ¶ added in v1.15.0
    
    
    type SummaryVecOpts struct {
    	SummaryOpts
    
    	// VariableLabels are used to partition the metric vector by the given set
    	// of labels. Each label value will be constrained with the optional Constraint
    	// function, if provided.
    	VariableLabels ConstrainableLabels
    }

SummaryVecOpts bundles the options to create a SummaryVec metric. It is mandatory to set SummaryOpts, see there for mandatory fields. VariableLabels is optional and can safely be left to its default value. 

####  type [Timer](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/timer.go#L20) ¶ added in v0.9.0
    
    
    type Timer struct {
    	// contains filtered or unexported fields
    }

Timer is a helper type to time functions. Use NewTimer to create new instances. 

Example ¶
    
    
    package main
    
    import (
    	"math/rand"
    	"time"
    
    	"github.com/prometheus/client_golang/prometheus"
    )
    
    var requestDuration = prometheus.NewHistogram(prometheus.HistogramOpts{
    	Name:    "example_request_duration_seconds",
    	Help:    "Histogram for the runtime of a simple example function.",
    	Buckets: prometheus.LinearBuckets(0.01, 0.01, 10),
    })
    
    func main() {
    	// timer times this example function. It uses a Histogram, but a Summary
    	// would also work, as both implement Observer. Check out
    	// https://prometheus.io/docs/practices/histograms/ for differences.
    	timer := prometheus.NewTimer(requestDuration)
    	defer timer.ObserveDuration()
    
    	// Do something here that takes time.
    	time.Sleep(time.Duration(rand.NormFloat64()*10000+50000) * time.Microsecond)
    }
    

Share Format Run

Example (Complex) ¶
    
    
    package main
    
    import (
    	"net/http"
    
    	"github.com/prometheus/client_golang/prometheus"
    )
    
    // apiRequestDuration tracks the duration separate for each HTTP status
    // class (1xx, 2xx, ...). This creates a fair amount of time series on
    // the Prometheus server. Usually, you would track the duration of
    // serving HTTP request without partitioning by outcome. Do something
    // like this only if needed. Also note how only status classes are
    // tracked, not every single status code. The latter would create an
    // even larger amount of time series. Request counters partitioned by
    // status code are usually OK as each counter only creates one time
    // series. Histograms are way more expensive, so partition with care and
    // only where you really need separate latency tracking. Partitioning by
    // status class is only an example. In concrete cases, other partitions
    // might make more sense.
    var apiRequestDuration = prometheus.NewHistogramVec(
    	prometheus.HistogramOpts{
    		Name:    "api_request_duration_seconds",
    		Help:    "Histogram for the request duration of the public API, partitioned by status class.",
    		Buckets: prometheus.ExponentialBuckets(0.1, 1.5, 5),
    	},
    	[]string{"status_class"},
    )
    
    func handler(w http.ResponseWriter, r *http.Request) {
    	status := http.StatusOK
    	// The ObserverFunc gets called by the deferred ObserveDuration and
    	// decides which Histogram's Observe method is called.
    	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
    		switch {
    		case status >= 500: // Server error.
    			apiRequestDuration.WithLabelValues("5xx").Observe(v)
    		case status >= 400: // Client error.
    			apiRequestDuration.WithLabelValues("4xx").Observe(v)
    		case status >= 300: // Redirection.
    			apiRequestDuration.WithLabelValues("3xx").Observe(v)
    		case status >= 200: // Success.
    			apiRequestDuration.WithLabelValues("2xx").Observe(v)
    		default: // Informational.
    			apiRequestDuration.WithLabelValues("1xx").Observe(v)
    		}
    	}))
    	defer timer.ObserveDuration()
    
    	// Handle the request. Set status accordingly.
    	// ...
    }
    
    func main() {
    	http.HandleFunc("/api", handler)
    }
    

Share Format Run

Example (Gauge) ¶
    
    
    package main
    
    import (
    	"os"
    
    	"github.com/prometheus/client_golang/prometheus"
    )
    
    // If a function is called rarely (i.e. not more often than scrapes
    // happen) or ideally only once (like in a batch job), it can make sense
    // to use a Gauge for timing the function call. For timing a batch job
    // and pushing the result to a Pushgateway, see also the comprehensive
    // example in the push package.
    var funcDuration = prometheus.NewGauge(prometheus.GaugeOpts{
    	Name: "example_function_duration_seconds",
    	Help: "Duration of the last call of an example function.",
    })
    
    func run() error {
    	// The Set method of the Gauge is used to observe the duration.
    	timer := prometheus.NewTimer(prometheus.ObserverFunc(funcDuration.Set))
    	defer timer.ObserveDuration()
    
    	// Do something. Return errors as encountered. The use of 'defer' above
    	// makes sure the function is still timed properly.
    	return nil
    }
    
    func main() {
    	if err := run(); err != nil {
    		os.Exit(1)
    	}
    }
    

Share Format Run

####  func [NewTimer](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/timer.go#L44) ¶ added in v0.9.0
    
    
    func NewTimer(o Observer) *Timer

NewTimer creates a new Timer. The provided Observer is used to observe a duration in seconds. If the Observer implements ExemplarObserver, passing exemplar later on will be also supported. Timer is usually used to time a function call in the following way: 
    
    
    func TimeMe() {
        timer := NewTimer(myHistogram)
        defer timer.ObserveDuration()
        // Do actual work.
    }
    

or 
    
    
    func TimeMeWithExemplar() {
    	    timer := NewTimer(myHistogram)
    	    defer timer.ObserveDurationWithExemplar(exemplar)
    	    // Do actual work.
    	}
    

####  func (*Timer) [ObserveDuration](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/timer.go#L59) ¶ added in v0.9.0
    
    
    func (t *Timer) ObserveDuration() [time](/time).[Duration](/time#Duration)

ObserveDuration records the duration passed since the Timer was created with NewTimer. It calls the Observe method of the Observer provided during construction with the duration in seconds as an argument. The observed duration is also returned. ObserveDuration is usually called with a defer statement. 

Note that this method is only guaranteed to never observe negative durations if used with Go1.9+. 

####  func (*Timer) [ObserveDurationWithExemplar](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/timer.go#L70) ¶ added in v1.15.0
    
    
    func (t *Timer) ObserveDurationWithExemplar(exemplar Labels) [time](/time).[Duration](/time#Duration)

ObserveDurationWithExemplar is like ObserveDuration, but it will also observe exemplar with the duration unless exemplar is nil or provided Observer can't be casted to ExemplarObserver. 

####  type [TransactionalGatherer](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/registry.go#L1038) ¶ added in v1.13.0
    
    
    type TransactionalGatherer interface {
    	// Gather returns metrics in a lexicographically sorted slice
    	// of uniquely named MetricFamily protobufs. Gather ensures that the
    	// returned slice is valid and self-consistent so that it can be used
    	// for valid exposition. As an exception to the strict consistency
    	// requirements described for metric.Desc, Gather will tolerate
    	// different sets of label names for metrics of the same metric family.
    	//
    	// Even if an error occurs, Gather attempts to gather as many metrics as
    	// possible. Hence, if a non-nil error is returned, the returned
    	// MetricFamily slice could be nil (in case of a fatal error that
    	// prevented any meaningful metric collection) or contain a number of
    	// MetricFamily protobufs, some of which might be incomplete, and some
    	// might be missing altogether. The returned error (which might be a
    	// MultiError) explains the details. Note that this is mostly useful for
    	// debugging purposes. If the gathered protobufs are to be used for
    	// exposition in actual monitoring, it is almost always better to not
    	// expose an incomplete result and instead disregard the returned
    	// MetricFamily protobufs in case the returned error is non-nil.
    	//
    	// Important: done is expected to be triggered (even if the error occurs!)
    	// once caller does not need returned slice of dto.MetricFamily.
    	Gather() (_ []*[dto](/github.com/prometheus/client_model/go).[MetricFamily](/github.com/prometheus/client_model/go#MetricFamily), done func(), err [error](/builtin#error))
    }

TransactionalGatherer represents transactional gatherer that can be triggered to notify gatherer that memory used by metric family is no longer used by a caller. This allows implementations with cache. 

####  func [ToTransactionalGatherer](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/registry.go#L1064) ¶ added in v1.13.0
    
    
    func ToTransactionalGatherer(g Gatherer) TransactionalGatherer

ToTransactionalGatherer transforms Gatherer to transactional one with noop as done function. 

####  type [UnconstrainedLabels](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/labels.go#L101) ¶ added in v1.15.0
    
    
    type UnconstrainedLabels [][string](/builtin#string)

UnconstrainedLabels represents collection of label without any constraint on their value. Thus, it is simply a collection of label names. 
    
    
    UnconstrainedLabels([]string{ "A", "B" })
    

is equivalent to 
    
    
    ConstrainedLabels {
      { Name: "A" },
      { Name: "B" },
    }
    

####  type [UntypedFunc](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/untyped.go#L24) ¶
    
    
    type UntypedFunc interface {
    	Metric
    	Collector
    }

UntypedFunc works like GaugeFunc but the collected metric is of type "Untyped". UntypedFunc is useful to mirror an external metric of unknown type. 

To create UntypedFunc instances, use NewUntypedFunc. 

####  func [NewUntypedFunc](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/untyped.go#L35) ¶
    
    
    func NewUntypedFunc(opts UntypedOpts, function func() [float64](/builtin#float64)) UntypedFunc

NewUntypedFunc creates a new UntypedFunc based on the provided UntypedOpts. The value reported is determined by calling the given function from within the Write method. Take into account that metric collection may happen concurrently. If that results in concurrent calls to Write, like in the case where an UntypedFunc is directly registered with Prometheus, the provided function must be concurrency-safe. 

####  type [UntypedOpts](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/untyped.go#L17) ¶
    
    
    type UntypedOpts Opts

UntypedOpts is an alias for Opts. See there for doc comments. 

####  type [ValueType](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/value.go#L31) ¶
    
    
    type ValueType [int](/builtin#int)

ValueType is an enumeration of metric types that represent a simple value. 
    
    
    const (
    	CounterValue ValueType
    	GaugeValue
    	UntypedValue
    )

Possible values for the ValueType enum. Use UntypedValue to mark a metric with an unknown type. 

####  func (ValueType) [ToDTO](https://github.com/prometheus/client_golang/blob/v1.23.2/prometheus/value.go#L48) ¶ added in v1.13.0
    
    
    func (v ValueType) ToDTO() *[dto](/github.com/prometheus/client_model/go).[MetricType](/github.com/prometheus/client_model/go#MetricType)
  *[v]: View this template
  *[t]: Discuss this template
  *[e]: Edit this template
