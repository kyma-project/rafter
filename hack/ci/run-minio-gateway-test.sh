#!/usr/bin/env bash

# Description: 
#   This scripts implements the flow of MinIO Gateway tests and MinIO Gateway Migration tests for CI.
#
# Required parameters:
#   - $1 - Absolute path to local Rafter repository
#   - $2 - Type of MinIO Gateway tests. Available values: basic, migration
# 
# Required env vars for gcs gateway:
#   - MINIO_GATEWAY_MODE - set to `gcs`
#   - CLOUDSDK_CORE_PROJECT - Name of the Google Cloud Platform (GCP) project for all GCP resources used in the tests
#   - GOOGLE_APPLICATION_CREDENTIALS - Absolute path to the Google Cloud Platform (GCP) Service Account Key file with the **Storage Admin** role
#
# Required env vars for azure gateway:
#   - MINIO_GATEWAY_MODE - set to `azure`
#   - BUILD_TYPE - Defines one of `pr/master/release`. This value is used to create the name of the Azure Storage Account.
#   - PULL_NUMBER - Defines pull request number. Required if BUILD_TYPE is set to `pr`.
#   - AZURE_RS_GROUP - Defines the name of the Azure Resource Group
#   - AZURE_REGION - Azure region code
#   - AZURE_SUBSCRIPTION_ID - ID of the the Azure Subscription
#   - AZURE_SUBSCRIPTION_APP_ID - App ID of the Azure Subscription
#   - AZURE_SUBSCRIPTION_SECRET - Credentials for the Azure Subscription
#   - AZURE_SUBSCRIPTION_TENANT - Tenant ID of the Azure Subscription

set -o errexit
set -o nounset
set -o pipefail

readonly ROOT_REPO_PATH="${1}"
readonly TEST_TYPE="${2}"

readonly MINIO_GATEWAY_TEST_BASIC="basic"
readonly MINIO_GATEWAY_TEST_MIGRATION="migration"

readonly MINIO_GATEWAY_PROVIDER_GCS="gcs"
readonly MINIO_GATEWAY_PROVIDER_AZURE="azure"

readonly CLUSTER_NAME="ci-minio-gateway-test"

readonly TMP_DIR="$(mktemp -d)"
readonly TMP_BIN_DIR="${TMP_DIR}/bin"
mkdir -p "${TMP_BIN_DIR}"
export PATH="${TMP_BIN_DIR}:${PATH}"

export ARTIFACTS_DIR="${ARTIFACTS:-"${TMP_DIR}/artifacts"}"
mkdir -p "${ARTIFACTS_DIR}"

init_environment() {
  if [[ -z ${MINIO_GATEWAY_MODE-} ]]; then
    echo '- $MINIO_GATEWAY_MODE variable is not set.'
    return 1
  fi
  local -r current_dir=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)

  source "${current_dir}/envs.sh" || {
    echo '- Cannot load environment variables.'
    return 1
  }
  source "${current_dir}/lib/test-helpers.sh" || {
    echo '- Cannot load test helpers.'
    return 1
  }
  source "${current_dir}/lib/minio/gateway-helpers.sh" || {
    echo '- Cannot load gateway helpers.'
    return 1
  }

  if [[ ${TEST_TYPE} = "${MINIO_GATEWAY_TEST_BASIC}" ]]; then
    source "${current_dir}/lib/minio/gateway-basic.sh" || {
      echo '- Cannot load gateway-basic test suite.'
      return 1
    }
  elif [[ ${TEST_TYPE} = "${MINIO_GATEWAY_TEST_MIGRATION}" ]] ; then
    source "${current_dir}/lib/minio/gateway-migration.sh" || {
      echo '- Cannot load gateway-migration test suite.'
      return 1
    }
  else
    log::error "- Not supported test type - ${TEST_TYPE}."
    return 1
  fi

  if [[ "${MINIO_GATEWAY_MODE}" = "${MINIO_GATEWAY_PROVIDER_GCS}" ]]; then
    command -v gsutil >/dev/null 2>&1 || { 
      log::error "- gsutil is reguired it's not installed. Aborting."
      return 1
    }
  elif [[ "${MINIO_GATEWAY_MODE}" = "${MINIO_GATEWAY_PROVIDER_AZURE}" ]]; then
    command -v az >/dev/null 2>&1 || { 
      log::error "- azure-cli is reguired but it's not installed. Aborting."
      return 1
    }
  else
    log::error "- Not supported MinIO Gateway mode - ${MINIO_GATEWAY_MODE}."
    return 1
  fi

  gatewayHelpers::check_gateway_mode "${MINIO_GATEWAY_MODE}"
  gateway::validate_environment
  if [[ "${MINIO_GATEWAY_MODE}" = "${MINIO_GATEWAY_PROVIDER_AZURE}" ]]; then
    azureGateway::create_storage_account_name
  fi

  docker::start
}

finalize() {
  local -r exit_status=$?
  local finalization_failed="false"

  gateway::after_test

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

  log::info "- Deleting directory with temporary binaries used in tests..."
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

  junit::suite_init "Rafter_Gateway"
  trap junit::test_fail ERR

  local -r minio_secret_name="rafter-minio"
  local -r release_name="rafter"
  local -r host_os="$(host::os)"

  local -r tmp_rafter_charts_dir="${TMP_DIR}/${RAFTER_CHART}"
  mkdir -p "${tmp_rafter_charts_dir}"

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

  if [[ ${TEST_TYPE} = "${MINIO_GATEWAY_TEST_BASIC}" ]]; then
    junit::test_start "MinIO_Gateway_Tests"
    gatewayBasic::run "${release_name}" "${minio_secret_name}" "${INGRESS_ADDRESS}" "${tmp_rafter_charts_dir}" 2>&1 | junit::test_output
    junit::test_pass
  elif [[ ${TEST_TYPE} = "${MINIO_GATEWAY_TEST_MIGRATION}" ]] ; then
    junit::test_start "Install_Rafter"
    testHelpers::install_rafter "${release_name}" "${minio_secret_name}" "${INGRESS_ADDRESS}" "${tmp_rafter_charts_dir}" 2>&1 | junit::test_output
    junit::test_pass
    
    junit::test_start "MinIO_Gateway_Migration_Tests"
    gatewayMigration::run "${release_name}" "${MINIO_ADDRESS}" "${minio_secret_name}" "${tmp_rafter_charts_dir}" 2>&1 | junit::test_output
    junit::test_pass
  else
    log::error "- Not supported test type - ${TEST_TYPE}."
    exit 1
  fi

  junit::test_start "Rafter_Integration_Test"
  testHelpers::run_integration_tests "${ROOT_REPO_PATH}" "${CLUSTER_NAME}" "${MINIO_ADDRESS}" "${UPLOAD_SERVICE_ENDPOINT}" "${MINIO_GATEWAY_SECRET_NAME}" 2>&1 | junit::test_output
  junit::test_pass
}
main
