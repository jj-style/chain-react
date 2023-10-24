//go:build wireinject
// +build wireinject

package cmd

import (
	"github.com/google/wire"
	"github.com/jj-style/chain-react/src/config"
	"github.com/jj-style/chain-react/src/server"
	"github.com/sirupsen/logrus"
)

// wireApp init application.
func wireApp(tConf *config.TmdbConfig, dbConf *config.DbConfig, sConf *config.MeilisearchConfig, rConf *config.RedisConfig, logger *logrus.Logger) (*server.Server, func(), error) {
	panic(wire.Build(server.ProviderSet))
}
