{
  testApplication():: {
    // set default values
    dataSources: {
      disabled: [],
      enabled: [],
    },
    platformHealthOnly: false,
    platformHealthOnlyShowOverride: false,
    providerSettings: {
      aws: {
        useAmiBlockDeviceMappings: false,
      },
      gce: {
        associatePublicIpAddress: false,
      },
    },
    trafficGuards: [],
    // set overrides
    withClusters(clusters):: self + if std.type(clusters) == 'array' then { clusters: clusters } else { clusters: [clusters] },
    withCloudProviders(cloudProviders):: self + { cloudProviders: cloudProviders },
    withDescription(description):: self + { description: description },
    withEmail(email):: self + { email: email },
    withName(name):: self + { name: name },
    withUser(user):: self + { user: user },
  },
}
