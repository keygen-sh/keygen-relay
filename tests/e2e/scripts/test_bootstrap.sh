#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BOOTSTRAP_SCRIPT="${SCRIPT_DIR}/bootstrap_keygen.sh"

export COMPOSE_PROJECT_NAME="e2e_bootstrap_test"
export KEYGEN_HTTP_PORT="63010"

cleanup() {
  "${BOOTSTRAP_SCRIPT}" down >/dev/null 2>&1 || true
}

trap cleanup EXIT

"${BOOTSTRAP_SCRIPT}" up >/dev/null

code="$(curl -sS -o /dev/null -w '%{http_code}' "http://127.0.0.1:${KEYGEN_HTTP_PORT}/v1/health" || true)"
if [[ "${code}" != "200" && "${code}" != "204" ]]; then
  echo "expected keygen health 200 or 204, got ${code}" >&2
  exit 1
fi

echo "bootstrap ok"
