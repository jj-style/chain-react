package server

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jj-style/chain-react/src/db"
	"github.com/jj-style/chain-react/src/graph"
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

func (s *Server) handleVerify(c *gin.Context) {
	type request struct {
		Chain []int `json:"chain"`
	}
	type response struct {
		Valid bool   `json:"valid"`
		Error string `json:"error"`
		Chain []struct {
			Src  db.Credit `json:"src"`
			Dest db.Credit `json:"dest"`
		} `json:"chain"`
	}

	var req request
	err := c.BindJSON(&req)
	if err != nil {
		if errors.Is(err, graph.ErrChainLength) {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid request: %v", err)})
		} else {
			c.JSON(http.StatusOK, response{
				Valid: false,
				Error: err.Error(),
			})
		}
		return
	}

	err = s.Graph.Verify(graph.Chain(req.Chain))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid request: %v", err)})
		return
	}

	c.JSON(http.StatusOK, response{Valid: true})
}
