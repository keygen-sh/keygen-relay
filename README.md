# Keygen Relay

A lightweight on-premise licensing server that relays data from Keygen's API to a local network.

Relay is run as a standalone binary.

## Commands

### Relay

Runs a local relay server that proxies data from Keygen's API to a local network. When internet
access is available, the relay server will periodically attempt to automatically sync with the
API unless the `--offline` flag is provided. When no internet connection is available, the
`sync` command should be run periodically.

```
relay serve --port=1337 [--offline]
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

Sync will perform a push, and then a pull. That is, it will send all queued requests to the
API, and then it will pull in all resources from the API and replace their current state
within the local datastore. If any of these steps fail, the failures will be logged but
the sync will continue.

The API is the authoritative source, i.e. the source of truth. If a conflict occurs, the
API's state wins and will overwrite any local data.

Data is stored in an encrypted SQLite database.
