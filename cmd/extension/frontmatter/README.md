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
APP_VERBOSE=true go run cmd/extension/frontmatter/main.go
```

The service listens on port `3000`.

### Build a production version

To build the production Docker image, use this command:

```bash
make build-frontmatter
```

### Extract metadata from files

For the full API documentation, including OpenAPI specification, see the [Asset Store docs](https://kyma-project.io/docs/master/components/asset-store#details-asset-metadata-service).

### Environment variables

Use these environment variables to configure the service:

| Name | Required | Default | Description |
|------|:----------:|---------|-------------|
| **APP_PORT** | No | `3000` | Port on which the HTTP server listens |
| **APP_HOST** | No | `127.0.0.1` | Host on which the HTTP server listens |
| **APP_VERBOSE** | No | `false` | Toggle used to enable detailed logs in the service |
| **APP_PROCESS_TIMEOUT** | No | `10m` | File process timeout |
| **APP_MAX_WORKERS** | No | `10` | Maximum number of concurrent metadata extraction workers |


### Configure the logger

This service uses `glog` to log messages. Pass command line arguments described in the [`glog.go`](https://github.com/golang/glog/blob/master/glog.go) file to customize the log parameters, such as the log level and output.

For example:

```bash
go run cmd/extension/frontmatter/main.go --stderrthreshold=INFO -logtostderr=false
```

## Development

There is a unified way of testing all changes in Rafter components. For details on how to run unit, integration, and MinIO Gateway tests, read [this](../../../docs/development-guide.md) development guide.
