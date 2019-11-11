#
# Library of useful utilities for CI purposes.
#

readonly LIB_DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )

# Prints first argument as header. Additionally prints current date.
shout() {
    echo -e "
#################################################################################################
# $(date)
# $1
#################################################################################################
"
}

# Workaround for release jobs which are using the build.sh script from test-infra
function ensure_docker() {
    set +e
    docker info &> /dev/null
    if [[ $? -eq 1 ]]; then
        start_docker
    fi
    set -e
}

# Copied from https://github.com/kyma-project/test-infra/blob/edde7967b87032ee494f0c00cfc60ae495d74ef4/prow/scripts/library.sh#L6
function start_docker() {
    shout "- Docker in Docker enabled, initializing..."
    printf '=%.0s' {1..80}; echo
    # If we have opted in to docker in docker, start the docker daemon,
    service docker start
    # the service can be started but the docker socket not ready, wait for ready
    local WAIT_N=0
    local MAX_WAIT=20
    while true; do
        # docker ps -q should only work if the daemon is ready
        docker ps -q > /dev/null 2>&1 && break
        if [[ ${WAIT_N} -lt ${MAX_WAIT} ]]; then
            WAIT_N=$((WAIT_N+1))
            echo "Waiting for docker to be ready, sleeping for ${WAIT_N} seconds."
            sleep ${WAIT_N}
        else
            echo "Reached maximum attempts, not waiting any longer..."
            exit 1
        fi
    done
    printf '=%.0s' {1..80}; echo

    echo "Done setting up docker in docker."
}

# Installs kind dependency locally.
# Required envs:
#  - KIND_VERSION
#  - INSTALL_DIR
#
# usage: env INSTALL_DIR=/tmp KIND_VERSION=v0.4.0 install::local::kind
install::kind() {
    mkdir -p "${INSTALL_DIR}/bin"
    export PATH="${INSTALL_DIR}/bin:${PATH}"

    pushd "${INSTALL_DIR}"

    os=$(host::os)
    arch=$(host::arch)

    shout "- Install kind ${KIND_VERSION} locally to a tempdir..."

    curl -sSLo kind "https://github.com/kubernetes-sigs/kind/releases/download/${KIND_VERSION}/kind-${os}-${arch}"
    chmod +x kind
    mv kind "${INSTALL_DIR}/bin"

    popd
}

host::os() {
  local host_os
  case "$(uname -s)" in
    Darwin)
      host_os=darwin
      ;;
    Linux)
      host_os=linux
      ;;
    *)
      kube::log::error "Unsupported host OS.  Must be Linux or Mac OS X."
      exit 1
      ;;
  esac
  echo "${host_os}"
}

host::arch() {
  local host_arch
  case "$(uname -m)" in
    x86_64*)
      host_arch=amd64
      ;;
    i?86_64*)
      host_arch=amd64
      ;;
    amd64*)
      host_arch=amd64
      ;;
    aarch64*)
      host_arch=arm64
      ;;
    arm64*)
      host_arch=arm64
      ;;
    arm*)
      host_arch=arm
      ;;
    ppc64le*)
      host_arch=ppc64le
      ;;
    *)
      kube::log::error "Unsupported host arch. Must be x86_64, arm, arm64, or ppc64le."
      exit 1
      ;;
  esac
  echo "${host_arch}"
}

#
# 'kind'(kubernetes-in-docker) functions
#
readonly KIND_CLUSTER_NAME="kind-ci"

kind::create_cluster() {
    shout "- Create k8s cluster..."
    kind create cluster --name=${KIND_CLUSTER_NAME} --image=kindest/node:${KUBERNETES_VERSION} --wait=5m
    export KUBECONFIG="$(kind get kubeconfig-path --name=${KIND_CLUSTER_NAME})"
}

kind::delete_cluster() {
    kind delete cluster --name=${KIND_CLUSTER_NAME}
}

# Arguments:
#   $1 - image name to copy into cluster nodes
kind::load_image() {
    kind load docker-image $1 --name=${KIND_CLUSTER_NAME}
}
