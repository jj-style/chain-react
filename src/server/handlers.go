package server

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jj-style/chain-react/src/db"
	"github.com/jj-style/chain-react/src/graph"
	graph2 "github.com/jj-style/chain-react/src/server/graph"

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
	var valid = true
	if err != nil {
		errs = err.Error()
		code = http.StatusBadRequest
		valid = false
	}

	c.JSON(code, response{
		Valid: valid,
		Chain: edgeResponses,
		Error: errs,
	})
}

func (s *Server) handleChains(c *gin.Context) {
	type request struct {
		Start int `json:"start"`
		End   int `json:"end"`
	}
	type responseEdge struct {
		Src  db.Credit `json:"src"`
		Dest db.Credit `json:"dest"`
	}
	type response struct {
		Valid  bool             `json:"valid"`
		Error  string           `json:"error"`
		Chains [][]responseEdge `json:"chains"`
	}

	var req request
	err := c.BindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid request: %v", err)})
		return
	}

	resp := response{
		Chains: make([][]responseEdge, 0),
	}

	s.Log.WithFields(log.Fields{"start": req.Start, "end": req.End}).Debug("searching for all chains")

	filterNeighbour := func(path []graph.Element[graph2.Vertex, graph2.Edge], e *graph.Edge[graph2.Vertex, graph2.Edge]) bool {
		if len(path) <= 1 {
			return false
		}
		if len(path) > 3 {
			return true
		}
		last := path[len(path)-1]
		return last.Edge.Weight.Dest.Movie.Id == e.Weight.Src.Movie.Id
	}

	found := make(chan []graph.Element[graph2.Vertex, graph2.Edge])
	go s.Graph.BfsWithNeighbourFilter(req.Start, req.End, found, filterNeighbour)
	for p := range found {
		edgeResponses := make([]responseEdge, 0, len(p))
		for _, e := range p[1:] {
			edgeResponses = append(edgeResponses, responseEdge{
				Src:  e.Edge.Weight.Src,
				Dest: e.Edge.Weight.Dest,
			})
		}
		resp.Chains = append(resp.Chains, edgeResponses)
	}

	if len(resp.Chains) == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid chain from %d to %d", req.Start, req.End)})
		return
	}

	c.JSON(http.StatusOK, resp)
}
