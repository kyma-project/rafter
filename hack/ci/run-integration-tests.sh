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

readonly ARTIFACTS_DIR="${ARTIFACTS:-"${TMP_DIR}/artifacts"}"
mkdir -p "${ARTIFACTS_DIR}"

source "${CURRENT_DIR}/envs.sh" || {
    echo 'Cannot load environment variables.'
    exit 1
}

source "${CURRENT_DIR}/test-helper.sh" || {
    echo 'Cannot load test helper.'
    exit 1
}


# finalize stores logs, saves JUnit report and removes cluster
function finalize {
    local -r exit_status=$?
    local finalization_failed="false"

    junit::test_start "Finalization"
    log::info "Finalizing job" 2>&1 | junit::test_output

    log::info "Printing all docker processes" 2>&1 | junit::test_output
    docker::print_processes 2>&1 | junit::test_output || finalization_failed="true"
    
    log::info "Exporting cluster logs to ${ARTIFACTS_DIR}" 2>&1 | junit::test_output
    kind::export_logs "${CLUSTER_NAME}" "${ARTIFACTS_DIR}" 2>&1 | junit::test_output || finalization_failed="true"

    log::info "Deleting cluster" 2>&1 | junit::test_output
    kind::delete_cluster "${CLUSTER_NAME}" 2>&1 | junit::test_output || finalization_failed="true"
    
    
    if [[ ${finalization_failed} = "true" ]]; then
        junit::test_fail || true
    else
        junit::test_pass
    fi

    junit::suite_save

    log::info "Deleting temporary dir ${TMP_DIR}"
    rm -rf "${TMP_DIR}" || true

    if [[ ${exit_status} -eq 0 ]]; then
        log::success "Job finished with success"
    else
        log::error "Job finished with error"
    fi

    return "${exit_status}"
}

trap finalize EXIT

main() {
    junit::suite_init "Rafter_Integration"
    trap junit::test_fail ERR
    # minio access key that will be used during rafter installation
    local -r MINIO_ACCESSKEY=4j4gEuRH96ZFjptUFeFm
    # minio secret key that will be used during the rafter installation
    local -r MINIO_SECRETKEY=UJnce86xA7hK01WblDdbmXg4gwjKwpFypdLJCvJ3
    # kind cluster configuration
    local -r CLUSTER_CONFIG=${CURRENT_DIR}/config/kind/cluster-config.yaml
    # the addres of the ingress that exposes upload and minio endpoints
    local -r INGRESS_ADDRESS=http://localhost:30080

    junit::test_start "Install_go_junit_report"
    testHelper::install_go_junit_report   2>&1 | junit::test_output
    junit::test_pass

    junit::test_start "Install_Helm_Tiller"
    infraHelper::install_helm_tiller "${STABLE_HELM_VERSION}" "$(host::os)" "${TMP_BIN_DIR}"  2>&1 | junit::test_output
    junit::test_pass

    junit::test_start "Install_Kind"
    infraHelper::install_kind "${STABLE_KIND_VERSION}" "$(host::os)" "${TMP_BIN_DIR}" 2>&1 | junit::test_output
    junit::test_pass

    junit::test_start "Install_Kubectl"
    kubernetes::ensure_kubectl "${STABLE_KUBERNETES_VERSION}" "$(host::os)" "${TMP_BIN_DIR}" 2>&1 | junit::test_output
    junit::test_pass

    junit::test_start "Create_Kind_Cluster"
    kind::create_cluster \
    "${CLUSTER_NAME}" \
    "${STABLE_KUBERNETES_VERSION}" \
    "${CLUSTER_CONFIG}" 2>&1 | junit::test_output
    junit::test_pass
    
    junit::test_start "Install_Tiller"
    testHelper::install_tiller 2>&1 | junit::test_output
    junit::test_pass

    junit::test_start "Helm_Add_Repo_And_Update"
    testHelper::add_repos_and_update 2>&1 | junit::test_output
    junit::test_pass

    junit::test_start "Install_Ingress"
    testHelper::install_ingress 2>&1 | junit::test_output
    junit::test_pass

    junit::test_start "Load_Images"
    testHelper::load_images "${CLUSTER_NAME}" "${UPLOADER_IMG_NAME}" "${MANAGER_IMG_NAME}" "${FRONT_MATTER_IMG_NAME}" "${ASYNCAPI_IMG_NAME}" 2>&1 | junit::test_output
    junit::test_pass

    junit::test_start "Install_Rafter"
    testHelper::install_rafter "${MINIO_ACCESSKEY}" "${MINIO_SECRETKEY}" "${INGRESS_ADDRESS}" 2>&1 | junit::test_output
    junit::test_pass

    junit::test_start "Rafter_Integration_Test"
    testHelper::start_integration_tests "${CLUSTER_NAME}" "${MINIO_ACCESSKEY}" "${MINIO_SECRETKEY}" "${INGRESS_ADDRESS}" "${ARTIFACTS_DIR}" 2>&1 | junit::test_output
    junit::test_pass
}

main