# Rafter Upload service

This helm chart installs Rafter Upload service version v1.0.0 https://github.com/kyma-project/rafter/tree/v1.0.0

## TL;DR;

``` bash
$ helm install incubator/rafter-upload-service
```

## Overview

This project contains the chart for the Rafter Upload service.

## Prerequisites

- Kubernetes 1.12+
- Helm 2.11+ or Helm 3.0-beta3+

## Installing the Chart

To install the chart with the release name `rafter-release`:

``` bash
$ helm install --name rafter-release incubator/rafter-upload-service
```

The command deploys Rafter Upload service on the Kubernetes cluster in the default configuration. The [configuration](#configuration) section lists the parameters that can be configured during installation.

## Uninstalling the Chart

To uninstall/delete the `rafter-release` deployment:

``` bash
$ helm delete rafter-release
```

The command removes all the Kubernetes components associated with the chart and deletes the release.

## Configuration

The following table lists the configurable parameters of the Rafter Upload service chart and their default values.

| Parameter | Description | Default |
| --- | ---| ---|
| `image.repository` | Rafter Upload service image repository. | `eu.gcr.io/kyma-project/rafter`  |
| `image.tag` | Rafter Upload service image tag. | `{TAG_NAME}` |
| `image.pullPolicy` | Rafter Upload service image pull policy. | `IfNotPresent` |
| `nameOverride` | String to partially override `rafterUploadService.name` template with a string (will prepend the release name). | `nil` |
| `fullnameOverride` | String to fully override `rafterUploadService.fullname` template with a string. | `nil` |
| `deployment.labels` | Custom labels for the `Deployment`. | `{}` |
| `deployment.annotations` | Custom annotations for the `Deployment`. | `{}` |
| `deployment.replicas` | Number of nodes. | `1` |
| `deployment.extraProperties` | Extra properties injected in the `Deployment`. | `{}` |
| `pod.labels` | Custom labels for the `Pod`. | `{}` |
| `pod.annotations` | Custom annotations for the `Pod`. | `{}` |
| `pod.extraProperties` | Extra properties injected in the `Pod`. | `{}` |
| `pod.extraContainerProperties` | Extra properties injected in the container. | `{}` |
| `service.name` | `Service` name. If not set a name is generated using the `rafterUploadService.fullname` template. | `nil` |
| `service.type` | `Service` type. | `ClusterIP` |
| `service.verbose` | Whether a logs from `Service` should be visible. | `true` |
| `service.host` | `Service` host. | `0.0.0.0` |
| `service.port.name` | Name of the port on the `Service`. | `http` |
| `service.port.external` | Port where the `Service` is exposed. | `80` |
| `service.port.internal` | Internal pod's port on the `Service`. | `3000` |
| `service.port.protocol` | Protocol of the port on the `Service`. | `TCP` |
| `service.annotations` | Custom annotations for the `Service`. | `{}` |
| `service.labels` | Custom labels for the `Service`. | `{}` |
| `serviceAccount.create` | Whether a new `ServiceAccount` resource that the Rafter Upload service will use should be created. | `true` |
| `serviceAccount.name` | `ServiceAccount` resource to be used for the Rafter Upload service. If not set and `serviceAccount.create` is `true` a name is generated using the `rafterUploadService.fullname` template. If not set and `serviceAccount.create` is `false` a name is `default`. | `nil` |
| `serviceAccount.labels` | Custom labels for the custom `ServiceAccount` resource. | `{}` |
| `serviceAccount.annotations` | Custom annotations for the custom `ServiceAccount` resource. | `{}` |
| `rbac.clusterScope.create` | Whether a new `ClusterRole` and `ClusterRoleBinding` resources that the Rafter will use should be created. | `true` |
| `rbac.clusterScope.role.name` | `ClusterRole` resource to be used for the Rafter. If not set and `rbac.clusterScope.create` is `true` a name is generated using the `rafterUploadService.fullname` template. If not set and `rbac.clusterScope.create` is `false` a name is `default`. | `nil` |
| `rbac.clusterScope.role.labels` | Custom labels for the custom `ClusterRole` resource. | `{}` |
| `rbac.clusterScope.role.annotations` | Custom annotations for the custom `ClusterRole` resource. | `{}` |
| `rbac.clusterScope.role.extraRules` | Extra rules for the custom `ClusterRole` resource. | `[]` |
| `rbac.clusterScope.roleBinding.name` | `ClusterRoleBinding` resource to be used for the Rafter. If not set and `rbac.clusterScope.create` is `true` a name is generated using the `rafterUploadService.fullname` template. If not set and `rbac.clusterScope.create` is `false` a name is `default`. | `nil` |
| `rbac.clusterScope.roleBinding.labels` | Custom labels for the custom `ClusterRoleBinding` resource. | `{}` |
| `rbac.clusterScope.roleBinding.annotations` | Custom annotations for the custom `ClusterRoleBinding` resource. | `{}` |
| `serviceMonitor.create` | Whether a new `ServiceMonitor` resource that the Prometheus operator will use should be created. | `false` |
| `serviceMonitor.name` | `ServiceMonitor` resource to be used for the Prometheus operator. If not set and `serviceMonitor.create` is `true` a name is generated using the `rafterUploadService.fullname` template. If not set and `serviceMonitor.create` is `false` a name is `default`. | `nil` |
| `serviceMonitor.scrapeInterval` | Scrape interval for the custom `ServiceMonitor` resource. | `30s` |
| `serviceMonitor.labels` | Custom labels for the custom `ServiceMonitor` resource. | `{}` |
| `serviceMonitor.annotations` | Custom annotations for the custom `ServiceMonitor` resource. | `{}` |

Specify each parameter using the `--set key=value[,key=value]` argument to `helm install`. For example:

``` bash
$ helm install --name rafter-release \
  --set serviceMonitor.create=true,serviceMonitor.name="rafter-service-monitor" \
    incubator/rafter-upload-service
```

The above command install release with custom `ServiceMonitor` resource with `rafter-service-monitor` name.

Alternatively, a YAML file that specifies the values for the above parameters can be provided while installing the chart. For example:

``` bash
$ helm install --name rafter-release -f values.yaml incubator/rafter-upload-service
```

> **Tip**: You can use the default [values.yaml](./values.yaml)
