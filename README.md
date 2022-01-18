# Keygen Relay

A lightweight on-premise licensing server that relays data from Keygen's API to a local network.

Relay is run as a standalone binary.

## Commands

### Relay

Runs a local relay server that proxies data from Keygen's API to a local network. When internet
access is available, the relay server can automatically sync with the API using the `--sync` flag.
When no internet connection is available, the `sync` command should be run periodically.

```
relay serve --port=1337 [--sync=<interval>]
```

### Sync

Syncs data from the local network to Keygen's API. This can be run periodically in air-gapped
environments to ensure licensing state is up-to-date. Requires internet access to run.

```
relay sync
```

## Architecture

Relay uses a write-ahead-log to keep a list of events that need to be synced to Keygen's API.
The `sync` command will perform the logged events in-order, while continuing to accept new
events. Any errors that occur during a sync will be logged to stderr. You may choose to pipe
the app's output into a log file for easier reference.

Data is stored in an encrypted SQLite database.
