{{- /* delete empty line */ -}}
package service

import(
	"github.com/google/wire"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(
	{{ range .Wires }}
	{{- /* delete empty line */ -}}
	{{ . }},
	{{ end }}
)