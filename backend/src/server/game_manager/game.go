package gamemanager

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/jj-style/chain-react/src/db"
	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
)

const (
	GameKey string = "gameOfTheDay"
)

type Game struct {
	Start *db.Actor `json:"start"`
	End   *db.Actor `json:"end"`
}

// NewGame creates a new game/challenge to complete and saves the result in redis.
// It loops until a new unique game can be created.
// It returns on any errors, or nil if there are none.
func NewGame(logger *log.Logger, repo db.Repository, cache *redis.Client) error {
	logger.Debug("updating game")
	var start, end *db.Actor
	var err error = nil
	for {
		start, err = repo.RandomActor()
		if err != nil {
			logger.Errorf("getting start actor for new game: %v", err)
			return err
		}
		end, err = repo.RandomActorNotId(start.Id)
		if err != nil {
			logger.Errorf("getting end actor for new game: %v", err)
			return err
		}

		game := Game{
			Start: start,
			End:   end,
		}
		gameB, err := json.Marshal(game)
		if err != nil {
			log.WithField("game", game).Errorf("serializing daily game: %v", err)
			return err
		}

		logger.WithFields(log.Fields{"start": start, "end": end}).Debug("new game")
		curr, err := cache.Get(context.TODO(), GameKey).Bytes()
		if err != nil && err != redis.Nil {
			log.Errorf("getting existing game: %v", err)
			return err
		}
		if same := bytes.Equal(curr, gameB); same {
			continue
		}

		if err := cache.Set(context.TODO(), GameKey, gameB, 0).Err(); err != nil && err != redis.Nil {
			return err
		}
		break
	}
	return nil
}
