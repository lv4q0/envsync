package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envsync/internal/envfile"
)

var snapshotDiffMask bool

var snapshotDiffCmd = &cobra.Command{
	Use:   "snapshot-diff <snapshot-file> <env-file>",
	Short: "Diff a saved snapshot against a current .env file",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		snapshotPath := args[0]
		envPath := args[1]

		base, err := envfile.LoadSnapshot(snapshotPath)
		if err != nil {
			return fmt.Errorf("loading snapshot: %w", err)
		}

		entries, err := envfile.Parse(envPath)
		if err != nil {
			return fmt.Errorf("parsing env file: %w", err)
		}

		current := make(map[string]string, len(entries))
		for _, e := range entries {
			current[e.Key] = e.Value
		}

		diff := envfile.DiffSnapshots(base, current)
		output := envfile.FormatSnapshotDiff(diff, snapshotDiffMask)
		fmt.Fprint(os.Stdout, output)

		if len(diff.Added)+len(diff.Removed)+len(diff.Changed) > 0 {
			os.Exit(1)
		}
		return nil
	},
}

func init() {
	snapshotDiffCmd.Flags().BoolVar(&snapshotDiffMask, "mask", true, "Mask secret values in output")
	rootCmd.AddCommand(snapshotDiffCmd)
}
