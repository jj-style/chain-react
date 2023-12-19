package server

import (
	"fmt"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

func newRouter(devMode bool, corsValue string) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gzip.Gzip(gzip.DefaultCompression))
	// configure logging middleware
	r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("[%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			// param.ClientIP, <============ DON'T LOG USERS IP ADDRESS
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))

	// Configure CORS
	config := cors.DefaultConfig()
	config.ExposeHeaders = []string{"Access-Control-Allow-Origin"}
	if devMode {
		config.AllowOrigins = []string{"*"}
	} else {
		config.AllowOrigins = []string{corsValue}
	}
	r.Use(cors.New(config))
	return r
}

func (s *Server) setupRoutes() {
	api := s.Router.Group("/api")

	api.GET("/randomActor", s.handleGetRandomActor)
	api.GET("/randomActorNot/:id", s.handleGetRandomActorNotId)
	api.POST("/verify", s.handleVerify)
	api.POST("/verifyEdges", s.handleVerifyEdges)
	api.POST("/graph", s.handleGetGraph)
}
