package cmd

import (
	"fmt"
	"os"

	"github.com/jj-style/chain-react/src/tmdb"
	go_tmdb "github.com/ryanbradynd05/go-tmdb"
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
		t := cmd.Context().Value("tmdb").(tmdb.TMDb)

		people := make(chan *go_tmdb.Person)
		reducer := func() error {
			// Reduce
			for p := range people {
				fmt.Printf("==> %d - %s\n", p.ID, p.Name)
			}
			return nil
		}
		var err error
		if update_missing {
			err = t.GetMissingActors(cmd.Context(), people, reducer)
		} else {
			err = t.GetAllActors(cmd.Context(), people, reducer)
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			return
		}
	},
}

func init() {
	syncCmd.AddCommand(actorsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// actorsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// actorsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	actorsCmd.Flags().BoolVarP(&update_missing, "update", "u", false, "Only get missing actors")
}
