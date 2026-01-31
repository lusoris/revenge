# go-faster/jx (JSON)

> Source: https://pkg.go.dev/github.com/go-faster/jx
> Fetched: 2026-01-30T23:49:12.861368+00:00
> Content-Hash: 0e4c6d9ad00b1e6d
> Type: html

---

Overview

¶

Package jx implements

RFC 7159

json encoding and decoding.

Example

¶

package main

import (
	"fmt"

	"github.com/go-faster/jx"
)

func main() {
	var e jx.Encoder
	e.Obj(func(e *jx.Encoder) {
		e.FieldStart("data")
		e.Base64([]byte("hello"))
	})
	fmt.Println(e)

	if err := jx.DecodeBytes(e.Bytes()).Obj(func(d *jx.Decoder, key string) error {
		v, err := d.Base64()
		fmt.Printf("%s: %s\n", key, v)
		return err
	}); err != nil {
		panic(err)
	}
}

Output:

{"data":"aGVsbG8="}
data: hello

Share

Format

Run

Index

¶

func PutDecoder(d *Decoder)

func PutEncoder(e *Encoder)

func PutWriter(e *Writer)

func Valid(data []byte) bool

type ArrIter

func (i *ArrIter) Err() error

func (i *ArrIter) Next() bool

type Decoder

func Decode(reader io.Reader, bufSize int) *Decoder

func DecodeBytes(input []byte) *Decoder

func DecodeStr(input string) *Decoder

func GetDecoder() *Decoder

func (d *Decoder) Arr(f func(d *Decoder) error) error

func (d *Decoder) ArrIter() (ArrIter, error)

func (d *Decoder) Base64() ([]byte, error)

func (d *Decoder) Base64Append(b []byte) ([]byte, error)

func (d *Decoder) BigFloat() (*big.Float, error)

func (d *Decoder) BigInt() (*big.Int, error)

func (d *Decoder) Bool() (bool, error)

func (d *Decoder) Capture(f func(d *Decoder) error) error

func (d *Decoder) Elem() (ok bool, err error)

func (d *Decoder) Float32() (float32, error)

func (d *Decoder) Float64() (float64, error)

func (d *Decoder) Int() (int, error)

func (d *Decoder) Int16() (int16, error)

func (d *Decoder) Int32() (int32, error)

func (d *Decoder) Int64() (int64, error)

func (d *Decoder) Int8() (int8, error)

func (d *Decoder) Next() Type

func (d *Decoder) Null() error

func (d *Decoder) Num() (Num, error)

func (d *Decoder) NumAppend(v Num) (Num, error)

func (d *Decoder) Obj(f func(d *Decoder, key string) error) error

func (d *Decoder) ObjBytes(f func(d *Decoder, key []byte) error) error

func (d *Decoder) ObjIter() (ObjIter, error)

func (d *Decoder) Raw() (Raw, error)

func (d *Decoder) RawAppend(buf Raw) (Raw, error)

func (d *Decoder) Reset(reader io.Reader)

func (d *Decoder) ResetBytes(input []byte)

func (d *Decoder) Skip() error

func (d *Decoder) Str() (string, error)

func (d *Decoder) StrAppend(b []byte) ([]byte, error)

func (d *Decoder) StrBytes() ([]byte, error)

func (d *Decoder) UInt() (uint, error)

func (d *Decoder) UInt16() (uint16, error)

func (d *Decoder) UInt32() (uint32, error)

func (d *Decoder) UInt64() (uint64, error)

func (d *Decoder) UInt8() (uint8, error)

func (d *Decoder) Validate() error

type Encoder

func GetEncoder() *Encoder

func NewStreamingEncoder(w io.Writer, bufSize int) *Encoder

func (e *Encoder) Arr(f func(e *Encoder)) (fail bool)

func (e *Encoder) ArrEmpty() bool

func (e *Encoder) ArrEnd() bool

func (e *Encoder) ArrStart() (fail bool)

func (e *Encoder) Base64(data []byte) bool

func (e *Encoder) Bool(v bool) bool

func (e *Encoder) ByteStr(v []byte) bool

func (e *Encoder) ByteStrEscape(v []byte) bool

func (e Encoder) Bytes() []byte

func (e *Encoder) Close() error

func (e *Encoder) Field(name string, f func(e *Encoder)) (fail bool)

func (e *Encoder) FieldStart(field string) (fail bool)

func (e *Encoder) Float32(v float32) bool

func (e *Encoder) Float64(v float64) bool

func (e *Encoder) Grow(n int)

func (e *Encoder) Int(v int) bool

func (e *Encoder) Int16(v int16) bool

func (e *Encoder) Int32(v int32) bool

func (e *Encoder) Int64(v int64) bool

func (e *Encoder) Int8(v int8) bool

func (e *Encoder) Null() bool

func (e *Encoder) Num(v Num) bool

func (e *Encoder) Obj(f func(e *Encoder)) (fail bool)

func (e *Encoder) ObjEmpty() bool

