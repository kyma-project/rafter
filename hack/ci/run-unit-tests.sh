#!/usr/bin/env bash

# Description: 
#   This scripts implements the flow of Rafter Unit tests for CI.
#
# Required parameters:
#   - $1 - Absolute path to local Rafter repository

set -o errexit
set -o nounset
set -o pipefail

readonly ROOT_REPO_PATH="${1}"

readonly TMP_DIR="$(mktemp -d)"
export ARTIFACTS_DIR="${ARTIFACTS:-"${TMP_DIR}/artifacts"}"
mkdir -p "${ARTIFACTS_DIR}"

init_environment() {
  local -r current_dir=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)

  source "${current_dir}/lib/test-helpers.sh" || {
    echo '- Cannot load test helpers.'
    return 1
  }

  # Default is 20s - available since controller-runtime 0.1.5
  export KUBEBUILDER_CONTROLPLANE_START_TIMEOUT=2m
  # Default is 20s - available since controller-runtime 0.1.5
  export KUBEBUILDER_CONTROLPLANE_STOP_TIMEOUT=2m

  testHelpers::install_go_junit_report
}

cleanup() {
  log::info "- Deleting directory with temporary binaries used in tests..."
  rm -rf "${TMP_DIR}" || true
  log::success "- Directory with temporary binaries used in tests deleted."
}
trap "cleanup" EXIT

main() {
  init_environment

  local -r log_file=unit_test_data.log
  local -r coverage_file="cover.out"
  local -r suite_name="Rafter_Unit_Tests"
  local test_failed="false"

  log::info "- Starting unit tests..."

  go test "${ROOT_REPO_PATH}"/... -count 1 -coverprofile="${ARTIFACTS_DIR}/${coverage_file}" -v 2>&1 | tee "${ARTIFACTS_DIR}/${log_file}" || test_failed="true"
  < "${ARTIFACTS_DIR}/${log_file}" go-junit-report > "${ARTIFACTS_DIR}/junit_${suite_name}_suite.xml"
  go tool cover -func="${ARTIFACTS_DIR}/${coverage_file}" \
		| grep total \
		| awk '{print "Total test coverage: " $3}'

  if [[ ${test_failed} = "true" ]]; then
    log::error "- Unit tests failed."
    return 1
  fi
  log::success "- Unit tests passed."
}
main
