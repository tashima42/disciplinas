package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/tashima42/disciplinas/CARS001/mom/hot-deals/crypto"
)

var (
	keyPairPath string
	keyPairName string
)

// cryptoCmd represents the crypto command
var cryptoCmd = &cobra.Command{
	Use:   "crypto",
	Short: "generates rsa keypairs",
	RunE: func(cmd *cobra.Command, args []string) error {
		private, public, err := crypto.GenerateKeyPair()
		if err != nil {
			return err
		}

		privateKeyPath := filepath.Join(keyPairPath, keyPairName+".key")
		publicKeyPath := filepath.Join(keyPairPath, keyPairName+".pub")

		if err := os.WriteFile(privateKeyPath, private, 0o400); err != nil {
			return err
		}

		if err := os.WriteFile(publicKeyPath, public, 0o400); err != nil {
			return err
		}

		fmt.Printf("private key: %s, public key: %s\n", privateKeyPath, publicKeyPath)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(cryptoCmd)
	cryptoCmd.Flags().StringVarP(&keyPairPath, "path", "p", ".", "directory path to write the keys to")
	cryptoCmd.Flags().StringVarP(&keyPairName, "name", "n", "id_rsa", "name for the keys: private (id_rsa.key), public (id_rsa.pub)")
}