func (e *Encoder) ObjEnd() bool

func (e *Encoder) ObjStart() (fail bool)

func (e *Encoder) Raw(b []byte) bool

func (e *Encoder) RawStr(v string) bool

func (e *Encoder) Reset()

func (e *Encoder) ResetWriter(out io.Writer)

func (e *Encoder) SetBytes(buf []byte)

func (e *Encoder) SetIdent(n int)

func (e *Encoder) Str(v string) bool

func (e *Encoder) StrEscape(v string) bool

func (e Encoder) String() string

func (e *Encoder) UInt(v uint) bool

func (e *Encoder) UInt16(v uint16) bool

func (e *Encoder) UInt32(v uint32) bool

func (e *Encoder) UInt64(v uint64) bool

func (e *Encoder) UInt8(v uint8) bool

func (e *Encoder) Write(p []byte) (n int, err error)

func (e *Encoder) WriteTo(w io.Writer) (n int64, err error)

type Num

func (n Num) Equal(v Num) bool

func (n Num) Float64() (float64, error)

func (n Num) Format(f fmt.State, verb rune)

func (n Num) Int64() (int64, error)

func (n Num) IsInt() bool

func (n Num) Negative() bool

func (n Num) Positive() bool

func (n Num) Sign() int

func (n Num) Str() bool

func (n Num) String() string

func (n Num) Uint64() (uint64, error)

func (n Num) Zero() bool

type ObjIter

func (i *ObjIter) Err() error

func (i *ObjIter) Key() []byte

func (i *ObjIter) Next() bool

type Raw

func (r Raw) String() string

func (r Raw) Type() Type

type Type

func (t Type) String() string

type Writer

func GetWriter() *Writer

func (w *Writer) ArrEnd() bool

func (w *Writer) ArrStart() bool

func (w *Writer) Base64(data []byte) bool

func (w *Writer) Bool(v bool) bool

func (w *Writer) ByteStr(v []byte) bool

func (w *Writer) ByteStrEscape(v []byte) bool

func (w *Writer) Close() error

func (w *Writer) Comma() bool

func (w *Writer) False() bool

func (w *Writer) FieldStart(field string) bool

func (w *Writer) Float(v float64, bits int) bool

func (w *Writer) Float32(v float32) bool

func (w *Writer) Float64(v float64) bool

func (w *Writer) Flush() (fail bool)

func (w *Writer) Grow(n int)

func (w *Writer) Int(v int) bool

func (w *Writer) Int16(v int16) (fail bool)

func (w *Writer) Int32(v int32) (fail bool)

func (w *Writer) Int64(v int64) (fail bool)

func (w *Writer) Int8(v int8) (fail bool)

func (w *Writer) Null() bool

func (w *Writer) Num(v Num) bool

func (w *Writer) ObjEnd() bool

func (w *Writer) ObjStart() bool

func (w *Writer) Raw(b []byte) bool

func (w *Writer) RawStr(v string) bool

func (w *Writer) Reset()

func (w *Writer) ResetWriter(out io.Writer)

func (w *Writer) Str(v string) bool

func (w *Writer) StrEscape(v string) bool

func (w Writer) String() string

func (w *Writer) True() bool

func (w *Writer) UInt(v uint) bool

func (w *Writer) UInt16(v uint16) (fail bool)

func (w *Writer) UInt32(v uint32) (fail bool)

func (w *Writer) UInt64(v uint64) (fail bool)

func (w *Writer) UInt8(v uint8) bool

func (w *Writer) Write(p []byte) (n int, err error)

func (w *Writer) WriteTo(t io.Writer) (n int64, err error)

Examples

¶

Package

DecodeStr

Decoder.Base64

Decoder.Capture

Decoder.Num

Decoder.Raw

Encoder.Base64

Encoder.SetIdent

Encoder.String

Valid

Constants

¶

This section is empty.

Variables

¶

This section is empty.

Functions

¶

func

PutDecoder

¶

func PutDecoder(d *

Decoder

)

PutDecoder puts *Decoder into pool.

func

PutEncoder

¶

func PutEncoder(e *

Encoder

)

PutEncoder puts *Encoder to pool

func

PutWriter

¶

added in

v0.32.0

func PutWriter(e *

Writer

)

PutWriter puts *Writer to pool

func

Valid

¶

func Valid(data []

byte

)

bool

Valid reports whether data is valid json.

Example

¶

package main

import (
	"fmt"

	"github.com/go-faster/jx"
)

func main() {
	fmt.Println(jx.Valid([]byte(`{"field": "value"}`)))
	fmt.Println(jx.Valid([]byte(`"Hello, world!"`)))
	fmt.Println(jx.Valid([]byte(`["foo"}`)))
}

Output:

true
true
false

Share

Format

Run

Types

¶

type

ArrIter

¶

added in

v0.33.0

type ArrIter struct {

// contains filtered or unexported fields

}

ArrIter is decoding array iterator.

func (*ArrIter)

Err

¶

added in

v0.33.0

