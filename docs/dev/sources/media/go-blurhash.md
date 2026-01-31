# go-blurhash

> Source: https://pkg.go.dev/github.com/bbrks/go-blurhash
> Fetched: 2026-01-30T23:53:56.156874+00:00
> Content-Hash: 4dcfa80ecb6616b3
> Type: html

---

Overview

¶

Package blurhash provides encoding and decoding of Blurhash image placeholders.

Blurhash is an algorithm that encodes an image into a short ASCII string
representing a gradient of colors. When decoded, this string produces a
blurred placeholder that approximates the original image's colors and structure.

For simple one-off operations, use the package-level

Encode

and

Decode

functions.
For batch processing with reduced allocations, use the reusable

Encoder

and

Decoder

types.

Index

¶

Variables

func Components(hash string) (x, y int, err error)

func Decode(hash string, width, height int, punch int) (image.Image, error)

func DecodeDraw(dst draw.Image, hash string, punch float64) error

func Encode(xComponents, yComponents int, img image.Image) (string, error)

type Decoder

func NewDecoder() *Decoder

func (d *Decoder) Decode(hash string, width, height, punch int) (image.Image, error)

func (d *Decoder) DecodeDraw(dst draw.Image, hash string, punch float64) error

type Encoder

func NewEncoder() *Encoder

func (e *Encoder) Encode(xComponents, yComponents int, img image.Image) (string, error)

Constants

¶

This section is empty.

Variables

¶

View Source

var (

// ErrInvalidComponents is returned when components passed to Encode are invalid.

ErrInvalidComponents =

errors

.

New

("blurhash: must have between 1 and 9 components")

// ErrInvalidHash is returned when the library encounters a hash it can't recognise.

ErrInvalidHash =

errors

.

New

("blurhash: invalid hash")

// ErrInvalidDimensions is returned when width or height is invalid.

ErrInvalidDimensions =

errors

.

New

("blurhash: width and height must be positive")
)

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

) (

string

,

error

)

Encode returns the blurhash for the given image.

Types

¶

type

Decoder

¶

added in

v1.2.0

type Decoder struct {

// contains filtered or unexported fields

}

Decoder is a reusable blurhash decoder that minimizes allocations
by reusing internal buffers across decode operations.

A Decoder is safe for sequential use but not for concurrent use.
For concurrent workloads, use a sync.Pool of Decoders.

The zero value is ready to use.

func

NewDecoder

¶

added in

v1.2.0

func NewDecoder() *

Decoder

NewDecoder creates a new reusable Decoder.
Buffers are allocated lazily on first use and grown as needed.

func (*Decoder)

Decode

¶

added in

v1.2.0

func (d *

Decoder

) Decode(hash

string

, width, height, punch

int

) (

image

.

Image

,

error

)

Decode decodes a blurhash to a new NRGBA image.
Internal buffers are reused across calls when possible.

func (*Decoder)

DecodeDraw

¶

added in

v1.2.0

func (d *

Decoder

) DecodeDraw(dst

draw

.

Image

, hash

string

, punch

float64

)

error

DecodeDraw decodes a blurhash into an existing image.
Internal buffers are reused across calls when possible.

type

Encoder

¶

added in

v1.2.0

type Encoder struct {

// contains filtered or unexported fields

}

Encoder is a reusable blurhash encoder that minimizes allocations
by reusing internal buffers across encode operations.

An Encoder is safe for sequential use but not for concurrent use.
For concurrent workloads, use a sync.Pool of Encoders.

The zero value is ready to use.

func

NewEncoder

¶

added in

v1.2.0

func NewEncoder() *

Encoder

NewEncoder creates a new reusable Encoder.
Buffers are allocated lazily on first use and grown as needed.

func (*Encoder)

Encode

¶

added in

v1.2.0

func (e *

Encoder

) Encode(xComponents, yComponents

int

, img

image

.

Image

) (

string

,

error

)

Encode returns the blurhash for the given image.
Internal buffers are reused across calls when possible.