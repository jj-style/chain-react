package cmd

import (
	"context"
	"errors"
	"fmt"
	"regexp"

	"github.com/jj-style/chain-react/src/config"
	"github.com/jj-style/chain-react/src/db"
	go_tmdb "github.com/jj-style/go-tmdb"
	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
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
		c := cmd.Context().Value(config.RConfig{}).(config.RConfig)
		runGetMovies(cmd.Context(), &c)
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
func runGetMovies(ctx context.Context, c *config.RConfig) {
	actors, err := c.Repo.AllActors()
	if err != nil {
		log.Fatalln("getting all actors: ", err)
	}
	actorIds := lo.Map(actors, func(a db.Actor, _ int) int { return a.Id })

	credits := make(chan *go_tmdb.PersonMovieCredits)
	reducer := func() error {
		for cr := range credits {
			// insert movie first (if not exist)
			for _, mc := range cr.Cast {
				m := db.Movie{Id: mc.ID, Title: mc.Title}
				// create movies once
				if exist, err := c.Redis.Get(ctx, fmt.Sprintf("movie:%d", m.Id)).Bool(); err != nil || !exist {
					_, err = c.Repo.CreateMovie(m)
					if err != nil {
						if errors.Is(err, db.ErrDuplicate) {
							c.Log.Debugf("skip creating existing movie(%v)", m)
						} else {
							c.Log.Errorf("creating movie(%v): %v", m, err)
							continue
						}
					}
					c.Redis.Set(ctx, fmt.Sprintf("movie:%d", m.Id), true, 0)
				}

				// movie exists so insert credit entry for this actor
				credit := db.CreditIn{ActorId: cr.ID, MovieId: mc.ID, CreditId: mc.CreditID, Character: mc.Character}

				// check if movie already entered for actor, if so skip
				if exist, err := c.Redis.SIsMember(ctx, fmt.Sprintf("actor:%d:movies", credit.ActorId), credit.MovieId).Result(); err == nil && exist {
					continue
				}
				// check for custom reason to skip credit
				// TODO - filter inserting credits based on "character" ("self"/"himself"/"voices", ...)
				// TODO - only insert actors if poss - not directors
				if shouldSkipCredit(credit) {
					continue
				}
				if _, err = c.Repo.CreateCredit(credit); err != nil {
					if errors.Is(err, db.ErrDuplicate) {
						c.Log.Debugf("skip creating existing credit(%v)", credit)
					} else {
						c.Log.Errorf("creating credit(%v): %v", credit, err)
					}
				} else {
					c.Redis.SAdd(ctx, fmt.Sprintf("actor:%d:movies", credit.ActorId), credit.MovieId)
				}
			}
		}
		return nil
	}
	err = c.Tmdb.GetAllActorMovieCredits(ctx, credits, reducer, actorIds...)
	if err != nil {
		log.Fatalln("getting all actor movie credits: ", err)
	}
}

func shouldSkipCredit(c db.CreditIn) bool {
	// TODO - compile regexes once
	skipCharactersRe := []string{`^$`, `^Self.*$`, "^Himself.*$", "Herself.*", `.*\(voices?\).*`, `archive footage`, `\(uncredited\)`, `^Host\b`, `^Narrator\b`}
	regexes := lo.Map(skipCharactersRe, func(s string, _ int) *regexp.Regexp { return regexp.MustCompile(s) })
	for _, re := range regexes {
		if re.MatchString(c.Character) {
			return true
		}
	}
	return false
}
