package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"

	c "github.com/codilime/floodgate/config"
	rm "github.com/codilime/floodgate/resourcemanager"
	"github.com/spf13/cobra"
)

// compareOptions store compare command options
type compareOptions struct {
}

// NewCompareCmd create new compare command
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
	changes, err := resourceManager.GetChanges()
	if err != nil {
		os.Exit(2)
	}
	if len(changes) == 0 {
		return nil
	}

	printCompareDiff(cmd.OutOrStdout(), changes)
	return errors.New("end diff")
}

func printCompareDiff(out io.Writer, changes []rm.ResourceChange) {
	for _, change := range changes {
		var line string
		if change.ID != "" {
			line = fmt.Sprintf("%s (%s) (%s)", change.ID, change.Name, change.Type)
		} else {
			line = fmt.Sprintf("%s (%s)", change.Name, change.Type)
		}
		fmt.Fprintln(out, line)
		fmt.Fprintln(out, change.Changes)
	}
}
