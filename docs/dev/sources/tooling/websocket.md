# coder/websocket

> Auto-fetched from [https://pkg.go.dev/github.com/coder/websocket](https://pkg.go.dev/github.com/coder/websocket)
> Last Updated: 2026-01-29T20:11:52.461378+00:00

---

Overview
¶
Wasm
Package websocket implements the
RFC 6455
WebSocket protocol.
https://tools.ietf.org/html/rfc6455
Use Dial to dial a WebSocket server.
Use Accept to accept a WebSocket client.
Conn represents the resulting WebSocket connection.
The examples are the best way to understand how to correctly use the library.
The wsjson subpackage contain helpers for JSON and protobuf messages.
More documentation at
https://github.com/coder/websocket
.
Wasm
¶
The client side supports compiling to Wasm.
It wraps the WebSocket browser API.
See
https://developer.mozilla.org/en-US/docs/Web/API/WebSocket
Some important caveats to be aware of:
Accept always errors out
Conn.Ping is no-op
Conn.CloseNow is Close(StatusGoingAway, "")
HTTPClient, HTTPHeader and CompressionMode in DialOptions are no-op
*http.Response from Dial is &http.Response{} with a 101 status code on success
Example (CrossOrigin)
¶
package main

import (
"log"
"net/http"

"github.com/coder/websocket"
)

func main() {
// This handler demonstrates how to safely accept cross origin WebSockets
// from the origin example.com.
fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
c, err := websocket.Accept(w, r, &websocket.AcceptOptions{
OriginPatterns: []string{"example.com"},
})
if err != nil {
log.Println(err)
return
}
c.Close(websocket.StatusNormalClosure, "cross origin WebSocket accepted")
})

err := http.ListenAndServe("localhost:8080", fn)
log.Fatal(err)
}
Share
Format
Run
Example (Echo)
¶
This example demonstrates a echo server.
package main

import ()

func main() {
// https://github.com/nhooyr/websocket/tree/master/internal/examples/echo
}
Share
Format
Run
Example (FullStackChat)
¶
This example demonstrates full stack chat with an automated test.
package main

import ()

func main() {
// https://github.com/nhooyr/websocket/tree/master/internal/examples/chat
}
Share
Format
Run
Example (WriteOnly)
¶
package main

import (
"context"
"log"
"net/http"
"time"

"github.com/coder/websocket"
"github.com/coder/websocket/wsjson"
)

func main() {
// This handler demonstrates how to correctly handle a write only WebSocket connection.
// i.e you only expect to write messages and do not expect to read any messages.
fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
c, err := websocket.Accept(w, r, nil)
if err != nil {
log.Println(err)
return
}
defer c.CloseNow()

ctx, cancel := context.WithTimeout(r.Context(), time.Minute*10)
defer cancel()

ctx = c.CloseRead(ctx)

t := time.NewTicker(time.Second * 30)
defer t.Stop()

for {
select {
case <-ctx.Done():
c.Close(websocket.StatusNormalClosure, "")
return
case <-t.C:
err = wsjson.Write(ctx, c, "hi")
if err != nil {
log.Println(err)
return
}
}
}
})

