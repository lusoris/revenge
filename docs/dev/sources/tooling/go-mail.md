# wneessen/go-mail

> Source: https://pkg.go.dev/github.com/wneessen/go-mail
> Fetched: 2026-02-01T11:43:12.556025+00:00
> Content-Hash: 766154dbb365a038
> Type: html

---

### Overview ¶

Package mail provides an easy to use interface for formating and sending mails. go-mail follows idiomatic Go style and best practice. It has a small dependency footprint by mainly relying on the Go Standard Library and the Go extended packages. It combines a lot of functionality from the standard library to give easy and convenient access to mail and SMTP related tasks. It works like a programmatic email client and provides lots of methods and functionalities you would consider standard in a MUA. 

### Index ¶

  * Constants
  * Variables
  * func IsAddrHeader(header string) bool
  * type AddrHeader
  *     * func (a AddrHeader) String() string
  * type AuthData
  *     * func NewAuthData(user, pass string) *AuthData
  * type Charset
  *     * func (c Charset) String() string
  * type Client
  *     * func NewClient(host string, opts ...Option) (*Client, error)
  *     * func (c *Client) Close() error
    * func (c *Client) CloseWithSMTPClient(client *smtp.Client) error
    * func (c *Client) DialAndSend(messages ...*Msg) error
    * func (c *Client) DialAndSendWithContext(ctx context.Context, messages ...*Msg) error
    * func (c *Client) DialToSMTPClientWithContext(ctxDial context.Context) (*smtp.Client, error)
    * func (c *Client) DialWithContext(ctxDial context.Context) error
    * func (c *Client) Reset() error
    * func (c *Client) ResetWithSMTPClient(client *smtp.Client) error
    * func (c *Client) Send(messages ...*Msg) (returnErr error)
    * func (c *Client) SendWithSMTPClient(client *smtp.Client, messages ...*Msg) (returnErr error)
    * func (c *Client) ServerAddr() string
    * func (c *Client) SetDebugLog(val bool)
    * func (c *Client) SetDebugLogWithSMTPClient(client *smtp.Client, val bool)
    * func (c *Client) SetLogAuthData(logAuth bool)
    * func (c *Client) SetLogger(logger log.Logger)
    * func (c *Client) SetLoggerWithSMTPClient(client *smtp.Client, logger log.Logger)
    * func (c *Client) SetPassword(password string)
    * func (c *Client) SetSMTPAuth(authtype SMTPAuthType)
    * func (c *Client) SetSMTPAuthCustom(smtpAuth smtp.Auth)
    * func (c *Client) SetSSL(ssl bool)
    * func (c *Client) SetSSLPort(ssl bool, fallback bool)
    * func (c *Client) SetTLSConfig(tlsconfig *tls.Config) error
    * func (c *Client) SetTLSPolicy(policy TLSPolicy)
    * func (c *Client) SetTLSPortPolicy(policy TLSPolicy)
    * func (c *Client) SetUsername(username string)
    * func (c *Client) TLSPolicy() string
  * type ContentType
  *     * func (c ContentType) String() string
  * type DSNMailReturnOption
  * type DSNRcptNotifyOption
  * type DialContextFunc
  * type Encoding
  *     * func (e Encoding) String() string
  * type File
  * type FileOption
  *     * func WithFileContentID(id string) FileOption
    * func WithFileContentType(contentType ContentType) FileOption
    * func WithFileDescription(description string) FileOption
    * func WithFileEncoding(encoding Encoding) FileOption
    * func WithFileName(name string) FileOption
  * type Header
  *     * func (h Header) String() string
  * type Importance
  *     * func (i Importance) NumString() string
    * func (i Importance) String() string
    * func (i Importance) XPrioString() string
  * type MIMEType
  *     * func (e MIMEType) String() string
  * type MIMEVersion
  * type Middleware
  * type MiddlewareType
  * type Msg
  *     * func EMLToMsgFromFile(filePath string) (*Msg, error)
    * func EMLToMsgFromReader(reader io.Reader) (*Msg, error)
    * func EMLToMsgFromString(emlString string) (*Msg, error)
    * func NewMsg(opts ...MsgOption) *Msg
    * func QuickSend(addr string, auth *AuthData, from string, rcpts []string, subject string, ...) (*Msg, error)
  *     * func (m *Msg) AddAlternativeHTMLTemplate(tpl *ht.Template, data interface{}, opts ...PartOption) error
    * func (m *Msg) AddAlternativeString(contentType ContentType, content string, opts ...PartOption)
    * func (m *Msg) AddAlternativeTextTemplate(tpl *tt.Template, data interface{}, opts ...PartOption) error
    * func (m *Msg) AddAlternativeWriter(contentType ContentType, writeFunc func(io.Writer) (int64, error), ...)
    * func (m *Msg) AddBcc(rcpt string) error
    * func (m *Msg) AddBccFormat(name, addr string) error
    * func (m *Msg) AddBccMailAddress(rcpt *mail.Address)
    * func (m *Msg) AddCc(rcpt string) error
    * func (m *Msg) AddCcFormat(name, addr string) error
    * func (m *Msg) AddCcMailAddress(rcpt *mail.Address)
    * func (m *Msg) AddTo(rcpt string) error
    * func (m *Msg) AddToFormat(name, addr string) error
    * func (m *Msg) AddToMailAddress(rcpt *mail.Address)
    * func (m *Msg) AttachFile(name string, opts ...FileOption)
    * func (m *Msg) AttachFromEmbedFS(name string, fs *embed.FS, opts ...FileOption) error
    * func (m *Msg) AttachFromIOFS(name string, iofs fs.FS, opts ...FileOption) error
    * func (m *Msg) AttachHTMLTemplate(name string, tpl *ht.Template, data interface{}, opts ...FileOption) error
    * func (m *Msg) AttachReadSeeker(name string, reader io.ReadSeeker, opts ...FileOption)
    * func (m *Msg) AttachReader(name string, reader io.Reader, opts ...FileOption) error
    * func (m *Msg) AttachTextTemplate(name string, tpl *tt.Template, data interface{}, opts ...FileOption) error
    * func (m *Msg) Bcc(rcpts ...string) error
    * func (m *Msg) BccFromString(rcpts string) error
    * func (m *Msg) BccIgnoreInvalid(rcpts ...string)
    * func (m *Msg) BccMailAddress(rcpts ...*mail.Address)
    * func (m *Msg) Cc(rcpts ...string) error
    * func (m *Msg) CcFromString(rcpts string) error
    * func (m *Msg) CcIgnoreInvalid(rcpts ...string)
    * func (m *Msg) CcMailAddress(rcpts ...*mail.Address)
    * func (m *Msg) Charset() string
    * func (m *Msg) EmbedFile(name string, opts ...FileOption)
    * func (m *Msg) EmbedFromEmbedFS(name string, fs *embed.FS, opts ...FileOption) error
    * func (m *Msg) EmbedFromIOFS(name string, iofs fs.FS, opts ...FileOption) error
    * func (m *Msg) EmbedHTMLTemplate(name string, tpl *ht.Template, data interface{}, opts ...FileOption) error
    * func (m *Msg) EmbedReadSeeker(name string, reader io.ReadSeeker, opts ...FileOption)
    * func (m *Msg) EmbedReader(name string, reader io.Reader, opts ...FileOption) error
    * func (m *Msg) EmbedTextTemplate(name string, tpl *tt.Template, data interface{}, opts ...FileOption) error
    * func (m *Msg) Encoding() string
    * func (m *Msg) EnvelopeFrom(from string) error
    * func (m *Msg) EnvelopeFromFormat(name, addr string) error
    * func (m *Msg) EnvelopeFromMailAddress(addr *mail.Address)
    * func (m *Msg) From(from string) error
    * func (m *Msg) FromFormat(name, addr string) error
    * func (m *Msg) FromMailAddress(from *mail.Address)
    * func (m *Msg) GetAddrHeader(header AddrHeader) []*mail.Address
    * func (m *Msg) GetAddrHeaderString(header AddrHeader) []string
    * func (m *Msg) GetAttachments() []*File
    * func (m *Msg) GetBcc() []*mail.Address
    * func (m *Msg) GetBccString() []string
    * func (m *Msg) GetBoundary() string
    * func (m *Msg) GetCc() []*mail.Address
    * func (m *Msg) GetCcString() []string
    * func (m *Msg) GetEmbeds() []*File
    * func (m *Msg) GetFrom() []*mail.Address
    * func (m *Msg) GetFromString() []string
    * func (m *Msg) GetGenHeader(header Header) []string
    * func (m *Msg) GetMessageID() string
    * func (m *Msg) GetParts() []*Part
    * func (m *Msg) GetRecipients() ([]string, error)
    * func (m *Msg) GetSender(useFullAddr bool) (string, error)
    * func (m *Msg) GetTo() []*mail.Address
    * func (m *Msg) GetToString() []string
    * func (m *Msg) HasSendError() bool
    * func (m *Msg) IsDelivered() bool
    * func (m *Msg) NewReader() *Reader
    * func (m *Msg) ReplyTo(addr string) error
    * func (m *Msg) ReplyToFormat(name, addr string) error
    * func (m *Msg) ReplyToMailAddress(addr *mail.Address)
    * func (m *Msg) RequestMDNAddTo(rcpt string) error
    * func (m *Msg) RequestMDNAddToFormat(name, addr string) error
    * func (m *Msg) RequestMDNTo(rcpts ...string) error
    * func (m *Msg) RequestMDNToFormat(name, addr string) error
    * func (m *Msg) Reset()
    * func (m *Msg) SendError() error
    * func (m *Msg) SendErrorIsTemp() bool
    * func (m *Msg) ServerResponse() string
    * func (m *Msg) SetAddrHeader(header AddrHeader, values ...string) error
    * func (m *Msg) SetAddrHeaderFromMailAddress(header AddrHeader, values ...*mail.Address)
    * func (m *Msg) SetAddrHeaderIgnoreInvalid(header AddrHeader, values ...string)
    * func (m *Msg) SetAttachements(files []*File)deprecated
    * func (m *Msg) SetAttachments(files []*File)
    * func (m *Msg) SetBodyHTMLTemplate(tpl *ht.Template, data interface{}, opts ...PartOption) error
    * func (m *Msg) SetBodyString(contentType ContentType, content string, opts ...PartOption)
    * func (m *Msg) SetBodyTextTemplate(tpl *tt.Template, data interface{}, opts ...PartOption) error
    * func (m *Msg) SetBodyWriter(contentType ContentType, writeFunc func(io.Writer) (int64, error), ...)
    * func (m *Msg) SetBoundary(boundary string)
    * func (m *Msg) SetBulk()
    * func (m *Msg) SetCharset(charset Charset)
    * func (m *Msg) SetDate()
    * func (m *Msg) SetDateWithValue(timeVal time.Time)
    * func (m *Msg) SetEmbeds(files []*File)
    * func (m *Msg) SetEncoding(encoding Encoding)
    * func (m *Msg) SetGenHeader(header Header, values ...string)
    * func (m *Msg) SetGenHeaderPreformatted(header Header, value string)
    * func (m *Msg) SetHeader(header Header, values ...string)deprecated
    * func (m *Msg) SetHeaderPreformatted(header Header, value string)deprecated
    * func (m *Msg) SetImportance(importance Importance)
    * func (m *Msg) SetMIMEVersion(version MIMEVersion)
    * func (m *Msg) SetMessageID()
    * func (m *Msg) SetMessageIDWithValue(messageID string)
    * func (m *Msg) SetOrganization(org string)
    * func (m *Msg) SetPGPType(pgptype PGPType)
    * func (m *Msg) SetUserAgent(userAgent string)
    * func (m *Msg) SignWithKeypair(privateKey crypto.PrivateKey, certificate *x509.Certificate, ...) error
    * func (m *Msg) SignWithTLSCertificate(keyPairTLS *tls.Certificate) error
    * func (m *Msg) Subject(subj string)
    * func (m *Msg) To(rcpts ...string) error
    * func (m *Msg) ToFromString(rcpts string) error
    * func (m *Msg) ToIgnoreInvalid(rcpts ...string)
    * func (m *Msg) ToMailAddress(rcpts ...*mail.Address)
    * func (m *Msg) UnsetAllAttachments()
    * func (m *Msg) UnsetAllEmbeds()
    * func (m *Msg) UnsetAllParts()
    * func (m *Msg) UpdateReader(reader *Reader)
    * func (m *Msg) Write(writer io.Writer) (int64, error)
    * func (m *Msg) WriteTo(writer io.Writer) (int64, error)
    * func (m *Msg) WriteToFile(name string) error
    * func (m *Msg) WriteToSendmail() error
    * func (m *Msg) WriteToSendmailWithCommand(sendmailPath string) error
    * func (m *Msg) WriteToSendmailWithContext(ctx context.Context, sendmailPath string, args ...string) error
    * func (m *Msg) WriteToSkipMiddleware(writer io.Writer, middleWareType MiddlewareType) (int64, error)
    * func (m *Msg) WriteToTempFile() (string, error)
  * type MsgOption
  *     * func WithBoundary(boundary string) MsgOption
    * func WithCharset(charset Charset) MsgOption
    * func WithEncoding(encoding Encoding) MsgOption
    * func WithMIMEVersion(version MIMEVersion) MsgOption
    * func WithMiddleware(middleware Middleware) MsgOption
    * func WithNoDefaultUserAgent() MsgOption
    * func WithPGPType(pgptype PGPType) MsgOption
  * type Option
  *     * func WithDSN() Option
    * func WithDSNMailReturnType(option DSNMailReturnOption) Option
    * func WithDSNRcptNotifyType(opts ...DSNRcptNotifyOption) Option
    * func WithDebugLog() Option
    * func WithDialContextFunc(dialCtxFunc DialContextFunc) Option
    * func WithHELO(helo string) Option
    * func WithLogAuthData() Option
    * func WithLogger(logger log.Logger) Option
    * func WithPassword(password string) Option
    * func WithPort(port int) Option
    * func WithSMTPAuth(authtype SMTPAuthType) Option
    * func WithSMTPAuthCustom(smtpAuth smtp.Auth) Option
    * func WithSSL() Option
    * func WithSSLPort(fallback bool) Option
    * func WithTLSConfig(tlsconfig *tls.Config) Option
    * func WithTLSPolicy(policy TLSPolicy) Option
    * func WithTLSPortPolicy(policy TLSPolicy) Option
    * func WithTimeout(timeout time.Duration) Option
    * func WithUsername(username string) Option
    * func WithoutNoop() Option
  * type PGPType
  * type Part
  *     * func (p *Part) Delete()
    * func (p *Part) GetCharset() Charset
    * func (p *Part) GetContent() ([]byte, error)
    * func (p *Part) GetContentType() ContentType
    * func (p *Part) GetDescription() string
    * func (p *Part) GetEncoding() Encoding
    * func (p *Part) GetWriteFunc() func(io.Writer) (int64, error)
    * func (p *Part) SetCharset(charset Charset)
    * func (p *Part) SetContent(content string)
    * func (p *Part) SetContentType(contentType ContentType)
    * func (p *Part) SetDescription(description string)
    * func (p *Part) SetEncoding(encoding Encoding)
    * func (p *Part) SetIsSMIMESigned(smime bool)
    * func (p *Part) SetWriteFunc(writeFunc func(io.Writer) (int64, error))
  * type PartOption
  *     * func WithPartCharset(charset Charset) PartOption
    * func WithPartContentDescription(description string) PartOption
    * func WithPartEncoding(encoding Encoding) PartOption
    * func WithSMIMESigning() PartOption
  * type Reader
  *     * func (r *Reader) Error() error
    * func (r *Reader) Read(payload []byte) (n int, err error)
    * func (r *Reader) Reset()
  * type SMIME
  * type SMTPAuthType
  *     * func (sa *SMTPAuthType) UnmarshalString(value string) error
  * type SendErrReason
  *     * func (r SendErrReason) String() string
  * type SendError
  *     * func (e *SendError) EnhancedStatusCode() string
    * func (e *SendError) Error() string
    * func (e *SendError) ErrorCode() int
    * func (e *SendError) Is(errType error) bool
    * func (e *SendError) IsTemp() bool
    * func (e *SendError) MessageID() string
    * func (e *SendError) Msg() *Msg
  * type TLSPolicy
  *     * func (p TLSPolicy) String() string



### Examples ¶

  * Client.DialAndSend
  * Client.SetTLSPolicy
  * Msg.SetBodyString (DifferentTypes)
  * Msg.SetBodyString (WithPartOption)
  * Msg.SetBodyTextTemplate
  * Msg.WriteToSendmail
  * Msg.WriteToSendmailWithContext
  * NewClient
  * NewMsg



### Constants ¶

