# PostgreSQL Benchmark

A simple tool to benchmark PostgreSQL queries.

## Getting Started

These instructions will guide you on how to set up and run the project on your
local machine.

### Prerequisites

To run this project, you will need the following tools:

* [Go](https://go.dev/) >= 1.19
* [Bash](https://www.gnu.org/software/bash/) >= 5.2
* [GNU Make](https://www.gnu.org/software/make/) >= 3.81
* [Docker](https://docs.docker.com/get-docker/) >= 20.10

### Build and Run

#### Getting the sources

First, Clone this repository locally:

```shell
git clone https://github.com/marcelohmariano/pg-bench
```

#### Build

Build the benchmark binary:

```shell
cd pg-bench
make all
```

The command above will generate the binary at `./bin/pg-bench`.

#### Run

You can run the benchmark binary like so:

```shell
./bin/bench -u "<user>:<pass>@<host>:<port>/<dbname>" -f <sql_file>
```

After running it, you should see an output similar to this:

```sh
Statements:
  Total: 200
  Succeeded: 200
  Failed: 0

Durations:
  Min: 79.543205ms
  Max: 344.725547ms
  Average: 181.459413ms
  Median: 168.113747ms
  Overall: 3.133904956s
```

## Testing

You can run the tests by running the `make test` from a terminal.

## Linting

This project uses [golangci-lint](https://golangci-lint.run/) for linting Go
source files. You can run the linters by running `make lint` from a terminal.
