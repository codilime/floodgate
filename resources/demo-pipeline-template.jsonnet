local pipelineTemplates = import 'v2PipelineTemplate.libsonnet';

local metadata = pipelineTemplates.metadata()
.withName('Demo pipeline template')
.withDescription('Demo Pipeline Template.')
.withOwner('demo-app-user@example.com')
.withScopes(['global']);

pipelineTemplates.pipelineTemplate()
.withId('demo-pipeline-template')
.withMetadata(metadata)
