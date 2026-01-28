# go-playground/validator

> Auto-fetched from [https://pkg.go.dev/github.com/go-playground/validator/v10](https://pkg.go.dev/github.com/go-playground/validator/v10)
> Last Updated: 2026-01-28T21:44:09.803271+00:00

---

Overview
¶
Singleton
Validation Functions Return Type error
Custom Validation Functions
Cross-Field Validation
Multiple Validators
Using Validator Tags
Baked In Validators and Tags
Skip Field
Or Operator
StructOnly
NoStructLevel
Omit Empty
Omit Nil
Omit Zero
Dive
Required
Required If
Required Unless
Required With
Required With All
Required Without
Required Without All
Excluded If
Excluded Unless
Is Default
Length
Maximum
Minimum
Equals
Not Equal
One Of
One Of Case Insensitive
Greater Than
Greater Than or Equal
Less Than
Less Than or Equal
Field Equals Another Field
Field Does Not Equal Another Field
Field Greater Than Another Field
Field Greater Than Another Relative Field
Field Greater Than or Equal To Another Field
Field Greater Than or Equal To Another Relative Field
Less Than Another Field
Less Than Another Relative Field
Less Than or Equal To Another Field
Less Than or Equal To Another Relative Field
Field Contains Another Field
Field Excludes Another Field
Unique
ValidateFn
Alpha Only
Alpha Space
Alphanumeric
Alphanumeric Space
Alpha Unicode
Alphanumeric Unicode
Boolean
Number
Numeric
Hexadecimal String
Hexcolor String
Lowercase String
Uppercase String
RGB String
RGBA String
HSL String
HSLA String
E.164 Phone Number String
E-mail String
JSON String
JWT String
File
Image path
File Path
URL String
URI String
Urn
RFC 2141
String
Base32 String
Base64 String
Base64URL String
Base64RawURL String
Bitcoin Address
Ethereum Address
Contains
Contains Any
Contains Rune
Excludes
Excludes All
Excludes Rune
Starts With
Ends With
Does Not Start With
Does Not End With
International Standard Book Number
International Standard Book Number 10
International Standard Book Number 13
Universally Unique Identifier UUID
Universally Unique Identifier UUID v3
Universally Unique Identifier UUID v4
Universally Unique Identifier UUID v5
Universally Unique Lexicographically Sortable Identifier ULID
ASCII
Printable ASCII
Multi-Byte Characters
Data URL
Latitude
Longitude
Employeer Identification Number EIN
Social Security Number SSN
Internet Protocol Address IP
Internet Protocol Address IPv4
Internet Protocol Address IPv6
Classless Inter-Domain Routing CIDR
Classless Inter-Domain Routing CIDRv4
Classless Inter-Domain Routing CIDRv6
Transmission Control Protocol Address TCP
Transmission Control Protocol Address TCPv4
Transmission Control Protocol Address TCPv6
User Datagram Protocol Address UDP
User Datagram Protocol Address UDPv4
User Datagram Protocol Address UDPv6
Internet Protocol Address IP
Internet Protocol Address IPv4
Internet Protocol Address IPv6
Unix domain socket end point Address
Unix Domain Socket Exists
Media Access Control Address MAC
Hostname
RFC 952
Hostname
RFC 1123
HTML Tags
HTML Encoded
URL Encoded
Directory
Directory Path
HostPort
Port
Datetime
Iso3166-1 alpha-2
Iso3166-1 alpha-3
Iso3166-1 alpha-numeric
BCP 47 Language Tag
RFC 1035
label
TimeZone
Semantic Version
CVE Identifier
Credit Card
Luhn Checksum
MongoDB
Cron
SpiceDb ObjectID/Permission/Object Type
Alias Validators and Tags
Non standard validators
Panics
Package validator implements value validations for structs and individual fields
based on tags.
It can also handle Cross-Field and Cross-Struct validation for nested structs
and has the ability to dive into arrays and maps of any type.
see more examples
https://github.com/go-playground/validator/tree/master/_examples
Singleton
¶
Validator is designed to be thread-safe and used as a singleton instance.
It caches information about your struct and validations,
in essence only parsing your validation tags once per struct type.
Using multiple instances neglects the benefit of caching.
The not thread-safe functions are explicitly marked as such in the documentation.
Validation Functions Return Type error
¶
Doing things this way is actually the way the standard library does, see the
file.Open method here:
https://golang.org/pkg/os/#Open.
The authors return type "error" to avoid the issue discussed in the following,
where err is always != nil:
http://stackoverflow.com/a/29138676/3158232
https://github.com/go-playground/validator/issues/134
Validator only InvalidValidationError for bad validation input, nil or
ValidationErrors as type error; so, in your code all you need to do is check
if the error returned is not nil, and if it's not check if error is
InvalidValidationError ( if necessary, most of the time it isn't ) type cast
it to type ValidationErrors like so err.(validator.ValidationErrors).
Custom Validation Functions
¶
Custom Validation functions can be added. Example:
// Structure
func customFunc(fl validator.FieldLevel) bool {

if fl.Field().String() == "invalid" {
return false
}

return true
}

validate.RegisterValidation("custom tag name", customFunc)
// NOTES: using the same tag name as an existing function
//        will overwrite the existing one
Cross-Field Validation
¶
Cross-Field Validation can be done via the following tags:
eqfield
nefield
gtfield
gtefield
ltfield
ltefield
eqcsfield
necsfield
gtcsfield
gtecsfield
ltcsfield
ltecsfield
If, however, some custom cross-field validation is required, it can be done
using a custom validation.
Why not just have cross-fields validation tags (i.e. only eqcsfield and not
eqfield)?
The reason is efficiency. If you want to check a field within the same struct
"eqfield" only has to find the field on the same struct (1 level). But, if we
used "eqcsfield" it could be multiple levels down. Example:
type Inner struct {
StartDate time.Time
}

type Outer struct {
InnerStructField *Inner
CreatedAt time.Time      `validate:"ltecsfield=InnerStructField.StartDate"`
}

now := time.Now()

inner := &Inner{
StartDate: now,
}

outer := &Outer{
InnerStructField: inner,
CreatedAt: now,
}

errs := validate.Struct(outer)

// NOTE: when calling validate.Struct(val) topStruct will be the top level struct passed
//       into the function
//       when calling validate.VarWithValue(val, field, tag) val will be
//       whatever you pass, struct, field...
//       when calling validate.Field(field, tag) val will be nil
Multiple Validators
¶
Multiple validators on a field will process in the order defined. Example:
type Test struct {
Field `validate:"max=10,min=1"`
}

// max will be checked then min
Bad Validator definitions are not handled by the library. Example:
type Test struct {
Field `validate:"min=10,max=0"`
}

// this definition of min max will never succeed
Using Validator Tags
¶
Baked In Cross-Field validation only compares fields on the same struct.
If Cross-Field + Cross-Struct validation is needed you should implement your
own custom validator.
Comma (",") is the default separator of validation tags. If you wish to
have a comma included within the parameter (i.e. excludesall=,) you will need to
use the UTF-8 hex representation 0x2C, which is replaced in the code as a comma,
so the above will become excludesall=0x2C.
type Test struct {
Field `validate:"excludesall=,"`    // BAD! Do not include a comma.
Field `validate:"excludesall=0x2C"` // GOOD! Use the UTF-8 hex representation.
}
Pipe ("|") is the 'or' validation tags deparator. If you wish to
have a pipe included within the parameter i.e. excludesall=| you will need to
use the UTF-8 hex representation 0x7C, which is replaced in the code as a pipe,
so the above will become excludesall=0x7C
type Test struct {
Field `validate:"excludesall=|"`    // BAD! Do not include a pipe!
Field `validate:"excludesall=0x7C"` // GOOD! Use the UTF-8 hex representation.
}
Baked In Validators and Tags
¶
Here is a list of the current built in validators:
Skip Field
¶
Tells the validation to skip this struct field; this is particularly
handy in ignoring embedded structs from being validated. (Usage: -)
Usage: -
Or Operator
¶
This is the 'or' operator allowing multiple validators to be used and
accepted. (Usage: rgb|rgba) <-- this would allow either rgb or rgba
colors to be accepted. This can also be combined with 'and' for example
( Usage: omitempty,rgb|rgba)
Usage: |
StructOnly
¶
When a field that is a nested struct is encountered, and contains this flag
any validation on the nested struct will be run, but none of the nested
struct fields will be validated. This is useful if inside of your program
you know the struct will be valid, but need to verify it has been assigned.
NOTE: only "required" and "omitempty" can be used on a struct itself.
Usage: structonly
NoStructLevel
¶
Same as structonly tag except that any struct level validations will not run.
Usage: nostructlevel
Omit Empty
¶
Allows conditional validation, for example, if a field is not set with
a value (Determined by the "required" validator) then other validation
such as min or max won't run, but if a value is set validation will run.
Usage: omitempty
Omit Nil
¶
Allows to skip the validation if the value is nil (same as omitempty, but
only for the nil-values).
Usage: omitnil
Omit Zero
¶
Allows to skip the validation if the value is a zero value.
For pointers, it checks if the pointer is nil or the underlying value is a zero value.
For slices and maps, it checks if the value is nil or empty.
Otherwise, behaves the same as omitempty.
Usage: omitzero
Dive
¶
This tells the validator to dive into a slice, array or map and validate that
level of the slice, array or map with the validation tags that follow.
Multidimensional nesting is also supported, each level you wish to dive will
require another dive tag. dive has some sub-tags, 'keys' & 'endkeys', please see
the Keys & EndKeys section just below.
Usage: dive
Example #1
[][]string with validation tag "gt=0,dive,len=1,dive,required"
// gt=0 will be applied to []
// len=1 will be applied to []string
// required will be applied to string
Example #2
[][]string with validation tag "gt=0,dive,dive,required"
// gt=0 will be applied to []
// []string will be spared validation
// required will be applied to string
Keys & EndKeys
These are to be used together directly after the dive tag and tells the validator
that anything between 'keys' and 'endkeys' applies to the keys of a map and not the
values; think of it like the 'dive' tag, but for map keys instead of values.
Multidimensional nesting is also supported, each level you wish to validate will
require another 'keys' and 'endkeys' tag. These tags are only valid for maps.
Usage: dive,keys,othertagvalidation(s),endkeys,valuevalidationtags
Example #1
map[string]string with validation tag "gt=0,dive,keys,eq=1|eq=2,endkeys,required"
// gt=0 will be applied to the map itself
// eq=1|eq=2 will be applied to the map keys
// required will be applied to map values
Example #2
map[[2]string]string with validation tag "gt=0,dive,keys,dive,eq=1|eq=2,endkeys,required"
// gt=0 will be applied to the map itself
// eq=1|eq=2 will be applied to each array element in the map keys
// required will be applied to map values
Required
¶
This validates that the value is not the data types default zero value.
For numbers ensures value is not zero. For strings ensures value is
not "". For booleans ensures value is not false. For slices, maps, pointers, interfaces, channels and functions
ensures the value is not nil. For structs ensures value is not the zero value when using WithRequiredStructEnabled.
Usage: required
Required If
¶
The field under validation must be present and not empty only if all
the other specified fields are equal to the value following the specified
field. For strings ensures value is not "". For slices, maps, pointers,
interfaces, channels and functions ensures the value is not nil. For structs ensures value is not the zero value.
Using the same field name multiple times in the parameters will result in a panic at runtime.
Usage: required_if
Examples:
// require the field if the Field1 is equal to the parameter given:
Usage: required_if=Field1 foobar

// require the field if the Field1 and Field2 is equal to the value respectively:
Usage: required_if=Field1 foo Field2 bar
Required Unless
¶
The field under validation must be present and not empty unless all
the other specified fields are equal to the value following the specified
field. For strings ensures value is not "". For slices, maps, pointers,
interfaces, channels and functions ensures the value is not nil. For structs ensures value is not the zero value.
Usage: required_unless
Examples:
// require the field unless the Field1 is equal to the parameter given:
Usage: required_unless=Field1 foobar

// require the field unless the Field1 and Field2 is equal to the value respectively:
Usage: required_unless=Field1 foo Field2 bar
Required With
¶
The field under validation must be present and not empty only if any
of the other specified fields are present. For strings ensures value is
not "". For slices, maps, pointers, interfaces, channels and functions
ensures the value is not nil. For structs ensures value is not the zero value.
Usage: required_with
Examples:
// require the field if the Field1 is present:
Usage: required_with=Field1

// require the field if the Field1 or Field2 is present:
Usage: required_with=Field1 Field2
Required With All
¶
The field under validation must be present and not empty only if all
of the other specified fields are present. For strings ensures value is
not "". For slices, maps, pointers, interfaces, channels and functions
ensures the value is not nil. For structs ensures value is not the zero value.
Usage: required_with_all
Example:
// require the field if the Field1 and Field2 is present:
Usage: required_with_all=Field1 Field2
Required Without
¶
The field under validation must be present and not empty only when any
of the other specified fields are not present. For strings ensures value is
not "". For slices, maps, pointers, interfaces, channels and functions
ensures the value is not nil. For structs ensures value is not the zero value.
Usage: required_without
Examples:
// require the field if the Field1 is not present:
Usage: required_without=Field1

// require the field if the Field1 or Field2 is not present:
Usage: required_without=Field1 Field2
Required Without All
¶
The field under validation must be present and not empty only when all
of the other specified fields are not present. For strings ensures value is
not "". For slices, maps, pointers, interfaces, channels and functions
ensures the value is not nil. For structs ensures value is not the zero value.
Usage: required_without_all
Example:
// require the field if the Field1 and Field2 is not present:
Usage: required_without_all=Field1 Field2
Excluded If
¶
The field under validation must not be present or not empty only if all
the other specified fields are equal to the value following the specified
field. For strings ensures value is not "". For slices, maps, pointers,
interfaces, channels and functions ensures the value is not nil. For structs ensures value is not the zero value.
Usage: excluded_if
Examples:
// exclude the field if the Field1 is equal to the parameter given:
Usage: excluded_if=Field1 foobar

// exclude the field if the Field1 and Field2 is equal to the value respectively:
Usage: excluded_if=Field1 foo Field2 bar
Excluded Unless
¶
The field under validation must not be present or empty unless all
the other specified fields are equal to the value following the specified
field. For strings ensures value is not "". For slices, maps, pointers,
interfaces, channels and functions ensures the value is not nil. For structs ensures value is not the zero value.
Usage: excluded_unless
Examples:
// exclude the field unless the Field1 is equal to the parameter given:
Usage: excluded_unless=Field1 foobar

// exclude the field unless the Field1 and Field2 is equal to the value respectively:
Usage: excluded_unless=Field1 foo Field2 bar
Is Default
¶
This validates that the value is the default value and is almost the
opposite of required.
Usage: isdefault
Length
¶
For numbers, length will ensure that the value is
equal to the parameter given. For strings, it checks that
the string length is exactly that number of characters. For slices,
arrays, and maps, validates the number of items.
Example #1
Usage: len=10
Example #2 (time.Duration)
For time.Duration, len will ensure that the value is equal to the duration given
in the parameter.
Usage: len=1h30m
Maximum
¶
For numbers, max will ensure that the value is
less than or equal to the parameter given. For strings, it checks
that the string length is at most that number of characters. For
slices, arrays, and maps, validates the number of items.
Example #1
Usage: max=10
Example #2 (time.Duration)
For time.Duration, max will ensure that the value is less than or equal to the
duration given in the parameter.
Usage: max=1h30m
Minimum
¶
For numbers, min will ensure that the value is
greater or equal to the parameter given. For strings, it checks that
the string length is at least that number of characters. For slices,
arrays, and maps, validates the number of items.
Example #1
Usage: min=10
Example #2 (time.Duration)
For time.Duration, min will ensure that the value is greater than or equal to
the duration given in the parameter.
Usage: min=1h30m
Equals
¶
For strings & numbers, eq will ensure that the value is
equal to the parameter given. For slices, arrays, and maps,
validates the number of items.
Example #1
Usage: eq=10
Example #2 (time.Duration)
For time.Duration, eq will ensure that the value is equal to the duration given
in the parameter.
Usage: eq=1h30m
Not Equal
¶
For strings & numbers, ne will ensure that the value is not
equal to the parameter given. For slices, arrays, and maps,
validates the number of items.
Example #1
Usage: ne=10
Example #2 (time.Duration)
For time.Duration, ne will ensure that the value is not equal to the duration
given in the parameter.
Usage: ne=1h30m
One Of
¶
For strings, ints, and uints, oneof will ensure that the value
is one of the values in the parameter.  The parameter should be
a list of values separated by whitespace. Values may be
strings or numbers. To match strings with spaces in them, include
the target string between single quotes. Kind of like an 'enum'.
Usage: oneof=red green
oneof='red green' 'blue yellow'
oneof=5 7 9
One Of Case Insensitive
¶
Works the same as oneof but is case insensitive and therefore only accepts strings.
Usage: oneofci=red green
oneofci='red green' 'blue yellow'
Greater Than
¶
For numbers, this will ensure that the value is greater than the
parameter given. For strings, it checks that the string length
is greater than that number of characters. For slices, arrays
and maps it validates the number of items.
Example #1
Usage: gt=10
Example #2 (time.Time)
For time.Time ensures the time value is greater than time.Now.UTC().
Usage: gt
Example #3 (time.Duration)
For time.Duration, gt will ensure that the value is greater than the duration
given in the parameter.
Usage: gt=1h30m
Greater Than or Equal
¶
Same as 'min' above. Kept both to make terminology with 'len' easier.
Example #1
Usage: gte=10
Example #2 (time.Time)
For time.Time ensures the time value is greater than or equal to time.Now.UTC().
Usage: gte
Example #3 (time.Duration)
For time.Duration, gte will ensure that the value is greater than or equal to
the duration given in the parameter.
Usage: gte=1h30m
Less Than
¶
For numbers, this will ensure that the value is less than the parameter given.
For strings, it checks that the string length is less than that number of
characters. For slices, arrays, and maps it validates the number of items.
Example #1
Usage: lt=10
Example #2 (time.Time)
For time.Time ensures the time value is less than time.Now.UTC().
Usage: lt
Example #3 (time.Duration)
For time.Duration, lt will ensure that the value is less than the duration given
in the parameter.
Usage: lt=1h30m
Less Than or Equal
¶
Same as 'max' above. Kept both to make terminology with 'len' easier.
Example #1
Usage: lte=10
Example #2 (time.Time)
For time.Time ensures the time value is less than or equal to time.Now.UTC().
Usage: lte
Example #3 (time.Duration)
For time.Duration, lte will ensure that the value is less than or equal to the
duration given in the parameter.
Usage: lte=1h30m
Field Equals Another Field
¶
This will validate the field value against another fields value either within
a struct or passed in field.
Example #1:
// Validation on Password field using:
Usage: eqfield=ConfirmPassword
Example #2:
// Validating by field:
validate.VarWithValue(password, confirmpassword, "eqfield")
Field Equals Another Field (relative)
This does the same as eqfield except that it validates the field provided relative
to the top level struct.
Usage: eqcsfield=InnerStructField.Field)
Field Does Not Equal Another Field
¶
This will validate the field value against another fields value either within
a struct or passed in field.
Examples:
// Confirm two colors are not the same:
//
// Validation on Color field:
Usage: nefield=Color2

// Validating by field:
validate.VarWithValue(color1, color2, "nefield")
Field Does Not Equal Another Field (relative)
This does the same as nefield except that it validates the field provided
relative to the top level struct.
Usage: necsfield=InnerStructField.Field
Field Greater Than Another Field
¶
Only valid for Numbers, time.Duration and time.Time types, this will validate
the field value against another fields value either within a struct or passed in
field. usage examples are for validation of a Start and End date:
Example #1:
// Validation on End field using:
validate.Struct Usage(gtfield=Start)
Example #2:
// Validating by field:
validate.VarWithValue(start, end, "gtfield")
Field Greater Than Another Relative Field
¶
This does the same as gtfield except that it validates the field provided
relative to the top level struct.
Usage: gtcsfield=InnerStructField.Field
Field Greater Than or Equal To Another Field
¶
Only valid for Numbers, time.Duration and time.Time types, this will validate
the field value against another fields value either within a struct or passed in
field. usage examples are for validation of a Start and End date:
Example #1:
// Validation on End field using:
validate.Struct Usage(gtefield=Start)
Example #2:
// Validating by field:
validate.VarWithValue(start, end, "gtefield")
Field Greater Than or Equal To Another Relative Field
¶
This does the same as gtefield except that it validates the field provided relative
to the top level struct.
Usage: gtecsfield=InnerStructField.Field
Less Than Another Field
¶
Only valid for Numbers, time.Duration and time.Time types, this will validate
the field value against another fields value either within a struct or passed in
field. usage examples are for validation of a Start and End date:
Example #1:
// Validation on End field using:
validate.Struct Usage(ltfield=Start)
Example #2:
// Validating by field:
validate.VarWithValue(start, end, "ltfield")
Less Than Another Relative Field
¶
This does the same as ltfield except that it validates the field provided relative
to the top level struct.
Usage: ltcsfield=InnerStructField.Field
Less Than or Equal To Another Field
¶
Only valid for Numbers, time.Duration and time.Time types, this will validate
the field value against another fields value either within a struct or passed in
field. usage examples are for validation of a Start and End date:
Example #1:
// Validation on End field using:
validate.Struct Usage(ltefield=Start)
Example #2:
// Validating by field:
validate.VarWithValue(start, end, "ltefield")
Less Than or Equal To Another Relative Field
¶
This does the same as ltefield except that it validates the field provided relative
to the top level struct.
Usage: ltecsfield=InnerStructField.Field
Field Contains Another Field
¶
This does the same as contains except for struct fields. It should only be used
with string types. See the behavior of reflect.Value.String() for behavior on
other types.
Usage: containsfield=InnerStructField.Field
Field Excludes Another Field
¶
This does the same as excludes except for struct fields. It should only be used
with string types. See the behavior of reflect.Value.String() for behavior on
other types.
Usage: excludesfield=InnerStructField.Field
Unique
¶
For arrays & slices, unique will ensure that there are no duplicates.
For maps, unique will ensure that there are no duplicate values.
For slices of struct, unique will ensure that there are no duplicate values
in a field of the struct specified via a parameter.
// For arrays, slices, and maps:
Usage: unique

// For slices of struct:
Usage: unique=field
ValidateFn
¶
This validates that an object responds to a method that can return error or bool.
By default it expects an interface `Validate() error` and check that the method
does not return an error. Other methods can be specified using two signatures:
If the method returns an error, it check if the return value is nil.
If the method returns a boolean, it checks if the value is true.
// to use the default method Validate() error
Usage: validateFn

// to use the custom method IsValid() bool (or error)
Usage: validateFn=IsValid
Alpha Only
¶
This validates that a string value contains ASCII alpha characters only
Usage: alpha
Alpha Space
¶
This validates that a string value contains ASCII alpha characters and spaces only
Usage: alphaspace
Alphanumeric
¶
This validates that a string value contains ASCII alphanumeric characters only
Usage: alphanum
Alphanumeric Space
¶
This validates that a string value contains ASCII alphanumeric characters and spaces only
Usage: alphanumspace
Alpha Unicode
¶
This validates that a string value contains unicode alpha characters only
Usage: alphaunicode
Alphanumeric Unicode
¶
This validates that a string value contains unicode alphanumeric characters only
Usage: alphanumunicode
Boolean
¶
This validates that a string value can successfully be parsed into a boolean with strconv.ParseBool
Usage: boolean
Number
¶
This validates that a string value contains number values only.
For integers or float it returns true.
Usage: number
Numeric
¶
This validates that a string value contains a basic numeric value.
basic excludes exponents etc...
for integers or float it returns true.
Usage: numeric
Hexadecimal String
¶
This validates that a string value contains a valid hexadecimal.
Usage: hexadecimal
Hexcolor String
¶
This validates that a string value contains a valid hex color including
hashtag (#)
Usage: hexcolor
Lowercase String
¶
This validates that a string value contains only lowercase characters. An empty string is not a valid lowercase string.
Usage: lowercase
Uppercase String
¶
This validates that a string value contains only uppercase characters. An empty string is not a valid uppercase string.
Usage: uppercase
RGB String
¶
This validates that a string value contains a valid rgb color
Usage: rgb
RGBA String
¶
This validates that a string value contains a valid rgba color
Usage: rgba
HSL String
¶
This validates that a string value contains a valid hsl color
Usage: hsl
HSLA String
¶
This validates that a string value contains a valid hsla color
Usage: hsla
E.164 Phone Number String
¶
This validates that a string value contains a valid E.164 Phone number
https://en.wikipedia.org/wiki/E.164
(ex. +1123456789)
Usage: e164
E-mail String
¶
This validates that a string value contains a valid email
This may not conform to all possibilities of any rfc standard, but neither
does any email provider accept all possibilities.
Usage: email
JSON String
¶
This validates that a string value is valid JSON
Usage: json
JWT String
¶
This validates that a string value is a valid JWT
Usage: jwt
File
¶
This validates that a string value contains a valid file path and that
the file exists on the machine.
This is done using os.Stat, which is a platform independent function.
Usage: file
Image path
¶
This validates that a string value contains a valid file path and that
the file exists on the machine and is an image.
This is done using os.Stat and github.com/gabriel-vasile/mimetype
Usage: image
File Path
¶
This validates that a string value contains a valid file path but does not
validate the existence of that file.
This is done using os.Stat, which is a platform independent function.
Usage: filepath
URL String
¶
This validates that a string value contains a valid url
This will accept any url the golang request uri accepts but must contain
a schema for example http:// or rtmp://
Usage: url
URI String
¶
This validates that a string value contains a valid uri
This will accept any uri the golang request uri accepts
Usage: uri
Urn
RFC 2141
String
RFC 2141
String" aria-label="Go to Urn
RFC 2141
String">¶
This validates that a string value contains a valid URN
according to the
RFC 2141
spec.
Usage: urn_rfc2141
Base32 String
¶
This validates that a string value contains a valid bas324 value.
Although an empty string is valid base32 this will report an empty string
as an error, if you wish to accept an empty string as valid you can use
this with the omitempty tag.
Usage: base32
Base64 String
¶
This validates that a string value contains a valid base64 value.
Although an empty string is valid base64 this will report an empty string
as an error, if you wish to accept an empty string as valid you can use
this with the omitempty tag.
Usage: base64
Base64URL String
¶
This validates that a string value contains a valid base64 URL safe value
according the RFC4648 spec.
Although an empty string is a valid base64 URL safe value, this will report
an empty string as an error, if you wish to accept an empty string as valid
you can use this with the omitempty tag.
Usage: base64url
Base64RawURL String
¶
This validates that a string value contains a valid base64 URL safe value,
but without = padding, according the RFC4648 spec, section 3.2.
Although an empty string is a valid base64 URL safe value, this will report
an empty string as an error, if you wish to accept an empty string as valid
you can use this with the omitempty tag.
Usage: base64rawurl
Bitcoin Address
¶
This validates that a string value contains a valid bitcoin address.
The format of the string is checked to ensure it matches one of the three formats
P2PKH, P2SH and performs checksum validation.
Usage: btc_addr
Bitcoin Bech32 Address (segwit)
This validates that a string value contains a valid bitcoin Bech32 address as defined
by bip-0173 (
https://github.com/bitcoin/bips/blob/master/bip-0173.mediawiki
)
Special thanks to Pieter Wuille for providing reference implementations.
Usage: btc_addr_bech32
Ethereum Address
¶
This validates that a string value contains a valid ethereum address.
The format of the string is checked to ensure it matches the standard Ethereum address format.
Usage: eth_addr
Contains
¶
This validates that a string value contains the substring value.
Usage: contains=@
Contains Any
¶
This validates that a string value contains any Unicode code points
in the substring value.
Usage: containsany=!@#?
Contains Rune
¶
This validates that a string value contains the supplied rune value.
Usage: containsrune=@
Excludes
¶
This validates that a string value does not contain the substring value.
Usage: excludes=@
Excludes All
¶
This validates that a string value does not contain any Unicode code
points in the substring value.
Usage: excludesall=!@#?
Excludes Rune
¶
This validates that a string value does not contain the supplied rune value.
Usage: excludesrune=@
Starts With
¶
This validates that a string value starts with the supplied string value
Usage: startswith=hello
Ends With
¶
This validates that a string value ends with the supplied string value
Usage: endswith=goodbye
Does Not Start With
¶
This validates that a string value does not start with the supplied string value
Usage: startsnotwith=hello
Does Not End With
¶
This validates that a string value does not end with the supplied string value
Usage: endsnotwith=goodbye
International Standard Book Number
¶
This validates that a string value contains a valid isbn10 or isbn13 value.
Usage: isbn
International Standard Book Number 10
¶
This validates that a string value contains a valid isbn10 value.
Usage: isbn10
International Standard Book Number 13
¶
This validates that a string value contains a valid isbn13 value.
Usage: isbn13
Universally Unique Identifier UUID
¶
This validates that a string value contains a valid UUID. Uppercase UUID values will not pass - use `uuid_rfc4122` instead.
Usage: uuid
Universally Unique Identifier UUID v3
¶
This validates that a string value contains a valid version 3 UUID.  Uppercase UUID values will not pass - use `uuid3_rfc4122` instead.
Usage: uuid3
Universally Unique Identifier UUID v4
¶
This validates that a string value contains a valid version 4 UUID.  Uppercase UUID values will not pass - use `uuid4_rfc4122` instead.
Usage: uuid4
Universally Unique Identifier UUID v5
¶
This validates that a string value contains a valid version 5 UUID.  Uppercase UUID values will not pass - use `uuid5_rfc4122` instead.
Usage: uuid5
Universally Unique Lexicographically Sortable Identifier ULID
¶
This validates that a string value contains a valid ULID value.
Usage: ulid
ASCII
¶
This validates that a string value contains only ASCII characters.
NOTE: if the string is blank, this validates as true.
Usage: ascii
Printable ASCII
¶
This validates that a string value contains only printable ASCII characters.
NOTE: if the string is blank, this validates as true.
Usage: printascii
Multi-Byte Characters
¶
This validates that a string value contains one or more multibyte characters.
NOTE: if the string is blank, this validates as true.
Usage: multibyte
Data URL
¶
This validates that a string value contains a valid DataURI.
NOTE: this will also validate that the data portion is valid base64
Usage: datauri
Latitude
¶
This validates that a string value contains a valid latitude.
Usage: latitude
Longitude
¶
This validates that a string value contains a valid longitude.
Usage: longitude
Employeer Identification Number EIN
¶
This validates that a string value contains a valid U.S. Employer Identification Number.
Usage: ein
Social Security Number SSN
¶
This validates that a string value contains a valid U.S. Social Security Number.
Usage: ssn
Internet Protocol Address IP
¶
This validates that a string value contains a valid IP Address.
Usage: ip
Internet Protocol Address IPv4
¶
This validates that a string value contains a valid v4 IP Address.
Usage: ipv4
Internet Protocol Address IPv6
¶
This validates that a string value contains a valid v6 IP Address.
Usage: ipv6
Classless Inter-Domain Routing CIDR
¶
This validates that a string value contains a valid CIDR Address.
Usage: cidr
Classless Inter-Domain Routing CIDRv4
¶
This validates that a string value contains a valid v4 CIDR Address.
Usage: cidrv4
Classless Inter-Domain Routing CIDRv6
¶
This validates that a string value contains a valid v6 CIDR Address.
Usage: cidrv6
Transmission Control Protocol Address TCP
¶
This validates that a string value contains a valid resolvable TCP Address.
Usage: tcp_addr
Transmission Control Protocol Address TCPv4
¶
This validates that a string value contains a valid resolvable v4 TCP Address.
Usage: tcp4_addr
Transmission Control Protocol Address TCPv6
¶
This validates that a string value contains a valid resolvable v6 TCP Address.
Usage: tcp6_addr
User Datagram Protocol Address UDP
¶
This validates that a string value contains a valid resolvable UDP Address.
Usage: udp_addr
User Datagram Protocol Address UDPv4
¶
This validates that a string value contains a valid resolvable v4 UDP Address.
Usage: udp4_addr
User Datagram Protocol Address UDPv6
¶
This validates that a string value contains a valid resolvable v6 UDP Address.
Usage: udp6_addr
Internet Protocol Address IP
¶
This validates that a string value contains a valid resolvable IP Address.
Usage: ip_addr
Internet Protocol Address IPv4
¶
This validates that a string value contains a valid resolvable v4 IP Address.
Usage: ip4_addr
Internet Protocol Address IPv6
¶
This validates that a string value contains a valid resolvable v6 IP Address.
Usage: ip6_addr
Unix domain socket end point Address
¶
This validates that a string value contains a valid Unix Address.
Usage: unix_addr
Unix Domain Socket Exists
¶
This validates that a Unix domain socket file exists at the specified path.
It checks both filesystem-based sockets and Linux abstract sockets (prefixed with @).
For filesystem sockets, it verifies the path exists and is a socket file.
For abstract sockets on Linux, it checks /proc/net/unix.
Usage: uds_exists
Media Access Control Address MAC
¶
This validates that a string value contains a valid MAC Address.
Usage: mac
Note: See Go's ParseMAC for accepted formats and types:
http://golang.org/src/net/mac.go?s=866:918#L29
Hostname
RFC 952
RFC 952
" aria-label="Go to Hostname
RFC 952
">¶
This validates that a string value is a valid Hostname according to
RFC 952
https://tools.ietf.org/html/rfc952
Usage: hostname
Hostname
RFC 1123
RFC 1123
" aria-label="Go to Hostname
RFC 1123
">¶
This validates that a string value is a valid Hostname according to
RFC 1123
https://tools.ietf.org/html/rfc1123
Usage: hostname_rfc1123 or if you want to continue to use 'hostname' in your tags, create an alias.
Full Qualified Domain Name (FQDN)
This validates that a string value contains a valid FQDN.
Usage: fqdn
HTML Tags
¶
This validates that a string value appears to be an HTML element tag
including those described at
https://developer.mozilla.org/en-US/docs/Web/HTML/Element
Usage: html
HTML Encoded
¶
This validates that a string value is a proper character reference in decimal
or hexadecimal format
Usage: html_encoded
URL Encoded
¶
This validates that a string value is percent-encoded (URL encoded) according
to
https://tools.ietf.org/html/rfc3986#section-2.1
Usage: url_encoded
Directory
¶
This validates that a string value contains a valid directory and that
it exists on the machine.
This is done using os.Stat, which is a platform independent function.
Usage: dir
Directory Path
¶
This validates that a string value contains a valid directory but does
not validate the existence of that directory.
This is done using os.Stat, which is a platform independent function.
It is safest to suffix the string with os.PathSeparator if the directory
may not exist at the time of validation.
Usage: dirpath
HostPort
¶
This validates that a string value contains a valid DNS hostname and port that
can be used to validate fields typically passed to sockets and connections.
Usage: hostname_port
Port
¶
This validates that the value falls within the valid port number range of 1 to 65,535.
Usage: port
Datetime
¶
This validates that a string value is a valid datetime based on the supplied datetime format.
Supplied format must match the official Go time format layout as documented in
https://golang.org/pkg/time/
Usage: datetime=2006-01-02
Iso3166-1 alpha-2
¶
This validates that a string value is a valid country code based on iso3166-1 alpha-2 standard.
see:
https://www.iso.org/iso-3166-country-codes.html
Usage: iso3166_1_alpha2
Iso3166-1 alpha-3
¶
This validates that a string value is a valid country code based on iso3166-1 alpha-3 standard.
see:
https://www.iso.org/iso-3166-country-codes.html
Usage: iso3166_1_alpha3
Iso3166-1 alpha-numeric
¶
This validates that a string value is a valid country code based on iso3166-1 alpha-numeric standard.
see:
https://www.iso.org/iso-3166-country-codes.html
Usage: iso3166_1_alpha3
BCP 47 Language Tag
¶
This validates that a string value is a valid BCP 47 language tag, as parsed by language.Parse.
More information on
https://pkg.go.dev/golang.org/x/text/language
Usage: bcp47_language_tag
BIC (SWIFT code - 2022 standard)
This validates that a string value is a valid Business Identifier Code (SWIFT code), defined in ISO 9362:2022.
More information on
https://www.iso.org/standard/84108.html
Usage: bic
BIC (SWIFT code - 2014 standard)
This validates that a string value is a valid Business Identifier Code (SWIFT code), defined in ISO 9362:2014.
More information on
https://www.iso.org/standard/60390.html
Usage: bic_iso_9362_2014
RFC 1035
label
RFC 1035
label" aria-label="Go to
RFC 1035
label">¶
This validates that a string value is a valid dns
RFC 1035
label, defined in
RFC 1035
.
More information on
https://datatracker.ietf.org/doc/html/rfc1035
Usage: dns_rfc1035_label
TimeZone
¶
This validates that a string value is a valid time zone based on the time zone database present on the system.
Although empty value and Local value are allowed by time.LoadLocation golang function, they are not allowed by this validator.
More information on
https://golang.org/pkg/time/#LoadLocation
Usage: timezone
Semantic Version
¶
This validates that a string value is a valid semver version, defined in Semantic Versioning 2.0.0.
More information on
https://semver.org/
Usage: semver
CVE Identifier
¶
This validates that a string value is a valid cve id, defined in cve mitre.
More information on
https://cve.mitre.org/
Usage: cve
Credit Card
¶
This validates that a string value contains a valid credit card number using Luhn algorithm.
Usage: credit_card
Luhn Checksum
¶
Usage: luhn_checksum
This validates that a string or (u)int value contains a valid checksum using the Luhn algorithm.
MongoDB
¶
This validates that a string is a valid 24 character hexadecimal string or valid connection string.
Usage: mongodb
mongodb_connection_string
Example:
type Test struct {
ObjectIdField         string `validate:"mongodb"`
ConnectionStringField string `validate:"mongodb_connection_string"`
}
Cron
¶
This validates that a string value contains a valid cron expression.
Usage: cron
SpiceDb ObjectID/Permission/Object Type
¶
This validates that a string is valid for use with SpiceDb for the indicated purpose. If no purpose is given, a purpose of 'id' is assumed.
Usage: spicedb=id|permission|type
Alias Validators and Tags
¶
Alias Validators and Tags
NOTE: When returning an error, the tag returned in "FieldError" will be
the alias tag unless the dive tag is part of the alias. Everything after the
dive tag is not reported as the alias tag. Also, the "ActualTag" in the before
case will be the actual tag within the alias that failed.
Here is a list of the current built in alias tags:
"iscolor"
alias is "hexcolor|rgb|rgba|hsl|hsla" (Usage: iscolor)
"country_code"
alias is "iso3166_1_alpha2|iso3166_1_alpha3|iso3166_1_alpha_numeric" (Usage: country_code)
Validator notes:
regex
a regex validator won't be added because commas and = signs can be part
of a regex which conflict with the validation definitions. Although
workarounds can be made, they take away from using pure regex's.
Furthermore it's quick and dirty but the regex's become harder to
maintain and are not reusable, so it's as much a programming philosophy
as anything.

In place of this new validator functions should be created; a regex can
be used within the validator function and even be precompiled for better
efficiency within regexes.go.

And the best reason, you can submit a pull request and we can keep on
adding to the validation library of this package!
Non standard validators
¶
A collection of validation rules that are frequently needed but are more
complex than the ones found in the baked in validators.
A non standard validator must be registered manually like you would
with your own custom validation functions.
Example of registration and use:
type Test struct {
TestField string `validate:"yourtag"`
}

t := &Test{
TestField: "Test"
}

validate := validator.New()
validate.RegisterValidation("yourtag", validators.NotBlank)
Here is a list of the current non standard validators:
NotBlank
This validates that the value is not blank or with length zero.
For strings ensures they do not contain only spaces. For channels, maps, slices and arrays
ensures they don't have zero length. For others, a non empty value is required.

Usage: notblank
Panics
¶
This package panics when bad input is provided, this is by design, bad code like
that should not make it to production.
type Test struct {
TestField string `validate:"nonexistentfunction=1"`
}

t := &Test{
TestField: "Test"
}

validate.Struct(t) // this will panic
Index
¶
type CustomTypeFunc
type FieldError
type FieldLevel
type FilterFunc
type Func
type FuncCtx
type InvalidValidationError
func (e *InvalidValidationError) Error() string
type Option
func WithPrivateFieldValidation() Option
func WithRequiredStructEnabled() Option
type RegisterTranslationsFunc
type StructLevel
type StructLevelFunc
type StructLevelFuncCtx
type TagNameFunc
type TranslationFunc
type Validate
func New(options ...Option) *Validate
func (v *Validate) RegisterAlias(alias, tags string)
func (v *Validate) RegisterCustomTypeFunc(fn CustomTypeFunc, types ...interface{})
func (v *Validate) RegisterStructValidation(fn StructLevelFunc, types ...interface{})
func (v *Validate) RegisterStructValidationCtx(fn StructLevelFuncCtx, types ...interface{})
func (v *Validate) RegisterStructValidationMapRules(rules map[string]string, types ...interface{})
func (v *Validate) RegisterTagNameFunc(fn TagNameFunc)
func (v *Validate) RegisterTranslation(tag string, trans ut.Translator, registerFn RegisterTranslationsFunc, ...) (err error)
func (v *Validate) RegisterValidation(tag string, fn Func, callValidationEvenIfNull ...bool) error
func (v *Validate) RegisterValidationCtx(tag string, fn FuncCtx, callValidationEvenIfNull ...bool) error
func (v *Validate) SetTagName(name string)
func (v *Validate) Struct(s interface{}) error
func (v *Validate) StructCtx(ctx context.Context, s interface{}) (err error)
func (v *Validate) StructExcept(s interface{}, fields ...string) error
func (v *Validate) StructExceptCtx(ctx context.Context, s interface{}, fields ...string) (err error)
func (v *Validate) StructFiltered(s interface{}, fn FilterFunc) error
func (v *Validate) StructFilteredCtx(ctx context.Context, s interface{}, fn FilterFunc) (err error)
func (v *Validate) StructPartial(s interface{}, fields ...string) error
func (v *Validate) StructPartialCtx(ctx context.Context, s interface{}, fields ...string) (err error)
func (v *Validate) ValidateMap(data map[string]interface{}, rules map[string]interface{}) map[string]interface{}
func (v Validate) ValidateMapCtx(ctx context.Context, data map[string]interface{}, rules map[string]interface{}) map[string]interface{}
func (v *Validate) Var(field interface{}, tag string) error
func (v *Validate) VarCtx(ctx context.Context, field interface{}, tag string) (err error)
func (v *Validate) VarWithKey(key string, field interface{}, tag string) error
func (v *Validate) VarWithKeyCtx(ctx context.Context, key string, field interface{}, tag string) (err error)
func (v *Validate) VarWithValue(field interface{}, other interface{}, tag string) error
func (v *Validate) VarWithValueCtx(ctx context.Context, field interface{}, other interface{}, tag string) (err error)
type ValidationErrors
func (ve ValidationErrors) Error() string
func (ve ValidationErrors) Translate(ut ut.Translator) ValidationErrorsTranslations
type ValidationErrorsTranslations
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
CustomTypeFunc
¶
type CustomTypeFunc func(field
reflect
.
Value
) interface{}
CustomTypeFunc allows for overriding or adding custom field type handler functions
field = field value of the type to return a value to be validated
example Valuer from sql drive see
https://golang.org/src/database/sql/driver/types.go?s=1210:1293#L29
type
FieldError
¶
type FieldError interface {
// Tag returns the validation tag that failed. if the
// validation was an alias, this will return the
// alias name and not the underlying tag that failed.
//
// eg. alias "iscolor": "hexcolor|rgb|rgba|hsl|hsla"
// will return "iscolor"
Tag()
string
// ActualTag returns the validation tag that failed, even if an
// alias the actual tag within the alias will be returned.
// If an 'or' validation fails the entire or will be returned.
//
// eg. alias "iscolor": "hexcolor|rgb|rgba|hsl|hsla"
// will return "hexcolor|rgb|rgba|hsl|hsla"
ActualTag()
string
// Namespace returns the namespace for the field error, with the tag
// name taking precedence over the field's actual name.
//
// eg. JSON name "User.fname"
//
// See StructNamespace() for a version that returns actual names.
//
// NOTE: this field can be blank when validating a single primitive field
// using validate.Field(...) as there is no way to extract it's name
Namespace()
string
// StructNamespace returns the namespace for the field error, with the field's
// actual name.
//
// eg. "User.FirstName" see Namespace for comparison
//
// NOTE: this field can be blank when validating a single primitive field
// using validate.Field(...) as there is no way to extract its name
StructNamespace()
string
// Field returns the field's name with the tag name taking precedence over the
// field's actual name.
//
// `RegisterTagNameFunc` must be registered to get tag value.
//
// eg. JSON name "fname"
// see StructField for comparison
Field()
string
// StructField returns the field's actual name from the struct, when able to determine.
//
// eg.  "FirstName"
// see Field for comparison
StructField()
string
// Value returns the actual field's value in case needed for creating the error
// message
Value() interface{}
// Param returns the param value, in string form for comparison; this will also
// help with generating an error message
Param()
string
// Kind returns the Field's reflect Kind
//
// eg. time.Time's kind is a struct
Kind()
reflect
.
Kind
// Type returns the Field's reflect Type
//
// eg. time.Time's type is time.Time
Type()
reflect
.
Type
// Translate returns the FieldError's translated error
// from the provided 'ut.Translator' and registered 'TranslationFunc'
//
// NOTE: if no registered translator can be found it returns the same as
// calling fe.Error()
Translate(ut
ut
.
Translator
)
string
// Error returns the FieldError's message
Error()
string
}
FieldError contains all functions to get error details
type
FieldLevel
¶
type FieldLevel interface {
// Top returns the top level struct, if any
Top()
reflect
.
Value
// Parent returns the current fields parent struct, if any or
// the comparison value if called 'VarWithValue'
Parent()
reflect
.
Value
// Field returns current field for validation
Field()
reflect
.
Value
// FieldName returns the field's name with the tag
// name taking precedence over the fields actual name.
FieldName()
string
// StructFieldName returns the struct field's name
StructFieldName()
string
// Param returns param for validation against current field
Param()
string
// GetTag returns the current validations tag name
GetTag()
string
// ExtractType gets the actual underlying type of field value.
// It will dive into pointers, customTypes and return you the
// underlying value and it's kind.
ExtractType(field
reflect
.
Value
) (value
reflect
.
Value
, kind
reflect
.
Kind
, nullable
bool
)
// GetStructFieldOK traverses the parent struct to retrieve a specific field denoted by the provided namespace
// in the param and returns the field, field kind and whether is was successful in retrieving
// the field at all.
//
// NOTE: when not successful ok will be false, this can happen when a nested struct is nil and so the field
// could not be retrieved because it didn't exist.
//
// Deprecated: Use GetStructFieldOK2() instead which also return if the value is nullable.
GetStructFieldOK() (
reflect
.
Value
,
reflect
.
Kind
,
bool
)
// GetStructFieldOKAdvanced is the same as GetStructFieldOK except that it accepts the parent struct to start looking for
// the field and namespace allowing more extensibility for validators.
//
// Deprecated: Use GetStructFieldOKAdvanced2() instead which also return if the value is nullable.
GetStructFieldOKAdvanced(val
reflect
.
Value
, namespace
string
) (
reflect
.
Value
,
reflect
.
Kind
,
bool
)
// GetStructFieldOK2 traverses the parent struct to retrieve a specific field denoted by the provided namespace
// in the param and returns the field, field kind, if it's a nullable type and whether is was successful in retrieving
// the field at all.
//
// NOTE: when not successful ok will be false, this can happen when a nested struct is nil and so the field
// could not be retrieved because it didn't exist.
GetStructFieldOK2() (
reflect
.
Value
,
reflect
.
Kind
,
bool
,
bool
)
// GetStructFieldOKAdvanced2 is the same as GetStructFieldOK except that it accepts the parent struct to start looking for
// the field and namespace allowing more extensibility for validators.
GetStructFieldOKAdvanced2(val
reflect
.
Value
, namespace
string
) (
reflect
.
Value
,
reflect
.
Kind
,
bool
,
bool
)
}
FieldLevel contains all the information and helper functions
to validate a field
type
FilterFunc
¶
type FilterFunc func(ns []
byte
)
bool
FilterFunc is the type used to filter fields using
StructFiltered(...) function.
returning true results in the field being filtered/skipped from
validation
type
Func
¶
type Func func(fl
FieldLevel
)
bool
Func accepts a FieldLevel interface for all validation needs. The return
value should be true when validation succeeds.
type
FuncCtx
¶
type FuncCtx func(ctx
context
.
Context
, fl
FieldLevel
)
bool
FuncCtx accepts a context.Context and FieldLevel interface for all
validation needs. The return value should be true when validation succeeds.
type
InvalidValidationError
¶
type InvalidValidationError struct {
Type
reflect
.
Type
}
InvalidValidationError describes an invalid argument passed to
`Struct`, `StructExcept`, StructPartial` or `Field`
func (*InvalidValidationError)
Error
¶
func (e *
InvalidValidationError
) Error()
string
Error returns InvalidValidationError message
type
Option
¶
added in
v10.15.2
type Option func(*
Validate
)
Option represents a configurations option to be applied to validator during initialization.
func
WithPrivateFieldValidation
¶
added in
v10.19.0
func WithPrivateFieldValidation()
Option
WithPrivateFieldValidation activates validation for unexported fields via the use of the `unsafe` package.
By opting into this feature you are acknowledging that you are aware of the risks and accept any current or future
consequences of using this feature.
func
WithRequiredStructEnabled
¶
added in
v10.15.2
func WithRequiredStructEnabled()
Option
WithRequiredStructEnabled enables required tag on non-pointer structs to be applied instead of ignored.
This was made opt-in behaviour in order to maintain backward compatibility with the behaviour previous
to being able to apply struct level validations on struct fields directly.
It is recommended you enabled this as it will be the default behaviour in v11+
type
RegisterTranslationsFunc
¶
type RegisterTranslationsFunc func(ut
ut
.
Translator
)
error
RegisterTranslationsFunc allows for registering of translations
for a 'ut.Translator' for use within the 'TranslationFunc'
type
StructLevel
¶
type StructLevel interface {
// Validator returns the main validation object, in case one wants to call validations internally.
// this is so you don't have to use anonymous functions to get access to the validate
// instance.
Validator() *
Validate
// Top returns the top level struct, if any
Top()
reflect
.
Value
// Parent returns the current fields parent struct, if any
Parent()
reflect
.
Value
// Current returns the current struct.
Current()
reflect
.
Value
// ExtractType gets the actual underlying type of field value.
// It will dive into pointers, customTypes and return you the
// underlying value and its kind.
ExtractType(field
reflect
.
Value
) (value
reflect
.
Value
, kind
reflect
.
Kind
, nullable
bool
)
// ReportError reports an error just by passing the field and tag information
//
// NOTES:
//
// fieldName and structFieldName get appended to the existing
// namespace that validator is on. e.g. pass 'FirstName' or
// 'Names[0]' depending on the nesting
//
// tag can be an existing validation tag or just something you make up
// and process on the flip side it's up to you.
ReportError(field interface{}, fieldName, structFieldName
string
, tag, param
string
)
// ReportValidationErrors reports an error just by passing ValidationErrors
//
// NOTES:
//
// relativeNamespace and relativeActualNamespace get appended to the
// existing namespace that validator is on.
// e.g. pass 'User.FirstName' or 'Users[0].FirstName' depending
// on the nesting. most of the time they will be blank, unless you validate
// at a level lower the current field depth
ReportValidationErrors(relativeNamespace, relativeActualNamespace
string
, errs
ValidationErrors
)
}
StructLevel contains all the information and helper functions
to validate a struct
type
StructLevelFunc
¶
type StructLevelFunc func(sl
StructLevel
)
StructLevelFunc accepts all values needed for struct level validation
type
StructLevelFuncCtx
¶
type StructLevelFuncCtx func(ctx
context
.
Context
, sl
StructLevel
)
StructLevelFuncCtx accepts all values needed for struct level validation
but also allows passing of contextual validation information via context.Context.
type
TagNameFunc
¶
type TagNameFunc func(field
reflect
.
StructField
)
string
TagNameFunc allows for adding of a custom tag name parser
type
TranslationFunc
¶
type TranslationFunc func(ut
ut
.
Translator
, fe
FieldError
)
string
TranslationFunc is the function type used to register or override
custom translations
type
Validate
¶
type Validate struct {
// contains filtered or unexported fields
}
Validate contains the validator settings and cache
func
New
¶
func New(options ...
Option
) *
Validate
New returns a new instance of 'validate' with sane defaults.
Validate is designed to be thread-safe and used as a singleton instance.
It caches information about your struct and validations,
in essence only parsing your validation tags once per struct type.
Using multiple instances neglects the benefit of caching.
func (*Validate)
RegisterAlias
¶
func (v *
Validate
) RegisterAlias(alias, tags
string
)
RegisterAlias registers a mapping of a single validation tag that
defines a common or complex set of validation(s) to simplify adding validation
to structs.
NOTE: this function is not thread-safe it is intended that these all be registered prior to any validation
func (*Validate)
RegisterCustomTypeFunc
¶
func (v *
Validate
) RegisterCustomTypeFunc(fn
CustomTypeFunc
, types ...interface{})
RegisterCustomTypeFunc registers a CustomTypeFunc against a number of types
NOTE: this method is not thread-safe it is intended that these all be registered prior to any validation
func (*Validate)
RegisterStructValidation
¶
func (v *
Validate
) RegisterStructValidation(fn
StructLevelFunc
, types ...interface{})
RegisterStructValidation registers a StructLevelFunc against a number of types.
NOTE:
- this method is not thread-safe it is intended that these all be registered prior to any validation
func (*Validate)
RegisterStructValidationCtx
¶
func (v *
Validate
) RegisterStructValidationCtx(fn
StructLevelFuncCtx
, types ...interface{})
RegisterStructValidationCtx registers a StructLevelFuncCtx against a number of types and allows passing
of contextual validation information via context.Context.
NOTE:
- this method is not thread-safe it is intended that these all be registered prior to any validation
func (*Validate)
RegisterStructValidationMapRules
¶
added in
v10.11.0
func (v *
Validate
) RegisterStructValidationMapRules(rules map[
string
]
string
, types ...interface{})
RegisterStructValidationMapRules registers validate map rules.
Be aware that map validation rules supersede those defined on a/the struct if present.
NOTE: this method is not thread-safe it is intended that these all be registered prior to any validation
func (*Validate)
RegisterTagNameFunc
¶
func (v *
Validate
) RegisterTagNameFunc(fn
TagNameFunc
)
RegisterTagNameFunc registers a function to get alternate names for StructFields.
eg. to use the names which have been specified for JSON representations of structs, rather than normal Go field names:
validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
// skip if tag key says it should be ignored
if name == "-" {
return ""
}
return name
})
func (*Validate)
RegisterTranslation
¶
func (v *
Validate
) RegisterTranslation(tag
string
, trans
ut
.
Translator
, registerFn
RegisterTranslationsFunc
, translationFn
TranslationFunc
) (err
error
)
RegisterTranslation registers translations against the provided tag.
func (*Validate)
RegisterValidation
¶
func (v *
Validate
) RegisterValidation(tag
string
, fn
Func
, callValidationEvenIfNull ...
bool
)
error
RegisterValidation adds a validation with the given tag
NOTES:
- if the key already exists, the previous validation function will be replaced.
- this method is not thread-safe it is intended that these all be registered prior to any validation
func (*Validate)
RegisterValidationCtx
¶
func (v *
Validate
) RegisterValidationCtx(tag
string
, fn
FuncCtx
, callValidationEvenIfNull ...
bool
)
error
RegisterValidationCtx does the same as RegisterValidation on accepts a FuncCtx validation
allowing context.Context validation support.
func (*Validate)
SetTagName
¶
func (v *
Validate
) SetTagName(name
string
)
SetTagName allows for changing of the default tag name of 'validate'
func (*Validate)
Struct
¶
func (v *
Validate
) Struct(s interface{})
error
Struct validates a structs exposed fields, and automatically validates nested structs, unless otherwise specified.
It returns InvalidValidationError for bad values passed in and nil or ValidationErrors as error otherwise.
You will need to assert the error if it's not nil eg. err.(validator.ValidationErrors) to access the array of errors.
func (*Validate)
StructCtx
¶
func (v *
Validate
) StructCtx(ctx
context
.
Context
, s interface{}) (err
error
)
StructCtx validates a structs exposed fields, and automatically validates nested structs, unless otherwise specified
and also allows passing of context.Context for contextual validation information.
It returns InvalidValidationError for bad values passed in and nil or ValidationErrors as error otherwise.
You will need to assert the error if it's not nil eg. err.(validator.ValidationErrors) to access the array of errors.
func (*Validate)
StructExcept
¶
func (v *
Validate
) StructExcept(s interface{}, fields ...
string
)
error
StructExcept validates all fields except the ones passed in.
Fields may be provided in a namespaced fashion relative to the  struct provided
i.e. NestedStruct.Field or NestedArrayField[0].Struct.Name
It returns InvalidValidationError for bad values passed in and nil or ValidationErrors as error otherwise.
You will need to assert the error if it's not nil eg. err.(validator.ValidationErrors) to access the array of errors.
func (*Validate)
StructExceptCtx
¶
func (v *
Validate
) StructExceptCtx(ctx
context
.
Context
, s interface{}, fields ...
string
) (err
error
)
StructExceptCtx validates all fields except the ones passed in and allows passing of contextual
validation information via context.Context
Fields may be provided in a namespaced fashion relative to the  struct provided
i.e. NestedStruct.Field or NestedArrayField[0].Struct.Name
It returns InvalidValidationError for bad values passed in and nil or ValidationErrors as error otherwise.
You will need to assert the error if it's not nil eg. err.(validator.ValidationErrors) to access the array of errors.
func (*Validate)
StructFiltered
¶
func (v *
Validate
) StructFiltered(s interface{}, fn
FilterFunc
)
error
StructFiltered validates a structs exposed fields, that pass the FilterFunc check and automatically validates
nested structs, unless otherwise specified.
It returns InvalidValidationError for bad values passed in and nil or ValidationErrors as error otherwise.
You will need to assert the error if it's not nil eg. err.(validator.ValidationErrors) to access the array of errors.
func (*Validate)
StructFilteredCtx
¶
func (v *
Validate
) StructFilteredCtx(ctx
context
.
Context
, s interface{}, fn
FilterFunc
) (err
error
)
StructFilteredCtx validates a structs exposed fields, that pass the FilterFunc check and automatically validates
nested structs, unless otherwise specified and also allows passing of contextual validation information via
context.Context
It returns InvalidValidationError for bad values passed in and nil or ValidationErrors as error otherwise.
You will need to assert the error if it's not nil eg. err.(validator.ValidationErrors) to access the array of errors.
func (*Validate)
StructPartial
¶
func (v *
Validate
) StructPartial(s interface{}, fields ...
string
)
error
StructPartial validates the fields passed in only, ignoring all others.
Fields may be provided in a namespaced fashion relative to the  struct provided
eg. NestedStruct.Field or NestedArrayField[0].Struct.Name
It returns InvalidValidationError for bad values passed in and nil or ValidationErrors as error otherwise.
You will need to assert the error if it's not nil eg. err.(validator.ValidationErrors) to access the array of errors.
func (*Validate)
StructPartialCtx
¶
func (v *
Validate
) StructPartialCtx(ctx
context
.
Context
, s interface{}, fields ...
string
) (err
error
)
StructPartialCtx validates the fields passed in only, ignoring all others and allows passing of contextual
validation information via context.Context
Fields may be provided in a namespaced fashion relative to the  struct provided
eg. NestedStruct.Field or NestedArrayField[0].Struct.Name
It returns InvalidValidationError for bad values passed in and nil or ValidationErrors as error otherwise.
You will need to assert the error if it's not nil eg. err.(validator.ValidationErrors) to access the array of errors.
func (*Validate)
ValidateMap
¶
added in
v10.6.0
func (v *
Validate
) ValidateMap(data map[
string
]interface{}, rules map[
string
]interface{}) map[
string
]interface{}
ValidateMap validates map data from a map of tags
func (Validate)
ValidateMapCtx
¶
added in
v10.6.0
func (v
Validate
) ValidateMapCtx(ctx
context
.
Context
, data map[
string
]interface{}, rules map[
string
]interface{}) map[
string
]interface{}
ValidateMapCtx validates a map using a map of validation rules and allows passing of contextual
validation information via context.Context.
func (*Validate)
Var
¶
func (v *
Validate
) Var(field interface{}, tag
string
)
error
Var validates a single variable using tag style validation.
eg.
var i int
validate.Var(i, "gt=1,lt=10")
WARNING: a struct can be passed for validation eg. time.Time is a struct or
if you have a custom type and have registered a custom type handler, so must
allow it; however unforeseen validations will occur if trying to validate a
struct that is meant to be passed to 'validate.Struct'
It returns InvalidValidationError for bad values passed in and nil or ValidationErrors as error otherwise.
You will need to assert the error if it's not nil eg. err.(validator.ValidationErrors) to access the array of errors.
validate Array, Slice and maps fields which may contain more than one error
func (*Validate)
VarCtx
¶
func (v *
Validate
) VarCtx(ctx
context
.
Context
, field interface{}, tag
string
) (err
error
)
VarCtx validates a single variable using tag style validation and allows passing of contextual
validation information via context.Context.
eg.
var i int
validate.Var(i, "gt=1,lt=10")
WARNING: a struct can be passed for validation eg. time.Time is a struct or
if you have a custom type and have registered a custom type handler, so must
allow it; however unforeseen validations will occur if trying to validate a
struct that is meant to be passed to 'validate.Struct'
It returns InvalidValidationError for bad values passed in and nil or ValidationErrors as error otherwise.
You will need to assert the error if it's not nil eg. err.(validator.ValidationErrors) to access the array of errors.
validate Array, Slice and maps fields which may contain more than one error
func (*Validate)
VarWithKey
¶
added in
v10.28.0
func (v *
Validate
) VarWithKey(key
string
, field interface{}, tag
string
)
error
VarWithKey validates a single variable with a key to be included in the returned error using tag style validation
eg.
var s string
validate.VarWithKey("email_address", s, "required,email")
WARNING: a struct can be passed for validation eg. time.Time is a struct or
if you have a custom type and have registered a custom type handler, so must
allow it; however unforeseen validations will occur if trying to validate a
struct that is meant to be passed to 'validate.Struct'
It returns InvalidValidationError for bad values passed in and nil or ValidationErrors as error otherwise.
You will need to assert the error if it's not nil eg. err.(validator.ValidationErrors) to access the array of errors.
validate Array, Slice and maps fields which may contain more than one error
func (*Validate)
VarWithKeyCtx
¶
added in
v10.28.0
func (v *
Validate
) VarWithKeyCtx(ctx
context
.
Context
, key
string
, field interface{}, tag
string
) (err
error
)
VarWithKeyCtx validates a single variable with a key to be included in the returned error using tag style validation
and allows passing of contextual validation information via context.Context.
eg.
var s string
validate.VarWithKeyCtx("email_address", s, "required,email")
WARNING: a struct can be passed for validation eg. time.Time is a struct or
if you have a custom type and have registered a custom type handler, so must
allow it; however unforeseen validations will occur if trying to validate a
struct that is meant to be passed to 'validate.Struct'
It returns InvalidValidationError for bad values passed in and nil or ValidationErrors as error otherwise.
You will need to assert the error if it's not nil eg. err.(validator.ValidationErrors) to access the array of errors.
validate Array, Slice and maps fields which may contain more than one error
func (*Validate)
VarWithValue
¶
func (v *
Validate
) VarWithValue(field interface{}, other interface{}, tag
string
)
error
VarWithValue validates a single variable, against another variable/field's value using tag style validation
eg.
s1 := "abcd"
s2 := "abcd"
validate.VarWithValue(s1, s2, "eqcsfield") // returns true
WARNING: a struct can be passed for validation eg. time.Time is a struct or
if you have a custom type and have registered a custom type handler, so must
allow it; however unforeseen validations will occur if trying to validate a
struct that is meant to be passed to 'validate.Struct'
It returns InvalidValidationError for bad values passed in and nil or ValidationErrors as error otherwise.
You will need to assert the error if it's not nil eg. err.(validator.ValidationErrors) to access the array of errors.
validate Array, Slice and maps fields which may contain more than one error
func (*Validate)
VarWithValueCtx
¶
func (v *
Validate
) VarWithValueCtx(ctx
context
.
Context
, field interface{}, other interface{}, tag
string
) (err
error
)
VarWithValueCtx validates a single variable, against another variable/field's value using tag style validation and
allows passing of contextual validation information via context.Context.
eg.
s1 := "abcd"
s2 := "abcd"
validate.VarWithValue(s1, s2, "eqcsfield") // returns true
WARNING: a struct can be passed for validation eg. time.Time is a struct or
if you have a custom type and have registered a custom type handler, so must
allow it; however unforeseen validations will occur if trying to validate a
struct that is meant to be passed to 'validate.Struct'
It returns InvalidValidationError for bad values passed in and nil or ValidationErrors as error otherwise.
You will need to assert the error if it's not nil eg. err.(validator.ValidationErrors) to access the array of errors.
validate Array, Slice and maps fields which may contain more than one error
type
ValidationErrors
¶
type ValidationErrors []
FieldError
ValidationErrors is an array of FieldError's
for use in custom error messages post validation.
func (ValidationErrors)
Error
¶
func (ve
ValidationErrors
) Error()
string
Error is intended for use in development + debugging and not intended to be a production error message.
It allows ValidationErrors to subscribe to the Error interface.
All information to create an error message specific to your application is contained within
the FieldError found within the ValidationErrors array
func (ValidationErrors)
Translate
¶
func (ve
ValidationErrors
) Translate(ut
ut
.
Translator
)
ValidationErrorsTranslations
Translate translates all of the ValidationErrors
type
ValidationErrorsTranslations
¶
type ValidationErrorsTranslations map[
string
]
string
ValidationErrorsTranslations is the translation return type