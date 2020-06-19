package resourcemanager

import (
	"bytes"
	spr "github.com/codilime/floodgate/spinnakerresource"
	"github.com/goccy/go-graphviz"
	"github.com/hashicorp/terraform/dag"
	"github.com/hashicorp/terraform/tfdiags"
	"net/url"
	"sync"
	"time"
)

type ResourceGraph struct {
	Resources spr.SpinnakerResources

	Graph dag.AcyclicGraph
}

func (rg *ResourceGraph) Create() {
	start := rg.Graph.Add("Start")

	for _, application := range rg.Resources.Applications {
		appVert := rg.Graph.Add(application.Name())
		rg.Graph.Connect(dag.BasicEdge(appVert, start))

		for _, pipeline := range rg.Resources.Pipelines {
			pipelineVert := rg.Graph.Add(pipeline.Name())

			if pipeline.Application() == application.Name() {
				rg.Graph.Connect(dag.BasicEdge(pipelineVert, appVert))

				templateRef, _ := url.Parse(pipeline.TemplateReference())
				for _, pipelineTemplate := range rg.Resources.PipelineTemplates {
					if templateRef.Host == pipelineTemplate.ID() {
						ptVert := rg.Graph.Add(pipelineTemplate.Name())
						rg.Graph.Connect(dag.BasicEdge(ptVert, start))
						rg.Graph.Connect(dag.BasicEdge(pipelineVert, ptVert))
					}
				}
			}
		}
	}
}

func (rg *ResourceGraph) Walk() error {
	var lock sync.Mutex
	w := &dag.Walker{
		Callback: func(vertex dag.Vertex) tfdiags.Diagnostics {
			lock.Lock()
			time.Sleep(1 * time.Second)
			lock.Unlock()

			return nil
		},
		Reverse: true,
	}
	w.Update(&rg.Graph)

	return w.Wait().Err()
}

func (rg *ResourceGraph) ExportGraphToFile(dot []byte, filename string) error {
	g := graphviz.New()
	graph, _ := graphviz.ParseBytes(dot)

	var buf bytes.Buffer
	if err := g.Render(graph, graphviz.PNG, &buf); err != nil {
		return err
	}

	if err := g.RenderFilename(graph, graphviz.PNG, filename); err != nil {
		return err
	}

	return nil
}
