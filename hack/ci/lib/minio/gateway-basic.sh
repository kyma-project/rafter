#!/usr/bin/env bash

# Arguments:
#   $1 - Release name
#   $2 - The MiniIO k8s secret name
#   $3 - The addres of the ingress that exposes upload and minio endpoints
#   $4 - Path to charts directory
gatewayBasic::run() {
  local -r release_name="${1}"
  local -r minio_secret_name="${2}"
  local -r ingress_address="${3}"
  local -r charts_path="${4}"

  # configure provider
  gateway::before_test "${minio_secret_name}"

  # install Rafter with MinIO gateway mode
  gateway::install "${release_name}" "${ingress_address}" "${charts_path}"
}
