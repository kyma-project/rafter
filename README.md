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

>**NOTE:** If you want to learn how Rafter is implemented in [Kyma](https://kyma-project.io), read the [Asset Store](https://kyma-project.io/docs/components/asset-store/) and [Headless CMS](https://kyma-project.io/docs/components/headless-cms/) documentation.

## Quick start

Try out [this](https://katacoda.com/rafter/scenarios/rafter) set of interactive tutorials to see Rafter in action on Minikube. These tutorials show how to:
- Quickly install Rafter with our Helm Chart.
- Host a simple static site.
- Use Rafter as headless CMS with support of Rafter metadata webhook and Front Matter service . This example is based on a use case of storing Markdown files.
- Use Rafter as headless CMS with the support of Rafter validation and conversion webhooks. This example is based on a use case of storing [AsyncAPI](https://asyncapi.org/) specifications.

>**NOTE:** Read [this](./docs/development-guide.md) development guide to start developing the project.
