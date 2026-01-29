# dhowden/tag (audio metadata)

> Auto-fetched from [https://pkg.go.dev/github.com/dhowden/tag](https://pkg.go.dev/github.com/dhowden/tag)
> Last Updated: 2026-01-28T21:46:53.499272+00:00

---

Overview
¶
Package tag provides MP3 (ID3: v1, 2.2, 2.3 and 2.4), MP4, FLAC and OGG metadata detection,
parsing and artwork extraction.
Detect and parse tag metadata from an io.ReadSeeker (i.e. an *os.File):
m, err := tag.ReadFrom(f)
if err != nil {
log.Fatal(err)
}
log.Print(m.Format()) // The detected format.
log.Print(m.Title())  // The title of the track (see Metadata interface for more details).
Index
¶
Variables
func Identify(r io.ReadSeeker) (format Format, fileType FileType, err error)
func Sum(r io.ReadSeeker) (string, error)
func SumAll(r io.ReadSeeker) (string, error)
func SumAtoms(r io.ReadSeeker) (string, error)
func SumFLAC(r io.ReadSeeker) (string, error)
func SumID3v1(r io.ReadSeeker) (string, error)
func SumID3v2(r io.ReadSeeker) (string, error)
type Comm
func (t Comm) String() string
type FileType
type Format
type Metadata
func ReadAtoms(r io.ReadSeeker) (Metadata, error)
func ReadDSFTags(r io.ReadSeeker) (Metadata, error)
func ReadFLACTags(r io.ReadSeeker) (Metadata, error)
func ReadFrom(r io.ReadSeeker) (Metadata, error)
func ReadID3v1Tags(r io.ReadSeeker) (Metadata, error)
func ReadID3v2Tags(r io.ReadSeeker) (Metadata, error)
func ReadOGGTags(r io.Reader) (Metadata, error)
type Picture
func (p Picture) String() string
type UFID
func (u UFID) String() string
Constants
¶
This section is empty.
Variables
¶
View Source
var DefaultUTF16WithBOMByteOrder
binary
.
ByteOrder
=
binary
.
LittleEndian
DefaultUTF16WithBOMByteOrder is the byte order used when the "UTF16 with BOM" encoding
is specified without a corresponding BOM in the data.
View Source
var ErrNoTagsFound =
errors
.
New
("no tags found")
ErrNoTagsFound is the error returned by ReadFrom when the metadata format
cannot be identified.
View Source
var ErrNotID3v1 =
errors
.
New
("invalid ID3v1 header")
ErrNotID3v1 is an error which is returned when no ID3v1 header is found.
Functions
¶
func
Identify
¶
func Identify(r
io
.
ReadSeeker
) (format
Format
, fileType
FileType
, err
error
)
Identify identifies the format and file type of the data in the ReadSeeker.
func
Sum
¶
func Sum(r
io
.
ReadSeeker
) (
string
,
error
)
Sum creates a checksum of the audio file data provided by the io.ReadSeeker which is metadata
(ID3, MP4) invariant.
func
SumAll
¶
func SumAll(r
io
.
ReadSeeker
) (
string
,
error
)
SumAll returns a checksum of the content from the reader (until EOF).
func
SumAtoms
¶
func SumAtoms(r
io
.
ReadSeeker
) (
string
,
error
)
SumAtoms constructs a checksum of MP4 audio file data provided by the io.ReadSeeker which is
metadata invariant.
func
SumFLAC
¶
func SumFLAC(r
io
.
ReadSeeker
) (
string
,
error
)
SumFLAC costructs a checksum of the FLAC audio file data provided by the io.ReadSeeker (ignores
metadata fields).
func
SumID3v1
¶
func SumID3v1(r
io
.
ReadSeeker
) (
string
,
error
)
SumID3v1 constructs a checksum of MP3 audio file data (assumed to have ID3v1 tags) provided
by the io.ReadSeeker which is metadata invariant.
func
SumID3v2
¶
func SumID3v2(r
io
.
ReadSeeker
) (
string
,
error
)
SumID3v2 constructs a checksum of MP3 audio file data (assumed to have ID3v2 tags) provided by the
io.ReadSeeker which is metadata invariant.
Types
¶
type
Comm
¶
type Comm struct {
Language
string
Description
string
Text
string
}
Comm is a type used in COMM, UFID, TXXX, WXXX and USLT tag.
It's a text with a description and a specified language
For WXXX, TXXX and UFID, we don't set a Language
func (Comm)
String
¶
func (t
Comm
) String()
string
String returns a string representation of the underlying Comm instance.
type
FileType
¶
type FileType
string
FileType is an enumeration of the audio file types supported by this package, in particular
there are audio file types which share metadata formats, and this type is used to distinguish
between them.
const (
UnknownFileType
FileType
= ""
// Unknown FileType.
MP3
FileType
= "MP3"
// MP3 file
M4A
FileType
= "M4A"
// M4A file Apple iTunes (ACC) Audio
M4B
FileType
= "M4B"
// M4A file Apple iTunes (ACC) Audio Book
M4P
FileType
= "M4P"
// M4A file Apple iTunes (ACC) AES Protected Audio
ALAC
FileType
= "ALAC"
// Apple Lossless file FIXME: actually detect this
FLAC
FileType
= "FLAC"
// FLAC file
OGG
FileType
= "OGG"
// OGG file
DSF
FileType
= "DSF"
// DSF file DSD Sony format see
https://dsd-guide.com/sites/default/files/white-papers/DSFFileFormatSpec_E.pdf
)
Supported file types.
type
Format
¶
type Format
string
Format is an enumeration of metadata types supported by this package.
const (
UnknownFormat
Format
= ""
// Unknown Format.
ID3v1
Format
= "ID3v1"
// ID3v1 tag format.
ID3v2_2
Format
= "ID3v2.2"
// ID3v2.2 tag format.
ID3v2_3
Format
= "ID3v2.3"
// ID3v2.3 tag format (most common).
ID3v2_4
Format
= "ID3v2.4"
// ID3v2.4 tag format.
MP4
Format
= "MP4"
// MP4 tag (atom) format (see
http://www.ftyps.com/
for a full file type list)
VORBIS
Format
= "VORBIS"
// Vorbis Comment tag format.
)
Supported tag formats.
type
Metadata
¶
type Metadata interface {
// Format returns the metadata Format used to encode the data.
Format()
Format
// FileType returns the file type of the audio file.
FileType()
FileType
// Title returns the title of the track.
Title()
string
// Album returns the album name of the track.
Album()
string
// Artist returns the artist name of the track.
Artist()
string
// AlbumArtist returns the album artist name of the track.
AlbumArtist()
string
// Composer returns the composer of the track.
Composer()
string
// Year returns the year of the track.
Year()
int
// Genre returns the genre of the track.
Genre()
string
// Track returns the track number and total tracks, or zero values if unavailable.
Track() (
int
,
int
)
// Disc returns the disc number and total discs, or zero values if unavailable.
Disc() (
int
,
int
)
// Picture returns a picture, or nil if not available.
Picture() *
Picture
// Lyrics returns the lyrics, or an empty string if unavailable.
Lyrics()
string
// Comment returns the comment, or an empty string if unavailable.
Comment()
string
// Raw returns the raw mapping of retrieved tag names and associated values.
// NB: tag/atom names are not standardised between formats.
Raw() map[
string
]interface{}
}
Metadata is an interface which is used to describe metadata retrieved by this package.
func
ReadAtoms
¶
func ReadAtoms(r
io
.
ReadSeeker
) (
Metadata
,
error
)
ReadAtoms reads MP4 metadata atoms from the io.ReadSeeker into a Metadata, returning
non-nil error if there was a problem.
func
ReadDSFTags
¶
func ReadDSFTags(r
io
.
ReadSeeker
) (
Metadata
,
error
)
ReadDSFTags reads DSF metadata from the io.ReadSeeker, returning the resulting
metadata in a Metadata implementation, or non-nil error if there was a problem.
samples:
http://www.2l.no/hires/index.html
func
ReadFLACTags
¶
func ReadFLACTags(r
io
.
ReadSeeker
) (
Metadata
,
error
)
ReadFLACTags reads FLAC metadata from the io.ReadSeeker, returning the resulting
metadata in a Metadata implementation, or non-nil error if there was a problem.
func
ReadFrom
¶
func ReadFrom(r
io
.
ReadSeeker
) (
Metadata
,
error
)
ReadFrom detects and parses audio file metadata tags (currently supports ID3v1,2.{2,3,4}, MP4, FLAC/OGG).
Returns non-nil error if the format of the given data could not be determined, or if there was a problem
parsing the data.
func
ReadID3v1Tags
¶
func ReadID3v1Tags(r
io
.
ReadSeeker
) (
Metadata
,
error
)
ReadID3v1Tags reads ID3v1 tags from the io.ReadSeeker.  Returns ErrNotID3v1
if there are no ID3v1 tags, otherwise non-nil error if there was a problem.
func
ReadID3v2Tags
¶
func ReadID3v2Tags(r
io
.
ReadSeeker
) (
Metadata
,
error
)
ReadID3v2Tags parses ID3v2.{2,3,4} tags from the io.ReadSeeker into a Metadata, returning
non-nil error on failure.
func
ReadOGGTags
¶
func ReadOGGTags(r
io
.
Reader
) (
Metadata
,
error
)
ReadOGGTags reads OGG metadata from the io.ReadSeeker, returning the resulting
metadata in a Metadata implementation, or non-nil error if there was a problem.
See
http://www.xiph.org/vorbis/doc/Vorbis_I_spec.html
and
http://www.xiph.org/ogg/doc/framing.html
for details.
For Opus see
https://tools.ietf.org/html/rfc7845
type
Picture
¶
type Picture struct {
Ext
string
// Extension of the picture file.
MIMEType
string
// MIMEType of the picture.
Type
string
// Type of the picture (see pictureTypes).
Description
string
// Description.
Data        []
byte
// Raw picture data.
}
Picture is a type which represents an attached picture extracted from metadata.
func (Picture)
String
¶
func (p
Picture
) String()
string
String returns a string representation of the underlying Picture instance.
type
UFID
¶
type UFID struct {
Provider
string
Identifier []
byte
}
UFID is composed of a provider (frequently a URL and a binary identifier)
The identifier can be a text (Musicbrainz use texts, but not necessary)
func (UFID)
String
¶
func (u
UFID
) String()
string