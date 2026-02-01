# go-astisub (subtitles)

> Source: https://pkg.go.dev/github.com/asticode/go-astisub
> Fetched: 2026-02-01T11:48:34.588398+00:00
> Content-Hash: 9a47a11349a26190
> Type: html

---

### Index ¶

  * Constants
  * Variables
  * type Color
  *     * func (c *Color) SSAString() string
    * func (c *Color) TTMLString() string
  * type Item
  *     * func (i Item) String() string
  * type Justification
  * type Line
  *     * func (l Line) String() string
  * type LineItem
  *     * func (li LineItem) STLString() string
  * type Metadata
  * type Options
  * type Region
  * type SSAOptions
  * type STLOptions
  * type STLPosition
  * type Style
  * type StyleAttributes
  * type Subtitles
  *     * func NewSubtitles() *Subtitles
    * func Open(o Options) (s *Subtitles, err error)
    * func OpenFile(filename string) (*Subtitles, error)
    * func ReadFromSRT(i io.Reader) (o *Subtitles, err error)
    * func ReadFromSSA(i io.Reader) (o *Subtitles, err error)
    * func ReadFromSSAWithOptions(i io.Reader, opts SSAOptions) (o *Subtitles, err error)
    * func ReadFromSTL(i io.Reader, opts STLOptions) (o *Subtitles, err error)
    * func ReadFromTTML(i io.Reader) (o *Subtitles, err error)
    * func ReadFromTeletext(r io.Reader, o TeletextOptions) (s *Subtitles, err error)
    * func ReadFromWebVTT(i io.Reader) (o *Subtitles, err error)
  *     * func (s *Subtitles) Add(d time.Duration)
    * func (s *Subtitles) ApplyLinearCorrection(actual1, desired1, actual2, desired2 time.Duration)
    * func (s Subtitles) Duration() time.Duration
    * func (s *Subtitles) ForceDuration(d time.Duration, addDummyItem bool)
    * func (s *Subtitles) Fragment(f time.Duration)
    * func (s Subtitles) IsEmpty() bool
    * func (s *Subtitles) Merge(i *Subtitles)
    * func (s *Subtitles) Optimize()
    * func (s *Subtitles) Order()
    * func (s *Subtitles) RemoveStyling()
    * func (s *Subtitles) Unfragment()
    * func (s Subtitles) Write(dst string) (err error)
    * func (s Subtitles) WriteToSRT(o io.Writer) (err error)
    * func (s Subtitles) WriteToSSA(o io.Writer) (err error)
    * func (s Subtitles) WriteToSTL(o io.Writer) (err error)
    * func (s Subtitles) WriteToTTML(o io.Writer, opts ...WriteToTTMLOption) (err error)
    * func (s Subtitles) WriteToWebVTT(o io.Writer) (err error)
  * type TTMLIn
  * type TTMLInBody
  * type TTMLInBodyDiv
  * type TTMLInDuration
  *     * func (d *TTMLInDuration) UnmarshalText(i []byte) (err error)
  * type TTMLInHeader
  * type TTMLInItem
  * type TTMLInItems
  *     * func (i *TTMLInItems) UnmarshalXML(d *xml.Decoder, start xml.StartElement) (err error)
  * type TTMLInMetadata
  * type TTMLInRegion
  * type TTMLInStyle
  * type TTMLInStyleAttributes
  * type TTMLInSubtitle
  * type TTMLOut
  * type TTMLOutDuration
  *     * func (t TTMLOutDuration) MarshalText() ([]byte, error)
  * type TTMLOutHeader
  * type TTMLOutItem
  * type TTMLOutMetadata
  * type TTMLOutRegion
  * type TTMLOutStyle
  * type TTMLOutStyleAttributes
  * type TTMLOutSubtitle
  * type TeletextOptions
  * type WebVTTPosition
  *     * func (p *WebVTTPosition) String() string
  * type WebVTTTag
  * type WebVTTTimestampMap
  *     * func (t *WebVTTTimestampMap) Offset() time.Duration
    * func (t *WebVTTTimestampMap) String() string
  * type WriteToTTMLOption
  *     * func WriteToTTMLWithIndentOption(indent string) WriteToTTMLOption
  * type WriteToTTMLOptions



### Constants ¶

