#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

__init_environment() {
  local -r current_dir=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)

  source "${current_dir}/lib/test-helpers.sh" || {
    echo '- Cannot load test helpers.'
    exit 1
  }

  # Default is 20s - available since controller-runtime 0.1.5
  export KUBEBUILDER_CONTROLPLANE_START_TIMEOUT=2m
  # Default is 20s - available since controller-runtime 0.1.5
  export KUBEBUILDER_CONTROLPLANE_STOP_TIMEOUT=2m

  testHelpers::install_go_junit_report
}

# Arguments:
#   $1 - Tmp directory with binaries used during tests
__cleanup() {
  local -r bin_dir="${1}"

  log::info "- Deleting directory with temporary binaries used in tests..."
  rm -rf "${bin_dir}" || true
  log::success "- Directory with temporary binaries used in tests deleted."
}

# Arguments:
#   $1 - Root path of repo
main() {
  local -r root_path="${1}"

  __init_environment

  local -r tmp_dir="$(mktemp -d)"
  export ARTIFACTS_DIR="${ARTIFACTS:-"${tmp_dir}/artifacts"}"
  mkdir -p "${ARTIFACTS_DIR}"

  trap "__cleanup ${tmp_dir}" EXIT

  local -r log_file=unit_test_data.log
  local -r coverage_file="cover.out"
  local -r suite_name="Rafter_Unit_Tests"
  local test_failed="false"

  log::info "- Starting unit tests..."

  go test "${root_path}"/... -count 1 -coverprofile="${ARTIFACTS_DIR}/${coverage_file}" -v 2>&1 | tee "${ARTIFACTS_DIR}/${log_file}" || test_failed="true"
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
