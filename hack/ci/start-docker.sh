#!/usr/bin/env bash

function start_docker() {
  local -r lib_dir="$(cd "${GOPATH}/src/github.com/kyma-project/test-infra/prow/scripts/lib" && pwd)"

  source "${lib_dir}/docker.sh" || {
    echo 'Cannot load docker utilities.'
    exit 1
  }

  docker::start
}
start_docker
