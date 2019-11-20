#!/usr/bin/env bash

__validateAzureGatewayEnvironment() {
    shout "Validating Azure Blob Gateway environment"; date

    for var in AZURE_RS_GROUP AZURE_REGION AZURE_SUBSCRIPTION_ID AZURE_SUBSCRIPTION_APP_ID AZURE_SUBSCRIPTION_SECRET AZURE_SUBSCRIPTION_TENANT AZURE_STORAGE_ACCOUNT_NAME; do
        if [ -z "${!var}" ] ; then
            echo "ERROR: $var is not set"
            local discoverUnsetVar=true
        fi
    done
    if [ "${discoverUnsetVar}" = true ] ; then
        exit 1
    fi

    echo "Environment validated"; date
}

__authenticateToAzure() {
    shout "Authenticating to Azure"; date

    az login --service-principal -u "${AZURE_SUBSCRIPTION_APP_ID}" -p "${AZURE_SUBSCRIPTION_SECRET}" --tenant "${AZURE_SUBSCRIPTION_TENANT}"
    az account set --subscription "${AZURE_SUBSCRIPTION_ID}"

    echo "Authenticated"; date
}

__createResourceGroup() {
    shout "Create Azure Resource Group ${AZURE_RS_GROUP}"; date

    if [[ $(az group exists --name "${AZURE_RS_GROUP}" -o json) == true ]]; then
        echo "Azure Resource Group ${AZURE_RS_GROUP} exists"; date
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
            echo -e "---\nAzure resource group ${AZURE_RS_GROUP} still not present after one minute wait.\n---"
            exit 1
        fi
    done

    echo "Resource Group created"; date
}

__createStorageAccount() {
    shout "Create ${AZURE_STORAGE_ACCOUNT_NAME} Storage Account"; date

    az storage account create \
        --name "${AZURE_STORAGE_ACCOUNT_NAME}" \
        --resource-group "${AZURE_RS_GROUP}" \
        --tags "created-at=$(date +%s)" "created-by=prow" "ttl=10800"

    echo "Storage Account created"; date
}

gateway::beforeTest() {
  __validateAzureGatewayEnvironment
  __authenticateToAzure
  __createResourceGroup
  __createStorageAccount
}

gateway::installOverrides() {
  shout "Installing Azure Minio Gateway overrides"; date

  local -r AZURE_ACCOUNT_KEY=$(az storage account keys list --account-name "${AZURE_STORAGE_ACCOUNT_NAME}" --resource-group "${AZURE_RS_GROUP}" --query "[?keyName=='key1'].value" --output tsv)

  local -r ASSET_STORE_RESOURCE_NAME="azure-minio-overrides"
  kubectl create -n kyma-installer secret generic "${ASSET_STORE_RESOURCE_NAME}" --from-literal "minio.secretKey=${AZURE_ACCOUNT_KEY}" --from-literal "minio.accessKey=${AZURE_STORAGE_ACCOUNT_NAME}"
  kubectl label -n kyma-installer secret "${ASSET_STORE_RESOURCE_NAME}" "installer=overrides" "component=assetstore" "kyma-project.io/installation="

  "${TEST_INFRA_CLUSTER_INTEGRATION_SCRIPTS}/create-config-map.sh" --name "${ASSET_STORE_RESOURCE_NAME}" \
    --data "minio.persistence.enabled=false" \
    --data "minio.azuregateway.enabled=true" \
    --data "minio.DeploymentUpdate.type=RollingUpdate" \
    --data "minio.DeploymentUpdate.maxSurge=0" \
    --data "minio.DeploymentUpdate.maxUnavailable=50%" \
    --label "component=assetstore"
    
  shout "Overrides installed"; date
}

gateway::afterTest() {
  shout "Delete ${AZURE_STORAGE_ACCOUNT_NAME} Storage Account"; date

  az storage account delete \
    --name "${AZURE_STORAGE_ACCOUNT_NAME}" \
    --resource-group "${AZURE_RS_GROUP}" \
    --yes

  echo "Storage Account deleted"; date
}
