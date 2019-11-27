#!/usr/bin/env bash

set -eo pipefail

CURRENT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)

source "${CURRENT_DIR}/test-helper.sh" || {
    echo 'Cannot load test helper.'
    exit 1
}

readonly TMP_DIR="$(mktemp -d)"
readonly TMP_BIN_DIR="${TMP_DIR}/bin"
mkdir -p "${TMP_BIN_DIR}"
readonly ARTIFACTS_DIR="${ARTIFACTS:-"${TMP_DIR}/artifacts"}"
mkdir -p "${ARTIFACTS_DIR}"

readonly ROOT="${1}"

cleanup(){
    log::info "- Cleaning up temporary directory..."
    rm -rf "${1}"
}

trap 'cleanup ${TMP_DIR}' EXIT

main(){    
    local test_failed="false"
    local -r JUNIT_SUITE_NAME="Rafter_Unit_Tests"
    local -r COVERAGE_FILENAME="cover.out"
    local -r LOG_FILE=unit_test_data.log
    testHelper::install_go_junit_report

    # Default is 20s - available since controller-runtime 0.1.5
    export KUBEBUILDER_CONTROLPLANE_START_TIMEOUT=2m
    # Default is 20s - available since controller-runtime 0.1.5
    export KUBEBUILDER_CONTROLPLANE_STOP_TIMEOUT=2m
    
    log::info "- Running unit tests..."
    go test "${ROOT}"/... -count 1 -coverprofile="${ARTIFACTS_DIR}/${COVERAGE_FILENAME}" -v 2>&1 | tee "${ARTIFACTS_DIR}/${LOG_FILE}" || test_failed="true"
    < "${ARTIFACTS_DIR}/${LOG_FILE}" go-junit-report > "${ARTIFACTS_DIR}/junit_${JUNIT_SUITE_NAME}_suite.xml"
    go tool cover -func="${ARTIFACTS_DIR}/${COVERAGE_FILENAME}" \
		| grep total \
		| awk '{print "Total test coverage: " $3}'

    if [[ ${test_failed} = "true" ]]; then
        log::error "Finished with error"
        return 1
    fi
    log::success "Finished with success"
}

main