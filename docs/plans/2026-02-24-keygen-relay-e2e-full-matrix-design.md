# Keygen Relay Comprehensive E2E Test Design (Full Matrix, Self-Hosted Keygen)

Date: 2026-02-24
Status: Approved

## 1. Context and Objectives

Keygen Relay already includes integration tests via `go test -tags=integration` and `testscript`, covering core claim/release and several pool/signing flows. The goal is to expand this into a comprehensive end-to-end suite that validates full behavior against real license files generated from a self-hosted Keygen API server.

Primary objectives:
- Validate real license generation and consumption end-to-end.
- Cover business-critical, security-critical, recovery, and resilience paths.
- Keep tests deterministic and reproducible in CI.
- Gate PR merges on the full approved matrix.

## 2. Chosen Approach

Selected approach: Docker Compose with scripted Keygen bootstrap/seed/export, integrated with existing test harness patterns.

Rationale:
- Lowest deployment and adoption cost.
- High compatibility with existing `testscript` integration style.
- Strong debuggability via plain container/process logs.

Not selected:
- Testcontainers-first: stronger isolation but higher migration/maintenance cost.
- Pre-seeded snapshots: faster but weaker realism and higher drift risk.

## 3. Execution Model (Updated)

Execution will be **serial**, not parallel.

Reason:
- Team accepts longer runtime.
- Reduces peak CI resource usage.
- Simplifies infrastructure and failure triage.

Pipeline shape:
- One `e2e-full-matrix` PR gate job.
- Groups run in order:
  1. `smoke`
  2. `core`
  3. `security`
  4. `recovery`
  5. `chaos`
- Fail-fast policy: stop on first failing group.

## 4. Self-Hosted Keygen Lifecycle

Each pipeline run executes:
1. `docker compose up -d` for Keygen API dependencies.
2. Health check loop until ready.
3. Seed script creates account/product/policy/license test data.
4. Export script writes real `.lic` and `.key` artifacts for Relay tests.
5. Relay tests consume artifacts via `relay add --file --key --public-key`.
6. On completion/failure, collect artifacts and tear down with `docker compose down -v`.

Key handling:
- Use test-only Ed25519 keypair(s).
- Store per-license public key and assert per-license verification behavior.
- Cover invalid/mismatched key scenarios in security group.

## 5. Test Matrix

### 5.1 Smoke
- Keygen boot/seed/export success.
- Relay add/serve/health success.
- Single claim (201) and release (204) happy path.

### 5.2 Core
- Single-seat exhaustion behavior.
- Floating-seat claim/release cycles.
- Same-fingerprint second claim semantics (202 vs 409 under no-heartbeats).
- Pool routing and header semantics.
- Status code contract coverage: 201/202/204/400/404/409/410.
- Strategy behavior checks for fifo/lifo/rand (rand asserted statistically).

### 5.3 Security
- Response signing header behavior with/without signing secret.
- Signature verification failure scenarios.
- Per-license public key correctness/incorrectness checks.
- Seat tampering detection via watchdog/verify interval.
- Audit logging on/off behavior checks.

### 5.4 Recovery
- Restart consistency with persisted database path.
- Lease/pool state correctness after restart.
- No ghost claims after restart.
- Migration compatibility read-path checks.

### 5.5 Chaos
- Temporary Keygen unavailability during Relay operations.
- Relay process interruption and restart recovery.
- Short TTL + high-frequency cull stability.
- Concurrent claim conflict behavior and idempotent retry semantics.
- Startup failure diagnosability (port conflicts, bad DB path).

## 6. Repository Layout

Planned structure:
- `tests/e2e/compose/docker-compose.keygen.yml`
- `tests/e2e/scripts/bootstrap_keygen.sh`
- `tests/e2e/scripts/seed_keygen.sh`
- `tests/e2e/scripts/export_license.sh`
- `tests/e2e/scripts/run_group.sh`
- `tests/e2e/scripts/collect_artifacts.sh`
- `tests/e2e/testscript/smoke/*`
- `tests/e2e/testscript/core/*`
- `tests/e2e/testscript/security/*`
- `tests/e2e/testscript/recovery/*`
- `tests/e2e/testscript/chaos/*`
- `tests/e2e/assets/keys/*`
- `tests/e2e/artifacts/*`

## 7. Gate Criteria

PR gate passes only when:
- All serial groups pass in order.
- No panic/data race/fatal process errors.
- Contracted status codes and response fields match expected behavior.
- Security/tamper tests trigger expected protective behavior.
- Recovery tests confirm persistent consistency after restart.

## 8. Reliability Rules

- Prefer readiness checks/retries over fixed sleeps.
- Bound all external calls with timeout and retry caps.
- Keep rand strategy assertions statistical (distribution threshold), not single-case deterministic.
- Always archive actionable logs/artifacts for failed runs.

## 9. Resource Profile (Serial Mode)

Expected peak while single group runs:
- CPU: ~1-2 vCPU
- Memory: ~1.0-2.5 GB
- Temporary disk: ~200 MB-1 GB

Trade-off:
- Lower peak resource pressure than parallel model.
- Longer wall-clock runtime accepted by team.

## 10. Milestones

1. M1: Bring up self-hosted Keygen and export real license artifacts.
2. M2: Port `smoke` and `core` groups to real-license flow.
3. M3: Implement `security`, `recovery`, and `chaos` groups.
4. M4: Wire serial PR gate and artifact reporting.
5. M5: Stabilize flake rate and lock as required merge gate.

## 11. Out of Scope

- Using keygen.sh cloud API as mandatory PR dependency.
- Parallel group execution in PR pipeline (can be revisited later).

## 12. Next Step

Create an implementation plan that breaks this design into concrete tasks, sequencing, and verification checkpoints.
