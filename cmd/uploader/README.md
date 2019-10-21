# Asset Upload Service

## Overview

The Asset Upload Service is an HTTP server used for hosting static files in Minio. It contains a simple HTTP endpoint which accepts `multipart/form-data` forms. It can upload files to the private and public system buckets so you can also use it for development purposes to host files, instead of relying on external providers.

## Prerequisites

Use these tools to set up the service:

- [Go](https://golang.org)
- [Docker](https://www.docker.com/)

## Usage

### Run a local version

To run the service against the local Kyma installation on Minikube without building the binary, use this command:

```bash
APP_KUBECONFIG_PATH=/Users/$USER/.kube/config APP_VERBOSE=true APP_UPLOAD_ACCESS_KEY={accessKey} APP_UPLOAD_SECRET_KEY={secretKey} go run main.go
```

Replace values in curly braces with proper details, where:

- `{accessKey}` is the access key required to sign in to the content storage server.
- `{secretKey}` is the secret key required to sign in to the content storage server.

The service listens on port `3000`.

### Access the service on a cluster

To run the Asset Upload Service on a cluster, use the command:

```bash
kubectl port-forward deployment/assetstore-asset-upload-service 3000:3000 -n kyma-system
```

You can access the service on port `3000`.

### Build a production version

To build the production Docker image, use this command:

```bash
docker build {image_name}:{image_tag}
```

The variables are:

- `{image_name}` that is the name of the output image. The default name is `asset-upload-service`.
- `{image_tag}` that is the tag of the output image. The default tag is `latest`.

### Upload files

For the full API documentation, including OpenAPI specification, see the [Asset Store docs](https://kyma-project.io/docs/master/components/asset-store#details-asset-upload-service).

### Environment variables

Use these environment variables to configure the service:

| Name | Required | Default | Description |
|------|----------|---------|-------------|
| **APP_BUCKET_PRIVATE_PREFIX** | No | `private` | The prefix of the private system bucket |
| **APP_BUCKET_PUBLIC_PREFIX** | No | `public` | The prefix of the public system bucket |
| **APP_BUCKET_REGION** | No | `us-east-1` | The region of system buckets |
| **APP_CONFIG_ENABLED** | No | `true` | The toggle used to save and load the configuration using the ConfigMap resource |
| **APP_CONFIG_NAME** | No | `asset-upload-service` | The name of the ConfigMap resource |
| **APP_CONFIG_NAMESPACE** | No | `kyma-system` | The Namespace in which the ConfigMap resource is created |
| **APP_HOST** | No | `127.0.0.1` | The host on which the HTTP server listens |
| **APP_KUBECONFIG_PATH** | No | None | The path to the kubeconfig file, needed to run the service outside of a cluster |
| **APP_MAX_UPLOAD_WORKERS** | No | `10` | The maximum number of concurrent upload workers |
| **APP_PORT** | No | `3000` | The port on which the HTTP server listens |
| **APP_UPLOAD_ACCESS_KEY** | Yes | None | The access key required to sign in to the content storage server |
| **APP_UPLOAD_ENDPOINT** | No | `minio.kyma.local` | The address of the content storage server |
| **APP_UPLOAD_EXTERNAL_ENDPOINT** | No | `https://minio.kyma.local` | The external address of the content storage server. If not set, the system uses the `APP_UPLOAD_ENDPOINT` variable. |
| **APP_UPLOAD_PORT** | No | `443` | The port on which the content storage server listens |
| **APP_UPLOAD_SECRET_KEY** | Yes | None | The secret key required to sign in to the content storage server |
| **APP_UPLOAD_SECURE** | No | `true` | The HTTPS connection with the content storage server |
| **APP_UPLOAD_TIMEOUT** | No | `30m` | The file upload timeout |
| **APP_VERBOSE** | No | None | The toggle used to enable detailed logs in the service |

### Configure the logger

This service uses `glog` to log messages. Pass command line arguments described in the [`glog.go`](https://github.com/golang/glog/blob/master/glog.go) file to customize the log parameters, such as the log level and output.

For example:
```bash
go run main.go --stderrthreshold=INFO -logtostderr=false
```

## Development

### Install dependencies

This project uses `dep` as a dependency manager. To install all required dependencies, use this command:

```bash
dep ensure -vendor-only
```

### Run tests

To run all unit tests, use this command:

```bash
go test ./...
```

### Verify the code

To check if the code is correct and you can push it, run the `before-commit.sh` script. It builds the service, runs tests, and checks the status of the vendored libraries. It also performs the static code analysis and ensures that code formatting is correct.
