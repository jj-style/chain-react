package config

import (
	"database/sql"

	"github.com/jj-style/chain-react/src/db"
	"github.com/jj-style/chain-react/src/search"
	"github.com/jj-style/chain-react/src/tmdb"
	_ "github.com/mattn/go-sqlite3"
	"github.com/meilisearch/meilisearch-go"
	log "github.com/sirupsen/logrus"
)

type RConfig struct {
	Log    *log.Logger
	Repo   db.Repository
	Tmdb   tmdb.TMDb
	Search search.Repository
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

	m := meilisearch.NewClient(meilisearch.ClientConfig{
		Host:   c.Meilisearch.Host,
		APIKey: c.Meilisearch.ApiKey,
	})
	searchRepo := search.NewMeilisearchRepository(m)

	return RConfig{
		Log:    log.StandardLogger(),
		Repo:   repo,
		Tmdb:   t,
		Search: &searchRepo,
	}
}
