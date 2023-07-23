package search

import (
	"github.com/meilisearch/meilisearch-go"
	log "github.com/sirupsen/logrus"
)

type MeilisearchRepository struct {
	client *meilisearch.Client
	log    *log.Logger
}

func NewMeilisearchRepository(c *meilisearch.Client) MeilisearchRepository {
	return MeilisearchRepository{
		client: c,
		log:    log.StandardLogger(),
	}
}

func (m *MeilisearchRepository) AddDocuments(docs interface{}, index string) error {
	idx := m.client.Index(index)
	task, err := idx.AddDocuments(docs)
	if err != nil {
		return err
	}
	m.log.WithFields(log.Fields{"index": index, "task": task}).Info("indexing documents to meilisearch")
	return nil
}