[View Source](https://github.com/asticode/go-astisub/blob/v0.38.0/language.go#L4)
    
    
    const (
    	LanguageChinese   = "chinese"
    	LanguageEnglish   = "english"
    	LanguageFrench    = "french"
    	LanguageJapanese  = "japanese"
    	LanguageNorwegian = "norwegian"
    )

Languages 

### Variables ¶

[View Source](https://github.com/asticode/go-astisub/blob/v0.38.0/subtitles.go#L29)
    
    
    var (
    	ColorBlack   = &Color{}
    	ColorBlue    = &Color{Blue: 255}
    	ColorCyan    = &Color{Blue: 255, Green: 255}
    	ColorGray    = &Color{Blue: 128, Green: 128, Red: 128}
    	ColorGreen   = &Color{Green: 128}
    	ColorLime    = &Color{Green: 255}
    	ColorMagenta = &Color{Blue: 255, Red: 255}
    	ColorMaroon  = &Color{Red: 128}
    	ColorNavy    = &Color{Blue: 128}
    	ColorOlive   = &Color{Green: 128, Red: 128}
    	ColorPurple  = &Color{Blue: 128, Red: 128}
    	ColorRed     = &Color{Red: 255}
    	ColorSilver  = &Color{Blue: 192, Green: 192, Red: 192}
    	ColorTeal    = &Color{Blue: 128, Green: 128}
    	ColorYellow  = &Color{Green: 255, Red: 255}
    	ColorWhite   = &Color{Blue: 255, Green: 255, Red: 255}
    )

Colors 

[View Source](https://github.com/asticode/go-astisub/blob/v0.38.0/subtitles.go#L49)
    
    
    var (
    	ErrInvalidExtension   = [errors](/errors).[New](/errors#New)("astisub: invalid extension")
    	ErrNoSubtitlesToWrite = [errors](/errors).[New](/errors#New)("astisub: no subtitles to write")
    )

Errors 

[View Source](https://github.com/asticode/go-astisub/blob/v0.38.0/subtitles.go#L177)
    
    
    var (
    	JustificationUnchanged = Justification(1)
    	JustificationLeft      = Justification(2)
    	JustificationCentered  = Justification(3)
    	JustificationRight     = Justification(4)
    )

[View Source](https://github.com/asticode/go-astisub/blob/v0.38.0/subtitles.go#L22)
    
    
    var (
    	BytesBOM = [][byte](/builtin#byte){239, 187, 191}
    )

Bytes 

[View Source](https://github.com/asticode/go-astisub/blob/v0.38.0/teletext.go#L19)
    
    
    var (
    	ErrNoValidTeletextPID = [errors](/errors).[New](/errors#New)("astisub: no valid teletext PID")
    )

Errors 

[View Source](https://github.com/asticode/go-astisub/blob/v0.38.0/subtitles.go#L61)
    
    
    var Now = func() [time](/time).[Time](/time#Time) {
    	return [time](/time).[Now](/time#Now)()
    }

Now allows testing functions using it 

### Functions ¶

This section is empty.

### Types ¶

####  type [Color](https://github.com/asticode/go-astisub/blob/v0.38.0/subtitles.go#L145) ¶
    
    
    type Color struct {
    	Alpha, Blue, Green, Red [uint8](/builtin#uint8)
    }

Color represents a color 

####  func (*Color) [SSAString](https://github.com/asticode/go-astisub/blob/v0.38.0/subtitles.go#L166) ¶
    
    
    func (c *Color) SSAString() [string](/builtin#string)

SSAString expresses the color as an SSA string 

####  func (*Color) [TTMLString](https://github.com/asticode/go-astisub/blob/v0.38.0/subtitles.go#L171) ¶
    
    
    func (c *Color) TTMLString() [string](/builtin#string)

TTMLString expresses the color as a TTML string 

####  type [Item](https://github.com/asticode/go-astisub/blob/v0.38.0/subtitles.go#L124) ¶
    
    
    type Item struct {
    	Comments    [][string](/builtin#string)
    	Index       [int](/builtin#int)
    	EndAt       [time](/time).[Duration](/time#Duration)
    	InlineStyle *StyleAttributes
    	Lines       []Line
    	Region      *Region
    	StartAt     [time](/time).[Duration](/time#Duration)
    	Style       *Style
    }

Item represents a text to show between 2 time boundaries with formatting 

####  func (Item) [String](https://github.com/asticode/go-astisub/blob/v0.38.0/subtitles.go#L136) ¶
    
    
    func (i Item) String() [string](/builtin#string)

String implements the Stringer interface 

####  type [Justification](https://github.com/asticode/go-astisub/blob/v0.38.0/subtitles.go#L175) ¶ added in v0.12.0
    
    
    type Justification [int](/builtin#int)

####  type [Line](https://github.com/asticode/go-astisub/blob/v0.38.0/subtitles.go#L659) ¶
    
    
    type Line struct {
    	Items     []LineItem
    	VoiceName [string](/builtin#string)
    }

Line represents a set of formatted line items 

####  func (Line) [String](https://github.com/asticode/go-astisub/blob/v0.38.0/subtitles.go#L665) ¶
    
    
    func (l Line) String() [string](/builtin#string)

String implement the Stringer interface 

####  type [LineItem](https://github.com/asticode/go-astisub/blob/v0.38.0/subtitles.go#L675) ¶
    
    
    type LineItem struct {
    	InlineStyle *StyleAttributes
    	StartAt     [time](/time).[Duration](/time#Duration)
    	Style       *Style
    	Text        [string](/builtin#string)
    }

LineItem represents a formatted line item 

####  func (LineItem) [STLString](https://github.com/asticode/go-astisub/blob/v0.38.0/stl.go#L720) ¶ added in v0.12.0
    
    
    func (li LineItem) STLString() [string](/builtin#string)

####  type [Metadata](https://github.com/asticode/go-astisub/blob/v0.38.0/subtitles.go#L605) ¶
    
    
    type Metadata struct {
    	Comments                                            [][string](/builtin#string)
    	Framerate                                           [int](/builtin#int)
    	Language                                            [string](/builtin#string)
    	SSACollisions                                       [string](/builtin#string)
    	SSAOriginalEditing                                  [string](/builtin#string)
    	SSAOriginalScript                                   [string](/builtin#string)
    	SSAOriginalTiming                                   [string](/builtin#string)
    	SSAOriginalTranslation                              [string](/builtin#string)
    	SSAPlayDepth                                        *[int](/builtin#int)
    	SSAPlayResX, SSAPlayResY                            *[int](/builtin#int)
    	SSAScriptType                                       [string](/builtin#string)
    	SSAScriptUpdatedBy                                  [string](/builtin#string)
    	SSASynchPoint                                       [string](/builtin#string)
    	SSATimer                                            *[float64](/builtin#float64)
    	SSAUpdateDetails                                    [string](/builtin#string)
    	SSAWrapStyle                                        [string](/builtin#string)
    	STLCountryOfOrigin                                  [string](/builtin#string)
    	STLCreationDate                                     *[time](/time).[Time](/time#Time)
    	STLDisplayStandardCode                              [string](/builtin#string)
    	STLEditorContactDetails                             [string](/builtin#string)
    	STLEditorName                                       [string](/builtin#string)
    	STLMaximumNumberOfDisplayableCharactersInAnyTextRow *[int](/builtin#int)
    	STLMaximumNumberOfDisplayableRows                   *[int](/builtin#int)
    	STLOriginalEpisodeTitle                             [string](/builtin#string)
    	STLPublisher                                        [string](/builtin#string)
    	STLRevisionDate                                     *[time](/time).[Time](/time#Time)
    	STLRevisionNumber                                   [int](/builtin#int)
    	STLSubtitleListReferenceCode                        [string](/builtin#string)
    	STLTimecodeStartOfProgramme                         [time](/time).[Duration](/time#Duration)
    	STLTranslatedEpisodeTitle                           [string](/builtin#string)
    	STLTranslatedProgramTitle                           [string](/builtin#string)
    	STLTranslatorContactDetails                         [string](/builtin#string)
    	STLTranslatorName                                   [string](/builtin#string)
    	Title                                               [string](/builtin#string)
    	TTMLCopyright                                       [string](/builtin#string)
    	WebVTTTimestampMap                                  *WebVTTTimestampMap
    }

Metadata represents metadata TODO Merge attributes 

####  type [Options](https://github.com/asticode/go-astisub/blob/v0.38.0/subtitles.go#L66) ¶
    
    
    type Options struct {
    	Filename [string](/builtin#string)
    	Teletext TeletextOptions
    	STL      STLOptions
    }

Options represents open or write options 

####  type [Region](https://github.com/asticode/go-astisub/blob/v0.38.0/subtitles.go#L645) ¶
    
    
    type Region struct {
    	ID          [string](/builtin#string)
    	InlineStyle *StyleAttributes
    	Style       *Style
    }

Region represents a subtitle's region 

####  type [SSAOptions](https://github.com/asticode/go-astisub/blob/v0.38.0/ssa.go#L1262) ¶ added in v0.20.0
    
    
    type SSAOptions struct {
    	OnUnknownSectionName func(name [string](/builtin#string))
    	OnInvalidLine        func(line [string](/builtin#string))
    }

SSAOptions 

####  type [STLOptions](https://github.com/asticode/go-astisub/blob/v0.38.0/stl.go#L180) ¶ added in v0.15.0
    
    
    type STLOptions struct {
    	// IgnoreTimecodeStartOfProgramme - set STLTimecodeStartOfProgramme to zero before parsing
    	IgnoreTimecodeStartOfProgramme [bool](/builtin#bool)
    }

STLOptions represents STL parsing options 

####  type [STLPosition](https://github.com/asticode/go-astisub/blob/v0.38.0/stl.go#L173) ¶ added in v0.12.0
    
    
    type STLPosition struct {
    	VerticalPosition [int](/builtin#int)
    	MaxRows          [int](/builtin#int)
    	Rows             [int](/builtin#int)
    }

####  type [Style](https://github.com/asticode/go-astisub/blob/v0.38.0/subtitles.go#L652) ¶
    
    
    type Style struct {
    	ID          [string](/builtin#string)
    	InlineStyle *StyleAttributes
    	Style       *Style
    }

Style represents a subtitle's style 

####  type [StyleAttributes](https://github.com/asticode/go-astisub/blob/v0.38.0/subtitles.go#L185) ¶
    
    
    type StyleAttributes struct {
    	SRTBold              [bool](/builtin#bool)
    	SRTColor             *[string](/builtin#string)
    	SRTItalics           [bool](/builtin#bool)
    	SRTPosition          [byte](/builtin#byte) // 1-9 numpad layout
    	SRTUnderline         [bool](/builtin#bool)
    	SSAAlignment         *[int](/builtin#int)
    	SSAAlphaLevel        *[float64](/builtin#float64)
    	SSAAngle             *[float64](/builtin#float64) // degrees
    	SSABackColour        *Color
    	SSABold              *[bool](/builtin#bool)
    	SSABorderStyle       *[int](/builtin#int)
    	SSAEffect            [string](/builtin#string)
    	SSAEncoding          *[int](/builtin#int)
    	SSAFontName          [string](/builtin#string)
    	SSAFontSize          *[float64](/builtin#float64)
    	SSAItalic            *[bool](/builtin#bool)
    	SSALayer             *[int](/builtin#int)
    	SSAMarginLeft        *[int](/builtin#int) // pixels
    	SSAMarginRight       *[int](/builtin#int) // pixels
    	SSAMarginVertical    *[int](/builtin#int) // pixels
    	SSAMarked            *[bool](/builtin#bool)
    	SSAOutline           *[float64](/builtin#float64) // pixels
    	SSAOutlineColour     *Color
    	SSAPrimaryColour     *Color
    	SSAScaleX            *[float64](/builtin#float64) // %
    	SSAScaleY            *[float64](/builtin#float64) // %
    	SSASecondaryColour   *Color
    	SSAShadow            *[float64](/builtin#float64) // pixels
    	SSASpacing           *[float64](/builtin#float64) // pixels
    	SSAStrikeout         *[bool](/builtin#bool)
    	SSAUnderline         *[bool](/builtin#bool)
    	STLBoxing            *[bool](/builtin#bool)
    	STLItalics           *[bool](/builtin#bool)
    	STLJustification     *Justification
    	STLPosition          *STLPosition
    	STLUnderline         *[bool](/builtin#bool)
    	TeletextColor        *Color
    	TeletextDoubleHeight *[bool](/builtin#bool)
    	TeletextDoubleSize   *[bool](/builtin#bool)
    	TeletextDoubleWidth  *[bool](/builtin#bool)
    	TeletextSpacesAfter  *[int](/builtin#int)
    	TeletextSpacesBefore *[int](/builtin#int)
    	// TODO Use pointers with real types below
    	TTMLBackgroundColor  *[string](/builtin#string) // <https://htmlcolorcodes.com/fr/>
    	TTMLColor            *[string](/builtin#string)
    	TTMLDirection        *[string](/builtin#string)
    	TTMLDisplay          *[string](/builtin#string)
    	TTMLDisplayAlign     *[string](/builtin#string)
    	TTMLExtent           *[string](/builtin#string)
    	TTMLFontFamily       *[string](/builtin#string)
    	TTMLFontSize         *[string](/builtin#string)
    	TTMLFontStyle        *[string](/builtin#string)
    	TTMLFontWeight       *[string](/builtin#string)
    	TTMLLineHeight       *[string](/builtin#string)
    	TTMLOpacity          *[string](/builtin#string)
    	TTMLOrigin           *[string](/builtin#string)
    	TTMLOverflow         *[string](/builtin#string)
    	TTMLPadding          *[string](/builtin#string)
    	TTMLShowBackground   *[string](/builtin#string)
    	TTMLTextAlign        *[string](/builtin#string)
    	TTMLTextDecoration   *[string](/builtin#string)
    	TTMLTextOutline      *[string](/builtin#string)
    	TTMLUnicodeBidi      *[string](/builtin#string)
    	TTMLVisibility       *[string](/builtin#string)
    	TTMLWrapOption       *[string](/builtin#string)
    	TTMLWritingMode      *[string](/builtin#string)
    	TTMLZIndex           *[int](/builtin#int)
    	WebVTTAlign          [string](/builtin#string)
    	WebVTTBold           [bool](/builtin#bool)
    	WebVTTItalics        [bool](/builtin#bool)
    	WebVTTLine           [string](/builtin#string)
    	WebVTTLines          [int](/builtin#int)
    	WebVTTPosition       *WebVTTPosition
    	WebVTTRegionAnchor   [string](/builtin#string)
    	WebVTTScroll         [string](/builtin#string)
    	WebVTTSize           [string](/builtin#string)
    	WebVTTStyles         [][string](/builtin#string)
    	WebVTTTags           []WebVTTTag
    	WebVTTUnderline      [bool](/builtin#bool)
    	WebVTTVertical       [string](/builtin#string)
    	WebVTTViewportAnchor [string](/builtin#string)
    	WebVTTWidth          [string](/builtin#string)
    }

StyleAttributes represents style attributes 

####  type [Subtitles](https://github.com/asticode/go-astisub/blob/v0.38.0/subtitles.go#L108) ¶
    
    
    type Subtitles struct {
    	Items    []*Item
    	Metadata *Metadata
    	Regions  map[[string](/builtin#string)]*Region
    	Styles   map[[string](/builtin#string)]*Style
    }

Subtitles represents an ordered list of items with formatting 

####  func [NewSubtitles](https://github.com/asticode/go-astisub/blob/v0.38.0/subtitles.go#L116) ¶
    
    
    func NewSubtitles() *Subtitles

NewSubtitles creates new subtitles 

####  func [Open](https://github.com/asticode/go-astisub/blob/v0.38.0/subtitles.go#L73) ¶
    
    
    func Open(o Options) (s *Subtitles, err [error](/builtin#error))

Open opens a subtitle reader based on options 

####  func [OpenFile](https://github.com/asticode/go-astisub/blob/v0.38.0/subtitles.go#L103) ¶
    
    
    func OpenFile(filename [string](/builtin#string)) (*Subtitles, [error](/builtin#error))

OpenFile opens a file regardless of other options 

####  func [ReadFromSRT](https://github.com/asticode/go-astisub/blob/v0.38.0/srt.go#L35) ¶
    
    
    func ReadFromSRT(i [io](/io).[Reader](/io#Reader)) (o *Subtitles, err [error](/builtin#error))

ReadFromSRT parses an .srt content 

####  func [ReadFromSSA](https://github.com/asticode/go-astisub/blob/v0.38.0/ssa.go#L127) ¶
    
    
    func ReadFromSSA(i [io](/io).[Reader](/io#Reader)) (o *Subtitles, err [error](/builtin#error))

ReadFromSSA parses an .ssa content 

####  func [ReadFromSSAWithOptions](https://github.com/asticode/go-astisub/blob/v0.38.0/ssa.go#L133) ¶ added in v0.20.0
    
    
    func ReadFromSSAWithOptions(i [io](/io).[Reader](/io#Reader), opts SSAOptions) (o *Subtitles, err [error](/builtin#error))

ReadFromSSAWithOptions parses an .ssa content 

####  func [ReadFromSTL](https://github.com/asticode/go-astisub/blob/v0.38.0/stl.go#L186) ¶
    
    
    func ReadFromSTL(i [io](/io).[Reader](/io#Reader), opts STLOptions) (o *Subtitles, err [error](/builtin#error))

ReadFromSTL parses an .stl content 

####  func [ReadFromTTML](https://github.com/asticode/go-astisub/blob/v0.38.0/ttml.go#L353) ¶
    
    
    func ReadFromTTML(i [io](/io).[Reader](/io#Reader)) (o *Subtitles, err [error](/builtin#error))

ReadFromTTML parses a .ttml content 

####  func [ReadFromTeletext](https://github.com/asticode/go-astisub/blob/v0.38.0/teletext.go#L335) ¶
    
    
    func ReadFromTeletext(r [io](/io).[Reader](/io#Reader), o TeletextOptions) (s *Subtitles, err [error](/builtin#error))

ReadFromTeletext parses a teletext content <http://www.etsi.org/deliver/etsi_en/300400_300499/300472/01.03.01_60/en_300472v010301p.pdf> <http://www.etsi.org/deliver/etsi_i_ets/300700_300799/300706/01_60/ets_300706e01p.pdf> TODO Update README TODO Add tests 

####  func [ReadFromWebVTT](https://github.com/asticode/go-astisub/blob/v0.38.0/webvtt.go#L150) ¶
    
    
    func ReadFromWebVTT(i [io](/io).[Reader](/io#Reader)) (o *Subtitles, err [error](/builtin#error))

ReadFromWebVTT parses a .vtt content TODO Tags (u, i, b) TODO Class 

####  func (*Subtitles) [Add](https://github.com/asticode/go-astisub/blob/v0.38.0/subtitles.go#L683) ¶
    
    
    func (s *Subtitles) Add(d [time](/time).[Duration](/time#Duration))

Add adds a duration to each time boundaries. As in the time package, duration can be negative. 

####  func (*Subtitles) [ApplyLinearCorrection](https://github.com/asticode/go-astisub/blob/v0.38.0/subtitles.go#L936) ¶ added in v0.21.0
    
    
    func (s *Subtitles) ApplyLinearCorrection(actual1, desired1, actual2, desired2 [time](/time).[Duration](/time#Duration))

ApplyLinearCorrection applies linear correction 

####  func (Subtitles) [Duration](https://github.com/asticode/go-astisub/blob/v0.38.0/subtitles.go#L697) ¶
    
    
    func (s Subtitles) Duration() [time](/time).[Duration](/time#Duration)

Duration returns the subtitles duration 

####  func (*Subtitles) [ForceDuration](https://github.com/asticode/go-astisub/blob/v0.38.0/subtitles.go#L707) ¶
    
    
    func (s *Subtitles) ForceDuration(d [time](/time).[Duration](/time#Duration), addDummyItem [bool](/builtin#bool))

ForceDuration updates the subtitles duration. If requested duration is bigger, then we create a dummy item. If requested duration is smaller, then we remove useless items and we cut the last item or add a dummy item. 

####  func (*Subtitles) [Fragment](https://github.com/asticode/go-astisub/blob/v0.38.0/subtitles.go#L740) ¶
    
    
    func (s *Subtitles) Fragment(f [time](/time).[Duration](/time#Duration))

Fragment fragments subtitles with a specific fragment duration 

####  func (Subtitles) [IsEmpty](https://github.com/asticode/go-astisub/blob/v0.38.0/subtitles.go#L795) ¶
    
    
    func (s Subtitles) IsEmpty() [bool](/builtin#bool)

IsEmpty returns whether the subtitles are empty 

####  func (*Subtitles) [Merge](https://github.com/asticode/go-astisub/blob/v0.38.0/subtitles.go#L800) ¶
    
    
    func (s *Subtitles) Merge(i *Subtitles)

Merge merges subtitles i into subtitles 

####  func (*Subtitles) [Optimize](https://github.com/asticode/go-astisub/blob/v0.38.0/subtitles.go#L821) ¶
    
    
    func (s *Subtitles) Optimize()

Optimize optimizes subtitles 

####  func (*Subtitles) [Order](https://github.com/asticode/go-astisub/blob/v0.38.0/subtitles.go#L878) ¶
    
    
    func (s *Subtitles) Order()

Order orders items 

####  func (*Subtitles) [RemoveStyling](https://github.com/asticode/go-astisub/blob/v0.38.0/subtitles.go#L891) ¶
    
    
    func (s *Subtitles) RemoveStyling()

RemoveStyling removes the styling from the subtitles 

####  func (*Subtitles) [Unfragment](https://github.com/asticode/go-astisub/blob/v0.38.0/subtitles.go#L908) ¶
    
    
    func (s *Subtitles) Unfragment()

Unfragment unfragments subtitles 

####  func (Subtitles) [Write](https://github.com/asticode/go-astisub/blob/v0.38.0/subtitles.go#L949) ¶
    
    
    func (s Subtitles) Write(dst [string](/builtin#string)) (err [error](/builtin#error))

Write writes subtitles to a file 

####  func (Subtitles) [WriteToSRT](https://github.com/asticode/go-astisub/blob/v0.38.0/srt.go#L216) ¶
    
    
    func (s Subtitles) WriteToSRT(o [io](/io).[Writer](/io#Writer)) (err [error](/builtin#error))

WriteToSRT writes subtitles in .srt format 

####  func (Subtitles) [WriteToSSA](https://github.com/asticode/go-astisub/blob/v0.38.0/ssa.go#L1169) ¶
    
    
    func (s Subtitles) WriteToSSA(o [io](/io).[Writer](/io#Writer)) (err [error](/builtin#error))

WriteToSSA writes subtitles in .ssa format 

####  func (Subtitles) [WriteToSTL](https://github.com/asticode/go-astisub/blob/v0.38.0/stl.go#L911) ¶
    
    
    func (s Subtitles) WriteToSTL(o [io](/io).[Writer](/io#Writer)) (err [error](/builtin#error))

WriteToSTL writes subtitles in .stl format 

####  func (Subtitles) [WriteToTTML](https://github.com/asticode/go-astisub/blob/v0.38.0/ttml.go#L667) ¶
    
    
    func (s Subtitles) WriteToTTML(o [io](/io).[Writer](/io#Writer), opts ...WriteToTTMLOption) (err [error](/builtin#error))

WriteToTTML writes subtitles in .ttml format 

####  func (Subtitles) [WriteToWebVTT](https://github.com/asticode/go-astisub/blob/v0.38.0/webvtt.go#L493) ¶
    
    
    func (s Subtitles) WriteToWebVTT(o [io](/io).[Writer](/io#Writer)) (err [error](/builtin#error))

WriteToWebVTT writes subtitles in .vtt format 

####  type [TTMLIn](https://github.com/asticode/go-astisub/blob/v0.38.0/ttml.go#L63) ¶
    
    
    type TTMLIn struct {
    	Framerate [int](/builtin#int)            `xml:"frameRate,attr"`
    	Lang      [string](/builtin#string)         `xml:"lang,attr"`
    	Metadata  TTMLInMetadata `xml:"head>metadata"`
    	Regions   []TTMLInRegion `xml:"head>layout>region"`
    	Styles    []TTMLInStyle  `xml:"head>styling>style"`
    	Body      TTMLInBody     `xml:"body"`
    	Tickrate  [int](/builtin#int)            `xml:"tickRate,attr"`
    	XMLName   [xml](/encoding/xml).[Name](/encoding/xml#Name)       `xml:"tt"`
    }

TTMLIn represents an input TTML that must be unmarshaled We split it from the output TTML as we can't add strict namespace without breaking retrocompatibility 

####  type [TTMLInBody](https://github.com/asticode/go-astisub/blob/v0.38.0/ttml.go#L52) ¶ added in v0.36.0
    
    
    type TTMLInBody struct {
    	XMLName [xml](/encoding/xml).[Name](/encoding/xml#Name)        `xml:"body"`
    	Divs    []TTMLInBodyDiv `xml:"div"`
    
    	Region [string](/builtin#string) `xml:"region,attr,omitempty"`
    	Style  [string](/builtin#string) `xml:"style,attr,omitempty"`
    	TTMLInStyleAttributes
    }

####  type [TTMLInBodyDiv](https://github.com/asticode/go-astisub/blob/v0.38.0/ttml.go#L44) ¶ added in v0.36.0
    
    
    type TTMLInBodyDiv struct {
    	XMLName   [xml](/encoding/xml).[Name](/encoding/xml#Name)         `xml:"div"`
    	Subtitles []TTMLInSubtitle `xml:"p"`
    
    	Region [string](/builtin#string) `xml:"region,attr,omitempty"`
    	Style  [string](/builtin#string) `xml:"style,attr,omitempty"`
    	TTMLInStyleAttributes
    }

####  type [TTMLInDuration](https://github.com/asticode/go-astisub/blob/v0.38.0/ttml.go#L265) ¶
    
    
    type TTMLInDuration struct {
    	// contains filtered or unexported fields
    }

TTMLInDuration represents an input TTML duration 

####  func (*TTMLInDuration) [UnmarshalText](https://github.com/asticode/go-astisub/blob/v0.38.0/ttml.go#L276) ¶
    
    
    func (d *TTMLInDuration) UnmarshalText(i [][byte](/builtin#byte)) (err [error](/builtin#error))

UnmarshalText implements the TextUnmarshaler interface Possible formats are: - hh:mm:ss.mmm - hh:mm:ss:fff (fff being frames) - [ticks]t ([ticks] being the tick amount) 

####  type [TTMLInHeader](https://github.com/asticode/go-astisub/blob/v0.38.0/ttml.go#L154) ¶
    
    
    type TTMLInHeader struct {
    	ID    [string](/builtin#string) `xml:"id,attr,omitempty"`
    	Style [string](/builtin#string) `xml:"style,attr,omitempty"`
    	TTMLInStyleAttributes
    }

TTMLInHeader represents an input TTML header 

####  type [TTMLInItem](https://github.com/asticode/go-astisub/blob/v0.38.0/ttml.go#L257) ¶
    
    
    type TTMLInItem struct {
    	Style [string](/builtin#string) `xml:"style,attr,omitempty"`
    	Text  [string](/builtin#string) `xml:",chardata"`
    	TTMLInStyleAttributes
    	XMLName [xml](/encoding/xml).[Name](/encoding/xml#Name)
    }

TTMLInItem represents an input TTML item 

####  type [TTMLInItems](https://github.com/asticode/go-astisub/blob/v0.38.0/ttml.go#L186) ¶
    
    
    type TTMLInItems []TTMLInItem

TTMLInItems represents input TTML items 

####  func (*TTMLInItems) [UnmarshalXML](https://github.com/asticode/go-astisub/blob/v0.38.0/ttml.go#L189) ¶
    
    
    func (i *TTMLInItems) UnmarshalXML(d *[xml](/encoding/xml).[Decoder](/encoding/xml#Decoder), start [xml](/encoding/xml).[StartElement](/encoding/xml#StartElement)) (err [error](/builtin#error))

UnmarshalXML implements the XML unmarshaler interface 

####  type [TTMLInMetadata](https://github.com/asticode/go-astisub/blob/v0.38.0/ttml.go#L88) ¶
    
    
    type TTMLInMetadata struct {
    	Copyright [string](/builtin#string) `xml:"copyright"`
    	Title     [string](/builtin#string) `xml:"title"`
    }

TTMLInMetadata represents an input TTML Metadata 

####  type [TTMLInRegion](https://github.com/asticode/go-astisub/blob/v0.38.0/ttml.go#L161) ¶
    
    
    type TTMLInRegion struct {
    	TTMLInHeader
    	XMLName [xml](/encoding/xml).[Name](/encoding/xml#Name) `xml:"region"`
    }

TTMLInRegion represents an input TTML region 

####  type [TTMLInStyle](https://github.com/asticode/go-astisub/blob/v0.38.0/ttml.go#L167) ¶
    
    
    type TTMLInStyle struct {
    	TTMLInHeader
    	XMLName [xml](/encoding/xml).[Name](/encoding/xml#Name) `xml:"style"`
    }

TTMLInStyle represents an input TTML style 

####  type [TTMLInStyleAttributes](https://github.com/asticode/go-astisub/blob/v0.38.0/ttml.go#L94) ¶
    
    
    type TTMLInStyleAttributes struct {
    	BackgroundColor *[string](/builtin#string) `xml:"backgroundColor,attr,omitempty"`
    	Color           *[string](/builtin#string) `xml:"color,attr,omitempty"`
    	Direction       *[string](/builtin#string) `xml:"direction,attr,omitempty"`
    	Display         *[string](/builtin#string) `xml:"display,attr,omitempty"`
    	DisplayAlign    *[string](/builtin#string) `xml:"displayAlign,attr,omitempty"`
    	Extent          *[string](/builtin#string) `xml:"extent,attr,omitempty"`
    	FontFamily      *[string](/builtin#string) `xml:"fontFamily,attr,omitempty"`
    	FontSize        *[string](/builtin#string) `xml:"fontSize,attr,omitempty"`
    	FontStyle       *[string](/builtin#string) `xml:"fontStyle,attr,omitempty"`
    	FontWeight      *[string](/builtin#string) `xml:"fontWeight,attr,omitempty"`
    	LineHeight      *[string](/builtin#string) `xml:"lineHeight,attr,omitempty"`
    	Opacity         *[string](/builtin#string) `xml:"opacity,attr,omitempty"`
    	Origin          *[string](/builtin#string) `xml:"origin,attr,omitempty"`
    	Overflow        *[string](/builtin#string) `xml:"overflow,attr,omitempty"`
    	Padding         *[string](/builtin#string) `xml:"padding,attr,omitempty"`
    	ShowBackground  *[string](/builtin#string) `xml:"showBackground,attr,omitempty"`
    	TextAlign       *[string](/builtin#string) `xml:"textAlign,attr,omitempty"`
    	TextDecoration  *[string](/builtin#string) `xml:"textDecoration,attr,omitempty"`
    	TextOutline     *[string](/builtin#string) `xml:"textOutline,attr,omitempty"`
    	UnicodeBidi     *[string](/builtin#string) `xml:"unicodeBidi,attr,omitempty"`
    	Visibility      *[string](/builtin#string) `xml:"visibility,attr,omitempty"`
    	WrapOption      *[string](/builtin#string) `xml:"wrapOption,attr,omitempty"`
    	WritingMode     *[string](/builtin#string) `xml:"writingMode,attr,omitempty"`
    	ZIndex          *[int](/builtin#int)    `xml:"zIndex,attr,omitempty"`
    }

TTMLInStyleAttributes represents input TTML style attributes 

####  type [TTMLInSubtitle](https://github.com/asticode/go-astisub/blob/v0.38.0/ttml.go#L173) ¶
    
    
    type TTMLInSubtitle struct {
    	Begin *TTMLInDuration `xml:"begin,attr,omitempty"`
    	End   *TTMLInDuration `xml:"end,attr,omitempty"`
    	ID    [string](/builtin#string)          `xml:"id,attr,omitempty"`
    	// We must store inner XML temporarily here since there's no tag to describe both any tag and chardata
    	// Real unmarshal will be done manually afterwards
    	Items  [string](/builtin#string) `xml:",innerxml"`
    	Region [string](/builtin#string) `xml:"region,attr,omitempty"`
    	Style  [string](/builtin#string) `xml:"style,attr,omitempty"`
    	TTMLInStyleAttributes
    }

TTMLInSubtitle represents an input TTML subtitle 

####  type [TTMLOut](https://github.com/asticode/go-astisub/blob/v0.38.0/ttml.go#L527) ¶
    
    
    type TTMLOut struct {
    	Lang            [string](/builtin#string)            `xml:"xml:lang,attr,omitempty"`
    	Metadata        *TTMLOutMetadata  `xml:"head>metadata,omitempty"`
    	Styles          []TTMLOutStyle    `xml:"head>styling>style,omitempty"` //!\\ Order is important! Keep Styling above Layout
    	Regions         []TTMLOutRegion   `xml:"head>layout>region,omitempty"`
    	Subtitles       []TTMLOutSubtitle `xml:"body>div>p,omitempty"`
    	XMLName         [xml](/encoding/xml).[Name](/encoding/xml#Name)          `xml:"http://www.w3.org/ns/ttml tt"`
    	XMLNamespaceTTM [string](/builtin#string)            `xml:"xmlns:ttm,attr"`
    	XMLNamespaceTTS [string](/builtin#string)            `xml:"xmlns:tts,attr"`
    }

TTMLOut represents an output TTML that must be marshaled We split it from the input TTML as this time we'll add strict namespaces 

####  type [TTMLOutDuration](https://github.com/asticode/go-astisub/blob/v0.38.0/ttml.go#L644) ¶
    
    
    type TTMLOutDuration [time](/time).[Duration](/time#Duration)

TTMLOutDuration represents an output TTML duration 

####  func (TTMLOutDuration) [MarshalText](https://github.com/asticode/go-astisub/blob/v0.38.0/ttml.go#L647) ¶
    
    
    func (t TTMLOutDuration) MarshalText() ([][byte](/builtin#byte), [error](/builtin#error))

MarshalText implements the TextMarshaler interface 

####  type [TTMLOutHeader](https://github.com/asticode/go-astisub/blob/v0.38.0/ttml.go#L606) ¶
    
    
    type TTMLOutHeader struct {
    	ID    [string](/builtin#string) `xml:"xml:id,attr,omitempty"`
    	Style [string](/builtin#string) `xml:"style,attr,omitempty"`
    	TTMLOutStyleAttributes
    }

TTMLOutHeader represents an output TTML header 

####  type [TTMLOutItem](https://github.com/asticode/go-astisub/blob/v0.38.0/ttml.go#L636) ¶
    
    
    type TTMLOutItem struct {
    	Style [string](/builtin#string) `xml:"style,attr,omitempty"`
    	Text  [string](/builtin#string) `xml:",chardata"`
    	TTMLOutStyleAttributes
    	XMLName [xml](/encoding/xml).[Name](/encoding/xml#Name)
    }

TTMLOutItem represents an output TTML Item 

####  type [TTMLOutMetadata](https://github.com/asticode/go-astisub/blob/v0.38.0/ttml.go#L539) ¶
    
    
    type TTMLOutMetadata struct {
    	Copyright [string](/builtin#string) `xml:"ttm:copyright,omitempty"`
    	Title     [string](/builtin#string) `xml:"ttm:title,omitempty"`
    }

TTMLOutMetadata represents an output TTML Metadata 

####  type [TTMLOutRegion](https://github.com/asticode/go-astisub/blob/v0.38.0/ttml.go#L613) ¶
    
    
    type TTMLOutRegion struct {
    	TTMLOutHeader
    	XMLName [xml](/encoding/xml).[Name](/encoding/xml#Name) `xml:"region"`
    }

TTMLOutRegion represents an output TTML region 

####  type [TTMLOutStyle](https://github.com/asticode/go-astisub/blob/v0.38.0/ttml.go#L619) ¶
    
    
    type TTMLOutStyle struct {
    	TTMLOutHeader
    	XMLName [xml](/encoding/xml).[Name](/encoding/xml#Name) `xml:"style"`
    }

TTMLOutStyle represents an output TTML style 

####  type [TTMLOutStyleAttributes](https://github.com/asticode/go-astisub/blob/v0.38.0/ttml.go#L545) ¶
    
    
    type TTMLOutStyleAttributes struct {
    	BackgroundColor *[string](/builtin#string) `xml:"tts:backgroundColor,attr,omitempty"`
    	Color           *[string](/builtin#string) `xml:"tts:color,attr,omitempty"`
    	Direction       *[string](/builtin#string) `xml:"tts:direction,attr,omitempty"`
    	Display         *[string](/builtin#string) `xml:"tts:display,attr,omitempty"`
    	DisplayAlign    *[string](/builtin#string) `xml:"tts:displayAlign,attr,omitempty"`
    	Extent          *[string](/builtin#string) `xml:"tts:extent,attr,omitempty"`
    	FontFamily      *[string](/builtin#string) `xml:"tts:fontFamily,attr,omitempty"`
    	FontSize        *[string](/builtin#string) `xml:"tts:fontSize,attr,omitempty"`
    	FontStyle       *[string](/builtin#string) `xml:"tts:fontStyle,attr,omitempty"`
    	FontWeight      *[string](/builtin#string) `xml:"tts:fontWeight,attr,omitempty"`
    	LineHeight      *[string](/builtin#string) `xml:"tts:lineHeight,attr,omitempty"`
    	Opacity         *[string](/builtin#string) `xml:"tts:opacity,attr,omitempty"`
    	Origin          *[string](/builtin#string) `xml:"tts:origin,attr,omitempty"`
    	Overflow        *[string](/builtin#string) `xml:"tts:overflow,attr,omitempty"`
    	Padding         *[string](/builtin#string) `xml:"tts:padding,attr,omitempty"`
    	ShowBackground  *[string](/builtin#string) `xml:"tts:showBackground,attr,omitempty"`
    	TextAlign       *[string](/builtin#string) `xml:"tts:textAlign,attr,omitempty"`
    	TextDecoration  *[string](/builtin#string) `xml:"tts:textDecoration,attr,omitempty"`
    	TextOutline     *[string](/builtin#string) `xml:"tts:textOutline,attr,omitempty"`
    	UnicodeBidi     *[string](/builtin#string) `xml:"tts:unicodeBidi,attr,omitempty"`
    	Visibility      *[string](/builtin#string) `xml:"tts:visibility,attr,omitempty"`
    	WrapOption      *[string](/builtin#string) `xml:"tts:wrapOption,attr,omitempty"`
    	WritingMode     *[string](/builtin#string) `xml:"tts:writingMode,attr,omitempty"`
    	ZIndex          *[int](/builtin#int)    `xml:"tts:zIndex,attr,omitempty"`
    }

TTMLOutStyleAttributes represents output TTML style attributes 

####  type [TTMLOutSubtitle](https://github.com/asticode/go-astisub/blob/v0.38.0/ttml.go#L625) ¶
    
    
    type TTMLOutSubtitle struct {
    	Begin  TTMLOutDuration `xml:"begin,attr"`
    	End    TTMLOutDuration `xml:"end,attr"`
    	ID     [string](/builtin#string)          `xml:"id,attr,omitempty"`
    	Items  []TTMLOutItem
    	Region [string](/builtin#string) `xml:"region,attr,omitempty"`
    	Style  [string](/builtin#string) `xml:"style,attr,omitempty"`
    	TTMLOutStyleAttributes
    }

TTMLOutSubtitle represents an output TTML subtitle 

####  type [TeletextOptions](https://github.com/asticode/go-astisub/blob/v0.38.0/teletext.go#L325) ¶
    
    
    type TeletextOptions struct {
    	Page [int](/builtin#int)
    	PID  [int](/builtin#int)
    }

TeletextOptions represents teletext options 

####  type [WebVTTPosition](https://github.com/asticode/go-astisub/blob/v0.38.0/webvtt.go#L40) ¶ added in v0.38.0
    
    
    type WebVTTPosition struct {
    	XPosition [string](/builtin#string)
    	Alignment [string](/builtin#string)
    }

####  func (*WebVTTPosition) [String](https://github.com/asticode/go-astisub/blob/v0.38.0/webvtt.go#L63) ¶ added in v0.38.0
    
    
    func (p *WebVTTPosition) String() [string](/builtin#string)

####  type [WebVTTTag](https://github.com/asticode/go-astisub/blob/v0.38.0/subtitles.go#L270) ¶ added in v0.26.1
    
    
    type WebVTTTag struct {
    	Name       [string](/builtin#string)
    	Annotation [string](/builtin#string)
    	Classes    [][string](/builtin#string)
    }

####  type [WebVTTTimestampMap](https://github.com/asticode/go-astisub/blob/v0.38.0/webvtt.go#L81) ¶ added in v0.27.0
    
    
    type WebVTTTimestampMap struct {
    	Local  [time](/time).[Duration](/time#Duration)
    	MpegTS [int64](/builtin#int64)
    }

WebVTTTimestampMap is a structure for storing timestamps for WEBVTT's X-TIMESTAMP-MAP feature commonly used for syncing cue times with MPEG-TS streams. 

####  func (*WebVTTTimestampMap) [Offset](https://github.com/asticode/go-astisub/blob/v0.38.0/webvtt.go#L88) ¶ added in v0.27.0
    
    
    func (t *WebVTTTimestampMap) Offset() [time](/time).[Duration](/time#Duration)

Offset calculates and returns the time offset described by the timestamp map. 

####  func (*WebVTTTimestampMap) [String](https://github.com/asticode/go-astisub/blob/v0.38.0/webvtt.go#L97) ¶ added in v0.27.0
    
    
    func (t *WebVTTTimestampMap) String() [string](/builtin#string)

String implements Stringer interface for TimestampMap, returning the fully formatted header string for the instance. 

####  type [WriteToTTMLOption](https://github.com/asticode/go-astisub/blob/v0.38.0/ttml.go#L657) ¶ added in v0.28.0
    
    
    type WriteToTTMLOption func(o *WriteToTTMLOptions)

WriteToTTMLOption represents a WriteToTTML option. 

####  func [WriteToTTMLWithIndentOption](https://github.com/asticode/go-astisub/blob/v0.38.0/ttml.go#L660) ¶ added in v0.28.0
    
    
    func WriteToTTMLWithIndentOption(indent [string](/builtin#string)) WriteToTTMLOption

WriteToTTMLWithIndentOption sets the indent option. 

####  type [WriteToTTMLOptions](https://github.com/asticode/go-astisub/blob/v0.38.0/ttml.go#L652) ¶ added in v0.28.0
    
    
    type WriteToTTMLOptions struct {
    	Indent [string](/builtin#string) // Default is 4 spaces.
    }

WriteToTTMLOptions represents TTML write options. 
  *[↑]: Back to Top
  *[v]: View this template
  *[t]: Discuss this template
  *[e]: Edit this template
