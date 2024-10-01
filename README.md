
# Keygen Relay

Relay is an offline-first on-premise licensing server backed by [Keygen](https://keygen.sh).
Use Relay to securely manage distribution of cryptographically signed and
encrypted license files across nodes in offline or air-gapped environments.

Relay has a vendor-facing CLI that can be used to onboard a customer's air-gap
environment. An admin can initialize Relay with N licenses to be distributed
across M nodes, ensuring that only N nodes are licensed at one time.

Relay has an app-facing REST API which nodes can communicate with to claim and
release licenses.

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

To install Relay, you can run the following command. Alternatively, you can
install manually by downloading a binary and following [the install
instructions here](https://keygen.sh/docs/relay/).

```bash
curl -sSL https://get.keygen.sh/keygen/relay/latest/install.sh | sh
```

## Usage

### CLI

The CLI can be used by the vendor to setup and manage customer environments.

#### Add license

You can add a new license to the system using the `add` command:

```bash
relay add --file license.lic --key xxx --public-key xxx
```

The `add` command supports the following flags:

| Flag                 | Description                                      |
|:---------------------|:-------------------------------------------------|
| `--file`             | Path to the license file to be added.            |
| `--key`              | License key for decryption.                      |
| `--public-key`       | Path to the public key for license verification. |

#### Delete license

To delete a license, use the `del` command:

```bash
relay del --id xxx
```

The `del` command supports the following flags:

| Flag   | Description                                      |
|:-------|:-------------------------------------------------|
| `--id` | The unique ID of the license to retrieve delete. |

#### List licenses

To list all the licenses in the database, use the `ls` command:

```bash
relay ls
```

#### Stat license

To retrieve the status of a specific license, use the `stat` command:

```bash
relay stat --id xxx
```

The `stat` command supports the following flags:

| Flag   | Description                                          |
|:-------|:-----------------------------------------------------|
| `--id` | The unique ID of the license to retrieve info about. |

### Server

To start the relay server, use the following command:

```bash
relay serve --port 8080
```

The `serve` command supports the following flags:

| Flag                 | Description                                                                                                                                                | Default          |
|:---------------------|:-----------------------------------------------------------------------------------------------------------------------------------------------------------|:-----------------|
| `--port`, `-p`       | Specifies the port on which the relay server will run.                                                                                                     | `8080`           |
| `--no-heartbeats`    | Disables the heartbeat mechanism. When this flag is enabled, the server will not track node activity.                                                      | `false`          |
| `--strategy`         | Specifies the license assignment strategy. Options: `fifo`, `lifo`, `rand`.                                                                                | `fifo`           |
| `--ttl`, `-t`        | Sets the Time-to-Live (TTL) for license claims. After this period, licenses will be automatically released. Accepts durations like `30s`, `1m`, `1h`, etc. | `30s`            |
| `--cleanup-interval` | Specifies how often the server should check for inactive nodes to clean up.                                                                                | `15s`            |
| `--database`         | Specify a custom database file for storing the license and node data.                                                                                      | `./relay.sqlite` |

E.g. to start the server on port `8080`, with a 30 second node TTL and FIFO
distribution strategy:

```bash
./bin/relay serve --port 8080 --ttl 30s --strategy fifo
```

### API

The API can be consumed by the vendor's application to claim and release a
license on behalf of a node.

#### Claim license

Nodes can claim a license by sending a `PUT` request to the
`/v1/nodes/{fingerprint}` endpoint:

```bash
curl -v -X PUT "http://localhost:8080/v1/nodes/$(cat /etc/machine-id)"
```

Accepts a `fingerprint`, an arbitrary string identifying the node.

Returns `201 Created` with a `license_file` and `license_key` for new nodes. If
a claim already exists for the node, the claim is extended by `--ttl` and the
server will return `202 Accepted`, unless heartbeats are disabled and in that
case a `409 Conflict` will be returned. If no licenses are available to be
claimed, i.e. no licenses exist or all have been claimed, the server will
return `410 Gone`.

#### Release license

Nodes can release a license by sending a `DELETE` request to the same endpoint:

```bash
curl -v -X DELETE "http://localhost:8080/v1/nodes/$(cat /etc/machine-id)"
```

Accepts a `fingerprint`, the node fingerprint used for the claim.

Returns `204 No Content` with no content. If a claim does not exist for the
node, the server will return a `404 Not Found`.

## Developing

### Building

To build the Keygen Relay from the source, clone this repository and run:

```bash
go build -o relay ./cmd/relay
```

Alternatively, you can build binaries for different platforms and architectures
using the provided `Makefile`:

```bash
make build # or make build-all
```

### Testing

Keygen Relay comes with a suite of tests, including integration tests that
verify behavior with the real server.

To run regular tests:

```bash
make test
```

To run integration tests (which are tagged with `// +build integrity`):

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
