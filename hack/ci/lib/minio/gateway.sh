#!/usr/bin/env bash

readonly MINIO_GATEWAY_MODE="${1}"

readonly CURRENT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)

gateway::loadSource() {
  if [[ "${MINIO_GATEWAY_MODE}" = "gcs" ]]; then
      source "${TEST_INFRA_SOURCES_DIR}/prow/scripts/cluster-integration/minio/gcs-gateway.sh"
  elif [[ "${MINIO_GATEWAY_MODE}" = "azure" ]]; then
      source "${TEST_INFRA_SOURCES_DIR}/prow/scripts/cluster-integration/minio/azure-gateway.sh"
  else
      shout "Not supported Minio Gateway mode - ${MINIO_GATEWAY_MODE}"
      exit 1
  fi
}

gateway::start() {

}