# corona10/goimagehash

> Source: https://pkg.go.dev/github.com/corona10/goimagehash
> Fetched: 2026-02-01T11:48:49.057225+00:00
> Content-Hash: b6375d2d2c0ee4d0
> Type: html

---

### Index ¶

  * type ExtImageHash
  *     * func ExtAverageHash(img image.Image, width, height int) (*ExtImageHash, error)
    * func ExtDifferenceHash(img image.Image, width, height int) (*ExtImageHash, error)
    * func ExtImageHashFromString(s string) (*ExtImageHash, error)deprecated
    * func ExtPerceptionHash(img image.Image, width, height int) (*ExtImageHash, error)
    * func LoadExtImageHash(b io.Reader) (*ExtImageHash, error)
    * func NewExtImageHash(hash []uint64, kind Kind, bits int) *ExtImageHash
  *     * func (h *ExtImageHash) Bits() int
    * func (h *ExtImageHash) Distance(other *ExtImageHash) (int, error)
    * func (h *ExtImageHash) Dump(w io.Writer) error
    * func (h *ExtImageHash) GetHash() []uint64
    * func (h *ExtImageHash) GetKind() Kind
    * func (h *ExtImageHash) ToString() string
  * type ImageHash
  *     * func AverageHash(img image.Image) (*ImageHash, error)
    * func DifferenceHash(img image.Image) (*ImageHash, error)
    * func ImageHashFromString(s string) (*ImageHash, error)deprecated
    * func LoadImageHash(b io.Reader) (*ImageHash, error)
    * func NewImageHash(hash uint64, kind Kind) *ImageHash
    * func PerceptionHash(img image.Image) (*ImageHash, error)
  *     * func (h *ImageHash) Bits() int
    * func (h *ImageHash) Distance(other *ImageHash) (int, error)
    * func (h *ImageHash) Dump(w io.Writer) error
    * func (h *ImageHash) GetHash() uint64
    * func (h *ImageHash) GetKind() Kind
    * func (h *ImageHash) ToString() string
  * type Kind



### Constants ¶

This section is empty.

### Variables ¶

This section is empty.

### Functions ¶

This section is empty.

### Types ¶

####  type [ExtImageHash](https://github.com/corona10/goimagehash/blob/v1.1.0/imagehash.go#L26) ¶ added in v0.3.0
    
    
    type ExtImageHash struct {
    	// contains filtered or unexported fields
    }

ExtImageHash is a struct of big hash computation. 

