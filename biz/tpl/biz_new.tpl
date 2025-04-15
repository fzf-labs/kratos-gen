{{- /* delete empty line */ -}}
package biz

import (
	{{- if .UseContext }}
	"context"
	{{- end }}
	{{- if .UseIO }}
	"io"
	{{- end }}
	pb "{{ .GoPackage }}"
	{{- if .GoogleEmpty }}
	"google.golang.org/protobuf/types/known/emptypb"
	{{- end }}
)

func New{{ .UpperName }}UseCase(
	logger log.Logger,
) *{{ .UpperName }}UseCase {
	l := log.NewHelper(log.With(logger, "module", "biz/{{ .LowerName }}"))
	return &{{ .UpperName }}UseCase{
		log:         l,
	}
}

type {{ .UpperName }}UseCase struct {
	log *log.Helper
}