err := http.ListenAndServe("localhost:8080", fn)
log.Fatal(err)
}
Share
Format
Run
Index
¶
Variables
func NetConn(ctx context.Context, c *Conn, msgType MessageType) net.Conn
type AcceptOptions
type CloseError
func (ce CloseError) Error() string
type CompressionMode
type Conn
func Accept(w http.ResponseWriter, r *http.Request, opts *AcceptOptions) (*Conn, error)
func Dial(ctx context.Context, u string, opts *DialOptions) (*Conn, *http.Response, error)
func (c *Conn) Close(code StatusCode, reason string) (err error)
func (c *Conn) CloseNow() (err error)
func (c *Conn) CloseRead(ctx context.Context) context.Context
func (c *Conn) Ping(ctx context.Context) error
func (c *Conn) Read(ctx context.Context) (MessageType, []byte, error)
func (c *Conn) Reader(ctx context.Context) (MessageType, io.Reader, error)
func (c *Conn) SetReadLimit(n int64)
func (c *Conn) Subprotocol() string
func (c *Conn) Write(ctx context.Context, typ MessageType, p []byte) error
func (c *Conn) Writer(ctx context.Context, typ MessageType) (io.WriteCloser, error)
type DialOptions
type MessageType
func (i MessageType) String() string
type StatusCode
func CloseStatus(err error) StatusCode
func (i StatusCode) String() string
Examples
¶
Package (CrossOrigin)
Package (Echo)
Package (FullStackChat)
Package (WriteOnly)
Accept
CloseStatus
Conn.Ping
Dial
Constants
¶
This section is empty.
Variables
¶
View Source
var ErrMessageTooBig =
errors
.
New
("websocket: message too big")
ErrMessageTooBig is returned when a message exceeds the read limit.
Functions
¶
func
NetConn
¶
func NetConn(ctx
context
.
Context
, c *
Conn
, msgType
MessageType
)
net
.
Conn
NetConn converts a *websocket.Conn into a net.Conn.
It's for tunneling arbitrary protocols over WebSockets.
Few users of the library will need this but it's tricky to implement
correctly and so provided in the library.
See
https://github.com/nhooyr/websocket/issues/100
.
Every Write to the net.Conn will correspond to a message write of
the given type on *websocket.Conn.
The passed ctx bounds the lifetime of the net.Conn. If cancelled,
all reads and writes on the net.Conn will be cancelled.
If a message is read that is not of the correct type, the connection
will be closed with StatusUnsupportedData and an error will be returned.
Close will close the *websocket.Conn with StatusNormalClosure.
When a deadline is hit and there is an active read or write goroutine, the
connection will be closed. This is different from most net.Conn implementations
where only the reading/writing goroutines are interrupted but the connection
is kept alive.
The Addr methods will return the real addresses for connections obtained
from websocket.Accept. But for connections obtained from websocket.Dial, a mock net.Addr
will be returned that gives "websocket" for Network() and "websocket/unknown-addr" for
String(). This is because websocket.Dial only exposes a io.ReadWriteCloser instead of the
full net.Conn to us.
When running as WASM, the Addr methods will always return the mock address described above.
A received StatusNormalClosure or StatusGoingAway close frame will be translated to
io.EOF when reading.
Furthermore, the ReadLimit is set to -1 to disable it.
Types
¶
type
AcceptOptions
¶
type AcceptOptions struct {
// Subprotocols lists the WebSocket subprotocols that Accept will negotiate with the client.
// The empty subprotocol will always be negotiated as per
RFC 6455
. If you would like to
// reject it, close the connection when c.Subprotocol() == "".
Subprotocols []
string
// InsecureSkipVerify is used to disable Accept's origin verification behaviour.
//
// You probably want to use OriginPatterns instead.
InsecureSkipVerify
bool
// OriginPatterns lists the host patterns for authorized origins.
// The request host is always authorized.
// Use this to enable cross origin WebSockets.
//
// i.e javascript running on example.com wants to access a WebSocket server at chat.example.com.
// In such a case, example.com is the origin and chat.example.com is the request host.
// One would set this field to []string{"example.com"} to authorize example.com to connect.
//
// Each pattern is matched case insensitively with path.Match (see
//
https://golang.org/pkg/path/#Match
). By default, it is matched
// against the request origin host. If the pattern contains a URI
// scheme ("://"), it will be matched against "scheme://host".
//
// Please ensure you understand the ramifications of enabling this.
// If used incorrectly your WebSocket server will be open to CSRF attacks.
//
// Do not use * as a pattern to allow any origin, prefer to use InsecureSkipVerify instead
// to bring attention to the danger of such a setting.
OriginPatterns []
string
// CompressionMode controls the compression mode.
// Defaults to CompressionDisabled.
//
// See docs on CompressionMode for details.
CompressionMode
CompressionMode
// CompressionThreshold controls the minimum size of a message before compression is applied.
//
// Defaults to 512 bytes for CompressionNoContextTakeover and 128 bytes
// for CompressionContextTakeover.
CompressionThreshold
int
// OnPingReceived is an optional callback invoked synchronously when a ping frame is received.
//
// The payload contains the application data of the ping frame.
// If the callback returns false, the subsequent pong frame will not be sent.
// To avoid blocking, any expensive processing should be performed asynchronously using a goroutine.
OnPingReceived func(ctx
context
.
Context
, payload []
byte
)
bool
// OnPongReceived is an optional callback invoked synchronously when a pong frame is received.
//
// The payload contains the application data of the pong frame.
// To avoid blocking, any expensive processing should be performed asynchronously using a goroutine.
//
// Unlike OnPingReceived, this callback does not return a value because a pong frame
// is a response to a ping and does not trigger any further frame transmission.
OnPongReceived func(ctx
context
.
Context
, payload []
byte
)
}
AcceptOptions represents Accept's options.
type
CloseError
¶
type CloseError struct {
Code
StatusCode
Reason
string
}
CloseError is returned when the connection is closed with a status and reason.
Use Go 1.13's errors.As to check for this error.
Also see the CloseStatus helper.
func (CloseError)
Error
¶
func (ce
CloseError
) Error()
string
type
CompressionMode
¶
type CompressionMode
int
CompressionMode represents the modes available to the permessage-deflate extension.
See
https://tools.ietf.org/html/rfc7692
Works in all modern browsers except Safari which does not implement the permessage-deflate extension.
Compression is only used if the peer supports the mode selected.
const (
// CompressionDisabled disables the negotiation of the permessage-deflate extension.
//
// This is the default. Do not enable compression without benchmarking for your particular use case first.
CompressionDisabled
CompressionMode
=
iota
// CompressionContextTakeover compresses each message greater than 128 bytes reusing the 32 KB sliding window from
// previous messages. i.e compression context across messages is preserved.
//
// As most WebSocket protocols are text based and repetitive, this compression mode can be very efficient.
//
// The memory overhead is a fixed 32 KB sliding window, a fixed 1.2 MB flate.Writer and a sync.Pool of 40 KB flate.Reader's
// that are used when reading and then returned.
//
// Thus, it uses more memory than CompressionNoContextTakeover but compresses more efficiently.
//
// If the peer does not support CompressionContextTakeover then we will fall back to CompressionNoContextTakeover.
CompressionContextTakeover
// CompressionNoContextTakeover compresses each message greater than 512 bytes. Each message is compressed with
// a new 1.2 MB flate.Writer pulled from a sync.Pool. Each message is read with a 40 KB flate.Reader pulled from
// a sync.Pool.
//
// This means less efficient compression as the sliding window from previous messages will not be used but the
// memory overhead will be lower as there will be no fixed cost for the flate.Writer nor the 32 KB sliding window.
// Especially if the connections are long lived and seldom written to.
//
// Thus, it uses less memory than CompressionContextTakeover but compresses less efficiently.
//
// If the peer does not support CompressionNoContextTakeover then we will fall back to CompressionDisabled.
CompressionNoContextTakeover
)
type
Conn
¶
type Conn struct {
// contains filtered or unexported fields
}
Conn represents a WebSocket connection.
All methods may be called concurrently except for Reader and Read.
You must always read from the connection. Otherwise control
frames will not be handled. See Reader and CloseRead.
Be sure to call Close on the connection when you
are finished with it to release associated resources.
On any error from any method, the connection is closed
with an appropriate reason.
This applies to context expirations as well unfortunately.
See
https://github.com/nhooyr/websocket/issues/242#issuecomment-633182220
func
Accept
¶
func Accept(w
http
.
ResponseWriter
, r *
http
.
Request
, opts *
AcceptOptions
) (*
Conn
,
error
)
Accept accepts a WebSocket handshake from a client and upgrades the
the connection to a WebSocket.
Accept will not allow cross origin requests by default.
See the InsecureSkipVerify and OriginPatterns options to allow cross origin requests.
Accept will write a response to w on all errors.
Note that using the http.Request Context after Accept returns may lead to
unexpected behavior (see http.Hijacker).
Example
¶
package main

