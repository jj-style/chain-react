package db

import (
	"context"
	"fmt"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	log "github.com/sirupsen/logrus"
)

// match p=allShortestPaths((a:Actor{id: 3})-[:ACTED_IN*..10]-(b:Actor{id:2})) return p
// https://gist.github.com/innat/e59b71ad58b095e618a780091189cd79
// match p=(a:Actor{id: 3})-[:ACTED_IN*1..4]-(b:Actor{id:4}) return p limit 3

type Neo4jRepository struct {
	log    *log.Logger
	driver neo4j.DriverWithContext
}

func NewNeo4jRepository(driver neo4j.DriverWithContext) *Neo4jRepository {
	return &Neo4jRepository{
		driver: driver,
		log:    log.StandardLogger(),
	}
}

// General
func (n *Neo4jRepository) Migrate() error {
	return nil
}

// Actors
func (n *Neo4jRepository) CreateActor(actor Actor) (*Actor, error) {
	ctx := context.TODO()
	session := n.driver.NewSession(ctx, neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeWrite,
	})
	defer func() {
		session.Close(ctx)
	}()

	created, err := session.
		ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
			result, err := tx.Run(ctx, "MERGE (a:Actor {id: $id, name: $name}) RETURN a", map[string]any{
				"id":   actor.Id,
				"name": actor.Name,
			})
			if err != nil {
				return nil, err
			}

			if result.Next(ctx) {
				itemNode, _, err := neo4j.GetRecordValue[neo4j.Node](result.Record(), "a")
				if err != nil {
					return nil, fmt.Errorf("could not find node a")
				}
				id, err := neo4j.GetProperty[int64](itemNode, "id")
				if err != nil {
					return nil, err
				}
				name, err := neo4j.GetProperty[string](itemNode, "name")
				if err != nil {
					return nil, err
				}
				return &Actor{Id: int(id), Name: name}, nil
			}

			return nil, result.Err()
		})
	if err != nil {
		return nil, err
	}
	n.log.Printf("====++> created actor %+v\n", created)

	return created.(*Actor), nil
}

func (n *Neo4jRepository) AllActors() ([]Actor, error) {
	ctx := context.TODO()
	session := n.driver.NewSession(ctx, neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeWrite,
	})
	defer func() {
		session.Close(ctx)
	}()

	results, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		res, err := tx.Run(
			ctx,
			`MATCH (a:Actor) RETURN a`,
			nil,
		)
		if err != nil {
			return nil, err
		}

		actors := make([]Actor, 0)

		for res.Next(ctx) {
			itemNode, _, err := neo4j.GetRecordValue[neo4j.Node](res.Record(), "a")
			if err != nil {
				return nil, fmt.Errorf("could not find node a")
			}
			id, err := neo4j.GetProperty[int64](itemNode, "id")
			if err != nil {
				return nil, err
			}
			name, err := neo4j.GetProperty[string](itemNode, "name")
			if err != nil {
				return nil, err
			}
			actors = append(actors, Actor{Id: int(id), Name: name})
		}
		return actors, nil
	})
	if err != nil {
		return []Actor{}, err
	}
	return results.([]Actor), nil
}

func (n *Neo4jRepository) DeleteActor(id int64) error {
	return nil
}

func (n *Neo4jRepository) LatestActor() (*Actor, error) {
	return &Actor{}, nil
}

func (n *Neo4jRepository) RandomActor() (*Actor, error) {
	return &Actor{}, nil
}

func (n *Neo4jRepository) RandomActorNotId(int) (*Actor, error) {
	return &Actor{}, nil
}

// Movies
func (n *Neo4jRepository) CreateMovie(movie Movie) (*Movie, error) {
	ctx := context.TODO()
	session := n.driver.NewSession(ctx, neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeWrite,
	})
	defer func() {
		session.Close(ctx)
	}()

	created, err := session.
		ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
			result, err := tx.Run(ctx, "MERGE (m:Movie {id: $id, title: $title}) RETURN m", map[string]any{
				"id":    movie.Id,
				"title": movie.Title,
			})
			if err != nil {
				return nil, err
			}

			if result.Next(ctx) {
				itemNode, _, err := neo4j.GetRecordValue[neo4j.Node](result.Record(), "m")
				if err != nil {
					return nil, fmt.Errorf("could not find node n")
				}
				id, err := neo4j.GetProperty[int64](itemNode, "id")
				if err != nil {
					return nil, err
				}
				title, err := neo4j.GetProperty[string](itemNode, "title")
				if err != nil {
					return nil, err
				}
				return &Movie{Id: int(id), Title: title}, nil
			}

			return nil, result.Err()
		})
	if err != nil {
		return nil, err
	}
	n.log.Printf("====++> created movie %+v\n", created)

	return created.(*Movie), nil
}

func (n *Neo4jRepository) AllMovies() ([]Movie, error) {
	return []Movie{}, nil
}

func (n *Neo4jRepository) DeleteMovie(id int64) error {
	return nil
}

// Credits
func (n *Neo4jRepository) CreateCredit(credit CreditIn) (*CreditIn, error) {
	ctx := context.TODO()
	session := n.driver.NewSession(ctx, neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeWrite,
	})
	defer func() {
		session.Close(ctx)
	}()

	_, err := session.
		ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
			_, err := tx.Run(
				ctx,
				`MATCH (a:Actor {id: $actorId}), (m:Movie {id: $movieId})
				 WHERE a.id IS NOT NULL AND m.id IS NOT NULL
				 CREATE (a)-[:ACTED_IN {character: $character, id: $creditId, cost: 1}]->(m)`,
				map[string]any{
					"actorId":   credit.ActorId,
					"movieId":   credit.MovieId,
					"creditId":  credit.CreditId,
					"character": credit.Character,
				},
			)
			if err != nil {
				return nil, err
			}
			return nil, err
		})
	if err != nil {
		return nil, err
	}
	n.log.Printf("====++> created credit %+v\n", credit)

	return &credit, nil
}

func (n *Neo4jRepository) AllCredits() ([]Credit, error) {
	return []Credit{}, nil
}
