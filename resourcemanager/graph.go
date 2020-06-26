package resourcemanager

import (
	"bytes"
	gc "github.com/codilime/floodgate/gateclient"
	spr "github.com/codilime/floodgate/spinnakerresource"
	"github.com/goccy/go-graphviz"
	"github.com/hashicorp/terraform/dag"
	"github.com/hashicorp/terraform/tfdiags"
	log "github.com/sirupsen/logrus"
	"net/url"
)

// ResourceGraph stores resources and graph to be able to save it in spinnaker
type ResourceGraph struct {
	Resources spr.SpinnakerResources

	Graph dag.AcyclicGraph
}

// Create creates graph from resources
func (rg *ResourceGraph) Create() {
	start := rg.Graph.Add("Start")

	for _, application := range rg.Resources.Applications {
		appVert := rg.Graph.Add(application)
		rg.Graph.Connect(dag.BasicEdge(appVert, start))

		for _, pipeline := range rg.Resources.Pipelines {
			pipelineVert := rg.Graph.Add(pipeline)

			if pipeline.Application() == application.Name() {
				rg.Graph.Connect(dag.BasicEdge(pipelineVert, appVert))

				templateRef, _ := url.Parse(pipeline.TemplateReference())
				for _, pipelineTemplate := range rg.Resources.PipelineTemplates {
					if templateRef.Host == pipelineTemplate.ID() {
						ptVert := rg.Graph.Add(pipelineTemplate)
						rg.Graph.Connect(dag.BasicEdge(ptVert, start))
						rg.Graph.Connect(dag.BasicEdge(pipelineVert, ptVert))
					}
				}
			}
		}
	}
}

// Apply walks through graph and saves local state in spinnaker
func (rg *ResourceGraph) Apply(spinnakerAPI *gc.GateapiClient, maxConcurrentConn int) error {
	var sem = make(chan int, maxConcurrentConn)

	w := &dag.Walker{
		Callback: func(v dag.Vertex) tfdiags.Diagnostics {
			sem <- 1

			switch vertex := v.(type) {
			case *spr.Application:
				log.Infof("Saving Application: %s", vertex.Name())
				err := vertex.SaveLocalState(spinnakerAPI)
				if err != nil {
					log.Errorf("Error while saving application: %s\n%v", vertex.Name(), err)
				}

			case *spr.Pipeline:
				log.Infof("Saving Pipeline: %s", vertex.Name())
				err := vertex.SaveLocalState(spinnakerAPI)
				if err != nil {
					log.Errorf("Error while saving pipeline: %s\n%v", vertex.Name(), err)
				}
			case *spr.PipelineTemplate:
				log.Infof("Saving PipelineTemplate: %s", vertex.Name())
				err := vertex.SaveLocalState(spinnakerAPI)
				if err != nil {
					log.Errorf("Error while saving pipeline template: %s\n%v", vertex.Name(), err)
				}
			case string:
				break
			default:
				log.Infof("Unsupported type %T!\n", v)
			}

			<-sem
			return nil
		},
		Reverse: true,
	}
	w.Update(&rg.Graph)

	return w.Wait().Err()
}

// ExportGraphToFile exports graph to png
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
