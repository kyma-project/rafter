#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

# docker images to load into kind
readonly UPLOADER_IMG_NAME="${1}"
readonly MANAGER_IMG_NAME="${2}"
readonly FRONT_MATTER_IMG_NAME="${3}"
readonly ASYNCAPI_IMG_NAME="${4}"

CURRENT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)

readonly TMP_DIR="$(mktemp -d)"
readonly TMP_BIN_DIR="${TMP_DIR}/bin"
mkdir -p "${TMP_BIN_DIR}"
export PATH="${TMP_BIN_DIR}:${PATH}"


source "${CURRENT_DIR}/test-helper.sh" || {
    echo 'Cannot load test helper.'
    exit 1
}

source "${CURRENT_DIR}/envs.sh" || {
    echo 'Cannot load environment variables.'
    exit 1
}

main() {
    trap "testHelper::cleanup ${TMP_DIR}" EXIT
    
    infraHelper::install_helm_tiller "${STABLE_HELM_VERSION}" "$(host::os)" "${TMP_BIN_DIR}"
    infraHelper::install_kind "${STABLE_KIND_VERSION}" "$(host::os)" "${TMP_BIN_DIR}"
    
    kubernetes::ensure_kubectl "${STABLE_KUBERNETES_VERSION}" "$(host::os)" "${TMP_BIN_DIR}"
    
    kind::create_cluster \
    "${CLUSTER_NAME}" \
    "${STABLE_KUBERNETES_VERSION}" \
    "${CLUSTER_CONFIG}"
    
    testHelper::install_tiller

    testHelper::add_repos_and_update
    
    testHelper::install_ingress
    
    testHelper::load_images "${UPLOADER_IMG_NAME}" "${MANAGER_IMG_NAME}" "${FRONT_MATTER_IMG_NAME}" "${ASYNCAPI_IMG_NAME}"
    
    testHelper::install_rafter
    
    testHelper::start_integration_tests
}

main