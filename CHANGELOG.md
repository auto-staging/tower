# [1.3.0](https://github.com/auto-staging/tower/compare/1.2.1...1.3.0) (2019-04-10)


### Features

* added go structs for component version informations ([3ab2595](https://github.com/auto-staging/tower/commit/3ab2595))
* added versions endpoint to get version information about all Auto Staging components ([2cee155](https://github.com/auto-staging/tower/commit/2cee155))

## [1.2.1](https://github.com/auto-staging/tower/compare/1.2.0...1.2.1) (2019-03-28)


### Bug Fixes

* compile binary for linux ([68296dc](https://github.com/auto-staging/tower/commit/68296dc))

# [1.2.0](https://github.com/auto-staging/tower/compare/1.1.0...1.2.0) (2019-02-22)


### Features

* made region dynamic based on the lambda region ([52b3d59](https://github.com/auto-staging/tower/commit/52b3d59))

# [1.1.0](https://github.com/auto-staging/tower/compare/1.0.1...1.1.0) (2019-02-18)


### Features

* added webhok secret token to tower configuration endpoint - fixes [#5](https://github.com/auto-staging/tower/issues/5) ([1b939f3](https://github.com/auto-staging/tower/commit/1b939f3))

## [1.0.1](https://github.com/auto-staging/tower/compare/1.0.0...1.0.1) (2019-01-26)


### Bug Fixes

* fixed [#2](https://github.com/auto-staging/tower/issues/2) missing check for webhook parameter in github webhook endpoints ([aafd90c](https://github.com/auto-staging/tower/commit/aafd90c))

# 1.0.0 (2019-01-10)


### Bug Fixes

* added check for current status in start and stop trigger - fixes [#9](https://github.com/auto-staging/tower/issues/9) ([2653c98](https://github.com/auto-staging/tower/commit/2653c98))
* added check for current status to webhook env delete - fixes [#9](https://github.com/auto-staging/tower/issues/9) ([bb76dc1](https://github.com/auto-staging/tower/commit/bb76dc1))
* added missing CodeBuidRoleARN to dynamodb update expression ([4d3069a](https://github.com/auto-staging/tower/commit/4d3069a))
* added missing status to check for update - fixes [#9](https://github.com/auto-staging/tower/issues/9) ([bfabcc3](https://github.com/auto-staging/tower/commit/bfabcc3))
* check current status before executing environment update - fixes [#9](https://github.com/auto-staging/tower/issues/9) ([d59c101](https://github.com/auto-staging/tower/commit/d59c101))
* check environment status before executing delete - fixes [#9](https://github.com/auto-staging/tower/issues/9) ([0d63257](https://github.com/auto-staging/tower/commit/0d63257))
* fixed condition expression for environment add ([80c2315](https://github.com/auto-staging/tower/commit/80c2315))
* fixed typo in json annotation ([3310349](https://github.com/auto-staging/tower/commit/3310349))
* only allow environment updates in "running" and "updating failed" status ([b92caa1](https://github.com/auto-staging/tower/commit/b92caa1))
* return success for github delete webhook if the environment is already deleted, - fixes [#8](https://github.com/auto-staging/tower/issues/8) ([c9e783b](https://github.com/auto-staging/tower/commit/c9e783b))


### Features

* added builder invokation for environment update and replaced RepoURL by InfrastructureRepoUrl ([7794da2](https://github.com/auto-staging/tower/commit/7794da2))
* added check if environments exist before deleting parent repository ([99f8bf0](https://github.com/auto-staging/tower/commit/99f8bf0))
* added check to prevent environment destroy while in init state ([d3e7f70](https://github.com/auto-staging/tower/commit/d3e7f70))
* added codebuildarn variable to environment model update and create functions ([73a8878](https://github.com/auto-staging/tower/commit/73a8878))
* added codeBuildRoleARN to repo global config and repository functions ([a4df1f1](https://github.com/auto-staging/tower/commit/a4df1f1))
* added controller to update the tower configuration ([e7290a5](https://github.com/auto-staging/tower/commit/e7290a5))
* added data structure for environments with add and get all  endpoints ([a3a7b5a](https://github.com/auto-staging/tower/commit/a3a7b5a))
* added delete single environment endpoint ([4f69d04](https://github.com/auto-staging/tower/commit/4f69d04))
* added endpoint to delete single repository ([3b828a5](https://github.com/auto-staging/tower/commit/3b828a5))
* added endpoint to get single environment ([39d75e6](https://github.com/auto-staging/tower/commit/39d75e6))
* added endpoint to get single environment status information ([ef7fb3d](https://github.com/auto-staging/tower/commit/ef7fb3d))
* added endpoint to get single repository and moved database logic into model ([64115cc](https://github.com/auto-staging/tower/commit/64115cc))
* added endpoint to get status information for all environments ([89fc848](https://github.com/auto-staging/tower/commit/89fc848))
* added endpoint to trigger scheduler directly ([9937326](https://github.com/auto-staging/tower/commit/9937326))
* added endpoint to update existing environment ([74a8cc1](https://github.com/auto-staging/tower/commit/74a8cc1))
* added endpoint to update repository ([b9d086f](https://github.com/auto-staging/tower/commit/b9d086f))
* added endpoints to get and set global environment configuraton ([4b70ba7](https://github.com/auto-staging/tower/commit/4b70ba7))
* added EnvironmentVariables object to Repository data structure ([d642180](https://github.com/auto-staging/tower/commit/d642180))
* added get and post endpoints for repositories ([e73f682](https://github.com/auto-staging/tower/commit/e73f682))
* added GitHub ping and create event webhook endpoint ([ad5bdce](https://github.com/auto-staging/tower/commit/ad5bdce))
* added hmac check for github webhooks ([6862a73](https://github.com/auto-staging/tower/commit/6862a73))
* added IAM Role validator for repository add and update ([c07f662](https://github.com/auto-staging/tower/commit/c07f662))
* added InfrastructureRepoURL to environment create and update functions ([bf07138](https://github.com/auto-staging/tower/commit/bf07138))
* added infrastructureRepoURL to structs ([6afe430](https://github.com/auto-staging/tower/commit/6afe430))
* added InfrastructureRepoURL to update repository model ([40e16f3](https://github.com/auto-staging/tower/commit/40e16f3))
* added invokation of builder to configure configure and remove schedules ([173691b](https://github.com/auto-staging/tower/commit/173691b))
* added lightning logger for configuration endpoints ([4f94f59](https://github.com/auto-staging/tower/commit/4f94f59))
* added omit empty to json annotations in structs ([780c141](https://github.com/auto-staging/tower/commit/780c141))
* added unique constraint condition for repository adding ([87caa9b](https://github.com/auto-staging/tower/commit/87caa9b))
* added validation for codeBuildRoleARN to environment update endpoint ([dfd7440](https://github.com/auto-staging/tower/commit/dfd7440))
* added variable for the codebuild role arn ([0006d9e](https://github.com/auto-staging/tower/commit/0006d9e))
* added webhook endpoint for github delete events ([061c43e](https://github.com/auto-staging/tower/commit/061c43e))
* configure lightning log based on lambda environment and reconfigure lambda environment ([954ac9a](https://github.com/auto-staging/tower/commit/954ac9a))
* invoke builder lambda for environemt add and delete ([46102e2](https://github.com/auto-staging/tower/commit/46102e2))
* limited dynamodb response to required parameters ([f121d0d](https://github.com/auto-staging/tower/commit/f121d0d))
* only update existing repositories and allow empty filters ([f2f3fe9](https://github.com/auto-staging/tower/commit/f2f3fe9))
* overwrite unset environment values with the parent repo configuration ([2e9c3cb](https://github.com/auto-staging/tower/commit/2e9c3cb))
* overwrite unset values in repository creation with global defaults ([0cb44a2](https://github.com/auto-staging/tower/commit/0cb44a2))
* project init with all used packages ([ecdcf2b](https://github.com/auto-staging/tower/commit/ecdcf2b))
* reflector app and get endpoint for configuration resource ([dd2e37f](https://github.com/auto-staging/tower/commit/dd2e37f))
* unfied and improvided error response and logging ([74491ae](https://github.com/auto-staging/tower/commit/74491ae))
