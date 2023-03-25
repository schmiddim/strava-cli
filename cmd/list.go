package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

type listOptions struct {
	format string
}

var listOpts listOptions

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List subjects of subcommand",
	Run: func(_ *cobra.Command, _ []string) {
		fmt.Println("Please specify the subject to be listed")
	},
	PersistentPreRunE: validateListCommand,
}

func init() {
	rootCmd.AddCommand(listCmd)

	flags := listCmd.PersistentFlags()
	flags.StringVar(&listOpts.format, "format", "table", "Output format (table, json)")
}

func validateListCommand(_ *cobra.Command, _ []string) error {
	return validateListOptions(&listOpts)
}

func validateListOptions(opts *listOptions) error {
	switch opts.format {
	case "table", "json":
		return nil
	default:
		return fmt.Errorf("invalid format: %s", opts.format)
	}
}
