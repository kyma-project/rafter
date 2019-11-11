#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

readonly TMP_DIR=$(mktemp -d)

readonly CURRENT_DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
readonly LOCAL_REPO_ROOT_DIR=$( cd ${CURRENT_DIR}/../../ && pwd )
readonly CONTAINER_REPO_ROOT_DIR="/workdir"

source "${CURRENT_DIR}/lib/utilities.sh" || { echo 'Cannot load CI utilities.'; exit 1; }
source "${CURRENT_DIR}/lib/deps_ver.sh" || { echo 'Cannot load dependencies versions.'; exit 1; }

cleanup() {
    shout '- Removing ct container...'
    docker kill ct > /dev/null 2>&1
    kind::delete_cluster || true

    rm -rf "${TMP_DIR}" > /dev/null 2>&1 || true
    shout 'Cleanup Done!'
}

run_ct_container() {
    shout '- Running ct container...'
    docker run --rm --interactive --detach --network host --name ct \
        --volume "$LOCAL_REPO_ROOT_DIR":"$CONTAINER_REPO_ROOT_DIR" \
        --workdir "$CONTAINER_REPO_ROOT_DIR" \
        "quay.io/helmpack/chart-testing:$CT_VERSION" \
        cat
}

docker_ct_exec() {
    docker exec --interactive ct "$@"
}


setup_kubectl_in_ct_container() {
    docker_ct_exec mkdir -p /root/.kube

    shout '- Copying KUBECONFIG to container...'
    docker cp "$KUBECONFIG" ct:/root/.kube/config

    shout '- Checking connection to cluster...'
    docker_ct_exec kubectl cluster-info
}

install::tiller() {
    shout '- Installing Tiller...'
    docker_ct_exec kubectl --namespace kube-system create sa tiller
    docker_ct_exec kubectl create clusterrolebinding tiller-cluster-rule --clusterrole=cluster-admin --serviceaccount=kube-system:tiller
    docker_ct_exec helm init --service-account tiller --upgrade --wait --history-max 200
}

main() {
    if [[ "${RUN_ON_PROW-no}" = "true" ]]; then
        # This is a workaround for our CI. More info you can find in this issue:
        # https://github.com/kyma-project/test-infra/issues/1499
        ensure_docker
    fi

    run_ct_container
    trap cleanup EXIT

    export INSTALL_DIR=${TMP_DIR} KIND_VERSION=${STABLE_KIND_VERSION} HELM_VERSION=${STABLE_HELM_VERSION}
    install::kind

    export KUBERNETES_VERSION=${STABLE_KUBERNETES_VERSION}
    kind::create_cluster
    setup_kubectl_in_ct_container
    install::tiller

}

main
