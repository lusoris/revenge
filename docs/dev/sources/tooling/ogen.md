# ogen OpenAPI Generator

> Source: https://pkg.go.dev/github.com/ogen-go/ogen
> Fetched: 2026-01-30T23:48:29.138535+00:00
> Content-Hash: fc2c9b30523560cc
> Type: html

---

Overview

¶

Package ogen implements OpenAPI v3 code generation.

Index

¶

type AdditionalProperties

func (p AdditionalProperties) MarshalJSON() ([]byte, error)

func (p AdditionalProperties) MarshalYAML() (any, error)

func (p *AdditionalProperties) ToJSONSchema() *jsonschema.AdditionalProperties

func (p *AdditionalProperties) UnmarshalJSON(data []byte) error

func (p *AdditionalProperties) UnmarshalYAML(node *yaml.Node) error

type Callback

type Components

func (c *Components) Init()

type Contact

func NewContact() *Contact

func (c *Contact) SetEmail(e string) *Contact

func (c *Contact) SetName(n string) *Contact

func (c *Contact) SetURL(url string) *Contact

type Default

type Discriminator

func (d *Discriminator) ToJSONSchema() *jsonschema.RawDiscriminator

type Encoding

type Enum

type Example

type ExampleValue

type Extensions

type ExternalDocumentation

type Header

type Info

func NewInfo() *Info

func (i *Info) SetContact(c *Contact) *Info

func (i *Info) SetDescription(d string) *Info

func (i *Info) SetLicense(l *License) *Info

func (i *Info) SetTermsOfService(t string) *Info

func (i *Info) SetTitle(t string) *Info

func (i *Info) SetVersion(v string) *Info

type Items

func (p Items) MarshalJSON() ([]byte, error)

func (p Items) MarshalYAML() (any, error)

func (p *Items) ToJSONSchema() *jsonschema.RawItems

func (p *Items) UnmarshalJSON(data []byte) error

func (p *Items) UnmarshalYAML(node *yaml.Node) error

type License

func NewLicense() *License

func (l *License) SetName(n string) *License

func (l *License) SetURL(url string) *License

type Link

type Locator

type Media

type NamedParameter

func NewNamedParameter(n string, p *Parameter) *NamedParameter

func (p *NamedParameter) AsLocalRef() *Parameter

type NamedPathItem

func NewNamedPathItem(n string, p *PathItem) *NamedPathItem

func (p *NamedPathItem) AsLocalRef() *PathItem

type NamedRequestBody

func NewNamedRequestBody(n string, p *RequestBody) *NamedRequestBody

func (p *NamedRequestBody) AsLocalRef() *RequestBody

type NamedResponse

func NewNamedResponse(n string, p *Response) *NamedResponse

func (p *NamedResponse) AsLocalRef() *Response

type NamedSchema

func NewNamedSchema(n string, p *Schema) *NamedSchema

func (p *NamedSchema) AsLocalRef() *Schema

type Num

type OAuthFlow

type OAuthFlows

type OpenAPICommon

type Operation

func NewOperation() *Operation

func (o *Operation) AddNamedResponses(ps ...*NamedResponse) *Operation

func (o *Operation) AddParameters(ps ...*Parameter) *Operation

func (o *Operation) AddResponse(n string, p *Response) *Operation

func (o *Operation) AddTags(ts ...string) *Operation

func (s *Operation) MarshalJSON() ([]byte, error)

func (o *Operation) SetDescription(d string) *Operation

func (o *Operation) SetOperationID(id string) *Operation

func (o *Operation) SetParameters(ps []*Parameter) *Operation

func (o *Operation) SetRequestBody(r *RequestBody) *Operation

func (o *Operation) SetResponses(r Responses) *Operation

func (o *Operation) SetSummary(s string) *Operation

func (o *Operation) SetTags(ts []string) *Operation

type Parameter

func NewParameter() *Parameter

func (p *Parameter) InCookie() *Parameter

func (p *Parameter) InHeader() *Parameter

func (p *Parameter) InPath() *Parameter

func (p *Parameter) InQuery() *Parameter

func (p *Parameter) SetContent(c map[string]Media) *Parameter

func (p *Parameter) SetDeprecated(d bool) *Parameter

func (p *Parameter) SetDescription(d string) *Parameter

func (p *Parameter) SetExplode(e bool) *Parameter

func (p *Parameter) SetIn(i string) *Parameter

func (p *Parameter) SetName(n string) *Parameter

func (p *Parameter) SetRef(r string) *Parameter

func (p *Parameter) SetRequired(r bool) *Parameter

func (p *Parameter) SetSchema(s *Schema) *Parameter

func (p *Parameter) SetStyle(s string) *Parameter

func (p *Parameter) ToNamed(n string) *NamedParameter

type PathItem

func NewPathItem() *PathItem

func (p *PathItem) AddParameters(ps ...*Parameter) *PathItem

func (p *PathItem) AddServers(srvs ...*Server) *PathItem

func (s *PathItem) MarshalJSON() ([]byte, error)

func (p *PathItem) SetDelete(o *Operation) *PathItem

func (p *PathItem) SetDescription(d string) *PathItem

func (p *PathItem) SetGet(o *Operation) *PathItem

func (p *PathItem) SetHead(o *Operation) *PathItem

func (p *PathItem) SetOptions(o *Operation) *PathItem

func (p *PathItem) SetParameters(ps []*Parameter) *PathItem

func (p *PathItem) SetPatch(o *Operation) *PathItem

func (p *PathItem) SetPost(o *Operation) *PathItem

func (p *PathItem) SetPut(o *Operation) *PathItem

func (p *PathItem) SetRef(r string) *PathItem

func (p *PathItem) SetServers(srvs []Server) *PathItem

func (p *PathItem) SetTrace(o *Operation) *PathItem

func (p *PathItem) ToNamed(n string) *NamedPathItem

type Paths

type PatternProperties

func (p PatternProperties) MarshalJSON() ([]byte, error)

func (p PatternProperties) MarshalYAML() (any, error)

func (p PatternProperties) ToJSONSchema() (result jsonschema.RawPatternProperties)

func (p *PatternProperties) UnmarshalJSON(data []byte) error

func (p *PatternProperties) UnmarshalYAML(node *yaml.Node) error

type PatternProperty

type Properties

func (p Properties) MarshalJSON() ([]byte, error)

func (p Properties) MarshalYAML() (any, error)

func (p Properties) ToJSONSchema() jsonschema.RawProperties

func (p *Properties) UnmarshalJSON(data []byte) error

func (p *Properties) UnmarshalYAML(node *yaml.Node) error

type Property

func NewProperty() *Property

func (p *Property) SetName(n string) *Property

func (p *Property) SetSchema(s *Schema) *Property

func (p Property) ToJSONSchema() jsonschema.RawProperty

type RawValue

type RequestBody

func NewRequestBody() *RequestBody

func (r *RequestBody) AddContent(mt string, s *Schema) *RequestBody

func (r *RequestBody) SetContent(c map[string]Media) *RequestBody

func (r *RequestBody) SetDescription(d string) *RequestBody

func (r *RequestBody) SetJSONContent(s *Schema) *RequestBody

func (r *RequestBody) SetRef(ref string) *RequestBody

func (r *RequestBody) SetRequired(req bool) *RequestBody

func (r *RequestBody) ToNamed(n string) *NamedRequestBody

type Response

func NewResponse() *Response

func (r *Response) AddContent(mt string, s *Schema) *Response

func (r *Response) SetContent(c map[string]Media) *Response

func (r *Response) SetDescription(d string) *Response

func (r *Response) SetHeaders(h map[string]*Header) *Response

func (r *Response) SetJSONContent(s *Schema) *Response

func (r *Response) SetLinks(l map[string]*Link) *Response

func (r *Response) SetRef(ref string) *Response

func (r *Response) ToNamed(n string) *NamedResponse

type Responses

type Schema

func Binary() *Schema

func Bool() *Schema

func Bytes() *Schema

func Date() *Schema

func DateTime() *Schema

func Double() *Schema

func Float() *Schema

func Int() *Schema

func Int32() *Schema

func Int64() *Schema

func NewSchema() *Schema

func Password() *Schema

func String() *Schema

func UUID() *Schema

func (s *Schema) AddOptionalProperties(ps ...*Property) *Schema

func (s *Schema) AddRequiredProperties(ps ...*Property) *Schema

func (s *Schema) AsArray() *Schema

func (s *Schema) AsEnum(def json.RawMessage, values ...json.RawMessage) *Schema

func (s *Schema) SetAllOf(a []*Schema) *Schema

func (s *Schema) SetAnyOf(a []*Schema) *Schema

func (s *Schema) SetDefault(d json.RawMessage) *Schema

func (s *Schema) SetDeprecated(d bool) *Schema

func (s *Schema) SetDescription(d string) *Schema

func (s *Schema) SetDiscriminator(d *Discriminator) *Schema

func (s *Schema) SetEnum(e []json.RawMessage) *Schema

func (s *Schema) SetExclusiveMaximum(e bool) *Schema

func (s *Schema) SetExclusiveMinimum(e bool) *Schema

func (s *Schema) SetFormat(f string) *Schema

func (s *Schema) SetItems(i *Schema) *Schema

func (s *Schema) SetMaxItems(m *uint64) *Schema

func (s *Schema) SetMaxLength(m *uint64) *Schema

func (s *Schema) SetMaxProperties(m *uint64) *Schema

func (s *Schema) SetMaximum(m *int64) *Schema

func (s *Schema) SetMinItems(m *uint64) *Schema

func (s *Schema) SetMinLength(m *uint64) *Schema

func (s *Schema) SetMinProperties(m *uint64) *Schema

func (s *Schema) SetMinimum(m *int64) *Schema

func (s *Schema) SetMultipleOf(m *uint64) *Schema

func (s *Schema) SetNullable(n bool) *Schema

func (s *Schema) SetOneOf(o []*Schema) *Schema

func (s *Schema) SetPattern(p string) *Schema

func (s *Schema) SetProperties(p *Properties) *Schema

func (s *Schema) SetRef(r string) *Schema

func (s *Schema) SetRequired(r []string) *Schema

func (s *Schema) SetSummary(smry string) *Schema

func (s *Schema) SetType(t string) *Schema

func (s *Schema) SetUniqueItems(u bool) *Schema

func (s *Schema) ToJSONSchema() *jsonschema.RawSchema

func (s *Schema) ToNamed(n string) *NamedSchema

func (s *Schema) ToProperty(n string) *Property

type SecurityRequirement

type SecurityRequirements

type SecurityScheme

type Server

func NewServer() *Server

func (s *Server) SetDescription(d string) *Server

func (s *Server) SetURL(url string) *Server

type ServerVariable

type Spec

func NewSpec() *Spec

func Parse(data []byte) (s *Spec, err error)

func (s *Spec) AddNamedParameters(ps ...*NamedParameter) *Spec

func (s *Spec) AddNamedPathItems(ps ...*NamedPathItem) *Spec

func (s *Spec) AddNamedRequestBodies(scs ...*NamedRequestBody) *Spec

func (s *Spec) AddNamedResponses(scs ...*NamedResponse) *Spec

func (s *Spec) AddNamedSchemas(scs ...*NamedSchema) *Spec

