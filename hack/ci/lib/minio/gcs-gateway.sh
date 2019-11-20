#!/usr/bin/env bash

gateway::beforeTest() {
    shout "Validating Google Cloud Storage Gateway environment"; date

    for var in GOOGLE_APPLICATION_CREDENTIALS CLOUDSDK_CORE_PROJECT TEST_INFRA_CLUSTER_INTEGRATION_SCRIPTS; do
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

gateway::installOverrides() {
    shout "Installing Google Cloud Storage Minio Gateway overrides"; date

    local -r ASSET_STORE_RESOURCE_NAME="gcs-minio-overrides"

    kubectl create -n kyma-installer secret generic "${ASSET_STORE_RESOURCE_NAME}" --from-file=minio.gcsgateway.gcsKeyJson="${GOOGLE_APPLICATION_CREDENTIALS}"
    kubectl label -n kyma-installer secret "${ASSET_STORE_RESOURCE_NAME}" "installer=overrides" "component=assetstore" "kyma-project.io/installation="

    "${TEST_INFRA_CLUSTER_INTEGRATION_SCRIPTS}/create-config-map.sh" --name "${ASSET_STORE_RESOURCE_NAME}" \
        --data "minio.persistence.enabled=false" \
        --data "minio.gcsgateway.enabled=true" \
        --data "minio.gcsgateway.projectId=${CLOUDSDK_CORE_PROJECT}" \
        --data "minio.DeploymentUpdate.type=RollingUpdate" \
        --data "minio.DeploymentUpdate.maxSurge=0" \
        --data "minio.DeploymentUpdate.maxUnavailable=50%" \
        --label "component=assetstore"
    
    echo "Overrides installed"; date
}

gateway::afterTest() {
    shout "Delete Google Cloud Storage Buckets"; date

    "${TEST_INFRA_CLUSTER_INTEGRATION_SCRIPTS}/delete-buckets.sh"

    echo "Buckets deleted"; date
}
