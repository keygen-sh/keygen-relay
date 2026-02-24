#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
E2E_DIR="$(cd "${SCRIPT_DIR}/.." && pwd)"
RUNTIME_DIR="${E2E_DIR}/.runtime"
ENV_FILE="${RUNTIME_DIR}/keygen.env"
COMPOSE_FILE="${E2E_DIR}/compose/docker-compose.keygen.yml"

KEYGEN_HTTP_PORT="${KEYGEN_HTTP_PORT:-63000}"
COMPOSE_PROJECT_NAME="${COMPOSE_PROJECT_NAME:-keygen_e2e}"
KEYGEN_API_URL="${KEYGEN_API_URL:-http://127.0.0.1:${KEYGEN_HTTP_PORT}}"

usage() {
  cat <<'USAGE'
Usage:
  bootstrap_keygen.sh up
  bootstrap_keygen.sh down
USAGE
}

ensure_runtime_env() {
  mkdir -p "${RUNTIME_DIR}"

  local account_id
  account_id="$(uuidgen | tr '[:upper:]' '[:lower:]')"

  cat >"${ENV_FILE}" <<EOF
POSTGRES_USER=keygen
POSTGRES_PASSWORD=keygen
POSTGRES_DB=keygen
REDIS_URL=redis://redis:6379
SECRET_KEY_BASE=$(openssl rand -hex 64)
ENCRYPTION_DETERMINISTIC_KEY=$(openssl rand -base64 32)
ENCRYPTION_PRIMARY_KEY=$(openssl rand -base64 32)
ENCRYPTION_KEY_DERIVATION_SALT=$(openssl rand -base64 32)
KEYGEN_EDITION=CE
KEYGEN_MODE=singleplayer
KEYGEN_ACCOUNT_ID=${account_id}
KEYGEN_ADMIN_EMAIL=e2e-admin@example.com
KEYGEN_ADMIN_PASSWORD=password123
KEYGEN_HOST=127.0.0.1:${KEYGEN_HTTP_PORT}
KEYGEN_HTTP_PORT=${KEYGEN_HTTP_PORT}
EOF
}

compose_cmd() {
  docker compose \
    --project-name "${COMPOSE_PROJECT_NAME}" \
    --env-file "${ENV_FILE}" \
    -f "${COMPOSE_FILE}" \
    "$@"
}

wait_for_keygen() {
  local retries=90
  local wait_secs=2
  local health_code

  for _ in $(seq 1 "${retries}"); do
    health_code="$(curl -sS -o /dev/null -w '%{http_code}' "${KEYGEN_API_URL}/v1/health" || true)"
    if [[ "${health_code}" == "200" || "${health_code}" == "204" ]]; then
      return 0
    fi
    sleep "${wait_secs}"
  done

  echo "keygen web health check failed at ${KEYGEN_API_URL}" >&2
  compose_cmd logs web >&2 || true
  return 1
}

up() {
  ensure_runtime_env
  compose_cmd up -d postgres redis
  compose_cmd --profile setup run --rm setup
  compose_cmd up -d web
  wait_for_keygen

  echo "KEYGEN_API_URL=${KEYGEN_API_URL}"
  echo "KEYGEN_ENV_FILE=${ENV_FILE}"
}

down() {
  if [[ ! -f "${ENV_FILE}" ]]; then
    echo "no runtime env found at ${ENV_FILE}, nothing to tear down"
    return 0
  fi

  compose_cmd down -v --remove-orphans
}

main() {
  local cmd="${1:-up}"
  case "${cmd}" in
  up)
    up
    ;;
  down)
    down
    ;;
  *)
    usage >&2
    exit 2
    ;;
  esac
}

main "$@"
