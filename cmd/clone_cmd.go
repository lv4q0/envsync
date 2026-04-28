package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"envsync/internal/envfile"
)

func init() {
	var overwrite bool
	var stripSecrets bool
	var profile string

	cloneCmd := &cobra.Command{
		Use:   "clone <source> <destination>",
		Short: "Clone an env file to a new location",
		Long: `Clone copies a .env file to a destination path.

Optionally strip secret values (keys matching patterns like *_SECRET, *_PASSWORD)
so the clone is safe to share or commit as a template.

Use --profile to write the clone as a named profile variant (e.g. .env.staging).`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			src := args[0]
			dest := args[1]

			opts := envfile.CloneOptions{
				Overwrite:    overwrite,
				StripSecrets: stripSecrets,
				Profile:      profile,
			}

			result, err := envfile.Clone(src, dest, opts)
			if err != nil {
				return fmt.Errorf("clone failed: %w", err)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Cloned %q → %q\n", src, result.Destination)
			fmt.Fprintf(cmd.OutOrStdout(), "  Keys written : %d\n", result.KeysWritten)
			if stripSecrets {
				fmt.Fprintf(cmd.OutOrStdout(), "  Secrets stripped: %d\n", result.SecretsStripped)
			}
			return nil
		},
	}

	cloneCmd.Flags().BoolVar(&overwrite, "overwrite", false, "Overwrite destination if it already exists")
	cloneCmd.Flags().BoolVar(&stripSecrets, "strip-secrets", false, "Replace secret values with empty strings in the clone")
	cloneCmd.Flags().StringVar(&profile, "profile", "", "Write clone as a named profile variant (appends .<profile> to destination)")

	rootCmd.AddCommand(cloneCmd)
}
