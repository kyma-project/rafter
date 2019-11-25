# Rafter

[![Go Report Card](https://goreportcard.com/badge/github.com/kyma-project/rafter)](https://goreportcard.com/report/github.com/kyma-project/rafter)
[![Slack](https://img.shields.io/badge/slack-%23rafter%20channel-yellow)](http://slack.kyma-project.io)

<p align="center">
  <img src="rafter.png" alt="rafter" width="300" />
</p>

## Overview

Rafter is a solution for storing and managing different types of files called assets. It uses [MinIO](https://min.io/) as object storage. The whole concept of Rafter relies on Kubernetes custom resources (CRs) managed by the [Rafter Controller Manager](./cmd/manager/README.md). These CRs include:

- Asset CR which manages a single asset or a package of assets
- Bucket CR which manages buckets
- AssetGroup CR which manages a group of Asset CRs of a specific type to make it easier to use and extract webhook information

Rafter enables you to manage assets using supported webhooks. For example, if you use Rafter to store a file such as a specification, you can additionally define a webhook service that Rafer should call before the file is sent to storage. The webhook service can:

- validate the file
- mutate the file
- extract some of the file information and put it in the status of the custom resource

Rafter comes with the following set of services and extensions compatible with Rafter webhooks:

- [Upload Service](./cmd/uploader/README.md) (optional service)
- [AsyncAPI Service](./cmd/extension/asyncapi/README.md) (extension)
- [Front Matter Service](./cmd/extension/frontmatter/README.md) (extension)

To see the implementation of Rafter in [Kyma](https://kyma-project.io), follow these links:

- [Asset Store](https://kyma-project.io/docs/components/asset-store/)
- [Headless CMS](https://kyma-project.io/docs/components/headless-cms/)

## Quick start

Try out [this](https://katacoda.com/rafter/scenarios/rafter) set of interactive tutorials to see Rafter in action on Minikube. These tutorials show how to:
- Quickly install Rafter with our Helm Chart
- Host a simple static site
- Use Rafter as headless CMS with support of Rafter metadata webhook and Front Matter service . This example is based on a use case of storing Markdown files
- Use Rafter as headless CMS with support of Rafter validation and conversion webhooks. This example is based on a use case of storing [AsyncAPI](https://asyncapi.org/) specification

## Project structure

The repository has the following structure:

```txt
├── .github                     # Pull request and issue templates
├── cmd                         # Rafter's applications
├── config                      # Configuration file templates
├── deploy                      # Dockerfiles for Rafter's applications
├── hack                        # Information, scripts, and files useful for development
├── internal                    # Private application and library code
├── pkg                         # Library code to be used by external applications
└── tests                       # Integration tests
```
