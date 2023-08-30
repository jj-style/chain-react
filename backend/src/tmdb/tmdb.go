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

type TMDbClient interface {
	GetPersonLatest() (*go_tmdb.PersonLatest, error)
	SearchPerson(name string, options map[string]string) (*go_tmdb.PersonSearchResults, error)
	GetPersonInfo(id int, options map[string]string) (*go_tmdb.Person, error)
	GetPersonMovieCredits(id int, options map[string]string) (*go_tmdb.PersonMovieCredits, error)
}

type tmdb struct {
	client TMDbClient
	log    *log.Logger
}

func NewClient(client TMDbClient, logger *log.Logger) TMDb {

	return &tmdb{
		client: client,
		log:    log.StandardLogger(),
	}
}
