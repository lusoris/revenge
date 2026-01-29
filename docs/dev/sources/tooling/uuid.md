# google/uuid

> Auto-fetched from [https://pkg.go.dev/github.com/google/uuid](https://pkg.go.dev/github.com/google/uuid)
> Last Updated: 2026-01-28T21:44:06.842537+00:00

---

Overview
¶
Package uuid generates and inspects UUIDs.
UUIDs are based on
RFC 4122
and DCE 1.1: Authentication and Security
Services.
A UUID is a 16 byte (128 bit) array.  UUIDs may be used as keys to
maps or compared directly.
Index
¶
Constants
Variables
func ClockSequence() int
func DisableRandPool()
func EnableRandPool()
func IsInvalidLengthError(err error) bool
func NewString() string
func NodeID() []byte
func NodeInterface() string
func SetClockSequence(seq int)
func SetNodeID(id []byte) bool
func SetNodeInterface(name string) bool
func SetRand(r io.Reader)
func Validate(s string) error
type Domain
func (d Domain) String() string
type NullUUID
func (nu NullUUID) MarshalBinary() ([]byte, error)
func (nu NullUUID) MarshalJSON() ([]byte, error)
func (nu NullUUID) MarshalText() ([]byte, error)
func (nu *NullUUID) Scan(value interface{}) error
func (nu *NullUUID) UnmarshalBinary(data []byte) error
func (nu *NullUUID) UnmarshalJSON(data []byte) error
func (nu *NullUUID) UnmarshalText(data []byte) error
func (nu NullUUID) Value() (driver.Value, error)
type Time
func GetTime() (Time, uint16, error)
func (t Time) UnixTime() (sec, nsec int64)
type UUID
func FromBytes(b []byte) (uuid UUID, err error)
func Must(uuid UUID, err error) UUID
func MustParse(s string) UUID
func New() UUID
func NewDCEGroup() (UUID, error)
func NewDCEPerson() (UUID, error)
func NewDCESecurity(domain Domain, id uint32) (UUID, error)
func NewHash(h hash.Hash, space UUID, data []byte, version int) UUID
func NewMD5(space UUID, data []byte) UUID
func NewRandom() (UUID, error)
func NewRandomFromReader(r io.Reader) (UUID, error)
func NewSHA1(space UUID, data []byte) UUID
func NewUUID() (UUID, error)
func NewV6() (UUID, error)
func NewV7() (UUID, error)
func NewV7FromReader(r io.Reader) (UUID, error)
func Parse(s string) (UUID, error)
func ParseBytes(b []byte) (UUID, error)
func (uuid UUID) ClockSequence() int
func (uuid UUID) Domain() Domain
func (uuid UUID) ID() uint32
func (uuid UUID) MarshalBinary() ([]byte, error)
func (uuid UUID) MarshalText() ([]byte, error)
func (uuid UUID) NodeID() []byte
func (uuid *UUID) Scan(src interface{}) error
func (uuid UUID) String() string
func (uuid UUID) Time() Time
func (uuid UUID) URN() string
func (uuid *UUID) UnmarshalBinary(data []byte) error
func (uuid *UUID) UnmarshalText(data []byte) error
func (uuid UUID) Value() (driver.Value, error)
func (uuid UUID) Variant() Variant
func (uuid UUID) Version() Version
type UUIDs
func (uuids UUIDs) Strings() []string
type Variant
func (v Variant) String() string
type Version
func (v Version) String() string
Constants
¶
View Source
const (
Person =
Domain
(0)
Group  =
Domain
(1)
Org    =
Domain
(2)
)
Domain constants for DCE Security (Version 2) UUIDs.
View Source
const (
Invalid   =
Variant
(
iota
)
// Invalid UUID
RFC4122
// The variant specified in RFC4122
Reserved
// Reserved, NCS backward compatibility.
Microsoft
// Reserved, Microsoft Corporation backward compatibility.
Future
// Reserved for future definition.
)
Constants returned by Variant.
Variables
¶
View Source
var (
NameSpaceDNS  =
Must
(
Parse
("6ba7b810-9dad-11d1-80b4-00c04fd430c8"))
NameSpaceURL  =
Must
(
Parse
("6ba7b811-9dad-11d1-80b4-00c04fd430c8"))
NameSpaceOID  =
Must
(
Parse
("6ba7b812-9dad-11d1-80b4-00c04fd430c8"))
NameSpaceX500 =
Must
(
Parse
("6ba7b814-9dad-11d1-80b4-00c04fd430c8"))
Nil
UUID
// empty UUID, all zeros
// The Max UUID is special form of UUID that is specified to have all 128 bits set to 1.
Max =
UUID
{
0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
}
)
Well known namespace IDs and UUIDs
Functions
¶
func
ClockSequence
¶
func ClockSequence()
int
ClockSequence returns the current clock sequence, generating one if not
already set.  The clock sequence is only used for Version 1 UUIDs.
The uuid package does not use global static storage for the clock sequence or
the last time a UUID was generated.  Unless SetClockSequence is used, a new
random clock sequence is generated the first time a clock sequence is
requested by ClockSequence, GetTime, or NewUUID.  (section 4.2.1.1)
func
DisableRandPool
¶
added in
v1.3.0
func DisableRandPool()
DisableRandPool disables the randomness pool if it was previously
enabled with EnableRandPool.
Both EnableRandPool and DisableRandPool are not thread-safe and should
only be called when there is no possibility that New or any other
UUID Version 4 generation function will be called concurrently.
func
EnableRandPool
¶
added in
v1.3.0
func EnableRandPool()
EnableRandPool enables internal randomness pool used for Random
(Version 4) UUID generation. The pool contains random bytes read from
the random number generator on demand in batches. Enabling the pool
may improve the UUID generation throughput significantly.
Since the pool is stored on the Go heap, this feature may be a bad fit
for security sensitive applications.
Both EnableRandPool and DisableRandPool are not thread-safe and should
only be called when there is no possibility that New or any other
UUID Version 4 generation function will be called concurrently.
func
IsInvalidLengthError
¶
added in
v1.3.0
func IsInvalidLengthError(err
error
)
bool
IsInvalidLengthError is matcher function for custom error invalidLengthError
func
NewString
¶
added in
v1.2.0
func NewString()
string
NewString creates a new random UUID and returns it as a string or panics.
NewString is equivalent to the expression
uuid.New().String()
func
NodeID
¶
func NodeID() []
byte
NodeID returns a slice of a copy of the current Node ID, setting the Node ID
if not already set.
func
NodeInterface
¶
func NodeInterface()
string
NodeInterface returns the name of the interface from which the NodeID was
derived.  The interface "user" is returned if the NodeID was set by
SetNodeID.
func
SetClockSequence
¶
func SetClockSequence(seq
int
)
SetClockSequence sets the clock sequence to the lower 14 bits of seq.  Setting to
-1 causes a new sequence to be generated.
func
SetNodeID
¶
func SetNodeID(id []
byte
)
bool
SetNodeID sets the Node ID to be used for Version 1 UUIDs.  The first 6 bytes
of id are used.  If id is less than 6 bytes then false is returned and the
Node ID is not set.
func
SetNodeInterface
¶
func SetNodeInterface(name
string
)
bool
SetNodeInterface selects the hardware address to be used for Version 1 UUIDs.
If name is "" then the first usable interface found will be used or a random
Node ID will be generated.  If a named interface cannot be found then false
is returned.
SetNodeInterface never fails when name is "".
func
SetRand
¶
func SetRand(r
io
.
Reader
)
SetRand sets the random number generator to r, which implements io.Reader.
If r.Read returns an error when the package requests random data then
a panic will be issued.
Calling SetRand with nil sets the random number generator to the default
generator.
func
Validate
¶
added in
v1.5.0
func Validate(s
string
)
error
Validate returns an error if s is not a properly formatted UUID in one of the following formats:
xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
urn:uuid:xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
{xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx}
It returns an error if the format is invalid, otherwise nil.
Types
¶
type
Domain
¶
type Domain
byte
A Domain represents a Version 2 domain
func (Domain)
String
¶
func (d
Domain
) String()
string
type
NullUUID
¶
added in
v1.3.0
type NullUUID struct {
UUID
UUID
Valid
bool
// Valid is true if UUID is not NULL
}
NullUUID represents a UUID that may be null.
NullUUID implements the SQL driver.Scanner interface so
it can be used as a scan destination:
var u uuid.NullUUID
err := db.QueryRow("SELECT name FROM foo WHERE id=?", id).Scan(&u)
...
if u.Valid {
// use u.UUID
} else {
// NULL value
}
func (NullUUID)
MarshalBinary
¶
added in
v1.3.0
func (nu
NullUUID
) MarshalBinary() ([]
byte
,
error
)
MarshalBinary implements encoding.BinaryMarshaler.
func (NullUUID)
MarshalJSON
¶
added in
v1.3.0
func (nu
NullUUID
) MarshalJSON() ([]
byte
,
error
)
MarshalJSON implements json.Marshaler.
func (NullUUID)
MarshalText
¶
added in
v1.3.0
func (nu
NullUUID
) MarshalText() ([]
byte
,
error
)
MarshalText implements encoding.TextMarshaler.
func (*NullUUID)
Scan
¶
added in
v1.3.0
func (nu *
NullUUID
) Scan(value interface{})
error
Scan implements the SQL driver.Scanner interface.
func (*NullUUID)
UnmarshalBinary
¶
added in
v1.3.0
func (nu *
NullUUID
) UnmarshalBinary(data []
byte
)
error
UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (*NullUUID)
UnmarshalJSON
¶
added in
v1.3.0
func (nu *
NullUUID
) UnmarshalJSON(data []
byte
)
error
UnmarshalJSON implements json.Unmarshaler.
func (*NullUUID)
UnmarshalText
¶
added in
v1.3.0
func (nu *
NullUUID
) UnmarshalText(data []
byte
)
error
UnmarshalText implements encoding.TextUnmarshaler.
func (NullUUID)
Value
¶
added in
v1.3.0
func (nu
NullUUID
) Value() (
driver
.
Value
,
error
)
Value implements the driver Valuer interface.
type
Time
¶
type Time
int64
A Time represents a time as the number of 100's of nanoseconds since 15 Oct
1582.
func
GetTime
¶
func GetTime() (
Time
,
uint16
,
error
)
GetTime returns the current Time (100s of nanoseconds since 15 Oct 1582) and
clock sequence as well as adjusting the clock sequence as needed.  An error
is returned if the current time cannot be determined.
func (Time)
UnixTime
¶
func (t
Time
) UnixTime() (sec, nsec
int64
)
UnixTime converts t the number of seconds and nanoseconds using the Unix
epoch of 1 Jan 1970.
type
UUID
¶
type UUID [16]
byte
A UUID is a 128 bit (16 byte) Universal Unique IDentifier as defined in
RFC
4122
.
func
FromBytes
¶
func FromBytes(b []
byte
) (uuid
UUID
, err
error
)
FromBytes creates a new UUID from a byte slice. Returns an error if the slice
does not have a length of 16. The bytes are copied from the slice.
func
Must
¶
func Must(uuid
UUID
, err
error
)
UUID
Must returns uuid if err is nil and panics otherwise.
func
MustParse
¶
added in
v1.1.0
func MustParse(s
string
)
UUID
MustParse is like Parse but panics if the string cannot be parsed.
It simplifies safe initialization of global variables holding compiled UUIDs.
func
New
¶
func New()
UUID
New creates a new random UUID or panics.  New is equivalent to
the expression
uuid.Must(uuid.NewRandom())
func
NewDCEGroup
¶
func NewDCEGroup() (
UUID
,
error
)
NewDCEGroup returns a DCE Security (Version 2) UUID in the group
domain with the id returned by os.Getgid.
NewDCESecurity(Group, uint32(os.Getgid()))
func
NewDCEPerson
¶
func NewDCEPerson() (
UUID
,
error
)
NewDCEPerson returns a DCE Security (Version 2) UUID in the person
domain with the id returned by os.Getuid.
NewDCESecurity(Person, uint32(os.Getuid()))
func
NewDCESecurity
¶
func NewDCESecurity(domain
Domain
, id
uint32
) (
UUID
,
error
)
NewDCESecurity returns a DCE Security (Version 2) UUID.
The domain should be one of Person, Group or Org.
On a POSIX system the id should be the users UID for the Person
domain and the users GID for the Group.  The meaning of id for
the domain Org or on non-POSIX systems is site defined.
For a given domain/id pair the same token may be returned for up to
7 minutes and 10 seconds.
func
NewHash
¶
func NewHash(h
hash
.
Hash
, space
UUID
, data []
byte
, version
int
)
UUID
NewHash returns a new UUID derived from the hash of space concatenated with
data generated by h.  The hash should be at least 16 byte in length.  The
first 16 bytes of the hash are used to form the UUID.  The version of the
UUID will be the lower 4 bits of version.  NewHash is used to implement
NewMD5 and NewSHA1.
func
NewMD5
¶
func NewMD5(space
UUID
, data []
byte
)
UUID
NewMD5 returns a new MD5 (Version 3) UUID based on the
supplied name space and data.  It is the same as calling:
NewHash(md5.New(), space, data, 3)
func
NewRandom
¶
func NewRandom() (
UUID
,
error
)
NewRandom returns a Random (Version 4) UUID.
The strength of the UUIDs is based on the strength of the crypto/rand
package.
Uses the randomness pool if it was enabled with EnableRandPool.
A note about uniqueness derived from the UUID Wikipedia entry:
Randomly generated UUIDs have 122 random bits.  One's annual risk of being
hit by a meteorite is estimated to be one chance in 17 billion, that
means the probability is about 0.00000000006 (6 × 10−11),
equivalent to the odds of creating a few tens of trillions of UUIDs in a
year and having one duplicate.
func
NewRandomFromReader
¶
added in
v1.1.2
func NewRandomFromReader(r
io
.
Reader
) (
UUID
,
error
)
NewRandomFromReader returns a UUID based on bytes read from a given io.Reader.
func
NewSHA1
¶
func NewSHA1(space
UUID
, data []
byte
)
UUID
NewSHA1 returns a new SHA1 (Version 5) UUID based on the
supplied name space and data.  It is the same as calling:
NewHash(sha1.New(), space, data, 5)
func
NewUUID
¶
func NewUUID() (
UUID
,
error
)
NewUUID returns a Version 1 UUID based on the current NodeID and clock
sequence, and the current time.  If the NodeID has not been set by SetNodeID
or SetNodeInterface then it will be set automatically.  If the NodeID cannot
be set NewUUID returns nil.  If clock sequence has not been set by
SetClockSequence then it will be set automatically.  If GetTime fails to
return the current NewUUID returns nil and an error.
In most cases, New should be used.
func
NewV6
¶
added in
v1.5.0
func NewV6() (
UUID
,
error
)
UUID version 6 is a field-compatible version of UUIDv1, reordered for improved DB locality.
It is expected that UUIDv6 will primarily be used in contexts where there are existing v1 UUIDs.
Systems that do not involve legacy UUIDv1 SHOULD consider using UUIDv7 instead.
see
https://datatracker.ietf.org/doc/html/draft-peabody-dispatch-new-uuid-format-03#uuidv6
NewV6 returns a Version 6 UUID based on the current NodeID and clock
sequence, and the current time. If the NodeID has not been set by SetNodeID
or SetNodeInterface then it will be set automatically. If the NodeID cannot
be set NewV6 set NodeID is random bits automatically . If clock sequence has not been set by
SetClockSequence then it will be set automatically. If GetTime fails to
return the current NewV6 returns Nil and an error.
func
NewV7
¶
added in
v1.5.0
func NewV7() (
UUID
,
error
)
UUID version 7 features a time-ordered value field derived from the widely
implemented and well known Unix Epoch timestamp source,
the number of milliseconds seconds since midnight 1 Jan 1970 UTC, leap seconds excluded.
As well as improved entropy characteristics over versions 1 or 6.
see
https://datatracker.ietf.org/doc/html/draft-peabody-dispatch-new-uuid-format-03#name-uuid-version-7
Implementations SHOULD utilize UUID version 7 over UUID version 1 and 6 if possible.
NewV7 returns a Version 7 UUID based on the current time(Unix Epoch).
Uses the randomness pool if it was enabled with EnableRandPool.
On error, NewV7 returns Nil and an error
func
NewV7FromReader
¶
added in
v1.5.0
func NewV7FromReader(r
io
.
Reader
) (
UUID
,
error
)
NewV7FromReader returns a Version 7 UUID based on the current time(Unix Epoch).
it use NewRandomFromReader fill random bits.
On error, NewV7FromReader returns Nil and an error.
func
Parse
¶
func Parse(s
string
) (
UUID
,
error
)
Parse decodes s into a UUID or returns an error if it cannot be parsed.  Both
the standard UUID forms defined in
RFC 4122
(xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx and
urn:uuid:xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx) are decoded.  In addition,
Parse accepts non-standard strings such as the raw hex encoding
xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx and 38 byte "Microsoft style" encodings,
e.g.  {xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx}.  Only the middle 36 bytes are
examined in the latter case.  Parse should not be used to validate strings as
it parses non-standard encodings as indicated above.
func
ParseBytes
¶
func ParseBytes(b []
byte
) (
UUID
,
error
)
ParseBytes is like Parse, except it parses a byte slice instead of a string.
func (UUID)
ClockSequence
¶
func (uuid
UUID
) ClockSequence()
int
ClockSequence returns the clock sequence encoded in uuid.
The clock sequence is only well defined for version 1 and 2 UUIDs.
func (UUID)
Domain
¶
func (uuid
UUID
) Domain()
Domain
Domain returns the domain for a Version 2 UUID.  Domains are only defined
for Version 2 UUIDs.
func (UUID)
ID
¶
func (uuid
UUID
) ID()
uint32
ID returns the id for a Version 2 UUID. IDs are only defined for Version 2
UUIDs.
func (UUID)
MarshalBinary
¶
func (uuid
UUID
) MarshalBinary() ([]
byte
,
error
)
MarshalBinary implements encoding.BinaryMarshaler.
func (UUID)
MarshalText
¶
func (uuid
UUID
) MarshalText() ([]
byte
,
error
)
MarshalText implements encoding.TextMarshaler.
func (UUID)
NodeID
¶
func (uuid
UUID
) NodeID() []
byte
NodeID returns the 6 byte node id encoded in uuid.  It returns nil if uuid is
not valid.  The NodeID is only well defined for version 1 and 2 UUIDs.
func (*UUID)
Scan
¶
func (uuid *
UUID
) Scan(src interface{})
error
Scan implements sql.Scanner so UUIDs can be read from databases transparently.
Currently, database types that map to string and []byte are supported. Please
consult database-specific driver documentation for matching types.
func (UUID)
String
¶
func (uuid
UUID
) String()
string
String returns the string form of uuid, xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
, or "" if uuid is invalid.
func (UUID)
Time
¶
func (uuid
UUID
) Time()
Time
Time returns the time in 100s of nanoseconds since 15 Oct 1582 encoded in
uuid.  The time is only defined for version 1, 2, 6 and 7 UUIDs.
func (UUID)
URN
¶
func (uuid
UUID
) URN()
string
URN returns the
RFC 2141
URN form of uuid,
urn:uuid:xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx,  or "" if uuid is invalid.
func (*UUID)
UnmarshalBinary
¶
func (uuid *
UUID
) UnmarshalBinary(data []
byte
)
error
UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (*UUID)
UnmarshalText
¶
func (uuid *
UUID
) UnmarshalText(data []
byte
)
error
UnmarshalText implements encoding.TextUnmarshaler.
func (UUID)
Value
¶
func (uuid
UUID
) Value() (
driver
.
Value
,
error
)
Value implements sql.Valuer so that UUIDs can be written to databases
transparently. Currently, UUIDs map to strings. Please consult
database-specific driver documentation for matching types.
func (UUID)
Variant
¶
func (uuid
UUID
) Variant()
Variant
Variant returns the variant encoded in uuid.
func (UUID)
Version
¶
func (uuid
UUID
) Version()
Version
Version returns the version of uuid.
type
UUIDs
¶
added in
v1.4.0
type UUIDs []
UUID
UUIDs is a slice of UUID types.
func (UUIDs)
Strings
¶
added in
v1.4.0
func (uuids
UUIDs
) Strings() []
string
Strings returns a string slice containing the string form of each UUID in uuids.
type
Variant
¶
type Variant
byte
A Variant represents a UUID's variant.
func (Variant)
String
¶
func (v
Variant
) String()
string
type
Version
¶
type Version
byte
A Version represents a UUID's version.
func (Version)
String
¶
func (v
Version
) String()
string