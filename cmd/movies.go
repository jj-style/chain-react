package cmd

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jj-style/chain-react/src/db"
	"github.com/jj-style/chain-react/src/tmdb"
	go_tmdb "github.com/jj-style/go-tmdb"
	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// moviesCmd represents the movies command
var moviesCmd = &cobra.Command{
	Use:   "movies",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		t := cmd.Context().Value("tmdb").(tmdb.TMDb)
		db, err := sql.Open("sqlite3", viper.GetString("db.file"))
		if err != nil {
			log.Fatalln(err)
		}
		_, err = db.Exec("PRAGMA foreign_keys = ON;")
		if err != nil {
			log.Fatalln(err)
		}
		runGetMovies(cmd.Context(), &CmdConfig{t: t, db: db, log: log.StandardLogger()})
	},
}

func init() {
	getCmd.AddCommand(moviesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// moviesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// moviesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
func runGetMovies(ctx context.Context, c *CmdConfig) {
	repo := db.NewSQLiteRepository(c.db)
	if err := repo.Migrate(); err != nil {
		log.Fatalln("migrating db: ", err)
	}

	actors, err := repo.AllActors()
	if err != nil {
		log.Fatalln("getting all actors: ", err)
	}
	actorIds := lo.Map(actors, func(a db.Actor, _ int) int { return a.Id })

	credits := make(chan *go_tmdb.PersonMovieCredits)
	reducer := func() error {
		for cr := range credits {
			// fmt.Println(cr)

			// insert movie first (if not exist)
			for _, mc := range cr.Cast {
				m := db.Movie{Id: mc.ID, Title: mc.Title}
				_, err = repo.CreateMovie(m)
				if err != nil {
					if errors.Is(err, db.ErrDuplicate) {
						c.log.Debugf("skip creating existing movie(%v)", m)
					} else {
						c.log.Errorf("creating movie(%v): %v", m, err)
						continue
					}
				}
				// movie exists so insert credit entry for this actor
				credit := db.CreditIn{ActorId: cr.ID, MovieId: mc.ID, CreditId: mc.CreditID, Character: mc.Character}
				// TODO - filter inserting credits based on "character" ("self"/"himself"/"voices", ...)
				// TODO - only insert actors if poss - not directors
				if _, err = repo.CreateCredit(credit); err != nil {
					if errors.Is(err, db.ErrDuplicate) {
						c.log.Debugf("skip creating existing credit(%v)", credit)
					} else {
						c.log.Errorf("creating credit(%v): %v", credit, err)
					}
				}
			}
		}
		return nil
	}
	err = c.t.GetAllActorMovieCredits(ctx, credits, reducer, actorIds...)
	if err != nil {
		log.Fatalln("getting all actor movie credits: ", err)
	}
}
