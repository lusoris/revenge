# gohlslib (HLS)

> Auto-fetched from [https://pkg.go.dev/github.com/bluenviron/gohlslib/v2](https://pkg.go.dev/github.com/bluenviron/gohlslib/v2)
> Last Updated: 2026-01-29T20:14:38.874237+00:00

---

Overview
¶
Package gohlslib is a HLS client and muxer library for the Go programming language.
Examples are available at
https://github.com/bluenviron/gohlslib/tree/main/examples
Index
¶
Variables
type Client
func (c *Client) AbsoluteTime(track *Track) (time.Time, bool)
func (c *Client) Close()
func (c *Client) OnDataAV1(track *Track, cb ClientOnDataAV1Func)
func (c *Client) OnDataH26x(track *Track, cb ClientOnDataH26xFunc)
func (c *Client) OnDataMPEG4Audio(track *Track, cb ClientOnDataMPEG4AudioFunc)
func (c *Client) OnDataOpus(track *Track, cb ClientOnDataOpusFunc)
func (c *Client) OnDataVP9(track *Track, cb ClientOnDataVP9Func)
func (c *Client) Start() error
func (c *Client) Wait() chan error
deprecated
func (c *Client) Wait2() error
type ClientOnDataAV1Func
type ClientOnDataH26xFunc
type ClientOnDataMPEG4AudioFunc
type ClientOnDataOpusFunc
type ClientOnDataVP9Func
type ClientOnDecodeErrorFunc
type ClientOnDownloadPartFunc
type ClientOnDownloadPrimaryPlaylistFunc
type ClientOnDownloadSegmentFunc
type ClientOnDownloadStreamPlaylistFunc
type ClientOnRequestFunc
type ClientOnTracksFunc
type Muxer
func (m *Muxer) Close()
func (m *Muxer) Handle(w http.ResponseWriter, r *http.Request)
func (m *Muxer) Start() error
func (m *Muxer) WriteAV1(track *Track, ntp time.Time, pts int64, tu [][]byte) error
func (m *Muxer) WriteH264(track *Track, ntp time.Time, pts int64, au [][]byte) error
func (m *Muxer) WriteH265(track *Track, ntp time.Time, pts int64, au [][]byte) error
func (m *Muxer) WriteMPEG4Audio(track *Track, ntp time.Time, pts int64, aus [][]byte) error
func (m *Muxer) WriteOpus(track *Track, ntp time.Time, pts int64, packets [][]byte) error
func (m *Muxer) WriteVP9(track *Track, ntp time.Time, pts int64, frame []byte) error
type MuxerOnEncodeErrorFunc
type MuxerVariant
type Track
Constants
¶
This section is empty.
Variables
¶
View Source
var ErrClientEOS =
errors
.
New
("end of stream")
ErrClientEOS is returned by Wait() when the stream has ended.
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
// parameters (all optional except URI)
//
// URI of the playlist.
URI
string
// Start distance from the end of the playlist,
// expressed as number of segments.
// It defaults to 3.
StartDistance
int
// Maximum distance from the end of the playlist,
// expressed as number of segments.
// It defaults to 5.
MaxDistance
int
// HTTP client.
// It defaults to http.DefaultClient.
HTTPClient *
http
.
Client
//
// callbacks (all optional)
//
// called when sending a request to the server.
OnRequest
ClientOnRequestFunc
// called when tracks are available.
OnTracks
ClientOnTracksFunc
// called before downloading a primary playlist.
OnDownloadPrimaryPlaylist
ClientOnDownloadPrimaryPlaylistFunc
// called before downloading a stream playlist.
OnDownloadStreamPlaylist
ClientOnDownloadStreamPlaylistFunc
// called before downloading a segment.
OnDownloadSegment
ClientOnDownloadSegmentFunc
// called before downloading a part.
OnDownloadPart
ClientOnDownloadPartFunc
// called when a non-fatal decode error occurs.
OnDecodeError
ClientOnDecodeErrorFunc
// contains filtered or unexported fields
}
Client is a HLS client.
func (*Client)
AbsoluteTime
¶
func (c *
Client
) AbsoluteTime(track *
Track
) (
time
.
Time
,
bool
)
AbsoluteTime returns the absolute timestamp of the last sample.
func (*Client)
Close
¶
func (c *
Client
) Close()
Close closes all the Client resources and waits for them to exit.
func (*Client)
OnDataAV1
¶
func (c *
Client
) OnDataAV1(track *
Track
, cb
ClientOnDataAV1Func
)
OnDataAV1 sets a callback that is called when data from an AV1 track is received.
func (*Client)
OnDataH26x
¶
func (c *
Client
) OnDataH26x(track *
Track
, cb
ClientOnDataH26xFunc
)
OnDataH26x sets a callback that is called when data from an H26x track is received.
func (*Client)
OnDataMPEG4Audio
¶
func (c *
Client
) OnDataMPEG4Audio(track *
Track
, cb
ClientOnDataMPEG4AudioFunc
)
OnDataMPEG4Audio sets a callback that is called when data from a MPEG-4 Audio track is received.
func (*Client)
OnDataOpus
¶
func (c *
Client
) OnDataOpus(track *
Track
, cb
ClientOnDataOpusFunc
)
OnDataOpus sets a callback that is called when data from an Opus track is received.
func (*Client)
OnDataVP9
¶
func (c *
Client
) OnDataVP9(track *
Track
, cb
ClientOnDataVP9Func
)
OnDataVP9 sets a callback that is called when data from a VP9 track is received.
func (*Client)
Start
¶
func (c *
Client
) Start()
error
Start starts the client.
func (*Client)
Wait
deprecated
func (c *
Client
) Wait() chan
error
Wait waits for any error of the Client.
Deprecated: replaced by Wait2.
func (*Client)
Wait2
¶
added in
v2.2.0
func (c *
Client
) Wait2()
error
Wait2 waits until all client resources are closed.
This can happen when a fatal error occurs or when Close() is called.
type
ClientOnDataAV1Func
¶
type ClientOnDataAV1Func func(pts
int64
, tu [][]
byte
)
ClientOnDataAV1Func is the prototype of the function passed to OnDataAV1().
type
ClientOnDataH26xFunc
¶
type ClientOnDataH26xFunc func(pts
int64
, dts
int64
, au [][]
byte
)
ClientOnDataH26xFunc is the prototype of the function passed to OnDataH26x().
type
ClientOnDataMPEG4AudioFunc
¶
type ClientOnDataMPEG4AudioFunc func(pts
int64
, aus [][]
byte
)
ClientOnDataMPEG4AudioFunc is the prototype of the function passed to OnDataMPEG4Audio().
type
ClientOnDataOpusFunc
¶
type ClientOnDataOpusFunc func(pts
int64
, packets [][]
byte
)
ClientOnDataOpusFunc is the prototype of the function passed to OnDataOpus().
type
ClientOnDataVP9Func
¶
type ClientOnDataVP9Func func(pts
int64
, frame []
byte
)
ClientOnDataVP9Func is the prototype of the function passed to OnDataVP9().
type
ClientOnDecodeErrorFunc
¶
type ClientOnDecodeErrorFunc func(err
error
)
ClientOnDecodeErrorFunc is the prototype of Client.OnDecodeError.
type
ClientOnDownloadPartFunc
¶
type ClientOnDownloadPartFunc func(url
string
)
ClientOnDownloadPartFunc is the prototype of Client.OnDownloadPart.
type
ClientOnDownloadPrimaryPlaylistFunc
¶
type ClientOnDownloadPrimaryPlaylistFunc func(url
string
)
ClientOnDownloadPrimaryPlaylistFunc is the prototype of Client.OnDownloadPrimaryPlaylist.
type
ClientOnDownloadSegmentFunc
¶
type ClientOnDownloadSegmentFunc func(url
string
)
ClientOnDownloadSegmentFunc is the prototype of Client.OnDownloadSegment.
type
ClientOnDownloadStreamPlaylistFunc
¶
type ClientOnDownloadStreamPlaylistFunc func(url
string
)
ClientOnDownloadStreamPlaylistFunc is the prototype of Client.OnDownloadStreamPlaylist.
type
ClientOnRequestFunc
¶
type ClientOnRequestFunc func(*
http
.
Request
)
ClientOnRequestFunc is the prototype of the function passed to OnRequest().
type
ClientOnTracksFunc
¶
type ClientOnTracksFunc func([]*
Track
)
error
ClientOnTracksFunc is the prototype of the function passed to OnTracks().
type
Muxer
¶
type Muxer struct {
//
// parameters (all optional except Tracks).
//
// tracks.
Tracks []*
Track
// Variant to use.
// It defaults to MuxerVariantLowLatency
Variant
MuxerVariant
// Number of HLS segments to keep on the server.
// Segments allow to seek through the stream.
// Their number doesn't influence latency.
// It defaults to 7.
SegmentCount
int
// Minimum duration of each segment.
// This is adjusted in order to include at least one IDR frame in each segment.
// A player usually puts 3 segments in a buffer before reproducing the stream.
// It defaults to 1sec.
SegmentMinDuration
time
.
Duration
// Minimum duration of each part.
// Parts are used in Low-Latency HLS in place of segments.
// This is adjusted in order to produce segments with a similar duration.
// A player usually puts 3 parts in a buffer before reproducing the stream.
// It defaults to 200ms.
PartMinDuration
time
.
Duration
// Maximum size of each segment.
// This prevents RAM exhaustion.
// It defaults to 50MB.
SegmentMaxSize
uint64
// Directory in which to save segments.
// This decreases performance, since saving segments on disk is less performant
// than saving them on RAM, but allows to preserve RAM.
Directory
string
//
// callbacks (all optional)
//
// called when a non-fatal encode error occurs.
OnEncodeError
MuxerOnEncodeErrorFunc
// contains filtered or unexported fields
}
Muxer is a HLS muxer.
func (*Muxer)
Close
¶
func (m *
Muxer
) Close()
Close closes a Muxer.
func (*Muxer)
Handle
¶
func (m *
Muxer
) Handle(w
http
.
ResponseWriter
, r *
http
.
Request
)
Handle handles a HTTP request.
func (*Muxer)
Start
¶
func (m *
Muxer
) Start()
error
Start initializes the muxer.
func (*Muxer)
WriteAV1
¶
func (m *
Muxer
) WriteAV1(
track *
Track
,
ntp
time
.
Time
,
pts
int64
,
tu [][]
byte
,
)
error
WriteAV1 writes an AV1 temporal unit.
func (*Muxer)
WriteH264
¶
func (m *
Muxer
) WriteH264(
track *
Track
,
ntp
time
.
Time
,
pts
int64
,
au [][]
byte
,
)
error
WriteH264 writes an H264 access unit.
func (*Muxer)
WriteH265
¶
func (m *
Muxer
) WriteH265(
track *
Track
,
ntp
time
.
Time
,
pts
int64
,
au [][]
byte
,
)
error
WriteH265 writes an H265 access unit.
func (*Muxer)
WriteMPEG4Audio
¶
func (m *
Muxer
) WriteMPEG4Audio(
track *
Track
,
ntp
time
.
Time
,
pts
int64
,
aus [][]
byte
,
)
error
WriteMPEG4Audio writes MPEG-4 Audio access units.
func (*Muxer)
WriteOpus
¶
func (m *
Muxer
) WriteOpus(
track *
Track
,
ntp
time
.
Time
,
pts
int64
,
packets [][]
byte
,
)
error
WriteOpus writes Opus packets.
func (*Muxer)
WriteVP9
¶
func (m *
Muxer
) WriteVP9(
track *
Track
,
ntp
time
.
Time
,
pts
int64
,
frame []
byte
,
)
error
WriteVP9 writes a VP9 frame.
type
MuxerOnEncodeErrorFunc
¶
type MuxerOnEncodeErrorFunc func(err
error
)
MuxerOnEncodeErrorFunc is the prototype of Muxer.OnEncodeError.
type
MuxerVariant
¶
type MuxerVariant
int
MuxerVariant is a muxer variant.
const (
MuxerVariantMPEGTS
MuxerVariant
=
iota
+ 1
MuxerVariantFMP4
MuxerVariantLowLatency
)
supported variants.
type
Track
¶
type Track struct {
// Codec
Codec
codecs
.
Codec
// Clock rate
ClockRate
int
// Name
// For audio renditions only.
Name
string
// Language
// For audio renditions only.
Language
string
// whether this is the default track.
// For audio renditions only.
IsDefault
bool
}
Track is a HLS track.