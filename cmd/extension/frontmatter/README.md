# Front Matter Service

## Overview

The Front Matter Service is an HTTP server used for extracting metadata from text files. It contains a simple HTTP endpoint which accepts `multipart/form-data` forms. The service extracts YAML front matter metadata from text files of all extensions.

## Prerequisites

Use these tools to set up the service:

- [Go](https://golang.org)
- [Docker](https://www.docker.com/)

## Usage

Read how to run and use the service.

### Run from sources

To run the service against the local installation on Minikube without building the binary, use this command:

```bash
APP_VERBOSE=true go run main.go
```

The service listens on port `3000`.

### Build a production version

To build the production Docker image, use this command:

```bash
docker build {image_name}:{image_tag}
```

The variables are:

- `{image_name}` that is the name of the output image. The default name is `asset-metadata-service`.
- `{image_tag}` that is the tag of the output image. The default tag is `latest`.

### Extract metadata from files

For the full API documentation, including OpenAPI specification, see the [Asset Store docs](https://kyma-project.io/docs/master/components/asset-store#details-asset-metadata-service).

### Environment variables

Use these environment variables to configure the service:

| Name | Required | Default | Description |
|------|:----------:|---------|-------------|
| **APP_HOST** | No | `127.0.0.1` | The host on which the HTTP server listens |
| **APP_MAX_WORKERS** | No | `10` | The maximum number of concurrent metadata extraction workers |
| **APP_PORT** | No | `3000` | The port on which the HTTP server listens |
| **APP_PROCESS_TIMEOUT** | No | `10m` | The file process timeout |
| **APP_VERBOSE** | No | None | The toggle used to enable detailed logs in the service |


### Configure the logger

This service uses `glog` to log messages. Pass command line arguments described in the [`glog.go`](https://github.com/golang/glog/blob/master/glog.go) file to customize the log parameters, such as the log level and output.

For example:

```bash
go run main.go --stderrthreshold=INFO -logtostderr=false
```

## Development

### Run tests

To run all unit tests, use this command:

```bash
go test ./...
```
