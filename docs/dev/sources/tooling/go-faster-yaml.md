# go-faster/yaml

> Source: https://pkg.go.dev/github.com/go-faster/yaml
> Fetched: 2026-01-30T23:49:15.067596+00:00
> Content-Hash: a1b05d845a59d7ec
> Type: html

---

Overview

¶

Package yaml implements YAML support for the Go language.

Source code and other details for the project are available at GitHub:

https://github.com/go-faster/yaml

Index

¶

func Marshal(in any) (out []byte, err error)

func Unmarshal(in []byte, out any) (err error)

type Decoder

func NewDecoder(r io.Reader) *Decoder

func (dec *Decoder) Decode(v any) (err error)

func (dec *Decoder) KnownFields(enable bool)

type DuplicateKeyError

func (d *DuplicateKeyError) Error() string

type Encoder

func NewEncoder(w io.Writer) *Encoder

func (e *Encoder) Close() (err error)

func (e *Encoder) Encode(v any) (err error)

func (e *Encoder) SetIndent(spaces int)

type IsZeroer

type Kind

func (s Kind) String() string

type MarshalError

func (s *MarshalError) Error() string

type Marshaler

type Node

func (n *Node) Decode(v any) (err error)

func (n *Node) Encode(v any) (err error)

func (n *Node) EncodeJSON(e *jx.Encoder) error

func (n *Node) IsZero() bool

func (n *Node) LongTag() string

func (n *Node) SetString(s string)

func (n *Node) ShortTag() string

type Style

func (s Style) String() string

type SyntaxError

func (s *SyntaxError) Error() string

type TypeError

func (e *TypeError) Error() string

func (e *TypeError) Unwrap() error

type UnknownFieldError

func (d *UnknownFieldError) Error() string

type UnmarshalError

func (s *UnmarshalError) Error() string

func (s *UnmarshalError) Unwrap() error

type Unmarshaler

Examples

¶

Unmarshal (Embedded)

Constants

¶

This section is empty.

Variables

¶

This section is empty.

Functions

¶

func

Marshal

¶

func Marshal(in

any

) (out []

byte

, err

error

)

Marshal serializes the value provided into a YAML document. The structure
of the generated document will reflect the structure of the value itself.
Maps and pointers (to struct, string, int, etc) are accepted as the in value.

Struct fields are only marshaled if they are exported (have an upper case
first letter), and are marshaled using the field name lowercased as the
default key. Custom keys may be defined via the "yaml" name in the field
tag: the content preceding the first comma is used as the key, and the
following comma-separated options are used to tweak the marshaling process.
Conflicting names result in a runtime error.

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

func

Unmarshal

¶

func Unmarshal(in []

byte

, out

any

) (err

error

)

Unmarshal decodes the first document found within the in byte slice
and assigns decoded values into the out value.

Maps and pointers (to a struct, string, int, etc) are accepted as out
values. If an internal pointer within a struct is not initialized,
the yaml package will initialize it if necessary for unmarshaling
the provided data. The out parameter must not be nil.

The type of the decoded values should be compatible with the respective
values in out. If one or more values cannot be decoded due to a type
mismatches, decoding continues partially until the end of the YAML
content, and a *yaml.TypeError is returned with details for all
missed values.

Struct fields are only unmarshaled if they are exported (have an
upper case first letter), and are unmarshaled using the field name
lowercased as the default key. Custom keys may be defined via the
"yaml" name in the field tag: the content preceding the first comma
is used as the key, and the following comma-separated options are
used to tweak the marshaling process (see Marshal).
Conflicting names result in a runtime error.

For example:

type T struct {
    F int `yaml:"a,omitempty"`
    B int
}
var t T
yaml.Unmarshal([]byte("a: 1\nb: 2"), &t)

See the documentation of Marshal for the format of tags and a list of
supported tag options.

Example (Embedded)

¶

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

Share

Format

Run

Types

¶

type

Decoder

¶

type Decoder struct {

// contains filtered or unexported fields

}

A Decoder reads and decodes YAML values from an input stream.

func

NewDecoder

¶

func NewDecoder(r

io

.

Reader

) *

Decoder

NewDecoder returns a new decoder that reads from r.

The decoder introduces its own buffering and may read
data from r beyond the YAML values requested.

func (*Decoder)

