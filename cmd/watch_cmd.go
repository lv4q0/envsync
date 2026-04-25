package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"

	"envsync/internal/diff"
	"envsync/internal/envfile"
)

var watchInterval int

var watchCmd = &cobra.Command{
	Use:   "watch <base> <target>",
	Short: "Watch a .env file and print a diff whenever it changes",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		basePath := args[0]
		targetPath := args[1]

		interval := time.Duration(watchInterval) * time.Millisecond

		done := make(chan struct{})
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		go func() {
			<-sigs
			close(done)
		}()

		fmt.Fprintf(cmd.OutOrStdout(), "Watching %s (interval: %v) — press Ctrl+C to stop\n", targetPath, interval)

		ch := envfile.Watch(targetPath, interval, done)
		for ev := range ch {
			if ev.Err != nil {
				return fmt.Errorf("watch error: %w", ev.Err)
			}
			if !ev.Changed {
				continue
			}
			fmt.Fprintf(cmd.OutOrStdout(), "\n[%s] Change detected in %s\n", time.Now().Format(time.RFC3339), ev.Path)

			base, err := envfile.Parse(basePath)
			if err != nil {
				fmt.Fprintf(cmd.ErrOrStderr(), "parse base: %v\n", err)
				continue
			}
			target, err := envfile.Parse(targetPath)
			if err != nil {
				fmt.Fprintf(cmd.ErrOrStderr(), "parse target: %v\n", err)
				continue
			}
			results := diff.Compare(base, target)
			diff.Report(cmd.OutOrStdout(), results)
		}
		fmt.Fprintln(cmd.OutOrStdout(), "\nWatch stopped.")
		return nil
	},
}

func init() {
	watchCmd.Flags().IntVar(&watchInterval, "interval", 500, "Poll interval in milliseconds")
	rootCmd.AddCommand(watchCmd)
}
