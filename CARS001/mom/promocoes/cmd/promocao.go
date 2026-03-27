package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tashima42/disciplinas/CARS001/mom/promocoes/promocao"
)

// promocaoCmd represents the promocao command
var promocaoCmd = &cobra.Command{
	Use:   "promocao",
	Short: "responsavel por receber promocoes, validar e publicar",
	RunE: func(cmd *cobra.Command, args []string) error {
		return promocao.RecebePromocao(rabbitMqURL, gatewayPublicKeyPath, promocaoPrivateKeyPath)
	},
}

func init() {
	rootCmd.AddCommand(promocaoCmd)
}
