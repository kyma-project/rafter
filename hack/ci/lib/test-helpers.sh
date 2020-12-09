#!/usr/bin/env bash

testHelpers::load_test_infra_utilities() {
  local -r lib_dir="$(cd "${GOPATH}/src/github.com/kyma-project/test-infra/prow/scripts/lib" && pwd)"

  source "${lib_dir}/docker.sh" || {
    echo 'Cannot load docker utilities.'
    exit 1
  }
  source "${lib_dir}/kind.sh" || {
    echo 'Cannot load kind utilities.'
    exit 1
  }
  source "${lib_dir}/log.sh" || {
    echo 'Cannot load log utilities.'
    exit 1
  }
  source "${lib_dir}/host.sh" || {
    echo 'Cannot load host utilities.'
    exit 1
  }
  source "${lib_dir}/kubernetes.sh" || {
    echo 'Cannot load kubernetes utilities.'
    exit 1
  }
  source "${lib_dir}/docker.sh" || {
    echo 'Cannot load docker utilities.'
    exit 1
  }
  source "${lib_dir}/junit.sh" || {	
    echo 'Cannot load JUnit utilities.'	
    exit 1	
  }
}
testHelpers::load_test_infra_utilities

testHelpers::install_go_junit_report() {
  log::info '- Installing go-junit-report...'

  if ! [ -x "$(command -v go-junit-report)" ]; then
      GO111MODULE=off go get -u github.com/jstemmer/go-junit-report
      log::success "- go-junit-reports installed."
      return 0
  fi

  log::info "- go-junit-reports already installed!"
}

# Arguments:
#   $1 - Helm version
#   $2 - Host OS
#   $3 - Destination directory
testHelpers::download_helm_tiller() {
  local -r helm_version="${1}"
  local -r host_os="${2}"
  local -r destination_dir="${3}"

  log::info "Downloading Helm and Tiller in version ${helm_version}..."

  curl -LO "https://get.helm.sh/helm-${helm_version}-${host_os}-amd64.tar.gz" --fail \
      && tar -xzvf "helm-${helm_version}-${host_os}-amd64.tar.gz" \
      && mv "./${host_os}-amd64/helm" "${destination_dir}/helm" \
      && mv "./${host_os}-amd64/tiller" "${destination_dir}/tiller" \
      && rm -rf "helm-${helm_version}-${host_os}-amd64.tar.gz" \
      && rm -rf "${host_os}-amd64"

  log::success "- Helm & Tiller downloaded in ${helm_version} version."
}

# Arguments:
#   $1 - Kind version
#   $2 - Host OS
#   $3 - Destination directory
testHelpers::download_kind() {
  local -r kind_version="${1}"
  local -r host_os="${2}"
  local -r destination_dir="${3}"

  log::info "Downloading kind in version ${kind_version}..."
  curl -LO "https://github.com/kubernetes-sigs/kind/releases/download/${kind_version}/kind-${host_os}-amd64" --fail \
      && chmod +x "kind-${host_os}-amd64" \
      && mv "kind-${host_os}-amd64" "${destination_dir}/kind" \
      && rm -rf "kind-${host_os}-amd64"

  log::success "- Kind downloaded."
}

testHelpers::install_tiller() {
  log::info '- Installing Tiller...'

  kubectl --namespace kube-system create sa tiller
  kubectl create clusterrolebinding tiller-cluster-rule \
      --clusterrole=cluster-admin \
      --serviceaccount=kube-system:tiller \

  helm init \
      --service-account tiller \
      --upgrade --wait  \
      --history-max 200

  log::success "- Tiller installed."
}

