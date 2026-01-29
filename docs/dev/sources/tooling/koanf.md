# koanf

> Auto-fetched from [https://pkg.go.dev/github.com/knadh/koanf/v2](https://pkg.go.dev/github.com/knadh/koanf/v2)
> Last Updated: 2026-01-29T20:11:21.210215+00:00

---

Index
¶
type Conf
type KeyMap
type Koanf
func New(delim string) *Koanf
func NewWithConf(conf Conf) *Koanf
func (ko *Koanf) All() map[string]any
func (ko *Koanf) Bool(path string) bool
func (ko *Koanf) BoolMap(path string) map[string]bool
func (ko *Koanf) Bools(path string) []bool
func (ko *Koanf) Bytes(path string) []byte
func (ko *Koanf) Copy() *Koanf
func (ko *Koanf) Cut(path string) *Koanf
func (ko *Koanf) Delete(path string)
func (ko *Koanf) Delim() string
func (ko *Koanf) Duration(path string) time.Duration
func (ko *Koanf) Exists(path string) bool
func (ko *Koanf) Float64(path string) float64
func (ko *Koanf) Float64Map(path string) map[string]float64
func (ko *Koanf) Float64s(path string) []float64
func (ko *Koanf) Get(path string) any
func (ko *Koanf) Int(path string) int
func (ko *Koanf) Int64(path string) int64
func (ko *Koanf) Int64Map(path string) map[string]int64
func (ko *Koanf) Int64s(path string) []int64
func (ko *Koanf) IntMap(path string) map[string]int
func (ko *Koanf) Ints(path string) []int
func (ko *Koanf) KeyMap() KeyMap
func (ko *Koanf) Keys() []string
func (ko *Koanf) Load(p Provider, pa Parser, opts ...Option) error
func (ko *Koanf) MapKeys(path string) []string
func (ko *Koanf) Marshal(p Parser) ([]byte, error)
func (ko *Koanf) Merge(in *Koanf) error
func (ko *Koanf) MergeAt(in *Koanf, path string) error
func (ko *Koanf) MustBoolMap(path string) map[string]bool
func (ko *Koanf) MustBools(path string) []bool
func (ko *Koanf) MustBytes(path string) []byte
func (ko *Koanf) MustDuration(path string) time.Duration
func (ko *Koanf) MustFloat64(path string) float64
func (ko *Koanf) MustFloat64Map(path string) map[string]float64
func (ko *Koanf) MustFloat64s(path string) []float64
func (ko *Koanf) MustInt(path string) int
func (ko *Koanf) MustInt64(path string) int64
func (ko *Koanf) MustInt64Map(path string) map[string]int64
func (ko *Koanf) MustInt64s(path string) []int64
func (ko *Koanf) MustIntMap(path string) map[string]int
func (ko *Koanf) MustInts(path string) []int
func (ko *Koanf) MustString(path string) string
func (ko *Koanf) MustStringMap(path string) map[string]string
func (ko *Koanf) MustStrings(path string) []string
func (ko *Koanf) MustStringsMap(path string) map[string][]string
func (ko *Koanf) MustTime(path, layout string) time.Time
func (ko *Koanf) Print()
func (ko *Koanf) Raw() map[string]any
func (ko *Koanf) Set(key string, val any) error
func (ko *Koanf) Slices(path string) []*Koanf
func (ko *Koanf) Sprint() string
func (ko *Koanf) String(path string) string
func (ko *Koanf) StringMap(path string) map[string]string
func (ko *Koanf) Strings(path string) []string
func (ko *Koanf) StringsMap(path string) map[string][]string
func (ko *Koanf) Time(path, layout string) time.Time
func (ko *Koanf) Unmarshal(path string, o any) error
func (ko *Koanf) UnmarshalWithConf(path string, o any, c UnmarshalConf) error
type Option
func WithMergeFunc(merge func(src, dest map[string]any) error) Option
type Parser
type Provider
type UnmarshalConf
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
Conf
¶
type Conf struct {
// Delim is the delimiter to use
// when specifying config key paths, for instance a . for `parent.child.key`
// or a / for `parent/child/key`.
Delim
string
// StrictMerge makes the merging behavior strict.
// Meaning when loading two files that have the same key,
// the first loaded file will define the desired type, and if the second file loads
// a different type will cause an error.
StrictMerge
bool
}
Conf is the Koanf configuration.
type
KeyMap
¶
type KeyMap map[
string
][]
string
KeyMap represents a map of flattened delimited keys and the non-delimited
parts as their slices. For nested keys, the map holds all levels of path combinations.
For example, the nested structure `parent -> child -> key` will produce the map:
parent.child.key => [parent, child, key]
parent.child => [parent, child]
parent => [parent]
type
Koanf
¶
type Koanf struct {
// contains filtered or unexported fields
}
Koanf is the configuration apparatus.
func
New
¶
func New(delim
string
) *
Koanf
New returns a new instance of Koanf. delim is the delimiter to use
when specifying config key paths, for instance a . for `parent.child.key`
or a / for `parent/child/key`.
func
NewWithConf
¶
func NewWithConf(conf
Conf
) *
Koanf
NewWithConf returns a new instance of Koanf based on the Conf.
func (*Koanf)
All
¶
func (ko *
Koanf
) All() map[
string
]
any
All returns a map of all flattened key paths and their values.
Note that it uses maps.Copy to create a copy that uses
json.Marshal which changes the numeric types to float64.
func (*Koanf)
Bool
¶
func (ko *
Koanf
) Bool(path
string
)
bool
Bool returns the bool value of a given key path or false if the path
does not exist or if the value is not a valid bool representation.
Accepted string representations of bool are the ones supported by strconv.ParseBool.
func (*Koanf)
BoolMap
¶
func (ko *
Koanf
) BoolMap(path
string
) map[
string
]
bool
BoolMap returns the map[string]bool value of a given key path
or an empty map[string]bool if the path does not exist or if the
value is not a valid bool map.
func (*Koanf)
Bools
¶
func (ko *
Koanf
) Bools(path
string
) []
bool
Bools returns the []bool slice value of a given key path or an
empty []bool slice if the path does not exist or if the value
is not a valid bool slice.
func (*Koanf)
Bytes
¶
func (ko *
Koanf
) Bytes(path
string
) []
byte
Bytes returns the []byte value of a given key path or an empty
[]byte slice if the path does not exist or if the value is not a valid string.
func (*Koanf)
Copy
¶
func (ko *
Koanf
) Copy() *
Koanf
Copy returns a copy of the Koanf instance.
func (*Koanf)
Cut
¶
func (ko *
Koanf
) Cut(path
string
) *
Koanf
Cut cuts the config map at a given key path into a sub map and
returns a new Koanf instance with the cut config map loaded.
For instance, if the loaded config has a path that looks like
parent.child.sub.a.b, `Cut("parent.child")` returns a new Koanf
instance with the config map `sub.a.b` where everything above
`parent.child` are cut out.
func (*Koanf)
Delete
¶
func (ko *
Koanf
) Delete(path
string
)
Delete removes all nested values from a given path.
Clears all keys/values if no path is specified.
Every empty, key on the path, is recursively deleted.
func (*Koanf)
Delim
¶
func (ko *
Koanf
) Delim()
string
Delim returns delimiter in used by this instance of Koanf.
func (*Koanf)
Duration
¶
func (ko *
Koanf
) Duration(path
string
)
time
.
Duration
Duration returns the time.Duration value of a given key path assuming
that the key contains a valid numeric value.
func (*Koanf)
Exists
¶
func (ko *
Koanf
) Exists(path
string
)
bool
Exists returns true if the given key path exists in the conf map.
func (*Koanf)
Float64
¶
func (ko *
Koanf
) Float64(path
string
)
float64
Float64 returns the float64 value of a given key path or 0 if the path
does not exist or if the value is not a valid float64.
func (*Koanf)
Float64Map
¶
func (ko *
Koanf
) Float64Map(path
string
) map[
string
]
float64
Float64Map returns the map[string]float64 value of a given key path
or an empty map[string]float64 if the path does not exist or if the
value is not a valid float64 map.
func (*Koanf)
Float64s
¶
func (ko *
Koanf
) Float64s(path
string
) []
float64
Float64s returns the []float64 slice value of a given key path or an
empty []float64 slice if the path does not exist or if the value
is not a valid float64 slice.
func (*Koanf)
Get
¶
func (ko *
Koanf
) Get(path
string
)
any
Get returns the raw, uncast any value of a given key path
in the config map. If the key path does not exist, nil is returned.
func (*Koanf)
Int
¶
func (ko *
Koanf
) Int(path
string
)
int
Int returns the int value of a given key path or 0 if the path
does not exist or if the value is not a valid int.
func (*Koanf)
Int64
¶
func (ko *
Koanf
) Int64(path
string
)
int64
Int64 returns the int64 value of a given key path or 0 if the path
does not exist or if the value is not a valid int64.
func (*Koanf)
Int64Map
¶
func (ko *
Koanf
) Int64Map(path
string
) map[
string
]
int64
Int64Map returns the map[string]int64 value of a given key path
or an empty map[string]int64 if the path does not exist or if the
value is not a valid int64 map.
func (*Koanf)
Int64s
¶
func (ko *
Koanf
) Int64s(path
string
) []
int64
Int64s returns the []int64 slice value of a given key path or an
empty []int64 slice if the path does not exist or if the value
is not a valid int slice.
func (*Koanf)
IntMap
¶
func (ko *
Koanf
) IntMap(path
string
) map[
string
]
int
IntMap returns the map[string]int value of a given key path
or an empty map[string]int if the path does not exist or if the
value is not a valid int map.
func (*Koanf)
Ints
¶
func (ko *
Koanf
) Ints(path
string
) []
int
Ints returns the []int slice value of a given key path or an
empty []int slice if the path does not exist or if the value
is not a valid int slice.
func (*Koanf)
KeyMap
¶
func (ko *
Koanf
) KeyMap()
KeyMap
KeyMap returns a map of flattened keys and the individual parts of the
key as slices. eg: "parent.child.key" => ["parent", "child", "key"].
func (*Koanf)
Keys
¶
func (ko *
Koanf
) Keys() []
string
Keys returns the slice of all flattened keys in the loaded configuration
sorted alphabetically.
func (*Koanf)
Load
¶
func (ko *
Koanf
) Load(p
Provider
, pa
Parser
, opts ...
Option
)
error
Load takes a Provider that either provides a parsed config map[string]any
in which case pa (Parser) can be nil, or raw bytes to be parsed, where a Parser
can be provided to parse. Additionally, options can be passed which modify the
load behavior, such as passing a custom merge function.
func (*Koanf)
MapKeys
¶
func (ko *
Koanf
) MapKeys(path
string
) []
string
MapKeys returns a sorted string list of keys in a map addressed by the
given path. If the path is not a map, an empty string slice is
returned.
func (*Koanf)
Marshal
¶
func (ko *
Koanf
) Marshal(p
Parser
) ([]
byte
,
error
)
Marshal takes a Parser implementation and marshals the config map into bytes,
for example, to TOML or JSON bytes.
func (*Koanf)
Merge
¶
func (ko *
Koanf
) Merge(in *
Koanf
)
error
Merge merges the config map of a given Koanf instance into
the current instance.
func (*Koanf)
MergeAt
¶
func (ko *
Koanf
) MergeAt(in *
Koanf
, path
string
)
error
MergeAt merges the config map of a given Koanf instance into
the current instance as a sub map, at the given key path.
If all or part of the key path is missing, it will be created.
If the key path is `""`, this is equivalent to Merge.
func (*Koanf)
MustBoolMap
¶
func (ko *
Koanf
) MustBoolMap(path
string
) map[
string
]
bool
MustBoolMap returns the map[string]bool value of a given key path or panics
if the value is not set or set to default value.
func (*Koanf)
MustBools
¶
func (ko *
Koanf
) MustBools(path
string
) []
bool
MustBools returns the []bool value of a given key path or panics
if the value is not set or set to default value.
func (*Koanf)
MustBytes
¶
func (ko *
Koanf
) MustBytes(path
string
) []
byte
MustBytes returns the []byte value of a given key path or panics
if the value is not set or set to default value.
func (*Koanf)
MustDuration
¶
func (ko *
Koanf
) MustDuration(path
string
)
time
.
Duration
MustDuration returns the time.Duration value of a given key path or panics
if it isn't set or set to default value 0.
func (*Koanf)
MustFloat64
¶
func (ko *
Koanf
) MustFloat64(path
string
)
float64
MustFloat64 returns the float64 value of a given key path or panics
if it isn't set or set to default value 0.
func (*Koanf)
MustFloat64Map
¶
func (ko *
Koanf
) MustFloat64Map(path
string
) map[
string
]
float64
MustFloat64Map returns the map[string]float64 value of a given key path or panics
if the value is not set or set to default value.
func (*Koanf)
MustFloat64s
¶
func (ko *
Koanf
) MustFloat64s(path
string
) []
float64
MustFloat64s returns the []Float64 slice value of a given key path or panics
if the value is not set or set to default value.
func (*Koanf)
MustInt
¶
func (ko *
Koanf
) MustInt(path
string
)
int
MustInt returns the int value of a given key path or panics
if it isn't set or set to default value of 0.
func (*Koanf)
MustInt64
¶
func (ko *
Koanf
) MustInt64(path
string
)
int64
MustInt64 returns the int64 value of a given key path or panics
if the value is not set or set to default value of 0.
func (*Koanf)
MustInt64Map
¶
func (ko *
Koanf
) MustInt64Map(path
string
) map[
string
]
int64
MustInt64Map returns the map[string]int64 value of a given key path
or panics if it isn't set or set to default value.
func (*Koanf)
MustInt64s
¶
func (ko *
Koanf
) MustInt64s(path
string
) []
int64
MustInt64s returns the []int64 slice value of a given key path or panics
if the value is not set or its default value.
func (*Koanf)
MustIntMap
¶
func (ko *
Koanf
) MustIntMap(path
string
) map[
string
]
int
MustIntMap returns the map[string]int value of a given key path or panics
if the value is not set or set to default value.
func (*Koanf)
MustInts
¶
func (ko *
Koanf
) MustInts(path
string
) []
int
MustInts returns the []int slice value of a given key path or panics
if the value is not set or set to default value.
func (*Koanf)
MustString
¶
func (ko *
Koanf
) MustString(path
string
)
string
MustString returns the string value of a given key path
or panics if it isn't set or set to default value "".
func (*Koanf)
MustStringMap
¶
func (ko *
Koanf
) MustStringMap(path
string
) map[
string
]
string
MustStringMap returns the map[string]string value of a given key path or panics
if the value is not set or set to default value.
func (*Koanf)
MustStrings
¶
func (ko *
Koanf
) MustStrings(path
string
) []
string
MustStrings returns the []string slice value of a given key path or panics
if the value is not set or set to default value.
func (*Koanf)
MustStringsMap
¶
func (ko *
Koanf
) MustStringsMap(path
string
) map[
string
][]
string
MustStringsMap returns the map[string][]string value of a given key path or panics
if the value is not set or set to default value.
func (*Koanf)
MustTime
¶
func (ko *
Koanf
) MustTime(path, layout
string
)
time
.
Time
MustTime attempts to parse the value of a given key path and return time.Time
representation. If the value is numeric, it is treated as a UNIX timestamp
and if it's string, a parse is attempted with the given layout. It panics if
the parsed time is zero.
func (*Koanf)
Print
¶
func (ko *
Koanf
) Print()
Print prints a key -> value string representation
of the config map with keys sorted alphabetically.
func (*Koanf)
Raw
¶
func (ko *
Koanf
) Raw() map[
string
]
any
Raw returns a copy of the full raw conf map.
Note that it uses maps.Copy to create a copy that uses
json.Marshal which changes the numeric types to float64.
func (*Koanf)
Set
¶
func (ko *
Koanf
) Set(key
string
, val
any
)
error
Set sets the value at a specific key.
func (*Koanf)
Slices
¶
func (ko *
Koanf
) Slices(path
string
) []*
Koanf
Slices returns a list of Koanf instances constructed out of a
[]map[string]any interface at the given path.
func (*Koanf)
Sprint
¶
func (ko *
Koanf
) Sprint()
string
Sprint returns a key -> value string representation
of the config map with keys sorted alphabetically.
func (*Koanf)
String
¶
func (ko *
Koanf
) String(path
string
)
string
String returns the string value of a given key path or "" if the path
does not exist or if the value is not a valid string.
func (*Koanf)
StringMap
¶
func (ko *
Koanf
) StringMap(path
string
) map[
string
]
string
StringMap returns the map[string]string value of a given key path
or an empty map[string]string if the path does not exist or if the
value is not a valid string map.
func (*Koanf)
Strings
¶
func (ko *
Koanf
) Strings(path
string
) []
string
Strings returns the []string slice value of a given key path or an
empty []string slice if the path does not exist or if the value
is not a valid string slice.
func (*Koanf)
StringsMap
¶
func (ko *
Koanf
) StringsMap(path
string
) map[
string
][]
string
StringsMap returns the map[string][]string value of a given key path
or an empty map[string][]string if the path does not exist or if the
value is not a valid strings map.
func (*Koanf)
Time
¶
func (ko *
Koanf
) Time(path, layout
string
)
time
.
Time
Time attempts to parse the value of a given key path and return time.Time
representation. If the value is numeric, it is treated as a UNIX timestamp
and if it's string, a parse is attempted with the given layout.
func (*Koanf)
Unmarshal
¶
func (ko *
Koanf
) Unmarshal(path
string
, o
any
)
error
Unmarshal unmarshals a given key path into the given struct using
the mapstructure lib. If no path is specified, the whole map is unmarshalled.
`koanf` is the struct field tag used to match field names. To customize,
use UnmarshalWithConf(). It uses the mitchellh/mapstructure package.
func (*Koanf)
UnmarshalWithConf
¶
func (ko *
Koanf
) UnmarshalWithConf(path
string
, o
any
, c
UnmarshalConf
)
error
UnmarshalWithConf is like Unmarshal but takes configuration params in UnmarshalConf.
See mitchellh/mapstructure's DecoderConfig for advanced customization
of the unmarshal behaviour.
type
Option
¶
type Option func(*options)
Option is a generic type used to modify the behavior of Koanf.Load.
func
WithMergeFunc
¶
func WithMergeFunc(merge func(src, dest map[
string
]
any
)
error
)
Option
WithMergeFunc is an option to modify the merge behavior of Koanf.Load.
If unset, the default merge function is used.
The merge function is expected to merge map src into dest (left to right).
type
Parser
¶
type Parser interface {
Unmarshal([]
byte
) (map[
string
]
any
,
error
)
Marshal(map[
string
]
any
) ([]
byte
,
error
)
}
Parser represents a configuration format parser.
type
Provider
¶
type Provider interface {
// ReadBytes returns the entire configuration as raw []bytes to be parsed.
// with a Parser.
ReadBytes() ([]
byte
,
error
)
// Read returns the parsed configuration as a nested map[string]any.
// It is important to note that the string keys should not be flat delimited
// keys like `parent.child.key`, but nested like `{parent: {child: {key: 1}}}`.
Read() (map[
string
]
any
,
error
)
}
Provider represents a configuration provider. Providers can
read configuration from a source (file, HTTP etc.)
type
UnmarshalConf
¶
type UnmarshalConf struct {
// Tag is the struct field tag to unmarshal.
// `koanf` is used if left empty.
Tag
string
// If this is set to true, instead of unmarshalling nested structures
// based on the key path, keys are taken literally to unmarshal into
// a flat struct. For example:
// “`
// type MyStuff struct {
// 	Child1Name string `koanf:"parent1.child1.name"`
// 	Child2Name string `koanf:"parent2.child2.name"`
// 	Type       string `koanf:"json"`
// }
// “`
FlatPaths
bool
DecoderConfig *
mapstructure
.
DecoderConfig
}
UnmarshalConf represents configuration options used by
Unmarshal() to unmarshal conf maps into arbitrary structs.