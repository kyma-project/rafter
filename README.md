# Rafter

[![Go Report Card](https://goreportcard.com/badge/github.com/kyma-project/rafter)](https://goreportcard.com/report/github.com/kyma-project/rafter)
[![Slack](https://img.shields.io/badge/slack-%23rafter%20channel-yellow)](http://slack.kyma-project.io)

<p align="center">
  <img src="rafter.png" alt="rafter" width="300" />
</p>

---

:warning: **Warning** :warning:

**Rafter is looking for new maintainers**

The project will no longer be developed within the `kyma-project` organization.
Contact us if you are interested in becoming a new maintainer.
If we fail to find new maintainers, the project will be archived.
Until then, no new features will be developed and maintenance activities will be limited to the bare minimum.

---

## Overview

Rafter is a solution for storing and managing different types of files called assets. It uses [MinIO](https://min.io/) as object storage. The whole concept of Rafter relies on [Kubernetes custom resources (CRs)](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/) managed by the [Rafter Controller Manager](./cmd/manager/README.md). These CRs include:

- Asset CR which manages single assets or asset packages from URLs or ConfigMaps
- Bucket CR which manages buckets
- AssetGroup CR which manages a group of Asset CRs of a specific type to make it easier to use and extract webhook information

Rafter enables you to manage assets using supported webhooks. For example, if you use Rafter to store a file such as a specification, you can additionally define a webhook service that Rafter should call before the file is sent to storage. The webhook service can:

- Validate the file
- Mutate the file
- Extract some of the file information and put it in the status of the custom resource

Rafter comes with the following set of services and extensions compatible with Rafter webhooks:

- [Upload Service](./cmd/uploader/README.md) (optional service)
- [AsyncAPI Service](./cmd/extension/asyncapi/README.md) (extension)
- [Front Matter Service](./cmd/extension/frontmatter/README.md) (extension)

> **NOTE:** To learn how Rafter is implemented in [Kyma](https://kyma-project.io), read [Rafter](https://kyma-project.io/docs/master/components/rafter) documentation.

### What Rafter is not

- Rafter is not a [Content Management System](https://en.wikipedia.org/wiki/Content_management_system) (Wordpress-like),
- Rafter is not a solution for [Enterprise Content Management](https://en.wikipedia.org/wiki/Enterprise_content_management),
- Rafter doesn't come with any out-of-the-box UI that allows you to modify or consume files managed by Rafter.

### What Rafter can be used for

- Rafter is based on [CRs](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/). Therefore, it is an extension of Kubernetes API and should be used mainly by developers building their solutions on top of Kubernetes,
- Rafter is a file store that allows you to programmatically modify, validate the files and/or extract their metadata before they go to storage. Content of those files can be fetched using an API. This is a basic functionality of the [headless CMS](https://en.wikipedia.org/wiki/Headless_content_management_system) concept. If you want to deploy an application to Kubernetes and enrich it with additional documentation or specifications, you can do it using Rafter,
- Rafter is an S3-like file store also for files written in HTML, CSS, and JS. It means that Rafter can be used as a hosting solution for client-side applications.

## Quick start

Try out [this](https://katacoda.com/rafter/) set of interactive tutorials to see Rafter in action on Minikube. These tutorials show how to:

- Quickly install Rafter with our Helm Chart.
- Host a simple static site.
- Use Rafter as [headless CMS](https://en.wikipedia.org/wiki/Headless_content_management_system) with the support of Rafter metadata webhook and Front Matter service. This example is based on a use case of storing Markdown files.
- Use Rafter as [headless CMS](https://en.wikipedia.org/wiki/Headless_content_management_system) with the support of Rafter validation and conversion webhooks. This example is based on a use case of storing [AsyncAPI](https://asyncapi.org/) specifications.

> **NOTE:** Read [this](./docs/development-guide.md) development guide to start developing the project.

## Installation

### Prerequisites

- Kubernetes 1.14 or higher / Minikube 1.3 or higher
- Helm 2.16.0 or higher

### Steps

1. Add a new chart's repository to Helm. Run:

   `helm repo add rafter-charts https://kyma-project.github.io/rafter`

2. Install Rafter:

   `helm install --name rafter --set rafter-controller-manager.minio.service.type=NodePort rafter-charts/rafter`
