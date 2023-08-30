package config

import (
	"github.com/jj-style/chain-react/src/db"
	"github.com/jj-style/chain-react/src/search"
	"github.com/jj-style/chain-react/src/tmdb"
	go_tmdb "github.com/jj-style/go-tmdb"
	_ "github.com/mattn/go-sqlite3"
	"github.com/meilisearch/meilisearch-go"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type RConfig struct {
	Log    *log.Logger
	Repo   db.Repository
	Tmdb   tmdb.TMDb
	Search search.Repository
	Redis  *redis.Client
}

func NewRuntimeConfig(c *Config) RConfig {
	var repo db.Repository
	switch c.Db.Driver {
	case "neo4j":
		driver, err := neo4j.NewDriverWithContext(c.Db.Uri, neo4j.BasicAuth(c.Db.Username, c.Db.Password, ""))
		if err != nil {
			log.Fatalln("opening neo4j database driver: ", err)
		}
		repo = db.NewNeo4jRepository(driver)
	default:
		log.Fatalf("unknown DB driver: %s", c.Db.Driver)
	}

	logger := log.StandardLogger()
	loglvl := log.InfoLevel
	if p, err := log.ParseLevel(viper.GetString("log.level")); err == nil {
		loglvl = p
	} else {
		log.Fatalf("error parsing log level: %v", err.Error())
	}
	logger.SetLevel(loglvl)

	tcfg := go_tmdb.Config{
		APIKey: c.Tmdb.ApiKey,
	}
	t := tmdb.NewClient(go_tmdb.Init(tcfg), logger)

	m := meilisearch.NewClient(meilisearch.ClientConfig{
		Host:   c.Meilisearch.Host,
		APIKey: c.Meilisearch.ApiKey,
	})
	searchRepo := search.NewMeilisearchRepository(m)

	redisClient := redis.NewClient(&redis.Options{
		Addr:     c.Redis.Address,
		Password: c.Redis.Password,
		DB:       c.Redis.Db,
	})

	return RConfig{
		Log:    logger,
		Repo:   repo,
		Tmdb:   t,
		Search: &searchRepo,
		Redis:  redisClient,
	}
}
