local sponnet = import 'pipeline.libsonnet';

sponnet.pipeline()
.withApplication('nginx')
.withId('nginx-demo-pipeline')
.withName('Nginx demo pipeline')