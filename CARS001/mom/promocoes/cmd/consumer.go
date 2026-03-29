package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tashima42/disciplinas/CARS001/mom/promocoes/consumer"
)

var categories []string

// consumerCmd represents the consumer command
var consumerCmd = &cobra.Command{
	Use:   "consumer",
	Short: "recebe notificacacoes",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := consumer.NewConsumer(rabbitMqURL, categories)
		if err != nil {
			return err
		}
		return c.Run()
	},
}

func init() {
	rootCmd.AddCommand(consumerCmd)

	consumerCmd.Flags().StringSliceVarP(&categories, "categories", "c", nil, "categories to listen notifications")
}
