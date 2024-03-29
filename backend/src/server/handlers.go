package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jj-style/chain-react/src/db"

	log "github.com/sirupsen/logrus"
)

func (s *Server) handleGetRandomActor(c *gin.Context) {
	actor, err := s.Repo.RandomActor()
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

	actor, err := s.Repo.RandomActorNotId(int(id))
	if err != nil {
		s.Log.WithField("error", err).Error("getting random actor from db")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to get random actor: %v", err)})
		return
	}
	s.Log.WithField("actor", actor).Debug("got random actor from db")
	c.JSON(http.StatusOK, actor)
}

func (s *Server) handleVerifyEdges(c *gin.Context) {
	type request struct {
		db.Chain `json:"chain"`
	}
	type response struct {
		Valid bool       `json:"valid"`
		Error string     `json:"error"`
		Chain []*db.Edge `json:"chain"`
	}

	var req request
	if err := c.BindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	edges, err := s.Repo.VerifyWithEdges(req.Chain)

	var errs = ""
	var code = http.StatusOK
	var valid = true
	if err != nil {
		errs = err.Error()
		code = http.StatusBadRequest
		valid = false
	}

	c.JSON(code, response{
		Valid: valid,
		Chain: edges,
		Error: errs,
	})
}

func (s *Server) handleVerify(c *gin.Context) {
	type request struct {
		db.Chain `json:"chain"`
	}
	type response struct {
		Valid bool   `json:"valid"`
		Error string `json:"error"`
	}

	var req request
	if err := c.BindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	valid, err := s.Repo.Verify(req.Chain)

	var errs = ""
	var code = http.StatusOK
	if err != nil {
		errs = err.Error()
		code = http.StatusBadRequest
		valid = false
	}

	c.JSON(code, response{
		Valid: valid,
		Error: errs,
	})
}

func (s *Server) handleGetGraph(c *gin.Context) {
	type request struct {
		Chain  []int `json:"chain" binding:"required"`
		Length int   `json:"length"`
	}
	var req request
	if err := c.BindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var length = 4
	if req.Length != 0 {
		// TODO - fix query otherwise this blows things up
		if req.Length != 4 {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("length '%d' not valid option", req.Length)})
			return
		}
		length = req.Length
	}

	g, err := s.Repo.GetGraph(length, req.Chain...)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": g})
}

func (s *Server) handleGetManagedGame(c *gin.Context) {
	game, err := s.GameManager.GetGame(c.Request.Context())
	if err != nil {
		s.Log.Errorf("getting current game from manager: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, game)
}
