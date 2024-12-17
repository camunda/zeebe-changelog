# Zeebe Changelog

[![Build Status](https://travis-ci.com/camunda/zeebe-changelog.svg?branch=master)](https://travis-ci.com/camunda/zeebe-changelog)
[![Go Report Card](https://goreportcard.com/badge/github.com/camunda/zeebe-changelog?style=flat-square)](https://goreportcard.com/report/github.com/camunda/zeebe-changelog)
[![Release](https://img.shields.io/github/release/camunda/zeebe-changelog.svg?style=flat-square)](https://github.com/camunda/zeebe-changelog/releases/latest)
[![codecov](https://codecov.io/gh/camunda/zeebe-changelog/branch/master/graph/badge.svg)](https://codecov.io/gh/camunda/zeebe-changelog)

Generate changelog for [Zeebe](github.com/camunda/camunda) project.


## Example usage

```sh
  export ZCL_FROM_REV=PREV_VERSION
  export ZCL_TARGET_REV=TARGET_VERSION

  # This will add labels to the issues in GitHub. You can verify this step by looking at closed issues. They should now be tagged with the release.
  zcl add-labels \
    --token=$GITHUB_TOKEN \
    --from=$ZCL_FROM_REV \
    --target=$ZCL_TARGET_REV \
    --label="version:$ZCL_TARGET_REV" \
    --org camunda --repo zeebe

  # This command will print markdown code to the console. You will need to manually insert this output into the release draft.
  zcl generate \
     --token=$GITHUB_TOKEN \
     --label="version:$ZCL_TARGET_REV" \
     --org camunda --repo zeebe
```
## Release ZCL

* [Prerequisite] Install [goreleaser](https://goreleaser.com/intro/#usage)
  * We have experienced issues with the recent versions (likely the project is not compatible with the recent versions)
  * To overcome this we used (in the last releases): `go install github.com/goreleaser/goreleaser@v1.0.0`
* Create a new tag with the latest changes:
  * Create tag local: `git tag <version>`
  * Push tag: `git push origin <tag>`
* Release ZCL
  * Run goreleaser: `$GOPATH/bin/goreleaser release`
  * Verify on [release page](https://github.com/camunda/zeebe-changelog/releases)

