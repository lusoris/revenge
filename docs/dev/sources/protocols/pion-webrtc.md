# Pion WebRTC (Go)

> Source: https://pkg.go.dev/github.com/pion/webrtc/v4
> Fetched: 2026-02-01T11:46:44.979329+00:00
> Content-Hash: 7011c90eb7f9ee48
> Type: html

---

### Overview ¶

Package webrtc implements the WebRTC 1.0 as defined in W3C WebRTC specification document.

### Index ¶

- Constants
- Variables
- func ConfigureCongestionControlFeedback(mediaEngine *MediaEngine, interceptorRegistry*interceptor.Registry) error
- func ConfigureFlexFEC03(payloadType PayloadType, mediaEngine *MediaEngine, ...) error
- func ConfigureNack(mediaEngine *MediaEngine, interceptorRegistry*interceptor.Registry) error
- func ConfigureRTCPReports(interceptorRegistry *interceptor.Registry) error
- func ConfigureSimulcastExtensionHeaders(mediaEngine *MediaEngine) error
- func ConfigureStatsInterceptor(interceptorRegistry *interceptor.Registry) error
- func ConfigureTWCCHeaderExtensionSender(mediaEngine *MediaEngine, interceptorRegistry*interceptor.Registry) error
- func ConfigureTWCCSender(mediaEngine *MediaEngine, interceptorRegistry*interceptor.Registry) error
- func GatheringCompletePromise(pc *PeerConnection) (gatherComplete <-chan struct{})
- func NewAudioPlayoutStatsProvider(id string) *defaultAudioPlayoutStatsProvider
- func NewICETCPMux(logger logging.LeveledLogger, listener net.Listener, readBufferSize int) ice.TCPMux
- func NewICEUDPMux(logger logging.LeveledLogger, udpConn net.PacketConn) ice.UDPMux
- func RegisterDefaultInterceptors(mediaEngine *MediaEngine, interceptorRegistry*interceptor.Registry) error
- func WithInterceptorRegistry(ir *interceptor.Registry) func(a*API)
- func WithMediaEngine(m *MediaEngine) func(a*API)
- func WithPayloader(h func(RTPCodecCapability) (rtp.Payloader, error)) func(*TrackLocalStaticRTP)
- func WithRTPSequenceNumber(sequenceNumber uint16) func(*TrackLocalStaticRTP)
- func WithRTPStreamID(rid string) func(*TrackLocalStaticRTP)
- func WithRTPTimestamp(timestamp uint32) func(*TrackLocalStaticRTP)
- func WithSettingEngine(s SettingEngine) func(a *API)
- type API
-     * func NewAPI(options ...func(*API)) *API
-     * func (api *API) NewDTLSTransport(transport *ICETransport, certificates []Certificate) (*DTLSTransport, error)
  - func (api *API) NewDataChannel(transport*SCTPTransport, params *DataChannelParameters) (*DataChannel, error)
  - func (api *API) NewICEGatherer(opts ICEGatherOptions) (*ICEGatherer, error)
  - func (api *API) NewICETransport(gatherer*ICEGatherer) *ICETransport
  - func (api *API) NewPeerConnection(configuration Configuration) (*PeerConnection, error)
  - func (api *API) NewRTPReceiver(kind RTPCodecType, transport*DTLSTransport) (*RTPReceiver, error)
  - func (api *API) NewRTPSender(track TrackLocal, transport*DTLSTransport) (*RTPSender, error)
  - func (api *API) NewSCTPTransport(dtls*DTLSTransport) *SCTPTransport
- type AnswerOptions
- type AudioPlayoutStats
- type AudioPlayoutStatsProvider
- type AudioReceiverStats
- type AudioSenderStats
- type AudioSourceStats
- type BundlePolicy
-     * func (t BundlePolicy) MarshalJSON() ([]byte, error)
  - func (t BundlePolicy) String() string
  - func (t *BundlePolicy) UnmarshalJSON(b []byte) error
- type Certificate
-     * func CertificateFromPEM(pems string) (*Certificate, error)
  - func CertificateFromX509(privateKey crypto.PrivateKey, certificate *x509.Certificate) Certificate
  - func GenerateCertificate(secretKey crypto.PrivateKey) (*Certificate, error)
  - func NewCertificate(key crypto.PrivateKey, tpl x509.Certificate) (*Certificate, error)
-     * func (c Certificate) Equals(cert Certificate) bool
  - func (c Certificate) Expires() time.Time
  - func (c Certificate) GetFingerprints() ([]DTLSFingerprint, error)
  - func (c Certificate) PEM() (string, error)
- type CertificateStats
- type CodecStats
- type CodecType
- type Configuration
- type DTLSFingerprint
- type DTLSParameters
- type DTLSRole
-     * func (r DTLSRole) String() string
- type DTLSTransport
-     * func (t *DTLSTransport) GetLocalParameters() (DTLSParameters, error)
  - func (t *DTLSTransport) GetRemoteCertificate() []byte
  - func (t *DTLSTransport) ICETransport()*ICETransport
  - func (t *DTLSTransport) OnStateChange(f func(DTLSTransportState))
  - func (t *DTLSTransport) Start(remoteParameters DTLSParameters) error
  - func (t *DTLSTransport) State() DTLSTransportState
  - func (t *DTLSTransport) Stop() error
  - func (t *DTLSTransport) WriteRTCP(pkts []rtcp.Packet) (int, error)
- type DTLSTransportState
-     * func (t DTLSTransportState) MarshalText() ([]byte, error)
  - func (t DTLSTransportState) String() string
  - func (t *DTLSTransportState) UnmarshalText(b []byte) error
- type DataChannel
-     * func (d *DataChannel) BufferedAmount() uint64
  - func (d *DataChannel) BufferedAmountLowThreshold() uint64
  - func (d *DataChannel) Close() error
  - func (d *DataChannel) Detach() (datachannel.ReadWriteCloser, error)
  - func (d *DataChannel) DetachWithDeadline() (datachannel.ReadWriteCloserDeadliner, error)
  - func (d *DataChannel) GracefulClose() error
  - func (d *DataChannel) ID()*uint16
  - func (d *DataChannel) Label() string
  - func (d *DataChannel) MaxPacketLifeTime()*uint16
  - func (d *DataChannel) MaxRetransmits()*uint16
  - func (d *DataChannel) Negotiated() bool
  - func (d *DataChannel) OnBufferedAmountLow(f func())
  - func (d *DataChannel) OnClose(f func())
  - func (d *DataChannel) OnDial(f func())
  - func (d *DataChannel) OnError(f func(err error))
  - func (d *DataChannel) OnMessage(f func(msg DataChannelMessage))
  - func (d *DataChannel) OnOpen(f func())
  - func (d *DataChannel) Ordered() bool
  - func (d *DataChannel) Protocol() string
  - func (d *DataChannel) ReadyState() DataChannelState
  - func (d *DataChannel) Send(data []byte) error
  - func (d *DataChannel) SendText(s string) error
  - func (d *DataChannel) SetBufferedAmountLowThreshold(th uint64)
  - func (d *DataChannel) Transport()*SCTPTransport
- type DataChannelInit
- type DataChannelMessage
- type DataChannelParameters
- type DataChannelState
-     * func (t DataChannelState) MarshalText() ([]byte, error)
  - func (t DataChannelState) String() string
  - func (t *DataChannelState) UnmarshalText(b []byte) error
- type DataChannelStats
- type ICEAddressRewriteMode
- type ICEAddressRewriteRule
- type ICECandidate
-     * func (c ICECandidate) String() string
  - func (c ICECandidate) ToICE() (cand ice.Candidate, err error)
  - func (c ICECandidate) ToJSON() ICECandidateInit
- type ICECandidateInit
- type ICECandidatePair
-     * func NewICECandidatePair(local, remote *ICECandidate) *ICECandidatePair
-     * func (p *ICECandidatePair) String() string
- type ICECandidatePairStats
- type ICECandidateStats
- type ICECandidateType
-     * func NewICECandidateType(raw string) (ICECandidateType, error)
-     * func (t ICECandidateType) MarshalText() ([]byte, error)
  - func (t ICECandidateType) String() string
  - func (t *ICECandidateType) UnmarshalText(b []byte) error
- type ICEComponent
-     * func (t ICEComponent) String() string
- type ICEConnectionState
-     * func NewICEConnectionState(raw string) ICEConnectionState
-     * func (c ICEConnectionState) String() string
- type ICECredentialType
-     * func (t ICECredentialType) MarshalJSON() ([]byte, error)
  - func (t ICECredentialType) String() string
  - func (t *ICECredentialType) UnmarshalJSON(b []byte) error
- type ICEGatherOptions
- type ICEGatherPolicy
- type ICEGatherer
-     * func (g *ICEGatherer) Close() error
  - func (g *ICEGatherer) Gather() error
  - func (g *ICEGatherer) GetLocalCandidates() ([]ICECandidate, error)
  - func (g *ICEGatherer) GetLocalParameters() (ICEParameters, error)
  - func (g *ICEGatherer) GracefulClose() error
  - func (g *ICEGatherer) OnLocalCandidate(f func(*ICECandidate))
  - func (g *ICEGatherer) OnStateChange(f func(ICEGathererState))
  - func (g *ICEGatherer) State() ICEGathererState
- type ICEGathererState
-     * func (s ICEGathererState) String() string
- type ICEGatheringState
-     * func NewICEGatheringState(raw string) ICEGatheringState
-     * func (t ICEGatheringState) String() string
- type ICEParameters
- type ICEProtocol
-     * func NewICEProtocol(raw string) (ICEProtocol, error)
-     * func (t ICEProtocol) String() string
- type ICERole
-     * func (t ICERole) MarshalText() ([]byte, error)
  - func (t ICERole) String() string
  - func (t *ICERole) UnmarshalText(b []byte) error
- type ICEServer
-     * func (s ICEServer) MarshalJSON() ([]byte, error)
  - func (s *ICEServer) UnmarshalJSON(b []byte) error
- type ICETransport
-     * func NewICETransport(gatherer *ICEGatherer, loggerFactory logging.LoggerFactory) *ICETransport
-     * func (t *ICETransport) AddRemoteCandidate(remoteCandidate *ICECandidate) error
  - func (t *ICETransport) GetLocalParameters() (ICEParameters, error)
  - func (t *ICETransport) GetRemoteParameters() (ICEParameters, error)
  - func (t *ICETransport) GetSelectedCandidatePair() (*ICECandidatePair, error)
  - func (t *ICETransport) GetSelectedCandidatePairStats() (ICECandidatePairStats, bool)
  - func (t *ICETransport) GracefulStop() error
  - func (t *ICETransport) OnConnectionStateChange(f func(ICETransportState))
  - func (t *ICETransport) OnSelectedCandidatePairChange(f func(*ICECandidatePair))
  - func (t *ICETransport) Role() ICERole
  - func (t *ICETransport) SetRemoteCandidates(remoteCandidates []ICECandidate) error
  - func (t *ICETransport) Start(gatherer*ICEGatherer, params ICEParameters, role *ICERole) error
  - func (t *ICETransport) State() ICETransportState
  - func (t *ICETransport) Stop() error
- type ICETransportPolicy
-     * func NewICETransportPolicy(raw string) ICETransportPolicy
-     * func (t ICETransportPolicy) MarshalJSON() ([]byte, error)
  - func (t ICETransportPolicy) String() string
  - func (t *ICETransportPolicy) UnmarshalJSON(b []byte) error
- type ICETransportState
-     * func (c ICETransportState) MarshalText() ([]byte, error)
  - func (c ICETransportState) String() string
  - func (c *ICETransportState) UnmarshalText(b []byte) error
- type ICETrickleCapability
-     * func (t ICETrickleCapability) String() string
- type InboundRTPStreamStats
- type MediaEngine
-     * func (m *MediaEngine) RegisterCodec(codec RTPCodecParameters, typ RTPCodecType) error
  - func (m *MediaEngine) RegisterDefaultCodecs() error
  - func (m *MediaEngine) RegisterFeedback(feedback RTCPFeedback, typ RTPCodecType)
  - func (m *MediaEngine) RegisterHeaderExtension(extension RTPHeaderExtensionCapability, typ RTPCodecType, ...) error
- type MediaKind
- type MediaStreamStats
- type NetworkType
-     * func NewNetworkType(raw string) (NetworkType, error)
-     * func (t NetworkType) Protocol() string
  - func (t NetworkType) String() string
- type NominationValueGenerator
- type OAuthCredential
- type OfferAnswerOptions
- type OfferOptions
- type OutboundRTPStreamStats
- type PayloadType
- type PeerConnection
-     * func NewPeerConnection(configuration Configuration) (*PeerConnection, error)
-     * func (pc *PeerConnection) AddICECandidate(candidate ICECandidateInit) error
  - func (pc *PeerConnection) AddTrack(track TrackLocal) (*RTPSender, error)
  - func (pc *PeerConnection) AddTransceiverFromKind(kind RTPCodecType, init ...RTPTransceiverInit) (t*RTPTransceiver, err error)
  - func (pc *PeerConnection) AddTransceiverFromTrack(track TrackLocal, init ...RTPTransceiverInit) (t*RTPTransceiver, err error)
  - func (pc *PeerConnection) CanTrickleICECandidates() ICETrickleCapability
  - func (pc *PeerConnection) Close() error
  - func (pc *PeerConnection) ConnectionState() PeerConnectionState
  - func (pc *PeerConnection) CreateAnswer(options*AnswerOptions) (SessionDescription, error)
  - func (pc *PeerConnection) CreateDataChannel(label string, options*DataChannelInit) (*DataChannel, error)
  - func (pc *PeerConnection) CreateOffer(options*OfferOptions) (SessionDescription, error)
  - func (pc *PeerConnection) CurrentLocalDescription()*SessionDescription
  - func (pc *PeerConnection) CurrentRemoteDescription()*SessionDescription
  - func (pc *PeerConnection) GetConfiguration() Configuration
  - func (pc *PeerConnection) GetReceivers() (receivers []*RTPReceiver)
  - func (pc *PeerConnection) GetSenders() (result []*RTPSender)
  - func (pc *PeerConnection) GetStats() StatsReport
  - func (pc *PeerConnection) GetTransceivers() []*RTPTransceiver
  - func (pc *PeerConnection) GracefulClose() error
  - func (pc *PeerConnection) ICEConnectionState() ICEConnectionState
  - func (pc *PeerConnection) ICEGatheringState() ICEGatheringState
  - func (pc *PeerConnection) ID() string
  - func (pc *PeerConnection) LocalDescription()*SessionDescription
  - func (pc *PeerConnection) OnConnectionStateChange(f func(PeerConnectionState))
  - func (pc *PeerConnection) OnDataChannel(f func(*DataChannel))
  - func (pc *PeerConnection) OnICECandidate(f func(*ICECandidate))
  - func (pc *PeerConnection) OnICEConnectionStateChange(f func(ICEConnectionState))
  - func (pc *PeerConnection) OnICEGatheringStateChange(f func(ICEGatheringState))
  - func (pc *PeerConnection) OnNegotiationNeeded(f func())
  - func (pc *PeerConnection) OnSignalingStateChange(f func(SignalingState))
  - func (pc *PeerConnection) OnTrack(f func(*TrackRemote, *RTPReceiver))
  - func (pc *PeerConnection) PendingLocalDescription()*SessionDescription
  - func (pc *PeerConnection) PendingRemoteDescription()*SessionDescription
  - func (pc *PeerConnection) RemoteDescription()*SessionDescription
  - func (pc *PeerConnection) RemoveTrack(sender*RTPSender) (err error)
  - func (pc *PeerConnection) SCTP()*SCTPTransport
  - func (pc *PeerConnection) SetConfiguration(configuration Configuration) error
  - func (pc *PeerConnection) SetIdentityProvider(string) error
  - func (pc *PeerConnection) SetLocalDescription(desc SessionDescription) error
  - func (pc *PeerConnection) SetRemoteDescription(desc SessionDescription) error
  - func (pc *PeerConnection) SignalingState() SignalingState
  - func (pc *PeerConnection) WriteRTCP(pkts []rtcp.Packet) error
- type PeerConnectionState
-     * func (t PeerConnectionState) String() string
- type PeerConnectionStats
- type QualityLimitationReason
- type RTCPFeedback
- type RTCPMuxPolicy
-     * func (t RTCPMuxPolicy) MarshalJSON() ([]byte, error)
  - func (t RTCPMuxPolicy) String() string
  - func (t *RTCPMuxPolicy) UnmarshalJSON(b []byte) error
- type RTPCapabilities
- type RTPCodecCapability
- type RTPCodecParameters
- type RTPCodecType
-     * func NewRTPCodecType(r string) RTPCodecType
-     * func (t RTPCodecType) String() string
- type RTPCodingParameters
- type RTPContributingSourceStats
- type RTPDecodingParameters
- type RTPEncodingParameters
- type RTPFecParameters
- type RTPHeaderExtensionCapability
- type RTPHeaderExtensionParameter
- type RTPParameters
- type RTPReceiveParameters
- type RTPReceiver
-     * func (r *RTPReceiver) GetParameters() RTPParameters
  - func (r *RTPReceiver) RTPTransceiver()*RTPTransceiver
  - func (r *RTPReceiver) Read(b []byte) (n int, a interceptor.Attributes, err error)
  - func (r *RTPReceiver) ReadRTCP() ([]rtcp.Packet, interceptor.Attributes, error)
  - func (r *RTPReceiver) ReadSimulcast(b []byte, rid string) (n int, a interceptor.Attributes, err error)
  - func (r *RTPReceiver) ReadSimulcastRTCP(rid string) ([]rtcp.Packet, interceptor.Attributes, error)
  - func (r *RTPReceiver) Receive(parameters RTPReceiveParameters) error
  - func (r *RTPReceiver) SetRTPParameters(params RTPParameters)
  - func (r *RTPReceiver) SetReadDeadline(t time.Time) error
  - func (r *RTPReceiver) SetReadDeadlineSimulcast(deadline time.Time, rid string) error
  - func (r *RTPReceiver) Stop() error
  - func (r *RTPReceiver) Track()*TrackRemote
  - func (r *RTPReceiver) Tracks() []*TrackRemote
  - func (r *RTPReceiver) Transport()*DTLSTransport
- type RTPRtxParameters
- type RTPSendParameters
- type RTPSender
-     * func (r *RTPSender) AddEncoding(track TrackLocal) error
  - func (r *RTPSender) GetParameters() RTPSendParameters
  - func (r *RTPSender) Read(b []byte) (n int, a interceptor.Attributes, err error)
  - func (r *RTPSender) ReadRTCP() ([]rtcp.Packet, interceptor.Attributes, error)
  - func (r *RTPSender) ReadSimulcast(b []byte, rid string) (n int, a interceptor.Attributes, err error)
  - func (r *RTPSender) ReadSimulcastRTCP(rid string) ([]rtcp.Packet, interceptor.Attributes, error)
  - func (r *RTPSender) ReplaceTrack(track TrackLocal) error
  - func (r *RTPSender) Send(parameters RTPSendParameters) error
  - func (r *RTPSender) SetReadDeadline(t time.Time) error
  - func (r *RTPSender) SetReadDeadlineSimulcast(deadline time.Time, rid string) error
  - func (r *RTPSender) Stop() error
  - func (r *RTPSender) Track() TrackLocal
  - func (r *RTPSender) Transport()*DTLSTransport
- type RTPTransceiver
-     * func (t *RTPTransceiver) Direction() RTPTransceiverDirection
  - func (t *RTPTransceiver) Kind() RTPCodecType
  - func (t *RTPTransceiver) Mid() string
  - func (t *RTPTransceiver) Receiver()*RTPReceiver
  - func (t *RTPTransceiver) Sender()*RTPSender
  - func (t *RTPTransceiver) SetCodecPreferences(codecs []RTPCodecParameters) error
  - func (t *RTPTransceiver) SetMid(mid string) error
  - func (t *RTPTransceiver) SetSender(s*RTPSender, track TrackLocal) error
  - func (t *RTPTransceiver) Stop() error
- type RTPTransceiverDirection
-     * func NewRTPTransceiverDirection(raw string) RTPTransceiverDirection
-     * func (t RTPTransceiverDirection) Revers() RTPTransceiverDirection
  - func (t RTPTransceiverDirection) String() string
- type RTPTransceiverInit
- type RemoteInboundRTPStreamStats
- type RemoteOutboundRTPStreamStats
- type RenominationOption
-     * func WithRenominationGenerator(generator NominationValueGenerator) RenominationOption
  - func WithRenominationInterval(interval time.Duration) RenominationOption
- type SCTPCapabilities
- type SCTPTransport
-     * func (r *SCTPTransport) BufferedAmount() int
  - func (r *SCTPTransport) GetCapabilities() SCTPCapabilities
  - func (r *SCTPTransport) MaxChannels() uint16
  - func (r *SCTPTransport) OnClose(f func(err error))
  - func (r *SCTPTransport) OnDataChannel(f func(*DataChannel))
  - func (r *SCTPTransport) OnDataChannelOpened(f func(*DataChannel))
  - func (r *SCTPTransport) OnError(f func(err error))
  - func (r *SCTPTransport) Start(capabilities SCTPCapabilities) error
  - func (r *SCTPTransport) State() SCTPTransportState
  - func (r *SCTPTransport) Stop() error
  - func (r *SCTPTransport) Transport()*DTLSTransport
- type SCTPTransportState
-     * func (s SCTPTransportState) String() string
- type SCTPTransportStats
- type SDPSemantics
-     * func (s SDPSemantics) MarshalJSON() ([]byte, error)
  - func (s SDPSemantics) String() string
  - func (s *SDPSemantics) UnmarshalJSON(b []byte) error
- type SDPType
-     * func NewSDPType(raw string) SDPType
-     * func (t SDPType) MarshalJSON() ([]byte, error)
  - func (t SDPType) String() string
  - func (t *SDPType) UnmarshalJSON(b []byte) error
- type SSRC
- type SenderAudioTrackAttachmentStats
- type SenderVideoTrackAttachmentStats
- type SessionDescription
-     * func (sd *SessionDescription) Unmarshal() (*sdp.SessionDescription, error)
- type SettingEngine
-     * func (e *SettingEngine) DetachDataChannels()
  - func (e *SettingEngine) DisableActiveTCP(isDisabled bool)
  - func (e *SettingEngine) DisableCertificateFingerprintVerification(isDisabled bool)
  - func (e *SettingEngine) DisableCloseByDTLS(isEnabled bool)
  - func (e *SettingEngine) DisableMediaEngineCopy(isDisabled bool)
  - func (e *SettingEngine) DisableMediaEngineMultipleCodecs(isDisabled bool)
  - func (e *SettingEngine) DisableSRTCPReplayProtection(isDisabled bool)
  - func (e *SettingEngine) DisableSRTPReplayProtection(isDisabled bool)
  - func (e *SettingEngine) EnableDataChannelBlockWrite(nonblockWrite bool)
  - func (e *SettingEngine) EnableSCTPZeroChecksum(isEnabled bool)
  - func (e *SettingEngine) SetAnsweringDTLSRole(role DTLSRole) error
  - func (e *SettingEngine) SetDTLSCertificateRequestMessageHook(hook func(handshake.MessageCertificateRequest) handshake.Message)
  - func (e *SettingEngine) SetDTLSCipherSuites(cipherSuites ...dtls.CipherSuiteID)
  - func (e *SettingEngine) SetDTLSClientAuth(clientAuth dtls.ClientAuthType)
  - func (e *SettingEngine) SetDTLSClientCAs(clientCAs*x509.CertPool)
  - func (e *SettingEngine) SetDTLSClientHelloMessageHook(hook func(handshake.MessageClientHello) handshake.Message)
  - func (e *SettingEngine) SetDTLSConnectContextMaker(connectContextMaker func() (context.Context, func()))
  - func (e *SettingEngine) SetDTLSCustomerCipherSuites(customCipherSuites func() []dtls.CipherSuite)
  - func (e *SettingEngine) SetDTLSDisableInsecureSkipVerify(disable bool)
  - func (e *SettingEngine) SetDTLSEllipticCurves(ellipticCurves ...dtlsElliptic.Curve)
  - func (e *SettingEngine) SetDTLSExtendedMasterSecret(extendedMasterSecret dtls.ExtendedMasterSecretType)
  - func (e *SettingEngine) SetDTLSInsecureSkipHelloVerify(skip bool)
  - func (e *SettingEngine) SetDTLSKeyLogWriter(writer io.Writer)
  - func (e *SettingEngine) SetDTLSReplayProtectionWindow(n uint)
  - func (e *SettingEngine) SetDTLSRetransmissionInterval(interval time.Duration)
  - func (e *SettingEngine) SetDTLSRootCAs(rootCAs*x509.CertPool)
  - func (e *SettingEngine) SetDTLSServerHelloMessageHook(hook func(handshake.MessageServerHello) handshake.Message)
  - func (e *SettingEngine) SetEphemeralUDPPortRange(portMin, portMax uint16) error
  - func (e *SettingEngine) SetFireOnTrackBeforeFirstRTP(fireOnTrackBeforeFirstRTP bool)
  - func (e *SettingEngine) SetHandleUndeclaredSSRCWithoutAnswer(handleUndeclaredSSRCWithoutAnswer bool)
  - func (e *SettingEngine) SetHostAcceptanceMinWait(t time.Duration)
  - func (e *SettingEngine) SetICEAddressRewriteRules(rules ...ICEAddressRewriteRule) error
  - func (e *SettingEngine) SetICEBindingRequestHandler(...)
  - func (e *SettingEngine) SetICECredentials(usernameFragment, password string)
  - func (e *SettingEngine) SetICEMaxBindingRequests(d uint16)
  - func (e *SettingEngine) SetICEMulticastDNSMode(multicastDNSMode ice.MulticastDNSMode)
  - func (e *SettingEngine) SetICEProxyDialer(d proxy.Dialer)
  - func (e *SettingEngine) SetICERenomination(options ...RenominationOption) error
  - func (e *SettingEngine) SetICETCPMux(tcpMux ice.TCPMux)
  - func (e *SettingEngine) SetICETimeouts(disconnectedTimeout, failedTimeout, keepAliveInterval time.Duration)
  - func (e *SettingEngine) SetICEUDPMux(udpMux ice.UDPMux)
  - func (e *SettingEngine) SetIPFilter(filter func(net.IP) (keep bool))
  - func (e *SettingEngine) SetIgnoreRidPauseForRecv(ignoreRidPauseForRecv bool)
  - func (e *SettingEngine) SetIncludeLoopbackCandidate(include bool)
  - func (e *SettingEngine) SetInterfaceFilter(filter func(string) (keep bool))
  - func (e *SettingEngine) SetLite(lite bool)
  - func (e *SettingEngine) SetMulticastDNSHostName(hostName string)
  - func (e *SettingEngine) SetNAT1To1IPs(ips []string, candidateType ICECandidateType)deprecated
  - func (e *SettingEngine) SetNet(net transport.Net)
  - func (e *SettingEngine) SetNetworkTypes(candidateTypes []NetworkType)
  - func (e *SettingEngine) SetPrflxAcceptanceMinWait(t time.Duration)
  - func (e *SettingEngine) SetReceiveMTU(receiveMTU uint)
  - func (e *SettingEngine) SetRelayAcceptanceMinWait(t time.Duration)
  - func (e *SettingEngine) SetSCTPCwndCAStep(cwndCAStep uint32)
  - func (e *SettingEngine) SetSCTPFastRtxWnd(fastRtxWnd uint32)
  - func (e *SettingEngine) SetSCTPMaxMessageSize(maxMessageSize uint32)
  - func (e *SettingEngine) SetSCTPMaxReceiveBufferSize(maxReceiveBufferSize uint32)
  - func (e *SettingEngine) SetSCTPMinCwnd(minCwnd uint32)
  - func (e *SettingEngine) SetSCTPRTOMax(rtoMax time.Duration)
  - func (e *SettingEngine) SetSDPMediaLevelFingerprints(sdpMediaLevelFingerprints bool)
  - func (e *SettingEngine) SetSRTCPReplayProtectionWindow(n uint)
  - func (e *SettingEngine) SetSRTPProtectionProfiles(profiles ...dtls.SRTPProtectionProfile)
  - func (e *SettingEngine) SetSRTPReplayProtectionWindow(n uint)
  - func (e *SettingEngine) SetSTUNGatherTimeout(t time.Duration)
  - func (e *SettingEngine) SetSrflxAcceptanceMinWait(t time.Duration)
- type SignalingState
-     * func (t *SignalingState) Get() SignalingState
  - func (t *SignalingState) Set(state SignalingState)
  - func (t SignalingState) String() string
- type Stats
-     * func UnmarshalStatsJSON(b []byte) (Stats, error)
- type StatsICECandidatePairState
- type StatsReport
-     * func (r StatsReport) GetCertificateStats(c *Certificate) (CertificateStats, bool)
  - func (r StatsReport) GetCodecStats(c *RTPCodecParameters) (CodecStats, bool)
  - func (r StatsReport) GetConnectionStats(conn *PeerConnection) (PeerConnectionStats, bool)
  - func (r StatsReport) GetDataChannelStats(dc *DataChannel) (DataChannelStats, bool)
  - func (r StatsReport) GetICECandidatePairStats(c *ICECandidatePair) (ICECandidatePairStats, bool)
  - func (r StatsReport) GetICECandidateStats(c *ICECandidate) (ICECandidateStats, bool)
- type StatsTimestamp
-     * func (s StatsTimestamp) Time() time.Time
- type StatsType
- type TrackLocal
- type TrackLocalContext
- type TrackLocalStaticRTP
-     * func NewTrackLocalStaticRTP(c RTPCodecCapability, id, streamID string, ...) (*TrackLocalStaticRTP, error)
-     * func (s *TrackLocalStaticRTP) Bind(trackContext TrackLocalContext) (RTPCodecParameters, error)
  - func (s *TrackLocalStaticRTP) Codec() RTPCodecCapability
  - func (s *TrackLocalStaticRTP) ID() string
  - func (s *TrackLocalStaticRTP) Kind() RTPCodecType
  - func (s *TrackLocalStaticRTP) RID() string
  - func (s *TrackLocalStaticRTP) StreamID() string
  - func (s *TrackLocalStaticRTP) Unbind(t TrackLocalContext) error
  - func (s *TrackLocalStaticRTP) Write(b []byte) (n int, err error)
  - func (s *TrackLocalStaticRTP) WriteRTP(p*rtp.Packet) error
- type TrackLocalStaticSample
-     * func NewTrackLocalStaticSample(c RTPCodecCapability, id, streamID string, ...) (*TrackLocalStaticSample, error)
-     * func (s *TrackLocalStaticSample) Bind(t TrackLocalContext) (RTPCodecParameters, error)
  - func (s *TrackLocalStaticSample) Codec() RTPCodecCapability
  - func (s *TrackLocalStaticSample) GeneratePadding(samples uint32) error
  - func (s *TrackLocalStaticSample) ID() string
  - func (s *TrackLocalStaticSample) Kind() RTPCodecType
  - func (s *TrackLocalStaticSample) RID() string
  - func (s *TrackLocalStaticSample) StreamID() string
  - func (s *TrackLocalStaticSample) Unbind(t TrackLocalContext) error
  - func (s *TrackLocalStaticSample) WriteSample(sample media.Sample) error
- type TrackLocalWriter
- type TrackRemote
-     * func (t *TrackRemote) Codec() RTPCodecParameters
  - func (t *TrackRemote) HasRTX() bool
  - func (t *TrackRemote) ID() string
  - func (t *TrackRemote) Kind() RTPCodecType
  - func (t *TrackRemote) Msid() string
  - func (t *TrackRemote) PayloadType() PayloadType
  - func (t *TrackRemote) RID() string
  - func (t *TrackRemote) Read(b []byte) (n int, attributes interceptor.Attributes, err error)
  - func (t *TrackRemote) ReadRTP() (*rtp.Packet, interceptor.Attributes, error)
  - func (t *TrackRemote) RtxSSRC() SSRC
  - func (t *TrackRemote) SSRC() SSRC
  - func (t *TrackRemote) SetReadDeadline(deadline time.Time) error
  - func (t *TrackRemote) StreamID() string
- type TransportStats
- type VideoReceiverStats
- type VideoSenderStats
- type VideoSourceStats

### Examples ¶

- GatheringCompletePromise
- SettingEngine.SetICEAddressRewriteRules (AppendSrflx)
- SettingEngine.SetICEAddressRewriteRules (ReplaceHost)

### Constants ¶

[View Source](https://github.com/pion/webrtc/blob/v4.2.3/constants.go#L12)

    const (
    
     // AttributeRtxPayloadType is the interceptor attribute added when Read()
     // returns an RTX packet containing the RTX stream payload type.
     AttributeRtxPayloadType = "rtx_payload_type"
     // AttributeRtxSsrc is the interceptor attribute added when Read()
     // returns an RTX packet containing the RTX stream SSRC.
     AttributeRtxSsrc = "rtx_ssrc"
     // AttributeRtxSequenceNumber is the interceptor attribute added when
     // Read() returns an RTX packet containing the RTX stream sequence number.
     AttributeRtxSequenceNumber = "rtx_sequence_number"
    )

[View Source](https://github.com/pion/webrtc/blob/v4.2.3/mimetype.go#L6)

    const (
     // MimeTypeH264 H264 MIME type.
     // Note: Matching should be case insensitive.
     MimeTypeH264 = "video/H264"
     // MimeTypeH265 H265 MIME type
     // Note: Matching should be case insensitive.
     MimeTypeH265 = "video/H265"
     // MimeTypeOpus Opus MIME type
     // Note: Matching should be case insensitive.
     MimeTypeOpus = "audio/opus"
     // MimeTypeVP8 VP8 MIME type
     // Note: Matching should be case insensitive.
     MimeTypeVP8 = "video/VP8"
     // MimeTypeVP9 VP9 MIME type
     // Note: Matching should be case insensitive.
     MimeTypeVP9 = "video/VP9"
     // MimeTypeAV1 AV1 MIME type
     // Note: Matching should be case insensitive.
     MimeTypeAV1 = "video/AV1"
     // MimeTypeG722 G722 MIME type
     // Note: Matching should be case insensitive.
     MimeTypeG722 = "audio/G722"
     // MimeTypePCMU PCMU MIME type
     // Note: Matching should be case insensitive.
     MimeTypePCMU = "audio/PCMU"
     // MimeTypePCMA PCMA MIME type
     // Note: Matching should be case insensitive.
     MimeTypePCMA = "audio/PCMA"
     // MimeTypeRTX RTX MIME type
     // Note: Matching should be case insensitive.
     MimeTypeRTX = "video/rtx"
     // MimeTypeFlexFEC FEC MIME Type
     // Note: Matching should be case insensitive.
     MimeTypeFlexFEC = "video/flexfec"
     // MimeTypeFlexFEC03 FlexFEC03 MIME Type
     // Note: Matching should be case insensitive.
     MimeTypeFlexFEC03 = "video/flexfec-03"
     // MimeTypeUlpFEC UlpFEC MIME Type
     // Note: Matching should be case insensitive.
     MimeTypeUlpFEC = "video/ulpfec"
    )

[View Source](https://github.com/pion/webrtc/blob/v4.2.3/rtcpfeedback.go#L6)

    const (
     // TypeRTCPFBTransportCC ..
     TypeRTCPFBTransportCC = "transport-cc"
    
     // TypeRTCPFBGoogREMB ..
     TypeRTCPFBGoogREMB = "goog-remb"
    
     // TypeRTCPFBACK ..
     TypeRTCPFBACK = "ack"
    
     // TypeRTCPFBCCM ..
     TypeRTCPFBCCM = "ccm"
    
     // TypeRTCPFBNACK ..
     TypeRTCPFBNACK = "nack"
    )

### Variables ¶

[View Source](https://github.com/pion/webrtc/blob/v4.2.3/errors.go#L10)

    var (
     // ErrUnknownType indicates an error with Unknown info.
     ErrUnknownType = [errors](/errors).[New](/errors#New)("unknown")
    
     // ErrConnectionClosed indicates an operation executed after connection
     // has already been closed.
     ErrConnectionClosed = [errors](/errors).[New](/errors#New)("connection closed")
    
     // ErrDataChannelNotOpen indicates an operation executed when the data
     // channel is not (yet) open.
     ErrDataChannelNotOpen = [errors](/errors).[New](/errors#New)("data channel not open")
    
     // ErrCertificateExpired indicates that an x509 certificate has expired.
     ErrCertificateExpired = [errors](/errors).[New](/errors#New)("x509Cert expired")
    
     // ErrNoTurnCredentials indicates that a TURN server URL was provided
     // without required credentials.
     ErrNoTurnCredentials = [errors](/errors).[New](/errors#New)("turn server credentials required")
    
     // ErrTurnCredentials indicates that provided TURN credentials are partial
     // or malformed.
     ErrTurnCredentials = [errors](/errors).[New](/errors#New)("invalid turn server credentials")
    
     // ErrExistingTrack indicates that a track already exists.
     ErrExistingTrack = [errors](/errors).[New](/errors#New)("track already exists")
    
     // ErrPrivateKeyType indicates that a particular private key encryption
     // chosen to generate a certificate is not supported.
     ErrPrivateKeyType = [errors](/errors).[New](/errors#New)("private key type not supported")
    
     // ErrModifyingPeerIdentity indicates that an attempt to modify
     // PeerIdentity was made after PeerConnection has been initialized.
     ErrModifyingPeerIdentity = [errors](/errors).[New](/errors#New)("peerIdentity cannot be modified")
    
     // ErrModifyingCertificates indicates that an attempt to modify
     // Certificates was made after PeerConnection has been initialized.
     ErrModifyingCertificates = [errors](/errors).[New](/errors#New)("certificates cannot be modified")
    
     // ErrModifyingBundlePolicy indicates that an attempt to modify
     // BundlePolicy was made after PeerConnection has been initialized.
     ErrModifyingBundlePolicy = [errors](/errors).[New](/errors#New)("bundle policy cannot be modified")
    
     // ErrModifyingRTCPMuxPolicy indicates that an attempt to modify
     // RTCPMuxPolicy was made after PeerConnection has been initialized.
     ErrModifyingRTCPMuxPolicy = [errors](/errors).[New](/errors#New)("rtcp mux policy cannot be modified")
    
     // ErrModifyingICECandidatePoolSize indicates that an attempt to modify
     // ICECandidatePoolSize was made after PeerConnection has been initialized.
     ErrModifyingICECandidatePoolSize = [errors](/errors).[New](/errors#New)("ice candidate pool size cannot be modified")
    
     // ErrStringSizeLimit indicates that the character size limit of string is
     // exceeded. The limit is hardcoded to 65535 according to specifications.
     ErrStringSizeLimit = [errors](/errors).[New](/errors#New)("data channel label exceeds size limit")
    
     // ErrMaxDataChannelID indicates that the maximum number ID that could be
     // specified for a data channel has been exceeded.
     ErrMaxDataChannelID = [errors](/errors).[New](/errors#New)("maximum number ID for datachannel specified")
    
     // ErrNegotiatedWithoutID indicates that an attempt to create a data channel
     // was made while setting the negotiated option to true without providing
     // the negotiated channel ID.
     ErrNegotiatedWithoutID = [errors](/errors).[New](/errors#New)("negotiated set without channel id")
    
     // ErrRetransmitsOrPacketLifeTime indicates that an attempt to create a data
     // channel was made with both options MaxPacketLifeTime and MaxRetransmits
     // set together. Such configuration is not supported by the specification
     // and is mutually exclusive.
     ErrRetransmitsOrPacketLifeTime = [errors](/errors).[New](/errors#New)("both MaxPacketLifeTime and MaxRetransmits was set")
    
     // ErrCodecNotFound is returned when a codec search to the Media Engine fails.
     ErrCodecNotFound = [errors](/errors).[New](/errors#New)("codec not found")
    
     // ErrNoRemoteDescription indicates that an operation was rejected because
     // the remote description is not set.
     ErrNoRemoteDescription = [errors](/errors).[New](/errors#New)("remote description is not set")
    
     // ErrIncorrectSDPSemantics indicates that the PeerConnection was configured to
     // generate SDP Answers with different SDP Semantics than the received Offer.
     ErrIncorrectSDPSemantics = [errors](/errors).[New](/errors#New)("remote SessionDescription semantics does not match configuration")
    
     // ErrIncorrectSignalingState indicates that the signaling state of PeerConnection is not correct.
     ErrIncorrectSignalingState = [errors](/errors).[New](/errors#New)("operation can not be run in current signaling state")
    
     // ErrProtocolTooLarge indicates that value given for a DataChannelInit protocol is
     // longer then 65535 bytes.
     ErrProtocolTooLarge = [errors](/errors).[New](/errors#New)("protocol is larger then 65535 bytes")
    
     // ErrSenderNotCreatedByConnection indicates RemoveTrack was called with a RtpSender not created
     // by this PeerConnection.
     ErrSenderNotCreatedByConnection = [errors](/errors).[New](/errors#New)("RtpSender not created by this PeerConnection")
    
     // ErrSessionDescriptionNoFingerprint indicates SetRemoteDescription was called with a SessionDescription that has no
     // fingerprint.
     ErrSessionDescriptionNoFingerprint = [errors](/errors).[New](/errors#New)("SetRemoteDescription called with no fingerprint")
    
     // ErrSessionDescriptionInvalidFingerprint indicates SetRemoteDescription was called with a SessionDescription that
     // has an invalid fingerprint.
     ErrSessionDescriptionInvalidFingerprint = [errors](/errors).[New](/errors#New)("SetRemoteDescription called with an invalid fingerprint")
    
     // ErrSessionDescriptionConflictingFingerprints indicates SetRemoteDescription was called with a SessionDescription
     // that has an conflicting fingerprints.
     ErrSessionDescriptionConflictingFingerprints = [errors](/errors).[New](/errors#New)(
      "SetRemoteDescription called with multiple conflicting fingerprint",
     )
    
     // ErrSessionDescriptionMissingIceUfrag indicates SetRemoteDescription was called with a SessionDescription that
     // is missing an ice-ufrag value.
     ErrSessionDescriptionMissingIceUfrag = [errors](/errors).[New](/errors#New)("SetRemoteDescription called with no ice-ufrag")
    
     // ErrSessionDescriptionMissingIcePwd indicates SetRemoteDescription was called with a SessionDescription that
     // is missing an ice-pwd value.
     ErrSessionDescriptionMissingIcePwd = [errors](/errors).[New](/errors#New)("SetRemoteDescription called with no ice-pwd")
    
     // ErrSessionDescriptionConflictingIceUfrag  indicates SetRemoteDescription was called with a SessionDescription
     // that contains multiple conflicting ice-ufrag values.
     ErrSessionDescriptionConflictingIceUfrag = [errors](/errors).[New](/errors#New)(
      "SetRemoteDescription called with multiple conflicting ice-ufrag values",
     )
    
     // ErrSessionDescriptionConflictingIcePwd indicates SetRemoteDescription was called with a SessionDescription
     // that contains multiple conflicting ice-pwd values.
     ErrSessionDescriptionConflictingIcePwd = [errors](/errors).[New](/errors#New)(
      "SetRemoteDescription called with multiple conflicting ice-pwd values",
     )
    
     // ErrNoSRTPProtectionProfile indicates that the DTLS handshake completed and no SRTP Protection Profile was chosen.
     ErrNoSRTPProtectionProfile = [errors](/errors).[New](/errors#New)("DTLS Handshake completed and no SRTP Protection Profile was chosen")
    
     // ErrFailedToGenerateCertificateFingerprint indicates that we failed to generate the fingerprint
     // used for comparing certificates.
     ErrFailedToGenerateCertificateFingerprint = [errors](/errors).[New](/errors#New)("failed to generate certificate fingerprint")
    
     // ErrNoCodecsAvailable indicates that operation isn't possible because the MediaEngine has no codecs available.
     ErrNoCodecsAvailable = [errors](/errors).[New](/errors#New)("operation failed no codecs are available")
    
     // ErrUnsupportedCodec indicates the remote peer doesn't support the requested codec.
     ErrUnsupportedCodec = [errors](/errors).[New](/errors#New)("unable to start track, codec is not supported by remote")
    
     // ErrSenderWithNoCodecs indicates that a RTPSender was created without any codecs. To send media the MediaEngine
     //  needs at least one configured codec.
     ErrSenderWithNoCodecs = [errors](/errors).[New](/errors#New)("unable to populate media section, RTPSender created with no codecs")
    
     // ErrCodecAlreadyRegistered indicates that a codec has already been registered for the same payload type.
     ErrCodecAlreadyRegistered = [errors](/errors).[New](/errors#New)("codec already registered for same payload type")
    
     // ErrRTPSenderNewTrackHasIncorrectKind indicates that the new track is of a different kind than the previous/original.
     ErrRTPSenderNewTrackHasIncorrectKind = [errors](/errors).[New](/errors#New)("new track must be of the same kind as previous")
    
     // ErrRTPSenderNewTrackHasIncorrectEnvelope indicates that the new track has a different envelope
     //  than the previous/original.
     ErrRTPSenderNewTrackHasIncorrectEnvelope = [errors](/errors).[New](/errors#New)("new track must have the same envelope as previous")
    
     // ErrUnbindFailed indicates that a TrackLocal was not able to be unbind.
     ErrUnbindFailed = [errors](/errors).[New](/errors#New)("failed to unbind TrackLocal from PeerConnection")
    
     // ErrNoPayloaderForCodec indicates that the requested codec does not have a payloader.
     ErrNoPayloaderForCodec = [errors](/errors).[New](/errors#New)("the requested codec does not have a payloader")
    
     // ErrRegisterHeaderExtensionInvalidDirection indicates that a extension was
     // registered with a direction besides `sendonly` or `recvonly`.
     ErrRegisterHeaderExtensionInvalidDirection = [errors](/errors).[New](/errors#New)(
      "a header extension must be registered as 'recvonly', 'sendonly' or both",
     )
    
     // ErrSimulcastProbeOverflow indicates that too many Simulcast probe streams are in flight
     // and the requested SSRC was ignored.
     ErrSimulcastProbeOverflow = [errors](/errors).[New](/errors#New)("simulcast probe limit has been reached, new SSRC has been discarded")
    
     // ErrSDPUnmarshalling indicates that the SDP could not be unmarshalled.
     ErrSDPUnmarshalling = [errors](/errors).[New](/errors#New)("failed to unmarshal SDP")
    )

### Functions ¶

#### func [ConfigureCongestionControlFeedback](https://github.com/pion/webrtc/blob/v4.2.3/interceptor.go#L172) ¶

    func ConfigureCongestionControlFeedback(mediaEngine *MediaEngine, interceptorRegistry *[interceptor](/github.com/pion/interceptor).[Registry](/github.com/pion/interceptor#Registry)) [error](/builtin#error)

ConfigureCongestionControlFeedback registers congestion control feedback as defined in [RFC 8888](https://rfc-editor.org/rfc/rfc8888.html) (<https://datatracker.ietf.org/doc/rfc8888/>)

#### func [ConfigureFlexFEC03](https://github.com/pion/webrtc/blob/v4.2.3/interceptor.go#L208) ¶ added in v4.1.2

    func ConfigureFlexFEC03(
     payloadType PayloadType,
     mediaEngine *MediaEngine,
     interceptorRegistry *[interceptor](/github.com/pion/interceptor).[Registry](/github.com/pion/interceptor#Registry),
     options ...[flexfec](/github.com/pion/interceptor/pkg/flexfec).[FecOption](/github.com/pion/interceptor/pkg/flexfec#FecOption),
    ) [error](/builtin#error)

ConfigureFlexFEC03 registers flexfec-03 codec with provided payloadType in mediaEngine and adds corresponding interceptor to the registry. Note that this function should be called before any other interceptor that modifies RTP packets (i.e. TWCCHeaderExtensionSender) is added to the registry, so that packets generated by flexfec interceptor are not modified.

#### func [ConfigureNack](https://github.com/pion/webrtc/blob/v4.2.3/interceptor.go#L99) ¶

    func ConfigureNack(mediaEngine *MediaEngine, interceptorRegistry *[interceptor](/github.com/pion/interceptor).[Registry](/github.com/pion/interceptor#Registry)) [error](/builtin#error)

ConfigureNack will setup everything necessary for handling generating/responding to nack messages.

#### func [ConfigureRTCPReports](https://github.com/pion/webrtc/blob/v4.2.3/interceptor.go#L81) ¶

    func ConfigureRTCPReports(interceptorRegistry *[interceptor](/github.com/pion/interceptor).[Registry](/github.com/pion/interceptor#Registry)) [error](/builtin#error)

ConfigureRTCPReports will setup everything necessary for generating Sender and Receiver Reports.

#### func [ConfigureSimulcastExtensionHeaders](https://github.com/pion/webrtc/blob/v4.2.3/interceptor.go#L185) ¶

    func ConfigureSimulcastExtensionHeaders(mediaEngine *MediaEngine) [error](/builtin#error)

ConfigureSimulcastExtensionHeaders enables the RTP Extension Headers needed for Simulcast.

#### func [ConfigureStatsInterceptor](https://github.com/pion/webrtc/blob/v4.2.3/interceptor.go#L48) ¶ added in v4.1.5

    func ConfigureStatsInterceptor(interceptorRegistry *[interceptor](/github.com/pion/interceptor).[Registry](/github.com/pion/interceptor#Registry)) [error](/builtin#error)

ConfigureStatsInterceptor will setup everything necessary for generating RTP stream statistics.

#### func [ConfigureTWCCHeaderExtensionSender](https://github.com/pion/webrtc/blob/v4.2.3/interceptor.go#L120) ¶

    func ConfigureTWCCHeaderExtensionSender(mediaEngine *MediaEngine, interceptorRegistry *[interceptor](/github.com/pion/interceptor).[Registry](/github.com/pion/interceptor#Registry)) [error](/builtin#error)

ConfigureTWCCHeaderExtensionSender will setup everything necessary for adding a TWCC header extension to outgoing RTP packets. This will allow the remote peer to generate TWCC reports.

#### func [ConfigureTWCCSender](https://github.com/pion/webrtc/blob/v4.2.3/interceptor.go#L145) ¶

    func ConfigureTWCCSender(mediaEngine *MediaEngine, interceptorRegistry *[interceptor](/github.com/pion/interceptor).[Registry](/github.com/pion/interceptor#Registry)) [error](/builtin#error)

ConfigureTWCCSender will setup everything necessary for generating TWCC reports. This must be called after registering codecs with the MediaEngine.

#### func [GatheringCompletePromise](https://github.com/pion/webrtc/blob/v4.2.3/gathering_complete_promise.go#L17) ¶

    func GatheringCompletePromise(pc *PeerConnection) (gatherComplete <-chan struct{})

GatheringCompletePromise is a Pion specific helper function that returns a channel that is closed when gathering is complete. This function may be helpful in cases where you are unable to trickle your ICE Candidates.

It is better to not use this function, and instead trickle candidates. If you use this function you will see longer connection startup times. When the call is connected you will see no impact however.

Example ¶

ExampleGatheringCompletePromise demonstrates how to implement non-trickle ICE in Pion, an older form of ICE that does not require an asynchronous side channel between peers: negotiation is just a single offer-answer exchange. It works by explicitly waiting for all local ICE candidates to have been gathered before sending an offer to the peer.

    // create a peer connection
    pc, err := NewPeerConnection(Configuration{})
    if err != nil {
     panic(err)
    }
    defer func() {
     closeErr := pc.Close()
     if closeErr != nil {
      panic(closeErr)
     }
    }()
    
    // add at least one transceiver to the peer connection, or nothing
    // interesting will happen.  This could use pc.AddTrack instead.
    _, err = pc.AddTransceiverFromKind(RTPCodecTypeVideo)
    if err != nil {
     panic(err)
    }
    
    // create a first offer that does not contain any local candidates
    offer, err := pc.CreateOffer(nil)
    if err != nil {
     panic(err)
    }
    
    // gatherComplete is a channel that will be closed when
    // the gathering of local candidates is complete.
    gatherComplete := GatheringCompletePromise(pc)
    
    // apply the offer
    err = pc.SetLocalDescription(offer)
    if err != nil {
     panic(err)
    }
    
    // wait for gathering of local candidates to complete
    <-gatherComplete
    
    // compute the local offer again
    offer2 := pc.LocalDescription()
    
    // this second offer contains all candidates, and may be sent to
    // the peer with no need for further communication.  In this
    // example, we simply check that it contains at least one
    // candidate.
    hasCandidate := strings.Contains(offer2.SDP, "\na=candidate:")
    if hasCandidate {
     fmt.Println("Ok!")
    }
    
    
    
    Output:
    
    Ok!
    

#### func [NewAudioPlayoutStatsProvider](https://github.com/pion/webrtc/blob/v4.2.3/stats_go.go#L137) ¶ added in v4.1.5

    func NewAudioPlayoutStatsProvider(id [string](/builtin#string)) *defaultAudioPlayoutStatsProvider

NewAudioPlayoutStatsProvider constructs a default provider with the supplied stats ID.

#### func [NewICETCPMux](https://github.com/pion/webrtc/blob/v4.2.3/icemux.go#L15) ¶

    func NewICETCPMux(logger [logging](/github.com/pion/logging).[LeveledLogger](/github.com/pion/logging#LeveledLogger), listener [net](/net).[Listener](/net#Listener), readBufferSize [int](/builtin#int)) [ice](/github.com/pion/ice/v4).[TCPMux](/github.com/pion/ice/v4#TCPMux)

NewICETCPMux creates a new instance of ice.TCPMuxDefault. It enables use of passive ICE TCP candidates.

#### func [NewICEUDPMux](https://github.com/pion/webrtc/blob/v4.2.3/icemux.go#L25) ¶

    func NewICEUDPMux(logger [logging](/github.com/pion/logging).[LeveledLogger](/github.com/pion/logging#LeveledLogger), udpConn [net](/net).[PacketConn](/net#PacketConn)) [ice](/github.com/pion/ice/v4).[UDPMux](/github.com/pion/ice/v4#UDPMux)

NewICEUDPMux creates a new instance of ice.UDPMuxDefault. It allows many PeerConnections to be served by a single UDP Port.

#### func [RegisterDefaultInterceptors](https://github.com/pion/webrtc/blob/v4.2.3/interceptor.go#L27) ¶

    func RegisterDefaultInterceptors(mediaEngine *MediaEngine, interceptorRegistry *[interceptor](/github.com/pion/interceptor).[Registry](/github.com/pion/interceptor#Registry)) [error](/builtin#error)

RegisterDefaultInterceptors will register some useful interceptors. If you want to customize which interceptors are loaded, you should copy the code from this method and remove unwanted interceptors.

#### func [WithInterceptorRegistry](https://github.com/pion/webrtc/blob/v4.2.3/api.go#L89) ¶

    func WithInterceptorRegistry(ir *[interceptor](/github.com/pion/interceptor).[Registry](/github.com/pion/interceptor#Registry)) func(a *API)

WithInterceptorRegistry allows providing Interceptors to the API. Settings should not be changed after passing the registry to an API.

#### func [WithMediaEngine](https://github.com/pion/webrtc/blob/v4.2.3/api.go#L70) ¶

    func WithMediaEngine(m *MediaEngine) func(a *API)

WithMediaEngine allows providing a MediaEngine to the API. Settings can be changed after passing the engine to an API. When a PeerConnection is created the MediaEngine is copied and no more changes can be made.

#### func [WithPayloader](https://github.com/pion/webrtc/blob/v4.2.3/track_local_static.go#L68) ¶ added in v4.0.2

    func WithPayloader(h func(RTPCodecCapability) ([rtp](/github.com/pion/rtp).[Payloader](/github.com/pion/rtp#Payloader), [error](/builtin#error))) func(*TrackLocalStaticRTP)

WithPayloader allows the user to override the Payloader.

#### func [WithRTPSequenceNumber](https://github.com/pion/webrtc/blob/v4.2.3/track_local_static.go#L82) ¶ added in v4.2.2

    func WithRTPSequenceNumber(sequenceNumber [uint16](/builtin#uint16)) func(*TrackLocalStaticRTP)

WithRTPSequenceNumber sets the initial RTP sequence number for the track.

#### func [WithRTPStreamID](https://github.com/pion/webrtc/blob/v4.2.3/track_local_static.go#L61) ¶

    func WithRTPStreamID(rid [string](/builtin#string)) func(*TrackLocalStaticRTP)

WithRTPStreamID sets the RTP stream ID for this TrackLocalStaticRTP.

#### func [WithRTPTimestamp](https://github.com/pion/webrtc/blob/v4.2.3/track_local_static.go#L75) ¶ added in v4.1.0

    func WithRTPTimestamp(timestamp [uint32](/builtin#uint32)) func(*TrackLocalStaticRTP)

WithRTPTimestamp set the initial RTP timestamp for the track.

#### func [WithSettingEngine](https://github.com/pion/webrtc/blob/v4.2.3/api.go#L81) ¶

    func WithSettingEngine(s SettingEngine) func(a *API)

WithSettingEngine allows providing a SettingEngine to the API. Settings should not be changed after passing the engine to an API.

### Types ¶

#### type [API](https://github.com/pion/webrtc/blob/v4.2.3/api.go#L19) ¶

    type API struct {
     // contains filtered or unexported fields
    }

API allows configuration of a PeerConnection with APIs that are available in the standard. This lets you set custom behavior via the SettingEngine, configure codecs via the MediaEngine and define custom media behaviors via Interceptors.

#### func [NewAPI](https://github.com/pion/webrtc/blob/v4.2.3/api.go#L31) ¶

    func NewAPI(options ...func(*API)) *API

NewAPI Creates a new API object for keeping semi-global settings to WebRTC objects

It uses the default Codecs and Interceptors unless you customize them using WithMediaEngine and WithInterceptorRegistry respectively.

#### func (*API) [NewDTLSTransport](https://github.com/pion/webrtc/blob/v4.2.3/dtlstransport.go#L78) ¶

    func (api *API) NewDTLSTransport(transport *ICETransport, certificates []Certificate) (*DTLSTransport, [error](/builtin#error))

NewDTLSTransport creates a new DTLSTransport. This constructor is part of the ORTC API. It is not meant to be used together with the basic WebRTC API.

#### func (*API) [NewDataChannel](https://github.com/pion/webrtc/blob/v4.2.3/datachannel.go#L72) ¶

    func (api *API) NewDataChannel(transport *SCTPTransport, params *DataChannelParameters) (*DataChannel, [error](/builtin#error))

NewDataChannel creates a new DataChannel. This constructor is part of the ORTC API. It is not meant to be used together with the basic WebRTC API.

#### func (*API) [NewICEGatherer](https://github.com/pion/webrtc/blob/v4.2.3/icegatherer.go#L94) ¶

    func (api *API) NewICEGatherer(opts ICEGatherOptions) (*ICEGatherer, [error](/builtin#error))

NewICEGatherer creates a new NewICEGatherer. This constructor is part of the ORTC API. It is not meant to be used together with the basic WebRTC API.

#### func (*API) [NewICETransport](https://github.com/pion/webrtc/blob/v4.2.3/ice_go.go#L12) ¶

    func (api *API) NewICETransport(gatherer *ICEGatherer) *ICETransport

NewICETransport creates a new NewICETransport. This constructor is part of the ORTC API. It is not meant to be used together with the basic WebRTC API.

#### func (*API) [NewPeerConnection](https://github.com/pion/webrtc/blob/v4.2.3/peerconnection.go#L116) ¶

    func (api *API) NewPeerConnection(configuration Configuration) (*PeerConnection, [error](/builtin#error))

NewPeerConnection creates a new PeerConnection with the provided configuration against the received API object. This method will attach a default set of codecs and interceptors to the resulting PeerConnection. If this behavior is not desired, set the set of codecs and interceptors explicitly by using WithMediaEngine and WithInterceptorRegistry when calling NewAPI.

#### func (*API) [NewRTPReceiver](https://github.com/pion/webrtc/blob/v4.2.3/rtpreceiver.go#L83) ¶

    func (api *API) NewRTPReceiver(kind RTPCodecType, transport *DTLSTransport) (*RTPReceiver, [error](/builtin#error))

NewRTPReceiver constructs a new RTPReceiver.

#### func (*API) [NewRTPSender](https://github.com/pion/webrtc/blob/v4.2.3/rtpsender.go#L61) ¶

    func (api *API) NewRTPSender(track TrackLocal, transport *DTLSTransport) (*RTPSender, [error](/builtin#error))

NewRTPSender constructs a new RTPSender.

#### func (*API) [NewSCTPTransport](https://github.com/pion/webrtc/blob/v4.2.3/sctptransport.go#L63) ¶

    func (api *API) NewSCTPTransport(dtls *DTLSTransport) *SCTPTransport

NewSCTPTransport creates a new SCTPTransport. This constructor is part of the ORTC API. It is not meant to be used together with the basic WebRTC API.

#### type [AnswerOptions](https://github.com/pion/webrtc/blob/v4.2.3/offeransweroptions.go#L20) ¶

    type AnswerOptions struct {
     OfferAnswerOptions
    }

AnswerOptions structure describes the options used to control the answer creation process.

#### type [AudioPlayoutStats](https://github.com/pion/webrtc/blob/v4.2.3/stats.go#L1311) ¶

    type AudioPlayoutStats struct {
     // Timestamp is the timestamp associated with this object.
     Timestamp StatsTimestamp `json:"timestamp"`
    
     // Type is the object's StatsType
     Type StatsType `json:"type"`
    
     // ID is a unique id that is associated with the component inspected to produce
     // this Stats object. Two Stats objects will have the same ID if they were produced
     // by inspecting the same underlying object.
     ID [string](/builtin#string) `json:"id"`
    
     // Kind is "audio"
     Kind [string](/builtin#string) `json:"kind"`
    
     // SynthesizedSamplesDuration is measured in seconds and is incremented each time an audio sample is synthesized by
     // this playout path. This metric can be used together with totalSamplesDuration to calculate the percentage of played
     // out media being synthesized. If the playout path is unable to produce audio samples on time for device playout,
     // samples are synthesized to be played out instead. Synthesization typically only happens if the pipeline is
     // underperforming. Samples synthesized by the RTCInboundRtpStreamStats are not counted for here, but in
     // InboundRtpStreamStats.concealedSamples.
     SynthesizedSamplesDuration [float64](/builtin#float64) `json:"synthesizedSamplesDuration"`
    
     // SynthesizedSamplesEvents is the number of synthesized samples events. This counter increases every time a sample
     // is synthesized after a non-synthesized sample. That is, multiple consecutive synthesized samples will increase
     // synthesizedSamplesDuration multiple times but is a single synthesization samples event.
     SynthesizedSamplesEvents [uint64](/builtin#uint64) `json:"synthesizedSamplesEvents"`
    
     // TotalSamplesDuration represents the total duration in seconds of all samples
     // that have sent or received (and thus counted by TotalSamplesSent or TotalSamplesReceived).
     // Can be used with TotalAudioEnergy to compute an average audio level over different intervals.
     TotalSamplesDuration [float64](/builtin#float64) `json:"totalSamplesDuration"`
    
     // When audio samples are pulled by the playout device, this counter is incremented with the estimated delay of the
     // playout path for that audio sample. The playout delay includes the delay from being emitted to the actual time of
     // playout on the device. This metric can be used together with totalSamplesCount to calculate the average
     // playout delay per sample.
     TotalPlayoutDelay [float64](/builtin#float64) `json:"totalPlayoutDelay"`
    
     // When audio samples are pulled by the playout device, this counter is incremented with the number of samples
     // emitted for playout.
     TotalSamplesCount [uint64](/builtin#uint64) `json:"totalSamplesCount"`
    }

AudioPlayoutStats represents one playout path - if the same playout stats object is referenced by multiple RTCInboundRtpStreamStats this is an indication that audio mixing is happening in which case sample counters in this stats object refer to the samples after mixing. Only applicable if the playout path represents an audio device.

#### type [AudioPlayoutStatsProvider](https://github.com/pion/webrtc/blob/v4.2.3/stats_go.go#L112) ¶ added in v4.1.5

    type AudioPlayoutStatsProvider interface {
     // AddTrack registers a track to report playout stats to this provider.
     AddTrack(track *TrackRemote) [error](/builtin#error)
    
     // RemoveTrack unregisters a track from this provider.
     RemoveTrack(track *TrackRemote)
    
     // Snapshot returns the accumulated stats at the given time.
     Snapshot(now [time](/time).[Time](/time#Time)) (AudioPlayoutStats, [bool](/builtin#bool))
    }

AudioPlayoutStatsProvider is an interface for getting audio playout metrics.

#### type [AudioReceiverStats](https://github.com/pion/webrtc/blob/v4.2.3/stats.go#L1728) ¶

    type AudioReceiverStats struct {
     // Timestamp is the timestamp associated with this object.
     Timestamp StatsTimestamp `json:"timestamp"`
    
     // Type is the object's StatsType
     Type StatsType `json:"type"`
    
     // ID is a unique id that is associated with the component inspected to produce
     // this Stats object. Two Stats objects will have the same ID if they were produced
     // by inspecting the same underlying object.
     ID [string](/builtin#string) `json:"id"`
    
     // Kind is "audio"
     Kind [string](/builtin#string) `json:"kind"`
    
     // AudioLevel represents the output audio level of the track.
     //
     // The value is a value between 0..1 (linear), where 1.0 represents 0 dBov,
     // 0 represents silence, and 0.5 represents approximately 6 dBSPL change in
     // the sound pressure level from 0 dBov.
     //
     // If the track is sourced from a Receiver, does no audio processing, has a
     // constant level, and has a volume setting of 1.0, the audio level is expected
     // to be the same as the audio level of the source SSRC, while if the volume setting
     // is 0.5, the AudioLevel is expected to be half that value.
     //
     // For outgoing audio tracks, the AudioLevel is the level of the audio being sent.
     AudioLevel [float64](/builtin#float64) `json:"audioLevel"`
    
     // TotalAudioEnergy is the total energy of all the audio samples sent/received
     // for this object, calculated by duration * Math.pow(energy/maxEnergy, 2) for
     // each audio sample seen.
     TotalAudioEnergy [float64](/builtin#float64) `json:"totalAudioEnergy"`
    
     // VoiceActivityFlag represents whether the last RTP packet sent or played out
     // by this track contained voice activity or not based on the presence of the
     // V bit in the extension header, as defined in [RFC6464].
     //
     // This value indicates the voice activity in the latest RTP packet played out
     // from a given SSRC, and is defined in RTPSynchronizationSource.voiceActivityFlag.
     VoiceActivityFlag [bool](/builtin#bool) `json:"voiceActivityFlag"`
    
     // TotalSamplesDuration represents the total duration in seconds of all samples
     // that have sent or received (and thus counted by TotalSamplesSent or TotalSamplesReceived).
     // Can be used with TotalAudioEnergy to compute an average audio level over different intervals.
     TotalSamplesDuration [float64](/builtin#float64) `json:"totalSamplesDuration"`
    
     // EstimatedPlayoutTimestamp is the estimated playout time of this receiver's
     // track. The playout time is the NTP timestamp of the last playable sample that
     // has a known timestamp (from an RTCP SR packet mapping RTP timestamps to NTP
     // timestamps), extrapolated with the time elapsed since it was ready to be played out.
     // This is the "current time" of the track in NTP clock time of the sender and
     // can be present even if there is no audio currently playing.
     //
     // This can be useful for estimating how much audio and video is out of
     // sync for two tracks from the same source:
     //   AudioTrackStats.EstimatedPlayoutTimestamp - VideoTrackStats.EstimatedPlayoutTimestamp
     EstimatedPlayoutTimestamp StatsTimestamp `json:"estimatedPlayoutTimestamp"`
    
     // JitterBufferDelay is the sum of the time, in seconds, each sample takes from
     // the time it is received and to the time it exits the jitter buffer.
     // This increases upon samples exiting, having completed their time in the buffer
     // (incrementing JitterBufferEmittedCount). The average jitter buffer delay can
     // be calculated by dividing the JitterBufferDelay with the JitterBufferEmittedCount.
     JitterBufferDelay [float64](/builtin#float64) `json:"jitterBufferDelay"`
    
     // JitterBufferEmittedCount is the total number of samples that have come out
     // of the jitter buffer (increasing JitterBufferDelay).
     JitterBufferEmittedCount [uint64](/builtin#uint64) `json:"jitterBufferEmittedCount"`
    
     // TotalSamplesReceived is the total number of samples that have been received
     // by this receiver. This includes ConcealedSamples.
     TotalSamplesReceived [uint64](/builtin#uint64) `json:"totalSamplesReceived"`
    
     // ConcealedSamples is the total number of samples that are concealed samples.
     // A concealed sample is a sample that is based on data that was synthesized
     // to conceal packet loss and does not represent incoming data.
     ConcealedSamples [uint64](/builtin#uint64) `json:"concealedSamples"`
    
     // ConcealmentEvents is the number of concealment events. This counter increases
     // every time a concealed sample is synthesized after a non-concealed sample.
     // That is, multiple consecutive concealed samples will increase the concealedSamples
     // count multiple times but is a single concealment event.
     ConcealmentEvents [uint64](/builtin#uint64) `json:"concealmentEvents"`
    }

AudioReceiverStats contains audio metrics related to a specific receiver.

#### type [AudioSenderStats](https://github.com/pion/webrtc/blob/v4.2.3/stats.go#L1506) ¶

    type AudioSenderStats struct {
     // Timestamp is the timestamp associated with this object.
     Timestamp StatsTimestamp `json:"timestamp"`
    
     // Type is the object's StatsType
     Type StatsType `json:"type"`
    
     // ID is a unique id that is associated with the component inspected to produce
     // this Stats object. Two Stats objects will have the same ID if they were produced
     // by inspecting the same underlying object.
     ID [string](/builtin#string) `json:"id"`
    
     // TrackIdentifier represents the id property of the track.
     TrackIdentifier [string](/builtin#string) `json:"trackIdentifier"`
    
     // RemoteSource is true if the source is remote, for instance if it is sourced
     // from another host via a PeerConnection. False otherwise. Only applicable for 'track' stats.
     RemoteSource [bool](/builtin#bool) `json:"remoteSource"`
    
     // Ended reflects the "ended" state of the track.
     Ended [bool](/builtin#bool) `json:"ended"`
    
     // Kind is "audio"
     Kind [string](/builtin#string) `json:"kind"`
    
     // AudioLevel represents the output audio level of the track.
     //
     // The value is a value between 0..1 (linear), where 1.0 represents 0 dBov,
     // 0 represents silence, and 0.5 represents approximately 6 dBSPL change in
     // the sound pressure level from 0 dBov.
     //
     // If the track is sourced from an Receiver, does no audio processing, has a
     // constant level, and has a volume setting of 1.0, the audio level is expected
     // to be the same as the audio level of the source SSRC, while if the volume setting
     // is 0.5, the AudioLevel is expected to be half that value.
     //
     // For outgoing audio tracks, the AudioLevel is the level of the audio being sent.
     AudioLevel [float64](/builtin#float64) `json:"audioLevel"`
    
     // TotalAudioEnergy is the total energy of all the audio samples sent/received
     // for this object, calculated by duration * Math.pow(energy/maxEnergy, 2) for
     // each audio sample seen.
     TotalAudioEnergy [float64](/builtin#float64) `json:"totalAudioEnergy"`
    
     // VoiceActivityFlag represents whether the last RTP packet sent or played out
     // by this track contained voice activity or not based on the presence of the
     // V bit in the extension header, as defined in [RFC6464].
     //
     // This value indicates the voice activity in the latest RTP packet played out
     // from a given SSRC, and is defined in RTPSynchronizationSource.voiceActivityFlag.
     VoiceActivityFlag [bool](/builtin#bool) `json:"voiceActivityFlag"`
    
     // TotalSamplesDuration represents the total duration in seconds of all samples
     // that have sent or received (and thus counted by TotalSamplesSent or TotalSamplesReceived).
     // Can be used with TotalAudioEnergy to compute an average audio level over different intervals.
     TotalSamplesDuration [float64](/builtin#float64) `json:"totalSamplesDuration"`
    
     // EchoReturnLoss is only present while the sender is sending a track sourced from
     // a microphone where echo cancellation is applied. Calculated in decibels.
     EchoReturnLoss [float64](/builtin#float64) `json:"echoReturnLoss"`
    
     // EchoReturnLossEnhancement is only present while the sender is sending a track
     // sourced from a microphone where echo cancellation is applied. Calculated in decibels.
     EchoReturnLossEnhancement [float64](/builtin#float64) `json:"echoReturnLossEnhancement"`
    
     // TotalSamplesSent is the total number of samples that have been sent by this sender.
     TotalSamplesSent [uint64](/builtin#uint64) `json:"totalSamplesSent"`
    }

AudioSenderStats represents the stats about one audio sender of a PeerConnection object for which one calls GetStats.

It appears in the stats as soon as the RTPSender is added by either AddTrack or AddTransceiver, or by media negotiation.

#### type [AudioSourceStats](https://github.com/pion/webrtc/blob/v4.2.3/stats.go#L1170) ¶

    type AudioSourceStats struct {
     // Timestamp is the timestamp associated with this object.
     Timestamp StatsTimestamp `json:"timestamp"`
    
     // Type is the object's StatsType
     Type StatsType `json:"type"`
    
     // ID is a unique id that is associated with the component inspected to produce
     // this Stats object. Two Stats objects will have the same ID if they were produced
     // by inspecting the same underlying object.
     ID [string](/builtin#string) `json:"id"`
    
     // TrackIdentifier represents the id property of the track.
     TrackIdentifier [string](/builtin#string) `json:"trackIdentifier"`
    
     // Kind is "audio"
     Kind [string](/builtin#string) `json:"kind"`
    
     // AudioLevel represents the output audio level of the track.
     //
     // The value is a value between 0..1 (linear), where 1.0 represents 0 dBov,
     // 0 represents silence, and 0.5 represents approximately 6 dBSPL change in
     // the sound pressure level from 0 dBov.
     //
     // If the track is sourced from an Receiver, does no audio processing, has a
     // constant level, and has a volume setting of 1.0, the audio level is expected
     // to be the same as the audio level of the source SSRC, while if the volume setting
     // is 0.5, the AudioLevel is expected to be half that value.
     AudioLevel [float64](/builtin#float64) `json:"audioLevel"`
    
     // TotalAudioEnergy is the total energy of all the audio samples sent/received
     // for this object, calculated by duration * Math.pow(energy/maxEnergy, 2) for
     // each audio sample seen.
     TotalAudioEnergy [float64](/builtin#float64) `json:"totalAudioEnergy"`
    
     // TotalSamplesDuration represents the total duration in seconds of all samples
     // that have sent or received (and thus counted by TotalSamplesSent or TotalSamplesReceived).
     // Can be used with TotalAudioEnergy to compute an average audio level over different intervals.
     TotalSamplesDuration [float64](/builtin#float64) `json:"totalSamplesDuration"`
    
     // EchoReturnLoss is only present while the sender is sending a track sourced from
     // a microphone where echo cancellation is applied. Calculated in decibels.
     EchoReturnLoss [float64](/builtin#float64) `json:"echoReturnLoss"`
    
     // EchoReturnLossEnhancement is only present while the sender is sending a track
     // sourced from a microphone where echo cancellation is applied. Calculated in decibels.
     EchoReturnLossEnhancement [float64](/builtin#float64) `json:"echoReturnLossEnhancement"`
    
     // DroppedSamplesDuration represents the total duration, in seconds, of samples produced by the device that got
     // dropped before reaching the media source. Only applicable if this media source is backed by an audio capture device.
     DroppedSamplesDuration [float64](/builtin#float64) `json:"droppedSamplesDuration"`
    
     // DroppedSamplesEvents is the number of dropped samples events. This counter increases every time a sample is
     // dropped after a non-dropped sample. That is, multiple consecutive dropped samples will increase
     // droppedSamplesDuration multiple times but is a single dropped samples event.
     DroppedSamplesEvents [uint64](/builtin#uint64) `json:"droppedSamplesEvents"`
    
     // TotalCaptureDelay is the total delay, in seconds, for each audio sample between the time the sample was emitted
     // by the capture device and the sample reaching the source. This can be used together with totalSamplesCaptured to
     // calculate the average capture delay per sample.
     // Only applicable if the audio source represents an audio capture device.
     TotalCaptureDelay [float64](/builtin#float64) `json:"totalCaptureDelay"`
    
     // TotalSamplesCaptured is the total number of captured samples reaching the audio source, i.e. that were not dropped
     // by the capture pipeline. The frequency of the media source is not necessarily the same as the frequency of encoders
     // later in the pipeline. Only applicable if the audio source represents an audio capture device.
     TotalSamplesCaptured [uint64](/builtin#uint64) `json:"totalSamplesCaptured"`
    }

AudioSourceStats represents an audio track that is attached to one or more senders.

#### type [BundlePolicy](https://github.com/pion/webrtc/blob/v4.2.3/bundlepolicy.go#L14) ¶

    type BundlePolicy [int](/builtin#int)

BundlePolicy affects which media tracks are negotiated if the remote endpoint is not bundle-aware, and what ICE candidates are gathered. If the remote endpoint is bundle-aware, all media tracks and data channels are bundled onto the same transport.

    const (
     // BundlePolicyUnknown is the enum's zero-value.
     BundlePolicyUnknown BundlePolicy = [iota](/builtin#iota)
    
     // BundlePolicyBalanced indicates to gather ICE candidates for each
     // media type in use (audio, video, and data). If the remote endpoint is
     // not bundle-aware, negotiate only one audio and video track on separate
     // transports.
     BundlePolicyBalanced
    
     // BundlePolicyMaxCompat indicates to gather ICE candidates for each
     // track. If the remote endpoint is not bundle-aware, negotiate all media
     // tracks on separate transports.
     BundlePolicyMaxCompat
    
     // BundlePolicyMaxBundle indicates to gather ICE candidates for only
     // one track. If the remote endpoint is not bundle-aware, negotiate only
     // one media track.
     BundlePolicyMaxBundle
    )

#### func (BundlePolicy) [MarshalJSON](https://github.com/pion/webrtc/blob/v4.2.3/bundlepolicy.go#L83) ¶

    func (t BundlePolicy) MarshalJSON() ([][byte](/builtin#byte), [error](/builtin#error))

MarshalJSON returns the JSON encoding.

#### func (BundlePolicy) [String](https://github.com/pion/webrtc/blob/v4.2.3/bundlepolicy.go#L57) ¶

    func (t BundlePolicy) String() [string](/builtin#string)

#### func (*BundlePolicy) [UnmarshalJSON](https://github.com/pion/webrtc/blob/v4.2.3/bundlepolicy.go#L71) ¶

    func (t *BundlePolicy) UnmarshalJSON(b [][byte](/builtin#byte)) [error](/builtin#error)

UnmarshalJSON parses the JSON-encoded data and stores the result.

#### type [Certificate](https://github.com/pion/webrtc/blob/v4.2.3/certificate.go#L25) ¶

    type Certificate struct {
     // contains filtered or unexported fields
    }

Certificate represents a x509Cert used to authenticate WebRTC communications.

#### func [CertificateFromPEM](https://github.com/pion/webrtc/blob/v4.2.3/certificate.go#L191) ¶

    func CertificateFromPEM(pems [string](/builtin#string)) (*Certificate, [error](/builtin#error))

CertificateFromPEM creates a fresh certificate based on a string containing pem blocks fort the private key and x509 certificate.

#### func [CertificateFromX509](https://github.com/pion/webrtc/blob/v4.2.3/certificate.go#L160) ¶

    func CertificateFromX509(privateKey [crypto](/crypto).[PrivateKey](/crypto#PrivateKey), certificate *[x509](/crypto/x509).[Certificate](/crypto/x509#Certificate)) Certificate

CertificateFromX509 creates a new WebRTC Certificate from a given PrivateKey and Certificate

This can be used if you want to share a certificate across multiple PeerConnections.

#### func [GenerateCertificate](https://github.com/pion/webrtc/blob/v4.2.3/certificate.go#L136) ¶

    func GenerateCertificate(secretKey [crypto](/crypto).[PrivateKey](/crypto#PrivateKey)) (*Certificate, [error](/builtin#error))

GenerateCertificate causes the creation of an X.509 certificate and corresponding private key.

#### func [NewCertificate](https://github.com/pion/webrtc/blob/v4.2.3/certificate.go#L35) ¶

    func NewCertificate(key [crypto](/crypto).[PrivateKey](/crypto#PrivateKey), tpl [x509](/crypto/x509).[Certificate](/crypto/x509#Certificate)) (*Certificate, [error](/builtin#error))

NewCertificate generates a new x509 compliant Certificate to be used by DTLS for encrypting data sent over the wire. This method differs from GenerateCertificate by allowing to specify a template x509.Certificate to be used in order to define certificate parameters.

#### func (Certificate) [Equals](https://github.com/pion/webrtc/blob/v4.2.3/certificate.go#L71) ¶

    func (c Certificate) Equals(cert Certificate) [bool](/builtin#bool)

Equals determines if two certificates are identical by comparing both the secretKeys and x509Certificates.

#### func (Certificate) [Expires](https://github.com/pion/webrtc/blob/v4.2.3/certificate.go#L99) ¶

    func (c Certificate) Expires() [time](/time).[Time](/time#Time)

Expires returns the timestamp after which this certificate is no longer valid.

#### func (Certificate) [GetFingerprints](https://github.com/pion/webrtc/blob/v4.2.3/certificate.go#L109) ¶

    func (c Certificate) GetFingerprints() ([]DTLSFingerprint, [error](/builtin#error))

GetFingerprints returns the list of certificate fingerprints, one of which is computed with the digest algorithm used in the certificate signature.

#### func (Certificate) [PEM](https://github.com/pion/webrtc/blob/v4.2.3/certificate.go#L244) ¶

    func (c Certificate) PEM() ([string](/builtin#string), [error](/builtin#error))

PEM returns the certificate encoded as two pem block: once for the X509 certificate and the other for the private key.

#### type [CertificateStats](https://github.com/pion/webrtc/blob/v4.2.3/stats.go#L2355) ¶

    type CertificateStats struct {
     // Timestamp is the timestamp associated with this object.
     Timestamp StatsTimestamp `json:"timestamp"`
    
     // Type is the object's StatsType
     Type StatsType `json:"type"`
    
     // ID is a unique id that is associated with the component inspected to produce
     // this Stats object. Two Stats objects will have the same ID if they were produced
     // by inspecting the same underlying object.
     ID [string](/builtin#string) `json:"id"`
    
     // Fingerprint is the fingerprint of the certificate.
     Fingerprint [string](/builtin#string) `json:"fingerprint"`
    
     // FingerprintAlgorithm is the hash function used to compute the certificate fingerprint. For instance, "sha-256".
     FingerprintAlgorithm [string](/builtin#string) `json:"fingerprintAlgorithm"`
    
     // Base64Certificate is the DER-encoded base-64 representation of the certificate.
     Base64Certificate [string](/builtin#string) `json:"base64Certificate"`
    
     // IssuerCertificateID refers to the stats object that contains the next certificate
     // in the certificate chain. If the current certificate is at the end of the chain
     // (i.e. a self-signed certificate), this will not be set.
     IssuerCertificateID [string](/builtin#string) `json:"issuerCertificateId"`
    }

CertificateStats contains information about a certificate used by an ICETransport.

#### type [CodecStats](https://github.com/pion/webrtc/blob/v4.2.3/stats.go#L225) ¶

    type CodecStats struct {
     // Timestamp is the timestamp associated with this object.
     Timestamp StatsTimestamp `json:"timestamp"`
    
     // Type is the object's StatsType
     Type StatsType `json:"type"`
    
     // ID is a unique id that is associated with the component inspected to produce
     // this Stats object. Two Stats objects will have the same ID if they were produced
     // by inspecting the same underlying object.
     ID [string](/builtin#string) `json:"id"`
    
     // PayloadType as used in RTP encoding or decoding
     PayloadType PayloadType `json:"payloadType"`
    
     // CodecType of this CodecStats
     CodecType CodecType `json:"codecType"`
    
     // TransportID is the unique identifier of the transport on which this codec is
     // being used, which can be used to look up the corresponding TransportStats object.
     TransportID [string](/builtin#string) `json:"transportId"`
    
     // MimeType is the codec MIME media type/subtype. e.g., video/vp8 or equivalent.
     MimeType [string](/builtin#string) `json:"mimeType"`
    
     // ClockRate represents the media sampling rate.
     ClockRate [uint32](/builtin#uint32) `json:"clockRate"`
    
     // Channels is 2 for stereo, missing for most other cases.
     Channels [uint8](/builtin#uint8) `json:"channels"`
    
     // SDPFmtpLine is the a=fmtp line in the SDP corresponding to the codec,
     // i.e., after the colon following the PT.
     SDPFmtpLine [string](/builtin#string) `json:"sdpFmtpLine"`
    
     // Implementation identifies the implementation used. This is useful for diagnosing
     // interoperability issues.
     Implementation [string](/builtin#string) `json:"implementation"`
    }

CodecStats contains statistics for a codec that is currently being used by RTP streams being sent or received by this PeerConnection object.

#### type [CodecType](https://github.com/pion/webrtc/blob/v4.2.3/stats.go#L211) ¶

    type CodecType [string](/builtin#string)

CodecType specifies whether a CodecStats objects represents a media format that is being encoded or decoded.

    const (
     // CodecTypeEncode means the attached CodecStats represents a media format that
     // is being encoded, or that the implementation is prepared to encode.
     CodecTypeEncode CodecType = "encode"
    
     // CodecTypeDecode means the attached CodecStats represents a media format
     // that the implementation is prepared to decode.
     CodecTypeDecode CodecType = "decode"
    )

#### type [Configuration](https://github.com/pion/webrtc/blob/v4.2.3/configuration.go#L14) ¶

    type Configuration struct {
     // ICEServers defines a slice describing servers available to be used by
     // ICE, such as STUN and TURN servers.
     ICEServers []ICEServer `json:"iceServers,omitempty"`
    
     // ICETransportPolicy indicates which candidates the ICEAgent is allowed
     // to use.
     ICETransportPolicy ICETransportPolicy `json:"iceTransportPolicy,omitempty"`
    
     // BundlePolicy indicates which media-bundling policy to use when gathering
     // ICE candidates.
     BundlePolicy BundlePolicy `json:"bundlePolicy,omitempty"`
    
     // RTCPMuxPolicy indicates which rtcp-mux policy to use when gathering ICE
     // candidates.
     RTCPMuxPolicy RTCPMuxPolicy `json:"rtcpMuxPolicy,omitempty"`
    
     // PeerIdentity sets the target peer identity for the PeerConnection.
     // The PeerConnection will not establish a connection to a remote peer
     // unless it can be successfully authenticated with the provided name.
     PeerIdentity [string](/builtin#string) `json:"peerIdentity,omitempty"`
    
     // Certificates describes a set of certificates that the PeerConnection
     // uses to authenticate. Valid values for this parameter are created
     // through calls to the GenerateCertificate function. Although any given
     // DTLS connection will use only one certificate, this attribute allows the
     // caller to provide multiple certificates that support different
     // algorithms. The final certificate will be selected based on the DTLS
     // handshake, which establishes which certificates are allowed. The
     // PeerConnection implementation selects which of the certificates is
     // used for a given connection; how certificates are selected is outside
     // the scope of this specification. If this value is absent, then a default
     // set of certificates is generated for each PeerConnection instance.
     Certificates []Certificate `json:"certificates,omitempty"`
    
     // ICECandidatePoolSize describes the size of the prefetched ICE pool.
     ICECandidatePoolSize [uint8](/builtin#uint8) `json:"iceCandidatePoolSize,omitempty"`
    
     // SDPSemantics controls the type of SDP offers accepted by and
     // SDP answers generated by the PeerConnection.
     SDPSemantics SDPSemantics `json:"sdpSemantics,omitempty"`
    }

A Configuration defines how peer-to-peer communication via PeerConnection is established or re-established. Configurations may be set up once and reused across multiple connections. Configurations are treated as readonly. As long as they are unmodified, they are safe for concurrent use.

#### type [DTLSFingerprint](https://github.com/pion/webrtc/blob/v4.2.3/dtlsfingerprint.go#L8) ¶

    type DTLSFingerprint struct {
     // Algorithm specifies one of the hash function algorithms defined in
     // the 'Hash function Textual Names' registry.
     Algorithm [string](/builtin#string) `json:"algorithm"`
    
     // Value specifies the value of the certificate fingerprint in lowercase
     // hex string as expressed utilizing the syntax of 'fingerprint' in
     // <https://tools.ietf.org/html/rfc4572#section-5>.
     Value [string](/builtin#string) `json:"value"`
    }

DTLSFingerprint specifies the hash function algorithm and certificate fingerprint as described in <https://tools.ietf.org/html/rfc4572>.

#### type [DTLSParameters](https://github.com/pion/webrtc/blob/v4.2.3/dtlsparameters.go#L7) ¶

    type DTLSParameters struct {
     Role         DTLSRole          `json:"role"`
     Fingerprints []DTLSFingerprint `json:"fingerprints"`
    }

DTLSParameters holds information relating to DTLS configuration.

#### type [DTLSRole](https://github.com/pion/webrtc/blob/v4.2.3/dtlsrole.go#L11) ¶

    type DTLSRole [byte](/builtin#byte)

DTLSRole indicates the role of the DTLS transport.

    const (
     // DTLSRoleUnknown is the enum's zero-value.
     DTLSRoleUnknown DTLSRole = [iota](/builtin#iota)
    
     // DTLSRoleAuto defines the DTLS role is determined based on
     // the resolved ICE role: the ICE controlled role acts as the DTLS
     // client and the ICE controlling role acts as the DTLS server.
     DTLSRoleAuto
    
     // DTLSRoleClient defines the DTLS client role.
     DTLSRoleClient
    
     // DTLSRoleServer defines the DTLS server role.
     DTLSRoleServer
    )

#### func (DTLSRole) [String](https://github.com/pion/webrtc/blob/v4.2.3/dtlsrole.go#L48) ¶

    func (r DTLSRole) String() [string](/builtin#string)

#### type [DTLSTransport](https://github.com/pion/webrtc/blob/v4.2.3/dtlstransport.go#L37) ¶

    type DTLSTransport struct {
     // contains filtered or unexported fields
    }

DTLSTransport allows an application access to information about the DTLS transport over which RTP and RTCP packets are sent and received by RTPSender and RTPReceiver, as well other data such as SCTP packets sent and received by data channels.

#### func (*DTLSTransport) [GetLocalParameters](https://github.com/pion/webrtc/blob/v4.2.3/dtlstransport.go#L168) ¶

    func (t *DTLSTransport) GetLocalParameters() (DTLSParameters, [error](/builtin#error))

GetLocalParameters returns the DTLS parameters of the local DTLSTransport upon construction.

#### func (*DTLSTransport) [GetRemoteCertificate](https://github.com/pion/webrtc/blob/v4.2.3/dtlstransport.go#L188) ¶

    func (t *DTLSTransport) GetRemoteCertificate() [][byte](/builtin#byte)

GetRemoteCertificate returns the certificate chain in use by the remote side returns an empty list prior to selection of the remote certificate.

#### func (*DTLSTransport) [ICETransport](https://github.com/pion/webrtc/blob/v4.2.3/dtlstransport.go#L113) ¶

    func (t *DTLSTransport) ICETransport() *ICETransport

ICETransport returns the currently-configured *ICETransport or nil if one has not been configured.

#### func (*DTLSTransport) [OnStateChange](https://github.com/pion/webrtc/blob/v4.2.3/dtlstransport.go#L131) ¶

    func (t *DTLSTransport) OnStateChange(f func(DTLSTransportState))

OnStateChange sets a handler that is fired when the DTLS connection state changes.

#### func (*DTLSTransport) [Start](https://github.com/pion/webrtc/blob/v4.2.3/dtlstransport.go#L304) ¶

    func (t *DTLSTransport) Start(remoteParameters DTLSParameters) [error](/builtin#error)

Start DTLS transport negotiation with the parameters of the remote DTLS transport.

#### func (*DTLSTransport) [State](https://github.com/pion/webrtc/blob/v4.2.3/dtlstransport.go#L138) ¶

    func (t *DTLSTransport) State() DTLSTransportState

State returns the current dtls transport state.

#### func (*DTLSTransport) [Stop](https://github.com/pion/webrtc/blob/v4.2.3/dtlstransport.go#L451) ¶

    func (t *DTLSTransport) Stop() [error](/builtin#error)

Stop stops and closes the DTLSTransport object.

#### func (*DTLSTransport) [WriteRTCP](https://github.com/pion/webrtc/blob/v4.2.3/dtlstransport.go#L147) ¶

    func (t *DTLSTransport) WriteRTCP(pkts [][rtcp](/github.com/pion/rtcp).[Packet](/github.com/pion/rtcp#Packet)) ([int](/builtin#int), [error](/builtin#error))

WriteRTCP sends a user provided RTCP packet to the connected peer. If no peer is connected the packet is discarded.

#### type [DTLSTransportState](https://github.com/pion/webrtc/blob/v4.2.3/dtlstransportstate.go#L7) ¶

    type DTLSTransportState [int](/builtin#int)

DTLSTransportState indicates the DTLS transport establishment state.

    const (
     // DTLSTransportStateUnknown is the enum's zero-value.
     DTLSTransportStateUnknown DTLSTransportState = [iota](/builtin#iota)
    
     // DTLSTransportStateNew indicates that DTLS has not started negotiating
     // yet.
     DTLSTransportStateNew
    
     // DTLSTransportStateConnecting indicates that DTLS is in the process of
     // negotiating a secure connection and verifying the remote fingerprint.
     DTLSTransportStateConnecting
    
     // DTLSTransportStateConnected indicates that DTLS has completed
     // negotiation of a secure connection and verified the remote fingerprint.
     DTLSTransportStateConnected
    
     // DTLSTransportStateClosed indicates that the transport has been closed
     // intentionally as the result of receipt of a close_notify alert, or
     // calling close().
     DTLSTransportStateClosed
    
     // DTLSTransportStateFailed indicates that the transport has failed as
     // the result of an error (such as receipt of an error alert or failure to
     // validate the remote fingerprint).
     DTLSTransportStateFailed
    )

#### func (DTLSTransportState) [MarshalText](https://github.com/pion/webrtc/blob/v4.2.3/dtlstransportstate.go#L80) ¶

    func (t DTLSTransportState) MarshalText() ([][byte](/builtin#byte), [error](/builtin#error))

MarshalText implements encoding.TextMarshaler.

#### func (DTLSTransportState) [String](https://github.com/pion/webrtc/blob/v4.2.3/dtlstransportstate.go#L62) ¶

    func (t DTLSTransportState) String() [string](/builtin#string)

#### func (*DTLSTransportState) [UnmarshalText](https://github.com/pion/webrtc/blob/v4.2.3/dtlstransportstate.go#L85) ¶

    func (t *DTLSTransportState) UnmarshalText(b [][byte](/builtin#byte)) [error](/builtin#error)

UnmarshalText implements encoding.TextUnmarshaler.

#### type [DataChannel](https://github.com/pion/webrtc/blob/v4.2.3/datachannel.go#L27) ¶

    type DataChannel struct {
     // contains filtered or unexported fields
    }

DataChannel represents a WebRTC DataChannel The DataChannel interface represents a network channel which can be used for bidirectional peer-to-peer transfers of arbitrary data.

#### func (*DataChannel) [BufferedAmount](https://github.com/pion/webrtc/blob/v4.2.3/datachannel.go#L662) ¶

    func (d *DataChannel) BufferedAmount() [uint64](/builtin#uint64)

BufferedAmount represents the number of bytes of application data (UTF-8 text and binary data) that have been queued using send(). Even though the data transmission can occur in parallel, the returned value MUST NOT be decreased before the current task yielded back to the event loop to prevent race conditions. The value does not include framing overhead incurred by the protocol, or buffering done by the operating system or network hardware. The value of BufferedAmount slot will only increase with each call to the send() method as long as the ReadyState is open; however, BufferedAmount does not reset to zero once the channel closes.

#### func (*DataChannel) [BufferedAmountLowThreshold](https://github.com/pion/webrtc/blob/v4.2.3/datachannel.go#L679) ¶

    func (d *DataChannel) BufferedAmountLowThreshold() [uint64](/builtin#uint64)

BufferedAmountLowThreshold represents the threshold at which the bufferedAmount is considered to be low. When the bufferedAmount decreases from above this threshold to equal or below it, the bufferedamountlow event fires. BufferedAmountLowThreshold is initially zero on each new DataChannel, but the application may change its value at any time. The threshold is set to 0 by default.

#### func (*DataChannel) [Close](https://github.com/pion/webrtc/blob/v4.2.3/datachannel.go#L533) ¶

    func (d *DataChannel) Close() [error](/builtin#error)

Close Closes the DataChannel. It may be called regardless of whether the DataChannel object was created by this peer or the remote peer.

#### func (*DataChannel) [Detach](https://github.com/pion/webrtc/blob/v4.2.3/datachannel.go#L484) ¶

    func (d *DataChannel) Detach() ([datachannel](/github.com/pion/datachannel).[ReadWriteCloser](/github.com/pion/datachannel#ReadWriteCloser), [error](/builtin#error))

Detach allows you to detach the underlying datachannel. This provides an idiomatic API to work with (`io.ReadWriteCloser` with its `.Read()` and `.Write()` methods, as opposed to `.Send()` and `.OnMessage`), however it disables the OnMessage callback. Before calling Detach you have to enable this behavior by calling webrtc.DetachDataChannels(). Combining detached and normal data channels is not supported. Please refer to the data-channels-detach example and the pion/datachannel documentation for the correct way to handle the resulting DataChannel object.

#### func (*DataChannel) [DetachWithDeadline](https://github.com/pion/webrtc/blob/v4.2.3/datachannel.go#L490) ¶ added in v4.0.6

    func (d *DataChannel) DetachWithDeadline() ([datachannel](/github.com/pion/datachannel).[ReadWriteCloserDeadliner](/github.com/pion/datachannel#ReadWriteCloserDeadliner), [error](/builtin#error))

DetachWithDeadline allows you to detach the underlying datachannel. It is the same as Detach but returns a ReadWriteCloserDeadliner.

#### func (*DataChannel) [GracefulClose](https://github.com/pion/webrtc/blob/v4.2.3/datachannel.go#L541) ¶

    func (d *DataChannel) GracefulClose() [error](/builtin#error)

GracefulClose Closes the DataChannel. It may be called regardless of whether the DataChannel object was created by this peer or the remote peer. It also waits for any goroutines it started to complete. This is only safe to call outside of DataChannel callbacks or if in a callback, in its own goroutine.

#### func (*DataChannel) [ID](https://github.com/pion/webrtc/blob/v4.2.3/datachannel.go#L636) ¶

    func (d *DataChannel) ID() *[uint16](/builtin#uint16)

ID represents the ID for this DataChannel. The value is initially null, which is what will be returned if the ID was not provided at channel creation time, and the DTLS role of the SCTP transport has not yet been negotiated. Otherwise, it will return the ID that was either selected by the script or generated. After the ID is set to a non-null value, it will not change.

#### func (*DataChannel) [Label](https://github.com/pion/webrtc/blob/v4.2.3/datachannel.go#L578) ¶

    func (d *DataChannel) Label() [string](/builtin#string)

Label represents a label that can be used to distinguish this DataChannel object from other DataChannel objects. Scripts are allowed to create multiple DataChannel objects with the same label.

#### func (*DataChannel) [MaxPacketLifeTime](https://github.com/pion/webrtc/blob/v4.2.3/datachannel.go#L596) ¶

    func (d *DataChannel) MaxPacketLifeTime() *[uint16](/builtin#uint16)

MaxPacketLifeTime represents the length of the time window (msec) during which transmissions and retransmissions may occur in unreliable mode.

#### func (*DataChannel) [MaxRetransmits](https://github.com/pion/webrtc/blob/v4.2.3/datachannel.go#L605) ¶

    func (d *DataChannel) MaxRetransmits() *[uint16](/builtin#uint16)

MaxRetransmits represents the maximum number of retransmissions that are attempted in unreliable mode.

#### func (*DataChannel) [Negotiated](https://github.com/pion/webrtc/blob/v4.2.3/datachannel.go#L623) ¶

    func (d *DataChannel) Negotiated() [bool](/builtin#bool)

Negotiated represents whether this DataChannel was negotiated by the application (true), or not (false).

#### func (*DataChannel) [OnBufferedAmountLow](https://github.com/pion/webrtc/blob/v4.2.3/datachannel.go#L706) ¶

    func (d *DataChannel) OnBufferedAmountLow(f func())

OnBufferedAmountLow sets an event handler which is invoked when the number of bytes of outgoing data becomes lower than or equal to the BufferedAmountLowThreshold.

#### func (*DataChannel) [OnClose](https://github.com/pion/webrtc/blob/v4.2.3/datachannel.go#L286) ¶

    func (d *DataChannel) OnClose(f func())

OnClose sets an event handler which is invoked when the underlying data transport has been closed. Note: Due to backwards compatibility, there is a chance that OnClose can be called, even if the GracefulClose is used. If this is the case for you, you can deregister OnClose prior to GracefulClose.

#### func (*DataChannel) [OnDial](https://github.com/pion/webrtc/blob/v4.2.3/datachannel.go#L253) ¶

    func (d *DataChannel) OnDial(f func())

OnDial sets an event handler which is invoked when the peer has been dialed, but before said peer has responded.

#### func (*DataChannel) [OnError](https://github.com/pion/webrtc/blob/v4.2.3/datachannel.go#L377) ¶

    func (d *DataChannel) OnError(f func(err [error](/builtin#error)))

OnError sets an event handler which is invoked when the underlying data transport cannot be read.

#### func (*DataChannel) [OnMessage](https://github.com/pion/webrtc/blob/v4.2.3/datachannel.go#L308) ¶

    func (d *DataChannel) OnMessage(f func(msg DataChannelMessage))

OnMessage sets an event handler which is invoked on a binary message arrival over the sctp transport from a remote peer. OnMessage can currently receive messages up to 16384 bytes in size. Check out the detach API if you want to use larger message sizes. Note that browser support for larger messages is also limited.

#### func (*DataChannel) [OnOpen](https://github.com/pion/webrtc/blob/v4.2.3/datachannel.go#L218) ¶

    func (d *DataChannel) OnOpen(f func())

OnOpen sets an event handler which is invoked when the underlying data transport has been established (or re-established).

#### func (*DataChannel) [Ordered](https://github.com/pion/webrtc/blob/v4.2.3/datachannel.go#L587) ¶

    func (d *DataChannel) Ordered() [bool](/builtin#bool)

Ordered returns true if the DataChannel is ordered, and false if out-of-order delivery is allowed.

#### func (*DataChannel) [Protocol](https://github.com/pion/webrtc/blob/v4.2.3/datachannel.go#L614) ¶

    func (d *DataChannel) Protocol() [string](/builtin#string)

Protocol represents the name of the sub-protocol used with this DataChannel.

#### func (*DataChannel) [ReadyState](https://github.com/pion/webrtc/blob/v4.2.3/datachannel.go#L644) ¶

    func (d *DataChannel) ReadyState() DataChannelState

ReadyState represents the state of the DataChannel object.

#### func (*DataChannel) [Send](https://github.com/pion/webrtc/blob/v4.2.3/datachannel.go#L440) ¶

    func (d *DataChannel) Send(data [][byte](/builtin#byte)) [error](/builtin#error)

Send sends the binary message to the DataChannel peer.

#### func (*DataChannel) [SendText](https://github.com/pion/webrtc/blob/v4.2.3/datachannel.go#L452) ¶

    func (d *DataChannel) SendText(s [string](/builtin#string)) [error](/builtin#error)

SendText sends the text message to the DataChannel peer.

#### func (*DataChannel) [SetBufferedAmountLowThreshold](https://github.com/pion/webrtc/blob/v4.2.3/datachannel.go#L692) ¶

    func (d *DataChannel) SetBufferedAmountLowThreshold(th [uint64](/builtin#uint64))

SetBufferedAmountLowThreshold is used to update the threshold. See BufferedAmountLowThreshold().

#### func (*DataChannel) [Transport](https://github.com/pion/webrtc/blob/v4.2.3/datachannel.go#L198) ¶

    func (d *DataChannel) Transport() *SCTPTransport

Transport returns the SCTPTransport instance the DataChannel is sending over.

#### type [DataChannelInit](https://github.com/pion/webrtc/blob/v4.2.3/datachannelinit.go#L8) ¶

    type DataChannelInit struct {
     // Ordered indicates if data is allowed to be delivered out of order. The
     // default value of true, guarantees that data will be delivered in order.
     Ordered *[bool](/builtin#bool)
    
     // MaxPacketLifeTime limits the time (in milliseconds) during which the
     // channel will transmit or retransmit data if not acknowledged. This value
     // may be clamped if it exceeds the maximum value supported.
     MaxPacketLifeTime *[uint16](/builtin#uint16)
    
     // MaxRetransmits limits the number of times a channel will retransmit data
     // if not successfully delivered. This value may be clamped if it exceeds
     // the maximum value supported.
     MaxRetransmits *[uint16](/builtin#uint16)
    
     // Protocol describes the subprotocol name used for this channel.
     Protocol *[string](/builtin#string)
    
     // Negotiated describes if the data channel is created by the local peer or
     // the remote peer. The default value of false tells the user agent to
     // announce the channel in-band and instruct the other peer to dispatch a
     // corresponding DataChannel. If set to true, it is up to the application
     // to negotiate the channel and create an DataChannel with the same id
     // at the other peer.
     Negotiated *[bool](/builtin#bool)
    
     // ID overrides the default selection of ID for this channel.
     ID *[uint16](/builtin#uint16)
    }

DataChannelInit can be used to configure properties of the underlying channel such as data reliability.

#### type [DataChannelMessage](https://github.com/pion/webrtc/blob/v4.2.3/datachannelmessage.go#L10) ¶

    type DataChannelMessage struct {
     IsString [bool](/builtin#bool)
     Data     [][byte](/builtin#byte)
    }

DataChannelMessage represents a message received from the data channel. IsString will be set to true if the incoming message is of the string type. Otherwise the message is of a binary type.

#### type [DataChannelParameters](https://github.com/pion/webrtc/blob/v4.2.3/datachannelparameters.go#L7) ¶

    type DataChannelParameters struct {
     Label             [string](/builtin#string)  `json:"label"`
     Protocol          [string](/builtin#string)  `json:"protocol"`
     ID                *[uint16](/builtin#uint16) `json:"id"`
     Ordered           [bool](/builtin#bool)    `json:"ordered"`
     MaxPacketLifeTime *[uint16](/builtin#uint16) `json:"maxPacketLifeTime"`
     MaxRetransmits    *[uint16](/builtin#uint16) `json:"maxRetransmits"`
     Negotiated        [bool](/builtin#bool)    `json:"negotiated"`
    }

DataChannelParameters describes the configuration of the DataChannel.

#### type [DataChannelState](https://github.com/pion/webrtc/blob/v4.2.3/datachannelstate.go#L7) ¶

    type DataChannelState [int](/builtin#int)

DataChannelState indicates the state of a data channel.

    const (
     // DataChannelStateUnknown is the enum's zero-value.
     DataChannelStateUnknown DataChannelState = [iota](/builtin#iota)
    
     // DataChannelStateConnecting indicates that the data channel is being
     // established. This is the initial state of DataChannel, whether created
     // with CreateDataChannel, or dispatched as a part of an DataChannelEvent.
     DataChannelStateConnecting
    
     // DataChannelStateOpen indicates that the underlying data transport is
     // established and communication is possible.
     DataChannelStateOpen
    
     // DataChannelStateClosing indicates that the procedure to close down the
     // underlying data transport has started.
     DataChannelStateClosing
    
     // DataChannelStateClosed indicates that the underlying data transport
     // has been closed or could not be established.
     DataChannelStateClosed
    )

#### func (DataChannelState) [MarshalText](https://github.com/pion/webrtc/blob/v4.2.3/datachannelstate.go#L70) ¶

    func (t DataChannelState) MarshalText() ([][byte](/builtin#byte), [error](/builtin#error))

MarshalText implements encoding.TextMarshaler.

#### func (DataChannelState) [String](https://github.com/pion/webrtc/blob/v4.2.3/datachannelstate.go#L54) ¶

    func (t DataChannelState) String() [string](/builtin#string)

#### func (*DataChannelState) [UnmarshalText](https://github.com/pion/webrtc/blob/v4.2.3/datachannelstate.go#L75) ¶

    func (t *DataChannelState) UnmarshalText(b [][byte](/builtin#byte)) [error](/builtin#error)

UnmarshalText implements encoding.TextUnmarshaler.

#### type [DataChannelStats](https://github.com/pion/webrtc/blob/v4.2.3/stats.go#L1414) ¶

    type DataChannelStats struct {
     // Timestamp is the timestamp associated with this object.
     Timestamp StatsTimestamp `json:"timestamp"`
    
     // Type is the object's StatsType
     Type StatsType `json:"type"`
    
     // ID is a unique id that is associated with the component inspected to produce
     // this Stats object. Two Stats objects will have the same ID if they were produced
     // by inspecting the same underlying object.
     ID [string](/builtin#string) `json:"id"`
    
     // Label is the "label" value of the DataChannel object.
     Label [string](/builtin#string) `json:"label"`
    
     // Protocol is the "protocol" value of the DataChannel object.
     Protocol [string](/builtin#string) `json:"protocol"`
    
     // DataChannelIdentifier is the "id" attribute of the DataChannel object.
     DataChannelIdentifier [int32](/builtin#int32) `json:"dataChannelIdentifier"`
    
     // TransportID the ID of the TransportStats object for transport used to carry this datachannel.
     TransportID [string](/builtin#string) `json:"transportId"`
    
     // State is the "readyState" value of the DataChannel object.
     State DataChannelState `json:"state"`
    
     // MessagesSent represents the total number of API "message" events sent.
     MessagesSent [uint32](/builtin#uint32) `json:"messagesSent"`
    
     // BytesSent represents the total number of payload bytes sent on this
     // datachannel not including headers or padding.
     BytesSent [uint64](/builtin#uint64) `json:"bytesSent"`
    
     // MessagesReceived represents the total number of API "message" events received.
     MessagesReceived [uint32](/builtin#uint32) `json:"messagesReceived"`
    
     // BytesReceived represents the total number of bytes received on this
     // datachannel not including headers or padding.
     BytesReceived [uint64](/builtin#uint64) `json:"bytesReceived"`
    }

DataChannelStats contains statistics related to each DataChannel ID.

#### type [ICEAddressRewriteMode](https://github.com/pion/webrtc/blob/v4.2.3/icegatherer.go#L50) ¶ added in v4.2.0

    type ICEAddressRewriteMode [byte](/builtin#byte)

ICEAddressRewriteMode controls whether a rule replaces or appends candidates.

    const (
     ICEAddressRewriteModeUnspecified ICEAddressRewriteMode = [iota](/builtin#iota)
     ICEAddressRewriteReplace
     ICEAddressRewriteAppend
    )

#### type [ICEAddressRewriteRule](https://github.com/pion/webrtc/blob/v4.2.3/icegatherer.go#L63) ¶ added in v4.2.0

    type ICEAddressRewriteRule struct {
     External        [][string](/builtin#string)
     Local           [string](/builtin#string)
     Iface           [string](/builtin#string)
     CIDR            [string](/builtin#string)
     AsCandidateType ICECandidateType
     Mode            ICEAddressRewriteMode
     Networks        []NetworkType
    }

ICEAddressRewriteRule represents a rule for remapping candidate addresses.

#### type [ICECandidate](https://github.com/pion/webrtc/blob/v4.2.3/icecandidate.go#L13) ¶

    type ICECandidate struct {
     Foundation     [string](/builtin#string)           `json:"foundation"`
     Priority       [uint32](/builtin#uint32)           `json:"priority"`
     Address        [string](/builtin#string)           `json:"address"`
     Protocol       ICEProtocol      `json:"protocol"`
     Port           [uint16](/builtin#uint16)           `json:"port"`
     Typ            ICECandidateType `json:"type"`
     Component      [uint16](/builtin#uint16)           `json:"component"`
     RelatedAddress [string](/builtin#string)           `json:"relatedAddress"`
     RelatedPort    [uint16](/builtin#uint16)           `json:"relatedPort"`
     TCPType        [string](/builtin#string)           `json:"tcpType"`
     SDPMid         [string](/builtin#string)           `json:"sdpMid"`
     SDPMLineIndex  [uint16](/builtin#uint16)           `json:"sdpMLineIndex"`
     // contains filtered or unexported fields
    }

ICECandidate represents a ice candidate.

#### func (ICECandidate) [String](https://github.com/pion/webrtc/blob/v4.2.3/icecandidate.go#L219) ¶

    func (c ICECandidate) String() [string](/builtin#string)

#### func (ICECandidate) [ToICE](https://github.com/pion/webrtc/blob/v4.2.3/icecandidate.go#L84) ¶ added in v4.0.15

    func (c ICECandidate) ToICE() (cand [ice](/github.com/pion/ice/v4).[Candidate](/github.com/pion/ice/v4#Candidate), err [error](/builtin#error))

ToICE converts ICECandidate to ice.Candidate.

#### func (ICECandidate) [ToJSON](https://github.com/pion/webrtc/blob/v4.2.3/icecandidate.go#L230) ¶

    func (c ICECandidate) ToJSON() ICECandidateInit

ToJSON returns an ICECandidateInit as indicated by the spec <https://w3c.github.io/webrtc-pc/#dom-rtcicecandidate-tojson>

#### type [ICECandidateInit](https://github.com/pion/webrtc/blob/v4.2.3/icecandidateinit.go#L7) ¶

    type ICECandidateInit struct {
     Candidate        [string](/builtin#string)  `json:"candidate"`
     SDPMid           *[string](/builtin#string) `json:"sdpMid"`
     SDPMLineIndex    *[uint16](/builtin#uint16) `json:"sdpMLineIndex"`
     UsernameFragment *[string](/builtin#string) `json:"usernameFragment"`
    }

ICECandidateInit is used to serialize ice candidates.

#### type [ICECandidatePair](https://github.com/pion/webrtc/blob/v4.2.3/icecandidatepair.go#L9) ¶

    type ICECandidatePair struct {
     Local  *ICECandidate
     Remote *ICECandidate
     // contains filtered or unexported fields
    }

ICECandidatePair represents an ICE Candidate pair.

#### func [NewICECandidatePair](https://github.com/pion/webrtc/blob/v4.2.3/icecandidatepair.go#L25) ¶

    func NewICECandidatePair(local, remote *ICECandidate) *ICECandidatePair

NewICECandidatePair returns an initialized *ICECandidatePair for the given pair of ICECandidate instances.

#### func (*ICECandidatePair) [String](https://github.com/pion/webrtc/blob/v4.2.3/icecandidatepair.go#L19) ¶

    func (p *ICECandidatePair) String() [string](/builtin#string)

#### type [ICECandidatePairStats](https://github.com/pion/webrtc/blob/v4.2.3/stats.go#L2103) ¶

    type ICECandidatePairStats struct {
     // Timestamp is the timestamp associated with this object.
     Timestamp StatsTimestamp `json:"timestamp"`
    
     // Type is the object's StatsType
     Type StatsType `json:"type"`
    
     // ID is a unique id that is associated with the component inspected to produce
     // this Stats object. Two Stats objects will have the same ID if they were produced
     // by inspecting the same underlying object.
     ID [string](/builtin#string) `json:"id"`
    
     // TransportID is a unique identifier that is associated to the object that
     // was inspected to produce the TransportStats associated with this candidate pair.
     TransportID [string](/builtin#string) `json:"transportId"`
    
     // LocalCandidateID is a unique identifier that is associated to the object
     // that was inspected to produce the ICECandidateStats for the local candidate
     // associated with this candidate pair.
     LocalCandidateID [string](/builtin#string) `json:"localCandidateId"`
    
     // RemoteCandidateID is a unique identifier that is associated to the object
     // that was inspected to produce the ICECandidateStats for the remote candidate
     // associated with this candidate pair.
     RemoteCandidateID [string](/builtin#string) `json:"remoteCandidateId"`
    
     // State represents the state of the checklist for the local and remote
     // candidates in a pair.
     State StatsICECandidatePairState `json:"state"`
    
     // Nominated is true when this valid pair that should be used for media
     // if it is the highest-priority one amongst those whose nominated flag is set
     Nominated [bool](/builtin#bool) `json:"nominated"`
    
     // PacketsSent represents the total number of packets sent on this candidate pair.
     PacketsSent [uint32](/builtin#uint32) `json:"packetsSent"`
    
     // PacketsReceived represents the total number of packets received on this candidate pair.
     PacketsReceived [uint32](/builtin#uint32) `json:"packetsReceived"`
    
     // BytesSent represents the total number of payload bytes sent on this candidate pair
     // not including headers or padding.
     BytesSent [uint64](/builtin#uint64) `json:"bytesSent"`
    
     // BytesReceived represents the total number of payload bytes received on this candidate pair
     // not including headers or padding.
     BytesReceived [uint64](/builtin#uint64) `json:"bytesReceived"`
    
     // LastPacketSentTimestamp represents the timestamp at which the last packet was
     // sent on this particular candidate pair, excluding STUN packets.
     LastPacketSentTimestamp StatsTimestamp `json:"lastPacketSentTimestamp"`
    
     // LastPacketReceivedTimestamp represents the timestamp at which the last packet
     // was received on this particular candidate pair, excluding STUN packets.
     LastPacketReceivedTimestamp StatsTimestamp `json:"lastPacketReceivedTimestamp"`
    
     // FirstRequestTimestamp represents the timestamp at which the first STUN request
     // was sent on this particular candidate pair.
     FirstRequestTimestamp StatsTimestamp `json:"firstRequestTimestamp"`
    
     // LastRequestTimestamp represents the timestamp at which the last STUN request
     // was sent on this particular candidate pair. The average interval between two
     // consecutive connectivity checks sent can be calculated with
     // (LastRequestTimestamp - FirstRequestTimestamp) / RequestsSent.
     LastRequestTimestamp StatsTimestamp `json:"lastRequestTimestamp"`
    
     // FirstResponseTimestamp represents the timestamp at which the first STUN response
     // was received on this particular candidate pair.
     FirstResponseTimestamp StatsTimestamp `json:"firstResponseTimestamp"`
    
     // LastResponseTimestamp represents the timestamp at which the last STUN response
     // was received on this particular candidate pair.
     LastResponseTimestamp StatsTimestamp `json:"lastResponseTimestamp"`
    
     // FirstRequestReceivedTimestamp represents the timestamp at which the first
     // connectivity check request was received.
     FirstRequestReceivedTimestamp StatsTimestamp `json:"firstRequestReceivedTimestamp"`
    
     // LastRequestReceivedTimestamp represents the timestamp at which the last
     // connectivity check request was received.
     LastRequestReceivedTimestamp StatsTimestamp `json:"lastRequestReceivedTimestamp"`
    
     // TotalRoundTripTime represents the sum of all round trip time measurements
     // in seconds since the beginning of the session, based on STUN connectivity
     // check responses (ResponsesReceived), including those that reply to requests
     // that are sent in order to verify consent. The average round trip time can
     // be computed from TotalRoundTripTime by dividing it by ResponsesReceived.
     TotalRoundTripTime [float64](/builtin#float64) `json:"totalRoundTripTime"`
    
     // CurrentRoundTripTime represents the latest round trip time measured in seconds,
     // computed from both STUN connectivity checks, including those that are sent
     // for consent verification.
     CurrentRoundTripTime [float64](/builtin#float64) `json:"currentRoundTripTime"`
    
     // AvailableOutgoingBitrate is calculated by the underlying congestion control
     // by combining the available bitrate for all the outgoing RTP streams using
     // this candidate pair. The bitrate measurement does not count the size of the
     // IP or other transport layers like TCP or UDP. It is similar to the TIAS defined
     // in [RFC 3890](https://rfc-editor.org/rfc/rfc3890.html), i.e., it is measured in bits per second and the bitrate is calculated
     // over a 1 second window.
     AvailableOutgoingBitrate [float64](/builtin#float64) `json:"availableOutgoingBitrate"`
    
     // AvailableIncomingBitrate is calculated by the underlying congestion control
     // by combining the available bitrate for all the incoming RTP streams using
     // this candidate pair. The bitrate measurement does not count the size of the
     // IP or other transport layers like TCP or UDP. It is similar to the TIAS defined
     // in  [RFC 3890](https://rfc-editor.org/rfc/rfc3890.html), i.e., it is measured in bits per second and the bitrate is
     // calculated over a 1 second window.
     AvailableIncomingBitrate [float64](/builtin#float64) `json:"availableIncomingBitrate"`
    
     // CircuitBreakerTriggerCount represents the number of times the circuit breaker
     // is triggered for this particular 5-tuple, ceasing transmission.
     CircuitBreakerTriggerCount [uint32](/builtin#uint32) `json:"circuitBreakerTriggerCount"`
    
     // RequestsReceived represents the total number of connectivity check requests
     // received (including retransmissions). It is impossible for the receiver to
     // tell whether the request was sent in order to check connectivity or check
     // consent, so all connectivity checks requests are counted here.
     RequestsReceived [uint64](/builtin#uint64) `json:"requestsReceived"`
    
     // RequestsSent represents the total number of connectivity check requests
     // sent (not including retransmissions).
     RequestsSent [uint64](/builtin#uint64) `json:"requestsSent"`
    
     // ResponsesReceived represents the total number of connectivity check responses received.
     ResponsesReceived [uint64](/builtin#uint64) `json:"responsesReceived"`
    
     // ResponsesSent represents the total number of connectivity check responses sent.
     // Since we cannot distinguish connectivity check requests and consent requests,
     // all responses are counted.
     ResponsesSent [uint64](/builtin#uint64) `json:"responsesSent"`
    
     // RetransmissionsReceived represents the total number of connectivity check
     // request retransmissions received.
     RetransmissionsReceived [uint64](/builtin#uint64) `json:"retransmissionsReceived"`
    
     // RetransmissionsSent represents the total number of connectivity check
     // request retransmissions sent.
     RetransmissionsSent [uint64](/builtin#uint64) `json:"retransmissionsSent"`
    
     // ConsentRequestsSent represents the total number of consent requests sent.
     ConsentRequestsSent [uint64](/builtin#uint64) `json:"consentRequestsSent"`
    
     // ConsentExpiredTimestamp represents the timestamp at which the latest valid
     // STUN binding response expired.
     ConsentExpiredTimestamp StatsTimestamp `json:"consentExpiredTimestamp"`
    
     // PacketsDiscardedOnSend represents the total number of packets for this candidate pair
     // that have been discarded due to socket errors, i.e. a socket error occurred
     // when handing the packets to the socket. This might happen due to various reasons,
     // including full buffer or no available memory.
     PacketsDiscardedOnSend [uint32](/builtin#uint32) `json:"packetsDiscardedOnSend"`
    
     // BytesDiscardedOnSend represents the total number of bytes for this candidate pair
     // that have been discarded due to socket errors, i.e. a socket error occurred
     // when handing the packets containing the bytes to the socket. This might happen due
     // to various reasons, including full buffer or no available memory.
     // Calculated as defined in [RFC3550] section 6.4.1.
     BytesDiscardedOnSend [uint32](/builtin#uint32) `json:"bytesDiscardedOnSend"`
    }

ICECandidatePairStats contains ICE candidate pair statistics related to the ICETransport objects.

#### type [ICECandidateStats](https://github.com/pion/webrtc/blob/v4.2.3/stats.go#L2277) ¶

    type ICECandidateStats struct {
     // Timestamp is the timestamp associated with this object.
     Timestamp StatsTimestamp `json:"timestamp"`
    
     // Type is the object's StatsType
     Type StatsType `json:"type"`
    
     // ID is a unique id that is associated with the component inspected to produce
     // this Stats object. Two Stats objects will have the same ID if they were produced
     // by inspecting the same underlying object.
     ID [string](/builtin#string) `json:"id"`
    
     // TransportID is a unique identifier that is associated to the object that
     // was inspected to produce the TransportStats associated with this candidate.
     TransportID [string](/builtin#string) `json:"transportId"`
    
     // NetworkType represents the type of network interface used by the base of a
     // local candidate (the address the ICE agent sends from). Only present for
     // local candidates; it's not possible to know what type of network interface
     // a remote candidate is using.
     //
     // Note:
     // This stat only tells you about the network interface used by the first "hop";
     // it's possible that a connection will be bottlenecked by another type of network.
     // For example, when using Wi-Fi tethering, the networkType of the relevant candidate
     // would be "wifi", even when the next hop is over a cellular connection.
     //
     // DEPRECATED. Although it may still work in some browsers, the networkType property was deprecated for
     // preserving privacy.
     NetworkType [string](/builtin#string) `json:"networkType,omitempty"`
    
     // IP is the IP address of the candidate, allowing for IPv4 addresses and
     // IPv6 addresses, but fully qualified domain names (FQDNs) are not allowed.
     IP [string](/builtin#string) `json:"ip"`
    
     // Port is the port number of the candidate.
     Port [int32](/builtin#int32) `json:"port"`
    
     // Protocol is one of udp and tcp.
     Protocol [string](/builtin#string) `json:"protocol"`
    
     // CandidateType is the "Type" field of the ICECandidate.
     CandidateType ICECandidateType `json:"candidateType"`
    
     // Priority is the "Priority" field of the ICECandidate.
     Priority [int32](/builtin#int32) `json:"priority"`
    
     // URL of the TURN or STUN server that produced this candidate
     // It is the URL address surfaced in an PeerConnectionICEEvent.
     URL [string](/builtin#string) `json:"url"`
    
     // RelayProtocol is the protocol used by the endpoint to communicate with the
     // TURN server. This is only present for local candidates. Valid values for
     // the TURN URL protocol is one of udp, tcp, or tls.
     RelayProtocol [string](/builtin#string) `json:"relayProtocol"`
    
     // Deleted is true if the candidate has been deleted/freed. For host candidates,
     // this means that any network resources (typically a socket) associated with the
     // candidate have been released. For TURN candidates, this means the TURN allocation
     // is no longer active.
     //
     // Only defined for local candidates. For remote candidates, this property is not applicable.
     Deleted [bool](/builtin#bool) `json:"deleted"`
    }

ICECandidateStats contains ICE candidate statistics related to the ICETransport objects.

#### type [ICECandidateType](https://github.com/pion/webrtc/blob/v4.2.3/icecandidatetype.go#L13) ¶

    type ICECandidateType [int](/builtin#int)

ICECandidateType represents the type of the ICE candidate used.

    const (
     // ICECandidateTypeUnknown is the enum's zero-value.
     ICECandidateTypeUnknown ICECandidateType = [iota](/builtin#iota)
    
     // ICECandidateTypeHost indicates that the candidate is of Host type as
     // described in <https://tools.ietf.org/html/rfc8445#section-5.1.1.1>. A
     // candidate obtained by binding to a specific port from an IP address on
     // the host. This includes IP addresses on physical interfaces and logical
     // ones, such as ones obtained through VPNs.
     ICECandidateTypeHost
    
     // ICECandidateTypeSrflx indicates the candidate is of Server
     // Reflexive type as described
     // <https://tools.ietf.org/html/rfc8445#section-5.1.1.2>. A candidate type
     // whose IP address and port are a binding allocated by a NAT for an ICE
     // agent after it sends a packet through the NAT to a server, such as a
     // STUN server.
     ICECandidateTypeSrflx
    
     // ICECandidateTypePrflx indicates that the candidate is of Peer
     // Reflexive type. A candidate type whose IP address and port are a binding
     // allocated by a NAT for an ICE agent after it sends a packet through the
     // NAT to its peer.
     ICECandidateTypePrflx
    
     // ICECandidateTypeRelay indicates the candidate is of Relay type as
     // described in <https://tools.ietf.org/html/rfc8445#section-5.1.1.2>. A
     // candidate type obtained from a relay server, such as a TURN server.
     ICECandidateTypeRelay
    )

#### func [NewICECandidateType](https://github.com/pion/webrtc/blob/v4.2.3/icecandidatetype.go#L55) ¶

    func NewICECandidateType(raw [string](/builtin#string)) (ICECandidateType, [error](/builtin#error))

NewICECandidateType takes a string and converts it into ICECandidateType.

#### func (ICECandidateType) [MarshalText](https://github.com/pion/webrtc/blob/v4.2.3/icecandidatetype.go#L104) ¶

    func (t ICECandidateType) MarshalText() ([][byte](/builtin#byte), [error](/builtin#error))

MarshalText implements the encoding.TextMarshaler interface.

#### func (ICECandidateType) [String](https://github.com/pion/webrtc/blob/v4.2.3/icecandidatetype.go#L70) ¶

    func (t ICECandidateType) String() [string](/builtin#string)

#### func (*ICECandidateType) [UnmarshalText](https://github.com/pion/webrtc/blob/v4.2.3/icecandidatetype.go#L109) ¶

    func (t *ICECandidateType) UnmarshalText(b [][byte](/builtin#byte)) [error](/builtin#error)

UnmarshalText implements the encoding.TextUnmarshaler interface.

#### type [ICEComponent](https://github.com/pion/webrtc/blob/v4.2.3/icecomponent.go#L8) ¶

    type ICEComponent [int](/builtin#int)

ICEComponent describes if the ice transport is used for RTP (or RTCP multiplexing).

    const (
     // ICEComponentUnknown is the enum's zero-value.
     ICEComponentUnknown ICEComponent = [iota](/builtin#iota)
    
     // ICEComponentRTP indicates that the ICE Transport is used for RTP (or
     // RTCP multiplexing), as defined in
     // <https://tools.ietf.org/html/rfc5245#section-4.1.1.1>. Protocols
     // multiplexed with RTP (e.g. data channel) share its component ID. This
     // represents the component-id value 1 when encoded in candidate-attribute.
     ICEComponentRTP
    
     // ICEComponentRTCP indicates that the ICE Transport is used for RTCP as
     // defined by <https://tools.ietf.org/html/rfc5245#section-4.1.1.1>. This
     // represents the component-id value 2 when encoded in candidate-attribute.
     ICEComponentRTCP
    )

#### func (ICEComponent) [String](https://github.com/pion/webrtc/blob/v4.2.3/icecomponent.go#L44) ¶

    func (t ICEComponent) String() [string](/builtin#string)

#### type [ICEConnectionState](https://github.com/pion/webrtc/blob/v4.2.3/iceconnectionstate.go#L7) ¶

    type ICEConnectionState [int](/builtin#int)

ICEConnectionState indicates signaling state of the ICE Connection.

    const (
     // ICEConnectionStateUnknown is the enum's zero-value.
     ICEConnectionStateUnknown ICEConnectionState = [iota](/builtin#iota)
    
     // ICEConnectionStateNew indicates that any of the ICETransports are
     // in the "new" state and none of them are in the "checking", "disconnected"
     // or "failed" state, or all ICETransports are in the "closed" state, or
     // there are no transports.
     ICEConnectionStateNew
    
     // ICEConnectionStateChecking indicates that any of the ICETransports
     // are in the "checking" state and none of them are in the "disconnected"
     // or "failed" state.
     ICEConnectionStateChecking
    
     // ICEConnectionStateConnected indicates that all ICETransports are
     // in the "connected", "completed" or "closed" state and at least one of
     // them is in the "connected" state.
     ICEConnectionStateConnected
    
     // ICEConnectionStateCompleted indicates that all ICETransports are
     // in the "completed" or "closed" state and at least one of them is in the
     // "completed" state.
     ICEConnectionStateCompleted
    
     // ICEConnectionStateDisconnected indicates that any of the
     // ICETransports are in the "disconnected" state and none of them are
     // in the "failed" state.
     ICEConnectionStateDisconnected
    
     // ICEConnectionStateFailed indicates that any of the ICETransports
     // are in the "failed" state.
     ICEConnectionStateFailed
    
     // ICEConnectionStateClosed indicates that the PeerConnection's
     // isClosed is true.
     ICEConnectionStateClosed
    )

#### func [NewICEConnectionState](https://github.com/pion/webrtc/blob/v4.2.3/iceconnectionstate.go#L60) ¶

    func NewICEConnectionState(raw [string](/builtin#string)) ICEConnectionState

NewICEConnectionState takes a string and converts it to ICEConnectionState.

#### func (ICEConnectionState) [String](https://github.com/pion/webrtc/blob/v4.2.3/iceconnectionstate.go#L81) ¶

    func (c ICEConnectionState) String() [string](/builtin#string)

#### type [ICECredentialType](https://github.com/pion/webrtc/blob/v4.2.3/icecredentialtype.go#L13) ¶

    type ICECredentialType [int](/builtin#int)

ICECredentialType indicates the type of credentials used to connect to an ICE server.

    const (
     // ICECredentialTypePassword describes username and password based
     // credentials as described in <https://tools.ietf.org/html/rfc5389>.
     ICECredentialTypePassword ICECredentialType = [iota](/builtin#iota)
    
     // ICECredentialTypeOauth describes token based credential as described
     // in <https://tools.ietf.org/html/rfc7635>.
     ICECredentialTypeOauth
    )

#### func (ICECredentialType) [MarshalJSON](https://github.com/pion/webrtc/blob/v4.2.3/icecredentialtype.go#L71) ¶

    func (t ICECredentialType) MarshalJSON() ([][byte](/builtin#byte), [error](/builtin#error))

MarshalJSON returns the JSON encoding.

#### func (ICECredentialType) [String](https://github.com/pion/webrtc/blob/v4.2.3/icecredentialtype.go#L42) ¶

    func (t ICECredentialType) String() [string](/builtin#string)

#### func (*ICECredentialType) [UnmarshalJSON](https://github.com/pion/webrtc/blob/v4.2.3/icecredentialtype.go#L54) ¶

    func (t *ICECredentialType) UnmarshalJSON(b [][byte](/builtin#byte)) [error](/builtin#error)

UnmarshalJSON parses the JSON-encoded data and stores the result.

#### type [ICEGatherOptions](https://github.com/pion/webrtc/blob/v4.2.3/icegatheroptions.go#L7) ¶

    type ICEGatherOptions struct {
     ICEServers      []ICEServer
     ICEGatherPolicy ICETransportPolicy
    }

ICEGatherOptions provides options relating to the gathering of ICE candidates.

#### type [ICEGatherPolicy](https://github.com/pion/webrtc/blob/v4.2.3/icetransportpolicy.go#L15) ¶

    type ICEGatherPolicy = ICETransportPolicy

ICEGatherPolicy is the ORTC equivalent of ICETransportPolicy.

#### type [ICEGatherer](https://github.com/pion/webrtc/blob/v4.2.3/icegatherer.go#L25) ¶

    type ICEGatherer struct {
     // contains filtered or unexported fields
    }

ICEGatherer gathers local host, server reflexive and relay candidates, as well as enabling the retrieval of local Interactive Connectivity Establishment (ICE) parameters which can be exchanged in signaling.

#### func (*ICEGatherer) [Close](https://github.com/pion/webrtc/blob/v4.2.3/icegatherer.go#L460) ¶

    func (g *ICEGatherer) Close() [error](/builtin#error)

Close prunes all local candidates, and closes the ports.

#### func (*ICEGatherer) [Gather](https://github.com/pion/webrtc/blob/v4.2.3/icegatherer.go#L401) ¶

    func (g *ICEGatherer) Gather() [error](/builtin#error)

Gather ICE candidates.

#### func (*ICEGatherer) [GetLocalCandidates](https://github.com/pion/webrtc/blob/v4.2.3/icegatherer.go#L526) ¶

    func (g *ICEGatherer) GetLocalCandidates() ([]ICECandidate, [error](/builtin#error))

GetLocalCandidates returns the sequence of valid local candidates associated with the ICEGatherer.

#### func (*ICEGatherer) [GetLocalParameters](https://github.com/pion/webrtc/blob/v4.2.3/icegatherer.go#L502) ¶

    func (g *ICEGatherer) GetLocalParameters() (ICEParameters, [error](/builtin#error))

GetLocalParameters returns the ICE parameters of the ICEGatherer.

#### func (*ICEGatherer) [GracefulClose](https://github.com/pion/webrtc/blob/v4.2.3/icegatherer.go#L467) ¶

    func (g *ICEGatherer) GracefulClose() [error](/builtin#error)

GracefulClose prunes all local candidates, and closes the ports. It also waits for any goroutines it started to complete. This is only safe to call outside of ICEGatherer callbacks or if in a callback, in its own goroutine.

#### func (*ICEGatherer) [OnLocalCandidate](https://github.com/pion/webrtc/blob/v4.2.3/icegatherer.go#L554) ¶

    func (g *ICEGatherer) OnLocalCandidate(f func(*ICECandidate))

OnLocalCandidate sets an event handler which fires when a new local ICE candidate is available Take note that the handler will be called with a nil pointer when gathering is finished.

#### func (*ICEGatherer) [OnStateChange](https://github.com/pion/webrtc/blob/v4.2.3/icegatherer.go#L559) ¶

    func (g *ICEGatherer) OnStateChange(f func(ICEGathererState))

OnStateChange fires any time the ICEGatherer changes.

#### func (*ICEGatherer) [State](https://github.com/pion/webrtc/blob/v4.2.3/icegatherer.go#L564) ¶

    func (g *ICEGatherer) State() ICEGathererState

State indicates the current state of the ICE gatherer.

#### type [ICEGathererState](https://github.com/pion/webrtc/blob/v4.2.3/icegathererstate.go#L11) ¶

    type ICEGathererState [uint32](/builtin#uint32)

ICEGathererState represents the current state of the ICE gatherer.

    const (
     // ICEGathererStateUnknown is the enum's zero-value.
     ICEGathererStateUnknown ICEGathererState = [iota](/builtin#iota)
    
     // ICEGathererStateNew indicates object has been created but
     // gather() has not been called.
     ICEGathererStateNew
    
     // ICEGathererStateGathering indicates gather() has been called,
     // and the ICEGatherer is in the process of gathering candidates.
     ICEGathererStateGathering
    
     // ICEGathererStateComplete indicates the ICEGatherer has completed gathering.
     ICEGathererStateComplete
    
     // ICEGathererStateClosed indicates the closed state can only be entered
     // when the ICEGatherer has been closed intentionally by calling close().
     ICEGathererStateClosed
    )

#### func (ICEGathererState) [String](https://github.com/pion/webrtc/blob/v4.2.3/icegathererstate.go#L33) ¶

    func (s ICEGathererState) String() [string](/builtin#string)

#### type [ICEGatheringState](https://github.com/pion/webrtc/blob/v4.2.3/icegatheringstate.go#L7) ¶

    type ICEGatheringState [int](/builtin#int)

ICEGatheringState describes the state of the candidate gathering process.

    const (
     // ICEGatheringStateUnknown is the enum's zero-value.
     ICEGatheringStateUnknown ICEGatheringState = [iota](/builtin#iota)
    
     // ICEGatheringStateNew indicates that any of the ICETransports are
     // in the "new" gathering state and none of the transports are in the
     // "gathering" state, or there are no transports.
     ICEGatheringStateNew
    
     // ICEGatheringStateGathering indicates that any of the ICETransports
     // are in the "gathering" state.
     ICEGatheringStateGathering
    
     // ICEGatheringStateComplete indicates that at least one ICETransport
     // exists, and all ICETransports are in the "completed" gathering state.
     ICEGatheringStateComplete
    )

#### func [NewICEGatheringState](https://github.com/pion/webrtc/blob/v4.2.3/icegatheringstate.go#L35) ¶

    func NewICEGatheringState(raw [string](/builtin#string)) ICEGatheringState

NewICEGatheringState takes a string and converts it to ICEGatheringState.

#### func (ICEGatheringState) [String](https://github.com/pion/webrtc/blob/v4.2.3/icegatheringstate.go#L48) ¶

    func (t ICEGatheringState) String() [string](/builtin#string)

#### type [ICEParameters](https://github.com/pion/webrtc/blob/v4.2.3/iceparameters.go#L8) ¶

    type ICEParameters struct {
     UsernameFragment [string](/builtin#string) `json:"usernameFragment"`
     Password         [string](/builtin#string) `json:"password"`
     ICELite          [bool](/builtin#bool)   `json:"iceLite"`
    }

ICEParameters includes the ICE username fragment and password and other ICE-related parameters.

#### type [ICEProtocol](https://github.com/pion/webrtc/blob/v4.2.3/iceprotocol.go#L13) ¶

    type ICEProtocol [int](/builtin#int)

ICEProtocol indicates the transport protocol type that is used in the ice.URL structure.

    const (
     // ICEProtocolUnknown is the enum's zero-value.
     ICEProtocolUnknown ICEProtocol = [iota](/builtin#iota)
    
     // ICEProtocolUDP indicates the URL uses a UDP transport.
     ICEProtocolUDP
    
     // ICEProtocolTCP indicates the URL uses a TCP transport.
     ICEProtocolTCP
    )

#### func [NewICEProtocol](https://github.com/pion/webrtc/blob/v4.2.3/iceprotocol.go#L33) ¶

    func NewICEProtocol(raw [string](/builtin#string)) (ICEProtocol, [error](/builtin#error))

NewICEProtocol takes a string and converts it to ICEProtocol.

#### func (ICEProtocol) [String](https://github.com/pion/webrtc/blob/v4.2.3/iceprotocol.go#L44) ¶

    func (t ICEProtocol) String() [string](/builtin#string)

#### type [ICERole](https://github.com/pion/webrtc/blob/v4.2.3/icerole.go#L8) ¶

    type ICERole [int](/builtin#int)

ICERole describes the role ice.Agent is playing in selecting the preferred the candidate pair.

    const (
     // ICERoleUnknown is the enum's zero-value.
     ICERoleUnknown ICERole = [iota](/builtin#iota)
    
     // ICERoleControlling indicates that the ICE agent that is responsible
     // for selecting the final choice of candidate pairs and signaling them
     // through STUN and an updated offer, if needed. In any session, one agent
     // is always controlling. The other is the controlled agent.
     ICERoleControlling
    
     // ICERoleControlled indicates that an ICE agent that waits for the
     // controlling agent to select the final choice of candidate pairs.
     ICERoleControlled
    )

#### func (ICERole) [MarshalText](https://github.com/pion/webrtc/blob/v4.2.3/icerole.go#L54) ¶

    func (t ICERole) MarshalText() ([][byte](/builtin#byte), [error](/builtin#error))

MarshalText implements encoding.TextMarshaler.

#### func (ICERole) [String](https://github.com/pion/webrtc/blob/v4.2.3/icerole.go#L42) ¶

    func (t ICERole) String() [string](/builtin#string)

#### func (*ICERole) [UnmarshalText](https://github.com/pion/webrtc/blob/v4.2.3/icerole.go#L59) ¶

    func (t *ICERole) UnmarshalText(b [][byte](/builtin#byte)) [error](/builtin#error)

UnmarshalText implements encoding.TextUnmarshaler.

#### type [ICEServer](https://github.com/pion/webrtc/blob/v4.2.3/iceserver.go#L18) ¶

    type ICEServer struct {
     URLs           [][string](/builtin#string)          `json:"urls"`
     Username       [string](/builtin#string)            `json:"username,omitempty"`
     Credential     [any](/builtin#any)               `json:"credential,omitempty"`
     CredentialType ICECredentialType `json:"credentialType,omitempty"`
    }

ICEServer describes a single STUN and TURN server that can be used by the ICEAgent to establish a connection with a peer.

#### func (ICEServer) [MarshalJSON](https://github.com/pion/webrtc/blob/v4.2.3/iceserver.go#L176) ¶

    func (s ICEServer) MarshalJSON() ([][byte](/builtin#byte), [error](/builtin#error))

MarshalJSON returns the JSON encoding.

#### func (*ICEServer) [UnmarshalJSON](https://github.com/pion/webrtc/blob/v4.2.3/iceserver.go#L162) ¶

    func (s *ICEServer) UnmarshalJSON(b [][byte](/builtin#byte)) [error](/builtin#error)

UnmarshalJSON parses the JSON-encoded data and stores the result.

#### type [ICETransport](https://github.com/pion/webrtc/blob/v4.2.3/icetransport.go#L24) ¶

    type ICETransport struct {
     // contains filtered or unexported fields
    }

ICETransport allows an application access to information about the ICE transport over which packets are sent and received.

#### func [NewICETransport](https://github.com/pion/webrtc/blob/v4.2.3/icetransport.go#L79) ¶

    func NewICETransport(gatherer *ICEGatherer, loggerFactory [logging](/github.com/pion/logging).[LoggerFactory](/github.com/pion/logging#LoggerFactory)) *ICETransport

NewICETransport creates a new NewICETransport.

#### func (*ICETransport) [AddRemoteCandidate](https://github.com/pion/webrtc/blob/v4.2.3/icetransport.go#L315) ¶

    func (t *ICETransport) AddRemoteCandidate(remoteCandidate *ICECandidate) [error](/builtin#error)

AddRemoteCandidate adds a candidate associated with the remote ICETransport.

#### func (*ICETransport) [GetLocalParameters](https://github.com/pion/webrtc/blob/v4.2.3/icetransport.go#L353) ¶

    func (t *ICETransport) GetLocalParameters() (ICEParameters, [error](/builtin#error))

GetLocalParameters returns an IceParameters object which provides information uniquely identifying the local peer for the duration of the ICE session.

#### func (*ICETransport) [GetRemoteParameters](https://github.com/pion/webrtc/blob/v4.2.3/icetransport.go#L363) ¶ added in v4.0.10

    func (t *ICETransport) GetRemoteParameters() (ICEParameters, [error](/builtin#error))

GetRemoteParameters returns an IceParameters object which provides information uniquely identifying the remote peer for the duration of the ICE session.

#### func (*ICETransport) [GetSelectedCandidatePair](https://github.com/pion/webrtc/blob/v4.2.3/icetransport.go#L48) ¶

    func (t *ICETransport) GetSelectedCandidatePair() (*ICECandidatePair, [error](/builtin#error))

GetSelectedCandidatePair returns the selected candidate pair on which packets are sent if there is no selected pair nil is returned.

#### func (*ICETransport) [GetSelectedCandidatePairStats](https://github.com/pion/webrtc/blob/v4.2.3/icetransport.go#L74) ¶

    func (t *ICETransport) GetSelectedCandidatePairStats() (ICECandidatePairStats, [bool](/builtin#bool))

GetSelectedCandidatePairStats returns the selected candidate pair stats on which packets are sent if there is no selected pair empty stats, false is returned to indicate stats not available.

#### func (*ICETransport) [GracefulStop](https://github.com/pion/webrtc/blob/v4.2.3/icetransport.go#L213) ¶

    func (t *ICETransport) GracefulStop() [error](/builtin#error)

GracefulStop irreversibly stops the ICETransport. It also waits for any goroutines it started to complete. This is only safe to call outside of ICETransport callbacks or if in a callback, in its own goroutine.

#### func (*ICETransport) [OnConnectionStateChange](https://github.com/pion/webrtc/blob/v4.2.3/icetransport.go#L265) ¶

    func (t *ICETransport) OnConnectionStateChange(f func(ICETransportState))

OnConnectionStateChange sets a handler that is fired when the ICE connection state changes.

#### func (*ICETransport) [OnSelectedCandidatePairChange](https://github.com/pion/webrtc/blob/v4.2.3/icetransport.go#L253) ¶

    func (t *ICETransport) OnSelectedCandidatePairChange(f func(*ICECandidatePair))

OnSelectedCandidatePairChange sets a handler that is invoked when a new ICE candidate pair is selected.

#### func (*ICETransport) [Role](https://github.com/pion/webrtc/blob/v4.2.3/icetransport.go#L279) ¶

    func (t *ICETransport) Role() ICERole

Role indicates the current role of the ICE transport.

#### func (*ICETransport) [SetRemoteCandidates](https://github.com/pion/webrtc/blob/v4.2.3/icetransport.go#L287) ¶

    func (t *ICETransport) SetRemoteCandidates(remoteCandidates []ICECandidate) [error](/builtin#error)

SetRemoteCandidates sets the sequence of candidates associated with the remote ICETransport.

#### func (*ICETransport) [Start](https://github.com/pion/webrtc/blob/v4.2.3/icetransport.go#L91) ¶

    func (t *ICETransport) Start(gatherer *ICEGatherer, params ICEParameters, role *ICERole) [error](/builtin#error)

Start incoming connectivity checks based on its configured role.

#### func (*ICETransport) [State](https://github.com/pion/webrtc/blob/v4.2.3/icetransport.go#L343) ¶

    func (t *ICETransport) State() ICETransportState

State returns the current ice transport state.

#### func (*ICETransport) [Stop](https://github.com/pion/webrtc/blob/v4.2.3/icetransport.go#L206) ¶

    func (t *ICETransport) Stop() [error](/builtin#error)

Stop irreversibly stops the ICETransport.

#### type [ICETransportPolicy](https://github.com/pion/webrtc/blob/v4.2.3/icetransportpolicy.go#L12) ¶

    type ICETransportPolicy [int](/builtin#int)

ICETransportPolicy defines the ICE candidate policy surface the permitted candidates. Only these candidates are used for connectivity checks.

    const (
     // ICETransportPolicyAll indicates any type of candidate is used.
     ICETransportPolicyAll ICETransportPolicy = [iota](/builtin#iota)
    
     // ICETransportPolicyRelay indicates only media relay candidates such
     // as candidates passing through a TURN server are used.
     ICETransportPolicyRelay
    
     // ICETransportPolicyNoHost indicates only non-host candidates are used.
     ICETransportPolicyNoHost
    )

#### func [NewICETransportPolicy](https://github.com/pion/webrtc/blob/v4.2.3/icetransportpolicy.go#L37) ¶

    func NewICETransportPolicy(raw [string](/builtin#string)) ICETransportPolicy

NewICETransportPolicy takes a string and converts it to ICETransportPolicy.

#### func (ICETransportPolicy) [MarshalJSON](https://github.com/pion/webrtc/blob/v4.2.3/icetransportpolicy.go#L73) ¶

    func (t ICETransportPolicy) MarshalJSON() ([][byte](/builtin#byte), [error](/builtin#error))

MarshalJSON returns the JSON encoding.

#### func (ICETransportPolicy) [String](https://github.com/pion/webrtc/blob/v4.2.3/icetransportpolicy.go#L48) ¶

    func (t ICETransportPolicy) String() [string](/builtin#string)

#### func (*ICETransportPolicy) [UnmarshalJSON](https://github.com/pion/webrtc/blob/v4.2.3/icetransportpolicy.go#L62) ¶

    func (t *ICETransportPolicy) UnmarshalJSON(b [][byte](/builtin#byte)) [error](/builtin#error)

UnmarshalJSON parses the JSON-encoded data and stores the result.

#### type [ICETransportState](https://github.com/pion/webrtc/blob/v4.2.3/icetransportstate.go#L9) ¶

    type ICETransportState [int](/builtin#int)

ICETransportState represents the current state of the ICE transport.

    const (
     // ICETransportStateUnknown is the enum's zero-value.
     ICETransportStateUnknown ICETransportState = [iota](/builtin#iota)
    
     // ICETransportStateNew indicates the ICETransport is waiting
     // for remote candidates to be supplied.
     ICETransportStateNew
    
     // ICETransportStateChecking indicates the ICETransport has
     // received at least one remote candidate, and a local and remote
     // ICECandidateComplete dictionary was not added as the last candidate.
     ICETransportStateChecking
    
     // ICETransportStateConnected indicates the ICETransport has
     // received a response to an outgoing connectivity check, or has
     // received incoming DTLS/media after a successful response to an
     // incoming connectivity check, but is still checking other candidate
     // pairs to see if there is a better connection.
     ICETransportStateConnected
    
     // ICETransportStateCompleted indicates the ICETransport tested
     // all appropriate candidate pairs and at least one functioning
     // candidate pair has been found.
     ICETransportStateCompleted
    
     // ICETransportStateFailed indicates the ICETransport the last
     // candidate was added and all appropriate candidate pairs have either
     // failed connectivity checks or have lost consent.
     ICETransportStateFailed
    
     // ICETransportStateDisconnected indicates the ICETransport has received
     // at least one local and remote candidate, but the final candidate was
     // received yet and all appropriate candidate pairs thus far have been
     // tested and failed.
     ICETransportStateDisconnected
    
     // ICETransportStateClosed indicates the ICETransport has shut down
     // and is no longer responding to STUN requests.
     ICETransportStateClosed
    )

#### func (ICETransportState) [MarshalText](https://github.com/pion/webrtc/blob/v4.2.3/icetransportstate.go#L147) ¶

    func (c ICETransportState) MarshalText() ([][byte](/builtin#byte), [error](/builtin#error))

MarshalText implements encoding.TextMarshaler.

#### func (ICETransportState) [String](https://github.com/pion/webrtc/blob/v4.2.3/icetransportstate.go#L83) ¶

    func (c ICETransportState) String() [string](/builtin#string)

#### func (*ICETransportState) [UnmarshalText](https://github.com/pion/webrtc/blob/v4.2.3/icetransportstate.go#L152) ¶

    func (c *ICETransportState) UnmarshalText(b [][byte](/builtin#byte)) [error](/builtin#error)

UnmarshalText implements encoding.TextUnmarshaler.

#### type [ICETrickleCapability](https://github.com/pion/webrtc/blob/v4.2.3/sessiondescription.go#L16) ¶ added in v4.1.7

    type ICETrickleCapability [int](/builtin#int)

ICETrickleCapability represents whether the remote endpoint accepts trickled ICE candidates.

    const (
     // ICETrickleCapabilityUnknown no remote peer has been established.
     ICETrickleCapabilityUnknown ICETrickleCapability = [iota](/builtin#iota)
     // ICETrickleCapabilitySupported remote peer can accept trickled ICE candidates.
     ICETrickleCapabilitySupported
     // ICETrickleCapabilitySupported remote peer didn't state that it can accept trickle ICE candidates.
     ICETrickleCapabilityUnsupported
    )

#### func (ICETrickleCapability) [String](https://github.com/pion/webrtc/blob/v4.2.3/sessiondescription.go#L28) ¶ added in v4.1.7

    func (t ICETrickleCapability) String() [string](/builtin#string)

String returns the string representation of ICETrickleCapability.

#### type [InboundRTPStreamStats](https://github.com/pion/webrtc/blob/v4.2.3/stats.go#L279) ¶

    type InboundRTPStreamStats struct {
     // Mid represents a mid value of RTPTransceiver owning this stream, if that value is not
     // null. Otherwise, this member is not present.
     Mid [string](/builtin#string) `json:"mid"`
    
     // Timestamp is the timestamp associated with this object.
     Timestamp StatsTimestamp `json:"timestamp"`
    
     // Type is the object's StatsType
     Type StatsType `json:"type"`
    
     // ID is a unique id that is associated with the component inspected to produce
     // this Stats object. Two Stats objects will have the same ID if they were produced
     // by inspecting the same underlying object.
     ID [string](/builtin#string) `json:"id"`
    
     // SSRC is the 32-bit unsigned integer value used to identify the source of the
     // stream of RTP packets that this stats object concerns.
     SSRC SSRC `json:"ssrc"`
    
     // Kind is either "audio" or "video"
     Kind [string](/builtin#string) `json:"kind"`
    
     // It is a unique identifier that is associated to the object that was inspected
     // to produce the TransportStats associated with this RTP stream.
     TransportID [string](/builtin#string) `json:"transportId"`
    
     // CodecID is a unique identifier that is associated to the object that was inspected
     // to produce the CodecStats associated with this RTP stream.
     CodecID [string](/builtin#string) `json:"codecId"`
    
     // FIRCount counts the total number of Full Intra Request (FIR) packets received
     // by the sender. This metric is only valid for video and is sent by receiver.
     FIRCount [uint32](/builtin#uint32) `json:"firCount"`
    
     // PLICount counts the total number of Picture Loss Indication (PLI) packets
     // received by the sender. This metric is only valid for video and is sent by receiver.
     PLICount [uint32](/builtin#uint32) `json:"pliCount"`
    
     // TotalProcessingDelay is the sum of the time, in seconds, each audio sample or video frame
     // takes from the time the first RTP packet is received (reception timestamp) and to the time
     // the corresponding sample or frame is decoded (decoded timestamp). At this point the audio
     // sample or video frame is ready for playout by the MediaStreamTrack. Typically ready for
     // playout here means after the audio sample or video frame is fully decoded by the decoder.
     TotalProcessingDelay [float64](/builtin#float64) `json:"totalProcessingDelay"`
    
     // NACKCount counts the total number of Negative ACKnowledgement (NACK) packets
     // received by the sender and is sent by receiver.
     NACKCount [uint32](/builtin#uint32) `json:"nackCount"`
    
     // JitterBufferDelay is the sum of the time, in seconds, each audio sample or a video frame
     // takes from the time the first packet is received by the jitter buffer (ingest timestamp)
     // to the time it exits the jitter buffer (emit timestamp). The average jitter buffer delay
     // can be calculated by dividing the JitterBufferDelay with the JitterBufferEmittedCount.
     JitterBufferDelay [float64](/builtin#float64) `json:"jitterBufferDelay"`
    
     // JitterBufferTargetDelay is increased by the target jitter buffer delay every time a sample is emitted
     // by the jitter buffer. The added target is the target delay, in seconds, at the time that
     // the sample was emitted from the jitter buffer. To get the average target delay,
     // divide by JitterBufferEmittedCount
     JitterBufferTargetDelay [float64](/builtin#float64) `json:"jitterBufferTargetDelay"`
    
     // JitterBufferEmittedCount is the total number of audio samples or video frames that
     // have come out of the jitter buffer (increasing jitterBufferDelay).
     JitterBufferEmittedCount [uint64](/builtin#uint64) `json:"jitterBufferEmittedCount"`
    
     // JitterBufferMinimumDelay works the same way as jitterBufferTargetDelay, except that
     // it is not affected by external mechanisms that increase the jitter buffer target delay,
     // such as  jitterBufferTarget, AV sync, or any other mechanisms. This metric is purely
     // based on the network characteristics such as jitter and packet loss, and can be seen
     // as the minimum obtainable jitter  buffer delay if no external factors would affect it.
     // The metric is updated every time JitterBufferEmittedCount is updated.
     JitterBufferMinimumDelay [float64](/builtin#float64) `json:"jitterBufferMinimumDelay"`
    
     // TotalSamplesReceived is the total number of samples that have been received on
     // this RTP stream. This includes concealedSamples. Does not exist for video.
     TotalSamplesReceived [uint64](/builtin#uint64) `json:"totalSamplesReceived"`
    
     // ConcealedSamples is the total number of samples that are concealed samples.
     // A concealed sample is a sample that was replaced with synthesized samples generated
     // locally before being played out. Examples of samples that have to be concealed are
     // samples from lost packets (reported in packetsLost) or samples from packets that
     // arrive too late to be played out (reported in packetsDiscarded). Does not exist for video.
     ConcealedSamples [uint64](/builtin#uint64) `json:"concealedSamples"`
    
     // SilentConcealedSamples is the total number of concealed samples inserted that
     // are "silent". Playing out silent samples results in silence or comfort noise.
     // This is a subset of concealedSamples. Does not exist for video.
     SilentConcealedSamples [uint64](/builtin#uint64) `json:"silentConcealedSamples"`
    
     // ConcealmentEvents increases every time a concealed sample is synthesized after
     // a non-concealed sample. That is, multiple consecutive concealed samples will increase
     // the concealedSamples count multiple times but is a single concealment event.
     // Does not exist for video.
     ConcealmentEvents [uint64](/builtin#uint64) `json:"concealmentEvents"`
    
     // InsertedSamplesForDeceleration is increased by the difference between the number of
     // samples received and the number of samples played out when playout is slowed down.
     // If playout is slowed down by inserting samples, this will be the number of inserted samples.
     // Does not exist for video.
     InsertedSamplesForDeceleration [uint64](/builtin#uint64) `json:"insertedSamplesForDeceleration"`
    
     // RemovedSamplesForAcceleration is increased by the difference between the number of
     // samples received and the number of samples played out when playout is sped up. If speedup
     // is achieved by removing samples, this will be the count of samples removed.
     // Does not exist for video.
     RemovedSamplesForAcceleration [uint64](/builtin#uint64) `json:"removedSamplesForAcceleration"`
    
     // AudioLevel represents the audio level of the receiving track..
     //
     // The value is a value between 0..1 (linear), where 1.0 represents 0 dBov,
     // 0 represents silence, and 0.5 represents approximately 6 dBSPL change in
     // the sound pressure level from 0 dBov. Does not exist for video.
     AudioLevel [float64](/builtin#float64) `json:"audioLevel"`
    
     // TotalAudioEnergy represents the audio energy of the receiving track. It is calculated
     // by duration * Math.pow(energy/maxEnergy, 2) for each audio sample received (and thus
     // counted by TotalSamplesReceived). Does not exist for video.
     TotalAudioEnergy [float64](/builtin#float64) `json:"totalAudioEnergy"`
    
     // TotalSamplesDuration represents the total duration in seconds of all samples that have been
     // received (and thus counted by TotalSamplesReceived). Can be used with totalAudioEnergy to
     // compute an average audio level over different intervals. Does not exist for video.
     TotalSamplesDuration [float64](/builtin#float64) `json:"totalSamplesDuration"`
    
     // SLICount counts the total number of Slice Loss Indication (SLI) packets received
     // by the sender. This metric is only valid for video and is sent by receiver.
     SLICount [uint32](/builtin#uint32) `json:"sliCount"`
    
     // QPSum is the sum of the QP values of frames passed. The count of frames is
     // in FramesDecoded for inbound stream stats, and in FramesEncoded for outbound stream stats.
     QPSum [uint64](/builtin#uint64) `json:"qpSum"`
    
     // TotalDecodeTime is the total number of seconds that have been spent decoding the FramesDecoded
     // frames of this stream. The average decode time can be calculated by dividing this value
     // with FramesDecoded. The time it takes to decode one frame is the time passed between
     // feeding the decoder a frame and the decoder returning decoded data for that frame.
     TotalDecodeTime [float64](/builtin#float64) `json:"totalDecodeTime"`
    
     // TotalInterFrameDelay is the sum of the interframe delays in seconds between consecutively
     // rendered frames, recorded just after a frame has been rendered. The interframe delay variance
     // be calculated from TotalInterFrameDelay, TotalSquaredInterFrameDelay, and FramesRendered according
     // to the formula: (TotalSquaredInterFrameDelay - TotalInterFrameDelay^2 / FramesRendered) / FramesRendered.
     // Does not exist for audio.
     TotalInterFrameDelay [float64](/builtin#float64) `json:"totalInterFrameDelay"`
    
     // TotalSquaredInterFrameDelay is the sum of the squared interframe delays in seconds
     // between consecutively rendered frames, recorded just after a frame has been rendered.
     // See TotalInterFrameDelay for details on how to calculate the interframe delay variance.
     // Does not exist for audio.
     TotalSquaredInterFrameDelay [float64](/builtin#float64) `json:"totalSquaredInterFrameDelay"`
    
     // PacketsReceived is the total number of RTP packets received for this SSRC.
     PacketsReceived [uint32](/builtin#uint32) `json:"packetsReceived"`
    
     // PacketsLost is the total number of RTP packets lost for this SSRC. Note that
     // because of how this is estimated, it can be negative if more packets are received than sent.
     PacketsLost [int32](/builtin#int32) `json:"packetsLost"`
    
     // Jitter is the packet jitter measured in seconds for this SSRC
     Jitter [float64](/builtin#float64) `json:"jitter"`
    
     // PacketsDiscarded is the cumulative number of RTP packets discarded by the jitter
     // buffer due to late or early-arrival, i.e., these packets are not played out.
     // RTP packets discarded due to packet duplication are not reported in this metric.
     PacketsDiscarded [uint32](/builtin#uint32) `json:"packetsDiscarded"`
    
     // PacketsRepaired is the cumulative number of lost RTP packets repaired after applying
     // an error-resilience mechanism. It is measured for the primary source RTP packets
     // and only counted for RTP packets that have no further chance of repair.
     PacketsRepaired [uint32](/builtin#uint32) `json:"packetsRepaired"`
    
     // BurstPacketsLost is the cumulative number of RTP packets lost during loss bursts.
     BurstPacketsLost [uint32](/builtin#uint32) `json:"burstPacketsLost"`
    
     // BurstPacketsDiscarded is the cumulative number of RTP packets discarded during discard bursts.
     BurstPacketsDiscarded [uint32](/builtin#uint32) `json:"burstPacketsDiscarded"`
    
     // BurstLossCount is the cumulative number of bursts of lost RTP packets.
     BurstLossCount [uint32](/builtin#uint32) `json:"burstLossCount"`
    
     // BurstDiscardCount is the cumulative number of bursts of discarded RTP packets.
     BurstDiscardCount [uint32](/builtin#uint32) `json:"burstDiscardCount"`
    
     // BurstLossRate is the fraction of RTP packets lost during bursts to the
     // total number of RTP packets expected in the bursts.
     BurstLossRate [float64](/builtin#float64) `json:"burstLossRate"`
    
     // BurstDiscardRate is the fraction of RTP packets discarded during bursts to
     // the total number of RTP packets expected in bursts.
     BurstDiscardRate [float64](/builtin#float64) `json:"burstDiscardRate"`
    
     // GapLossRate is the fraction of RTP packets lost during the gap periods.
     GapLossRate [float64](/builtin#float64) `json:"gapLossRate"`
    
     // GapDiscardRate is the fraction of RTP packets discarded during the gap periods.
     GapDiscardRate [float64](/builtin#float64) `json:"gapDiscardRate"`
    
     // TrackID is the identifier of the stats object representing the receiving track,
     // a ReceiverAudioTrackAttachmentStats or ReceiverVideoTrackAttachmentStats.
     TrackID [string](/builtin#string) `json:"trackId"`
    
     // ReceiverID is the stats ID used to look up the AudioReceiverStats or VideoReceiverStats
     // object receiving this stream.
     ReceiverID [string](/builtin#string) `json:"receiverId"`
    
     // RemoteID is used for looking up the remote RemoteOutboundRTPStreamStats object
     // for the same SSRC.
     RemoteID [string](/builtin#string) `json:"remoteId"`
    
     // FramesDecoded represents the total number of frames correctly decoded for this SSRC,
     // i.e., frames that would be displayed if no frames are dropped. Only valid for video.
     FramesDecoded [uint32](/builtin#uint32) `json:"framesDecoded"`
    
     // KeyFramesDecoded represents the total number of key frames, such as key frames in
     // VP8 [RFC6386] or IDR-frames in H.264 [RFC6184], successfully decoded for this RTP
     // media stream. This is a subset of FramesDecoded. FramesDecoded - KeyFramesDecoded
     // gives you the number of delta frames decoded. Does not exist for audio.
     KeyFramesDecoded [uint32](/builtin#uint32) `json:"keyFramesDecoded"`
    
     // FramesRendered represents the total number of frames that have been rendered.
     // It is incremented just after a frame has been rendered. Does not exist for audio.
     FramesRendered [uint32](/builtin#uint32) `json:"framesRendered"`
    
     // FramesDropped is the total number of frames dropped prior to decode or dropped
     // because the frame missed its display deadline for this receiver's track.
     // The measurement begins when the receiver is created and is a cumulative metric
     // as defined in Appendix A (g) of [RFC7004]. Does not exist for audio.
     FramesDropped [uint32](/builtin#uint32) `json:"framesDropped"`
    
     // FrameWidth represents the width of the last decoded frame. Before the first
     // frame is decoded this member does not exist. Does not exist for audio.
     FrameWidth [uint32](/builtin#uint32) `json:"frameWidth"`
    
     // FrameHeight represents the height of the last decoded frame. Before the first
     // frame is decoded this member does not exist. Does not exist for audio.
     FrameHeight [uint32](/builtin#uint32) `json:"frameHeight"`
    
     // LastPacketReceivedTimestamp represents the timestamp at which the last packet was
     // received for this SSRC. This differs from Timestamp, which represents the time
     // at which the statistics were generated by the local endpoint.
     LastPacketReceivedTimestamp StatsTimestamp `json:"lastPacketReceivedTimestamp"`
    
     // HeaderBytesReceived is the total number of RTP header and padding bytes received for this SSRC.
     // This includes retransmissions. This does not include the size of transport layer headers such
     // as IP or UDP. headerBytesReceived + bytesReceived equals the number of bytes received as
     // payload over the transport.
     HeaderBytesReceived [uint64](/builtin#uint64) `json:"headerBytesReceived"`
    
     // AverageRTCPInterval is the average RTCP interval between two consecutive compound RTCP packets.
     // This is calculated by the sending endpoint when sending compound RTCP reports.
     // Compound packets must contain at least a RTCP RR or SR packet and an SDES packet
     // with the CNAME item.
     AverageRTCPInterval [float64](/builtin#float64) `json:"averageRtcpInterval"`
    
     // FECPacketsReceived is the total number of RTP FEC packets received for this SSRC.
     // This counter can also be incremented when receiving FEC packets in-band with media packets (e.g., with Opus).
     FECPacketsReceived [uint32](/builtin#uint32) `json:"fecPacketsReceived"`
    
     // FECPacketsDiscarded is the total number of RTP FEC packets received for this SSRC where the
     // error correction payload was discarded by the application. This may happen
     // 1. if all the source packets protected by the FEC packet were received or already
     // recovered by a separate FEC packet, or
     // 2. if the FEC packet arrived late, i.e., outside the recovery window, and the
     // lost RTP packets have already been skipped during playout.
     // This is a subset of FECPacketsReceived.
     FECPacketsDiscarded [uint64](/builtin#uint64) `json:"fecPacketsDiscarded"`
    
     // BytesReceived is the total number of bytes received for this SSRC.
     BytesReceived [uint64](/builtin#uint64) `json:"bytesReceived"`
    
     // FramesReceived represents the total number of complete frames received on this RTP stream.
     // This metric is incremented when the complete frame is received. Does not exist for audio.
     FramesReceived [uint32](/builtin#uint32) `json:"framesReceived"`
    
     // PacketsFailedDecryption is the cumulative number of RTP packets that failed
     // to be decrypted. These packets are not counted by PacketsDiscarded.
     PacketsFailedDecryption [uint32](/builtin#uint32) `json:"packetsFailedDecryption"`
    
     // PacketsDuplicated is the cumulative number of packets discarded because they
     // are duplicated. Duplicate packets are not counted in PacketsDiscarded.
     //
     // Duplicated packets have the same RTP sequence number and content as a previously
     // received packet. If multiple duplicates of a packet are received, all of them are counted.
     // An improved estimate of lost packets can be calculated by adding PacketsDuplicated to PacketsLost.
     PacketsDuplicated [uint32](/builtin#uint32) `json:"packetsDuplicated"`
    
     // PerDSCPPacketsReceived is the total number of packets received for this SSRC,
     // per Differentiated Services code point (DSCP) [RFC2474]. DSCPs are identified
     // as decimal integers in string form. Note that due to network remapping and bleaching,
     // these numbers are not expected to match the numbers seen on sending. Not all
     // OSes make this information available.
     PerDSCPPacketsReceived map[[string](/builtin#string)][uint32](/builtin#uint32) `json:"perDscpPacketsReceived"`
    
     // Identifies the decoder implementation used. This is useful for diagnosing interoperability issues.
     // Does not exist for audio.
     DecoderImplementation [string](/builtin#string) `json:"decoderImplementation"`
    
     // PauseCount is the total number of video pauses experienced by this receiver.
     // Video is considered to be paused if time passed since last rendered frame exceeds 5 seconds.
     // PauseCount is incremented when a frame is rendered after such a pause. Does not exist for audio.
     PauseCount [uint32](/builtin#uint32) `json:"pauseCount"`
    
     // TotalPausesDuration is the total duration of pauses (for definition of pause see PauseCount), in seconds.
     // Does not exist for audio.
     TotalPausesDuration [float64](/builtin#float64) `json:"totalPausesDuration"`
    
     // FreezeCount is the total number of video freezes experienced by this receiver.
     // It is a freeze if frame duration, which is time interval between two consecutively rendered frames,
     // is equal or exceeds Max(3 * avg_frame_duration_ms, avg_frame_duration_ms + 150),
     // where avg_frame_duration_ms is linear average of durations of last 30 rendered frames.
     // Does not exist for audio.
     FreezeCount [uint32](/builtin#uint32) `json:"freezeCount"`
    
     // TotalFreezesDuration is the total duration of rendered frames which are considered as frozen
     // (for definition of freeze see freezeCount), in seconds. Does not exist for audio.
     TotalFreezesDuration [float64](/builtin#float64) `json:"totalFreezesDuration"`
    
     // PowerEfficientDecoder indicates whether the decoder currently used is considered power efficient
     // by the user agent. Does not exist for audio.
     PowerEfficientDecoder [bool](/builtin#bool) `json:"powerEfficientDecoder"`
    }

InboundRTPStreamStats contains statistics for an inbound RTP stream that is currently received with this PeerConnection object.

#### type [MediaEngine](https://github.com/pion/webrtc/blob/v4.2.3/mediaengine.go#L33) ¶

    type MediaEngine struct {
     // contains filtered or unexported fields
    }

A MediaEngine defines the codecs supported by a PeerConnection, and the configuration of those codecs.

#### func (*MediaEngine) [RegisterCodec](https://github.com/pion/webrtc/blob/v4.2.3/mediaengine.go#L258) ¶

    func (m *MediaEngine) RegisterCodec(codec RTPCodecParameters, typ RTPCodecType) [error](/builtin#error)

RegisterCodec adds codec to the MediaEngine These are the list of codecs supported by this PeerConnection.

#### func (*MediaEngine) [RegisterDefaultCodecs](https://github.com/pion/webrtc/blob/v4.2.3/mediaengine.go#L65) ¶

    func (m *MediaEngine) RegisterDefaultCodecs() [error](/builtin#error)

RegisterDefaultCodecs registers the default codecs supported by Pion WebRTC. RegisterDefaultCodecs is not safe for concurrent use.

#### func (*MediaEngine) [RegisterFeedback](https://github.com/pion/webrtc/blob/v4.2.3/mediaengine.go#L327) ¶

    func (m *MediaEngine) RegisterFeedback(feedback RTCPFeedback, typ RTPCodecType)

RegisterFeedback adds feedback mechanism to already registered codecs.

#### func (*MediaEngine) [RegisterHeaderExtension](https://github.com/pion/webrtc/blob/v4.2.3/mediaengine.go#L280) ¶

    func (m *MediaEngine) RegisterHeaderExtension(
     extension RTPHeaderExtensionCapability,
     typ RTPCodecType,
     allowedDirections ...RTPTransceiverDirection,
    ) [error](/builtin#error)

RegisterHeaderExtension adds a header extension to the MediaEngine To determine the negotiated value use `GetHeaderExtensionID` after signaling is complete.

#### type [MediaKind](https://github.com/pion/webrtc/blob/v4.2.3/stats.go#L143) ¶

    type MediaKind [string](/builtin#string)

MediaKind indicates the kind of media (audio or video).

    const (
     // MediaKindAudio indicates this is audio stats.
     MediaKindAudio MediaKind = "audio"
     // MediaKindVideo indicates this is video stats.
     MediaKindVideo MediaKind = "video"
    )

#### type [MediaStreamStats](https://github.com/pion/webrtc/blob/v4.2.3/stats.go#L1469) ¶

    type MediaStreamStats struct {
     // Timestamp is the timestamp associated with this object.
     Timestamp StatsTimestamp `json:"timestamp"`
    
     // Type is the object's StatsType
     Type StatsType `json:"type"`
    
     // ID is a unique id that is associated with the component inspected to produce
     // this Stats object. Two Stats objects will have the same ID if they were produced
     // by inspecting the same underlying object.
     ID [string](/builtin#string) `json:"id"`
    
     // StreamIdentifier is the "id" property of the MediaStream
     StreamIdentifier [string](/builtin#string) `json:"streamIdentifier"`
    
     // TrackIDs is a list of the identifiers of the stats object representing the
     // stream's tracks, either ReceiverAudioTrackAttachmentStats or ReceiverVideoTrackAttachmentStats.
     TrackIDs [][string](/builtin#string) `json:"trackIds"`
    }

MediaStreamStats contains statistics related to a specific MediaStream.

#### type [NetworkType](https://github.com/pion/webrtc/blob/v4.2.3/networktype.go#L22) ¶

    type NetworkType [int](/builtin#int)

NetworkType represents the type of network.

    const (
     // NetworkTypeUnknown is the enum's zero-value.
     NetworkTypeUnknown NetworkType = [iota](/builtin#iota)
    
     // NetworkTypeUDP4 indicates UDP over IPv4.
     NetworkTypeUDP4
    
     // NetworkTypeUDP6 indicates UDP over IPv6.
     NetworkTypeUDP6
    
     // NetworkTypeTCP4 indicates TCP over IPv4.
     NetworkTypeTCP4
    
     // NetworkTypeTCP6 indicates TCP over IPv6.
     NetworkTypeTCP6
    )

#### func [NewNetworkType](https://github.com/pion/webrtc/blob/v4.2.3/networktype.go#L82) ¶

    func NewNetworkType(raw [string](/builtin#string)) (NetworkType, [error](/builtin#error))

NewNetworkType allows create network type from string It will be useful for getting custom network types from external config.

#### func (NetworkType) [Protocol](https://github.com/pion/webrtc/blob/v4.2.3/networktype.go#L65) ¶

    func (t NetworkType) Protocol() [string](/builtin#string)

Protocol returns udp or tcp.

#### func (NetworkType) [String](https://github.com/pion/webrtc/blob/v4.2.3/networktype.go#L49) ¶

    func (t NetworkType) String() [string](/builtin#string)

#### type [NominationValueGenerator](https://github.com/pion/webrtc/blob/v4.2.3/settingengine.go#L128) ¶ added in v4.2.0

    type NominationValueGenerator func() [uint32](/builtin#uint32)

NominationValueGenerator generates nomination values for ICE renomination.

#### type [OAuthCredential](https://github.com/pion/webrtc/blob/v4.2.3/oauthcredential.go#L10) ¶

    type OAuthCredential struct {
     // MACKey is a base64-url encoded format. It is used in STUN message
     // integrity hash calculation.
     MACKey [string](/builtin#string)
    
     // AccessToken is a base64-encoded format. This is an encrypted
     // self-contained token that is opaque to the application.
     AccessToken [string](/builtin#string)
    }

OAuthCredential represents OAuth credential information which is used by the STUN/TURN client to connect to an ICE server as defined in <https://tools.ietf.org/html/rfc7635>. Note that the kid parameter is not located in OAuthCredential, but in ICEServer's username member.

#### type [OfferAnswerOptions](https://github.com/pion/webrtc/blob/v4.2.3/offeransweroptions.go#L8) ¶

    type OfferAnswerOptions struct {
     // VoiceActivityDetection allows the application to provide information
     // about whether it wishes voice detection feature to be enabled or disabled.
     VoiceActivityDetection [bool](/builtin#bool)
     // ICETricklingSupported indicates whether the ICE agent should use trickle ICE
     // If set, the "a=ice-options:trickle" attribute is added to the generated SDP payload.
     // (See <https://datatracker.ietf.org/doc/html/rfc9725#section-4.3.3>)
     ICETricklingSupported [bool](/builtin#bool)
    }

OfferAnswerOptions is a base structure which describes the options that can be used to control the offer/answer creation process.

#### type [OfferOptions](https://github.com/pion/webrtc/blob/v4.2.3/offeransweroptions.go#L26) ¶

    type OfferOptions struct {
     OfferAnswerOptions
    
     // ICERestart forces the underlying ice gathering process to be restarted.
     // When this value is true, the generated description will have ICE
     // credentials that are different from the current credentials
     ICERestart [bool](/builtin#bool)
    }

OfferOptions structure describes the options used to control the offer creation process.

#### type [OutboundRTPStreamStats](https://github.com/pion/webrtc/blob/v4.2.3/stats.go#L638) ¶

    type OutboundRTPStreamStats struct {
     // Mid represents a mid value of RTPTransceiver owning this stream, if that value is not
     // null. Otherwise, this member is not present.
     Mid [string](/builtin#string) `json:"mid"`
    
     // Rid only exists if a rid has been set for this RTP stream.
     // Must not exist for audio.
     Rid [string](/builtin#string) `json:"rid"`
    
     // MediaSourceID is the identifier of the stats object representing the track currently
     // attached to the sender of this stream, an RTCMediaSourceStats.
     MediaSourceID [string](/builtin#string) `json:"mediaSourceId"`
    
     // Timestamp is the timestamp associated with this object.
     Timestamp StatsTimestamp `json:"timestamp"`
    
     // Type is the object's StatsType
     Type StatsType `json:"type"`
    
     // ID is a unique id that is associated with the component inspected to produce
     // this Stats object. Two Stats objects will have the same ID if they were produced
     // by inspecting the same underlying object.
     ID [string](/builtin#string) `json:"id"`
    
     // SSRC is the 32-bit unsigned integer value used to identify the source of the
     // stream of RTP packets that this stats object concerns.
     SSRC SSRC `json:"ssrc"`
    
     // Kind is either "audio" or "video"
     Kind [string](/builtin#string) `json:"kind"`
    
     // It is a unique identifier that is associated to the object that was inspected
     // to produce the TransportStats associated with this RTP stream.
     TransportID [string](/builtin#string) `json:"transportId"`
    
     // CodecID is a unique identifier that is associated to the object that was inspected
     // to produce the CodecStats associated with this RTP stream.
     CodecID [string](/builtin#string) `json:"codecId"`
    
     // HeaderBytesSent is the total number of RTP header and padding bytes sent for this SSRC. This does not
     // include the size of transport layer headers such as IP or UDP.
     // HeaderBytesSent + BytesSent equals the number of bytes sent as payload over the transport.
     HeaderBytesSent [uint64](/builtin#uint64) `json:"headerBytesSent"`
    
     // RetransmittedPacketsSent is the total number of packets that were retransmitted for this SSRC.
     // This is a subset of packetsSent. If RTX is not negotiated, retransmitted packets are sent
     // over this ssrc. If RTX was negotiated, retransmitted packets are sent over a separate SSRC
     // but is still accounted for here.
     RetransmittedPacketsSent [uint64](/builtin#uint64) `json:"retransmittedPacketsSent"`
    
     // RetransmittedBytesSent is the total number of bytes that were retransmitted for this SSRC,
     // only including payload bytes. This is a subset of bytesSent. If RTX is not negotiated,
     // retransmitted bytes are sent over this ssrc. If RTX was negotiated, retransmitted bytes
     // are sent over a separate SSRC but is still accounted for here.
     RetransmittedBytesSent [uint64](/builtin#uint64) `json:"retransmittedBytesSent"`
    
     // FIRCount counts the total number of Full Intra Request (FIR) packets received
     // by the sender. This metric is only valid for video and is sent by receiver.
     FIRCount [uint32](/builtin#uint32) `json:"firCount"`
    
     // PLICount counts the total number of Picture Loss Indication (PLI) packets
     // received by the sender. This metric is only valid for video and is sent by receiver.
     PLICount [uint32](/builtin#uint32) `json:"pliCount"`
    
     // NACKCount counts the total number of Negative ACKnowledgement (NACK) packets
     // received by the sender and is sent by receiver.
     NACKCount [uint32](/builtin#uint32) `json:"nackCount"`
    
     // SLICount counts the total number of Slice Loss Indication (SLI) packets received
     // by the sender. This metric is only valid for video and is sent by receiver.
     SLICount [uint32](/builtin#uint32) `json:"sliCount"`
    
     // QPSum is the sum of the QP values of frames passed. The count of frames is
     // in FramesDecoded for inbound stream stats, and in FramesEncoded for outbound stream stats.
     QPSum [uint64](/builtin#uint64) `json:"qpSum"`
    
     // PacketsSent is the total number of RTP packets sent for this SSRC.
     PacketsSent [uint32](/builtin#uint32) `json:"packetsSent"`
    
     // PacketsDiscardedOnSend is the total number of RTP packets for this SSRC that
     // have been discarded due to socket errors, i.e. a socket error occurred when handing
     // the packets to the socket. This might happen due to various reasons, including
     // full buffer or no available memory.
     PacketsDiscardedOnSend [uint32](/builtin#uint32) `json:"packetsDiscardedOnSend"`
    
     // FECPacketsSent is the total number of RTP FEC packets sent for this SSRC.
     // This counter can also be incremented when sending FEC packets in-band with
     // media packets (e.g., with Opus).
     FECPacketsSent [uint32](/builtin#uint32) `json:"fecPacketsSent"`
    
     // BytesSent is the total number of bytes sent for this SSRC.
     BytesSent [uint64](/builtin#uint64) `json:"bytesSent"`
    
     // BytesDiscardedOnSend is the total number of bytes for this SSRC that have
     // been discarded due to socket errors, i.e. a socket error occurred when handing
     // the packets containing the bytes to the socket. This might happen due to various
     // reasons, including full buffer or no available memory.
     BytesDiscardedOnSend [uint64](/builtin#uint64) `json:"bytesDiscardedOnSend"`
    
     // TrackID is the identifier of the stats object representing the current track
     // attachment to the sender of this stream, a SenderAudioTrackAttachmentStats
     // or SenderVideoTrackAttachmentStats.
     TrackID [string](/builtin#string) `json:"trackId"`
    
     // SenderID is the stats ID used to look up the AudioSenderStats or VideoSenderStats
     // object sending this stream.
     SenderID [string](/builtin#string) `json:"senderId"`
    
     // RemoteID is used for looking up the remote RemoteInboundRTPStreamStats object
     // for the same SSRC.
     RemoteID [string](/builtin#string) `json:"remoteId"`
    
     // LastPacketSentTimestamp represents the timestamp at which the last packet was
     // sent for this SSRC. This differs from timestamp, which represents the time at
     // which the statistics were generated by the local endpoint.
     LastPacketSentTimestamp StatsTimestamp `json:"lastPacketSentTimestamp"`
    
     // TargetBitrate is the current target bitrate configured for this particular SSRC
     // and is the Transport Independent Application Specific (TIAS) bitrate [RFC3890].
     // Typically, the target bitrate is a configuration parameter provided to the codec's
     // encoder and does not count the size of the IP or other transport layers like TCP or UDP.
     // It is measured in bits per second and the bitrate is calculated over a 1 second window.
     TargetBitrate [float64](/builtin#float64) `json:"targetBitrate"`
    
     // TotalEncodedBytesTarget is increased by the target frame size in bytes every time
     // a frame has been encoded. The actual frame size may be bigger or smaller than this number.
     // This value goes up every time framesEncoded goes up.
     TotalEncodedBytesTarget [uint64](/builtin#uint64) `json:"totalEncodedBytesTarget"`
    
     // FrameWidth represents the width of the last encoded frame. The resolution of the
     // encoded frame may be lower than the media source. Before the first frame is encoded
     // this member does not exist. Does not exist for audio.
     FrameWidth [uint32](/builtin#uint32) `json:"frameWidth"`
    
     // FrameHeight represents the height of the last encoded frame. The resolution of the
     // encoded frame may be lower than the media source. Before the first frame is encoded
     // this member does not exist. Does not exist for audio.
     FrameHeight [uint32](/builtin#uint32) `json:"frameHeight"`
    
     // FramesPerSecond is the number of encoded frames during the last second. This may be
     // lower than the media source frame rate. Does not exist for audio.
     FramesPerSecond [float64](/builtin#float64) `json:"framesPerSecond"`
    
     // FramesSent represents the total number of frames sent on this RTP stream. Does not exist for audio.
     FramesSent [uint32](/builtin#uint32) `json:"framesSent"`
    
     // HugeFramesSent represents the total number of huge frames sent by this RTP stream.
     // Huge frames, by definition, are frames that have an encoded size at least 2.5 times
     // the average size of the frames. The average size of the frames is defined as the
     // target bitrate per second divided by the target FPS at the time the frame was encoded.
     // These are usually complex to encode frames with a lot of changes in the picture.
     // This can be used to estimate, e.g slide changes in the streamed presentation.
     // Does not exist for audio.
     HugeFramesSent [uint32](/builtin#uint32) `json:"hugeFramesSent"`
    
     // FramesEncoded represents the total number of frames successfully encoded for this RTP media stream.
     // Only valid for video.
     FramesEncoded [uint32](/builtin#uint32) `json:"framesEncoded"`
    
     // KeyFramesEncoded represents the total number of key frames, such as key frames in VP8 [RFC6386] or
     // IDR-frames in H.264 [RFC6184], successfully encoded for this RTP media stream. This is a subset of
     // FramesEncoded. FramesEncoded - KeyFramesEncoded gives you the number of delta frames encoded.
     // Does not exist for audio.
     KeyFramesEncoded [uint32](/builtin#uint32) `json:"keyFramesEncoded"`
    
     // TotalEncodeTime is the total number of seconds that has been spent encoding the
     // framesEncoded frames of this stream. The average encode time can be calculated by
     // dividing this value with FramesEncoded. The time it takes to encode one frame is the
     // time passed between feeding the encoder a frame and the encoder returning encoded data
     // for that frame. This does not include any additional time it may take to packetize the resulting data.
     TotalEncodeTime [float64](/builtin#float64) `json:"totalEncodeTime"`
    
     // TotalPacketSendDelay is the total number of seconds that packets have spent buffered
     // locally before being transmitted onto the network. The time is measured from when
     // a packet is emitted from the RTP packetizer until it is handed over to the OS network socket.
     // This measurement is added to totalPacketSendDelay when packetsSent is incremented.
     TotalPacketSendDelay [float64](/builtin#float64) `json:"totalPacketSendDelay"`
    
     // AverageRTCPInterval is the average RTCP interval between two consecutive compound RTCP
     // packets. This is calculated by the sending endpoint when sending compound RTCP reports.
     // Compound packets must contain at least a RTCP RR or SR packet and an SDES packet with the CNAME item.
     AverageRTCPInterval [float64](/builtin#float64) `json:"averageRtcpInterval"`
    
     // QualityLimitationReason is the current reason for limiting the resolution and/or framerate,
     // or "none" if not limited. Only valid for video.
     QualityLimitationReason QualityLimitationReason `json:"qualityLimitationReason"`
    
     // QualityLimitationDurations is record of the total time, in seconds, that this
     // stream has spent in each quality limitation state. The record includes a mapping
     // for all QualityLimitationReason types, including "none". Only valid for video.
     QualityLimitationDurations map[[string](/builtin#string)][float64](/builtin#float64) `json:"qualityLimitationDurations"`
    
     // QualityLimitationResolutionChanges is the number of times that the resolution has changed
     // because we are quality limited (qualityLimitationReason has a value other than "none").
     // The counter is initially zero and increases when the resolution goes up or down.
     // For example, if a 720p track is sent as 480p for some time and then recovers to 720p,
     // qualityLimitationResolutionChanges will have the value 2. Does not exist for audio.
     QualityLimitationResolutionChanges [uint32](/builtin#uint32) `json:"qualityLimitationResolutionChanges"`
    
     // PerDSCPPacketsSent is the total number of packets sent for this SSRC, per DSCP.
     // DSCPs are identified as decimal integers in string form.
     PerDSCPPacketsSent map[[string](/builtin#string)][uint32](/builtin#uint32) `json:"perDscpPacketsSent"`
    
     // Active indicates whether this RTP stream is configured to be sent or disabled. Note that an
     // active stream can still not be sending, e.g. when being limited by network conditions.
     Active [bool](/builtin#bool) `json:"active"`
    
     // Identifies the encoder implementation used. This is useful for diagnosing interoperability issues.
     // Does not exist for audio.
     EncoderImplementation [string](/builtin#string) `json:"encoderImplementation"`
    
     // PowerEfficientEncoder indicates whether the encoder currently used is considered power efficient.
     // by the user agent. Does not exist for audio.
     PowerEfficientEncoder [bool](/builtin#bool) `json:"powerEfficientEncoder"`
    
     // ScalabilityMode identifies the layering mode used for video encoding. Does not exist for audio.
     ScalabilityMode [string](/builtin#string) `json:"scalabilityMode"`
    }

OutboundRTPStreamStats contains statistics for an outbound RTP stream that is currently sent with this PeerConnection object.

#### type [PayloadType](https://github.com/pion/webrtc/blob/v4.2.3/webrtc.go#L20) ¶

    type PayloadType [uint8](/builtin#uint8)

PayloadType identifies the format of the RTP payload and determines its interpretation by the application. Each codec in a RTP Session will have a different PayloadType

<https://tools.ietf.org/html/rfc3550#section-3>

#### type [PeerConnection](https://github.com/pion/webrtc/blob/v4.2.3/peerconnection.go#L36) ¶

    type PeerConnection struct {
     // contains filtered or unexported fields
    }

PeerConnection represents a WebRTC connection that establishes a peer-to-peer communications with another PeerConnection instance in a browser, or to another endpoint implementing the required protocols.

#### func [NewPeerConnection](https://github.com/pion/webrtc/blob/v4.2.3/peerconnection.go#L105) ¶

    func NewPeerConnection(configuration Configuration) (*PeerConnection, [error](/builtin#error))

NewPeerConnection creates a PeerConnection with the default codecs and interceptors.

If you wish to customize the set of available codecs and/or the set of active interceptors, create an API with a custom MediaEngine and/or interceptor.Registry, then call [(*API).NewPeerConnection] instead of this function.

#### func (*PeerConnection) [AddICECandidate](https://github.com/pion/webrtc/blob/v4.2.3/peerconnection.go#L2040) ¶

    func (pc *PeerConnection) AddICECandidate(candidate ICECandidateInit) [error](/builtin#error)

AddICECandidate accepts an ICE candidate string and adds it to the existing set of candidates.

#### func (*PeerConnection) [AddTrack](https://github.com/pion/webrtc/blob/v4.2.3/peerconnection.go#L2151) ¶

    func (pc *PeerConnection) AddTrack(track TrackLocal) (*RTPSender, [error](/builtin#error))

AddTrack adds a Track to the PeerConnection.

#### func (*PeerConnection) [AddTransceiverFromKind](https://github.com/pion/webrtc/blob/v4.2.3/peerconnection.go#L2254) ¶

    func (pc *PeerConnection) AddTransceiverFromKind(
     kind RTPCodecType,
     init ...RTPTransceiverInit,
    ) (t *RTPTransceiver, err [error](/builtin#error))

AddTransceiverFromKind Create a new RtpTransceiver and adds it to the set of transceivers.

#### func (*PeerConnection) [AddTransceiverFromTrack](https://github.com/pion/webrtc/blob/v4.2.3/peerconnection.go#L2299) ¶

    func (pc *PeerConnection) AddTransceiverFromTrack(
     track TrackLocal,
     init ...RTPTransceiverInit,
    ) (t *RTPTransceiver, err [error](/builtin#error))

AddTransceiverFromTrack Create a new RtpTransceiver(SendRecv or SendOnly) and add it to the set of transceivers.

#### func (*PeerConnection) [CanTrickleICECandidates](https://github.com/pion/webrtc/blob/v4.2.3/peerconnection.go#L2628) ¶ added in v4.1.7

    func (pc *PeerConnection) CanTrickleICECandidates() ICETrickleCapability

CanTrickleICECandidates reports whether the remote endpoint indicated support for receiving trickled ICE candidates.

#### func (*PeerConnection) [Close](https://github.com/pion/webrtc/blob/v4.2.3/peerconnection.go#L2429) ¶

    func (pc *PeerConnection) Close() [error](/builtin#error)

Close ends the PeerConnection.

#### func (*PeerConnection) [ConnectionState](https://github.com/pion/webrtc/blob/v4.2.3/peerconnection.go#L2660) ¶

    func (pc *PeerConnection) ConnectionState() PeerConnectionState

ConnectionState attribute returns the connection state of the PeerConnection instance.

#### func (*PeerConnection) [CreateAnswer](https://github.com/pion/webrtc/blob/v4.2.3/peerconnection.go#L865) ¶

    func (pc *PeerConnection) CreateAnswer(options *AnswerOptions) (SessionDescription, [error](/builtin#error))

CreateAnswer starts the PeerConnection and generates the localDescription.

#### func (*PeerConnection) [CreateDataChannel](https://github.com/pion/webrtc/blob/v4.2.3/peerconnection.go#L2329) ¶

    func (pc *PeerConnection) CreateDataChannel(label [string](/builtin#string), options *DataChannelInit) (*DataChannel, [error](/builtin#error))

CreateDataChannel creates a new DataChannel object with the given label and optional DataChannelInit used to configure properties of the underlying channel such as data reliability.

#### func (*PeerConnection) [CreateOffer](https://github.com/pion/webrtc/blob/v4.2.3/peerconnection.go#L639) ¶

    func (pc *PeerConnection) CreateOffer(options *OfferOptions) (SessionDescription, [error](/builtin#error))

CreateOffer starts the PeerConnection and generates the localDescription <https://w3c.github.io/webrtc-pc/#dom-rtcpeerconnection-createoffer>

#### func (*PeerConnection) [CurrentLocalDescription](https://github.com/pion/webrtc/blob/v4.2.3/peerconnection.go#L2577) ¶

    func (pc *PeerConnection) CurrentLocalDescription() *SessionDescription

CurrentLocalDescription represents the local description that was successfully negotiated the last time the PeerConnection transitioned into the stable state plus any local candidates that have been generated by the ICEAgent since the offer or answer was created.

#### func (*PeerConnection) [CurrentRemoteDescription](https://github.com/pion/webrtc/blob/v4.2.3/peerconnection.go#L2607) ¶

    func (pc *PeerConnection) CurrentRemoteDescription() *SessionDescription

CurrentRemoteDescription represents the last remote description that was successfully negotiated the last time the PeerConnection transitioned into the stable state plus any remote candidates that have been supplied via AddICECandidate() since the offer or answer was created.

#### func (*PeerConnection) [GetConfiguration](https://github.com/pion/webrtc/blob/v4.2.3/peerconnection.go#L607) ¶

    func (pc *PeerConnection) GetConfiguration() Configuration

GetConfiguration returns a Configuration object representing the current configuration of this PeerConnection object. The returned object is a copy and direct mutation on it will not take affect until SetConfiguration has been called with Configuration passed as its only argument. <https://www.w3.org/TR/webrtc/#dom-rtcpeerconnection-getconfiguration>

#### func (*PeerConnection) [GetReceivers](https://github.com/pion/webrtc/blob/v4.2.3/peerconnection.go#L2127) ¶

    func (pc *PeerConnection) GetReceivers() (receivers []*RTPReceiver)

GetReceivers returns the RTPReceivers that are currently attached to this PeerConnection.

#### func (*PeerConnection) [GetSenders](https://github.com/pion/webrtc/blob/v4.2.3/peerconnection.go#L2113) ¶

    func (pc *PeerConnection) GetSenders() (result []*RTPSender)

GetSenders returns the RTPSender that are currently attached to this PeerConnection.

#### func (*PeerConnection) [GetStats](https://github.com/pion/webrtc/blob/v4.2.3/peerconnection.go#L2669) ¶

    func (pc *PeerConnection) GetStats() StatsReport

GetStats return data providing statistics about the overall connection.

#### func (*PeerConnection) [GetTransceivers](https://github.com/pion/webrtc/blob/v4.2.3/peerconnection.go#L2141) ¶

    func (pc *PeerConnection) GetTransceivers() []*RTPTransceiver

GetTransceivers returns the RtpTransceiver that are currently attached to this PeerConnection.

#### func (*PeerConnection) [GracefulClose](https://github.com/pion/webrtc/blob/v4.2.3/peerconnection.go#L2436) ¶

    func (pc *PeerConnection) GracefulClose() [error](/builtin#error)

GracefulClose ends the PeerConnection. It also waits for any goroutines it started to complete. This is only safe to call outside of PeerConnection callbacks or if in a callback, in its own goroutine.

#### func (*PeerConnection) [ICEConnectionState](https://github.com/pion/webrtc/blob/v4.2.3/peerconnection.go#L2104) ¶

    func (pc *PeerConnection) ICEConnectionState() ICEConnectionState

ICEConnectionState returns the ICE connection state of the PeerConnection instance.

#### func (*PeerConnection) [ICEGatheringState](https://github.com/pion/webrtc/blob/v4.2.3/peerconnection.go#L2643) ¶

    func (pc *PeerConnection) ICEGatheringState() ICEGatheringState

ICEGatheringState attribute returns the ICE gathering state of the PeerConnection instance.

#### func (*PeerConnection) [ID](https://github.com/pion/webrtc/blob/v4.2.3/peerconnection.go#L611) ¶ added in v4.1.7

    func (pc *PeerConnection) ID() [string](/builtin#string)

#### func (*PeerConnection) [LocalDescription](https://github.com/pion/webrtc/blob/v4.2.3/peerconnection.go#L1119) ¶

    func (pc *PeerConnection) LocalDescription() *SessionDescription

LocalDescription returns PendingLocalDescription if it is not null and otherwise it returns CurrentLocalDescription. This property is used to determine if SetLocalDescription has already been called. <https://www.w3.org/TR/webrtc/#dom-rtcpeerconnection-localdescription>

#### func (*PeerConnection) [OnConnectionStateChange](https://github.com/pion/webrtc/blob/v4.2.3/peerconnection.go#L509) ¶

    func (pc *PeerConnection) OnConnectionStateChange(f func(PeerConnectionState))

OnConnectionStateChange sets an event handler which is called when the PeerConnectionState has changed.

#### func (*PeerConnection) [OnDataChannel](https://github.com/pion/webrtc/blob/v4.2.3/peerconnection.go#L288) ¶

    func (pc *PeerConnection) OnDataChannel(f func(*DataChannel))

OnDataChannel sets an event handler which is invoked when a data channel message arrives from a remote peer.

#### func (*PeerConnection) [OnICECandidate](https://github.com/pion/webrtc/blob/v4.2.3/peerconnection.go#L450) ¶

    func (pc *PeerConnection) OnICECandidate(f func(*ICECandidate))

OnICECandidate sets an event handler which is invoked when a new ICE candidate is found. ICE candidate gathering only begins when SetLocalDescription or SetRemoteDescription is called. Take note that the handler will be called with a nil pointer when gathering is finished.

#### func (*PeerConnection) [OnICEConnectionStateChange](https://github.com/pion/webrtc/blob/v4.2.3/peerconnection.go#L495) ¶

    func (pc *PeerConnection) OnICEConnectionStateChange(f func(ICEConnectionState))

OnICEConnectionStateChange sets an event handler which is called when an ICE connection state is changed.

#### func (*PeerConnection) [OnICEGatheringStateChange](https://github.com/pion/webrtc/blob/v4.2.3/peerconnection.go#L456) ¶

    func (pc *PeerConnection) OnICEGatheringStateChange(f func(ICEGatheringState))

OnICEGatheringStateChange sets an event handler which is invoked when the ICE candidate gathering state has changed.

#### func (*PeerConnection) [OnNegotiationNeeded](https://github.com/pion/webrtc/blob/v4.2.3/peerconnection.go#L296) ¶

    func (pc *PeerConnection) OnNegotiationNeeded(f func())

OnNegotiationNeeded sets an event handler which is invoked when a change has occurred which requires session negotiation.

#### func (*PeerConnection) [OnSignalingStateChange](https://github.com/pion/webrtc/blob/v4.2.3/peerconnection.go#L269) ¶

    func (pc *PeerConnection) OnSignalingStateChange(f func(SignalingState))

OnSignalingStateChange sets an event handler which is invoked when the peer connection's signaling state changes.

#### func (*PeerConnection) [OnTrack](https://github.com/pion/webrtc/blob/v4.2.3/peerconnection.go#L472) ¶

    func (pc *PeerConnection) OnTrack(f func(*TrackRemote, *RTPReceiver))

OnTrack sets an event handler which is called when remote track arrives from a remote peer.

#### func (*PeerConnection) [PendingLocalDescription](https://github.com/pion/webrtc/blob/v4.2.3/peerconnection.go#L2592) ¶

    func (pc *PeerConnection) PendingLocalDescription() *SessionDescription

PendingLocalDescription represents a local description that is in the process of being negotiated plus any local candidates that have been generated by the ICEAgent since the offer or answer was created. If the PeerConnection is in the stable state, the value is null.

#### func (*PeerConnection) [PendingRemoteDescription](https://github.com/pion/webrtc/blob/v4.2.3/peerconnection.go#L2619) ¶

    func (pc *PeerConnection) PendingRemoteDescription() *SessionDescription

PendingRemoteDescription represents a remote description that is in the process of being negotiated, complete with any remote candidates that have been supplied via AddICECandidate() since the offer or answer was created. If the PeerConnection is in the stable state, the value is null.

#### func (*PeerConnection) [RemoteDescription](https://github.com/pion/webrtc/blob/v4.2.3/peerconnection.go#L2027) ¶

    func (pc *PeerConnection) RemoteDescription() *SessionDescription

RemoteDescription returns pendingRemoteDescription if it is not null and otherwise it returns currentRemoteDescription. This property is used to determine if setRemoteDescription has already been called. <https://www.w3.org/TR/webrtc/#dom-rtcpeerconnection-remotedescription>

#### func (*PeerConnection) [RemoveTrack](https://github.com/pion/webrtc/blob/v4.2.3/peerconnection.go#L2189) ¶

    func (pc *PeerConnection) RemoveTrack(sender *RTPSender) (err [error](/builtin#error))

RemoveTrack removes a Track from the PeerConnection.

#### func (*PeerConnection) [SCTP](https://github.com/pion/webrtc/blob/v4.2.3/peerconnection.go#L3071) ¶

    func (pc *PeerConnection) SCTP() *SCTPTransport

SCTP returns the SCTPTransport for this PeerConnection

The SCTP transport over which SCTP data is sent and received. If SCTP has not been negotiated, the value is nil. <https://www.w3.org/TR/webrtc/#attributes-15>

#### func (*PeerConnection) [SetConfiguration](https://github.com/pion/webrtc/blob/v4.2.3/peerconnection.go#L523) ¶

    func (pc *PeerConnection) SetConfiguration(configuration Configuration) [error](/builtin#error)

SetConfiguration updates the configuration of this PeerConnection object. <https://www.w3.org/TR/webrtc/#dom-rtcpeerconnection-setconfiguration>

#### func (*PeerConnection) [SetIdentityProvider](https://github.com/pion/webrtc/blob/v4.2.3/peerconnection.go#L2412) ¶

    func (pc *PeerConnection) SetIdentityProvider([string](/builtin#string)) [error](/builtin#error)

SetIdentityProvider is used to configure an identity provider to generate identity assertions.

#### func (*PeerConnection) [SetLocalDescription](https://github.com/pion/webrtc/blob/v4.2.3/peerconnection.go#L1059) ¶

    func (pc *PeerConnection) SetLocalDescription(desc SessionDescription) [error](/builtin#error)

SetLocalDescription sets the SessionDescription of the local peer

#### func (*PeerConnection) [SetRemoteDescription](https://github.com/pion/webrtc/blob/v4.2.3/peerconnection.go#L1130) ¶

    func (pc *PeerConnection) SetRemoteDescription(desc SessionDescription) [error](/builtin#error)

SetRemoteDescription sets the SessionDescription of the remote peer

#### func (*PeerConnection) [SignalingState](https://github.com/pion/webrtc/blob/v4.2.3/peerconnection.go#L2637) ¶

    func (pc *PeerConnection) SignalingState() SignalingState

SignalingState attribute returns the signaling state of the PeerConnection instance.

#### func (*PeerConnection) [WriteRTCP](https://github.com/pion/webrtc/blob/v4.2.3/peerconnection.go#L2418) ¶

    func (pc *PeerConnection) WriteRTCP(pkts [][rtcp](/github.com/pion/rtcp).[Packet](/github.com/pion/rtcp#Packet)) [error](/builtin#error)

WriteRTCP sends a user provided RTCP packet to the connected peer. If no peer is connected the packet is discarded. It also runs any configured interceptors.

#### type [PeerConnectionState](https://github.com/pion/webrtc/blob/v4.2.3/peerconnectionstate.go#L7) ¶

    type PeerConnectionState [int](/builtin#int)

PeerConnectionState indicates the state of the PeerConnection.

    const (
     // PeerConnectionStateUnknown is the enum's zero-value.
     PeerConnectionStateUnknown PeerConnectionState = [iota](/builtin#iota)
    
     // PeerConnectionStateNew indicates that any of the ICETransports or
     // DTLSTransports are in the "new" state and none of the transports are
     // in the "connecting", "checking", "failed" or "disconnected" state, or
     // all transports are in the "closed" state, or there are no transports.
     PeerConnectionStateNew
    
     // PeerConnectionStateConnecting indicates that any of the
     // ICETransports or DTLSTransports are in the "connecting" or
     // "checking" state and none of them is in the "failed" state.
     PeerConnectionStateConnecting
    
     // PeerConnectionStateConnected indicates that all ICETransports and
     // DTLSTransports are in the "connected", "completed" or "closed" state
     // and at least one of them is in the "connected" or "completed" state.
     PeerConnectionStateConnected
    
     // PeerConnectionStateDisconnected indicates that any of the
     // ICETransports or DTLSTransports are in the "disconnected" state
     // and none of them are in the "failed" or "connecting" or "checking" state.
     PeerConnectionStateDisconnected
    
     // PeerConnectionStateFailed indicates that any of the ICETransports
     // or DTLSTransports are in a "failed" state.
     PeerConnectionStateFailed
    
     // PeerConnectionStateClosed indicates the peer connection is closed
     // and the isClosed member variable of PeerConnection is true.
     PeerConnectionStateClosed
    )

#### func (PeerConnectionState) [String](https://github.com/pion/webrtc/blob/v4.2.3/peerconnectionstate.go#L72) ¶

    func (t PeerConnectionState) String() [string](/builtin#string)

#### type [PeerConnectionStats](https://github.com/pion/webrtc/blob/v4.2.3/stats.go#L1368) ¶

    type PeerConnectionStats struct {
     // Timestamp is the timestamp associated with this object.
     Timestamp StatsTimestamp `json:"timestamp"`
    
     // Type is the object's StatsType
     Type StatsType `json:"type"`
    
     // ID is a unique id that is associated with the component inspected to produce
     // this Stats object. Two Stats objects will have the same ID if they were produced
     // by inspecting the same underlying object.
     ID [string](/builtin#string) `json:"id"`
    
     // DataChannelsOpened represents the number of unique DataChannels that have
     // entered the "open" state during their lifetime.
     DataChannelsOpened [uint32](/builtin#uint32) `json:"dataChannelsOpened"`
    
     // DataChannelsClosed represents the number of unique DataChannels that have
     // left the "open" state during their lifetime (due to being closed by either
     // end or the underlying transport being closed). DataChannels that transition
     // from "connecting" to "closing" or "closed" without ever being "open"
     // are not counted in this number.
     DataChannelsClosed [uint32](/builtin#uint32) `json:"dataChannelsClosed"`
    
     // DataChannelsRequested Represents the number of unique DataChannels returned
     // from a successful createDataChannel() call on the PeerConnection. If the
     // underlying data transport is not established, these may be in the "connecting" state.
     DataChannelsRequested [uint32](/builtin#uint32) `json:"dataChannelsRequested"`
    
     // DataChannelsAccepted represents the number of unique DataChannels signaled
     // in a "datachannel" event on the PeerConnection.
     DataChannelsAccepted [uint32](/builtin#uint32) `json:"dataChannelsAccepted"`
    }

PeerConnectionStats contains statistics related to the PeerConnection object.

#### type [QualityLimitationReason](https://github.com/pion/webrtc/blob/v4.2.3/stats.go#L616) ¶

    type QualityLimitationReason [string](/builtin#string)

QualityLimitationReason lists the reason for limiting the resolution and/or framerate. Only valid for video.

    const (
     // QualityLimitationReasonNone means the resolution and/or framerate is not limited.
     QualityLimitationReasonNone QualityLimitationReason = "none"
    
     // QualityLimitationReasonCPU means the resolution and/or framerate is primarily limited due to CPU load.
     QualityLimitationReasonCPU QualityLimitationReason = "cpu"
    
     // QualityLimitationReasonBandwidth means the resolution and/or framerate is primarily limited
     // due to congestion cues during bandwidth estimation.
     // Typical, congestion control algorithms use inter-arrival time, round-trip time,
     //  packet or other congestion cues to perform bandwidth estimation.
     QualityLimitationReasonBandwidth QualityLimitationReason = "bandwidth"
    
     // QualityLimitationReasonOther means the resolution and/or framerate is primarily limited
     //  for a reason other than the above.
     QualityLimitationReasonOther QualityLimitationReason = "other"
    )

#### type [RTCPFeedback](https://github.com/pion/webrtc/blob/v4.2.3/rtcpfeedback.go#L25) ¶

    type RTCPFeedback struct {
     // Type is the type of feedback.
     // see: <https://draft.ortc.org/#dom-rtcrtcpfeedback>
     // valid: ack, ccm, nack, goog-remb, transport-cc
     Type [string](/builtin#string)
    
     // The parameter value depends on the type.
     // For example, type="nack" parameter="pli" will send Picture Loss Indicator packets.
     Parameter [string](/builtin#string)
    }

RTCPFeedback signals the connection to use additional RTCP packet types. <https://draft.ortc.org/#dom-rtcrtcpfeedback>

#### type [RTCPMuxPolicy](https://github.com/pion/webrtc/blob/v4.2.3/rtcpmuxpolicy.go#L12) ¶

    type RTCPMuxPolicy [int](/builtin#int)

RTCPMuxPolicy affects what ICE candidates are gathered to support non-multiplexed RTCP.

    const (
     // RTCPMuxPolicyUnknown is the enum's zero-value.
     RTCPMuxPolicyUnknown RTCPMuxPolicy = [iota](/builtin#iota)
    
     // RTCPMuxPolicyNegotiate indicates to gather ICE candidates for both
     // RTP and RTCP candidates. If the remote-endpoint is capable of
     // multiplexing RTCP, multiplex RTCP on the RTP candidates. If it is not,
     // use both the RTP and RTCP candidates separately.
     RTCPMuxPolicyNegotiate
    
     // RTCPMuxPolicyRequire indicates to gather ICE candidates only for
     // RTP and multiplex RTCP on the RTP candidates. If the remote endpoint is
     // not capable of rtcp-mux, session negotiation will fail.
     RTCPMuxPolicyRequire
    )

#### func (RTCPMuxPolicy) [MarshalJSON](https://github.com/pion/webrtc/blob/v4.2.3/rtcpmuxpolicy.go#L71) ¶

    func (t RTCPMuxPolicy) MarshalJSON() ([][byte](/builtin#byte), [error](/builtin#error))

MarshalJSON returns the JSON encoding.

#### func (RTCPMuxPolicy) [String](https://github.com/pion/webrtc/blob/v4.2.3/rtcpmuxpolicy.go#L47) ¶

    func (t RTCPMuxPolicy) String() [string](/builtin#string)

#### func (*RTCPMuxPolicy) [UnmarshalJSON](https://github.com/pion/webrtc/blob/v4.2.3/rtcpmuxpolicy.go#L59) ¶

    func (t *RTCPMuxPolicy) UnmarshalJSON(b [][byte](/builtin#byte)) [error](/builtin#error)

UnmarshalJSON parses the JSON-encoded data and stores the result.

#### type [RTPCapabilities](https://github.com/pion/webrtc/blob/v4.2.3/rtpcapabilities.go#L9) ¶

    type RTPCapabilities struct {
     Codecs           []RTPCodecCapability
     HeaderExtensions []RTPHeaderExtensionCapability
    }

RTPCapabilities represents the capabilities of a transceiver

<https://w3c.github.io/webrtc-pc/#rtcrtpcapabilities>

#### type [RTPCodecCapability](https://github.com/pion/webrtc/blob/v4.2.3/rtpcodec.go#L54) ¶

    type RTPCodecCapability struct {
     MimeType     [string](/builtin#string)
     ClockRate    [uint32](/builtin#uint32)
     Channels     [uint16](/builtin#uint16)
     SDPFmtpLine  [string](/builtin#string)
     RTCPFeedback []RTCPFeedback
    }

RTPCodecCapability provides information about codec capabilities.

<https://w3c.github.io/webrtc-pc/#dictionary-rtcrtpcodeccapability-members>

#### type [RTPCodecParameters](https://github.com/pion/webrtc/blob/v4.2.3/rtpcodec.go#L82) ¶

    type RTPCodecParameters struct {
     RTPCodecCapability
     PayloadType PayloadType
     // contains filtered or unexported fields
    }

RTPCodecParameters is a sequence containing the media codecs that an RtpSender will choose from, as well as entries for RTX, RED and FEC mechanisms. This also includes the PayloadType that has been negotiated

<https://w3c.github.io/webrtc-pc/#rtcrtpcodecparameters>

#### type [RTPCodecType](https://github.com/pion/webrtc/blob/v4.2.3/rtpcodec.go#L15) ¶

    type RTPCodecType [int](/builtin#int)

RTPCodecType determines the type of a codec.

    const (
     // RTPCodecTypeUnknown is the enum's zero-value.
     RTPCodecTypeUnknown RTPCodecType = [iota](/builtin#iota)
    
     // RTPCodecTypeAudio indicates this is an audio codec.
     RTPCodecTypeAudio
    
     // RTPCodecTypeVideo indicates this is a video codec.
     RTPCodecTypeVideo
    )

#### func [NewRTPCodecType](https://github.com/pion/webrtc/blob/v4.2.3/rtpcodec.go#L40) ¶

    func NewRTPCodecType(r [string](/builtin#string)) RTPCodecType

NewRTPCodecType creates a RTPCodecType from a string.

#### func (RTPCodecType) [String](https://github.com/pion/webrtc/blob/v4.2.3/rtpcodec.go#L28) ¶

    func (t RTPCodecType) String() [string](/builtin#string)

#### type [RTPCodingParameters](https://github.com/pion/webrtc/blob/v4.2.3/rtpcodingparameters.go#L21) ¶

    type RTPCodingParameters struct {
     RID         [string](/builtin#string)           `json:"rid"`
     SSRC        SSRC             `json:"ssrc"`
     PayloadType PayloadType      `json:"payloadType"`
     RTX         RTPRtxParameters `json:"rtx"`
     FEC         RTPFecParameters `json:"fec"`
    }

RTPCodingParameters provides information relating to both encoding and decoding. This is a subset of the RFC since Pion WebRTC doesn't implement encoding/decoding itself <http://draft.ortc.org/#dom-rtcrtpcodingparameters>

#### type [RTPContributingSourceStats](https://github.com/pion/webrtc/blob/v4.2.3/stats.go#L1123) ¶

    type RTPContributingSourceStats struct {
     // Timestamp is the timestamp associated with this object.
     Timestamp StatsTimestamp `json:"timestamp"`
    
     // Type is the object's StatsType
     Type StatsType `json:"type"`
    
     // ID is a unique id that is associated with the component inspected to produce
     // this Stats object. Two Stats objects will have the same ID if they were produced
     // by inspecting the same underlying object.
     ID [string](/builtin#string) `json:"id"`
    
     // ContributorSSRC is the SSRC identifier of the contributing source represented
     // by this stats object. It is a 32-bit unsigned integer that appears in the CSRC
     // list of any packets the relevant source contributed to.
     ContributorSSRC SSRC `json:"contributorSsrc"`
    
     // InboundRTPStreamID is the ID of the InboundRTPStreamStats object representing
     // the inbound RTP stream that this contributing source is contributing to.
     InboundRTPStreamID [string](/builtin#string) `json:"inboundRtpStreamId"`
    
     // PacketsContributedTo is the total number of RTP packets that this contributing
     // source contributed to. This value is incremented each time a packet is counted
     // by InboundRTPStreamStats.packetsReceived, and the packet's CSRC list contains
     // the SSRC identifier of this contributing source, ContributorSSRC.
     PacketsContributedTo [uint32](/builtin#uint32) `json:"packetsContributedTo"`
    
     // AudioLevel is present if the last received RTP packet that this source contributed
     // to contained an [RFC6465] mixer-to-client audio level header extension. The value
     // of audioLevel is between 0..1 (linear), where 1.0 represents 0 dBov, 0 represents
     // silence, and 0.5 represents approximately 6 dBSPL change in the sound pressure level from 0 dBov.
     AudioLevel [float64](/builtin#float64) `json:"audioLevel"`
    }

RTPContributingSourceStats contains statistics for a contributing source (CSRC) that contributed to an inbound RTP stream.

#### type [RTPDecodingParameters](https://github.com/pion/webrtc/blob/v4.2.3/rtpdecodingparameters.go#L9) ¶

    type RTPDecodingParameters struct {
     RTPCodingParameters
    }

RTPDecodingParameters provides information relating to both encoding and decoding. This is a subset of the RFC since Pion WebRTC doesn't implement decoding itself <http://draft.ortc.org/#dom-rtcrtpdecodingparameters>

#### type [RTPEncodingParameters](https://github.com/pion/webrtc/blob/v4.2.3/rtpencodingparameters.go#L9) ¶

    type RTPEncodingParameters struct {
     RTPCodingParameters
    }

RTPEncodingParameters provides information relating to both encoding and decoding. This is a subset of the RFC since Pion WebRTC doesn't implement encoding itself <http://draft.ortc.org/#dom-rtcrtpencodingparameters>

#### type [RTPFecParameters](https://github.com/pion/webrtc/blob/v4.2.3/rtpcodingparameters.go#L14) ¶

    type RTPFecParameters struct {
     SSRC SSRC `json:"ssrc"`
    }

RTPFecParameters dictionary contains information relating to forward error correction (FEC) settings. <https://draft.ortc.org/#dom-rtcrtpfecparameters>

#### type [RTPHeaderExtensionCapability](https://github.com/pion/webrtc/blob/v4.2.3/rtpcodec.go#L65) ¶

    type RTPHeaderExtensionCapability struct {
     URI [string](/builtin#string)
    }

RTPHeaderExtensionCapability is used to define a RFC5285 RTP header extension supported by the codec.

<https://w3c.github.io/webrtc-pc/#dom-rtcrtpcapabilities-headerextensions>

#### type [RTPHeaderExtensionParameter](https://github.com/pion/webrtc/blob/v4.2.3/rtpcodec.go#L72) ¶

    type RTPHeaderExtensionParameter struct {
     URI [string](/builtin#string)
     ID  [int](/builtin#int)
    }

RTPHeaderExtensionParameter represents a negotiated RFC5285 RTP header extension.

<https://w3c.github.io/webrtc-pc/#dictionary-rtcrtpheaderextensionparameters-members>

#### type [RTPParameters](https://github.com/pion/webrtc/blob/v4.2.3/rtpcodec.go#L92) ¶

    type RTPParameters struct {
     HeaderExtensions []RTPHeaderExtensionParameter
     Codecs           []RTPCodecParameters
    }

RTPParameters is a list of negotiated codecs and header extensions

<https://w3c.github.io/webrtc-pc/#dictionary-rtcrtpparameters-members>

#### type [RTPReceiveParameters](https://github.com/pion/webrtc/blob/v4.2.3/rtpreceiveparameters.go#L7) ¶

    type RTPReceiveParameters struct {
     Encodings []RTPDecodingParameters
    }

RTPReceiveParameters contains the RTP stack settings used by receivers.

#### type [RTPReceiver](https://github.com/pion/webrtc/blob/v4.2.3/rtpreceiver.go#L62) ¶

    type RTPReceiver struct {
     // contains filtered or unexported fields
    }

RTPReceiver allows an application to inspect the receipt of a TrackRemote.

#### func (*RTPReceiver) [GetParameters](https://github.com/pion/webrtc/blob/v4.2.3/rtpreceiver.go#L133) ¶

    func (r *RTPReceiver) GetParameters() RTPParameters

GetParameters describes the current configuration for the encoding and transmission of media on the receiver's track.

#### func (*RTPReceiver) [RTPTransceiver](https://github.com/pion/webrtc/blob/v4.2.3/rtpreceiver.go#L168) ¶

    func (r *RTPReceiver) RTPTransceiver() *RTPTransceiver

RTPTransceiver returns the RTPTransceiver this RTPReceiver belongs too, or nil if none.

#### func (*RTPReceiver) [Read](https://github.com/pion/webrtc/blob/v4.2.3/rtpreceiver.go#L287) ¶

    func (r *RTPReceiver) Read(b [][byte](/builtin#byte)) (n [int](/builtin#int), a [interceptor](/github.com/pion/interceptor).[Attributes](/github.com/pion/interceptor#Attributes), err [error](/builtin#error))

Read reads incoming RTCP for this RTPReceiver.

#### func (*RTPReceiver) [ReadRTCP](https://github.com/pion/webrtc/blob/v4.2.3/rtpreceiver.go#L327) ¶

    func (r *RTPReceiver) ReadRTCP() ([][rtcp](/github.com/pion/rtcp).[Packet](/github.com/pion/rtcp#Packet), [interceptor](/github.com/pion/interceptor).[Attributes](/github.com/pion/interceptor#Attributes), [error](/builtin#error))

ReadRTCP is a convenience method that wraps Read and unmarshal for you. It also runs any configured interceptors.

#### func (*RTPReceiver) [ReadSimulcast](https://github.com/pion/webrtc/blob/v4.2.3/rtpreceiver.go#L301) ¶

    func (r *RTPReceiver) ReadSimulcast(b [][byte](/builtin#byte), rid [string](/builtin#string)) (n [int](/builtin#int), a [interceptor](/github.com/pion/interceptor).[Attributes](/github.com/pion/interceptor#Attributes), err [error](/builtin#error))

ReadSimulcast reads incoming RTCP for this RTPReceiver for given rid.

#### func (*RTPReceiver) [ReadSimulcastRTCP](https://github.com/pion/webrtc/blob/v4.2.3/rtpreceiver.go#L343) ¶

    func (r *RTPReceiver) ReadSimulcastRTCP(rid [string](/builtin#string)) ([][rtcp](/github.com/pion/rtcp).[Packet](/github.com/pion/rtcp#Packet), [interceptor](/github.com/pion/interceptor).[Attributes](/github.com/pion/interceptor#Attributes), [error](/builtin#error))

ReadSimulcastRTCP is a convenience method that wraps ReadSimulcast and unmarshal for you.

#### func (*RTPReceiver) [Receive](https://github.com/pion/webrtc/blob/v4.2.3/rtpreceiver.go#L280) ¶

    func (r *RTPReceiver) Receive(parameters RTPReceiveParameters) [error](/builtin#error)

Receive initialize the track and starts all the transports.

#### func (*RTPReceiver) [SetRTPParameters](https://github.com/pion/webrtc/blob/v4.2.3/rtpreceiver_go.go#L17) ¶

    func (r *RTPReceiver) SetRTPParameters(params RTPParameters)

SetRTPParameters applies provided RTPParameters the RTPReceiver's tracks.

This method is part of the ORTC API. It is not meant to be used together with the basic WebRTC API.

The amount of provided codecs must match the number of tracks on the receiver.

#### func (*RTPReceiver) [SetReadDeadline](https://github.com/pion/webrtc/blob/v4.2.3/rtpreceiver.go#L710) ¶

    func (r *RTPReceiver) SetReadDeadline(t [time](/time).[Time](/time#Time)) [error](/builtin#error)

SetReadDeadline sets the max amount of time the RTCP stream will block before returning. 0 is forever.

#### func (*RTPReceiver) [SetReadDeadlineSimulcast](https://github.com/pion/webrtc/blob/v4.2.3/rtpreceiver.go#L719) ¶

    func (r *RTPReceiver) SetReadDeadlineSimulcast(deadline [time](/time).[Time](/time#Time), rid [string](/builtin#string)) [error](/builtin#error)

SetReadDeadlineSimulcast sets the max amount of time the RTCP stream for a given rid will block before returning. 0 is forever.

#### func (*RTPReceiver) [Stop](https://github.com/pion/webrtc/blob/v4.2.3/rtpreceiver.go#L369) ¶

    func (r *RTPReceiver) Stop() [error](/builtin#error)

Stop irreversibly stops the RTPReceiver.

#### func (*RTPReceiver) [Track](https://github.com/pion/webrtc/blob/v4.2.3/rtpreceiver.go#L141) ¶

    func (r *RTPReceiver) Track() *TrackRemote

Track returns the RtpTransceiver TrackRemote.

#### func (*RTPReceiver) [Tracks](https://github.com/pion/webrtc/blob/v4.2.3/rtpreceiver.go#L154) ¶

    func (r *RTPReceiver) Tracks() []*TrackRemote

Tracks returns the RtpTransceiver tracks A RTPReceiver to support Simulcast may now have multiple tracks.

#### func (*RTPReceiver) [Transport](https://github.com/pion/webrtc/blob/v4.2.3/rtpreceiver.go#L112) ¶

    func (r *RTPReceiver) Transport() *DTLSTransport

Transport returns the currently-configured *DTLSTransport or nil if one has not yet been configured.

#### type [RTPRtxParameters](https://github.com/pion/webrtc/blob/v4.2.3/rtpcodingparameters.go#L8) ¶

    type RTPRtxParameters struct {
     SSRC SSRC `json:"ssrc"`
    }

RTPRtxParameters dictionary contains information relating to retransmission (RTX) settings. <https://draft.ortc.org/#dom-rtcrtprtxparameters>

#### type [RTPSendParameters](https://github.com/pion/webrtc/blob/v4.2.3/rtpsendparameters.go#L7) ¶

    type RTPSendParameters struct {
     RTPParameters
     Encodings []RTPEncodingParameters
    }

RTPSendParameters contains the RTP stack settings used by receivers.

#### type [RTPSender](https://github.com/pion/webrtc/blob/v4.2.3/rtpsender.go#L36) ¶

    type RTPSender struct {
     // contains filtered or unexported fields
    }

RTPSender allows an application to control how a given Track is encoded and transmitted to a remote peer.

#### func (*RTPSender) [AddEncoding](https://github.com/pion/webrtc/blob/v4.2.3/rtpsender.go#L154) ¶

    func (r *RTPSender) AddEncoding(track TrackLocal) [error](/builtin#error)

AddEncoding adds an encoding to RTPSender. Used by simulcast senders.

#### func (*RTPSender) [GetParameters](https://github.com/pion/webrtc/blob/v4.2.3/rtpsender.go#L117) ¶

    func (r *RTPSender) GetParameters() RTPSendParameters

GetParameters describes the current configuration for the encoding and transmission of media on the sender's track.

#### func (*RTPSender) [Read](https://github.com/pion/webrtc/blob/v4.2.3/rtpsender.go#L411) ¶

    func (r *RTPSender) Read(b [][byte](/builtin#byte)) (n [int](/builtin#int), a [interceptor](/github.com/pion/interceptor).[Attributes](/github.com/pion/interceptor#Attributes), err [error](/builtin#error))

Read reads incoming RTCP for this RTPSender.

#### func (*RTPSender) [ReadRTCP](https://github.com/pion/webrtc/blob/v4.2.3/rtpsender.go#L421) ¶

    func (r *RTPSender) ReadRTCP() ([][rtcp](/github.com/pion/rtcp).[Packet](/github.com/pion/rtcp#Packet), [interceptor](/github.com/pion/interceptor).[Attributes](/github.com/pion/interceptor#Attributes), [error](/builtin#error))

ReadRTCP is a convenience method that wraps Read and unmarshals for you.

#### func (*RTPSender) [ReadSimulcast](https://github.com/pion/webrtc/blob/v4.2.3/rtpsender.go#L437) ¶

    func (r *RTPSender) ReadSimulcast(b [][byte](/builtin#byte), rid [string](/builtin#string)) (n [int](/builtin#int), a [interceptor](/github.com/pion/interceptor).[Attributes](/github.com/pion/interceptor#Attributes), err [error](/builtin#error))

ReadSimulcast reads incoming RTCP for this RTPSender for given rid.

#### func (*RTPSender) [ReadSimulcastRTCP](https://github.com/pion/webrtc/blob/v4.2.3/rtpsender.go#L458) ¶

    func (r *RTPSender) ReadSimulcastRTCP(rid [string](/builtin#string)) ([][rtcp](/github.com/pion/rtcp).[Packet](/github.com/pion/rtcp#Packet), [interceptor](/github.com/pion/interceptor).[Attributes](/github.com/pion/interceptor#Attributes), [error](/builtin#error))

ReadSimulcastRTCP is a convenience method that wraps ReadSimulcast and unmarshal for you.

#### func (*RTPSender) [ReplaceTrack](https://github.com/pion/webrtc/blob/v4.2.3/rtpsender.go#L233) ¶

    func (r *RTPSender) ReplaceTrack(track TrackLocal) [error](/builtin#error)

ReplaceTrack replaces the track currently being used as the sender's source with a new TrackLocal. The new track must be of the same media kind (audio, video, etc) and switching the track should not require negotiation.

#### func (*RTPSender) [Send](https://github.com/pion/webrtc/blob/v4.2.3/rtpsender.go#L302) ¶

    func (r *RTPSender) Send(parameters RTPSendParameters) [error](/builtin#error)

Send Attempts to set the parameters controlling the sending of media.

#### func (*RTPSender) [SetReadDeadline](https://github.com/pion/webrtc/blob/v4.2.3/rtpsender.go#L472) ¶

    func (r *RTPSender) SetReadDeadline(t [time](/time).[Time](/time#Time)) [error](/builtin#error)

SetReadDeadline sets the deadline for the Read operation. Setting to zero means no deadline.

#### func (*RTPSender) [SetReadDeadlineSimulcast](https://github.com/pion/webrtc/blob/v4.2.3/rtpsender.go#L482) ¶

    func (r *RTPSender) SetReadDeadlineSimulcast(deadline [time](/time).[Time](/time#Time), rid [string](/builtin#string)) [error](/builtin#error)

SetReadDeadlineSimulcast sets the max amount of time the RTCP stream for a given rid will block before returning. 0 is forever.

#### func (*RTPSender) [Stop](https://github.com/pion/webrtc/blob/v4.2.3/rtpsender.go#L379) ¶

    func (r *RTPSender) Stop() [error](/builtin#error)

Stop irreversibly stops the RTPSender.

#### func (*RTPSender) [Track](https://github.com/pion/webrtc/blob/v4.2.3/rtpsender.go#L219) ¶

    func (r *RTPSender) Track() TrackLocal

Track returns the RTCRtpTransceiver track, or nil.

#### func (*RTPSender) [Transport](https://github.com/pion/webrtc/blob/v4.2.3/rtpsender.go#L108) ¶

    func (r *RTPSender) Transport() *DTLSTransport

Transport returns the currently-configured *DTLSTransport or nil if one has not yet been configured.

#### type [RTPTransceiver](https://github.com/pion/webrtc/blob/v4.2.3/rtptransceiver.go#L21) ¶

    type RTPTransceiver struct {
     // contains filtered or unexported fields
    }

RTPTransceiver represents a combination of an RTPSender and an RTPReceiver that share a common mid.

#### func (*RTPTransceiver) [Direction](https://github.com/pion/webrtc/blob/v4.2.3/rtptransceiver.go#L248) ¶

    func (t *RTPTransceiver) Direction() RTPTransceiverDirection

Direction returns the RTPTransceiver's current direction.

#### func (*RTPTransceiver) [Kind](https://github.com/pion/webrtc/blob/v4.2.3/rtptransceiver.go#L243) ¶

    func (t *RTPTransceiver) Kind() RTPCodecType

Kind returns RTPTransceiver's kind.

#### func (*RTPTransceiver) [Mid](https://github.com/pion/webrtc/blob/v4.2.3/rtptransceiver.go#L234) ¶

    func (t *RTPTransceiver) Mid() [string](/builtin#string)

Mid gets the Transceiver's mid value. When not already set, this value will be set in CreateOffer or CreateAnswer.

#### func (*RTPTransceiver) [Receiver](https://github.com/pion/webrtc/blob/v4.2.3/rtptransceiver.go#L215) ¶

    func (t *RTPTransceiver) Receiver() *RTPReceiver

Receiver returns the RTPTransceiver's RTPReceiver if it has one.

#### func (*RTPTransceiver) [Sender](https://github.com/pion/webrtc/blob/v4.2.3/rtptransceiver.go#L187) ¶

    func (t *RTPTransceiver) Sender() *RTPSender

Sender returns the RTPTransceiver's RTPSender if it has one.

#### func (*RTPTransceiver) [SetCodecPreferences](https://github.com/pion/webrtc/blob/v4.2.3/rtptransceiver.go#L55) ¶

    func (t *RTPTransceiver) SetCodecPreferences(codecs []RTPCodecParameters) [error](/builtin#error)

SetCodecPreferences sets preferred list of supported codecs if codecs is empty or nil we reset to default from MediaEngine.

#### func (*RTPTransceiver) [SetMid](https://github.com/pion/webrtc/blob/v4.2.3/rtptransceiver.go#L224) ¶

    func (t *RTPTransceiver) SetMid(mid [string](/builtin#string)) [error](/builtin#error)

SetMid sets the RTPTransceiver's mid. If it was already set, will return an error.

#### func (*RTPTransceiver) [SetSender](https://github.com/pion/webrtc/blob/v4.2.3/rtptransceiver.go#L196) ¶

    func (t *RTPTransceiver) SetSender(s *RTPSender, track TrackLocal) [error](/builtin#error)

SetSender sets the RTPSender and Track to current transceiver.

#### func (*RTPTransceiver) [Stop](https://github.com/pion/webrtc/blob/v4.2.3/rtptransceiver.go#L257) ¶

    func (t *RTPTransceiver) Stop() [error](/builtin#error)

Stop irreversibly stops the RTPTransceiver.

#### type [RTPTransceiverDirection](https://github.com/pion/webrtc/blob/v4.2.3/rtptransceiverdirection.go#L9) ¶

    type RTPTransceiverDirection [int](/builtin#int)

RTPTransceiverDirection indicates the direction of the RTPTransceiver.

    const (
     // RTPTransceiverDirectionUnknown is the enum's zero-value.
     RTPTransceiverDirectionUnknown RTPTransceiverDirection = [iota](/builtin#iota)
    
     // RTPTransceiverDirectionSendrecv indicates the RTPSender will offer
     // to send RTP and the RTPReceiver will offer to receive RTP.
     RTPTransceiverDirectionSendrecv
    
     // RTPTransceiverDirectionSendonly indicates the RTPSender will offer
     // to send RTP.
     RTPTransceiverDirectionSendonly
    
     // RTPTransceiverDirectionRecvonly indicates the RTPReceiver will
     // offer to receive RTP.
     RTPTransceiverDirectionRecvonly
    
     // RTPTransceiverDirectionInactive indicates the RTPSender won't offer
     // to send RTP and the RTPReceiver won't offer to receive RTP.
     RTPTransceiverDirectionInactive
    )

#### func [NewRTPTransceiverDirection](https://github.com/pion/webrtc/blob/v4.2.3/rtptransceiverdirection.go#L42) ¶

    func NewRTPTransceiverDirection(raw [string](/builtin#string)) RTPTransceiverDirection

NewRTPTransceiverDirection defines a procedure for creating a new RTPTransceiverDirection from a raw string naming the transceiver direction.

#### func (RTPTransceiverDirection) [Revers](https://github.com/pion/webrtc/blob/v4.2.3/rtptransceiverdirection.go#L73) ¶

    func (t RTPTransceiverDirection) Revers() RTPTransceiverDirection

Revers indicate the opposite direction.

#### func (RTPTransceiverDirection) [String](https://github.com/pion/webrtc/blob/v4.2.3/rtptransceiverdirection.go#L57) ¶

    func (t RTPTransceiverDirection) String() [string](/builtin#string)

#### type [RTPTransceiverInit](https://github.com/pion/webrtc/blob/v4.2.3/rtptransceiverinit.go#L8) ¶

    type RTPTransceiverInit struct {
     Direction     RTPTransceiverDirection
     SendEncodings []RTPEncodingParameters
    }

RTPTransceiverInit dictionary is used when calling the WebRTC function addTransceiver() to provide configuration options for the new transceiver.

#### type [RemoteInboundRTPStreamStats](https://github.com/pion/webrtc/blob/v4.2.3/stats.go#L873) ¶

    type RemoteInboundRTPStreamStats struct {
     // Timestamp is the timestamp associated with this object.
     Timestamp StatsTimestamp `json:"timestamp"`
    
     // Type is the object's StatsType
     Type StatsType `json:"type"`
    
     // ID is a unique id that is associated with the component inspected to produce
     // this Stats object. Two Stats objects will have the same ID if they were produced
     // by inspecting the same underlying object.
     ID [string](/builtin#string) `json:"id"`
    
     // SSRC is the 32-bit unsigned integer value used to identify the source of the
     // stream of RTP packets that this stats object concerns.
     SSRC SSRC `json:"ssrc"`
    
     // Kind is either "audio" or "video"
     Kind [string](/builtin#string) `json:"kind"`
    
     // It is a unique identifier that is associated to the object that was inspected
     // to produce the TransportStats associated with this RTP stream.
     TransportID [string](/builtin#string) `json:"transportId"`
    
     // CodecID is a unique identifier that is associated to the object that was inspected
     // to produce the CodecStats associated with this RTP stream.
     CodecID [string](/builtin#string) `json:"codecId"`
    
     // FIRCount counts the total number of Full Intra Request (FIR) packets received
     // by the sender. This metric is only valid for video and is sent by receiver.
     FIRCount [uint32](/builtin#uint32) `json:"firCount"`
    
     // PLICount counts the total number of Picture Loss Indication (PLI) packets
     // received by the sender. This metric is only valid for video and is sent by receiver.
     PLICount [uint32](/builtin#uint32) `json:"pliCount"`
    
     // NACKCount counts the total number of Negative ACKnowledgement (NACK) packets
     // received by the sender and is sent by receiver.
     NACKCount [uint32](/builtin#uint32) `json:"nackCount"`
    
     // SLICount counts the total number of Slice Loss Indication (SLI) packets received
     // by the sender. This metric is only valid for video and is sent by receiver.
     SLICount [uint32](/builtin#uint32) `json:"sliCount"`
    
     // QPSum is the sum of the QP values of frames passed. The count of frames is
     // in FramesDecoded for inbound stream stats, and in FramesEncoded for outbound stream stats.
     QPSum [uint64](/builtin#uint64) `json:"qpSum"`
    
     // PacketsReceived is the total number of RTP packets received for this SSRC.
     PacketsReceived [uint32](/builtin#uint32) `json:"packetsReceived"`
    
     // PacketsLost is the total number of RTP packets lost for this SSRC. Note that
     // because of how this is estimated, it can be negative if more packets are received than sent.
     PacketsLost [int32](/builtin#int32) `json:"packetsLost"`
    
     // Jitter is the packet jitter measured in seconds for this SSRC
     Jitter [float64](/builtin#float64) `json:"jitter"`
    
     // PacketsDiscarded is the cumulative number of RTP packets discarded by the jitter
     // buffer due to late or early-arrival, i.e., these packets are not played out.
     // RTP packets discarded due to packet duplication are not reported in this metric.
     PacketsDiscarded [uint32](/builtin#uint32) `json:"packetsDiscarded"`
    
     // PacketsRepaired is the cumulative number of lost RTP packets repaired after applying
     // an error-resilience mechanism. It is measured for the primary source RTP packets
     // and only counted for RTP packets that have no further chance of repair.
     PacketsRepaired [uint32](/builtin#uint32) `json:"packetsRepaired"`
    
     // BurstPacketsLost is the cumulative number of RTP packets lost during loss bursts.
     BurstPacketsLost [uint32](/builtin#uint32) `json:"burstPacketsLost"`
    
     // BurstPacketsDiscarded is the cumulative number of RTP packets discarded during discard bursts.
     BurstPacketsDiscarded [uint32](/builtin#uint32) `json:"burstPacketsDiscarded"`
    
     // BurstLossCount is the cumulative number of bursts of lost RTP packets.
     BurstLossCount [uint32](/builtin#uint32) `json:"burstLossCount"`
    
     // BurstDiscardCount is the cumulative number of bursts of discarded RTP packets.
     BurstDiscardCount [uint32](/builtin#uint32) `json:"burstDiscardCount"`
    
     // BurstLossRate is the fraction of RTP packets lost during bursts to the
     // total number of RTP packets expected in the bursts.
     BurstLossRate [float64](/builtin#float64) `json:"burstLossRate"`
    
     // BurstDiscardRate is the fraction of RTP packets discarded during bursts to
     // the total number of RTP packets expected in bursts.
     BurstDiscardRate [float64](/builtin#float64) `json:"burstDiscardRate"`
    
     // GapLossRate is the fraction of RTP packets lost during the gap periods.
     GapLossRate [float64](/builtin#float64) `json:"gapLossRate"`
    
     // GapDiscardRate is the fraction of RTP packets discarded during the gap periods.
     GapDiscardRate [float64](/builtin#float64) `json:"gapDiscardRate"`
    
     // LocalID is used for looking up the local OutboundRTPStreamStats object for the same SSRC.
     LocalID [string](/builtin#string) `json:"localId"`
    
     // RoundTripTime is the estimated round trip time for this SSRC based on the
     // RTCP timestamps in the RTCP Receiver Report (RR) and measured in seconds.
     RoundTripTime [float64](/builtin#float64) `json:"roundTripTime"`
    
     // TotalRoundTripTime represents the cumulative sum of all round trip time measurements
     // in seconds since the beginning of the session. The individual round trip time is calculated
     // based on the RTCP timestamps in the RTCP Receiver Report (RR) [RFC3550], hence requires
     // a DLSR value other than 0. The average round trip time can be computed from
     // TotalRoundTripTime by dividing it by RoundTripTimeMeasurements.
     TotalRoundTripTime [float64](/builtin#float64) `json:"totalRoundTripTime"`
    
     // FractionLost is the fraction packet loss reported for this SSRC.
     FractionLost [float64](/builtin#float64) `json:"fractionLost"`
    
     // RoundTripTimeMeasurements represents the total number of RTCP RR blocks received for this SSRC
     // that contain a valid round trip time. This counter will not increment if the RoundTripTime can
     // not be calculated because no RTCP Receiver Report with a DLSR value other than 0 has been received.
     RoundTripTimeMeasurements [uint64](/builtin#uint64) `json:"roundTripTimeMeasurements"`
    }

RemoteInboundRTPStreamStats contains statistics for the remote endpoint's inbound RTP stream corresponding to an outbound stream that is currently sent with this PeerConnection object. It is measured at the remote endpoint and reported in an RTCP Receiver Report (RR) or RTCP Extended Report (XR).

#### type [RemoteOutboundRTPStreamStats](https://github.com/pion/webrtc/blob/v4.2.3/stats.go#L1005) ¶

    type RemoteOutboundRTPStreamStats struct {
     // Timestamp is the timestamp associated with this object.
     Timestamp StatsTimestamp `json:"timestamp"`
    
     // Type is the object's StatsType
     Type StatsType `json:"type"`
    
     // ID is a unique id that is associated with the component inspected to produce
     // this Stats object. Two Stats objects will have the same ID if they were produced
     // by inspecting the same underlying object.
     ID [string](/builtin#string) `json:"id"`
    
     // SSRC is the 32-bit unsigned integer value used to identify the source of the
     // stream of RTP packets that this stats object concerns.
     SSRC SSRC `json:"ssrc"`
    
     // Kind is either "audio" or "video"
     Kind [string](/builtin#string) `json:"kind"`
    
     // It is a unique identifier that is associated to the object that was inspected
     // to produce the TransportStats associated with this RTP stream.
     TransportID [string](/builtin#string) `json:"transportId"`
    
     // CodecID is a unique identifier that is associated to the object that was inspected
     // to produce the CodecStats associated with this RTP stream.
     CodecID [string](/builtin#string) `json:"codecId"`
    
     // FIRCount counts the total number of Full Intra Request (FIR) packets received
     // by the sender. This metric is only valid for video and is sent by receiver.
     FIRCount [uint32](/builtin#uint32) `json:"firCount"`
    
     // PLICount counts the total number of Picture Loss Indication (PLI) packets
     // received by the sender. This metric is only valid for video and is sent by receiver.
     PLICount [uint32](/builtin#uint32) `json:"pliCount"`
    
     // NACKCount counts the total number of Negative ACKnowledgement (NACK) packets
     // received by the sender and is sent by receiver.
     NACKCount [uint32](/builtin#uint32) `json:"nackCount"`
    
     // SLICount counts the total number of Slice Loss Indication (SLI) packets received
     // by the sender. This metric is only valid for video and is sent by receiver.
     SLICount [uint32](/builtin#uint32) `json:"sliCount"`
    
     // QPSum is the sum of the QP values of frames passed. The count of frames is
     // in FramesDecoded for inbound stream stats, and in FramesEncoded for outbound stream stats.
     QPSum [uint64](/builtin#uint64) `json:"qpSum"`
    
     // PacketsSent is the total number of RTP packets sent for this SSRC.
     PacketsSent [uint32](/builtin#uint32) `json:"packetsSent"`
    
     // PacketsDiscardedOnSend is the total number of RTP packets for this SSRC that
     // have been discarded due to socket errors, i.e. a socket error occurred when handing
     // the packets to the socket. This might happen due to various reasons, including
     // full buffer or no available memory.
     PacketsDiscardedOnSend [uint32](/builtin#uint32) `json:"packetsDiscardedOnSend"`
    
     // FECPacketsSent is the total number of RTP FEC packets sent for this SSRC.
     // This counter can also be incremented when sending FEC packets in-band with
     // media packets (e.g., with Opus).
     FECPacketsSent [uint32](/builtin#uint32) `json:"fecPacketsSent"`
    
     // BytesSent is the total number of bytes sent for this SSRC.
     BytesSent [uint64](/builtin#uint64) `json:"bytesSent"`
    
     // BytesDiscardedOnSend is the total number of bytes for this SSRC that have
     // been discarded due to socket errors, i.e. a socket error occurred when handing
     // the packets containing the bytes to the socket. This might happen due to various
     // reasons, including full buffer or no available memory.
     BytesDiscardedOnSend [uint64](/builtin#uint64) `json:"bytesDiscardedOnSend"`
    
     // LocalID is used for looking up the local InboundRTPStreamStats object for the same SSRC.
     LocalID [string](/builtin#string) `json:"localId"`
    
     // RemoteTimestamp represents the remote timestamp at which these statistics were
     // sent by the remote endpoint. This differs from timestamp, which represents the
     // time at which the statistics were generated or received by the local endpoint.
     // The RemoteTimestamp, if present, is derived from the NTP timestamp in an RTCP
     // Sender Report (SR) packet, which reflects the remote endpoint's clock.
     // That clock may not be synchronized with the local clock.
     RemoteTimestamp StatsTimestamp `json:"remoteTimestamp"`
    
     // ReportsSent represents the total number of RTCP Sender Report (SR) blocks sent for this SSRC.
     ReportsSent [uint64](/builtin#uint64) `json:"reportsSent"`
    
     // RoundTripTime is estimated round trip time for this SSRC based on the latest
     // RTCP Sender Report (SR) that contains a DLRR report block as defined in [RFC3611].
     // The Calculation of the round trip time is defined in section 4.5. of [RFC3611].
     // Does not exist if the latest SR does not contain the DLRR report block, or if the last RR timestamp
     // in the DLRR report block is zero, or if the delay since last RR value in the DLRR report block is zero.
     RoundTripTime [float64](/builtin#float64) `json:"roundTripTime"`
    
     // TotalRoundTripTime represents the cumulative sum of all round trip time measurements in seconds
     // since the beginning of the session. The individual round trip time is calculated based on the DLRR
     // report block in the RTCP Sender Report (SR) [RFC3611]. This counter will not increment if the
     // RoundTripTime can not be calculated. The average round trip time can be computed from
     // TotalRoundTripTime by dividing it by RoundTripTimeMeasurements.
     TotalRoundTripTime [float64](/builtin#float64) `json:"totalRoundTripTime"`
    
     // RoundTripTimeMeasurements represents the total number of RTCP Sender Report (SR) blocks
     // received for this SSRC that contain a DLRR report block that can derive a valid round trip time
     // according to [RFC3611]. This counter will not increment if the RoundTripTime can not be calculated.
     RoundTripTimeMeasurements [uint64](/builtin#uint64) `json:"roundTripTimeMeasurements"`
    }

RemoteOutboundRTPStreamStats contains statistics for the remote endpoint's outbound RTP stream corresponding to an inbound stream that is currently received with this PeerConnection object. It is measured at the remote endpoint and reported in an RTCP Sender Report (SR).

#### type [RenominationOption](https://github.com/pion/webrtc/blob/v4.2.3/settingengine.go#L135) ¶ added in v4.2.0

    type RenominationOption func(*renominationSettings)

RenominationOption allows configuring ICE renomination behavior.

#### func [WithRenominationGenerator](https://github.com/pion/webrtc/blob/v4.2.3/settingengine.go#L138) ¶ added in v4.2.0

    func WithRenominationGenerator(generator NominationValueGenerator) RenominationOption

WithRenominationGenerator overrides the default nomination value generator.

#### func [WithRenominationInterval](https://github.com/pion/webrtc/blob/v4.2.3/settingengine.go#L146) ¶ added in v4.2.0

    func WithRenominationInterval(interval [time](/time).[Duration](/time#Duration)) RenominationOption

WithRenominationInterval sets the interval for automatic renomination checks. Passing zero or a negative duration returns an error from SetICERenomination.

#### type [SCTPCapabilities](https://github.com/pion/webrtc/blob/v4.2.3/sctpcapabilities.go#L7) ¶

    type SCTPCapabilities struct {
     MaxMessageSize [uint32](/builtin#uint32) `json:"maxMessageSize"`
    }

SCTPCapabilities indicates the capabilities of the SCTPTransport.

#### type [SCTPTransport](https://github.com/pion/webrtc/blob/v4.2.3/sctptransport.go#L24) ¶

    type SCTPTransport struct {
     // contains filtered or unexported fields
    }

SCTPTransport provides details about the SCTP transport.

#### func (*SCTPTransport) [BufferedAmount](https://github.com/pion/webrtc/blob/v4.2.3/sctptransport.go#L450) ¶ added in v4.0.15

    func (r *SCTPTransport) BufferedAmount() [int](/builtin#int)

BufferedAmount returns total amount (in bytes) of currently buffered user data.

#### func (*SCTPTransport) [GetCapabilities](https://github.com/pion/webrtc/blob/v4.2.3/sctptransport.go#L86) ¶

    func (r *SCTPTransport) GetCapabilities() SCTPCapabilities

GetCapabilities returns the SCTPCapabilities of the SCTPTransport.

#### func (*SCTPTransport) [MaxChannels](https://github.com/pion/webrtc/blob/v4.2.3/sctptransport.go#L373) ¶

    func (r *SCTPTransport) MaxChannels() [uint16](/builtin#uint16)

MaxChannels is the maximum number of RTCDataChannels that can be open simultaneously.

#### func (*SCTPTransport) [OnClose](https://github.com/pion/webrtc/blob/v4.2.3/sctptransport.go#L304) ¶

    func (r *SCTPTransport) OnClose(f func(err [error](/builtin#error)))

OnClose sets an event handler which is invoked when the SCTP Association closes.

#### func (*SCTPTransport) [OnDataChannel](https://github.com/pion/webrtc/blob/v4.2.3/sctptransport.go#L322) ¶

    func (r *SCTPTransport) OnDataChannel(f func(*DataChannel))

OnDataChannel sets an event handler which is invoked when a data channel message arrives from a remote peer.

#### func (*SCTPTransport) [OnDataChannelOpened](https://github.com/pion/webrtc/blob/v4.2.3/sctptransport.go#L330) ¶

    func (r *SCTPTransport) OnDataChannelOpened(f func(*DataChannel))

OnDataChannelOpened sets an event handler which is invoked when a data channel is opened.

#### func (*SCTPTransport) [OnError](https://github.com/pion/webrtc/blob/v4.2.3/sctptransport.go#L287) ¶

    func (r *SCTPTransport) OnError(f func(err [error](/builtin#error)))

OnError sets an event handler which is invoked when the SCTP Association errors.

#### func (*SCTPTransport) [Start](https://github.com/pion/webrtc/blob/v4.2.3/sctptransport.go#L100) ¶

    func (r *SCTPTransport) Start(capabilities SCTPCapabilities) [error](/builtin#error)

Start the SCTPTransport. Since both local and remote parties must mutually create an SCTPTransport, SCTP SO (Simultaneous Open) is used to establish a connection over SCTP.

#### func (*SCTPTransport) [State](https://github.com/pion/webrtc/blob/v4.2.3/sctptransport.go#L385) ¶

    func (r *SCTPTransport) State() SCTPTransportState

State returns the current state of the SCTPTransport.

#### func (*SCTPTransport) [Stop](https://github.com/pion/webrtc/blob/v4.2.3/sctptransport.go#L161) ¶

    func (r *SCTPTransport) Stop() [error](/builtin#error)

Stop stops the SCTPTransport.

#### func (*SCTPTransport) [Transport](https://github.com/pion/webrtc/blob/v4.2.3/sctptransport.go#L78) ¶

    func (r *SCTPTransport) Transport() *DTLSTransport

Transport returns the DTLSTransport instance the SCTPTransport is sending over.

#### type [SCTPTransportState](https://github.com/pion/webrtc/blob/v4.2.3/sctptransportstate.go#L7) ¶

    type SCTPTransportState [int](/builtin#int)

SCTPTransportState indicates the state of the SCTP transport.

    const (
     // SCTPTransportStateUnknown is the enum's zero-value.
     SCTPTransportStateUnknown SCTPTransportState = [iota](/builtin#iota)
    
     // SCTPTransportStateConnecting indicates the SCTPTransport is in the
     // process of negotiating an association. This is the initial state of the
     // SCTPTransportState when an SCTPTransport is created.
     SCTPTransportStateConnecting
    
     // SCTPTransportStateConnected indicates the negotiation of an
     // association is completed.
     SCTPTransportStateConnected
    
     // SCTPTransportStateClosed indicates a SHUTDOWN or ABORT chunk is
     // received or when the SCTP association has been closed intentionally,
     // such as by closing the peer connection or applying a remote description
     // that rejects data or changes the SCTP port.
     SCTPTransportStateClosed
    )

#### func (SCTPTransportState) [String](https://github.com/pion/webrtc/blob/v4.2.3/sctptransportstate.go#L49) ¶

    func (s SCTPTransportState) String() [string](/builtin#string)

#### type [SCTPTransportStats](https://github.com/pion/webrtc/blob/v4.2.3/stats.go#L2395) ¶

    type SCTPTransportStats struct {
     // Timestamp is the timestamp associated with this object.
     Timestamp StatsTimestamp `json:"timestamp"`
    
     // Type is the object's StatsType
     Type StatsType `json:"type"`
    
     // ID is a unique id that is associated with the component inspected to produce
     // this Stats object. Two Stats objects will have the same ID if they were produced
     // by inspecting the same underlying object.
     ID [string](/builtin#string) `json:"id"`
    
     // TransportID is the identifier of the object that was inspected to produce the
     // RTCTransportStats for the DTLSTransport and ICETransport supporting the SCTP transport.
     TransportID [string](/builtin#string) `json:"transportId"`
    
     // SmoothedRoundTripTime is the latest smoothed round-trip time value,
     // corresponding to spinfo_srtt defined in [RFC6458] but converted to seconds.
     // If there has been no round-trip time measurements yet, this value is undefined.
     SmoothedRoundTripTime [float64](/builtin#float64) `json:"smoothedRoundTripTime"`
    
     // CongestionWindow is the latest congestion window, corresponding to spinfo_cwnd defined in [RFC6458].
     CongestionWindow [uint32](/builtin#uint32) `json:"congestionWindow"`
    
     // ReceiverWindow is the latest receiver window, corresponding to sstat_rwnd defined in [RFC6458].
     ReceiverWindow [uint32](/builtin#uint32) `json:"receiverWindow"`
    
     // MTU is the latest maximum transmission unit, corresponding to spinfo_mtu defined in [RFC6458].
     MTU [uint32](/builtin#uint32) `json:"mtu"`
    
     // UNACKData is the number of unacknowledged DATA chunks, corresponding to sstat_unackdata defined in [RFC6458].
     UNACKData [uint32](/builtin#uint32) `json:"unackData"`
    
     // BytesSent represents the total number of bytes sent on this SCTPTransport
     BytesSent [uint64](/builtin#uint64) `json:"bytesSent"`
    
     // BytesReceived represents the total number of bytes received on this SCTPTransport
     BytesReceived [uint64](/builtin#uint64) `json:"bytesReceived"`
    }

SCTPTransportStats contains information about a certificate used by an SCTPTransport.

#### type [SDPSemantics](https://github.com/pion/webrtc/blob/v4.2.3/sdpsemantics.go#L12) ¶

    type SDPSemantics [int](/builtin#int)

SDPSemantics determines which style of SDP offers and answers can be used.

    const (
     // SDPSemanticsUnifiedPlan uses unified-plan offers and answers
     // (the default in Chrome since M72)
     // <https://tools.ietf.org/html/draft-roach-mmusic-unified-plan-00>
     SDPSemanticsUnifiedPlan SDPSemantics = [iota](/builtin#iota)
    
     // SDPSemanticsPlanB uses plan-b offers and answers
     // NB: This format should be considered deprecated
     // <https://tools.ietf.org/html/draft-uberti-rtcweb-plan-00>
     SDPSemanticsPlanB
    
     // SDPSemanticsUnifiedPlanWithFallback prefers unified-plan
     // offers and answers, but will respond to a plan-b offer
     // with a plan-b answer.
     SDPSemanticsUnifiedPlanWithFallback
    )

#### func (SDPSemantics) [MarshalJSON](https://github.com/pion/webrtc/blob/v4.2.3/sdpsemantics.go#L74) ¶

    func (s SDPSemantics) MarshalJSON() ([][byte](/builtin#byte), [error](/builtin#error))

MarshalJSON returns the JSON encoding.

#### func (SDPSemantics) [String](https://github.com/pion/webrtc/blob/v4.2.3/sdpsemantics.go#L48) ¶

    func (s SDPSemantics) String() [string](/builtin#string)

#### func (*SDPSemantics) [UnmarshalJSON](https://github.com/pion/webrtc/blob/v4.2.3/sdpsemantics.go#L62) ¶

    func (s *SDPSemantics) UnmarshalJSON(b [][byte](/builtin#byte)) [error](/builtin#error)

UnmarshalJSON parses the JSON-encoded data and stores the result.

#### type [SDPType](https://github.com/pion/webrtc/blob/v4.2.3/sdptype.go#L12) ¶

    type SDPType [int](/builtin#int)

SDPType describes the type of an SessionDescription.

    const (
     // SDPTypeUnknown is the enum's zero-value.
     SDPTypeUnknown SDPType = [iota](/builtin#iota)
    
     // SDPTypeOffer indicates that a description MUST be treated as an SDP offer.
     SDPTypeOffer
    
     // SDPTypePranswer indicates that a description MUST be treated as an
     // SDP answer, but not a final answer. A description used as an SDP
     // pranswer may be applied as a response to an SDP offer, or an update to
     // a previously sent SDP pranswer.
     SDPTypePranswer
    
     // SDPTypeAnswer indicates that a description MUST be treated as an SDP
     // final answer, and the offer-answer exchange MUST be considered complete.
     // A description used as an SDP answer may be applied as a response to an
     // SDP offer or as an update to a previously sent SDP pranswer.
     SDPTypeAnswer
    
     // SDPTypeRollback indicates that a description MUST be treated as
     // canceling the current SDP negotiation and moving the SDP offer and
     // answer back to what it was in the previous stable state. Note the
     // local or remote SDP descriptions in the previous stable state could be
     // null if there has not yet been a successful offer-answer negotiation.
     SDPTypeRollback
    )

#### func [NewSDPType](https://github.com/pion/webrtc/blob/v4.2.3/sdptype.go#L50) ¶

    func NewSDPType(raw [string](/builtin#string)) SDPType

NewSDPType creates an SDPType from a string.

#### func (SDPType) [MarshalJSON](https://github.com/pion/webrtc/blob/v4.2.3/sdptype.go#L81) ¶

    func (t SDPType) MarshalJSON() ([][byte](/builtin#byte), [error](/builtin#error))

MarshalJSON enables JSON marshaling of a SDPType.

#### func (SDPType) [String](https://github.com/pion/webrtc/blob/v4.2.3/sdptype.go#L65) ¶

    func (t SDPType) String() [string](/builtin#string)

#### func (*SDPType) [UnmarshalJSON](https://github.com/pion/webrtc/blob/v4.2.3/sdptype.go#L86) ¶

    func (t *SDPType) UnmarshalJSON(b [][byte](/builtin#byte)) [error](/builtin#error)

UnmarshalJSON enables JSON unmarshaling of a SDPType.

#### type [SSRC](https://github.com/pion/webrtc/blob/v4.2.3/webrtc.go#L13) ¶

    type SSRC [uint32](/builtin#uint32)

SSRC represents a synchronization source A synchronization source is a randomly chosen value meant to be globally unique within a particular RTP session. Used to identify a single stream of media.

<https://tools.ietf.org/html/rfc3550#section-3>

#### type [SenderAudioTrackAttachmentStats](https://github.com/pion/webrtc/blob/v4.2.3/stats.go#L1589) ¶

    type SenderAudioTrackAttachmentStats AudioSenderStats

SenderAudioTrackAttachmentStats object represents the stats about one attachment of an audio MediaStreamTrack to the PeerConnection object for which one calls GetStats.

It appears in the stats as soon as it is attached (via AddTrack, via AddTransceiver, via ReplaceTrack on an RTPSender object).

If an audio track is attached twice (via AddTransceiver or ReplaceTrack), there will be two SenderAudioTrackAttachmentStats objects, one for each attachment. They will have the same "TrackIdentifier" attribute, but different "ID" attributes.

If the track is detached from the PeerConnection (via removeTrack or via replaceTrack), it continues to appear, but with the "ObjectDeleted" member set to true.

#### type [SenderVideoTrackAttachmentStats](https://github.com/pion/webrtc/blob/v4.2.3/stats.go#L1657) ¶

    type SenderVideoTrackAttachmentStats VideoSenderStats

SenderVideoTrackAttachmentStats represents the stats about one attachment of a video MediaStreamTrack to the PeerConnection object for which one calls GetStats.

It appears in the stats as soon as it is attached (via AddTrack, via AddTransceiver, via ReplaceTrack on an RTPSender object).

If a video track is attached twice (via AddTransceiver or ReplaceTrack), there will be two SenderVideoTrackAttachmentStats objects, one for each attachment. They will have the same "TrackIdentifier" attribute, but different "ID" attributes.

If the track is detached from the PeerConnection (via RemoveTrack or via ReplaceTrack), it continues to appear, but with the "ObjectDeleted" member set to true.

#### type [SessionDescription](https://github.com/pion/webrtc/blob/v4.2.3/sessiondescription.go#L40) ¶

    type SessionDescription struct {
     Type SDPType `json:"type"`
     SDP  [string](/builtin#string)  `json:"sdp"`
     // contains filtered or unexported fields
    }

SessionDescription is used to expose local and remote session descriptions.

#### func (*SessionDescription) [Unmarshal](https://github.com/pion/webrtc/blob/v4.2.3/sessiondescription.go#L49) ¶

    func (sd *SessionDescription) Unmarshal() (*[sdp](/github.com/pion/sdp/v3).[SessionDescription](/github.com/pion/sdp/v3#SessionDescription), [error](/builtin#error))

Unmarshal is a helper to deserialize the sdp.

#### type [SettingEngine](https://github.com/pion/webrtc/blob/v4.2.3/settingengine.go#L31) ¶

    type SettingEngine struct {
     BufferFactory func(packetType [packetio](/github.com/pion/transport/v4/packetio).[BufferPacketType](/github.com/pion/transport/v4/packetio#BufferPacketType), ssrc [uint32](/builtin#uint32)) [io](/io).[ReadWriteCloser](/io#ReadWriteCloser)
     LoggerFactory [logging](/github.com/pion/logging).[LoggerFactory](/github.com/pion/logging#LoggerFactory)
     // contains filtered or unexported fields
    }

SettingEngine allows influencing behavior in ways that are not supported by the WebRTC API. This allows us to support additional use-cases without deviating from the WebRTC API elsewhere.

#### func (*SettingEngine) [DetachDataChannels](https://github.com/pion/webrtc/blob/v4.2.3/settingengine.go#L202) ¶

    func (e *SettingEngine) DetachDataChannels()

DetachDataChannels enables detaching data channels. When enabled data channels have to be detached in the OnOpen callback using the DataChannel.Detach method.

#### func (*SettingEngine) [DisableActiveTCP](https://github.com/pion/webrtc/blob/v4.2.3/settingengine.go#L492) ¶

    func (e *SettingEngine) DisableActiveTCP(isDisabled [bool](/builtin#bool))

DisableActiveTCP disables using active TCP for ICE. Active TCP is enabled by default.

#### func (*SettingEngine) [DisableCertificateFingerprintVerification](https://github.com/pion/webrtc/blob/v4.2.3/settingengine.go#L428) ¶

    func (e *SettingEngine) DisableCertificateFingerprintVerification(isDisabled [bool](/builtin#bool))

DisableCertificateFingerprintVerification disables fingerprint verification after DTLS Handshake has finished.

#### func (*SettingEngine) [DisableCloseByDTLS](https://github.com/pion/webrtc/blob/v4.2.3/settingengine.go#L677) ¶ added in v4.0.5

    func (e *SettingEngine) DisableCloseByDTLS(isEnabled [bool](/builtin#bool))

DisableCloseByDTLS sets if the connection should be closed when dtls transport is closed. Setting this to true will keep the connection open when dtls transport is closed and relies on the ice failed state to detect the connection is interrupted.

#### func (*SettingEngine) [DisableMediaEngineCopy](https://github.com/pion/webrtc/blob/v4.2.3/settingengine.go#L499) ¶

    func (e *SettingEngine) DisableMediaEngineCopy(isDisabled [bool](/builtin#bool))

DisableMediaEngineCopy stops the MediaEngine from being copied. This allows a user to modify the MediaEngine after the PeerConnection has been constructed. This is useful if you wish to modify codecs after signaling. Make sure not to share MediaEngines between PeerConnections.

#### func (*SettingEngine) [DisableMediaEngineMultipleCodecs](https://github.com/pion/webrtc/blob/v4.2.3/settingengine.go#L509) ¶ added in v4.0.16

    func (e *SettingEngine) DisableMediaEngineMultipleCodecs(isDisabled [bool](/builtin#bool))

DisableMediaEngineMultipleCodecs disables the MediaEngine negotiating different codecs. With the default value multiple media sections in the SDP can each negotiate different codecs. This is the new default behvior, because it makes Pion more spec compliant. The value of this setting will get copied to every copy of the MediaEngine generated for new PeerConnections (assuming DisableMediaEngineCopy is set to false). Note: this setting is targeted to be removed in release 4.2.0 (or later).

#### func (*SettingEngine) [DisableSRTCPReplayProtection](https://github.com/pion/webrtc/blob/v4.2.3/settingengine.go#L455) ¶

    func (e *SettingEngine) DisableSRTCPReplayProtection(isDisabled [bool](/builtin#bool))

DisableSRTCPReplayProtection disables SRTCP replay protection.

#### func (*SettingEngine) [DisableSRTPReplayProtection](https://github.com/pion/webrtc/blob/v4.2.3/settingengine.go#L450) ¶

    func (e *SettingEngine) DisableSRTPReplayProtection(isDisabled [bool](/builtin#bool))

DisableSRTPReplayProtection disables SRTP replay protection.

#### func (*SettingEngine) [EnableDataChannelBlockWrite](https://github.com/pion/webrtc/blob/v4.2.3/settingengine.go#L208) ¶ added in v4.0.6

    func (e *SettingEngine) EnableDataChannelBlockWrite(nonblockWrite [bool](/builtin#bool))

EnableDataChannelBlockWrite allows data channels to block on write, it only works if DetachDataChannels is enabled.

#### func (*SettingEngine) [EnableSCTPZeroChecksum](https://github.com/pion/webrtc/blob/v4.2.3/settingengine.go#L589) ¶

    func (e *SettingEngine) EnableSCTPZeroChecksum(isEnabled [bool](/builtin#bool))

EnableSCTPZeroChecksum controls the zero checksum feature in SCTP. This removes the need to checksum every incoming/outgoing packet and will reduce latency and CPU usage. This feature is not backwards compatible so is disabled by default.

#### func (*SettingEngine) [SetAnsweringDTLSRole](https://github.com/pion/webrtc/blob/v4.2.3/settingengine.go#L387) ¶

    func (e *SettingEngine) SetAnsweringDTLSRole(role DTLSRole) [error](/builtin#error)

SetAnsweringDTLSRole sets the DTLS role that is selected when offering The DTLS role controls if the WebRTC Client as a client or server. This may be useful when interacting with non-compliant clients or debugging issues.

DTLSRoleActive:

    Act as DTLS Client, send the ClientHello and starts the handshake
    

DTLSRolePassive:

    Act as DTLS Server, wait for ClientHello
    

#### func (*SettingEngine) [SetDTLSCertificateRequestMessageHook](https://github.com/pion/webrtc/blob/v4.2.3/settingengine.go#L626) ¶

    func (e *SettingEngine) SetDTLSCertificateRequestMessageHook(
     hook func([handshake](/github.com/pion/dtls/v3/pkg/protocol/handshake).[MessageCertificateRequest](/github.com/pion/dtls/v3/pkg/protocol/handshake#MessageCertificateRequest)) [handshake](/github.com/pion/dtls/v3/pkg/protocol/handshake).[Message](/github.com/pion/dtls/v3/pkg/protocol/handshake#Message),
    )

SetDTLSCertificateRequestMessageHook if not nil, is called when a DTLS Certificate Request message is sent from a client. The returned handshake message replaces the original message.

#### func (*SettingEngine) [SetDTLSCipherSuites](https://github.com/pion/webrtc/blob/v4.2.3/settingengine.go#L602) ¶ added in v4.1.7

    func (e *SettingEngine) SetDTLSCipherSuites(cipherSuites ...[dtls](/github.com/pion/dtls/v3).[CipherSuiteID](/github.com/pion/dtls/v3#CipherSuiteID))

SetDTLSCipherSuites allows the user to specify a list of DTLS CipherSuites. This allow to control which ciphers implemented by pion/dtls are used during the DTLS handshake. It can be used for DTLS connection hardening.

#### func (*SettingEngine) [SetDTLSClientAuth](https://github.com/pion/webrtc/blob/v4.2.3/settingengine.go#L560) ¶

    func (e *SettingEngine) SetDTLSClientAuth(clientAuth [dtls](/github.com/pion/dtls/v3).[ClientAuthType](/github.com/pion/dtls/v3#ClientAuthType))

SetDTLSClientAuth sets the client auth type for DTLS.

#### func (*SettingEngine) [SetDTLSClientCAs](https://github.com/pion/webrtc/blob/v4.2.3/settingengine.go#L565) ¶

    func (e *SettingEngine) SetDTLSClientCAs(clientCAs *[x509](/crypto/x509).[CertPool](/crypto/x509#CertPool))

SetDTLSClientCAs sets the client CA certificate pool for DTLS certificate verification.

#### func (*SettingEngine) [SetDTLSClientHelloMessageHook](https://github.com/pion/webrtc/blob/v4.2.3/settingengine.go#L614) ¶

    func (e *SettingEngine) SetDTLSClientHelloMessageHook(hook func([handshake](/github.com/pion/dtls/v3/pkg/protocol/handshake).[MessageClientHello](/github.com/pion/dtls/v3/pkg/protocol/handshake#MessageClientHello)) [handshake](/github.com/pion/dtls/v3/pkg/protocol/handshake).[Message](/github.com/pion/dtls/v3/pkg/protocol/handshake#Message))

SetDTLSClientHelloMessageHook if not nil, is called when a DTLS Client Hello message is sent from a client. The returned handshake message replaces the original message.

#### func (*SettingEngine) [SetDTLSConnectContextMaker](https://github.com/pion/webrtc/blob/v4.2.3/settingengine.go#L550) ¶

    func (e *SettingEngine) SetDTLSConnectContextMaker(connectContextMaker func() ([context](/context).[Context](/context#Context), func()))

SetDTLSConnectContextMaker sets the context used during the DTLS Handshake. It can be used to extend or reduce the timeout on the DTLS Handshake. If nil, the default dtls.ConnectContextMaker is used. It can be implemented as following.

    func ConnectContextMaker() (context.Context, func()) {
     return context.WithTimeout(context.Background(), 30*time.Second)
    }
    

#### func (*SettingEngine) [SetDTLSCustomerCipherSuites](https://github.com/pion/webrtc/blob/v4.2.3/settingengine.go#L608) ¶

    func (e *SettingEngine) SetDTLSCustomerCipherSuites(customCipherSuites func() [][dtls](/github.com/pion/dtls/v3).[CipherSuite](/github.com/pion/dtls/v3#CipherSuite))

SetDTLSCustomerCipherSuites allows the user to specify a list of custom DTLS CipherSuites. It allows to use custom/private DTLS CipherSuites in addition to the ones implemented by pion/dtls.

#### func (*SettingEngine) [SetDTLSDisableInsecureSkipVerify](https://github.com/pion/webrtc/blob/v4.2.3/settingengine.go#L534) ¶

    func (e *SettingEngine) SetDTLSDisableInsecureSkipVerify(disable [bool](/builtin#bool))

SetDTLSDisableInsecureSkipVerify sets the disable skip insecure verify flag for DTLS. This controls whether a client verifies the server's certificate chain and host name.

#### func (*SettingEngine) [SetDTLSEllipticCurves](https://github.com/pion/webrtc/blob/v4.2.3/settingengine.go#L539) ¶

    func (e *SettingEngine) SetDTLSEllipticCurves(ellipticCurves ...[dtlsElliptic](/github.com/pion/dtls/v3/pkg/crypto/elliptic).[Curve](/github.com/pion/dtls/v3/pkg/crypto/elliptic#Curve))

SetDTLSEllipticCurves sets the elliptic curves for DTLS.

#### func (*SettingEngine) [SetDTLSExtendedMasterSecret](https://github.com/pion/webrtc/blob/v4.2.3/settingengine.go#L555) ¶

    func (e *SettingEngine) SetDTLSExtendedMasterSecret(extendedMasterSecret [dtls](/github.com/pion/dtls/v3).[ExtendedMasterSecretType](/github.com/pion/dtls/v3#ExtendedMasterSecretType))

SetDTLSExtendedMasterSecret sets the extended master secret type for DTLS.

#### func (*SettingEngine) [SetDTLSInsecureSkipHelloVerify](https://github.com/pion/webrtc/blob/v4.2.3/settingengine.go#L528) ¶

    func (e *SettingEngine) SetDTLSInsecureSkipHelloVerify(skip [bool](/builtin#bool))

SetDTLSInsecureSkipHelloVerify sets the skip HelloVerify flag for DTLS. If true and when acting as DTLS server, will allow client to skip hello verify phase and receive ServerHello after initial ClientHello. This will mean faster connect times, but will have lower DoS attack resistance.

#### func (*SettingEngine) [SetDTLSKeyLogWriter](https://github.com/pion/webrtc/blob/v4.2.3/settingengine.go#L576) ¶

    func (e *SettingEngine) SetDTLSKeyLogWriter(writer [io](/io).[Writer](/io#Writer))

SetDTLSKeyLogWriter sets the destination of the TLS key material for debugging. Logging key material compromises security and should only be use for debugging.

#### func (*SettingEngine) [SetDTLSReplayProtectionWindow](https://github.com/pion/webrtc/blob/v4.2.3/settingengine.go#L433) ¶

    func (e *SettingEngine) SetDTLSReplayProtectionWindow(n [uint](/builtin#uint))

SetDTLSReplayProtectionWindow sets a replay attack protection window size of DTLS connection.

#### func (*SettingEngine) [SetDTLSRetransmissionInterval](https://github.com/pion/webrtc/blob/v4.2.3/settingengine.go#L520) ¶

    func (e *SettingEngine) SetDTLSRetransmissionInterval(interval [time](/time).[Duration](/time#Duration))

SetDTLSRetransmissionInterval sets the retranmission interval for DTLS.

#### func (*SettingEngine) [SetDTLSRootCAs](https://github.com/pion/webrtc/blob/v4.2.3/settingengine.go#L570) ¶

    func (e *SettingEngine) SetDTLSRootCAs(rootCAs *[x509](/crypto/x509).[CertPool](/crypto/x509#CertPool))

SetDTLSRootCAs sets the root CA certificate pool for DTLS certificate verification.

#### func (*SettingEngine) [SetDTLSServerHelloMessageHook](https://github.com/pion/webrtc/blob/v4.2.3/settingengine.go#L620) ¶

    func (e *SettingEngine) SetDTLSServerHelloMessageHook(hook func([handshake](/github.com/pion/dtls/v3/pkg/protocol/handshake).[MessageServerHello](/github.com/pion/dtls/v3/pkg/protocol/handshake#MessageServerHello)) [handshake](/github.com/pion/dtls/v3/pkg/protocol/handshake).[Message](/github.com/pion/dtls/v3/pkg/protocol/handshake#Message))

SetDTLSServerHelloMessageHook if not nil, is called when a DTLS Server Hello message is sent from a client. The returned handshake message replaces the original message.

#### func (*SettingEngine) [SetEphemeralUDPPortRange](https://github.com/pion/webrtc/blob/v4.2.3/settingengine.go#L270) ¶

    func (e *SettingEngine) SetEphemeralUDPPortRange(portMin, portMax [uint16](/builtin#uint16)) [error](/builtin#error)

SetEphemeralUDPPortRange limits the pool of ephemeral ports that ICE UDP connections can allocate from. This affects both host candidates, and the local address of server reflexive candidates.

When portMin and portMax are left to the 0 default value, pion/ice candidate gatherer replaces them and uses 1 for portMin and 65535 for portMax.

#### func (*SettingEngine) [SetFireOnTrackBeforeFirstRTP](https://github.com/pion/webrtc/blob/v4.2.3/settingengine.go#L670) ¶ added in v4.0.1

    func (e *SettingEngine) SetFireOnTrackBeforeFirstRTP(fireOnTrackBeforeFirstRTP [bool](/builtin#bool))

SetFireOnTrackBeforeFirstRTP sets if firing the OnTrack event should happen before any RTP packets are received. Setting this to true will have the Track's Codec and PayloadTypes be initially set to their zero values in the OnTrack handler. Note: This does not yet affect simulcast tracks.

#### func (*SettingEngine) [SetHandleUndeclaredSSRCWithoutAnswer](https://github.com/pion/webrtc/blob/v4.2.3/settingengine.go#L683) ¶ added in v4.1.4

    func (e *SettingEngine) SetHandleUndeclaredSSRCWithoutAnswer(handleUndeclaredSSRCWithoutAnswer [bool](/builtin#bool))

SetHandleUndeclaredSSRCWithoutAnswer controls if an SDP answer is required for processing early media of non-simulcast tracks.

#### func (*SettingEngine) [SetHostAcceptanceMinWait](https://github.com/pion/webrtc/blob/v4.2.3/settingengine.go#L240) ¶

    func (e *SettingEngine) SetHostAcceptanceMinWait(t [time](/time).[Duration](/time#Duration))

SetHostAcceptanceMinWait sets the ICEHostAcceptanceMinWait.

#### func (*SettingEngine) [SetICEAddressRewriteRules](https://github.com/pion/webrtc/blob/v4.2.3/settingengine.go#L349) ¶ added in v4.2.0

    func (e *SettingEngine) SetICEAddressRewriteRules(rules ...ICEAddressRewriteRule) [error](/builtin#error)

SetICEAddressRewriteRules configures address rewrite rules for candidate publication. These rules provide fine-grained control over which local addresses are replaced or supplemented with external IPs. This replaces the legacy NAT1To1 settings, which will be deprecated in the future.

Example (AppendSrflx) ¶

ExampleSettingEngine_SetICEAddressRewriteRules_appendSrflx demonstrates appending a server reflexive candidate that advertises a public address while still keeping the host candidate.

    var se SettingEngine
    
    _ = se.SetICEAddressRewriteRules(
     ICEAddressRewriteRule{
      External:        []string{"198.51.100.2"},
      AsCandidateType: ICECandidateTypeSrflx,
      Mode:            ICEAddressRewriteAppend,
     },
    )
    

Example (ReplaceHost) ¶

ExampleSettingEngine_SetICEAddressRewriteRules_replaceHost demonstrates replacing host candidates with a fixed public address using the rewrite API.

    var se SettingEngine
    
    _ = se.SetICEAddressRewriteRules(
     ICEAddressRewriteRule{
      External:        []string{"198.51.100.1"},
      AsCandidateType: ICECandidateTypeHost,
      Mode:            ICEAddressRewriteReplace,
     },
    )
    

#### func (*SettingEngine) [SetICEBindingRequestHandler](https://github.com/pion/webrtc/blob/v4.2.3/settingengine.go#L659) ¶

    func (e *SettingEngine) SetICEBindingRequestHandler(
     bindingRequestHandler func(m *[stun](/github.com/pion/stun/v3).[Message](/github.com/pion/stun/v3#Message), local, remote [ice](/github.com/pion/ice/v4).[Candidate](/github.com/pion/ice/v4#Candidate), pair *[ice](/github.com/pion/ice/v4).[CandidatePair](/github.com/pion/ice/v4#CandidatePair)) [bool](/builtin#bool),
    )

SetICEBindingRequestHandler sets a callback that is fired on a STUN BindingRequest This allows users to do things like \- Log incoming Binding Requests for debugging \- Implement draft-thatcher-ice-renomination \- Implement custom CandidatePair switching logic.

#### func (*SettingEngine) [SetICECredentials](https://github.com/pion/webrtc/blob/v4.2.3/settingengine.go#L422) ¶

    func (e *SettingEngine) SetICECredentials(usernameFragment, password [string](/builtin#string))

SetICECredentials sets a staic uFrag/uPwd to be used by pion/ice

This is useful if you want to do signalless WebRTC session, or having a reproducible environment with static credentials.

#### func (*SettingEngine) [SetICEMaxBindingRequests](https://github.com/pion/webrtc/blob/v4.2.3/settingengine.go#L487) ¶

    func (e *SettingEngine) SetICEMaxBindingRequests(d [uint16](/builtin#uint16))

SetICEMaxBindingRequests sets the maximum amount of binding requests that can be sent on a candidate before it is considered invalid.

#### func (*SettingEngine) [SetICEMulticastDNSMode](https://github.com/pion/webrtc/blob/v4.2.3/settingengine.go#L406) ¶

    func (e *SettingEngine) SetICEMulticastDNSMode(multicastDNSMode [ice](/github.com/pion/ice/v4).[MulticastDNSMode](/github.com/pion/ice/v4#MulticastDNSMode))

SetICEMulticastDNSMode controls if pion/ice queries and generates mDNS ICE Candidates.

#### func (*SettingEngine) [SetICEProxyDialer](https://github.com/pion/webrtc/blob/v4.2.3/settingengine.go#L481) ¶

    func (e *SettingEngine) SetICEProxyDialer(d [proxy](/golang.org/x/net/proxy).[Dialer](/golang.org/x/net/proxy#Dialer))

SetICEProxyDialer sets the proxy dialer interface based on golang.org/x/net/proxy.

#### func (*SettingEngine) [SetICERenomination](https://github.com/pion/webrtc/blob/v4.2.3/settingengine.go#L158) ¶ added in v4.2.0

    func (e *SettingEngine) SetICERenomination(options ...RenominationOption) [error](/builtin#error)

SetICERenomination configures ICE renomination using options for generator and scheduling. Manual control is not exposed yet. This always enables automatic renomination with the default generator unless a custom one is provided.

#### func (*SettingEngine) [SetICETCPMux](https://github.com/pion/webrtc/blob/v4.2.3/settingengine.go#L469) ¶

    func (e *SettingEngine) SetICETCPMux(tcpMux [ice](/github.com/pion/ice/v4).[TCPMux](/github.com/pion/ice/v4#TCPMux))

SetICETCPMux enables ICE-TCP when set to a non-nil value. Make sure that NetworkTypeTCP4 or NetworkTypeTCP6 is enabled as well.

#### func (*SettingEngine) [SetICETimeouts](https://github.com/pion/webrtc/blob/v4.2.3/settingengine.go#L233) ¶

    func (e *SettingEngine) SetICETimeouts(disconnectedTimeout, failedTimeout, keepAliveInterval [time](/time).[Duration](/time#Duration))

SetICETimeouts sets the behavior around ICE Timeouts

disconnectedTimeout:

    Duration without network activity before an Agent is considered disconnected. Default is 5 Seconds
    

failedTimeout:

    Duration without network activity before an Agent is considered failed after disconnected. Default is 25 Seconds
    

keepAliveInterval:

    How often the ICE Agent sends extra traffic if there is no activity, if media is flowing no traffic will be sent.
    

Default is 2 seconds.

#### func (*SettingEngine) [SetICEUDPMux](https://github.com/pion/webrtc/blob/v4.2.3/settingengine.go#L476) ¶

    func (e *SettingEngine) SetICEUDPMux(udpMux [ice](/github.com/pion/ice/v4).[UDPMux](/github.com/pion/ice/v4#UDPMux))

SetICEUDPMux allows ICE traffic to come through a single UDP port, drastically simplifying deployments where ports will need to be opened/forwarded. UDPMux should be started prior to creating PeerConnections.

#### func (*SettingEngine) [SetIPFilter](https://github.com/pion/webrtc/blob/v4.2.3/settingengine.go#L304) ¶

    func (e *SettingEngine) SetIPFilter(filter func([net](/net).[IP](/net#IP)) (keep [bool](/builtin#bool)))

SetIPFilter sets the filtering functions when gathering ICE candidates This can be used to exclude certain ip from ICE. Which may be useful if you know a certain ip will never succeed, or if you wish to reduce the amount of information you wish to expose to the remote peer.

#### func (*SettingEngine) [SetIgnoreRidPauseForRecv](https://github.com/pion/webrtc/blob/v4.2.3/settingengine.go#L689) ¶ added in v4.1.7

    func (e *SettingEngine) SetIgnoreRidPauseForRecv(ignoreRidPauseForRecv [bool](/builtin#bool))

SetIgnoreRidPauseForRecv controls if SDP `a=simulcast:recv` will include the paused attribute of a RID (simulcast layer).

#### func (*SettingEngine) [SetIncludeLoopbackCandidate](https://github.com/pion/webrtc/blob/v4.2.3/settingengine.go#L372) ¶

    func (e *SettingEngine) SetIncludeLoopbackCandidate(include [bool](/builtin#bool))

SetIncludeLoopbackCandidate enable pion to gather loopback candidates, it is useful for some VM have public IP mapped to loopback interface.

#### func (*SettingEngine) [SetInterfaceFilter](https://github.com/pion/webrtc/blob/v4.2.3/settingengine.go#L296) ¶

    func (e *SettingEngine) SetInterfaceFilter(filter func([string](/builtin#string)) (keep [bool](/builtin#bool)))

SetInterfaceFilter sets the filtering functions when gathering ICE candidates This can be used to exclude certain network interfaces from ICE. Which may be useful if you know a certain interface will never succeed, or if you wish to reduce the amount of information you wish to expose to the remote peer.

#### func (*SettingEngine) [SetLite](https://github.com/pion/webrtc/blob/v4.2.3/settingengine.go#L282) ¶

    func (e *SettingEngine) SetLite(lite [bool](/builtin#bool))

SetLite configures whether or not the ice agent should be a lite agent.

#### func (*SettingEngine) [SetMulticastDNSHostName](https://github.com/pion/webrtc/blob/v4.2.3/settingengine.go#L414) ¶

    func (e *SettingEngine) SetMulticastDNSHostName(hostName [string](/builtin#string))

SetMulticastDNSHostName sets a static HostName to be used by pion/ice instead of generating one on startup

This should only be used for a single PeerConnection. Having multiple PeerConnections with the same HostName will cause undefined behavior.

#### func (*SettingEngine) [SetNAT1To1IPs](https://github.com/pion/webrtc/blob/v4.2.3/settingengine.go#L340) deprecated

    func (e *SettingEngine) SetNAT1To1IPs(ips [][string](/builtin#string), candidateType ICECandidateType)

SetNAT1To1IPs sets a list of external IP addresses of 1:1 (D)NAT and a candidate type for which the external IP address is used. This is useful when you host a server using Pion on an AWS EC2 instance which has a private address, behind a 1:1 DNAT with a public IP (e.g. Elastic IP). In this case, you can give the public IP address so that Pion will use the public IP address in its candidate instead of the private IP address. The second argument, candidateType, is used to tell Pion which type of candidate should use the given public IP address. Two types of candidates are supported:

ICECandidateTypeHost:

    The public IP address will be used for the host candidate in the SDP.
    

ICECandidateTypeSrflx:

    A server reflexive candidate with the given public IP address will be added to the SDP.
    

Please note that if you choose ICECandidateTypeHost, then the private IP address won't be advertised with the peer. Also, this option cannot be used along with mDNS.

If you choose ICECandidateTypeSrflx, it simply adds a server reflexive candidate with the public IP. The host candidate is still available along with mDNS capabilities unaffected. Also, you cannot give STUN server URL at the same time. It will result in an error otherwise.

Deprecated: Use SetICEAddressRewriteRules instead. To mirror the legacy behavior, supply ICEAddressRewriteRule with External set to ips, AsCandidateType set to candidateType, and Mode set to ICEAddressRewriteReplace for host candidates or ICEAddressRewriteAppend for server reflexive candidates. Or leave Mode unspecified to use the default behavior; replace for host candidates and append for server reflexive candidates.

#### func (*SettingEngine) [SetNet](https://github.com/pion/webrtc/blob/v4.2.3/settingengine.go#L401) ¶

    func (e *SettingEngine) SetNet(net [transport](/github.com/pion/transport/v4).[Net](/github.com/pion/transport/v4#Net))

SetNet sets the Net instance that is passed to pion/ice

Net is an network interface layer for Pion, allowing users to replace Pions network stack with a custom implementation.

#### func (*SettingEngine) [SetNetworkTypes](https://github.com/pion/webrtc/blob/v4.2.3/settingengine.go#L288) ¶

    func (e *SettingEngine) SetNetworkTypes(candidateTypes []NetworkType)

SetNetworkTypes configures what types of candidate networks are supported during local and server reflexive gathering.

#### func (*SettingEngine) [SetPrflxAcceptanceMinWait](https://github.com/pion/webrtc/blob/v4.2.3/settingengine.go#L250) ¶

    func (e *SettingEngine) SetPrflxAcceptanceMinWait(t [time](/time).[Duration](/time#Duration))

SetPrflxAcceptanceMinWait sets the ICEPrflxAcceptanceMinWait.

#### func (*SettingEngine) [SetReceiveMTU](https://github.com/pion/webrtc/blob/v4.2.3/settingengine.go#L515) ¶

    func (e *SettingEngine) SetReceiveMTU(receiveMTU [uint](/builtin#uint))

SetReceiveMTU sets the size of read buffer that copies incoming packets. This is optional. Leave this 0 for the default receiveMTU.

#### func (*SettingEngine) [SetRelayAcceptanceMinWait](https://github.com/pion/webrtc/blob/v4.2.3/settingengine.go#L255) ¶

    func (e *SettingEngine) SetRelayAcceptanceMinWait(t [time](/time).[Duration](/time#Duration))

SetRelayAcceptanceMinWait sets the ICERelayAcceptanceMinWait.

#### func (*SettingEngine) [SetSCTPCwndCAStep](https://github.com/pion/webrtc/blob/v4.2.3/settingengine.go#L650) ¶ added in v4.1.1

    func (e *SettingEngine) SetSCTPCwndCAStep(cwndCAStep [uint32](/builtin#uint32))

SetSCTPCwndCAStep sets congestion window adjustment step size during congestion avoidance.

#### func (*SettingEngine) [SetSCTPFastRtxWnd](https://github.com/pion/webrtc/blob/v4.2.3/settingengine.go#L645) ¶ added in v4.1.1

    func (e *SettingEngine) SetSCTPFastRtxWnd(fastRtxWnd [uint32](/builtin#uint32))

SetSCTPFastRtxWnd sets the fast retransmission window size.

#### func (*SettingEngine) [SetSCTPMaxMessageSize](https://github.com/pion/webrtc/blob/v4.2.3/settingengine.go#L595) ¶ added in v4.0.13

    func (e *SettingEngine) SetSCTPMaxMessageSize(maxMessageSize [uint32](/builtin#uint32))

SetSCTPMaxMessageSize sets the largest message we are willing to accept. Leave this 0 for the default max message size.

#### func (*SettingEngine) [SetSCTPMaxReceiveBufferSize](https://github.com/pion/webrtc/blob/v4.2.3/settingengine.go#L582) ¶

    func (e *SettingEngine) SetSCTPMaxReceiveBufferSize(maxReceiveBufferSize [uint32](/builtin#uint32))

SetSCTPMaxReceiveBufferSize sets the maximum receive buffer size. Leave this 0 for the default maxReceiveBufferSize.

#### func (*SettingEngine) [SetSCTPMinCwnd](https://github.com/pion/webrtc/blob/v4.2.3/settingengine.go#L640) ¶ added in v4.1.1

    func (e *SettingEngine) SetSCTPMinCwnd(minCwnd [uint32](/builtin#uint32))

SetSCTPMinCwnd sets the minimum congestion window size. The congestion window will not be smaller than this value during congestion control.

#### func (*SettingEngine) [SetSCTPRTOMax](https://github.com/pion/webrtc/blob/v4.2.3/settingengine.go#L634) ¶

    func (e *SettingEngine) SetSCTPRTOMax(rtoMax [time](/time).[Duration](/time#Duration))

SetSCTPRTOMax sets the maximum retransmission timeout. Leave this 0 for the default timeout.

#### func (*SettingEngine) [SetSDPMediaLevelFingerprints](https://github.com/pion/webrtc/blob/v4.2.3/settingengine.go#L463) ¶

    func (e *SettingEngine) SetSDPMediaLevelFingerprints(sdpMediaLevelFingerprints [bool](/builtin#bool))

SetSDPMediaLevelFingerprints configures the logic for DTLS Fingerprint insertion If true, fingerprints will be inserted in the sdp at the fingerprint level, instead of the session level. This helps with compatibility with some webrtc implementations.

#### func (*SettingEngine) [SetSRTCPReplayProtectionWindow](https://github.com/pion/webrtc/blob/v4.2.3/settingengine.go#L444) ¶

    func (e *SettingEngine) SetSRTCPReplayProtectionWindow(n [uint](/builtin#uint))

SetSRTCPReplayProtectionWindow sets a replay attack protection window size of SRTCP session.

#### func (*SettingEngine) [SetSRTPProtectionProfiles](https://github.com/pion/webrtc/blob/v4.2.3/settingengine.go#L214) ¶

    func (e *SettingEngine) SetSRTPProtectionProfiles(profiles ...[dtls](/github.com/pion/dtls/v3).[SRTPProtectionProfile](/github.com/pion/dtls/v3#SRTPProtectionProfile))

SetSRTPProtectionProfiles allows the user to override the default SRTP Protection Profiles The default srtp protection profiles are provided by the function `defaultSrtpProtectionProfiles`.

#### func (*SettingEngine) [SetSRTPReplayProtectionWindow](https://github.com/pion/webrtc/blob/v4.2.3/settingengine.go#L438) ¶

    func (e *SettingEngine) SetSRTPReplayProtectionWindow(n [uint](/builtin#uint))

SetSRTPReplayProtectionWindow sets a replay attack protection window size of SRTP session.

#### func (*SettingEngine) [SetSTUNGatherTimeout](https://github.com/pion/webrtc/blob/v4.2.3/settingengine.go#L260) ¶

    func (e *SettingEngine) SetSTUNGatherTimeout(t [time](/time).[Duration](/time#Duration))

SetSTUNGatherTimeout sets the ICESTUNGatherTimeout.

#### func (*SettingEngine) [SetSrflxAcceptanceMinWait](https://github.com/pion/webrtc/blob/v4.2.3/settingengine.go#L245) ¶

    func (e *SettingEngine) SetSrflxAcceptanceMinWait(t [time](/time).[Duration](/time#Duration))

SetSrflxAcceptanceMinWait sets the ICESrflxAcceptanceMinWait.

#### type [SignalingState](https://github.com/pion/webrtc/blob/v4.2.3/signalingstate.go#L32) ¶

    type SignalingState [int32](/builtin#int32)

SignalingState indicates the signaling state of the offer/answer process.

    const (
     // SignalingStateUnknown is the enum's zero-value.
     SignalingStateUnknown SignalingState = [iota](/builtin#iota)
    
     // SignalingStateStable indicates there is no offer/answer exchange in
     // progress. This is also the initial state, in which case the local and
     // remote descriptions are nil.
     SignalingStateStable
    
     // SignalingStateHaveLocalOffer indicates that a local description, of
     // type "offer", has been successfully applied.
     SignalingStateHaveLocalOffer
    
     // SignalingStateHaveRemoteOffer indicates that a remote description, of
     // type "offer", has been successfully applied.
     SignalingStateHaveRemoteOffer
    
     // SignalingStateHaveLocalPranswer indicates that a remote description
     // of type "offer" has been successfully applied and a local description
     // of type "pranswer" has been successfully applied.
     SignalingStateHaveLocalPranswer
    
     // SignalingStateHaveRemotePranswer indicates that a local description
     // of type "offer" has been successfully applied and a remote description
     // of type "pranswer" has been successfully applied.
     SignalingStateHaveRemotePranswer
    
     // SignalingStateClosed indicates The PeerConnection has been closed.
     SignalingStateClosed
    )

#### func (*SignalingState) [Get](https://github.com/pion/webrtc/blob/v4.2.3/signalingstate.go#L114) ¶

    func (t *SignalingState) Get() SignalingState

Get thread safe read value.

#### func (*SignalingState) [Set](https://github.com/pion/webrtc/blob/v4.2.3/signalingstate.go#L119) ¶

    func (t *SignalingState) Set(state SignalingState)

Set thread safe write value.

#### func (SignalingState) [String](https://github.com/pion/webrtc/blob/v4.2.3/signalingstate.go#L94) ¶

    func (t SignalingState) String() [string](/builtin#string)

#### type [Stats](https://github.com/pion/webrtc/blob/v4.2.3/stats.go#L17) ¶

    type Stats interface {
     // contains filtered or unexported methods
    }

A Stats object contains a set of statistics copies out of a monitored component of the WebRTC stack at a specific time.

#### func [UnmarshalStatsJSON](https://github.com/pion/webrtc/blob/v4.2.3/stats.go#L22) ¶

    func UnmarshalStatsJSON(b [][byte](/builtin#byte)) (Stats, [error](/builtin#error))

UnmarshalStatsJSON unmarshals a Stats object from JSON.

#### type [StatsICECandidatePairState](https://github.com/pion/webrtc/blob/v4.2.3/stats.go#L2013) ¶

    type StatsICECandidatePairState [string](/builtin#string)

StatsICECandidatePairState is the state of an ICE candidate pair used in the ICECandidatePairStats object.

    const (
     // StatsICECandidatePairStateFrozen means a check for this pair hasn't been
     // performed, and it can't yet be performed until some other check succeeds,
     // allowing this pair to unfreeze and move into the Waiting state.
     StatsICECandidatePairStateFrozen StatsICECandidatePairState = "frozen"
    
     // StatsICECandidatePairStateWaiting means a check has not been performed for
     // this pair, and can be performed as soon as it is the highest-priority Waiting
     // pair on the check list.
     StatsICECandidatePairStateWaiting StatsICECandidatePairState = "waiting"
    
     // StatsICECandidatePairStateInProgress means a check has been sent for this pair,
     // but the transaction is in progress.
     StatsICECandidatePairStateInProgress StatsICECandidatePairState = "in-progress"
    
     // StatsICECandidatePairStateFailed means a check for this pair was already done
     // and failed, either never producing any response or producing an unrecoverable
     // failure response.
     StatsICECandidatePairStateFailed StatsICECandidatePairState = "failed"
    
     // StatsICECandidatePairStateSucceeded means a check for this pair was already
     // done and produced a successful result.
     StatsICECandidatePairStateSucceeded StatsICECandidatePairState = "succeeded"
    )

#### type [StatsReport](https://github.com/pion/webrtc/blob/v4.2.3/stats.go#L173) ¶

    type StatsReport map[[string](/builtin#string)]Stats

StatsReport collects Stats objects indexed by their ID.

#### func (StatsReport) [GetCertificateStats](https://github.com/pion/webrtc/blob/v4.2.3/stats_go.go#L80) ¶

    func (r StatsReport) GetCertificateStats(c *Certificate) (CertificateStats, [bool](/builtin#bool))

GetCertificateStats is a helper method to return the associated stats for a given Certificate.

#### func (StatsReport) [GetCodecStats](https://github.com/pion/webrtc/blob/v4.2.3/stats_go.go#L96) ¶

    func (r StatsReport) GetCodecStats(c *RTPCodecParameters) (CodecStats, [bool](/builtin#bool))

GetCodecStats is a helper method to return the associated stats for a given Codec.

#### func (StatsReport) [GetConnectionStats](https://github.com/pion/webrtc/blob/v4.2.3/stats_go.go#L16) ¶

    func (r StatsReport) GetConnectionStats(conn *PeerConnection) (PeerConnectionStats, [bool](/builtin#bool))

GetConnectionStats is a helper method to return the associated stats for a given PeerConnection.

#### func (StatsReport) [GetDataChannelStats](https://github.com/pion/webrtc/blob/v4.2.3/stats_go.go#L32) ¶

    func (r StatsReport) GetDataChannelStats(dc *DataChannel) (DataChannelStats, [bool](/builtin#bool))

GetDataChannelStats is a helper method to return the associated stats for a given DataChannel.

#### func (StatsReport) [GetICECandidatePairStats](https://github.com/pion/webrtc/blob/v4.2.3/stats_go.go#L64) ¶

    func (r StatsReport) GetICECandidatePairStats(c *ICECandidatePair) (ICECandidatePairStats, [bool](/builtin#bool))

GetICECandidatePairStats is a helper method to return the associated stats for a given ICECandidatePair.

#### func (StatsReport) [GetICECandidateStats](https://github.com/pion/webrtc/blob/v4.2.3/stats_go.go#L48) ¶

    func (r StatsReport) GetICECandidateStats(c *ICECandidate) (ICECandidateStats, [bool](/builtin#bool))

GetICECandidateStats is a helper method to return the associated stats for a given ICECandidate.

#### type [StatsTimestamp](https://github.com/pion/webrtc/blob/v4.2.3/stats.go#L154) ¶

    type StatsTimestamp [float64](/builtin#float64)

StatsTimestamp is a timestamp represented by the floating point number of milliseconds since the epoch.

#### func (StatsTimestamp) [Time](https://github.com/pion/webrtc/blob/v4.2.3/stats.go#L157) ¶

    func (s StatsTimestamp) Time() [time](/time).[Time](/time#Time)

Time returns the time.Time represented by this timestamp.

#### type [StatsType](https://github.com/pion/webrtc/blob/v4.2.3/stats.go#L78) ¶

    type StatsType [string](/builtin#string)

StatsType indicates the type of the object that a Stats object represents.

    const (
     // StatsTypeCodec is used by CodecStats.
     StatsTypeCodec StatsType = "codec"
    
     // StatsTypeInboundRTP is used by InboundRTPStreamStats.
     StatsTypeInboundRTP StatsType = "inbound-rtp"
    
     // StatsTypeOutboundRTP is used by OutboundRTPStreamStats.
     StatsTypeOutboundRTP StatsType = "outbound-rtp"
    
     // StatsTypeRemoteInboundRTP is used by RemoteInboundRTPStreamStats.
     StatsTypeRemoteInboundRTP StatsType = "remote-inbound-rtp"
    
     // StatsTypeRemoteOutboundRTP is used by RemoteOutboundRTPStreamStats.
     StatsTypeRemoteOutboundRTP StatsType = "remote-outbound-rtp"
    
     // StatsTypeCSRC is used by RTPContributingSourceStats.
     StatsTypeCSRC StatsType = "csrc"
    
     // StatsTypeMediaSource is used by AudioSourceStats or VideoSourceStats depending on kind.
     StatsTypeMediaSource = "media-source"
    
     // StatsTypeMediaPlayout is used by AudioPlayoutStats.
     StatsTypeMediaPlayout StatsType = "media-playout"
    
     // StatsTypePeerConnection used by PeerConnectionStats.
     StatsTypePeerConnection StatsType = "peer-connection"
    
     // StatsTypeDataChannel is used by DataChannelStats.
     StatsTypeDataChannel StatsType = "data-channel"
    
     // StatsTypeStream is used by MediaStreamStats.
     StatsTypeStream StatsType = "stream"
    
     // StatsTypeTrack is used by SenderVideoTrackAttachmentStats and SenderAudioTrackAttachmentStats depending on kind.
     StatsTypeTrack StatsType = "track"
    
     // StatsTypeSender is used by the AudioSenderStats or VideoSenderStats depending on kind.
     StatsTypeSender StatsType = "sender"
    
     // StatsTypeReceiver is used by the AudioReceiverStats or VideoReceiverStats depending on kind.
     StatsTypeReceiver StatsType = "receiver"
    
     // StatsTypeTransport is used by TransportStats.
     StatsTypeTransport StatsType = "transport"
    
     // StatsTypeCandidatePair is used by ICECandidatePairStats.
     StatsTypeCandidatePair StatsType = "candidate-pair"
    
     // StatsTypeLocalCandidate is used by ICECandidateStats for the local candidate.
     StatsTypeLocalCandidate StatsType = "local-candidate"
    
     // StatsTypeRemoteCandidate is used by ICECandidateStats for the remote candidate.
     StatsTypeRemoteCandidate StatsType = "remote-candidate"
    
     // StatsTypeCertificate is used by CertificateStats.
     StatsTypeCertificate StatsType = "certificate"
    
     // StatsTypeSCTPTransport is used by SCTPTransportStats.
     StatsTypeSCTPTransport StatsType = "sctp-transport"
    )

#### type [TrackLocal](https://github.com/pion/webrtc/blob/v4.2.3/track_local.go#L105) ¶

    type TrackLocal interface {
     // Bind should implement the way how the media data flows from the Track to the PeerConnection
     // This will be called internally after signaling is complete and the list of available
     // codecs has been determined
     Bind(TrackLocalContext) (RTPCodecParameters, [error](/builtin#error))
    
     // Unbind should implement the teardown logic when the track is no longer needed. This happens
     // because a track has been stopped.
     Unbind(TrackLocalContext) [error](/builtin#error)
    
     // ID is the unique identifier for this Track. This should be unique for the
     // stream, but doesn't have to globally unique. A common example would be 'audio' or 'video'
     // and StreamID would be 'desktop' or 'webcam'
     ID() [string](/builtin#string)
    
     // RID is the RTP Stream ID for this track.
     RID() [string](/builtin#string)
    
     // StreamID is the group this track belongs too. This must be unique
     StreamID() [string](/builtin#string)
    
     // Kind controls if this TrackLocal is audio or video
     Kind() RTPCodecType
    }

TrackLocal is an interface that controls how the user can send media The user can provide their own TrackLocal implementations, or use the implementations in pkg/media.

#### type [TrackLocalContext](https://github.com/pion/webrtc/blob/v4.2.3/track_local.go#L22) ¶

    type TrackLocalContext interface {
     // CodecParameters returns the negotiated RTPCodecParameters. These are the codecs supported by both
     // PeerConnections and the PayloadTypes
     CodecParameters() []RTPCodecParameters
    
     // HeaderExtensions returns the negotiated RTPHeaderExtensionParameters. These are the header extensions supported by
     // both PeerConnections and the URI/IDs
     HeaderExtensions() []RTPHeaderExtensionParameter
    
     // SSRC returns the negotiated SSRC of this track
     SSRC() SSRC
    
     // SSRCRetransmission returns the negotiated SSRC used to send retransmissions for this track
     SSRCRetransmission() SSRC
    
     // SSRCForwardErrorCorrection returns the negotiated SSRC to send forward error correction for this track
     SSRCForwardErrorCorrection() SSRC
    
     // WriteStream returns the WriteStream for this TrackLocal. The implementer writes the outbound
     // media packets to it
     WriteStream() TrackLocalWriter
    
     // ID is a unique identifier that is used for both Bind/Unbind
     ID() [string](/builtin#string)
    
     // RTCPReader returns the RTCP interceptor for this TrackLocal. Used to read RTCP of this TrackLocal.
     RTCPReader() [interceptor](/github.com/pion/interceptor).[RTCPReader](/github.com/pion/interceptor#RTCPReader)
    }

TrackLocalContext is the Context passed when a TrackLocal has been Binded/Unbinded from a PeerConnection, and used in Interceptors.

#### type [TrackLocalStaticRTP](https://github.com/pion/webrtc/blob/v4.2.3/track_local_static.go#L30) ¶

    type TrackLocalStaticRTP struct {
     // contains filtered or unexported fields
    }

TrackLocalStaticRTP is a TrackLocal that has a pre-set codec and accepts RTP Packets. If you wish to send a media.Sample use TrackLocalStaticSample.

#### func [NewTrackLocalStaticRTP](https://github.com/pion/webrtc/blob/v4.2.3/track_local_static.go#L41) ¶

    func NewTrackLocalStaticRTP(
     c RTPCodecCapability,
     id, streamID [string](/builtin#string),
     options ...func(*TrackLocalStaticRTP),
    ) (*TrackLocalStaticRTP, [error](/builtin#error))

NewTrackLocalStaticRTP returns a TrackLocalStaticRTP.

#### func (*TrackLocalStaticRTP) [Bind](https://github.com/pion/webrtc/blob/v4.2.3/track_local_static.go#L91) ¶

    func (s *TrackLocalStaticRTP) Bind(trackContext TrackLocalContext) (RTPCodecParameters, [error](/builtin#error))

Bind is called by the PeerConnection after negotiation is complete This asserts that the code requested is supported by the remote peer. If so it sets up all the state (SSRC and PayloadType) to have a call.

#### func (*TrackLocalStaticRTP) [Codec](https://github.com/pion/webrtc/blob/v4.2.3/track_local_static.go#L158) ¶

    func (s *TrackLocalStaticRTP) Codec() RTPCodecCapability

Codec gets the Codec of the track.

#### func (*TrackLocalStaticRTP) [ID](https://github.com/pion/webrtc/blob/v4.2.3/track_local_static.go#L137) ¶

    func (s *TrackLocalStaticRTP) ID() [string](/builtin#string)

ID is the unique identifier for this Track. This should be unique for the stream, but doesn't have to globally unique. A common example would be 'audio' or 'video' and StreamID would be 'desktop' or 'webcam'.

#### func (*TrackLocalStaticRTP) [Kind](https://github.com/pion/webrtc/blob/v4.2.3/track_local_static.go#L146) ¶

    func (s *TrackLocalStaticRTP) Kind() RTPCodecType

Kind controls if this TrackLocal is audio or video.

#### func (*TrackLocalStaticRTP) [RID](https://github.com/pion/webrtc/blob/v4.2.3/track_local_static.go#L143) ¶

    func (s *TrackLocalStaticRTP) RID() [string](/builtin#string)

RID is the RTP stream identifier.

#### func (*TrackLocalStaticRTP) [StreamID](https://github.com/pion/webrtc/blob/v4.2.3/track_local_static.go#L140) ¶

    func (s *TrackLocalStaticRTP) StreamID() [string](/builtin#string)

StreamID is the group this track belongs too. This must be unique.

#### func (*TrackLocalStaticRTP) [Unbind](https://github.com/pion/webrtc/blob/v4.2.3/track_local_static.go#L118) ¶

    func (s *TrackLocalStaticRTP) Unbind(t TrackLocalContext) [error](/builtin#error)

Unbind implements the teardown logic when the track is no longer needed. This happens because a track has been stopped.

#### func (*TrackLocalStaticRTP) [Write](https://github.com/pion/webrtc/blob/v4.2.3/track_local_static.go#L222) ¶

    func (s *TrackLocalStaticRTP) Write(b [][byte](/builtin#byte)) (n [int](/builtin#int), err [error](/builtin#error))

Write writes a RTP Packet as a buffer to the TrackLocalStaticRTP If one PeerConnection fails the packets will still be sent to all PeerConnections. The error message will contain the ID of the failed PeerConnections so you can remove them.

#### func (*TrackLocalStaticRTP) [WriteRTP](https://github.com/pion/webrtc/blob/v4.2.3/track_local_static.go#L185) ¶

    func (s *TrackLocalStaticRTP) WriteRTP(p *[rtp](/github.com/pion/rtp).[Packet](/github.com/pion/rtp#Packet)) [error](/builtin#error)

WriteRTP writes a RTP Packet to the TrackLocalStaticRTP If one PeerConnection fails the packets will still be sent to all PeerConnections. The error message will contain the ID of the failed PeerConnections so you can remove them.

#### type [TrackLocalStaticSample](https://github.com/pion/webrtc/blob/v4.2.3/track_local_static.go#L236) ¶

    type TrackLocalStaticSample struct {
     // contains filtered or unexported fields
    }

TrackLocalStaticSample is a TrackLocal that has a pre-set codec and accepts Samples. If you wish to send a RTP Packet use TrackLocalStaticRTP.

#### func [NewTrackLocalStaticSample](https://github.com/pion/webrtc/blob/v4.2.3/track_local_static.go#L246) ¶

    func NewTrackLocalStaticSample(
     c RTPCodecCapability,
     id, streamID [string](/builtin#string),
     options ...func(*TrackLocalStaticRTP),
    ) (*TrackLocalStaticSample, [error](/builtin#error))

NewTrackLocalStaticSample returns a TrackLocalStaticSample.

#### func (*TrackLocalStaticSample) [Bind](https://github.com/pion/webrtc/blob/v4.2.3/track_local_static.go#L283) ¶

    func (s *TrackLocalStaticSample) Bind(t TrackLocalContext) (RTPCodecParameters, [error](/builtin#error))

Bind is called by the PeerConnection after negotiation is complete This asserts that the code requested is supported by the remote peer. If so it setups all the state (SSRC and PayloadType) to have a call.

#### func (*TrackLocalStaticSample) [Codec](https://github.com/pion/webrtc/blob/v4.2.3/track_local_static.go#L276) ¶

    func (s *TrackLocalStaticSample) Codec() RTPCodecCapability

Codec gets the Codec of the track.

#### func (*TrackLocalStaticSample) [GeneratePadding](https://github.com/pion/webrtc/blob/v4.2.3/track_local_static.go#L393) ¶

    func (s *TrackLocalStaticSample) GeneratePadding(samples [uint32](/builtin#uint32)) [error](/builtin#error)

GeneratePadding writes padding-only samples to the TrackLocalStaticSample If one PeerConnection fails the packets will still be sent to all PeerConnections. The error message will contain the ID of the failed PeerConnections so you can remove them.

#### func (*TrackLocalStaticSample) [ID](https://github.com/pion/webrtc/blob/v4.2.3/track_local_static.go#L264) ¶

    func (s *TrackLocalStaticSample) ID() [string](/builtin#string)

ID is the unique identifier for this Track. This should be unique for the stream, but doesn't have to globally unique. A common example would be 'audio' or 'video' and StreamID would be 'desktop' or 'webcam'.

#### func (*TrackLocalStaticSample) [Kind](https://github.com/pion/webrtc/blob/v4.2.3/track_local_static.go#L273) ¶

    func (s *TrackLocalStaticSample) Kind() RTPCodecType

Kind controls if this TrackLocal is audio or video.

#### func (*TrackLocalStaticSample) [RID](https://github.com/pion/webrtc/blob/v4.2.3/track_local_static.go#L270) ¶

    func (s *TrackLocalStaticSample) RID() [string](/builtin#string)

RID is the RTP stream identifier.

#### func (*TrackLocalStaticSample) [StreamID](https://github.com/pion/webrtc/blob/v4.2.3/track_local_static.go#L267) ¶

    func (s *TrackLocalStaticSample) StreamID() [string](/builtin#string)

StreamID is the group this track belongs too. This must be unique.

#### func (*TrackLocalStaticSample) [Unbind](https://github.com/pion/webrtc/blob/v4.2.3/track_local_static.go#L336) ¶

    func (s *TrackLocalStaticSample) Unbind(t TrackLocalContext) [error](/builtin#error)

Unbind implements the teardown logic when the track is no longer needed. This happens because a track has been stopped.

#### func (*TrackLocalStaticSample) [WriteSample](https://github.com/pion/webrtc/blob/v4.2.3/track_local_static.go#L344) ¶

    func (s *TrackLocalStaticSample) WriteSample(sample [media](/github.com/pion/webrtc/v4@v4.2.3/pkg/media).[Sample](/github.com/pion/webrtc/v4@v4.2.3/pkg/media#Sample)) [error](/builtin#error)

WriteSample writes a Sample to the TrackLocalStaticSample If one PeerConnection fails the packets will still be sent to all PeerConnections. The error message will contain the ID of the failed PeerConnections so you can remove them.

#### type [TrackLocalWriter](https://github.com/pion/webrtc/blob/v4.2.3/track_local.go#L12) ¶

    type TrackLocalWriter interface {
     // WriteRTP encrypts a RTP packet and writes to the connection
     WriteRTP(header *[rtp](/github.com/pion/rtp).[Header](/github.com/pion/rtp#Header), payload [][byte](/builtin#byte)) ([int](/builtin#int), [error](/builtin#error))
    
     // Write encrypts and writes a full RTP packet
     Write(b [][byte](/builtin#byte)) ([int](/builtin#int), [error](/builtin#error))
    }

TrackLocalWriter is the Writer for outbound RTP Packets.

#### type [TrackRemote](https://github.com/pion/webrtc/blob/v4.2.3/track_remote.go#L25) ¶

    type TrackRemote struct {
     // contains filtered or unexported fields
    }

TrackRemote represents a single inbound source of media.

#### func (*TrackRemote) [Codec](https://github.com/pion/webrtc/blob/v4.2.3/track_remote.go#L114) ¶

    func (t *TrackRemote) Codec() RTPCodecParameters

Codec gets the Codec of the track.

#### func (*TrackRemote) [HasRTX](https://github.com/pion/webrtc/blob/v4.2.3/track_remote.go#L236) ¶

    func (t *TrackRemote) HasRTX() [bool](/builtin#bool)

HasRTX returns true if the track has a separate RTX stream.

#### func (*TrackRemote) [ID](https://github.com/pion/webrtc/blob/v4.2.3/track_remote.go#L59) ¶

    func (t *TrackRemote) ID() [string](/builtin#string)

ID is the unique identifier for this Track. This should be unique for the stream, but doesn't have to globally unique. A common example would be 'audio' or 'video' and StreamID would be 'desktop' or 'webcam'.

#### func (*TrackRemote) [Kind](https://github.com/pion/webrtc/blob/v4.2.3/track_remote.go#L85) ¶

    func (t *TrackRemote) Kind() RTPCodecType

Kind gets the Kind of the track.

#### func (*TrackRemote) [Msid](https://github.com/pion/webrtc/blob/v4.2.3/track_remote.go#L109) ¶

    func (t *TrackRemote) Msid() [string](/builtin#string)

Msid gets the Msid of the track.

#### func (*TrackRemote) [PayloadType](https://github.com/pion/webrtc/blob/v4.2.3/track_remote.go#L77) ¶

    func (t *TrackRemote) PayloadType() PayloadType

PayloadType gets the PayloadType of the track.

#### func (*TrackRemote) [RID](https://github.com/pion/webrtc/blob/v4.2.3/track_remote.go#L69) ¶

    func (t *TrackRemote) RID() [string](/builtin#string)

RID gets the RTP Stream ID of this Track With Simulcast you will have multiple tracks with the same ID, but different RID values. In many cases a TrackRemote will not have an RID, so it is important to assert it is non-zero.

#### func (*TrackRemote) [Read](https://github.com/pion/webrtc/blob/v4.2.3/track_remote.go#L122) ¶

    func (t *TrackRemote) Read(b [][byte](/builtin#byte)) (n [int](/builtin#int), attributes [interceptor](/github.com/pion/interceptor).[Attributes](/github.com/pion/interceptor#Attributes), err [error](/builtin#error))

Read reads data from the track.

#### func (*TrackRemote) [ReadRTP](https://github.com/pion/webrtc/blob/v4.2.3/track_remote.go#L188) ¶

    func (t *TrackRemote) ReadRTP() (*[rtp](/github.com/pion/rtp).[Packet](/github.com/pion/rtp#Packet), [interceptor](/github.com/pion/interceptor).[Attributes](/github.com/pion/interceptor#Attributes), [error](/builtin#error))

ReadRTP is a convenience method that wraps Read and unmarshals for you.

#### func (*TrackRemote) [RtxSSRC](https://github.com/pion/webrtc/blob/v4.2.3/track_remote.go#L228) ¶

    func (t *TrackRemote) RtxSSRC() SSRC

RtxSSRC returns the RTX SSRC for a track, or 0 if track does not have a separate RTX stream.

#### func (*TrackRemote) [SSRC](https://github.com/pion/webrtc/blob/v4.2.3/track_remote.go#L101) ¶

    func (t *TrackRemote) SSRC() SSRC

SSRC gets the SSRC of the track.

#### func (*TrackRemote) [SetReadDeadline](https://github.com/pion/webrtc/blob/v4.2.3/track_remote.go#L223) ¶

    func (t *TrackRemote) SetReadDeadline(deadline [time](/time).[Time](/time#Time)) [error](/builtin#error)

SetReadDeadline sets the max amount of time the RTP stream will block before returning. 0 is forever.

#### func (*TrackRemote) [StreamID](https://github.com/pion/webrtc/blob/v4.2.3/track_remote.go#L93) ¶

    func (t *TrackRemote) StreamID() [string](/builtin#string)

StreamID is the group this track belongs too. This must be unique.

#### type [TransportStats](https://github.com/pion/webrtc/blob/v4.2.3/stats.go#L1935) ¶

    type TransportStats struct {
     // Timestamp is the timestamp associated with this object.
     Timestamp StatsTimestamp `json:"timestamp"`
    
     // Type is the object's StatsType
     Type StatsType `json:"type"`
    
     // ID is a unique id that is associated with the component inspected to produce
     // this Stats object. Two Stats objects will have the same ID if they were produced
     // by inspecting the same underlying object.
     ID [string](/builtin#string) `json:"id"`
    
     // PacketsSent represents the total number of packets sent over this transport.
     PacketsSent [uint32](/builtin#uint32) `json:"packetsSent"`
    
     // PacketsReceived represents the total number of packets received on this transport.
     PacketsReceived [uint32](/builtin#uint32) `json:"packetsReceived"`
    
     // BytesSent represents the total number of payload bytes sent on this PeerConnection
     // not including headers or padding.
     BytesSent [uint64](/builtin#uint64) `json:"bytesSent"`
    
     // BytesReceived represents the total number of bytes received on this PeerConnection
     // not including headers or padding.
     BytesReceived [uint64](/builtin#uint64) `json:"bytesReceived"`
    
     // RTCPTransportStatsID is the ID of the transport that gives stats for the RTCP
     // component If RTP and RTCP are not multiplexed and this record has only
     // the RTP component stats.
     RTCPTransportStatsID [string](/builtin#string) `json:"rtcpTransportStatsId"`
    
     // ICERole is set to the current value of the "role" attribute of the underlying
     // DTLSTransport's "iceTransport".
     ICERole ICERole `json:"iceRole"`
    
     // DTLSState is set to the current value of the "state" attribute of the underlying DTLSTransport.
     DTLSState DTLSTransportState `json:"dtlsState"`
    
     // ICEState is set to the current value of the "state" attribute of the underlying
     // RTCIceTransport's "state".
     ICEState ICETransportState `json:"iceState"`
    
     // SelectedCandidatePairID is a unique identifier that is associated to the object
     // that was inspected to produce the ICECandidatePairStats associated with this transport.
     SelectedCandidatePairID [string](/builtin#string) `json:"selectedCandidatePairId"`
    
     // LocalCertificateID is the ID of the CertificateStats for the local certificate.
     // Present only if DTLS is negotiated.
     LocalCertificateID [string](/builtin#string) `json:"localCertificateId"`
    
     // RemoteCertificateID is the ID of the CertificateStats for the remote certificate.
     // Present only if DTLS is negotiated.
     RemoteCertificateID [string](/builtin#string) `json:"remoteCertificateId"`
    
     // DTLSCipher is the descriptive name of the cipher suite used for the DTLS transport,
     // as defined in the "Description" column of the IANA cipher suite registry.
     DTLSCipher [string](/builtin#string) `json:"dtlsCipher"`
    
     // SRTPCipher is the descriptive name of the protection profile used for the SRTP
     // transport, as defined in the "Profile" column of the IANA DTLS-SRTP protection
     // profile registry.
     SRTPCipher [string](/builtin#string) `json:"srtpCipher"`
    }

TransportStats contains transport statistics related to the PeerConnection object.

#### type [VideoReceiverStats](https://github.com/pion/webrtc/blob/v4.2.3/stats.go#L1817) ¶

    type VideoReceiverStats struct {
     // Timestamp is the timestamp associated with this object.
     Timestamp StatsTimestamp `json:"timestamp"`
    
     // Type is the object's StatsType
     Type StatsType `json:"type"`
    
     // ID is a unique id that is associated with the component inspected to produce
     // this Stats object. Two Stats objects will have the same ID if they were produced
     // by inspecting the same underlying object.
     ID [string](/builtin#string) `json:"id"`
    
     // Kind is "video"
     Kind [string](/builtin#string) `json:"kind"`
    
     // FrameWidth represents the width of the last processed frame for this track.
     // Before the first frame is processed this attribute is missing.
     FrameWidth [uint32](/builtin#uint32) `json:"frameWidth"`
    
     // FrameHeight represents the height of the last processed frame for this track.
     // Before the first frame is processed this attribute is missing.
     FrameHeight [uint32](/builtin#uint32) `json:"frameHeight"`
    
     // FramesPerSecond represents the nominal FPS value before the degradation preference
     // is applied. It is the number of complete frames in the last second. For sending
     // tracks it is the current captured FPS and for the receiving tracks it is the
     // current decoding framerate.
     FramesPerSecond [float64](/builtin#float64) `json:"framesPerSecond"`
    
     // EstimatedPlayoutTimestamp is the estimated playout time of this receiver's
     // track. The playout time is the NTP timestamp of the last playable sample that
     // has a known timestamp (from an RTCP SR packet mapping RTP timestamps to NTP
     // timestamps), extrapolated with the time elapsed since it was ready to be played out.
     // This is the "current time" of the track in NTP clock time of the sender and
     // can be present even if there is no audio currently playing.
     //
     // This can be useful for estimating how much audio and video is out of
     // sync for two tracks from the same source:
     //   AudioTrackStats.EstimatedPlayoutTimestamp - VideoTrackStats.EstimatedPlayoutTimestamp
     EstimatedPlayoutTimestamp StatsTimestamp `json:"estimatedPlayoutTimestamp"`
    
     // JitterBufferDelay is the sum of the time, in seconds, each sample takes from
     // the time it is received and to the time it exits the jitter buffer.
     // This increases upon samples exiting, having completed their time in the buffer
     // (incrementing JitterBufferEmittedCount). The average jitter buffer delay can
     // be calculated by dividing the JitterBufferDelay with the JitterBufferEmittedCount.
     JitterBufferDelay [float64](/builtin#float64) `json:"jitterBufferDelay"`
    
     // JitterBufferEmittedCount is the total number of samples that have come out
     // of the jitter buffer (increasing JitterBufferDelay).
     JitterBufferEmittedCount [uint64](/builtin#uint64) `json:"jitterBufferEmittedCount"`
    
     // FramesReceived Represents the total number of complete frames received for
     // this receiver. This metric is incremented when the complete frame is received.
     FramesReceived [uint32](/builtin#uint32) `json:"framesReceived"`
    
     // KeyFramesReceived represents the total number of complete key frames received
     // for this MediaStreamTrack, such as Intra-frames in VP8 [RFC6386] or I-frames
     // in H.264 [RFC6184]. This is a subset of framesReceived. `framesReceived - keyFramesReceived`
     // gives you the number of delta frames received. This metric is incremented when
     // the complete key frame is received. It is not incremented if a partial key
     // frame is received and sent for decoding, i.e., the frame could not be recovered
     // via retransmission or FEC.
     KeyFramesReceived [uint32](/builtin#uint32) `json:"keyFramesReceived"`
    
     // FramesDecoded represents the total number of frames correctly decoded for this
     // SSRC, i.e., frames that would be displayed if no frames are dropped.
     FramesDecoded [uint32](/builtin#uint32) `json:"framesDecoded"`
    
     // FramesDropped is the total number of frames dropped predecode or dropped
     // because the frame missed its display deadline for this receiver's track.
     FramesDropped [uint32](/builtin#uint32) `json:"framesDropped"`
    
     // The cumulative number of partial frames lost. This metric is incremented when
     // the frame is sent to the decoder. If the partial frame is received and recovered
     // via retransmission or FEC before decoding, the FramesReceived counter is incremented.
     PartialFramesLost [uint32](/builtin#uint32) `json:"partialFramesLost"`
    
     // FullFramesLost is the cumulative number of full frames lost.
     FullFramesLost [uint32](/builtin#uint32) `json:"fullFramesLost"`
    }

VideoReceiverStats contains video metrics related to a specific receiver.

#### type [VideoSenderStats](https://github.com/pion/webrtc/blob/v4.2.3/stats.go#L1598) ¶

    type VideoSenderStats struct {
     // Timestamp is the timestamp associated with this object.
     Timestamp StatsTimestamp `json:"timestamp"`
    
     // Type is the object's StatsType
     Type StatsType `json:"type"`
    
     // ID is a unique id that is associated with the component inspected to produce
     // this Stats object. Two Stats objects will have the same ID if they were produced
     // by inspecting the same underlying object.
     ID [string](/builtin#string) `json:"id"`
    
     // Kind is "video"
     Kind [string](/builtin#string) `json:"kind"`
    
     // FramesCaptured represents the total number of frames captured, before encoding,
     // for this RTPSender (or for this MediaStreamTrack, if type is "track"). For example,
     // if type is "sender" and this sender's track represents a camera, then this is the
     // number of frames produced by the camera for this track while being sent by this sender,
     // combined with the number of frames produced by all tracks previously attached to this
     // sender while being sent by this sender. Framerates can vary due to hardware limitations
     // or environmental factors such as lighting conditions.
     FramesCaptured [uint32](/builtin#uint32) `json:"framesCaptured"`
    
     // FramesSent represents the total number of frames sent by this RTPSender
     // (or for this MediaStreamTrack, if type is "track").
     FramesSent [uint32](/builtin#uint32) `json:"framesSent"`
    
     // HugeFramesSent represents the total number of huge frames sent by this RTPSender
     // (or for this MediaStreamTrack, if type is "track"). Huge frames, by definition,
     // are frames that have an encoded size at least 2.5 times the average size of the frames.
     // The average size of the frames is defined as the target bitrate per second divided
     // by the target fps at the time the frame was encoded. These are usually complex
     // to encode frames with a lot of changes in the picture. This can be used to estimate,
     // e.g slide changes in the streamed presentation. If a huge frame is also a key frame,
     // then both counters HugeFramesSent and KeyFramesSent are incremented.
     HugeFramesSent [uint32](/builtin#uint32) `json:"hugeFramesSent"`
    
     // KeyFramesSent represents the total number of key frames sent by this RTPSender
     // (or for this MediaStreamTrack, if type is "track"), such as Infra-frames in
     // VP8 [RFC6386] or I-frames in H.264 [RFC6184]. This is a subset of FramesSent.
     // FramesSent - KeyFramesSent gives you the number of delta frames sent.
     KeyFramesSent [uint32](/builtin#uint32) `json:"keyFramesSent"`
    }

VideoSenderStats represents the stats about one video sender of a PeerConnection object for which one calls GetStats.

It appears in the stats as soon as the sender is added by either AddTrack or AddTransceiver, or by media negotiation.

#### type [VideoSourceStats](https://github.com/pion/webrtc/blob/v4.2.3/stats.go#L1242) ¶

    type VideoSourceStats struct {
     // Timestamp is the timestamp associated with this object.
     Timestamp StatsTimestamp `json:"timestamp"`
    
     // Type is the object's StatsType
     Type StatsType `json:"type"`
    
     // ID is a unique id that is associated with the component inspected to produce
     // this Stats object. Two Stats objects will have the same ID if they were produced
     // by inspecting the same underlying object.
     ID [string](/builtin#string) `json:"id"`
    
     // TrackIdentifier represents the id property of the track.
     TrackIdentifier [string](/builtin#string) `json:"trackIdentifier"`
    
     // Kind is "video"
     Kind [string](/builtin#string) `json:"kind"`
    
     // Width is width of the last frame originating from this source in pixels.
     Width [uint32](/builtin#uint32) `json:"width"`
    
     // Height is height of the last frame originating from this source in pixels.
     Height [uint32](/builtin#uint32) `json:"height"`
    
     // Frames is the total number of frames originating from this source.
     Frames [uint32](/builtin#uint32) `json:"frames"`
    
     // FramesPerSecond is the number of frames originating from this source, measured during the last second.
     FramesPerSecond [float64](/builtin#float64) `json:"framesPerSecond"`
    }

VideoSourceStats represents a video track that is attached to one or more senders.
  *[↑]: Back to Top
