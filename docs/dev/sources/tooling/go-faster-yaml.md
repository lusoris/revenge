# go-faster/yaml

> Source: https://pkg.go.dev/github.com/go-faster/yaml
> Fetched: 2026-01-31T10:56:34.813885+00:00
> Content-Hash: 9f556e8437c8a825
> Type: html

---

### Overview ¶

Package yaml implements YAML support for the Go language. 

Source code and other details for the project are available at GitHub: 
    
    
    https://github.com/go-faster/yaml
    

### Index ¶

  * func Marshal(in any) (out []byte, err error)
  * func Unmarshal(in []byte, out any) (err error)
  * type Decoder
  *     * func NewDecoder(r io.Reader) *Decoder
  *     * func (dec *Decoder) Decode(v any) (err error)
    * func (dec *Decoder) KnownFields(enable bool)
  * type DuplicateKeyError
  *     * func (d *DuplicateKeyError) Error() string
  * type Encoder
  *     * func NewEncoder(w io.Writer) *Encoder
  *     * func (e *Encoder) Close() (err error)
    * func (e *Encoder) Encode(v any) (err error)
    * func (e *Encoder) SetIndent(spaces int)
  * type IsZeroer
  * type Kind
  *     * func (s Kind) String() string
  * type MarshalError
  *     * func (s *MarshalError) Error() string
  * type Marshaler
  * type Node
  *     * func (n *Node) Decode(v any) (err error)
    * func (n *Node) Encode(v any) (err error)
    * func (n *Node) EncodeJSON(e *jx.Encoder) error
    * func (n *Node) IsZero() bool
    * func (n *Node) LongTag() string
    * func (n *Node) SetString(s string)
    * func (n *Node) ShortTag() string
  * type Style
  *     * func (s Style) String() string
  * type SyntaxError
  *     * func (s *SyntaxError) Error() string
  * type TypeError
  *     * func (e *TypeError) Error() string
    * func (e *TypeError) Unwrap() error
  * type UnknownFieldError
  *     * func (d *UnknownFieldError) Error() string
  * type UnmarshalError
  *     * func (s *UnmarshalError) Error() string
    * func (s *UnmarshalError) Unwrap() error
  * type Unmarshaler



### Examples ¶

  * Unmarshal (Embedded)



### Constants ¶

This section is empty.

### Variables ¶

This section is empty.

### Functions ¶

