package server

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/jj-style/chain-react/src/config"
	"github.com/jj-style/chain-react/src/db"
	redis2 "github.com/jj-style/chain-react/src/redis"
	"github.com/jj-style/chain-react/src/search"
	gamemanager "github.com/jj-style/chain-react/src/server/game_manager"
	"github.com/jj-style/chain-react/src/tmdb"
)

var ProviderSet = wire.NewSet(db.NewRepository, tmdb.NewTMDb, search.NewRepository, redis2.NewRedis, gamemanager.NewCronGameManager, NewServer)

type Server struct {
	Router      *gin.Engine
	Log         *log.Logger
	Repo        db.Repository
	Tmdb        tmdb.TMDb
	Search      search.Repository
	Cache       *redis.Client
	GameManager gamemanager.GameManager
}

func NewServer(logger *log.Logger, conf *config.Server, repo db.Repository, tmdb tmdb.TMDb, search search.Repository, redis *redis.Client, gameManager gamemanager.GameManager) *Server {

	s := &Server{
		Router:      newRouter(viper.GetBool("devMode"), conf.Cors),
		Log:         logger,
		Repo:        repo,
		Tmdb:        tmdb,
		Search:      search,
		Cache:       redis,
		GameManager: gameManager,
	}
	s.setupRoutes()
	return s
}
