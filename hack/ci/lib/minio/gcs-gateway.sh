#!/usr/bin/env bash

readonly MINIO_GATEWAY_SECRET_NAME="gcs-minio-secret"

gcsGateway::authenticate_to_gcp() {
  log::info "- Authenticating to GCP..."

  gcloud config set project "${CLOUDSDK_CORE_PROJECT}"
  gcloud auth activate-service-account --key-file "${GOOGLE_APPLICATION_CREDENTIALS}"

  log::success "- Authenticated."
}

# Delete given bucket.
# Arguments:
#   $1 - The name of bucket
gcsGateway::delete_bucket() {
  local -r bucket_name="${1}"

  log::info "- Deleting ${bucket_name} bucket..."
  gsutil rm -r "gs://${bucket_name}"
  log::success "- ${bucket_name} bucket deleted."
}

gcsGateway::delete_gcp_buckets() {
  log::info "- Deleting Google Cloud Storage Buckets..."

  local -r cluster_buckets=$(kubectl get clusterbuckets -o jsonpath="{.items[*].status.remoteName}" | xargs -n1 echo)
  local -r buckets=$(kubectl get buckets --all-namespaces -o jsonpath="{.items[*].status.remoteName}" | xargs -n1 echo)

  local public_bucket=""
  local private_bucket=""
  read public_bucket private_bucket < <(gatewayHelpers::get_bucket_names)

  for clusterBucket in ${cluster_buckets}
  do
    gcsGateway::delete_bucket "${clusterBucket}"
  done

  for bucket in ${buckets}
  do
    gcsGateway::delete_bucket "${bucket}"
  done

  if [ -n "${public_bucket}" ]; then
    gcsGateway::delete_bucket "${public_bucket}"
  fi

  if [ -n "${private_bucket}" ]; then
    gcsGateway::delete_bucket "${private_bucket}"
  fi

  log::success "- Buckets deleted."
}

gateway::validate_environment() {
  log::info "- Validating Google Cloud Storage Gateway environment..."

  local discoverUnsetVar=false
  for var in GOOGLE_APPLICATION_CREDENTIALS CLOUDSDK_CORE_PROJECT; do
    if [ -n "${var-}" ] ; then
      continue
    else
      log::error "- ERROR: $var is not set"
      discoverUnsetVar=true
    fi
  done
  if [ "${discoverUnsetVar}" = true ] ; then
    return 1
  fi

  log::success "- Google Cloud Storage Gateway environment validated."
}

# Arguments:
#   $1 - Minio access key
#   $2 - Minio secret key
gateway::before_test() {
  local -r minio_secret_name="${1}"
  local minio_accessKey=""
  local minio_secretKey=""
  read minio_accessKey minio_secretKey < <(testHelpers::get_k8s_secret_data ${minio_secret_name})

  gcsGateway::authenticate_to_gcp
  testHelpers::create_k8s_secret "${MINIO_GATEWAY_SECRET_NAME}" "${minio_accessKey}" "${minio_secretKey}" "${GOOGLE_APPLICATION_CREDENTIALS}"
}

