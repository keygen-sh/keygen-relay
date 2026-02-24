#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
groups=(smoke core security recovery chaos)
results=()
failed=0

if [[ "${1:-}" == "--dry-run" ]]; then
  printf '%s\n' "${groups[@]}"
  exit 0
fi

for group in "${groups[@]}"; do
  echo "== running ${group} =="
  if "${SCRIPT_DIR}/run_group.sh" "${group}"; then
    results+=("${group}: PASS")
  else
    results+=("${group}: FAIL")
    failed=1
    break
  fi
done

echo "== summary =="
for result in "${results[@]}"; do
  echo "${result}"
done

if [[ "${failed}" -ne 0 ]]; then
  exit 1
fi
