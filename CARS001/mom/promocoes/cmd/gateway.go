package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tashima42/disciplinas/CARS001/mom/promocoes/gateway"
)

// gatewayCmd represents the gateway command
var gatewayCmd = &cobra.Command{
	Use:   "gateway",
	Short: "responsavel pela interacao com usuarios para cadastrar promocoes, listar e votar em promocoes",
	RunE: func(cmd *cobra.Command, args []string) error {
		return gateway.Run(rabbitMqURL, gatewayPrivateKeyPath, promocaoPublicKeyPath)
	},
}

func init() {
	rootCmd.AddCommand(gatewayCmd)
}
