package cmd

import (
	"database/sql"

	"github.com/jj-style/chain-react/src/tmdb"
	log "github.com/sirupsen/logrus"
)

type CmdConfig struct {
	t   tmdb.TMDb
	db  *sql.DB
	log *log.Logger
}
