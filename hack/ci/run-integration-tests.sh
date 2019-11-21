#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail
set -e

CURRENT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)

readonly TMP_DIR="$(mktemp -d)"
readonly TMP_BIN_DIR="${TMP_DIR}/bin"
mkdir -p "${TMP_BIN_DIR}"
export PATH="${TMP_BIN_DIR}:${PATH}"


source "${CURRENT_DIR}/test-helper.sh" || {
    echo 'Cannot load test helper.'
    exit 1
}

main() {
    trap testHelper::cleanup EXIT
    
    infraHelper::install_helm_tiller "$(host::os)" "${TMP_BIN_DIR}"
    infraHelper::install_kind "$(host::os)" "${TMP_BIN_DIR}"
    
    kubernetes::ensure_kubectl "${STABLE_KUBERNETES_VERSION}" "$(host::os)" "${TMP_BIN_DIR}"
    
    kind::create_cluster \
    "${CLUSTER_NAME}" \
    "${STABLE_KUBERNETES_VERSION}" \
    "${CLUSTER_CONFIG}" 2>&1
    
    testHelper::install_tiller

    testHelper::add_repos_and_update
    
    testHelper::install_ingress
    
    testHelper::load_images
    
    testHelper::install_rafter
    
    testHelper::start_integration_tests
}

main