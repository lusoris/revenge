# bogem/id3v2

> Source: https://pkg.go.dev/github.com/bogem/id3v2/v2
> Fetched: 2026-01-31T16:01:16.110096+00:00
> Content-Hash: 68f108ceab998213
> Type: html

---

### Overview ¶

Package id3v2 is the ID3 parsing and writing library for Go. 

Example ¶
    
    
    // Open file and parse tag in it.
    tag, err := id3v2.Open("file.mp3", id3v2.Options{Parse: true})
    if err != nil {
    	log.Fatal("Error while opening mp3 file: ", err)
    }
    defer tag.Close()
    
    // Read frames.
    fmt.Println(tag.Artist())
    fmt.Println(tag.Title())
    
    // Set simple text frames.
    tag.SetArtist("Artist")
    tag.SetTitle("Title")
    
    // Set comment frame.
    comment := id3v2.CommentFrame{
    	Encoding:    id3v2.EncodingUTF8,
    	Language:    "eng",
    	Description: "My opinion",
    	Text:        "Very good song",
    }
    tag.AddCommentFrame(comment)
    
    // Write tag to file.
    if err = tag.Save(); err != nil {
    	log.Fatal("Error while saving a tag: ", err)
    }
    

Example (Concurrent) ¶
    
    
    tagPool := sync.Pool{New: func() interface{} { return id3v2.NewEmptyTag() }}
    
    var wg sync.WaitGroup
    wg.Add(100)
    for i := 0; i < 100; i++ {
    	go func() {
    		defer wg.Done()
    
    		tag := tagPool.Get().(*id3v2.Tag)
    		defer tagPool.Put(tag)
    
    		file, err := os.Open("file.mp3")
    		if err != nil {
    			log.Fatal("Error while opening file:", err)
    		}
    		defer file.Close()
    
    		if err := tag.Reset(file, id3v2.Options{Parse: true}); err != nil {
    			log.Fatal("Error while reseting tag to file:", err)
    		}
    
    		fmt.Println(tag.Artist() + " - " + tag.Title())
    	}()
    }
    wg.Wait()
    

### Index ¶

  * Constants
  * Variables
  * type ChapterFrame
  *     * func (cf ChapterFrame) Size() int
    * func (cf ChapterFrame) UniqueIdentifier() string
    * func (cf ChapterFrame) WriteTo(w io.Writer) (n int64, err error)
  * type CommentFrame
  *     * func (cf CommentFrame) Size() int
    * func (cf CommentFrame) UniqueIdentifier() string
    * func (cf CommentFrame) WriteTo(w io.Writer) (n int64, err error)
  * type Encoding
  *     * func (e Encoding) Equals(other Encoding) bool
    * func (e Encoding) String() string
  * type Framer
  * type Options
  * type PictureFrame
  *     * func (pf PictureFrame) Size() int
    * func (pf PictureFrame) UniqueIdentifier() string
    * func (pf PictureFrame) WriteTo(w io.Writer) (n int64, err error)
  * type PopularimeterFrame
  *     * func (pf PopularimeterFrame) Size() int
    * func (pf PopularimeterFrame) UniqueIdentifier() string
    * func (pf PopularimeterFrame) WriteTo(w io.Writer) (n int64, err error)
  * type Tag
  *     * func NewEmptyTag() *Tag
    * func Open(name string, opts Options) (*Tag, error)
    * func ParseReader(rd io.Reader, opts Options) (*Tag, error)
  *     * func (tag *Tag) AddAttachedPicture(pf PictureFrame)
    * func (tag *Tag) AddChapterFrame(cf ChapterFrame)
    * func (tag *Tag) AddCommentFrame(cf CommentFrame)
    * func (tag *Tag) AddFrame(id string, f Framer)
    * func (tag *Tag) AddTextFrame(id string, encoding Encoding, text string)
    * func (tag *Tag) AddUFIDFrame(ufid UFIDFrame)
    * func (tag *Tag) AddUnsynchronisedLyricsFrame(uslf UnsynchronisedLyricsFrame)
    * func (tag *Tag) AddUserDefinedTextFrame(udtf UserDefinedTextFrame)
    * func (tag *Tag) Album() string
    * func (tag *Tag) AllFrames() map[string][]Framer
    * func (tag *Tag) Artist() string
    * func (tag *Tag) Close() error
    * func (tag *Tag) CommonID(description string) string
    * func (tag *Tag) Count() int
    * func (tag *Tag) DefaultEncoding() Encoding
    * func (tag *Tag) DeleteAllFrames()
    * func (tag *Tag) DeleteFrames(id string)
    * func (tag *Tag) Genre() string
    * func (tag *Tag) GetFrames(id string) []Framer
    * func (tag *Tag) GetLastFrame(id string) Framer
    * func (tag *Tag) GetTextFrame(id string) TextFrame
    * func (tag *Tag) HasFrames() bool
    * func (tag *Tag) Reset(rd io.Reader, opts Options) error
    * func (tag *Tag) Save() error
    * func (tag *Tag) SetAlbum(album string)
    * func (tag *Tag) SetArtist(artist string)
    * func (tag *Tag) SetDefaultEncoding(encoding Encoding)
    * func (tag *Tag) SetGenre(genre string)
    * func (tag *Tag) SetTitle(title string)
    * func (tag *Tag) SetVersion(version byte)
    * func (tag *Tag) SetYear(year string)
    * func (tag *Tag) Size() int
    * func (tag *Tag) Title() string
    * func (tag *Tag) Version() byte
    * func (tag *Tag) WriteTo(w io.Writer) (n int64, err error)
    * func (tag *Tag) Year() string
  * type TextFrame
  *     * func (tf TextFrame) Size() int
    * func (tf TextFrame) UniqueIdentifier() string
    * func (tf TextFrame) WriteTo(w io.Writer) (int64, error)
  * type UFIDFrame
  *     * func (ufid UFIDFrame) Size() int
    * func (ufid UFIDFrame) UniqueIdentifier() string
    * func (ufid UFIDFrame) WriteTo(w io.Writer) (n int64, err error)
  * type UnknownFrame
  *     * func (uf UnknownFrame) Size() int
    * func (uf UnknownFrame) UniqueIdentifier() string
    * func (uf UnknownFrame) WriteTo(w io.Writer) (n int64, err error)
  * type UnsynchronisedLyricsFrame
  *     * func (uslf UnsynchronisedLyricsFrame) Size() int
    * func (uslf UnsynchronisedLyricsFrame) UniqueIdentifier() string
    * func (uslf UnsynchronisedLyricsFrame) WriteTo(w io.Writer) (n int64, err error)
  * type UserDefinedTextFrame
  *     * func (udtf UserDefinedTextFrame) Size() int
    * func (udtf UserDefinedTextFrame) UniqueIdentifier() string
    * func (udtf UserDefinedTextFrame) WriteTo(w io.Writer) (n int64, err error)



