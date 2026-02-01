# WebRTC API

> Source: https://developer.mozilla.org/en-US/docs/Web/API/WebRTC_API
> Fetched: 2026-02-01T11:46:42.257930+00:00
> Content-Hash: b1975205e41b4120
> Type: html

---

# WebRTC API

**WebRTC** (Web Real-Time Communication) is a technology that enables Web applications and sites to capture and optionally stream audio and/or video media, as well as to exchange arbitrary data between browsers without requiring an intermediary. The set of standards that comprise WebRTC makes it possible to share data and perform teleconferencing peer-to-peer, without requiring that the user install plug-ins or any other third-party software.

WebRTC consists of several interrelated APIs and protocols which work together to achieve this. The documentation you'll find here will help you understand the fundamentals of WebRTC, how to set up and use both data and media connections, and more.

## WebRTC concepts and usage

WebRTC serves multiple purposes; together with the [Media Capture and Streams API](/en-US/docs/Web/API/Media_Capture_and_Streams_API), they provide powerful multimedia capabilities to the Web, including support for audio and video conferencing, file exchange, screen sharing, identity management, and interfacing with legacy telephone systems including support for sending [DTMF](/en-US/docs/Glossary/DTMF) (touch-tone dialing) signals. Connections between peers can be made without requiring any special drivers or plug-ins, and can often be made without any intermediary servers.

Connections between two peers are represented by the [`RTCPeerConnection`](/en-US/docs/Web/API/RTCPeerConnection) interface. Once a connection has been established and opened using `RTCPeerConnection`, media streams ([`MediaStream`](/en-US/docs/Web/API/MediaStream)s) and/or data channels ([`RTCDataChannel`](/en-US/docs/Web/API/RTCDataChannel)s) can be added to the connection.

Media streams can consist of any number of tracks of media information; tracks, which are represented by objects based on the [`MediaStreamTrack`](/en-US/docs/Web/API/MediaStreamTrack) interface, may contain one of a number of types of media data, including audio, video, and text (such as subtitles or even chapter names). Most streams consist of at least one audio track and likely also a video track, and can be used to send and receive both live media or stored media information (such as a streamed movie).

You can also use the connection between two peers to exchange arbitrary binary data using the [`RTCDataChannel`](/en-US/docs/Web/API/RTCDataChannel) interface. This can be used for back-channel information, metadata exchange, game status packets, file transfers, or even as a primary channel for data transfer.

### Interoperability

