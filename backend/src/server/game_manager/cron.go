package gamemanager

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/jj-style/chain-react/src/config"
	"github.com/jj-style/chain-react/src/db"
	"github.com/redis/go-redis/v9"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
)

// CronGameManager is a GameManager that manages games according to a given cron schedule.
type CronGameManager struct {
	repo  db.Repository
	redis *redis.Client
	log   *log.Logger
	cron  *cron.Cron
}

// Creates a new CronGameManager
func NewCronGameManager(logger *log.Logger, conf *config.Server, repo db.Repository, cache *redis.Client) (GameManager, func(), error) {
	c := cron.New(cron.WithLocation(time.UTC))

	c.AddFunc(conf.GameSchedule, func() {
		NewGame(logger, repo, cache)
	})

	c.Start()
	logger.WithField("schedule", conf.GameSchedule).Info("started game manager")

	return &CronGameManager{
			log:   logger,
			repo:  repo,
			redis: cache,
			cron:  c,
		}, func() {
			c.Stop()
		}, nil
}

func (c *CronGameManager) GetGame(ctx context.Context) (*Game, error) {
	got, err := c.redis.Get(ctx, GameKey).Bytes()
	if err == redis.Nil {
		return nil, errors.New("current game not set")
	}
	if err != nil {
		c.log.Errorf("getting current game: %v", err)
		return nil, err
	}

	var g Game
	if err := json.Unmarshal(got, &g); err != nil {
		c.log.Errorf("parsing current game: %v", err)
		return nil, err
	}

	return &g, nil
}
