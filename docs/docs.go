// Package docs provides Swagger documentation for the Nyumbani Care API.
package docs

import "embed"

//go:embed swagger.json
var SwaggerJSON embed.FS
