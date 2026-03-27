// Package cmd executes the cmd
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	rabbitMqURL               string
	gatewayPrivateKeyPath     string
	gatewayPublicKeyPath      string
	rankingPrivateKeyPath     string
	rankingPublicKeyPath      string
	promocaoPrivateKeyPath    string
	promocaoPublicKeyPath     string
	notificacaoPrivateKeyPath string
	notificacaoPublicKeyPath  string
)

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

	rootCmd.PersistentFlags().StringVar(&rabbitMqURL, "rabbitmq-url", "amqp://user:password@localhost:5672/", "RabbitMQ connection URL")

	rootCmd.PersistentFlags().StringVar(&gatewayPublicKeyPath, "gateway-public-key-path", "./gateway/gateway.pub", "Gateway public key path")
	rootCmd.PersistentFlags().StringVar(&gatewayPrivateKeyPath, "gateway-private-key-path", "./gateway/gateway.key", "Gateway private key path")

	rootCmd.PersistentFlags().StringVar(&rankingPrivateKeyPath, "ranking-public-key-path", "./ranking/ranking.pub", "Ranking public key path")
	rootCmd.PersistentFlags().StringVar(&rankingPublicKeyPath, "ranking-private-key-path", "./ranking/ranking.key", "Ranking private key path")

	rootCmd.PersistentFlags().StringVar(&promocaoPublicKeyPath, "promocao-public-key-path", "./promocao/promocao.pub", "Promocao public key path")
	rootCmd.PersistentFlags().StringVar(&promocaoPrivateKeyPath, "promocao-private-key-path", "./promocao/promocao.key", "Promocao private key path")

	rootCmd.PersistentFlags().StringVar(&notificacaoPublicKeyPath, "notificacao-public-key-path", "./notificacao/notificacao.pub", "Notificacao public key path")
	rootCmd.PersistentFlags().StringVar(&notificacaoPrivateKeyPath, "notificacao-private-key-path", "./notificacao/notificacao.key", "Notificacao private key path")
}
