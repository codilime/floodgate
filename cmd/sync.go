package cmd

import (
	"fmt"
	"io"

	c "github.com/codilime/floodgate/config"
	"github.com/codilime/floodgate/gateclient"
	"github.com/codilime/floodgate/parser"
	"github.com/codilime/floodgate/sync"
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
	client := gateclient.NewGateapiClient(config)
	p := parser.CreateParser(config.Libraries)
	if err := p.LoadObjectsFromDirectories(config.Resources); err != nil {
		return err
	}
	sync := &sync.Sync{}
	if err := sync.Init(client, &p.Resources); err != nil {
		return err
	}
	if options.dryRun {
		changes := sync.GetChanges()
		printChangedResources(changes)
	} else {
		if err := sync.SyncResources(); err != nil {
			return err
		}
	}
	return nil
}

func printChangedResources(changes []sync.ResourceChange) {
	fmt.Println("Following resources are changed:")
	for _, change := range changes {
		var line string
		if change.ID != "" {
			line = fmt.Sprintf("Resource: %s (%s)", change.ID, change.Name)
		} else {
			line = fmt.Sprintf("Resource: %s", change.Name)
		}
		fmt.Println(line)
		fmt.Println("Type:", change.Type)
		fmt.Println("Changes:")
		fmt.Println(change.Changes)
	}
}
