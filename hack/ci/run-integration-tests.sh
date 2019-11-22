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

source "${CURRENT_DIR}/envs.sh" || {
    echo 'Cannot load environment variables.'
    exit 1
}

source "${CURRENT_DIR}/test-helper.sh" || {
    echo 'Cannot load test helper.'
    exit 1
}

main() {
    # minio access key that will be used during rafter installation
    local -r MINIO_ACCESSKEY=4j4gEuRH96ZFjptUFeFm
    # minio secret key that will be used during the rafter installation
    local -r MINIO_SECRETKEY=UJnce86xA7hK01WblDdbmXg4gwjKwpFypdLJCvJ3
    # kind cluster configuration
    local -r CLUSTER_CONFIG=${CURRENT_DIR}/config/kind/cluster-config.yaml
    # the addres of the ingress that exposes upload and minio endpoints
    local -r INGRESS_ADDRESS=http://localhost:30080

    trap "testHelper::cleanup ${CLUSTER_NAME} ${TMP_DIR}" EXIT
    
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
    
    testHelper::load_images "${CLUSTER_NAME}" "${UPLOADER_IMG_NAME}" "${MANAGER_IMG_NAME}" "${FRONT_MATTER_IMG_NAME}" "${ASYNCAPI_IMG_NAME}"
    
    testHelper::install_rafter "${MINIO_ACCESSKEY}" "${MINIO_SECRETKEY}" "${INGRESS_ADDRESS}"
    
    testHelper::start_integration_tests "${CLUSTER_NAME}" "${MINIO_ACCESSKEY}" "${MINIO_SECRETKEY}" "${INGRESS_ADDRESS}"
}

main