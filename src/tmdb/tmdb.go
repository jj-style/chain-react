package tmdb

import (
	"context"

	go_tmdb "github.com/ryanbradynd05/go-tmdb"
	log "github.com/sirupsen/logrus"
)

type TMDb interface {
	GetAllActors(ctx context.Context, c chan<- *go_tmdb.Person, r func() error) error
	GetActorsFrom(ctx context.Context, c chan<- *go_tmdb.Person, r func() error, id int) error
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
