# gohlslib (HLS)

> Source: https://pkg.go.dev/github.com/bluenviron/gohlslib/v2
> Fetched: 2026-02-01T11:48:29.726557+00:00
> Content-Hash: 4a8879e1124b1e59
> Type: html

---

### Overview ¶

Package gohlslib is a HLS client and muxer library for the Go programming language. 

Examples are available at <https://github.com/bluenviron/gohlslib/tree/main/examples>

### Index ¶

  * Variables
  * type Client
  *     * func (c *Client) AbsoluteTime(track *Track) (time.Time, bool)
    * func (c *Client) Close()
    * func (c *Client) OnDataAV1(track *Track, cb ClientOnDataAV1Func)
    * func (c *Client) OnDataH26x(track *Track, cb ClientOnDataH26xFunc)
    * func (c *Client) OnDataMPEG4Audio(track *Track, cb ClientOnDataMPEG4AudioFunc)
    * func (c *Client) OnDataOpus(track *Track, cb ClientOnDataOpusFunc)
    * func (c *Client) OnDataVP9(track *Track, cb ClientOnDataVP9Func)
    * func (c *Client) Start() error
    * func (c *Client) Wait() chan errordeprecated
    * func (c *Client) Wait2() error
  * type ClientOnDataAV1Func
  * type ClientOnDataH26xFunc
  * type ClientOnDataMPEG4AudioFunc
  * type ClientOnDataOpusFunc
  * type ClientOnDataVP9Func
  * type ClientOnDecodeErrorFunc
  * type ClientOnDownloadPartFunc
  * type ClientOnDownloadPrimaryPlaylistFunc
  * type ClientOnDownloadSegmentFunc
  * type ClientOnDownloadStreamPlaylistFunc
  * type ClientOnRequestFunc
  * type ClientOnTracksFunc
  * type Muxer
  *     * func (m *Muxer) Close()
    * func (m *Muxer) Handle(w http.ResponseWriter, r *http.Request)
    * func (m *Muxer) Start() error
    * func (m *Muxer) WriteAV1(track *Track, ntp time.Time, pts int64, tu [][]byte) error
    * func (m *Muxer) WriteH264(track *Track, ntp time.Time, pts int64, au [][]byte) error
    * func (m *Muxer) WriteH265(track *Track, ntp time.Time, pts int64, au [][]byte) error
    * func (m *Muxer) WriteMPEG4Audio(track *Track, ntp time.Time, pts int64, aus [][]byte) error
    * func (m *Muxer) WriteOpus(track *Track, ntp time.Time, pts int64, packets [][]byte) error
    * func (m *Muxer) WriteVP9(track *Track, ntp time.Time, pts int64, frame []byte) error
  * type MuxerOnEncodeErrorFunc
  * type MuxerVariant
  * type Track



### Constants ¶

This section is empty.

### Variables ¶

