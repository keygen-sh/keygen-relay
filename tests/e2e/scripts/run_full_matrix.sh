#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
E2E_DIR="$(cd "${SCRIPT_DIR}/.." && pwd)"
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

echo ""
echo "== detailed results =="
total_pass=0
total_fail=0
for group in "${groups[@]}"; do
  status_file="${E2E_DIR}/artifacts/${group}/test-status.tsv"
  if [[ ! -f "${status_file}" ]]; then
    continue
  fi

  group_pass=0
  group_fail=0
  group_total=0

  echo ""
  echo "[${group}]"
  while IFS=$'\t' read -r test_name status; do
    ((group_total++))
    case "${status}" in
      PASS)
        ((group_pass++))
        ((total_pass++))
        ;;
      FAIL|NOT_RUN)
        ((group_fail++))
        ((total_fail++))
        ;;
    esac
    echo "  ${test_name}: ${status}"
  done < "${status_file}"
  echo "  ${group_pass}/${group_total} PASS"
done

echo ""
echo "total: ${total_pass} PASS, ${total_fail} FAIL"

if [[ "${failed}" -ne 0 ]]; then
  exit 1
fi
