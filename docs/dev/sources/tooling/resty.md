# resty HTTP Client

> Source: https://pkg.go.dev/github.com/go-resty/resty/v2
> Fetched: 2026-02-01T11:42:20.828103+00:00
> Content-Hash: ead4b89ad7971197
> Type: html

---

### Overview ¶

Package resty provides Simple HTTP and REST client library for Go.

Example (ClientCertificates) ¶

    // Parsing public/private key pair from a pair of files. The files must contain PEM encoded data.
    cert, err := tls.LoadX509KeyPair("certs/client.pem", "certs/client.key")
    if err != nil {
     log.Fatalf("ERROR client certificate: %s", err)
    }
    
    // Create a resty client
    client := resty.New()
    
    client.SetCertificates(cert)
    

Example (CustomRootCertificate) ¶

    // Create a resty client
    client := resty.New()
    client.SetRootCertificate("/path/to/root/pemFile.pem")
    

Example (DropboxUpload) ¶

    // For example: upload file to Dropbox
    // POST of raw bytes for file upload.
    fileBytes, _ := os.ReadFile("/Users/jeeva/mydocument.pdf")
    
    // Create a resty client
    client := resty.New()
    
    // See we are not setting content-type header, since go-resty automatically detects Content-Type for you
    resp, err := client.R().
     SetBody(fileBytes).     // resty autodetects content type
     SetContentLength(true). // Dropbox expects this value
     SetAuthToken("<your-auth-token>").
     SetError(DropboxError{}).
     Post("https://content.dropboxapi.com/1/files_put/auto/resty/mydocument.pdf") // you can use PUT method too dropbox supports it
    
    // Output print
    fmt.Printf("\nError: %v\n", err)
    fmt.Printf("Time: %v\n", resp.Time())
    fmt.Printf("Body: %v\n", resp)
    

Example (EnhancedGet) ¶

    // Create a resty client
    client := resty.New()
    
    resp, err := client.R().
     SetQueryParams(map[string]string{
      "page_no": "1",
      "limit":   "20",
      "sort":    "name",
      "order":   "asc",
      "random":  strconv.FormatInt(time.Now().Unix(), 10),
     }).
     SetHeader("Accept", "application/json").
     SetAuthToken("BC594900518B4F7EAC75BD37F019E08FBC594900518B4F7EAC75BD37F019E08F").
     Get("/search_result")
    
    printOutput(resp, err)
    

Example (Get) ¶

    // Create a resty client
    client := resty.New()
    
    resp, err := client.R().Get("http://httpbin.org/get")
    
    fmt.Printf("\nError: %v", err)
    fmt.Printf("\nResponse Status Code: %v", resp.StatusCode())
    fmt.Printf("\nResponse Status: %v", resp.Status())
    fmt.Printf("\nResponse Body: %v", resp)
    fmt.Printf("\nResponse Time: %v", resp.Time())
    fmt.Printf("\nResponse Received At: %v", resp.ReceivedAt())
    

Example (Post) ¶

    // Create a resty client
    client := resty.New()
    
    // POST JSON string
    // No need to set content type, if you have client level setting
    resp, err := client.R().
     SetHeader("Content-Type", "application/json").
     SetBody(`{"username":"testuser", "password":"testpass"}`).
     SetResult(AuthSuccess{}). // or SetResult(&AuthSuccess{}).
     Post("https://myapp.com/login")
    
    printOutput(resp, err)
    
    // POST []byte array
    // No need to set content type, if you have client level setting
    resp1, err1 := client.R().
     SetHeader("Content-Type", "application/json").
     SetBody([]byte(`{"username":"testuser", "password":"testpass"}`)).
     SetResult(AuthSuccess{}). // or SetResult(&AuthSuccess{}).
     Post("https://myapp.com/login")
    
    printOutput(resp1, err1)
    
    // POST Struct, default is JSON content type. No need to set one
    resp2, err2 := client.R().
     SetBody(resty.User{Username: "testuser", Password: "testpass"}).
     SetResult(&AuthSuccess{}). // or SetResult(AuthSuccess{}).
     SetError(&AuthError{}).    // or SetError(AuthError{}).
     Post("https://myapp.com/login")
    
    printOutput(resp2, err2)
    
    // POST Map, default is JSON content type. No need to set one
    resp3, err3 := client.R().
     SetBody(map[string]interface{}{"username": "testuser", "password": "testpass"}).
     SetResult(&AuthSuccess{}). // or SetResult(AuthSuccess{}).
     SetError(&AuthError{}).    // or SetError(AuthError{}).
     Post("https://myapp.com/login")
    
    printOutput(resp3, err3)
    

Example (Put) ¶

    // Create a resty client
    client := resty.New()
    
    // Just one sample of PUT, refer POST for more combination
    // request goes as JSON content type
    // No need to set auth token, error, if you have client level settings
    resp, err := client.R().
     SetBody(Article{
      Title:   "go-resty",
      Content: "This is my article content, oh ya!",
      Author:  "Jeevanandam M",
      Tags:    []string{"article", "sample", "resty"},
     }).
     SetAuthToken("C6A79608-782F-4ED0-A11D-BD82FAD829CD").
     SetError(&Error{}). // or SetError(Error{}).
     Put("https://myapp.com/article/1234")
    
    printOutput(resp, err)
    

Example (Socks5Proxy) ¶

    // create a dialer
    dialer, err := proxy.SOCKS5("tcp", "127.0.0.1:9150", nil, proxy.Direct)
    if err != nil {
     log.Fatalf("Unable to obtain proxy dialer: %v\n", err)
    }
    
    // create a transport
    ptransport := &http.Transport{Dial: dialer.Dial}
    
    // Create a resty client
    client := resty.New()
    
    // set transport into resty
    client.SetTransport(ptransport)
    
    resp, err := client.R().Get("http://check.torproject.org")
    fmt.Println(err, resp)
    

### Index ¶

- Constants
- Variables
- func Backoff(operation func() (*Response, error), options ...Option) error
- func DetectContentType(body interface{}) string
- func IsJSONType(ct string) bool
- func IsStringEmpty(str string) bool
- func IsXMLType(ct string) bool
- func Unmarshalc(c *Client, ct string, b []byte, d interface{}) (err error)
- type Client
-     * func New() *Client
  - func NewWithClient(hc *http.Client)*Client
  - func NewWithLocalAddr(localAddr net.Addr) *Client
-     * func (c *Client) AddRetryAfterErrorCondition() *Client
  - func (c *Client) AddRetryCondition(condition RetryConditionFunc)*Client
  - func (c *Client) AddRetryHook(hook OnRetryFunc)*Client
  - func (c *Client) Clone()*Client
  - func (c *Client) DisableGenerateCurlOnDebug()*Client
  - func (c *Client) DisableTrace()*Client
  - func (c *Client) EnableGenerateCurlOnDebug()*Client
  - func (c *Client) EnableTrace()*Client
  - func (c *Client) GetClient()*http.Client
  - func (c *Client) IsProxySet() bool
  - func (c *Client) NewRequest()*Request
  - func (c *Client) OnAfterResponse(m ResponseMiddleware)*Client
  - func (c *Client) OnBeforeRequest(m RequestMiddleware)*Client
  - func (c *Client) OnError(h ErrorHook)*Client
  - func (c *Client) OnInvalid(h ErrorHook)*Client
  - func (c *Client) OnPanic(h ErrorHook)*Client
  - func (c *Client) OnRequestLog(rl RequestLogCallback)*Client
  - func (c *Client) OnResponseLog(rl ResponseLogCallback)*Client
  - func (c *Client) OnSuccess(h SuccessHook)*Client
  - func (c *Client) R()*Request
  - func (c *Client) RemoveProxy()*Client
  - func (c *Client) SetAllowGetMethodPayload(a bool)*Client
  - func (c *Client) SetAuthScheme(scheme string)*Client
  - func (c *Client) SetAuthToken(token string)*Client
  - func (c *Client) SetBaseURL(url string)*Client
  - func (c *Client) SetBasicAuth(username, password string)*Client
  - func (c *Client) SetCertificates(certs ...tls.Certificate)*Client
  - func (c *Client) SetClientRootCertificate(pemFilePath string)*Client
  - func (c *Client) SetClientRootCertificateFromString(pemCerts string)*Client
  - func (c *Client) SetCloseConnection(close bool)*Client
  - func (c *Client) SetContentLength(l bool)*Client
  - func (c *Client) SetCookie(hc*http.Cookie) *Client
  - func (c *Client) SetCookieJar(jar http.CookieJar)*Client
  - func (c *Client) SetCookies(cs []*http.Cookie) *Client
  - func (c *Client) SetDebug(d bool)*Client
  - func (c *Client) SetDebugBodyLimit(sl int64)*Client
  - func (c *Client) SetDigestAuth(username, password string)*Client
  - func (c *Client) SetDisableWarn(d bool)*Client
  - func (c *Client) SetDoNotParseResponse(notParse bool)*Client
  - func (c *Client) SetError(err interface{})*Client
  - func (c *Client) SetFormData(data map[string]string)*Client
  - func (c *Client) SetHeader(header, value string)*Client
  - func (c *Client) SetHeaderVerbatim(header, value string)*Client
  - func (c *Client) SetHeaders(headers map[string]string)*Client
  - func (c *Client) SetHostURL(url string)*Clientdeprecated
  - func (c *Client) SetJSONEscapeHTML(b bool)*Client
  - func (c *Client) SetJSONMarshaler(marshaler func(v interface{}) ([]byte, error))*Client
  - func (c *Client) SetJSONUnmarshaler(unmarshaler func(data []byte, v interface{}) error)*Client
  - func (c *Client) SetLogger(l Logger)*Client
  - func (c *Client) SetOutputDirectory(dirPath string)*Client
  - func (c *Client) SetPathParam(param, value string)*Client
  - func (c *Client) SetPathParams(params map[string]string)*Client
  - func (c *Client) SetPreRequestHook(h PreRequestHook)*Client
  - func (c *Client) SetProxy(proxyURL string)*Client
  - func (c *Client) SetQueryParam(param, value string)*Client
  - func (c *Client) SetQueryParams(params map[string]string)*Client
  - func (c *Client) SetRateLimiter(rl RateLimiter)*Client
  - func (c *Client) SetRawPathParam(param, value string)*Client
  - func (c *Client) SetRawPathParams(params map[string]string)*Client
  - func (c *Client) SetRedirectPolicy(policies ...interface{})*Client
  - func (c *Client) SetResponseBodyLimit(v int)*Client
  - func (c *Client) SetRetryAfter(callback RetryAfterFunc)*Client
  - func (c *Client) SetRetryCount(count int)*Client
  - func (c *Client) SetRetryMaxWaitTime(maxWaitTime time.Duration)*Client
  - func (c *Client) SetRetryResetReaders(b bool)*Client
  - func (c *Client) SetRetryWaitTime(waitTime time.Duration)*Client
  - func (c *Client) SetRootCertificate(pemFilePath string)*Client
  - func (c *Client) SetRootCertificateFromString(pemCerts string)*Client
  - func (c *Client) SetScheme(scheme string)*Client
  - func (c *Client) SetTLSClientConfig(config*tls.Config) *Client
  - func (c *Client) SetTimeout(timeout time.Duration)*Client
  - func (c *Client) SetTransport(transport http.RoundTripper)*Client
  - func (c *Client) SetUnescapeQueryParams(unescape bool)*Client
  - func (c *Client) SetXMLMarshaler(marshaler func(v interface{}) ([]byte, error))*Client
  - func (c *Client) SetXMLUnmarshaler(unmarshaler func(data []byte, v interface{}) error)*Client
  - func (c *Client) Transport() (*http.Transport, error)
- type ErrorHook
- type File
-     * func (f *File) String() string
- type Logger
- type MultipartField
- type OnRetryFunc
- type Option
-     * func MaxWaitTime(value time.Duration) Option
  - func ResetMultipartReaders(value bool) Option
  - func Retries(value int) Option
  - func RetryConditions(conditions []RetryConditionFunc) Option
  - func RetryHooks(hooks []OnRetryFunc) Option
  - func WaitTime(value time.Duration) Option
- type Options
- type PreRequestHook
- type RateLimiter
- type RedirectPolicy
-     * func DomainCheckRedirectPolicy(hostnames ...string) RedirectPolicy
  - func FlexibleRedirectPolicy(noOfRedirect int) RedirectPolicy
  - func NoRedirectPolicy() RedirectPolicy
