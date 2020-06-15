local pipelines = import 'pipeline.libsonnet';

local artifact = pipelines.artifacts.front50PipelineTemplate()
.withReference('spinnaker://jsonnetpt');

pipelines.pipeline()
.withApplication('jsonnetapp')
.withId('jsonnetpipelinept')
.withName('Example pipeline v2 created from Jsonnet file.')
.withTemplate(artifact)
.withSchema('v2')
.withInherit([])
.withNotifications([])
.withTriggers([])