func (i *

ArrIter

) Err()

error

Err returns the error, if any, that was encountered during iteration.

func (*ArrIter)

Next

¶

added in

v0.33.0

func (i *

ArrIter

) Next()

bool

Next consumes element and returns false, if there is no elements anymore.

type

Decoder

¶

type Decoder struct {

// contains filtered or unexported fields

}

Decoder decodes json.

Can decode from io.Reader or byte slice directly.

func

Decode

¶

func Decode(reader

io

.

Reader

, bufSize

int

) *

Decoder

Decode creates a Decoder that reads json from io.Reader.

func

DecodeBytes

¶

func DecodeBytes(input []

byte

) *

Decoder

DecodeBytes creates a Decoder that reads json from byte slice.

func

DecodeStr

¶

func DecodeStr(input

string

) *

Decoder

DecodeStr creates a Decoder that reads string as json.

Example

¶

package main

import (
	"fmt"

	"github.com/go-faster/jx"
)

func main() {
	d := jx.DecodeStr(`{"values":[4,8,15,16,23,42]}`)

	// Save all integers from "values" array to slice.
	var values []int

	// Iterate over each object field.
	if err := d.Obj(func(d *jx.Decoder, key string) error {
		switch key {
		case "values":
			// Iterate over each array element.
			return d.Arr(func(d *jx.Decoder) error {
				v, err := d.Int()
				if err != nil {
					return err
				}
				values = append(values, v)
				return nil
			})
		default:
			// Skip unknown fields if any.
			return d.Skip()
		}
	}); err != nil {
		panic(err)
	}

	fmt.Println(values)
}

Output:

[4 8 15 16 23 42]

Share

Format

Run

func

GetDecoder

¶

func GetDecoder() *

Decoder

GetDecoder gets *Decoder from pool.

func (*Decoder)

Arr

¶

func (d *

Decoder

) Arr(f func(d *

Decoder

)

error

)

error

Arr decodes array and invokes callback on each array element.

func (*Decoder)

ArrIter

¶

added in

v0.33.0

func (d *

Decoder

) ArrIter() (

ArrIter

,

error

)

ArrIter creates new array iterator.

func (*Decoder)

Base64

¶

func (d *

Decoder

) Base64() ([]

byte

,

error

)

Base64 decodes base64 encoded data from string.

Same as encoding/json, base64.StdEncoding or

RFC 4648

.

Example

¶

package main

import (
	"fmt"

	"github.com/go-faster/jx"
)

func main() {
	data, _ := jx.DecodeStr(`"SGVsbG8="`).Base64()
	fmt.Printf("%s", data)
}

Output:

Hello

Share

Format

Run

func (*Decoder)

Base64Append

¶

func (d *

Decoder

) Base64Append(b []

byte

) ([]

byte

,

error

)

Base64Append appends base64 encoded data from string.

Same as encoding/json, base64.StdEncoding or

RFC 4648

.

func (*Decoder)

BigFloat

¶

func (d *

Decoder

) BigFloat() (*

big

.

Float

,

error

)

BigFloat read big.Float

func (*Decoder)

BigInt

¶

func (d *

Decoder

) BigInt() (*

big

.

Int

,

error

)

BigInt read big.Int

func (*Decoder)

Bool

¶

func (d *

Decoder

) Bool() (

bool

,

error

)

Bool reads a json object as Bool

func (*Decoder)

Capture

¶

func (d *

Decoder

) Capture(f func(d *

Decoder

)

error

)

error

Capture calls f and then rolls back to state before call.

Example

¶

package main

import (
	"fmt"

	"github.com/go-faster/jx"
)

func main() {
	d := jx.DecodeStr(`["foo", "bar", "baz"]`)
	var elems int
	// NB: Currently Capture does not support io.Reader, only buffers.
	if err := d.Capture(func(d *jx.Decoder) error {
		// Everything decoded in this callback will be rolled back.
		return d.Arr(func(d *jx.Decoder) error {
			elems++
			return d.Skip()
		})
	}); err != nil {
		panic(err)
	}
	// Decoder is rolled back to state before "Capture" call.
	fmt.Println("Read", elems, "elements on first pass")
	fmt.Println("Next element is", d.Next(), "again")

}

Output:

Read 3 elements on first pass
Next element is array again

Share

Format

Run

func (*Decoder)

Elem

¶

func (d *

Decoder

) Elem() (ok

bool

, err

error

)

Elem skips to the start of next array element, returning true boolean
if element exists.

Can be called before or in Array.

func (*Decoder)

Float32

¶

func (d *

Decoder

) Float32() (

float32

,

error

)

Float32 reads float32 value.

func (*Decoder)

Float64

¶

func (d *

Decoder

) Float64() (

float64

,

error

)

Float64 read float64

func (*Decoder)

Int

¶

func (d *

Decoder

) Int() (

int

,

error

)

Int reads int.

func (*Decoder)

Int16

¶

added in

v0.40.0

func (d *

Decoder

) Int16() (

int16

,

error

)

