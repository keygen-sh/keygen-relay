#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
E2E_DIR="$(cd "${SCRIPT_DIR}/.." && pwd)"
RUNTIME_DIR="${E2E_DIR}/.runtime"
STATE_FILE="${RUNTIME_DIR}/seed-state.json"

OUT_DIR="${1:-}"
if [[ -z "${OUT_DIR}" ]]; then
  echo "usage: export_license.sh <output-dir>" >&2
  exit 2
fi

if [[ ! -f "${STATE_FILE}" ]]; then
  echo "missing seed state at ${STATE_FILE}; run seed_keygen.sh first" >&2
  exit 1
fi

mkdir -p "${OUT_DIR}"

jq -er '.certificate' "${STATE_FILE}" >"${OUT_DIR}/license.lic"
jq -er '.license_key' "${STATE_FILE}" >"${OUT_DIR}/license.key"
jq -er '.public_key' "${STATE_FILE}" >"${OUT_DIR}/public.key"

echo "exported:"
echo "  ${OUT_DIR}/license.lic"
echo "  ${OUT_DIR}/license.key"
echo "  ${OUT_DIR}/public.key"