func (s *Spec) AddParameter(n string, p *Parameter) *Spec

func (s *Spec) AddPathItem(n string, p *PathItem) *Spec

func (s *Spec) AddRequestBody(n string, sc *RequestBody) *Spec

func (s *Spec) AddResponse(n string, sc *Response) *Spec

func (s *Spec) AddSchema(n string, sc *Schema) *Spec

func (s *Spec) AddServers(srvs ...*Server) *Spec

func (s *Spec) Init()

func (s *Spec) RefRequestBody(n string) *NamedRequestBody

func (s *Spec) RefResponse(n string) *NamedResponse

func (s *Spec) RefSchema(n string) *NamedSchema

func (s *Spec) SetComponents(c *Components) *Spec

func (s *Spec) SetInfo(i *Info) *Spec

func (s *Spec) SetOpenAPI(v string) *Spec

func (s *Spec) SetPaths(p Paths) *Spec

func (s *Spec) SetServers(srvs []Server) *Spec

func (s *Spec) UnmarshalJSON(bytes []byte) error

func (s *Spec) UnmarshalYAML(n *yaml.Node) error

type Tag

type XML

func (d *XML) ToJSONSchema() *jsonschema.XML

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

AdditionalProperties

¶

added in

v0.9.0

type AdditionalProperties struct {

Bool   *

bool

Schema

Schema

}

AdditionalProperties is JSON Schema additionalProperties validator description.

func (AdditionalProperties)

MarshalJSON

¶

added in

v0.9.0

func (p

AdditionalProperties

) MarshalJSON() ([]

byte

,

error

)

MarshalJSON implements json.Marshaler.

func (AdditionalProperties)

MarshalYAML

¶

added in

v0.44.0

func (p

AdditionalProperties

) MarshalYAML() (

any

,

error

)

MarshalYAML implements yaml.Marshaler.

func (*AdditionalProperties)

ToJSONSchema

¶

added in

v0.13.0

func (p *

AdditionalProperties

) ToJSONSchema() *

jsonschema

.

AdditionalProperties

ToJSONSchema converts AdditionalProperties to jsonschema.AdditionalProperties.

func (*AdditionalProperties)

UnmarshalJSON

¶

added in

v0.9.0

func (p *

AdditionalProperties

) UnmarshalJSON(data []

byte

)

error

UnmarshalJSON implements json.Unmarshaler.

func (*AdditionalProperties)

UnmarshalYAML

¶

added in

v0.43.0

func (p *

AdditionalProperties

) UnmarshalYAML(node *

yaml

.

Node

)

error

UnmarshalYAML implements yaml.Unmarshaler.

type

Callback

¶

added in

v0.44.0

type Callback map[

string

]*

PathItem

Callback is a map of possible out-of band callbacks related to the parent operation.

Each value in the map is a Path Item Object that describes a set of requests that may be
initiated by the API provider and the expected responses.

The key value used to identify the path item object is an expression, evaluated at runtime,
that identifies a URL to use for the callback operation.

To describe incoming requests from the API provider independent from another
API call, use the `webhooks` field.

See

https://spec.openapis.org/oas/v3.1.0#callback-object

.

type

Components

¶

type Components struct {

// An object to hold reusable Schema Objects.

Schemas map[

string

]*

Schema

`json:"schemas,omitempty" yaml:"schemas,omitempty"`

// An object to hold reusable Response Objects.

Responses map[

string

]*

Response

`json:"responses,omitempty" yaml:"responses,omitempty"`

// An object to hold reusable Parameter Objects.

Parameters map[

string

]*

Parameter

`json:"parameters,omitempty" yaml:"parameters,omitempty"`

// An object to hold reusable Example Objects.

Examples map[

string

]*

Example

`json:"examples,omitempty" yaml:"examples,omitempty"`

// An object to hold reusable Request Body Objects.

RequestBodies map[

string

]*

RequestBody

`json:"requestBodies,omitempty" yaml:"requestBodies,omitempty"`

// An object to hold reusable Header Objects.

Headers map[

string

]*

Header

`json:"headers,omitempty" yaml:"headers,omitempty"`

// An object to hold reusable Security Scheme Objects.

SecuritySchemes map[

string

]*

SecurityScheme

`json:"securitySchemes,omitempty" yaml:"securitySchemes,omitempty"`

// An object to hold reusable Link Objects.

Links map[

string

]*

Link

`json:"links,omitempty" yaml:"links,omitempty"`

// An object to hold reusable Callback Objects.

Callbacks map[

string

]*

Callback

`json:"callbacks,omitempty" yaml:"callbacks,omitempty"`

// An object to hold reusable Path Item Objects.

PathItems map[

string

]*

PathItem

`json:"pathItems,omitempty" yaml:"pathItems,omitempty"`

Common

OpenAPICommon

`json:"-" yaml:",inline"`

}

Components Holds a set of reusable objects for different aspects of the OAS.
All objects defined within the components object will have no effect on the API unless
they are explicitly referenced from properties outside the components object.

See

https://spec.openapis.org/oas/v3.1.0#components-object

.

func (*Components)

Init

¶

added in

v0.42.0

func (c *

Components

) Init()

Init initializes all fields.

type

Contact

¶

type Contact struct {

// The identifying name of the contact person/organization.

Name

string

`json:"name,omitempty" yaml:"name,omitempty"`

// The URL pointing to the contact information.

URL

string

`json:"url,omitempty" yaml:"url,omitempty"`

// The email address of the contact person/organization.

Email

string

`json:"email,omitempty" yaml:"email,omitempty"`

// Specification extensions.

Extensions

Extensions

`json:"-" yaml:",inline"`
}

Contact information for the exposed API.

See

https://spec.openapis.org/oas/v3.1.0#contact-object

.

func

NewContact

¶

func NewContact() *

Contact

NewContact returns a new Contact.

func (*Contact)

SetEmail

¶

func (c *

Contact

) SetEmail(e

string

) *

Contact

SetEmail sets the Email of the Contact.

func (*Contact)

SetName

¶

func (c *

Contact

) SetName(n

string

) *

Contact

SetName sets the Name of the Contact.

func (*Contact)

SetURL

¶

func (c *

Contact

) SetURL(url

string

) *

Contact

SetURL sets the URL of the Contact.

type

Default

¶

added in

v0.43.0

type Default =

jsonschema

.

Default

Default is a default value.

type

Discriminator

¶

type Discriminator struct {

// REQUIRED. The name of the property in the payload that will hold the discriminator value.

PropertyName

string

`json:"propertyName" yaml:"propertyName"`

// An object to hold mappings between payload values and schema names or references.

Mapping map[

string

]

string

`json:"mapping,omitempty" yaml:"mapping,omitempty"`

Common

OpenAPICommon

`json:"-" yaml:",inline"`

}

Discriminator discriminates types for OneOf, AllOf, AnyOf.

See

https://spec.openapis.org/oas/v3.1.0#discriminator-object

.

func (*Discriminator)

ToJSONSchema

¶

added in

v0.13.0

func (d *

Discriminator

) ToJSONSchema() *

jsonschema

.

RawDiscriminator

ToJSONSchema converts Discriminator to jsonschema.RawDiscriminator.

type

Encoding

¶

added in

v0.33.0

type Encoding struct {

// The Content-Type for encoding a specific property.

ContentType

string

`json:"contentType,omitempty" yaml:"contentType,omitempty"`

// A map allowing additional information to be provided as headers, for example Content-Disposition.

// Content-Type is described separately and SHALL be ignored in this section. This property SHALL be

// ignored if the request body media type is not a multipart.

Headers map[

string

]*

Header

`json:"headers,omitempty" yaml:"headers,omitempty"`

// Describes how the parameter value will be serialized

// depending on the type of the parameter value.

Style

string

`json:"style,omitempty" yaml:"style,omitempty"`

// When this is true, parameter values of type array or object

// generate separate parameters for each value of the array

// or key-value pair of the map.

// For other types of parameters this property has no effect.

Explode *

bool

`json:"explode,omitempty" yaml:"explode,omitempty"`

// Determines whether the parameter value SHOULD allow reserved characters, as defined by

// RFC3986 :/?#[]@!$&'()*+,;= to be included without percent-encoding.

// The default value is false. This property SHALL be ignored if the request body media type

// is not application/x-www-form-urlencoded.

AllowReserved

bool

`json:"allowReserved,omitempty" yaml:"allowReserved,omitempty"`

Common

OpenAPICommon

`json:"-" yaml:",inline"`

}

Encoding describes single encoding definition applied to a single schema property.

See

https://spec.openapis.org/oas/v3.1.0#encoding-object

.

type

Enum

¶

added in

v0.38.0

type Enum =

jsonschema

.

Enum

Enum is JSON Schema enum validator description.

type

Example

¶

type Example struct {

Ref

string

`json:"$ref,omitempty" yaml:"$ref,omitempty"`

// ref object

// Short description for the example.

Summary

string

`json:"summary,omitempty" yaml:"summary,omitempty"`

// Long description for the example.

// CommonMark syntax MAY be used for rich text representation.

Description

string

`json:"description,omitempty" yaml:"description,omitempty"`

// Embedded literal example.

Value

ExampleValue

`json:"value,omitempty" yaml:"value,omitempty"`

// A URI that points to the literal example.

ExternalValue

string

`json:"externalValue,omitempty" yaml:"externalValue,omitempty"`

Common

OpenAPICommon

`json:"-" yaml:",inline"`

}

Example object.

See

https://spec.openapis.org/oas/v3.1.0#example-object

.

type

ExampleValue

¶

added in

v0.43.0

type ExampleValue =

jsonschema

.

Example

ExampleValue is an example value.

type

Extensions

¶

added in

v0.49.0

type Extensions =

jsonschema

.

Extensions

Extensions is a map of OpenAPI extensions.

See

https://spec.openapis.org/oas/v3.1.0#specification-extensions

.

type

ExternalDocumentation

¶

added in

v0.40.0

type ExternalDocumentation struct {

// A description of the target documentation. CommonMark syntax MAY be used for rich text representation.

Description

string

`json:"description,omitempty" yaml:"description,omitempty"`

// REQUIRED. The URL for the target documentation. This MUST be in the form of a URL.

URL

string

`json:"url" yaml:"url"`

// Specification extensions.

Extensions

Extensions

`json:"-" yaml:",inline"`
}

ExternalDocumentation describes a reference to external resource for extended documentation.

See

https://spec.openapis.org/oas/v3.1.0#external-documentation-object

.

type

Header

¶

added in

v0.33.0

type Header =

Parameter

Header describes header response.

Header Object follows the structure of the Parameter Object with the following changes:

`name` MUST NOT be specified, it is given in the corresponding headers map.

`in` MUST NOT be specified, it is implicitly in header.

All traits that are affected by the location MUST be applicable to a location of header.

See

https://spec.openapis.org/oas/v3.1.0#header-object

.

type

Info

¶

