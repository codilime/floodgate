package cli

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

// renderOptions store sync command options
type renderOptions struct {
}

// NewRenderCmd ?
func NewRenderCmd(out io.Writer) *cobra.Command {
	options := renderOptions{}
	cmd := &cobra.Command{
		Use:   "render",
		Short: "Render Jsonnet files",
		Long:  "Render Jsonnet files",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRender(cmd, options)
		},
	}
	return cmd
}

func runRender(cmd *cobra.Command, options renderOptions) error {
	return fmt.Errorf("not implemented")
}
