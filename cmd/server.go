package cmd

import (
	"fmt"

	"github.com/jj-style/chain-react/src/config"
	"github.com/jj-style/chain-react/src/db"
	"github.com/jj-style/chain-react/src/graph"
	_ "github.com/mattn/go-sqlite3"
	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := cmd.Context().Value(config.RConfig{}).(config.RConfig)

		credits, err := cfg.Repo.AllCredits()
		if err != nil {
			log.Fatalln("getting all credits: ", err)
		}

		type edge struct {
			src  db.Credit
			dest db.Credit
		}

		g := graph.NewGraph[db.Actor, edge]()

		// add all actors as vertices in graph
		actors, err := cfg.Repo.AllActors()
		if err != nil {
			log.Fatalln("getting all actors: ", err)
		}
		for _, a := range actors {
			g.AddVertex(a.Id, a)
		}

		creditsByMovieId := lo.GroupBy(credits, func(c db.Credit) int { return c.Movie.Id })
		for _, c := range credits {
			// build edge between actor in credit, and every other actor in that movie
			a := c.Actor
			for _, oc := range creditsByMovieId[c.Movie.Id] {
				if oc.Actor == a {
					// avoid cycles
					continue
				}
				g.AddEdge(a.Id, oc.Actor.Id, edge{src: c, dest: oc})
			}
		}

		paths := make(chan []graph.Element[db.Actor, edge])

		go g.Bfs(31, 48, paths)
		for p := range paths {
			// TODO - ignore paths that hop via the same movie
			fmt.Println("***********")
			for idx, ve := range p {
				if idx == 0 {
					continue
				}
				fmt.Printf("%s:%s <=%s=> %s:%s\n", ve.Edge.Weight.src.Name, ve.Edge.Weight.src.Character, ve.Edge.Weight.src.Title, ve.Edge.Weight.dest.Name, ve.Edge.Weight.dest.Character)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serverCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serverCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	serverCmd.Flags().Int("port", 8080, "Port to run chain-react server on")
	viper.BindPFlag("port", serverCmd.Flags().Lookup("port"))

}
