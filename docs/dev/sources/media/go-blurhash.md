# go-blurhash

> Source: https://pkg.go.dev/github.com/bbrks/go-blurhash
> Fetched: 2026-02-01T11:48:36.898228+00:00
> Content-Hash: 6f35fec4f26d30a2
> Type: html

---

### Overview ¶

Package blurhash provides encoding and decoding of Blurhash image placeholders. 

Blurhash is an algorithm that encodes an image into a short ASCII string representing a gradient of colors. When decoded, this string produces a blurred placeholder that approximates the original image's colors and structure. 

For simple one-off operations, use the package-level Encode and Decode functions. For batch processing with reduced allocations, use the reusable Encoder and Decoder types. 

### Index ¶

  * Variables
  * func Components(hash string) (x, y int, err error)
  * func Decode(hash string, width, height int, punch int) (image.Image, error)
  * func DecodeDraw(dst draw.Image, hash string, punch float64) error
  * func Encode(xComponents, yComponents int, img image.Image) (string, error)
  * type Decoder
  *     * func NewDecoder() *Decoder
  *     * func (d *Decoder) Decode(hash string, width, height, punch int) (image.Image, error)
    * func (d *Decoder) DecodeDraw(dst draw.Image, hash string, punch float64) error
  * type Encoder
  *     * func NewEncoder() *Encoder
  *     * func (e *Encoder) Encode(xComponents, yComponents int, img image.Image) (string, error)



### Constants ¶

This section is empty.

### Variables ¶

[View Source](https://github.com/bbrks/go-blurhash/blob/v1.2.0/error.go#L7)
    
    
    var (
    	// ErrInvalidComponents is returned when components passed to Encode are invalid.
    	ErrInvalidComponents = [errors](/errors).[New](/errors#New)("blurhash: must have between 1 and 9 components")
    	// ErrInvalidHash is returned when the library encounters a hash it can't recognise.
    	ErrInvalidHash = [errors](/errors).[New](/errors#New)("blurhash: invalid hash")
    	// ErrInvalidDimensions is returned when width or height is invalid.
    	ErrInvalidDimensions = [errors](/errors).[New](/errors#New)("blurhash: width and height must be positive")
    )

### Functions ¶

####  func [Components](https://github.com/bbrks/go-blurhash/blob/v1.2.0/decode.go#L14) ¶
    
    
    func Components(hash [string](/builtin#string)) (x, y [int](/builtin#int), err [error](/builtin#error))

Components returns the X and Y components of a blurhash. 

####  func [Decode](https://github.com/bbrks/go-blurhash/blob/v1.2.0/decode.go#L177) ¶
    
    
    func Decode(hash [string](/builtin#string), width, height [int](/builtin#int), punch [int](/builtin#int)) ([image](/image).[Image](/image#Image), [error](/builtin#error))

Decode returns an NRGBA image of the given hash with the given size. 

####  func [DecodeDraw](https://github.com/bbrks/go-blurhash/blob/v1.2.0/decode.go#L183) ¶
    
    
    func DecodeDraw(dst [draw](/image/draw).[Image](/image/draw#Image), hash [string](/builtin#string), punch [float64](/builtin#float64)) [error](/builtin#error)

DecodeDraw decodes the given hash into the given image. 

####  func [Encode](https://github.com/bbrks/go-blurhash/blob/v1.2.0/encode.go#L176) ¶
    
    
    func Encode(xComponents, yComponents [int](/builtin#int), img [image](/image).[Image](/image#Image)) ([string](/builtin#string), [error](/builtin#error))

Encode returns the blurhash for the given image. 

### Types ¶

####  type [Decoder](https://github.com/bbrks/go-blurhash/blob/v1.2.0/decode.go#L47) ¶ added in v1.2.0
    
    
    type Decoder struct {
    	// contains filtered or unexported fields
    }

Decoder is a reusable blurhash decoder that minimizes allocations by reusing internal buffers across decode operations. 

A Decoder is safe for sequential use but not for concurrent use. For concurrent workloads, use a sync.Pool of Decoders. 

The zero value is ready to use. 

####  func [NewDecoder](https://github.com/bbrks/go-blurhash/blob/v1.2.0/decode.go#L54) ¶ added in v1.2.0
    
    
    func NewDecoder() *Decoder

NewDecoder creates a new reusable Decoder. Buffers are allocated lazily on first use and grown as needed. 

####  func (*Decoder) [Decode](https://github.com/bbrks/go-blurhash/blob/v1.2.0/decode.go#L60) ¶ added in v1.2.0
    
    
    func (d *Decoder) Decode(hash [string](/builtin#string), width, height, punch [int](/builtin#int)) ([image](/image).[Image](/image#Image), [error](/builtin#error))

Decode decodes a blurhash to a new NRGBA image. Internal buffers are reused across calls when possible. 

####  func (*Decoder) [DecodeDraw](https://github.com/bbrks/go-blurhash/blob/v1.2.0/decode.go#L73) ¶ added in v1.2.0
    
    
    func (d *Decoder) DecodeDraw(dst [draw](/image/draw).[Image](/image/draw#Image), hash [string](/builtin#string), punch [float64](/builtin#float64)) [error](/builtin#error)

DecodeDraw decodes a blurhash into an existing image. Internal buffers are reused across calls when possible. 

####  type [Encoder](https://github.com/bbrks/go-blurhash/blob/v1.2.0/encode.go#L33) ¶ added in v1.2.0
    
    
    type Encoder struct {
    	// contains filtered or unexported fields
    }

Encoder is a reusable blurhash encoder that minimizes allocations by reusing internal buffers across encode operations. 

An Encoder is safe for sequential use but not for concurrent use. For concurrent workloads, use a sync.Pool of Encoders. 

The zero value is ready to use. 

####  func [NewEncoder](https://github.com/bbrks/go-blurhash/blob/v1.2.0/encode.go#L42) ¶ added in v1.2.0
    
    
    func NewEncoder() *Encoder

NewEncoder creates a new reusable Encoder. Buffers are allocated lazily on first use and grown as needed. 

####  func (*Encoder) [Encode](https://github.com/bbrks/go-blurhash/blob/v1.2.0/encode.go#L48) ¶ added in v1.2.0
    
    
    func (e *Encoder) Encode(xComponents, yComponents [int](/builtin#int), img [image](/image).[Image](/image#Image)) ([string](/builtin#string), [error](/builtin#error))

Encode returns the blurhash for the given image. Internal buffers are reused across calls when possible. 
  *[↑]: Back to Top
  *[v]: View this template
  *[t]: Discuss this template
  *[e]: Edit this template
