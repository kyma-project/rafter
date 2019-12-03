#!/usr/bin/env bash

readonly STABLE_HELM_VERSION="v2.16.1"
readonly STABLE_KIND_VERSION="v0.5.1"
readonly STABLE_KUBERNETES_VERSION="v1.16.3"

readonly ENVS_FILE_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
readonly CLUSTER_CONFIG_FILE="${ENVS_FILE_DIR}/config/kind/cluster-config.yaml"
readonly INGRESS_YAML_FILE="${ENVS_FILE_DIR}/config/kind/ingress.yaml"
readonly MINIO_ADDRESS="localhost:30080"
readonly INGRESS_ADDRESS="http://${MINIO_ADDRESS}"
readonly UPLOAD_SERVICE_ENDPOINT="${INGRESS_ADDRESS}/v1/upload"

readonly DEFAULT_MINIO_ACCESS_KEY="4j4gEuRH96ZFjptUFeFm"
readonly DEFAULT_MINIO_SECRET_KEY="UJnce86xA7hK01WblDdbmXg4gwjKwpFypdLJCvJ3"

readonly RAFTER_CHART="rafter"
readonly RAFTER_CONTROLLER_MANAGER_CHART="rafter-controller-manager"
readonly RAFTER_UPLOAD_SERVICE_CHART="rafter-upload-service"
readonly RAFTER_FRONT_MATTER_SERVICE_CHART="rafter-front-matter-service"
readonly RAFTER_ASYNCAPI_SERVICE_CHART="rafter-asyncapi-service"

# needed globally by run-minio-gateway-test.sh in `azure`` mode
export AZURE_STORAGE_ACCOUNT_NAME=""