#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
E2E_DIR="$(cd "${SCRIPT_DIR}/.." && pwd)"
GROUP="${1:-}"

if [[ -z "${GROUP}" ]]; then
  echo "usage: run_group.sh <smoke|core|security|recovery|chaos>" >&2
  exit 2
fi

case "${GROUP}" in
smoke) port=63100 ;;
core) port=63110 ;;
security) port=63120 ;;
recovery) port=63130 ;;
chaos) port=63140 ;;
*)
  echo "invalid group: ${GROUP}" >&2
  exit 2
  ;;
esac

export COMPOSE_PROJECT_NAME="keygen_e2e_${GROUP}"
export KEYGEN_HTTP_PORT="${KEYGEN_HTTP_PORT:-$port}"

group_dir="${E2E_DIR}/testscript/${GROUP}"
artifacts_dir="${E2E_DIR}/artifacts/${GROUP}"

cleanup() {
  "${SCRIPT_DIR}/bootstrap_keygen.sh" down >/dev/null 2>&1 || true
}
trap cleanup EXIT

mkdir -p "${artifacts_dir}"

"${SCRIPT_DIR}/bootstrap_keygen.sh" up
"${SCRIPT_DIR}/seed_keygen.sh"
"${SCRIPT_DIR}/export_license.sh" "${artifacts_dir}"

test_count="$(find "${group_dir}" -type f -name '*.test.txt' | wc -l | tr -d ' ')"
echo "group=${GROUP} tests=${test_count}"

if [[ "${test_count}" == "0" ]]; then
  echo "no testscript files in ${group_dir}; skipping execution"
  exit 0
fi

TESTSCRIPT_DIR="${group_dir}" E2E_EXPORT_DIR="${artifacts_dir}" go test -v -tags=integration ./cli
