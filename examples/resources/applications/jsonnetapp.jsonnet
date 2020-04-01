local sponnet = import 'application.libsonnet';

sponnet.application()
.withName('jsonnetapp')
.withDescription('Example application created from Jsonnet file.')
.withUser('floodgate@example.com')
.withEmail('example@example.com')
.withCloudProviders('kubernetes')
