{{- /* delete empty line */ -}}
package data

import (
	"context"
)

var _ biz.{{ .upperName }}Repo = (*{{ .upperName }}Repo)(nil)

func New{{ .upperName }}Repo(
	logger log.Logger,
	data *Data,
) biz.{{ .upperName }}Repo {
	l := log.NewHelper(log.With(logger, "module", "data/{{ .lowerName }}"), log.WithMessageKey("message"))
	return &{{ .upperName }}Repo{
		log:         l,
		data:     data,
	}
}

type {{ .upperName }}Repo struct {
	log *log.Helper
	data *Data
}