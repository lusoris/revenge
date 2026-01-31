# gortsplib (RTSP)

> Source: https://pkg.go.dev/github.com/bluenviron/gortsplib/v4
> Fetched: 2026-01-31T11:01:05.397094+00:00
> Content-Hash: 0142b1d687eb613e
> Type: html

---

### Overview ¶

Package gortsplib is a RTSP library for the Go programming language. 

Examples are available at <https://github.com/bluenviron/gortsplib/tree/main/examples>

### Index ¶

  * type Client
  *     * func (c *Client) Announce(u *base.URL, desc *description.Session) (*base.Response, error)
    * func (c *Client) Close()
    * func (c *Client) Describe(u *base.URL) (*description.Session, *base.Response, error)
    * func (c *Client) OnPacketRTCP(medi *description.Media, cb OnPacketRTCPFunc)
    * func (c *Client) OnPacketRTCPAny(cb OnPacketRTCPAnyFunc)
    * func (c *Client) OnPacketRTP(medi *description.Media, forma format.Format, cb OnPacketRTPFunc)
    * func (c *Client) OnPacketRTPAny(cb OnPacketRTPAnyFunc)
    * func (c *Client) Options(u *base.URL) (*base.Response, error)
    * func (c *Client) PacketNTP(medi *description.Media, pkt *rtp.Packet) (time.Time, bool)
    * func (c *Client) PacketPTS(medi *description.Media, pkt *rtp.Packet) (time.Duration, bool)deprecated
    * func (c *Client) PacketPTS2(medi *description.Media, pkt *rtp.Packet) (int64, bool)
    * func (c *Client) Pause() (*base.Response, error)
    * func (c *Client) Play(ra *headers.Range) (*base.Response, error)
    * func (c *Client) Record() (*base.Response, error)
    * func (c *Client) Seek(ra *headers.Range) (*base.Response, error)deprecated
    * func (c *Client) Setup(baseURL *base.URL, media *description.Media, rtpPort int, rtcpPort int) (*base.Response, error)
    * func (c *Client) SetupAll(baseURL *base.URL, medias []*description.Media) error
    * func (c *Client) Start(scheme string, host string) errordeprecated
    * func (c *Client) Start2() error
    * func (c *Client) StartRecording(address string, desc *description.Session) error
    * func (c *Client) Stats() *ClientStats
    * func (c *Client) Wait() error
    * func (c *Client) WritePacketRTCP(medi *description.Media, pkt rtcp.Packet) error
    * func (c *Client) WritePacketRTP(medi *description.Media, pkt *rtp.Packet) error
    * func (c *Client) WritePacketRTPWithNTP(medi *description.Media, pkt *rtp.Packet, ntp time.Time) error
  * type ClientOnDecodeErrorFunc
  * type ClientOnPacketLostFuncdeprecated
  * type ClientOnPacketsLostFunc
  * type ClientOnRequestFunc
  * type ClientOnResponseFunc
  * type ClientOnTransportSwitchFunc
  * type ClientStats
  * type OnPacketRTCPAnyFunc
  * type OnPacketRTCPFunc
  * type OnPacketRTPAnyFunc
  * type OnPacketRTPFunc
  * type Server
  *     * func (s *Server) Close()
    * func (s *Server) Start() error
    * func (s *Server) StartAndWait() error
    * func (s *Server) Wait() error
  * type ServerConn
  *     * func (sc *ServerConn) BytesReceived() uint64deprecated
    * func (sc *ServerConn) BytesSent() uint64deprecated
    * func (sc *ServerConn) Close()
    * func (sc *ServerConn) NetConn() net.Conn
    * func (sc *ServerConn) Session() *ServerSession
    * func (sc *ServerConn) SetUserData(v interface{})
    * func (sc *ServerConn) Stats() *StatsConn
    * func (sc *ServerConn) UserData() interface{}
    * func (sc *ServerConn) VerifyCredentials(req *base.Request, expectedUser string, expectedPass string) bool
  * type ServerHandler
  * type ServerHandlerOnAnnounce
  * type ServerHandlerOnAnnounceCtx
  * type ServerHandlerOnConnClose
  * type ServerHandlerOnConnCloseCtx
  * type ServerHandlerOnConnOpen
  * type ServerHandlerOnConnOpenCtx
  * type ServerHandlerOnDecodeError
  * type ServerHandlerOnDecodeErrorCtx
  * type ServerHandlerOnDescribe
  * type ServerHandlerOnDescribeCtx
  * type ServerHandlerOnGetParameter
  * type ServerHandlerOnGetParameterCtx
  * type ServerHandlerOnPacketLostdeprecated
  * type ServerHandlerOnPacketLostCtxdeprecated
  * type ServerHandlerOnPacketsLost
  * type ServerHandlerOnPacketsLostCtx
  * type ServerHandlerOnPause
  * type ServerHandlerOnPauseCtx
  * type ServerHandlerOnPlay
  * type ServerHandlerOnPlayCtx
  * type ServerHandlerOnRecord
  * type ServerHandlerOnRecordCtx
  * type ServerHandlerOnRequest
  * type ServerHandlerOnResponse
  * type ServerHandlerOnSessionClose
  * type ServerHandlerOnSessionCloseCtx
  * type ServerHandlerOnSessionOpen
  * type ServerHandlerOnSessionOpenCtx
  * type ServerHandlerOnSetParameter
  * type ServerHandlerOnSetParameterCtx
  * type ServerHandlerOnSetup
  * type ServerHandlerOnSetupCtx
  * type ServerHandlerOnStreamWriteError
  * type ServerHandlerOnStreamWriteErrorCtx
  * type ServerSession
  *     * func (ss *ServerSession) AnnouncedDescription() *description.Session
    * func (ss *ServerSession) BytesReceived() uint64deprecated
    * func (ss *ServerSession) BytesSent() uint64deprecated
    * func (ss *ServerSession) Close()
    * func (ss *ServerSession) OnPacketRTCP(medi *description.Media, cb OnPacketRTCPFunc)
    * func (ss *ServerSession) OnPacketRTCPAny(cb OnPacketRTCPAnyFunc)
    * func (ss *ServerSession) OnPacketRTP(medi *description.Media, forma format.Format, cb OnPacketRTPFunc)
    * func (ss *ServerSession) OnPacketRTPAny(cb OnPacketRTPAnyFunc)
    * func (ss *ServerSession) PacketNTP(medi *description.Media, pkt *rtp.Packet) (time.Time, bool)
    * func (ss *ServerSession) PacketPTS(medi *description.Media, pkt *rtp.Packet) (time.Duration, bool)deprecated
    * func (ss *ServerSession) PacketPTS2(medi *description.Media, pkt *rtp.Packet) (int64, bool)
    * func (ss *ServerSession) SetUserData(v interface{})
    * func (ss *ServerSession) SetuppedMedias() []*description.Media
    * func (ss *ServerSession) SetuppedPath() string
    * func (ss *ServerSession) SetuppedQuery() string
    * func (ss *ServerSession) SetuppedSecure() bool
    * func (ss *ServerSession) SetuppedStream() *ServerStream
    * func (ss *ServerSession) SetuppedTransport() *Transport
    * func (ss *ServerSession) State() ServerSessionState
    * func (ss *ServerSession) Stats() *StatsSession
    * func (ss *ServerSession) UserData() interface{}
    * func (ss *ServerSession) WritePacketRTCP(medi *description.Media, pkt rtcp.Packet) error
    * func (ss *ServerSession) WritePacketRTP(medi *description.Media, pkt *rtp.Packet) error
  * type ServerSessionState
  *     * func (s ServerSessionState) String() string
  * type ServerStream
  *     * func NewServerStream(s *Server, desc *description.Session) *ServerStreamdeprecated
  *     * func (st *ServerStream) BytesSent() uint64deprecated
    * func (st *ServerStream) Close()
    * func (st *ServerStream) Description() *description.Sessiondeprecated
    * func (st *ServerStream) Initialize() error
    * func (st *ServerStream) Stats() *ServerStreamStats
    * func (st *ServerStream) WritePacketRTCP(medi *description.Media, pkt rtcp.Packet) error
    * func (st *ServerStream) WritePacketRTP(medi *description.Media, pkt *rtp.Packet) error
    * func (st *ServerStream) WritePacketRTPWithNTP(medi *description.Media, pkt *rtp.Packet, ntp time.Time) error
  * type ServerStreamStats
  * type ServerStreamStatsFormat
  * type ServerStreamStatsMedia
  * type StatsConn
  * type StatsSession
  * type StatsSessionFormat
  * type StatsSessionMedia
  * type Transport
  *     * func (t Transport) String() string



### Constants ¶

This section is empty.

### Variables ¶

This section is empty.

