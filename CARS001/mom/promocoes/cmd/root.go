// Package cmd executes the cmd
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rabbitMQURL string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "promocoes",
	Short: "Find, rank and send promocoes",
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
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&rabbitMQURL, "rabbitmq-url", "", "RabbitMQ connection URL")
}
