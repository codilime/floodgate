package cmd

import (
	"bytes"
	c "github.com/codilime/floodgate/config"
	rm "github.com/codilime/floodgate/resourcemanager"
	"github.com/goccy/go-graphviz"
	"github.com/spf13/cobra"
	"io"
	"log"
	"net/url"
)

// renderOptions store render command options
type graphOptions struct {
}

// NewVersionCmd create new render command
func NewGraphCmd(out io.Writer) *cobra.Command {
	options := renderOptions{}
	cmd := &cobra.Command{
		Use:   "graph",
		Short: "Print dependency graph",
		Long:  "Print dependency graph",
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

func runGraph(cmd *cobra.Command, options renderOptions) error {
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

	//d := dag.NewDAG()
	g := graphviz.New()
	graph, err := g.Graph()
	if err != nil {
		log.Fatal(err)
	}
	//dot, sfdp, neato
	graph.SetLayout("sfdp")

	defer func() {
		if err := graph.Close(); err != nil {
			log.Fatal(err)
		}
		g.Close()
	}()

	for _, application := range resources.Applications {
		appNodeG, _ := graph.CreateNode(application.Name())

		for _, pipeline := range resources.Pipelines {
			pipelineNodeN, err := graph.CreateNode(pipeline.Name())
			if err != nil {
				log.Fatal(err)
			}

			if pipeline.Application() == application.Name() {
				_, err := graph.CreateEdge("dependency", pipelineNodeN, appNodeG)
				if err != nil {
					log.Fatal(err)
				}
			}

			templateReference, _ := url.Parse(pipeline.TemplateReference())
			for _, pipelineTemplate := range resources.PipelineTemplates {
				if templateReference.Host == pipelineTemplate.ID() {
					pipelineTemplateNode, err := graph.CreateNode(pipelineTemplate.Name())
					if err != nil {
						log.Fatal(err)
					}

					_, err = graph.CreateEdge("dependency", pipelineTemplateNode, pipelineNodeN)
					if err != nil {
						log.Fatal(err)
					}
				}
			}
		}
	}
	var buf bytes.Buffer
	if err := g.Render(graph, "dot", &buf); err != nil {
		log.Fatal(err)
	}

	if err := g.RenderFilename(graph, graphviz.PNG, "graph.png"); err != nil {
		log.Fatal(err)
	}

	return nil
}
