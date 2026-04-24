package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envsync/internal/envfile"
)

func init() {
	var keyHex string
	var decrypt bool
	var outputFormat string

	encryptCmd := &cobra.Command{
		Use:   "encrypt <file>",
		Short: "Encrypt or decrypt secret values in an env file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			filePath := args[0]

			if len(keyHex) == 0 {
				return fmt.Errorf("--key is required (16, 24, or 32 ASCII characters)")
			}
			key := []byte(keyHex)
			if len(key) != 16 && len(key) != 24 && len(key) != 32 {
				return fmt.Errorf("--key must be exactly 16, 24, or 32 characters, got %d", len(key))
			}

			env, err := envfile.Parse(filePath)
			if err != nil {
				return fmt.Errorf("failed to parse env file: %w", err)
			}

			var processed map[string]string
			if decrypt {
				processed, err = envfile.DecryptSecrets(env, key)
				if err != nil {
					return fmt.Errorf("decryption failed: %w", err)
				}
			} else {
				processed, err = envfile.EncryptSecrets(env, key)
				if err != nil {
					return fmt.Errorf("encryption failed: %w", err)
				}
			}

			out, err := envfile.Serialize(processed, outputFormat)
			if err != nil {
				return fmt.Errorf("serialization failed: %w", err)
			}

			_, err = fmt.Fprint(os.Stdout, out)
			return err
		},
	}

	encryptCmd.Flags().StringVar(&keyHex, "key", "", "Encryption key (16, 24, or 32 ASCII characters)")
	encryptCmd.Flags().BoolVar(&decrypt, "decrypt", false, "Decrypt secret values instead of encrypting")
	encryptCmd.Flags().StringVar(&outputFormat, "format", "dotenv", "Output format: dotenv, export, json")

	rootCmd.AddCommand(encryptCmd)
}
