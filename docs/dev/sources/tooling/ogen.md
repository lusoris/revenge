# ogen OpenAPI Generator

> Source: https://pkg.go.dev/github.com/ogen-go/ogen
> Fetched: 2026-02-01T11:41:51.863618+00:00
> Content-Hash: abfde520a6beeadb
> Type: html

---

### Overview ¶

Package ogen implements OpenAPI v3 code generation.

### Index ¶

- type AdditionalProperties
-     * func (p AdditionalProperties) MarshalJSON() ([]byte, error)
  - func (p AdditionalProperties) MarshalYAML() (any, error)
  - func (p *AdditionalProperties) ToJSONSchema()*jsonschema.AdditionalProperties
  - func (p *AdditionalProperties) UnmarshalJSON(data []byte) error
  - func (p *AdditionalProperties) UnmarshalYAML(node*yaml.Node) error
- type Callback
- type Components
-     * func (c *Components) Init()
- type Contact
-     * func NewContact() *Contact
-     * func (c *Contact) SetEmail(e string) *Contact
  - func (c *Contact) SetName(n string)*Contact
  - func (c *Contact) SetURL(url string)*Contact
- type Default
- type Discriminator
-     * func (d *Discriminator) ToJSONSchema() *jsonschema.RawDiscriminator
- type Encoding
- type Enum
- type Example
- type ExampleValue
- type Extensions
- type ExternalDocumentation
- type Header
- type Info
-     * func NewInfo() *Info
-     * func (i *Info) SetContact(c *Contact) *Info
  - func (i *Info) SetDescription(d string)*Info
  - func (i *Info) SetLicense(l*License) *Info
  - func (i *Info) SetTermsOfService(t string)*Info
  - func (i *Info) SetTitle(t string)*Info
  - func (i *Info) SetVersion(v string)*Info
- type Items
-     * func (p Items) MarshalJSON() ([]byte, error)
  - func (p Items) MarshalYAML() (any, error)
  - func (p *Items) ToJSONSchema()*jsonschema.RawItems
  - func (p *Items) UnmarshalJSON(data []byte) error
  - func (p *Items) UnmarshalYAML(node*yaml.Node) error
- type License
-     * func NewLicense() *License
-     * func (l *License) SetName(n string) *License
  - func (l *License) SetURL(url string)*License
- type Link
- type Locator
- type Media
- type NamedParameter
-     * func NewNamedParameter(n string, p *Parameter) *NamedParameter
-     * func (p *NamedParameter) AsLocalRef() *Parameter
- type NamedPathItem
-     * func NewNamedPathItem(n string, p *PathItem) *NamedPathItem
-     * func (p *NamedPathItem) AsLocalRef() *PathItem
- type NamedRequestBody
-     * func NewNamedRequestBody(n string, p *RequestBody) *NamedRequestBody
-     * func (p *NamedRequestBody) AsLocalRef() *RequestBody
- type NamedResponse
-     * func NewNamedResponse(n string, p *Response) *NamedResponse
-     * func (p *NamedResponse) AsLocalRef() *Response
- type NamedSchema
-     * func NewNamedSchema(n string, p *Schema) *NamedSchema
-     * func (p *NamedSchema) AsLocalRef() *Schema
- type Num
- type OAuthFlow
- type OAuthFlows
- type OpenAPICommon
- type Operation
-     * func NewOperation() *Operation
-     * func (o *Operation) AddNamedResponses(ps ...*NamedResponse) *Operation
  - func (o *Operation) AddParameters(ps ...*Parameter) *Operation
  - func (o *Operation) AddResponse(n string, p*Response) *Operation
  - func (o *Operation) AddTags(ts ...string)*Operation
  - func (s *Operation) MarshalJSON() ([]byte, error)
  - func (o *Operation) SetDescription(d string)*Operation
  - func (o *Operation) SetOperationID(id string)*Operation
  - func (o *Operation) SetParameters(ps []*Parameter) *Operation
  - func (o *Operation) SetRequestBody(r*RequestBody) *Operation
  - func (o *Operation) SetResponses(r Responses)*Operation
  - func (o *Operation) SetSummary(s string)*Operation
  - func (o *Operation) SetTags(ts []string)*Operation
- type Parameter
-     * func NewParameter() *Parameter
-     * func (p *Parameter) InCookie() *Parameter
  - func (p *Parameter) InHeader()*Parameter
  - func (p *Parameter) InPath()*Parameter
  - func (p *Parameter) InQuery()*Parameter
  - func (p *Parameter) SetContent(c map[string]Media)*Parameter
  - func (p *Parameter) SetDeprecated(d bool)*Parameter
  - func (p *Parameter) SetDescription(d string)*Parameter
  - func (p *Parameter) SetExplode(e bool)*Parameter
  - func (p *Parameter) SetIn(i string)*Parameter
  - func (p *Parameter) SetName(n string)*Parameter
  - func (p *Parameter) SetRef(r string)*Parameter
  - func (p *Parameter) SetRequired(r bool)*Parameter
  - func (p *Parameter) SetSchema(s*Schema) *Parameter
  - func (p *Parameter) SetStyle(s string)*Parameter
  - func (p *Parameter) ToNamed(n string)*NamedParameter
- type PathItem
-     * func NewPathItem() *PathItem
-     * func (p *PathItem) AddParameters(ps ...*Parameter) *PathItem
  - func (p *PathItem) AddServers(srvs ...*Server) *PathItem
  - func (s *PathItem) MarshalJSON() ([]byte, error)
  - func (p *PathItem) SetDelete(o*Operation) *PathItem
  - func (p *PathItem) SetDescription(d string)*PathItem
  - func (p *PathItem) SetGet(o*Operation) *PathItem
  - func (p *PathItem) SetHead(o*Operation) *PathItem
  - func (p *PathItem) SetOptions(o*Operation) *PathItem
  - func (p *PathItem) SetParameters(ps []*Parameter) *PathItem
  - func (p *PathItem) SetPatch(o*Operation) *PathItem
  - func (p *PathItem) SetPost(o*Operation) *PathItem
  - func (p *PathItem) SetPut(o*Operation) *PathItem
  - func (p *PathItem) SetRef(r string)*PathItem
  - func (p *PathItem) SetServers(srvs []Server)*PathItem
  - func (p *PathItem) SetTrace(o*Operation) *PathItem
  - func (p *PathItem) ToNamed(n string)*NamedPathItem
- type Paths
- type PatternProperties
-     * func (p PatternProperties) MarshalJSON() ([]byte, error)
  - func (p PatternProperties) MarshalYAML() (any, error)
  - func (p PatternProperties) ToJSONSchema() (result jsonschema.RawPatternProperties)
  - func (p *PatternProperties) UnmarshalJSON(data []byte) error
  - func (p *PatternProperties) UnmarshalYAML(node*yaml.Node) error
- type PatternProperty
- type Properties
-     * func (p Properties) MarshalJSON() ([]byte, error)
  - func (p Properties) MarshalYAML() (any, error)
  - func (p Properties) ToJSONSchema() jsonschema.RawProperties
  - func (p *Properties) UnmarshalJSON(data []byte) error
  - func (p *Properties) UnmarshalYAML(node*yaml.Node) error
- type Property
-     * func NewProperty() *Property
-     * func (p *Property) SetName(n string) *Property
  - func (p *Property) SetSchema(s*Schema) *Property
  - func (p Property) ToJSONSchema() jsonschema.RawProperty
- type RawValue
- type RequestBody
-     * func NewRequestBody() *RequestBody
-     * func (r *RequestBody) AddContent(mt string, s *Schema) *RequestBody
  - func (r *RequestBody) SetContent(c map[string]Media)*RequestBody
  - func (r *RequestBody) SetDescription(d string)*RequestBody
  - func (r *RequestBody) SetJSONContent(s*Schema) *RequestBody
  - func (r *RequestBody) SetRef(ref string)*RequestBody
  - func (r *RequestBody) SetRequired(req bool)*RequestBody
  - func (r *RequestBody) ToNamed(n string)*NamedRequestBody
- type Response
-     * func NewResponse() *Response
-     * func (r *Response) AddContent(mt string, s *Schema) *Response
  - func (r *Response) SetContent(c map[string]Media)*Response
  - func (r *Response) SetDescription(d string)*Response
  - func (r *Response) SetHeaders(h map[string]*Header) *Response
  - func (r *Response) SetJSONContent(s*Schema) *Response
  - func (r *Response) SetLinks(l map[string]*Link) *Response
  - func (r *Response) SetRef(ref string)*Response
  - func (r *Response) ToNamed(n string)*NamedResponse
- type Responses
- type Schema
-     * func Binary() *Schema
  - func Bool() *Schema
  - func Bytes() *Schema
  - func Date() *Schema
  - func DateTime() *Schema
  - func Double() *Schema
  - func Float() *Schema
  - func Int() *Schema
  - func Int32() *Schema
  - func Int64() *Schema
  - func NewSchema() *Schema
  - func Password() *Schema
  - func String() *Schema
  - func UUID() *Schema
-     * func (s *Schema) AddOptionalProperties(ps ...*Property) *Schema
  - func (s *Schema) AddRequiredProperties(ps ...*Property) *Schema
  - func (s *Schema) AsArray()*Schema
  - func (s *Schema) AsEnum(def json.RawMessage, values ...json.RawMessage)*Schema
  - func (s *Schema) SetAllOf(a []*Schema) *Schema
  - func (s *Schema) SetAnyOf(a []*Schema) *Schema
  - func (s *Schema) SetDefault(d json.RawMessage)*Schema
  - func (s *Schema) SetDeprecated(d bool)*Schema
  - func (s *Schema) SetDescription(d string)*Schema
  - func (s *Schema) SetDiscriminator(d*Discriminator) *Schema
  - func (s *Schema) SetEnum(e []json.RawMessage)*Schema
  - func (s *Schema) SetExclusiveMaximum(e bool)*Schema
  - func (s *Schema) SetExclusiveMinimum(e bool)*Schema
  - func (s *Schema) SetFormat(f string)*Schema
  - func (s *Schema) SetItems(i*Schema) *Schema
  - func (s *Schema) SetMaxItems(m*uint64) *Schema
  - func (s *Schema) SetMaxLength(m*uint64) *Schema
  - func (s *Schema) SetMaxProperties(m*uint64) *Schema
  - func (s *Schema) SetMaximum(m*int64) *Schema
  - func (s *Schema) SetMinItems(m*uint64) *Schema
  - func (s *Schema) SetMinLength(m*uint64) *Schema
  - func (s *Schema) SetMinProperties(m*uint64) *Schema
  - func (s *Schema) SetMinimum(m*int64) *Schema
  - func (s *Schema) SetMultipleOf(m*uint64) *Schema
  - func (s *Schema) SetNullable(n bool)*Schema
  - func (s *Schema) SetOneOf(o []*Schema) *Schema
  - func (s *Schema) SetPattern(p string)*Schema
  - func (s *Schema) SetProperties(p*Properties) *Schema
  - func (s *Schema) SetRef(r string)*Schema
  - func (s *Schema) SetRequired(r []string)*Schema
  - func (s *Schema) SetSummary(smry string)*Schema
  - func (s *Schema) SetType(t string)*Schema
  - func (s *Schema) SetUniqueItems(u bool)*Schema
  - func (s *Schema) ToJSONSchema()*jsonschema.RawSchema
  - func (s *Schema) ToNamed(n string)*NamedSchema
  - func (s *Schema) ToProperty(n string)*Property
- type SecurityRequirement
- type SecurityRequirements
- type SecurityScheme
- type Server
-     * func NewServer() *Server
-     * func (s *Server) SetDescription(d string) *Server
  - func (s *Server) SetURL(url string)*Server
- type ServerVariable
- type Spec
-     * func NewSpec() *Spec
  - func Parse(data []byte) (s *Spec, err error)
-     * func (s *Spec) AddNamedParameters(ps ...*NamedParameter) *Spec
  - func (s *Spec) AddNamedPathItems(ps ...*NamedPathItem) *Spec
  - func (s *Spec) AddNamedRequestBodies(scs ...*NamedRequestBody) *Spec
  - func (s *Spec) AddNamedResponses(scs ...*NamedResponse) *Spec
  - func (s *Spec) AddNamedSchemas(scs ...*NamedSchema) *Spec
  - func (s *Spec) AddParameter(n string, p*Parameter) *Spec
  - func (s *Spec) AddPathItem(n string, p*PathItem) *Spec
  - func (s *Spec) AddRequestBody(n string, sc*RequestBody) *Spec
  - func (s *Spec) AddResponse(n string, sc*Response) *Spec
  - func (s *Spec) AddSchema(n string, sc*Schema) *Spec
  - func (s *Spec) AddServers(srvs ...*Server) *Spec
  - func (s *Spec) Init()
  - func (s *Spec) RefRequestBody(n string)*NamedRequestBody
  - func (s *Spec) RefResponse(n string)*NamedResponse
  - func (s *Spec) RefSchema(n string)*NamedSchema
  - func (s *Spec) SetComponents(c*Components) *Spec
  - func (s *Spec) SetInfo(i*Info) *Spec
  - func (s *Spec) SetOpenAPI(v string)*Spec
  - func (s *Spec) SetPaths(p Paths)*Spec
  - func (s *Spec) SetServers(srvs []Server)*Spec
  - func (s *Spec) UnmarshalJSON(bytes []byte) error
  - func (s *Spec) UnmarshalYAML(n*yaml.Node) error
- type Tag
- type XML
-     * func (d *XML) ToJSONSchema() *jsonschema.XML

### Constants ¶

This section is empty.

### Variables ¶

This section is empty.

### Functions ¶

