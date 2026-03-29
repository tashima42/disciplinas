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
		v, err := promocao.NewVerificador(rabbitMqURL, gatewayPublicKeyPath, promocaoPrivateKeyPath)
		if err != nil {
			return err
		}
		return v.Run()
	},
}

func init() {
	rootCmd.AddCommand(promocaoCmd)
}
