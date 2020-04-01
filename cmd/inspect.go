package cmd

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

// inspectOptions store inspect command options
type inspectOptions struct {
}

// NewInspectCmd create new inspect command
func NewInspectCmd(out io.Writer) *cobra.Command {
	options := inspectOptions{}
	cmd := &cobra.Command{
		Use:   "inspect",
		Short: "Inspect resources' status on Spinnaker",
		Long:  "Inspect resources' status on Spinnaker",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runInspect(cmd, options)
		},
	}
	return cmd
}

func runInspect(cmd *cobra.Command, options inspectOptions) error {
	return fmt.Errorf("not implemented")
}
