# Rafter

## Overview

Rafter is a Kubernetes-native solution for storing assets and managing content. The whole concept of Rafter relies on custom resources (CRs) such as AssetGroup CR, Bucket CR, or Asset CR. The custom resources can apply to a given Namespace or be cluster-wide. They are handled by the Rafter Controller Manager which consists of three seperate components:

- AssetGroup Controller
- Asset Controller
- Bucket Controller

Using Rafter, you can create AssetGroup CR for a particular content type, such as images, Markdown documents, AsyncAPI, OData, and OpenAPI specification files. Once the AssetsGroup Controller reads the AssetGroup CR definition, it creates a new Bucket CR and Asset CRs. Then, the controller monitors the status of the Asset CR and updates the status of the AssetsGroup CR accordingly.

## Prerequisites

Use the following tools to set up the project:

- [Go](https://golang.org)
- [Docker](https://www.docker.com/)
- [Kubebuilder](https://github.com/kubernetes-sigs/kubebuilder)

## Installation

Explain the steps to install your project. Include the instructions or provide links to the related documentation.

## Usage

Explain how to use the project. You can create multiple subsections. Include the instructions or provide links to the related documentation.

## Development

Add instructions on how to develop the project. It must be clear what to do and how to trigger the tests so that other contributors know how to make their pull requests acceptable. Include the instructions or provide links to related documentation.

## Project structure

The repository has the following structure:

```txt
├── .github                     # Pull request and issue templates
├── cmd                         # Rafter's applications
├── config                      # Configuration file templates
├── deploy                      # Dockerfiles for Rafter's applications
├── hack                        # Information useful for development
├── internal                    # Private application and library code
├── pkg                         # Library code ready-to-use by external applications
└── tests                       # Integration tests
```
