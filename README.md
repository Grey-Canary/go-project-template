# Go Project Template

This is the sole mono-repo containing a raw template for creating API, worker, and scheduled command services.

## Requirements

* [Golang](https://go.dev/doc/install) >= 1.23
* [Docker](https://docs.docker.com/engine/install/)

## Quick Start: Running Docker Locally

This project contains a Makefile that summarizes most of the frequently used build commands. To use it simply run:
```bash
make up
```

Under the hood this will generate an `.env` file from the `.env.local` and run `docker compose -f docker-compose.local.yaml up -d` which runs docker in the background in detached mode.

To stop all containers run:
```bash
make down
```

This will stop all containers and additionally delete the api container for cleanup.

## Quick Start: Running the CLI

You can install the project CLI and run it manually using the build command:
```bash
go build -o ./cmd/project
```

This will create a binary at the root of the project which can take commands like `./project api` to run the API.
Alternatively the project CLI can be directly installed on your machine using `go install ...` instead.

## Manual Testing

Tests can be run individually or for the whole repository. To run the full suite of tests using make:
```bash
make test
```

Which runs the equivalent of `go test -v -race -cover -count=1 -failfast ./...`

## Documentation

Any infrastructure, code, or API documentation shall live in the `/docs` directory. This includes images, media, or text files that are relative to the explanation of this project and it's environment.

API Documentation can be manually generated using the command:
```bash
make swaggo
```

This uses the [Swaggo/Swag](https://github.com/swaggo/swag) library to parse the API Handler and domain model annotations. A Go, JSON, and YAML file in OpenAPI v2 spec will be upserted into the `/docs` directory.