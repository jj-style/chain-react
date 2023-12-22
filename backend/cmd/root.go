package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/jj-style/chain-react/src/config"
	"github.com/jj-style/chain-react/src/server"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var cleanup func()

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "chain-react",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		cfg := &config.Config{}
		err := viper.Unmarshal(cfg)
		if err != nil {
			log.Fatalln("unmarshalling config: ", err)
		}
		var s *server.Server

		logger := log.New()
		logger.SetFormatter(&log.TextFormatter{
			DisableTimestamp: false,
			FullTimestamp:    true,
		})
		logger.SetReportCaller(true)
		if lvl, err := log.ParseLevel(viper.GetString("LOG.LEVEL")); err == nil {
			log.WithField("level", lvl.String()).Info("setting log level")
			logger.SetLevel(lvl)
		}

		s, cleanup, err = wireApp(&cfg.Server, &cfg.Tmdb, &cfg.Db, &cfg.Meilisearch, &cfg.Redis, logger)
		if err != nil {
			panic(err)
		}
		cmd.SetContext(context.WithValue(cmd.Context(), server.Server{}, s))
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		cleanup()
	},
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ~/chain-react.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".chain-react" (without extension).
		viper.AddConfigPath(home)
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")
		viper.SetConfigName(".chain-react")
	}

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