type Info struct {

// REQUIRED. The title of the API.

Title

string

`json:"title" yaml:"title"`

// A short summary of the API.

Summary

string

`json:"summary,omitempty" yaml:"summary,omitempty"`

// A short description of the API.

// CommonMark syntax MAY be used for rich text representation.

Description

string

`json:"description,omitempty" yaml:"description,omitempty"`

// A URL to the Terms of Service for the API. MUST be in the format of a URL.

TermsOfService

string

`json:"termsOfService,omitempty" yaml:"termsOfService,omitempty"`

// The contact information for the exposed API.

Contact *

Contact

`json:"contact,omitempty" yaml:"contact,omitempty"`

// The license information for the exposed API.

License *

License

`json:"license,omitempty" yaml:"license,omitempty"`

// REQUIRED. The version of the OpenAPI document.

Version

string

`json:"version" yaml:"version"`

// Specification extensions.

Extensions

Extensions

`json:"-" yaml:",inline"`
}

Info provides metadata about the API.

The metadata MAY be used by the clients if needed,
and MAY be presented in editing or documentation generation tools for convenience.

See

https://spec.openapis.org/oas/v3.1.0#info-object

.

func

NewInfo

¶

func NewInfo() *

Info

NewInfo returns a new Info.

func (*Info)

SetContact

¶

func (i *

Info

) SetContact(c *

Contact

) *

Info

SetContact sets the Contact of the Info.

func (*Info)

SetDescription

¶

func (i *

Info

) SetDescription(d

string

) *

Info

SetDescription sets the description of the Info.

func (*Info)

SetLicense

¶

func (i *

Info

) SetLicense(l *

License

) *

Info

SetLicense sets the License of the Info.

func (*Info)

SetTermsOfService

¶

func (i *

Info

) SetTermsOfService(t

string

) *

Info

SetTermsOfService sets the terms of service of the Info.

func (*Info)

SetTitle

¶

func (i *

Info

) SetTitle(t

string

) *

Info

SetTitle sets the title of the Info.

func (*Info)

SetVersion

¶

func (i *

Info

) SetVersion(v

string

) *

Info

SetVersion sets the version of the Info.

type

Items

¶

added in

v0.69.0

type Items struct {

Item  *

Schema

Items []*

Schema

}

Items is unparsed JSON Schema items validator description.

func (Items)

MarshalJSON

¶

added in

v0.69.0

func (p

Items

) MarshalJSON() ([]

byte

,

error

)

MarshalJSON implements json.Marshaler.

func (Items)

MarshalYAML

¶

added in

v0.69.0

func (p

Items

) MarshalYAML() (

any

,

error

)

MarshalYAML implements yaml.Marshaler.

func (*Items)

ToJSONSchema

¶

added in

v0.69.0

func (p *

Items

) ToJSONSchema() *

jsonschema

.

RawItems

ToJSONSchema converts Items to jsonschema.RawItems.

func (*Items)

UnmarshalJSON

¶

added in

v0.69.0

func (p *

Items

) UnmarshalJSON(data []

byte

)

error

UnmarshalJSON implements json.Unmarshaler.

func (*Items)

UnmarshalYAML

¶

added in

v0.69.0

func (p *

Items

) UnmarshalYAML(node *

yaml

.

Node

)

error

UnmarshalYAML implements yaml.Unmarshaler.

type

License

¶

type License struct {

// REQUIRED. The license name used for the API.

Name

string

`json:"name" yaml:"name"`

// An SPDX license expression for the API.

Identifier

string

`json:"identifier,omitempty" yaml:"identifier,omitempty"`

// A URL to the license used for the API.

URL

string

`json:"url,omitempty" yaml:"url,omitempty"`

// Specification extensions.

Extensions

Extensions

`json:"-" yaml:",inline"`
}

License information for the exposed API.

See

https://spec.openapis.org/oas/v3.1.0#license-object

.

func

NewLicense

¶

func NewLicense() *

License

NewLicense returns a new License.

func (*License)

SetName

¶

func (l *

License

) SetName(n

string

) *

License

SetName sets the Name of the License.

func (*License)

SetURL

¶

func (l *

License

) SetURL(url

string

) *

License

SetURL sets the URL of the License.

type

Link

¶

added in

v0.44.0

type Link struct {

Ref

string

`json:"$ref,omitempty" yaml:"$ref,omitempty"`

// ref object

// A relative or absolute URI reference to an OAS operation.

//

// This field is mutually exclusive of the operationId field, and MUST point to an Operation Object.

OperationRef

string

`json:"operationRef,omitempty" yaml:"operationRef,omitempty"`

// The name of an existing, resolvable OAS operation, as defined with a unique operationId.

//

// This field is mutually exclusive of the operationRef field.

OperationID

string

`json:"operationId,omitempty" yaml:"operationId,omitempty"`

// A map representing parameters to pass to an operation as specified with operationId or identified

// via operationRef.

//

// The key is the parameter name to be used, whereas the value can be a constant or an expression to be

// evaluated and passed to the linked operation.

Parameters map[

string

]

RawValue

`json:"parameters,omitempty" yaml:"parameters,omitempty"`

Common

OpenAPICommon

`json:"-" yaml:",inline"`

}

Link describes a possible design-time link for a response.

See

https://spec.openapis.org/oas/v3.1.0#link-object

.

type

Locator

¶

added in

v0.42.0

type Locator =

location

.

Locator

Locator stores location of JSON value.

type

Media

¶

type Media struct {

// The schema defining the content of the request, response, or parameter.

Schema *

Schema

`json:"schema,omitempty" yaml:"schema,omitempty"`

// Example of the media type.

Example

ExampleValue

`json:"example,omitempty" yaml:"example,omitempty"`

// Examples of the media type.

Examples map[

string

]*

Example

`json:"examples,omitempty" yaml:"examples,omitempty"`

// A map between a property name and its encoding information.

//

// The key, being the property name, MUST exist in the schema as a property.

//

// The encoding object SHALL only apply to requestBody objects when the media

// type is multipart or application/x-www-form-urlencoded.

Encoding map[

string

]

Encoding

`json:"encoding,omitempty" yaml:"encoding,omitempty"`

Common

OpenAPICommon

`json:"-" yaml:",inline"`

}

Media provides schema and examples for the media type identified by its key.

See

https://spec.openapis.org/oas/v3.1.0#media-type-object

.

type

NamedParameter

¶

type NamedParameter struct {

Parameter *

Parameter

Name

string

}

NamedParameter can be used to construct a reference to the wrapped Parameter.

func

NewNamedParameter

¶

func NewNamedParameter(n

string

, p *

Parameter

) *

NamedParameter

NewNamedParameter returns a new NamedParameter.

func (*NamedParameter)

AsLocalRef

¶

func (p *

NamedParameter

) AsLocalRef() *

Parameter

AsLocalRef returns a new Parameter referencing the wrapped Parameter in the local document.

type

NamedPathItem

¶

type NamedPathItem struct {

PathItem *

PathItem

Name

string

}

NamedPathItem can be used to construct a reference to the wrapped PathItem.

func

NewNamedPathItem

¶

func NewNamedPathItem(n

string

, p *

PathItem

) *

NamedPathItem

NewNamedPathItem returns a new NamedPathItem.

func (*NamedPathItem)

AsLocalRef

¶

func (p *

NamedPathItem

) AsLocalRef() *

PathItem

AsLocalRef returns a new PathItem referencing the wrapped PathItem in the local document.

type

NamedRequestBody

¶

type NamedRequestBody struct {

RequestBody *

RequestBody

Name

string

}

NamedRequestBody can be used to construct a reference to the wrapped RequestBody.

func

NewNamedRequestBody

¶

func NewNamedRequestBody(n

string

, p *

RequestBody

) *

NamedRequestBody

NewNamedRequestBody returns a new NamedRequestBody.

func (*NamedRequestBody)

AsLocalRef

¶

func (p *

NamedRequestBody

) AsLocalRef() *

RequestBody

AsLocalRef returns a new RequestBody referencing the wrapped RequestBody in the local document.

type

NamedResponse

¶

type NamedResponse struct {

Response *

Response

Name

string

}

NamedResponse can be used to construct a reference to the wrapped Response.

func

NewNamedResponse

¶

func NewNamedResponse(n

string

, p *

Response

) *

NamedResponse

NewNamedResponse returns a new NamedResponse.

func (*NamedResponse)

AsLocalRef

¶

func (p *

NamedResponse

) AsLocalRef() *

Response

AsLocalRef returns a new Response referencing the wrapped Response in the local document.

type

NamedSchema

¶

type NamedSchema struct {

Schema *

Schema

Name

string

}

NamedSchema can be used to construct a reference to the wrapped Schema.

func

NewNamedSchema

¶

func NewNamedSchema(n

string

, p *

Schema

) *

NamedSchema

NewNamedSchema returns a new NamedSchema.

func (*NamedSchema)

AsLocalRef

¶

func (p *

NamedSchema

) AsLocalRef() *

Schema

AsLocalRef returns a new Schema referencing the wrapped Schema in the local document.

type

Num

¶

added in

v0.16.0

type Num =

jsonschema

.

Num

Num represents JSON number.

type

OAuthFlow

¶

added in

v0.19.0

type OAuthFlow struct {

// The authorization URL to be used for this flow.

// This MUST be in the form of a URL. The OAuth2 standard requires the use of TLS.

AuthorizationURL

string

`json:"authorizationUrl,omitempty" yaml:"authorizationUrl,omitempty"`

// The token URL to be used for this flow.

// This MUST be in the form of a URL. The OAuth2 standard requires the use of TLS.

TokenURL

string

`json:"tokenUrl,omitempty" yaml:"tokenUrl,omitempty"`

// The URL to be used for obtaining refresh tokens.

// This MUST be in the form of a URL. The OAuth2 standard requires the use of TLS.

RefreshURL

string

`json:"refreshUrl,omitempty" yaml:"refreshUrl,omitempty"`

// The available scopes for the OAuth2 security scheme.

// A map between the scope name and a short description for it. The map MAY be empty.

Scopes map[

string

]

string

`json:"scopes" yaml:"scopes"`

Common

OpenAPICommon

`json:"-" yaml:",inline"`

}

OAuthFlow is configuration details for a supported OAuth Flow.

See

https://spec.openapis.org/oas/v3.1.0#oauth-flow-object

.

type

OAuthFlows

¶

added in

v0.19.0

type OAuthFlows struct {

// Configuration for the OAuth Implicit flow.

Implicit *

OAuthFlow

`json:"implicit,omitempty" yaml:"implicit,omitempty"`

// Configuration for the OAuth Resource Owner Password flow.

Password *

OAuthFlow

`json:"password,omitempty" yaml:"password,omitempty"`

// Configuration for the OAuth Client Credentials flow. Previously called application in OpenAPI 2.0.

ClientCredentials *

OAuthFlow

`json:"clientCredentials,omitempty" yaml:"clientCredentials,omitempty"`

// Configuration for the OAuth Authorization Code flow. Previously called accessCode in OpenAPI 2.0.

AuthorizationCode *

OAuthFlow

`json:"authorizationCode,omitempty" yaml:"authorizationCode,omitempty"`

Common

OpenAPICommon

`json:"-" yaml:",inline"`

}

OAuthFlows allows configuration of the supported OAuth Flows.

See

https://spec.openapis.org/oas/v3.1.0#oauth-flows-object

.

type

OpenAPICommon

¶

added in

v0.49.0

type OpenAPICommon =

jsonschema

.

OpenAPICommon

