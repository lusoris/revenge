# resty HTTP Client

> Auto-fetched from [https://pkg.go.dev/resty.dev/v3](https://pkg.go.dev/resty.dev/v3)
> Last Updated: 2026-01-28T21:43:56.768443+00:00

---

Overview
¶
Package resty provides Simple HTTP, REST, and SSE client library for Go.
Index
¶
Constants
Variables
func AutoParseResponseMiddleware(c *Client, res *Response) (err error)
func CircuitBreaker5xxPolicy(resp *http.Response) bool
func DebugLogFormatter(dl *DebugLog) string
func DebugLogJSONFormatter(dl *DebugLog) string
func PrepareRequestMiddleware(c *Client, r *Request) (err error)
func SaveToFileResponseMiddleware(c *Client, res *Response) error
type CertWatcherOptions
type CircuitBreaker
func NewCircuitBreakerWithCount(failureThreshold uint64, successThreshold uint64, resetTimeout time.Duration, ...) *CircuitBreaker
func NewCircuitBreakerWithRatio(failureRatio float64, minRequests uint64, resetTimeout time.Duration, ...) *CircuitBreaker
func (cb *CircuitBreaker) OnStateChange(hooks ...CircuitBreakerStateChangeHook) *CircuitBreaker
func (cb *CircuitBreaker) OnTrigger(hooks ...CircuitBreakerTriggerHook) *CircuitBreaker
type CircuitBreakerPolicy
type CircuitBreakerState
type CircuitBreakerStateChangeHook
type CircuitBreakerTriggerHook
type Client
func New() *Client
func NewWithClient(hc *http.Client) *Client
func NewWithDialer(dialer *net.Dialer) *Client
func NewWithDialerAndTransportSettings(dialer *net.Dialer, transportSettings *TransportSettings) *Client
func NewWithLocalAddr(localAddr net.Addr) *Client
func NewWithTransportSettings(transportSettings *TransportSettings) *Client
func (c *Client) AddContentDecompresser(k string, d ContentDecompresser) *Client
func (c *Client) AddContentTypeDecoder(ct string, d ContentTypeDecoder) *Client
func (c *Client) AddContentTypeEncoder(ct string, e ContentTypeEncoder) *Client
func (c *Client) AddRequestMiddleware(m RequestMiddleware) *Client
func (c *Client) AddResponseMiddleware(m ResponseMiddleware) *Client
func (c *Client) AddRetryConditions(conditions ...RetryConditionFunc) *Client
func (c *Client) AddRetryHooks(hooks ...RetryHookFunc) *Client
func (c *Client) AllowMethodDeletePayload() bool
func (c *Client) AllowMethodGetPayload() bool
func (c *Client) AllowNonIdempotentRetry() bool
func (c *Client) AuthScheme() string
func (c *Client) AuthToken() string
func (c *Client) BaseURL() string
func (c *Client) Client() *http.Client
func (c *Client) Clone(ctx context.Context) *Client
func (c *Client) Close() error
func (c *Client) ContentDecompresserKeys() string
func (c *Client) ContentDecompressers() map[string]ContentDecompresser
func (c *Client) ContentTypeDecoders() map[string]ContentTypeDecoder
func (c *Client) ContentTypeEncoders() map[string]ContentTypeEncoder
func (c *Client) Context() context.Context
func (c *Client) CookieJar() http.CookieJar
func (c *Client) Cookies() []*http.Cookie
func (c *Client) DebugBodyLimit() int
func (c *Client) DisableDebug() *Client
func (c *Client) DisableGenerateCurlCmd() *Client
func (c *Client) DisableRetryDefaultConditions() *Client
func (c *Client) DisableTrace() *Client
func (c *Client) EnableDebug() *Client
func (c *Client) EnableGenerateCurlCmd() *Client
func (c *Client) EnableRetryDefaultConditions() *Client
func (c *Client) EnableTrace() *Client
func (c *Client) Error() reflect.Type
func (c *Client) FormData() url.Values
func (c *Client) HTTPTransport() (*http.Transport, error)
func (c *Client) Header() http.Header
func (c *Client) HeaderAuthorizationKey() string
func (c *Client) IsDebug() bool
func (c *Client) IsDisableWarn() bool
func (c *Client) IsProxySet() bool
func (c *Client) IsRetryDefaultConditions() bool
func (c *Client) IsSaveResponse() bool
func (c *Client) IsTrace() bool
func (c *Client) LoadBalancer() LoadBalancer
func (c *Client) Logger() Logger
func (c *Client) NewRequest() *Request
func (c *Client) OnClose(h CloseHook) *Client
func (c *Client) OnDebugLog(dlc DebugLogCallbackFunc) *Client
func (c *Client) OnError(h ErrorHook) *Client
func (c *Client) OnInvalid(h ErrorHook) *Client
func (c *Client) OnPanic(h ErrorHook) *Client
func (c *Client) OnSuccess(h SuccessHook) *Client
func (c *Client) OutputDirectory() string
func (c *Client) PathParams() map[string]string
func (c *Client) ProxyURL() *url.URL
func (c *Client) QueryParams() url.Values
func (c *Client) R() *Request
func (c *Client) RemoveProxy() *Client
func (c *Client) ResponseBodyLimit() int64
func (c *Client) ResponseBodyUnlimitedReads() bool
func (c *Client) RetryConditions() []RetryConditionFunc
func (c *Client) RetryCount() int
func (c *Client) RetryHooks() []RetryHookFunc
func (c *Client) RetryMaxWaitTime() time.Duration
func (c *Client) RetryStrategy() RetryStrategyFunc
func (c *Client) RetryWaitTime() time.Duration
func (c *Client) Scheme() string
func (c *Client) SetAllowMethodDeletePayload(allow bool) *Client
func (c *Client) SetAllowMethodGetPayload(allow bool) *Client
func (c *Client) SetAllowNonIdempotentRetry(b bool) *Client
func (c *Client) SetAuthScheme(scheme string) *Client
func (c *Client) SetAuthToken(token string) *Client
func (c *Client) SetBaseURL(url string) *Client
func (c *Client) SetBasicAuth(username, password string) *Client
func (c *Client) SetCertificateFromFile(certFilePath, certKeyFilePath string) *Client
func (c *Client) SetCertificateFromString(certStr, certKeyStr string) *Client
func (c *Client) SetCertificates(certs ...tls.Certificate) *Client
func (c *Client) SetCircuitBreaker(b *CircuitBreaker) *Client
func (c *Client) SetClientRootCertificateFromString(pemCerts string) *Client
func (c *Client) SetClientRootCertificates(pemFilePaths ...string) *Client
func (c *Client) SetClientRootCertificatesWatcher(options *CertWatcherOptions, pemFilePaths ...string) *Client
func (c *Client) SetCloseConnection(close bool) *Client
func (c *Client) SetContentDecompresserKeys(keys []string) *Client
func (c *Client) SetContext(ctx context.Context) *Client
func (c *Client) SetCookie(hc *http.Cookie) *Client
func (c *Client) SetCookieJar(jar http.CookieJar) *Client
func (c *Client) SetCookies(cs []*http.Cookie) *Client
func (c *Client) SetDebug(d bool) *Client
func (c *Client) SetDebugBodyLimit(sl int) *Client
func (c *Client) SetDebugLogCurlCmd(b bool) *Client
func (c *Client) SetDebugLogFormatter(df DebugLogFormatterFunc) *Client
func (c *Client) SetDigestAuth(username, password string) *Client
func (c *Client) SetDisableWarn(d bool) *Client
func (c *Client) SetDoNotParseResponse(notParse bool) *Client
func (c *Client) SetError(v any) *Client
func (c *Client) SetFormData(data map[string]string) *Client
func (c *Client) SetGenerateCurlCmd(b bool) *Client
func (c *Client) SetHeader(header, value string) *Client
func (c *Client) SetHeaderAuthorizationKey(k string) *Client
func (c *Client) SetHeaderVerbatim(header, value string) *Client
func (c *Client) SetHeaders(headers map[string]string) *Client
func (c *Client) SetJSONEscapeHTML(b bool) *Client
func (c *Client) SetLoadBalancer(b LoadBalancer) *Client
func (c *Client) SetLogger(l Logger) *Client
func (c *Client) SetOutputDirectory(dirPath string) *Client
func (c *Client) SetPathParam(param, value string) *Client
func (c *Client) SetPathParams(params map[string]string) *Client
func (c *Client) SetProxy(proxyURL string) *Client
func (c *Client) SetQueryParam(param, value string) *Client
func (c *Client) SetQueryParams(params map[string]string) *Client
func (c *Client) SetRawPathParam(param, value string) *Client
func (c *Client) SetRawPathParams(params map[string]string) *Client
func (c *Client) SetRedirectPolicy(policies ...RedirectPolicy) *Client
func (c *Client) SetRequestMiddlewares(middlewares ...RequestMiddleware) *Client
func (c *Client) SetResponseBodyLimit(v int64) *Client
func (c *Client) SetResponseBodyUnlimitedReads(b bool) *Client
func (c *Client) SetResponseMiddlewares(middlewares ...ResponseMiddleware) *Client
func (c *Client) SetRetryCount(count int) *Client
func (c *Client) SetRetryDefaultConditions(b bool) *Client
func (c *Client) SetRetryMaxWaitTime(maxWaitTime time.Duration) *Client
func (c *Client) SetRetryStrategy(rs RetryStrategyFunc) *Client
func (c *Client) SetRetryWaitTime(waitTime time.Duration) *Client
func (c *Client) SetRootCertificateFromString(pemCerts string) *Client
func (c *Client) SetRootCertificates(pemFilePaths ...string) *Client
func (c *Client) SetRootCertificatesWatcher(options *CertWatcherOptions, pemFilePaths ...string) *Client
func (c *Client) SetSaveResponse(save bool) *Client
func (c *Client) SetScheme(scheme string) *Client
func (c *Client) SetTLSClientConfig(tlsConfig *tls.Config) *Client
func (c *Client) SetTimeout(timeout time.Duration) *Client
func (c *Client) SetTrace(t bool) *Client
func (c *Client) SetTransport(transport http.RoundTripper) *Client
func (c *Client) SetUnescapeQueryParams(unescape bool) *Client
func (c *Client) TLSClientConfig() *tls.Config
func (c *Client) Timeout() time.Duration
func (c *Client) Transport() http.RoundTripper
type CloseHook
type ContentDecompresser
type ContentTypeDecoder
type ContentTypeEncoder
type DebugLog
type DebugLogCallbackFunc
type DebugLogFormatterFunc
type DebugLogRequest
type DebugLogResponse
type ErrorHook
type Event
type EventErrorFunc
type EventMessageFunc
type EventOpenFunc
type EventRequestFailureFunc
type EventSource
func NewEventSource() *EventSource
func (es *EventSource) AddEventListener(eventName string, ef EventMessageFunc, result any) *EventSource
func (es *EventSource) AddHeader(header, value string) *EventSource
func (es *EventSource) Close()
func (es *EventSource) Get() error
func (es *EventSource) Logger() Logger
func (es *EventSource) OnError(ef EventErrorFunc) *EventSource
func (es *EventSource) OnMessage(ef EventMessageFunc, result any) *EventSource
func (es *EventSource) OnOpen(ef EventOpenFunc) *EventSource
func (es *EventSource) OnRequestFailure(ef EventRequestFailureFunc) *EventSource
func (es *EventSource) SetBody(body io.Reader) *EventSource
func (es *EventSource) SetHeader(header, value string) *EventSource
func (es *EventSource) SetLogger(l Logger) *EventSource
func (es *EventSource) SetMaxBufSize(bufSize int) *EventSource
func (es *EventSource) SetMethod(method string) *EventSource
func (es *EventSource) SetRetryCount(count int) *EventSource
func (es *EventSource) SetRetryMaxWaitTime(maxWaitTime time.Duration) *EventSource
func (es *EventSource) SetRetryWaitTime(waitTime time.Duration) *EventSource
func (es *EventSource) SetTLSClientConfig(tlsConfig *tls.Config) *EventSource
func (es *EventSource) SetURL(url string) *EventSource
func (es *EventSource) TLSClientConfig() *tls.Config
type Host
type HostState
type HostStateChangeFunc
type LoadBalancer
type Logger
type MultipartField
func (mf *MultipartField) Clone() *MultipartField
type MultipartFieldCallbackFunc
type MultipartFieldProgress
func (mfp MultipartFieldProgress) String() string
type RedirectInfo
type RedirectPolicy
func DomainCheckRedirectPolicy(hostnames ...string) RedirectPolicy
func FlexibleRedirectPolicy(noOfRedirect int) RedirectPolicy
func NoRedirectPolicy() RedirectPolicy
type RedirectPolicyFunc
func (f RedirectPolicyFunc) Apply(req *http.Request, via []*http.Request) error
type Request
func (r *Request) AddRetryConditions(conditions ...RetryConditionFunc) *Request
func (r *Request) AddRetryHooks(hooks ...RetryHookFunc) *Request
func (r *Request) Clone(ctx context.Context) *Request
func (r *Request) Context() context.Context
func (r *Request) CurlCmd() string
func (r *Request) Delete(url string) (*Response, error)
func (r *Request) DisableDebug() *Request
func (r *Request) DisableGenerateCurlCmd() *Request
func (r *Request) DisableRetryDefaultConditions() *Request
func (r *Request) DisableTrace() *Request
func (r *Request) EnableDebug() *Request
func (r *Request) EnableGenerateCurlCmd() *Request
func (r *Request) EnableRetryDefaultConditions() *Request
func (r *Request) EnableTrace() *Request
func (r *Request) Execute(method, url string) (res *Response, err error)
func (r *Request) Funcs(funcs ...RequestFunc) *Request
func (r *Request) Get(url string) (*Response, error)
func (r *Request) Head(url string) (*Response, error)
func (r *Request) Options(url string) (*Response, error)
func (r *Request) Patch(url string) (*Response, error)
func (r *Request) Post(url string) (*Response, error)
func (r *Request) Put(url string) (*Response, error)
func (r *Request) Send() (*Response, error)
func (r *Request) SetAllowMethodDeletePayload(allow bool) *Request
func (r *Request) SetAllowMethodGetPayload(allow bool) *Request
func (r *Request) SetAllowNonIdempotentRetry(b bool) *Request
func (r *Request) SetAuthScheme(scheme string) *Request
func (r *Request) SetAuthToken(authToken string) *Request
func (r *Request) SetBasicAuth(username, password string) *Request
func (r *Request) SetBody(body any) *Request
func (r *Request) SetCloseConnection(close bool) *Request
func (r *Request) SetContentLength(v int64) *Request
func (r *Request) SetContentType(ct string) *Request
func (r *Request) SetContext(ctx context.Context) *Request
func (r *Request) SetCookie(hc *http.Cookie) *Request
func (r *Request) SetCookies(rs []*http.Cookie) *Request
func (r *Request) SetDebug(d bool) *Request
func (r *Request) SetDebugLogCurlCmd(b bool) *Request
func (r *Request) SetDoNotParseResponse(notParse bool) *Request
func (r *Request) SetError(err any) *Request
func (r *Request) SetExpectResponseContentType(contentType string) *Request
func (r *Request) SetFile(fieldName, filePath string) *Request
func (r *Request) SetFileReader(fieldName, fileName string, reader io.Reader) *Request
func (r *Request) SetFiles(files map[string]string) *Request
func (r *Request) SetForceResponseContentType(contentType string) *Request
func (r *Request) SetFormData(data map[string]string) *Request
func (r *Request) SetFormDataFromValues(data url.Values) *Request
func (r *Request) SetGenerateCurlCmd(b bool) *Request
func (r *Request) SetHeader(header, value string) *Request
func (r *Request) SetHeaderAuthorizationKey(k string) *Request
func (r *Request) SetHeaderMultiValues(headers map[string][]string) *Request
func (r *Request) SetHeaderVerbatim(header, value string) *Request
func (r *Request) SetHeaders(headers map[string]string) *Request
func (r *Request) SetJSONEscapeHTML(b bool) *Request
func (r *Request) SetLogger(l Logger) *Request
func (r *Request) SetMethod(m string) *Request
func (r *Request) SetMultipartBoundary(boundary string) *Request
func (r *Request) SetMultipartField(fieldName, fileName, contentType string, reader io.Reader) *Request
func (r *Request) SetMultipartFields(fields ...*MultipartField) *Request
func (r *Request) SetMultipartFormData(data map[string]string) *Request
func (r *Request) SetMultipartOrderedFormData(name string, values []string) *Request
func (r *Request) SetOutputFileName(file string) *Request
func (r *Request) SetPathParam(param, value string) *Request
func (r *Request) SetPathParams(params map[string]string) *Request
func (r *Request) SetQueryParam(param, value string) *Request
func (r *Request) SetQueryParams(params map[string]string) *Request
func (r *Request) SetQueryParamsFromValues(params url.Values) *Request
func (r *Request) SetQueryString(query string) *Request
func (r *Request) SetRawPathParam(param, value string) *Request
func (r *Request) SetRawPathParams(params map[string]string) *Request
func (r *Request) SetResponseBodyLimit(v int64) *Request
func (r *Request) SetResponseBodyUnlimitedReads(b bool) *Request
func (r *Request) SetResult(v any) *Request
func (r *Request) SetRetryConditions(conditions ...RetryConditionFunc) *Request
func (r *Request) SetRetryCount(count int) *Request
func (r *Request) SetRetryDefaultConditions(b bool) *Request
func (r *Request) SetRetryHooks(hooks ...RetryHookFunc) *Request
func (r *Request) SetRetryMaxWaitTime(maxWaitTime time.Duration) *Request
func (r *Request) SetRetryStrategy(rs RetryStrategyFunc) *Request
func (r *Request) SetRetryWaitTime(waitTime time.Duration) *Request
func (r *Request) SetSaveResponse(save bool) *Request
func (r *Request) SetTimeout(timeout time.Duration) *Request
func (r *Request) SetTrace(t bool) *Request
func (r *Request) SetURL(url string) *Request
func (r *Request) SetUnescapeQueryParams(unescape bool) *Request
func (r *Request) Trace(url string) (*Response, error)
func (r *Request) TraceInfo() TraceInfo
func (r *Request) WithContext(ctx context.Context) *Request
type RequestFeedback
type RequestFunc
type RequestMiddleware
type Response
func (r *Response) Bytes() []byte
func (r *Response) Cookies() []*http.Cookie
func (r *Response) Duration() time.Duration
func (r *Response) Error() any
func (r *Response) Header() http.Header
func (r *Response) IsError() bool
func (r *Response) IsSuccess() bool
func (r *Response) Proto() string
func (r *Response) ReceivedAt() time.Time
func (r *Response) RedirectHistory() []*RedirectInfo
func (r *Response) Result() any
func (r *Response) Size() int64
func (r *Response) Status() string
func (r *Response) StatusCode() int
func (r *Response) String() string
type ResponseError
func (e *ResponseError) Error() string
func (e *ResponseError) Unwrap() error
type ResponseMiddleware
type RetryConditionFunc
type RetryHookFunc
type RetryStrategyFunc
type RoundRobin
func NewRoundRobin(baseURLs ...string) (*RoundRobin, error)
func (rr *RoundRobin) Close() error
func (rr *RoundRobin) Feedback(_ *RequestFeedback)
func (rr *RoundRobin) Next() (string, error)
func (rr *RoundRobin) Refresh(baseURLs ...string) error
type SRVWeightedRoundRobin
func NewSRVWeightedRoundRobin(service, proto, domainName, httpScheme string) (*SRVWeightedRoundRobin, error)
func (swrr *SRVWeightedRoundRobin) Close() error
func (swrr *SRVWeightedRoundRobin) Feedback(f *RequestFeedback)
func (swrr *SRVWeightedRoundRobin) Next() (string, error)
func (swrr *SRVWeightedRoundRobin) Refresh() error
func (swrr *SRVWeightedRoundRobin) SetOnStateChange(fn HostStateChangeFunc)
func (swrr *SRVWeightedRoundRobin) SetRecoveryDuration(d time.Duration)
func (swrr *SRVWeightedRoundRobin) SetRefreshDuration(d time.Duration)
type SuccessHook
type TLSClientConfiger
type TraceInfo
func (ti TraceInfo) Clone() *TraceInfo
func (ti TraceInfo) JSON() string
func (ti TraceInfo) String() string
type TransportSettings
type WeightedRoundRobin
func NewWeightedRoundRobin(recovery time.Duration, hosts ...*Host) (*WeightedRoundRobin, error)
func (wrr *WeightedRoundRobin) Close() error
func (wrr *WeightedRoundRobin) Feedback(f *RequestFeedback)
func (wrr *WeightedRoundRobin) Next() (string, error)
func (wrr *WeightedRoundRobin) Refresh(hosts ...*Host) error
func (wrr *WeightedRoundRobin) SetOnStateChange(fn HostStateChangeFunc)
func (wrr *WeightedRoundRobin) SetRecoveryDuration(d time.Duration)
Constants
¶
View Source
const (
// MethodGet HTTP method
MethodGet = "GET"
// MethodPost HTTP method
MethodPost = "POST"
// MethodPut HTTP method
MethodPut = "PUT"
// MethodDelete HTTP method
MethodDelete = "DELETE"
// MethodPatch HTTP method
MethodPatch = "PATCH"
// MethodHead HTTP method
MethodHead = "HEAD"
// MethodOptions HTTP method
MethodOptions = "OPTIONS"
// MethodTrace HTTP method
MethodTrace = "TRACE"
)
View Source
const Version = "3.0.0-beta.6"
Version # of resty
Variables
¶
View Source
var (
ErrNotHttpTransportType       =
errors
.
New
("resty: not a http.Transport type")
ErrUnsupportedRequestBodyKind =
errors
.
New
("resty: unsupported request body kind")
)
View Source
var (
ErrDigestBadChallenge    =
errors
.
New
("resty: digest: challenge is bad")
ErrDigestInvalidCharset  =
errors
.
New
("resty: digest: invalid charset")
ErrDigestAlgNotSupported =
errors
.
New
("resty: digest: algorithm is not supported")
ErrDigestQopNotSupported =
errors
.
New
("resty: digest: qop is not supported")
)
View Source
var (
// InMemoryJSONMarshal function performs the JSON marshalling completely in memory.
//
//	c := resty.New()
//	defer c.Close()
//
//	c.AddContentTypeEncoder("application/json", resty.InMemoryJSONMarshal)
InMemoryJSONMarshal = func(w
io
.
Writer
, v
any
)
error
{
jsonData, err :=
json
.
Marshal
(v)
if err !=
nil
{
return err
}
_, err = w.Write(jsonData)
return err
}
// InMemoryJSONUnmarshal function performs the JSON unmarshalling completely in memory.
//
//	c := resty.New()
//	defer c.Close()
//
//	c.AddContentTypeDecoder("application/json", resty.InMemoryJSONUnmarshal)
InMemoryJSONUnmarshal = func(r
io
.
Reader
, v
any
)
error
{
byteData, err :=
io
.
ReadAll
(r)
if err !=
nil
{
return err
}
return
json
.
Unmarshal
(byteData, v)
}
// InMemoryXMLMarshal function performs the XML marshalling completely in memory.
//
//	c := resty.New()
//	defer c.Close()
//
//	c.AddContentTypeEncoder("application/xml", resty.InMemoryXMLMarshal)
InMemoryXMLMarshal = func(w
io
.
Writer
, v
any
)
error
{
xmlData, err :=
xml
.
Marshal
(v)
if err !=
nil
{
return err
}
_, err = w.Write(xmlData)
return err
}
// InMemoryJSONUnmarshal function performs the XML unmarshalling completely in memory.
//
//	c := resty.New()
//	defer c.Close()
//
//	c.AddContentTypeDecoder("application/xml", resty.InMemoryXMLUnmarshal)
InMemoryXMLUnmarshal = func(r
io
.
Reader
, v
any
)
error
{
byteData, err :=
io
.
ReadAll
(r)
if err !=
nil
{
return err
}
return
xml
.
Unmarshal
(byteData, v)
}
)
View Source
var ErrCircuitBreakerOpen =
errors
.
New
("resty: circuit breaker open")
ErrCircuitBreakerOpen is returned when the circuit breaker is open.
View Source
var (
ErrContentDecompresserNotFound =
errors
.
New
("resty: content decoder not found")
)
View Source
var ErrNoActiveHost =
errors
.
New
("resty: no active host")
ErrNoActiveHost error returned when all hosts are inactive on the load balancer
View Source
var ErrReadExceedsThresholdLimit =
errors
.
New
("resty: read exceeds the threshold limit")
Functions
¶
func
AutoParseResponseMiddleware
¶
func AutoParseResponseMiddleware(c *
Client
, res *
Response
) (err
error
)
AutoParseResponseMiddleware method is used to parse the response body automatically
based on registered HTTP response `Content-Type` decoder, see
Client.AddContentTypeDecoder
;
if
Request.SetResult
,
Request.SetError
, or
Client.SetError
is used
func
CircuitBreaker5xxPolicy
¶
func CircuitBreaker5xxPolicy(resp *
http
.
Response
)
bool
CircuitBreaker5xxPolicy is a
CircuitBreakerPolicy
that trips the
CircuitBreaker
if
the response status code is 500 or greater.
func
DebugLogFormatter
¶
func DebugLogFormatter(dl *
DebugLog
)
string
DebugLogFormatter function formats the given debug log info in human readable
format.
This is the default debug log formatter in the Resty.
func
DebugLogJSONFormatter
¶
func DebugLogJSONFormatter(dl *
DebugLog
)
string
DebugLogJSONFormatter function formats the given debug log info in JSON format.
func
PrepareRequestMiddleware
¶
func PrepareRequestMiddleware(c *
Client
, r *
Request
) (err
error
)
PrepareRequestMiddleware method is used to prepare HTTP requests from
user provides request values. Request preparation fails if any error occurs
func
SaveToFileResponseMiddleware
¶
func SaveToFileResponseMiddleware(c *
Client
, res *
Response
)
error
SaveToFileResponseMiddleware method used to write HTTP response body into
file. The filename is determined in the following order -
Request.SetOutputFileName
Content-Disposition header
Request URL using
path.Base
Types
¶
type
CertWatcherOptions
¶
type CertWatcherOptions struct {
// PoolInterval is the frequency at which resty will check if the PEM file needs to be reloaded.
// Default is 24 hours.
PoolInterval
time
.
Duration
}
CertWatcherOptions allows configuring a watcher that reloads dynamically TLS certs.
type
CircuitBreaker
¶
type CircuitBreaker struct {
// contains filtered or unexported fields
}
CircuitBreaker struct implements a state machine to monitor and manage the
states of circuit breakers. The three states are:
Closed: requests are allowed
Open: requests are blocked
Half-Open: a single request is allowed to determine
Transitions
To Closed State: when the success count reaches the success threshold.
To Open State: when the failure count reaches the failure threshold.
Half-Open Check: when the specified timeout reaches, a single request is allowed
to determine the transition state; if failed, it goes back to the open state.
Use
NewCircuitBreakerWithCount
or
NewCircuitBreakerWithRatio
to create a new
CircuitBreaker
instance accordingly.
func
NewCircuitBreakerWithCount
¶
func NewCircuitBreakerWithCount(failureThreshold
uint64
, successThreshold
uint64
,
resetTimeout
time
.
Duration
, policies ...
CircuitBreakerPolicy
) *
CircuitBreaker
NewCircuitBreakerWithCount method creates a new
CircuitBreaker
instance with Count settings.
The default settings are:
Policies: CircuitBreaker5xxPolicy
func
NewCircuitBreakerWithRatio
¶
func NewCircuitBreakerWithRatio(failureRatio
float64
, minRequests
uint64
,
resetTimeout
time
.
Duration
, policies ...
CircuitBreakerPolicy
) *
CircuitBreaker
NewCircuitBreakerWithRatio method creates a new
CircuitBreaker
instance with Ratio settings.
The default settings are:
Policies: CircuitBreaker5xxPolicy
func (*CircuitBreaker)
OnStateChange
¶
func (cb *
CircuitBreaker
) OnStateChange(hooks ...
CircuitBreakerStateChangeHook
) *
CircuitBreaker
OnStateChange method adds a
CircuitBreakerStateChangeHook
to the
CircuitBreaker
instance.
func (*CircuitBreaker)
OnTrigger
¶
func (cb *
CircuitBreaker
) OnTrigger(hooks ...
CircuitBreakerTriggerHook
) *
CircuitBreaker
OnTrigger method adds a
CircuitBreakerTriggerHook
to the
CircuitBreaker
instance.
type
CircuitBreakerPolicy
¶
type CircuitBreakerPolicy func(resp *
http
.
Response
)
bool
CircuitBreakerPolicy is a function type that determines whether a response should
trip the
CircuitBreaker
.
type
CircuitBreakerState
¶
type CircuitBreakerState
uint32
CircuitBreakerState type represents the state of the circuit breaker.
const (
// CircuitBreakerStateClosed represents the closed state of the circuit breaker.
CircuitBreakerStateClosed
CircuitBreakerState
=
iota
// CircuitBreakerStateOpen represents the open state of the circuit breaker.
CircuitBreakerStateOpen
// CircuitBreakerStateHalfOpen represents the half-open state of the circuit breaker.
CircuitBreakerStateHalfOpen
)
type
CircuitBreakerStateChangeHook
¶
type CircuitBreakerStateChangeHook func(oldState, newState
CircuitBreakerState
)
CircuitBreakerStateChangeHook type is for reacting to circuit breaker state change hooks.
type
CircuitBreakerTriggerHook
¶
type CircuitBreakerTriggerHook func(*
Request
,
error
)
CircuitBreakerTriggerHook type is for reacting to circuit breaker trigger hooks.
type
Client
¶
type Client struct {
// contains filtered or unexported fields
}
Client struct is used to create a Resty client with client-level settings,
these settings apply to all the requests raised from the client.
Resty also provides an option to override most of the client settings
at
Request
level.
func
New
¶
func New() *
Client
New method creates a new Resty client.
func
NewWithClient
¶
func NewWithClient(hc *
http
.
Client
) *
Client
NewWithClient method creates a new Resty client with given
http.Client
.
func
NewWithDialer
¶
func NewWithDialer(dialer *
net
.
Dialer
) *
Client
NewWithDialer method creates a new Resty client with given Local Address
to dial from.
func
NewWithDialerAndTransportSettings
¶
func NewWithDialerAndTransportSettings(dialer *
net
.
Dialer
, transportSettings *
TransportSettings
) *
Client
NewWithDialerAndTransportSettings method creates a new Resty client with given Local Address
to dial from.
func
NewWithLocalAddr
¶
func NewWithLocalAddr(localAddr
net
.
Addr
) *
Client
NewWithLocalAddr method creates a new Resty client with the given Local Address.
func
NewWithTransportSettings
¶
func NewWithTransportSettings(transportSettings *
TransportSettings
) *
Client
NewWithTransportSettings method creates a new Resty client with provided
timeout values.
func (*Client)
AddContentDecompresser
¶
func (c *
Client
) AddContentDecompresser(k
string
, d
ContentDecompresser
) *
Client
AddContentDecompresser method adds the user-provided Content-Encoding (
RFC 9110
) Decompresser
and directive into a client.
NOTE: It overwrites the Decompresser function if the given Content-Encoding directive already exists.
func (*Client)
AddContentTypeDecoder
¶
func (c *
Client
) AddContentTypeDecoder(ct
string
, d
ContentTypeDecoder
) *
Client
AddContentTypeDecoder method adds the user-provided Content-Type decoder into a client.
NOTE: It overwrites the decoder function if the given Content-Type key already exists.
func (*Client)
AddContentTypeEncoder
¶
func (c *
Client
) AddContentTypeEncoder(ct
string
, e
ContentTypeEncoder
) *
Client
AddContentTypeEncoder method adds the user-provided Content-Type encoder into a client.
NOTE: It overwrites the encoder function if the given Content-Type key already exists.
func (*Client)
AddRequestMiddleware
¶
func (c *
Client
) AddRequestMiddleware(m
RequestMiddleware
) *
Client
AddRequestMiddleware method appends a request middleware to the before request chain.
After all requests, middlewares are applied, and the request is sent to the host server.
client.AddRequestMiddleware(func(c *resty.Client, r *resty.Request) error {
// Now you have access to the Client and Request instance
// manipulate it as per your need

return nil 	// if its successful otherwise return error
})
func (*Client)
AddResponseMiddleware
¶
func (c *
Client
) AddResponseMiddleware(m
ResponseMiddleware
) *
Client
AddResponseMiddleware method appends response middleware to the after-response chain.
All the response middlewares are applied; once we receive a response
from the host server.
client.AddResponseMiddleware(func(c *resty.Client, r *resty.Response) error {
// Now you have access to the Client and Response instance
// Also, you could access request via Response.Request i.e., r.Request
// manipulate it as per your need

return nil 	// if its successful otherwise return error
})
func (*Client)
AddRetryConditions
¶
func (c *
Client
) AddRetryConditions(conditions ...
RetryConditionFunc
) *
Client
AddRetryConditions method adds one or more retry condition functions into the request.
These retry conditions are executed to determine if the request can be retried.
The request will retry if any functions return `true`, otherwise return `false`.
NOTE:
The default retry conditions are applied first.
The client-level retry conditions are applied to all requests.
The request-level retry conditions are executed first before the client-level
retry conditions. See
Request.AddRetryConditions
,
Request.SetRetryConditions
func (*Client)
AddRetryHooks
¶
func (c *
Client
) AddRetryHooks(hooks ...
RetryHookFunc
) *
Client
AddRetryHooks method adds one or more side-effecting retry hooks to an array
of hooks that will be executed on each retry.
NOTE:
All the retry hooks are executed on request retry.
The request-level retry hooks are executed first before client-level hooks.
func (*Client)
AllowMethodDeletePayload
¶
func (c *
Client
) AllowMethodDeletePayload()
bool
AllowMethodDeletePayload method returns `true` if the client is enabled to allow
payload with DELETE method; otherwise, it is `false`.
More info, refer to GH#881
func (*Client)
AllowMethodGetPayload
¶
func (c *
Client
) AllowMethodGetPayload()
bool
AllowMethodGetPayload method returns `true` if the client is enabled to allow
payload with GET method; otherwise, it is `false`.
func (*Client)
AllowNonIdempotentRetry
¶
func (c *
Client
) AllowNonIdempotentRetry()
bool
AllowNonIdempotentRetry method returns true if the client is enabled to allow
non-idempotent HTTP methods retry; otherwise, it is `false`
Default value is `false`
func (*Client)
AuthScheme
¶
func (c *
Client
) AuthScheme()
string
AuthScheme method returns the auth scheme name set in the client instance.
See
Client.SetAuthScheme
,
Request.SetAuthScheme
.
func (*Client)
AuthToken
¶
func (c *
Client
) AuthToken()
string
AuthToken method returns the auth token value registered in the client instance.
func (*Client)
BaseURL
¶
func (c *
Client
) BaseURL()
string
BaseURL method returns the Base URL value from the client instance.
func (*Client)
Client
¶
func (c *
Client
) Client() *
http
.
Client
Client method returns the underlying Go
http.Client
used by the Resty.
func (*Client)
Clone
¶
func (c *
Client
) Clone(ctx
context
.
Context
) *
Client
Clone method returns a clone of the original client.
NOTE: Use with care:
Interface values are not deeply cloned. Thus, both the original and the
clone will use the same value.
It is not safe for concurrent use. You should only use this method
when you are sure that any other concurrent process is not using the client
or client instance is protected by a mutex.
func (*Client)
Close
¶
func (c *
Client
) Close()
error
Close method performs cleanup and closure activities on the client instance
func (*Client)
ContentDecompresserKeys
¶
func (c *
Client
) ContentDecompresserKeys()
string
ContentDecompresserKeys method returns all the registered content-encoding Decompressers
keys as comma-separated string.
func (*Client)
ContentDecompressers
¶
func (c *
Client
) ContentDecompressers() map[
string
]
ContentDecompresser
ContentDecompressers method returns all the registered content-encoding Decompressers.
func (*Client)
ContentTypeDecoders
¶
func (c *
Client
) ContentTypeDecoders() map[
string
]
ContentTypeDecoder
ContentTypeDecoders method returns all the registered content type decoders.
func (*Client)
ContentTypeEncoders
¶
func (c *
Client
) ContentTypeEncoders() map[
string
]
ContentTypeEncoder
ContentTypeEncoders method returns all the registered content type encoders.
func (*Client)
Context
¶
func (c *
Client
) Context()
context
.
Context
Context method returns the
context.Context
from the client instance.
func (*Client)
CookieJar
¶
func (c *
Client
) CookieJar()
http
.
CookieJar
CookieJar method returns the HTTP cookie jar instance from the underlying Go HTTP Client.
func (*Client)
Cookies
¶
func (c *
Client
) Cookies() []*
http
.
Cookie
Cookies method returns all cookies registered in the client instance.
func (*Client)
DebugBodyLimit
¶
func (c *
Client
) DebugBodyLimit()
int
DebugBodyLimit method returns the debug body limit value set on the client instance
func (*Client)
DisableDebug
¶
func (c *
Client
) DisableDebug() *
Client
DisableDebug method is a helper method for
Client.SetDebug
func (*Client)
DisableGenerateCurlCmd
¶
func (c *
Client
) DisableGenerateCurlCmd() *
Client
DisableGenerateCurlCmd method disables the option set by
Client.EnableGenerateCurlCmd
or
Client.SetGenerateCurlCmd
.
func (*Client)
DisableRetryDefaultConditions
¶
func (c *
Client
) DisableRetryDefaultConditions() *
Client
DisableRetryDefaultConditions method disables the Resty's default retry conditions
func (*Client)
DisableTrace
¶
func (c *
Client
) DisableTrace() *
Client
DisableTrace method disables the Resty client trace. Refer to
Client.EnableTrace
.
func (*Client)
EnableDebug
¶
func (c *
Client
) EnableDebug() *
Client
EnableDebug method is a helper method for
Client.SetDebug
func (*Client)
EnableGenerateCurlCmd
¶
func (c *
Client
) EnableGenerateCurlCmd() *
Client
EnableGenerateCurlCmd method enables the generation of curl command at the
client instance level.
By default, Resty does not log the curl command in the debug log since it has the potential
to leak sensitive data unless explicitly enabled via
Client.SetDebugLogCurlCmd
or
Request.SetDebugLogCurlCmd
.
NOTE: Use with care.
Potential to leak sensitive data from
Request
and
Response
in the debug log
when the debug log option is enabled.
Additional memory usage since the request body was reread.
curl body is not generated for
io.Reader
and multipart request flow.
func (*Client)
EnableRetryDefaultConditions
¶
func (c *
Client
) EnableRetryDefaultConditions() *
Client
EnableRetryDefaultConditions method enables the Resty's default retry conditions
func (*Client)
EnableTrace
¶
func (c *
Client
) EnableTrace() *
Client
EnableTrace method enables the Resty client trace for the requests fired from
the client using
httptrace.ClientTrace
and provides insights.
client := resty.New().EnableTrace()

resp, err := client.R().Get("https://httpbin.org/get")
fmt.Println("error:", err)
fmt.Println("Trace Info:", resp.Request.TraceInfo())
The method
Request.EnableTrace
is also available to get trace info for a single request.
func (*Client)
Error
¶
func (c *
Client
) Error()
reflect
.
Type
Error method returns the global or client common `Error` object type registered in the Resty.
func (*Client)
FormData
¶
func (c *
Client
) FormData()
url
.
Values
FormData method returns the form parameters and their values from the client instance.
func (*Client)
HTTPTransport
¶
func (c *
Client
) HTTPTransport() (*
http
.
Transport
,
error
)
HTTPTransport method does type assertion and returns
http.Transport
from the client instance, if type assertion fails it returns an error
func (*Client)
Header
¶
func (c *
Client
) Header()
http
.
Header
Header method returns the headers from the client instance.
func (*Client)
HeaderAuthorizationKey
¶
func (c *
Client
) HeaderAuthorizationKey()
string
HeaderAuthorizationKey method returns the HTTP header name for Authorization from the client instance.
func (*Client)
IsDebug
¶
func (c *
Client
) IsDebug()
bool
IsDebug method returns `true` if the client is in debug mode; otherwise, it is `false`.
func (*Client)
IsDisableWarn
¶
func (c *
Client
) IsDisableWarn()
bool
IsDisableWarn method returns `true` if the warning message is disabled; otherwise, it is `false`.
func (*Client)
IsProxySet
¶
func (c *
Client
) IsProxySet()
bool
IsProxySet method returns the true is proxy is set from the Resty client; otherwise
false. By default, the proxy is set from the environment variable; refer to
http.ProxyFromEnvironment
.
func (*Client)
IsRetryDefaultConditions
¶
func (c *
Client
) IsRetryDefaultConditions()
bool
IsRetryDefaultConditions method returns true if Resty's default retry conditions
are enabled otherwise false
Default value is `true`
func (*Client)
IsSaveResponse
¶
func (c *
Client
) IsSaveResponse()
bool
IsSaveResponse method returns true if the save response is set to true; otherwise, false
func (*Client)
IsTrace
¶
func (c *
Client
) IsTrace()
bool
IsTrace method returns true if the trace is enabled on the client instance; otherwise, it returns false.
func (*Client)
LoadBalancer
¶
func (c *
Client
) LoadBalancer()
LoadBalancer
LoadBalancer method returns the request load balancer instance from the client
instance. Otherwise returns nil.
func (*Client)
Logger
¶
func (c *
Client
) Logger()
Logger
Logger method returns the logger instance used by the client instance.
func (*Client)
NewRequest
¶
func (c *
Client
) NewRequest() *
Request
NewRequest method is an alias for method `R()`.
func (*Client)
OnClose
¶
func (c *
Client
) OnClose(h
CloseHook
) *
Client
OnClose method adds a callback that will be run whenever the client is closed.
The hooks are executed in the order they were registered.
func (*Client)
OnDebugLog
¶
func (c *
Client
) OnDebugLog(dlc
DebugLogCallbackFunc
) *
Client
OnDebugLog method sets the debug log callback function to the client instance.
Registered callback gets called before the Resty logs the information.
func (*Client)
OnError
¶
func (c *
Client
) OnError(h
ErrorHook
) *
Client
OnError method adds a callback that will be run whenever a request execution fails.
This is called after all retries have been attempted (if any).
If there was a response from the server, the error will be wrapped in
ResponseError
which has the last response received from the server.
client.OnError(func(req *resty.Request, err error) {
if v, ok := err.(*resty.ResponseError); ok {
// Do something with v.Response
}
// Log the error, increment a metric, etc...
})
Out of the
Client.OnSuccess
,
Client.OnError
,
Client.OnInvalid
,
Client.OnPanic
callbacks, exactly one set will be invoked for each call to
Request.Execute
that completes.
NOTE:
Do not use
Client
setter methods within OnError hooks; deadlock will happen.
func (*Client)
OnInvalid
¶
func (c *
Client
) OnInvalid(h
ErrorHook
) *
Client
OnInvalid method adds a callback that will be run whenever a request execution
fails before it starts because the request is invalid.
Out of the
Client.OnSuccess
,
Client.OnError
,
Client.OnInvalid
,
Client.OnPanic
callbacks, exactly one set will be invoked for each call to
Request.Execute
that completes.
NOTE:
Do not use
Client
setter methods within OnInvalid hooks; deadlock will happen.
func (*Client)
OnPanic
¶
func (c *
Client
) OnPanic(h
ErrorHook
) *
Client
OnPanic method adds a callback that will be run whenever a request execution
panics.
Out of the
Client.OnSuccess
,
Client.OnError
,
Client.OnInvalid
,
Client.OnPanic
callbacks, exactly one set will be invoked for each call to
Request.Execute
that completes.
If an
Client.OnSuccess
,
Client.OnError
, or
Client.OnInvalid
callback panics,
then exactly one rule can be violated.
NOTE:
Do not use
Client
setter methods within OnPanic hooks; deadlock will happen.
func (*Client)
OnSuccess
¶
func (c *
Client
) OnSuccess(h
SuccessHook
) *
Client
OnSuccess method adds a callback that will be run whenever a request execution
succeeds.  This is called after all retries have been attempted (if any).
Out of the
Client.OnSuccess
,
Client.OnError
,
Client.OnInvalid
,
Client.OnPanic
callbacks, exactly one set will be invoked for each call to
Request.Execute
that completes.
NOTE:
Do not use
Client
setter methods within OnSuccess hooks; deadlock will happen.
func (*Client)
OutputDirectory
¶
func (c *
Client
) OutputDirectory()
string
OutputDirectory method returns the output directory value from the client.
func (*Client)
PathParams
¶
func (c *
Client
) PathParams() map[
string
]
string
PathParams method returns the path parameters from the client.
pathParams := client.PathParams()
func (*Client)
ProxyURL
¶
func (c *
Client
) ProxyURL() *
url
.
URL
ProxyURL method returns the proxy URL if set otherwise nil.
func (*Client)
QueryParams
¶
func (c *
Client
) QueryParams()
url
.
Values
QueryParams method returns all query parameters and their values from the client instance.
func (*Client)
R
¶
func (c *
Client
) R() *
Request
R method creates a new request instance; it's used for Get, Post, Put, Delete, Patch, Head, Options, etc.
func (*Client)
RemoveProxy
¶
func (c *
Client
) RemoveProxy() *
Client
RemoveProxy method removes the proxy configuration from the Resty client
client.RemoveProxy()
func (*Client)
ResponseBodyLimit
¶
func (c *
Client
) ResponseBodyLimit()
int64
ResponseBodyLimit method returns the value max body size limit in bytes from
the client instance.
func (*Client)
ResponseBodyUnlimitedReads
¶
func (c *
Client
) ResponseBodyUnlimitedReads()
bool
ResponseBodyUnlimitedReads method returns true if enabled. Otherwise, it returns false
func (*Client)
RetryConditions
¶
func (c *
Client
) RetryConditions() []
RetryConditionFunc
RetryConditions method returns all the retry condition functions.
func (*Client)
RetryCount
¶
func (c *
Client
) RetryCount()
int
RetryCount method returns the retry count value from the client instance.
func (*Client)
RetryHooks
¶
func (c *
Client
) RetryHooks() []
RetryHookFunc
RetryHooks method returns all the retry hook functions.
func (*Client)
RetryMaxWaitTime
¶
func (c *
Client
) RetryMaxWaitTime()
time
.
Duration
RetryMaxWaitTime method returns the retry max wait time that is used to sleep
before retrying the request.
func (*Client)
RetryStrategy
¶
func (c *
Client
) RetryStrategy()
RetryStrategyFunc
RetryStrategy method returns the retry strategy function; otherwise, it is nil.
See
Client.SetRetryStrategy
func (*Client)
RetryWaitTime
¶
func (c *
Client
) RetryWaitTime()
time
.
Duration
RetryWaitTime method returns the retry wait time that is used to sleep before
retrying the request.
func (*Client)
Scheme
¶
func (c *
Client
) Scheme()
string
Scheme method returns custom scheme value from the client.
scheme := client.Scheme()
func (*Client)
SetAllowMethodDeletePayload
¶
func (c *
Client
) SetAllowMethodDeletePayload(allow
bool
) *
Client
SetAllowMethodDeletePayload method allows the DELETE method with payload on the Resty client.
By default, Resty does not allow.
client.SetAllowMethodDeletePayload(true)
More info, refer to GH#881
It can be overridden at the request level. See
Request.SetAllowMethodDeletePayload
func (*Client)
SetAllowMethodGetPayload
¶
func (c *
Client
) SetAllowMethodGetPayload(allow
bool
) *
Client
SetAllowMethodGetPayload method allows the GET method with payload on the Resty client.
By default, Resty does not allow.
client.SetAllowMethodGetPayload(true)
It can be overridden at the request level. See
Request.SetAllowMethodGetPayload
func (*Client)
SetAllowNonIdempotentRetry
¶
func (c *
Client
) SetAllowNonIdempotentRetry(b
bool
) *
Client
SetAllowNonIdempotentRetry method is used to enable/disable non-idempotent HTTP
methods retry. By default, Resty only allows idempotent HTTP methods, see
RFC 9110 Section 9.2.2
,
RFC 9110 Section 18.2
It can be overridden at request level, see
Request.SetAllowNonIdempotentRetry
func (*Client)
SetAuthScheme
¶
func (c *
Client
) SetAuthScheme(scheme
string
) *
Client
SetAuthScheme method sets the auth scheme type in the HTTP request. For Example:
Authorization: <auth-scheme-value> <auth-token-value>
For Example: To set the scheme to use OAuth
client.SetAuthScheme("OAuth")
This auth scheme gets added to all the requests raised from this client instance.
Also, it can be overridden at the request level.
Information about auth schemes can be found in
RFC 7235
, IANA
HTTP Auth schemes
.
See
Request.SetAuthScheme
.
func (*Client)
SetAuthToken
¶
func (c *
Client
) SetAuthToken(token
string
) *
Client
SetAuthToken method sets the auth token of the `Authorization` header for all HTTP requests.
The default auth scheme is `Bearer`; it can be customized with the method
Client.SetAuthScheme
. For Example:
Authorization: <auth-scheme> <auth-token-value>
For Example: To set auth token BC594900518B4F7EAC75BD37F019E08FBC594900518B4F7EAC75BD37F019E08F
client.SetAuthToken("BC594900518B4F7EAC75BD37F019E08FBC594900518B4F7EAC75BD37F019E08F")
This auth token gets added to all the requests raised from this client instance.
Also, it can be overridden at the request level.
See
Request.SetAuthToken
.
func (*Client)
SetBaseURL
¶
func (c *
Client
) SetBaseURL(url
string
) *
Client
SetBaseURL method sets the Base URL in the client instance. It will be used with a request
raised from this client with a relative URL
// Setting HTTP address
client.SetBaseURL("http://myjeeva.com")

// Setting HTTPS address
client.SetBaseURL("https://myjeeva.com")
func (*Client)
SetBasicAuth
¶
func (c *
Client
) SetBasicAuth(username, password
string
) *
Client
SetBasicAuth method sets the basic authentication header in the HTTP request. For Example:
Authorization: Basic <base64-encoded-value>
For Example: To set the header for username "go-resty" and password "welcome"
client.SetBasicAuth("go-resty", "welcome")
This basic auth information is added to all requests from this client instance.
It can also be overridden at the request level.
See
Request.SetBasicAuth
.
func (*Client)
SetCertificateFromFile
¶
func (c *
Client
) SetCertificateFromFile(certFilePath, certKeyFilePath
string
) *
Client
SetCertificateFromFile method helps to set client certificates into Resty
from cert and key files to perform SSL client authentication
client.SetCertificateFromFile("certs/client.pem", "certs/client.key")
func (*Client)
SetCertificateFromString
¶
func (c *
Client
) SetCertificateFromString(certStr, certKeyStr
string
) *
Client
SetCertificateFromString method helps to set client certificates into Resty
from string to perform SSL client authentication
myClientCertStr := `-----BEGIN CERTIFICATE-----
... cert content ...
-----END CERTIFICATE-----`

myClientCertKeyStr := `-----BEGIN PRIVATE KEY-----
... cert key content ...
-----END PRIVATE KEY-----`

client.SetCertificateFromString(myClientCertStr, myClientCertKeyStr)
func (*Client)
SetCertificates
¶
func (c *
Client
) SetCertificates(certs ...
tls
.
Certificate
) *
Client
SetCertificates method helps to conveniently set a slice of client certificates
into Resty to perform SSL client authentication
cert, err := tls.LoadX509KeyPair("certs/client.pem", "certs/client.key")
if err != nil {
log.Printf("ERROR client certificate/key parsing error: %v", err)
return
}

client.SetCertificates(cert)
func (*Client)
SetCircuitBreaker
¶
func (c *
Client
) SetCircuitBreaker(b *
CircuitBreaker
) *
Client
SetCircuitBreaker method sets the Circuit Breaker instance into the client.
It is used to prevent the client from sending requests that are likely to fail.
For Example: To use the default Circuit Breaker:
client.SetCircuitBreaker(NewCircuitBreaker())
func (*Client)
SetClientRootCertificateFromString
¶
func (c *
Client
) SetClientRootCertificateFromString(pemCerts
string
) *
Client
SetClientRootCertificateFromString method helps to add a client root certificate
from the string into the Resty client
myClientRootCertStr := `-----BEGIN CERTIFICATE-----
... cert content ...
-----END CERTIFICATE-----`

client.SetClientRootCertificateFromString(myClientRootCertStr)
func (*Client)
SetClientRootCertificates
¶
func (c *
Client
) SetClientRootCertificates(pemFilePaths ...
string
) *
Client
SetClientRootCertificates method helps to add one or more client root
certificate files into the Resty client
// one pem file path
client.SetClientRootCertificates("/path/to/client-root/pemFile.pem")

// one or more pem file path(s)
client.SetClientRootCertificates(
"/path/to/client-root/pemFile1.pem",
"/path/to/client-root/pemFile2.pem"
"/path/to/client-root/pemFile3.pem"
)

// if you happen to have string slices
client.SetClientRootCertificates(certs...)
func (*Client)
SetClientRootCertificatesWatcher
¶
func (c *
Client
) SetClientRootCertificatesWatcher(options *
CertWatcherOptions
, pemFilePaths ...
string
) *
Client
SetClientRootCertificatesWatcher method enables dynamic reloading of one or more client root certificate files.
It is designed for scenarios involving long-running Resty clients where certificates may be renewed.
client.SetClientRootCertificatesWatcher(
&resty.CertWatcherOptions{
PoolInterval: 24 * time.Hour,
},
"client-root-ca.pem",
)
func (*Client)
SetCloseConnection
¶
func (c *
Client
) SetCloseConnection(close
bool
) *
Client
SetCloseConnection method sets variable `Close` in HTTP request struct with the given
value. More info:
https://golang.org/src/net/http/request.go
It can be overridden at the request level, see
Request.SetCloseConnection
func (*Client)
SetContentDecompresserKeys
¶
func (c *
Client
) SetContentDecompresserKeys(keys []
string
) *
Client
SetContentDecompresserKeys method sets given Content-Encoding (
RFC 9110
) directives into the client instance.
It checks the given Content-Encoding exists in the
ContentDecompresser
list before assigning it,
if it does not exist, it will skip that directive.
Use this method to overwrite the default order. If a new content Decompresser is added,
that directive will be the first.
func (*Client)
SetContext
¶
func (c *
Client
) SetContext(ctx
context
.
Context
) *
Client
SetContext method sets the given
context.Context
in the client instance and
it gets added to
Request
raised from this instance.
func (*Client)
SetCookie
¶
func (c *
Client
) SetCookie(hc *
http
.
Cookie
) *
Client
SetCookie method appends a single cookie to the client instance.
These cookies will be added to all the requests from this client instance.
client.SetCookie(&http.Cookie{
Name:"go-resty",
Value:"This is cookie value",
})
func (*Client)
SetCookieJar
¶
func (c *
Client
) SetCookieJar(jar
http
.
CookieJar
) *
Client
SetCookieJar method sets custom
http.CookieJar
in the resty client. It's a way to override the default.
For Example, sometimes we don't want to save cookies in API mode so that we can remove the default
CookieJar in resty client.
client.SetCookieJar(nil)
func (*Client)
SetCookies
¶
func (c *
Client
) SetCookies(cs []*
http
.
Cookie
) *
Client
SetCookies method sets an array of cookies in the client instance.
These cookies will be added to all the requests from this client instance.
cookies := []*http.Cookie{
&http.Cookie{
Name:"go-resty-1",
Value:"This is cookie 1 value",
},
&http.Cookie{
Name:"go-resty-2",
Value:"This is cookie 2 value",
},
}

// Setting a cookies into resty
client.SetCookies(cookies)
func (*Client)
SetDebug
¶
func (c *
Client
) SetDebug(d
bool
) *
Client
SetDebug method enables the debug mode on the Resty client. The client logs details
of every request and response.
client.SetDebug(true)
// OR
client.EnableDebug()
Also, it can be enabled at the request level for a particular request; see
Request.SetDebug
.
For
Request
, it logs information such as HTTP verb, Relative URL path,
Host, Headers, and Body if it has one.
For
Response
, it logs information such as Status, Response Time, Headers,
and Body if it has one.
func (*Client)
SetDebugBodyLimit
¶
func (c *
Client
) SetDebugBodyLimit(sl
int
) *
Client
SetDebugBodyLimit sets the maximum size in bytes for which the response and
request body will be logged in debug mode.
client.SetDebugBodyLimit(1000000)
func (*Client)
SetDebugLogCurlCmd
¶
func (c *
Client
) SetDebugLogCurlCmd(b
bool
) *
Client
SetDebugLogCurlCmd method enables the curl command to be logged in the debug log.
It can be overridden at the request level; see
Request.SetDebugLogCurlCmd
func (*Client)
SetDebugLogFormatter
¶
func (c *
Client
) SetDebugLogFormatter(df
DebugLogFormatterFunc
) *
Client
SetDebugLogFormatter method sets the Resty debug log formatter to the client instance.
func (*Client)
SetDigestAuth
¶
func (c *
Client
) SetDigestAuth(username, password
string
) *
Client
SetDigestAuth method sets the Digest Auth transport with provided credentials in the client.
If a server responds with 401 and sends a Digest challenge in the header `WWW-Authenticate`,
the request will be resent with the appropriate digest `Authorization` header.
For Example: To set the Digest scheme with user "Mufasa" and password "Circle Of Life"
client.SetDigestAuth("Mufasa", "Circle Of Life")
Information about Digest Access Authentication can be found in
RFC 7616
.
NOTE:
On the QOP `auth-int` scenario, the request body is read into memory to
compute the body hash that increases memory usage.
Create a dedicated client instance to use digest auth,
as it does digest auth for all the requests raised by the client.
func (*Client)
SetDisableWarn
¶
func (c *
Client
) SetDisableWarn(d
bool
) *
Client
SetDisableWarn method disables the warning log message on the Resty client.
For example, Resty warns users when BasicAuth is used in non-TLS mode.
client.SetDisableWarn(true)
func (*Client)
SetDoNotParseResponse
¶
func (c *
Client
) SetDoNotParseResponse(notParse
bool
) *
Client
SetDoNotParseResponse method instructs Resty not to parse the response body automatically.
Resty exposes the raw response body as
io.ReadCloser
. If you use it, do not
forget to close the body, otherwise, you might get into connection leaks, and connection
reuse may not happen.
NOTE: The default
Response
middlewares are not executed when using this option. User
takes over the control of handling response body from Resty.
func (*Client)
SetError
¶
func (c *
Client
) SetError(v
any
) *
Client
SetError method registers the global or client common `Error` object into Resty.
It is used for automatic unmarshalling if the response status code is greater than 399 and
content type is JSON or XML. It can be a pointer or a non-pointer.
client.SetError(&Error{})
// OR
client.SetError(Error{})
func (*Client)
SetFormData
¶
func (c *
Client
) SetFormData(data map[
string
]
string
) *
Client
SetFormData method sets Form parameters and their values in the client instance.
The request content type would be set as `application/x-www-form-urlencoded`.
The client-level form data gets added to all the requests. Also, it can be
overridden at the request level.
See
Request.SetFormData
.
client.SetFormData(map[string]string{
"access_token": "BC594900-518B-4F7E-AC75-BD37F019E08F",
"user_id": "3455454545",
})
func (*Client)
SetGenerateCurlCmd
¶
func (c *
Client
) SetGenerateCurlCmd(b
bool
) *
Client
SetGenerateCurlCmd method is used to turn on/off the generate curl command at the
client instance level.
By default, Resty does not log the curl command in the debug log since it has the potential
to leak sensitive data unless explicitly enabled via
Client.SetDebugLogCurlCmd
or
Request.SetDebugLogCurlCmd
.
NOTE: Use with care.
Potential to leak sensitive data from
Request
and
Response
in the debug log
when the debug log option is enabled.
Additional memory usage since the request body was reread.
curl body is not generated for
io.Reader
and multipart request flow.
It can be overridden at the request level; see
Request.SetGenerateCurlCmd
func (*Client)
SetHeader
¶
func (c *
Client
) SetHeader(header, value
string
) *
Client
SetHeader method sets a single header and its value in the client instance.
These headers will be applied to all requests raised from the client instance.
Also, it can be overridden by request-level header options.
For Example: To set `Content-Type` and `Accept` as `application/json`
client.
SetHeader("Content-Type", "application/json").
SetHeader("Accept", "application/json")
See
Request.SetHeader
or
Request.SetHeaders
.
func (*Client)
SetHeaderAuthorizationKey
¶
func (c *
Client
) SetHeaderAuthorizationKey(k
string
) *
Client
SetHeaderAuthorizationKey method sets the given HTTP header name for Authorization in the client instance.
It can be overridden at the request level; see
Request.SetHeaderAuthorizationKey
.
client.SetHeaderAuthorizationKey("X-Custom-Authorization")
func (*Client)
SetHeaderVerbatim
¶
func (c *
Client
) SetHeaderVerbatim(header, value
string
) *
Client
SetHeaderVerbatim method is used to set the HTTP header key and value verbatim in the current request.
It is typically helpful for legacy applications or servers that require HTTP headers in a certain way
For Example: To set header key as `all_lowercase`, `UPPERCASE`, and `x-cloud-trace-id`
client.
SetHeaderVerbatim("all_lowercase", "available").
SetHeaderVerbatim("UPPERCASE", "available").
SetHeaderVerbatim("x-cloud-trace-id", "798e94019e5fc4d57fbb8901eb4c6cae")
See
Request.SetHeaderVerbatim
.
func (*Client)
SetHeaders
¶
func (c *
Client
) SetHeaders(headers map[
string
]
string
) *
Client
SetHeaders method sets multiple headers and their values at one go, and
these headers will be applied to all requests raised from the client instance.
Also, it can be overridden at request-level headers options.
For Example: To set `Content-Type` and `Accept` as `application/json`
client.SetHeaders(map[string]string{
"Content-Type": "application/json",
"Accept": "application/json",
})
See
Request.SetHeaders
or
Request.SetHeader
.
func (*Client)
SetJSONEscapeHTML
¶
func (c *
Client
) SetJSONEscapeHTML(b
bool
) *
Client
SetJSONEscapeHTML method enables or disables the HTML escape on JSON marshal.
By default, escape HTML is `true`.
NOTE: This option only applies to the standard JSON Marshaller used by Resty.
It can be overridden at the request level, see
Request.SetJSONEscapeHTML
func (*Client)
SetLoadBalancer
¶
func (c *
Client
) SetLoadBalancer(b
LoadBalancer
) *
Client
SetLoadBalancer method is used to set the new request load balancer into the client.
func (*Client)
SetLogger
¶
func (c *
Client
) SetLogger(l
Logger
) *
Client
SetLogger method sets given writer for logging Resty request and response details.
Compliant to interface
resty.Logger
func (*Client)
SetOutputDirectory
¶
func (c *
Client
) SetOutputDirectory(dirPath
string
) *
Client
SetOutputDirectory method sets the output directory for saving HTTP responses in a file.
Resty creates one if the output directory does not exist. This setting is optional,
if you plan to use the absolute path in
Request.SetOutputFileName
and can used together.
client.SetOutputDirectory("/save/http/response/here")
func (*Client)
SetPathParam
¶
func (c *
Client
) SetPathParam(param, value
string
) *
Client
SetPathParam method sets a single URL path key-value pair in the
Resty client instance.
client.SetPathParam("userId", "sample@sample.com")

Result:
URL - /v1/users/{userId}/details
Composed URL - /v1/users/sample@sample.com/details
It replaces the value of the key while composing the request URL.
The value will be escaped using
url.PathEscape
function.
It can be overridden at the request level,
see
Request.SetPathParam
or
Request.SetPathParams
func (*Client)
SetPathParams
¶
func (c *
Client
) SetPathParams(params map[
string
]
string
) *
Client
SetPathParams method sets multiple URL path key-value pairs at one go in the
Resty client instance.
client.SetPathParams(map[string]string{
"userId":       "sample@sample.com",
"subAccountId": "100002",
"path":         "groups/developers",
})

Result:
URL - /v1/users/{userId}/{subAccountId}/{path}/details
Composed URL - /v1/users/sample@sample.com/100002/groups%2Fdevelopers/details
It replaces the value of the key while composing the request URL.
The values will be escaped using
url.PathEscape
function.
It can be overridden at the request level,
see
Request.SetPathParam
or
Request.SetPathParams
func (*Client)
SetProxy
¶
func (c *
Client
) SetProxy(proxyURL
string
) *
Client
SetProxy method sets the Proxy URL and Port for the Resty client.
// HTTP/HTTPS proxy
client.SetProxy("http://proxyserver:8888")

// SOCKS5 Proxy
client.SetProxy("socks5://127.0.0.1:1080")
OR you could also set Proxy via environment variable, refer to
http.ProxyFromEnvironment
func (*Client)
SetQueryParam
¶
func (c *
Client
) SetQueryParam(param, value
string
) *
Client
SetQueryParam method sets a single parameter and its value in the client instance.
It will be formed as a query string for the request.
For Example: `search=kitchen%20papers&size=large`
In the URL after the `?` mark. These query params will be added to all the requests raised from
this client instance. Also, it can be overridden at the request level.
See
Request.SetQueryParam
or
Request.SetQueryParams
.
client.
SetQueryParam("search", "kitchen papers").
SetQueryParam("size", "large")
func (*Client)
SetQueryParams
¶
func (c *
Client
) SetQueryParams(params map[
string
]
string
) *
Client
SetQueryParams method sets multiple parameters and their values at one go in the client instance.
It will be formed as a query string for the request.
For Example: `search=kitchen%20papers&size=large`
In the URL after the `?` mark. These query params will be added to all the requests raised from this
client instance. Also, it can be overridden at the request level.
See
Request.SetQueryParams
or
Request.SetQueryParam
.
client.SetQueryParams(map[string]string{
"search": "kitchen papers",
"size": "large",
})
func (*Client)
SetRawPathParam
¶
func (c *
Client
) SetRawPathParam(param, value
string
) *
Client
SetRawPathParam method sets a single URL path key-value pair in the
Resty client instance without path escape.
client.SetRawPathParam("path", "groups/developers")

Result:
URL - /v1/users/{path}/details
Composed URL - /v1/users/groups/developers/details
It replaces the value of the key while composing the request URL.
The value will be used as-is, no path escape applied.
It can be overridden at the request level,
see
Request.SetRawPathParam
or
Request.SetRawPathParams
func (*Client)
SetRawPathParams
¶
func (c *
Client
) SetRawPathParams(params map[
string
]
string
) *
Client
SetRawPathParams method sets multiple URL path key-value pairs at one go in the
Resty client instance without path escape.
client.SetRawPathParams(map[string]string{
"userId":       "sample@sample.com",
"subAccountId": "100002",
"path":         "groups/developers",
})

Result:
URL - /v1/users/{userId}/{subAccountId}/{path}/details
Composed URL - /v1/users/sample@sample.com/100002/groups/developers/details
It replaces the value of the key while composing the request URL.
The value will be used as-is, no path escape applied.
It can be overridden at the request level,
see
Request.SetRawPathParam
or
Request.SetRawPathParams
func (*Client)
SetRedirectPolicy
¶
func (c *
Client
) SetRedirectPolicy(policies ...
RedirectPolicy
) *
Client
SetRedirectPolicy method sets the redirect policy for the client. Resty provides ready-to-use
redirect policies. Wanna create one for yourself, refer to `redirect.go`.
client.SetRedirectPolicy(resty.FlexibleRedirectPolicy(20))

// Need multiple redirect policies together
client.SetRedirectPolicy(resty.FlexibleRedirectPolicy(20), resty.DomainCheckRedirectPolicy("host1.com", "host2.net"))
NOTE: It overwrites the previous redirect policies in the client instance.
func (*Client)
SetRequestMiddlewares
¶
func (c *
Client
) SetRequestMiddlewares(middlewares ...
RequestMiddleware
) *
Client
SetRequestMiddlewares method allows Resty users to override the default request
middlewares sequence
client.SetRequestMiddlewares(
Custom1RequestMiddleware,
Custom2RequestMiddleware,
resty.PrepareRequestMiddleware, // after this, `Request.RawRequest` instance is available
Custom3RequestMiddleware,
Custom4RequestMiddleware,
)
See,
Client.AddRequestMiddleware
NOTE:
It overwrites the existing request middleware list.
Be sure to include Resty request middlewares in the request chain at the appropriate spot.
func (*Client)
SetResponseBodyLimit
¶
func (c *
Client
) SetResponseBodyLimit(v
int64
) *
Client
SetResponseBodyLimit method sets a maximum body size limit in bytes on response,
avoid reading too much data to memory.
Client will return
resty.ErrResponseBodyTooLarge
if the body size of the body
in the uncompressed response is larger than the limit.
Body size limit will not be enforced in the following cases:
ResponseBodyLimit <= 0, which is the default behavior.
Request.SetOutputFileName
is called to save response data to the file.
"DoNotParseResponse" is set for client or request.
It can be overridden at the request level; see
Request.SetResponseBodyLimit
func (*Client)
SetResponseBodyUnlimitedReads
¶
func (c *
Client
) SetResponseBodyUnlimitedReads(b
bool
) *
Client
SetResponseBodyUnlimitedReads method is to turn on/off the response body in memory
that provides an ability to do unlimited reads.
It can be overridden at the request level; see
Request.SetResponseBodyUnlimitedReads
Unlimited reads are possible in a few scenarios, even without enabling it.
When debug mode is enabled
NOTE: Use with care
Turning on this feature keeps the response body in memory, which might cause additional memory usage.
func (*Client)
SetResponseMiddlewares
¶
func (c *
Client
) SetResponseMiddlewares(middlewares ...
ResponseMiddleware
) *
Client
SetResponseMiddlewares method allows Resty users to override the default response
middlewares sequence
client.SetResponseMiddlewares(
Custom1ResponseMiddleware,
Custom2ResponseMiddleware,
resty.AutoParseResponseMiddleware, // before this, the body is not read except on the debug flow
Custom3ResponseMiddleware,
resty.SaveToFileResponseMiddleware, // See, Request.SetOutputFileName, Request.SetSaveResponse
Custom4ResponseMiddleware,
Custom5ResponseMiddleware,
)
See,
Client.AddResponseMiddleware
NOTE:
It overwrites the existing response middleware list.
Be sure to include Resty response middlewares in the response chain at the appropriate spot.
func (*Client)
SetRetryCount
¶
func (c *
Client
) SetRetryCount(count
int
) *
Client
SetRetryCount method enables retry on Resty client and allows you
to set no. of retry count.
first attempt + retry count = total attempts
See
Request.SetRetryStrategy
NOTE:
By default, Resty only does retry on idempotent HTTP verb,
RFC 9110 Section 9.2.2
,
RFC 9110 Section 18.2
func (*Client)
SetRetryDefaultConditions
¶
func (c *
Client
) SetRetryDefaultConditions(b
bool
) *
Client
SetRetryDefaultConditions method is used to enable/disable the Resty's default
retry conditions
It can be overridden at request level, see
Request.SetRetryDefaultConditions
func (*Client)
SetRetryMaxWaitTime
¶
func (c *
Client
) SetRetryMaxWaitTime(maxWaitTime
time
.
Duration
) *
Client
SetRetryMaxWaitTime method sets the max wait time for sleep before retrying
Default is 2 seconds.
func (*Client)
SetRetryStrategy
¶
func (c *
Client
) SetRetryStrategy(rs
RetryStrategyFunc
) *
Client
SetRetryStrategy method used to set the custom Retry strategy into Resty client,
it is used to get wait time before each retry. It can be overridden at request
level, see
Request.SetRetryStrategy
Default (nil) implies exponential backoff with a jitter strategy
func (*Client)
SetRetryWaitTime
¶
func (c *
Client
) SetRetryWaitTime(waitTime
time
.
Duration
) *
Client
SetRetryWaitTime method sets the default wait time for sleep before retrying
Default is 100 milliseconds.
func (*Client)
SetRootCertificateFromString
¶
func (c *
Client
) SetRootCertificateFromString(pemCerts
string
) *
Client
SetRootCertificateFromString method helps to add root certificate from the string
into the Resty client
myRootCertStr := `-----BEGIN CERTIFICATE-----
... cert content ...
-----END CERTIFICATE-----`

client.SetRootCertificateFromString(myRootCertStr)
func (*Client)
SetRootCertificates
¶
func (c *
Client
) SetRootCertificates(pemFilePaths ...
string
) *
Client
SetRootCertificates method helps to add one or more root certificate files
into the Resty client
// one pem file path
client.SetRootCertificates("/path/to/root/pemFile.pem")

// one or more pem file path(s)
client.SetRootCertificates(
"/path/to/root/pemFile1.pem",
"/path/to/root/pemFile2.pem"
"/path/to/root/pemFile3.pem"
)

// if you happen to have string slices
client.SetRootCertificates(certs...)
func (*Client)
SetRootCertificatesWatcher
¶
func (c *
Client
) SetRootCertificatesWatcher(options *
CertWatcherOptions
, pemFilePaths ...
string
) *
Client
SetRootCertificatesWatcher method enables dynamic reloading of one or more root certificate files.
It is designed for scenarios involving long-running Resty clients where certificates may be renewed.
client.SetRootCertificatesWatcher(
&resty.CertWatcherOptions{
PoolInterval: 24 * time.Hour,
},
"root-ca.pem",
)
func (*Client)
SetSaveResponse
¶
func (c *
Client
) SetSaveResponse(save
bool
) *
Client
SetSaveResponse method used to enable the save response option at the client level for
all requests
client.SetSaveResponse(true)
Resty determines the save filename in the following order -
Request.SetOutputFileName
Content-Disposition header
Request URL using
path.Base
Request URL hostname if path is empty or "/"
It can be overridden at request level, see
Request.SetSaveResponse
func (*Client)
SetScheme
¶
func (c *
Client
) SetScheme(scheme
string
) *
Client
SetScheme method sets a custom scheme for the Resty client. It's a way to override the default.
client.SetScheme("http")
func (*Client)
SetTLSClientConfig
¶
func (c *
Client
) SetTLSClientConfig(tlsConfig *
tls
.
Config
) *
Client
SetTLSClientConfig method sets TLSClientConfig for underlying client Transport.
Values supported by
https://pkg.go.dev/crypto/tls#Config
can be configured.
// Disable SSL cert verification for local development
client.SetTLSClientConfig(&tls.Config{
InsecureSkipVerify: true
})
NOTE: This method overwrites existing
http.Transport.TLSClientConfig
func (*Client)
SetTimeout
¶
func (c *
Client
) SetTimeout(timeout
time
.
Duration
) *
Client
SetTimeout method is used to set a timeout for a request raised by the client.
client.SetTimeout(1 * time.Minute)
It can be overridden at the request level. See
Request.SetTimeout
NOTE: Resty uses
context.WithTimeout
on the request, it does not use
http.Client
.Timeout
func (*Client)
SetTrace
¶
func (c *
Client
) SetTrace(t
bool
) *
Client
SetTrace method is used to turn on/off the trace capability in the Resty client
Refer to
Client.EnableTrace
or
Client.DisableTrace
.
Also, see
Request.SetTrace
func (*Client)
SetTransport
¶
func (c *
Client
) SetTransport(transport
http
.
RoundTripper
) *
Client
SetTransport method sets custom
http.Transport
or any
http.RoundTripper
compatible interface implementation in the Resty client.
transport := &http.Transport{
// something like Proxying to httptest.Server, etc...
Proxy: func(req *http.Request) (*url.URL, error) {
return url.Parse(server.URL)
},
}
client.SetTransport(transport)
NOTE:
If transport is not the type of
http.Transport
, you may lose the
ability to set a few Resty client settings. However, if you implement
TLSClientConfiger
interface, then TLS client config is possible to set.
It overwrites the Resty client transport instance and its configurations.
func (*Client)
SetUnescapeQueryParams
¶
func (c *
Client
) SetUnescapeQueryParams(unescape
bool
) *
Client
SetUnescapeQueryParams method sets the choice of unescape query parameters for the request URL.
To prevent broken URL, Resty replaces space (" ") with "+" in the query parameters.
See
Request.SetUnescapeQueryParams
NOTE: Request failure is possible due to non-standard usage of Unescaped Query Parameters.
func (*Client)
TLSClientConfig
¶
func (c *
Client
) TLSClientConfig() *
tls
.
Config
TLSClientConfig method returns the
tls.Config
from underlying client transport
otherwise returns nil
func (*Client)
Timeout
¶
func (c *
Client
) Timeout()
time
.
Duration
Timeout method returns the timeout duration value from the client
func (*Client)
Transport
¶
func (c *
Client
) Transport()
http
.
RoundTripper
Transport method returns underlying client transport referance as-is
i.e.,
http.RoundTripper
type
CloseHook
¶
type CloseHook func()
CloseHook type is for reacting to client closing
type
ContentDecompresser
¶
type ContentDecompresser func(
io
.
ReadCloser
) (
io
.
ReadCloser
,
error
)
ContentDecompresser type is for decompressing response body based on header Content-Encoding
(
RFC 9110
)
For example, gzip, deflate, etc.
type
ContentTypeDecoder
¶
type ContentTypeDecoder func(
io
.
Reader
,
any
)
error
ContentTypeDecoder type is for decoding the response body based on header Content-Type
type
ContentTypeEncoder
¶
type ContentTypeEncoder func(
io
.
Writer
,
any
)
error
ContentTypeEncoder type is for encoding the request body based on header Content-Type
type
DebugLog
¶
type DebugLog struct {
Request   *
DebugLogRequest
`json:"request"`
Response  *
DebugLogResponse
`json:"response"`
TraceInfo *
TraceInfo
`json:"trace_info"`
}
DebugLog struct is used to collect details from Resty request and response
for debug logging callback purposes.
type
DebugLogCallbackFunc
¶
type DebugLogCallbackFunc func(*
DebugLog
)
DebugLogCallbackFunc function type is for request and response debug log callback purposes.
It gets called before Resty logs it
type
DebugLogFormatterFunc
¶
type DebugLogFormatterFunc func(*
DebugLog
)
string
DebugLogFormatterFunc function type is used to implement debug log formatting.
See out of the box [DebugLogStringFormatter],
DebugLogJSONFormatter
type
DebugLogRequest
¶
type DebugLogRequest struct {
Host
string
`json:"host"`
URI
string
`json:"uri"`
Method
string
`json:"method"`
Proto
string
`json:"proto"`
Header
http
.
Header
`json:"header"`
CurlCmd
string
`json:"curl_cmd"`
RetryTraceID
string
`json:"retry_trace_id"`
Attempt
int
`json:"attempt"`
Body
string
`json:"body"`
}
DebugLogRequest type used to capture debug info about the
Request
.
type
DebugLogResponse
¶
type DebugLogResponse struct {
StatusCode
int
`json:"status_code"`
Status
string
`json:"status"`
Proto
string
`json:"proto"`
ReceivedAt
time
.
Time
`json:"received_at"`
Duration
time
.
Duration
`json:"duration"`
Size
int64
`json:"size"`
Header
http
.
Header
`json:"header"`
Body
string
`json:"body"`
}
DebugLogResponse type used to capture debug info about the
Response
.
type
ErrorHook
¶
type ErrorHook func(*
Request
,
error
)
ErrorHook type is for reacting to request errors, called after all retries were attempted
type
Event
¶
type Event struct {
ID
string
Name
string
Data
string
}
Event struct represents the event details from the Server-Sent Events(SSE) stream
type
EventErrorFunc
¶
type EventErrorFunc func(
error
)
EventErrorFunc is a callback function type used to receive notification
when an error occurs with
EventSource
processing
type
EventMessageFunc
¶
type EventMessageFunc func(
any
)
EventMessageFunc is a callback function type used to receive event details
from the Server-Sent Events(SSE) stream
type
EventOpenFunc
¶
type EventOpenFunc func(url
string
, respHdr
http
.
Header
)
EventOpenFunc is a callback function type used to receive notification
when Resty establishes a connection with the server for the
Server-Sent Events(SSE)
type
EventRequestFailureFunc
¶
type EventRequestFailureFunc func(err
error
, res *
http
.
Response
)
EventRequestFailureFunc is a callback function type used to receive event
details from the Server-Sent Events(SSE) request failure
type
EventSource
¶
type EventSource struct {
// contains filtered or unexported fields
}
EventSource struct implements the Server-Sent Events(SSE)
specification
to receive
stream from the server
func
NewEventSource
¶
func NewEventSource() *
EventSource
NewEventSource method creates a new instance of
EventSource
with default values for Server-Sent Events(SSE)
es := NewEventSource().
SetURL("https://sse.dev/test").
OnMessage(
func(e any) {
event := e.(*Event)
fmt.Println(event)
},
nil, // see method godoc
)

err := es.Connect()
fmt.Println(err)
See
EventSource.OnMessage
,
EventSource.AddEventListener
func (*EventSource)
AddEventListener
¶
func (es *
EventSource
) AddEventListener(eventName
string
, ef
EventMessageFunc
, result
any
) *
EventSource
AddEventListener method registers a callback to consume a specific event type
messages from the server. The second result argument is optional; it can be used
to register the data type for JSON data.
es.AddEventListener(
"friend_logged_in",
func(e any) {
event := e.(*Event)
fmt.Println(event)
},
nil,
)

// Receiving JSON data from the server, you can set result type
// to do auto-unmarshal
es.AddEventListener(
"friend_logged_in",
func(e any) {
event := e.(*UserLoggedIn)
fmt.Println(event)
},
UserLoggedIn{},
)
func (*EventSource)
AddHeader
¶
func (es *
EventSource
) AddHeader(header, value
string
) *
EventSource
AddHeader method adds a header and its value to the
EventSource
instance.
If the header key already exists, it appends. These headers will be sent in
the request while establishing a connection to the event source
es.AddHeader("Authorization", "token here").
AddHeader("X-Header", "value")
func (*EventSource)
Close
¶
func (es *
EventSource
) Close()
Close method used to close SSE connection explicitly
func (*EventSource)
Get
¶
func (es *
EventSource
) Get()
error
Get method establishes the connection with the server.
es := NewEventSource().
SetURL("https://sse.dev/test").
OnMessage(
func(e any) {
event := e.(*Event)
fmt.Println(event)
},
nil, // see method godoc
)

err := es.Get()
fmt.Println(err)
func (*EventSource)
Logger
¶
func (es *
EventSource
) Logger()
Logger
Logger method returns the logger instance used by the event source instance.
func (*EventSource)
OnError
¶
func (es *
EventSource
) OnError(ef
EventErrorFunc
) *
EventSource
OnError registered callback gets triggered when the error occurred
in the process
es.OnError(func(err error) {
fmt.Println("Error occurred:", err)
})
func (*EventSource)
OnMessage
¶
func (es *
EventSource
) OnMessage(ef
EventMessageFunc
, result
any
) *
EventSource
OnMessage method registers a callback to emit every SSE event message
from the server. The second result argument is optional; it can be used
to register the data type for JSON data.
es.OnMessage(
func(e any) {
event := e.(*Event)
fmt.Println("Event message", event)
},
nil,
)

// Receiving JSON data from the server, you can set result type
// to do auto-unmarshal
es.OnMessage(
func(e any) {
event := e.(*MyData)
fmt.Println(event)
},
MyData{},
)
func (*EventSource)
OnOpen
¶
func (es *
EventSource
) OnOpen(ef
EventOpenFunc
) *
EventSource
OnOpen registered callback gets triggered when the connection is
established with the server
es.OnOpen(func(url string) {
fmt.Println("I'm connected:", url)
})
func (*EventSource)
OnRequestFailure
¶
func (es *
EventSource
) OnRequestFailure(ef
EventRequestFailureFunc
) *
EventSource
OnRequestFailure registered callback gets triggered when the HTTP request
failure while establishing a SSE connection.
es.OnRequestFailure(func(err error, res *http.Response) {
fmt.Println("Error and response:", err, res)
})
Note:
Do not forget to close the HTTP response body.
HTTP response may be nil.
func (*EventSource)
SetBody
¶
func (es *
EventSource
) SetBody(body
io
.
Reader
) *
EventSource
SetBody method sets body value to the
EventSource
instance
Example:
es.SetBody(bytes.NewReader([]byte(`{"test":"put_data"}`)))
func (*EventSource)
SetHeader
¶
func (es *
EventSource
) SetHeader(header, value
string
) *
EventSource
SetHeader method sets a header and its value to the
EventSource
instance.
It overwrites the header value if the key already exists. These headers will be sent in
the request while establishing a connection to the event source
es.SetHeader("Authorization", "token here").
SetHeader("X-Header", "value")
func (*EventSource)
SetLogger
¶
func (es *
EventSource
) SetLogger(l
Logger
) *
EventSource
SetLogger method sets given writer for logging
Compliant to interface
resty.Logger
func (*EventSource)
SetMaxBufSize
¶
func (es *
EventSource
) SetMaxBufSize(bufSize
int
) *
EventSource
SetMaxBufSize method sets the given buffer size into the SSE client
Default is 32kb
es.SetMaxBufSize(64 * 1024) // 64kb
func (*EventSource)
SetMethod
¶
func (es *
EventSource
) SetMethod(method
string
) *
EventSource
SetMethod method sets a
EventSource
connection HTTP method in the instance
es.SetMethod("POST"), or es.SetMethod(resty.MethodPost)
func (*EventSource)
SetRetryCount
¶
func (es *
EventSource
) SetRetryCount(count
int
) *
EventSource
SetRetryCount method enables retry attempts on the SSE client while establishing
connection with the server
first attempt + retry count = total attempts
Default is 3
es.SetRetryCount(10)
func (*EventSource)
SetRetryMaxWaitTime
¶
func (es *
EventSource
) SetRetryMaxWaitTime(maxWaitTime
time
.
Duration
) *
EventSource
SetRetryMaxWaitTime method sets the max wait time for sleep before retrying
the request
Default is 2 seconds.
NOTE: The server-sent retry value takes precedence if present.
es.SetRetryMaxWaitTime(3 * time.Second)
func (*EventSource)
SetRetryWaitTime
¶
func (es *
EventSource
) SetRetryWaitTime(waitTime
time
.
Duration
) *
EventSource
SetRetryWaitTime method sets the default wait time for sleep before retrying
the request
Default is 100 milliseconds.
NOTE: The server-sent retry value takes precedence if present.
es.SetRetryWaitTime(1 * time.Second)
func (*EventSource)
SetTLSClientConfig
¶
func (es *
EventSource
) SetTLSClientConfig(tlsConfig *
tls
.
Config
) *
EventSource
SetTLSClientConfig method sets TLSClientConfig for underlying client Transport.
Values supported by
https://pkg.go.dev/crypto/tls#Config
can be configured.
// Disable SSL cert verification for local development
es.SetTLSClientConfig(&tls.Config{
InsecureSkipVerify: true
})
NOTE: This method overwrites existing
http.Transport.TLSClientConfig
func (*EventSource)
SetURL
¶
func (es *
EventSource
) SetURL(url
string
) *
EventSource
SetURL method sets a
EventSource
connection URL in the instance
es.SetURL("https://sse.dev/test")
func (*EventSource)
TLSClientConfig
¶
func (es *
EventSource
) TLSClientConfig() *
tls
.
Config
TLSClientConfig method returns the
tls.Config
from underlying client transport
otherwise returns nil
type
Host
¶
type Host struct {
// BaseURL represents the targeted host base URL
//
https://resty.dev
BaseURL
string
// Weight represents the host weight to determine
// the percentage of requests to send
Weight
int
// MaxFailures represents the value to mark the host as
// not usable until it reaches the Recovery duration
//	Default value is 5
MaxFailures
int
// contains filtered or unexported fields
}
Host struct used to represent the host information and its weight
to load balance the requests
type
HostState
¶
type HostState
int
const (
HostStateInActive
HostState
=
iota
HostStateActive
)
Host transition states
type
HostStateChangeFunc
¶
type HostStateChangeFunc func(baseURL
string
, from, to
HostState
)
HostStateChangeFunc type provides feedback on host state transitions
type
LoadBalancer
¶
type LoadBalancer interface {
Next() (
string
,
error
)
Feedback(*
RequestFeedback
)
Close()
error
}
LoadBalancer is the interface that wraps the HTTP client load-balancing
algorithm that returns the "Next" Base URL for the request to target
type
Logger
¶
type Logger interface {
Errorf(format
string
, v ...
any
)
Warnf(format
string
, v ...
any
)
Debugf(format
string
, v ...
any
)
}
Logger interface is to abstract the logging from Resty. Gives control to
the Resty users, choice of the logger.
type
MultipartField
¶
type MultipartField struct {
// Name of the multipart field name that the server expects it
Name
string
// FileName is used to set the file name we have to send to the server
FileName
string
// ContentType is a multipart file content-type value. It is highly
// recommended setting it if you know the content-type so that Resty
// don't have to do additional computing to auto-detect (Optional)
ContentType
string
// Reader is an input of [io.Reader] for multipart upload. It
// is optional if you set the FilePath value
Reader
io
.
Reader
// FilePath is a file path for multipart upload. It
// is optional if you set the Reader value
FilePath
string
// FileSize in bytes is used just for the information purpose of
// sharing via [MultipartFieldCallbackFunc] (Optional)
FileSize
int64
// ProgressCallback function is used to provide live progress details
// during a multipart upload (Optional)
//
// NOTE: It is recommended to set the FileSize value when using `MultipartField.Reader`
// with `ProgressCallback` feature so that Resty sends the FileSize
// value via [MultipartFieldProgress]
ProgressCallback
MultipartFieldCallbackFunc
// Values field is used to provide form field value. (Optional, unless it's a form-data field)
//
// It is primarily added for ordered multipart form-data field use cases
Values []
string
}
MultipartField struct represents the multipart field to compose
all
io.Reader
capable input for multipart form request
func (*MultipartField)
Clone
¶
func (mf *
MultipartField
) Clone() *
MultipartField
Clone method returns the deep copy of m except
io.Reader
.
type
MultipartFieldCallbackFunc
¶
type MultipartFieldCallbackFunc func(
MultipartFieldProgress
)
MultipartFieldCallbackFunc function used to transmit live multipart upload
progress in bytes count
type
MultipartFieldProgress
¶
type MultipartFieldProgress struct {
Name
string
FileName
string
FileSize
int64
Written
int64
}
MultipartFieldProgress struct used to provide multipart field upload progress
details via callback function
func (MultipartFieldProgress)
String
¶
func (mfp
MultipartFieldProgress
) String()
string
String method creates the string representation of
MultipartFieldProgress
type
RedirectInfo
¶
type RedirectInfo struct {
URL
string
StatusCode
int
}
RedirectInfo struct is used to capture the URL and status code for the redirect history
type
RedirectPolicy
¶
type RedirectPolicy interface {
Apply(*
http
.
Request
, []*
http
.
Request
)
error
}
RedirectPolicy to regulate the redirects in the Resty client.
Objects implementing the
RedirectPolicy
interface can be registered as
Apply function should return nil to continue the redirect journey; otherwise
return error to stop the redirect.
func
DomainCheckRedirectPolicy
¶
func DomainCheckRedirectPolicy(hostnames ...
string
)
RedirectPolicy
DomainCheckRedirectPolicy method is convenient for defining domain name redirect rules in Resty clients.
Redirect is allowed only for the host mentioned in the policy.
resty.SetRedirectPolicy(resty.DomainCheckRedirectPolicy("host1.com", "host2.org", "host3.net"))
func
FlexibleRedirectPolicy
¶
func FlexibleRedirectPolicy(noOfRedirect
int
)
RedirectPolicy
FlexibleRedirectPolicy method is convenient for creating several redirect policies for Resty clients.
resty.SetRedirectPolicy(FlexibleRedirectPolicy(20))
func
NoRedirectPolicy
¶
func NoRedirectPolicy()
RedirectPolicy
NoRedirectPolicy is used to disable the redirects in the Resty client
resty.SetRedirectPolicy(resty.NoRedirectPolicy())
type
RedirectPolicyFunc
¶
type RedirectPolicyFunc func(*
http
.
Request
, []*
http
.
Request
)
error
The
RedirectPolicyFunc
type is an adapter to allow the use of ordinary
functions as
RedirectPolicy
. If `f` is a function with the appropriate
signature, RedirectPolicyFunc(f) is a RedirectPolicy object that calls `f`.
func (RedirectPolicyFunc)
Apply
¶
func (f
RedirectPolicyFunc
) Apply(req *
http
.
Request
, via []*
http
.
Request
)
error
Apply calls f(req, via).
type
Request
¶
type Request struct {
URL
string
Method
string
AuthToken
string
AuthScheme
string
QueryParams
url
.
Values
FormData
url
.
Values
PathParams                 map[
string
]
string
Header
http
.
Header
Time
time
.
Time
Body
any
Result
any
Error
any
RawRequest                 *
http
.
Request
Cookies                    []*
http
.
Cookie
Debug
bool
CloseConnection
bool
DoNotParseResponse
bool
OutputFileName
string
ExpectResponseContentType
string
ForceResponseContentType
string
DebugBodyLimit
int
ResponseBodyLimit
int64
ResponseBodyUnlimitedReads
bool
IsTrace
bool
AllowMethodGetPayload
bool
AllowMethodDeletePayload
bool
IsDone
bool
IsSaveResponse
bool
Timeout
time
.
Duration
HeaderAuthorizationKey
string
RetryCount
int
RetryWaitTime
time
.
Duration
RetryMaxWaitTime
time
.
Duration
RetryStrategy
RetryStrategyFunc
IsRetryDefaultConditions
bool
AllowNonIdempotentRetry
bool
// RetryTraceID provides GUID for retry count > 0
RetryTraceID
string
// Attempt provides insights into no. of attempts
// Resty made.
//
//	first attempt + retry count = total attempts
Attempt
int
// contains filtered or unexported fields
}
Request struct is used to compose and fire individual requests from
Resty client. The
Request
provides an option to override client-level
settings and also an option for the request composition.
func (*Request)
AddRetryConditions
¶
func (r *
Request
) AddRetryConditions(conditions ...
RetryConditionFunc
) *
Request
AddRetryConditions method adds one or more retry condition functions into the request.
These retry conditions are executed to determine if the request can be retried.
The request will retry if any functions return `true`, otherwise return `false`.
NOTE:
The default retry conditions are applied first.
The client-level retry conditions are applied to all requests.
The request-level retry conditions are executed first before the client-level
retry conditions. See
Request.SetRetryConditions
func (*Request)
AddRetryHooks
¶
func (r *
Request
) AddRetryHooks(hooks ...
RetryHookFunc
) *
Request
AddRetryHooks method adds one or more side-effecting retry hooks in the request.
NOTE:
All the retry hooks are executed on each request retry.
The request-level retry hooks are executed first before the client-level hooks.
func (*Request)
Clone
¶
func (r *
Request
) Clone(ctx
context
.
Context
) *
Request
Clone returns a deep copy of r with its context changed to ctx.
It does clone appropriate fields, reset, and reinitialize, so
Request
can be used again.
The body is not copied, but it's a reference to the original body.
req := client.R().
SetBody("body").
SetHeader("header", "value")
clonedRequest := req.Clone(context.Background())
func (*Request)
Context
¶
func (r *
Request
) Context()
context
.
Context
Context method returns the request's
context.Context
. To change the context, use
Request.Clone
or
Request.WithContext
.
The returned context is always non-nil; it defaults to the
background context.
func (*Request)
CurlCmd
¶
func (r *
Request
) CurlCmd()
string
CurlCmd method generates the curl command for the request.
func (*Request)
Delete
¶
func (r *
Request
) Delete(url
string
) (*
Response
,
error
)
Delete method does DELETE HTTP request. It's defined in section 9.3.5 of
RFC 9110
.
func (*Request)
DisableDebug
¶
func (r *
Request
) DisableDebug() *
Request
DisableDebug method is a helper method for
Request.SetDebug
func (*Request)
DisableGenerateCurlCmd
¶
func (r *
Request
) DisableGenerateCurlCmd() *
Request
DisableGenerateCurlCmd method disables the option set by
Request.EnableGenerateCurlCmd
or
Request.SetGenerateCurlCmd
.
It overrides the options set in the
Client
.
func (*Request)
DisableRetryDefaultConditions
¶
func (r *
Request
) DisableRetryDefaultConditions() *
Request
DisableRetryDefaultConditions method disables the Resty's default retry
conditions on request level
func (*Request)
DisableTrace
¶
func (r *
Request
) DisableTrace() *
Request
DisableTrace method disables the request trace for the current request
func (*Request)
EnableDebug
¶
func (r *
Request
) EnableDebug() *
Request
EnableDebug method is a helper method for
Request.SetDebug
func (*Request)
EnableGenerateCurlCmd
¶
func (r *
Request
) EnableGenerateCurlCmd() *
Request
EnableGenerateCurlCmd method enables the generation of curl commands for the current request.
By default, Resty does not log the curl command in the debug log since it has the potential
to leak sensitive data unless explicitly enabled via
Request.SetDebugLogCurlCmd
or
Client.SetDebugLogCurlCmd
.
It overrides the options set in the
Client
.
NOTE: Use with care.
Potential to leak sensitive data from
Request
and
Response
in the debug log
when the debug log option is enabled.
Additional memory usage since the request body was reread.
curl body is not generated for
io.Reader
and multipart request flow.
func (*Request)
EnableRetryDefaultConditions
¶
func (r *
Request
) EnableRetryDefaultConditions() *
Request
EnableRetryDefaultConditions method enables the Resty's default retry
conditions on request level
func (*Request)
EnableTrace
¶
func (r *
Request
) EnableTrace() *
Request
EnableTrace method enables trace for the current request
using
httptrace.ClientTrace
and provides insights.
client := resty.New()

resp, err := client.R().EnableTrace().Get("https://httpbin.org/get")
fmt.Println("Error:", err)
fmt.Println("Trace Info:", resp.Request.TraceInfo())
See
Client.EnableTrace
,
Client.SetTrace
are also available to
get trace info for all requests.
func (*Request)
Execute
¶
func (r *
Request
) Execute(method, url
string
) (res *
Response
, err
error
)
Execute method performs the HTTP request with the given HTTP method and URL
for current
Request
.
resp, err := client.R().Execute(resty.MethodGet, "http://httpbin.org/get")
func (*Request)
Funcs
¶
func (r *
Request
) Funcs(funcs ...
RequestFunc
) *
Request
Funcs method gets executed on request composition that passes the
current request instance to provided
RequestFunc
, which could be
used to apply common/reusable logic to the given request instance.
func addRequestContentType(r *Request) *Request {
return r.SetHeader("Content-Type", "application/json").
SetHeader("Accept", "application/json")
}

func addRequestQueryParams(page, size int) func(r *Request) *Request {
return func(r *Request) *Request {
return r.SetQueryParam("page", strconv.Itoa(page)).
SetQueryParam("size", strconv.Itoa(size)).
SetQueryParam("request_no", strconv.Itoa(int(time.Now().Unix())))
}
}

client.R().
Funcs(addRequestContentType, addRequestQueryParams(1, 100)).
Get("https://localhost:8080/foobar")
func (*Request)
Get
¶
func (r *
Request
) Get(url
string
) (*
Response
,
error
)
Get method does GET HTTP request. It's defined in section 9.3.1 of
RFC 9110
.
func (*Request)
Head
¶
func (r *
Request
) Head(url
string
) (*
Response
,
error
)
Head method does HEAD HTTP request. It's defined in section 9.3.2 of
RFC 9110
.
func (*Request)
Options
¶
func (r *
Request
) Options(url
string
) (*
Response
,
error
)
Options method does OPTIONS HTTP request. It's defined in section 9.3.7 of
RFC 9110
.
func (*Request)
Patch
¶
func (r *
Request
) Patch(url
string
) (*
Response
,
error
)
Patch method does PATCH HTTP request. It's defined in section 2 of
RFC 5789
.
func (*Request)
Post
¶
func (r *
Request
) Post(url
string
) (*
Response
,
error
)
Post method does POST HTTP request. It's defined in section 9.3.3 of
RFC 9110
.
func (*Request)
Put
¶
func (r *
Request
) Put(url
string
) (*
Response
,
error
)
Put method does PUT HTTP request. It's defined in section 9.3.4 of
RFC 9110
.
func (*Request)
Send
¶
func (r *
Request
) Send() (*
Response
,
error
)
Send method performs the HTTP request using the method and URL already defined
for current
Request
.
res, err := client.R().
SetMethod(resty.MethodGet).
SetURL("http://httpbin.org/get").
Send()
func (*Request)
SetAllowMethodDeletePayload
¶
func (r *
Request
) SetAllowMethodDeletePayload(allow
bool
) *
Request
SetAllowMethodDeletePayload method allows the DELETE method with payload on the request level.
By default, Resty does not allow.
client.R().SetAllowMethodDeletePayload(true)
More info, refer to GH#881
It overrides the option set by the
Client.SetAllowMethodDeletePayload
func (*Request)
SetAllowMethodGetPayload
¶
func (r *
Request
) SetAllowMethodGetPayload(allow
bool
) *
Request
SetAllowMethodGetPayload method allows the GET method with payload on the request level.
By default, Resty does not allow.
client.R().SetAllowMethodGetPayload(true)
It overrides the option set by the
Client.SetAllowMethodGetPayload
func (*Request)
SetAllowNonIdempotentRetry
¶
func (r *
Request
) SetAllowNonIdempotentRetry(b
bool
) *
Request
SetAllowNonIdempotentRetry method is used to enable/disable non-idempotent HTTP
methods retry. By default, Resty only allows idempotent HTTP methods, see
RFC 9110 Section 9.2.2
,
RFC 9110 Section 18.2
It overrides value set at the client instance level, see
Client.SetAllowNonIdempotentRetry
func (*Request)
SetAuthScheme
¶
func (r *
Request
) SetAuthScheme(scheme
string
) *
Request
SetAuthScheme method sets the auth token scheme type in the HTTP request.
Example Header value structure:
Authorization: <auth-scheme-value-set-here> <auth-token-value>
For Example: To set the scheme to use OAuth
client.R().SetAuthScheme("OAuth")

// The outcome will be -
Authorization: OAuth <auth-token-value>
Information about Auth schemes can be found in
RFC 7235
, IANA
HTTP Auth schemes
It overrides the `Authorization` scheme set by method
Client.SetAuthScheme
.
func (*Request)
SetAuthToken
¶
func (r *
Request
) SetAuthToken(authToken
string
) *
Request
SetAuthToken method sets the auth token header(Default Scheme: Bearer) in the current HTTP request. Header example:
Authorization: Bearer <auth-token-value-comes-here>
For Example: To set auth token BC594900518B4F7EAC75BD37F019E08FBC594900518B4F7EAC75BD37F019E08F
client.R().SetAuthToken("BC594900518B4F7EAC75BD37F019E08FBC594900518B4F7EAC75BD37F019E08F")
It overrides the Auth token set by method
Client.SetAuthToken
.
func (*Request)
SetBasicAuth
¶
func (r *
Request
) SetBasicAuth(username, password
string
) *
Request
SetBasicAuth method sets the basic authentication header in the current HTTP request.
For Example:
Authorization: Basic <base64-encoded-value>
To set the header for username "go-resty" and password "welcome"
client.R().SetBasicAuth("go-resty", "welcome")
It overrides the credentials set by method
Client.SetBasicAuth
.
func (*Request)
SetBody
¶
func (r *
Request
) SetBody(body
any
) *
Request
SetBody method sets the request body for the request. It supports various practical needs as easy.
It's quite handy and powerful. Supported request body data types are `string`,
`[]byte`, `struct`, `map`, `slice` and
io.Reader
.
Body value can be pointer or non-pointer. Automatic marshalling for JSON and XML content type, if it is `struct`, `map`, or `slice`.
NOTE:
io.Reader
is processed in bufferless mode while sending a request.
For Example:
`struct` gets marshaled based on the request header `Content-Type`.
client.R().
SetBody(User{
Username: "jeeva@myjeeva.com",
Password: "welcome2resty",
})
'map` gets marshaled based on the request header `Content-Type`.
client.R().
SetBody(map[string]any{
"username": "jeeva@myjeeva.com",
"password": "welcome2resty",
"address": &Address{
Address1: "1111 This is my street",
Address2: "Apt 201",
City: "My City",
State: "My State",
ZipCode: 00000,
},
})
`string` as a body input. Suitable for any need as a string input.
client.R().
SetBody(`{
"username": "jeeva@getrightcare.com",
"password": "admin"
}`)
`[]byte` as a body input. Suitable for raw requests such as file upload, serialize & deserialize, etc.
client.R().
SetBody([]byte("This is my raw request, sent as-is"))
and so on.
func (*Request)
SetCloseConnection
¶
func (r *
Request
) SetCloseConnection(close
bool
) *
Request
SetCloseConnection method sets variable `Close` in HTTP request struct with the given
value. More info:
https://golang.org/src/net/http/request.go
It overrides the value set at the client instance level, see
Client.SetCloseConnection
func (*Request)
SetContentLength
¶
func (r *
Request
) SetContentLength(v
int64
) *
Request
SetContentLength method sets the given content length value in the HTTP request.
By default, Resty won't set `Content-Length`.
client.R().SetContentLength(3486547657)
func (*Request)
SetContentType
¶
func (r *
Request
) SetContentType(ct
string
) *
Request
SetContentType method is a convenient way to set the header Content-Type in the request
client.R().SetContentType("application/json")
func (*Request)
SetContext
¶
func (r *
Request
) SetContext(ctx
context
.
Context
) *
Request
SetContext method sets the
context.Context
for current
Request
.
It overwrites the current context in the Request instance; it does not
affect the
Request
.RawRequest that was already created.
If you want this method to take effect, use this method before invoking
Request.Send
or
Request
.HTTPVerb methods.
See
Request.WithContext
,
Request.Clone
func (*Request)
SetCookie
¶
func (r *
Request
) SetCookie(hc *
http
.
Cookie
) *
Request
SetCookie method appends a single cookie in the current request instance.
client.R().SetCookie(&http.Cookie{
Name:"go-resty",
Value:"This is cookie value",
})
NOTE: Method appends the Cookie value into existing Cookie even if its already existing.
func (*Request)
SetCookies
¶
func (r *
Request
) SetCookies(rs []*
http
.
Cookie
) *
Request
SetCookies method sets an array of cookies in the current request instance.
cookies := []*http.Cookie{
&http.Cookie{
Name:"go-resty-1",
Value:"This is cookie 1 value",
},
&http.Cookie{
Name:"go-resty-2",
Value:"This is cookie 2 value",
},
}

// Setting a cookies into resty's current request
client.R().SetCookies(cookies)
NOTE: Method appends the Cookie value into existing Cookie even if its already existing.
func (*Request)
SetDebug
¶
func (r *
Request
) SetDebug(d
bool
) *
Request
SetDebug method enables the debug mode on the current request. It logs
the details current request and response.
client.R().SetDebug(true)
// OR
client.R().EnableDebug()
Also, it can be enabled at the request level for a particular request; see
Request.SetDebug
.
For
Request
, it logs information such as HTTP verb, Relative URL path,
Host, Headers, and Body if it has one.
For
Response
, it logs information such as Status, Response Time, Headers,
and Body if it has one.
func (*Request)
SetDebugLogCurlCmd
¶
func (r *
Request
) SetDebugLogCurlCmd(b
bool
) *
Request
SetDebugLogCurlCmd method enables the curl command to be logged in the debug log
for the current request.
It can be overridden at the request level; see
Client.SetDebugLogCurlCmd
func (*Request)
SetDoNotParseResponse
¶
func (r *
Request
) SetDoNotParseResponse(notParse
bool
) *
Request
SetDoNotParseResponse method instructs Resty not to parse the response body automatically.
Resty exposes the raw response body as
io.ReadCloser
. If you use it, do not
forget to close the body, otherwise, you might get into connection leaks, and connection
reuse may not happen.
NOTE: The default
Response
middlewares are not executed when using this option. User
takes over the control of handling response body from Resty.
func (*Request)
SetError
¶
func (r *
Request
) SetError(err
any
) *
Request
SetError method is to register the request `Error` object for automatic unmarshalling for the request,
if the response status code is greater than 399 and the content type is either JSON or XML.
NOTE:
Request.SetError
input can be a pointer or non-pointer.
client.R().SetError(&AuthError{})
// OR
client.R().SetError(AuthError{})
Accessing an error value from response instance.
response.Error().(*AuthError)
If this request Error object is nil, Resty will use the client-level error object Type if it is set.
func (*Request)
SetExpectResponseContentType
¶
func (r *
Request
) SetExpectResponseContentType(contentType
string
) *
Request
SetExpectResponseContentType method allows to provide fallback `Content-Type`
for automatic unmarshalling when the `Content-Type` response header is unavailable.
func (*Request)
SetFile
¶
func (r *
Request
) SetFile(fieldName, filePath
string
) *
Request
SetFile method sets a single file field name and its path for multipart upload.
Resty provides an optional multipart live upload progress callback;
see method
Request.SetMultipartFields
client.R().
SetFile("my_file", "/Users/jeeva/Gas Bill - Sep.pdf")
func (*Request)
SetFileReader
¶
func (r *
Request
) SetFileReader(fieldName, fileName
string
, reader
io
.
Reader
) *
Request
SetFileReader method is to set a file using
io.Reader
for multipart upload.
Resty provides an optional multipart live upload progress callback;
see method
Request.SetMultipartFields
client.R().
SetFileReader("profile_img", "my-profile-img.png", bytes.NewReader(profileImgBytes)).
SetFileReader("notes", "user-notes.txt", bytes.NewReader(notesBytes))
func (*Request)
SetFiles
¶
func (r *
Request
) SetFiles(files map[
string
]
string
) *
Request
SetFiles method sets multiple file field names and their paths for multipart uploads.
Resty provides an optional multipart live upload progress callback;
see method
Request.SetMultipartFields
client.R().
SetFiles(map[string]string{
"my_file1": "/Users/jeeva/Gas Bill - Sep.pdf",
"my_file2": "/Users/jeeva/Electricity Bill - Sep.pdf",
"my_file3": "/Users/jeeva/Water Bill - Sep.pdf",
})
func (*Request)
SetForceResponseContentType
¶
func (r *
Request
) SetForceResponseContentType(contentType
string
) *
Request
SetForceResponseContentType method provides a strong sense of response `Content-Type` for
automatic unmarshalling. Resty gives this a higher priority than the `Content-Type`
response header.
This means that if both
Request.SetForceResponseContentType
is set and
the response `Content-Type` is available, `SetForceResponseContentType` value will win.
func (*Request)
SetFormData
¶
func (r *
Request
) SetFormData(data map[
string
]
string
) *
Request
SetFormData method sets form parameters and their values in the current request.
The request content type would be set as `application/x-www-form-urlencoded`.
client.R().
SetFormData(map[string]string{
"access_token": "BC594900-518B-4F7E-AC75-BD37F019E08F",
"user_id": "3455454545",
})
It overrides the form data value set at the client instance level.
See
Request.SetFormDataFromValues
for the same field name with multiple values.
func (*Request)
SetFormDataFromValues
¶
func (r *
Request
) SetFormDataFromValues(data
url
.
Values
) *
Request
SetFormDataFromValues method appends multiple form parameters with multi-value
(
url.Values
) at one go in the current request.
client.R().
SetFormDataFromValues(url.Values{
"search_criteria": []string{"book", "glass", "pencil"},
})
It overrides the form data value set at the client instance level.
func (*Request)
SetGenerateCurlCmd
¶
func (r *
Request
) SetGenerateCurlCmd(b
bool
) *
Request
SetGenerateCurlCmd method is used to turn on/off the generate curl command for the current request.
By default, Resty does not log the curl command in the debug log since it has the potential
to leak sensitive data unless explicitly enabled via
Request.SetDebugLogCurlCmd
or
Client.SetDebugLogCurlCmd
.
It overrides the options set by the
Client.SetGenerateCurlCmd
NOTE: Use with care.
Potential to leak sensitive data from
Request
and
Response
in the debug log
when the debug log option is enabled.
Additional memory usage since the request body was reread.
curl body is not generated for
io.Reader
and multipart request flow.
func (*Request)
SetHeader
¶
func (r *
Request
) SetHeader(header, value
string
) *
Request
SetHeader method sets a single header field and its value in the current request.
For Example: To set `Content-Type` and `Accept` as `application/json`.
client.R().
SetHeader("Content-Type", "application/json").
SetHeader("Accept", "application/json")
It overrides the header value set at the client instance level.
func (*Request)
SetHeaderAuthorizationKey
¶
func (r *
Request
) SetHeaderAuthorizationKey(k
string
) *
Request
SetHeaderAuthorizationKey method sets the given HTTP header name for Authorization in the request.
It overrides the `Authorization` header name set by method
Client.SetHeaderAuthorizationKey
.
client.R().SetHeaderAuthorizationKey("X-Custom-Authorization")
func (*Request)
SetHeaderMultiValues
¶
func (r *
Request
) SetHeaderMultiValues(headers map[
string
][]
string
) *
Request
SetHeaderMultiValues sets multiple header fields and their values as a list of strings in the current request.
For Example: To set `Accept` as `text/html, application/xhtml+xml, application/xml;q=0.9, image/webp, */*;q=0.8`
client.R().
SetHeaderMultiValues(map[string][]string{
"Accept": []string{"text/html", "application/xhtml+xml", "application/xml;q=0.9", "image/webp", "*/*;q=0.8"},
})
It overrides the header value set at the client instance level.
func (*Request)
SetHeaderVerbatim
¶
func (r *
Request
) SetHeaderVerbatim(header, value
string
) *
Request
SetHeaderVerbatim method is used to set the HTTP header key and value verbatim in the current request.
It is typically helpful for legacy applications or servers that require HTTP headers in a certain way
For Example: To set header key as `all_lowercase`, `UPPERCASE`, and `x-cloud-trace-id`
client.R().
SetHeaderVerbatim("all_lowercase", "available").
SetHeaderVerbatim("UPPERCASE", "available").
SetHeaderVerbatim("x-cloud-trace-id", "798e94019e5fc4d57fbb8901eb4c6cae")
It overrides the header value set at the client instance level.
func (*Request)
SetHeaders
¶
func (r *
Request
) SetHeaders(headers map[
string
]
string
) *
Request
SetHeaders method sets multiple header fields and their values at one go in the current request.
For Example: To set `Content-Type` and `Accept` as `application/json`
client.R().
SetHeaders(map[string]string{
"Content-Type": "application/json",
"Accept": "application/json",
})
It overrides the header value set at the client instance level.
func (*Request)
SetJSONEscapeHTML
¶
func (r *
Request
) SetJSONEscapeHTML(b
bool
) *
Request
SetJSONEscapeHTML method enables or disables the HTML escape on JSON marshal.
By default, escape HTML is `true`.
NOTE: This option only applies to the standard JSON Marshaller used by Resty.
It overrides the value set at the client instance level, see
Client.SetJSONEscapeHTML
func (*Request)
SetLogger
¶
func (r *
Request
) SetLogger(l
Logger
) *
Request
SetLogger method sets given writer for logging Resty request and response details.
By default, requests and responses inherit their logger from the client.
Compliant to interface
resty.Logger
.
It overrides the logger value set at the client instance level.
func (*Request)
SetMethod
¶
func (r *
Request
) SetMethod(m
string
) *
Request
SetMethod method used to set the HTTP verb for the request
func (*Request)
SetMultipartBoundary
¶
func (r *
Request
) SetMultipartBoundary(boundary
string
) *
Request
SetMultipartBoundary method sets the custom multipart boundary for the multipart request.
Typically, the `mime/multipart` package generates a random multipart boundary if not provided.
func (*Request)
SetMultipartField
¶
func (r *
Request
) SetMultipartField(fieldName, fileName, contentType
string
, reader
io
.
Reader
) *
Request
SetMultipartField method sets custom data with Content-Type using
io.Reader
for multipart upload.
Resty provides an optional multipart live upload progress callback;
see method
Request.SetMultipartFields
func (*Request)
SetMultipartFields
¶
func (r *
Request
) SetMultipartFields(fields ...*
MultipartField
) *
Request
SetMultipartFields method sets multiple data fields using
io.Reader
for multipart upload.
Resty provides an optional multipart live upload progress count in bytes; see
MultipartField
.ProgressCallback and
MultipartFieldProgress
For Example:
client.R().SetMultipartFields(
&resty.MultipartField{
Name:        "uploadManifest1",
FileName:    "upload-file-1.json",
ContentType: "application/json",
Reader:      strings.NewReader(`{"input": {"name": "Uploaded document 1", "_filename" : ["file1.txt"]}}`),
},
&resty.MultipartField{
Name:        "uploadManifest2",
FileName:    "upload-file-2.json",
ContentType: "application/json",
FilePath:    "/path/to/upload-file-2.json",
},
&resty.MultipartField{
Name:             "image-file1",
FileName:         "image-file1.png",
ContentType:      "image/png",
Reader:           bytes.NewReader(fileBytes),
ProgressCallback: func(mp MultipartFieldProgress) {
// use the progress details
},
},
&resty.MultipartField{
Name:             "image-file2",
FileName:         "image-file2.png",
ContentType:      "image/png",
Reader:           imageFile2, // instance of *os.File
ProgressCallback: func(mp MultipartFieldProgress) {
// use the progress details
},
})
If you have a `slice` of fields already, then call-
client.R().SetMultipartFields(fields...)
func (*Request)
SetMultipartFormData
¶
func (r *
Request
) SetMultipartFormData(data map[
string
]
string
) *
Request
SetMultipartFormData method allows simple form data to be attached to the request
as `multipart:form-data`
func (*Request)
SetMultipartOrderedFormData
¶
func (r *
Request
) SetMultipartOrderedFormData(name
string
, values []
string
) *
Request
SetMultipartOrderedFormData method allows add ordered form data to be attached to the request
as `multipart:form-data`
func (*Request)
SetOutputFileName
¶
func (r *
Request
) SetOutputFileName(file
string
) *
Request
SetOutputFileName method sets the output file for the current HTTP request. The current
HTTP response will be saved in the given file. It is similar to the `curl -o` flag.
Absolute path or relative path can be used.
If it is a relative path, then the output file goes under the output directory, as mentioned
in the
Client.SetOutputDirectory
.
client.R().
SetOutputFileName("/Users/jeeva/Downloads/ReplyWithHeader-v5.1-beta.zip").
Get("http://bit.ly/1LouEKr")
NOTE: In this scenario
[Response.BodyBytes] might be nil.
Response
.Body might have been already read.
func (*Request)
SetPathParam
¶
func (r *
Request
) SetPathParam(param, value
string
) *
Request
SetPathParam method sets a single URL path key-value pair in the
Resty current request instance.
client.R().SetPathParam("userId", "sample@sample.com")

Result:
URL - /v1/users/{userId}/details
Composed URL - /v1/users/sample@sample.com/details

client.R().SetPathParam("path", "groups/developers")

Result:
URL - /v1/users/{path}/details
Composed URL - /v1/users/groups%2Fdevelopers/details
It replaces the value of the key while composing the request URL.
The values will be escaped using function
url.PathEscape
.
It overrides the path parameter set at the client instance level.
func (*Request)
SetPathParams
¶
func (r *
Request
) SetPathParams(params map[
string
]
string
) *
Request
SetPathParams method sets multiple URL path key-value pairs at one go in the
Resty current request instance.
client.R().SetPathParams(map[string]string{
"userId":       "sample@sample.com",
"subAccountId": "100002",
"path":         "groups/developers",
})

Result:
URL - /v1/users/{userId}/{subAccountId}/{path}/details
Composed URL - /v1/users/sample@sample.com/100002/groups%2Fdevelopers/details
It replaces the value of the key while composing the request URL.
The values will be escaped using function
url.PathEscape
.
It overrides the path parameter set at the client instance level.
func (*Request)
SetQueryParam
¶
func (r *
Request
) SetQueryParam(param, value
string
) *
Request
SetQueryParam method sets a single parameter and its value in the current request.
It will be formed as a query string for the request.
For Example: `search=kitchen%20papers&size=large` in the URL after the `?` mark.
client.R().
SetQueryParam("search", "kitchen papers").
SetQueryParam("size", "large")
It overrides the query parameter value set at the client instance level.
func (*Request)
SetQueryParams
¶
func (r *
Request
) SetQueryParams(params map[
string
]
string
) *
Request
SetQueryParams method sets multiple parameters and their values at one go in the current request.
It will be formed as a query string for the request.
For Example: `search=kitchen%20papers&size=large` in the URL after the `?` mark.
client.R().
SetQueryParams(map[string]string{
"search": "kitchen papers",
"size": "large",
})
It overrides the query parameter value set at the client instance level.
func (*Request)
SetQueryParamsFromValues
¶
func (r *
Request
) SetQueryParamsFromValues(params
url
.
Values
) *
Request
SetQueryParamsFromValues method appends multiple parameters with multi-value
(
url.Values
) at one go in the current request. It will be formed as
query string for the request.
For Example: `status=pending&status=approved&status=open` in the URL after the `?` mark.
client.R().
SetQueryParamsFromValues(url.Values{
"status": []string{"pending", "approved", "open"},
})
It overrides the query parameter value set at the client instance level.
func (*Request)
SetQueryString
¶
func (r *
Request
) SetQueryString(query
string
) *
Request
SetQueryString method provides the ability to use string as an input to set URL query string for the request.
client.R().
SetQueryString("productId=232&template=fresh-sample&cat=resty&source=google&kw=buy a lot more")
It overrides the query parameter value set at the client instance level.
func (*Request)
SetRawPathParam
¶
func (r *
Request
) SetRawPathParam(param, value
string
) *
Request
SetRawPathParam method sets a single URL path key-value pair in the
Resty current request instance without path escape.
client.R().SetRawPathParam("userId", "sample@sample.com")

Result:
URL - /v1/users/{userId}/details
Composed URL - /v1/users/sample@sample.com/details

client.R().SetRawPathParam("path", "groups/developers")

Result:
URL - /v1/users/{path}/details
Composed URL - /v1/users/groups/developers/details
It replaces the value of the key while composing the request URL.
The value will be used as-is, no path escape applied.
It overrides the raw path parameter set at the client instance level.
func (*Request)
SetRawPathParams
¶
func (r *
Request
) SetRawPathParams(params map[
string
]
string
) *
Request
SetRawPathParams method sets multiple URL path key-value pairs at one go in the
Resty current request instance without path escape.
client.R().SetPathParams(map[string]string{
"userId": "sample@sample.com",
"subAccountId": "100002",
"path":         "groups/developers",
})

Result:
URL - /v1/users/{userId}/{subAccountId}/{path}/details
Composed URL - /v1/users/sample@sample.com/100002/groups/developers/details
It replaces the value of the key while composing the request URL.
The value will be used as-is, no path escape applied.
It overrides the raw path parameter set at the client instance level.
func (*Request)
SetResponseBodyLimit
¶
func (r *
Request
) SetResponseBodyLimit(v
int64
) *
Request
SetResponseBodyLimit method sets a maximum body size limit in bytes on response,
avoid reading too much data to memory.
Client will return
resty.ErrResponseBodyTooLarge
if the body size of the body
in the uncompressed response is larger than the limit.
Body size limit will not be enforced in the following cases:
ResponseBodyLimit <= 0, which is the default behavior.
Request.SetOutputFileName
is called to save response data to the file.
"DoNotParseResponse" is set for client or request.
It overrides the value set at the client instance level, see
Client.SetResponseBodyLimit
func (*Request)
SetResponseBodyUnlimitedReads
¶
func (r *
Request
) SetResponseBodyUnlimitedReads(b
bool
) *
Request
SetResponseBodyUnlimitedReads method is to turn on/off the response body in memory
that provides an ability to do unlimited reads.
It overrides the value set at the client level; see
Client.SetResponseBodyUnlimitedReads
Unlimited reads are possible in a few scenarios, even without enabling it.
When debug mode is enabled
NOTE: Use with care
Turning on this feature keeps the response body in memory, which might cause additional memory usage.
func (*Request)
SetResult
¶
func (r *
Request
) SetResult(v
any
) *
Request
SetResult method is to register the response `Result` object for automatic
unmarshalling of the HTTP response if the response status code is
between 200 and 299, and the content type is JSON or XML.
Note:
Request.SetResult
input can be a pointer or non-pointer.
The pointer with handle
authToken := &AuthToken{}
client.R().SetResult(authToken)

// Can be accessed via -
fmt.Println(authToken) OR fmt.Println(response.Result().(*AuthToken))
OR -
The pointer without handle or non-pointer
client.R().SetResult(&AuthToken{})
// OR
client.R().SetResult(AuthToken{})

// Can be accessed via -
fmt.Println(response.Result().(*AuthToken))
func (*Request)
SetRetryConditions
¶
func (r *
Request
) SetRetryConditions(conditions ...
RetryConditionFunc
) *
Request
SetRetryConditions method overwrites the retry conditions in the request.
These retry conditions are executed to determine if the request can be retried.
The request will retry if any function returns `true`, otherwise return `false`.
func (*Request)
SetRetryCount
¶
func (r *
Request
) SetRetryCount(count
int
) *
Request
SetRetryCount method enables retry on Resty client and allows you
to set no. of retry count.
first attempt + retry count = total attempts
See
Request.SetRetryStrategy
NOTE:
By default, Resty only does retry on idempotent HTTP verb,
RFC 9110 Section 9.2.2
,
RFC 9110 Section 18.2
func (*Request)
SetRetryDefaultConditions
¶
func (r *
Request
) SetRetryDefaultConditions(b
bool
) *
Request
SetRetryDefaultConditions method is used to enable/disable the Resty's default
retry conditions on request level
It overrides value set at the client instance level, see
Client.SetRetryDefaultConditions
func (*Request)
SetRetryHooks
¶
func (r *
Request
) SetRetryHooks(hooks ...
RetryHookFunc
) *
Request
SetRetryHooks method overwrites side-effecting retry hooks in the request.
NOTE:
All the retry hooks are executed on each request retry.
func (*Request)
SetRetryMaxWaitTime
¶
func (r *
Request
) SetRetryMaxWaitTime(maxWaitTime
time
.
Duration
) *
Request
SetRetryMaxWaitTime method sets the max wait time for sleep before retrying
Default is 2 seconds.
func (*Request)
SetRetryStrategy
¶
func (r *
Request
) SetRetryStrategy(rs
RetryStrategyFunc
) *
Request
SetRetryStrategy method used to set the custom Retry strategy on request,
it is used to get wait time before each retry. It overrides the retry
strategy set at the client instance level, see
Client.SetRetryStrategy
Default (nil) implies capped exponential backoff with a jitter strategy
func (*Request)
SetRetryWaitTime
¶
func (r *
Request
) SetRetryWaitTime(waitTime
time
.
Duration
) *
Request
SetRetryWaitTime method sets the default wait time for sleep before retrying
Default is 100 milliseconds.
func (*Request)
SetSaveResponse
¶
func (r *
Request
) SetSaveResponse(save
bool
) *
Request
SetSaveResponse method used to enable the save response option for the current requests
client.R().SetSaveResponse(true)
Resty determines the save filename in the following order -
Request.SetOutputFileName
Content-Disposition header
Request URL using
path.Base
Request URL hostname if path is empty or "/"
It overrides the value set at the client instance level, see
Client.SetSaveResponse
func (*Request)
SetTimeout
¶
func (r *
Request
) SetTimeout(timeout
time
.
Duration
) *
Request
SetTimeout method is used to set a timeout for the current request
client.R().SetTimeout(1 * time.Minute)
It overrides the timeout set at the client instance level, See
Client.SetTimeout
NOTE: Resty uses
context.WithTimeout
on the request, it does not use
http.Client.Timeout
func (*Request)
SetTrace
¶
func (r *
Request
) SetTrace(t
bool
) *
Request
SetTrace method is used to turn on/off the trace capability at the request level
See
Request.EnableTrace
or
Client.SetTrace
func (*Request)
SetURL
¶
func (r *
Request
) SetURL(url
string
) *
Request
SetURL method used to set the request URL for the request
func (*Request)
SetUnescapeQueryParams
¶
func (r *
Request
) SetUnescapeQueryParams(unescape
bool
) *
Request
SetUnescapeQueryParams method sets the choice of unescape query parameters for the request URL.
To prevent broken URL, Resty replaces space (" ") with "+" in the query parameters.
This method overrides the value set by
Client.SetUnescapeQueryParams
NOTE: Request failure is possible due to non-standard usage of Unescaped Query Parameters.
func (*Request)
Trace
¶
func (r *
Request
) Trace(url
string
) (*
Response
,
error
)
Trace method does TRACE HTTP request. It's defined in section 9.3.8 of
RFC 9110
.
func (*Request)
TraceInfo
¶
func (r *
Request
) TraceInfo()
TraceInfo
TraceInfo method returns the trace info for the request.
If either the
Client.EnableTrace
or
Request.EnableTrace
function has not been called
before the request is made, an empty
resty.TraceInfo
object is returned.
func (*Request)
WithContext
¶
func (r *
Request
) WithContext(ctx
context
.
Context
) *
Request
WithContext method returns a shallow copy of r with its context changed
to ctx. The provided ctx must be non-nil. It does not
affect the
Request
.RawRequest that was already created.
If you want this method to take effect, use this method before invoking
Request.Send
or
Request
.HTTPVerb methods.
See
Request.SetContext
,
Request.Clone
type
RequestFeedback
¶
type RequestFeedback struct {
BaseURL
string
Success
bool
Attempt
int
}
RequestFeedback struct is used to send the request feedback to load balancing
algorithm
type
RequestFunc
¶
type RequestFunc func(*
Request
) *
Request
RequestFunc type is for extended manipulation of the Request instance
type
RequestMiddleware
¶
type RequestMiddleware func(*
Client
, *
Request
)
error
RequestMiddleware type is for request middleware, called before a request is sent
type
Response
¶
type Response struct {
Request     *
Request
Body
io
.
ReadCloser
RawResponse *
http
.
Response
IsRead
bool
// Err field used to cascade the response middleware error
// in the chain
Err
error
// contains filtered or unexported fields
}
Response struct holds response values of executed requests.
func (*Response)
Bytes
¶
func (r *
Response
) Bytes() []
byte
Bytes method returns the body of the HTTP response as a byte slice.
It returns an empty byte slice if it is nil or the body is zero length.
NOTE:
Returns an empty byte slice on auto-unmarshal scenarios, unless
Client.SetResponseBodyUnlimitedReads
or
Request.SetResponseBodyUnlimitedReads
set.
Returns an empty byte slice when
Client.SetDoNotParseResponse
or
Request.SetDoNotParseResponse
set.
func (*Response)
Cookies
¶
func (r *
Response
) Cookies() []*
http
.
Cookie
Cookies method to returns all the response cookies
func (*Response)
Duration
¶
func (r *
Response
) Duration()
time
.
Duration
Duration method returns the duration of HTTP response time from the request we sent
and received a request.
See
Response.ReceivedAt
to know when the client received a response and see
`Response.Request.Time` to know when the client sent a request.
func (*Response)
Error
¶
func (r *
Response
) Error()
any
Error method returns the error object if it has one
See
Request.SetError
,
Client.SetError
func (*Response)
Header
¶
func (r *
Response
) Header()
http
.
Header
Header method returns the response headers
func (*Response)
IsError
¶
func (r *
Response
) IsError()
bool
IsError method returns true if HTTP status `code >= 400` otherwise false.
func (*Response)
IsSuccess
¶
func (r *
Response
) IsSuccess()
bool
IsSuccess method returns true if HTTP status `code >= 200 and <= 299` otherwise false.
func (*Response)
Proto
¶
func (r *
Response
) Proto()
string
Proto method returns the HTTP response protocol used for the request.
func (*Response)
ReceivedAt
¶
func (r *
Response
) ReceivedAt()
time
.
Time
ReceivedAt method returns the time we received a response from the server for the request.
func (*Response)
RedirectHistory
¶
func (r *
Response
) RedirectHistory() []*
RedirectInfo
RedirectHistory method returns a redirect history slice with the URL and status code
func (*Response)
Result
¶
func (r *
Response
) Result()
any
Result method returns the response value as an object if it has one
See
Request.SetResult
func (*Response)
Size
¶
func (r *
Response
) Size()
int64
Size method returns the HTTP response size in bytes. Yeah, you can rely on HTTP `Content-Length`
header, however it won't be available for chucked transfer/compressed response.
Since Resty captures response size details when processing the response body
when possible. So that users get the actual size of response bytes.
func (*Response)
Status
¶
func (r *
Response
) Status()
string
Status method returns the HTTP status string for the executed request.
Example: 200 OK
func (*Response)
StatusCode
¶
func (r *
Response
) StatusCode()
int
StatusCode method returns the HTTP status code for the executed request.
Example: 200
func (*Response)
String
¶
func (r *
Response
) String()
string
String method returns the body of the HTTP response as a `string`.
It returns an empty string if it is nil or the body is zero length.
NOTE:
Returns an empty string on auto-unmarshal scenarios, unless
Client.SetResponseBodyUnlimitedReads
or
Request.SetResponseBodyUnlimitedReads
set.
Returns an empty string when
Client.SetDoNotParseResponse
or
Request.SetDoNotParseResponse
set.
type
ResponseError
¶
type ResponseError struct {
Response *
Response
Err
error
}
ResponseError is a wrapper that includes the server response with an error.
Neither the err nor the response should be nil.
func (*ResponseError)
Error
¶
func (e *
ResponseError
) Error()
string
func (*ResponseError)
Unwrap
¶
func (e *
ResponseError
) Unwrap()
error
type
ResponseMiddleware
¶
type ResponseMiddleware func(*
Client
, *
Response
)
error
ResponseMiddleware type is for response middleware, called after a response has been received
type
RetryConditionFunc
¶
type RetryConditionFunc func(*
Response
,
error
)
bool
RetryConditionFunc type is for the retry condition function
input: non-nil Response OR request execution error
type
RetryHookFunc
¶
type RetryHookFunc func(*
Response
,
error
)
RetryHookFunc is for side-effecting functions triggered on retry
type
RetryStrategyFunc
¶
type RetryStrategyFunc func(*
Response
,
error
) (
time
.
Duration
,
error
)
RetryStrategyFunc type is for custom retry strategy implementation
By default Resty uses the capped exponential backoff with a jitter strategy
type
RoundRobin
¶
type RoundRobin struct {
// contains filtered or unexported fields
}
RoundRobin struct used to implement the Round-Robin(RR) request
load balancer algorithm
func
NewRoundRobin
¶
func NewRoundRobin(baseURLs ...
string
) (*
RoundRobin
,
error
)
NewRoundRobin method creates the new Round-Robin(RR) request load balancer
instance with given base URLs
func (*RoundRobin)
Close
¶
func (rr *
RoundRobin
) Close()
error
Close method does nothing in Round-Robin(RR) request load balancer
func (*RoundRobin)
Feedback
¶
func (rr *
RoundRobin
) Feedback(_ *
RequestFeedback
)
Feedback method does nothing in Round-Robin(RR) request load balancer
func (*RoundRobin)
Next
¶
func (rr *
RoundRobin
) Next() (
string
,
error
)
Next method returns the next Base URL based on the Round-Robin(RR) algorithm
func (*RoundRobin)
Refresh
¶
func (rr *
RoundRobin
) Refresh(baseURLs ...
string
)
error
Refresh method reset the existing Base URLs with the given Base URLs slice to refresh it
type
SRVWeightedRoundRobin
¶
type SRVWeightedRoundRobin struct {
Service
string
Proto
string
DomainName
string
HttpScheme
string
// contains filtered or unexported fields
}
SRVWeightedRoundRobin struct used to implement SRV Weighted Round-Robin(RR) algorithm
func
NewSRVWeightedRoundRobin
¶
func NewSRVWeightedRoundRobin(service, proto, domainName, httpScheme
string
) (*
SRVWeightedRoundRobin
,
error
)
NewSRVWeightedRoundRobin method creates a new Weighted Round-Robin(WRR) load balancer instance
with given SRV values
func (*SRVWeightedRoundRobin)
Close
¶
func (swrr *
SRVWeightedRoundRobin
) Close()
error
Close method does the cleanup by stopping the
time.Ticker
SRV Base URL based
on Weighted Round-Robin(WRR) request load balancer
func (*SRVWeightedRoundRobin)
Feedback
¶
func (swrr *
SRVWeightedRoundRobin
) Feedback(f *
RequestFeedback
)
Feedback method does nothing in SRV Base URL based on Weighted Round-Robin(WRR)
request load balancer
func (*SRVWeightedRoundRobin)
Next
¶
func (swrr *
SRVWeightedRoundRobin
) Next() (
string
,
error
)
Next method returns the next SRV Base URL based on Weighted Round-Robin(RR)
func (*SRVWeightedRoundRobin)
Refresh
¶
func (swrr *
SRVWeightedRoundRobin
) Refresh()
error
Refresh method reset the values based
net.LookupSRV
values to refresh it
func (*SRVWeightedRoundRobin)
SetOnStateChange
¶
func (swrr *
SRVWeightedRoundRobin
) SetOnStateChange(fn
HostStateChangeFunc
)
SetOnStateChange method used to set a callback for the host transition state
func (*SRVWeightedRoundRobin)
SetRecoveryDuration
¶
func (swrr *
SRVWeightedRoundRobin
) SetRecoveryDuration(d
time
.
Duration
)
SetRecoveryDuration method is used to change the existing recovery duration for the host
func (*SRVWeightedRoundRobin)
SetRefreshDuration
¶
func (swrr *
SRVWeightedRoundRobin
) SetRefreshDuration(d
time
.
Duration
)
SetRefreshDuration method assists in changing the default (180 seconds) refresh duration
type
SuccessHook
¶
type SuccessHook func(*
Client
, *
Response
)
SuccessHook type is for reacting to request success
type
TLSClientConfiger
¶
type TLSClientConfiger interface {
TLSClientConfig() *
tls
.
Config
SetTLSClientConfig(*
tls
.
Config
)
error
}
TLSClientConfiger interface is to configure TLS Client configuration on custom transport
implemented using
http.RoundTripper
type
TraceInfo
¶
type TraceInfo struct {
// DNSLookup is the duration that transport took to perform
// DNS lookup.
DNSLookup
time
.
Duration
`json:"dns_lookup_time"`
// ConnTime is the duration it took to obtain a successful connection.
ConnTime
time
.
Duration
`json:"connection_time"`
// TCPConnTime is the duration it took to obtain the TCP connection.
TCPConnTime
time
.
Duration
`json:"tcp_connection_time"`
// TLSHandshake is the duration of the TLS handshake.
TLSHandshake
time
.
Duration
`json:"tls_handshake_time"`
// ServerTime is the server's duration for responding to the first byte.
ServerTime
time
.
Duration
`json:"server_time"`
// ResponseTime is the duration since the first response byte from the server to
// request completion.
ResponseTime
time
.
Duration
`json:"response_time"`
// TotalTime is the duration of the total time request taken end-to-end.
TotalTime
time
.
Duration
`json:"total_time"`
// IsConnReused is whether this connection has been previously
// used for another HTTP request.
IsConnReused
bool
`json:"is_connection_reused"`
// IsConnWasIdle is whether this connection was obtained from an
// idle pool.
IsConnWasIdle
bool
`json:"is_connection_was_idle"`
// ConnIdleTime is the duration how long the connection that was previously
// idle, if IsConnWasIdle is true.
ConnIdleTime
time
.
Duration
`json:"connection_idle_time"`
// RequestAttempt is to represent the request attempt made during a Resty
// request execution flow, including retry count.
RequestAttempt
int
`json:"request_attempt"`
// RemoteAddr returns the remote network address.
RemoteAddr
string
`json:"remote_address"`
}
TraceInfo struct is used to provide request trace info such as DNS lookup
duration, Connection obtain duration, Server processing duration, etc.
func (TraceInfo)
Clone
¶
func (ti
TraceInfo
) Clone() *
TraceInfo
Clone method returns the clone copy of
TraceInfo
func (TraceInfo)
JSON
¶
func (ti
TraceInfo
) JSON()
string
JSON method returns the JSON string of request trace information
func (TraceInfo)
String
¶
func (ti
TraceInfo
) String()
string
String method returns string representation of request trace information.
type
TransportSettings
¶
type TransportSettings struct {
// DialerTimeout, default value is `30` seconds.
DialerTimeout
time
.
Duration
// DialerKeepAlive, default value is `30` seconds.
DialerKeepAlive
time
.
Duration
// IdleConnTimeout, default value is `90` seconds.
IdleConnTimeout
time
.
Duration
// TLSHandshakeTimeout, default value is `10` seconds.
TLSHandshakeTimeout
time
.
Duration
// ExpectContinueTimeout, default value is `1` seconds.
ExpectContinueTimeout
time
.
Duration
// ResponseHeaderTimeout, added to provide ability to
// set value. No default value in Resty, the Go
// HTTP client default value applies.
ResponseHeaderTimeout
time
.
Duration
// MaxIdleConns, default value is `100`.
MaxIdleConns
int
// MaxIdleConnsPerHost, default value is `runtime.GOMAXPROCS(0) + 1`.
MaxIdleConnsPerHost
int
// MaxConnsPerHost, default value is no limit.
MaxConnsPerHost
int
// DisableKeepAlives, default value is `false`.
DisableKeepAlives
bool
// MaxResponseHeaderBytes, added to provide ability to
// set value. No default value in Resty, the Go
// HTTP client default value applies.
MaxResponseHeaderBytes
int64
// WriteBufferSize, added to provide ability to
// set value. No default value in Resty, the Go
// HTTP client default value applies.
WriteBufferSize
int
// ReadBufferSize, added to provide ability to
// set value. No default value in Resty, the Go
// HTTP client default value applies.
ReadBufferSize
int
}
TransportSettings struct is used to define custom dialer and transport
values for the Resty client. Please refer to individual
struct fields to know the default values.
Also, refer to
https://pkg.go.dev/net/http#Transport
for more details.
type
WeightedRoundRobin
¶
type WeightedRoundRobin struct {
// contains filtered or unexported fields
}
WeightedRoundRobin struct used to represent the host details for
Weighted Round-Robin(WRR) algorithm implementation
func
NewWeightedRoundRobin
¶
func NewWeightedRoundRobin(recovery
time
.
Duration
, hosts ...*
Host
) (*
WeightedRoundRobin
,
error
)
NewWeightedRoundRobin method creates the new Weighted Round-Robin(WRR)
request load balancer instance with given recovery duration and hosts slice
func (*WeightedRoundRobin)
Close
¶
func (wrr *
WeightedRoundRobin
) Close()
error
Close method does the cleanup by stopping the
time.Ticker
on
Weighted Round-Robin(WRR) request load balancer
func (*WeightedRoundRobin)
Feedback
¶
func (wrr *
WeightedRoundRobin
) Feedback(f *
RequestFeedback
)
Feedback method process the request feedback for Weighted Round-Robin(WRR)
request load balancer
func (*WeightedRoundRobin)
Next
¶
func (wrr *
WeightedRoundRobin
) Next() (
string
,
error
)
Next method returns the next Base URL based on Weighted Round-Robin(WRR)
func (*WeightedRoundRobin)
Refresh
¶
func (wrr *
WeightedRoundRobin
) Refresh(hosts ...*
Host
)
error
Refresh method reset the existing values with the given
Host
slice to refresh it
func (*WeightedRoundRobin)
SetOnStateChange
¶
func (wrr *
WeightedRoundRobin
) SetOnStateChange(fn
HostStateChangeFunc
)
SetOnStateChange method used to set a callback for the host transition state
func (*WeightedRoundRobin)
SetRecoveryDuration
¶
func (wrr *
WeightedRoundRobin
) SetRecoveryDuration(d
time
.
Duration
)
SetRecoveryDuration method is used to change the existing recovery duration for the host