Int16 reads int16.

func (*Decoder)

Int32

¶

func (d *

Decoder

) Int32() (

int32

,

error

)

Int32 reads int32.

func (*Decoder)

Int64

¶

func (d *

Decoder

) Int64() (

int64

,

error

)

Int64 reads int64.

func (*Decoder)

Int8

¶

added in

v0.40.0

func (d *

Decoder

) Int8() (

int8

,

error

)

Int8 reads int8.

func (*Decoder)

Next

¶

func (d *

Decoder

) Next()

Type

Next gets Type of relatively next json element

func (*Decoder)

Null

¶

func (d *

Decoder

) Null()

error

Null reads a json object as null and
returns whether it's a null or not.

func (*Decoder)

Num

¶

func (d *

Decoder

) Num() (

Num

,

error

)

Num decodes number.

Do not retain returned value, it references underlying buffer.

Example

¶

package main

import (
	"fmt"

	"github.com/go-faster/jx"
)

func main() {
	// Can decode numbers and number strings.
	d := jx.DecodeStr(`{"foo": "10531.0"}`)

	var n jx.Num
	if err := d.Obj(func(d *jx.Decoder, key string) error {
		v, err := d.Num()
		if err != nil {
			return err
		}
		n = v
		return nil
	}); err != nil {
		panic(err)
	}

	fmt.Println(n)
	fmt.Println("positive:", n.Positive())

	// Can decode floats with zero fractional part as integers:
	v, err := n.Int64()
	if err != nil {
		panic(err)
	}
	fmt.Println("int64:", v)
}

Output:

"10531.0"
positive: true
int64: 10531

Share

Format

Run

func (*Decoder)

NumAppend

¶

func (d *

Decoder

) NumAppend(v

Num

) (

Num

,

error

)

NumAppend appends number.

func (*Decoder)

Obj

¶

func (d *

Decoder

) Obj(f func(d *

Decoder

, key

string

)

error

)

error

Obj reads json object, calling f on each field.

Use ObjBytes to reduce heap allocations for keys.

func (*Decoder)

ObjBytes

¶

func (d *

Decoder

) ObjBytes(f func(d *

Decoder

, key []

byte

)

error

)

error

ObjBytes calls f for every key in object, using byte slice as key.

The key value is valid only until f is not returned.

func (*Decoder)

ObjIter

¶

added in

v0.35.0

func (d *

Decoder

) ObjIter() (

ObjIter

,

error

)

ObjIter creates new object iterator.

func (*Decoder)

Raw

¶

func (d *

Decoder

) Raw() (

Raw

,

error

)

Raw is like Skip(), but saves and returns skipped value as raw json.

Do not retain returned value, it references underlying buffer.

Example

¶

package main

import (
	"fmt"

	"github.com/go-faster/jx"
)

func main() {
	d := jx.DecodeStr(`{"foo": [1, 2, 3]}`)

	var raw jx.Raw
	if err := d.Obj(func(d *jx.Decoder, key string) error {
		v, err := d.Raw()
		if err != nil {
			return err
		}
		raw = v
		return nil
	}); err != nil {
		panic(err)
	}

	fmt.Println(raw.Type(), raw)
}

Output:

array [1, 2, 3]

Share

Format

Run

func (*Decoder)

RawAppend

¶

func (d *

Decoder

) RawAppend(buf

Raw

) (

Raw

,

error

)

RawAppend is Raw that appends saved raw json value to buf.

func (*Decoder)

Reset

¶

func (d *

Decoder

) Reset(reader

io

.

Reader

)

Reset resets reader and underlying state, next reads will use provided io.Reader.

func (*Decoder)

ResetBytes

¶

func (d *

Decoder

) ResetBytes(input []

byte

)

ResetBytes resets underlying state, next reads will use provided buffer.

func (*Decoder)

Skip

¶

func (d *

Decoder

) Skip()

error

Skip skips a json object and positions to relatively the next json object.

func (*Decoder)

Str

¶

func (d *

Decoder

) Str() (

string

,

error

)

Str reads string.

func (*Decoder)

StrAppend

¶

func (d *

Decoder

) StrAppend(b []

byte

) ([]

byte

,

error

)

StrAppend reads string and appends it to byte slice.

func (*Decoder)

StrBytes

¶

func (d *

Decoder

) StrBytes() ([]

byte

,

error

)

StrBytes returns string value as sub-slice of internal buffer.

Bytes are valid only until next call to any Decoder method.

func (*Decoder)

UInt

¶

added in

v0.26.0

func (d *

Decoder

) UInt() (

uint

,

error

)

UInt reads uint.

func (*Decoder)

UInt16

¶

added in

v0.40.0

func (d *

Decoder

) UInt16() (

uint16

,

error

)

UInt16 reads uint16.

func (*Decoder)

UInt32

¶

added in

v0.26.0

func (d *

Decoder

) UInt32() (

uint32

,

error

)

UInt32 reads uint32.

func (*Decoder)

UInt64

¶

added in

v0.26.0

