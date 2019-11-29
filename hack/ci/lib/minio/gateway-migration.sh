#!/usr/bin/env bash

readonly MINIO_GATEWAY_MIGRATION_SAMPLE_CONTENT="sample"
readonly MINIO_GATEWAY_MIGRATION_CONTENT_TYPE="application/octet-stream"
readonly MINIO_GATEWAY_MIGRATION_SAMPLE_FILE="sample"
readonly MINIO_GATEWAY_MIGRATION_SAMPLE_FILE_WITH_DIR="sampledir/sample"

# Creates sample file and uploads it to minio.
# Arguments:
#   $1 - The name of the bucket the file will be uploaded to
#   $2 - The name of the created file
#   $3 - MiniIO host
#   $4 - MiniIO access key
#   $5 - MiniIO secret key
gatewayMigration::upload_sample_file_to_minio() {
  local -r bucket_name="${1}"
  local -r file_name="${2}"
  local -r minio_host="${3}"
  local -r minio_accessKey="${4}"
  local -r minio_secretKey="${5}"

  local -r resource="${bucket_name}"/"${file_name}"

  local -r date=$(date -R)
  local -r signature=PUT"\n\n${MINIO_GATEWAY_MIGRATION_CONTENT_TYPE}\n${date}\n/${resource}"
  local -r checksum=$(echo -en "${signature}" | openssl sha1 -hmac "${minio_secretKey}" -binary | base64)

  log::info "- Uploading ${resource} to minio ${minio_host}..."

  echo "${MINIO_GATEWAY_MIGRATION_SAMPLE_CONTENT}" | curl -X PUT -d @- \
    -H "Date: ${date}" \
    -H "Content-Type: ${MINIO_GATEWAY_MIGRATION_CONTENT_TYPE}" \
    -H "Authorization: AWS ${minio_accessKey}:${checksum}" \
    --insecure \
    --silent \
    --fail \
    "${minio_host}"/"${resource}"

  log::success "- Uploaded ${resource}."
}

# Download samples from minIO to verify if the migration was successful.
# Arguments:
#   $1 - The name of the bucket the file will be uploaded to
#   $2 - The name of the created file
#   $3 - MiniIO host
#   $4 - MiniIO access key
#   $5 - MiniIO secret key
gatewayMigration::download_sample_file_from_minio() {
  local -r bucket_name="${1}"
  local -r file_name="${2}"
  local -r minio_host="${3}"
  local -r minio_accessKey="${4}"
  local -r minio_secretKey="${5}"

  local -r resource="${bucket_name}"/"${file_name}"
  
  local -r date=$(date -R)
  local -r signature=GET"\n\n${MINIO_GATEWAY_MIGRATION_CONTENT_TYPE}\n${date}\n/${resource}"
  local -r checksum=$(echo -en "${signature}" | openssl sha1 -hmac "${minio_secretKey}" -binary | base64)
  
  log::info "- Downloading ${resource} from minio ${minio_host}..."

  curl -H "Date: ${date}" \
    -H "Content-Type: ${MINIO_GATEWAY_MIGRATION_CONTENT_TYPE}" \
    -H "Authorization: AWS ${minio_accessKey}:${checksum}" \
	  --silent \
	  --insecure \
	  --fail \
    "${minio_host}"/"${resource}" > /dev/null

  log::success "- Downloaded ${resource}."
}

# Arguments:
#   $1 - Minio host
#   $2 - The MiniIO k8s secret name
gatewayMigration::before_migration() {
  local -r minio_host="${1}"
  local -r minio_secret_name="${2}"

  local public_bucket=""
  local private_bucket=""
  read public_bucket private_bucket < <(gatewayHelpers::get_bucket_names)

  local minio_accessKey=""
  local minio_secretKey=""
  read minio_accessKey minio_secretKey < <(testHelpers::get_k8s_secret_data ${minio_secret_name})
  
  log::info "- Uploading files to MinIO.."

  gatewayMigration::upload_sample_file_to_minio "${public_bucket}" "${MINIO_GATEWAY_MIGRATION_SAMPLE_FILE}" "${minio_host}" "${minio_accessKey}" "${minio_secretKey}"
  gatewayMigration::upload_sample_file_to_minio "${public_bucket}" "${MINIO_GATEWAY_MIGRATION_SAMPLE_FILE_WITH_DIR}" "${minio_host}" "${minio_accessKey}" "${minio_secretKey}"

  gatewayMigration::upload_sample_file_to_minio "${private_bucket}" "${MINIO_GATEWAY_MIGRATION_SAMPLE_FILE}" "${minio_host}" "${minio_accessKey}" "${minio_secretKey}"
  gatewayMigration::upload_sample_file_to_minio "${private_bucket}" "${MINIO_GATEWAY_MIGRATION_SAMPLE_FILE_WITH_DIR}" "${minio_host}" "${minio_accessKey}" "${minio_secretKey}"

  log::success "- Uploaded files to MinIO."
}

# Arguments:
#   $1 - Minio host
#   $2 - The MiniIO k8s secret name
gatewayMigration::after_migration() {
  local -r minio_host="${1}"
  local -r minio_secret_name="${2}"

  local public_bucket=""
  local private_bucket=""
  read public_bucket private_bucket < <(gatewayHelpers::get_bucket_names)

  local minio_accessKey=""
  local minio_secretKey=""
  read minio_accessKey minio_secretKey < <(testHelpers::get_k8s_secret_data ${minio_secret_name})

  log::info "- Verifying MinIO bucket migration..."

  gatewayMigration::download_sample_file_from_minio "${public_bucket}" "${MINIO_GATEWAY_MIGRATION_SAMPLE_FILE}" "${minio_host}" "${minio_accessKey}" "${minio_secretKey}"
  gatewayMigration::download_sample_file_from_minio "${public_bucket}" "${MINIO_GATEWAY_MIGRATION_SAMPLE_FILE_WITH_DIR}" "${minio_host}" "${minio_accessKey}" "${minio_secretKey}"

  gatewayMigration::download_sample_file_from_minio "${private_bucket}" "${MINIO_GATEWAY_MIGRATION_SAMPLE_FILE}" "${minio_host}" "${minio_accessKey}" "${minio_secretKey}"
  gatewayMigration::download_sample_file_from_minio "${private_bucket}" "${MINIO_GATEWAY_MIGRATION_SAMPLE_FILE_WITH_DIR}" "${minio_host}" "${minio_accessKey}" "${minio_secretKey}"

  log::success "- Verified MinIO bucket migration."
}

# Arguments:
#   $1 - Release name
#   $2 - MiniIO host
#   $3 - The MiniIO k8s secret name
#   $4 - Path to charts directory
gatewayMigration::run() {
  local -r release_name="${1}"
  local -r minio_host="${2}"
  local -r minio_secret_name="${3}"
  local -r charts_path="${4}"

  # configure provider
  gatewayMigration::before_migration "${minio_host}" "${minio_secret_name}"
  gateway::before_test "${minio_secret_name}"

  # switch to MinIO gateway mode
  gateway::switch "${release_name}" "${charts_path}"

  # test migration
  gatewayMigration::after_migration "${minio_host}" "${MINIO_GATEWAY_SECRET_NAME}"
}
