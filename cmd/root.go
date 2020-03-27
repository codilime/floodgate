package cmd

import (
	"io"

	"github.com/codilime/floodgate/cmd/cli"
	"github.com/codilime/floodgate/version"
	"github.com/spf13/cobra"
)

// RootOptions store root options
type RootOptions struct {
	configFile string
	quiet      bool
}

// Execute execute command
func Execute(out io.Writer) error {
	rootCmd := NewRootCmd(out)
	return rootCmd.Execute()
}

// TODO:
// save to spinnaker option
// 	- dry-run option
// print diff option
// print raw and hydrated pipeline
// print states from spinnaker

// NewRootCmd create new root command
func NewRootCmd(out io.Writer) *cobra.Command {
	options := RootOptions{}

	cmd := &cobra.Command{
		SilenceUsage:  true,
		SilenceErrors: true,
		Version:       version.String(),
	}
	cmd.PersistentFlags().StringVar(&options.configFile, "config", "", "path to config file (default ./config.yaml)")
	cmd.PersistentFlags().BoolVarP(&options.quiet, "quiet", "q", false, "squelch non-essential output")

	cmd.AddCommand(cli.NewSyncCmd(out))
	cmd.AddCommand(cli.NewCompareCmd(out))
	cmd.AddCommand(cli.NewHydrateCmd(out))
	cmd.AddCommand(cli.NewInspectCmd(out))
	cmd.AddCommand(cli.NewRenderCmd(out))

	return cmd
}
