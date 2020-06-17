package cmd

import (
	c "github.com/codilime/floodgate/config"
	rm "github.com/codilime/floodgate/resourcemanager"
	"github.com/spf13/cobra"
	"io"
	"log"
)

// graphOptions store render command options
type graphOptions struct {
}

// NewGraphCmd create new graph command
func NewGraphCmd(out io.Writer) *cobra.Command {
	options := graphOptions{}
	cmd := &cobra.Command{
		Use:   "graph",
		Short: "Create resources' graph",
		Long:  "Create resources' graph",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGraph(cmd, options)
		},
	}
	return cmd
}

type ResourceVertex struct {
	id    string
	Label string
}

func (rv ResourceVertex) Id() string {
	return rv.id
}

func (rv ResourceVertex) String() string {
	return rv.Label
}

func runGraph(cmd *cobra.Command, options graphOptions) error {
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
	resourceGraph.CreateResGraph()

	dot := resourceGraph.DependencyGraph.Dot(nil)

	err = resourceGraph.ExportGraphToFile(dot, "resources-graph.png")
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
