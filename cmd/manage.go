package cmd

import (
	"context"

	"github.com/jj-style/chain-react/src/tmdb"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// manageCmd represents the manage command
var manageCmd = &cobra.Command{
	Use:   "manage",
	Short: "manage the TMDB database",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		api_key := viper.GetString("tmdb.api_key")
		t := tmdb.NewClient(api_key)
		cmd.SetContext(context.WithValue(cmd.Context(), "tmdb", t))
	},
}

func init() {
	rootCmd.AddCommand(manageCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// syncCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// syncCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