OpenAPICommon is a common OpenAPI object fields (extensions and locator).

type

Operation

¶

type Operation struct {

// A list of tags for API documentation control.

// Tags can be used for logical grouping of operations by resources or any other qualifier.

Tags []

string

`json:"tags,omitempty" yaml:"tags,omitempty"`

// A short summary of what the operation does.

Summary

string

`json:"summary,omitempty" yaml:"summary,omitempty"`

// A verbose explanation of the operation behavior.

// CommonMark syntax MAY be used for rich text representation.

Description

string

`json:"description,omitempty" yaml:"description,omitempty"`

// Additional external documentation for this operation.

ExternalDocs *

ExternalDocumentation

`json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`

// Unique string used to identify the operation.

//

// The id MUST be unique among all operations described in the API.

//

// The operationId value is case-sensitive.

OperationID

string

`json:"operationId,omitempty" yaml:"operationId,omitempty"`

// A list of parameters that are applicable for this operation.

//

// If a parameter is already defined at the Path Item, the new definition will override it but

// can never remove it.

//

// The list MUST NOT include duplicated parameters. A unique parameter is defined by

// a combination of a name and location.

Parameters []*

Parameter

`json:"parameters,omitempty" yaml:"parameters,omitempty"`

// The request body applicable for this operation.

RequestBody *

RequestBody

`json:"requestBody,omitempty" yaml:"requestBody,omitempty"`

// The list of possible responses as they are returned from executing this operation.

Responses

Responses

`json:"responses,omitempty" yaml:"responses,omitempty"`

// A map of possible out-of band callbacks related to the parent operation.

//

// The key is a unique identifier for the Callback Object.

Callbacks map[

string

]*

Callback

`json:"callbacks,omitempty" yaml:"callbacks,omitempty"`

// Declares this operation to be deprecated

Deprecated

bool

`json:"deprecated,omitempty" yaml:"deprecated,omitempty"`

// A declaration of which security mechanisms can be used for this operation.

//

// The list of values includes alternative security requirement objects that can be used.

//

// Only one of the security requirement objects need to be satisfied to authorize a request.

Security

SecurityRequirements

`json:"security,omitempty" yaml:"security,omitempty"`

// An alternative server array to service this operation.

//

// If an alternative server object is specified at the Path Item Object or Root level,

// it will be overridden by this value.

Servers []

Server

`json:"servers,omitempty" yaml:"servers,omitempty"`

Common

OpenAPICommon

`json:"-" yaml:",inline"`

}

Operation describes a single API operation on a path.

See

https://spec.openapis.org/oas/v3.1.0#operation-object

.

func

NewOperation

¶

func NewOperation() *

Operation

NewOperation returns a new Operation.

func (*Operation)

AddNamedResponses

¶

func (o *

Operation

) AddNamedResponses(ps ...*

NamedResponse

) *

Operation

AddNamedResponses adds the given namedResponses to the Responses of the Operation.

func (*Operation)

AddParameters

¶

func (o *

Operation

) AddParameters(ps ...*

Parameter

) *

Operation

AddParameters adds Parameters to the Parameters of the Operation.

func (*Operation)

AddResponse

¶

func (o *

Operation

) AddResponse(n

string

, p *

Response

) *

Operation

AddResponse adds the given Response under the given Name to the Responses of the Operation.

func (*Operation)

AddTags

¶

func (o *

Operation

) AddTags(ts ...

string

) *

Operation

AddTags adds Tags to the Tags of the Operation.

func (*Operation)

MarshalJSON

¶

added in

v1.3.0

func (s *

Operation

) MarshalJSON() ([]

byte

,

error

)

MarshalJSON implements

json.Marshaler

.

func (*Operation)

SetDescription

¶

func (o *

Operation

) SetDescription(d

string

) *

Operation

SetDescription sets the Description of the Operation.

func (*Operation)

SetOperationID

¶

func (o *

Operation

) SetOperationID(id

string

) *

Operation

SetOperationID sets the OperationID of the Operation.

func (*Operation)

SetParameters

¶

func (o *

Operation

) SetParameters(ps []*

Parameter

) *

Operation

SetParameters sets the Parameters of the Operation.

func (*Operation)

SetRequestBody

¶

func (o *

Operation

) SetRequestBody(r *

RequestBody

) *

Operation

SetRequestBody sets the RequestBody of the Operation.

func (*Operation)

SetResponses

¶

func (o *

Operation

) SetResponses(r

Responses

) *

Operation

SetResponses sets the Responses of the Operation.

func (*Operation)

SetSummary

¶

func (o *

Operation

) SetSummary(s

string

) *

Operation

SetSummary sets the Summary of the Operation.

func (*Operation)

SetTags

¶

func (o *

Operation

) SetTags(ts []

string

) *

Operation

SetTags sets the Tags of the Operation.

type

Parameter

¶

type Parameter struct {

Ref

string

`json:"$ref,omitempty" yaml:"$ref,omitempty"`

// ref object

// REQUIRED. The name of the parameter. Parameter names are case sensitive.

Name

string

`json:"name,omitempty" yaml:"name,omitempty"`

// REQUIRED. The location of the parameter. Possible values are "query", "header", "path" or "cookie".

In

string

`json:"in,omitempty" yaml:"in,omitempty"`

// A brief description of the parameter. This could contain examples of use.

// CommonMark syntax MAY be used for rich text representation.

Description

string

`json:"description,omitempty" yaml:"description,omitempty"`

// Determines whether this parameter is mandatory.

// If the parameter location is "path", this property is REQUIRED

// and its value MUST be true.

// Otherwise, the property MAY be included and its default value is false.

Required

bool

`json:"required,omitempty" yaml:"required,omitempty"`

// Specifies that a parameter is deprecated and SHOULD be transitioned out of usage.

// Default value is false.

Deprecated

bool

`json:"deprecated,omitempty" yaml:"deprecated,omitempty"`

// Describes how the parameter value will be serialized

// depending on the type of the parameter value.

Style

string

`json:"style,omitempty" yaml:"style,omitempty"`

// When this is true, parameter values of type array or object

// generate separate parameters for each value of the array

// or key-value pair of the map.

// For other types of parameters this property has no effect.

Explode *

bool

`json:"explode,omitempty" yaml:"explode,omitempty"`

// Determines whether the parameter value SHOULD allow reserved characters, as defined by

RFC 3986

.

//

// This property only applies to parameters with an in value of query.

//

// The default value is false.

AllowReserved

bool

`json:"allowReserved,omitempty" yaml:"allowReserved,omitempty"`

// The schema defining the type used for the parameter.

Schema *

Schema

`json:"schema,omitempty" yaml:"schema,omitempty"`

// Example of the parameter's potential value.

Example

ExampleValue

`json:"example,omitempty" yaml:"example,omitempty"`

// Examples of the parameter's potential value.

Examples map[

string

]*

Example

`json:"examples,omitempty" yaml:"examples,omitempty"`

// For more complex scenarios, the content property can define the media type and schema of the parameter.

// A parameter MUST contain either a schema property, or a content property, but not both.

// When example or examples are provided in conjunction with the schema object,

// the example MUST follow the prescribed serialization strategy for the parameter.

//

// A map containing the representations for the parameter.

// The key is the media type and the value describes it.

// The map MUST only contain one entry.

Content map[

string

]

Media

`json:"content,omitempty" yaml:"content,omitempty"`

Common

OpenAPICommon

`json:"-" yaml:",inline"`

}

Parameter describes a single operation parameter.
A unique parameter is defined by a combination of a name and location.

See

https://spec.openapis.org/oas/v3.1.0#parameter-object

.

func

NewParameter

¶

func NewParameter() *

Parameter

NewParameter returns a new Parameter.

func (*Parameter)

InCookie

¶

func (p *

Parameter

) InCookie() *

Parameter

InCookie sets the In of the Parameter to "cookie".

func (*Parameter)

InHeader

¶

func (p *

Parameter

) InHeader() *

Parameter

InHeader sets the In of the Parameter to "header".

func (*Parameter)

InPath

¶

func (p *

Parameter

) InPath() *

Parameter

InPath sets the In of the Parameter to "path".

func (*Parameter)

InQuery

¶

func (p *

Parameter

) InQuery() *

Parameter

InQuery sets the In of the Parameter to "query".

func (*Parameter)

SetContent

¶

func (p *

Parameter

) SetContent(c map[

string

]

Media

) *

Parameter

SetContent sets the Content of the Parameter.

func (*Parameter)

SetDeprecated

¶

func (p *

Parameter

) SetDeprecated(d

bool

) *

Parameter

SetDeprecated sets the Deprecated of the Parameter.

func (*Parameter)

SetDescription

¶

func (p *

Parameter

) SetDescription(d

string

) *

Parameter

SetDescription sets the Description of the Parameter.

func (*Parameter)

SetExplode

¶

func (p *

Parameter

) SetExplode(e

bool

) *

Parameter

SetExplode sets the Explode of the Parameter.

func (*Parameter)

SetIn

¶

func (p *

Parameter

) SetIn(i

string

) *

Parameter

SetIn sets the In of the Parameter.

func (*Parameter)

SetName

¶

func (p *

Parameter

) SetName(n

string

) *

Parameter

SetName sets the Name of the Parameter.

func (*Parameter)

SetRef

¶

func (p *

Parameter

) SetRef(r

string

) *

Parameter

SetRef sets the Ref of the Parameter.

func (*Parameter)

SetRequired

¶

func (p *

Parameter

) SetRequired(r

bool

) *

Parameter

SetRequired sets the Required of the Parameter.

func (*Parameter)

SetSchema

¶

func (p *

Parameter

) SetSchema(s *

Schema

) *

Parameter

SetSchema sets the Schema of the Parameter.

func (*Parameter)

SetStyle

¶

func (p *

Parameter

) SetStyle(s

string

) *

Parameter

SetStyle sets the Style of the Parameter.

func (*Parameter)

ToNamed

¶

func (p *

Parameter

) ToNamed(n

string

) *

NamedParameter

ToNamed returns a NamedParameter wrapping the receiver.

type

PathItem

¶

type PathItem struct {

// Allows for an external definition of this path item.

// The referenced structure MUST be in the format of a Path Item Object.

// In case a Path Item Object field appears both

// in the defined object and the referenced object, the behavior is undefined.

Ref

string

`json:"$ref,omitempty" yaml:"$ref,omitempty"`

// An optional, string summary, intended to apply to all operations in this path.

Summary

string

`json:"summary,omitempty" yaml:"summary,omitempty"`

// An optional, string description, intended to apply to all operations in this path.

// CommonMark syntax MAY be used for rich text representation.

Description

string

`json:"description,omitempty" yaml:"description,omitempty"`

// A definition of a GET operation on this path.

Get *

Operation

`json:"get,omitempty" yaml:"get,omitempty"`

// A definition of a PUT operation on this path.

Put *

Operation

`json:"put,omitempty" yaml:"put,omitempty"`

// A definition of a POST operation on this path.

Post *

Operation

`json:"post,omitempty" yaml:"post,omitempty"`

// A definition of a DELETE operation on this path.

Delete *

Operation

`json:"delete,omitempty" yaml:"delete,omitempty"`

// A definition of a OPTIONS operation on this path.

Options *

Operation

`json:"options,omitempty" yaml:"options,omitempty"`

// A definition of a HEAD operation on this path.

Head *

Operation

`json:"head,omitempty" yaml:"head,omitempty"`

// A definition of a PATCH operation on this path.

Patch *

Operation

`json:"patch,omitempty" yaml:"patch,omitempty"`

// A definition of a TRACE operation on this path.

Trace *

Operation

`json:"trace,omitempty" yaml:"trace,omitempty"`

// An alternative server array to service all operations in this path.

Servers []

Server

`json:"servers,omitempty" yaml:"servers,omitempty"`

// A list of parameters that are applicable for all the operations described under this path.

//

// These parameters can be overridden at the operation level, but cannot be removed there.

//

// The list MUST NOT include duplicated parameters. A unique parameter is defined by

// a combination of a name and location.

Parameters []*

Parameter

`json:"parameters,omitempty" yaml:"parameters,omitempty"`

Common

OpenAPICommon

`json:"-" yaml:",inline"`

}