### Functions ¶

This section is empty.

### Types ¶

####  type [Client](https://github.com/bluenviron/gortsplib/blob/v4.16.2/client.go#L399) ¶
    
    
    type Client struct {
    	//
    	// Target
    	//
    	// Scheme. Either "rtsp" or "rtsps".
    	Scheme [string](/builtin#string)
    	// Host and port.
    	Host [string](/builtin#string)
    
    	//
    	// RTSP parameters (all optional)
    	//
    	// timeout of read operations.
    	// It defaults to 10 seconds.
    	ReadTimeout [time](/time).[Duration](/time#Duration)
    	// timeout of write operations.
    	// It defaults to 10 seconds.
    	WriteTimeout [time](/time).[Duration](/time#Duration)
    	// a TLS configuration to connect to TLS (RTSPS) servers.
    	// It defaults to nil.
    	TLSConfig *[tls](/crypto/tls).[Config](/crypto/tls#Config)
    	// enable communication with servers which don't provide UDP server ports
    	// or use different server ports than the announced ones.
    	// This can be a security issue.
    	// It defaults to false.
    	AnyPortEnable [bool](/builtin#bool)
    	// transport protocol (UDP, Multicast or TCP).
    	// If nil, it is chosen automatically (first UDP, then, if it fails, TCP).
    	// It defaults to nil.
    	Transport *Transport
    	// If the client is reading with UDP, it must receive
    	// at least a packet within this timeout, otherwise it switches to TCP.
    	// It defaults to 3 seconds.
    	InitialUDPReadTimeout [time](/time).[Duration](/time#Duration)
    	// Size of the UDP read buffer.
    	// This can be increased to reduce packet losses.
    	// It defaults to the operating system default value.
    	UDPReadBufferSize [int](/builtin#int)
    	// Size of the queue of outgoing packets.
    	// It defaults to 256.
    	WriteQueueSize [int](/builtin#int)
    	// maximum size of outgoing RTP / RTCP packets.
    	// This must be less than the UDP MTU (1472 bytes).
    	// It defaults to 1472.
    	MaxPacketSize [int](/builtin#int)
    	// user agent header.
    	// It defaults to "gortsplib"
    	UserAgent [string](/builtin#string)
    	// disable automatic RTCP sender reports.
    	DisableRTCPSenderReports [bool](/builtin#bool)
    	// explicitly request back channels to the server.
    	RequestBackChannels [bool](/builtin#bool)
    	// pointer to a variable that stores received bytes.
    	//
    	// Deprecated: use Client.Stats()
    	BytesReceived *[uint64](/builtin#uint64)
    	// pointer to a variable that stores sent bytes.
    	//
    	// Deprecated: use Client.Stats()
    	BytesSent *[uint64](/builtin#uint64)
    
    	//
    	// system functions (all optional)
    	//
    	// function used to initialize the TCP client.
    	// It defaults to (&net.Dialer{}).DialContext.
    	DialContext func(ctx [context](/context).[Context](/context#Context), network, address [string](/builtin#string)) ([net](/net).[Conn](/net#Conn), [error](/builtin#error))
    	// function used to initialize UDP listeners.
    	// It defaults to net.ListenPacket.
    	ListenPacket func(network, address [string](/builtin#string)) ([net](/net).[PacketConn](/net#PacketConn), [error](/builtin#error))
    
    	//
    	// callbacks (all optional)
    	//
    	// called when sending a request to the server.
    	OnRequest ClientOnRequestFunc
    	// called when receiving a response from the server.
    	OnResponse ClientOnResponseFunc
    	// called when receiving a request from the server.
    	OnServerRequest ClientOnRequestFunc
    	// called when sending a response to the server.
    	OnServerResponse ClientOnResponseFunc
    	// called when the transport protocol changes.
    	OnTransportSwitch ClientOnTransportSwitchFunc
    	// called when the client detects lost packets.
    	//
    	// Deprecated: replaced by OnPacketsLost
    	OnPacketLost ClientOnPacketLostFunc
    	// called when the client detects lost packets.
    	OnPacketsLost ClientOnPacketsLostFunc
    	// called when a non-fatal decode error occurs.
    	OnDecodeError ClientOnDecodeErrorFunc
    	// contains filtered or unexported fields
    }

Client is a RTSP client. 

####  func (*Client) [Announce](https://github.com/bluenviron/gortsplib/blob/v4.16.2/client.go#L1541) ¶
    
    
    func (c *Client) Announce(u *[base](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/base).[URL](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/base#URL), desc *[description](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/description).[Session](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/description#Session)) (*[base](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/base).[Response](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/base#Response), [error](/builtin#error))

Announce sends an ANNOUNCE request. 

####  func (*Client) [Close](https://github.com/bluenviron/gortsplib/blob/v4.16.2/client.go#L724) ¶
    
    
    func (c *Client) Close()

Close closes all client resources and waits for them to exit. 

####  func (*Client) [Describe](https://github.com/bluenviron/gortsplib/blob/v4.16.2/client.go#L1471) ¶
    
    
    func (c *Client) Describe(u *[base](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/base).[URL](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/base#URL)) (*[description](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/description).[Session](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/description#Session), *[base](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/base).[Response](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/base#Response), [error](/builtin#error))

Describe sends a DESCRIBE request. 

####  func (*Client) [OnPacketRTCP](https://github.com/bluenviron/gortsplib/blob/v4.16.2/client.go#L2284) ¶
    
    
    func (c *Client) OnPacketRTCP(medi *[description](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/description).[Media](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/description#Media), cb OnPacketRTCPFunc)

OnPacketRTCP sets a callback that is called when a RTCP packet is read. 

####  func (*Client) [OnPacketRTCPAny](https://github.com/bluenviron/gortsplib/blob/v4.16.2/client.go#L2267) ¶
    
    
    func (c *Client) OnPacketRTCPAny(cb OnPacketRTCPAnyFunc)

OnPacketRTCPAny sets a callback that is called when a RTCP packet is read from any setupped media. 

####  func (*Client) [OnPacketRTP](https://github.com/bluenviron/gortsplib/blob/v4.16.2/client.go#L2277) ¶
    
    
    func (c *Client) OnPacketRTP(medi *[description](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/description).[Media](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/description#Media), forma [format](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/format).[Format](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/format#Format), cb OnPacketRTPFunc)

OnPacketRTP sets a callback that is called when a RTP packet is read. 

####  func (*Client) [OnPacketRTPAny](https://github.com/bluenviron/gortsplib/blob/v4.16.2/client.go#L2255) ¶
    
    
    func (c *Client) OnPacketRTPAny(cb OnPacketRTPAnyFunc)

OnPacketRTPAny sets a callback that is called when a RTP packet is read from any setupped media. 

####  func (*Client) [Options](https://github.com/bluenviron/gortsplib/blob/v4.16.2/client.go#L1360) ¶
    
    
    func (c *Client) Options(u *[base](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/base).[URL](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/base#URL)) (*[base](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/base).[Response](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/base#Response), [error](/builtin#error))

Options sends an OPTIONS request. 

####  func (*Client) [PacketNTP](https://github.com/bluenviron/gortsplib/blob/v4.16.2/client.go#L2360) ¶
    
    
    func (c *Client) PacketNTP(medi *[description](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/description).[Media](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/description#Media), pkt *[rtp](/github.com/pion/rtp).[Packet](/github.com/pion/rtp#Packet)) ([time](/time).[Time](/time#Time), [bool](/builtin#bool))

PacketNTP returns the NTP (absolute timestamp) of an incoming RTP packet. The NTP is computed from RTCP sender reports. 

####  func (*Client) [PacketPTS](https://github.com/bluenviron/gortsplib/blob/v4.16.2/client.go#L2338) deprecated
    
    
    func (c *Client) PacketPTS(medi *[description](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/description).[Media](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/description#Media), pkt *[rtp](/github.com/pion/rtp).[Packet](/github.com/pion/rtp#Packet)) ([time](/time).[Duration](/time#Duration), [bool](/builtin#bool))

PacketPTS returns the PTS (presentation timestamp) of an incoming RTP packet. It is computed by decoding the packet timestamp and sychronizing it with other tracks. 

Deprecated: replaced by PacketPTS2. 

####  func (*Client) [PacketPTS2](https://github.com/bluenviron/gortsplib/blob/v4.16.2/client.go#L2352) ¶ added in v4.11.0
    
    
    func (c *Client) PacketPTS2(medi *[description](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/description).[Media](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/description#Media), pkt *[rtp](/github.com/pion/rtp).[Packet](/github.com/pion/rtp#Packet)) ([int64](/builtin#int64), [bool](/builtin#bool))

PacketPTS2 returns the PTS (presentation timestamp) of an incoming RTP packet. It is computed by decoding the packet timestamp and sychronizing it with other tracks. 

####  func (*Client) [Pause](https://github.com/bluenviron/gortsplib/blob/v4.16.2/client.go#L2230) ¶
    
    
    func (c *Client) Pause() (*[base](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/base).[Response](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/base#Response), [error](/builtin#error))

Pause sends a PAUSE request. This can be called only after Play() or Record(). 

####  func (*Client) [Play](https://github.com/bluenviron/gortsplib/blob/v4.16.2/client.go#L2124) ¶
    
    
    func (c *Client) Play(ra *[headers](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/headers).[Range](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/headers#Range)) (*[base](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/base).[Response](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/base#Response), [error](/builtin#error))

Play sends a PLAY request. This can be called only after Setup(). 

####  func (*Client) [Record](https://github.com/bluenviron/gortsplib/blob/v4.16.2/client.go#L2175) ¶
    
    
    func (c *Client) Record() (*[base](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/base).[Response](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/base#Response), [error](/builtin#error))

Record sends a RECORD request. This can be called only after Announce() and Setup(). 

####  func (*Client) [Seek](https://github.com/bluenviron/gortsplib/blob/v4.16.2/client.go#L2245) deprecated
    
    
    func (c *Client) Seek(ra *[headers](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/headers).[Range](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/headers#Range)) (*[base](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/base).[Response](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/base#Response), [error](/builtin#error))

Seek asks the server to re-start the stream from a specific timestamp. 

Deprecated: will be removed in next version. Equivalent to using Pause() followed by Play(). 

####  func (*Client) [Setup](https://github.com/bluenviron/gortsplib/blob/v4.16.2/client.go#L1994) ¶
    
    
    func (c *Client) Setup(
    	baseURL *[base](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/base).[URL](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/base#URL),
    	media *[description](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/description).[Media](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/description#Media),
    	rtpPort [int](/builtin#int),
    	rtcpPort [int](/builtin#int),
    ) (*[base](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/base).[Response](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/base#Response), [error](/builtin#error))

Setup sends a SETUP request. rtpPort and rtcpPort are used only if transport is UDP. if rtpPort and rtcpPort are zero, they are chosen automatically. 

####  func (*Client) [SetupAll](https://github.com/bluenviron/gortsplib/blob/v4.16.2/client.go#L2018) ¶
    
    
    func (c *Client) SetupAll(baseURL *[base](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/base).[URL](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/base#URL), medias []*[description](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/description).[Media](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/description#Media)) [error](/builtin#error)

SetupAll setups all the given medias. 

####  func (*Client) [Start](https://github.com/bluenviron/gortsplib/blob/v4.16.2/client.go#L554) deprecated
    
    
    func (c *Client) Start(scheme [string](/builtin#string), host [string](/builtin#string)) [error](/builtin#error)

Start initializes the connection to a server. 

Deprecated: replaced by Start2. 

####  func (*Client) [Start2](https://github.com/bluenviron/gortsplib/blob/v4.16.2/client.go#L561) ¶ added in v4.15.0
    
    
    func (c *Client) Start2() [error](/builtin#error)

Start2 initializes the connection to a server. 

####  func (*Client) [StartRecording](https://github.com/bluenviron/gortsplib/blob/v4.16.2/client.go#L688) ¶
    
    
    func (c *Client) StartRecording(address [string](/builtin#string), desc *[description](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/description).[Session](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/description#Session)) [error](/builtin#error)

StartRecording connects to the address and starts publishing given media. 

####  func (*Client) [Stats](https://github.com/bluenviron/gortsplib/blob/v4.16.2/client.go#L2367) ¶ added in v4.12.0
    
    
    func (c *Client) Stats() *ClientStats

Stats returns client statistics. 

####  func (*Client) [Wait](https://github.com/bluenviron/gortsplib/blob/v4.16.2/client.go#L731) ¶
    
    
    func (c *Client) Wait() [error](/builtin#error)

Wait waits until all client resources are closed. This can happen when a fatal error occurs or when Close() is called. 

####  func (*Client) [WritePacketRTCP](https://github.com/bluenviron/gortsplib/blob/v4.16.2/client.go#L2316) ¶
    
    
    func (c *Client) WritePacketRTCP(medi *[description](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/description).[Media](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/description#Media), pkt [rtcp](/github.com/pion/rtcp).[Packet](/github.com/pion/rtcp#Packet)) [error](/builtin#error)

WritePacketRTCP writes a RTCP packet to the server. 

####  func (*Client) [WritePacketRTP](https://github.com/bluenviron/gortsplib/blob/v4.16.2/client.go#L2290) ¶
    
    
    func (c *Client) WritePacketRTP(medi *[description](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/description).[Media](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/description#Media), pkt *[rtp](/github.com/pion/rtp).[Packet](/github.com/pion/rtp#Packet)) [error](/builtin#error)

WritePacketRTP writes a RTP packet to the server. 

####  func (*Client) [WritePacketRTPWithNTP](https://github.com/bluenviron/gortsplib/blob/v4.16.2/client.go#L2296) ¶
    
    
    func (c *Client) WritePacketRTPWithNTP(medi *[description](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/description).[Media](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/description#Media), pkt *[rtp](/github.com/pion/rtp).[Packet](/github.com/pion/rtp#Packet), ntp [time](/time).[Time](/time#Time)) [error](/builtin#error)

WritePacketRTPWithNTP writes a RTP packet to the server. ntp is the absolute timestamp of the packet, and is sent with periodic RTCP sender reports. 

####  type [ClientOnDecodeErrorFunc](https://github.com/bluenviron/gortsplib/blob/v4.16.2/client.go#L384) ¶
    
    
    type ClientOnDecodeErrorFunc func(err [error](/builtin#error))

ClientOnDecodeErrorFunc is the prototype of Client.OnDecodeError. 

####  type [ClientOnPacketLostFunc](https://github.com/bluenviron/gortsplib/blob/v4.16.2/client.go#L378) deprecated
    
    
    type ClientOnPacketLostFunc func(err [error](/builtin#error))

ClientOnPacketLostFunc is the prototype of Client.OnPacketLost. 

Deprecated: replaced by ClientOnPacketsLostFunc 

####  type [ClientOnPacketsLostFunc](https://github.com/bluenviron/gortsplib/blob/v4.16.2/client.go#L381) ¶ added in v4.13.0
    
    
    type ClientOnPacketsLostFunc func(lost [uint64](/builtin#uint64))

ClientOnPacketsLostFunc is the prototype of Client.OnPacketsLost. 

####  type [ClientOnRequestFunc](https://github.com/bluenviron/gortsplib/blob/v4.16.2/client.go#L367) ¶
    
    
    type ClientOnRequestFunc func(*[base](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/base).[Request](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/base#Request))

ClientOnRequestFunc is the prototype of Client.OnRequest. 

####  type [ClientOnResponseFunc](https://github.com/bluenviron/gortsplib/blob/v4.16.2/client.go#L370) ¶
    
    
    type ClientOnResponseFunc func(*[base](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/base).[Response](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/base#Response))

ClientOnResponseFunc is the prototype of Client.OnResponse. 

####  type [ClientOnTransportSwitchFunc](https://github.com/bluenviron/gortsplib/blob/v4.16.2/client.go#L373) ¶
    
    
    type ClientOnTransportSwitchFunc func(err [error](/builtin#error))

ClientOnTransportSwitchFunc is the prototype of Client.OnTransportSwitch. 

####  type [ClientStats](https://github.com/bluenviron/gortsplib/blob/v4.16.2/client_stats.go#L4) ¶ added in v4.12.0
    
    
    type ClientStats struct {
    	Conn    StatsConn
    	Session StatsSession
    }

ClientStats are client statistics 

####  type [OnPacketRTCPAnyFunc](https://github.com/bluenviron/gortsplib/blob/v4.16.2/client.go#L396) ¶
    
    
    type OnPacketRTCPAnyFunc func(*[description](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/description).[Media](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/description#Media), [rtcp](/github.com/pion/rtcp).[Packet](/github.com/pion/rtcp#Packet))

OnPacketRTCPAnyFunc is the prototype of the callback passed to OnPacketRTCPAny(). 

####  type [OnPacketRTCPFunc](https://github.com/bluenviron/gortsplib/blob/v4.16.2/client.go#L393) ¶
    
    
    type OnPacketRTCPFunc func([rtcp](/github.com/pion/rtcp).[Packet](/github.com/pion/rtcp#Packet))

OnPacketRTCPFunc is the prototype of the callback passed to OnPacketRTCP(). 

####  type [OnPacketRTPAnyFunc](https://github.com/bluenviron/gortsplib/blob/v4.16.2/client.go#L390) ¶
    
    
    type OnPacketRTPAnyFunc func(*[description](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/description).[Media](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/description#Media), [format](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/format).[Format](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/format#Format), *[rtp](/github.com/pion/rtp).[Packet](/github.com/pion/rtp#Packet))

OnPacketRTPAnyFunc is the prototype of the callback passed to OnPacketRTP(Any). 

####  type [OnPacketRTPFunc](https://github.com/bluenviron/gortsplib/blob/v4.16.2/client.go#L387) ¶
    
    
    type OnPacketRTPFunc func(*[rtp](/github.com/pion/rtp).[Packet](/github.com/pion/rtp#Packet))

OnPacketRTPFunc is the prototype of the callback passed to OnPacketRTP(). 

####  type [Server](https://github.com/bluenviron/gortsplib/blob/v4.16.2/server.go#L55) ¶
    
    
    type Server struct {
    	//
    	// RTSP parameters (all optional except RTSPAddress)
    	//
    	// the RTSP address of the server, to accept connections and send and receive
    	// packets with the TCP transport.
    	RTSPAddress [string](/builtin#string)
    	// a port to send and receive RTP packets with the UDP transport.
    	// If UDPRTPAddress and UDPRTCPAddress are filled, the server can support the UDP transport.
    	UDPRTPAddress [string](/builtin#string)
    	// a port to send and receive RTCP packets with the UDP transport.
    	// If UDPRTPAddress and UDPRTCPAddress are filled, the server can support the UDP transport.
    	UDPRTCPAddress [string](/builtin#string)
    	// a range of multicast IPs to use with the UDP-multicast transport.
    	// If MulticastIPRange, MulticastRTPPort, MulticastRTCPPort are filled, the server
    	// can support the UDP-multicast transport.
    	MulticastIPRange [string](/builtin#string)
    	// a port to send RTP packets with the UDP-multicast transport.
    	// If MulticastIPRange, MulticastRTPPort, MulticastRTCPPort are filled, the server
    	// can support the UDP-multicast transport.
    	MulticastRTPPort [int](/builtin#int)
    	// a port to send RTCP packets with the UDP-multicast transport.
    	// If MulticastIPRange, MulticastRTPPort, MulticastRTCPPort are filled, the server
    	// can support the UDP-multicast transport.
    	MulticastRTCPPort [int](/builtin#int)
    	// timeout of read operations.
    	// It defaults to 10 seconds
    	ReadTimeout [time](/time).[Duration](/time#Duration)
    	// timeout of write operations.
    	// It defaults to 10 seconds
    	WriteTimeout [time](/time).[Duration](/time#Duration)
    	// a TLS configuration to accept TLS (RTSPS) connections.
    	TLSConfig *[tls](/crypto/tls).[Config](/crypto/tls#Config)
    	// Size of the UDP read buffer.
    	// This can be increased to reduce packet losses.
    	// It defaults to the operating system default value.
    	UDPReadBufferSize [int](/builtin#int)
    	// Size of the queue of outgoing packets.
    	// It defaults to 256.
    	WriteQueueSize [int](/builtin#int)
    	// maximum size of outgoing RTP / RTCP packets.
    	// This must be less than the UDP MTU (1472 bytes).
    	// It defaults to 1472.
    	MaxPacketSize [int](/builtin#int)
    	// disable automatic RTCP sender reports.
    	DisableRTCPSenderReports [bool](/builtin#bool)
    	// authentication methods.
    	// It defaults to plain and digest+MD5.
    	AuthMethods [][auth](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/auth).[VerifyMethod](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/auth#VerifyMethod)
    
    	//
    	// handler (optional)
    	//
    	// an handler to handle server events.
    	// It may implement one or more of the ServerHandler* interfaces.
    	Handler ServerHandler
    
    	//
    	// system functions (all optional)
    	//
    	// function used to initialize the TCP listener.
    	// It defaults to net.Listen.
    	Listen func(network [string](/builtin#string), address [string](/builtin#string)) ([net](/net).[Listener](/net#Listener), [error](/builtin#error))
    	// function used to initialize UDP listeners.
    	// It defaults to net.ListenPacket.
    	ListenPacket func(network, address [string](/builtin#string)) ([net](/net).[PacketConn](/net#PacketConn), [error](/builtin#error))
    	// contains filtered or unexported fields
    }

Server is a RTSP server. 

####  func (*Server) [Close](https://github.com/bluenviron/gortsplib/blob/v4.16.2/server.go#L339) ¶
    
    
    func (s *Server) Close()

Close closes all the server resources and waits for them to exit. 

####  func (*Server) [Start](https://github.com/bluenviron/gortsplib/blob/v4.16.2/server.go#L154) ¶
    
    
    func (s *Server) Start() [error](/builtin#error)

Start starts the server. 

####  func (*Server) [StartAndWait](https://github.com/bluenviron/gortsplib/blob/v4.16.2/server.go#L471) ¶
    
    
    func (s *Server) StartAndWait() [error](/builtin#error)

StartAndWait starts the server and waits until a fatal error. 

####  func (*Server) [Wait](https://github.com/bluenviron/gortsplib/blob/v4.16.2/server.go#L346) ¶
    
    
    func (s *Server) Wait() [error](/builtin#error)

Wait waits until all server resources are closed. This can happen when a fatal error occurs or when Close() is called. 

####  type [ServerConn](https://github.com/bluenviron/gortsplib/blob/v4.16.2/server_conn.go#L195) ¶
    
    
    type ServerConn struct {
    	// contains filtered or unexported fields
    }

ServerConn is a server-side RTSP connection. 

####  func (*ServerConn) [BytesReceived](https://github.com/bluenviron/gortsplib/blob/v4.16.2/server_conn.go#L247) deprecated
    
    
    func (sc *ServerConn) BytesReceived() [uint64](/builtin#uint64)

BytesReceived returns the number of read bytes. 

Deprecated: replaced by Stats() 

####  func (*ServerConn) [BytesSent](https://github.com/bluenviron/gortsplib/blob/v4.16.2/server_conn.go#L254) deprecated
    
    
    func (sc *ServerConn) BytesSent() [uint64](/builtin#uint64)

BytesSent returns the number of written bytes. 

Deprecated: replaced by Stats() 

####  func (*ServerConn) [Close](https://github.com/bluenviron/gortsplib/blob/v4.16.2/server_conn.go#L235) ¶
    
    
    func (sc *ServerConn) Close()

Close closes the ServerConn. 

####  func (*ServerConn) [NetConn](https://github.com/bluenviron/gortsplib/blob/v4.16.2/server_conn.go#L240) ¶
    
    
    func (sc *ServerConn) NetConn() [net](/net).[Conn](/net#Conn)

NetConn returns the underlying net.Conn. 

####  func (*ServerConn) [Session](https://github.com/bluenviron/gortsplib/blob/v4.16.2/server_conn.go#L269) ¶ added in v4.12.1
    
    
    func (sc *ServerConn) Session() *ServerSession

Session returns the associated session. 

####  func (*ServerConn) [SetUserData](https://github.com/bluenviron/gortsplib/blob/v4.16.2/server_conn.go#L259) ¶
    
    
    func (sc *ServerConn) SetUserData(v interface{})

SetUserData sets some user data associated with the connection. 

####  func (*ServerConn) [Stats](https://github.com/bluenviron/gortsplib/blob/v4.16.2/server_conn.go#L274) ¶ added in v4.12.0
    
    
    func (sc *ServerConn) Stats() *StatsConn

Stats returns connection statistics. 

####  func (*ServerConn) [UserData](https://github.com/bluenviron/gortsplib/blob/v4.16.2/server_conn.go#L264) ¶
    
    
    func (sc *ServerConn) UserData() interface{}

UserData returns some user data associated with the connection. 

####  func (*ServerConn) [VerifyCredentials](https://github.com/bluenviron/gortsplib/blob/v4.16.2/server_conn.go#L282) ¶ added in v4.13.0
    
    
    func (sc *ServerConn) VerifyCredentials(
    	req *[base](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/base).[Request](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/base#Request),
    	expectedUser [string](/builtin#string),
    	expectedPass [string](/builtin#string),
    ) [bool](/builtin#bool)

VerifyCredentials verifies credentials provided by the user. 

####  type [ServerHandler](https://github.com/bluenviron/gortsplib/blob/v4.16.2/server_handler.go#L9) ¶
    
    
    type ServerHandler interface{}

ServerHandler is the interface implemented by all the server handlers. 

####  type [ServerHandlerOnAnnounce](https://github.com/bluenviron/gortsplib/blob/v4.16.2/server_handler.go#L95) ¶
    
    
    type ServerHandlerOnAnnounce interface {
    	// called when receiving an ANNOUNCE request.
    	OnAnnounce(*ServerHandlerOnAnnounceCtx) (*[base](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/base).[Response](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/base#Response), [error](/builtin#error))
    }

ServerHandlerOnAnnounce can be implemented by a ServerHandler. 

####  type [ServerHandlerOnAnnounceCtx](https://github.com/bluenviron/gortsplib/blob/v4.16.2/server_handler.go#L85) ¶
    
    
    type ServerHandlerOnAnnounceCtx struct {
    	Session     *ServerSession
    	Conn        *ServerConn
    	Request     *[base](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/base).[Request](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/base#Request)
    	Path        [string](/builtin#string)
    	Query       [string](/builtin#string)
    	Description *[description](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/description).[Session](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/description#Session)
    }

ServerHandlerOnAnnounceCtx is the context of OnAnnounce. 

####  type [ServerHandlerOnConnClose](https://github.com/bluenviron/gortsplib/blob/v4.16.2/server_handler.go#L29) ¶
    
    
    type ServerHandlerOnConnClose interface {
    	// called when a connection is closed.
    	OnConnClose(*ServerHandlerOnConnCloseCtx)
    }

ServerHandlerOnConnClose can be implemented by a ServerHandler. 

####  type [ServerHandlerOnConnCloseCtx](https://github.com/bluenviron/gortsplib/blob/v4.16.2/server_handler.go#L23) ¶
    
    
    type ServerHandlerOnConnCloseCtx struct {
    	Conn  *ServerConn
    	Error [error](/builtin#error)
    }

ServerHandlerOnConnCloseCtx is the context of OnConnClose. 

####  type [ServerHandlerOnConnOpen](https://github.com/bluenviron/gortsplib/blob/v4.16.2/server_handler.go#L17) ¶
    
    
    type ServerHandlerOnConnOpen interface {
    	// called when a connection is opened.
    	OnConnOpen(*ServerHandlerOnConnOpenCtx)
    }

ServerHandlerOnConnOpen can be implemented by a ServerHandler. 

####  type [ServerHandlerOnConnOpenCtx](https://github.com/bluenviron/gortsplib/blob/v4.16.2/server_handler.go#L12) ¶
    
    
    type ServerHandlerOnConnOpenCtx struct {
    	Conn *ServerConn
    }

ServerHandlerOnConnOpenCtx is the context of OnConnOpen. 

####  type [ServerHandlerOnDecodeError](https://github.com/bluenviron/gortsplib/blob/v4.16.2/server_handler.go#L230) ¶
    
    
    type ServerHandlerOnDecodeError interface {
    	// called when a non-fatal decode error occurs.
    	OnDecodeError(*ServerHandlerOnDecodeErrorCtx)
    }

ServerHandlerOnDecodeError can be implemented by a ServerHandler. 

####  type [ServerHandlerOnDecodeErrorCtx](https://github.com/bluenviron/gortsplib/blob/v4.16.2/server_handler.go#L224) ¶
    
    
    type ServerHandlerOnDecodeErrorCtx struct {
    	Session *ServerSession
    	Error   [error](/builtin#error)
    }

ServerHandlerOnDecodeErrorCtx is the context of OnDecodeError. 

####  type [ServerHandlerOnDescribe](https://github.com/bluenviron/gortsplib/blob/v4.16.2/server_handler.go#L79) ¶
    
    
    type ServerHandlerOnDescribe interface {
    	// called when receiving a DESCRIBE request.
    	OnDescribe(*ServerHandlerOnDescribeCtx) (*[base](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/base).[Response](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/base#Response), *ServerStream, [error](/builtin#error))
    }

ServerHandlerOnDescribe can be implemented by a ServerHandler. 

####  type [ServerHandlerOnDescribeCtx](https://github.com/bluenviron/gortsplib/blob/v4.16.2/server_handler.go#L71) ¶
    
    
    type ServerHandlerOnDescribeCtx struct {
    	Conn    *ServerConn
    	Request *[base](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/base).[Request](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/base#Request)
    	Path    [string](/builtin#string)
    	Query   [string](/builtin#string)
    }

ServerHandlerOnDescribeCtx is the context of OnDescribe. 

####  type [ServerHandlerOnGetParameter](https://github.com/bluenviron/gortsplib/blob/v4.16.2/server_handler.go#L175) ¶
    
    
    type ServerHandlerOnGetParameter interface {
    	// called when receiving a GET_PARAMETER request.
    	OnGetParameter(*ServerHandlerOnGetParameterCtx) (*[base](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/base).[Response](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/base#Response), [error](/builtin#error))
    }

ServerHandlerOnGetParameter can be implemented by a ServerHandler. 

####  type [ServerHandlerOnGetParameterCtx](https://github.com/bluenviron/gortsplib/blob/v4.16.2/server_handler.go#L166) ¶
    
    
    type ServerHandlerOnGetParameterCtx struct {
    	Session *ServerSession
    	Conn    *ServerConn
    	Request *[base](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/base).[Request](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/base#Request)
    	Path    [string](/builtin#string)
    	Query   [string](/builtin#string)
    }

ServerHandlerOnGetParameterCtx is the context of OnGetParameter. 

####  type [ServerHandlerOnPacketLost](https://github.com/bluenviron/gortsplib/blob/v4.16.2/server_handler.go#L206) deprecated
    
    
    type ServerHandlerOnPacketLost interface {
    	// called when the server detects lost packets.
    	OnPacketLost(*ServerHandlerOnPacketLostCtx)
    }

ServerHandlerOnPacketLost can be implemented by a ServerHandler. 

Deprecated: replaced by ServerHandlerOnPacketsLost 

####  type [ServerHandlerOnPacketLostCtx](https://github.com/bluenviron/gortsplib/blob/v4.16.2/server_handler.go#L198) deprecated
    
    
    type ServerHandlerOnPacketLostCtx struct {
    	Session *ServerSession
    	Error   [error](/builtin#error)
    }

ServerHandlerOnPacketLostCtx is the context of OnPacketLost. 

Deprecated: replaced by ServerHandlerOnPacketsLostCtx 

####  type [ServerHandlerOnPacketsLost](https://github.com/bluenviron/gortsplib/blob/v4.16.2/server_handler.go#L218) ¶ added in v4.13.0
    
    
    type ServerHandlerOnPacketsLost interface {
    	// called when the server detects lost packets.
    	OnPacketsLost(*ServerHandlerOnPacketsLostCtx)
    }

ServerHandlerOnPacketsLost can be implemented by a ServerHandler. 

####  type [ServerHandlerOnPacketsLostCtx](https://github.com/bluenviron/gortsplib/blob/v4.16.2/server_handler.go#L212) ¶ added in v4.13.0
    
    
    type ServerHandlerOnPacketsLostCtx struct {
    	Session *ServerSession
    	Lost    [uint64](/builtin#uint64)
    }

ServerHandlerOnPacketsLostCtx is the context of OnPacketsLost. 

####  type [ServerHandlerOnPause](https://github.com/bluenviron/gortsplib/blob/v4.16.2/server_handler.go#L160) ¶
    
    
    type ServerHandlerOnPause interface {
    	// called when receiving a PAUSE request.
    	OnPause(*ServerHandlerOnPauseCtx) (*[base](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/base).[Response](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/base#Response), [error](/builtin#error))
    }

ServerHandlerOnPause can be implemented by a ServerHandler. 

####  type [ServerHandlerOnPauseCtx](https://github.com/bluenviron/gortsplib/blob/v4.16.2/server_handler.go#L151) ¶
    
    
    type ServerHandlerOnPauseCtx struct {
    	Session *ServerSession
    	Conn    *ServerConn
    	Request *[base](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/base).[Request](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/base#Request)
    	Path    [string](/builtin#string)
    	Query   [string](/builtin#string)
    }

ServerHandlerOnPauseCtx is the context of OnPause. 

####  type [ServerHandlerOnPlay](https://github.com/bluenviron/gortsplib/blob/v4.16.2/server_handler.go#L130) ¶
    
    
    type ServerHandlerOnPlay interface {
    	// called when receiving a PLAY request.
    	OnPlay(*ServerHandlerOnPlayCtx) (*[base](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/base).[Response](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/base#Response), [error](/builtin#error))
    }

ServerHandlerOnPlay can be implemented by a ServerHandler. 

####  type [ServerHandlerOnPlayCtx](https://github.com/bluenviron/gortsplib/blob/v4.16.2/server_handler.go#L121) ¶
    
    
    type ServerHandlerOnPlayCtx struct {
    	Session *ServerSession
    	Conn    *ServerConn
    	Request *[base](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/base).[Request](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/base#Request)
    	Path    [string](/builtin#string)
    	Query   [string](/builtin#string)
    }

ServerHandlerOnPlayCtx is the context of OnPlay. 

####  type [ServerHandlerOnRecord](https://github.com/bluenviron/gortsplib/blob/v4.16.2/server_handler.go#L145) ¶
    
    
    type ServerHandlerOnRecord interface {
    	// called when receiving a RECORD request.
    	OnRecord(*ServerHandlerOnRecordCtx) (*[base](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/base).[Response](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/base#Response), [error](/builtin#error))
    }

ServerHandlerOnRecord can be implemented by a ServerHandler. 

####  type [ServerHandlerOnRecordCtx](https://github.com/bluenviron/gortsplib/blob/v4.16.2/server_handler.go#L136) ¶
    
    
    type ServerHandlerOnRecordCtx struct {
    	Session *ServerSession
    	Conn    *ServerConn
    	Request *[base](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/base).[Request](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/base#Request)
    	Path    [string](/builtin#string)
    	Query   [string](/builtin#string)
    }

ServerHandlerOnRecordCtx is the context of OnRecord. 

####  type [ServerHandlerOnRequest](https://github.com/bluenviron/gortsplib/blob/v4.16.2/server_handler.go#L59) ¶
    
    
    type ServerHandlerOnRequest interface {
    	// called when receiving a request from a connection.
    	OnRequest(*ServerConn, *[base](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/base).[Request](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/base#Request))
    }

ServerHandlerOnRequest can be implemented by a ServerHandler. 

####  type [ServerHandlerOnResponse](https://github.com/bluenviron/gortsplib/blob/v4.16.2/server_handler.go#L65) ¶
    
    
    type ServerHandlerOnResponse interface {
    	// called when sending a response to a connection.
    	OnResponse(*ServerConn, *[base](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/base).[Response](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/base#Response))
    }

ServerHandlerOnResponse can be implemented by a ServerHandler. 

####  type [ServerHandlerOnSessionClose](https://github.com/bluenviron/gortsplib/blob/v4.16.2/server_handler.go#L53) ¶
    
    
    type ServerHandlerOnSessionClose interface {
    	// called when a session is closed.
    	OnSessionClose(*ServerHandlerOnSessionCloseCtx)
    }

ServerHandlerOnSessionClose can be implemented by a ServerHandler. 

####  type [ServerHandlerOnSessionCloseCtx](https://github.com/bluenviron/gortsplib/blob/v4.16.2/server_handler.go#L47) ¶
    
    
    type ServerHandlerOnSessionCloseCtx struct {
    	Session *ServerSession
    	Error   [error](/builtin#error)
    }

ServerHandlerOnSessionCloseCtx is the context of ServerHandlerOnSessionClose. 

####  type [ServerHandlerOnSessionOpen](https://github.com/bluenviron/gortsplib/blob/v4.16.2/server_handler.go#L41) ¶
    
    
    type ServerHandlerOnSessionOpen interface {
    	// called when a session is opened.
    	OnSessionOpen(*ServerHandlerOnSessionOpenCtx)
    }

ServerHandlerOnSessionOpen can be implemented by a ServerHandler. 

####  type [ServerHandlerOnSessionOpenCtx](https://github.com/bluenviron/gortsplib/blob/v4.16.2/server_handler.go#L35) ¶
    
    
    type ServerHandlerOnSessionOpenCtx struct {
    	Session *ServerSession
    	Conn    *ServerConn
    }

ServerHandlerOnSessionOpenCtx is the context OnSessionOpen. 

####  type [ServerHandlerOnSetParameter](https://github.com/bluenviron/gortsplib/blob/v4.16.2/server_handler.go#L190) ¶
    
    
    type ServerHandlerOnSetParameter interface {
    	// called when receiving a SET_PARAMETER request.
    	OnSetParameter(*ServerHandlerOnSetParameterCtx) (*[base](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/base).[Response](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/base#Response), [error](/builtin#error))
    }

ServerHandlerOnSetParameter can be implemented by a ServerHandler. 

####  type [ServerHandlerOnSetParameterCtx](https://github.com/bluenviron/gortsplib/blob/v4.16.2/server_handler.go#L181) ¶
    
    
    type ServerHandlerOnSetParameterCtx struct {
    	Session *ServerSession
    	Conn    *ServerConn
    	Request *[base](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/base).[Request](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/base#Request)
    	Path    [string](/builtin#string)
    	Query   [string](/builtin#string)
    }

ServerHandlerOnSetParameterCtx is the context of OnSetParameter. 

####  type [ServerHandlerOnSetup](https://github.com/bluenviron/gortsplib/blob/v4.16.2/server_handler.go#L111) ¶
    
    
    type ServerHandlerOnSetup interface {
    	// called when receiving a SETUP request.
    	// must return a Response and a stream.
    	// the stream is needed to
    	// - add the session the the stream's readers
    	// - send the stream SSRC to the session
    	OnSetup(*ServerHandlerOnSetupCtx) (*[base](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/base).[Response](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/base#Response), *ServerStream, [error](/builtin#error))
    }

ServerHandlerOnSetup can be implemented by a ServerHandler. 

####  type [ServerHandlerOnSetupCtx](https://github.com/bluenviron/gortsplib/blob/v4.16.2/server_handler.go#L101) ¶
    
    
    type ServerHandlerOnSetupCtx struct {
    	Session   *ServerSession
    	Conn      *ServerConn
    	Request   *[base](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/base).[Request](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/base#Request)
    	Path      [string](/builtin#string)
    	Query     [string](/builtin#string)
    	Transport Transport
    }

ServerHandlerOnSetupCtx is the context of OnSetup. 

####  type [ServerHandlerOnStreamWriteError](https://github.com/bluenviron/gortsplib/blob/v4.16.2/server_handler.go#L242) ¶
    
    
    type ServerHandlerOnStreamWriteError interface {
    	// called when a ServerStream is unable to write packets to a session.
    	OnStreamWriteError(*ServerHandlerOnStreamWriteErrorCtx)
    }

ServerHandlerOnStreamWriteError can be implemented by a ServerHandler. 

####  type [ServerHandlerOnStreamWriteErrorCtx](https://github.com/bluenviron/gortsplib/blob/v4.16.2/server_handler.go#L236) ¶
    
    
    type ServerHandlerOnStreamWriteErrorCtx struct {
    	Session *ServerSession
    	Error   [error](/builtin#error)
    }

ServerHandlerOnStreamWriteErrorCtx is the context of OnStreamWriteError. 

####  type [ServerSession](https://github.com/bluenviron/gortsplib/blob/v4.16.2/server_session.go#L423) ¶
    
    
    type ServerSession struct {
    	// contains filtered or unexported fields
    }

ServerSession is a server-side RTSP session. 

####  func (*ServerSession) [AnnouncedDescription](https://github.com/bluenviron/gortsplib/blob/v4.16.2/server_session.go#L539) ¶
    
    
    func (ss *ServerSession) AnnouncedDescription() *[description](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/description).[Session](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/description#Session)

AnnouncedDescription returns the announced stream description. 

####  func (*ServerSession) [BytesReceived](https://github.com/bluenviron/gortsplib/blob/v4.16.2/server_session.go#L487) deprecated
    
    
    func (ss *ServerSession) BytesReceived() [uint64](/builtin#uint64)

BytesReceived returns the number of read bytes. 

Deprecated: replaced by Stats() 

####  func (*ServerSession) [BytesSent](https://github.com/bluenviron/gortsplib/blob/v4.16.2/server_session.go#L498) deprecated
    
    
    func (ss *ServerSession) BytesSent() [uint64](/builtin#uint64)

BytesSent returns the number of written bytes. 

Deprecated: replaced by Stats() 

####  func (*ServerSession) [Close](https://github.com/bluenviron/gortsplib/blob/v4.16.2/server_session.go#L480) ¶
    
    
    func (ss *ServerSession) Close()

Close closes the ServerSession. 

####  func (*ServerSession) [OnPacketRTCP](https://github.com/bluenviron/gortsplib/blob/v4.16.2/server_session.go#L1738) ¶
    
    
    func (ss *ServerSession) OnPacketRTCP(medi *[description](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/description).[Media](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/description#Media), cb OnPacketRTCPFunc)

OnPacketRTCP sets a callback that is called when a RTCP packet is read. 

####  func (*ServerSession) [OnPacketRTCPAny](https://github.com/bluenviron/gortsplib/blob/v4.16.2/server_session.go#L1721) ¶
    
    
    func (ss *ServerSession) OnPacketRTCPAny(cb OnPacketRTCPAnyFunc)

OnPacketRTCPAny sets a callback that is called when a RTCP packet is read from any setupped media. 

####  func (*ServerSession) [OnPacketRTP](https://github.com/bluenviron/gortsplib/blob/v4.16.2/server_session.go#L1731) ¶
    
    
    func (ss *ServerSession) OnPacketRTP(medi *[description](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/description).[Media](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/description#Media), forma [format](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/format).[Format](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/format#Format), cb OnPacketRTPFunc)

OnPacketRTP sets a callback that is called when a RTP packet is read. 

####  func (*ServerSession) [OnPacketRTPAny](https://github.com/bluenviron/gortsplib/blob/v4.16.2/server_session.go#L1709) ¶
    
    
    func (ss *ServerSession) OnPacketRTPAny(cb OnPacketRTPAnyFunc)

OnPacketRTPAny sets a callback that is called when a RTP packet is read from any setupped media. 

####  func (*ServerSession) [PacketNTP](https://github.com/bluenviron/gortsplib/blob/v4.16.2/server_session.go#L1782) ¶
    
    
    func (ss *ServerSession) PacketNTP(medi *[description](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/description).[Media](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/description#Media), pkt *[rtp](/github.com/pion/rtp).[Packet](/github.com/pion/rtp#Packet)) ([time](/time).[Time](/time#Time), [bool](/builtin#bool))

PacketNTP returns the NTP (absolute timestamp) of an incoming RTP packet. The NTP is computed from RTCP sender reports. 

####  func (*ServerSession) [PacketPTS](https://github.com/bluenviron/gortsplib/blob/v4.16.2/server_session.go#L1760) deprecated
    
    
    func (ss *ServerSession) PacketPTS(medi *[description](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/description).[Media](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/description#Media), pkt *[rtp](/github.com/pion/rtp).[Packet](/github.com/pion/rtp#Packet)) ([time](/time).[Duration](/time#Duration), [bool](/builtin#bool))

PacketPTS returns the PTS (presentation timestamp) of an incoming RTP packet. It is computed by decoding the packet timestamp and sychronizing it with other tracks. 

Deprecated: replaced by PacketPTS2. 

####  func (*ServerSession) [PacketPTS2](https://github.com/bluenviron/gortsplib/blob/v4.16.2/server_session.go#L1774) ¶ added in v4.11.0
    
    
    func (ss *ServerSession) PacketPTS2(medi *[description](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/description).[Media](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/description#Media), pkt *[rtp](/github.com/pion/rtp).[Packet](/github.com/pion/rtp#Packet)) ([int64](/builtin#int64), [bool](/builtin#bool))

PacketPTS2 returns the PTS (presentation timestamp) of an incoming RTP packet. It is computed by decoding the packet timestamp and sychronizing it with other tracks. 

####  func (*ServerSession) [SetUserData](https://github.com/bluenviron/gortsplib/blob/v4.16.2/server_session.go#L553) ¶
    
    
    func (ss *ServerSession) SetUserData(v interface{})

SetUserData sets some user data associated with the session. 

####  func (*ServerSession) [SetuppedMedias](https://github.com/bluenviron/gortsplib/blob/v4.16.2/server_session.go#L544) ¶
    
    
    func (ss *ServerSession) SetuppedMedias() []*[description](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/description).[Media](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/description#Media)

SetuppedMedias returns the setupped medias. 

####  func (*ServerSession) [SetuppedPath](https://github.com/bluenviron/gortsplib/blob/v4.16.2/server_session.go#L529) ¶ added in v4.3.0
    
    
    func (ss *ServerSession) SetuppedPath() [string](/builtin#string)

SetuppedPath returns the path sent during SETUP or ANNOUNCE. 

####  func (*ServerSession) [SetuppedQuery](https://github.com/bluenviron/gortsplib/blob/v4.16.2/server_session.go#L534) ¶ added in v4.3.0
    
    
    func (ss *ServerSession) SetuppedQuery() [string](/builtin#string)

SetuppedQuery returns the query sent during SETUP or ANNOUNCE. 

####  func (*ServerSession) [SetuppedSecure](https://github.com/bluenviron/gortsplib/blob/v4.16.2/server_session.go#L519) ¶ added in v4.15.0
    
    
    func (ss *ServerSession) SetuppedSecure() [bool](/builtin#bool)

SetuppedSecure returns whether a secure profile is in use. If this is false, it does not mean that the stream is not secure, since there are some combinations that are secure nonetheless, like RTSPS+TCP+unsecure. 

####  func (*ServerSession) [SetuppedStream](https://github.com/bluenviron/gortsplib/blob/v4.16.2/server_session.go#L524) ¶ added in v4.3.0
    
    
    func (ss *ServerSession) SetuppedStream() *ServerStream

SetuppedStream returns the stream associated with the session. 

####  func (*ServerSession) [SetuppedTransport](https://github.com/bluenviron/gortsplib/blob/v4.16.2/server_session.go#L512) ¶
    
    
    func (ss *ServerSession) SetuppedTransport() *Transport

SetuppedTransport returns the transport negotiated during SETUP. 

####  func (*ServerSession) [State](https://github.com/bluenviron/gortsplib/blob/v4.16.2/server_session.go#L507) ¶
    
    
    func (ss *ServerSession) State() ServerSessionState

State returns the state of the session. 

####  func (*ServerSession) [Stats](https://github.com/bluenviron/gortsplib/blob/v4.16.2/server_session.go#L563) ¶ added in v4.12.0
    
    
    func (ss *ServerSession) Stats() *StatsSession

Stats returns server session statistics. 

####  func (*ServerSession) [UserData](https://github.com/bluenviron/gortsplib/blob/v4.16.2/server_session.go#L558) ¶
    
    
    func (ss *ServerSession) UserData() interface{}

UserData returns some user data associated with the session. 

####  func (*ServerSession) [WritePacketRTCP](https://github.com/bluenviron/gortsplib/blob/v4.16.2/server_session.go#L1751) ¶
    
    
    func (ss *ServerSession) WritePacketRTCP(medi *[description](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/description).[Media](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/description#Media), pkt [rtcp](/github.com/pion/rtcp).[Packet](/github.com/pion/rtcp#Packet)) [error](/builtin#error)

WritePacketRTCP writes a RTCP packet to the session. 

####  func (*ServerSession) [WritePacketRTP](https://github.com/bluenviron/gortsplib/blob/v4.16.2/server_session.go#L1744) ¶
    
    
    func (ss *ServerSession) WritePacketRTP(medi *[description](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/description).[Media](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/description#Media), pkt *[rtp](/github.com/pion/rtp).[Packet](/github.com/pion/rtp#Packet)) [error](/builtin#error)

WritePacketRTP writes a RTP packet to the session. 

####  type [ServerSessionState](https://github.com/bluenviron/gortsplib/blob/v4.16.2/server_session.go#L394) ¶
    
    
    type ServerSessionState [int](/builtin#int)

ServerSessionState is a state of a ServerSession. 
    
    
    const (
    	ServerSessionStateInitial ServerSessionState = [iota](/builtin#iota)
    	ServerSessionStatePrePlay
    	ServerSessionStatePlay
    	ServerSessionStatePreRecord
    	ServerSessionStateRecord
    )

states. 

####  func (ServerSessionState) [String](https://github.com/bluenviron/gortsplib/blob/v4.16.2/server_session.go#L406) ¶
    
    
    func (s ServerSessionState) String() [string](/builtin#string)

String implements fmt.Stringer. 

####  type [ServerStream](https://github.com/bluenviron/gortsplib/blob/v4.16.2/server_stream.go#L37) ¶
    
    
    type ServerStream struct {
    	Server *Server
    	Desc   *[description](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/description).[Session](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/description#Session)
    	// contains filtered or unexported fields
    }

ServerStream represents a data stream. This is in charge of - storing stream description and statistics - distributing the stream to each reader - allocating multicast listeners 

####  func [NewServerStream](https://github.com/bluenviron/gortsplib/blob/v4.16.2/server_stream.go#L20) deprecated
    
    
    func NewServerStream(s *Server, desc *[description](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/description).[Session](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/description#Session)) *ServerStream

NewServerStream allocates a ServerStream. 

Deprecated: replaced by ServerStream.Initialize(). 

####  func (*ServerStream) [BytesSent](https://github.com/bluenviron/gortsplib/blob/v4.16.2/server_stream.go#L97) deprecated added in v4.4.0
    
    
    func (st *ServerStream) BytesSent() [uint64](/builtin#uint64)

BytesSent returns the number of written bytes. 

Deprecated: replaced by Stats() 

####  func (*ServerStream) [Close](https://github.com/bluenviron/gortsplib/blob/v4.16.2/server_stream.go#L80) ¶
    
    
    func (st *ServerStream) Close()

Close closes a ServerStream. 

####  func (*ServerStream) [Description](https://github.com/bluenviron/gortsplib/blob/v4.16.2/server_stream.go#L108) deprecated
    
    
    func (st *ServerStream) Description() *[description](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/description).[Session](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/description#Session)

Description returns the description of the stream. 

Deprecated: use ServerStream.Desc. 

####  func (*ServerStream) [Initialize](https://github.com/bluenviron/gortsplib/blob/v4.16.2/server_stream.go#L50) ¶ added in v4.13.0
    
    
    func (st *ServerStream) Initialize() [error](/builtin#error)

Initialize initializes a ServerStream. 

####  func (*ServerStream) [Stats](https://github.com/bluenviron/gortsplib/blob/v4.16.2/server_stream.go#L113) ¶ added in v4.12.0
    
    
    func (st *ServerStream) Stats() *ServerStreamStats

Stats returns stream statistics. 

####  func (*ServerStream) [WritePacketRTCP](https://github.com/bluenviron/gortsplib/blob/v4.16.2/server_stream.go#L293) ¶
    
    
    func (st *ServerStream) WritePacketRTCP(medi *[description](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/description).[Media](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/description#Media), pkt [rtcp](/github.com/pion/rtcp).[Packet](/github.com/pion/rtcp#Packet)) [error](/builtin#error)

WritePacketRTCP writes a RTCP packet to all the readers of the stream. 

####  func (*ServerStream) [WritePacketRTP](https://github.com/bluenviron/gortsplib/blob/v4.16.2/server_stream.go#L273) ¶
    
    
    func (st *ServerStream) WritePacketRTP(medi *[description](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/description).[Media](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/description#Media), pkt *[rtp](/github.com/pion/rtp).[Packet](/github.com/pion/rtp#Packet)) [error](/builtin#error)

WritePacketRTP writes a RTP packet to all the readers of the stream. 

####  func (*ServerStream) [WritePacketRTPWithNTP](https://github.com/bluenviron/gortsplib/blob/v4.16.2/server_stream.go#L279) ¶
    
    
    func (st *ServerStream) WritePacketRTPWithNTP(medi *[description](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/description).[Media](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/description#Media), pkt *[rtp](/github.com/pion/rtp).[Packet](/github.com/pion/rtp#Packet), ntp [time](/time).[Time](/time#Time)) [error](/builtin#error)

WritePacketRTPWithNTP writes a RTP packet to all the readers of the stream. ntp is the absolute timestamp of the packet, and is sent with periodic RTCP sender reports. 

####  type [ServerStreamStats](https://github.com/bluenviron/gortsplib/blob/v4.16.2/server_stream_stats.go#L28) ¶ added in v4.12.0
    
    
    type ServerStreamStats struct {
    	// sent bytes
    	BytesSent [uint64](/builtin#uint64)
    	// number of sent RTP packets
    	RTPPacketsSent [uint64](/builtin#uint64)
    	// number of sent RTCP packets
    	RTCPPacketsSent [uint64](/builtin#uint64)
    
    	// media statistics
    	Medias map[*[description](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/description).[Media](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/description#Media)]ServerStreamStatsMedia
    }

ServerStreamStats are stream statistics. 

####  type [ServerStreamStatsFormat](https://github.com/bluenviron/gortsplib/blob/v4.16.2/server_stream_stats.go#L9) ¶ added in v4.12.0
    
    
    type ServerStreamStatsFormat struct {
    	// number of sent RTP packets
    	RTPPacketsSent [uint64](/builtin#uint64)
    	// local SSRC
    	LocalSSRC [uint32](/builtin#uint32)
    }

ServerStreamStatsFormat are stream format statistics. 

####  type [ServerStreamStatsMedia](https://github.com/bluenviron/gortsplib/blob/v4.16.2/server_stream_stats.go#L17) ¶ added in v4.12.0
    
    
    type ServerStreamStatsMedia struct {
    	// sent bytes
    	BytesSent [uint64](/builtin#uint64)
    	// number of sent RTCP packets
    	RTCPPacketsSent [uint64](/builtin#uint64)
    
    	// format statistics
    	Formats map[[format](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/format).[Format](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/format#Format)]ServerStreamStatsFormat
    }

ServerStreamStatsMedia are stream media statistics. 

####  type [StatsConn](https://github.com/bluenviron/gortsplib/blob/v4.16.2/stats_conn.go#L4) ¶ added in v4.12.0
    
    
    type StatsConn struct {
    	// received bytes
    	BytesReceived [uint64](/builtin#uint64)
    	// sent bytes
    	BytesSent [uint64](/builtin#uint64)
    }

StatsConn are connection statistics. 

####  type [StatsSession](https://github.com/bluenviron/gortsplib/blob/v4.16.2/stats_session.go#L52) ¶ added in v4.12.0
    
    
    type StatsSession struct {
    	// received bytes
    	BytesReceived [uint64](/builtin#uint64)
    	// sent bytes
    	BytesSent [uint64](/builtin#uint64)
    	// number of RTP packets correctly received and processed
    	RTPPacketsReceived [uint64](/builtin#uint64)
    	// number of sent RTP packets
    	RTPPacketsSent [uint64](/builtin#uint64)
    	// number of lost RTP packets
    	RTPPacketsLost [uint64](/builtin#uint64)
    	// number of RTP packets that could not be processed
    	RTPPacketsInError [uint64](/builtin#uint64)
    	// mean jitter of received RTP packets
    	RTPPacketsJitter [float64](/builtin#float64)
    	// number of RTCP packets correctly received and processed
    	RTCPPacketsReceived [uint64](/builtin#uint64)
    	// number of sent RTCP packets
    	RTCPPacketsSent [uint64](/builtin#uint64)
    	// number of RTCP packets that could not be processed
    	RTCPPacketsInError [uint64](/builtin#uint64)
    
    	// media statistics
    	Medias map[*[description](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/description).[Media](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/description#Media)]StatsSessionMedia
    }

StatsSession are session statistics. 

####  type [StatsSessionFormat](https://github.com/bluenviron/gortsplib/blob/v4.16.2/stats_session.go#L11) ¶ added in v4.12.0
    
    
    type StatsSessionFormat struct {
    	// number of RTP packets correctly received and processed
    	RTPPacketsReceived [uint64](/builtin#uint64)
    	// number of sent RTP packets
    	RTPPacketsSent [uint64](/builtin#uint64)
    	// number of lost RTP packets
    	RTPPacketsLost [uint64](/builtin#uint64)
    	// mean jitter of received RTP packets
    	RTPPacketsJitter [float64](/builtin#float64)
    	// local SSRC
    	LocalSSRC [uint32](/builtin#uint32)
    	// remote SSRC
    	RemoteSSRC [uint32](/builtin#uint32)
    	// last sequence number of incoming/outgoing RTP packets
    	RTPPacketsLastSequenceNumber [uint16](/builtin#uint16)
    	// last RTP time of incoming/outgoing RTP packets
    	RTPPacketsLastRTP [uint32](/builtin#uint32)
    	// last NTP time of incoming/outgoing NTP packets
    	RTPPacketsLastNTP [time](/time).[Time](/time#Time)
    }

StatsSessionFormat are session format statistics. 

####  type [StatsSessionMedia](https://github.com/bluenviron/gortsplib/blob/v4.16.2/stats_session.go#L33) ¶ added in v4.12.0
    
    
    type StatsSessionMedia struct {
    	// received bytes
    	BytesReceived [uint64](/builtin#uint64)
    	// sent bytes
    	BytesSent [uint64](/builtin#uint64)
    	// number of RTP packets that could not be processed
    	RTPPacketsInError [uint64](/builtin#uint64)
    	// number of RTCP packets correctly received and processed
    	RTCPPacketsReceived [uint64](/builtin#uint64)
    	// number of sent RTCP packets
    	RTCPPacketsSent [uint64](/builtin#uint64)
    	// number of RTCP packets that could not be processed
    	RTCPPacketsInError [uint64](/builtin#uint64)
    
    	// format statistics
    	Formats map[[format](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/format).[Format](/github.com/bluenviron/gortsplib/v4@v4.16.2/pkg/format#Format)]StatsSessionFormat
    }

StatsSessionMedia are session media statistics. 

####  type [Transport](https://github.com/bluenviron/gortsplib/blob/v4.16.2/transport.go#L4) ¶
    
    
    type Transport [int](/builtin#int)

Transport is a RTSP transport protocol. 
    
    
    const (
    	TransportUDP Transport = [iota](/builtin#iota)
    	TransportUDPMulticast
    	TransportTCP
    )

transport protocols. 

####  func (Transport) [String](https://github.com/bluenviron/gortsplib/blob/v4.16.2/transport.go#L20) ¶
    
    
    func (t Transport) String() [string](/builtin#string)

String implements fmt.Stringer. 
  *[v]: View this template
  *[t]: Discuss this template
  *[e]: Edit this template
