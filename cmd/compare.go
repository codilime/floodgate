package cmd

import (
	"errors"
	"fmt"
	"io"

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
	resourceManager := &rm.ResourceManager{}
	if err := resourceManager.Init(configPath); err != nil {
		return err
	}
	changes := resourceManager.GetChanges()
	if len(changes) == 0 {
		return nil
	}
	printCompareDiff(changes)
	return errors.New("end diff")
}

func printCompareDiff(changes []rm.ResourceChange) {
	for _, change := range changes {
		var line string
		if change.ID != "" {
			line = fmt.Sprintf("%s (%s) (%s)", change.ID, change.Name, change.Type)
		} else {
			line = fmt.Sprintf("%s (%s)", change.Name, change.Type)
		}
		fmt.Println(line)
		fmt.Println(change.Changes)
	}
}
