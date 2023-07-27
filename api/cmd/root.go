package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/caraml-dev/mlp/api/log"
)

var (
	rootCmd = &cobra.Command{
		Use:   "mlp",
		Short: "CaraML Machine Learning Platform Console",
		Long: "CaraML Machine Learning Platform Console, which provides a web UI to interact with different CaraML " +
			"services. If no subcommand are provided, serve command will be run as default.",
	}
)

func init() {
	rootCmd.AddCommand(serveCmd)
	rootCmd.AddCommand(bootstrapCmd)
}

func Execute() {
	cmd, _, err := rootCmd.Find(os.Args[1:])
	// use serve as default cmd if no cmd is given
	if err == nil && cmd.Use == rootCmd.Use && cmd.Flags().Parse(os.Args[1:]) != pflag.ErrHelp {
		args := append([]string{serveCmd.Use}, os.Args[1:]...)
		rootCmd.SetArgs(args)
	}

	err = rootCmd.Execute()
	if err != nil {
		log.Fatalf("failed executing root command: %v", err)
	}
}