import (
"context"
"log"
"net/http"
"time"

"github.com/coder/websocket"
"github.com/coder/websocket/wsjson"
)

func main() {
// This handler accepts a WebSocket connection, reads a single JSON
// message from the client and then closes the connection.

fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
c, err := websocket.Accept(w, r, nil)
if err != nil {
log.Println(err)
return
}
defer c.CloseNow()

ctx, cancel := context.WithTimeout(r.Context(), time.Second*10)
defer cancel()

var v any
err = wsjson.Read(ctx, c, &v)
if err != nil {
log.Println(err)
return
}

c.Close(websocket.StatusNormalClosure, "")
})

err := http.ListenAndServe("localhost:8080", fn)
log.Fatal(err)
}
Share
Format
Run
func
Dial
¶
func Dial(ctx
context
.
Context
, u
string
, opts *
DialOptions
) (*
Conn
, *
http
.
Response
,
error
)
Dial performs a WebSocket handshake on url.
The response is the WebSocket handshake response from the server.
You never need to close resp.Body yourself.
If an error occurs, the returned response may be non nil.
However, you can only read the first 1024 bytes of the body.
This function requires at least Go 1.12 as it uses a new feature
in net/http to perform WebSocket handshakes.
See docs on the HTTPClient option and
https://github.com/golang/go/issues/26937#issuecomment-415855861
URLs with http/https schemes will work and are interpreted as ws/wss.
Example
¶
package main

