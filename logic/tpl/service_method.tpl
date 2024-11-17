{{- /* delete empty line */ -}}
{{- $s1 := "google.protobuf.Empty" }}
{{- if eq .Type 1 }}
// {{ .Name }} {{ .Comment }}
func ({{.FirstChar}} *{{ .UpperName }}Service) {{ .Name }}(ctx context.Context, req {{ if eq .Request $s1 }}*emptypb.Empty{{ else }}*pb.{{ .Request }}{{ end }}) ({{ if eq .Reply $s1 }}*emptypb.Empty{{ else }}*pb.{{ .Reply }}{{ end }}, error) {
	return {{ if eq .Reply $s1 }}&emptypb.Empty{},nil{{ else }}{{.FirstChar}}.{{ .LowerName }}UseCase.{{ .Name }}(ctx, req){{ end }}
}

{{- else if eq .Type 2 }}
// {{ .Name }} {{ .Comment }}
func ({{.FirstChar}} *{{ .UpperName }}Service) {{ .Name }}(conn pb.{{ .UpperName }}_{{ .UpperName }}Server) error {
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
func ({{.FirstChar}} *{{ .UpperName }}Service) {{ .Name }}(conn pb.{{ .UpperName }}_{{ .UpperName }}Server) error {
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
func ({{.FirstChar}} *{{ .UpperName }}Service) {{ .Name }}(req {{ if eq .Request $s1 }}*emptypb.Empty
{{ else }}*pb.{{ .Request }}{{ end }}, conn pb.{{ .UpperName }}_{{ .UpperName }}Server) error {
	for {
		err := conn.Send(&pb.{{ .Reply }}{})
		if err != nil {
			return err
		}
	}
}
{{- end }}