WebRTC is in general well supported in modern browsers, but some incompatibilities remain. The [adapter.js](https://github.com/webrtcHacks/adapter) library is a shim to insulate apps from these incompatibilities.

## WebRTC reference

Because WebRTC provides interfaces that work together to accomplish a variety of tasks, we have divided up the reference by category. Please see the sidebar for an alphabetical list.

### Connection setup and management

These interfaces, dictionaries, and types are used to set up, open, and manage WebRTC connections. Included are interfaces representing peer media connections, data channels, and interfaces used when exchanging information on the capabilities of each peer in order to select the best possible configuration for a two-way media connection.

#### Interfaces

[`RTCPeerConnection`](/en-US/docs/Web/API/RTCPeerConnection)

Represents a WebRTC connection between the local computer and a remote peer. It is used to handle efficient streaming of data between the two peers.

[`RTCDataChannel`](/en-US/docs/Web/API/RTCDataChannel)

Represents a bi-directional data channel between two peers of a connection.

[`RTCDataChannelEvent`](/en-US/docs/Web/API/RTCDataChannelEvent)

Represents events that occur while attaching a [`RTCDataChannel`](/en-US/docs/Web/API/RTCDataChannel) to a [`RTCPeerConnection`](/en-US/docs/Web/API/RTCPeerConnection). The only event sent with this interface is [`datachannel`](/en-US/docs/Web/API/RTCPeerConnection/datachannel_event "datachannel").

[`RTCSessionDescription`](/en-US/docs/Web/API/RTCSessionDescription)

Represents the parameters of a session. Each `RTCSessionDescription` consists of a description [`type`](/en-US/docs/Web/API/RTCSessionDescription/type "type") indicating which part of the offer/answer negotiation process it describes and of the [SDP](/en-US/docs/Glossary/SDP) descriptor of the session.

[`RTCStatsReport`](/en-US/docs/Web/API/RTCStatsReport)

Provides information detailing statistics for a connection or for an individual track on the connection; the report can be obtained by calling [`RTCPeerConnection.getStats()`](/en-US/docs/Web/API/RTCPeerConnection/getStats).

[`RTCIceCandidate`](/en-US/docs/Web/API/RTCIceCandidate)

Represents a candidate Interactive Connectivity Establishment ([ICE](/en-US/docs/Glossary/ICE)) server for establishing an [`RTCPeerConnection`](/en-US/docs/Web/API/RTCPeerConnection).

[`RTCIceTransport`](/en-US/docs/Web/API/RTCIceTransport)

Represents information about an [ICE](/en-US/docs/Glossary/ICE) transport.

[`RTCPeerConnectionIceEvent`](/en-US/docs/Web/API/RTCPeerConnectionIceEvent)

Represents events that occur in relation to ICE candidates with the target, usually an [`RTCPeerConnection`](/en-US/docs/Web/API/RTCPeerConnection). Only one event is of this type: [`icecandidate`](/en-US/docs/Web/API/RTCPeerConnection/icecandidate_event "icecandidate").

[`RTCRtpSender`](/en-US/docs/Web/API/RTCRtpSender)

Manages the encoding and transmission of data for a [`MediaStreamTrack`](/en-US/docs/Web/API/MediaStreamTrack) on an [`RTCPeerConnection`](/en-US/docs/Web/API/RTCPeerConnection).

[`RTCRtpReceiver`](/en-US/docs/Web/API/RTCRtpReceiver)

Manages the reception and decoding of data for a [`MediaStreamTrack`](/en-US/docs/Web/API/MediaStreamTrack) on an [`RTCPeerConnection`](/en-US/docs/Web/API/RTCPeerConnection).

[`RTCTrackEvent`](/en-US/docs/Web/API/RTCTrackEvent)

The interface used to represent a [`track`](/en-US/docs/Web/API/RTCPeerConnection/track_event "track") event, which indicates that an [`RTCRtpReceiver`](/en-US/docs/Web/API/RTCRtpReceiver) object was added to the [`RTCPeerConnection`](/en-US/docs/Web/API/RTCPeerConnection) object, indicating that a new incoming [`MediaStreamTrack`](/en-US/docs/Web/API/MediaStreamTrack) was created and added to the `RTCPeerConnection`.

[`RTCSctpTransport`](/en-US/docs/Web/API/RTCSctpTransport)

Provides information which describes a Stream Control Transmission Protocol (**[SCTP](/en-US/docs/Glossary/SCTP)**) transport and also provides a way to access the underlying Datagram Transport Layer Security (**[DTLS](/en-US/docs/Glossary/DTLS)**) transport over which SCTP packets for all of an [`RTCPeerConnection`](/en-US/docs/Web/API/RTCPeerConnection)'s data channels are sent and received.

#### Events

[`bufferedamountlow`](/en-US/docs/Web/API/RTCDataChannel/bufferedamountlow_event "bufferedamountlow")

The amount of data currently buffered by the data channelâas indicated by its [`bufferedAmount`](/en-US/docs/Web/API/RTCDataChannel/bufferedAmount "bufferedAmount") propertyâhas decreased to be at or below the channel's minimum buffered data size, as specified by [`bufferedAmountLowThreshold`](/en-US/docs/Web/API/RTCDataChannel/bufferedAmountLowThreshold "bufferedAmountLowThreshold").

[`close`](/en-US/docs/Web/API/RTCDataChannel/close_event "close")

The data channel has completed the closing process and is now in the `closed` state. Its underlying data transport is completely closed at this point. You can be notified _before_ closing completes by watching for the `closing` event instead.

[`closing`](/en-US/docs/Web/API/RTCDataChannel/closing_event "closing")

The `RTCDataChannel` has transitioned to the `closing` state, indicating that it will be closed soon. You can detect the completion of the closing process by watching for the `close` event.

[`connectionstatechange`](/en-US/docs/Web/API/RTCPeerConnection/connectionstatechange_event "connectionstatechange")

The connection's state, which can be accessed in [`connectionState`](/en-US/docs/Web/API/RTCPeerConnection/connectionState "connectionState"), has changed.

[`datachannel`](/en-US/docs/Web/API/RTCPeerConnection/datachannel_event "datachannel")

A new [`RTCDataChannel`](/en-US/docs/Web/API/RTCDataChannel) is available following the remote peer opening a new data channel. This event's type is [`RTCDataChannelEvent`](/en-US/docs/Web/API/RTCDataChannelEvent).

[`error`](/en-US/docs/Web/API/RTCDataChannel/error_event "error")

An [`RTCErrorEvent`](/en-US/docs/Web/API/RTCErrorEvent) indicating that an error occurred on the data channel.

[`error`](/en-US/docs/Web/API/RTCDtlsTransport/error_event "error")

An [`RTCErrorEvent`](/en-US/docs/Web/API/RTCErrorEvent) indicating that an error occurred on the [`RTCDtlsTransport`](/en-US/docs/Web/API/RTCDtlsTransport). This error will be either `dtls-failure` or `fingerprint-failure`.

[`gatheringstatechange`](/en-US/docs/Web/API/RTCIceTransport/gatheringstatechange_event "gatheringstatechange")

The [`RTCIceTransport`](/en-US/docs/Web/API/RTCIceTransport)'s gathering state has changed.

[`icecandidate`](/en-US/docs/Web/API/RTCPeerConnection/icecandidate_event "icecandidate")

An [`RTCPeerConnectionIceEvent`](/en-US/docs/Web/API/RTCPeerConnectionIceEvent) which is sent whenever the local device has identified a new ICE candidate which needs to be added to the local peer by calling [`setLocalDescription()`](/en-US/docs/Web/API/RTCPeerConnection/setLocalDescription "setLocalDescription\(\)").

[`icecandidateerror`](/en-US/docs/Web/API/RTCPeerConnection/icecandidateerror_event "icecandidateerror")

An [`RTCPeerConnectionIceErrorEvent`](/en-US/docs/Web/API/RTCPeerConnectionIceErrorEvent) indicating that an error has occurred while gathering ICE candidates.

[`iceconnectionstatechange`](/en-US/docs/Web/API/RTCPeerConnection/iceconnectionstatechange_event "iceconnectionstatechange")

Sent to an [`RTCPeerConnection`](/en-US/docs/Web/API/RTCPeerConnection) when its ICE connection's stateâfound in the [`iceConnectionState`](/en-US/docs/Web/API/RTCPeerConnection/iceConnectionState "iceConnectionState") propertyâchanges.

[`icegatheringstatechange`](/en-US/docs/Web/API/RTCPeerConnection/icegatheringstatechange_event "icegatheringstatechange")

Sent to an [`RTCPeerConnection`](/en-US/docs/Web/API/RTCPeerConnection) when its ICE gathering stateâfound in the [`iceGatheringState`](/en-US/docs/Web/API/RTCPeerConnection/iceGatheringState "iceGatheringState") propertyâchanges.

[`message`](/en-US/docs/Web/API/RTCDataChannel/message_event "message")

A message has been received on the data channel. The event is of type [`MessageEvent`](/en-US/docs/Web/API/MessageEvent).

[`negotiationneeded`](/en-US/docs/Web/API/RTCPeerConnection/negotiationneeded_event "negotiationneeded")

Informs the `RTCPeerConnection` that it needs to perform session negotiation by calling [`createOffer()`](/en-US/docs/Web/API/RTCPeerConnection/createOffer "createOffer\(\)") followed by [`setLocalDescription()`](/en-US/docs/Web/API/RTCPeerConnection/setLocalDescription "setLocalDescription\(\)").

[`open`](/en-US/docs/Web/API/RTCDataChannel/open_event "open")

The underlying data transport for the `RTCDataChannel` has been successfully opened or re-opened.

[`selectedcandidatepairchange`](/en-US/docs/Web/API/RTCIceTransport/selectedcandidatepairchange_event "selectedcandidatepairchange")

The currently-selected pair of ICE candidates has changed for the `RTCIceTransport` on which the event is fired.

[`track`](/en-US/docs/Web/API/RTCPeerConnection/track_event "track")

The `track` event, of type [`RTCTrackEvent`](/en-US/docs/Web/API/RTCTrackEvent) is sent to an [`RTCPeerConnection`](/en-US/docs/Web/API/RTCPeerConnection) when a new track is added to the connection following the successful negotiation of the media's streaming.

[`signalingstatechange`](/en-US/docs/Web/API/RTCPeerConnection/signalingstatechange_event "signalingstatechange")

Sent to the peer connection when its [`signalingState`](/en-US/docs/Web/API/RTCPeerConnection/signalingState "signalingState") has changed. This happens as a result of a call to either [`setLocalDescription()`](/en-US/docs/Web/API/RTCPeerConnection/setLocalDescription "setLocalDescription\(\)") or [`setRemoteDescription()`](/en-US/docs/Web/API/RTCPeerConnection/setRemoteDescription "setRemoteDescription\(\)").

[`statechange`](/en-US/docs/Web/API/RTCDtlsTransport/statechange_event "statechange")

The state of the `RTCDtlsTransport` has changed.

[`statechange`](/en-US/docs/Web/API/RTCIceTransport/statechange_event "statechange")

The state of the `RTCIceTransport` has changed.

[`statechange`](/en-US/docs/Web/API/RTCSctpTransport/statechange_event "statechange")

The state of the `RTCSctpTransport` has changed.

[`rtctransform`](/en-US/docs/Web/API/DedicatedWorkerGlobalScope/rtctransform_event "rtctransform")

An encoded video or audio frame is ready to process using a transform stream in a worker.

#### Types

[`RTCSctpTransport.state`](/en-US/docs/Web/API/RTCSctpTransport/state)

Indicates the state of an [`RTCSctpTransport`](/en-US/docs/Web/API/RTCSctpTransport) instance.

### Identity and security

These APIs are used to manage user identity and security, in order to authenticate the user for a connection.

`RTCIdentityProvider`

Enables a user agent is able to request that an identity assertion be generated or validated.

[`RTCIdentityAssertion`](/en-US/docs/Web/API/RTCIdentityAssertion)

Represents the identity of the remote peer of the current connection. If no peer has yet been set and verified this interface returns `null`. Once set it can't be changed.

`RTCIdentityProviderRegistrar`

Registers an identity provider (idP).

[`RTCCertificate`](/en-US/docs/Web/API/RTCCertificate)

Represents a certificate that an [`RTCPeerConnection`](/en-US/docs/Web/API/RTCPeerConnection) uses to authenticate.

### Telephony

These interfaces and events are related to interactivity with Public-Switched Telephone Networks (PSTNs). They're primarily used to send tone dialing soundsâor packets representing those tonesâacross the network to the remote peer.

#### Interfaces

[`RTCDTMFSender`](/en-US/docs/Web/API/RTCDTMFSender)

Manages the encoding and transmission of Dual-Tone Multi-Frequency ([DTMF](/en-US/docs/Glossary/DTMF)) signaling for an [`RTCPeerConnection`](/en-US/docs/Web/API/RTCPeerConnection).

[`RTCDTMFToneChangeEvent`](/en-US/docs/Web/API/RTCDTMFToneChangeEvent)

Used by the [`tonechange`](/en-US/docs/Web/API/RTCDTMFSender/tonechange_event "tonechange") event to indicate that a DTMF tone has either begun or ended. This event does not bubble (except where otherwise stated) and is not cancelable (except where otherwise stated).

#### Events

[`tonechange`](/en-US/docs/Web/API/RTCDTMFSender/tonechange_event "tonechange")

Either a new [DTMF](/en-US/docs/Glossary/DTMF) tone has begun to play over the connection, or the last tone in the `RTCDTMFSender`'s [`toneBuffer`](/en-US/docs/Web/API/RTCDTMFSender/toneBuffer "toneBuffer") has been sent and the buffer is now empty. The event's type is [`RTCDTMFToneChangeEvent`](/en-US/docs/Web/API/RTCDTMFToneChangeEvent).

### Encoded Transforms

These interfaces and events are used to process incoming and outgoing encoded video and audio frames using a transform stream running in a worker.

#### Interfaces

[`RTCRtpScriptTransform`](/en-US/docs/Web/API/RTCRtpScriptTransform)

An interface for inserting transform stream(s) running in a worker into the RTC pipeline.

[`RTCRtpScriptTransformer`](/en-US/docs/Web/API/RTCRtpScriptTransformer)

The worker-side counterpart of an `RTCRtpScriptTransform` that passes options from the main thread, along with a readable stream and writeable stream that can be used to pipe encoded frames through a [`TransformStream`](/en-US/docs/Web/API/TransformStream).

[`RTCEncodedVideoFrame`](/en-US/docs/Web/API/RTCEncodedVideoFrame)

Represents an encoded video frame to be transformed in the RTC pipeline.

[`RTCEncodedAudioFrame`](/en-US/docs/Web/API/RTCEncodedAudioFrame)

Represents an encoded audio frame to be transformed in the RTC pipeline.

#### Properties

[`RTCRtpReceiver.transform`](/en-US/docs/Web/API/RTCRtpReceiver/transform)

A property used to insert a transform stream into the receiver pipeline for incoming encoded video and audio frames.

[`RTCRtpSender.transform`](/en-US/docs/Web/API/RTCRtpSender/transform)

A property used to insert a transform stream into the sender pipeline for outgoing encoded video and audio frames.

#### Events

[`rtctransform`](/en-US/docs/Web/API/DedicatedWorkerGlobalScope/rtctransform_event "rtctransform")

An RTC transform is ready to run in the worker, or an encoded video or audio frame is ready to process.

## Guides

[Introduction to the Real-time Transport Protocol (RTP)](/en-US/docs/Web/API/WebRTC_API/Intro_to_RTP)

The Real-time Transport Protocol (RTP), defined in [RFC 3550](https://datatracker.ietf.org/doc/html/rfc3550), is an IETF standard protocol to enable real-time connectivity for exchanging data that needs real-time priority. This article provides an overview of what RTP is and how it functions in the context of WebRTC.

[Introduction to WebRTC protocols](/en-US/docs/Web/API/WebRTC_API/Protocols)

This article introduces the protocols on top of which the WebRTC API is built.

[WebRTC connectivity](/en-US/docs/Web/API/WebRTC_API/Connectivity)

A guide to how WebRTC connections work and how the various protocols and interfaces can be used together to build powerful communication apps.

[Lifetime of a WebRTC session](/en-US/docs/Web/API/WebRTC_API/Session_lifetime)

WebRTC lets you build peer-to-peer communication of arbitrary data, audio, or videoâor any combination thereofâinto a browser application. In this article, we'll look at the lifetime of a WebRTC session, from establishing the connection all the way through closing the connection when it's no longer needed.

[Establishing a connection: The perfect negotiation pattern](/en-US/docs/Web/API/WebRTC_API/Perfect_negotiation)

**Perfect negotiation** is a design pattern which is recommended for your signaling process to follow, which provides transparency in negotiation while allowing both sides to be either the offerer or the answerer, without significant coding needed to differentiate the two.

[Signaling and two-way video calling](/en-US/docs/Web/API/WebRTC_API/Signaling_and_video_calling)

A tutorial and example which turns a WebSocket-based chat system created for a previous example and adds support for opening video calls among participants. The chat server's WebSocket connection is used for WebRTC signaling.

[Codecs used by WebRTC](/en-US/docs/Web/Media/Guides/Formats/WebRTC_codecs)

A guide to the codecs which WebRTC requires browsers to support as well as the optional ones supported by various popular browsers. Included is a guide to help you choose the best codecs for your needs.

[Using WebRTC data channels](/en-US/docs/Web/API/WebRTC_API/Using_data_channels)

This guide covers how you can use a peer connection and an associated [`RTCDataChannel`](/en-US/docs/Web/API/RTCDataChannel) to exchange arbitrary data between two peers.

[Using DTMF with WebRTC](/en-US/docs/Web/API/WebRTC_API/Using_DTMF)

WebRTC's support for interacting with gateways that link to old-school telephone systems includes support for sending DTMF tones using the [`RTCDTMFSender`](/en-US/docs/Web/API/RTCDTMFSender) interface. This guide shows how to do so.

[Using WebRTC Encoded Transforms](/en-US/docs/Web/API/WebRTC_API/Using_Encoded_Transforms)

This guide shows how a web application can modify incoming and outgoing WebRTC encoded video and audio frames, using a [`TransformStream`](/en-US/docs/Web/API/TransformStream) running into a worker.

## Tutorials

Improving compatibility using WebRTC adapter.js

The WebRTC organization [provides on GitHub the WebRTC adapter](https://github.com/webrtc/adapter/) to work around compatibility issues in different browsers' WebRTC implementations. The adapter is a JavaScript shim which lets your code to be written to the specification so that it will "just work" in all browsers with WebRTC support.

[A simple RTCDataChannel sample](/en-US/docs/Web/API/WebRTC_API/Simple_RTCDataChannel_sample)

The [`RTCDataChannel`](/en-US/docs/Web/API/RTCDataChannel) interface is a feature which lets you open a channel between two peers over which you may send and receive arbitrary data. The API is intentionally similar to the [WebSocket API](/en-US/docs/Web/API/WebSockets_API), so that the same programming model can be used for each.

[Building an internet connected phone with Peer.js](/en-US/docs/Web/API/WebRTC_API/Build_a_phone_with_peerjs)

This tutorial is a step-by-step guide on how to build a phone using Peer.js

## Specifications

Specification  
---  
[WebRTC: Real-Time Communication in Browsers](https://w3c.github.io/webrtc-pc/)  
[Media Capture and Streams](https://w3c.github.io/mediacapture-main/)  
[Media Capture from DOM Elements](https://w3c.github.io/mediacapture-fromelement/)  
  
### WebRTC-proper protocols

- [Application Layer Protocol Negotiation for Web Real-Time Communications](https://datatracker.ietf.org/doc/rfc8833/)
- [WebRTC Audio Codec and Processing Requirements](https://datatracker.ietf.org/doc/rfc7874/)
- [RTCWeb Data Channels](https://datatracker.ietf.org/doc/rfc8831/)
- [RTCWeb Data Channel Protocol](https://datatracker.ietf.org/doc/rfc8832/)
- [Web Real-Time Communication (WebRTC): Media Transport and Use of RTP](https://datatracker.ietf.org/doc/rfc8834/)
- [WebRTC Security Architecture](https://datatracker.ietf.org/doc/rfc8827/)
- [Transports for RTCWEB](https://datatracker.ietf.org/doc/rfc8835/)

### Related supporting protocols

- [Interactive Connectivity Establishment (ICE): A Protocol for Network Address Translator (NAT) Traversal for Offer/Answer Protocol](https://datatracker.ietf.org/doc/html/rfc5245)
- [Session Traversal Utilities for NAT (STUN)](https://datatracker.ietf.org/doc/html/rfc5389)
- [URI Scheme for the Session Traversal Utilities for NAT (STUN) Protocol](https://datatracker.ietf.org/doc/html/rfc7064)
- [Traversal Using Relays around NAT (TURN) Uniform Resource Identifiers](https://datatracker.ietf.org/doc/html/rfc7065)
- [An Offer/Answer Model with Session Description Protocol (SDP)](https://datatracker.ietf.org/doc/html/rfc3264)
- [Session Traversal Utilities for NAT (STUN) Extension for Third Party Authorization](https://datatracker.ietf.org/doc/rfc7635/)

## See also

- [`MediaDevices`](/en-US/docs/Web/API/MediaDevices)
- [`MediaStreamEvent`](/en-US/docs/Web/API/MediaStreamEvent)
- [`MediaStreamTrack`](/en-US/docs/Web/API/MediaStreamTrack)
- [`MessageEvent`](/en-US/docs/Web/API/MessageEvent)
- [`MediaStream`](/en-US/docs/Web/API/MediaStream)
- [Media Capture and Streams API](/en-US/docs/Web/API/Media_Capture_and_Streams_API)
- [Firefox multistream and renegotiation for Jitsi Videobridge](https://hacks.mozilla.org/2015/06/firefox-multistream-and-renegotiation-for-jitsi-videobridge/)
- [Peering Through the WebRTC Fog with SocketPeer](https://hacks.mozilla.org/2015/04/peering-through-the-webrtc-fog-with-socketpeer/)
- [Inside the Party Bus: Building a Web App with Multiple Live Video Streams + Interactive Graphics](https://hacks.mozilla.org/2014/04/inside-the-party-bus-building-a-web-app-with-multiple-live-video-streams-interactive-graphics/)
- [Web media technologies](/en-US/docs/Web/Media)

## Help improve MDN

Was this page helpful to you?

Yes No

[Learn how to contribute](/en-US/docs/MDN/Community/Getting_started)

This page was last modified on Jun 26, 2025 by [MDN contributors](/en-US/docs/Web/API/WebRTC_API/contributors.txt).

[View this page on GitHub](https://github.com/mdn/content/blob/main/files/en-us/web/api/webrtc_api/index.md?plain=1 "Folder: en-us/web/api/webrtc_api \(Opens in a new tab\)") â¢ [Report a problem with this content](https://github.com/mdn/content/issues/new?template=page-report.yml&mdn-url=https%3A%2F%2Fdeveloper.mozilla.org%2Fen-US%2Fdocs%2FWeb%2FAPI%2FWebRTC_API&metadata=%3C%21--+Do+not+make+changes+below+this+line+--%3E%0A%3Cdetails%3E%0A%3Csummary%3EPage+report+details%3C%2Fsummary%3E%0A%0A*+Folder%3A+%60en-us%2Fweb%2Fapi%2Fwebrtc_api%60%0A*+MDN+URL%3A+https%3A%2F%2Fdeveloper.mozilla.org%2Fen-US%2Fdocs%2FWeb%2FAPI%2FWebRTC_API%0A*+GitHub+URL%3A+https%3A%2F%2Fgithub.com%2Fmdn%2Fcontent%2Fblob%2Fmain%2Ffiles%2Fen-us%2Fweb%2Fapi%2Fwebrtc_api%2Findex.md%0A*+Last+commit%3A+https%3A%2F%2Fgithub.com%2Fmdn%2Fcontent%2Fcommit%2Fd0ed4906719465102739e604bdb35213fb19f251%0A*+Document+last+modified%3A+2025-06-26T01%3A52%3A11.000Z%0A%0A%3C%2Fdetails%3E "This will take you to GitHub to file a new issue.")
  *[↑]: Back to Top
