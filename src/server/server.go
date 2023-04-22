package server

import (
	"github.com/gin-gonic/gin"
	"github.com/jj-style/chain-react/src/config"
	"github.com/jj-style/chain-react/src/graph"
	graph2 "github.com/jj-style/chain-react/src/server/graph"
	log "github.com/sirupsen/logrus"
)

type Server struct {
	Router *gin.Engine
	Config *config.RConfig
	Log    *log.Logger
	Graph  *graph.Graph[graph2.Vertex, graph2.Edge]
}

func NewServer(config *config.RConfig) *Server {
	s := &Server{
		Router: setupRouter(),
		Config: config,
		Log:    log.StandardLogger(),
		Graph:  graph2.LoadGraph(config),
	}
	s.routes()
	return s
}
