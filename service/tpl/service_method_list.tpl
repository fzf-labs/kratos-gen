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
	resp := &pb.{{ .Reply }}{
		Total: 0,
		List:  []*pb.{{ .UpperName }}Info{},
	}
	param := &condition.Req{
		Page:     req.GetPage(),
		PageSize: req.GetPageSize(),
		Query:    []*condition.QueryParam{},
		Order: []*condition.OrderParam{
			{
				Field: "created_at",
				Order: condition.DESC,
			},
		},
	}
	list, p, err := {{.FirstChar}}.{{ .LowerName }}Repo.FindMultiCacheByCondition(ctx, param)
	if err != nil {
		return nil, pb.ErrorReasonDataSQLError(pb.WithError(err))
	}
	resp.Total = int32(p.Total)
	if len(list) > 0 {
		for _, v := range list {
			resp.List = append(resp.List, &pb.{{ .UpperName }}Info{
				Id: v.ID,
				{{- range $field := $infoFields }}
				{{- if ne $field.Name "id" }}
				{{ $field.Name | UCFirst }}: v.{{ $field.Name | UCFirst }},
				{{- end }}
				{{- end }}
			})
		}
	}
	return resp,nil
}