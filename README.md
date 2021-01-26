# expressive

[![CircleCI](https://circleci.com/gh/CarlCui/expressive/tree/master.svg?style=svg)](https://circleci.com/gh/CarlCui/expressive/tree/master)

## LLVM
Current llvm version is v10.0.0.

### Upgrade llvm locally

1. go to https://releases.llvm.org/download.html
1. Download binary for MacOS
1. extract tar ball
1. copy to `/usr/local/llvm`

### Upgrade pre-built lli

1. go to https://releases.llvm.org/download.html
1. Download binary for linux ubuntu
1. extract tar ball
1. copy `bin/lli` to `root/llvm-prebuilt`
1. commit the change
## To run e2e

1. Switch to local development llvm: `LLI_PATH=lli # for local development`
1. `bash build.sh`
1. `bash test_e2e.sh`

## To debug circleci build locally

* circleci config validate
* circleci local execute --job JOB_NAME (build)
