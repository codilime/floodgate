local pipelineTemplates = import 'v2PipelineTemplate.libsonnet';

local metadata = pipelineTemplates.metadata()
.withName('Example pipeline template from Jsonnet')
.withDescription('Example pipeline template created from Jsonnet file.')
.withOwner('floodgate@example.com')
.withScopes(['global']);

pipelineTemplates.pipelineTemplate()
.withId('jsonnetpt')
.withMetadata(metadata)
