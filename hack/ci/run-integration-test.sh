#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

__init_environment() {
  local -r current_dir=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)

  source "${current_dir}/envs.sh" || {
    echo '- Cannot load environment variables.'
    exit 1
  }
  source "${current_dir}/lib/test-helpers.sh" || {
    echo '- Cannot load test helpers.'
    exit 1
  }
}

# Arguments:
#   $1 - Name of the kind cluster
#   $2 - Tmp directory with binaries used during test
#   $3 - Artifacts dir used to store JUnit reports
__finalize() {
  local -r kind_cluster_name="${1}"
  local -r bin_dir="${2}"
  local -r artifacts_dir="${3}"
  local -r exit_status=$?
  local finalization_failed="false"

  junit::test_start "Finalization"
  log::info "Finalizing job" 2>&1 | junit::test_output

  log::info "- Printing all docker processes..." 2>&1 | junit::test_output
  docker::print_processes 2>&1 | junit::test_output || finalization_failed="true"
    
  log::info "- Exporting cluster logs to ${artifacts_dir}..." 2>&1 | junit::test_output
  kind::export_logs "${kind_cluster_name}" "${artifacts_dir}" 2>&1 | junit::test_output || finalization_failed="true"

  log::info "- Cleaning up cluster ${kind_cluster_name}..." | junit::test_output
  kind::delete_cluster "${kind_cluster_name}" 2>&1 | junit::test_output || finalization_failed="true"

  if [[ ${finalization_failed} = "true" ]]; then
    junit::test_fail || true
  else
    junit::test_pass
  fi
  junit::suite_save

  log::info "- Deleting directory with temporary binaries and charts used in tests..."
  rm -rf "${bin_dir}" || true

  if [[ ${exit_status} -eq 0 ]]; then
    log::success "- Job finished with success"
  else
    log::error "- Job finished with error"
  fi

  return "${exit_status}"
}

# Arguments:
#   $1 - $4 - Name of images to load in kind cluster
main() {
  __init_environment

  junit::suite_init "Rafter_Integration"
  trap junit::test_fail ERR

  local -r controller_manager_img="${1}"
  local -r upload_service_img="${2}"
  local -r front_matter_service_img="${3}"
  local -r asyncapi_service_img="${4}"

  local -r tmp_dir="$(mktemp -d)"
  local -r tmp_bin_dir="${tmp_dir}/bin"
  mkdir -p "${tmp_bin_dir}"
  export PATH="${tmp_bin_dir}:${GOPATH}:${PATH}"

  export ARTIFACTS_DIR="${ARTIFACTS:-"${tmp_dir}/artifacts"}"
  mkdir -p "${ARTIFACTS_DIR}"

  local -r temp_rafter_charts_dir="${tmp_dir}/${__RAFTER__}"
  mkdir -p "${temp_rafter_charts_dir}"

  local -r cluster_name="ci-integration-test"
  local -r host_os="$(host::os)"
  local -r release_name="rafter"
  local -r minio_secret_name="rafter-minio"

  trap "__finalize ${cluster_name} ${tmp_dir} ${ARTIFACTS_DIR}" EXIT

  junit::test_start "Install_go_junit_report"
  testHelpers::install_go_junit_report 2>&1 | junit::test_output
  junit::test_pass

  junit::test_start "Install_Helm_Tiller"
  testHelpers::download_helm_tiller "${__STABLE_HELM_VERSION__}" "${host_os}" "${tmp_bin_dir}" 2>&1 | junit::test_output
  junit::test_pass

  junit::test_start "Install_Kind"
  testHelpers::download_kind "${__STABLE_KIND_VERSION__}" "${host_os}" "${tmp_bin_dir}" 2>&1 | junit::test_output
  junit::test_pass
  
  junit::test_start "Install_Kubectl"
  testHelpers::download_kubectl "${__STABLE_KUBERNETES_VERSION__}" "${host_os}" "${tmp_bin_dir}" 2>&1 | junit::test_output
  junit::test_pass

  junit::test_start "Create_Kind_Cluster"
  testHelpers::create_cluster "${cluster_name}" "${__STABLE_KUBERNETES_VERSION__}" "${__CLUSTER_CONFIG_FILE__}" 2>&1 | junit::test_output
  junit::test_pass

  junit::test_start "Install_Tiller"
  testHelpers::install_tiller 2>&1 | junit::test_output
  junit::test_pass

  junit::test_start "Prepare_Local_Helm_Charts"
  testHelpers::prepare_local_helm_charts "${temp_rafter_charts_dir}" 2>&1 | junit::test_output
  junit::test_pass
  
  junit::test_start "Install_Ingress"
  testHelpers::install_ingress "${__INGRESS_YAML_FILE__}" 2>&1 | junit::test_output
  junit::test_pass
  
  junit::test_start "Load_Images"
  testHelpers::load_images "${cluster_name}" "${controller_manager_img}" "${upload_service_img}" "${front_matter_service_img}" "${asyncapi_service_img}" 2>&1 | junit::test_output
  junit::test_pass
  
  junit::test_start "Create_MinIO_K8S_Secret"
  testHelpers::create_minio_k8s_secret "${minio_secret_name}" "${__DEFAULT_MINIO_ACCESS_KEY__}" "${__DEFAULT_MINIO_SECRET_KEY__}" 2>&1 | junit::test_output
  junit::test_pass

  junit::test_start "Install_Rafter"
  testHelpers::install_rafter "${release_name}" "${minio_secret_name}" "${__INGRESS_ADDRESS__}" "${temp_rafter_charts_dir}" 2>&1 | junit::test_output
  junit::test_pass 

  junit::test_start "Rafter_Integration_Test"
  testHelpers::run_integration_tests "${cluster_name}" "${__MINIO_ADDRESS__}" "${__UPLOAD_SERVICE_ENDPOINT__}" "${minio_secret_name}" "${ARTIFACTS_DIR}" 2>&1 | junit::test_output
  junit::test_pass
}
main "$@"
