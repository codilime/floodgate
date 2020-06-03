local pipelines = import 'pipeline.libsonnet';

local deployment = {
  "apiVersion": "apps/v1",
  "kind": "Deployment",
  "metadata": {
    "name": "nginx-deployment2"
  },
  "spec": {
    "replicas": 3,
    "selector": {
      "matchLabels": {
        "app": "nginx"
      }
    },
    "template": {
      "metadata": {
        "labels": {
          "app": "nginx",
	  "lb": "nginx"
        }
      },
      "spec": {
        "containers": [
          {
            "name": "nginx",
            "image": "nginx:latest",
            "ports": [
              {
                "containerPort": 80,
		"": ""
	      }
            ]
          }
        ]
      }
    }
  }
};

local deploymentProd = import 'deployProd.json';

local webhookTrigger(name, source) = pipelines.triggers
                       .webhook(name)
                       .withSource(source);

local deployManifest (env, manifest = deployment)  = pipelines.stages
                              .deployManifest('Deploy nginx to ' + env)
                              .withAccount('inner-kind')
                              .withManifests(manifest)
                              .withNamespaceOverride(env);

local deployDev = deployManifest('spinnaker');
local triggerPipeline = webhookTrigger("test-pipeline", "test-pipeline");

pipelines.pipeline()
.withName('Test Pipeline 1')
.withId('test-pipeline-1')
.withApplication('testapp')
.withStages(deployDev)
.withTriggers(triggerPipeline)
