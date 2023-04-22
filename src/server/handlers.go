package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func (s *Server) handleGetRandomActor(c *gin.Context) {
	actor, err := s.Config.Repo.RandomActor()
	if err != nil {
		s.Log.WithField("error", err).Error("getting random actor from db")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to get random actor: %v", err)})
		return
	}
	s.Log.WithField("actor", actor).Debug("got random actor from db")
	c.JSON(http.StatusOK, actor)
}

func (s *Server) handleGetRandomActorNotId(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		s.Log.WithFields(log.Fields{"id": idStr, "error": err}).Error("could not parse id to int")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("failed to parse id: %v", err)})
		return
	}

	actor, err := s.Config.Repo.RandomActorNotId(int(id))
	if err != nil {
		s.Log.WithField("error", err).Error("getting random actor from db")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to get random actor: %v", err)})
		return
	}
	s.Log.WithField("actor", actor).Debug("got random actor from db")
	c.JSON(http.StatusOK, actor)
}
