package graph

import (
	"log"

	"github.com/jj-style/chain-react/src/config"
	"github.com/jj-style/chain-react/src/db"
	"github.com/jj-style/chain-react/src/graph"
	"github.com/samber/lo"
)

type Vertex db.Actor
type Edge struct {
	Src  db.Credit
	Dest db.Credit
}

func LoadGraph(cfg *config.RConfig) *graph.Graph[Vertex, Edge] {
	credits, err := cfg.Repo.AllCredits()
	if err != nil {
		log.Fatalln("getting all credits: ", err)
	}

	g := graph.NewGraph[Vertex, Edge]()

	// add all actors as vertices in graph
	actors, err := cfg.Repo.AllActors()
	if err != nil {
		log.Fatalln("getting all actors: ", err)
	}
	for _, a := range actors {
		g.AddVertex(a.Id, Vertex(a))
	}

	creditsByMovieId := lo.GroupBy(credits, func(c db.Credit) int { return c.Movie.Id })
	for _, c := range credits {
		// build edge between actor in credit, and every other actor in that movie
		a := c.Actor
		for _, oc := range creditsByMovieId[c.Movie.Id] {
			if oc.Actor == a {
				// avoid cycles
				continue
			}
			g.AddEdge(a.Id, oc.Actor.Id, Edge{Src: c, Dest: oc})
		}
	}
	return g
}
