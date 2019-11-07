# AsyncAPI Service

## Overview

The AsyncAPI Service is an HTTP server used to process AsyncAPI specifications. It contains the `/validate` and `/convert` HTTP endpoints which accept `multipart/form-data` forms:
- The `/validate` endpoint validates the AsyncAPI specification against the AsyncAPI schema in version 2.0.0., using the [AsyncAPI Parser](https://github.com/asyncapi/parser).
- The `/convert` endpoint converts the version and format of the AsyncAPI files.

This service uses the [AsyncAPI Converter](https://github.com/asyncapi/converter-go) to change the AsyncAPI specifications from older versions to version 2.0.0, and convert any YAML input files to the JSON format that is required to render the specifications in the Console UI.

## Prerequisites

Use the following tools to set up the service:

- [Go](https://golang.org)
- [Docker](https://www.docker.com/)

## Usage

Read how to run and use the service.

### API

See the [OpenAPI specification](openapi.yaml) for the full API documentation. You can use the [Swagger Editor](https://editor.swagger.io/) to preview and test the API service.

### Run from sources

To run the local version of the AsyncAPI Service without building the binary, use this command:

```bash
go run cmd/extension/asyncapi/main.go
```

The service listens on port `3000`.

### Build a production version

To build the production Docker image, use this command:

```bash
make build-asyncapi
```

### Environment variables

Use these environment variables to configure the service:

| Name | Required | Default | Description |
|------|----------|---------|-------------|
| **APP_SERVICE_HOST** | No | `127.0.0.1` | The host on which the HTTP server listens |
| **APP_SERVICE_PORT** | No | `3000` | The port on which the HTTP server listens |
| **APP_VERBOSE** | No | `false` | The toggle used to enable detailed logs in the service |
