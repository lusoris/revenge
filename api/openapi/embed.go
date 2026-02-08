// Package openapi provides the embedded OpenAPI specification.
package openapi

import _ "embed"

//go:embed openapi.yaml
var Spec []byte
