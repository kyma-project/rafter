#!/usr/bin/env bash

readonly STABLE_KUBERNETES_VERSION=v1.16.2
readonly STABLE_KIND_VERSION=v0.5.1
readonly STABLE_HELM_VERSION=v2.16.0
readonly CT_VERSION=v2.3.3
readonly CLUSTER_NAME=ci-test-cluster
# docker images to load into kind
readonly ARTIFACTS_DIR="${ARTIFACTS:-"${TMP_DIR}/artifacts"}"
readonly UPLOADER_IMG_NAME="${1}"
readonly MANAGER_IMG_NAME="${2}"
readonly FRONT_MATTER_IMG_NAME="${3}"
readonly ASYNCAPI_IMG_NAME="${4}"
readonly TMP_DIR="$(mktemp -d)"
readonly TMP_BIN_DIR="${TMP_DIR}/bin"
mkdir -p "${TMP_BIN_DIR}"
export PATH="${TMP_BIN_DIR}:${PATH}"

# external dependencies
readonly LIB_DIR="$(cd "${GOPATH}/src/github.com/kyma-project/test-infra/prow/scripts/lib" && pwd)"
CURRENT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)

source "${LIB_DIR}/kind.sh" || {
    echo 'Cannot load kind utilities.'
    exit 1
}
source "${LIB_DIR}/log.sh" || {
    echo 'Cannot load log utilities.'
    exit 1
}

source "${LIB_DIR}/host.sh" || {
    echo 'Cannot load host utilities.'
    exit 1
}

source "${LIB_DIR}/docker.sh" || {
    echo 'Cannot load docker utilities.'
    exit 1
}

source "${LIB_DIR}/docker.sh" || {
    echo 'Cannot load docker utilities.'
    exit 1
}


# kind cluster configuration
readonly CLUSTER_CONFIG=${CURRENT_DIR}/config/kind/cluster-config.yaml

readonly INSTALL_TIMEOUT=180
# minio access key that will be used during rafter installation
export APP_TEST_MINIO_ACCESSKEY=4j4gEuRH96ZFjptUFeFm
# minio secret key that will be used during the rafter installation
export APP_TEST_MINIO_SECRETKEY=UJnce86xA7hK01WblDdbmXg4gwjKwpFypdLJCvJ3

# the addres of the ingress that exposes upload and minio endpoints
export INGRESS_ADDRESS=http://localhost:30080
# URL of the uploader that will be used to upload test data in tests,
# it must be visible from outside of the cluster 
export APP_TEST_UPLOAD_SERVICE_URL=${INGRESS_ADDRESS}/v1/upload
export APP_TEST_MINIO_USE_SSL="false"
export APP_TEST_MINIO_ENDPOINT=localhost:30080

testHelper::install_rafter() {
    log::info '- Installing rafter...'
    helm install --name rafter \
    rafter-charts/rafter \
    --set rafter-controller-manager.minio.accessKey=${APP_TEST_MINIO_ACCESSKEY} \
    --set rafter-controller-manager.minio.secretKey=${APP_TEST_MINIO_SECRETKEY} \
    --set rafter-controller-manager.envs.store.externalEndpoint.value=${INGRESS_ADDRESS} \
    --wait \
    --timeout ${INSTALL_TIMEOUT}
}

testHelper::install_tiller() {
    log::info '- Installing Tiller...'
    kubectl --namespace kube-system create sa tiller
    kubectl create clusterrolebinding tiller-cluster-rule \
    --clusterrole=cluster-admin \
    --serviceaccount=kube-system:tiller \
    
    helm init \
    --service-account tiller \
    --upgrade --wait  \
    --history-max 200
}

# ingress http port
readonly NODE_PORT_HTTP=30080
# ingress https port
readonly NODE_PORT_HTTPS=30443

testHelper::install_ingress() {
    log::info '- Installing ingress...'
    helm install --name my-ingress stable/nginx-ingress \
    --set controller.service.type=NodePort \
    --set controller.service.nodePorts.http=${NODE_PORT_HTTP} \
    --set controller.service.nodePorts.https=${NODE_PORT_HTTPS} \
    --wait

    log::info '- Applying ingress...'
    kubectl apply -f ${CURRENT_DIR}/config/kind/ingress.yaml
}

testHelper::add_repos_and_update() {
    log::info '- Adding helm repositories and updating helm...'
    helm repo add rafter-charts https://rafter-charts.storage.googleapis.com
    helm repo update
}

testHelper::cleanup() {
    log::info "- Cleaning up cluster ${CLUSTER_NAME}..."
    kind::delete_cluster "${CLUSTER_NAME}" 2>&1
}

testHelper::start_integration_tests() {
    # required by integration suite
    export APP_KUBECONFIG_PATH="$(kind get kubeconfig-path --name=${CLUSTER_NAME})"
    log::info "Starting integration tests..."
    go test ${CURRENT_DIR}/../../tests/asset-store/main_test.go -count 1
}

testHelper::load_images() {
    log::info "- Loading image ${UPLOADER_IMG_NAME}..."
    kind::load_image "${CLUSTER_NAME}" "${UPLOADER_IMG_NAME}"
    
    log::info "- Loading image ${MANAGER_IMG_NAME}..."
    kind::load_image "${CLUSTER_NAME}" "${MANAGER_IMG_NAME}"
    
    log::info "- Loading image ${FRONT_MATTER_IMG_NAME}..."
    kind::load_image "${CLUSTER_NAME}" "${FRONT_MATTER_IMG_NAME}"
    
    log::info "- Loading image ${ASYNCAPI_IMG_NAME}..."
    kind::load_image "${CLUSTER_NAME}" "${ASYNCAPI_IMG_NAME}"
}

testHelper::check_helm_version() {
    readonly HELM_VERSION=$(helm version --short)
    if [[ $(helm version 2>/dev/null| sed 's/.*v\([0-9][0-9]*\)\..*$/\1/g') > 2 ]];
    then
        log::error "Invalid helm version ${HELM_VERSION}, required version is ${STABLE_HELM_VERSION}!"
        exit 1
    fi
}

testHelper::check_kind_version(){
    readonly KIND_VERSION=$(kind version)
    if [ "$KIND_VERSION" != "$STABLE_KIND_VERSION" ]; 
    then
        log::error "Invalid kind version ${KIND_VERSION}, required version is ${STABLE_KIND_VERSION}!"
        exit 1
    fi

}

infraHelper::install_helm_tiller(){
    log::info "Installing Helm and Tiller..."
    curl -LO https://get.helm.sh/helm-"${STABLE_HELM_VERSION}"-linux-amd64.tar.gz
    tar -xzvf helm-"${STABLE_HELM_VERSION}"-linux-amd64.tar.gz
    mv ./linux-amd64/{helm,tiller} /usr/local/bin
}