# Arguments:
#   $1 - Absolute path of repo
#   $2 - Temporary folder path for local charts
testHelpers::prepare_local_helm_charts() {
  local -r root_repo_path="${1}"
  local -r temp_rafter_dir="${2}"
  local -r charts_path="${root_repo_path}/charts"

  log::info "- Preparing local charts..."

  cp -r "${charts_path}/${RAFTER_CHART}/." "${temp_rafter_dir}/"

  local -r temp_rafter_charts_dir="${temp_rafter_dir}/charts"
  mkdir -p "${temp_rafter_charts_dir}"

  for img in RAFTER_CONTROLLER_MANAGER_CHART RAFTER_UPLOAD_SERVICE_CHART RAFTER_FRONT_MATTER_SERVICE_CHART RAFTER_ASYNCAPI_SERVICE_CHART; do
    cp -r "${charts_path}/${!img}" "${temp_rafter_charts_dir}"
  done

  helm dependency update "${temp_rafter_charts_dir}/${RAFTER_CONTROLLER_MANAGER_CHART}" # to fetch minio

  log::success "- Local charts prepared."
}

# Arguments:
#   $1 - The absolute path of custom ingress yaml definition
testHelpers::install_ingress() {
  local -r ingress_version="${1}"
  local -r ingress_yaml="${2}"
  local -r node_port_http=30080
  local -r node_port_https=30443

  log::info '- Installing ingress...'

  helm repo add stable https://kubernetes-charts.storage.googleapis.com/

  helm fetch stable/nginx-ingress --version ${ingress_version}

  helm install --name my-ingress stable/nginx-ingress --version ${ingress_version} \
      --set controller.service.type=NodePort \
      --set controller.service.nodePorts.http=${node_port_http} \
      --set controller.service.nodePorts.https=${node_port_https} \
      --wait

  log::info '- Applying custom ingress...'
  kubectl apply -f "${ingress_yaml}"

  log::success "- Ingress installed."
}

# Arguments:
#   $1 - Name of the cluser to load images to
testHelpers::load_rafter_images() {
  local -r cluster_name="${1}"

  log::info "- Loading Rafter images..."

  for img in RAFTER_CONTROLLER_MANAGER_CHART RAFTER_UPLOAD_SERVICE_CHART RAFTER_FRONT_MATTER_SERVICE_CHART RAFTER_ASYNCAPI_SERVICE_CHART; do
    log::info "- Loading image ${!img}..."
    kind::load_image "${cluster_name}" "${!img}"
  done

  log::success "- Rafter images loaded."
}

# Arguments:
#   $1 - Release name
#   $2 - The MiniIO k8s secret name
#   $3 - The addres of the ingress that exposes upload and minio endpoints
#   $4 - Path to charts directory
testHelpers::install_rafter() {
  local -r release_name="${1}"
  local -r minio_secret_name="${2}"
  local -r ingress_address="${3}"
  local -r charts_path="${4}"

  local -r tag="latest"
  local -r pull_policy="Never"
  local -r timeout=180

  log::info "- Installing Rafter in ${release_name} release from local charts..."

  helm install --name "${release_name}" "${charts_path}" \
    --set rafter-controller-manager.minio.existingSecret="${minio_secret_name}" \
    --set rafter-controller-manager.envs.store.externalEndpoint.value="${ingress_address}" \
    --set rafter-controller-manager.image.pullPolicy="${pull_policy}" \
    --set rafter-upload-service.image.pullPolicy="${pull_policy}" \
    --set rafter-front-matter-service.image.pullPolicy="${pull_policy}" \
    --set rafter-asyncapi-service.image.pullPolicy="${pull_policy}" \
    --set rafter-controller-manager.image.tag="${tag}" \
    --set rafter-upload-service.image.tag="${tag}" \
    --set rafter-front-matter-service.image.tag="${tag}" \
    --set rafter-asyncapi-service.image.tag="${tag}" \
    --set rafter-controller-manager.image.repository="${RAFTER_CONTROLLER_MANAGER_CHART}" \
    --set rafter-upload-service.image.repository="${RAFTER_UPLOAD_SERVICE_CHART}" \
    --set rafter-front-matter-service.image.repository="${RAFTER_FRONT_MATTER_SERVICE_CHART}" \
    --set rafter-asyncapi-service.image.repository="${RAFTER_ASYNCAPI_SERVICE_CHART}" \
    --wait \
    --timeout ${timeout}

  log::success "- Rafter installed."
}

