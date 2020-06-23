package cmd

import (
	c "github.com/codilime/floodgate/config"
	rm "github.com/codilime/floodgate/resourcemanager"
	"github.com/spf13/cobra"
	"io"
	"log"
)

// applyOptions store render command options
type applyOptions struct {
	graph bool
}

// NewApplyCmd create new apply command
func NewApplyCmd(out io.Writer) *cobra.Command {
	options := applyOptions{}
	cmd := &cobra.Command{
		Use:   "apply",
		Short: "Apply resources' by going through dependency graph",
		Long:  "Apply resources' by going through dependency graph",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runApply(cmd, options)
		},
	}
	cmd.Flags().BoolVarP(&options.graph, "graph", "g", false, "export dependency graph to png")
	return cmd
}

func runApply(cmd *cobra.Command, options applyOptions) error {
	flags := cmd.InheritedFlags()
	configPath, err := flags.GetString("config")
	if err != nil {
		return err
	}
	config, err := c.LoadConfig(configPath)
	if err != nil {
		return err
	}
	resourceManager := &rm.ResourceManager{}

	if err := resourceManager.Init(config); err != nil {
		return err
	}

	resources := resourceManager.GetResources()
	resourceGraph := &rm.ResourceGraph{Resources: resources}
	resourceGraph.Create()

	if options.graph {
		dot := resourceGraph.Graph.Dot(nil)

		err = resourceGraph.ExportGraphToFile(dot, "graph.png")
		if err != nil {
			log.Fatal(err)
		}
	} else {
		if err := resourceGraph.Apply(resourceManager.GetClient()); err != nil {
			log.Fatal(err)
		}
	}

	return nil
}
