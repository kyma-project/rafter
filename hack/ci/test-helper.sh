#!/usr/bin/env bash

# external dependencies
readonly LIB_DIR="$(cd "${GOPATH}/src/github.com/kyma-project/test-infra/prow/scripts/lib" && pwd)"
readonly CURRENT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)

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

source "${LIB_DIR}/kubernetes.sh" || {
    echo 'Cannot load kubernetes utilities.'
    exit 1
}

# Arguments:
#   $1 - minio access key that will be used during rafter installation
#   $2 - minio secret key that will be used during the rafter installation
#   $3 - the addres of the ingress that exposes upload and minio endpoints
testHelper::install_rafter() {
    local -r TAG=latest
    local -r PULL_POLICY=Never
    local -r INSTALL_TIMEOUT=180
    log::info '- Installing rafter...'
    helm install --name rafter \
    rafter-charts/rafter \
    --set rafter-controller-manager.minio.accessKey=${1} \
    --set rafter-controller-manager.minio.secretKey=${2} \
    --set rafter-controller-manager.envs.store.externalEndpoint.value=${3} \
    --set rafter-controller-manager.image.pullPolicy="${PULL_POLICY}" \
    --set rafter-upload-service.image.pullPolicy="${PULL_POLICY}" \
    --set rafter-asyncapi-service.image.pullPolicy="${PULL_POLICY}" \
    --set rafter-front-matter-service.image.pullPolicy="${PULL_POLICY}" \
    --set rafter-controller-manager.image.tag="${TAG}" \
    --set rafter-asyncapi-service.image.tag="${TAG}" \
    --set rafter-front-matter-service.image.tag="${TAG}" \
    --set rafter-upload-service.image.tag="${TAG}" \
    --set rafter-controller-manager.image.repository=rafter-controller-manager \
    --set rafter-asyncapi-service.image.repository=rafter-asyncapi-service \
    --set rafter-front-matter-service.image.repository=rafter-front-matter-service \
    --set rafter-upload-service.image.repository=rafter-upload-service \
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


testHelper::install_ingress() {
    # ingress http port
    local -r NODE_PORT_HTTP=30080
    # ingress https port
    local -r NODE_PORT_HTTPS=30443

    log::info '- Installing ingress...'
    helm install --name my-ingress stable/nginx-ingress \
    --set controller.service.type=NodePort \
    --set controller.service.nodePorts.http=${NODE_PORT_HTTP} \
    --set controller.service.nodePorts.https=${NODE_PORT_HTTPS} \
    --wait

    kubectl apply -f ${CURRENT_DIR}/config/kind/ingress.yaml
}

testHelper::add_repos_and_update() {
    log::info '- Adding helm repositories and updating helm...'
    helm repo add rafter-charts https://rafter-charts.storage.googleapis.com
    helm repo update
}

# Arguments:
#   $1 - name of the kind cluster
#   $2 - tmp directory with binaries used during test
testHelper::cleanup() {
    log::info "- Cleaning up cluster ${1}..."
    kind::delete_cluster "${1}"
    log::info "- Deleting directory with temporary binaries used in tests..."
    rm -rf "${2}"
}

# Arguments:
#   $1 - name of the kind cluster
#   $2 - minio access key that will be used during rafter installation
#   $3 - minio secret key that will be used during the rafter installation
#   $4 - the addres of the ingress that exposes upload and minio endpoints
testHelper::start_integration_tests() {
    # required by integration suite
    export APP_KUBECONFIG_PATH="$(kind get kubeconfig-path --name=${1})"
    export APP_TEST_MINIO_USE_SSL="false"
    # port same as ingress http port
    export APP_TEST_MINIO_ENDPOINT=localhost:30080
    # URL of the uploader that will be used to upload test data in tests,
    # it must be visible from outside of the cluster 
    export APP_TEST_MINIO_ACCESSKEY="${2}"
    export APP_TEST_MINIO_SECRETKEY="${3}"
    export APP_TEST_UPLOAD_SERVICE_URL="${4}/v1/upload"
    log::info "Starting integration tests..."
    go test ${CURRENT_DIR}/../../tests/asset-store/main_test.go -count 1
}

# Arguments:
#   $1 - name of the cluser to load images to
#   $2 - $5 - images to be loaded into cluster
testHelper::load_images() {
    log::info "- Loading image ${2}..."
    kind::load_image "${1}" "${2}"
    
    log::info "- Loading image ${3}..."
    kind::load_image "${1}" "${3}"
    
    log::info "- Loading image ${4}..."
    kind::load_image "${1}" "${4}"
    
    log::info "- Loading image ${5}..."
    kind::load_image "${1}" "${5}"
}

# Arguments:
#   $1 - Helm version
#   $2 - Host OS
#   $3 - Destination directory
infraHelper::install_helm_tiller(){
    log::info "Installing Helm and Tiller in version ${1}"
    curl -LO "https://get.helm.sh/helm-${1}-${2}-amd64.tar.gz" --fail \
        && tar -xzvf "helm-${1}-${2}-amd64.tar.gz" \
        && mv "./${2}-amd64/helm" "${3}/helm" \
        && mv "./${2}-amd64/tiller" "${3}/tiller" \
        && rm -rf "helm-${1}-${2}-amd64.tar.gz" \
        && rm -rf "${2}-amd64"
}

# Arguments:
#   $1 - Kind version
#   $2 - Host OS
#   $3 - Destination directory
infraHelper::install_kind(){
    log::info "Installing kind..."
    curl -LO "https://github.com/kubernetes-sigs/kind/releases/download/${1}/kind-${2}-amd64" --fail \
        && chmod +x "kind-${2}-amd64" \
        && mv "kind-${2}-amd64" "${3}/kind"
}