func (d *

Decoder

) UInt64() (

uint64

,

error

)

UInt64 reads uint64.

func (*Decoder)

UInt8

¶

added in

v0.40.0

func (d *

Decoder

) UInt8() (

uint8

,

error

)

UInt8 reads uint8.

func (*Decoder)

Validate

¶

added in

v0.18.0

func (d *

Decoder

) Validate()

error

Validate consumes all input, validating that input is a json object
without any trialing data.

type

Encoder

¶

type Encoder struct {

// contains filtered or unexported fields

}

Encoder encodes json to underlying buffer.

Zero value is valid.

func

GetEncoder

¶

func GetEncoder() *

Encoder

GetEncoder returns *Encoder from pool.

func

NewStreamingEncoder

¶

added in

v1.0.0

func NewStreamingEncoder(w

io

.

Writer

, bufSize

int

) *

Encoder

NewStreamingEncoder creates new streaming encoder.

func (*Encoder)

Arr

¶

added in

v0.21.0

func (e *

Encoder

) Arr(f func(e *

Encoder

)) (fail

bool

)

Arr writes start of array, invokes callback and writes end of array.

If callback is nil, writes empty array.

func (*Encoder)

ArrEmpty

¶

func (e *

Encoder

) ArrEmpty()

bool

ArrEmpty writes empty array.

func (*Encoder)

ArrEnd

¶

func (e *

Encoder

) ArrEnd()

bool

ArrEnd writes end of array, performing indentation if needed.

Use Arr as convenience helper for writing arrays.

func (*Encoder)

ArrStart

¶

func (e *

Encoder

) ArrStart() (fail

bool

)

ArrStart writes start of array, performing indentation if needed.

Use Arr as convenience helper for writing arrays.

func (*Encoder)

Base64

¶

func (e *

Encoder

) Base64(data []

byte

)

bool

Base64 encodes data as standard base64 encoded string.

Same as encoding/json, base64.StdEncoding or

RFC 4648

.

Example

¶

package main

import (
	"fmt"

	"github.com/go-faster/jx"
)

func main() {
	var e jx.Encoder
	e.Base64([]byte("Hello"))
	fmt.Println(e)

	data, _ := jx.DecodeBytes(e.Bytes()).Base64()
	fmt.Printf("%s", data)
}

Output:

"SGVsbG8="
Hello

Share

Format

Run

func (*Encoder)

Bool

¶

func (e *

Encoder

) Bool(v

bool

)

bool

Bool encodes boolean.

func (*Encoder)

ByteStr

¶

added in

v0.34.0

func (e *

Encoder

) ByteStr(v []

byte

)

bool

ByteStr encodes byte slice without html escaping.

Use ByteStrEscape to escape html, this is default for encoding/json and
should be used by default for untrusted strings.

func (*Encoder)

ByteStrEscape

¶

added in

v0.34.0

func (e *

Encoder

) ByteStrEscape(v []

byte

)

bool

ByteStrEscape encodes string with html special characters escaping.

func (Encoder)

Bytes

¶

func (e

Encoder

) Bytes() []

byte

Bytes returns underlying buffer.

func (*Encoder)

Close

¶

added in

v1.0.0

func (e *

Encoder

) Close()

error

Close flushes underlying buffer to writer in streaming mode.
Otherwise, it does nothing.

func (*Encoder)

Field

¶

added in

v0.19.0

func (e *

Encoder

) Field(name

string

, f func(e *

Encoder

)) (fail

bool

)

Field encodes field start and then invokes callback.

Has ~5ns overhead over FieldStart.

func (*Encoder)

FieldStart

¶

added in

v0.22.0

func (e *

Encoder

) FieldStart(field

string

) (fail

bool

)

FieldStart encodes field name and writes colon.

For non-zero indentation also writes single space after colon.

Use Field as convenience helper for encoding fields.

func (*Encoder)

Float32

¶

func (e *

Encoder

) Float32(v

float32

)

bool

Float32 encodes float32.

NB: Infinities and NaN are represented as null.

func (*Encoder)

Float64

¶

func (e *

Encoder

) Float64(v

float64

)

bool

Float64 encodes float64.

NB: Infinities and NaN are represented as null.

func (*Encoder)

Grow

¶

added in

v1.1.0

func (e *

Encoder

) Grow(n

int

)

Grow grows the underlying buffer

func (*Encoder)

Int

¶

func (e *

Encoder

) Int(v

int

)

bool

Int encodes int.

func (*Encoder)

Int16

¶

added in

v0.25.0

func (e *

Encoder

) Int16(v

int16

)

bool

Int16 encodes int16.

func (*Encoder)

Int32

¶

func (e *

Encoder

) Int32(v

int32

)

bool

Int32 encodes int32.

func (*Encoder)

Int64

¶

func (e *

Encoder

) Int64(v

int64

)

bool

Int64 encodes int64.

func (*Encoder)

Int8

¶

added in

v0.25.0

func (e *

Encoder

) Int8(v

int8

)

