package db

import (
	"github.com/jj-style/go-tmdb"
)

type Actor struct {
	Id         int     `json:"id"`
	Name       string  `json:"name"`
	Popularity float32 `json:"popularity"`
}

func ActorFromTmdbPerson(p *tmdb.Person) Actor {
	return Actor{
		Id:         p.ID,
		Name:       p.Name,
		Popularity: p.Popularity,
	}
}
