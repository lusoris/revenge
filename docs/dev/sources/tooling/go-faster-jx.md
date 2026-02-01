# go-faster/jx (JSON)

> Source: https://pkg.go.dev/github.com/go-faster/jx
> Fetched: 2026-02-01T11:42:34.727883+00:00
> Content-Hash: aeec1650fbf6cad9
> Type: html

---

### Overview ¶

Package jx implements [RFC 7159](https://rfc-editor.org/rfc/rfc7159.html) json encoding and decoding.

Example ¶

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
    

Share Format Run

### Index ¶

- func PutDecoder(d *Decoder)
- func PutEncoder(e *Encoder)
- func PutWriter(e *Writer)
- func Valid(data []byte) bool
- type ArrIter
-     * func (i *ArrIter) Err() error
  - func (i *ArrIter) Next() bool
- type Decoder
-     * func Decode(reader io.Reader, bufSize int) *Decoder
  - func DecodeBytes(input []byte) *Decoder
  - func DecodeStr(input string) *Decoder
  - func GetDecoder() *Decoder
-     * func (d *Decoder) Arr(f func(d *Decoder) error) error
  - func (d *Decoder) ArrIter() (ArrIter, error)
  - func (d *Decoder) Base64() ([]byte, error)
  - func (d *Decoder) Base64Append(b []byte) ([]byte, error)
  - func (d *Decoder) BigFloat() (*big.Float, error)
  - func (d *Decoder) BigInt() (*big.Int, error)
  - func (d *Decoder) Bool() (bool, error)
  - func (d *Decoder) Capture(f func(d*Decoder) error) error
  - func (d *Decoder) Elem() (ok bool, err error)
  - func (d *Decoder) Float32() (float32, error)
  - func (d *Decoder) Float64() (float64, error)
  - func (d *Decoder) Int() (int, error)
  - func (d *Decoder) Int16() (int16, error)
  - func (d *Decoder) Int32() (int32, error)
  - func (d *Decoder) Int64() (int64, error)
  - func (d *Decoder) Int8() (int8, error)
  - func (d *Decoder) Next() Type
  - func (d *Decoder) Null() error
  - func (d *Decoder) Num() (Num, error)
  - func (d *Decoder) NumAppend(v Num) (Num, error)
  - func (d *Decoder) Obj(f func(d*Decoder, key string) error) error
  - func (d *Decoder) ObjBytes(f func(d*Decoder, key []byte) error) error
  - func (d *Decoder) ObjIter() (ObjIter, error)
  - func (d *Decoder) Raw() (Raw, error)
  - func (d *Decoder) RawAppend(buf Raw) (Raw, error)
  - func (d *Decoder) Reset(reader io.Reader)
  - func (d *Decoder) ResetBytes(input []byte)
  - func (d *Decoder) Skip() error
  - func (d *Decoder) Str() (string, error)
  - func (d *Decoder) StrAppend(b []byte) ([]byte, error)
  - func (d *Decoder) StrBytes() ([]byte, error)
  - func (d *Decoder) UInt() (uint, error)
  - func (d *Decoder) UInt16() (uint16, error)
  - func (d *Decoder) UInt32() (uint32, error)
  - func (d *Decoder) UInt64() (uint64, error)
  - func (d *Decoder) UInt8() (uint8, error)
  - func (d *Decoder) Validate() error
- type Encoder
-     * func GetEncoder() *Encoder
  - func NewStreamingEncoder(w io.Writer, bufSize int) *Encoder
-     * func (e *Encoder) Arr(f func(e *Encoder)) (fail bool)
  - func (e *Encoder) ArrEmpty() bool
  - func (e *Encoder) ArrEnd() bool
  - func (e *Encoder) ArrStart() (fail bool)
  - func (e *Encoder) Base64(data []byte) bool
  - func (e *Encoder) Bool(v bool) bool
  - func (e *Encoder) ByteStr(v []byte) bool
  - func (e *Encoder) ByteStrEscape(v []byte) bool
  - func (e Encoder) Bytes() []byte
  - func (e *Encoder) Close() error
  - func (e *Encoder) Field(name string, f func(e*Encoder)) (fail bool)
  - func (e *Encoder) FieldStart(field string) (fail bool)
  - func (e *Encoder) Float32(v float32) bool
  - func (e *Encoder) Float64(v float64) bool
  - func (e *Encoder) Grow(n int)
  - func (e *Encoder) Int(v int) bool
  - func (e *Encoder) Int16(v int16) bool
  - func (e *Encoder) Int32(v int32) bool
  - func (e *Encoder) Int64(v int64) bool
  - func (e *Encoder) Int8(v int8) bool
  - func (e *Encoder) Null() bool
  - func (e *Encoder) Num(v Num) bool
  - func (e *Encoder) Obj(f func(e*Encoder)) (fail bool)
  - func (e *Encoder) ObjEmpty() bool
  - func (e *Encoder) ObjEnd() bool
  - func (e *Encoder) ObjStart() (fail bool)
  - func (e *Encoder) Raw(b []byte) bool
  - func (e *Encoder) RawStr(v string) bool
  - func (e *Encoder) Reset()
  - func (e *Encoder) ResetWriter(out io.Writer)
  - func (e *Encoder) SetBytes(buf []byte)
  - func (e *Encoder) SetIdent(n int)
  - func (e *Encoder) Str(v string) bool
  - func (e *Encoder) StrEscape(v string) bool
  - func (e Encoder) String() string
  - func (e *Encoder) UInt(v uint) bool
  - func (e *Encoder) UInt16(v uint16) bool
  - func (e *Encoder) UInt32(v uint32) bool
  - func (e *Encoder) UInt64(v uint64) bool
  - func (e *Encoder) UInt8(v uint8) bool
  - func (e *Encoder) Write(p []byte) (n int, err error)
  - func (e *Encoder) WriteTo(w io.Writer) (n int64, err error)
- type Num
-     * func (n Num) Equal(v Num) bool
  - func (n Num) Float64() (float64, error)
  - func (n Num) Format(f fmt.State, verb rune)
  - func (n Num) Int64() (int64, error)
  - func (n Num) IsInt() bool
  - func (n Num) Negative() bool
  - func (n Num) Positive() bool
  - func (n Num) Sign() int
  - func (n Num) Str() bool
  - func (n Num) String() string
  - func (n Num) Uint64() (uint64, error)
  - func (n Num) Zero() bool
- type ObjIter
-     * func (i *ObjIter) Err() error
  - func (i *ObjIter) Key() []byte
  - func (i *ObjIter) Next() bool
- type Raw
-     * func (r Raw) String() string
  - func (r Raw) Type() Type
- type Type
-     * func (t Type) String() string
- type Writer
-     * func GetWriter() *Writer
-     * func (w *Writer) ArrEnd() bool
  - func (w *Writer) ArrStart() bool
  - func (w *Writer) Base64(data []byte) bool
  - func (w *Writer) Bool(v bool) bool
  - func (w *Writer) ByteStr(v []byte) bool
  - func (w *Writer) ByteStrEscape(v []byte) bool
  - func (w *Writer) Close() error
  - func (w *Writer) Comma() bool
  - func (w *Writer) False() bool
  - func (w *Writer) FieldStart(field string) bool
  - func (w *Writer) Float(v float64, bits int) bool
  - func (w *Writer) Float32(v float32) bool
  - func (w *Writer) Float64(v float64) bool
  - func (w *Writer) Flush() (fail bool)
  - func (w *Writer) Grow(n int)
  - func (w *Writer) Int(v int) bool
  - func (w *Writer) Int16(v int16) (fail bool)
  - func (w *Writer) Int32(v int32) (fail bool)
  - func (w *Writer) Int64(v int64) (fail bool)
  - func (w *Writer) Int8(v int8) (fail bool)
  - func (w *Writer) Null() bool
  - func (w *Writer) Num(v Num) bool
  - func (w *Writer) ObjEnd() bool
  - func (w *Writer) ObjStart() bool
  - func (w *Writer) Raw(b []byte) bool
  - func (w *Writer) RawStr(v string) bool
  - func (w *Writer) Reset()
  - func (w *Writer) ResetWriter(out io.Writer)
  - func (w *Writer) Str(v string) bool
  - func (w *Writer) StrEscape(v string) bool
  - func (w Writer) String() string
  - func (w *Writer) True() bool
  - func (w *Writer) UInt(v uint) bool
  - func (w *Writer) UInt16(v uint16) (fail bool)
  - func (w *Writer) UInt32(v uint32) (fail bool)
  - func (w *Writer) UInt64(v uint64) (fail bool)
  - func (w *Writer) UInt8(v uint8) bool
  - func (w *Writer) Write(p []byte) (n int, err error)
  - func (w *Writer) WriteTo(t io.Writer) (n int64, err error)

### Examples ¶

- Package
- DecodeStr
- Decoder.Base64
- Decoder.Capture
- Decoder.Num
- Decoder.Raw
- Encoder.Base64
- Encoder.SetIdent
- Encoder.String
- Valid

### Constants ¶

This section is empty.

### Variables ¶

This section is empty.

### Functions ¶

#### func [PutDecoder](https://github.com/go-faster/jx/blob/v1.2.0/jx.go#L40) ¶

    func PutDecoder(d *Decoder)

PutDecoder puts *Decoder into pool.

#### func [PutEncoder](https://github.com/go-faster/jx/blob/v1.2.0/jx.go#L51) ¶

    func PutEncoder(e *Encoder)

PutEncoder puts *Encoder to pool

#### func [PutWriter](https://github.com/go-faster/jx/blob/v1.2.0/jx.go#L63) ¶ added in v0.32.0

    func PutWriter(e *Writer)

PutWriter puts *Writer to pool

#### func [Valid](https://github.com/go-faster/jx/blob/v1.2.0/jx.go#L9) ¶

    func Valid(data [][byte](/builtin#byte)) [bool](/builtin#bool)

Valid reports whether data is valid json.

Example ¶

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
    

Share Format Run

### Types ¶

#### type [ArrIter](https://github.com/go-faster/jx/blob/v1.2.0/dec_arr_iter.go#L8) ¶ added in v0.33.0

    type ArrIter struct {
     // contains filtered or unexported fields
    }

ArrIter is decoding array iterator.

#### func (*ArrIter) [Err](https://github.com/go-faster/jx/blob/v1.2.0/dec_arr_iter.go#L61) ¶ added in v0.33.0

    func (i *ArrIter) Err() [error](/builtin#error)

Err returns the error, if any, that was encountered during iteration.

#### func (*ArrIter) [Next](https://github.com/go-faster/jx/blob/v1.2.0/dec_arr_iter.go#L31) ¶ added in v0.33.0

    func (i *ArrIter) Next() [bool](/builtin#bool)

Next consumes element and returns false, if there is no elements anymore.

#### type [Decoder](https://github.com/go-faster/jx/blob/v1.2.0/dec.go#L77) ¶

    type Decoder struct {
     // contains filtered or unexported fields
    }

Decoder decodes json.

Can decode from io.Reader or byte slice directly.

#### func [Decode](https://github.com/go-faster/jx/blob/v1.2.0/dec.go#L95) ¶

    func Decode(reader [io](/io).[Reader](/io#Reader), bufSize [int](/builtin#int)) *Decoder

Decode creates a Decoder that reads json from io.Reader.

#### func [DecodeBytes](https://github.com/go-faster/jx/blob/v1.2.0/dec.go#L106) ¶

    func DecodeBytes(input [][byte](/builtin#byte)) *Decoder

DecodeBytes creates a Decoder that reads json from byte slice.

#### func [DecodeStr](https://github.com/go-faster/jx/blob/v1.2.0/dec.go#L114) ¶

    func DecodeStr(input [string](/builtin#string)) *Decoder

DecodeStr creates a Decoder that reads string as json.

Example ¶

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
    

Share Format Run

#### func [GetDecoder](https://github.com/go-faster/jx/blob/v1.2.0/jx.go#L35) ¶

    func GetDecoder() *Decoder

GetDecoder gets *Decoder from pool.

#### func (*Decoder) [Arr](https://github.com/go-faster/jx/blob/v1.2.0/dec_arr.go#L37) ¶

    func (d *Decoder) Arr(f func(d *Decoder) [error](/builtin#error)) [error](/builtin#error)

Arr decodes array and invokes callback on each array element.

#### func (*Decoder) [ArrIter](https://github.com/go-faster/jx/blob/v1.2.0/dec_arr_iter.go#L16) ¶ added in v0.33.0

    func (d *Decoder) ArrIter() (ArrIter, [error](/builtin#error))

ArrIter creates new array iterator.

#### func (*Decoder) [Base64](https://github.com/go-faster/jx/blob/v1.2.0/dec_b64.go#L12) ¶

    func (d *Decoder) Base64() ([][byte](/builtin#byte), [error](/builtin#error))

Base64 decodes base64 encoded data from string.

Same as encoding/json, base64.StdEncoding or [RFC 4648](https://rfc-editor.org/rfc/rfc4648.html).

Example ¶

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
    

Share Format Run

#### func (*Decoder) [Base64Append](https://github.com/go-faster/jx/blob/v1.2.0/dec_b64.go#L25) ¶

    func (d *Decoder) Base64Append(b [][byte](/builtin#byte)) ([][byte](/builtin#byte), [error](/builtin#error))

Base64Append appends base64 encoded data from string.

Same as encoding/json, base64.StdEncoding or [RFC 4648](https://rfc-editor.org/rfc/rfc4648.html).

#### func (*Decoder) [BigFloat](https://github.com/go-faster/jx/blob/v1.2.0/dec_float_big.go#L11) ¶

    func (d *Decoder) BigFloat() (*[big](/math/big).[Float](/math/big#Float), [error](/builtin#error))

BigFloat read big.Float

#### func (*Decoder) [BigInt](https://github.com/go-faster/jx/blob/v1.2.0/dec_float_big.go#L28) ¶

    func (d *Decoder) BigInt() (*[big](/math/big).[Int](/math/big#Int), [error](/builtin#error))

BigInt read big.Int

#### func (*Decoder) [Bool](https://github.com/go-faster/jx/blob/v1.2.0/dec_bool.go#L4) ¶

    func (d *Decoder) Bool() ([bool](/builtin#bool), [error](/builtin#error))

Bool reads a json object as Bool

#### func (*Decoder) [Capture](https://github.com/go-faster/jx/blob/v1.2.0/dec_capture.go#L9) ¶

    func (d *Decoder) Capture(f func(d *Decoder) [error](/builtin#error)) [error](/builtin#error)

Capture calls f and then rolls back to state before call.

Example ¶

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
    

Share Format Run

#### func (*Decoder) [Elem](https://github.com/go-faster/jx/blob/v1.2.0/dec_arr.go#L11) ¶

    func (d *Decoder) Elem() (ok [bool](/builtin#bool), err [error](/builtin#error))

Elem skips to the start of next array element, returning true boolean if element exists.

Can be called before or in Array.

#### func (*Decoder) [Float32](https://github.com/go-faster/jx/blob/v1.2.0/dec_float.go#L49) ¶

    func (d *Decoder) Float32() ([float32](/builtin#float32), [error](/builtin#error))

Float32 reads float32 value.

#### func (*Decoder) [Float64](https://github.com/go-faster/jx/blob/v1.2.0/dec_float.go#L149) ¶

    func (d *Decoder) Float64() ([float64](/builtin#float64), [error](/builtin#error))

Float64 read float64

#### func (*Decoder) [Int](https://github.com/go-faster/jx/blob/v1.2.0/dec_int.go#L25) ¶

    func (d *Decoder) Int() ([int](/builtin#int), [error](/builtin#error))

Int reads int.

#### func (*Decoder) [Int16](https://github.com/go-faster/jx/blob/v1.2.0/dec_int.gen.go#L360) ¶ added in v0.40.0

    func (d *Decoder) Int16() ([int16](/builtin#int16), [error](/builtin#error))

Int16 reads int16.

#### func (*Decoder) [Int32](https://github.com/go-faster/jx/blob/v1.2.0/dec_int.gen.go#L663) ¶

    func (d *Decoder) Int32() ([int32](/builtin#int32), [error](/builtin#error))

Int32 reads int32.

#### func (*Decoder) [Int64](https://github.com/go-faster/jx/blob/v1.2.0/dec_int.gen.go#L966) ¶

    func (d *Decoder) Int64() ([int64](/builtin#int64), [error](/builtin#error))

Int64 reads int64.

#### func (*Decoder) [Int8](https://github.com/go-faster/jx/blob/v1.2.0/dec_int.gen.go#L155) ¶ added in v0.40.0

    func (d *Decoder) Int8() ([int8](/builtin#int8), [error](/builtin#error))

Int8 reads int8.

#### func (*Decoder) [Next](https://github.com/go-faster/jx/blob/v1.2.0/dec_read.go#L9) ¶

    func (d *Decoder) Next() Type

Next gets Type of relatively next json element

#### func (*Decoder) [Null](https://github.com/go-faster/jx/blob/v1.2.0/dec_null.go#L5) ¶

    func (d *Decoder) Null() [error](/builtin#error)

Null reads a json object as null and returns whether it's a null or not.

#### func (*Decoder) [Num](https://github.com/go-faster/jx/blob/v1.2.0/dec_num.go#L10) ¶

    func (d *Decoder) Num() (Num, [error](/builtin#error))

Num decodes number.

Do not retain returned value, it references underlying buffer.

Example ¶

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
    

Share Format Run

#### func (*Decoder) [NumAppend](https://github.com/go-faster/jx/blob/v1.2.0/dec_num.go#L15) ¶

    func (d *Decoder) NumAppend(v Num) (Num, [error](/builtin#error))

NumAppend appends number.

#### func (*Decoder) [Obj](https://github.com/go-faster/jx/blob/v1.2.0/dec_obj.go#L85) ¶

    func (d *Decoder) Obj(f func(d *Decoder, key [string](/builtin#string)) [error](/builtin#error)) [error](/builtin#error)

Obj reads json object, calling f on each field.

Use ObjBytes to reduce heap allocations for keys.

#### func (*Decoder) [ObjBytes](https://github.com/go-faster/jx/blob/v1.2.0/dec_obj.go#L10) ¶

    func (d *Decoder) ObjBytes(f func(d *Decoder, key [][byte](/builtin#byte)) [error](/builtin#error)) [error](/builtin#error)

ObjBytes calls f for every key in object, using byte slice as key.

The key value is valid only until f is not returned.

#### func (*Decoder) [ObjIter](https://github.com/go-faster/jx/blob/v1.2.0/dec_obj_iter.go#L16) ¶ added in v0.35.0

    func (d *Decoder) ObjIter() (ObjIter, [error](/builtin#error))

ObjIter creates new object iterator.

#### func (*Decoder) [Raw](https://github.com/go-faster/jx/blob/v1.2.0/dec_raw.go#L33) ¶

    func (d *Decoder) Raw() (Raw, [error](/builtin#error))

Raw is like Skip(), but saves and returns skipped value as raw json.

Do not retain returned value, it references underlying buffer.

Example ¶

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
    

Share Format Run

#### func (*Decoder) [RawAppend](https://github.com/go-faster/jx/blob/v1.2.0/dec_raw.go#L64) ¶

    func (d *Decoder) RawAppend(buf Raw) (Raw, [error](/builtin#error))

RawAppend is Raw that appends saved raw json value to buf.

#### func (*Decoder) [Reset](https://github.com/go-faster/jx/blob/v1.2.0/dec.go#L123) ¶

    func (d *Decoder) Reset(reader [io](/io).[Reader](/io#Reader))

Reset resets reader and underlying state, next reads will use provided io.Reader.

#### func (*Decoder) [ResetBytes](https://github.com/go-faster/jx/blob/v1.2.0/dec.go#L141) ¶

    func (d *Decoder) ResetBytes(input [][byte](/builtin#byte))

ResetBytes resets underlying state, next reads will use provided buffer.

#### func (*Decoder) [Skip](https://github.com/go-faster/jx/blob/v1.2.0/dec_skip.go#L10) ¶

    func (d *Decoder) Skip() [error](/builtin#error)

Skip skips a json object and positions to relatively the next json object.

#### func (*Decoder) [Str](https://github.com/go-faster/jx/blob/v1.2.0/dec_str.go#L239) ¶

    func (d *Decoder) Str() ([string](/builtin#string), [error](/builtin#error))

Str reads string.

#### func (*Decoder) [StrAppend](https://github.com/go-faster/jx/blob/v1.2.0/dec_str.go#L12) ¶

    func (d *Decoder) StrAppend(b [][byte](/builtin#byte)) ([][byte](/builtin#byte), [error](/builtin#error))

StrAppend reads string and appends it to byte slice.

#### func (*Decoder) [StrBytes](https://github.com/go-faster/jx/blob/v1.2.0/dec_str.go#L230) ¶

    func (d *Decoder) StrBytes() ([][byte](/builtin#byte), [error](/builtin#error))

StrBytes returns string value as sub-slice of internal buffer.

Bytes are valid only until next call to any Decoder method.

#### func (*Decoder) [UInt](https://github.com/go-faster/jx/blob/v1.2.0/dec_int.go#L47) ¶ added in v0.26.0

    func (d *Decoder) UInt() ([uint](/builtin#uint), [error](/builtin#error))

UInt reads uint.

#### func (*Decoder) [UInt16](https://github.com/go-faster/jx/blob/v1.2.0/dec_int.gen.go#L185) ¶ added in v0.40.0

    func (d *Decoder) UInt16() ([uint16](/builtin#uint16), [error](/builtin#error))

UInt16 reads uint16.

#### func (*Decoder) [UInt32](https://github.com/go-faster/jx/blob/v1.2.0/dec_int.gen.go#L390) ¶ added in v0.26.0

    func (d *Decoder) UInt32() ([uint32](/builtin#uint32), [error](/builtin#error))

UInt32 reads uint32.

#### func (*Decoder) [UInt64](https://github.com/go-faster/jx/blob/v1.2.0/dec_int.gen.go#L693) ¶ added in v0.26.0

    func (d *Decoder) UInt64() ([uint64](/builtin#uint64), [error](/builtin#error))

UInt64 reads uint64.

#### func (*Decoder) [UInt8](https://github.com/go-faster/jx/blob/v1.2.0/dec_int.gen.go#L23) ¶ added in v0.40.0

    func (d *Decoder) UInt8() ([uint8](/builtin#uint8), [error](/builtin#error))

UInt8 reads uint8.

#### func (*Decoder) [Validate](https://github.com/go-faster/jx/blob/v1.2.0/dec_validate.go#L11) ¶ added in v0.18.0

    func (d *Decoder) Validate() [error](/builtin#error)

Validate consumes all input, validating that input is a json object without any trialing data.

#### type [Encoder](https://github.com/go-faster/jx/blob/v1.2.0/enc.go#L8) ¶

    type Encoder struct {
     // contains filtered or unexported fields
    }

Encoder encodes json to underlying buffer.

Zero value is valid.

#### func [GetEncoder](https://github.com/go-faster/jx/blob/v1.2.0/jx.go#L46) ¶

    func GetEncoder() *Encoder

GetEncoder returns *Encoder from pool.

#### func [NewStreamingEncoder](https://github.com/go-faster/jx/blob/v1.2.0/enc_stream.go#L11) ¶ added in v1.0.0

    func NewStreamingEncoder(w [io](/io).[Writer](/io#Writer), bufSize [int](/builtin#int)) *Encoder

NewStreamingEncoder creates new streaming encoder.

#### func (*Encoder) [Arr](https://github.com/go-faster/jx/blob/v1.2.0/enc.go#L193) ¶ added in v0.21.0

    func (e *Encoder) Arr(f func(e *Encoder)) (fail [bool](/builtin#bool))

Arr writes start of array, invokes callback and writes end of array.

If callback is nil, writes empty array.

#### func (*Encoder) [ArrEmpty](https://github.com/go-faster/jx/blob/v1.2.0/enc.go#L175) ¶

    func (e *Encoder) ArrEmpty() [bool](/builtin#bool)

ArrEmpty writes empty array.

#### func (*Encoder) [ArrEnd](https://github.com/go-faster/jx/blob/v1.2.0/enc.go#L184) ¶

    func (e *Encoder) ArrEnd() [bool](/builtin#bool)

ArrEnd writes end of array, performing indentation if needed.

Use Arr as convenience helper for writing arrays.

#### func (*Encoder) [ArrStart](https://github.com/go-faster/jx/blob/v1.2.0/enc.go#L168) ¶

    func (e *Encoder) ArrStart() (fail [bool](/builtin#bool))

ArrStart writes start of array, performing indentation if needed.

Use Arr as convenience helper for writing arrays.

#### func (*Encoder) [Base64](https://github.com/go-faster/jx/blob/v1.2.0/enc_b64.go#L6) ¶

    func (e *Encoder) Base64(data [][byte](/builtin#byte)) [bool](/builtin#bool)

Base64 encodes data as standard base64 encoded string.

Same as encoding/json, base64.StdEncoding or [RFC 4648](https://rfc-editor.org/rfc/rfc4648.html).

Example ¶

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
    

Share Format Run

#### func (*Encoder) [Bool](https://github.com/go-faster/jx/blob/v1.2.0/enc.go#L97) ¶

    func (e *Encoder) Bool(v [bool](/builtin#bool)) [bool](/builtin#bool)

Bool encodes boolean.

#### func (*Encoder) [ByteStr](https://github.com/go-faster/jx/blob/v1.2.0/enc_str.go#L16) ¶ added in v0.34.0

    func (e *Encoder) ByteStr(v [][byte](/builtin#byte)) [bool](/builtin#bool)

ByteStr encodes byte slice without html escaping.

Use ByteStrEscape to escape html, this is default for encoding/json and should be used by default for untrusted strings.

#### func (*Encoder) [ByteStrEscape](https://github.com/go-faster/jx/blob/v1.2.0/enc_str_escape.go#L10) ¶ added in v0.34.0

    func (e *Encoder) ByteStrEscape(v [][byte](/builtin#byte)) [bool](/builtin#bool)

ByteStrEscape encodes string with html special characters escaping.

#### func (Encoder) [Bytes](https://github.com/go-faster/jx/blob/v1.2.0/enc.go#L68) ¶

    func (e Encoder) Bytes() [][byte](/builtin#byte)

Bytes returns underlying buffer.

#### func (*Encoder) [Close](https://github.com/go-faster/jx/blob/v1.2.0/enc_stream.go#L28) ¶ added in v1.0.0

    func (e *Encoder) Close() [error](/builtin#error)

Close flushes underlying buffer to writer in streaming mode. Otherwise, it does nothing.

#### func (*Encoder) [Field](https://github.com/go-faster/jx/blob/v1.2.0/enc.go#L130) ¶ added in v0.19.0

    func (e *Encoder) Field(name [string](/builtin#string), f func(e *Encoder)) (fail [bool](/builtin#bool))

Field encodes field start and then invokes callback.

Has ~5ns overhead over FieldStart.

#### func (*Encoder) [FieldStart](https://github.com/go-faster/jx/blob/v1.2.0/enc.go#L116) ¶ added in v0.22.0

    func (e *Encoder) FieldStart(field [string](/builtin#string)) (fail [bool](/builtin#bool))

FieldStart encodes field name and writes colon.

For non-zero indentation also writes single space after colon.

Use Field as convenience helper for encoding fields.

#### func (*Encoder) [Float32](https://github.com/go-faster/jx/blob/v1.2.0/enc_float.go#L6) ¶

    func (e *Encoder) Float32(v [float32](/builtin#float32)) [bool](/builtin#bool)

Float32 encodes float32.

NB: Infinities and NaN are represented as null.

#### func (*Encoder) [Float64](https://github.com/go-faster/jx/blob/v1.2.0/enc_float.go#L14) ¶

    func (e *Encoder) Float64(v [float64](/builtin#float64)) [bool](/builtin#bool)

Float64 encodes float64.

NB: Infinities and NaN are represented as null.

#### func (*Encoder) [Grow](https://github.com/go-faster/jx/blob/v1.2.0/enc.go#L63) ¶ added in v1.1.0

    func (e *Encoder) Grow(n [int](/builtin#int))

Grow grows the underlying buffer

#### func (*Encoder) [Int](https://github.com/go-faster/jx/blob/v1.2.0/enc_int.go#L4) ¶

    func (e *Encoder) Int(v [int](/builtin#int)) [bool](/builtin#bool)

Int encodes int.

#### func (*Encoder) [Int16](https://github.com/go-faster/jx/blob/v1.2.0/w_int.gen.go#L69) ¶ added in v0.25.0

    func (e *Encoder) Int16(v [int16](/builtin#int16)) [bool](/builtin#bool)

Int16 encodes int16.

#### func (*Encoder) [Int32](https://github.com/go-faster/jx/blob/v1.2.0/w_int.gen.go#L126) ¶

    func (e *Encoder) Int32(v [int32](/builtin#int32)) [bool](/builtin#bool)

Int32 encodes int32.

#### func (*Encoder) [Int64](https://github.com/go-faster/jx/blob/v1.2.0/w_int.gen.go#L219) ¶

    func (e *Encoder) Int64(v [int64](/builtin#int64)) [bool](/builtin#bool)

Int64 encodes int64.

#### func (*Encoder) [Int8](https://github.com/go-faster/jx/blob/v1.2.0/enc_int.go#L22) ¶ added in v0.25.0

    func (e *Encoder) Int8(v [int8](/builtin#int8)) [bool](/builtin#bool)

Int8 encodes int8.

#### func (*Encoder) [Null](https://github.com/go-faster/jx/blob/v1.2.0/enc.go#L91) ¶

    func (e *Encoder) Null() [bool](/builtin#bool)

Null writes null.

#### func (*Encoder) [Num](https://github.com/go-faster/jx/blob/v1.2.0/enc_num.go#L4) ¶

    func (e *Encoder) Num(v Num) [bool](/builtin#bool)

Num encodes number.

#### func (*Encoder) [Obj](https://github.com/go-faster/jx/blob/v1.2.0/enc.go#L155) ¶ added in v0.21.0

    func (e *Encoder) Obj(f func(e *Encoder)) (fail [bool](/builtin#bool))

Obj writes start of object, invokes callback and writes end of object.

If callback is nil, writes empty object.

#### func (*Encoder) [ObjEmpty](https://github.com/go-faster/jx/blob/v1.2.0/enc.go#L146) ¶

    func (e *Encoder) ObjEmpty() [bool](/builtin#bool)

ObjEmpty writes empty object.

#### func (*Encoder) [ObjEnd](https://github.com/go-faster/jx/blob/v1.2.0/enc.go#L140) ¶

    func (e *Encoder) ObjEnd() [bool](/builtin#bool)

ObjEnd writes end of object token, performing indentation if needed.

Use Obj as convenience helper for writing objects.

#### func (*Encoder) [ObjStart](https://github.com/go-faster/jx/blob/v1.2.0/enc.go#L105) ¶

    func (e *Encoder) ObjStart() (fail [bool](/builtin#bool))

ObjStart writes object start, performing indentation if needed.

Use Obj as convenience helper for writing objects.

#### func (*Encoder) [Raw](https://github.com/go-faster/jx/blob/v1.2.0/enc.go#L85) ¶

    func (e *Encoder) Raw(b [][byte](/builtin#byte)) [bool](/builtin#bool)

Raw writes byte slice as raw json.

#### func (*Encoder) [RawStr](https://github.com/go-faster/jx/blob/v1.2.0/enc.go#L79) ¶ added in v0.20.0

    func (e *Encoder) RawStr(v [string](/builtin#string)) [bool](/builtin#bool)

RawStr writes string as raw json.

#### func (*Encoder) [Reset](https://github.com/go-faster/jx/blob/v1.2.0/enc.go#L51) ¶

    func (e *Encoder) Reset()

Reset resets underlying buffer.

If e is in streaming mode, it is reset to non-streaming mode.

#### func (*Encoder) [ResetWriter](https://github.com/go-faster/jx/blob/v1.2.0/enc.go#L57) ¶ added in v1.0.0

    func (e *Encoder) ResetWriter(out [io](/io).[Writer](/io#Writer))

ResetWriter resets underlying buffer and sets output writer.

#### func (*Encoder) [SetBytes](https://github.com/go-faster/jx/blob/v1.2.0/enc.go#L71) ¶

    func (e *Encoder) SetBytes(buf [][byte](/builtin#byte))

SetBytes sets underlying buffer.

#### func (*Encoder) [SetIdent](https://github.com/go-faster/jx/blob/v1.2.0/enc.go#L39) ¶

    func (e *Encoder) SetIdent(n [int](/builtin#int))

SetIdent sets length of single indentation step.

Example ¶

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
    

Share Format Run

#### func (*Encoder) [Str](https://github.com/go-faster/jx/blob/v1.2.0/enc_str.go#L7) ¶

    func (e *Encoder) Str(v [string](/builtin#string)) [bool](/builtin#bool)

Str encodes string without html escaping.

Use StrEscape to escape html, this is default for encoding/json and should be used by default for untrusted strings.

#### func (*Encoder) [StrEscape](https://github.com/go-faster/jx/blob/v1.2.0/enc_str_escape.go#L4) ¶

    func (e *Encoder) StrEscape(v [string](/builtin#string)) [bool](/builtin#bool)

StrEscape encodes string with html special characters escaping.

#### func (Encoder) [String](https://github.com/go-faster/jx/blob/v1.2.0/enc.go#L44) ¶

    func (e Encoder) String() [string](/builtin#string)

String returns string of underlying buffer.

Example ¶

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
    

Share Format Run

#### func (*Encoder) [UInt](https://github.com/go-faster/jx/blob/v1.2.0/enc_int.go#L10) ¶ added in v0.26.0

    func (e *Encoder) UInt(v [uint](/builtin#uint)) [bool](/builtin#bool)

UInt encodes uint.

#### func (*Encoder) [UInt16](https://github.com/go-faster/jx/blob/v1.2.0/w_int.gen.go#L52) ¶ added in v0.26.0

    func (e *Encoder) UInt16(v [uint16](/builtin#uint16)) [bool](/builtin#bool)

UInt16 encodes uint16.

#### func (*Encoder) [UInt32](https://github.com/go-faster/jx/blob/v1.2.0/w_int.gen.go#L109) ¶ added in v0.26.0

    func (e *Encoder) UInt32(v [uint32](/builtin#uint32)) [bool](/builtin#bool)

UInt32 encodes uint32.

#### func (*Encoder) [UInt64](https://github.com/go-faster/jx/blob/v1.2.0/w_int.gen.go#L202) ¶ added in v0.26.0

    func (e *Encoder) UInt64(v [uint64](/builtin#uint64)) [bool](/builtin#bool)

UInt64 encodes uint64.

#### func (*Encoder) [UInt8](https://github.com/go-faster/jx/blob/v1.2.0/enc_int.go#L16) ¶ added in v0.26.0

    func (e *Encoder) UInt8(v [uint8](/builtin#uint8)) [bool](/builtin#bool)

UInt8 encodes uint8.

#### func (*Encoder) [Write](https://github.com/go-faster/jx/blob/v1.2.0/enc.go#L29) ¶

    func (e *Encoder) Write(p [][byte](/builtin#byte)) (n [int](/builtin#int), err [error](/builtin#error))

Write implements io.Writer.

#### func (*Encoder) [WriteTo](https://github.com/go-faster/jx/blob/v1.2.0/enc.go#L34) ¶

    func (e *Encoder) WriteTo(w [io](/io).[Writer](/io#Writer)) (n [int64](/builtin#int64), err [error](/builtin#error))

WriteTo implements io.WriterTo.

#### type [Num](https://github.com/go-faster/jx/blob/v1.2.0/num.go#L21) ¶

    type Num [][byte](/builtin#byte)

Num represents number, which can be raw json number or number string.

Same as Raw, but with number invariants.

Examples:

    123.45   // Str: false, IsInt: false
    "123.45" // Str: true,  IsInt: false
    "12345"  // Str: true,  IsInt: true
    12345    // Str: false, IsInt: true
    

#### func (Num) [Equal](https://github.com/go-faster/jx/blob/v1.2.0/num.go#L119) ¶

    func (n Num) Equal(v Num) [bool](/builtin#bool)

Equal reports whether numbers are strictly equal, including their formats.

#### func (Num) [Float64](https://github.com/go-faster/jx/blob/v1.2.0/num.go#L113) ¶

    func (n Num) Float64() ([float64](/builtin#float64), [error](/builtin#error))

Float64 decodes number as 64-bit floating point.

#### func (Num) [Format](https://github.com/go-faster/jx/blob/v1.2.0/num.go#L131) ¶ added in v0.23.2

    func (n Num) Format(f [fmt](/fmt).[State](/fmt#State), verb [rune](/builtin#rune))

Format implements fmt.Formatter.

#### func (Num) [Int64](https://github.com/go-faster/jx/blob/v1.2.0/num.go#L64) ¶

    func (n Num) Int64() ([int64](/builtin#int64), [error](/builtin#error))

Int64 decodes number as a signed 64-bit integer. Works on floats with zero fractional part.

#### func (Num) [IsInt](https://github.com/go-faster/jx/blob/v1.2.0/num.go#L77) ¶

    func (n Num) IsInt() [bool](/builtin#bool)

IsInt reports whether number is integer.

#### func (Num) [Negative](https://github.com/go-faster/jx/blob/v1.2.0/num.go#L182) ¶

    func (n Num) Negative() [bool](/builtin#bool)

Negative reports whether number is negative.

#### func (Num) [Positive](https://github.com/go-faster/jx/blob/v1.2.0/num.go#L179) ¶

    func (n Num) Positive() [bool](/builtin#bool)

Positive reports whether number is positive.

#### func (Num) [Sign](https://github.com/go-faster/jx/blob/v1.2.0/num.go#L157) ¶

    func (n Num) Sign() [int](/builtin#int)

Sign reports sign of number.

0 is zero, 1 is positive, -1 is negative.

#### func (Num) [Str](https://github.com/go-faster/jx/blob/v1.2.0/num.go#L38) ¶

    func (n Num) Str() [bool](/builtin#bool)

Str reports whether Num is string number.

#### func (Num) [String](https://github.com/go-faster/jx/blob/v1.2.0/num.go#L123) ¶

    func (n Num) String() [string](/builtin#string)

#### func (Num) [Uint64](https://github.com/go-faster/jx/blob/v1.2.0/num.go#L100) ¶

    func (n Num) Uint64() ([uint64](/builtin#uint64), [error](/builtin#error))

Uint64 decodes number as an unsigned 64-bit integer. Works on floats with zero fractional part.

#### func (Num) [Zero](https://github.com/go-faster/jx/blob/v1.2.0/num.go#L185) ¶

    func (n Num) Zero() [bool](/builtin#bool)

Zero reports whether number is zero.

#### type [ObjIter](https://github.com/go-faster/jx/blob/v1.2.0/dec_obj_iter.go#L6) ¶ added in v0.35.0

    type ObjIter struct {
     // contains filtered or unexported fields
    }

ObjIter is decoding object iterator.

#### func (*ObjIter) [Err](https://github.com/go-faster/jx/blob/v1.2.0/dec_obj_iter.go#L88) ¶ added in v0.35.0

    func (i *ObjIter) Err() [error](/builtin#error)

Err returns the error, if any, that was encountered during iteration.

#### func (*ObjIter) [Key](https://github.com/go-faster/jx/blob/v1.2.0/dec_obj_iter.go#L33) ¶ added in v0.35.0

    func (i *ObjIter) Key() [][byte](/builtin#byte)

Key returns current key.

Key call must be preceded by a call to Next.

#### func (*ObjIter) [Next](https://github.com/go-faster/jx/blob/v1.2.0/dec_obj_iter.go#L38) ¶ added in v0.35.0

    func (i *ObjIter) Next() [bool](/builtin#bool)

Next consumes element and returns false, if there is no elements anymore.

#### type [Raw](https://github.com/go-faster/jx/blob/v1.2.0/dec_raw.go#L73) ¶

    type Raw [][byte](/builtin#byte)

Raw json value.

#### func (Raw) [String](https://github.com/go-faster/jx/blob/v1.2.0/dec_raw.go#L81) ¶

    func (r Raw) String() [string](/builtin#string)

#### func (Raw) [Type](https://github.com/go-faster/jx/blob/v1.2.0/dec_raw.go#L76) ¶

    func (r Raw) Type() Type

Type of Raw json value.

#### type [Type](https://github.com/go-faster/jx/blob/v1.2.0/dec.go#L8) ¶

    type Type [int](/builtin#int)

Type of json value.

    const (
     // Invalid json value.
     Invalid Type = [iota](/builtin#iota)
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

#### func (Type) [String](https://github.com/go-faster/jx/blob/v1.2.0/dec.go#L10) ¶

    func (t Type) String() [string](/builtin#string)

#### type [Writer](https://github.com/go-faster/jx/blob/v1.2.0/w.go#L11) ¶ added in v0.26.0

    type Writer struct {
     Buf [][byte](/builtin#byte) // underlying buffer
     // contains filtered or unexported fields
    }

Writer writes json tokens to underlying buffer.

Zero value is valid.

#### func [GetWriter](https://github.com/go-faster/jx/blob/v1.2.0/jx.go#L58) ¶ added in v0.32.0

    func GetWriter() *Writer

GetWriter returns *Writer from pool.

#### func (*Writer) [ArrEnd](https://github.com/go-faster/jx/blob/v1.2.0/w.go#L158) ¶ added in v0.26.0

    func (w *Writer) ArrEnd() [bool](/builtin#bool)

ArrEnd writes end of array.

#### func (*Writer) [ArrStart](https://github.com/go-faster/jx/blob/v1.2.0/w.go#L153) ¶ added in v0.26.0

    func (w *Writer) ArrStart() [bool](/builtin#bool)

ArrStart writes start of array.

#### func (*Writer) [Base64](https://github.com/go-faster/jx/blob/v1.2.0/w_b64.go#L12) ¶ added in v0.26.0

    func (w *Writer) Base64(data [][byte](/builtin#byte)) [bool](/builtin#bool)

Base64 encodes data as standard base64 encoded string.

Same as encoding/json, base64.StdEncoding or [RFC 4648](https://rfc-editor.org/rfc/rfc4648.html).

#### func (*Writer) [Bool](https://github.com/go-faster/jx/blob/v1.2.0/w.go#L129) ¶ added in v0.26.0

    func (w *Writer) Bool(v [bool](/builtin#bool)) [bool](/builtin#bool)

Bool encodes boolean.

#### func (*Writer) [ByteStr](https://github.com/go-faster/jx/blob/v1.2.0/w_str.go#L37) ¶ added in v0.34.0

    func (w *Writer) ByteStr(v [][byte](/builtin#byte)) [bool](/builtin#bool)

ByteStr encodes string without html escaping.

Use ByteStrEscape to escape html, this is default for encoding/json and should be used by default for untrusted strings.

#### func (*Writer) [ByteStrEscape](https://github.com/go-faster/jx/blob/v1.2.0/w_str_escape.go#L121) ¶ added in v0.34.0

    func (w *Writer) ByteStrEscape(v [][byte](/builtin#byte)) [bool](/builtin#bool)

ByteStrEscape encodes string with html special characters escaping.

#### func (*Writer) [Close](https://github.com/go-faster/jx/blob/v1.2.0/w_stream.go#L12) ¶ added in v1.0.0

    func (w *Writer) Close() [error](/builtin#error)

Close flushes underlying buffer to writer in streaming mode. Otherwise, it does nothing.

#### func (*Writer) [Comma](https://github.com/go-faster/jx/blob/v1.2.0/w.go#L163) ¶ added in v0.26.0

    func (w *Writer) Comma() [bool](/builtin#bool)

Comma writes comma.

#### func (*Writer) [False](https://github.com/go-faster/jx/blob/v1.2.0/w.go#L124) ¶ added in v0.26.0

    func (w *Writer) False() [bool](/builtin#bool)

False writes false.

#### func (*Writer) [FieldStart](https://github.com/go-faster/jx/blob/v1.2.0/w.go#L142) ¶ added in v0.26.0

    func (w *Writer) FieldStart(field [string](/builtin#string)) [bool](/builtin#bool)

FieldStart encodes field name and writes colon.

#### func (*Writer) [Float](https://github.com/go-faster/jx/blob/v1.2.0/w_float_bits.go#L13) ¶ added in v0.26.0

    func (w *Writer) Float(v [float64](/builtin#float64), bits [int](/builtin#int)) [bool](/builtin#bool)

Float writes float value to buffer.

#### func (*Writer) [Float32](https://github.com/go-faster/jx/blob/v1.2.0/w_float.go#L6) ¶ added in v0.26.0

    func (w *Writer) Float32(v [float32](/builtin#float32)) [bool](/builtin#bool)

Float32 encodes float32.

NB: Infinities and NaN are represented as null.

#### func (*Writer) [Float64](https://github.com/go-faster/jx/blob/v1.2.0/w_float.go#L11) ¶ added in v0.26.0

    func (w *Writer) Float64(v [float64](/builtin#float64)) [bool](/builtin#bool)

Float64 encodes float64.

NB: Infinities and NaN are represented as null.

#### func (*Writer) [Flush](https://github.com/go-faster/jx/blob/v1.2.0/w.go#L72) ¶ added in v1.2.0

    func (w *Writer) Flush() (fail [bool](/builtin#bool))

Flush flushes the stream. It does nothing if not in streaming mode

#### func (*Writer) [Grow](https://github.com/go-faster/jx/blob/v1.2.0/w.go#L65) ¶ added in v1.1.0

    func (w *Writer) Grow(n [int](/builtin#int))

Grow grows the underlying buffer.

Calls (*bytes.Buffer).Grow(n int) on w.Buf.

#### func (*Writer) [Int](https://github.com/go-faster/jx/blob/v1.2.0/w_int.go#L4) ¶ added in v0.26.0

    func (w *Writer) Int(v [int](/builtin#int)) [bool](/builtin#bool)

Int encodes int.

#### func (*Writer) [Int16](https://github.com/go-faster/jx/blob/v1.2.0/w_int.gen.go#L57) ¶ added in v0.26.0

    func (w *Writer) Int16(v [int16](/builtin#int16)) (fail [bool](/builtin#bool))

Int16 encodes int16.

#### func (*Writer) [Int32](https://github.com/go-faster/jx/blob/v1.2.0/w_int.gen.go#L114) ¶ added in v0.26.0

    func (w *Writer) Int32(v [int32](/builtin#int32)) (fail [bool](/builtin#bool))

Int32 encodes int32.

#### func (*Writer) [Int64](https://github.com/go-faster/jx/blob/v1.2.0/w_int.gen.go#L207) ¶ added in v0.26.0

    func (w *Writer) Int64(v [int64](/builtin#int64)) (fail [bool](/builtin#bool))

Int64 encodes int64.

#### func (*Writer) [Int8](https://github.com/go-faster/jx/blob/v1.2.0/w_int.go#L20) ¶ added in v0.26.0

    func (w *Writer) Int8(v [int8](/builtin#int8)) (fail [bool](/builtin#bool))

Int8 encodes int8.

#### func (*Writer) [Null](https://github.com/go-faster/jx/blob/v1.2.0/w.go#L114) ¶ added in v0.26.0

    func (w *Writer) Null() [bool](/builtin#bool)

Null writes null.

#### func (*Writer) [Num](https://github.com/go-faster/jx/blob/v1.2.0/w_num.go#L4) ¶ added in v0.26.0

    func (w *Writer) Num(v Num) [bool](/builtin#bool)

Num encodes number.

#### func (*Writer) [ObjEnd](https://github.com/go-faster/jx/blob/v1.2.0/w.go#L148) ¶ added in v0.26.0

    func (w *Writer) ObjEnd() [bool](/builtin#bool)

ObjEnd writes end of object token.

#### func (*Writer) [ObjStart](https://github.com/go-faster/jx/blob/v1.2.0/w.go#L137) ¶ added in v0.26.0

    func (w *Writer) ObjStart() [bool](/builtin#bool)

ObjStart writes object start.

#### func (*Writer) [Raw](https://github.com/go-faster/jx/blob/v1.2.0/w.go#L109) ¶ added in v0.26.0

    func (w *Writer) Raw(b [][byte](/builtin#byte)) [bool](/builtin#bool)

Raw writes byte slice as raw json.

#### func (*Writer) [RawStr](https://github.com/go-faster/jx/blob/v1.2.0/w.go#L100) ¶ added in v0.26.0

    func (w *Writer) RawStr(v [string](/builtin#string)) [bool](/builtin#bool)

RawStr writes string as raw json.

#### func (*Writer) [Reset](https://github.com/go-faster/jx/blob/v1.2.0/w.go#L48) ¶ added in v0.26.0

    func (w *Writer) Reset()

Reset resets underlying buffer.

If w is in streaming mode, it is reset to non-streaming mode.

#### func (*Writer) [ResetWriter](https://github.com/go-faster/jx/blob/v1.2.0/w.go#L54) ¶ added in v1.0.0

    func (w *Writer) ResetWriter(out [io](/io).[Writer](/io#Writer))

ResetWriter resets underlying buffer and sets output writer.

#### func (*Writer) [Str](https://github.com/go-faster/jx/blob/v1.2.0/w_str.go#L29) ¶ added in v0.26.0

    func (w *Writer) Str(v [string](/builtin#string)) [bool](/builtin#bool)

Str encodes string without html escaping.

Use StrEscape to escape html, this is default for encoding/json and should be used by default for untrusted strings.

#### func (*Writer) [StrEscape](https://github.com/go-faster/jx/blob/v1.2.0/w_str_escape.go#L116) ¶ added in v0.26.0

    func (w *Writer) StrEscape(v [string](/builtin#string)) [bool](/builtin#bool)

StrEscape encodes string with html special characters escaping.

#### func (Writer) [String](https://github.com/go-faster/jx/blob/v1.2.0/w.go#L40) ¶ added in v0.26.0

    func (w Writer) String() [string](/builtin#string)

String returns string of underlying buffer.

#### func (*Writer) [True](https://github.com/go-faster/jx/blob/v1.2.0/w.go#L119) ¶ added in v0.26.0

    func (w *Writer) True() [bool](/builtin#bool)

True writes true.

#### func (*Writer) [UInt](https://github.com/go-faster/jx/blob/v1.2.0/w_int.go#L9) ¶ added in v0.26.0

    func (w *Writer) UInt(v [uint](/builtin#uint)) [bool](/builtin#bool)

UInt encodes uint.

#### func (*Writer) [UInt16](https://github.com/go-faster/jx/blob/v1.2.0/w_int.gen.go#L36) ¶ added in v0.26.0

    func (w *Writer) UInt16(v [uint16](/builtin#uint16)) (fail [bool](/builtin#bool))

UInt16 encodes uint16.

#### func (*Writer) [UInt32](https://github.com/go-faster/jx/blob/v1.2.0/w_int.gen.go#L74) ¶ added in v0.26.0

    func (w *Writer) UInt32(v [uint32](/builtin#uint32)) (fail [bool](/builtin#bool))

UInt32 encodes uint32.

#### func (*Writer) [UInt64](https://github.com/go-faster/jx/blob/v1.2.0/w_int.gen.go#L131) ¶ added in v0.26.0

    func (w *Writer) UInt64(v [uint64](/builtin#uint64)) (fail [bool](/builtin#bool))

UInt64 encodes uint64.

#### func (*Writer) [UInt8](https://github.com/go-faster/jx/blob/v1.2.0/w_int.go#L14) ¶ added in v0.26.0

    func (w *Writer) UInt8(v [uint8](/builtin#uint8)) [bool](/builtin#bool)

UInt8 encodes uint8.

#### func (*Writer) [Write](https://github.com/go-faster/jx/blob/v1.2.0/w.go#L17) ¶ added in v0.26.0

    func (w *Writer) Write(p [][byte](/builtin#byte)) (n [int](/builtin#int), err [error](/builtin#error))

Write implements io.Writer.

#### func (*Writer) [WriteTo](https://github.com/go-faster/jx/blob/v1.2.0/w.go#L31) ¶ added in v0.26.0

    func (w *Writer) WriteTo(t [io](/io).[Writer](/io#Writer)) (n [int64](/builtin#int64), err [error](/builtin#error))

WriteTo implements io.WriterTo.
