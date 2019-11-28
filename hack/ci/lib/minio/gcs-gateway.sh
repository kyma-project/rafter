#!/usr/bin/env bash

readonly __MINIO_GATEWAY_SECRET_NAME__="gcs-minio-secret"

__validate_GCP_gateway_environment() {
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
        exit 1
    fi

    log::success "- Google Cloud Storage Gateway environment validated."
}

__authenticate_to_GCP() {
    log::info "- Authenticating to GCP..."

    gcloud config set project "${CLOUDSDK_CORE_PROJECT}"
    gcloud auth activate-service-account --key-file "${GOOGLE_APPLICATION_CREDENTIALS}"

    log::success "- Authenticated."
}

# Delete given bucket.
# Arguments:
#   $1 - The name of bucket
__delete_bucket() {
  local -r bucket_name="${1}"

  log::info "- Deleting ${bucket_name} bucket..."
  gsutil rm -r "gs://${bucket_name}"
  log::success "- ${bucket_name} bucket deleted."
}

__delete_GCP_buckets() {
  log::info "- Deleting Google Cloud Storage Buckets..."

  local -r cluster_buckets=$(kubectl get clusterbuckets -o jsonpath="{.items[*].status.remoteName}" | xargs -n1 echo)
  local -r buckets=$(kubectl get buckets --all-namespaces -o jsonpath="{.items[*].status.remoteName}" | xargs -n1 echo)

  local public_bucket=""
  local private_bucket=""
  read public_bucket private_bucket < <(gatewayHelpers::get_bucket_names)

  for clusterBucket in ${cluster_buckets}
  do
    __delete_bucket "${clusterBucket}"
  done

  for bucket in ${buckets}
  do
    __delete_bucket "${bucket}"
  done

  if [ -n "${public_bucket}" ]; then
    __delete_bucket "${public_bucket}"
  fi

  if [ -n "${private_bucket}" ]; then
    __delete_bucket "${private_bucket}"
  fi

  log::success "- Buckets deleted."
}

# Arguments:
#   $1 - Minio access key
#   $2 - Minio secret key
gateway::before_test() {
    local -r minio_secret_name="${1}"
    local minio_accessKey=""
    local minio_secretKey=""
    read minio_accessKey minio_secretKey < <(testHelpers::get_minio_k8s_secret ${minio_secret_name})

    junit::test_start "MinIO_Gateway_GCP_Validate_GCP_Gateway_Environment"
    __validate_GCP_gateway_environment 2>&1 | junit::test_output
    junit::test_pass

    junit::test_start "MinIO_Gateway_GCP_Authenticate_To_GCP"
    __authenticate_to_GCP 2>&1 | junit::test_output
    junit::test_pass
    
    junit::test_start "MinIO_Gateway_GCP_Create_MinIO_K8S_Secret"
    testHelpers::create_minio_k8s_secret "${__MINIO_GATEWAY_SECRET_NAME__}" "${minio_accessKey}" "${minio_secretKey}" "${GOOGLE_APPLICATION_CREDENTIALS}" 2>&1 | junit::test_output
    junit::test_pass
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
        --set rafter-controller-manager.minio.existingSecret="${__MINIO_GATEWAY_SECRET_NAME__}" \
        --set rafter-controller-manager.minio.gcsgateway.enabled="true" \
        --set rafter-controller-manager.minio.gcsgateway.projectId="${CLOUDSDK_CORE_PROJECT}" \
        --set rafter-controller-manager.minio.DeploymentUpdate.type="RollingUpdate" \
        --set rafter-controller-manager.minio.DeploymentUpdate.maxSurge="0" \
        --set rafter-controller-manager.minio.DeploymentUpdate.maxUnavailable="\"50%\"" \
        --set rafter-controller-manager.envs.store.accessKey.valueFrom.secretKeyRef.name="${__MINIO_GATEWAY_SECRET_NAME__}" \
        --set rafter-controller-manager.envs.store.secretKey.valueFrom.secretKeyRef.name="${__MINIO_GATEWAY_SECRET_NAME__}" \
        --set rafter-upload-service.minio.persistence.enabled="false" \
        --set rafter-upload-service.envs.upload.accessKey.valueFrom.secretKeyRef.name="${__MINIO_GATEWAY_SECRET_NAME__}" \
        --set rafter-upload-service.envs.upload.secretKey.valueFrom.secretKeyRef.name="${__MINIO_GATEWAY_SECRET_NAME__}" \
        --set rafter-controller-manager.image.pullPolicy="${pull_policy}" \
        --set rafter-upload-service.image.pullPolicy="${pull_policy}" \
        --set rafter-front-matter-service.image.pullPolicy="${pull_policy}" \
        --set rafter-asyncapi-service.image.pullPolicy="${pull_policy}" \
        --set rafter-controller-manager.image.tag="${tag}" \
        --set rafter-upload-service.image.tag="${tag}" \
        --set rafter-front-matter-service.image.tag="${tag}" \
        --set rafter-asyncapi-service.image.tag="${tag}" \
        --set rafter-controller-manager.image.repository="${__RAFTER_CONTROLLER_MANAGER__}" \
        --set rafter-upload-service.image.repository="${__RAFTER_UPLOAD_SERVICE__}" \
        --set rafter-front-matter-service.image.repository="${__RAFTER_FRONT_MATTER_SERVICE__}" \
        --set rafter-asyncapi-service.image.repository="${__RAFTER_ASYNCAPI_SERVICE__}" \
        --wait \
        --timeout ${timeout}
    
    log::success "- Rafter in ${release_name} release installed."
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
        --set rafter-controller-manager.minio.podAnnotations.persistence="\"false\"" \
        --set rafter-controller-manager.minio.existingSecret="${__MINIO_GATEWAY_SECRET_NAME__}" \
        --set rafter-controller-manager.minio.gcsgateway.enabled="true" \
        --set rafter-controller-manager.minio.gcsgateway.projectId="${CLOUDSDK_CORE_PROJECT}" \
        --set rafter-controller-manager.minio.DeploymentUpdate.type="RollingUpdate" \
        --set rafter-controller-manager.minio.DeploymentUpdate.maxSurge="0" \
        --set rafter-controller-manager.minio.DeploymentUpdate.maxUnavailable="\"50%\"" \
        --set rafter-controller-manager.envs.store.accessKey.valueFrom.secretKeyRef.name="${__MINIO_GATEWAY_SECRET_NAME__}" \
        --set rafter-controller-manager.envs.store.secretKey.valueFrom.secretKeyRef.name="${__MINIO_GATEWAY_SECRET_NAME__}" \
        --set rafter-upload-service.minio.persistence.enabled="false" \
        --set rafter-upload-service.minio.podAnnotations.persistence="\"false\"" \
        --set rafter-upload-service.envs.upload.accessKey.valueFrom.secretKeyRef.name="${__MINIO_GATEWAY_SECRET_NAME__}" \
        --set rafter-upload-service.envs.upload.secretKey.valueFrom.secretKeyRef.name="${__MINIO_GATEWAY_SECRET_NAME__}" \
        --set rafter-upload-service.migrator.post.minioSecretRefName="${__MINIO_GATEWAY_SECRET_NAME__}" \
        --wait \
        --timeout ${timeout}
    
    log::success "- Switched to Google Cloud Storage Minio Gateway mode."
}

gateway::after_test() {
    junit::test_start "MinIO_Gateway_GCP_Delete_GCP_Buckets"
    __delete_GCP_buckets 2>&1 | junit::test_output
    junit::test_pass
}
