#!/usr/bin/env bash

# Init environment for gateway test
# Arguments:
#   $1 - The name of provider
gatewayHelpers::init_environment() {
  local -r mode="${1}"
  
  __check_gateway_mode "${mode}"
}

# Outputs:
#   $1 - The public bucket name
#   $1 - The private bucket name
gatewayHelpers::get_bucket_names() {
  local -r cm_name="rafter-upload-service"

  local -r public_bucket=$(kubectl -n default get configmap ${cm_name} -o jsonpath="{.data.public}" | xargs -n1 echo)
  local -r private_bucket=$(kubectl -n default get configmap ${cm_name} -o jsonpath="{.data.private}" | xargs -n1 echo)

  echo "${public_bucket}" "${private_bucket}"
}

# Check kind of gateway mode.
# Arguments:
#   $1 - The name of provider
__check_gateway_mode() {
  local -r current_dir=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
  local -r mode="${1}"

  log::info "- Checking MinIO Gateway mode..."

  if [[ "${mode}" = "${__MINIO_GATEWAY_PROVIDER_GCS__}" ]]; then
    source "${current_dir}/gcs-gateway.sh" || {
      log::error '- Cannot load gcs-gateway test suite.'
      exit 1
    }
  elif [[ "${mode}" = "${__MINIO_GATEWAY_PROVIDER_AZURE__}" ]]; then
    source "${current_dir}/azure-gateway.sh" || {
      log::error '- Cannot load azure-gateway test suite.'
      exit 1
    }
  else
    log::error "- Not supported MinIO Gateway mode - ${mode}."
    exit 1
  fi

  log::success "- Running MinIO on ${mode} Gateway mode."
}
