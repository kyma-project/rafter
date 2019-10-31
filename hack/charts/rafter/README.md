# Rafter

This helm chart installs Rafter version v1.0.0 https://github.com/kyma-project/rafter/tree/v1.0.0

## TL;DR;

``` bash
$ helm install incubator/rafter
```

## Overview

This project contains the chart for the Rafter Controller Manager.

## Prerequisites

- Kubernetes 1.15+
- Helm 2.15+ or Helm 3.0-beta3+

## Installing the Chart

To install the chart with the release name `rafter-release`:

``` bash
$ helm install --name rafter-release incubator/rafter
```

The command deploys Rafter on the Kubernetes cluster in the default configuration. The [parameters](#[arameters) section lists the parameters that can be configured during installation.

> **Tip**: List all releases using `helm list`.

## Uninstalling the Chart

To uninstall/delete the `rafter-release` deployment:

``` bash
$ helm delete rafter-release
```

The command removes all the Kubernetes components associated with the chart and deletes the release.

CRDs created by this chart are not removed by default and should be manually cleaned up:

``` bash
kubectl delete crd clusterassetgroups.rafter.kyma-project.io
kubectl delete crd assetgroups.rafter.kyma-project.io
kubectl delete crd clusterbuckets.rafter.kyma-project.io
kubectl delete crd buckets.rafter.kyma-project.io
kubectl delete crd clusterassets.rafter.kyma-project.io
kubectl delete crd assets.rafter.kyma-project.io
```

## Parameters

The following table lists the configurable parameters of the Rafter chart and their default values.

| Parameter | Description | Default |
| --- | ---| ---|
| `image.repository` | Rafter image repository. | `eu.gcr.io/kyma-project/rafter`  |
| `image.tag` | Rafter image tag. | `{TAG_NAME}` |
| `image.pullPolicy` | Rafter image pull policy. | `IfNotPresent` |
| `nameOverride` | String to partially override `rafter.name` template with a string (will prepend the release name). | `nil` |
| `fullnameOverride` | String to fully override `rafter.fullname` template with a string. | `nil` |
| `installCRDs` | If true, create CRDs managed by the Rafter. | `true` |
| `deployment.labels` | Custom labels for the `Deployment`. | `{}` |
| `deployment.annotations` | Custom annotations for the `Deployment`. | `{}` |
| `deployment.replicas` | Number of nodes. If value is great than 1, then controllers have enabled `leader-election`. | `1` |
| `deployment.extraProperties` | Extra properties injected in the `Deployment`. | `{}` |
| `pod.labels` | Custom labels for the `Pod`. | `{}` |
| `pod.annotations` | Custom annotations for the `Pod`. | `{}` |
| `pod.resources` | Pod's resource requests and limits. | `{}` |
| `pod.volumes` | Volumes for the `Pod`. | `{}` |
| `pod.volumeMounts` | Volume mounts for the container. | `{}` |
| `pod.extraProperties` | Extra properties injected in the `Pod`. | `{}` |
| `pod.extraContainerProperties` | Extra properties injected in the container. | `{}` |
| `serviceAccount.create` | Whether a new `ServiceAccount` resource that the Rafter will use should be created. | `true` |
| `serviceAccount.name` | `ServiceAccount` resource to be used for the Rafter. If not set and `serviceAccount.create` is `true` a name is generated using the `rafter.fullname` template. If not set and `serviceAccount.create` is `false` a name is `default`. | `nil` |
| `serviceAccount.labels` | Custom labels for the custom `ServiceAccount` resource. | `{}` |
| `serviceAccount.annotations` | Custom annotations for the custom `ServiceAccount` resource. | `{}` |
| `rbac.clusterScope.create` | Whether a new `ClusterRole` and `ClusterRoleBinding` resources that the Rafter will use should be created. | `true` |
| `rbac.clusterScope.role.name` | `ClusterRole` resource to be used for the Rafter. If not set and `rbac.clusterScope.create` is `true` a name is generated using the `rafter.fullname` template. If not set and `rbac.clusterScope.create` is `false` a name is `default`. | `nil` |
| `rbac.clusterScope.role.labels` | Custom labels for the custom `ClusterRole` resource. | `{}` |
| `rbac.clusterScope.role.annotations` | Custom annotations for the custom `ClusterRole` resource. | `{}` |
| `rbac.clusterScope.role.extraRules` | Extra rules for the custom `ClusterRole` resource. | `[]` |
| `rbac.clusterScope.roleBinding.name` | `ClusterRoleBinding` resource to be used for the Rafter. If not set and `rbac.clusterScope.create` is `true` a name is generated using the `rafter.fullname` template. If not set and `rbac.clusterScope.create` is `false` a name is `default`. | `nil` |
| `rbac.clusterScope.roleBinding.labels` | Custom labels for the custom `ClusterRoleBinding` resource. | `{}` |
| `rbac.clusterScope.roleBinding.annotations` | Custom annotations for the custom `ClusterRoleBinding` resource. | `{}` |
| `rbac.namespaced.create` | Whether a new `Role` and `RoleBinding` resources that the Rafter will use should be created. | `true` |
| `rbac.namespaced.role.name` | `Role` resource to be used for the Rafter. If not set and `rbac.namespaced.create` is `true` a name is generated using the `rafter.fullname` template. If not set and `rbac.namespaced.create` is `false` a name is `default`. | `nil` |
| `rbac.namespaced.role.labels` | Custom labels for the custom `Role` resource. | `{}` |
| `rbac.namespaced.role.annotations` | Custom annotations for the custom `Role` resource. | `{}` |
| `rbac.namespaced.role.extraRules` | Extra rules for the custom `Role` resource. | `[]` |
| `rbac.namespaced.roleBinding.name` | `RoleBinding` resource to be used for the Rafter. If not set and `rbac.namespaced.create` is `true` a name is generated using the `rafter.fullname` template. If not set and `rbac.namespaced.create` is `false` a name is `default`. | `nil` |
| `rbac.namespaced.roleBinding.labels` | Custom labels for the custom `RoleBinding` resource. | `{}` |
| `rbac.namespaced.roleBinding.annotations` | Custom annotations for the custom `RoleBinding` resource. | `{}` |
| `rbac.namespaced.roleBinding.name` | `RoleBinding` resource to be used for the Rafter. If not set and `rbac.namespaced.create` is `true` a name is generated using the `rafter.fullname` template. If not set and `rbac.namespaced.create` is `false` a name is `default`. | `nil` |
| `rbac.namespaced.roleBinding.labels` | Custom labels for the custom `RoleBinding` resource. | `{}` |
| `rbac.namespaced.roleBinding.annotations` | Custom annotations for the custom `RoleBinding` resource. | `{}` |
| `webhooksConfigMap.create` | Whether a new `ConfigMap` resource that the Rafter will use should be created. | `false` |
| `webhooksConfigMap.name` | `ConfigMap` resource to be used for the Rafter. If not set and `webhooksConfigMap.create` is `true` a name is generated using the `rafter.fullname` template. If not set and `webhooksConfigMap.create` is `false` a name is `default`. | `nil` |
| `webhooksConfigMap.hooks` | Data passed to custom `ConfigMap` resource. | `{}` |
| `webhooksConfigMap.labels` | Custom labels for the custom `ConfigMap` resource. | `{}` |
| `webhooksConfigMap.annotations` | Custom annotations for the custom `ConfigMap` resource. | `{}` |
| `metrics.enabled` | Set this to `true` to enable exporting the Prometheus monitoring metrics. | `true` |
| `metrics.service.name` | `Service` to be used for the exposing metrics. If not set and `metrics.enabled` is `true` a name is generated using the `rafter.fullname` template. If not set and `metrics.enabled` is `false` a name is `default`. | `nil` |
| `metrics.service.type` | `Service` type. | `ClusterIP` |
| `metrics.service.port.name` | Name of the port on the metrics `Service`. | `metrics` |
| `metrics.service.port.external` | Port where the metrics `Service` is exposed. | `8080` |
| `metrics.service.port.internal` | Internal pod's port on the metrics `Service`. | `metrics` |
| `metrics.service.port.protocol` | Protocol of the port on the metrics `Service`. | `TCP` |
| `metrics.service.annotations` | Custom annotations for the metrics `Service`. | `{}` |
| `metrics.service.labels` | Custom labels for the metrics `Service`. | `{}` |
| `metrics.serviceMonitor.create` | Whether a new `ServiceMonitor` resource that the Prometheus operator will use should be created. | `false` |
| `metrics.serviceMonitor.name` | `ServiceMonitor` resource to be used for the Prometheus operator. If not set and `metrics.serviceMonitor.create` is `true` a name is generated using the `rafter.fullname` template. If not set and `metrics.serviceMonitor.create` is `false` a name is `default`. | `nil` |
| `metrics.serviceMonitor.scrapeInterval` | Scrape interval for the custom `ServiceMonitor` resource. | `30s` |
| `metrics.serviceMonitor.labels` | Custom labels for the custom `ServiceMonitor` resource. | `{}` |
| `metrics.serviceMonitor.annotations` | Custom annotations for the custom `ServiceMonitor` resource. | `{}` |
| `metrics.pod.labels` | Custom labels for the Rafter `Pod`, when `metrics.enabled` is set to `true`. | `{}` |
| `metrics.pod.annotations` | Custom annotations for the Rafter `Pod`, when `metrics.enabled` is set to `true` | `{}` |
| `envs.clusterAssetGroup.relistInterval` | The period of time after which the controller refreshes the status of a `ClusterAssetGroup` CR. | `5m` |
| `envs.assetGroup.relistInterval` | The period of time after which the controller refreshes the status of an `AssetGroup` CR. | `5m` |
| `envs.clusterBucket.relistInterval` | The period of time after which the controller refreshes the status of a `ClusterBucket` CR. | `30s` |
| `envs.clusterBucket.maxConcurrentReconciles` | The maximum number of `ClusterBucket` reconciles that can run in parallel. | `1` |
| `envs.clusterBucket.region` | The location of the region in which the controller creates a `ClusterBucket` CR. If the field is empty, the controller creates the bucket under the default location. | `us-east-1` |
| `envs.bucket.relistInterval` | The period of time after which the controller refreshes the status of a `Bucket` CR. | `30s` |
| `envs.bucket.maxConcurrentReconciles` | The maximum number of `Bucket` reconciles that can run in parallel. | `1` |
| `envs.bucket.region` | The location of the region in which the controller creates a `Bucket` CR. If the field is empty, the controller creates the bucket under the default location. | `us-east-1` |
| `envs.clusterAsset.relistInterval` | The period of time after which the controller refreshes the status of a `ClusterAsset` CR. | `30s` |
| `envs.clusterAsset.maxConcurrentReconciles` | The maximum number of `ClusterAsset` reconciles that can run in parallel. | `1` |
| `envs.asset.relistInterval` | The period of time after which the controller refreshes the status of a `Asset` CR. | `30s` |
| `envs.asset.maxConcurrentReconciles` | The maximum number of `Asset` reconciles that can run in parallel. | `1` |
| `envs.store.endpoint` | The address of the content storage server. | `minio.kyma.local` |
| `envs.store.externalEndpoint` | The external address of the content storage server. If not set, the system uses the `APP_UPLOAD_ENDPOINT` variable. | `https://minio.kyma.local` |
| `envs.store.accessKey` | The access key required to sign in to the content storage server. | `nil` |
| `envs.store.secretKey` | The secret key required to sign in to the content storage server. | `nil` |
| `envs.store.useSSL` | The HTTPS connection with the content storage server. | `true` |
| `envs.store.uploadWorkers` | The number of workers used in parallel to upload files to the storage server. | `true` |
| `envs.loader.verifySSL` | The variable that verifies the SSL certificate before downloading source files. | `true` |
| `envs.loader.tempDir` | The path to the directory used to store data temporarily. | `/tmp` |
| `envs.webhooks.validation.timeout` | The period of time after which validation is canceled. | `1m` |
| `envs.webhooks.validation.workers` | The number of workers used in parallel to validate files. | `10` |
| `envs.webhooks.mutation.timeout` | The period of time after which mutation is canceled. | `1m` |
| `envs.webhooks.mutation.workers` | The number of workers used in parallel to mutate files. | `10` |
| `envs.webhooks.metadata.timeout` | The period of time after which metadata extraction is canceled. | `1m` |

Specify each parameter using the `--set key=value[,key=value]` argument to `helm install`. For example:

``` bash
$ helm install --name rafter-release \
  --set serviceAccount.create=true,serviceAccount.name="rafter-service-account" \
    incubator/rafter
```

The above command install release with custom `ServiceAccount` resource with `rafter-service-account` name.

Alternatively, a YAML file that specifies the values for the above parameters can be provided while installing the chart. For example:

``` bash
$ helm install --name rafter-release -f values.yaml incubator/rafter
```

> **Tip**: You can use the default [values.yaml](./values.yaml).

### Templating values.yaml

The Rafter chart has possibility to templating `values.yaml`. This means that you can use, for example, `.Chart.*`, `.Values.*` or other defined by Helm variables. For example:

``` yaml
pod:
  annotations:
    sidecar.istio.io/inject: "false"
    recreate: "{{ .Release.Time.Seconds }}"
``` 

### Change values for `envs.*` parameters

All `envs.*` parameters have possibility to define their values as object, so parameters can be provided as inline `value` or `valueFrom`. For example:

``` yaml
envs:
  clusterAssetGroup:
    relistInterval: 
      value: 5m
  store:
    endpoint: 
      valueFrom:
        configMapKeyRef:
          name: assetstore-minio-docs-upload
          key: APP_UPLOAD_ENDPOINT_WITH_PORT
    accessKey:
      valueFrom:
        secretKeyRef:
          name: assetstore-minio
          key: accesskey
```
