
# Keygen Relay

Keygen Relay is a server that helps manage the distribution of software licenses across nodes.
It provides an API for claiming and releasing licenses, and comes with features such as license expiration (TTL), license claiming strategies, and heartbeats for tracking active nodes.

## Features
- License distribution and management
- Heartbeat mechanism for license renewal
- Configurable TTL (Time-to-Live) for licenses
- Support for different distribution strategies: FIFO, LIFO, and Random
- REST API for interaction

## Installation

### Build from source

To build the Keygen Relay from the source, clone this repository and run:

```bash
go build -o relay ./cmd/relay
```

Alternatively, you can build binaries for different platforms and architectures using the provided `Makefile`:

```bash
make build # or make build-all
```

## Usage

To start the relay server, use the following command:

```bash
./bin/relay serve --port 8080
```

### Available Flags

`server` command comes with several flags that allow you to customize its behavior

| Flag         | Description                                                                                                                                       | Default        |
|--------------|---------------------------------------------------------------------------------------------------------------------------------------------------|----------------|
| `--port`, `-p` | Specifies the port on which the relay server will run.                                                                                            | `8080`         |
| `--no-heartbeats` | Disables the heartbeat mechanism. When this flag is enabled, the server will not track node activity.                                             | `false`        |
| `--strategy` | Specifies the license assignment strategy. Options: `fifo`, `lifo`, `rand`.                                                                       | `fifo`         |
| `--ttl`, `-t` | Sets the Time-to-Live (TTL) for license claims. After this period, licenses will be automatically released. Accepts values like `5s`, `30s`, `1m` |    `30s`
| `--cleanup-interval` | Specifies how often the server should check for inactive nodes to clean up.                                                                       | `15s`          |
| `--database` | Specify a custom database file for storing the license and node data.                                                                             | `relay.sqlite` |

### Example Usage
#### Start the server

To start the server on port 8080 with a TTL of 30 seconds for licenses and FIFO license assignment:

```bash
./bin/relay serve --port 8080 --ttl 30s --strategy fifo
```

#### Claim a license
Nodes can claim a license by sending a PUT request to the `/v1/nodes/{fingerprint}` endpoint:

```bash
curl -X PUT "http://localhost:8080/v1/nodes/$(cat /etc/machine-id)"
```

#### Release a license
Nodes can release a license by sending a DELETE request to the same endpoint:

```bash
curl -X DELETE "http://localhost:8080/v1/nodes/$(cat /etc/machine-id)"
```

## CLI Commands

### Add a License
You can add a new license to the system using the add command:

| Flag                 | Description                                                                                                                                      |
|----------------------|--------------------------------------------------------------------------------------------------------------------------------------------------|
| `--file`             | Path to the license file to be added.                                                                                                            | 
| `--key`              | License key for decryption                                                                                                                       | 
| `--public-key`       | Path to the public key for license verification                                                                                                  |


```bash
./bin/relay add --file license.lic --key xxx --public-key xxx
```

### List licenses
To list all the licenses in the database, use the ls command:

```bash
./bin/relay ls
```

### License status
To retrieve the status of a specific license, use the stat command:

```bash
./bin/relay stat --id xxx
```

| Flag   | Description                                          |
|--------|------------------------------------------------------|
| `--id` | The unique ID of the license to retrieve info about  |

### Delete a license
To delete a license, use the del command:

```bash
./bin/relay del --id xxx
```

## Development
### Running Tests

Keygen Relay comes with a suite of tests, including integration tests that verify behavior with the real server.

- To run regular tests:

  ```bash
  make test
  ```

- To run integration tests (which are tagged with `// +build integrity`):

  ```bash
  make test-integration
  ```

## License

This project is licensed under the MIT License. See the [LICENSE](https://github.com/keygen-sh/keygen-relay/blob/master/LICENSE) file for details.
