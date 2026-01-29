# go-blurhash

> Auto-fetched from [https://pkg.go.dev/github.com/bbrks/go-blurhash](https://pkg.go.dev/github.com/bbrks/go-blurhash)
> Last Updated: 2026-01-29T20:14:48.595501+00:00

---

Index
¶
Variables
func Components(hash string) (x, y int, err error)
func Decode(hash string, width, height int, punch int) (image.Image, error)
func DecodeDraw(dst draw.Image, hash string, punch float64) error
func Encode(xComponents, yComponents int, img image.Image) (hash string, err error)
Constants
¶
This section is empty.
Variables
¶
View Source
var ErrInvalidComponents =
errors
.
New
("blurhash: must have between 1 and 9 components")
ErrInvalidComponents is returned when components passed to Encode are invalid.
View Source
var ErrInvalidHash =
errors
.
New
("blurhash: invalid hash")
ErrInvalidHash is returned when the library encounters a hash it can't recognise.
Functions
¶
func
Components
¶
func Components(hash
string
) (x, y
int
, err
error
)
Components returns the X and Y components of a blurhash.
func
Decode
¶
func Decode(hash
string
, width, height
int
, punch
int
) (
image
.
Image
,
error
)
Decode returns an NRGBA image of the given hash with the given size.
func
DecodeDraw
¶
func DecodeDraw(dst
draw
.
Image
, hash
string
, punch
float64
)
error
DecodeDraw decodes the given hash into the given image.
func
Encode
¶
func Encode(xComponents, yComponents
int
, img
image
.
Image
) (hash
string
, err
error
)
Encode returns the blurhash for the given image.
Types
¶
This section is empty.