# bimg (libvips)

> Auto-fetched from [https://pkg.go.dev/github.com/h2non/bimg](https://pkg.go.dev/github.com/h2non/bimg)
> Last Updated: 2026-01-28T21:46:45.834707+00:00

---

Index
¶
Constants
Variables
func ColourspaceIsSupported(buf []byte) (bool, error)
func DetermineImageTypeName(buf []byte) string
func ImageTypeName(t ImageType) string
func Initialize()
func IsSVGImage(buf []byte) bool
func IsTypeNameSupported(t string) bool
func IsTypeNameSupportedSave(t string) bool
func IsTypeSupported(t ImageType) bool
func IsTypeSupportedSave(t ImageType) bool
func MaxSize() int
func Read(path string) ([]byte, error)
func Resize(buf []byte, o Options) ([]byte, error)
func SetMaxsize(s int) error
func Shutdown()
func VipsCacheDropAll()
func VipsCacheSetMax(maxCacheSize int)
func VipsCacheSetMaxMem(maxCacheMem int)
func VipsDebugInfo()
func VipsIsTypeSupported(t ImageType) bool
func VipsIsTypeSupportedSave(t ImageType) bool
func VipsVectorSetEnabled(enable bool)
func Write(path string, buf []byte) error
type Angle
type Color
type Direction
type EXIF
type Extend
type GaussianBlur
type Gravity
type Image
func NewImage(buf []byte) *Image
func (i *Image) AutoRotate() ([]byte, error)
func (i *Image) Colourspace(c Interpretation) ([]byte, error)
func (i *Image) ColourspaceIsSupported() (bool, error)
func (i *Image) Convert(t ImageType) ([]byte, error)
func (i *Image) Crop(width, height int, gravity Gravity) ([]byte, error)
func (i *Image) CropByHeight(height int) ([]byte, error)
func (i *Image) CropByWidth(width int) ([]byte, error)
func (i *Image) Enlarge(width, height int) ([]byte, error)
func (i *Image) EnlargeAndCrop(width, height int) ([]byte, error)
func (i *Image) Extract(top, left, width, height int) ([]byte, error)
func (i *Image) Flip() ([]byte, error)
func (i *Image) Flop() ([]byte, error)
func (i *Image) ForceResize(width, height int) ([]byte, error)
func (i *Image) Gamma(exponent float64) ([]byte, error)
func (i *Image) Image() []byte
func (i *Image) Interpretation() (Interpretation, error)
func (i *Image) Length() int
func (i *Image) Metadata() (ImageMetadata, error)
func (i *Image) Process(o Options) ([]byte, error)
func (i *Image) Resize(width, height int) ([]byte, error)
func (i *Image) ResizeAndCrop(width, height int) ([]byte, error)
func (i *Image) Rotate(a Angle) ([]byte, error)
func (i *Image) Size() (ImageSize, error)
func (i *Image) SmartCrop(width, height int) ([]byte, error)
func (i *Image) Thumbnail(pixels int) ([]byte, error)
func (i *Image) Trim() ([]byte, error)
func (i *Image) Type() string
func (i *Image) Watermark(w Watermark) ([]byte, error)
func (i *Image) WatermarkImage(w WatermarkImage) ([]byte, error)
func (i *Image) Zoom(factor int) ([]byte, error)
type ImageMetadata
func Metadata(buf []byte) (ImageMetadata, error)
type ImageSize
func Size(buf []byte) (ImageSize, error)
type ImageType
func DetermineImageType(buf []byte) ImageType
type Interpolator
func (i Interpolator) String() string
type Interpretation
func ImageInterpretation(buf []byte) (Interpretation, error)
type Options
type Sharpen
type SupportedImageType
func IsImageTypeSupportedByVips(t ImageType) SupportedImageType
type VipsMemoryInfo
func VipsMemory() VipsMemoryInfo
type Watermark
type WatermarkImage
Constants
¶
View Source
const (
Make                    = "exif-ifd0-Make"
Model                   = "exif-ifd0-Model"
Orientation             = "exif-ifd0-Orientation"
XResolution             = "exif-ifd0-XResolution"
YResolution             = "exif-ifd0-YResolution"
ResolutionUnit          = "exif-ifd0-ResolutionUnit"
Software                = "exif-ifd0-Software"
Datetime                = "exif-ifd0-DateTime"
YCbCrPositioning        = "exif-ifd0-YCbCrPositioning"
Compression             = "exif-ifd1-Compression"
ExposureTime            = "exif-ifd2-ExposureTime"
FNumber                 = "exif-ifd2-FNumber"
ExposureProgram         = "exif-ifd2-ExposureProgram"
ISOSpeedRatings         = "exif-ifd2-ISOSpeedRatings"
ExifVersion             = "exif-ifd2-ExifVersion"
DateTimeOriginal        = "exif-ifd2-DateTimeOriginal"
DateTimeDigitized       = "exif-ifd2-DateTimeDigitized"
ComponentsConfiguration = "exif-ifd2-ComponentsConfiguration"
ShutterSpeedValue       = "exif-ifd2-ShutterSpeedValue"
ApertureValue           = "exif-ifd2-ApertureValue"
BrightnessValue         = "exif-ifd2-BrightnessValue"
ExposureBiasValue       = "exif-ifd2-ExposureBiasValue"
MeteringMode            = "exif-ifd2-MeteringMode"
Flash                   = "exif-ifd2-Flash"
FocalLength             = "exif-ifd2-FocalLength"
SubjectArea             = "exif-ifd2-SubjectArea"
MakerNote               = "exif-ifd2-MakerNote"
SubSecTimeOriginal      = "exif-ifd2-SubSecTimeOriginal"
SubSecTimeDigitized     = "exif-ifd2-SubSecTimeDigitized"
ColorSpace              = "exif-ifd2-ColorSpace"
PixelXDimension         = "exif-ifd2-PixelXDimension"
PixelYDimension         = "exif-ifd2-PixelYDimension"
SensingMethod           = "exif-ifd2-SensingMethod"
SceneType               = "exif-ifd2-SceneType"
ExposureMode            = "exif-ifd2-ExposureMode"
WhiteBalance            = "exif-ifd2-WhiteBalance"
FocalLengthIn35mmFilm   = "exif-ifd2-FocalLengthIn35mmFilm"
SceneCaptureType        = "exif-ifd2-SceneCaptureType"
GPSLatitudeRef          = "exif-ifd3-GPSLatitudeRef"
GPSLatitude             = "exif-ifd3-GPSLatitude"
GPSLongitudeRef         = "exif-ifd3-GPSLongitudeRef"
GPSLongitude            = "exif-ifd3-GPSLongitude"
GPSAltitudeRef          = "exif-ifd3-GPSAltitudeRef"
GPSAltitude             = "exif-ifd3-GPSAltitude"
GPSSpeedRef             = "exif-ifd3-GPSSpeedRef"
GPSSpeed                = "exif-ifd3-GPSSpeed"
GPSImgDirectionRef      = "exif-ifd3-GPSImgDirectionRef"
GPSImgDirection         = "exif-ifd3-GPSImgDirection"
GPSDestBearingRef       = "exif-ifd3-GPSDestBearingRef"
GPSDestBearing          = "exif-ifd3-GPSDestBearing"
GPSDateStamp            = "exif-ifd3-GPSDateStamp"
)
Common EXIF fields for data extraction
View Source
const (
// Quality defines the default JPEG quality to be used.
Quality = 75
)
View Source
const Version = "1.1.9"
Version represents the current package semantic version.
View Source
const VipsMajorVersion =
int
(
C
.
VIPS_MAJOR_VERSION
)
VipsMajorVersion exposes the current libvips major version number
View Source
const VipsMinorVersion =
int
(
C
.
VIPS_MINOR_VERSION
)
VipsMinorVersion exposes the current libvips minor version number
View Source
const VipsVersion =
string
(
C
.
VIPS_VERSION
)
VipsVersion exposes the current libvips semantic version
Variables
¶
View Source
var ColorBlack =
Color
{0, 0, 0}
ColorBlack is a shortcut to black RGB color representation.
View Source
var (
// ErrExtractAreaParamsRequired defines a generic extract area error
ErrExtractAreaParamsRequired =
errors
.
New
("extract area width/height params are required")
)
View Source
var ImageTypes = map[
ImageType
]
string
{
JPEG
:   "jpeg",
PNG
:    "png",
WEBP
:   "webp",
TIFF
:   "tiff",
GIF
:    "gif",
PDF
:    "pdf",
SVG
:    "svg",
MAGICK
: "magick",
HEIF
:   "heif",
AVIF
:   "avif",
}
ImageTypes stores as pairs of image types supported and its alias names.
View Source
var SupportedImageTypes = map[
ImageType
]
SupportedImageType
{}
SupportedImageTypes stores the optional image type supported
by the current libvips compilation.
Note: lazy evaluation as demand is required due
to bootstrap runtime limitation with C/libvips world.
View Source
var WatermarkFont = "sans 10"
WatermarkFont defines the default watermark font to be used.
Functions
¶
func
ColourspaceIsSupported
¶
func ColourspaceIsSupported(buf []
byte
) (
bool
,
error
)
ColourspaceIsSupported checks if the image colourspace is supported by libvips.
func
DetermineImageTypeName
¶
func DetermineImageTypeName(buf []
byte
)
string
DetermineImageTypeName determines the image type format by name (jpeg, png, webp or tiff)
func
ImageTypeName
¶
func ImageTypeName(t
ImageType
)
string
ImageTypeName is used to get the human friendly name of an image format.
func
Initialize
¶
func Initialize()
Initialize is used to explicitly start libvips in thread-safe way.
Only call this function if you have previously turned off libvips.
func
IsSVGImage
¶
added in
v1.0.3
func IsSVGImage(buf []
byte
)
bool
IsSVGImage returns true if the given buffer is a valid SVG image.
func
IsTypeNameSupported
¶
func IsTypeNameSupported(t
string
)
bool
IsTypeNameSupported checks if a given image type name is supported
func
IsTypeNameSupportedSave
¶
added in
v1.0.7
func IsTypeNameSupportedSave(t
string
)
bool
IsTypeNameSupportedSave checks if a given image type name is supported for
saving
func
IsTypeSupported
¶
func IsTypeSupported(t
ImageType
)
bool
IsTypeSupported checks if a given image type is supported
func
IsTypeSupportedSave
¶
added in
v1.0.7
func IsTypeSupportedSave(t
ImageType
)
bool
IsTypeSupportedSave checks if a given image type is support for saving
func
MaxSize
¶
added in
v1.0.0
func MaxSize()
int
MaxSize returns maxSize.
func
Read
¶
func Read(path
string
) ([]
byte
,
error
)
Read reads all the content of the given file path
and returns it as byte buffer.
func
Resize
¶
func Resize(buf []
byte
, o
Options
) ([]
byte
,
error
)
Resize is used to transform a given image as byte buffer
with the passed options.
func
SetMaxsize
¶
added in
v1.1.8
func SetMaxsize(s
int
)
error
SetMaxSize sets maxSize.
func
Shutdown
¶
func Shutdown()
Shutdown is used to shutdown libvips in a thread-safe way.
You can call this to drop caches as well.
If libvips was already initialized, the function is no-op
func
VipsCacheDropAll
¶
added in
v1.0.10
func VipsCacheDropAll()
VipsCacheDropAll drops the vips operation cache, freeing the allocated memory.
func
VipsCacheSetMax
¶
added in
v1.0.10
func VipsCacheSetMax(maxCacheSize
int
)
VipsCacheSetMax sets the maximum number of operations to keep in the vips operation cache.
func
VipsCacheSetMaxMem
¶
added in
v1.0.10
func VipsCacheSetMaxMem(maxCacheMem
int
)
VipsCacheSetMaxMem Sets the maximum amount of tracked memory allowed before the vips operation cache
begins to drop entries.
func
VipsDebugInfo
¶
func VipsDebugInfo()
VipsDebugInfo outputs to stdout libvips collected data. Useful for debugging.
func
VipsIsTypeSupported
¶
added in
v1.0.3
func VipsIsTypeSupported(t
ImageType
)
bool
VipsIsTypeSupported returns true if the given image type
is supported by the current libvips compilation.
func
VipsIsTypeSupportedSave
¶
added in
v1.0.7
func VipsIsTypeSupportedSave(t
ImageType
)
bool
VipsIsTypeSupportedSave returns true if the given image type
is supported by the current libvips compilation for the
save operation.
func
VipsVectorSetEnabled
¶
added in
v1.1.6
func VipsVectorSetEnabled(enable
bool
)
VipsVectorSetEnabled enables or disables SIMD vector instructions. This can give speed-up,
but can also be unstable on some systems and versions.
func
Write
¶
func Write(path
string
, buf []
byte
)
error
Write writes the given byte buffer into disk
to the given file path.
Types
¶
type
Angle
¶
type Angle
int
Angle represents the image rotation angle value.
const (
// D0 represents the rotation angle 0 degrees.
D0
Angle
= 0
// D45 represents the rotation angle 45 degrees.
D45
Angle
= 45
// D90 represents the rotation angle 90 degrees.
D90
Angle
= 90
// D135 represents the rotation angle 135 degrees.
D135
Angle
= 135
// D180 represents the rotation angle 180 degrees.
D180
Angle
= 180
// D235 represents the rotation angle 235 degrees.
D235
Angle
= 235
// D270 represents the rotation angle 270 degrees.
D270
Angle
= 270
// D315 represents the rotation angle 315 degrees.
D315
Angle
= 315
)
type
Color
¶
type Color struct {
R, G, B
uint8
}
Color represents a traditional RGB color scheme.
type
Direction
¶
type Direction
int
Direction represents the image direction value.
const (
// Horizontal represents the orizontal image direction value.
Horizontal
Direction
=
C
.
VIPS_DIRECTION_HORIZONTAL
// Vertical represents the vertical image direction value.
Vertical
Direction
=
C
.
VIPS_DIRECTION_VERTICAL
)
type
EXIF
¶
added in
v1.1.3
type EXIF struct {
Make
string
Model
string
Orientation
int
XResolution
string
YResolution
string
ResolutionUnit
int
Software
string
Datetime
string
YCbCrPositioning
int
Compression
int
ExposureTime
string
FNumber
string
ExposureProgram
int
ISOSpeedRatings
int
ExifVersion
string
DateTimeOriginal
string
DateTimeDigitized
string
ComponentsConfiguration
string
ShutterSpeedValue
string
ApertureValue
string
BrightnessValue
string
ExposureBiasValue
string
MeteringMode
int
Flash
int
FocalLength
string
SubjectArea
string
MakerNote
string
SubSecTimeOriginal
string
SubSecTimeDigitized
string
ColorSpace
int
PixelXDimension
int
PixelYDimension
int
SensingMethod
int
SceneType
string
ExposureMode
int
WhiteBalance
int
FocalLengthIn35mmFilm
int
SceneCaptureType
int
GPSLatitudeRef
string
GPSLatitude
string
GPSLongitudeRef
string
GPSLongitude
string
GPSAltitudeRef
string
GPSAltitude
string
GPSSpeedRef
string
GPSSpeed
string
GPSImgDirectionRef
string
GPSImgDirection
string
GPSDestBearingRef
string
GPSDestBearing
string
GPSDateStamp
string
}
EXIF image metadata
type
Extend
¶
added in
v1.0.5
type Extend
int
Extend represents the image extend mode, used when the edges
of an image are extended, you can specify how you want the extension done.
See:
https://libvips.github.io/libvips/API/current/libvips-conversion.html#VIPS-EXTEND-BACKGROUND:CAPS
const (
// ExtendBlack extend with black (all 0) pixels mode.
ExtendBlack
Extend
=
C
.
VIPS_EXTEND_BLACK
// ExtendCopy copy the image edges.
ExtendCopy
Extend
=
C
.
VIPS_EXTEND_COPY
// ExtendRepeat repeat the whole image.
ExtendRepeat
Extend
=
C
.
VIPS_EXTEND_REPEAT
// ExtendMirror mirror the whole image.
ExtendMirror
Extend
=
C
.
VIPS_EXTEND_MIRROR
// ExtendWhite extend with white (all bits set) pixels.
ExtendWhite
Extend
=
C
.
VIPS_EXTEND_WHITE
// ExtendBackground with colour from the background property.
ExtendBackground
Extend
=
C
.
VIPS_EXTEND_BACKGROUND
// ExtendLast extend with last pixel.
ExtendLast
Extend
=
C
.
VIPS_EXTEND_LAST
)
type
GaussianBlur
¶
type GaussianBlur struct {
Sigma
float64
MinAmpl
float64
}
GaussianBlur represents the gaussian image transformation values.
type
Gravity
¶
type Gravity
int
Gravity represents the image gravity value.
const (
// GravityCentre represents the centre value used for image gravity orientation.
GravityCentre
Gravity
=
iota
// GravityNorth represents the north value used for image gravity orientation.
GravityNorth
// GravityEast represents the east value used for image gravity orientation.
GravityEast
// GravitySouth represents the south value used for image gravity orientation.
GravitySouth
// GravityWest represents the west value used for image gravity orientation.
GravityWest
// GravitySmart enables libvips Smart Crop algorithm for image gravity orientation.
GravitySmart
)
type
Image
¶
type Image struct {
// contains filtered or unexported fields
}
Image provides a simple method DSL to transform a given image as byte buffer.
func
NewImage
¶
func NewImage(buf []
byte
) *
Image
NewImage creates a new Image struct with method DSL.
func (*Image)
AutoRotate
¶
added in
v1.1.3
func (i *
Image
) AutoRotate() ([]
byte
,
error
)
AutoRotate automatically rotates the image with no additional transformation based on the EXIF oritentation metadata, if available.
func (*Image)
Colourspace
¶
func (i *
Image
) Colourspace(c
Interpretation
) ([]
byte
,
error
)
Colourspace performs a color space conversion bsaed on the given interpretation.
func (*Image)
ColourspaceIsSupported
¶
func (i *
Image
) ColourspaceIsSupported() (
bool
,
error
)
ColourspaceIsSupported checks if the current image
color space is supported.
func (*Image)
Convert
¶
func (i *
Image
) Convert(t
ImageType
) ([]
byte
,
error
)
Convert converts image to another format.
func (*Image)
Crop
¶
func (i *
Image
) Crop(width, height
int
, gravity
Gravity
) ([]
byte
,
error
)
Crop crops the image to the exact size specified.
func (*Image)
CropByHeight
¶
func (i *
Image
) CropByHeight(height
int
) ([]
byte
,
error
)
CropByHeight crops an image by height (auto width).
func (*Image)
CropByWidth
¶
func (i *
Image
) CropByWidth(width
int
) ([]
byte
,
error
)
CropByWidth crops an image by width only param (auto height).
func (*Image)
Enlarge
¶
func (i *
Image
) Enlarge(width, height
int
) ([]
byte
,
error
)
Enlarge enlarges the image by width and height. Aspect ratio is maintained.
func (*Image)
EnlargeAndCrop
¶
func (i *
Image
) EnlargeAndCrop(width, height
int
) ([]
byte
,
error
)
EnlargeAndCrop enlarges the image by width and height with additional crop transformation.
func (*Image)
Extract
¶
func (i *
Image
) Extract(top, left, width, height
int
) ([]
byte
,
error
)
Extract area from the by X/Y axis in the current image.
func (*Image)
Flip
¶
func (i *
Image
) Flip() ([]
byte
,
error
)
Flip flips the image about the vertical Y axis.
func (*Image)
Flop
¶
func (i *
Image
) Flop() ([]
byte
,
error
)
Flop flops the image about the horizontal X axis.
func (*Image)
ForceResize
¶
func (i *
Image
) ForceResize(width, height
int
) ([]
byte
,
error
)
ForceResize resizes with custom size (aspect ratio won't be maintained).
func (*Image)
Gamma
¶
added in
v1.1.0
func (i *
Image
) Gamma(exponent
float64
) ([]
byte
,
error
)
Gamma returns the gamma filtered image buffer.
func (*Image)
Image
¶
func (i *
Image
) Image() []
byte
Image returns the current resultant image buffer.
func (*Image)
Interpretation
¶
func (i *
Image
) Interpretation() (
Interpretation
,
error
)
Interpretation gets the image interpretation type.
See:
https://libvips.github.io/libvips/API/current/VipsImage.html#VipsInterpretation
func (*Image)
Length
¶
added in
v1.0.10
func (i *
Image
) Length()
int
Length returns the size in bytes of the image buffer.
func (*Image)
Metadata
¶
func (i *
Image
) Metadata() (
ImageMetadata
,
error
)
Metadata returns the image metadata (size, alpha channel, profile, EXIF rotation).
func (*Image)
Process
¶
func (i *
Image
) Process(o
Options
) ([]
byte
,
error
)
Process processes the image based on the given transformation options,
talking with libvips bindings accordingly and returning the resultant
image buffer.
func (*Image)
Resize
¶
func (i *
Image
) Resize(width, height
int
) ([]
byte
,
error
)
Resize resizes the image to fixed width and height.
func (*Image)
ResizeAndCrop
¶
func (i *
Image
) ResizeAndCrop(width, height
int
) ([]
byte
,
error
)
ResizeAndCrop resizes the image to fixed width and height with additional crop transformation.
func (*Image)
Rotate
¶
func (i *
Image
) Rotate(a
Angle
) ([]
byte
,
error
)
Rotate rotates the image by given angle degrees (0, 90, 180 or 270).
func (*Image)
Size
¶
func (i *
Image
) Size() (
ImageSize
,
error
)
Size returns the image size as form of width and height pixels.
func (*Image)
SmartCrop
¶
added in
v1.0.8
func (i *
Image
) SmartCrop(width, height
int
) ([]
byte
,
error
)
SmartCrop produces a thumbnail aiming at focus on the interesting part.
func (*Image)
Thumbnail
¶
func (i *
Image
) Thumbnail(pixels
int
) ([]
byte
,
error
)
Thumbnail creates a thumbnail of the image by the a given width by aspect ratio 4:4.
func (*Image)
Trim
¶
added in
v1.0.14
func (i *
Image
) Trim() ([]
byte
,
error
)
Trim removes the background from the picture. It can result in a 0x0 output
if the image is all background.
func (*Image)
Type
¶
func (i *
Image
) Type()
string
Type returns the image type format (jpeg, png, webp, tiff).
func (*Image)
Watermark
¶
func (i *
Image
) Watermark(w
Watermark
) ([]
byte
,
error
)
Watermark adds text as watermark on the given image.
func (*Image)
WatermarkImage
¶
added in
v1.0.8
func (i *
Image
) WatermarkImage(w
WatermarkImage
) ([]
byte
,
error
)
WatermarkImage adds image as watermark on the given image.
func (*Image)
Zoom
¶
func (i *
Image
) Zoom(factor
int
) ([]
byte
,
error
)
Zoom zooms the image by the given factor.
You should probably call Extract() before.
type
ImageMetadata
¶
type ImageMetadata struct {
Orientation
int
Channels
int
Alpha
bool
Profile
bool
Type
string
Space
string
Colourspace
string
Size
ImageSize
EXIF
EXIF
}
ImageMetadata represents the basic metadata fields
func
Metadata
¶
func Metadata(buf []
byte
) (
ImageMetadata
,
error
)
Metadata returns the image metadata (size, type, alpha channel, profile, EXIF orientation...).
type
ImageSize
¶
type ImageSize struct {
Width
int
Height
int
}
ImageSize represents the image width and height values
func
Size
¶
func Size(buf []
byte
) (
ImageSize
,
error
)
Size returns the image size by width and height pixels.
type
ImageType
¶
type ImageType
int
ImageType represents an image type value.
const (
// UNKNOWN represents an unknow image type value.
UNKNOWN
ImageType
=
iota
// JPEG represents the JPEG image type.
JPEG
// WEBP represents the WEBP image type.
WEBP
// PNG represents the PNG image type.
PNG
// TIFF represents the TIFF image type.
TIFF
// GIF represents the GIF image type.
GIF
// PDF represents the PDF type.
PDF
// SVG represents the SVG image type.
SVG
// MAGICK represents the libmagick compatible genetic image type.
MAGICK
// HEIF represents the HEIC/HEIF/HVEC image type
HEIF
// AVIF represents the AVIF image type.
AVIF
)
func
DetermineImageType
¶
func DetermineImageType(buf []
byte
)
ImageType
DetermineImageType determines the image type format (jpeg, png, webp or tiff)
type
Interpolator
¶
type Interpolator
int
Interpolator represents the image interpolation value.
const (
// Bicubic interpolation value.
Bicubic
Interpolator
=
iota
// Bilinear interpolation value.
Bilinear
// Nohalo interpolation value.
Nohalo
// Nearest neighbour interpolation value.
Nearest
)
func (Interpolator)
String
¶
func (i
Interpolator
) String()
string
type
Interpretation
¶
type Interpretation
int
Interpretation represents the image interpretation type.
See:
https://libvips.github.io/libvips/API/current/VipsImage.html#VipsInterpretation
const (
// InterpretationError points to the libvips interpretation error type.
InterpretationError
Interpretation
=
C
.
VIPS_INTERPRETATION_ERROR
// InterpretationMultiband points to its libvips interpretation equivalent type.
InterpretationMultiband
Interpretation
=
C
.
VIPS_INTERPRETATION_MULTIBAND
// InterpretationBW points to its libvips interpretation equivalent type.
InterpretationBW
Interpretation
=
C
.
VIPS_INTERPRETATION_B_W
// InterpretationCMYK points to its libvips interpretation equivalent type.
InterpretationCMYK
Interpretation
=
C
.
VIPS_INTERPRETATION_CMYK
// InterpretationRGB points to its libvips interpretation equivalent type.
InterpretationRGB
Interpretation
=
C
.
VIPS_INTERPRETATION_RGB
// InterpretationSRGB points to its libvips interpretation equivalent type.
InterpretationSRGB
Interpretation
=
C
.
VIPS_INTERPRETATION_sRGB
// InterpretationRGB16 points to its libvips interpretation equivalent type.
InterpretationRGB16
Interpretation
=
C
.
VIPS_INTERPRETATION_RGB16
// InterpretationGREY16 points to its libvips interpretation equivalent type.
InterpretationGREY16
Interpretation
=
C
.
VIPS_INTERPRETATION_GREY16
// InterpretationScRGB points to its libvips interpretation equivalent type.
InterpretationScRGB
Interpretation
=
C
.
VIPS_INTERPRETATION_scRGB
// InterpretationLAB points to its libvips interpretation equivalent type.
InterpretationLAB
Interpretation
=
C
.
VIPS_INTERPRETATION_LAB
// InterpretationXYZ points to its libvips interpretation equivalent type.
InterpretationXYZ
Interpretation
=
C
.
VIPS_INTERPRETATION_XYZ
)
func
ImageInterpretation
¶
func ImageInterpretation(buf []
byte
) (
Interpretation
,
error
)
ImageInterpretation returns the image interpretation type.
See:
https://libvips.github.io/libvips/API/current/VipsImage.html#VipsInterpretation
type
Options
¶
type Options struct {
Height
int
Width
int
AreaHeight
int
AreaWidth
int
Top
int
Left
int
Quality
int
Compression
int
Zoom
int
Crop
bool
SmartCrop
bool
// Deprecated, use: bimg.Options.Gravity = bimg.GravitySmart
Enlarge
bool
Embed
bool
Flip
bool
Flop
bool
Force
bool
NoAutoRotate
bool
NoProfile
bool
Interlace
bool
StripMetadata
bool
Trim
bool
Lossless
bool
Extend
Extend
Rotate
Angle
Background
Color
Gravity
Gravity
Watermark
Watermark
WatermarkImage
WatermarkImage
Type
ImageType
Interpolator
Interpolator
Interpretation
Interpretation
GaussianBlur
GaussianBlur
Sharpen
Sharpen
Threshold
float64
Gamma
float64
Brightness
float64
Contrast
float64
OutputICC
string
InputICC
string
Palette
bool
// Speed defines the AVIF encoders CPU effort. Valid values are:
// 0-8 for AVIF encoding.
// 0-9 for PNG encoding.
Speed
int
// contains filtered or unexported fields
}
Options represents the supported image transformation options.
type
Sharpen
¶
type Sharpen struct {
Radius
int
X1
float64
Y2
float64
Y3
float64
M1
float64
M2
float64
}
Sharpen represents the image sharp transformation options.
type
SupportedImageType
¶
added in
v1.0.7
type SupportedImageType struct {
Load
bool
Save
bool
}
SupportedImageType represents whether a type can be loaded and/or saved by
the current libvips compilation.
func
IsImageTypeSupportedByVips
¶
added in
v1.0.3
func IsImageTypeSupportedByVips(t
ImageType
)
SupportedImageType
IsImageTypeSupportedByVips returns true if the given image type
is supported by current libvips compilation.
type
VipsMemoryInfo
¶
type VipsMemoryInfo struct {
Memory
int64
MemoryHighwater
int64
Allocations
int64
}
VipsMemoryInfo represents the memory stats provided by libvips.
func
VipsMemory
¶
func VipsMemory()
VipsMemoryInfo
VipsMemory gets memory info stats from libvips (cache size, memory allocs...)
type
Watermark
¶
type Watermark struct {
Width
int
DPI
int
Margin
int
Opacity
float32
NoReplicate
bool
Text
string
Font
string
Background
Color
}
Watermark represents the text-based watermark supported options.
type
WatermarkImage
¶
added in
v1.0.8
type WatermarkImage struct {
Left
int
Top
int
Buf     []
byte
Opacity
float32
}
WatermarkImage represents the image-based watermark supported options.