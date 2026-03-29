package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tashima42/disciplinas/CARS001/mom/promocoes/notificacao"
)

// notificacaoCmd represents the notificacao command
var notificacaoCmd = &cobra.Command{
	Use:   "notificacao",
	Short: "responsavel por receber promocoes e publicar em filas por categoria",
	RunE: func(cmd *cobra.Command, args []string) error {
		n, err := notificacao.NewNotificacao(rabbitMqURL, promocaoPublicKeyPath, rankingPublicKeyPath)
		if err != nil {
			return err
		}
		return n.Run()
	},
}

func init() {
	rootCmd.AddCommand(notificacaoCmd)
}
