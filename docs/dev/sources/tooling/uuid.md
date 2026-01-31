# google/uuid

> Source: https://pkg.go.dev/github.com/google/uuid
> Fetched: 2026-01-31T10:56:26.952852+00:00
> Content-Hash: aa90e02bc836df81
> Type: html

---

### Overview ¶

Package uuid generates and inspects UUIDs. 

UUIDs are based on [RFC 4122](https://rfc-editor.org/rfc/rfc4122.html) and DCE 1.1: Authentication and Security Services. 

A UUID is a 16 byte (128 bit) array. UUIDs may be used as keys to maps or compared directly. 

### Index ¶

  * Constants
  * Variables
  * func ClockSequence() int
  * func DisableRandPool()
  * func EnableRandPool()
  * func IsInvalidLengthError(err error) bool
  * func NewString() string
  * func NodeID() []byte
  * func NodeInterface() string
  * func SetClockSequence(seq int)
  * func SetNodeID(id []byte) bool
  * func SetNodeInterface(name string) bool
  * func SetRand(r io.Reader)
  * func Validate(s string) error
  * type Domain
  *     * func (d Domain) String() string
  * type NullUUID
  *     * func (nu NullUUID) MarshalBinary() ([]byte, error)
    * func (nu NullUUID) MarshalJSON() ([]byte, error)
    * func (nu NullUUID) MarshalText() ([]byte, error)
    * func (nu *NullUUID) Scan(value interface{}) error
    * func (nu *NullUUID) UnmarshalBinary(data []byte) error
    * func (nu *NullUUID) UnmarshalJSON(data []byte) error
    * func (nu *NullUUID) UnmarshalText(data []byte) error
    * func (nu NullUUID) Value() (driver.Value, error)
  * type Time
  *     * func GetTime() (Time, uint16, error)
  *     * func (t Time) UnixTime() (sec, nsec int64)
  * type UUID
  *     * func FromBytes(b []byte) (uuid UUID, err error)
    * func Must(uuid UUID, err error) UUID
    * func MustParse(s string) UUID
    * func New() UUID
    * func NewDCEGroup() (UUID, error)
    * func NewDCEPerson() (UUID, error)
    * func NewDCESecurity(domain Domain, id uint32) (UUID, error)
    * func NewHash(h hash.Hash, space UUID, data []byte, version int) UUID
    * func NewMD5(space UUID, data []byte) UUID
    * func NewRandom() (UUID, error)
    * func NewRandomFromReader(r io.Reader) (UUID, error)
    * func NewSHA1(space UUID, data []byte) UUID
    * func NewUUID() (UUID, error)
    * func NewV6() (UUID, error)
    * func NewV7() (UUID, error)
    * func NewV7FromReader(r io.Reader) (UUID, error)
    * func Parse(s string) (UUID, error)
    * func ParseBytes(b []byte) (UUID, error)
  *     * func (uuid UUID) ClockSequence() int
    * func (uuid UUID) Domain() Domain
    * func (uuid UUID) ID() uint32
    * func (uuid UUID) MarshalBinary() ([]byte, error)
    * func (uuid UUID) MarshalText() ([]byte, error)
    * func (uuid UUID) NodeID() []byte
    * func (uuid *UUID) Scan(src interface{}) error
    * func (uuid UUID) String() string
    * func (uuid UUID) Time() Time
    * func (uuid UUID) URN() string
    * func (uuid *UUID) UnmarshalBinary(data []byte) error
    * func (uuid *UUID) UnmarshalText(data []byte) error
    * func (uuid UUID) Value() (driver.Value, error)
    * func (uuid UUID) Variant() Variant
    * func (uuid UUID) Version() Version
  * type UUIDs
  *     * func (uuids UUIDs) Strings() []string
  * type Variant
  *     * func (v Variant) String() string
  * type Version
  *     * func (v Version) String() string



### Constants ¶

[View Source](https://github.com/google/uuid/blob/v1.6.0/dce.go#L17)
    
    
    const (
    	Person = Domain(0)
    	Group  = Domain(1)
    	Org    = Domain(2)
    )

Domain constants for DCE Security (Version 2) UUIDs. 

[View Source](https://github.com/google/uuid/blob/v1.6.0/uuid.go#L29)
    
    
    const (
    	Invalid   = Variant([iota](/builtin#iota)) // Invalid UUID
    	RFC4122                   // The variant specified in RFC4122
    	Reserved                  // Reserved, NCS backward compatibility.
    	Microsoft                 // Reserved, Microsoft Corporation backward compatibility.
    	Future                    // Reserved for future definition.
    )

Constants returned by Variant. 

### Variables ¶

[View Source](https://github.com/google/uuid/blob/v1.6.0/hash.go#L14)
    
    
    var (
    	NameSpaceDNS  = Must(Parse("6ba7b810-9dad-11d1-80b4-00c04fd430c8"))
    	NameSpaceURL  = Must(Parse("6ba7b811-9dad-11d1-80b4-00c04fd430c8"))
    	NameSpaceOID  = Must(Parse("6ba7b812-9dad-11d1-80b4-00c04fd430c8"))
    	NameSpaceX500 = Must(Parse("6ba7b814-9dad-11d1-80b4-00c04fd430c8"))
    	Nil           UUID // empty UUID, all zeros
    
    	// The Max UUID is special form of UUID that is specified to have all 128 bits set to 1.
    	Max = UUID{
    		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
    		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
    	}
    )

Well known namespace IDs and UUIDs 

### Functions ¶

####  func [ClockSequence](https://github.com/google/uuid/blob/v1.6.0/time.go#L76) ¶
    
    
    func ClockSequence() [int](/builtin#int)

ClockSequence returns the current clock sequence, generating one if not already set. The clock sequence is only used for Version 1 UUIDs. 

The uuid package does not use global static storage for the clock sequence or the last time a UUID was generated. Unless SetClockSequence is used, a new random clock sequence is generated the first time a clock sequence is requested by ClockSequence, GetTime, or NewUUID. (section 4.2.1.1) 

####  func [DisableRandPool](https://github.com/google/uuid/blob/v1.6.0/uuid.go#L348) ¶ added in v1.3.0
    
    
    func DisableRandPool()

DisableRandPool disables the randomness pool if it was previously enabled with EnableRandPool. 

Both EnableRandPool and DisableRandPool are not thread-safe and should only be called when there is no possibility that New or any other UUID Version 4 generation function will be called concurrently. 

####  func [EnableRandPool](https://github.com/google/uuid/blob/v1.6.0/uuid.go#L338) ¶ added in v1.3.0
    
    
    func EnableRandPool()

EnableRandPool enables internal randomness pool used for Random (Version 4) UUID generation. The pool contains random bytes read from the random number generator on demand in batches. Enabling the pool may improve the UUID generation throughput significantly. 

Since the pool is stored on the Go heap, this feature may be a bad fit for security sensitive applications. 

Both EnableRandPool and DisableRandPool are not thread-safe and should only be called when there is no possibility that New or any other UUID Version 4 generation function will be called concurrently. 

####  func [IsInvalidLengthError](https://github.com/google/uuid/blob/v1.6.0/uuid.go#L54) ¶ added in v1.3.0
    
    
    func IsInvalidLengthError(err [error](/builtin#error)) [bool](/builtin#bool)

IsInvalidLengthError is matcher function for custom error invalidLengthError 

####  func [NewString](https://github.com/google/uuid/blob/v1.6.0/version4.go#L21) ¶ added in v1.2.0
    
    
    func NewString() [string](/builtin#string)

NewString creates a new random UUID and returns it as a string or panics. NewString is equivalent to the expression 
    
    
    uuid.New().String()
    

####  func [NodeID](https://github.com/google/uuid/blob/v1.6.0/node.go#L60) ¶
    
    
    func NodeID() [][byte](/builtin#byte)

NodeID returns a slice of a copy of the current Node ID, setting the Node ID if not already set. 

####  func [NodeInterface](https://github.com/google/uuid/blob/v1.6.0/node.go#L21) ¶
    
    
    func NodeInterface() [string](/builtin#string)

NodeInterface returns the name of the interface from which the NodeID was derived. The interface "user" is returned if the NodeID was set by SetNodeID. 

####  func [SetClockSequence](https://github.com/google/uuid/blob/v1.6.0/time.go#L91) ¶
    
    
    func SetClockSequence(seq [int](/builtin#int))

SetClockSequence sets the clock sequence to the lower 14 bits of seq. Setting to -1 causes a new sequence to be generated. 

####  func [SetNodeID](https://github.com/google/uuid/blob/v1.6.0/node.go#L73) ¶
    
    
    func SetNodeID(id [][byte](/builtin#byte)) [bool](/builtin#bool)

SetNodeID sets the Node ID to be used for Version 1 UUIDs. The first 6 bytes of id are used. If id is less than 6 bytes then false is returned and the Node ID is not set. 

####  func [SetNodeInterface](https://github.com/google/uuid/blob/v1.6.0/node.go#L33) ¶
    
    
    func SetNodeInterface(name [string](/builtin#string)) [bool](/builtin#bool)

SetNodeInterface selects the hardware address to be used for Version 1 UUIDs. If name is "" then the first usable interface found will be used or a random Node ID will be generated. If a named interface cannot be found then false is returned. 

SetNodeInterface never fails when name is "". 

####  func [SetRand](https://github.com/google/uuid/blob/v1.6.0/uuid.go#L319) ¶
    
    
    func SetRand(r [io](/io).[Reader](/io#Reader))

SetRand sets the random number generator to r, which implements io.Reader. If r.Read returns an error when the package requests random data then a panic will be issued. 

Calling SetRand with nil sets the random number generator to the default generator. 

####  func [Validate](https://github.com/google/uuid/blob/v1.6.0/uuid.go#L195) ¶ added in v1.5.0
    
    
    func Validate(s [string](/builtin#string)) [error](/builtin#error)

Validate returns an error if s is not a properly formatted UUID in one of the following formats: 
    
    
    xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
    urn:uuid:xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
    xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
    {xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx}
    

It returns an error if the format is invalid, otherwise nil. 

### Types ¶

####  type [Domain](https://github.com/google/uuid/blob/v1.6.0/dce.go#L14) ¶
    
    
    type Domain [byte](/builtin#byte)

A Domain represents a Version 2 domain 

####  func (Domain) [String](https://github.com/google/uuid/blob/v1.6.0/dce.go#L70) ¶
    
    
    func (d Domain) String() [string](/builtin#string)

####  type [NullUUID](https://github.com/google/uuid/blob/v1.6.0/null.go#L29) ¶ added in v1.3.0
    
    
    type NullUUID struct {
    	UUID  UUID
    	Valid [bool](/builtin#bool) // Valid is true if UUID is not NULL
    }

NullUUID represents a UUID that may be null. NullUUID implements the SQL driver.Scanner interface so it can be used as a scan destination: 
    
    
    var u uuid.NullUUID
    err := db.QueryRow("SELECT name FROM foo WHERE id=?", id).Scan(&u)
    ...
    if u.Valid {
       // use u.UUID
    } else {
       // NULL value
    }
    

####  func (NullUUID) [MarshalBinary](https://github.com/google/uuid/blob/v1.6.0/null.go#L61) ¶ added in v1.3.0
    
    
    func (nu NullUUID) MarshalBinary() ([][byte](/builtin#byte), [error](/builtin#error))

MarshalBinary implements encoding.BinaryMarshaler. 

####  func (NullUUID) [MarshalJSON](https://github.com/google/uuid/blob/v1.6.0/null.go#L101) ¶ added in v1.3.0
    
    
    func (nu NullUUID) MarshalJSON() ([][byte](/builtin#byte), [error](/builtin#error))

MarshalJSON implements json.Marshaler. 

####  func (NullUUID) [MarshalText](https://github.com/google/uuid/blob/v1.6.0/null.go#L80) ¶ added in v1.3.0
    
    
    func (nu NullUUID) MarshalText() ([][byte](/builtin#byte), [error](/builtin#error))

MarshalText implements encoding.TextMarshaler. 

####  func (*NullUUID) [Scan](https://github.com/google/uuid/blob/v1.6.0/null.go#L35) ¶ added in v1.3.0
    
    
    func (nu *NullUUID) Scan(value interface{}) [error](/builtin#error)

Scan implements the SQL driver.Scanner interface. 

####  func (*NullUUID) [UnmarshalBinary](https://github.com/google/uuid/blob/v1.6.0/null.go#L70) ¶ added in v1.3.0
    
    
    func (nu *NullUUID) UnmarshalBinary(data [][byte](/builtin#byte)) [error](/builtin#error)

UnmarshalBinary implements encoding.BinaryUnmarshaler. 

####  func (*NullUUID) [UnmarshalJSON](https://github.com/google/uuid/blob/v1.6.0/null.go#L110) ¶ added in v1.3.0
    
    
    func (nu *NullUUID) UnmarshalJSON(data [][byte](/builtin#byte)) [error](/builtin#error)

UnmarshalJSON implements json.Unmarshaler. 

####  func (*NullUUID) [UnmarshalText](https://github.com/google/uuid/blob/v1.6.0/null.go#L89) ¶ added in v1.3.0
    
    
    func (nu *NullUUID) UnmarshalText(data [][byte](/builtin#byte)) [error](/builtin#error)

UnmarshalText implements encoding.TextUnmarshaler. 

####  func (NullUUID) [Value](https://github.com/google/uuid/blob/v1.6.0/null.go#L52) ¶ added in v1.3.0
    
    
    func (nu NullUUID) Value() ([driver](/database/sql/driver).[Value](/database/sql/driver#Value), [error](/builtin#error))

Value implements the driver Valuer interface. 

####  type [Time](https://github.com/google/uuid/blob/v1.6.0/time.go#L15) ¶
    
    
    type Time [int64](/builtin#int64)

A Time represents a time as the number of 100's of nanoseconds since 15 Oct 1582\. 

####  func [GetTime](https://github.com/google/uuid/blob/v1.6.0/time.go#L45) ¶
    
    
    func GetTime() (Time, [uint16](/builtin#uint16), [error](/builtin#error))

GetTime returns the current Time (100s of nanoseconds since 15 Oct 1582) and clock sequence as well as adjusting the clock sequence as needed. An error is returned if the current time cannot be determined. 

####  func (Time) [UnixTime](https://github.com/google/uuid/blob/v1.6.0/time.go#L35) ¶
    
    
    func (t Time) UnixTime() (sec, nsec [int64](/builtin#int64))

UnixTime converts t the number of seconds and nanoseconds using the Unix epoch of 1 Jan 1970. 

####  type [UUID](https://github.com/google/uuid/blob/v1.6.0/uuid.go#L20) ¶
    
    
    type UUID [16][byte](/builtin#byte)

A UUID is a 128 bit (16 byte) Universal Unique IDentifier as defined in [RFC 4122](https://rfc-editor.org/rfc/rfc4122.html). 

####  func [FromBytes](https://github.com/google/uuid/blob/v1.6.0/uuid.go#L176) ¶
    
    
    func FromBytes(b [][byte](/builtin#byte)) (uuid UUID, err [error](/builtin#error))

FromBytes creates a new UUID from a byte slice. Returns an error if the slice does not have a length of 16. The bytes are copied from the slice. 

####  func [Must](https://github.com/google/uuid/blob/v1.6.0/uuid.go#L182) ¶
    
    
    func Must(uuid UUID, err [error](/builtin#error)) UUID

Must returns uuid if err is nil and panics otherwise. 

####  func [MustParse](https://github.com/google/uuid/blob/v1.6.0/uuid.go#L166) ¶ added in v1.1.0
    
    
    func MustParse(s [string](/builtin#string)) UUID

MustParse is like Parse but panics if the string cannot be parsed. It simplifies safe initialization of global variables holding compiled UUIDs. 

####  func [New](https://github.com/google/uuid/blob/v1.6.0/version4.go#L13) ¶
    
    
    func New() UUID

New creates a new random UUID or panics. New is equivalent to the expression 
    
    
    uuid.Must(uuid.NewRandom())
    

####  func [NewDCEGroup](https://github.com/google/uuid/blob/v1.6.0/dce.go#L54) ¶
    
    
    func NewDCEGroup() (UUID, [error](/builtin#error))

NewDCEGroup returns a DCE Security (Version 2) UUID in the group domain with the id returned by os.Getgid. 
    
    
    NewDCESecurity(Group, uint32(os.Getgid()))
    

####  func [NewDCEPerson](https://github.com/google/uuid/blob/v1.6.0/dce.go#L46) ¶
    
    
    func NewDCEPerson() (UUID, [error](/builtin#error))

NewDCEPerson returns a DCE Security (Version 2) UUID in the person domain with the id returned by os.Getuid. 
    
    
    NewDCESecurity(Person, uint32(os.Getuid()))
    

####  func [NewDCESecurity](https://github.com/google/uuid/blob/v1.6.0/dce.go#L32) ¶
    
    
    func NewDCESecurity(domain Domain, id [uint32](/builtin#uint32)) (UUID, [error](/builtin#error))

NewDCESecurity returns a DCE Security (Version 2) UUID. 

The domain should be one of Person, Group or Org. On a POSIX system the id should be the users UID for the Person domain and the users GID for the Group. The meaning of id for the domain Org or on non-POSIX systems is site defined. 

For a given domain/id pair the same token may be returned for up to 7 minutes and 10 seconds. 

####  func [NewHash](https://github.com/google/uuid/blob/v1.6.0/hash.go#L33) ¶
    
    
    func NewHash(h [hash](/hash).[Hash](/hash#Hash), space UUID, data [][byte](/builtin#byte), version [int](/builtin#int)) UUID

NewHash returns a new UUID derived from the hash of space concatenated with data generated by h. The hash should be at least 16 byte in length. The first 16 bytes of the hash are used to form the UUID. The version of the UUID will be the lower 4 bits of version. NewHash is used to implement NewMD5 and NewSHA1. 

####  func [NewMD5](https://github.com/google/uuid/blob/v1.6.0/hash.go#L49) ¶
    
    
    func NewMD5(space UUID, data [][byte](/builtin#byte)) UUID

NewMD5 returns a new MD5 (Version 3) UUID based on the supplied name space and data. It is the same as calling: 
    
    
    NewHash(md5.New(), space, data, 3)
    

####  func [NewRandom](https://github.com/google/uuid/blob/v1.6.0/version4.go#L39) ¶
    
    
    func NewRandom() (UUID, [error](/builtin#error))

NewRandom returns a Random (Version 4) UUID. 

The strength of the UUIDs is based on the strength of the crypto/rand package. 

Uses the randomness pool if it was enabled with EnableRandPool. 

A note about uniqueness derived from the UUID Wikipedia entry: 
    
    
    Randomly generated UUIDs have 122 random bits.  One's annual risk of being
    hit by a meteorite is estimated to be one chance in 17 billion, that
    means the probability is about 0.00000000006 (6 × 10−11),
    equivalent to the odds of creating a few tens of trillions of UUIDs in a
    year and having one duplicate.
    

####  func [NewRandomFromReader](https://github.com/google/uuid/blob/v1.6.0/version4.go#L47) ¶ added in v1.1.2
    
    
    func NewRandomFromReader(r [io](/io).[Reader](/io#Reader)) (UUID, [error](/builtin#error))

NewRandomFromReader returns a UUID based on bytes read from a given io.Reader. 

####  func [NewSHA1](https://github.com/google/uuid/blob/v1.6.0/hash.go#L57) ¶
    
    
    func NewSHA1(space UUID, data [][byte](/builtin#byte)) UUID

NewSHA1 returns a new SHA1 (Version 5) UUID based on the supplied name space and data. It is the same as calling: 
    
    
    NewHash(sha1.New(), space, data, 5)
    

####  func [NewUUID](https://github.com/google/uuid/blob/v1.6.0/version1.go#L19) ¶
    
    
    func NewUUID() (UUID, [error](/builtin#error))

NewUUID returns a Version 1 UUID based on the current NodeID and clock sequence, and the current time. If the NodeID has not been set by SetNodeID or SetNodeInterface then it will be set automatically. If the NodeID cannot be set NewUUID returns nil. If clock sequence has not been set by SetClockSequence then it will be set automatically. If GetTime fails to return the current NewUUID returns nil and an error. 

In most cases, New should be used. 

####  func [NewV6](https://github.com/google/uuid/blob/v1.6.0/version6.go#L21) ¶ added in v1.5.0
    
    
    func NewV6() (UUID, [error](/builtin#error))

UUID version 6 is a field-compatible version of UUIDv1, reordered for improved DB locality. It is expected that UUIDv6 will primarily be used in contexts where there are existing v1 UUIDs. Systems that do not involve legacy UUIDv1 SHOULD consider using UUIDv7 instead. 

see <https://datatracker.ietf.org/doc/html/draft-peabody-dispatch-new-uuid-format-03#uuidv6>

NewV6 returns a Version 6 UUID based on the current NodeID and clock sequence, and the current time. If the NodeID has not been set by SetNodeID or SetNodeInterface then it will be set automatically. If the NodeID cannot be set NewV6 set NodeID is random bits automatically . If clock sequence has not been set by SetClockSequence then it will be set automatically. If GetTime fails to return the current NewV6 returns Nil and an error. 

####  func [NewV7](https://github.com/google/uuid/blob/v1.6.0/version7.go#L23) ¶ added in v1.5.0
    
    
    func NewV7() (UUID, [error](/builtin#error))

UUID version 7 features a time-ordered value field derived from the widely implemented and well known Unix Epoch timestamp source, the number of milliseconds seconds since midnight 1 Jan 1970 UTC, leap seconds excluded. As well as improved entropy characteristics over versions 1 or 6. 

see <https://datatracker.ietf.org/doc/html/draft-peabody-dispatch-new-uuid-format-03#name-uuid-version-7>

Implementations SHOULD utilize UUID version 7 over UUID version 1 and 6 if possible. 

NewV7 returns a Version 7 UUID based on the current time(Unix Epoch). Uses the randomness pool if it was enabled with EnableRandPool. On error, NewV7 returns Nil and an error 

####  func [NewV7FromReader](https://github.com/google/uuid/blob/v1.6.0/version7.go#L35) ¶ added in v1.5.0
    
    
    func NewV7FromReader(r [io](/io).[Reader](/io#Reader)) (UUID, [error](/builtin#error))

NewV7FromReader returns a Version 7 UUID based on the current time(Unix Epoch). it use NewRandomFromReader fill random bits. On error, NewV7FromReader returns Nil and an error. 

####  func [Parse](https://github.com/google/uuid/blob/v1.6.0/uuid.go#L68) ¶
    
    
    func Parse(s [string](/builtin#string)) (UUID, [error](/builtin#error))

Parse decodes s into a UUID or returns an error if it cannot be parsed. Both the standard UUID forms defined in [RFC 4122](https://rfc-editor.org/rfc/rfc4122.html) (xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx and urn:uuid:xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx) are decoded. In addition, Parse accepts non-standard strings such as the raw hex encoding xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx and 38 byte "Microsoft style" encodings, e.g. {xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx}. Only the middle 36 bytes are examined in the latter case. Parse should not be used to validate strings as it parses non-standard encodings as indicated above. 

####  func [ParseBytes](https://github.com/google/uuid/blob/v1.6.0/uuid.go#L120) ¶
    
    
    func ParseBytes(b [][byte](/builtin#byte)) (UUID, [error](/builtin#error))

ParseBytes is like Parse, except it parses a byte slice instead of a string. 

####  func (UUID) [ClockSequence](https://github.com/google/uuid/blob/v1.6.0/time.go#L132) ¶
    
    
    func (uuid UUID) ClockSequence() [int](/builtin#int)

ClockSequence returns the clock sequence encoded in uuid. The clock sequence is only well defined for version 1 and 2 UUIDs. 

####  func (UUID) [Domain](https://github.com/google/uuid/blob/v1.6.0/dce.go#L60) ¶
    
    
    func (uuid UUID) Domain() Domain

Domain returns the domain for a Version 2 UUID. Domains are only defined for Version 2 UUIDs. 

####  func (UUID) [ID](https://github.com/google/uuid/blob/v1.6.0/dce.go#L66) ¶
    
    
    func (uuid UUID) ID() [uint32](/builtin#uint32)

ID returns the id for a Version 2 UUID. IDs are only defined for Version 2 UUIDs. 

####  func (UUID) [MarshalBinary](https://github.com/google/uuid/blob/v1.6.0/marshal.go#L27) ¶
    
    
    func (uuid UUID) MarshalBinary() ([][byte](/builtin#byte), [error](/builtin#error))

MarshalBinary implements encoding.BinaryMarshaler. 

####  func (UUID) [MarshalText](https://github.com/google/uuid/blob/v1.6.0/marshal.go#L10) ¶
    
    
    func (uuid UUID) MarshalText() ([][byte](/builtin#byte), [error](/builtin#error))

MarshalText implements encoding.TextMarshaler. 

####  func (UUID) [NodeID](https://github.com/google/uuid/blob/v1.6.0/node.go#L86) ¶
    
    
    func (uuid UUID) NodeID() [][byte](/builtin#byte)

NodeID returns the 6 byte node id encoded in uuid. It returns nil if uuid is not valid. The NodeID is only well defined for version 1 and 2 UUIDs. 

####  func (*UUID) [Scan](https://github.com/google/uuid/blob/v1.6.0/sql.go#L15) ¶
    
    
    func (uuid *UUID) Scan(src interface{}) [error](/builtin#error)

Scan implements sql.Scanner so UUIDs can be read from databases transparently. Currently, database types that map to string and []byte are supported. Please consult database-specific driver documentation for matching types. 

####  func (UUID) [String](https://github.com/google/uuid/blob/v1.6.0/uuid.go#L244) ¶
    
    
    func (uuid UUID) String() [string](/builtin#string)

String returns the string form of uuid, xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx , or "" if uuid is invalid. 

####  func (UUID) [Time](https://github.com/google/uuid/blob/v1.6.0/time.go#L112) ¶
    
    
    func (uuid UUID) Time() Time

Time returns the time in 100s of nanoseconds since 15 Oct 1582 encoded in uuid. The time is only defined for version 1, 2, 6 and 7 UUIDs. 

####  func (UUID) [URN](https://github.com/google/uuid/blob/v1.6.0/uuid.go#L252) ¶
    
    
    func (uuid UUID) URN() [string](/builtin#string)

URN returns the [RFC 2141](https://rfc-editor.org/rfc/rfc2141.html) URN form of uuid, urn:uuid:xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx, or "" if uuid is invalid. 

####  func (*UUID) [UnmarshalBinary](https://github.com/google/uuid/blob/v1.6.0/marshal.go#L32) ¶
    
    
    func (uuid *UUID) UnmarshalBinary(data [][byte](/builtin#byte)) [error](/builtin#error)

UnmarshalBinary implements encoding.BinaryUnmarshaler. 

####  func (*UUID) [UnmarshalText](https://github.com/google/uuid/blob/v1.6.0/marshal.go#L17) ¶
    
    
    func (uuid *UUID) UnmarshalText(data [][byte](/builtin#byte)) [error](/builtin#error)

UnmarshalText implements encoding.TextUnmarshaler. 

####  func (UUID) [Value](https://github.com/google/uuid/blob/v1.6.0/sql.go#L57) ¶
    
    
    func (uuid UUID) Value() ([driver](/database/sql/driver).[Value](/database/sql/driver#Value), [error](/builtin#error))

Value implements sql.Valuer so that UUIDs can be written to databases transparently. Currently, UUIDs map to strings. Please consult database-specific driver documentation for matching types. 

####  func (UUID) [Variant](https://github.com/google/uuid/blob/v1.6.0/uuid.go#L272) ¶
    
    
    func (uuid UUID) Variant() Variant

Variant returns the variant encoded in uuid. 

####  func (UUID) [Version](https://github.com/google/uuid/blob/v1.6.0/uuid.go#L286) ¶
    
    
    func (uuid UUID) Version() Version

Version returns the version of uuid. 

####  type [UUIDs](https://github.com/google/uuid/blob/v1.6.0/uuid.go#L356) ¶ added in v1.4.0
    
    
    type UUIDs []UUID

UUIDs is a slice of UUID types. 

####  func (UUIDs) [Strings](https://github.com/google/uuid/blob/v1.6.0/uuid.go#L359) ¶ added in v1.4.0
    
    
    func (uuids UUIDs) Strings() [][string](/builtin#string)

Strings returns a string slice containing the string form of each UUID in uuids. 

####  type [Variant](https://github.com/google/uuid/blob/v1.6.0/uuid.go#L26) ¶
    
    
    type Variant [byte](/builtin#byte)

A Variant represents a UUID's variant. 

####  func (Variant) [String](https://github.com/google/uuid/blob/v1.6.0/uuid.go#L297) ¶
    
    
    func (v Variant) String() [string](/builtin#string)

####  type [Version](https://github.com/google/uuid/blob/v1.6.0/uuid.go#L23) ¶
    
    
    type Version [byte](/builtin#byte)

A Version represents a UUID's version. 

####  func (Version) [String](https://github.com/google/uuid/blob/v1.6.0/uuid.go#L290) ¶
    
    
    func (v Version) String() [string](/builtin#string)
