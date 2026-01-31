# bogem/id3v2

> Source: https://pkg.go.dev/github.com/bogem/id3v2/v2
> Fetched: 2026-01-30T23:54:00.988089+00:00
> Content-Hash: ec4b274083e9b730
> Type: html

---

Overview

¶

Package id3v2 is the ID3 parsing and writing library for Go.

Example

¶

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

Example (Concurrent)

¶

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

Index

¶

Constants

Variables

type ChapterFrame

func (cf ChapterFrame) Size() int

func (cf ChapterFrame) UniqueIdentifier() string

func (cf ChapterFrame) WriteTo(w io.Writer) (n int64, err error)

type CommentFrame

func (cf CommentFrame) Size() int

func (cf CommentFrame) UniqueIdentifier() string

func (cf CommentFrame) WriteTo(w io.Writer) (n int64, err error)

type Encoding

func (e Encoding) Equals(other Encoding) bool

func (e Encoding) String() string

type Framer

type Options

type PictureFrame

func (pf PictureFrame) Size() int

func (pf PictureFrame) UniqueIdentifier() string

func (pf PictureFrame) WriteTo(w io.Writer) (n int64, err error)

type PopularimeterFrame

func (pf PopularimeterFrame) Size() int

func (pf PopularimeterFrame) UniqueIdentifier() string

func (pf PopularimeterFrame) WriteTo(w io.Writer) (n int64, err error)

type Tag

func NewEmptyTag() *Tag

func Open(name string, opts Options) (*Tag, error)

func ParseReader(rd io.Reader, opts Options) (*Tag, error)

func (tag *Tag) AddAttachedPicture(pf PictureFrame)

func (tag *Tag) AddChapterFrame(cf ChapterFrame)

func (tag *Tag) AddCommentFrame(cf CommentFrame)

func (tag *Tag) AddFrame(id string, f Framer)

func (tag *Tag) AddTextFrame(id string, encoding Encoding, text string)

func (tag *Tag) AddUFIDFrame(ufid UFIDFrame)

func (tag *Tag) AddUnsynchronisedLyricsFrame(uslf UnsynchronisedLyricsFrame)

func (tag *Tag) AddUserDefinedTextFrame(udtf UserDefinedTextFrame)

func (tag *Tag) Album() string

func (tag *Tag) AllFrames() map[string][]Framer

func (tag *Tag) Artist() string

func (tag *Tag) Close() error

func (tag *Tag) CommonID(description string) string

func (tag *Tag) Count() int

func (tag *Tag) DefaultEncoding() Encoding

func (tag *Tag) DeleteAllFrames()

func (tag *Tag) DeleteFrames(id string)

func (tag *Tag) Genre() string

func (tag *Tag) GetFrames(id string) []Framer

func (tag *Tag) GetLastFrame(id string) Framer

func (tag *Tag) GetTextFrame(id string) TextFrame

func (tag *Tag) HasFrames() bool

func (tag *Tag) Reset(rd io.Reader, opts Options) error

func (tag *Tag) Save() error

func (tag *Tag) SetAlbum(album string)

func (tag *Tag) SetArtist(artist string)

func (tag *Tag) SetDefaultEncoding(encoding Encoding)

func (tag *Tag) SetGenre(genre string)

func (tag *Tag) SetTitle(title string)

func (tag *Tag) SetVersion(version byte)

func (tag *Tag) SetYear(year string)

func (tag *Tag) Size() int

func (tag *Tag) Title() string

func (tag *Tag) Version() byte

func (tag *Tag) WriteTo(w io.Writer) (n int64, err error)

func (tag *Tag) Year() string

type TextFrame

func (tf TextFrame) Size() int

func (tf TextFrame) UniqueIdentifier() string

func (tf TextFrame) WriteTo(w io.Writer) (int64, error)

type UFIDFrame

func (ufid UFIDFrame) Size() int

func (ufid UFIDFrame) UniqueIdentifier() string

func (ufid UFIDFrame) WriteTo(w io.Writer) (n int64, err error)

type UnknownFrame

func (uf UnknownFrame) Size() int

func (uf UnknownFrame) UniqueIdentifier() string

func (uf UnknownFrame) WriteTo(w io.Writer) (n int64, err error)

