package server

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/jj-style/chain-react/src/config"
	"github.com/jj-style/chain-react/src/db"
	"github.com/jj-style/chain-react/src/search"
	"github.com/jj-style/chain-react/src/tmdb"
)

var ProviderSet = wire.NewSet(db.NewRepository, tmdb.NewTMDb, search.NewRepository, NewServer)

type Server struct {
	Router *gin.Engine
	Log    *log.Logger
	Repo   db.Repository
	Tmdb   tmdb.TMDb
	Search search.Repository
	Cache  *redis.Client
}

func NewServer(logger *log.Logger, repo db.Repository, tmdb tmdb.TMDb, search search.Repository, rconf *config.RedisConfig) *Server {
	if lvl, err := log.ParseLevel(viper.GetString("LOG_LEVEL")); err == nil {
		logger.SetLevel(lvl)
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     rconf.Address,
		Password: rconf.Password,
		DB:       rconf.Db,
	})

	s := &Server{
		Router: newRouter(viper.GetBool("devMode")),
		Log:    logger,
		Repo:   repo,
		Tmdb:   tmdb,
		Search: search,
		Cache:  redisClient,
	}
	s.setupRoutes()
	return s
}
