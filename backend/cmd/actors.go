package cmd

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strings"

	"github.com/jj-style/chain-react/src/config"
	"github.com/jj-style/chain-react/src/db"
	go_tmdb "github.com/jj-style/go-tmdb"
	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	update_missing bool
	actor_file     string
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
	actorsCmd.Flags().StringVarP(&actor_file, "file", "f", "", "File to load actors from")
}

func runGetActors(ctx context.Context, c *config.RConfig) {
	people := make(chan *go_tmdb.Person)

	// listen for interrup signals to tear down the errgroup
	ctx, cancel := context.WithCancel(ctx)
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	go func() {
		for sig := range signals {
			// sig is a ^C, handle it
			c.Log.WithField("signal", sig).Info("received signal, disposing resources")
			cancel()
		}
	}()

	reducer := func() error {
		// Reduce
		newPeople := make([]db.Actor, 0)
		for p := range people {
			if p.Popularity < 10 {
				c.Log.Warnf("skipping person(%d - %s) as popularity(%f) < 10", p.ID, p.Name, p.Popularity)
				continue
			}
			c.Log.Infof("==> saving person(%d - %s)\n", p.ID, p.Name)
			actor := db.ActorFromTmdbPerson(p)
			_, err := c.Repo.CreateActor(actor)
			if err != nil {
				fmt.Printf("error storing %v: %v\n", actor, err)
			}
			newPeople = append(newPeople, actor)
		}

		// Index every fetched actor to the search DB
		if err := c.Search.AddDocuments(newPeople, "actors"); err != nil {
			fmt.Printf("error indexing new people: %v\n", err)
		}

		return nil
	}
	var err error
	if update_missing {
		err = updateMissingActors(ctx, c, people, reducer)
	} else if actor_file != "" {
		err = getActorsFromFile(ctx, c, people, reducer, actor_file)
	} else {
		err = c.Tmdb.GetAllActors(ctx, people, reducer)
	}
	if err != nil {
		log.Fatalln(err)
	}
}

func updateMissingActors(ctx context.Context, c *config.RConfig, people chan *go_tmdb.Person, reducer func() error) error {
	var err error = nil
	if latest, errLatest := c.Repo.LatestActor(); errLatest != nil {
		if errors.Is(errLatest, db.ErrNotExists) {
			err = c.Tmdb.GetAllActors(ctx, people, reducer)
		} else {
			log.Fatalln("getting latest actor: ", errLatest)
		}
	} else {
		err = c.Tmdb.GetActorsFrom(ctx, people, reducer, latest.Id)
	}
	return err
}

func getActorsFromFile(ctx context.Context, c *config.RConfig, people chan *go_tmdb.Person, reducer func() error, file string) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	actorNames := make([]string, 0)

	scan := bufio.NewScanner(f)
	scan.Split(bufio.ScanLines)
	for scan.Scan() {
		name := strings.TrimSpace(scan.Text())
		actorNames = append(actorNames, name)
	}
	err = c.Tmdb.GetActorsByName(ctx, people, reducer, actorNames...)
	if err != nil {
		return err
	}
	return f.Close()
}