# Arguments:
#   $1 - Name of new secret
#   $2 - The Minio access key
#   $3 - The Minio secret key
#   $4 - The path to GCS secret key - optional
testHelpers::create_k8s_secret() {
  local -r secret_name="${1}"
  local -r minio_accessKey="${2}"
  local -r minio_secretKey="${3}"

  log::info "- Creating ${secret_name} k8s secret..."

  if [ -n "${4-}" ]; then
    kubectl create secret generic "${secret_name}" \
      --from-literal=accesskey="${minio_accessKey}" \
      --from-literal=secretkey="${minio_secretKey}" \
      --from-file=gcs_key.json="${4}"
  else
    kubectl create secret generic "${secret_name}" \
      --from-literal=accesskey="${minio_accessKey}" \
      --from-literal=secretkey="${minio_secretKey}"
  fi

  log::success "- Secret created."
}

# Arguments:
#   $1 - The MiniIO k8s secret name.
# Outputs:
#   $1 - The MiniIO accessKey
#   $1 - The MiniIO secretKey
testHelpers::get_k8s_secret_data() {
  local -r secret_name="${1}"

  local -r host_os="$(host::os)"
  local minio_accessKey=""
  local minio_secretKey=""

  if [[ "${host_os}" == "darwin"* ]]; then
    # Mac OSX
    minio_accessKey=$(kubectl -n default get secret ${secret_name} -o jsonpath="{.data.accesskey}" | base64 -D | xargs -n1 echo)
    minio_secretKey=$(kubectl -n default get secret ${secret_name} -o jsonpath="{.data.secretkey}" | base64 -D | xargs -n1 echo)
  else 
    # Linux
    minio_accessKey=$(kubectl -n default get secret ${secret_name} -o jsonpath="{.data.accesskey}" | base64 -d | xargs -n1 echo)
    minio_secretKey=$(kubectl -n default get secret ${secret_name} -o jsonpath="{.data.secretkey}" | base64 -d | xargs -n1 echo)
  fi

  echo "${minio_accessKey}" "${minio_secretKey}"
}

# Arguments:
#   $1 - Absolute path of repo
#   $2 - Name of the kind cluster
#   $3 - The addres of the ingress that exposes minio endpoint
#   $4 - The addres of the ingress that exposes upload service endpoint
#   $5 - The MiniIO k8s secret name
#   $6 - Artifacts dir used to store JUnit reports. Optional
testHelpers::run_integration_tests() {
  local -r root_repo_path="${1}"
  local -r cluster_name="${2}"
  local -r minio_host="${3}"
  local -r upload_service_address="${4}"
  local -r minio_secret_name="${5}"
  local test_failed="false"

  local minio_accessKey=""
  local minio_secretKey=""
  read minio_accessKey minio_secretKey < <(testHelpers::get_k8s_secret_data ${minio_secret_name})

  log::info "- Starting integration tests..."

  export APP_KUBECONFIG_PATH="$(kind get kubeconfig-path --name=${cluster_name})"
  export APP_TEST_MINIO_USE_SSL="false"
  export APP_TEST_MINIO_ACCESSKEY="${minio_accessKey}"
  export APP_TEST_MINIO_SECRETKEY="${minio_secretKey}"
  export APP_TEST_MINIO_ENDPOINT="${minio_host}"
  export APP_TEST_UPLOAD_SERVICE_URL="${upload_service_address}"

  if [ -n "${6-}" ] ; then
    local -r artifacts_dir="${6}"
    local -r log_file=unit_test_data.log
    local -r suite_name="Rafter_Integration_Go_Test"

    go test ${root_repo_path}/tests/main_test.go -count 1 -v 2>&1 | tee "${log_file}" || test_failed="true"
    < "${log_file}" go-junit-report > "${artifacts_dir}/junit_${suite_name}_suite.xml"
    rm -rf "${log_file}"
  else
    go test ${root_repo_path}/tests/main_test.go -count 1
  fi
    
  if [[ ${test_failed} = "true" ]]; then
    log::error "- Integration test failed."
    return 1
  fi

  log::success "- Integration tests passed."
}
