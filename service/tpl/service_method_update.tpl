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
	// TODO
	err = {{.FirstChar}}.{{ .LowerName }}Repo.UpdateOneCache(ctx, data, oldData)
	if err != nil {
		return nil, pb.ErrorReasonDataSQLError(pb.WithError(err))
	}
	return resp,nil
}