package tmdb

import (
	go_tmdb "github.com/ryanbradynd05/go-tmdb"
	log "github.com/sirupsen/logrus"
)

type TMDb interface {
	GetAllActors() ([]*go_tmdb.Person, error)
}

type tmdb struct {
	client *go_tmdb.TMDb
	log    *log.Logger
}

func NewClient(api_key string) TMDb {
	config := go_tmdb.Config{
		APIKey: api_key,
	}
	client := go_tmdb.Init(config)
	return &tmdb{
		client: client,
		log:    log.StandardLogger(),
	}
}
