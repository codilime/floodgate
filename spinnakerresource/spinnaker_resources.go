package spinnakerresource

// SpinnakerResources Spinnaker resources collection
type SpinnakerResources struct {
	Applications      []*Application
	Pipelines         []*Pipeline
	PipelineTemplates []*PipelineTemplate
}
