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

func New{{ .UpperName }}Service(
	logger log.Logger,
	{{ .LowerName }}UseCase *biz.{{ .UpperName }}UseCase,
) *{{ .UpperName }}Service {
	l := log.NewHelper(log.With(logger, "module", "service/{{ .LowerName }}"))
	return &{{ .UpperName }}Service{
		log:         l,
		{{ .LowerName }}UseCase:{{ .LowerName }}UseCase,
	}
}

type {{ .UpperName }}Service struct {
	pb.Unimplemented{{ .UpperServiceName }}Server
	log *log.Helper
	{{ .LowerName }}UseCase *biz.{{ .UpperName }}UseCase
}
