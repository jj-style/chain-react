package db

import (
	"database/sql"
	"errors"

	"github.com/mattn/go-sqlite3"
)

var (
	ErrDuplicate    = errors.New("record already exists")
	ErrNotExists    = errors.New("row not exists")
	ErrDeleteFailed = errors.New("delete failed")
)

type SQLiteRepository struct {
	db *sql.DB
}

func NewSQLiteRepository(db *sql.DB) *SQLiteRepository {
	return &SQLiteRepository{
		db: db,
	}
}

func (r *SQLiteRepository) Migrate() error {
	query := `
    CREATE TABLE IF NOT EXISTS actors(
        id INTEGER PRIMARY KEY,
        name TEXT NOT NULL UNIQUE
    );
    `

	_, err := r.db.Exec(query)
	return err
}

func (r *SQLiteRepository) CreateActor(actor Actor) (*Actor, error) {
	_, err := r.db.Exec("INSERT INTO actors(id, name) values(?,?)", actor.Id, actor.Name)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {
			if errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintUnique) {
				return nil, ErrDuplicate
			}
		}
		return nil, err
	}

	return &actor, nil
}

func (r *SQLiteRepository) AllActors() ([]Actor, error) {
	rows, err := r.db.Query("SELECT * FROM actors")
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
	res, err := r.db.Exec("DELETE FROM actors WHERE id = ?", id)
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
	row := r.db.QueryRow("SELECT * FROM actors ORDER BY id DESC LIMIT 1")

	var actor Actor
	if err := row.Scan(&actor.Id, &actor.Name); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotExists
		}
		return nil, err
	}
	return &actor, nil

}
