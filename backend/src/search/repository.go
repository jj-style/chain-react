package search

import (
	"github.com/jj-style/chain-react/src/config"
	"github.com/meilisearch/meilisearch-go"
	log "github.com/sirupsen/logrus"
)

type Repository interface {
	AddDocuments(docs interface{}, index string) error
}

func NewRepository(conf *config.MeilisearchConfig, logger *log.Logger) Repository {
	m := meilisearch.New(conf.Host, meilisearch.WithAPIKey(conf.ApiKey))
	meili := NewMeilisearchRepository(m, logger)
	return &meili
}