type UnsynchronisedLyricsFrame

func (uslf UnsynchronisedLyricsFrame) Size() int

func (uslf UnsynchronisedLyricsFrame) UniqueIdentifier() string

func (uslf UnsynchronisedLyricsFrame) WriteTo(w io.Writer) (n int64, err error)

type UserDefinedTextFrame

func (udtf UserDefinedTextFrame) Size() int

func (udtf UserDefinedTextFrame) UniqueIdentifier() string

func (udtf UserDefinedTextFrame) WriteTo(w io.Writer) (n int64, err error)

Examples

¶

Package

Package (Concurrent)

CommentFrame (Add)

CommentFrame (Get)

PictureFrame (Add)

PictureFrame (Get)

PopularimeterFrame (Add)

PopularimeterFrame (Get)

Tag.AddAttachedPicture

Tag.AddCommentFrame

Tag.AddUnsynchronisedLyricsFrame

Tag.GetFrames

Tag.GetLastFrame

TextFrame (Add)

TextFrame (Get)

UnsynchronisedLyricsFrame (Add)

UnsynchronisedLyricsFrame (Get)

Constants

¶

View Source

const (

PTOther =

iota

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

View Source

const (

IgnoredOffset = 0xFFFFFFFF

)

Variables

¶

View Source

var (

V23CommonIDs = map[

string

]

string

{

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

V24CommonIDs = map[

string

]

string

{

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

View Source

var (

// EncodingISO is ISO-8859-1 encoding.

EncodingISO =

Encoding

{
		Name:             "ISO-8859-1",
		Key:              0,
		TerminationBytes: []

byte

{0},
	}

// EncodingUTF16 is UTF-16 encoded Unicode with BOM.

EncodingUTF16 =

Encoding

{
		Name:             "UTF-16 encoded Unicode with BOM",
		Key:              1,
		TerminationBytes: []

byte

{0, 0},
	}

// EncodingUTF16BE is UTF-16BE encoded Unicode without BOM.

EncodingUTF16BE =

Encoding

{
		Name:             "UTF-16BE encoded Unicode without BOM",
		Key:              2,
		TerminationBytes: []

byte

{0, 0},
	}

// EncodingUTF8 is UTF-8 encoded Unicode.

EncodingUTF8 =

Encoding

{
		Name:             "UTF-8 encoded Unicode",
		Key:              3,
		TerminationBytes: []

byte

{0},
	}
)

Available encodings.

View Source

var ErrBodyOverflow =

errors

.

New

("frame went over tag area")

ErrBodyOverflow is returned when a frame has greater size than the remaining tag size

View Source

var ErrInvalidLanguageLength =

errors

.

New

("language code must consist of three letters according to ISO 639-2")

View Source

var ErrInvalidSizeFormat =

errors

.

New

("invalid format of tag's/frame's size")

View Source

var ErrNoFile =

errors

.

New

("tag was not initialized with file")

View Source

var ErrSizeOverflow =

errors

.

New

("size of tag/frame is greater than allowed in id3 tag")

View Source

var ErrSmallHeaderSize =

errors

.

New

("size of tag header is less than expected")

View Source

var ErrUnsupportedVersion =

errors

.

New

("unsupported version of ID3 tag")

Functions

¶

This section is empty.

Types

¶

type

ChapterFrame

¶

type ChapterFrame struct {

ElementID

string

StartTime

time

.

Duration

EndTime

time

.

Duration

StartOffset

uint32

EndOffset

uint32

Title       *

TextFrame

Description *

TextFrame

}

ChapterFrame is used to work with CHAP frames
according to spec from

http://id3.org/id3v2-chapters-1.0

This implementation only supports single TIT2 subframe (Title field).
All other subframes are ignored.
If StartOffset or EndOffset == id3v2.IgnoredOffset, then it should be ignored
and StartTime or EndTime should be utilized

func (ChapterFrame)

Size

¶

func (cf

ChapterFrame

) Size()

int

func (ChapterFrame)

UniqueIdentifier

¶

func (cf

ChapterFrame

) UniqueIdentifier()

string

func (ChapterFrame)

WriteTo

¶

func (cf

ChapterFrame

) WriteTo(w

io

.

Writer

) (n

int64

, err

error

)

type

CommentFrame

¶

type CommentFrame struct {

Encoding

Encoding

Language

string

Description

string

Text

string

}

CommentFrame is used to work with COMM frames.
The information about how to add comment frame to tag you can
see in the docs to tag.AddCommentFrame function.

You must choose a three-letter language code from
ISO 639-2 code list:

https://www.loc.gov/standards/iso639-2/php/code_list.php

Example (Add)

¶

tag := id3v2.NewEmptyTag()
comment := id3v2.CommentFrame{
	Encoding:    id3v2.EncodingUTF8,
	Language:    "eng",
	Description: "My opinion",
	Text:        "Very good song",
}
tag.AddCommentFrame(comment)

Example (Get)

¶

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

func (CommentFrame)

Size

¶

func (cf

CommentFrame

) Size()

int

func (CommentFrame)

UniqueIdentifier

¶

func (cf

CommentFrame

) UniqueIdentifier()

string

func (CommentFrame)

WriteTo

¶

func (cf

CommentFrame

) WriteTo(w

io

.

Writer

) (n

int64

, err

error

)

type

Encoding

¶

type Encoding struct {

Name

string

Key

byte

TerminationBytes []

byte

}

Encoding is a struct for encodings.

func (Encoding)

Equals

¶

func (e

Encoding

) Equals(other

Encoding

)

bool

func (Encoding)

String

¶

func (e

Encoding

) String()

string

type

Framer

¶

type Framer interface {

// Size returns the size of frame body.

Size()

int

// UniqueIdentifier returns the string that makes this frame unique from others.

// For example, some frames with same id can be added in tag, but they should be differ in other properties.

// E.g. It would be "Description" for TXXX and APIC.

//

// Frames that can be added only once with same id (e.g. all text frames) should return just "".

UniqueIdentifier()

string

// WriteTo writes body slice into w.

WriteTo(w

io

.

Writer

) (n

int64

, err

error

)
}

Framer provides a generic interface for frames.
You can create your own frames. They must implement only this interface.

type

Options

¶

type Options struct {

// Parse defines, if tag will be parsed.

Parse

bool

// ParseFrames defines, that frames do you only want to parse. For example,

// `ParseFrames: []string{"Artist", "Title"}` will only parse artist

// and title frames. You can specify IDs ("TPE1", "TIT2") as well as

// descriptions ("Artist", "Title"). If ParseFrame is blank or nil,

// id3v2 will parse all frames in tag. It works only if Parse is true.

//

// It's very useful for performance, so for example

// if you want to get only some text frames,

// id3v2 will not parse huge picture or unknown frames.

ParseFrames []

string

}

Options influence on processing the tag.

type

PictureFrame

¶

type PictureFrame struct {

Encoding

Encoding

MimeType

string

PictureType

byte

Description

string

Picture     []

byte

}

PictureFrame structure is used for picture frames (APIC).
The information about how to add picture frame to tag you can
see in the docs to tag.AddAttachedPicture function.

Available picture types you can see in constants.

Example (Add)

¶

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

Example (Get)

¶

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

func (PictureFrame)

Size

¶

func (pf

PictureFrame

) Size()

int

func (PictureFrame)

UniqueIdentifier

¶

func (pf

PictureFrame

) UniqueIdentifier()

string

func (PictureFrame)

WriteTo

¶

func (pf

PictureFrame

) WriteTo(w

io

.

Writer

) (n

int64

, err

error

)

type

PopularimeterFrame

¶

type PopularimeterFrame struct {

// Email is the identifier for a POPM frame.

Email

string

// The rating is 1-255 where 1 is worst and 255 is best. 0 is unknown.

Rating

uint8

// Counter is the number of times this file has been played by this email.

Counter *

big

.

Int

}

PopularimeterFrame structure is used for Popularimeter (POPM).

https://id3.org/id3v2.3.0#Popularimeter

Example (Add)

¶

tag := id3v2.NewEmptyTag()

popmFrame := id3v2.PopularimeterFrame{
	Email:   "foo@bar.com",
	Rating:  128,
	Counter: big.NewInt(10000000000000000),
}
tag.AddFrame(tag.CommonID("Popularimeter"), popmFrame)

Example (Get)

¶

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

func (PopularimeterFrame)

Size

¶

func (pf

PopularimeterFrame

) Size()

int

func (PopularimeterFrame)

UniqueIdentifier

¶

func (pf

PopularimeterFrame

) UniqueIdentifier()

string

func (PopularimeterFrame)

WriteTo

¶

func (pf

PopularimeterFrame

) WriteTo(w

io

.

Writer

) (n

int64

, err

error

)

type

Tag

¶

type Tag struct {

// contains filtered or unexported fields

}

Tag stores all information about opened tag.

func

NewEmptyTag

¶

func NewEmptyTag() *

Tag

NewEmptyTag returns an empty ID3v2.4 tag without any frames and reader.

func

Open

¶

func Open(name

string

, opts

Options

) (*

Tag

,

error

)

Open opens file with name and passes it to OpenFile.
If there is no tag in file, it will create new one with version ID3v2.4.

func

ParseReader

¶

func ParseReader(rd

io

.

Reader

, opts

Options

) (*

Tag

,

error

)

ParseReader parses rd and finds tag in it considering opts.
If there is no tag in rd, it will create new one with version ID3v2.4.

func (*Tag)

AddAttachedPicture

¶

func (tag *

Tag

) AddAttachedPicture(pf

PictureFrame

)

AddAttachedPicture adds the picture frame to tag.

Example

¶

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

func (*Tag)

AddChapterFrame

¶

func (tag *

Tag

) AddChapterFrame(cf

ChapterFrame

)

AddChapterFrame adds the chapter frame to tag.

func (*Tag)

AddCommentFrame

¶

func (tag *

Tag

) AddCommentFrame(cf

CommentFrame

)

AddCommentFrame adds the comment frame to tag.

Example

¶

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

func (*Tag)

AddFrame

¶

func (tag *

Tag

) AddFrame(id

string

, f

Framer

)

AddFrame adds f to tag with appropriate id. If id is "" or f is nil,
AddFrame will not add it to tag.

If you want to add attached picture, comment or unsynchronised lyrics/text
transcription frames, better use AddAttachedPicture, AddCommentFrame
or AddUnsynchronisedLyricsFrame methods respectively.

func (*Tag)

AddTextFrame

¶

func (tag *

Tag

) AddTextFrame(id

string

, encoding

Encoding

, text

string

)

AddTextFrame creates the text frame with provided encoding and text
and adds to tag.

func (*Tag)

AddUFIDFrame

¶

func (tag *

Tag

) AddUFIDFrame(ufid

UFIDFrame

)

AddUFIDFrame adds the unique file identifier frame (UFID) to tag.

func (*Tag)

AddUnsynchronisedLyricsFrame

¶

func (tag *

Tag

) AddUnsynchronisedLyricsFrame(uslf

UnsynchronisedLyricsFrame

)

AddUnsynchronisedLyricsFrame adds the unsynchronised lyrics/text frame
to tag.

Example

¶

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

func (*Tag)

AddUserDefinedTextFrame

¶

func (tag *

Tag

) AddUserDefinedTextFrame(udtf

UserDefinedTextFrame

)

AddUserDefinedTextFrame adds the custom frame (TXXX) to tag.

func (*Tag)

Album

¶

func (tag *

Tag

) Album()

string

func (*Tag)

AllFrames

¶

func (tag *

Tag

) AllFrames() map[

string

][]

Framer

AllFrames returns map, that contains all frames in tag, that could be parsed.
The key of this map is an ID of frame and value is an array of frames.

func (*Tag)

Artist

¶

func (tag *

Tag

) Artist()

string

func (*Tag)

Close

¶

func (tag *

Tag

) Close()

error

Close closes tag's file, if tag was opened with a file.
If tag was initiliazed not with file, it returns ErrNoFile.

func (*Tag)

CommonID

¶

func (tag *

Tag

) CommonID(description

string

)

string

CommonID returns frame ID from given description.
For example, CommonID("Language") will return "TLAN".
If it can't find the ID with given description, it returns the description.

All descriptions you can find in file common_ids.go
or in id3 documentation.
v2.3:

http://id3.org/id3v2.3.0#Declared_ID3v2_frames

v2.4:

http://id3.org/id3v2.4.0-frames

func (*Tag)

Count

¶

func (tag *

Tag

) Count()

int

Count returns the number of frames in tag.

func (*Tag)

DefaultEncoding

¶

func (tag *

Tag

) DefaultEncoding()

Encoding

DefaultEncoding returns default encoding of tag.
Default encoding is used in methods (e.g. SetArtist, SetAlbum ...) for
setting text frames without the explicit providing of encoding.

func (*Tag)

DeleteAllFrames

¶

func (tag *

Tag

) DeleteAllFrames()

DeleteAllFrames deletes all frames in tag.

func (*Tag)

DeleteFrames

¶

func (tag *

Tag

) DeleteFrames(id

string

)

DeleteFrames deletes frames in tag with given id.

func (*Tag)

Genre

¶

func (tag *

Tag

) Genre()

string

func (*Tag)

GetFrames

¶

func (tag *

Tag

) GetFrames(id

string

) []

Framer

GetFrames returns frames with corresponding id.
It returns nil if there is no frames with given id.

Example

¶

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

func (*Tag)

GetLastFrame

¶

func (tag *

Tag

) GetLastFrame(id

string

)

Framer

GetLastFrame returns last frame from slice, that is returned from GetFrames function.
GetLastFrame is suitable for frames, that can be only one in whole tag.
For example, for text frames.

Example

¶

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

func (*Tag)

GetTextFrame

¶

func (tag *

Tag

) GetTextFrame(id

string

)

TextFrame

GetTextFrame returns text frame with corresponding id.

func (*Tag)

HasFrames

¶

func (tag *

Tag

) HasFrames()

bool

HasFrames checks if there is at least one frame in tag.
It's much faster than tag.Count() > 0.

func (*Tag)

Reset

¶

func (tag *

Tag

) Reset(rd

io

.

Reader

, opts

Options

)

error

Reset deletes all frames in tag and parses rd considering opts.

func (*Tag)

Save

¶

func (tag *

Tag

) Save()

error

Save writes tag to the file, if tag was opened with a file.
If there are no frames in tag, Save will write
only music part without any ID3v2 information.
If tag was initiliazed not with file, it returns ErrNoFile.

func (*Tag)

SetAlbum

¶

func (tag *

Tag

) SetAlbum(album

string

)

func (*Tag)

SetArtist

¶

func (tag *

Tag

) SetArtist(artist

string

)

func (*Tag)

SetDefaultEncoding

¶

func (tag *

Tag

) SetDefaultEncoding(encoding

Encoding

)

SetDefaultEncoding sets default encoding for tag.
Default encoding is used in methods (e.g. SetArtist, SetAlbum ...) for
setting text frames without explicit providing encoding.

func (*Tag)

SetGenre

¶

func (tag *

Tag

) SetGenre(genre

string

)

func (*Tag)

SetTitle

¶

func (tag *

Tag

) SetTitle(title

string

)

func (*Tag)

SetVersion

¶

func (tag *

Tag

) SetVersion(version

byte

)

SetVersion sets given ID3v2 version to tag.
If version is less than 3 or greater than 4, then this method will do nothing.
If tag has some frames, which are deprecated or changed in given version,
then to your notice you can delete, change or just stay them.

func (*Tag)

SetYear

¶

func (tag *

Tag

) SetYear(year

string

)

func (*Tag)

Size

¶

func (tag *

Tag

) Size()

int

Size returns the size of tag (tag header + size of all frames) in bytes.

func (*Tag)

Title

¶

func (tag *

Tag

) Title()

string

func (*Tag)

Version

¶

func (tag *

Tag

) Version()

byte

Version returns current ID3v2 version of tag.

func (*Tag)

WriteTo

¶

func (tag *

Tag

) WriteTo(w

io

.

Writer

) (n

int64

, err

error

)

WriteTo writes whole tag in w if there is at least one frame.
It returns the number of bytes written and error during the write.
It returns nil as error if the write was successful.

func (*Tag)

Year

¶

func (tag *

Tag

) Year()

string

type

TextFrame

¶

type TextFrame struct {

Encoding

Encoding

Text

string

}

TextFrame is used to work with all text frames
(all T*** frames like TIT2 (title), TALB (album) and so on).

Example (Add)

¶

tag := id3v2.NewEmptyTag()
textFrame := id3v2.TextFrame{
	Encoding: id3v2.EncodingUTF8,
	Text:     "Happy",
}
tag.AddFrame(tag.CommonID("Mood"), textFrame)

Example (Get)

¶

tag, err := id3v2.Open("file.mp3", id3v2.Options{Parse: true})
if tag == nil || err != nil {
	log.Fatal("Error while opening mp3 file: ", err)
}

tf := tag.GetTextFrame(tag.CommonID("Mood"))
fmt.Println(tf.Text)

func (TextFrame)

Size

¶

func (tf

TextFrame

) Size()

int

func (TextFrame)

UniqueIdentifier

¶

func (tf

TextFrame

) UniqueIdentifier()

string

func (TextFrame)

WriteTo

¶

func (tf

TextFrame

) WriteTo(w

io

.

Writer

) (

int64

,

error

)

type

UFIDFrame

¶

type UFIDFrame struct {

OwnerIdentifier

string

Identifier      []

byte

}

UFIDFrame is used for "Unique file identifier"

func (UFIDFrame)

Size

¶

func (ufid

UFIDFrame

) Size()

int

func (UFIDFrame)

UniqueIdentifier

¶

func (ufid

UFIDFrame

) UniqueIdentifier()

string

func (UFIDFrame)

WriteTo

¶

func (ufid

UFIDFrame

) WriteTo(w

io

.

Writer

) (n

int64

, err

error

)

type

UnknownFrame

¶

type UnknownFrame struct {

Body []

byte

}

UnknownFrame is used for frames, which id3v2 so far doesn't know how to
parse and write it. It just contains an unparsed byte body of the frame.

func (UnknownFrame)

Size

¶

func (uf

UnknownFrame

) Size()

int

func (UnknownFrame)

UniqueIdentifier

¶

func (uf

UnknownFrame

) UniqueIdentifier()

string

func (UnknownFrame)

WriteTo

¶

func (uf

UnknownFrame

) WriteTo(w

io

.

Writer

) (n

int64

, err

error

)

type

UnsynchronisedLyricsFrame

¶

type UnsynchronisedLyricsFrame struct {

Encoding

Encoding

Language

string

ContentDescriptor

string

Lyrics

string

}

UnsynchronisedLyricsFrame is used to work with USLT frames.
The information about how to add unsynchronised lyrics/text frame to tag
you can see in the docs to tag.AddUnsynchronisedLyricsFrame function.

You must choose a three-letter language code from
ISO 639-2 code list:

https://www.loc.gov/standards/iso639-2/php/code_list.php

Example (Add)

¶

tag := id3v2.NewEmptyTag()
uslt := id3v2.UnsynchronisedLyricsFrame{
	Encoding:          id3v2.EncodingUTF8,
	Language:          "ger",
	ContentDescriptor: "Deutsche Nationalhymne",
	Lyrics:            "Einigkeit und Recht und Freiheit...",
}
tag.AddUnsynchronisedLyricsFrame(uslt)

Example (Get)

¶

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

func (UnsynchronisedLyricsFrame)

Size

¶

func (uslf

UnsynchronisedLyricsFrame

) Size()

int

func (UnsynchronisedLyricsFrame)

UniqueIdentifier

¶

func (uslf

UnsynchronisedLyricsFrame

) UniqueIdentifier()

string

func (UnsynchronisedLyricsFrame)

WriteTo

¶

func (uslf

UnsynchronisedLyricsFrame

) WriteTo(w

io

.

Writer

) (n

int64

, err

error

)

type

UserDefinedTextFrame

¶

type UserDefinedTextFrame struct {

Encoding

Encoding

Description

string

Value

string

}

UserDefinedTextFrame is used to work with TXXX frames.
There can be many UserDefinedTextFrames but the Description fields need to be unique.

func (UserDefinedTextFrame)

Size

¶

func (udtf

UserDefinedTextFrame

) Size()

int

func (UserDefinedTextFrame)

UniqueIdentifier

¶

func (udtf

UserDefinedTextFrame

) UniqueIdentifier()

string

func (UserDefinedTextFrame)

WriteTo

¶

func (udtf

UserDefinedTextFrame

) WriteTo(w

io

.

Writer

) (n

int64

, err

error

)