#!/usr/bin/env bash
set -euo pipefail

for d in \
  tests/e2e/compose \
  tests/e2e/scripts \
  tests/e2e/testscript/smoke \
  tests/e2e/testscript/core \
  tests/e2e/testscript/security \
  tests/e2e/testscript/recovery \
  tests/e2e/testscript/chaos
  do
  test -d "$d"
done

for f in \
  tests/e2e/compose/docker-compose.keygen.yml \
  tests/e2e/scripts/bootstrap_keygen.sh \
  tests/e2e/scripts/seed_keygen.sh \
  tests/e2e/scripts/export_license.sh \
  tests/e2e/scripts/run_group.sh \
  tests/e2e/scripts/collect_artifacts.sh
  do
  test -f "$f"
done

echo "layout ok"
