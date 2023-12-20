package gamemanager

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"

	"github.com/jj-style/chain-react/src/db"
	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
)

const (
	GameKey string = "gameOfTheDay"
)

type GameManager interface {
	// GetGame gets the current game
	GetGame(context.Context) (*Game, error)
	// Generate a new game
	CreateGame(context.Context) (*Game, error)
}

type Game struct {
	Start *db.Actor `json:"start"`
	End   *db.Actor `json:"end"`
}

// NewGame creates a new game/challenge to complete and saves the result in redis.
// It loops until a new unique game can be created.
// It returns the created game if successful, and any errors, or nil if there are none.
func NewGame(ctx context.Context, logger *log.Logger, repo db.Repository, cache *redis.Client) (*Game, error) {
	logger.Debug("updating game")
	var start, end *db.Actor
	var game Game
	var err error = nil
	for {
		start, err = repo.RandomActor()
		if err != nil {
			logger.Errorf("getting start actor for new game: %v", err)
			return nil, err
		}
		end, err = repo.RandomActorNotId(start.Id)
		if err != nil {
			logger.Errorf("getting end actor for new game: %v", err)
			return nil, err
		}

		game = Game{
			Start: start,
			End:   end,
		}
		gameB, err := json.Marshal(game)
		if err != nil {
			log.WithField("game", game).Errorf("serializing daily game: %v", err)
			return nil, err
		}

		logger.WithFields(log.Fields{"start": start, "end": end}).Debug("new game")
		curr, err := cache.Get(ctx, GameKey).Bytes()
		if err != nil && err != redis.Nil {
			log.Errorf("getting existing game: %v", err)
			return nil, err
		}
		if same := bytes.Equal(curr, gameB); same {
			continue
		}

		if err := cache.Set(context.TODO(), GameKey, gameB, 0).Err(); err != nil && err != redis.Nil {
			return nil, err
		}
		break
	}
	return &game, nil
}

// GetGame gets the current game, or errors retrieved while doing so.
func GetGame(ctx context.Context, cache *redis.Client) (*Game, error) {
	got, err := cache.Get(ctx, GameKey).Bytes()
	if err == redis.Nil {
		return nil, errors.New("current game not set")
	}
	if err != nil {
		return nil, err
	}

	var g Game
	if err := json.Unmarshal(got, &g); err != nil {
		return nil, err
	}

	return &g, nil
}
