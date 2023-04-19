/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"

	"github.com/jj-style/chain-react/src/config"
	"github.com/spf13/cobra"
)

// indexCmd represents the index command
var indexCmd = &cobra.Command{
	Use:   "index",
	Short: "Index data into the search database",
	Run: func(cmd *cobra.Command, args []string) {
		c := cmd.Context().Value(config.RConfig{}).(config.RConfig)
		runIndex(cmd.Context(), &c)
	},
}

func init() {
	manageCmd.AddCommand(indexCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// indexCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// indexCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func runIndex(ctx context.Context, c *config.RConfig) {
	actors, err := c.Repo.AllActors()
	if err != nil {
		c.Log.Fatalln("getting actors: ", err)
	}
	if len(actors) == 0 {
		c.Log.Fatalln("No actors to index")
	}
	if err = c.Search.AddDocuments(actors, "actors"); err != nil {
		c.Log.Fatalln(err)
	}
}
