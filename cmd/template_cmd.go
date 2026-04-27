package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"envsync/internal/envfile"
)

func init() {
	var templatePath string
	var applyDefaults bool

	templateCmd := &cobra.Command{
		Use:   "template",
		Short: "Validate an env file against a .env.template",
		Long:  "Check that all required keys from a .env.template are present in the target env file, optionally filling defaults.",
		RunE: func(cmd *cobra.Command, args []string) error {
			envPath, _ := cmd.Flags().GetString("file")
			if envPath == "" {
				envPath = ".env"
			}
			if templatePath == "" {
				templatePath = ".env.template"
			}

			tmpl, err := envfile.LoadTemplate(templatePath)
			if err != nil {
				return fmt.Errorf("load template: %w", err)
			}

			entries, err := envfile.Parse(envPath)
			if err != nil {
				return fmt.Errorf("parse env: %w", err)
			}

			env := make(map[string]string, len(entries))
			for _, e := range entries {
				env[e.Key] = e.Value
			}

			missing := envfile.ApplyTemplate(env, tmpl)

			if applyDefaults && len(missing) == 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "All required keys present. Defaults applied where needed.\n")
				return nil
			}

			if len(missing) > 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "Missing required keys:\n  %s\n", strings.Join(missing, "\n  "))
				return fmt.Errorf("%d required key(s) missing", len(missing))
			}

			fmt.Fprintf(cmd.OutOrStdout(), "OK: all required template keys are present in %s\n", envPath)
			return nil
		},
	}

	templateCmd.Flags().StringVarP(&templatePath, "template", "t", ".env.template", "Path to the .env.template file")
	templateCmd.Flags().BoolVar(&applyDefaults, "apply-defaults", false, "Report after applying template defaults")
	rootCmd.AddCommand(templateCmd)
}
