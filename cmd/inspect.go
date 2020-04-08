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
	flags := cmd.InheritedFlags()
	configPath, err := flags.GetString("config")
	if err != nil {
		return err
	}
	resourceHandler, err := getResourceHandler(configPath)
	if err != nil {
		return err
	}
	fmt.Println("Current Spinnaker resource status:")
	fmt.Println("\nApplications:")
	fmt.Println(resourceHandler.GetAllApplicationsRemoteState())
	fmt.Println("\nPipelines:")
	fmt.Println(resourceHandler.GetAllPipelinesRemoteState())
	fmt.Println("\nPipeline templates:")
	fmt.Println(resourceHandler.GetAllPipelineTemplatesRemoteState())
	return nil
}
