#!/usr/bin/env bash

set -o errexit

echo "$(date +"%Y/%m/%d %T %Z"): Validating environment"
discoverUnsetVar=false

for var in AZURE_RS_GROUP AZURE_SUBSCRIPTION_ID AZURE_SUBSCRIPTION_APP_ID AZURE_SUBSCRIPTION_SECRET AZURE_SUBSCRIPTION_TENANT; do
    if [ -z "${!var}" ] ; then
        echo "ERROR: $var is not set"
        discoverUnsetVar=true
    fi
done
if [ "${discoverUnsetVar}" = true ] ; then
    exit 1
fi

echo "$(date +"%Y/%m/%d %T %Z"): Authenticating to Azure"
az login --service-principal -u "${AZURE_SUBSCRIPTION_APP_ID}" -p "${AZURE_SUBSCRIPTION_SECRET}" --tenant "${AZURE_SUBSCRIPTION_TENANT}"
az account set --subscription "${AZURE_SUBSCRIPTION_ID}"

echo "$(date +"%Y/%m/%d %T %Z"): Removing orphaned Storage Accounts from ${AZURE_RS_GROUP} Resource Group"
while read -r account; do
    if [[ -z "${account}" ]]; then
        # Nothing to delete
        continue
    fi

    echo "....Removing ${account}"
    az storage account delete \
        --name "${account}" \
        --resource-group "${AZURE_RS_GROUP}" \
        --yes
done <<< "$(az storage account list \
    --query "[?tags.\"created-by\"=='prow' && sum([to_number(tags.\"created-at\"),to_number(tags.\"ttl\")]) < to_number('$(date +%s)')].name" \
    --resource-group "${AZURE_RS_GROUP}" \
    --output tsv)"

echo "$(date +"%Y/%m/%d %T %Z"): Finished"