- type RedirectPolicyFunc
-     * func (f RedirectPolicyFunc) Apply(req *http.Request, via []*http.Request) error
- type Request
-     * func (r *Request) AddRetryCondition(condition RetryConditionFunc) *Request
  - func (r *Request) Context() context.Context
  - func (r *Request) Delete(url string) (*Response, error)
  - func (r *Request) DisableGenerateCurlOnDebug()*Request
  - func (r *Request) EnableGenerateCurlOnDebug()*Request
  - func (r *Request) EnableTrace()*Request
  - func (r *Request) Execute(method, url string) (*Response, error)
  - func (r *Request) ExpectContentType(contentType string)*Request
  - func (r *Request) ForceContentType(contentType string)*Request
  - func (r *Request) GenerateCurlCommand() string
  - func (r *Request) Get(url string) (*Response, error)
  - func (r *Request) Head(url string) (*Response, error)
  - func (r *Request) Options(url string) (*Response, error)
  - func (r *Request) Patch(url string) (*Response, error)
  - func (r *Request) Post(url string) (*Response, error)
  - func (r *Request) Put(url string) (*Response, error)
  - func (r *Request) Send() (*Response, error)
  - func (r *Request) SetAuthScheme(scheme string)*Request
  - func (r *Request) SetAuthToken(token string)*Request
  - func (r *Request) SetBasicAuth(username, password string)*Request
  - func (r *Request) SetBody(body interface{})*Request
  - func (r *Request) SetContentLength(l bool)*Request
  - func (r *Request) SetContext(ctx context.Context)*Request
  - func (r *Request) SetCookie(hc*http.Cookie) *Request
  - func (r *Request) SetCookies(rs []*http.Cookie) *Request
  - func (r *Request) SetDebug(d bool)*Request
  - func (r *Request) SetDigestAuth(username, password string)*Request
  - func (r *Request) SetDoNotParseResponse(parse bool)*Request
  - func (r *Request) SetError(err interface{})*Request
  - func (r *Request) SetFile(param, filePath string)*Request
  - func (r *Request) SetFileReader(param, fileName string, reader io.Reader)*Request
  - func (r *Request) SetFiles(files map[string]string)*Request
  - func (r *Request) SetFormData(data map[string]string)*Request
  - func (r *Request) SetFormDataFromValues(data url.Values)*Request
  - func (r *Request) SetHeader(header, value string)*Request
  - func (r *Request) SetHeaderMultiValues(headers map[string][]string)*Request
  - func (r *Request) SetHeaderVerbatim(header, value string)*Request
  - func (r *Request) SetHeaders(headers map[string]string)*Request
  - func (r *Request) SetJSONEscapeHTML(b bool)*Request
  - func (r *Request) SetLogger(l Logger)*Request
  - func (r *Request) SetMultipartBoundary(boundary string)*Request
  - func (r *Request) SetMultipartField(param, fileName, contentType string, reader io.Reader)*Request
  - func (r *Request) SetMultipartFields(fields ...*MultipartField) *Request
  - func (r *Request) SetMultipartFormData(data map[string]string)*Request
  - func (r *Request) SetOutput(file string)*Request
  - func (r *Request) SetPathParam(param, value string)*Request
  - func (r *Request) SetPathParams(params map[string]string)*Request
  - func (r *Request) SetQueryParam(param, value string)*Request
  - func (r *Request) SetQueryParams(params map[string]string)*Request
  - func (r *Request) SetQueryParamsFromValues(params url.Values)*Request
  - func (r *Request) SetQueryString(query string)*Request
  - func (r *Request) SetRawPathParam(param, value string)*Request
  - func (r *Request) SetRawPathParams(params map[string]string)*Request
  - func (r *Request) SetResponseBodyLimit(v int)*Request
  - func (r *Request) SetResult(res interface{})*Request
  - func (r *Request) SetSRV(srv*SRVRecord) *Request
  - func (r *Request) SetUnescapeQueryParams(unescape bool)*Request
  - func (r *Request) TraceInfo() TraceInfo
- type RequestLog
- type RequestLogCallback
- type RequestMiddleware
- type Response
-     * func (r *Response) Body() []byte
  - func (r *Response) Cookies() []*http.Cookie
  - func (r *Response) Error() interface{}
  - func (r *Response) Header() http.Header
  - func (r *Response) IsError() bool
  - func (r *Response) IsSuccess() bool
  - func (r *Response) Proto() string
  - func (r *Response) RawBody() io.ReadCloser
  - func (r *Response) ReceivedAt() time.Time
  - func (r *Response) Result() interface{}
  - func (r *Response) SetBody(b []byte)*Response
  - func (r *Response) Size() int64
  - func (r *Response) Status() string
  - func (r *Response) StatusCode() int
  - func (r *Response) String() string
  - func (r *Response) Time() time.Duration
- type ResponseError
-     * func (e *ResponseError) Error() string
  - func (e *ResponseError) Unwrap() error
- type ResponseLog
- type ResponseLogCallback
- type ResponseMiddleware
- type RetryAfterFunc
- type RetryConditionFunc
- type SRVRecord
- type SuccessHook
- type TraceInfo
- type User

### Examples ¶

- Package (ClientCertificates)
- Package (CustomRootCertificate)
- Package (DropboxUpload)
- Package (EnhancedGet)
- Package (Get)
- Package (Post)
- Package (Put)
- Package (Socks5Proxy)
- Client.SetCertificates
- New

### Constants ¶

[View Source](https://github.com/go-resty/resty/blob/v2.17.1/client.go#L28)

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
    )

