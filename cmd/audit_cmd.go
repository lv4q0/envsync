package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/user/envsync/internal/envfile"
	"github.com/user/envsync/internal/sync"
)

var auditCmd = &cobra.Command{
	Use:   "audit <base> <target>",
	Short: "Show an audit log of changes needed to sync target from base",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		basePath := args[0]
		targetPath := args[1]

		base, err := envfile.Parse(basePath)
		if err != nil {
			return fmt.Errorf("parsing base file: %w", err)
		}

		target, err := envfile.Parse(targetPath)
		if err != nil {
			return fmt.Errorf("parsing target file: %w", err)
		}

		overwrite, _ := cmd.Flags().GetBool("overwrite")
		opts := sync.Options{Overwrite: overwrite}

		_, log, err := sync.AuditedSync(base, target, opts)
		if err != nil {
			return fmt.Errorf("auditing sync: %w", err)
		}

		fmt.Fprint(os.Stdout, log.Summary())
		return nil
	},
}

func init() {
	auditCmd.Flags().Bool("overwrite", false, "Include changed keys in the audit (as if overwrite were applied)")
	rootCmd.AddCommand(auditCmd)
}