import (
"context"
"log"
"time"

"github.com/coder/websocket"
"github.com/coder/websocket/wsjson"
)

func main() {
// Dials a server, writes a single JSON message and then
// closes the connection.

ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
defer cancel()

c, _, err := websocket.Dial(ctx, "ws://localhost:8080", nil)
if err != nil {
log.Fatal(err)
}
defer c.CloseNow()

err = wsjson.Write(ctx, c, "hi")
if err != nil {
log.Fatal(err)
}

c.Close(websocket.StatusNormalClosure, "")
}
Share
Format
Run
func (*Conn)
Close
¶
func (c *
Conn
) Close(code
StatusCode
, reason
string
) (err
error
)
Close performs the WebSocket close handshake with the given status code and reason.
It will write a WebSocket close frame with a timeout of 5s and then wait 5s for
the peer to send a close frame.
All data messages received from the peer during the close handshake will be discarded.
The connection can only be closed once. Additional calls to Close
are no-ops.
The maximum length of reason must be 125 bytes. Avoid sending a dynamic reason.
Close will unblock all goroutines interacting with the connection once
complete.
func (*Conn)
CloseNow
¶
func (c *
Conn
) CloseNow() (err
error
)
CloseNow closes the WebSocket connection without attempting a close handshake.
Use when you do not want the overhead of the close handshake.
func (*Conn)
CloseRead
¶
func (c *
Conn
) CloseRead(ctx
context
.
Context
)
context
.
Context
CloseRead starts a goroutine to read from the connection until it is closed
or a data message is received.
Once CloseRead is called you cannot read any messages from the connection.
The returned context will be cancelled when the connection is closed.
If a data message is received, the connection will be closed with StatusPolicyViolation.
Call CloseRead when you do not expect to read any more messages.
Since it actively reads from the connection, it will ensure that ping, pong and close
frames are responded to. This means c.Ping and c.Close will still work as expected.
This function is idempotent.
func (*Conn)
Ping
¶
func (c *
Conn
) Ping(ctx
context
.
Context
)
error
Ping sends a ping to the peer and waits for a pong.
Use this to measure latency or ensure the peer is responsive.
Ping must be called concurrently with Reader as it does
not read from the connection but instead waits for a Reader call
to read the pong.
TCP Keepalives should suffice for most use cases.
Example
¶
package main

