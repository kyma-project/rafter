#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

readonly STABLE_KUBERNETES_VERSION=v1.16.3
readonly STABLE_KIND_VERSION=v0.5.1
readonly STABLE_HELM_VERSION=v2.16.0
readonly CLUSTER_NAME=ci-test-cluster