# Arguments:
#   $1 - Release name
#   $2 - The addres of the ingress that exposes upload and minio endpoints
#   $3 - Path to charts directory
gateway::install() {
  local -r release_name="${1}"
  local -r ingress_address="${2}"
  local -r charts_path="${3}"

  local -r tag="latest"
  local -r pull_policy="Never"
  local -r timeout=180

  log::info "- Installing Rafter with Google Cloud Storage Minio Gateway mode in ${release_name} release..."

  helm install --name "${release_name}" "${charts_path}" \
    --set rafter-controller-manager.minio.persistence.enabled="false" \
    --set rafter-controller-manager.envs.store.externalEndpoint.value="${ingress_address}" \
    --set rafter-controller-manager.minio.existingSecret="${MINIO_GATEWAY_SECRET_NAME}" \
    --set rafter-controller-manager.minio.gcsgateway.enabled="true" \
    --set rafter-controller-manager.minio.gcsgateway.projectId="${CLOUDSDK_CORE_PROJECT}" \
    --set rafter-controller-manager.minio.DeploymentUpdate.type="RollingUpdate" \
    --set rafter-controller-manager.minio.DeploymentUpdate.maxSurge="0" \
    --set rafter-controller-manager.minio.DeploymentUpdate.maxUnavailable="\"50%\"" \
    --set rafter-controller-manager.envs.store.accessKey.valueFrom.secretKeyRef.name="${MINIO_GATEWAY_SECRET_NAME}" \
    --set rafter-controller-manager.envs.store.secretKey.valueFrom.secretKeyRef.name="${MINIO_GATEWAY_SECRET_NAME}" \
    --set rafter-upload-service.minio.persistence.enabled="false" \
    --set rafter-upload-service.envs.upload.accessKey.valueFrom.secretKeyRef.name="${MINIO_GATEWAY_SECRET_NAME}" \
    --set rafter-upload-service.envs.upload.secretKey.valueFrom.secretKeyRef.name="${MINIO_GATEWAY_SECRET_NAME}" \
    --set rafter-controller-manager.image.pullPolicy="${pull_policy}" \
    --set rafter-upload-service.image.pullPolicy="${pull_policy}" \
    --set rafter-front-matter-service.image.pullPolicy="${pull_policy}" \
    --set rafter-asyncapi-service.image.pullPolicy="${pull_policy}" \
    --set rafter-controller-manager.image.tag="${tag}" \
    --set rafter-upload-service.image.tag="${tag}" \
    --set rafter-front-matter-service.image.tag="${tag}" \
    --set rafter-asyncapi-service.image.tag="${tag}" \
    --set rafter-controller-manager.image.repository="${RAFTER_CONTROLLER_MANAGER_CHART}" \
    --set rafter-upload-service.image.repository="${RAFTER_UPLOAD_SERVICE_CHART}" \
    --set rafter-front-matter-service.image.repository="${RAFTER_FRONT_MATTER_SERVICE_CHART}" \
    --set rafter-asyncapi-service.image.repository="${RAFTER_ASYNCAPI_SERVICE_CHART}" \
    --wait \
    --timeout ${timeout}
    
  log::success "- Rafter installed."
}

# Arguments:
#   $1 - Release name
#   $2 - Path to charts directory
gateway::switch() {
  local -r release_name="${1}"
  local -r charts_path="${2}"

  local -r timeout=180

  log::info "- Switching to Google Cloud Storage Minio Gateway mode..."

  helm upgrade "${release_name}" "${charts_path}" \
    --reuse-values \
    --set rafter-controller-manager.minio.persistence.enabled="false" \
    --set rafter-controller-manager.minio.podAnnotations.persistence="off" \
    --set rafter-controller-manager.minio.existingSecret="${MINIO_GATEWAY_SECRET_NAME}" \
    --set rafter-controller-manager.minio.gcsgateway.enabled="true" \
    --set rafter-controller-manager.minio.gcsgateway.projectId="${CLOUDSDK_CORE_PROJECT}" \
    --set rafter-controller-manager.minio.DeploymentUpdate.type="RollingUpdate" \
    --set rafter-controller-manager.minio.DeploymentUpdate.maxSurge="0" \
    --set rafter-controller-manager.minio.DeploymentUpdate.maxUnavailable="\"50%\"" \
    --set rafter-controller-manager.envs.store.accessKey.valueFrom.secretKeyRef.name="${MINIO_GATEWAY_SECRET_NAME}" \
    --set rafter-controller-manager.envs.store.secretKey.valueFrom.secretKeyRef.name="${MINIO_GATEWAY_SECRET_NAME}" \
    --set rafter-upload-service.minio.persistence.enabled="false" \
    --set rafter-upload-service.minio.podAnnotations.persistence="off" \
    --set rafter-upload-service.envs.upload.accessKey.valueFrom.secretKeyRef.name="${MINIO_GATEWAY_SECRET_NAME}" \
    --set rafter-upload-service.envs.upload.secretKey.valueFrom.secretKeyRef.name="${MINIO_GATEWAY_SECRET_NAME}" \
    --set rafter-upload-service.migrator.post.minioSecretRefName="${MINIO_GATEWAY_SECRET_NAME}" \
    --wait \
    --timeout ${timeout}
    
  log::success "- Switched to Google Cloud Storage Minio Gateway mode."
}

gateway::after_test() {
  gcsGateway::delete_gcp_buckets
}
