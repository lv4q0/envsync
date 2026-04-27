package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envsync/internal/envfile"
)

func init() {
	var rollbackDir string
	var label string
	var restore string
	var list bool

	rollbackCmd := &cobra.Command{
		Use:   "rollback",
		Short: "Save or restore rollback points for an .env file",
		RunE: func(cmd *cobra.Command, args []string) error {
			base, _ := cmd.Flags().GetString("file")
			if base == "" {
				base = ".env"
			}

			if list {
				points, err := envfile.ListRollbackPoints(rollbackDir)
				if err != nil {
					return err
				}
				if len(points) == 0 {
					fmt.Fprintln(cmd.OutOrStdout(), "No rollback points found.")
					return nil
				}
				for i, p := range points {
					fmt.Fprintf(cmd.OutOrStdout(), "[%d] %s  label=%q  path=%s\n",
						i, p.Timestamp.Format("2006-01-02T15:04:05Z"), p.Label, p.Path)
				}
				return nil
			}

			if restore != "" {
				points, err := envfile.ListRollbackPoints(rollbackDir)
				if err != nil {
					return err
				}
				for _, p := range points {
					if p.Path == restore {
						if err := envfile.RestoreRollbackPoint(p, base); err != nil {
							return err
						}
						fmt.Fprintf(cmd.OutOrStdout(), "Restored %s from %s\n", base, restore)
						return nil
					}
				}
				return fmt.Errorf("rollback point not found: %s", restore)
			}

			// default: save a new rollback point
			entry, err := envfile.SaveRollbackPoint(base, rollbackDir, label)
			if err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Rollback point saved: %s\n", entry.Path)
			return nil
		},
	}

	rollbackCmd.Flags().StringVar(&rollbackDir, "dir", ".envsync/rollbacks", "Directory to store rollback points")
	rollbackCmd.Flags().StringVar(&label, "label", "manual", "Label for this rollback point")
	rollbackCmd.Flags().StringVar(&restore, "restore", "", "Path of rollback point to restore")
	rollbackCmd.Flags().BoolVar(&list, "list", false, "List available rollback points")
	rollbackCmd.Flags().String("file", ".env", "Target .env file")

	_ = os.MkdirAll(".envsync/rollbacks", 0700)
	rootCmd.AddCommand(rollbackCmd)
}
