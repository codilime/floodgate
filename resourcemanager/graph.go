package resourcemanager

import (
	"bytes"
	spr "github.com/codilime/floodgate/spinnakerresource"
	"github.com/goccy/go-graphviz"
	"github.com/hashicorp/terraform/dag"
)

type ResourceGraph struct {
	Resources spr.SpinnakerResources

	DependencyGraph dag.AcyclicGraph
}

type vertex struct {
	Name               string
	ConnectedTemplates []string
}

func (rg *ResourceGraph) CreateDepGraph() {
	start := rg.DependencyGraph.Add("Start")

	var ptVertex = make(map[string]dag.Vertex)
	for _, pt := range rg.Resources.PipelineTemplates {
		ptVertex[pt.ID()] = rg.DependencyGraph.Add(pt.Name())
	}

	for _, application := range rg.Resources.Applications {
		appVert := rg.DependencyGraph.Add(vertex{
			Name:               application.Name(),
			ConnectedTemplates: []string{},
		})

		rg.DependencyGraph.Connect(dag.BasicEdge(start, appVert))

		for _, pipeline := range rg.Resources.Pipelines {
			pipelineVert := rg.DependencyGraph.Add(pipeline.Name())

			if pipeline.Application() == application.Name() && pipeline.TemplateReference() == "" {
				rg.DependencyGraph.Connect(dag.BasicEdge(appVert, pipelineVert))
			}

			//if pipeline.TemplateReference() != "" {
			//templateRef, _ := url.Parse(pipeline.TemplateReference())
			//
			//applicationVert, ok := appVert.(vertex)
			//if !ok {
			//	continue
			//}

			//for _, t := range applicationVert.ConnectedTemplates {
			//	if templateRef.Host != t {
			//		rg.DependencyGraph.Connect(dag.BasicEdge(ptVertex[templateRef.Host], pipelineVert))
			//		applicationVert.ConnectedTemplates = append(applicationVert.ConnectedTemplates, templateRef.Host)
			//	}
			//}

			//templateReference, _ := url.Parse(pipeline.TemplateReference())
			//for _, pipelineTemplate := range resources.PipelineTemplates {
			//	if templateReference.Host == pipelineTemplate.ID() {
			//		pipelineTemplateG := g.Add(pipelineTemplate.Name())
			//		g.Connect(dag.BasicEdge(appG, pipelineTemplateG))
			//		g.Connect(dag.BasicEdge(pipelineTemplateG, pipelineG))
			//	}
			//}
			//}
		}
	}
}

func (rg *ResourceGraph) ExportDepGraphToFile(filename string) error {
	dot := rg.DependencyGraph.Dot(nil)
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
