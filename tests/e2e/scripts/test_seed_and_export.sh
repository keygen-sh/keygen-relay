#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BOOTSTRAP_SCRIPT="${SCRIPT_DIR}/bootstrap_keygen.sh"
SEED_SCRIPT="${SCRIPT_DIR}/seed_keygen.sh"
EXPORT_SCRIPT="${SCRIPT_DIR}/export_license.sh"

export COMPOSE_PROJECT_NAME="e2e_seed_test"
export KEYGEN_HTTP_PORT="63012"

outdir="$(mktemp -d)"

cleanup() {
  "${BOOTSTRAP_SCRIPT}" down >/dev/null 2>&1 || true
  rm -rf "${outdir}"
}

trap cleanup EXIT

"${BOOTSTRAP_SCRIPT}" up >/dev/null
"${SEED_SCRIPT}" >/dev/null
"${EXPORT_SCRIPT}" "${outdir}" >/dev/null

test -s "${outdir}/license.lic"
test -s "${outdir}/license.key"
test -s "${outdir}/public.key"

grep -q "BEGIN LICENSE FILE" "${outdir}/license.lic"
rg -q "^[A-Z0-9-]+$" "${outdir}/license.key"
rg -q "^[a-f0-9]{64}$" "${outdir}/public.key"

echo "seed and export ok"
