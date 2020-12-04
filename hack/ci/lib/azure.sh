#!/usr/bin/env bash
# source: https://github.com/kyma-project/test-infra/blob/bb9d70fc141529637e099279f6de0418c4fa6b3b/prow/scripts/lib/azure.sh

LIBDIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" || exit; pwd)"

# shellcheck source=prow/scripts/lib/log.sh
source "${LIBDIR}"/log.sh

function az::verify_deps {
  if ! [[ -x $(command -v az) ]]; then
    log::error "'az' command not found in \$PATH. Exiting..."
    exit 1
  else
    echo "Azure CLI Version:"
    az version
  fi
  if ! [[ -x $(command -v jq) ]]; then
    log::error "'jq' command not found in \$PATH. Exiting..."
    exit 1
  else
    echo "jq version:"
    jq --version
  fi
}

# az::login logs in to the azure service using provided credentials file in the function argument.
# Function accepts JSON file formatted below:
# {
#   "tenant_id": "tenant_id",
#   "app_id": "subscription_app_id",
#   "secret": "subscription_secret"
# }
function az::login {
  local AZURE_SUBSCRIPTION_TENANT
  local AZURE_SUBSCRIPTION_APP_ID
  local AZURE_SUBSCRIPTION_SECRET

  # Check the provided credentials in the argument.
  # Use arguments to avoid exporting sensitive values.
  if [[ -z "$1" ]]; then
    log::error "Azure credentials file not provided. please provide azure credentials filepath in the argument. Exiting..."
    exit 1
  elif [[ ! -f "$1" ]]; then
    log::error "Azure credentials file not found. Make sure it is present under the provided filepath. Exiting..."
    exit 1
  fi
  AZURE_CREDENTIALS_FILE="$1"
  AZURE_SUBSCRIPTION_TENANT=$(jq -r '.tenant_id' "$AZURE_CREDENTIALS_FILE")
  AZURE_SUBSCRIPTION_APP_ID=$(jq -r '.app_id' "$AZURE_CREDENTIALS_FILE")
  AZURE_SUBSCRIPTION_SECRET=$(jq -r '.secret' "$AZURE_CREDENTIALS_FILE")

  # login
  log::info "Logging in to Azure"
  az login --service-principal -u "${AZURE_SUBSCRIPTION_APP_ID}" -p "${AZURE_SUBSCRIPTION_SECRET}" --tenant "${AZURE_SUBSCRIPTION_TENANT}"
  log::info "Successfully logged-in!"
}

# az::set_subscription sets the subscription using provided subscription ID in the argument.
function az::set_subscription {
  if [[ -z "$1" ]]; then
    log::error "Azure Subscription ID is not found. Please provide azure subscription ID in the argument. Exiting..."
    exit 1
  fi
  az account set \
    --subscription "$1"
}

# This check will trigger everytime the file is sourced.
# This should allow easy checking for the related requirements without copying loads of code.
# If more checks will be needed we should add another function there.
az::verify_deps