Decode

¶

func (dec *

Decoder

) Decode(v

any

) (err

error

)

Decode reads the next YAML-encoded value from its input
and stores it in the value pointed to by v.

See the documentation for Unmarshal for details about the
conversion of YAML into a Go value.

func (*Decoder)

KnownFields

¶

func (dec *

Decoder

) KnownFields(enable

bool

)

KnownFields ensures that the keys in decoded mappings to
exist as fields in the struct being decoded into.

type

DuplicateKeyError

¶

type DuplicateKeyError struct {

First, Second *

Node

}

DuplicateKeyError reports a duplicate key.

func (*DuplicateKeyError)

Error

¶

func (d *

DuplicateKeyError

) Error()

string

Error returns the error message.

type

Encoder

¶

type Encoder struct {

// contains filtered or unexported fields

}

An Encoder writes YAML values to an output stream.

func

NewEncoder

¶

func NewEncoder(w

io

.

Writer

) *

Encoder

NewEncoder returns a new encoder that writes to w.
The Encoder should be closed after use to flush all data
to w.

func (*Encoder)

Close

¶

func (e *

Encoder

) Close() (err

error

)

Close closes the encoder by writing any remaining data.
It does not write a stream terminating string "...".

func (*Encoder)

Encode

¶

func (e *

Encoder

) Encode(v

any

) (err

error

)

Encode writes the YAML encoding of v to the stream.
If multiple items are encoded to the stream, the
second and subsequent document will be preceded
with a "---" document separator, but the first will not.

See the documentation for Marshal for details about the conversion of Go
values to YAML.

func (*Encoder)

SetIndent

¶

func (e *

Encoder

) SetIndent(spaces

int

)

SetIndent changes the used indentation used when encoding.

type

IsZeroer

¶

type IsZeroer interface {

IsZero()

bool

}

IsZeroer is used to check whether an object is zero to
determine whether it should be omitted when marshaling
with the omitempty flag. One notable implementation
is time.Time.

type

Kind

¶

type Kind

uint32

Kind defines the kind of node.

const (

DocumentNode

Kind

= 1 <<

iota

SequenceNode

MappingNode

ScalarNode

AliasNode

)

func (Kind)

String

¶

func (s

Kind

) String()

string

String implements fmt.Stringer.

type

MarshalError

¶

type MarshalError struct {

Msg

string

}

MarshalError is an error that occurs during marshaling.

func (*MarshalError)

Error

¶

func (s *

MarshalError

) Error()

string

Error returns the error message.

type

Marshaler

¶

type Marshaler interface {

MarshalYAML() (

any

,

error

)

}

The Marshaler interface may be implemented by types to customize their
behavior when being marshaled into a YAML document. The returned value
is marshaled in place of the original value implementing Marshaler.

If an error is returned by MarshalYAML, the marshaling procedure stops
and returns with the provided error.

type

Node

¶

type Node struct {

// Kind defines whether the node is a document, a mapping, a sequence,

// a scalar value, or an alias to another node. The specific data type of

// scalar nodes may be obtained via the ShortTag and LongTag methods.

Kind

Kind

// Style allows customizing the appearance of the node in the tree.

Style

Style

// Tag holds the YAML tag defining the data type for the value.

// When decoding, this field will always be set to the resolved tag,

// even when it wasn't explicitly provided in the YAML content.

// When encoding, if this field is unset the value type will be

// implied from the node properties, and if it is set, it will only

// be serialized into the representation if TaggedStyle is used or

// the implicit tag diverges from the provided one.

Tag

string

// Value holds the unescaped and unquoted representation of the value.

Value

string

// Anchor holds the anchor name for this node, which allows aliases to point to it.

Anchor

string

// Alias holds the node that this alias points to. Only valid when Kind is AliasNode.

Alias *

Node

// Content holds contained nodes for documents, mappings, and sequences.

Content []*

Node

// HeadComment holds any comments in the lines preceding the node and

// not separated by an empty line.

HeadComment

string

// LineComment holds any comments at the end of the line where the node is in.

LineComment

string

// FootComment holds any comments following the node and before empty lines.

FootComment

string

// Line and Column hold the node position in the decoded YAML text.

// These fields are not respected when encoding the node.

Line

int

Column

int

}

