package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tashima42/disciplinas/CARS001/mom/promocoes/ranking"
)

// rankingCmd represents the ranking command
var rankingCmd = &cobra.Command{
	Use:   "ranking",
	Short: "responsavel por receber promocoes, validar e publicar",
	RunE: func(cmd *cobra.Command, args []string) error {
		r, err := ranking.NewRanking(rabbitMqURL, rankingPrivateKeyPath, gatewayPublicKeyPath)
		if err != nil {
			return err
		}
		return r.Run()
	},
}

func init() {
	rootCmd.AddCommand(rankingCmd)
}
