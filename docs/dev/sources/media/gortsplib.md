# gortsplib (RTSP)

> Auto-fetched from [https://pkg.go.dev/github.com/bluenviron/gortsplib/v4](https://pkg.go.dev/github.com/bluenviron/gortsplib/v4)
> Last Updated: 2026-01-29T20:14:36.429979+00:00

---

Overview
¶
Package gortsplib is a RTSP library for the Go programming language.
Examples are available at
https://github.com/bluenviron/gortsplib/tree/main/examples
Index
¶
type Client
func (c *Client) Announce(u *base.URL, desc *description.Session) (*base.Response, error)
func (c *Client) Close()
func (c *Client) Describe(u *base.URL) (*description.Session, *base.Response, error)
func (c *Client) OnPacketRTCP(medi *description.Media, cb OnPacketRTCPFunc)
func (c *Client) OnPacketRTCPAny(cb OnPacketRTCPAnyFunc)
func (c *Client) OnPacketRTP(medi *description.Media, forma format.Format, cb OnPacketRTPFunc)
func (c *Client) OnPacketRTPAny(cb OnPacketRTPAnyFunc)
func (c *Client) Options(u *base.URL) (*base.Response, error)
func (c *Client) PacketNTP(medi *description.Media, pkt *rtp.Packet) (time.Time, bool)
func (c *Client) PacketPTS(medi *description.Media, pkt *rtp.Packet) (time.Duration, bool)
deprecated
func (c *Client) PacketPTS2(medi *description.Media, pkt *rtp.Packet) (int64, bool)
func (c *Client) Pause() (*base.Response, error)
func (c *Client) Play(ra *headers.Range) (*base.Response, error)
func (c *Client) Record() (*base.Response, error)
func (c *Client) Seek(ra *headers.Range) (*base.Response, error)
deprecated
func (c *Client) Setup(baseURL *base.URL, media *description.Media, rtpPort int, rtcpPort int) (*base.Response, error)
func (c *Client) SetupAll(baseURL *base.URL, medias []*description.Media) error
func (c *Client) Start(scheme string, host string) error
deprecated
func (c *Client) Start2() error
func (c *Client) StartRecording(address string, desc *description.Session) error
func (c *Client) Stats() *ClientStats
func (c *Client) Wait() error
func (c *Client) WritePacketRTCP(medi *description.Media, pkt rtcp.Packet) error
func (c *Client) WritePacketRTP(medi *description.Media, pkt *rtp.Packet) error
func (c *Client) WritePacketRTPWithNTP(medi *description.Media, pkt *rtp.Packet, ntp time.Time) error
type ClientOnDecodeErrorFunc
type ClientOnPacketLostFunc
deprecated
type ClientOnPacketsLostFunc
type ClientOnRequestFunc
type ClientOnResponseFunc
type ClientOnTransportSwitchFunc
type ClientStats
type OnPacketRTCPAnyFunc
type OnPacketRTCPFunc
type OnPacketRTPAnyFunc
type OnPacketRTPFunc
type Server
func (s *Server) Close()
func (s *Server) Start() error
func (s *Server) StartAndWait() error
func (s *Server) Wait() error
type ServerConn
func (sc *ServerConn) BytesReceived() uint64
deprecated
func (sc *ServerConn) BytesSent() uint64
deprecated
func (sc *ServerConn) Close()
func (sc *ServerConn) NetConn() net.Conn
func (sc *ServerConn) Session() *ServerSession
func (sc *ServerConn) SetUserData(v interface{})
func (sc *ServerConn) Stats() *StatsConn
func (sc *ServerConn) UserData() interface{}
func (sc *ServerConn) VerifyCredentials(req *base.Request, expectedUser string, expectedPass string) bool
type ServerHandler
type ServerHandlerOnAnnounce
type ServerHandlerOnAnnounceCtx
type ServerHandlerOnConnClose
type ServerHandlerOnConnCloseCtx
type ServerHandlerOnConnOpen
type ServerHandlerOnConnOpenCtx
type ServerHandlerOnDecodeError
type ServerHandlerOnDecodeErrorCtx
type ServerHandlerOnDescribe
type ServerHandlerOnDescribeCtx
type ServerHandlerOnGetParameter
type ServerHandlerOnGetParameterCtx
type ServerHandlerOnPacketLost
deprecated
type ServerHandlerOnPacketLostCtx
deprecated
type ServerHandlerOnPacketsLost
type ServerHandlerOnPacketsLostCtx
type ServerHandlerOnPause
type ServerHandlerOnPauseCtx
type ServerHandlerOnPlay
type ServerHandlerOnPlayCtx
type ServerHandlerOnRecord
type ServerHandlerOnRecordCtx
type ServerHandlerOnRequest
type ServerHandlerOnResponse
type ServerHandlerOnSessionClose
type ServerHandlerOnSessionCloseCtx
type ServerHandlerOnSessionOpen
type ServerHandlerOnSessionOpenCtx
type ServerHandlerOnSetParameter
type ServerHandlerOnSetParameterCtx
type ServerHandlerOnSetup
type ServerHandlerOnSetupCtx
type ServerHandlerOnStreamWriteError
type ServerHandlerOnStreamWriteErrorCtx
type ServerSession
func (ss *ServerSession) AnnouncedDescription() *description.Session
func (ss *ServerSession) BytesReceived() uint64
deprecated
func (ss *ServerSession) BytesSent() uint64
deprecated
func (ss *ServerSession) Close()
func (ss *ServerSession) OnPacketRTCP(medi *description.Media, cb OnPacketRTCPFunc)
func (ss *ServerSession) OnPacketRTCPAny(cb OnPacketRTCPAnyFunc)
func (ss *ServerSession) OnPacketRTP(medi *description.Media, forma format.Format, cb OnPacketRTPFunc)
func (ss *ServerSession) OnPacketRTPAny(cb OnPacketRTPAnyFunc)
func (ss *ServerSession) PacketNTP(medi *description.Media, pkt *rtp.Packet) (time.Time, bool)
func (ss *ServerSession) PacketPTS(medi *description.Media, pkt *rtp.Packet) (time.Duration, bool)
deprecated
func (ss *ServerSession) PacketPTS2(medi *description.Media, pkt *rtp.Packet) (int64, bool)
func (ss *ServerSession) SetUserData(v interface{})
func (ss *ServerSession) SetuppedMedias() []*description.Media
func (ss *ServerSession) SetuppedPath() string
func (ss *ServerSession) SetuppedQuery() string
func (ss *ServerSession) SetuppedSecure() bool
func (ss *ServerSession) SetuppedStream() *ServerStream
func (ss *ServerSession) SetuppedTransport() *Transport
func (ss *ServerSession) State() ServerSessionState
func (ss *ServerSession) Stats() *StatsSession
func (ss *ServerSession) UserData() interface{}
func (ss *ServerSession) WritePacketRTCP(medi *description.Media, pkt rtcp.Packet) error
func (ss *ServerSession) WritePacketRTP(medi *description.Media, pkt *rtp.Packet) error
type ServerSessionState
func (s ServerSessionState) String() string
type ServerStream
func NewServerStream(s *Server, desc *description.Session) *ServerStream
deprecated
func (st *ServerStream) BytesSent() uint64
deprecated
func (st *ServerStream) Close()
func (st *ServerStream) Description() *description.Session
deprecated
func (st *ServerStream) Initialize() error
func (st *ServerStream) Stats() *ServerStreamStats
func (st *ServerStream) WritePacketRTCP(medi *description.Media, pkt rtcp.Packet) error
func (st *ServerStream) WritePacketRTP(medi *description.Media, pkt *rtp.Packet) error
func (st *ServerStream) WritePacketRTPWithNTP(medi *description.Media, pkt *rtp.Packet, ntp time.Time) error
type ServerStreamStats
type ServerStreamStatsFormat
type ServerStreamStatsMedia
type StatsConn
type StatsSession
type StatsSessionFormat
type StatsSessionMedia
type Transport
func (t Transport) String() string
Constants
¶
This section is empty.
Variables
¶
This section is empty.
Functions
¶
This section is empty.
Types
¶
type
Client
¶
type Client struct {
//
// Target
//
// Scheme. Either "rtsp" or "rtsps".
Scheme
string
// Host and port.
Host
string
//
// RTSP parameters (all optional)
//
// timeout of read operations.
// It defaults to 10 seconds.
ReadTimeout
time
.
Duration
// timeout of write operations.
// It defaults to 10 seconds.
WriteTimeout
time
.
Duration
// a TLS configuration to connect to TLS (RTSPS) servers.
// It defaults to nil.
TLSConfig *
tls
.
Config
// enable communication with servers which don't provide UDP server ports
// or use different server ports than the announced ones.
// This can be a security issue.
// It defaults to false.
AnyPortEnable
bool
// transport protocol (UDP, Multicast or TCP).
// If nil, it is chosen automatically (first UDP, then, if it fails, TCP).
// It defaults to nil.
Transport *
Transport
// If the client is reading with UDP, it must receive
// at least a packet within this timeout, otherwise it switches to TCP.
// It defaults to 3 seconds.
InitialUDPReadTimeout
time
.
Duration
// Size of the UDP read buffer.
// This can be increased to reduce packet losses.
// It defaults to the operating system default value.
UDPReadBufferSize
int
// Size of the queue of outgoing packets.
// It defaults to 256.
WriteQueueSize
int
// maximum size of outgoing RTP / RTCP packets.
// This must be less than the UDP MTU (1472 bytes).
// It defaults to 1472.
MaxPacketSize
int
// user agent header.
// It defaults to "gortsplib"
UserAgent
string
// disable automatic RTCP sender reports.
DisableRTCPSenderReports
bool
// explicitly request back channels to the server.
RequestBackChannels
bool
// pointer to a variable that stores received bytes.
//
// Deprecated: use Client.Stats()
BytesReceived *
uint64
// pointer to a variable that stores sent bytes.
//
// Deprecated: use Client.Stats()
BytesSent *
uint64
//
// system functions (all optional)
//
// function used to initialize the TCP client.
// It defaults to (&net.Dialer{}).DialContext.
DialContext func(ctx
context
.
Context
, network, address
string
) (
net
.
Conn
,
error
)
// function used to initialize UDP listeners.
// It defaults to net.ListenPacket.
ListenPacket func(network, address
string
) (
net
.
PacketConn
,
error
)
//
// callbacks (all optional)
//
// called when sending a request to the server.
OnRequest
ClientOnRequestFunc
// called when receiving a response from the server.
OnResponse
ClientOnResponseFunc
// called when receiving a request from the server.
OnServerRequest
ClientOnRequestFunc
// called when sending a response to the server.
OnServerResponse
ClientOnResponseFunc
// called when the transport protocol changes.
OnTransportSwitch
ClientOnTransportSwitchFunc
// called when the client detects lost packets.
//
// Deprecated: replaced by OnPacketsLost
OnPacketLost
ClientOnPacketLostFunc
// called when the client detects lost packets.
OnPacketsLost
ClientOnPacketsLostFunc
// called when a non-fatal decode error occurs.
OnDecodeError
ClientOnDecodeErrorFunc
// contains filtered or unexported fields
}
Client is a RTSP client.
func (*Client)
Announce
¶
func (c *
Client
) Announce(u *
base
.
URL
, desc *
description
.
Session
) (*
base
.
Response
,
error
)
Announce sends an ANNOUNCE request.
func (*Client)
Close
¶
func (c *
Client
) Close()
Close closes all client resources and waits for them to exit.
func (*Client)
Describe
¶
func (c *
Client
) Describe(u *
base
.
URL
) (*
description
.
Session
, *
base
.
Response
,
error
)
Describe sends a DESCRIBE request.
func (*Client)
OnPacketRTCP
¶
func (c *
Client
) OnPacketRTCP(medi *
description
.
Media
, cb
OnPacketRTCPFunc
)
OnPacketRTCP sets a callback that is called when a RTCP packet is read.
func (*Client)
OnPacketRTCPAny
¶
func (c *
Client
) OnPacketRTCPAny(cb
OnPacketRTCPAnyFunc
)
OnPacketRTCPAny sets a callback that is called when a RTCP packet is read from any setupped media.
func (*Client)
OnPacketRTP
¶
func (c *
Client
) OnPacketRTP(medi *
description
.
Media
, forma
format
.
Format
, cb
OnPacketRTPFunc
)
OnPacketRTP sets a callback that is called when a RTP packet is read.
func (*Client)
OnPacketRTPAny
¶
func (c *
Client
) OnPacketRTPAny(cb
OnPacketRTPAnyFunc
)
OnPacketRTPAny sets a callback that is called when a RTP packet is read from any setupped media.
func (*Client)
Options
¶
func (c *
Client
) Options(u *
base
.
URL
) (*
base
.
Response
,
error
)
Options sends an OPTIONS request.
func (*Client)
PacketNTP
¶
func (c *
Client
) PacketNTP(medi *
description
.
Media
, pkt *
rtp
.
Packet
) (
time
.
Time
,
bool
)
PacketNTP returns the NTP (absolute timestamp) of an incoming RTP packet.
The NTP is computed from RTCP sender reports.
func (*Client)
PacketPTS
deprecated
func (c *
Client
) PacketPTS(medi *
description
.
Media
, pkt *
rtp
.
Packet
) (
time
.
Duration
,
bool
)
PacketPTS returns the PTS (presentation timestamp) of an incoming RTP packet.
It is computed by decoding the packet timestamp and sychronizing it with other tracks.
Deprecated: replaced by PacketPTS2.
func (*Client)
PacketPTS2
¶
added in
v4.11.0
func (c *
Client
) PacketPTS2(medi *
description
.
Media
, pkt *
rtp
.
Packet
) (
int64
,
bool
)
PacketPTS2 returns the PTS (presentation timestamp) of an incoming RTP packet.
It is computed by decoding the packet timestamp and sychronizing it with other tracks.
func (*Client)
Pause
¶
func (c *
Client
) Pause() (*
base
.
Response
,
error
)
Pause sends a PAUSE request.
This can be called only after Play() or Record().
func (*Client)
Play
¶
func (c *
Client
) Play(ra *
headers
.
Range
) (*
base
.
Response
,
error
)
Play sends a PLAY request.
This can be called only after Setup().
func (*Client)
Record
¶
func (c *
Client
) Record() (*
base
.
Response
,
error
)
Record sends a RECORD request.
This can be called only after Announce() and Setup().
func (*Client)
Seek
deprecated
func (c *
Client
) Seek(ra *
headers
.
Range
) (*
base
.
Response
,
error
)
Seek asks the server to re-start the stream from a specific timestamp.
Deprecated: will be removed in next version. Equivalent to using Pause() followed by Play().
func (*Client)
Setup
¶
func (c *
Client
) Setup(
baseURL *
base
.
URL
,
media *
description
.
Media
,
rtpPort
int
,
rtcpPort
int
,
) (*
base
.
Response
,
error
)
Setup sends a SETUP request.
rtpPort and rtcpPort are used only if transport is UDP.
if rtpPort and rtcpPort are zero, they are chosen automatically.
func (*Client)
SetupAll
¶
func (c *
Client
) SetupAll(baseURL *
base
.
URL
, medias []*
description
.
Media
)
error
SetupAll setups all the given medias.
func (*Client)
Start
deprecated
func (c *
Client
) Start(scheme
string
, host
string
)
error
Start initializes the connection to a server.
Deprecated: replaced by Start2.
func (*Client)
Start2
¶
added in
v4.15.0
func (c *
Client
) Start2()
error
Start2 initializes the connection to a server.
func (*Client)
StartRecording
¶
func (c *
Client
) StartRecording(address
string
, desc *
description
.
Session
)
error
StartRecording connects to the address and starts publishing given media.
func (*Client)
Stats
¶
added in
v4.12.0
func (c *
Client
) Stats() *
ClientStats
Stats returns client statistics.
func (*Client)
Wait
¶
func (c *
Client
) Wait()
error
Wait waits until all client resources are closed.
This can happen when a fatal error occurs or when Close() is called.
func (*Client)
WritePacketRTCP
¶
func (c *
Client
) WritePacketRTCP(medi *
description
.
Media
, pkt
rtcp
.
Packet
)
error
WritePacketRTCP writes a RTCP packet to the server.
func (*Client)
WritePacketRTP
¶
func (c *
Client
) WritePacketRTP(medi *
description
.
Media
, pkt *
rtp
.
Packet
)
error
WritePacketRTP writes a RTP packet to the server.
func (*Client)
WritePacketRTPWithNTP
¶
func (c *
Client
) WritePacketRTPWithNTP(medi *
description
.
Media
, pkt *
rtp
.
Packet
, ntp
time
.
Time
)
error
WritePacketRTPWithNTP writes a RTP packet to the server.
ntp is the absolute timestamp of the packet, and is sent with periodic RTCP sender reports.
type
ClientOnDecodeErrorFunc
¶
type ClientOnDecodeErrorFunc func(err
error
)
ClientOnDecodeErrorFunc is the prototype of Client.OnDecodeError.
type
ClientOnPacketLostFunc
deprecated
type ClientOnPacketLostFunc func(err
error
)
ClientOnPacketLostFunc is the prototype of Client.OnPacketLost.
Deprecated: replaced by ClientOnPacketsLostFunc
type
ClientOnPacketsLostFunc
¶
added in
v4.13.0
type ClientOnPacketsLostFunc func(lost
uint64
)
ClientOnPacketsLostFunc is the prototype of Client.OnPacketsLost.
type
ClientOnRequestFunc
¶
type ClientOnRequestFunc func(*
base
.
Request
)
ClientOnRequestFunc is the prototype of Client.OnRequest.
type
ClientOnResponseFunc
¶
type ClientOnResponseFunc func(*
base
.
Response
)
ClientOnResponseFunc is the prototype of Client.OnResponse.
type
ClientOnTransportSwitchFunc
¶
type ClientOnTransportSwitchFunc func(err
error
)
ClientOnTransportSwitchFunc is the prototype of Client.OnTransportSwitch.
type
ClientStats
¶
added in
v4.12.0
type ClientStats struct {
Conn
StatsConn
Session
StatsSession
}
ClientStats are client statistics
type
OnPacketRTCPAnyFunc
¶
type OnPacketRTCPAnyFunc func(*
description
.
Media
,
rtcp
.
Packet
)
OnPacketRTCPAnyFunc is the prototype of the callback passed to OnPacketRTCPAny().
type
OnPacketRTCPFunc
¶
type OnPacketRTCPFunc func(
rtcp
.
Packet
)
OnPacketRTCPFunc is the prototype of the callback passed to OnPacketRTCP().
type
OnPacketRTPAnyFunc
¶
type OnPacketRTPAnyFunc func(*
description
.
Media
,
format
.
Format
, *
rtp
.
Packet
)
OnPacketRTPAnyFunc is the prototype of the callback passed to OnPacketRTP(Any).
type
OnPacketRTPFunc
¶
type OnPacketRTPFunc func(*
rtp
.
Packet
)
OnPacketRTPFunc is the prototype of the callback passed to OnPacketRTP().
type
Server
¶
type Server struct {
//
// RTSP parameters (all optional except RTSPAddress)
//
// the RTSP address of the server, to accept connections and send and receive
// packets with the TCP transport.
RTSPAddress
string
// a port to send and receive RTP packets with the UDP transport.
// If UDPRTPAddress and UDPRTCPAddress are filled, the server can support the UDP transport.
UDPRTPAddress
string
// a port to send and receive RTCP packets with the UDP transport.
// If UDPRTPAddress and UDPRTCPAddress are filled, the server can support the UDP transport.
UDPRTCPAddress
string
// a range of multicast IPs to use with the UDP-multicast transport.
// If MulticastIPRange, MulticastRTPPort, MulticastRTCPPort are filled, the server
// can support the UDP-multicast transport.
MulticastIPRange
string
// a port to send RTP packets with the UDP-multicast transport.
// If MulticastIPRange, MulticastRTPPort, MulticastRTCPPort are filled, the server
// can support the UDP-multicast transport.
MulticastRTPPort
int
// a port to send RTCP packets with the UDP-multicast transport.
// If MulticastIPRange, MulticastRTPPort, MulticastRTCPPort are filled, the server
// can support the UDP-multicast transport.
MulticastRTCPPort
int
// timeout of read operations.
// It defaults to 10 seconds
ReadTimeout
time
.
Duration
// timeout of write operations.
// It defaults to 10 seconds
WriteTimeout
time
.
Duration
// a TLS configuration to accept TLS (RTSPS) connections.
TLSConfig *
tls
.
Config
// Size of the UDP read buffer.
// This can be increased to reduce packet losses.
// It defaults to the operating system default value.
UDPReadBufferSize
int
// Size of the queue of outgoing packets.
// It defaults to 256.
WriteQueueSize
int
// maximum size of outgoing RTP / RTCP packets.
// This must be less than the UDP MTU (1472 bytes).
// It defaults to 1472.
MaxPacketSize
int
// disable automatic RTCP sender reports.
DisableRTCPSenderReports
bool
// authentication methods.
// It defaults to plain and digest+MD5.
AuthMethods []
auth
.
VerifyMethod
//
// handler (optional)
//
// an handler to handle server events.
// It may implement one or more of the ServerHandler* interfaces.
Handler
ServerHandler
//
// system functions (all optional)
//
// function used to initialize the TCP listener.
// It defaults to net.Listen.
Listen func(network
string
, address
string
) (
net
.
Listener
,
error
)
// function used to initialize UDP listeners.
// It defaults to net.ListenPacket.
ListenPacket func(network, address
string
) (
net
.
PacketConn
,
error
)
// contains filtered or unexported fields
}
Server is a RTSP server.
func (*Server)
Close
¶
func (s *
Server
) Close()
Close closes all the server resources and waits for them to exit.
func (*Server)
Start
¶
func (s *
Server
) Start()
error
Start starts the server.
func (*Server)
StartAndWait
¶
func (s *
Server
) StartAndWait()
error
StartAndWait starts the server and waits until a fatal error.
func (*Server)
Wait
¶
func (s *
Server
) Wait()
error
Wait waits until all server resources are closed.
This can happen when a fatal error occurs or when Close() is called.
type
ServerConn
¶
type ServerConn struct {
// contains filtered or unexported fields
}
ServerConn is a server-side RTSP connection.
func (*ServerConn)
BytesReceived
deprecated
func (sc *
ServerConn
) BytesReceived()
uint64
BytesReceived returns the number of read bytes.
Deprecated: replaced by Stats()
func (*ServerConn)
BytesSent
deprecated
func (sc *
ServerConn
) BytesSent()
uint64
BytesSent returns the number of written bytes.
Deprecated: replaced by Stats()
func (*ServerConn)
Close
¶
func (sc *
ServerConn
) Close()
Close closes the ServerConn.
func (*ServerConn)
NetConn
¶
func (sc *
ServerConn
) NetConn()
net
.
Conn
NetConn returns the underlying net.Conn.
func (*ServerConn)
Session
¶
added in
v4.12.1
func (sc *
ServerConn
) Session() *
ServerSession
Session returns the associated session.
func (*ServerConn)
SetUserData
¶
func (sc *
ServerConn
) SetUserData(v interface{})
SetUserData sets some user data associated with the connection.
func (*ServerConn)
Stats
¶
added in
v4.12.0
func (sc *
ServerConn
) Stats() *
StatsConn
Stats returns connection statistics.
func (*ServerConn)
UserData
¶
func (sc *
ServerConn
) UserData() interface{}
UserData returns some user data associated with the connection.
func (*ServerConn)
VerifyCredentials
¶
added in
v4.13.0
func (sc *
ServerConn
) VerifyCredentials(
req *
base
.
Request
,
expectedUser
string
,
expectedPass
string
,
)
bool
VerifyCredentials verifies credentials provided by the user.
type
ServerHandler
¶
type ServerHandler interface{}
ServerHandler is the interface implemented by all the server handlers.
type
ServerHandlerOnAnnounce
¶
type ServerHandlerOnAnnounce interface {
// called when receiving an ANNOUNCE request.
OnAnnounce(*
ServerHandlerOnAnnounceCtx
) (*
base
.
Response
,
error
)
}
ServerHandlerOnAnnounce can be implemented by a ServerHandler.
type
ServerHandlerOnAnnounceCtx
¶
type ServerHandlerOnAnnounceCtx struct {
Session     *
ServerSession
Conn        *
ServerConn
Request     *
base
.
Request
Path
string
Query
string
Description *
description
.
Session
}
ServerHandlerOnAnnounceCtx is the context of OnAnnounce.
type
ServerHandlerOnConnClose
¶
type ServerHandlerOnConnClose interface {
// called when a connection is closed.
OnConnClose(*
ServerHandlerOnConnCloseCtx
)
}
ServerHandlerOnConnClose can be implemented by a ServerHandler.
type
ServerHandlerOnConnCloseCtx
¶
type ServerHandlerOnConnCloseCtx struct {
Conn  *
ServerConn
Error
error
}
ServerHandlerOnConnCloseCtx is the context of OnConnClose.
type
ServerHandlerOnConnOpen
¶
type ServerHandlerOnConnOpen interface {
// called when a connection is opened.
OnConnOpen(*
ServerHandlerOnConnOpenCtx
)
}
ServerHandlerOnConnOpen can be implemented by a ServerHandler.
type
ServerHandlerOnConnOpenCtx
¶
type ServerHandlerOnConnOpenCtx struct {
Conn *
ServerConn
}
ServerHandlerOnConnOpenCtx is the context of OnConnOpen.
type
ServerHandlerOnDecodeError
¶
type ServerHandlerOnDecodeError interface {
// called when a non-fatal decode error occurs.
OnDecodeError(*
ServerHandlerOnDecodeErrorCtx
)
}
ServerHandlerOnDecodeError can be implemented by a ServerHandler.
type
ServerHandlerOnDecodeErrorCtx
¶
type ServerHandlerOnDecodeErrorCtx struct {
Session *
ServerSession
Error
error
}
ServerHandlerOnDecodeErrorCtx is the context of OnDecodeError.
type
ServerHandlerOnDescribe
¶
type ServerHandlerOnDescribe interface {
// called when receiving a DESCRIBE request.
OnDescribe(*
ServerHandlerOnDescribeCtx
) (*
base
.
Response
, *
ServerStream
,
error
)
}
ServerHandlerOnDescribe can be implemented by a ServerHandler.
type
ServerHandlerOnDescribeCtx
¶
type ServerHandlerOnDescribeCtx struct {
Conn    *
ServerConn
Request *
base
.
Request
Path
string
Query
string
}
ServerHandlerOnDescribeCtx is the context of OnDescribe.
type
ServerHandlerOnGetParameter
¶
type ServerHandlerOnGetParameter interface {
// called when receiving a GET_PARAMETER request.
OnGetParameter(*
ServerHandlerOnGetParameterCtx
) (*
base
.
Response
,
error
)
}
ServerHandlerOnGetParameter can be implemented by a ServerHandler.
type
ServerHandlerOnGetParameterCtx
¶
type ServerHandlerOnGetParameterCtx struct {
Session *
ServerSession
Conn    *
ServerConn
Request *
base
.
Request
Path
string
Query
string
}
ServerHandlerOnGetParameterCtx is the context of OnGetParameter.
type
ServerHandlerOnPacketLost
deprecated
type ServerHandlerOnPacketLost interface {
// called when the server detects lost packets.
OnPacketLost(*
ServerHandlerOnPacketLostCtx
)
}
ServerHandlerOnPacketLost can be implemented by a ServerHandler.
Deprecated: replaced by ServerHandlerOnPacketsLost
type
ServerHandlerOnPacketLostCtx
deprecated
type ServerHandlerOnPacketLostCtx struct {
Session *
ServerSession
Error
error
}
ServerHandlerOnPacketLostCtx is the context of OnPacketLost.
Deprecated: replaced by ServerHandlerOnPacketsLostCtx
type
ServerHandlerOnPacketsLost
¶
added in
v4.13.0
type ServerHandlerOnPacketsLost interface {
// called when the server detects lost packets.
OnPacketsLost(*
ServerHandlerOnPacketsLostCtx
)
}
ServerHandlerOnPacketsLost can be implemented by a ServerHandler.
type
ServerHandlerOnPacketsLostCtx
¶
added in
v4.13.0
type ServerHandlerOnPacketsLostCtx struct {
Session *
ServerSession
Lost
uint64
}
ServerHandlerOnPacketsLostCtx is the context of OnPacketsLost.
type
ServerHandlerOnPause
¶
type ServerHandlerOnPause interface {
// called when receiving a PAUSE request.
OnPause(*
ServerHandlerOnPauseCtx
) (*
base
.
Response
,
error
)
}
ServerHandlerOnPause can be implemented by a ServerHandler.
type
ServerHandlerOnPauseCtx
¶
type ServerHandlerOnPauseCtx struct {
Session *
ServerSession
Conn    *
ServerConn
Request *
base
.
Request
Path
string
Query
string
}
ServerHandlerOnPauseCtx is the context of OnPause.
type
ServerHandlerOnPlay
¶
type ServerHandlerOnPlay interface {
// called when receiving a PLAY request.
OnPlay(*
ServerHandlerOnPlayCtx
) (*
base
.
Response
,
error
)
}
ServerHandlerOnPlay can be implemented by a ServerHandler.
type
ServerHandlerOnPlayCtx
¶
type ServerHandlerOnPlayCtx struct {
Session *
ServerSession
Conn    *
ServerConn
Request *
base
.
Request
Path
string
Query
string
}
ServerHandlerOnPlayCtx is the context of OnPlay.
type
ServerHandlerOnRecord
¶
type ServerHandlerOnRecord interface {
// called when receiving a RECORD request.
OnRecord(*
ServerHandlerOnRecordCtx
) (*
base
.
Response
,
error
)
}
ServerHandlerOnRecord can be implemented by a ServerHandler.
type
ServerHandlerOnRecordCtx
¶
type ServerHandlerOnRecordCtx struct {
Session *
ServerSession
Conn    *
ServerConn
Request *
base
.
Request
Path
string
Query
string
}
ServerHandlerOnRecordCtx is the context of OnRecord.
type
ServerHandlerOnRequest
¶
type ServerHandlerOnRequest interface {
// called when receiving a request from a connection.
OnRequest(*
ServerConn
, *
base
.
Request
)
}
ServerHandlerOnRequest can be implemented by a ServerHandler.
type
ServerHandlerOnResponse
¶
type ServerHandlerOnResponse interface {
// called when sending a response to a connection.
OnResponse(*
ServerConn
, *
base
.
Response
)
}
ServerHandlerOnResponse can be implemented by a ServerHandler.
type
ServerHandlerOnSessionClose
¶
type ServerHandlerOnSessionClose interface {
// called when a session is closed.
OnSessionClose(*
ServerHandlerOnSessionCloseCtx
)
}
ServerHandlerOnSessionClose can be implemented by a ServerHandler.
type
ServerHandlerOnSessionCloseCtx
¶
type ServerHandlerOnSessionCloseCtx struct {
Session *
ServerSession
Error
error
}
ServerHandlerOnSessionCloseCtx is the context of ServerHandlerOnSessionClose.
type
ServerHandlerOnSessionOpen
¶
type ServerHandlerOnSessionOpen interface {
// called when a session is opened.
OnSessionOpen(*
ServerHandlerOnSessionOpenCtx
)
}
ServerHandlerOnSessionOpen can be implemented by a ServerHandler.
type
ServerHandlerOnSessionOpenCtx
¶
type ServerHandlerOnSessionOpenCtx struct {
Session *
ServerSession
Conn    *
ServerConn
}
ServerHandlerOnSessionOpenCtx is the context OnSessionOpen.
type
ServerHandlerOnSetParameter
¶
type ServerHandlerOnSetParameter interface {
// called when receiving a SET_PARAMETER request.
OnSetParameter(*
ServerHandlerOnSetParameterCtx
) (*
base
.
Response
,
error
)
}
ServerHandlerOnSetParameter can be implemented by a ServerHandler.
type
ServerHandlerOnSetParameterCtx
¶
type ServerHandlerOnSetParameterCtx struct {
Session *
ServerSession
Conn    *
ServerConn
Request *
base
.
Request
Path
string
Query
string
}
ServerHandlerOnSetParameterCtx is the context of OnSetParameter.
type
ServerHandlerOnSetup
¶
type ServerHandlerOnSetup interface {
// called when receiving a SETUP request.
// must return a Response and a stream.
// the stream is needed to
// - add the session the the stream's readers
// - send the stream SSRC to the session
OnSetup(*
ServerHandlerOnSetupCtx
) (*
base
.
Response
, *
ServerStream
,
error
)
}
ServerHandlerOnSetup can be implemented by a ServerHandler.
type
ServerHandlerOnSetupCtx
¶
type ServerHandlerOnSetupCtx struct {
Session   *
ServerSession
Conn      *
ServerConn
Request   *
base
.
Request
Path
string
Query
string
Transport
Transport
}
ServerHandlerOnSetupCtx is the context of OnSetup.
type
ServerHandlerOnStreamWriteError
¶
type ServerHandlerOnStreamWriteError interface {
// called when a ServerStream is unable to write packets to a session.
OnStreamWriteError(*
ServerHandlerOnStreamWriteErrorCtx
)
}
ServerHandlerOnStreamWriteError can be implemented by a ServerHandler.
type
ServerHandlerOnStreamWriteErrorCtx
¶
type ServerHandlerOnStreamWriteErrorCtx struct {
Session *
ServerSession
Error
error
}
ServerHandlerOnStreamWriteErrorCtx is the context of OnStreamWriteError.
type
ServerSession
¶
type ServerSession struct {
// contains filtered or unexported fields
}
ServerSession is a server-side RTSP session.
func (*ServerSession)
AnnouncedDescription
¶
func (ss *
ServerSession
) AnnouncedDescription() *
description
.
Session
AnnouncedDescription returns the announced stream description.
func (*ServerSession)
BytesReceived
deprecated
func (ss *
ServerSession
) BytesReceived()
uint64
BytesReceived returns the number of read bytes.
Deprecated: replaced by Stats()
func (*ServerSession)
BytesSent
deprecated
func (ss *
ServerSession
) BytesSent()
uint64
BytesSent returns the number of written bytes.
Deprecated: replaced by Stats()
func (*ServerSession)
Close
¶
func (ss *
ServerSession
) Close()
Close closes the ServerSession.
func (*ServerSession)
OnPacketRTCP
¶
func (ss *
ServerSession
) OnPacketRTCP(medi *
description
.
Media
, cb
OnPacketRTCPFunc
)
OnPacketRTCP sets a callback that is called when a RTCP packet is read.
func (*ServerSession)
OnPacketRTCPAny
¶
func (ss *
ServerSession
) OnPacketRTCPAny(cb
OnPacketRTCPAnyFunc
)
OnPacketRTCPAny sets a callback that is called when a RTCP packet is read from any setupped media.
func (*ServerSession)
OnPacketRTP
¶
func (ss *
ServerSession
) OnPacketRTP(medi *
description
.
Media
, forma
format
.
Format
, cb
OnPacketRTPFunc
)
OnPacketRTP sets a callback that is called when a RTP packet is read.
func (*ServerSession)
OnPacketRTPAny
¶
func (ss *
ServerSession
) OnPacketRTPAny(cb
OnPacketRTPAnyFunc
)
OnPacketRTPAny sets a callback that is called when a RTP packet is read from any setupped media.
func (*ServerSession)
PacketNTP
¶
func (ss *
ServerSession
) PacketNTP(medi *
description
.
Media
, pkt *
rtp
.
Packet
) (
time
.
Time
,
bool
)
PacketNTP returns the NTP (absolute timestamp) of an incoming RTP packet.
The NTP is computed from RTCP sender reports.
func (*ServerSession)
PacketPTS
deprecated
func (ss *
ServerSession
) PacketPTS(medi *
description
.
Media
, pkt *
rtp
.
Packet
) (
time
.
Duration
,
bool
)
PacketPTS returns the PTS (presentation timestamp) of an incoming RTP packet.
It is computed by decoding the packet timestamp and sychronizing it with other tracks.
Deprecated: replaced by PacketPTS2.
func (*ServerSession)
PacketPTS2
¶
added in
v4.11.0
func (ss *
ServerSession
) PacketPTS2(medi *
description
.
Media
, pkt *
rtp
.
Packet
) (
int64
,
bool
)
PacketPTS2 returns the PTS (presentation timestamp) of an incoming RTP packet.
It is computed by decoding the packet timestamp and sychronizing it with other tracks.
func (*ServerSession)
SetUserData
¶
func (ss *
ServerSession
) SetUserData(v interface{})
SetUserData sets some user data associated with the session.
func (*ServerSession)
SetuppedMedias
¶
func (ss *
ServerSession
) SetuppedMedias() []*
description
.
Media
SetuppedMedias returns the setupped medias.
func (*ServerSession)
SetuppedPath
¶
added in
v4.3.0
func (ss *
ServerSession
) SetuppedPath()
string
SetuppedPath returns the path sent during SETUP or ANNOUNCE.
func (*ServerSession)
SetuppedQuery
¶
added in
v4.3.0
func (ss *
ServerSession
) SetuppedQuery()
string
SetuppedQuery returns the query sent during SETUP or ANNOUNCE.
func (*ServerSession)
SetuppedSecure
¶
added in
v4.15.0
func (ss *
ServerSession
) SetuppedSecure()
bool
SetuppedSecure returns whether a secure profile is in use.
If this is false, it does not mean that the stream is not secure, since
there are some combinations that are secure nonetheless, like RTSPS+TCP+unsecure.
func (*ServerSession)
SetuppedStream
¶
added in
v4.3.0
func (ss *
ServerSession
) SetuppedStream() *
ServerStream
SetuppedStream returns the stream associated with the session.
func (*ServerSession)
SetuppedTransport
¶
func (ss *
ServerSession
) SetuppedTransport() *
Transport
SetuppedTransport returns the transport negotiated during SETUP.
func (*ServerSession)
State
¶
func (ss *
ServerSession
) State()
ServerSessionState
State returns the state of the session.
func (*ServerSession)
Stats
¶
added in
v4.12.0
func (ss *
ServerSession
) Stats() *
StatsSession
Stats returns server session statistics.
func (*ServerSession)
UserData
¶
func (ss *
ServerSession
) UserData() interface{}
UserData returns some user data associated with the session.
func (*ServerSession)
WritePacketRTCP
¶
func (ss *
ServerSession
) WritePacketRTCP(medi *
description
.
Media
, pkt
rtcp
.
Packet
)
error
WritePacketRTCP writes a RTCP packet to the session.
func (*ServerSession)
WritePacketRTP
¶
func (ss *
ServerSession
) WritePacketRTP(medi *
description
.
Media
, pkt *
rtp
.
Packet
)
error
WritePacketRTP writes a RTP packet to the session.
type
ServerSessionState
¶
type ServerSessionState
int
ServerSessionState is a state of a ServerSession.
const (
ServerSessionStateInitial
ServerSessionState
=
iota
ServerSessionStatePrePlay
ServerSessionStatePlay
ServerSessionStatePreRecord
ServerSessionStateRecord
)
states.
func (ServerSessionState)
String
¶
func (s
ServerSessionState
) String()
string
String implements fmt.Stringer.
type
ServerStream
¶
type ServerStream struct {
Server *
Server
Desc   *
description
.
Session
// contains filtered or unexported fields
}
ServerStream represents a data stream.
This is in charge of
- storing stream description and statistics
- distributing the stream to each reader
- allocating multicast listeners
func
NewServerStream
deprecated
func NewServerStream(s *
Server
, desc *
description
.
Session
) *
ServerStream
NewServerStream allocates a ServerStream.
Deprecated: replaced by ServerStream.Initialize().
func (*ServerStream)
BytesSent
deprecated
added in
v4.4.0
func (st *
ServerStream
) BytesSent()
uint64
BytesSent returns the number of written bytes.
Deprecated: replaced by Stats()
func (*ServerStream)
Close
¶
func (st *
ServerStream
) Close()
Close closes a ServerStream.
func (*ServerStream)
Description
deprecated
func (st *
ServerStream
) Description() *
description
.
Session
Description returns the description of the stream.
Deprecated: use ServerStream.Desc.
func (*ServerStream)
Initialize
¶
added in
v4.13.0
func (st *
ServerStream
) Initialize()
error
Initialize initializes a ServerStream.
func (*ServerStream)
Stats
¶
added in
v4.12.0
func (st *
ServerStream
) Stats() *
ServerStreamStats
Stats returns stream statistics.
func (*ServerStream)
WritePacketRTCP
¶
func (st *
ServerStream
) WritePacketRTCP(medi *
description
.
Media
, pkt
rtcp
.
Packet
)
error
WritePacketRTCP writes a RTCP packet to all the readers of the stream.
func (*ServerStream)
WritePacketRTP
¶
func (st *
ServerStream
) WritePacketRTP(medi *
description
.
Media
, pkt *
rtp
.
Packet
)
error
WritePacketRTP writes a RTP packet to all the readers of the stream.
func (*ServerStream)
WritePacketRTPWithNTP
¶
func (st *
ServerStream
) WritePacketRTPWithNTP(medi *
description
.
Media
, pkt *
rtp
.
Packet
, ntp
time
.
Time
)
error
WritePacketRTPWithNTP writes a RTP packet to all the readers of the stream.
ntp is the absolute timestamp of the packet, and is sent with periodic RTCP sender reports.
type
ServerStreamStats
¶
added in
v4.12.0
type ServerStreamStats struct {
// sent bytes
BytesSent
uint64
// number of sent RTP packets
RTPPacketsSent
uint64
// number of sent RTCP packets
RTCPPacketsSent
uint64
// media statistics
Medias map[*
description
.
Media
]
ServerStreamStatsMedia
}
ServerStreamStats are stream statistics.
type
ServerStreamStatsFormat
¶
added in
v4.12.0
type ServerStreamStatsFormat struct {
// number of sent RTP packets
RTPPacketsSent
uint64
// local SSRC
LocalSSRC
uint32
}
ServerStreamStatsFormat are stream format statistics.
type
ServerStreamStatsMedia
¶
added in
v4.12.0
type ServerStreamStatsMedia struct {
// sent bytes
BytesSent
uint64
// number of sent RTCP packets
RTCPPacketsSent
uint64
// format statistics
Formats map[
format
.
Format
]
ServerStreamStatsFormat
}
ServerStreamStatsMedia are stream media statistics.
type
StatsConn
¶
added in
v4.12.0
type StatsConn struct {
// received bytes
BytesReceived
uint64
// sent bytes
BytesSent
uint64
}
StatsConn are connection statistics.
type
StatsSession
¶
added in
v4.12.0
type StatsSession struct {
// received bytes
BytesReceived
uint64
// sent bytes
BytesSent
uint64
// number of RTP packets correctly received and processed
RTPPacketsReceived
uint64
// number of sent RTP packets
RTPPacketsSent
uint64
// number of lost RTP packets
RTPPacketsLost
uint64
// number of RTP packets that could not be processed
RTPPacketsInError
uint64
// mean jitter of received RTP packets
RTPPacketsJitter
float64
// number of RTCP packets correctly received and processed
RTCPPacketsReceived
uint64
// number of sent RTCP packets
RTCPPacketsSent
uint64
// number of RTCP packets that could not be processed
RTCPPacketsInError
uint64
// media statistics
Medias map[*
description
.
Media
]
StatsSessionMedia
}
StatsSession are session statistics.
type
StatsSessionFormat
¶
added in
v4.12.0
type StatsSessionFormat struct {
// number of RTP packets correctly received and processed
RTPPacketsReceived
uint64
// number of sent RTP packets
RTPPacketsSent
uint64
// number of lost RTP packets
RTPPacketsLost
uint64
// mean jitter of received RTP packets
RTPPacketsJitter
float64
// local SSRC
LocalSSRC
uint32
// remote SSRC
RemoteSSRC
uint32
// last sequence number of incoming/outgoing RTP packets
RTPPacketsLastSequenceNumber
uint16
// last RTP time of incoming/outgoing RTP packets
RTPPacketsLastRTP
uint32
// last NTP time of incoming/outgoing NTP packets
RTPPacketsLastNTP
time
.
Time
}
StatsSessionFormat are session format statistics.
type
StatsSessionMedia
¶
added in
v4.12.0
type StatsSessionMedia struct {
// received bytes
BytesReceived
uint64
// sent bytes
BytesSent
uint64
// number of RTP packets that could not be processed
RTPPacketsInError
uint64
// number of RTCP packets correctly received and processed
RTCPPacketsReceived
uint64
// number of sent RTCP packets
RTCPPacketsSent
uint64
// number of RTCP packets that could not be processed
RTCPPacketsInError
uint64
// format statistics
Formats map[
format
.
Format
]
StatsSessionFormat
}
StatsSessionMedia are session media statistics.
type
Transport
¶
type Transport
int
Transport is a RTSP transport protocol.
const (
TransportUDP
Transport
=
iota
TransportUDPMulticast
TransportTCP
)
transport protocols.
func (Transport)
String
¶
func (t
Transport
) String()
string
String implements fmt.Stringer.