bool

Int8 encodes int8.

func (*Encoder)

Null

¶

func (e *

Encoder

) Null()

bool

Null writes null.

func (*Encoder)

Num

¶

func (e *

Encoder

) Num(v

Num

)

bool

Num encodes number.

func (*Encoder)

Obj

¶

added in

v0.21.0

func (e *

Encoder

) Obj(f func(e *

Encoder

)) (fail

bool

)

Obj writes start of object, invokes callback and writes end of object.

If callback is nil, writes empty object.

func (*Encoder)

ObjEmpty

¶

func (e *

Encoder

) ObjEmpty()

bool

ObjEmpty writes empty object.

func (*Encoder)

ObjEnd

¶

func (e *

Encoder

) ObjEnd()

bool

ObjEnd writes end of object token, performing indentation if needed.

Use Obj as convenience helper for writing objects.

func (*Encoder)

ObjStart

¶

func (e *

Encoder

) ObjStart() (fail

bool

)

ObjStart writes object start, performing indentation if needed.

Use Obj as convenience helper for writing objects.

func (*Encoder)

Raw

¶

func (e *

Encoder

) Raw(b []

byte

)

bool

Raw writes byte slice as raw json.

func (*Encoder)

RawStr

¶

added in

v0.20.0

func (e *

Encoder

) RawStr(v

string

)

bool

RawStr writes string as raw json.

func (*Encoder)

Reset

¶

func (e *

Encoder

) Reset()

Reset resets underlying buffer.

If e is in streaming mode, it is reset to non-streaming mode.

func (*Encoder)

ResetWriter

¶

added in

v1.0.0

func (e *

Encoder

) ResetWriter(out

io

.

Writer

)

ResetWriter resets underlying buffer and sets output writer.

func (*Encoder)

SetBytes

¶

func (e *

Encoder

) SetBytes(buf []

byte

)

SetBytes sets underlying buffer.

func (*Encoder)

SetIdent

¶

func (e *

Encoder

) SetIdent(n

int

)

SetIdent sets length of single indentation step.

Example

¶

package main

import (
	"fmt"

	"github.com/go-faster/jx"
)

func main() {
	var e jx.Encoder
	e.SetIdent(2)
	e.ObjStart()

	e.FieldStart("data")
	e.ArrStart()
	e.Int(1)
	e.Int(2)
	e.ArrEnd()

	e.ObjEnd()
	fmt.Println(e)

}

Output:

{
  "data": [
    1,
    2
  ]
}

Share

Format

Run

func (*Encoder)

Str

¶

func (e *

Encoder

) Str(v

string

)

bool

Str encodes string without html escaping.

Use StrEscape to escape html, this is default for encoding/json and
should be used by default for untrusted strings.

func (*Encoder)

StrEscape

¶

func (e *

Encoder

) StrEscape(v

string

)

bool

StrEscape encodes string with html special characters escaping.

func (Encoder)

String

¶

func (e

Encoder

) String()

string

String returns string of underlying buffer.

Example

¶

package main

import (
	"fmt"

	"github.com/go-faster/jx"
)

func main() {
	var e jx.Encoder
	e.ObjStart()           // {
	e.FieldStart("values") // "values":
	e.ArrStart()           // [
	for _, v := range []int{4, 8, 15, 16, 23, 42} {
		e.Int(v)
	}
	e.ArrEnd() // ]
	e.ObjEnd() // }
	fmt.Println(e)
	fmt.Println("Buffer len:", len(e.Bytes()))
}

Output:

{"values":[4,8,15,16,23,42]}
Buffer len: 28

Share

Format

Run

func (*Encoder)

UInt

¶

added in

v0.26.0

func (e *

Encoder

) UInt(v

uint

)

bool

UInt encodes uint.

func (*Encoder)

UInt16

¶

added in

v0.26.0

func (e *

Encoder

) UInt16(v

uint16

)

bool

UInt16 encodes uint16.

func (*Encoder)

UInt32

¶

added in

v0.26.0

func (e *

Encoder

) UInt32(v

uint32

)

bool

UInt32 encodes uint32.

func (*Encoder)

UInt64

¶

added in

v0.26.0

func (e *

Encoder

) UInt64(v

uint64

)

bool

UInt64 encodes uint64.

func (*Encoder)

UInt8

¶

added in

v0.26.0

func (e *

Encoder

) UInt8(v

uint8

)

bool

UInt8 encodes uint8.

func (*Encoder)

Write

¶

func (e *

Encoder

) Write(p []

byte

) (n

int

, err

error

)

Write implements io.Writer.

func (*Encoder)

WriteTo

¶

func (e *

Encoder

) WriteTo(w

io

.

Writer

) (n

int64

, err

error

)

WriteTo implements io.WriterTo.

type

Num

¶

type Num []

byte

Num represents number, which can be raw json number or number string.

Same as Raw, but with number invariants.

Examples:

123.45   // Str: false, IsInt: false
"123.45" // Str: true,  IsInt: false
"12345"  // Str: true,  IsInt: true
12345    // Str: false, IsInt: true