import (
"context"
"log"
"time"

"github.com/coder/websocket"
)

func main() {
// Dials a server and pings it 5 times.

ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
defer cancel()

c, _, err := websocket.Dial(ctx, "ws://localhost:8080", nil)
if err != nil {
log.Fatal(err)
}
defer c.CloseNow()

// Required to read the Pongs from the server.
ctx = c.CloseRead(ctx)

for range 5 {
err = c.Ping(ctx)
if err != nil {
log.Fatal(err)
}
}

c.Close(websocket.StatusNormalClosure, "")
}
Share
Format
Run
func (*Conn)
Read
¶
func (c *
Conn
) Read(ctx
context
.
Context
) (
MessageType
, []
byte
,
error
)
Read is a convenience method around Reader to read a single message
from the connection.
func (*Conn)
Reader
¶
func (c *
Conn
) Reader(ctx
context
.
Context
) (
MessageType
,
io
.
Reader
,
error
)
Reader reads from the connection until there is a WebSocket
data message to be read. It will handle ping, pong and close frames as appropriate.
It returns the type of the message and an io.Reader to read it.
The passed context will also bound the reader.
Ensure you read to EOF otherwise the connection will hang.
Call CloseRead if you do not expect any data messages from the peer.
Only one Reader may be open at a time.
If you need a separate timeout on the Reader call and the Read itself,
use time.AfterFunc to cancel the context passed in.
See
https://github.com/nhooyr/websocket/issues/87#issue-451703332
Most users should not need this.
func (*Conn)
SetReadLimit
¶
func (c *
Conn
) SetReadLimit(n
int64
)
SetReadLimit sets the max number of bytes to read for a single message.
It applies to the Reader and Read methods.
By default, the connection has a message read limit of 32768 bytes.
When the limit is hit, reads return an error wrapping ErrMessageTooBig and
the connection is closed with StatusMessageTooBig.
Set to -1 to disable.
func (*Conn)
Subprotocol
¶
func (c *
Conn
) Subprotocol()
string
Subprotocol returns the negotiated subprotocol.
An empty string means the default protocol.
func (*Conn)
Write
¶
func (c *
Conn
) Write(ctx
context
.
Context
, typ
MessageType
, p []
byte
)
error
Write writes a message to the connection.
See the Writer method if you want to stream a message.
If compression is disabled or the compression threshold is not met, then it
will write the message in a single frame.
func (*Conn)
Writer
¶
func (c *
Conn
) Writer(ctx
context
.
Context
, typ
MessageType
) (
io
.
WriteCloser
,
error
)
Writer returns a writer bounded by the context that will write
a WebSocket message of type dataType to the connection.
You must close the writer once you have written the entire message.
Only one writer can be open at a time, multiple calls will block until the previous writer
is closed.
type
DialOptions
¶
type DialOptions struct {
// HTTPClient is used for the connection.
// Its Transport must return writable bodies for WebSocket handshakes.
// http.Transport does beginning with Go 1.12.
HTTPClient *
http
.
Client
// HTTPHeader specifies the HTTP headers included in the handshake request.
HTTPHeader
http
.
Header
// Host optionally overrides the Host HTTP header to send. If empty, the value
// of URL.Host will be used.
Host
string
// Subprotocols lists the WebSocket subprotocols to negotiate with the server.
Subprotocols []
string
// CompressionMode controls the compression mode.
// Defaults to CompressionDisabled.
//
// See docs on CompressionMode for details.
CompressionMode
CompressionMode
// CompressionThreshold controls the minimum size of a message before compression is applied.
//
// Defaults to 512 bytes for CompressionNoContextTakeover and 128 bytes
// for CompressionContextTakeover.
CompressionThreshold
int
// OnPingReceived is an optional callback invoked synchronously when a ping frame is received.
//
// The payload contains the application data of the ping frame.
// If the callback returns false, the subsequent pong frame will not be sent.
// To avoid blocking, any expensive processing should be performed asynchronously using a goroutine.
OnPingReceived func(ctx
context
.
Context
, payload []
byte
)
bool
// OnPongReceived is an optional callback invoked synchronously when a pong frame is received.
//
// The payload contains the application data of the pong frame.
// To avoid blocking, any expensive processing should be performed asynchronously using a goroutine.
//
// Unlike OnPingReceived, this callback does not return a value because a pong frame
// is a response to a ping and does not trigger any further frame transmission.
OnPongReceived func(ctx
context
.
Context
, payload []
byte
)
}
DialOptions represents Dial's options.
type
MessageType
¶
type MessageType
int
MessageType represents the type of a WebSocket message.
See
https://tools.ietf.org/html/rfc6455#section-5.6
const (
// MessageText is for UTF-8 encoded text messages like JSON.
MessageText
MessageType
=
iota
+ 1
// MessageBinary is for binary messages like protobufs.
MessageBinary
)
MessageType constants.
func (MessageType)
String
¶
func (i
MessageType
) String()
string
type
StatusCode
¶
type StatusCode
int
StatusCode represents a WebSocket status code.
https://tools.ietf.org/html/rfc6455#section-7.4
const (
StatusNormalClosure
StatusCode
= 1000
StatusGoingAway
StatusCode
= 1001
StatusProtocolError
StatusCode
= 1002
StatusUnsupportedData
StatusCode
= 1003
// StatusNoStatusRcvd cannot be sent in a close message.
// It is reserved for when a close message is received without
// a status code.
StatusNoStatusRcvd
StatusCode
= 1005
// StatusAbnormalClosure is exported for use only with Wasm.
// In non Wasm Go, the returned error will indicate whether the
// connection was closed abnormally.
StatusAbnormalClosure
StatusCode
= 1006
StatusInvalidFramePayloadData
StatusCode
= 1007
StatusPolicyViolation
StatusCode
= 1008
StatusMessageTooBig
StatusCode
= 1009
StatusMandatoryExtension
StatusCode
= 1010
StatusInternalError
StatusCode
= 1011
StatusServiceRestart
StatusCode
= 1012
StatusTryAgainLater
StatusCode
= 1013
StatusBadGateway
StatusCode
= 1014
// StatusTLSHandshake is only exported for use with Wasm.
// In non Wasm Go, the returned error will indicate whether there was
// a TLS handshake failure.
StatusTLSHandshake
StatusCode
= 1015
)
https://www.iana.org/assignments/websocket/websocket.xhtml#close-code-number
These are only the status codes defined by the protocol.
You can define custom codes in the 3000-4999 range.
The 3000-3999 range is reserved for use by libraries, frameworks and applications.
The 4000-4999 range is reserved for private use.
func
CloseStatus
¶
func CloseStatus(err
error
)
StatusCode
CloseStatus is a convenience wrapper around Go 1.13's errors.As to grab
the status code from a CloseError.
-1 will be returned if the passed error is nil or not a CloseError.
Example
¶
package main

import (
"context"
"log"
"time"

"github.com/coder/websocket"
)

func main() {
// Dials a server and then expects to be disconnected with status code
// websocket.StatusNormalClosure.

ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
defer cancel()

c, _, err := websocket.Dial(ctx, "ws://localhost:8080", nil)
if err != nil {
log.Fatal(err)
}
defer c.CloseNow()

_, _, err = c.Reader(ctx)
if websocket.CloseStatus(err) != websocket.StatusNormalClosure {
log.Fatalf("expected to be disconnected with StatusNormalClosure but got: %v", err)
}
}
Share
Format
Run
func (StatusCode)
String
¶
func (i
StatusCode
) String()
string