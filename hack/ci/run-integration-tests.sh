#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

readonly CURRENT_DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )

source "${CURRENT_DIR}/lib/utilities.sh" || { echo 'Cannot load CI utilities.'; exit 1; }
source "${CURRENT_DIR}/lib/deps_ver.sh" || { echo 'Cannot load dependencies versions.'; exit 1; }

cleanup() {
    shout '- Removing kind cluster...'
    kind::delete_cluster || true
    shout 'Cleanup Done!'
}


install::tiller() {
    shout '- Installing Tiller...'
    kubectl --namespace kube-system create sa tiller
    kubectl create clusterrolebinding tiller-cluster-rule --clusterrole=cluster-admin --serviceaccount=kube-system:tiller
    helm init --service-account tiller --upgrade --wait --history-max 200
}

main() {
    if [[ "${RUN_ON_PROW-no}" = "true" ]]; then
        # This is a workaround for our CI. More info you can find in this issue:
        # https://github.com/kyma-project/test-infra/issues/1499
        ensure_docker
    fi

#    run_container
    trap cleanup EXIT
#    export INSTALL_DIR=${TMP_DIR} KIND_VERSION=${STABLE_KIND_VERSION} HELM_VERSION=${STABLE_HELM_VERSION}
    export KUBERNETES_VERSION=${STABLE_KUBERNETES_VERSION}
    kind::create_cluster
#    setup_kubectl_in_container
    install::tiller
}

main
