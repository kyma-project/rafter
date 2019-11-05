# Rafter AsyncAPI service

This helm chart installs Rafter AsyncAPI service version v1.0.0 https://github.com/kyma-project/rafter/tree/v1.0.0

## TL;DR;

``` bash
$ helm install incubator/rafter-asyncapi-service
```

## Overview

This project contains the chart for the Rafter AsyncAPI service.

## Prerequisites

- Kubernetes 1.12+
- Helm 2.11+ or Helm 3.0-beta3+

## Installing the Chart

To install the chart with the release name `rafter-release`:

``` bash
$ helm install --name rafter-release incubator/rafter-asyncapi-service
```

The command deploys Rafter AsyncAPI service on the Kubernetes cluster in the default configuration. The [configuration](#configuration) section lists the parameters that can be configured during installation.

> **Tip**: List all releases using `helm list`.

## Uninstalling the Chart

To uninstall/delete the `rafter-release` deployment:

``` bash
$ helm delete rafter-release
```

The command removes all the Kubernetes components associated with the chart and deletes the release.

## Configuration

The following table lists the configurable parameters of the Rafter AsyncAPI service chart and their default values.

| Parameter | Description | Default |
| --- | ---| ---|
| `image.repository` | Rafter AsyncAPI service image repository. | `eu.gcr.io/kyma-project/rafter`  |
| `image.tag` | Rafter AsyncAPI service image tag. | `{TAG_NAME}` |
| `image.pullPolicy` | Rafter AsyncAPI service image pull policy. | `IfNotPresent` |
| `nameOverride` | String to partially override `rafterAsyncAPIService.name` template with a string (will prepend the release name). | `nil` |
| `fullnameOverride` | String to fully override `rafterAsyncAPIService.fullname` template with a string. | `nil` |
| `deployment.labels` | Custom labels for the `Deployment`. | `{}` |
| `deployment.annotations` | Custom annotations for the `Deployment`. | `{}` |
| `deployment.replicas` | Number of nodes. | `1` |
| `deployment.extraProperties` | Extra properties injected in the `Deployment`. | `{}` |
| `pod.labels` | Custom labels for the `Pod`. | `{}` |
| `pod.annotations` | Custom annotations for the `Pod`. | `{}` |
| `pod.extraProperties` | Extra properties injected in the `Pod`. | `{}` |
| `pod.extraContainerProperties` | Extra properties injected in the container. | `{}` |
| `service.name` | `Service` name. If not set a name is generated using the `rafterAsyncAPIService.fullname` template. | `nil` |
| `service.type` | `Service` type. | `ClusterIP` |
| `service.port.name` | Name of the port on the `Service`. | `http` |
| `service.port.external` | Port where the `Service` is exposed. | `80` |
| `service.port.internal` | Internal pod's port on the `Service`. | `3000` |
| `service.port.protocol` | Protocol of the port on the `Service`. | `TCP` |
| `service.annotations` | Custom annotations for the `Service`. | `{}` |
| `service.labels` | Custom labels for the `Service`. | `{}` |
| `serviceMonitor.create` | Whether a new `ServiceMonitor` resource that the Prometheus operator will use should be created. | `false` |
| `serviceMonitor.name` | `ServiceMonitor` resource to be used for the Prometheus operator. If not set and `serviceMonitor.create` is `true` a name is generated using the `rafterAsyncAPIService.fullname` template. If not set and `serviceMonitor.create` is `false` a name is `default`. | `nil` |
| `serviceMonitor.scrapeInterval` | Scrape interval for the custom `ServiceMonitor` resource. | `30s` |
| `serviceMonitor.labels` | Custom labels for the custom `ServiceMonitor` resource. | `{}` |
| `serviceMonitor.annotations` | Custom annotations for the custom `ServiceMonitor` resource. | `{}` |
| `envs.host` | App host. | `0.0.0.0` |
| `envs.verbose` | Whether a logs from app should be visible. | `true` |

Specify each parameter using the `--set key=value[,key=value]` argument to `helm install`. For example:

``` bash
$ helm install --name rafter-release \
  --set serviceMonitor.create=true,serviceMonitor.name="rafter-service-monitor" \
    incubator/rafter-asyncapi-service
```

The above command install release with custom `ServiceMonitor` resource with `rafter-service-monitor` name.

Alternatively, a YAML file that specifies the values for the above parameters can be provided while installing the chart. For example:

``` bash
$ helm install --name rafter-release -f values.yaml incubator/rafter-asyncapi-service
```

> **Tip**: You can use the default [values.yaml](./values.yaml).

### Templating values.yaml

The Rafter AsyncAPI service chart has possibility to templating `values.yaml`. This means that you can use, for example, `.Chart.*`, `.Values.*` or other defined by Helm variables. For example:

``` yaml
pod:
  annotations:
    sidecar.istio.io/inject: "false"
    recreate: "{{ .Release.Time.Seconds }}"
``` 
