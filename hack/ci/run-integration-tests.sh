#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail
set -e

CURRENT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)

source "${CURRENT_DIR}/test-helper.sh" || {
    echo 'Cannot load test helper.'
    exit 1
}

main() {
    trap testHelper::cleanup EXIT
    
    kind::create_cluster \
    "${CLUSTER_NAME}" \
    "${STABLE_KUBERNETES_VERSION}" \
    "${CLUSTER_CONFIG}" 2>&1
    
    testHelper::check_version

    testHelper::install_tiller

    testHelper::add_repos_and_update
    
    testHelper::install_ingress
    
    testHelper::load_images
    
    testHelper::install_rafter
    
    testHelper::apply_ingress
    
    testHelper::start_integration_tests
}

main