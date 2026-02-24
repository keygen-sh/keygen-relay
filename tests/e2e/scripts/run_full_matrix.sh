#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
groups=(smoke core security recovery chaos)

if [[ "${1:-}" == "--dry-run" ]]; then
  printf '%s\n' "${groups[@]}"
  exit 0
fi

for group in "${groups[@]}"; do
  echo "== running ${group} =="
  "${SCRIPT_DIR}/run_group.sh" "${group}"
done
