# corona10/goimagehash

> Source: https://pkg.go.dev/github.com/corona10/goimagehash
> Fetched: 2026-01-30T23:54:07.849051+00:00
> Content-Hash: 7de82478ff8e3ca6
> Type: html

---

Index

¶

type ExtImageHash

func ExtAverageHash(img image.Image, width, height int) (*ExtImageHash, error)

func ExtDifferenceHash(img image.Image, width, height int) (*ExtImageHash, error)

func ExtImageHashFromString(s string) (*ExtImageHash, error)

deprecated

func ExtPerceptionHash(img image.Image, width, height int) (*ExtImageHash, error)

func LoadExtImageHash(b io.Reader) (*ExtImageHash, error)

func NewExtImageHash(hash []uint64, kind Kind, bits int) *ExtImageHash

func (h *ExtImageHash) Bits() int

func (h *ExtImageHash) Distance(other *ExtImageHash) (int, error)

func (h *ExtImageHash) Dump(w io.Writer) error

func (h *ExtImageHash) GetHash() []uint64

func (h *ExtImageHash) GetKind() Kind

func (h *ExtImageHash) ToString() string

type ImageHash

func AverageHash(img image.Image) (*ImageHash, error)

func DifferenceHash(img image.Image) (*ImageHash, error)

func ImageHashFromString(s string) (*ImageHash, error)

deprecated

func LoadImageHash(b io.Reader) (*ImageHash, error)

func NewImageHash(hash uint64, kind Kind) *ImageHash

func PerceptionHash(img image.Image) (*ImageHash, error)

func (h *ImageHash) Bits() int

func (h *ImageHash) Distance(other *ImageHash) (int, error)

func (h *ImageHash) Dump(w io.Writer) error

func (h *ImageHash) GetHash() uint64

func (h *ImageHash) GetKind() Kind

func (h *ImageHash) ToString() string

type Kind

Constants

¶

This section is empty.

Variables

¶

This section is empty.

Functions

¶

This section is empty.

Types

¶

type

ExtImageHash

¶

added in

v0.3.0

type ExtImageHash struct {

// contains filtered or unexported fields

}

ExtImageHash is a struct of big hash computation.

func

ExtAverageHash

¶

added in

v1.0.0

func ExtAverageHash(img

image

.

Image

, width, height

int

) (*

ExtImageHash

,

error

)

ExtAverageHash function returns ahash of which the size can be set larger than uint64
Support 64bits ahash (width=8, height=8) and 256bits ahash (width=16, height=16)

func

ExtDifferenceHash

¶

added in

v1.0.0

func ExtDifferenceHash(img

image

.

Image

, width, height

int

) (*

ExtImageHash

,

error

)

ExtDifferenceHash function returns dhash of which the size can be set larger than uint64
Support 64bits dhash (width=8, height=8) and 256bits dhash (width=16, height=16)

func

ExtImageHashFromString

deprecated

added in

v0.3.0

func ExtImageHashFromString(s

string

) (*

ExtImageHash

,

error

)

ExtImageHashFromString returns a big hash from a hex representation

Deprecated: Use goimagehash.LoadExtImageHash instead.

func

ExtPerceptionHash

¶

added in

v1.0.0

func ExtPerceptionHash(img

image

.

Image

, width, height

int

) (*

ExtImageHash

,

error

)

ExtPerceptionHash function returns phash of which the size can be set larger than uint64
Some variable name refer to

https://github.com/JohannesBuchner/imagehash/blob/master/imagehash/__init__.py

Support 64bits phash (width=8, height=8) and 256bits phash (width=16, height=16)
Important: width * height should be the power of 2

func

LoadExtImageHash

¶

added in

v1.0.0

func LoadExtImageHash(b

io

.

Reader

) (*

ExtImageHash

,

error

)

LoadExtImageHash method loads a ExtImageHash from io.Reader.

func

NewExtImageHash

¶

added in

v0.3.0

func NewExtImageHash(hash []

uint64

, kind

Kind

, bits

int

) *

ExtImageHash

