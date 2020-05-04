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

## Master version update
* application version on master should be updated after first stable release for last version
  * if we have released version `v0.1.0` then we should update master version to `v0.2.0`

## Bug fixing
* if bug exists on master branch and supported released versions:
  * first we should fix master branch
  * fix from master should be cherry-picked to all supported versions branches
  * after CI pipelines are successfully finished, new patch version should be released
* if bug exists only in certain version:
  * fix should be merge directly to version on which bug exists
  * if bug exists on more than one released and supported versions
    * bug should be fixed on latest supported version and fix should be cherry-picked to rest of supported versions
  * after CI pipelines are successfully finished, new patch version should be released
