package gamemanager

import (
	"time"

	"github.com/jj-style/chain-react/src/config"
	"github.com/jj-style/chain-react/src/db"
	"github.com/redis/go-redis/v9"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
)

type CronGameManager struct {
	repo  db.Repository
	redis *redis.Client
	log   *log.Logger
	cron  *cron.Cron
}

func NewCronGameManager(logger *log.Logger, conf *config.Server, repo db.Repository, redis *redis.Client) (*CronGameManager, func(), error) {
	c := cron.New(cron.WithLocation(time.UTC))
	c.AddFunc(conf.GameSchedule, func() {
		logger.Debug("updating game")
		var start, end *db.Actor
		var err error = nil
		for {
			start, err = repo.RandomActor()
			if err != nil {
				logger.Errorf("getting start actor for new game: %v", err)
				return
			}
			end, err = repo.RandomActorNotId(start.Id)
			if err != nil {
				logger.Errorf("getting end actor for new game: %v", err)
				return
			}

			// TODO: check redis, if okay - set redis and break
			// if not in redis {...
			logger.WithFields(log.Fields{"start": start, "end": end}).Debug("new game")
			break
		}

	})
	logger.WithField("schedule", conf.GameSchedule).Info("starting game manager")
	c.Start()
	return &CronGameManager{
			log:   logger,
			repo:  repo,
			redis: redis,
			cron:  c,
		}, func() {
			c.Stop()
		}, nil
}
