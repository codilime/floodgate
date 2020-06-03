package cmd

import (
	"fmt"
	"io"

	"github.com/codilime/floodgate/version"
	"github.com/spf13/cobra"
)

// renderOptions store render command options
type versionOptions struct {
}

// NewVersionCmd create new render command
func NewVersionCmd(out io.Writer) *cobra.Command {
	options := renderOptions{}
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print info about version",
		Long:  "Prints information about version and build details",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runVersion(cmd, options)
		},
	}
	return cmd
}

func runVersion(cmd *cobra.Command, options renderOptions) error {
	fmt.Fprintf(cmd.OutOrStdout(), version.BuildInfo())
	return nil
}