PathItem Describes the operations available on a single path.
A Path Item MAY be empty, due to ACL constraints. The path itself is still exposed to the
documentation viewer, but they will not know which operations and parameters are available.

See

https://spec.openapis.org/oas/v3.1.0#path-item-object

.

func

NewPathItem

¶

func NewPathItem() *

PathItem

NewPathItem returns a new PathItem.

func (*PathItem)

AddParameters

¶

func (p *

PathItem

) AddParameters(ps ...*

Parameter

) *

PathItem

AddParameters adds Parameters to the Parameters of the PathItem.

func (*PathItem)

AddServers

¶

func (p *

PathItem

) AddServers(srvs ...*

Server

) *

PathItem

AddServers adds Servers to the Servers of the PathItem.

func (*PathItem)

MarshalJSON

¶

added in

v1.3.0

func (s *

PathItem

) MarshalJSON() ([]

byte

,

error

)

MarshalJSON implements

json.Marshaler

.

func (*PathItem)

SetDelete

¶

func (p *

PathItem

) SetDelete(o *

Operation

) *

PathItem

SetDelete sets the Delete of the PathItem.

func (*PathItem)

SetDescription

¶

func (p *

PathItem

) SetDescription(d

string

) *

PathItem

SetDescription sets the Description of the PathItem.

func (*PathItem)

SetGet

¶

func (p *

PathItem

) SetGet(o *

Operation

) *

PathItem

SetGet sets the Get of the PathItem.

func (*PathItem)

SetHead

¶

func (p *

PathItem

) SetHead(o *

Operation

) *

PathItem

SetHead sets the Head of the PathItem.

func (*PathItem)

SetOptions

¶

func (p *

PathItem

) SetOptions(o *

Operation

) *

PathItem

SetOptions sets the Options of the PathItem.

func (*PathItem)

SetParameters

¶

func (p *

PathItem

) SetParameters(ps []*

Parameter

) *

PathItem

SetParameters sets the Parameters of the PathItem.

func (*PathItem)

SetPatch

¶

func (p *

PathItem

) SetPatch(o *

Operation

) *

PathItem

SetPatch sets the Patch of the PathItem.

func (*PathItem)

SetPost

¶

func (p *

PathItem

) SetPost(o *

Operation

) *

PathItem

SetPost sets the Post of the PathItem.

func (*PathItem)

SetPut

¶

func (p *

PathItem

) SetPut(o *

Operation

) *

PathItem

SetPut sets the Put of the PathItem.

func (*PathItem)

SetRef

¶

func (p *

PathItem

) SetRef(r

string

) *

PathItem

SetRef sets the Ref of the PathItem.

func (*PathItem)

SetServers

¶

func (p *

PathItem

) SetServers(srvs []

Server

) *

PathItem

SetServers sets the Servers of the PathItem.

func (*PathItem)

SetTrace

¶

func (p *

PathItem

) SetTrace(o *

Operation

) *

PathItem

SetTrace sets the Trace of the PathItem.

func (*PathItem)

ToNamed

¶

func (p *

PathItem

) ToNamed(n

string

) *

NamedPathItem

ToNamed returns a NamedPathItem wrapping the receiver.

type

Paths

¶

type Paths map[

string

]*

PathItem

Paths holds the relative paths to the individual endpoints and their operations.
The path is appended to the URL from the Server Object in order to construct the full URL.
The Paths MAY be empty, due to ACL constraints.

See

https://spec.openapis.org/oas/v3.1.0#paths-object

.

type

PatternProperties

¶

added in

v0.23.0

type PatternProperties []

PatternProperty

PatternProperties is unparsed JSON Schema patternProperties validator description.

func (PatternProperties)

MarshalJSON

¶

added in

v0.23.0

func (p

PatternProperties

) MarshalJSON() ([]

byte

,

error

)

MarshalJSON implements json.Marshaler.

func (PatternProperties)

MarshalYAML

¶

added in

v0.44.0

func (p

PatternProperties

) MarshalYAML() (

any

,

error

)

MarshalYAML implements yaml.Marshaler.

func (PatternProperties)

ToJSONSchema

¶

added in

v0.23.0

func (p

PatternProperties

) ToJSONSchema() (result

jsonschema

.

RawPatternProperties

)

ToJSONSchema converts PatternProperties to jsonschema.RawPatternProperties.

func (*PatternProperties)

UnmarshalJSON

¶

added in

v0.23.0

func (p *

PatternProperties

) UnmarshalJSON(data []

byte

)

error

UnmarshalJSON implements json.Unmarshaler.

func (*PatternProperties)

UnmarshalYAML

¶

added in

v0.43.0

func (p *

PatternProperties

) UnmarshalYAML(node *

yaml

.

Node

)

error

UnmarshalYAML implements yaml.Unmarshaler.

type

PatternProperty

¶

added in

v0.23.0

type PatternProperty struct {

Pattern

string

Schema  *

Schema

}

PatternProperty is item of PatternProperties.

type

Properties

¶

type Properties []

Property

Properties is unparsed JSON Schema properties validator description.

func (Properties)

MarshalJSON

¶

func (p

Properties

) MarshalJSON() ([]

byte

,

error

)

MarshalJSON implements json.Marshaler.

func (Properties)

MarshalYAML

¶

added in

v0.44.0

func (p

Properties

) MarshalYAML() (

any

,

error

)

MarshalYAML implements yaml.Marshaler.

func (Properties)

ToJSONSchema

¶

added in

v0.13.0

func (p

Properties

) ToJSONSchema()

jsonschema

.

RawProperties

ToJSONSchema converts Properties to jsonschema.RawProperties.

func (*Properties)

UnmarshalJSON

¶

func (p *

Properties

) UnmarshalJSON(data []

byte

)

error

UnmarshalJSON implements json.Unmarshaler.

func (*Properties)

UnmarshalYAML

¶

added in

v0.43.0

func (p *

Properties

) UnmarshalYAML(node *

yaml

.

Node

)

error

UnmarshalYAML implements yaml.Unmarshaler.

type

Property

¶

type Property struct {

Name

string

Schema *

Schema

}

Property is item of Properties.

func

NewProperty

¶

func NewProperty() *

Property

NewProperty returns a new Property.

func (*Property)

SetName

¶

func (p *

Property

) SetName(n

string

) *

Property

SetName sets the Name of the Property.

func (*Property)

SetSchema

¶

func (p *

Property

) SetSchema(s *

Schema

) *

Property

SetSchema sets the Schema of the Property.

func (Property)

ToJSONSchema

¶

added in

v0.13.0

func (p

Property

) ToJSONSchema()

jsonschema

.

RawProperty

ToJSONSchema converts Property to jsonschema.Property.

type

RawValue

¶

added in

v0.44.0

type RawValue =

jsonschema

.

RawValue

RawValue is a raw JSON value.

type

RequestBody

¶

type RequestBody struct {

Ref

string

`json:"$ref,omitempty" yaml:"$ref,omitempty"`

// ref object

// A brief description of the request body. This could contain examples of use.

// CommonMark syntax MAY be used for rich text representation.

Description

string

`json:"description,omitempty" yaml:"description,omitempty"`

// REQUIRED. The content of the request body.

//

// The key is a media type or media type range and the value describes it.

//

// For requests that match multiple keys, only the most specific key is applicable.

// e.g. text/plain overrides text/*

Content map[

string

]

Media

`json:"content,omitempty" yaml:"content,omitempty"`

// Determines if the request body is required in the request. Defaults to false.

Required

bool

`json:"required,omitempty" yaml:"required,omitempty"`

Common

OpenAPICommon

`json:"-" yaml:",inline"`

}

RequestBody describes a single request body.

See

https://spec.openapis.org/oas/v3.1.0#request-body-object

.

func

NewRequestBody

¶

func NewRequestBody() *

RequestBody

NewRequestBody returns a new RequestBody.

func (*RequestBody)

AddContent

¶

func (r *

RequestBody

) AddContent(mt

string

, s *

Schema

) *

RequestBody

AddContent adds the given Schema under the MediaType to the Content of the Response.

func (*RequestBody)

SetContent

¶

func (r *

RequestBody

) SetContent(c map[

string

]

Media

) *

RequestBody

SetContent sets the Content of the RequestBody.

func (*RequestBody)

SetDescription

¶

func (r *

RequestBody

) SetDescription(d

string

) *

RequestBody

SetDescription sets the Description of the RequestBody.

func (*RequestBody)

SetJSONContent

¶

func (r *

RequestBody

) SetJSONContent(s *

Schema

) *

RequestBody

SetJSONContent sets the given Schema under the JSON MediaType to the Content of the Response.

func (*RequestBody)

SetRef

¶

func (r *

RequestBody

) SetRef(ref

string

) *

RequestBody

SetRef sets the Ref of the RequestBody.

func (*RequestBody)

SetRequired

¶

func (r *

RequestBody

) SetRequired(req

bool

) *

RequestBody

SetRequired sets the Required of the RequestBody.

func (*RequestBody)

ToNamed

¶

func (r *

RequestBody

) ToNamed(n

string

) *

NamedRequestBody

ToNamed returns a NamedRequestBody wrapping the receiver.

type

Response

¶

type Response struct {

Ref

string

`json:"$ref,omitempty" yaml:"$ref,omitempty"`

// ref object

// REQUIRED. A description of the response.

// CommonMark syntax MAY be used for rich text representation.

Description

string

`json:"description,omitempty" yaml:"description,omitempty"`

// Maps a header name to its definition.

//

// RFC7230 states header names are case insensitive.

//

// If a response header is defined with the name "Content-Type", it SHALL be ignored.

Headers map[

string

]*

Header

`json:"headers,omitempty" yaml:"headers,omitempty"`

// A map containing descriptions of potential response payloads.

//

// The key is a media type or media type range and the value describes it.

//

// For requests that match multiple keys, only the most specific key is applicable.

// e.g. text/plain overrides text/*

Content map[

string

]

Media

`json:"content,omitempty" yaml:"content,omitempty"`

// A map of operations links that can be followed from the response.

//

// The key of the map is a short name for the link, following the naming constraints

// of the names for Component Objects.

Links map[

string

]*

Link

`json:"links,omitempty" yaml:"links,omitempty"`

Common

OpenAPICommon

`json:"-" yaml:",inline"`

}

