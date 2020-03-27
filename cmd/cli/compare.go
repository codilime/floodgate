package cli

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

// compareOptions store sync command options
type compareOptions struct {
}

// NewCompareCmd ?
func NewCompareCmd(out io.Writer) *cobra.Command {
	options := compareOptions{}
	cmd := &cobra.Command{
		Use:   "compare",
		Short: "Compare local resources' definitions with Spinnaker and show discrepancies",
		Long:  "Compare local resources' definitions with Spinnaker and show discrepancies",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCompare(cmd, options)
		},
	}
	return cmd
}

func runCompare(cmd *cobra.Command, options compareOptions) error {
	return fmt.Errorf("not implemented")
}
