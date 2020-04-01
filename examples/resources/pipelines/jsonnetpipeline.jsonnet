local pipelines = import 'pipeline.libsonnet';

pipelines.pipeline()
.withName('Example pipeline from Jsonnet')
.withId('jsonnetpipeline')
.withApplication('jsonnetapp')
