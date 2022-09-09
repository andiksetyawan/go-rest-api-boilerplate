# go-rest-api-boilerplate

#### TODO:
- [X] Clean Arch
- [X] Unit Test
- [X] Metrics, tracer, logger integration
- [ ] Integration Test
- [ ] http req/res metrics
- [ ] Docker-compose services: open-telemetry collector, uptrace, dll
- [ ] Doc Readme
- [ ] Openapi spec
- [ ] Increase test coverage
- [ ] .etc

### Table of contents
- [Installation](#installation)
- [Usage](#usage)
- [Documentation](#documentation)

## Installation

### Requirements

1. [Go](https://golang.org/doc/install) 1.16+
2. [Docker](https://docs.docker.com/engine/install/)

### Setting up environment
Default env:
```
SERVICE_NAME=go-rest-api-boilerplate
SERVICE_ENVIRONMENT=production
SERVICE_ADDRESS=:8080
DB_HOST=127.0.0.1
DB_PORT=5436
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=svc-go-rest-api-boilerplate
OTEL_UPTRACE_DSN=
```

you can setup the environment config using .env file or environment variables (OS). Set SERVICE ENVIRONMENT=production if you want sent metrics, trackers, logs to opentelemetry-collector

## Getting Started
## Usage
### Development
```
go mod tidy
```
After all installed properly, you can Build binary file and then migrate the database:
#### Build binary
```
make build
```
### Run database migration:
```
go run cmd/main.go migrate up
```
or with binary
```
./go-rest-api-boilerplate migrate up
```
### Run api server:
```
go run cmd/main.go server
```
or with binary
```
./go-rest-api-boilerplate server
```
### Run Api server using docker container:
Run api server service, database service, migrator with docker-compose:

```
docker-compose up
```

### Make requests:
HTTP/1.1 POST API with curl
```
```

reponse:
```
{"error":false,"message":"OK"}
```

### Make test:
```
go test ./...
```
or
```
make test
```

## Documentation
### Api specs:
openapi:
[```api/v1/api-specs.json```](./api/v1/api-specs.json)


