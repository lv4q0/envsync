package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envsync/internal/envfile"
)

var compareFlagFormat string

var compareCmd = &cobra.Command{
	Use:   "compare <base> <target>",
	Short: "Compare two .env files and show a structured diff",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		basePath := args[0]
		targetPath := args[1]

		baseEntries, err := envfile.Parse(basePath)
		if err != nil {
			return fmt.Errorf("parsing base file: %w", err)
		}

		targetEntries, err := envfile.Parse(targetPath)
		if err != nil {
			return fmt.Errorf("parsing target file: %w", err)
		}

		baseMap := entriesToMap(baseEntries)
		targetMap := entriesToMap(targetEntries)

		diff := envfile.CompareEnvMaps(baseMap, targetMap)

		if len(diff.Added)+len(diff.Removed)+len(diff.Changed) == 0 {
			fmt.Fprintln(cmd.OutOrStdout(), "No differences found.")
			return nil
		}

		out := envfile.FormatEnvDiff(diff)
		fmt.Fprint(cmd.OutOrStdout(), out)
		return nil
	},
}

func entriesToMap(entries []envfile.Entry) map[string]string {
	m := make(map[string]string, len(entries))
	for _, e := range entries {
		m[e.Key] = e.Value
	}
	return m
}

func init() {
	compareCmd.Flags().StringVar(&compareFlagFormat, "format", "text", "Output format: text")
	rootCmd.AddCommand(compareCmd)
	_ = os.Stderr // suppress unused import
}
