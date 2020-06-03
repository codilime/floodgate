local pipelines = import 'pipeline.libsonnet';

local deployment = import 'deploy.json';

local app = "jsonnetapp";
local webhookTrigger(name, source) = pipelines.triggers
                       .webhook(name)
                       .withSource(source);

local moniker = pipelines.moniker(app);

local deployManifest (env, manifest = deployment)  = pipelines.stages
                              .deployManifest('Deploy nginx to ' + env)
                              .withAccount('inner-kind')
                              .withManifests(manifest)
                              .withNamespaceOverride(env)
                              .withMoniker(moniker);

local deployNginx = deployManifest('spinnaker');
local triggerPipeline = webhookTrigger("deploy_nginx", "deploy_nginx");

pipelines.pipeline()
.withName('Test Deployment')
.withId('test_deployment')
.withApplication(app)
.withStages(deployNginx)
.withTriggers(triggerPipeline)
