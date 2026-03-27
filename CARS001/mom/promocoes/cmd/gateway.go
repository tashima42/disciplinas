package cmd

import (
	"errors"

	"github.com/spf13/cobra"
	"github.com/tashima42/disciplinas/CARS001/mom/promocoes/gateway"
)

// gatewayCmd represents the gateway command
var gatewayCmd = &cobra.Command{
	Use:   "gateway",
	Short: "responsavel pela interacao com usuarios para cadastrar promocoes, listar e votar em promocoes",
}

var cadastrarSubCmd = &cobra.Command{
	Use: "cadastrar",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			return errors.New("invalid args")
		}

		titulo := args[0]
		categoria := args[1]
		return gateway.Cadastrar(rabbitMqURL, gatewayPrivateKeyPath, titulo, categoria)
	},
}

var votarSubCmd = &cobra.Command{
	Use: "votar",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("invalid args")
		}

		id := args[0]
		return gateway.Votar(rabbitMqURL, gatewayPrivateKeyPath, id)
	},
}

func init() {
	gatewayCmd.AddCommand(cadastrarSubCmd)
	gatewayCmd.AddCommand(votarSubCmd)

	rootCmd.AddCommand(gatewayCmd)
}
