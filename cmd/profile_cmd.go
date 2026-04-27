package cmd

import (
	"fmt"
	"os"
	"sort"

	"github.com/spf13/cobra"

	"envsync/internal/envfile"
)

func init() {
	profileCmd := &cobra.Command{
		Use:   "profile",
		Short: "Manage and compare environment profiles",
	}

	listCmd := &cobra.Command{
		Use:   "list [dir]",
		Short: "List available environment profiles in a directory",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			dir := "."
			if len(args) == 1 {
				dir = args[0]
			}
			profiles, err := envfile.ListProfiles(dir)
			if err != nil {
				return err
			}
			if len(profiles) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "No profiles found.")
				return nil
			}
			for _, p := range profiles {
				fmt.Fprintln(cmd.OutOrStdout(), p)
			}
			return nil
		},
	}

	diffCmd := &cobra.Command{
		Use:   "diff <base-profile> <target-profile>",
		Short: "Diff two environment profiles",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			dir, _ := cmd.Flags().GetString("dir")

			base, err := envfile.LoadProfile(dir, args[0])
			if err != nil {
				return fmt.Errorf("loading base profile: %w", err)
			}
			target, err := envfile.LoadProfile(dir, args[1])
			if err != nil {
				return fmt.Errorf("loading target profile: %w", err)
			}

			diffs := envfile.DiffProfiles(base, target)
			if len(diffs) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "Profiles are identical.")
				return nil
			}

			keys := make([]string, 0, len(diffs))
			for k := range diffs {
				keys = append(keys, k)
			}
			sort.Strings(keys)

			for _, k := range keys {
				pair := diffs[k]
				baseVal, targetVal := pair[0], pair[1]
				if envfile.IsSecret(k) {
					if baseVal != "" {
						baseVal = "***"
					}
					if targetVal != "" {
						targetVal = "***"
					}
				}
				switch {
				case baseVal == "":
					fmt.Fprintf(cmd.OutOrStdout(), "+ %s=%s\n", k, targetVal)
				case targetVal == "":
					fmt.Fprintf(cmd.OutOrStdout(), "- %s=%s\n", k, baseVal)
				default:
					fmt.Fprintf(cmd.OutOrStdout(), "~ %s: %s -> %s\n", k, baseVal, targetVal)
				}
			}
			return nil
		},
	}
	diffCmd.Flags().String("dir", ".", "Directory containing profile files")

	profileCmd.AddCommand(listCmd, diffCmd)
	rootCmd.AddCommand(profileCmd)

	_ = os.Getenv // suppress unused import
}