[View Source](https://github.com/go-resty/resty/blob/v2.17.1/resty.go#L17)

    const Version = "2.17.1"

Version # of resty

### Variables ¶

[View Source](https://github.com/go-resty/resty/blob/v2.17.1/digest.go#L23)

    var (
     ErrDigestBadChallenge    = [errors](/errors).[New](/errors#New)("digest: challenge is bad")
     ErrDigestCharset         = [errors](/errors).[New](/errors#New)("digest: unsupported charset")
     ErrDigestAlgNotSupported = [errors](/errors).[New](/errors#New)("digest: algorithm is not supported")
     ErrDigestQopNotSupported = [errors](/errors).[New](/errors#New)("digest: no supported qop in list")
     ErrDigestNoQop           = [errors](/errors).[New](/errors#New)("digest: qop must be specified")
    )

[View Source](https://github.com/go-resty/resty/blob/v2.17.1/redirect.go#L15)

    var (
     ErrAutoRedirectDisabled = [errors](/errors).[New](/errors#New)("auto redirect is disabled")
    )

[View Source](https://github.com/go-resty/resty/blob/v2.17.1/util.go#L75)

    var ErrRateLimitExceeded = [errors](/errors).[New](/errors#New)("rate limit exceeded")

[View Source](https://github.com/go-resty/resty/blob/v2.17.1/client.go#L1312)

    var ErrResponseBodyTooLarge = [errors](/errors).[New](/errors#New)("resty: response body too large")

### Functions ¶

#### func [Backoff](https://github.com/go-resty/resty/blob/v2.17.1/retry.go#L96) ¶

    func Backoff(operation func() (*Response, [error](/builtin#error)), options ...Option) [error](/builtin#error)

Backoff retries with increasing timeout duration up until X amount of retries (Default is 3 attempts, Override with option Retries(n))

#### func [DetectContentType](https://github.com/go-resty/resty/blob/v2.17.1/util.go#L87) ¶

    func DetectContentType(body interface{}) [string](/builtin#string)

DetectContentType method is used to figure out `Request.Body` content type for request header

#### func [IsJSONType](https://github.com/go-resty/resty/blob/v2.17.1/util.go#L107) ¶

    func IsJSONType(ct [string](/builtin#string)) [bool](/builtin#bool)

IsJSONType method is to check JSON content type or not

#### func [IsStringEmpty](https://github.com/go-resty/resty/blob/v2.17.1/util.go#L82) ¶

    func IsStringEmpty(str [string](/builtin#string)) [bool](/builtin#bool)

IsStringEmpty method tells whether given string is empty or not

#### func [IsXMLType](https://github.com/go-resty/resty/blob/v2.17.1/util.go#L112) ¶

    func IsXMLType(ct [string](/builtin#string)) [bool](/builtin#bool)

IsXMLType method is to check XML content type or not

#### func [Unmarshalc](https://github.com/go-resty/resty/blob/v2.17.1/util.go#L117) ¶

    func Unmarshalc(c *Client, ct [string](/builtin#string), b [][byte](/builtin#byte), d interface{}) (err [error](/builtin#error))

Unmarshalc content into object from JSON or XML

### Types ¶

#### type [Client](https://github.com/go-resty/resty/blob/v2.17.1/client.go#L100) ¶

    type Client struct {
     BaseURL               [string](/builtin#string)
     HostURL               [string](/builtin#string) // Deprecated: use BaseURL instead. To be removed in v3.0.0 release.
     QueryParam            [url](/net/url).[Values](/net/url#Values)
     FormData              [url](/net/url).[Values](/net/url#Values)
     PathParams            map[[string](/builtin#string)][string](/builtin#string)
     RawPathParams         map[[string](/builtin#string)][string](/builtin#string)
     Header                [http](/net/http).[Header](/net/http#Header)
     UserInfo              *User
     Token                 [string](/builtin#string)
     AuthScheme            [string](/builtin#string)
     Cookies               []*[http](/net/http).[Cookie](/net/http#Cookie)
     Error                 [reflect](/reflect).[Type](/reflect#Type)
     Debug                 [bool](/builtin#bool)
     DisableWarn           [bool](/builtin#bool)
     AllowGetMethodPayload [bool](/builtin#bool)
     RetryCount            [int](/builtin#int)
     RetryWaitTime         [time](/time).[Duration](/time#Duration)
     RetryMaxWaitTime      [time](/time).[Duration](/time#Duration)
     RetryConditions       []RetryConditionFunc
     RetryHooks            []OnRetryFunc
     RetryAfter            RetryAfterFunc
     RetryResetReaders     [bool](/builtin#bool)
     JSONMarshal           func(v interface{}) ([][byte](/builtin#byte), [error](/builtin#error))
     JSONUnmarshal         func(data [][byte](/builtin#byte), v interface{}) [error](/builtin#error)
     XMLMarshal            func(v interface{}) ([][byte](/builtin#byte), [error](/builtin#error))
     XMLUnmarshal          func(data [][byte](/builtin#byte), v interface{}) [error](/builtin#error)
    
     // HeaderAuthorizationKey is used to set/access Request Authorization header
     // value when `SetAuthToken` option is used.
     HeaderAuthorizationKey [string](/builtin#string)
     ResponseBodyLimit      [int](/builtin#int)
     // contains filtered or unexported fields
    }

Client struct is used to create a Resty client with client-level settings, these settings apply to all the requests raised from the client.

Resty also provides an option to override most of the client settings at Request level.

#### func [New](https://github.com/go-resty/resty/blob/v2.17.1/resty.go#L20) ¶

    func New() *Client

New method creates a new Resty client.

Example ¶

    // Creating client1
    client1 := resty.New()
    resp1, err1 := client1.R().Get("http://httpbin.org/get")
    fmt.Println(resp1, err1)
    
    // Creating client2
    client2 := resty.New()
    resp2, err2 := client2.R().Get("http://httpbin.org/get")
    fmt.Println(resp2, err2)
    

#### func [NewWithClient](https://github.com/go-resty/resty/blob/v2.17.1/resty.go#L28) ¶

    func NewWithClient(hc *[http](/net/http).[Client](/net/http#Client)) *Client

NewWithClient method creates a new Resty client with given [http.Client](/net/http#Client).

#### func [NewWithLocalAddr](https://github.com/go-resty/resty/blob/v2.17.1/resty.go#L34) ¶

    func NewWithLocalAddr(localAddr [net](/net).[Addr](/net#Addr)) *Client

NewWithLocalAddr method creates a new Resty client with the given Local Address. to dial from.

#### func (*Client) [AddRetryAfterErrorCondition](https://github.com/go-resty/resty/blob/v2.17.1/client.go#L779) ¶ added in v2.6.0

    func (c *Client) AddRetryAfterErrorCondition() *Client

AddRetryAfterErrorCondition adds the basic condition of retrying after encountering an error from the HTTP response

#### func (*Client) [AddRetryCondition](https://github.com/go-resty/resty/blob/v2.17.1/client.go#L772) ¶

    func (c *Client) AddRetryCondition(condition RetryConditionFunc) *Client

AddRetryCondition method adds a retry condition function to an array of functions that are checked to determine if the request is retried. The request will retry if any functions return true and the error is nil.

NOTE: These retry conditions are applied on all requests made using this Client. For Request specific retry conditions, check Request.AddRetryCondition

#### func (*Client) [AddRetryHook](https://github.com/go-resty/resty/blob/v2.17.1/client.go#L788) ¶ added in v2.6.0

    func (c *Client) AddRetryHook(hook OnRetryFunc) *Client

AddRetryHook adds a side-effecting retry hook to an array of hooks that will be executed on each retry.

#### func (*Client) [Clone](https://github.com/go-resty/resty/blob/v2.17.1/client.go#L1180) ¶ added in v2.12.0

    func (c *Client) Clone() *Client

Clone returns a clone of the original client.

NOTE: Use with care:

- Interface values are not deeply cloned. Thus, both the original and the clone will use the same value.
- This function is not safe for concurrent use. You should only use this method when you are sure that any other goroutine is not using the client.

#### func (*Client) [DisableGenerateCurlOnDebug](https://github.com/go-resty/resty/blob/v2.17.1/client.go#L1157) ¶ added in v2.15.0

    func (c *Client) DisableGenerateCurlOnDebug() *Client

DisableGenerateCurlOnDebug method disables the option set by Client.EnableGenerateCurlOnDebug.

#### func (*Client) [DisableTrace](https://github.com/go-resty/resty/blob/v2.17.1/client.go#L1140) ¶

    func (c *Client) DisableTrace() *Client

DisableTrace method disables the Resty client trace. Refer to Client.EnableTrace.

#### func (*Client) [EnableGenerateCurlOnDebug](https://github.com/go-resty/resty/blob/v2.17.1/client.go#L1151) ¶ added in v2.15.0

    func (c *Client) EnableGenerateCurlOnDebug() *Client

EnableGenerateCurlOnDebug method enables the generation of CURL commands in the debug log. It works in conjunction with debug mode.

NOTE: Use with care.

- Potential to leak sensitive data from Request and Response in the debug log.
- Beware of memory usage since the request body is reread.

#### func (*Client) [EnableTrace](https://github.com/go-resty/resty/blob/v2.17.1/client.go#L1134) ¶

    func (c *Client) EnableTrace() *Client

EnableTrace method enables the Resty client trace for the requests fired from the client using [httptrace.ClientTrace](/net/http/httptrace#ClientTrace) and provides insights.

    client := resty.New().EnableTrace()
    
    resp, err := client.R().Get("https://httpbin.org/get")
    fmt.Println("Error:", err)
    fmt.Println("Trace Info:", resp.Request.TraceInfo())
    

The method Request.EnableTrace is also available to get trace info for a single request.

#### func (*Client) [GetClient](https://github.com/go-resty/resty/blob/v2.17.1/client.go#L1169) ¶

    func (c *Client) GetClient() *[http](/net/http).[Client](/net/http#Client)

GetClient method returns the underlying [http.Client](/net/http#Client) used by the Resty.

#### func (*Client) [IsProxySet](https://github.com/go-resty/resty/blob/v2.17.1/client.go#L1164) ¶

    func (c *Client) IsProxySet() [bool](/builtin#bool)

IsProxySet method returns the true is proxy is set from the Resty client; otherwise false. By default, the proxy is set from the environment variable; refer to [http.ProxyFromEnvironment](/net/http#ProxyFromEnvironment).

#### func (*Client) [NewRequest](https://github.com/go-resty/resty/blob/v2.17.1/client.go#L468) ¶

    func (c *Client) NewRequest() *Request

NewRequest method is an alias for method `R()`.

#### func (*Client) [OnAfterResponse](https://github.com/go-resty/resty/blob/v2.17.1/client.go#L501) ¶

    func (c *Client) OnAfterResponse(m ResponseMiddleware) *Client

OnAfterResponse method appends response middleware to the after-response chain. Once we receive a response from the host server, the default Resty response middleware gets applied, and then the user-assigned response middleware is applied.

    client.OnAfterResponse(func(c *resty.Client, r *resty.Response) error {
      // Now you have access to the Client and Response instance
      // manipulate it as per your need
    
      return nil  // if its successful otherwise return error
     })
    

#### func (*Client) [OnBeforeRequest](https://github.com/go-resty/resty/blob/v2.17.1/client.go#L482) ¶

    func (c *Client) OnBeforeRequest(m RequestMiddleware) *Client

OnBeforeRequest method appends a request middleware to the before request chain. The user-defined middlewares are applied before the default Resty request middlewares. After all middlewares have been applied, the request is sent from Resty to the host server.

    client.OnBeforeRequest(func(c *resty.Client, r *resty.Request) error {
      // Now you have access to the Client and Request instance
      // manipulate it as per your need
    
      return nil  // if its successful otherwise return error
     })
    

#### func (*Client) [OnError](https://github.com/go-resty/resty/blob/v2.17.1/client.go#L524) ¶ added in v2.4.0

    func (c *Client) OnError(h ErrorHook) *Client

OnError method adds a callback that will be run whenever a request execution fails. This is called after all retries have been attempted (if any). If there was a response from the server, the error will be wrapped in ResponseError which has the last response received from the server.

    client.OnError(func(req *resty.Request, err error) {
     if v, ok := err.(*resty.ResponseError); ok {
      // Do something with v.Response
     }
     // Log the error, increment a metric, etc...
    })
    

Out of the Client.OnSuccess, Client.OnError, Client.OnInvalid, Client.OnPanic callbacks, exactly one set will be invoked for each call to Request.Execute that completes.

#### func (*Client) [OnInvalid](https://github.com/go-resty/resty/blob/v2.17.1/client.go#L544) ¶ added in v2.8.0

    func (c *Client) OnInvalid(h ErrorHook) *Client

OnInvalid method adds a callback that will be run whenever a request execution fails before it starts because the request is invalid.

Out of the Client.OnSuccess, Client.OnError, Client.OnInvalid, Client.OnPanic callbacks, exactly one set will be invoked for each call to Request.Execute that completes.

#### func (*Client) [OnPanic](https://github.com/go-resty/resty/blob/v2.17.1/client.go#L557) ¶ added in v2.8.0

    func (c *Client) OnPanic(h ErrorHook) *Client

OnPanic method adds a callback that will be run whenever a request execution panics.

Out of the Client.OnSuccess, Client.OnError, Client.OnInvalid, Client.OnPanic callbacks, exactly one set will be invoked for each call to Request.Execute that completes.

If an Client.OnSuccess, Client.OnError, or Client.OnInvalid callback panics, then exactly one rule can be violated.

#### func (*Client) [OnRequestLog](https://github.com/go-resty/resty/blob/v2.17.1/client.go#L600) ¶

    func (c *Client) OnRequestLog(rl RequestLogCallback) *Client

OnRequestLog method sets the request log callback to Resty. Registered callback gets called before the resty logs the information.

#### func (*Client) [OnResponseLog](https://github.com/go-resty/resty/blob/v2.17.1/client.go#L611) ¶

    func (c *Client) OnResponseLog(rl ResponseLogCallback) *Client

OnResponseLog method sets the response log callback to Resty. Registered callback gets called before the resty logs the information.

#### func (*Client) [OnSuccess](https://github.com/go-resty/resty/blob/v2.17.1/client.go#L534) ¶ added in v2.8.0

    func (c *Client) OnSuccess(h SuccessHook) *Client

OnSuccess method adds a callback that will be run whenever a request execution succeeds. This is called after all retries have been attempted (if any).

Out of the Client.OnSuccess, Client.OnError, Client.OnInvalid, Client.OnPanic callbacks, exactly one set will be invoked for each call to Request.Execute that completes.

#### func (*Client) [R](https://github.com/go-resty/resty/blob/v2.17.1/client.go#L444) ¶

    func (c *Client) R() *Request

R method creates a new request instance; it's used for Get, Post, Put, Delete, Patch, Head, Options, etc.

#### func (*Client) [RemoveProxy](https://github.com/go-resty/resty/blob/v2.17.1/client.go#L847) ¶

    func (c *Client) RemoveProxy() *Client

RemoveProxy method removes the proxy configuration from the Resty client

    client.RemoveProxy()
    

#### func (*Client) [SetAllowGetMethodPayload](https://github.com/go-resty/resty/blob/v2.17.1/client.go#L635) ¶

    func (c *Client) SetAllowGetMethodPayload(a [bool](/builtin#bool)) *Client

SetAllowGetMethodPayload method allows the GET method with payload on the Resty client.

For example, Resty allows the user to send a request with a payload using the HTTP GET method.

    client.SetAllowGetMethodPayload(true)
    

#### func (*Client) [SetAuthScheme](https://github.com/go-resty/resty/blob/v2.17.1/client.go#L410) ¶ added in v2.3.0

    func (c *Client) SetAuthScheme(scheme [string](/builtin#string)) *Client

SetAuthScheme method sets the auth scheme type in the HTTP request. For Example:

    Authorization: <auth-scheme-value> <auth-token-value>
    

For Example: To set the scheme to use OAuth

    client.SetAuthScheme("OAuth")
    

This auth scheme gets added to all the requests raised from this client instance. Also, it can be overridden at the request level.

Information about auth schemes can be found in [](https://tools.ietf.org/html/rfc7235)[RFC 7235](https://rfc-editor.org/rfc/rfc7235.html), IANA [HTTP Auth schemes](https://www.iana.org/assignments/http-authschemes/http-authschemes.xhtml#authschemes).

See Request.SetAuthToken.

#### func (*Client) [SetAuthToken](https://github.com/go-resty/resty/blob/v2.17.1/client.go#L388) ¶

    func (c *Client) SetAuthToken(token [string](/builtin#string)) *Client

SetAuthToken method sets the auth token of the `Authorization` header for all HTTP requests. The default auth scheme is `Bearer`; it can be customized with the method Client.SetAuthScheme. For Example:

    Authorization: <auth-scheme> <auth-token-value>
    

For Example: To set auth token BC594900518B4F7EAC75BD37F019E08FBC594900518B4F7EAC75BD37F019E08F

    client.SetAuthToken("BC594900518B4F7EAC75BD37F019E08FBC594900518B4F7EAC75BD37F019E08F")
    

This auth token gets added to all the requests raised from this client instance. Also, it can be overridden at the request level.

See Request.SetAuthToken.

#### func (*Client) [SetBaseURL](https://github.com/go-resty/resty/blob/v2.17.1/client.go#L193) ¶ added in v2.7.0

    func (c *Client) SetBaseURL(url [string](/builtin#string)) *Client

SetBaseURL method sets the Base URL in the client instance. It will be used with a request raised from this client with a relative URL

    // Setting HTTP address
    client.SetBaseURL("http://myjeeva.com")
    
    // Setting HTTPS address
    client.SetBaseURL("https://myjeeva.com")
    

#### func (*Client) [SetBasicAuth](https://github.com/go-resty/resty/blob/v2.17.1/client.go#L370) ¶

    func (c *Client) SetBasicAuth(username, password [string](/builtin#string)) *Client

SetBasicAuth method sets the basic authentication header in the HTTP request. For Example:

    Authorization: Basic <base64-encoded-value>
    

For Example: To set the header for username "go-resty" and password "welcome"

    client.SetBasicAuth("go-resty", "welcome")
    

This basic auth information is added to all requests from this client instance. It can also be overridden at the request level.

See Request.SetBasicAuth.

#### func (*Client) [SetCertificates](https://github.com/go-resty/resty/blob/v2.17.1/client.go#L859) ¶

    func (c *Client) SetCertificates(certs ...[tls](/crypto/tls).[Certificate](/crypto/tls#Certificate)) *Client

SetCertificates method helps to conveniently set client certificates into Resty.

Example ¶

    // Parsing public/private key pair from a pair of files. The files must contain PEM encoded data.
    cert, err := tls.LoadX509KeyPair("certs/client.pem", "certs/client.key")
    if err != nil {
     log.Fatalf("ERROR client certificate: %s", err)
    }
    
    // Create a resty client
    client := resty.New()
    
    client.SetCertificates(cert)
    

#### func (*Client) [SetClientRootCertificate](https://github.com/go-resty/resty/blob/v2.17.1/client.go#L895) ¶ added in v2.15.0

    func (c *Client) SetClientRootCertificate(pemFilePath [string](/builtin#string)) *Client

SetClientRootCertificate method helps to add one or more client's root certificates into the Resty client

    client.SetClientRootCertificate("/path/to/root/pemFile.pem")
    

#### func (*Client) [SetClientRootCertificateFromString](https://github.com/go-resty/resty/blob/v2.17.1/client.go#L909) ¶ added in v2.15.0

    func (c *Client) SetClientRootCertificateFromString(pemCerts [string](/builtin#string)) *Client

SetClientRootCertificateFromString method helps to add one or more clients root certificates into the Resty client

    client.SetClientRootCertificateFromString("pem certs content")
    

#### func (*Client) [SetCloseConnection](https://github.com/go-resty/resty/blob/v2.17.1/client.go#L986) ¶

    func (c *Client) SetCloseConnection(close [bool](/builtin#bool)) *Client

SetCloseConnection method sets variable `Close` in HTTP request struct with the given value. More info: <https://golang.org/src/net/http/request.go>

#### func (*Client) [SetContentLength](https://github.com/go-resty/resty/blob/v2.17.1/client.go#L654) ¶

    func (c *Client) SetContentLength(l [bool](/builtin#bool)) *Client

SetContentLength method enables the HTTP header `Content-Length` value for every request. By default, Resty won't set `Content-Length`.

    client.SetContentLength(true)
    

Also, you have the option to enable a particular request. See Request.SetContentLength

#### func (*Client) [SetCookie](https://github.com/go-resty/resty/blob/v2.17.1/client.go#L264) ¶

    func (c *Client) SetCookie(hc *[http](/net/http).[Cookie](/net/http#Cookie)) *Client

SetCookie method appends a single cookie to the client instance. These cookies will be added to all the requests from this client instance.

    client.SetCookie(&http.Cookie{
       Name:"go-resty",
       Value:"This is cookie value",
      })
    

#### func (*Client) [SetCookieJar](https://github.com/go-resty/resty/blob/v2.17.1/client.go#L252) ¶

    func (c *Client) SetCookieJar(jar [http](/net/http).[CookieJar](/net/http#CookieJar)) *Client

SetCookieJar method sets custom [http.CookieJar](/net/http#CookieJar) in the resty client. It's a way to override the default.

For Example, sometimes we don't want to save cookies in API mode so that we can remove the default CookieJar in resty client.

    client.SetCookieJar(nil)
    

#### func (*Client) [SetCookies](https://github.com/go-resty/resty/blob/v2.17.1/client.go#L285) ¶

    func (c *Client) SetCookies(cs []*[http](/net/http).[Cookie](/net/http#Cookie)) *Client

SetCookies method sets an array of cookies in the client instance. These cookies will be added to all the requests from this client instance.

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
    

#### func (*Client) [SetDebug](https://github.com/go-resty/resty/blob/v2.17.1/client.go#L584) ¶

    func (c *Client) SetDebug(d [bool](/builtin#bool)) *Client

SetDebug method enables the debug mode on the Resty client. The client logs details of every request and response.

    client.SetDebug(true)
    

Also, it can be enabled at the request level for a particular request; see Request.SetDebug.

- For Request, it logs information such as HTTP verb, Relative URL path, Host, Headers, and Body if it has one.
- For Response, it logs information such as Status, Response Time, Headers, and Body if it has one.

#### func (*Client) [SetDebugBodyLimit](https://github.com/go-resty/resty/blob/v2.17.1/client.go#L593) ¶

    func (c *Client) SetDebugBodyLimit(sl [int64](/builtin#int64)) *Client

SetDebugBodyLimit sets the maximum size in bytes for which the response and request body will be logged in debug mode.

    client.SetDebugBodyLimit(1000000)
    

#### func (*Client) [SetDigestAuth](https://github.com/go-resty/resty/blob/v2.17.1/client.go#L427) ¶ added in v2.8.0

    func (c *Client) SetDigestAuth(username, password [string](/builtin#string)) *Client

SetDigestAuth method sets the Digest Access auth scheme for the client. If a server responds with 401 and sends a Digest challenge in the WWW-Authenticate Header, requests will be resent with the appropriate Authorization Header.

For Example: To set the Digest scheme with user "Mufasa" and password "Circle Of Life"

    client.SetDigestAuth("Mufasa", "Circle Of Life")
    

Information about Digest Access Authentication can be found in [](https://datatracker.ietf.org/doc/html/rfc7616)[RFC 7616](https://rfc-editor.org/rfc/rfc7616.html).

See Request.SetDigestAuth.

#### func (*Client) [SetDisableWarn](https://github.com/go-resty/resty/blob/v2.17.1/client.go#L625) ¶

    func (c *Client) SetDisableWarn(d [bool](/builtin#bool)) *Client

SetDisableWarn method disables the warning log message on the Resty client.

For example, Resty warns users when BasicAuth is used in non-TLS mode.

    client.SetDisableWarn(true)
    

#### func (*Client) [SetDoNotParseResponse](https://github.com/go-resty/resty/blob/v2.17.1/client.go#L998) ¶

    func (c *Client) SetDoNotParseResponse(notParse [bool](/builtin#bool)) *Client

SetDoNotParseResponse method instructs Resty not to parse the response body automatically. Resty exposes the raw response body as [io.ReadCloser](/io#ReadCloser). If you use it, do not forget to close the body, otherwise, you might get into connection leaks, and connection reuse may not happen.

NOTE: Response middlewares are not executed using this option. You have taken over the control of response parsing from Resty.

#### func (*Client) [SetError](https://github.com/go-resty/resty/blob/v2.17.1/client.go#L674) ¶

    func (c *Client) SetError(err interface{}) *Client

SetError method registers the global or client common `Error` object into Resty. It is used for automatic unmarshalling if the response status code is greater than 399 and content type is JSON or XML. It can be a pointer or a non-pointer.

    client.SetError(&Error{})
    // OR
    client.SetError(Error{})
    

#### func (*Client) [SetFormData](https://github.com/go-resty/resty/blob/v2.17.1/client.go#L351) ¶

    func (c *Client) SetFormData(data map[[string](/builtin#string)][string](/builtin#string)) *Client

SetFormData method sets Form parameters and their values in the client instance. It applies only to HTTP methods `POST` and `PUT`, and the request content type would be set as `application/x-www-form-urlencoded`. These form data will be added to all the requests raised from this client instance. Also, it can be overridden at the request level.

See Request.SetFormData.

    client.SetFormData(map[string]string{
      "access_token": "BC594900-518B-4F7E-AC75-BD37F019E08F",
      "user_id": "3455454545",
     })
    

#### func (*Client) [SetHeader](https://github.com/go-resty/resty/blob/v2.17.1/client.go#L210) ¶

    func (c *Client) SetHeader(header, value [string](/builtin#string)) *Client

SetHeader method sets a single header field and its value in the client instance. These headers will be applied to all requests from this client instance. Also, it can be overridden by request-level header options.

See Request.SetHeader or Request.SetHeaders.

For Example: To set `Content-Type` and `Accept` as `application/json`

    client.
     SetHeader("Content-Type", "application/json").
     SetHeader("Accept", "application/json")
    

#### func (*Client) [SetHeaderVerbatim](https://github.com/go-resty/resty/blob/v2.17.1/client.go#L241) ¶ added in v2.6.0

    func (c *Client) SetHeaderVerbatim(header, value [string](/builtin#string)) *Client

SetHeaderVerbatim method sets a single header field and its value verbatim in the current request.

For Example: To set `all_lowercase` and `UPPERCASE` as `available`.

    client.
     SetHeaderVerbatim("all_lowercase", "available").
     SetHeaderVerbatim("UPPERCASE", "available")
    

#### func (*Client) [SetHeaders](https://github.com/go-resty/resty/blob/v2.17.1/client.go#L227) ¶

    func (c *Client) SetHeaders(headers map[[string](/builtin#string)][string](/builtin#string)) *Client

SetHeaders method sets multiple header fields and their values at one go in the client instance. These headers will be applied to all requests from this client instance. Also, it can be overridden at request level headers options.

See Request.SetHeaders or Request.SetHeader.

For Example: To set `Content-Type` and `Accept` as `application/json`

    client.SetHeaders(map[string]string{
      "Content-Type": "application/json",
      "Accept": "application/json",
     })
    

#### func (*Client) [SetHostURL](https://github.com/go-resty/resty/blob/v2.17.1/client.go#L180) deprecated

    func (c *Client) SetHostURL(url [string](/builtin#string)) *Client

SetHostURL method sets the Host URL in the client instance. It will be used with a request raised from this client with a relative URL

    // Setting HTTP address
    client.SetHostURL("http://myjeeva.com")
    
    // Setting HTTPS address
    client.SetHostURL("https://myjeeva.com")
    

Deprecated: use Client.SetBaseURL instead. To be removed in the v3.0.0 release.

#### func (*Client) [SetJSONEscapeHTML](https://github.com/go-resty/resty/blob/v2.17.1/client.go#L1103) ¶

    func (c *Client) SetJSONEscapeHTML(b [bool](/builtin#bool)) *Client

SetJSONEscapeHTML method enables or disables the HTML escape on JSON marshal. By default, escape HTML is false.

NOTE: This option only applies to the standard JSON Marshaller used by Resty.

It can be overridden at the request level, see Client.SetJSONEscapeHTML

#### func (*Client) [SetJSONMarshaler](https://github.com/go-resty/resty/blob/v2.17.1/client.go#L740) ¶ added in v2.8.0

    func (c *Client) SetJSONMarshaler(marshaler func(v interface{}) ([][byte](/builtin#byte), [error](/builtin#error))) *Client

SetJSONMarshaler method sets the JSON marshaler function to marshal the request body. By default, Resty uses [encoding/json](/encoding/json) package to marshal the request body.

#### func (*Client) [SetJSONUnmarshaler](https://github.com/go-resty/resty/blob/v2.17.1/client.go#L747) ¶ added in v2.8.0

    func (c *Client) SetJSONUnmarshaler(unmarshaler func(data [][byte](/builtin#byte), v interface{}) [error](/builtin#error)) *Client

SetJSONUnmarshaler method sets the JSON unmarshaler function to unmarshal the response body. By default, Resty uses [encoding/json](/encoding/json) package to unmarshal the response body.

#### func (*Client) [SetLogger](https://github.com/go-resty/resty/blob/v2.17.1/client.go#L643) ¶

    func (c *Client) SetLogger(l Logger) *Client

SetLogger method sets given writer for logging Resty request and response details.

Compliant to interface resty.Logger

#### func (*Client) [SetOutputDirectory](https://github.com/go-resty/resty/blob/v2.17.1/client.go#L940) ¶

    func (c *Client) SetOutputDirectory(dirPath [string](/builtin#string)) *Client

SetOutputDirectory method sets the output directory for saving HTTP responses in a file. Resty creates one if the output directory does not exist. This setting is optional, if you plan to use the absolute path in Request.SetOutput and can used together.

    client.SetOutputDirectory("/save/http/response/here")
    

#### func (*Client) [SetPathParam](https://github.com/go-resty/resty/blob/v2.17.1/client.go#L1017) ¶ added in v2.4.0

    func (c *Client) SetPathParam(param, value [string](/builtin#string)) *Client

SetPathParam method sets a single URL path key-value pair in the Resty client instance.

    client.SetPathParam("userId", "sample@sample.com")
    
    Result:
       URL - /v1/users/{userId}/details
       Composed URL - /v1/users/sample@sample.com/details
    

It replaces the value of the key while composing the request URL. The value will be escaped using [url.PathEscape](/net/url#PathEscape) function.

It can be overridden at the request level, see Request.SetPathParam or Request.SetPathParams

#### func (*Client) [SetPathParams](https://github.com/go-resty/resty/blob/v2.17.1/client.go#L1040) ¶

    func (c *Client) SetPathParams(params map[[string](/builtin#string)][string](/builtin#string)) *Client

SetPathParams method sets multiple URL path key-value pairs at one go in the Resty client instance.

    client.SetPathParams(map[string]string{
     "userId":       "sample@sample.com",
     "subAccountId": "100002",
     "path":         "groups/developers",
    })
    
    Result:
       URL - /v1/users/{userId}/{subAccountId}/{path}/details
       Composed URL - /v1/users/sample@sample.com/100002/groups%2Fdevelopers/details
    

It replaces the value of the key while composing the request URL. The values will be escaped using [url.PathEscape](/net/url#PathEscape) function.

It can be overridden at the request level, see Request.SetPathParam or Request.SetPathParams

#### func (*Client) [SetPreRequestHook](https://github.com/go-resty/resty/blob/v2.17.1/client.go#L566) ¶

    func (c *Client) SetPreRequestHook(h PreRequestHook) *Client

SetPreRequestHook method sets the given pre-request function into a resty client. It is called right before the request is fired.

NOTE: Only one pre-request hook can be registered. Use Client.OnBeforeRequest for multiple.

#### func (*Client) [SetProxy](https://github.com/go-resty/resty/blob/v2.17.1/client.go#L826) ¶

    func (c *Client) SetProxy(proxyURL [string](/builtin#string)) *Client

SetProxy method sets the Proxy URL and Port for the Resty client.

    client.SetProxy("http://proxyserver:8888")
    

OR you could also set Proxy via environment variable, refer to [http.ProxyFromEnvironment](/net/http#ProxyFromEnvironment)

#### func (*Client) [SetQueryParam](https://github.com/go-resty/resty/blob/v2.17.1/client.go#L303) ¶

    func (c *Client) SetQueryParam(param, value [string](/builtin#string)) *Client

SetQueryParam method sets a single parameter and its value in the client instance. It will be formed as a query string for the request.

    For Example: `search=kitchen%20papers&size=large`
    

In the URL after the `?` mark. These query params will be added to all the requests raised from this client instance. Also, it can be overridden at the request level.

See Request.SetQueryParam or Request.SetQueryParams.

    client.
     SetQueryParam("search", "kitchen papers").
     SetQueryParam("size", "large")
    

#### func (*Client) [SetQueryParams](https://github.com/go-resty/resty/blob/v2.17.1/client.go#L322) ¶

    func (c *Client) SetQueryParams(params map[[string](/builtin#string)][string](/builtin#string)) *Client

SetQueryParams method sets multiple parameters and their values at one go in the client instance. It will be formed as a query string for the request.

    For Example: `search=kitchen%20papers&size=large`
    

In the URL after the `?` mark. These query params will be added to all the requests raised from this client instance. Also, it can be overridden at the request level.

See Request.SetQueryParams or Request.SetQueryParam.

    client.SetQueryParams(map[string]string{
      "search": "kitchen papers",
      "size": "large",
     })
    

#### func (*Client) [SetRateLimiter](https://github.com/go-resty/resty/blob/v2.17.1/client.go#L947) ¶ added in v2.9.0

    func (c *Client) SetRateLimiter(rl RateLimiter) *Client

SetRateLimiter sets an optional RateLimiter. If set, the rate limiter will control all requests were made by this client.

#### func (*Client) [SetRawPathParam](https://github.com/go-resty/resty/blob/v2.17.1/client.go#L1067) ¶ added in v2.8.0

    func (c *Client) SetRawPathParam(param, value [string](/builtin#string)) *Client

SetRawPathParam method sets a single URL path key-value pair in the Resty client instance.

    client.SetPathParam("userId", "sample@sample.com")
    
    Result:
       URL - /v1/users/{userId}/details
       Composed URL - /v1/users/sample@sample.com/details
    
    client.SetPathParam("path", "groups/developers")
    
    Result:
       URL - /v1/users/{userId}/details
       Composed URL - /v1/users/groups%2Fdevelopers/details
    

It replaces the value of the key while composing the request URL. The value will be used as it is and will not be escaped.

It can be overridden at the request level, see Request.SetRawPathParam or Request.SetRawPathParams

#### func (*Client) [SetRawPathParams](https://github.com/go-resty/resty/blob/v2.17.1/client.go#L1090) ¶ added in v2.8.0

    func (c *Client) SetRawPathParams(params map[[string](/builtin#string)][string](/builtin#string)) *Client

SetRawPathParams method sets multiple URL path key-value pairs at one go in the Resty client instance.

    client.SetPathParams(map[string]string{
     "userId":       "sample@sample.com",
     "subAccountId": "100002",
     "path":         "groups/developers",
    })
    
    Result:
       URL - /v1/users/{userId}/{subAccountId}/{path}/details
       Composed URL - /v1/users/sample@sample.com/100002/groups/developers/details
    

It replaces the value of the key while composing the request URL. The values will be used as they are and will not be escaped.

It can be overridden at the request level, see Request.SetRawPathParam or Request.SetRawPathParams

#### func (*Client) [SetRedirectPolicy](https://github.com/go-resty/resty/blob/v2.17.1/client.go#L686) ¶

    func (c *Client) SetRedirectPolicy(policies ...interface{}) *Client

SetRedirectPolicy method sets the redirect policy for the client. Resty provides ready-to-use redirect policies. Wanna create one for yourself, refer to `redirect.go`.

    client.SetRedirectPolicy(FlexibleRedirectPolicy(20))
    
    // Need multiple redirect policies together
    client.SetRedirectPolicy(FlexibleRedirectPolicy(20), DomainCheckRedirectPolicy("host1.com", "host2.net"))
    

#### func (*Client) [SetResponseBodyLimit](https://github.com/go-resty/resty/blob/v2.17.1/client.go#L1119) ¶ added in v2.15.0

    func (c *Client) SetResponseBodyLimit(v [int](/builtin#int)) *Client

SetResponseBodyLimit method sets a maximum body size limit in bytes on response, avoid reading too much data to memory.

Client will return resty.ErrResponseBodyTooLarge if the body size of the body in the uncompressed response is larger than the limit. Body size limit will not be enforced in the following cases:

- ResponseBodyLimit <= 0, which is the default behavior.
- Request.SetOutput is called to save response data to the file.
- "DoNotParseResponse" is set for client or request.

It can be overridden at the request level; see Request.SetResponseBodyLimit

#### func (*Client) [SetRetryAfter](https://github.com/go-resty/resty/blob/v2.17.1/client.go#L733) ¶

    func (c *Client) SetRetryAfter(callback RetryAfterFunc) *Client

SetRetryAfter sets a callback to calculate the wait time between retries. Default (nil) implies exponential backoff with jitter

#### func (*Client) [SetRetryCount](https://github.com/go-resty/resty/blob/v2.17.1/client.go#L708) ¶

    func (c *Client) SetRetryCount(count [int](/builtin#int)) *Client

SetRetryCount method enables retry on Resty client and allows you to set no. of retry count. Resty uses a Backoff mechanism.

#### func (*Client) [SetRetryMaxWaitTime](https://github.com/go-resty/resty/blob/v2.17.1/client.go#L726) ¶

    func (c *Client) SetRetryMaxWaitTime(maxWaitTime [time](/time).[Duration](/time#Duration)) *Client

SetRetryMaxWaitTime method sets the max wait time for sleep before retrying request.

Default is 2 seconds.

#### func (*Client) [SetRetryResetReaders](https://github.com/go-resty/resty/blob/v2.17.1/client.go#L795) ¶ added in v2.8.0

    func (c *Client) SetRetryResetReaders(b [bool](/builtin#bool)) *Client

SetRetryResetReaders method enables the Resty client to seek the start of all file readers are given as multipart files if the object implements [io.ReadSeeker](/io#ReadSeeker).

#### func (*Client) [SetRetryWaitTime](https://github.com/go-resty/resty/blob/v2.17.1/client.go#L717) ¶

    func (c *Client) SetRetryWaitTime(waitTime [time](/time).[Duration](/time#Duration)) *Client

SetRetryWaitTime method sets the default wait time for sleep before retrying request.

Default is 100 milliseconds.

#### func (*Client) [SetRootCertificate](https://github.com/go-resty/resty/blob/v2.17.1/client.go#L872) ¶

    func (c *Client) SetRootCertificate(pemFilePath [string](/builtin#string)) *Client

SetRootCertificate method helps to add one or more root certificates into the Resty client

    client.SetRootCertificate("/path/to/root/pemFile.pem")
    

#### func (*Client) [SetRootCertificateFromString](https://github.com/go-resty/resty/blob/v2.17.1/client.go#L886) ¶ added in v2.2.0

    func (c *Client) SetRootCertificateFromString(pemCerts [string](/builtin#string)) *Client

SetRootCertificateFromString method helps to add one or more root certificates into the Resty client

    client.SetRootCertificateFromString("pem certs content")
    

#### func (*Client) [SetScheme](https://github.com/go-resty/resty/blob/v2.17.1/client.go#L977) ¶

    func (c *Client) SetScheme(scheme [string](/builtin#string)) *Client

SetScheme method sets a custom scheme for the Resty client. It's a way to override the default.

    client.SetScheme("http")
    

#### func (*Client) [SetTLSClientConfig](https://github.com/go-resty/resty/blob/v2.17.1/client.go#L811) ¶

    func (c *Client) SetTLSClientConfig(config *[tls](/crypto/tls).[Config](/crypto/tls#Config)) *Client

SetTLSClientConfig method sets TLSClientConfig for underlying client Transport.

For Example:

    // One can set a custom root certificate. Refer: http://golang.org/pkg/crypto/tls/#example_Dial
    client.SetTLSClientConfig(&tls.Config{ RootCAs: roots })
    
    // or One can disable security check (https)
    client.SetTLSClientConfig(&tls.Config{ InsecureSkipVerify: true })
    

NOTE: This method overwrites existing [http.Transport.TLSClientConfig](/net/http#Transport.TLSClientConfig)

#### func (*Client) [SetTimeout](https://github.com/go-resty/resty/blob/v2.17.1/client.go#L662) ¶

    func (c *Client) SetTimeout(timeout [time](/time).[Duration](/time#Duration)) *Client

SetTimeout method sets the timeout for a request raised by the client.

    client.SetTimeout(time.Duration(1 * time.Minute))
    

#### func (*Client) [SetTransport](https://github.com/go-resty/resty/blob/v2.17.1/client.go#L967) ¶

    func (c *Client) SetTransport(transport [http](/net/http).[RoundTripper](/net/http#RoundTripper)) *Client

SetTransport method sets custom [http.Transport](/net/http#Transport) or any [http.RoundTripper](/net/http#RoundTripper) compatible interface implementation in the Resty client.

    transport := &http.Transport{
     // something like Proxying to httptest.Server, etc...
     Proxy: func(req *http.Request) (*url.URL, error) {
      return url.Parse(server.URL)
     },
    }
    client.SetTransport(transport)
    

NOTE:

- If transport is not the type of `*http.Transport`, then you may not be able to take advantage of some of the Resty client settings.
- It overwrites the Resty client transport instance and its configurations.

#### func (*Client) [SetUnescapeQueryParams](https://github.com/go-resty/resty/blob/v2.17.1/client.go#L335) ¶ added in v2.16.0

    func (c *Client) SetUnescapeQueryParams(unescape [bool](/builtin#bool)) *Client

SetUnescapeQueryParams method sets the unescape query parameters choice for request URL. To prevent broken URL, resty replaces space (" ") with "+" in the query parameters.

See Request.SetUnescapeQueryParams

NOTE: Request failure is possible due to non-standard usage of Unescaped Query Parameters.

#### func (*Client) [SetXMLMarshaler](https://github.com/go-resty/resty/blob/v2.17.1/client.go#L754) ¶ added in v2.8.0

    func (c *Client) SetXMLMarshaler(marshaler func(v interface{}) ([][byte](/builtin#byte), [error](/builtin#error))) *Client

SetXMLMarshaler method sets the XML marshaler function to marshal the request body. By default, Resty uses [encoding/xml](/encoding/xml) package to marshal the request body.

#### func (*Client) [SetXMLUnmarshaler](https://github.com/go-resty/resty/blob/v2.17.1/client.go#L761) ¶ added in v2.8.0

    func (c *Client) SetXMLUnmarshaler(unmarshaler func(data [][byte](/builtin#byte), v interface{}) [error](/builtin#error)) *Client

SetXMLUnmarshaler method sets the XML unmarshaler function to unmarshal the response body. By default, Resty uses [encoding/xml](/encoding/xml) package to unmarshal the response body.

#### func (*Client) [Transport](https://github.com/go-resty/resty/blob/v2.17.1/client.go#L1361) ¶ added in v2.8.0

    func (c *Client) Transport() (*[http](/net/http).[Transport](/net/http#Transport), [error](/builtin#error))

Transport method returns [http.Transport](/net/http#Transport) currently in use or error in case the currently used `transport` is not a [http.Transport](/net/http#Transport).

Since v2.8.0 has become exported method.

#### type [ErrorHook](https://github.com/go-resty/resty/blob/v2.17.1/client.go#L89) ¶ added in v2.4.0

    type ErrorHook func(*Request, [error](/builtin#error))

ErrorHook type is for reacting to request errors, called after all retries were attempted

#### type [File](https://github.com/go-resty/resty/blob/v2.17.1/client.go#L1426) ¶

    type File struct {
     Name      [string](/builtin#string)
     ParamName [string](/builtin#string)
     [io](/io).[Reader](/io#Reader)
    }

File struct represents file information for multipart request

#### func (*File) [String](https://github.com/go-resty/resty/blob/v2.17.1/client.go#L1433) ¶

    func (f *File) String() [string](/builtin#string)

String method returns the string value of current file details

#### type [Logger](https://github.com/go-resty/resty/blob/v2.17.1/util.go#L30) ¶

    type Logger interface {
     Errorf(format [string](/builtin#string), v ...interface{})
     Warnf(format [string](/builtin#string), v ...interface{})
     Debugf(format [string](/builtin#string), v ...interface{})
    }

Logger interface is to abstract the logging from Resty. Gives control to the Resty users, choice of the logger.

#### type [MultipartField](https://github.com/go-resty/resty/blob/v2.17.1/client.go#L1442) ¶

    type MultipartField struct {
     Param       [string](/builtin#string)
     FileName    [string](/builtin#string)
     ContentType [string](/builtin#string)
     [io](/io).[Reader](/io#Reader)
    }

MultipartField struct represents the custom data part for a multipart request

#### type [OnRetryFunc](https://github.com/go-resty/resty/blob/v2.17.1/retry.go#L31) ¶ added in v2.6.0

    type OnRetryFunc func(*Response, [error](/builtin#error))

OnRetryFunc is for side-effecting functions triggered on retry

#### type [Option](https://github.com/go-resty/resty/blob/v2.17.1/retry.go#L24) ¶

    type Option func(*Options)

Option is to create convenient retry options like wait time, max retries, etc.

#### func [MaxWaitTime](https://github.com/go-resty/resty/blob/v2.17.1/retry.go#L66) ¶

    func MaxWaitTime(value [time](/time).[Duration](/time#Duration)) Option

MaxWaitTime sets the max wait time to sleep between requests

#### func [ResetMultipartReaders](https://github.com/go-resty/resty/blob/v2.17.1/retry.go#L88) ¶ added in v2.8.0

    func ResetMultipartReaders(value [bool](/builtin#bool)) Option

ResetMultipartReaders sets a boolean value which will lead the start being seeked out on all multipart file readers if they implement [io.ReadSeeker](/io#ReadSeeker)

#### func [Retries](https://github.com/go-resty/resty/blob/v2.17.1/retry.go#L52) ¶

    func Retries(value [int](/builtin#int)) Option

Retries sets the max number of retries

#### func [RetryConditions](https://github.com/go-resty/resty/blob/v2.17.1/retry.go#L73) ¶

    func RetryConditions(conditions []RetryConditionFunc) Option

RetryConditions sets the conditions that will be checked for retry

#### func [RetryHooks](https://github.com/go-resty/resty/blob/v2.17.1/retry.go#L80) ¶ added in v2.6.0

    func RetryHooks(hooks []OnRetryFunc) Option

RetryHooks sets the hooks that will be executed after each retry

#### func [WaitTime](https://github.com/go-resty/resty/blob/v2.17.1/retry.go#L59) ¶

    func WaitTime(value [time](/time).[Duration](/time#Duration)) Option

WaitTime sets the default wait time to sleep between requests

#### type [Options](https://github.com/go-resty/resty/blob/v2.17.1/retry.go#L41) ¶

    type Options struct {
     // contains filtered or unexported fields
    }

Options struct is used to hold retry settings.

#### type [PreRequestHook](https://github.com/go-resty/resty/blob/v2.17.1/client.go#L80) ¶ added in v2.3.0

    type PreRequestHook func(*Client, *[http](/net/http).[Request](/net/http#Request)) [error](/builtin#error)

PreRequestHook type is for the request hook, called right before the request is sent

#### type [RateLimiter](https://github.com/go-resty/resty/blob/v2.17.1/util.go#L71) ¶ added in v2.9.0

    type RateLimiter interface {
     Allow() [bool](/builtin#bool)
    }

#### type [RedirectPolicy](https://github.com/go-resty/resty/blob/v2.17.1/redirect.go#L25) ¶

    type RedirectPolicy interface {
     Apply(req *[http](/net/http).[Request](/net/http#Request), via []*[http](/net/http).[Request](/net/http#Request)) [error](/builtin#error)
    }

RedirectPolicy to regulate the redirects in the Resty client. Objects implementing the RedirectPolicy interface can be registered as

Apply function should return nil to continue the redirect journey; otherwise return error to stop the redirect.

#### func [DomainCheckRedirectPolicy](https://github.com/go-resty/resty/blob/v2.17.1/redirect.go#L66) ¶

    func DomainCheckRedirectPolicy(hostnames ...[string](/builtin#string)) RedirectPolicy

DomainCheckRedirectPolicy method is convenient for defining domain name redirect rules in Resty clients. Redirect is allowed only for the host mentioned in the policy.

    resty.SetRedirectPolicy(DomainCheckRedirectPolicy("host1.com", "host2.org", "host3.net"))
    

#### func [FlexibleRedirectPolicy](https://github.com/go-resty/resty/blob/v2.17.1/redirect.go#L52) ¶

    func FlexibleRedirectPolicy(noOfRedirect [int](/builtin#int)) RedirectPolicy

FlexibleRedirectPolicy method is convenient for creating several redirect policies for Resty clients.

    resty.SetRedirectPolicy(FlexibleRedirectPolicy(20))
    

#### func [NoRedirectPolicy](https://github.com/go-resty/resty/blob/v2.17.1/redirect.go#L43) ¶

    func NoRedirectPolicy() RedirectPolicy

NoRedirectPolicy is used to disable redirects in the Resty client

    resty.SetRedirectPolicy(NoRedirectPolicy())
    

#### type [RedirectPolicyFunc](https://github.com/go-resty/resty/blob/v2.17.1/redirect.go#L32) ¶

    type RedirectPolicyFunc func(*[http](/net/http).[Request](/net/http#Request), []*[http](/net/http).[Request](/net/http#Request)) [error](/builtin#error)

The RedirectPolicyFunc type is an adapter to allow the use of ordinary functions as RedirectPolicy. If `f` is a function with the appropriate signature, RedirectPolicyFunc(f) is a RedirectPolicy object that calls `f`.

#### func (RedirectPolicyFunc) [Apply](https://github.com/go-resty/resty/blob/v2.17.1/redirect.go#L36) ¶

    func (f RedirectPolicyFunc) Apply(req *[http](/net/http).[Request](/net/http#Request), via []*[http](/net/http).[Request](/net/http#Request)) [error](/builtin#error)

Apply calls f(req, via).

#### type [Request](https://github.com/go-resty/resty/blob/v2.17.1/request.go#L29) ¶

    type Request struct {
     URL           [string](/builtin#string)
     Method        [string](/builtin#string)
     Token         [string](/builtin#string)
     AuthScheme    [string](/builtin#string)
     QueryParam    [url](/net/url).[Values](/net/url#Values)
     FormData      [url](/net/url).[Values](/net/url#Values)
     PathParams    map[[string](/builtin#string)][string](/builtin#string)
     RawPathParams map[[string](/builtin#string)][string](/builtin#string)
     Header        [http](/net/http).[Header](/net/http#Header)
     Time          [time](/time).[Time](/time#Time)
     Body          interface{}
     Result        interface{}
    
     Error      interface{}
     RawRequest *[http](/net/http).[Request](/net/http#Request)
     SRV        *SRVRecord
     UserInfo   *User
     Cookies    []*[http](/net/http).[Cookie](/net/http#Cookie)
     Debug      [bool](/builtin#bool)
    
     // Attempt is to represent the request attempt made during a Resty
     // request execution flow, including retry count.
     Attempt [int](/builtin#int)
     // contains filtered or unexported fields
    }

Request struct is used to compose and fire individual requests from Resty client. The Request provides an option to override client-level settings and also an option for the request composition.

#### func (*Request) [AddRetryCondition](https://github.com/go-resty/resty/blob/v2.17.1/request.go#L848) ¶ added in v2.7.0

    func (r *Request) AddRetryCondition(condition RetryConditionFunc) *Request

AddRetryCondition method adds a retry condition function to the request's array of functions is checked to determine if the request can be retried. The request will retry if any functions return true and the error is nil.

NOTE: The request level retry conditions are checked before all retry conditions from the client instance.

#### func (*Request) [Context](https://github.com/go-resty/resty/blob/v2.17.1/request.go#L101) ¶

    func (r *Request) Context() [context](/context).[Context](/context#Context)

Context method returns the Context if it is already set in the Request otherwise, it creates a new one using [context.Background](/context#Background).

#### func (*Request) [Delete](https://github.com/go-resty/resty/blob/v2.17.1/request.go#L984) ¶

    func (r *Request) Delete(url [string](/builtin#string)) (*Response, [error](/builtin#error))

Delete method does DELETE HTTP request. It's defined in section 4.3.5 of RFC7231.

#### func (*Request) [DisableGenerateCurlOnDebug](https://github.com/go-resty/resty/blob/v2.17.1/request.go#L885) ¶ added in v2.15.0

    func (r *Request) DisableGenerateCurlOnDebug() *Request

DisableGenerateCurlOnDebug method disables the option set by Request.EnableGenerateCurlOnDebug. It overrides the options set by the Client.

#### func (*Request) [EnableGenerateCurlOnDebug](https://github.com/go-resty/resty/blob/v2.17.1/request.go#L878) ¶ added in v2.15.0

    func (r *Request) EnableGenerateCurlOnDebug() *Request

EnableGenerateCurlOnDebug method enables the generation of CURL commands in the debug log. It works in conjunction with debug mode. It overrides the options set by the Client.

NOTE: Use with care.

- Potential to leak sensitive data from Request and Response in the debug log.
- Beware of memory usage since the request body is reread.

#### func (*Request) [EnableTrace](https://github.com/go-resty/resty/blob/v2.17.1/request.go#L867) ¶

    func (r *Request) EnableTrace() *Request

EnableTrace method enables trace for the current request using [httptrace.ClientTrace](/net/http/httptrace#ClientTrace) and provides insights.

    client := resty.New()
    
    resp, err := client.R().EnableTrace().Get("https://httpbin.org/get")
    fmt.Println("Error:", err)
    fmt.Println("Trace Info:", resp.Request.TraceInfo())
    

See Client.EnableTrace is also available to get trace info for all requests.

#### func (*Request) [Execute](https://github.com/go-resty/resty/blob/v2.17.1/request.go#L1013) ¶

    func (r *Request) Execute(method, url [string](/builtin#string)) (*Response, [error](/builtin#error))

Execute method performs the HTTP request with the given HTTP method and URL for current Request.

    resp, err := client.R().Execute(resty.MethodGet, "http://httpbin.org/get")
    

#### func (*Request) [ExpectContentType](https://github.com/go-resty/resty/blob/v2.17.1/request.go#L754) ¶

    func (r *Request) ExpectContentType(contentType [string](/builtin#string)) *Request

ExpectContentType method allows to provide fallback `Content-Type` for automatic unmarshalling when the `Content-Type` response header is unavailable.

#### func (*Request) [ForceContentType](https://github.com/go-resty/resty/blob/v2.17.1/request.go#L765) ¶ added in v2.3.0

    func (r *Request) ForceContentType(contentType [string](/builtin#string)) *Request

ForceContentType method provides a strong sense of response `Content-Type` for automatic unmarshalling. Resty gives this a higher priority than the `Content-Type` response header.

This means that if both Request.ForceContentType is set and the response `Content-Type` is available, `ForceContentType` will win.

#### func (*Request) [GenerateCurlCommand](https://github.com/go-resty/resty/blob/v2.17.1/request.go#L80) ¶ added in v2.14.0

    func (r *Request) GenerateCurlCommand() [string](/builtin#string)

GenerateCurlCommand method generates the CURL command for the request.

#### func (*Request) [Get](https://github.com/go-resty/resty/blob/v2.17.1/request.go#L964) ¶

    func (r *Request) Get(url [string](/builtin#string)) (*Response, [error](/builtin#error))

Get method does GET HTTP request. It's defined in section 4.3.1 of RFC7231.

#### func (*Request) [Head](https://github.com/go-resty/resty/blob/v2.17.1/request.go#L969) ¶

    func (r *Request) Head(url [string](/builtin#string)) (*Response, [error](/builtin#error))

Head method does HEAD HTTP request. It's defined in section 4.3.2 of RFC7231.

#### func (*Request) [Options](https://github.com/go-resty/resty/blob/v2.17.1/request.go#L989) ¶

    func (r *Request) Options(url [string](/builtin#string)) (*Response, [error](/builtin#error))

Options method does OPTIONS HTTP request. It's defined in section 4.3.7 of RFC7231.

#### func (*Request) [Patch](https://github.com/go-resty/resty/blob/v2.17.1/request.go#L994) ¶

    func (r *Request) Patch(url [string](/builtin#string)) (*Response, [error](/builtin#error))

Patch method does PATCH HTTP request. It's defined in section 2 of RFC5789.

#### func (*Request) [Post](https://github.com/go-resty/resty/blob/v2.17.1/request.go#L974) ¶

    func (r *Request) Post(url [string](/builtin#string)) (*Response, [error](/builtin#error))

Post method does POST HTTP request. It's defined in section 4.3.3 of RFC7231.

#### func (*Request) [Put](https://github.com/go-resty/resty/blob/v2.17.1/request.go#L979) ¶

    func (r *Request) Put(url [string](/builtin#string)) (*Response, [error](/builtin#error))

Put method does PUT HTTP request. It's defined in section 4.3.4 of RFC7231.

#### func (*Request) [Send](https://github.com/go-resty/resty/blob/v2.17.1/request.go#L1005) ¶ added in v2.2.0

    func (r *Request) Send() (*Response, [error](/builtin#error))

Send method performs the HTTP request using the method and URL already defined for current Request.

    req := client.R()
    req.Method = resty.MethodGet
    req.URL = "http://httpbin.org/get"
    resp, err := req.Send()
    

#### func (*Request) [SetAuthScheme](https://github.com/go-resty/resty/blob/v2.17.1/request.go#L563) ¶ added in v2.3.0

    func (r *Request) SetAuthScheme(scheme [string](/builtin#string)) *Request

SetAuthScheme method sets the auth token scheme type in the HTTP request.

Example Header value structure:

    Authorization: <auth-scheme-value-set-here> <auth-token-value>
    

For Example: To set the scheme to use OAuth

    client.R().SetAuthScheme("OAuth")
    
    // The outcome will be -
    Authorization: OAuth <auth-token-value>
    

Information about Auth schemes can be found in [](https://tools.ietf.org/html/rfc7235)[RFC 7235](https://rfc-editor.org/rfc/rfc7235.html), IANA [HTTP Auth schemes](https://www.iana.org/assignments/http-authschemes/http-authschemes.xhtml#authschemes)

It overrides the `Authorization` scheme set by method Client.SetAuthScheme.

#### func (*Request) [SetAuthToken](https://github.com/go-resty/resty/blob/v2.17.1/request.go#L539) ¶

    func (r *Request) SetAuthToken(token [string](/builtin#string)) *Request

SetAuthToken method sets the auth token header(Default Scheme: Bearer) in the current HTTP request. Header example:

    Authorization: Bearer <auth-token-value-comes-here>
    

For Example: To set auth token BC594900518B4F7EAC75BD37F019E08FBC594900518B4F7EAC75BD37F019E08F

    client.R().SetAuthToken("BC594900518B4F7EAC75BD37F019E08FBC594900518B4F7EAC75BD37F019E08F")
    

It overrides the Auth token set by method Client.SetAuthToken.

#### func (*Request) [SetBasicAuth](https://github.com/go-resty/resty/blob/v2.17.1/request.go#L525) ¶

    func (r *Request) SetBasicAuth(username, password [string](/builtin#string)) *Request

SetBasicAuth method sets the basic authentication header in the current HTTP request.

For Example:

    Authorization: Basic <base64-encoded-value>
    

To set the header for username "go-resty" and password "welcome"

    client.R().SetBasicAuth("go-resty", "welcome")
    

It overrides the credentials set by method Client.SetBasicAuth.

#### func (*Request) [SetBody](https://github.com/go-resty/resty/blob/v2.17.1/request.go#L349) ¶

    func (r *Request) SetBody(body interface{}) *Request

SetBody method sets the request body for the request. It supports various practical needs as easy. It's quite handy and powerful. Supported request body data types are `string`, `[]byte`, `struct`, `map`, `slice` and [io.Reader](/io#Reader).

Body value can be pointer or non-pointer. Automatic marshalling for JSON and XML content type, if it is `struct`, `map`, or `slice`.

NOTE: [io.Reader](/io#Reader) is processed in bufferless mode while sending a request.

For Example:

`struct` gets marshaled based on the request header `Content-Type`.

    client.R().
     SetBody(User{
      Username: "jeeva@myjeeva.com",
      Password: "welcome2resty",
     })
    

'map` gets marshaled based on the request header `Content-Type`.

    client.R().
     SetBody(map[string]interface{}{
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

#### func (*Request) [SetContentLength](https://github.com/go-resty/resty/blob/v2.17.1/request.go#L509) ¶

    func (r *Request) SetContentLength(l [bool](/builtin#bool)) *Request

SetContentLength method sets the current request's HTTP header `Content-Length` value. By default, Resty won't set `Content-Length`.

See Client.SetContentLength

    client.R().SetContentLength(true)
    

It overrides the value set at the client instance level.

#### func (*Request) [SetContext](https://github.com/go-resty/resty/blob/v2.17.1/request.go#L112) ¶

    func (r *Request) SetContext(ctx [context](/context).[Context](/context#Context)) *Request

SetContext method sets the [context.Context](/context#Context) for current Request. It allows to interrupt the request execution if `ctx.Done()` channel is closed. See <https://blog.golang.org/context> article and the package [context](/context) documentation.

#### func (*Request) [SetCookie](https://github.com/go-resty/resty/blob/v2.17.1/request.go#L789) ¶ added in v2.1.0

    func (r *Request) SetCookie(hc *[http](/net/http).[Cookie](/net/http#Cookie)) *Request

SetCookie method appends a single cookie in the current request instance.

    client.R().SetCookie(&http.Cookie{
       Name:"go-resty",
       Value:"This is cookie value",
      })
    

NOTE: Method appends the Cookie value into existing Cookie even if its already existing.

#### func (*Request) [SetCookies](https://github.com/go-resty/resty/blob/v2.17.1/request.go#L811) ¶ added in v2.1.0

    func (r *Request) SetCookies(rs []*[http](/net/http).[Cookie](/net/http#Cookie)) *Request

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

#### func (*Request) [SetDebug](https://github.com/go-resty/resty/blob/v2.17.1/request.go#L837) ¶ added in v2.8.0

    func (r *Request) SetDebug(d [bool](/builtin#bool)) *Request

SetDebug method enables the debug mode on the current request. It logs the details current request and response.

    client.SetDebug(true)
    

Also, it can be enabled at the request level for a particular request; see Request.SetDebug.

- For Request, it logs information such as HTTP verb, Relative URL path, Host, Headers, and Body if it has one.
- For Response, it logs information such as Status, Response Time, Headers, and Body if it has one.

#### func (*Request) [SetDigestAuth](https://github.com/go-resty/resty/blob/v2.17.1/request.go#L581) ¶ added in v2.8.0

    func (r *Request) SetDigestAuth(username, password [string](/builtin#string)) *Request

SetDigestAuth method sets the Digest Access auth scheme for the HTTP request. If a server responds with 401 and sends a Digest challenge in the WWW-Authenticate Header, the request will be resent with the appropriate Authorization Header.

For Example: To set the Digest scheme with username "Mufasa" and password "Circle Of Life"

    client.R().SetDigestAuth("Mufasa", "Circle Of Life")
    

Information about Digest Access Authentication can be found in [](https://datatracker.ietf.org/doc/html/rfc7616)[RFC 7616](https://rfc-editor.org/rfc/rfc7616.html)

It overrides the digest username and password set by method Client.SetDigestAuth.

#### func (*Request) [SetDoNotParseResponse](https://github.com/go-resty/resty/blob/v2.17.1/request.go#L635) ¶

    func (r *Request) SetDoNotParseResponse(parse [bool](/builtin#bool)) *Request

SetDoNotParseResponse method instructs Resty not to parse the response body automatically. Resty exposes the raw response body as [io.ReadCloser](/io#ReadCloser). If you use it, do not forget to close the body, otherwise, you might get into connection leaks, and connection reuse may not happen.

NOTE: Response middlewares are not executed using this option. You have taken over the control of response parsing from Resty.

#### func (*Request) [SetError](https://github.com/go-resty/resty/blob/v2.17.1/request.go#L399) ¶

    func (r *Request) SetError(err interface{}) *Request

SetError method is to register the request `Error` object for automatic unmarshalling for the request, if the response status code is greater than 399 and the content type is either JSON or XML.

NOTE: Request.SetError input can be a pointer or non-pointer.

    client.R().SetError(&AuthError{})
    // OR
    client.R().SetError(AuthError{})
    

Accessing an error value from response instance.

    response.Error().(*AuthError)
    

If this request Error object is nil, Resty will use the client-level error object Type if it is set.

#### func (*Request) [SetFile](https://github.com/go-resty/resty/blob/v2.17.1/request.go#L408) ¶

    func (r *Request) SetFile(param, filePath [string](/builtin#string)) *Request

SetFile method sets a single file field name and its path for multipart upload.

    client.R().
     SetFile("my_file", "/Users/jeeva/Gas Bill - Sep.pdf")
    

#### func (*Request) [SetFileReader](https://github.com/go-resty/resty/blob/v2.17.1/request.go#L435) ¶

    func (r *Request) SetFileReader(param, fileName [string](/builtin#string), reader [io](/io).[Reader](/io#Reader)) *Request

SetFileReader method is to set a file using [io.Reader](/io#Reader) for multipart upload.

    client.R().
     SetFileReader("profile_img", "my-profile-img.png", bytes.NewReader(profileImgBytes)).
     SetFileReader("notes", "user-notes.txt", bytes.NewReader(notesBytes))
    

#### func (*Request) [SetFiles](https://github.com/go-resty/resty/blob/v2.17.1/request.go#L422) ¶

    func (r *Request) SetFiles(files map[[string](/builtin#string)][string](/builtin#string)) *Request

SetFiles method sets multiple file field names and their paths for multipart uploads.

    client.R().
     SetFiles(map[string]string{
       "my_file1": "/Users/jeeva/Gas Bill - Sep.pdf",
       "my_file2": "/Users/jeeva/Electricity Bill - Sep.pdf",
       "my_file3": "/Users/jeeva/Water Bill - Sep.pdf",
      })
    

#### func (*Request) [SetFormData](https://github.com/go-resty/resty/blob/v2.17.1/request.go#L277) ¶

    func (r *Request) SetFormData(data map[[string](/builtin#string)][string](/builtin#string)) *Request

SetFormData method sets Form parameters and their values for the current request. It applies only to HTTP methods `POST` and `PUT`, and by default requests content type would be set as `application/x-www-form-urlencoded`.

    client.R().
     SetFormData(map[string]string{
      "access_token": "BC594900-518B-4F7E-AC75-BD37F019E08F",
      "user_id": "3455454545",
     })
    

It overrides the form data value set at the client instance level.

#### func (*Request) [SetFormDataFromValues](https://github.com/go-resty/resty/blob/v2.17.1/request.go#L293) ¶

    func (r *Request) SetFormDataFromValues(data [url](/net/url).[Values](/net/url#Values)) *Request

SetFormDataFromValues method appends multiple form parameters with multi-value ([url.Values](/net/url#Values)) at one go in the current request.

    client.R().
     SetFormDataFromValues(url.Values{
      "search_criteria": []string{"book", "glass", "pencil"},
     })
    

It overrides the form data value set at the client instance level.

#### func (*Request) [SetHeader](https://github.com/go-resty/resty/blob/v2.17.1/request.go#L126) ¶

    func (r *Request) SetHeader(header, value [string](/builtin#string)) *Request

SetHeader method sets a single header field and its value in the current request.

For Example: To set `Content-Type` and `Accept` as `application/json`.

    client.R().
     SetHeader("Content-Type", "application/json").
     SetHeader("Accept", "application/json")
    

It overrides the header value set at the client instance level.

#### func (*Request) [SetHeaderMultiValues](https://github.com/go-resty/resty/blob/v2.17.1/request.go#L159) ¶ added in v2.7.0

    func (r *Request) SetHeaderMultiValues(headers map[[string](/builtin#string)][][string](/builtin#string)) *Request

SetHeaderMultiValues sets multiple header fields and their values as a list of strings in the current request.

For Example: To set `Accept` as `text/html, application/xhtml+xml, application/xml;q=0.9, image/webp, */*;q=0.8`

    client.R().
     SetHeaderMultiValues(map[string][]string{
      "Accept": []string{"text/html", "application/xhtml+xml", "application/xml;q=0.9", "image/webp", "*/*;q=0.8"},
     })
    

It overrides the header value set at the client instance level.

#### func (*Request) [SetHeaderVerbatim](https://github.com/go-resty/resty/blob/v2.17.1/request.go#L175) ¶ added in v2.6.0

    func (r *Request) SetHeaderVerbatim(header, value [string](/builtin#string)) *Request

SetHeaderVerbatim method sets a single header field and its value verbatim in the current request.

For Example: To set `all_lowercase` and `UPPERCASE` as `available`.

    client.R().
     SetHeaderVerbatim("all_lowercase", "available").
     SetHeaderVerbatim("UPPERCASE", "available")
    

It overrides the header value set at the client instance level.

#### func (*Request) [SetHeaders](https://github.com/go-resty/resty/blob/v2.17.1/request.go#L142) ¶

    func (r *Request) SetHeaders(headers map[[string](/builtin#string)][string](/builtin#string)) *Request

SetHeaders method sets multiple header fields and their values at one go in the current request.

For Example: To set `Content-Type` and `Accept` as `application/json`

    client.R().
     SetHeaders(map[string]string{
      "Content-Type": "application/json",
      "Accept": "application/json",
     })
    

It overrides the header value set at the client instance level.

#### func (*Request) [SetJSONEscapeHTML](https://github.com/go-resty/resty/blob/v2.17.1/request.go#L776) ¶

    func (r *Request) SetJSONEscapeHTML(b [bool](/builtin#bool)) *Request

SetJSONEscapeHTML method enables or disables the HTML escape on JSON marshal. By default, escape HTML is false.

NOTE: This option only applies to the standard JSON Marshaller used by Resty.

It overrides the value set at the client instance level, see Client.SetJSONEscapeHTML

#### func (*Request) [SetLogger](https://github.com/go-resty/resty/blob/v2.17.1/request.go#L822) ¶ added in v2.8.0

    func (r *Request) SetLogger(l Logger) *Request

SetLogger method sets given writer for logging Resty request and response details. By default, requests and responses inherit their logger from the client.

Compliant to interface resty.Logger.

It overrides the logger value set at the client instance level.

#### func (*Request) [SetMultipartBoundary](https://github.com/go-resty/resty/blob/v2.17.1/request.go#L496) ¶ added in v2.15.0

    func (r *Request) SetMultipartBoundary(boundary [string](/builtin#string)) *Request

SetMultipartBoundary method sets the custom multipart boundary for the multipart request. Typically, the `mime/multipart` package generates a random multipart boundary if not provided.

#### func (*Request) [SetMultipartField](https://github.com/go-resty/resty/blob/v2.17.1/request.go#L456) ¶

    func (r *Request) SetMultipartField(param, fileName, contentType [string](/builtin#string), reader [io](/io).[Reader](/io#Reader)) *Request

SetMultipartField method sets custom data with Content-Type using [io.Reader](/io#Reader) for multipart upload.

#### func (*Request) [SetMultipartFields](https://github.com/go-resty/resty/blob/v2.17.1/request.go#L488) ¶

    func (r *Request) SetMultipartFields(fields ...*MultipartField) *Request

SetMultipartFields method sets multiple data fields using [io.Reader](/io#Reader) for multipart upload.

For Example:

    client.R().SetMultipartFields(
     &resty.MultipartField{
      Param:       "uploadManifest1",
      FileName:    "upload-file-1.json",
      ContentType: "application/json",
      Reader:      strings.NewReader(`{"input": {"name": "Uploaded document 1", "_filename" : ["file1.txt"]}}`),
     },
     &resty.MultipartField{
      Param:       "uploadManifest2",
      FileName:    "upload-file-2.json",
      ContentType: "application/json",
      Reader:      strings.NewReader(`{"input": {"name": "Uploaded document 2", "_filename" : ["file2.txt"]}}`),
     })
    

If you have a `slice` of fields already, then call-

    client.R().SetMultipartFields(fields...)
    

#### func (*Request) [SetMultipartFormData](https://github.com/go-resty/resty/blob/v2.17.1/request.go#L447) ¶ added in v2.3.0

    func (r *Request) SetMultipartFormData(data map[[string](/builtin#string)][string](/builtin#string)) *Request

SetMultipartFormData method allows simple form data to be attached to the request as `multipart:form-data`

#### func (*Request) [SetOutput](https://github.com/go-resty/resty/blob/v2.17.1/request.go#L611) ¶

    func (r *Request) SetOutput(file [string](/builtin#string)) *Request

SetOutput method sets the output file for the current HTTP request. The current HTTP response will be saved in the given file. It is similar to the `curl -o` flag.

Absolute path or relative path can be used.

If it is a relative path, then the output file goes under the output directory, as mentioned in the Client.SetOutputDirectory.

    client.R().
     SetOutput("/Users/jeeva/Downloads/ReplyWithHeader-v5.1-beta.zip").
     Get("http://bit.ly/1LouEKr")
    

NOTE: In this scenario Response.Body might be nil.

#### func (*Request) [SetPathParam](https://github.com/go-resty/resty/blob/v2.17.1/request.go#L675) ¶ added in v2.4.0

    func (r *Request) SetPathParam(param, value [string](/builtin#string)) *Request

SetPathParam method sets a single URL path key-value pair in the Resty current request instance.

    client.R().SetPathParam("userId", "sample@sample.com")
    
    Result:
       URL - /v1/users/{userId}/details
       Composed URL - /v1/users/sample@sample.com/details
    
    client.R().SetPathParam("path", "groups/developers")
    
    Result:
       URL - /v1/users/{userId}/details
       Composed URL - /v1/users/groups%2Fdevelopers/details
    

It replaces the value of the key while composing the request URL. The values will be escaped using function [url.PathEscape](/net/url#PathEscape).

It overrides the path parameter set at the client instance level.

#### func (*Request) [SetPathParams](https://github.com/go-resty/resty/blob/v2.17.1/request.go#L697) ¶

    func (r *Request) SetPathParams(params map[[string](/builtin#string)][string](/builtin#string)) *Request

SetPathParams method sets multiple URL path key-value pairs at one go in the Resty current request instance.

    client.R().SetPathParams(map[string]string{
     "userId":       "sample@sample.com",
     "subAccountId": "100002",
     "path":         "groups/developers",
    })
    
    Result:
       URL - /v1/users/{userId}/{subAccountId}/{path}/details
       Composed URL - /v1/users/sample@sample.com/100002/groups%2Fdevelopers/details
    

It replaces the value of the key while composing the request URL. The values will be escaped using function [url.PathEscape](/net/url#PathEscape).

It overrides the path parameter set at the client instance level.

#### func (*Request) [SetQueryParam](https://github.com/go-resty/resty/blob/v2.17.1/request.go#L190) ¶

    func (r *Request) SetQueryParam(param, value [string](/builtin#string)) *Request

SetQueryParam method sets a single parameter and its value in the current request. It will be formed as a query string for the request.

For Example: `search=kitchen%20papers&size=large` in the URL after the `?` mark.

    client.R().
     SetQueryParam("search", "kitchen papers").
     SetQueryParam("size", "large")
    

It overrides the query parameter value set at the client instance level.

#### func (*Request) [SetQueryParams](https://github.com/go-resty/resty/blob/v2.17.1/request.go#L207) ¶

    func (r *Request) SetQueryParams(params map[[string](/builtin#string)][string](/builtin#string)) *Request

SetQueryParams method sets multiple parameters and their values at one go in the current request. It will be formed as a query string for the request.

For Example: `search=kitchen%20papers&size=large` in the URL after the `?` mark.

    client.R().
     SetQueryParams(map[string]string{
      "search": "kitchen papers",
      "size": "large",
     })
    

It overrides the query parameter value set at the client instance level.

#### func (*Request) [SetQueryParamsFromValues](https://github.com/go-resty/resty/blob/v2.17.1/request.go#L237) ¶

    func (r *Request) SetQueryParamsFromValues(params [url](/net/url).[Values](/net/url#Values)) *Request

SetQueryParamsFromValues method appends multiple parameters with multi-value ([url.Values](/net/url#Values)) at one go in the current request. It will be formed as query string for the request.

For Example: `status=pending&status=approved&status=open` in the URL after the `?` mark.

    client.R().
     SetQueryParamsFromValues(url.Values{
      "status": []string{"pending", "approved", "open"},
     })
    

It overrides the query parameter value set at the client instance level.

#### func (*Request) [SetQueryString](https://github.com/go-resty/resty/blob/v2.17.1/request.go#L252) ¶

    func (r *Request) SetQueryString(query [string](/builtin#string)) *Request

SetQueryString method provides the ability to use string as an input to set URL query string for the request.

    client.R().
     SetQueryString("productId=232&template=fresh-sample&cat=resty&source=google&kw=buy a lot more")
    

It overrides the query parameter value set at the client instance level.

#### func (*Request) [SetRawPathParam](https://github.com/go-resty/resty/blob/v2.17.1/request.go#L723) ¶ added in v2.8.0

    func (r *Request) SetRawPathParam(param, value [string](/builtin#string)) *Request

SetRawPathParam method sets a single URL path key-value pair in the Resty current request instance.

    client.R().SetPathParam("userId", "sample@sample.com")
    
    Result:
       URL - /v1/users/{userId}/details
       Composed URL - /v1/users/sample@sample.com/details
    
    client.R().SetPathParam("path", "groups/developers")
    
    Result:
       URL - /v1/users/{userId}/details
       Composed URL - /v1/users/groups/developers/details
    

It replaces the value of the key while composing the request URL. The value will be used as-is and has not been escaped.

It overrides the raw path parameter set at the client instance level.

#### func (*Request) [SetRawPathParams](https://github.com/go-resty/resty/blob/v2.17.1/request.go#L745) ¶ added in v2.8.0

    func (r *Request) SetRawPathParams(params map[[string](/builtin#string)][string](/builtin#string)) *Request

SetRawPathParams method sets multiple URL path key-value pairs at one go in the Resty current request instance.

    client.R().SetPathParams(map[string]string{
     "userId": "sample@sample.com",
     "subAccountId": "100002",
     "path":         "groups/developers",
    })
    
    Result:
       URL - /v1/users/{userId}/{subAccountId}/{path}/details
       Composed URL - /v1/users/sample@sample.com/100002/groups/developers/details
    

It replaces the value of the key while composing the request URL. The value will be used as-is and has not been escaped.

It overrides the raw path parameter set at the client instance level.

#### func (*Request) [SetResponseBodyLimit](https://github.com/go-resty/resty/blob/v2.17.1/request.go#L651) ¶ added in v2.15.0

    func (r *Request) SetResponseBodyLimit(v [int](/builtin#int)) *Request

SetResponseBodyLimit method sets a maximum body size limit in bytes on response, avoid reading too much data to memory.

Client will return resty.ErrResponseBodyTooLarge if the body size of the body in the uncompressed response is larger than the limit. Body size limit will not be enforced in the following cases:

- ResponseBodyLimit <= 0, which is the default behavior.
- Request.SetOutput is called to save response data to the file.
- "DoNotParseResponse" is set for client or request.

It overrides the value set at the client instance level. see Client.SetResponseBodyLimit

#### func (*Request) [SetResult](https://github.com/go-resty/resty/blob/v2.17.1/request.go#L378) ¶

    func (r *Request) SetResult(res interface{}) *Request

SetResult method is to register the response `Result` object for automatic unmarshalling of the HTTP response if the response status code is between 200 and 299, and the content type is JSON or XML.

Note: Request.SetResult input can be a pointer or non-pointer.

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
    

#### func (*Request) [SetSRV](https://github.com/go-resty/resty/blob/v2.17.1/request.go#L623) ¶

    func (r *Request) SetSRV(srv *SRVRecord) *Request

SetSRV method sets the details to query the service SRV record and execute the request.

    client.R().
     SetSRV(SRVRecord{"web", "testservice.com"}).
     Get("/get")
    

#### func (*Request) [SetUnescapeQueryParams](https://github.com/go-resty/resty/blob/v2.17.1/request.go#L220) ¶ added in v2.16.0

    func (r *Request) SetUnescapeQueryParams(unescape [bool](/builtin#bool)) *Request

SetUnescapeQueryParams method sets the unescape query parameters choice for request URL. To prevent broken URL, resty replaces space (" ") with "+" in the query parameters.

This method overrides the value set by Client.SetUnescapeQueryParams

NOTE: Request failure is possible due to non-standard usage of Unescaped Query Parameters.

#### func (*Request) [TraceInfo](https://github.com/go-resty/resty/blob/v2.17.1/request.go#L893) ¶

    func (r *Request) TraceInfo() TraceInfo

TraceInfo method returns the trace info for the request. If either the Client.EnableTrace or Request.EnableTrace function has not been called before the request is made, an empty resty.TraceInfo object is returned.

#### type [RequestLog](https://github.com/go-resty/resty/blob/v2.17.1/util.go#L134) ¶

    type RequestLog struct {
     Header [http](/net/http).[Header](/net/http#Header)
     Body   [string](/builtin#string)
    }

RequestLog struct is used to collected information from resty request instance for debug logging. It sent to request log callback before resty actually logs the information.

#### type [RequestLogCallback](https://github.com/go-resty/resty/blob/v2.17.1/client.go#L83) ¶ added in v2.3.0

    type RequestLogCallback func(*RequestLog) [error](/builtin#error)

RequestLogCallback type is for request logs, called before the request is logged

#### type [RequestMiddleware](https://github.com/go-resty/resty/blob/v2.17.1/client.go#L74) ¶ added in v2.3.0

    type RequestMiddleware func(*Client, *Request) [error](/builtin#error)

RequestMiddleware type is for request middleware, called before a request is sent

#### type [Response](https://github.com/go-resty/resty/blob/v2.17.1/response.go#L21) ¶

    type Response struct {
     Request     *Request
     RawResponse *[http](/net/http).[Response](/net/http#Response)
     // contains filtered or unexported fields
    }

Response struct holds response values of executed requests.

#### func (*Response) [Body](https://github.com/go-resty/resty/blob/v2.17.1/response.go#L34) ¶

    func (r *Response) Body() [][byte](/builtin#byte)

Body method returns the HTTP response as `[]byte` slice for the executed request.

NOTE: Response.Body might be nil if Request.SetOutput is used. Also see Request.SetDoNotParseResponse, Client.SetDoNotParseResponse

#### func (*Response) [Cookies](https://github.com/go-resty/resty/blob/v2.17.1/response.go#L102) ¶

    func (r *Response) Cookies() []*[http](/net/http).[Cookie](/net/http#Cookie)

Cookies method to returns all the response cookies

#### func (*Response) [Error](https://github.com/go-resty/resty/blob/v2.17.1/response.go#L89) ¶

    func (r *Response) Error() interface{}

Error method returns the error object if it has one

See Request.SetError, Client.SetError

#### func (*Response) [Header](https://github.com/go-resty/resty/blob/v2.17.1/response.go#L94) ¶

    func (r *Response) Header() [http](/net/http).[Header](/net/http#Header)

Header method returns the response headers

#### func (*Response) [IsError](https://github.com/go-resty/resty/blob/v2.17.1/response.go#L162) ¶

    func (r *Response) IsError() [bool](/builtin#bool)

IsError method returns true if HTTP status `code >= 400` otherwise false.

#### func (*Response) [IsSuccess](https://github.com/go-resty/resty/blob/v2.17.1/response.go#L157) ¶

    func (r *Response) IsSuccess() [bool](/builtin#bool)

IsSuccess method returns true if HTTP status `code >= 200 and <= 299` otherwise false.

#### func (*Response) [Proto](https://github.com/go-resty/resty/blob/v2.17.1/response.go#L72) ¶ added in v2.3.0

    func (r *Response) Proto() [string](/builtin#string)

Proto method returns the HTTP response protocol used for the request.

#### func (*Response) [RawBody](https://github.com/go-resty/resty/blob/v2.17.1/response.go#L149) ¶

    func (r *Response) RawBody() [io](/io).[ReadCloser](/io#ReadCloser)

RawBody method exposes the HTTP raw response body. Use this method in conjunction with Client.SetDoNotParseResponse or Request.SetDoNotParseResponse option; otherwise, you get an error as `read err: http: read on closed response body.`

Do not forget to close the body, otherwise you might get into connection leaks, no connection reuse. You have taken over the control of response parsing from Resty.

#### func (*Response) [ReceivedAt](https://github.com/go-resty/resty/blob/v2.17.1/response.go#L131) ¶

    func (r *Response) ReceivedAt() [time](/time).[Time](/time#Time)

ReceivedAt method returns the time we received a response from the server for the request.

#### func (*Response) [Result](https://github.com/go-resty/resty/blob/v2.17.1/response.go#L82) ¶

    func (r *Response) Result() interface{}

Result method returns the response value as an object if it has one

See Request.SetResult

#### func (*Response) [SetBody](https://github.com/go-resty/resty/blob/v2.17.1/response.go#L46) ¶ added in v2.10.0

    func (r *Response) SetBody(b [][byte](/builtin#byte)) *Response

SetBody method sets Response body in byte slice. Typically, It is helpful for test cases.

    resp.SetBody([]byte("This is test body content"))
    resp.SetBody(nil)
    

#### func (*Response) [Size](https://github.com/go-resty/resty/blob/v2.17.1/response.go#L139) ¶

    func (r *Response) Size() [int64](/builtin#int64)

Size method returns the HTTP response size in bytes. Yeah, you can rely on HTTP `Content-Length` header, however it won't be available for chucked transfer/compressed response. Since Resty captures response size details when processing the response body when possible. So that users get the actual size of response bytes.

#### func (*Response) [Status](https://github.com/go-resty/resty/blob/v2.17.1/response.go#L54) ¶

    func (r *Response) Status() [string](/builtin#string)

Status method returns the HTTP status string for the executed request.

    Example: 200 OK
    

#### func (*Response) [StatusCode](https://github.com/go-resty/resty/blob/v2.17.1/response.go#L64) ¶

    func (r *Response) StatusCode() [int](/builtin#int)

StatusCode method returns the HTTP status code for the executed request.

    Example: 200
    

#### func (*Response) [String](https://github.com/go-resty/resty/blob/v2.17.1/response.go#L111) ¶

    func (r *Response) String() [string](/builtin#string)

String method returns the body of the HTTP response as a `string`. It returns an empty string if it is nil or the body is zero length.

#### func (*Response) [Time](https://github.com/go-resty/resty/blob/v2.17.1/response.go#L123) ¶

    func (r *Response) Time() [time](/time).[Duration](/time#Duration)

Time method returns the duration of HTTP response time from the request we sent and received a request.

See Response.ReceivedAt to know when the client received a response and see `Response.Request.Time` to know when the client sent a request.

#### type [ResponseError](https://github.com/go-resty/resty/blob/v2.17.1/client.go#L1376) ¶ added in v2.4.0

    type ResponseError struct {
     Response *Response
     Err      [error](/builtin#error)
    }

ResponseError is a wrapper that includes the server response with an error. Neither the err nor the response should be nil.

#### func (*ResponseError) [Error](https://github.com/go-resty/resty/blob/v2.17.1/client.go#L1381) ¶ added in v2.4.0

    func (e *ResponseError) Error() [string](/builtin#string)

#### func (*ResponseError) [Unwrap](https://github.com/go-resty/resty/blob/v2.17.1/client.go#L1385) ¶ added in v2.4.0

    func (e *ResponseError) Unwrap() [error](/builtin#error)

#### type [ResponseLog](https://github.com/go-resty/resty/blob/v2.17.1/util.go#L142) ¶

    type ResponseLog struct {
     Header [http](/net/http).[Header](/net/http#Header)
     Body   [string](/builtin#string)
    }

ResponseLog struct is used to collected information from resty response instance for debug logging. It sent to response log callback before resty actually logs the information.

#### type [ResponseLogCallback](https://github.com/go-resty/resty/blob/v2.17.1/client.go#L86) ¶ added in v2.3.0

    type ResponseLogCallback func(*ResponseLog) [error](/builtin#error)

ResponseLogCallback type is for response logs, called before the response is logged

#### type [ResponseMiddleware](https://github.com/go-resty/resty/blob/v2.17.1/client.go#L77) ¶ added in v2.3.0

    type ResponseMiddleware func(*Client, *Response) [error](/builtin#error)

ResponseMiddleware type is for response middleware, called after a response has been received

#### type [RetryAfterFunc](https://github.com/go-resty/resty/blob/v2.17.1/retry.go#L38) ¶

    type RetryAfterFunc func(*Client, *Response) ([time](/time).[Duration](/time#Duration), [error](/builtin#error))

RetryAfterFunc returns time to wait before retry For example, it can parse HTTP Retry-After header <https://www.w3.org/Protocols/rfc2616/rfc2616-sec14.html> Non-nil error is returned if it is found that the request is not retryable (0, nil) is a special result that means 'use default algorithm'

#### type [RetryConditionFunc](https://github.com/go-resty/resty/blob/v2.17.1/retry.go#L28) ¶

    type RetryConditionFunc func(*Response, [error](/builtin#error)) [bool](/builtin#bool)

RetryConditionFunc type is for the retry condition function input: non-nil Response OR request execution error

#### type [SRVRecord](https://github.com/go-resty/resty/blob/v2.17.1/request.go#L1091) ¶

    type SRVRecord struct {
     Service [string](/builtin#string)
     Domain  [string](/builtin#string)
    }

SRVRecord struct holds the data to query the SRV record for the following service.

#### type [SuccessHook](https://github.com/go-resty/resty/blob/v2.17.1/client.go#L92) ¶ added in v2.8.0

    type SuccessHook func(*Client, *Response)

SuccessHook type is for reacting to request success

#### type [TraceInfo](https://github.com/go-resty/resty/blob/v2.17.1/trace.go#L22) ¶

    type TraceInfo struct {
     // DNSLookup is the duration that transport took to perform
     // DNS lookup.
     DNSLookup [time](/time).[Duration](/time#Duration)
    
     // ConnTime is the duration it took to obtain a successful connection.
     ConnTime [time](/time).[Duration](/time#Duration)
    
     // TCPConnTime is the duration it took to obtain the TCP connection.
     TCPConnTime [time](/time).[Duration](/time#Duration)
    
     // TLSHandshake is the duration of the TLS handshake.
     TLSHandshake [time](/time).[Duration](/time#Duration)
    
     // ServerTime is the server's duration for responding to the first byte.
     ServerTime [time](/time).[Duration](/time#Duration)
    
     // ResponseTime is the duration since the first response byte from the server to
     // request completion.
     ResponseTime [time](/time).[Duration](/time#Duration)
    
     // TotalTime is the duration of the total time request taken end-to-end.
     TotalTime [time](/time).[Duration](/time#Duration)
    
     // IsConnReused is whether this connection has been previously
     // used for another HTTP request.
     IsConnReused [bool](/builtin#bool)
    
     // IsConnWasIdle is whether this connection was obtained from an
     // idle pool.
     IsConnWasIdle [bool](/builtin#bool)
    
     // ConnIdleTime is the duration how long the connection that was previously
     // idle, if IsConnWasIdle is true.
     ConnIdleTime [time](/time).[Duration](/time#Duration)
    
     // RequestAttempt is to represent the request attempt made during a Resty
     // request execution flow, including retry count.
     RequestAttempt [int](/builtin#int)
    
     // RemoteAddr returns the remote network address.
     RemoteAddr [net](/net).[Addr](/net#Addr)
    }

TraceInfo struct is used to provide request trace info such as DNS lookup duration, Connection obtain duration, Server processing duration, etc.

#### type [User](https://github.com/go-resty/resty/blob/v2.17.1/client.go#L162) ¶

    type User struct {
     Username, Password [string](/builtin#string)
    }

User type is to hold an username and password information