NewExtImageHash function creates a new big hash

func (*ExtImageHash)

Bits

¶

added in

v1.0.0

func (h *

ExtImageHash

) Bits()

int

Bits method returns an actual hash bit size

func (*ExtImageHash)

Distance

¶

added in

v0.3.0

func (h *

ExtImageHash

) Distance(other *

ExtImageHash

) (

int

,

error

)

Distance method returns a distance between two big hashes

func (*ExtImageHash)

Dump

¶

added in

v1.0.0

func (h *

ExtImageHash

) Dump(w

io

.

Writer

)

error

Dump method writes a binary serialization into w io.Writer.

func (*ExtImageHash)

GetHash

¶

added in

v0.3.0

func (h *

ExtImageHash

) GetHash() []

uint64

GetHash method returns a big hash value

func (*ExtImageHash)

GetKind

¶

added in

v0.3.0

func (h *

ExtImageHash

) GetKind()

Kind

GetKind method returns a kind of big hash

func (*ExtImageHash)

ToString

¶

added in

v0.3.0

func (h *

ExtImageHash

) ToString()

string

ToString returns a hex representation of big hash

type

ImageHash

¶

type ImageHash struct {

// contains filtered or unexported fields

}

ImageHash is a struct of hash computation.

func

AverageHash

¶

func AverageHash(img

image

.

Image

) (*

ImageHash

,

error

)

AverageHash fuction returns a hash computation of average hash.
Implementation follows

http://www.hackerfactor.com/blog/index.php?/archives/432-Looks-Like-It.html

func

DifferenceHash

¶

func DifferenceHash(img

image

.

Image

) (*

ImageHash

,

error

)

DifferenceHash function returns a hash computation of difference hash.
Implementation follows

http://www.hackerfactor.com/blog/?/archives/529-Kind-of-Like-That.html

func

ImageHashFromString

deprecated

func ImageHashFromString(s

string

) (*

ImageHash

,

error

)

ImageHashFromString returns an image hash from a hex representation

Deprecated: Use goimagehash.LoadImageHash instead.

func

LoadImageHash

¶

added in

v1.0.0

func LoadImageHash(b

io

.

Reader

) (*

ImageHash

,

error

)

LoadImageHash method loads a ImageHash from io.Reader.

func

NewImageHash

¶

func NewImageHash(hash

uint64

, kind

Kind

) *

ImageHash

NewImageHash function creates a new image hash.

func

PerceptionHash

¶

func PerceptionHash(img

image

.

Image

) (*

ImageHash

,

error

)

PerceptionHash function returns a hash computation of phash.
Implementation follows

http://www.hackerfactor.com/blog/index.php?/archives/432-Looks-Like-It.html

func (*ImageHash)

Bits

¶

added in

v1.0.0

func (h *

ImageHash

) Bits()

int

Bits method returns an actual hash bit size

func (*ImageHash)

Distance

¶

func (h *

ImageHash

) Distance(other *

ImageHash

) (

int

,

error

)

Distance method returns a distance between two hashes.

func (*ImageHash)

Dump

¶

added in

v1.0.0

func (h *

ImageHash

) Dump(w

io

.

Writer

)

error

Dump method writes a binary serialization into w io.Writer.

func (*ImageHash)

GetHash

¶

func (h *

ImageHash

) GetHash()

uint64

GetHash method returns a 64bits hash value.

func (*ImageHash)

GetKind

¶

func (h *

ImageHash

) GetKind()

Kind

GetKind method returns a kind of image hash.

func (*ImageHash)

ToString

¶

func (h *

ImageHash

) ToString()

string

ToString returns a hex representation of the hash

type

Kind

¶

type Kind

int

Kind describes the kinds of hash.

const (

// Unknown is a enum value of the unknown hash.

Unknown

Kind

=

iota

// AHash is a enum value of the average hash.

AHash

//PHash is a enum value of the perceptual hash.

PHash

// DHash is a enum value of the difference hash.

DHash

// WHash is a enum value of the wavelet hash.

WHash
)