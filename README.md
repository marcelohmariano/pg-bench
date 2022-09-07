# TimescaleDB Benchmark

A simple tool to benchmark TimescaleDB queries.

## Getting Started

These instructions will guide you on how to setup and run the project on your
local machine.

### Prerequisites

In order to run this project, you will need the following tools:

* [Bash](https://www.gnu.org/software/bash/) >= 3.2
* [GNU Make](https://www.gnu.org/software/make/) >= 3.81
* [Docker](https://docs.docker.com/get-docker/) >= 20.10
* [Docker Compose](https://docs.docker.com/compose/) >= 2.7.0

### Build and Run

Before performing the following steps, ensure that Docker is running so that the
`make` commands can work.

#### Setup the development environment

First, clone this repository:

```shell
git clone https://github.com/marcelohmariano/timescaledb-benchmark
```

Then, initialize the development environment:

```shell
cd timescaledb-benchmark
make init
```

Note that `make init` will build the `build-env` image with `UID` and `GID` from
the host to avoid permission problems.

#### Build

Once the development environment is ready, you can build the benchmark image:

```shell
docker build -t benchmark --target release .
```

Optionally, you can build a binary:

```shell
make all
```

The command above will produce a binary according to the host OS and CPU architecture
and will place it at `./bin/benchmark`. You can manually select an OS and CPU
architecture by setting the `GOOS` and `GOARCH` environment variables respectively.

#### Setup a TimescaleDB instance

Before running the benchmarks, setup a TimescaleDB instance set up with the
sample data provided [here](./data):

```shell
make seed
```

The command above will create a database named `homework` in a local TimescaleDB
instance configured according to the definitions in [docker-compose.yml](docker-compose.yml)
file.

#### Running

Now you can run the benchmark image:

```shell
docker run -it --rm --network host -v "$PWD/data:/data" benchmark \
  --db-url "postgres://postgres:pass@localhost:5432/homework" \
  --query ./query_template.sql \
  --query-args ./query_params.csv
```

Or, if you've chosen to build a binary in the previous steps, you can run
it like this:

```shell
./bin/benchmark \
  --db-url "postgres://postgres:pass@localhost:5432/homework" \
  --query ./data/query_template.sql \
  --query-args ./data/query_params.csv
```

In both cases, the ouput should be similar to this:

```shell
Queries processed: 200
Queries with success: 200
Queries with error: 0

Minimum query time: 8.679291ms
Maximum query time: 98.003167ms

Median query time: 8.990333ms
Average query time: 10.850226ms

Overall processing time: 573.229875ms
```
## Testing

You can run the tests by running the `make test` command.

## Linting

This project uses [golangci-lint](https://golangci-lint.run/) and [shellcheck](
https://github.com/koalaman/shellcheck) for linting Go and Bash sources
respectively. You can run the linters by running `make lint` from a terminal.

## Troubleshooting

### Cannot connect to the Docker daemon at unix:///var/run/docker.sock. Is the docker daemon running?

The message above indicates that Docker is not running. Start it and try what you
was doing again.

### Service "XXXXX" is not running container #1

The message above indicates that the development environment was not started.
Run `docker compose down && docker-compose up -d` or `make init` to start it.
