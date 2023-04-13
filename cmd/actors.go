package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
)

// actorsCmd represents the actors command
var actorsCmd = &cobra.Command{
	Use:   "actors",
	Short: "fetch all actors from TMDB",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("actors called")
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		actors, err := TMDb.GetAllActors(ctx)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
		}

		for _, a := range actors {
			fmt.Printf("==> %d - %s\n", a.ID, a.Name)
		}
		cancel()
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
}
