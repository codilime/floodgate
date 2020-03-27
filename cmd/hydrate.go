package cmd

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

// hydrateOptions store hydrate command options
type hydrateOptions struct {
}

// NewHydrateCmd create new hydrate command
func NewHydrateCmd(out io.Writer) *cobra.Command {
	options := hydrateOptions{}
	cmd := &cobra.Command{
		Use:   "hydrate",
		Short: "Hydrate pipeline templates with configurations and preview the result",
		Long:  "Hydrate pipeline templates with configurations, without creating actual pipelines in Spinnaker",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runHydrate(cmd, options)
		},
	}
	return cmd
}

func runHydrate(cmd *cobra.Command, options hydrateOptions) error {
	return fmt.Errorf("not implemented")
}
