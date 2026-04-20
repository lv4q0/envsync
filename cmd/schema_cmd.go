package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envsync/internal/envfile"
)

var schemaFile string

var schemaCmd = &cobra.Command{
	Use:   "schema [env-file]",
	Short: "Validate an .env file against a JSON schema",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		envPath := args[0]

		entries, err := envfile.Parse(envPath)
		if err != nil {
			return fmt.Errorf("parse env file: %w", err)
		}

		schema, err := envfile.LoadSchema(schemaFile)
		if err != nil {
			return fmt.Errorf("load schema: %w", err)
		}

		// Convert []Entry to map[string]string for schema validation.
		envMap := make(map[string]string, len(entries))
		for _, e := range entries {
			envMap[e.Key] = e.Value
		}

		violations := envfile.ValidateAgainstSchema(envMap, schema)
		if len(violations) == 0 {
			fmt.Fprintln(cmd.OutOrStdout(), "✓ No schema violations found.")
			return nil
		}

		fmt.Fprintf(cmd.OutOrStdout(), "✗ %d schema violation(s):\n", len(violations))
		for _, v := range violations {
			fmt.Fprintf(cmd.OutOrStdout(), "  - %s\n", v.Error())
		}

		// Exit non-zero to signal failure to shell / CI.
		os.Exit(1)
		return nil
	},
}

func init() {
	schemaCmd.Flags().StringVarP(&schemaFile, "schema", "s", ".env.schema.json", "path to JSON schema file")
	rootCmd.AddCommand(schemaCmd)
}
