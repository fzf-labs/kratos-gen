{{- /* delete empty line */ -}}
package biz

import(
	"github.com/google/wire"
)

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(
	{{ range .Wires }}
	{{- /* delete empty line */ -}}
	{{ . }},
	{{ end }}
)