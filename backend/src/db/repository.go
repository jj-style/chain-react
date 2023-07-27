package db

import (
	"errors"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
)

type Repository interface {
	// General
	Migrate() error

	// Actors
	CreateActor(actor Actor) (*Actor, error)
	AllActors() ([]Actor, error)
	DeleteActor(id int64) error
	LatestActor() (*Actor, error)
	RandomActor() (*Actor, error)
	RandomActorNotId(int) (*Actor, error)

	// Movies
	CreateMovie(movie Movie) (*Movie, error)
	AllMovies() ([]Movie, error)
	DeleteMovie(id int64) error

	// Credits
	CreateCredit(credit CreditIn) (*CreditIn, error)
	AllCredits() ([]Credit, error)

	// Graph Functions
	Verify(c Chain) (bool, error)
	VerifyWithEdges(c Chain) ([]*Edge, error)
	GetGraph(length int, nodes ...int) ([]dbtype.Path, error)
}

var (
	ErrNotExists error = errors.New("not exist")
	ErrDuplicate error = errors.New("unique constraint")
)
