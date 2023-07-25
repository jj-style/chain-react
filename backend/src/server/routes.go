package server

import (
	"fmt"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func setupRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	// LoggerWithFormatter middleware will write the logs to gin.DefaultWriter
	// By default gin.DefaultWriter = os.Stdout
	r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// your custom format
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
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
	// - No origin allowed by default
	// - GET,POST, PUT, HEAD methods
	// - Credentials share disabled
	// - Preflight requests cached for 12 hours
	config := cors.DefaultConfig()
	config.ExposeHeaders = []string{"Access-Control-Allow-Origin"}
	if viper.GetBool("devMode") {
		config.AllowOrigins = []string{"*"}
	} else {
		config.AllowOrigins = []string{"http://localhost:3000"}
	}
	r.Use(cors.New(config))
	return r
}

func (s *Server) routes() {
	s.Router.GET("/randomActor", s.handleGetRandomActor)
	s.Router.GET("/randomActorNot/:id", s.handleGetRandomActorNotId)
	s.Router.POST("/verify", s.handleVerify)
	s.Router.POST("/verifyEdges", s.handleVerifyEdges)
}
