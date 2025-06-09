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
{{- $request := .Request }}
{{- $requestFields := slice .Messages 0 0 }}
{{- range $index, $message := .Messages }}
{{- if eq $message.Name $request }}
{{- $requestFields = $message.Fields }}
{{- end }}
{{- end }}

// {{ .Name }} {{ .Comment }}
func ({{.FirstChar}} *{{ .UpperServiceName }}Service) {{ .Name }}(ctx context.Context, req *pb.{{ .Request }}) (*pb.{{ .Reply }}, error) {
	resp := &pb.{{ .Reply }}{}
	data, err := {{.FirstChar}}.{{ .LowerName }}Repo.FindOneCacheByID(ctx, req.GetId())
	if err != nil {
		return nil, pb.ErrorReasonDataSQLError(pb.WithError(err))
	}
	if data == nil || data.ID == "" {
		return nil, pb.ErrorReasonDataRecordNotFound()
	}
	oldData := {{.FirstChar}}.{{ .LowerName }}Repo.DeepCopy(data)
	{{- if $requestFields }}
	{{- range $requestFields }}
	{{- if and (ne .Name "id") (ne .Name "ID") }}
	data.{{ .Name | ToCamel }} = req.Get{{ .Name | ToCamel }}()
	{{- end }}
	{{- end }}
	{{- else }}
	// No request fields detected
	{{- end }}
	err = {{.FirstChar}}.{{ .LowerName }}Repo.UpdateOneCacheWithZero(ctx, data, oldData)
	if err != nil {
		return nil, pb.ErrorReasonDataSQLError(pb.WithError(err))
	}
	return resp,nil
}