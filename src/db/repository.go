package db

type Repository interface {
	// General
	Migrate() error

	// Actors
	CreateActor(actor Actor) (*Actor, error)
	AllActors() ([]Actor, error)
	DeleteActor(id int64) error
	LatestActor() (*Actor, error)

	// Movies
	CreateMovie(movie Movie) (*Movie, error)
	AllMovies() ([]Movie, error)
	DeleteMovie(id int64) error

	// Credits
	CreateCredit(credit CreditIn) (*CreditIn, error)
	AllCredits() ([]Credit, error)
}
