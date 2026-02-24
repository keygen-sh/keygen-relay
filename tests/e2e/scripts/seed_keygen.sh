#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
E2E_DIR="$(cd "${SCRIPT_DIR}/.." && pwd)"
RUNTIME_DIR="${E2E_DIR}/.runtime"
ENV_FILE="${RUNTIME_DIR}/keygen.env"
STATE_FILE="${RUNTIME_DIR}/seed-state.json"
COMPOSE_FILE="${E2E_DIR}/compose/docker-compose.keygen.yml"

COMPOSE_PROJECT_NAME="${COMPOSE_PROJECT_NAME:-keygen_e2e}"

if [[ ! -f "${ENV_FILE}" ]]; then
  echo "missing runtime env file at ${ENV_FILE}; run bootstrap_keygen.sh up first" >&2
  exit 1
fi

# shellcheck source=/dev/null
source "${ENV_FILE}"

KEYGEN_API_URL="${KEYGEN_API_URL:-http://127.0.0.1:${KEYGEN_HTTP_PORT:-63000}}"
KEYGEN_E2E_MAX_MACHINES="${KEYGEN_E2E_MAX_MACHINES:-3}"

compose_cmd() {
  docker compose \
    --project-name "${COMPOSE_PROJECT_NAME}" \
    --env-file "${ENV_FILE}" \
    -f "${COMPOSE_FILE}" \
    "$@"
}

require_cmds() {
  local cmd
  for cmd in curl jq date; do
    command -v "${cmd}" >/dev/null 2>&1 || {
      echo "required command not found: ${cmd}" >&2
      exit 1
    }
  done
}

json_api_post() {
  local url="$1"
  local token="$2"
  local body="$3"
  curl -fsS -X POST "${url}" \
    -H "Authorization: Bearer ${token}" \
    -H "X-Forwarded-Proto: https" \
    -H "Accept: application/vnd.api+json" \
    -H "Content-Type: application/vnd.api+json" \
    -d "${body}"
}

get_admin_token() {
  curl -fsS -X POST "${KEYGEN_API_URL}/v1/tokens" \
    -u "${KEYGEN_ADMIN_EMAIL}:${KEYGEN_ADMIN_PASSWORD}" \
    -H "X-Forwarded-Proto: https" \
    -H "Accept: application/vnd.api+json" \
    -H "Content-Type: application/vnd.api+json" \
    -d "{}" | jq -er '.data.attributes.token'
}

main() {
  require_cmds

  local token now suffix product_code product_name policy_name license_json
  local product_id policy_id license_id license_key certificate public_key

  token="$(get_admin_token)"
  now="$(date +%s)"
  suffix="${now}"
  product_code="RELAY-E2E-${suffix}"
  product_name="Relay E2E Product ${suffix}"
  policy_name="Relay E2E Policy ${suffix}"

  product_id="$(json_api_post \
    "${KEYGEN_API_URL}/v1/products" \
    "${token}" \
    "{\"data\":{\"type\":\"products\",\"attributes\":{\"name\":\"${product_name}\",\"code\":\"${product_code}\"}}}" \
    | jq -er '.data.id')"

  policy_id="$(json_api_post \
    "${KEYGEN_API_URL}/v1/policies" \
    "${token}" \
    "{\"data\":{\"type\":\"policies\",\"attributes\":{\"name\":\"${policy_name}\",\"floating\":true,\"maxMachines\":${KEYGEN_E2E_MAX_MACHINES}},\"relationships\":{\"product\":{\"data\":{\"type\":\"products\",\"id\":\"${product_id}\"}}}}}" \
    | jq -er '.data.id')"

  license_json="$(json_api_post \
    "${KEYGEN_API_URL}/v1/licenses" \
    "${token}" \
    "{\"data\":{\"type\":\"licenses\",\"relationships\":{\"policy\":{\"data\":{\"type\":\"policies\",\"id\":\"${policy_id}\"}}}}}")"
  license_id="$(printf '%s' "${license_json}" | jq -er '.data.id')"
  license_key="$(printf '%s' "${license_json}" | jq -er '.data.attributes.key')"

  certificate="$(json_api_post \
    "${KEYGEN_API_URL}/v1/licenses/${license_id}/actions/check-out" \
    "${token}" \
    "{}" | jq -er '.data.attributes.certificate')"

  public_key="$(compose_cmd exec -T postgres psql -U "${POSTGRES_USER}" -d "${POSTGRES_DB}" -At -c "select ed25519_public_key from accounts where id='${KEYGEN_ACCOUNT_ID}';" | tr -d '\r')"

  jq -n \
    --arg api_url "${KEYGEN_API_URL}" \
    --arg account_id "${KEYGEN_ACCOUNT_ID}" \
    --arg product_id "${product_id}" \
    --arg policy_id "${policy_id}" \
    --arg license_id "${license_id}" \
    --arg license_key "${license_key}" \
    --arg certificate "${certificate}" \
    --arg public_key "${public_key}" \
    '{
      api_url: $api_url,
      account_id: $account_id,
      product_id: $product_id,
      policy_id: $policy_id,
      license_id: $license_id,
      license_key: $license_key,
      certificate: $certificate,
      public_key: $public_key
    }' >"${STATE_FILE}"

  echo "seeded keygen resources:"
  echo "  account_id=${KEYGEN_ACCOUNT_ID}"
  echo "  product_id=${product_id}"
  echo "  policy_id=${policy_id}"
  echo "  license_id=${license_id}"
  echo "  state_file=${STATE_FILE}"
}

main "$@"
