// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package cmd

import (
	"github.com/jj-style/chain-react/src/config"
	"github.com/jj-style/chain-react/src/db"
	"github.com/jj-style/chain-react/src/redis"
	"github.com/jj-style/chain-react/src/search"
	"github.com/jj-style/chain-react/src/server"
	"github.com/jj-style/chain-react/src/tmdb"
	"github.com/sirupsen/logrus"
)

import (
	_ "github.com/mattn/go-sqlite3"
)

// Injectors from wire.go:

// wireApp init application.
func wireApp(tConf *config.TmdbConfig, dbConf *config.DbConfig, sConf *config.MeilisearchConfig, rConf *config.RedisConfig, logger *logrus.Logger) (*server.Server, func(), error) {
	repository, err := db.NewRepository(dbConf)
	if err != nil {
		return nil, nil, err
	}
	tmDb := tmdb.NewTMDb(tConf, logger)
	searchRepository := search.NewRepository(sConf, logger)
	client := cache.NewRedis(rConf)
	serverServer := server.NewServer(logger, repository, tmDb, searchRepository, client)
	return serverServer, func() {
	}, nil
}
