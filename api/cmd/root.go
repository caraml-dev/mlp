package cmd

import (
	"github.com/spf13/cobra"

	"github.com/caraml-dev/mlp/api/config"
	"github.com/caraml-dev/mlp/api/log"
)

var (
	configFiles  []string
	globalConfig *config.Config
	rootCmd      = &cobra.Command{
		Use:   "mlp",
		Short: "CaraML Machine Learning Platform Console",
		Long: "CaraML Machine Learning Platform Console, which provides a web UI to interact with different CaraML " +
			"services.",
	}
)

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringSliceVarP(&configFiles, "config", "c", []string{},
		"Comma separated list of config files to load. The last config file will take precedence over the "+
			"previous ones.")
	rootCmd.AddCommand(serveCmd)
	rootCmd.AddCommand(bootstrapCmd)
}

func initConfig() {
	var err error
	globalConfig, err = config.LoadAndValidate(configFiles...)
	if err != nil {
		log.Fatalf("failed initializing config: %v", err)
	}
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		log.Fatalf("failed executing root command: %v", err)
	}
}
