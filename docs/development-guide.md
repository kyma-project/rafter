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
├── hack                        # Information, scripts, and files useful for development
├── internal                    # Private application and library code
├── pkg                         # Library code to be used by external applications
└── tests                       # Integration tests
```

## Usage

After you make changes to a given Rafter component, build it to see if it works as expected. The commands you must run differ depending on the application you develop.

Follow these links for details:

- [AsyncAPI Service](./cmd/extension/asyncapi#usage)
- [Front Matter Service](./cmd/extension/frontmatter#usage)
- [Rafter Controller Manager](./cmd/manager/README.md#usage)
- [Upload Service](./cmd/uploader#usage)

## Unit tests

To perform unit tests, run this command from the root of the `rafter` repository:

```bash
make test
```

## Integration tests

_TODO_
