{{- /* delete empty line */ -}}
package data

import (
	"context"
)

func New{{ .upperName }}Repo(
	logger log.Logger,
	data *Data,
	{{ .lowerName }}Repo *{{ .dbName }}_repo.{{ .upperName }}Repo,
) *{{ .upperName }}Repo {
	l := log.NewHelper(log.With(logger, "module", "data/{{ .lowerName }}"))
	return &{{ .upperName }}Repo{
		log:         		l,
		data:     			data,
		{{ .upperName }}Repo:{{ .lowerName }}Repo,
	}
}

type {{ .upperName }}Repo struct {
	log *log.Helper
	data *Data
	*{{ .dbName }}_repo.{{ .upperName }}Repo
}