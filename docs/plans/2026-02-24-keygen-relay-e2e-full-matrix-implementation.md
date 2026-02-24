# Keygen Relay Full-Matrix E2E Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Add a serial, PR-blocking, full-matrix E2E gate that uses a self-hosted Keygen API server to generate real license files consumed by Relay.

**Architecture:** Reuse the existing Go integration harness and `testscript` style, but add a dedicated `tests/e2e` orchestration layer for Keygen bootstrap, seeding, license export, grouped E2E execution, and artifact collection. CI runs one serial gate job in group order `smoke -> core -> security -> recovery -> chaos` with fail-fast behavior.

**Tech Stack:** Go (`go test`, `testscript`), Bash, Docker Compose, GitHub Actions, SQLite, curl/jq.

---

### Task 1: Create E2E scaffold and group directories

**Files:**
- Create: `tests/e2e/compose/docker-compose.keygen.yml`
- Create: `tests/e2e/scripts/bootstrap_keygen.sh`
- Create: `tests/e2e/scripts/seed_keygen.sh`
- Create: `tests/e2e/scripts/export_license.sh`
- Create: `tests/e2e/scripts/run_group.sh`
- Create: `tests/e2e/scripts/collect_artifacts.sh`
- Create: `tests/e2e/testscript/smoke/.keep`
- Create: `tests/e2e/testscript/core/.keep`
- Create: `tests/e2e/testscript/security/.keep`
- Create: `tests/e2e/testscript/recovery/.keep`
- Create: `tests/e2e/testscript/chaos/.keep`

**Step 1: Write the failing test (structure check)**

```bash
#!/usr/bin/env bash
set -euo pipefail
for d in tests/e2e/compose tests/e2e/scripts tests/e2e/testscript/{smoke,core,security,recovery,chaos}; do
  test -d "$d"
done
```

**Step 2: Run test to verify it fails**

Run: `bash tests/e2e/scripts/check_layout.sh`
Expected: FAIL because files/directories do not exist yet.

**Step 3: Write minimal implementation**

Create the above files and directories with executable bits for scripts.

**Step 4: Run test to verify it passes**

Run: `bash tests/e2e/scripts/check_layout.sh`
Expected: PASS with exit code 0.

**Step 5: Commit**

```bash
git add tests/e2e
git commit -m "test(e2e): scaffold self-hosted keygen test layout"
```

### Task 2: Add Keygen compose bootstrap with health checks

**Files:**
- Modify: `tests/e2e/compose/docker-compose.keygen.yml`
- Modify: `tests/e2e/scripts/bootstrap_keygen.sh`
- Test: `tests/e2e/scripts/bootstrap_keygen.sh`

**Step 1: Write the failing test**

```bash
#!/usr/bin/env bash
set -euo pipefail
COMPOSE_PROJECT_NAME=e2e_bootstrap_test tests/e2e/scripts/bootstrap_keygen.sh
curl -fsS "${KEYGEN_API_URL}/health" >/dev/null
```

**Step 2: Run test to verify it fails**

Run: `bash tests/e2e/scripts/test_bootstrap.sh`
Expected: FAIL until compose and wait logic are implemented.

**Step 3: Write minimal implementation**

- Define compose services and required env wiring.
- In `bootstrap_keygen.sh`: `docker compose up -d`, then poll health endpoint with timeout.

**Step 4: Run test to verify it passes**

Run: `bash tests/e2e/scripts/test_bootstrap.sh`
Expected: PASS and health endpoint is reachable.

**Step 5: Commit**

```bash
git add tests/e2e/compose/docker-compose.keygen.yml tests/e2e/scripts/bootstrap_keygen.sh tests/e2e/scripts/test_bootstrap.sh
git commit -m "test(e2e): add keygen compose bootstrap and readiness checks"
```

### Task 3: Seed Keygen and export real license artifacts

**Files:**
- Modify: `tests/e2e/scripts/seed_keygen.sh`
- Modify: `tests/e2e/scripts/export_license.sh`
- Create: `tests/e2e/assets/keys/test_ed25519_public.key`
- Create: `tests/e2e/assets/keys/test_ed25519_private.key`
- Test: `tests/e2e/scripts/test_seed_and_export.sh`

