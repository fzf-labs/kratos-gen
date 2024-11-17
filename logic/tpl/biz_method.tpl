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
{{- $s1 := "google.protobuf.Empty" }}
{{- if eq .Type 1 }}
// {{ .Name }} {{ .Comment }}
func ({{.FirstChar}} *{{ .UpperName }}UseCase) {{ .Name }}(ctx context.Context, req {{ if eq .Request $s1 }}*emptypb.Empty{{ else }}*pb.{{ .Request }}{{ end }}) ({{ if eq .Reply $s1 }}*emptypb.Empty{{ else }}*pb.{{ .Reply }}{{ end }}, error) {
	{{- if eq .Reply $s1 }}
 	resp :=	&emptypb.Empty{}
	{{- else }}
 	resp :=	&pb.{{ .Reply }}{}
	{{- end }}
	return resp, nil
}

{{- else if eq .Type 2 }}
// {{ .Name }} {{ .Comment }}
func ({{.FirstChar}} *{{ .UpperName }}UseCase) {{ .Name }}(conn pb.{{ .UpperName }}_{{ .UpperName }}Server) error {
	for {
		req, err := conn.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		err = conn.Send(&pb.{{ .Reply }}{})
		if err != nil {
			return err
		}
	}
}

{{- else if eq .Type 3 }}
// {{ .Name }} {{ .Comment }}
func ({{.FirstChar}} *{{ .UpperName }}UseCase) {{ .Name }}(conn pb.{{ .UpperName }}_{{ .UpperName }}Server) error {
	for {
		req, err := conn.Recv()
		if err == io.EOF {
			return conn.SendAndClose(&pb.{{ .Reply }}{})
		}
		if err != nil {
			return err
		}
	}
}

{{- else if eq .Type 4 }}
// {{ .Name }} {{ .Comment }}
func ({{.FirstChar}} *{{ .UpperName }}UseCase) {{ .Name }}(req {{ if eq .Request $s1 }}*emptypb.Empty
{{ else }}*pb.{{ .Request }}{{ end }}, conn pb.{{ .UpperName }}_{{ .UpperName }}Server) error {
	for {
		err := conn.Send(&pb.{{ .Reply }}{})
		if err != nil {
			return err
		}
	}
}

{{- end }}