####  func [Marshal](https://github.com/go-faster/yaml/blob/v0.4.6/marshal.go#L60) ¶
    
    
    func Marshal(in [any](/builtin#any)) (out [][byte](/builtin#byte), err [error](/builtin#error))

Marshal serializes the value provided into a YAML document. The structure of the generated document will reflect the structure of the value itself. Maps and pointers (to struct, string, int, etc) are accepted as the in value. 

Struct fields are only marshaled if they are exported (have an upper case first letter), and are marshaled using the field name lowercased as the default key. Custom keys may be defined via the "yaml" name in the field tag: the content preceding the first comma is used as the key, and the following comma-separated options are used to tweak the marshaling process. Conflicting names result in a runtime error. 

The field tag format accepted is: 
    
    
    `(...) yaml:"[<key>][,<flag1>[,<flag2>]]" (...)`
    

The following flags are currently supported: 
    
    
    omitempty    Only include the field if it's not set to the zero
                 value for the type or to empty slices or maps.
                 Zero valued structs will be omitted if all their public
                 fields are zero, unless they implement an IsZero
                 method (see the IsZeroer interface type), in which
                 case the field will be excluded if IsZero returns true.
    
    flow         Marshal using a flow style (useful for structs,
                 sequences and maps).
    
    inline       Inline the field, which must be a struct or a map,
                 causing all of its fields or keys to be processed as if
                 they were part of the outer struct. For maps, keys must
                 not conflict with the yaml keys of other struct fields.
    

In addition, if the key is "-", the field is ignored. 

For example: 
    
    
    type T struct {
        F int `yaml:"a,omitempty"`
        B int
    }
    yaml.Marshal(&T{B: 2}) // Returns "b: 2\n"
    yaml.Marshal(&T{F: 1}} // Returns "a: 1\nb: 0\n"
    

####  func [Unmarshal](https://github.com/go-faster/yaml/blob/v0.4.6/unmarshal.go#L53) ¶
    
    
    func Unmarshal(in [][byte](/builtin#byte), out [any](/builtin#any)) (err [error](/builtin#error))

Unmarshal decodes the first document found within the in byte slice and assigns decoded values into the out value. 

Maps and pointers (to a struct, string, int, etc) are accepted as out values. If an internal pointer within a struct is not initialized, the yaml package will initialize it if necessary for unmarshaling the provided data. The out parameter must not be nil. 

The type of the decoded values should be compatible with the respective values in out. If one or more values cannot be decoded due to a type mismatches, decoding continues partially until the end of the YAML content, and a *yaml.TypeError is returned with details for all missed values. 

Struct fields are only unmarshaled if they are exported (have an upper case first letter), and are unmarshaled using the field name lowercased as the default key. Custom keys may be defined via the "yaml" name in the field tag: the content preceding the first comma is used as the key, and the following comma-separated options are used to tweak the marshaling process (see Marshal). Conflicting names result in a runtime error. 

For example: 
    
    
    type T struct {
        F int `yaml:"a,omitempty"`
        B int
    }
    var t T
    yaml.Unmarshal([]byte("a: 1\nb: 2"), &t)
    

See the documentation of Marshal for the format of tags and a list of supported tag options. 

Example (Embedded) ¶
    
    
    package main
    
    import (
    	"fmt"
    	"log"
    
    	"github.com/go-faster/yaml"
    )
    
    // An example showing how to unmarshal embedded
    // structs from YAML.
    
    type StructA struct {
    	A string `yaml:"a"`
    }
    
    type StructB struct {
    	// Embedded structs are not treated as embedded in YAML by default. To do that,
    	// add the ",inline" annotation below
    	StructA `yaml:",inline"`
    	B       string `yaml:"b"`
    }
    
    func main() {
    	var (
    		b    StructB
    		data = `
    a: a string from struct A
    b: a string from struct B
    `
    	)
    
    	err := yaml.Unmarshal([]byte(data), &b)
    	if err != nil {
    		log.Fatalf("cannot unmarshal data: %v", err)
    	}
    	fmt.Println(b.A)
    	fmt.Println(b.B)
    }
    
    
    
    Output:
    
    a string from struct A
    a string from struct B
    

Share Format Run

### Types ¶

####  type [Decoder](https://github.com/go-faster/yaml/blob/v0.4.6/unmarshal.go#L58) ¶
    
    
    type Decoder struct {
    	// contains filtered or unexported fields
    }

A Decoder reads and decodes YAML values from an input stream. 

####  func [NewDecoder](https://github.com/go-faster/yaml/blob/v0.4.6/unmarshal.go#L67) ¶
    
    
    func NewDecoder(r [io](/io).[Reader](/io#Reader)) *Decoder

NewDecoder returns a new decoder that reads from r. 

The decoder introduces its own buffering and may read data from r beyond the YAML values requested. 

####  func (*Decoder) [Decode](https://github.com/go-faster/yaml/blob/v0.4.6/unmarshal.go#L84) ¶
    
    
    func (dec *Decoder) Decode(v [any](/builtin#any)) (err [error](/builtin#error))

Decode reads the next YAML-encoded value from its input and stores it in the value pointed to by v. 

See the documentation for Unmarshal for details about the conversion of YAML into a Go value. 

####  func (*Decoder) [KnownFields](https://github.com/go-faster/yaml/blob/v0.4.6/unmarshal.go#L75) ¶
    
    
    func (dec *Decoder) KnownFields(enable [bool](/builtin#bool))

KnownFields ensures that the keys in decoded mappings to exist as fields in the struct being decoded into. 

####  type [DuplicateKeyError](https://github.com/go-faster/yaml/blob/v0.4.6/errors.go#L68) ¶
    
    
    type DuplicateKeyError struct {
    	First, Second *Node
    }

DuplicateKeyError reports a duplicate key. 

####  func (*DuplicateKeyError) [Error](https://github.com/go-faster/yaml/blob/v0.4.6/errors.go#L81) ¶
    
    
    func (d *DuplicateKeyError) Error() [string](/builtin#string)

Error returns the error message. 

####  type [Encoder](https://github.com/go-faster/yaml/blob/v0.4.6/marshal.go#L71) ¶
    
    
    type Encoder struct {
    	// contains filtered or unexported fields
    }

An Encoder writes YAML values to an output stream. 

####  func [NewEncoder](https://github.com/go-faster/yaml/blob/v0.4.6/marshal.go#L78) ¶
    
    
    func NewEncoder(w [io](/io).[Writer](/io#Writer)) *Encoder

NewEncoder returns a new encoder that writes to w. The Encoder should be closed after use to flush all data to w. 

####  func (*Encoder) [Close](https://github.com/go-faster/yaml/blob/v0.4.6/marshal.go#L125) ¶
    
    
    func (e *Encoder) Close() (err [error](/builtin#error))

Close closes the encoder by writing any remaining data. It does not write a stream terminating string "...". 

####  func (*Encoder) [Encode](https://github.com/go-faster/yaml/blob/v0.4.6/marshal.go#L91) ¶
    
    
    func (e *Encoder) Encode(v [any](/builtin#any)) (err [error](/builtin#error))

Encode writes the YAML encoding of v to the stream. If multiple items are encoded to the stream, the second and subsequent document will be preceded with a "---" document separator, but the first will not. 

See the documentation for Marshal for details about the conversion of Go values to YAML. 

####  func (*Encoder) [SetIndent](https://github.com/go-faster/yaml/blob/v0.4.6/marshal.go#L116) ¶
    
    
    func (e *Encoder) SetIndent(spaces [int](/builtin#int))

SetIndent changes the used indentation used when encoding. 

####  type [IsZeroer](https://github.com/go-faster/yaml/blob/v0.4.6/zero.go#L9) ¶
    
    
    type IsZeroer interface {
    	IsZero() [bool](/builtin#bool)
    }

IsZeroer is used to check whether an object is zero to determine whether it should be omitted when marshaling with the omitempty flag. One notable implementation is time.Time. 

####  type [Kind](https://github.com/go-faster/yaml/blob/v0.4.6/node.go#L10) ¶
    
    
    type Kind [uint32](/builtin#uint32)

Kind defines the kind of node. 
    
    
    const (
    	DocumentNode Kind = 1 << [iota](/builtin#iota)
    	SequenceNode
    	MappingNode
    	ScalarNode
    	AliasNode
    )

####  func (Kind) [String](https://github.com/go-faster/yaml/blob/v0.4.6/node.go#L21) ¶
    
    
    func (s Kind) String() [string](/builtin#string)

String implements fmt.Stringer. 

####  type [MarshalError](https://github.com/go-faster/yaml/blob/v0.4.6/errors.go#L124) ¶
    
    
    type MarshalError struct {
    	Msg [string](/builtin#string)
    }

MarshalError is an error that occurs during marshaling. 

####  func (*MarshalError) [Error](https://github.com/go-faster/yaml/blob/v0.4.6/errors.go#L129) ¶
    
    
    func (s *MarshalError) Error() [string](/builtin#string)

Error returns the error message. 

####  type [Marshaler](https://github.com/go-faster/yaml/blob/v0.4.6/marshal.go#L14) ¶
    
    
    type Marshaler interface {
    	MarshalYAML() ([any](/builtin#any), [error](/builtin#error))
    }

The Marshaler interface may be implemented by types to customize their behavior when being marshaled into a YAML document. The returned value is marshaled in place of the original value implementing Marshaler. 

If an error is returned by MarshalYAML, the marshaling procedure stops and returns with the provided error. 

####  type [Node](https://github.com/go-faster/yaml/blob/v0.4.6/node.go#L99) ¶
    
    
    type Node struct {
    	// Kind defines whether the node is a document, a mapping, a sequence,
    	// a scalar value, or an alias to another node. The specific data type of
    	// scalar nodes may be obtained via the ShortTag and LongTag methods.
    	Kind Kind
    
    	// Style allows customizing the appearance of the node in the tree.
    	Style Style
    
    	// Tag holds the YAML tag defining the data type for the value.
    	// When decoding, this field will always be set to the resolved tag,
    	// even when it wasn't explicitly provided in the YAML content.
    	// When encoding, if this field is unset the value type will be
    	// implied from the node properties, and if it is set, it will only
    	// be serialized into the representation if TaggedStyle is used or
    	// the implicit tag diverges from the provided one.
    	Tag [string](/builtin#string)
    
    	// Value holds the unescaped and unquoted representation of the value.
    	Value [string](/builtin#string)
    
    	// Anchor holds the anchor name for this node, which allows aliases to point to it.
    	Anchor [string](/builtin#string)
    
    	// Alias holds the node that this alias points to. Only valid when Kind is AliasNode.
    	Alias *Node
    
    	// Content holds contained nodes for documents, mappings, and sequences.
    	Content []*Node
    
    	// HeadComment holds any comments in the lines preceding the node and
    	// not separated by an empty line.
    	HeadComment [string](/builtin#string)
    
    	// LineComment holds any comments at the end of the line where the node is in.
    	LineComment [string](/builtin#string)
    
    	// FootComment holds any comments following the node and before empty lines.
    	FootComment [string](/builtin#string)
    
    	// Line and Column hold the node position in the decoded YAML text.
    	// These fields are not respected when encoding the node.
    	Line   [int](/builtin#int)
    	Column [int](/builtin#int)
    }

Node represents an element in the YAML document hierarchy. While documents are typically encoded and decoded into higher level types, such as structs and maps, Node is an intermediate representation that allows detailed control over the content being decoded or encoded. 

It's worth noting that although Node offers access into details such as line numbers, colums, and comments, the content when re-encoded will not have its original textual representation preserved. An effort is made to render the data plesantly, and to preserve comments near the data they describe, though. 

Values that make use of the Node type interact with the yaml package in the same way any other type would do, by encoding and decoding yaml data directly or indirectly into them. 

For example: 
    
    
    var person struct {
            Name    string
            Address yaml.Node
    }
    err := yaml.Unmarshal(data, &person)
    

Or by itself: 
    
    
    var person Node
    err := yaml.Unmarshal(data, &person)
    

####  func (*Node) [Decode](https://github.com/go-faster/yaml/blob/v0.4.6/unmarshal.go#L109) ¶
    
    
    func (n *Node) Decode(v [any](/builtin#any)) (err [error](/builtin#error))

Decode decodes the node and stores its data into the value pointed to by v. 

See the documentation for Unmarshal for details about the conversion of YAML into a Go value. 

####  func (*Node) [Encode](https://github.com/go-faster/yaml/blob/v0.4.6/marshal.go#L101) ¶
    
    
    func (n *Node) Encode(v [any](/builtin#any)) (err [error](/builtin#error))

Encode encodes value v and stores its representation in n. 

See the documentation for Marshal for details about the conversion of Go values into YAML. 

####  func (*Node) [EncodeJSON](https://github.com/go-faster/yaml/blob/v0.4.6/json.go#L126) ¶
    
    
    func (n *Node) EncodeJSON(e *[jx](/github.com/go-faster/jx).[Encoder](/github.com/go-faster/jx#Encoder)) [error](/builtin#error)

EncodeJSON writes the JSON representation of the node to given encoder. 

####  func (*Node) [IsZero](https://github.com/go-faster/yaml/blob/v0.4.6/node.go#L146) ¶
    
    
    func (n *Node) IsZero() [bool](/builtin#bool)

IsZero returns whether the node has all of its fields unset. 

####  func (*Node) [LongTag](https://github.com/go-faster/yaml/blob/v0.4.6/node.go#L154) ¶
    
    
    func (n *Node) LongTag() [string](/builtin#string)

LongTag returns the long form of the tag that indicates the data type for the node. If the Tag field isn't explicitly defined, one will be computed based on the node properties. 

####  func (*Node) [SetString](https://github.com/go-faster/yaml/blob/v0.4.6/node.go#L197) ¶
    
    
    func (n *Node) SetString(s [string](/builtin#string))

SetString is a convenience function that sets the node to a string value and defines its style in a pleasant way depending on its content. 

####  func (*Node) [ShortTag](https://github.com/go-faster/yaml/blob/v0.4.6/node.go#L161) ¶
    
    
    func (n *Node) ShortTag() [string](/builtin#string)

ShortTag returns the short form of the YAML tag that indicates data type for the node. If the Tag field isn't explicitly defined, one will be computed based on the node properties. 

####  type [Style](https://github.com/go-faster/yaml/blob/v0.4.6/node.go#L39) ¶
    
    
    type Style [uint32](/builtin#uint32)

Style describes the style of a node. 
    
    
    const (
    	TaggedStyle Style = 1 << [iota](/builtin#iota)
    	DoubleQuotedStyle
    	SingleQuotedStyle
    	LiteralStyle
    	FoldedStyle
    	FlowStyle
    )

####  func (Style) [String](https://github.com/go-faster/yaml/blob/v0.4.6/node.go#L51) ¶
    
    
    func (s Style) String() [string](/builtin#string)

String implements fmt.Stringer. 

####  type [SyntaxError](https://github.com/go-faster/yaml/blob/v0.4.6/errors.go#L18) ¶
    
    
    type SyntaxError struct {
    	Offset [int](/builtin#int)
    	Line   [int](/builtin#int)
    	Column [int](/builtin#int)
    	Msg    [string](/builtin#string)
    }

SyntaxError is an error that occurs during parsing. 

####  func (*SyntaxError) [Error](https://github.com/go-faster/yaml/blob/v0.4.6/errors.go#L35) ¶
    
    
    func (s *SyntaxError) Error() [string](/builtin#string)

Error returns the error message. 

####  type [TypeError](https://github.com/go-faster/yaml/blob/v0.4.6/errors.go#L140) ¶
    
    
    type TypeError struct {
    	Group [error](/builtin#error)
    }

A TypeError is returned by Unmarshal when one or more fields in the YAML document cannot be properly decoded into the requested types. When this error is returned, the value is still unmarshaled partially. 

Group is a multi-error which contains all errors that occurred. Use multierr.Errors to get a list of all errors. 

####  func (*TypeError) [Error](https://github.com/go-faster/yaml/blob/v0.4.6/errors.go#L150) ¶
    
    
    func (e *TypeError) Error() [string](/builtin#string)

Error returns the error message. 

####  func (*TypeError) [Unwrap](https://github.com/go-faster/yaml/blob/v0.4.6/errors.go#L145) ¶
    
    
    func (e *TypeError) Unwrap() [error](/builtin#error)

Unwrap returns the underlying error. 

####  type [UnknownFieldError](https://github.com/go-faster/yaml/blob/v0.4.6/errors.go#L49) ¶
    
    
    type UnknownFieldError struct {
    	Field [string](/builtin#string)
    	Type  [reflect](/reflect).[Type](/reflect#Type)
    }

UnknownFieldError reports an unknown field. 

####  func (*UnknownFieldError) [Error](https://github.com/go-faster/yaml/blob/v0.4.6/errors.go#L55) ¶
    
    
    func (d *UnknownFieldError) Error() [string](/builtin#string)

Error returns the error message. 

####  type [UnmarshalError](https://github.com/go-faster/yaml/blob/v0.4.6/errors.go#L95) ¶
    
    
    type UnmarshalError struct {
    	Node *Node
    	Type [reflect](/reflect).[Type](/reflect#Type)
    	Err  [error](/builtin#error)
    }

UnmarshalError is an error that occurs during unmarshaling. 

####  func (*UnmarshalError) [Error](https://github.com/go-faster/yaml/blob/v0.4.6/errors.go#L115) ¶
    
    
    func (s *UnmarshalError) Error() [string](/builtin#string)

Error returns the error message. 

####  func (*UnmarshalError) [Unwrap](https://github.com/go-faster/yaml/blob/v0.4.6/errors.go#L110) ¶
    
    
    func (s *UnmarshalError) Unwrap() [error](/builtin#error)

Unwrap returns the underlying error. 

####  type [Unmarshaler](https://github.com/go-faster/yaml/blob/v0.4.6/unmarshal.go#L12) ¶
    
    
    type Unmarshaler interface {
    	UnmarshalYAML(value *Node) [error](/builtin#error)
    }

The Unmarshaler interface may be implemented by types to customize their behavior when being unmarshaled from a YAML document. 
