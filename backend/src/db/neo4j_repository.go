package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
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
	return created.(*Actor), nil
}

func (n *Neo4jRepository) AllActors() ([]Actor, error) {
	ctx := context.TODO()
	session := n.driver.NewSession(ctx, neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeRead,
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
			node, _, err := neo4j.GetRecordValue[neo4j.Node](res.Record(), "a")
			if err != nil {
				return nil, fmt.Errorf("could not find actor node a")
			}
			actor, err := actorFromNode(node)
			if err != nil {
				return nil, err
			}
			actors = append(actors, *actor)
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
	ctx := context.TODO()
	session := n.driver.NewSession(ctx, neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeRead,
	})
	defer func() {
		session.Close(ctx)
	}()

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		res, err := tx.Run(
			ctx,
			`match (a:Actor) return a order by a.id desc limit 1`,
			nil,
		)
		if err != nil {
			return nil, err
		}

		if res.Next(ctx) {
			node, _, err := neo4j.GetRecordValue[neo4j.Node](res.Record(), "a")
			if err != nil {
				return nil, err
			}
			return actorFromNode(node)
		} else {
			return nil, ErrNotExists
		}
	})
	if err != nil {
		return &Actor{}, err
	}
	return result.(*Actor), nil
}

func (n *Neo4jRepository) RandomActor() (*Actor, error) {
	ctx := context.TODO()
	session := n.driver.NewSession(ctx, neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeRead,
	})
	defer func() {
		session.Close(ctx)
	}()

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		res, err := tx.Run(
			ctx,
			`match (a:Actor) RETURN a ORDER BY rand() LIMIT 1`,
			nil,
		)
		if err != nil {
			return nil, err
		}

		if res.Next(ctx) {
			node, _, err := neo4j.GetRecordValue[neo4j.Node](res.Record(), "a")
			if err != nil {
				return nil, err
			}
			return actorFromNode(node)
		} else {
			return nil, errors.New("no actors found")
		}
	})
	if err != nil {
		return &Actor{}, err
	}
	return result.(*Actor), nil
}

func (n *Neo4jRepository) RandomActorNotId(id int) (*Actor, error) {
	ctx := context.TODO()
	session := n.driver.NewSession(ctx, neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeRead,
	})
	defer func() {
		session.Close(ctx)
	}()

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		res, err := tx.Run(
			ctx,
			`MATCH (a:Actor) WHERE a.id <> $id RETURN a ORDER BY rand() LIMIT 1`,
			map[string]any{
				"id": id,
			},
		)
		if err != nil {
			return nil, err
		}

		if res.Next(ctx) {
			node, _, err := neo4j.GetRecordValue[neo4j.Node](res.Record(), "a")
			if err != nil {
				return nil, err
			}
			return actorFromNode(node)
		} else {
			return nil, errors.New("no actors found")
		}
	})
	if err != nil {
		return &Actor{}, err
	}
	return result.(*Actor), nil
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

func (n *Neo4jRepository) Verify(c Chain) (bool, error) {
	_, err := n.VerifyWithEdges(c)
	if err != nil {
		return false, err
	}
	return err == nil, err
}

func (n *Neo4jRepository) VerifyWithEdges(c Chain) ([]*Edge, error) {
	ctx := context.TODO()
	session := n.driver.NewSession(ctx, neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeRead,
	})
	defer func() {
		session.Close(ctx)
	}()

	edges := make([]*Edge, 0, len(c))

	for i := 0; i < len(c)-1; i++ {
		from := c[i]
		to := c[i+1]

		n.log.Infof("verifying edge between %d and %d", from, to)

		resp, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
			res, err := tx.Run(ctx, `MATCH (src:Actor {id: $src}) MATCH (dest:Actor {id: $dest}) 
				MATCH p=(src)-[r1:ACTED_IN]-(movie:Movie)-[r2:ACTED_IN]-(dest)
				return src, r1 as srcEdge, movie, dest, r2 as destEdge limit 1`, map[string]any{
				"src":  from,
				"dest": to,
			})
			if err != nil {
				return false, err
			}
			if res.Next(ctx) {
				rec := res.Record()

				srcEdge, _, _ := neo4j.GetRecordValue[neo4j.Relationship](rec, "srcEdge")
				destEdge, _, _ := neo4j.GetRecordValue[neo4j.Relationship](rec, "destEdge")
				movieNode, _, _ := neo4j.GetRecordValue[neo4j.Node](res.Record(), "movie")
				srcANode, _, _ := neo4j.GetRecordValue[neo4j.Node](res.Record(), "src")
				destANode, _, _ := neo4j.GetRecordValue[neo4j.Node](res.Record(), "src")

				srcActor, err := actorFromNode(srcANode)
				if err != nil {
					return false, err
				}
				destActor, err := actorFromNode(destANode)
				if err != nil {
					return false, err
				}
				movie, err := movieFromNode(movieNode)
				if err != nil {
					return false, err
				}
				srcCredit, err := creditFromRelation(srcEdge)
				if err != nil {
					return false, err
				}
				destCredit, err := creditFromRelation(destEdge)
				if err != nil {
					return false, err
				}
				srcCredit.Movie = *movie
				srcCredit.Actor = *srcActor
				destCredit.Movie = *movie
				destCredit.Actor = *destActor

				edges = append(edges, &Edge{
					Src:  *srcCredit,
					Dest: *destCredit,
				})
				return true, nil
			}
			return false, nil
		})

		if err != nil {
			return []*Edge{}, err
		}
		if !resp.(bool) {
			return edges, fmt.Errorf("neighbour %d not found adjacent to %d", to, from)
		}
	}

	return edges, nil
}

func actorFromNode(node dbtype.Node) (*Actor, error) {
	id, err := neo4j.GetProperty[int64](node, "id")
	if err != nil {
		return nil, err
	}
	name, err := neo4j.GetProperty[string](node, "name")
	if err != nil {
		return nil, err
	}
	return &Actor{
		Id:   int(id),
		Name: name,
	}, nil
}

func movieFromNode(node dbtype.Node) (*Movie, error) {
	id, err := neo4j.GetProperty[int64](node, "id")
	if err != nil {
		return nil, err
	}
	title, err := neo4j.GetProperty[string](node, "title")
	if err != nil {
		return nil, err
	}
	return &Movie{
		Id:    int(id),
		Title: title,
	}, nil
}

func creditFromRelation(rel dbtype.Relationship) (*Credit, error) {
	return &Credit{
		Character: rel.GetProperties()["character"].(string),
		CreditId:  rel.GetProperties()["id"].(string),
	}, nil
}
