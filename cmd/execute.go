package cmd

import (
	"github.com/spf13/cobra"
	"io"
)

type executeOptions struct {
	webhookName string
}

func NewExecuteCmd(out io.Writer) *cobra.Command {
	options := executeOptions{}
	cmd := &cobra.Command{
		Use:   "execute [webhook name]",
		Short: "Trigger spinnaker webhook and track execution status",
		Long:  "Trigger spinnaker webhook and track execution status",
		RunE: func(cmd *cobra.Command, args []string) error {
			options.webhookName = args[0]

			return runExecute(cmd, options)
		},
		Args: cobra.MinimumNArgs(1),
	}

	return cmd
}

func runExecute(cmd *cobra.Command, options executeOptions) error {
	return nil
}
