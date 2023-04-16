package cmd

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jj-style/chain-react/src/db"
	"github.com/jj-style/chain-react/src/tmdb"
	go_tmdb "github.com/jj-style/go-tmdb"
	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	update_missing bool
)

// actorsCmd represents the actors command
var actorsCmd = &cobra.Command{
	Use:   "actors",
	Short: "fetch all actors from TMDB",
	Run: func(cmd *cobra.Command, args []string) {
		t := cmd.Context().Value("tmdb").(tmdb.TMDb)
		db, err := sql.Open("sqlite3", viper.GetString("db.file"))
		if err != nil {
			log.Fatalln(err)
		}
		runGetActors(cmd.Context(), &CmdConfig{t: t, db: db, log: log.StandardLogger()})
	},
}

func init() {
	getCmd.AddCommand(actorsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// actorsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// actorsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	actorsCmd.Flags().BoolVarP(&update_missing, "update", "u", false, "Only get missing actors")
}

func runGetActors(ctx context.Context, c *CmdConfig) {
	people := make(chan *go_tmdb.Person)
	repo := db.NewSQLiteRepository(c.db)
	if err := repo.Migrate(); err != nil {
		log.Fatalln("migrating db: ", err)
	}

	reducer := func() error {
		// Reduce
		for p := range people {
			if p.Popularity < 10 {
				c.log.Warnf("skipping person(%d - %s) as popularity(%f) < 10", p.ID, p.Name, p.Popularity)
				continue
			}
			c.log.Infof("==> saving person(%d - %s)\n", p.ID, p.Name)
			actor := db.Actor{Id: p.ID, Name: p.Name}
			_, err := repo.CreateActor(actor)
			if err != nil {
				fmt.Printf("error storing %v: %v\n", actor, err)
			}
		}
		return nil
	}
	var err error
	if update_missing {
		a, err2 := repo.LatestActor()
		if err2 != nil {
			if errors.Is(err2, db.ErrNotExists) {
				err = c.t.GetAllActors(ctx, people, reducer)
			} else {
				log.Fatalln("getting latest actor: ", err2)
			}
		} else {
			err = c.t.GetActorsFrom(ctx, people, reducer, a.Id)
		}
	} else {
		err = c.t.GetAllActors(ctx, people, reducer)
	}
	if err != nil {
		log.Fatalln(err)
	}
}
