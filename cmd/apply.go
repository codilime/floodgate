package cmd

import (
	c "github.com/codilime/floodgate/config"
	rm "github.com/codilime/floodgate/resourcemanager"
	"github.com/spf13/cobra"
	"io"
	"io/ioutil"
	"log"
	"os"
)

// applyOptions store apply command options
type applyOptions struct {
	graph                    bool
	outputPath               string
	maxConcurrentConnections int
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
	cmd.Flags().BoolVarP(&options.graph, "dot", "d", false, "export dependency graph to dot format")
	cmd.Flags().StringVarP(&options.outputPath, "output-path", "o", "graph.dot", "dot output path")
	cmd.Flags().IntVarP(&options.maxConcurrentConnections, "maxConcurrentConnections", "c", 5, "max concurrent connections to spinnaker")
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
	config.Merge(cfg)

	resourceManager := &rm.ResourceManager{}

	if err := resourceManager.Init(config); err != nil {
		os.Exit(2)
	}

	resources := resourceManager.GetResources()
	resourceGraph := &rm.ResourceGraph{Resources: resources}
	resourceGraph.Create()

	if options.graph {
		dot := resourceGraph.Graph.Dot(nil)

		err = ioutil.WriteFile(options.outputPath, dot, 0644)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		if err := resourceGraph.Apply(resourceManager.GetClient(), options.maxConcurrentConnections); err != nil {
			log.Fatal(err)
		}
	}

	return nil
}
