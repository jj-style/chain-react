package cmd

import (
	"github.com/jj-style/chain-react/src/tmdb"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	TMDb tmdb.TMDb
)

// syncCmd represents the sync command
var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "re-sync database with TMDB database",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		api_key := viper.GetString("tmdb.api_key")
		TMDb = tmdb.NewClient(api_key)
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// syncCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// syncCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
