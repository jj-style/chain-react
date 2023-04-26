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
	type responseEdge struct {
		Src  db.Credit `json:"src"`
		Dest db.Credit `json:"dest"`
	}
	type response struct {
		Valid bool           `json:"valid"`
		Error string         `json:"error"`
		Chain []responseEdge `json:"chain"`
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

	edges, err := s.Graph.VerifyWithEdges(graph.Chain(req.Chain))
	edgeResponses := make([]responseEdge, 0, len(edges))
	for _, e := range edges {
		edgeResponses = append(edgeResponses, responseEdge{
			Src:  e.Weight.Src,
			Dest: e.Weight.Dest,
		})
	}

	var errs = ""
	var code = http.StatusOK
	if err != nil {
		errs = err.Error()
		code = http.StatusBadRequest
	}

	c.JSON(code, response{
		Valid: true,
		Chain: edgeResponses,
		Error: errs,
	})
}
