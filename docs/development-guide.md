# Development guide

This short guide describes the repository structure, and explains how you can run and test changes you make while developing specific Rafter components.

## Project structure

The project's main component is the Rafter Controller Manager that embraces all controllers handling custom resources in Rafter. It also contains services that mutate, validate, and extract metadata from assets. The source code of these applications is located under the `cmd` folder.

The whole structure of the repository looks as follows:

```txt
├── .github                     # Pull request and issue templates
├── charts                      # Configuration of component charts
├── cmd                         # Rafter's applications
├── config                      # Configuration file templates
├── deploy                      # Dockerfiles for Rafter's applications
├── docs                        # Rafter-related documentation
├── hack                        # Information, scripts, and files useful for development
├── internal                    # Private application and library code
├── pkg                         # Library code to be used by external applications
└── tests                       # Integration tests
```

## Usage

After you make changes to a given Rafter component, build it to see if it works as expected. The commands you must run differ depending on the application you develop.

Follow these links for details:

- [AsyncAPI Service](../cmd/extension/asyncapi#usage)
- [Front Matter Service](../cmd/extension/frontmatter#usage)
- [Rafter Controller Manager](../cmd/manager/README.md#usage)
- [Upload Service](../cmd/uploader#usage)

## Unit tests

>**NOTE:** Install [Go](https://golang.org) before you run unit tests.

To perform unit tests, run this command from the root of the `rafter` repository:

```bash
make test
```

## Integration tests

>**NOTE:** Install [Go](https://golang.org) and [Docker](https://www.docker.com/) before you run integration tests.

You can run the integration tests against a cluster on which Rafter is installed. To perform the tests, copy the [`test-infra`](https://github.com/kyma-project/test-infra) repository under your `$GOPATH` workspace as `${GOPATH}/src/github.com/kyma-project/test-infra/` and run this command from the root of the `rafter` repository:

```bash
make integration-test
```

## MinIO Gateway tests

> **NOTE:** Install [Go](https://golang.org), [Docker](https://www.docker.com/), [gsutil](https://cloud.google.com/storage/docs/gsutil), [azure-cli](https://docs.microsoft.com/en-us/cli/azure/install-azure-cli?view=azure-cli-latest) before you run MinIO Gateway tests.

There are two types of tests you can run for Rafter to check if the Gateway mode works as expected:

- `minio-gateway-test` tests Rafter that already has MinIO in the Gateway mode.
- `minio-gateway-migration-test` tests if Rafter runs as expected by switching it from MinIO in stand-alone mode to MinIO Gateway.

MinIO Gateway tests can run on [`GCS`](https://cloud.google.com/storage/) and [`Azure`](https://azure.microsoft.com/en-us/) platforms. See the [**MinIO Gateway environments**](#minio-gateway-environments) section to know which environment variables you must define for a given platform before you run MinIO Gateway tests.

You can run MinIO Gateway tests against a cluster on which Rafter is installed. Before you start, copy the [`test-infra`](https://github.com/kyma-project/test-infra) repository under your `$GOPATH` workspace as `${GOPATH}/src/github.com/kyma-project/test-infra/`. Run one of these commands from the root of the `rafter` repository:

- for MinIO already in the Gateway mode

```bash
make minio-gateway-test
```

- to switch from MinIO in stand-alone mode to MinIO Gateway

```bash
make minio-gateway-migration-test
```

### MinIO Gateway environments

See the required GCP variables:

| Variable | Description |
| --- | --- |
| **MINIO_GATEWAY_MODE** | Platform for MinIO Gateway tests |
| **CLOUDSDK_CORE_PROJECT** | Name of the [Google Cloud Platform (GCP)](https://cloud.google.com/) project for all GCP resources used in the tests |
| **GOOGLE_APPLICATION_CREDENTIALS** | Absolute path to the [Google Cloud Platform (GCP)](https://cloud.google.com/) Service Account Key File with the **Storage Admin** role |

>**NOTE:** **MINIO_GATEWAY_MODE** must be set to `gcs`.

See the required Azure variables:

| Variable | Description |
| --- | --- |
| **MINIO_GATEWAY_MODE** | Platform for MinIO Gateway tests |
| **BUILD_TYPE** | Defines one of `pr/master/release`. This value is used to create the name of the Azure Storage Account. |
| **PULL_NUMBER** | Defines the pull request number. Required if **BUILD_TYPE** is set to `pr`. |
| **AZURE_RS_GROUP** | Defines the name of the Azure Resource Group |
| **AZURE_REGION** | Azure region code |
| **AZURE_SUBSCRIPTION_ID** | ID of the the Azure Subscription |
| **AZURE_CREDENTIALS_FILE** | Path to the credentials for the Azure Subscription |

>**NOTE:** **MINIO_GATEWAY_MODE** must be set to `azure`.
