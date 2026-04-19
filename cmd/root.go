package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envsync/internal/diff"
	"envsync/internal/envfile"
	"envsync/internal/sync"
)

var (
	overwrite bool
	dryRun    bool
)

var rootCmd = &cobra.Command{
	Use:   "envsync",
	Short: "Diff and sync .env files across environments",
}

var diffCmd = &cobra.Command{
	Use:   "diff <base> <target>",
	Short: "Show differences between two .env files",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		base, err := envfile.Parse(args[0])
		if err != nil {
			return fmt.Errorf("parsing base: %w", err)
		}
		target, err := envfile.Parse(args[1])
		if err != nil {
			return fmt.Errorf("parsing target: %w", err)
		}
		results := diff.Compare(base, target)
		fmt.Print(diff.Report(results))
		return nil
	},
}

var syncCmd = &cobra.Command{
	Use:   "sync <source> <destination>",
	Short: "Sync keys from source into destination .env file",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		source, err := envfile.Parse(args[0])
		if err != nil {
			return fmt.Errorf("parsing source: %w", err)
		}
		dest, err := envfile.Parse(args[1])
		if err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("parsing destination: %w", err)
		}
		added, skipped, err := sync.Sync(source, dest, args[1], overwrite, dryRun)
		if err != nil {
			return fmt.Errorf("syncing: %w", err)
		}
		fmt.Printf("Added: %d  Skipped: %d", added, skipped)
		if dryRun {
			fmt.Print("  (dry run)")
		}
		fmt.Println()
		return nil
	},
}

func init() {
	syncCmd.Flags().BoolVar(&overwrite, "overwrite", false, "Overwrite changed keys in destination")
	syncCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Preview changes without writing")
	rootCmd.AddCommand(diffCmd, syncCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
