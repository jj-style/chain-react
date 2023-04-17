package cmd

import (
	"fmt"

	"github.com/jj-style/chain-react/src/config"
	_ "github.com/mattn/go-sqlite3"
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
		fmt.Println(credits[0])
		// TODO - put into graph represented by matrix
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