**Step 1: Write the failing test**

```bash
#!/usr/bin/env bash
set -euo pipefail
workdir="$(mktemp -d)"
COMPOSE_PROJECT_NAME=e2e_seed_test tests/e2e/scripts/bootstrap_keygen.sh
tests/e2e/scripts/seed_keygen.sh
tests/e2e/scripts/export_license.sh "$workdir"
test -s "$workdir/license.lic"
test -s "$workdir/license.key"
```

**Step 2: Run test to verify it fails**

Run: `bash tests/e2e/scripts/test_seed_and_export.sh`
Expected: FAIL because no API seed/export implementation yet.

**Step 3: Write minimal implementation**

- In `seed_keygen.sh`, call Keygen API endpoints to create product/policy/license test fixtures.
- In `export_license.sh`, fetch license payload and key and write deterministic output filenames.

**Step 4: Run test to verify it passes**

Run: `bash tests/e2e/scripts/test_seed_and_export.sh`
Expected: PASS with non-empty `.lic` and `.key` outputs.

**Step 5: Commit**

```bash
git add tests/e2e/scripts/seed_keygen.sh tests/e2e/scripts/export_license.sh tests/e2e/scripts/test_seed_and_export.sh tests/e2e/assets/keys
git commit -m "test(e2e): seed self-hosted keygen and export real license artifacts"
```

### Task 4: Wire grouped testscript runner and serial gate order

**Files:**
- Modify: `tests/e2e/scripts/run_group.sh`
- Create: `tests/e2e/scripts/run_full_matrix.sh`
- Modify: `Makefile`
- Test: `tests/e2e/scripts/test_group_order.sh`

**Step 1: Write the failing test**

```bash
#!/usr/bin/env bash
set -euo pipefail
out="$(mktemp)"
tests/e2e/scripts/run_full_matrix.sh --dry-run >"$out"
grep -n "^smoke$" "$out"
grep -n "^core$" "$out"
grep -n "^security$" "$out"
grep -n "^recovery$" "$out"
grep -n "^chaos$" "$out"
```

**Step 2: Run test to verify it fails**

Run: `bash tests/e2e/scripts/test_group_order.sh`
Expected: FAIL until serial ordering script exists.

**Step 3: Write minimal implementation**

- Add `run_group.sh` interface: `run_group.sh <smoke|core|security|recovery|chaos>`.
- Add `run_full_matrix.sh` to execute groups in approved order with fail-fast.
- Add `make test-e2e-full-matrix` target invoking `run_full_matrix.sh`.

**Step 4: Run test to verify it passes**

Run: `bash tests/e2e/scripts/test_group_order.sh && make test-e2e-full-matrix`
Expected: Dry-run order PASS; full command executes runner contract.

**Step 5: Commit**

```bash
git add tests/e2e/scripts/run_group.sh tests/e2e/scripts/run_full_matrix.sh tests/e2e/scripts/test_group_order.sh Makefile
git commit -m "test(e2e): add serial full-matrix group runner and make target"
```

### Task 5: Implement smoke group with real license flow

**Files:**
- Create: `tests/e2e/testscript/smoke/smoke_real_license_claim_release.test.txt`
- Modify: `cli/cli_test.go`
- Test: `tests/e2e/testscript/smoke/smoke_real_license_claim_release.test.txt`

**Step 1: Write the failing test**

```txt
# smoke: bootstrap and export artifacts must exist
exec test -s $WORK/license.lic
exec test -s $WORK/license.key
# import and serve
exec relay add --file $WORK/license.lic --key $(cat $WORK/license.key) --public-key $(cat $WORK/public.key)
exec relay serve --port $PORT &relay_smoke&
exec curl -s -o claim.json -w "%{http_code}" -X PUT http://localhost:$PORT/v1/nodes/smoke-node
stdout '201'
exec curl -s -o /dev/null -w "%{http_code}" -X DELETE http://localhost:$PORT/v1/nodes/smoke-node
stdout '204'
kill relay_smoke
```

**Step 2: Run test to verify it fails**

