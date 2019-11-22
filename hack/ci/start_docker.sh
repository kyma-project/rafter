#!/usr/bin/env bash

set -eo pipefail

# external dependencies
readonly LIB_DIR="$(cd "${GOPATH}/src/github.com/kyma-project/test-infra/prow/scripts/lib" && pwd)"
CURRENT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)

source "${LIB_DIR}/docker.sh" || {
    echo 'Cannot load docker utilities.'
    exit 1
}

docker::start
