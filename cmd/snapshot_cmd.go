package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"

	"envsync/internal/envfile"
)

func init() {
	var outputPath string

	snapshotCmd := &cobra.Command{
		Use:   "snapshot <env-file>",
		Short: "Capture a snapshot of an env file for later diffing",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			srcPath := args[0]

			snap, err := envfile.TakeSnapshot(srcPath)
			if err != nil {
				return fmt.Errorf("snapshot: %w", err)
			}

			if outputPath == "" {
				base := filepath.Base(srcPath)
				ts := time.Now().UTC().Format("20060102T150405")
				outputPath = fmt.Sprintf("%s.%s.snapshot.json", base, ts)
			}

			if err := envfile.SaveSnapshot(snap, outputPath); err != nil {
				return fmt.Errorf("snapshot: %w", err)
			}

			fmt.Fprintf(os.Stdout, "Snapshot saved to %s (%d keys)\n", outputPath, len(snap.Entries))
			return nil
		},
	}

	snapshotCmd.Flags().StringVarP(&outputPath, "output", "o", "", "destination path for the snapshot JSON file")

	compareSnapshotCmd := &cobra.Command{
		Use:   "compare-snapshot <snapshot-file> <env-file>",
		Short: "Diff a saved snapshot against a current env file",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			snap, err := envfile.LoadSnapshot(args[0])
			if err != nil {
				return fmt.Errorf("compare-snapshot: %w", err)
			}

			current, err := envfile.Parse(args[1])
			if err != nil {
				return fmt.Errorf("compare-snapshot: %w", err)
			}

			fmt.Fprintf(os.Stdout, "Comparing snapshot (%s) → %s\n",
				snap.Timestamp.Format(time.RFC3339), args[1])

			added, removed, changed := 0, 0, 0
			for k, v := range current {
				if old, ok := snap.Entries[k]; !ok {
					fmt.Fprintf(os.Stdout, "  + %s\n", k)
					added++
				} else if old != v {
					fmt.Fprintf(os.Stdout, "  ~ %s\n", k)
					changed++
				}
			}
			for k := range snap.Entries {
				if _, ok := current[k]; !ok {
					fmt.Fprintf(os.Stdout, "  - %s\n", k)
					removed++
				}
			}
			fmt.Fprintf(os.Stdout, "Summary: +%d -%d ~%d\n", added, removed, changed)
			return nil
		},
	}

	rootCmd.AddCommand(snapshotCmd)
	rootCmd.AddCommand(compareSnapshotCmd)
}
