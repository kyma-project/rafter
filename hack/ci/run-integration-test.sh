#!/usr/bin/env bash

# Description: 
#   This scripts implements the flow of Rafter Integration tests for CI.
#
# Required parameters:
#   - $1 - Absolute path to local Rafter repository

set -o errexit
set -o nounset
set -o pipefail

readonly ROOT_REPO_PATH="${1}"

readonly CLUSTER_NAME="ci-integration-test"

readonly TMP_DIR="$(mktemp -d)"
readonly TMP_BIN_DIR="${TMP_DIR}/bin"
mkdir -p "${TMP_BIN_DIR}"
export PATH="${TMP_BIN_DIR}:${PATH}"

export ARTIFACTS_DIR="${ARTIFACTS:-"${TMP_DIR}/artifacts"}"
mkdir -p "${ARTIFACTS_DIR}"

init_environment() {
  local -r current_dir=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)

  source "${current_dir}/envs.sh" || {
    echo '- Cannot load environment variables.'
    return 1
  }
  source "${current_dir}/lib/test-helpers.sh" || {
    echo '- Cannot load test helpers.'
    return 1
  }

  docker::start
}

finalize() {
  local -r exit_status=$?
  local finalization_failed="false"

  junit::test_start "Finalization"
  log::info "Finalizing job" 2>&1 | junit::test_output

  log::info "- Printing all docker processes..." 2>&1 | junit::test_output
  docker::print_processes 2>&1 | junit::test_output || finalization_failed="true"
    
  log::info "- Exporting cluster logs to ${ARTIFACTS_DIR}..." 2>&1 | junit::test_output
  kind::export_logs "${CLUSTER_NAME}" "${ARTIFACTS_DIR}" 2>&1 | junit::test_output || finalization_failed="true"

  log::info "- Cleaning up cluster ${CLUSTER_NAME}..." | junit::test_output
  kind::delete_cluster "${CLUSTER_NAME}" 2>&1 | junit::test_output || finalization_failed="true"

  if [[ ${finalization_failed} = "true" ]]; then
    junit::test_fail || true
  else
    junit::test_pass
  fi
  junit::suite_save

  log::info "- Deleting directory with temporary binaries and charts used in tests..."
  rm -rf "${TMP_DIR}" || true

  if [[ ${exit_status} -eq 0 ]]; then
    log::success "- Job finished with success"
  else
    log::error "- Job finished with error"
  fi

  return "${exit_status}"
}
trap "finalize" EXIT

main() {
  init_environment

  junit::suite_init "Rafter_Integration"
  trap junit::test_fail ERR

  local -r host_os="$(host::os)"
  local -r minio_secret_name="rafter-minio"

  local -r tmp_rafter_charts_dir="${TMP_DIR}/${RAFTER_CHART}"
  mkdir -p "${tmp_rafter_charts_dir}"

  junit::test_start "Install_go_junit_report"
  testHelpers::install_go_junit_report 2>&1 | junit::test_output
  junit::test_pass

  junit::test_start "Install_Helm_Tiller"
  testHelpers::download_helm_tiller "${STABLE_HELM_VERSION}" "${host_os}" "${TMP_BIN_DIR}" 2>&1 | junit::test_output
  junit::test_pass

  junit::test_start "Install_Kind"
  testHelpers::download_kind "${STABLE_KIND_VERSION}" "${host_os}" "${TMP_BIN_DIR}" 2>&1 | junit::test_output
  junit::test_pass
  
  junit::test_start "Install_Kubectl"
  kubernetes::ensure_kubectl "${STABLE_KUBERNETES_VERSION}" "${host_os}" "${TMP_BIN_DIR}" 2>&1 | junit::test_output
  junit::test_pass

  junit::test_start "Create_Kind_Cluster"
  kind::create_cluster \
    "${CLUSTER_NAME}" \
    "${STABLE_KUBERNETES_VERSION}" \
    "${CLUSTER_CONFIG_FILE}" 2>&1 | junit::test_output
  junit::test_pass

  junit::test_start "Install_Tiller"
  testHelpers::install_tiller 2>&1 | junit::test_output
  junit::test_pass

  junit::test_start "Prepare_Local_Helm_Charts"
  testHelpers::prepare_local_helm_charts "${ROOT_REPO_PATH}" "${tmp_rafter_charts_dir}" 2>&1 | junit::test_output
  junit::test_pass
  
  junit::test_start "Install_Ingress"
  testHelpers::install_ingress "${INGRESS_YAML_FILE}" 2>&1 | junit::test_output
  junit::test_pass
  
  junit::test_start "Load_Images"
  testHelpers::load_rafter_images "${CLUSTER_NAME}" 2>&1 | junit::test_output
  junit::test_pass
  
  junit::test_start "Create_K8S_Secret_For_MinIO"
  testHelpers::create_k8s_secret "${minio_secret_name}" "${DEFAULT_MINIO_ACCESS_KEY}" "${DEFAULT_MINIO_SECRET_KEY}" 2>&1 | junit::test_output
  junit::test_pass

  junit::test_start "Install_Rafter"
  testHelpers::install_rafter "rafter" "${minio_secret_name}" "${INGRESS_ADDRESS}" "${tmp_rafter_charts_dir}" 2>&1 | junit::test_output
  junit::test_pass 

  junit::test_start "Rafter_Integration_Test"
  testHelpers::run_integration_tests "${ROOT_REPO_PATH}" "${CLUSTER_NAME}" "${MINIO_ADDRESS}" "${UPLOAD_SERVICE_ENDPOINT}" "${minio_secret_name}" "${ARTIFACTS_DIR}" 2>&1 | junit::test_output
  junit::test_pass
}
main
