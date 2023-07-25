package tmdb

import (
	"context"

	go_tmdb "github.com/jj-style/go-tmdb"
	log "github.com/sirupsen/logrus"
)

type TMDb interface {
	// Get all actors up to the latest actor in TMDB database.
	// Calls the reducer function `r` for each actor successfully retrieved.
	GetAllActors(ctx context.Context, c chan<- *go_tmdb.Person, r func() error) error
	// Get all actors from the id given up to the latest actor in TMDB database from.
	// Calls the reducer function `r` for each actor successfully retrieved.
	GetActorsFrom(ctx context.Context, c chan<- *go_tmdb.Person, r func() error, id int) error
	// TODO - add GetActorByID
	GetActorsByName(ctx context.Context, c chan<- *go_tmdb.Person, r func() error, names ...string) error
	// Gets the movie credits for the actor given by the id.
	GetActorMovieCredits(ctx context.Context, id int) (*go_tmdb.PersonMovieCredits, error)
	GetAllActorMovieCredits(ctx context.Context, c chan<- *go_tmdb.PersonMovieCredits, r func() error, ids ...int) error
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
