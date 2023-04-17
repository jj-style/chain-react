package config

import (
	"database/sql"

	"github.com/jj-style/chain-react/src/db"
	"github.com/jj-style/chain-react/src/tmdb"
	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
)

type RConfig struct {
	Log  *log.Logger
	Repo db.Repository
	Tmdb tmdb.TMDb
}

func NewRuntimeConfig(c *Config) RConfig {
	dbs, err := sql.Open(c.Db.Driver, c.Db.Uri)
	if err != nil {
		log.Fatalln(err)
	}

	_, err = dbs.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		log.Fatalln(err)
	}

	repo := db.NewSQLiteRepository(dbs)
	if err := repo.Migrate(); err != nil {
		log.Fatalln("migrating db: ", err)
	}

	t := tmdb.NewClient(c.Tmdb.ApiKey)

	return RConfig{
		Log:  log.StandardLogger(),
		Repo: repo,
		Tmdb: t,
	}
}
