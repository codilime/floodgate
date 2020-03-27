package cli

import (
	"io"

	"github.com/codilime/floodgate/cmd/gateclient"
	"github.com/codilime/floodgate/cmd/parser"
	"github.com/codilime/floodgate/cmd/sync"
	"github.com/spf13/cobra"
)

// syncOptions store sync command options
type syncOptions struct {
	dryRun bool
}

// NewSyncCmd create sync command
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
	config, err := LoadConfig(configPath)
	if err != nil {
		return err
	}
	// TODO(mlembke): check if there's no error with client
	client := gateclient.NewGateapiClient(config)
	p := parser.CreateParser(config.Libraries)
	if err := p.LoadObjectsFromDirectories(config.Resources); err != nil {
		return err
	}
	sync := &sync.Sync{}
	sync.Init(p, client)
	if err := sync.SyncResources(); err != nil {
		return err
	}
	return nil
}
