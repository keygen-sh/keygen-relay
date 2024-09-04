# Keygen Relay

**Synopsis:** `relay` is a small command line utility that distributes license
files to nodes on a local network.

## Distribution

Relay will be distributed as a standalone cross-platform binary.

## Architecture

Relay will use a SQLite database for storage.

Relay will be written in Go.

Relay will have 2 parts: a vendor-facing CLI and an app-facing server. The CLI
can be used to manage the data served by the server.

Relay will be responsible for storing [License Files](https://keygen.sh/docs/api/licenses/#licenses-actions-check-out)
and distributing license files to nodes on a local network, offering a way to
implement a leasing model in an offline or air-gapped environment.

Alongide license files, Relay will also store license keys, which will be
distributed along side license files for decryption purposes.

There should be `license_files`, `license_keys`, and `nodes` tables. For
claiming, I'd like to use `FOR UPDATE SKIP LOCKED` to atomically move a license
from "unclaimed" to "claimed." This may necessitate a `claims` table or
similar, though I haven't fully fleshed this behavior out yet.

Applications will interact with Relay through the server, by requesting a
license during an application's boot lifecycle event, and returning it before
the shutdown lifecyle event. In the case of crashes, a heartbeat system coupled
with a TTL is used to reap zombie nodes, where the higher the TTL the longer
zombies may sit around.

For example, a vendor offers a node-based licensing model to their customers.
They have a customer that has purchased 250 nodes, but this customer runs an
air-gapped environment. The vendor can use Relay by storing 250 license files
and distributing them to nodes on the local network in a FIFO order.

In this example, Relay runs on the customer's local network, and is loaded by
the vendor with 250 license files and keys. Each time the vendor's application
boots, it requests a license file from Relay, and each time it shuts down, it
returns the license file.

This request-return lifecycle ensures that no more than 250 licenses are
"claimed" at any point in time.

## Security

All license files are signed with an account's private key, and so they cannot
be tampered with or forged. Any tampered, forged, or otherwise invalid license
files added via the CLI will be rejected by the application during signature
verification. All license files are unique.

As always, the application is responsible for verifying the license file's
signature, the license file's expiry, and the license's expiry.

All actions will be logged to an `audit_logs` table.

## CLI

Relay can be managed via the following CLI commands.

### `serve`

Runs a local relay server accessible at `--port`.

```sh
relay serve [--port=1337 --ttl=30s --lifo --fifo --rand]

relay serve             # serve on default port
relay serve --port 1337 # serve on custom port
```

### `add`

Push a license to the local relay server's pool.

```sh
relay add --file xxx.lic --key xxx [--file=yyy.lic --key=yyy]
```

Prints identifers of the added licenses.

This is an atomic operation.

### `del`

Delete a license from the local relay server's pool.

```sh
relay del --id xxx [--id=yyy]
```

Prints identifers of the deleted licenses.

This is an atomic operation.

### `ls`

Print the local relay server's license pool, with stats for each license.

```sh
relay ls
```

### `stat`

Print stats for a license in the local relay server's pool.

```sh
relay stat --id xxx [--id=yyy]
```

## API

Relay can be used via the following API endpoints.

### `PUT /v1/nodes/:node_id`

Claim a license from the relay server for a node, blocking others from claiming
it. The TTL for the claim respects the license file's TTL. This is an atomic
operation.

Accepts a `node_id`, where `node_id` is some fingerprint identifying the node.

Returns `200 OK` with a `license_file` and `license_key`. If no licenses are
available to be claimed, i.e. no licenses exist or all have been claimed,
the server will return `410 Gone`. If a claim already exists for the node,
a `409 Conflict` will be returned.

### `DELETE /v1/nodes/:node_id`

Release a claim on a license. This allows the license to be claimed and used by
other nodes. This is an atomic operation.

Accepts a `node_id`, where `node_id` is a fingerprint used for the claim.

Returns `204 No Content` with no content. If a claim does not exist for the
node, the server will return a `404 Not Found`.

### `PATCH /v1/nodes/:node_id`

Keep a claim on a license, extending the claim TTL past the default.

Accepts a `node_id`, where `node_id` is a fingerprint used for the claim.

Returns `204 No Content` with no content. If a claim does not exist for the
node, the server will return a `404 Not Found`.

## SDK

No immediate plans right now. For the time being, integrations will be similar
to those with Keygen's flagship licensing API.

Eventually, I'd like to offer a reference SDK in Go for Relay.

## Etc.

Relay will be backed by Keygen's [Go SDK](https://github.com/keygen-sh/keygen-go).
For example, when adding a license file, one could e.g. do:

```go
import (
  "github.com/keygen-sh/jsonapi-go"
  "github.com/keygen-sh/keygen-go/v3"
)

func AddLicense(licenseFile string, licenseKey string) {
  lic := &keygen.LicenseFile{Certificate: licenseFile}

  // crytographically verify the license file
  err = lic.Verify()
  switch {
  case err == keygen.ErrLicenseFileNotGenuine:
    panic("license file is not genuine!")
  case err != nil:
    panic(err)
  }

  // decrypt the license file
  dataset, err := lic.Decrypt(licenseKey)
  switch {
  case err == keygen.ErrSystemClockUnsynced:
    panic("system clock tampering detected!")
  case err == keygen.ErrLicenseFileExpired:
    panic("license file is expired!")
  case err != nil:
    panic(err)
  }

  // unmarshal the license
  license := keygen.License{}
  err = jsonapi.Unmarshal(dataset, license)
  if err != nil {
    panic(err)
  }

  // store everything
  save(
    license.ID,
    license.Key,
    lic,
  )
}
```
