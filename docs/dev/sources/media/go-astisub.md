# go-astisub (subtitles)

> Auto-fetched from [https://pkg.go.dev/github.com/asticode/go-astisub](https://pkg.go.dev/github.com/asticode/go-astisub)
> Last Updated: 2026-01-29T20:14:46.305147+00:00

---

Index
¶
Constants
Variables
type Color
func (c *Color) SSAString() string
func (c *Color) TTMLString() string
type Item
func (i Item) String() string
type Justification
type Line
func (l Line) String() string
type LineItem
func (li LineItem) STLString() string
type Metadata
type Options
type Region
type SSAOptions
type STLOptions
type STLPosition
type Style
type StyleAttributes
type Subtitles
func NewSubtitles() *Subtitles
func Open(o Options) (s *Subtitles, err error)
func OpenFile(filename string) (*Subtitles, error)
func ReadFromSRT(i io.Reader) (o *Subtitles, err error)
func ReadFromSSA(i io.Reader) (o *Subtitles, err error)
func ReadFromSSAWithOptions(i io.Reader, opts SSAOptions) (o *Subtitles, err error)
func ReadFromSTL(i io.Reader, opts STLOptions) (o *Subtitles, err error)
func ReadFromTTML(i io.Reader) (o *Subtitles, err error)
func ReadFromTeletext(r io.Reader, o TeletextOptions) (s *Subtitles, err error)
func ReadFromWebVTT(i io.Reader) (o *Subtitles, err error)
func (s *Subtitles) Add(d time.Duration)
func (s *Subtitles) ApplyLinearCorrection(actual1, desired1, actual2, desired2 time.Duration)
func (s Subtitles) Duration() time.Duration
func (s *Subtitles) ForceDuration(d time.Duration, addDummyItem bool)
func (s *Subtitles) Fragment(f time.Duration)
func (s Subtitles) IsEmpty() bool
func (s *Subtitles) Merge(i *Subtitles)
func (s *Subtitles) Optimize()
func (s *Subtitles) Order()
func (s *Subtitles) RemoveStyling()
func (s *Subtitles) Unfragment()
func (s Subtitles) Write(dst string) (err error)
func (s Subtitles) WriteToSRT(o io.Writer) (err error)
func (s Subtitles) WriteToSSA(o io.Writer) (err error)
func (s Subtitles) WriteToSTL(o io.Writer) (err error)
func (s Subtitles) WriteToTTML(o io.Writer, opts ...WriteToTTMLOption) (err error)
func (s Subtitles) WriteToWebVTT(o io.Writer) (err error)
type TTMLIn
type TTMLInBody
type TTMLInBodyDiv
type TTMLInDuration
func (d *TTMLInDuration) UnmarshalText(i []byte) (err error)
type TTMLInHeader
type TTMLInItem
type TTMLInItems
func (i *TTMLInItems) UnmarshalXML(d *xml.Decoder, start xml.StartElement) (err error)
type TTMLInMetadata
type TTMLInRegion
type TTMLInStyle
type TTMLInStyleAttributes
type TTMLInSubtitle
type TTMLOut
type TTMLOutDuration
func (t TTMLOutDuration) MarshalText() ([]byte, error)
type TTMLOutHeader
type TTMLOutItem
type TTMLOutMetadata
type TTMLOutRegion
type TTMLOutStyle
type TTMLOutStyleAttributes
type TTMLOutSubtitle
type TeletextOptions
type WebVTTPosition
func (p *WebVTTPosition) String() string
type WebVTTTag
type WebVTTTimestampMap
func (t *WebVTTTimestampMap) Offset() time.Duration
func (t *WebVTTTimestampMap) String() string
type WriteToTTMLOption
func WriteToTTMLWithIndentOption(indent string) WriteToTTMLOption
type WriteToTTMLOptions
Constants
¶
View Source
const (
LanguageChinese   = "chinese"
LanguageEnglish   = "english"
LanguageFrench    = "french"
LanguageJapanese  = "japanese"
LanguageNorwegian = "norwegian"
)
Languages
Variables
¶
View Source
var (
ColorBlack   = &
Color
{}
ColorBlue    = &
Color
{Blue: 255}
ColorCyan    = &
Color
{Blue: 255, Green: 255}
ColorGray    = &
Color
{Blue: 128, Green: 128, Red: 128}
ColorGreen   = &
Color
{Green: 128}
ColorLime    = &
Color
{Green: 255}
ColorMagenta = &
Color
{Blue: 255, Red: 255}
ColorMaroon  = &
Color
{Red: 128}
ColorNavy    = &
Color
{Blue: 128}
ColorOlive   = &
Color
{Green: 128, Red: 128}
ColorPurple  = &
Color
{Blue: 128, Red: 128}
ColorRed     = &
Color
{Red: 255}
ColorSilver  = &
Color
{Blue: 192, Green: 192, Red: 192}
ColorTeal    = &
Color
{Blue: 128, Green: 128}
ColorYellow  = &
Color
{Green: 255, Red: 255}
ColorWhite   = &
Color
{Blue: 255, Green: 255, Red: 255}
)
Colors
View Source
var (
ErrInvalidExtension   =
errors
.
New
("astisub: invalid extension")
ErrNoSubtitlesToWrite =
errors
.
New
("astisub: no subtitles to write")
)
Errors
View Source
var (
JustificationUnchanged =
Justification
(1)
JustificationLeft      =
Justification
(2)
JustificationCentered  =
Justification
(3)
JustificationRight     =
Justification
(4)
)
View Source
var (
BytesBOM = []
byte
{239, 187, 191}
)
Bytes
View Source
var (
ErrNoValidTeletextPID =
errors
.
New
("astisub: no valid teletext PID")
)
Errors
View Source
var Now = func()
time
.
Time
{
return
time
.
Now
()
}
Now allows testing functions using it
Functions
¶
This section is empty.
Types
¶
type
Color
¶
type Color struct {
Alpha, Blue, Green, Red
uint8
}
Color represents a color
func (*Color)
SSAString
¶
func (c *
Color
) SSAString()
string
SSAString expresses the color as an SSA string
func (*Color)
TTMLString
¶
func (c *
Color
) TTMLString()
string
TTMLString expresses the color as a TTML string
type
Item
¶
type Item struct {
Comments    []
string
Index
int
EndAt
time
.
Duration
InlineStyle *
StyleAttributes
Lines       []
Line
Region      *
Region
StartAt
time
.
Duration
Style       *
Style
}
Item represents a text to show between 2 time boundaries with formatting
func (Item)
String
¶
func (i
Item
) String()
string
String implements the Stringer interface
type
Justification
¶
added in
v0.12.0
type Justification
int
type
Line
¶
type Line struct {
Items     []
LineItem
VoiceName
string
}
Line represents a set of formatted line items
func (Line)
String
¶
func (l
Line
) String()
string
String implement the Stringer interface
type
LineItem
¶
type LineItem struct {
InlineStyle *
StyleAttributes
StartAt
time
.
Duration
Style       *
Style
Text
string
}
LineItem represents a formatted line item
func (LineItem)
STLString
¶
added in
v0.12.0
func (li
LineItem
) STLString()
string
type
Metadata
¶
type Metadata struct {
Comments                                            []
string
Framerate
int
Language
string
SSACollisions
string
SSAOriginalEditing
string
SSAOriginalScript
string
SSAOriginalTiming
string
SSAOriginalTranslation
string
SSAPlayDepth                                        *
int
SSAPlayResX, SSAPlayResY                            *
int
SSAScriptType
string
SSAScriptUpdatedBy
string
SSASynchPoint
string
SSATimer                                            *
float64
SSAUpdateDetails
string
SSAWrapStyle
string
STLCountryOfOrigin
string
STLCreationDate                                     *
time
.
Time
STLDisplayStandardCode
string
STLEditorContactDetails
string
STLEditorName
string
STLMaximumNumberOfDisplayableCharactersInAnyTextRow *
int
STLMaximumNumberOfDisplayableRows                   *
int
STLOriginalEpisodeTitle
string
STLPublisher
string
STLRevisionDate                                     *
time
.
Time
STLRevisionNumber
int
STLSubtitleListReferenceCode
string
STLTimecodeStartOfProgramme
time
.
Duration
STLTranslatedEpisodeTitle
string
STLTranslatedProgramTitle
string
STLTranslatorContactDetails
string
STLTranslatorName
string
Title
string
TTMLCopyright
string
WebVTTTimestampMap                                  *
WebVTTTimestampMap
}
Metadata represents metadata
TODO Merge attributes
type
Options
¶
type Options struct {
Filename
string
Teletext
TeletextOptions
STL
STLOptions
}
Options represents open or write options
type
Region
¶
type Region struct {
ID
string
InlineStyle *
StyleAttributes
Style       *
Style
}
Region represents a subtitle's region
type
SSAOptions
¶
added in
v0.20.0
type SSAOptions struct {
OnUnknownSectionName func(name
string
)
OnInvalidLine        func(line
string
)
}
SSAOptions
type
STLOptions
¶
added in
v0.15.0
type STLOptions struct {
// IgnoreTimecodeStartOfProgramme - set STLTimecodeStartOfProgramme to zero before parsing
IgnoreTimecodeStartOfProgramme
bool
}
STLOptions represents STL parsing options
type
STLPosition
¶
added in
v0.12.0
type STLPosition struct {
VerticalPosition
int
MaxRows
int
Rows
int
}
type
Style
¶
type Style struct {
ID
string
InlineStyle *
StyleAttributes
Style       *
Style
}
Style represents a subtitle's style
type
StyleAttributes
¶
type StyleAttributes struct {
SRTBold
bool
SRTColor             *
string
SRTItalics
bool
SRTPosition
byte
// 1-9 numpad layout
SRTUnderline
bool
SSAAlignment         *
int
SSAAlphaLevel        *
float64
SSAAngle             *
float64
// degrees
SSABackColour        *
Color
SSABold              *
bool
SSABorderStyle       *
int
SSAEffect
string
SSAEncoding          *
int
SSAFontName
string
SSAFontSize          *
float64
SSAItalic            *
bool
SSALayer             *
int
SSAMarginLeft        *
int
// pixels
SSAMarginRight       *
int
// pixels
SSAMarginVertical    *
int
// pixels
SSAMarked            *
bool
SSAOutline           *
float64
// pixels
SSAOutlineColour     *
Color
SSAPrimaryColour     *
Color
SSAScaleX            *
float64
// %
SSAScaleY            *
float64
// %
SSASecondaryColour   *
Color
SSAShadow            *
float64
// pixels
SSASpacing           *
float64
// pixels
SSAStrikeout         *
bool
SSAUnderline         *
bool
STLBoxing            *
bool
STLItalics           *
bool
STLJustification     *
Justification
STLPosition          *
STLPosition
STLUnderline         *
bool
TeletextColor        *
Color
TeletextDoubleHeight *
bool
TeletextDoubleSize   *
bool
TeletextDoubleWidth  *
bool
TeletextSpacesAfter  *
int
TeletextSpacesBefore *
int
// TODO Use pointers with real types below
TTMLBackgroundColor  *
string
//
https://htmlcolorcodes.com/fr/
TTMLColor            *
string
TTMLDirection        *
string
TTMLDisplay          *
string
TTMLDisplayAlign     *
string
TTMLExtent           *
string
TTMLFontFamily       *
string
TTMLFontSize         *
string
TTMLFontStyle        *
string
TTMLFontWeight       *
string
TTMLLineHeight       *
string
TTMLOpacity          *
string
TTMLOrigin           *
string
TTMLOverflow         *
string
TTMLPadding          *
string
TTMLShowBackground   *
string
TTMLTextAlign        *
string
TTMLTextDecoration   *
string
TTMLTextOutline      *
string
TTMLUnicodeBidi      *
string
TTMLVisibility       *
string
TTMLWrapOption       *
string
TTMLWritingMode      *
string
TTMLZIndex           *
int
WebVTTAlign
string
WebVTTBold
bool
WebVTTItalics
bool
WebVTTLine
string
WebVTTLines
int
WebVTTPosition       *
WebVTTPosition
WebVTTRegionAnchor
string
WebVTTScroll
string
WebVTTSize
string
WebVTTStyles         []
string
WebVTTTags           []
WebVTTTag
WebVTTUnderline
bool
WebVTTVertical
string
WebVTTViewportAnchor
string
WebVTTWidth
string
}
StyleAttributes represents style attributes
type
Subtitles
¶
type Subtitles struct {
Items    []*
Item
Metadata *
Metadata
Regions  map[
string
]*
Region
Styles   map[
string
]*
Style
}
Subtitles represents an ordered list of items with formatting
func
NewSubtitles
¶
func NewSubtitles() *
Subtitles
NewSubtitles creates new subtitles
func
Open
¶
func Open(o
Options
) (s *
Subtitles
, err
error
)
Open opens a subtitle reader based on options
func
OpenFile
¶
func OpenFile(filename
string
) (*
Subtitles
,
error
)
OpenFile opens a file regardless of other options
func
ReadFromSRT
¶
func ReadFromSRT(i
io
.
Reader
) (o *
Subtitles
, err
error
)
ReadFromSRT parses an .srt content
func
ReadFromSSA
¶
func ReadFromSSA(i
io
.
Reader
) (o *
Subtitles
, err
error
)
ReadFromSSA parses an .ssa content
func
ReadFromSSAWithOptions
¶
added in
v0.20.0
func ReadFromSSAWithOptions(i
io
.
Reader
, opts
SSAOptions
) (o *
Subtitles
, err
error
)
ReadFromSSAWithOptions parses an .ssa content
func
ReadFromSTL
¶
func ReadFromSTL(i
io
.
Reader
, opts
STLOptions
) (o *
Subtitles
, err
error
)
ReadFromSTL parses an .stl content
func
ReadFromTTML
¶
func ReadFromTTML(i
io
.
Reader
) (o *
Subtitles
, err
error
)
ReadFromTTML parses a .ttml content
func
ReadFromTeletext
¶
func ReadFromTeletext(r
io
.
Reader
, o
TeletextOptions
) (s *
Subtitles
, err
error
)
ReadFromTeletext parses a teletext content
http://www.etsi.org/deliver/etsi_en/300400_300499/300472/01.03.01_60/en_300472v010301p.pdf
http://www.etsi.org/deliver/etsi_i_ets/300700_300799/300706/01_60/ets_300706e01p.pdf
TODO Update README
TODO Add tests
func
ReadFromWebVTT
¶
func ReadFromWebVTT(i
io
.
Reader
) (o *
Subtitles
, err
error
)
ReadFromWebVTT parses a .vtt content
TODO Tags (u, i, b)
TODO Class
func (*Subtitles)
Add
¶
func (s *
Subtitles
) Add(d
time
.
Duration
)
Add adds a duration to each time boundaries. As in the time package, duration can be negative.
func (*Subtitles)
ApplyLinearCorrection
¶
added in
v0.21.0
func (s *
Subtitles
) ApplyLinearCorrection(actual1, desired1, actual2, desired2
time
.
Duration
)
ApplyLinearCorrection applies linear correction
func (Subtitles)
Duration
¶
func (s
Subtitles
) Duration()
time
.
Duration
Duration returns the subtitles duration
func (*Subtitles)
ForceDuration
¶
func (s *
Subtitles
) ForceDuration(d
time
.
Duration
, addDummyItem
bool
)
ForceDuration updates the subtitles duration.
If requested duration is bigger, then we create a dummy item.
If requested duration is smaller, then we remove useless items and we cut the last item or add a dummy item.
func (*Subtitles)
Fragment
¶
func (s *
Subtitles
) Fragment(f
time
.
Duration
)
Fragment fragments subtitles with a specific fragment duration
func (Subtitles)
IsEmpty
¶
func (s
Subtitles
) IsEmpty()
bool
IsEmpty returns whether the subtitles are empty
func (*Subtitles)
Merge
¶
func (s *
Subtitles
) Merge(i *
Subtitles
)
Merge merges subtitles i into subtitles
func (*Subtitles)
Optimize
¶
func (s *
Subtitles
) Optimize()
Optimize optimizes subtitles
func (*Subtitles)
Order
¶
func (s *
Subtitles
) Order()
Order orders items
func (*Subtitles)
RemoveStyling
¶
func (s *
Subtitles
) RemoveStyling()
RemoveStyling removes the styling from the subtitles
func (*Subtitles)
Unfragment
¶
func (s *
Subtitles
) Unfragment()
Unfragment unfragments subtitles
func (Subtitles)
Write
¶
func (s
Subtitles
) Write(dst
string
) (err
error
)
Write writes subtitles to a file
func (Subtitles)
WriteToSRT
¶
func (s
Subtitles
) WriteToSRT(o
io
.
Writer
) (err
error
)
WriteToSRT writes subtitles in .srt format
func (Subtitles)
WriteToSSA
¶
func (s
Subtitles
) WriteToSSA(o
io
.
Writer
) (err
error
)
WriteToSSA writes subtitles in .ssa format
func (Subtitles)
WriteToSTL
¶
func (s
Subtitles
) WriteToSTL(o
io
.
Writer
) (err
error
)
WriteToSTL writes subtitles in .stl format
func (Subtitles)
WriteToTTML
¶
func (s
Subtitles
) WriteToTTML(o
io
.
Writer
, opts ...
WriteToTTMLOption
) (err
error
)
WriteToTTML writes subtitles in .ttml format
func (Subtitles)
WriteToWebVTT
¶
func (s
Subtitles
) WriteToWebVTT(o
io
.
Writer
) (err
error
)
WriteToWebVTT writes subtitles in .vtt format
type
TTMLIn
¶
type TTMLIn struct {
Framerate
int
`xml:"frameRate,attr"`
Lang
string
`xml:"lang,attr"`
Metadata
TTMLInMetadata
`xml:"head>metadata"`
Regions   []
TTMLInRegion
`xml:"head>layout>region"`
Styles    []
TTMLInStyle
`xml:"head>styling>style"`
Body
TTMLInBody
`xml:"body"`
Tickrate
int
`xml:"tickRate,attr"`
XMLName
xml
.
Name
`xml:"tt"`
}
TTMLIn represents an input TTML that must be unmarshaled
We split it from the output TTML as we can't add strict namespace without breaking retrocompatibility
type
TTMLInBody
¶
added in
v0.36.0
type TTMLInBody struct {
XMLName
xml
.
Name
`xml:"body"`
Divs    []
TTMLInBodyDiv
`xml:"div"`
Region
string
`xml:"region,attr,omitempty"`
Style
string
`xml:"style,attr,omitempty"`
TTMLInStyleAttributes
}
type
TTMLInBodyDiv
¶
added in
v0.36.0
type TTMLInBodyDiv struct {
XMLName
xml
.
Name
`xml:"div"`
Subtitles []
TTMLInSubtitle
`xml:"p"`
Region
string
`xml:"region,attr,omitempty"`
Style
string
`xml:"style,attr,omitempty"`
TTMLInStyleAttributes
}
type
TTMLInDuration
¶
type TTMLInDuration struct {
// contains filtered or unexported fields
}
TTMLInDuration represents an input TTML duration
func (*TTMLInDuration)
UnmarshalText
¶
func (d *
TTMLInDuration
) UnmarshalText(i []
byte
) (err
error
)
UnmarshalText implements the TextUnmarshaler interface
Possible formats are:
- hh:mm:ss.mmm
- hh:mm:ss:fff (fff being frames)
- [ticks]t ([ticks] being the tick amount)
type
TTMLInHeader
¶
type TTMLInHeader struct {
ID
string
`xml:"id,attr,omitempty"`
Style
string
`xml:"style,attr,omitempty"`
TTMLInStyleAttributes
}
TTMLInHeader represents an input TTML header
type
TTMLInItem
¶
type TTMLInItem struct {
Style
string
`xml:"style,attr,omitempty"`
Text
string
`xml:",chardata"`
TTMLInStyleAttributes
XMLName
xml
.
Name
}
TTMLInItem represents an input TTML item
type
TTMLInItems
¶
type TTMLInItems []
TTMLInItem
TTMLInItems represents input TTML items
func (*TTMLInItems)
UnmarshalXML
¶
func (i *
TTMLInItems
) UnmarshalXML(d *
xml
.
Decoder
, start
xml
.
StartElement
) (err
error
)
UnmarshalXML implements the XML unmarshaler interface
type
TTMLInMetadata
¶
type TTMLInMetadata struct {
Copyright
string
`xml:"copyright"`
Title
string
`xml:"title"`
}
TTMLInMetadata represents an input TTML Metadata
type
TTMLInRegion
¶
type TTMLInRegion struct {
TTMLInHeader
XMLName
xml
.
Name
`xml:"region"`
}
TTMLInRegion represents an input TTML region
type
TTMLInStyle
¶
type TTMLInStyle struct {
TTMLInHeader
XMLName
xml
.
Name
`xml:"style"`
}
TTMLInStyle represents an input TTML style
type
TTMLInStyleAttributes
¶
type TTMLInStyleAttributes struct {
BackgroundColor *
string
`xml:"backgroundColor,attr,omitempty"`
Color           *
string
`xml:"color,attr,omitempty"`
Direction       *
string
`xml:"direction,attr,omitempty"`
Display         *
string
`xml:"display,attr,omitempty"`
DisplayAlign    *
string
`xml:"displayAlign,attr,omitempty"`
Extent          *
string
`xml:"extent,attr,omitempty"`
FontFamily      *
string
`xml:"fontFamily,attr,omitempty"`
FontSize        *
string
`xml:"fontSize,attr,omitempty"`
FontStyle       *
string
`xml:"fontStyle,attr,omitempty"`
FontWeight      *
string
`xml:"fontWeight,attr,omitempty"`
LineHeight      *
string
`xml:"lineHeight,attr,omitempty"`
Opacity         *
string
`xml:"opacity,attr,omitempty"`
Origin          *
string
`xml:"origin,attr,omitempty"`
Overflow        *
string
`xml:"overflow,attr,omitempty"`
Padding         *
string
`xml:"padding,attr,omitempty"`
ShowBackground  *
string
`xml:"showBackground,attr,omitempty"`
TextAlign       *
string
`xml:"textAlign,attr,omitempty"`
TextDecoration  *
string
`xml:"textDecoration,attr,omitempty"`
TextOutline     *
string
`xml:"textOutline,attr,omitempty"`
UnicodeBidi     *
string
`xml:"unicodeBidi,attr,omitempty"`
Visibility      *
string
`xml:"visibility,attr,omitempty"`
WrapOption      *
string
`xml:"wrapOption,attr,omitempty"`
WritingMode     *
string
`xml:"writingMode,attr,omitempty"`
ZIndex          *
int
`xml:"zIndex,attr,omitempty"`
}
TTMLInStyleAttributes represents input TTML style attributes
type
TTMLInSubtitle
¶
type TTMLInSubtitle struct {
Begin *
TTMLInDuration
`xml:"begin,attr,omitempty"`
End   *
TTMLInDuration
`xml:"end,attr,omitempty"`
ID
string
`xml:"id,attr,omitempty"`
// We must store inner XML temporarily here since there's no tag to describe both any tag and chardata
// Real unmarshal will be done manually afterwards
Items
string
`xml:",innerxml"`
Region
string
`xml:"region,attr,omitempty"`
Style
string
`xml:"style,attr,omitempty"`
TTMLInStyleAttributes
}
TTMLInSubtitle represents an input TTML subtitle
type
TTMLOut
¶
type TTMLOut struct {
Lang
string
`xml:"xml:lang,attr,omitempty"`
Metadata        *
TTMLOutMetadata
`xml:"head>metadata,omitempty"`
Styles          []
TTMLOutStyle
`xml:"head>styling>style,omitempty"`
//!\\ Order is important! Keep Styling above Layout
Regions         []
TTMLOutRegion
`xml:"head>layout>region,omitempty"`
Subtitles       []
TTMLOutSubtitle
`xml:"body>div>p,omitempty"`
XMLName
xml
.
Name
`xml:"http://www.w3.org/ns/ttml tt"`
XMLNamespaceTTM
string
`xml:"xmlns:ttm,attr"`
XMLNamespaceTTS
string
`xml:"xmlns:tts,attr"`
}
TTMLOut represents an output TTML that must be marshaled
We split it from the input TTML as this time we'll add strict namespaces
type
TTMLOutDuration
¶
type TTMLOutDuration
time
.
Duration
TTMLOutDuration represents an output TTML duration
func (TTMLOutDuration)
MarshalText
¶
func (t
TTMLOutDuration
) MarshalText() ([]
byte
,
error
)
MarshalText implements the TextMarshaler interface
type
TTMLOutHeader
¶
type TTMLOutHeader struct {
ID
string
`xml:"xml:id,attr,omitempty"`
Style
string
`xml:"style,attr,omitempty"`
TTMLOutStyleAttributes
}
TTMLOutHeader represents an output TTML header
type
TTMLOutItem
¶
type TTMLOutItem struct {
Style
string
`xml:"style,attr,omitempty"`
Text
string
`xml:",chardata"`
TTMLOutStyleAttributes
XMLName
xml
.
Name
}
TTMLOutItem represents an output TTML Item
type
TTMLOutMetadata
¶
type TTMLOutMetadata struct {
Copyright
string
`xml:"ttm:copyright,omitempty"`
Title
string
`xml:"ttm:title,omitempty"`
}
TTMLOutMetadata represents an output TTML Metadata
type
TTMLOutRegion
¶
type TTMLOutRegion struct {
TTMLOutHeader
XMLName
xml
.
Name
`xml:"region"`
}
TTMLOutRegion represents an output TTML region
type
TTMLOutStyle
¶
type TTMLOutStyle struct {
TTMLOutHeader
XMLName
xml
.
Name
`xml:"style"`
}
TTMLOutStyle represents an output TTML style
type
TTMLOutStyleAttributes
¶
type TTMLOutStyleAttributes struct {
BackgroundColor *
string
`xml:"tts:backgroundColor,attr,omitempty"`
Color           *
string
`xml:"tts:color,attr,omitempty"`
Direction       *
string
`xml:"tts:direction,attr,omitempty"`
Display         *
string
`xml:"tts:display,attr,omitempty"`
DisplayAlign    *
string
`xml:"tts:displayAlign,attr,omitempty"`
Extent          *
string
`xml:"tts:extent,attr,omitempty"`
FontFamily      *
string
`xml:"tts:fontFamily,attr,omitempty"`
FontSize        *
string
`xml:"tts:fontSize,attr,omitempty"`
FontStyle       *
string
`xml:"tts:fontStyle,attr,omitempty"`
FontWeight      *
string
`xml:"tts:fontWeight,attr,omitempty"`
LineHeight      *
string
`xml:"tts:lineHeight,attr,omitempty"`
Opacity         *
string
`xml:"tts:opacity,attr,omitempty"`
Origin          *
string
`xml:"tts:origin,attr,omitempty"`
Overflow        *
string
`xml:"tts:overflow,attr,omitempty"`
Padding         *
string
`xml:"tts:padding,attr,omitempty"`
ShowBackground  *
string
`xml:"tts:showBackground,attr,omitempty"`
TextAlign       *
string
`xml:"tts:textAlign,attr,omitempty"`
TextDecoration  *
string
`xml:"tts:textDecoration,attr,omitempty"`
TextOutline     *
string
`xml:"tts:textOutline,attr,omitempty"`
UnicodeBidi     *
string
`xml:"tts:unicodeBidi,attr,omitempty"`
Visibility      *
string
`xml:"tts:visibility,attr,omitempty"`
WrapOption      *
string
`xml:"tts:wrapOption,attr,omitempty"`
WritingMode     *
string
`xml:"tts:writingMode,attr,omitempty"`
ZIndex          *
int
`xml:"tts:zIndex,attr,omitempty"`
}
TTMLOutStyleAttributes represents output TTML style attributes
type
TTMLOutSubtitle
¶
type TTMLOutSubtitle struct {
Begin
TTMLOutDuration
`xml:"begin,attr"`
End
TTMLOutDuration
`xml:"end,attr"`
ID
string
`xml:"id,attr,omitempty"`
Items  []
TTMLOutItem
Region
string
`xml:"region,attr,omitempty"`
Style
string
`xml:"style,attr,omitempty"`
TTMLOutStyleAttributes
}
TTMLOutSubtitle represents an output TTML subtitle
type
TeletextOptions
¶
type TeletextOptions struct {
Page
int
PID
int
}
TeletextOptions represents teletext options
type
WebVTTPosition
¶
added in
v0.38.0
type WebVTTPosition struct {
XPosition
string
Alignment
string
}
func (*WebVTTPosition)
String
¶
added in
v0.38.0
func (p *
WebVTTPosition
) String()
string
type
WebVTTTag
¶
added in
v0.26.1
type WebVTTTag struct {
Name
string
Annotation
string
Classes    []
string
}
type
WebVTTTimestampMap
¶
added in
v0.27.0
type WebVTTTimestampMap struct {
Local
time
.
Duration
MpegTS
int64
}
WebVTTTimestampMap is a structure for storing timestamps for WEBVTT's
X-TIMESTAMP-MAP feature commonly used for syncing cue times with
MPEG-TS streams.
func (*WebVTTTimestampMap)
Offset
¶
added in
v0.27.0
func (t *
WebVTTTimestampMap
) Offset()
time
.
Duration
Offset calculates and returns the time offset described by the
timestamp map.
func (*WebVTTTimestampMap)
String
¶
added in
v0.27.0
func (t *
WebVTTTimestampMap
) String()
string
String implements Stringer interface for TimestampMap, returning
the fully formatted header string for the instance.
type
WriteToTTMLOption
¶
added in
v0.28.0
type WriteToTTMLOption func(o *
WriteToTTMLOptions
)
WriteToTTMLOption represents a WriteToTTML option.
func
WriteToTTMLWithIndentOption
¶
added in
v0.28.0
func WriteToTTMLWithIndentOption(indent
string
)
WriteToTTMLOption
WriteToTTMLWithIndentOption sets the indent option.
type
WriteToTTMLOptions
¶
added in
v0.28.0
type WriteToTTMLOptions struct {
Indent
string
// Default is 4 spaces.
}
WriteToTTMLOptions represents TTML write options.