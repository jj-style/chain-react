package db

import (
	"errors"
	"fmt"
	"log"

	"github.com/jj-style/chain-react/src/config"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
)

func NewRepository(conf *config.DbConfig) (Repository, error) {
	var repo Repository
	switch conf.Driver {
	case "neo4j":
		driver, err := neo4j.NewDriverWithContext(conf.Uri, neo4j.BasicAuth(conf.Username, conf.Password, ""))
		if err != nil {
			log.Fatalln("opening neo4j database driver: ", err)
		}
		repo = NewNeo4jRepository(driver)
	default:
		return nil, fmt.Errorf("unknown DB driver: %s", conf.Driver)
	}
	return repo, nil
}

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
	SearchActorName(name string) (*Actor, error)

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