Response describes a single response from an API Operation,
including design-time, static links to operations based on the response.

See

https://spec.openapis.org/oas/v3.1.0#response-object

.

func

NewResponse

¶

func NewResponse() *

Response

NewResponse returns a new Response.

func (*Response)

AddContent

¶

func (r *

Response

) AddContent(mt

string

, s *

Schema

) *

Response

AddContent adds the given Schema under the MediaType to the Content of the Response.

func (*Response)

SetContent

¶

func (r *

Response

) SetContent(c map[

string

]

Media

) *

Response

SetContent sets the Content of the Response.

func (*Response)

SetDescription

¶

func (r *

Response

) SetDescription(d

string

) *

Response

SetDescription sets the Description of the Response.

func (*Response)

SetHeaders

¶

added in

v0.33.0

func (r *

Response

) SetHeaders(h map[

string

]*

Header

) *

Response

SetHeaders sets the Headers of the Response.

func (*Response)

SetJSONContent

¶

func (r *

Response

) SetJSONContent(s *

Schema

) *

Response

SetJSONContent sets the given Schema under the JSON MediaType to the Content of the Response.

func (*Response)

SetLinks

¶

func (r *

Response

) SetLinks(l map[

string

]*

Link

) *

Response

SetLinks sets the Links of the Response.

func (*Response)

SetRef

¶

func (r *

Response

) SetRef(ref

string

) *

Response

SetRef sets the Ref of the Response.

func (*Response)

ToNamed

¶

func (r *

Response

) ToNamed(n

string

) *

NamedResponse

ToNamed returns a NamedResponse wrapping the receiver.

type

Responses

¶

type Responses map[

string

]*

Response

Responses is a container for the expected responses of an operation.

The container maps the HTTP response code to the expected response.

The `default` MAY be used as a default response object for all HTTP
codes that are not covered individually by the Responses Object.

The Responses Object MUST contain at least one response code, and if only one
response code is provided it SHOULD be the response for a successful operation call.

See

https://spec.openapis.org/oas/v3.1.0#responses-object

.

type

Schema

¶

type Schema struct {

Ref

string

`json:"$ref,omitempty" yaml:"$ref,omitempty"`

// ref object

Summary

string

`json:"summary,omitempty" yaml:"summary,omitempty"`

Description

string

`json:"description,omitempty" yaml:"description,omitempty"`

// Additional external documentation for this schema.

ExternalDocs *

ExternalDocumentation

`json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`

// Value MUST be a string. Multiple types via an array are not supported.

Type

string

`json:"type,omitempty" yaml:"type,omitempty"`

// See Data Type Formats for further details (

https://swagger.io/specification/#data-type-format

).

// While relying on JSON Schema's defined formats,

// the OAS offers a few additional predefined formats.

Format

string

`json:"format,omitempty" yaml:"format,omitempty"`

// Property definitions MUST be a Schema Object and not a standard JSON Schema

// (inline or referenced).

Properties

Properties

`json:"properties,omitempty" yaml:"properties,omitempty"`

// Value can be boolean or object. Inline or referenced schema MUST be of a Schema Object

// and not a standard JSON Schema. Consistent with JSON Schema, additionalProperties defaults to true.

AdditionalProperties *

AdditionalProperties

`json:"additionalProperties,omitempty" yaml:"additionalProperties,omitempty"`

// The value of "patternProperties" MUST be an object. Each property

// name of this object SHOULD be a valid regular expression, according

// to the ECMA-262 regular expression dialect. Each property value of

// this object MUST be a valid JSON Schema.

PatternProperties

PatternProperties

`json:"patternProperties,omitempty" yaml:"patternProperties,omitempty"`

// The value of this keyword MUST be an array.

// This array MUST have at least one element.

// Elements of this array MUST be strings, and MUST be unique.

Required []

string

`json:"required,omitempty" yaml:"required,omitempty"`

// Value MUST be an object and not an array.

// Inline or referenced schema MUST be of a Schema Object and not a standard JSON Schema.

// MUST be present if the Type is "array".

Items *

Items

`json:"items,omitempty" yaml:"items,omitempty"`

// A true value adds "null" to the allowed type specified by the type keyword,

// only if type is explicitly defined within the same Schema Object.

// Other Schema Object constraints retain their defined behavior,

// and therefore may disallow the use of null as a value.

// A false value leaves the specified or default type unmodified.

// The default value is false.

Nullable

bool

`json:"nullable,omitempty" yaml:"nullable,omitempty"`

// AllOf takes an array of object definitions that are used

// for independent validation but together compose a single object.

// Still, it does not imply a hierarchy between the models.

// For that purpose, you should include the discriminator.

AllOf []*

Schema

`json:"allOf,omitempty" yaml:"allOf,omitempty"`

// OneOf validates the value against exactly one of the subschemas

OneOf []*

Schema

`json:"oneOf,omitempty" yaml:"oneOf,omitempty"`

// AnyOf validates the value against any (one or more) of the subschemas

AnyOf []*

Schema

`json:"anyOf,omitempty" yaml:"anyOf,omitempty"`

// Discriminator for subschemas.

Discriminator *

Discriminator

`json:"discriminator,omitempty" yaml:"discriminator,omitempty"`

// Adds additional metadata to describe the XML representation of this property.

//

// This MAY be used only on properties schemas. It has no effect on root schemas

XML *

XML

`json:"xml,omitempty" yaml:"xml,omitempty"`

// The value of this keyword MUST be an array.

// This array SHOULD have at least one element.

// Elements in the array SHOULD be unique.

Enum

Enum

`json:"enum,omitempty" yaml:"enum,omitempty"`

// The value of "multipleOf" MUST be a number, strictly greater than 0.

//

// A numeric instance is only valid if division by this keyword's value

// results in an integer.

MultipleOf

Num

`json:"multipleOf,omitempty" yaml:"multipleOf,omitempty"`

// The value of "maximum" MUST be a number, representing an upper limit

// for a numeric instance.

//

// If the instance is a number, then this keyword validates if

// "exclusiveMaximum" is true and instance is less than the provided

// value, or else if the instance is less than or exactly equal to the

// provided value.

Maximum

Num

`json:"maximum,omitempty" yaml:"maximum,omitempty"`

// The value of "exclusiveMaximum" MUST be a boolean, representing

// whether the limit in "maximum" is exclusive or not.  An undefined

// value is the same as false.

//

// If "exclusiveMaximum" is true, then a numeric instance SHOULD NOT be

// equal to the value specified in "maximum".  If "exclusiveMaximum" is

// false (or not specified), then a numeric instance MAY be equal to the

// value of "maximum".

ExclusiveMaximum

bool

`json:"exclusiveMaximum,omitempty" yaml:"exclusiveMaximum,omitempty"`

// The value of "minimum" MUST be a number, representing a lower limit

// for a numeric instance.

//

// If the instance is a number, then this keyword validates if

// "exclusiveMinimum" is true and instance is greater than the provided

// value, or else if the instance is greater than or exactly equal to

// the provided value.

Minimum

Num

`json:"minimum,omitempty" yaml:"minimum,omitempty"`

// The value of "exclusiveMinimum" MUST be a boolean, representing

// whether the limit in "minimum" is exclusive or not.  An undefined

// value is the same as false.

//

// If "exclusiveMinimum" is true, then a numeric instance SHOULD NOT be

// equal to the value specified in "minimum".  If "exclusiveMinimum" is

// false (or not specified), then a numeric instance MAY be equal to the

// value of "minimum".

ExclusiveMinimum

bool

`json:"exclusiveMinimum,omitempty" yaml:"exclusiveMinimum,omitempty"`

// The value of this keyword MUST be a non-negative integer.

//

// The value of this keyword MUST be an integer.  This integer MUST be

// greater than, or equal to, 0.

//

// A string instance is valid against this keyword if its length is less

// than, or equal to, the value of this keyword.

//

// The length of a string instance is defined as the number of its

// characters as defined by

RFC 7159

[RFC7159].

MaxLength *

uint64

`json:"maxLength,omitempty" yaml:"maxLength,omitempty"`

// A string instance is valid against this keyword if its length is

// greater than, or equal to, the value of this keyword.

//

// The length of a string instance is defined as the number of its

// characters as defined by

RFC 7159

[RFC7159].

//

// The value of this keyword MUST be an integer.  This integer MUST be

// greater than, or equal to, 0.

//

// "minLength", if absent, may be considered as being present with

// integer value 0.

MinLength *

uint64

`json:"minLength,omitempty" yaml:"minLength,omitempty"`

// The value of this keyword MUST be a string.  This string SHOULD be a

// valid regular expression, according to the ECMA 262 regular

// expression dialect.

//

// A string instance is considered valid if the regular expression

// matches the instance successfully. Recall: regular expressions are

// not implicitly anchored.

Pattern

string

`json:"pattern,omitempty" yaml:"pattern,omitempty"`

// The value of this keyword MUST be an integer.  This integer MUST be

// greater than, or equal to, 0.

//

// An array instance is valid against "maxItems" if its size is less

// than, or equal to, the value of this keyword.

MaxItems *

uint64

`json:"maxItems,omitempty" yaml:"maxItems,omitempty"`

// The value of this keyword MUST be an integer.  This integer MUST be

// greater than, or equal to, 0.

//

// An array instance is valid against "minItems" if its size is greater

// than, or equal to, the value of this keyword.

//

// If this keyword is not present, it may be considered present with a

// value of 0.

MinItems *

uint64

`json:"minItems,omitempty" yaml:"minItems,omitempty"`

// The value of this keyword MUST be a boolean.

//

// If this keyword has boolean value false, the instance validates

// successfully.  If it has boolean value true, the instance validates

// successfully if all of its elements are unique.

//

// If not present, this keyword may be considered present with boolean

// value false.

UniqueItems

bool

`json:"uniqueItems,omitempty" yaml:"uniqueItems,omitempty"`

// The value of this keyword MUST be an integer.  This integer MUST be

// greater than, or equal to, 0.

//

// An object instance is valid against "maxProperties" if its number of

// properties is less than, or equal to, the value of this keyword.

MaxProperties *

uint64

`json:"maxProperties,omitempty" yaml:"maxProperties,omitempty"`

// The value of this keyword MUST be an integer.  This integer MUST be

// greater than, or equal to, 0.

//

// An object instance is valid against "minProperties" if its number of

// properties is greater than, or equal to, the value of this keyword.

//

// If this keyword is not present, it may be considered present with a

// value of 0.

MinProperties *

uint64

`json:"minProperties,omitempty" yaml:"minProperties,omitempty"`

// Default value.

Default

Default

`json:"default,omitempty" yaml:"default,omitempty"`

// A free-form property to include an example of an instance for this schema.

// To represent examples that cannot be naturally represented in JSON or YAML,

// a string value can be used to contain the example with escaping where necessary.

Example

ExampleValue

`json:"example,omitempty" yaml:"example,omitempty"`

// Specifies that a schema is deprecated and SHOULD be transitioned out

// of usage.

Deprecated

bool

`json:"deprecated,omitempty" yaml:"deprecated,omitempty"`

// If the instance value is a string, this property defines that the

// string SHOULD be interpreted as binary data and decoded using the

// encoding named by this property.

RFC 2045, Section 6.1

