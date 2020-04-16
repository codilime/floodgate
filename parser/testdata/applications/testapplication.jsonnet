local testlib = import 'testlib.libsonnet';

testlib.testApplication()
.withName('testappjsonnet')
.withUser('example@example.com')
.withDescription('Test application from Jsonnet file.')
.withEmail('example@example.com')
.withCloudProviders('kubernetes')