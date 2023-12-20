package gamemanager

import (
	"context"
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
	manager := &CronGameManager{
		log:   logger,
		repo:  repo,
		redis: cache,
		cron:  c,
	}

	c.AddFunc(conf.GameSchedule, func() {
		manager.CreateGame(context.TODO())
	})

	c.Start()
	logger.WithField("schedule", conf.GameSchedule).Info("started game manager")

	return manager, func() {
		c.Stop()
	}, nil
}

func (c *CronGameManager) GetGame(ctx context.Context) (*Game, error) {
	return GetGame(ctx, c.redis)
}

func (c *CronGameManager) CreateGame(ctx context.Context) (*Game, error) {
	return NewGame(ctx, c.log, c.repo, c.redis)
}
