#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
out="$(mktemp)"
trap 'rm -f "$out"' EXIT

"${SCRIPT_DIR}/run_full_matrix.sh" --dry-run >"$out"

actual="$(cat "$out")"
expected=$'smoke\ncore\nsecurity\nrecovery\nchaos'

if [[ "$actual" != "$expected" ]]; then
  echo "unexpected group order:" >&2
  cat "$out" >&2
  exit 1
fi

echo "group order ok"
