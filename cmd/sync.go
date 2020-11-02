package cmd

import (
	"fmt"
	"io"
	"os"

	c "github.com/codilime/floodgate/config"
	rm "github.com/codilime/floodgate/resourcemanager"
	"github.com/spf13/cobra"
)

// syncOptions store synchronize command options
type syncOptions struct {
	dryRun bool
}

// NewSyncCmd create synchronize command
func NewSyncCmd(out io.Writer) *cobra.Command {
	options := syncOptions{}
	cmd := &cobra.Command{
		Use:     "synchronize",
		Aliases: []string{"sync"},
		Short:   "Synchronize resources to Spinnaker",
		Long:    "Synchronize resources to Spinnaker",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSync(cmd, options)
		},
	}
	cmd.Flags().BoolVarP(&options.dryRun, "dry-run", "d", false, "process resources and preview the result but don't sync with Spinnaker")
	return cmd
}

func runSync(cmd *cobra.Command, options syncOptions) error {
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
		return err
	}
	if options.dryRun {
		changes, err := resourceManager.GetChanges()
		if err != nil {
			os.Exit(1)
		}
		printChangedResources(cmd.OutOrStdout(), changes)
	} else {
		if err := resourceManager.SyncResources(); err != nil {
			return err
		}
	}
	return nil
}

func printChangedResources(out io.Writer, changes []rm.ResourceChange) {
	fmt.Fprintln(out, "Following resources are changed:")
	for _, change := range changes {
		var line string
		if change.ID != "" {
			line = fmt.Sprintf("Resource: %s (%s)", change.ID, change.Name)
		} else {
			line = fmt.Sprintf("Resource: %s", change.Name)
		}
		fmt.Fprintln(out, line)
		fmt.Fprintln(out, "Type:", change.Type)
		fmt.Fprintln(out, "Changes:")
		fmt.Fprintln(out, change.Changes)
	}
}
