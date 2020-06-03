local sponnet = import 'application.libsonnet';

sponnet.application()
.withName('testapp')
.withDescription('Example application created from Jsonnet file.')
.withUser('admin')
.withEmail('rafal.paradowski@codilime.com')
.withCloudProviders('kubernetes')
