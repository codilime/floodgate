{
  test():: {
    // set default values
    variable: false,
    // set overrides
    withName(name):: self + { name: name },
  },
}