lists

// the possible values for this property.

//

// The value of this property MUST be a string.

//

// The value of this property SHOULD be ignored if the instance

// described is not a string.

ContentEncoding

string

`json:"contentEncoding,omitempty" yaml:"contentEncoding,omitempty"`

// The value of this property must be a media type, as defined by RFC

// 2046. This property defines the media type of instances

// which this schema defines.

//

// The value of this property MUST be a string.

//

// The value of this property SHOULD be ignored if the instance

// described is not a string.

ContentMediaType

string

`json:"contentMediaType,omitempty" yaml:"contentMediaType,omitempty"`

Common

jsonschema

.

OpenAPICommon

`json:"-" yaml:",inline"`

}

The Schema Object allows the definition of input and output data types.
These types can be objects, but also primitives and arrays.

func

Binary

¶

func Binary() *

Schema

Binary returns a sequence of octets OAS data type (Schema).

func

Bool

¶

func Bool() *

Schema

Bool returns a boolean OAS data type (Schema).

func

Bytes

¶

func Bytes() *

Schema

Bytes returns a base64 encoded OAS data type (Schema).

func

Date

¶

func Date() *

Schema

Date returns a date as defined by full-date - RFC3339 OAS data type (Schema).

func

DateTime

¶

func DateTime() *

Schema

DateTime returns a date as defined by date-time - RFC3339 OAS data type (Schema).

func

Double

¶

func Double() *

Schema

Double returns a double OAS data type (Schema).

func

Float

¶

func Float() *

Schema

Float returns a float OAS data type (Schema).

func

Int

¶

added in

v0.2.0

func Int() *

Schema

Int returns an integer OAS data type (Schema).

func

Int32

¶

func Int32() *

Schema

Int32 returns an 32-bit integer OAS data type (Schema).

func

Int64

¶

func Int64() *

Schema

Int64 returns an 64-bit integer OAS data type (Schema).

func

NewSchema

¶

func NewSchema() *

Schema

NewSchema returns a new Schema.

func

Password

¶

func Password() *

Schema

Password returns an obscured OAS data type (Schema).

func

String

¶

func String() *

Schema

String returns a string OAS data type (Schema).

func

UUID

¶

added in

v0.2.0

func UUID() *

Schema

UUID returns a UUID OAS data type (Schema).

func (*Schema)

AddOptionalProperties

¶

func (s *

Schema

) AddOptionalProperties(ps ...*

Property

) *

Schema

AddOptionalProperties adds the Properties to the Properties of the Schema.

func (*Schema)

AddRequiredProperties

¶

func (s *

Schema

) AddRequiredProperties(ps ...*

Property

) *

Schema

AddRequiredProperties adds the Properties to the Properties of the Schema and marks them as required.

func (*Schema)

AsArray

¶

func (s *

Schema

) AsArray() *

Schema

AsArray returns a new "array" Schema wrapping the receiver.

func (*Schema)

AsEnum

¶

func (s *

Schema

) AsEnum(def

json

.

RawMessage

, values ...

json

.

RawMessage

) *

Schema

AsEnum returns a new "enum" Schema wrapping the receiver.

func (*Schema)

SetAllOf

¶

func (s *

Schema

) SetAllOf(a []*

Schema

) *

Schema

SetAllOf sets the AllOf of the Schema.

func (*Schema)

SetAnyOf

¶

func (s *

Schema

) SetAnyOf(a []*

Schema

) *

Schema

SetAnyOf sets the AnyOf of the Schema.

func (*Schema)

SetDefault

¶

func (s *

Schema

) SetDefault(d

json

.

RawMessage

) *

Schema

SetDefault sets the Default of the Schema.

func (*Schema)

SetDeprecated

¶

added in

v0.68.0

func (s *

Schema

) SetDeprecated(d

bool

) *

Schema

SetDeprecated sets the Deprecated of the Schema.

func (*Schema)

SetDescription

¶

func (s *

Schema

) SetDescription(d

string

) *

Schema

SetDescription sets the Description of the Schema.

func (*Schema)

SetDiscriminator

¶

func (s *

Schema

) SetDiscriminator(d *

Discriminator

) *

Schema

SetDiscriminator sets the Discriminator of the Schema.

func (*Schema)

SetEnum

¶

func (s *

Schema

) SetEnum(e []

json

.

RawMessage

) *

Schema

SetEnum sets the Enum of the Schema.

func (*Schema)

SetExclusiveMaximum

¶

func (s *

Schema

) SetExclusiveMaximum(e

bool

) *

Schema

SetExclusiveMaximum sets the ExclusiveMaximum of the Schema.

func (*Schema)

SetExclusiveMinimum

¶

func (s *

Schema

) SetExclusiveMinimum(e

bool

) *

Schema

SetExclusiveMinimum sets the ExclusiveMinimum of the Schema.

func (*Schema)

SetFormat

¶

func (s *

Schema

) SetFormat(f

string

) *

Schema

SetFormat sets the Format of the Schema.

func (*Schema)

SetItems

¶

func (s *

Schema

) SetItems(i *

Schema

) *

Schema

SetItems sets the Items of the Schema.

func (*Schema)

SetMaxItems

¶

func (s *

Schema

) SetMaxItems(m *

uint64

) *

Schema

SetMaxItems sets the MaxItems of the Schema.

func (*Schema)

SetMaxLength

¶

func (s *

Schema

) SetMaxLength(m *

uint64

) *

Schema

SetMaxLength sets the MaxLength of the Schema.

func (*Schema)

SetMaxProperties

¶

func (s *

Schema

) SetMaxProperties(m *

uint64

) *

Schema

SetMaxProperties sets the MaxProperties of the Schema.

func (*Schema)

SetMaximum

¶

func (s *

Schema

) SetMaximum(m *

int64

) *

Schema

SetMaximum sets the Maximum of the Schema.

func (*Schema)

SetMinItems

¶

func (s *

Schema

) SetMinItems(m *

uint64

) *

Schema

SetMinItems sets the MinItems of the Schema.

func (*Schema)

SetMinLength

¶

func (s *

Schema

) SetMinLength(m *

uint64

) *

Schema

SetMinLength sets the MinLength of the Schema.

func (*Schema)

SetMinProperties

¶

func (s *

Schema

) SetMinProperties(m *

uint64

) *

Schema

SetMinProperties sets the MinProperties of the Schema.

func (*Schema)

SetMinimum

¶

func (s *

Schema

) SetMinimum(m *

int64

) *

Schema

SetMinimum sets the Minimum of the Schema.

func (*Schema)

SetMultipleOf

¶

func (s *

Schema

) SetMultipleOf(m *

uint64

) *

Schema

SetMultipleOf sets the MultipleOf of the Schema.

func (*Schema)

SetNullable

¶

func (s *

Schema

) SetNullable(n

bool

) *

Schema

SetNullable sets the Nullable of the Schema.

func (*Schema)

SetOneOf

¶

func (s *

Schema

) SetOneOf(o []*

Schema

) *

Schema

SetOneOf sets the OneOf of the Schema.

func (*Schema)

SetPattern

¶

func (s *

Schema

) SetPattern(p

string

) *

Schema

SetPattern sets the Pattern of the Schema.

func (*Schema)

SetProperties

¶

func (s *

Schema

) SetProperties(p *

Properties

) *

Schema

SetProperties sets the Properties of the Schema.

func (*Schema)

SetRef

¶

func (s *

Schema

) SetRef(r

string

) *

Schema

SetRef sets the Ref of the Schema.

func (*Schema)

SetRequired

¶

func (s *

Schema

) SetRequired(r []

string

) *

Schema

SetRequired sets the Required of the Schema.

func (*Schema)

SetSummary

¶

added in

v0.68.3

func (s *

Schema

) SetSummary(smry

string

) *

Schema

SetSummary sets the Summary of the Schema.

func (*Schema)

SetType

¶

func (s *

Schema

) SetType(t

string

) *

Schema

SetType sets the Type of the Schema.

func (*Schema)

SetUniqueItems

¶

func (s *

Schema

) SetUniqueItems(u

bool

) *

Schema

SetUniqueItems sets the UniqueItems of the Schema.

func (*Schema)

ToJSONSchema

¶

added in

v0.13.0

func (s *

Schema

) ToJSONSchema() *

jsonschema

.

RawSchema

ToJSONSchema converts Schema to jsonschema.Schema.

func (*Schema)

ToNamed

¶

func (s *

Schema

) ToNamed(n

string

) *

NamedSchema

ToNamed returns a NamedSchema wrapping the receiver.

func (*Schema)

ToProperty

¶

func (s *

Schema

) ToProperty(n

string

) *

Property

ToProperty returns a Property with the given name and with this Schema.

type

SecurityRequirement

¶

added in

v0.57.0

type SecurityRequirement = map[

string

][]

string

SecurityRequirement lists the required security schemes to execute this operation.

See

https://spec.openapis.org/oas/v3.1.0#security-requirement-object

.

type

SecurityRequirements

¶

added in

v0.19.0

type SecurityRequirements []

SecurityRequirement

SecurityRequirements lists the security requirements of the operation.

type

SecurityScheme

¶

added in

v0.42.0

type SecurityScheme struct {

Ref

string

`json:"$ref,omitempty" yaml:"$ref,omitempty"`

// The type of the security scheme. Valid values are "apiKey", "http", "mutualTLS", "oauth2", "openIdConnect".

Type

string

`json:"type" yaml:"type,omitempty"`

// A description for security scheme. CommonMark syntax MAY be used for rich text representation.

Description

string

`json:"description,omitempty" yaml:"description,omitempty"`

// The name of the header, query or cookie parameter to be used.

Name

string

`json:"name,omitempty" yaml:"name,omitempty"`

// The location of the API key. Valid values are "query", "header" or "cookie".

In

string

`json:"in,omitempty" yaml:"in,omitempty"`

// The name of the HTTP Authorization scheme to be used in the Authorization header as defined in RFC7235.

// The values used SHOULD be registered in the IANA Authentication Scheme registry.

Scheme

string

`json:"scheme,omitempty" yaml:"scheme,omitempty"`

// A hint to the client to identify how the bearer token is formatted. Bearer tokens are usually generated

// by an authorization server, so this information is primarily for documentation purposes.

BearerFormat

string

`json:"bearerFormat,omitempty" yaml:"bearerFormat,omitempty"`

// An object containing configuration information for the flow types supported.

Flows *

OAuthFlows

`json:"flows,omitempty" yaml:"flows,omitempty"`

// OpenId Connect URL to discover OAuth2 configuration values.

// This MUST be in the form of a URL. The OpenID Connect standard requires the use of TLS.

OpenIDConnectURL

string

`json:"openIdConnectUrl,omitempty" yaml:"openIdConnectUrl,omitempty"`

Common

OpenAPICommon

`json:"-" yaml:",inline"`

}

SecurityScheme defines a security scheme that can be used by the operations.

See

https://spec.openapis.org/oas/v3.1.0#security-scheme-object

.

type

Server

¶

