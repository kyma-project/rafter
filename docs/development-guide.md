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
├── Docs                        # Rafter-related documentation
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

You can run the integration tests for Rafter with the same command both locally and on a cluster. To perform integration tests, copy the [`test-infra`](https://github.com/kyma-project/test-infra) repository under your `$GOPATH` workspace as `${GOPATH}/src/github.com/kyma-project/test-infra/` and run this command from the root of the `rafter` repository:

```bash
make integration-test
```