####  func [ExtAverageHash](https://github.com/corona10/goimagehash/blob/v1.1.0/hashcompute.go#L139) ¶ added in v1.0.0
    
    
    func ExtAverageHash(img [image](/image).[Image](/image#Image), width, height [int](/builtin#int)) (*ExtImageHash, [error](/builtin#error))

ExtAverageHash function returns ahash of which the size can be set larger than uint64 Support 64bits ahash (width=8, height=8) and 256bits ahash (width=16, height=16) 

####  func [ExtDifferenceHash](https://github.com/corona10/goimagehash/blob/v1.1.0/hashcompute.go#L169) ¶ added in v1.0.0
    
    
    func ExtDifferenceHash(img [image](/image).[Image](/image#Image), width, height [int](/builtin#int)) (*ExtImageHash, [error](/builtin#error))

ExtDifferenceHash function returns dhash of which the size can be set larger than uint64 Support 64bits dhash (width=8, height=8) and 256bits dhash (width=16, height=16) 

####  func [ExtImageHashFromString](https://github.com/corona10/goimagehash/blob/v1.1.0/imagehash.go#L236) deprecated added in v0.3.0
    
    
    func ExtImageHashFromString(s [string](/builtin#string)) (*ExtImageHash, [error](/builtin#error))

ExtImageHashFromString returns a big hash from a hex representation 

Deprecated: Use goimagehash.LoadExtImageHash instead. 

####  func [ExtPerceptionHash](https://github.com/corona10/goimagehash/blob/v1.1.0/hashcompute.go#L106) ¶ added in v1.0.0
    
    
    func ExtPerceptionHash(img [image](/image).[Image](/image#Image), width, height [int](/builtin#int)) (*ExtImageHash, [error](/builtin#error))

ExtPerceptionHash function returns phash of which the size can be set larger than uint64 Some variable name refer to <https://github.com/JohannesBuchner/imagehash/blob/master/imagehash/__init__.py> Support 64bits phash (width=8, height=8) and 256bits phash (width=16, height=16) Important: width * height should be the power of 2 

####  func [LoadExtImageHash](https://github.com/corona10/goimagehash/blob/v1.1.0/imagehash.go#L216) ¶ added in v1.0.0
    
    
    func LoadExtImageHash(b [io](/io).[Reader](/io#Reader)) (*ExtImageHash, [error](/builtin#error))

LoadExtImageHash method loads a ExtImageHash from io.Reader. 

####  func [NewExtImageHash](https://github.com/corona10/goimagehash/blob/v1.1.0/imagehash.go#L155) ¶ added in v0.3.0
    
    
    func NewExtImageHash(hash [][uint64](/builtin#uint64), kind Kind, bits [int](/builtin#int)) *ExtImageHash

NewExtImageHash function creates a new big hash 

####  func (*ExtImageHash) [Bits](https://github.com/corona10/goimagehash/blob/v1.1.0/imagehash.go#L160) ¶ added in v1.0.0
    
    
    func (h *ExtImageHash) Bits() [int](/builtin#int)

Bits method returns an actual hash bit size 

####  func (*ExtImageHash) [Distance](https://github.com/corona10/goimagehash/blob/v1.1.0/imagehash.go#L165) ¶ added in v0.3.0
    
    
    func (h *ExtImageHash) Distance(other *ExtImageHash) ([int](/builtin#int), [error](/builtin#error))

Distance method returns a distance between two big hashes 

####  func (*ExtImageHash) [Dump](https://github.com/corona10/goimagehash/blob/v1.1.0/imagehash.go#L201) ¶ added in v1.0.0
    
    
    func (h *ExtImageHash) Dump(w [io](/io).[Writer](/io#Writer)) [error](/builtin#error)

Dump method writes a binary serialization into w io.Writer. 

####  func (*ExtImageHash) [GetHash](https://github.com/corona10/goimagehash/blob/v1.1.0/imagehash.go#L191) ¶ added in v0.3.0
    
    
    func (h *ExtImageHash) GetHash() [][uint64](/builtin#uint64)

GetHash method returns a big hash value 

####  func (*ExtImageHash) [GetKind](https://github.com/corona10/goimagehash/blob/v1.1.0/imagehash.go#L196) ¶ added in v0.3.0
    
    
    func (h *ExtImageHash) GetKind() Kind

GetKind method returns a kind of big hash 

####  func (*ExtImageHash) [ToString](https://github.com/corona10/goimagehash/blob/v1.1.0/imagehash.go#L273) ¶ added in v0.3.0
    
    
    func (h *ExtImageHash) ToString() [string](/builtin#string)

ToString returns a hex representation of big hash 

####  type [ImageHash](https://github.com/corona10/goimagehash/blob/v1.1.0/imagehash.go#L20) ¶
    
    
    type ImageHash struct {
    	// contains filtered or unexported fields
    }

ImageHash is a struct of hash computation. 

####  func [AverageHash](https://github.com/corona10/goimagehash/blob/v1.1.0/hashcompute.go#L20) ¶
    
    
    func AverageHash(img [image](/image).[Image](/image#Image)) (*ImageHash, [error](/builtin#error))

AverageHash fuction returns a hash computation of average hash. Implementation follows <http://www.hackerfactor.com/blog/index.php?/archives/432-Looks-Like-It.html>

####  func [DifferenceHash](https://github.com/corona10/goimagehash/blob/v1.1.0/hashcompute.go#L44) ¶
    
    
    func DifferenceHash(img [image](/image).[Image](/image#Image)) (*ImageHash, [error](/builtin#error))

DifferenceHash function returns a hash computation of difference hash. Implementation follows <http://www.hackerfactor.com/blog/?/archives/529-Kind-of-Like-That.html>

####  func [ImageHashFromString](https://github.com/corona10/goimagehash/blob/v1.1.0/imagehash.go#L116) deprecated
    
    
    func ImageHashFromString(s [string](/builtin#string)) (*ImageHash, [error](/builtin#error))

ImageHashFromString returns an image hash from a hex representation 

Deprecated: Use goimagehash.LoadImageHash instead. 

####  func [LoadImageHash](https://github.com/corona10/goimagehash/blob/v1.1.0/imagehash.go#L99) ¶ added in v1.0.0
    
    
    func LoadImageHash(b [io](/io).[Reader](/io#Reader)) (*ImageHash, [error](/builtin#error))

LoadImageHash method loads a ImageHash from io.Reader. 

####  func [NewImageHash](https://github.com/corona10/goimagehash/blob/v1.1.0/imagehash.go#L46) ¶
    
    
    func NewImageHash(hash [uint64](/builtin#uint64), kind Kind) *ImageHash

NewImageHash function creates a new image hash. 

####  func [PerceptionHash](https://github.com/corona10/goimagehash/blob/v1.1.0/hashcompute.go#L68) ¶
    
    
    func PerceptionHash(img [image](/image).[Image](/image#Image)) (*ImageHash, [error](/builtin#error))

PerceptionHash function returns a hash computation of phash. Implementation follows <http://www.hackerfactor.com/blog/index.php?/archives/432-Looks-Like-It.html>

####  func (*ImageHash) [Bits](https://github.com/corona10/goimagehash/blob/v1.1.0/imagehash.go#L51) ¶ added in v1.0.0
    
    
    func (h *ImageHash) Bits() [int](/builtin#int)

Bits method returns an actual hash bit size 

####  func (*ImageHash) [Distance](https://github.com/corona10/goimagehash/blob/v1.1.0/imagehash.go#L56) ¶
    
    
    func (h *ImageHash) Distance(other *ImageHash) ([int](/builtin#int), [error](/builtin#error))

Distance method returns a distance between two hashes. 

####  func (*ImageHash) [Dump](https://github.com/corona10/goimagehash/blob/v1.1.0/imagehash.go#L85) ¶ added in v1.0.0
    
    
    func (h *ImageHash) Dump(w [io](/io).[Writer](/io#Writer)) [error](/builtin#error)

Dump method writes a binary serialization into w io.Writer. 

####  func (*ImageHash) [GetHash](https://github.com/corona10/goimagehash/blob/v1.1.0/imagehash.go#L69) ¶
    
    
    func (h *ImageHash) GetHash() [uint64](/builtin#uint64)

GetHash method returns a 64bits hash value. 

####  func (*ImageHash) [GetKind](https://github.com/corona10/goimagehash/blob/v1.1.0/imagehash.go#L74) ¶
    
    
    func (h *ImageHash) GetKind() Kind

GetKind method returns a kind of image hash. 

####  func (*ImageHash) [ToString](https://github.com/corona10/goimagehash/blob/v1.1.0/imagehash.go#L139) ¶
    
    
    func (h *ImageHash) ToString() [string](/builtin#string)

ToString returns a hex representation of the hash 

####  type [Kind](https://github.com/corona10/goimagehash/blob/v1.1.0/imagehash.go#L17) ¶
    
    
    type Kind [int](/builtin#int)

Kind describes the kinds of hash. 
    
    
    const (
    	// Unknown is a enum value of the unknown hash.
    	Unknown Kind = [iota](/builtin#iota)
    	// AHash is a enum value of the average hash.
    	AHash
    	//PHash is a enum value of the perceptual hash.
    	PHash
    	// DHash is a enum value of the difference hash.
    	DHash
    	// WHash is a enum value of the wavelet hash.
    	WHash
    )
  *[↑]: Back to Top
  *[v]: View this template
  *[t]: Discuss this template
  *[e]: Edit this template