type Server struct {

// REQUIRED. A URL to the target host. This URL supports Server Variables and MAY be relative,

// to indicate that the host location is relative to the location where the OpenAPI document is being served.

// Variable substitutions will be made when a variable is named in {brackets}.

URL

string

`json:"url" yaml:"url"`

// An optional string describing the host designated by the URL.

// CommonMark syntax MAY be used for rich text representation.

Description

string

`json:"description,omitempty" yaml:"description,omitempty"`

// A map between a variable name and its value. The value is used for substitution in the server's URL template.

Variables map[

string

]

ServerVariable

`json:"variables,omitempty" yaml:"variables,omitempty"`

Common

OpenAPICommon

`json:"-" yaml:",inline"`

}

Server represents a Server.

See

https://spec.openapis.org/oas/v3.1.0#server-object

.

func

NewServer

¶

func NewServer() *

Server

NewServer returns a new Server.

func (*Server)

SetDescription

¶

func (s *

Server

) SetDescription(d

string

) *

Server

SetDescription sets the Description of the Server.

func (*Server)

SetURL

¶

func (s *

Server

) SetURL(url

string

) *

Server

SetURL sets the URL of the Server.

type

ServerVariable

¶

added in

v0.40.0

type ServerVariable struct {

// An enumeration of string values to be used if the substitution options are from a limited set.

//

// The array MUST NOT be empty.

Enum []

string

`json:"enum,omitempty" yaml:"enum,omitempty"`

// REQUIRED. The default value to use for substitution, which SHALL be sent if an alternate value is not supplied.

// Note this behavior is different than the Schema Object's treatment of default values, because in those

// cases parameter values are optional. If the enum is defined, the value MUST exist in the enum's values.

Default

string

`json:"default" yaml:"default"`

// An optional description for the server variable. CommonMark syntax MAY be used for rich text representation.

Description

string

`json:"description,omitempty" yaml:"description,omitempty"`

Common

OpenAPICommon

`json:"-" yaml:",inline"`

}

ServerVariable describes an object representing a Server Variable for server URL template substitution.

See

https://spec.openapis.org/oas/v3.1.0#server-variable-object

type

Spec

¶

type Spec struct {

// REQUIRED. This string MUST be the version number of the OpenAPI Specification

// that the OpenAPI document uses.

OpenAPI

string

`json:"openapi" yaml:"openapi"`

// Added just to detect v2 openAPI specifications and to pretty print version error.

Swagger

string

`json:"swagger,omitempty" yaml:"swagger,omitempty"`

// REQUIRED. Provides metadata about the API.

//

// The metadata MAY be used by tooling as required.

Info

Info

`json:"info" yaml:"info"`

// The default value for the `$schema` keyword within Schema Objects contained within this OAS document.

JSONSchemaDialect

string

`json:"jsonSchemaDialect,omitempty" yaml:"jsonSchemaDialect,omitempty"`

// An array of Server Objects, which provide connectivity information to a target server.

Servers []

Server

`json:"servers,omitempty" yaml:"servers,omitempty"`

// The available paths and operations for the API.

Paths

Paths

`json:"paths,omitempty" yaml:"paths,omitempty"`

// The incoming webhooks that MAY be received as part of this API and that

// the API consumer MAY choose to implement.

//

// Closely related to the `callbacks` feature, this section describes requests initiated other

// than by an API call, for example by an out of band registration.

//

// The key name is a unique string to refer to each webhook, while the (optionally referenced)

// PathItem Object describes a request that may be initiated by the API provider and the expected responses.

Webhooks map[

string

]*

PathItem

`json:"webhooks,omitempty" yaml:"webhooks,omitempty"`

// An element to hold various schemas for the document.

Components *

Components

`json:"components,omitempty" yaml:"components,omitempty"`

// A declaration of which security mechanisms can be used across the API.

// The list of values includes alternative security requirement objects that can be used.

//

// Only one of the security requirement objects need to be satisfied to authorize a request.

//

// Individual operations can override this definition.

Security

SecurityRequirements

`json:"security,omitempty" yaml:"security,omitempty"`

// A list of tags used by the specification with additional metadata.

// The order of the tags can be used to reflect on their order by the parsing

// tools. Not all tags that are used by the Operation Object must be declared.

// The tags that are not declared MAY be organized randomly or based on the tools' logic.

// Each tag name in the list MUST be unique.

Tags []

Tag

`json:"tags,omitempty" yaml:"tags,omitempty"`

// Additional external documentation.

ExternalDocs *

ExternalDocumentation

`json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`

// Specification extensions.

Extensions

Extensions

`json:"-" yaml:",inline"`

// Raw YAML node. Used by '$ref' resolvers.

Raw *

yaml

.

Node

`json:"-" yaml:"-"`
}

Spec is the root document object of the OpenAPI document.

See

https://spec.openapis.org/oas/v3.1.0#openapi-object

.

func

NewSpec

¶

func NewSpec() *

Spec

NewSpec returns a new Spec.

func

Parse

¶

func Parse(data []

byte

) (s *

Spec

, err

error

)

Parse parses JSON/YAML into OpenAPI Spec.

func (*Spec)

AddNamedParameters

¶

func (s *

Spec

) AddNamedParameters(ps ...*

NamedParameter

) *

Spec

AddNamedParameters adds the given namedParameters to the Components of the Spec.

func (*Spec)

AddNamedPathItems

¶

func (s *

Spec

) AddNamedPathItems(ps ...*

NamedPathItem

) *

Spec

AddNamedPathItems adds the given namedPaths to the Paths of the Spec.

func (*Spec)

AddNamedRequestBodies

¶

func (s *

Spec

) AddNamedRequestBodies(scs ...*

NamedRequestBody

) *

Spec

AddNamedRequestBodies adds the given namedRequestBodies to the Components of the Spec.

func (*Spec)

AddNamedResponses

¶

func (s *

Spec

) AddNamedResponses(scs ...*

NamedResponse

) *

Spec

AddNamedResponses adds the given namedResponses to the Components of the Spec.

func (*Spec)

AddNamedSchemas

¶

func (s *

Spec

) AddNamedSchemas(scs ...*

NamedSchema

) *

Spec

AddNamedSchemas adds the given namedSchemas to the Components of the Spec.

func (*Spec)

AddParameter

¶

func (s *

Spec

) AddParameter(n

string

, p *

Parameter

) *

Spec

AddParameter adds the given Parameter under the given Name to the Components of the Spec.

func (*Spec)

AddPathItem

¶

func (s *

Spec

) AddPathItem(n

string

, p *

PathItem

) *

Spec

AddPathItem adds the given PathItem under the given Name to the Paths of the Spec.

func (*Spec)

AddRequestBody

¶

func (s *

Spec

) AddRequestBody(n

string

, sc *

RequestBody

) *

Spec

AddRequestBody adds the given RequestBody under the given Name to the Components of the Spec.

func (*Spec)

AddResponse

¶

func (s *

Spec

) AddResponse(n

string

, sc *

Response

) *

Spec

AddResponse adds the given Response under the given Name to the Components of the Spec.

func (*Spec)

AddSchema

¶

func (s *

Spec

) AddSchema(n

string

, sc *

Schema

) *

Spec

AddSchema adds the given Schema under the given Name to the Components of the Spec.

func (*Spec)

AddServers

¶

func (s *

Spec

) AddServers(srvs ...*

Server

) *

Spec

AddServers adds Servers to the Servers of the Spec.

func (*Spec)

Init

¶

func (s *

Spec

) Init()

Init components of schema.

func (*Spec)

RefRequestBody

¶

func (s *

Spec

) RefRequestBody(n

string

) *

NamedRequestBody

RefRequestBody returns a new RequestBody referencing the given name.

func (*Spec)

RefResponse

¶

func (s *

Spec

) RefResponse(n

string

) *

NamedResponse

RefResponse returns a new Response referencing the given name.

func (*Spec)

RefSchema

¶

func (s *

Spec

) RefSchema(n

string

) *

NamedSchema

RefSchema returns a new Schema referencing the given name.

func (*Spec)

SetComponents

¶

func (s *

Spec

) SetComponents(c *

Components

) *

Spec

SetComponents sets the Components of the Spec.

func (*Spec)

SetInfo

¶

func (s *

Spec

) SetInfo(i *

Info

) *

Spec

SetInfo sets the Info of the Spec.

func (*Spec)

SetOpenAPI

¶

func (s *

Spec

) SetOpenAPI(v

string

) *

Spec

SetOpenAPI sets the OpenAPI Specification version of the document.

func (*Spec)

SetPaths

¶

func (s *

Spec

) SetPaths(p

Paths

) *

Spec

SetPaths sets the Paths of the Spec.

func (*Spec)

SetServers

¶

func (s *

Spec

) SetServers(srvs []

Server

) *

Spec

SetServers sets the Servers of the Spec.

func (*Spec)

UnmarshalJSON

¶

added in

v0.22.0

func (s *

Spec

) UnmarshalJSON(bytes []

byte

)

error

UnmarshalJSON implements json.Unmarshaler.

func (*Spec)

UnmarshalYAML

¶

added in

v0.44.0

func (s *

Spec

) UnmarshalYAML(n *

yaml

.

Node

)

error

UnmarshalYAML implements yaml.Unmarshaler.

type

Tag

¶

type Tag struct {

// REQUIRED. The name of the tag.

Name

string

`json:"name" yaml:"name"`

// A description for the tag. CommonMark syntax MAY be used for rich text representation.

Description

string

`json:"description,omitempty" yaml:"description,omitempty"`

// Additional external documentation for this tag.

ExternalDocs *

ExternalDocumentation

`json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`

// Specification extensions.

Extensions

Extensions

`json:"-" yaml:",inline"`
}

Tag adds metadata to a single tag that is used by the Operation Object.

See

https://spec.openapis.org/oas/v3.1.0#tag-object

type

XML

¶

added in

v0.44.0

type XML struct {

// Replaces the name of the element/attribute used for the described schema property.

//

// When defined within items, it will affect the name of the individual XML elements within the list.

//

// When defined alongside type being array (outside the items), it will affect the wrapping element

// and only if wrapped is true.

//

// If wrapped is false, it will be ignored.

Name

string

`json:"name,omitempty" yaml:"name,omitempty"`

// The URI of the namespace definition.

//

// This MUST be in the form of an absolute URI.

Namespace

string

`json:"namespace,omitempty" yaml:"namespace,omitempty"`

// The prefix to be used for the name.

Prefix

string

`json:"prefix,omitempty" yaml:"prefix,omitempty"`

// Declares whether the property definition translates to an attribute instead of an element.

//

// Default value is false.

Attribute

bool

`json:"attribute,omitempty" yaml:"attribute,omitempty"`

// MAY be used only for an array definition. Signifies whether the array is wrapped

// (for example, `<books><book/><book/></books>`) or unwrapped (`<book/><book/>`).

//

// The definition takes effect only when defined alongside type being array (outside the items).

//

// Default value is false.

Wrapped

bool

`json:"wrapped,omitempty" yaml:"wrapped,omitempty"`

Common

OpenAPICommon

`json:"-" yaml:",inline"`

}

XML is a metadata object that allows for more fine-tuned XML model definitions.

See

https://spec.openapis.org/oas/v3.1.0#xml-object

.

func (*XML)

ToJSONSchema

¶

added in

v0.44.0

func (d *

XML

) ToJSONSchema() *

jsonschema

.

XML

ToJSONSchema converts XML to jsonschema.XML.