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
{{- $s1 := "google.protobuf.Empty" }}
{{- $upperName := .UpperName }}
{{- $infoName := printf "%sInfo" $upperName }}
{{- $infoFields := slice .Messages 0 0 }}
{{- range $index, $message := .Messages }}
{{- if eq $message.Name $infoName }}
{{- $infoFields = $message.Fields }}
{{- end }}
{{- end }}

// {{ .Name }} {{ .Comment }}
func ({{.FirstChar}} *{{ .UpperServiceName }}Service) {{ .Name }}(ctx context.Context, req *pb.{{ .Request }}) (*pb.{{ .Reply }}, error) {
	resp := &pb.{{ .Reply }}{}
	data, err := {{.FirstChar}}.{{ .LowerName }}Repo.FindOneCacheByID(ctx, req.GetId())
	if err != nil {
		return nil, pb.ErrorReasonDataSQLError(pb.WithError(err))
	}
	resp.Info = &pb.{{ .UpperName }}Info{
		Id: data.ID,
		{{- range $field := $infoFields }}
		{{- if ne $field.Name "id" }}
		{{ $field.Name | ToCamel }}: data.{{ $field.Name | ToCamel }},
		{{- end }}
		{{- end }}
	}
	return resp,nil
}