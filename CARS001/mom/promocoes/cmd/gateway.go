package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tashima42/disciplinas/CARS001/mom/promocoes/gateway"
)

// gatewayCmd represents the gateway command
var gatewayCmd = &cobra.Command{
	Use:   "gateway",
	Short: "responsavel pela interacao com usuarios para cadastrar promocoes, listar e votar em promocoes",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("gateway called")
	},
}

var cadastrarSubCmd = &cobra.Command{
	Use: "cadastrar-promocao",
	RunE: func(cmd *cobra.Command, args []string) error {
		return gateway.Cadastrar()
	},
}

func init() {
	gatewayCmd.AddCommand(cadastrarSubCmd)

	rootCmd.AddCommand(gatewayCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// gatewayCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// gatewayCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
