#!/usr/bin/env bash

readonly __MINIO_GATEWAY_SECRET_NAME__="az-minio-secret"
__AZURE_STORAGE_ACCOUNT_NAME__=""

__validate_Azure_gateway_environment() {
    log::info "- Validating Azure Blob Gateway environment..."

    local discoverUnsetVar=false
    for var in AZURE_RS_GROUP AZURE_REGION AZURE_SUBSCRIPTION_ID AZURE_SUBSCRIPTION_APP_ID AZURE_SUBSCRIPTION_SECRET AZURE_SUBSCRIPTION_TENANT BUILD_TYPE; do
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

    log::success "- Azure Blob Gateway environment validated."
}

__authenticate_to_Azure() {
    log::info "- Authenticating to Azure..."

    az login --service-principal -u "${AZURE_SUBSCRIPTION_APP_ID}" -p "${AZURE_SUBSCRIPTION_SECRET}" --tenant "${AZURE_SUBSCRIPTION_TENANT}"
    az account set --subscription "${AZURE_SUBSCRIPTION_ID}"

    log::success "- Authenticated."
}

__create_resource_group() {
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

__create_storage_account_name() {
    log::info "- Creating Azure Storage Account Name..."

    local -r random_name_suffix=$(LC_ALL=C tr -dc 'a-z0-9' < /dev/urandom | head -c10)

    if [[ "$BUILD_TYPE" == "pr" ]]; then
        # In case of PR, operate on PR number
        __AZURE_STORAGE_ACCOUNT_NAME__=$(echo "mimpr${PULL_NUMBER}${random_name_suffix}" | tr "[:upper:]" "[:lower:]")
    elif [[ "$BUILD_TYPE" == "release" ]]; then
        # In case of release
        __AZURE_STORAGE_ACCOUNT_NAME__=$(echo "mimrel${random_name_suffix}" | tr "[:upper:]" "[:lower:]")
    else
        # Otherwise (master), operate on triggering commit id
        __AZURE_STORAGE_ACCOUNT_NAME__=$(echo "mim${COMMIT_ID}${random_name_suffix}" | tr "[:upper:]" "[:lower:]")
    fi

    log::success "- ${__AZURE_STORAGE_ACCOUNT_NAME__} storage Account Name created."
}

__create_storage_account() {
    log::info "- Creating ${__AZURE_STORAGE_ACCOUNT_NAME__} Storage Account..."

    az storage account create \
        --name "${__AZURE_STORAGE_ACCOUNT_NAME__}" \
        --resource-group "${AZURE_RS_GROUP}" \
        --tags "created-at=$(date +%s)" "created-by=prow" "ttl=10800"

    log::success "- Storage Account created."
}

__delete_storage_account() {
  if [ -z "${__AZURE_STORAGE_ACCOUNT_NAME__}" ]; then
    return 0
  fi

  log::info "- Deleting ${__AZURE_STORAGE_ACCOUNT_NAME__} Storage Account..."

  az storage account delete \
    --name "${__AZURE_STORAGE_ACCOUNT_NAME__}" \
    --resource-group "${AZURE_RS_GROUP}" \
    --yes

  log::success "- Storage Account deleted."
}

gateway::before_test() {
  junit::test_start "MinIO_Gateway_Azure_Validate_Azure_Gateway_Environment"
  __validate_Azure_gateway_environment 2>&1 | junit::test_output
  junit::test_pass

  junit::test_start "MinIO_Gateway_Azure_Authenticate_To_Azure"
  __authenticate_to_Azure 2>&1 | junit::test_output
  junit::test_pass

  junit::test_start "MinIO_Gateway_Azure_Create_Resource_Group"  
  __create_resource_group 2>&1 | junit::test_output
  junit::test_pass

  __create_storage_account_name

  junit::test_start "MinIO_Gateway_Azure_Create_Storage_Account"
  __create_storage_account 2>&1 | junit::test_output
  junit::test_pass

  local -r azure_account_key=$(az storage account keys list --account-name "${__AZURE_STORAGE_ACCOUNT_NAME__}" --resource-group "${AZURE_RS_GROUP}" --query "[0].value" --output tsv)
  junit::test_start "MinIO_Gateway_Azure_Create_MinIO_K8S_Secret"
  testHelpers::create_minio_k8s_secret "${__MINIO_GATEWAY_SECRET_NAME__}" "${__AZURE_STORAGE_ACCOUNT_NAME__}" "${azure_account_key}" 2>&1 | junit::test_output
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

  log::info "- Installing Rafter with Azure Minio Gateway mode in ${release_name} release..."

  helm install --name "${release_name}" "${charts_path}" \
    --set rafter-controller-manager.minio.persistence.enabled="false" \
    --set rafter-controller-manager.envs.store.externalEndpoint.value="${ingress_address}" \
    --set rafter-controller-manager.minio.existingSecret="${__MINIO_GATEWAY_SECRET_NAME__}" \
    --set rafter-controller-manager.minio.azuregateway.enabled="true" \
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

  log::info "- Switching to Azure Minio Gateway mode..."

  helm upgrade "${release_name}" "${charts_path}" \
    --reuse-values \
    --set rafter-controller-manager.minio.persistence.enabled="false" \
    --set rafter-controller-manager.minio.podAnnotations.persistence="\"false\"" \
    --set rafter-controller-manager.minio.existingSecret="${__MINIO_GATEWAY_SECRET_NAME__}" \
    --set rafter-controller-manager.minio.azuregateway.enabled="true" \
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
    
  log::success "- Switched to Azure Minio Gateway mode."
}

gateway::after_test() {
  junit::test_start "MinIO_Gateway_Azure_Delete_Storage_Account"
  __delete_storage_account 2>&1 | junit::test_output
  junit::test_pass
}