### Examples ¶

  * Package
  * Package (Concurrent)
  * CommentFrame (Add)
  * CommentFrame (Get)
  * PictureFrame (Add)
  * PictureFrame (Get)
  * PopularimeterFrame (Add)
  * PopularimeterFrame (Get)
  * Tag.AddAttachedPicture
  * Tag.AddCommentFrame
  * Tag.AddUnsynchronisedLyricsFrame
  * Tag.GetFrames
  * Tag.GetLastFrame
  * TextFrame (Add)
  * TextFrame (Get)
  * UnsynchronisedLyricsFrame (Add)
  * UnsynchronisedLyricsFrame (Get)



### Constants ¶

[View Source](https://github.com/bogem/id3v2/blob/v2.1.4/v2/id3v2.go#L14)
    
    
    const (
    	PTOther = [iota](/builtin#iota)
    	PTFileIcon
    	PTOtherFileIcon
    	PTFrontCover
    	PTBackCover
    	PTLeafletPage
    	PTMedia
    	PTLeadArtistSoloist
    	PTArtistPerformer
    	PTConductor
    	PTBandOrchestra
    	PTComposer
    	PTLyricistTextWriter
    	PTRecordingLocation
    	PTDuringRecording
    	PTDuringPerformance
    	PTMovieScreenCapture
    	PTBrightColouredFish
    	PTIllustration
    	PTBandArtistLogotype
    	PTPublisherStudioLogotype
    )

Available picture types for picture frame. 

[View Source](https://github.com/bogem/id3v2/blob/v2.1.4/v2/chapter_frame.go#L9)
    
    
    const (
    	IgnoredOffset = 0xFFFFFFFF
    )

### Variables ¶

[View Source](https://github.com/bogem/id3v2/blob/v2.1.4/v2/common_ids.go#L10)
    
    
    var (
    	V23CommonIDs = map[[string](/builtin#string)][string](/builtin#string){
    		"Attached picture":                   "APIC",
    		"Chapters":                           "CHAP",
    		"Comments":                           "COMM",
    		"Album/Movie/Show title":             "TALB",
    		"BPM":                                "TBPM",
    		"Composer":                           "TCOM",
    		"Content type":                       "TCON",
    		"Copyright message":                  "TCOP",
    		"Date":                               "TDAT",
    		"Playlist delay":                     "TDLY",
    		"Encoded by":                         "TENC",
    		"Lyricist/Text writer":               "TEXT",
    		"File type":                          "TFLT",
    		"Time":                               "TIME",
    		"Content group description":          "TIT1",
    		"Title/Songname/Content description": "TIT2",
    		"Subtitle/Description refinement":    "TIT3",
    		"Initial key":                        "TKEY",
    		"Language":                           "TLAN",
    		"Length":                             "TLEN",
    		"Media type":                         "TMED",
    		"Original album/movie/show title":    "TOAL",
    		"Original filename":                  "TOFN",
    		"Original lyricist/text writer":      "TOLY",
    		"Original artist/performer":          "TOPE",
    		"Original release year":              "TORY",
    		"Popularimeter":                      "POPM",
    		"File owner/licensee":                "TOWN",
    		"Lead artist/Lead performer/Soloist/Performing group": "TPE1",
    		"Band/Orchestra/Accompaniment":                        "TPE2",
    		"Conductor/performer refinement":                      "TPE3",
    		"Interpreted, remixed, or otherwise modified by":      "TPE4",
    		"Part of a set":                "TPOS",
    		"Publisher":                    "TPUB",
    		"Track number/Position in set": "TRCK",
    		"Recording dates":              "TRDA",
    		"Internet radio station name":  "TRSN",
    		"Internet radio station owner": "TRSO",
    		"Size":                         "TSIZ",
    		"ISRC":                         "TSRC",
    		"Software/Hardware and settings used for encoding": "TSSE",
    		"Year":                                     "TYER",
    		"User defined text information frame":      "TXXX",
    		"Unique file identifier":                   "UFID",
    		"Unsynchronised lyrics/text transcription": "USLT",
    
    		"Artist": "TPE1",
    		"Title":  "TIT2",
    		"Genre":  "TCON",
    	}
    
    	V24CommonIDs = map[[string](/builtin#string)][string](/builtin#string){
    		"Attached picture":                   "APIC",
    		"Chapters":                           "CHAP",
    		"Comments":                           "COMM",
    		"Album/Movie/Show title":             "TALB",
    		"BPM":                                "TBPM",
    		"Composer":                           "TCOM",
    		"Content type":                       "TCON",
    		"Copyright message":                  "TCOP",
    		"Encoding time":                      "TDEN",
    		"Playlist delay":                     "TDLY",
    		"Original release time":              "TDOR",
    		"Recording time":                     "TDRC",
    		"Release time":                       "TDRL",
    		"Tagging time":                       "TDTG",
    		"Encoded by":                         "TENC",
    		"Lyricist/Text writer":               "TEXT",
    		"File type":                          "TFLT",
    		"Involved people list":               "TIPL",
    		"Content group description":          "TIT1",
    		"Title/Songname/Content description": "TIT2",
    		"Subtitle/Description refinement":    "TIT3",
    		"Initial key":                        "TKEY",
    		"Language":                           "TLAN",
    		"Length":                             "TLEN",
    		"Musician credits list":              "TMCL",
    		"Media type":                         "TMED",
    		"Mood":                               "TMOO",
    		"Original album/movie/show title":    "TOAL",
    		"Original filename":                  "TOFN",
    		"Original lyricist/text writer":      "TOLY",
    		"Original artist/performer":          "TOPE",
    		"Popularimeter":                      "POPM",
    		"File owner/licensee":                "TOWN",
    		"Lead artist/Lead performer/Soloist/Performing group": "TPE1",
    		"Band/Orchestra/Accompaniment":                        "TPE2",
    		"Conductor/performer refinement":                      "TPE3",
    		"Interpreted, remixed, or otherwise modified by":      "TPE4",
    		"Part of a set":                "TPOS",
    		"Produced notice":              "TPRO",
    		"Publisher":                    "TPUB",
    		"Track number/Position in set": "TRCK",
    		"Internet radio station name":  "TRSN",
    		"Internet radio station owner": "TRSO",
    		"Album sort order":             "TSOA",
    		"Performer sort order":         "TSOP",
    		"Title sort order":             "TSOT",
    		"ISRC":                         "TSRC",
    		"Software/Hardware and settings used for encoding": "TSSE",
    		"Set subtitle":                             "TSST",
    		"User defined text information frame":      "TXXX",
    		"Unique file identifier":                   "UFID",
    		"Unsynchronised lyrics/text transcription": "USLT",
    
    		"Date":                  "TDRC",
    		"Time":                  "TDRC",
    		"Original release year": "TDOR",
    		"Recording dates":       "TDRC",
    		"Size":                  "",
    		"Year":                  "TDRC",
    
    		"Artist": "TPE1",
    		"Title":  "TIT2",
    		"Genre":  "TCON",
    	}
    )

Common IDs for ID3v2.3 and ID3v2.4. 

[View Source](https://github.com/bogem/id3v2/blob/v2.1.4/v2/encoding.go#L28)
    
    
    var (
    	// EncodingISO is ISO-8859-1 encoding.
    	EncodingISO = Encoding{
    		Name:             "ISO-8859-1",
    		Key:              0,
    		TerminationBytes: [][byte](/builtin#byte){0},
    	}
    
    	// EncodingUTF16 is UTF-16 encoded Unicode with BOM.
    	EncodingUTF16 = Encoding{
    		Name:             "UTF-16 encoded Unicode with BOM",
    		Key:              1,
    		TerminationBytes: [][byte](/builtin#byte){0, 0},
    	}
    
    	// EncodingUTF16BE is UTF-16BE encoded Unicode without BOM.
    	EncodingUTF16BE = Encoding{
    		Name:             "UTF-16BE encoded Unicode without BOM",
    		Key:              2,
    		TerminationBytes: [][byte](/builtin#byte){0, 0},
    	}
    
    	// EncodingUTF8 is UTF-8 encoded Unicode.
    	EncodingUTF8 = Encoding{
    		Name:             "UTF-8 encoded Unicode",
    		Key:              3,
    		TerminationBytes: [][byte](/builtin#byte){0},
    	}
    )

Available encodings. 

[View Source](https://github.com/bogem/id3v2/blob/v2.1.4/v2/parse.go#L19)
    
    
    var ErrBodyOverflow = [errors](/errors).[New](/errors#New)("frame went over tag area")

ErrBodyOverflow is returned when a frame has greater size than the remaining tag size 

[View Source](https://github.com/bogem/id3v2/blob/v2.1.4/v2/framer.go#L12)
    
    
    var ErrInvalidLanguageLength = [errors](/errors).[New](/errors#New)("language code must consist of three letters according to ISO 639-2")

[View Source](https://github.com/bogem/id3v2/blob/v2.1.4/v2/size.go#L22)
    
    
    var ErrInvalidSizeFormat = [errors](/errors).[New](/errors#New)("invalid format of tag's/frame's size")

[View Source](https://github.com/bogem/id3v2/blob/v2.1.4/v2/tag.go#L13)
    
    
    var ErrNoFile = [errors](/errors).[New](/errors#New)("tag was not initialized with file")

[View Source](https://github.com/bogem/id3v2/blob/v2.1.4/v2/size.go#L23)
    
    
    var ErrSizeOverflow = [errors](/errors).[New](/errors#New)("size of tag/frame is greater than allowed in id3 tag")

[View Source](https://github.com/bogem/id3v2/blob/v2.1.4/v2/header.go#L20)
    
    
    var ErrSmallHeaderSize = [errors](/errors).[New](/errors#New)("size of tag header is less than expected")

[View Source](https://github.com/bogem/id3v2/blob/v2.1.4/v2/parse.go#L15)
    
    
    var ErrUnsupportedVersion = [errors](/errors).[New](/errors#New)("unsupported version of ID3 tag")

### Functions ¶

This section is empty.

### Types ¶

####  type [ChapterFrame](https://github.com/bogem/id3v2/blob/v2.1.4/v2/chapter_frame.go#L20) ¶
    
    
    type ChapterFrame struct {
    	ElementID   [string](/builtin#string)
    	StartTime   [time](/time).[Duration](/time#Duration)
    	EndTime     [time](/time).[Duration](/time#Duration)
    	StartOffset [uint32](/builtin#uint32)
    	EndOffset   [uint32](/builtin#uint32)
    	Title       *TextFrame
    	Description *TextFrame
    }

ChapterFrame is used to work with CHAP frames according to spec from <http://id3.org/id3v2-chapters-1.0> This implementation only supports single TIT2 subframe (Title field). All other subframes are ignored. If StartOffset or EndOffset == id3v2.IgnoredOffset, then it should be ignored and StartTime or EndTime should be utilized 

####  func (ChapterFrame) [Size](https://github.com/bogem/id3v2/blob/v2.1.4/v2/chapter_frame.go#L30) ¶
    
    
    func (cf ChapterFrame) Size() [int](/builtin#int)

####  func (ChapterFrame) [UniqueIdentifier](https://github.com/bogem/id3v2/blob/v2.1.4/v2/chapter_frame.go#L47) ¶
    
    
    func (cf ChapterFrame) UniqueIdentifier() [string](/builtin#string)

####  func (ChapterFrame) [WriteTo](https://github.com/bogem/id3v2/blob/v2.1.4/v2/chapter_frame.go#L51) ¶
    
    
    func (cf ChapterFrame) WriteTo(w [io](/io).[Writer](/io#Writer)) (n [int64](/builtin#int64), err [error](/builtin#error))

####  type [CommentFrame](https://github.com/bogem/id3v2/blob/v2.1.4/v2/comment_frame.go#L15) ¶
    
    
    type CommentFrame struct {
    	Encoding    Encoding
    	Language    [string](/builtin#string)
    	Description [string](/builtin#string)
    	Text        [string](/builtin#string)
    }

CommentFrame is used to work with COMM frames. The information about how to add comment frame to tag you can see in the docs to tag.AddCommentFrame function. 

You must choose a three-letter language code from ISO 639-2 code list: <https://www.loc.gov/standards/iso639-2/php/code_list.php>

Example (Add) ¶
    
    
    tag := id3v2.NewEmptyTag()
    comment := id3v2.CommentFrame{
    	Encoding:    id3v2.EncodingUTF8,
    	Language:    "eng",
    	Description: "My opinion",
    	Text:        "Very good song",
    }
    tag.AddCommentFrame(comment)
    

Example (Get) ¶
    
    
    tag, err := id3v2.Open("file.mp3", id3v2.Options{Parse: true})
    if tag == nil || err != nil {
    	log.Fatal("Error while opening mp3 file: ", err)
    }
    
    comments := tag.GetFrames(tag.CommonID("Comments"))
    for _, f := range comments {
    	comment, ok := f.(id3v2.CommentFrame)
    	if !ok {
    		log.Fatal("Couldn't assert comment frame")
    	}
    
    	// Do something with comment frame.
    	// For example, print the text:
    	fmt.Println(comment.Text)
    }
    

####  func (CommentFrame) [Size](https://github.com/bogem/id3v2/blob/v2.1.4/v2/comment_frame.go#L22) ¶
    
    
    func (cf CommentFrame) Size() [int](/builtin#int)

####  func (CommentFrame) [UniqueIdentifier](https://github.com/bogem/id3v2/blob/v2.1.4/v2/comment_frame.go#L27) ¶
    
    
    func (cf CommentFrame) UniqueIdentifier() [string](/builtin#string)

####  func (CommentFrame) [WriteTo](https://github.com/bogem/id3v2/blob/v2.1.4/v2/comment_frame.go#L31) ¶
    
    
    func (cf CommentFrame) WriteTo(w [io](/io).[Writer](/io#Writer)) (n [int64](/builtin#int64), err [error](/builtin#error))

####  type [Encoding](https://github.com/bogem/id3v2/blob/v2.1.4/v2/encoding.go#L13) ¶
    
    
    type Encoding struct {
    	Name             [string](/builtin#string)
    	Key              [byte](/builtin#byte)
    	TerminationBytes [][byte](/builtin#byte)
    }

Encoding is a struct for encodings. 

####  func (Encoding) [Equals](https://github.com/bogem/id3v2/blob/v2.1.4/v2/encoding.go#L19) ¶
    
    
    func (e Encoding) Equals(other Encoding) [bool](/builtin#bool)

####  func (Encoding) [String](https://github.com/bogem/id3v2/blob/v2.1.4/v2/encoding.go#L23) ¶
    
    
    func (e Encoding) String() [string](/builtin#string)

####  type [Framer](https://github.com/bogem/id3v2/blob/v2.1.4/v2/framer.go#L16) ¶
    
    
    type Framer interface {
    	// Size returns the size of frame body.
    	Size() [int](/builtin#int)
    
    	// UniqueIdentifier returns the string that makes this frame unique from others.
    	// For example, some frames with same id can be added in tag, but they should be differ in other properties.
    	// E.g. It would be "Description" for TXXX and APIC.
    	//
    	// Frames that can be added only once with same id (e.g. all text frames) should return just "".
    	UniqueIdentifier() [string](/builtin#string)
    
    	// WriteTo writes body slice into w.
    	WriteTo(w [io](/io).[Writer](/io#Writer)) (n [int64](/builtin#int64), err [error](/builtin#error))
    }

Framer provides a generic interface for frames. You can create your own frames. They must implement only this interface. 

####  type [Options](https://github.com/bogem/id3v2/blob/v2.1.4/v2/options.go#L8) ¶
    
    
    type Options struct {
    	// Parse defines, if tag will be parsed.
    	Parse [bool](/builtin#bool)
    
    	// ParseFrames defines, that frames do you only want to parse. For example,
    	// `ParseFrames: []string{"Artist", "Title"}` will only parse artist
    	// and title frames. You can specify IDs ("TPE1", "TIT2") as well as
    	// descriptions ("Artist", "Title"). If ParseFrame is blank or nil,
    	// id3v2 will parse all frames in tag. It works only if Parse is true.
    	//
    	// It's very useful for performance, so for example
    	// if you want to get only some text frames,
    	// id3v2 will not parse huge picture or unknown frames.
    	ParseFrames [][string](/builtin#string)
    }

Options influence on processing the tag. 

####  type [PictureFrame](https://github.com/bogem/id3v2/blob/v2.1.4/v2/picture_frame.go#L17) ¶
    
    
    type PictureFrame struct {
    	Encoding    Encoding
    	MimeType    [string](/builtin#string)
    	PictureType [byte](/builtin#byte)
    	Description [string](/builtin#string)
    	Picture     [][byte](/builtin#byte)
    }

PictureFrame structure is used for picture frames (APIC). The information about how to add picture frame to tag you can see in the docs to tag.AddAttachedPicture function. 

Available picture types you can see in constants. 

Example (Add) ¶
    
    
    tag := id3v2.NewEmptyTag()
    artwork, err := ioutil.ReadFile("artwork.jpg")
    if err != nil {
    	log.Fatal("Error while reading artwork file", err)
    }
    
    pic := id3v2.PictureFrame{
    	Encoding:    id3v2.EncodingUTF8,
    	MimeType:    "image/jpeg",
    	PictureType: id3v2.PTFrontCover,
    	Description: "Front cover",
    	Picture:     artwork,
    }
    tag.AddAttachedPicture(pic)
    

Example (Get) ¶
    
    
    tag, err := id3v2.Open("file.mp3", id3v2.Options{Parse: true})
    if tag == nil || err != nil {
    	log.Fatal("Error while opening mp3 file: ", err)
    }
    
    pictures := tag.GetFrames(tag.CommonID("Attached picture"))
    for _, f := range pictures {
    	pic, ok := f.(id3v2.PictureFrame)
    	if !ok {
    		log.Fatal("Couldn't assert picture frame")
    	}
    
    	// Do something with picture frame.
    	// For example, print the description:
    	fmt.Println(pic.Description)
    }
    

####  func (PictureFrame) [Size](https://github.com/bogem/id3v2/blob/v2.1.4/v2/picture_frame.go#L29) ¶
    
    
    func (pf PictureFrame) Size() [int](/builtin#int)

####  func (PictureFrame) [UniqueIdentifier](https://github.com/bogem/id3v2/blob/v2.1.4/v2/picture_frame.go#L25) ¶
    
    
    func (pf PictureFrame) UniqueIdentifier() [string](/builtin#string)

####  func (PictureFrame) [WriteTo](https://github.com/bogem/id3v2/blob/v2.1.4/v2/picture_frame.go#L34) ¶
    
    
    func (pf PictureFrame) WriteTo(w [io](/io).[Writer](/io#Writer)) (n [int64](/builtin#int64), err [error](/builtin#error))

####  type [PopularimeterFrame](https://github.com/bogem/id3v2/blob/v2.1.4/v2/popularimeter_frame.go#L10) ¶
    
    
    type PopularimeterFrame struct {
    	// Email is the identifier for a POPM frame.
    	Email [string](/builtin#string)
    
    	// The rating is 1-255 where 1 is worst and 255 is best. 0 is unknown.
    	Rating [uint8](/builtin#uint8)
    
    	// Counter is the number of times this file has been played by this email.
    	Counter *[big](/math/big).[Int](/math/big#Int)
    }

PopularimeterFrame structure is used for Popularimeter (POPM). <https://id3.org/id3v2.3.0#Popularimeter>

Example (Add) ¶
    
    
    tag := id3v2.NewEmptyTag()
    
    popmFrame := id3v2.PopularimeterFrame{
    	Email:   "foo@bar.com",
    	Rating:  128,
    	Counter: big.NewInt(10000000000000000),
    }
    tag.AddFrame(tag.CommonID("Popularimeter"), popmFrame)
    

Example (Get) ¶
    
    
    tag, err := id3v2.Open("file.mp3", id3v2.Options{Parse: true})
    if tag == nil || err != nil {
    	log.Fatal("Error while opening mp3 file: ", err)
    }
    
    f := tag.GetLastFrame(tag.CommonID("Popularimeter"))
    popm, ok := f.(id3v2.PopularimeterFrame)
    if !ok {
    	log.Fatal("Couldn't assert POPM frame")
    }
    
    // do something with POPM Frame
    fmt.Printf("Email: %s, Rating: %d, Counter: %d", popm.Email, popm.Rating, popm.Counter)
    

####  func (PopularimeterFrame) [Size](https://github.com/bogem/id3v2/blob/v2.1.4/v2/popularimeter_frame.go#L25) ¶
    
    
    func (pf PopularimeterFrame) Size() [int](/builtin#int)

####  func (PopularimeterFrame) [UniqueIdentifier](https://github.com/bogem/id3v2/blob/v2.1.4/v2/popularimeter_frame.go#L21) ¶
    
    
    func (pf PopularimeterFrame) UniqueIdentifier() [string](/builtin#string)

####  func (PopularimeterFrame) [WriteTo](https://github.com/bogem/id3v2/blob/v2.1.4/v2/popularimeter_frame.go#L44) ¶
    
    
    func (pf PopularimeterFrame) WriteTo(w [io](/io).[Writer](/io#Writer)) (n [int64](/builtin#int64), err [error](/builtin#error))

####  type [Tag](https://github.com/bogem/id3v2/blob/v2.1.4/v2/tag.go#L16) ¶
    
    
    type Tag struct {
    	// contains filtered or unexported fields
    }

Tag stores all information about opened tag. 

####  func [NewEmptyTag](https://github.com/bogem/id3v2/blob/v2.1.4/v2/id3v2.go#L57) ¶
    
    
    func NewEmptyTag() *Tag

NewEmptyTag returns an empty ID3v2.4 tag without any frames and reader. 

####  func [Open](https://github.com/bogem/id3v2/blob/v2.1.4/v2/id3v2.go#L40) ¶
    
    
    func Open(name [string](/builtin#string), opts Options) (*Tag, [error](/builtin#error))

Open opens file with name and passes it to OpenFile. If there is no tag in file, it will create new one with version ID3v2.4. 

####  func [ParseReader](https://github.com/bogem/id3v2/blob/v2.1.4/v2/id3v2.go#L50) ¶
    
    
    func ParseReader(rd [io](/io).[Reader](/io#Reader), opts Options) (*Tag, [error](/builtin#error))

ParseReader parses rd and finds tag in it considering opts. If there is no tag in rd, it will create new one with version ID3v2.4. 

####  func (*Tag) [AddAttachedPicture](https://github.com/bogem/id3v2/blob/v2.1.4/v2/tag.go#L50) ¶
    
    
    func (tag *Tag) AddAttachedPicture(pf PictureFrame)

AddAttachedPicture adds the picture frame to tag. 

Example ¶
    
    
    tag, err := id3v2.Open("file.mp3", id3v2.Options{Parse: true})
    if tag == nil || err != nil {
    	log.Fatal("Error while opening mp3 file: ", err)
    }
    
    artwork, err := ioutil.ReadFile("artwork.jpg")
    if err != nil {
    	log.Fatal("Error while reading artwork file", err)
    }
    
    pic := id3v2.PictureFrame{
    	Encoding:    id3v2.EncodingUTF8,
    	MimeType:    "image/jpeg",
    	PictureType: id3v2.PTFrontCover,
    	Description: "Front cover",
    	Picture:     artwork,
    }
    tag.AddAttachedPicture(pic)
    

####  func (*Tag) [AddChapterFrame](https://github.com/bogem/id3v2/blob/v2.1.4/v2/tag.go#L55) ¶
    
    
    func (tag *Tag) AddChapterFrame(cf ChapterFrame)

AddChapterFrame adds the chapter frame to tag. 

####  func (*Tag) [AddCommentFrame](https://github.com/bogem/id3v2/blob/v2.1.4/v2/tag.go#L60) ¶
    
    
    func (tag *Tag) AddCommentFrame(cf CommentFrame)

AddCommentFrame adds the comment frame to tag. 

Example ¶
    
    
    tag, err := id3v2.Open("file.mp3", id3v2.Options{Parse: true})
    if tag == nil || err != nil {
    	log.Fatal("Error while opening mp3 file: ", err)
    }
    
    comment := id3v2.CommentFrame{
    	Encoding:    id3v2.EncodingUTF8,
    	Language:    "eng",
    	Description: "My opinion",
    	Text:        "Very good song",
    }
    tag.AddCommentFrame(comment)
    

####  func (*Tag) [AddFrame](https://github.com/bogem/id3v2/blob/v2.1.4/v2/tag.go#L32) ¶
    
    
    func (tag *Tag) AddFrame(id [string](/builtin#string), f Framer)

AddFrame adds f to tag with appropriate id. If id is "" or f is nil, AddFrame will not add it to tag. 

If you want to add attached picture, comment or unsynchronised lyrics/text transcription frames, better use AddAttachedPicture, AddCommentFrame or AddUnsynchronisedLyricsFrame methods respectively. 

####  func (*Tag) [AddTextFrame](https://github.com/bogem/id3v2/blob/v2.1.4/v2/tag.go#L66) ¶
    
    
    func (tag *Tag) AddTextFrame(id [string](/builtin#string), encoding Encoding, text [string](/builtin#string))

AddTextFrame creates the text frame with provided encoding and text and adds to tag. 

####  func (*Tag) [AddUFIDFrame](https://github.com/bogem/id3v2/blob/v2.1.4/v2/tag.go#L82) ¶
    
    
    func (tag *Tag) AddUFIDFrame(ufid UFIDFrame)

AddUFIDFrame adds the unique file identifier frame (UFID) to tag. 

####  func (*Tag) [AddUnsynchronisedLyricsFrame](https://github.com/bogem/id3v2/blob/v2.1.4/v2/tag.go#L72) ¶
    
    
    func (tag *Tag) AddUnsynchronisedLyricsFrame(uslf UnsynchronisedLyricsFrame)

AddUnsynchronisedLyricsFrame adds the unsynchronised lyrics/text frame to tag. 

Example ¶
    
    
    tag, err := id3v2.Open("file.mp3", id3v2.Options{Parse: true})
    if tag == nil || err != nil {
    	log.Fatal("Error while opening mp3 file: ", err)
    }
    
    uslt := id3v2.UnsynchronisedLyricsFrame{
    	Encoding:          id3v2.EncodingUTF8,
    	Language:          "ger",
    	ContentDescriptor: "Deutsche Nationalhymne",
    	Lyrics:            "Einigkeit und Recht und Freiheit...",
    }
    tag.AddUnsynchronisedLyricsFrame(uslt)
    

####  func (*Tag) [AddUserDefinedTextFrame](https://github.com/bogem/id3v2/blob/v2.1.4/v2/tag.go#L77) ¶
    
    
    func (tag *Tag) AddUserDefinedTextFrame(udtf UserDefinedTextFrame)

AddUserDefinedTextFrame adds the custom frame (TXXX) to tag. 

####  func (*Tag) [Album](https://github.com/bogem/id3v2/blob/v2.1.4/v2/tag.go#L241) ¶
    
    
    func (tag *Tag) Album() [string](/builtin#string)

####  func (*Tag) [AllFrames](https://github.com/bogem/id3v2/blob/v2.1.4/v2/tag.go#L109) ¶
    
    
    func (tag *Tag) AllFrames() map[[string](/builtin#string)][]Framer

AllFrames returns map, that contains all frames in tag, that could be parsed. The key of this map is an ID of frame and value is an array of frames. 

####  func (*Tag) [Artist](https://github.com/bogem/id3v2/blob/v2.1.4/v2/tag.go#L233) ¶
    
    
    func (tag *Tag) Artist() [string](/builtin#string)

####  func (*Tag) [Close](https://github.com/bogem/id3v2/blob/v2.1.4/v2/tag.go#L446) ¶
    
    
    func (tag *Tag) Close() [error](/builtin#error)

Close closes tag's file, if tag was opened with a file. If tag was initiliazed not with file, it returns ErrNoFile. 

####  func (*Tag) [CommonID](https://github.com/bogem/id3v2/blob/v2.1.4/v2/tag.go#L94) ¶
    
    
    func (tag *Tag) CommonID(description [string](/builtin#string)) [string](/builtin#string)

CommonID returns frame ID from given description. For example, CommonID("Language") will return "TLAN". If it can't find the ID with given description, it returns the description. 

All descriptions you can find in file common_ids.go or in id3 documentation. v2.3: <http://id3.org/id3v2.3.0#Declared_ID3v2_frames> v2.4: <http://id3.org/id3v2.4.0-frames>

####  func (*Tag) [Count](https://github.com/bogem/id3v2/blob/v2.1.4/v2/tag.go#L211) ¶
    
    
    func (tag *Tag) Count() [int](/builtin#int)

Count returns the number of frames in tag. 

####  func (*Tag) [DefaultEncoding](https://github.com/bogem/id3v2/blob/v2.1.4/v2/tag.go#L191) ¶
    
    
    func (tag *Tag) DefaultEncoding() Encoding

DefaultEncoding returns default encoding of tag. Default encoding is used in methods (e.g. SetArtist, SetAlbum ...) for setting text frames without the explicit providing of encoding. 

####  func (*Tag) [DeleteAllFrames](https://github.com/bogem/id3v2/blob/v2.1.4/v2/tag.go#L123) ¶
    
    
    func (tag *Tag) DeleteAllFrames()

DeleteAllFrames deletes all frames in tag. 

####  func (*Tag) [DeleteFrames](https://github.com/bogem/id3v2/blob/v2.1.4/v2/tag.go#L136) ¶
    
    
    func (tag *Tag) DeleteFrames(id [string](/builtin#string))

DeleteFrames deletes frames in tag with given id. 

####  func (*Tag) [Genre](https://github.com/bogem/id3v2/blob/v2.1.4/v2/tag.go#L257) ¶
    
    
    func (tag *Tag) Genre() [string](/builtin#string)

####  func (*Tag) [GetFrames](https://github.com/bogem/id3v2/blob/v2.1.4/v2/tag.go#L152) ¶
    
    
    func (tag *Tag) GetFrames(id [string](/builtin#string)) []Framer

GetFrames returns frames with corresponding id. It returns nil if there is no frames with given id. 

Example ¶
    
    
    tag, err := id3v2.Open("file.mp3", id3v2.Options{Parse: true})
    if tag == nil || err != nil {
    	log.Fatal("Error while opening mp3 file: ", err)
    }
    
    pictures := tag.GetFrames(tag.CommonID("Attached picture"))
    for _, f := range pictures {
    	pic, ok := f.(id3v2.PictureFrame)
    	if !ok {
    		log.Fatal("Couldn't assert picture frame")
    	}
    	// Do something with picture frame.
    	// For example, print description of picture frame:
    	fmt.Println(pic.Description)
    }
    

####  func (*Tag) [GetLastFrame](https://github.com/bogem/id3v2/blob/v2.1.4/v2/tag.go#L164) ¶
    
    
    func (tag *Tag) GetLastFrame(id [string](/builtin#string)) Framer

GetLastFrame returns last frame from slice, that is returned from GetFrames function. GetLastFrame is suitable for frames, that can be only one in whole tag. For example, for text frames. 

Example ¶
    
    
    tag, err := id3v2.Open("file.mp3", id3v2.Options{Parse: true})
    if tag == nil || err != nil {
    	log.Fatal("Error while opening mp3 file: ", err)
    }
    
    bpmFramer := tag.GetLastFrame(tag.CommonID("BPM"))
    if bpmFramer != nil {
    	bpm, ok := bpmFramer.(id3v2.TextFrame)
    	if !ok {
    		log.Fatal("Couldn't assert bpm frame")
    	}
    	fmt.Println(bpm.Text)
    }
    

####  func (*Tag) [GetTextFrame](https://github.com/bogem/id3v2/blob/v2.1.4/v2/tag.go#L179) ¶
    
    
    func (tag *Tag) GetTextFrame(id [string](/builtin#string)) TextFrame

GetTextFrame returns text frame with corresponding id. 

####  func (*Tag) [HasFrames](https://github.com/bogem/id3v2/blob/v2.1.4/v2/tag.go#L221) ¶
    
    
    func (tag *Tag) HasFrames() [bool](/builtin#bool)

HasFrames checks if there is at least one frame in tag. It's much faster than tag.Count() > 0. 

####  func (*Tag) [Reset](https://github.com/bogem/id3v2/blob/v2.1.4/v2/tag.go#L145) ¶
    
    
    func (tag *Tag) Reset(rd [io](/io).[Reader](/io#Reader), opts Options) [error](/builtin#error)

Reset deletes all frames in tag and parses rd considering opts. 

####  func (*Tag) [Save](https://github.com/bogem/id3v2/blob/v2.1.4/v2/tag.go#L321) ¶
    
    
    func (tag *Tag) Save() [error](/builtin#error)

Save writes tag to the file, if tag was opened with a file. If there are no frames in tag, Save will write only music part without any ID3v2 information. If tag was initiliazed not with file, it returns ErrNoFile. 

####  func (*Tag) [SetAlbum](https://github.com/bogem/id3v2/blob/v2.1.4/v2/tag.go#L245) ¶
    
    
    func (tag *Tag) SetAlbum(album [string](/builtin#string))

####  func (*Tag) [SetArtist](https://github.com/bogem/id3v2/blob/v2.1.4/v2/tag.go#L237) ¶
    
    
    func (tag *Tag) SetArtist(artist [string](/builtin#string))

####  func (*Tag) [SetDefaultEncoding](https://github.com/bogem/id3v2/blob/v2.1.4/v2/tag.go#L198) ¶
    
    
    func (tag *Tag) SetDefaultEncoding(encoding Encoding)

SetDefaultEncoding sets default encoding for tag. Default encoding is used in methods (e.g. SetArtist, SetAlbum ...) for setting text frames without explicit providing encoding. 

####  func (*Tag) [SetGenre](https://github.com/bogem/id3v2/blob/v2.1.4/v2/tag.go#L261) ¶
    
    
    func (tag *Tag) SetGenre(genre [string](/builtin#string))

####  func (*Tag) [SetTitle](https://github.com/bogem/id3v2/blob/v2.1.4/v2/tag.go#L229) ¶
    
    
    func (tag *Tag) SetTitle(title [string](/builtin#string))

####  func (*Tag) [SetVersion](https://github.com/bogem/id3v2/blob/v2.1.4/v2/tag.go#L309) ¶
    
    
    func (tag *Tag) SetVersion(version [byte](/builtin#byte))

SetVersion sets given ID3v2 version to tag. If version is less than 3 or greater than 4, then this method will do nothing. If tag has some frames, which are deprecated or changed in given version, then to your notice you can delete, change or just stay them. 

####  func (*Tag) [SetYear](https://github.com/bogem/id3v2/blob/v2.1.4/v2/tag.go#L253) ¶
    
    
    func (tag *Tag) SetYear(year [string](/builtin#string))

####  func (*Tag) [Size](https://github.com/bogem/id3v2/blob/v2.1.4/v2/tag.go#L285) ¶
    
    
    func (tag *Tag) Size() [int](/builtin#int)

Size returns the size of tag (tag header + size of all frames) in bytes. 

####  func (*Tag) [Title](https://github.com/bogem/id3v2/blob/v2.1.4/v2/tag.go#L225) ¶
    
    
    func (tag *Tag) Title() [string](/builtin#string)

####  func (*Tag) [Version](https://github.com/bogem/id3v2/blob/v2.1.4/v2/tag.go#L301) ¶
    
    
    func (tag *Tag) Version() [byte](/builtin#byte)

Version returns current ID3v2 version of tag. 

####  func (*Tag) [WriteTo](https://github.com/bogem/id3v2/blob/v2.1.4/v2/tag.go#L395) ¶
    
    
    func (tag *Tag) WriteTo(w [io](/io).[Writer](/io#Writer)) (n [int64](/builtin#int64), err [error](/builtin#error))

WriteTo writes whole tag in w if there is at least one frame. It returns the number of bytes written and error during the write. It returns nil as error if the write was successful. 

####  func (*Tag) [Year](https://github.com/bogem/id3v2/blob/v2.1.4/v2/tag.go#L249) ¶
    
    
    func (tag *Tag) Year() [string](/builtin#string)

####  type [TextFrame](https://github.com/bogem/id3v2/blob/v2.1.4/v2/text_frame.go#L11) ¶
    
    
    type TextFrame struct {
    	Encoding Encoding
    	Text     [string](/builtin#string)
    }

TextFrame is used to work with all text frames (all T*** frames like TIT2 (title), TALB (album) and so on). 

Example (Add) ¶
    
    
    tag := id3v2.NewEmptyTag()
    textFrame := id3v2.TextFrame{
    	Encoding: id3v2.EncodingUTF8,
    	Text:     "Happy",
    }
    tag.AddFrame(tag.CommonID("Mood"), textFrame)
    

Example (Get) ¶
    
    
    tag, err := id3v2.Open("file.mp3", id3v2.Options{Parse: true})
    if tag == nil || err != nil {
    	log.Fatal("Error while opening mp3 file: ", err)
    }
    
    tf := tag.GetTextFrame(tag.CommonID("Mood"))
    fmt.Println(tf.Text)
    

####  func (TextFrame) [Size](https://github.com/bogem/id3v2/blob/v2.1.4/v2/text_frame.go#L16) ¶
    
    
    func (tf TextFrame) Size() [int](/builtin#int)

####  func (TextFrame) [UniqueIdentifier](https://github.com/bogem/id3v2/blob/v2.1.4/v2/text_frame.go#L20) ¶
    
    
    func (tf TextFrame) UniqueIdentifier() [string](/builtin#string)

####  func (TextFrame) [WriteTo](https://github.com/bogem/id3v2/blob/v2.1.4/v2/text_frame.go#L24) ¶
    
    
    func (tf TextFrame) WriteTo(w [io](/io).[Writer](/io#Writer)) ([int64](/builtin#int64), [error](/builtin#error))

####  type [UFIDFrame](https://github.com/bogem/id3v2/blob/v2.1.4/v2/ufid_frame.go#L6) ¶
    
    
    type UFIDFrame struct {
    	OwnerIdentifier [string](/builtin#string)
    	Identifier      [][byte](/builtin#byte)
    }

UFIDFrame is used for "Unique file identifier" 

####  func (UFIDFrame) [Size](https://github.com/bogem/id3v2/blob/v2.1.4/v2/ufid_frame.go#L15) ¶
    
    
    func (ufid UFIDFrame) Size() [int](/builtin#int)

####  func (UFIDFrame) [UniqueIdentifier](https://github.com/bogem/id3v2/blob/v2.1.4/v2/ufid_frame.go#L11) ¶
    
    
    func (ufid UFIDFrame) UniqueIdentifier() [string](/builtin#string)

####  func (UFIDFrame) [WriteTo](https://github.com/bogem/id3v2/blob/v2.1.4/v2/ufid_frame.go#L19) ¶
    
    
    func (ufid UFIDFrame) WriteTo(w [io](/io).[Writer](/io#Writer)) (n [int64](/builtin#int64), err [error](/builtin#error))

####  type [UnknownFrame](https://github.com/bogem/id3v2/blob/v2.1.4/v2/unknown_frame.go#L20) ¶
    
    
    type UnknownFrame struct {
    	Body [][byte](/builtin#byte)
    }

UnknownFrame is used for frames, which id3v2 so far doesn't know how to parse and write it. It just contains an unparsed byte body of the frame. 

####  func (UnknownFrame) [Size](https://github.com/bogem/id3v2/blob/v2.1.4/v2/unknown_frame.go#L29) ¶
    
    
    func (uf UnknownFrame) Size() [int](/builtin#int)

####  func (UnknownFrame) [UniqueIdentifier](https://github.com/bogem/id3v2/blob/v2.1.4/v2/unknown_frame.go#L24) ¶
    
    
    func (uf UnknownFrame) UniqueIdentifier() [string](/builtin#string)

####  func (UnknownFrame) [WriteTo](https://github.com/bogem/id3v2/blob/v2.1.4/v2/unknown_frame.go#L33) ¶
    
    
    func (uf UnknownFrame) WriteTo(w [io](/io).[Writer](/io#Writer)) (n [int64](/builtin#int64), err [error](/builtin#error))

####  type [UnsynchronisedLyricsFrame](https://github.com/bogem/id3v2/blob/v2.1.4/v2/unsynchronised_lyrics_frame.go#L15) ¶
    
    
    type UnsynchronisedLyricsFrame struct {
    	Encoding          Encoding
    	Language          [string](/builtin#string)
    	ContentDescriptor [string](/builtin#string)
    	Lyrics            [string](/builtin#string)
    }

UnsynchronisedLyricsFrame is used to work with USLT frames. The information about how to add unsynchronised lyrics/text frame to tag you can see in the docs to tag.AddUnsynchronisedLyricsFrame function. 

You must choose a three-letter language code from ISO 639-2 code list: <https://www.loc.gov/standards/iso639-2/php/code_list.php>

Example (Add) ¶
    
    
    tag := id3v2.NewEmptyTag()
    uslt := id3v2.UnsynchronisedLyricsFrame{
    	Encoding:          id3v2.EncodingUTF8,
    	Language:          "ger",
    	ContentDescriptor: "Deutsche Nationalhymne",
    	Lyrics:            "Einigkeit und Recht und Freiheit...",
    }
    tag.AddUnsynchronisedLyricsFrame(uslt)
    

Example (Get) ¶
    
    
    tag, err := id3v2.Open("file.mp3", id3v2.Options{Parse: true})
    if tag == nil || err != nil {
    	log.Fatal("Error while opening mp3 file: ", err)
    }
    
    uslfs := tag.GetFrames(tag.CommonID("Unsynchronised lyrics/text transcription"))
    for _, f := range uslfs {
    	uslf, ok := f.(id3v2.UnsynchronisedLyricsFrame)
    	if !ok {
    		log.Fatal("Couldn't assert USLT frame")
    	}
    
    	// Do something with USLT frame.
    	// For example, print the lyrics:
    	fmt.Println(uslf.Lyrics)
    }
    

####  func (UnsynchronisedLyricsFrame) [Size](https://github.com/bogem/id3v2/blob/v2.1.4/v2/unsynchronised_lyrics_frame.go#L22) ¶
    
    
    func (uslf UnsynchronisedLyricsFrame) Size() [int](/builtin#int)

####  func (UnsynchronisedLyricsFrame) [UniqueIdentifier](https://github.com/bogem/id3v2/blob/v2.1.4/v2/unsynchronised_lyrics_frame.go#L27) ¶
    
    
    func (uslf UnsynchronisedLyricsFrame) UniqueIdentifier() [string](/builtin#string)

####  func (UnsynchronisedLyricsFrame) [WriteTo](https://github.com/bogem/id3v2/blob/v2.1.4/v2/unsynchronised_lyrics_frame.go#L31) ¶
    
    
    func (uslf UnsynchronisedLyricsFrame) WriteTo(w [io](/io).[Writer](/io#Writer)) (n [int64](/builtin#int64), err [error](/builtin#error))

####  type [UserDefinedTextFrame](https://github.com/bogem/id3v2/blob/v2.1.4/v2/user_defined_text_frame.go#L7) ¶
    
    
    type UserDefinedTextFrame struct {
    	Encoding    Encoding
    	Description [string](/builtin#string)
    	Value       [string](/builtin#string)
    }

UserDefinedTextFrame is used to work with TXXX frames. There can be many UserDefinedTextFrames but the Description fields need to be unique. 

####  func (UserDefinedTextFrame) [Size](https://github.com/bogem/id3v2/blob/v2.1.4/v2/user_defined_text_frame.go#L13) ¶
    
    
    func (udtf UserDefinedTextFrame) Size() [int](/builtin#int)

####  func (UserDefinedTextFrame) [UniqueIdentifier](https://github.com/bogem/id3v2/blob/v2.1.4/v2/user_defined_text_frame.go#L17) ¶
    
    
    func (udtf UserDefinedTextFrame) UniqueIdentifier() [string](/builtin#string)

####  func (UserDefinedTextFrame) [WriteTo](https://github.com/bogem/id3v2/blob/v2.1.4/v2/user_defined_text_frame.go#L21) ¶
    
    
    func (udtf UserDefinedTextFrame) WriteTo(w [io](/io).[Writer](/io#Writer)) (n [int64](/builtin#int64), err [error](/builtin#error))
  *[↑]: Back to Top
  *[v]: View this template
  *[t]: Discuss this template
  *[e]: Edit this template