[View Source](https://github.com/wneessen/go-mail/blob/v0.7.2/client.go#L22)
    
    
    const (
    	// DefaultPort is the default connection port to the SMTP server.
    	DefaultPort = 25
    
    	// DefaultPortSSL is the default connection port for SSL/TLS to the SMTP server.
    	DefaultPortSSL = 465
    
    	// DefaultPortTLS is the default connection port for STARTTLS to the SMTP server.
    	DefaultPortTLS = 587
    
    	// DefaultTimeout is the default connection timeout.
    	DefaultTimeout = [time](/time).[Second](/time#Second) * 15
    
    	// DefaultTLSPolicy specifies the default TLS policy for connections.
    	DefaultTLSPolicy = TLSMandatory
    
    	// DefaultTLSMinVersion defines the minimum TLS version to be used for secure connections.
    	// Nowadays TLS 1.2 is assumed be a sane default.
    	DefaultTLSMinVersion = [tls](/crypto/tls).[VersionTLS12](/crypto/tls#VersionTLS12)
    )

[View Source](https://github.com/wneessen/go-mail/blob/v0.7.2/client.go#L43)
    
    
    const (
    	// DSNMailReturnHeadersOnly requests that only the message headers of the mail message are returned in
    	// a DSN (Delivery Status Notification).
    	//
    	// <https://datatracker.ietf.org/doc/html/rfc1891#section-5.3>
    	DSNMailReturnHeadersOnly DSNMailReturnOption = "HDRS"
    
    	// DSNMailReturnFull requests that the entire mail message is returned in any failed  DSN
    	// (Delivery Status Notification) issued for this recipient.
    	//
    	// <https://datatracker.ietf.org/doc/html/rfc1891/#section-5.3>
    	DSNMailReturnFull DSNMailReturnOption = "FULL"
    
    	// DSNRcptNotifyNever indicates that no DSN (Delivery Status Notifications) should be sent for the
    	// recipient under any condition.
    	//
    	// <https://datatracker.ietf.org/doc/html/rfc1891/#section-5.1>
    	DSNRcptNotifyNever DSNRcptNotifyOption = "NEVER"
    
    	// DSNRcptNotifySuccess indicates that the sender requests a DSN (Delivery Status Notification) if the
    	// message is successfully delivered.
    	//
    	// <https://datatracker.ietf.org/doc/html/rfc1891/#section-5.1>
    	DSNRcptNotifySuccess DSNRcptNotifyOption = "SUCCESS"
    
    	// DSNRcptNotifyFailure requests that a DSN (Delivery Status Notification) is issued if delivery of
    	// a message fails.
    	//
    	// <https://datatracker.ietf.org/doc/html/rfc1891/#section-5.1>
    	DSNRcptNotifyFailure DSNRcptNotifyOption = "FAILURE"
    
    	// DSNRcptNotifyDelay indicates the sender's willingness to receive "delayed" DSNs.
    	//
    	// Delayed DSNs may be issued if delivery of a message has been delayed for an unusual amount of time
    	// (as determined by the MTA at which the message is delayed), but the final delivery status (whether
    	// successful or failure) cannot be determined. The absence of the DELAY keyword in a NOTIFY parameter
    	// requests that a "delayed" DSN NOT be issued under any conditions.
    	//
    	// <https://datatracker.ietf.org/doc/html/rfc1891/#section-5.1>
    	DSNRcptNotifyDelay DSNRcptNotifyOption = "DELAY"
    )

[View Source](https://github.com/wneessen/go-mail/blob/v0.7.2/msgwriter.go#L21)
    
    
    const (
    	// MaxHeaderLength defines the maximum line length for a mail header.
    	//
    	// This constant follows the recommendation of [RFC 2047](https://rfc-editor.org/rfc/rfc2047.html), which suggests a maximum length of 76 characters.
    	//
    	// References:
    	//   - <https://datatracker.ietf.org/doc/html/rfc2047>
    	MaxHeaderLength = 76
    
    	// MaxBodyLength defines the maximum line length for the mail body.
    	//
    	// This constant follows the recommendation of [RFC 2047](https://rfc-editor.org/rfc/rfc2047.html), which suggests a maximum length of 76 characters.
    	//
    	// References:
    	//   - <https://datatracker.ietf.org/doc/html/rfc2047>
    	MaxBodyLength = 76
    
    	// SingleNewLine represents a single newline character sequence ("\r\n").
    	//
    	// This constant can be used by the msgWriter to issue a carriage return when writing mail content.
    	SingleNewLine = "\r\n"
    
    	// DoubleNewLine represents a double newline character sequence ("\r\n\r\n").
    	//
    	// This constant can be used by the msgWriter to indicate a new segment of the mail when writing mail content.
    	DoubleNewLine = "\r\n\r\n"
    )

[View Source](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L171)
    
    
    const SendmailPath = "/usr/sbin/sendmail"

SendmailPath is the default system path to the sendmail binary - at least on standard Unix-like OS. 

[View Source](https://github.com/wneessen/go-mail/blob/v0.7.2/doc.go#L14)
    
    
    const VERSION = "0.7.2"

VERSION indicates the current version of the package. It is also attached to the default user agent string. 

### Variables ¶

[View Source](https://github.com/wneessen/go-mail/blob/v0.7.2/auth.go#L157)
    
    
    var (
    	// ErrPlainAuthNotSupported is returned when the server does not support the "PLAIN" SMTP
    	// authentication type.
    	ErrPlainAuthNotSupported = [errors](/errors).[New](/errors#New)("server does not support SMTP AUTH type: PLAIN")
    
    	// ErrLoginAuthNotSupported is returned when the server does not support the "LOGIN" SMTP
    	// authentication type.
    	ErrLoginAuthNotSupported = [errors](/errors).[New](/errors#New)("server does not support SMTP AUTH type: LOGIN")
    
    	// ErrCramMD5AuthNotSupported is returned when the server does not support the "CRAM-MD5" SMTP
    	// authentication type.
    	ErrCramMD5AuthNotSupported = [errors](/errors).[New](/errors#New)("server does not support SMTP AUTH type: CRAM-MD5")
    
    	// ErrXOauth2AuthNotSupported is returned when the server does not support the "XOAUTH2" schema.
    	ErrXOauth2AuthNotSupported = [errors](/errors).[New](/errors#New)("server does not support SMTP AUTH type: XOAUTH2")
    
    	// ErrSCRAMSHA1AuthNotSupported is returned when the server does not support the "SCRAM-SHA-1" SMTP
    	// authentication type.
    	ErrSCRAMSHA1AuthNotSupported = [errors](/errors).[New](/errors#New)("server does not support SMTP AUTH type: SCRAM-SHA-1")
    
    	// ErrSCRAMSHA1PLUSAuthNotSupported is returned when the server does not support the "SCRAM-SHA-1-PLUS" SMTP
    	// authentication type.
    	ErrSCRAMSHA1PLUSAuthNotSupported = [errors](/errors).[New](/errors#New)("server does not support SMTP AUTH type: SCRAM-SHA-1-PLUS")
    
    	// ErrSCRAMSHA256AuthNotSupported is returned when the server does not support the "SCRAM-SHA-256" SMTP
    	// authentication type.
    	ErrSCRAMSHA256AuthNotSupported = [errors](/errors).[New](/errors#New)("server does not support SMTP AUTH type: SCRAM-SHA-256")
    
    	// ErrSCRAMSHA256PLUSAuthNotSupported is returned when the server does not support the "SCRAM-SHA-256-PLUS" SMTP
    	// authentication type.
    	ErrSCRAMSHA256PLUSAuthNotSupported = [errors](/errors).[New](/errors#New)("server does not support SMTP AUTH type: SCRAM-SHA-256-PLUS")
    
    	// ErrNoSupportedAuthDiscovered is returned when the SMTP Auth AutoDiscover process fails to identify
    	// any supported authentication mechanisms offered by the server.
    	ErrNoSupportedAuthDiscovered = [errors](/errors).[New](/errors#New)("SMTP Auth autodiscover was not able to detect a supported " +
    		"authentication mechanism")
    )

SMTP Auth related static errors 

[View Source](https://github.com/wneessen/go-mail/blob/v0.7.2/client.go#L218)
    
    
    var (
    	// ErrInvalidPort is returned when the specified port for the SMTP connection is not valid
    	ErrInvalidPort = [errors](/errors).[New](/errors#New)("invalid port number")
    
    	// ErrInvalidTimeout is returned when the specified timeout is zero or negative.
    	ErrInvalidTimeout = [errors](/errors).[New](/errors#New)("timeout cannot be zero or negative")
    
    	// ErrInvalidHELO is returned when the HELO/EHLO value is invalid due to being empty.
    	ErrInvalidHELO = [errors](/errors).[New](/errors#New)("invalid HELO/EHLO value - must not be empty")
    
    	// ErrInvalidTLSConfig is returned when the provided TLS configuration is invalid or nil.
    	ErrInvalidTLSConfig = [errors](/errors).[New](/errors#New)("invalid TLS config")
    
    	// ErrNoHostname is returned when the hostname for the client is not provided or empty.
    	ErrNoHostname = [errors](/errors).[New](/errors#New)("hostname for client cannot be empty")
    
    	// ErrDeadlineExtendFailed is returned when an attempt to extend the connection deadline fails.
    	ErrDeadlineExtendFailed = [errors](/errors).[New](/errors#New)("connection deadline extension failed")
    
    	// ErrNoActiveConnection indicates that there is no active connection to the SMTP server.
    	ErrNoActiveConnection = [errors](/errors).[New](/errors#New)("not connected to SMTP server")
    
    	// ErrServerNoUnencoded indicates that the server does not support 8BITMIME for unencoded 8-bit messages.
    	ErrServerNoUnencoded = [errors](/errors).[New](/errors#New)("message is 8bit unencoded, but server does not support 8BITMIME")
    
    	// ErrInvalidDSNMailReturnOption is returned when an invalid DSNMailReturnOption is provided as argument
    	// to the WithDSN Option.
    	ErrInvalidDSNMailReturnOption = [errors](/errors).[New](/errors#New)("DSN mail return option can only be HDRS or FULL")
    
    	// ErrInvalidDSNRcptNotifyOption is returned when an invalid DSNRcptNotifyOption is provided as argument
    	// to the WithDSN Option.
    	ErrInvalidDSNRcptNotifyOption = [errors](/errors).[New](/errors#New)("DSN rcpt notify option can only be: NEVER, " +
    		"SUCCESS, FAILURE or DELAY")
    
    	// ErrInvalidDSNRcptNotifyCombination is returned when an invalid combination of DSNRcptNotifyOption is
    	// provided as argument to the WithDSN Option.
    	ErrInvalidDSNRcptNotifyCombination = [errors](/errors).[New](/errors#New)("DSN rcpt notify option NEVER cannot be " +
    		"combined with any of SUCCESS, FAILURE or DELAY")
    
    	// ErrSMTPAuthMethodIsNil indicates that the SMTP authentication method provided is nil
    	ErrSMTPAuthMethodIsNil = [errors](/errors).[New](/errors#New)("SMTP auth method is nil")
    
    	// ErrDialContextFuncIsNil indicates that a required dial context function is not provided.
    	ErrDialContextFuncIsNil = [errors](/errors).[New](/errors#New)("dial context function is nil")
    )

[View Source](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L32)
    
    
    var (
    	// ErrNoFromAddress indicates that the FROM address is not set, which is required.
    	ErrNoFromAddress = [errors](/errors).[New](/errors#New)("no FROM address set")
    
    	// ErrNoRcptAddresses indicates that no recipient addresses have been set.
    	ErrNoRcptAddresses = [errors](/errors).[New](/errors#New)("no recipient addresses set")
    )

[View Source](https://github.com/wneessen/go-mail/blob/v0.7.2/smime.go#L17)
    
    
    var (
    	// ErrPrivateKeyMissing should be used if private key is invalid
    	ErrPrivateKeyMissing = [errors](/errors).[New](/errors#New)("private key is missing")
    
    	// ErrCertificateMissing should be used if the certificate is invalid
    	ErrCertificateMissing = [errors](/errors).[New](/errors#New)("certificate is missing")
    )

### Functions ¶

####  func [IsAddrHeader](https://github.com/wneessen/go-mail/blob/v0.7.2/header.go#L233) ¶ added in v0.7.0
    
    
    func IsAddrHeader(header [string](/builtin#string)) [bool](/builtin#bool)

IsAddrHeader checks if the provided string is an address header. It returns true on a valid AddrHeader and false for any other string. 

Parameters: 

  * header: The string to check.



### Types ¶

####  type [AddrHeader](https://github.com/wneessen/go-mail/blob/v0.7.2/header.go#L13) ¶
    
    
    type AddrHeader [string](/builtin#string)

AddrHeader is a type wrapper for a string and represents email address headers fields in a Msg. 
    
    
    const (
    	// HeaderBcc is the "Blind Carbon Copy" header field.
    	HeaderBcc AddrHeader = "Bcc"
    
    	// HeaderCc is the "Carbon Copy" header field.
    	HeaderCc AddrHeader = "Cc"
    
    	// HeaderEnvelopeFrom is the envelope FROM header field.
    	//
    	// It is generally not included in the mail body but only used by the Client for the communication with the
    	// SMTP server. If the Msg has no "FROM" address set in the mail body, the msgWriter will try to use the
    	// envelope from address, if this has been set for the Msg.
    	HeaderEnvelopeFrom AddrHeader = "EnvelopeFrom"
    
    	// HeaderFrom is the "From" header field.
    	HeaderFrom AddrHeader = "From"
    
    	// HeaderReplyTo is the "Reply-To" header field.
    	HeaderReplyTo AddrHeader = "Reply-To"
    
    	// HeaderTo is the "Recipient" header field.
    	HeaderTo AddrHeader = "To"
    )

####  func (AddrHeader) [String](https://github.com/wneessen/go-mail/blob/v0.7.2/header.go#L224) ¶ added in v0.1.4
    
    
    func (a AddrHeader) String() [string](/builtin#string)

String satisfies the fmt.Stringer interface for the AddrHeader type and returns the string representation of the AddrHeader. 

Returns: 

  * A string representing the AddrHeader.



####  type [AuthData](https://github.com/wneessen/go-mail/blob/v0.7.2/quicksend.go#L15) ¶ added in v0.6.0
    
    
    type AuthData struct {
    	Auth     [bool](/builtin#bool)
    	Username [string](/builtin#string)
    	Password [string](/builtin#string)
    }

####  func [NewAuthData](https://github.com/wneessen/go-mail/blob/v0.7.2/quicksend.go#L105) ¶ added in v0.6.0
    
    
    func NewAuthData(user, pass [string](/builtin#string)) *AuthData

NewAuthData creates a new AuthData instance with the provided username and password. 

This function initializes an AuthData struct with authentication enabled and sets the username and password fields. 

Parameters: 

  * user: The username for authentication.
  * pass: The password for authentication.



Returns: 

  * A pointer to the initialized AuthData instance.



####  type [Charset](https://github.com/wneessen/go-mail/blob/v0.7.2/encoding.go#L8) ¶
    
    
    type Charset [string](/builtin#string)

Charset is a type wrapper for a string representing different character encodings. 
    
    
    const (
    	// CharsetUTF7 represents the "UTF-7" charset.
    	CharsetUTF7 Charset = "UTF-7"
    
    	// CharsetUTF8 represents the "UTF-8" charset.
    	CharsetUTF8 Charset = "UTF-8"
    
    	// CharsetASCII represents the "US-ASCII" charset.
    	CharsetASCII Charset = "US-ASCII"
    
    	// CharsetISO88591 represents the "ISO-8859-1" charset.
    	CharsetISO88591 Charset = "ISO-8859-1"
    
    	// CharsetISO88592 represents the "ISO-8859-2" charset.
    	CharsetISO88592 Charset = "ISO-8859-2"
    
    	// CharsetISO88593 represents the "ISO-8859-3" charset.
    	CharsetISO88593 Charset = "ISO-8859-3"
    
    	// CharsetISO88594 represents the "ISO-8859-4" charset.
    	CharsetISO88594 Charset = "ISO-8859-4"
    
    	// CharsetISO88595 represents the "ISO-8859-5" charset.
    	CharsetISO88595 Charset = "ISO-8859-5"
    
    	// CharsetISO88596 represents the "ISO-8859-6" charset.
    	CharsetISO88596 Charset = "ISO-8859-6"
    
    	// CharsetISO88597 represents the "ISO-8859-7" charset.
    	CharsetISO88597 Charset = "ISO-8859-7"
    
    	// CharsetISO88599 represents the "ISO-8859-9" charset.
    	CharsetISO88599 Charset = "ISO-8859-9"
    
    	// CharsetISO885913 represents the "ISO-8859-13" charset.
    	CharsetISO885913 Charset = "ISO-8859-13"
    
    	// CharsetISO885914 represents the "ISO-8859-14" charset.
    	CharsetISO885914 Charset = "ISO-8859-14"
    
    	// CharsetISO885915 represents the "ISO-8859-15" charset.
    	CharsetISO885915 Charset = "ISO-8859-15"
    
    	// CharsetISO885916 represents the "ISO-8859-16" charset.
    	CharsetISO885916 Charset = "ISO-8859-16"
    
    	// CharsetISO2022JP represents the "ISO-2022-JP" charset.
    	CharsetISO2022JP Charset = "ISO-2022-JP"
    
    	// CharsetISO2022KR represents the "ISO-2022-KR" charset.
    	CharsetISO2022KR Charset = "ISO-2022-KR"
    
    	// CharsetWindows1250 represents the "windows-1250" charset.
    	CharsetWindows1250 Charset = "windows-1250"
    
    	// CharsetWindows1251 represents the "windows-1251" charset.
    	CharsetWindows1251 Charset = "windows-1251"
    
    	// CharsetWindows1252 represents the "windows-1252" charset.
    	CharsetWindows1252 Charset = "windows-1252"
    
    	// CharsetWindows1255 represents the "windows-1255" charset.
    	CharsetWindows1255 Charset = "windows-1255"
    
    	// CharsetWindows1256 represents the "windows-1256" charset.
    	CharsetWindows1256 Charset = "windows-1256"
    
    	// CharsetKOI8R represents the "KOI8-R" charset.
    	CharsetKOI8R Charset = "KOI8-R"
    
    	// CharsetKOI8U represents the "KOI8-U" charset.
    	CharsetKOI8U Charset = "KOI8-U"
    
    	// CharsetBig5 represents the "Big5" charset.
    	CharsetBig5 Charset = "Big5"
    
    	// CharsetGB18030 represents the "GB18030" charset.
    	CharsetGB18030 Charset = "GB18030"
    
    	// CharsetGB2312 represents the "GB2312" charset.
    	CharsetGB2312 Charset = "GB2312"
    
    	// CharsetTIS620 represents the "TIS-620" charset.
    	CharsetTIS620 Charset = "TIS-620"
    
    	// CharsetEUCKR represents the "EUC-KR" charset.
    	CharsetEUCKR Charset = "EUC-KR"
    
    	// CharsetShiftJIS represents the "Shift_JIS" charset.
    	CharsetShiftJIS Charset = "Shift_JIS"
    
    	// CharsetUnknown represents the "Unknown" charset.
    	CharsetUnknown Charset = "Unknown"
    
    	// CharsetGBK represents the "GBK" charset.
    	CharsetGBK Charset = "GBK"
    )

####  func (Charset) [String](https://github.com/wneessen/go-mail/blob/v0.7.2/encoding.go#L201) ¶
    
    
    func (c Charset) String() [string](/builtin#string)

String satisfies the fmt.Stringer interface for the Charset type. It converts a Charset into a printable format. 

This method returns the string representation of the Charset, allowing it to be easily printed or logged. 

Returns: 

  * A string representation of the Charset.



####  type [Client](https://github.com/wneessen/go-mail/blob/v0.7.2/client.go#L116) ¶
    
    
    type Client struct {
    	// ErrorHandlerRegistry provides access to the smtp.Client's custom error handlers for SMTP
    	// host-command pairs which are based on the smtp.ResponseErrorHandler interface.
    	//
    	// The smtp.ResponseErrorHandler interface defines a method for handling SMTP responses that do not
    	// comply with expected formats or behaviors. It is useful for implementing retry logic, logging,
    	// or error handling logic for non-compliant SMTP responses.
    	ErrorHandlerRegistry *[smtp](/github.com/wneessen/go-mail@v0.7.2/smtp).[ErrorHandlerRegistry](/github.com/wneessen/go-mail@v0.7.2/smtp#ErrorHandlerRegistry)
    	// contains filtered or unexported fields
    }

Client is responsible for connecting and interacting with an SMTP server. 

This struct represents the go-mail client, which manages the connection, authentication, and communication with an SMTP server. It contains various configuration options, including connection timeouts, encryption settings, authentication methods, and Delivery Status Notification (DSN) preferences. 

References: 

  * <https://datatracker.ietf.org/doc/html/rfc3207#section-2>
  * <https://datatracker.ietf.org/doc/html/rfc8314>



####  func [NewClient](https://github.com/wneessen/go-mail/blob/v0.7.2/client.go#L279) ¶
    
    
    func NewClient(host [string](/builtin#string), opts ...Option) (*Client, [error](/builtin#error))

NewClient creates a new Client instance with the provided host and optional configuration Option functions. 

This function initializes a Client with default values, such as connection timeout, port, TLS settings, and the HELO/EHLO hostname. Option functions, if provided, can override the default configuration. It ensures that essential values, like the host, are set. The function also supports connections to UNIX domain sockets by recognizing a "unix://" prefix in the host string and adjusting the configuration accordingly. An error is returned if critical defaults are unset. 

Parameters: 

  * host: The hostname of the SMTP server to connect to, or a UNIX domain socket prefixed with "unix://".
  * opts: Optional configuration functions to override default settings.



Returns: 

  * A pointer to the initialized Client.
  * An error if any critical default values are missing or options fail to apply.

Example ¶

Code example for the NewClient method 
    
    
    c, err := mail.NewClient("mail.example.com")
    if err != nil {
    	panic(err)
    }
    _ = c
    

####  func (*Client) [Close](https://github.com/wneessen/go-mail/blob/v0.7.2/client.go#L1098) ¶
    
    
    func (c *Client) Close() [error](/builtin#error)

Close terminates the connection to the SMTP server, returning an error if the disconnection fails. If the connection is already closed, this method is a no-op and disregards any error. 

This function checks if the Client's SMTP connection is active. If not, it simply returns without any action. If the connection is active, it attempts to gracefully close the connection using the Quit method. 

Returns: 

  * An error if the disconnection fails; otherwise, returns nil.



####  func (*Client) [CloseWithSMTPClient](https://github.com/wneessen/go-mail/blob/v0.7.2/client.go#L1115) ¶ added in v0.6.0
    
    
    func (c *Client) CloseWithSMTPClient(client *[smtp](/github.com/wneessen/go-mail@v0.7.2/smtp).[Client](/github.com/wneessen/go-mail@v0.7.2/smtp#Client)) [error](/builtin#error)

CloseWithSMTPClient terminates the connection of the provided smtp.Client to the SMTP server, returning an error if the disconnection fails. If the connection is already closed, this method is a no-op and disregards any error. 

This function checks if the smtp.Client connection is active. If not, it simply returns without any action. If the connection is active, it attempts to gracefully close the connection using the Quit method. 

Parameters: 

  * client: A pointer to the smtp.Client that handles the connection to the server.



Returns: 

  * An error if the disconnection fails; otherwise, returns nil.



####  func (*Client) [DialAndSend](https://github.com/wneessen/go-mail/blob/v0.7.2/client.go#L1174) ¶
    
    
    func (c *Client) DialAndSend(messages ...*Msg) [error](/builtin#error)

DialAndSend establishes a connection to the server and sends out the provided Msg. It calls DialAndSendWithContext with an empty Context.Background. 

This method simplifies the process of connecting to the SMTP server and sending messages by using a default context. It prepares the messages for sending and ensures the connection is established before attempting to send them. 

Parameters: 

  * messages: A variadic list of pointers to Msg objects to be sent.



Returns: 

  * An error if the connection fails or if sending the messages fails; otherwise, returns nil.

Example ¶

Code example for the Client.DialAndSend method 
    
    
    from := "Toni Tester <toni@example.com>"
    to := "Alice <alice@example.com>"
    server := "mail.example.com"
    
    m := mail.NewMsg()
    if err := m.From(from); err != nil {
    	fmt.Printf("failed to set FROM address: %s", err)
    	os.Exit(1)
    }
    if err := m.To(to); err != nil {
    	fmt.Printf("failed to set TO address: %s", err)
    	os.Exit(1)
    }
    m.Subject("This is a great subject")
    
    c, err := mail.NewClient(server)
    if err != nil {
    	fmt.Printf("failed to create mail client: %s", err)
    	os.Exit(1)
    }
    if err := c.DialAndSend(m); err != nil {
    	fmt.Printf("failed to send mail: %s", err)
    	os.Exit(1)
    }
    

####  func (*Client) [DialAndSendWithContext](https://github.com/wneessen/go-mail/blob/v0.7.2/client.go#L1194) ¶ added in v0.2.8
    
    
    func (c *Client) DialAndSendWithContext(ctx [context](/context).[Context](/context#Context), messages ...*Msg) [error](/builtin#error)

DialAndSendWithContext establishes a connection to the SMTP server using DialWithContext with the provided context.Context, then sends out the given Msg. After successful delivery, the Client will close the connection to the server. 

This method first attempts to connect to the SMTP server using the provided context. Upon successful connection, it sends the specified messages and ensures that the connection is closed after the operation, regardless of success or failure in sending the messages. 

Parameters: 

  * ctx: The context.Context to control the connection timeout and cancellation.
  * messages: A variadic list of pointers to Msg objects to be sent.



Returns: 

  * An error if the connection fails, if sending the messages fails, or if closing the connection fails; otherwise, returns nil.



####  func (*Client) [DialToSMTPClientWithContext](https://github.com/wneessen/go-mail/blob/v0.7.2/client.go#L1026) ¶ added in v0.6.0
    
    
    func (c *Client) DialToSMTPClientWithContext(ctxDial [context](/context).[Context](/context#Context)) (*[smtp](/github.com/wneessen/go-mail@v0.7.2/smtp).[Client](/github.com/wneessen/go-mail@v0.7.2/smtp#Client), [error](/builtin#error))

DialToSMTPClientWithContext establishes and configures a smtp.Client connection using the provided context. 

This function uses the provided context to manage the connection deadline and cancellation. It dials the SMTP server using the Client's configured DialContextFunc or a default dialer. If SSL is enabled, it uses a TLS connection. After successfully connecting, it initializes an smtp.Client, sends the HELO/EHLO command, and optionally performs STARTTLS and SMTP AUTH based on the Client's configuration. Debug and authentication logging are enabled if configured. 

Parameters: 

  * ctxDial: The context used to control the connection timeout and cancellation.



Returns: 

  * A pointer to the initialized smtp.Client.
  * An error if the connection fails, the smtp.Client cannot be created, or any subsequent commands fail.



####  func (*Client) [DialWithContext](https://github.com/wneessen/go-mail/blob/v0.7.2/client.go#L999) ¶
    
    
    func (c *Client) DialWithContext(ctxDial [context](/context).[Context](/context#Context)) [error](/builtin#error)

DialWithContext establishes a connection to the server using the provided context.Context. 

This function adds a deadline based on the Client's timeout to the provided context.Context before connecting to the server. After dialing the defined DialContextFunc and successfully establishing the connection, it sends the HELO/EHLO SMTP command, followed by optional STARTTLS and SMTP AUTH commands. If debug logging is enabled, it attaches the log.Logger. 

After this method is called, the Client will have an active (cancelable) connection to the SMTP server. 

Parameters: 

  * ctxDial: The context.Context used to control the connection timeout and cancellation.



Returns: 

  * An error if the connection to the SMTP server fails or any subsequent command fails.



####  func (*Client) [Reset](https://github.com/wneessen/go-mail/blob/v0.7.2/client.go#L1135) ¶
    
    
    func (c *Client) Reset() [error](/builtin#error)

Reset sends an SMTP RSET command to reset the state of the current SMTP session. 

This method checks the connection to the SMTP server and, if the connection is valid, it sends an RSET command to reset the session state. If the connection is invalid or the command fails, an error is returned. 

Returns: 

  * An error if the connection check fails or if sending the RSET command fails; otherwise, returns nil.



####  func (*Client) [ResetWithSMTPClient](https://github.com/wneessen/go-mail/blob/v0.7.2/client.go#L1151) ¶ added in v0.6.0
    
    
    func (c *Client) ResetWithSMTPClient(client *[smtp](/github.com/wneessen/go-mail@v0.7.2/smtp).[Client](/github.com/wneessen/go-mail@v0.7.2/smtp#Client)) [error](/builtin#error)

ResetWithSMTPClient sends an SMTP RSET command to the provided smtp.Client, to reset the state of the current SMTP session. 

This method checks the connection to the SMTP server and, if the connection is valid, it sends an RSET command to reset the session state. If the connection is invalid or the command fails, an error is returned. 

Parameters: 

  * client: A pointer to the smtp.Client that handles the connection to the server.



Returns: 

  * An error if the connection check fails or if sending the RSET command fails; otherwise, returns nil.



####  func (*Client) [Send](https://github.com/wneessen/go-mail/blob/v0.7.2/client.go#L1230) ¶
    
    
    func (c *Client) Send(messages ...*Msg) (returnErr [error](/builtin#error))

Send attempts to send one or more Msg using the SMTP client that is assigned to the Client. If the Client has no active connection to the server, Send will fail with an error. For each of the provided Msg, it will associate a SendError with the Msg in case of a transmission or delivery error. 

This method first checks for an active connection to the SMTP server. If the connection is not valid, it returns a SendError. It then iterates over the provided messages, attempting to send each one. If an error occurs during sending, the method records the error and associates it with the corresponding Msg. If multiple errors are encountered, it aggregates them into a single SendError to be returned. 

Parameters: 

  * client: A pointer to the smtp.Client that holds the connection to the SMTP server
  * messages: A variadic list of pointers to Msg objects to be sent.



Returns: 

  * An error that represents the sending result, which may include multiple SendErrors if any occurred; otherwise, returns nil.



####  func (*Client) [SendWithSMTPClient](https://github.com/wneessen/go-mail/blob/v0.7.2/client.go#L1254) ¶ added in v0.6.0
    
    
    func (c *Client) SendWithSMTPClient(client *[smtp](/github.com/wneessen/go-mail@v0.7.2/smtp).[Client](/github.com/wneessen/go-mail@v0.7.2/smtp#Client), messages ...*Msg) (returnErr [error](/builtin#error))

SendWithSMTPClient attempts to send one or more Msg using a provided smtp.Client with an established connection to the SMTP server. If the smtp.Client has no active connection to the server, SendWithSMTPClient will fail with an error. For each of the provided Msg, it will associate a SendError with the Msg in case of a transmission or delivery error. 

This method first checks for an active connection to the SMTP server. If the connection is not valid, it returns a SendError. It then iterates over the provided messages, attempting to send each one. If an error occurs during sending, the method records the error and associates it with the corresponding Msg. If multiple errors are encountered, it aggregates them into a single SendError to be returned. 

Parameters: 

  * client: A pointer to the smtp.Client that holds the connection to the SMTP server
  * messages: A variadic list of pointers to Msg objects to be sent.



Returns: 

  * An error that represents the sending result, which may include multiple SendErrors if any occurred; otherwise, returns nil.



####  func (*Client) [ServerAddr](https://github.com/wneessen/go-mail/blob/v0.7.2/client.go#L758) ¶
    
    
    func (c *Client) ServerAddr() [string](/builtin#string)

ServerAddr returns the server address that is currently set on the Client in the format "host:port". 

This method constructs and returns the server address using the host and port currently configured for the Client. 

Returns: 

  * A string representing the server address in the format "host:port".



####  func (*Client) [SetDebugLog](https://github.com/wneessen/go-mail/blob/v0.7.2/client.go#L851) ¶ added in v0.3.9
    
    
    func (c *Client) SetDebugLog(val [bool](/builtin#bool))

SetDebugLog sets or overrides whether the Client is using debug logging. The debug logger will log incoming and outgoing communication between the client and the server to log.Logger that is defined on the Client. 

Note: The SMTP communication might include unencrypted authentication data, depending on whether you are using SMTP authentication and the type of authentication mechanism. This could pose a data protection risk. Use debug logging with caution. 

Parameters: 

  * val: A boolean value indicating whether to enable (true) or disable (false) debug logging.



####  func (*Client) [SetDebugLogWithSMTPClient](https://github.com/wneessen/go-mail/blob/v0.7.2/client.go#L866) ¶ added in v0.6.0
    
    
    func (c *Client) SetDebugLogWithSMTPClient(client *[smtp](/github.com/wneessen/go-mail@v0.7.2/smtp).[Client](/github.com/wneessen/go-mail@v0.7.2/smtp#Client), val [bool](/builtin#bool))

SetDebugLogWithSMTPClient sets or overrides whether the provided smtp.Client is using debug logging. The debug logger will log incoming and outgoing communication between the client and the server to log.Logger that is defined on the Client. 

Note: The SMTP communication might include unencrypted authentication data, depending on whether you are using SMTP authentication and the type of authentication mechanism. This could pose a data protection risk. Use debug logging with caution. 

Parameters: 

  * client: A pointer to the smtp.Client that handles the connection to the server.
  * val: A boolean value indicating whether to enable (true) or disable (false) debug logging.



####  func (*Client) [SetLogAuthData](https://github.com/wneessen/go-mail/blob/v0.7.2/client.go#L980) ¶ added in v0.5.1
    
    
    func (c *Client) SetLogAuthData(logAuth [bool](/builtin#bool))

SetLogAuthData sets or overrides the logging of SMTP authentication data for the Client. 

This function sets the logAuthData field of the Client to true, enabling the logging of authentication data. 

Be cautious when using this option, as the logs may include unencrypted authentication data, depending on the SMTP authentication method in use, which could pose a data protection risk. 

Parameters: 

  * logAuth: Set wether or not to log SMTP authentication data for the Client.



####  func (*Client) [SetLogger](https://github.com/wneessen/go-mail/blob/v0.7.2/client.go#L883) ¶ added in v0.3.9
    
    
    func (c *Client) SetLogger(logger [log](/github.com/wneessen/go-mail@v0.7.2/log).[Logger](/github.com/wneessen/go-mail@v0.7.2/log#Logger))

SetLogger sets or overrides the custom logger currently used by the Client. The logger must satisfy the log.Logger interface and is only utilized when debug logging is enabled on the Client. 

By default, log.Stdlog is used if no custom logger is provided. 

Parameters: 

  * logger: A logger that satisfies the log.Logger interface to be set for the Client.



####  func (*Client) [SetLoggerWithSMTPClient](https://github.com/wneessen/go-mail/blob/v0.7.2/client.go#L896) ¶ added in v0.6.0
    
    
    func (c *Client) SetLoggerWithSMTPClient(client *[smtp](/github.com/wneessen/go-mail@v0.7.2/smtp).[Client](/github.com/wneessen/go-mail@v0.7.2/smtp#Client), logger [log](/github.com/wneessen/go-mail@v0.7.2/log).[Logger](/github.com/wneessen/go-mail@v0.7.2/log#Logger))

SetLoggerWithSMTPClient sets or overrides the custom logger currently used by the provided smtp.Client. The logger must satisfy the log.Logger interface and is only utilized when debug logging is enabled on the provided smtp.Client. 

By default, log.Stdlog is used if no custom logger is provided. 

Parameters: 

  * client: A pointer to the smtp.Client that handles the connection to the server.
  * logger: A logger that satisfies the log.Logger interface to be set for the Client.



####  func (*Client) [SetPassword](https://github.com/wneessen/go-mail/blob/v0.7.2/client.go#L940) ¶
    
    
    func (c *Client) SetPassword(password [string](/builtin#string))

SetPassword sets or overrides the password that the Client will use for SMTP authentication. 

This method updates the password used by the Client for authenticating with the SMTP server. 

Parameters: 

  * password: The password to be set for SMTP authentication.



####  func (*Client) [SetSMTPAuth](https://github.com/wneessen/go-mail/blob/v0.7.2/client.go#L952) ¶
    
    
    func (c *Client) SetSMTPAuth(authtype SMTPAuthType)

SetSMTPAuth sets or overrides the SMTPAuthType currently configured on the Client for SMTP authentication. 

This method updates the authentication type used by the Client for authenticating with the SMTP server and resets any custom SMTP authentication mechanism. 

Parameters: 

  * authtype: The SMTPAuthType to be set for the Client.



####  func (*Client) [SetSMTPAuthCustom](https://github.com/wneessen/go-mail/blob/v0.7.2/client.go#L966) ¶
    
    
    func (c *Client) SetSMTPAuthCustom(smtpAuth [smtp](/github.com/wneessen/go-mail@v0.7.2/smtp).[Auth](/github.com/wneessen/go-mail@v0.7.2/smtp#Auth))

SetSMTPAuthCustom sets or overrides the custom SMTP authentication mechanism currently configured for the Client. The provided authentication mechanism must satisfy the smtp.Auth interface. 

This method updates the authentication mechanism used by the Client for authenticating with the SMTP server and sets the authentication type to SMTPAuthCustom. 

Parameters: 

  * smtpAuth: The custom SMTP authentication mechanism to be set for the Client.



####  func (*Client) [SetSSL](https://github.com/wneessen/go-mail/blob/v0.7.2/client.go#L810) ¶ added in v0.1.4
    
    
    func (c *Client) SetSSL(ssl [bool](/builtin#bool))

SetSSL sets or overrides whether the Client should use implicit SSL/TLS. 

This method configures the Client to either enable or disable implicit SSL/TLS for secure communication. 

Parameters: 

  * ssl: A boolean value indicating whether to enable (true) or disable (false) implicit SSL/TLS.



####  func (*Client) [SetSSLPort](https://github.com/wneessen/go-mail/blob/v0.7.2/client.go#L827) ¶ added in v0.4.1
    
    
    func (c *Client) SetSSLPort(ssl [bool](/builtin#bool), fallback [bool](/builtin#bool))

SetSSLPort sets or overrides whether the Client should use implicit SSL/TLS with optional fallback. The correct port is automatically set. 

If ssl is set to true, the default port 25 will be overridden with port 465. If fallback is set to true and the SSL/TLS connection fails, the Client will attempt to connect on port 25 using an unencrypted connection. 

Note: If a different port has already been set using WithPort, that port takes precedence and is used to establish the SSL/TLS connection, skipping the automatic fallback mechanism. 

Parameters: 

  * ssl: A boolean value indicating whether to enable implicit SSL/TLS.
  * fallback: A boolean value indicating whether to enable fallback to an unencrypted connection.



####  func (*Client) [SetTLSConfig](https://github.com/wneessen/go-mail/blob/v0.7.2/client.go#L916) ¶
    
    
    func (c *Client) SetTLSConfig(tlsconfig *[tls](/crypto/tls).[Config](/crypto/tls#Config)) [error](/builtin#error)

SetTLSConfig sets or overrides the tls.Config currently configured for the Client with the given value. An error is returned if the provided tls.Config is invalid. 

This method ensures that the provided tls.Config is not nil before updating the Client's TLS configuration. 

Parameters: 

  * tlsconfig: A pointer to the tls.Config struct to be set for the Client. Must not be nil.



Returns: 

  * An error if the provided tls.Config is invalid or nil.



####  func (*Client) [SetTLSPolicy](https://github.com/wneessen/go-mail/blob/v0.7.2/client.go#L772) ¶
    
    
    func (c *Client) SetTLSPolicy(policy TLSPolicy)

SetTLSPolicy sets or overrides the TLSPolicy currently configured on the Client with the given TLSPolicy. 

This method allows the user to set a new TLSPolicy for the Client. For best practices regarding SMTP TLS connections, it is recommended to use SetTLSPortPolicy instead. 

Parameters: 

  * policy: The TLSPolicy to be set for the Client.

Example ¶

Code example for the Client.SetTLSPolicy method 
    
    
    c, err := mail.NewClient("mail.example.com")
    if err != nil {
    	panic(err)
    }
    c.SetTLSPolicy(mail.TLSMandatory)
    fmt.Println(c.TLSPolicy())
    
    
    
    Output:
    
    TLSMandatory
    

####  func (*Client) [SetTLSPortPolicy](https://github.com/wneessen/go-mail/blob/v0.7.2/client.go#L788) ¶ added in v0.4.1
    
    
    func (c *Client) SetTLSPortPolicy(policy TLSPolicy)

SetTLSPortPolicy sets or overrides the TLSPolicy currently configured on the Client with the given TLSPolicy. The correct port is automatically set based on the specified policy. 

If TLSMandatory or TLSOpportunistic is provided as the TLSPolicy, port 587 will be used for the connection. If the connection fails with TLSOpportunistic, the Client will attempt to connect on port 25 using an unencrypted connection as a fallback. If NoTLS is provided, the Client will always use port 25. 

Note: If a different port has already been set using WithPort, that port takes precedence and is used to establish the SSL/TLS connection, skipping the automatic fallback mechanism. 

Parameters: 

  * policy: The TLSPolicy to be set for the Client.



####  func (*Client) [SetUsername](https://github.com/wneessen/go-mail/blob/v0.7.2/client.go#L930) ¶
    
    
    func (c *Client) SetUsername(username [string](/builtin#string))

SetUsername sets or overrides the username that the Client will use for SMTP authentication. 

This method updates the username used by the Client for authenticating with the SMTP server. 

Parameters: 

  * username: The username to be set for SMTP authentication.



####  func (*Client) [TLSPolicy](https://github.com/wneessen/go-mail/blob/v0.7.2/client.go#L747) ¶
    
    
    func (c *Client) TLSPolicy() [string](/builtin#string)

TLSPolicy returns the TLSPolicy that is currently set on the Client as a string. 

This method retrieves the current TLSPolicy configured for the Client and returns it as a string representation. 

Returns: 

  * A string representing the currently set TLSPolicy for the Client.



####  type [ContentType](https://github.com/wneessen/go-mail/blob/v0.7.2/encoding.go#L11) ¶
    
    
    type ContentType [string](/builtin#string)

ContentType is a type wrapper for a string and represents the MIME type of the content being handled. 
    
    
    const (
    	// TypeAppOctetStream represents the MIME type for arbitrary binary data.
    	TypeAppOctetStream ContentType = "application/octet-stream"
    
    	// TypeMultipartAlternative represents the MIME type for a message body that can contain multiple alternative
    	// formats.
    	TypeMultipartAlternative ContentType = "multipart/alternative"
    
    	// TypeMultipartMixed represents the MIME type for a multipart message containing different parts.
    	TypeMultipartMixed ContentType = "multipart/mixed"
    
    	// TypeMultipartRelated represents the MIME type for a multipart message where each part is a related file
    	// or resource.
    	TypeMultipartRelated ContentType = "multipart/related"
    
    	// TypePGPSignature represents the MIME type for PGP signed messages.
    	TypePGPSignature ContentType = "application/pgp-signature"
    
    	// TypePGPEncrypted represents the MIME type for PGP encrypted messages.
    	TypePGPEncrypted ContentType = "application/pgp-encrypted"
    
    	// TypeTextHTML represents the MIME type for HTML text content.
    	TypeTextHTML ContentType = "text/html"
    
    	// TypeTextPlain represents the MIME type for plain text content.
    	TypeTextPlain ContentType = "text/plain"
    
    	// TypeSMIMESigned represents the MIME type for S/MIME singed messages.
    	TypeSMIMESigned ContentType = `application/pkcs7-signature; name="smime.p7s"`
    )

####  func (ContentType) [String](https://github.com/wneessen/go-mail/blob/v0.7.2/encoding.go#L213) ¶ added in v0.4.2
    
    
    func (c ContentType) String() [string](/builtin#string)

String satisfies the fmt.Stringer interface for the ContentType type. It converts a ContentType into a printable format. 

This method returns the string representation of the ContentType, enabling its use in formatted output such as logging or displaying information to the user. 

Returns: 

  * A string representation of the ContentType.



####  type [DSNMailReturnOption](https://github.com/wneessen/go-mail/blob/v0.7.2/client.go#L96) ¶ added in v0.2.7
    
    
    type DSNMailReturnOption [string](/builtin#string)

DSNMailReturnOption is a type wrapper for a string and specifies the type of return content requested in a Delivery Status Notification (DSN). 

<https://datatracker.ietf.org/doc/html/rfc1891/>

####  type [DSNRcptNotifyOption](https://github.com/wneessen/go-mail/blob/v0.7.2/client.go#L102) ¶ added in v0.2.7
    
    
    type DSNRcptNotifyOption [string](/builtin#string)

DSNRcptNotifyOption is a type wrapper for a string and specifies the notification options for a recipient in DSNs. 

<https://datatracker.ietf.org/doc/html/rfc1891/>

####  type [DialContextFunc](https://github.com/wneessen/go-mail/blob/v0.7.2/client.go#L90) ¶ added in v0.4.0
    
    
    type DialContextFunc func(ctx [context](/context).[Context](/context#Context), network, address [string](/builtin#string)) ([net](/net).[Conn](/net#Conn), [error](/builtin#error))

DialContextFunc defines a function type for establishing a network connection using context, network type, and address. It is used to specify custom DialContext function. 

By default we use net.Dial or tls.Dial respectively. 

####  type [Encoding](https://github.com/wneessen/go-mail/blob/v0.7.2/encoding.go#L15) ¶
    
    
    type Encoding [string](/builtin#string)

Encoding is a type wrapper for a string and represents the type of encoding used for email messages and/or parts. 
    
    
    const (
    	// EncodingB64 represents the Base64 encoding as specified in [RFC 2045](https://rfc-editor.org/rfc/rfc2045.html).
    	//
    	// <https://datatracker.ietf.org/doc/html/rfc2045#section-6.8>
    	EncodingB64 Encoding = "base64"
    
    	// EncodingQP represents the "quoted-printable" encoding as specified in [RFC 2045](https://rfc-editor.org/rfc/rfc2045.html).
    	//
    	// <https://datatracker.ietf.org/doc/html/rfc2045#section-6.7>
    	EncodingQP Encoding = "quoted-printable"
    
    	// EncodingUSASCII represents encoding with only US-ASCII characters (aka 7Bit)
    	//
    	// <https://datatracker.ietf.org/doc/html/rfc2045#section-2.7>
    	EncodingUSASCII Encoding = "7bit"
    
    	// NoEncoding represents 8-bit encoding for email messages as specified in [RFC 6152](https://rfc-editor.org/rfc/rfc6152.html).
    	//
    	// <https://datatracker.ietf.org/doc/html/rfc2045#section-2.8>
    	//
    	// <https://datatracker.ietf.org/doc/html/rfc6152>
    	NoEncoding Encoding = "8bit"
    )

####  func (Encoding) [String](https://github.com/wneessen/go-mail/blob/v0.7.2/encoding.go#L225) ¶
    
    
    func (e Encoding) String() [string](/builtin#string)

String satisfies the fmt.Stringer interface for the Encoding type. It converts an Encoding into a printable format. 

This method returns the string representation of the Encoding, which can be used for displaying or logging purposes. 

Returns: 

  * A string representation of the Encoding.



####  type [File](https://github.com/wneessen/go-mail/blob/v0.7.2/file.go#L21) ¶ added in v0.1.1
    
    
    type File struct {
    	ContentType ContentType
    	Desc        [string](/builtin#string)
    	Enc         Encoding
    	Header      [textproto](/net/textproto).[MIMEHeader](/net/textproto#MIMEHeader)
    	Name        [string](/builtin#string)
    	Writer      func(w [io](/io).[Writer](/io#Writer)) ([int64](/builtin#int64), [error](/builtin#error))
    }

File represents a file with properties such as content type, description, encoding, headers, name, and a writer function. 

This struct can represent either an attachment or an embedded file in a Msg, and it stores relevant metadata such as content type and encoding, as well as a function to write the file's content to an io.Writer. 

####  type [FileOption](https://github.com/wneessen/go-mail/blob/v0.7.2/file.go#L13) ¶ added in v0.1.1
    
    
    type FileOption func(*File)

FileOption is a function type used to modify properties of a File 

####  func [WithFileContentID](https://github.com/wneessen/go-mail/blob/v0.7.2/file.go#L40) ¶ added in v0.4.2
    
    
    func WithFileContentID(id [string](/builtin#string)) FileOption

WithFileContentID sets the "Content-ID" header in the File's MIME headers to the specified ID. 

This function updates the File's MIME headers by setting the "Content-ID" to the provided string value, allowing the file to be referenced by this ID within the MIME structure. 

Parameters: 

  * id: A string representing the content ID to be set in the "Content-ID" header.



Returns: 

  * A FileOption function that updates the File's "Content-ID" header.



####  func [WithFileContentType](https://github.com/wneessen/go-mail/blob/v0.7.2/file.go#L112) ¶ added in v0.3.9
    
    
    func WithFileContentType(contentType ContentType) FileOption

WithFileContentType sets the content type of the File. 

By default, the content type is guessed based on the file type, and if no matching type is identified, the default "application/octet-stream" is used. This FileOption allows overriding the guessed content type with a specific one if required. 

Parameters: 

  * contentType: The ContentType to be assigned to the File.



Returns: 

  * A FileOption function that sets the File's content type.



####  func [WithFileDescription](https://github.com/wneessen/go-mail/blob/v0.7.2/file.go#L72) ¶ added in v0.3.9
    
    
    func WithFileDescription(description [string](/builtin#string)) FileOption

WithFileDescription sets an optional description for the File, which is used in the Content-Description header of the MIME output. 

This function updates the File's description, allowing an additional text description to be added to the MIME headers for the file. 

Parameters: 

  * description: A string representing the description to be set in the Content-Description header.



Returns: 

  * A FileOption function that sets the File's description.



####  func [WithFileEncoding](https://github.com/wneessen/go-mail/blob/v0.7.2/file.go#L92) ¶ added in v0.3.9
    
    
    func WithFileEncoding(encoding Encoding) FileOption

WithFileEncoding sets the encoding type for a File. 

This function allows the specification of an encoding type for the file, typically used for attachments or embedded files. By default, Base64 encoding should be used, but this function can override the default if needed. 

Note: Quoted-printable encoding (EncodingQP) must never be used for attachments or embeds. If EncodingQP is passed to this function, it will be ignored and the encoding will remain unchanged. 

Parameters: 

  * encoding: The Encoding type to be assigned to the File, unless it's EncodingQP.



Returns: 

  * A FileOption function that sets the File's encoding.



####  func [WithFileName](https://github.com/wneessen/go-mail/blob/v0.7.2/file.go#L55) ¶ added in v0.1.1
    
    
    func WithFileName(name [string](/builtin#string)) FileOption

WithFileName sets the name of a File to the provided value. 

This function assigns the specified name to the File, updating its Name field. 

Parameters: 

  * name: A string representing the name to be assigned to the File.



Returns: 

  * A FileOption function that sets the File's name.



####  type [Header](https://github.com/wneessen/go-mail/blob/v0.7.2/header.go#L10) ¶
    
    
    type Header [string](/builtin#string)

Header is a type wrapper for a string and represents email header fields in a Msg. 
    
    
    const (
    	// HeaderContentDescription is the "Content-Description" header.
    	HeaderContentDescription Header = "Content-Description"
    
    	// HeaderContentDisposition is the "Content-Disposition" header.
    	HeaderContentDisposition Header = "Content-Disposition"
    
    	// HeaderContentID is the "Content-ID" header.
    	HeaderContentID Header = "Content-ID"
    
    	// HeaderContentLang is the "Content-Language" header.
    	HeaderContentLang Header = "Content-Language"
    
    	// HeaderContentLocation is the "Content-Location" header ([RFC 2110](https://rfc-editor.org/rfc/rfc2110.html)).
    	// <https://datatracker.ietf.org/doc/html/rfc2110#section-4.3>
    	HeaderContentLocation Header = "Content-Location"
    
    	// HeaderContentTransferEnc is the "Content-Transfer-Encoding" header.
    	HeaderContentTransferEnc Header = "Content-Transfer-Encoding"
    
    	// HeaderContentType is the "Content-Type" header.
    	HeaderContentType Header = "Content-Type"
    
    	// HeaderDate represents the "Date" field.
    	// <https://datatracker.ietf.org/doc/html/rfc822#section-5.1>
    	HeaderDate Header = "Date"
    
    	// HeaderDispositionNotificationTo is the MDN header as described in [RFC 8098](https://rfc-editor.org/rfc/rfc8098.html).
    	// <https://datatracker.ietf.org/doc/html/rfc8098#section-2.1>
    	HeaderDispositionNotificationTo Header = "Disposition-Notification-To"
    
    	// HeaderImportance represents the "Importance" field.
    	HeaderImportance Header = "Importance"
    
    	// HeaderInReplyTo represents the "In-Reply-To" field.
    	HeaderInReplyTo Header = "In-Reply-To"
    
    	// HeaderListUnsubscribe is the "List-Unsubscribe" header field.
    	HeaderListUnsubscribe Header = "List-Unsubscribe"
    
    	// HeaderListUnsubscribePost is the "List-Unsubscribe-Post" header field.
    	HeaderListUnsubscribePost Header = "List-Unsubscribe-Post"
    
    	// HeaderMessageID represents the "Message-ID" field for message identification.
    	// <https://datatracker.ietf.org/doc/html/rfc1036#section-2.1.5>
    	HeaderMessageID Header = "Message-ID"
    
    	// HeaderMIMEVersion represents the "MIME-Version" field as per [RFC 2045](https://rfc-editor.org/rfc/rfc2045.html).
    	// <https://datatracker.ietf.org/doc/html/rfc2045#section-4>
    	HeaderMIMEVersion Header = "MIME-Version"
    
    	// HeaderOrganization is the "Organization" header field.
    	HeaderOrganization Header = "Organization"
    
    	// HeaderPrecedence is the "Precedence" header field.
    	HeaderPrecedence Header = "Precedence"
    
    	// HeaderPriority represents the "Priority" field.
    	HeaderPriority Header = "Priority"
    
    	// HeaderReferences is the "References" header field.
    	HeaderReferences Header = "References"
    
    	// HeaderSubject is the "Subject" header field.
    	HeaderSubject Header = "Subject"
    
    	// HeaderUserAgent is the "User-Agent" header field.
    	HeaderUserAgent Header = "User-Agent"
    
    	// HeaderXAutoResponseSuppress is the "X-Auto-Response-Suppress" header field.
    	HeaderXAutoResponseSuppress Header = "X-Auto-Response-Suppress"
    
    	// HeaderXMailer is the "X-Mailer" header field.
    	HeaderXMailer Header = "X-Mailer"
    
    	// HeaderXMSMailPriority is the "X-MSMail-Priority" header field.
    	HeaderXMSMailPriority Header = "X-MSMail-Priority"
    
    	// HeaderXPriority is the "X-Priority" header field.
    	HeaderXPriority Header = "X-Priority"
    )

####  func (Header) [String](https://github.com/wneessen/go-mail/blob/v0.7.2/header.go#L215) ¶ added in v0.1.4
    
    
    func (h Header) String() [string](/builtin#string)

String satisfies the fmt.Stringer interface for the Header type and returns the string representation of the Header. 

Returns: 

  * A string representing the Header.



####  type [Importance](https://github.com/wneessen/go-mail/blob/v0.7.2/header.go#L16) ¶
    
    
    type Importance [int](/builtin#int)

Importance is a type wrapper for an int and represents the level of importance or priority for a Msg. 
    
    
    const (
    	// ImportanceLow indicates a low level of importance or priority in a Msg.
    	ImportanceLow Importance = [iota](/builtin#iota)
    
    	// ImportanceNormal indicates a standard level of importance or priority for a Msg.
    	ImportanceNormal
    
    	// ImportanceHigh indicates a high level of importance or priority in a Msg.
    	ImportanceHigh
    
    	// ImportanceNonUrgent indicates a non-urgent level of importance or priority in a Msg.
    	ImportanceNonUrgent
    
    	// ImportanceUrgent indicates an urgent level of importance or priority in a Msg.
    	ImportanceUrgent
    )

####  func (Importance) [NumString](https://github.com/wneessen/go-mail/blob/v0.7.2/header.go#L149) ¶
    
    
    func (i Importance) NumString() [string](/builtin#string)

NumString returns a numerical string representation of the Importance level. 

This method maps ImportanceHigh and ImportanceUrgent to "1", while ImportanceNonUrgent and ImportanceLow are mapped to "0". Other values return an empty string. 

Returns: 

  * A string representing the numerical value of the Importance level ("1" or "0"), or an empty string if the Importance level is unrecognized.



####  func (Importance) [String](https://github.com/wneessen/go-mail/blob/v0.7.2/header.go#L195) ¶
    
    
    func (i Importance) String() [string](/builtin#string)

String satisfies the fmt.Stringer interface for the Importance type and returns the string representation of the Importance level. 

This method provides a human-readable string for each Importance level. 

Returns: 

  * A string representing the Importance level ("non-urgent", "low", "high", or "urgent"), or an empty string if the Importance level is unrecognized.



####  func (Importance) [XPrioString](https://github.com/wneessen/go-mail/blob/v0.7.2/header.go#L172) ¶
    
    
    func (i Importance) XPrioString() [string](/builtin#string)

XPrioString returns the X-Priority string representation of the Importance level. 

This method maps ImportanceHigh and ImportanceUrgent to "1", while ImportanceNonUrgent and ImportanceLow are mapped to "5". Other values return an empty string. 

Returns: 

  * A string representing the X-Priority value of the Importance level ("1" or "5"), or an empty string if the Importance level is unrecognized.



####  type [MIMEType](https://github.com/wneessen/go-mail/blob/v0.7.2/encoding.go#L21) ¶
    
    
    type MIMEType [string](/builtin#string)

MIMEType is a type wrapper for a string and represents the MIME type for the Msg content or parts. 
    
    
    const (
    	// MIMEAlternative MIMEType represents a MIME multipart/alternative type, used for emails with multiple versions.
    	MIMEAlternative MIMEType = "alternative"
    
    	// MIMEMixed MIMEType represents a MIME multipart/mixed type used fork emails containing different types of content.
    	MIMEMixed MIMEType = "mixed"
    
    	// MIMERelated MIMEType represents a MIME multipart/related type, used for emails with related content entities.
    	MIMERelated MIMEType = "related"
    
    	// MIMESMIMESigned MIMEType represents a MIME multipart/signed type, used for siging emails with S/MIME.
    	MIMESMIMESigned MIMEType = `signed; protocol="application/pkcs7-signature"; micalg=sha-256`
    )

####  func (MIMEType) [String](https://github.com/wneessen/go-mail/blob/v0.7.2/encoding.go#L237) ¶ added in v0.6.0
    
    
    func (e MIMEType) String() [string](/builtin#string)

String satisfies the fmt.Stringer interface for the MIMEType type. It converts an MIMEType into a printable format. 

This method returns the string representation of the MIMEType, which can be used for displaying or logging purposes. 

Returns: 

  * A string representation of the MIMEType.



####  type [MIMEVersion](https://github.com/wneessen/go-mail/blob/v0.7.2/encoding.go#L18) ¶
    
    
    type MIMEVersion [string](/builtin#string)

MIMEVersion is a type wrapper for a string nad represents the MIME version used in email messages. 
    
    
    const MIME10 MIMEVersion = "1.0"

MIME10 represents the MIME version "1.0" used in email messages. 

####  type [Middleware](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L77) ¶ added in v0.2.8
    
    
    type Middleware interface {
    	Handle(*Msg) *Msg
    	Type() MiddlewareType
    }

Middleware represents the interface for modifying or handling email messages. A Middleware allows the user to alter a Msg before it is finally processed. Multiple Middleware can be applied to a Msg. 

Type returns a unique MiddlewareType. It describes the type of Middleware and makes sure that a Middleware is only applied once. Handle performs all the processing to the Msg. It always needs to return a Msg back. 

####  type [MiddlewareType](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L69) ¶ added in v0.3.3
    
    
    type MiddlewareType [string](/builtin#string)

MiddlewareType is a type wrapper for a string. It describes the type of the Middleware and needs to be returned by the Middleware.Type method to satisfy the Middleware interface. 

####  type [Msg](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L89) ¶
    
    
    type Msg struct {
    	// contains filtered or unexported fields
    }

Msg represents an email message with various headers, attachments, and encoding settings. 

The Msg is the central part of go-mail. It provided a lot of methods that you would expect in a mail user agent (MUA). Msg satisfies the io.WriterTo and io.Reader interfaces. 

####  func [EMLToMsgFromFile](https://github.com/wneessen/go-mail/blob/v0.7.2/eml.go#L77) ¶ added in v0.4.2
    
    
    func EMLToMsgFromFile(filePath [string](/builtin#string)) (*Msg, [error](/builtin#error))

EMLToMsgFromFile opens and parses a .eml file at a provided file path and returns a pre-filled Msg pointer. 

This function attempts to read and parse an EML file located at the specified file path. It initializes a Msg object and populates it with the parsed headers and body. Any errors encountered during the file operations or parsing are returned. 

Parameters: 

  * filePath: The path to the .eml file to be parsed.



Returns: 

  * A pointer to the Msg object populated with the parsed data, and an error if parsing fails.



####  func [EMLToMsgFromReader](https://github.com/wneessen/go-mail/blob/v0.7.2/eml.go#L50) ¶ added in v0.4.2
    
    
    func EMLToMsgFromReader(reader [io](/io).[Reader](/io#Reader)) (*Msg, [error](/builtin#error))

EMLToMsgFromReader parses a reader that holds EML content and returns a pre-filled Msg pointer. 

This function reads EML content from the provided io.Reader and populates a Msg object with the parsed data. It initializes the Msg and extracts headers and body parts from the EML content. Any errors encountered during parsing are returned. 

Parameters: 

  * reader: An io.Reader containing the EML formatted message.



Returns: 

  * A pointer to the Msg object populated with the parsed data, and an error if parsing fails.



####  func [EMLToMsgFromString](https://github.com/wneessen/go-mail/blob/v0.7.2/eml.go#L33) ¶ added in v0.4.2
    
    
    func EMLToMsgFromString(emlString [string](/builtin#string)) (*Msg, [error](/builtin#error))

EMLToMsgFromString parses a given EML string and returns a pre-filled Msg pointer. 

This function takes an EML formatted string, converts it into a bytes buffer, and then calls EMLToMsgFromReader to parse the buffer and create a Msg object. This provides a convenient way to convert EML strings directly into Msg objects. 

Parameters: 

  * emlString: A string containing the EML formatted message.



Returns: 

  * A pointer to the Msg object populated with the parsed data, and an error if parsing fails.



####  func [NewMsg](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L192) ¶
    
    
    func NewMsg(opts ...MsgOption) *Msg

NewMsg creates a new email message with optional MsgOption functions that customize various aspects of the message. 

This function initializes a new Msg instance with default values for address headers, character set, encoding, general headers, and MIME version. It then applies any provided MsgOption functions to customize the message according to the user's needs. If an option is nil, it will be ignored. After applying the options, the function sets the appropriate MIME WordEncoder for the message. 

Parameters: 

  * opts: A variadic list of MsgOption functions that can be used to customize the Msg instance.



Returns: 

  * A pointer to the newly created Msg instance.



References: 

  * <https://datatracker.ietf.org/doc/html/rfc5321>

Example ¶

Code example for the NewMsg method 
    
    
    m := mail.NewMsg(mail.WithEncoding(mail.EncodingQP), mail.WithCharset(mail.CharsetASCII))
    fmt.Printf("%s // %s\n", m.Encoding(), m.Charset())
    
    
    
    Output:
    
    quoted-printable // US-ASCII
    

####  func [QuickSend](https://github.com/wneessen/go-mail/blob/v0.7.2/quicksend.go#L48) ¶ added in v0.6.0
    
    
    func QuickSend(addr [string](/builtin#string), auth *AuthData, from [string](/builtin#string), rcpts [][string](/builtin#string), subject [string](/builtin#string), content [][byte](/builtin#byte)) (*Msg, [error](/builtin#error))

QuickSend is an all-in-one method for quickly sending simple text mails in go-mail. 

This method will create a new client that connects to the server at addr, switches to TLS if possible, authenticates with the optional AuthData provided in auth and create a new simple Msg with the provided subject string and message bytes as body. The message will be sent using from as sender address and will be delivered to every address in rcpts. QuickSend will always send as text/plain ContentType. 

For the SMTP authentication, if auth is not nil and AuthData.Auth is set to true, it will try to autodiscover the best SMTP authentication mechanism supported by the server. If auth is set to true but autodiscover is not able to find a suitable authentication mechanism or if the authentication fails, the mail delivery will fail completely. 

The content parameter should be an [RFC 822](https://rfc-editor.org/rfc/rfc822.html)-style email body. The lines of content should be CRLF terminated. 

Parameters: 

  * addr: The hostname and port of the mail server, it must include a port, as in "mail.example.com:smtp".
  * auth: A AuthData pointer. If nil or if AuthData.Auth is set to false, not SMTP authentication will be performed.
  * from: The from address of the sender as string.
  * rcpts: A slice of strings of receipient addresses.
  * subject: The subject line as string.
  * content: A byte slice of the mail content



Returns: 

  * A pointer to the generated Msg.
  * An error if any step in the process of mail generation or delivery failed.



####  func (*Msg) [AddAlternativeHTMLTemplate](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L2011) ¶ added in v0.2.2
    
    
    func (m *Msg) AddAlternativeHTMLTemplate(tpl *[ht](/html/template).[Template](/html/template#Template), data interface{}, opts ...PartOption) [error](/builtin#error)

AddAlternativeHTMLTemplate sets the alternative body of the message to an html/template.Template output. 

The content type will be set to "text/html" automatically. This method executes the provided HTML template with the given data and adds the result as an alternative version of the message body. If the template is nil or fails to execute, an error will be returned. 

Parameters: 

  * tpl: A pointer to the html/template.Template to be used for the alternative body.
  * data: The data to populate the template.
  * opts: Optional parameters for customizing the alternative body part.



Returns: 

  * An error if the template is nil or fails to execute, otherwise nil.



References: 

  * <https://datatracker.ietf.org/doc/html/rfc2045>
  * <https://datatracker.ietf.org/doc/html/rfc2046>



####  func (*Msg) [AddAlternativeString](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L1964) ¶
    
    
    func (m *Msg) AddAlternativeString(contentType ContentType, content [string](/builtin#string), opts ...PartOption)

AddAlternativeString sets the alternative body of the message. 

This method adds an alternative representation of the message body using the specified content type and string content. This is typically used to provide both plain text and HTML versions of the email. Optional part settings can be provided via PartOption to further customize the message. 

Parameters: 

  * contentType: The content type of the alternative body (e.g., plain text, HTML).
  * content: The string content to set as the alternative body.
  * opts: Optional parameters for customizing the alternative body part.



References: 

  * <https://datatracker.ietf.org/doc/html/rfc2045>
  * <https://datatracker.ietf.org/doc/html/rfc2046>



####  func (*Msg) [AddAlternativeTextTemplate](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L2041) ¶ added in v0.2.2
    
    
    func (m *Msg) AddAlternativeTextTemplate(tpl *[tt](/text/template).[Template](/text/template#Template), data interface{}, opts ...PartOption) [error](/builtin#error)

AddAlternativeTextTemplate sets the alternative body of the message to a text/template.Template output. 

The content type will be set to "text/plain" automatically. This method executes the provided text template with the given data and adds the result as an alternative version of the message body. If the template is nil or fails to execute, an error will be returned. 

Parameters: 

  * tpl: A pointer to the text/template.Template to be used for the alternative body.
  * data: The data to populate the template.
  * opts: Optional parameters for customizing the alternative body part.



Returns: 

  * An error if the template is nil or fails to execute, otherwise nil.



References: 

  * <https://datatracker.ietf.org/doc/html/rfc2045>
  * <https://datatracker.ietf.org/doc/html/rfc2046>



####  func (*Msg) [AddAlternativeWriter](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L1985) ¶
    
    
    func (m *Msg) AddAlternativeWriter(
    	contentType ContentType, writeFunc func([io](/io).[Writer](/io#Writer)) ([int64](/builtin#int64), [error](/builtin#error)),
    	opts ...PartOption,
    )

AddAlternativeWriter sets the alternative body of the message. 

This method adds an alternative representation of the message body using a write function, allowing content to be written directly to the body. This is typically used to provide different formats, such as plain text and HTML. Optional part settings can be provided via PartOption to customize the message part. 

Parameters: 

  * contentType: The content type of the alternative body (e.g., plain text, HTML).
  * writeFunc: A function that writes content to an io.Writer and returns the number of bytes written and an error, if any.
  * opts: Optional parameters for customizing the alternative body part.



References: 

  * <https://datatracker.ietf.org/doc/html/rfc2045>
  * <https://datatracker.ietf.org/doc/html/rfc2046>



####  func (*Msg) [AddBcc](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L1074) ¶
    
    
    func (m *Msg) AddBcc(rcpt [string](/builtin#string)) [error](/builtin#error)

AddBcc adds a single "BCC" (blind carbon copy) address to the existing list of "BCC" recipients in the mail body for the Msg. 

This method allows you to add a single recipient to the "BCC" field without replacing any previously set "BCC" addresses. The "BCC" address specifies recipient(s) of the message who will receive a copy without other recipients being aware of it. The provided address is validated according to [RFC 5322](https://rfc-editor.org/rfc/rfc5322.html), and an error will be returned if the validation fails. 

Parameters: 

  * rcpt: The BCC address to add to the existing list of recipients in the Msg.



References: 

  * <https://datatracker.ietf.org/doc/html/rfc5322#section-3.6.3>



####  func (*Msg) [AddBccFormat](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L1108) ¶
    
    
    func (m *Msg) AddBccFormat(name, addr [string](/builtin#string)) [error](/builtin#error)

AddBccFormat adds a single "BCC" (blind carbon copy) address with the provided name and email to the existing list of "BCC" recipients in the mail body for the Msg. 

This method allows you to add a recipient's name and email address to the "BCC" field without replacing any previously set "BCC" addresses. The "BCC" address specifies recipient(s) of the message who will receive a copy without other recipients being aware of it. The provided name and address are validated according to [RFC 5322](https://rfc-editor.org/rfc/rfc5322.html), and an error will be returned if the validation fails. 

Parameters: 

  * name: The name of the recipient to add to the BCC field.
  * addr: The email address of the recipient to add to the BCC field.



References: 

  * <https://datatracker.ietf.org/doc/html/rfc5322#section-3.6.3>



####  func (*Msg) [AddBccMailAddress](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L1089) ¶ added in v0.7.0
    
    
    func (m *Msg) AddBccMailAddress(rcpt *[mail](/net/mail).[Address](/net/mail#Address))

AddBccMailAddress adds a single "BCC" address to the existing list of recipients in the mail body for the Msg. 

This method allows you to add a single recipient to the "BCC" field without replacing any previously set "BCC" addresses. The "BCC" address specifies recipient(s) of the message who will receive a copy without other recipients being aware of it. 

Parameters: 

  * rcpt: The recipient email address as *mail.Address to add to the "BCC" field.



References: 

  * <https://datatracker.ietf.org/doc/html/rfc5322#section-3.6.3>



####  func (*Msg) [AddCc](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L946) ¶
    
    
    func (m *Msg) AddCc(rcpt [string](/builtin#string)) [error](/builtin#error)

AddCc adds a single "CC" (carbon copy) address to the existing list of "CC" recipients in the mail body for the Msg. 

This method allows you to add a single recipient to the "CC" field without replacing any previously set "CC" addresses. The "CC" address specifies secondary recipient(s) and is visible to all recipients, including those in the "TO" field. The provided address is validated according to [RFC 5322](https://rfc-editor.org/rfc/rfc5322.html), and an error will be returned if the validation fails. 

Parameters: 

  * rcpt: The recipient address to be added to the "CC" field.



References: 

  * <https://datatracker.ietf.org/doc/html/rfc5322#section-3.6.3>



####  func (*Msg) [AddCcFormat](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L981) ¶
    
    
    func (m *Msg) AddCcFormat(name, addr [string](/builtin#string)) [error](/builtin#error)

AddCcFormat adds a single "CC" (carbon copy) address with the provided name and email to the existing list of "CC" recipients in the mail body for the Msg. 

This method allows you to add a recipient's name and email address to the "CC" field without replacing any previously set "CC" addresses. The "CC" address specifies secondary recipient(s) and is visible to all recipients, including those in the "TO" field. The provided name and address are validated according to [RFC 5322](https://rfc-editor.org/rfc/rfc5322.html), and an error will be returned if the validation fails. 

Parameters: 

  * name: The name of the recipient to be added to the "CC" field.
  * addr: The email address of the recipient to be added to the "CC" field.



References: 

  * <https://datatracker.ietf.org/doc/html/rfc5322#section-3.6.3>



####  func (*Msg) [AddCcMailAddress](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L962) ¶ added in v0.7.0
    
    
    func (m *Msg) AddCcMailAddress(rcpt *[mail](/net/mail).[Address](/net/mail#Address))

AddCcMailAddress adds a single "CC" address to the existing list of recipients in the mail body for the Msg. 

This method allows you to add a single recipient to the "CC" field without replacing any previously set "CC" addresses. The "CC" address specifies secondary recipient(s) and is visible to all recipients, including those in the "CC" field. Since the provided mail.Address has already been validated, no further validation is performed in this method and the values are used as given. 

Parameters: 

  * rcpt: The recipient email address as *mail.Address to add to the "CC" field.



References: 

  * <https://datatracker.ietf.org/doc/html/rfc5322#section-3.6.3>



####  func (*Msg) [AddTo](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L820) ¶
    
    
    func (m *Msg) AddTo(rcpt [string](/builtin#string)) [error](/builtin#error)

AddTo adds a single "TO" address to the existing list of recipients in the mail body for the Msg. 

This method allows you to add a single recipient to the "TO" field without replacing any previously set "TO" addresses. The "TO" address specifies the primary recipient(s) of the message and is visible in the mail client. The provided address is validated according to [RFC 5322](https://rfc-editor.org/rfc/rfc5322.html), and an error will be returned if the validation fails. 

Parameters: 

  * rcpt: The recipient email address to add to the "TO" field.



References: 

  * <https://datatracker.ietf.org/doc/html/rfc5322#section-3.6.3>



####  func (*Msg) [AddToFormat](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L855) ¶
    
    
    func (m *Msg) AddToFormat(name, addr [string](/builtin#string)) [error](/builtin#error)

AddToFormat adds a single "TO" address with the provided name and email to the existing list of recipients in the mail body for the Msg. 

This method allows you to add a recipient's name and email address to the "TO" field without replacing any previously set "TO" addresses. The "TO" address specifies the primary recipient(s) of the message and is visible in the mail client. The provided name and address are validated according to [RFC 5322](https://rfc-editor.org/rfc/rfc5322.html), and an error will be returned if the validation fails. 

Parameters: 

  * name: The name of the recipient to add to the "TO" field.
  * addr: The email address of the recipient to add to the "TO" field.



References: 

  * <https://datatracker.ietf.org/doc/html/rfc5322#section-3.6.3>



####  func (*Msg) [AddToMailAddress](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L836) ¶ added in v0.7.0
    
    
    func (m *Msg) AddToMailAddress(rcpt *[mail](/net/mail).[Address](/net/mail#Address))

AddToMailAddress adds a single "TO" address to the existing list of recipients in the mail body for the Msg. 

This method allows you to add a single recipient to the "TO" field without replacing any previously set "TO" addresses. The "TO" address specifies the primary recipient(s) of the message and is visible in the mail client. Since the provided mail.Address has already been validated, no further validation is performed in this method and the values are used as given. 

Parameters: 

  * rcpt: The recipient email address as *mail.Address to add to the "TO" field.



References: 

  * <https://datatracker.ietf.org/doc/html/rfc5322#section-3.6.3>



####  func (*Msg) [AttachFile](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L2066) ¶ added in v0.1.1
    
    
    func (m *Msg) AttachFile(name [string](/builtin#string), opts ...FileOption)

AttachFile adds an attachment File to the Msg. 

This method attaches a file to the message by specifying the file name. The file is retrieved from the filesystem and added to the list of attachments. Optional FileOption parameters can be provided to customize the attachment, such as setting its content type or encoding. 

Parameters: 

  * name: The name of the file to be attached.
  * opts: Optional parameters for customizing the attachment.



References: 

  * <https://datatracker.ietf.org/doc/html/rfc2183>



####  func (*Msg) [AttachFromEmbedFS](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L2189) ¶ added in v0.2.5
    
    
    func (m *Msg) AttachFromEmbedFS(name [string](/builtin#string), fs *[embed](/embed).[FS](/embed#FS), opts ...FileOption) [error](/builtin#error)

AttachFromEmbedFS adds an attachment File from an embed.FS to the Msg. 

This method allows you to attach a file from an embedded filesystem (embed.FS) to the message. The file is retrieved from the provided embed.FS and attached to the email. If the embedded filesystem is nil or the file cannot be retrieved, an error will be returned. 

Parameters: 

  * name: The name of the file to be attached.
  * fs: A pointer to the embed.FS from which the file will be retrieved.
  * opts: Optional parameters for customizing the attachment.



Returns: 

  * An error if the embed.FS is nil or the file cannot be retrieved, otherwise nil.



References: 

  * <https://datatracker.ietf.org/doc/html/rfc2183>



####  func (*Msg) [AttachFromIOFS](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L2208) ¶ added in v0.6.0
    
    
    func (m *Msg) AttachFromIOFS(name [string](/builtin#string), iofs [fs](/io/fs).[FS](/io/fs#FS), opts ...FileOption) [error](/builtin#error)

AttachFromIOFS attaches a file from a generic file system to the message. 

This function retrieves a file by name from an fs.FS instance, processes it, and appends it to the message's attachment collection. Additional file options can be provided for further customization. 

Parameters: 

  * name: The name of the file to retrieve from the file system.
  * iofs: The file system (must not be nil).
  * opts: Optional file options to customize the attachment process.



Returns: 

  * An error if the file cannot be retrieved, the fs.FS is nil, or any other issue occurs.



####  func (*Msg) [AttachHTMLTemplate](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L2134) ¶ added in v0.2.2
    
    
    func (m *Msg) AttachHTMLTemplate(
    	name [string](/builtin#string), tpl *[ht](/html/template).[Template](/html/template#Template), data interface{}, opts ...FileOption,
    ) [error](/builtin#error)

AttachHTMLTemplate adds the output of a html/template.Template pointer as a File attachment to the Msg. 

This method allows you to attach the rendered output of an HTML template as a file to the message. The template is executed with the provided data, and its output is attached as a file. If the template fails to execute, an error will be returned. 

Parameters: 

  * name: The name of the file to be attached.
  * tpl: A pointer to the html/template.Template to be executed for the attachment.
  * data: The data to populate the template.
  * opts: Optional parameters for customizing the attachment.



Returns: 

  * An error if the template fails to execute or cannot be attached, otherwise nil.



References: 

  * <https://datatracker.ietf.org/doc/html/rfc2183>



####  func (*Msg) [AttachReadSeeker](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L2112) ¶ added in v0.3.9
    
    
    func (m *Msg) AttachReadSeeker(name [string](/builtin#string), reader [io](/io).[ReadSeeker](/io#ReadSeeker), opts ...FileOption)

AttachReadSeeker adds an attachment File via io.ReadSeeker to the Msg. 

This method allows you to attach a file to the message using an io.ReadSeeker, which is more efficient for larger files compared to AttachReader, as it allows for seeking through the data without needing to load the entire content into memory. 

Parameters: 

  * name: The name of the file to be attached.
  * reader: The io.ReadSeeker providing the file data to be attached.
  * opts: Optional parameters for customizing the attachment.



References: 

  * <https://datatracker.ietf.org/doc/html/rfc2183>



####  func (*Msg) [AttachReader](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L2090) ¶ added in v0.1.1
    
    
    func (m *Msg) AttachReader(name [string](/builtin#string), reader [io](/io).[Reader](/io#Reader), opts ...FileOption) [error](/builtin#error)

AttachReader adds an attachment File via io.Reader to the Msg. 

This method allows you to attach a file to the message using an io.Reader. It reads all data from the io.Reader into memory before attaching the file, which may not be suitable for large data sources. For larger files, it is recommended to use AttachFile or AttachReadSeeker instead. 

Parameters: 

  * name: The name of the file to be attached.
  * reader: The io.Reader providing the file data to be attached.
  * opts: Optional parameters for customizing the attachment.



Returns: 

  * An error if the file could not be read from the io.Reader, otherwise nil.



References: 

  * <https://datatracker.ietf.org/doc/html/rfc2183>



####  func (*Msg) [AttachTextTemplate](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L2162) ¶ added in v0.2.2
    
    
    func (m *Msg) AttachTextTemplate(
    	name [string](/builtin#string), tpl *[tt](/text/template).[Template](/text/template#Template), data interface{}, opts ...FileOption,
    ) [error](/builtin#error)

AttachTextTemplate adds the output of a text/template.Template pointer as a File attachment to the Msg. 

This method allows you to attach the rendered output of a text template as a file to the message. The template is executed with the provided data, and its output is attached as a file. If the template fails to execute, an error will be returned. 

Parameters: 

  * name: The name of the file to be attached.
  * tpl: A pointer to the text/template.Template to be executed for the attachment.
  * data: The data to populate the template.
  * opts: Optional parameters for customizing the attachment.



Returns: 

  * An error if the template fails to execute or cannot be attached, otherwise nil.



References: 

  * <https://datatracker.ietf.org/doc/html/rfc2183>



####  func (*Msg) [Bcc](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L1041) ¶
    
    
    func (m *Msg) Bcc(rcpts ...[string](/builtin#string)) [error](/builtin#error)

Bcc sets one or more "BCC" (blind carbon copy) addresses in the mail body for the Msg. 

The "BCC" address specifies recipient(s) of the message who will receive a copy without other recipients being aware of it. These addresses are not visible in the mail body or to any other recipients, ensuring the privacy of BCC'd recipients. Multiple "BCC" addresses can be set by passing them as variadic arguments to this method. Each provided address is validated according to [RFC 5322](https://rfc-editor.org/rfc/rfc5322.html), and an error will be returned if ANY validation fails. 

Parameters: 

  * rcpts: One or more string values representing the BCC addresses to set in the Msg.



References: 

  * <https://datatracker.ietf.org/doc/html/rfc5322#section-3.6.3>



####  func (*Msg) [BccFromString](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L1142) ¶ added in v0.4.1
    
    
    func (m *Msg) BccFromString(rcpts [string](/builtin#string)) [error](/builtin#error)

BccFromString takes a string of comma-separated email addresses, validates each, and sets them as the "BCC" addresses for the Msg. 

This method allows you to pass a single string containing multiple email addresses separated by commas. Each address is validated according to [RFC 5322](https://rfc-editor.org/rfc/rfc5322.html) and set as a recipient in the "BCC" field. If any validation fails, an error will be returned. The addresses are not visible in the mail body and ensure the privacy of BCC'd recipients. 

Parameters: 

  * rcpts: A string of comma-separated email addresses to set as BCC recipients.



References: 

  * <https://datatracker.ietf.org/doc/html/rfc5322#section-3.6.3>



####  func (*Msg) [BccIgnoreInvalid](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L1125) ¶
    
    
    func (m *Msg) BccIgnoreInvalid(rcpts ...[string](/builtin#string))

BccIgnoreInvalid sets one or more "BCC" (blind carbon copy) addresses in the mail body for the Msg, ignoring any invalid addresses. 

This method allows you to add multiple "BCC" recipients to the message body. Unlike the standard `Bcc` method, any invalid addresses are ignored, and no error is returned for those addresses. Valid addresses will still be included in the "BCC" field, which ensures the privacy of the BCC'd recipients. Use this method with caution if address validation is critical, as invalid addresses are determined according to [RFC 5322](https://rfc-editor.org/rfc/rfc5322.html). 

Parameters: 

  * rcpts: One or more string values representing the BCC email addresses to set.



References: 

  * <https://datatracker.ietf.org/doc/html/rfc5322#section-3.6.3>



####  func (*Msg) [BccMailAddress](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L1057) ¶ added in v0.7.0
    
    
    func (m *Msg) BccMailAddress(rcpts ...*[mail](/net/mail).[Address](/net/mail#Address))

BccMailAddress sets one or more "BCC" (blind carbon copy) addresses in the mail body for the Msg. 

The "BCC" address specifies recipient(s) of the message who will receive a copy without other recipients being aware of it. These addresses are not visible in the mail body or to any other recipients, ensuring the privacy of BCC'd recipients. Multiple "BCC" addresses can be set by passing them as variadic arguments arguments to this method. 

Parameters: 

  * rcpts: One or more recipient email addresses as mail.Address instance to include in the "BCC" field.



References: 

  * <https://datatracker.ietf.org/doc/html/rfc5322#section-3.6.3>



####  func (*Msg) [Cc](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L913) ¶
    
    
    func (m *Msg) Cc(rcpts ...[string](/builtin#string)) [error](/builtin#error)

Cc sets one or more "CC" (carbon copy) addresses in the mail body for the Msg. 

The "CC" address specifies secondary recipient(s) of the message, and is included in the mail body. These addresses are visible to all recipients, including those listed in the "TO" and other "CC" fields. Multiple "CC" addresses can be set by passing them as variadic arguments to this method. Each provided address is validated according to [RFC 5322](https://rfc-editor.org/rfc/rfc5322.html), and an error will be returned if ANY validation fails. 

Parameters: 

  * rcpts: One or more recipient addresses to be included in the "CC" field.



References: 

  * <https://datatracker.ietf.org/doc/html/rfc5322#section-3.6.3>



####  func (*Msg) [CcFromString](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L1015) ¶ added in v0.4.1
    
    
    func (m *Msg) CcFromString(rcpts [string](/builtin#string)) [error](/builtin#error)

CcFromString takes a string of comma-separated email addresses, validates each, and sets them as the "CC" addresses for the Msg. 

This method allows you to pass a single string containing multiple email addresses separated by commas. Each address is validated according to [RFC 5322](https://rfc-editor.org/rfc/rfc5322.html) and set as a recipient in the "CC" field. If any validation fails, an error will be returned. The addresses are visible in the mail body and displayed to recipients in the mail client. 

Parameters: 

  * rcpts: A string containing multiple email addresses separated by commas.



References: 

  * <https://datatracker.ietf.org/doc/html/rfc5322#section-3.6.3>



####  func (*Msg) [CcIgnoreInvalid](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L998) ¶
    
    
    func (m *Msg) CcIgnoreInvalid(rcpts ...[string](/builtin#string))

CcIgnoreInvalid sets one or more "CC" (carbon copy) addresses in the mail body for the Msg, ignoring any invalid addresses. 

This method allows you to add multiple "CC" recipients to the message body. Unlike the standard `Cc` method, any invalid addresses are ignored, and no error is returned for those addresses. Valid addresses will still be included in the "CC" field, which is visible to all recipients in the mail client. Use this method with caution if address validation is critical, as invalid addresses are determined according to [RFC 5322](https://rfc-editor.org/rfc/rfc5322.html). 

Parameters: 

  * rcpts: One or more recipient email addresses to be added to the "CC" field.



References: 

  * <https://datatracker.ietf.org/doc/html/rfc5322#section-3.6.3>



####  func (*Msg) [CcMailAddress](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L929) ¶ added in v0.7.0
    
    
    func (m *Msg) CcMailAddress(rcpts ...*[mail](/net/mail).[Address](/net/mail#Address))

CcMailAddress sets one or more "CC" (carbon copy) addresses in the mail body for the Msg. 

The "CC" address specifies secondary recipient(s) of the message, and is included in the mail body. This address is visible to the recipient and any other recipients of the message. Multiple "CC" addresses can be set by passing them as variadic arguments to this method. 

Parameters: 

  * rcpts: One or more recipient email addresses as mail.Address instance to include in the "CC" field.



References: 

  * <https://datatracker.ietf.org/doc/html/rfc5322#section-3.6.3>



####  func (*Msg) [Charset](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L462) ¶
    
    
    func (m *Msg) Charset() [string](/builtin#string)

Charset returns the currently set Charset of the Msg as a string. 

This method retrieves the character set that is currently applied to the message. The charset defines the encoding for the text content of the message, ensuring that characters are displayed correctly across different email clients and platforms. The returned string will reflect the specific charset in use, such as UTF-8 or ISO-8859-1. 

Returns: 

  * A string representation of the current Charset of the Msg.



####  func (*Msg) [EmbedFile](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L2232) ¶ added in v0.1.1
    
    
    func (m *Msg) EmbedFile(name [string](/builtin#string), opts ...FileOption)

EmbedFile adds an embedded File to the Msg. 

This method embeds a file from the filesystem directly into the email message. The embedded file, typically an image or media file, can be referenced within the email's content (such as inline in HTML). If the file is not found or cannot be loaded, it will not be added. 

Parameters: 

  * name: The name of the file to be embedded.
  * opts: Optional parameters for customizing the embedded file.



References: 

  * <https://datatracker.ietf.org/doc/html/rfc2183>



####  func (*Msg) [EmbedFromEmbedFS](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L2354) ¶ added in v0.2.5
    
    
    func (m *Msg) EmbedFromEmbedFS(name [string](/builtin#string), fs *[embed](/embed).[FS](/embed#FS), opts ...FileOption) [error](/builtin#error)

EmbedFromEmbedFS adds an embedded File from an embed.FS to the Msg. 

This method embeds a file from an embedded filesystem (embed.FS) into the email message. If the embedded filesystem is nil or the file cannot be retrieved, an error will be returned. 

Parameters: 

  * name: The name of the file to be embedded.
  * fs: A pointer to the embed.FS from which the file will be retrieved.
  * opts: Optional parameters for customizing the embedded file.



Returns: 

  * An error if the embed.FS is nil or the file cannot be retrieved, otherwise nil.



References: 

  * <https://datatracker.ietf.org/doc/html/rfc2183>



####  func (*Msg) [EmbedFromIOFS](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L2373) ¶ added in v0.6.0
    
    
    func (m *Msg) EmbedFromIOFS(name [string](/builtin#string), iofs [fs](/io/fs).[FS](/io/fs#FS), opts ...FileOption) [error](/builtin#error)

EmbedFromIOFS embeds a file from a generic file system into the message. 

This function retrieves a file by name from an fs.FS instance, processes it, and appends it to the message's embed collection. Additional file options can be provided for further customization. 

Parameters: 

  * name: The name of the file to retrieve from the file system.
  * iofs: The file system (must not be nil).
  * opts: Optional file options to customize the embedding process.



Returns: 

  * An error if the file cannot be retrieved, the fs.FS is nil, or any other issue occurs.



####  func (*Msg) [EmbedHTMLTemplate](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L2300) ¶ added in v0.2.2
    
    
    func (m *Msg) EmbedHTMLTemplate(
    	name [string](/builtin#string), tpl *[ht](/html/template).[Template](/html/template#Template), data interface{}, opts ...FileOption,
    ) [error](/builtin#error)

EmbedHTMLTemplate adds the output of a html/template.Template pointer as an embedded File to the Msg. 

This method embeds the rendered output of an HTML template into the email message. The template is executed with the provided data, and its output is embedded as a file. If the template fails to execute, an error will be returned. 

Parameters: 

  * name: The name of the embedded file.
  * tpl: A pointer to the html/template.Template to be executed for the embedded content.
  * data: The data to populate the template.
  * opts: Optional parameters for customizing the embedded file.



Returns: 

  * An error if the template fails to execute or cannot be embedded, otherwise nil.



References: 

  * <https://datatracker.ietf.org/doc/html/rfc2183>



####  func (*Msg) [EmbedReadSeeker](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L2278) ¶ added in v0.3.9
    
    
    func (m *Msg) EmbedReadSeeker(name [string](/builtin#string), reader [io](/io).[ReadSeeker](/io#ReadSeeker), opts ...FileOption)

EmbedReadSeeker adds an embedded File from an io.ReadSeeker to the Msg. 

This method embeds a file into the email message by reading its content from an io.ReadSeeker. Using io.ReadSeeker allows for more efficient handling of large files since it can seek through the data without loading the entire content into memory. 

Parameters: 

  * name: The name of the file to be embedded.
  * reader: The io.ReadSeeker providing the file data to be embedded.
  * opts: Optional parameters for customizing the embedded file.



References: 

  * <https://datatracker.ietf.org/doc/html/rfc2183>



####  func (*Msg) [EmbedReader](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L2256) ¶ added in v0.1.1
    
    
    func (m *Msg) EmbedReader(name [string](/builtin#string), reader [io](/io).[Reader](/io#Reader), opts ...FileOption) [error](/builtin#error)

EmbedReader adds an embedded File from an io.Reader to the Msg. 

This method embeds a file into the email message by reading its content from an io.Reader. It reads all data into memory before embedding the file, which may not be efficient for large data sources. For larger files, it is recommended to use EmbedFile or EmbedReadSeeker instead. 

Parameters: 

  * name: The name of the file to be embedded.
  * reader: The io.Reader providing the file data to be embedded.
  * opts: Optional parameters for customizing the embedded file.



Returns: 

  * An error if the file could not be read from the io.Reader, otherwise nil.



References: 

  * <https://datatracker.ietf.org/doc/html/rfc2183>



####  func (*Msg) [EmbedTextTemplate](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L2328) ¶ added in v0.2.2
    
    
    func (m *Msg) EmbedTextTemplate(
    	name [string](/builtin#string), tpl *[tt](/text/template).[Template](/text/template#Template), data interface{}, opts ...FileOption,
    ) [error](/builtin#error)

EmbedTextTemplate adds the output of a text/template.Template pointer as an embedded File to the Msg. 

This method embeds the rendered output of a text template into the email message. The template is executed with the provided data, and its output is embedded as a file. If the template fails to execute, an error will be returned. 

Parameters: 

  * name: The name of the embedded file.
  * tpl: A pointer to the text/template.Template to be executed for the embedded content.
  * data: The data to populate the template.
  * opts: Optional parameters for customizing the embedded file.



Returns: 

  * An error if the template fails to execute or cannot be embedded, otherwise nil.



References: 

  * <https://datatracker.ietf.org/doc/html/rfc2183>



####  func (*Msg) [Encoding](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L449) ¶
    
    
    func (m *Msg) Encoding() [string](/builtin#string)

Encoding returns the currently set Encoding of the Msg as a string. 

This method retrieves the encoding type that is currently applied to the message. The encoding type determines how the message content is encoded for transmission. Common encoding types include quoted-printable and base64, and the returned string will reflect the specific encoding method in use. 

Returns: 

  * A string representation of the current Encoding of the Msg.



####  func (*Msg) [EnvelopeFrom](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L686) ¶ added in v0.2.4
    
    
    func (m *Msg) EnvelopeFrom(from [string](/builtin#string)) [error](/builtin#error)

EnvelopeFrom sets the envelope from address for the Msg. 

The HeaderEnvelopeFrom address is generally not included in the mail body but only used by the Client for communication with the SMTP server. If the Msg has no "FROM" address set in the mail body, the msgWriter will try to use the envelope from address if it has been set for the Msg. The provided address is validated according to [RFC 5322](https://rfc-editor.org/rfc/rfc5322.html) and will return an error if the validation fails. 

Parameters: 

  * from: The envelope from address to set in the Msg.



References: 

  * <https://datatracker.ietf.org/doc/html/rfc5322#section-3.4>



####  func (*Msg) [EnvelopeFromFormat](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L704) ¶ added in v0.2.4
    
    
    func (m *Msg) EnvelopeFromFormat(name, addr [string](/builtin#string)) [error](/builtin#error)

EnvelopeFromFormat sets the provided name and mail address as HeaderEnvelopeFrom for the Msg. 

The HeaderEnvelopeFrom address is generally not included in the mail body but only used by the Client for communication with the SMTP server. If the Msg has no "FROM" address set in the mail body, the msgWriter will try to use the envelope from address if it has been set for the Msg. The provided name and address are validated according to [RFC 5322](https://rfc-editor.org/rfc/rfc5322.html) and will return an error if the validation fails. 

Parameters: 

  * name: The name to associate with the envelope from address.
  * addr: The mail address to set as the envelope from address.



References: 

  * <https://datatracker.ietf.org/doc/html/rfc5322#section-3.4>



####  func (*Msg) [EnvelopeFromMailAddress](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L721) ¶ added in v0.7.0
    
    
    func (m *Msg) EnvelopeFromMailAddress(addr *[mail](/net/mail).[Address](/net/mail#Address))

EnvelopeFromMailAddress sets the "FROM" address in the mail body for the Msg using a mail.Address instance. 

The HeaderEnvelopeFrom address is generally not included in the mail body but only used by the Client for communication with the SMTP server. If the Msg has no "FROM" address set in the mail body, the msgWriter will try to use the envelope from address if it has been set for the Msg. The provided name and address are validated according to [RFC 5322](https://rfc-editor.org/rfc/rfc5322.html) and will return an error if the validation fails. 

Parameters: 

  * addr: The address as mail.Address instance to be set as envelope from address.



References: 

  * <https://datatracker.ietf.org/doc/html/rfc5322#section-3.6.2>



####  func (*Msg) [From](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L738) ¶
    
    
    func (m *Msg) From(from [string](/builtin#string)) [error](/builtin#error)

From sets the "FROM" address in the mail body for the Msg. 

The "FROM" address is included in the mail body and indicates the sender of the message to the recipient. This address is visible in the email client and is typically displayed to the recipient. If the "FROM" address is not set, the msgWriter may attempt to use the envelope from address (if available) for sending. The provided address is validated according to [RFC 5322](https://rfc-editor.org/rfc/rfc5322.html) and will return an error if the validation fails. 

Parameters: 

  * from: The "FROM" address to set in the mail body.



References: 

  * <https://datatracker.ietf.org/doc/html/rfc5322#section-3.6.2>



####  func (*Msg) [FromFormat](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L772) ¶
    
    
    func (m *Msg) FromFormat(name, addr [string](/builtin#string)) [error](/builtin#error)

FromFormat sets the provided name and mail address as the "FROM" address in the mail body for the Msg. 

The "FROM" address is included in the mail body and indicates the sender of the message to the recipient, and is visible in the email client. If the "FROM" address is not explicitly set, the msgWriter may use the envelope from address (if provided) when sending the message. The provided name and address are validated according to [RFC 5322](https://rfc-editor.org/rfc/rfc5322.html) and will return an error if the validation fails. 

Parameters: 

  * name: The name of the sender to include in the "FROM" address.
  * addr: The email address of the sender to include in the "FROM" address.



References: 

  * <https://datatracker.ietf.org/doc/html/rfc5322#section-3.6.2>



####  func (*Msg) [FromMailAddress](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L754) ¶ added in v0.7.0
    
    
    func (m *Msg) FromMailAddress(from *[mail](/net/mail).[Address](/net/mail#Address))

FromMailAddress sets the "FROM" address in the mail body for the Msg using a mail.Address instance. 

The "FROM" address is included in the mail body and indicates the sender of the message to the recipient. This address is visible in the email client and is typically displayed to the recipient. If the "FROM" address is not set, the msgWriter may attempt to use the envelope from address (if available) for sending. 

Parameters: 

  * from: The "FROM" address to set in the mail body as *mail.Address.



References: 

  * <https://datatracker.ietf.org/doc/html/rfc5322#section-3.6.2>



####  func (*Msg) [GetAddrHeader](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L1552) ¶ added in v0.3.5
    
    
    func (m *Msg) GetAddrHeader(header AddrHeader) []*[mail](/net/mail).[Address](/net/mail#Address)

GetAddrHeader returns the content of the requested address header for the Msg. 

This method retrieves the addresses associated with the specified address header. It returns a slice of pointers to mail.Address structures representing the addresses found in the header. If the requested header does not exist or contains no addresses, it will return nil. 

Parameters: 

  * header: The AddrHeader enum value indicating which address header to retrieve (e.g., "TO", "CC", "BCC", etc.).



Returns: 

  * A slice of pointers to mail.Address structures containing the addresses from the specified header.



References: 

  * <https://datatracker.ietf.org/doc/html/rfc5322#section-3.6>



####  func (*Msg) [GetAddrHeaderString](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L1572) ¶ added in v0.3.5
    
    
    func (m *Msg) GetAddrHeaderString(header AddrHeader) [][string](/builtin#string)

GetAddrHeaderString returns the address strings of the requested address header for the Msg. 

This method retrieves the addresses associated with the specified address header and returns them as a slice of strings. Each address is formatted as a string, which includes both the name (if available) and the email address. If the requested header does not exist or contains no addresses, it will return an empty slice. 

Parameters: 

  * header: The AddrHeader enum value indicating which address header to retrieve (e.g., "TO", "CC", "BCC", etc.).



Returns: 

  * A slice of strings containing the formatted addresses from the specified header.



References: 

  * <https://datatracker.ietf.org/doc/html/rfc5322#section-3.6>



####  func (*Msg) [GetAttachments](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L1724) ¶ added in v0.3.1
    
    
    func (m *Msg) GetAttachments() []*File

GetAttachments returns the attachments of the Msg. 

This method retrieves the list of files that have been attached to the email message. Each attachment includes details about the file, such as its name, content type, and data. 

Returns: 

  * A slice of File pointers representing the attachments of the email.



####  func (*Msg) [GetBcc](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L1674) ¶ added in v0.3.5
    
    
    func (m *Msg) GetBcc() []*[mail](/net/mail).[Address](/net/mail#Address)

GetBcc returns the content of the "Bcc" address header of the Msg. 

This method retrieves the list of email addresses set in the "Bcc" (blind carbon copy) header of the message. It returns a slice of pointers to `mail.Address` objects representing the Bcc recipient(s) of the email. 

Returns: 

  * A slice of `*mail.Address` containing the "Bcc" header addresses.



References: 

  * <https://datatracker.ietf.org/doc/html/rfc5322#section-3.6.3>



####  func (*Msg) [GetBccString](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L1688) ¶ added in v0.3.5
    
    
    func (m *Msg) GetBccString() [][string](/builtin#string)

GetBccString returns the content of the "Bcc" address header of the Msg as a string slice. 

This method retrieves the list of email addresses set in the "Bcc" (blind carbon copy) header of the message and returns them as a slice of strings, with each entry representing a formatted email address. 

Returns: 

  * A slice of strings containing the "Bcc" header addresses.



References: 

  * <https://datatracker.ietf.org/doc/html/rfc5322#section-3.6.3>



####  func (*Msg) [GetBoundary](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L1743) ¶ added in v0.4.3
    
    
    func (m *Msg) GetBoundary() [string](/builtin#string)

GetBoundary returns the boundary of the Msg. 

This method retrieves the MIME boundary that is used to separate different parts of the message, particularly in multipart emails. The boundary helps to differentiate between various sections such as plain text, HTML content, and attachments. 

NOTE: By default, random MIME boundaries are created. Using a predefined boundary will only work with messages that hold a single multipart part. Using a predefined boundary with several multipart parts will render the mail unreadable to the mail client. 

Returns: 

  * A string representing the boundary of the message.



References: 

  * <https://datatracker.ietf.org/doc/html/rfc2046#section-5.1.1>



####  func (*Msg) [GetCc](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L1646) ¶ added in v0.3.5
    
    
    func (m *Msg) GetCc() []*[mail](/net/mail).[Address](/net/mail#Address)

GetCc returns the content of the "Cc" address header of the Msg. 

This method retrieves the list of email addresses set in the "Cc" (carbon copy) header of the message. It returns a slice of pointers to `mail.Address` objects representing the secondary recipient(s) of the email. 

Returns: 

  * A slice of `*mail.Address` containing the "Cc" header addresses.



References: 

  * <https://datatracker.ietf.org/doc/html/rfc5322#section-3.6.3>



####  func (*Msg) [GetCcString](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L1660) ¶ added in v0.3.5
    
    
    func (m *Msg) GetCcString() [][string](/builtin#string)

GetCcString returns the content of the "Cc" address header of the Msg as a string slice. 

This method retrieves the list of email addresses set in the "Cc" (carbon copy) header of the message and returns them as a slice of strings, with each entry representing a formatted email address. 

Returns: 

  * A slice of strings containing the "Cc" header addresses.



References: 

  * <https://datatracker.ietf.org/doc/html/rfc5322#section-3.6.3>



####  func (*Msg) [GetEmbeds](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L1803) ¶ added in v0.3.9
    
    
    func (m *Msg) GetEmbeds() []*File

GetEmbeds returns the embedded files of the Msg. 

This method retrieves the list of files that have been embedded in the message. Embeds are typically images or other media files that are referenced directly in the content of the email, such as inline images in HTML emails. 

Returns: 

  * A slice of pointers to File structures representing the embedded files in the message.



References: 

  * <https://datatracker.ietf.org/doc/html/rfc2183>



####  func (*Msg) [GetFrom](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L1590) ¶ added in v0.3.5
    
    
    func (m *Msg) GetFrom() []*[mail](/net/mail).[Address](/net/mail#Address)

GetFrom returns the content of the "From" address header of the Msg. 

This method retrieves the list of email addresses set in the "From" header of the message. It returns a slice of pointers to `mail.Address` objects representing the sender(s) of the email. 

Returns: 

  * A slice of `*mail.Address` containing the "From" header addresses.



References: 

  * <https://datatracker.ietf.org/doc/html/rfc5322#section-3.6.2>



####  func (*Msg) [GetFromString](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L1604) ¶ added in v0.3.5
    
    
    func (m *Msg) GetFromString() [][string](/builtin#string)

GetFromString returns the content of the "From" address header of the Msg as a string slice. 

This method retrieves the list of email addresses set in the "From" header of the message and returns them as a slice of strings, with each entry representing a formatted email address. 

Returns: 

  * A slice of strings containing the "From" header addresses.



References: 

  * <https://datatracker.ietf.org/doc/html/rfc5322#section-3.6.2>



####  func (*Msg) [GetGenHeader](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L1702) ¶ added in v0.2.9
    
    
    func (m *Msg) GetGenHeader(header Header) [][string](/builtin#string)

GetGenHeader returns the content of the requested generic header of the Msg. 

This method retrieves the list of string values associated with the specified generic header of the message. It returns a slice of strings representing the header's values. 

Parameters: 

  * header: The Header field whose values are being retrieved.



Returns: 

  * A slice of strings containing the values of the specified generic header.



####  func (*Msg) [GetMessageID](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L1247) ¶ added in v0.5.0
    
    
    func (m *Msg) GetMessageID() [string](/builtin#string)

GetMessageID retrieves the "Message-ID" header from the Msg. 

This method checks if a "Message-ID" has been set in the message's generated headers. If a valid "Message-ID" exists in the Msg, it returns the first occurrence of the header. If the "Message-ID" has not been set or is empty, it returns an empty string. This allows other components to access the unique identifier for the message, which is useful for tracking and referencing in email systems. 

References: 

  * <https://datatracker.ietf.org/doc/html/rfc5322#section-3.6.4>



####  func (*Msg) [GetParts](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L1713) ¶ added in v0.3.0
    
    
    func (m *Msg) GetParts() []*Part

GetParts returns the message parts of the Msg. 

This method retrieves the list of parts that make up the email message. Each part may represent a different section of the email, such as a plain text body, HTML body, or attachments. 

Returns: 

  * A slice of Part pointers representing the message parts of the email.



####  func (*Msg) [GetRecipients](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L1519) ¶
    
    
    func (m *Msg) GetRecipients() ([][string](/builtin#string), [error](/builtin#error))

GetRecipients returns a list of the currently set "TO", "CC", and "BCC" addresses for the Msg. 

This method aggregates recipients from the "TO", "CC", and "BCC" headers and returns them as a slice of strings. If no recipients are found in these headers, it will return an error indicating that no recipient addresses are present. 

Returns: 

  * A slice of strings containing the recipients' addresses and an error if applicable.
  * If there are no recipient addresses set, it will return an error indicating no recipient addresses are available.



References: 

  * <https://datatracker.ietf.org/doc/html/rfc5322#section-3.6.3>



####  func (*Msg) [GetSender](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L1490) ¶
    
    
    func (m *Msg) GetSender(useFullAddr [bool](/builtin#bool)) ([string](/builtin#string), [error](/builtin#error))

GetSender returns the currently set envelope "FROM" address for the Msg. If no envelope "FROM" address is set, it will use the first "FROM" address from the mail body. If the useFullAddr parameter is true, it will return the full address string, including the name if it is set. 

If neither the envelope "FROM" nor the body "FROM" addresses are available, it will return an error indicating that no "FROM" address is present. 

Parameters: 

  * useFullAddr: A boolean indicating whether to return the full address string (including the name) or just the email address.



Returns: 

  * The sender's address as a string and an error if applicable.



References: 

  * <https://datatracker.ietf.org/doc/html/rfc5322#section-3.6.2>



####  func (*Msg) [GetTo](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L1618) ¶ added in v0.3.5
    
    
    func (m *Msg) GetTo() []*[mail](/net/mail).[Address](/net/mail#Address)

GetTo returns the content of the "To" address header of the Msg. 

This method retrieves the list of email addresses set in the "To" header of the message. It returns a slice of pointers to `mail.Address` objects representing the primary recipient(s) of the email. 

Returns: 

  * A slice of `*mail.Address` containing the "To" header addresses.



References: 

  * <https://datatracker.ietf.org/doc/html/rfc5322#section-3.6.3>



####  func (*Msg) [GetToString](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L1632) ¶ added in v0.3.5
    
    
    func (m *Msg) GetToString() [][string](/builtin#string)

GetToString returns the content of the "To" address header of the Msg as a string slice. 

This method retrieves the list of email addresses set in the "To" header of the message and returns them as a slice of strings, with each entry representing a formatted email address. 

Returns: 

  * A slice of strings containing the "To" header addresses.



References: 

  * <https://datatracker.ietf.org/doc/html/rfc5322#section-3.6.3>



####  func (*Msg) [HasSendError](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L2683) ¶ added in v0.3.7
    
    
    func (m *Msg) HasSendError() [bool](/builtin#bool)

HasSendError returns true if the Msg experienced an error during message delivery and the sendError field of the Msg is not nil. 

This method checks whether the message has encountered a delivery error by verifying if the sendError field is populated. 

Returns: 

  * A boolean value indicating whether a send error occurred (true if an error is present).



####  func (*Msg) [IsDelivered](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L1379) ¶ added in v0.4.1
    
    
    func (m *Msg) IsDelivered() [bool](/builtin#bool)

IsDelivered indicates whether the Msg has been delivered. 

This method checks the internal state of the message to determine if it has been successfully delivered. It returns true if the message is marked as delivered and false otherwise. This can be useful for tracking the status of the email communication. 

Returns: 

  * A boolean value indicating the delivery status of the message (true if delivered, false otherwise).



####  func (*Msg) [NewReader](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L2645) ¶ added in v0.3.2
    
    
    func (m *Msg) NewReader() *Reader

NewReader returns a Reader type that satisfies the io.Reader interface. 

This method creates a new Reader for the Msg, capturing the current state of the message. Any subsequent changes made to the Msg after creating the Reader will not be reflected in the Reader's buffer. To reflect these changes in the Reader, you must call Msg.UpdateReader to update the Reader's content with the current state of the Msg. 

Returns: 

  * A pointer to a Reader, which allows the Msg to be read as a stream of bytes.



IMPORTANT: Any changes made to the Msg after creating the Reader will not be reflected in the Reader unless Msg.UpdateReader is called. 

References: 

  * <https://datatracker.ietf.org/doc/html/rfc5322>



####  func (*Msg) [ReplyTo](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L1167) ¶
    
    
    func (m *Msg) ReplyTo(addr [string](/builtin#string)) [error](/builtin#error)

ReplyTo sets the "Reply-To" address for the Msg, specifying where replies should be sent. 

This method takes a single email address as input and attempts to parse it. If the address is valid, it sets the "Reply-To" header in the message. The "Reply-To" address can be different from the "From" address, allowing the sender to specify an alternate address for responses. If the provided address cannot be parsed, an error will be returned, indicating the parsing failure. 

Parameters: 

  * addr: The email address to set as the "Reply-To" address.



References: 

  * <https://datatracker.ietf.org/doc/html/rfc5322#section-3.6.2>



####  func (*Msg) [ReplyToFormat](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L1199) ¶
    
    
    func (m *Msg) ReplyToFormat(name, addr [string](/builtin#string)) [error](/builtin#error)

ReplyToFormat sets the "Reply-To" address for the Msg using the provided name and email address, specifying where replies should be sent. 

This method formats the name and email address into a single "Reply-To" header. If the formatted address is valid, it sets the "Reply-To" header in the message. This allows the sender to specify a display name along with the reply address, providing clarity for recipients. If the constructed address cannot be parsed, an error will be returned, indicating the parsing failure. 

Parameters: 

  * name: The display name associated with the reply address.
  * addr: The email address to set as the "Reply-To" address.



References: 

  * <https://datatracker.ietf.org/doc/html/rfc5322#section-3.6.2>



####  func (*Msg) [ReplyToMailAddress](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L1181) ¶ added in v0.7.0
    
    
    func (m *Msg) ReplyToMailAddress(addr *[mail](/net/mail).[Address](/net/mail#Address))

ReplyToMailAddress sets one or more "BCC" (blind carbon copy) addresses in the mail body for the Msg. 

The "Reply-To" address can be different from the "From" address, allowing the sender to specify an alternate address for responses. 

Parameters: 

  * addr: The mail.Address instance to set as the "Reply-To" address.



References: 

  * <https://datatracker.ietf.org/doc/html/rfc5322#section-3.6.3>



####  func (*Msg) [RequestMDNAddTo](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L1441) ¶ added in v0.2.7
    
    
    func (m *Msg) RequestMDNAddTo(rcpt [string](/builtin#string)) [error](/builtin#error)

RequestMDNAddTo adds an additional recipient to the "Disposition-Notification-To" header for the Msg. 

This method allows you to append a new recipient address to the existing list of recipients for the MDN. The provided address is validated according to [RFC 5322](https://rfc-editor.org/rfc/rfc5322.html) standards. If the address is invalid, an error will be returned indicating the parsing failure. If the "Disposition-Notification-To" header is already set, the new recipient will be added to the existing list. 

Parameters: 

  * rcpt: The recipient email address to add to the "Disposition-Notification-To" header.



References: 

  * <https://datatracker.ietf.org/doc/html/rfc8098>



####  func (*Msg) [RequestMDNAddToFormat](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L1469) ¶ added in v0.2.7
    
    
    func (m *Msg) RequestMDNAddToFormat(name, addr [string](/builtin#string)) [error](/builtin#error)

RequestMDNAddToFormat adds an additional formatted recipient to the "Disposition-Notification-To" header for the Msg. 

This method allows you to specify a recipient address along with a name, formatting it appropriately before adding it to the existing list of recipients for the MDN. The formatted address is validated according to [RFC 5322](https://rfc-editor.org/rfc/rfc5322.html) standards. If the provided address is invalid, an error will be returned. This method internally calls RequestMDNAddTo to handle the actual addition of the recipient. 

Parameters: 

  * name: The name of the recipient to add to the "Disposition-Notification-To" header.
  * addr: The email address of the recipient to add to the "Disposition-Notification-To" header.



References: 

  * <https://datatracker.ietf.org/doc/html/rfc8098>



####  func (*Msg) [RequestMDNTo](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L1396) ¶ added in v0.2.7
    
    
    func (m *Msg) RequestMDNTo(rcpts ...[string](/builtin#string)) [error](/builtin#error)

RequestMDNTo adds the "Disposition-Notification-To" header to the Msg to request a Message Disposition Notification (MDN) from the receiving end, as specified in [RFC 8098](https://rfc-editor.org/rfc/rfc8098.html). 

This method allows you to provide a list of recipient addresses to receive the MDN. Each address is validated according to [RFC 5322](https://rfc-editor.org/rfc/rfc5322.html) standards. If ANY address is invalid, an error will be returned indicating the parsing failure. If the "Disposition-Notification-To" header is already set, it will be updated with the new list of addresses. 

Parameters: 

  * rcpts: One or more recipient email addresses to request the MDN from.



References: 

  * <https://datatracker.ietf.org/doc/html/rfc8098>



####  func (*Msg) [RequestMDNToFormat](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L1425) ¶ added in v0.2.7
    
    
    func (m *Msg) RequestMDNToFormat(name, addr [string](/builtin#string)) [error](/builtin#error)

RequestMDNToFormat adds the "Disposition-Notification-To" header to the Msg to request a Message Disposition Notification (MDN) from the receiving end, as specified in [RFC 8098](https://rfc-editor.org/rfc/rfc8098.html). 

This method allows you to provide a recipient address along with a name, formatting it appropriately. Address validation is performed according to [RFC 5322](https://rfc-editor.org/rfc/rfc5322.html) standards. If the provided address is invalid, an error will be returned. This method internally calls RequestMDNTo to handle the actual setting of the header. 

Parameters: 

  * name: The name of the recipient for the MDN request.
  * addr: The email address of the recipient for the MDN request.



References: 

  * <https://datatracker.ietf.org/doc/html/rfc8098>



####  func (*Msg) [Reset](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L2393) ¶ added in v0.1.2
    
    
    func (m *Msg) Reset()

Reset resets all headers, body parts, attachments, and embeds of the Msg. 

This method clears all address headers, attachments, embeds, generic headers, and body parts of the message. However, it preserves the existing encoding, charset, boundary, and other message-level settings. Use this method to reset the message content while keeping certain configurations intact. 

References: 

  * <https://datatracker.ietf.org/doc/html/rfc5322>



####  func (*Msg) [SendError](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L2711) ¶ added in v0.3.7
    
    
    func (m *Msg) SendError() [error](/builtin#error)

SendError returns the sendError field of the Msg. 

This method retrieves the error that occurred during the message delivery process, if any. It returns the sendError field, which holds the error encountered during sending. 

Returns: 

  * The error encountered during message delivery, or nil if no error occurred.



####  func (*Msg) [SendErrorIsTemp](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L2696) ¶ added in v0.3.7
    
    
    func (m *Msg) SendErrorIsTemp() [bool](/builtin#bool)

SendErrorIsTemp returns true if the Msg experienced a delivery error, and the corresponding error was of a temporary nature, meaning it can be retried later. 

This method checks whether the encountered sendError is a temporary error that can be retried. It uses the errors.As function to determine if the error is of type SendError and checks if the error is marked as temporary. 

Returns: 

  * A boolean value indicating whether the send error is temporary (true if the error is temporary).



####  func (*Msg) [ServerResponse](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L1756) ¶ added in v0.7.0
    
    
    func (m *Msg) ServerResponse() [string](/builtin#string)

ServerResponse returns the server's response after queuing the mail. 

This function retrieves the value of m.serverResponse, which typically contains information such as the queue ID returned by the mail server once a message has been queued. Unfortunately different mail server software returns different server responses, therefore you have to parse the output yourself. 

Returns: 

  * The server response string, usually containing the queue ID or status.



####  func (*Msg) [SetAddrHeader](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L575) ¶
    
    
    func (m *Msg) SetAddrHeader(header AddrHeader, values ...[string](/builtin#string)) [error](/builtin#error)

SetAddrHeader sets the specified AddrHeader for the Msg to the given values. 

Addresses are parsed according to [RFC 5322](https://rfc-editor.org/rfc/rfc5322.html). If parsing any of the provided values fails, an error is returned. If you cannot guarantee that all provided values are valid, you can use SetAddrHeaderIgnoreInvalid instead, which will silently skip any parsing errors. 

This method allows you to set address-related headers for the message, ensuring that the provided addresses are properly formatted and parsed. Using this method helps maintain the integrity of the email addresses within the message. 

Parameters: 

  * header: The AddrHeader to set in the Msg (e.g., "From", "To", "Cc", "Bcc").
  * values: One or more string values representing the email addresses to associate with the specified header.



Returns: 

  * An error if parsing the address according to [RFC 5322](https://rfc-editor.org/rfc/rfc5322.html) fails



References: 

  * <https://datatracker.ietf.org/doc/html/rfc5322#section-3.4>



####  func (*Msg) [SetAddrHeaderFromMailAddress](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L614) ¶ added in v0.7.0
    
    
    func (m *Msg) SetAddrHeaderFromMailAddress(header AddrHeader, values ...*[mail](/net/mail).[Address](/net/mail#Address))

SetAddrHeaderFromMailAddress sets the specified AddrHeader for the Msg to the given mail.Address values. 

This method allows you to set address-related headers for the message, with mail.Address instances as input. Using this method helps maintain the integrity of the email addresses within the message. 

Since we expect the mail.Address instances to be already parsed according to [RFC 5322](https://rfc-editor.org/rfc/rfc5322.html), this method will not attempt to perform any sanity checks except of nil pointers and therefore no error will be returned. Nil pointers will be silently ignored. 

Parameters: 

  * header: The AddrHeader to set in the Msg (e.g., "From", "To", "Cc", "Bcc").
  * addresses: One or more mail.Address pointers representing the email addresses to associate with the specified header.



References: 

  * <https://datatracker.ietf.org/doc/html/rfc5322#section-3.4>



####  func (*Msg) [SetAddrHeaderIgnoreInvalid](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L652) ¶
    
    
    func (m *Msg) SetAddrHeaderIgnoreInvalid(header AddrHeader, values ...[string](/builtin#string))

SetAddrHeaderIgnoreInvalid sets the specified AddrHeader for the Msg to the given values. 

Addresses are parsed according to [RFC 5322](https://rfc-editor.org/rfc/rfc5322.html). If parsing of any of the provided values fails, the error is ignored and the address is omitted from the address list. 

This method allows for setting address headers while ignoring invalid addresses. It is useful in scenarios where you want to ensure that only valid addresses are included without halting execution due to parsing errors. 

Parameters: 

  * header: The AddrHeader field to set in the Msg.
  * values: One or more string values representing email addresses.



References: 

  * <https://datatracker.ietf.org/doc/html/rfc5322#section-3.4>



####  func (*Msg) [SetAttachements](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L1777) deprecated added in v0.3.1
    
    
    func (m *Msg) SetAttachements(files []*File)

SetAttachements sets the attachments of the message. 

Deprecated: use SetAttachments instead. 

####  func (*Msg) [SetAttachments](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L1770) ¶ added in v0.4.3
    
    
    func (m *Msg) SetAttachments(files []*File)

SetAttachments sets the attachments of the message. 

This method allows you to specify the attachments for the message by providing a slice of File pointers. Each file represents an attachment that will be included in the email. 

Parameters: 

  * files: A slice of pointers to File structures representing the attachments to set for the message.



References: 

  * <https://datatracker.ietf.org/doc/html/rfc2183>



####  func (*Msg) [SetBodyHTMLTemplate](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L1906) ¶ added in v0.2.2
    
    
    func (m *Msg) SetBodyHTMLTemplate(tpl *[ht](/html/template).[Template](/html/template#Template), data interface{}, opts ...PartOption) [error](/builtin#error)

SetBodyHTMLTemplate sets the body of the message from a given html/template.Template pointer. 

This method sets the body of the message using the provided HTML template and data. The content type will be set to "text/html" automatically. The method executes the template with the provided data and writes the output to the message body. If the template is nil or fails to execute, an error will be returned. 

Parameters: 

  * tpl: A pointer to the html/template.Template to be used for the message body.
  * data: The data to populate the template.
  * opts: Optional parameters for customizing the body part.



Returns: 

  * An error if the template is nil or fails to execute, otherwise nil.



References: 

  * <https://datatracker.ietf.org/doc/html/rfc2045>
  * <https://datatracker.ietf.org/doc/html/rfc2046>



####  func (*Msg) [SetBodyString](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L1858) ¶
    
    
    func (m *Msg) SetBodyString(contentType ContentType, content [string](/builtin#string), opts ...PartOption)

SetBodyString sets the body of the message. 

This method sets the body of the message using the provided content type and string content. The body can be set as plain text, HTML, or other formats based on the specified content type. Optional part settings can be passed through PartOption to customize the message body further. 

Parameters: 

  * contentType: The ContentType of the body (e.g., plain text, HTML).
  * content: The string content to set as the body of the message.
  * opts: Optional parameters for customizing the body part.



References: 

  * <https://datatracker.ietf.org/doc/html/rfc2045>
  * <https://datatracker.ietf.org/doc/html/rfc2046>

Example (DifferentTypes) ¶

This code example shows how to use Msg.SetBodyString to set a string as message body with different content types 
    
    
    m := mail.NewMsg()
    m.SetBodyString(mail.TypeTextPlain, "This is a mail body that with content type: text/plain")
    m.SetBodyString(mail.TypeTextHTML, "<p>This is a mail body that with content type: text/html</p>")
    

Example (WithPartOption) ¶

This code example shows how to use Msg.SetBodyString to set a string as message body a PartOption to override the default encoding 
    
    
    m := mail.NewMsg(mail.WithEncoding(mail.EncodingB64))
    m.SetBodyString(mail.TypeTextPlain, "This is a mail body that with content type: text/plain",
    	mail.WithPartEncoding(mail.EncodingQP))
    

####  func (*Msg) [SetBodyTextTemplate](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L1937) ¶ added in v0.2.2
    
    
    func (m *Msg) SetBodyTextTemplate(tpl *[tt](/text/template).[Template](/text/template#Template), data interface{}, opts ...PartOption) [error](/builtin#error)

SetBodyTextTemplate sets the body of the message from a given text/template.Template pointer. 

This method sets the body of the message using the provided text template and data. The content type will be set to "text/plain" automatically. The method executes the template with the provided data and writes the output to the message body. If the template is nil or fails to execute, an error will be returned. 

Parameters: 

  * tpl: A pointer to the text/template.Template to be used for the message body.
  * data: The data to populate the template.
  * opts: Optional parameters for customizing the body part.



Returns: 

  * An error if the template is nil or fails to execute, otherwise nil.



References: 

  * <https://datatracker.ietf.org/doc/html/rfc2045>
  * <https://datatracker.ietf.org/doc/html/rfc2046>

Example ¶

This code example shows how to use a text/template as message Body. Msg.SetBodyHTMLTemplate works anolog to this just with html/template instead 
    
    
    type MyStruct struct {
    	Placeholder string
    }
    data := MyStruct{Placeholder: "Teststring"}
    tpl, err := template.New("test").Parse("This is a {{.Placeholder}}")
    if err != nil {
    	panic(err)
    }
    
    m := mail.NewMsg()
    if err := m.SetBodyTextTemplate(tpl, data); err != nil {
    	panic(err)
    }
    

####  func (*Msg) [SetBodyWriter](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L1879) ¶
    
    
    func (m *Msg) SetBodyWriter(
    	contentType ContentType, writeFunc func([io](/io).[Writer](/io#Writer)) ([int64](/builtin#int64), [error](/builtin#error)),
    	opts ...PartOption,
    )

SetBodyWriter sets the body of the message. 

This method sets the body of the message using a write function, allowing content to be written directly to the body. The content type determines the format (e.g., plain text, HTML). Optional part settings can be provided via PartOption to customize the body further. 

Parameters: 

  * contentType: The ContentType of the body (e.g., plain text, HTML).
  * writeFunc: A function that writes content to an io.Writer and returns the number of bytes written and an error, if any.
  * opts: Optional parameters for customizing the body part.



References: 

  * <https://datatracker.ietf.org/doc/html/rfc2045>
  * <https://datatracker.ietf.org/doc/html/rfc2046>



####  func (*Msg) [SetBoundary](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L401) ¶
    
    
    func (m *Msg) SetBoundary(boundary [string](/builtin#string))

SetBoundary sets or overrides the currently set boundary of the Msg. 

This method allows you to specify a custom boundary string for the MIME message. The boundary is used to separate different parts of the message, especially when dealing with multipart messages. 

NOTE: By default, random MIME boundaries are created. This option should only be used if a specific boundary is required for the email message. Using a predefined boundary will only work with messages that hold a single multipart part. Using a predefined boundary with several multipart parts will render the mail unreadable to the mail client. 

Parameters: 

  * boundary: The string value representing the boundary to set for the Msg, used in multipart messages to delimit different sections.



####  func (*Msg) [SetBulk](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L1282) ¶
    
    
    func (m *Msg) SetBulk()

SetBulk sets the "Precedence: bulk" and "X-Auto-Response-Suppress: All" headers for the Msg, which are recommended for automated emails such as out-of-office replies. 

The "Precedence: bulk" header indicates that the message is a bulk email, and the "X-Auto-Response-Suppress: All" header instructs mail servers and clients to suppress automatic responses to this message. This is particularly useful for reducing unnecessary replies to automated notifications or replies. 

References: 

  * <https://www.rfc-editor.org/rfc/rfc2076#section-3.9>
  * <https://learn.microsoft.com/en-us/openspecs/exchange_server_protocols/ms-oxcmail/ced68690-498a-4567-9d14-5c01f974d8b1#Appendix_A_Target_51>



####  func (*Msg) [SetCharset](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L368) ¶
    
    
    func (m *Msg) SetCharset(charset Charset)

SetCharset sets or overrides the currently set encoding charset of the Msg. 

This method allows you to specify a character set for the email message. The charset is important for ensuring that the content of the message is correctly interpreted by mail clients. Common charset values include UTF-8, ISO-8859-1, and others. If a charset is not explicitly set, CharsetUTF8 is used as default. 

Parameters: 

  * charset: The Charset value to set for the Msg, determining the encoding used for the message content.



####  func (*Msg) [SetDate](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L1296) ¶
    
    
    func (m *Msg) SetDate()

SetDate sets the "Date" header for the Msg to the current time in a valid [RFC 1123](https://rfc-editor.org/rfc/rfc1123.html) format. 

This method retrieves the current time and formats it according to [RFC 1123](https://rfc-editor.org/rfc/rfc1123.html), ensuring that the "Date" header is compliant with email standards. The "Date" header indicates when the message was created, providing recipients with context for the timing of the email. 

References: 

  * <https://datatracker.ietf.org/doc/html/rfc5322#section-3.3>
  * <https://datatracker.ietf.org/doc/html/rfc1123>



####  func (*Msg) [SetDateWithValue](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L1313) ¶ added in v0.1.2
    
    
    func (m *Msg) SetDateWithValue(timeVal [time](/time).[Time](/time#Time))

SetDateWithValue sets the "Date" header for the Msg using the provided time value in a valid [RFC 1123](https://rfc-editor.org/rfc/rfc1123.html) format. 

This method takes a `time.Time` value as input and formats it according to [RFC 1123](https://rfc-editor.org/rfc/rfc1123.html), ensuring that the "Date" header is compliant with email standards. The "Date" header indicates when the message was created, providing recipients with context for the timing of the email. This allows for setting a custom date rather than using the current time. 

Parameters: 

  * timeVal: The time value used to set the "Date" header.



References: 

  * <https://datatracker.ietf.org/doc/html/rfc5322#section-3.3>
  * <https://datatracker.ietf.org/doc/html/rfc1123>



####  func (*Msg) [SetEmbeds](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L1817) ¶ added in v0.3.9
    
    
    func (m *Msg) SetEmbeds(files []*File)

SetEmbeds sets the embedded files of the message. 

This method allows you to specify the files to be embedded in the message by providing a slice of File pointers. Embedded files, such as images or media, are typically used for inline content in HTML emails. 

Parameters: 

  * files: A slice of pointers to File structures representing the embedded files to set for the message.



References: 

  * <https://datatracker.ietf.org/doc/html/rfc2183>



####  func (*Msg) [SetEncoding](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L382) ¶
    
    
    func (m *Msg) SetEncoding(encoding Encoding)

SetEncoding sets or overrides the currently set Encoding of the Msg. 

This method allows you to specify the encoding type for the email message. The encoding determines how the message content is represented and can affect the size and compatibility of the email. Common encoding types include Base64 and Quoted-Printable. Setting a new encoding may also adjust how the message content is processed and transmitted. 

Parameters: 

  * encoding: The Encoding value to set for the Msg, determining the method used to encode the message content.



####  func (*Msg) [SetGenHeader](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L506) ¶ added in v0.3.5
    
    
    func (m *Msg) SetGenHeader(header Header, values ...[string](/builtin#string))

SetGenHeader sets a generic header field of the Msg to the provided list of values. 

This method is intended for setting generic headers in the email message. It takes a header name and a variadic list of string values, encoding them as necessary before storing them in the message's internal header map. 

Note: For adding email address-related headers (like "To:", "From", "Cc", etc.), use SetAddrHeader instead to ensure proper formatting and validation. 

Parameters: 

  * header: The header field to set in the Msg.
  * values: One or more string values to associate with the header field.



This method ensures that all values are appropriately encoded for email transmission, adhering to the necessary standards. 

References: 

  * <https://datatracker.ietf.org/doc/html/rfc5322#section-3>
  * <https://datatracker.ietf.org/doc/html/rfc2047>



####  func (*Msg) [SetGenHeaderPreformatted](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L548) ¶ added in v0.3.5
    
    
    func (m *Msg) SetGenHeaderPreformatted(header Header, value [string](/builtin#string))

SetGenHeaderPreformatted sets a generic header field of the Msg which content is already preformatted. 

This method does not take a slice of values but only a single value. The reason for this is that we do not perform any content alteration on these kinds of headers and expect the user to have already taken care of any kind of formatting required for the header. 

Note: This method should be used only as a last resort. Since the user is responsible for the formatting of the message header, we cannot guarantee any compliance with [RFC 2822](https://rfc-editor.org/rfc/rfc2822.html). It is advised to use SetGenHeader instead for general header fields. 

Parameters: 

  * header: The header field to set in the Msg.
  * value: The preformatted string value to associate with the header field.



References: 

  * <https://datatracker.ietf.org/doc/html/rfc2822>



####  func (*Msg) [SetHeader](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L483) deprecated
    
    
    func (m *Msg) SetHeader(header Header, values ...[string](/builtin#string))

SetHeader sets a generic header field of the Msg. 

Deprecated: This method only exists for compatibility reasons. Please use SetGenHeader instead. For adding address headers like "To:" or "From", use SetAddrHeader instead. 

This method allows you to set a header field for the message, providing the header name and its corresponding values. However, it is recommended to utilize the newer methods for better clarity and functionality. Using SetGenHeader or SetAddrHeader is preferred for more specific header types, ensuring proper handling of the message headers. 

Parameters: 

  * header: The header field to set in the Msg.
  * values: One or more string values to associate with the header field.



References: 

  * <https://datatracker.ietf.org/doc/html/rfc5322#section-3>
  * <https://datatracker.ietf.org/doc/html/rfc2047>



####  func (*Msg) [SetHeaderPreformatted](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L528) deprecated added in v0.3.4
    
    
    func (m *Msg) SetHeaderPreformatted(header Header, value [string](/builtin#string))

SetHeaderPreformatted sets a generic header field of the Msg, which content is already preformatted. 

Deprecated: This method only exists for compatibility reasons. Please use SetGenHeaderPreformatted instead for setting preformatted generic header fields. 

Parameters: 

  * header: The header field to set in the Msg.
  * value: The preformatted string value to associate with the header field.



References: 

  * <https://datatracker.ietf.org/doc/html/rfc5322#section-3>
  * <https://datatracker.ietf.org/doc/html/rfc2047>



####  func (*Msg) [SetImportance](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L1329) ¶
    
    
    func (m *Msg) SetImportance(importance Importance)

SetImportance sets the "Importance" and "Priority" headers for the Msg to the specified Importance level. 

This method adjusts the email's importance based on the provided Importance value. If the importance level is set to `ImportanceNormal`, no headers are modified. Otherwise, it sets the "Importance", "Priority", "X-Priority", and "X-MSMail-Priority" headers accordingly, providing email clients with information on how to prioritize the message. This allows the sender to indicate the significance of the email to recipients. 

Parameters: 

  * importance: The Importance value that determines the priority of the email message.



References: 

  * <https://datatracker.ietf.org/doc/html/rfc2156>



####  func (*Msg) [SetMIMEVersion](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L421) ¶
    
    
    func (m *Msg) SetMIMEVersion(version MIMEVersion)

SetMIMEVersion sets or overrides the currently set MIME version of the Msg. 

In the context of email, MIME Version 1.0 is the only officially standardized and supported version. Although MIME has been updated and extended over time through various RFCs, these updates do not introduce new MIME versions; they refine or add features within the framework of MIME 1.0. Therefore, there is generally no need to use this function to set a different MIME version. 

Parameters: 

  * version: The MIMEVersion value to set for the Msg, which determines the MIME version used in the email message.



References: 

  * <https://datatracker.ietf.org/doc/html/rfc1521>
  * <https://datatracker.ietf.org/doc/html/rfc2045>
  * <https://datatracker.ietf.org/doc/html/rfc2049>



####  func (*Msg) [SetMessageID](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L1228) ¶
    
    
    func (m *Msg) SetMessageID()

SetMessageID generates and sets a unique "Message-ID" header for the Msg. 

This method creates a "Message-ID" string using a randomly generated string and the hostname of the machine. The generated ID helps uniquely identify the message in email systems, facilitating tracking and preventing duplication. If the hostname cannot be retrieved, it defaults to "localhost.localdomain". 

The generated Message-ID follows the format "<randomString@hostname>". 

References: 

  * <https://datatracker.ietf.org/doc/html/rfc5322#section-3.6.4>



####  func (*Msg) [SetMessageIDWithValue](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L1268) ¶
    
    
    func (m *Msg) SetMessageIDWithValue(messageID [string](/builtin#string))

SetMessageIDWithValue sets the "Message-ID" header for the Msg using the provided messageID string. 

This method formats the input messageID by enclosing it in angle brackets ("<>") and sets it as the "Message-ID" header in the message. The "Message-ID" is a unique identifier for the email, helping email clients and servers to track and reference the message. There are no validations performed on the input messageID, so it should be in a suitable format for use as a Message-ID. 

Parameters: 

  * messageID: The string to set as the "Message-ID" in the message header.



References: 

  * <https://datatracker.ietf.org/doc/html/rfc5322#section-3.6.4>



####  func (*Msg) [SetOrganization](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L1350) ¶ added in v0.1.3
    
    
    func (m *Msg) SetOrganization(org [string](/builtin#string))

SetOrganization sets the "Organization" header for the Msg to the specified organization string. 

This method allows you to specify the organization associated with the email sender. The "Organization" header provides recipients with information about the organization that is sending the message. This can help establish context and credibility for the email communication. 

Parameters: 

  * org: The name of the organization to be set in the "Organization" header.



References: 

  * <https://datatracker.ietf.org/doc/html/rfc5322#section-3.4>



####  func (*Msg) [SetPGPType](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L436) ¶ added in v0.3.9
    
    
    func (m *Msg) SetPGPType(pgptype PGPType)

SetPGPType sets or overrides the currently set PGP type for the Msg, determining the encryption or signature method. 

This method allows you to specify the PGP type that will be used when encrypting or signing the message. Different PGP types correspond to various encryption and signing algorithms, and selecting the appropriate type is essential for ensuring the security and integrity of the message content. 

Parameters: 

  * pgptype: The PGPType value to set for the Msg, which determines the encryption or signature method used for the email message.



####  func (*Msg) [SetUserAgent](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L1366) ¶ added in v0.1.3
    
    
    func (m *Msg) SetUserAgent(userAgent [string](/builtin#string))

SetUserAgent sets the "User-Agent" and "X-Mailer" headers for the Msg to the specified user agent string. 

This method allows you to specify the user agent or mailer software used to send the email. The "User-Agent" and "X-Mailer" headers provide recipients with information about the email client or application that generated the message. This can be useful for identifying the source of the email, particularly for troubleshooting or filtering purposes. 

Parameters: 

  * userAgent: The user agent or mailer software to be set in the "User-Agent" and "X-Mailer" headers.



References: 

  * <https://datatracker.ietf.org/doc/html/rfc5322#section-3.6.7>



####  func (*Msg) [SignWithKeypair](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L2751) ¶ added in v0.6.0
    
    
    func (m *Msg) SignWithKeypair(privateKey [crypto](/crypto).[PrivateKey](/crypto#PrivateKey), certificate *[x509](/crypto/x509).[Certificate](/crypto/x509#Certificate),
    	intermediateCert *[x509](/crypto/x509).[Certificate](/crypto/x509#Certificate),
    ) [error](/builtin#error)

SignWithKeypair configures the Msg to be signed with S/MIME using RSA or ECDSA public-/private keypair. 

This function sets up S/MIME signing for the Msg by associating it with the provided private key, certificate, and intermediate certificate. 

Parameters: 

  * privateKey: The RSA private key used for signing.
  * certificate: The x509 certificate associated with the private key.
  * intermediateCert: An optional intermediate x509 certificate for chain validation.



Returns: 

  * An error if any issue occurred while configuring S/MIME signing; otherwise nil.



####  func (*Msg) [SignWithTLSCertificate](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L2772) ¶ added in v0.6.0
    
    
    func (m *Msg) SignWithTLSCertificate(keyPairTLS *[tls](/crypto/tls).[Certificate](/crypto/tls#Certificate)) [error](/builtin#error)

SignWithTLSCertificate signs the Msg with the provided *tls.Certificate. 

This function configures the Msg for S/MIME signing using the private key and certificates from the provided TLS certificate. It supports both RSA and ECDSA private keys. 

Parameters: 

  * keyPairTlS: The *tls.Certificate containing the private key and associated certificate chain.



Returns: 

  * An error if any issue occurred during parsing, signing configuration, or unsupported private key type.



####  func (*Msg) [Subject](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L1213) ¶
    
    
    func (m *Msg) Subject(subj [string](/builtin#string))

Subject sets the "Subject" header for the Msg, specifying the topic of the message. 

This method takes a single string as input and sets it as the "Subject" of the email. The subject line provides a brief summary of the content of the message, allowing recipients to quickly understand its purpose. 

Parameters: 

  * subj: The subject line of the email.



References: 

  * <https://datatracker.ietf.org/doc/html/rfc5322#section-3.6.5>



####  func (*Msg) [To](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L788) ¶
    
    
    func (m *Msg) To(rcpts ...[string](/builtin#string)) [error](/builtin#error)

To sets one or more "TO" addresses in the mail body for the Msg. 

The "TO" address specifies the primary recipient(s) of the message and is included in the mail body. This address is visible to the recipient and any other recipients of the message. Multiple "TO" addresses can be set by passing them as variadic arguments to this method. Each provided address is validated according to [RFC 5322](https://rfc-editor.org/rfc/rfc5322.html), and an error will be returned if ANY validation fails. 

Parameters: 

  * rcpts: One or more recipient email addresses to include in the "TO" field.



References: 

  * <https://datatracker.ietf.org/doc/html/rfc5322#section-3.6.3>



####  func (*Msg) [ToFromString](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L888) ¶ added in v0.4.1
    
    
    func (m *Msg) ToFromString(rcpts [string](/builtin#string)) [error](/builtin#error)

ToFromString takes a string of comma-separated email addresses, validates each, and sets them as the "TO" addresses for the Msg. 

This method allows you to pass a single string containing multiple email addresses separated by commas. Each address is validated according to [RFC 5322](https://rfc-editor.org/rfc/rfc5322.html) and set as a recipient in the "TO" field. If any validation fails, an error will be returned. The addresses are visible in the mail body and displayed to recipients in the mail client. Any "TO" address applied previously will be overwritten. 

Parameters: 

  * rcpts: A string containing multiple recipient addresses separated by commas.



References: 

  * <https://datatracker.ietf.org/doc/html/rfc5322#section-3.6.3>



####  func (*Msg) [ToIgnoreInvalid](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L871) ¶
    
    
    func (m *Msg) ToIgnoreInvalid(rcpts ...[string](/builtin#string))

ToIgnoreInvalid sets one or more "TO" addresses in the mail body for the Msg, ignoring any invalid addresses. 

This method allows you to add multiple "TO" recipients to the message body. Unlike the standard `To` method, any invalid addresses are ignored, and no error is returned for those addresses. Valid addresses will still be included in the "TO" field, which is visible in the recipient's mail client. Use this method with caution if address validation is critical. Invalid addresses are determined according to [RFC 5322](https://rfc-editor.org/rfc/rfc5322.html). 

Parameters: 

  * rcpts: One or more recipient addresses to add to the "TO" field.



References: 

  * <https://datatracker.ietf.org/doc/html/rfc5322#section-3.6.3>



####  func (*Msg) [ToMailAddress](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L804) ¶ added in v0.7.0
    
    
    func (m *Msg) ToMailAddress(rcpts ...*[mail](/net/mail).[Address](/net/mail#Address))

ToMailAddress sets one or more "TO" addresses in the mail body for the Msg. 

The "TO" address specifies the primary recipient(s) of the message and is included in the mail body. This address is visible to the recipient and any other recipients of the message. Multiple "TO" addresses can be set by passing them as variadic arguments to this method. 

Parameters: 

  * rcpts: One or more recipient email addresses as mail.Address instance to include in the "TO" field.



References: 

  * <https://datatracker.ietf.org/doc/html/rfc5322#section-3.6.3>



####  func (*Msg) [UnsetAllAttachments](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L1788) ¶ added in v0.4.1
    
    
    func (m *Msg) UnsetAllAttachments()

UnsetAllAttachments unsets the attachments of the message. 

This method removes all attachments from the message by setting the attachments to nil, effectively clearing any previously set attachments. 

References: 

  * <https://datatracker.ietf.org/doc/html/rfc2183>



####  func (*Msg) [UnsetAllEmbeds](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L1828) ¶ added in v0.4.1
    
    
    func (m *Msg) UnsetAllEmbeds()

UnsetAllEmbeds unsets the embedded files of the message. 

This method removes all embedded files from the message by setting the embeds to nil, effectively clearing any previously set embedded files. 

References: 

  * <https://datatracker.ietf.org/doc/html/rfc2183>



####  func (*Msg) [UnsetAllParts](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L1839) ¶ added in v0.4.1
    
    
    func (m *Msg) UnsetAllParts()

UnsetAllParts unsets the embeds and attachments of the message. 

This method removes all embedded files and attachments from the message by unsetting both the embeds and attachments, effectively clearing all previously set message parts. 

References: 

  * <https://datatracker.ietf.org/doc/html/rfc2183>



####  func (*Msg) [UpdateReader](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L2667) ¶ added in v0.3.2
    
    
    func (m *Msg) UpdateReader(reader *Reader)

UpdateReader updates a Reader with the current content of the Msg and resets the Reader's position to the start. 

This method rewrites the content of the provided Reader to reflect any changes made to the Msg. It resets the Reader's position to the beginning and updates the buffer with the latest message content. 

Parameters: 

  * reader: A pointer to the Reader that will be updated with the Msg's current content.



References: 

  * <https://datatracker.ietf.org/doc/html/rfc5322>



####  func (*Msg) [Write](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L2497) ¶
    
    
    func (m *Msg) Write(writer [io](/io).[Writer](/io#Writer)) ([int64](/builtin#int64), [error](/builtin#error))

Write is an alias method to WriteTo for compatibility reasons. 

This method provides a backward-compatible way to write the formatted Msg to the provided io.Writer by calling the WriteTo method. It writes the email message, including headers, body, and attachments, to the io.Writer and returns the number of bytes written and any error encountered. 

Parameters: 

  * writer: The io.Writer to which the formatted message will be written.



Returns: 

  * The total number of bytes written.
  * An error if any occurred during the writing process, otherwise nil.



References: 

  * <https://datatracker.ietf.org/doc/html/rfc5322>



####  func (*Msg) [WriteTo](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L2434) ¶ added in v0.1.9
    
    
    func (m *Msg) WriteTo(writer [io](/io).[Writer](/io#Writer)) ([int64](/builtin#int64), [error](/builtin#error))

WriteTo writes the formatted Msg into the given io.Writer and satisfies the io.WriterTo interface. 

This method writes the email message, including its headers, body, and attachments, to the provided io.Writer. It applies any middlewares to the message before writing it. The total number of bytes written and any error encountered during the writing process are returned. 

Parameters: 

  * writer: The io.Writer to which the formatted message will be written.



Returns: 

  * The total number of bytes written.
  * An error if any occurred during the writing process, otherwise nil.



References: 

  * <https://datatracker.ietf.org/doc/html/rfc5322>



####  func (*Msg) [WriteToFile](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L2515) ¶ added in v0.2.3
    
    
    func (m *Msg) WriteToFile(name [string](/builtin#string)) [error](/builtin#error)

WriteToFile stores the Msg as a file on disk. It will try to create the given filename, and if the file already exists, it will be overwritten. 

This method writes the email message, including its headers, body, and attachments, to a file on disk. If the file cannot be created or an error occurs during writing, an error is returned. 

Parameters: 

  * name: The name of the file to be created or overwritten.



Returns: 

  * An error if the file cannot be created or if writing to the file fails, otherwise nil.



References: 

  * <https://datatracker.ietf.org/doc/html/rfc5322>



####  func (*Msg) [WriteToSendmail](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L2538) ¶ added in v0.1.2
    
    
    func (m *Msg) WriteToSendmail() [error](/builtin#error)

WriteToSendmail returns WriteToSendmailWithCommand with a default sendmail path. 

This method sends the email message using the default sendmail path. It calls WriteToSendmailWithCommand using the standard SendmailPath. If sending via sendmail fails, an error is returned. 

Returns: 

  * An error if sending the message via sendmail fails, otherwise nil.



References: 

  * <https://datatracker.ietf.org/doc/html/rfc5321>

Example ¶

This code example shows how to utilize the Msg.WriteToSendmail method to send generated mails using a local sendmail installation 
    
    
    m := mail.NewMsg()
    m.SetBodyString(mail.TypeTextPlain, "This is the mail body string")
    if err := m.FromFormat("Toni Tester", "toni.tester@example.com"); err != nil {
    	panic(err)
    }
    if err := m.To("gandalf.tester@example.com"); err != nil {
    	panic(err)
    }
    if err := m.WriteToSendmail(); err != nil {
    	panic(err)
    }
    

####  func (*Msg) [WriteToSendmailWithCommand](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L2556) ¶ added in v0.1.5
    
    
    func (m *Msg) WriteToSendmailWithCommand(sendmailPath [string](/builtin#string)) [error](/builtin#error)

WriteToSendmailWithCommand returns WriteToSendmailWithContext with a default timeout of 5 seconds and a given sendmail path. 

This method sends the email message using the provided sendmail path, with a default timeout of 5 seconds. It creates a context with the specified timeout and then calls WriteToSendmailWithContext to send the message. 

Parameters: 

  * sendmailPath: The path to the sendmail executable to be used for sending the message.



Returns: 

  * An error if sending the message via sendmail fails, otherwise nil.



References: 

  * <https://datatracker.ietf.org/doc/html/rfc5321>



####  func (*Msg) [WriteToSendmailWithContext](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L2580) ¶ added in v0.1.2
    
    
    func (m *Msg) WriteToSendmailWithContext(ctx [context](/context).[Context](/context#Context), sendmailPath [string](/builtin#string), args ...[string](/builtin#string)) [error](/builtin#error)

WriteToSendmailWithContext opens a pipe to the local sendmail binary and tries to send the email through it. It takes a context.Context, the path to the sendmail binary, and additional arguments for the sendmail binary as parameters. 

This method establishes a pipe to the sendmail executable using the provided context and arguments. It writes the email message to the sendmail process via STDIN. If any errors occur during the communication with the sendmail binary, they will be captured and returned. 

Parameters: 

  * ctx: The context to control the timeout and cancellation of the sendmail process.
  * sendmailPath: The path to the sendmail executable.
  * args: Additional arguments for the sendmail binary.



Returns: 

  * An error if sending the message via sendmail fails, otherwise nil.



References: 

  * <https://datatracker.ietf.org/doc/html/rfc5321>

Example ¶

This code example shows how to send generated mails using a custom context and sendmail-compatbile command using the Msg.WriteToSendmailWithContext method 
    
    
    sendmailPath := "/opt/sendmail/sbin/sendmail"
    ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
    defer cancel()
    
    m := mail.NewMsg()
    m.SetBodyString(mail.TypeTextPlain, "This is the mail body string")
    if err := m.FromFormat("Toni Tester", "toni.tester@example.com"); err != nil {
    	panic(err)
    }
    if err := m.To("gandalf.tester@example.com"); err != nil {
    	panic(err)
    }
    if err := m.WriteToSendmailWithContext(ctx, sendmailPath); err != nil {
    	panic(err)
    }
    

####  func (*Msg) [WriteToSkipMiddleware](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L2466) ¶ added in v0.3.3
    
    
    func (m *Msg) WriteToSkipMiddleware(writer [io](/io).[Writer](/io#Writer), middleWareType MiddlewareType) ([int64](/builtin#int64), [error](/builtin#error))

WriteToSkipMiddleware writes the formatted Msg into the given io.Writer, but skips the specified middleware type. 

This method writes the email message to the provided io.Writer after applying all middlewares, except for the specified middleware type, which will be skipped. It temporarily removes the middleware of the given type, writes the message, and then restores the original middleware list. 

Parameters: 

  * writer: The io.Writer to which the formatted message will be written.
  * middleWareType: The MiddlewareType that should be skipped during the writing process.



Returns: 

  * The total number of bytes written.
  * An error if any occurred during the writing process, otherwise nil.



References: 

  * <https://datatracker.ietf.org/doc/html/rfc5322>



####  func (*Msg) [WriteToTempFile](https://github.com/wneessen/go-mail/blob/v0.7.2/msg_totmpfile.go#L20) ¶ added in v0.2.3
    
    
    func (m *Msg) WriteToTempFile() ([string](/builtin#string), [error](/builtin#error))

WriteToTempFile creates a temporary file and writes the Msg content to this file. 

This method generates a temporary file with a ".eml" extension, writes the Msg to it, and returns the filename of the created temporary file. 

Returns: 

  * A string representing the filename of the temporary file.
  * An error if the file creation or writing process fails.



####  type [MsgOption](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L174) ¶
    
    
    type MsgOption func(*Msg)

MsgOption is a function type that modifies a Msg instance during its creation or initialization. 

####  func [WithBoundary](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L295) ¶ added in v0.1.4
    
    
    func WithBoundary(boundary [string](/builtin#string)) MsgOption

WithBoundary sets the boundary of a Msg to the provided string value during its creation or initialization. 

NOTE: By default, random MIME boundaries are created. This option should only be used if a specific boundary is required for the email message. Using a predefined boundary will only work with messages that hold a single multipart part. Using a predefined boundary with several multipart parts will render the mail unreadable to the mail client. 

Parameters: 

  * boundary: The string value that specifies the desired boundary for the Msg.



Returns: 

  * A MsgOption function that can be used to customize the Msg instance.



####  func [WithCharset](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L232) ¶
    
    
    func WithCharset(charset Charset) MsgOption

WithCharset sets the Charset type for a Msg during its creation or initialization. 

This MsgOption function allows you to specify the character set to be used in the email message. The charset defines how the text in the message is encoded and interpreted by the email client. This option should be called when creating a new Msg instance to ensure that the desired charset is set correctly. 

Parameters: 

  * charset: The Charset value that specifies the desired character set for the Msg.



Returns: 

  * A MsgOption function that can be used to customize the Msg instance.



References: 

  * <https://datatracker.ietf.org/doc/html/rfc2047#section-5>



####  func [WithEncoding](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L253) ¶
    
    
    func WithEncoding(encoding Encoding) MsgOption

WithEncoding sets the Encoding type for a Msg during its creation or initialization. 

This MsgOption function allows you to specify the encoding type to be used in the email message. The encoding defines how the message content is encoded, which affects how it is transmitted and decoded by email clients. This option should be called when creating a new Msg instance to ensure that the desired encoding is set correctly. 

Parameters: 

  * encoding: The Encoding value that specifies the desired encoding type for the Msg.



Returns: 

  * A MsgOption function that can be used to customize the Msg instance.



References: 

  * <https://datatracker.ietf.org/doc/html/rfc2047#section-6>



####  func [WithMIMEVersion](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L276) ¶
    
    
    func WithMIMEVersion(version MIMEVersion) MsgOption

WithMIMEVersion sets the MIMEVersion type for a Msg during its creation or initialization. 

Note that in the context of email, MIME Version 1.0 is the only officially standardized and supported version. While MIME has been updated and extended over time via various RFCs, these updates and extensions do not introduce new MIME versions; they refine or add features within the framework of MIME 1.0. Therefore, there should be no reason to ever use this MsgOption. 

Parameters: 

  * version: The MIMEVersion value that specifies the desired MIME version for the Msg.



Returns: 

  * A MsgOption function that can be used to customize the Msg instance.



References: 

  * <https://datatracker.ietf.org/doc/html/rfc1521>
  * <https://datatracker.ietf.org/doc/html/rfc2045>
  * <https://datatracker.ietf.org/doc/html/rfc2049>



####  func [WithMiddleware](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L314) ¶ added in v0.2.8
    
    
    func WithMiddleware(middleware Middleware) MsgOption

WithMiddleware adds the given Middleware to the end of the list of the Client middlewares slice. Middleware are processed in FIFO order. 

This MsgOption function allows you to specify custom middleware that will be applied during the message handling process. Middleware can be used to modify the message, perform logging, or implement additional functionality as the message flows through the system. Each middleware is executed in the order it was added. 

Parameters: 

  * middleware: The Middleware to be added to the list for processing.



Returns: 

  * A MsgOption function that can be used to customize the Msg instance.



####  func [WithNoDefaultUserAgent](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L353) ¶ added in v0.4.2
    
    
    func WithNoDefaultUserAgent() MsgOption

WithNoDefaultUserAgent disables the inclusion of a default User-Agent header in the Msg during its creation or initialization. 

This MsgOption function allows you to customize the Msg instance by omitting the default User-Agent header, which is typically included to provide information about the software sending the email. This option can be useful when you want to have more control over the headers included in the message, such as when sending from a custom application or for privacy reasons. 

Returns: 

  * A MsgOption function that can be used to customize the Msg instance.



####  func [WithPGPType](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L336) ¶ added in v0.3.9
    
    
    func WithPGPType(pgptype PGPType) MsgOption

WithPGPType sets the PGP type for the Msg during its creation or initialization, determining the encryption or signature method. 

This MsgOption function allows you to specify the PGP (Pretty Good Privacy) type to be used for securing the message. The chosen PGP type influences how the message is encrypted or signed, ensuring confidentiality and integrity of the content. This option should be called when creating a new Msg instance to set the desired PGP type appropriately. 

Parameters: 

  * pgptype: The PGPType value that specifies the desired PGP type for the Msg.



Returns: 

  * A MsgOption function that can be used to customize the Msg instance.



References: 

  * <https://datatracker.ietf.org/doc/html/rfc4880>



####  type [Option](https://github.com/wneessen/go-mail/blob/v0.7.2/client.go#L105) ¶
    
    
    type Option func(*Client) [error](/builtin#error)

Option is a function type that modifies the configuration or behavior of a Client instance. 

####  func [WithDSN](https://github.com/wneessen/go-mail/blob/v0.7.2/client.go#L605) ¶ added in v0.2.7
    
    
    func WithDSN() Option

WithDSN enables DSN (Delivery Status Notifications) for the Client as described in [RFC 1891](https://rfc-editor.org/rfc/rfc1891.html). 

This function configures the Client to request DSN, which provides status notifications for email delivery. DSN is only effective if the SMTP server supports it. By default, DSNMailReturnOption is set to DSNMailReturnFull, and DSNRcptNotifyOption is set to DSNRcptNotifySuccess and DSNRcptNotifyFailure. 

Returns: 

  * An Option function that enables DSN for the Client.



References: 

  * <https://datatracker.ietf.org/doc/html/rfc1891>



####  func [WithDSNMailReturnType](https://github.com/wneessen/go-mail/blob/v0.7.2/client.go#L629) ¶ added in v0.2.7
    
    
    func WithDSNMailReturnType(option DSNMailReturnOption) Option

WithDSNMailReturnType enables DSN (Delivery Status Notifications) for the Client as described in [RFC 1891](https://rfc-editor.org/rfc/rfc1891.html). 

This function configures the Client to request DSN and sets the DSNMailReturnOption to the provided value. DSN is only effective if the SMTP server supports it. The provided option must be either DSNMailReturnHeadersOnly or DSNMailReturnFull; otherwise, an error is returned. 

Parameters: 

  * option: The DSNMailReturnOption to be used (DSNMailReturnHeadersOnly or DSNMailReturnFull).



Returns: 

  * An Option function that sets the DSNMailReturnOption for the Client.
  * An error if an invalid DSNMailReturnOption is provided.



References: 

  * <https://datatracker.ietf.org/doc/html/rfc1891>



####  func [WithDSNRcptNotifyType](https://github.com/wneessen/go-mail/blob/v0.7.2/client.go#L660) ¶ added in v0.2.7
    
    
    func WithDSNRcptNotifyType(opts ...DSNRcptNotifyOption) Option

WithDSNRcptNotifyType enables DSN (Delivery Status Notifications) for the Client as described in [RFC 1891](https://rfc-editor.org/rfc/rfc1891.html). 

This function configures the Client to request DSN and sets the DSNRcptNotifyOption to the provided values. The provided options must be valid DSNRcptNotifyOption types. If DSNRcptNotifyNever is combined with any other notification type (such as DSNRcptNotifySuccess, DSNRcptNotifyFailure, or DSNRcptNotifyDelay), an error is returned. 

Parameters: 

  * opts: A variadic list of DSNRcptNotifyOption values (e.g., DSNRcptNotifySuccess, DSNRcptNotifyFailure).



Returns: 

  * An Option function that sets the DSNRcptNotifyOption for the Client.
  * An error if invalid DSNRcptNotifyOption values are provided or incompatible combinations are used.



References: 

  * <https://datatracker.ietf.org/doc/html/rfc1891>



####  func [WithDebugLog](https://github.com/wneessen/go-mail/blob/v0.7.2/client.go#L404) ¶ added in v0.3.9
    
    
    func WithDebugLog() Option

WithDebugLog enables debug logging for the Client. 

This function activates debug logging, which logs incoming and outgoing communication between the Client and the SMTP server to os.Stderr. By default the debug logging will redact any kind of SMTP authentication data. If you need access to the actual authentication data in your logs, you can enable authentication data logging with the WithLogAuthData option or by setting it with the Client.SetLogAuthData method. 

Returns: 

  * An Option function that enables debug logging for the Client.



####  func [WithDialContextFunc](https://github.com/wneessen/go-mail/blob/v0.7.2/client.go#L715) ¶ added in v0.4.0
    
    
    func WithDialContextFunc(dialCtxFunc DialContextFunc) Option

WithDialContextFunc sets the provided DialContextFunc as the DialContext for connecting to the SMTP server. 

This function overrides the default DialContext function used by the Client when establishing a connection to the SMTP server with the provided DialContextFunc. 

Parameters: 

  * dialCtxFunc: The custom DialContextFunc to be used for connecting to the SMTP server.



Returns: 

  * An Option function that sets the custom DialContextFunc for the Client.



####  func [WithHELO](https://github.com/wneessen/go-mail/blob/v0.7.2/client.go#L439) ¶
    
    
    func WithHELO(helo [string](/builtin#string)) Option

WithHELO sets the HELO/EHLO string used by the Client. 

This function configures the HELO/EHLO string sent by the Client when initiating communication with the SMTP server. By default, os.Hostname is used to identify the HELO/EHLO string. 

Parameters: 

  * helo: The string to be used for the HELO/EHLO greeting. Must not be empty.



Returns: 

  * An Option function that sets the HELO/EHLO string for the Client.
  * An error if the provided HELO string is empty.



####  func [WithLogAuthData](https://github.com/wneessen/go-mail/blob/v0.7.2/client.go#L734) ¶ added in v0.5.1
    
    
    func WithLogAuthData() Option

WithLogAuthData enables logging of authentication data. 

This function sets the logAuthData field of the Client to true, enabling the logging of authentication data. 

Be cautious when using this option, as the logs may include unencrypted authentication data, depending on the SMTP authentication method in use, which could pose a data protection risk. 

Returns: 

  * An Option function that configures the Client to enable authentication data logging.



####  func [WithLogger](https://github.com/wneessen/go-mail/blob/v0.7.2/client.go#L421) ¶ added in v0.3.9
    
    
    func WithLogger(logger [log](/github.com/wneessen/go-mail@v0.7.2/log).[Logger](/github.com/wneessen/go-mail@v0.7.2/log#Logger)) Option

WithLogger defines a custom logger for the Client. 

This function sets a custom logger for the Client, which must satisfy the log.Logger interface. The custom logger is used only when debug logging is enabled. By default, log.Stdlog is used if no custom logger is provided. 

Parameters: 

  * logger: A logger that satisfies the log.Logger interface.



Returns: 

  * An Option function that sets the custom logger for the Client.



####  func [WithPassword](https://github.com/wneessen/go-mail/blob/v0.7.2/client.go#L587) ¶
    
    
    func WithPassword(password [string](/builtin#string)) Option

WithPassword sets the password that the Client will use for SMTP authentication. 

This function configures the Client with the specified password for SMTP authentication. 

Important: 

  * Specifying a password with this option alone does NOT enable SMTP authentication.
  * To actually perform authentication with the server, you must also configure an authentication mechanism by using either WithSMTPAuth() or WithSMTPAuthCustom().
  * If you only call WithPassword() without setting an SMTP authentication method, the provided password will be stored but never used.



Parameters: 

  * password: The password to be used for SMTP authentication.



Returns: 

  * An Option function that sets the password for the Client.



####  func [WithPort](https://github.com/wneessen/go-mail/blob/v0.7.2/client.go#L330) ¶
    
    
    func WithPort(port [int](/builtin#int)) Option

WithPort sets the port number for the Client and overrides the default port. 

This function sets the specified port number for the Client, ensuring that the port number is valid (between 1 and 65535). If the provided port number is invalid, an error is returned. 

Parameters: 

  * port: The port number to be used by the Client. Must be between 1 and 65535.



Returns: 

  * An Option function that applies the port setting to the Client.
  * An error if the port number is outside the valid range.



####  func [WithSMTPAuth](https://github.com/wneessen/go-mail/blob/v0.7.2/client.go#L520) ¶
    
    
    func WithSMTPAuth(authtype SMTPAuthType) Option

WithSMTPAuth configures the Client to use the specified SMTPAuthType for SMTP authentication. 

This function sets the Client to use the specified SMTPAuthType for authenticating with the SMTP server. 

Parameters: 

  * authtype: The SMTPAuthType to be used for SMTP authentication.



Returns: 

  * An Option function that configures the Client to use the specified SMTPAuthType.



####  func [WithSMTPAuthCustom](https://github.com/wneessen/go-mail/blob/v0.7.2/client.go#L537) ¶
    
    
    func WithSMTPAuthCustom(smtpAuth [smtp](/github.com/wneessen/go-mail@v0.7.2/smtp).[Auth](/github.com/wneessen/go-mail@v0.7.2/smtp#Auth)) Option

WithSMTPAuthCustom sets a custom SMTP authentication mechanism for the Client. 

This function configures the Client to use a custom SMTP authentication mechanism. The provided mechanism must satisfy the smtp.Auth interface. 

Parameters: 

  * smtpAuth: The custom SMTP authentication mechanism, which must implement the smtp.Auth interface.



Returns: 

  * An Option function that sets the custom SMTP authentication for the Client.



####  func [WithSSL](https://github.com/wneessen/go-mail/blob/v0.7.2/client.go#L367) ¶
    
    
    func WithSSL() Option

WithSSL enables implicit SSL/TLS for the Client. 

This function configures the Client to use implicit SSL/TLS for secure communication. 

Returns: 

  * An Option function that enables SSL/TLS for the Client.



####  func [WithSSLPort](https://github.com/wneessen/go-mail/blob/v0.7.2/client.go#L387) ¶ added in v0.4.1
    
    
    func WithSSLPort(fallback [bool](/builtin#bool)) Option

WithSSLPort enables implicit SSL/TLS with an optional fallback for the Client. The correct port is automatically set. 

When this option is used with NewClient, the default port 25 is overridden with port 465 for SSL/TLS connections. If fallback is set to true and the SSL/TLS connection fails, the Client attempts to connect on port 25 using an unencrypted connection. If WithPort has already been used to set a different port, that port takes precedence, and the automatic fallback mechanism is skipped. 

Parameters: 

  * fallback: A boolean indicating whether to fall back to port 25 without SSL/TLS if the connection fails.



Returns: 

  * An Option function that enables SSL/TLS and configures the fallback mechanism for the Client.



####  func [WithTLSConfig](https://github.com/wneessen/go-mail/blob/v0.7.2/client.go#L501) ¶
    
    
    func WithTLSConfig(tlsconfig *[tls](/crypto/tls).[Config](/crypto/tls#Config)) Option

WithTLSConfig sets the tls.Config for the Client and overrides the default configuration. 

This function configures the Client with a custom tls.Config. It overrides the default TLS settings. An error is returned if the provided tls.Config is nil or invalid. 

Parameters: 

  * tlsconfig: A pointer to a tls.Config struct to be used for the Client. Must not be nil.



Returns: 

  * An Option function that sets the tls.Config for the Client.
  * An error if the provided tls.Config is invalid.



####  func [WithTLSPolicy](https://github.com/wneessen/go-mail/blob/v0.7.2/client.go#L462) ¶
    
    
    func WithTLSPolicy(policy TLSPolicy) Option

WithTLSPolicy sets the TLSPolicy of the Client and overrides the DefaultTLSPolicy. 

This function configures the Client's TLSPolicy, specifying how the Client handles TLS for SMTP connections. It overrides the default policy. For best practices regarding SMTP TLS connections, it is recommended to use WithTLSPortPolicy instead. 

Parameters: 

  * policy: The TLSPolicy to be applied to the Client.



Returns: 

  * An Option function that sets the TLSPolicy for the Client.



WithTLSPortPolicy instead. 

####  func [WithTLSPortPolicy](https://github.com/wneessen/go-mail/blob/v0.7.2/client.go#L483) ¶ added in v0.4.1
    
    
    func WithTLSPortPolicy(policy TLSPolicy) Option

WithTLSPortPolicy enables explicit TLS via STARTTLS for the Client using the provided TLSPolicy. The correct port is automatically set. 

When TLSMandatory or TLSOpportunistic is provided as the TLSPolicy, port 587 is used for the connection. If the connection fails with TLSOpportunistic, the Client attempts to connect on port 25 using an unencrypted connection as a fallback. If NoTLS is specified, the Client will always use port 25. If WithPort has already been used to set a different port, that port takes precedence, and the automatic fallback mechanism is skipped. 

Parameters: 

  * policy: The TLSPolicy to be used for STARTTLS communication.



Returns: 

  * An Option function that sets the TLSPortPolicy for the Client.



####  func [WithTimeout](https://github.com/wneessen/go-mail/blob/v0.7.2/client.go#L351) ¶
    
    
    func WithTimeout(timeout [time](/time).[Duration](/time#Duration)) Option

WithTimeout sets the connection timeout for the Client and overrides the default timeout. 

This function configures the Client with a specified connection timeout duration. It validates that the provided timeout is greater than zero. If the timeout is invalid, an error is returned. 

Parameters: 

  * timeout: The duration to be set as the connection timeout. Must be greater than zero.



Returns: 

  * An Option function that applies the timeout setting to the Client.
  * An error if the timeout duration is invalid.



####  func [WithUsername](https://github.com/wneessen/go-mail/blob/v0.7.2/client.go#L564) ¶
    
    
    func WithUsername(username [string](/builtin#string)) Option

WithUsername sets the username that the Client will use for SMTP authentication. 

This function configures the Client with the specified username for SMTP authentication. 

Important: 

  * Specifying a username with this option alone does NOT enable SMTP authentication.
  * To actually perform authentication with the server, you must also configure an authentication mechanism by using either WithSMTPAuth() or WithSMTPAuthCustom().
  * If you only call WithUsername() without setting an SMTP authentication method, the provided username will be stored but never used.



Parameters: 

  * username: The username to be used for SMTP authentication.



Returns: 

  * An Option function that sets the username for the Client.



####  func [WithoutNoop](https://github.com/wneessen/go-mail/blob/v0.7.2/client.go#L698) ¶ added in v0.3.6
    
    
    func WithoutNoop() Option

WithoutNoop indicates that the Client should skip the "NOOP" command during the dial. 

This option is useful for servers that delay potentially unwanted clients when they perform commands other than AUTH, such as Microsoft's Exchange Tarpit. 

Returns: 

  * An Option function that configures the Client to skip the "NOOP" command.



####  type [PGPType](https://github.com/wneessen/go-mail/blob/v0.7.2/msg.go#L83) ¶ added in v0.3.9
    
    
    type PGPType [int](/builtin#int)

PGPType is a type wrapper for an int, representing a type of PGP encryption or signature. 
    
    
    const (
    	// NoPGP indicates that a message should not be treated as PGP encrypted or signed and is the default value
    	// for a message
    	NoPGP PGPType = [iota](/builtin#iota)
    
    	// PGPEncrypt indicates that a message should be treated as PGP encrypted. This works closely together with
    	// the corresponding go-mail-middleware.
    	PGPEncrypt
    
    	// PGPSignature indicates that a message should be treated as PGP signed. This works closely together with
    	// the corresponding go-mail-middleware.
    	PGPSignature
    )

####  type [Part](https://github.com/wneessen/go-mail/blob/v0.7.2/part.go#L20) ¶
    
    
    type Part struct {
    	// contains filtered or unexported fields
    }

Part is a part of the Msg. 

This struct represents a single part of a multipart message. Each part has a content type, charset, optional description, encoding, and a function to write its content to an io.Writer. It also includes a flag to mark the part as deleted. 

####  func (*Part) [Delete](https://github.com/wneessen/go-mail/blob/v0.7.2/part.go#L176) ¶ added in v0.3.9
    
    
    func (p *Part) Delete()

Delete removes the current part from the parts list of the Msg by setting the isDeleted flag to true. 

This function marks the Part as deleted by setting the isDeleted flag to true. The msgWriter will skip over this Part during processing. 

####  func (*Part) [GetCharset](https://github.com/wneessen/go-mail/blob/v0.7.2/part.go#L52) ¶ added in v0.4.1
    
    
    func (p *Part) GetCharset() Charset

GetCharset returns the currently set Charset of the Part. 

This function returns the Charset that is currently set for the Part. 

Returns: 

  * The Charset of the Part.



####  func (*Part) [GetContent](https://github.com/wneessen/go-mail/blob/v0.7.2/part.go#L38) ¶ added in v0.3.0
    
    
    func (p *Part) GetContent() ([][byte](/builtin#byte), [error](/builtin#error))

GetContent executes the WriteFunc of the Part and returns the content as a byte slice. 

This function runs the part's writeFunc to write its content into a buffer and then returns the content as a byte slice. If an error occurs during the writing process, it is returned. 

Returns: 

  * A byte slice containing the part's content.
  * An error if the writeFunc encounters an issue.



####  func (*Part) [GetContentType](https://github.com/wneessen/go-mail/blob/v0.7.2/part.go#L62) ¶ added in v0.3.0
    
    
    func (p *Part) GetContentType() ContentType

GetContentType returns the currently set ContentType of the Part. 

This function returns the ContentType that is currently set for the Part. 

Returns: 

  * The ContentType of the Part.



####  func (*Part) [GetDescription](https://github.com/wneessen/go-mail/blob/v0.7.2/part.go#L94) ¶ added in v0.3.9
    
    
    func (p *Part) GetDescription() [string](/builtin#string)

GetDescription returns the currently set Content-Description of the Part. 

This function returns the Content-Description that is currently set for the Part. 

Returns: 

  * The Content-Description of the Part as a string.



####  func (*Part) [GetEncoding](https://github.com/wneessen/go-mail/blob/v0.7.2/part.go#L72) ¶ added in v0.3.0
    
    
    func (p *Part) GetEncoding() Encoding

GetEncoding returns the currently set Encoding of the Part. 

This function returns the Encoding that is currently set for the Part. 

Returns: 

  * The Encoding of the Part.



####  func (*Part) [GetWriteFunc](https://github.com/wneessen/go-mail/blob/v0.7.2/part.go#L84) ¶ added in v0.3.0
    
    
    func (p *Part) GetWriteFunc() func([io](/io).[Writer](/io#Writer)) ([int64](/builtin#int64), [error](/builtin#error))

GetWriteFunc returns the currently set WriteFunc of the Part. 

This function returns the WriteFunc that is currently set for the Part, which writes the part's content to an io.Writer. 

Returns: 

  * The WriteFunc of the Part, which is a function that takes an io.Writer and returns the number of bytes written and an error (if any).



####  func (*Part) [SetCharset](https://github.com/wneessen/go-mail/blob/v0.7.2/part.go#L126) ¶ added in v0.4.1
    
    
    func (p *Part) SetCharset(charset Charset)

SetCharset overrides the Charset of the Part. 

This function sets a new Charset for the Part, replacing the existing one. 

Parameters: 

  * charset: The new Charset to be set for the Part.



####  func (*Part) [SetContent](https://github.com/wneessen/go-mail/blob/v0.7.2/part.go#L105) ¶ added in v0.3.0
    
    
    func (p *Part) SetContent(content [string](/builtin#string))

SetContent overrides the content of the Part with the given string. 

This function sets the content of the Part by creating a new writeFunc that writes the provided string content to an io.Writer. 

Parameters: 

  * content: The string that will replace the current content of the Part.



####  func (*Part) [SetContentType](https://github.com/wneessen/go-mail/blob/v0.7.2/part.go#L116) ¶ added in v0.3.0
    
    
    func (p *Part) SetContentType(contentType ContentType)

SetContentType overrides the ContentType of the Part. 

This function sets a new ContentType for the Part, replacing the existing one. 

Parameters: 

  * contentType: The new ContentType to be set for the Part.



####  func (*Part) [SetDescription](https://github.com/wneessen/go-mail/blob/v0.7.2/part.go#L146) ¶ added in v0.3.9
    
    
    func (p *Part) SetDescription(description [string](/builtin#string))

SetDescription overrides the Content-Description of the Part. 

This function sets a new Content-Description for the Part, replacing the existing one. 

Parameters: 

  * description: The new Content-Description to be set for the Part.



####  func (*Part) [SetEncoding](https://github.com/wneessen/go-mail/blob/v0.7.2/part.go#L136) ¶
    
    
    func (p *Part) SetEncoding(encoding Encoding)

SetEncoding creates a new mime.WordEncoder based on the encoding setting of the message. 

This function sets a new Encoding for the Part, replacing the existing one. 

Parameters: 

  * encoding: The new Encoding to be set for the Part.



####  func (*Part) [SetIsSMIMESigned](https://github.com/wneessen/go-mail/blob/v0.7.2/part.go#L156) ¶ added in v0.6.0
    
    
    func (p *Part) SetIsSMIMESigned(smime [bool](/builtin#bool))

SetIsSMIMESigned sets the flag for signing the Part with S/MIME. 

This function updates the S/MIME signing flag for the Part. 

Parameters: 

  * smime: A boolean indicating whether the Part should be signed with S/MIME.



####  func (*Part) [SetWriteFunc](https://github.com/wneessen/go-mail/blob/v0.7.2/part.go#L168) ¶ added in v0.3.0
    
    
    func (p *Part) SetWriteFunc(writeFunc func([io](/io).[Writer](/io#Writer)) ([int64](/builtin#int64), [error](/builtin#error)))

SetWriteFunc overrides the WriteFunc of the Part. 

This function sets a new WriteFunc for the Part, replacing the existing one. The WriteFunc is responsible for writing the Part's content to an io.Writer. 

Parameters: 

  * writeFunc: A function that writes the Part's content to an io.Writer and returns the number of bytes written and an error (if any).



####  type [PartOption](https://github.com/wneessen/go-mail/blob/v0.7.2/part.go#L13) ¶
    
    
    type PartOption func(*Part)

PartOption returns a function that can be used for grouping Part options 

####  func [WithPartCharset](https://github.com/wneessen/go-mail/blob/v0.7.2/part.go#L190) ¶ added in v0.4.1
    
    
    func WithPartCharset(charset Charset) PartOption

WithPartCharset overrides the default Part charset. 

This function returns a PartOption that allows the charset of a Part to be overridden with the specified Charset. 

Parameters: 

  * charset: The Charset to be set for the Part.



Returns: 

  * A PartOption function that sets the Part's charset.



####  func [WithPartContentDescription](https://github.com/wneessen/go-mail/blob/v0.7.2/part.go#L222) ¶ added in v0.3.9
    
    
    func WithPartContentDescription(description [string](/builtin#string)) PartOption

WithPartContentDescription overrides the default Part Content-Description. 

This function returns a PartOption that allows the Content-Description of a Part to be overridden with the specified description. 

Parameters: 

  * description: The Content-Description to be set for the Part.



Returns: 

  * A PartOption function that sets the Part's Content-Description.



####  func [WithPartEncoding](https://github.com/wneessen/go-mail/blob/v0.7.2/part.go#L206) ¶
    
    
    func WithPartEncoding(encoding Encoding) PartOption

WithPartEncoding overrides the default Part encoding. 

This function returns a PartOption that allows the encoding of a Part to be overridden with the specified Encoding. 

Parameters: 

  * encoding: The Encoding to be set for the Part.



Returns: 

  * A PartOption function that sets the Part's encoding.



####  func [WithSMIMESigning](https://github.com/wneessen/go-mail/blob/v0.7.2/part.go#L234) ¶ added in v0.6.0
    
    
    func WithSMIMESigning() PartOption

WithSMIMESigning enables the S/MIME signing flag for a Part. 

This function provides a PartOption that overrides the S/MIME signing flag to enable signing. 

Returns: 

  * A PartOption that sets the S/MIME signing flag to true.



####  type [Reader](https://github.com/wneessen/go-mail/blob/v0.7.2/reader.go#L16) ¶ added in v0.3.2
    
    
    type Reader struct {
    	// contains filtered or unexported fields
    }

Reader is a type that implements the io.Reader interface for a Msg. 

This struct represents a reader that reads from a byte slice buffer. It keeps track of the current read position (offset) and any initialization error. The buffer holds the data to be read from the message. 

####  func (*Reader) [Error](https://github.com/wneessen/go-mail/blob/v0.7.2/reader.go#L29) ¶ added in v0.3.2
    
    
    func (r *Reader) Error() [error](/builtin#error)

Error returns an error if the Reader err field is not nil. 

This function checks the Reader's err field and returns it if it is not nil. If no error occurred during initialization, it returns nil. 

Returns: 

  * The error stored in the err field, or nil if no error is present.



####  func (*Reader) [Read](https://github.com/wneessen/go-mail/blob/v0.7.2/reader.go#L46) ¶ added in v0.3.2
    
    
    func (r *Reader) Read(payload [][byte](/builtin#byte)) (n [int](/builtin#int), err [error](/builtin#error))

Read reads the content of the Msg buffer into the provided payload to satisfy the io.Reader interface. 

This function reads data from the Reader's buffer into the provided byte slice (payload). It checks for errors or an empty buffer and resets the Reader if necessary. If no data is available, it returns io.EOF. Otherwise, it copies the content from the buffer into the payload and updates the read offset. 

Parameters: 

  * payload: A byte slice where the data will be copied.



Returns: 

  * n: The number of bytes copied into the payload.
  * err: An error if any issues occurred during the read operation or io.EOF if the buffer is empty.



####  func (*Reader) [Reset](https://github.com/wneessen/go-mail/blob/v0.7.2/reader.go#L66) ¶ added in v0.3.2
    
    
    func (r *Reader) Reset()

Reset resets the Reader buffer to be empty, but it retains the underlying storage for future use. 

This function clears the Reader's buffer by setting its length to 0 and resets the read offset to the beginning. The underlying storage is retained, allowing future writes to reuse the buffer. 

####  type [SMIME](https://github.com/wneessen/go-mail/blob/v0.7.2/smime.go#L35) ¶ added in v0.6.0
    
    
    type SMIME struct {
    	// contains filtered or unexported fields
    }

SMIME represents the configuration and state for S/MIME signing. 

This struct encapsulates the private key, certificate, optional intermediate certificate, and a flag indicating whether a signing process is currently in progress. 

Fields: 

  * privateKey: The private key used for signing (implements crypto.PrivateKey).
  * certificate: The x509 certificate associated with the private key.
  * intermediateCert: An optional x509 intermediate certificate for chain validation.
  * inProgress: A boolean flag indicating if a signing operation is currently active.



####  type [SMTPAuthType](https://github.com/wneessen/go-mail/blob/v0.7.2/auth.go#L15) ¶
    
    
    type SMTPAuthType [string](/builtin#string)

SMTPAuthType is a type wrapper for a string type. It represents the type of SMTP authentication mechanism to be used. 
    
    
    const (
    	// SMTPAuthCramMD5 is the "CRAM-MD5" SASL authentication mechanism as described in [RFC 4954](https://rfc-editor.org/rfc/rfc4954.html).
    	// <https://datatracker.ietf.org/doc/html/rfc4954/>
    	//
    	// CRAM-MD5 is not secure by modern standards. The vulnerabilities of MD5 and the lack of
    	// advanced security features make it inappropriate for protecting sensitive communications
    	// today.
    	//
    	// It was recommended to deprecate the standard in 20 November 2008. As an alternative it
    	// recommends e.g. SCRAM or SASL Plain protected by TLS instead.
    	//
    	// <https://datatracker.ietf.org/doc/html/draft-ietf-sasl-crammd5-to-historic-00.html>
    	SMTPAuthCramMD5 SMTPAuthType = "CRAM-MD5"
    
    	// SMTPAuthCustom is a custom SMTP AUTH mechanism provided by the user. If a user provides
    	// a custom smtp.Auth function to the Client, the Client will its smtpAuthType to this type.
    	//
    	// Do not use this SMTPAuthType without setting a custom smtp.Auth function on the Client.
    	SMTPAuthCustom SMTPAuthType = "CUSTOM"
    
    	// SMTPAuthLogin is the "LOGIN" SASL authentication mechanism. This authentication mechanism
    	// does not have an official RFC that could be followed. There is a spec by Microsoft and an
    	// IETF draft. The IETF draft is more lax than the MS spec, therefore we follow the I-D, which
    	// automatically matches the MS spec.
    	//
    	// Since the "LOGIN" SASL authentication mechanism transmits the username and password in
    	// plaintext over the internet connection, we only allow this mechanism over a TLS secured
    	// connection.
    	//
    	// <https://msopenspecs.azureedge.net/files/MS-XLOGIN/%5bMS-XLOGIN%5d.pdf>
    	//
    	// <https://datatracker.ietf.org/doc/html/draft-murchison-sasl-login-00>
    	SMTPAuthLogin SMTPAuthType = "LOGIN"
    
    	// SMTPAuthLoginNoEnc is the "LOGIN" SASL authentication mechanism. This authentication mechanism
    	// does not have an official RFC that could be followed. There is a spec by Microsoft and an
    	// IETF draft. The IETF draft is more lax than the MS spec, therefore we follow the I-D, which
    	// automatically matches the MS spec.
    	//
    	// Since the "LOGIN" SASL authentication mechanism transmits the username and password in
    	// plaintext over the internet connection, by default we only allow this mechanism over
    	// a TLS secured connection. This authentiation mechanism overrides this default and will
    	// allow LOGIN authentication via an unencrypted channel. This can be useful if the
    	// connection has already been secured in a different way (e. g. a SSH tunnel)
    	//
    	// Note: Use this authentication method with caution. If used in the wrong way, you might
    	// expose your authentication information over unencrypted channels!
    	//
    	// <https://msopenspecs.azureedge.net/files/MS-XLOGIN/%5bMS-XLOGIN%5d.pdf>
    	//
    	// <https://datatracker.ietf.org/doc/html/draft-murchison-sasl-login-00>
    	SMTPAuthLoginNoEnc SMTPAuthType = "LOGIN-NOENC"
    
    	// SMTPAuthNoAuth is equivalent to performing no authentication at all. It is a convenience
    	// option and should not be used. Instead, for mail servers that do no support/require
    	// authentication, the Client should not be passed the WithSMTPAuth option at all.
    	SMTPAuthNoAuth SMTPAuthType = "NOAUTH"
    
    	// SMTPAuthPlain is the "PLAIN" authentication mechanism as described in [RFC 4616](https://rfc-editor.org/rfc/rfc4616.html).
    	//
    	// Since the "PLAIN" SASL authentication mechanism transmits the username and password in
    	// plaintext over the internet connection, we only allow this mechanism over a TLS secured
    	// connection.
    	//
    	// <https://datatracker.ietf.org/doc/html/rfc4616/>
    	SMTPAuthPlain SMTPAuthType = "PLAIN"
    
    	// SMTPAuthPlainNoEnc is the "PLAIN" authentication mechanism as described in [RFC 4616](https://rfc-editor.org/rfc/rfc4616.html).
    	//
    	// Since the "PLAIN" SASL authentication mechanism transmits the username and password in
    	// plaintext over the internet connection, by default we only allow this mechanism over
    	// a TLS secured connection. This authentiation mechanism overrides this default and will
    	// allow PLAIN authentication via an unencrypted channel. This can be useful if the
    	// connection has already been secured in a different way (e. g. a SSH tunnel)
    	//
    	// Note: Use this authentication method with caution. If used in the wrong way, you might
    	// expose your authentication information over unencrypted channels!
    	//
    	// <https://datatracker.ietf.org/doc/html/rfc4616/>
    	SMTPAuthPlainNoEnc SMTPAuthType = "PLAIN-NOENC"
    
    	// SMTPAuthXOAUTH2 is the "XOAUTH2" SASL authentication mechanism.
    	// <https://developers.google.com/gmail/imap/xoauth2-protocol>
    	SMTPAuthXOAUTH2 SMTPAuthType = "XOAUTH2"
    
    	// SMTPAuthSCRAMSHA1 is the "SCRAM-SHA-1" SASL authentication mechanism as described in [RFC 5802](https://rfc-editor.org/rfc/rfc5802.html).
    	//
    	// SCRAM-SHA-1 is still considered secure for certain applications, particularly when used as part
    	// of a challenge-response authentication mechanism (as we use it). However, it is generally
    	// recommended to prefer stronger alternatives like SCRAM-SHA-256(-PLUS), as SHA-1 has known
    	// vulnerabilities in other contexts, although it remains effective in HMAC constructions.
    	//
    	// <https://datatracker.ietf.org/doc/html/rfc5802>
    	SMTPAuthSCRAMSHA1 SMTPAuthType = "SCRAM-SHA-1"
    
    	// SMTPAuthSCRAMSHA1PLUS is the "SCRAM-SHA-1-PLUS" SASL authentication mechanism as described in [RFC 5802](https://rfc-editor.org/rfc/rfc5802.html).
    	//
    	// SCRAM-SHA-X-PLUS authentication require TLS channel bindings to protect against MitM attacks and
    	// to guarantee that the integrity of the transport layer is preserved throughout the authentication
    	// process. Therefore we only allow this mechanism over a TLS secured connection.
    	//
    	// SCRAM-SHA-1-PLUS is still considered secure for certain applications, particularly when used as part
    	// of a challenge-response authentication mechanism (as we use it). However, it is generally
    	// recommended to prefer stronger alternatives like SCRAM-SHA-256(-PLUS), as SHA-1 has known
    	// vulnerabilities in other contexts, although it remains effective in HMAC constructions.
    	//
    	// <https://datatracker.ietf.org/doc/html/rfc5802>
    	SMTPAuthSCRAMSHA1PLUS SMTPAuthType = "SCRAM-SHA-1-PLUS"
    
    	// SMTPAuthSCRAMSHA256 is the "SCRAM-SHA-256" SASL authentication mechanism as described in [RFC 7677](https://rfc-editor.org/rfc/rfc7677.html).
    	//
    	// <https://datatracker.ietf.org/doc/html/rfc7677>
    	SMTPAuthSCRAMSHA256 SMTPAuthType = "SCRAM-SHA-256"
    
    	// SMTPAuthSCRAMSHA256PLUS is the "SCRAM-SHA-256-PLUS" SASL authentication mechanism as described in [RFC 7677](https://rfc-editor.org/rfc/rfc7677.html).
    	//
    	// SCRAM-SHA-X-PLUS authentication require TLS channel bindings to protect against MitM attacks and
    	// to guarantee that the integrity of the transport layer is preserved throughout the authentication
    	// process. Therefore we only allow this mechanism over a TLS secured connection.
    	//
    	// <https://datatracker.ietf.org/doc/html/rfc7677>
    	SMTPAuthSCRAMSHA256PLUS SMTPAuthType = "SCRAM-SHA-256-PLUS"
    
    	// SMTPAuthAutoDiscover is a mechanism that dynamically discovers all authentication mechanisms
    	// supported by the SMTP server and selects the strongest available one.
    	//
    	// This type simplifies authentication by automatically negotiating the most secure mechanism
    	// offered by the server, based on a predefined security ranking. For instance, mechanisms like
    	// SCRAM-SHA-256(-PLUS) or XOAUTH2 are prioritized over weaker mechanisms such as CRAM-MD5 or PLAIN.
    	//
    	// The negotiation process ensures that mechanisms requiring additional capabilities (e.g.,
    	// SCRAM-SHA-X-PLUS with TLS channel binding) are only selected when the necessary prerequisites
    	// are in place, such as an active TLS-secured connection.
    	//
    	// By automating mechanism selection, SMTPAuthAutoDiscover minimizes configuration effort while
    	// maximizing security and compatibility with a wide range of SMTP servers.
    	SMTPAuthAutoDiscover SMTPAuthType = "AUTODISCOVER"
    )

####  func (*SMTPAuthType) [UnmarshalString](https://github.com/wneessen/go-mail/blob/v0.7.2/auth.go#L197) ¶ added in v0.5.2
    
    
    func (sa *SMTPAuthType) UnmarshalString(value [string](/builtin#string)) [error](/builtin#error)

UnmarshalString satisfies the fig.StringUnmarshaler interface for the SMTPAuthType type <https://pkg.go.dev/github.com/kkyr/fig#StringUnmarshaler>

####  type [SendErrReason](https://github.com/wneessen/go-mail/blob/v0.7.2/senderror.go#L75) ¶ added in v0.3.7
    
    
    type SendErrReason [int](/builtin#int)

SendErrReason represents a comparable reason on why the delivery failed 
    
    
    const (
    	// ErrGetSender is returned if the Msg.GetSender method fails during a Client.Send
    	ErrGetSender SendErrReason = [iota](/builtin#iota)
    
    	// ErrGetRcpts is returned if the Msg.GetRecipients method fails during a Client.Send
    	ErrGetRcpts
    
    	// ErrSMTPMailFrom is returned if the Msg delivery failed when sending the MAIL FROM command
    	// to the sending SMTP server
    	ErrSMTPMailFrom
    
    	// ErrSMTPRcptTo is returned if the Msg delivery failed when sending the RCPT TO command
    	// to the sending SMTP server
    	ErrSMTPRcptTo
    
    	// ErrSMTPData is returned if the Msg delivery failed when sending the DATA command
    	// to the sending SMTP server
    	ErrSMTPData
    
    	// ErrSMTPDataClose is returned if the Msg delivery failed when trying to close the
    	// Client data writer
    	ErrSMTPDataClose
    
    	// ErrSMTPReset is returned if the Msg delivery failed when sending the RSET command
    	// to the sending SMTP server
    	ErrSMTPReset
    
    	// ErrWriteContent is returned if the Msg delivery failed when sending Msg content
    	// to the Client writer
    	ErrWriteContent
    
    	// ErrConnCheck is returned if the Msg delivery failed when checking if the SMTP
    	// server connection is still working
    	ErrConnCheck
    
    	// ErrNoUnencoded is returned if the Msg delivery failed when the Msg is configured for
    	// unencoded delivery but the server does not support this
    	ErrNoUnencoded
    
    	// ErrAmbiguous is a generalized delivery error for the SendError type that is
    	// returned if the exact reason for the delivery failure is ambiguous
    	ErrAmbiguous
    )

List of SendError reasons 

####  func (SendErrReason) [String](https://github.com/wneessen/go-mail/blob/v0.7.2/senderror.go#L226) ¶ added in v0.3.7
    
    
    func (r SendErrReason) String() [string](/builtin#string)

String satisfies the fmt.Stringer interface for the SendErrReason type. 

This function converts the SendErrReason into a human-readable string representation based on the error type. If the error reason does not match any predefined case, it returns "unknown reason". 

Returns: 

  * A string representation of the SendErrReason.



####  type [SendError](https://github.com/wneessen/go-mail/blob/v0.7.2/senderror.go#L64) ¶ added in v0.3.7
    
    
    type SendError struct {
    	Reason SendErrReason
    	// contains filtered or unexported fields
    }

SendError is an error wrapper for delivery errors of the Msg. 

This struct represents an error that occurs during the delivery of a message. It holds details about the affected message, a list of errors, the recipient list, and whether the error is temporary or permanent. It also includes a reason code for the error. 

####  func (*SendError) [EnhancedStatusCode](https://github.com/wneessen/go-mail/blob/v0.7.2/senderror.go#L196) ¶ added in v0.6.0
    
    
    func (e *SendError) EnhancedStatusCode() [string](/builtin#string)

EnhancedStatusCode returns the enhanced status code of the server response if the server supports it, as described in [RFC 2034](https://rfc-editor.org/rfc/rfc2034.html). 

This function retrieves the enhanced status code of an error returned by the server. This requires that the receiving server supports this SMTP extension as described in [RFC 2034](https://rfc-editor.org/rfc/rfc2034.html). Since this is the SendError interface, we only collect status codes for error responses, meaning 4xx or 5xx. If the server does not support the ENHANCEDSTATUSCODES extension or the error did not include an enhanced status code, it will return an empty string. 

Returns: 

  * The enhanced status code as returned by the server, or an empty string is not supported.



References: 

  * <https://datatracker.ietf.org/doc/html/rfc2034>



####  func (*SendError) [Error](https://github.com/wneessen/go-mail/blob/v0.7.2/senderror.go#L87) ¶ added in v0.3.7
    
    
    func (e *SendError) Error() [string](/builtin#string)

Error implements the error interface for the SendError type. 

This function returns a detailed error message string for the SendError, including the reason for failure, list of errors, affected recipients, and the message ID of the affected message (if available). If the reason is unknown (greater than 10), it returns "unknown reason". The error message is built dynamically based on the content of the error list, recipient list, and message ID. 

Returns: 

  * A string representing the error message.



####  func (*SendError) [ErrorCode](https://github.com/wneessen/go-mail/blob/v0.7.2/senderror.go#L211) ¶ added in v0.6.0
    
    
    func (e *SendError) ErrorCode() [int](/builtin#int)

ErrorCode returns the error code of the server response. 

This function retrieves the error code the error returned by the server. The error code will start with 5 on permanent errors and with 4 on a temporary error. If the error is not returned by the server, but is generated by go-mail, the code will be 0. 

Returns: 

  * The error code as returned by the server, or 0 if not a server error.



####  func (*SendError) [Is](https://github.com/wneessen/go-mail/blob/v0.7.2/senderror.go#L132) ¶ added in v0.3.7
    
    
    func (e *SendError) Is(errType [error](/builtin#error)) [bool](/builtin#bool)

Is implements the errors.Is functionality and compares the SendErrReason. 

This function allows for comparison between two errors by checking if the provided error matches the SendError type and, if so, compares the SendErrReason and the temporary status (isTemp) of both errors. 

Parameters: 

  * errType: The error to compare against the current SendError.



Returns: 

  * true if the errors have the same reason and temporary status, false otherwise.



####  func (*SendError) [IsTemp](https://github.com/wneessen/go-mail/blob/v0.7.2/senderror.go#L147) ¶ added in v0.3.7
    
    
    func (e *SendError) IsTemp() [bool](/builtin#bool)

IsTemp returns true if the delivery error is of a temporary nature and can be retried. 

This function checks whether the SendError indicates a temporary error, which suggests that the delivery can be retried. If the SendError is nil, it returns false. 

Returns: 

  * true if the error is temporary, false otherwise.



####  func (*SendError) [MessageID](https://github.com/wneessen/go-mail/blob/v0.7.2/senderror.go#L161) ¶ added in v0.5.0
    
    
    func (e *SendError) MessageID() [string](/builtin#string)

MessageID returns the message ID of the affected Msg that caused the error. 

This function retrieves the message ID of the Msg associated with the SendError. If no message ID was set or if the SendError or Msg is nil, it returns an empty string. 

Returns: 

  * The message ID as a string, or an empty string if no ID is available.



####  func (*SendError) [Msg](https://github.com/wneessen/go-mail/blob/v0.7.2/senderror.go#L175) ¶ added in v0.5.0
    
    
    func (e *SendError) Msg() *Msg

Msg returns the pointer to the affected message that caused the error. 

This function retrieves the Msg associated with the SendError. If the SendError or the affectedMsg is nil, it returns nil. 

Returns: 

  * A pointer to the Msg that caused the error, or nil if not available.



####  type [TLSPolicy](https://github.com/wneessen/go-mail/blob/v0.7.2/tls.go#L8) ¶
    
    
    type TLSPolicy [int](/builtin#int)

TLSPolicy is a type wrapper for an int type and describes the different TLS policies we allow. 
    
    
    const (
    	// TLSMandatory requires that the connection to the server is
    	// encrypting using STARTTLS. If the server does not support STARTTLS
    	// the connection will be terminated with an error.
    	TLSMandatory TLSPolicy = [iota](/builtin#iota)
    
    	// TLSOpportunistic tries to establish an encrypted connection via the
    	// STARTTLS protocol. If the server does not support this, it will fall
    	// back to non-encrypted plaintext transmission.
    	TLSOpportunistic
    
    	// NoTLS forces the transaction to be not encrypted.
    	NoTLS
    )

####  func (TLSPolicy) [String](https://github.com/wneessen/go-mail/blob/v0.7.2/tls.go#L33) ¶
    
    
    func (p TLSPolicy) String() [string](/builtin#string)

String satisfies the fmt.Stringer interface for the TLSPolicy type. 

This function returns a string representation of the TLSPolicy. It matches the policy value to predefined constants and returns the corresponding string. If the policy does not match any known values, it returns "UnknownPolicy". 

Returns: 

  * A string representing the TLSPolicy.