[View Source](https://github.com/bluenviron/gohlslib/blob/v2.2.5/client.go#L25)
    
    
    var ErrClientEOS = [errors](/errors).[New](/errors#New)("end of stream")

ErrClientEOS is returned by Wait() when the stream has ended. 

### Functions ¶

This section is empty.

### Types ¶

####  type [Client](https://github.com/bluenviron/gohlslib/blob/v2.2.5/client.go#L72) ¶
    
    
    type Client struct {
    	//
    	// parameters (all optional except URI)
    	//
    	// URI of the playlist.
    	URI [string](/builtin#string)
    	// Start distance from the end of the playlist,
    	// expressed as number of segments.
    	// It defaults to 3.
    	StartDistance [int](/builtin#int)
    	// Maximum distance from the end of the playlist,
    	// expressed as number of segments.
    	// It defaults to 5.
    	MaxDistance [int](/builtin#int)
    	// HTTP client.
    	// It defaults to http.DefaultClient.
    	HTTPClient *[http](/net/http).[Client](/net/http#Client)
    
    	//
    	// callbacks (all optional)
    	//
    	// called when sending a request to the server.
    	OnRequest ClientOnRequestFunc
    	// called when tracks are available.
    	OnTracks ClientOnTracksFunc
    	// called before downloading a primary playlist.
    	OnDownloadPrimaryPlaylist ClientOnDownloadPrimaryPlaylistFunc
    	// called before downloading a stream playlist.
    	OnDownloadStreamPlaylist ClientOnDownloadStreamPlaylistFunc
    	// called before downloading a segment.
    	OnDownloadSegment ClientOnDownloadSegmentFunc
    	// called before downloading a part.
    	OnDownloadPart ClientOnDownloadPartFunc
    	// called when a non-fatal decode error occurs.
    	OnDecodeError ClientOnDecodeErrorFunc
    	// contains filtered or unexported fields
    }

Client is a HLS client. 

####  func (*Client) [AbsoluteTime](https://github.com/bluenviron/gohlslib/blob/v2.2.5/client.go#L249) ¶
    
    
    func (c *Client) AbsoluteTime(track *Track) ([time](/time).[Time](/time#Time), [bool](/builtin#bool))

AbsoluteTime returns the absolute timestamp of the last sample. 

####  func (*Client) [Close](https://github.com/bluenviron/gohlslib/blob/v2.2.5/client.go#L187) ¶
    
    
    func (c *Client) Close()

Close closes all the Client resources and waits for them to exit. 

####  func (*Client) [OnDataAV1](https://github.com/bluenviron/gohlslib/blob/v2.2.5/client.go#L212) ¶
    
    
    func (c *Client) OnDataAV1(track *Track, cb ClientOnDataAV1Func)

OnDataAV1 sets a callback that is called when data from an AV1 track is received. 

####  func (*Client) [OnDataH26x](https://github.com/bluenviron/gohlslib/blob/v2.2.5/client.go#L226) ¶
    
    
    func (c *Client) OnDataH26x(track *Track, cb ClientOnDataH26xFunc)

OnDataH26x sets a callback that is called when data from an H26x track is received. 

####  func (*Client) [OnDataMPEG4Audio](https://github.com/bluenviron/gohlslib/blob/v2.2.5/client.go#L233) ¶
    
    
    func (c *Client) OnDataMPEG4Audio(track *Track, cb ClientOnDataMPEG4AudioFunc)

OnDataMPEG4Audio sets a callback that is called when data from a MPEG-4 Audio track is received. 

####  func (*Client) [OnDataOpus](https://github.com/bluenviron/gohlslib/blob/v2.2.5/client.go#L240) ¶
    
    
    func (c *Client) OnDataOpus(track *Track, cb ClientOnDataOpusFunc)

OnDataOpus sets a callback that is called when data from an Opus track is received. 

####  func (*Client) [OnDataVP9](https://github.com/bluenviron/gohlslib/blob/v2.2.5/client.go#L219) ¶
    
    
    func (c *Client) OnDataVP9(track *Track, cb ClientOnDataVP9Func)

OnDataVP9 sets a callback that is called when data from a VP9 track is received. 

####  func (*Client) [Start](https://github.com/bluenviron/gohlslib/blob/v2.2.5/client.go#L126) ¶
    
    
    func (c *Client) Start() [error](/builtin#error)

Start starts the client. 

####  func (*Client) [Wait](https://github.com/bluenviron/gohlslib/blob/v2.2.5/client.go#L195) deprecated
    
    
    func (c *Client) Wait() chan [error](/builtin#error)

Wait waits for any error of the Client. 

Deprecated: replaced by Wait2. 

####  func (*Client) [Wait2](https://github.com/bluenviron/gohlslib/blob/v2.2.5/client.go#L206) ¶ added in v2.2.0
    
    
    func (c *Client) Wait2() [error](/builtin#error)

Wait2 waits until all client resources are closed. This can happen when a fatal error occurs or when Close() is called. 

####  type [ClientOnDataAV1Func](https://github.com/bluenviron/gohlslib/blob/v2.2.5/client.go#L49) ¶
    
    
    type ClientOnDataAV1Func func(pts [int64](/builtin#int64), tu [][][byte](/builtin#byte))

ClientOnDataAV1Func is the prototype of the function passed to OnDataAV1(). 

####  type [ClientOnDataH26xFunc](https://github.com/bluenviron/gohlslib/blob/v2.2.5/client.go#L55) ¶
    
    
    type ClientOnDataH26xFunc func(pts [int64](/builtin#int64), dts [int64](/builtin#int64), au [][][byte](/builtin#byte))

ClientOnDataH26xFunc is the prototype of the function passed to OnDataH26x(). 

####  type [ClientOnDataMPEG4AudioFunc](https://github.com/bluenviron/gohlslib/blob/v2.2.5/client.go#L58) ¶
    
    
    type ClientOnDataMPEG4AudioFunc func(pts [int64](/builtin#int64), aus [][][byte](/builtin#byte))

ClientOnDataMPEG4AudioFunc is the prototype of the function passed to OnDataMPEG4Audio(). 

####  type [ClientOnDataOpusFunc](https://github.com/bluenviron/gohlslib/blob/v2.2.5/client.go#L61) ¶
    
    
    type ClientOnDataOpusFunc func(pts [int64](/builtin#int64), packets [][][byte](/builtin#byte))

ClientOnDataOpusFunc is the prototype of the function passed to OnDataOpus(). 

####  type [ClientOnDataVP9Func](https://github.com/bluenviron/gohlslib/blob/v2.2.5/client.go#L52) ¶
    
    
    type ClientOnDataVP9Func func(pts [int64](/builtin#int64), frame [][byte](/builtin#byte))

ClientOnDataVP9Func is the prototype of the function passed to OnDataVP9(). 

####  type [ClientOnDecodeErrorFunc](https://github.com/bluenviron/gohlslib/blob/v2.2.5/client.go#L40) ¶
    
    
    type ClientOnDecodeErrorFunc func(err [error](/builtin#error))

ClientOnDecodeErrorFunc is the prototype of Client.OnDecodeError. 

####  type [ClientOnDownloadPartFunc](https://github.com/bluenviron/gohlslib/blob/v2.2.5/client.go#L37) ¶
    
    
    type ClientOnDownloadPartFunc func(url [string](/builtin#string))

ClientOnDownloadPartFunc is the prototype of Client.OnDownloadPart. 

####  type [ClientOnDownloadPrimaryPlaylistFunc](https://github.com/bluenviron/gohlslib/blob/v2.2.5/client.go#L28) ¶
    
    
    type ClientOnDownloadPrimaryPlaylistFunc func(url [string](/builtin#string))

ClientOnDownloadPrimaryPlaylistFunc is the prototype of Client.OnDownloadPrimaryPlaylist. 

####  type [ClientOnDownloadSegmentFunc](https://github.com/bluenviron/gohlslib/blob/v2.2.5/client.go#L34) ¶
    
    
    type ClientOnDownloadSegmentFunc func(url [string](/builtin#string))

ClientOnDownloadSegmentFunc is the prototype of Client.OnDownloadSegment. 

####  type [ClientOnDownloadStreamPlaylistFunc](https://github.com/bluenviron/gohlslib/blob/v2.2.5/client.go#L31) ¶
    
    
    type ClientOnDownloadStreamPlaylistFunc func(url [string](/builtin#string))

ClientOnDownloadStreamPlaylistFunc is the prototype of Client.OnDownloadStreamPlaylist. 

####  type [ClientOnRequestFunc](https://github.com/bluenviron/gohlslib/blob/v2.2.5/client.go#L43) ¶
    
    
    type ClientOnRequestFunc func(*[http](/net/http).[Request](/net/http#Request))

ClientOnRequestFunc is the prototype of the function passed to OnRequest(). 

####  type [ClientOnTracksFunc](https://github.com/bluenviron/gohlslib/blob/v2.2.5/client.go#L46) ¶
    
    
    type ClientOnTracksFunc func([]*Track) [error](/builtin#error)

ClientOnTracksFunc is the prototype of the function passed to OnTracks(). 

####  type [Muxer](https://github.com/bluenviron/gohlslib/blob/v2.2.5/muxer.go#L151) ¶
    
    
    type Muxer struct {
    	//
    	// parameters (all optional except Tracks).
    	//
    	// tracks.
    	Tracks []*Track
    	// Variant to use.
    	// It defaults to MuxerVariantLowLatency
    	Variant MuxerVariant
    	// Number of HLS segments to keep on the server.
    	// Segments allow to seek through the stream.
    	// Their number doesn't influence latency.
    	// It defaults to 7.
    	SegmentCount [int](/builtin#int)
    	// Minimum duration of each segment.
    	// This is adjusted in order to include at least one IDR frame in each segment.
    	// A player usually puts 3 segments in a buffer before reproducing the stream.
    	// It defaults to 1sec.
    	SegmentMinDuration [time](/time).[Duration](/time#Duration)
    	// Minimum duration of each part.
    	// Parts are used in Low-Latency HLS in place of segments.
    	// This is adjusted in order to produce segments with a similar duration.
    	// A player usually puts 3 parts in a buffer before reproducing the stream.
    	// It defaults to 200ms.
    	PartMinDuration [time](/time).[Duration](/time#Duration)
    	// Maximum size of each segment.
    	// This prevents RAM exhaustion.
    	// It defaults to 50MB.
    	SegmentMaxSize [uint64](/builtin#uint64)
    	// Directory in which to save segments.
    	// This decreases performance, since saving segments on disk is less performant
    	// than saving them on RAM, but allows to preserve RAM.
    	Directory [string](/builtin#string)
    
    	//
    	// callbacks (all optional)
    	//
    	// called when a non-fatal encode error occurs.
    	OnEncodeError MuxerOnEncodeErrorFunc
    	// contains filtered or unexported fields
    }

Muxer is a HLS muxer. 

####  func (*Muxer) [Close](https://github.com/bluenviron/gohlslib/blob/v2.2.5/muxer.go#L436) ¶
    
    
    func (m *Muxer) Close()

Close closes a Muxer. 

####  func (*Muxer) [Handle](https://github.com/bluenviron/gohlslib/blob/v2.2.5/muxer.go#L511) ¶
    
    
    func (m *Muxer) Handle(w [http](/net/http).[ResponseWriter](/net/http#ResponseWriter), r *[http](/net/http).[Request](/net/http#Request))

Handle handles a HTTP request. 

####  func (*Muxer) [Start](https://github.com/bluenviron/gohlslib/blob/v2.2.5/muxer.go#L209) ¶
    
    
    func (m *Muxer) Start() [error](/builtin#error)

Start initializes the muxer. 

####  func (*Muxer) [WriteAV1](https://github.com/bluenviron/gohlslib/blob/v2.2.5/muxer.go#L451) ¶
    
    
    func (m *Muxer) WriteAV1(
    	track *Track,
    	ntp [time](/time).[Time](/time#Time),
    	pts [int64](/builtin#int64),
    	tu [][][byte](/builtin#byte),
    ) [error](/builtin#error)

WriteAV1 writes an AV1 temporal unit. 

####  func (*Muxer) [WriteH264](https://github.com/bluenviron/gohlslib/blob/v2.2.5/muxer.go#L481) ¶
    
    
    func (m *Muxer) WriteH264(
    	track *Track,
    	ntp [time](/time).[Time](/time#Time),
    	pts [int64](/builtin#int64),
    	au [][][byte](/builtin#byte),
    ) [error](/builtin#error)

WriteH264 writes an H264 access unit. 

####  func (*Muxer) [WriteH265](https://github.com/bluenviron/gohlslib/blob/v2.2.5/muxer.go#L471) ¶
    
    
    func (m *Muxer) WriteH265(
    	track *Track,
    	ntp [time](/time).[Time](/time#Time),
    	pts [int64](/builtin#int64),
    	au [][][byte](/builtin#byte),
    ) [error](/builtin#error)

WriteH265 writes an H265 access unit. 

####  func (*Muxer) [WriteMPEG4Audio](https://github.com/bluenviron/gohlslib/blob/v2.2.5/muxer.go#L501) ¶
    
    
    func (m *Muxer) WriteMPEG4Audio(
    	track *Track,
    	ntp [time](/time).[Time](/time#Time),
    	pts [int64](/builtin#int64),
    	aus [][][byte](/builtin#byte),
    ) [error](/builtin#error)

WriteMPEG4Audio writes MPEG-4 Audio access units. 

####  func (*Muxer) [WriteOpus](https://github.com/bluenviron/gohlslib/blob/v2.2.5/muxer.go#L491) ¶
    
    
    func (m *Muxer) WriteOpus(
    	track *Track,
    	ntp [time](/time).[Time](/time#Time),
    	pts [int64](/builtin#int64),
    	packets [][][byte](/builtin#byte),
    ) [error](/builtin#error)

WriteOpus writes Opus packets. 

####  func (*Muxer) [WriteVP9](https://github.com/bluenviron/gohlslib/blob/v2.2.5/muxer.go#L461) ¶
    
    
    func (m *Muxer) WriteVP9(
    	track *Track,
    	ntp [time](/time).[Time](/time#Time),
    	pts [int64](/builtin#int64),
    	frame [][byte](/builtin#byte),
    ) [error](/builtin#error)

WriteVP9 writes a VP9 frame. 

####  type [MuxerOnEncodeErrorFunc](https://github.com/bluenviron/gohlslib/blob/v2.2.5/muxer.go#L148) ¶
    
    
    type MuxerOnEncodeErrorFunc func(err [error](/builtin#error))

MuxerOnEncodeErrorFunc is the prototype of Muxer.OnEncodeError. 

####  type [MuxerVariant](https://github.com/bluenviron/gohlslib/blob/v2.2.5/muxer_variant.go#L4) ¶
    
    
    type MuxerVariant [int](/builtin#int)

MuxerVariant is a muxer variant. 
    
    
    const (
    	MuxerVariantMPEGTS MuxerVariant = [iota](/builtin#iota) + 1
    	MuxerVariantFMP4
    	MuxerVariantLowLatency
    )

supported variants. 

####  type [Track](https://github.com/bluenviron/gohlslib/blob/v2.2.5/track.go#L8) ¶
    
    
    type Track struct {
    	// Codec
    	Codec [codecs](/github.com/bluenviron/gohlslib/v2@v2.2.5/pkg/codecs).[Codec](/github.com/bluenviron/gohlslib/v2@v2.2.5/pkg/codecs#Codec)
    
    	// Clock rate
    	ClockRate [int](/builtin#int)
    
    	// Name
    	// For audio renditions only.
    	Name [string](/builtin#string)
    
    	// Language
    	// For audio renditions only.
    	Language [string](/builtin#string)
    
    	// whether this is the default track.
    	// For audio renditions only.
    	IsDefault [bool](/builtin#bool)
    }

Track is a HLS track. 
  *[↑]: Back to Top
  *[v]: View this template
  *[t]: Discuss this template
  *[e]: Edit this template
