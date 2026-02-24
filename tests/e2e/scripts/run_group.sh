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
  echo "== group summary (${GROUP}) =="
  echo "0/0 PASS"
  exit 0
fi

log_file="${artifacts_dir}/go-test.log"
status_file="${artifacts_dir}/test-status.tsv"

set +e
TESTSCRIPT_DIR="${group_dir}" E2E_EXPORT_DIR="${artifacts_dir}" go test -v -tags=integration ./cli 2>&1 | tee "${log_file}"
test_exit="${PIPESTATUS[0]}"
set -e

awk '
/^[[:space:]]*--- PASS: TestIntegration\// {
  if (match($0, /TestIntegration\/[^ ]+/)) {
    name = substr($0, RSTART, RLENGTH)
    sub(/^TestIntegration\//, "", name)
    status[name] = "PASS"
  }
  next
}
/^[[:space:]]*--- FAIL: TestIntegration\// {
  if (match($0, /TestIntegration\/[^ ]+/)) {
    name = substr($0, RSTART, RLENGTH)
    sub(/^TestIntegration\//, "", name)
    status[name] = "FAIL"
  }
  next
}
END {
  for (name in status) {
    print name "\t" status[name]
  }
}
' "${log_file}" > "${status_file}"

pass_count=0
fail_count=0

echo "== group summary (${GROUP}) =="
while IFS= read -r test_file; do
  test_name="$(basename "${test_file}" .txt)"
  test_result="$(awk -F'\t' -v t="${test_name}" '$1==t { print $2; found=1; exit } END { if (!found) print "NOT_RUN" }' "${status_file}")"

  case "${test_result}" in
  PASS)
    ((pass_count++))
    ;;
  FAIL | NOT_RUN)
    ((fail_count++))
    ;;
  esac

  echo "${test_name}: ${test_result}"
done < <(find "${group_dir}" -type f -name '*.test.txt' | sort)

echo "${pass_count}/${test_count} PASS"

if [[ "${test_exit}" -ne 0 ]]; then
  exit "${test_exit}"
fi
