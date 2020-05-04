# Release process for Floodgate

## Preparations before creating first release candidate for new version
* make sure that all required tasks are completed
* create branch with proper name:
  * prefix should be `release-` phrase
  * middle part should be composed from `v` letters and major and minor version number: `v0.1`
  * suffix should always be `.x`
  * example proper release branch name: `release-v0.1.x`
* make sure that all CI tasks are executed on with success on new branch

## Creating new release candidate
* create new release in GitHub
  * release name should include full version number and following release candidate number
  * we are starting counting releases from `1`
  * example of proper release candidate name: `v0.1.0-rc1`
  * release should be marked as `pre-release` in GitHub
* compiled binaries should be build using CI system and added after CI system completes all task started after creating new release

## Creating proper release
* create new release in GitHub
  * release name should start with letter `v` and proper full version number, for example: `v0.1.0`
* binaries should be added after all tasks started in CI system are completed