Node represents an element in the YAML document hierarchy. While documents
are typically encoded and decoded into higher level types, such as structs
and maps, Node is an intermediate representation that allows detailed
control over the content being decoded or encoded.

It's worth noting that although Node offers access into details such as
line numbers, colums, and comments, the content when re-encoded will not
have its original textual representation preserved. An effort is made to
render the data plesantly, and to preserve comments near the data they
describe, though.

Values that make use of the Node type interact with the yaml package in the
same way any other type would do, by encoding and decoding yaml data
directly or indirectly into them.

For example:

var person struct {
        Name    string
        Address yaml.Node
}
err := yaml.Unmarshal(data, &person)

Or by itself:

var person Node
err := yaml.Unmarshal(data, &person)

func (*Node)

Decode

¶

func (n *

Node

) Decode(v

any

) (err

error

)

Decode decodes the node and stores its data into the value pointed to by v.

See the documentation for Unmarshal for details about the
conversion of YAML into a Go value.

func (*Node)

Encode

¶

func (n *

Node

) Encode(v

any

) (err

error

)

Encode encodes value v and stores its representation in n.

See the documentation for Marshal for details about the
conversion of Go values into YAML.

func (*Node)

EncodeJSON

¶

func (n *

Node

) EncodeJSON(e *

jx

.

Encoder

)

error

EncodeJSON writes the JSON representation of the node to given encoder.

func (*Node)

IsZero

¶

func (n *

Node

) IsZero()

bool

IsZero returns whether the node has all of its fields unset.

func (*Node)

LongTag

¶

func (n *

Node

) LongTag()

string

LongTag returns the long form of the tag that indicates the data type for
the node. If the Tag field isn't explicitly defined, one will be computed
based on the node properties.

func (*Node)

SetString

¶

func (n *

Node

) SetString(s

string

)

SetString is a convenience function that sets the node to a string value
and defines its style in a pleasant way depending on its content.

func (*Node)

ShortTag

¶

func (n *

Node

) ShortTag()

string

ShortTag returns the short form of the YAML tag that indicates data type for
the node. If the Tag field isn't explicitly defined, one will be computed
based on the node properties.

type

Style

¶

type Style

uint32

Style describes the style of a node.

const (

TaggedStyle

Style

= 1 <<

iota

DoubleQuotedStyle

SingleQuotedStyle

LiteralStyle

FoldedStyle

FlowStyle

)

func (Style)

String

¶

func (s

Style

) String()

string

String implements fmt.Stringer.

type

SyntaxError

¶

type SyntaxError struct {

Offset

int

Line

int

Column

int

Msg

string

}

SyntaxError is an error that occurs during parsing.

func (*SyntaxError)

Error

¶

func (s *

SyntaxError

) Error()

string

Error returns the error message.

type

TypeError

¶

type TypeError struct {

Group

error

}

A TypeError is returned by Unmarshal when one or more fields in
the YAML document cannot be properly decoded into the requested
types. When this error is returned, the value is still
unmarshaled partially.

Group is a multi-error which contains all errors that occurred.
Use multierr.Errors to get a list of all errors.

func (*TypeError)

Error

¶

func (e *

TypeError

) Error()

string

Error returns the error message.

func (*TypeError)

Unwrap

¶

func (e *

TypeError

) Unwrap()

error

Unwrap returns the underlying error.

type

UnknownFieldError

¶

type UnknownFieldError struct {

Field

string

Type

reflect

.

Type

}

UnknownFieldError reports an unknown field.

func (*UnknownFieldError)

Error

¶

func (d *

UnknownFieldError

) Error()

string

Error returns the error message.

type

UnmarshalError

¶

type UnmarshalError struct {

Node *

Node

Type

reflect

.

Type

Err

error

}

UnmarshalError is an error that occurs during unmarshaling.

func (*UnmarshalError)

Error

¶

func (s *

UnmarshalError

) Error()

string

Error returns the error message.

func (*UnmarshalError)

Unwrap

¶

func (s *

UnmarshalError

) Unwrap()

error

Unwrap returns the underlying error.

type

Unmarshaler

¶

type Unmarshaler interface {

UnmarshalYAML(value *

Node

)

error

}

The Unmarshaler interface may be implemented by types to customize their
behavior when being unmarshaled from a YAML document.