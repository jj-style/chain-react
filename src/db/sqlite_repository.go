package db

import (
	"database/sql"
	_ "embed"
	"errors"
	"io/fs"
	"path"

	"github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
)

var (
	ErrDuplicate    = errors.New("record already exists")
	ErrNotExists    = errors.New("row not exists")
	ErrDeleteFailed = errors.New("delete failed")
)

var (
	//go:embed ddl/actors/insert.sql
	INSERT_ACTOR_SQL string
	//go:embed ddl/actors/all.sql
	GET_ALL_ACTORS_SQL string
	//go:embed ddl/actors/delete.sql
	DELETE_ACTOR_SQL string
	//go:embed ddl/actors/latest.sql
	LATEST_ACTOR_SQL string
	//go:embed ddl/actors/random.sql
	RANDOM_ACTOR_SQL string
	//go:embed ddl/actors/random_not.sql
	RANDOM_ACTOR_NOT_ID_SQL string

	//go:embed ddl/movies/insert.sql
	INSERT_MOVIE_SQL string
	//go:embed ddl/movies/all.sql
	GET_ALL_MOVIES_SQL string
	//go:embed ddl/movies/delete.sql
	DELETE_MOVIE_SQL string

	//go:embed ddl/credits/insert.sql
	INSERT_CREDITS_SQL string
	//go:embed ddl/credits/all.sql
	GET_ALL_CREDITS_SQL string
)

type SQLiteRepository struct {
	db  *sql.DB
	log *log.Logger
}

func NewSQLiteRepository(db *sql.DB) *SQLiteRepository {
	return &SQLiteRepository{
		db:  db,
		log: log.StandardLogger(),
	}
}

func (r *SQLiteRepository) Migrate() error {
	it, err := migrationFs.ReadDir(path.Join("ddl", "migrations"))
	if err != nil {
		log.Fatalln("opening migration directory: ", err)
	}
	for _, dirent := range it {
		if !dirent.IsDir() {
			continue
		}
		up_migration, err := fs.ReadFile(migrationFs, path.Join("ddl", "migrations", dirent.Name(), "up.sql"))
		if err != nil {
			log.Fatalf("reading migration script from %s: %v", dirent.Name(), err)
		}
		if _, err := r.db.Exec(string(up_migration)); err != nil {
			log.Fatalf("executing migration from %s: %v", dirent.Name(), err)
		}
	}
	return nil
}

func (r *SQLiteRepository) CreateActor(actor Actor) (*Actor, error) {
	_, err := r.db.Exec(INSERT_ACTOR_SQL, actor.Id, actor.Name)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {
			if errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintUnique) || errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintPrimaryKey) {
				return nil, ErrDuplicate
			}
		}
		return nil, err
	}

	return &actor, nil
}

func (r *SQLiteRepository) AllActors() ([]Actor, error) {
	rows, err := r.db.Query(GET_ALL_ACTORS_SQL)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var all []Actor
	for rows.Next() {
		var actor Actor
		if err := rows.Scan(&actor.Id, &actor.Name); err != nil {
			return nil, err
		}
		all = append(all, actor)
	}
	return all, nil
}

func (r *SQLiteRepository) DeleteActor(id int64) error {
	res, err := r.db.Exec(DELETE_ACTOR_SQL, id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrDeleteFailed
	}

	return err
}

func (r *SQLiteRepository) LatestActor() (*Actor, error) {
	row := r.db.QueryRow(LATEST_ACTOR_SQL)

	var actor Actor
	if err := row.Scan(&actor.Id, &actor.Name); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotExists
		}
		return nil, err
	}
	return &actor, nil

}

func (r *SQLiteRepository) RandomActor() (*Actor, error) {
	row := r.db.QueryRow(RANDOM_ACTOR_SQL)

	var actor Actor
	if err := row.Scan(&actor.Id, &actor.Name); err != nil {
		return nil, err
	}
	return &actor, nil
}

func (r *SQLiteRepository) RandomActorNotId(id int) (*Actor, error) {
	row := r.db.QueryRow(RANDOM_ACTOR_NOT_ID_SQL, id)

	var actor Actor
	if err := row.Scan(&actor.Id, &actor.Name); err != nil {
		return nil, err
	}
	return &actor, nil
}

func (r *SQLiteRepository) CreateMovie(movie Movie) (*Movie, error) {
	_, err := r.db.Exec(INSERT_MOVIE_SQL, movie.Id, movie.Title)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {
			if errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintUnique) || errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintPrimaryKey) {
				return nil, ErrDuplicate
			}
		}
		return nil, err
	}

	return &movie, nil
}

func (r *SQLiteRepository) AllMovies() ([]Movie, error) {
	rows, err := r.db.Query(GET_ALL_MOVIES_SQL)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var all []Movie
	for rows.Next() {
		var movie Movie
		if err := rows.Scan(&movie.Id, &movie.Title); err != nil {
			return nil, err
		}
		all = append(all, movie)
	}
	return all, nil
}

func (r *SQLiteRepository) DeleteMovie(id int64) error {
	res, err := r.db.Exec(DELETE_MOVIE_SQL, id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrDeleteFailed
	}

	return err
}

func (r *SQLiteRepository) CreateCredit(credit CreditIn) (*CreditIn, error) {
	_, err := r.db.Exec(INSERT_CREDITS_SQL, credit.ActorId, credit.MovieId, credit.CreditId, credit.Character)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {
			if errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintUnique) || errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintPrimaryKey) || errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintForeignKey) {
				return nil, ErrDuplicate
			}
		}
		return nil, err
	}

	return &credit, nil
}

func (r *SQLiteRepository) AllCredits() ([]Credit, error) {
	rows, err := r.db.Query(GET_ALL_CREDITS_SQL)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var all []Credit
	for rows.Next() {
		var credit Credit
		if err := rows.Scan(&credit.Actor.Id, &credit.Actor.Name, &credit.Movie.Id, &credit.Movie.Title, &credit.CreditId, &credit.Character); err != nil {
			return nil, err
		}
		all = append(all, credit)
	}
	return all, nil
}