Run: `go test -v -tags=integration ./cli -run TestIntegration/smoke`
Expected: FAIL until setup injects generated artifacts.

**Step 3: Write minimal implementation**

- Extend integration setup to place generated artifacts into each script workdir.
- Add smoke script to new group path and run path mapping.

**Step 4: Run test to verify it passes**

Run: `go test -v -tags=integration ./cli -run TestIntegration/smoke`
Expected: PASS.

**Step 5: Commit**

```bash
git add tests/e2e/testscript/smoke cli/cli_test.go
git commit -m "test(e2e): add smoke real-license claim/release flow"
```

### Task 6: Implement core group

**Files:**
- Create: `tests/e2e/testscript/core/core_single_seat_exhaustion.test.txt`
- Create: `tests/e2e/testscript/core/core_floating_seat_cycle.test.txt`
- Create: `tests/e2e/testscript/core/core_pool_header_contracts.test.txt`
- Create: `tests/e2e/testscript/core/core_status_code_contracts.test.txt`
- Test: `tests/e2e/testscript/core/*.test.txt`

**Step 1: Write the failing test**

```txt
# assert single-seat exhaustion
exec curl -s -o /dev/null -w "%{http_code}" -X PUT http://localhost:$PORT/v1/nodes/node-a
stdout '201'
exec curl -s -o /dev/null -w "%{http_code}" -X PUT http://localhost:$PORT/v1/nodes/node-b
stdout '410'
```

**Step 2: Run test to verify it fails**

Run: `tests/e2e/scripts/run_group.sh core`
Expected: FAIL until all core scripts are implemented and wired.

**Step 3: Write minimal implementation**

Add core scripts for seat, pool, strategy, and status-code contract checks using generated real artifacts.

**Step 4: Run test to verify it passes**

Run: `tests/e2e/scripts/run_group.sh core`
Expected: PASS.

**Step 5: Commit**

```bash
git add tests/e2e/testscript/core
git commit -m "test(e2e): add core behavior matrix with real keygen licenses"
```

### Task 7: Implement security group

**Files:**
- Create: `tests/e2e/testscript/security/security_signing_header_and_verification.test.txt`
- Create: `tests/e2e/testscript/security/security_per_license_public_key.test.txt`
- Create: `tests/e2e/testscript/security/security_tamper_watchdog_shutdown.test.txt`
- Create: `tests/e2e/testscript/security/security_audit_on_off.test.txt`
- Test: `tests/e2e/testscript/security/*.test.txt`

**Step 1: Write the failing test**

```txt
# signing enabled should emit Relay-Signature
exec relay serve --port $PORT --signing-secret hunter2 &relay_security&
exec curl -s -D headers.txt -o /dev/null -w "%{http_code}" -X PUT http://localhost:$PORT/v1/nodes/sec-node
stdout '201'
exec grep 'Relay-Signature:' headers.txt
kill relay_security
```

**Step 2: Run test to verify it fails**

Run: `tests/e2e/scripts/run_group.sh security`
Expected: FAIL until scripts and tamper checks are fully implemented.

**Step 3: Write minimal implementation**

- Add signature header/verification scripts.
- Add per-license key mismatch scripts.
- Add tamper simulation and verify-interval watchdog shutdown check.

**Step 4: Run test to verify it passes**

Run: `tests/e2e/scripts/run_group.sh security`
Expected: PASS.

**Step 5: Commit**

```bash
git add tests/e2e/testscript/security
git commit -m "test(e2e): add security matrix for signing and tamper detection"
```

### Task 8: Implement recovery and chaos groups

**Files:**
- Create: `tests/e2e/testscript/recovery/recovery_restart_persistence.test.txt`
- Create: `tests/e2e/testscript/recovery/recovery_no_ghost_claims.test.txt`
- Create: `tests/e2e/testscript/chaos/chaos_relay_restart_recovery.test.txt`
- Create: `tests/e2e/testscript/chaos/chaos_short_ttl_cull_stability.test.txt`
- Create: `tests/e2e/testscript/chaos/chaos_startup_failure_diagnostics.test.txt`
- Test: `tests/e2e/testscript/recovery/*.test.txt`, `tests/e2e/testscript/chaos/*.test.txt`

