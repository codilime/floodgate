package cmd

import (
	"fmt"
	"io"
	"os"

	c "github.com/codilime/floodgate/config"
	rm "github.com/codilime/floodgate/resourcemanager"
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
	config, err := c.LoadConfig(configPath)
	if err != nil {
		return err
	}
	config.Merge(cfg)
	resourceManager := &rm.ResourceManager{}
	if err := resourceManager.Init(config); err != nil {
		os.Exit(2)
	}
	fmt.Fprintln(cmd.OutOrStdout(), "Current Spinnaker resource status:")
	fmt.Fprintln(cmd.OutOrStdout(), "\nApplications:")
	fmt.Fprintln(cmd.OutOrStdout(), resourceManager.GetAllApplicationsRemoteState())
	fmt.Fprintln(cmd.OutOrStdout(), "\nPipelines:")
	fmt.Fprintln(cmd.OutOrStdout(), resourceManager.GetAllPipelinesRemoteState())
	fmt.Fprintln(cmd.OutOrStdout(), "\nPipeline templates:")
	fmt.Fprintln(cmd.OutOrStdout(), resourceManager.GetAllPipelineTemplatesRemoteState())
	return nil
}
