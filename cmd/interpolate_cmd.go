package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envsync/internal/envfile"
)

var interpolateCmd = &cobra.Command{
	Use:   "interpolate <file>",
	Short: "Resolve variable references in a .env file and print the result",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		filePath := args[0]

		env, err := envfile.Parse(filePath)
		if err != nil {
			return fmt.Errorf("parsing %s: %w", filePath, err)
		}

		errs := envfile.Interpolate(env)
		if len(errs) > 0 {
			fmt.Fprintln(os.Stderr, "interpolation errors:")
			for _, e := range errs {
				fmt.Fprintf(os.Stderr, "  - %s\n", e)
			}
			return fmt.Errorf("%d interpolation error(s) found", len(errs))
		}

		format, _ := cmd.Flags().GetString("format")
		out, serErr := envfile.Serialize(env, format)
		if serErr != nil {
			return fmt.Errorf("serializing output: %w", serErr)
		}

		fmt.Print(out)
		return nil
	},
}

func init() {
	interpolateCmd.Flags().StringP("format", "f", "dotenv", "Output format: dotenv, export, json")
	rootCmd.AddCommand(interpolateCmd)
}
