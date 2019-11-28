#!/usr/bin/env bash

readonly __ENVS_FILE_DIR__=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)

readonly __STABLE_HELM_VERSION__="v2.16.1"
readonly __STABLE_KIND_VERSION__="v0.5.1"
readonly __STABLE_KUBERNETES_VERSION__="v1.16.3"

readonly __RAFTER__="rafter"
readonly __RAFTER_CONTROLLER_MANAGER__="rafter-controller-manager"
readonly __RAFTER_UPLOAD_SERVICE__="rafter-upload-service"
readonly __RAFTER_FRONT_MATTER_SERVICE__="rafter-front-matter-service"
readonly __RAFTER_ASYNCAPI_SERVICE__="rafter-asyncapi-service"

readonly __DEFAULT_MINIO_ACCESS_KEY__="4j4gEuRH96ZFjptUFeFm"
readonly __DEFAULT_MINIO_SECRET_KEY__="UJnce86xA7hK01WblDdbmXg4gwjKwpFypdLJCvJ3"

readonly __CLUSTER_CONFIG_FILE__="${__ENVS_FILE_DIR__}/config/kind/cluster-config.yaml"
readonly __INGRESS_YAML_FILE__="${__ENVS_FILE_DIR__}/config/kind/ingress.yaml"
readonly __MINIO_ADDRESS__="localhost:30080"
readonly __INGRESS_ADDRESS__="http://${__MINIO_ADDRESS__}"
readonly __UPLOAD_SERVICE_ENDPOINT__="${__INGRESS_ADDRESS__}/v1/upload"

readonly __MINIO_GATEWAY_TEST_BASIC__="basic"
readonly __MINIO_GATEWAY_TEST_MIGRATION__="migration"

readonly __MINIO_GATEWAY_PROVIDER_GCS__="gcs"
readonly __MINIO_GATEWAY_PROVIDER_AZURE__="azure"