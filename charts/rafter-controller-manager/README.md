# Rafter Controller Manager

This project contains the Helm chart for the Rafter Controller Manager.

## Prerequisites

- Kubernetes v1.14 or higher
- Helm v2.15 or higher

## Details

Read how to install, uninstall, and configure the chart.

### Install the chart

Use this command to install the chart:

``` bash
helm install incubator/rafter-controller-manager
```

To install the chart with the `rafter-release` release name, use:

``` bash
helm install --name rafter-release incubator/rafter-controller-manager
```

The command deploys the Rafter Controller Manager on the Kubernetes cluster with the default configuration. The [**Configuration**](#configuration) section lists the parameters that you can configure during installation.

> **TIP:** To list all releases, use `helm list`.

### Uninstall the chart

To uninstall the `rafter-release` release, run:

``` bash
helm delete rafter-release
```

That command removes all the Kubernetes components associated with the chart and deletes the release.

### Configuration

The following table lists the configurable parameters of the Rafter Controller Manager chart and their default values.

| Parameter | Description | Default |
| --- | ---| ---|
| `image.repository` | Rafter Controller Manager image repository | `eu.gcr.io/kyma-project/rafter-controller-manager` |
| `image.tag` | Rafter Controller Manager image tag | `{TAG_NAME}` |
| `image.pullPolicy` | Pull policy for the Rafter Controller Manager image | `IfNotPresent` |
| `nameOverride` | String that partially overrides the `rafter.name` template | `nil` |
| `fullnameOverride` | String that fully overrides the `rafter.fullname` template | `nil` |
| `deployment.labels` | Custom labels for the Deployment | `{}` |
| `deployment.annotations` | Custom annotations for the Deployment | `{}` |
| `deployment.replicas` | Number of Rafter Controller Manager nodes | `1` |
| `deployment.extraProperties` | Additional properties injected in the Deployment | `{}` |
| `pod.labels` | Custom labels for the Pod | `{}` |
| `pod.annotations` | Custom annotations for the Pod | `{}` |
| `pod.extraProperties` | Additional properties injected in the Pod | `{}` |
| `pod.extraContainerProperties` | Additional properties injected in the container | `{}` |
| `serviceAccount.create` | Parameter that defines whether to create a new ServiceAccount for the Rafter Controller Manager | `true` |
| `serviceAccount.name` | ServiceAccount resource that the Rafter Controller Manager uses. If not set and the `serviceAccount.create` parameter is set to `true`, the name is generated using the `rafter.fullname` template. If not set and `serviceAccount.create` is set to `false`, the name is set to `default`. | `nil` |
| `serviceAccount.labels` | Custom labels for the ServiceAccount | `{}` |
| `serviceAccount.annotations` | Custom annotations for the ServiceAccount | `{}` |
| `rbac.clusterScope.create` | Parameter that defines whether to create a new ClusterRole and ClusterRoleBinding for the Rafter Controller Manager | `true` |
| `rbac.clusterScope.role.name` | ClusterRole resource that the Rafter Controller Manager uses. If not set and the `rbac.clusterScope.create` parameter is set to `true`, the name is generated using the `rafter.fullname` template. If not set and `rbac.clusterScope.create` is set to `false`, the name is set to `default`. | `nil` |
| `rbac.clusterScope.role.labels` | Custom labels for the ClusterRole | `{}` |
| `rbac.clusterScope.role.annotations` | Custom annotations for the ClusterRole | `{}` |
| `rbac.clusterScope.role.extraRules` | Additional rules injected in the ClusterRole | `[]` |
| `rbac.clusterScope.roleBinding.name` | ClusterRoleBinding resource that the Rafter Controller Manager uses. If not set and the `rbac.clusterScope.create` parameter is set to `true`, the name is generated using the `rafter.fullname` template. If not set and `rbac.clusterScope.create` is set to `false`, the name is set to `default`. | `nil` |
| `rbac.clusterScope.roleBinding.labels` | Custom labels for the ClusterRoleBinding | `{}` |
| `rbac.clusterScope.roleBinding.annotations` | Custom annotations for the ClusterRoleBinding | `{}` |
| `rbac.namespaced.create` | Parameter that defines whether to create a new Role and RoleBinding for the Rafter Controller Manager | `true` |
| `rbac.namespaced.role.name` | Role resource that the Rafter Controller Manager uses. If not set and the `rbac.namespaced.create` parameter is set to `true`, the name is generated using the `rafter.fullname` template. If not set and `rbac.namespaced.create` is set to `false`, the name is set to `default`.  | `nil` |
| `rbac.namespaced.role.labels` | Custom annotations for the Role | `{}` |
| `rbac.namespaced.role.annotations` | Custom annotations for the Role | `{}` |
| `rbac.namespaced.role.extraRules` | Additional rules injected in the Role | `[]` |
| `rbac.namespaced.roleBinding.name` | RoleBinding resource that the Rafter Controller Manager uses. If not set and the `rbac.namespaced.create` parameter is set to `true`, the name is generated using the `rafter.fullname` template. If not set and `rbac.namespaced.create` is set to `false`, the name is set to `default`. | `nil` |
| `rbac.namespaced.roleBinding.labels` | Custom annotations for the RoleBinding | `{}` |
| `rbac.namespaced.roleBinding.annotations` | Custom annotations for the RoleBinding | `{}` |
| `webhooksConfigMap.create` | Parameter that defines whether to create a new ConfigMap with the Webhooks data for the Rafter Controller Manager | `false` |
| `webhooksConfigMap.name` | ConfigMap resource that the Rafter Controller Manager uses. If not set and the `webhooksConfigMap.create` parameter is set to `true`, the name is generated using the `rafter.fullname` template. If not set and `webhooksConfigMap.create` is set to `false`, the name is set to `default`. | `nil` |
| `webhooksConfigMap.hooks` | Data passed to ConfigMap | `{}` |
| `webhooksConfigMap.labels` | Custom labels for the ConfigMap | `{}` |
| `webhooksConfigMap.annotations` | Custom annotations for the ConfigMap | `{}` |
| `metrics.enabled` | Parameter that defines whether to enable exporting the Prometheus monitoring metrics | `true` |
| `metrics.service.name` | Name of Service to be used for the exposing metrics. If not set and `metrics.enabled` is `true` a name is generated using the `rafter.fullname` template. If not set and `metrics.enabled` is `false` a name is `default`. | `nil` |
| `metrics.service.type` |  Service type | `ClusterIP` |
| `metrics.service.port.name` | Name of the port on the metrics `Service`. | `metrics` |
| `metrics.service.port.internal` | Internal port of the Service in the Pod | `metrics` |
| `metrics.service.port.external` | Port on which the Service is exposed in Kubernetes | `8080` |
| `metrics.service.port.protocol` | Protocol of the Service port | `TCP` |
| `metrics.service.labels` | Custom labels for the Service | `{}` |
| `metrics.service.annotations` | Custom annotations for the Service | `{}` |
| `metrics.serviceMonitor.create` | Parameter that defines whether to create a new ServiceMonitor custom resource for the Prometheus Operator | `false` |
| `metrics.serviceMonitor.name` | ServiceMonitor resource that the Prometheus Operator uses. If not set and the `serviceMonitor.create` parameter is set to `true`, the name is generated using the `rafter.fullname` template. If not set and `serviceMonitor.create` is set to `false`, the name is set to `default`. | `nil` |
| `metrics.serviceMonitor.scrapeInterval` | Scrape interval for the ServiceMonitor custom resource | `30s` |
| `metrics.serviceMonitor.labels` | Custom labels for the ServiceMonitor custom resource | `{}` |
| `metrics.serviceMonitor.annotations` | Custom annotations for the ServiceMonitor custom resource | `{}` |
| `metrics.pod.labels` | Custom labels for the Pod, when `metrics.enabled` is set to `true` | `{}` |
| `metrics.pod.annotations` | Custom annotations for the Pod, when `metrics.enabled` is set to `true` | `{}` |
| `envs.clusterAssetGroup.relistInterval` | The period of time after which the controller refreshes the status of a ClusterAssetGroup CR | `5m` |
| `envs.assetGroup.relistInterval` | The period of time after which the controller refreshes the status of an AssetGroup CR | `5m` |
| `envs.clusterBucket.relistInterval` | The period of time after which the controller refreshes the status of a ClusterBucket CR | `30s` |
| `envs.clusterBucket.maxConcurrentReconciles` | The maximum number of ClusterBucket reconciles that can run in parallel | `1` |
| `envs.clusterBucket.region` | The location of the region in which the controller creates a ClusterBucket CR. If the field is empty, the controller creates the bucket under the default location. | `us-east-1` |
| `envs.bucket.relistInterval` | The period of time after which the controller refreshes the status of a Bucket CR | `30s` |
| `envs.bucket.maxConcurrentReconciles` | The maximum number of Bucket reconciles that can run in parallel | `1` |
| `envs.bucket.region` | The location of the region in which the controller creates a Bucket CR. If the field is empty, the controller creates the bucket under the default location. | `us-east-1` |
| `envs.clusterAsset.relistInterval` | The period of time after which the controller refreshes the status of a ClusterAsset CR | `30s` |
| `envs.clusterAsset.maxConcurrentReconciles` | The maximum number of ClusterAsset reconciles that can run in parallel | `1` |
| `envs.asset.relistInterval` | The period of time after which the controller refreshes the status of a Asset CR | `30s` |
| `envs.asset.maxConcurrentReconciles` | The maximum number of Asset reconciles that can run in parallel | `1` |
| `envs.store.endpoint` | The address of the content storage server | `{{ include "rafter.fullname" . }}-minio.{{ .Release.Namespace }}.svc.cluster.local:9000` |
| `envs.store.externalEndpoint` | The external address of the content storage server | `nil` |
| `envs.store.accessKey` | The access key required to sign in to the content storage server | Value from `{{ include "rafter.fullname" . }}-minio` ConfigMap |
| `envs.store.secretKey` | The secret key required to sign in to the content storage server | Value from `{{ include "rafter.fullname" . }}-minio` ConfigMap |
| `envs.store.useSSL` | The HTTPS connection with the content storage server | `false` |
| `envs.store.uploadWorkers` | The number of workers used in parallel to upload files to the storage server | `10` |
| `envs.loader.verifySSL` | The variable that verifies the SSL certificate before downloading source files | `false` |
| `envs.loader.tempDir` | The path to the directory used to store data temporarily | `/tmp` |
| `envs.webhooks.validation.timeout` | The period of time after which validation is canceled | `1m` |
| `envs.webhooks.validation.workers` | The number of workers used in parallel to validate files | `10` |
| `envs.webhooks.mutation.timeout` | The period of time after which mutation is canceled | `1m` |
| `envs.webhooks.mutation.workers` | The number of workers used in parallel to mutate files | `10` |
| `envs.webhooks.metadata.timeout` | The period of time after which metadata extraction is canceled | `1m` |

Specify each parameter using the `--set key=value[,key=value]` argument for `helm install`. See this example:

``` bash
helm install --name rafter-release \
  --set serviceMonitor.create=true,serviceMonitor.name="rafter-service-monitor" \
    incubator/rafter-controller-manager
```

That command installs the release with the `rafter-service-monitor` name for the ServiceMonitor custom resource.

Alternatively, use the default values in [values.yaml](./values.yaml) or provide a YAML file while installing the chart to specify the values for configurable parameters. See this example:

``` bash
helm install --name rafter-release -f values.yaml incubator/rafter-controller-manager
```

### values.yaml as a template

The `values.yaml` for the Rafter Controller Manager chart serves as a template. Use such Helm variables as `.Release.*`, or `.Values.*`. See this example:

``` yaml
pod:
  annotations:
    sidecar.istio.io/inject: "{{ .Values.injectIstio }}"
    recreate: "{{ .Release.Time.Seconds }}"
``` 

### Change values for envs.* parameters

You can define values for all `envs.*` parameters as objects by providing the parameters as the inline `value` or the `valueFrom` object. See the following example:

``` yaml
envs:
  clusterAssetGroup:
    relistInterval: 
      value: 5m
  assetGroup:
    valueFrom:
      configMapKeyRef:
        name: rafter-config
        key: RAFTER_ASSET_GROUP_RELIST_INTERVALL
```