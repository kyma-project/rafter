#!/usr/bin/env bash

readonly MINIO_GATEWAY_SECRET_NAME="az-minio-secret"

azureGateway::authenticate_to_azure() {
  log::info "- Authenticating to Azure..."

  az::login "$AZURE_CREDENTIALS_FILE"
  az::set_subscription "$AZURE_SUBSCRIPTION_ID"

  log::success "- Authenticated."
}

azureGateway::create_resource_group() {
  log::info "- Creating Azure Resource Group ${AZURE_RS_GROUP}..."

  if [[ $(az group exists --name "${AZURE_RS_GROUP}" -o json) == true ]]; then
    log::warn "- Azure Resource Group ${AZURE_RS_GROUP} exists"
    return
  fi

  az group create \
    --name "${AZURE_RS_GROUP}" \
    --location "${AZURE_REGION}" \
    --tags "created-by=prow"

  # Wait until resource group will be visible in azure.
  counter=0
  until [[ $(az group exists --name "${AZURE_RS_GROUP}" -o json) == true ]]; do
    sleep 15
    counter=$(( counter + 1 ))
    if (( counter == 5 )); then
      log::error -e "---\nAzure resource group ${AZURE_RS_GROUP} still not present after one minute wait.\n---"
      exit 1
    fi
  done

  log::success "- Resource Group created."
}

azureGateway::create_storage_account_name() {
  log::info "- Creating Azure Storage Account Name..."

  local -r random_name_suffix=$(LC_ALL=C tr -dc 'a-z0-9' < /dev/urandom | head -c10)

  if [[ "$BUILD_TYPE" == "pr" ]]; then
    # In case of PR, operate on PR number
    AZURE_STORAGE_ACCOUNT_NAME=$(echo "mimpr${PULL_NUMBER}${random_name_suffix}" | tr "[:upper:]" "[:lower:]")
  elif [[ "$BUILD_TYPE" == "release" ]]; then
    # In case of release
    AZURE_STORAGE_ACCOUNT_NAME=$(echo "mimrel${random_name_suffix}" | tr "[:upper:]" "[:lower:]")
  else
    # Otherwise (master), operate on triggering commit id
    local -r commit_id=$(git rev-parse --short HEAD)
    AZURE_STORAGE_ACCOUNT_NAME=$(echo "mim${commit_id}${random_name_suffix}" | tr "[:upper:]" "[:lower:]")
  fi

  log::success "- ${AZURE_STORAGE_ACCOUNT_NAME} storage Account Name created."
}

azureGateway::create_storage_account() {
  log::info "- Creating ${AZURE_STORAGE_ACCOUNT_NAME} Storage Account..."

  az storage account create \
    --name "${AZURE_STORAGE_ACCOUNT_NAME}" \
    --resource-group "${AZURE_RS_GROUP}" \
    --tags "created-at=$(date +%s)" "created-by=prow" "ttl=10800"

  log::success "- Storage Account created."
}

azureGateway::delete_storage_account() {
  if [ -z "${AZURE_STORAGE_ACCOUNT_NAME}" ]; then
    return 0
  fi

  log::info "- Deleting ${AZURE_STORAGE_ACCOUNT_NAME} Storage Account..."

  az storage account delete \
    --name "${AZURE_STORAGE_ACCOUNT_NAME}" \
    --resource-group "${AZURE_RS_GROUP}" \
    --yes

  log::success "- Storage Account deleted."
}

gateway::validate_environment() {
  log::info "- Validating Azure Blob Gateway environment..."

  local discoverUnsetVar=false
  for var in AZURE_RS_GROUP AZURE_REGION AZURE_SUBSCRIPTION_ID AZURE_CREDENTIALS_FILE BUILD_TYPE; do
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

  log::success "- Azure Blob Gateway environment validated."
}

gateway::before_test() {
  azureGateway::authenticate_to_azure
  azureGateway::create_resource_group
  azureGateway::create_storage_account_name
  azureGateway::create_storage_account

  local -r azure_account_key=$(az storage account keys list --account-name "${AZURE_STORAGE_ACCOUNT_NAME}" --resource-group "${AZURE_RS_GROUP}" --query "[0].value" --output tsv)
  testHelpers::create_k8s_secret "${MINIO_GATEWAY_SECRET_NAME}" "${AZURE_STORAGE_ACCOUNT_NAME}" "${azure_account_key}"
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

  log::info "- Installing Rafter with Azure Minio Gateway mode in ${release_name} release..."

  helm install --name "${release_name}" "${charts_path}" \
    --set rafter-controller-manager.minio.persistence.enabled="false" \
    --set rafter-controller-manager.envs.store.externalEndpoint.value="${ingress_address}" \
    --set rafter-controller-manager.minio.existingSecret="${MINIO_GATEWAY_SECRET_NAME}" \
    --set rafter-controller-manager.minio.azuregateway.enabled="true" \
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

  log::info "- Switching to Azure Minio Gateway mode..."

  helm upgrade "${release_name}" "${charts_path}" \
    --reuse-values \
    --set rafter-controller-manager.minio.persistence.enabled="false" \
    --set rafter-controller-manager.minio.podAnnotations.persistence="off" \
    --set rafter-controller-manager.minio.existingSecret="${MINIO_GATEWAY_SECRET_NAME}" \
    --set rafter-controller-manager.minio.azuregateway.enabled="true" \
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
    
  log::success "- Switched to Azure Minio Gateway mode."
}

gateway::after_test() {
  azureGateway::delete_storage_account
}
