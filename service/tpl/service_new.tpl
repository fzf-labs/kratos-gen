{{- /* delete empty line */ -}}
package service

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

func New{{ .UpperServiceName }}Service(
	logger log.Logger,
	{{ .LowerName }}Repo *data.{{ .UpperName }}Repo,
) *{{ .UpperServiceName }}Service {
	l := log.NewHelper(log.With(logger, "module", "service/{{ .LowerName }}"))
	return &{{ .UpperServiceName }}Service{
		log:         l,
		{{ .LowerName }}Repo: {{ .LowerName }}Repo,
	}
}

type {{ .UpperServiceName }}Service struct {
	pb.Unimplemented{{ .UpperName }}Server
	log *log.Helper
	{{ .LowerName }}Repo *data.{{ .UpperName }}Repo
}