**Step 1: Write the failing test**

```txt
# recovery: restart and ensure claim/release semantics remain correct
exec relay serve --port $PORT --database $DB &relay_recovery_1&
exec curl -s -o /dev/null -w "%{http_code}" -X PUT http://localhost:$PORT/v1/nodes/recovery-node
stdout '201'
kill relay_recovery_1
exec relay serve --port $PORT --database $DB &relay_recovery_2&
exec curl -s -o /dev/null -w "%{http_code}" -X DELETE http://localhost:$PORT/v1/nodes/recovery-node
stdout '204'
kill relay_recovery_2
```

**Step 2: Run test to verify it fails**

Run: `tests/e2e/scripts/run_group.sh recovery && tests/e2e/scripts/run_group.sh chaos`
Expected: FAIL until scripts are complete.

**Step 3: Write minimal implementation**

Add persistence, interruption, short-TTL, and startup-diagnostic scripts with bounded retries/timeouts.

**Step 4: Run test to verify it passes**

Run: `tests/e2e/scripts/run_group.sh recovery && tests/e2e/scripts/run_group.sh chaos`
Expected: PASS.

**Step 5: Commit**

```bash
git add tests/e2e/testscript/recovery tests/e2e/testscript/chaos
git commit -m "test(e2e): add recovery and chaos full-matrix coverage"
```

### Task 9: Add artifact collection and CI PR gate

**Files:**
- Modify: `tests/e2e/scripts/collect_artifacts.sh`
- Modify: `.github/workflows/test.yml`
- Modify: `Makefile`
- Test: `.github/workflows/test.yml`

**Step 1: Write the failing test**

```bash
#!/usr/bin/env bash
set -euo pipefail
rg -n "test-e2e-full-matrix" .github/workflows/test.yml
```

**Step 2: Run test to verify it fails**

Run: `bash tests/e2e/scripts/test_ci_gate.sh`
Expected: FAIL because workflow does not yet call full matrix target.

**Step 3: Write minimal implementation**

- Add PR step invoking `make test-e2e-full-matrix`.
- Ensure artifacts are uploaded on failure.
- Keep existing build/vet/fmt/unit/integration checks intact.

**Step 4: Run test to verify it passes**

Run: `bash tests/e2e/scripts/test_ci_gate.sh`
Expected: PASS and workflow contains serial gate invocation.

**Step 5: Commit**

```bash
git add .github/workflows/test.yml Makefile tests/e2e/scripts/collect_artifacts.sh tests/e2e/scripts/test_ci_gate.sh
git commit -m "ci: enforce serial full-matrix e2e gate with artifacts"
```

### Task 10: Verification and documentation updates

**Files:**
- Modify: `README.md`
- Create: `tests/e2e/README.md`
- Test: Full verification command list

**Step 1: Write the failing test**

```bash
#!/usr/bin/env bash
set -euo pipefail
rg -n "test-e2e-full-matrix|self-hosted Keygen|full matrix" README.md tests/e2e/README.md
```

**Step 2: Run test to verify it fails**

Run: `bash tests/e2e/scripts/test_docs_refs.sh`
Expected: FAIL until docs are added/updated.

**Step 3: Write minimal implementation**

- Add `tests/e2e/README.md` with local run instructions and troubleshooting.
- Update root `README.md` testing section with new E2E gate command.

**Step 4: Run test to verify it passes**

Run:
- `make test`
- `make test-integration`
- `make test-e2e-full-matrix`
Expected: PASS locally or in CI runner environment.

**Step 5: Commit**

```bash
git add README.md tests/e2e/README.md tests/e2e/scripts/test_docs_refs.sh
git commit -m "docs: document self-hosted keygen full-matrix e2e workflow"
```

## Final Verification Checklist

Run in order:
1. `make fmt`
2. `make vet`
3. `make test`
4. `make test-integration`
5. `make test-e2e-full-matrix`

Expected outcome:
- All commands pass.
- Full matrix runs in serial group order.
- Failed groups produce complete artifacts for diagnosis.

## Skill References

- @test-driven-development
- @verification-before-completion
- @requesting-code-review
