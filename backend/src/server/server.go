package server

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/jj-style/chain-react/src/config"
)

type Server struct {
	Router *gin.Engine
	Config *config.RConfig
	Log    *log.Logger
}

func NewServer(config *config.RConfig) *Server {
	logger := log.StandardLogger()
	if lvl, err := log.ParseLevel(viper.GetString("LOG_LEVEL")); err == nil {
		logger.SetLevel(lvl)
	}

	s := &Server{
		Router: setupRouter(),
		Config: config,
		Log:    logger,
	}
	s.routes()
	return s
}
