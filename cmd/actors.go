package cmd

import (
	"context"
	"errors"
	"fmt"

	"github.com/jj-style/chain-react/src/config"
	"github.com/jj-style/chain-react/src/db"
	go_tmdb "github.com/jj-style/go-tmdb"
	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	update_missing bool
)

// actorsCmd represents the actors command
var actorsCmd = &cobra.Command{
	Use:   "actors",
	Short: "fetch all actors from TMDB",
	Run: func(cmd *cobra.Command, args []string) {
		c := cmd.Context().Value(config.RConfig{}).(config.RConfig)
		runGetActors(cmd.Context(), &c)
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

func runGetActors(ctx context.Context, c *config.RConfig) {
	people := make(chan *go_tmdb.Person)

	reducer := func() error {
		// Reduce
		for p := range people {
			if p.Popularity < 10 {
				c.Log.Warnf("skipping person(%d - %s) as popularity(%f) < 10", p.ID, p.Name, p.Popularity)
				continue
			}
			c.Log.Infof("==> saving person(%d - %s)\n", p.ID, p.Name)
			actor := db.Actor{Id: p.ID, Name: p.Name}
			_, err := c.Repo.CreateActor(actor)
			if err != nil {
				fmt.Printf("error storing %v: %v\n", actor, err)
			}
		}
		return nil
	}
	var err error
	if update_missing {
		a, err2 := c.Repo.LatestActor()
		if err2 != nil {
			if errors.Is(err2, db.ErrNotExists) {
				err = c.Tmdb.GetAllActors(ctx, people, reducer)
			} else {
				log.Fatalln("getting latest actor: ", err2)
			}
		} else {
			err = c.Tmdb.GetActorsFrom(ctx, people, reducer, a.Id)
		}
	} else {
		err = c.Tmdb.GetAllActors(ctx, people, reducer)
	}
	if err != nil {
		log.Fatalln(err)
	}
}
