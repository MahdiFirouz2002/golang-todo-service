package docs

import "embed"

// OpenAPI specification embedded for runtime serving.
//
//go:embed openapi.yaml
var OpenAPI embed.FS
