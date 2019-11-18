#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail
set -e

readonly UPLOADER_IMG_NAME="${1}"
readonly MANAGER_IMG_NAME="${2}"
readonly FRONT_MATTER_IMG_NAME="${3}"
readonly ASYNCAPI_IMG_NAME="${4}"

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
source "${CURRENT_DIR}/var.sh" || {
    echo 'Cannot load variables.'
    exit 1
}

readonly CLUSTER_CONFIG=${CURRENT_DIR}/config/kind/cluster-config.yaml
readonly KUBECONFIG="$(kind get kubeconfig-path --name=${CLUSTER_NAME})"

readonly HELM_CHARTS_RELATIVE_PATH=../../charts
readonly RAFTER_CHART_DIR=${CURRENT_DIR}/${HELM_CHARTS_RELATIVE_PATH}/rafter
readonly UPLOAD_SERVICE_CHART_DIR=${CURRENT_DIR}/${HELM_CHARTS_RELATIVE_PATH}/rafter-upload-service
INSTALL_TIMEOUT=180

readonly APP_TEST_MINIO_ACCESSKEY=AKIAIOSFODNN7EXAMPLE 
readonly APP_TEST_MINIO_SECRETKEY=wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY
readonly APP_KUBECONFIG_PATH=$KUBECONFIG
readonly APP_TEST_UPLOAD_SERVICE_URL=http://localhost:30080/rafter-upload-service

helm::installRafter() {
    log::info '- Installing rafter'
    helm repo add rafter-charts https://rafter-charts.storage.googleapis.com
    helm install --name rafter rafter-charts/rafter --wait --timeout ${INSTALL_TIMEOUT}
#    helm install ${RAFTER_CHART_DIR} --wait --timeout ${INSTALL_TIMEOUT}
}

kubectl::installTiller() {
    log::info '- Installing Tiller...'
    kubectl --namespace kube-system create sa tiller
    kubectl create clusterrolebinding tiller-cluster-rule --clusterrole=cluster-admin --serviceaccount=kube-system:tiller
    helm init --service-account tiller --upgrade --wait --history-max 200
}

helm::installIngress() {
    log::info '- Installing ingress...'
    helm install -n my-ingress stable/nginx-ingress \
        --set controller.service.type=NodePort \
        --set controller.service.nodePorts.http=30080 \
        --set controller.service.nodePorts.https=30443 \
        --wait
}

kubectl::applyIngress() {
    log::info '- Applying ingress...'
    kubectl apply -f ./config/ingress.yaml
}

installation::cleanup() {
    kind::delete_cluster "${CLUSTER_NAME}" 2>&1
}

kind::loadImages() {
    log::info "- Loading image ${UPLOADER_IMG_NAME}..."
    kind::load_image "${CLUSTER_NAME}" "${UPLOADER_IMG_NAME}"

    log::info "- Loading image ${MANAGER_IMG_NAME}..."
    kind::load_image "${CLUSTER_NAME}" "${MANAGER_IMG_NAME}"

    log::info "- Loading image ${FRONT_MATTER_IMG_NAME}..."
    kind::load_image "${CLUSTER_NAME}" "${FRONT_MATTER_IMG_NAME}"

    log::info "- Loading image ${ASYNCAPI_IMG_NAME}..."
    kind::load_image "${CLUSTER_NAME}" "${ASYNCAPI_IMG_NAME}"
}

main() {
    # trap installation::cleanup EXIT

     kind::create_cluster \
         "${CLUSTER_NAME}" \
         "${STABLE_KUBERNETES_VERSION}" \
         "${CLUSTER_CONFIG}" 2>&1

     kubectl::installTiller

     helm::installIngress

     kind::loadImages

     helm::installRafter

     kubectl::applyIngress

    go test ${CURRENT_DIR}/../../tests/asset-store/main_test.go
}

main