func (Num)

Equal

¶

func (n

Num

) Equal(v

Num

)

bool

Equal reports whether numbers are strictly equal, including their formats.

func (Num)

Float64

¶

func (n

Num

) Float64() (

float64

,

error

)

Float64 decodes number as 64-bit floating point.

func (Num)

Format

¶

added in

v0.23.2

func (n

Num

) Format(f

fmt

.

State

, verb

rune

)

Format implements fmt.Formatter.

func (Num)

Int64

¶

func (n

Num

) Int64() (

int64

,

error

)

Int64 decodes number as a signed 64-bit integer.
Works on floats with zero fractional part.

func (Num)

IsInt

¶

func (n

Num

) IsInt()

bool

IsInt reports whether number is integer.

func (Num)

Negative

¶

func (n

Num

) Negative()

bool

Negative reports whether number is negative.

func (Num)

Positive

¶

func (n

Num

) Positive()

bool

Positive reports whether number is positive.

func (Num)

Sign

¶

func (n

Num

) Sign()

int

Sign reports sign of number.

0 is zero, 1 is positive, -1 is negative.

func (Num)

Str

¶

func (n

Num

) Str()

bool

Str reports whether Num is string number.

func (Num)

String

¶

func (n

Num

) String()

string

func (Num)

Uint64

¶

func (n

Num

) Uint64() (

uint64

,

error

)

Uint64 decodes number as an unsigned 64-bit integer.
Works on floats with zero fractional part.

func (Num)

Zero

¶

func (n

Num

) Zero()

bool

Zero reports whether number is zero.

type

ObjIter

¶

added in

v0.35.0

type ObjIter struct {

// contains filtered or unexported fields

}

ObjIter is decoding object iterator.

func (*ObjIter)

Err

¶

added in

v0.35.0

func (i *

ObjIter

) Err()

error

Err returns the error, if any, that was encountered during iteration.

func (*ObjIter)

Key

¶

added in

v0.35.0

func (i *

ObjIter

) Key() []

byte

Key returns current key.

Key call must be preceded by a call to Next.

func (*ObjIter)

Next

¶

added in

v0.35.0

func (i *

ObjIter

) Next()

bool

Next consumes element and returns false, if there is no elements anymore.

type

Raw

¶

type Raw []

byte

Raw json value.

func (Raw)

String

¶

func (r

Raw

) String()

string

func (Raw)

Type

¶

func (r

Raw

) Type()

Type

Type of Raw json value.

type

Type

¶

type Type

int

Type of json value.

const (

// Invalid json value.

Invalid

Type

=

iota

// String json value, like "foo".

String

// Number json value, like 100 or 1.01.

Number

// Null json value.

Null

// Bool json value, true or false.

Bool

// Array json value, like [1, 2, 3].

Array

// Object json value, like {"foo": 1}.

Object
)

func (Type)

String

¶

func (t

Type

) String()

string

type

Writer

¶

added in

v0.26.0

type Writer struct {

Buf []

byte

// underlying buffer

// contains filtered or unexported fields

}

Writer writes json tokens to underlying buffer.

Zero value is valid.

func

GetWriter

¶

added in

v0.32.0

func GetWriter() *

Writer

GetWriter returns *Writer from pool.

func (*Writer)

ArrEnd

¶

added in

v0.26.0

func (w *

Writer

) ArrEnd()

bool

ArrEnd writes end of array.

func (*Writer)

ArrStart

¶

added in

v0.26.0

func (w *

Writer

) ArrStart()

bool

ArrStart writes start of array.

func (*Writer)

Base64

¶

added in

v0.26.0

func (w *

Writer

) Base64(data []

byte

)

bool

Base64 encodes data as standard base64 encoded string.

Same as encoding/json, base64.StdEncoding or

RFC 4648

.

func (*Writer)

Bool

¶

added in

v0.26.0

func (w *

Writer

) Bool(v

bool

)

bool

Bool encodes boolean.

func (*Writer)

ByteStr

¶

added in

v0.34.0

func (w *

Writer

) ByteStr(v []

byte

)

bool

ByteStr encodes string without html escaping.

Use ByteStrEscape to escape html, this is default for encoding/json and
should be used by default for untrusted strings.

func (*Writer)

ByteStrEscape

¶

added in

v0.34.0

func (w *

Writer

) ByteStrEscape(v []

byte

)

bool

ByteStrEscape encodes string with html special characters escaping.

func (*Writer)

Close

¶

added in

v1.0.0

func (w *

Writer

) Close()

error

Close flushes underlying buffer to writer in streaming mode.
Otherwise, it does nothing.

func (*Writer)

Comma

¶

added in

v0.26.0

func (w *

Writer

) Comma()

bool

Comma writes comma.

func (*Writer)

False

¶

added in

v0.26.0

func (w *

Writer

) False()

bool

False writes false.

func (*Writer)

FieldStart

¶

added in

v0.26.0