This section is empty.

### Types ¶

#### type [AdditionalProperties](https://github.com/ogen-go/ogen/blob/v1.18.0/schema.go#L342) ¶ added in v0.9.0

    type AdditionalProperties struct {
     Bool   *[bool](/builtin#bool)
     Schema Schema
    }

AdditionalProperties is JSON Schema additionalProperties validator description.

#### func (AdditionalProperties) [MarshalJSON](https://github.com/ogen-go/ogen/blob/v1.18.0/schema.go#L368) ¶ added in v0.9.0

    func (p AdditionalProperties) MarshalJSON() ([][byte](/builtin#byte), [error](/builtin#error))

MarshalJSON implements json.Marshaler.

#### func (AdditionalProperties) [MarshalYAML](https://github.com/ogen-go/ogen/blob/v1.18.0/schema.go#L348) ¶ added in v0.44.0

    func (p AdditionalProperties) MarshalYAML() ([any](/builtin#any), [error](/builtin#error))

MarshalYAML implements yaml.Marshaler.

#### func (*AdditionalProperties) [ToJSONSchema](https://github.com/ogen-go/ogen/blob/v1.18.0/schema_backcomp.go#L79) ¶ added in v0.13.0

    func (p *AdditionalProperties) ToJSONSchema() *[jsonschema](/github.com/ogen-go/ogen@v1.18.0/jsonschema).[AdditionalProperties](/github.com/ogen-go/ogen@v1.18.0/jsonschema#AdditionalProperties)

ToJSONSchema converts AdditionalProperties to jsonschema.AdditionalProperties.

#### func (*AdditionalProperties) [UnmarshalJSON](https://github.com/ogen-go/ogen/blob/v1.18.0/schema.go#L376) ¶ added in v0.9.0

    func (p *AdditionalProperties) UnmarshalJSON(data [][byte](/builtin#byte)) [error](/builtin#error)

UnmarshalJSON implements json.Unmarshaler.

#### func (*AdditionalProperties) [UnmarshalYAML](https://github.com/ogen-go/ogen/blob/v1.18.0/schema.go#L356) ¶ added in v0.43.0

    func (p *AdditionalProperties) UnmarshalYAML(node *[yaml](/github.com/go-faster/yaml).[Node](/github.com/go-faster/yaml#Node)) [error](/builtin#error)

UnmarshalYAML implements yaml.Unmarshaler.

#### type [Callback](https://github.com/ogen-go/ogen/blob/v1.18.0/spec.go#L658) ¶ added in v0.44.0

    type Callback map[[string](/builtin#string)]*PathItem

Callback is a map of possible out-of band callbacks related to the parent operation.

Each value in the map is a Path Item Object that describes a set of requests that may be initiated by the API provider and the expected responses.

The key value used to identify the path item object is an expression, evaluated at runtime, that identifies a URL to use for the callback operation.

To describe incoming requests from the API provider independent from another API call, use the `webhooks` field.

See <https://spec.openapis.org/oas/v3.1.0#callback-object>.

#### type [Components](https://github.com/ogen-go/ogen/blob/v1.18.0/spec.go#L229) ¶

    type Components struct {
     // An object to hold reusable Schema Objects.
     Schemas map[[string](/builtin#string)]*Schema `json:"schemas,omitempty" yaml:"schemas,omitempty"`
     // An object to hold reusable Response Objects.
     Responses map[[string](/builtin#string)]*Response `json:"responses,omitempty" yaml:"responses,omitempty"`
     // An object to hold reusable Parameter Objects.
     Parameters map[[string](/builtin#string)]*Parameter `json:"parameters,omitempty" yaml:"parameters,omitempty"`
     // An object to hold reusable Example Objects.
     Examples map[[string](/builtin#string)]*Example `json:"examples,omitempty" yaml:"examples,omitempty"`
     // An object to hold reusable Request Body Objects.
     RequestBodies map[[string](/builtin#string)]*RequestBody `json:"requestBodies,omitempty" yaml:"requestBodies,omitempty"`
     // An object to hold reusable Header Objects.
     Headers map[[string](/builtin#string)]*Header `json:"headers,omitempty" yaml:"headers,omitempty"`
     // An object to hold reusable Security Scheme Objects.
     SecuritySchemes map[[string](/builtin#string)]*SecurityScheme `json:"securitySchemes,omitempty" yaml:"securitySchemes,omitempty"`
     // An object to hold reusable Link Objects.
     Links map[[string](/builtin#string)]*Link `json:"links,omitempty" yaml:"links,omitempty"`
     // An object to hold reusable Callback Objects.
     Callbacks map[[string](/builtin#string)]*Callback `json:"callbacks,omitempty" yaml:"callbacks,omitempty"`
     // An object to hold reusable Path Item Objects.
     PathItems map[[string](/builtin#string)]*PathItem `json:"pathItems,omitempty" yaml:"pathItems,omitempty"`
    
     Common OpenAPICommon `json:"-" yaml:",inline"`
    }

Components Holds a set of reusable objects for different aspects of the OAS. All objects defined within the components object will have no effect on the API unless they are explicitly referenced from properties outside the components object.

See <https://spec.openapis.org/oas/v3.1.0#components-object>.

#### func (*Components) [Init](https://github.com/ogen-go/ogen/blob/v1.18.0/spec.go#L262) ¶ added in v0.42.0

    func (c *Components) Init()

Init initializes all fields.

#### type [Contact](https://github.com/ogen-go/ogen/blob/v1.18.0/spec.go#L162) ¶

    type Contact struct {
     // The identifying name of the contact person/organization.
     Name [string](/builtin#string) `json:"name,omitempty" yaml:"name,omitempty"`
     // The URL pointing to the contact information.
     URL [string](/builtin#string) `json:"url,omitempty" yaml:"url,omitempty"`
     // The email address of the contact person/organization.
     Email [string](/builtin#string) `json:"email,omitempty" yaml:"email,omitempty"`
    
     // Specification extensions.
     Extensions Extensions `json:"-" yaml:",inline"`
    }

Contact information for the exposed API.

See <https://spec.openapis.org/oas/v3.1.0#contact-object>.

#### func [NewContact](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L325) ¶

    func NewContact() *Contact

NewContact returns a new Contact.

#### func (*Contact) [SetEmail](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L342) ¶

    func (c *Contact) SetEmail(e [string](/builtin#string)) *Contact

SetEmail sets the Email of the Contact.

#### func (*Contact) [SetName](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L330) ¶

    func (c *Contact) SetName(n [string](/builtin#string)) *Contact

SetName sets the Name of the Contact.

#### func (*Contact) [SetURL](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L336) ¶

    func (c *Contact) SetURL(url [string](/builtin#string)) *Contact

SetURL sets the URL of the Contact.

#### type [Default](https://github.com/ogen-go/ogen/blob/v1.18.0/spec.go#L19) ¶ added in v0.43.0

    type Default = [jsonschema](/github.com/ogen-go/ogen@v1.18.0/jsonschema).[Default](/github.com/ogen-go/ogen@v1.18.0/jsonschema#Default)

Default is a default value.

#### type [Discriminator](https://github.com/ogen-go/ogen/blob/v1.18.0/schema.go#L556) ¶

    type Discriminator struct {
     // REQUIRED. The name of the property in the payload that will hold the discriminator value.
     PropertyName [string](/builtin#string) `json:"propertyName" yaml:"propertyName"`
     // An object to hold mappings between payload values and schema names or references.
     Mapping map[[string](/builtin#string)][string](/builtin#string) `json:"mapping,omitempty" yaml:"mapping,omitempty"`
    
     Common OpenAPICommon `json:"-" yaml:",inline"`
    }

Discriminator discriminates types for OneOf, AllOf, AnyOf.

See <https://spec.openapis.org/oas/v3.1.0#discriminator-object>.

#### func (*Discriminator) [ToJSONSchema](https://github.com/ogen-go/ogen/blob/v1.18.0/schema_backcomp.go#L127) ¶ added in v0.13.0

    func (d *Discriminator) ToJSONSchema() *[jsonschema](/github.com/ogen-go/ogen@v1.18.0/jsonschema).[RawDiscriminator](/github.com/ogen-go/ogen@v1.18.0/jsonschema#RawDiscriminator)

ToJSONSchema converts Discriminator to jsonschema.RawDiscriminator.

#### type [Encoding](https://github.com/ogen-go/ogen/blob/v1.18.0/spec.go#L573) ¶ added in v0.33.0

    type Encoding struct {
     // The Content-Type for encoding a specific property.
     ContentType [string](/builtin#string) `json:"contentType,omitempty" yaml:"contentType,omitempty"`
    
     // A map allowing additional information to be provided as headers, for example Content-Disposition.
     // Content-Type is described separately and SHALL be ignored in this section. This property SHALL be
     // ignored if the request body media type is not a multipart.
     Headers map[[string](/builtin#string)]*Header `json:"headers,omitempty" yaml:"headers,omitempty"`
    
     // Describes how the parameter value will be serialized
     // depending on the type of the parameter value.
     Style [string](/builtin#string) `json:"style,omitempty" yaml:"style,omitempty"`
    
     // When this is true, parameter values of type array or object
     // generate separate parameters for each value of the array
     // or key-value pair of the map.
     // For other types of parameters this property has no effect.
     Explode *[bool](/builtin#bool) `json:"explode,omitempty" yaml:"explode,omitempty"`
    
     // Determines whether the parameter value SHOULD allow reserved characters, as defined by
     // RFC3986 :/?#[]@!$&'()*+,;= to be included without percent-encoding.
     // The default value is false. This property SHALL be ignored if the request body media type
     // is not application/x-www-form-urlencoded.
     AllowReserved [bool](/builtin#bool) `json:"allowReserved,omitempty" yaml:"allowReserved,omitempty"`
    
     Common OpenAPICommon `json:"-" yaml:",inline"`
    }

Encoding describes single encoding definition applied to a single schema property.

See <https://spec.openapis.org/oas/v3.1.0#encoding-object>.

#### type [Enum](https://github.com/ogen-go/ogen/blob/v1.18.0/spec.go#L17) ¶ added in v0.38.0

    type Enum = [jsonschema](/github.com/ogen-go/ogen@v1.18.0/jsonschema).[Enum](/github.com/ogen-go/ogen@v1.18.0/jsonschema#Enum)

Enum is JSON Schema enum validator description.

#### type [Example](https://github.com/ogen-go/ogen/blob/v1.18.0/spec.go#L663) ¶

    type Example struct {
     Ref [string](/builtin#string) `json:"$ref,omitempty" yaml:"$ref,omitempty"` // ref object
    
     // Short description for the example.
     Summary [string](/builtin#string) `json:"summary,omitempty" yaml:"summary,omitempty"`
     // Long description for the example.
     // CommonMark syntax MAY be used for rich text representation.
     Description [string](/builtin#string) `json:"description,omitempty" yaml:"description,omitempty"`
     // Embedded literal example.
     Value ExampleValue `json:"value,omitempty" yaml:"value,omitempty"`
     // A URI that points to the literal example.
     ExternalValue [string](/builtin#string) `json:"externalValue,omitempty" yaml:"externalValue,omitempty"`
    
     Common OpenAPICommon `json:"-" yaml:",inline"`
    }

Example object.

See <https://spec.openapis.org/oas/v3.1.0#example-object>.

#### type [ExampleValue](https://github.com/ogen-go/ogen/blob/v1.18.0/spec.go#L21) ¶ added in v0.43.0

    type ExampleValue = [jsonschema](/github.com/ogen-go/ogen@v1.18.0/jsonschema).[Example](/github.com/ogen-go/ogen@v1.18.0/jsonschema#Example)

ExampleValue is an example value.

#### type [Extensions](https://github.com/ogen-go/ogen/blob/v1.18.0/spec.go#L28) ¶ added in v0.49.0

    type Extensions = [jsonschema](/github.com/ogen-go/ogen@v1.18.0/jsonschema).[Extensions](/github.com/ogen-go/ogen@v1.18.0/jsonschema#Extensions)

Extensions is a map of OpenAPI extensions.

See <https://spec.openapis.org/oas/v3.1.0#specification-extensions>.

#### type [ExternalDocumentation](https://github.com/ogen-go/ogen/blob/v1.18.0/spec.go#L457) ¶ added in v0.40.0

    type ExternalDocumentation struct {
     // A description of the target documentation. CommonMark syntax MAY be used for rich text representation.
     Description [string](/builtin#string) `json:"description,omitempty" yaml:"description,omitempty"`
     // REQUIRED. The URL for the target documentation. This MUST be in the form of a URL.
     URL [string](/builtin#string) `json:"url" yaml:"url"`
    
     // Specification extensions.
     Extensions Extensions `json:"-" yaml:",inline"`
    }

ExternalDocumentation describes a reference to external resource for extended documentation.

See <https://spec.openapis.org/oas/v3.1.0#external-documentation-object>.

#### type [Header](https://github.com/ogen-go/ogen/blob/v1.18.0/spec.go#L712) ¶ added in v0.33.0

    type Header = Parameter

Header describes header response.

Header Object follows the structure of the Parameter Object with the following changes:

  1. `name` MUST NOT be specified, it is given in the corresponding headers map.
  2. `in` MUST NOT be specified, it is implicitly in header.
  3. All traits that are affected by the location MUST be applicable to a location of header.

See <https://spec.openapis.org/oas/v3.1.0#header-object>.

#### type [Info](https://github.com/ogen-go/ogen/blob/v1.18.0/spec.go#L138) ¶

    type Info struct {
     // REQUIRED. The title of the API.
     Title [string](/builtin#string) `json:"title" yaml:"title"`
     // A short summary of the API.
     Summary [string](/builtin#string) `json:"summary,omitempty" yaml:"summary,omitempty"`
     // A short description of the API.
     // CommonMark syntax MAY be used for rich text representation.
     Description [string](/builtin#string) `json:"description,omitempty" yaml:"description,omitempty"`
     // A URL to the Terms of Service for the API. MUST be in the format of a URL.
     TermsOfService [string](/builtin#string) `json:"termsOfService,omitempty" yaml:"termsOfService,omitempty"`
     // The contact information for the exposed API.
     Contact *Contact `json:"contact,omitempty" yaml:"contact,omitempty"`
     // The license information for the exposed API.
     License *License `json:"license,omitempty" yaml:"license,omitempty"`
     // REQUIRED. The version of the OpenAPI document.
     Version [string](/builtin#string) `json:"version" yaml:"version"`
    
     // Specification extensions.
     Extensions Extensions `json:"-" yaml:",inline"`
    }

Info provides metadata about the API.

The metadata MAY be used by the clients if needed, and MAY be presented in editing or documentation generation tools for convenience.

See <https://spec.openapis.org/oas/v3.1.0#info-object>.

#### func [NewInfo](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L284) ¶

    func NewInfo() *Info

NewInfo returns a new Info.

#### func (*Info) [SetContact](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L307) ¶

    func (i *Info) SetContact(c *Contact) *Info

SetContact sets the Contact of the Info.

#### func (*Info) [SetDescription](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L295) ¶

    func (i *Info) SetDescription(d [string](/builtin#string)) *Info

SetDescription sets the description of the Info.

#### func (*Info) [SetLicense](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L313) ¶

    func (i *Info) SetLicense(l *License) *Info

SetLicense sets the License of the Info.

#### func (*Info) [SetTermsOfService](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L301) ¶

    func (i *Info) SetTermsOfService(t [string](/builtin#string)) *Info

SetTermsOfService sets the terms of service of the Info.

#### func (*Info) [SetTitle](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L289) ¶

    func (i *Info) SetTitle(t [string](/builtin#string)) *Info

SetTitle sets the title of the Info.

#### func (*Info) [SetVersion](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L319) ¶

    func (i *Info) SetVersion(v [string](/builtin#string)) *Info

SetVersion sets the version of the Info.

#### type [Items](https://github.com/ogen-go/ogen/blob/v1.18.0/schema.go#L498) ¶ added in v0.69.0

    type Items struct {
     Item  *Schema
     Items []*Schema
    }

Items is unparsed JSON Schema items validator description.

#### func (Items) [MarshalJSON](https://github.com/ogen-go/ogen/blob/v1.18.0/schema.go#L524) ¶ added in v0.69.0

    func (p Items) MarshalJSON() ([][byte](/builtin#byte), [error](/builtin#error))

MarshalJSON implements json.Marshaler.

#### func (Items) [MarshalYAML](https://github.com/ogen-go/ogen/blob/v1.18.0/schema.go#L504) ¶ added in v0.69.0

    func (p Items) MarshalYAML() ([any](/builtin#any), [error](/builtin#error))

MarshalYAML implements yaml.Marshaler.

#### func (*Items) [ToJSONSchema](https://github.com/ogen-go/ogen/blob/v1.18.0/schema_backcomp.go#L107) ¶ added in v0.69.0

    func (p *Items) ToJSONSchema() *[jsonschema](/github.com/ogen-go/ogen@v1.18.0/jsonschema).[RawItems](/github.com/ogen-go/ogen@v1.18.0/jsonschema#RawItems)

ToJSONSchema converts Items to jsonschema.RawItems.

#### func (*Items) [UnmarshalJSON](https://github.com/ogen-go/ogen/blob/v1.18.0/schema.go#L532) ¶ added in v0.69.0

    func (p *Items) UnmarshalJSON(data [][byte](/builtin#byte)) [error](/builtin#error)

UnmarshalJSON implements json.Unmarshaler.

#### func (*Items) [UnmarshalYAML](https://github.com/ogen-go/ogen/blob/v1.18.0/schema.go#L512) ¶ added in v0.69.0

    func (p *Items) UnmarshalYAML(node *[yaml](/github.com/go-faster/yaml).[Node](/github.com/go-faster/yaml#Node)) [error](/builtin#error)

UnmarshalYAML implements yaml.Unmarshaler.

#### type [License](https://github.com/ogen-go/ogen/blob/v1.18.0/spec.go#L177) ¶

    type License struct {
     // REQUIRED. The license name used for the API.
     Name [string](/builtin#string) `json:"name" yaml:"name"`
     // An SPDX license expression for the API.
     Identifier [string](/builtin#string) `json:"identifier,omitempty" yaml:"identifier,omitempty"`
     // A URL to the license used for the API.
     URL [string](/builtin#string) `json:"url,omitempty" yaml:"url,omitempty"`
    
     // Specification extensions.
     Extensions Extensions `json:"-" yaml:",inline"`
    }

License information for the exposed API.

See <https://spec.openapis.org/oas/v3.1.0#license-object>.

#### func [NewLicense](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L348) ¶

    func NewLicense() *License

NewLicense returns a new License.

#### func (*License) [SetName](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L353) ¶

    func (l *License) SetName(n [string](/builtin#string)) *License

SetName sets the Name of the License.

#### func (*License) [SetURL](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L359) ¶

    func (l *License) SetURL(url [string](/builtin#string)) *License

SetURL sets the URL of the License.

#### type [Link](https://github.com/ogen-go/ogen/blob/v1.18.0/spec.go#L682) ¶ added in v0.44.0

    type Link struct {
     Ref [string](/builtin#string) `json:"$ref,omitempty" yaml:"$ref,omitempty"` // ref object
    
     // A relative or absolute URI reference to an OAS operation.
     //
     // This field is mutually exclusive of the operationId field, and MUST point to an Operation Object.
     OperationRef [string](/builtin#string) `json:"operationRef,omitempty" yaml:"operationRef,omitempty"`
     // The name of an existing, resolvable OAS operation, as defined with a unique operationId.
     //
     // This field is mutually exclusive of the operationRef field.
     OperationID [string](/builtin#string) `json:"operationId,omitempty" yaml:"operationId,omitempty"`
     // A map representing parameters to pass to an operation as specified with operationId or identified
     // via operationRef.
     //
     // The key is the parameter name to be used, whereas the value can be a constant or an expression to be
     // evaluated and passed to the linked operation.
     Parameters map[[string](/builtin#string)]RawValue `json:"parameters,omitempty" yaml:"parameters,omitempty"`
    
     Common OpenAPICommon `json:"-" yaml:",inline"`
    }

Link describes a possible design-time link for a response.

See <https://spec.openapis.org/oas/v3.1.0#link-object>.

#### type [Locator](https://github.com/ogen-go/ogen/blob/v1.18.0/spec.go#L33) ¶ added in v0.42.0

    type Locator = [location](/github.com/ogen-go/ogen@v1.18.0/location).[Locator](/github.com/ogen-go/ogen@v1.18.0/location#Locator)

Locator stores location of JSON value.

#### type [Media](https://github.com/ogen-go/ogen/blob/v1.18.0/spec.go#L551) ¶

    type Media struct {
     // The schema defining the content of the request, response, or parameter.
     Schema *Schema `json:"schema,omitempty" yaml:"schema,omitempty"`
     // Example of the media type.
     Example ExampleValue `json:"example,omitempty" yaml:"example,omitempty"`
     // Examples of the media type.
     Examples map[[string](/builtin#string)]*Example `json:"examples,omitempty" yaml:"examples,omitempty"`
    
     // A map between a property name and its encoding information.
     //
     // The key, being the property name, MUST exist in the schema as a property.
     //
     // The encoding object SHALL only apply to requestBody objects when the media
     // type is multipart or application/x-www-form-urlencoded.
     Encoding map[[string](/builtin#string)]Encoding `json:"encoding,omitempty" yaml:"encoding,omitempty"`
    
     Common OpenAPICommon `json:"-" yaml:",inline"`
    }

Media provides schema and examples for the media type identified by its key.

See <https://spec.openapis.org/oas/v3.1.0#media-type-object>.

#### type [NamedParameter](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L671) ¶

    type NamedParameter struct {
     Parameter *Parameter
     Name      [string](/builtin#string)
    }

NamedParameter can be used to construct a reference to the wrapped Parameter.

#### func [NewNamedParameter](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L677) ¶

    func NewNamedParameter(n [string](/builtin#string), p *Parameter) *NamedParameter

NewNamedParameter returns a new NamedParameter.

#### func (*NamedParameter) [AsLocalRef](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L682) ¶

    func (p *NamedParameter) AsLocalRef() *Parameter

AsLocalRef returns a new Parameter referencing the wrapped Parameter in the local document.

#### type [NamedPathItem](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L480) ¶

    type NamedPathItem struct {
     PathItem *PathItem
     Name     [string](/builtin#string)
    }

NamedPathItem can be used to construct a reference to the wrapped PathItem.

#### func [NewNamedPathItem](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L486) ¶

    func NewNamedPathItem(n [string](/builtin#string), p *PathItem) *NamedPathItem

NewNamedPathItem returns a new NamedPathItem.

#### func (*NamedPathItem) [AsLocalRef](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L491) ¶

    func (p *NamedPathItem) AsLocalRef() *PathItem

AsLocalRef returns a new PathItem referencing the wrapped PathItem in the local document.

#### type [NamedRequestBody](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L268) ¶

    type NamedRequestBody struct {
     RequestBody *RequestBody
     Name        [string](/builtin#string)
    }

NamedRequestBody can be used to construct a reference to the wrapped RequestBody.

#### func [NewNamedRequestBody](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L274) ¶

    func NewNamedRequestBody(n [string](/builtin#string), p *RequestBody) *NamedRequestBody

NewNamedRequestBody returns a new NamedRequestBody.

#### func (*NamedRequestBody) [AsLocalRef](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L279) ¶

    func (p *NamedRequestBody) AsLocalRef() *RequestBody

AsLocalRef returns a new RequestBody referencing the wrapped RequestBody in the local document.

#### type [NamedResponse](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L748) ¶

    type NamedResponse struct {
     Response *Response
     Name     [string](/builtin#string)
    }

NamedResponse can be used to construct a reference to the wrapped Response.

#### func [NewNamedResponse](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L754) ¶

    func NewNamedResponse(n [string](/builtin#string), p *Response) *NamedResponse

NewNamedResponse returns a new NamedResponse.

#### func (*NamedResponse) [AsLocalRef](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L759) ¶

    func (p *NamedResponse) AsLocalRef() *Response

AsLocalRef returns a new Response referencing the wrapped Response in the local document.

#### type [NamedSchema](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L1062) ¶

    type NamedSchema struct {
     Schema *Schema
     Name   [string](/builtin#string)
    }

NamedSchema can be used to construct a reference to the wrapped Schema.

#### func [NewNamedSchema](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L1068) ¶

    func NewNamedSchema(n [string](/builtin#string), p *Schema) *NamedSchema

NewNamedSchema returns a new NamedSchema.

#### func (*NamedSchema) [AsLocalRef](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L1073) ¶

    func (p *NamedSchema) AsLocalRef() *Schema

AsLocalRef returns a new Schema referencing the wrapped Schema in the local document.

#### type [Num](https://github.com/ogen-go/ogen/blob/v1.18.0/spec.go#L15) ¶ added in v0.16.0

    type Num = [jsonschema](/github.com/ogen-go/ogen@v1.18.0/jsonschema).[Num](/github.com/ogen-go/ogen@v1.18.0/jsonschema#Num)

Num represents JSON number.

#### type [OAuthFlow](https://github.com/ogen-go/ogen/blob/v1.18.0/security_scheme.go#L50) ¶ added in v0.19.0

    type OAuthFlow struct {
     // The authorization URL to be used for this flow.
     // This MUST be in the form of a URL. The OAuth2 standard requires the use of TLS.
     AuthorizationURL [string](/builtin#string) `json:"authorizationUrl,omitempty" yaml:"authorizationUrl,omitempty"`
     // The token URL to be used for this flow.
     // This MUST be in the form of a URL. The OAuth2 standard requires the use of TLS.
     TokenURL [string](/builtin#string) `json:"tokenUrl,omitempty" yaml:"tokenUrl,omitempty"`
     // The URL to be used for obtaining refresh tokens.
     // This MUST be in the form of a URL. The OAuth2 standard requires the use of TLS.
     RefreshURL [string](/builtin#string) `json:"refreshUrl,omitempty" yaml:"refreshUrl,omitempty"`
     // The available scopes for the OAuth2 security scheme.
     // A map between the scope name and a short description for it. The map MAY be empty.
     Scopes map[[string](/builtin#string)][string](/builtin#string) `json:"scopes" yaml:"scopes"`
    
     Common OpenAPICommon `json:"-" yaml:",inline"`
    }

OAuthFlow is configuration details for a supported OAuth Flow.

See <https://spec.openapis.org/oas/v3.1.0#oauth-flow-object>.

#### type [OAuthFlows](https://github.com/ogen-go/ogen/blob/v1.18.0/security_scheme.go#L34) ¶ added in v0.19.0

    type OAuthFlows struct {
     // Configuration for the OAuth Implicit flow.
     Implicit *OAuthFlow `json:"implicit,omitempty" yaml:"implicit,omitempty"`
     // Configuration for the OAuth Resource Owner Password flow.
     Password *OAuthFlow `json:"password,omitempty" yaml:"password,omitempty"`
     // Configuration for the OAuth Client Credentials flow. Previously called application in OpenAPI 2.0.
     ClientCredentials *OAuthFlow `json:"clientCredentials,omitempty" yaml:"clientCredentials,omitempty"`
     // Configuration for the OAuth Authorization Code flow. Previously called accessCode in OpenAPI 2.0.
     AuthorizationCode *OAuthFlow `json:"authorizationCode,omitempty" yaml:"authorizationCode,omitempty"`
    
     Common OpenAPICommon `json:"-" yaml:",inline"`
    }

OAuthFlows allows configuration of the supported OAuth Flows.

See <https://spec.openapis.org/oas/v3.1.0#oauth-flows-object>.

#### type [OpenAPICommon](https://github.com/ogen-go/ogen/blob/v1.18.0/spec.go#L30) ¶ added in v0.49.0

    type OpenAPICommon = [jsonschema](/github.com/ogen-go/ogen@v1.18.0/jsonschema).[OpenAPICommon](/github.com/ogen-go/ogen@v1.18.0/jsonschema#OpenAPICommon)

OpenAPICommon is a common OpenAPI object fields (extensions and locator).

#### type [Operation](https://github.com/ogen-go/ogen/blob/v1.18.0/spec.go#L369) ¶

    type Operation struct {
     // A list of tags for API documentation control.
     // Tags can be used for logical grouping of operations by resources or any other qualifier.
     Tags [][string](/builtin#string) `json:"tags,omitempty" yaml:"tags,omitempty"`
     // A short summary of what the operation does.
     Summary [string](/builtin#string) `json:"summary,omitempty" yaml:"summary,omitempty"`
     // A verbose explanation of the operation behavior.
     // CommonMark syntax MAY be used for rich text representation.
     Description [string](/builtin#string) `json:"description,omitempty" yaml:"description,omitempty"`
     // Additional external documentation for this operation.
     ExternalDocs *ExternalDocumentation `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`
    
     // Unique string used to identify the operation.
     //
     // The id MUST be unique among all operations described in the API.
     //
     // The operationId value is case-sensitive.
     OperationID [string](/builtin#string) `json:"operationId,omitempty" yaml:"operationId,omitempty"`
     // A list of parameters that are applicable for this operation.
     //
     // If a parameter is already defined at the Path Item, the new definition will override it but
     // can never remove it.
     //
     // The list MUST NOT include duplicated parameters. A unique parameter is defined by
     // a combination of a name and location.
     Parameters []*Parameter `json:"parameters,omitempty" yaml:"parameters,omitempty"`
     // The request body applicable for this operation.
     RequestBody *RequestBody `json:"requestBody,omitempty" yaml:"requestBody,omitempty"`
     // The list of possible responses as they are returned from executing this operation.
     Responses Responses `json:"responses,omitempty" yaml:"responses,omitempty"`
     // A map of possible out-of band callbacks related to the parent operation.
     //
     // The key is a unique identifier for the Callback Object.
     Callbacks map[[string](/builtin#string)]*Callback `json:"callbacks,omitempty" yaml:"callbacks,omitempty"`
     // Declares this operation to be deprecated
     Deprecated [bool](/builtin#bool) `json:"deprecated,omitempty" yaml:"deprecated,omitempty"`
     // A declaration of which security mechanisms can be used for this operation.
     //
     // The list of values includes alternative security requirement objects that can be used.
     //
     // Only one of the security requirement objects need to be satisfied to authorize a request.
     Security SecurityRequirements `json:"security,omitempty" yaml:"security,omitempty"`
     // An alternative server array to service this operation.
     //
     // If an alternative server object is specified at the Path Item Object or Root level,
     // it will be overridden by this value.
     Servers []Server `json:"servers,omitempty" yaml:"servers,omitempty"`
    
     Common OpenAPICommon `json:"-" yaml:",inline"`
    }

Operation describes a single API operation on a path.

See <https://spec.openapis.org/oas/v3.1.0#operation-object>.

#### func [NewOperation](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L496) ¶

    func NewOperation() *Operation

NewOperation returns a new Operation.

#### func (*Operation) [AddNamedResponses](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L562) ¶

    func (o *Operation) AddNamedResponses(ps ...*NamedResponse) *Operation

AddNamedResponses adds the given namedResponses to the Responses of the Operation.

#### func (*Operation) [AddParameters](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L537) ¶

    func (o *Operation) AddParameters(ps ...*Parameter) *Operation

AddParameters adds Parameters to the Parameters of the Operation.

#### func (*Operation) [AddResponse](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L555) ¶

    func (o *Operation) AddResponse(n [string](/builtin#string), p *Response) *Operation

AddResponse adds the given Response under the given Name to the Responses of the Operation.

#### func (*Operation) [AddTags](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L507) ¶

    func (o *Operation) AddTags(ts ...[string](/builtin#string)) *Operation

AddTags adds Tags to the Tags of the Operation.

#### func (*Operation) [MarshalJSON](https://github.com/ogen-go/ogen/blob/v1.18.0/spec.go#L421) ¶ added in v1.3.0

    func (s *Operation) MarshalJSON() ([][byte](/builtin#byte), [error](/builtin#error))

MarshalJSON implements [json.Marshaler](/encoding/json#Marshaler).

#### func (*Operation) [SetDescription](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L519) ¶

    func (o *Operation) SetDescription(d [string](/builtin#string)) *Operation

SetDescription sets the Description of the Operation.

#### func (*Operation) [SetOperationID](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L525) ¶

    func (o *Operation) SetOperationID(id [string](/builtin#string)) *Operation

SetOperationID sets the OperationID of the Operation.

#### func (*Operation) [SetParameters](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L531) ¶

    func (o *Operation) SetParameters(ps []*Parameter) *Operation

SetParameters sets the Parameters of the Operation.

#### func (*Operation) [SetRequestBody](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L543) ¶

    func (o *Operation) SetRequestBody(r *RequestBody) *Operation

SetRequestBody sets the RequestBody of the Operation.

#### func (*Operation) [SetResponses](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L549) ¶

    func (o *Operation) SetResponses(r Responses) *Operation

SetResponses sets the Responses of the Operation.

#### func (*Operation) [SetSummary](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L513) ¶

    func (o *Operation) SetSummary(s [string](/builtin#string)) *Operation

SetSummary sets the Summary of the Operation.

#### func (*Operation) [SetTags](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L501) ¶

    func (o *Operation) SetTags(ts [][string](/builtin#string)) *Operation

SetTags sets the Tags of the Operation.

#### type [Parameter](https://github.com/ogen-go/ogen/blob/v1.18.0/spec.go#L471) ¶

    type Parameter struct {
     Ref [string](/builtin#string) `json:"$ref,omitempty" yaml:"$ref,omitempty"` // ref object
    
     // REQUIRED. The name of the parameter. Parameter names are case sensitive.
     Name [string](/builtin#string) `json:"name,omitempty" yaml:"name,omitempty"`
     // REQUIRED. The location of the parameter. Possible values are "query", "header", "path" or "cookie".
     In [string](/builtin#string) `json:"in,omitempty" yaml:"in,omitempty"`
     // A brief description of the parameter. This could contain examples of use.
     // CommonMark syntax MAY be used for rich text representation.
     Description [string](/builtin#string) `json:"description,omitempty" yaml:"description,omitempty"`
     // Determines whether this parameter is mandatory.
     // If the parameter location is "path", this property is REQUIRED
     // and its value MUST be true.
     // Otherwise, the property MAY be included and its default value is false.
     Required [bool](/builtin#bool) `json:"required,omitempty" yaml:"required,omitempty"`
     // Specifies that a parameter is deprecated and SHOULD be transitioned out of usage.
     // Default value is false.
     Deprecated [bool](/builtin#bool) `json:"deprecated,omitempty" yaml:"deprecated,omitempty"`
    
     // Describes how the parameter value will be serialized
     // depending on the type of the parameter value.
     Style [string](/builtin#string) `json:"style,omitempty" yaml:"style,omitempty"`
     // When this is true, parameter values of type array or object
     // generate separate parameters for each value of the array
     // or key-value pair of the map.
     // For other types of parameters this property has no effect.
     Explode *[bool](/builtin#bool) `json:"explode,omitempty" yaml:"explode,omitempty"`
     // Determines whether the parameter value SHOULD allow reserved characters, as defined by [RFC 3986](https://rfc-editor.org/rfc/rfc3986.html).
     //
     // This property only applies to parameters with an in value of query.
     //
     // The default value is false.
     AllowReserved [bool](/builtin#bool) `json:"allowReserved,omitempty" yaml:"allowReserved,omitempty"`
     // The schema defining the type used for the parameter.
     Schema *Schema `json:"schema,omitempty" yaml:"schema,omitempty"`
     // Example of the parameter's potential value.
     Example ExampleValue `json:"example,omitempty" yaml:"example,omitempty"`
     // Examples of the parameter's potential value.
     Examples map[[string](/builtin#string)]*Example `json:"examples,omitempty" yaml:"examples,omitempty"`
    
     // For more complex scenarios, the content property can define the media type and schema of the parameter.
     // A parameter MUST contain either a schema property, or a content property, but not both.
     // When example or examples are provided in conjunction with the schema object,
     // the example MUST follow the prescribed serialization strategy for the parameter.
     //
     // A map containing the representations for the parameter.
     // The key is the media type and the value describes it.
     // The map MUST only contain one entry.
     Content map[[string](/builtin#string)]Media `json:"content,omitempty" yaml:"content,omitempty"`
    
     Common OpenAPICommon `json:"-" yaml:",inline"`
    }

Parameter describes a single operation parameter. A unique parameter is defined by a combination of a name and location.

See <https://spec.openapis.org/oas/v3.1.0#parameter-object>.

#### func [NewParameter](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L577) ¶

    func NewParameter() *Parameter

NewParameter returns a new Parameter.

#### func (*Parameter) [InCookie](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L615) ¶

    func (p *Parameter) InCookie() *Parameter

InCookie sets the In of the Parameter to "cookie".

#### func (*Parameter) [InHeader](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L610) ¶

    func (p *Parameter) InHeader() *Parameter

InHeader sets the In of the Parameter to "header".

#### func (*Parameter) [InPath](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L600) ¶

    func (p *Parameter) InPath() *Parameter

InPath sets the In of the Parameter to "path".

#### func (*Parameter) [InQuery](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L605) ¶

    func (p *Parameter) InQuery() *Parameter

InQuery sets the In of the Parameter to "query".

#### func (*Parameter) [SetContent](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L646) ¶

    func (p *Parameter) SetContent(c map[[string](/builtin#string)]Media) *Parameter

SetContent sets the Content of the Parameter.

#### func (*Parameter) [SetDeprecated](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L640) ¶

    func (p *Parameter) SetDeprecated(d [bool](/builtin#bool)) *Parameter

SetDeprecated sets the Deprecated of the Parameter.

#### func (*Parameter) [SetDescription](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L620) ¶

    func (p *Parameter) SetDescription(d [string](/builtin#string)) *Parameter

SetDescription sets the Description of the Parameter.

#### func (*Parameter) [SetExplode](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L660) ¶

    func (p *Parameter) SetExplode(e [bool](/builtin#bool)) *Parameter

SetExplode sets the Explode of the Parameter.

#### func (*Parameter) [SetIn](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L594) ¶

    func (p *Parameter) SetIn(i [string](/builtin#string)) *Parameter

SetIn sets the In of the Parameter.

#### func (*Parameter) [SetName](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L588) ¶

    func (p *Parameter) SetName(n [string](/builtin#string)) *Parameter

SetName sets the Name of the Parameter.

#### func (*Parameter) [SetRef](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L582) ¶

    func (p *Parameter) SetRef(r [string](/builtin#string)) *Parameter

SetRef sets the Ref of the Parameter.

#### func (*Parameter) [SetRequired](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L634) ¶

    func (p *Parameter) SetRequired(r [bool](/builtin#bool)) *Parameter

SetRequired sets the Required of the Parameter.

#### func (*Parameter) [SetSchema](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L626) ¶

    func (p *Parameter) SetSchema(s *Schema) *Parameter

SetSchema sets the Schema of the Parameter.

#### func (*Parameter) [SetStyle](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L654) ¶

    func (p *Parameter) SetStyle(s [string](/builtin#string)) *Parameter

SetStyle sets the Style of the Parameter.

#### func (*Parameter) [ToNamed](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L666) ¶

    func (p *Parameter) ToNamed(n [string](/builtin#string)) *NamedParameter

ToNamed returns a NamedParameter wrapping the receiver.

#### type [PathItem](https://github.com/ogen-go/ogen/blob/v1.18.0/spec.go#L290) ¶

    type PathItem struct {
     // Allows for an external definition of this path item.
     // The referenced structure MUST be in the format of a Path Item Object.
     // In case a Path Item Object field appears both
     // in the defined object and the referenced object, the behavior is undefined.
     Ref [string](/builtin#string) `json:"$ref,omitempty" yaml:"$ref,omitempty"`
    
     // An optional, string summary, intended to apply to all operations in this path.
     Summary [string](/builtin#string) `json:"summary,omitempty" yaml:"summary,omitempty"`
     // An optional, string description, intended to apply to all operations in this path.
     // CommonMark syntax MAY be used for rich text representation.
     Description [string](/builtin#string) `json:"description,omitempty" yaml:"description,omitempty"`
     // A definition of a GET operation on this path.
     Get *Operation `json:"get,omitempty" yaml:"get,omitempty"`
     // A definition of a PUT operation on this path.
     Put *Operation `json:"put,omitempty" yaml:"put,omitempty"`
     // A definition of a POST operation on this path.
     Post *Operation `json:"post,omitempty" yaml:"post,omitempty"`
     // A definition of a DELETE operation on this path.
     Delete *Operation `json:"delete,omitempty" yaml:"delete,omitempty"`
     // A definition of a OPTIONS operation on this path.
     Options *Operation `json:"options,omitempty" yaml:"options,omitempty"`
     // A definition of a HEAD operation on this path.
     Head *Operation `json:"head,omitempty" yaml:"head,omitempty"`
     // A definition of a PATCH operation on this path.
     Patch *Operation `json:"patch,omitempty" yaml:"patch,omitempty"`
     // A definition of a TRACE operation on this path.
     Trace *Operation `json:"trace,omitempty" yaml:"trace,omitempty"`
     // An alternative server array to service all operations in this path.
     Servers []Server `json:"servers,omitempty" yaml:"servers,omitempty"`
     // A list of parameters that are applicable for all the operations described under this path.
     //
     // These parameters can be overridden at the operation level, but cannot be removed there.
     //
     // The list MUST NOT include duplicated parameters. A unique parameter is defined by
     // a combination of a name and location.
     Parameters []*Parameter `json:"parameters,omitempty" yaml:"parameters,omitempty"`
    
     Common OpenAPICommon `json:"-" yaml:",inline"`
    }

PathItem Describes the operations available on a single path. A Path Item MAY be empty, due to ACL constraints. The path itself is still exposed to the documentation viewer, but they will not know which operations and parameters are available.

See <https://spec.openapis.org/oas/v3.1.0#path-item-object>.

#### func [NewPathItem](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L382) ¶

    func NewPathItem() *PathItem

NewPathItem returns a new PathItem.

#### func (*PathItem) [AddParameters](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L469) ¶

    func (p *PathItem) AddParameters(ps ...*Parameter) *PathItem

AddParameters adds Parameters to the Parameters of the PathItem.

#### func (*PathItem) [AddServers](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L453) ¶

    func (p *PathItem) AddServers(srvs ...*Server) *PathItem

AddServers adds Servers to the Servers of the PathItem.

#### func (*PathItem) [MarshalJSON](https://github.com/ogen-go/ogen/blob/v1.18.0/spec.go#L332) ¶ added in v1.3.0

    func (s *PathItem) MarshalJSON() ([][byte](/builtin#byte), [error](/builtin#error))

MarshalJSON implements [json.Marshaler](/encoding/json#Marshaler).

#### func (*PathItem) [SetDelete](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L417) ¶

    func (p *PathItem) SetDelete(o *Operation) *PathItem

SetDelete sets the Delete of the PathItem.

#### func (*PathItem) [SetDescription](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L393) ¶

    func (p *PathItem) SetDescription(d [string](/builtin#string)) *PathItem

SetDescription sets the Description of the PathItem.

#### func (*PathItem) [SetGet](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L399) ¶

    func (p *PathItem) SetGet(o *Operation) *PathItem

SetGet sets the Get of the PathItem.

#### func (*PathItem) [SetHead](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L429) ¶

    func (p *PathItem) SetHead(o *Operation) *PathItem

SetHead sets the Head of the PathItem.

#### func (*PathItem) [SetOptions](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L423) ¶

    func (p *PathItem) SetOptions(o *Operation) *PathItem

SetOptions sets the Options of the PathItem.

#### func (*PathItem) [SetParameters](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L463) ¶

    func (p *PathItem) SetParameters(ps []*Parameter) *PathItem

SetParameters sets the Parameters of the PathItem.

#### func (*PathItem) [SetPatch](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L435) ¶

    func (p *PathItem) SetPatch(o *Operation) *PathItem

SetPatch sets the Patch of the PathItem.

#### func (*PathItem) [SetPost](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L411) ¶

    func (p *PathItem) SetPost(o *Operation) *PathItem

SetPost sets the Post of the PathItem.

#### func (*PathItem) [SetPut](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L405) ¶

    func (p *PathItem) SetPut(o *Operation) *PathItem

SetPut sets the Put of the PathItem.

#### func (*PathItem) [SetRef](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L387) ¶

    func (p *PathItem) SetRef(r [string](/builtin#string)) *PathItem

SetRef sets the Ref of the PathItem.

#### func (*PathItem) [SetServers](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L447) ¶

    func (p *PathItem) SetServers(srvs []Server) *PathItem

SetServers sets the Servers of the PathItem.

#### func (*PathItem) [SetTrace](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L441) ¶

    func (p *PathItem) SetTrace(o *Operation) *PathItem

SetTrace sets the Trace of the PathItem.

#### func (*PathItem) [ToNamed](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L475) ¶

    func (p *PathItem) ToNamed(n [string](/builtin#string)) *NamedPathItem

ToNamed returns a NamedPathItem wrapping the receiver.

#### type [Paths](https://github.com/ogen-go/ogen/blob/v1.18.0/spec.go#L283) ¶

    type Paths map[[string](/builtin#string)]*PathItem

Paths holds the relative paths to the individual endpoints and their operations. The path is appended to the URL from the Server Object in order to construct the full URL. The Paths MAY be empty, due to ACL constraints.

See <https://spec.openapis.org/oas/v3.1.0#paths-object>.

#### type [PatternProperties](https://github.com/ogen-go/ogen/blob/v1.18.0/schema.go#L410) ¶ added in v0.23.0

    type PatternProperties []PatternProperty

PatternProperties is unparsed JSON Schema patternProperties validator description.

#### func (PatternProperties) [MarshalJSON](https://github.com/ogen-go/ogen/blob/v1.18.0/schema.go#L459) ¶ added in v0.23.0

    func (p PatternProperties) MarshalJSON() ([][byte](/builtin#byte), [error](/builtin#error))

MarshalJSON implements json.Marshaler.

#### func (PatternProperties) [MarshalYAML](https://github.com/ogen-go/ogen/blob/v1.18.0/schema.go#L413) ¶ added in v0.44.0

    func (p PatternProperties) MarshalYAML() ([any](/builtin#any), [error](/builtin#error))

MarshalYAML implements yaml.Marshaler.

#### func (PatternProperties) [ToJSONSchema](https://github.com/ogen-go/ogen/blob/v1.18.0/schema_backcomp.go#L96) ¶ added in v0.23.0

    func (p PatternProperties) ToJSONSchema() (result [jsonschema](/github.com/ogen-go/ogen@v1.18.0/jsonschema).[RawPatternProperties](/github.com/ogen-go/ogen@v1.18.0/jsonschema#RawPatternProperties))

ToJSONSchema converts PatternProperties to jsonschema.RawPatternProperties.

#### func (*PatternProperties) [UnmarshalJSON](https://github.com/ogen-go/ogen/blob/v1.18.0/schema.go#L476) ¶ added in v0.23.0

    func (p *PatternProperties) UnmarshalJSON(data [][byte](/builtin#byte)) [error](/builtin#error)

UnmarshalJSON implements json.Unmarshaler.

#### func (*PatternProperties) [UnmarshalYAML](https://github.com/ogen-go/ogen/blob/v1.18.0/schema.go#L433) ¶ added in v0.43.0

    func (p *PatternProperties) UnmarshalYAML(node *[yaml](/github.com/go-faster/yaml).[Node](/github.com/go-faster/yaml#Node)) [error](/builtin#error)

UnmarshalYAML implements yaml.Unmarshaler.

#### type [PatternProperty](https://github.com/ogen-go/ogen/blob/v1.18.0/schema.go#L404) ¶ added in v0.23.0

    type PatternProperty struct {
     Pattern [string](/builtin#string)
     Schema  *Schema
    }

PatternProperty is item of PatternProperties.

#### type [Properties](https://github.com/ogen-go/ogen/blob/v1.18.0/schema.go#L254) ¶

    type Properties []Property

Properties is unparsed JSON Schema properties validator description.

#### func (Properties) [MarshalJSON](https://github.com/ogen-go/ogen/blob/v1.18.0/schema.go#L303) ¶

    func (p Properties) MarshalJSON() ([][byte](/builtin#byte), [error](/builtin#error))

MarshalJSON implements json.Marshaler.

#### func (Properties) [MarshalYAML](https://github.com/ogen-go/ogen/blob/v1.18.0/schema.go#L257) ¶ added in v0.44.0

    func (p Properties) MarshalYAML() ([any](/builtin#any), [error](/builtin#error))

MarshalYAML implements yaml.Marshaler.

#### func (Properties) [ToJSONSchema](https://github.com/ogen-go/ogen/blob/v1.18.0/schema_backcomp.go#L62) ¶ added in v0.13.0

    func (p Properties) ToJSONSchema() [jsonschema](/github.com/ogen-go/ogen@v1.18.0/jsonschema).[RawProperties](/github.com/ogen-go/ogen@v1.18.0/jsonschema#RawProperties)

ToJSONSchema converts Properties to jsonschema.RawProperties.

#### func (*Properties) [UnmarshalJSON](https://github.com/ogen-go/ogen/blob/v1.18.0/schema.go#L320) ¶

    func (p *Properties) UnmarshalJSON(data [][byte](/builtin#byte)) [error](/builtin#error)

UnmarshalJSON implements json.Unmarshaler.

#### func (*Properties) [UnmarshalYAML](https://github.com/ogen-go/ogen/blob/v1.18.0/schema.go#L277) ¶ added in v0.43.0

    func (p *Properties) UnmarshalYAML(node *[yaml](/github.com/go-faster/yaml).[Node](/github.com/go-faster/yaml#Node)) [error](/builtin#error)

UnmarshalYAML implements yaml.Unmarshaler.

#### type [Property](https://github.com/ogen-go/ogen/blob/v1.18.0/schema.go#L248) ¶

    type Property struct {
     Name   [string](/builtin#string)
     Schema *Schema
    }

Property is item of Properties.

#### func [NewProperty](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L1078) ¶

    func NewProperty() *Property

NewProperty returns a new Property.

#### func (*Property) [SetName](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L1083) ¶

    func (p *Property) SetName(n [string](/builtin#string)) *Property

SetName sets the Name of the Property.

#### func (*Property) [SetSchema](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L1089) ¶

    func (p *Property) SetSchema(s *Schema) *Property

SetSchema sets the Schema of the Property.

#### func (Property) [ToJSONSchema](https://github.com/ogen-go/ogen/blob/v1.18.0/schema_backcomp.go#L71) ¶ added in v0.13.0

    func (p Property) ToJSONSchema() [jsonschema](/github.com/ogen-go/ogen@v1.18.0/jsonschema).[RawProperty](/github.com/ogen-go/ogen@v1.18.0/jsonschema#RawProperty)

ToJSONSchema converts Property to jsonschema.Property.

#### type [RawValue](https://github.com/ogen-go/ogen/blob/v1.18.0/spec.go#L23) ¶ added in v0.44.0

    type RawValue = [jsonschema](/github.com/ogen-go/ogen@v1.18.0/jsonschema).[RawValue](/github.com/ogen-go/ogen@v1.18.0/jsonschema#RawValue)

RawValue is a raw JSON value.

#### type [RequestBody](https://github.com/ogen-go/ogen/blob/v1.18.0/spec.go#L527) ¶

    type RequestBody struct {
     Ref [string](/builtin#string) `json:"$ref,omitempty" yaml:"$ref,omitempty"` // ref object
    
     // A brief description of the request body. This could contain examples of use.
     // CommonMark syntax MAY be used for rich text representation.
     Description [string](/builtin#string) `json:"description,omitempty" yaml:"description,omitempty"`
    
     // REQUIRED. The content of the request body.
     //
     // The key is a media type or media type range and the value describes it.
     //
     // For requests that match multiple keys, only the most specific key is applicable.
     // e.g. text/plain overrides text/*
     Content map[[string](/builtin#string)]Media `json:"content,omitempty" yaml:"content,omitempty"`
    
     // Determines if the request body is required in the request. Defaults to false.
     Required [bool](/builtin#bool) `json:"required,omitempty" yaml:"required,omitempty"`
    
     Common OpenAPICommon `json:"-" yaml:",inline"`
    }

RequestBody describes a single request body.

See <https://spec.openapis.org/oas/v3.1.0#request-body-object>.

#### func [NewRequestBody](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L213) ¶

    func NewRequestBody() *RequestBody

NewRequestBody returns a new RequestBody.

#### func (*RequestBody) [AddContent](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L236) ¶

    func (r *RequestBody) AddContent(mt [string](/builtin#string), s *Schema) *RequestBody

AddContent adds the given Schema under the MediaType to the Content of the Response.

#### func (*RequestBody) [SetContent](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L230) ¶

    func (r *RequestBody) SetContent(c map[[string](/builtin#string)]Media) *RequestBody

SetContent sets the Content of the RequestBody.

#### func (*RequestBody) [SetDescription](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L224) ¶

    func (r *RequestBody) SetDescription(d [string](/builtin#string)) *RequestBody

SetDescription sets the Description of the RequestBody.

#### func (*RequestBody) [SetJSONContent](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L245) ¶

    func (r *RequestBody) SetJSONContent(s *Schema) *RequestBody

SetJSONContent sets the given Schema under the JSON MediaType to the Content of the Response.

#### func (*RequestBody) [SetRef](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L218) ¶

    func (r *RequestBody) SetRef(ref [string](/builtin#string)) *RequestBody

SetRef sets the Ref of the RequestBody.

#### func (*RequestBody) [SetRequired](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L257) ¶

    func (r *RequestBody) SetRequired(req [bool](/builtin#bool)) *RequestBody

SetRequired sets the Required of the RequestBody.

#### func (*RequestBody) [ToNamed](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L263) ¶

    func (r *RequestBody) ToNamed(n [string](/builtin#string)) *NamedRequestBody

ToNamed returns a NamedRequestBody wrapping the receiver.

#### type [Response](https://github.com/ogen-go/ogen/blob/v1.18.0/spec.go#L618) ¶

    type Response struct {
     Ref [string](/builtin#string) `json:"$ref,omitempty" yaml:"$ref,omitempty"` // ref object
    
     // REQUIRED. A description of the response.
     // CommonMark syntax MAY be used for rich text representation.
     Description [string](/builtin#string) `json:"description,omitempty" yaml:"description,omitempty"`
     // Maps a header name to its definition.
     //
     // RFC7230 states header names are case insensitive.
     //
     // If a response header is defined with the name "Content-Type", it SHALL be ignored.
     Headers map[[string](/builtin#string)]*Header `json:"headers,omitempty" yaml:"headers,omitempty"`
     // A map containing descriptions of potential response payloads.
     //
     // The key is a media type or media type range and the value describes it.
     //
     // For requests that match multiple keys, only the most specific key is applicable.
     // e.g. text/plain overrides text/*
     Content map[[string](/builtin#string)]Media `json:"content,omitempty" yaml:"content,omitempty"`
     // A map of operations links that can be followed from the response.
     //
     // The key of the map is a short name for the link, following the naming constraints
     // of the names for Component Objects.
     Links map[[string](/builtin#string)]*Link `json:"links,omitempty" yaml:"links,omitempty"`
    
     Common OpenAPICommon `json:"-" yaml:",inline"`
    }

Response describes a single response from an API Operation, including design-time, static links to operations based on the response.

See <https://spec.openapis.org/oas/v3.1.0#response-object>.

#### func [NewResponse](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L687) ¶

    func NewResponse() *Response

NewResponse returns a new Response.

#### func (*Response) [AddContent](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L716) ¶

    func (r *Response) AddContent(mt [string](/builtin#string), s *Schema) *Response

AddContent adds the given Schema under the MediaType to the Content of the Response.

#### func (*Response) [SetContent](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L710) ¶

    func (r *Response) SetContent(c map[[string](/builtin#string)]Media) *Response

SetContent sets the Content of the Response.

#### func (*Response) [SetDescription](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L698) ¶

    func (r *Response) SetDescription(d [string](/builtin#string)) *Response

SetDescription sets the Description of the Response.

#### func (*Response) [SetHeaders](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L704) ¶ added in v0.33.0

    func (r *Response) SetHeaders(h map[[string](/builtin#string)]*Header) *Response

SetHeaders sets the Headers of the Response.

#### func (*Response) [SetJSONContent](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L725) ¶

    func (r *Response) SetJSONContent(s *Schema) *Response

SetJSONContent sets the given Schema under the JSON MediaType to the Content of the Response.

#### func (*Response) [SetLinks](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L737) ¶

    func (r *Response) SetLinks(l map[[string](/builtin#string)]*Link) *Response

SetLinks sets the Links of the Response.

#### func (*Response) [SetRef](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L692) ¶

    func (r *Response) SetRef(ref [string](/builtin#string)) *Response

SetRef sets the Ref of the Response.

#### func (*Response) [ToNamed](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L743) ¶

    func (r *Response) ToNamed(n [string](/builtin#string)) *NamedResponse

ToNamed returns a NamedResponse wrapping the receiver.

#### type [Responses](https://github.com/ogen-go/ogen/blob/v1.18.0/spec.go#L612) ¶

    type Responses map[[string](/builtin#string)]*Response

Responses is a container for the expected responses of an operation.

The container maps the HTTP response code to the expected response.

The `default` MAY be used as a default response object for all HTTP codes that are not covered individually by the Responses Object.

The Responses Object MUST contain at least one response code, and if only one response code is provided it SHOULD be the response for a successful operation call.

See <https://spec.openapis.org/oas/v3.1.0#responses-object>.

#### type [Schema](https://github.com/ogen-go/ogen/blob/v1.18.0/schema.go#L16) ¶

    type Schema struct {
     Ref         [string](/builtin#string) `json:"$ref,omitempty" yaml:"$ref,omitempty"` // ref object
     Summary     [string](/builtin#string) `json:"summary,omitempty" yaml:"summary,omitempty"`
     Description [string](/builtin#string) `json:"description,omitempty" yaml:"description,omitempty"`
    
     // Additional external documentation for this schema.
     ExternalDocs *ExternalDocumentation `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`
    
     // Value MUST be a string. Multiple types via an array are not supported.
     Type [string](/builtin#string) `json:"type,omitempty" yaml:"type,omitempty"`
    
     // See Data Type Formats for further details (<https://swagger.io/specification/#data-type-format>).
     // While relying on JSON Schema's defined formats,
     // the OAS offers a few additional predefined formats.
     Format [string](/builtin#string) `json:"format,omitempty" yaml:"format,omitempty"`
    
     // Property definitions MUST be a Schema Object and not a standard JSON Schema
     // (inline or referenced).
     Properties Properties `json:"properties,omitempty" yaml:"properties,omitempty"`
    
     // Value can be boolean or object. Inline or referenced schema MUST be of a Schema Object
     // and not a standard JSON Schema. Consistent with JSON Schema, additionalProperties defaults to true.
     AdditionalProperties *AdditionalProperties `json:"additionalProperties,omitempty" yaml:"additionalProperties,omitempty"`
    
     // The value of "patternProperties" MUST be an object. Each property
     // name of this object SHOULD be a valid regular expression, according
     // to the ECMA-262 regular expression dialect. Each property value of
     // this object MUST be a valid JSON Schema.
     PatternProperties PatternProperties `json:"patternProperties,omitempty" yaml:"patternProperties,omitempty"`
    
     // The value of this keyword MUST be an array.
     // This array MUST have at least one element.
     // Elements of this array MUST be strings, and MUST be unique.
     Required [][string](/builtin#string) `json:"required,omitempty" yaml:"required,omitempty"`
    
     // Value MUST be an object and not an array.
     // Inline or referenced schema MUST be of a Schema Object and not a standard JSON Schema.
     // MUST be present if the Type is "array".
     Items *Items `json:"items,omitempty" yaml:"items,omitempty"`
    
     // A true value adds "null" to the allowed type specified by the type keyword,
     // only if type is explicitly defined within the same Schema Object.
     // Other Schema Object constraints retain their defined behavior,
     // and therefore may disallow the use of null as a value.
     // A false value leaves the specified or default type unmodified.
     // The default value is false.
     Nullable [bool](/builtin#bool) `json:"nullable,omitempty" yaml:"nullable,omitempty"`
    
     // AllOf takes an array of object definitions that are used
     // for independent validation but together compose a single object.
     // Still, it does not imply a hierarchy between the models.
     // For that purpose, you should include the discriminator.
     AllOf []*Schema `json:"allOf,omitempty" yaml:"allOf,omitempty"`
    
     // OneOf validates the value against exactly one of the subschemas
     OneOf []*Schema `json:"oneOf,omitempty" yaml:"oneOf,omitempty"`
    
     // AnyOf validates the value against any (one or more) of the subschemas
     AnyOf []*Schema `json:"anyOf,omitempty" yaml:"anyOf,omitempty"`
    
     // Discriminator for subschemas.
     Discriminator *Discriminator `json:"discriminator,omitempty" yaml:"discriminator,omitempty"`
    
     // Adds additional metadata to describe the XML representation of this property.
     //
     // This MAY be used only on properties schemas. It has no effect on root schemas
     XML *XML `json:"xml,omitempty" yaml:"xml,omitempty"`
    
     // The value of this keyword MUST be an array.
     // This array SHOULD have at least one element.
     // Elements in the array SHOULD be unique.
     Enum Enum `json:"enum,omitempty" yaml:"enum,omitempty"`
    
     // The value of "multipleOf" MUST be a number, strictly greater than 0.
     //
     // A numeric instance is only valid if division by this keyword's value
     // results in an integer.
     MultipleOf Num `json:"multipleOf,omitempty" yaml:"multipleOf,omitempty"`
    
     // The value of "maximum" MUST be a number, representing an upper limit
     // for a numeric instance.
     //
     // If the instance is a number, then this keyword validates if
     // "exclusiveMaximum" is true and instance is less than the provided
     // value, or else if the instance is less than or exactly equal to the
     // provided value.
     Maximum Num `json:"maximum,omitempty" yaml:"maximum,omitempty"`
    
     // The value of "exclusiveMaximum" MUST be a boolean, representing
     // whether the limit in "maximum" is exclusive or not.  An undefined
     // value is the same as false.
     //
     // If "exclusiveMaximum" is true, then a numeric instance SHOULD NOT be
     // equal to the value specified in "maximum".  If "exclusiveMaximum" is
     // false (or not specified), then a numeric instance MAY be equal to the
     // value of "maximum".
     ExclusiveMaximum [bool](/builtin#bool) `json:"exclusiveMaximum,omitempty" yaml:"exclusiveMaximum,omitempty"`
    
     // The value of "minimum" MUST be a number, representing a lower limit
     // for a numeric instance.
     //
     // If the instance is a number, then this keyword validates if
     // "exclusiveMinimum" is true and instance is greater than the provided
     // value, or else if the instance is greater than or exactly equal to
     // the provided value.
     Minimum Num `json:"minimum,omitempty" yaml:"minimum,omitempty"`
    
     // The value of "exclusiveMinimum" MUST be a boolean, representing
     // whether the limit in "minimum" is exclusive or not.  An undefined
     // value is the same as false.
     //
     // If "exclusiveMinimum" is true, then a numeric instance SHOULD NOT be
     // equal to the value specified in "minimum".  If "exclusiveMinimum" is
     // false (or not specified), then a numeric instance MAY be equal to the
     // value of "minimum".
     ExclusiveMinimum [bool](/builtin#bool) `json:"exclusiveMinimum,omitempty" yaml:"exclusiveMinimum,omitempty"`
    
     // The value of this keyword MUST be a non-negative integer.
     //
     // The value of this keyword MUST be an integer.  This integer MUST be
     // greater than, or equal to, 0.
     //
     // A string instance is valid against this keyword if its length is less
     // than, or equal to, the value of this keyword.
     //
     // The length of a string instance is defined as the number of its
     // characters as defined by [RFC 7159](https://rfc-editor.org/rfc/rfc7159.html) [RFC7159].
     MaxLength *[uint64](/builtin#uint64) `json:"maxLength,omitempty" yaml:"maxLength,omitempty"`
    
     // A string instance is valid against this keyword if its length is
     // greater than, or equal to, the value of this keyword.
     //
     // The length of a string instance is defined as the number of its
     // characters as defined by [RFC 7159](https://rfc-editor.org/rfc/rfc7159.html) [RFC7159].
     //
     // The value of this keyword MUST be an integer.  This integer MUST be
     // greater than, or equal to, 0.
     //
     // "minLength", if absent, may be considered as being present with
     // integer value 0.
     MinLength *[uint64](/builtin#uint64) `json:"minLength,omitempty" yaml:"minLength,omitempty"`
    
     // The value of this keyword MUST be a string.  This string SHOULD be a
     // valid regular expression, according to the ECMA 262 regular
     // expression dialect.
     //
     // A string instance is considered valid if the regular expression
     // matches the instance successfully. Recall: regular expressions are
     // not implicitly anchored.
     Pattern [string](/builtin#string) `json:"pattern,omitempty" yaml:"pattern,omitempty"`
    
     // The value of this keyword MUST be an integer.  This integer MUST be
     // greater than, or equal to, 0.
     //
     // An array instance is valid against "maxItems" if its size is less
     // than, or equal to, the value of this keyword.
     MaxItems *[uint64](/builtin#uint64) `json:"maxItems,omitempty" yaml:"maxItems,omitempty"`
    
     // The value of this keyword MUST be an integer.  This integer MUST be
     // greater than, or equal to, 0.
     //
     // An array instance is valid against "minItems" if its size is greater
     // than, or equal to, the value of this keyword.
     //
     // If this keyword is not present, it may be considered present with a
     // value of 0.
     MinItems *[uint64](/builtin#uint64) `json:"minItems,omitempty" yaml:"minItems,omitempty"`
    
     // The value of this keyword MUST be a boolean.
     //
     // If this keyword has boolean value false, the instance validates
     // successfully.  If it has boolean value true, the instance validates
     // successfully if all of its elements are unique.
     //
     // If not present, this keyword may be considered present with boolean
     // value false.
     UniqueItems [bool](/builtin#bool) `json:"uniqueItems,omitempty" yaml:"uniqueItems,omitempty"`
    
     // The value of this keyword MUST be an integer.  This integer MUST be
     // greater than, or equal to, 0.
     //
     // An object instance is valid against "maxProperties" if its number of
     // properties is less than, or equal to, the value of this keyword.
     MaxProperties *[uint64](/builtin#uint64) `json:"maxProperties,omitempty" yaml:"maxProperties,omitempty"`
    
     // The value of this keyword MUST be an integer.  This integer MUST be
     // greater than, or equal to, 0.
     //
     // An object instance is valid against "minProperties" if its number of
     // properties is greater than, or equal to, the value of this keyword.
     //
     // If this keyword is not present, it may be considered present with a
     // value of 0.
     MinProperties *[uint64](/builtin#uint64) `json:"minProperties,omitempty" yaml:"minProperties,omitempty"`
    
     // Default value.
     Default Default `json:"default,omitempty" yaml:"default,omitempty"`
    
     // A free-form property to include an example of an instance for this schema.
     // To represent examples that cannot be naturally represented in JSON or YAML,
     // a string value can be used to contain the example with escaping where necessary.
     Example ExampleValue `json:"example,omitempty" yaml:"example,omitempty"`
    
     // Specifies that a schema is deprecated and SHOULD be transitioned out
     // of usage.
     Deprecated [bool](/builtin#bool) `json:"deprecated,omitempty" yaml:"deprecated,omitempty"`
    
     // If the instance value is a string, this property defines that the
     // string SHOULD be interpreted as binary data and decoded using the
     // encoding named by this property.  [RFC 2045, Section 6.1](https://rfc-editor.org/rfc/rfc2045.html#section-6.1) lists
     // the possible values for this property.
     //
     // The value of this property MUST be a string.
     //
     // The value of this property SHOULD be ignored if the instance
     // described is not a string.
     ContentEncoding [string](/builtin#string) `json:"contentEncoding,omitempty" yaml:"contentEncoding,omitempty"`
    
     // The value of this property must be a media type, as defined by RFC
     // 2046. This property defines the media type of instances
     // which this schema defines.
     //
     // The value of this property MUST be a string.
     //
     // The value of this property SHOULD be ignored if the instance
     // described is not a string.
     ContentMediaType [string](/builtin#string) `json:"contentMediaType,omitempty" yaml:"contentMediaType,omitempty"`
    
     Common [jsonschema](/github.com/ogen-go/ogen@v1.18.0/jsonschema).[OpenAPICommon](/github.com/ogen-go/ogen@v1.18.0/jsonschema#OpenAPICommon) `json:"-" yaml:",inline"`
    }

The Schema Object allows the definition of input and output data types. These types can be objects, but also primitives and arrays.

#### func [Binary](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L1018) ¶

    func Binary() *Schema

Binary returns a sequence of octets OAS data type (Schema).

#### func [Bool](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L1021) ¶

    func Bool() *Schema

Bool returns a boolean OAS data type (Schema).

#### func [Bytes](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L1015) ¶

    func Bytes() *Schema

Bytes returns a base64 encoded OAS data type (Schema).

#### func [Date](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L1024) ¶

    func Date() *Schema

Date returns a date as defined by full-date - RFC3339 OAS data type (Schema).

#### func [DateTime](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L1027) ¶

    func DateTime() *Schema

DateTime returns a date as defined by date-time - RFC3339 OAS data type (Schema).

#### func [Double](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L1006) ¶

    func Double() *Schema

Double returns a double OAS data type (Schema).

#### func [Float](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L1003) ¶

    func Float() *Schema

Float returns a float OAS data type (Schema).

#### func [Int](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L994) ¶ added in v0.2.0

    func Int() *Schema

Int returns an integer OAS data type (Schema).

#### func [Int32](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L997) ¶

    func Int32() *Schema

Int32 returns an 32-bit integer OAS data type (Schema).

#### func [Int64](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L1000) ¶

    func Int64() *Schema

Int64 returns an 64-bit integer OAS data type (Schema).

#### func [NewSchema](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L766) ¶

    func NewSchema() *Schema

NewSchema returns a new Schema.

#### func [Password](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L1030) ¶

    func Password() *Schema

Password returns an obscured OAS data type (Schema).

#### func [String](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L1009) ¶

    func String() *Schema

String returns a string OAS data type (Schema).

#### func [UUID](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L1012) ¶ added in v0.2.0

    func UUID() *Schema

UUID returns a UUID OAS data type (Schema).

#### func (*Schema) [AddOptionalProperties](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L810) ¶

    func (s *Schema) AddOptionalProperties(ps ...*Property) *Schema

AddOptionalProperties adds the Properties to the Properties of the Schema.

#### func (*Schema) [AddRequiredProperties](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L821) ¶

    func (s *Schema) AddRequiredProperties(ps ...*Property) *Schema

AddRequiredProperties adds the Properties to the Properties of the Schema and marks them as required.

#### func (*Schema) [AsArray](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L1038) ¶

    func (s *Schema) AsArray() *Schema

AsArray returns a new "array" Schema wrapping the receiver.

#### func (*Schema) [AsEnum](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L1048) ¶

    func (s *Schema) AsEnum(def [json](/encoding/json).[RawMessage](/encoding/json#RawMessage), values ...[json](/encoding/json).[RawMessage](/encoding/json#RawMessage)) *Schema

AsEnum returns a new "enum" Schema wrapping the receiver.

#### func (*Schema) [SetAllOf](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L852) ¶

    func (s *Schema) SetAllOf(a []*Schema) *Schema

SetAllOf sets the AllOf of the Schema.

#### func (*Schema) [SetAnyOf](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L864) ¶

    func (s *Schema) SetAnyOf(a []*Schema) *Schema

SetAnyOf sets the AnyOf of the Schema.

#### func (*Schema) [SetDefault](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L977) ¶

    func (s *Schema) SetDefault(d [json](/encoding/json).[RawMessage](/encoding/json#RawMessage)) *Schema

SetDefault sets the Default of the Schema.

#### func (*Schema) [SetDeprecated](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L983) ¶ added in v0.68.0

    func (s *Schema) SetDeprecated(d [bool](/builtin#bool)) *Schema

SetDeprecated sets the Deprecated of the Schema.

#### func (*Schema) [SetDescription](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L783) ¶

    func (s *Schema) SetDescription(d [string](/builtin#string)) *Schema

SetDescription sets the Description of the Schema.

#### func (*Schema) [SetDiscriminator](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L870) ¶

    func (s *Schema) SetDiscriminator(d *Discriminator) *Schema

SetDiscriminator sets the Discriminator of the Schema.

#### func (*Schema) [SetEnum](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L876) ¶

    func (s *Schema) SetEnum(e [][json](/encoding/json).[RawMessage](/encoding/json#RawMessage)) *Schema

SetEnum sets the Enum of the Schema.

#### func (*Schema) [SetExclusiveMaximum](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L906) ¶

    func (s *Schema) SetExclusiveMaximum(e [bool](/builtin#bool)) *Schema

SetExclusiveMaximum sets the ExclusiveMaximum of the Schema.

#### func (*Schema) [SetExclusiveMinimum](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L923) ¶

    func (s *Schema) SetExclusiveMinimum(e [bool](/builtin#bool)) *Schema

SetExclusiveMinimum sets the ExclusiveMinimum of the Schema.

#### func (*Schema) [SetFormat](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L795) ¶

    func (s *Schema) SetFormat(f [string](/builtin#string)) *Schema

SetFormat sets the Format of the Schema.

#### func (*Schema) [SetItems](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L838) ¶

    func (s *Schema) SetItems(i *Schema) *Schema

SetItems sets the Items of the Schema.

#### func (*Schema) [SetMaxItems](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L947) ¶

    func (s *Schema) SetMaxItems(m *[uint64](/builtin#uint64)) *Schema

SetMaxItems sets the MaxItems of the Schema.

#### func (*Schema) [SetMaxLength](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L929) ¶

    func (s *Schema) SetMaxLength(m *[uint64](/builtin#uint64)) *Schema

SetMaxLength sets the MaxLength of the Schema.

#### func (*Schema) [SetMaxProperties](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L965) ¶

    func (s *Schema) SetMaxProperties(m *[uint64](/builtin#uint64)) *Schema

SetMaxProperties sets the MaxProperties of the Schema.

#### func (*Schema) [SetMaximum](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L895) ¶

    func (s *Schema) SetMaximum(m *[int64](/builtin#int64)) *Schema

SetMaximum sets the Maximum of the Schema.

#### func (*Schema) [SetMinItems](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L953) ¶

    func (s *Schema) SetMinItems(m *[uint64](/builtin#uint64)) *Schema

SetMinItems sets the MinItems of the Schema.

#### func (*Schema) [SetMinLength](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L935) ¶

    func (s *Schema) SetMinLength(m *[uint64](/builtin#uint64)) *Schema

SetMinLength sets the MinLength of the Schema.

#### func (*Schema) [SetMinProperties](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L971) ¶

    func (s *Schema) SetMinProperties(m *[uint64](/builtin#uint64)) *Schema

SetMinProperties sets the MinProperties of the Schema.

#### func (*Schema) [SetMinimum](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L912) ¶

    func (s *Schema) SetMinimum(m *[int64](/builtin#int64)) *Schema

SetMinimum sets the Minimum of the Schema.

#### func (*Schema) [SetMultipleOf](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L884) ¶

    func (s *Schema) SetMultipleOf(m *[uint64](/builtin#uint64)) *Schema

SetMultipleOf sets the MultipleOf of the Schema.

#### func (*Schema) [SetNullable](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L846) ¶

    func (s *Schema) SetNullable(n [bool](/builtin#bool)) *Schema

SetNullable sets the Nullable of the Schema.

#### func (*Schema) [SetOneOf](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L858) ¶

    func (s *Schema) SetOneOf(o []*Schema) *Schema

SetOneOf sets the OneOf of the Schema.

#### func (*Schema) [SetPattern](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L941) ¶

    func (s *Schema) SetPattern(p [string](/builtin#string)) *Schema

SetPattern sets the Pattern of the Schema.

#### func (*Schema) [SetProperties](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L801) ¶

    func (s *Schema) SetProperties(p *Properties) *Schema

SetProperties sets the Properties of the Schema.

#### func (*Schema) [SetRef](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L771) ¶

    func (s *Schema) SetRef(r [string](/builtin#string)) *Schema

SetRef sets the Ref of the Schema.

#### func (*Schema) [SetRequired](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L832) ¶

    func (s *Schema) SetRequired(r [][string](/builtin#string)) *Schema

SetRequired sets the Required of the Schema.

#### func (*Schema) [SetSummary](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L777) ¶ added in v0.68.3

    func (s *Schema) SetSummary(smry [string](/builtin#string)) *Schema

SetSummary sets the Summary of the Schema.

#### func (*Schema) [SetType](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L789) ¶

    func (s *Schema) SetType(t [string](/builtin#string)) *Schema

SetType sets the Type of the Schema.

#### func (*Schema) [SetUniqueItems](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L959) ¶

    func (s *Schema) SetUniqueItems(u [bool](/builtin#bool)) *Schema

SetUniqueItems sets the UniqueItems of the Schema.

#### func (*Schema) [ToJSONSchema](https://github.com/ogen-go/ogen/blob/v1.18.0/schema_backcomp.go#L8) ¶ added in v0.13.0

    func (s *Schema) ToJSONSchema() *[jsonschema](/github.com/ogen-go/ogen@v1.18.0/jsonschema).[RawSchema](/github.com/ogen-go/ogen@v1.18.0/jsonschema#RawSchema)

ToJSONSchema converts Schema to jsonschema.Schema.

#### func (*Schema) [ToNamed](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L989) ¶

    func (s *Schema) ToNamed(n [string](/builtin#string)) *NamedSchema

ToNamed returns a NamedSchema wrapping the receiver.

#### func (*Schema) [ToProperty](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L1057) ¶

    func (s *Schema) ToProperty(n [string](/builtin#string)) *Property

ToProperty returns a Property with the given name and with this Schema.

#### type [SecurityRequirement](https://github.com/ogen-go/ogen/blob/v1.18.0/security_scheme.go#L71) ¶ added in v0.57.0

    type SecurityRequirement = map[[string](/builtin#string)][][string](/builtin#string)

SecurityRequirement lists the required security schemes to execute this operation.

See <https://spec.openapis.org/oas/v3.1.0#security-requirement-object>.

#### type [SecurityRequirements](https://github.com/ogen-go/ogen/blob/v1.18.0/security_scheme.go#L73) ¶ added in v0.19.0

    type SecurityRequirements []SecurityRequirement

SecurityRequirements lists the security requirements of the operation.

#### type [SecurityScheme](https://github.com/ogen-go/ogen/blob/v1.18.0/security_scheme.go#L6) ¶ added in v0.42.0

    type SecurityScheme struct {
     Ref [string](/builtin#string) `json:"$ref,omitempty" yaml:"$ref,omitempty"`
     // The type of the security scheme. Valid values are "apiKey", "http", "mutualTLS", "oauth2", "openIdConnect".
     Type [string](/builtin#string) `json:"type" yaml:"type,omitempty"`
     // A description for security scheme. CommonMark syntax MAY be used for rich text representation.
     Description [string](/builtin#string) `json:"description,omitempty" yaml:"description,omitempty"`
     // The name of the header, query or cookie parameter to be used.
     Name [string](/builtin#string) `json:"name,omitempty" yaml:"name,omitempty"`
     // The location of the API key. Valid values are "query", "header" or "cookie".
     In [string](/builtin#string) `json:"in,omitempty" yaml:"in,omitempty"`
     // The name of the HTTP Authorization scheme to be used in the Authorization header as defined in RFC7235.
     // The values used SHOULD be registered in the IANA Authentication Scheme registry.
     Scheme [string](/builtin#string) `json:"scheme,omitempty" yaml:"scheme,omitempty"`
     // A hint to the client to identify how the bearer token is formatted. Bearer tokens are usually generated
     // by an authorization server, so this information is primarily for documentation purposes.
     BearerFormat [string](/builtin#string) `json:"bearerFormat,omitempty" yaml:"bearerFormat,omitempty"`
     // An object containing configuration information for the flow types supported.
     Flows *OAuthFlows `json:"flows,omitempty" yaml:"flows,omitempty"`
     // OpenId Connect URL to discover OAuth2 configuration values.
     // This MUST be in the form of a URL. The OpenID Connect standard requires the use of TLS.
     OpenIDConnectURL [string](/builtin#string) `json:"openIdConnectUrl,omitempty" yaml:"openIdConnectUrl,omitempty"`
    
     Common OpenAPICommon `json:"-" yaml:",inline"`
    }

SecurityScheme defines a security scheme that can be used by the operations.

See <https://spec.openapis.org/oas/v3.1.0#security-scheme-object>.

#### type [Server](https://github.com/ogen-go/ogen/blob/v1.18.0/spec.go#L192) ¶

    type Server struct {
     // REQUIRED. A URL to the target host. This URL supports Server Variables and MAY be relative,
     // to indicate that the host location is relative to the location where the OpenAPI document is being served.
     // Variable substitutions will be made when a variable is named in {brackets}.
     URL [string](/builtin#string) `json:"url" yaml:"url"`
     // An optional string describing the host designated by the URL.
     // CommonMark syntax MAY be used for rich text representation.
     Description [string](/builtin#string) `json:"description,omitempty" yaml:"description,omitempty"`
     // A map between a variable name and its value. The value is used for substitution in the server's URL template.
     Variables map[[string](/builtin#string)]ServerVariable `json:"variables,omitempty" yaml:"variables,omitempty"`
    
     Common OpenAPICommon `json:"-" yaml:",inline"`
    }

Server represents a Server.

See <https://spec.openapis.org/oas/v3.1.0#server-object>.

#### func [NewServer](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L365) ¶

    func NewServer() *Server

NewServer returns a new Server.

#### func (*Server) [SetDescription](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L370) ¶

    func (s *Server) SetDescription(d [string](/builtin#string)) *Server

SetDescription sets the Description of the Server.

#### func (*Server) [SetURL](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L376) ¶

    func (s *Server) SetURL(url [string](/builtin#string)) *Server

SetURL sets the URL of the Server.

#### type [ServerVariable](https://github.com/ogen-go/ogen/blob/v1.18.0/spec.go#L209) ¶ added in v0.40.0

    type ServerVariable struct {
     // An enumeration of string values to be used if the substitution options are from a limited set.
     //
     // The array MUST NOT be empty.
     Enum [][string](/builtin#string) `json:"enum,omitempty" yaml:"enum,omitempty"`
     // REQUIRED. The default value to use for substitution, which SHALL be sent if an alternate value is not supplied.
     // Note this behavior is different than the Schema Object's treatment of default values, because in those
     // cases parameter values are optional. If the enum is defined, the value MUST exist in the enum's values.
     Default [string](/builtin#string) `json:"default" yaml:"default"`
     // An optional description for the server variable. CommonMark syntax MAY be used for rich text representation.
     Description [string](/builtin#string) `json:"description,omitempty" yaml:"description,omitempty"`
    
     Common OpenAPICommon `json:"-" yaml:",inline"`
    }

ServerVariable describes an object representing a Server Variable for server URL template substitution.

See <https://spec.openapis.org/oas/v3.1.0#server-variable-object>

#### type [Spec](https://github.com/ogen-go/ogen/blob/v1.18.0/spec.go#L39) ¶

    type Spec struct {
     // REQUIRED. This string MUST be the version number of the OpenAPI Specification
     // that the OpenAPI document uses.
     OpenAPI [string](/builtin#string) `json:"openapi" yaml:"openapi"`
     // Added just to detect v2 openAPI specifications and to pretty print version error.
     Swagger [string](/builtin#string) `json:"swagger,omitempty" yaml:"swagger,omitempty"`
     // REQUIRED. Provides metadata about the API.
     //
     // The metadata MAY be used by tooling as required.
     Info Info `json:"info" yaml:"info"`
     // The default value for the `$schema` keyword within Schema Objects contained within this OAS document.
     JSONSchemaDialect [string](/builtin#string) `json:"jsonSchemaDialect,omitempty" yaml:"jsonSchemaDialect,omitempty"`
     // An array of Server Objects, which provide connectivity information to a target server.
     Servers []Server `json:"servers,omitempty" yaml:"servers,omitempty"`
     // The available paths and operations for the API.
     Paths Paths `json:"paths,omitempty" yaml:"paths,omitempty"`
     // The incoming webhooks that MAY be received as part of this API and that
     // the API consumer MAY choose to implement.
     //
     // Closely related to the `callbacks` feature, this section describes requests initiated other
     // than by an API call, for example by an out of band registration.
     //
     // The key name is a unique string to refer to each webhook, while the (optionally referenced)
     // PathItem Object describes a request that may be initiated by the API provider and the expected responses.
     Webhooks map[[string](/builtin#string)]*PathItem `json:"webhooks,omitempty" yaml:"webhooks,omitempty"`
     // An element to hold various schemas for the document.
     Components *Components `json:"components,omitempty" yaml:"components,omitempty"`
     // A declaration of which security mechanisms can be used across the API.
     // The list of values includes alternative security requirement objects that can be used.
     //
     // Only one of the security requirement objects need to be satisfied to authorize a request.
     //
     // Individual operations can override this definition.
     Security SecurityRequirements `json:"security,omitempty" yaml:"security,omitempty"`
    
     // A list of tags used by the specification with additional metadata.
     // The order of the tags can be used to reflect on their order by the parsing
     // tools. Not all tags that are used by the Operation Object must be declared.
     // The tags that are not declared MAY be organized randomly or based on the tools' logic.
     // Each tag name in the list MUST be unique.
     Tags []Tag `json:"tags,omitempty" yaml:"tags,omitempty"`
    
     // Additional external documentation.
     ExternalDocs *ExternalDocumentation `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`
    
     // Specification extensions.
     Extensions Extensions `json:"-" yaml:",inline"`
    
     // Raw YAML node. Used by '$ref' resolvers.
     Raw *[yaml](/github.com/go-faster/yaml).[Node](/github.com/go-faster/yaml#Node) `json:"-" yaml:"-"`
    }

Spec is the root document object of the OpenAPI document.

See <https://spec.openapis.org/oas/v3.1.0#openapi-object>.

#### func [NewSpec](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L15) ¶

    func NewSpec() *Spec

NewSpec returns a new Spec.

#### func [Parse](https://github.com/ogen-go/ogen/blob/v1.18.0/ogen.go#L9) ¶

    func Parse(data [][byte](/builtin#byte)) (s *Spec, err [error](/builtin#error))

Parse parses JSON/YAML into OpenAPI Spec.

#### func (*Spec) [AddNamedParameters](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L121) ¶

    func (s *Spec) AddNamedParameters(ps ...*NamedParameter) *Spec

AddNamedParameters adds the given namedParameters to the Components of the Spec.

#### func (*Spec) [AddNamedPathItems](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L63) ¶

    func (s *Spec) AddNamedPathItems(ps ...*NamedPathItem) *Spec

AddNamedPathItems adds the given namedPaths to the Paths of the Spec.

#### func (*Spec) [AddNamedRequestBodies](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L136) ¶

    func (s *Spec) AddNamedRequestBodies(scs ...*NamedRequestBody) *Spec

AddNamedRequestBodies adds the given namedRequestBodies to the Components of the Spec.

#### func (*Spec) [AddNamedResponses](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L106) ¶

    func (s *Spec) AddNamedResponses(scs ...*NamedResponse) *Spec

AddNamedResponses adds the given namedResponses to the Components of the Spec.

#### func (*Spec) [AddNamedSchemas](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L91) ¶

    func (s *Spec) AddNamedSchemas(scs ...*NamedSchema) *Spec

AddNamedSchemas adds the given namedSchemas to the Components of the Spec.

#### func (*Spec) [AddParameter](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L114) ¶

    func (s *Spec) AddParameter(n [string](/builtin#string), p *Parameter) *Spec

AddParameter adds the given Parameter under the given Name to the Components of the Spec.

#### func (*Spec) [AddPathItem](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L56) ¶

    func (s *Spec) AddPathItem(n [string](/builtin#string), p *PathItem) *Spec

AddPathItem adds the given PathItem under the given Name to the Paths of the Spec.

#### func (*Spec) [AddRequestBody](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L129) ¶

    func (s *Spec) AddRequestBody(n [string](/builtin#string), sc *RequestBody) *Spec

AddRequestBody adds the given RequestBody under the given Name to the Components of the Spec.

#### func (*Spec) [AddResponse](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L99) ¶

    func (s *Spec) AddResponse(n [string](/builtin#string), sc *Response) *Spec

AddResponse adds the given Response under the given Name to the Components of the Spec.

#### func (*Spec) [AddSchema](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L84) ¶

    func (s *Spec) AddSchema(n [string](/builtin#string), sc *Schema) *Spec

AddSchema adds the given Schema under the given Name to the Components of the Spec.

#### func (*Spec) [AddServers](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L40) ¶

    func (s *Spec) AddServers(srvs ...*Server) *Spec

AddServers adds Servers to the Servers of the Spec.

#### func (*Spec) [Init](https://github.com/ogen-go/ogen/blob/v1.18.0/spec.go#L125) ¶

    func (s *Spec) Init()

Init components of schema.

#### func (*Spec) [RefRequestBody](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L164) ¶

    func (s *Spec) RefRequestBody(n [string](/builtin#string)) *NamedRequestBody

RefRequestBody returns a new RequestBody referencing the given name.

#### func (*Spec) [RefResponse](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L154) ¶

    func (s *Spec) RefResponse(n [string](/builtin#string)) *NamedResponse

RefResponse returns a new Response referencing the given name.

#### func (*Spec) [RefSchema](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L144) ¶

    func (s *Spec) RefSchema(n [string](/builtin#string)) *NamedSchema

RefSchema returns a new Schema referencing the given name.

#### func (*Spec) [SetComponents](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L71) ¶

    func (s *Spec) SetComponents(c *Components) *Spec

SetComponents sets the Components of the Spec.

#### func (*Spec) [SetInfo](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L26) ¶

    func (s *Spec) SetInfo(i *Info) *Spec

SetInfo sets the Info of the Spec.

#### func (*Spec) [SetOpenAPI](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L20) ¶

    func (s *Spec) SetOpenAPI(v [string](/builtin#string)) *Spec

SetOpenAPI sets the OpenAPI Specification version of the document.

#### func (*Spec) [SetPaths](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L50) ¶

    func (s *Spec) SetPaths(p Paths) *Spec

SetPaths sets the Paths of the Spec.

#### func (*Spec) [SetServers](https://github.com/ogen-go/ogen/blob/v1.18.0/dsl.go#L34) ¶

    func (s *Spec) SetServers(srvs []Server) *Spec

SetServers sets the Servers of the Spec.

#### func (*Spec) [UnmarshalJSON](https://github.com/ogen-go/ogen/blob/v1.18.0/spec.go#L106) ¶ added in v0.22.0

    func (s *Spec) UnmarshalJSON(bytes [][byte](/builtin#byte)) [error](/builtin#error)

UnmarshalJSON implements json.Unmarshaler.

#### func (*Spec) [UnmarshalYAML](https://github.com/ogen-go/ogen/blob/v1.18.0/spec.go#L92) ¶ added in v0.44.0

    func (s *Spec) UnmarshalYAML(n *[yaml](/github.com/go-faster/yaml).[Node](/github.com/go-faster/yaml#Node)) [error](/builtin#error)

UnmarshalYAML implements yaml.Unmarshaler.

#### type [Tag](https://github.com/ogen-go/ogen/blob/v1.18.0/spec.go#L717) ¶

    type Tag struct {
     // REQUIRED. The name of the tag.
     Name [string](/builtin#string) `json:"name" yaml:"name"`
     // A description for the tag. CommonMark syntax MAY be used for rich text representation.
     Description [string](/builtin#string) `json:"description,omitempty" yaml:"description,omitempty"`
     // Additional external documentation for this tag.
     ExternalDocs *ExternalDocumentation `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`
    
     // Specification extensions.
     Extensions Extensions `json:"-" yaml:",inline"`
    }

Tag adds metadata to a single tag that is used by the Operation Object.

See <https://spec.openapis.org/oas/v3.1.0#tag-object>

#### type [XML](https://github.com/ogen-go/ogen/blob/v1.18.0/schema.go#L568) ¶ added in v0.44.0

    type XML struct {
     // Replaces the name of the element/attribute used for the described schema property.
     //
     // When defined within items, it will affect the name of the individual XML elements within the list.
     //
     // When defined alongside type being array (outside the items), it will affect the wrapping element
     // and only if wrapped is true.
     //
     // If wrapped is false, it will be ignored.
     Name [string](/builtin#string) `json:"name,omitempty" yaml:"name,omitempty"`
     // The URI of the namespace definition.
     //
     // This MUST be in the form of an absolute URI.
     Namespace [string](/builtin#string) `json:"namespace,omitempty" yaml:"namespace,omitempty"`
     // The prefix to be used for the name.
     Prefix [string](/builtin#string) `json:"prefix,omitempty" yaml:"prefix,omitempty"`
     // Declares whether the property definition translates to an attribute instead of an element.
     //
     // Default value is false.
     Attribute [bool](/builtin#bool) `json:"attribute,omitempty" yaml:"attribute,omitempty"`
     // MAY be used only for an array definition. Signifies whether the array is wrapped
     // (for example, `<books><book/><book/></books>`) or unwrapped (`<book/><book/>`).
     //
     // The definition takes effect only when defined alongside type being array (outside the items).
     //
     // Default value is false.
     Wrapped [bool](/builtin#bool) `json:"wrapped,omitempty" yaml:"wrapped,omitempty"`
    
     Common OpenAPICommon `json:"-" yaml:",inline"`
    }

XML is a metadata object that allows for more fine-tuned XML model definitions.

See <https://spec.openapis.org/oas/v3.1.0#xml-object>.

#### func (*XML) [ToJSONSchema](https://github.com/ogen-go/ogen/blob/v1.18.0/schema_backcomp.go#L140) ¶ added in v0.44.0

    func (d *XML) ToJSONSchema() *[jsonschema](/github.com/ogen-go/ogen@v1.18.0/jsonschema).[XML](/github.com/ogen-go/ogen@v1.18.0/jsonschema#XML)

ToJSONSchema converts XML to jsonschema.XML.
