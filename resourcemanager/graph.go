package resourcemanager

import (
	"bytes"
	"fmt"
	spr "github.com/codilime/floodgate/spinnakerresource"
	"github.com/goccy/go-graphviz"
	"github.com/hashicorp/terraform/dag"
	"net/url"
)

type ResourceGraph struct {
	Resources spr.SpinnakerResources

	ResourcesGraph  dag.AcyclicGraph
	DependencyGraph dag.AcyclicGraph
}

func (rg *ResourceGraph) CreateDepGraph() {
	worker1 := rg.DependencyGraph.Add("Worker 1")

	//var ptVertex = make(map[string]dag.Vertex)
	//for _, pt := range rg.Resources.PipelineTemplates {
	//	ptVertex[pt.ID()] = rg.DependencyGraph.Add(pt.Name())
	//}

	for _, application := range rg.Resources.Applications {
		appVert := rg.DependencyGraph.Add(application.Name())
		rg.DependencyGraph.Connect(dag.BasicEdge(worker1, appVert))

		rg.handleApplicationDep(application, appVert)
	}

	rg.DependencyGraph.DepthFirstWalk([]dag.Vertex{worker1}, func(vertex dag.Vertex, i int) error {
		fmt.Println(vertex)
		return nil
	})
}

func (rg *ResourceGraph) handleApplicationDep(application *spr.Application, appVert dag.Vertex) {
	for _, pipeline := range rg.Resources.Pipelines {
		pipelineVert := rg.DependencyGraph.Add(pipeline.Name())

		if pipeline.Application() == application.Name() {
			if pipeline.TemplateReference() == "" {
				rg.DependencyGraph.Connect(dag.BasicEdge(appVert, pipelineVert))
			} else {
				templateReference, _ := url.Parse(pipeline.TemplateReference())
				for _, pipelineTemplate := range rg.Resources.PipelineTemplates {
					if templateReference.Host == pipelineTemplate.ID() {
						ptVert := rg.DependencyGraph.Add(pipelineTemplate.Name())
						rg.DependencyGraph.Connect(dag.BasicEdge(appVert, ptVert))
						rg.DependencyGraph.Connect(dag.BasicEdge(ptVert, pipelineVert))
					}
				}
			}
		}
	}
}

func (rg *ResourceGraph) CreateResGraph() {
	start := rg.ResourcesGraph.Add("Start")

	for _, application := range rg.Resources.Applications {
		appVert := rg.ResourcesGraph.Add(application.Name())
		rg.ResourcesGraph.Connect(dag.BasicEdge(start, appVert))

		for _, pipeline := range rg.Resources.Pipelines {
			pipelineVert := rg.ResourcesGraph.Add(pipeline.Name())

			if pipeline.Application() == application.Name() {
				rg.ResourcesGraph.Connect(dag.BasicEdge(pipelineVert, appVert))
			}

			templateRef, _ := url.Parse(pipeline.TemplateReference())
			for _, pipelineTemplate := range rg.Resources.PipelineTemplates {
				if templateRef.Host == pipelineTemplate.ID() {
					ptVert := rg.ResourcesGraph.Add(pipelineTemplate.Name())
					rg.ResourcesGraph.Connect(dag.BasicEdge(ptVert, pipelineVert))
				}
			}
		}
	}
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