func (w *

Writer

) FieldStart(field

string

)

bool

FieldStart encodes field name and writes colon.

func (*Writer)

Float

¶

added in

v0.26.0

func (w *

Writer

) Float(v

float64

, bits

int

)

bool

Float writes float value to buffer.

func (*Writer)

Float32

¶

added in

v0.26.0

func (w *

Writer

) Float32(v

float32

)

bool

Float32 encodes float32.

NB: Infinities and NaN are represented as null.

func (*Writer)

Float64

¶

added in

v0.26.0

func (w *

Writer

) Float64(v

float64

)

bool

Float64 encodes float64.

NB: Infinities and NaN are represented as null.

func (*Writer)

Flush

¶

added in

v1.2.0

func (w *

Writer

) Flush() (fail

bool

)

Flush flushes the stream. It does nothing if not in streaming mode

func (*Writer)

Grow

¶

added in

v1.1.0

func (w *

Writer

) Grow(n

int

)

Grow grows the underlying buffer.

Calls (*bytes.Buffer).Grow(n int) on w.Buf.

func (*Writer)

Int

¶

added in

v0.26.0

func (w *

Writer

) Int(v

int

)

bool

Int encodes int.

func (*Writer)

Int16

¶

added in

v0.26.0

func (w *

Writer

) Int16(v

int16

) (fail

bool

)

Int16 encodes int16.

func (*Writer)

Int32

¶

added in

v0.26.0

func (w *

Writer

) Int32(v

int32

) (fail

bool

)

Int32 encodes int32.

func (*Writer)

Int64

¶

added in

v0.26.0

func (w *

Writer

) Int64(v

int64

) (fail

bool

)

Int64 encodes int64.

func (*Writer)

Int8

¶

added in

v0.26.0

func (w *

Writer

) Int8(v

int8

) (fail

bool

)

Int8 encodes int8.

func (*Writer)

Null

¶

added in

v0.26.0

func (w *

Writer

) Null()

bool

Null writes null.

func (*Writer)

Num

¶

added in

v0.26.0

func (w *

Writer

) Num(v

Num

)

bool

Num encodes number.

func (*Writer)

ObjEnd

¶

added in

v0.26.0

func (w *

Writer

) ObjEnd()

bool

ObjEnd writes end of object token.

func (*Writer)

ObjStart

¶

added in

v0.26.0

func (w *

Writer

) ObjStart()

bool

ObjStart writes object start.

func (*Writer)

Raw

¶

added in

v0.26.0

func (w *

Writer

) Raw(b []

byte

)

bool

Raw writes byte slice as raw json.

func (*Writer)

RawStr

¶

added in

v0.26.0

func (w *

Writer

) RawStr(v

string

)

bool

RawStr writes string as raw json.

func (*Writer)

Reset

¶

added in

v0.26.0

func (w *

Writer

) Reset()

Reset resets underlying buffer.

If w is in streaming mode, it is reset to non-streaming mode.

func (*Writer)

ResetWriter

¶

added in

v1.0.0

func (w *

Writer

) ResetWriter(out

io

.

Writer

)

ResetWriter resets underlying buffer and sets output writer.

func (*Writer)

Str

¶

added in

v0.26.0

func (w *

Writer

) Str(v

string

)

bool

Str encodes string without html escaping.

Use StrEscape to escape html, this is default for encoding/json and
should be used by default for untrusted strings.

func (*Writer)

StrEscape

¶

added in

v0.26.0

func (w *

Writer

) StrEscape(v

string

)

bool

StrEscape encodes string with html special characters escaping.

func (Writer)

String

¶

added in

v0.26.0

func (w

Writer

) String()

string

String returns string of underlying buffer.

func (*Writer)

True

¶

added in

v0.26.0

func (w *

Writer

) True()

bool

True writes true.

func (*Writer)

UInt

¶

added in

v0.26.0

func (w *

Writer

) UInt(v

uint

)

bool

UInt encodes uint.

func (*Writer)

UInt16

¶

added in

v0.26.0

func (w *

Writer

) UInt16(v

uint16

) (fail

bool

)

UInt16 encodes uint16.

func (*Writer)

UInt32

¶

added in

v0.26.0

func (w *

Writer

) UInt32(v

uint32

) (fail

bool

)

UInt32 encodes uint32.

func (*Writer)

UInt64

¶

added in

v0.26.0

func (w *

Writer

) UInt64(v

uint64

) (fail

bool

)

UInt64 encodes uint64.

func (*Writer)

UInt8

¶

added in

v0.26.0

func (w *

Writer

) UInt8(v

uint8

)

bool

UInt8 encodes uint8.

func (*Writer)

Write

¶

added in

v0.26.0

func (w *

Writer

) Write(p []

byte

) (n

int

, err

error

)

Write implements io.Writer.

func (*Writer)

WriteTo

¶

added in

v0.26.0

func (w *

Writer

) WriteTo(t

io

.

Writer

) (n

int64

, err

error

)

WriteTo implements io.WriterTo.