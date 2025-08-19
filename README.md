# Keygen Relay

[![CI](https://github.com/keygen-sh/keygen-relay/actions/workflows/test.yml/badge.svg)](https://github.com/keygen-sh/keygen-relay/actions)

Relay is an offline-first on-premise licensing server backed by [Keygen](https://keygen.sh).
Use Relay to securely manage distribution of cryptographically signed and
encrypted license files across nodes in an offline or air-gapped environment.
Relay does not require or utilize an internet connection — it is meant to be
used stand-alone in an offline or air-gapped network.

Relay has a vendor-facing CLI that can be used to onboard a customer's air-gap
environment. An admin can initialize Relay with N licenses to be distributed
across M nodes, ensuring that only N nodes are licensed at one time.

Relay provides an app-facing REST API that allows nodes to claim a lease on a
license and release it when no longer needed.

In low-trust environments, Relay itself can be [node-locked](#node-locking).

## Background

Relay was born out of a limitation in Keygen — and really, a limitation in all
licensing APIs — the limitation being that implementing a node-based licensing
model, e.g. floating licenses, is hard in an air-gapped or otherwise offline
environment using an external API, due to the nature of APIs needing an
internet connection. Whether self-hosting Keygen EE, or using Keygen Cloud,
the issue remains the same.

Since Keygen is an API, it can't communicate to the nodes inside these isolated
environments, and that means it can't easily track which nodes are being used
and which nodes are not. It also has no visibility into how many nodes there
are currently vs how many nodes are allowed in total. Some vendors may be able
to whitelist Keygen in the customer's firewall, but that's rare.

In the past, we've seen workarounds for this problem. Most of them consist of
using an intermediary between the offline world and the online world —
typically a mobile device or a tablet. In this case, the intermediary acts on
behalf of the offline node, activating it via an online portal, and passing
on a signed payload, e.g. a license file, for verification.

As an alternative, some customers have even asked if they can self-host Keygen
on-premise for customers — but that's inherently unsafe, since customers would
have full access to Keygen, thus access to granting themselves licenses,
adjusting policy rules, etc.

While the aforementioned intermediary-based workaround can work — it's brittle.
And it requires human intervention, which just doesn't really work in the age
of cloud computing and autoscaling. For example, you couldn't use this
workaround to license on-premise software, where you wanted to only allow the
customer to use 20 concurrent processes at one time — it just wouldn't be
feasible to ask a human to hop on their phone and activate nodes in an
autoscaling k8s cluster as it autoscales.

Thus, the idea for Relay was born — a bridge between Keygen and the offline
universe, secured via cryptography.

## Installation

To install Relay, you can follow the instructions and run the command below.
Alternatively, you can install manually by downloading a prebuilt binary and
following [the install instructions here](https://keygen.sh/docs/relay/).

Automatically detect and install `relay` on the current platform:

```bash
curl -sSL https://raw.pkg.keygen.sh/keygen/relay/latest/install.sh | sh
```

This will install `relay` in `/usr/local/bin`.

Missing a platform? Open an [issue](https://github.com/keygen-sh/keygen-relay/issues).

## Usage

For all available commands and flags, run `relay --help`.

### CLI

The CLI can be used by the vendor to setup and manage customer environments.

#### Add license

You can add a new license to the pool using the `add` command:

```bash
relay add --file license.lic --key xxx --public-key xxx
```

The `add` command supports the following flags:

| Flag                 | Description                                                                                                 |
|:---------------------|:------------------------------------------------------------------------------------------------------------|
| `--file`             | Path to the license file to add to the pool.                                                                |
| `--key`              | License key for decryption.                                                                                 |
| `--public-key`       | Your account's public key for license file verification. (Not available when [node-locked](#node-locking).) |
| `--pool`             | Add the license to a specific named pool.                                                                   |

The `add` command supports multiple `--file` and `--key` pairs.

#### Delete license

To delete a license from the pool, use the `del` command:

```bash
relay del --license xxx
```

The `del` command supports the following flags:

| Flag        | Description                                           |
|:------------|:------------------------------------------------------|
| `--license` | The unique ID of the license to delete from the pool. |

The `del` command supports multiple `--license` flags.

#### List licenses

To list all the licenses in the pool, use the `ls` command:

```bash
relay ls
```

The `ls` command supports the following flags:

| Flag      | Description                                   |
|:----------|:----------------------------------------------|
| `--plain` | Print results non-interactively in plaintext. |
| `--pool`  | Print licenses from a specific pool.          |

#### Stat license

To retrieve the status of a specific license, use the `stat` command:

```bash
relay stat --license xxx
```

The `stat` command supports the following flags:

| Flag        | Description                                          |
|:------------|:-----------------------------------------------------|
| `--license` | The unique ID of the license to retrieve info about. |
| `--plain`   | Print results non-interactively in plaintext.        |

### Server

To start the relay server, use the following command:

```bash
relay serve --port 6349
```

The `serve` command supports the following flags:

| Flag                 | Description                                                                                                                                                                           | Default          |
|:---------------------|:--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|:-----------------|
| `--port`, `-p`       | Specifies the port on which the relay server will run.                                                                                                                                | `6349`           |
| `--no-heartbeats`    | Disables the heartbeat system. When this flag is enabled, the server will not automatically release inactive or dead nodes, and leases cannot be extended.                            | `false`          |
| `--strategy`         | Specifies the license distribution strategy. Options: `fifo`, `lifo`, `rand`.                                                                                                         | `fifo`           |
| `--ttl`, `-t`        | Sets the time-to-live for leases. Licenses will be automatically released after the time-to-live if a node heartbeat is not maintained. Options: e.g. `30s`, `1m`, `1h`, etc.         | `60s`            |
| `--cull-interval`    | Specifies how often the server should check for and deactivate inactive or dead nodes.                                                                                                | `15s`            |
| `--database`         | Specify a custom database file for storing the license and node data.                                                                                                                 | `./relay.sqlite` |
| `--pool`             | Specify a specific pool to serve licenses from.                                                                                                                                       |                  |

E.g. to start the server on port `8080`, with a 30 second node TTL and FIFO
distribution strategy:

```bash
relay serve --port 8080 --ttl 30s --strategy fifo
```

### API

The API can be consumed by the vendor's application to claim a lease on a
license, and also release the lease, on behalf of a node.

#### Health check

The Relay server's health can be checked with the following endpoint:

```bash
curl -v -X GET "http://localhost:6349/v1/health"
```

Returns a `200 OK` status code.

#### Claim license

Nodes can claim a lease on a license by sending a `PUT` request to the
`/v1/nodes/{fingerprint}` endpoint:

```bash
curl -v -X PUT "http://localhost:6349/v1/nodes/$(cat /etc/machine-id)"
```

Accepts a `fingerprint`, an arbitrary string identifying the node.

Returns `201 Created` with a `license_file` and `license_key` for new nodes. If
a lease already exists for the node, the lease is extended by `--ttl` and the
server will return `202 Accepted`, unless heartbeats are disabled and in that
case a `409 Conflict` will be returned. If no licenses are available to be
leased, i.e. no licenses exist or all are being actively leased, the server
will return `410 Gone`.

```json
{
  "license_file": "LS0tLS1CRUdJTiBMSUNFTlNFIEZJTEUtLS0tL...S0NCg0K",
  "license_key": "9A96B8-FD08CD-8C433B-7657C8-8A8655-V3"
}
```

The `license_file` will be base64 encoded.

#### Release license

Nodes can release a license when no longer needed by sending a `DELETE` request
to the same endpoint:

```bash
curl -v -X DELETE "http://localhost:6349/v1/nodes/$(cat /etc/machine-id)"
```

Accepts a `fingerprint`, the node fingerprint used for the lease.

Returns `204 No Content` with no content. If a lease does not exist for the
node, the server will return a `404 Not Found`.

## Pools

Relay supports a concept called "pools," where, via the `--pool` flag, licenses
can be added into distinct license pools, effectively allowing you to run
multiple Relays under a single Relay instance. Pools can be used to separate
licenses between environments, products, etc.

Alternatively, the `serve` command can also be configured to serve from a
specific pool:

```bash
relay serve --pool "prod" -vvvv
```

Consumers of Relay can specify a pool using the `Relay-Pool` header:

```bash
fingerprint=$(cat /etc/machine-id | openssl dgst -sha256 -hmac "prod" -binary | xxd -p -c 256)

curl -v -X PUT -H "Relay-Pool: prod" "http://localhost:6349/v1/nodes/$fingerprint"
```

> [!WARNING]
> As demonstrated above, we recommend using an HMAC on node fingerprints, to
> prevent node collisions, especially in situations where a node could be used
> across pools.
>
> Cross-pool node collisions will result in a `409 Conflict`.

If Relay is serving from a specific pool, i.e. via the `--pool` flag, the
`Relay-Pool` header MUST match the configured pool, or it can be omitted to
default to the configured pool. If Relay is configured to serve from the `prod`
pool and a request comes in for the `dev` pool, a `400 Bad Request` will be
returned.

Otherwise, if Relay is serving from all pools, the `Relay-Pool` header can be
provided to interact with a specific pool, or omitted to consume from the
global pool.

## Logs

Relay comes equipped with audit logs out-of-the-box, allowing the full history
of the Relay server to be audited. They can be viewed using a `sqlite3` client,
providing the path to the Relay database file:

```bash
sqlite3 ./relay.sqlite
```

```bash
.mode box --wrap 60 --wordwrap on --noquote
.headers on
```

```sql
-- recent events
SELECT
  audit_logs.*
FROM
  audit_logs
ORDER BY
  audit_logs.created_at DESC
LIMIT
  25;

-- recently leased licenses
SELECT
  audit_logs.*
FROM
  audit_logs
INNER JOIN
  event_types ON event_types.id = audit_logs.event_type_id
WHERE
  event_types.name = 'license.leased'
ORDER BY
  audit_logs.created_at DESC
LIMIT
  5;

-- entire history in chronological order
SELECT
  datetime(audit_logs.created_at, 'unixepoch') AS created_at,
  event_types.name AS event_type,
  entity_types.name AS entity_type,
  audit_logs.entity_id,
  audit_logs.pool_id,
  pools.name AS pool_name
FROM
  audit_logs
INNER JOIN
  event_types ON event_types.id = audit_logs.event_type_id
INNER JOIN
  entity_types ON entity_types.id = audit_logs.entity_type_id
LEFT OUTER JOIN
  pools ON pools.id = audit_logs.pool_id
ORDER BY
  audit_logs.created_at ASC;
```

If you have concerns about storage, or do not wish to keep audit logs, use
Relay's `--no-audit` flag to disable them.

## Building

To build Keygen Relay from source, clone this repository and run:

```bash
go build -o relay ./cmd/relay

# or...
make build
```

Alternatively, you can build binaries for specific platforms and architectures
using the provided `make` commands:

```bash
make build

# or specific platform...
make build-linux-amd64

# or all platforms...
make build-all
```

### Node-locking

Node-locking is a powerful way to secure Relay in low-trust environments by
tying it to a specific machine. To enable this, you'll need to build Relay
with a few extra flags:

```bash
# Your account's Ed25519 public key for verifying license files (required)
export BUILD_NODE_LOCKED_PUBLIC_KEY='e8601e48b69383ba520245fd07971e983d06d22c4257cfd82304601479cee788'

# Machine fingerprint to lock to (required)
export BUILD_NODE_LOCKED_FINGERPRINT='364646b45d9b732f1baaeee9382f4e7e541e65e8a9fd4aa72e4853477d85bf08'

# Platform identifier (optional)
export BUILD_NODE_LOCKED_PLATFORM='linux/amd64'

# Hostname (optional)
export BUILD_NODE_LOCKED_HOSTNAME='relay'

# Local IP address (optional)
export BUILD_NODE_LOCKED_IP='192.168.1.1'

# Relay bind addr (optional)
export BUILD_NODE_LOCKED_ADDR='0.0.0.0'

# Relay port (optional)
export BUILD_NODE_LOCKED_PORT='6349'

# Build the node-locked binary using the above constraints
BUILD_NODE_LOCKED=1 make build-linux-amd64
```

Relay fingerprints machines using [`keygen-sh/machineid`](https://github.com/keygen-sh/machineid),
which calculates the HMAC-SHA256 of the app ID, `keygen-relay`, keyed by the
underlying machine ID.

To determine a machine’s fingerprint before building, use the `machineid` CLI:

```bash
go install github.com/keygen-sh/machineid/cmd/machineid

machineid --appid keygen-relay
# => 364646b45d9b732f1baaeee9382f4e7e541e65e8a9fd4aa72e4853477d85bf08
```

(If you need an alternate fingerprinting method, open an [issue](https://github.com/keygen-sh/keygen-relay/issues)
or a [pull request](https://github.com/keygen-sh/keygen-relay/pulls)!)

Once built, you can [checkout](https://keygen.sh/docs/api/machines/#machines-actions-check-out)
and distribute an encrypted machine file along with a license key to the
end-user. Relay can then validate the machine before running, ensuring
it matches the locked configuration.

To provide the machine file and license key, use the `--node-locked-machine-file-path`
and `--node-locked-license-key` flags, respectively:

```bash
relay serve -vvvv --node-locked-machine-file-path /etc/keygen/relay.lic \
  --node-locked-license-key '73F7DA-19BCBF-30B806-2F4C7D-3C2ACE-V3'
```

Relay will verify and decrypt the machine file, then check that the system
matches the expected values. This makes it easy to enforce machine-level
security and prevent unauthorized use.

## Developing

### Building

To build Relay for the current platform, run the following `make` command:

```bash
make build
```

### Releasing

To cut and publish a new release of Relay, update the `VERSION` file and run
the following `make` command:

```
make release
```

Releases are uploaded and published using the [Keygen CLI](https://keygen.sh/docs/cli/).
Releases are hosted and distributed by [Keygen Cloud](https://keygen.sh). You
will need credentials and permission to upload to our production Keygen Cloud
account.

### Testing

Keygen Relay comes with a suite of tests, including integration tests that
verify behavior with the real server.

To run regular tests:

```bash
make test
```

To run integration tests, tagged with `// +build integration`:

```bash
make test-integration
```

## License

This project is licensed under the MIT License. See the [LICENSE](https://github.com/keygen-sh/keygen-relay/blob/master/LICENSE) file for details.

## Contributing

If you discover an issue, or are interested in a new feature, please open an
issue. If you want to contribute code, feel free to open a pull request. If the
PR is substantial, it may be beneficial to open an issue beforehand to discuss.

The CLA is available [here](https://keygen.sh/cla/).

## Security

We take security at Keygen very seriously. If you believe you've found a
vulnerability, please see our [`SECURITY.md`](https://github.com/keygen-sh/keygen-relay/blob/master/SECURITY.md)
file.

## Is it any good?

[Yes.](https://news.ycombinator.com/item?id=3067434)
