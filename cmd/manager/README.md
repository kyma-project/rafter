# Rafter Controller Manager

## Overview

The Rafter Controller Manager runs the AssetGroup, Asset, and Bucket Controllers that manage AssetGroup, Asset, and Bucket custom resources (CR).

## Prerequisites

Use these tools to set up the controller:

* [Go](https://golang.org)
* [Docker](https://www.docker.com/)
* [Kubebuilder](https://github.com/kubernetes-sigs/kubebuilder)

## Usage

Read how to run and use the controller manager.

### Run from sources

To run the controller manager from sources, use this command:

```bash
kubectl apply -k config/default \
    && APP_STORE_ACCESSKEY='{minio-accesskey}' APP_STORE_SECRETKEY='{minio-secretkey}' go run cmd/manager/main.go
```

### Build a production version

To build the production Docker image, use this command:

```bash
make build-manager
```

### Environment variables

Use these environment variables to configure the controller manager:

| Name | Required | Default | Description |
|------|----------|---------|-------------|
| **APP_CLUSTER_ASSET_GROUP_RELIST_INTERVAL** | No | `5m` | Period of time after which the controller refreshes the status of a ClusterAssetGroup CR. |
| **APP_ASSET_GROUP_RELIST_INTERVAL** | No | `5m` | Period of time after which the controller refreshes the status of an AssetGroup CR |
| **APP_CLUSTER_BUCKET_RELIST_INTERVAL** | No | `30s` | Period of time after which the controller refreshes the status of a ClusterBucket |
| **APP_CLUSTER_BUCKET_MAX_CONCURRENT_RECONCILES** | No | `1` | Maximum number of cluster bucket reconciles that can run in parallel |
| **APP_CLUSTER_BUCKET_REGION** | No | `us-east-1` | Location of the region in which the controller creates a ClusterBucket CR. If the field is empty, the controller creates the bucket under the default location. |
| **APP_BUCKET_RELIST_INTERVAL** | No | `30s` | Period of time after which the controller refreshes the status of a Bucket CR |
| **APP_BUCKET_MAX_CONCURRENT_RECONCILES** | No | `1` | Maximum number of bucket reconciles that can run in parallel |
| **APP_BUCKET_REGION** | No | `us-east-1` | Location of the region in which the controller creates a Bucket CR. If the field is empty, the controller creates the bucket under the default location. |
| **APP_CLUSTER_ASSET_RELIST_INTERVAL** | No | `30s` | Period of time after which the controller refreshes the status of a ClusterAsset CR |
| **APP_CLUSTER_ASSET_MAX_CONCURRENT_RECONCILES** | No | `1` | Maximum number of cluster asset reconciles that can run in parallel |
| **APP_ASSET_RELIST_INTERVAL** | No | `30s` | Period of time after which the controller refreshes the status of an Asset CR |
| **APP_ASSET_MAX_CONCURRENT_RECONCILES** | No | `1` | Maximum number of asset reconciles that can run in parallel |
| **APP_STORE_ENDPOINT** | No | `minio.kyma.local` | Address of the content storage server |
| **APP_STORE_EXTERNAL_ENDPOINT** | No | `https://minio.kyma.local` | External address of the content storage server |
| **APP_STORE_ACCESS_KEY** | Yes | None | Access key required to sign in to the content storage server |
| **APP_STORE_SECRET_KEY** | Yes | None | Secret key required to sign in to the content storage server |
| **APP_STORE_USE_SSL** | No | `true` | Variable that enforces the use of HTTPS for the connection with the content storage server |
| **APP_STORE_UPLOAD_WORKERS_COUNT** | No | `10` | Number of workers used in parallel to upload files to the storage bucket |
| **APP_LOADER_VERIFY_SSL** | No | `true` | Variable that verifies the SSL certificate before downloading source files |
| **APP_LOADER_TEMPORARY_DIRECTORY** | No | `/tmp` | Path to the directory used to store data temporarily |
| **APP_WEBHOOK_VALIDATION_TIMEOUT** | No | `1m` | Period of time after which validation is canceled |
| **APP_WEBHOOK_VALIDATION_WORKERS_COUNT** | No | `10` | Number of workers used in parallel to validate files |
| **APP_WEBHOOK_MUTATION_TIMEOUT** | No | `1m` | Period of time after which mutation is canceled |
| **APP_WEBHOOK_MUTATION_WORKERS_COUNT** | No | `10` | Number of workers used in parallel to mutate files |
| **APP_WEBHOOK_METADATA_EXTRACTION_TIMEOUT** | No | `1m` | Period of time after which metadata extraction is canceled |
| **APP_WEBHOOK_CONFIG_MAP_CFG_MAP_NAME** | No | `webhook-configmap` | Name of the ConfigMap that contains webhook definitions |
| **APP_WEBHOOK_CONFIG_MAP_CFG_MAP_NAMESPACE** | No | `kyma-system` | Namespace of the ConfigMap that contains webhook definitions |

## Development

There is a unified way of testing all changes in Rafter components. For details on how to run unit, integration, and MinIO Gateway tests, read [this](../../docs/development-guide.md) development guide.
